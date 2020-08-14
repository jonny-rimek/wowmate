import cdk = require('@aws-cdk/core');
import targets = require('@aws-cdk/aws-route53-targets');
import ec2 = require('@aws-cdk/aws-ec2');
import rds = require('@aws-cdk/aws-rds');
import ecs = require('@aws-cdk/aws-ecs');
import elbv2 = require('@aws-cdk/aws-elasticloadbalancingv2');
import route53= require('@aws-cdk/aws-route53');
import ecsPatterns = require('@aws-cdk/aws-ecs-patterns');
import * as lambda from '@aws-cdk/aws-lambda';
import apigateway = require('@aws-cdk/aws-apigateway');
import acm = require('@aws-cdk/aws-certificatemanager');
import { BaseLoadBalancer } from '@aws-cdk/aws-elasticloadbalancingv2';
import s3 = require('@aws-cdk/aws-s3');
import { RemovalPolicy, Duration } from '@aws-cdk/core';

interface VpcProps extends cdk.StackProps {
	vpc: ec2.IVpc;
}

export class Api extends cdk.Construct {
	constructor(scope: cdk.Construct, id: string, props: VpcProps) {
		super(scope, id)

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
			},
		});
	}
}
