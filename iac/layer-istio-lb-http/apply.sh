#!/bin/bash

source ../../env.sh

IAP_DOMAIN="iap.gcp-wescale.slavayssiere.fr"

terraform apply -auto-approve

#gcloud compute addresses create istio-lb-http --global
# gcloud beta compute ssl-certificates create "gcp-wescale-iap-cert" --domains $IAP_DOMAIN

kubectl apply -f istio-ingress.yaml

sleep 20

hcistio=$(gcloud compute health-checks list | grep k8s | cut -d ' ' -f1)
prevhc=$(gcloud compute health-checks list | grep HTTP | grep -v k8s-be | cut -d ' ' -f1)
port=$(gcloud compute health-checks describe $prevhc --format="value(httpHealthCheck.port)")

gcloud compute health-checks update http $hcistio --request-path=/healthz --port=$port 

id_service=$(gcloud compute backend-services describe $hcistio --global --format="value(id)")

audience="$IAP_AUDIENCE$id_service"

cp policy.yaml policy-audience.yaml
sed -i.bak "s/IAP_AUDIENCE/$audience/g" policy-audience.yaml

# kubectl apply -f policy-audience.yaml

# rm policy-audience.yaml
rm policy-audience.yaml.bak

#ASSET_DOMAIN="test-asset-bucket"
#gcloud compute backend-buckets create backend-asset-bucket --gcs-bucket-name=$ASSET_DOMAIN

#urlmap=$(gcloud compute url-maps list | grep $hcistio | cut -d ' ' -f1)

# gcloud compute url-maps add-path-matcher $urlmap \
#    --default-backend-bucket backend-asset-bucket \
#    --path-matcher-name "asset-path" \
#    --backend-bucket-path-rules "/=backend-asset-bucket,/*=backend-asset-bucket"

# gcloud compute url-maps add-path-matcher $urlmap \
#    --default-service $hcistio \
#    --path-matcher-name "webservice-path" \
#    --backend-service-path-rules "/events=$hcistio,/aggregator=$hcistio,/aggregator/*=$hcistio" \
#    --new-hosts="iap.gcp-wescale.slavayssiere.fr"

# gcloud compute url-maps remove-path-matcher $urlmap --path-matcher-name="asset-path"
# gcloud compute url-maps remove-path-matcher $urlmap --path-matcher-name="webservice-path"

# gcloud compute url-maps add-path-matcher $urlmap \
#    --default-backend-bucket backend-asset-bucket \
#    --path-matcher-name "bucket-path" \
#    --new-hosts="iap.gcp-wescale.slavayssiere.fr"