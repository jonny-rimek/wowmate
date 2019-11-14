#!/bin/bash

git add --all
git commit --verbose
git push

rm lambda/size/main
go build -o lambda/size/main lambda/size/size.go
echo built size

cd lambda/parquet
rm main
go build -o main parquet.go
cd ../..
echo built parquet

rm lambda/athena/main
go build -o lambda/athena/main lambda/athena/athena.go
echo built athena

rm lambda/check/main
go build -o lambda/check/main lambda/check/check.go
echo built check

rm lambda/import/main
go build -o lambda/import/main lambda/import/import.go
echo built import

rm lambda/import2/main
go build -o lambda/import2/main lambda/import2/import2.go
echo built import2


tsc
echo compiled typescript to javascript

cdk diff | less

cdk deploy --require-approval=never

aws s3 cp test.csv s3://wowmatestack-upload51c4d210-18wofa313p49y


