module github.com/jonny-rimek/wowmate/services/upload/convert

go 1.16

replace github.com/wowmate/jonny-rimek/wowmate/services/upload/convert/normalize => ./normalize

replace github.com/jonny-rimek/wowmate/services/common/golib => ./../../common/golib

require (
	github.com/andybalholm/brotli v1.0.3 // indirect
	github.com/aws/aws-lambda-go v1.27.0
	github.com/aws/aws-sdk-go v1.40.59
	github.com/aws/aws-xray-sdk-go v1.6.0
	github.com/jonny-rimek/wowmate/services/common/golib v0.0.0-20210608175304-608f34f03462
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/mitchellh/hashstructure/v2 v2.0.2
	github.com/valyala/fasthttp v1.30.0 // indirect
	github.com/wowmate/jonny-rimek/wowmate/services/upload/convert/normalize v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.0.0-20211008194852-3b03d305991f
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20211008145708-270636b82663 // indirect
	google.golang.org/grpc v1.41.0 // indirect
)
