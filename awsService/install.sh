#!bin/bash
########
# This script build the code and push the code to the lambda api which is bitmark-node-info
#####
GOOS=linux go build
rm *.zip
zip -r node-updater-ver.zip node-updater-ver main.go
aws lambda update-function-code --function-name bitmark-node-info --zip-file fileb://node-updater-ver.zip