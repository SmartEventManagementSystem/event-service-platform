# ArgoCD + Argo Rollouts Setup Guide

## Cách Setup trên ArgoCD UI

### Bước 1: Mở ArgoCD

```
http://localhost:9080
Username: admin
Password: fdPUgRtl0TK1yIhl
```

### Bước 2: Tạo Application mới

1. Click **"+ New App"** (góc trên phải)
2. Điền thông tin:

```
Application Name: event-service-rollout
Project: default
Sync Policy: Manual (hoặc Automatic)

SOURCE:
  Repository URL: https://github.com/iamhuutho/event-management-system.git
  Revision: HEAD
  Path: infra/kubernetes/argocd

DESTINATION:
  Cluster URL: https://kubernetes.default.svc
  Namespace: default
```

3. Click **"Create"**

### Bước 3: Sync Application

1. Click vào app **event-service-rollout**
2. Click **"Sync"** button
3. Chọn resources để sync
4. Click **"Synchronize"**

---

## Các Files trong `infra/kubernetes/argocd/`

| File | Mục đích |
|------|----------|
| `rollout.yaml` | Canary deployment config |
| `services.yaml` | Stable & Canary services |
| `experiment.yaml` | Test version mới trước khi deploy |
| `kustomization.yaml` | Kustomize build file |
| `argocd-application.yaml` | ArgoCD Application manifest |

---

## Files YAML

### 1. rollout.yaml - Canary Deployment

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Rollout
metadata:
  name: event-service
spec:
  replicas: 3
  strategy:
    canary:
      steps:
        - setWeight: 10    # 10% traffic sang version mới
        - pause: {}         # Dừng để review
        - setWeight: 30
        - pause: {duration: 5}  # Chờ 5 phút
        - setWeight: 100   # Full traffic
      canaryService: event-service-canary
      stableService: event-service-stable
```

### 2. services.yaml - Services

```yaml
# Stable Service - Version production
apiVersion: v1
kind: Service
metadata:
  name: event-service-stable
spec:
  selector:
    app: event-service
    role: stable

# Canary Service - Version mới
apiVersion: v1
kind: Service
metadata:
  name: event-service-canary
spec:
  selector:
    app: event-service
    role: canary
```

### 3. experiment.yaml - Experiment

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Experiment
metadata:
  name: event-service-test
spec:
  duration: 30m
  components:
    - replicas: 1
      name: event-service-new
      template:
        spec:
          containers:
            - name: event-service
              image: nginx:new-version
```

---

## Commands để test

### Apply trực tiếp (không qua Git)

```bash
# Apply Rollout
kubectl apply -f infra/kubernetes/argocd/rollout.yaml

# Xem trạng thái
kubectl argo rollouts get rollout event-service --watch

# Upgrade version
kubectl argo rollouts set image event-service event-service=nginx:1.21

# Pause/Resume/Rollback
kubectl argo rollouts pause event-service
kubectl argo rollouts resume event-service
kubectl argo rollouts abort event-service
```

### Experiment Commands

```bash
# Run experiment
kubectl apply -f infra/kubernetes/argocd/experiment.yaml

# Xem experiment
kubectl get experiment
kubectl describe experiment event-service-test
```

---

## Argo Rollouts Dashboard

```bash
# Mở dashboard
kubectl argo rollouts dashboard

# Truy cập
open http://localhost:3100
```

---

## Canary Flow

```
1. Deploy v1 (100% traffic)
2. Apply v2 → 10% traffic sang v2
3. ArgoCD pause → Chờ review
4. Manual approve → Tiếp tục
5. 30% → Pause
6. Review → Continue
7. 50% → 100%
8. v2 = 100% traffic
```

## Rollback Flow

```bash
# Nếu có vấn đề
kubectl argo rollouts abort event-service

# Hoặc undo về version trước
kubectl argo rollouts undo event-service
```

---

## Monitoring

```bash
# Xem rollout status
kubectl argo rollouts status event-service

# Xem history
kubectl argo rollouts history event-service

# Xem chi tiết
kubectl argo rollouts get rollout event-service -o yaml
```
