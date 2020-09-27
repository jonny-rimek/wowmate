module github.com/jonny-rimek/wowmate/services/upload/import

go 1.15

replace github.com/jonny-rimek/wowmate/services/common/golib => ./../../common/golib

require (
	github.com/aws/aws-lambda-go v1.19.1
	github.com/aws/aws-sdk-go v1.34.32
	github.com/jonny-rimek/wowmate/services/common/golib v0.0.0-00010101000000-000000000000
	github.com/lib/pq v1.8.0
)
