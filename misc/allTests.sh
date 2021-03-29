#!/bin/bash

./misc/damagePerLog.sh && \
  ./misc/insertAdvancedDamageDynamoDB.sh && \
  ./misc/topkeys.sh && \
  ./misc/topkeysPerDungeon.sh && \
  ./misc/convert.sh skipBuild && \
  echo "integration tests done"

