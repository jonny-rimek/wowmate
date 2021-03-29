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
	"log"
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
	log.Printf("SNS EVENT: %+v", e)

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

	/* THIS DOESNT HELP
	tr := &http.Transport{
		ResponseHeaderTimeout: 20 * time.Second,
		// Using DefaultTransport values for other parameters: https://golang.org/pkg/net/http/#RoundTripper
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			KeepAlive: 30 * time.Second,
			DualStack: true,
			Timeout:   30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	// So client makes HTTP/2 requests
	http2.ConfigureTransport(tr)

	sess, err := session.NewSession(&aws.Config{
		MaxRetries: aws.Int(10),
		HTTPClient: &http.Client{Transport: tr},
	})

	*/
	sess, err := session.NewSession()
	if err != nil {
		logrus.Info(fmt.Sprintf("Error creating session: %v", err.Error()))
		return
	}

	snsSvc = sns.New(sess)
	xray.AWS(snsSvc.Client)
	querySvc = timestreamquery.New(sess)
	xray.AWS(querySvc.Client)

	lambda.Start(handler)
}
