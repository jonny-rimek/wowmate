import cdk = require('@aws-cdk/core');
import s3 = require('@aws-cdk/aws-s3');
import { RemovalPolicy } from '@aws-cdk/core';

export class Buckets extends cdk.Construct {
	// public readonly csvBucket: s3.Bucket;
	public readonly uploadBucket: s3.Bucket;
	public readonly accessLogBucket: s3.Bucket;

	constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
		super(scope, id)

        this.accessLogBucket = new s3.Bucket(this, 'S3AccessLog', {
			blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL,
			encryption: s3.BucketEncryption.S3_MANAGED,
			serverAccessLogsPrefix: "accessLogBucket",
		})

		this.uploadBucket = new s3.Bucket(this, 'Upload', {
			serverAccessLogsPrefix: "uploadBucket",
			serverAccessLogsBucket: this.accessLogBucket,
			blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL,
		    encryption: s3.BucketEncryption.S3_MANAGED,
			removalPolicy: RemovalPolicy.DESTROY, //TODO: remove
			cors: [
				{
					allowedOrigins: [
						"*", //setting only to wowmate.io didnt work, maybe agw domain
					],
					allowedMethods: [
						s3.HttpMethods.POST,
						s3.HttpMethods.HEAD,
					],
					allowedHeaders: [
						"*", //check if we can remove this
					],
				}
			],
			metrics: [{
				id: 'metric',
			}],
			lifecycleRules: [{
			    //deletes error content
				//these are failed uploads only kept for debugging, but mostly faulty uploads
				expiration: cdk.Duration.days(14),
				prefix: 'error/',
			},
			{
			    //deletes partial uploads that failed
				abortIncompleteMultipartUploadAfter: cdk.Duration.days(1),
			},
			{
				//changes the storage tier
				//I never really need the files again, I could just save em to glacier
				//but I don't want to prematurely optimize, as it is not as ez to get data out of glacier again
				prefix: 'upload/',
			    transitions: [{
			    	storageClass: s3.StorageClass.INFREQUENT_ACCESS,
					transitionAfter: cdk.Duration.days(30)//minimum duration
				}]
			}]
		})

		const cfnAccessLogBucket = this.accessLogBucket.node.defaultChild as s3.CfnBucket
		cfnAccessLogBucket.cfnOptions.metadata = {
			cfn_nag: {
				rules_to_suppress: [
					{
						id: 'W51',
						reason: "suppress 'S3 bucket should likely have a bucket policy', because this bucket doesn't need a policy",
					},
				]
			}
		}
		const cfnBucket = this.uploadBucket.node.defaultChild as s3.CfnBucket
		cfnBucket.cfnOptions.metadata = {
			cfn_nag: {
				rules_to_suppress: [
					{
						id: 'W51',
						reason: "suppress 'S3 bucket should likely have a bucket policy', because this bucket doesn't need a policy",
					},
				]
			}
		}

	}
}