import cdk = require('@aws-cdk/core');
import sqs = require('@aws-cdk/aws-sqs');
import lambda = require('@aws-cdk/aws-lambda');
import { RetentionDays } from '@aws-cdk/aws-logs';
import * as dynamodb from '@aws-cdk/aws-dynamodb';
import * as sns from "@aws-cdk/aws-sns";
import * as subs from "@aws-cdk/aws-sns-subscriptions";
import * as destinations from '@aws-cdk/aws-lambda-destinations';

interface Props extends cdk.StackProps {
	dynamoDB: dynamodb.Table,
	topic: sns.Topic;
	topicDLQ: sqs.Queue;
	codePath: string
	lambdaDescription: string
	envVars: {[key: string]: string}
}

export class InsertResult extends cdk.Construct {
	public readonly lambda: lambda.Function;
	public readonly lambdaDLQ: sqs.Queue;

	constructor(scope: cdk.Construct, id: string, props: Props) {
		super(scope, id)

		this.lambdaDLQ = new sqs.Queue(this, 'LambdaDLQ')

		this.lambda = new lambda.Function(this, 'Lambda', {
			description: props.description,
			code: lambda.Code.fromAsset(props.codePath),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 128,
			timeout: cdk.Duration.seconds(2),
			environment: {
				DYNAMODB_TABLE_NAME: props.dynamoDB.tableName,
				LOCAL: "false",
				...props.envVars
			},
			reservedConcurrentExecutions: 30,
			logRetention: RetentionDays.ONE_MONTH,
			tracing: lambda.Tracing.ACTIVE,
			retryAttempts: 2, //default
			onFailure: new destinations.SqsDestination(this.lambdaDLQ),
			//Fails will be retried twice without landing in the DLQ, if the last retry also fails the message
			//lands in the DLQ
			//TODO: on success destination doesn't work with xray DON'T USE until fixed
		})
        //temp queue to get message content easily
		//this is not the whole event one needs to invoke it locally only the SNS part
		const q = new sqs.Queue(this, 'Q')
		props.topic.addSubscription(new subs.SqsSubscription(q))

		props.topic.addSubscription(new subs.LambdaSubscription(this.lambda, {
			deadLetterQueue: props.topicDLQ,
		}))

		//IMPROVE: this also adds permission to delete data, which we don't need
		props.dynamoDB.grantWriteData(this.lambda)
	}
}
