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
import glue = require('@aws-cdk/aws-glue');
import events = require('@aws-cdk/aws-events');
import { RetentionDays } from '@aws-cdk/aws-logs';

interface DatabaseProps extends cdk.StackProps {
	dynamoDB: ddb.ITable;
}

export class Upload extends cdk.Construct {
	constructor(scope: cdk.Construct, id: string, props: DatabaseProps) {
		super(scope, id);

		const db = props.dynamoDB
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

		trail.addS3EventSelector([uploadBucket.bucketArn + "/new/"], {
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

		const athenaDB = new glue.Database(this, 'Wowmate', {
			databaseName: 'wowmate',
		})

		new glue.Table(this, 'Combatlogs', {
			database: athenaDB,
			tableName: 'combatlogs',
			s3Prefix: '',
			bucket: parquetBucket,
			dataFormat: glue.DataFormat.PARQUET,
			storedAsSubDirectories: true,
			// compressed: true,
			columns: [{
				name: 'upload_uuid',
				type: glue.Schema.STRING,
			}, {
				name: 'unsupported',
				type: glue.Schema.BOOLEAN,
			}, {
				name: 'combatlog_uuid',
				type: glue.Schema.STRING,
			}, {
				name: 'boss_fight_uuid',
				type: glue.Schema.STRING,
			}, {
				name: 'mythicplus_uuid',
				type: glue.Schema.STRING,
			}, {
				name: 'column_uuid',
				type: glue.Schema.STRING,
			}, {
				name: 'timestamp',
				type: glue.Schema.TIMESTAMP,
			}, {
				name: 'event_type',
				type: glue.Schema.STRING,
			}, {
				name: 'version',
				type: glue.Schema.INTEGER,
			}, {
				name: 'advanced_log_enabled',
				type: glue.Schema.INTEGER,
			}, {
				name: 'dungeon_name',
				type: glue.Schema.STRING,
			}, {
				name: 'dungeon_id',
				type: glue.Schema.INTEGER,
			}, {
				name: 'key_unknown_1',
				type: glue.Schema.INTEGER,
			}, {
				name: 'key_level',
				type: glue.Schema.INTEGER,
			}, {
				name: 'key_array',
				type: glue.Schema.STRING,
			}, {
				name: 'key_duration',
				type: glue.Schema.BIG_INT,
			}, {
				name: 'encounter_id',
				type: glue.Schema.INTEGER,
			}, {
				name: 'encounter_name',
				type: glue.Schema.STRING,
			}, {
				name: 'encounter_unkown_1',
				type: glue.Schema.INTEGER,
			}, {
				name: 'encounter_unkown_2',
				type: glue.Schema.INTEGER,
			}, {
				name: 'killed',
				type: glue.Schema.INTEGER,
			}, {
				name: 'caster_id',
				type: glue.Schema.STRING,
			}, {
				name: 'caster_name',
				type: glue.Schema.STRING,
			}, {
				name: 'caster_type',
				type: glue.Schema.STRING,
			}, {
				name: 'source_flag',
				type: glue.Schema.STRING,
			}, {
				name: 'target_id',
				type: glue.Schema.STRING,
			}, {
				name: 'target_name',
				type: glue.Schema.STRING,
			}, {
				name: 'target_type',
				type: glue.Schema.STRING,
			}, {
				name: 'dest_flag',
				type: glue.Schema.STRING,
			}, {
				name: 'spell_id',
				type: glue.Schema.INTEGER,
			}, {
				name: 'spell_name',
				type: glue.Schema.STRING,
			}, {
				name: 'spell_type',
				type: glue.Schema.STRING,
			}, {
				name: 'extra_spell_id',
				type: glue.Schema.INTEGER,
			}, {
				name: 'extra_spell_name',
				type: glue.Schema.STRING,
			}, {
				name: 'extra_school',
				type: glue.Schema.STRING,
			}, {
				name: 'aura_type',
				type: glue.Schema.STRING,
			}, {
				name: 'another_player_id',
				type: glue.Schema.STRING,
			}, {
				name: 'd0',
				type: glue.Schema.STRING,
			}, {
				name: 'd1',
				type: glue.Schema.BIG_INT,
			}, {
				name: 'd2',
				type: glue.Schema.BIG_INT,
			}, {
				name: 'd3',
				type: glue.Schema.BIG_INT,
			}, {
				name: 'd4',
				type: glue.Schema.BIG_INT,
			}, {
				name: 'd5',
				type: glue.Schema.BIG_INT,
			}, {
				name: 'd6',
				type: glue.Schema.BIG_INT,
			}, {
				name: 'd7',
				type: glue.Schema.BIG_INT,
			}, {
				name: 'd8',
				type: glue.Schema.BIG_INT,
			}, {
				name: 'd9',
				type: glue.Schema.STRING,
			}, {
				name: 'd10',
				type: glue.Schema.STRING,
			}, {
				name: 'd11',
				type: glue.Schema.STRING,
			}, {
				name: 'd12',
				type: glue.Schema.STRING,
			}, {
				name: 'd13',
				type: glue.Schema.STRING,
			}, {
				name: 'damage_unknown_14',
				type: glue.Schema.STRING,
			}, {
				name: 'actual_amount',
				type: glue.Schema.BIG_INT,
			}, {
				name: 'base_amount',
				type: glue.Schema.BIG_INT,
			}, {
				name: 'overhealing',
				type: glue.Schema.BIG_INT,
			}, {
				name: 'overkill',
				type: glue.Schema.STRING,
			}, {
				name: 'school',
				type: glue.Schema.STRING,
			}, {
				name: 'resisted',
				type: glue.Schema.STRING,
			}, {
				name: 'blocked',
				type: glue.Schema.STRING,
			}, {
				name: 'absorbed',
				type: glue.Schema.BIG_INT,
			}, {
				name: 'critical',
				type: glue.Schema.STRING,
			}, {
				name: 'glancing',
				type: glue.Schema.STRING,
			}, {
				name: 'crushing',
				type: glue.Schema.STRING,
			}, {
				name: 'is_offhand',
				type: glue.Schema.STRING,
			}],
			partitionKeys: [{
				name: 'year',
				type: glue.Schema.STRING,
			}, {
				name: 'month',
				type: glue.Schema.STRING,
			}, {
				name: 'day',
				type: glue.Schema.STRING,
			}, {
				name: 'hour',
				type: glue.Schema.STRING,
			}, {
				name: 'minute',
				type: glue.Schema.STRING,
			}],
		});
		
		//IAM for UPLOAD LAMBDA

		//TODO: tighten athena IAM permissions
		//permission are from the aws docs https://docs.aws.amazon.com/athena/latest/ug/example-policies-workgroup.html
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
				"glue:GetTable",
				"glue:GetDatabase",
				"glue:GetPartition",
				"glue:GetPartitions",
				//below I added all for partitions lambda, not sure which of those I actually need
				"glue:UpdateDatabase",
				"glue:UpdatePartition",
				"glue:UpdateTable",
				"glue:BatchCreatePartition",
				"s3:ListAllMyBuckets",
				"s3:ListBucket",
				"s3:PutObject",
				"s3:GetObject",
				"s3:AbortMultipartUpload",
				"s3:ListMultipartUploadParts",
				"s3:GetBucketLocation",
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
				'arn:aws:athena:us-east-1:940880032268:workgroup/primary'
			],
		})

		//UPLOAD LAMBDA
		const sizeFunc = new lambda.Function(this, 'Size', {
			code: lambda.Code.asset('services/upload/size'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 128,
			timeout: Duration.seconds(3),
			logRetention: RetentionDays.ONE_MONTH,
			tracing: lambda.Tracing.ACTIVE,
		})
		uploadBucket.grantRead(sizeFunc)

		const parquetFunc = new lambda.Function(this, 'ParquetFunc', {
			code: lambda.Code.asset('services/upload/parquet'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: Duration.seconds(720),
			environment: {TARGET_BUCKET_NAME: parquetBucket.bucketName},
			logRetention: RetentionDays.ONE_MONTH,
			tracing: lambda.Tracing.ACTIVE,
		})
		uploadBucket.grantReadWrite(parquetFunc)
		parquetBucket.grantPut(parquetFunc)

		const athenaRepairFunc = new lambda.Function(this, 'AthenaRepairFunc', {
			code: lambda.Code.asset('services/upload/athena-repair'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: Duration.seconds(3),
			environment: {
				"RESULT_BUCKET": athenaBucket.bucketName,
				"REGION": "us-east-1",
				"ATHENA_DATABASE": "wowmate",//TODO: don't hardcode values
			},
			logRetention: RetentionDays.ONE_MONTH,
			tracing: lambda.Tracing.ACTIVE,
		})
		parquetBucket.grantRead(athenaRepairFunc)
		athenaBucket.grantReadWrite(athenaRepairFunc)
		athenaRepairFunc.addToRolePolicy(athenaGeneralPolicy)
		athenaRepairFunc.addToRolePolicy(athenaWorkgroupPolicy)

		const athenaFunc = new lambda.Function(this, 'AthenaFunc', {
			code: lambda.Code.asset('services/upload/athena'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: Duration.seconds(3),
			environment: {
				"RESULT_BUCKET": athenaBucket.bucketName,
				"REGION": "us-east-1",
				"ATHENA_DATABASE": "wowmate", //TODO: don't hardcode values
			},
			logRetention: RetentionDays.ONE_MONTH,
			tracing: lambda.Tracing.ACTIVE,
		})
		parquetBucket.grantRead(athenaFunc)
		athenaBucket.grantReadWrite(athenaFunc)
		athenaFunc.addToRolePolicy(athenaGeneralPolicy)
		athenaFunc.addToRolePolicy(athenaWorkgroupPolicy)

		const checkFunc = new lambda.Function(this, 'Check', {
			code: lambda.Code.asset('services/upload/check'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: Duration.seconds(3),
			logRetention: RetentionDays.ONE_MONTH,
			tracing: lambda.Tracing.ACTIVE,
		})
		athenaBucket.grantWrite(checkFunc)
		checkFunc.addToRolePolicy(athenaGeneralPolicy)
		checkFunc.addToRolePolicy(athenaWorkgroupPolicy)
	
		const impFunc = new lambda.Function(this, 'Import', {
			code: lambda.Code.asset('services/upload/import'),
			handler: 'main',
			runtime: lambda.Runtime.GO_1_X,
			memorySize: 3008,
			timeout: Duration.seconds(10),//NOTE: its so long because im waiting for delete confirmation
			environment: {
				DDB_NAME: db.tableName,
				// LOG_LEVEL: 'prod',
			},
			logRetention: RetentionDays.ONE_MONTH,
			tracing: lambda.Tracing.ACTIVE,
		})
		athenaBucket.grantRead(impFunc)
		parquetBucket.grantDelete(impFunc)
		db.grantReadWriteData(impFunc)

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

		const athenaRepairJob = new sfn.Task(this, 'AthenaRepairJob', {
			task: new tasks.InvokeFunction(athenaRepairFunc),
		});
		const setWaitTimeJobY = new sfn.Pass(this, 'Set wait timeY', {
			inputPath: '$',
			result: sfn.Result.fromNumber(3),
			resultPath: '$.wait_time',
			outputPath: '$'
		})

		const waitY = new sfn.Wait(this, 'Wait Y Seconds', {
			time: sfn.WaitTime.secondsPath('$.wait_time')
		});
		const checkJobY = new sfn.Task(this, 'Check Athena Status JobY', {
			task: new tasks.InvokeFunction(checkFunc)
		})
		checkJobY.addRetry({
			interval: Duration.seconds(2),
			maxAttempts: 10,
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
			error: 'combatlog is already in the database',
		});

		const dynamodbJob = new sfn.Task(this, 'DynamoDB Job', {
			task: new tasks.InvokeFunction(impFunc),
		});

		//to athenaBucket stuff as state machine fragment
		//https://docs.aws.amazon.com/cdk/api/latest/docs/aws-stepfunctions-readme.html#state-machine-fragments
		const sfunc = new sfn.StateMachine(this, 'StateMachine', {
			definition: sizeJob
			.next(new sfn.Choice(this, 'Check file sizeFunc')
				//NOTE: I know I can handle 11MB(gzipped), but 21MB already fails with out of memory fatal error
				.when(sfn.Condition.numberGreaterThan('$.file_size', 16), fileTooBig)
				.otherwise(parquetJob
					.next(athenaRepairJob)
					.next(setWaitTimeJobY)
					.next(waitY)
					.next(checkJobY)
					.next(athenaJob)
					.next(setWaitTimeJob)
					.next(waitX)
					.next(checkJob)
					.next(duplicateJob)
					.next(new sfn.Choice(this, 'Check if log already exists')
						.when(sfn.Condition.stringEquals('$.duplicate', 'true'), duplicateLog)
						.otherwise(dynamodbJob
						)
					)
				)
			),
			// IMPROVE: use express functions
			// caveats, need to activate extra logging
			// stateMachineType: StateMachineType.EXPRESS,
		});

		uploadBucket.onCloudTrailPutObject('cwEvent', {
			target: new targets.SfnStateMachine(sfunc),	
		}).addEventPattern({
            detail: {
                eventName: ['CompleteMultipartUpload'],
            },
		});
		
		//Client S3 upload IAM
		const role = new iam.Role(this, 'ClientRole', {
			assumedBy: new iam.AccountPrincipal('940880032268'),
		})

		role.addToPolicy(new iam.PolicyStatement({
			resources: [uploadBucket.bucketArn],
			actions: ['s3:PutObject'],
		}))
	}
}
