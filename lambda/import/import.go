package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
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

//CSV2 test ..
type CSV2 struct {
	BossFightUUID string `json:"boss_fight_uuid"`
	Damage        string `json:"damage"`
	CasterName    string `json:"caster_name"`
	CasterID      string `json:"caster_id"`
	EncounterID   string `json:"encounter_id"`
}

func handler(e Event) error {
	ddbTableName := os.Getenv("DDB_NAME")
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)

	downloader := s3manager.NewDownloader(sess)

	file2 := &aws.WriteAtBuffer{}

	numBytes2, err := downloader.Download(
		file2,
		&s3.GetObjectInput{
			Bucket: aws.String(e.BucketName),
			Key:    aws.String(e.Key),
		})

	if err != nil {
		log.Fatalf("Unable to download item %q: %v", e.Key, err)
	}

	log.Println("DEBUG: Downloaded", numBytes2, "bytes")

	var records2 []CSV2

	r2 := bytes.NewReader(file2.Bytes())
	s2 := bufio.NewScanner(r2)

	for s2.Scan() {
		//TODO add GSI
		//TODO skip first line
		//TODO convert to int
		row2 := strings.Split(s2.Text(), ",")

		if err != nil {
			log.Fatalf("Failed to convert 2nd row to int32")
		}

		r2 := CSV2{
			row2[0],
			row2[1],
			row2[2],
			row2[3],
			row2[4],
		}

		records2 = append(records2, r2)
	}

	log.Print("DEBUG: read CSV into structs")

	svcdb := dynamodb.New(sess)
	var wcuConsumed float64

	for _, s := range records2 {
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

func main() {
	lambda.Start(handler)
}
