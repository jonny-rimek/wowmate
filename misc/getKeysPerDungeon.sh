#!/bin/bash

echo '{"pathParameters": {"dungeon_id": "2291"}}' | \
sam local invoke ApiGetKeysPerDungeonLambda073DE524 \
  --template cdk.out/wm-dev.template.json \
  --event - \
  --profile default \
  --env-vars=misc/env.json
