# Level 37: CI/CD - Exercises

> ⚠️ **Verification scope for this level.** Every Go command shown in these exercises (`go build`, `go vet`, `gofmt`, `go test -race -cover`, `golangci-lint run`) was **genuinely executed** in a real scratch Go project while writing this level - the "Expected Output" blocks are real, captured output, not invented. `golangci-lint` (v1.64.8) happened to be installed in this environment, so its output is real too. Every `.github/workflows/*.yml` file shown was validated for real syntactic and structural correctness with a YAML parser (`python3 -c "import yaml; yaml.safe_load(...)"`). What was **not** exercised: an actual GitHub-hosted run of any workflow (no git repository exists in this course, and GitHub's own runner infrastructure isn't reachable from here), and the Docker/Kubernetes portions of Exercises 8 and 9 - which follow the same illustrative disclaimer Level 35 (Docker) and Level 36 (Kubernetes) already gave you: `docker` and `kubectl` client binaries exist here, but no daemon/cluster is reachable, so those two exercises are conceptual sketches, not verified end-to-end.
>
> This environment also has exactly **one** Go toolchain installed (`go1.26.5`) and no second version to install alongside it - Exercise 6 (matrix builds) is explicit about exactly what that does and doesn't limit.

Complete all exercises in order. Each exercise builds on previous knowledge.

---

## Exercise 1: A Basic Build + Vet + Test Workflow

**Objective:** Write a real Go package, verify it locally with the exact commands a CI workflow would run, then write the workflow that automates them.

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level37-exercise1
cd ~/projects/level37-exercise1
go mod init level37.example/exercise1
mkdir -p mathutil
```

2. Create `mathutil/mathutil.go`:

```bash
cat > mathutil/mathutil.go << 'EOF'
package mathutil

// Add returns the sum of two integers.
func Add(a, b int) int {
    return a + b
}

// Multiply returns the product of two integers.
func Multiply(a, b int) int {
    return a * b
}
EOF
```

3. Create `mathutil/mathutil_test.go`:

```bash
cat > mathutil/mathutil_test.go << 'EOF'
package mathutil

import "testing"

func TestAdd(t *testing.T) {
    if got := Add(2, 3); got != 5 {
        t.Errorf("Add(2, 3) = %d; want 5", got)
    }
}

func TestMultiply(t *testing.T) {
    if got := Multiply(4, 5); got != 20 {
        t.Errorf("Multiply(4, 5) = %d; want 20", got)
    }
}
EOF
```

4. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"

    "level37.example/exercise1/mathutil"
)

func main() {
    fmt.Println("2 + 3 =", mathutil.Add(2, 3))
    fmt.Println("4 * 5 =", mathutil.Multiply(4, 5))
}
EOF
```

5. Run exactly what a CI job would run, locally:

```bash
go build ./...
go vet ./...
go test ./...
go test -v ./...
```

**Expected Output (real, verified with `go build`/`go vet`/`go test` in this environment):**

```
$ go build ./...
$ echo $?
0

$ go vet ./...
$ echo $?
0

$ go test ./...
?   	level37.example/exercise1	[no test files]
ok  	level37.example/exercise1/mathutil	0.525s

$ go test -v ./...
?   	level37.example/exercise1	[no test files]
=== RUN   TestAdd
--- PASS: TestAdd (0.00s)
=== RUN   TestMultiply
--- PASS: TestMultiply (0.00s)
PASS
ok  	level37.example/exercise1/mathutil	0.393s
```

6. Create the workflow file:

```bash
mkdir -p .github/workflows
cat > .github/workflows/ci.yml << 'EOF'
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
EOF
```

**Validated (syntax only - not run on a GitHub runner in this environment):**

```
$ python3 -c "import yaml; yaml.safe_load(open('.github/workflows/ci.yml'))"
(no output = valid YAML)
```

**Learning Objectives:**
- ✅ Run the exact three commands (`go build`, `go vet`, `go test`) a CI job runs, locally, and see real output
- ✅ Write a correct GitHub Actions workflow using `actions/checkout` and `actions/setup-go`
- ✅ Understand the difference between "this Go code was verified" and "this YAML is syntactically valid" - both true here, for different reasons

---

## Exercise 2: The gofmt Quality Gate (and Its Exit-Code Gotcha)

**Objective:** Discover, first-hand, that `gofmt -l` does not fail on unformatted code by exit code alone - and write a CI step that actually catches it

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level37-exercise2
cd ~/projects/level37-exercise2
go mod init level37.example/exercise2
```

2. Create a deliberately badly-formatted `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
        x :=    1
	y := 2
    fmt.Println(  "sum:", x+y )
}
EOF
```

3. Run `gofmt -l` and check its exit code:

```bash
gofmt -l .
echo "exit code: $?"
```

**Expected Output (real):**

```
main.go
exit code: 0
```

🚨 Notice: `main.go` WAS listed as needing formatting, and the exit code is still `0`. A CI step that just runs `gofmt -l .` and relies on its own exit code will **never fail**, no matter how mangled the formatting is.

4. See the actual diff `gofmt` would apply, then fix it for real:

```bash
gofmt -d .
gofmt -w .
gofmt -l .
echo "exit code: $?"
cat main.go
```

**Expected Output (real):**

```
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

exit code: 0
package main

import "fmt"

func main() {
	x := 1
	y := 2
	fmt.Println("sum:", x+y)
}
```

After `gofmt -w .`, `gofmt -l .` prints nothing at all (still exit `0` - but now correctly meaning "nothing to report").

5. Write a workflow step that checks for **empty output**, not exit code:

```bash
mkdir -p .github/workflows
cat > .github/workflows/ci.yml << 'EOF'
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
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Check formatting (gofmt)
        run: |
          unformatted=$(gofmt -l .)
          if [ -n "$unformatted" ]; then
            echo "The following files are not gofmt-formatted:"
            echo "$unformatted"
            exit 1
          fi

      - run: go build ./...
      - run: go vet ./...
      - run: go test ./...
EOF
python3 -c "import yaml; yaml.safe_load(open('.github/workflows/ci.yml'))"
```

**Learning Objectives:**
- ✅ Discover `gofmt -l`'s real, documented exit-code behavior first-hand instead of assuming it "just works" as a gate
- ✅ Fix formatting with `gofmt -w` and confirm the fix with `gofmt -l`
- ✅ Write a CI step that checks for empty output explicitly, the correct way to gate on `gofmt`

---

## Exercise 3: Race Detector and Coverage

**Objective:** See `go test -race -cover` both pass cleanly on correct code and genuinely catch a data race on buggy code

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level37-exercise3
cd ~/projects/level37-exercise3
go mod init level37.example/exercise3
mkdir -p counter
```

2. Create a concurrency-safe `counter/counter.go`:

```bash
cat > counter/counter.go << 'EOF'
package counter

import "sync"

// Counter is a concurrency-safe counter.
type Counter struct {
    mu    sync.Mutex
    value int
}

// Increment adds 1 to the counter's value.
func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

// Value returns the current count.
func (c *Counter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}
EOF
```

3. Create `counter/counter_test.go`:

```bash
cat > counter/counter_test.go << 'EOF'
package counter

import (
    "sync"
    "testing"
)

func TestCounterConcurrentIncrement(t *testing.T) {
    c := &Counter{}
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            c.Increment()
        }()
    }
    wg.Wait()

    if got := c.Value(); got != 100 {
        t.Errorf("Value() = %d; want 100", got)
    }
}

func TestCounterStartsAtZero(t *testing.T) {
    c := &Counter{}
    if got := c.Value(); got != 0 {
        t.Errorf("Value() = %d; want 0", got)
    }
}
EOF
```

4. Run with the race detector and coverage:

```bash
go test -v -race -cover ./...
```

**Expected Output (real):**

```
=== RUN   TestCounterConcurrentIncrement
--- PASS: TestCounterConcurrentIncrement (0.00s)
=== RUN   TestCounterStartsAtZero
--- PASS: TestCounterStartsAtZero (0.00s)
PASS
coverage: 100.0% of statements
ok  	level37.example/exercise3/counter	1.269s	coverage: 100.0% of statements
```

5. Now see what `-race` actually catches. In a **separate** scratch directory, recreate `counter.go` **without** the mutex:

```bash
mkdir -p ~/projects/level37-exercise3-racedemo/counter
cd ~/projects/level37-exercise3-racedemo
go mod init level37.example/racedemo

cat > counter/counter.go << 'EOF'
package counter

// Counter is deliberately NOT concurrency-safe, to demonstrate what
// go test -race catches.
type Counter struct {
    value int
}

func (c *Counter) Increment() {
    c.value++ // BUG: unsynchronized read-modify-write
}

func (c *Counter) Value() int {
    return c.value
}
EOF
```

Reuse the same `counter_test.go` from step 3, then run:

```bash
go test -race ./...
```

**Expected Output (real, machine-dependent addresses/goroutine numbers/final count - the SHAPE is what matters):**

```
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

**Learning Objectives:**
- ✅ Run `go test -race -cover` against correct, mutex-protected concurrent code
- ✅ See a real, genuine data-race report from the unsynchronized version - not a description of one
- ✅ Understand why `-race` belongs in CI even though it makes tests slower: it catches bugs that pass silently without it

---

## Exercise 4: Adding a Linter (golangci-lint)

**Objective:** Run `golangci-lint` for real, see it catch a genuine issue, then fix it

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level37-exercise4
cd ~/projects/level37-exercise4
go mod init level37.example/exercise4
```

2. Create `main.go` with a real, easy-to-miss bug - an unchecked error return:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "os"
)

func writeGreeting(name string) {
    f, _ := os.Create("/tmp/greeting.txt")
    f.WriteString("Hello, " + name)
    f.Close()
}

func greet(name string) string {
    return fmt.Sprintf("Hello, %s!", name)
}

func main() {
    result := greet("Gopher")
    fmt.Println(result)
    writeGreeting("Gopher")
}
EOF
```

3. This compiles fine - `go build` and `go vet` won't catch it. Run the linter:

```bash
go build ./...
golangci-lint run ./...
echo "lint exit code: $?"
```

**Expected Output (real - if you don't have `golangci-lint` installed, see https://golangci-lint.run/welcome/install/ to add it):**

```
main.go:10:18: Error return value of `f.WriteString` is not checked (errcheck)
    f.WriteString("Hello, " + name)
                 ^
lint exit code: 1
```

4. Fix it properly - check and propagate the error, the idiomatic pattern from `Level 16: Error Handling`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "log"
    "os"
)

func writeGreeting(name string) error {
    f, err := os.Create("/tmp/greeting.txt")
    if err != nil {
        return err
    }
    defer f.Close()

    if _, err := f.WriteString("Hello, " + name); err != nil {
        return err
    }
    return nil
}

func greet(name string) string {
    return fmt.Sprintf("Hello, %s!", name)
}

func main() {
    result := greet("Gopher")
    fmt.Println(result)
    if err := writeGreeting("Gopher"); err != nil {
        log.Fatal(err)
    }
}
EOF
go build ./...
go vet ./...
golangci-lint run ./...
echo "lint exit code: $?"
go run main.go
```

**Expected Output (real):**

```
lint exit code: 0
Hello, Gopher!
```

5. Add the linter to the workflow:

```bash
mkdir -p .github/workflows
cat > .github/workflows/ci.yml << 'EOF'
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
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - run: go build ./...
      - run: go vet ./...
      - name: Lint
        uses: golangci/golangci-lint-action@v6
        with:
          version: latest
      - run: go test ./...
EOF
python3 -c "import yaml; yaml.safe_load(open('.github/workflows/ci.yml'))"
```

**Learning Objectives:**
- ✅ See `golangci-lint` catch a real bug that `go build` and `go vet` both miss
- ✅ Fix an unchecked-error lint finding the idiomatic way
- ✅ Wire `golangci-lint-action` into a workflow as its own quality-gate step

---

## Exercise 5: Triggers - on: push vs on: pull_request

**Objective:** Configure a workflow to run on both push and pull_request events, and understand what each one means

**Instructions:**

1. Create working directory (reuse Exercise 1's project, or a fresh one):

```bash
mkdir -p ~/projects/level37-exercise5
cd ~/projects/level37-exercise5
go mod init level37.example/exercise5
mkdir -p .github/workflows
```

2. Write a workflow with both triggers:

```bash
cat > .github/workflows/ci.yml << 'EOF'
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
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - run: go build ./...
      - run: go vet ./...
      - run: go test ./...
EOF
```

3. Validate the YAML:

```bash
python3 -c "
import yaml
data = yaml.safe_load(open('.github/workflows/ci.yml'))
print('top-level keys:', list(data.keys()))
print('jobs:', list(data['jobs'].keys()))
"
```

**Expected Output (real):**

```
top-level keys: ['name', True, 'jobs']
jobs: ['build-and-test']
```

🚨 Notice `True` in that key list, not `'on'`. This is a genuine YAML 1.1 quirk: a generic parser (like PyYAML's `safe_load`) reads the bare word `on` as the boolean `true`. GitHub's own Actions parser is not confused by this - it always treats a workflow's top-level `on:` key as the trigger definition - but it's a real trap if you ever try to read a workflow's trigger with a plain YAML library instead of GitHub's own tooling.

**What each trigger means:**
- `push: branches: [main]` - runs whenever a commit lands directly on `main` (including a merge).
- `pull_request: branches: [main]` - runs whenever a PR **targeting** `main` is opened, or updated with a new commit - this is what lets you see CI results **before** a merge happens at all.

**Learning Objectives:**
- ✅ Configure both `push` and `pull_request` triggers correctly
- ✅ Understand the practical difference: `pull_request` catches problems before merge, `push` alone catches them only after
- ✅ Encounter the real `on:` → `True` YAML parsing quirk instead of just reading about it

---

## Exercise 6: Matrix Builds

**Objective:** Write a correct matrix-build workflow, and be honest about what could and couldn't be verified locally

**Instructions:**

1. Check what Go versions are actually available in your environment:

```bash
go version
```

**Expected Output (real, in the environment that authored this course):**

```
go version go1.26.5 darwin/arm64
```

Only **one** Go toolchain (`1.26.5`) was available here - there was no second version installed to genuinely test a real multi-version matrix locally. If you have `gvm`, `asdf`, or another Go version manager, you can install a second version and repeat step 2 below against each one to get closer to a real matrix verification than this environment could.

2. Write the matrix workflow (the concept - "run this job once per listed Go version" - is exactly what GitHub Actions does per its documentation; this file's syntax is validated, but a real multi-version run was not exercised here):

```bash
mkdir -p .github/workflows
cat > .github/workflows/ci.yml << 'EOF'
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  build-and-test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: ['1.22', '1.23']
    steps:
      - uses: actions/checkout@v4
      - name: Set up Go ${{ matrix.go-version }}
        uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go-version }}
      - run: go build ./...
      - run: go vet ./...
      - run: go test ./...
EOF
python3 -c "
import yaml
data = yaml.safe_load(open('.github/workflows/ci.yml'))
job = data['jobs']['build-and-test']
print('matrix go-version list:', job['strategy']['matrix']['go-version'])
"
```

**Expected Output (real - confirms the matrix structure parses correctly):**

```
matrix go-version list: ['1.22', '1.23']
```

3. Run the build/vet/test sequence with the one Go version genuinely available here, so at least that slice of the matrix has real output behind it:

```bash
go build ./...
go vet ./...
go test ./...
```

(Reuse Exercise 1's `mathutil` project for this - the exact commands and output are already shown there.)

**Learning Objectives:**
- ✅ Write a syntactically correct `strategy.matrix` block
- ✅ Understand that a matrix runs the identical job once per entry, in parallel, with `${{ matrix.<key> }}` substituted in
- ✅ Practice being explicit about verification scope: what was genuinely tested (one Go version, locally) vs. what the matrix concept promises (multiple versions, on GitHub's runners) but this environment couldn't confirm

---

## Exercise 7: Understanding Dependency Caching

**Objective:** Identify exactly what a Go CI cache holds, using your own machine's real cache paths

**Instructions:**

1. Find your own real cache directories:

```bash
go env GOCACHE
go env GOMODCACHE
```

**Expected Output (real, from the environment that authored this course - your own paths will differ by OS/user, but the two cache KINDS are universal):**

```
/Users/<you>/Library/Caches/go-build
/Users/<you>/go/pkg/mod
```

(On a GitHub Actions `ubuntu-latest` runner, the build cache path is `~/.cache/go-build` instead - same purpose, different OS convention.)

2. Write a workflow using `actions/setup-go`'s built-in caching:

```bash
mkdir -p .github/workflows
cat > .github/workflows/ci.yml << 'EOF'
name: CI

on:
  push:
    branches: [main]

jobs:
  build-and-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'
          cache: true
      - run: go build ./...
      - run: go test ./...
EOF
```

3. Now write the manual equivalent with `actions/cache`, to see exactly what's being cached under the hood:

```bash
cat > .github/workflows/ci-manual-cache.yml << 'EOF'
name: CI

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
          cache: false

      - name: Cache Go modules and build cache
        uses: actions/cache@v4
        with:
          path: |
            ~/.cache/go-build
            ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
          restore-keys: |
            ${{ runner.os }}-go-

      - run: go build ./...
      - run: go test ./...
EOF
python3 -c "import yaml; yaml.safe_load(open('.github/workflows/ci.yml'))"
python3 -c "import yaml; yaml.safe_load(open('.github/workflows/ci-manual-cache.yml'))"
```

**Learning Objectives:**
- ✅ Locate your own machine's real Go module and build cache directories
- ✅ Understand that `cache: true` and the manual `actions/cache` block do the same thing - restore the module download cache and the compiled-package build cache, keyed by `go.sum`'s hash
- ✅ Recognize that this cache speeds up *repeated* runs, not the first one - the first run on a new key always starts cold

---

## Exercise 8: Docker Build-and-Push Job (Bridging Level 35)

**Objective:** Write a correct, syntactically valid job that builds and pushes a Docker image in CI

> ⚠️ Per `Level 35: Docker`'s own disclaimer, `docker` is present in this environment as a client binary but its daemon is not reachable (`docker version`'s server section fails to connect) - so, exactly like that level, this exercise's Docker portion is illustrative: correct per Docker's and GitHub Actions' documented behavior, but not run end-to-end here.

**Instructions:**

1. Reuse Exercise 1's Go project (it already has a working `main.go` and passes build/vet/test).

2. Add a multi-stage `Dockerfile` (from `Level 35: Docker`, Section 3):

```bash
cat > Dockerfile << 'EOF'
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
EOF
```

3. Write the workflow with a `docker` job gated behind the test job:

```bash
mkdir -p .github/workflows
cat > .github/workflows/docker-build-push.yml << 'EOF'
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
EOF
python3 -c "
import yaml
data = yaml.safe_load(open('.github/workflows/docker-build-push.yml'))
print('jobs:', list(data['jobs'].keys()))
print('docker job needs:', data['jobs']['docker']['needs'])
"
```

**Expected Output (real - confirms the job graph parses as intended):**

```
jobs: ['build-and-test', 'docker']
docker job needs: build-and-test
```

**Learning Objectives:**
- ✅ Gate a packaging job behind a testing job with `needs:`
- ✅ Reference registry credentials via `${{ secrets.* }}`, never hardcoded
- ✅ Recognize which claims here are verified (the YAML's structure and job graph) vs. illustrative (an actual image being built and pushed)

---

## Exercise 9: A Conceptual Kubernetes Deploy Job (Bridging Level 36)

**Objective:** Write a correct, illustrative deploy job that updates a Kubernetes Deployment's image after a successful build

> ⚠️ Per `Level 36: Kubernetes`'s own disclaimer, no real cluster (and no working `kubectl`-to-cluster connection) exists in this environment - this exercise is a conceptual sketch, correct per Kubernetes' documented `kubectl` behavior, not verified against a live cluster.

**Instructions:**

1. Extend Exercise 8's workflow with a `deploy` job:

```bash
cat >> .github/workflows/docker-build-push.yml << 'EOF'

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
EOF
python3 -c "
import yaml
data = yaml.safe_load(open('.github/workflows/docker-build-push.yml'))
print('jobs:', list(data['jobs'].keys()))
print('deploy job needs:', data['jobs']['deploy']['needs'])
"
```

**Expected Output (real - the YAML parses correctly with all three jobs chained):**

```
jobs: ['build-and-test', 'docker', 'deploy']
deploy job needs: docker
```

2. Compare `kubectl set image` to what you did by hand in `Level 36: Kubernetes`, Section 9 (Rolling Updates) - this is the exact same operation, just scripted instead of typed interactively.

**Learning Objectives:**
- ✅ Chain a three-job pipeline (test → build image → deploy) with `needs:`
- ✅ See how `kubectl set image` + `kubectl rollout status` automate what Level 36 did by hand
- ✅ Understand exactly why this exercise is labeled conceptual: no cluster exists here to actually apply it against

---

## Exercise 10: Comprehensive Practice - A Complete CI Pipeline for a REST API

**Objective:** Build a small, real REST API (in the style of `Level 27: HTTP & REST APIs`), verify it with every quality gate from this level, and write the one `ci.yml` that would run all of them

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level37-exercise10
cd ~/projects/level37-exercise10
go mod init level37.example/exercise10
```

2. Create `tasks.go` - a tiny in-memory task-list API:

```bash
cat > tasks.go << 'EOF'
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
)

// Task is a single to-do item.
type Task struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Done bool   `json:"done"`
}

// TaskStore is a concurrency-safe in-memory task list.
type TaskStore struct {
    mu     sync.Mutex
    tasks  map[int]Task
    nextID int
}

// NewTaskStore returns an empty TaskStore.
func NewTaskStore() *TaskStore {
    return &TaskStore{
        tasks:  make(map[int]Task),
        nextID: 1,
    }
}

// Add creates a new task and returns it.
func (s *TaskStore) Add(name string) Task {
    s.mu.Lock()
    defer s.mu.Unlock()

    t := Task{ID: s.nextID, Name: name}
    s.tasks[t.ID] = t
    s.nextID++
    return t
}

// List returns all tasks.
func (s *TaskStore) List() []Task {
    s.mu.Lock()
    defer s.mu.Unlock()

    result := make([]Task, 0, len(s.tasks))
    for _, t := range s.tasks {
        result = append(result, t)
    }
    return result
}

// tasksHandler handles GET (list) and POST (create) on /tasks.
func tasksHandler(store *TaskStore) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            w.Header().Set("Content-Type", "application/json")
            if err := json.NewEncoder(w).Encode(store.List()); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
            }
        case http.MethodPost:
            var body struct {
                Name string `json:"name"`
            }
            if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
                http.Error(w, "invalid request body", http.StatusBadRequest)
                return
            }
            if body.Name == "" {
                http.Error(w, "name is required", http.StatusBadRequest)
                return
            }
            t := store.Add(body.Name)
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusCreated)
            if err := json.NewEncoder(w).Encode(t); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
            }
        default:
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    }
}

func main() {
    store := NewTaskStore()
    http.HandleFunc("/tasks", tasksHandler(store))
    fmt.Println("listening on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println("server error:", err)
    }
}
EOF
```

3. Create `tasks_test.go`:

```bash
cat > tasks_test.go << 'EOF'
package main

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

func TestTaskStoreAddAndList(t *testing.T) {
    store := NewTaskStore()
    store.Add("write tests")
    store.Add("ship it")

    got := store.List()
    if len(got) != 2 {
        t.Fatalf("List() returned %d tasks; want 2", len(got))
    }
}

func TestTasksHandlerPostCreatesTask(t *testing.T) {
    store := NewTaskStore()
    handler := tasksHandler(store)

    req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{"name":"learn CI/CD"}`))
    rec := httptest.NewRecorder()

    handler(rec, req)

    if rec.Code != http.StatusCreated {
        t.Fatalf("status = %d; want %d", rec.Code, http.StatusCreated)
    }

    var got Task
    if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
        t.Fatalf("decoding response: %v", err)
    }
    if got.Name != "learn CI/CD" {
        t.Errorf("Name = %q; want %q", got.Name, "learn CI/CD")
    }
    if got.ID != 1 {
        t.Errorf("ID = %d; want 1", got.ID)
    }
}

func TestTasksHandlerPostRejectsEmptyName(t *testing.T) {
    store := NewTaskStore()
    handler := tasksHandler(store)

    req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{"name":""}`))
    rec := httptest.NewRecorder()

    handler(rec, req)

    if rec.Code != http.StatusBadRequest {
        t.Errorf("status = %d; want %d", rec.Code, http.StatusBadRequest)
    }
}

func TestTasksHandlerGetListsTasks(t *testing.T) {
    store := NewTaskStore()
    store.Add("first task")

    handler := tasksHandler(store)
    req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
    rec := httptest.NewRecorder()

    handler(rec, req)

    if rec.Code != http.StatusOK {
        t.Fatalf("status = %d; want %d", rec.Code, http.StatusOK)
    }

    var got []Task
    if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
        t.Fatalf("decoding response: %v", err)
    }
    if len(got) != 1 {
        t.Fatalf("got %d tasks; want 1", len(got))
    }
}

func TestTasksHandlerRejectsUnsupportedMethod(t *testing.T) {
    store := NewTaskStore()
    handler := tasksHandler(store)

    req := httptest.NewRequest(http.MethodDelete, "/tasks", nil)
    rec := httptest.NewRecorder()

    handler(rec, req)

    if rec.Code != http.StatusMethodNotAllowed {
        t.Errorf("status = %d; want %d", rec.Code, http.StatusMethodNotAllowed)
    }
}
EOF
```

4. Run every quality gate from this level, in fail-fast order:

```bash
gofmt -l .
gofmt -w .
gofmt -l .
go build ./...
go vet ./...
golangci-lint run ./...
go test -v -race -cover ./...
```

**Expected Output (real, exactly as captured while writing this level):**

```
$ gofmt -l .
tasks.go
tasks_test.go

$ gofmt -w .

$ gofmt -l .
(empty - nothing left to format)

$ go build ./...
(exit code: 0)

$ go vet ./...
(exit code: 0)

$ golangci-lint run ./...
(exit code: 0)

$ go test -v -race -cover ./...
=== RUN   TestTaskStoreAddAndList
--- PASS: TestTaskStoreAddAndList (0.00s)
=== RUN   TestTasksHandlerPostCreatesTask
--- PASS: TestTasksHandlerPostCreatesTask (0.00s)
=== RUN   TestTasksHandlerPostRejectsEmptyName
--- PASS: TestTasksHandlerPostRejectsEmptyName (0.00s)
=== RUN   TestTasksHandlerGetListsTasks
--- PASS: TestTasksHandlerGetListsTasks (0.00s)
=== RUN   TestTasksHandlerRejectsUnsupportedMethod
--- PASS: TestTasksHandlerRejectsUnsupportedMethod (0.00s)
PASS
coverage: 75.0% of statements
ok  	level37.example/exercise10	1.753s	coverage: 75.0% of statements
```

Notice step 4 reproduces this level's Exercise 2 lesson on real code: the heredoc-created files needed `gofmt -w` before they were clean - exactly the kind of thing a real `gofmt` CI gate exists to catch automatically, on every change, without anyone having to remember to run it by hand.

5. Write the complete `ci.yml` covering every gate:

```bash
mkdir -p .github/workflows
cat > .github/workflows/ci.yml << 'EOF'
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
          cache: true

      - name: Check formatting (gofmt)
        run: |
          unformatted=$(gofmt -l .)
          if [ -n "$unformatted" ]; then
            echo "The following files are not gofmt-formatted:"
            echo "$unformatted"
            exit 1
          fi

      - name: Build
        run: go build ./...

      - name: Vet
        run: go vet ./...

      - name: Lint
        uses: golangci/golangci-lint-action@v6
        with:
          version: latest

      - name: Test with race detector and coverage
        run: go test -v -race -cover ./...
EOF
python3 -c "
import yaml
data = yaml.safe_load(open('.github/workflows/ci.yml'))
steps = data['jobs']['build-and-test']['steps']
print('number of steps:', len(steps))
print('step names:', [s.get('name', s.get('uses', s.get('run'))) for s in steps])
"
```

**Expected Output (real - confirms all seven steps parsed correctly, in the intended order):**

```
number of steps: 7
step names: ['Check out code', 'Set up Go', 'Check formatting (gofmt)', 'Build', 'Vet', 'Lint', 'Test with race detector and coverage']
```

**Learning Objectives:**
- ✅ Build and fully test a small, real REST API using only the standard library and `httptest` (Level 27 patterns)
- ✅ Run every quality gate from this level - gofmt, build, vet, lint, race+coverage - against one real project, in the fail-fast order recommended in the README
- ✅ Write and validate the single, complete `ci.yml` a real small service would actually use

---

## Bonus Challenges

### Challenge 1: Different Steps for push vs pull_request

Using `github.event_name` in an `if:` condition, write a workflow where a plain `push` to `main` runs the fast test suite (`go test ./...`), while a `pull_request` run additionally runs the slower race-and-coverage suite (`go test -race -cover ./...`).

```bash
mkdir -p ~/projects/level37-bonus1/.github/workflows
cd ~/projects/level37-bonus1
go mod init level37.example/bonus1
```

**Hints:**
- A step's `if:` field can reference `github.event_name` directly: `if: github.event_name == 'pull_request'`
- You'll need two separate `Test` steps, each gated by a different `if:`, rather than one step that branches internally
- Validate your result the same way every exercise in this level did: `python3 -c "import yaml; yaml.safe_load(open('.github/workflows/ci.yml'))"`

### Challenge 2: A Release-Tagging Workflow

Sketch a workflow that triggers only on version tags (`v1.2.3`-style) and cross-compiles the binary for multiple `GOOS`/`GOARCH` combinations using a matrix, uploading each as a build artifact.

```bash
mkdir -p ~/projects/level37-bonus2/.github/workflows
cd ~/projects/level37-bonus2
go mod init level37.example/bonus2
```

**Hints:**
- `on: push: tags: ['v*.*.*']` triggers only on tags matching that pattern, never on ordinary branch pushes
- `strategy.matrix` can carry two independent lists (`goos: [linux, darwin, windows]` and `goarch: [amd64, arm64]`) - Actions runs the full cross-product automatically
- Set `GOOS`/`GOARCH` as step-level `env:` values sourced from `${{ matrix.goos }}` / `${{ matrix.goarch }}`
- `actions/upload-artifact@v4` is the standard way to save a build output from a job

### Challenge 3: A Status Badge for Your README

Add a Markdown snippet to a project's `README.md` that shows a live "build passing/failing" badge sourced from a workflow named `ci.yml`.

```bash
mkdir -p ~/projects/level37-bonus3
cd ~/projects/level37-bonus3
go mod init level37.example/bonus3
```

**Hints:**
- GitHub generates a badge image automatically at `https://github.com/<owner>/<repo>/actions/workflows/ci.yml/badge.svg`
- Wrap it in a Markdown image inside a link to the workflow's Actions page: `[![CI](.../ci.yml/badge.svg)](.../actions/workflows/ci.yml)`
- This only renders a real "passing"/"failing" badge once the workflow has actually run at least once on a real GitHub repository - there is nothing to verify locally here, which is exactly why this is a bonus, hints-only challenge rather than one with "Expected Output"

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Write a real Go project and verify it locally with exactly the commands a CI workflow would run
✅ Know `gofmt -l`'s true exit-code behavior and gate on empty output correctly
✅ Read and trust real `go test -race -cover` output, including a genuine data-race report
✅ Add `golangci-lint` as a quality gate and fix a real finding
✅ Configure `push` and `pull_request` triggers, and recognize the `on:`/YAML-boolean quirk
✅ Write a syntactically correct matrix-build workflow, while being explicit about what a single-Go-version environment can and can't verify
✅ Explain exactly what a Go CI cache holds and why it speeds up repeated runs
✅ Sketch a gated, multi-job pipeline: test → build & push a Docker image → deploy to Kubernetes
✅ Compare rolling, blue-green, and canary deployment strategies
✅ Assemble every gate into one complete, validated `ci.yml` for a real small REST API

---

## Next Level

Level 38: Production Debugging
- Diagnosing issues in a system that's already running and already shipped
- Reading logs, metrics, and traces under real production pressure
- Everything from Levels 33 (Logging) through 37 (CI/CD) comes together once something goes wrong at 3 AM

Great work! Your code now proves itself automatically, on every single change! 🚀
