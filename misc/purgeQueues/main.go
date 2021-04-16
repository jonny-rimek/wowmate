package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		log.Printf("Error creating session: %v", err.Error())
	}

	svc := sqs.New(sess)
	queues, err := svc.ListQueues(&sqs.ListQueuesInput{
		MaxResults:      aws.Int64(100),
	})
	if err != nil {
		log.Printf("failed listing queues %v", err)
		return
	}
	for _, q := range queues.QueueUrls {
		_, err := svc.PurgeQueue(&sqs.PurgeQueueInput{
			QueueUrl: q,
		})
		if err != nil {
			log.Printf("failed purging queue %p: %v", q, err)
			return
		}
		log.Printf("purged %s", *q)
	}
}
