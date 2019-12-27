package main

import (
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

//DamageSummary is the query output from Athena
type DamageSummary struct {
	BossFightUUID string `json:"pk"`
	Damage        int64  `json:"sk"`
	CasterName    string `json:"caster_name"`
	CasterID      string `json:"gsi2pk"`
	EncounterID   int    `json:"gsi1pk"`
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	bossFightUUID:= req.PathParameters["boss-fight-uuid"]
	ddbTableName := os.Getenv("DDB_NAME")
	svc := dynamodb.New(session.New())

	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v1": {
				S: aws.String(bossFightUUID),
			},
		},
		KeyConditionExpression: aws.String("pk = :v1"),
		TableName:              aws.String(ddbTableName),
		ReturnConsumedCapacity: aws.String("TOTAL"),
	}

	result, err := svc.QueryWithContext(ctx, input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				log.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				log.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				log.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}
		// return awsh.ServerError(err)
	}

	log.Printf("Consumed RCU: %f", *result.ConsumedCapacity.CapacityUnits)

	summaries := []DamageSummary{}

	for _, item := range result.Items {
		s := DamageSummary{}
		err = dynamodbattribute.UnmarshalMap(item, &s)
		if err != nil {
			log.Printf("Failed to unmarshal Record, %v", err)
			// return awsh.ServerError(err)
		}
		summaries = append(summaries, s)
	}

	js, err := json.Marshal(summaries)


	return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body: string(js),
		}, nil
}

func main() {
	lambda.Start(handler)
}