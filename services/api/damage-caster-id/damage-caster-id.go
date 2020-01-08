package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jonny-rimek/wowmate/services/golib"
	"log"
	"os"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v1": {
				S: aws.String(req.PathParameters["caster-id"]),
			},
		},
		KeyConditionExpression: aws.String("sk = :v1"),
		TableName:              aws.String(os.Getenv("DDB_NAME")),
		ReturnConsumedCapacity: aws.String("TOTAL"),
		IndexName:              aws.String("GSI2"),
		Limit:                  aws.Int64(50),
		ScanIndexForward:       aws.Bool(false),
	}

	rcu, apiGwRes, err := golib.DDBQuery(ctx, input)
	golib.CanonicalLog(map[string]interface{}{
		"caster-id":        req.PathParameters["caster-id"],
		"rcu":              rcu,
		"http-status-code": apiGwRes.StatusCode,
	})
	return apiGwRes, err
}

func main() {
	log.Println("tet1")
	lambda.Start(handler)
}