package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/jonny-rimek/wowmate/services/common/golib"
	"github.com/sirupsen/logrus"
)

type logData struct {
	Wcu float64
}

var svc *dynamodb.DynamoDB

func handler(ctx aws.Context, e events.SNSEvent) error {
	logData, err := handle(ctx, e)
	if err != nil {
		//goland:noinspection GoNilness
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

// TODO: tests
func convertQueryResult(queryResult *timestreamquery.QueryOutput) (golib.DynamoDBKeys, error) {
	resp := golib.DynamoDBKeys{}

	var summaries []golib.PlayerDamage

	for i := 0; i < len(queryResult.Rows); i++ {
		dam, err := strconv.Atoi(*queryResult.Rows[i].Data[0].ScalarValue)
		if err != nil {
			return resp, err
		}

		d := golib.PlayerDamage{
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
	// converts duration to date 1970 + duration, of which I only display the minutes and seconds
	// time.Duration, doesn't allow mixed formatting like min:seconds
	t := time.Unix(0, durationInMilliseconds*1e6) // milliseconds > nanoseconds

	durAsPercent, intime, err := golib.TimedAsPercent(dungeonID, float64(durationInMilliseconds))
	if err != nil {
		return resp, err
	}

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

	patch := *queryResult.Rows[0].Data[13].ScalarValue

	resp = golib.DynamoDBKeys{
		// hardcoding the patch like that might be too granular, maybe it makes more sense that e.g. 9.0.2 and 9.0.5 are both S1
		Pk: fmt.Sprintf("LOG#KEY#%s", patch),
		Sk: fmt.Sprintf("%02d#%3.6f#%v", keyLevel, durAsPercent, combatlogUUID),
		// sorting in dynamoDB is achieved via the sort key, in order to sort by key level and within the key level by
		// time I'm printing the value as string and sort the string.
		// As I'm sorting descending I can't just print the duration in milliseconds.
		// instead I print the duration as percent in relation to the intime duration
		Damage:        summaries,
		Gsi1pk:        fmt.Sprintf("LOG#KEY#%s#%v", patch, dungeonID),
		Gsi1sk:        fmt.Sprintf("%02d#%3.6f#%v", keyLevel, durAsPercent, combatlogUUID),
		Duration:      t.Format("04:05"), // formats to minutes:seconds
		Deaths:        0,                 // TODO:
		Affixes:       golib.AffixIDsToString(twoAffixID, fourAffixID, sevenAffixID, tenAffixID),
		Keylevel:      keyLevel,
		DungeonName:   dungeonName,
		DungeonID:     dungeonID,
		CombatlogUUID: combatlogUUID,
		Finished:      finished != 0, // if 0 false, else true
		Intime:        intime,
	}
	return resp, err
}

func main() {
	golib.InitLogging()

	sess, err := session.NewSession()
	if err != nil {
		logrus.Info(fmt.Sprintf("Error creating session: %v", err.Error()))
		return
	}

	svc = dynamodb.New(sess)
	if os.Getenv("LOCAL") == "false" {
		xray.AWS(svc.Client)
	}

	lambda.Start(handler)
}
