import cdk = require('@aws-cdk/core');
import cloudwatch = require('@aws-cdk/aws-cloudwatch');
import { GraphWidget, IMetric, Metric, Row } from "@aws-cdk/aws-cloudwatch";
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
				new GraphWidget({
					title: 'Convert Lambda Invocations',
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
				new GraphWidget({
					title: 'Convert Queue visible messages',
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
					title: 'Number of messages not visible',
					left: [
						props.convertQueue.metricApproximateNumberOfMessagesNotVisible(),
						props.importQueue.metricApproximateNumberOfMessagesNotVisible(),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Number of messages recieved',
					left: [
						props.convertQueue.metricNumberOfMessagesReceived(),
						props.importQueue.metricNumberOfMessagesReceived(),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Convert DLQ messages',
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
