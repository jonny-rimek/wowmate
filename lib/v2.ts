import cdk = require('@aws-cdk/core');
import ec2 = require('@aws-cdk/aws-ec2');
import rds = require('@aws-cdk/aws-rds');
import ecs = require('@aws-cdk/aws-ecs');
import ecsPatterns = require('@aws-cdk/aws-ecs-patterns');

export class V2 extends cdk.Construct {
	constructor(scope: cdk.Construct, id: string) {
		super(scope, id)

		const vpc = new ec2.Vpc(this, 'Vpc', {
			subnetConfiguration: [{
				name: 'publicSubnet',
				subnetType: ec2.SubnetType.PUBLIC,
			}],
			natGateways: 0,
		})

		const postgres = new rds.DatabaseInstance(this, 'Postgres', {
			engine: rds.DatabaseInstanceEngine.POSTGRES,
			instanceClass: ec2.InstanceType.of(ec2.InstanceClass.BURSTABLE2, ec2.InstanceSize.MICRO),
			masterUsername: 'postgres',
			deletionProtection: false,
			vpc,
			vpcPlacement: { subnetType: ec2.SubnetType.PUBLIC }
		})

		postgres.connections.allowFromAnyIpv4(ec2.Port.tcp(5432))

		const loadBalancedFargateService = new ecsPatterns.ApplicationLoadBalancedFargateService(this, 'Service', {
		vpc,
		memoryLimitMiB: 512,
		cpu: 256,
		desiredCount: 1,
		taskImageOptions: {
			image: ecs.ContainerImage.fromAsset('services/api'),
		},
		});
	}
}
