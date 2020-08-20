import cdk = require('@aws-cdk/core');
import ec2 = require('@aws-cdk/aws-ec2');
import s3 = require('@aws-cdk/aws-s3');
import sqs = require('@aws-cdk/aws-sqs');
import s3n = require('@aws-cdk/aws-s3-notifications');
import lambda = require('@aws-cdk/aws-lambda');
import { RetentionDays } from '@aws-cdk/aws-logs';

interface VpcProps extends cdk.StackProps {
	vpc: ec2.IVpc
	bucket: s3.Bucket
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
		});

		bucket.addEventNotification(s3.EventType.OBJECT_CREATED, new s3n.SqsDestination(q))

		const importFunc = new lambda.Function(this, 'F', {
			code: lambda.Code.asset('services/import'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: cdk.Duration.seconds(180),
			environment: {
				CSV_BUCKET_NAME: bucket.bucketName,
			},
			logRetention: RetentionDays.ONE_MONTH,
			tracing: lambda.Tracing.ACTIVE,
		})
		bucket.grantRead(importFunc)
	}
}
