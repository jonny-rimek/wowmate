import { Frontend } from './frontend';
import { Api } from './api';
import { Construct, Stack, StackProps } from '@aws-cdk/core';
import { Convert } from './convert';
import { Import } from './import';
import { Presign } from './presign';

export class Wowmate extends Stack {
	constructor(scope: Construct, id: string, props?: StackProps) {
		super(scope, id, props);

		new Frontend(this, 'Frontend')

		const presign = new Presign(this, 'Presign')

		const api = new Api(this, 'Api')

		new Convert(this, 'Convert', {
			vpc: api.vpc,
			convertBucket: api.bucket,
			uploadBucket: presign.bucket
		})

		new Import(this, 'Import', {
			vpc: api.vpc,
			bucket: api.bucket,
			securityGroup: api.securityGrp,
			secret: api.dbCreds,
		})
	}
}