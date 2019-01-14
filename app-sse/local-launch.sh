# Localy do:

version=$1

cd src

export GCP_PROJECT=$(gcloud config get-value project)
export SUB_NAME="projects/slavayssiere-sandbox/subscriptions/messages-normalized-sub"
export SECRET_PATH="/Users/slavayssiere/Code/slavayssiere-sandbox-gcp/iac/sa-pubsub-subscriber.json"

go build && ./app-sse

# docker run -it -p 8080:8080 \
#     -e SUB_NAME="projects/slavayssiere-sandbox/subscriptions/messages-normalized-sub" \
#     -e SECRET_PATH="/secret/sa-pubsub-subscriber.json" \
#     -v /Users/slavayssiere/Code/slavayssiere-sandbox-gcp/iac/sa-pubsub-subscriber.json:/secret/sa-pubsub-subscriber.json \
#     --entrypoint sh \
#     eu.gcr.io/slavayssiere-sandbox/app-sse:0.0.4


# docker run -d -p 8080:8080 \
#     -e SUB_NAME="projects/slavayssiere-sandbox/subscriptions/messages-normalized-sub" \
#     -e SECRET_PATH="/secret/sa-pubsub-subscriber.json" \
#     -v /Users/slavayssiere/Code/slavayssiere-sandbox-gcp/iac/sa-pubsub-subscriber.json:/secret/sa-pubsub-subscriber.json \
#     eu.gcr.io/slavayssiere-sandbox/app-sse:0.0.4
