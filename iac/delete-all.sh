#!/bin/bash

cd layer-kubernetes
terraform destroy \
    -auto-approve
cd -

cd layer-bastion
terraform destroy \
    -auto-approve
cd -

cd layer-base
./destroy.sh
cd -
