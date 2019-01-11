REGION="europe-west1"
SUBNET_EUROPE1_CIDR="192.168.0.0/20"

gcloud config set project slavayssiere-sandbox
gcloud config set compute/region $REGION

gcloud -q beta compute routers nats delete nat-$REGION --router=nat-$REGION
gcloud -q beta compute routers delete nat-$REGION
gcloud -q beta dns managed-zones delete private-dns-zone

terraform destroy \
    -auto-approve

gcloud -q iam service-accounts delete "sa-pubsub-publisher"
gcloud -q iam service-accounts delete "sa-pubsub-subscriber"
gcloud -q iam service-accounts delete "sa-pubsub-full"
gcloud -q iam service-accounts delete "sa-aggregator"
