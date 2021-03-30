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

npm install
echo "cdk installed"

for i in "${goDirs[@]}"
do
  cd "$i" || exit
  pwd
  go get
  echo "go installed"
  cd "$wowmateDir" || exit
done

echo "start presign install"
cd $presignDir || exit
npm install
cd "$wowmateDir" || exit
echo "presign installed"

echo "start frontend install"
cd $frontendDir || exit
yarn install
cd "$wowmateDir" || exit
echo "frontend installed"