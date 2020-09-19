import cdk = require('@aws-cdk/core');
import ec2 = require('@aws-cdk/aws-ec2');
import rds = require('@aws-cdk/aws-rds');
import * as secretsmanager from '@aws-cdk/aws-secretsmanager';
import { RetentionDays } from '@aws-cdk/aws-logs';
import { HttpProxyIntegration, HttpApi, LambdaProxyIntegration, HttpMethod } from '@aws-cdk/aws-apigatewayv2';
import s3 = require('@aws-cdk/aws-s3');
import iam = require('@aws-cdk/aws-iam');
import lambda = require('@aws-cdk/aws-lambda');
import { CfnOutput } from '@aws-cdk/core';

export class Api extends cdk.Construct {
	public readonly vpc: ec2.Vpc;
	public readonly securityGrp: ec2.SecurityGroup;
	public readonly dbCreds: secretsmanager.ISecret;
	public readonly bucket: s3.Bucket;
	public readonly lambda: lambda.Function;
	public readonly api: HttpApi;
	public readonly dbEndpoint: string;
	public readonly cluster: rds.DatabaseCluster;

	constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
		super(scope, id)

		const csvBucket = new s3.Bucket(this, 'CSV', {
			removalPolicy: cdk.RemovalPolicy.DESTROY,
			blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL,
			metrics: [{ //enables advanced s3metrics
				id: 'metric',
			}]
		})
		this.bucket = csvBucket

		const role = new iam.Role(this, "Role", {
			assumedBy: new iam.ServicePrincipal("rds.amazonaws.com"), // required
		});

		role.addToPolicy(
			new iam.PolicyStatement({
				effect: iam.Effect.ALLOW,
				resources: [csvBucket.bucketArn, `${csvBucket.bucketArn}/*`],
				actions: ["s3:GetObject", "s3:ListBucket"],
			})
		);

		const vpc = new ec2.Vpc(this, 'WowmateVpc', {
			natGateways: 1,
		})

		let dbGroup = new ec2.SecurityGroup(this, 'DBAccess', {
			vpc
		})
		dbGroup.addIngressRule(dbGroup, ec2.Port.tcp(5432), 'allow db connection')

		this.vpc = vpc
		this.securityGrp = dbGroup

		this.cluster = new rds.DatabaseCluster(this, 'ImportDB', {
			engine: rds.DatabaseClusterEngine.auroraPostgres({
				// version: rds.AuroraPostgresEngineVersion.VER_10_11,
				version: rds.AuroraPostgresEngineVersion.VER_11_6,
			}),
			masterUser: {
				username: 'clusteradmin'
			},
			instanceProps: {
				vpc: vpc,
				securityGroups: [dbGroup],
				instanceType: ec2.InstanceType.of(ec2.InstanceClass.BURSTABLE3, ec2.InstanceSize.MEDIUM),
				enablePerformanceInsights: true,
				performanceInsightRetention: rds.PerformanceInsightRetention.DEFAULT, //7days
			},
			instances: 1,
			defaultDatabaseName: 'wm',

			cloudwatchLogsExports: ['postgresql'],
			cloudwatchLogsRetention: RetentionDays.TWO_WEEKS,
			monitoringInterval: cdk.Duration.seconds(60),

			//NOTE: remove in production
			removalPolicy: cdk.RemovalPolicy.DESTROY,
			deletionProtection: false,
			//NOTE: remove in production
			// s3ImportBuckets: [csvBucket],
		})
		this.dbCreds = this.cluster.secret!

		//NOTE: 11.6 works with the proxy, just activate and remove the old this.dbEndpoint
		//		every lambda should still work
		// const proxy = this.cluster.addProxy('DBProxy', {
		// 	secrets: [this.cluster.secret!],
		// 	vpc: vpc,
		// 	securityGroups: [dbGroup],
		// })
		// this.dbEndpoint = proxy.endpoint
		this.dbEndpoint = this.cluster.clusterEndpoint.hostname

		new CfnOutput(this, 'DBEndpoint', {
			value: this.dbEndpoint
		}),

		new ec2.BastionHostLinux(this, 'BastionHost', { 
			vpc,
			securityGroup: dbGroup,
		});

		const migrateLambda = new lambda.Function(this, 'MigrateLambda', {
			code: lambda.Code.fromAsset('services/migrate'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: cdk.Duration.seconds(30),
			environment: {
				SECRET_ARN: this.cluster.secret!.secretArn,
			},
			logRetention: RetentionDays.ONE_WEEK,
			tracing: lambda.Tracing.ACTIVE,
			vpc: vpc,
			securityGroups: [dbGroup],
			//TODO: add DLQ
		})
		this.cluster.secret?.grantRead(migrateLambda)


		const topDamageLambda = new lambda.Function(this, 'TopDamageLambda', {
			code: lambda.Code.fromAsset('services/api'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: cdk.Duration.seconds(30),
			environment: {
				DB_ENDPOINT: this.dbEndpoint,
				SECRET_ARN: this.cluster.secret!.secretArn,
			},
			logRetention: RetentionDays.ONE_WEEK,
			tracing: lambda.Tracing.ACTIVE,
			vpc: vpc,
			securityGroups: [dbGroup],
		})
		this.lambda = topDamageLambda
		this.cluster.secret?.grantRead(topDamageLambda)

		const topDamageIntegration = new LambdaProxyIntegration({
			handler: topDamageLambda
		})

		const httpApi = new HttpApi(this, 'Api')

		new CfnOutput(this, 'HttpApiEndpoint', {
			value: httpApi.url!,
		}),

		httpApi.addRoutes({
			path: '/api/combatlog/summary/{combatlog_uuid}/damage',
			methods: [HttpMethod.GET],
			integration: topDamageIntegration,
		})
		this.api = httpApi
	}
}