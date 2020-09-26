import { Frontend } from './frontend';
import { Api } from './api';
import { Construct, Stack, StackProps } from '@aws-cdk/core';
import { Convert } from './convert';
import { Import } from './import';
import { Presign } from './presign';
import { EtlDashboard } from './etl-dashboard';
import { FrontendDashboard } from './frontend-dashboard';
import { Common } from './common';

export class Wowmate extends Stack {
	constructor(scope: Construct, id: string, props?: StackProps) {
		super(scope, id, props);

		const common = new Common(this, 'Common-')

		const api = new Api(this, 'Api-', {
			dbSecret: common.dbSecret,
			dbEndpoint: common.dbEndpoint,
			vpc: common.vpc,
			dbSecGrp: common.dbSecGrp,
		})

		const frontend = new Frontend(this, 'Frontend-', {
			api: api.api,
		})

		const presign = new Presign(this, 'Presign-', {
			uploadBucket: common.uploadBucket,
		})

		const convert = new Convert(this, 'Convert-', {
			vpc: common.vpc,
			csvBucket: common.csvBucket,
			uploadBucket: common.uploadBucket,
		})

		//NOTE: import is a saved keyword
		const importz = new Import(this, 'Import-', {
			vpc: common.vpc,
			csvBucket: common.csvBucket,
			securityGroup: common.dbSecGrp,
			dbSecret: common.dbSecret,
			dbEndpoint: common.dbEndpoint,
		})

		new FrontendDashboard(this, 'ApiFrontendDashboard-', {
			topDamageLambda: api.lambda,
			api: api.api,
			s3: frontend.bucket,
			cloudfront: frontend.cloudfront,
		})

		new EtlDashboard(this, 'EtlDashboard-', {
			convertLambda: convert.lambda,
			convertQueue: convert.queue,
			convertDLQ: convert.DLQ,
			importLambda: importz.importLambda,
			importQueue: importz.queue,
			importDLQ: importz.DLQ,
			summaryLambda: importz.summaryLambda,
			summaryDLQ: importz.summaryDLQ,
			presignLambda: presign.lambda,
			uploadBucket: common.uploadBucket,
			csvBucket: common.csvBucket,
			presignApiGateway: presign.apiGateway,
			cluster: common.cluster,
		})
	}
}