import cdk = require('@aws-cdk/core');
import { CfnOutput, cfnTagToCloudFormation } from '@aws-cdk/core';
import timestream = require('@aws-cdk/aws-timestream');

interface Props extends cdk.StackProps {
}

export class Timestream extends cdk.Construct {
	public readonly timestreamArn: string;

	constructor(scope: cdk.Construct, id: string, props?: Props) {
		super(scope, id)

		const db = new timestream.CfnDatabase(this, "db", {
			databaseName: "wowmate-analytics",
		})

		const table = new timestream.CfnTable(this, "table", {
			databaseName: "wowmate-analytics",
			tableName: "combatlogs",
			retentionProperties: {
				MemoryStoreRetentionPeriodInHours: "4380", //6month, can't ingest data older, than this value
				//NOTE: should be lowered in the future to save cost as memory storage is very expensive,
				//		won't be a problem now as I'm not storing a lot of data
				//		users won't be able to upload logs older than this value
				MagneticStoreRetentionPeriodInDays: "73000", //200years = max
				//data will be deleted after 200years
			}
		})

		table.node.addDependency(db)

		this.timestreamArn = table.attrArn

		new CfnOutput(this, 'HttpApiEndpoint', {
			value: table.attrArn,
		})
	}
}