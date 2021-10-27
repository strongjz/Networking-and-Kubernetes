
Tools Needed
* Docker
* Kind
* Helm

Kind install can be found [here](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)

Helm install can be found [here](https://helm.sh/docs/intro/install/)

Cilium Rules basic are available [here](https://docs.cilium.io/en/v1.9/policy/intro/#rule-basics)

Steps
1. Create Kind cluster
2. Add Cilium images to kind cluster
3. Install Cilium in the cluster
4. Test connectivity
5. Test Webserver and Database NetworkPolicies 

# 1. Create Kind cluster

With the kind cluster configuration yaml, we can use kind to create that cluster with the below command. If this is the first time running it, it will take some time to download all the docker images for the working and control plane docker images.

```bash
kind create cluster --config=kind-config.yaml
Creating cluster "kind" ...
✓ Ensuring node image (kindest/node:v1.18.2) 🖼
✓ Preparing nodes 📦 📦 📦 📦
✓ Writing configuration 📜
✓ Starting control-plane 🕹️
✓ Installing StorageClass 💾
✓ Joining worker nodes 🚜
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind

Have a question, bug, or feature request? Let us know! https://kind.sigs.k8s.io/#community 🙂
```

Always verify that the cluster is up and running with kubectl.

```bash
kubectl cluster-info --context kind-kind
Kubernetes master is running at https://127.0.0.1:59511
KubeDNS is running at https://127.0.0.1:59511/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
```

# 2. Add Cilium images to kind cluster
Now that our cluster is running locally we can begin installing Cilium using helm, a kubernetes deployment tool. This is the prefered way to install Cilium. First, we need to add the helm repo for Cilium. Then download the docker images for cilium, and finally instruct kind to load the cilium images into the cluster.

```bash
helm repo add cilium https://helm.cilium.io/
docker pull cilium/cilium:v1.9.1
kind load docker-image cilium/cilium:v1.9.1
```

# 3. Install Cilium in the cluster

Now the pre-requisites for Cilium are completed we can install Cilium in our cluster with helm. There are many configuration options for Ciluim, and they are set with the helm options --set.

```bash
helm install cilium cilium/cilium --version 1.9.1 \
 --namespace kube-system \
 --set nodeinit.enabled=true \
 --set kubeProxyReplacement=partial \
 --set hostServices.enabled=false \
 --set externalIPs.enabled=true \
 --set nodePort.enabled=true \
 --set hostPort.enabled=true \
 --set bpf.masquerade=false \
 --set image.pullPolicy=IfNotPresent \
--set ipam.mode=kubernetes
```

# 4. Test connectivity

Now that Cilium is deployed we can run the connectivity check from Cilium to ensure the CNI is installed in the cluster correctly.

```bash
kubectl create ns cilium-test
namespace/cilium-test created

kubectl apply -n cilium-test -f https://raw.githubusercontent.com/strongjz/advanced_networking_code_examples/master/chapter-4/connectivity-check.yaml
deployment.apps/echo-a created
deployment.apps/echo-b created
deployment.apps/echo-b-host created
deployment.apps/pod-to-a created
deployment.apps/pod-to-external-1111 created
deployment.apps/pod-to-a-denied-cnp created
deployment.apps/pod-to-a-allowed-cnp created
deployment.apps/pod-to-external-fqdn-allow-google-cnp created
deployment.apps/pod-to-b-multi-node-clusterip created
deployment.apps/pod-to-b-multi-node-headless created
deployment.apps/host-to-b-multi-node-clusterip created
deployment.apps/host-to-b-multi-node-headless created
deployment.apps/pod-to-b-multi-node-nodeport created
deployment.apps/pod-to-b-intra-node-nodeport created
service/echo-a created
service/echo-b created
service/echo-b-headless created
service/echo-b-host-headless created
ciliumnetworkpolicy.cilium.io/pod-to-a-denied-cnp created
ciliumnetworkpolicy.cilium.io/pod-to-a-allowed-cnp created
ciliumnetworkpolicy.cilium.io/pod-to-external-fqdn-allow-google-cnp created
```

Cilium installs several pieces in the cluster, the agent, the client, operator and the cilium-cni plugin.

Agent - The Cilium agent, cilium-agent, runs on each node in the cluster. The agent accepts configuration via Kubernetes or APIs that describes networking, service load-balancing, network policies, and visibility & monitoring requirements.

Client (CLI) - The Cilium CLI client (cilium) is a command-line tool that is installed along with the Cilium agent. It interacts with the REST API of the Cilium agent running on the same node. The CLI allows inspecting the state and status of the local agent. It also provides tooling to directly access the eBPF maps to validate their state.

Operator - The Cilium Operator is responsible for managing duties in the cluster which should logically be handled once for the entire cluster, rather than once for each node in the cluster.

CNI Plugin - The CNI plugin (cilium-cni) interacts with the Cilium API of the node to trigger the configuration to provide networking, load-balancing and network policies for the pod.


We can observe all these components being deployed in the cluster with the kubectl -n kube-system get pods --watch command.

```bash
kubectl get pods -n cilium-test -w
NAME                                                     READY   STATUS    RESTARTS   AGE
echo-a-57cbbd9b8b-szn94                                  1/1     Running   0          34m
echo-b-6db5fc8ff8-wkcr6                                  1/1     Running   0          34m
echo-b-host-76d89978c-dsjm8                              1/1     Running   0          34m
host-to-b-multi-node-clusterip-fd6868749-7zkcr           1/1     Running   2          34m
host-to-b-multi-node-headless-54fbc4659f-z4rtd           1/1     Running   2          34m
pod-to-a-648fd74787-x27hc                                1/1     Running   1          34m
pod-to-a-allowed-cnp-7776c879f-6rq7z                     1/1     Running   0          34m
pod-to-a-denied-cnp-b5ff897c7-qp5kp                      1/1     Running   0          34m
pod-to-b-intra-node-nodeport-6546644d59-qkmck            1/1     Running   2          34m
pod-to-b-multi-node-clusterip-7d54c74c5f-4j7pm           1/1     Running   2          34m
pod-to-b-multi-node-headless-76db68d547-fhlz7            1/1     Running   2          34m
pod-to-b-multi-node-nodeport-7496df84d7-5z872            1/1     Running   2          34m
pod-to-external-1111-6d4f9d9645-kfl4x                    1/1     Running   0          34m
pod-to-external-fqdn-allow-google-cnp-5bc496897c-bnlqs   1/1     Running   0          34m
```

# 5. Test Webserver and Database NetworkPolicies 

Now that the Cilium CNI is deployed into our cluster we can begin exploring the power of its Network policies. We 
will deploy our golang webserver that now connects to a database. Using a network utility pod we will test connectivity 
without the network policies in place, then deploy network policies that will restrict connectivity to the web 
server and database. 

1. Deploy Containers, Web, DB and Utils
2. Test open connectivity 
3. Deploy Network policies 
4. Test Closed Network Connectivity

#### 1. Deploy Golang web server

Our Golang web server has been updated to connect to a postgres database. Let's deploy the Postgres database with 
the following yaml and commands. 

1.1 Deploy Database

```bash
kubectl apply -f database.yaml 
service/postgres created
configmap/postgres-config created
statefulset.apps/postgres created
```

Deploying our Webserver as a kubernetes deployment to our kind cluster. 

1.2  Deploy Web Server

```bash
 kubectl apply -f web.yml 
deployment.apps/app created 
```

To run connectivity tests inside the cluster network we will deploy and use a dns utils pod that has basic 
networking tools like ping and curl. 

1.3 Deploy Dns Utils pod

```bash
kubectl apply -f dnsutils.yaml
pod/dnsutils created
```

#### 2. Test open connectivity

Since we are not deploying A service with an ingress, we can use kubectl port forward to test connectivity to our 
webserver 

More information about kubectl port-forward can be found [here](https://kubernetes.io/docs/tasks/access-application-cluster/port-forward-access-application-cluster/)
```bash
kubectl port-forward app-5878d69796-j889q 8080:8080
```

Now from our local terminal we can reach out API. 

```bash
curl localhost:8080/
Hello
curl localhost:8080/healthz
Healthy
curl localhost:8080/data
Database Connected
```

Let's test connectivity to our web server inside the cluster from other pods. In order to do that we need to get the 
IP address of our web server pod. 

```bash
kubectl get pods -l app=app -o wide
NAME                   READY   STATUS    RESTARTS   AGE   IP             NODE           NOMINATED NODE   READINESS GATES
app-5878d69796-j889q   1/1     Running   0          87m   10.244.1.188   kind-worker3   <none>           <none>
```

Now we can test layer4 and 7 connectivity to the web server from the DNS utils pod. 

```bash
kubectl exec dnsutils -- nc -z -vv 10.244.1.188 8080
10.244.1.188 (10.244.1.188:8080) open
sent 0, rcvd 0
```

Layer 7 HTTP API Access
```bash
kubectl exec dnsutils -- wget -qO- 10.244.1.188:8080/
Hello

kubectl exec dnsutils -- wget -qO- 10.244.1.188:8080/data
Database Connected

kubectl exec dnsutils -- wget -qO- 10.244.1.188:8080/healthz
Healthy
```

We can also test the same to the database pod. 

Retrieve the IP Address of database pod.
```bash
kubectl get pods -l app=postgres -o wide
NAME         READY   STATUS    RESTARTS   AGE   IP             NODE          NOMINATED NODE   READINESS GATES
postgres-0   1/1     Running   0          98m   10.244.2.189   kind-worker   <none>           <none>
```

DNS Utils Connectivity
```bash
kubectl exec dnsutils -- nc -z -vv 10.244.2.189 5432
10.244.2.189 (10.244.2.189:5432) open
sent 0, rcvd 0
```


#### 3. Deploy Network policies and Test Closed Network Connectivity

Let's first restrict access to the database pod to only the Web server.

The postgress port 5432 is open from dnsutils to database. 

```bash
kubectl exec dnsutils -- nc -z -vv -w 5 10.244.2.189 5432
10.244.2.189 (10.244.2.189:5432) open
sent 0, rcvd 0
```

Apply the Network policy that only allows traffic from the Web Server pod to the database.

```bash
kubectl apply -f layer_3_net_pol.yaml
ciliumnetworkpolicy.cilium.io/l3-rule-app-to-db created
```

With the network policy applied, the dnsutils pod can no longer reach the database pod. 

```bash
kubectl exec dnsutils -- nc -z -vv -w 5 10.244.2.189 5432
nc: 10.244.2.189 (10.244.2.189:5432): Operation timed out
sent 0, rcvd 0
command terminated with exit code 1
```

But we the Web server is still connected to the Database. 

```bash
kubectl exec dnsutils -- wget -qO- 10.244.1.188:8080/data
Database Connected

curl localhost:8080/data
Database Connected
```

The Cilium install and deploy of cilium objects creates resources that can retrieved just like pods with kubectl. 

```bash 
kubectl describe ciliumnetworkpolicies.cilium.io l3-rule-app-to-db
Name:         l3-rule-app-to-db
Namespace:    default
Labels:       <none>
Annotations:  API Version:  cilium.io/v2
Kind:         CiliumNetworkPolicy
Metadata:
Creation Timestamp:  2021-01-10T01:06:13Z
Generation:          1
Managed Fields:
API Version:  cilium.io/v2
Fields Type:  FieldsV1
fieldsV1:
f:metadata:
f:annotations:
.:
f:kubectl.kubernetes.io/last-applied-configuration:
f:spec:
.:
f:endpointSelector:
.:
f:matchLabels:
.:
f:app:
f:ingress:
Manager:         kubectl
Operation:       Update
Time:            2021-01-10T01:06:13Z
Resource Version:  47377
Self Link:         /apis/cilium.io/v2/namespaces/default/ciliumnetworkpolicies/l3-rule-app-to-db
UID:               71ee6571-9551-449d-8f3e-c177becda35a
Spec:
Endpoint Selector:
Match Labels:
App:  postgres
Ingress:
From Endpoints:
Match Labels:
App:  app
Events:       <none>
```

Now let us apply the Layer 7 policy. Cilium is layer 7 aware, so we can block or allow certain base on HTTP URI paths. 
In our example policy we allow HTTP GETs on / and /data but not allow on /healthz, lets test that out. 

```bash
kubectl apply -f layer_7_netpol.yml
ciliumnetworkpolicy.cilium.io/l7-rule created
```

```bash
kubectl get ciliumnetworkpolicies.cilium.io
NAME      AGE
l7-rule   6m54s

kubectl describe ciliumnetworkpolicies.cilium.io l7-rule
Name:         l7-rule
Namespace:    default
Labels:       <none>
Annotations:  API Version:  cilium.io/v2
Kind:         CiliumNetworkPolicy
Metadata:
  Creation Timestamp:  2021-01-10T00:49:34Z
  Generation:          1
  Managed Fields:
    API Version:  cilium.io/v2
    Fields Type:  FieldsV1
    fieldsV1:
      f:metadata:
        f:annotations:
          .:
          f:kubectl.kubernetes.io/last-applied-configuration:
      f:spec:
        .:
        f:egress:
        f:endpointSelector:
          .:
          f:matchLabels:
            .:
            f:app:
    Manager:         kubectl
    Operation:       Update
    Time:            2021-01-10T00:49:34Z
  Resource Version:  43869
  Self Link:         /apis/cilium.io/v2/namespaces/default/ciliumnetworkpolicies/l7-rule
  UID:               0162c16e-dd55-4020-83b9-464bb625b164
Spec:
  Egress:
    To Ports:
      Ports:
        Port:      8080
        Protocol:  TCP
      Rules:
        Http:
          Method:  GET
          Path:    /
          Method:  GET
          Path:    /data
  Endpoint Selector:
    Match Labels:
      App:  app
Events:     <none>
```

As you can see, / and /data are available by not /healthz

```bash
kubectl exec dnsutils -- wget -qO- 10.244.1.188:8080/data
Database Connected

kubectl exec dnsutils -- wget -qO- 10.244.1.188:8080/
Hello

kubectl exec dnsutils -- wget -qO- -T 5 10.244.1.188:8080/healthz
wget: error getting response
command terminated with exit code 1
```
