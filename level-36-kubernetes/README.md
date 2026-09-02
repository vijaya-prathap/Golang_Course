# Level 36: Kubernetes - Complete Guide

> ⚠️ **Environment note:** This level's manifests could not be applied to or tested against a real Kubernetes cluster in the environment that authored this course, because no cluster (and possibly no `kubectl`) is available here. The `kubectl` CLI binary itself was present, but even `kubectl apply --dry-run=client` requires contacting a live API server (or a locally cached OpenAPI schema) for discovery in this environment, so not even client-side dry-run validation was possible. Every manifest below is written correctly per the current stable Kubernetes API and was checked for well-formed YAML syntax, but — unlike every Go-code level in this course — none of it was verified against a live, running cluster. When you run these yourself with a cluster (even a local `kind` or `minikube` one), trust your own real output over anything shown here.

## Introduction

Welcome to Level 36! In Level 35 (Docker) you learned to package a single Go application into a container image: a `Dockerfile`, a multi-stage build, and `docker-compose` to run one or two containers together on your own machine.

That works great until you need more than one machine, or your one container crashes at 3 AM and nobody restarts it, or traffic spikes and you need five copies of your app instead of one. **Kubernetes** (often shortened to "K8s") is the system that takes the container image you built in Level 35 and runs it reliably across a whole fleet of machines - restarting it when it dies, scaling it up and down, rolling out new versions without downtime, and load-balancing traffic across every running copy.

This level is different from every other level in this course: there is no Go code to write. Kubernetes is configured almost entirely through YAML manifests that describe the *desired state* of your system - "I want 3 copies of this container running" - and Kubernetes' job is to continuously make reality match that description. Difficulty: ⭐⭐⭐.

---

## Table of Contents

1. [Why Kubernetes](#why-kubernetes)
2. [Core Building Blocks: Pod, Deployment, Service](#core-building-blocks-pod-deployment-service)
3. [A Deployment Manifest for a Go App](#a-deployment-manifest-for-a-go-app)
4. [A Service Manifest to Expose the Deployment](#a-service-manifest-to-expose-the-deployment)
5. [ConfigMaps and Secrets](#configmaps-and-secrets)
6. [Readiness and Liveness Probes](#readiness-and-liveness-probes)
7. [Resource Requests and Limits](#resource-requests-and-limits)
8. [kubectl Basics](#kubectl-basics)
9. [Rolling Updates and Rollbacks](#rolling-updates-and-rollbacks)
10. [Best Practices](#best-practices)
11. [Common Mistakes](#common-mistakes)

---

## Why Kubernetes

In Level 35, `docker-compose up` was enough because everything ran on one machine, under one person's supervision. That breaks down at real-world scale:

- **A container crashes.** On a single machine with plain Docker, it stays dead until a human (or a cron job) notices and restarts it. Kubernetes watches every container it manages and restarts it automatically.
- **Traffic grows.** One container handling all requests becomes a bottleneck. You want 3, 10, or 50 identical copies (**replicas**) sharing the load, and you want to change that number in one command, not by SSH-ing into servers.
- **You have many machines, not one.** A real production system runs across a cluster of servers (called **nodes**). Something has to decide which node runs which container, move containers off a node that dies, and let you treat the whole cluster as one resource pool instead of managing each machine by hand.
- **You need to deploy a new version without downtime.** Stopping all old containers and starting new ones creates a gap where nothing serves traffic. Kubernetes can replace them gradually - a few at a time - so the service never fully drops.
- **Traffic needs to reach the *right* containers, wherever they are.** Replicas get created and destroyed, and they get new internal IP addresses each time. Something needs to give clients (or other services) one stable address that always routes to whichever replicas are currently healthy.

Kubernetes is the industry-standard answer to all five problems at once: **container orchestration** - automatically running, healing, scaling, updating, and load-balancing containers across a cluster of machines. It builds directly on what Level 35 gave you: Kubernetes doesn't replace Docker images, it runs them, at scale, with self-healing built in.

---

## Core Building Blocks: Pod, Deployment, Service

Three objects cover almost everything you need to get a Go application running in a cluster. They stack on top of each other:

```
Service  (stable network address)
   │
   │ routes traffic to
   ▼
Deployment  (manages a set of replicas)
   │
   │ creates and supervises
   ▼
Pod  (one or more containers)  Pod  (one or more containers)  Pod  ...
```

### Pod: the smallest deployable unit

A **Pod** wraps one or more containers that are always scheduled together, on the same node, sharing the same network namespace (localhost) and, optionally, storage. Almost always, in this course, a Pod holds exactly **one** container - your Go application's image from Level 35. You will rarely create bare Pods directly in real work (see Deployment below), but every other object ultimately manages Pods, so you must understand them first.

### Deployment: manages a set of identical Pod replicas

You almost never create Pods by hand. Instead you create a **Deployment**, and you tell it two things: "run this container image" and "keep exactly N copies of it running." The Deployment then:

- Creates N Pods matching that template.
- Watches them continuously - if one crashes or its node dies, the Deployment creates a replacement automatically (**self-healing**).
- Handles **rolling updates**: change the image tag in the Deployment, and it replaces old Pods with new ones gradually, keeping the service available the whole time (covered in [Section 9](#rolling-updates-and-rollbacks)).
- Lets you scale by changing one number (`replicas: 3` → `replicas: 10`).

### Service: a stable network endpoint in front of a changing set of Pods

Pods are disposable - a Deployment can kill and recreate them at any time, and each new Pod gets a brand-new internal IP address. Nothing that talks to your application should ever hard-code a Pod's IP. A **Service** solves this: it gives you one stable name and IP address that automatically load-balances traffic across whichever Pods currently match its label selector, regardless of how many times those Pods have been replaced.

### How they relate

1. You write a **Deployment** describing your Go app's container and how many replicas you want.
2. The Deployment creates and supervises the matching **Pods**.
3. You write a **Service** that selects those same Pods (by label) and exposes them under one stable address.
4. Clients (or other services in the cluster) talk to the Service; the Service quietly load-balances across however many Pods are currently healthy.

---

## A Deployment Manifest for a Go App

Here is a complete, correct Deployment for a Go application whose image you built in Level 35 (`books-api:1.0.0`), running 3 replicas:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: books-api
  labels:
    app: books-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: books-api
  template:
    metadata:
      labels:
        app: books-api
    spec:
      containers:
        - name: books-api
          image: books-api:1.0.0
          ports:
            - containerPort: 8080
          resources:
            requests:
              cpu: "100m"
              memory: "64Mi"
            limits:
              cpu: "500m"
              memory: "256Mi"
```

### Anatomy

- **`apiVersion: apps/v1`** - Deployments live in the `apps/v1` API group. This is stable and has been the correct value for years; you will not see `apps/v1beta1` or similar in current Kubernetes.
- **`kind: Deployment`** - the type of object being described.
- **`metadata.name`** - the Deployment's own name (`books-api`). Kubernetes objects are identified by name within a namespace.
- **`metadata.labels`** - arbitrary key/value tags used for organizing and selecting objects. `app: books-api` here is just a label on the Deployment object itself.
- **`spec.replicas: 3`** - the desired number of identical Pods. Change this one number to scale.
- **`spec.selector.matchLabels`** - **critical**: this tells the Deployment which Pods belong to it. It **must** match `spec.template.metadata.labels` exactly, or Kubernetes rejects the manifest.
- **`spec.template`** - the **Pod template**. Everything under here describes the Pod that gets stamped out `replicas` times. Note it has its own nested `metadata` and `spec` - a Pod template is a full Pod spec embedded inside the Deployment.
- **`spec.template.spec.containers[].image`** - the exact image (and tag) built and pushed in Level 35. Referencing `books-api:1.0.0` here is exactly like `docker run books-api:1.0.0`, except Kubernetes runs it 3 times and keeps it running.
- **`resources.requests` / `resources.limits`** - covered fully in [Section 7](#resource-requests-and-limits); shown here because a real Deployment manifest should always set them.

---

## A Service Manifest to Expose the Deployment

A Deployment alone gives you 3 running Pods, each with its own internal, unstable IP. A **Service** gives you one stable way to reach them.

### Service Types, Briefly Compared

| Type | Reachable from | Typical use |
|------|-----------------|-------------|
| `ClusterIP` (default) | Inside the cluster only | Internal traffic between services (e.g., a frontend calling a backend API) |
| `NodePort` | Any cluster node's IP, on a fixed high port (30000-32767) | Simple external access, local development, demos |
| `LoadBalancer` | The public internet, via a cloud provider's load balancer | Production external access on a cloud platform (AWS/GCP/Azure provision a real load balancer for you) |

### A Complete ClusterIP Example

```yaml
apiVersion: v1
kind: Service
metadata:
  name: books-api
spec:
  type: ClusterIP
  selector:
    app: books-api
  ports:
    - port: 80
      targetPort: 8080
      protocol: TCP
```

### Anatomy

- **`apiVersion: v1`** - Services live in the core `v1` API group (same as Pods, ConfigMaps, and Secrets).
- **`spec.type: ClusterIP`** - internal-only; this is also the default if `type` is omitted entirely.
- **`spec.selector`** - matches Pods by label, exactly like a Deployment's selector. Any Pod with the label `app: books-api` - regardless of which Deployment created it - receives traffic from this Service.
- **`spec.ports[].port`** - the port the *Service* listens on (what clients connect to: `books-api:80`).
- **`spec.ports[].targetPort`** - the port the *container* is actually listening on (`8080`, matching `containerPort` in the Deployment). `port` and `targetPort` do not have to match, and frequently don't - it's common to expose `80` externally while the app listens on `8080` internally.

A `NodePort` variant is identical except for `type: NodePort` and an optional `nodePort` field - see Exercise 4.

---

## ConfigMaps and Secrets

Level 34 (Configuration) covered externalizing configuration out of your Go source code - environment variables, flags, config files - instead of hardcoding values. Kubernetes gives you two cluster-native objects for exactly that: **ConfigMaps** for non-sensitive configuration, and **Secrets** for sensitive values. Neither is baked into the image built in Level 35; both are injected into the Pod at runtime, so the same image can run identically in every environment with different configuration attached.

### ConfigMap Example

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: books-api-config
data:
  LOG_LEVEL: "info"
  MAX_BOOKS_PER_PAGE: "20"
```

### Secret Example

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: books-api-secret
type: Opaque
stringData:
  DB_PASSWORD: "s3cr3t-value"
```

A few important details about Secrets:

- **`type: Opaque`** means "arbitrary user-defined data" - the most common Secret type. Kubernetes also has built-in types for TLS certificates, Docker registry credentials, and more.
- **`stringData`** lets you write plain-text values in the manifest for readability; Kubernetes stores them internally as `data`, which is **base64-encoded**, not encrypted. Base64 is an encoding, not encryption - anyone who can read the Secret object can trivially decode it. Treat Secret manifests themselves as sensitive files (don't commit real secret values to version control), and rely on your cluster's RBAC and, ideally, encryption-at-rest for real protection.

### Referencing Them From a Pod as Environment Variables

```yaml
      containers:
        - name: books-api
          image: books-api:1.0.0
          env:
            - name: LOG_LEVEL
              valueFrom:
                configMapKeyRef:
                  name: books-api-config
                  key: LOG_LEVEL
            - name: DB_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: books-api-secret
                  key: DB_PASSWORD
```

Inside the Go program, this is retrieved exactly the way Level 34 taught: `os.Getenv("LOG_LEVEL")`, `os.Getenv("DB_PASSWORD")`. The application code never knows or cares whether the value came from a ConfigMap, a Secret, or a plain shell environment variable - Kubernetes just populates the process environment before your `main()` runs.

---

## Readiness and Liveness Probes

A running container is not automatically a *working* container. Kubernetes needs two different questions answered continuously, and it asks them with two different probes:

- **Liveness probe** - "Is this container still alive, or has it deadlocked/hung and needs to be killed and restarted?" If the liveness probe fails, Kubernetes kills the container and starts a fresh one.
- **Readiness probe** - "Is this container ready to receive traffic *right now*?" If the readiness probe fails, Kubernetes stops sending it traffic through the Service (without killing it) until it passes again.

The distinction matters most at startup: a Go app that needs a second or two to connect to a database is *alive* the instant the process starts, but not yet *ready*. Without a readiness probe, the Service would route real user requests to a Pod before it can handle them, causing errors during every rollout or restart.

This bridges directly to Level 27 (HTTP & REST APIs): add a simple health endpoint to your Go server, and probe it with HTTP.

```yaml
      containers:
        - name: books-api
          image: books-api:1.0.0
          ports:
            - containerPort: 8080
          readinessProbe:
            httpGet:
              path: /ready
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 10
          livenessProbe:
            httpGet:
              path: /healthz
              port: 8080
            initialDelaySeconds: 15
            periodSeconds: 20
```

### Anatomy

- **`httpGet.path` / `httpGet.port`** - Kubernetes makes an HTTP GET request to this path and port on the Pod's own IP. Any response in the `200`-`399` range counts as success; anything else (or a timeout/connection refused) counts as failure. A minimal Go handler is enough: `mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })`.
- **`initialDelaySeconds`** - how long to wait after the container starts before the first probe. Liveness typically gets a longer delay than readiness, so a slow-starting app isn't killed before it's had a fair chance to come up.
- **`periodSeconds`** - how often to repeat the probe after that.
- It is common (and often sufficient) to point both probes at the same `/healthz` endpoint. Splitting them into `/ready` (checks dependencies like a database connection) and `/healthz` (checks only that the process itself is responsive) is the more precise pattern shown here, and matters once your readiness condition is more than "the process is up."

---

## Resource Requests and Limits

Every container in a Pod can (and, in a real cluster, should) declare how much CPU and memory it needs:

```yaml
          resources:
            requests:
              cpu: "100m"
              memory: "64Mi"
            limits:
              cpu: "500m"
              memory: "256Mi"
```

- **`requests`** - what the container is **guaranteed** to get, and what the scheduler uses to decide which node has room for this Pod. A Pod is only scheduled onto a node that has at least this much CPU/memory unclaimed.
- **`limits`** - the **hard cap**. If a container tries to use more memory than its limit, it is killed (an "OOMKilled" event) and restarted. If it tries to use more CPU than its limit, it isn't killed - it's throttled instead.
- **Units:** CPU is measured in cores, where `"1"` is one full core and `"100m"` means "100 millicores" - one-tenth of a core. Memory uses `Mi` (mebibytes) or `Gi` (gibibytes) - `"256Mi"` is 256 mebibytes.

Setting both matters for two independent reasons: **scheduling** (requests tell Kubernetes how tightly it can safely pack Pods onto nodes) and **stability** (limits stop one misbehaving Pod - a memory leak, a runaway loop - from starving every other Pod on the same node). A Pod with no limits at all can consume an entire node's memory and take down unrelated workloads sharing that node.

---

## kubectl Basics

`kubectl` is the command-line tool for talking to a Kubernetes cluster's API server. These five commands cover the vast majority of day-to-day work. They are explained precisely here even though this authoring environment has no live cluster to run them against - the syntax and behavior described are standard, stable `kubectl` behavior.

| Command | Purpose |
|---------|---------|
| `kubectl apply -f <file>.yaml` | **Declarative**: create the objects described in the file if they don't exist, or update them to match the file if they do. This is the command you use almost always. |
| `kubectl get <kind>` | List objects of a given kind (e.g. `kubectl get pods`, `kubectl get deployments`, `kubectl get services`). Add `-o wide` for more columns, or a specific name (`kubectl get pod books-api-xyz`) for one object. |
| `kubectl describe <kind> <name>` | Show full details of one object, including recent events - the first place to look when something isn't working (e.g. a Pod stuck in `Pending` or `CrashLoopBackOff`). |
| `kubectl logs <pod-name>` | Stream or dump the stdout/stderr of a container in a Pod - your Go program's `log.Println` and `fmt.Println` output ends up here. Add `-f` to follow it live, or `-c <container>` if the Pod has more than one container. |
| `kubectl delete -f <file>.yaml` (or `kubectl delete <kind> <name>`) | Remove the object(s). |

### apply Is Declarative, Not Imperative

`kubectl apply -f deployment.yaml` does not mean "run this command once." It means "make the cluster's actual state match this file, whatever that requires" - creating objects that don't exist yet, patching fields that changed, and leaving untouched fields alone. Running the exact same `apply` command twice in a row is safe and does nothing the second time, because the cluster already matches the file. This is fundamentally different from an imperative command like `kubectl create deployment books-api --image=books-api:1.0.0`, which describes an *action* to take once and fails if you run it again (the object already exists). Real-world Kubernetes workflows are almost entirely `apply`-based, with YAML files as the source of truth, often checked into version control.

---

## Rolling Updates and Rollbacks

When you change the container image (or almost any Pod-template field) in a Deployment and re-`apply` it, Kubernetes does **not** stop all 3 old Pods and start 3 new ones at once - that would cause a visible gap in service. Instead, by default, it performs a **rolling update**:

1. Kubernetes creates a new Pod running the new image.
2. Once that new Pod passes its **readiness probe** (Section 6), Kubernetes considers it able to take traffic.
3. Kubernetes then terminates one old Pod.
4. Repeat, a few Pods at a time, until every old Pod has been replaced.

Throughout the process, the Service keeps routing traffic only to Pods that are currently ready - so `replicas: 3` never actually drops to zero available Pods, even mid-rollout. The exact pace is controlled by `strategy.rollingUpdate.maxSurge` / `maxUnavailable` (see Bonus Challenge 3).

To trigger a rollout, you typically change the image tag and re-apply:

```bash
kubectl set image deployment/books-api books-api=books-api:1.1.0
# or: edit the image field in the YAML file and run kubectl apply -f deployment.yaml again
```

You can watch it happen with:

```bash
kubectl rollout status deployment/books-api
```

And if the new version turns out to be broken, revert to the previous version with a single command - no need to remember or re-apply the old YAML by hand:

```bash
kubectl rollout undo deployment/books-api
```

Kubernetes keeps a revision history of a Deployment's Pod templates specifically to make this possible; `kubectl rollout history deployment/books-api` lists past revisions, and `kubectl rollout undo deployment/books-api --to-revision=<N>` can jump back further than one step.

---

## Best Practices

### 1. Always Set Resource Requests and Limits

```yaml
# ✅ Good
resources:
  requests:
    cpu: "100m"
    memory: "64Mi"
  limits:
    cpu: "500m"
    memory: "256Mi"

# ❌ Risky - no requests/limits at all
resources: {}
```

### 2. Always Define Readiness and Liveness Probes

A container with no readiness probe is assumed ready the instant it starts - even if it's still connecting to a database. Define both probes on anything that serves traffic.

### 3. Never Store Secrets in a Plain ConfigMap

```yaml
# ❌ WRONG - password sitting in plain text in a ConfigMap
apiVersion: v1
kind: ConfigMap
metadata:
  name: books-api-config
data:
  DB_PASSWORD: "s3cr3t-value"

# ✅ RIGHT - use a Secret instead
apiVersion: v1
kind: Secret
metadata:
  name: books-api-secret
type: Opaque
stringData:
  DB_PASSWORD: "s3cr3t-value"
```

ConfigMaps are not encrypted or access-restricted any differently from other ordinary objects; Secrets at least signal intent and get special handling (e.g., not shown in plain text by some tooling) and can be integrated with external secret managers.

### 4. Pin Image Tags - Never Use `latest`

```yaml
# ✅ Good - a specific, reproducible version
image: books-api:1.0.0

# ❌ WRONG - "latest" means a different image every time it's pulled,
# and rolling back becomes guesswork
image: books-api:latest
```

### 5. Use Namespaces to Separate Environments

```bash
kubectl create namespace staging
kubectl create namespace production
kubectl apply -f deployment.yaml -n staging
```

Namespaces give you separate, isolated pools of objects (a `staging` Deployment named `books-api` doesn't collide with a `production` one of the same name), and let you apply different resource quotas and access controls per environment.

---

## Common Mistakes

### Mistake 1: No Resource Limits, So One Pod Starves Its Neighbors

```yaml
# ❌ WRONG - a memory leak in this container can consume the
# entire node, starving every other Pod scheduled on it
resources: {}

# ✅ RIGHT - a hard cap contains the damage to this one container
resources:
  limits:
    memory: "256Mi"
```

### Mistake 2: Missing Readiness Probe Sends Traffic to a Pod Too Early

```yaml
# ❌ WRONG - no readinessProbe means "ready" the instant the process starts,
# even if it hasn't finished connecting to its database yet
containers:
  - name: books-api
    image: books-api:1.0.0

# ✅ RIGHT
containers:
  - name: books-api
    image: books-api:1.0.0
    readinessProbe:
      httpGet:
        path: /ready
        port: 8080
```

### Mistake 3: Hardcoding Config That Should Be a ConfigMap/Secret

```yaml
# ❌ WRONG - config value baked directly into the manifest with no
# ability to change per environment without editing the Deployment
env:
  - name: LOG_LEVEL
    value: "info"

# ✅ RIGHT - externalized, so staging and production can each supply
# a different value without touching the Deployment spec
env:
  - name: LOG_LEVEL
    valueFrom:
      configMapKeyRef:
        name: books-api-config
        key: LOG_LEVEL
```

### Mistake 4: Confusing Declarative `apply` With One-Off Imperative Commands

```bash
# ❌ Fragile - imperative, fails if the object already exists,
# and the exact state isn't captured anywhere reusable
kubectl create deployment books-api --image=books-api:1.0.0

# ✅ RIGHT - declarative, safe to re-run, and the YAML file is the
# single source of truth you can check into version control
kubectl apply -f deployment.yaml
```

---

## Summary

**Core Objects:**
- `Pod` - one or more containers, the smallest deployable unit
- `Deployment` - manages N identical Pod replicas; handles self-healing and rolling updates
- `Service` - a stable network endpoint load-balancing across a changing set of Pods

**Configuration:**
- `ConfigMap` - non-sensitive config, injected as environment variables (or files)
- `Secret` - sensitive config, same injection mechanism, base64-encoded (not encrypted) at rest

**Health:**
- `readinessProbe` - is this Pod ready for traffic right now?
- `livenessProbe` - is this Pod still alive, or does it need restarting?

**Resources:**
- `requests` - guaranteed minimum, used for scheduling
- `limits` - hard cap, enforced at runtime

**Workflow:**
- `kubectl apply -f` is declarative and safe to re-run
- Rolling updates replace Pods gradually; `kubectl rollout undo` reverts a bad release

---

## Next Steps

You now understand:
- ✅ Why Kubernetes exists and what problems it solves beyond single-machine Docker
- ✅ Pod, Deployment, and Service, and how they relate
- ✅ Writing Deployment and Service manifests for a Go application
- ✅ Externalizing configuration with ConfigMaps and Secrets
- ✅ Readiness and liveness probes, and resource requests/limits
- ✅ Core kubectl commands and the declarative `apply` model
- ✅ How rolling updates and `kubectl rollout undo` work

**Next level:** Level 37 - CI/CD
- Automating the build → test → containerize → deploy pipeline
- Connecting the Docker image from Level 35 and the manifests from this level into an automated workflow
- Continuous integration and continuous delivery concepts

You're now equipped to describe, deploy, and reason about a real production topology! Keep going! 🚀
