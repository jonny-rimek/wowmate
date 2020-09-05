package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	uuid "github.com/gofrs/uuid"
)

//https://mholt.github.io/json-to-go/ best tool EVER

type Request struct {
	Records []struct {
		S3 struct {
			Bucket struct {
				Name string `json:"name"`
			} `json:"bucket"`
			Object struct {
				Key string `json:"key"`
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
	csvBucket := os.Getenv("CSV_BUCKET_NAME")
	if csvBucket == "" {
		return fmt.Errorf("csv bucket env var is empty")
	}

	sess, _ := session.NewSession()

	if len(e.Records) == 0 {
		return fmt.Errorf("SQS Event doesn't contain any messages")
	}
	log.Printf("amount of messages %v", len(e.Records))

	for j, msg := range e.Records {
		log.Printf("index j: %v", j+1)
		i, err := strconv.ParseInt(msg.Attributes.ApproximateFirstReceiveTimestamp, 10, 64)
		if err != nil {
			log.Printf("Failed to parse int: %v", err)
			return err
		}
		tm1 := time.Unix(0, i*int64(1000000))

		ii, err := strconv.ParseInt(msg.Attributes.SentTimestamp, 10, 64)
		if err != nil {
			return err
		}
		tm2 := time.Unix(0, ii*int64(1000000))

		log.Printf("seconds the message was unprocessed in the queue: %v", tm1.Sub(tm2).Seconds())

		body := msg.Body

		req := Request{}
		err = json.Unmarshal([]byte(body), &req)
		if err != nil {
			log.Printf("Failed Unmarshal: %v", err.Error())
			return err
		}

		downloader := s3manager.NewDownloader(sess)

		fileContent := &aws.WriteAtBuffer{}

		if len(req.Records) > 1 {
			log.Printf("Failed: the S3 event contains more than 1 element, not sure how that would happen")
			return err
		}

		_, err = downloader.Download(
			fileContent,
			&s3.GetObjectInput{
				Bucket: aws.String(req.Records[0].S3.Bucket.Name),
				Key:    aws.String(req.Records[0].S3.Object.Key),
			},
		)
		if err != nil {
			log.Printf("Failed to download item from bucket")
			return err
		}
		log.Println("downloaded from s3")

		s := bufio.NewScanner(bytes.NewReader(fileContent.Bytes()))
		uploadUUID := uuid.Must(uuid.NewV4()).String()

		events, err := Import(s, uploadUUID)
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		var ss [][]string
		w := csv.NewWriter(&buf)

		_ = EventsAsStringSlices(&events, &ss)
		if err := w.WriteAll(ss); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}

		r := io.Reader(&buf)

		log.Println("converted to csv")

		uploader := s3manager.NewUploader(sess)
		result, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(csvBucket),
			Key:    aws.String(fmt.Sprintf("%v.csv", uploadUUID)),
			Body:   r,
		})
		if err != nil {
			log.Println("Failed to upload to S3")
			return err
		}
		log.Println("Upload finished! location: " + result.Location)
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
