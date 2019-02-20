#!/bin/bash

GCP_PROJECT="slavayssiere-sandbox"

gcloud -q functions delete laststat \
    --region europe-west1

gcloud -q functions delete getstat \
    --region europe-west1

gcloud -q projects remove-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-cloudfunction@$GCP_PROJECT.iam.gserviceaccount.com --role roles/datastore.owner  >> /dev/null

gcloud -q iam service-accounts delete "sa-cloudfunction@$GCP_PROJECT.iam.gserviceaccount.com"  >> /dev/null

