import cdk = require('@aws-cdk/core');
import ec2 = require('@aws-cdk/aws-ec2');
import s3 = require('@aws-cdk/aws-s3');
import lambda = require('@aws-cdk/aws-lambda');
import { RetentionDays } from '@aws-cdk/aws-logs';
import * as secretsmanager from '@aws-cdk/aws-secretsmanager';
import * as events from '@aws-cdk/aws-events';
import * as eventTargets from '@aws-cdk/aws-events-targets';

interface VpcProps extends cdk.StackProps {
	vpc: ec2.IVpc
	dbSecGrp: ec2.SecurityGroup
	dbSecret : secretsmanager.ISecret
	dbEndpoint: string
}

export class Partition extends cdk.Construct {
	public readonly partitionLambda: lambda.Function;

	constructor(scope: cdk.Construct, id: string, props: VpcProps) {
		super(scope, id)

		const partitionLambda = new lambda.Function(this, 'PartitionLambda', {
			code: lambda.Code.fromAsset('services/partition'),
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
		})
		props.dbSecret?.grantRead(partitionLambda)

		const partitionTarget = new eventTargets.LambdaFunction(partitionLambda)

		new events.Rule(this, 'PartitionSchedule', {
			//run every day at 4:30 am
			schedule: events.Schedule.cron({ 
				minute: '30',
				hour: '4',
			}),
			targets: [partitionTarget],
		})
	}
}
