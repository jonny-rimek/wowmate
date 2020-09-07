import cdk = require('@aws-cdk/core');
import cloudwatch = require('@aws-cdk/aws-cloudwatch');
import { GraphWidget, IMetric, Metric, Row, TextWidget } from "@aws-cdk/aws-cloudwatch";
import { SnsAction } from '@aws-cdk/aws-cloudwatch-actions';
import sns = require('@aws-cdk/aws-sns');
import * as subscriptions from '@aws-cdk/aws-sns-subscriptions';
import lambda = require('@aws-cdk/aws-lambda');
import sqs = require('@aws-cdk/aws-sqs');

interface Props extends cdk.StackProps {
	convertLambda: lambda.Function
	convertQueue: sqs.Queue
	convertDLQ: sqs.Queue
	importLambda: lambda.Function
	importQueue: sqs.Queue
	importDLQ: sqs.Queue
}

export class EtlDashboard extends cdk.Construct {
	constructor(scope: cdk.Construct, id: string, props: Props) {
		super(scope, id)

		const errorTopic = new sns.Topic(this, 'errorTopic');
		//TODO: reactivate in prod
		errorTopic.addSubscription(new subscriptions.EmailSubscription('jimbo.db@protonmail.com'));

		new cloudwatch.Alarm(this, 'Failed to import to DB', {
			metric: props.importDLQ.metricApproximateNumberOfMessagesVisible({period: cdk.Duration.minutes(5)}),
			threshold: 1,
			evaluationPeriods: 2,
			datapointsToAlarm: 1,
			treatMissingData: cloudwatch.TreatMissingData.NOT_BREACHING,
		}).addAlarmAction(new SnsAction(errorTopic))

		//NOTE widget height viable values: 3, 6, ?
		new cloudwatch.Dashboard(this, 'Dashboard').addWidgets(
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
						props.convertLambda.metricInvocations({period: cdk.Duration.minutes(5)}),
						props.importLambda.metricInvocations({period: cdk.Duration.minutes(5)}),
					],
					stacked: false,
					width: 4,
				}),
				new GraphWidget({
					title: 'Errors',
					left: [
						props.convertLambda.metricErrors({period: cdk.Duration.minutes(5)}),
						props.importLambda.metricErrors({period: cdk.Duration.minutes(5)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Throttles',
					left: [
						props.convertLambda.metricThrottles({period: cdk.Duration.minutes(5)}),
						props.importLambda.metricThrottles({period: cdk.Duration.minutes(5)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Duration',
					left: [
						props.convertLambda.metricDuration({period: cdk.Duration.minutes(5)}),
						props.importLambda.metricDuration({period: cdk.Duration.minutes(5)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Concurrent Executions',
					left: [
						props.convertLambda.metric('ConcurrentExecutions',{ statistic: 'Maximum', period: cdk.Duration.minutes(5) }),
						props.importLambda.metric('ConcurrentExecutions',{ statistic: 'Maximum', period: cdk.Duration.minutes(5) }),
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
						props.convertQueue.metricApproximateNumberOfMessagesVisible({period: cdk.Duration.minutes(5)}),
						props.importQueue.metricApproximateNumberOfMessagesVisible({period: cdk.Duration.minutes(5)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Age of oldest message',
					left: [
						props.convertQueue.metricApproximateAgeOfOldestMessage({period: cdk.Duration.minutes(5)}),
						props.importQueue.metricApproximateAgeOfOldestMessage({period: cdk.Duration.minutes(5)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Messages not visible',
					left: [
						props.convertQueue.metricApproximateNumberOfMessagesNotVisible({period: cdk.Duration.minutes(5)}),
						props.importQueue.metricApproximateNumberOfMessagesNotVisible({period: cdk.Duration.minutes(5)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Messages recieved',
					left: [
						props.convertQueue.metricNumberOfMessagesReceived({period: cdk.Duration.minutes(5)}),
						props.importQueue.metricNumberOfMessagesReceived({period: cdk.Duration.minutes(5)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'DLQ messages',
					left: [
						props.convertDLQ.metricApproximateNumberOfMessagesVisible({period: cdk.Duration.minutes(5)}),
						props.importDLQ.metricApproximateNumberOfMessagesVisible({period: cdk.Duration.minutes(5)}),
					],
					stacked: false,
					width: 4
				}),
			)
		)
	}
}
