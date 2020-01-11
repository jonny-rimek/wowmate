package golib

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"io"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sirupsen/logrus"
)

//DamageSummary is the DynamoDB schema for all damage summaries
type DamageSummary struct {
	BossFightUUID string `json:"pk"`
	CasterID      string `json:"sk"`
	EncounterID   int    `json:"gsi3pk"`
	Damage        int64  `json:"gsi123sk"`
	CasterName    string `json:"caster_name"`
}

//CanonicalLog writes a structured message to stdout if the log level is atleast INFO
func CanonicalLog(msg map[string]interface{}) {
	logrus.WithFields(msg).Info()
}

//DDBQuery runs a DynamoDB query and returns an APIGateway Response
//IMPROVE: the APIGateway stuff and the marshalling into a struct should be moved to a seperate function
func DDBQuery(ctx context.Context, queryInput *dynamodb.QueryInput) (float64, events.APIGatewayProxyResponse, error) {
	svc := dynamodb.New(session.New())
	var rcu float64
	result, err := svc.QueryWithContext(ctx, queryInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				logrus.Error(dynamodb.ErrCodeProvisionedThroughputExceededException)
				return rcu, APIGwError(429), nil
			case dynamodb.ErrCodeResourceNotFoundException:
				return rcu, APIGwError(500), err
			case dynamodb.ErrCodeInternalServerError:
				return rcu, APIGwError(500), err
			default:
				return rcu, APIGwError(500), err
			}
		} else {
			return rcu, APIGwError(500), err
		}
	}
	//WISHLIST: check result.LastEvaluatedKey if it is not nil return it
	//			to enable pagination

	rcu = *result.ConsumedCapacity.CapacityUnits
	if len(result.Items) == 0 {
		logrus.Error("no records returned from DynamoDB")
		return rcu, APIGwError(404), nil
	}

	summaries := []DamageSummary{}
	for _, item := range result.Items {
		s := DamageSummary{}
		err = dynamodbattribute.UnmarshalMap(item, &s)
		if err != nil {
			err = fmt.Errorf("Failed to unmarshal Record, %v", err)
			return rcu, APIGwError(500), err
		}
		summaries = append(summaries, s)
	}

	js, err := json.Marshal(summaries)
	if err != nil {
		return rcu, APIGwError(500), fmt.Errorf("failed to marshal data to JSON: %v", err)
	}

	return rcu, APIGwOK(js), nil
}

//APIGwOK return http status code 200 plus the body
func APIGwOK(body []byte) events.APIGatewayProxyResponse {
	h := make(map[string]string)
	h["Access-Control-Allow-Origin"] = "*"

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    h,
		Body:       string(body),
	}
}

//APIGwError is a helper function for APIGatewayProxyResponse
func APIGwError(status int) events.APIGatewayProxyResponse {
	h := make(map[string]string)
	h["Access-Control-Allow-Origin"] = "*"

	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers:    h,
		Body:       http.StatusText(status),
	}
}

//DownloadFileFromS3 writes a file from S3 to memory
func DownloadFileFromS3(bucket string, key string, sess *session.Session) ([]byte, int64, error) {
	downloader := s3manager.NewDownloader(sess)

	file := &aws.WriteAtBuffer{}

	bytes, err := downloader.Download(
		file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		return nil, bytes, fmt.Errorf("Unable to download item %v from bucket %v: %v", key, bucket, err)
	}
	return file.Bytes(), bytes, nil
}

//UploadFileToS3 writes files to S3 and cleans up the local storage
func UploadFileToS3(fileContent io.Reader, bucket string, key string, sess *session.Session) error {
	s3Svc := s3.New(sess)

	uploader := s3manager.NewUploaderWithClient(s3Svc)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    &key,
		Body:   fileContent,
	})
	if err != nil {
		logrus.Error("Failed to upload to S3: " + err.Error())
		return err
	}

	logrus.Debug("Upload finished! location: " + result.Location)

	return nil
}

//InitLogging sets up the logging for every lambda and should be called before the handler
func InitLogging() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "prod" {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.SetFormatter(&logrus.JSONFormatter{})
}
