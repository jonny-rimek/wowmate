package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	lambdaService "github.com/aws/aws-sdk-go/service/lambda"
)

var lambdaSvc *lambdaService.Lambda

func invokeConvert() error {
	functionName := "ConvertLambda3540DCCB"
	// invoke lambda via sdk
	input := &lambdaService.InvokeInput{
		FunctionName:   &functionName,
		Payload:        nil,
		InvocationType: aws.String(lambdaService.InvocationTypeRequestResponse), // synchronous - default

		// DOESN'T WORK LOCALLY
		// LogType:        aws.String(lambdaService.LogTypeTail),                   // returns the log in the response
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

	// decodeString, err := base64.StdEncoding.DecodeString(*resp.LogResult)
	// if err != nil {
	// 	return fmt.Errorf("failed to decode the log: %v", err)
	// }
	// log.Printf("log result: %s", decodeString)

	if resp.FunctionError != nil {
		return fmt.Errorf("lambda was invoked but returned error: %s", *resp.FunctionError)
	}
	return nil
}

func main() {
	sess, err := session.NewSession(
		&aws.Config{
			Endpoint:   aws.String("http://127.0.0.1:3001"),
			DisableSSL: aws.Bool(true),
			Region:     aws.String("us-east-1"),
			MaxRetries: aws.Int(0),
		})
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	lambdaSvc = lambdaService.New(sess) // , aws.NewConfig().WithEndpoint("http://localhost:3001").WithRegion("us-east-1"))

	err = invokeConvert()
	if err != nil {
		log.Printf("failed to invoke lambda: %s", err)
		os.Exit(1)
	}
	log.Println("invoke finished")
}
