import cdk = require('@aws-cdk/core');
import s3 = require('@aws-cdk/aws-s3');
import targets = require('@aws-cdk/aws-route53-targets');
import iam = require('@aws-cdk/aws-iam');
import cloudtrail = require('@aws-cdk/aws-cloudtrail');
import apigateway = require('@aws-cdk/aws-apigateway');
import s3deploy = require('@aws-cdk/aws-s3-deployment');
import cloudfront = require('@aws-cdk/aws-cloudfront');
import route53= require('@aws-cdk/aws-route53');
import acm = require('@aws-cdk/aws-certificatemanager');
import { SSLMethod, SecurityPolicyProtocol } from '@aws-cdk/aws-cloudfront';
import { HttpApi } from '@aws-cdk/aws-apigatewayv2';

interface Props extends cdk.StackProps {
	api: HttpApi
	presignApi: apigateway.LambdaRestApi
}

export class Frontend extends cdk.Construct {
	public readonly cloudfront: cloudfront.CloudFrontWebDistribution;
	public readonly bucket: s3.Bucket;

	constructor(scope: cdk.Construct, id: string, props: Props) {
		super(scope, id);

		const hostedZone = route53.HostedZone.fromHostedZoneAttributes(this, 'HostedZone', {
			zoneName: 'wowmate.io',
			hostedZoneId: 'Z3LVG9ZF2H87DX',
		});

		const cert = new acm.DnsValidatedCertificate(this, 'Certificate', {
			domainName: 'wowmate.io',
			hostedZone,
		});

		//if the bucket would be private we can't make it a website bucket
		//which results in strange behaviour
		//e.g. you have to call the file exactly /page.html simply /page won't work
		//the disadvantage is that people could access the bucket directly which would
		//result in a slower website, but more importantly it could be used as a denial of wallet
		//as it doesn't get cached and data transfer out is billed every time. as the bucket name
		//is random and not having smart redirect seems worse, I'm going with a public website bucket
		const bucket = new s3.Bucket(this, 'Bucket', {
			websiteIndexDocument: 'index.html',
			publicReadAccess: true,
			removalPolicy: cdk.RemovalPolicy.DESTROY,
			metrics: [{
				id: 'metric',
			}]
		});
		this.bucket = bucket

		const distribution = new cloudfront.CloudFrontWebDistribution(this, 'Distribution', {
			originConfigs: [
				{
					customOriginSource: {
						domainName: props.api.url!.replace('https://','').replace('/',''),
					},
					behaviors: [{
						pathPattern: '/api/*',
						compress: true,
						// allowedMethods
						// cachedMethods
						// defaultTtl
						// forwardedValues
						// maxTtl
						// minTtl
					}]
				},
				{
					customOriginSource: {
						domainName: props.presignApi.url!.replace('https://','').replace('/',''),
					},
					behaviors: [{
						pathPattern: '/presign',
						compress: true,
						allowedMethods: cloudfront.CloudFrontAllowedMethods.ALL,
					}],
				},
				{
					customOriginSource: {
						domainName: bucket.bucketWebsiteDomainName,
						originProtocolPolicy: cloudfront.OriginProtocolPolicy.HTTP_ONLY,
					},
					behaviors : [ {
						isDefaultBehavior: true,
						compress: true,
					}]
				}
			],
			errorConfigurations: [
				// {
				// errorCode
				// errorCachingMinTtl
				// responseCode
				// responsePagePath
				// }
			],
			aliasConfiguration: {
				names: ['wowmate.io'],
				acmCertRef: cert.certificateArn,
				sslMethod: SSLMethod.SNI,
				securityPolicy: SecurityPolicyProtocol.TLS_V1_2_2018,
			},
			// comment
			// defaultRootObject
			// enableIpV6
			// httpVersion
			// loggingConfig
			// priceClass
			// viewerCertificate
			// viewerProtocolPolicy //redirect to https
			// webACLId //WAF config
		});
		this.cloudfront = distribution

		new s3deploy.BucketDeployment(this, 'DeployWebsite', {
			sources: [s3deploy.Source.asset('services/frontend/dist')],
			destinationBucket: bucket,
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
