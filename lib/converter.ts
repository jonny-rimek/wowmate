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

export class Converter extends cdk.Construct {
	constructor(scope: cdk.Construct, id: string) {
		super(scope, id)

		const queueProcessingFargateService = new ecsPatterns.QueueProcessingFargateService(this, 'Service', {
			memoryLimitMiB: 512,
			cpu: 256,
			image: ecs.ContainerImage.fromAsset('services/converter'),
			// (optional, default: CMD value built into container image.)
			// command: ["-c", "4", "amazon.com"],
			desiredTaskCount: 1,
			environment: {
				TEST_ENVIRONMENT_VARIABLE1: "test environment variable 1 value",
				TEST_ENVIRONMENT_VARIABLE2: "test environment variable 2 value"
			},
		});
	}
}
