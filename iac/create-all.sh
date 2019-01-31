#!/bin/bash

source ../env.sh

cd layer-base
./apply.sh
cd -

cd layer-bastion
./apply.sh
cd -

cd layer-kubernetes
./apply.sh
cd -

cd layer-data
./apply.sh
cd -

cd layer-services
./apply.sh
cd -

cd ../visualizer
./apply.sh
cd -

cd ../functions
./apply.sh
cd -
