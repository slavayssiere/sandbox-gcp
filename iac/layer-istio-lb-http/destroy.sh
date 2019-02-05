#!/bin/bash

HEATH_NAME="health-gke"
PORT_NAME="http-port"
PORT_NUM="31380"
BACKEND_NAME="my-http-backend-service"
URL_MAP_NAME="my-lb"

# listGroup=$(gcloud container clusters describe test-cluster --zone europe-west1 --format 'value(instanceGroupUrls)')
# export IFS=";"
# for groupUrl in $listGroup; do
#   group=$(echo $groupUrl | awk -F "/" '{print $11}')
#   zone=$(echo $groupUrl | awk -F "/" '{print $9}')

#   gcloud compute instance-groups set-named-ports $group \
#      --named-ports $PORT_NAME:$PORT_NUM \
#      --zone $zone
# done

# gcloud compute http-health-checks create $HEATH_NAME \
#     --port $PORT_NUM \
#      --request-path=/
