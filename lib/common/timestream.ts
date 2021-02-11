import cdk = require('@aws-cdk/core');
import { CfnOutput } from '@aws-cdk/core';
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
		})

		table.node.addDependency(db)

		this.timestreamArn = table.attrArn
	}
}