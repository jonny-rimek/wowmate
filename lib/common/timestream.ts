import cdk = require('@aws-cdk/core');
import { CfnOutput } from '@aws-cdk/core';

interface Props extends cdk.StackProps {
}

export class Database extends cdk.Construct {
	constructor(scope: cdk.Construct, id: string, props?: Props) {
		super(scope, id)
		//not yet implemented as CFN
	}
}