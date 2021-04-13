import cdk = require('@aws-cdk/core');
import s3 = require('@aws-cdk/aws-s3');
import sqs = require('@aws-cdk/aws-sqs');
import * as sns from '@aws-cdk/aws-sns';
import * as subs from '@aws-cdk/aws-sns-subscriptions';
import lambda = require('@aws-cdk/aws-lambda');
import { RetentionDays } from '@aws-cdk/aws-logs';
import s3n = require('@aws-cdk/aws-s3-notifications');
import { SqsEventSource } from '@aws-cdk/aws-lambda-event-sources';
import * as iam from "@aws-cdk/aws-iam"

interface Props extends cdk.StackProps {
	uploadBucket: s3.Bucket
	timestreamArn: string
	queryTimestreamLambdas: lambda.Function[] //to get summary lambda array
	envVars: {[key: string]: string}
}

export class Convert extends cdk.Construct {
	public readonly lambda: lambda.Function;
	public readonly queue: sqs.Queue;
	public readonly DLQ: sqs.Queue;
	public readonly topic: sns.Topic;

	constructor(scope: cdk.Construct, id: string, props: Props) {
		super(scope, id)

		this.DLQ = new sqs.Queue(this, 'DeadLetterQueue', {
			retentionPeriod: cdk.Duration.days(14),
		});

		this.queue = new sqs.Queue(this, 'ProcessingQueue', {
			deadLetterQueue: {
				queue: this.DLQ,
				//maxReceiveCount: 3,
				maxReceiveCount: 1, //no need during dev
			},
			visibilityTimeout: cdk.Duration.minutes(18) //6x lambda duration
		});

        const topic = new sns.Topic(this, 'Topic', {})
        this.topic = topic

		props.queryTimestreamLambdas.forEach(function(l){
			topic.addSubscription(new subs.LambdaSubscription(l))
		})

		this.lambda = new lambda.Function(this, 'Lambda', {
			description: "takes combatlog file and uploads it to amazon timestream",
			code: lambda.Code.fromAsset('services/upload/convert'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			// memorySize: 3584, //exactly 2 core
			memorySize: 1792, //exactly 1 core
			timeout: cdk.Duration.seconds(30),
			environment: {
				TOPIC_ARN: topic.topicArn,
				LOCAL: "false",
				...props.envVars,
			},
			reservedConcurrentExecutions: 60,
			logRetention: RetentionDays.ONE_WEEK,
			tracing: lambda.Tracing.ACTIVE,
			retryAttempts: 0, 	//0 in dev, but it has sqs as target, afaik this is only for async.
								// sqs invokes would be retried via the sqs maxReceiveCount
			/* the source is sqs, which invokes the lambda synchronously, ergo no onFailure or onSuccess =(
            onFailure
			onSuccess:
			 */
		})
		this.lambda.addEventSource(new SqsEventSource(this.queue, {
			batchSize: 1,
			//leave at one, simplifies the code and invocation costs of lambda are very likely not gonna matter
		}))

		topic.grantPublish(this.lambda)

		//only trigger convert lambda if file end on one of these suffixes
        //in theory files with a wrong ending could linger in the bucket forever without being processed
		//but the presign lambda refuses uploads if the ending is not one of the mentioned
		props.uploadBucket.addEventNotification(s3.EventType.OBJECT_CREATED, new s3n.SqsDestination(this.queue), {
			suffix: ".txt",
			prefix: "upload/"
		})
		props.uploadBucket.addEventNotification(s3.EventType.OBJECT_CREATED, new s3n.SqsDestination(this.queue), {
			suffix: ".txt.gz",
			prefix: "upload/"
		})
		props.uploadBucket.addEventNotification(s3.EventType.OBJECT_CREATED, new s3n.SqsDestination(this.queue), {
			suffix: ".zip",
			prefix: "upload/"
		})
		props.uploadBucket.grantReadWrite(this.lambda)

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
