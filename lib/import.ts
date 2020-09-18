import cdk = require('@aws-cdk/core');
import ec2 = require('@aws-cdk/aws-ec2');
import s3 = require('@aws-cdk/aws-s3');
import sqs = require('@aws-cdk/aws-sqs');
import s3n = require('@aws-cdk/aws-s3-notifications');
import lambda = require('@aws-cdk/aws-lambda');
import { RetentionDays } from '@aws-cdk/aws-logs';
import { SqsEventSource } from '@aws-cdk/aws-lambda-event-sources';
import * as destinations from '@aws-cdk/aws-lambda-destinations';
import * as secretsmanager from '@aws-cdk/aws-secretsmanager';
import rds = require('@aws-cdk/aws-rds');
import { DatabaseProxy } from '@aws-cdk/aws-rds';

interface VpcProps extends cdk.StackProps {
	vpc: ec2.IVpc
	bucket: s3.Bucket
	securityGroup: ec2.SecurityGroup
	secret : secretsmanager.ISecret
	dbEndpoint: string
}

export class Import extends cdk.Construct {
	public readonly importLambda: lambda.Function;
	public readonly summaryLambda: lambda.Function;
	public readonly queue: sqs.Queue;
	public readonly dlq: sqs.Queue;
	public readonly summaryDlq: sqs.Queue;

	constructor(scope: cdk.Construct, id: string, props: VpcProps) {
		super(scope, id)

		const bucket = props.bucket

		this.summaryDlq = new sqs.Queue(this, 'SummaryDeadLetterQueue', {
			retentionPeriod: cdk.Duration.days(14)
		})

		this.dlq = new sqs.Queue(this, 'DeadLetterQueue', {
			retentionPeriod: cdk.Duration.days(14),
		});

		//NOTE: sometimes the db import fails, thats why the maxReceiveCount is so high
		//		the error fixes itself on the next try or two
		this.queue = new sqs.Queue(this, 'ProcessingQueue', {
			deadLetterQueue: {
				queue: this.dlq,
				maxReceiveCount: 10,
			},
			visibilityTimeout: cdk.Duration.minutes(10)
		});

		bucket.addEventNotification(s3.EventType.OBJECT_CREATED, new s3n.SqsDestination(this.queue))

		this.summaryLambda = new lambda.Function(this, 'SummaryLambda', {
			code: lambda.Code.fromAsset('services/summary'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: cdk.Duration.seconds(60),
			environment: {
				SECRET_ARN: props.secret.secretArn,
				DB_ENDPOINT: props.dbEndpoint,
			},
			reservedConcurrentExecutions: 10, 
			logRetention: RetentionDays.ONE_WEEK,
			tracing: lambda.Tracing.ACTIVE,
			vpc: props.vpc,
			securityGroups: [props.securityGroup],
			onFailure: new destinations.SqsDestination(this.summaryDlq)
		})
		props.secret?.grantRead(this.summaryLambda)
		
		this.importLambda = new lambda.Function(this, 'Lambda', {
			code: lambda.Code.fromAsset('services/import'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: cdk.Duration.seconds(180),
			environment: {
				SECRET_ARN: props.secret.secretArn,
				DB_ENDPOINT: props.dbEndpoint,
				SUMMARY_LAMBDA_NAME: this.summaryLambda.functionName,
			},
			reservedConcurrentExecutions: 1, 
			logRetention: RetentionDays.ONE_WEEK,
			tracing: lambda.Tracing.ACTIVE,
			vpc: props.vpc,
			securityGroups: [props.securityGroup],
			// onSuccess: new destinations.LambdaDestination(this.summaryLambda),
			//NOTE: SQS invokes lambda synchronously and thus lambda Destinations
			//		don't work. I have to call the summary lambda in the code
			//		https://docs.aws.amazon.com/sdk-for-go/api/service/lambda/#Lambda.Invoke
			//		second example
		})
		this.summaryLambda.grantInvoke(this.importLambda)

		this.importLambda.addEventSource(new SqsEventSource(this.queue))
		props.secret?.grantRead(this.importLambda)
	}
}
