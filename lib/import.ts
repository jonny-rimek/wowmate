import cdk = require('@aws-cdk/core');
import ecs = require('@aws-cdk/aws-ecs');
import ecsPatterns = require('@aws-cdk/aws-ecs-patterns');
import ec2 = require('@aws-cdk/aws-ec2');
import s3 = require('@aws-cdk/aws-s3');
import sqs = require('@aws-cdk/aws-sqs');
import s3n = require('@aws-cdk/aws-s3-notifications');

interface VpcProps extends cdk.StackProps {
	vpc: ec2.IVpc
	bucket: s3.Bucket
}

export class Import extends cdk.Construct {
	constructor(scope: cdk.Construct, id: string, props: VpcProps) {
		super(scope, id)

		const b = props.bucket

		const dlq = new sqs.Queue(this, 'DeadLetterQueue', {
			retentionPeriod: cdk.Duration.days(14),
		});

		const q = new sqs.Queue(this, 'ProcessingQueue', {
			deadLetterQueue: {
				queue: dlq,
				maxReceiveCount: 3,
			},
		});

		// CSV_BUCKET_NAME: csvBucket.bucketName,

		// uploadBucket.addEventNotification(s3.EventType.OBJECT_CREATED, new s3n.SqsDestination(queueFargate.sqsQueue))
		// csvBucket.grantRead(queueFargate.service.taskDefinition.taskRole)
	}
}
