import cdk = require('@aws-cdk/core');
import { RetentionDays } from '@aws-cdk/aws-logs';
import { HttpApi, HttpMethod } from '@aws-cdk/aws-apigatewayv2';
import { LambdaProxyIntegration } from '@aws-cdk/aws-apigatewayv2-integrations';
import lambda = require('@aws-cdk/aws-lambda');
import * as dynamodb from '@aws-cdk/aws-dynamodb';
import { CfnOutput } from '@aws-cdk/core';

interface Props extends cdk.StackProps {
	dynamoDB: dynamodb.Table,
}

export class Api extends cdk.Construct {
	public readonly getKeysLambda: lambda.Function;
	public readonly getKeysPerDungeonLambda: lambda.Function;
	public readonly getPlayerDamageDoneLambda: lambda.Function;
	public readonly api: HttpApi;

	constructor(scope: cdk.Construct, id: string, props: Props) {
		super(scope, id)

		//TODO: rename the id
		this.getKeysLambda = new lambda.Function(this, 'DamageSummariesLambda', {
			code: lambda.Code.fromAsset('services/api/combatlogs/keys/index/get'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 10240,
			//for max memory the cold start duration is ~90ms, the init duration is ~110ms
			//for min memory the cold start duration is ~800ms, this is because the dns resolution takes longer, because it
			//has lower network bandwidth, init duration is the same.
            //warm start only take ~6ms so leaving them at max memory won't be that expensive, optimize later if cost
			//of api lambdas becomes a factor
			timeout: cdk.Duration.seconds(3),
			environment: {
				DYNAMODB_TABLE_NAME: props.dynamoDB.tableName,
				LOG_LEVEL: "info" //only info or debug are support
			},
			logRetention: RetentionDays.ONE_MONTH,
			tracing: lambda.Tracing.ACTIVE,
			reservedConcurrentExecutions: 10,
		})
		props.dynamoDB.grantReadData(this.getKeysLambda)
		const topOverallDamageIntegration = new LambdaProxyIntegration({
			handler: this.getKeysLambda
		})

		this.getKeysPerDungeonLambda = new lambda.Function(this, 'DamageDungeonSummariesLambda', {
			code: lambda.Code.fromAsset('services/api/combatlogs/keys/_dungeon_id/get'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 10240,
			timeout: cdk.Duration.seconds(3),
			environment: {
				DYNAMODB_TABLE_NAME: props.dynamoDB.tableName,
				LOG_LEVEL: "info" //only info or debug are support
			},
			logRetention: RetentionDays.ONE_MONTH,
			tracing: lambda.Tracing.ACTIVE,
			reservedConcurrentExecutions: 10,
		})
		props.dynamoDB.grantReadData(this.getKeysPerDungeonLambda)
		const topDungeonDamageIntegration = new LambdaProxyIntegration({
			handler: this.getKeysPerDungeonLambda
		})

		this.getPlayerDamageDoneLambda = new lambda.Function(this, 'CombatlogPlayerDamageAdvanced', {
			code: lambda.Code.fromAsset('services/api/combatlogs/keys/_combatlog_uuid/player-damage-done/get'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 10240,
			timeout: cdk.Duration.seconds(3),
			environment: {
				DYNAMODB_TABLE_NAME: props.dynamoDB.tableName,
				LOG_LEVEL: "info" //only info or debug are support
			},
			logRetention: RetentionDays.ONE_MONTH,
			tracing: lambda.Tracing.ACTIVE,
			reservedConcurrentExecutions: 10,
		})
		props.dynamoDB.grantReadData(this.getPlayerDamageDoneLambda)
		const playerDamageAdvancedIntegration = new LambdaProxyIntegration({
			handler: this.getPlayerDamageDoneLambda
		})

		const httpApi = new HttpApi(this, 'Api', {
			/*
			corsPreflight: { //might need this if I break the api into its own domain
				allowOrigins: ["wowmate.io"],
			},
			 */
			description: "wowmate combatlog api",
		})

		new CfnOutput(this, 'HttpApiEndpoint', {
			value: httpApi.url!,
		})

		httpApi.addRoutes({
			//TODO: remove api part https://github.com/jonny-rimek/wowmate/issues/235
			path: '/api/combatlogs/keys',
			methods: [HttpMethod.GET],
			integration: topOverallDamageIntegration,
		})
		httpApi.addRoutes({
			path: '/api/combatlogs/keys/{dungeon_id}',
			methods: [HttpMethod.GET],
			integration: topDungeonDamageIntegration,
		})
		httpApi.addRoutes({
			path: '/api/combatlogs/keys/{combatlog_uuid}/player-damage-done',
			methods: [HttpMethod.GET],
			integration: playerDamageAdvancedIntegration,
		})

		this.api = httpApi
	}
}