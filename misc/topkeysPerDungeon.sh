#!/bin/bash

echo '{"pathParameters": {"dungeon_id": "2291"}}' | \
sam local invoke ApiDamageDungeonSummariesLambda00443886 \
  --template cdk.out/wm.template.json \
  --event - \
  --profile default \
  --env-vars=misc/env.json
