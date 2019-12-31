package main

import (
	"github.com/sirupsen/logrus"
	"bufio"
	"bytes"
	"os"
	"strconv"
	"strings"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jonny-rimek/wowmate/services/golib"
)

//Event is the data from StepFunctions
type Event struct {
	BucketName string `json:"result_bucket"`
	Key        string `json:"file_name"`
}

func handler(e Event) error {
	bytes, wcu, err := handle(e)
	writeCanonicalLog(e.BucketName, e.Key, bytes, wcu)
	return err
}

func handle(e Event) (int64, float64, error) {
	sess, _ := session.NewSession()
	var bytes int64
	var wcu float64

	file, bytes, err := downloadFileFromS3(e.BucketName, e.Key, sess)
	if err != nil {
		return bytes, 0, err
	}

	records, err := parseCSV(file)
	if err != nil {
		return bytes, 0, err
	}

	wcu, err = writeDynamoDB(records, sess)
	return bytes, wcu, err
}

func writeCanonicalLog(bucketName string, objectKey string, bytes int64, wcu float64){
	logrus.WithFields(logrus.Fields{
		"bucket":        bucketName, 
		"key":           objectKey,
		"downloaded KB": bytes/1024,
		"wcu":           wcu,
	}).Info()
}

func writeDynamoDB(records []golib.DamageSummary, sess *session.Session) (float64, error) {
	writeRequests, err := createDynamoDBWriteRequest(records)
	var writes []*dynamodb.WriteRequest

	var consumedWCU float64
	for _, value := range writeRequests {
		writes = append(writes, value)
		if len(writes) == 25 {
			logrus.Debug("writing batch to dynamodb")
			wcu, err := writeBatchDynamoDB(writes, sess)
			if err != nil {
				return consumedWCU, err
			}
			consumedWCU += wcu
			writes = nil
		}
	}
	//NOTE: if the size was exactly 25 this will still execute with 
	//an empty array, not sure how it will behave
	wcu, err := writeBatchDynamoDB(writes, sess)
	if err != nil {
		return consumedWCU, err
	}
	consumedWCU += wcu

	return consumedWCU, nil
}

func createDynamoDBWriteRequest(records []golib.DamageSummary) ([]*dynamodb.WriteRequest, error) {
	writesRequets := []*dynamodb.WriteRequest{}

	for _, s := range records {
		av, err := dynamodbattribute.MarshalMap(s)
		if err != nil {
			return nil, fmt.Errorf("got error marshalling csv struct into dynamoDB element: %v", err)
		}

		wr := &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: av,
			},
		}
		writesRequets = append(writesRequets, wr)
	}
	return writesRequets, nil
}

func writeBatchDynamoDB(writeRequests[]*dynamodb.WriteRequest, sess *session.Session) (float64, error) {
	svcdb := dynamodb.New(sess)
	ddbTableName := os.Getenv("DDB_NAME")

	input:= &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			ddbTableName: writeRequests,
		},
		ReturnConsumedCapacity: aws.String("TOTAL"),
	}

	result, err := svcdb.BatchWriteItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				return 0, fmt.Errorf("%v -- %v", dynamodb.ErrCodeProvisionedThroughputExceededException, err)
			case dynamodb.ErrCodeResourceNotFoundException:
				return 0, fmt.Errorf("%v -- %v", dynamodb.ErrCodeResourceNotFoundException, err)
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				return 0, fmt.Errorf("%v -- %v", dynamodb.ErrCodeItemCollectionSizeLimitExceededException, err)
			case dynamodb.ErrCodeRequestLimitExceeded:
				return 0, fmt.Errorf("%v -- %v", dynamodb.ErrCodeRequestLimitExceeded, err)
			case dynamodb.ErrCodeInternalServerError:
				return 0, fmt.Errorf("%v -- %v", dynamodb.ErrCodeInternalServerError, err)
			case dynamodb.ErrCodeTransactionCanceledException:
				return 0, err
			default:
				return 0, fmt.Errorf("default error: %v", err)
			}
		} else {
			return 0, fmt.Errorf("non aws error: %v", err)
		}
	}
	//NOTE: unprocessed items of result are never check, if it is not empty the lambda will
	//		fail and thus alert me, when the case arrises
	//		when does this occur, if I get an error I believe non in the batch got written to DDB
	if len(result.UnprocessedItems) > 0 {
		return 0, fmt.Errorf("handle unprocessed items")
	}
	return *result.ConsumedCapacity[0].CapacityUnits, nil
}

func parseCSV(file []byte) ([]golib.DamageSummary, error){
	var records []golib.DamageSummary

	reader := bytes.NewReader(file)
	scanner := bufio.NewScanner(reader)

	scanner.Scan() //skips the first line, which is the header of the csv
	for scanner.Scan() {
		row := strings.Split(scanner.Text(), ",")

		damage, err := strconv.ParseInt(trimQuotes(row[0]), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to convert damage column to int64: %v", err)
		}

		encounterID, err := strconv.Atoi(trimQuotes(row[1]))
		if err != nil {
			return nil, fmt.Errorf("Failed to convert encounter id column to int: %v", err)
		}

		r := golib.DamageSummary{
			PK:            fmt.Sprintf("%v#%v",trimQuotes(row[0]), trimQuotes(row[3])),
			Damage:        damage,
			EncounterID:   encounterID,
			BossFightUUID: trimQuotes(row[2]), //boss fight uuid
			CasterID:      trimQuotes(row[3]), //caster id
			CasterName:    trimQuotes(row[4]), //caster name
		}

		records = append(records, r)
	}

	logrus.Debug("read CSV into structs")

	return records, nil
}

func downloadFileFromS3(bucket string, key string, sess *session.Session) ([]byte, int64, error) {
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

func trimQuotes(input string) string {
	output := strings.TrimSuffix(input, "\"")
	output = strings.TrimPrefix(output, "\"")
	return output
}

func main() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "prod" {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.SetFormatter(&logrus.JSONFormatter{})

	lambda.Start(handler)
}
