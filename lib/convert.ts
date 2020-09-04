import cdk = require('@aws-cdk/core');
import ecs = require('@aws-cdk/aws-ecs');
import ecsPatterns = require('@aws-cdk/aws-ecs-patterns');
import ec2 = require('@aws-cdk/aws-ec2');
import s3 = require('@aws-cdk/aws-s3');
import sqs = require('@aws-cdk/aws-sqs');
import lambda = require('@aws-cdk/aws-lambda');
import { RetentionDays } from '@aws-cdk/aws-logs';
import s3n = require('@aws-cdk/aws-s3-notifications');
import { SqsEventSource } from '@aws-cdk/aws-lambda-event-sources';

interface VpcProps extends cdk.StackProps {
	vpc: ec2.IVpc;
	csvBucket: s3.Bucket
	uploadBucket: s3.Bucket
}

export class Convert extends cdk.Construct {
	constructor(scope: cdk.Construct, id: string, props: VpcProps) {
		super(scope, id)

		const dlq = new sqs.Queue(this, 'DeadLetterQueue', {
			retentionPeriod: cdk.Duration.days(14),
		});

		const q = new sqs.Queue(this, 'ProcessingQueue', {
			deadLetterQueue: {
				queue: dlq,
				maxReceiveCount: 3,
			},
			visibilityTimeout: cdk.Duration.minutes(2)
		});

		const convertFunc = new lambda.Function(this, 'F', {
			code: lambda.Code.asset('services/convert'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: cdk.Duration.seconds(60),
			environment: {
				CSV_BUCKET_NAME: props.csvBucket.bucketName,
			},
			reservedConcurrentExecutions: 10, 
			logRetention: RetentionDays.ONE_WEEK,
			tracing: lambda.Tracing.ACTIVE,
			//NOTE: not in VPC by design, because I don't have an S3 endpoint and it would incur
			//		additional charges
		})
		convertFunc.addEventSource(new SqsEventSource(q))

		props.uploadBucket.addEventNotification(s3.EventType.OBJECT_CREATED, new s3n.SqsDestination(q))
		props.uploadBucket.grantRead(convertFunc)
		props.csvBucket.grantWrite(convertFunc)
	}
}
