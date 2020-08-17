import cdk = require('@aws-cdk/core');
import ecs = require('@aws-cdk/aws-ecs');
import ecsPatterns = require('@aws-cdk/aws-ecs-patterns');
import ec2 = require('@aws-cdk/aws-ec2');
import s3 = require('@aws-cdk/aws-s3');
import sqs = require('@aws-cdk/aws-sqs');
import s3n = require('@aws-cdk/aws-s3-notifications');

interface VpcProps extends cdk.StackProps {
	vpc: ec2.IVpc;
}

export class Converter extends cdk.Construct {
	constructor(scope: cdk.Construct, id: string, props: VpcProps) {
		super(scope, id)

		const dlq = new sqs.Queue(this, 'EcsProcessingDeadLetterQueue', {
			retentionPeriod: cdk.Duration.days(14),
		});

		const q = new sqs.Queue(this, 'EcsProcessingQueue', {
			deadLetterQueue: {
				queue: dlq,
				maxReceiveCount: 3,
			},
		});

		const queueFargate = new ecsPatterns.QueueProcessingFargateService(this, 'Service', {
			queue: q,
			vpc: props.vpc,
			memoryLimitMiB: 512,
			cpu: 256,
			image: ecs.ContainerImage.fromAsset('services/converter'),
			platformVersion: ecs.FargatePlatformVersion.VERSION1_4,
			desiredTaskCount: 1,
			// this is created by default
			// environment: {
			// 	QUEUE_URL: q.queueUrl
			// },
		});

		const uploadBucket = new s3.Bucket(this, 'Upload', {
			removalPolicy: cdk.RemovalPolicy.DESTROY,
			blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL,
		})

		uploadBucket.addEventNotification(s3.EventType.OBJECT_CREATED, new s3n.SqsDestination(queueFargate.sqsQueue))
		uploadBucket.grantRead(queueFargate.service.taskDefinition.taskRole)

		//add read from upload bucket
		// queueFargate.taskDefinition.addToTaskRolePolicy
	}
}
