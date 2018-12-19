#!/bin/bash

kubectl create -f ../consumer/consumer-deploy.yaml
kubectl create -f ../normalizer/normalizer-deploy.yaml
kubectl create -f ../bigtable/bigtable-load-deploy.yaml

kubectl create -f ../injector/