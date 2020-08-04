import cdk = require('@aws-cdk/core');
import ec2 = require('@aws-cdk/aws-ec2');
import rds = require('@aws-cdk/aws-rds');
import ecs = require('@aws-cdk/aws-ecs');
import elbv2 = require('@aws-cdk/aws-elasticloadbalancingv2');
import route53= require('@aws-cdk/aws-route53');
import ecsPatterns = require('@aws-cdk/aws-ecs-patterns');

export class V2 extends cdk.Construct {
	constructor(scope: cdk.Construct, id: string) {
		super(scope, id)
		const vpc = new ec2.Vpc(this, 'WowmateVpc', {
			// subnetConfiguration: [{
			// 	name: 'publicSubnet',
			// 	subnetType: ec2.SubnetType.PUBLIC,
			// }],
			natGateways: 1,
		})

		/*
		const postgres = new rds.DatabaseInstance(this, 'Postgres', {
			engine: rds.DatabaseInstanceEngine.postgres({
				version: rds.PostgresEngineVersion.VER_11_7,
			}),
			instanceType: ec2.InstanceType.of(ec2.InstanceClass.BURSTABLE2, ec2.InstanceSize.MICRO),
			masterUsername: 'postgres',
			vpc,
			// vpcPlacement: { subnetType: ec2.SubnetType.PUBLIC },
			//NOTE: remove in production
			removalPolicy: cdk.RemovalPolicy.DESTROY,
			deletionProtection: false,
		})
		postgres.connections.allowFromAnyIpv4(ec2.Port.tcp(5432))
		*/
		/*
		const fargateTask = new ecs.FargateTaskDefinition(this, 'FargateTask', {
			cpu: 256,
			memoryLimitMiB: 512,
		})

		const container = fargateTask.addContainer("GinContainer", {
			image: ecs.ContainerImage.fromAsset('services/api'),
			cpu: 256,
			memoryLimitMiB: 512
		})
		container.addPortMappings({containerPort: 80})

		const cluster = new ecs.Cluster(this, 'Cluster', {
			containerInsights: true,
			vpc
		})

		const fargateService = new ecs.FargateService(this, 'FargateService', {
			cluster,
			taskDefinition: fargateTask,
			desiredCount: 1,
			platformVersion: ecs.FargatePlatformVersion.VERSION1_4,
			assignPublicIp: true,
		})

		const lb = new elbv2.ApplicationLoadBalancer(this, 'LB', { vpc, internetFacing: true });
		const listener = lb.addListener('Listener', { port: 80 });
		const targetGroup1 = listener.addTargets('ECS1', {
			port: 80,
			targets: [fargateService]
		});
		*/
		// /*

		const hostedZone = route53.HostedZone.fromHostedZoneAttributes(this, 'HostedZone', {
			zoneName: 'wowmate.io',
			hostedZoneId: 'Z3LVG9ZF2H87DX',
		});

		//TODO: add api.wowmate.io domain and https
		//need to define the cluster seperately and in it the VPC i think
		const loadBalancedFargateService = new ecsPatterns.ApplicationLoadBalancedFargateService(this, 'Service', {
			vpc,
			domainName: 'api.wowmate.io',
			domainZone: hostedZone,
			memoryLimitMiB: 512,
			protocol: elbv2.ApplicationProtocol.HTTPS,
			cpu: 256,
			desiredCount: 1,
			publicLoadBalancer: true,
			platformVersion: ecs.FargatePlatformVersion.VERSION1_4,
			taskImageOptions: {
				image: ecs.ContainerImage.fromAsset('services/api'),
			},
		});
		// */
	}
}
