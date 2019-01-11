# Localy do:


export GCP_PROJECT=$(gcloud config get-value project)
export TOPIC_NAME="projects/$GCP_PROJECT/topics/raw-mastodon"
export SECRET_PATH="$HOME/gcp-secret/bigdata-project.json"

docker build -t eu.gcr.io/$GCP_PROJECT/mastodon-injector:latest .

# docker run -d \
#     -v $SECRET_PATH:/secret/secret-sa-gcp-pubsub.json \
#     -e CLIENT_ID="" \
#     -e CLIENT_SECRET="" \
#     -e SERVER="https://linuxrocks.online" \
#     -e LOGIN="" \
#     -e PASSWORD="" \
#     -e HASHTAG="#CAC40" \
#     -e TOPIC_NAME=$TOPIC_NAME \
#     -e SECRET_PATH=$SECRET_PATH \
#     eu.gcr.io/$GCP_PROJECT/mastodon-injector:latest
