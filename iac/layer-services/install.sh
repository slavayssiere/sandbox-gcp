#!/bin/bash

cd traefik-consul
kubectl apply -f .
cd ..

cd traefik-app
kubectl apply -f .
cd ..

cd traefik-admin
kubectl apply -f .
cd ..

