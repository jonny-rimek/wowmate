package main

import (
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
	queueURL := os.Getenv("QUEUE_URL")
	log.Println(queueURL)

	for {
		log.Println("lol")
		time.Sleep(time.Duration(10)*time.Second)

		sess,_ := session.NewSession()
		//TODO: check and handle error
		svc := sqs.New(sess)
		
		msgResult, _ := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl: aws.String(queueURL),
		})

		if len(msgResult.Messages) > 0 {
			log.Println(*msgResult.Messages[0].Body)
			//TODO: download from s3

			_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl: aws.String(queueURL),
				ReceiptHandle: msgResult.Messages[0].ReceiptHandle,
			})

			if err != nil {
				log.Println("delete failed")
				continue;
			}
			log.Println("delete succeeded, unless failed")
		}
	}
}

