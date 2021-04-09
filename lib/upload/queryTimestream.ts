import cdk = require('@aws-cdk/core');
import sqs = require('@aws-cdk/aws-sqs');
import lambda = require('@aws-cdk/aws-lambda');
import { RetentionDays } from '@aws-cdk/aws-logs';
import * as dynamodb from '@aws-cdk/aws-dynamodb';
import * as iam from "@aws-cdk/aws-iam"
import * as sns from "@aws-cdk/aws-sns";
import * as destinations from '@aws-cdk/aws-lambda-destinations';

interface Props extends cdk.StackProps {
	dynamoDB: dynamodb.Table,
	timestreamArn: string
	codePath: string
	lambdaDescription: string
	envVars: {[key: string]: string}
}

export class QueryTimestream extends cdk.Construct {
	public readonly lambda: lambda.Function;
	public readonly topicDLQ: sqs.Queue;
	public readonly lambdaDLQ: sqs.Queue;
	public readonly topic: sns.Topic;

	constructor(scope: cdk.Construct, id: string, props: Props) {
		super(scope, id)

		this.topicDLQ = new sqs.Queue(this, 'TopicDLQ')
		this.lambdaDLQ = new sqs.Queue(this, 'LambdaDLQ')

		//the message to the topic is send inside the lambda via the SDK
		//the topic in return is subscribed to by the insert summary lambdas
		this.topic = new sns.Topic(this, 'Topic')

		this.lambda = new lambda.Function(this, 'Lambda', {
			description: props.lambdaDescription,
			code: lambda.Code.fromAsset(props.codePath),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 128,
			timeout: cdk.Duration.seconds(40),
			environment: {
				TOPIC_ARN: this.topic.topicArn,
				LOCAL: "false",
				...props.envVars,
			},
			reservedConcurrentExecutions: 75,
			logRetention: RetentionDays.ONE_MONTH,
			tracing: lambda.Tracing.ACTIVE,
			retryAttempts: 2, //default
			onFailure: new destinations.SqsDestination(this.lambdaDLQ),
			//Fails will be retried twice without landing in the DLQ, if the last retry also fails the message
			//lands in the DLQ
			//TODO: on success destination doesn't work with xray DON'T USE until fixed
		})

		this.topic.grantPublish(this.lambda)

        //extract to function in timestream construct
		this.lambda.addToRolePolicy(new iam.PolicyStatement({
			actions: [
				"timestream:DescribeEndpoints",
				"timestream:SelectValues",
			],
			resources: ["*"],
			effect: iam.Effect.ALLOW,
		}))

		this.lambda.addToRolePolicy(new iam.PolicyStatement({
			actions: [
				"timestream:Select",
				"timestream:ListMeasures"
			],
			resources: [props.timestreamArn],
			effect: iam.Effect.ALLOW,
		}))
	}
}
