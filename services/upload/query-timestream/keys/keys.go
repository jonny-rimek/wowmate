package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/aws/aws-xray-sdk-go/xraylog"
	"github.com/jonny-rimek/wowmate/services/common/golib"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var snsSvc *sns.SNS
var querySvc *timestreamquery.TimestreamQuery

type logData struct {
	ScannedMegabytes int64
	BilledMegabytes  int64
	CombatlogHash    string
	QueryID          string
}

func handler(ctx aws.Context, e events.SNSEvent) error {
	logData, err := handle(ctx, e)
	if err != nil {
		//goland:noinspection GoNilness
		golib.CanonicalLog(map[string]interface{}{
			"combatlog_hash":    logData.CombatlogHash,
			"billed_megabytes":  logData.BilledMegabytes,
			"scanned_megabytes": logData.ScannedMegabytes,
			"query_id":          logData.QueryID,
			"err":               err.Error(),
			"event":             e,
		})
		return err
	}

	golib.CanonicalLog(map[string]interface{}{
		"combatlog_hash":    logData.CombatlogHash,
		"billed_megabytes":  logData.BilledMegabytes,
		"scanned_megabytes": logData.ScannedMegabytes,
		"query_id":          logData.QueryID,
	})
	return err
}

func handle(ctx aws.Context, e events.SNSEvent) (logData, error) {
	var logData logData

	topicArn, combatlogHash, err := validateInput(e)
	if err != nil {
		return logData, err
	}
	logData.CombatlogHash = combatlogHash

	// NOTE: AND caster_id LIKE 'Player-%' doesnt work, sprintf tries to format the %
	query := fmt.Sprintf(`
		WITH dungeon AS (
		    SELECT
				dungeon_name,
		        measure_value::bigint AS dungeon_id,
		        combatlog_hash
			FROM
				"wowmate-analytics"."combatlogs"
			WHERE
				combatlog_hash = '%v'  AND
		        time between ago(3m) and now() AND
		        measure_name = 'dungeon_id'
		),
		key_level AS (
		    SELECT
		        measure_value::bigint AS key_level,
		        combatlog_hash
			FROM
				"wowmate-analytics"."combatlogs"
			WHERE
				combatlog_hash = '%v'  AND
		        time between ago(3m) and now() AND
		        measure_name = 'key_level'
		),
		duration AS (
		    SELECT
		        measure_value::bigint AS duration,
		        combatlog_hash
			FROM
				"wowmate-analytics"."combatlogs"
			WHERE
				combatlog_hash = '%v'  AND
		        time between ago(3m) and now() AND
		        measure_name = 'duration'
		),
		finished AS (
		    SELECT
		        measure_value::bigint AS finished,
		        combatlog_hash
			FROM
				"wowmate-analytics"."combatlogs"
			WHERE
				combatlog_hash = '%v'  AND
		        time between ago(3m) and now() AND
		        measure_name = 'finished' AND
				measure_value::bigint != 0 -- don't add non finished keys to the top keys list'
		),
        two_affix_id AS (
		    SELECT
		        measure_value::bigint AS two_affix_id,
		        combatlog_hash
			FROM
				"wowmate-analytics"."combatlogs"
			WHERE
				combatlog_hash = '%v'  AND
		        time between ago(3m) and now() AND
		        measure_name = 'two_affix_id'
		),
        four_affix_id AS (
		    SELECT
		        measure_value::bigint AS four_affix_id,
		        combatlog_hash
			FROM
				"wowmate-analytics"."combatlogs"
			WHERE
				combatlog_hash = '%v'  AND
		        time between ago(3m) and now() AND
		        measure_name = 'four_affix_id'
		),
        seven_affix_id AS (
		    SELECT
		        measure_value::bigint AS seven_affix_id,
		        combatlog_hash
			FROM
				"wowmate-analytics"."combatlogs"
			WHERE
				combatlog_hash = '%v'  AND
		        time between ago(3m) and now() AND
		        measure_name = 'seven_affix_id'
		),
        ten_affix_id AS (
		    SELECT
		        measure_value::bigint AS ten_affix_id,
		        combatlog_hash
			FROM
				"wowmate-analytics"."combatlogs"
			WHERE
				combatlog_hash = '%v'  AND
		        time between ago(3m) and now() AND
		        measure_name = 'ten_affix_id'
		),        
        patch AS (
		    SELECT
		        measure_value::varchar AS patch,
		        combatlog_hash
			FROM
				"wowmate-analytics"."combatlogs"
			WHERE
				combatlog_hash = '%v'  AND
		        time between ago(3m) and now() AND
		        measure_name = 'patch'
		),        
        date AS (
		    SELECT
		        measure_value::bigint AS date,
		        combatlog_hash
			FROM
				"wowmate-analytics"."combatlogs"
			WHERE
				combatlog_hash = '%v'  AND
		        time between ago(3m) and now() AND
		        measure_name = 'date'
		),        
		damage as (
			SELECT
				SUM(measure_value::bigint) AS damage,
				caster_name,
				caster_id,
				combatlog_hash
			FROM
				"wowmate-analytics"."combatlogs"
			WHERE
				combatlog_hash = '%v' AND
				(caster_type = '0x512' OR caster_type = '0x511') AND
				caster_id != '0000000000000000' AND -- sometime the caster_id is empty, dunno why
		  		time between ago(3m) and now()
			GROUP BY
				caster_name, caster_id, combatlog_hash
			ORDER BY
				damage DESC
		)
		SELECT
			damage, 
			caster_name, 
			caster_id, 
			damage.combatlog_hash, 
			dungeon_name, 
			dungeon_id,
			key_level, 
			duration, 
			finished, 
			two_affix_id, 
			four_affix_id, 
			seven_affix_id, 
			ten_affix_id,
			patch,
			date
		FROM
			damage
		JOIN
			dungeon
			ON damage.combatlog_hash = dungeon.combatlog_hash
		JOIN
			key_level
		    ON key_level.combatlog_hash = dungeon.combatlog_hash
		JOIN
			duration
		    ON duration.combatlog_hash = dungeon.combatlog_hash
		JOIN
			finished
		    ON finished.combatlog_hash = dungeon.combatlog_hash
		JOIN
			two_affix_id
		    ON two_affix_id.combatlog_hash = dungeon.combatlog_hash
        JOIN
			four_affix_id
		    ON four_affix_id.combatlog_hash = dungeon.combatlog_hash
        JOIN
			seven_affix_id
		    ON seven_affix_id.combatlog_hash = dungeon.combatlog_hash
        JOIN
			ten_affix_id
		    ON ten_affix_id.combatlog_hash = dungeon.combatlog_hash            
        JOIN
			patch
		    ON patch.combatlog_hash = dungeon.combatlog_hash            
        JOIN
			date
		    ON date.combatlog_hash = dungeon.combatlog_hash            
		`,
		combatlogHash,
		combatlogHash,
		combatlogHash,
		combatlogHash,
		combatlogHash,
		combatlogHash,
		combatlogHash,
		combatlogHash,
		combatlogHash,
		combatlogHash,
		combatlogHash,
	)

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
	logData.BilledMegabytes = *queryResult.QueryStatus.CumulativeBytesMetered / 1e6 // 1.000.000
	logData.ScannedMegabytes = *queryResult.QueryStatus.CumulativeBytesScanned / 1e6

	if len(queryResult.Rows) == 0 {
		// query was empty which is fine, but don't push it to sns
		return logData, nil
	}

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

func validateInput(e events.SNSEvent) (topicArn string, combatlogHash string, err error) {
	topicArn = os.Getenv("TOPIC_ARN")
	if topicArn == "" {
		return "", "", fmt.Errorf("arn topic env var is empty")
	}
	logrus.Debug("topicArn: " + topicArn)

	combatlogHash = e.Records[0].SNS.Message
	if combatlogHash == "" {
		return topicArn, "", fmt.Errorf("combatlog hash is empty")
	}

	return topicArn, combatlogHash, nil
}

func main() {
	golib.InitLogging()

	sess, err := session.NewSession()
	if err != nil {
		logrus.Info(fmt.Sprintf("Error creating session: %v", err.Error()))
		return
	}

	xray.SetLogger(xraylog.NullLogger)

	snsSvc = sns.New(sess)
	if os.Getenv("LOCAL") != "true" {
		xray.AWS(snsSvc.Client)
	}

	querySvc = timestreamquery.New(sess)
	if os.Getenv("LOCAL") != "true" {
		xray.AWS(querySvc.Client)
	}

	lambda.Start(handler)
}
