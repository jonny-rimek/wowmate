import cdk = require('@aws-cdk/core');
import ddb = require('@aws-cdk/aws-dynamodb');
import s3 = require('@aws-cdk/aws-s3');
import lambda = require('@aws-cdk/aws-lambda');
import sfn = require('@aws-cdk/aws-stepfunctions');
import tasks = require('@aws-cdk/aws-stepfunctions-tasks');
import { RemovalPolicy, Duration } from '@aws-cdk/core';
import { BlockPublicAccess } from '@aws-cdk/aws-s3';
// import targets = require('@aws-cdk/aws-events-targets');
import targets = require('@aws-cdk/aws-route53-targets');
import iam = require('@aws-cdk/aws-iam');
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
import { CfnClientCertificate } from '@aws-cdk/aws-apigateway';

export class Frontend extends cdk.Construct {
	constructor(scope: cdk.Construct, id: string) {
		super(scope, id);

		const hostedZone = route53.HostedZone.fromHostedZoneAttributes(this, 'HostedZone', {
			zoneName: 'wowmate.io',
			hostedZoneId: 'Z3LVG9ZF2H87DX',
		});
		
		const cert = new acm.DnsValidatedCertificate(this, 'Certificate', {
			domainName: 'wowmate.io',
			hostedZone,
		});

		const frontendBucket = new s3.Bucket(this, 'FrontendBucket', {
			websiteIndexDocument: 'index.html',
			publicReadAccess: true,
		});

		const distribution = new cloudfront.CloudFrontWebDistribution(this, 'Distribution', {
			originConfigs: [
				{
					customOriginSource: {
						domainName: 'api.wowmate.io',
					},
					behaviors: [{
						pathPattern: '/api/*',
						compress: true,
					}]
				},
				{
					customOriginSource: {
						domainName: frontendBucket.bucketWebsiteDomainName,
						originProtocolPolicy: cloudfront.OriginProtocolPolicy.HTTP_ONLY,
					},
					behaviors : [ {
						isDefaultBehavior: true,
						compress: true,
					}]
				}
			],
			errorConfigurations: [
			],
			aliasConfiguration: {
				names: ['wowmate.io'],
				acmCertRef: cert.certificateArn,
				sslMethod: SSLMethod.SNI,
				securityPolicy: SecurityPolicyProtocol.TLS_V1_2_2018,
			}
		});

		new s3deploy.BucketDeployment(this, 'DeployWebsite', {
			sources: [s3deploy.Source.asset('services/frontend/dist')],
			destinationBucket: frontendBucket,
			distribution,
		});

		new route53.ARecord(this, 'Alias', {
			zone: hostedZone,
			target: route53.RecordTarget.fromAlias(new targets.CloudFrontTarget(distribution)),
		});

		new route53.AaaaRecord(this, 'AliasAAA', {
			zone: hostedZone,
			target: route53.RecordTarget.fromAlias(new targets.CloudFrontTarget(distribution))
		});

	}
}