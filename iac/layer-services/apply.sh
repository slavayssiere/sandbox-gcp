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
helm install --name prometheuses stable/prometheus-operator --namespace monitoring

# wget https://storage.googleapis.com/gke-release/istio/release/1.0.3-gke.0/stackdriver/stackdriver-tracing.yaml
# kubectl apply -f stackdriver-tracing.yaml
# rm stackdriver-tracing.yaml

############################ replace by istio ############################
# apply_kubectl "traefik-consul"
# apply_kubectl "traefik-app" 
# apply_kubectl "traefik-admin"
############################ /replace by istio ############################
apply_kubectl "external-dns"

# if [ ! -d "istio-1.0.3" ]; then
#     wget https://github.com/istio/istio/releases/download/1.0.3/istio-1.0.3-osx.tar.gz
#     tar -xvf istio-1.0.3-osx.tar.gz
#     rm istio-1.0.3-osx.tar.gz
# fi

# cd istio-1.0.3 
#     helm install install/kubernetes/helm/istio --name istio --namespace istio-system -f ../istio/values-istio-1.0.3.yaml
# cd -

if [ ! -d "istio-1.0.5" ]; then
    wget https://github.com/istio/istio/releases/download/1.0.5/istio-1.0.5-osx.tar.gz
    tar -xvf istio-1.0.5-osx.tar.gz
    rm istio-1.0.5-osx.tar.gz
fi

cd istio-1.0.5
    helm install install/kubernetes/helm/istio --name istio --namespace istio-system -f ../istio/values-istio-1.0.5.yaml
cd -

apply_kubectl "monitoring"

## delete istio CRD
# kubectl delete customresourcedefinitions $(kubectl get customresourcedefinitions | cut -d ' ' -f1 | grep istio.io)