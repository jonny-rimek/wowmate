import cdk = require('@aws-cdk/core');
import ec2 = require('@aws-cdk/aws-ec2');

export class Vpc extends cdk.Construct {
	public readonly vpc: ec2.Vpc;

	constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
		super(scope, id)

		this.vpc = new ec2.Vpc(this, 'Vpc', {
			natGateways: 1,
		})
	}
}