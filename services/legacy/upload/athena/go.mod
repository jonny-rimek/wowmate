module github.com/jonny-rimek/wowmate/services/upload/athena

go 1.13

replace github.com/jonny-rimek/wowmate/services/golib => ./../../golib

require (
	github.com/aws/aws-lambda-go v1.15.0
	github.com/aws/aws-sdk-go v1.29.34
	github.com/jmespath/go-jmespath v0.3.0 // indirect
)
