# Localy do:


export GCP_PROJECT=$(gcloud config get-value project)
export SUB_NAME="projects/$GCP_PROJECT/subscriptions/raw-normalizer"
export SECRET_PATH="$HOME/gcp-secret/bigdata-project.json"

docker build -t eu.gcr.io/$GCP_PROJECT/app-sse:latest .

# docker run -d \
#     -v $SECRET_PATH:/secret/secret-sa-gcp-pubsub.json \
#     -e SUB_NAME=$TOPIC_NAME \
#     -e SECRET_PATH=$SECRET_PATH \
#     eu.gcr.io/$GCP_PROJECT/mastodon-injector:latest
