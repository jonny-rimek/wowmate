import cdk = require('@aws-cdk/core');
import ec2 = require('@aws-cdk/aws-ec2');
import rds = require('@aws-cdk/aws-rds');

export class V2 extends cdk.Construct {
	constructor(scope: cdk.Construct, id: string) {
		super(scope, id)

		const vpc = new ec2.Vpc(this, 'Vpc', {
			cidr: '10.0.0.0/16',
			maxAzs: 2,
			subnetConfiguration: [{
				cidrMask: 26,
				name: 'publicSubnet',
				subnetType: ec2.SubnetType.PUBLIC,
			}],
			natGateways: 0
		})

		const postgres = new rds.DatabaseInstance(this, 'Postgres', {
			engine: rds.DatabaseInstanceEngine.POSTGRES,
			instanceClass: ec2.InstanceType.of(ec2.InstanceClass.BURSTABLE2, ec2.InstanceSize.MICRO),
			masterUsername: 'postgres',
			vpc,
			vpcPlacement: { subnetType: ec2.SubnetType.PUBLIC }
		})

		postgres.connections.allowDefaultPortFromAnyIpv4
	}
}
