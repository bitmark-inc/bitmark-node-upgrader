#!/usr/bin/env sh
## How to run bitmarkNodeUpgrader as a docker container
## Setup your nase mount directory (here is staging directory)
# -v $nodeDir/watcherlog:/var/log \
nodeDir=$HOME/bitmark-node-data
docker run -d --name bitmarkNodeUpgrader \
    -e DOCKER_HOST="unix:///var/run/docker.sock" \
    -e NODE_IMAGE="bitmark/bitmark-node" \
    -e NODE_NAME="bitmarkNode" \
    -e USER_NODE_BASE_DIR=$nodeDir \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v $nodeDir/data:/.config/bitmark-node/bitmarkd/bitmark/data \
    -v $nodeDir/data-test:/.config/bitmark-node/bitmarkd/testing/data \
    bitmark/bitmark-node-upgrader
