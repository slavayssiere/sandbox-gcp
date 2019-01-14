#!/bin/bash

export PROJECT_ID="slavayssiere-sandbox"
export TABLE_ID="test-table"
export INSTANCE_ID="test-instance"
export SUB_NAME="projects/slavayssiere-sandbox/subscriptions/messages-normalized-sub-bigtable"
export SECRET_PATH="/Users/slavayssiere/Code/slavayssiere-sandbox-gcp/iac/sa-pubsub-bigtable.json"


docker run -it -p 8080:8080 \
    -e SUB_NAME="projects/slavayssiere-sandbox/subscriptions/messages-normalized-sub-bigtable" \
    -e SECRET_PATH="/secret/sa-pubsub-bigtable.json" \
    -v /Users/slavayssiere/Code/slavayssiere-sandbox-gcp/iac/sa-pubsub-bigtable.json:/secret/sa-pubsub-bigtable.json \
    --entrypoint sh \
    eu.gcr.io/slavayssiere-sandbox/consumer-bigtable:0.0.3


docker run -d -p 8080:8080 \
    -e SUB_NAME="projects/slavayssiere-sandbox/subscriptions/messages-normalized-sub-bigtable" \
    -e SECRET_PATH="/secret/sa-pubsub-bigtable.json" \
    -v /Users/slavayssiere/Code/slavayssiere-sandbox-gcp/iac/sa-pubsub-bigtable.json:/secret/sa-pubsub-bigtable.json \
    eu.gcr.io/slavayssiere-sandbox/consumer-bigtable:0.0.3
