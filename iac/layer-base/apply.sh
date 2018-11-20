#!/bin/bash

REGION="europe-west1"
SUBNET_EUROPE1_CIDR="192.168.0.0/20"

gcloud config set project slavayssiere-sandbox
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

# --nat-custom-subnet-ip-ranges=[SUBNET_1],[SUBNET_3]

# # network
# gcloud compute networks create net-0 \
#     --subnet-mode custom

# gcloud compute networks subnets create subnet-0 \
#     --network net-0 \
#     --region $REGION \
#     --range $SUBNET_EUROPE1_CIDR \
#     --secondary-range c0-pods=10.4.0.0/14,c0-services=10.0.32.0/20 \
#     --enable-private-ip-google-access

# gcloud compute firewall-rules create fw-0 \
#     --network net-0 \
#     --allow tcp:22

# # bastion

# gcloud compute instances create nat-test-1 \
#     --image-family debian-9 \
#     --image-project debian-cloud \
#     --network net-0 \
#     --subnet subnet-0 \
#     --zone $REGION-b

# # kubernetes

# gcloud container clusters create private-cluster-1 \
#     --zone $REGION-b \
#     --enable-ip-alias \
#     --network net-0 \
#     --subnetwork subnet-0 \
#     --cluster-secondary-range-name c0-pods \
#     --services-secondary-range-name c0-services \
#     --enable-private-nodes \
#     --master-ipv4-cidr 172.16.0.0/28 \
#     --enable-master-authorized-networks \
#     --master-authorized-networks $SUBNET_EUROPE1_CIDR \
#     --no-enable-basic-auth \
#     --no-issue-client-certificate \
#     --enable-private-endpoint

