# README for bitmark-node-updater

## Purpose

The purpose of bitmark-node-updater is to

+ update the bitmark-node docker image 

+ bitmarkd database

## Architecture

![BitmarkNodeUpdater](https://i.imgur.com/iWNNhWf.jpg)

## Information on project source tree

### bitmark-node-updater/awsService/data

+ latestDBVers: Json file which records information of latest database

+ uploads3.sh: auto upload json file to S3

### bitmark-node-updater/awsService

+ install.sh: install aws lambda code

+ go code: code for lambda in go

### bitmark-node-updater

+ go code: updater service code

+ Dockerfile: build docker image

