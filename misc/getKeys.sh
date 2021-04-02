#!/bin/bash

sam local invoke ApiGetKeysLambdaFDF1A526 \
  --no-event \
  --template cdk.out/wm.template.json \
  --profile default \
  --env-vars=misc/env.json
