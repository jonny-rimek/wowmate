module github.com/jonny-rimek/wowmate/services/upload/parquet

go 1.13

replace github.com/jonny-rimek/wowmate/services/golib => ./../../golib

require (
	github.com/apache/thrift v0.13.0 // indirect
	github.com/aws/aws-lambda-go v1.15.0
	github.com/aws/aws-sdk-go v1.29.34
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/golang/snappy v0.0.1 // indirect
	github.com/jonny-rimek/wowmate/services/golib v0.0.0-20200327084839-1b7797e0fe61
	github.com/klauspost/compress v1.10.3 // indirect
	github.com/xitongsys/parquet-go v1.5.1
	github.com/xitongsys/parquet-go-source v0.0.0-20200326031722-42b453e70c3b
)
