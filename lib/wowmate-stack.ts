import { Frontend } from './frontend';
import { Api } from './api';
import { Construct, Stack, StackProps } from '@aws-cdk/core';
import { Convert } from './convert';
import { Import } from './import';
import { Presign } from './presign';
import { EtlDashboard } from './etl-dashboard';
import { FrontendDashboard } from './frontend-dashboard';

export class Wowmate extends Stack {
	constructor(scope: Construct, id: string, props?: StackProps) {
		super(scope, id, props);

		//TODO: add - at the end of each name for better readability
		const api = new Api(this, 'Api')

		const frontend = new Frontend(this, 'Frontend', {
			api: api.api,
		})

		const presign = new Presign(this, 'Presign')

		const convert = new Convert(this, 'Convert', {
			vpc: api.vpc,
			csvBucket: api.bucket,
			uploadBucket: presign.bucket
		})

		//NOTE: import is a saved keyword
		const importz = new Import(this, 'Import', {
			vpc: api.vpc,
			bucket: api.bucket,
			securityGroup: api.securityGrp,
			secret: api.dbCreds,
			dbEndpoint: api.dbEndpoint,
		})

		new FrontendDashboard(this, 'Namingthingsishard', {
			topDamageLambda: api.lambda,
			api: api.api,
			s3: frontend.bucket,
			cloudfront: frontend.cloudfront,
		})

		new EtlDashboard(this, 'ETL', {
			convertLambda: convert.lambda,
			convertQueue: convert.queue,
			convertDLQ: convert.dlq,
			importLambda: importz.importLambda,
			importQueue: importz.queue,
			importDLQ: importz.dlq,
			//sumaryLambda
			//summaryDlq
			presignLambda: presign.lambda,
			uploadBucket: presign.bucket,
			presignApiGateway: presign.apiGateway,
			cluster: api.cluster,
		})
	}
}