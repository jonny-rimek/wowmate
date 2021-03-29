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
				MemoryStoreRetentionPeriodInHours: "1", //can't ingest data older, than this value
				// 26€/GB/month
				MagneticStoreRetentionPeriodInDays: "365", //data will be deleted after 1year
				// 0,03€/GB/month = 30€/TB/month
				//73000 = 200years = max
			}
		})

		table.node.addDependency(db)

		this.timestreamArn = table.attrArn

		new CfnOutput(this, 'HttpApiEndpoint', {
			value: table.attrArn,
		})
	}
}