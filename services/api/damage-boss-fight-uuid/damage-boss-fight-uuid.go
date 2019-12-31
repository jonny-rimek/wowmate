package main

import (
	"github.com/Sirupsen/logrus"
	"fmt"
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
	"net/http"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"context"
)

//DamageSummary format of the dynamodb data
//FIXME: update for new ddb schema - add gsi3pk
type DamageSummary struct {
	BossFightUUID string `json:"pk"`
	Damage        int64  `json:"sk"`
	CasterName    string `json:"caster_name"`
	CasterID      string `json:"gsi2pk"`
	EncounterID   int    `json:"gsi1pk"`
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//TODO: log req.PathParameters
	bossFightUUID:= req.PathParameters["boss-fight-uuid"]
	ddbTableName := os.Getenv("DDB_NAME")
	svc := dynamodb.New(session.New())

	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v1": {
				S: aws.String(bossFightUUID),
			},
		},
		KeyConditionExpression: aws.String("gsi2pk = :v1"),
		TableName:              aws.String(ddbTableName),
		ReturnConsumedCapacity: aws.String("TOTAL"),
		IndexName:              aws.String("GSI2"),
	}

	result, err := svc.QueryWithContext(ctx, input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				return APIGwError(429, aerr)
			case dynamodb.ErrCodeResourceNotFoundException:
				return APIGwError(500, aerr)
			case dynamodb.ErrCodeInternalServerError:
				return APIGwError(500, aerr)
			default:
				return APIGwError(500, aerr)
			}
		} else {
			 return APIGwError(500, err)
		}
	}
	log.Printf("Consumed RCU: %f", *result.ConsumedCapacity.CapacityUnits)
	if len(result.Items) == 0 {
		logrus.Error("no records returned from DynamoDB")
		return APIGwError(404, nil)
	}
	summaries := []DamageSummary{}
	for _, item := range result.Items {
		s := DamageSummary{}
		err = dynamodbattribute.UnmarshalMap(item, &s)
		if err != nil {
			 err = fmt.Errorf("Failed to unmarshal Record, %v", err)
			 return APIGwError(500, err)
		}
		summaries = append(summaries, s)
	}

	js, err := json.Marshal(summaries)

	return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body: string(js),
		}, nil
}

//APIGwError is a helper function for APIGatewayProxyResponse
func APIGwError(status int, err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, err
}

func main() {
	lambda.Start(handler)
}