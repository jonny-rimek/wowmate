module github.com/jonny-rimek/wowmate/services/upload/convert

go 1.16

replace github.com/wowmate/jonny-rimek/wowmate/services/upload/convert/normalize => ./normalize

replace github.com/jonny-rimek/wowmate/services/common/golib => ./../../common/golib

require (
	github.com/aws/aws-lambda-go v1.23.0
	github.com/aws/aws-sdk-go v1.38.28
	github.com/aws/aws-xray-sdk-go v1.3.0
	github.com/jonny-rimek/wowmate/services/common/golib v0.0.0-20210421002005-d547e9073201
	github.com/wowmate/jonny-rimek/wowmate/services/upload/convert/normalize v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.0.0-20210428140749-89ef3d95e781
)
