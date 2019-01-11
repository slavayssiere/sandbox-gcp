#!/bin/bash

source ../env.sh

cd layer-base
./apply.sh
cd -

cd layer-kubernetes
./apply.sh
cd -

cd layer-services
./apply.sh
cd -
