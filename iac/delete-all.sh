#!/bin/bash

cd ../visualizer
echo "Destroy dataset..."
./destroy.sh
cd -

cd layer-bastion
terraform destroy \
    -auto-approve
cd -

cd layer-data
./destroy.sh
cd -

cd layer-kubernetes
terraform destroy \
    -auto-approve
cd -

cd layer-base
./destroy.sh
cd -

cd ../functions
./destroy.sh
cd -


ASSET_DOMAIN="assets.gcp-wescale.slavayssiere.fr"

gsutil rm -r $ASSET_DOMAIN

gcloud beta compute ssl-certificates delete "gcp-wescale-asset-cert"


