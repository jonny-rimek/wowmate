import cdk = require('@aws-cdk/core');
import s3 = require('@aws-cdk/aws-s3');
import sqs = require('@aws-cdk/aws-sqs');
import lambda = require('@aws-cdk/aws-lambda');
import s3n = require('@aws-cdk/aws-s3-notifications');
import * as sns from '@aws-cdk/aws-sns';
import * as subs from '@aws-cdk/aws-sns-subscriptions';
import {RetentionDays} from '@aws-cdk/aws-logs';
import {SqsEventSource} from '@aws-cdk/aws-lambda-event-sources';
import * as iam from "@aws-cdk/aws-iam"
import {Effect} from "@aws-cdk/aws-iam"
import * as kms from '@aws-cdk/aws-kms';
import * as dynamodb from '@aws-cdk/aws-dynamodb';

interface Props extends cdk.StackProps {
	uploadBucket: s3.Bucket
    dynamodb: dynamodb.Table
	timestreamArn: string
	queryTimestreamLambdas: lambda.Function[] //to get summary lambda array
	envVars: {[key: string]: string}
	key: kms.Key
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
			encryption: sqs.QueueEncryption.KMS_MANAGED,
		});

		this.queue = new sqs.Queue(this, 'ProcessingQueue', {
			deadLetterQueue: {
				queue: this.DLQ,
				//maxReceiveCount: 3,
				maxReceiveCount: 1, //no need during dev
			},
			visibilityTimeout: cdk.Duration.seconds(150*6), //6x lambda duration, it's an aws best practice
			encryption: sqs.QueueEncryption.KMS,
			encryptionMasterKey: props.key,
		});

        const topic = new sns.Topic(this, 'Topic', {
			masterKey: props.key,
			// publishing encrypted messages to SNS doesn't work from SAM+CDK, I suspect that SAM uses an incomplete name
			// CDK, auto generates a name, but those names aren't used locally
			// e.g. wm-dev-DynamoDBtableF8E87752-HSV525WR7KN3 is the name of the ddb in the cloud
			// locally it the name it knows is wm-dev-DynamoDBtableF8E87752 the last bit is missing
            // the same is gonna be the problem for the KMS key, and I don't know how or if I can pass in the complete key
			// my solution is to skip the sns publishing locally
			// this will probably be fixed by sam for cdk which is in beta
		})
        this.topic = topic

		props.queryTimestreamLambdas.forEach(function(l){
			topic.addSubscription(new subs.LambdaSubscription(l))
		})

		this.lambda = new lambda.Function(this, 'Lambda', {
			description: "takes combatlog file and uploads it to amazon timestream",
			code: lambda.Code.fromAsset('dist/upload/convert'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			// with a file limit of 450MB uncompressed 1792 would be enough, but on repeated invokes
			// of the lambda the max memory increases.
			// e.g. first invoke 1535MB max memory, second 2238MB
			// I'm not sure what the reason is, I'll leave it at two cores for now!
			memorySize: 3584, //exactly 2 core
			// memorySize: 1792, //exactly 1 core
			timeout: cdk.Duration.seconds(150),
			// timestream write api has some sort of cold start, where at the beginning
			// it's super slow, that's why the max duration needs to be way higher than
			// the median duration
			// for my example log with cold start it has a duration of ~120sec and warm it is 25-40sec
			// synchronously making all the timestream api calls takes always around 120sec, 4-5x longer
			environment: {
				TOPIC_ARN: topic.topicArn,
				DYNAMODB_TABLE_NAME: props.dynamodb.tableName,
				...props.envVars,
			},
			reservedConcurrentExecutions: 150,
			logRetention: RetentionDays.ONE_WEEK,
			tracing: lambda.Tracing.ACTIVE,
			retryAttempts: 0, 	// it has sqs as target, this is only for async.
								// sqs invokes would be retried via the sqs maxReceiveCount

			/* the source is sqs, which invokes the lambda synchronously, ergo no onFailure or onSuccess =(
            onFailure
			onSuccess:
			 */
		})
		this.lambda.addEventSource(new SqsEventSource(this.queue, {
			batchSize: 1,
			//leave at one, simplifies the code and invocation costs of lambda are negligible compared to the rest
		}))

		props.dynamodb.grantReadWriteData(this.lambda)
        props.key.grantEncryptDecrypt(this.lambda)
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
			effect: iam.Effect.ALLOW,
			actions: [
				"timestream:DescribeEndpoints",
			],
			resources: ["*"], // resource needs to be *, doesn't work with the ARN
		}))

		this.lambda.addToRolePolicy(new iam.PolicyStatement({
			effect: iam.Effect.ALLOW,
			actions: [
				"timestream:WriteRecords",
			],
			resources: [props.timestreamArn], 
		}))
	}
}

// this is code to use lambda canary deployments to run integration tests during deployments and potentially
// rollback a deployment in case of errors. It doesn't work very well for long workflows and I don't want my
// integration tests to write data in my prod account, especially because you can't delete data in timestream
/*
const versionAlias = new lambda.Alias(this, 'Alias', {
    aliasName: "alias",
    version: this.lambda.currentVersion,
})

const preHook = new lambda.Function(this, 'LambdaPreHook', {
    description: "pre hook",
    code: lambda.Code.fromAsset('dist/upload/convert-pre-hook'),
    handler: 'main',
    runtime: lambda.Runtime.GO_1_X,
    memorySize: 128,
    timeout: cdk.Duration.minutes(1),
    environment: {
        FUNCTION_NAME: this.lambda.currentVersion.functionName,
    },
    reservedConcurrentExecutions: 5,
    logRetention: RetentionDays.ONE_WEEK,
})
// this.lambda.grantInvoke(preHook) // this doesn't work, I need to grant invoke to all functions :s
preHook.addToRolePolicy(new iam.PolicyStatement({
    actions: [
        "lambda:InvokeFunction",
    ],
    resources: ["*"],
    effect: iam.Effect.ALLOW,
}))

const application = new codedeploy.LambdaApplication(this, 'CodeDeployApplication')
new codedeploy.LambdaDeploymentGroup(this, 'CanaryDeployment', {
    application: application,
    alias: versionAlias,
    deploymentConfig: codedeploy.LambdaDeploymentConfig.ALL_AT_ONCE,
    preHook: preHook,
    autoRollback: {
        failedDeployment: true,
        stoppedDeployment: true,
        deploymentInAlarm: false,
    },
    ignorePollAlarmsFailure: false,
    // alarms:
    // autoRollback: codedeploy.A
    // postHook:
})

 */
