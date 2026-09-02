# Level 36: Study Guide & Visual Reference

> ⚠️ **Caveat:** No live Kubernetes cluster was available while authoring this level, and even `kubectl apply --dry-run=client` could not run without cluster connectivity in this environment. Manifests here are correct per the stable Kubernetes API and passed a plain-YAML well-formedness check, but were not verified against a real cluster. Trust your own `kubectl` output over anything shown here.

## 📚 Learning Path

### Week 1: Core Objects
```
Day 1:  Why Kubernetes - the problems it solves beyond single-machine Docker
Day 2:  Pod, Deployment, Service and how they relate
Day 3:  Writing a Deployment manifest for a Go app
Day 4:  Writing ClusterIP and NodePort Service manifests
Day 5:  ConfigMaps and Secrets
Day 6:  Readiness and liveness probes
Day 7:  Resource requests and limits
```

### Week 2: Operations & Practice
```
Day 1:  kubectl basics (apply, get, describe, logs, delete)
Day 2:  Rolling updates and kubectl rollout undo
Day 3:  Exercises 1-4 (Pod, Deployment, ClusterIP, NodePort)
Day 4:  Exercises 5-6 (ConfigMap, Secret)
Day 5:  Exercises 7-8 (probes, resources)
Day 6:  Exercises 9-10 (kubectl walkthrough, comprehensive manifest set)
Day 7:  Bonus challenges & review
```

---

## 🗺️ Pod / Deployment / Service Relationship

```
                          ┌─────────────────────────┐
  Client / other Pod ───▶ │   Service: books-api     │   stable name + IP,
                          │   (ClusterIP / NodePort  │   load-balances across
                          │    / LoadBalancer)       │   whatever Pods match
                          └────────────┬─────────────┘   its label selector
                                       │  selector: app=books-api
                    ┌──────────────────┼──────────────────┐
                    ▼                  ▼                  ▼
            ┌───────────────┐ ┌───────────────┐ ┌───────────────┐
            │ Pod (replica) │ │ Pod (replica) │ │ Pod (replica) │
            │ app=books-api │ │ app=books-api │ │ app=books-api │  ◀── created &
            │ books-api:1.0 │ │ books-api:1.0 │ │ books-api:1.0 │      supervised by
            └───────────────┘ └───────────────┘ └───────────────┘
                    ▲                  ▲                  ▲
                    └──────────────────┴──────────────────┘
                                       │
                          ┌────────────┴─────────────┐
                          │  Deployment: books-api    │  desired state:
                          │  replicas: 3              │  "always keep 3 Pods
                          └───────────────────────────┘   matching this template
                                                           running, self-heal,
                                                           roll out updates"
```

**Read it top to bottom:** a Deployment's only job is keeping N Pods matching its template alive. A Service's only job is giving those (disposable, IP-changing) Pods one stable address. Neither object talks to the other directly - they connect purely through matching labels (`app: books-api`).

---

## 🔑 ConfigMap / Secret → Pod Injection

```
┌─────────────────────┐     ┌─────────────────────┐
│  ConfigMap           │     │  Secret              │
│  books-api-config    │     │  books-api-secret    │
│  ─────────────────   │     │  ─────────────────   │
│  LOG_LEVEL: "info"    │     │  DB_PASSWORD: (b64)  │
│  MAX_BOOKS_PER_PAGE   │     │  type: Opaque         │
└──────────┬───────────┘     └──────────┬───────────┘
           │ configMapKeyRef             │ secretKeyRef
           ▼                             ▼
        ┌─────────────────────────────────────┐
        │  Pod: books-api                       │
        │  containers:                          │
        │    - env:                             │
        │        LOG_LEVEL         (from CM)    │
        │        MAX_BOOKS_PER_PAGE (from CM)   │
        │        DB_PASSWORD       (from Secret)│
        └──────────────────┬────────────────────┘
                            ▼
                  Go process environment
                  os.Getenv("LOG_LEVEL")   ◀── Level 34 pattern,
                  os.Getenv("DB_PASSWORD")     unchanged in the container
```

Neither value is baked into the `books-api:1.0.0` image built in Level 35 - the same image runs unmodified in every environment, with a different ConfigMap/Secret supplying different values per cluster or namespace.

---

## ⚖️ Readiness vs Liveness Probes

| | Readiness Probe | Liveness Probe |
|---|---|---|
| **Question it answers** | "Can this Pod take traffic *right now*?" | "Is this Pod alive, or hung/deadlocked?" |
| **On failure** | Removed from the Service's endpoints (no traffic), **not** restarted | Container is killed and restarted |
| **Typical check** | Dependencies ready (DB connection, cache warm) | Process itself still responds at all |
| **Typical `initialDelaySeconds`** | Short (app-startup dependent, e.g. `5`) | Longer, gives startup time before risking a kill (e.g. `15`) |
| **Example path** | `/ready` | `/healthz` |
| **Consequence if missing** | Traffic hits a Pod before it can handle it (errors during rollout/restart) | A hung Pod stays `Running` forever, silently serving nothing |

```
Pod lifecycle timeline:

  t=0s   container starts
  t=5s   ─┬─ readinessProbe first check (initialDelaySeconds: 5)
          └─ fails until DB connection is up → Pod excluded from Service
  t=8s   readinessProbe passes → Pod added to Service, starts receiving traffic
  t=15s  ─┬─ livenessProbe first check (initialDelaySeconds: 15)
          └─ passes → nothing happens (this is the common case)
  ...    both probes repeat every periodSeconds indefinitely
```

---

## 🌐 Service Type Comparison

| Type | Reachable From | Gets a Fixed Cluster-Internal IP? | External Access | Typical Use |
|------|-----------------|-----------------------------------|------------------|-------------|
| `ClusterIP` (default) | Inside the cluster only | Yes | None | Service-to-service traffic (e.g. frontend → backend) |
| `NodePort` | Any node's IP on a fixed high port | Yes (plus the port) | Yes, via `<node-ip>:<nodePort>` | Local dev (`kind`/`minikube`), quick demos |
| `LoadBalancer` | The public internet | Yes | Yes, via a cloud-provisioned load balancer's IP/hostname | Production external access on a cloud platform |

```
ClusterIP:     [Pod A] ──▶ [Service:80] ──▶ {Pod X, Pod Y, Pod Z}   (internal only)

NodePort:      [Outside] ──▶ [Any Node:30080] ──▶ [Service:80] ──▶ {Pod X, Pod Y, Pod Z}

LoadBalancer:  [Internet] ──▶ [Cloud LB] ──▶ [Service:80] ──▶ {Pod X, Pod Y, Pod Z}
```

Every `LoadBalancer` Service is also a `NodePort` Service is also a `ClusterIP` Service under the hood - each type is a superset of the one before it.

---

## ⌨️ kubectl Command Reference

| Command | What It Does | Declarative or Imperative? |
|---------|---------------|------------------------------|
| `kubectl apply -f <file>` | Create or update objects to match the file; safe to re-run | Declarative |
| `kubectl get <kind>` | List objects of a kind (`pods`, `deployments`, `services`, ...) | Read-only |
| `kubectl get <kind> <name> -o yaml` | Dump one object's full current state as YAML | Read-only |
| `kubectl describe <kind> <name>` | Full detail + recent Events for one object (first stop when debugging) | Read-only |
| `kubectl logs <pod>` | Container stdout/stderr (add `-f` to follow, `-c` for multi-container Pods) | Read-only |
| `kubectl delete -f <file>` | Remove everything the file describes | Declarative |
| `kubectl set image deployment/<name> <container>=<image>` | Trigger a rolling update to a new image | Imperative (one-off action) |
| `kubectl rollout status deployment/<name>` | Watch a rollout in progress | Read-only |
| `kubectl rollout undo deployment/<name>` | Revert to the previous Deployment revision | Imperative (one-off action) |
| `kubectl exec <pod> -- <cmd>` | Run a command inside a running container | Imperative (one-off action) |

---

## 🚨 Common Mistakes

### Mistake 1: No Resource Limits
```yaml
# ❌ resources: {}  → one leaking Pod can starve its whole node
# ✅ always set requests AND limits
```

### Mistake 2: Missing Readiness Probe
```yaml
# ❌ no readinessProbe → traffic hits the Pod before it's ready
# ✅ readinessProbe.httpGet checking a real dependency-aware endpoint
```

### Mistake 3: Config Hardcoded Instead of Externalized
```yaml
# ❌ env: [{name: LOG_LEVEL, value: "info"}]   (fixed forever)
# ✅ env: [{name: LOG_LEVEL, valueFrom: {configMapKeyRef: {...}}}]
```

### Mistake 4: Treating apply Like an Imperative One-Off
```
# ❌ assuming kubectl apply only "does something" the first time
# ✅ apply is declarative - it reconciles current state to match the file,
#    every single time you run it
```

---

## 📈 Progression Summary

### Understanding Level 36

Level 36 teaches how a single Docker image (Level 35) becomes a resilient, scalable, load-balanced service:

1. **Pod → Deployment → Service** - the three-object chain that runs and exposes an app
2. **ConfigMap/Secret** - injecting Level 34-style configuration at the cluster level instead of the image
3. **Probes** - telling Kubernetes when a Pod (built on a Level 27 HTTP server) is ready vs alive
4. **Resources** - requests for scheduling, limits for stability
5. **kubectl + rollouts** - the declarative workflow for applying, inspecting, and safely updating all of the above

### Prerequisites for Level 37

Before moving to Level 37 (CI/CD), you need:

- ✅ Comfortable writing a Deployment and a Service manifest from scratch
- ✅ Understand how ConfigMaps and Secrets get injected into a Pod's environment
- ✅ Understand why both readiness and liveness probes matter, and how they differ
- ✅ Understand resource requests vs limits
- ✅ Comfortable with `kubectl apply`, `get`, `describe`, `logs`, `delete`
- ✅ Understand what a rolling update does and how `kubectl rollout undo` works

### Ready for Level 37?

Level 37 teaches automating everything you just did by hand:
- Building and testing the Go app automatically on every commit
- Building and pushing the Level 35 Docker image automatically
- Applying the Level 36 manifests to a cluster automatically as part of a pipeline

---

## ✅ Checklist Before Level 37

- [ ] Can write a Pod, Deployment, and Service manifest without a reference
- [ ] Can explain the difference between ClusterIP, NodePort, and LoadBalancer
- [ ] Can wire a ConfigMap and a Secret into a container's environment variables
- [ ] Can explain readiness vs liveness probes and why both matter
- [ ] Can explain requests vs limits for CPU and memory
- [ ] Comfortable with `apply`, `get`, `describe`, `logs`, `delete`
- [ ] Can explain what a rolling update does and how to undo one
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The Chain
Service → Deployment → Pods. Each layer only knows about the one below it through label selectors, not direct references.

### The Externalization
Nothing environment-specific belongs in the container image. ConfigMaps and Secrets carry it in at runtime, exactly like Level 34's philosophy, now applied at the cluster level.

### The Health Contract
A Pod that doesn't declare readiness/liveness probes is assumed healthy the instant it starts - which is rarely true. Bridge from Level 27: your HTTP server needs real health endpoints.

### The Declarative Model
`kubectl apply -f` describes desired state, not a one-time action. Re-running it is always safe.

---

## 📚 Next Level

Level 37: CI/CD
- Automating build → test → containerize → deploy
- Wiring Level 35's Docker image and this level's manifests into a pipeline
- Continuous integration and continuous delivery concepts

You've got Kubernetes fundamentals down! Keep going! 🚀
