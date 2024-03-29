package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type messageSender func(*sqs.SQS, *string, *string) error

func sendMessage(ms messageSender, svc *sqs.SQS, messageBody *string, queueURL *string, amount int) error {
	errorChannel := make(chan error)

	for i := 0; i < amount; i++ {
		go func() {
			errorChannel <- ms(svc, messageBody, queueURL)
		}()
	}

	var errs []error

	for i := 0; i < amount; i++ {
		err := <-errorChannel
		errs = append(errs, err)
	}

	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	log.Println("done")
	return nil
}

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		log.Printf("Error creating session: %v", err.Error())
	}

	svc := sqs.New(sess)

	arg := os.Args

	data, err := ioutil.ReadFile(arg[1])
	if err != nil {
		log.Println("File reading error", err)
		return
	}

	s3Event := string(data)
	queueURL := "https://sqs.us-east-1.amazonaws.com/461497339039/wm-dev-ConvertProcessingQueueE8D6E023-17QQELE95GJC8"

	// one batch contains 10 messages, to send 7,5k events set amount to 750
	var amount int

	// if no argument is provided to the cli send 1 batch of messages aka 10
	if len(arg) == 2 {
		amount = 1
	} else {
		amount, err = strconv.Atoi(arg[2])
		if err != nil {
			log.Println("failed to convert cli argument to int")
			return
		}
	}


	err = sendMessage(sqsBatchMessageSender, svc, &s3Event, &queueURL, amount)
	if err != nil {
		log.Println(err)
		return
	}
}

func sqsBatchMessageSender(svc *sqs.SQS, messageBody *string, queueURL *string) error {
	var entries []*sqs.SendMessageBatchRequestEntry
	for i := 0; i < 10; i++ {
		entry := sqs.SendMessageBatchRequestEntry{
			Id:          aws.String(strconv.Itoa(i)),
			MessageBody: messageBody,
		}
		entries = append(entries, &entry)
	}

	message, err := svc.SendMessageBatch(&sqs.SendMessageBatchInput{
		Entries:  entries,
		QueueUrl: queueURL,
	})
	if err != nil {
		log.Printf("failed to send message to sqs: %v", err)
		return nil
	}
	log.Println(message)

	return nil
}
