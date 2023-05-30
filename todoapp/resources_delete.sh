#!/bin/bash
set -e

kubectl delete deployment todo-app-deployment -n space3
kubectl delete namespace space3
