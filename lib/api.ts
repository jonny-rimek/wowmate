import cdk = require('@aws-cdk/core');
import ddb = require('@aws-cdk/aws-dynamodb');
import s3 = require('@aws-cdk/aws-s3');
import lambda = require('@aws-cdk/aws-lambda');
import sfn = require('@aws-cdk/aws-stepfunctions');
import tasks = require('@aws-cdk/aws-stepfunctions-tasks');
import { RemovalPolicy, Duration } from '@aws-cdk/core';
import { BlockPublicAccess } from '@aws-cdk/aws-s3';
import iam = require('@aws-cdk/aws-iam');
import targets = require('@aws-cdk/aws-route53-targets');
import { Effect } from '@aws-cdk/aws-iam';
import cloudtrail = require('@aws-cdk/aws-cloudtrail');
import apigateway = require('@aws-cdk/aws-apigateway');
import s3deploy = require('@aws-cdk/aws-s3-deployment');
import cloudfront = require('@aws-cdk/aws-cloudfront');
import route53= require('@aws-cdk/aws-route53');
import acm = require('@aws-cdk/aws-certificatemanager');
import { SSLMethod, SecurityPolicyProtocol, OriginAccessIdentity } from '@aws-cdk/aws-cloudfront';
import { StateMachineType } from '@aws-cdk/aws-stepfunctions';
import events = require('@aws-cdk/aws-events');
import { LogRetention } from '@aws-cdk/aws-lambda';
import { LogGroupLogDestination } from '@aws-cdk/aws-apigateway';
import { RetentionDays } from '@aws-cdk/aws-logs';

interface DatabaseProps extends cdk.StackProps {
	dynamoDB: ddb.ITable;
}

export class Api extends cdk.Construct {
	constructor(scope: cdk.Construct, id: string, props: DatabaseProps) {
		super(scope, id);

		const db = props.dynamoDB

		const damageBossFightUuidFunc = new lambda.Function(this, 'DamageBossFightUuid', {
			code: lambda.Code.asset('services/api/damage-boss-fight-uuid'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: Duration.seconds(3),
			environment: {DDB_NAME: db.tableName},
			logRetention: RetentionDays.ONE_MONTH,
			tracing: lambda.Tracing.ACTIVE,
		})

		const damageEncounterIdFunc = new lambda.Function(this, 'DamageEncounterId', {
			code: lambda.Code.asset('services/api/damage-encounter-id'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: Duration.seconds(3),
			environment: {DDB_NAME: db.tableName},
			logRetention: RetentionDays.ONE_MONTH,
			tracing: lambda.Tracing.ACTIVE,
		})

		db.grantReadData(damageBossFightUuidFunc)
		db.grantReadData(damageEncounterIdFunc)
		
		const hostedZone = route53.HostedZone.fromHostedZoneAttributes(this, 'HostedZone', {
			zoneName: 'wowmate.io',
			hostedZoneId: 'Z3LVG9ZF2H87DX',
		});
		const cert = new acm.DnsValidatedCertificate(this, 'Certificate', {
			domainName: 'api.wowmate.io',
			hostedZone,
		});

		const api = new apigateway.LambdaRestApi(this, 'api', {
			handler: damageBossFightUuidFunc,
			proxy: false,
			endpointTypes: [apigateway.EndpointType.REGIONAL],
			domainName: {
				domainName: 'api.wowmate.io',
				certificate: cert,
				securityPolicy: apigateway.SecurityPolicy.TLS_1_2,
			}
		});

		//NOTE: does it make sense to an aaaa record?
		new route53.ARecord(this, 'CustomDomainAliasRecord', {
			zone: hostedZone,
			target: route53.RecordTarget.fromAlias(new targets.ApiGateway(api)),
			recordName: 'api.wowmate.io',
		});

		const basePath = api.root.addResource('api');
		const bossFightPath= basePath.addResource('bossfights');
		const bossFightUuidParam = bossFightPath.addResource('{boss-fight-uuid}');
		const bossFightDamagePath  = bossFightUuidParam.addResource('damage');
		bossFightDamagePath.addMethod('GET')

		const encounterIdPath = basePath.addResource('encounters');
		const encounterIdParam = encounterIdPath.addResource('{encounter-id}');
		const damageEncounterIdIntegration = new apigateway.LambdaIntegration(damageEncounterIdFunc);
		const encounterDamagePath  = encounterIdParam.addResource('damage');
		encounterDamagePath.addMethod('GET', damageEncounterIdIntegration)
	}
}
