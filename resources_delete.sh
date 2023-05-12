#!/bin/bash
set -e

kubectl delete deployment web-demo-deployment -n space2
kubectl delete secret web-demo-secret -n space2
kubectl delete deployment mysql-deployment -n space2
kubectl delete persistentvolumeclaim mysql-pv-claim -n space2
kubectl delete persistentvolume mysql-pv-volume

