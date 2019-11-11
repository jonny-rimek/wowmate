#!/bin/bash

git add --all
git commit --verbose
git push

rm handler/size/main
go build -o handler/size/main handler/size/size.go
echo built size

rm handler/parquet/main
go build -o handler/parquet/main handler/parquet/parquet.go
echo built parquet

rm handler/athena/main
go build -o handler/athena/main handler/athena/athena.go
echo built athena

rm handler/check/main
go build -o handler/check/main handler/check/check.go
echo built check

rm handler/import/main
go build -o handler/import/main handler/import/import.go
echo built import

rm handler/import2/main
go build -o handler/import2/main handler/import2/import2.go
echo built import2


tsc
echo compiled typescript to javascript

cdk deploy --require-approval=never

aws s3 cp test.csv s3://wowmatestack-upload51c4d210-18wofa313p49y


