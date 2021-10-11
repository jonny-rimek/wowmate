module github.com/jonny-rimek/wowmate/services/upload/insert/summary/player-damage-simple-dynamodb

go 1.16

replace github.com/jonny-rimek/wowmate/services/common/golib => ./../../../../common/golib

require (
	github.com/aws/aws-lambda-go v1.27.0
	github.com/aws/aws-sdk-go v1.40.59
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.5.1
	github.com/aws/aws-sdk-go-v2/service/timestreamquery v1.4.1
	github.com/jonny-rimek/wowmate/services/common/golib v0.0.0-20210608175304-608f34f03462
)
