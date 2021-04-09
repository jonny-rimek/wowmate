#!/bin/bash

echo '{"pathParameters": {"combatlog_uuid": "06778555-45d4-404b-a79b-5f03ef2d6b37"}}' | \
sam local invoke ApiGetPlayerDamageDoneLambda7F052E98 \
  --template cdk.out/wm-dev.template.json \
  --event - \
  --profile default \
  --env-vars=misc/env.json