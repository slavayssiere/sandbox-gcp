#!/bin/bash

source ../../env.sh

IAP_DOMAIN="iap.gcp-wescale.slavayssiere.fr"

gcloud compute addresses delete istio-lb-http
# gcloud beta compute ssl-certificates create "gcp-wescale-iap-cert" --domains $IAP_DOMAIN
