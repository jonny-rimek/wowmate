import cdk = require('@aws-cdk/core');
import { RetentionDays } from '@aws-cdk/aws-logs';
import { HttpApi, HttpMethod } from '@aws-cdk/aws-apigatewayv2';
import * as agwv2 from '@aws-cdk/aws-apigatewayv2';
import { LambdaProxyIntegration } from '@aws-cdk/aws-apigatewayv2-integrations';
import lambda = require('@aws-cdk/aws-lambda');
import * as dynamodb from '@aws-cdk/aws-dynamodb';
import * as logs from '@aws-cdk/aws-logs';
import s3 = require('@aws-cdk/aws-s3');
import { CfnOutput } from '@aws-cdk/core';
import * as origins from "@aws-cdk/aws-cloudfront-origins";
import cloudfront = require('@aws-cdk/aws-cloudfront');
import route53= require('@aws-cdk/aws-route53');
import acm = require('@aws-cdk/aws-certificatemanager');
import targets = require('@aws-cdk/aws-route53-targets');
import * as kms from "@aws-cdk/aws-kms";
import * as iam from "@aws-cdk/aws-iam";

interface Props extends cdk.StackProps {
	dynamoDB: dynamodb.Table,
	uploadBucket: s3.Bucket;
	hostedZoneId: string
	hostedZoneName: string
	apiDomainName: string
	accessLogBucket: s3.Bucket
}

export class Api extends cdk.Construct {
	public readonly getKeysLambda: lambda.Function;
	public readonly getKeysPerDungeonLambda: lambda.Function;
	public readonly getPlayerDamageDoneLambda: lambda.Function;
	public readonly presignLambda: lambda.Function;
	public readonly api: HttpApi;
	public readonly cloudfront: cloudfront.Distribution;

	constructor(scope: cdk.Construct, id: string, props: Props) {
		super(scope, id)

		this.getKeysLambda = new lambda.Function(this, 'GetKeysLambda', {
			code: lambda.Code.fromAsset('dist/api/combatlogs/keys/index/get'),
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
				LOG_LEVEL: "info", //only info or debug are support
				LOCAL: "false",
			},
			logRetention: RetentionDays.ONE_MONTH,
			tracing: lambda.Tracing.ACTIVE,
			reservedConcurrentExecutions: 10,
		})
		props.dynamoDB.grantReadData(this.getKeysLambda)
		const topOverallDamageIntegration = new LambdaProxyIntegration({
			handler: this.getKeysLambda
		})

		this.getKeysPerDungeonLambda = new lambda.Function(this, 'GetKeysPerDungeonLambda', {
			code: lambda.Code.fromAsset('dist/api/combatlogs/keys/_dungeon_id/get'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 10240,
			timeout: cdk.Duration.seconds(3),
			environment: {
				DYNAMODB_TABLE_NAME: props.dynamoDB.tableName,
				LOCAL: "false",
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

		this.getPlayerDamageDoneLambda = new lambda.Function(this, 'GetPlayerDamageDoneLambda', {
			code: lambda.Code.fromAsset('dist/api/combatlogs/keys/_combatlog_uuid/player-damage-done/get'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 10240,
			timeout: cdk.Duration.seconds(3),
			environment: {
				DYNAMODB_TABLE_NAME: props.dynamoDB.tableName,
				LOCAL: "false",
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

		this.presignLambda = new lambda.Function(this, 'Lambda', {
			runtime: lambda.Runtime.NODEJS_14_X,
			description: "allows to upload combatlogs to private s3 bucket",
			code: lambda.Code.fromAsset('services/upload/presign'),
			handler: 'index.handler',
			environment: {
				LOCAL: "false",
				BUCKET_NAME: props.uploadBucket.bucketName,
				AWS_NODEJS_CONNECTION_REUSE_ENABLED: "1",
			},
			memorySize: 128,
			reservedConcurrentExecutions: 100,
			tracing: lambda.Tracing.ACTIVE,
		});
		const presignIntegration = new LambdaProxyIntegration({
			handler: this.presignLambda
		})
		props.uploadBucket.grantPut(this.presignLambda);

		const httpApi = new HttpApi(this, 'Api', {
			corsPreflight: {
				allowOrigins: ["*"],
			},
			description: "wowmate API",
		})

		new CfnOutput(this, 'HttpApiEndpoint', {
			value: httpApi.url!,
		})

		httpApi.addRoutes({
			path: '/combatlogs/keys',
			methods: [HttpMethod.GET],
			integration: topOverallDamageIntegration,
		})
		httpApi.addRoutes({
			path: '/combatlogs/keys/{dungeon_id}',
			methods: [HttpMethod.GET],
			integration: topDungeonDamageIntegration,
		})
		httpApi.addRoutes({
			path: '/combatlogs/keys/{combatlog_uuid}/player-damage-done',
			methods: [HttpMethod.GET],
			integration: playerDamageAdvancedIntegration,
		})
		httpApi.addRoutes({
			path: '/presign/{filename}',
			methods: [HttpMethod.POST],
			integration: presignIntegration,
		})

		this.api = httpApi

		const hostedZone = route53.HostedZone.fromHostedZoneAttributes(this, 'HostedZone', {
			zoneName: props.hostedZoneName,
			hostedZoneId: props.hostedZoneId,
		});

		const cert = new acm.DnsValidatedCertificate(this, 'Certificate', {
			domainName: props.apiDomainName,
			hostedZone,
		});

		const allowCorsAndQueryString = new cloudfront.OriginRequestPolicy(this, 'AllowCorsAndQueryStringParam2', {
			originRequestPolicyName: 'AllowCorsAndQueryStringParam2',
			cookieBehavior: cloudfront.OriginRequestCookieBehavior.none(),
			queryStringBehavior: cloudfront.OriginRequestQueryStringBehavior.all(),
			headerBehavior: cloudfront.OriginRequestHeaderBehavior.allowList('origin'),
		})

		//make sure enhanced metrics is enabled via the GUI no CF support =(
		//https://console.aws.amazon.com/cloudfront/v2/home#/monitoring
		this.cloudfront = new cloudfront.Distribution(this, 'Distribution', {
			defaultBehavior: {
				origin: new origins.HttpOrigin(this.api.url!.replace('https://','').replace('/',''), {}),
				cachePolicy: cloudfront.CachePolicy.CACHING_DISABLED,
				originRequestPolicy: allowCorsAndQueryString,
				viewerProtocolPolicy: cloudfront.ViewerProtocolPolicy.HTTPS_ONLY,
				allowedMethods: cloudfront.AllowedMethods.ALLOW_ALL,
				cachedMethods: cloudfront.CachedMethods.CACHE_GET_HEAD_OPTIONS,
			},
			certificate: cert,
			domainNames: [props.apiDomainName],
			comment: "wowmate.io API",
			logBucket: props.accessLogBucket,
			logFilePrefix: 'apiCloudfront',
		})

		const log = new logs.LogGroup(this, 'HttpApiLog')
		const cfnLog = log.node.defaultChild as logs.CfnLogGroup
		cfnLog.cfnOptions.metadata = {
			cfn_nag: {
				rules_to_suppress: [
					{
						id: 'W84',
						reason: "Log group is encrypted with the default kms key by default, no need to specify one",
					},
				]
			}
		}

		const stage = this.api.defaultStage?.node.defaultChild as agwv2.CfnStage;
		stage.accessLogSettings = {
			destinationArn: log.logGroupArn,
			format: JSON.stringify({
				"requestId": "$context.requestId",
				"ip": "$context.identity.sourceIp",
				"caller": "$context.identity.caller",
				"user": "$context.identity.user",
				"requestTime": "$context.requestTime",
				"httpMethod": "$context.httpMethod",
				"resourcePath": "$context.resourcePath",
				"status": "$context.status",
				"protocol": "$context.protocol",
				"responseLength": "$context.responseLength"
			})
		}

		const cfnDist = this.cloudfront.node.defaultChild as cloudfront.CfnDistribution;
		cfnDist.addPropertyOverride('DistributionConfig.Origins.0.OriginShield', {
			Enabled: true,
			OriginShieldRegion: 'us-east-1',
		});

		new route53.ARecord(this, 'Alias', {
			zone: hostedZone,
			target: route53.RecordTarget.fromAlias(new targets.CloudFrontTarget(this.cloudfront)),
			recordName: props.apiDomainName,
		});

		new route53.AaaaRecord(this, 'AliasAAA', {
			zone: hostedZone,
			target: route53.RecordTarget.fromAlias(new targets.CloudFrontTarget(this.cloudfront)),
			recordName: props.apiDomainName,
		});
	}
}