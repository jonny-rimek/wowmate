package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/lambda"
)

type dynamodbDedup struct {
	Pk string `json:"pk"`
	Sk string `json:"sk"`
}

type convertOutput struct {
	Hashes []uint64
}

var lambdaSvc *lambda.Lambda
var ddbSvc *dynamodb.DynamoDB

func invokeConvert(local bool) (*convertOutput, error) {
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

	response, err := invokeLambda(local, functionName, payload, logType)
	if err != nil {
		return nil, err
	}
	// log.Printf("%s", response)

	// parse output and convert
	var r convertOutput

	err = json.Unmarshal(response, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %s", err.Error())
	}

	log.Println("convert finished")
	return &r, nil
}

func invokeQueryKeys(combatlogHash string, local bool) error {
	log.Println("invoke query keys")
	var functionName string
	var logType *string

	if local == true {
		functionName = "QueryKeysLambda58DE9A2E"
		logType = nil
	} else {
		functionName = "wm-preprod-QueryKeysLambda58DE9A2E-1BFBY8IVLR3BX"
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
		`, combatlogHash))

	_, err := invokeLambda(local, functionName, payload, logType)
	if err != nil {
		return err
	}
	log.Println("query keys finished")

	return nil
}

func invokeInsertKeys(local bool) error {
	log.Println("invoke insert keys")
	var functionName string
	var logType *string

	if local == true {
		functionName = "InsertKeysToDynamodbLambda15825024"
		logType = nil
	} else {
		functionName = "wm-preprod-InsertKeysToDynamodbLambda15825024-81XD2VFVV093"
		logType = aws.String(lambda.LogTypeTail)
	}

	payload, err := ioutil.ReadFile("insertKeysToDynamodbEvent.json")
	if err != nil {
		return fmt.Errorf("failed reading the file: %s", err.Error())
	}

	_, err = invokeLambda(local, functionName, payload, logType)
	if err != nil {
		return err
	}
	log.Println("insert keys finished")

	return nil
}

func invokeQueryPlayerDamageDone(combatlogHash string, local bool) error {
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
		`, combatlogHash))

	_, err := invokeLambda(local, functionName, payload, logType)
	if err != nil {
		return err
	}
	log.Println("query player damage finished")

	return nil
}

func invokeInsertPlayerDamageDone(local bool) error {
	log.Println("invoke insert player damage")
	var functionName string
	var logType *string

	if local == true {
		functionName = "InsertPlayerDamageDoneToDynamodbLambda659E0DA7"
		logType = nil
	} else {
		functionName = "wm-preprod-InsertPlayerDamageDoneToDynamodbLambda6-1AULMHM0IL8N3"
		logType = aws.String(lambda.LogTypeTail)
	}

	payload, err := ioutil.ReadFile("insertPlayerDamageDoneToDynamodbEvent.json")
	if err != nil {
		return fmt.Errorf("failed reading the file: %s", err.Error())
	}

	_, err = invokeLambda(local, functionName, payload, logType)
	if err != nil {
		return err
	}
	log.Println("insert player damage finished")

	return nil
}

func invokeLambda(local bool, functionName string, payload []byte, logType *string) (responsePayload []byte, err error) {
	input := &lambda.InvokeInput{
		FunctionName:   &functionName,
		Payload:        payload,
		InvocationType: aws.String(lambda.InvocationTypeRequestResponse), // synchronous - default
		// it needs to be synchronous so I can assert on the response if an error was returned or not
		// even if in aws it's called async via SNS
		LogType: logType, // returns the log in the response
	}

	err = input.Validate()
	if err != nil {
		return nil, fmt.Errorf("validating the input failed: %s", err.Error())
	}

	resp, err := lambdaSvc.Invoke(input)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke lambda: %s", err.Error())
	}
	log.Printf("%s", resp.Payload)

	if local != true {
		// log doesn't work locally
		decodeString, err := base64.StdEncoding.DecodeString(*resp.LogResult)
		if err != nil {
			return nil, fmt.Errorf("failed to decode the log: %s", err.Error())
		}
		log.Printf("log result: %s", decodeString)
	}

	if resp.FunctionError != nil {
		return nil, fmt.Errorf("lambda was invoked but returned error: %s", *resp.FunctionError)
	}
	return resp.Payload, nil
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

	// i dont want to use a local ddb
	sess2, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
	ddbSvc = dynamodb.New(sess2)

	c, err := invokeConvert(local)
	if err != nil {
		handleError(err)
	}
	hash := strconv.FormatUint(c.Hashes[0], 10)

	// this is not an ideal solution, it only deletes the record inside the test
	// if the same file is uploaded from somewhere else it will fail the test
	// but that's a problem for future me
	err = ddbDelete(c.Hashes, local)
	if err != nil {
		handleError(err)
	}

	err = invokeQueryKeys(hash, local)
	if err != nil {
		handleError(err)
	}

	err = invokeInsertKeys(local)
	if err != nil {
		handleError(err)
	}

	err = invokeQueryPlayerDamageDone(hash, local)
	if err != nil {
		handleError(err)
	}

	err = invokeInsertPlayerDamageDone(local)
	if err != nil {
		handleError(err)
	}
}

func ddbDelete(hashes []uint64, local bool) error {
	var ddbTable string
	if local == true {
		ddbTable = "wm-dev-DynamoDBtableF8E87752-HSV525WR7KN3"
	} else {
		ddbTable = "wm-preprod-DynamoDBtableF8E87752-XIQBZHCM8YN4"
	}

	for _, hash := range hashes {

		dd := dynamodbDedup{
			Pk: fmt.Sprintf("DEDUP#%d", hash),
			Sk: fmt.Sprintf("DEDUP#%d", hash),
		}
		d, err := dynamodbattribute.MarshalMap(dd)
		if err != nil {
			return fmt.Errorf("failed to marshal delete input %v", err)
		}

		_, err = ddbSvc.DeleteItem(&dynamodb.DeleteItemInput{
			Key:       d,
			TableName: &ddbTable,
		})
		if err != nil {
			return fmt.Errorf("failed to delete from ddb %v", err)
		}
	}
	return nil
}

func handleError(err error) {
	log.Printf("%s", err)
	os.Exit(1)
}
