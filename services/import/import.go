package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
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
		log.Println(err.Error())
		return err
	}

	connStr := fmt.Sprintf(
		"user=%v dbname=%v sslmode=verify-full host=%v password=%v port=5432",
		creds.UserName,
		creds.DatabaseName,
		// creds.Host,
		dbEndpoint,
		creds.Password,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println(err.Error())
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
					'',
					'(format csv, DELIMITER '','', HEADER true)',
					'(%v,%v,us-east-1)');
			`, s3.Records[0].S3.Bucket.Name, s3.Records[0].S3.Object.Key)

		rows, err := db.Query(q)
		if err != nil {
			log.Println(err.Error())
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

			svc := lambdaservice.New(session.New())
input := &lambdaservice.InvokeInput{
    FunctionName:   aws.String("my-function"),
    InvocationType: aws.String("Event"),
    Payload:        []byte("{}"),
    Qualifier:      aws.String("$LATEST"),
}

result, err := svc.Invoke(input)
if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
        switch aerr.Code() {
        case lambdaservice.ErrCodeServiceException:
            fmt.Println(lambdaservice.ErrCodeServiceException, aerr.Error())
        case lambdaservice.ErrCodeResourceNotFoundException:
            fmt.Println(lambdaservice.ErrCodeResourceNotFoundException, aerr.Error())
        case lambdaservice.ErrCodeInvalidRequestContentException:
            fmt.Println(lambdaservice.ErrCodeInvalidRequestContentException, aerr.Error())
        case lambdaservice.ErrCodeRequestTooLargeException:
            fmt.Println(lambdaservice.ErrCodeRequestTooLargeException, aerr.Error())
        case lambdaservice.ErrCodeUnsupportedMediaTypeException:
            fmt.Println(lambdaservice.ErrCodeUnsupportedMediaTypeException, aerr.Error())
        case lambdaservice.ErrCodeTooManyRequestsException:
            fmt.Println(lambdaservice.ErrCodeTooManyRequestsException, aerr.Error())
        case lambdaservice.ErrCodeInvalidParameterValueException:
            fmt.Println(lambdaservice.ErrCodeInvalidParameterValueException, aerr.Error())
        case lambdaservice.ErrCodeEC2UnexpectedException:
            fmt.Println(lambdaservice.ErrCodeEC2UnexpectedException, aerr.Error())
        case lambdaservice.ErrCodeSubnetIPAddressLimitReachedException:
            fmt.Println(lambdaservice.ErrCodeSubnetIPAddressLimitReachedException, aerr.Error())
        case lambdaservice.ErrCodeENILimitReachedException:
            fmt.Println(lambdaservice.ErrCodeENILimitReachedException, aerr.Error())
        case lambdaservice.ErrCodeEFSMountConnectivityException:
            fmt.Println(lambdaservice.ErrCodeEFSMountConnectivityException, aerr.Error())
        case lambdaservice.ErrCodeEFSMountFailureException:
            fmt.Println(lambdaservice.ErrCodeEFSMountFailureException, aerr.Error())
        case lambdaservice.ErrCodeEFSMountTimeoutException:
            fmt.Println(lambdaservice.ErrCodeEFSMountTimeoutException, aerr.Error())
        case lambdaservice.ErrCodeEFSIOException:
            fmt.Println(lambdaservice.ErrCodeEFSIOException, aerr.Error())
        case lambdaservice.ErrCodeEC2ThrottledException:
            fmt.Println(lambdaservice.ErrCodeEC2ThrottledException, aerr.Error())
        case lambdaservice.ErrCodeEC2AccessDeniedException:
            fmt.Println(lambdaservice.ErrCodeEC2AccessDeniedException, aerr.Error())
        case lambdaservice.ErrCodeInvalidSubnetIDException:
            fmt.Println(lambdaservice.ErrCodeInvalidSubnetIDException, aerr.Error())
        case lambdaservice.ErrCodeInvalidSecurityGroupIDException:
            fmt.Println(lambdaservice.ErrCodeInvalidSecurityGroupIDException, aerr.Error())
        case lambdaservice.ErrCodeInvalidZipFileException:
            fmt.Println(lambdaservice.ErrCodeInvalidZipFileException, aerr.Error())
        case lambdaservice.ErrCodeKMSDisabledException:
            fmt.Println(lambdaservice.ErrCodeKMSDisabledException, aerr.Error())
        case lambdaservice.ErrCodeKMSInvalidStateException:
            fmt.Println(lambdaservice.ErrCodeKMSInvalidStateException, aerr.Error())
        case lambdaservice.ErrCodeKMSAccessDeniedException:
            fmt.Println(lambdaservice.ErrCodeKMSAccessDeniedException, aerr.Error())
        case lambdaservice.ErrCodeKMSNotFoundException:
            fmt.Println(lambdaservice.ErrCodeKMSNotFoundException, aerr.Error())
        case lambdaservice.ErrCodeInvalidRuntimeException:
            fmt.Println(lambdaservice.ErrCodeInvalidRuntimeException, aerr.Error())
        case lambdaservice.ErrCodeResourceConflictException:
            fmt.Println(lambdaservice.ErrCodeResourceConflictException, aerr.Error())
        case lambdaservice.ErrCodeResourceNotReadyException:
            fmt.Println(lambdaservice.ErrCodeResourceNotReadyException, aerr.Error())
        default:
            fmt.Println(aerr.Error())
        }
    } else {
        // Print the error, cast err to awserr.Error to get the Code and
        // Message from an error.
        fmt.Println(err.Error())
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
