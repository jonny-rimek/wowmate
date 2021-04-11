package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/jonny-rimek/wowmate/services/common/golib"
	"github.com/sirupsen/logrus"
)

// this could just be a map[string]interface{}, but I kinda prefer to have a set structure
type logData struct {
	Rcu           float64
	DungeonID     int
	EmptyQuery    bool
	FirstPage     bool
	SortAscending bool
	InputNextSk   string
	InputPrevSK   string
}

type paginatedQueryInput struct {
	request      events.APIGatewayV2HTTPRequest
	ddbTableName *string
	dungeonID    int
}

type paginatedQueryOutput struct {
	queryInput    dynamodb.QueryInput
	firstPage     bool
	sortAscending bool
	inputNextSK   string
	inputPrevSK   string
}

var svc *dynamodb.DynamoDB

func handler(ctx aws.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	response, logData, err := handle(ctx, request)
	if err != nil {
		//goland:noinspection GoNilness
		golib.CanonicalLog(map[string]interface{}{
			"rcu":            logData.Rcu,
			"dungeon_id":     logData.DungeonID,
			"err":            err.Error(),
			"empty_query":    logData.EmptyQuery,
			"first_page":     logData.FirstPage,
			"sort_ascending": logData.SortAscending,
			"input_next_sk":  logData.InputNextSk,
			"input_prev_sk":  logData.InputPrevSK,
			"event":          request,
		})
		return response, err
	}

	golib.CanonicalLog(map[string]interface{}{
		"rcu":            logData.Rcu,
		"dungeon_id":     logData.DungeonID,
		"empty_query":    logData.EmptyQuery,
		"first_page":     logData.FirstPage,
		"sort_ascending": logData.SortAscending,
		"input_next_sk":  logData.InputNextSk,
		"input_prev_sk":  logData.InputPrevSK,
	})
	return response, err
}

func handle(ctx aws.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, logData, error) {
	var logData logData
	ddbTableName := os.Getenv("DYNAMODB_TABLE_NAME")
	if ddbTableName == "" {
		return golib.AGW500(),
			logData, fmt.Errorf("failed dynamodb table name env var is empty")
	}

	dungeonID, err := checkInput(request.PathParameters["dungeon_id"])
	if err != nil {
		return golib.AGW200("", nil), logData, nil
	}
	logData.DungeonID = dungeonID

	paginatedQuery, err := paginatedQuery(paginatedQueryInput{
		request:      request,
		ddbTableName: &ddbTableName,
		dungeonID:    dungeonID,
	})
	if err != nil {
		return golib.AGW500(), logData, err
	}
	logData.FirstPage = paginatedQuery.firstPage
	logData.SortAscending = paginatedQuery.sortAscending
	logData.InputNextSk = paginatedQuery.inputNextSK
	logData.InputPrevSK = paginatedQuery.inputPrevSK

	result, err := golib.DynamoDBQuery(ctx, svc, paginatedQuery.queryInput)
	if err != nil {
		return golib.AGW500(), logData, err
	}
	logData.Rcu = *result.ConsumedCapacity.CapacityUnits

	if len(result.Items) == 0 {
		logData.EmptyQuery = true
		return golib.AGW200("", nil), logData, nil
	}

	logrus.Debug(result.Items)

	json, err := golib.KeysResponseToJson(result, paginatedQuery.sortAscending, paginatedQuery.firstPage)
	if err != nil {
		return golib.AGW500(), logData, err
	}

	return golib.AGW200(json, map[string]string{
		"Content-type": "application/json",
	}), logData, nil
}

func paginatedQuery(input paginatedQueryInput) (paginatedQueryOutput, error) {
	var expressionAttributeValues = make(map[string]*dynamodb.AttributeValue)

	expressionAttributeValues[":v1"] = &dynamodb.AttributeValue{
		S: aws.String(fmt.Sprintf("LOG#KEY#S2#%v", input.dungeonID)),
	}

	// this is the default query aka no pagination
	resp := paginatedQueryOutput{
		queryInput: dynamodb.QueryInput{
			ExpressionAttributeValues: expressionAttributeValues,
			KeyConditionExpression:    aws.String("gsi1pk = :v1"),
			IndexName:                 aws.String("gsi1"),
			TableName:                 input.ddbTableName,
			ScanIndexForward:          aws.Bool(false),
			ReturnConsumedCapacity:    aws.String("TOTAL"),
			Limit:                     aws.Int64(20 + 1),
		},
		sortAscending: false,
		firstPage:     true,
	}

	// # becomes %23 inside the query parameter, needs to be transformed back
	next, err := url.QueryUnescape(input.request.QueryStringParameters["next"])
	if err != nil {
		return resp, err
	}
	prev, err := url.QueryUnescape(input.request.QueryStringParameters["prev"])
	if err != nil {
		return resp, err
	}
	resp.inputNextSK = next
	resp.inputPrevSK = prev

	if next != "" { // this is the next page
		resp.firstPage = false

		resp.queryInput.KeyConditionExpression = aws.String("gsi1pk = :v1 AND gsi1sk < :v2")

		expressionAttributeValues[":v2"] = &dynamodb.AttributeValue{
			S: aws.String(next),
		}
	} else if prev != "" { // this is the previous page, note the reversed ordering
		resp.firstPage = false
		resp.sortAscending = true
		resp.queryInput.ScanIndexForward = &resp.sortAscending

		resp.queryInput.KeyConditionExpression = aws.String("gsi1pk = :v1 AND gsi1sk > :v2")
		expressionAttributeValues[":v2"] = &dynamodb.AttributeValue{
			S: aws.String(prev),
		}

	}

	resp.queryInput.ExpressionAttributeValues = expressionAttributeValues

	return resp, nil
}

func checkInput(input string) (int, error) {
	if input == "" {
		return 0, fmt.Errorf("dungeon id parameter is empty")
	}

	dungeonID, err := strconv.Atoi(input)
	if err != nil {
		return 0, err
	}
	// IMPROVE: we could check it against a list of known dungeon ids

	return dungeonID, nil
}

func main() {
	golib.InitLogging()

	sess, err := session.NewSession()
	if err != nil {
		return
	}
	svc = dynamodb.New(sess)
	xray.AWS(svc.Client)

	lambda.Start(handler)
}
