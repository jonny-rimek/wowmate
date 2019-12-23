package main

import (
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
	CasterID      string `json:"caster_id"`
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
		//TODO add GSI
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
	var wcuConsumed float64

	for _, s := range records {
		av, err := dynamodbattribute.MarshalMap(s)
		if err != nil {
			log.Println("Got error marshalling map:")
			log.Fatalf(err.Error())
		}

		input := &dynamodb.PutItemInput{
			Item:                   av,
			ReturnConsumedCapacity: aws.String("TOTAL"),
			TableName:              aws.String(ddbTableName),
		}

		oup, err := svcdb.PutItem(input)
		if err != nil {
			log.Fatalf("Got error calling PutItem: %v", err)
		}

		wcuConsumed = wcuConsumed + *oup.ConsumedCapacity.CapacityUnits
	}
	log.Printf("Consumed WCU: %f: ", wcuConsumed)

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
