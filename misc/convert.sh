#!/bin/bash

clear && echo "building cdk and go"
# watcher cli can't handle executing 2 commands at once with && in between
# and during testing I don't want to rebuild for every test that's why I need to
# add an extra flag
if [ -z "$1" ]
then
  ./misc/build.sh skipFrontend >/dev/null
fi

sam local invoke ConvertLambda3540DCCB \
  --template cdk.out/wm-dev.template.json \
  --event misc/convertInput.json \
  --profile default \
  --env-vars=misc/env.json
