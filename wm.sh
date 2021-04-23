#!/bin/bash

wowmateDir=$(pwd)
presignDir="services/upload/presign"
frontendDir="services/frontend"
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

local-api() {
  build_go
  build_cdk
  cdk_synth

  sam local start-api \
    --template cdk.out/wm.template.json \
    --profile default \
    --env-vars=misc/env.json \
    || exit 1
}

get-keys-per-dungeon() {
  build_go
  build_cdk
  cdk_synth

  echo '{"pathParameters": {"dungeon_id": "2291"}}' | \
  sam local invoke ApiGetKeysPerDungeonLambda073DE524 \
    --template cdk.out/wm-dev.template.json \
    --event - \
    --profile default \
    --env-vars=misc/env.json \
    || exit 1
}

insert-keys-to-dynamodb() {
  build_go
  build_cdk
  cdk_synth

  sam local invoke InsertKeysToDynamodbLambda15825024 \
    --template cdk.out/wm-dev.template.json \
    --event misc/insertPlayerDamageDoneToDynamodbEvent.json \
    --profile default \
    --env-vars=misc/env.json \
    || exit 1
}
get-keys() {
  build_go
  build_cdk
  cdk_synth

  sam local invoke ApiGetKeysLambdaFDF1A526 \
    --no-event \
    --template cdk.out/wm-dev.template.json \
    --profile default \
    --env-vars=misc/env.json \
    || exit 1
}

get-player-damage-done() {
  build_go
  build_cdk
  cdk_synth

  echo '{"pathParameters": {"combatlog_uuid": "06778555-45d4-404b-a79b-5f03ef2d6b37"}}' | \
  sam local invoke ApiGetPlayerDamageDoneLambda7F052E98 \
    --template cdk.out/wm-dev.template.json \
    --event - \
    --profile default \
    --env-vars=misc/env.json \
    || exit 1
}

insert-player-damage-done() {
  build_go
  build_cdk
  cdk_synth

  sam local invoke InsertPlayerDamageDoneToDynamodbLambda659E0DA7 \
    --template cdk.out/wm-dev.template.json \
    --event misc/insertPlayerDamageDoneToDynamodbEvent.json \
    --profile default \
    --env-vars=misc/env.json \
    || exit 1
}

query-player-damage-done() {
  build_go
  build_cdk
  cdk_synth

  sam local invoke QueryPlayerDamageDoneLambda98AFC037 \
    --template cdk.out/wm-dev.template.json \
    --event misc/queryTimestreamInput.json \
    --profile default \
    --env-vars=misc/env.json \
    || exit 1
}

convert() {
  build_go
  build_cdk
  cdk_synth

  sam local invoke ConvertLambda3540DCCB \
    --template cdk.out/wm-dev.template.json \
    --event misc/convertInput.json \
    --profile default \
    --env-vars=misc/env.json \
    || exit 1
}

cdk_synth() {
  npx cdk synth "wm-dev" -q || exit 1
}

# comments the diff on the current PR
# source https://github.com/youyo/aws-cdk-github-actions/blob/master/entrypoint.sh#L63
# won't work locally, but it's not supposed to be.
cdk_diff_prod() {
	output=$(cdk_diff "wm" 2>&1)
	exitCode=${?}
	echo ::set-output name=status_code::${exitCode}
	echo "${output}"

	commentStatus="Failed"
	if [ "${exitCode}" == "0" ]; then
		commentStatus="Success"
	elif [ "${exitCode}" != "0" ]; then
		echo "CDK diff for stack wm has failed. See above console output for more details."
		exit 1
	fi

  commentWrapper="
\`\`\`
${output}
\`\`\`
"

  payload=$(echo "${commentWrapper}" | jq -R --slurp '{body: .}')
  commentsURL=$(cat "${GITHUB_EVENT_PATH}" | jq -r .pull_request.comments_url)

  echo "${payload}" | curl -s -S -H "Authorization: token ${GITHUB_TOKEN}" --header "Content-Type: application/json" --data @- "${commentsURL}"  > /dev/null
}

cdk_deploy() {
 npm run cdk deploy -- --require-approval=never "$1" || exit 1
}

cdk_diff() {
 npx cdk diff "$1" || exit 1
}

update_go() {
  for i in "${goDirs[@]}"
  do
    cd "$wowmateDir" || exit 1
    cd "$i" || exit 1
    go get -u || exit 1
    go mod tidy || exit 1
  done
echo "go updated"
}

update_cdk() {
  cd "$wowmateDir" || exit 1
  npm update || exit 1
  echo "cdk updated"
}

update_frontend() {
#  echo "start frontend update"
  cd "$wowmateDir" || exit 1
  cd $frontendDir || exit 1
  yarn upgrade || exit 1
  echo "frontend updated"
}

update_presign() {
#  echo "start presign update"
  cd "$wowmateDir" || exit 1
  cd $presignDir || exit 1
  npm update || exit 1
  echo "presign updated"
}

install_go() {
  for i in "${goDirs[@]}"
  do
    cd "$wowmateDir" || exit 1
    cd "$i" || exit 1
    go get || exit 1
  done
  echo "go installed"
}

install_cdk() {
  cd "$wowmateDir" || exit 1
  npm ci || exit 1
  echo "cdk installed"
}

install_frontend() {
#  echo "start frontend install"
  cd "$wowmateDir" || exit 1
  cd "$frontendDir" || exit 1
  yarn install || exit 1
  echo "frontend installed"
}

install_presign() {
#  echo "start presign install"
  cd "$wowmateDir" || exit 1
  cd "$presignDir" || exit 1
  npm ci || exit 1
  echo "presign installed"
}

lint_go() {
  go get -u github.com/kisielk/errcheck || exit 1
  go get -u honnef.co/go/tools/cmd/staticcheck || exit 1

  for i in "${goDirs[@]}"
  do
    cd "$wowmateDir" || exit 1
    cd "$i" || exit 1
    diff -u <(echo -n) <(gofmt -d .) || exit 1
    go vet || exit 1
    errcheck || exit 1
#    TODO: fix all problems
    staticcheck # || exit 1
    go test . || exit 1
  done
  echo "compiled go"
}

build_go() {
  for i in "${goDirs[@]}"
  do
    cd "$wowmateDir" || exit 1
    cd "$i" || exit 1
    # build all binaries into the dist folder
    path=$(pwd)
    go build -ldflags='-s -w' -o "${path/services/dist}"/main . || exit 1
  done
  echo "compiled go"
}

build_cdk() {
  cd "$wowmateDir" || exit 1
  echo "start typescript compile"
  npm run tsc || exit 1 #compile cdk typescript
}

build_frontend() {
  cd "$wowmateDir" || exit 1
  cd "$frontendDir" || exit 1

  yarn build || exit 1
  echo "frontend built"
}

main() {
  if [ "$1" == "build" ]
  then
    if [ "$2" == "frontend" ]
    then
      build_frontend
    elif [ "$2" == "go" ]
    then
      build_go
    elif [ "$2" == "cdk" ]
    then
      build_cdk
    else
      build_cdk
      build_go
      build_frontend
    fi
  elif [ "$1" == "install" ]
  then
    if [ "$2" == "frontend" ]
    then
      install_frontend
    elif [ "$2" == "presign" ]
    then
      install_presign
    elif [ "$2" == "go" ]
    then
      install_go
    elif [ "$2" == "cdk" ]
    then
      install_cdk
    else
      install_cdk
      install_go
      install_frontend
    fi
  elif [ "$1" == "update" ]
  then
    if [ "$2" == "frontend" ]
    then
      update_frontend
    elif [ "$2" == "presign" ]
    then
      update_presign
    elif [ "$2" == "go" ]
    then
      update_go
    elif [ "$2" == "cdk" ]
    then
      update_cdk
    else
      update_cdk
      update_go
      update_frontend
    fi
  elif [ "$1" == "deploy" ]
  then
    if [ "$2" == "prod" ]
    then
      cdk_diff "wm"
      cdk_deploy "wm"
    elif [ "$2" == "dev" ]
    then
      cdk_diff "wm-dev"
      cdk_deploy "wm-dev"
    fi
  elif [ "$1" == "diff" ]
  then
    if [ "$2" == "prod" ]
    then
      cdk_diff_prod
    elif [ "$2" == "dev" ]
    then
      cdk_diff "wm-dev"
    elif [ "$2" == "preprod" ]
    then
      cdk_diff "wm-preprod"
    fi
  elif [ "$1" == "synth" ]
  then
    cdk_synth
  fi
  if [ "$1" == "lint" ]
  then
    if [ "$2" == "go" ]
    then
      lint_go
    fi
  # invoking lambdas locally
  elif [ "$1" == "invoke" ]
  then
    if [ "$2" == "convert" ]
    then
      convert
    elif [ "$2" == "query-player-damage-done" ]
    then
      query-player-damage-done
    elif [ "$2" == "get-keys-per-dungeon" ]
    then
      get-keys-per-dungeon
    elif [ "$2" == "local-api" ]
    then
      local-api
    elif [ "$2" == "insert-keys-to-dynamodb" ]
    then
      insert-keys-to-dynamodb
    elif [ "$2" == "get-keys" ]
    then
      get-keys
    elif [ "$2" == "get-player-damage-done" ]
    then
      get-player-damage-done
    elif [ "$2" == "insert-player-damage-done" ]
    then
      insert-player-damage-done
    fi
  fi
}

main "$@" || exit 1
