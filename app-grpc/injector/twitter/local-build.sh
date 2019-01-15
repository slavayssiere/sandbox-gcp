# Localy do:

version=$1

export GCP_PROJECT=$(gcloud config get-value project)
export TOPIC_NAME="projects/$GCP_PROJECT/topics/twitter-raw"
export SECRET_PATH="/Users/slavayssiere/Code/slavayssiere-sandbox-gcp/iac/sa-pubsub-publisher.json"

docker build -t eu.gcr.io/$GCP_PROJECT/twitter-injector:$version .
docker push eu.gcr.io/$GCP_PROJECT/twitter-injector:$version

# docker run -d \
#     -v $SECRET_PATH:/secret/secret-sa-gcp-pubsub.json \
#     -e PROM_PORT=8080 \
#     -e CONSUMER_KEY=$CONSUMER_KEY \
#     -e CONSUMER_SECRET=$CONSUMER_SECRET \
#     -e ACCESS_TOKEN=$ACCESS_TOKEN \
#     -e ACCESS_SECRET=$ACCESS_SECRET \
#     -e HASHTAG="#CAC40" \
#     -e TOPIC_NAME=$TOPIC_NAME \
#     -e SECRET_PATH=$SECRET_PATH \
#     eu.gcr.io/$GCP_PROJECT/twitter-injector:latest
