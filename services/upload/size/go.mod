module github.com/jonny-rimek/wowmate/services/upload/size

go 1.13

replace github.com/jonny-rimek/wowmate/services/golib => ../../golib

require (
	github.com/aws/aws-lambda-go v1.13.3
	github.com/aws/aws-sdk-go v1.27.0
	github.com/jonny-rimek/wowmate/services/golib v0.0.0-00010101000000-000000000000
)
