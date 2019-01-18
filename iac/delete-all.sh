#!/bin/bash

cd ../visualizer
echo "Destroy dataset..."
./destroy.sh
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
