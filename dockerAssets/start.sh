#!/usr/bin/env sh
# script to run after docker is started
/go/bin/bitmark-node-upgrader --host=$DOCKER_HOST --image=$NODE_IMAGE --name=$NODE_NAME --verbose=true

# Run Example when not using docker
# export USER_NODE_BASE_DIR=$HOME/bitmark-node-data
# /go/bin/bitmark-node-upgrader --host="unix:///var/run/docker.sock" --image="bitmark/bitmark-node" --name="bitmarkNode" --verbose=true