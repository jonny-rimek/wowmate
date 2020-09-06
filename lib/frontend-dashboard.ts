import cdk = require('@aws-cdk/core');
import cloudwatch = require('@aws-cdk/aws-cloudwatch');
import { GraphWidget, IMetric, Metric, Row, TextWidget } from "@aws-cdk/aws-cloudwatch";
import { SnsAction } from '@aws-cdk/aws-cloudwatch-actions';
import sns = require('@aws-cdk/aws-sns');
import * as subscriptions from '@aws-cdk/aws-sns-subscriptions';
import lambda = require('@aws-cdk/aws-lambda');
import sqs = require('@aws-cdk/aws-sqs');
import { HttpApi } from '@aws-cdk/aws-apigatewayv2';
import s3 = require('@aws-cdk/aws-s3');
import cloudfront = require('@aws-cdk/aws-cloudfront');

interface Props extends cdk.StackProps {
	topDamageLambda: lambda.Function
	api: HttpApi
	s3: s3.Bucket
	cloudfront: cloudfront.CloudFrontWebDistribution,
}

export class FrontendDashboard extends cdk.Construct {
	constructor(scope: cdk.Construct, id: string, props: Props) {
		super(scope, id)

		const errorTopic = new sns.Topic(this, 'errorTopic');
		//TODO: reactivate in prod
		// errorTopic.addSubscription(new subscriptions.EmailSubscription('jimbo.db@protonmail.com'));

		/*
		new cloudwatch.Alarm(this, 'Failed to import to DB', {
			metric: props.importDLQ.metricApproximateNumberOfMessagesVisible({period: cdk.Duration.minutes(1)}),
			threshold: 1,
			evaluationPeriods: 2,
			datapointsToAlarm: 1,
			treatMissingData: cloudwatch.TreatMissingData.NOT_BREACHING,
		}).addAlarmAction(new SnsAction(errorTopic))
*/
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
						props.topDamageLambda.metricInvocations({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4,
				}),
				new GraphWidget({
					title: 'Errors',
					left: [
						props.topDamageLambda.metricErrors({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Throttles',
					left: [
						props.topDamageLambda.metricThrottles({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Duration',
					left: [
						props.topDamageLambda.metricDuration({period: cdk.Duration.minutes(1)}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Concurrent Executions',
					left: [
						props.topDamageLambda.metric('ConcurrentExecutions',{ statistic: 'Maximum', period: cdk.Duration.minutes(1) }),
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
					title: 'CloudFront Error Rate',
					left: [
						new cloudwatch.Metric({
							metricName: 'TotalErrorRate',
							namespace: 'AWS/Cloudfront',
							statistic: 'Average',
							period: cdk.Duration.minutes(1),
						}),
					],
					stacked: false,
					width: 4
				}),
	/* 			new GraphWidget({
					title: 'Age of oldest message',
					left: [
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Messages not visible',
					left: [
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Messages recieved',
					left: [
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'DLQ messages',
					left: [
					],
					stacked: false,
					width: 4 */
				// }),
			)
		)
	}
}
