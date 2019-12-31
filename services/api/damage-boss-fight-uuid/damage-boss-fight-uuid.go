package main

import (
	"github.com/sirupsen/logrus"
	"fmt"
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
	"net/http"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"context"
	"github.com/jonny-rimek/wowmate/services/golib"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v1": {
				S: aws.String(req.PathParameters["boss-fight-uuid"]),
			},
		},
		KeyConditionExpression: aws.String("gsi2pk = :v1"),
		TableName:              aws.String(os.Getenv("DDB_NAME")),
		ReturnConsumedCapacity: aws.String("TOTAL"),
		IndexName:              aws.String("GSI2"),
	}

	rcu, apiGwRes, err := ddbQuery(ctx, input)
	CanonicalLog(logrus.Fields{
		"boss-fight-uuid":  req.PathParameters["boss-fight-uuid"], 
		"rcu":              rcu,
		"http-status-code": apiGwRes.StatusCode,
	})
	return apiGwRes, err
}

//CanonicalLog IMPROVE:_
func CanonicalLog(msg map[string]interface{}){
	logrus.WithFields(msg).Info()
}

func ddbQuery(ctx context.Context, queryInput *dynamodb.QueryInput) (float64, events.APIGatewayProxyResponse, error) {
	svc := dynamodb.New(session.New())
	result, err := svc.QueryWithContext(ctx, queryInput)
	rcu := *result.ConsumedCapacity.CapacityUnits
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				return rcu, APIGwError(429), err
			case dynamodb.ErrCodeResourceNotFoundException:
				return rcu, APIGwError(500), err
			case dynamodb.ErrCodeInternalServerError:
				return rcu, APIGwError(500), err
			default:
				return rcu, APIGwError(500), err
			}
		} else {
			return rcu, APIGwError(500), err
		}
	}

	if len(result.Items) == 0 {
		logrus.Error("no records returned from DynamoDB")
		return rcu, APIGwError(404), nil
	}

	summaries := []golib.DamageSummary{}
	for _, item := range result.Items {
		s := golib.DamageSummary{}
		err = dynamodbattribute.UnmarshalMap(item, &s)
		if err != nil {
			 err = fmt.Errorf("Failed to unmarshal Record, %v", err)
			 return rcu, APIGwError(500), err
		}
		summaries = append(summaries, s)
	}

	js, err := json.Marshal(summaries)
	if err != nil {
		return rcu, APIGwError(500), fmt.Errorf("failed to marshal data to JSON: %v", err)
	}

	return rcu, APIGwOK(js), nil
}

//APIGwOK TODO: add comments
func APIGwOK(body []byte) (events.APIGatewayProxyResponse) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body: string(body),
	}
}

//APIGwError is a helper function for APIGatewayProxyResponse
func APIGwError(status int) (events.APIGatewayProxyResponse) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}
}

func main() {
	lambda.Start(handler)
}