# THIS FEATURE JUST A EXPERIMENTS, DON'T APPLY IT IN OFFWORK SYSTEM

# VELERO

Tool backup k8s resource.
- Support store file in both local and cloud storage.
- Support schedule cron job to backup

## INSTALL 
Need to install both client machine and server  
**Step 1: Install local CLI** [Ref](https://velero.io/docs/v1.11/basic-install/#install-the-cli) 

**Step 2: Install server component** [Ref](https://github.com/vmware-tanzu/helm-charts/blob/main/charts/velero/README.md#option-2-yaml-file) 

1. In client, add repo by helm: [Ref](https://github.com/vmware-tanzu/helm-charts#usage)

``` bash
helm repo add vmware-tanzu https://vmware-tanzu.github.io/helm-charts
```

Install velero:  

``` bash
kubectl create namespace velero
helm install vmware-tanzu/velero --namespace velero -f ./values.demo.yaml --generate-name
```

## CONFIG BACKUP LOCATION
In case setup aws s3 storage location, please follow this [Ref](https://github.com/vmware-tanzu/velero-plugin-for-aws#setup)  

``` bash
kubectl apply -f resources/backupLocation.yaml
```


## USAGE
Create backup of mysql pod with persistent volume contain mysql data.  

``` bash
velero backup create test-backup-2 --include-namespaces=offwork-v2 --selector app=mysql --include-resources=persistentvolumeclaims/data-offwork-v2-mysql-master-0 --storage-location bsl-aws-s3
```



