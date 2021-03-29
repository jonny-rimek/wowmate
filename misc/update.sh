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
  cd $wowmateDir || exit
done

echo "start presign update"
cd $presignDir || exit
npm update
echo "presign updated"

echo "start frontend update"
cd $frontendDir || exit
yarn upgrade
echo "frontend updated"