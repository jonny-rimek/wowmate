#!/bin/bash

sam local start-api \
  --template cdk.out/wm.template.json \
  --profile default \
  --env-vars=misc/env.json