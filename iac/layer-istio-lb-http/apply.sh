#!/bin/bash

source ../../env.sh

IAP_DOMAIN="iap.gcp-wescale.slavayssiere.fr"

gcloud compute addresses create istio-lb-http --global
gcloud beta compute ssl-certificates create "gcp-wescale-iap-cert" --domains $IAP_DOMAIN

cp policy.yaml policy-audience.yaml
sed -i.bak "s/IAP_AUDIENCE/$IAP_AUDIENCE/g" policy-audience.yaml

kubectl apply -f istio-ingress.yaml
kubectl apply -f policy-audience.yaml

rm policy-audience.yaml
