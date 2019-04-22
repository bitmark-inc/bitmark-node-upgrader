#!/usr/bin/env sh
########
# This script build the code and push latestDBVers.Json to S3
# You need to have your aws-cli profile to use aws-cli
aws s3 cp latestDBVers.json s3://bitmark-node-update
aws s3api put-object-acl --bucket bitmark-node-update --key latestDBVers.json --acl public-read