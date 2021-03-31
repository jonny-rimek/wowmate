import cdk = require('@aws-cdk/core');
import s3 = require('@aws-cdk/aws-s3');
import targets = require('@aws-cdk/aws-route53-targets');
import s3deploy = require('@aws-cdk/aws-s3-deployment');
import cloudfront = require('@aws-cdk/aws-cloudfront');
import route53= require('@aws-cdk/aws-route53');
import acm = require('@aws-cdk/aws-certificatemanager');
import { SSLMethod, SecurityPolicyProtocol } from '@aws-cdk/aws-cloudfront';
import { HttpApi } from '@aws-cdk/aws-apigatewayv2';
import * as origins from "@aws-cdk/aws-cloudfront-origins"

interface Props extends cdk.StackProps {
	api: HttpApi
	presignApi: HttpApi
	hostedZoneId: string
	hostedZoneName: string
}

export class Frontend extends cdk.Construct {
	public readonly cloudfront: cloudfront.Distribution;
	public readonly bucket: s3.Bucket;

	constructor(scope: cdk.Construct, id: string, props: Props) {
		super(scope, id);

		const hostedZone = route53.HostedZone.fromHostedZoneAttributes(this, 'HostedZone', {
			zoneName: props.hostedZoneName,
			hostedZoneId: props.hostedZoneId,
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
			}],
		});
		this.bucket = bucket

		const allowCorsAndQueryString = new cloudfront.OriginRequestPolicy(this, 'AllowCorsAndQueryStringParam', {
			originRequestPolicyName: 'AllowCorsAndQueryStringParam',
			cookieBehavior: cloudfront.OriginRequestCookieBehavior.none(),
			queryStringBehavior: cloudfront.OriginRequestQueryStringBehavior.all(),
			headerBehavior: cloudfront.OriginRequestHeaderBehavior.allowList('origin')
		})

		//make sure enhanced metrics is enabled via the GUI no CF support =(
		//https://console.aws.amazon.com/cloudfront/v2/home#/monitoring
		this.cloudfront = new cloudfront.Distribution(this, 'Distribution', {
			defaultBehavior: {
				origin: new origins.S3Origin(bucket),
				cachePolicy: cloudfront.CachePolicy.CACHING_OPTIMIZED,
				originRequestPolicy: cloudfront.OriginRequestPolicy.CORS_S3_ORIGIN,
			},
            additionalBehaviors: {
				"/api/*": {
					origin: new origins.HttpOrigin(props.api.url!.replace('https://','').replace('/',''), {}),
					cachePolicy: cloudfront.CachePolicy.CACHING_DISABLED,
					originRequestPolicy: allowCorsAndQueryString,
					viewerProtocolPolicy: cloudfront.ViewerProtocolPolicy.HTTPS_ONLY,
					allowedMethods: cloudfront.AllowedMethods.ALLOW_ALL,
					cachedMethods: cloudfront.CachedMethods.CACHE_GET_HEAD_OPTIONS,
				},
				"/presign/*": {
					origin: new origins.HttpOrigin(props.presignApi.url!.replace('https://','').replace('/','')),
					cachePolicy: cloudfront.CachePolicy.CACHING_DISABLED,
					originRequestPolicy: cloudfront.OriginRequestPolicy.CORS_CUSTOM_ORIGIN,
					viewerProtocolPolicy: cloudfront.ViewerProtocolPolicy.HTTPS_ONLY,
					allowedMethods: cloudfront.AllowedMethods.ALLOW_ALL,
				}
			},
            errorResponses: [{
				httpStatus: 404,
				responseHttpStatus: 200,
				responsePagePath: '/index.html',
				ttl: cdk.Duration.seconds(0),
			}],
			certificate: cert,
			domainNames: ["wowmate.io"],
			comment: "wowmate.io frontend, log api and presign api",
		})

		/*
		const distribution = new cloudfront.CloudFrontWebDistribution(this, 'Distribution', {
			originConfigs: [
				{
					customOriginSource: {
						domainName: props.api.url!.replace('https://','').replace('/',''),
						// allowedOriginSSLVersions
						// httpPort
						// httpsPort
						// originHeaders:
						// originKeepaliveTimeout
						// originPath
						// originProtocolPolicy
						// originReadTimeout
					},
					behaviors: [{
						pathPattern: '/api/*',
						compress: true,
						allowedMethods: cloudfront.CloudFrontAllowedMethods.ALL,
						// cachedMethods
						// defaultTtl
						// forwardedValues:
						// maxTtl
						// minTtl
					}]
				},
				{
					customOriginSource: {
						domainName: props.presignApi.url!.replace('https://','').replace('/',''),
					},
					behaviors: [{
						pathPattern: '/presign/*',
						compress: true,
						allowedMethods: cloudfront.CloudFrontAllowedMethods.ALL,
					}],
				},
				{
					customOriginSource: {
						domainName: bucket.bucketWebsiteDomainName,
						originProtocolPolicy: cloudfront.OriginProtocolPolicy.HTTP_ONLY, //doesn't work with https
					},
					behaviors : [{
						isDefaultBehavior: true,
						compress: true,
					}]
				}
			],
			errorConfigurations: [
				{
					errorCode: 404,
					errorCachingMinTtl: 0,
					responseCode: 200,
					responsePagePath: '/index.html',
				}
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

		 */

		const cfnDist = this.cloudfront.node.defaultChild as cloudfront.CfnDistribution;
		cfnDist.addPropertyOverride('DistributionConfig.Origins.0.OriginShield', {
			Enabled: true,
			OriginShieldRegion: 'us-east-1',
		});
		cfnDist.addPropertyOverride('DistributionConfig.Origins.1.OriginShield', {
			Enabled: true,
			OriginShieldRegion: 'us-east-1',
		});

		new s3deploy.BucketDeployment(this, 'DeployWebsite', {
			sources: [s3deploy.Source.asset('services/frontend/dist')],
			destinationBucket: bucket,
			distribution: this.cloudfront,
		});

		new route53.ARecord(this, 'Alias', {
			zone: hostedZone,
			target: route53.RecordTarget.fromAlias(new targets.CloudFrontTarget(this.cloudfront)),
			// recordName
			// ttl: cdk.Duration.minutes(30)
		});

		new route53.AaaaRecord(this, 'AliasAAA', {
			zone: hostedZone,
			target: route53.RecordTarget.fromAlias(new targets.CloudFrontTarget(this.cloudfront))
		});
	}
}
