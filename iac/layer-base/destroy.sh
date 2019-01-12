#!/bin/bash

REGION="europe-west1"
SUBNET_EUROPE1_CIDR="192.168.0.0/20"
GCP_PROJECT="slavayssiere-sandbox"

gcloud config set project $GCP_PROJECT
gcloud config set compute/region $REGION

gcloud -q beta compute routers nats delete nat-$REGION --router=nat-$REGION
gcloud -q beta compute routers delete nat-$REGION
gcloud -q beta dns managed-zones delete private-dns-zone

terraform destroy \
    -auto-approve

gcloud -q iam service-accounts delete "sa-pubsub-publisher@$GCP_PROJECT.iam.gserviceaccount.com"
gcloud -q iam service-accounts delete "sa-pubsub-subscriber@$GCP_PROJECT.iam.gserviceaccount.com"
gcloud -q iam service-accounts delete "sa-pubsub-full@$GCP_PROJECT.iam.gserviceaccount.com"
gcloud -q iam service-accounts delete "sa-aggregator@$GCP_PROJECT.iam.gserviceaccount.com"
