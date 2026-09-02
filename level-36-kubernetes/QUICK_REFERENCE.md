# Level 36: Quick Reference Card

> ⚠️ **Unverified against a live cluster:** No Kubernetes cluster was reachable while authoring this level, and even `kubectl apply --dry-run=client` requires cluster connectivity that wasn't available here. Manifests below are correct per the stable Kubernetes API but were checked only for plain YAML well-formedness, not applied or schema-validated. Trust your own `kubectl` output over anything shown here.

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir -p ~/k8s/myapp && cd ~/k8s/myapp

# Minimal Deployment + Service
cat > app.yaml << 'EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
        - name: myapp
          image: myapp:1.0.0
          ports:
            - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: myapp
spec:
  selector:
    app: myapp
  ports:
    - port: 80
      targetPort: 8080
EOF

# Apply
kubectl apply -f app.yaml
kubectl get pods -l app=myapp
```

---

## 📋 The Three Core Objects

```yaml
# Pod - one or more containers, smallest deployable unit
apiVersion: v1
kind: Pod

# Deployment - manages N identical Pod replicas, self-heals, handles rollouts
apiVersion: apps/v1
kind: Deployment

# Service - stable network endpoint load-balancing across matching Pods
apiVersion: v1
kind: Service
```

---

## 🌐 Service Types

```yaml
spec:
  type: ClusterIP     # internal only (default)
  type: NodePort      # internal + <node-ip>:<nodePort> (30000-32767)
  type: LoadBalancer  # internal + cloud-provisioned public LB
```

---

## 🔑 ConfigMap & Secret Injection

```yaml
# ConfigMap
apiVersion: v1
kind: ConfigMap
metadata: { name: myapp-config }
data:
  LOG_LEVEL: "info"

# Secret
apiVersion: v1
kind: Secret
metadata: { name: myapp-secret }
type: Opaque
stringData:
  DB_PASSWORD: "changeme"

# Referenced in a container:
env:
  - name: LOG_LEVEL
    valueFrom:
      configMapKeyRef: { name: myapp-config, key: LOG_LEVEL }
  - name: DB_PASSWORD
    valueFrom:
      secretKeyRef: { name: myapp-secret, key: DB_PASSWORD }
```

---

## 🩺 Probes

```yaml
readinessProbe:      # ready for traffic? (removed from Service if it fails)
  httpGet: { path: /ready, port: 8080 }
  initialDelaySeconds: 5
  periodSeconds: 10

livenessProbe:        # still alive? (container killed + restarted if it fails)
  httpGet: { path: /healthz, port: 8080 }
  initialDelaySeconds: 15
  periodSeconds: 20
```

---

## ⚖️ Resource Requests & Limits

```yaml
resources:
  requests:      # guaranteed minimum; used for scheduling
    cpu: "100m"    # 0.1 core
    memory: "64Mi"
  limits:        # hard cap; CPU throttled, memory OOMKilled if exceeded
    cpu: "500m"
    memory: "256Mi"
```

---

## ⌨️ kubectl Essentials

```bash
kubectl apply -f app.yaml          # create or update (declarative, safe to re-run)
kubectl get pods                   # list objects
kubectl get pods -o wide           # list with more detail
kubectl describe pod <name>        # full detail + Events (debug here first)
kubectl logs <pod>                 # container stdout/stderr
kubectl logs -f <pod>              # follow logs live
kubectl delete -f app.yaml         # remove everything in the file
```

---

## 🔄 Rolling Updates & Rollback

```bash
kubectl set image deployment/myapp myapp=myapp:1.1.0   # trigger a rollout
kubectl rollout status deployment/myapp                 # watch it progress
kubectl rollout undo deployment/myapp                    # revert to previous version
kubectl rollout history deployment/myapp                 # list past revisions
```

```yaml
# Rolling update pace (Bonus Challenge 3)
spec:
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1          # extra Pods allowed above `replicas` during rollout
      maxUnavailable: 0    # Pods allowed to be unavailable during rollout
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| No resource limits | `resources: {}` | Always set `requests` AND `limits` |
| Missing readiness probe | no `readinessProbe` | Add one backed by a real health endpoint |
| Secrets in a ConfigMap | `data: { DB_PASSWORD: "..." }` in a ConfigMap | Use a `Secret` with `stringData` |
| Using `latest` | `image: myapp:latest` | Pin an exact tag: `image: myapp:1.0.0` |
| Confusing apply with one-off commands | `kubectl create deployment ...` | `kubectl apply -f deployment.yaml` |

---

## 🎓 Before Next Level

Can you:
- [ ] Write a Pod, Deployment, and Service manifest from memory?
- [ ] Explain ClusterIP vs NodePort vs LoadBalancer?
- [ ] Wire a ConfigMap and Secret into a Pod's environment?
- [ ] Explain readiness vs liveness probes?
- [ ] Explain resource requests vs limits?
- [ ] Use `apply`, `get`, `describe`, `logs`, `delete` confidently?
- [ ] Explain what `kubectl rollout undo` does?

If YES → You're ready for Level 37!

---

## 📚 Next Level

Level 37: CI/CD
- Automating build → test → containerize → deploy
- Wiring Level 35's Docker image and this level's manifests into a pipeline

You've got Kubernetes fundamentals down! 💪
