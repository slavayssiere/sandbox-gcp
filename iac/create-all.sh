#!/bin/bash

cd layer-base
./apply.sh
cd -

cd layer-bastion
terraform apply \
    --var "region=europe-west1" \
    -auto-approve
cd -

cd layer-kubernetes
terraform apply \
    --var "region=europe-west1" \
    -auto-approve
cd -