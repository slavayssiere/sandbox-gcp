#!/bin/bash

list_groups=$(gcloud container clusters describe test-cluster --zone europe-west1 --format 'value(instanceGroupUrls)' | sed 's/instanceGroupManagers/instanceGroups/g')

echo $list_groups

terraform apply \
    --var "list_groups=$list_groups"