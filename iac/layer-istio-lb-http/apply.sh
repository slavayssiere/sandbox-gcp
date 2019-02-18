#!/bin/bash

source ../../env.sh

terraform apply -auto-approve

kubectl apply -f istio-ingress.yaml

until (gcloud compute health-checks list | grep k8s)
do
    echo "Wait for healthcheck..."
    sleep 10
done

hcs=$(gcloud compute health-checks list --protocol=HTTP | grep k8s | cut -d ' ' -f1)
porthc=$(kubectl -n istio-system get svc admin-ingressgateway -o jsonpath='{.spec.healthCheckNodePort}')

while read -r line; do   
    echo "change health for '$line' on port: $porthc "
    gcloud compute health-checks update http $line --request-path=/healthz --port=$porthc
    # id_service=$(gcloud compute backend-services describe $hcistio --global --format="value(id)")
done <<< "$hcs"


# audience="$IAP_AUDIENCE$id_service"

# cp policy.yaml policy-audience.yaml
# sed -i.bak "s/IAP_AUDIENCE/$audience/g" policy-audience.yaml

# # kubectl apply -f policy-audience.yaml

# # rm policy-audience.yaml
# rm policy-audience.yaml.bak

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