package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/jonny-rimek/wowmate/services/common/golib"
	"github.com/sirupsen/logrus"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var snsSvc *sns.SNS
var querySvc *timestreamquery.TimestreamQuery

type logData struct {
	ScannedMegabytes int64
	BilledMegabytes  int64
	CombatlogUUID    string
	QueryID          string
}

func handler(ctx aws.Context, e events.SNSEvent) error {
	logData, err := handle(ctx, e)
	if err != nil {
		golib.CanonicalLog(map[string]interface{}{
			"combatlog_uuid":    logData.CombatlogUUID,
			"billed_megabytes":  logData.BilledMegabytes,
			"scanned_megabytes": logData.ScannedMegabytes,
			"query_id":          logData.QueryID,
			"err":               err.Error(),
			"event":             e,
		})
		return err
	}

	golib.CanonicalLog(map[string]interface{}{
		"combatlog_uuid":    logData.CombatlogUUID,
		"billed_megabytes":  logData.BilledMegabytes,
		"scanned_megabytes": logData.ScannedMegabytes,
		"query_id":          logData.QueryID,
	})
	return err
}

func handle(ctx aws.Context, e events.SNSEvent) (logData, error) {
	var logData logData

	topicArn, combatlogUUID, err := validateInput(e)
	if err != nil {
		return logData, err
	}
	logData.CombatlogUUID = combatlogUUID

	//NOTE: AND caster_id LIKE 'Player-%' doesnt work, sprintf tries to format the %
	query := fmt.Sprintf(`
	SELECT 
		SUM(measure_value::bigint) AS damage,
		caster_name,
		caster_id,
		combatlog_uuid
	FROM 
		"wowmate-analytics"."combatlogs" 
	WHERE
		combatlog_uuid = '%v' AND
		(caster_type = '0x512' OR caster_type = '0x511')
	GROUP BY
		caster_name, caster_id, combatlog_uuid
	ORDER BY
		damage DESC
	`, combatlogUUID)
	/*
		WITH dungeon AS (
		    SELECT
				dungeon_name,
		        measure_value::bigint AS dungeon_id,
		        combatlog_uuid
			FROM
				"wowmate-analytics"."combatlogs"
			WHERE
				combatlog_uuid = '08408a8a-1c1d-4ded-9c8e-34a63ad7cdd5'  AND
		        time between ago(15m) and now() AND
		        measure_name = 'dungeon_id'
		),
		key_level AS (
		    SELECT
		        measure_value::bigint AS key_level,
		        combatlog_uuid
			FROM
				"wowmate-analytics"."combatlogs"
			WHERE
				combatlog_uuid = '08408a8a-1c1d-4ded-9c8e-34a63ad7cdd5'  AND
		        time between ago(15m) and now() AND
		        measure_name = 'key_level'
		),
		duration AS (
		    SELECT
		        measure_value::bigint AS duration,
		        combatlog_uuid
			FROM
				"wowmate-analytics"."combatlogs"
			WHERE
				combatlog_uuid = '08408a8a-1c1d-4ded-9c8e-34a63ad7cdd5'  AND
		        time between ago(15m) and now() AND
		        measure_name = 'duration'
		),
		finished AS (
		    SELECT
		        measure_value::bigint AS finished,
		        combatlog_uuid
			FROM
				"wowmate-analytics"."combatlogs"
			WHERE
				combatlog_uuid = '08408a8a-1c1d-4ded-9c8e-34a63ad7cdd5'  AND
		        time between ago(15m) and now() AND
		        measure_name = 'finished'
		),
		damage as (
			SELECT
				SUM(measure_value::bigint) AS damage,
				caster_name,
				caster_id,
				combatlog_uuid
			FROM
				"wowmate-analytics"."combatlogs"
			WHERE
				combatlog_uuid = '08408a8a-1c1d-4ded-9c8e-34a63ad7cdd5' AND
				(caster_type = '0x512' OR caster_type = '0x511') AND
		  		time between ago(15m) and now()
			GROUP BY
				caster_name, caster_id, combatlog_uuid
			ORDER BY
				damage DESC
		)
		SELECT
			*
		FROM
			damage
		JOIN
			dungeon
			ON damage.combatlog_uuid = dungeon.combatlog_uuid
		JOIN
			key_level
		    ON key_level.combatlog_uuid = dungeon.combatlog_uuid
		JOIN
			duration
		    ON duration.combatlog_uuid = dungeon.combatlog_uuid
		JOIN
			finished
		    ON finished.combatlog_uuid = dungeon.combatlog_uuid
	*/

	//CHALLENGE_MODE_END contains duration in milliseconds in the last field

	queryResult, err := golib.TimestreamQuery(ctx, &query, querySvc)
	if err != nil {
		if queryResult == nil {
			logData.QueryID = "queryResult=nil"
			return logData, err
		}
		logData.QueryID = *queryResult.QueryId
		return logData, err
	}

	logData.QueryID = *queryResult.QueryId
	logData.BilledMegabytes = *queryResult.QueryStatus.CumulativeBytesMetered / 1e6 //1.000.000
	logData.ScannedMegabytes = *queryResult.QueryStatus.CumulativeBytesScanned / 1e6

	input, err := json.Marshal(queryResult)
	if err != nil {
		return logData, err
	}

	err = golib.SNSPublishMsg(ctx, snsSvc, string(input), &topicArn)
	if err != nil {
		return logData, err
	}
	return logData, nil
}

func validateInput(e events.SNSEvent) (topicArn string, combatlogUUID string, err error) {
	topicArn = os.Getenv("TOPIC_ARN")
	if topicArn == "" {
		return "", "", fmt.Errorf("arn topic env var is empty")
	}
	logrus.Debug("topicArn: " + topicArn)

	combatlogUUID = e.Records[0].SNS.Message
	/*
		if err != nil {
			return topicArn, "", fmt.Errorf("json unmarschal of sns message failed: %v", err.Error())
		}
	*/

	if combatlogUUID == "" {
		return topicArn, "", fmt.Errorf("combatlog uuid is empty")
	}

	return topicArn, combatlogUUID, nil
}

func main() {
	golib.InitLogging()

	sess, err := session.NewSession()
	if err != nil {
		logrus.Info(fmt.Sprintf("Error creating session: %v", err.Error()))
		return
	}

	snsSvc = sns.New(sess)
	if os.Getenv("LOCAL") == "false" {
		xray.AWS(snsSvc.Client)
	}
	querySvc = timestreamquery.New(sess)

	if os.Getenv("LOCAL") == "false" {
		xray.AWS(querySvc.Client)
	}
	lambda.Start(handler)
}
