module github.com/jonny-rimek/wowmate/services/upload/insert/summary/player-damage-simple-dynamodb

go 1.16

replace github.com/jonny-rimek/wowmate/services/common/golib => ./../../../../common/golib

require (
	github.com/aws/aws-lambda-go v1.23.0
	github.com/aws/aws-sdk-go v1.38.21
	github.com/aws/aws-sdk-go-v2 v1.3.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.2.2
	github.com/aws/aws-sdk-go-v2/service/timestreamquery v1.1.5
	github.com/jonny-rimek/wowmate/services/common/golib v0.0.0-20210415101058-d1136c33be1a
)
