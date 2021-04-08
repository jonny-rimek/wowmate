package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/jonny-rimek/wowmate/services/common/golib"
	"log"
	"os"
	"strconv"
	"time"
)

type logData struct {
	Wcu float64
}

var svc *dynamodb.DynamoDB

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

	ddbTableName := os.Getenv("DYNAMODB_TABLE_NAME")
	if ddbTableName == "" {
		return logData, fmt.Errorf("dynamo db table name env var is empty")
	}

	queryResult, err := extractQueryResult(e)

	record, err := convertQueryResult(queryResult)
	if err != nil {
		return logData, err
	}

	response, err := golib.DynamoDBPutItem(ctx, svc, &ddbTableName, record)
	if err != nil {
		return logData, err
	}

	logData.Wcu = *response.ConsumedCapacity.CapacityUnits

	return logData, nil
}

func extractQueryResult(e events.SNSEvent) (*timestreamquery.QueryOutput, error) {
	message := e.Records[0].SNS.Message
	if message == "" {
		return nil, fmt.Errorf("message is empty")
	}

	var result *timestreamquery.QueryOutput

	err := json.Unmarshal([]byte(message), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal sns message which contains the query result: %v", err)
	}

	return result, err
}

//TODO: tests
func convertQueryResult(queryResult *timestreamquery.QueryOutput) (golib.DynamoDBPlayerDamageDone, error) {
	resp := golib.DynamoDBPlayerDamageDone{}

	var summaries []golib.KeysResult

	for i := 0; i < len(queryResult.Rows); i++ {
		dam, err := strconv.Atoi(*queryResult.Rows[i].Data[0].ScalarValue)
		if err != nil {
			return resp, err
		}

		d := golib.KeysResult{
			Damage:   dam,
			Name:     *queryResult.Rows[i].Data[1].ScalarValue,
			PlayerID: *queryResult.Rows[i].Data[2].ScalarValue,
			Class:    "unsupported",
			Specc:    "unsupported",
		}

		summaries = append(summaries, d)
	}
	combatlogUUID := *queryResult.Rows[0].Data[3].ScalarValue

	dungeonName := *queryResult.Rows[0].Data[4].ScalarValue

	dungeonID, err := strconv.Atoi(*queryResult.Rows[0].Data[5].ScalarValue)
	if err != nil {
		return resp, err
	}

	keyLevel, err := strconv.Atoi(*queryResult.Rows[0].Data[6].ScalarValue)
	if err != nil {
		return resp, err
	}

	durationInMilliseconds, err := golib.Atoi64(*queryResult.Rows[0].Data[7].ScalarValue)
	if err != nil {
		return resp, err
	}
	//converts duration to date 1970 + duration, of which I only display the minutes and seconds
	//time duration, doesn't allow mixed formatting like min:seconds
	t := time.Unix(0, durationInMilliseconds*1e6) //milliseconds > nanoseconds

	finished, err := strconv.Atoi(*queryResult.Rows[0].Data[8].ScalarValue)
	if err != nil {
		return resp, err
	}

	twoAffixID, err := strconv.Atoi(*queryResult.Rows[0].Data[9].ScalarValue)
	if err != nil {
		return resp, err
	}

	fourAffixID, err := strconv.Atoi(*queryResult.Rows[0].Data[10].ScalarValue)
	if err != nil {
		return resp, err
	}

	sevenAffixID, err := strconv.Atoi(*queryResult.Rows[0].Data[11].ScalarValue)
	if err != nil {
		return resp, err
	}

	tenAffixID, err := strconv.Atoi(*queryResult.Rows[0].Data[12].ScalarValue)
	if err != nil {
		return resp, err
	}

	resp = golib.DynamoDBPlayerDamageDone{
		Pk:            fmt.Sprintf("LOG#SPECIFIC#%v#OVERALL_PLAYER_DAMAGE", combatlogUUID),
		Sk:            fmt.Sprintf("LOG#SPECIFIC#%v#OVERALL_PLAYER_DAMAGE", combatlogUUID),
		Damage:        summaries,
		Duration:      t.Format("04:05"), //formats to minutes:seconds
		Deaths:        0,                 //TODO:
		Affixes:       fmt.Sprintf("%v, %v, %v, %v", twoAffixID, fourAffixID, sevenAffixID, tenAffixID),
		Keylevel:      keyLevel,
		DungeonName:   dungeonName,
		DungeonID:     dungeonID,
		CombatlogUUID: combatlogUUID,
		Finished:      finished != 0, //if 0 false, else 1
	}
	return resp, err
}

func main() {
	golib.InitLogging()

	sess, err := session.NewSession()
	if err != nil {
		log.Printf("Error creating session: %v", err.Error())
		return
	}

	svc = dynamodb.New(sess)
	xray.AWS(svc.Client)

	lambda.Start(handler)
}
