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

func invokeQueryPlayerDamageDone(local bool) error {
	var functionName string
	var logType *string
	var payload []byte
	payload = nil
	var combatlogUUID string

	if local == true {
		functionName = "QueryPlayerDamageDoneLambda98AFC037"
		logType = nil
		combatlogUUID = "e408d5b8-2c98-43bc-9a38-83f4e3c94383"
		// TODO: the result is probably empty because of the time predicate ago(15m) upload new log
		// 	and test again, figure out a way to handle it locally and in ci
	} else {
		combatlogUUID = "" // TODO
		functionName = ""  // TODO
		logType = aws.String(lambdaService.LogTypeTail)
	}
	payload = []byte(fmt.Sprintf(`
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

	input := &lambdaService.InvokeInput{
		FunctionName:   &functionName,
		Payload:        payload,
		InvocationType: aws.String(lambdaService.InvocationTypeRequestResponse), // synchronous - default
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
		log.Printf("failed creating a session %s", err.Error())
		os.Exit(1)
	}
	lambdaSvc = lambdaService.New(sess)

	err = invokeQueryPlayerDamageDone(local)
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
	log.Println("invoke finished")
}
