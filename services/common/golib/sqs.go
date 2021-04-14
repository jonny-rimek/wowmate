package golib

import (
	"log"
	"strconv"
	"time"
)

// SQSEvent is all the data that gets passed into the lambda from the q
// IMPROVE: use events.SQSEvent, see summary lambda
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
		EventSource    string `json:"eventSource"`
		EventSourceARN string `json:"eventSourceARN"`
		AwsRegion      string `json:"awsRegion"`
	} `json:"Records"`
}

// TimeMessageInQueue
// IMPROVE: only pass in pointer to SQSEvent
// we can see an overall trend of this in the SQS metrics, I created this to double check how fast messages are polled
// with fargate, but since I use lambda the lambda service takes care of this and messages are usually only 10ms old
// before they are pushed into lambda
func TimeMessageInQueue(e SQSEvent, i int) error {
	j, err := strconv.ParseInt(e.Records[i].Attributes.ApproximateFirstReceiveTimestamp, 10, 64)
	if err != nil {
		log.Printf("Failed to parse int: %v", err)
		return err
	}
	tm1 := time.Unix(0, j*int64(1000000))

	ii, err := strconv.ParseInt(e.Records[i].Attributes.SentTimestamp, 10, 64)
	if err != nil {
		return err
	}
	tm2 := time.Unix(0, ii*int64(1000000))

	log.Printf("seconds the message was unprocessed in the queue: %v", tm1.Sub(tm2).Seconds())

	return nil
}
