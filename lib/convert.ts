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
				maxReceiveCount: 1, //NOTE: I want failed messages to directly land in dlq
			},
			visibilityTimeout: cdk.Duration.minutes(20)
		});
		this.queue = q

		this.efs = new efs.FileSystem(this, 'Efs', {
			vpc: props.vpc,
			encrypted: false,
			// performanceMode: efs.PerformanceMode.GENERAL_PURPOSE,
			performanceMode: efs.PerformanceMode.MAX_IO,
			throughputMode: efs.ThroughputMode.BURSTING,
			//TODO: remove in prod
			removalPolicy: RemovalPolicy.DESTROY,
		})

		const accessPoint = this.efs.addAccessPoint('ConvertAccessPoint', {
			path: '/convert',
			createAcl: {
				ownerGid: '1001',
				ownerUid: '1001',
				permissions: '755',
			},
			posixUser: {
				uid: '1001',
				gid: '1001',
			},
		})

		const convertLambda = new lambda.Function(this, 'F', {
			code: lambda.Code.fromAsset('services/convert'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: cdk.Duration.minutes(15),
			environment: {
				CSV_BUCKET_NAME: props.csvBucket.bucketName,
			},
			reservedConcurrentExecutions: 50, 
			logRetention: RetentionDays.ONE_WEEK,
			tracing: lambda.Tracing.ACTIVE,
			filesystem: lambda.FileSystem.fromEfsAccessPoint(accessPoint, '/mnt/efs'),
			vpc: props.vpc,
			//NOTE: not in VPC by design, because I don't have an S3 endpoint and it would incur
			//		additional charges
			//		if I endup using EFS I need to add it to the VPC tho
		})
		this.lambda = convertLambda
		convertLambda.addEventSource(new SqsEventSource(q))

		props.uploadBucket.addEventNotification(s3.EventType.OBJECT_CREATED, new s3n.SqsDestination(q))
		props.uploadBucket.grantRead(convertLambda)
		props.csvBucket.grantWrite(convertLambda)
	}
}
