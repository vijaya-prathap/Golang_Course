# Level 37: Study Guide & Visual Reference

> ⚠️ **Verification scope:** The Go commands in this guide (`go build`, `go vet`, `gofmt`, `go test -race -cover`, `golangci-lint run`) were genuinely run in this environment - their output is real. Every workflow YAML shown or referenced was validated with a real YAML parser for syntax and structure. No workflow was ever executed on an actual GitHub runner (no git repository exists in this course), and the Docker-push/Kubernetes-deploy portions are illustrative, exactly like Levels 35 and 36. See README.md's introduction for the full breakdown.

## 📚 Learning Path

### Week 1: The Pipeline Itself
```
Day 1:  What CI/CD is - the problem it solves beyond "remember to test"
Day 2:  GitHub Actions anatomy - workflow, job, step, runner
Day 3:  A basic build+vet+test workflow
Day 4:  Quality gates - gofmt (and its exit-code gotcha!)
Day 5:  Quality gates - race detector, coverage, golangci-lint
Day 6:  Matrix builds across Go versions
Day 7:  Caching dependencies
```

### Week 2: Shipping & Practice
```
Day 1:  Docker build-and-push job (bridging Level 35)
Day 2:  Kubernetes deploy job sketch (bridging Level 36)
Day 3:  Deployment strategies - rolling, blue-green, canary
Day 4:  Exercises 1-4 (basic workflow, gofmt, race/cover, lint)
Day 5:  Exercises 5-7 (triggers, matrix, caching)
Day 6:  Exercises 8-10 (Docker, K8s, comprehensive REST API pipeline)
Day 7:  Bonus challenges & review
```

---

## 🔄 CI/CD Pipeline Flow

```
  ┌─────────┐     ┌─────────┐     ┌─────────┐     ┌──────────┐     ┌─────────┐
  │ COMMIT  │ ──▶ │  BUILD  │ ──▶ │  TEST   │ ──▶ │ PACKAGE  │ ──▶ │ DEPLOY  │
  │ (push/  │     │ go build│     │ go test │     │ (Docker  │     │ (apply  │
  │  PR)    │     │ ./...   │     │ + gates │     │  image)  │     │  to K8s)│
  └─────────┘     └─────────┘     └─────────┘     └──────────┘     └─────────┘
       │               │               │                │                │
   triggers      catches syntax   catches bugs,    Level 35's        Level 36's
   the whole     errors, first    races, style     multi-stage       Deployment,
   workflow      and fastest      violations       build, pushed     rolling
                 check                              to a registry     update
```

**The core idea:** every stage is a **gate** - if it fails, everything after it never runs. A broken build never gets to "test." A test failure never gets packaged into an image. An untested image never gets deployed. This is what `needs:` enforces between jobs, and what a failing step enforces within a job (later steps simply don't execute).

---

## 🏗️ GitHub Actions Anatomy

```
.github/workflows/ci.yml                              ← ONE FILE = one workflow
│
├─ name: CI                                            (label shown in GitHub's UI)
│
├─ on:                                                  ← TRIGGER
│   push:        { branches: [main] }                    runs on push to main
│   pull_request: { branches: [main] }                    runs on PRs targeting main
│
└─ jobs:                                                ← one or more, PARALLEL by default
    │
    ├─ build-and-test:                                  ← a JOB
    │    runs-on: ubuntu-latest                            ← the RUNNER (a fresh VM)
    │    steps:                                          ← ordered, sequential within a job
    │      - uses: actions/checkout@v4                     STEP: run a marketplace action
    │      - uses: actions/setup-go@v5                     STEP: run a marketplace action
    │        with: { go-version: '1.23' }
    │      - run: go build ./...                           STEP: run a shell command
    │      - run: go vet ./...
    │      - run: go test ./...
    │
    └─ docker:                                           ← a SECOND job
         needs: build-and-test                             ← waits for the first job
         runs-on: ubuntu-latest
         steps: [ ... ]
```

**Key relationships:**
- Jobs run in **parallel** unless linked with `needs:` - `needs: build-and-test` means "don't even start until that job finishes successfully."
- Steps run **sequentially**, top to bottom, within one job - and stop at the first failing step.
- A runner is a **fresh, disposable machine** every single run - nothing persists between runs except what you explicitly cache (see below).

---

## ✅ Quality Gates Checklist Diagram

```
                    ┌─────────────────────────────────────┐
                    │   Does the code even compile?        │
                    │   go build ./...                     │
                    └───────────────┬───────────────────────┘
                                    │ pass
                                    ▼
                    ┌─────────────────────────────────────┐
                    │   Is it formatted correctly?         │
                    │   test -z "$(gofmt -l .)"             │  🚨 NOT gofmt -l's own
                    │   (check for EMPTY output!)           │     exit code - see below
                    └───────────────┬───────────────────────┘
                                    │ pass
                                    ▼
                    ┌─────────────────────────────────────┐
                    │   Any suspicious-but-compiling code? │
                    │   go vet ./...                       │
                    └───────────────┬───────────────────────┘
                                    │ pass
                                    ▼
                    ┌─────────────────────────────────────┐
                    │   Any style/correctness issues       │
                    │   beyond vet's fixed checklist?      │
                    │   golangci-lint run ./...             │
                    └───────────────┬───────────────────────┘
                                    │ pass
                                    ▼
                    ┌─────────────────────────────────────┐
                    │   Does it work, including under      │
                    │   real concurrency?                  │
                    │   go test -race -cover ./...          │
                    └───────────────┬───────────────────────┘
                                    │ pass
                                    ▼
                              ✅ MERGE / SHIP
```

🚨 **The one gotcha that breaks this whole diagram if you get it wrong:** `gofmt -l .` exits `0` even when it lists files that need formatting - confirmed for real in this environment (`gofmt -l .` printed `main.go` and still reported exit code `0`). The gate must check whether the **output string** is empty, never trust the exit code alone:

```bash
unformatted=$(gofmt -l .)
if [ -n "$unformatted" ]; then
  echo "$unformatted"
  exit 1
fi
```

---

## 🎁 Local Command → CI Step Cheat Sheet

| You'd run locally | The CI step that does the same thing |
|---|---|
| `go build ./...` | `- run: go build ./...` |
| `go vet ./...` | `- run: go vet ./...` |
| `gofmt -l .` | `- run: \| \n    unformatted=$(gofmt -l .) \n    if [ -n "$unformatted" ]; then exit 1; fi` |
| `go test ./...` | `- run: go test ./...` |
| `go test -race -cover ./...` | `- run: go test -race -cover ./...` |
| `golangci-lint run ./...` | `- uses: golangci/golangci-lint-action@v6` |
| installing Go itself | `- uses: actions/setup-go@v5` `with: { go-version: '1.23' }` |
| cloning the repo | `- uses: actions/checkout@v4` (implicit locally - you're already in the repo!) |

---

## 🚦 Deployment Strategy Comparison

| Strategy | How it works | Extra infrastructure | Rollback speed | Blast radius of a bad release |
|----------|---------------|------------------------|------------------|-------------------------------|
| **Rolling** | Replace old instances with new ones a few at a time | None beyond the orchestrator | Moderate - roll old ones back in | Grows gradually toward 100% |
| **Blue-Green** | Deploy fully to an idle second environment, then flip all traffic at once | Double (two full live-capable environments) | Instant - flip traffic back to blue | 100% the instant you flip |
| **Canary** | Send a small % of real traffic to the new version, increase gradually if healthy | Traffic-splitting + monitoring | Fast - route traffic away from canary | Small, by design |

```
Rolling:     [old][old][old] → [new][old][old] → [new][new][old] → [new][new][new]
                                (gradual replacement, both versions live mid-rollout)

Blue-Green:  BLUE (live) ──────────────┐
             GREEN (idle, new version) ┴─── flip router ───▶ GREEN now live, BLUE idle

Canary:      95% traffic → old version  ┐
              5% traffic → new version  ┴─ watch metrics ─▶ shift more traffic if healthy
```

---

## 🚨 Common Mistakes

### Mistake 1: Trusting gofmt -l's Exit Code
```bash
# ❌ WRONG - gofmt -l ALWAYS exits 0 unless there's a parse error
gofmt -l .

# ✅ RIGHT - check the OUTPUT, not the exit code
test -z "$(gofmt -l .)"
```

### Mistake 2: Not Pinning Action Versions
```yaml
# ❌ WRONG - unreproducible; behavior can change with zero changes to your repo
- uses: actions/checkout@main

# ✅ RIGHT
- uses: actions/checkout@v4
```

### Mistake 3: Secrets in the Workflow File
```yaml
# ❌ WRONG - visible in the repo forever, in every past commit
- run: docker login -u admin -p hunter2

# ✅ RIGHT
- run: docker login -u ${{ secrets.DOCKERHUB_USERNAME }} -p ${{ secrets.DOCKERHUB_TOKEN }}
```

### Mistake 4: No Caching
```yaml
# ❌ WRONG - every run re-downloads every module from zero
cache: false

# ✅ RIGHT
cache: true
```

### Mistake 5: Testing Only on main
```yaml
# ❌ WRONG - broken code is already merged by the time anyone finds out
on:
  push:
    branches: [main]

# ✅ RIGHT - catch it before the merge is even possible
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
```

---

## 📈 Progression Summary

### Understanding Level 37

Level 37 automates everything Levels 26 (Testing), 33 (Logging), 34 (Configuration), 35 (Docker), and 36 (Kubernetes) taught you how to do by hand:

1. **CI** - build and test every change automatically (Sections 1-4)
2. **Scale and speed** - matrix builds, caching (Sections 5-6)
3. **CD** - package (Docker) and ship (Kubernetes) automatically, gated on CI passing (Sections 7-8)
4. **Strategy** - how a new version actually replaces an old one safely (Section 9)

### Prerequisites for Level 38

Before moving to Level 38 (Production Debugging), you need:

- ✅ Comfortable reading and writing a GitHub Actions workflow file
- ✅ Understand why `gofmt -l .` needs an output check, not an exit-code check
- ✅ Comfortable with `go test -race -cover` and what a real data race report looks like
- ✅ Understand `needs:` for job ordering and `secrets.*` for credentials
- ✅ Can compare rolling, blue-green, and canary deployment strategies

### Ready for Level 38?

Level 38 picks up right where a shipped, CI/CD-automated pipeline leaves off: something is running in production, right now, and it's misbehaving.
- Reading logs, metrics, and traces under real pressure
- Correlating a CI/CD deploy with a production incident
- Techniques for diagnosing issues you can't reproduce locally

---

## ✅ Checklist Before Level 38

- [ ] Can explain the difference between CI and CD in one sentence each
- [ ] Can write a correct `on:` / `jobs:` / `steps:` / `runs-on:` workflow from memory
- [ ] Know why `gofmt -l .`'s exit code alone can't gate a build
- [ ] Can explain what `go test -race` catches that a normal test run misses
- [ ] Can explain what a Go CI cache actually holds (module cache + build cache)
- [ ] Understand why secrets are referenced, never hardcoded
- [ ] Can compare rolling, blue-green, and canary deployments
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The One Sentence Version
CI/CD turns "please remember to run the checks" into "the checks run themselves, on every change, whether anyone remembers or not."

### The Gate Chain
Build → format → vet → lint → test(race+cover) → package → deploy - each one blocks everything after it.

### The Two Sharpest Gotchas
`gofmt -l .` exits `0` even when it finds problems - check its output, not its exit code. And a generic YAML parser reads `on:` as the boolean `true`, not the string `"on"` - GitHub's own tooling isn't confused by this, but your own scripts might be.

### The Bridge
This level doesn't replace Levels 35-36 - it automates them. The Dockerfile you wrote by hand in Level 35 and the Kubernetes manifests you wrote by hand in Level 36 are exactly what a CI/CD pipeline builds and applies for you, automatically, on every change that passes its gates.

---

## 📚 Next Level

Level 38: Production Debugging
- Diagnosing issues in a system that's already running and already shipped
- Reading logs, metrics, and traces under real production pressure
- Everything from Levels 33 (Logging) through 37 (CI/CD) comes together once something goes wrong at 3 AM

You've automated the pipeline - now let's talk about what happens when it ships something that breaks anyway! 🚀
