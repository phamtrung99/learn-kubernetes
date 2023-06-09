# I. Overview

## 1. Components
- Clusters -> Nodes (one node can be one machine in reality) -> Pods -> Containers

![image][pasted-2023.05.11-11.10.37.png]

 [Original document link](https://kubernetes.io/docs/concepts/overview/components/#control-plane-components) 
-  **Clusters** : EACH cluster have at least 1 control plane components , which contains control manager, deployment, schedules. Normally, control plane components can be setup in same machine or not.
    -  **kubeapi-server** : also called Kubernetes API server, expose API for user or others component in same system. Scaling by horizontally, create more instance of kube-api. The case you need scale is when more user interacted with Kubernetes API, time of response is slow.
    -  **etcd** : a key-value store used by the Kubernetes control plane components to store configuration data and the state of the cluster.
    -  **kube-scheduler** : watch for newly Pod was created and then assigning it to a node.
    -  **kube-controller-manager** : controller processes, but run in single processes. Contain: Node controller (notice when one node down), job controller (create new pod to handle task until completion if that task was be off ) 
    -  ** **cloud-controller-manager** ** : interact with cloud provider API like S3. Example: My node run web app need interact with aws s3. Deep flow of request is my app => Kubernetes control plane => kube-controller-manager => s3 (response flow is opposite)
-  **Node Components** : 
    -  **kubelet** : an agent of node, purpose: manage the deployment of pods on that node, include ensure that the desired state of the pods and containers on the node is in sync with the actual state, responsible for starting, stopping, and monitoring container processes within a pod, and it reports the status of these processes back to the control plane. Totally, its manage pods and containers.
    -  **kube-proxy** : handle and route network traffic between services and pods running on the node.
    - Besides, an software called Container Runtime is responsible for running containers ( like containerd, CRI-O, docker )
-  **Cluster DNS** : All cluster must have cluster DNS, is designed to support communication between nodes and services within a SINGLE Kubernetes cluster, it resolves IP addresses within the cluster by mapping service names to their corresponding IP addresses. Similar like docker DNS, we can interact between docker container with container name easily, without use container IP address. But kubernetes DNS support more advanced feature.

## 2. Features
- Decrease container downtime: Auto start others replicate container when some container down. 
- Set limit resources of container can be usage, like docker.
- Self health check: auto restart, kill, start containers when health check fail without inform to client.
- Roll-based access control: similar dockers (enterprises version only), which user can create, delete container, ….  
- Kubernetes is not monolithic: logging, monitoring, and alerting solutions is optionals and pluggable.

# II. Install.
Firstly, we need two server in order to setup two k8s node. Unless you can install any VMware in your local machine. In this tutorial, I will use  [Multipass](https://multipass.run) to run 2 virtual ubuntu server.
In case you use Mac os apple silicon chip like M1 or M2, please DON'T try to install VMware tool like Virtual box, or use docker container to run many ubuntu (some error will happened related to `systemd`) . 

Step to install Multipass: 
  ***Step 1* : Install** 

``` shell
brew install multipass
```

  ***Step 2* : Setup 2 ubuntu server**  
A master server with 1 core cpu and 1gb ram

``` shell
multipass launch --name master-server --cpus 1 --memory 1024M --disk 5G
```

And worker server with 1 core cpu and 500MB ram

``` shell
multipass launch --name worker-server --cpus 1 --memory 512M --disk 5G
```

Use command `multipass shell [ server-name]` to access of those server.
For example: 

``` shell
multipass shell master-server 
multipass shell worker-server 
```

In case you want to delete any virtual machine, use this command: 

``` shell
multipass delete --purge master-server
multipass delete --purge worker-server
```

Uninstall Multipass: 

``` shell
brew uninstall multipass
```

## 1. Install K8s
Somewhere in the internet, there are many way to setup an k8s, if you read [k8s official document](https://kubernetes.io/docs/tasks/tools/), they are recommend install by Kubeadm but its so complicated for beginners. However, currently we have some tool will make this mission more easily like K3s or k3d. In this article, I use k3s to setup. 

Access to shell of both server
  ***Step 1* : install Containerd** 

``` bash
sudo apt-get update
sudo apt install -y containerd
```

 *Check result* : 

``` bash
systemctl status containerd
```
![image][pasted-2023.05.11-10.47.19.png]

  ***Step 2* : Setup the Master k3s Node** 
Access to shell of Master server
On master server 

``` bash 
curl -sfL https://get.k3s.io | sh -
```
If you need install with k8s specify version, please read  [Here](https://www.rancher.co.jp/docs/k3s/latest/en/installation/#install-script)

Set permission for kubectl 

``` bash
sudo chmod 644 /etc/rancher/k3s/k3s.yaml
```

 *Check result* 

``` bash
systemctl status k3s
```
![image][pasted-2023.05.11-10.54.43.png]

 *Check running node* 

``` bash
kubectl get node -o wide
```
![image][pasted-2023.05.11-10.58.11.png]

 *Note* : We already have a master-node with `IP 192.168.64.9` and container runtime is `containerd://1.6.19-k3s1`, if you install docker instead of containerd, this field will show `docker://`

 ***Step 3*: Allow port on firewall** 
Allow ports that will will be used to communicate between the master and the worker nodes.

On both master and worker server:

``` bash
sudo ufw allow 6443/tcp
sudo ufw allow 443/tcp
```

 *Result* :
![image][pasted-2023.05.11-11.03.25.png]


 ***Step 4*: Extract token from master node** 
This token will be used to join the others to the master node.

On master server:

``` bash
sudo cat /var/lib/rancher/k3s/server/node-token
```

 *Result* : 
![image][pasted-2023.05.11-11.07.55.png]

 ***Step 5*: Install k3s on worker nodes and connect them to the master** 

On worker server:

``` bash 
curl -sfL http://get.k3s.io | K3S_URL=https://<master_IP>:6443 K3S_TOKEN=<join_token> sh -s
```

 **join_token** : take from result of step 4
 **master_IP**: IP address of master node (take from result of step 2)

 *Example* :

``` bash
curl -sfL http://get.k3s.io | K3S_URL=https://192.168.64.9:6443 K3S_TOKEN=K10e95973bef83ad436c24dd649ff5137b89a44049283d9906e9625ff6819083314::server:a55d94ad789f735eec72f5ebce359abd sh -s
```

*Check again in master server:* 

``` bash
kubectl get nodes -o wide
```
![image][pasted-2023.05.11-11.22.33.png]

Explain: Because range IP of cluster is 192.168.64.x so we have a master node (IP 192.168.64.9) and worker node (IP 192.168.64.7) 


# III. Deploy first app to k8s

Flow to deploy an app in k8s:
- Write docker file => build app to image => push to registry (dockerhub, aws, ….) => Access to node, define k8s api resources deployment (yml) and apply.

## 1. Sample project
In this example, i deployed two container in two separated node. One container run a basic restful app golang will retrieve data from mysql database which ran in the others container. For simplicity, I put all step needed into a bash scripting file and if you want to know details step, please view this file. 
Ssh to master-node and clone project in [here](https://github.com/ExecutionLab/k8s_config/tree/master/k8s_demo) 
Run command:

``` bash
bash resources_deploy.sh
```

To check again: 

``` bash
kubectl get all -o wide -n space2
```

Example result: 

![image][pasted-2023.05.15-11.25.20.png]

Access to mysql db and run migrate data.

```  bash
kubectl exec -i -t -n [space-name] [pod-name] -c [container-name] -- sh -c bash
```

Then type: `mysql -p ` with password `1234`
Copy content of `seed.sql` file to above shell session and run.

Example: 
![image][pasted-2023.05.15-11.52.44.png]

Check db connection from app container

``` bash
kubectl logs <pod-name> -c <container-name> 
```

Example result: 
![image][pasted-2023.05.15-14.04.14.png]

To delete all resources: 

``` bash
bash resources_delete.sh
```

# V. Usual command
Below list is some usual basic command when interact with k8s. For more information about any command, please type help suffix `-h` after that command, for example: kubectl get -h. Some resources must specify namespace when retrieve info, so remember use options `-n [space-name]`.
 
 **kubectl get [api-resource-name]:**  get list deployed resources. To get all api-resources name, type `kubectl api-resources`
Ex: kubectl get pods:  

 **kubectl config view --minify --raw** : view config file of kubernetes, use this config to connect from local UI

 **kubectl create namespace [space-name]**: create namespace

 **kubectl describe [resource-type] [resource-name-prefix]** : detailed description of the selected resources
Ex: kubectl describe pods/web-demo-development

 **kubectl logs <pod-name> -c <container-name>**:  Access container logs. 

 **k3s kubectl apply -f [api-resources-definition].yaml** : apply deploy resource.

 **kubectl exec -i -t -n [space-name] [pod-name] -c [container-name] -- sh -c bash**: Exec specific container shell.
 

# VI. Networks

![image][pasted-2023.05.17-11.10.40.png]

 *Network flow: Client => Ingress => Cluster service => Pods* 
The default Kubernetes network model is based on overlay networking, which uses a virtual network on top of the underlying physical network to connect containers and pods across multiple nodes in the cluster


# VII. Storage

 **1. Persistent volume:**  a cluster resource as a storage to save data, file, …. It's can be a remote cloud-storage. PV was be accessed to whole cluster (not depend on name-space,)

 **2. Persistent volume claim:**  a volume was be extracted from a part of Persistent Volume for pod. For example: PV have 10GB, and we define PVC with 2Gb, pvc will take from PV. Note
    - Claim must be same namespace with pod which use claim.
    - Claim can be define in order to use by specify container or pods.

 **3. Storage class:** is a way to describe the types and capabilities of available storage in a cluster. 
- Dynamic provisioning: use this when admin no need to setup storage by hand for each application, we can define storage class development, and kubernetes create will create storage depend on your definition. 

# VIII. Helm
[Official documents](https://helm.sh/docs/#welcome) 
Imaging that helm look like docker-compose, terraform,...We can create and share helm charts in public registres.
Why to use helm:
- In microservice, instead of define many api resrouces definition yaml file, with helm we just create one configure file.
- We can define only 1 template file, many service which have the same configuration but different value only, can refer to a same template file and each service have their own values.yaml file. 

[Install tutorial](https://helm.sh/docs/intro/install/#from-apt-debianubuntu)

After installed, need specify kubectl config to helm, run command:

``` bash
export KUBECONFIG=/path_to_your_kubeconfig_file
```

**1. Structure**  
Chart file structure: [Ref](https://helm.sh/docs/topics/charts/#the-chart-file-structure)  
Helm infra structure (helm version 2): cleint (helm CLI) -> (Server (Tiller) -> K8s API server) inside a k8s cluster 
Helm server track history of charts, so we can apply upgrade or rollback previous charts version. [Ref](https://helm.sh/docs/helm/helm_history/#helm-history) (look like goose in golang manage migrate version)

**2. Usual command**  
[helm history](https://helm.sh/docs/helm/helm_history/#helm-history): Show chart execution history.

helm template .: show all final template after apply value reference and will show syntax error if occurs.  

[helm install](https://helm.sh/docs/helm/helm_install/#helm-install): install kubernetes package, browse kubernetes packages in [artifacthub](https://artifacthub.io), when type helm install, helm will find in templates directory and run kubectl apply.

[helm delete](): delete an installed chart, terminate all resources related to


**3. Template**  
Start learn how to write helm chart template from scratch [Ref](https://helm.sh/docs/chart_template_guide/getting_started/#charts)  
Some keyword need must to know:   
[Built-in Objects](https://helm.sh/docs/chart_template_guide/builtin_objects/#helm)  
[Template function](https://helm.sh/docs/chart_template_guide/functions_and_pipelines/): some function to transform data which get from value before add to template, list all template function in [ref](https://helm.sh/docs/chart_template_guide/function_list/#helm)

Optionals knowlegde:
[Flow control](https://helm.sh/docs/chart_template_guide/control_structures/): how to use if/else, range, ect in helm template.



# IV. TLS certificate

Because I use an global DNS service to resgister new domain, so we can't setup https for an service in local, we need an server which be atached with public IPv4 address.
For more visualization, I will run an simple todoapp (git repo links) and setup k8s and resgister a certifcate.
(continue)


#IV
