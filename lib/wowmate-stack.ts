import { Frontend } from './frontend';
import { V2 } from './v2';
import { Construct, Stage, Stack, StackProps, StageProps, SecretValue } from '@aws-cdk/core';

export class Wowmate extends Stack {
	constructor(scope: Construct, id: string, props?: StageProps) {
		super(scope, id, props);

		new V2(this, 'V2')
		new Frontend(this, 'frontend')
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