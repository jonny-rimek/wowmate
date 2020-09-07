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

		const agwLatency = new cloudwatch.Metric({
								metricName: 'Latency',
								namespace: 'AWS/ApiGateway',
								dimensions: { ApiId: props.api.httpApiId },
								statistic: 'Avg',
								period: cdk.Duration.minutes(1),
							})
		
		const agwIntegrationLatency = new cloudwatch.Metric({
								metricName: 'IntegrationLatency',
								namespace: 'AWS/ApiGateway',
								dimensions: { ApiId: props.api.httpApiId },
								statistic: 'Avg',
								period: cdk.Duration.minutes(1),
							})

		let apiGatewayOverhead = new cloudwatch.MathExpression({
			expression: 'm1 - m2',
			label: 'API Gateway overhead',
			usingMetrics: {
				m1: agwLatency,
				m2: agwIntegrationLatency,
			},
			period: cdk.Duration.minutes(5)
		});

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
					markdown: `# CloudFront
					
CloudFront is the CDN that is a globally distributed cache, to speed up the users requests.

Due to the fact that log data is static and never changes we can cache all requests about combatlogs indefintely.
For metrics about combatlogs (damage per specc in all dungeons) the TTL is lower.

Origin Latency is the time the request spends after hitting CloudFront from the user, until the request is sent back
and leaves the AWS network and should be as low as possible.
					`,
					width: 4,
					height: 6,
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
					title: 'Error Rate',
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
/* 						new cloudwatch.Metric({
							metricName: '401ErrorRate',
							namespace: 'AWS/CloudFront',
							dimensions: { DistributionId: props.cloudfront.distributionId, Region: 'Global' },
							statistic: 'Average',
							period: cdk.Duration.minutes(1),
						}),
						new cloudwatch.Metric({
							metricName: '403ErrorRate',
							namespace: 'AWS/CloudFront',
							dimensions: { DistributionId: props.cloudfront.distributionId, Region: 'Global' },
							statistic: 'Average',
							period: cdk.Duration.minutes(1),
						}),
						new cloudwatch.Metric({
							metricName: '404ErrorRate',
							namespace: 'AWS/CloudFront',
							dimensions: { DistributionId: props.cloudfront.distributionId, Region: 'Global' },
							statistic: 'Average',
							period: cdk.Duration.minutes(1),
						}),
						new cloudwatch.Metric({
							metricName: '502ErrorRate',
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
	 */				],
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
			),
			new Row(
				new TextWidget({
					markdown: `# HTTP API

API Gatewayv2 is fronting all public lambdas.
					`,
					width: 4,
					height: 6,
				}),
				new GraphWidget({
					title: 'Count',
					left: [
						new cloudwatch.Metric({
							metricName: 'Count',
							namespace: 'AWS/ApiGateway',
							dimensions: { ApiId: props.api.httpApiId },
							statistic: 'Sum',
							period: cdk.Duration.minutes(1),
						}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Error Rates',
					left: [
						new cloudwatch.Metric({
							metricName: '5xx',
							namespace: 'AWS/ApiGateway',
							dimensions: { ApiId: props.api.httpApiId },
							statistic: 'Average',
							period: cdk.Duration.minutes(1),
						}),
						new cloudwatch.Metric({
							metricName: '4xx',
							namespace: 'AWS/ApiGateway',
							dimensions: { ApiId: props.api.httpApiId },
							statistic: 'Average',
							period: cdk.Duration.minutes(1),
						}),
				],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					left: [
						agwLatency,
/* 
						new cloudwatch.Metric({
							metricName: 'Latency',
							namespace: 'AWS/ApiGateway',
							dimensions: { ApiId: props.api.httpApiId },
							statistic: 'p90',
							period: cdk.Duration.minutes(1),
						}),
						new cloudwatch.Metric({
							metricName: 'Latency',
							namespace: 'AWS/ApiGateway',
							dimensions: { ApiId: props.api.httpApiId },
							statistic: 'p95',
							period: cdk.Duration.minutes(1),
						}),
						new cloudwatch.Metric({
							metricName: 'Latency',
							namespace: 'AWS/ApiGateway',
							dimensions: { ApiId: props.api.httpApiId },
							statistic: 'p99',
							period: cdk.Duration.minutes(1),
						}),
						new cloudwatch.Metric({
							metricName: 'Latency',
							namespace: 'AWS/ApiGateway',
							dimensions: { ApiId: props.api.httpApiId },
							statistic: 'p10',
							period: cdk.Duration.minutes(1),
						}),
	 */				],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Data processed',
					left: [
						new cloudwatch.Metric({
							metricName: 'DataProcessed',
							namespace: 'AWS/ApiGateway',
							dimensions: { ApiId: props.api.httpApiId },
							statistic: 'Sum',
							period: cdk.Duration.minutes(1),
						}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Time spent in AGW (ms)',
					left: [
						apiGatewayOverhead
					],
					stacked: false,
					width: 4
				}),
			),
			new Row(
				new TextWidget({
					markdown: `# Lambda metrics

Those are all customer facing lambdas that access the database, ergo errors should be 0, no throttles and duration should be as low as possible.
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
					title: 'Duration',
					left: [
						props.topDamageLambda.metricDuration({period: cdk.Duration.minutes(1)}),
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
					markdown: `# S3 metrics

					`,
					width: 4,
					height: 6,
				}),
				new GraphWidget({
					title: 'Invocations',
					left: [
						new cloudwatch.Metric({
							metricName: 'GetRequests',
							namespace: 'AWS/S3',
							dimensions: { BucketName: props.s3.bucketName, FilterId: 'metric' },
							statistic: 'Sum',
							period: cdk.Duration.minutes(1),
						}),
					],
					stacked: false,
					width: 4,
				}),
				new GraphWidget({
					title: 'Errors',
					left: [
						new cloudwatch.Metric({
							metricName: '4xxErrors',
							namespace: 'AWS/S3',
							dimensions: { BucketName: props.s3.bucketName, FilterId: 'metric' },
							statistic: 'Sum',
							period: cdk.Duration.minutes(1),
						}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Duration',
					left: [
						new cloudwatch.Metric({
							metricName: 'TotalRequestLatency',
							namespace: 'AWS/S3',
							dimensions: { BucketName: props.s3.bucketName, FilterId: 'metric' },
							statistic: 'Average',
							period: cdk.Duration.minutes(1),
						}),
						new cloudwatch.Metric({
							metricName: 'FirstByteLatency',
							namespace: 'AWS/S3',
							dimensions: { BucketName: props.s3.bucketName, FilterId: 'metric' },
							statistic: 'Average',
							period: cdk.Duration.minutes(1),
						}),
					],
					stacked: false,
					width: 4
				}),
				new GraphWidget({
					title: 'Throttles',
					left: [
						new cloudwatch.Metric({
							metricName: 'BytesDownloaded',
							namespace: 'AWS/S3',
							dimensions: { BucketName: props.s3.bucketName, FilterId: 'metric' },
							statistic: 'Average',
							period: cdk.Duration.minutes(1),
						}),
					],
					stacked: false,
					width: 4
				}),
			),
		)
	}
}