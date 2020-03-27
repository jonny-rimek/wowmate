module github.com/jonny-rimek/wowmate/services/upload/import

go 1.13

require (
	github.com/aws/aws-lambda-go v1.15.0
	github.com/aws/aws-sdk-go v1.29.34
	github.com/jonny-rimek/wowmate/services/golib v0.0.0-20200327084839-1b7797e0fe61
	github.com/sirupsen/logrus v1.5.0
)

replace github.com/jonny-rimek/wowmate/services/golib => ../../golib
