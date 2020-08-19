package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
	queueURL := os.Getenv("QUEUE_URL")
	bucket := os.Getenv("BUCKET_NAME")
	key := os.Getenv("FILE_NAME")
	log.Println(queueURL)

	for {
		log.Println("lol")
		time.Sleep(time.Duration(10) * time.Second)

		sess, _ := session.NewSession()
		//TODO: check and handle error
		svc := sqs.New(sess)

		msgResult, _ := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl: aws.String(queueURL),
		})

		if len(msgResult.Messages) > 0 {
			log.Println(*msgResult.Messages[0].Body)
			downloader := s3manager.NewDownloader(sess)

			file := &aws.WriteAtBuffer{}

			//bytes unused
			_, err := downloader.Download(
				file,
				&s3.GetObjectInput{
					Bucket: aws.String(bucket),
					Key:    aws.String(key),
				},
			)

			if err != nil {
				fmt.Printf("Unable to download item %v from bucket %v: %v", key, bucket, err)
				return
			}

			// _, err = svc.DeleteMessage(&sqs.DeleteMessageInput{
			// 	QueueUrl:      aws.String(queueURL),
			// 	ReceiptHandle: msgResult.Messages[0].ReceiptHandle,
			// })

			if err != nil {
				log.Println("delete failed")
				continue
			}
			log.Println("delete succeeded, unless failed")
		}
	}
}
