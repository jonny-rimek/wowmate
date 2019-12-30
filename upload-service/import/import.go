package main

import (
	"bufio"
	"bytes"
	"log"
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
)

//Event is the data from StepFunctions
type Event struct {
	BucketName string `json:"result_bucket"`
	Key        string `json:"file_name"`
}

//CSV is the query output from Athena
type CSV struct {
	CasterID      string `json:"pk"`
	Damage        int64  `json:"sk"`
	CasterName    string `json:"caster_name"`
	BossFightUUID string `json:"gsi2pk"`
	EncounterID   int    `json:"gsi1pk"`
}

func handler(e Event) error {
	sess, _ := session.NewSession()
	//the file from s3 is read directly into memory
	file, err := downloadFileFromS3(e.BucketName, e.Key, sess)
	if err != nil {
		return err
	}

	records, err := parseCSV(file)
	if err != nil {
		return err
	}
	writeRequests, err := createDynamoDBWriteRequest(records)
	var writes []*dynamodb.WriteRequest

	//TODO: extra function and sum WCU consumed
	log.Println("DEBUG: starting to loop through write request array")
	for _, value := range writeRequests {
		writes = append(writes, value)
		if len(writes) == 25 {
			log.Println("DEBUG: writing batch to dynamodb")
			err = writeBatchDynamoDB(writes, sess)
			if err != nil {
				log.Println(err)
				return err
			}
			writes = nil
		}
	}
	err = writeBatchDynamoDB(writes, sess)
	if err != nil {
		return err
	}
//TODO: add logrus and log level, read log level from env

	return nil
}

func createDynamoDBWriteRequest(records []CSV) ([]*dynamodb.WriteRequest, error) {
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

func writeBatchDynamoDB(writeRequests[]*dynamodb.WriteRequest, sess *session.Session) error {
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
				return fmt.Errorf("%v -- %v", dynamodb.ErrCodeProvisionedThroughputExceededException, err)
			case dynamodb.ErrCodeResourceNotFoundException:
				return fmt.Errorf("%v -- %v", dynamodb.ErrCodeResourceNotFoundException, err)
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				return fmt.Errorf("%v -- %v", dynamodb.ErrCodeItemCollectionSizeLimitExceededException, err)
			case dynamodb.ErrCodeRequestLimitExceeded:
				return fmt.Errorf("%v -- %v", dynamodb.ErrCodeRequestLimitExceeded, err)
			case dynamodb.ErrCodeInternalServerError:
				return fmt.Errorf("%v -- %v", dynamodb.ErrCodeInternalServerError, err)
			case dynamodb.ErrCodeTransactionCanceledException:
				return err
			default:
				return fmt.Errorf("default error: %v", err)
			}
		} else {
			return fmt.Errorf("non aws error: %v", err)
		}
	}
	log.Printf("Consumed WCU: %v", *result.ConsumedCapacity[0].CapacityUnits)
	//TODO: check unprocessed items of resuilt
	return nil
}

func parseCSV(file []byte) ([]CSV, error){
	var records []CSV

	reader := bytes.NewReader(file)
	scanner := bufio.NewScanner(reader)

	scanner.Scan() //skips the first line, which is the header of the csv
	for scanner.Scan() {
		row := strings.Split(scanner.Text(), ",")

		damage, err := strconv.ParseInt(trimQuotes(row[1]), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to convert damage column to int64: %v", err)
		}

		encounterID, err := strconv.Atoi(trimQuotes(row[4]))
		if err != nil {
			return nil, fmt.Errorf("Failed to convert encounter id column to int: %v", err)
		}

		r := CSV{
			trimQuotes(row[0]),
			damage,
			trimQuotes(row[2]),
			trimQuotes(row[3]),
			encounterID,
		}

		records = append(records, r)
	}

	log.Println("DEBUG: read CSV into structs")

	return records, nil
}

func downloadFileFromS3(bucket string, key string, sess *session.Session) ([]byte, error) {
	downloader := s3manager.NewDownloader(sess)

	file := &aws.WriteAtBuffer{}

	numBytes, err := downloader.Download(
		file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		return nil, fmt.Errorf("Unable to download item %v from bucket %v: %v", key, bucket, err)
	}

	log.Printf("DEBUG: Downloaded %v bytes %v/%v", numBytes, bucket, key)

	return file.Bytes(), nil
}

func trimQuotes(input string) (output string) {
	output = strings.TrimSuffix(input, "\"")
	output = strings.TrimPrefix(output, "\"")
	return output
}

func main() {
	lambda.Start(handler)
}
