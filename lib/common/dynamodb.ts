import cdk = require('@aws-cdk/core');
import * as dynamodb from '@aws-cdk/aws-dynamodb';

interface Props extends cdk.StackProps {
	// csvBucket: s3.Bucket
}

export class DynamoDB extends cdk.Construct {
	public readonly table: dynamodb.Table;

	constructor(scope: cdk.Construct, id: string, props?: Props) {
		super(scope, id)

		this.table = new dynamodb.Table(this, 'table', {
			partitionKey: {
				name: 'pk',
				type: dynamodb.AttributeType.STRING,
			},
			sortKey: {
				name: 'sk',
				type: dynamodb.AttributeType.STRING,
			},
			billingMode: dynamodb.BillingMode.PAY_PER_REQUEST, //on demand
			pointInTimeRecovery: true,
			removalPolicy: cdk.RemovalPolicy.DESTROY, //NOTE: remove in prod
			// activate if I need dynamoDB streams
			// stream: dynamodb.StreamViewType.NEW_AND_OLD_IMAGES,
		})

		// this.table.addGlobalSecondaryIndex({
		// 	indexName: 'GSI1',
		// 	partitionKey: {
		// 		name: 'GSI1PK',
		// 		type: dynamodb.AttributeType.STRING,
		// 	},
		// 	sortKey: {
		// 		name: 'GSI1SK',
		// 		type: dynamodb.AttributeType.STRING,
		// 	},
		// 	projectionType: dynamodb.ProjectionType.ALL,
		// 	nonKeyAttributes: ["",""],
		// })
	}
}