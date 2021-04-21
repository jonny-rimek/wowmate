package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codedeploy"
	lambdaService "github.com/aws/aws-sdk-go/service/lambda"
)

var svc *codedeploy.CodeDeploy
var lambdaSvc *lambdaService.Lambda

type codeDeployEvent struct {
	DeploymentId                  string `json:"deploymentId"`
	LifecycleEventHookExecutionId string `json:"lifecycleEventHookExecutionId"`
}

func handler(e codeDeployEvent) error {
	params := &codedeploy.PutLifecycleEventHookExecutionStatusInput{
		DeploymentId:                  &e.DeploymentId,
		LifecycleEventHookExecutionId: &e.LifecycleEventHookExecutionId,
	}
	err := handle()
	if err != nil {
		log.Println(err)
		params.Status = aws.String(codedeploy.LifecycleEventStatusFailed)
	} else {
		params.Status = aws.String(codedeploy.LifecycleEventStatusSucceeded)
	}

	_, err = svc.PutLifecycleEventHookExecutionStatus(params)
	if err != nil {
		return fmt.Errorf("failed putting the lifecycle event hook execution status. the status was %s", *params.Status)
	}

	return nil
}

func handle() error {
	functionName := os.Getenv("FUNCTION_NAME")
	if functionName == "" {
		return fmt.Errorf("FUNCTION_NAME not set")
	}
	log.Printf("function name: %s", functionName)

	// invoke lambda via sdk
	input := &lambdaService.InvokeInput{
		FunctionName: &functionName,
		Payload:      nil,
		LogType:      aws.String(lambdaService.LogTypeTail), // returns the log in the response
		// Qualifier:    &functionVersion,

		// ClientContext:  nil,
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
		return
	}
	svc = codedeploy.New(sess)
	lambdaSvc = lambdaService.New(sess)
	lambda.Start(handler)
}
