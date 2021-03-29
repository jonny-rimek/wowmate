#!/bin/bash

sam local invoke SummaryInsertPlayerDamageAdvancedDynamodbLambdaDE5C23E4 \
  --template cdk.out/wm.template.json \
  --event misc/insertSummaryEvent.json \
  --profile default \
  --env-vars=misc/env.json
