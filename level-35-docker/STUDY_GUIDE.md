# Level 35: Study Guide & Visual Reference

> ⚠️ **Reminder:** Docker's daemon was not reachable in the environment that authored this course (`docker version` returned a client but no server connection). The Go source in this level was genuinely built/run and verified; the diagrams, tables, and command output below are correct per Docker's documented behavior but are **illustrative, not captured from a real build**. See EXERCISES.md's banner for the full explanation.

## 📚 Learning Path

### Week 1: Dockerfiles and Multi-Stage Builds
```
Day 1:  Why containerize a Go app; Dockerfile basics (FROM/WORKDIR/COPY/RUN/EXPOSE/CMD)
Day 2:  CMD vs ENTRYPOINT; a first single-stage Dockerfile
Day 3:  Multi-stage builds; choosing a runtime base (scratch/distroless/alpine)
Day 4:  .dockerignore and build-context hygiene
Day 5:  Layer caching order (go.mod first, source second)
Day 6:  Building and running; docker build / docker run flags
Day 7:  Environment variables in containers (ENV, -e, Level 34 tie-in)
```

### Week 2: Compose and Practice
```
Day 1:  Exercises 1-3 (single-stage, multi-stage, .dockerignore)
Day 2:  Exercises 4-6 (layer caching, port mapping, ENV defaults)
Day 3:  Exercises 7-8 (dockerignore audit, docker-compose with a database)
Day 4:  Exercises 9-10 (HEALTHCHECK, comprehensive REST API containerization)
Day 5:  Bonus challenges
Day 6-7: Practice & review
```

---

## 🐳 Multi-Stage Build Diagram

```
┌─────────────────────────────────────────────────────────────┐
│  STAGE 1: "builder"  (FROM golang:1.23 AS builder)           │
│  ──────────────────────────────────────────────────────────  │
│  Contains:                                                    │
│    • Full Go toolchain (compiler, stdlib source, go tool)     │
│    • Downloaded modules (go mod download)                     │
│    • Your source code (COPY . .)                               │
│    • Compiled binary (go build -o server .)                   │
│                                                                 │
│  Size: often 800MB - 1GB+                                     │
└───────────────────────────┬─────────────────────────────────┘
                             │
                    COPY --from=builder
                    /app/server  ->  /app/server
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│  STAGE 2: "runtime"  (FROM alpine:3.20  or  scratch  or        │
│                        gcr.io/distroless/static)               │
│  ──────────────────────────────────────────────────────────  │
│  Contains ONLY:                                                │
│    • The compiled binary                                       │
│    • Whatever minimal OS pieces the base provides               │
│      (alpine: shell + package manager; scratch: nothing;        │
│       distroless: CA certs + timezone data + non-root user)     │
│                                                                 │
│  Does NOT contain:                                              │
│    ✗ Go compiler          ✗ go.mod / go.sum                    │
│    ✗ Source code (.go)    ✗ Build cache / module cache         │
│                                                                 │
│  Size: single-digit to tens of MB                               │
└─────────────────────────────────────────────────────────────┘

This is the final image that gets pushed, pulled, and run in production.
```

---

## 📊 Image Size Comparison (Conceptual — this environment could not measure real sizes)

| Approach | What's in the final image | Typical size category |
|----------|---------------------------|------------------------|
| Single-stage (`FROM golang:1.23`, build and run in one stage) | Full Go toolchain + module cache + source code + binary | Very large (comparable to the base `golang` image itself, commonly several hundred MB to 1GB+) |
| Multi-stage → `alpine` runtime | Binary + shell + package manager + minimal Linux userland | Small (commonly single-digit to low tens of MB) |
| Multi-stage → `gcr.io/distroless/static` runtime | Binary + CA certs + timezone data + non-root user entry, no shell | Very small (commonly a few MB) |
| Multi-stage → `scratch` runtime | Binary only, nothing else | Smallest possible (the image size ≈ the binary's own size) |

**Why the categories differ this much:** a language runtime/toolchain (compiler, standard library sources, package manager) is fundamentally a different order of magnitude of content than a single statically-linked executable. This is true regardless of the specific project - it follows directly from Go producing one self-contained binary with no external runtime dependency.

---

## 🔁 Layer Caching Order Diagram

```
❌ INEFFICIENT ORDER                     ✅ EFFICIENT ORDER
─────────────────────                    ──────────────────
FROM golang:1.23 AS builder              FROM golang:1.23 AS builder
WORKDIR /app                             WORKDIR /app
COPY . .          ← invalidated by       COPY go.mod go.sum ./  ← invalidated
                     ANY file change        ONLY when dependencies change
RUN go mod download  ← re-downloads      RUN go mod download    ← cached across
                     on EVERY source        source-only edits!
                     change!
                                          COPY . .               ← invalidated by
RUN go build -o server .                    source changes (expected, cheap)

                                          RUN go build -o server .

Result: every code edit re-downloads     Result: every code edit reuses the
all dependencies from the network.       cached dependency layer; only the
                                          compile step re-runs.
```

```
Docker's caching rule:
  A layer is reused from cache ONLY IF:
    1. The instruction text is identical to the cached build, AND
    2. Every layer BEFORE it was also reused from cache, AND
    3. (for COPY/ADD) the copied files' contents are unchanged

  The moment one layer misses the cache, EVERY layer after it
  must be rebuilt - even if their own inputs didn't change.
```

---

## 📋 Dockerfile Instruction Reference Table

| Instruction | Purpose | Runs at | Notes |
|-------------|---------|---------|-------|
| `FROM` | Set the base image | Build | Can appear multiple times for multi-stage builds; `AS name` labels a stage |
| `WORKDIR` | Set/create the working directory | Build & runtime | Persists for all later instructions and the running container |
| `COPY` | Copy files from build context into the image | Build | Respects `.dockerignore`; `--from=<stage>` copies from an earlier stage instead |
| `ADD` | Like `COPY`, plus URL fetching and auto-extracting archives | Build | Prefer `COPY` unless you specifically need `ADD`'s extra behavior |
| `RUN` | Execute a command, bake its filesystem changes into a layer | Build | Each `RUN` is its own layer |
| `ENV` | Set an environment variable, baked into the image as a default | Build & runtime | Overridable at `docker run -e` without rebuilding |
| `ARG` | Define a build-time-only variable | Build only | Not present in the running container unless also assigned to an `ENV` |
| `EXPOSE` | Document which port(s) the container listens on | Metadata only | Does not publish the port - `docker run -p` does that |
| `USER` | Set the user the container process runs as | Runtime | Should be a non-root user in the final stage ([Best Practices](./README.md#best-practices)) |
| `HEALTHCHECK` | Define a command Docker runs periodically to check container health | Runtime | Needs a shell/utility present - not usable on `scratch` |
| `CMD` | Default command/arguments when the container starts | Runtime | Overridable entirely by `docker run <image> <other-command>` |
| `ENTRYPOINT` | Fixed command that always runs | Runtime | Arguments passed to `docker run` are appended, not substituted |

---

## 🚨 Common Mistakes

### Mistake 1: Skipping Multi-Stage Builds
```dockerfile
❌ FROM golang:1.23
   COPY . .
   RUN go build -o server .
   CMD ["./server"]        # ships the whole toolchain to production

✅ FROM golang:1.23 AS builder ... RUN go build ...
   FROM alpine:3.20
   COPY --from=builder /app/server .
```

### Mistake 2: Bloating the Build Context
```dockerfile
❌ COPY . .   (no .dockerignore — .git, .env, node_modules-equivalent clutter all sent)
✅ Write .dockerignore FIRST, before your first COPY . .
```

### Mistake 3: Hardcoding Secrets
```dockerfile
❌ ENV DATABASE_PASSWORD=hunter2     # baked into the image layer permanently
✅ docker run -e DATABASE_PASSWORD=hunter2 myimage   # injected at runtime
```

### Mistake 4: Not Pinning Image Tags
```dockerfile
❌ FROM golang:latest        # moving target, unreproducible builds
✅ FROM golang:1.23-alpine   # exact, reproducible
```

---

## 📈 Progression Summary

### Understanding Level 35

Level 35 teaches how to package a Go program the way it actually ships:

1. **Why containerize** — consistent environments, and Go's static binaries are a natural fit
2. **Dockerfile basics** — FROM/WORKDIR/COPY/RUN/EXPOSE/CMD/ENTRYPOINT
3. **Multi-stage builds** — separate compiling from running, ship only the binary
4. **.dockerignore** — keep the build context clean and secret-free
5. **Layer caching** — order instructions so dependency downloads are cached
6. **Environment variables** — the runtime bridge to Level 34's configuration patterns
7. **docker-compose** — multi-service applications, bridging Level 29's databases

### Prerequisites for Level 36

Before moving to Level 36 (Kubernetes), you need:

- ✅ Comfortable writing a multi-stage Dockerfile from scratch
- ✅ Understand why the runtime stage should be minimal, and which base to reach for
- ✅ Comfortable with `.dockerignore` and layer-caching order
- ✅ Comfortable running containers with `-p` and `-e`
- ✅ Understand how `ENV`/`docker run -e` connects to Level 34's configuration approach
- ✅ Can read and write a basic `docker-compose.yml`

### Ready for Level 36?

Level 36 teaches how to run containers like the ones you just built, at scale, across a cluster:
- Deployments, Pods, and Services
- ConfigMaps and Secrets (the orchestrator-level equivalent of `docker run -e`)
- Scaling and self-healing beyond what a single `docker run` can do

---

## ✅ Checklist Before Level 36

- [ ] Can write a multi-stage Dockerfile for a Go program from memory
- [ ] Can explain why the final image excludes the Go toolchain and source code
- [ ] Can write a `.dockerignore` covering `.git`, secrets, and editor files
- [ ] Can order Dockerfile instructions for layer-caching efficiency
- [ ] Can run a container with a mapped port and an injected environment variable
- [ ] Can write a `docker-compose.yml` combining an app and a database service
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The Fit
Go's statically-linked binaries need no interpreter or runtime in the image — a uniquely good match for minimal containers.

### The Split
Multi-stage builds separate "what it takes to build the program" from "what it takes to run it" — only the second survives into the shipped image.

### The Bridge
`ENV` and `docker run -e` are the runtime mechanism behind Level 34's "configure through the environment" principle — the same image behaves differently per environment with no rebuild.

### The Order
Layer caching rewards putting rarely-changing instructions (base image, dependency manifests) before frequently-changing ones (source code).

---

## 📚 Next Level

Level 36: Kubernetes
- Running containers at scale across a cluster
- Deployments, Services, and Pods
- ConfigMaps/Secrets building on this level's environment variables

You've learned to package a Go program the way it really ships! Keep going! 🚀
