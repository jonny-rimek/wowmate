import cdk = require('@aws-cdk/core');
import ddb = require('@aws-cdk/aws-dynamodb');
import s3 = require('@aws-cdk/aws-s3');
import lambda = require('@aws-cdk/aws-lambda');
import sfn = require('@aws-cdk/aws-stepfunctions');
import tasks = require('@aws-cdk/aws-stepfunctions-tasks');
import { RemovalPolicy, Duration } from '@aws-cdk/core';
import { BlockPublicAccess } from '@aws-cdk/aws-s3';
import targets = require('@aws-cdk/aws-events-targets');
import iam = require('@aws-cdk/aws-iam');
import { Effect } from '@aws-cdk/aws-iam';
import cloudtrail = require('@aws-cdk/aws-cloudtrail');
import apigateway = require('@aws-cdk/aws-apigateway');
import s3deploy = require('@aws-cdk/aws-s3-deployment');
import cloudfront = require('@aws-cdk/aws-cloudfront');
// import events = require('@aws-cdk/aws-events');
// import { Result } from '@aws-cdk/aws-stepfunctions';

export class WowmateStack extends cdk.Stack {
	constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
		super(scope, id, props);

		//FRONTEND
		const frontendBucket = new s3.Bucket(this, 'FrontendBucket', {
		websiteIndexDocument: 'index.html',
		publicReadAccess: true
		});

		const distribution = new cloudfront.CloudFrontWebDistribution(this, 'Distribution', {
		originConfigs: [
			{
			s3OriginSource: {
				s3BucketSource: frontendBucket
			},
			behaviors : [ {isDefaultBehavior: true}]
			}
		]
		});

		new s3deploy.BucketDeployment(this, 'DeployWebsite', {
		sources: [s3deploy.Source.asset('./frontend/dist')],
		destinationBucket: frontendBucket,
		distribution,
		});

		//DYNAMODB
		const db = new ddb.Table(this, 'DDB', {
			partitionKey: { name: 'pk', type: ddb.AttributeType.STRING },
			sortKey: {name: 'sk', type: ddb.AttributeType.NUMBER},
			removalPolicy: RemovalPolicy.DESTROY,
            billingMode: ddb.BillingMode.PAY_PER_REQUEST
		})

		db.addGlobalSecondaryIndex({
			indexName: 'GSI1',
			partitionKey: {name: 'gsi1pk', type: ddb.AttributeType.NUMBER},
			sortKey: {name: 'sk', type: ddb.AttributeType.NUMBER}
		})

		db.addGlobalSecondaryIndex({
			indexName: 'GSI2',
			partitionKey: {name: 'gsi2pk', type: ddb.AttributeType.STRING},
			sortKey: {name: 'sk', type: ddb.AttributeType.NUMBER}
		})

		//API
		const bossFightDamageFunc = new lambda.Function(this, 'BossFightDamage', {
			code: lambda.Code.asset("api-service/boss-fight-damage"),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: Duration.seconds(3),
			environment: {DDB_NAME: db.tableName}
		})

		db.grantReadData(bossFightDamageFunc)
		
		const api = new apigateway.LambdaRestApi(this, 'api', {
			handler: bossFightDamageFunc,
			proxy: false,
		});

		const basePath = api.root.addResource('api');
		const bossFightPath = basePath.addResource('boss-fight');
		const damagePath = bossFightPath.addResource('damage');
		const encounterId = damagePath.addResource('{boss-fight-uuid}');
		encounterId.addMethod('GET')

		//CLOUDTRAIL
		const trail = new cloudtrail.Trail(this, 'CloudTrail', {
			sendToCloudWatchLogs: true,
			managementEvents: cloudtrail.ReadWriteType.WRITE_ONLY,
		});

		//S3 BUCKETS
		const uploadBucket = new s3.Bucket(this, 'Upload', {
			removalPolicy: RemovalPolicy.DESTROY,
			blockPublicAccess: BlockPublicAccess.BLOCK_ALL,
		})

		trail.addS3EventSelector([uploadBucket.bucketArn + "/"], {
			readWriteType: cloudtrail.ReadWriteType.WRITE_ONLY,
		})

		const parquetBucket = new s3.Bucket(this, 'Parquet', {
			removalPolicy: RemovalPolicy.DESTROY,
			blockPublicAccess: BlockPublicAccess.BLOCK_ALL,
		})

		const athenaBucket = new s3.Bucket(this, 'Athena', {
			removalPolicy: RemovalPolicy.DESTROY,
			blockPublicAccess: BlockPublicAccess.BLOCK_ALL,
		})

		//LAMBDA
		const sizeFunc = new lambda.Function(this, 'Size', {
			code: lambda.Code.asset("upload-service/size"),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: Duration.seconds(3),
		})

		const parquetFunc = new lambda.Function(this, 'ParquetFunc', {
			code: lambda.Code.asset("upload-service/parquet"),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: Duration.seconds(10),
			environment: {TARGET_BUCKET_NAME: parquetBucket.bucketName}
		})

		const athenaFunc = new lambda.Function(this, 'AthenaFunc', {
			code: lambda.Code.asset("upload-service/athena"),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: Duration.seconds(3),
		})

		const checkFunc = new lambda.Function(this, 'Check', {
			code: lambda.Code.asset("upload-service/check"),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: Duration.seconds(3),
		})
	
		const impFunc = new lambda.Function(this, 'Import', {
			code: lambda.Code.asset("upload-service/import"),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: Duration.seconds(10),
			environment: {DDB_NAME: db.tableName}
		})

		const imp2Func = new lambda.Function(this, 'Import2', {
			code: lambda.Code.asset("upload-service/import2"),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: Duration.seconds(10),
			environment: {DDB_NAME: db.tableName}
		})

		//IAM
		uploadBucket.grantRead(sizeFunc)

		uploadBucket.grantRead(parquetFunc)
		parquetBucket.grantPut(parquetFunc)

		parquetBucket.grantRead(athenaFunc)

		athenaBucket.grantReadWrite(athenaFunc)
		athenaBucket.grantWrite(checkFunc)

		//permission are from the aws docs https://docs.aws.amazon.com/athena/latest/ug/example-policies-workgroup.html
		//they could probably be tightend a bit, especially differences between checkFunc and athenaFunc 
		const athenaGeneralPolicy = new iam.PolicyStatement({
			effect: Effect.ALLOW,
			actions: [ 
				"athena:ListWorkGroups",
                "athena:GetExecutionEngine",
                "athena:GetExecutionEngines",
                "athena:GetNamespace",
                "athena:GetCatalogs",
                "athena:GetNamespaces",
                "athena:GetTables",
				"athena:GetTable",
				"glue:GetTable"
			],
			resources: [
				'*'
			],
		})

		const athenaWorkgroupPolicy = new iam.PolicyStatement({
			effect: Effect.ALLOW,
			actions: [ 
				"athena:StartQueryExecution",
                "athena:GetQueryResults",
                "athena:DeleteNamedQuery",
                "athena:GetNamedQuery",
                "athena:ListQueryExecutions",
                "athena:StopQueryExecution",
                "athena:GetQueryResultsStream",
                "athena:ListNamedQueries",
                "athena:CreateNamedQuery",
                "athena:GetQueryExecution",
                "athena:BatchGetNamedQuery",
                "athena:BatchGetQueryExecution", 
                "athena:GetWorkGroup" 
			],
			resources: [
				'arn:aws:athena:eu-central-1:940880032268:workgroup/primary'
			],
		})

		athenaFunc.addToRolePolicy(athenaGeneralPolicy)
		checkFunc.addToRolePolicy(athenaGeneralPolicy)

		athenaFunc.addToRolePolicy(athenaWorkgroupPolicy)
		checkFunc.addToRolePolicy(athenaWorkgroupPolicy)

		athenaBucket.grantRead(impFunc)
		db.grantWriteData(impFunc)

		athenaBucket.grantRead(imp2Func)
		db.grantWriteData(imp2Func)


		//STEP FUNCTION
		const sizeJob = new sfn.Task(this, 'Size Job', {
			task: new tasks.InvokeFunction(sizeFunc),
		});

		const parquetJob = new sfn.Task(this, 'Parquet Job', {
			task: new tasks.InvokeFunction(parquetFunc),
		});

		const fileTooBig = new sfn.Fail(this, 'File too big', {
			cause: 'File too big',
			error: 'parquetJob returned FAILED',
		});

		const athenaInput = new sfn.Pass(this, 'Input for Athena', {
			result: sfn.Result.fromArray([{
				"result_bucket": athenaBucket.bucketName,
                "query": `
                    SELECT cl.boss_fight_uuid, cl.damage, cl.caster_name, cl.caster_id, ei.encounter_id
                    FROM (SELECT SUM(actual_amount) AS damage, caster_name, caster_id, boss_fight_uuid
                          FROM  "wowmate"."combatlogs"
                          WHERE caster_type LIKE '0x5%' AND caster_name != 'nil' 
                          GROUP BY caster_name, caster_id, boss_fight_uuid
                          ) AS cl
                    JOIN (SELECT encounter_id, boss_fight_uuid
                          FROM "wowmate"."combatlogs"
                          WHERE event_type = 'ENCOUNTER_START'
                          GROUP BY encounter_id, boss_fight_uuid) AS ei
                          ON cl.boss_fight_uuid = ei.boss_fight_uuid
                          
                    ORDER BY encounter_id, damage DESC
                `,
				"region": "eu-central-1",
				"database": "wowmate"
			}]),
		})

		const athenaJob = new sfn.Task(this, 'Athena Job', {
			task: new tasks.InvokeFunction(athenaFunc),
		});

		const setWaitTimeJob = new sfn.Pass(this, 'Set wait time', {
			inputPath: '$',
			result: sfn.Result.fromNumber(3),
			resultPath: '$.wait_time',
			outputPath: '$'
		})

		const waitX = new sfn.Wait(this, 'Wait X Seconds', {
			time: sfn.WaitTime.secondsPath('$.wait_time')
		});
	
		const checkJob = new sfn.Task(this, 'Check Athena Status Job', {
			task: new tasks.InvokeFunction(checkFunc)
		})
		
		checkJob.addRetry({
			interval: Duration.seconds(2),
			maxAttempts: 10,
		})

		const duplicateJob = new sfn.Pass(this, 'Detect duplicate combatlogs', {
			inputPath: '$',
			result: sfn.Result.fromBoolean(false),
			resultPath: '$.duplicate',
			outputPath: '$'
		})

		const duplicateLog = new sfn.Fail(this, 'Combatlog duplicate', {
			//cause: 'File too big',
			error: 'combatlog is already in the database',
		});

		const dynamodbJob = new sfn.Task(this, 'DynamoDB Job', {
			task: new tasks.InvokeFunction(impFunc),
		});

		// const parallel = new sfn.Parallel(this, 'Parallel Queries');

		// const sumQueryInput = new sfn.Pass(this, 'Input for Sum', {
		// 	result: sfn.Result.fromArray([{
		// 		"result_bucket": athenaBucket.bucketName,
		// 		"query": `SELECT sum(n1) as sum FROM combatlogs;`,
		// 		"region": "eu-central-1",
		// 		"table": "wowmate"
		// 	}]),
		// })

		// const athenaJob2 = new sfn.Task(this, 'Athena Job2', {
		// 	task: new tasks.InvokeFunction(athenaFunc),
		// });

		// const setWaitTimeJob2 = new sfn.Pass(this, 'Set wait time2', {
		// 	inputPath: '$',
		// 	result: sfn.Result.fromNumber(3),
		// 	resultPath: '$.wait_time',
		// 	outputPath: '$'
		// })

		// const waitX2 = new sfn.Wait(this, 'Wait X Seconds2', {
		// 	time: sfn.WaitTime.secondsPath('$.wait_time')
		// });
	
		// const checkJob2 = new sfn.Task(this, 'Check Athena Status Job2', {
		// 	task: new tasks.InvokeFunction(checkFunc),
		// })
		
		// checkJob2.addRetry({
		// 	interval: Duration.seconds(2),
		// 	maxAttempts: 10,
		// })

		// const dynamodbJob2 = new sfn.Task(this, 'DynamoDB Job2', {
		// 	task: new tasks.InvokeFunction(imp2Func),
		// });

		// const averageQueryInput = new sfn.Pass(this, 'Input for Average', {
		// 	result: sfn.Result.fromArray([{
		// 		"result_bucket": athenaBucket.bucketName,
		// 		"query": `SELECT avg(n1) as avg FROM combatlogs;`,
		// 		"region": "eu-central-1",
		// 		"database": "wowmate"
		// 	}]),
		// })

		// const athenaJob3 = new sfn.Task(this, 'Athena Job3', {
		// 	task: new tasks.InvokeFunction(athenaFunc),
		// });

		// const setWaitTimeJob3 = new sfn.Pass(this, 'Set wait time3', {
		// 	inputPath: '$',
		// 	result: sfn.Result.fromNumber(3),
		// 	resultPath: '$.wait_time',
		// 	outputPath: '$'
		// })

		// const waitX3 = new sfn.Wait(this, 'Wait X Seconds3', {
		// 	time: sfn.WaitTime.secondsPath('$.wait_time')
		// });
	
		// const checkJob3 = new sfn.Task(this, 'Check Athena Status Job3', {
		// 	task: new tasks.InvokeFunction(checkFunc),
		// })
		
		// checkJob3.addRetry({
		// 	interval: Duration.seconds(2),
		// 	maxAttempts: 10,
		// })

		// const dynamodbJob3 = new sfn.Task(this, 'DynamoDB Job3', {
		// 	task: new tasks.InvokeFunction(imp2Func),
		// });

		//to athenaBucket stuff as state machine fragment
		//https://docs.aws.amazon.com/cdk/api/latest/docs/aws-stepfunctions-readme.html#state-machine-fragments
		const sfunc = new sfn.StateMachine(this, 'StateMachine', {
			definition: sizeJob
			.next(new sfn.Choice(this, 'Check file sizeFunc')
				.when(sfn.Condition.numberGreaterThan('$.file_size', 400), fileTooBig)
				.otherwise(parquetJob
					.next(athenaInput)
					.next(athenaJob)
					.next(setWaitTimeJob)
					.next(waitX)
					.next(checkJob)
					.next(duplicateJob)
					.next(new sfn.Choice(this, 'Check if log already exists')
						.when(sfn.Condition.stringEquals('$.duplicate', 'true'), duplicateLog)
						.otherwise(dynamodbJob
							// .next(parallel
							// 	.branch(averageQueryInput
							// 		.next(athenaJob2)
							// 		.next(setWaitTimeJob2)
							// 		.next(waitX2)
							// 		.next(checkJob2)
							// 		.next(dynamodbJob2)									
							// 	).branch(sumQueryInput
							// 		.next(athenaJob3)
							// 		.next(setWaitTimeJob3)
							// 		.next(waitX3)
							// 		.next(checkJob3)
							// 		.next(dynamodbJob3)
							// 	)
							// )
						)
					)
				)
			),
		});

		uploadBucket.onCloudTrailPutObject('cwEvent', {
			target: new targets.SfnStateMachine(sfunc),	
		}).addEventPattern({
            detail: {
                eventName: ['CompleteMultipartUpload'],
            },
        });
	}
}
