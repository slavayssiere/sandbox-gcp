#!/bin/bash

kubectl delete -f ../consumer/consumer-deploy.yaml
kubectl delete -f ../normalizer/normalizer-deploy.yaml
kubectl delete -f ../bigtable/bigtable-load-deploy.yaml

kubectl delete -f ../injector/

