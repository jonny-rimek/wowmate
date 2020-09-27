import { Construct, Stack, StackProps } from '@aws-cdk/core';
import { Frontend } from './frontend/frontend';
import { Api } from './api/api';
import { Convert } from './upload/convert';
import { Import } from './upload/import';
import { Summary } from './upload/summary';
import { Presign } from './upload/presign';
import { EtlDashboard } from './upload/etl-dashboard';
import { ApiFrontendDashboard } from './common/api-frontend-dashboard';
import { Vpc } from './common/vpc';
import { Database } from './common/database';
import { Buckets } from './common/buckets';
import { Migrate } from './common/migrate';
import { Partition } from './common/partition';
import { Cloudtrail } from './common/cloudtrail';

export class Wowmate extends Stack {
	constructor(scope: Construct, id: string, props?: StackProps) {
		super(scope, id, props);

		const buckets = new Buckets(this, 'Buckets-')

		const vpc = new Vpc(this, 'Vpc-')

		new Cloudtrail(this, 'Cloudtrail-', {
			uploadBucket: buckets.uploadBucket,
		})

		const db = new Database(this, 'Database-',{
			vpc: vpc.vpc,
			csvBucket: buckets.csvBucket,
		})

		//lambda is exported, metrics could be displayed somewhere
		new Migrate(this, 'Migrate-',{
			dbSecret: db.dbSecret,
			vpc: vpc.vpc,
			dbSecGrp: db.dbSecGrp,
		})

		//lambda is exported, metrics could be displayed somewhere
		new Partition(this, 'Partition-',{
			dbSecret: db.dbSecret,
			vpc: vpc.vpc,
			dbSecGrp: db.dbSecGrp,
			dbEndpoint: db.dbEndpoint,
		})

		const api = new Api(this, 'Api-', {
			dbSecret: db.dbSecret,
			dbEndpoint: db.dbEndpoint,
			vpc: vpc.vpc,
			dbSecGrp: db.dbSecGrp,
		})

		const presign = new Presign(this, 'Presign-', {
			uploadBucket: buckets.uploadBucket,
		})

		const frontend = new Frontend(this, 'Frontend-', {
			api: api.api,
			presignApi: presign.api,
		})

		const convert = new Convert(this, 'Convert-', {
			vpc: vpc.vpc,
			csvBucket: buckets.csvBucket,
			uploadBucket: buckets.uploadBucket,
		})

		const summary = new Summary(this, 'Summary-', {
			vpc: vpc.vpc,
			csvBucket: buckets.csvBucket,
			dbSecGrp: db.dbSecGrp,
			dbSecret: db.dbSecret,
			dbEndpoint: db.dbEndpoint,
		})

		//NOTE: import is a saved keyword
		const importz = new Import(this, 'Import-', {
			vpc: vpc.vpc,
			csvBucket: buckets.csvBucket,
			dbSecGrp: db.dbSecGrp,
			dbSecret: db.dbSecret,
			dbEndpoint: db.dbEndpoint,
			summaryTopic: summary.summaryTopic,
		})

		new ApiFrontendDashboard(this, 'UserFacing-', {
			topDamageLambda: api.topDamageLambda,
			api: api.api,
			s3: frontend.bucket,
			cloudfront: frontend.cloudfront,
		})

		new EtlDashboard(this, 'Etl-', {
			convertLambda: convert.lambda,
			convertQueue: convert.queue,
			convertDLQ: convert.DLQ,
			importLambda: importz.importLambda,
			importQueue: importz.importQueue,
			importDLQ: importz.ImportDLQ,
			summaryLambda: summary.summaryLambda,
			summaryDLQ: summary.summaryDLQ,
			presignLambda: presign.lambda,
			uploadBucket: buckets.uploadBucket,
			csvBucket: buckets.csvBucket,
			presignApi: presign.api,
			cluster: db.cluster,
		})
	}
}