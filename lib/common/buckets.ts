import cdk = require('@aws-cdk/core');
import s3 = require('@aws-cdk/aws-s3');
import { RemovalPolicy } from '@aws-cdk/core';

export class Buckets extends cdk.Construct {
	public readonly csvBucket: s3.Bucket;
	public readonly uploadBucket: s3.Bucket;

	constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
		super(scope, id)

		this.uploadBucket = new s3.Bucket(this, 'Upload', {
			removalPolicy: RemovalPolicy.DESTROY,
			//I think presigning a link didn't work with block all
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

		this.csvBucket = new s3.Bucket(this, 'CSV', {
			removalPolicy: cdk.RemovalPolicy.DESTROY,
			blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL,
			metrics: [{ //enables advanced s3metrics
				id: 'metric',
			}]
		})
	}
}