import cdk = require('@aws-cdk/core');
import ecs = require('@aws-cdk/aws-ecs');
import ecsPatterns = require('@aws-cdk/aws-ecs-patterns');
import ec2 = require('@aws-cdk/aws-ec2');

interface VpcProps extends cdk.StackProps {
	vpc: ec2.IVpc;
}

export class Converter extends cdk.Construct {
	constructor(scope: cdk.Construct, id: string, props: VpcProps) {
		super(scope, id)

		const queueProcessingFargateService = new ecsPatterns.QueueProcessingFargateService(this, 'Service', {
			vpc: props.vpc,
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
