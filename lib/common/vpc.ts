import cdk = require('@aws-cdk/core');
import ec2 = require('@aws-cdk/aws-ec2');

export class Vpc extends cdk.Construct {
	public readonly vpc: ec2.Vpc;
	public readonly natGateways: ec2.GatewayConfig[];

	constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
		super(scope, id)

		const nat = ec2.NatProvider.gateway()
		this.natGateways = nat.configuredGateways

		this.vpc = new ec2.Vpc(this, 'Vpc', {
			natGateways: 1,
			natGatewayProvider: nat,
			//gateway endpoint will lead to no datatransfer cost for s3 to vpc/NATgateway
			gatewayEndpoints: {
				S3: {
					service: ec2.GatewayVpcEndpointAwsService.S3,
				}
			}
		})
	}
}