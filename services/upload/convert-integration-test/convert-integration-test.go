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
	if local == true {
		functionName = "ConvertLambda3540DCCB"
		logType = nil
	} else {
		functionName = "wm-preprod-ConvertLambda3540DCCB-7AWDGGCRUN1G"
		logType = aws.String(lambdaService.LogTypeTail)
	}

	// TODO: add payload s3 event
	input := &lambdaService.InvokeInput{
		FunctionName:   &functionName,
		Payload:        nil,
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
