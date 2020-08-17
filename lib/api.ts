import cdk = require('@aws-cdk/core');
import targets = require('@aws-cdk/aws-route53-targets');
import ec2 = require('@aws-cdk/aws-ec2');
import rds = require('@aws-cdk/aws-rds');
import ecs = require('@aws-cdk/aws-ecs');
import ecsPatterns = require('@aws-cdk/aws-ecs-patterns');

interface VpcProps extends cdk.StackProps {
	vpc: ec2.IVpc;
}

export class Api extends cdk.Construct {
	constructor(scope: cdk.Construct, id: string, props: VpcProps) {
		super(scope, id)

		//IMPROVE: add https redirect
		//need to define the cluster seperately and in it the VPC i think
		const loadBalancedFargateService = new ecsPatterns.ApplicationLoadBalancedFargateService(this, 'Service', {
			vpc: props.vpc,
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
