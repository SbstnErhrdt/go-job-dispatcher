# Kubernetes deployment

### 0. Create new namespace

```
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Namespace
metadata:
  name: go-job-dispatcher
EOF
```

### 1. Login for docker registry

```
kubectl -n go-job-dispatcher create secret docker-registry registry \
--docker-server=registry.url.com \
--docker-username=docker \
--docker-password=xxxx \
--docker-email=sebastian@url.com
```

### 2. Environment variables for pods

```
kubectl -n go-job-dispatcher create secret generic go-job-dispatcher-env-secret --from-env-file=production.env
```

### 3. Deploy resources

```
kubectl apply -f deployment.yml
```

### 4. Update resources

```
kubectl rollout restart deployment/go-job-dispatcher-deployment -n go-job-dispatcher
```