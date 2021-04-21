package main

import (
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
	lambdaARN := os.Getenv("LAMBDA_ARN")
	if lambdaARN == "" {
		return fmt.Errorf("lambda_arn not set")
	}
	lambdaVersion := os.Getenv("LAMBDA_VERSION")
	if lambdaVersion == "" {
		return fmt.Errorf("lambda_version not set")
	}

	// invoke lambda via sdk
	lambdaSvc.Invoke(&lambdaService.InvokeInput{
		FunctionName: &lambdaARN,
		Payload:      nil,
		Qualifier:    &lambdaVersion,

		// ClientContext:  nil,
		// InvocationType: nil, // default synchronous
		// LogType:        nil, // returns the log in the response
	})
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
