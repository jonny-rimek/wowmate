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
	"math/rand"
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

	summaries, combatlogUUID, err := extractInput(e)
	if err != nil {
		return logData, err
	}

	record := convertInput(combatlogUUID, summaries)

	response, err := golib.DynamoDBPutItem(ctx, svc, &ddbTableName, record)
	if err != nil {
		return logData, err
	}

	logData.Wcu = *response.ConsumedCapacity.CapacityUnits

	return logData, nil
}

//TODO: tests
func convertInput(combatlogUUID string, summaries []golib.KeysResult) golib.DynamoDBKeys {
	rand.Seed(time.Now().UnixNano())
	min := 2
	max := 26
	keylevel := rand.Intn(max-min+1) + min

	min = 50
	max = 150
	timePercent := rand.Intn(max-min+1) + min

	record := golib.DynamoDBKeys{
		Pk:            "LOG#S2", //IMPROVE: dynamic season
		Sk:            fmt.Sprintf("%02d#%v#%v", keylevel, timePercent, combatlogUUID),
		Damage:        summaries,
		Gsi1pk:        "LOG#S2#2291",
		Gsi1sk:        fmt.Sprintf("%02d#%v#%v", keylevel, timePercent, combatlogUUID),
		Duration:      "34:59 +0:01",
		Deaths:        1,
		Affixes:       "tyrannical, explosive, storming, prideful",
		Keylevel:      keylevel,
		DungeonName:   "De Other Site",
		DungeonID:     2291,
		CombatlogUUID: combatlogUUID,
		//TODO: get real keylevel, dungeon id, dungeon name, deaths, deplete, duration
	}
	return record
}

//TODO: tests
func extractInput(e events.SNSEvent) ([]golib.KeysResult, string, error) {
	message := e.Records[0].SNS.Message
	if message == "" {
		return nil, "", fmt.Errorf("message is empty")
	}

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

	sess, err := session.NewSession()
	if err != nil {
		log.Printf("Error creating session: %v", err.Error())
		return
	}

	svc = dynamodb.New(sess)
	xray.AWS(svc.Client)

	lambda.Start(handler)
}
