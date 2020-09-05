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

export class Presign extends cdk.Construct {
	public readonly bucket: s3.Bucket;

	constructor(scope: cdk.Construct, id: string) {
		super(scope, id)
		
		const uploadBucket = new s3.Bucket(this, 'Upload', {
			removalPolicy: RemovalPolicy.DESTROY,
			// blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL,
			cors: [
				{
					allowedOrigins: [
						"*", //if I only do wowmate.io it will not work locally
					],
					allowedMethods: [
						s3.HttpMethods.POST,
						s3.HttpMethods.GET, //should be removable
						s3.HttpMethods.PUT,//should be removable
						s3.HttpMethods.HEAD,
					],
					allowedHeaders: [
						"*", //dunno
					],
				}
			],
			metrics: [{
				id: 'metric',
			}]
		})
		this.bucket = uploadBucket

		const trail = new cloudtrail.Trail(this, 'Cloudtrail', {
			managementEvents: ReadWriteType.WRITE_ONLY,

			sendToCloudWatchLogs: true,
			// cloudWatchLogsRetention:^
		});

		trail.addS3EventSelector([{
			bucket: uploadBucket, 
		}]);
		
		//TODO: create bucket and pass to lambda
		const presignLambda = new lambda.Function(this, 'PresignLambda', {
			runtime: lambda.Runtime.NODEJS_12_X,
			code: lambda.Code.asset('services/presign'),
			handler: 'index.handler',
			environment: {BUCKET_NAME: uploadBucket.bucketName},
		});

		// uploadBucket.grantWrite(presignLambda);
		uploadBucket.grantPut(presignLambda);
		// uploadBucket.grantReadWrite(presignLambda);
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
			},
			//TODO: test if i need cors after I activated CORS on the bucket
			defaultCorsPreflightOptions: {
				allowOrigins: apigateway.Cors.ALL_ORIGINS,
			}
		});
		const presign = api.root.addResource('presign');
		presign.addMethod('POST');

		//NOTE: does it make sense to an aaaa record?
		new route53.ARecord(this, 'CustomDomainAliasRecord', {
			zone: hostedZone,
			target: route53.RecordTarget.fromAlias(new targets.ApiGateway(api)),
			recordName: 'presign.wowmate.io',
		});
		new cdk.CfnOutput(this, 'HTTP API Url', {
			value: api.url ?? 'Something went wrong with the deploy'
		});
	}
}
