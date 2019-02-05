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

gsutil mb gs://assets.gcp-wescale.slavayssiere.fr

gsutil cp ../app-sse/src/templates/index.html gs://assets.gcp-wescale.slavayssiere.fr
gsutil cp ../app-sse/src/templates/twitter.png gs://assets.gcp-wescale.slavayssiere.fr
gsutil iam ch allUsers:objectViewer gs://assets.gcp-wescale.slavayssiere.fr
gsutil web set -m index.html -e 404.html gs://assets.gcp-wescale.slavayssiere.fr
