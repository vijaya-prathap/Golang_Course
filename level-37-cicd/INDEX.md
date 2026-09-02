# Level 37: CI/CD - INDEX

Welcome to **Level 37: CI/CD**! This is where every check you've been running by hand since Level 26 - build, vet, format, test, lint - starts running itself, automatically, on every single change.

> ⚠️ **Verification scope:** Every Go command in this level's materials was genuinely run in the environment that authored this course; every workflow YAML file was validated with a real YAML parser for syntax and structure. No workflow was ever executed on an actual GitHub runner (there's no git repository in this course), and the Docker-push/Kubernetes-deploy sketches follow the same illustrative disclaimer Levels 35-36 already gave you. Full breakdown in README.md's introduction.

---

## 📖 What You'll Learn

- ✅ What CI (Continuous Integration) and CD (Continuous Delivery/Deployment) each mean, and why they matter
- ✅ GitHub Actions fundamentals: workflows, triggers, jobs, steps, runners
- ✅ A complete, correct build+vet+test workflow for Go
- ✅ Quality gates: gofmt (and its real exit-code gotcha), race detection, coverage, golangci-lint
- ✅ Matrix builds across multiple Go versions
- ✅ Dependency caching - what's cached and why it matters
- ✅ Building and pushing a Docker image in CI (bridging Level 35)
- ✅ Sketching a Kubernetes deploy job (bridging Level 36)
- ✅ Rolling, blue-green, and canary deployment strategies

---

## 🗂️ Level 37 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview, honesty note on verification scope, and welcome message. Start here first!

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- CI/CD concepts and the problem they solve
- GitHub Actions anatomy (workflow/job/step/runner)
- A basic, real, working CI workflow
- Quality gates and the real gofmt exit-code gotcha
- Matrix builds, caching, Docker, Kubernetes bridging
- Deployment strategies
- Best practices
- Common mistakes

**Read Time:** 50-65 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 detailed exercises + 3 bonus challenges - each one builds a real Go project, runs the exact commands a CI workflow would run (with genuinely captured output), and writes/validates the matching workflow YAML:
1. A basic build+vet+test workflow
2. The gofmt quality gate (and its exit-code gotcha)
3. Race detector and coverage
4. Adding a linter (golangci-lint)
5. Triggers - on: push vs on: pull_request
6. Matrix builds
7. Understanding dependency caching
8. Docker build-and-push job (bridging Level 35)
9. A conceptual Kubernetes deploy job (bridging Level 36)
10. Comprehensive practice - a complete CI pipeline for a REST API

**Time Commitment:** 5-6 hours
**When:** After reading README

### 4. **STUDY_GUIDE.md** (Visual Learning)
- CI/CD pipeline flow diagram (commit → build → test → deploy)
- GitHub Actions anatomy diagram
- Quality-gates checklist diagram
- Deployment-strategy comparison table
- Common mistakes

**Read Time:** 30-40 minutes
**When:** Use while doing exercises

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- Workflow anatomy at a glance
- The core gate sequence, ready to copy
- Common mistakes table

**Read Time:** 10-15 minutes
**When:** Quick lookup during practice

---

## 🎯 Recommended Learning Path

### Day 1: CI/CD Concepts & GitHub Actions Anatomy (2 hours)
1. Read **START_HERE.txt** (10 min)
2. Read **README.md** Sections 1-3 (35 min)
3. Complete Exercise 1 (1 hour 15 min)

### Day 2: Quality Gates (2.5 hours)
1. Read **README.md** Section 4 (25 min)
2. Complete Exercises 2-4 (2 hours)

### Day 3: Triggers, Scale, and Speed (2 hours)
1. Read **README.md** Sections 5-6 (25 min)
2. Complete Exercises 5-7 (1.5 hours)

### Day 4: Shipping - Docker & Kubernetes Bridging (2 hours)
1. Read **README.md** Sections 7-9 (30 min)
2. Complete Exercises 8-9 (1.5 hours)

### Day 5: Comprehensive Practice & Consolidation (1.5 hours)
1. Read **README.md** Sections 10-11 (15 min)
2. Complete Exercise 10 (1 hour)
3. Try bonus challenges

---

## 💡 Key Concepts At A Glance

### Workflow Anatomy
```yaml
on: { push: { branches: [main] } }     # TRIGGER
jobs:
  build:
    runs-on: ubuntu-latest              # RUNNER
    steps:                              # sequential STEPS
      - uses: actions/checkout@v4
      - run: go build ./...
```

### The Gate Sequence
```
gofmt (check OUTPUT, not exit code!) → build → vet → lint → test -race -cover
```

### Job Ordering
```yaml
docker:
  needs: build-and-test    # waits for that job to succeed first
```

---

## ✅ Prerequisites

Make sure you've completed **Level 36: Kubernetes**

You need:
- ✅ Comfortable writing a Dockerfile and understanding multi-stage builds (Level 35)
- ✅ Comfortable with Kubernetes Deployment/Service manifests and rolling updates (Level 36)
- ✅ Comfortable with `go test`, table-driven tests, and `-race`/`-cover` (Level 26)
- ✅ Basic familiarity with YAML syntax (used throughout Levels 35-36)

If you haven't done Level 36, go back and complete it first!

---

## 🎓 Learning Objectives

By the end of Level 37, you'll be able to:

- ✅ Explain CI and CD as distinct but related ideas
- ✅ Write a correct GitHub Actions workflow from scratch
- ✅ Explain why `gofmt -l .`'s exit code alone can't gate a build, and write a gate that works
- ✅ Add race detection, coverage, and a linter as CI quality gates
- ✅ Write a matrix build across multiple Go versions
- ✅ Explain exactly what a Go CI cache holds and why
- ✅ Sketch a gated pipeline: test → build & push a Docker image → deploy to Kubernetes
- ✅ Compare rolling, blue-green, and canary deployment strategies

---

## 📊 Statistics

- **Main Theory:** README.md covering CI/CD concepts, GitHub Actions, quality gates, matrix builds, caching, Docker/Kubernetes bridging, and deployment strategies
- **Exercises:** 10 detailed exercises + 3 bonus challenges
- **Study Guide:** Pipeline flow diagram, Actions anatomy diagram, quality-gates checklist, deployment-strategy table
- **Quick Reference:** One-page cheat sheet
- **Practice Time:** 5-6 hours
- **Difficulty:** ⭐⭐⭐ (Advanced)

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for comprehensive theory
2. Use STUDY_GUIDE.md alongside for diagrams and the quality-gates checklist
3. Reference QUICK_REFERENCE.md for fast lookup

### For Practice
1. Read exercise description
2. Copy provided bash commands (or type them)
3. Follow step-by-step instructions - every Go command genuinely runs on your machine
4. Verify output matches expected results
5. Validate every workflow YAML file with the same `python3 -c "import yaml; ..."` command shown throughout

### For Review
1. Use QUICK_REFERENCE.md for quick recap
2. Refer to the common mistakes section
3. Practice writing the gate sequence (gofmt → build → vet → lint → test) from memory

---

## 🆘 Common Questions

**Q: Why does `gofmt -l .` exit 0 even when it lists files?**
A: `gofmt -l` only returns non-zero on a genuine parse error - never merely because a file needs reformatting. Confirmed for real in this environment. A CI gate must check whether its *output* is empty, not trust its exit code.

**Q: Could this level's workflows actually be tested end-to-end?**
A: No - that requires a real GitHub repository and GitHub's own runner infrastructure, neither of which exists in this sandbox. Every Go command a workflow would run was genuinely executed locally, and every workflow file was validated as syntactically correct YAML - but the full, GitHub-hosted execution (triggering, caching across runs, real secrets, real registry/cluster auth) was not exercised.

**Q: Is Docker or Kubernetes required to complete this level?**
A: No. Like Levels 35 and 36 themselves, the Docker-push and Kubernetes-deploy portions of this level are illustrative sketches - correct per their documented behavior, not run against a live daemon or cluster here. Everything else (the Go commands and the YAML files) is fully real and verifiable on your own machine right now.

**Q: What's the difference between CI and CD?**
A: CI (Continuous Integration) is "automatically build and test every change." CD is what happens after CI passes - Continuous *Delivery* usually still has a human click "release"; Continuous *Deployment* ships automatically with no human in the loop.

**Q: Why does `strategy.matrix` matter if I only have one Go version installed?**
A: It matters more on GitHub's infrastructure, where any Go version can be installed on demand by `actions/setup-go`. Locally, you're limited to whatever's actually installed on your machine - which is exactly the honesty note this level gives in Exercise 6 and Section 5 of README.md.

---

## 🎯 Before Moving to Level 38

Make sure you can answer these questions:

- [ ] What's the difference between CI and CD?
- [ ] What are the four core GitHub Actions concepts (workflow, trigger, job, step/runner)?
- [ ] Why can't you trust `gofmt -l .`'s exit code alone?
- [ ] What does `go test -race` catch that a normal test run misses?
- [ ] What does a Go CI cache actually hold?
- [ ] How does `needs:` change job execution order?
- [ ] What's the difference between rolling, blue-green, and canary deployments?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (10 min)
   - Gets you oriented, including the verification-scope honesty note

2. **README.md** Sections 1-4 (40 min)
   - CI/CD concepts, GitHub Actions anatomy, a basic workflow, quality gates

3. **EXERCISES.md** Exercises 1-4 (2.5 hours)
   - Build a real workflow, hit the real gofmt gotcha, see a real race, add a real linter

4. **README.md** Sections 5-9 (35 min)
   - Matrix builds, caching, Docker/Kubernetes bridging, deployment strategies

5. **STUDY_GUIDE.md** (30 min)
   - Study the pipeline diagram and quality-gates checklist

6. **EXERCISES.md** Exercises 5-10 (3+ hours)
   - Triggers, matrix, caching, Docker/K8s sketches, the comprehensive REST API pipeline

7. **QUICK_REFERENCE.md** (10 min)
   - Create your mental cheat sheet

---

## 🏆 Success Indicators

You've mastered Level 37 when:

- ✅ You can write a correct `on:`/`jobs:`/`steps:`/`runs-on:` workflow without looking it up
- ✅ You instinctively check `gofmt -l .`'s *output*, never its exit code, when gating a build
- ✅ You can explain what `-race` and `-cover` each add to a test run, and why the order (fastest checks first) matters
- ✅ You can explain what a Go CI cache holds and why `go.sum`'s hash is the right cache key
- ✅ You can sketch a `needs:`-gated pipeline: test → build image → deploy
- ✅ You can compare rolling, blue-green, and canary deployments without hesitation
- ✅ You've completed 8+ exercises
- ✅ You can explain CI/CD to someone else

---

## 🚀 What's Next?

After Level 37, you're ready for:

**Level 38: Production Debugging**
- Diagnosing issues in a system that's already running and already shipped
- Reading logs, metrics, and traces under real production pressure
- Everything from Levels 33 (Logging) through 37 (CI/CD) comes together once something goes wrong at 3 AM

---

## 💬 Key Takeaway

> **CI/CD turns "please remember to run the checks" into "the checks run themselves, on every single change, whether anyone remembers or not" - and `gofmt -l .`'s silent exit-code gotcha is exactly the kind of mistake that automation exists to catch for good.**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Exercises
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Patterns
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Happy learning! Level 37 turns everything you've built so far into a pipeline that proves itself, automatically, forever! 🎉

*Estimated time to complete Level 37: 5-6 hours*
*Difficulty: ⭐⭐⭐ (Advanced)*
*Next Level: Level 38 - Production Debugging*
