import cdk = require('@aws-cdk/core');
import ec2 = require('@aws-cdk/aws-ec2');
import s3 = require('@aws-cdk/aws-s3');
import sqs = require('@aws-cdk/aws-sqs');
import lambda = require('@aws-cdk/aws-lambda');
import { RetentionDays } from '@aws-cdk/aws-logs';
import * as destinations from '@aws-cdk/aws-lambda-destinations';
import * as secretsmanager from '@aws-cdk/aws-secretsmanager';
import * as sns from '@aws-cdk/aws-sns';
import * as subs from '@aws-cdk/aws-sns-subscriptions';

interface VpcProps extends cdk.StackProps {
	// vpc: ec2.IVpc
	// csvBucket: s3.Bucket
	// dbSecret : secretsmanager.ISecret
	// dbEndpoint: string
}

export class Summary extends cdk.Construct {
	public readonly lambda: lambda.Function;
	public readonly lambdaDLQ: sqs.Queue;
	public readonly queue: sqs.Queue;
	public readonly queueDLQ: sqs.Queue;
	public readonly topic: sns.Topic

	constructor(scope: cdk.Construct, id: string, props: VpcProps) {
		super(scope, id)

		this.queueDLQ = new sqs.Queue(this, 'QueueDLQ', {
			retentionPeriod: cdk.Duration.days(14),
		});

		this.queue = new sqs.Queue(this, 'Queue', {
			deadLetterQueue: {
				queue: this.queueDLQ,
				maxReceiveCount: 1, //NOTE: I want failed messages to directly land in dlq during dev
			},
			visibilityTimeout: cdk.Duration.minutes(10)
		});

		this.lambdaDLQ = new sqs.Queue(this, 'LambdaDLQ', {
			retentionPeriod: cdk.Duration.days(14)
		})

		this.lambda = new lambda.Function(this, 'Lambda', {
			code: lambda.Code.fromAsset('services/upload/summary'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: cdk.Duration.seconds(60),
			// environment: {},
			reservedConcurrentExecutions: 10, 
			logRetention: RetentionDays.ONE_WEEK,
			tracing: lambda.Tracing.ACTIVE,
			onFailure: new destinations.SqsDestination(this.lambdaDLQ)
		})
	}
}
