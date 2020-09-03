package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sqs"
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

func main() {
	queueURL := os.Getenv("QUEUE_URL")
	csvBucket := os.Getenv("CSV_BUCKET_NAME")

	sess, _ := session.NewSession()
	//TODO: check and handle error
	svc := sqs.New(sess)

	for {
		time.Sleep(time.Duration(2) * time.Second)

		msgResult, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
			AttributeNames: []*string{
				aws.String(sqs.MessageSystemAttributeNameSenderId),
				aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
				aws.String(sqs.MessageSystemAttributeNameApproximateReceiveCount),
				aws.String(sqs.MessageSystemAttributeNameApproximateFirstReceiveTimestamp),
				aws.String(sqs.MessageSystemAttributeNameSequenceNumber),
				aws.String(sqs.MessageSystemAttributeNameMessageDeduplicationId),
				aws.String(sqs.MessageSystemAttributeNameMessageGroupId),
				aws.String(sqs.MessageSystemAttributeNameAwstraceHeader),
			},
			MessageAttributeNames: []*string{
				aws.String(sqs.QueueAttributeNameAll),
			},
			QueueUrl:            aws.String(queueURL),
			MaxNumberOfMessages: aws.Int64(10),
			VisibilityTimeout:   aws.Int64(30), // 60 seconds
			WaitTimeSeconds:     aws.Int64(0),
		})

		if err != nil {
			log.Printf("recieve message failed: %v", err)
			return
		}

		if len(msgResult.Messages) == 0 {
			continue
		}
		//TODO: process all results
		// fmt.Printf("Success: %+v\n", msgResult.Messages)
		log.Printf("amount of messages %v", len(msgResult.Messages))

		i, err := strconv.ParseInt(*msgResult.Messages[0].Attributes["ApproximateFirstReceiveTimestamp"], 10, 64)
		if err != nil {
			log.Printf("failed to parse int: %v", err)
			return
		}
		tm1 := time.Unix(0, i*int64(1000000))

		ii, err := strconv.ParseInt(*msgResult.Messages[0].Attributes["SentTimestamp"], 10, 64)
		if err != nil {
			panic(err)
		}
		tm2 := time.Unix(0, ii*int64(1000000))

		log.Printf("seconds the message was unprocessed in the queue: %v", tm1.Sub(tm2).Seconds())

		body := *msgResult.Messages[0].Body

		req := Request{}
		err = json.Unmarshal([]byte(body), &req)
		if err != nil {
			log.Println("Unmarshal failed")
			return
		}

		downloader := s3manager.NewDownloader(sess)

		fileContent := &aws.WriteAtBuffer{}

		_, err = downloader.Download(
			fileContent,
			&s3.GetObjectInput{
				//TODO: check for more than 1 records
				Bucket: aws.String(req.Records[0].S3.Bucket.Name),
				Key:    aws.String(req.Records[0].S3.Object.Key),
			},
		)
		if err != nil {
			fmt.Printf("Unable to download item from bucket")
			// fmt.Printf("Unable to download item %v from bucket %v: %v", key, bucket, err)
			return
		}
		log.Println("downloaded from s3")

		uploader := s3manager.NewUploader(sess)
		result, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(csvBucket),
			Key:    aws.String("converted.csv"),
			Body:   bytes.NewReader(fileContent.Bytes()),
		})
		if err != nil {
			log.Println("Failed to upload to S3: " + err.Error())
			return
		}
		log.Println("Upload finished! location: " + result.Location)

		_, err = svc.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      aws.String(queueURL),
			ReceiptHandle: msgResult.Messages[0].ReceiptHandle,
		})
		if err != nil {
			log.Println("delete failed")
			continue
		}
		log.Println("message delete succeeded")
	}

}
