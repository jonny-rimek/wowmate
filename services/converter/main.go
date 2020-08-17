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

		log.Println(msgResult.Messages[0].Body)
	}
}

