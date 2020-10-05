import cdk = require('@aws-cdk/core');
import ec2 = require('@aws-cdk/aws-ec2');
import s3 = require('@aws-cdk/aws-s3');
import s3n = require('@aws-cdk/aws-s3-notifications');
import sqs = require('@aws-cdk/aws-sqs');
import lambda = require('@aws-cdk/aws-lambda');
import { RetentionDays } from '@aws-cdk/aws-logs';
import { SqsEventSource } from '@aws-cdk/aws-lambda-event-sources';
import * as secretsmanager from '@aws-cdk/aws-secretsmanager';
import * as sns from '@aws-cdk/aws-sns';

interface VpcProps extends cdk.StackProps {
	vpc: ec2.IVpc
	csvBucket: s3.Bucket
	dbSecGrp: ec2.SecurityGroup
	dbSecret : secretsmanager.ISecret
	dbEndpoint: string
	summaryTopic: sns.Topic
}

export class Import extends cdk.Construct {
	public readonly importLambda: lambda.Function;
	public readonly importQueue: sqs.Queue;
	public readonly ImportDLQ: sqs.Queue;

	constructor(scope: cdk.Construct, id: string, props: VpcProps) {
		super(scope, id)

		this.ImportDLQ= new sqs.Queue(this, 'DeadLetterQueue', {
			retentionPeriod: cdk.Duration.days(14),
		});

		//NOTE: sometimes the db import fails, thats why the maxReceiveCount is so high
		//		the error fixes itself on the next try or two
		this.importQueue= new sqs.Queue(this, 'ProcessingQueue', {
			deadLetterQueue: {
				queue: this.ImportDLQ,
				maxReceiveCount: 10,
			},
			visibilityTimeout: cdk.Duration.minutes(10),
		});

		props.csvBucket.addEventNotification(s3.EventType.OBJECT_CREATED, new s3n.SqsDestination(this.importQueue))

		this.importLambda = new lambda.Function(this, 'Lambda', {
			code: lambda.Code.fromAsset('services/upload/import'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: cdk.Duration.seconds(180),
			environment: {
				SECRET_ARN: props.dbSecret.secretArn,
				DB_ENDPOINT: props.dbEndpoint,
				SUMMARY_TOPIC_NAME: props.summaryTopic.topicArn,
			},
			reservedConcurrentExecutions: 1, 
			logRetention: RetentionDays.ONE_WEEK,
			tracing: lambda.Tracing.ACTIVE,
			vpc: props.vpc,
			securityGroups: [props.dbSecGrp],
			// onSuccess: new destinations.LambdaDestination(this.summaryLambda),
			//NOTE: SQS invokes lambda synchronously and thus lambda Destinations
			//		don't work. I have to post the message to the topic in the code
			//		https://docs.aws.amazon.com/sdk-for-go/api/service/lambda/#Lambda.Invoke
			//		second example
		})
		this.importLambda.addEventSource(new SqsEventSource(this.importQueue, {
			batchSize: 1,
		}))

		props.dbSecret?.grantRead(this.importLambda)

		props.summaryTopic.grantPublish(this.importLambda)
	}
}
