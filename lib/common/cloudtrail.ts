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

export class Cloudtrail extends cdk.Construct {

	constructor(scope: cdk.Construct, id: string, props: Props) {
		super(scope, id)
		
		const trail = new cloudtrail.Trail(this, 'Cloudtrail', {
			managementEvents: ReadWriteType.WRITE_ONLY,
			sendToCloudWatchLogs: true,
			// cloudWatchLogsRetention:^
		});

		trail.addS3EventSelector([{
			bucket: props.uploadBucket, 
		}]);
	}
}
