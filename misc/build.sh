#!/bin/bash

wowmateDir="/home/jonny/dev/wowmate/"

#TODO: measure time total and substeps
#TODO: get them dynamically with log list ./... or smth
goDirs=(
  "services/common/golib"
  "services/api/combatlogs/keys/_combatlog_uuid/player-damage-done/get"
  "services/api/combatlogs/keys/_dungeon_id/get"
  "services/api/combatlogs/keys/index/get"
  "services/upload/convert/normalize"
  "services/upload/convert"
  "services/upload/query-timestream/keys"
  "services/upload/query-timestream/player-damage-done"
  "services/upload/insert/dynamodb/keys"
  "services/upload/insert/dynamodb/player-damage-done"
  "services/upload/insert/timestream/keys"
)

frontendDir="services/frontend"

for i in "${goDirs[@]}"
do
  cd $wowmateDir || exit
  cd "$i" || exit
  pwd
  go mod tidy
  gofmt -w -s .
  go test .
  go build -ldflags='-s -w' -o main .
  echo "compiled go"
done

#if there is a 2nd argument skip cdk build step
if [ -z "$2" ]
then
  cd $wowmateDir || exit
  echo "start typescript compile"
  npm run tsc #compile cdk typescript
  echo "start cdk synth"
  npm run cdk synth >/dev/null #this suppresses the output because it just spits out endless cfn yaml
  echo "built cdk"
fi

#if there is an argument skip frontend build step
if [ -z "$1" ]
then
  cd $frontendDir || exit

  yarn build
  echo "frontend built"
fi
