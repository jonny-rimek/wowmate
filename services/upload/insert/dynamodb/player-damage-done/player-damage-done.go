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
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"strconv"
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

	record, err := extractInput(e)
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

//TODO: tests
func extractInput(e events.SNSEvent) (golib.DynamoDBPlayerDamageDone, error) {
	resp := golib.DynamoDBPlayerDamageDone{}

	//extract sns part into extra func
	message := e.Records[0].SNS.Message
	if message == "" {
		return resp, fmt.Errorf("message is empty")
	}
	logrus.Debug("message:" + message)

	var result timestreamquery.QueryOutput

	err := json.Unmarshal([]byte(message), &result)
	if err != nil {
		return resp, err
	}

	var summaries []golib.KeysResult

	for i := 0; i < len(result.Rows); i++ {
		dam, err := strconv.Atoi(*result.Rows[i].Data[0].ScalarValue)
		if err != nil {
			return resp, err
		}

		d := golib.KeysResult{
			Damage:   dam,
			Name:     *result.Rows[i].Data[1].ScalarValue,
			PlayerID: *result.Rows[i].Data[2].ScalarValue,
		}
		//TODO:
		//	- convert affix ids to values
		//	- update query keys and insert keys lambdas

		summaries = append(summaries, d)
	}
	combatlogUUID := *result.Rows[0].Data[3].ScalarValue
	dungeonName := *result.Rows[0].Data[4].ScalarValue
	dungeonID, err := strconv.Atoi(*result.Rows[0].Data[5].ScalarValue)
	if err != nil {
		return resp, err
	}
	keyLevel, err := strconv.Atoi(*result.Rows[0].Data[6].ScalarValue)
	if err != nil {
		return resp, err
	}
	durationInMilliseconds, err := strconv.Atoi(*result.Rows[0].Data[7].ScalarValue)
	if err != nil {
		return resp, err
	}
	finished, err := strconv.Atoi(*result.Rows[0].Data[8].ScalarValue)
	if err != nil {
		return resp, err
	}
	twoAffixID, err := strconv.Atoi(*result.Rows[0].Data[9].ScalarValue)
	if err != nil {
		return resp, err
	}
	fourAffixID, err := strconv.Atoi(*result.Rows[0].Data[10].ScalarValue)
	if err != nil {
		return resp, err
	}
	sevenAffixID, err := strconv.Atoi(*result.Rows[0].Data[11].ScalarValue)
	if err != nil {
		return resp, err
	}
	tenAffixID, err := strconv.Atoi(*result.Rows[0].Data[12].ScalarValue)
	if err != nil {
		return resp, err
	}

	resp = golib.DynamoDBPlayerDamageDone{
		Pk:            fmt.Sprintf("LOG#SPECIFIC#%v#OVERALL_PLAYER_DAMAGE", combatlogUUID),
		Sk:            fmt.Sprintf("LOG#SPECIFIC#%v#OVERALL_PLAYER_DAMAGE", combatlogUUID),
		Damage:        summaries,
		Duration:      strconv.Itoa(durationInMilliseconds),
		Deaths:        0, //TODO:
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
