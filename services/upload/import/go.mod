module github.com/jonny-rimek/wowmate/services/upload/import

go 1.13

require (
	github.com/aws/aws-lambda-go v1.13.3
	github.com/aws/aws-sdk-go v1.26.8
	github.com/jonny-rimek/wowmate/services/golib v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.4.2
)

replace github.com/jonny-rimek/wowmate/services/golib => ../../golib
