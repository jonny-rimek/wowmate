#!/bin/bash
cd services

rm api/damage-boss-fight-uuid/main
go build -ldflags "-s -w" -o api/damage-boss-fight-uuid/main api/damage-boss-fight-uuid/damage-boss-fight-uuid.go
echo built damage-boss-fight-uuid

rm api/damage-encounter-id/main
go build -ldflags "-s -w" -o api/damage-encounter-id/main api/damage-encounter-id/damage-encounter-id.go
echo built damage-encounter-id

rm upload/size/main
go build -ldflags "-s -w" -o upload/size/main upload/size/size.go
echo built size

cd upload/parquet
rm main
go build -ldflags "-s -w" -o main . #this is needed when there is more than one file in the dir
cd ../..
echo built parquet

rm upload/athena/main
go build -ldflags "-s -w" -o upload/athena/main upload/athena/athena.go
echo built athena

rm upload/check/main
go build -ldflags "-s -w" -o upload/check/main upload/check/check.go
echo built check

rm upload/import/main
go build -ldflags "-s -w" -o upload/import/main upload/import/import.go
echo built import

rm upload/import2/main
go build -ldflags "-s -w" -o upload/import2/main upload/import2/import2.go
echo built import2

cd frontend
npm run build
echo built frontend
cd ../..

tsc
echo compiled typescript to javascript

cdk diff

cdk deploy --require-approval=never

# aws s3 cp WoWCombatLog.txt s3://wowmatestack-upload51c4d210-18wofa313p49y


