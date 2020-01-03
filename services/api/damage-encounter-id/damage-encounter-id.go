package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"context"
	"github.com/jonny-rimek/wowmate/services/golib"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v1": {
				N: aws.String(req.PathParameters["encounter-id"]),
			},
		},
		KeyConditionExpression: aws.String("gsi1pk = :v1"),
		TableName:              aws.String(os.Getenv("DDB_NAME")),
		ReturnConsumedCapacity: aws.String("TOTAL"),
		IndexName:              aws.String("GSI1"),
		Limit:                  aws.Int64(50),
		ScanIndexForward:       aws.Bool(false),
	}

	rcu, apiGwRes, err := golib.DDBQuery(ctx, input)
	golib.CanonicalLog(map[string]interface{}{
		"encounter-id":     req.PathParameters["encounter-id"], 
		"rcu":              rcu,
		"http-status-code": apiGwRes.StatusCode,
	})
	return apiGwRes, err
}

func main() {
	lambda.Start(handler)
}