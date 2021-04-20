import cdk = require('@aws-cdk/core');
import s3 = require('@aws-cdk/aws-s3');
import targets = require('@aws-cdk/aws-route53-targets');
import s3deploy = require('@aws-cdk/aws-s3-deployment');
import cloudfront = require('@aws-cdk/aws-cloudfront');
import route53= require('@aws-cdk/aws-route53');
import acm = require('@aws-cdk/aws-certificatemanager');
import * as origins from "@aws-cdk/aws-cloudfront-origins"
import * as iam from "@aws-cdk/aws-iam"

interface Props extends cdk.StackProps {
	hostedZoneId: string
	hostedZoneName: string
	domainName: string
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
			domainName: props.domainName,
			hostedZone,
		});

		//if the bucket would be private we can't make it a website bucket
		//which results in strange behaviour
		//e.g. you have to call the file exactly /page.html simply /page won't work
		//the disadvantage is that people could access the bucket directly which would
		//result in a slower website, but more importantly it could be used as a denial of wallet
		//as it doesn't get cached and data transfer out is billed every time. as the bucket name
		//is random and not having smart redirect seems worse, I'm going with a public website bucket
		this.bucket = new s3.Bucket(this, 'Bucket', {
			websiteIndexDocument: 'index.html',
			publicReadAccess: true,
			removalPolicy: cdk.RemovalPolicy.DESTROY,
			metrics: [{
				id: 'metric',
			}],
			encryption: s3.BucketEncryption.S3_MANAGED,
			// encrypting a bucket that can be publicly read is probably not the most useful thing to do
			// but it's a best practice
		});

		const cfnBucket = this.bucket.node.defaultChild as s3.CfnBucket
		cfnBucket.cfnOptions.metadata = {
			cfn_nag: {
				rules_to_suppress: [
					{
						id: 'W35',
						reason: "this is a website bucket, so there is no point tracking access to it",
					},
					{
						id: 'W41',
						reason: "this is a website bucket, it needs to be public, so there is no point in encrypting it",
					},
				]
			}
		}

		const cfnBucketPolicy = this.bucket.policy?.node.defaultChild as iam.CfnPolicy
		cfnBucketPolicy.cfnOptions.metadata = {
			cfn_nag: {
				rules_to_suppress: [
					{
						id: 'F16',
						reason: "this is a website bucket, it needs to be public",
					},
				]
			}
		}

		//make sure enhanced metrics is enabled via the GUI no CF support =(
		//https://console.aws.amazon.com/cloudfront/v2/home#/monitoring
		this.cloudfront = new cloudfront.Distribution(this, 'Distribution', {
			defaultBehavior: {
				origin: new origins.S3Origin(this.bucket),
				cachePolicy: cloudfront.CachePolicy.CACHING_OPTIMIZED,
				originRequestPolicy: cloudfront.OriginRequestPolicy.CORS_S3_ORIGIN,
				viewerProtocolPolicy: cloudfront.ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
			},
            errorResponses: [{
				httpStatus: 404,
				responseHttpStatus: 200,
				responsePagePath: '/index.html',
				ttl: cdk.Duration.seconds(0),
			}],
			certificate: cert,
			domainNames: [props.domainName],
			comment: "wowmate.io frontend",
		})

		const cfnDist = this.cloudfront.node.defaultChild as cloudfront.CfnDistribution;
		cfnDist.addPropertyOverride('DistributionConfig.Origins.0.OriginShield', {
			Enabled: true,
			OriginShieldRegion: 'us-east-1',
		});

		new s3deploy.BucketDeployment(this, 'DeployWebsite', {
			sources: [s3deploy.Source.asset('services/frontend/dist')],
			destinationBucket: this.bucket,
			distribution: this.cloudfront,
		});

		new route53.ARecord(this, 'Alias', {
			zone: hostedZone,
			target: route53.RecordTarget.fromAlias(new targets.CloudFrontTarget(this.cloudfront)),
			recordName: props.domainName,
		});

		new route53.AaaaRecord(this, 'AliasAAA', {
			zone: hostedZone,
			target: route53.RecordTarget.fromAlias(new targets.CloudFrontTarget(this.cloudfront)),
			recordName: props.domainName,
		});
	}
}
