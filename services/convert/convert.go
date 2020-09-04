package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

/*
{
  "Records": [
    {
      "eventVersion": "2.1",
      "eventSource": "aws:s3",
      "awsRegion": "us-east-1",
      "eventTime": "2020-08-19T01:03:36.746Z",
      "eventName": "ObjectCreated:Put",
      "userIdentity": {
        "principalId": "AHRIC0SLDQ6UK"
      },
      "requestParameters": {
        "sourceIPAddress": "37.120.217.85"
      },
      "responseElements": {
        "x-amz-request-id": "2F5DA7D78A70E81D",
        "x-amz-id-2": "X5jUdNbjmiHrPnZj8rc4yOzwQoDR5nqw0H+15B7wm8kpjCxqCQouG3XQ94f3Fe1nM5vh3yBL5PCHdNcOq1UpFmNB5MH9x6ut"
      },
      "s3": {
        "s3SchemaVersion": "1.0",
        "configurationId": "OTI2NzljZTEtZmIxNi00N2I1LWFiNTMtNDNkOTY5MDc5MTIw",
        "bucket": {
          "name": "wm-converteruploadde59095e-akevvaglcv61",
          "ownerIdentity": {
            "principalId": "AHRIC0SLDQ6UK"
          },
          "arn": "arn:aws:s3:::wm-converteruploadde59095e-akevvaglcv61"
        },
        "object": {
          "key": "myFile",
          "size": 5,
          "eTag": "d8e8fca2dc0f896fd7cb4cb0031ba249",
          "sequencer": "005F3C7A6D4755F771"
        }
      }
    }
  ]
}
*/
//https://mholt.github.io/json-to-go/ best tool EVER

type Request struct {
	Records []struct {
		// EventVersion string    `json:"eventVersion"`
		// EventSource  string    `json:"eventSource"`
		// AwsRegion    string    `json:"awsRegion"`
		// EventTime    time.Time `json:"eventTime"`
		// EventName    string    `json:"eventName"`
		// UserIdentity struct {
		// 	PrincipalID string `json:"principalId"`
		// } `json:"userIdentity"`
		// RequestParameters struct {
		// 	SourceIPAddress string `json:"sourceIPAddress"`
		// } `json:"requestParameters"`
		// ResponseElements struct {
		// 	XAmzRequestID string `json:"x-amz-request-id"`
		// 	XAmzID2       string `json:"x-amz-id-2"`
		// } `json:"responseElements"`
		S3 struct {
			// S3SchemaVersion string `json:"s3SchemaVersion"`
			// ConfigurationID string `json:"configurationId"`
			Bucket struct {
				Name string `json:"name"`
				// OwnerIdentity struct {
				// 	PrincipalID string `json:"principalId"`
				// } `json:"ownerIdentity"`
				// Arn string `json:"arn"`
			} `json:"bucket"`
			Object struct {
				Key string `json:"key"`
				// Size      int    `json:"size"`
				// ETag      string `json:"eTag"`
				// Sequencer string `json:"sequencer"`
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
	// queueURL := os.Getenv("QUEUE_URL")
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

		//TODO: gzip before upload
		var buf bytes.Buffer
		w := gzip.NewWriter(&buf)

		_, err = w.Write(fileContent.Bytes())
		if err != nil {
			log.Println("Failed to create a gzip writer")
			return err
		}

		r, err := gzip.NewReader(&buf)
		if err != nil {
			log.Println("Failed to create a gzip reader")
			return err
		}

		uploader := s3manager.NewUploader(sess)
		result, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(csvBucket),
			Key:    aws.String("converted.csv.gz"),
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
