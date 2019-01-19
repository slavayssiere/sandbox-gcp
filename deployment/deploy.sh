#!/bin/bash

kubectl apply -f namespaces-injectors.yaml
kubectl apply -f namespaces-normalizers.yaml
kubectl apply -f namespaces-consumers.yaml

################################ Injectors ################################
source ../env.sh
kubectl create secret generic twitter-secrets \
    --from-literal=CONSUMER_KEY=$CONSUMER_KEY \
    --from-literal=CONSUMER_SECRET=$CONSUMER_SECRET \
    --from-literal=ACCESS_TOKEN=$ACCESS_TOKEN \
    --from-literal=ACCESS_SECRET=$ACCESS_SECRET \
    -n injectors

kubectl create secret generic sa-pubsub-publisher \
    --from-file=../iac/sa-pubsub-publisher.json \
    -n injectors

################################ Normalizers ################################
kubectl create secret generic sa-pubsub-full \
    --from-file=../iac/sa-pubsub-full.json \
    -n normalizers

################################ Consumers ################################
kubectl create secret generic sa-pubsub-subscriber \
    --from-file=../iac/sa-pubsub-subscriber.json \
    -n consumers

kubectl create secret generic sa-pubsub-bigtable \
    --from-file=../iac/sa-pubsub-bigtable.json \
    -n consumers


kubectl create secret generic sa-pubsub-datastore \
    --from-file=../iac/sa-pubsub-datastore.json \
    -n consumers

################################ Aggregators ################################
kubectl create secret generic sa-aggregator \
    --from-file=../iac/sa-aggregator.json \
    -n aggregators

kubectl apply -f .
