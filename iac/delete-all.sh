#!/bin/bash

cd layer-data
terraform destroy \
    -auto-approve
cd -

cd layer-kubernetes
terraform destroy \
    -auto-approve
cd -

cd layer-base
./destroy.sh
cd -

GCP_PROJECT="slavayssiere-sandbox"
cbt -project $PROJECT deleteinstance "test-instance"
