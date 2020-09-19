package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	lambdaservice "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	_ "github.com/lib/pq"
)

//DatabasesCredentials are the data to log into the db
type DatabasesCredentials struct {
	DatabaseName string `json:"dbname"`
	Password     string `json:"password"`
	UserName     string `json:"username"`
	Host         string `json:"host"`
}

//S3Event is the data that come from s3 and contains all the information about the event
type S3Event struct {
	Records []struct {
		EventVersion string    `json:"eventVersion"`
		EventSource  string    `json:"eventSource"`
		AwsRegion    string    `json:"awsRegion"`
		EventTime    time.Time `json:"eventTime"`
		EventName    string    `json:"eventName"`
		UserIdentity struct {
			PrincipalID string `json:"principalId"`
		} `json:"userIdentity"`
		RequestParameters struct {
			SourceIPAddress string `json:"sourceIPAddress"`
		} `json:"requestParameters"`
		ResponseElements struct {
			XAmzRequestID string `json:"x-amz-request-id"`
			XAmzID2       string `json:"x-amz-id-2"`
		} `json:"responseElements"`
		S3 struct {
			S3SchemaVersion string `json:"s3SchemaVersion"`
			ConfigurationID string `json:"configurationId"`
			Bucket          struct {
				Name          string `json:"name"`
				OwnerIdentity struct {
					PrincipalID string `json:"principalId"`
				} `json:"ownerIdentity"`
				Arn string `json:"arn"`
			} `json:"bucket"`
			Object struct {
				Key       string `json:"key"`
				Size      int    `json:"size"`
				ETag      string `json:"eTag"`
				Sequencer string `json:"sequencer"`
			} `json:"object"`
		} `json:"s3"`
	} `json:"Records"`
}

//SQSEvent is all the data that gets passed into the lambda from the q
type SQSEvent struct {
	Records []struct {
		MessageID     string `json:"messageId"`
		ReceiptHandle string `json:"receiptHandle"`
		Body          string `json:"body"`
		Attributes    struct {
			ApproximateReceiveCount          string `json:"ApproximateReceiveCount"`
			SentTimestamp                    string `json:"SentTimestamp"`
			SenderID                         string `json:"SenderId"`
			ApproximateFirstReceiveTimestamp string `json:"ApproximateFirstReceiveTimestamp"`
		} `json:"attributes"`
		MessageAttributes struct {
		} `json:"messageAttributes"`
		Md5OfBody      string `json:"md5OfBody"`
		EventSource    string `json:"eventSource"`
		EventSourceARN string `json:"eventSourceARN"`
		AwsRegion      string `json:"awsRegion"`
	} `json:"Records"`
}

func handler(e SQSEvent) error {
	secretArn := os.Getenv("SECRET_ARN")
	if secretArn == "" {
		return fmt.Errorf("secret arn env var is empty")
	}

	dbEndpoint := os.Getenv("DB_ENDPOINT")
	if dbEndpoint == "" {
		return fmt.Errorf("db endpoint env var is empty")
	}

	summaryLambdaName := os.Getenv("SUMMARY_LAMBDA_NAME")
	if summaryLambdaName == "" {
		return fmt.Errorf("summary lambda name env var is empty")
	}

	sess, err := session.NewSession()
	if err != nil {
		log.Println(err.Error())
		log.Println("failed to create new session")
		return err
	}

	//TODO: should move get secret outside of handler, because it dosn't need to run on every invocation
	svc := secretsmanager.New(sess)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretArn),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeDecryptionFailure:
				// Secrets Manager can't decrypt the protected secret text using the provided KMS key.
				log.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
				return err

			case secretsmanager.ErrCodeInternalServiceError:
				// An error occurred on the server side.
				log.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
				return err

			case secretsmanager.ErrCodeInvalidParameterException:
				// You provided an invalid value for a parameter.
				log.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
				return err

			case secretsmanager.ErrCodeInvalidRequestException:
				// You provided a parameter value that is not valid for the current state of the resource.
				log.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
				return err

			case secretsmanager.ErrCodeResourceNotFoundException:
				// We can't find the resource that you asked for.
				log.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
				return err
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
			return err
		}
	}

	var creds = DatabasesCredentials{}
	err = json.Unmarshal([]byte(*result.SecretString), &creds)
	if err != nil {
		return err
	}

	connStr := fmt.Sprintf(
		"user=%v dbname=%v sslmode=verify-full host=%v password=%v port=5432",
		creds.UserName,
		creds.DatabaseName,
		// creds.Host,
		//NOTE:
		//we aren't using the host from the secret, because we are passing it in
		//through the env var, because it could be a proxy or a read only endpoint
		dbEndpoint,
		creds.Password,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()
	log.Println("openend connection")
	log.Printf("number of messages: %v", len(e.Records))

	for _, record := range e.Records {
		s3 := S3Event{}
		err = json.Unmarshal([]byte(record.Body), &s3)
		if err != nil {
			log.Println("Unmarshal sqs json failed")
			return err
		}

		if len(s3.Records) > 1 {
			return fmt.Errorf("Failed: the S3 event contains more than 1 element, not sure how that would happen")
		}
		q := fmt.Sprintf(`
				SELECT aws_s3.table_import_from_s3(
					'combatlogs',
					'column_uuid,upload_uuid,unsupported,combatlog_uuid,boss_fight_uuid,mythicplus_uuid,event_type,version,advanced_log_enabled,dungeon_name,dungeon_id,key_unkown_1,key_level,key_array,key_duration,encounter_id,encounter_name,encounter_unkown_1,encounter_unkown_2,killed,caster_id,caster_name,caster_type,source_flag,target_id,target_name,target_type,dest_flag,spell_id,spell_name,spell_type,extra_spell_id,extra_spell_name,extra_school,aura_type,another_player_id,d0,d1,d2,d3,d4,d5,d6,d7,d8,d9,d10,d11,d12,d13,damage_unknown_14,actual_amount,base_amount,overhealing,overkill,school,resisted,blocked,absorbed,critical,glancing,crushing,is_offhand',
					'(format csv, DELIMITER '','', HEADER true)',
					'(%v,%v,us-east-1)');
			`, s3.Records[0].S3.Bucket.Name, s3.Records[0].S3.Object.Key)

		rows, err := db.Query(q)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				//NOTE: this can happen because sometimes the import takes 10x the normal time and
				//		the impport finished, after the lambda timed out
				//		which causes the lambda to be triggered again from SQS, instead of failing
				//		till it lands in the DLQ we will just delete it from the queue,
				//		because it got actually imported to the DB
				log.Printf("duplicate key for %v, error msg: %v ", s3.Records[0].S3.Object.Key, err.Error())
				return nil
			}
			//IMPROVE: check if error is:
			//pq: duplicate key value violates unique constraint "combatlogs_pkey": Error
			//if yes don't fail just continue
			return err
		}

		defer rows.Close()

		for rows.Next() {
			var s string

			err = rows.Scan(&s)
			if err != nil {
				log.Println(err.Error())
				return err
			}
			log.Printf("import query successfull: %v", s)

			svc := lambdaservice.New(sess)
			input := &lambdaservice.InvokeInput{
				FunctionName:   aws.String(summaryLambdaName),
				InvocationType: aws.String("Event"),
				Payload:        []byte(fmt.Sprintf("{\"filename\":\"%v\"}", s3.Records[0].S3.Object.Key)),
			}

			result, err := svc.Invoke(input)
			if err != nil {
				if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					case lambdaservice.ErrCodeServiceException:
						log.Println(lambdaservice.ErrCodeServiceException, aerr.Error())
					case lambdaservice.ErrCodeResourceNotFoundException:
						log.Println(lambdaservice.ErrCodeResourceNotFoundException, aerr.Error())
					case lambdaservice.ErrCodeInvalidRequestContentException:
						log.Println(lambdaservice.ErrCodeInvalidRequestContentException, aerr.Error())
					case lambdaservice.ErrCodeRequestTooLargeException:
						log.Println(lambdaservice.ErrCodeRequestTooLargeException, aerr.Error())
					case lambdaservice.ErrCodeUnsupportedMediaTypeException:
						log.Println(lambdaservice.ErrCodeUnsupportedMediaTypeException, aerr.Error())
					case lambdaservice.ErrCodeTooManyRequestsException:
						log.Println(lambdaservice.ErrCodeTooManyRequestsException, aerr.Error())
					case lambdaservice.ErrCodeInvalidParameterValueException:
						log.Println(lambdaservice.ErrCodeInvalidParameterValueException, aerr.Error())
					case lambdaservice.ErrCodeEC2UnexpectedException:
						log.Println(lambdaservice.ErrCodeEC2UnexpectedException, aerr.Error())
					case lambdaservice.ErrCodeSubnetIPAddressLimitReachedException:
						log.Println(lambdaservice.ErrCodeSubnetIPAddressLimitReachedException, aerr.Error())
					case lambdaservice.ErrCodeENILimitReachedException:
						log.Println(lambdaservice.ErrCodeENILimitReachedException, aerr.Error())
					case lambdaservice.ErrCodeEFSMountConnectivityException:
						log.Println(lambdaservice.ErrCodeEFSMountConnectivityException, aerr.Error())
					case lambdaservice.ErrCodeEFSMountFailureException:
						log.Println(lambdaservice.ErrCodeEFSMountFailureException, aerr.Error())
					case lambdaservice.ErrCodeEFSMountTimeoutException:
						log.Println(lambdaservice.ErrCodeEFSMountTimeoutException, aerr.Error())
					case lambdaservice.ErrCodeEFSIOException:
						log.Println(lambdaservice.ErrCodeEFSIOException, aerr.Error())
					case lambdaservice.ErrCodeEC2ThrottledException:
						log.Println(lambdaservice.ErrCodeEC2ThrottledException, aerr.Error())
					case lambdaservice.ErrCodeEC2AccessDeniedException:
						log.Println(lambdaservice.ErrCodeEC2AccessDeniedException, aerr.Error())
					case lambdaservice.ErrCodeInvalidSubnetIDException:
						log.Println(lambdaservice.ErrCodeInvalidSubnetIDException, aerr.Error())
					case lambdaservice.ErrCodeInvalidSecurityGroupIDException:
						log.Println(lambdaservice.ErrCodeInvalidSecurityGroupIDException, aerr.Error())
					case lambdaservice.ErrCodeInvalidZipFileException:
						log.Println(lambdaservice.ErrCodeInvalidZipFileException, aerr.Error())
					case lambdaservice.ErrCodeKMSDisabledException:
						log.Println(lambdaservice.ErrCodeKMSDisabledException, aerr.Error())
					case lambdaservice.ErrCodeKMSInvalidStateException:
						log.Println(lambdaservice.ErrCodeKMSInvalidStateException, aerr.Error())
					case lambdaservice.ErrCodeKMSAccessDeniedException:
						log.Println(lambdaservice.ErrCodeKMSAccessDeniedException, aerr.Error())
					case lambdaservice.ErrCodeKMSNotFoundException:
						log.Println(lambdaservice.ErrCodeKMSNotFoundException, aerr.Error())
					case lambdaservice.ErrCodeInvalidRuntimeException:
						log.Println(lambdaservice.ErrCodeInvalidRuntimeException, aerr.Error())
					case lambdaservice.ErrCodeResourceConflictException:
						log.Println(lambdaservice.ErrCodeResourceConflictException, aerr.Error())
					case lambdaservice.ErrCodeResourceNotReadyException:
						log.Println(lambdaservice.ErrCodeResourceNotReadyException, aerr.Error())
					default:
						log.Println(aerr.Error())
					}
				} else {
					log.Println(err.Error())
				}
				return err
			}

			log.Println(result)
		}
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
