module github.com/jonny-rimek/wowmate/services/upload/convert

go 1.15

replace github.com/wowmate/jonny-rimek/wowmate/services/upload/convert/normalize => ./normalize

require (
	github.com/aws/aws-lambda-go v1.19.1
	github.com/aws/aws-sdk-go v1.35.2
	github.com/wowmate/jonny-rimek/wowmate/services/upload/convert/normalize v0.0.0-00010101000000-000000000000
	golang.org/x/text v0.3.3 // indirect
)
