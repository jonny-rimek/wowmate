module github.com/jonny-rimek/wowmate/services/upload/insert/summary/player-damage-advanced-dynamodb

go 1.16

replace github.com/jonny-rimek/wowmate/services/common/golib => ./../../../../common/golib

require (
	github.com/aws/aws-lambda-go v1.24.0
	github.com/aws/aws-sdk-go v1.38.45
	github.com/aws/aws-xray-sdk-go v1.4.0
	github.com/jonny-rimek/wowmate/services/common/golib v0.0.0-20210520093338-ff37f7a209e1
	github.com/sirupsen/logrus v1.8.1
)
