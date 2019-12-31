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

//DamageSummary format of the dynamodb data
type DamageSummary struct {
	BossFightUUID string `json:"pk"`
	Damage        int64  `json:"sk"`
	CasterName    string `json:"caster_name"`
	CasterID      string `json:"gsi2pk"`
	EncounterID   int    `json:"gsi1pk"`
}

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
	}

	rcu, apiGwRes, err := golib.DDBQuery(ctx, input)
	golib.CanonicalLog(map[string]interface{}{
		"encounter-id":     req.PathParameters["encounter-id"], 
		"rcu":              rcu,
		"http-status-code": apiGwRes.StatusCode,
	})
	return apiGwRes, err

	// result, err := svc.QueryWithContext(ctx, input)
	// if err != nil {
	// 	if aerr, ok := err.(awserr.Error); ok {
	// 		switch aerr.Code() {
	// 		case dynamodb.ErrCodeProvisionedThroughputExceededException:
	// 			return APIGwError(429, aerr)
	// 		case dynamodb.ErrCodeResourceNotFoundException:
	// 			return APIGwError(404, aerr)
	// 		case dynamodb.ErrCodeInternalServerError:
	// 			return APIGwError(500, aerr)
	// 		default:
	// 			return APIGwError(500, aerr)
	// 		}
	// 	} else {
	// 		 return APIGwError(500, err)
	// 	}
	// }
	// log.Printf("Consumed RCU: %f", *result.ConsumedCapacity.CapacityUnits)

	// summaries := []DamageSummary{}
	// for _, item := range result.Items {
	// 	s := DamageSummary{}
	// 	err = dynamodbattribute.UnmarshalMap(item, &s)
	// 	if err != nil {
	// 		 err = fmt.Errorf("Failed to unmarshal Record, %v", err)
	// 		 return APIGwError(500, err)
	// 	}
	// 	summaries = append(summaries, s)
	// }

	// js, err := json.Marshal(summaries)

	// return events.APIGatewayProxyResponse{
	// 		StatusCode: http.StatusOK,
	// 		Body: string(js),
	// 	}, nil
}

//APIGwError is a helper function for APIGatewayProxyResponse
// func APIGwError(status int, err error) (events.APIGatewayProxyResponse, error) {
// 	return events.APIGatewayProxyResponse{
// 		StatusCode: status,
// 		Body:       http.StatusText(status),
// 	}, err
// }

func main() {
	lambda.Start(handler)
}