package main

import (
	"log"
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
				S: aws.String(req.PathParameters["caster-id"]),
			},
		},
		KeyConditionExpression: aws.String("gsi3pk = :v1"),
		TableName:              aws.String(os.Getenv("DDB_NAME")),
		ReturnConsumedCapacity: aws.String("TOTAL"),
		IndexName:              aws.String("GSI3"),
		Limit:                  aws.Int64(50),
	}

	rcu, apiGwRes, err := golib.DDBQuery(ctx, input)
	golib.CanonicalLog(map[string]interface{}{
		"caster-id":     req.PathParameters["caster-id"], 
		"rcu":              rcu,
		"http-status-code": apiGwRes.StatusCode,
	})
	return apiGwRes, err
}

func main() {
	log.Println("tet1")
	lambda.Start(handler)
}