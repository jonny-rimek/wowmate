import cdk = require('@aws-cdk/core');
import ec2 = require('@aws-cdk/aws-ec2');
import * as secretsmanager from '@aws-cdk/aws-secretsmanager';
import { RetentionDays } from '@aws-cdk/aws-logs';
import lambda = require('@aws-cdk/aws-lambda');

interface Props extends cdk.StackProps {
	vpc: ec2.IVpc;
	dbSecGrp: ec2.SecurityGroup
	dbSecret: secretsmanager.ISecret
	// dbEndpoint: string
}

export class Migrate extends cdk.Construct {
	public readonly migrateLambda: lambda.Function;

	constructor(scope: cdk.Construct, id: string, props: Props) {
		super(scope, id)

		const migrateLambda = new lambda.Function(this, 'MigrateLambda', {
			code: lambda.Code.fromAsset('services/migrate'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: cdk.Duration.seconds(30),
			environment: {
				SECRET_ARN: props.dbSecret.secretArn,
				//not using dbEndpoint because we get the endpoint from the secret
				//and it only runs once a day so connections aren't a problem
			},
			logRetention: RetentionDays.ONE_WEEK,
			tracing: lambda.Tracing.ACTIVE,
			vpc: props.vpc,
			securityGroups: [props.dbSecGrp],
			//TODO: add DLQ
		})
		props.dbSecret.grantRead(migrateLambda)
	}
}