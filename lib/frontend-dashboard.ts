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
import { cfnTagToCloudFormation } from '@aws-cdk/core';

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
							namespace: 'AWS/CloudFront',
							dimensions: { DistributionId: props.cloudfront.distributionId, Region: 'Global' },
							statistic: 'Average',
							period: cdk.Duration.minutes(1),
						}),
						new cloudwatch.Metric({
							metricName: '5xxErrorRate',
							namespace: 'AWS/CloudFront',
							dimensions: { DistributionId: props.cloudfront.distributionId, Region: 'Global' },
							statistic: 'Average',
							period: cdk.Duration.minutes(1),
						}),
						new cloudwatch.Metric({
							metricName: '4xxErrorRate',
							namespace: 'AWS/CloudFront',
							dimensions: { DistributionId: props.cloudfront.distributionId, Region: 'Global' },
							statistic: 'Average',
							period: cdk.Duration.minutes(1),
						}),
						new cloudwatch.Metric({
							metricName: '401',
							namespace: 'AWS/CloudFront',
							dimensions: { DistributionId: props.cloudfront.distributionId, Region: 'Global' },
							statistic: 'Average',
							period: cdk.Duration.minutes(1),
						}),
						new cloudwatch.Metric({
							metricName: '403',
							namespace: 'AWS/CloudFront',
							dimensions: { DistributionId: props.cloudfront.distributionId, Region: 'Global' },
							statistic: 'Average',
							period: cdk.Duration.minutes(1),
						}),
						new cloudwatch.Metric({
							metricName: '502',
							namespace: 'AWS/CloudFront',
							dimensions: { DistributionId: props.cloudfront.distributionId, Region: 'Global' },
							statistic: 'Average',
							period: cdk.Duration.minutes(1),
						}),
						new cloudwatch.Metric({
							metricName: '503ErrorRate',
							namespace: 'AWS/CloudFront',
							dimensions: { DistributionId: props.cloudfront.distributionId, Region: 'Global' },
							statistic: 'Average',
							period: cdk.Duration.minutes(1),
						}),
						new cloudwatch.Metric({
							metricName: '504ErrorRate',
							namespace: 'AWS/CloudFront',
							dimensions: { DistributionId: props.cloudfront.distributionId, Region: 'Global' },
							statistic: 'Average',
							period: cdk.Duration.minutes(1),
						}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Origin Latency',
					left: [
						new cloudwatch.Metric({
							metricName: 'OriginLatency',
							namespace: 'AWS/CloudFront',
							dimensions: { DistributionId: props.cloudfront.distributionId, Region: 'Global' },
							statistic: 'Average',
							period: cdk.Duration.minutes(1),
						}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Cache Hit Rate',
					left: [
						new cloudwatch.Metric({
							metricName: 'CacheHitRate',
							namespace: 'AWS/CloudFront',
							dimensions: { DistributionId: props.cloudfront.distributionId, Region: 'Global' },
							statistic: 'Average',
							period: cdk.Duration.minutes(1),
						}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Requests',
					left: [
						new cloudwatch.Metric({
							metricName: 'Requests',
							namespace: 'AWS/CloudFront',
							dimensions: { DistributionId: props.cloudfront.distributionId, Region: 'Global' },
							statistic: 'Sum',
							period: cdk.Duration.minutes(1),
						}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Bytes Downloaded/Uploaded ',
					left: [
						new cloudwatch.Metric({
							metricName: 'BytesDownloaded',
							namespace: 'AWS/CloudFront',
							dimensions: { DistributionId: props.cloudfront.distributionId, Region: 'Global' },
							statistic: 'Sum',
							period: cdk.Duration.minutes(1),
						}),
						new cloudwatch.Metric({
							metricName: 'BytesUploaded',
							namespace: 'AWS/CloudFront',
							dimensions: { DistributionId: props.cloudfront.distributionId, Region: 'Global' },
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