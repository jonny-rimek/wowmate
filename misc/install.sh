#!/bin/bash

wowmateDir="/home/jonny/dev/wowmate/"

goDirs=(
  "services/common/golib"
  "services/api/combatlogs/list/summaries/index/get"
  "services/api/combatlogs/list/summaries/_dungeon_id/get"
  "services/api/combatlogs/advanced-damage/_combatlog_uuid/get"
  "services/upload/convert/normalize"
  "services/upload/convert"
  "services/upload/get-summary/player-damage-simple"
  "services/upload/get-summary/player-damage-advanced"
  "services/upload/insert-summary/player-damage-simple-dynamodb"
  "services/upload/insert-summary/player-damage-advanced-dynamodb"
)

frontendDir="services/frontend"
presignDir="services/upload/presign"

cd $wowmateDir || exit
pwd

npm install
echo "cdk installed"

for i in "${goDirs[@]}"
do
  cd "$i" || exit
  pwd
  go get
  echo "go installed"
  cd $wowmateDir || exit
done

echo "start presign install"
cd $presignDir || exit
npm update
echo "presign installed"

echo "start frontend install"
cd $frontendDir || exit
yarn install
echo "frontend installed"