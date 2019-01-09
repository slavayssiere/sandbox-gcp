#!/bin/bash

REGION="europe-west1"

terraform apply --var "region=$REGION"

gcloud beta container clusters update test-cluster \
    --update-addons=Istio=ENABLED \
    --istio-config=auth=MTLS_STRICT \
    --region $REGION
