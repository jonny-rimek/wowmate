package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	resp := events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       "waddup, yo",
	}
	return resp, nil
}

func main() {
	lambda.Start(handler)
}
