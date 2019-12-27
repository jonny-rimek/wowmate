package main

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"bufio"
	"bytes"
	"log"
	"os"
	"strconv"
	"strings"

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
	BossFightUUID string `json:"pk"`
	Damage        int64  `json:"sk"`
	CasterName    string `json:"caster_name"`
	CasterID      string `json:"gsi2pk"`
	EncounterID   int    `json:"gsi1pk"`
}

func handler(e Event) error {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)
	//the file from s3 is read directly into memory
	file, err := downloadFileFromS3(e.BucketName, e.Key, sess)
	if err != nil {
		return err
	}

	records, err := parseCSV(file)
	if err != nil {
		return err
	}

	err = writeBatchDynamoDB(records, sess)
//TODO: add logrus and log level, read log level from env

	return nil
}

func writeBatchDynamoDB(records []CSV, sess *session.Session) error {
	svcdb := dynamodb.New(sess)
	ddbTableName := os.Getenv("DDB_NAME")

	writesRequets := []*dynamodb.WriteRequest{}

	for _, s := range records {
		av, err := dynamodbattribute.MarshalMap(s)
		if err != nil {
			log.Println("Got error marshalling map:")
			return err
		}

		wr := &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: av,
			},
		}

		writesRequets = append(writesRequets, wr)

	}
	input:= &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			ddbTableName: writesRequets,
		},
		ReturnConsumedCapacity: aws.String("TOTAL"),
		
	}

	result, err := svcdb.BatchWriteItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				log.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				log.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				log.Println(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				log.Println(dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				log.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}
		return err
	}
	log.Printf("Consumed WCU: %v: ", *result.ConsumedCapacity[0].CapacityUnits)
	//TODO: check unprocessed items of resuilt
	//TODO: check size of input, if > 25 loop
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
			log.Println("Failed to convert damage column to int64")
			return nil, err
		}

		encounterID, err := strconv.Atoi(trimQuotes(row[4]))
		if err != nil {
			log.Println("Failed to convert damage column to int64")
			return nil, err
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

	log.Print("DEBUG: read CSV into structs")

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
		log.Printf("Unable to download item %v from bucket %v: %v", key, bucket, err)
		return nil, err
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
