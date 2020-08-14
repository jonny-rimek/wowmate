import cdk = require('@aws-cdk/core');
import rds = require('@aws-cdk/aws-rds');
import ec2 = require('@aws-cdk/aws-ec2');

interface VpcProps extends cdk.StackProps {
	vpc: ec2.IVpc;
}

export class Postgres extends cdk.Construct {
	public readonly postgres: rds.DatabaseInstance;

	constructor(scope: cdk.Construct, id: string, props: VpcProps) {
		super(scope, id);

		const postgres = new rds.DatabaseInstance(this, 'Postgres', {
			vpc: props.vpc,
			engine: rds.DatabaseInstanceEngine.postgres({
				version: rds.PostgresEngineVersion.VER_11_7,
			}),
			instanceType: ec2.InstanceType.of(ec2.InstanceClass.BURSTABLE2, ec2.InstanceSize.MICRO),
			masterUsername: 'postgres',
			// vpcPlacement: { subnetType: ec2.SubnetType.PUBLIC },
			//NOTE: remove in production
			removalPolicy: cdk.RemovalPolicy.DESTROY,
			deletionProtection: false,
		})
		postgres.connections.allowFromAnyIpv4(ec2.Port.tcp(5432))

		this.postgres = postgres

	}
}
