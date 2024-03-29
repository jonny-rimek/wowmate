import { Construct, Stack, StackProps } from "@aws-cdk/core";
import { Frontend } from "./frontend/frontend";
import { Api } from "./api/api";
import { Convert } from "./upload/convert";
import { QueryTimestream } from "./upload/queryTimestream";
import { InsertResult } from "./upload/insert-result";
import { EtlDashboard } from "./upload/etl-dashboard";
import { ApiFrontendDashboard } from "./common/api-frontend-dashboard";
import { DynamoDB } from "./common/dynamodb";
import { Buckets } from "./common/buckets";
import { Kms } from "./common/kms";
import { Cloudtrail } from "./common/cloudtrail";
import { Timestream } from "./common/timestream";
import { Synthetics } from "./synthetics/synthetics";

interface Props extends StackProps {
	hostedZoneId: string
	hostedZoneName: string
	domainName: string
	apiDomainName: string
	stage: string
}

export class Wowmate extends Stack {
	constructor(scope: Construct, id: string, props: Props) {
		super(scope, id, props);
		const errorMail = 'hi@wowmate.io'

		const buckets = new Buckets(this, "Buckets-");

		new Cloudtrail(this, "Cloudtrail-", {
			uploadBucket: buckets.uploadBucket,
		});

		const timestream = new Timestream(this, "Timestream-");

		const dynamoDB = new DynamoDB(this, 'DynamoDB-')

		const kms = new Kms(this, 'Kms-')

		const api = new Api(this, 'Api-', {
			dynamoDB: dynamoDB.table,
			uploadBucket: buckets.uploadBucket,
			hostedZoneId: props.hostedZoneId,
			hostedZoneName: props.hostedZoneName,
			apiDomainName: props.apiDomainName,
			accessLogBucket: buckets.accessLogBucket,
		})

		const frontend = new Frontend(this, "Frontend-", {
			hostedZoneId: props.hostedZoneId,
			hostedZoneName: props.hostedZoneName,
            domainName: props.domainName,
			accessLogBucket: buckets.accessLogBucket,
		});

		const queryKeys = new QueryTimestream(this, "QueryKeys-", {
			dynamoDB: dynamoDB.table,
			timestreamArn: timestream.timestreamArn,
			lambdaDescription: 'queries timestream for the keys, keys per dungeon etc and posts it to SNS',
			codePath: 'dist/upload/query-timestream/keys',
			envVars: {
				LOG_LEVEL: "info", //only info or debug are support
				LOCAL: "false",
			},
			key: kms.key,
		});

		const insertKeysToDynamoDB = new InsertResult(this, "InsertKeysToDynamodb-", {
			dynamoDB: dynamoDB.table,
            topic: queryKeys.topic,
			topicDLQ: queryKeys.topicDLQ,
			lambdaDescription: 'writes the keys, keys per dungeon etc to dynamodb for access by the frontend',
			codePath: 'dist/upload/insert/dynamodb/keys',
            envVars: {
				LOG_LEVEL: "info", //only info or debug are support
				LOCAL: "false",
			}
		});

		const insertKeysToTimestream = new InsertResult(this, "InsertKeysToTimestream-", {
			dynamoDB: dynamoDB.table,
			topic: queryKeys.topic,
			topicDLQ: queryKeys.topicDLQ,
			lambdaDescription: 'writes the simple damage summary to timestream for later analyzing e.g. statistics per dungeon',
			codePath: 'dist/upload/insert/timestream/keys',
			envVars: {
				LOG_LEVEL: "info", //only info or debug are support
				LOCAL: "false",
			}
		});

		//query player damage done
		const queryPlayerDamageDone = new QueryTimestream(this, "QueryPlayerDamageDone-", {
			dynamoDB: dynamoDB.table,
			timestreamArn: timestream.timestreamArn,
			lambdaDescription: 'queries timestream and creates the advanced damage summary for the combatlog specific page',
			//upload/query/player-damage-done
			//upload/query/keys
			codePath: 'dist/upload/query-timestream/player-damage-done',
			envVars: {
				LOG_LEVEL: "info", //only info or debug are support
				LOCAL: "false",
			},
			key: kms.key,
		});

		const insertPlayerDamageDoneToDynamodb = new InsertResult(this, "InsertPlayerDamageDoneToDynamodb-", {
			dynamoDB: dynamoDB.table,
			topic: queryPlayerDamageDone.topic,
			topicDLQ: queryPlayerDamageDone.topicDLQ,
			lambdaDescription: 'writes the player damage done summary for combatlog specific page to dynamodb for access by the frontend',
			codePath: 'dist/upload/insert/dynamodb/player-damage-done',
			envVars: {
				LOG_LEVEL: "info", //only info or debug are support
				LOCAL: "false",
			}
		});

		const convert = new Convert(this, "Convert-", {
			timestreamArn: timestream.timestreamArn,
			uploadBucket: buckets.uploadBucket,
			queryTimestreamLambdas: [
				//every lambda that subscribes gets a notification when a combatlog is processed
				queryKeys.lambda,
				queryPlayerDamageDone.lambda,
			],
			envVars: {
				LOG_LEVEL: "info", //only info or debug are support
				LOCAL: "false",
			},
			dynamodb: dynamoDB.table,
			key: kms.key,
		});

		new ApiFrontendDashboard(this, 'UserFacing-', {
			getKeysLambda: api.getKeysLambda,
			getKeysPerDungeonLambda: api.getKeysPerDungeonLambda,
			api: api.api,
			s3: frontend.bucket,
			cloudfront: frontend.cloudfront,
			errorMail: errorMail,
		})

		new EtlDashboard(this, 'Etl-', {
			convertLambda: convert.lambda,
			convertQueue: convert.queue,
			convertDLQ: convert.DLQ,
			queryKeys: queryKeys.lambda,
			insertKeysToDynamoDB: insertKeysToDynamoDB.lambda,
			insertKeysToTimestream: insertKeysToTimestream.lambda,
			queryPlayerDamageDone: queryPlayerDamageDone.lambda,
			insertPlayerDamageDoneToDynamodb: insertPlayerDamageDoneToDynamodb.lambda,
			queryKeysTopicDLQ: queryKeys.topicDLQ,
			queryPlayerDamageDoneTopicDLQ: queryPlayerDamageDone.topicDLQ,
			queryPlayerDamageDoneLambdaDLQ: queryPlayerDamageDone.lambdaDLQ,
			queryKeysLambdaDLQ: queryKeys.lambdaDLQ,
			insertKeysToDynamoDBLambdaDLQ: insertKeysToDynamoDB.lambdaDLQ,
			insertKeysToTimestreamLambdaDLQ: insertKeysToTimestream.lambdaDLQ,
			insertPlayerDamageDoneDynamoDBLambdaDLQ: insertPlayerDamageDoneToDynamodb.lambdaDLQ,
			presignLambda: api.presignLambda,
			uploadBucket: buckets.uploadBucket,
			presignApi: api.api,
			dynamoDB: dynamoDB.table,
			convertTopic: convert.topic,
            queryKeysTopic: queryKeys.topic,
			queryPlayerDamageDoneTopic: queryPlayerDamageDone.topic,
			errorMail: errorMail,
		})

		new Synthetics(this, 'Canary-',  {
			errorMail: errorMail,
			stage: props.stage,
			apiUrl: api.api.url!,
			key: kms.key,
			domainName: props.domainName,
		})
	}
}