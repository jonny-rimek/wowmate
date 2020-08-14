import { Frontend } from './frontend';
import { Api } from './api';
import { Construct, Stack, StackProps } from '@aws-cdk/core';
import { Converter } from './converter';
import { Vpc } from './vpc';
import { Postgres } from './postgres';

export class Wowmate extends Stack {
	constructor(scope: Construct, id: string, props?: StackProps) {
		super(scope, id, props);

		const vpc = new Vpc(this, 'Vpc')

		new Postgres(this, 'Postgres', {
			vpc: vpc.vpc
		})

		new Api(this, 'Api', {
			vpc: vpc.vpc
		})

		new Converter(this, 'Converter', {
			vpc: vpc.vpc
		})
		// new Frontend(this, 'frontend')
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