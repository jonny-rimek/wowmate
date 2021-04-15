module github.com/jonny-rimek/wowmate/services/upload/insert/summary/player-damage-advanced-dynamodb

go 1.16

replace github.com/jonny-rimek/wowmate/services/common/golib => ./../../../../common/golib

require (
	github.com/aws/aws-lambda-go v1.23.0
	github.com/aws/aws-sdk-go v1.38.21
	github.com/aws/aws-xray-sdk-go v1.3.0
	github.com/jonny-rimek/wowmate/services/common/golib v0.0.0-20210415101058-d1136c33be1a
	github.com/sirupsen/logrus v1.8.1
)
