#!/bin/bash

# ./dataflow \
#     --input-type "text" \ 
#     --input test.json \
#     --output "gs://dataflow-bigdata-test/count"

bucketname="dataflow-bigdata-test"

./dataflow --input-type "pubsub" \
          --topic messages-normalized \
          --subscription messages-normalized-sub-dataproc \
          --output "gs://$bucketname/count" \
          --runner dataflow \
          --project "slavayssiere-sandbox" \
          --temp_location "gs://$bucketname/tmp/" \
          --staging_location "gs://$bucketname/binaries/" \
          --worker_harness_container_image=eu.gcr.io/slavayssiere-sandbox/beam-go
          
# eu.gcr.io/slavayssiere-sandbox/beam-go:2.10.0