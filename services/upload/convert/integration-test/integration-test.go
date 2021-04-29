package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	lambdaService "github.com/aws/aws-sdk-go/service/lambda"
)

var lambdaSvc *lambdaService.Lambda

func invokeConvert(local bool) error {
	var functionName string
	var logType *string
	var payload []byte
	payload = nil

	if local == true {
		functionName = "ConvertLambda3540DCCB"
		logType = nil
		payload = []byte(`
			{
			  "Records": [
				{
				  "messageId": "19dd0b57-b21e-4ac1-bd88-01bbb068cb78",
				  "receiptHandle": "MessageReceiptHandle",
				  "body": "{\"Records\":[{\"eventVersion\":\"2.1\",\"eventSource\":\"aws:s3\",\"awsRegion\":\"us-east-1\",\"eventTime\":\"2020-09-03T22:04:40.595Z\",\"eventName\":\"ObjectCreated:Put\",\"userIdentity\":{\"principalId\":\"AWS:AIDA5WEGR6IGFG5BIU7QW\"},\"requestParameters\":{\"sourceIPAddress\":\"37.120.217.169\"},\"responseElements\":{\"x-amz-request-id\":\"2F5A942F7BB7532F\",\"x-amz-id-2\":\"fjLaREBoWXqzZymVY2hclfWmjutD45WbiKWoxCBeGkYnESkJp+yaSiAy5WcTX9yaup1xERIJbXPUKLmQzg2tZHAS5v8cdcXI\"},\"s3\":{\"s3SchemaVersion\":\"1.0\",\"configurationId\":\"NGJhMDAzOWEtNTliZS00NTY4LThhZjgtNjAyODc1NmJiMjQy\",\"bucket\":{\"name\":\"wm-dev-bucketsupload0b7f8f15-f6hlcjt6x88p\",\"ownerIdentity\":{\"principalId\":\"AHRIC0SLDQ6UK\"},\"arn\":\"arn:aws:s3:::wm-bucketsupload0b7f8f15-1k8l1923idej5\"},\"object\":{\"key\":\"upload/2021/4/14/8/eb148ae5-52a7-4be0-b026-dfab25da00d9.zip\",\"size\":6873965,\"eTag\":\"cd1a984dd8323ef1e6e5541d363d5e4d\",\"sequencer\":\"005F5168797266420C\"}}}]}",
				  "attributes": {
					"ApproximateReceiveCount": "1",
					"SentTimestamp": "1523232000000",
					"SenderId": "123456789012",
					"ApproximateFirstReceiveTimestamp": "1523232000001"
				  },
				  "eventSource": "aws:sqs",
				  "eventSourceARN": "arn:aws:sqs:us-east-1:123456789012:MyQueue",
				  "awsRegion": "us-east-1"
				}
			  ]
			}
		`)
	} else {
		payload = []byte(`
			{
			  "Records": [
				{
				  "messageId": "19dd0b57-b21e-4ac1-bd88-01bbb068cb78",
				  "receiptHandle": "MessageReceiptHandle",
				  "body": "{\"Records\":[{\"eventVersion\":\"2.1\",\"eventSource\":\"aws:s3\",\"awsRegion\":\"us-east-1\",\"eventTime\":\"2020-09-03T22:04:40.595Z\",\"eventName\":\"ObjectCreated:Put\",\"userIdentity\":{\"principalId\":\"AWS:AIDA5WEGR6IGFG5BIU7QW\"},\"requestParameters\":{\"sourceIPAddress\":\"37.120.217.169\"},\"responseElements\":{\"x-amz-request-id\":\"2F5A942F7BB7532F\",\"x-amz-id-2\":\"fjLaREBoWXqzZymVY2hclfWmjutD45WbiKWoxCBeGkYnESkJp+yaSiAy5WcTX9yaup1xERIJbXPUKLmQzg2tZHAS5v8cdcXI\"},\"s3\":{\"s3SchemaVersion\":\"1.0\",\"configurationId\":\"NGJhMDAzOWEtNTliZS00NTY4LThhZjgtNjAyODc1NmJiMjQy\",\"bucket\":{\"name\":\"wm-preprod-bucketsupload0b7f8f15-14ny22audh3ba\",\"ownerIdentity\":{\"principalId\":\"AHRIC0SLDQ6UK\"},\"arn\":\"arn:aws:s3:::wm-bucketsupload0b7f8f15-1k8l1923idej5\"},\"object\":{\"key\":\"upload/2021/4/26/18/34d55f0c-fadf-4234-9c99-87a855f20ef0.zip\",\"size\":6873965,\"eTag\":\"cd1a984dd8323ef1e6e5541d363d5e4d\",\"sequencer\":\"005F5168797266420C\"}}}]}",
				  "attributes": {
					"ApproximateReceiveCount": "1",
					"SentTimestamp": "1523232000000",
					"SenderId": "123456789012",
					"ApproximateFirstReceiveTimestamp": "1523232000001"
				  },
				  "eventSource": "aws:sqs",
				  "eventSourceARN": "arn:aws:sqs:us-east-1:123456789012:MyQueue",
				  "awsRegion": "us-east-1"
				}
			  ]
			}
		`)
		functionName = "wm-preprod-ConvertLambda3540DCCB-1898ZUCBEU62K"
		logType = aws.String(lambdaService.LogTypeTail)
	}

	// TODO: add payload s3 event
	input := &lambdaService.InvokeInput{
		FunctionName:   &functionName,
		Payload:        payload,
		InvocationType: aws.String(lambdaService.InvocationTypeRequestResponse), // synchronous - default
		LogType:        logType,                                                 // returns the log in the response
	}
	err := input.Validate()
	if err != nil {
		return fmt.Errorf("validating the input failed: %v", err)
	}

	resp, err := lambdaSvc.Invoke(input)
	if err != nil {
		return fmt.Errorf("failed to invoke lambda: %v", err)
	}
	log.Printf("%s", resp.Payload)

	if local != true {
		// log doesn't work locally
		decodeString, err := base64.StdEncoding.DecodeString(*resp.LogResult)
		if err != nil {
			return fmt.Errorf("failed to decode the log: %v", err)
		}
		log.Printf("log result: %s", decodeString)
	}

	if resp.FunctionError != nil {
		return fmt.Errorf("lambda was invoked but returned error: %s", *resp.FunctionError)
	}
	return nil
}

func main() {
	local := true
	// CI is always set to true if it runs in a github action, if the env is not present it's ran locally
	if os.Getenv("CI") == "true" {
		local = false
	}
	var sess *session.Session
	var err error

	if local == true {
		sess, err = session.NewSession(
			&aws.Config{
				Endpoint:   aws.String("http://127.0.0.1:3001"),
				DisableSSL: aws.Bool(true),
				Region:     aws.String("us-east-1"),
				MaxRetries: aws.Int(0),
			})
	} else {
		sess, err = session.NewSession()
	}
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	lambdaSvc = lambdaService.New(sess)

	err = invokeConvert(local)
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
	log.Println("invoke finished")
}
