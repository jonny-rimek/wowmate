import cdk = require('@aws-cdk/core');
import ec2 = require('@aws-cdk/aws-ec2');
import rds = require('@aws-cdk/aws-rds');
import * as secretsmanager from '@aws-cdk/aws-secretsmanager';
import { RetentionDays } from '@aws-cdk/aws-logs';
import { HttpProxyIntegration, HttpApi, LambdaProxyIntegration, HttpMethod } from '@aws-cdk/aws-apigatewayv2';
import iam = require('@aws-cdk/aws-iam');
import lambda = require('@aws-cdk/aws-lambda');
import { CfnOutput } from '@aws-cdk/core';

interface Props extends cdk.StackProps {
	vpc: ec2.IVpc;
	dbSecGrp: ec2.SecurityGroup
	dbSecret: secretsmanager.ISecret
	dbEndpoint: string
	// cluster: rds.DatabaseCluster
}

export class Api extends cdk.Construct {
	public readonly lambda: lambda.Function;
	public readonly api: HttpApi;

	constructor(scope: cdk.Construct, id: string, props: Props) {
		super(scope, id)

		const migrateLambda = new lambda.Function(this, 'MigrateLambda', {
			code: lambda.Code.fromAsset('services/migrate'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: cdk.Duration.seconds(30),
			environment: {
				SECRET_ARN: props.dbSecret.secretArn,
			},
			logRetention: RetentionDays.ONE_WEEK,
			tracing: lambda.Tracing.ACTIVE,
			vpc: props.vpc,
			securityGroups: [props.dbSecGrp],
			//TODO: add DLQ
		})
		props.dbSecret.grantRead(migrateLambda)


		const topDamageLambda = new lambda.Function(this, 'TopDamageLambda', {
			code: lambda.Code.fromAsset('services/api'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: cdk.Duration.seconds(30),
			environment: {
				DB_ENDPOINT: props.dbEndpoint,
				SECRET_ARN: props.dbSecret.secretArn,
			},
			logRetention: RetentionDays.ONE_WEEK,
			tracing: lambda.Tracing.ACTIVE,
			vpc: props.vpc,
			securityGroups: [props.dbSecGrp],
		})
		this.lambda = topDamageLambda
		props.dbSecret.grantRead(topDamageLambda)

		const topDamageIntegration = new LambdaProxyIntegration({
			handler: topDamageLambda
		})

		const httpApi = new HttpApi(this, 'Api')

		new CfnOutput(this, 'HttpApiEndpoint', {
			value: httpApi.url!,
		}),

		httpApi.addRoutes({
			path: '/api/combatlog/summary/{combatlog_uuid}/damage',
			methods: [HttpMethod.GET],
			integration: topDamageIntegration,
		})
		this.api = httpApi
	}
}