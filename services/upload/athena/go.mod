module github.com/jonny-rimek/wowmate/services/upload/athena

go 1.13

replace github.com/jonny-rimek/wowmate/services/golib => ../../golib

require (
	github.com/aws/aws-lambda-go v1.13.3
	github.com/aws/aws-sdk-go v1.27.0
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553 // indirect
)