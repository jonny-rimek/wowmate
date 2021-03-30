#!/bin/bash

wowmateDir=$(pwd)

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
presignDir="services/upload/presign"

cd "$wowmateDir" || exit
pwd

npm update
echo "cdk updated"

for i in "${goDirs[@]}"
do
  cd "$i" || exit
  pwd
  go mod tidy
  gofmt -w -s .
  go get -u
  echo "go updated"
  cd "$wowmateDir" || exit
done

echo "start presign update"
cd $presignDir || exit
npm update
cd "$wowmateDir" || exit
echo "presign updated"

echo "start frontend update"
cd $frontendDir || exit
yarn upgrade
cd "$wowmateDir" || exit
echo "frontend updated"