#!/bin/bash

cd layer-kubernetes
terraform destroy
cd -

cd layer-bastion
terraform destroy
cd -

cd layer-base
./destroy.sh
cd -
