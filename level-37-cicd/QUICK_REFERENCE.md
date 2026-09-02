# Level 37: Quick Reference Card

> ⚠️ **Verification scope:** Go commands below (`go build`/`vet`/`test`, `gofmt`, `golangci-lint`) were genuinely run in this environment - real output. Workflow YAML was validated for real syntax/structure with a YAML parser. No workflow ran on an actual GitHub runner, and the Docker-push/K8s-deploy snippets are illustrative (same disclaimer as Levels 35-36). Full details in README.md's introduction.

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir myapp && cd myapp
go mod init myapp
mkdir -p .github/workflows

# The minimal real CI workflow
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

# Validate the YAML syntax locally (this is real, always run this yourself)
python3 -c "import yaml; yaml.safe_load(open('.github/workflows/ci.yml'))"
```

---

## 📋 Workflow Anatomy

```yaml
name: CI                        # label shown in GitHub's UI

on:                              # TRIGGER
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:                            # jobs run in PARALLEL unless linked with needs:
  <job-name>:
    runs-on: ubuntu-latest        # the RUNNER
    steps:                        # steps run SEQUENTIALLY
      - uses: owner/action@v1     # a marketplace action
      - run: some-shell-command   # a raw shell command
```

---

## 🔁 The Core Gate Sequence (Fail Fast Order)

```yaml
- uses: actions/checkout@v4
- uses: actions/setup-go@v5
  with:
    go-version: '1.23'
    cache: true

- name: Check formatting (gofmt)
  run: |
    unformatted=$(gofmt -l .)
    if [ -n "$unformatted" ]; then
      echo "$unformatted"
      exit 1
    fi

- run: go build ./...
- run: go vet ./...

- uses: golangci/golangci-lint-action@v6
  with:
    version: latest

- run: go test -race -cover ./...
```

---

## 🚨 The gofmt Gotcha (Memorize This)

```bash
gofmt -l .           # lists unformatted files... and STILL exits 0!
echo $?               # → 0, even when files were listed

# ✅ Correct gate: check the OUTPUT, not the exit code
test -z "$(gofmt -l .)"    # fails (non-zero) only when gofmt -l printed something
```

---

## 🧩 Matrix Builds

```yaml
strategy:
  matrix:
    go-version: ['1.22', '1.23']
steps:
  - uses: actions/setup-go@v5
    with:
      go-version: ${{ matrix.go-version }}
```
Runs the whole job once per list entry, in parallel. Add a second key (`os: [...]`) for the full cross-product.

---

## 📦 Caching

```yaml
- uses: actions/setup-go@v5
  with:
    go-version: '1.23'
    cache: true          # caches GOMODCACHE (~/go/pkg/mod) and GOCACHE (~/.cache/go-build)
                          # keyed by go.sum's hash
```

---

## 🐳 Docker Job (Bridging Level 35)

```yaml
docker:
  needs: build-and-test          # only runs after tests pass
  runs-on: ubuntu-latest
  steps:
    - uses: docker/login-action@v3
      with:
        username: ${{ secrets.DOCKERHUB_USERNAME }}
        password: ${{ secrets.DOCKERHUB_TOKEN }}
    - uses: docker/build-push-action@v6
      with:
        context: .
        push: true
        tags: mycompany/app:${{ github.sha }}
```

## ☸️ Deploy Job Sketch (Bridging Level 36)

```yaml
deploy:
  needs: docker
  runs-on: ubuntu-latest
  steps:
    - run: kubectl set image deployment/app app=mycompany/app:${{ github.sha }}
    - run: kubectl rollout status deployment/app --timeout=120s
```

---

## 🚦 Deployment Strategies

| Strategy | Rollback | Extra infra |
|----------|----------|-------------|
| Rolling | Moderate | None |
| Blue-Green | Instant | Double |
| Canary | Fast | Traffic-splitting |

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Trusting gofmt's exit code | `gofmt -l .` alone | `test -z "$(gofmt -l .)"` |
| Unpinned actions | `uses: actions/checkout@main` | `uses: actions/checkout@v4` |
| Hardcoded secrets | `-p hunter2` in the YAML | `-p ${{ secrets.TOKEN }}` |
| No caching | `cache: false` | `cache: true` |
| Testing only on merge | `on: push` only | `on: push` + `on: pull_request` |

---

## 🎓 Before Next Level

Can you:
- [ ] Write a correct `on:`/`jobs:`/`steps:`/`runs-on:` workflow from memory?
- [ ] Explain why `gofmt -l .`'s exit code alone can't gate a build?
- [ ] Explain what `go test -race` catches that a normal run misses?
- [ ] Explain what a matrix build does and what it costs (parallel runner-minutes)?
- [ ] Explain what a Go CI cache actually holds?
- [ ] Compare rolling, blue-green, and canary deployments?

If YES → You're ready for Level 38!

---

## 📚 Next Level

Level 38: Production Debugging
- Diagnosing issues in a system that's already running and already shipped
- Reading logs, metrics, and traces under real production pressure

Your pipeline ships itself now - onward! 💪
