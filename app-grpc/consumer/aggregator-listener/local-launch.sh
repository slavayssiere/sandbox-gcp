#!/bin/bash

export PROJECT_ID="slavayssiere-sandbox"
export SUB_NAME="projects/slavayssiere-sandbox/subscriptions/messages-normalized-sub-datastore"
export SECRET_PATH="/Users/slavayssiere/Code/slavayssiere-sandbox-gcp/iac/sa-pubsub-datastore.json"
export REDIS_HOST="localhost"

docker run -it -p 8080:8080 \
    -e SUB_NAME="projects/slavayssiere-sandbox/subscriptions/messages-normalized-sub-datastore" \
    -e SECRET_PATH="/secret/sa-pubsub-datastore.json" \
    -v /Users/slavayssiere/Code/slavayssiere-sandbox-gcp/iac/sa-pubsub-datastore.json:/secret/sa-pubsub-datastore.json \
    --entrypoint sh \
    eu.gcr.io/slavayssiere-sandbox/consumer-datastore:0.0.3


docker run -d -p 8080:8080 \
    -e SUB_NAME="projects/slavayssiere-sandbox/subscriptions/messages-normalized-sub-datastore" \
    -e SECRET_PATH="/secret/sa-pubsub-datastore.json" \
    -v /Users/slavayssiere/Code/slavayssiere-sandbox-gcp/iac/sa-pubsub-datastore.json:/secret/sa-pubsub-datastore.json \
    eu.gcr.io/slavayssiere-sandbox/consumer-datastore:0.0.3
