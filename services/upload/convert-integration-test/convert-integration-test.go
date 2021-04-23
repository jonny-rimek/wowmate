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

func invokeConvert() error {
	functionName := os.Getenv("FUNCTION_NAME")
	if functionName == "" {
		return fmt.Errorf("FUNCTION_NAME not set")
	}
	log.Printf("function name: %s", functionName)

	// invoke lambda via sdk
	input := &lambdaService.InvokeInput{
		FunctionName:   &functionName,
		Payload:        nil,
		LogType:        aws.String(lambdaService.LogTypeTail),                   // returns the log in the response
		InvocationType: aws.String(lambdaService.InvocationTypeRequestResponse), // synchronous - default
	}
	err := input.Validate()
	if err != nil {
		return fmt.Errorf("validating the input failed: %v", err)
	}

	resp, err := lambdaSvc.Invoke(input)
	if err != nil {
		return fmt.Errorf("failed to invoke lambda: %v", err)
	}

	decodeString, err := base64.StdEncoding.DecodeString(*resp.LogResult)
	if err != nil {
		return fmt.Errorf("failed to decode the log: %v", err)
	}
	log.Printf("log result: %s", decodeString)

	if resp.FunctionError != nil {
		return fmt.Errorf("lambda was invoked but returned error: %s", *resp.FunctionError)
	}
	return nil
}

func main() {
	sess, err := session.NewSession()
	if err != nil {
		log.Printf("failed to create a session: %v", err)
		return
	}
	lambdaSvc = lambdaService.New(sess)

	err = invokeConvert()
	if err != nil {
		return
	}
}
