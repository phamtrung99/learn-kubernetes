#!/bin/bash
set -e

kubectl create namespace space2
kubectl apply -f ./resources/mysql-storage.yaml
kubectl apply -f ./resources/mysql-deployment.yaml
kubectl apply -f ./resources/web-demo-secret.yaml -n space2
kubectl apply -f ./resources/web-demo-deployment.yaml

