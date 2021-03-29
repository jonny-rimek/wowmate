#!/bin/bash

sam local invoke ApiDamageSummariesLambdaAB40B084 \
  --no-event \
  --template cdk.out/wm.template.json \
  --profile default \
  --env-vars=misc/env.json
