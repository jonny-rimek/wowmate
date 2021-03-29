#!/bin/bash

echo '{"pathParameters": {"combatloguuid": "0d13c057-77a0-4660-9840-520449326b31"}}' | \
sam local invoke ApiCombatlogPlayerDamageAdvancedDBC10651 \
  --template cdk.out/wm.template.json \
  --event - \
  --profile default \
  --env-vars=misc/env.json