#!/bin/bash

rm api-service/damage-boss-fight-uuid/main
go build -ldflags "-s -w" -o api-service/damage-boss-fight-uuid/main api-service/damage-boss-fight-uuid/damage-boss-fight-uuid.go
echo built damage-boss-fight-uuid

rm api-service/damage-encounter-id/main
go build -ldflags "-s -w" -o api-service/damage-encounter-id/main api-service/damage-encounter-id/damage-encounter-id.go
echo built damage-encounter-id

rm upload-service/size/main
go build -ldflags "-s -w" -o upload-service/size/main upload-service/size/size.go
echo built size

cd upload-service/parquet
rm main
go build -ldflags "-s -w" -o main . #this is needed when there is more than one file in the dir
cd ../..
echo built parquet

rm upload-service/athena/main
go build -ldflags "-s -w" -o upload-service/athena/main upload-service/athena/athena.go
echo built athena

rm upload-service/check/main
go build -ldflags "-s -w" -o upload-service/check/main upload-service/check/check.go
echo built check

rm upload-service/import/main
go build -ldflags "-s -w" -o upload-service/import/main upload-service/import/import.go
echo built import

rm upload-service/import2/main
go build -ldflags "-s -w" -o upload-service/import2/main upload-service/import2/import2.go
echo built import2

cd frontend
npm run build
echo built frontend
cd ..

tsc
echo compiled typescript to javascript

cdk diff

cdk deploy --require-approval=never

# aws s3 cp WoWCombatLog.txt s3://wowmatestack-upload51c4d210-18wofa313p49y


