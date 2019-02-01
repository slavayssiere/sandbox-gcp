#!/bin/bash

REGION="europe-west1"
MYIP=$(curl ifconfig.me)

terraform apply \
    --var "region=$REGION" \
    --var "myip=$MYIP" \
    -auto-approve

# gcloud beta container clusters update test-cluster \
#     --update-addons=Istio=ENABLED \
#     --istio-config=auth=MTLS_STRICT \
#     --region $REGION
