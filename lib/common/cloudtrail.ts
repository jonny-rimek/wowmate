import cdk = require('@aws-cdk/core');
import s3 = require('@aws-cdk/aws-s3');
import * as cloudtrail from '@aws-cdk/aws-cloudtrail';
import {ReadWriteType} from '@aws-cdk/aws-cloudtrail';

interface Props extends cdk.StackProps {
	uploadBucket: s3.Bucket;
}

export class Cloudtrail extends cdk.Construct {

	constructor(scope: cdk.Construct, id: string, props: Props) {
		super(scope, id)

		/*
		//https://github.com/aws/aws-cdk/issues/13966
		//can't create a bucket only trail
		const trail = new cloudtrail.Trail(this, 't', {
			managementEvents: ReadWriteType.NONE,
			sendToCloudWatchLogs: true,
		});

		trail.addS3EventSelector([{
			bucket: props.uploadBucket,
			objectPrefix: "upload",
		}], {
			includeManagementEvents: false,
			readWriteType: ReadWriteType.WRITE_ONLY,
		});
		 */

	}
}
