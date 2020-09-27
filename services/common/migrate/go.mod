module github.com/jonny-rimek/wowmate/services/common/migrate

go 1.15

replace github.com/jonny-rimek/wowmate/services/common/golib => ./../golib

require (
	github.com/aws/aws-lambda-go v1.19.1
	github.com/aws/aws-sdk-go v1.34.32
	github.com/golang-migrate/migrate/v4 v4.13.0
	github.com/jonny-rimek/wowmate/services/common/golib v0.0.0-00010101000000-000000000000
	github.com/lib/pq v1.8.0 // indirect
)
