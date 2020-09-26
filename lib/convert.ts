import cdk = require('@aws-cdk/core');
import ec2 = require('@aws-cdk/aws-ec2');
import s3 = require('@aws-cdk/aws-s3');
import sqs = require('@aws-cdk/aws-sqs');
import lambda = require('@aws-cdk/aws-lambda');
import { RetentionDays } from '@aws-cdk/aws-logs';
import s3n = require('@aws-cdk/aws-s3-notifications');
import { SqsEventSource } from '@aws-cdk/aws-lambda-event-sources';
import * as efs from '@aws-cdk/aws-efs';
import { Vpc } from './vpc';
import { RemovalPolicy } from '@aws-cdk/core';

interface VpcProps extends cdk.StackProps {
	vpc: ec2.IVpc;
	csvBucket: s3.Bucket
	uploadBucket: s3.Bucket
}

export class Convert extends cdk.Construct {
	public readonly lambda: lambda.Function;
	public readonly queue: sqs.Queue;
	public readonly DLQ: sqs.Queue;
	public readonly efs: efs.FileSystem;

	constructor(scope: cdk.Construct, id: string, props: VpcProps) {
		super(scope, id)

		this.DLQ = new sqs.Queue(this, 'DeadLetterQueue', {
			retentionPeriod: cdk.Duration.days(14),
		});

		const q = new sqs.Queue(this, 'ProcessingQueue', {
			deadLetterQueue: {
				queue: this.DLQ,
				maxReceiveCount: 1, //NOTE: I want failed messages to directly land in dlq during dev
			},
			visibilityTimeout: cdk.Duration.minutes(15)
		});
		this.queue = q

		const convertLambda = new lambda.Function(this, 'F', {
			code: lambda.Code.fromAsset('services/convert'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: cdk.Duration.minutes(3),
			environment: {
				CSV_BUCKET_NAME: props.csvBucket.bucketName,
			},
			reservedConcurrentExecutions: 50, 
			logRetention: RetentionDays.ONE_WEEK,
			tracing: lambda.Tracing.ACTIVE,
		})
		this.lambda = convertLambda
		convertLambda.addEventSource(new SqsEventSource(q))

		props.uploadBucket.addEventNotification(s3.EventType.OBJECT_CREATED, new s3n.SqsDestination(q))
		props.uploadBucket.grantRead(convertLambda)
		props.csvBucket.grantWrite(convertLambda)
	}
}
