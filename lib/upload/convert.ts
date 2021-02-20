import cdk = require('@aws-cdk/core');
import ec2 = require('@aws-cdk/aws-ec2');
import s3 = require('@aws-cdk/aws-s3');
import sqs = require('@aws-cdk/aws-sqs');
import lambda = require('@aws-cdk/aws-lambda');
import { RetentionDays } from '@aws-cdk/aws-logs';
import s3n = require('@aws-cdk/aws-s3-notifications');
import { SqsEventSource } from '@aws-cdk/aws-lambda-event-sources';
import * as iam from "@aws-cdk/aws-iam"

interface Props extends cdk.StackProps {
	uploadBucket: s3.Bucket
	timestreamArn: string
	summaryQueue: sqs.Queue
}

export class Convert extends cdk.Construct {
	public readonly lambda: lambda.Function;
	public readonly queue: sqs.Queue;
	public readonly DLQ: sqs.Queue;

	constructor(scope: cdk.Construct, id: string, props: Props) {
		super(scope, id)

		this.DLQ = new sqs.Queue(this, 'DeadLetterQueue', {
			retentionPeriod: cdk.Duration.days(14),
		});

		this.queue = new sqs.Queue(this, 'ProcessingQueue', {
			deadLetterQueue: {
				queue: this.DLQ,
				maxReceiveCount: 1, //NOTE: I want failed messages to directly land in dlq during dev
			},
			visibilityTimeout: cdk.Duration.minutes(10)
		});

		this.lambda = new lambda.Function(this, 'Lambda', {
			code: lambda.Code.fromAsset('services/upload/convert'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 10240,
			timeout: cdk.Duration.minutes(2),
			environment: {
				SUMMARY_QUEUE: props.summaryQueue.queueName,
			},
			reservedConcurrentExecutions: 50, 
			logRetention: RetentionDays.ONE_WEEK,
			tracing: lambda.Tracing.ACTIVE,
		})
		this.lambda.addEventSource(new SqsEventSource(this.queue, {
			batchSize: 1,
			//the should be able to handle multiple events at once.
			//but for now limiting it to one element per invocation makes it easier to reason about
			//the expected duration etc.
			//I'll optimize this later, potentially to save invocation costs
			//also I'm not sure how the memory is garbage collected between different elements in the batch
		}))

		props.summaryQueue.grantSendMessages(this.lambda)

		props.uploadBucket.addEventNotification(s3.EventType.OBJECT_CREATED, new s3n.SqsDestination(this.queue))
		props.uploadBucket.grantRead(this.lambda)

		this.lambda.addToRolePolicy(new iam.PolicyStatement({
			actions: [
				"timestream:DescribeEndpoints",
			],
			resources: ["*"], 
			effect: iam.Effect.ALLOW,
		}))

		this.lambda.addToRolePolicy(new iam.PolicyStatement({
			actions: [
				"timestream:WriteRecords",
			],
			resources: [props.timestreamArn], 
			effect: iam.Effect.ALLOW,
		}))
	}
}
