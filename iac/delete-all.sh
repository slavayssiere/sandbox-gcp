#!/bin/bash

cd ../visualizer
./destroy.sh
cd -

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

cbt deleteinstance "test-instance"
