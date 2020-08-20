import { Frontend } from './frontend';
import { Api } from './api';
import { Construct, Stack, StackProps } from '@aws-cdk/core';
import { Convert } from './convert';
import { Import } from './import';

export class Wowmate extends Stack {
	constructor(scope: Construct, id: string, props?: StackProps) {
		super(scope, id, props);

		const api = new Api(this, 'Api')

		const convert = new Convert(this, 'Converter', {
			vpc: api.vpc
		})
		new Frontend(this, 'frontend')

		new Import(this, 'Import', {
			vpc: api.vpc,
			bucket: convert.bucket,
		})
	}
}
/* 
export class Wowmate extends cdk.Stack {
	constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
		super(scope, id, props);

		// const db = new Database(this, 'Database')

		// new Api(this, 'Api', {
		// 	dynamoDB: db.dynamoDB,
		// })

		// new Frontend(this, 'Frontend')
		
		// new Upload(this, 'Upload', {
		// 	dynamoDB: db.dynamoDB,
		// })

		new V2(this, 'V2')
		new Frontend(this, 'frontend')
	}
}
 */