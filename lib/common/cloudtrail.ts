import cdk = require('@aws-cdk/core');
import s3 = require('@aws-cdk/aws-s3');
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
