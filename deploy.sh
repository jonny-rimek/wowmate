#!/bin/bash
cd ~/dev/wowmate/services/api/damage-boss-fight-uuid
rm main
go mod tidy
go build -ldflags "-s -w" -o main . 
echo built damage-boss-fight-uuid

cd ../damage-encounter-id
rm main
go mod tidy
go build -ldflags "-s -w" -o main . 
echo built damage-encounter-id

cd ../damage-caster-id
rm main
go mod tidy
go build -ldflags "-s -w" -o main . 
echo built damage-caster-id

cd ../../upload/size
rm main
#TODO: go mod tidy 
go build -ldflags "-s -w" -o main . 
echo built size

cd ../parquet
rm main
go mod tidy
go build -ldflags "-s -w" -o main . 
echo built parquet

cd ../athena
rm main
#TODO: go mod tidy
go build -ldflags "-s -w" -o main . 
echo built athena

cd ../check
rm main
#TODO: go mod tidy
go build -ldflags "-s -w" -o main . 
echo built check

cd ../import
rm main
go mod tidy
go build -ldflags "-s -w" -o main .
echo built import

cd ../../frontend
npm run build
echo built frontend
cd ../..

tsc
echo compiled typescript to javascript

cdk diff

cdk deploy --require-approval=never

# aws s3 cp WoWCombatLog.txt s3://wowmatestack-upload51c4d210-18wofa313p49y


