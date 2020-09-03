import cdk = require('@aws-cdk/core');
import route53= require('@aws-cdk/aws-route53');
import targets = require('@aws-cdk/aws-route53-targets');
import ec2 = require('@aws-cdk/aws-ec2');
import rds = require('@aws-cdk/aws-rds');
import ecs = require('@aws-cdk/aws-ecs');
import ecsPatterns = require('@aws-cdk/aws-ecs-patterns');
import * as elbv2 from '@aws-cdk/aws-elasticloadbalancingv2';
import * as secretsmanager from '@aws-cdk/aws-secretsmanager';
import { RetentionDays } from '@aws-cdk/aws-logs';
import s3 = require('@aws-cdk/aws-s3');
import iam = require('@aws-cdk/aws-iam');

export class Api extends cdk.Construct {
	public readonly vpc: ec2.Vpc;
	public readonly securityGrp: ec2.SecurityGroup;
	public readonly dbCreds: secretsmanager.ISecret;
	public readonly bucket: s3.Bucket;

	constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
		super(scope, id)

		const csvBucket = new s3.Bucket(this, 'CSV', {
			removalPolicy: cdk.RemovalPolicy.DESTROY,
			blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL,
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

		const auroraPostgres = new rds.DatabaseCluster(this, 'ImportDB', {
			engine: rds.DatabaseClusterEngine.auroraPostgres({
				version: rds.AuroraPostgresEngineVersion.VER_11_7,
			}),
			masterUser: {
				username: 'clusteradmin'
			},
			instanceProps: {
				vpc: vpc,
				securityGroups: [dbGroup],
				//TODO: try small even tho it doesn't seem to exist
				instanceType: ec2.InstanceType.of(ec2.InstanceClass.BURSTABLE3, ec2.InstanceSize.MEDIUM),
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
		this.dbCreds = auroraPostgres.secret!

		new ec2.BastionHostLinux(this, 'BastionHost', { 
			vpc,
			securityGroup: dbGroup,
		});
		
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