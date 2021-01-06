
Tools Needed
* Docker
* Kind
* Helm

Kind install can be found [here](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)

Helm install can be found [here](https://helm.sh/docs/helm/helm_install/)


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

```bash
kubectl cluster-info --context kind-kind
Kubernetes master is running at https://127.0.0.1:59511
KubeDNS is running at https://127.0.0.1:59511/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
```

```bash
helm repo add cilium https://helm.cilium.io/
docker pull cilium/cilium:v1.9.1
kind load docker-image cilium/cilium:v1.9.1
```

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

