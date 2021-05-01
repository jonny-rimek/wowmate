import * as synthetics from '@aws-cdk/aws-synthetics';
import cdk = require('@aws-cdk/core');
import lambda = require('@aws-cdk/aws-lambda');
import cloudwatch = require('@aws-cdk/aws-cloudwatch');

interface Props extends cdk.StackProps {
}

export class Synthetics extends cdk.Construct {

    constructor(scope: cdk.Construct, id: string, props: Props) {
        super(scope, id)

        const canary = new synthetics.Canary(this, 'Canary', {
            schedule: synthetics.Schedule.rate(cdk.Duration.minutes(60)),
            test: synthetics.Test.custom({
                code: lambda.Code.fromAsset("services/synthetics"),
                handler: 'index.handler',
            }),
            runtime: synthetics.Runtime.SYNTHETICS_NODEJS_PUPPETEER_3_0,
        });

        new cloudwatch.Alarm(this, 'CanaryAlarm', {
            metric: canary.metricSuccessPercent(),
            evaluationPeriods: 2,
            threshold: 100,
            comparisonOperator: cloudwatch.ComparisonOperator.LESS_THAN_THRESHOLD,
        });
        //addAlarmAction(new SnsAction(errorTopic))


    }
}
