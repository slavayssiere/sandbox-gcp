
PROJECT="cool-wharf-207907"
KUBE_VERSION="1.10.4-gke.2"
CLUSTER_NAME="cluster-test-pubsub"
BT_INSTANCE="bt-instance-test"
MONITORING_POOL_NAME="monitoring"

gcloud beta container clusters delete $CLUSTER_NAME
gcloud container node-pools create $MONITORING_POOL_NAME

gcloud alpha pubsub subscriptions delete "normalized_subcription"
gcloud alpha pubsub subscriptions delete "trade_subcription"

gcloud alpha pubsub topics delete "normalized_topic" 

cbt -project $PROJECT deleteinstance $BT_INSTANCE 
