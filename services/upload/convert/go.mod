module github.com/jonny-rimek/wowmate/services/upload/convert

go 1.16

replace github.com/wowmate/jonny-rimek/wowmate/services/upload/convert/normalize => ./normalize

replace github.com/jonny-rimek/wowmate/services/common/golib => ./../../common/golib

require (
	github.com/aws/aws-lambda-go v1.24.0
	github.com/aws/aws-sdk-go v1.38.45
	github.com/aws/aws-xray-sdk-go v1.4.0
	github.com/jonny-rimek/wowmate/services/common/golib v0.0.0-20210520093338-ff37f7a209e1
	github.com/mitchellh/hashstructure/v2 v2.0.1
	github.com/wowmate/jonny-rimek/wowmate/services/upload/convert/normalize v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.0.0-20210521195947-fe42d452be8f
)
