# scale-operator

Kubernetes time based scale operator

## Commands

1. Create project

```bash
kubebuilder init --domain scaler.cloudsteak.com --owner "CloudSteak" --repo github.com/cloudsteak/scale-operator.git
```

2. Create API

```bash
kubebuilder create api --kind Scaler --version v1alpha1 --group api
```

3. Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.

```bash
make generate
```

4. Create manifests (CRD, RBAC and Controller)

```bash
make manifests
```

5. Develop your code

6. Install CRDs

```bash
make install
```

7. Create deployment

```bash
kubectl create deploy nginx --image=nginx
```

8. Configure scaler

```bash
nano config/samples/api_v1alpha1_scaler.yaml
```

```yaml
spec:
  start: 5 # UTC time
  end: 11 # UTC time
  replicas: 2
  deployments:
    - name: nginx
      namespace: nginx
```

9. Create scaler

```bash
kubectl apply -f config/samples/api_v1alpha1_scaler.yaml
```

10. Run the controller

```bash
make run
```

11. Check output

\*\* Scale up

```bash
2024-07-08T07:41:02+02:00       INFO    --- Scaling up deployments      {"controller": "scaler", "controllerGroup": "api.scaler.cloudsteak.com", "controllerKind": "Scaler", "Scaler": {"name":"scaler-sample","namespace":"default"}, "namespace": "default", "name": "scaler-sample", "reconcileID": "d113f2c7-e8d2-4d9a-a3f9-3b51134e7dd1"}
2024-07-08T07:41:02+02:00       INFO    Scaling Deployment      {"name": "nginx", "namespace": "default", "replicas_to": 2, "replicas_from": 1}
```

\*\* Scale down

```bash
2024-07-08T16:40:02+02:00       INFO    --- Scaling down deployments    {"controller": "scaler", "controllerGroup": "api.scaler.cloudsteak.com", "controllerKind": "Scaler", "Scaler": {"name":"scaler-sample","namespace":"default"}, "namespace": "default", "name": "scaler-sample", "reconcileID": "c0668843-f7ee-4c7a-99d6-7cf6e5fbeb0f"}
2024-07-08T16:40:02+02:00       INFO    Scaling Deployment      {"name": "nginx", "namespace": "default", "replicas_to": 1, "replicas_from": 2}
```
