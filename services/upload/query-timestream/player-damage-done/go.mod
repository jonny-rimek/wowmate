module github.com/jonny-rimek/wowmate/services/upload/get-summary/player-damage-advanced

go 1.16

replace github.com/jonny-rimek/wowmate/services/common/golib => ./../../../common/golib

require (
	github.com/aws/aws-lambda-go v1.23.0
	github.com/aws/aws-sdk-go v1.38.1
	github.com/aws/aws-xray-sdk-go v1.3.0
	github.com/jonny-rimek/wowmate/services/common/golib v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.8.1
)
