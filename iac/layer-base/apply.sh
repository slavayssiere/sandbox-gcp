#!/bin/bash

REGION="europe-west1"
SUBNET_EUROPE1_CIDR="192.168.0.0/20"
GCP_PROJECT="slavayssiere-sandbox"

gcloud config set project $GCP_PROJECT
gcloud config set compute/region $REGION

terraform apply \
    --var "region=europe-west1" \
    -auto-approve

gcloud -q beta compute routers create nat-$REGION  \
    --network demo-net \
    --region $REGION

gcloud -q beta compute routers nats create nat-$REGION \
    --router-region $REGION \
    --router nat-$REGION \
    --nat-all-subnet-ip-ranges \
    --auto-allocate-nat-external-ips

gcloud -q beta dns managed-zones create private-dns-zone \
    --dns-name="gcp.wescale" \
    --description="A private zone" \
    --visibility=private \
    --networks=demo-net


gcloud iam service-accounts create "sa-pubsub-publisher"
    --display-name "SA for pubsub publish apps"

gcloud iam service-accounts create "sa-pubsub-subscriber"
    --display-name "SA for pubsub publish apps"

gcloud iam service-accounts create "sa-pubsub-full"
    --display-name "SA for pubsub publish apps"

gcloud projects add-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-pubsub-publisher@$GCP_PROJECT.iam.gserviceaccount.com --role roles/pubsub.publisher
gcloud projects add-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-pubsub-subscriber@$GCP_PROJECT.iam.gserviceaccount.com --role roles/pubsub.subscriber
gcloud projects add-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-pubsub-full@$GCP_PROJECT.iam.gserviceaccount.com --role roles/pubsub.publisher
gcloud projects add-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-pubsub-full@$GCP_PROJECT.iam.gserviceaccount.com --role roles/pubsub.subscriber

gcloud iam service-accounts keys create ../sa-pubsub-publisher.json \
  --iam-account sa-pubsub-publisher@$GCP_PROJECT.iam.gserviceaccount.com

gcloud iam service-accounts keys create ../sa-pubsub-subscriber.json \
  --iam-account sa-pubsub-subscriber@$GCP_PROJECT.iam.gserviceaccount.com

gcloud iam service-accounts keys create ../sa-pubsub-full.json \
  --iam-account sa-pubsub-full@$GCP_PROJECT.iam.gserviceaccount.com

