package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

var lambdaSvc *lambda.Lambda

func invokeConvert(local bool) ([]string, error) {
	log.Println("invoked convert lambda")

	var functionName, bucketName, objectKey string
	var logType *string

	if local == true {
		functionName = "ConvertLambda3540DCCB"
		logType = nil
		bucketName = "wm-dev-bucketsupload0b7f8f15-f6hlcjt6x88p"
		objectKey = "upload/2021/4/14/8/eb148ae5-52a7-4be0-b026-dfab25da00d9.zip"
	} else {
		bucketName = "wm-preprod-bucketsupload0b7f8f15-14ny22audh3ba"
		objectKey = "upload/2021/4/26/18/34d55f0c-fadf-4234-9c99-87a855f20ef0.zip"
		functionName = "wm-preprod-ConvertLambda3540DCCB-1898ZUCBEU62K"
		logType = aws.String(lambda.LogTypeTail)
	}
	payload := []byte(fmt.Sprintf(`
			{
			  "Records": [
				{
				  "messageId": "19dd0b57-b21e-4ac1-bd88-01bbb068cb78",
				  "receiptHandle": "MessageReceiptHandle",
				  "body": "{\"Records\":[{\"eventVersion\":\"2.1\",\"eventSource\":\"aws:s3\",\"awsRegion\":\"us-east-1\",\"eventTime\":\"2020-09-03T22:04:40.595Z\",\"eventName\":\"ObjectCreated:Put\",\"userIdentity\":{\"principalId\":\"AWS:AIDA5WEGR6IGFG5BIU7QW\"},\"requestParameters\":{\"sourceIPAddress\":\"37.120.217.169\"},\"responseElements\":{\"x-amz-request-id\":\"2F5A942F7BB7532F\",\"x-amz-id-2\":\"fjLaREBoWXqzZymVY2hclfWmjutD45WbiKWoxCBeGkYnESkJp+yaSiAy5WcTX9yaup1xERIJbXPUKLmQzg2tZHAS5v8cdcXI\"},\"s3\":{\"s3SchemaVersion\":\"1.0\",\"configurationId\":\"NGJhMDAzOWEtNTliZS00NTY4LThhZjgtNjAyODc1NmJiMjQy\",\"bucket\":{\"name\":\"%s\",\"ownerIdentity\":{\"principalId\":\"AHRIC0SLDQ6UK\"},\"arn\":\"arn:aws:s3:::wm-bucketsupload0b7f8f15-1k8l1923idej5\"},\"object\":{\"key\":\"%s\",\"size\":6873965,\"eTag\":\"cd1a984dd8323ef1e6e5541d363d5e4d\",\"sequencer\":\"005F5168797266420C\"}}}]}",
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
		`, bucketName, objectKey))

	input := &lambda.InvokeInput{
		FunctionName:   &functionName,
		Payload:        payload,
		InvocationType: aws.String(lambda.InvocationTypeRequestResponse), // synchronous - default
		LogType:        logType,                                          // returns the log in the response
	}
	err := input.Validate()
	if err != nil {
		return nil, fmt.Errorf("validating the input failed: %v", err)
	}

	resp, err := lambdaSvc.Invoke(input)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke lambda: %v", err)
	}

	if local != true {
		// log doesn't work locally
		decodeString, err := base64.StdEncoding.DecodeString(*resp.LogResult)
		if err != nil {
			return nil, fmt.Errorf("failed to decode the log: %v", err)
		}
		log.Printf("log result: %s", decodeString)
	}

	if resp.FunctionError != nil {
		return nil, fmt.Errorf("lambda was invoked but returned error: %s", *resp.FunctionError)
	}

	// parse output and convert to string array
	var r []string
	err = json.Unmarshal(resp.Payload, &r)
	if err != nil {
		return nil, err
	}

	log.Println("convert finished")
	return r, nil
}

func invokeQueryPlayerDamageDone(combatlogUUID string, local bool) error {
	log.Println("invoke query player damage")
	var functionName string
	var logType *string

	if local == true {
		functionName = "QueryPlayerDamageDoneLambda98AFC037"
		logType = nil
	} else {
		functionName = "wm-preprod-QueryPlayerDamageDoneLambda98AFC037-1UAP6TNCDJR2U"
		logType = aws.String(lambda.LogTypeTail)
	}
	payload := []byte(fmt.Sprintf(`
		{
		  "Records": [
			{
			  "Sns": {
				"Subject": "unused",
				"Message": "%s",
				"MessageAttributes": {
				  "Test": {
					"Type": "unused",
					"Value": "unused"
				  }
				}
			  }
			}
		  ]
		}
		`, combatlogUUID))

	input := &lambda.InvokeInput{
		FunctionName:   &functionName,
		Payload:        payload,
		InvocationType: aws.String(lambda.InvocationTypeRequestResponse), // synchronous - default
		// it needs to be synchronous so I can assert on the response if an error was returned or not
		// even if in aws it's called async via SNS
		LogType: logType, // returns the log in the response
	}
	err := input.Validate()
	if err != nil {
		return fmt.Errorf("validating the input failed: %s", err.Error())
	}

	resp, err := lambdaSvc.Invoke(input)
	if err != nil {
		return fmt.Errorf("failed to invoke lambda: %s", err.Error())
	}
	log.Printf("%s", resp.Payload)

	if local != true {
		// log doesn't work locally
		decodeString, err := base64.StdEncoding.DecodeString(*resp.LogResult)
		if err != nil {
			return fmt.Errorf("failed to decode the log: %s", err.Error())
		}
		log.Printf("log result: %s", decodeString)
	}

	if resp.FunctionError != nil {
		return fmt.Errorf("lambda was invoked but returned error: %s", *resp.FunctionError)
	}
	log.Println("query player damage finished")

	return nil
}
func invokeInsertPlayerDamageDone(local bool) error {
	log.Println("invoke insert player damage")
	var functionName string
	var logType *string

	if local == true {
		functionName = ""
		logType = nil
	} else {
		functionName = "" // TODO
		logType = aws.String(lambda.LogTypeTail)
	}
	payload := []byte(fmt.Sprintf(`
		`)) // TODO

	input := &lambda.InvokeInput{
		FunctionName:   &functionName,
		Payload:        payload,
		InvocationType: aws.String(lambda.InvocationTypeRequestResponse), // synchronous - default
		// it needs to be synchronous so I can assert on the response if an error was returned or not
		// even if in aws it's called async via SNS
		LogType: logType, // returns the log in the response
	}
	err := input.Validate()
	if err != nil {
		return fmt.Errorf("validating the input failed: %s", err.Error())
	}

	resp, err := lambdaSvc.Invoke(input)
	if err != nil {
		return fmt.Errorf("failed to invoke lambda: %s", err.Error())
	}
	log.Printf("%s", resp.Payload)

	if local != true {
		// log doesn't work locally
		decodeString, err := base64.StdEncoding.DecodeString(*resp.LogResult)
		if err != nil {
			return fmt.Errorf("failed to decode the log: %s", err.Error())
		}
		log.Printf("log result: %s", decodeString)
	}

	if resp.FunctionError != nil {
		return fmt.Errorf("lambda was invoked but returned error: %s", *resp.FunctionError)
	}
	log.Println("query player damage finished")

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
	lambdaSvc = lambda.New(sess)

	combatlogUUIDs, err := invokeConvert(local)
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}

	err = invokeQueryPlayerDamageDone(combatlogUUIDs[0], local)
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
}
