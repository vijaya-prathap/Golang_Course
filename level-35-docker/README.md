# Level 35: Docker - Complete Guide

## Introduction

Welcome to Level 35! You've mastered `Level 34: Configuration` - reading settings from environment variables, flags, and config files instead of hardcoding them. Now it's time to learn **Docker**: packaging your Go program, together with everything it needs to run, into a single portable image that behaves identically on your laptop, a teammate's machine, and a production server.

> ⚠️ **This level is different from every Go-code level in this course.** Levels 0-31 (and most levels since) run entirely through `go run` / `go build` / `go test` - nothing but the Go toolchain. Docker is a **separate tool** from a separate vendor (Docker Inc.), with its own daemon that has to be installed and running. This level's Go source code is written and verified the normal way (`go build`, `go run`, real captured output), but the Docker commands themselves (`docker build`, `docker run`, `docker compose up`) could not be executed end-to-end in the environment that authored this course, because no Docker daemon was reachable here. See the banner at the top of **EXERCISES.md** for the full explanation, and treat every "Illustrative Output" block in this level as **correct per Docker's documented behavior, but not something you can trust byte-for-byte** the way you can trust the Go output in every other level. Run the commands yourself once you have Docker installed - your own output is the one that counts.

Why does this matter enough to be its own level? Because a Go program that works perfectly with `go run main.go` on your machine can still fail to run on a colleague's machine, or in production, for reasons that have nothing to do with your code: a missing environment variable, a different OS, a library version that isn't installed, a port that's already taken. Docker exists to eliminate that entire class of problem by shipping the *environment* along with the *program*.

---

## Table of Contents

1. [Why Containerize a Go App](#why-containerize-a-go-app)
2. [Dockerfile Basics](#dockerfile-basics)
3. [Multi-Stage Builds](#multi-stage-builds)
4. [.dockerignore](#dockerignore)
5. [Build Context and Layer Caching](#build-context-and-layer-caching)
6. [Building and Running](#building-and-running)
7. [Environment Variables in Containers](#environment-variables-in-containers)
8. [A docker-compose.yml Example](#a-docker-composeyml-example)
9. [Best Practices](#best-practices)
10. [Common Mistakes](#common-mistakes)

---

## Why Containerize a Go App

### The "Works on My Machine" Problem

You write a program, it runs fine locally. You hand it to someone else (or deploy it to a server) and it breaks - wrong OS, missing dependency, different environment variable, different version of some shared library. A **container** solves this by packaging the program together with a minimal, precisely-specified slice of an operating system (just enough filesystem, libraries, and configuration to run it) into a single artifact called an **image**. Anyone who runs that image - on their laptop, in CI, on a cloud server - gets the exact same environment every time.

This is different from a virtual machine. A VM virtualizes an entire operating system, including its own kernel, which is heavyweight (gigabytes, minutes to boot). A container shares the host machine's kernel and only isolates the process, filesystem, and network - so containers are lightweight (megabytes, milliseconds to start).

### Why Go Is Unusually Well-Suited to This

Most languages need their runtime shipped alongside your code: a Python container needs a Python interpreter, a Node container needs the Node runtime and often thousands of files in `node_modules`, a Java container needs a JVM. Go is different:

- **Go compiles to a single, statically-linked binary.** `go build` produces one executable file containing your code, the Go runtime (goroutine scheduler, garbage collector, everything) and (with `CGO_ENABLED=0`) no dependency on the host's C libraries either.
- **No interpreter, no VM, no runtime to install in the image.** The container just needs to be able to execute one file.
- **This means a Go container's runtime stage can be nearly empty** - as you'll see in the [Multi-Stage Builds](#multi-stage-builds) section, it can be just the binary and nothing else, not even an operating system in the traditional sense.

A minimal Python or Node image is tens to hundreds of megabytes before your code is even added, because the interpreter and its standard library have to be there. A minimal Go image can be the size of the binary itself, often single-digit to tens of megabytes total. This is one of the most concrete, practical payoffs of choosing a compiled, statically-linked language for backend services.

---

## Dockerfile Basics

A `Dockerfile` is a text file of instructions describing how to build an image, layer by layer. Here is a first, simple (deliberately **not yet optimized** - that comes in the next section) Dockerfile for a small Go HTTP server:

```dockerfile
FROM golang:1.23

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o server .

EXPOSE 8080

CMD ["./server"]
```

### Instruction by Instruction

**`FROM golang:1.23`** - every Dockerfile starts from a **base image**. `golang:1.23` is an official image that already has the Go toolchain (compiler, `go` command, standard library) installed on top of a Linux distribution. Everything you do afterward happens "inside" a container built from this starting point.

**`WORKDIR /app`** - sets the working directory for every instruction that follows, and creates it if it doesn't exist. Equivalent to `cd /app`, except it persists for the rest of the file (and becomes the working directory of the running container too).

**`COPY . .`** - copies files from the **build context** (the directory tree you point `docker build` at, typically your project's source) into the image filesystem. The first `.` means "everything in the build context"; the second `.` means "into the current `WORKDIR`."

**`RUN go mod download`** and **`RUN go build -o server .`** - `RUN` executes a shell command *while building the image*, and whatever it changes on disk becomes part of the resulting image layer. This is where you compile your program.

**`EXPOSE 8080`** - documents that the container listens on port 8080. This is metadata, not a firewall rule - it doesn't actually publish the port (that happens with `docker run -p`, covered later); it's there so anyone reading the Dockerfile (or tooling that inspects the image) knows which port matters.

**`CMD ["./server"]`** - the default command run when a container starts from this image. Written in **exec form** (a JSON array), which runs the command directly rather than through a shell - the recommended form, since it avoids an unnecessary shell process and handles OS signals (like `SIGTERM` on `docker stop`) correctly.

### CMD vs ENTRYPOINT

Both specify what runs when the container starts, but they compose differently:

- **`CMD`** provides a *default* that's easy to override: `docker run myimage ./other-binary` replaces the `CMD` entirely.
- **`ENTRYPOINT`** sets the command that *always* runs; anything passed on `docker run myimage <args>` is appended as arguments to it, not a replacement.

A common pattern combines them - `ENTRYPOINT ["./server"]` with `CMD ["--port=8080"]` - so the binary always runs, but its default flags can be overridden by anything passed after the image name on the command line. For a single-purpose container like this level's examples, a plain `CMD ["./server"]` is simplest and is what you'll use throughout.

This first Dockerfile works, but it ships the entire `golang` base image - the whole Go toolchain, your source code, and the compiled binary - to production. The next section fixes that.

---

## Multi-Stage Builds

A **multi-stage build** uses more than one `FROM` in the same Dockerfile. Earlier stages can do heavyweight work (like compiling); only the files you explicitly `COPY --from=<stage>` make it into the final image. Everything else - the compiler, the source tree, intermediate build artifacts - is discarded when the build finishes.

```dockerfile
# ---- Stage 1: build ----
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# ---- Stage 2: runtime ----
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
```

### Why This Shrinks the Image Dramatically

The **build stage** (`FROM golang:1.23 AS builder`) is large - it contains the full Go toolchain, all downloaded modules, your source files, and whatever intermediate object files the compiler produced. None of that is needed to *run* the finished program; it's only needed to *produce* it.

The **runtime stage** starts fresh from a small base image and copies in exactly one file: the compiled binary. Because `CGO_ENABLED=0` disables cgo, the resulting binary is statically linked and has no dependency on C libraries provided by the base OS, which is what makes it possible to run it on such a minimal base at all.

The final image contains **no Go toolchain, no source code, and none of the build stage's cache or layers** - only the binary and whatever minimal OS pieces the runtime base provides. This is the single biggest lever for reducing image size in a compiled-language Dockerfile, and it directly reflects the point made in [Why Containerize a Go App](#why-containerize-a-go-app): because the binary is self-contained, the runtime stage barely needs to contain anything else.

### Choosing a Runtime Base

| Base | Size | Contains | Trade-off |
|------|------|----------|-----------|
| `scratch` | 0 MB (empty) | Literally nothing - not even a shell, `ls`, or CA certificates | Smallest possible image; you must add CA certs yourself if your binary makes HTTPS calls, and you get no shell for debugging (`docker exec` won't have anything to run) |
| `gcr.io/distroless/static` | ~2 MB | CA certificates, timezone data, `/etc/passwd` entry for a non-root user - no shell, no package manager | Nearly as small as `scratch`, but solves the "my HTTPS client can't verify certificates" problem out of the box; still no shell for debugging |
| `alpine` | ~5-8 MB | A minimal but real Linux distribution: shell (`sh`), package manager (`apk`), basic utilities | Slightly larger, but you can `docker exec -it <container> sh` to poke around, and `apk add` anything else you need |

`scratch` and `distroless` are the leanest choices for a pure static binary with no debugging needs; `alpine` is the common middle ground when you want a shell available for troubleshooting. All three are dramatically smaller than shipping the full `golang` image, because none of them carry a compiler or source tree.

---

## .dockerignore

`docker build` sends the entire build context (everything in the directory you point it at) to the Docker daemon *before* any `COPY` instruction even runs - so files you never intend to `COPY` still slow down every build and could still end up copied in by an overly broad `COPY . .`. A `.dockerignore` file, using the same pattern syntax as `.gitignore`, excludes files from the build context entirely.

```dockerignore
# Version control
.git
.gitignore

# Local development files
*.md
.env
.env.local

# Editor/IDE files
.vscode/
.idea/
*.swp

# Build artifacts that shouldn't be baked in
*.log
tmp/
bin/

# Never ship compiled binaries you don't intend to run directly
server
```

Two concrete reasons this matters for a Go project specifically:

- **`.git`** can be enormous (full history) and has no reason to be inside a build context that only needs to compile the current source tree.
- **`.env`** (and similar local secret/config files) must never end up inside an image layer - see [Common Mistakes](#common-mistakes) for why baking secrets into an image is dangerous even if you later delete the file in a later instruction.

---

## Build Context and Layer Caching

Every instruction in a Dockerfile produces a **layer**, and Docker caches layers: if an instruction and everything before it are unchanged since the last build, Docker reuses the cached layer instead of re-running it. The moment one instruction's inputs change, that layer *and every layer after it* must be rebuilt.

This is why instruction **order** matters. Compare:

```dockerfile
# ❌ Less efficient: any source change invalidates the dependency download too
FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o server .
```

```dockerfile
# ✅ More efficient: dependency download is cached separately from source
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server .
```

In the second version, `COPY go.mod go.sum ./` only changes when your *dependencies* change (you added, removed, or upgraded a module). `RUN go mod download` therefore only re-runs - and re-downloads every dependency from the network - when `go.mod`/`go.sum` actually changed. Editing a `.go` source file (the overwhelmingly common case during development) only invalidates the cache starting at `COPY . .`, so the (often slow) dependency download step is skipped entirely and reused from cache.

```
Build 1 (first build):
  COPY go.mod go.sum ./   → cache MISS (nothing cached yet)
  RUN go mod download     → cache MISS → downloads everything
  COPY . .                → cache MISS
  RUN go build            → cache MISS → compiles

Build 2 (you edited main.go, dependencies unchanged):
  COPY go.mod go.sum ./   → cache HIT (files identical)
  RUN go mod download     → cache HIT (skipped entirely, no network call!)
  COPY . .                → cache MISS (source changed)
  RUN go build            → cache MISS → recompiles
```

The general rule: **order instructions from least-frequently-changing to most-frequently-changing.** Dependency manifests change rarely; source code changes constantly. Put the rare thing first.

---

## Building and Running

```bash
docker build -t go-server:1.0 .
```

- `-t go-server:1.0` tags the resulting image with a name and version (`docker build` alone produces an unnamed image identified only by its hash).
- `.` is the build context - the directory whose contents (minus anything in `.dockerignore`) are sent to the Docker daemon.

```bash
docker run -p 8080:8080 go-server:1.0
```

- `-p 8080:8080` maps port 8080 **on the host** to port 8080 **inside the container** (`-p <host>:<container>` - they don't have to match; `-p 9000:8080` would expose the container's port 8080 on the host's port 9000).

> ⚠️ **Illustrative output (unverified in this environment)** - Docker's daemon was not reachable here (see the banner in EXERCISES.md), so the following is what these commands are documented to produce, not a real captured build log:
>
> ```
> $ docker build -t go-server:1.0 .
> [+] Building 12.4s (14/14) FINISHED
>  => [internal] load build definition from Dockerfile           0.0s
>  => [internal] load .dockerignore                               0.0s
>  => [builder 1/5] FROM docker.io/library/golang:1.23            4.1s
>  => [internal] load build context                               0.1s
>  => [builder 2/5] WORKDIR /app                                  0.0s
>  => [builder 3/5] COPY go.mod go.sum ./                         0.0s
>  => [builder 4/5] RUN go mod download                           3.2s
>  => [builder 5/5] RUN CGO_ENABLED=0 GOOS=linux go build -o server .   4.8s
>  => [stage-1 1/3] FROM docker.io/library/alpine:3.20             0.2s
>  => [stage-1 2/3] WORKDIR /app                                  0.0s
>  => [stage-1 3/3] COPY --from=builder /app/server .              0.0s
>  => exporting to image                                          0.1s
>  => => naming to docker.io/library/go-server:1.0                0.0s
>
> $ docker images go-server
> REPOSITORY   TAG    IMAGE ID       CREATED         SIZE
> go-server    1.0    a1b2c3d4e5f6   5 seconds ago   14.2MB
>
> $ docker run -p 8080:8080 go-server:1.0
> 2026/01/01 00:00:00 listening on :8080
> ```
>
> The build log's exact step numbering, timings, and the `14.2MB` figure will differ on your machine - the shape of the output (build stage runs first, runtime stage runs second, final image contains only the binary) is what's documented and reliable; the specific bytes and seconds are not.

---

## Environment Variables in Containers

A container should never have secrets or environment-specific settings compiled into its binary or baked into its image - it should read them from the **environment**, exactly the pattern `Level 34: Configuration` establishes with `os.Getenv` and friends. Docker gives you two ways to set that environment:

**`ENV` in the Dockerfile** - bakes a *default* value into the image itself:

```dockerfile
ENV PORT=8080
```

**`-e` on `docker run`** - overrides it (or sets a variable with no Dockerfile default at all) at container-start time, without rebuilding the image:

```bash
docker run -p 9090:9090 -e PORT=9090 -e MESSAGE="Hello from a container" go-server:1.0
```

This is exactly why Level 34's configuration pattern (read settings from the environment rather than hardcoding them) matters so much for containers: the *same image*, built once, can run differently in development, staging, and production purely by changing which environment variables are injected at `docker run` time - no rebuild required. This is also precisely how orchestrators like Kubernetes (`Level 36: Kubernetes`) configure containers: a `Deployment` manifest sets environment variables on the container spec, and your Go program reads them the same way it would from a plain `docker run -e`.

> ⚠️ **Illustrative output (unverified in this environment):**
>
> ```
> $ docker run -p 9090:9090 -e PORT=9090 -e MESSAGE="Hello from a container" go-server:1.0
> 2026/01/01 00:00:00 listening on :9090
>
> $ curl http://localhost:9090/
> Hello from a container
> ```

---

## A docker-compose.yml Example

Real applications rarely run alone - they need a database, a cache, maybe other services. `docker compose` describes a multi-container application in one YAML file and starts/stops the whole group together. Here's the Go app from this level paired with a Postgres database, bridging back to `Level 29: Database/SQL`:

```yaml
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      PORT: "8080"
      DATABASE_URL: "postgres://appuser:apppass@db:5432/appdb?sslmode=disable"
    depends_on:
      - db

  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: appuser
      POSTGRES_PASSWORD: apppass
      POSTGRES_DB: appdb
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data

volumes:
  pgdata:
```

### What's Happening Here

- **`build: .`** tells Compose to build the `app` service's image from the Dockerfile in the current directory, rather than pulling a pre-built image.
- **`db: image: postgres:16-alpine`** pulls Postgres's official image directly - no Dockerfile needed for a service you're not writing yourself.
- **`environment:`** under `app` sets `DATABASE_URL` pointing at `db` - inside a Compose network, services reach each other **by service name** (`db`), not `localhost`; Compose provides DNS resolution between services automatically.
- **`depends_on: - db`** tells Compose to start `db` before `app` (it controls start *order*, not full "wait until Postgres is ready to accept queries" readiness - a real app should still retry its first connection, the same defensive pattern `Level 29: Database/SQL` uses around `sql.Open` being lazy).
- **`volumes: pgdata:/var/lib/postgresql/data`** persists Postgres's data directory outside the container's own filesystem, so `docker compose down` (without `-v`) doesn't wipe the database.

```bash
docker compose up --build
docker compose down
```

`up --build` builds the `app` image (if needed) and starts both services; `down` stops and removes the containers (add `-v` to also delete the named volume and its data).

---

## Best Practices

### 1. Always Use a Multi-Stage Build

```dockerfile
# ✅ Good - build stage discarded, only the binary ships
FROM golang:1.23 AS builder
# ... build ...
FROM alpine:3.20
COPY --from=builder /app/server .
```

```dockerfile
# ❌ Avoid - ships the entire toolchain and source tree to production
FROM golang:1.23
# ... build and run in the same stage ...
```

### 2. Never Run as Root in the Final Image

By default, a container runs as `root` unless told otherwise. If an attacker breaks out of your application, running as root hands them root inside the container (and a much easier path to escaping it). Create and switch to an unprivileged user in the runtime stage:

```dockerfile
FROM alpine:3.20
RUN adduser -D -u 10001 appuser
WORKDIR /app
COPY --from=builder /app/server .
USER appuser
CMD ["./server"]
```

### 3. Pin Base Image Versions - Never `latest`

```dockerfile
# ✅ Good - reproducible; you control when the base image changes
FROM golang:1.23-alpine

# ❌ Avoid - "latest" is a moving target that silently changes over time
FROM golang:latest
```

### 4. Keep the Runtime Image Minimal

Prefer `scratch` or `distroless` when your binary needs nothing else; use `alpine` only when you specifically need a shell or package manager for debugging. Every extra package is more attack surface and more bytes to pull and store.

### 5. Order Instructions for Cache Efficiency

Put the least-frequently-changing instructions first ([Build Context and Layer Caching](#build-context-and-layer-caching)): base image → dependency manifests → dependency download → source code → build.

---

## Common Mistakes

### Mistake 1: Skipping Multi-Stage Builds

```dockerfile
# ❌ WRONG - single stage ships the full golang image (900MB+) to production
FROM golang:1.23
WORKDIR /app
COPY . .
RUN go build -o server .
CMD ["./server"]
```

```dockerfile
# ✅ RIGHT - multi-stage ships only the compiled binary
FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o server .

FROM alpine:3.20
COPY --from=builder /app/server /server
CMD ["/server"]
```

### Mistake 2: Copying Unnecessary Files Into the Build Context

```dockerfile
# ❌ WRONG - no .dockerignore means .git, local .env files, and editor
# clutter are all sent to the Docker daemon and could end up in a COPY . .
COPY . .
```

Add a `.dockerignore` ([see above](#dockerignore)) before your first `COPY . .` - not after you notice the build is slow or the image is bloated.

### Mistake 3: Hardcoding Secrets Into the Image

```dockerfile
# ❌ WRONG - the secret is baked into an image layer permanently, even if
# a later instruction "deletes" it - previous layers are still inspectable
# with `docker history` and `docker save`
ENV DATABASE_PASSWORD=hunter2
```

```bash
# ✅ RIGHT - inject secrets at runtime, never at build time
docker run -e DATABASE_PASSWORD=hunter2 go-server:1.0
```

Anyone with access to the image can extract a value baked in at build time, even from an intermediate layer that a later instruction appears to overwrite. Secrets belong in runtime environment variables, mounted secret files, or a secrets manager - never in the Dockerfile or the image itself.

### Mistake 4: Not Pinning Image Tags

```dockerfile
# ❌ WRONG - "latest" (or no tag at all) means the same Dockerfile can
# produce a different image every time it's built, with no warning
FROM golang
FROM alpine:latest
```

```dockerfile
# ✅ RIGHT - the exact same Dockerfile always builds from the exact
# same base image until you deliberately change the tag
FROM golang:1.23-alpine
FROM alpine:3.20
```

An unpinned tag makes builds unreproducible: a build that worked yesterday can fail (or silently behave differently) today because `latest` moved out from under you, with no corresponding change to your own code or Dockerfile.

---

## Summary

**Why Containerize:**
- Eliminates "works on my machine" by shipping the exact runtime environment with the code
- Go's statically-linked binaries need no interpreter/runtime in the image, unlike Python, Node, or Java

**Dockerfile Basics:**
- `FROM` (base image), `WORKDIR` (working directory), `COPY` (files in), `RUN` (build-time commands), `EXPOSE` (documented port), `CMD`/`ENTRYPOINT` (what runs at container start)

**Multi-Stage Builds:**
- A `golang` build stage compiles; a minimal runtime stage (`scratch`, `distroless`, or `alpine`) ships only the binary
- Dramatically smaller images because the toolchain, source, and build cache never reach the final image

**.dockerignore:**
- Excludes `.git`, local dev files, and secrets from the build context before any `COPY` runs

**Layer Caching:**
- `COPY go.mod go.sum ./` + `RUN go mod download` **before** `COPY . .` lets dependency downloads be cached across builds when only source code changes

**Environment Variables:**
- `ENV` sets an image default; `docker run -e` overrides it at runtime - the same container image behaves differently per environment, connecting directly to `Level 34: Configuration`'s environment-variable-driven config pattern

**docker-compose.yml:**
- Describes multi-service applications (app + database) in one file; services reach each other by service name on the Compose network

**Best Practices:**
- Always multi-stage, never run as root, pin versions, keep images minimal, order instructions for cache efficiency

**Common Mistakes:**
- Skipping multi-stage builds, bloating the build context, hardcoding secrets, not pinning tags

---

## Next Steps

You now understand:
- ✅ Why Go's static binaries are unusually well-suited to small containers
- ✅ Core Dockerfile instructions and multi-stage builds
- ✅ .dockerignore and build-context hygiene
- ✅ Layer caching order for fast, repeatable builds
- ✅ Injecting environment variables into containers, connecting to Level 34's configuration patterns
- ✅ A basic docker-compose.yml combining an app and a database

**Next level:** Level 36 - Kubernetes
- Running containers at scale across a cluster
- Deployments, Services, and Pods
- How Kubernetes builds on everything you just learned about images and environment variables

You're one level away from orchestrating containers, not just building them! 🚀
