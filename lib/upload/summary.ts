import cdk = require('@aws-cdk/core');
import ec2 = require('@aws-cdk/aws-ec2');
import s3 = require('@aws-cdk/aws-s3');
import sqs = require('@aws-cdk/aws-sqs');
import lambda = require('@aws-cdk/aws-lambda');
import { RetentionDays } from '@aws-cdk/aws-logs';
import * as destinations from '@aws-cdk/aws-lambda-destinations';
import * as secretsmanager from '@aws-cdk/aws-secretsmanager';

interface VpcProps extends cdk.StackProps {
	vpc: ec2.IVpc
	csvBucket: s3.Bucket
	dbSecGrp: ec2.SecurityGroup
	dbSecret : secretsmanager.ISecret
	dbEndpoint: string
}

export class Summary extends cdk.Construct {
	public readonly summaryLambda: lambda.Function;
	public readonly summaryDLQ: sqs.Queue;

	constructor(scope: cdk.Construct, id: string, props: VpcProps) {
		super(scope, id)

		this.summaryDLQ = new sqs.Queue(this, 'SummaryDeadLetterQueue', {
			retentionPeriod: cdk.Duration.days(14)
		})

		this.summaryLambda = new lambda.Function(this, 'SummaryLambda', {
			code: lambda.Code.fromAsset('services/summary'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: cdk.Duration.seconds(60),
			environment: {
				SECRET_ARN: props.dbSecret.secretArn,
				DB_ENDPOINT: props.dbEndpoint,
			},
			reservedConcurrentExecutions: 10, 
			logRetention: RetentionDays.ONE_WEEK,
			tracing: lambda.Tracing.ACTIVE,
			vpc: props.vpc,
			securityGroups: [props.dbSecGrp],
			onFailure: new destinations.SqsDestination(this.summaryDLQ)
		})
		props.dbSecret?.grantRead(this.summaryLambda)
	}
}
