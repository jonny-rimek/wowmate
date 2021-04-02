#!/bin/bash

./misc/getPlayerDamageDone.sh && \
  ./misc/insertPlayerDamageDoneToDynamodb.sh && \
  ./misc/getKeys.sh && \
  ./misc/getKeysPerDungeon.sh && \
  ./misc/convert.sh skipBuild && \
  echo "integration tests done"

