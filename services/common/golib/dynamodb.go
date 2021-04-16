package golib

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// DynamoDBPutItem writes a single item to dynamodb and always returns consumed capacity
func DynamoDBPutItem(ctx aws.Context, svc *dynamodb.DynamoDB, ddbTableName *string, record interface{}) (*dynamodb.PutItemOutput, error) {
	av, err := dynamodbattribute.MarshalMap(record)
	if err != nil {
		return nil, err
	}

	var response *dynamodb.PutItemOutput
	if os.Getenv("LOCAL") == "true" {
		response, err = svc.PutItem(&dynamodb.PutItemInput{
			Item:                   av,
			TableName:              ddbTableName,
			ReturnConsumedCapacity: aws.String("TOTAL"),
		})
	} else {
		response, err = svc.PutItemWithContext(ctx, &dynamodb.PutItemInput{
			Item:                   av,
			TableName:              ddbTableName,
			ReturnConsumedCapacity: aws.String("TOTAL"),
		})
	}
	if err != nil {
		return nil, fmt.Errorf("failed to put item to dynamodb: %v", err)
	}

	return response, nil
}

/*
func DynamoDBPutItem2(client *dynamodb2.Client, ddbTableName *string, record interface{}) (*dynamodb2.PutItemOutput, error) {
	av, err := attributevalue2.MarshalMap(record)
	if err != nil {
		return nil, err
	}

	response, err := client.PutItem(context.TODO(), &dynamodb2.PutItemInput{
		Item:                   av,
		TableName:              ddbTableName,
		ReturnConsumedCapacity: types2.ReturnConsumedCapacityTotal,
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}
*/

// DynamoDBQuery is a helper to simplify querying a dynamo db table
func DynamoDBQuery(ctx aws.Context, svc *dynamodb.DynamoDB, input dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	var err error
	var result *dynamodb.QueryOutput

	if os.Getenv("LOCAL") == "true" {
		result, err = svc.Query(&input)
	} else {
		// result, err = svc.QueryWithContext(ctx, &input)
	}
	if err != nil {
		return nil, err
	}

	return result, nil
}

/*
//DynamoDBQuery2 is a helper to simplify querying a dynamo db table
func DynamoDBQuery2(client *dynamodb2.Client, input dynamodb2.QueryInput) (*dynamodb2.QueryOutput, error) {

	result, err := client.Query(context.TODO(), &input)
	if err != nil {
		return nil, err
	}

	return result, nil
}
*/
