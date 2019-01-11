# Localy do:
export GCP_PROJECT=$(gcloud config get-value project)
export TOPIC_NAME="projects/$GCP_PROJECT/topics/raw-twitter"
export SECRET_PATH="$HOME/gcp-secret/bigdata-project.json"

docker build -t eu.gcr.io/$GCP_PROJECT/twitter-injector:0.0.1 .

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
