
import cdk = require('@aws-cdk/core');
import ec2 = require('@aws-cdk/aws-ec2');
import rds = require('@aws-cdk/aws-rds');
import * as secretsmanager from '@aws-cdk/aws-secretsmanager';
import { RetentionDays } from '@aws-cdk/aws-logs';
import s3 = require('@aws-cdk/aws-s3');
import { CfnOutput } from '@aws-cdk/core';
import { RemovalPolicy } from '@aws-cdk/core';

export class Common extends cdk.Construct {
	public readonly vpc: ec2.Vpc;
	public readonly dbSecGrp: ec2.SecurityGroup;
	public readonly dbSecret: secretsmanager.ISecret;
	public readonly dbEndpoint: string;
	public readonly cluster: rds.DatabaseCluster;
	public readonly csvBucket: s3.Bucket;
	public readonly uploadBucket: s3.Bucket;

	constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
		super(scope, id)

		//EXTRACT
		this.uploadBucket = new s3.Bucket(this, 'Upload', {
			removalPolicy: RemovalPolicy.DESTROY,
			// blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL,
			cors: [
				{
					allowedOrigins: [
						"*", //if I only do wowmate.io it will not work locally
					],
					allowedMethods: [
						s3.HttpMethods.POST,
						s3.HttpMethods.GET, //should be removable
						s3.HttpMethods.PUT,//should be removable
						s3.HttpMethods.HEAD,
					],
					allowedHeaders: [
						"*", //dunno
					],
				}
			],
			metrics: [{
				id: 'metric',
			}]
		})

		this.csvBucket = new s3.Bucket(this, 'CSV', {
			removalPolicy: cdk.RemovalPolicy.DESTROY,
			blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL,
			metrics: [{ //enables advanced s3metrics
				id: 'metric',
			}]
		})

		this.vpc = new ec2.Vpc(this, 'Vpc', {
			natGateways: 1,
		})
		let vpc = this.vpc

		this.dbSecGrp = new ec2.SecurityGroup(this, 'DBAccess', {
			vpc,
		})
		this.dbSecGrp.addIngressRule(this.dbSecGrp, ec2.Port.tcp(5432), 'allow db connection')


		this.cluster = new rds.DatabaseCluster(this, 'DB', {
			engine: rds.DatabaseClusterEngine.auroraPostgres({
				version: rds.AuroraPostgresEngineVersion.VER_11_6,
			}),
			masterUser: {
				username: 'clusteradmin'
			},
			instanceProps: {
				vpc: this.vpc,
				securityGroups: [this.dbSecGrp],
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
			s3ImportBuckets: [this.csvBucket],
		})
		this.dbSecret = this.cluster.secret!

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
			securityGroup: this.dbSecGrp,
		});
	}
}