#!/usr/bin/env sh
# script to run after docker is started
/go/bin/bitmark-node-updater --host=$DOCKER_HOST --image=$NODE_IMAGE --name=$NODE_NAME --verbose=true