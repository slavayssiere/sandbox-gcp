#!/bin/bash

test_tiller_present() {
    kubectl get pod -n kube-system -l app=helm,name=tiller | grep Running | wc -l | tr -d ' '
}

apply_kubectl() {
    cd $1
    kubectl apply -f .
    cd ..
}

gcloud container clusters get-credentials test-cluster --region europe-west1
username=$(gcloud config get-value account)
kubectl create clusterrolebinding cluster-admin-binding --clusterrole=cluster-admin --user=$username
    
# install prometheus operator
kubectl apply -f helm/rbac.yaml
helm init  --service-account tiller


test_tiller=$(test_tiller_present)
while [ $test_tiller -lt 1 ]; do
    echo "Wait for Tiller: $test_tiller"
    test_tiller=$(test_tiller_present)
    sleep 1
done

sleep 10

kubectl create ns monitoring
helm install stable/prometheus-operator --namespace monitoring

# wget https://storage.googleapis.com/gke-release/istio/release/1.0.3-gke.0/stackdriver/stackdriver-tracing.yaml
# kubectl apply -f stackdriver-tracing.yaml
# rm stackdriver-tracing.yaml

apply_kubectl "traefik-consul"
apply_kubectl "traefik-app"
apply_kubectl "traefik-admin"
apply_kubectl "external-dns"