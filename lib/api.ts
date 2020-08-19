import cdk = require('@aws-cdk/core');
import targets = require('@aws-cdk/aws-route53-targets');
import ec2 = require('@aws-cdk/aws-ec2');
import rds = require('@aws-cdk/aws-rds');
import ecs = require('@aws-cdk/aws-ecs');
import ecsPatterns = require('@aws-cdk/aws-ecs-patterns');

// interface VpcProps extends cdk.StackProps {
// 	vpc: ec2.IVpc;
// }

export class Api extends cdk.Construct {
	public readonly vpc: ec2.Vpc;
	constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
		super(scope, id)

		const vpc = new ec2.Vpc(this, 'WowmateVpc', {
			natGateways: 1,
		});

		const postgres = new rds.DatabaseInstance(this, 'Postgres', {
			vpc: vpc,
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

		//IMPROVE: add https redirect
		//need to define the cluster seperately and in it the VPC i think
		const loadBalancedFargateService = new ecsPatterns.ApplicationLoadBalancedFargateService(this, 'Service', {
			vpc: vpc,
			//TODO: reactivate
			// domainName: 'api.wowmate.io',
			// domainZone: hostedZone,
			memoryLimitMiB: 512,
			// protocol: elbv2.ApplicationProtocol.HTTPS,
			cpu: 256,
			desiredCount: 1,
			publicLoadBalancer: true,
			platformVersion: ecs.FargatePlatformVersion.VERSION1_4,
			taskImageOptions: {
				image: ecs.ContainerImage.fromAsset('services/api'),
				environment: {
					GIN_MODE: "release"
				}
			},
		});
	}
}
