#!/bin/bash

# watcher cli can't handle executing 2 cmds at once with && inbetween
# and during testing I don't want to rebuild for every test that's why I need to
# add an extra flag
if [ -z "$1" ]
then
  ./misc/build.sh skipFrontend
fi

sam local invoke ConvertLambda3540DCCB \
  --template cdk.out/wm.template.json \
  --event misc/convertInput.json \
  --profile default \
  --env-vars=misc/env.json