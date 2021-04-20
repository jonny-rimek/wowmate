import cdk = require('@aws-cdk/core');
import cloudwatch = require('@aws-cdk/aws-cloudwatch');
import { GraphWidget, IMetric, Metric, Row, TextWidget } from "@aws-cdk/aws-cloudwatch";
import sns = require('@aws-cdk/aws-sns');
import * as subscriptions from '@aws-cdk/aws-sns-subscriptions';
import lambda = require('@aws-cdk/aws-lambda');
import sqs = require('@aws-cdk/aws-sqs');
import s3 = require('@aws-cdk/aws-s3');
import { HttpApi } from '@aws-cdk/aws-apigatewayv2';
import * as dynamodb from '@aws-cdk/aws-dynamodb';
import * as kms from "@aws-cdk/aws-kms";

interface Props extends cdk.StackProps {
	convertLambda: lambda.Function
	convertQueue: sqs.Queue
	convertDLQ: sqs.Queue
	queryKeys: lambda.Function
	insertKeysToDynamoDB: lambda.Function
	insertKeysToTimestream: lambda.Function
	queryPlayerDamageDone: lambda.Function
	insertPlayerDamageDoneToDynamodb: lambda.Function
	queryKeysTopicDLQ: sqs.Queue
	queryPlayerDamageDoneTopicDLQ: sqs.Queue
	queryPlayerDamageDoneLambdaDLQ: sqs.Queue
	queryKeysLambdaDLQ: sqs.Queue
	insertKeysToDynamoDBLambdaDLQ: sqs.Queue
	insertKeysToTimestreamLambdaDLQ: sqs.Queue
	insertPlayerDamageDoneDynamoDBLambdaDLQ: sqs.Queue
	presignLambda: lambda.Function
	uploadBucket: s3.Bucket
	presignApi: HttpApi
	dynamoDB: dynamodb.Table
	convertTopic: sns.Topic
	queryKeysTopic: sns.Topic
	queryPlayerDamageDoneTopic: sns.Topic
    errorMail: string
}

export class EtlDashboard extends cdk.Construct {
	constructor(scope: cdk.Construct, id: string, props: Props) {
		super(scope, id)
		
		const presignApiId = props.presignApi.httpApiId

		// const key = new kms.Key(this, 'SnsKmsKey', {
		// 	enableKeyRotation: true,
		// })
		// const errorTopic = new sns.Topic(this, 'errorTopic', {
		// 	masterKey: key,
		// });
		// errorTopic.addSubscription(new subscriptions.EmailSubscription(props.errorMail));

		//NOTE widget height viable values: 3, 6, ?
		new cloudwatch.Dashboard(this, 'Dashboard', {
			start: "-P3H",
		}).addWidgets(
			new Row(
				new TextWidget({
					markdown: `# Lambda metrics

**Convert** takes the combatlog, normal or compressed, and converts it to a format that can be imported to amazon timestream.

**SummaryGet** queries the freshly imported data and publishes it to an sns topic

**SummariesInsert** get invoked by the topic and store the data
					`,
					width: 4,
					height: 6,
				}),
				new GraphWidget({
					title: 'Invocations',
					left: [
						props.convertLambda.metricInvocations({period: cdk.Duration.minutes(1)}),
						props.queryKeys.metricInvocations({period: cdk.Duration.minutes(1)}),
						props.insertKeysToDynamoDB.metricInvocations({period: cdk.Duration.minutes(1)}),
						props.insertKeysToTimestream.metricInvocations({period: cdk.Duration.minutes(1)}),
						props.insertPlayerDamageDoneToDynamodb.metricInvocations({period: cdk.Duration.minutes(1)}),
						props.queryPlayerDamageDone.metricInvocations({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4,
				}),
				new GraphWidget({
					title: 'Errors',
					left: [
						props.convertLambda.metricErrors({period: cdk.Duration.minutes(1)}),
						props.queryKeys.metricErrors({period: cdk.Duration.minutes(1)}),
						props.insertKeysToDynamoDB.metricErrors({period: cdk.Duration.minutes(1)}),
						props.insertKeysToTimestream.metricErrors({period: cdk.Duration.minutes(1)}),
						props.insertPlayerDamageDoneToDynamodb.metricErrors({period: cdk.Duration.minutes(1)}),
						props.queryPlayerDamageDone.metricErrors({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Throttles',
					left: [
						props.convertLambda.metricThrottles({period: cdk.Duration.minutes(1)}),
						props.queryKeys.metricThrottles({period: cdk.Duration.minutes(1)}),
						props.insertKeysToDynamoDB.metricThrottles({period: cdk.Duration.minutes(1)}),
						props.insertKeysToTimestream.metricThrottles({period: cdk.Duration.minutes(1)}),
						props.insertPlayerDamageDoneToDynamodb.metricThrottles({period: cdk.Duration.minutes(1)}),
						props.queryPlayerDamageDone.metricThrottles({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Average Duration',
					left: [
						props.convertLambda.metricDuration({period: cdk.Duration.minutes(1)}),
						props.queryKeys.metricDuration({period: cdk.Duration.minutes(1)}),
						props.insertKeysToDynamoDB.metricDuration({period: cdk.Duration.minutes(1)}),
						props.insertKeysToTimestream.metricDuration({period: cdk.Duration.minutes(1)}),
						props.insertPlayerDamageDoneToDynamodb.metricDuration({period: cdk.Duration.minutes(1)}),
						props.queryPlayerDamageDone.metricDuration({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Concurrent Executions',
					left: [
						props.convertLambda.metric('ConcurrentExecutions',{ statistic: 'Maximum', period: cdk.Duration.minutes(1) }),
						props.queryKeys.metric('ConcurrentExecutions',{ statistic: 'Maximum', period: cdk.Duration.minutes(1) }),
						props.insertKeysToDynamoDB.metric('ConcurrentExecutions',{ statistic: 'Maximum', period: cdk.Duration.minutes(1) }),
						props.insertKeysToTimestream.metric('ConcurrentExecutions',{ statistic: 'Maximum', period: cdk.Duration.minutes(1) }),
						props.insertPlayerDamageDoneToDynamodb.metric('ConcurrentExecutions',{ statistic: 'Maximum', period: cdk.Duration.minutes(1) }),
						props.queryPlayerDamageDone.metric('ConcurrentExecutions',{ statistic: 'Maximum', period: cdk.Duration.minutes(1) }),
					],
					stacked: false,
					width: 4
				}),
			),
			new Row(
				new TextWidget({
					markdown: `# SQS metrics
					
visible messages and message age *should* be as low as possible

messages in convert DLQ *should* be 0, the import and summary DLQ *must* be 0
					`,
					width: 4,
					height: 6,
				}),
				new GraphWidget({
					title: 'Visible messages',
					left: [
						props.convertQueue.metricApproximateNumberOfMessagesVisible({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Messages not visible',
					left: [
						props.convertQueue.metricApproximateNumberOfMessagesNotVisible({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Messages received',
					left: [
						props.convertQueue.metricNumberOfMessagesReceived({period: cdk.Duration.minutes(1)}),
						props.convertQueue.metricNumberOfMessagesSent({period: cdk.Duration.minutes(1)}),
						props.convertQueue.metricNumberOfMessagesDeleted({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Age of oldest message',
					left: [
						props.convertQueue.metricApproximateAgeOfOldestMessage({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'DLQ messages',
					left: [
						props.convertDLQ.metricApproximateNumberOfMessagesVisible({period: cdk.Duration.minutes(1)}),
						props.queryPlayerDamageDoneLambdaDLQ.metricApproximateNumberOfMessagesVisible({period: cdk.Duration.minutes(1)}),
						props.queryPlayerDamageDoneTopicDLQ.metricApproximateNumberOfMessagesVisible({period: cdk.Duration.minutes(1)}),
						props.insertPlayerDamageDoneDynamoDBLambdaDLQ.metricApproximateNumberOfMessagesVisible({period: cdk.Duration.minutes(1)}),
						props.queryKeysLambdaDLQ.metricApproximateNumberOfMessagesVisible({period: cdk.Duration.minutes(1)}),
						props.queryKeysTopicDLQ.metricApproximateNumberOfMessagesVisible({period: cdk.Duration.minutes(1)}),
						props.insertKeysToDynamoDBLambdaDLQ.metricApproximateNumberOfMessagesVisible({period: cdk.Duration.minutes(1)}),
						props.insertKeysToTimestreamLambdaDLQ.metricApproximateNumberOfMessagesVisible({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4
				}),
			),
			new Row(
				new TextWidget({
					markdown: `# Timestream

all combatlog data is ingested into timestream to process.
					`,
					width: 4,
					height: 6,
				}),
				new GraphWidget({
					title: 'Number of WriteRecords',
					left: [
						new cloudwatch.Metric({
							metricName: 'SuccessfulRequestLatency',
							namespace: 'AWS/Timestream',
							dimensions: { DatabaseName: 'wowmate-analytics', TableName: 'combatlogs', Operation: 'WriteRecords' },
							statistic: 'SampleCount',
							period: cdk.Duration.minutes(1),
						}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'WriteRecords Latency',
					left: [
						new cloudwatch.Metric({
							metricName: 'SuccessfulRequestLatency',
							namespace: 'AWS/Timestream',
							dimensions: { DatabaseName: 'wowmate-analytics', TableName: 'combatlogs', Operation: 'WriteRecords' },
							statistic: 'Avg',
							period: cdk.Duration.minutes(1),
						}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Query Latency',
					left: [
						new cloudwatch.Metric({
							metricName: 'SuccessfulRequestLatency',
							namespace: 'AWS/Timestream',
							dimensions: { Operation: 'Query' },
							statistic: 'Avg',
							period: cdk.Duration.minutes(1),
						}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Number of Queries',
					left: [
						new cloudwatch.Metric({
							metricName: 'SuccessfulRequestLatency',
							namespace: 'AWS/Timestream',
							dimensions: { Operation: 'Query' },
							statistic: 'SampleCount',
							label: 'Query',
							period: cdk.Duration.minutes(1),
						}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Errors',
					left: [
						new cloudwatch.Metric({
							metricName: 'UserErrors',
							namespace: 'AWS/Timestream',
							dimensions: { Operation: 'Query' },
							statistic: 'Sum',
							period: cdk.Duration.minutes(1),
						}),
						new cloudwatch.Metric({
							metricName: 'UserErrors',
							namespace: 'AWS/Timestream',
							dimensions: { Operation: 'WriteRecords' },
							statistic: 'Sum',
							period: cdk.Duration.minutes(1),
						}),
						new cloudwatch.Metric({
							metricName: 'SystemErrors',
							namespace: 'AWS/Timestream',
							dimensions: { Operation: 'Query' },
							statistic: 'Sum',
							period: cdk.Duration.minutes(1),
						}),
						new cloudwatch.Metric({
							metricName: 'SystemErrors',
							namespace: 'AWS/Timestream',
							dimensions: { Operation: 'WriteRecords' },
							statistic: 'Sum',
							period: cdk.Duration.minutes(1),
						}),
					],
					stacked: false,
					width: 4
				}),
			),
			new Row(
				new TextWidget({
					markdown: `# DynamoDB

contains the summaries from timestream
					`,
					width: 4,
					height: 6,
				}),
				new GraphWidget({
					// title: 'sasd',
					left: [
						props.dynamoDB.metricConsumedWriteCapacityUnits({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4,
				}),
				new GraphWidget({
					// title: '',
					left: [
						props.dynamoDB.metricConsumedReadCapacityUnits({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4,
				}),
				new GraphWidget({
					// title: '',
					left: [
						props.dynamoDB.metricThrottledRequests({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4,
				}),
				new GraphWidget({
					// title: '',
					left: [
						props.dynamoDB.metricSystemErrorsForOperations({period: cdk.Duration.minutes(1)}),
					],
					width: 4,
					stacked: false,
				}),
				new GraphWidget({
					// title: 's',
					left: [
						props.dynamoDB.metricUserErrors({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4,
				}),
			),
			new Row(
				new TextWidget({
					markdown: `# Presign/Upload Metrics
					
These components (AGW, Lambda and s3 bucket) are responsible to allow users to upload to a private bucket.
					`,
					width: 4,
					height: 6,
				}),
				new GraphWidget({
					title: 'Lambda throttles/errors',
					left: [
						props.presignLambda.metricThrottles({period: cdk.Duration.minutes(1)}),
						props.presignLambda.metricErrors({period: cdk.Duration.minutes(1)}),
						props.presignLambda.metric('ConcurrentExecutions',{ statistic: 'Maximum', period: cdk.Duration.minutes(1) }),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Uploaded',
					left: [
						new cloudwatch.Metric({
							metricName: 'BytesUploaded',
							namespace: 'AWS/S3',
							dimensions: { BucketName: props.uploadBucket.bucketName, FilterId: 'metric' },
							statistic: 'Sum',
							period: cdk.Duration.minutes(1),
						}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Bucket objects + size',
					left: [
						new cloudwatch.Metric({
							metricName: 'NumberOfObjects',
							namespace: 'AWS/S3',
							dimensions: { BucketName: props.uploadBucket.bucketName, StorageType: 'AllStorageTypes' },
							statistic: 'Sum',
							period: cdk.Duration.minutes(1),
						}),
					],
					right: [
						new cloudwatch.Metric({
							metricName: 'BucketSizeBytes',
							namespace: 'AWS/S3',
							dimensions: { BucketName: props.uploadBucket.bucketName, StorageType: 'StandardStorage' },
							statistic: 'Sum',
							period: cdk.Duration.minutes(1),
						}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Bucket errors',
					left: [
						new cloudwatch.Metric({
							metricName: '4xxErrors',
							namespace: 'AWS/S3',
							dimensions: { BucketName: props.uploadBucket.bucketName, FilterId: 'metric' },
							statistic: 'Sum',
							period: cdk.Duration.minutes(1),
						}),
						new cloudwatch.Metric({
							metricName: '5xxErros',
							namespace: 'AWS/S3',
							dimensions: { BucketName: props.uploadBucket.bucketName, FilterId: 'metric' },
							statistic: 'Sum',
							period: cdk.Duration.minutes(1),
						}),
					],
					stacked: false,
					width: 4
				}),
				//TODO: fix metrics for agw v2
				new GraphWidget({
					title: 'ApiGateway errors',
					left: [
						new cloudwatch.Metric({
							metricName: '4xx',
							namespace: 'AWS/ApiGateway',
							dimensions: { ApiId: presignApiId }, //NOTE: ApiName is not exposed on the object
							statistic: 'Sum',
							period: cdk.Duration.minutes(1),
						}),
						new cloudwatch.Metric({
							metricName: '5xx',
							namespace: 'AWS/ApiGateway',
							dimensions: { ApiId: presignApiId }, //NOTE: ApiName is not exposed on the object
							statistic: 'Sum',
							period: cdk.Duration.minutes(1),
						}),
					],
					stacked: false,
					width: 4
				}),
			),
			new Row(
				new TextWidget({
					markdown: `# SNS
					
fans out processing of the combatlogs
					`,
					width: 4,
					height: 6,
				}),
				new GraphWidget({
					title: 'Number of messages published',
					left: [
						props.convertTopic.metricNumberOfMessagesPublished({period: cdk.Duration.minutes(1)}),
						props.queryKeysTopic.metricNumberOfMessagesPublished({period: cdk.Duration.minutes(1)}),
						props.queryPlayerDamageDoneTopic.metricNumberOfMessagesPublished({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Number Notifications delivered',
					left: [
						props.convertTopic.metricNumberOfNotificationsDelivered({period: cdk.Duration.minutes(1)}),
						props.queryKeysTopic.metricNumberOfNotificationsDelivered({period: cdk.Duration.minutes(1)}),
						props.queryPlayerDamageDoneTopic.metricNumberOfNotificationsDelivered({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Number Notifications failed',
					left: [
						props.convertTopic.metricNumberOfNotificationsFailed({period: cdk.Duration.minutes(1)}),
						props.queryKeysTopic.metricNumberOfNotificationsFailed({period: cdk.Duration.minutes(1)}),
						props.queryPlayerDamageDoneTopic.metricNumberOfNotificationsFailed({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Avg Size of messages published',
					left: [
						props.convertTopic.metricPublishSize({period: cdk.Duration.minutes(1)}),
						props.queryKeysTopic.metricPublishSize({period: cdk.Duration.minutes(1)}),
						props.queryPlayerDamageDoneTopic.metricPublishSize({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4
				}),
			)
		)
	}
}