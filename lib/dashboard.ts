import cdk = require('@aws-cdk/core');
import cloudwatch = require('@aws-cdk/aws-cloudwatch');
import { GraphWidget, IMetric, Metric } from "@aws-cdk/aws-cloudwatch";
import { SnsAction } from '@aws-cdk/aws-cloudwatch-actions';
import sns = require('@aws-cdk/aws-sns');
import lambda = require('@aws-cdk/aws-lambda');
import { DynamoProjectionExpression } from '@aws-cdk/aws-stepfunctions-tasks';

interface Props extends cdk.StackProps {
	convertLambda: lambda.Function
}

export class Dashboard extends cdk.Construct {
	constructor(scope: cdk.Construct, id: string, props: Props) {
		super(scope, id)

		new cloudwatch.Dashboard(this, 'Dashboard').addWidgets(
			new GraphWidget({
				title: 'somename',
				left: [props.convertLambda.metricInvocations()],
				// new cloudwatch.Metric({
				// 	// dimensions:  {"TableName":"sdakjdjs", "Operation": "GetItem"},
				// 	namespace: 'namespacename',
				// 	metricName: 'asdasdas',
				// 	period: cdk.Duration.minutes(5),
				// 	statistic: 'somenameforgrap',
				// }),
				stacked: false,
				width: 8
			})
		)
	}
}
