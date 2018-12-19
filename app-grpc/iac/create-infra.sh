
PROJECT="cool-wharf-207907"
KUBE_VERSION="1.10.4-gke.2"
CLUSTER_NAME="cluster-test-pubsub"
MONITORING_POOL_NAME="monitoring"
BT_INSTANCE="bt-instance-test"

gcloud config set project $PROJECT
gcloud config set compute/zone europe-west1-b

gcloud components install cbt
gcloud components install kubectl

NB_NODE=90
NB_NODE_MONITORING=3

gcloud beta container clusters create $CLUSTER_NAME \
  --cluster-version $KUBE_VERSION \
  --machine-type "n1-standard-2" \
  --image-type "COS" \
  --scopes "https://www.googleapis.com/auth/cloud-platform","https://www.googleapis.com/auth/compute","https://www.googleapis.com/auth/devstorage.read_write","https://www.googleapis.com/auth/logging.write","https://www.googleapis.com/auth/monitoring","https://www.googleapis.com/auth/pubsub","https://www.googleapis.com/auth/servicecontrol","https://www.googleapis.com/auth/service.management.readonly","https://www.googleapis.com/auth/sqlservice.admin","https://www.googleapis.com/auth/trace.append" \
  --num-nodes $NB_NODE \
  --enable-cloud-logging \
  --enable-cloud-monitoring \
  --enable-autorepair \
  --enable-ip-alias \
  --network "default" \
  --subnetwork "default" \
  --addons HorizontalPodAutoscaling,HttpLoadBalancing

  # --enable-autoscaling \
  #   --min-nodes $NB_NODE \
  #   --max-nodes "200" \

gcloud container node-pools create $MONITORING_POOL_NAME \
  --cluster $CLUSTER_NAME \
  --image-type "COS" \
  --enable-autorepair \
  --num-nodes $NB_NODE_MONITORING \
  --machine-type=n1-standard-4 \
  --node-version $KUBE_VERSION \
  --zone=europe-west1-b \
  --node-taints dedicated=monitoring:NoExecute \
  --scopes "https://www.googleapis.com/auth/cloud-platform","https://www.googleapis.com/auth/compute","https://www.googleapis.com/auth/devstorage.read_write","https://www.googleapis.com/auth/logging.write","https://www.googleapis.com/auth/monitoring","https://www.googleapis.com/auth/pubsub","https://www.googleapis.com/auth/servicecontrol","https://www.googleapis.com/auth/service.management.readonly","https://www.googleapis.com/auth/sqlservice.admin","https://www.googleapis.com/auth/trace.append"

  # --enable-autoscaling \
  #   --max-nodes=4 \
  #   --min-nodes=1 \

gcloud container clusters get-credentials $CLUSTER_NAME

kubectl create clusterrolebinding cluster-admin-binding \
  --clusterrole cluster-admin \
  --user "$(gcloud config get-value account)"

kubectl create -f ./pubsub-monitoring.yaml
kubectl create -f ./secret-sa.yaml

kubectl apply -f ./traefik-ic/traefik-rbac.yaml
kubectl apply -f ./traefik-ic/traefik-ds.yaml

kubectl taint nodes $(kubectl get nodes | grep "^gke-$CLUSTER_NAME-$MONITORING_POOL_NAME" | cut -f 1 -d' ') dedicated=monitoring:NoExecute
kubectl label nodes $(kubectl get nodes | grep "^gke-$CLUSTER_NAME-$MONITORING_POOL_NAME" | cut -f 1 -d' ') dedicated=monitoring
kubectl label nodes $(kubectl get nodes | grep "^gke-cluster-test-pubsub-default-pool" | cut -f 1 -d' ') dedicated=compute --overwrite


kubectl create -f ./monitoring/namespace.yaml
kubectl create -f ./monitoring/prometheus-k8s.yaml
kubectl create -f ./monitoring/prometheus-injectors.yaml
kubectl create -f ./monitoring/prometheus-consumers.yaml
kubectl create -f ./monitoring/prometheus-bigtable.yaml
kubectl create -f ./monitoring/grafana-dashboard.yaml
kubectl create -f ./monitoring/grafana.yaml

gcloud alpha pubsub topics create "normalized_topic"
gcloud alpha pubsub subscriptions create "normalized_subcription" \
  --ack-deadline 60 \
  --topic "normalized_topic" \
  --topic-project $PROJECT


# Create subscription
gcloud alpha pubsub subscriptions create "trade_subcription" \
  --ack-deadline 60 \
  --topic "normalized_topic" \
  --topic-project $PROJECT

gcloud beta bigtable instances create $BT_INSTANCE \
  --cluster=$BT_INSTANCE \
  --cluster-zone="europe-west1-b" \
  --display-name=$BT_INSTANCE \
  --cluster-num-nodes=3

# Sample code to bootstrap BigTable structure and inserting a sample value
cbt -project $PROJECT -instance $BT_INSTANCE createtable my-table

cbt -instance $BT_INSTANCE createfamily my-table tests

# kubectl logs `kubectl get pods | grep '^bigtable' | grep Running  | head -n 1 | cut -f 1 -d' '`
# kubectl delete pod --force --grace-period=0 -n monitoring $(kubectl get pods | grep "^Evicted" | cut -f 1 -d' ')


