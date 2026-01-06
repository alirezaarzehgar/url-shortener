# Install scylladb operator
**NOTE:** This part of project is not complete

```bash
# Cert Manager
kubectl apply --server-side -f=https://raw.githubusercontent.com/scylladb/scylla-operator/v1.19/examples/third-party/cert-manager.yaml

kubectl wait --for condition=established --timeout=60s crd/certificates.cert-manager.io crd/issuers.cert-manager.io

for deploy in cert-manager{,-cainjector,-webhook}; do
    kubectl -n=cert-manager rollout status --timeout=10m deployment.apps/"${deploy}"
done

for i in {1..30}; do
    { kubectl -n=cert-manager get secret/cert-manager-webhook-ca && break; } || sleep 1
done

# Prometheus
kubectl apply -n prometheus-operator --server-side -f=https://raw.githubusercontent.com/scylladb/scylla-operator/v1.19/examples/third-party/prometheus-operator.yaml

kubectl wait --for='condition=established' crd/prometheuses.monitoring.coreos.com crd/prometheusrules.monitoring.coreos.com crd/servicemonitors.monitoring.coreos.com

kubectl -n=prometheus-operator rollout status --timeout=10m deployment.apps/prometheus-operator

# ScyllaDB Operator
kubectl -n=scylla-operator apply --server-side -f=https://raw.githubusercontent.com/scylladb/scylla-operator/v1.19/deploy/operator.yaml

# Apply ScyllaDB CRD
kubectl apply -f ./k8s/scylla/operator.yaml

# Storage
kubectl -n=scylla-operator apply --server-side -f=https://raw.githubusercontent.com/scylladb/scylla-operator/v1.19/examples/generic/nodeconfig-alpha.yaml

kubectl wait --timeout=10m --for='condition=Progressing=False' nodeconfigs.scylla.scylladb.com/scylladb-nodepool-1
kubectl wait --timeout=10m --for='condition=Degraded=False' nodeconfigs.scylla.scylladb.com/scylladb-nodepool-1
kubectl wait --timeout=10m --for='condition=Available=True' nodeconfigs.scylla.scylladb.com/scylladb-nodepool-1

kubectl -n=local-csi-driver apply --server-side -f=https://raw.githubusercontent.com/scylladb/scylla-operator/v1.19/examples/common/local-volume-provisioner/local-csi-driver/{00_clusterrole_def,00_clusterrole_def_openshift,00_clusterrole,00_namespace,00_scylladb-local-xfs.storageclass,10_csidriver,10_serviceaccount,20_clusterrolebinding,50_daemonset}.yaml

kubectl -n=local-csi-driver rollout status --timeout=10m daemonset.apps/local-csi-driver

# ScyllaDB Manager
kubectl -n=scylla-manager apply --server-side -f=https://raw.githubusercontent.com/scylladb/scylla-operator/v1.19/deploy/manager-dev.yaml

kubectl -n=scylla-manager rollout status --timeout=10m deployment.apps/scylla-manager

# Set affinity and fix space for KinD
kubectl label nodes kind-control-plane scylla.scylladb.com/node-type=scylla
```
