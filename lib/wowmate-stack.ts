import cdk = require('@aws-cdk/core');
import { Database } from './database';
import { Upload } from './upload';
import { Api } from './api';
import { Frontend } from './frontend';
import { V2 } from './v2';

export class Wowmate extends cdk.Stack {
	constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
		super(scope, id, props);

		const db = new Database(this, 'Database')

		new Api(this, 'Api', {
			dynamoDB: db.dynamoDB,
		})

		new Frontend(this, 'Frontend')
		
		new Upload(this, 'Upload', {
			dynamoDB: db.dynamoDB,
		})

		new V2(this, 'V2')
	}
}
