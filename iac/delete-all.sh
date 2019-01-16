#!/bin/bash

cd ../visualizer
echo "Destroy dataset..."
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

gcloud dataflow jobs run delete-datastore \
    --gcs-location gs://dataflow-templates/latest/Datastore_to_Datastore_Delete \
    --parameters \
datastoreReadGqlQuery="SELECT * FROM userstats",\
datastoreReadProjectId="slavayssiere-sandbox",\
datastoreDeleteProjectId="slavayssiere-sandbox"
