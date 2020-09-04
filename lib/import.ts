import cdk = require('@aws-cdk/core');
import ec2 = require('@aws-cdk/aws-ec2');
import s3 = require('@aws-cdk/aws-s3');
import sqs = require('@aws-cdk/aws-sqs');
import s3n = require('@aws-cdk/aws-s3-notifications');
import lambda = require('@aws-cdk/aws-lambda');
import { RetentionDays } from '@aws-cdk/aws-logs';
import { SqsEventSource } from '@aws-cdk/aws-lambda-event-sources';
import * as secretsmanager from '@aws-cdk/aws-secretsmanager';

interface VpcProps extends cdk.StackProps {
	vpc: ec2.IVpc
	bucket: s3.Bucket
	securityGroup: ec2.SecurityGroup
	secret : secretsmanager.ISecret
}

export class Import extends cdk.Construct {
	constructor(scope: cdk.Construct, id: string, props: VpcProps) {
		super(scope, id)

		const bucket = props.bucket

		const dlq = new sqs.Queue(this, 'DeadLetterQueue', {
			retentionPeriod: cdk.Duration.days(14),
		});

		const q = new sqs.Queue(this, 'ProcessingQueue', {
			deadLetterQueue: {
				queue: dlq,
				maxReceiveCount: 3,
			},
			visibilityTimeout: cdk.Duration.seconds(30)
		});

		bucket.addEventNotification(s3.EventType.OBJECT_CREATED, new s3n.SqsDestination(q))

		const importFunc = new lambda.Function(this, 'F', {
			code: lambda.Code.asset('services/import'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: cdk.Duration.seconds(5),
			environment: {
				SECRET_ARN: props.secret.secretArn,
			},
			//only execute one lambda at a time, to not overwhelm the db with imports
			reservedConcurrentExecutions: 1, 
			logRetention: RetentionDays.ONE_WEEK,
			tracing: lambda.Tracing.ACTIVE,
			vpc: props.vpc,
			securityGroups: [props.securityGroup],
		})
		importFunc.addEventSource(new SqsEventSource(q))
		props.secret?.grantRead(importFunc)
	}
}
