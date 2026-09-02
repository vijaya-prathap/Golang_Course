# Level 37: CI/CD - Complete Guide

## Introduction

Welcome to Level 37! You've mastered `Level 36: Kubernetes` - describing the desired state of your running application in YAML and letting a cluster keep reality matching it. Now it's time to close the loop: **CI/CD** is what automatically builds, tests, and (optionally) ships every change you make, so that "it works on my machine" gets checked by a machine other than yours before anyone else ever sees the change.

> ⚠️ **Environment note, read this first.** This level is not pure Go code, and like Levels 35 (Docker) and 36 (Kubernetes) it needed an honest adaptation of this course's normal "run it for real" rule. Here is exactly what that means:
>
> - **Every Go command a workflow would run - genuinely executed.** `go build ./...`, `go vet ./...`, `gofmt -l .`, `go test ./... -race -cover`, and `golangci-lint run ./...` were all actually run, in real scratch Go projects, in the environment that authored this course. `golangci-lint` (v1.64.8) happened to be installed here, so its output in this level is real captured output too, not a guess. Go `1.26.5` was the only Go toolchain available in this environment - see the honesty note in [Matrix Builds](#matrix-builds) for exactly what that does and doesn't limit.
> - **Every workflow YAML file in this level was validated for real** with a YAML parser (`python3 -c "import yaml; yaml.safe_load(...)"`), confirming it is syntactically valid YAML and structurally sane (has `on`, `jobs`, and each job has `runs-on` and `steps`) per GitHub Actions' documented, stable schema.
> - **What was NOT exercised:** an actual GitHub-hosted run of any workflow. That requires pushing to a real GitHub repository and letting GitHub's own infrastructure schedule a runner, resolve `uses:` actions from the marketplace, cache across runs, and authenticate against a real container registry or cluster - none of which exists in this sandbox (confirmed: there is no git repository at all in this course's directory). `docker` and `kubectl` client binaries are present here (checked: `docker version`, `kubectl version --client` both returned real client info), but `docker version`'s server section and any real cluster connection are not - the same situation Level 35 and Level 36 were honest about, so the Docker-push and Kubernetes-deploy jobs in this level are **illustrative, per that same disclaimer**, not independently re-verified.
>
> Bottom line: trust the Go command output and the YAML-validity claims in this level completely - they're real. Treat anything described as "illustrative" as correct per GitHub Actions', Docker's, and Kubernetes' documented behavior, but verify it yourself the first time you actually push to a real repository.

Why does this level exist at all, separate from just "writing good Go code"? Because code that passes on your laptop can still break the build for everyone else: you forgot to run `go vet`, a teammate's change didn't get tested before merging, someone's `.go` file was never `gofmt`-formatted and now every diff in the repo is noisy. **CI/CD turns "please remember to run the checks" into "the checks run themselves, on every single change, whether anyone remembers or not."**

---

## Table of Contents

1. [What CI/CD Is and Why It Matters](#what-cicd-is-and-why-it-matters)
2. [GitHub Actions Fundamentals](#github-actions-fundamentals)
3. [A Basic Go CI Workflow](#a-basic-go-ci-workflow)
4. [Adding Quality Gates](#adding-quality-gates)
5. [Matrix Builds](#matrix-builds)
6. [Caching Dependencies](#caching-dependencies)
7. [Building and Pushing a Docker Image in CI](#building-and-pushing-a-docker-image-in-ci)
8. [Triggering Deployment](#triggering-deployment)
9. [Deployment Strategies Overview](#deployment-strategies-overview)
10. [Best Practices](#best-practices)
11. [Common Mistakes](#common-mistakes)

---

## What CI/CD Is and Why It Matters

**CI/CD** is actually two related ideas, and it's worth keeping them separate in your head:

- **Continuous Integration (CI)** - every time someone pushes a change (or opens a pull request), a machine automatically builds the project and runs its tests. If anything is broken, everyone finds out within minutes, on that one change - not weeks later when ten broken changes are tangled together and nobody can tell which one caused the failure.
- **Continuous Delivery / Deployment (CD)** - once a change has passed CI, it's automatically packaged (built into a container image, for example) and, depending on how far you take it, automatically shipped to production. "Delivery" usually implies a human still clicks a button to release; "Deployment" means it ships with no human in the loop at all, as long as every automated check passed.

### The Problem This Solves

Picture a team without CI. Everyone works on their own branch for a few days. When it's finally time to merge, five different people's changes collide, nobody's sure whose code broke the build, and someone spends an afternoon doing archaeology through `git log` to find out. Multiply this by every release, forever.

Now picture the same team with CI. The instant anyone pushes, a workflow builds and tests *that exact change* in isolation. A broken change gets flagged in minutes, while the mistake is still fresh in the author's mind and easy to fix - not after it's buried under three days of everyone else's unrelated work.

### The Core Insight

> **Catch problems at the moment they're introduced, not at the moment you try to ship.**

Every section that follows is really just one idea applied at increasing levels of ambition: start by automatically compiling and testing every change (Sections 3-4), scale that check across environments (Section 5), make it fast enough that nobody's tempted to skip it (Section 6), and then extend the same "automate it so it can't be forgotten" philosophy all the way through packaging (Section 7) and shipping (Sections 8-9).

---

## GitHub Actions Fundamentals

**GitHub Actions** is GitHub's built-in CI/CD system. You configure it entirely with YAML files - no separate service to sign up for, no separate UI to configure by hand. Every workflow file lives in one specific place in your repository:

```
your-repo/
└── .github/
    └── workflows/
        ├── ci.yml
        └── deploy.yml
```

Any `.yml` (or `.yaml`) file directly inside `.github/workflows/` is a **workflow** - GitHub watches that folder and automatically picks up anything you put there. You can have as many workflow files as you want; each one is independent.

### The Four Concepts You Need

```
Workflow                     the whole YAML file - "what should happen and when"
  │
  ├─ on:                     the TRIGGER - what event starts this workflow
  │
  └─ jobs:                   one or more independent units of work
       │
       └─ <job-name>:
            ├─ runs-on:      which machine (RUNNER) executes this job
            └─ steps:        an ordered list of actions to run, one after another
```

**Triggers (`on:`)** - the event that starts a workflow run. The two you'll use constantly:

```yaml
on:
  push:
    branches: [main]        # runs whenever someone pushes to main
  pull_request:
    branches: [main]        # runs whenever a PR targeting main is opened/updated
```

**Jobs** - a workflow can have one job or many. By default, every job in a workflow runs **in parallel**, on its own fresh runner, unless you tell one job to wait for another with `needs:` (you'll see this in [Section 7](#building-and-pushing-a-docker-image-in-ci)).

**Steps** - the ordered list of things a job does. A step is either:
- `run:` - execute a shell command directly (`run: go build ./...`)
- `uses:` - run a pre-built, reusable **action** from GitHub's marketplace (`uses: actions/checkout@v4`)

**Runners (`runs-on:`)** - the actual machine that executes a job. `ubuntu-latest` is by far the most common choice for backend/Go projects: it's fast to provision, free for public repositories, and Linux is what almost every Go service ultimately deploys to anyway.

```yaml
runs-on: ubuntu-latest
```

Other options exist (`windows-latest`, `macos-latest`, self-hosted runners you provide yourself) but `ubuntu-latest` is the default assumption for every example in this level unless stated otherwise.

### Minimal Skeleton

Putting the four concepts together, here is the smallest workflow that actually does something:

```yaml
name: Hello CI

on:
  push:
    branches: [main]

jobs:
  say-hello:
    runs-on: ubuntu-latest
    steps:
      - name: Print a greeting
        run: echo "Hello from GitHub Actions!"
```

- `name:` at the top level is just a human-readable label shown in GitHub's UI - optional, but always worth setting.
- Each step's own `name:` is likewise a label; the field that actually does the work is `run:` or `uses:`.

This workflow is syntactically real and was validated the same way as every other file in this level - but it has never executed on a GitHub runner in this environment, for the reasons explained in the introduction.

---

## A Basic Go CI Workflow

Here is a complete, correct workflow for a real Go project: check out the code, install Go, then run the three commands that catch the overwhelming majority of "I forgot to..." mistakes.

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  build-and-test:
    runs-on: ubuntu-latest
    steps:
      - name: Check out code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Build
        run: go build ./...

      - name: Vet
        run: go vet ./...

      - name: Test
        run: go test ./...
```

### Step by Step

**`actions/checkout@v4`** - almost every workflow's first step. Without it, the runner is an empty machine with no copy of your repository at all; this action clones your repo into the runner's working directory. The `@v4` pins an exact major version of the action itself (more on why this matters in [Best Practices](#best-practices)).

**`actions/setup-go@v5`** - downloads and installs the requested Go toolchain version onto the runner, and puts `go` on `PATH` for every step that follows. `with: go-version: '1.23'` is quoted deliberately - YAML would otherwise parse `1.23` as a floating-point number and could silently normalize it (e.g., dropping a trailing zero), so version strings are always quoted in this level's examples.

**`go build ./...`** - compiles every package in the module (the `./...` pattern means "this directory and everything under it"), catching syntax errors and type errors before anything else runs. This is the fastest possible check and the first one that should ever fail.

**`go vet ./...`** - Go's built-in static analyzer, catching mistakes that compile fine but are almost certainly bugs: a `Printf` format string that doesn't match its arguments, a struct copied by value when it contains a `sync.Mutex`, an unreachable code path. `go vet` ships with every Go installation - no extra tool to install.

**`go test ./...`** - runs every `_test.go` file's tests across the whole module, exactly like you've done locally since `Level 26: Testing`. If any test fails, this step - and the whole job - fails.

### This Exact Sequence, Run for Real

Every one of these three commands was actually run against a real scratch Go module while writing this level, to confirm the output shown throughout this course is genuine and not just "what should happen in theory." Here is that real run, against a tiny `mathutil` package:

```
$ go build ./...
$ echo "exit code: $?"
exit code: 0

$ go vet ./...
$ echo "exit code: $?"
exit code: 0

$ go test ./...
?   	level37.example/exercise1	[no test files]
ok  	level37.example/exercise1/mathutil	0.525s
```

Notice the `?` line - that's `go test` telling you package `level37.example/exercise1` (the one holding just `main.go`, no `_test.go` file) was skipped, not failed. That's expected and fine; only the `mathutil` subpackage has tests, and it reports `ok`.

### Order Matters: Fail Fast

Notice the order: build, then vet, then test. This is deliberate, not arbitrary - see [Best Practices](#best-practices) for why the fastest, cheapest checks should always run first.

---

## Adding Quality Gates

The workflow above catches broken code and failing tests. A **quality gate** goes further: it blocks a merge for problems that aren't outright bugs but that the team has decided matter - inconsistent formatting, races that only show up under concurrency, style violations a human reviewer shouldn't have to point out by hand.

### Gate 1: gofmt - Fail the Build on Unformatted Code

`gofmt -l .` lists every file that isn't canonically formatted (empty output means "everything's clean"). Here's the one genuine gotcha worth knowing before you wire this into CI:

> 🚨 **`gofmt -l` exits `0` even when it lists files.** It only returns a non-zero exit code on a genuine parse error - never merely for "this file needs reformatting." A naive `gofmt -l .` step in CI will print the offending filenames and then report success anyway, silently defeating the entire point of the gate.

This was confirmed for real in this environment, not assumed:

```
$ gofmt -l .
main.go
$ echo "gofmt -l exit code: $?"
gofmt -l exit code: 0
```

`main.go` was listed as needing formatting, and the shell still reports exit code `0`. The fix is to check whether the output was **empty**, not to trust the exit code:

```yaml
- name: Check formatting (gofmt)
  run: |
    unformatted=$(gofmt -l .)
    if [ -n "$unformatted" ]; then
      echo "The following files are not gofmt-formatted:"
      echo "$unformatted"
      exit 1
    fi
```

This step explicitly `exit 1`s only when `gofmt -l .`'s output is non-empty, which is what actually fails the CI job. Here's the same file after fixing it with `gofmt -w .` (which *does* rewrite files in place, real output):

```
$ gofmt -d .
diff main.go.orig main.go
--- main.go.orig
+++ main.go
@@ -3,7 +3,7 @@
 import "fmt"
 
 func main() {
-        x :=    1
+	x := 1
 	y := 2
-    fmt.Println(  "sum:", x+y )
+	fmt.Println("sum:", x+y)
 }

$ gofmt -w .
$ gofmt -l .
$ echo "gofmt -l exit code: $?"
gofmt -l exit code: 0
```

After `gofmt -w .`, `gofmt -l .` prints nothing at all - exactly what the CI step above needs to see to let the build proceed.

### Gate 2: go test -race -cover

`-race` enables Go's race detector, instrumenting the binary to catch unsynchronized concurrent access to shared memory - the class of bug from Levels 20-23 that can pass every normal test run and still corrupt data in production under real concurrent load. `-cover` reports what percentage of statements your tests actually exercised (a floor, not a target - see `Level 26: Testing`'s coverage section for the full nuance).

```yaml
- name: Test with race detector and coverage
  run: go test -race -cover ./...
```

Real output, on a `Counter` type that correctly guards its state with a `sync.Mutex`:

```
$ go test -race -cover ./...
ok  	level37.example/exercise3/counter	1.617s	coverage: 100.0% of statements
```

To make the value of `-race` concrete, here is the *exact same test* run against a deliberately unsynchronized version of the same `Counter` (plain `c.value++`, no mutex) - this is real, captured race-detector output, not a mockup:

```
$ go test -race ./...
==================
WARNING: DATA RACE
Read at 0x00c0000122a8 by goroutine 11:
  level37.example/racedemo/counter.(*Counter).Increment()
      counter.go:10 +0x70
  ...

Previous write at 0x00c0000122a8 by goroutine 9:
  level37.example/racedemo/counter.(*Counter).Increment()
      counter.go:10 +0x80
  ...
==================
--- FAIL: TestCounterConcurrentIncrement (0.00s)
    counter_test.go:21: Value() = 88; want 100
    testing.go:1712: race detected during execution of test
FAIL
FAIL	level37.example/racedemo/counter	0.291s
FAIL
```

(Exact memory addresses, goroutine numbers, and the final count of `88` are machine- and scheduling-dependent - run it yourself and you'll likely see different numbers, but the shape - a `WARNING: DATA RACE`, the wrong final value, and a `FAIL` - is exactly what an unsynchronized counter under concurrent load produces.) Without `-race`, this exact bug can pass silently for months, because `c.value++` usually *happens* to produce the right answer even when it's technically undefined behavior - `-race` is what turns "usually right" into "provably wrong, immediately."

### Gate 3: A Linter - golangci-lint

`go vet` catches a fixed, deliberately conservative set of correctness bugs. **`golangci-lint`** is the standard third-party choice for everything beyond that: unchecked error return values, ineffectual assignments, unused code, and dozens of other checks, run through one fast, parallelized command that aggregates many individual linters at once.

```yaml
- name: Lint
  uses: golangci/golangci-lint-action@v6
  with:
    version: latest
```

`golangci-lint` (v1.64.8) happened to be installed in the environment that authored this course, so unlike Docker and Kubernetes, its output below is genuinely real - not illustrative. Here it is catching a real, deliberately unchecked error:

```
$ golangci-lint run ./...
main.go:10:18: Error return value of `f.WriteString` is not checked (errcheck)
    f.WriteString("Hello, " + name)
                 ^
```

And after fixing it (checking the error and returning it up the call stack, the idiomatic Go pattern from `Level 16: Error Handling`):

```
$ go build ./...
$ go vet ./...
$ golangci-lint run ./...
$ echo "lint exit code: $?"
lint exit code: 0
```

No output at all is `golangci-lint`'s way of saying "found nothing to complain about" - the same convention `gofmt -l` uses, but (unlike `gofmt -l`) `golangci-lint run` *does* exit non-zero when it finds something, which is exactly why it can be used directly as a gate without the `-n` trick `gofmt` needed above.

### Putting the Gates in Order

```yaml
steps:
  - uses: actions/checkout@v4
  - uses: actions/setup-go@v5
    with:
      go-version: '1.23'

  - name: Check formatting (gofmt)
    run: |
      unformatted=$(gofmt -l .)
      if [ -n "$unformatted" ]; then
        echo "$unformatted"
        exit 1
      fi

  - run: go build ./...
  - run: go vet ./...

  - name: Lint
    uses: golangci/golangci-lint-action@v6
    with:
      version: latest

  - name: Test with race detector and coverage
    run: go test -race -cover ./...
```

Formatting and vet run first because they're nearly instant and catch the most common mistakes; the race-enabled test run goes last because `-race` makes tests measurably slower (instrumented binaries do more work per memory access), so it should only run once everything cheaper has already passed.

---

## Matrix Builds

A **matrix** runs the *same* job multiple times with different input values substituted in - most commonly, multiple Go versions, to prove your code works on more than just whatever version happens to be on your own laptop.

```yaml
jobs:
  build-and-test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: ['1.22', '1.23']
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go-version }}
      - run: go build ./...
      - run: go vet ./...
      - run: go test ./...
```

`strategy.matrix.go-version` defines a list; GitHub Actions runs the entire job once **per entry** - here, that's two independent job runs, one pinned to Go 1.22 and one to Go 1.23, executing in parallel by default. `${{ matrix.go-version }}` interpolates whichever value that particular run was given. Add a second matrix dimension (say, `os: [ubuntu-latest, macos-latest]`) and Actions runs the full cross-product - 2 Go versions × 2 operating systems = 4 job runs - with zero extra job definitions to write by hand.

> ⚠️ **Honesty note on this section specifically.** This environment has exactly **one** Go toolchain installed: `go1.26.5`. There was no second Go version available to actually install and test against locally, so the "two different Go versions genuinely produced the same passing result" claim that a real matrix build demonstrates could **not** be verified end-to-end here. What *was* verified: the single available Go version (`1.26.5`) builds, vets, and tests this level's exercise code cleanly (see the real output throughout Sections 3-4), and the matrix YAML above is syntactically valid and uses the documented, stable `strategy.matrix` schema. The concept - "run the identical job once per listed version, in parallel" - is exactly what GitHub Actions does per its own documentation; only the specific claim "and both versions actually passed on a real runner" is untested here. If you have `gvm`, `asdf`, or multiple Go installations available, running the same build/vet/test sequence against each locally is the closest you can get to this without pushing to GitHub.

### Why This Matters in Practice

Go's [release policy](https://go.dev/doc/devel/release) officially supports the two most recent major versions at any time. A library used by many other projects typically tests against both the current and previous Go version, so a user who hasn't upgraded yet still gets a signal their environment is supported. An internal application usually only needs to test against whatever version production actually runs - but even then, testing against "current" and "one version behind" catches upgrade-readiness problems before they become an emergency.

---

## Caching Dependencies

Every job starts on a **fresh runner** with nothing pre-installed beyond the base image - no downloaded Go modules, no compiled build artifacts from last time. Without caching, every single run re-downloads every dependency and recompiles every package from scratch, even when nothing changed since the last run.

### What Actually Gets Cached

Two distinct directories matter for a Go project:

| Directory | What it holds | Real path (confirmed with `go env`) |
|-----------|---------------|--------------------------------------|
| Module download cache | Downloaded `.zip`/`.info`/`.mod` files for every module version your `go.mod` has ever referenced | `go env GOMODCACHE` → `/Users/<you>/go/pkg/mod` |
| Build cache | Compiled package object files, keyed by source hash - lets `go build`/`go test` skip recompiling packages that haven't changed | `go env GOCACHE` → `~/Library/Caches/go-build` (macOS) / `~/.cache/go-build` (Linux, including GitHub's `ubuntu-latest` runners) |

Both paths above are genuinely what this author's machine reports via `go env GOMODCACHE` and `go env GOCACHE` - the point isn't the exact path (it differs by OS), it's that these two caches are what caching in CI is actually caching.

### The Easy Way: setup-go's Built-In Caching

`actions/setup-go@v5` can cache both directories for you with a single option - no separate action needed:

```yaml
- name: Set up Go
  uses: actions/setup-go@v5
  with:
    go-version: '1.23'
    cache: true    # this is the default when a go.sum file exists - shown explicitly here
```

Internally, this uses `go.sum`'s hash as the cache key: as long as your dependencies haven't changed, later runs restore the module and build caches instead of starting from zero, which is normally the single biggest speedup available to a Go CI pipeline. When `go.sum` *does* change (a dependency was added, removed, or upgraded), the key changes too and a fresh cache is built - it can never silently serve a stale cache for the wrong dependency set.

### The Explicit Way: actions/cache

Understanding what `cache: true` does under the hood is easier by seeing the manual equivalent, using the general-purpose `actions/cache` action directly:

```yaml
- name: Set up Go
  uses: actions/setup-go@v5
  with:
    go-version: '1.23'
    cache: false   # disabled here since we're caching manually below

- name: Cache Go modules and build cache
  uses: actions/cache@v4
  with:
    path: |
      ~/.cache/go-build
      ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
    restore-keys: |
      ${{ runner.os }}-go-
```

- **`path:`** - exactly the two directories from the table above.
- **`key:`** - the cache is saved and restored under this exact string. Including `hashFiles('**/go.sum')` means a new key (and therefore a fresh cache) is generated automatically whenever dependencies change.
- **`restore-keys:`** - a fallback prefix. If no cache exactly matches `key` (say, `go.sum` changed since the last run), Actions falls back to the most recent cache whose key starts with this prefix, which is still a highly relevant partial cache rather than nothing at all.

Prefer `actions/setup-go`'s built-in `cache: true` in real workflows - it's simpler and does the same thing. The manual version above exists in this level purely so you understand what's happening rather than treating it as a black box.

---

## Building and Pushing a Docker Image in CI

`Level 35: Docker` covered writing a multi-stage `Dockerfile` and building an image by hand with `docker build`. In CI, that same build should happen automatically after tests pass, and the resulting image should be pushed to a registry other machines (or a Kubernetes cluster) can pull from.

```yaml
name: Build and Push Docker Image

on:
  push:
    branches: [main]

jobs:
  build-and-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - run: go build ./...
      - run: go vet ./...
      - run: go test ./...

  docker:
    needs: build-and-test
    runs-on: ubuntu-latest
    steps:
      - name: Check out code
        uses: actions/checkout@v4

      - name: Log in to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      - name: Build and push image
        uses: docker/build-push-action@v6
        with:
          context: .
          push: true
          tags: |
            mycompany/books-api:${{ github.sha }}
            mycompany/books-api:latest
```

### needs: - Making a Job Wait

`needs: build-and-test` on the `docker` job means: don't start this job until `build-and-test` has finished successfully. Without `needs:`, every job in a workflow runs in parallel and independently - which would mean a broken build could still get packaged and pushed as an image. `needs:` is what actually enforces "only ship things that passed CI."

### Secrets - Never Hardcode Credentials

`${{ secrets.DOCKERHUB_USERNAME }}` and `${{ secrets.DOCKERHUB_TOKEN }}` reference **GitHub Actions secrets**: encrypted values configured once in the repository's (or organization's) settings, never visible in the workflow file, in logs, or in the repository's history. GitHub automatically redacts a secret's value from workflow logs even if a step accidentally prints it.

```yaml
# ❌ NEVER do this - a real password committed straight into version control,
# visible to anyone with read access to the repo, forever, in every past commit
- run: docker login -u myuser -p SuperSecret123!

# ✅ Reference a secret by name - the actual value lives only in
# repository settings, injected at runtime, never written to the file
- run: docker login -u ${{ secrets.DOCKERHUB_USERNAME }} -p ${{ secrets.DOCKERHUB_TOKEN }}
```

Conceptually, a secret is just a named, encrypted environment variable scoped to the repository (or organization): you set it once in **Settings → Secrets and variables → Actions** on GitHub, and any workflow in that repository can reference it by name via the `secrets` context - the workflow file itself never contains, and never needs to contain, the actual credential.

### Registry Choice

`docker/login-action` and `docker/build-push-action` work identically against Docker Hub, GitHub Container Registry (`ghcr.io`), AWS ECR, or any other registry - only the `registry:` input and the secret's contents change. GitHub Container Registry is a common default for projects already hosted on GitHub, since it can authenticate using GitHub's own built-in `GITHUB_TOKEN` instead of a separate registry account.

---

## Triggering Deployment

`Level 36: Kubernetes` covered writing `Deployment` and `Service` manifests and applying them by hand with `kubectl apply`. A deploy job automates the last mile: once a new image has been built and pushed (Section 7), tell the cluster to start using it.

```yaml
  deploy:
    needs: docker
    runs-on: ubuntu-latest
    steps:
      - name: Check out manifests
        uses: actions/checkout@v4

      - name: Configure kubectl
        uses: azure/k8s-set-context@v4
        with:
          kubeconfig: ${{ secrets.KUBE_CONFIG }}

      - name: Set the new image on the Deployment
        run: |
          kubectl set image deployment/books-api \
            books-api=mycompany/books-api:${{ github.sha }} \
            --record

      - name: Wait for the rollout to finish
        run: kubectl rollout status deployment/books-api --timeout=120s
```

Reading this job the same way you'd read the manifests from Level 36:

- **`needs: docker`** - don't attempt to deploy an image that was never successfully built and pushed.
- **`kubeconfig: ${{ secrets.KUBE_CONFIG }}`** - the credentials `kubectl` needs to reach your specific cluster, stored as a secret exactly like the registry credentials in Section 7. A `kubeconfig` file is itself effectively a credential (it can contain a client certificate or token) - it must never be committed to the repository.
- **`kubectl set image`** - the direct, scriptable equivalent of editing a Deployment's `image:` field and re-applying it (Level 36, Section 9: Rolling Updates). Pointing it at `${{ github.sha }}` - the exact commit that triggered this run - means the running Deployment always traces back to one specific, auditable commit.
- **`kubectl rollout status`** - blocks the job until the rolling update actually finishes (or fails), so a broken deploy shows up as a failed CI job instead of silently leaving old and new Pods mixed together.

This job is **kept deliberately conceptual**, exactly as flagged in this level's introduction: there is no real cluster in this environment to apply it against, `azure/k8s-set-context` is one of several valid ways to authenticate `kubectl` in CI (not the only one), and a real setup would also need the Deployment's manifest to already exist in the cluster (from Level 36) before this job could update it. Treat this as a correct sketch of the *shape* of a deploy job, to be adapted to your actual cluster and CI provider's exact authentication mechanism.

---

## Deployment Strategies Overview

Once CI/CD can ship a new version automatically, *how* it replaces the old version running in production matters as much as the automation itself. This is general DevOps knowledge, not specific to GitHub Actions - the same three strategies apply whether you're deploying with Kubernetes, a load balancer and VMs, or a serverless platform.

### Rolling Deployment

Replace old instances with new ones gradually, a few at a time, keeping the service available throughout. This is exactly what a Kubernetes `Deployment`'s default update strategy does (Level 36, Section 9): a handful of new Pods start, old ones are terminated once the new ones are healthy, repeat until every Pod is running the new version.

- **Pro:** No extra infrastructure needed; built into Kubernetes by default.
- **Con:** Old and new versions run simultaneously, mid-rollout - your API must tolerate both versions existing at once (a real constraint on database schema changes, in particular).

### Blue-Green Deployment

Run two complete, identical environments - "blue" (currently live) and "green" (the new version) - side by side. Deploy the new version entirely to green, verify it's healthy, then flip a router or load balancer to send all traffic to green at once. Blue stays fully intact and idle, ready for an instant rollback by flipping traffic back.

- **Pro:** Rollback is immediate - just flip traffic back to blue, no waiting for a reverse rollout.
- **Con:** Needs double the infrastructure (two full environments) running at the switchover moment.

### Canary Deployment

Send a small percentage of real traffic (say, 5%) to the new version while the rest keeps going to the old one. Watch error rates and latency on that small slice; if it looks healthy, gradually increase the percentage until the new version handles 100% of traffic - if it doesn't, roll back after only a small fraction of users were ever affected.

- **Pro:** Limits the *blast radius* of a bad release to a small slice of real traffic, with real production signals guiding the rollout speed.
- **Con:** Needs traffic-splitting infrastructure (a service mesh or a smart load balancer) and more sophisticated monitoring to actually interpret the canary's health.

### Comparing All Three

| Strategy | Extra infrastructure needed | Rollback speed | Blast radius of a bad release |
|----------|------------------------------|-----------------|-------------------------------|
| Rolling | None beyond the orchestrator itself | Moderate (roll old Pods back in) | Up to 100% eventually, but gradual |
| Blue-Green | Double (two full environments) | Instant (flip traffic back) | 100% the instant you flip |
| Canary | Traffic-splitting + monitoring | Fast (route traffic away from canary) | Small, by design |

None of these is universally "best" - a small internal tool may not justify blue-green's doubled infrastructure cost, while a payments system handling millions of requests a day may consider canary's operational complexity well worth the reduced blast radius of a bad release.

---

## Best Practices

### 1. Fail Fast - Cheapest Checks First

```yaml
# ✅ Good - format check and build fail in seconds; race-enabled tests
# (the slowest check) only run once everything cheaper already passed
steps:
  - run: gofmt -l .        # fastest
  - run: go build ./...
  - run: go vet ./...
  - run: go test -race ./...   # slowest
```

Ordering checks from fastest to slowest means a broken build fails in seconds instead of minutes, and it means expensive runner time isn't spent on a race-detector-instrumented test run for code that doesn't even compile.

### 2. Pin Action Versions

```yaml
# ✅ Good - reproducible: this exact major version, forever
- uses: actions/checkout@v4

# ❌ Avoid - "whatever the latest tag currently points to," which can
# change (and silently break your build) with no change to your own repo
- uses: actions/checkout@main
```

Pinning to a tagged version (`@v4`) means your workflow's behavior only changes when you deliberately bump the version. Pointing at a branch (`@main`) means the action's maintainers can change behavior underneath you at any time.

### 3. Never Hardcode Secrets

Covered in depth in [Section 7](#building-and-pushing-a-docker-image-in-ci) - always `${{ secrets.NAME }}`, never a literal credential in the YAML file, ever, even "temporarily."

### 4. Keep Workflows Fast With Caching

A slow CI pipeline gets skipped, worked around, or ignored - see [Caching Dependencies](#caching-dependencies). `cache: true` on `actions/setup-go` costs one line and typically saves the majority of a run's wall-clock time on the module-download step alone.

### 5. Run Checks on Every Pull Request, Not Just main

```yaml
# ✅ Good - every PR gets checked BEFORE it merges
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
```

Testing only on `push` to `main` means broken code is already merged by the time anyone finds out. Testing on `pull_request` as well catches it *before* the merge button is even clickable in most repository configurations (via required status checks).

---

## Common Mistakes

### Mistake 1: Not Pinning Action Versions

```yaml
# ❌ WRONG - unreproducible: this can build differently tomorrow with
# zero changes to your own repository
- uses: actions/setup-go@main

# ✅ RIGHT - pinned to an exact, stable major version
- uses: actions/setup-go@v5
```

### Mistake 2: Forgetting to Cache Dependencies

```yaml
# ❌ WRONG - every single run re-downloads every module from scratch
- uses: actions/setup-go@v5
  with:
    go-version: '1.23'
    cache: false

# ✅ RIGHT - restores the module and build cache when go.sum is unchanged
- uses: actions/setup-go@v5
  with:
    go-version: '1.23'
    cache: true
```

### Mistake 3: Secrets Directly in Workflow YAML

```yaml
# ❌ WRONG - visible to anyone who can read the repo, forever, in git history
- run: docker login -u admin -p hunter2

# ✅ RIGHT - reference a secret configured in repository settings
- run: docker login -u ${{ secrets.DOCKERHUB_USERNAME }} -p ${{ secrets.DOCKERHUB_TOKEN }}
```

### Mistake 4: Not Running go vet / gofmt Checks

```yaml
# ❌ WRONG - only tests run; a Printf format mismatch or unformatted
# file merges cleanly because nothing ever checked for it
- run: go test ./...

# ✅ RIGHT - vet and a real gofmt gate (checking for EMPTY output,
# not trusting gofmt -l's exit code - see Section 4) run too
- run: gofmt -l . | (! grep -q '.') # fails if gofmt -l printed anything
- run: go vet ./...
- run: go test ./...
```

### Mistake 5: Testing Only on Merge Instead of Every PR

```yaml
# ❌ WRONG - broken code is only discovered AFTER it's already on main
on:
  push:
    branches: [main]

# ✅ RIGHT - every PR is checked before it can be merged at all
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
```

### A Real Quirk Found While Writing This Level: `on:` and Generic YAML Parsers

Validating every workflow in this level with a plain YAML parser (`python3 -c "import yaml; yaml.safe_load(...)"`, as described in this level's introduction) surfaced a genuine, documented YAML 1.1 quirk worth knowing:

```python
>>> import yaml
>>> yaml.safe_load(open('ci.yml')).keys()
dict_keys(['name', True, 'jobs'])
```

A generic YAML 1.1 parser (PyYAML included) treats the bare word `on` as the **boolean `true`**, not the string `"on"` - the same rule that makes bare `yes`/`no`/`off` special in YAML 1.1. This is real, confirmed output from validating this level's own files, not a hypothetical. It does **not** affect GitHub Actions itself - GitHub's own workflow parser always treats a top-level `on:` key as the trigger definition, regardless of this generic-YAML ambiguity - but it's exactly the kind of surprise you'd hit if you ever tried to programmatically read a workflow file's trigger with a naive YAML library instead of a purpose-built Actions tool.

---

## Summary

**The Two Halves:**
- **CI** - every push/PR automatically built and tested, catching problems immediately
- **CD** - validated changes automatically packaged and (optionally) shipped

**GitHub Actions Anatomy:**
- `.github/workflows/*.yml` - workflow files, auto-discovered by GitHub
- `on:` (trigger) → `jobs:` (parallel by default, or `needs:`-ordered) → `steps:` (`run:` or `uses:`)
- `runs-on: ubuntu-latest` - the runner

**A Real Go CI Workflow:**
- `actions/checkout@v4` → `actions/setup-go@v5` → `go build ./...` → `go vet ./...` → `go test ./...`

**Quality Gates:**
- `gofmt -l .` - check for **empty output**, not exit code (it's always `0`)
- `go test -race -cover ./...` - catch concurrency bugs and measure statement coverage
- `golangci-lint run ./...` - the standard linter beyond what `go vet` covers

**Scaling and Speed:**
- `strategy.matrix` - the same job, run once per Go version (or OS, or both)
- `actions/setup-go`'s `cache: true` (or `actions/cache` manually) - skip re-downloading/recompiling unchanged dependencies

**Shipping (bridging Levels 35-36):**
- A `docker` job (gated by `needs:`) builds and pushes an image using `secrets.*` for registry credentials, never hardcoded
- A `deploy` job (gated by `needs:` on the docker job) updates a Kubernetes Deployment's image and waits for the rollout - conceptual/illustrative here, real once you have a real cluster

**Deployment Strategies:**
- Rolling (gradual, no extra infra) · Blue-Green (instant rollback, double infra) · Canary (smallest blast radius, needs traffic-splitting)

---

## Next Steps

You now understand:
- ✅ What CI and CD each mean, and the problem they solve
- ✅ GitHub Actions' workflow/job/step/runner anatomy
- ✅ A complete, correct build+vet+test workflow - and why the order matters
- ✅ gofmt, race+coverage, and linter quality gates - including gofmt -l's real exit-code gotcha
- ✅ Matrix builds across Go versions
- ✅ What dependency caching actually caches, and why
- ✅ How a Docker build-and-push job and a Kubernetes deploy job fit into the same pipeline
- ✅ Rolling, blue-green, and canary deployment strategies at a conceptual level

**Next level:** Level 38 - Production Debugging
- Diagnosing issues in a system that's already running and already shipped
- Reading logs, metrics, and traces under real production pressure
- Everything from Levels 33 (Logging) through 37 (CI/CD) comes together once something goes wrong at 3 AM

Your code doesn't just work - now it proves itself automatically, on every single change. Keep going! 🚀
