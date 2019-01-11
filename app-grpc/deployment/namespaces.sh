#!/bin/bash

GCP_PROJECT="slavayssiere-sandbox"

################################ Injectors ################################
kubectl create ns injectors
kubectl create secret generic twitter-secrets \
    --from-literal=CONSUMER_KEY=$CONSUMER_KEY \
    --from-literal=CONSUMER_SECRET=$CONSUMER_SECRET \
    --from-literal=ACCESS_TOKEN=$ACCESS_TOKEN \
    --from-literal=ACCESS_SECRET=$ACCESS_SECRET \
    -n injectors

kubectl create secret generic sa-pubsub-publisher \
    --from-file=../../iac/sa-pubsub-publisher.json \
    -n injectors

kubectl apply -f injector-twitter.yaml

################################ Normalizers ################################
kubectl create ns normalizers

kubectl create secret generic sa-pubsub-full \
    --from-file=../../iac/sa-pubsub-full.json \
    -n normalizers

################################ Consumers ################################
kubectl create ns consumers

kubectl create secret generic sa-pubsub-subscriber \
    --from-file=../../iac/sa-pubsub-subscriber.json \
    -n consumers
