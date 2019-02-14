# Localy do:

version=$1

export GCP_PROJECT=$(gcloud config get-value project)

docker build -t eu.gcr.io/$GCP_PROJECT/app-public:$version .
docker push eu.gcr.io/$GCP_PROJECT/app-public:$version

# docker run -d \
#     -v $SECRET_PATH:/secret/secret-sa-gcp-pubsub.json \
#     -e SUB_NAME=$TOPIC_NAME \
#     -e SECRET_PATH=$SECRET_PATH \
#     eu.gcr.io/$GCP_PROJECT/mastodon-injector:latest
