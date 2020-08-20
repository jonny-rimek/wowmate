import cdk = require('@aws-cdk/core');
import route53= require('@aws-cdk/aws-route53');
import targets = require('@aws-cdk/aws-route53-targets');
import ec2 = require('@aws-cdk/aws-ec2');
import rds = require('@aws-cdk/aws-rds');
import ecs = require('@aws-cdk/aws-ecs');
import ecsPatterns = require('@aws-cdk/aws-ecs-patterns');
import * as elbv2 from '@aws-cdk/aws-elasticloadbalancingv2';

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

			//NOTE: remove in production
			vpcPlacement: { subnetType: ec2.SubnetType.PUBLIC },
			removalPolicy: cdk.RemovalPolicy.DESTROY,
			deletionProtection: false,
			//NOTE: remove in production

			enablePerformanceInsights: true,
			monitoringInterval: cdk.Duration.seconds(60),
			cloudwatchLogsExports: ['postgresql'],
			//improve set max duration of log
			// cloudwatchLogsRetention
		})
		postgres.connections.allowFromAnyIpv4(ec2.Port.tcp(5432))

		const hostedZone = route53.HostedZone.fromHostedZoneAttributes(this, 'HostedZone', {
			zoneName: 'wowmate.io',
			hostedZoneId: 'Z3LVG9ZF2H87DX',
		});

		//IMPROVE: add https redirect
		new ecsPatterns.ApplicationLoadBalancedFargateService(this, 'Service', {
			vpc: vpc,
			domainName: 'api.wowmate.io',
			domainZone: hostedZone,
			redirectHTTP: true,
			protocol: elbv2.ApplicationProtocol.HTTPS,
			memoryLimitMiB: 512,
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
