import { Construct, Stack, StackProps } from '@aws-cdk/core';
import { Frontend } from './frontend/frontend';
import { Api } from './api/api';
import { Convert } from './upload/convert';
import { Import } from './upload/import';
import { Summary } from './upload/summary';
import { Presign } from './upload/presign';
import { EtlDashboard } from './upload/etl-dashboard';
import { ApiFrontendDashboard } from './common/api-frontend-dashboard';
import { Common } from './common/common';
import { Migrate } from './common/migrate';
import { Partition } from './common/partition';

export class Wowmate extends Stack {
	constructor(scope: Construct, id: string, props?: StackProps) {
		super(scope, id, props);

		const common = new Common(this, 'Common-')

		//lambda is exported, metrics could be displayed somewhere
		new Migrate(this, 'Migrate-',{
			dbSecret: common.dbSecret,
			vpc: common.vpc,
			dbSecGrp: common.dbSecGrp,
		})

		//lambda is exported, metrics could be displayed somewhere
		new Partition(this, 'Partition-',{
			dbSecret: common.dbSecret,
			vpc: common.vpc,
			dbSecGrp: common.dbSecGrp,
			dbEndpoint: common.dbEndpoint,
		})

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

		const summary = new Summary(this, 'Import-', {
			vpc: common.vpc,
			csvBucket: common.csvBucket,
			dbSecGrp: common.dbSecGrp,
			dbSecret: common.dbSecret,
			dbEndpoint: common.dbEndpoint,
		})

		//NOTE: import is a saved keyword
		const importz = new Import(this, 'Import-', {
			vpc: common.vpc,
			csvBucket: common.csvBucket,
			dbSecGrp: common.dbSecGrp,
			dbSecret: common.dbSecret,
			dbEndpoint: common.dbEndpoint,
			summaryLambda: summary.summaryLambda,
		})

		new ApiFrontendDashboard(this, 'ApiFrontendDashboard-', {
			topDamageLambda: api.topDamageLambda,
			api: api.api,
			s3: frontend.bucket,
			cloudfront: frontend.cloudfront,
		})

		new EtlDashboard(this, 'EtlDashboard-', {
			convertLambda: convert.lambda,
			convertQueue: convert.queue,
			convertDLQ: convert.DLQ,
			importLambda: importz.importLambda,
			importQueue: importz.importQueue,
			importDLQ: importz.ImportDLQ,
			summaryLambda: summary.summaryLambda,
			summaryDLQ: summary.summaryDLQ,
			presignLambda: presign.lambda,
			uploadBucket: common.uploadBucket,
			csvBucket: common.csvBucket,
			presignApiGateway: presign.apiGateway,
			cluster: common.cluster,
		})
	}
}