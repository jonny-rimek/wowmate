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
// import events = require('@aws-cdk/aws-events');
// import { Result } from '@aws-cdk/aws-stepfunctions';

export class WowmateStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

		const trail = new cloudtrail.Trail(this, 'CloudTrail', {
			sendToCloudWatchLogs: true,
			managementEvents: cloudtrail.ReadWriteType.WRITE_ONLY,
		});

		const upload = new s3.Bucket(this, 'Upload', {
			removalPolicy: RemovalPolicy.DESTROY,
			blockPublicAccess: BlockPublicAccess.BLOCK_ALL,
		})

		trail.addS3EventSelector([upload.bucketArn + "/"], {
			readWriteType: cloudtrail.ReadWriteType.WRITE_ONLY,
		})

		const db = new ddb.Table(this, 'DDB', {
			partitionKey: { name: 'name', type: ddb.AttributeType.STRING },
			removalPolicy: RemovalPolicy.DESTROY,
		})

		const parquet = new s3.Bucket(this, 'Parquet', {
			removalPolicy: RemovalPolicy.DESTROY,
			blockPublicAccess: BlockPublicAccess.BLOCK_ALL,
		})

		const athena = new s3.Bucket(this, 'Athena', {
			removalPolicy: RemovalPolicy.DESTROY,
			blockPublicAccess: BlockPublicAccess.BLOCK_ALL,
		})

		const size = new lambda.Function(this, 'Size', {
			code: lambda.Code.asset("handler/size"),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: Duration.seconds(3),
        })

		const func = new lambda.Function(this, 'test', {
			code: lambda.Code.asset("handler/parquet"),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: Duration.seconds(10),
			environment: {TARGET_BUCKET_NAME: parquet.bucketName}
            //TODO: add parquet bucket name as env variable
		})

		const athenaFunc = new lambda.Function(this, 'AthenaFunc', {
			code: lambda.Code.asset("handler/athena"),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: Duration.seconds(3),
		})

		const check = new lambda.Function(this, 'Check', {
			code: lambda.Code.asset("handler/check"),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: Duration.seconds(3),
		})
	
		const imp = new lambda.Function(this, 'Import', {
			code: lambda.Code.asset("handler/import"),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: Duration.seconds(10),
            //TODO: add parquet bucket name as env variable
		})

		const imp2 = new lambda.Function(this, 'Import2', {
			code: lambda.Code.asset("handler/import2"),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: Duration.seconds(10),
            //TODO: add parquet bucket name as env variable
		})

		upload.grantRead(size)

		upload.grantRead(func)
		parquet.grantPut(func)

		parquet.grantRead(athenaFunc)
		athena.grantWrite(athenaFunc)

		const athenaPolicy = new iam.PolicyStatement({
			effect: Effect.ALLOW,
			actions: [ 
				'athena:*',
				'glue:*'
			],
			resources: [
				'*'
			],
		})
		athenaFunc.addToRolePolicy(athenaPolicy)

		check.addToRolePolicy(athenaPolicy)

		athena.grantRead(imp)
		db.grantWriteData(imp)

		athena.grantRead(imp2)
		db.grantWriteData(imp2)

		const sizeJob = new sfn.Task(this, 'Size Job', {
			task: new tasks.InvokeFunction(size),
		});

		const parquetJob = new sfn.Task(this, 'Parquet Job', {
			task: new tasks.InvokeFunction(func),
		});

		const fileTooBig = new sfn.Fail(this, 'File too big', {
			cause: 'File too big',
			error: 'parquetJob returned FAILED',
		});

		const athenaInput = new sfn.Pass(this, 'Input for Athena', {
			result: sfn.Result.fromArray([{
				"result_bucket": "sfnstack-athenabc8ba882-d8g4by741f1i",
				"query": `SELECT name, sum(n1) as sum FROM sfntest GROUP BY name;`,
				"region": "eu-central-1",
				"table": "sfn"
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
			task: new tasks.InvokeFunction(check),
		});
		
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
			task: new tasks.InvokeFunction(imp),
		});

		const parallel = new sfn.Parallel(this, 'Parallel Queries');

		const sumQueryInput = new sfn.Pass(this, 'Input for Sum', {
			result: sfn.Result.fromArray([{
				"result_bucket": "sfnstack-athenabc8ba882-d8g4by741f1i",
				"query": `SELECT sum(n1) as sum FROM sfntest;`,
				"region": "eu-central-1",
				"table": "sfn"
			}]),
		})

		const athenaJob2 = new sfn.Task(this, 'Athena Job2', {
			task: new tasks.InvokeFunction(athenaFunc),
		});

		const setWaitTimeJob2 = new sfn.Pass(this, 'Set wait time2', {
			inputPath: '$',
			result: sfn.Result.fromNumber(3),
			resultPath: '$.wait_time',
			outputPath: '$'
		})

		const waitX2 = new sfn.Wait(this, 'Wait X Seconds2', {
			time: sfn.WaitTime.secondsPath('$.wait_time')
		});
	
		const checkJob2 = new sfn.Task(this, 'Check Athena Status Job2', {
			task: new tasks.InvokeFunction(check),
		});
		
		checkJob2.addRetry({
			interval: Duration.seconds(2),
			maxAttempts: 10,
		})

		const dynamodbJob2 = new sfn.Task(this, 'DynamoDB Job2', {
			task: new tasks.InvokeFunction(imp2),
		});

		const averageQueryInput = new sfn.Pass(this, 'Input for Average', {
			result: sfn.Result.fromArray([{
				"result_bucket": "sfnstack-athenabc8ba882-d8g4by741f1i",
				"query": `SELECT avg(n1) as avg FROM sfntest;`,
				"region": "eu-central-1",
				"table": "sfn"
			}]),
		})

		const athenaJob3 = new sfn.Task(this, 'Athena Job3', {
			task: new tasks.InvokeFunction(athenaFunc),
		});

		const setWaitTimeJob3 = new sfn.Pass(this, 'Set wait time3', {
			inputPath: '$',
			result: sfn.Result.fromNumber(3),
			resultPath: '$.wait_time',
			outputPath: '$'
		})

		const waitX3 = new sfn.Wait(this, 'Wait X Seconds3', {
			time: sfn.WaitTime.secondsPath('$.wait_time')
		});
	
		const checkJob3 = new sfn.Task(this, 'Check Athena Status Job3', {
			task: new tasks.InvokeFunction(check),
		});
		
		checkJob3.addRetry({
			interval: Duration.seconds(2),
			maxAttempts: 10,
		})

		const dynamodbJob3 = new sfn.Task(this, 'DynamoDB Job3', {
			task: new tasks.InvokeFunction(imp2),
		});

		//to athena stuff as state machine fragment
		//https://docs.aws.amazon.com/cdk/api/latest/docs/aws-stepfunctions-readme.html#state-machine-fragments
		const sfunc = new sfn.StateMachine(this, 'StateMachine', {
			definition: sizeJob
			.next(new sfn.Choice(this, 'Check file size')
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
							.next(parallel
								.branch(averageQueryInput
									.next(athenaJob2)
									.next(setWaitTimeJob2)
									.next(waitX2)
									.next(checkJob2)
									.next(dynamodbJob2)									
								).branch(sumQueryInput
									.next(athenaJob3)
									.next(setWaitTimeJob3)
									.next(waitX3)
									.next(checkJob3)
									.next(dynamodbJob3)
								)
							)
						)
					)
				)
			),
		});

		upload.onCloudTrailPutObject('cwEvent', {
			target: new targets.SfnStateMachine(sfunc),	
		}).addEventPattern({
            detail: {
                eventName: ['CompleteMultipartUpload'],
            },
        });

/*
*/
  }
}
