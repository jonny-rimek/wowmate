#!/bin/bash
rm upload-service/size/main
go build -o upload-service/size/main upload-service/size/size.go
echo built size

cd upload-service/parquet
rm main
go build -o main .
cd ../..
echo built parquet

rm upload-service/athena/main
go build -o upload-service/athena/main upload-service/athena/athena.go
echo built athena

rm upload-service/check/main
go build -o upload-service/check/main upload-service/check/check.go
echo built check

rm upload-service/import/main
go build -o upload-service/import/main upload-service/import/import.go
echo built import

rm upload-service/import2/main
go build -o upload-service/import2/main upload-service/import2/import2.go
echo built import2


tsc
echo compiled typescript to javascript

cdk diff

cdk deploy --require-approval=never

aws s3 cp WoWCombatLog.txt s3://wowmatestack-upload51c4d210-18wofa313p49y


