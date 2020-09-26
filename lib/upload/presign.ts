import cdk = require('@aws-cdk/core');
import targets = require('@aws-cdk/aws-route53-targets');
import route53= require('@aws-cdk/aws-route53');
import * as lambda from '@aws-cdk/aws-lambda';
import apigateway = require('@aws-cdk/aws-apigateway');
import acm = require('@aws-cdk/aws-certificatemanager');
import s3 = require('@aws-cdk/aws-s3');
import { RemovalPolicy } from '@aws-cdk/core';
import * as cloudtrail from '@aws-cdk/aws-cloudtrail';
import { ReadWriteType } from '@aws-cdk/aws-cloudtrail';

interface Props extends cdk.StackProps {
	uploadBucket: s3.Bucket;
}

export class Presign extends cdk.Construct {
	public readonly lambda: lambda.Function
	public readonly apiGateway: apigateway.LambdaRestApi


	constructor(scope: cdk.Construct, id: string, props: Props) {
		super(scope, id)
		
		//EXRACT
		const trail = new cloudtrail.Trail(this, '-Cloudtrail', {
			managementEvents: ReadWriteType.WRITE_ONLY,
			sendToCloudWatchLogs: true,
			// cloudWatchLogsRetention:^
		});

		trail.addS3EventSelector([{
			bucket: props.uploadBucket, 
		}]);
		
		const presignLambda = new lambda.Function(this, '-PresignLambda', {
			runtime: lambda.Runtime.NODEJS_12_X,
			code: lambda.Code.fromAsset('services/presign'),
			handler: 'index.handler',
			environment: {BUCKET_NAME: props.uploadBucket.bucketName},
		});
		this.lambda = presignLambda

		props.uploadBucket.grantPut(presignLambda);

		const hostedZone = route53.HostedZone.fromHostedZoneAttributes(this, '-HostedZone', {
			zoneName: 'wowmate.io',
			hostedZoneId: 'Z3LVG9ZF2H87DX',
		});

		const cert = new acm.DnsValidatedCertificate(this, '-Certificate', {
			domainName: 'presign.wowmate.io',
			hostedZone: hostedZone,
		});

		this.apiGateway = new apigateway.LambdaRestApi(this, '-PresignApi', {
			handler: presignLambda,
			proxy: false,
			endpointTypes: [apigateway.EndpointType.REGIONAL],
			domainName: {
				domainName: 'presign.wowmate.io',
				certificate: cert,
				securityPolicy: apigateway.SecurityPolicy.TLS_1_2,
			},
			//TODO: test if i need cors after I activated CORS on the bucket
			defaultCorsPreflightOptions: {
				allowOrigins: apigateway.Cors.ALL_ORIGINS,
			}
		});
		const presign = this.apiGateway.root.addResource('presign');
		presign.addMethod('POST');

		//NOTE: does it make sense to an aaaa record?
		new route53.ARecord(this, 'CustomDomainAliasRecord', {
			zone: hostedZone,
			target: route53.RecordTarget.fromAlias(new targets.ApiGateway(this.apiGateway)),
			recordName: 'presign.wowmate.io',
		});
		new cdk.CfnOutput(this, 'HTTP API Url', {
			value: this.apiGateway.url ?? 'Something went wrong with the deploy'
		});
	}
}
