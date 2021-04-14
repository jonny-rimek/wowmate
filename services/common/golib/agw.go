package golib

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// AGW404 returns a error agi gw v2 response
func AGW404() events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 404,
		Body:       http.StatusText(404),
	}
}

// AGW400 returns a error agi gw v2 response
func AGW400() events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 400,
		Body:       http.StatusText(400),
	}
}

// AGW500 returns a error agi gw v2 response
func AGW500() events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 500,
		Body:       http.StatusText(500),
	}
}

// AGW200 returns a agi gw v2 success response
func AGW200(body string, headers map[string]string) events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		Headers:    headers,
		StatusCode: 200,
		Body:       body,
	}
}
