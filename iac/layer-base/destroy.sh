REGION="europe-west1"
SUBNET_EUROPE1_CIDR="192.168.0.0/20"

gcloud config set project slavayssiere-sandbox
gcloud config set compute/region $REGION

gcloud -q beta compute routers nats delete nat-$REGION --router=nat-$REGION
gcloud -q beta compute routers delete nat-$REGION

terraform destroy \
    -auto-approve

# # bastion

# gcloud compute instances delete nat-test-1

# # kubernetes

# gcloud container clusters delete private-cluster-1 


# # network

# gcloud compute firewall-rules delete fw-0

# gcloud compute networks subnets delete subnet-0

# gcloud compute networks delete net-0 
