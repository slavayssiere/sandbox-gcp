#!/bin/bash

PROJECT="cool-wharf-207907"

while IFS=';' read -r NAME_INJECTOR SIZE FREQUENCY
do
    echo "Delete ${NAME_INJECTOR}_topic and ${NAME_INJECTOR}_subcription"
    gcloud alpha pubsub topics delete "${NAME_INJECTOR}_topic"
    gcloud alpha pubsub subscriptions delete "${NAME_INJECTOR}_subcription"

done < ../iac/list_injector.csv

kubectl delete -f ./tpl/
rm -f ./tpl/*
