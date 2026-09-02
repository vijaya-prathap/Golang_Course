# Level 36: Kubernetes - INDEX

Welcome to **Level 36: Kubernetes**! This is where the single-container Docker world of Level 35 becomes a resilient, scalable, self-healing system running across a whole cluster of machines.

> ⚠️ **Environment note:** No live Kubernetes cluster (and no working `kubectl` cluster connection) was available in the environment that authored this level. Every manifest is written correctly per the current stable Kubernetes API and checked for well-formed YAML, but none of it was verified against a real cluster. See the banner at the top of README.md and EXERCISES.md for details.

---

## 📖 What You'll Learn

- ✅ Why Kubernetes exists - orchestrating containers across many machines
- ✅ Pod, Deployment, and Service, and how they relate
- ✅ Writing Deployment and Service manifests for a Go application
- ✅ ConfigMaps and Secrets for cluster-native configuration
- ✅ Readiness and liveness probes
- ✅ Resource requests and limits
- ✅ Core kubectl commands (apply, get, describe, logs, delete)
- ✅ Rolling updates and rollbacks

---

## 🗂️ Level 36 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview and welcome message. Start here first!

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- Why Kubernetes, and how it builds on Level 35's Docker foundation
- Pod, Deployment, Service and their relationship
- Deployment and Service manifests for a Go app
- ConfigMaps, Secrets, probes, and resource requests/limits
- kubectl basics and the declarative `apply` model
- Rolling updates and `kubectl rollout undo`
- Best practices
- Common mistakes

**Read Time:** 45-60 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 detailed exercises + 3 bonus challenges:
1. A basic Pod manifest
2. A Deployment with 3 replicas
3. A ClusterIP Service
4. A NodePort Service variant
5. A ConfigMap and a Pod that uses it
6. A Secret and a Pod that uses it
7. Adding readiness and liveness probes
8. Adding resource requests and limits
9. kubectl commands walkthrough
10. Comprehensive practice - a complete manifest set for books-api

**Time Commitment:** 4-6 hours
**When:** After reading README

### 4. **STUDY_GUIDE.md** (Visual Learning)
- Pod/Deployment/Service relationship diagram
- ConfigMap/Secret → Pod injection diagram
- Readiness vs liveness probe comparison table
- Service type comparison table
- kubectl command reference table
- Common mistakes guide

**Read Time:** 30-40 minutes
**When:** Use while doing exercises

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- Essential manifest syntax
- kubectl command essentials
- Rolling update / rollback commands
- Common mistakes table

**Read Time:** 10-15 minutes
**When:** Quick lookup during practice

---

## 🎯 Recommended Learning Path

### Day 1: Why Kubernetes & Core Objects (2 hours)
1. Read **START_HERE.txt** (10 min)
2. Read **README.md** sections 1-2 (30 min)
3. Complete Exercises 1-2 (1 hour 20 min)

### Day 2: Services & Configuration (2 hours)
1. Read **README.md** sections 3-5 (35 min)
2. Complete Exercises 3-6 (1.5 hours)

### Day 3: Health & Resources (1.5 hours)
1. Read **README.md** sections 6-7 (20 min)
2. Complete Exercises 7-8 (1 hour 10 min)

### Day 4: Operations & Real-World Practice (2 hours)
1. Read **README.md** sections 8-9 (25 min)
2. Complete Exercises 9-10 (1.5 hours)

### Day 5: Consolidation (1 hour)
1. Try bonus challenges
2. Review with **QUICK_REFERENCE.md**
3. Make sure the Pod → Deployment → Service chain feels automatic

---

## 💡 Key Concepts At A Glance

### Core Objects
```
Service  →  Deployment  →  Pods
(stable      (manages       (containers,
 address)     replicas)      smallest unit)
```

### A Minimal Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  replicas: 3
  selector: { matchLabels: { app: myapp } }
  template:
    metadata: { labels: { app: myapp } }
    spec:
      containers:
        - name: myapp
          image: myapp:1.0.0
```

### kubectl Essentials
```bash
kubectl apply -f app.yaml     # declarative create/update
kubectl get pods              # list
kubectl describe pod <name>   # full detail + events
kubectl logs <pod>            # stdout/stderr
kubectl delete -f app.yaml    # remove
```

---

## ✅ Prerequisites

Make sure you've completed **Level 35: Docker**

You need:
- ✅ Comfort building a Docker image from a Dockerfile
- ✅ Understanding of multi-stage builds
- ✅ Familiarity with docker-compose for running containers together
- ✅ From Level 27: a Go HTTP server with request handlers
- ✅ From Level 34: the concept of externalizing configuration via environment variables

---

## 🎓 Learning Objectives

By the end of Level 36, you'll be able to:

- ✅ Explain why Kubernetes is needed beyond single-machine Docker
- ✅ Write Pod, Deployment, and Service manifests for a Go application
- ✅ Compare ClusterIP, NodePort, and LoadBalancer Service types
- ✅ Externalize configuration into ConfigMaps and Secrets
- ✅ Add readiness and liveness probes backed by a real health endpoint
- ✅ Set resource requests and limits and explain why both matter
- ✅ Use kubectl's core commands and understand the declarative `apply` model
- ✅ Explain how rolling updates work and how to roll one back

---

## 📊 Statistics

- **Main Theory:** README.md covering core objects, manifests, configuration, health, resources, and operations
- **Exercises:** 10 detailed exercises + 3 bonus challenges
- **Study Guide:** relationship diagrams, probe/service comparison tables, kubectl reference
- **Quick Reference:** One-page cheat sheet
- **Practice Time:** 4-6 hours
- **Difficulty:** ⭐⭐⭐ (Advanced)

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for comprehensive theory
2. Use STUDY_GUIDE.md alongside for diagrams and comparison tables
3. Reference QUICK_REFERENCE.md for fast lookup

### For Practice
1. Read exercise description
2. Copy provided bash/YAML (or type them)
3. Follow step-by-step instructions
4. Compare against the labeled illustrative output (this level's manifests were not run against a live cluster - see the banner)
5. Understand what each field in the manifest does

### For Review
1. Use QUICK_REFERENCE.md for quick recap
2. Refer to the common mistakes section
3. Practice writing a Deployment + Service pair from memory

---

## 🆘 Common Questions

**Q: Do I need a real Kubernetes cluster to learn from this level?**
A: To truly verify behavior, yes - a local cluster like `kind` or `minikube` is free and sufficient. This course's own authoring environment did not have one, which is why every output in EXERCISES.md is labeled illustrative rather than verified.

**Q: Why don't I just create bare Pods instead of Deployments?**
A: A bare Pod that crashes stays dead - nothing recreates it. A Deployment supervises its Pods continuously and replaces any that disappear, plus gives you replica scaling and rolling updates for free.

**Q: When do I use NodePort vs LoadBalancer vs ClusterIP?**
A: ClusterIP for internal service-to-service traffic (the default, and most common inside a cluster). NodePort for simple local/demo external access. LoadBalancer for real production external access on a cloud provider.

**Q: Are Secrets actually secure just because they're a different object than ConfigMap?**
A: Not by themselves - Secret values are base64-encoded, not encrypted, by default. Real protection comes from your cluster's RBAC, encryption-at-rest, and not committing real secret values into manifests you check into version control.

**Q: What's the difference between kubectl apply and kubectl create?**
A: `apply` is declarative - it reconciles the cluster to match your file, and is safe to run repeatedly. `create` is imperative - it performs a one-time action and fails if the object already exists.

---

## 🎯 Before Moving to Level 37

Make sure you can answer these questions:

- [ ] What problem does a Deployment solve that a bare Pod doesn't?
- [ ] What's the difference between ClusterIP, NodePort, and LoadBalancer?
- [ ] How does a Pod get a ConfigMap or Secret value as an environment variable?
- [ ] What's the difference between a readiness probe and a liveness probe?
- [ ] What's the difference between a resource request and a resource limit?
- [ ] Why is `kubectl apply -f` safe to run more than once?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (10 min)
   - Gets you oriented

2. **README.md** Sections 1-4 (30 min)
   - Why Kubernetes, core objects, Deployment and Service manifests

3. **README.md** Sections 5-9 (35 min)
   - ConfigMaps/Secrets, probes, resources, kubectl, rollouts

4. **EXERCISES.md** Exercises 1-5 (2.5 hours)
   - Pod, Deployment, both Service types, ConfigMap

5. **STUDY_GUIDE.md** (30 min)
   - Study the relationship diagrams and comparison tables

6. **EXERCISES.md** Exercises 6-10 (2.5+ hours)
   - Secrets, probes, resources, kubectl walkthrough, comprehensive practice

7. **QUICK_REFERENCE.md** (10 min)
   - Create your mental cheat sheet

---

## 🏆 Success Indicators

You've mastered Level 36 when:

- ✅ You can explain why a Deployment, not a bare Pod, is the normal way to run an app
- ✅ You can write a Deployment + Service pair for a Go app without a reference
- ✅ You never store a sensitive value in a ConfigMap
- ✅ You reach for both readiness and liveness probes by default
- ✅ You always set resource requests and limits
- ✅ You understand `kubectl apply` is declarative, not a one-off action
- ✅ You've completed 8+ exercises
- ✅ You can explain the whole Pod → Deployment → Service chain to someone else

---

## 🚀 What's Next?

After Level 36, you're ready for:

**Level 37: CI/CD**
- Automating build → test → containerize → deploy
- Wiring Level 35's Docker image and this level's manifests into a pipeline
- Continuous integration and continuous delivery concepts

---

## 💬 Key Takeaway

> **A Deployment keeps N copies of your container alive and heals them automatically; a Service gives clients one stable address no matter which copies are currently running - together they turn a single Docker image into a resilient, scalable service.**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Exercises
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Patterns
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Happy learning! Level 36 turns your Level 35 Docker image into a resilient, scalable, production-shaped service! 🎉

*Estimated time to complete Level 36: 4-6 hours*
*Difficulty: ⭐⭐⭐ (Advanced)*
*Next Level: Level 37 - CI/CD*
