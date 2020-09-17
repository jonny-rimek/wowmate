import cdk = require('@aws-cdk/core');
import cloudwatch = require('@aws-cdk/aws-cloudwatch');
import { GraphWidget, IMetric, Metric, Row, TextWidget } from "@aws-cdk/aws-cloudwatch";
import { SnsAction } from '@aws-cdk/aws-cloudwatch-actions';
import sns = require('@aws-cdk/aws-sns');
import * as subscriptions from '@aws-cdk/aws-sns-subscriptions';
import lambda = require('@aws-cdk/aws-lambda');
import sqs = require('@aws-cdk/aws-sqs');
import s3 = require('@aws-cdk/aws-s3');
import apigateway = require('@aws-cdk/aws-apigateway');
import rds = require('@aws-cdk/aws-rds');

interface Props extends cdk.StackProps {
	convertLambda: lambda.Function
	convertQueue: sqs.Queue
	convertDLQ: sqs.Queue
	importLambda: lambda.Function
	importQueue: sqs.Queue
	importDLQ: sqs.Queue
	presignLambda: lambda.Function
	uploadBucket: s3.Bucket
	presignApiGateway: apigateway.LambdaRestApi
	cluster: rds.DatabaseCluster;
}

export class EtlDashboard extends cdk.Construct {
	constructor(scope: cdk.Construct, id: string, props: Props) {
		super(scope, id)

		const errorTopic = new sns.Topic(this, 'errorTopic');
		errorTopic.addSubscription(new subscriptions.EmailSubscription('jimbo.db@protonmail.com'));

		new cloudwatch.Alarm(this, 'Failed to import to DB', {
			metric: props.importDLQ.metricApproximateNumberOfMessagesVisible({period: cdk.Duration.minutes(1)}),
			threshold: 1,
			evaluationPeriods: 2,
			datapointsToAlarm: 1,
			treatMissingData: cloudwatch.TreatMissingData.NOT_BREACHING,
		}).addAlarmAction(new SnsAction(errorTopic))

		//NOTE widget height viable values: 3, 6, ?
		new cloudwatch.Dashboard(this, 'Dashboard').addWidgets(
			new Row(
				new TextWidget({
					markdown: `# Postgres Aurora Cluster

Crucial is the write IOPS, because we are ingesting a ton of data
					`,
					width: 4,
					height: 6,
				}),
				new GraphWidget({
					title: 'Write IOPS',
					left: [
						props.cluster.metricVolumeWriteIOPs({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4,
				}),
				new GraphWidget({
					title: 'Network recieve througput',
					left: [
						props.cluster.metricNetworkReceiveThroughput({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'CPU utilization',
					left: [
						props.cluster.metricCPUUtilization({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Storage',
					left: [
						props.cluster.metricVolumeBytesUsed({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Connections',
					left: [
						props.cluster.metricDatabaseConnections({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4
				}),
			),
			new Row(
				new TextWidget({
					markdown: `# Lambda metrics

**Convert** takes the combatlog, normal or compressed, and converts it to a csv and passes it into the csvBucket

**Import** takes the processed combat log and loads it into postgres aurora
					`,
					width: 4,
					height: 6,
				}),
				new GraphWidget({
					title: 'Invocations',
					left: [
						props.convertLambda.metricInvocations({period: cdk.Duration.minutes(1)}),
						props.importLambda.metricInvocations({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4,
				}),
				new GraphWidget({
					title: 'Errors',
					left: [
						props.convertLambda.metricErrors({period: cdk.Duration.minutes(1)}),
						props.importLambda.metricErrors({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Throttles',
					left: [
						props.convertLambda.metricThrottles({period: cdk.Duration.minutes(1)}),
						props.importLambda.metricThrottles({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Duration',
					left: [
						props.convertLambda.metricDuration({period: cdk.Duration.minutes(1)}),
						props.importLambda.metricDuration({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Concurrent Executions',
					left: [
						props.convertLambda.metric('ConcurrentExecutions',{ statistic: 'Maximum', period: cdk.Duration.minutes(1) }),
						props.importLambda.metric('ConcurrentExecutions',{ statistic: 'Maximum', period: cdk.Duration.minutes(1) }),
					],
					stacked: false,
					width: 4
				}),
			),
			new Row(
				new TextWidget({
					markdown: `# SQS metrics
					
visible messages and message age *should* be as low as possible

messages in convert DLQ *should* be 0, the import DLQ _must_ be 0
					`,
					width: 4,
					height: 6,
				}),
				new GraphWidget({
					title: 'Visible messages',
					left: [
						props.convertQueue.metricApproximateNumberOfMessagesVisible({period: cdk.Duration.minutes(1)}),
						props.importQueue.metricApproximateNumberOfMessagesVisible({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Age of oldest message',
					left: [
						props.convertQueue.metricApproximateAgeOfOldestMessage({period: cdk.Duration.minutes(1)}),
						props.importQueue.metricApproximateAgeOfOldestMessage({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Messages not visible',
					left: [
						props.convertQueue.metricApproximateNumberOfMessagesNotVisible({period: cdk.Duration.minutes(1)}),
						props.importQueue.metricApproximateNumberOfMessagesNotVisible({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Messages recieved',
					left: [
						props.convertQueue.metricNumberOfMessagesReceived({period: cdk.Duration.minutes(1)}),
						props.importQueue.metricNumberOfMessagesReceived({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'DLQ messages',
					left: [
						props.convertDLQ.metricApproximateNumberOfMessagesVisible({period: cdk.Duration.minutes(1)}),
						props.importDLQ.metricApproximateNumberOfMessagesVisible({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4
				}),
			),
			new Row(
				new TextWidget({
					markdown: `# Presign/Upload Metrics
					
These 3 components (AGW, Lambda and s3 bucket) are responsible to allow users to upload to a private bucket.
					`,
					width: 4,
					height: 6,
				}),
				new GraphWidget({
					title: 'Lambda throttles/errors',
					left: [
						props.presignLambda.metricThrottles({period: cdk.Duration.minutes(1)}),
						props.presignLambda.metricErrors({period: cdk.Duration.minutes(1)}),
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
				new GraphWidget({
					title: 'ApiGateway errors',
					left: [
						new cloudwatch.Metric({
							metricName: '4XXError',
							namespace: 'AWS/ApiGateway',
							dimensions: { ApiName: 'PresignApi' }, //NOTE: ApiName is not exposed on the object
							statistic: 'Sum',
							period: cdk.Duration.minutes(1),
						}),
						new cloudwatch.Metric({
							metricName: '5XXError',
							namespace: 'AWS/ApiGateway',
							dimensions: { ApiName: 'PresignApi' }, //NOTE: ApiName is not exposed on the object
							statistic: 'Sum',
							period: cdk.Duration.minutes(1),
						}),
					],
					stacked: false,
					width: 4
				}),
			)
		)
	}
}
