#!/bin/bash
set -e

kubectl create namespace space3
kubectl apply -f ./resources/todo-app-deployment.yaml

