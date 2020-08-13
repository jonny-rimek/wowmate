import cdk = require('@aws-cdk/core');
import targets = require('@aws-cdk/aws-route53-targets');
import ec2 = require('@aws-cdk/aws-ec2');
import rds = require('@aws-cdk/aws-rds');
import ecs = require('@aws-cdk/aws-ecs');
import elbv2 = require('@aws-cdk/aws-elasticloadbalancingv2');
import route53= require('@aws-cdk/aws-route53');
import ecsPatterns = require('@aws-cdk/aws-ecs-patterns');
import * as lambda from '@aws-cdk/aws-lambda';
import apigateway = require('@aws-cdk/aws-apigateway');
import acm = require('@aws-cdk/aws-certificatemanager');
import { BaseLoadBalancer } from '@aws-cdk/aws-elasticloadbalancingv2';
import s3 = require('@aws-cdk/aws-s3');
import { RemovalPolicy, Duration } from '@aws-cdk/core';

export class V2 extends cdk.Construct {
	constructor(scope: cdk.Construct, id: string) {
		super(scope, id)
		
		const uploadBucket = new s3.Bucket(this, 'Upload', {
			removalPolicy: RemovalPolicy.DESTROY,
			blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL,
		})
/* 
		//TODO: create bucket and pass to lambda
		const presignLambda = new lambda.Function(this, 'PresignLambda', {
			runtime: lambda.Runtime.NODEJS_12_X,
			code: lambda.Code.asset('services/presign'),
			handler: 'index.handler',
			environment: {BUCKET_NAME: uploadBucket.bucketName},
		});

		uploadBucket.grantWrite(presignLambda);
		const hostedZone = route53.HostedZone.fromHostedZoneAttributes(this, 'HostedZone', {
			zoneName: 'wowmate.io',
			hostedZoneId: 'Z3LVG9ZF2H87DX',
		});

		const cert = new acm.DnsValidatedCertificate(this, 'Certificate', {
			domainName: 'presign.wowmate.io',
			hostedZone: hostedZone,
		});

		const api = new apigateway.LambdaRestApi(this, 'PresignApi', {
			handler: presignLambda,
			proxy: false,
			endpointTypes: [apigateway.EndpointType.REGIONAL],
			domainName: {
				domainName: 'presign.wowmate.io',
				certificate: cert,
				securityPolicy: apigateway.SecurityPolicy.TLS_1_2,
			}
		});
		api.root.addMethod('POST');

		//NOTE: does it make sense to an aaaa record?
		new route53.ARecord(this, 'CustomDomainAliasRecord', {
			zone: hostedZone,
			target: route53.RecordTarget.fromAlias(new targets.ApiGateway(api)),
			recordName: 'presign.wowmate.io',
		});

		new cdk.CfnOutput(this, 'HTTP API Url', {
			value: api.url ?? 'Something went wrong with the deploy'
		});
		const vpc = new ec2.Vpc(this, 'WowmateVpc', {
			natGateways: 1,
		});

		/*
		const postgres = new rds.DatabaseInstance(this, 'Postgres', {
			engine: rds.DatabaseInstanceEngine.postgres({
				version: rds.PostgresEngineVersion.VER_11_7,
			}),
			instanceType: ec2.InstanceType.of(ec2.InstanceClass.BURSTABLE2, ec2.InstanceSize.MICRO),
			masterUsername: 'postgres',
			vpc,
			// vpcPlacement: { subnetType: ec2.SubnetType.PUBLIC },
			//NOTE: remove in production
			removalPolicy: cdk.RemovalPolicy.DESTROY,
			deletionProtection: false,
		})
		postgres.connections.allowFromAnyIpv4(ec2.Port.tcp(5432))

		//IMPROVE: add https redirect
		//need to define the cluster seperately and in it the VPC i think
		const loadBalancedFargateService = new ecsPatterns.ApplicationLoadBalancedFargateService(this, 'Service', {
			vpc,
			// domainName: 'api.wowmate.io',
			// domainZone: hostedZone,
			memoryLimitMiB: 512,
			// protocol: elbv2.ApplicationProtocol.HTTPS,
			cpu: 256,
			desiredCount: 1,
			publicLoadBalancer: true,
			platformVersion: ecs.FargatePlatformVersion.VERSION1_4,
			taskImageOptions: {
				image: ecs.ContainerImage.fromAsset('services/api'),
			},
		});
		// */
	}
}
