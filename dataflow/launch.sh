# ./dataflow \
#     --topic messages-normalized \
#     --subscription messages-normalized-sub-dataproc \
#     --output counts

bucketname="dataflow-bigdata-test"

./dataflow --topic messages-normalized \
          --subscription messages-normalized-sub-dataproc \
          --output "gs://$bucketname/count" \
          --runner dataflow \
          --project "slavayssiere-sandbox" \
          --temp_location "gs://$bucketname/tmp/" \
          --staging_location "gs://$bucketname/binaries/" \
          --worker_harness_container_image=apache-docker-beam-snapshots-docker.bintray.io/beam/go:20180515