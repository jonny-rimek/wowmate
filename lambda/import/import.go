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

//Event ..
type Event struct {
	BucketName string `json:"result_bucket"`
	Key        string `json:"file_name"`
}

//CSV ..
type CSV struct {
	BossFightUUID string `json:"pk"`
	Damage        int64  `json:"sk"`
	CasterName    string `json:"caster_name"`
	CasterID      string `json:"gsi2pk"`
	EncounterID   int    `json:"gsi1pk"`
}

func handler(e Event) error {
	ddbTableName := os.Getenv("DDB_NAME")
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)

	downloader := s3manager.NewDownloader(sess)

	file := &aws.WriteAtBuffer{}

	numBytes, err := downloader.Download(
		file,
		&s3.GetObjectInput{
			Bucket: aws.String(e.BucketName),
			Key:    aws.String(e.Key),
		})

	if err != nil {
		log.Fatalf("Unable to download item %q: %v", e.Key, err)
	}

	log.Println("DEBUG: Downloaded", numBytes, "bytes")

	var records []CSV

	reader := bytes.NewReader(file.Bytes())
	scanner := bufio.NewScanner(reader)

	scanner.Scan() //skips the first line, which is the header of the csv
	for scanner.Scan() {
		row := strings.Split(scanner.Text(), ",")

		damage, err := strconv.ParseInt(trimQuotes(row[1]), 10, 64)
		if err != nil {
			log.Fatalf("Failed to convert damage column to int64")
		}

		encounterID, err := strconv.Atoi(trimQuotes(row[4]))
		if err != nil {
			log.Fatalf("Failed to convert damage column to int64")
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

	svcdb := dynamodb.New(sess)

	writesRequets := []*dynamodb.WriteRequest{}

	for _, s := range records {
		av, err := dynamodbattribute.MarshalMap(s)
		if err != nil {
			log.Println("Got error marshalling map:")
			log.Fatalf(err.Error())
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

func trimQuotes(input string) (output string) {
	output = strings.TrimSuffix(input, "\"")
	output = strings.TrimPrefix(output, "\"")
	return output
}

func main() {
	lambda.Start(handler)
}
