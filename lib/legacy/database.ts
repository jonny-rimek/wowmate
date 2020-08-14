import cdk = require('@aws-cdk/core');
import ddb = require('@aws-cdk/aws-dynamodb');
import { RemovalPolicy } from '@aws-cdk/core';

export class Database extends cdk.Construct {
	public readonly dynamoDB: ddb.Table;

	constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
		super(scope, id);

		const db = new ddb.Table(this, 'DDB', {
			partitionKey: { name: 'pk', type: ddb.AttributeType.STRING },
			sortKey: {name: 'sk', type: ddb.AttributeType.STRING},
			removalPolicy: RemovalPolicy.DESTROY,
            billingMode: ddb.BillingMode.PAY_PER_REQUEST
		})

		db.addGlobalSecondaryIndex({
			indexName: 'GSI1',
			partitionKey: {name: 'gsi1pk', type: ddb.AttributeType.NUMBER},
			sortKey: {name: 'gsi1sk', type: ddb.AttributeType.NUMBER}
		})
		this.dynamoDB = db;
	}
}
