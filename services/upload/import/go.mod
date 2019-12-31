module github.com/jonny-rimek/wowmate/services/upload/import

go 1.13

require (
	github.com/aws/aws-lambda-go v1.13.3
	github.com/aws/aws-sdk-go v1.26.8
	github.com/jonny-rimek/wowmate/services/golib/ddb v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.4.2
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553 // indirect
)

replace github.com/jonny-rimek/wowmate/services/golib/ddb => ../../golib/ddb