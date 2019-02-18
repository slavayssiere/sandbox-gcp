#!/bin/bash


bucketname="dataflow-bigdata-test"

# python -m apache_beam.examples.wordcount \
#           --input_subscription messages-normalized-sub-dataproc \
#           --output "gs://$bucketname/counts" \
#           --runner DataflowRunner \
#           --project "slavayssiere-sandbox" \
#           --temp_location "gs://$bucketname/tmp/""


python main.py \
          --input_subscription "projects/slavayssiere-sandbox/subscriptions/messages-normalized-sub-dataproc" \
          --output "gs://$bucketname/count" \
          --runner DataflowRunner \
          --project "slavayssiere-sandbox" \
          --temp_location "gs://$bucketname/tmp/" \
          --experiments=allow_non_updatable_job