package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/timestreamquery"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/jonny-rimek/wowmate/services/common/golib"
	"log"
	"strconv"
)

var client *dynamodb.Client

type logData struct {
	Wcu float64
}

func handler(ctx aws.Context, e events.SNSEvent) error {
	logData, err := handle(ctx, e)
	if err != nil {
		golib.CanonicalLog(map[string]interface{}{
			"wcu":   logData.Wcu,
			"err":   err.Error(),
			"event": e,
		})
		return err
	}

	golib.CanonicalLog(map[string]interface{}{
		"wcu": logData.Wcu,
	})
	return err
}

func handle(ctx aws.Context, e events.SNSEvent) (logData, error) {
	var logData logData

	/*
		ddbTableName := os.Getenv("DYNAMODB_TABLE_NAME")
		if ddbTableName == "" {
			return logData, fmt.Errorf("dynamo db table name env var is empty")
		}

		summaries, combatlogUUID, err := extractInput(e)
		if err != nil {
			return logData, err
		}

		_ = convertInput(combatlogUUID, summaries)

		//TODO:
		//convert input
		//write to timestream


	*/

	return logData, nil
}

func convertInput(combatlogUUID string, summaries []golib.KeysResult) {
}

func extractInput(e events.SNSEvent) ([]golib.KeysResult, string, error) {
	message := e.Records[0].SNS.Message
	if message == "" {
		return nil, "", fmt.Errorf("message is empty")
	}
	log.Println("message:" + message)

	var result timestreamquery.QueryOutput

	err := json.Unmarshal([]byte(message), &result)
	if err != nil {
		return nil, "", err
	}

	var combatlogUUID string
	var summaries []golib.KeysResult

	for i := 0; i < len(result.Rows); i++ {
		dam, err := strconv.Atoi(*result.Rows[i].Data[0].ScalarValue)
		if err != nil {
			return nil, "", err
		}

		d := golib.KeysResult{
			Damage:   dam,
			Name:     *result.Rows[i].Data[1].ScalarValue,
			PlayerID: *result.Rows[i].Data[2].ScalarValue,
		}
		combatlogUUID = *result.Rows[i].Data[3].ScalarValue

		summaries = append(summaries, d)
	}

	return summaries, combatlogUUID, err
}

func main() {
	golib.InitLogging()

	lambda.Start(handler)
}
