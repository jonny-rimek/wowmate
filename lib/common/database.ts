import cdk = require('@aws-cdk/core');
import ec2 = require('@aws-cdk/aws-ec2');
import rds = require('@aws-cdk/aws-rds');
import * as secretsmanager from '@aws-cdk/aws-secretsmanager';
import { RetentionDays } from '@aws-cdk/aws-logs';
import s3 = require('@aws-cdk/aws-s3');
import { CfnOutput } from '@aws-cdk/core';

interface Props extends cdk.StackProps {
	csvBucket: s3.Bucket
	vpc: ec2.Vpc
}

export class Database extends cdk.Construct {
	public readonly dbSecGrp: ec2.SecurityGroup;
	public readonly dbSecret: secretsmanager.ISecret;
	public readonly dbEndpoint: string;
	public readonly cluster: rds.DatabaseCluster;

	constructor(scope: cdk.Construct, id: string, props: Props) {
		super(scope, id)

		this.dbSecGrp = new ec2.SecurityGroup(this, 'Access', {
			vpc: props.vpc,
		})
		this.dbSecGrp.addIngressRule(this.dbSecGrp, ec2.Port.tcp(5432), 'allow db connection')

		this.cluster = new rds.DatabaseCluster(this, 'DB', {
			engine: rds.DatabaseClusterEngine.auroraPostgres({
				version: rds.AuroraPostgresEngineVersion.VER_11_6,
			}),
			masterUser: {
				username: 'clusteradmin' //admin is reserved and can't be used
			},
			instanceProps: {
				vpc: props.vpc,
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
			s3ImportBuckets: [props.csvBucket],
		})
		this.dbSecret = this.cluster.secret!

		//NOTE: 11.6 works with the proxy, just activate and remove the old this.dbEndpoint
		//		every lambda will still work
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
			vpc: props.vpc,
			securityGroup: this.dbSecGrp,
		});
	}
}