package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/jonny-rimek/wowmate/services/common/golib"
	"github.com/sirupsen/logrus"
)

type logData struct {
	Rcu           float64
	CombatlogUUID string
	EmptyQuery    bool
}

var svc *dynamodb.DynamoDB

func handler(ctx aws.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	response, logData, err := handle(ctx, request)
	if err != nil {
		//goland:noinspection GoNilness
		golib.CanonicalLog(map[string]interface{}{
			"rcu":            logData.Rcu,
			"combatlog_uuid": logData.CombatlogUUID,
			"err":            err.Error(),
			"empty_query":    logData.EmptyQuery,
			"event":          request,
		})
		return response, err
	}

	golib.CanonicalLog(map[string]interface{}{
		"rcu":            logData.Rcu,
		"combatlog_uuid": logData.CombatlogUUID,
		"empty_query":    logData.EmptyQuery,
	})

	return response, err
}

func handle(ctx aws.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, logData, error) {
	var logData logData
	ddbTableName := os.Getenv("DYNAMODB_TABLE_NAME")
	if ddbTableName == "" {
		return golib.AGW500(), logData, fmt.Errorf("failed dynamodb table name env var is empty")
	}

	combatlogUUID, err := checkInput(request.PathParameters)
	if err != nil {
		return golib.AGW400(), logData, err
	}
	logData.CombatlogUUID = combatlogUUID

	input := &dynamodb.GetItemInput{
		TableName: &ddbTableName,
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(fmt.Sprintf("LOG#KEY#%v#OVERALL_PLAYER_DAMAGE", combatlogUUID)),
			},
			"sk": {
				S: aws.String(fmt.Sprintf("LOG#KEY#%v#OVERALL_PLAYER_DAMAGE", combatlogUUID)),
			},
		},
		ReturnConsumedCapacity: aws.String("TOTAL"),
	}

	result, err := golib.DynamoDBGetItem(ctx, svc, input)
	if err != nil {
		return golib.AGW500(), logData, err
	}
	logData.Rcu = *result.ConsumedCapacity.CapacityUnits

	// check if query from dynamodb is empty and return 404
	if result.Item == nil {
		logData.EmptyQuery = true
		return golib.AGW404(), logData, nil
	}

	json, err := golib.PlayerDamageDoneToJson(result)
	if err != nil {
		return golib.AGW500(), logData, err
	}

	return golib.AGW200(json, map[string]string{
		"Content-type": "application/json",
	}), logData, nil
}

func checkInput(input map[string]string) (string, error) {
	if input["combatlog_uuid"] == "" {
		return "", fmt.Errorf("combatloguuid is empty")
	}
	return input["combatlog_uuid"], nil
}

func init() {
	// I don't get when it makes sense to use init the docs doesnt explain it
	// I tried starting the session here, but no performance difference
	// https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html#golang-handler-state
}

func main() {
	golib.InitLogging()

	sess, err := session.NewSession()
	if err != nil {
		logrus.Infof("failed to start a session: %v", err)
		return
	}
	svc = dynamodb.New(sess)

	if os.Getenv("LOCAL") == "false" {
		xray.AWS(svc.Client)
	}

	lambda.Start(handler)
}
