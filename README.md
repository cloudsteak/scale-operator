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

7. Create scaler

```bash
kubectl apply -f config/samples/api_v1alpha1_scaler.yaml
```

8. Run the controller

```bash
make run
```
