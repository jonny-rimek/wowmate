module github.com/jonny-rimek/wowmate/services/upload/get-summary/player-damage-simple

go 1.16

replace github.com/jonny-rimek/wowmate/services/common/golib => ./../../../common/golib

require (
	github.com/aws/aws-lambda-go v1.23.0
	github.com/aws/aws-sdk-go v1.38.34
	github.com/aws/aws-xray-sdk-go v1.4.0
	github.com/jonny-rimek/wowmate/services/common/golib v0.0.0-20210505192330-927229c4ea83
	github.com/sirupsen/logrus v1.8.1
)
