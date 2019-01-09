#!/bin/bash

# install prometheus operator

wget https://storage.googleapis.com/gke-release/istio/release/1.0.3-gke.0/stackdriver/stackdriver-tracing.yaml -o stackdriver-tracing.yaml
kubectl apply -f stackdriver-tracing.yaml

cd traefik-consul
kubectl apply -f .
cd ..

cd traefik-app
kubectl apply -f .
cd ..

cd traefik-admin
kubectl apply -f .
cd ..

