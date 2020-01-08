module github.com/jonny-rimek/wowmate/services/upload/parquet

go 1.13

replace github.com/jonny-rimek/wowmate/services/golib => ../../golib

require (
	github.com/DataDog/zstd v1.4.4 // indirect
	github.com/apache/thrift v0.13.0 // indirect
	github.com/aws/aws-lambda-go v1.13.3
	github.com/aws/aws-sdk-go v1.26.8
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/golang/snappy v0.0.1 // indirect
	github.com/jonny-rimek/wowmate/services/golib v0.0.0-00010101000000-000000000000
	github.com/xitongsys/parquet-go v1.4.0
	github.com/xitongsys/parquet-go-source v0.0.0-20191104003508-ecfa341356a6
)
