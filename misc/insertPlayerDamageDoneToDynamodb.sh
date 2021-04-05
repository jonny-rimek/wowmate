#!/bin/bash

sam local invoke InsertPlayerDamageDoneToDynamodbLambda659E0DA7 \
  --template cdk.out/wm-dev.template.json \
  --event misc/insertPlayerDamageDoneToDynamodbEvent.json \
  --profile default \
  --env-vars=misc/env.json
