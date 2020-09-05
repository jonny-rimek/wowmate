import cdk = require('@aws-cdk/core');
import cloudwatch = require('@aws-cdk/aws-cloudwatch');
import { GraphWidget, IMetric, Metric, Row, TextWidget } from "@aws-cdk/aws-cloudwatch";
import { SnsAction } from '@aws-cdk/aws-cloudwatch-actions';
import sns = require('@aws-cdk/aws-sns');
import lambda = require('@aws-cdk/aws-lambda');
import sqs = require('@aws-cdk/aws-sqs');
import s3 = require('@aws-cdk/aws-s3');

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
					title: 'Lambda Invocations',
					left: [
						props.convertLambda.metricInvocations(),
						props.importLambda.metricInvocations(),
					],
					stacked: false,
					width: 4,
				}),
				new GraphWidget({
					title: 'Invocations',
					left: [
						props.convertLambda.metricErrors(),
						props.importLambda.metricErrors(),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Throttles',
					left: [
						props.convertLambda.metricThrottles(),
						props.importLambda.metricThrottles(),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Duration',
					left: [
						props.convertLambda.metricDuration(),
						props.importLambda.metricDuration(),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Concurrent Executions',
					left: [
						props.convertLambda.metric('ConcurrentExecutions',{ statistic: 'Maximum' }),
						props.importLambda.metric('ConcurrentExecutions',{ statistic: 'Maximum' }),
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
						props.convertQueue.metricApproximateNumberOfMessagesVisible(),
						props.importQueue.metricApproximateNumberOfMessagesVisible(),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Age of oldest message',
					left: [
						props.convertQueue.metricApproximateAgeOfOldestMessage(),
						props.importQueue.metricApproximateAgeOfOldestMessage(),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Messages not visible',
					left: [
						props.convertQueue.metricApproximateNumberOfMessagesNotVisible(),
						props.importQueue.metricApproximateNumberOfMessagesNotVisible(),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Messages recieved',
					left: [
						props.convertQueue.metricNumberOfMessagesReceived(),
						props.importQueue.metricNumberOfMessagesReceived(),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'DLQ messages',
					left: [
						props.convertDLQ.metricApproximateNumberOfMessagesVisible(),
						props.importDLQ.metricApproximateNumberOfMessagesVisible(),
					],
					stacked: false,
					width: 4
				}),
			)
		)
	}
}
