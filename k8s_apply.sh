#!/bin/bash
set -e

kubectl create namespace space2
kubectl apply -f ./mysql-storage.yaml
kubectl apply -f ./mysql-deployment.yaml
kubectl apply -f ./web-demo-secret.yaml
kubectl apply -f ./web-demo-deployment.yaml

