# Level 35: Docker - Exercises

> ⚠️ **This level's Docker commands could not be executed in the environment that authored this course**, because Docker is not available here (checked: `docker version` showed the client but `docker info`/`docker version`'s server section failed with "Cannot connect to the Docker daemon at unix:///Users/macbookpro/.docker/run/docker.sock. Is the docker daemon running?" - the daemon was not running/reachable). Every Dockerfile and command shown is written correctly per Docker's documented behavior, but - unlike every Go-code level in this course - none of it was actually run and verified end-to-end. When you run these commands yourself with Docker installed and running, **trust your own real output over any numbers or exact log text shown here.**
>
> The Go source code in these exercises **was** genuinely compiled and run with `go build`/`go run` in the authoring environment (that part has nothing to do with Docker), so the application logic itself is verified - only the containerization step is illustrative.

Complete all exercises in order. Each exercise builds on previous knowledge.

---

## Exercise 1: A Basic Single-Stage Dockerfile for a Go HTTP Server

**Objective:** Write and understand a first, unoptimized Dockerfile for a small Go program

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level35-exercise1
cd ~/projects/level35-exercise1
go mod init level35.example/exercise1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    message := os.Getenv("MESSAGE")
    if message == "" {
        message = "Hello from a Dockerized Go server!"
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, message)
    })

    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintln(w, "ok")
    })

    log.Printf("listening on :%s", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal(err)
    }
}
EOF
```

3. Verify the Go program itself works before containerizing anything:

```bash
go run main.go &
sleep 1
curl http://localhost:8080/
curl -i http://localhost:8080/health
kill %1
```

**Expected Output (real, verified with `go run` in this environment):**

```
2026/01/01 00:00:00 listening on :8080
Hello from a Dockerized Go server!
HTTP/1.1 200 OK
Content-Length: 3
Content-Type: text/plain; charset=utf-8

ok
```

4. Create a basic, single-stage `Dockerfile`:

```bash
cat > Dockerfile << 'EOF'
FROM golang:1.23

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o server .

EXPOSE 8080

CMD ["./server"]
EOF
```

5. Build and run the image:

```bash
docker build -t level35-exercise1 .
docker run -p 8080:8080 level35-exercise1
```

**Illustrative Output (unverified in this environment):**

```
$ docker build -t level35-exercise1 .
[+] Building 8.7s (10/10) FINISHED
 => [1/4] FROM docker.io/library/golang:1.23
 => [internal] load build context
 => [2/4] WORKDIR /app
 => [3/4] COPY . .
 => [4/4] RUN go build -o server .
 => exporting to image
 => => naming to docker.io/library/level35-exercise1

$ docker run -p 8080:8080 level35-exercise1
2026/01/01 00:00:00 listening on :8080
```

**Learning Objectives:**
- ✅ Write `main.go` for a minimal, environment-configurable HTTP server
- ✅ Write a first, unoptimized Dockerfile using FROM, WORKDIR, COPY, RUN, EXPOSE, CMD
- ✅ Understand that this Dockerfile ships the entire `golang` toolchain image to production

---

## Exercise 2: Converting to a Multi-Stage Build

**Objective:** Shrink the image by separating the build stage from the runtime stage

**Instructions:**

1. Reuse the working directory from Exercise 1, or create a fresh one:

```bash
mkdir -p ~/projects/level35-exercise2
cd ~/projects/level35-exercise2
go mod init level35.example/exercise2
```

2. Copy the same `main.go` from Exercise 1 into this directory.

3. Replace `Dockerfile` with a multi-stage version:

```bash
cat > Dockerfile << 'EOF'
# ---- Stage 1: build ----
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# ---- Stage 2: runtime ----
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
EOF
```

4. Build both versions and compare their sizes:

```bash
docker build -t level35-single -f- . <<'SINGLE'
FROM golang:1.23
WORKDIR /app
COPY . .
RUN go build -o server .
CMD ["./server"]
SINGLE

docker build -t level35-multistage .

docker images | grep level35
```

**Illustrative Output (unverified in this environment):**

```
$ docker images | grep level35
level35-multistage   latest   ...   14.1MB
level35-single       latest   ...   912MB
```

**Conceptual size comparison (why, not just "trust these numbers"):** the single-stage image contains the full `golang:1.23` base (Go compiler, standard library source, module cache) *plus* your source code *plus* the compiled binary. The multi-stage image's runtime stage starts fresh from `alpine` and copies in exactly one file - the already-compiled binary. The compiler, your `.go` source files, and any intermediate build cache from the `builder` stage never exist in the final image at all. The real ratio you see on your machine will differ from any specific number above, but it will be dramatic - commonly 10-50x smaller - because the categories of content being removed (a whole language toolchain vs. one binary) differ by that much in kind, not just in this configuration.

**Learning Objectives:**
- ✅ Convert a single-stage Dockerfile into a multi-stage build
- ✅ Understand `AS builder` naming and `COPY --from=builder`
- ✅ Reason about *why* the size difference is dramatic, not just observe that it is

---

## Exercise 3: Writing a .dockerignore File

**Objective:** Exclude files that shouldn't be sent to the Docker build context

**Instructions:**

1. In the Exercise 2 directory, add a few files that should never end up in an image:

```bash
mkdir -p .git
echo "fake git internals" > .git/HEAD
echo "SECRET_KEY=do-not-ship-this" > .env
echo "notes to self" > NOTES.md
```

2. Create `.dockerignore`:

```bash
cat > .dockerignore << 'EOF'
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

# Compiled output - rebuilt inside the image, not copied in
server
EOF
```

3. Confirm what a build context would include by listing what's *not* ignored:

```bash
# There is no built-in "dry run" flag for build context contents, but you
# can approximate it by listing files git would track, or manually reviewing
# against the .dockerignore patterns:
ls -la
cat .dockerignore
```

**Illustrative Output (unverified in this environment) - showing that `.git`, `.env`, and `NOTES.md` are excluded from what `COPY . .` would copy:**

```
$ docker build -t level35-exercise3 .
[+] Building 6.2s (9/9) FINISHED
 => [internal] load .dockerignore
 => => transferring context: 3 patterns
 ...
```

**Learning Objectives:**
- ✅ Write a `.dockerignore` using gitignore-style patterns
- ✅ Understand that `.dockerignore` applies to the entire build context, before any `COPY` runs
- ✅ Recognize why `.git` and `.env` files must never reach the build context

---

## Exercise 4: Layer Caching With COPY go.mod Before COPY . .

**Objective:** Demonstrate the dependency-caching pattern that speeds up repeated builds

**Instructions:**

1. Create a fresh working directory with a real dependency (so `go.sum` exists):

```bash
mkdir -p ~/projects/level35-exercise4
cd ~/projects/level35-exercise4
go mod init level35.example/exercise4
```

2. Create `main.go` using only the standard library (kept dependency-free here so the exercise runs offline; the pattern below applies identically whether or not you have third-party dependencies):

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("Layer caching demo")
}
EOF
go build -o app . && ./app
```

**Expected Output (real, verified with `go build`/`go run` in this environment):**

```
Layer caching demo
```

3. Write a Dockerfile using the cache-efficient ordering:

```bash
cat > Dockerfile << 'EOF'
FROM golang:1.23 AS builder
WORKDIR /app

# Copy dependency manifests FIRST - this layer only invalidates when
# dependencies change, not when source code changes
COPY go.mod ./
RUN go mod download

# Copy source code AFTER - this is what changes on every edit
COPY . .
RUN CGO_ENABLED=0 go build -o app .

FROM alpine:3.20
COPY --from=builder /app/app /app
CMD ["/app"]
EOF
```

4. Simulate two builds: an initial build, then a source-only change:

```bash
docker build -t level35-exercise4 .

# Now change ONLY the source code, not go.mod:
echo '// a harmless comment' >> main.go

docker build -t level35-exercise4 .
```

**Illustrative Output (unverified in this environment) - the second build's dependency-download step is cached:**

```
$ docker build -t level35-exercise4 .          # first build
 => [builder 2/5] COPY go.mod ./                CACHED
 => [builder 3/5] RUN go mod download           CACHED
 => [builder 4/5] COPY . .                      2.1s   (source changed)
 => [builder 5/5] RUN go build -o app .         1.8s
```

**Learning Objectives:**
- ✅ Order `COPY go.mod ./` + `RUN go mod download` before `COPY . .`
- ✅ Understand that Docker invalidates a layer (and everything after it) only when that layer's own inputs change
- ✅ Recognize why this ordering matters most on a real dependency-download step, which can be slow over the network

---

## Exercise 5: Running the Container With a Mapped Port and an Injected Environment Variable

**Objective:** Run a container with `-p` port mapping and `-e` environment injection

**Instructions:**

1. Reuse the `main.go` and Dockerfile from Exercise 1 (the `PORT`/`MESSAGE`-aware server).

2. Build the image:

```bash
docker build -t level35-exercise5 .
```

3. Run it with default settings, then again with overrides:

```bash
# Default: reads its own fallback values (PORT=8080, default message)
docker run -d --name level35-default -p 8080:8080 level35-exercise5
curl http://localhost:8080/
docker stop level35-default && docker rm level35-default

# Overridden: host port 9090 -> container port 9090, custom message
docker run -d --name level35-custom -p 9090:9090 \
    -e PORT=9090 \
    -e MESSAGE="Hello from an injected environment variable!" \
    level35-exercise5
curl http://localhost:9090/
docker stop level35-custom && docker rm level35-custom
```

**Illustrative Output (unverified in this environment) - shape matches the real `go run` output captured in Exercise 1, since the container runs the identical binary:**

```
$ curl http://localhost:8080/
Hello from a Dockerized Go server!

$ curl http://localhost:9090/
Hello from an injected environment variable!
```

**Learning Objectives:**
- ✅ Map a host port to a container port with `docker run -p`
- ✅ Inject environment variables at runtime with `docker run -e`, overriding the image's built-in defaults
- ✅ Recognize this as the same environment-driven configuration pattern from Level 34, applied to a running container

---

## Exercise 6: Adding an ENV Default in the Dockerfile

**Objective:** Set a default environment variable in the image itself, then override it at runtime

**Instructions:**

1. In the Exercise 5 project, update the Dockerfile's runtime stage to declare a default:

```bash
cat > Dockerfile << 'EOF'
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/server .
ENV PORT=8080
ENV MESSAGE="Default message baked into the image"
EXPOSE 8080
CMD ["./server"]
EOF
```

2. Run it three ways: as-is, with one override, and with both overridden:

```bash
docker build -t level35-exercise6 .

docker run --rm -p 8080:8080 level35-exercise6 &
sleep 1 && curl http://localhost:8080/ && kill %1

docker run --rm -p 8080:8080 -e MESSAGE="Overridden at runtime" level35-exercise6 &
sleep 1 && curl http://localhost:8080/ && kill %1
```

**Illustrative Output (unverified in this environment):**

```
$ curl http://localhost:8080/       # using the image's ENV default
Default message baked into the image

$ curl http://localhost:8080/       # after -e MESSAGE=... override
Overridden at runtime
```

**Learning Objectives:**
- ✅ Set a default value with `ENV` in the Dockerfile
- ✅ Confirm that `docker run -e` overrides an `ENV` default without rebuilding the image
- ✅ Understand why baking a *default* (not a secret) into the image with `ENV` is safe, unlike hardcoding a real credential

---

## Exercise 7: A .dockerignore Audit

**Objective:** Practice reviewing a project for files that should never enter a build context

**Instructions:**

1. Create a project directory simulating a realistic (slightly messy) Go project:

```bash
mkdir -p ~/projects/level35-exercise7/.git
mkdir -p ~/projects/level35-exercise7/.vscode
cd ~/projects/level35-exercise7
go mod init level35.example/exercise7
echo 'package main; func main() {}' > main.go
echo "DB_PASSWORD=supersecret" > .env
echo "some notes" > TODO.md
echo '{"editor.tabSize": 4}' > .vscode/settings.json
touch .git/HEAD
mkdir -p bin && echo "old binary" > bin/app-old
```

2. Write a `.dockerignore` that excludes everything that shouldn't reach the build context, but keeps `go.mod`, `go.sum` (if present), and `*.go` files:

```bash
cat > .dockerignore << 'EOF'
.git
.gitignore
.vscode/
.idea/
*.md
.env
.env.*
bin/
*.log
EOF
```

3. Verify by listing the directory and cross-checking each entry against your `.dockerignore` patterns:

```bash
find . -not -path './.git*' | sort
cat .dockerignore
```

**Expected reasoning (this is an audit exercise, not a runtime check):** every file EXCEPT `main.go`, `go.mod`, and the `.dockerignore` file itself should match one of your exclusion patterns. If any secret-bearing or oversized file (`.env`, `.git`, `bin/`) is not covered by a pattern, the audit has found a real gap - fix the `.dockerignore` before it matters.

**Learning Objectives:**
- ✅ Recognize which project files are safe to enter a Docker build context and which are not
- ✅ Write `.dockerignore` patterns covering version control, secrets, editor files, and stale build artifacts
- ✅ Practice auditing a project structure rather than only copying a template

---

## Exercise 8: A docker-compose.yml With an App and a Database

**Objective:** Combine the Go app with a Postgres database service using Compose

**Instructions:**

1. In a fresh directory, create `main.go` that reads a `DATABASE_URL` the way `Level 29: Database/SQL` and `Level 34: Configuration` would expect (kept simple - it only reads and prints the value, since actually dialing Postgres would require a real running database and the `postgres` driver as a dependency):

```bash
mkdir -p ~/projects/level35-exercise8
cd ~/projects/level35-exercise8
go mod init level35.example/exercise8
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        dbURL = "(not set)"
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "app is up; DATABASE_URL=%s\n", dbURL)
    })

    log.Printf("listening on :%s (DATABASE_URL=%s)", port, dbURL)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}
EOF
go build -o app . && DATABASE_URL="postgres://appuser:apppass@db:5432/appdb" PORT=8080 ./app &
sleep 1
curl http://localhost:8080/
kill %1
```

**Expected Output (real, verified with `go build`/`go run` in this environment - note this confirms the Go program reads `DATABASE_URL` correctly; it does not confirm connectivity to an actual database):**

```
app is up; DATABASE_URL=postgres://appuser:apppass@db:5432/appdb
```

2. Write the Dockerfile (multi-stage, matching Exercise 2's pattern):

```bash
cat > Dockerfile << 'EOF'
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["./app"]
EOF
```

3. Write `docker-compose.yml`:

```bash
cat > docker-compose.yml << 'EOF'
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
EOF
```

4. Bring the stack up:

```bash
docker compose up --build
```

**Illustrative Output (unverified in this environment):**

```
$ docker compose up --build
[+] Running 2/2
 ✔ Container level35-exercise8-db-1   Started
 ✔ Container level35-exercise8-app-1  Started
app-1  | 2026/01/01 00:00:00 listening on :8080 (DATABASE_URL=postgres://appuser:apppass@db:5432/appdb?sslmode=disable)

$ curl http://localhost:8080/
app is up; DATABASE_URL=postgres://appuser:apppass@db:5432/appdb?sslmode=disable
```

**Learning Objectives:**
- ✅ Write a `docker-compose.yml` combining an application service and a database service
- ✅ Understand service-name-based DNS resolution between Compose services (`db`, not `localhost`)
- ✅ Connect environment-variable configuration (Level 34) to a multi-service Compose stack (bridging Level 29)

---

## Exercise 9: HEALTHCHECK and Graceful Container Lifecycle

**Objective:** Add a `HEALTHCHECK` instruction so Docker can detect whether the app is actually serving requests

**Instructions:**

1. Reuse the Exercise 1 server (it already exposes `GET /health`).

2. Add a `HEALTHCHECK` to the Dockerfile's runtime stage:

```bash
cat > Dockerfile << 'EOF'
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
HEALTHCHECK --interval=10s --timeout=2s --start-period=5s --retries=3 \
    CMD wget -q -O- http://localhost:8080/health || exit 1
CMD ["./server"]
EOF
```

3. Build and observe the reported health state:

```bash
docker build -t level35-exercise9 .
docker run -d --name level35-health -p 8080:8080 level35-exercise9
sleep 12
docker ps --filter name=level35-health --format "table {{.Names}}\t{{.Status}}"
docker stop level35-health && docker rm level35-health
```

**Illustrative Output (unverified in this environment):**

```
$ docker ps --filter name=level35-health --format "table {{.Names}}\t{{.Status}}"
NAMES              STATUS
level35-health     Up 15 seconds (healthy)
```

Note: `alpine`'s minimal image does include `wget` by default, which is why it's used here rather than `curl` (not installed by default on `alpine`). On `scratch` or `distroless`, there is no shell or utility to run at all, so `HEALTHCHECK` isn't usable in the same way - an orchestrator's own liveness/readiness probes (as in `Level 36: Kubernetes`) become the equivalent mechanism there.

**Learning Objectives:**
- ✅ Add a `HEALTHCHECK` instruction with interval, timeout, and retry configuration
- ✅ Understand why `HEALTHCHECK` needs a shell/utility present, which constrains which runtime base images support it directly
- ✅ Connect container health checks conceptually to Kubernetes liveness/readiness probes (Level 36)

---

## Exercise 10: Comprehensive Practice — Fully Containerizing a REST API

**Objective:** Combine everything into a production-quality multi-stage Dockerfile for a small REST API, echoing Level 27's HTTP & REST APIs patterns

**Instructions:**

1. Create the project:

```bash
mkdir -p ~/projects/level35-exercise10
cd ~/projects/level35-exercise10
go mod init level35.example/exercise10
```

2. Create `main.go` - an in-memory Task API in the style of Level 27's task manager exercise:

```bash
cat > main.go << 'EOF'
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "sync"
)

type Task struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Done bool   `json:"done"`
}

type store struct {
    mu     sync.Mutex
    nextID int
    tasks  map[int]Task
}

func newStore() *store {
    return &store{nextID: 1, tasks: make(map[int]Task)}
}

func (s *store) create(name string) Task {
    s.mu.Lock()
    defer s.mu.Unlock()
    t := Task{ID: s.nextID, Name: name}
    s.tasks[t.ID] = t
    s.nextID++
    return t
}

func (s *store) list() []Task {
    s.mu.Lock()
    defer s.mu.Unlock()
    out := make([]Task, 0, len(s.tasks))
    for i := 1; i < s.nextID; i++ {
        if t, ok := s.tasks[i]; ok {
            out = append(out, t)
        }
    }
    return out
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    s := newStore()
    s.create("Write Dockerfile")
    s.create("Build image")

    mux := http.NewServeMux()
    mux.HandleFunc("GET /tasks", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(s.list())
    })
    mux.HandleFunc("POST /tasks", func(w http.ResponseWriter, r *http.Request) {
        var body struct {
            Name string `json:"name"`
        }
        if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
            http.Error(w, "invalid body", http.StatusBadRequest)
            return
        }
        t := s.create(body.Name)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(t)
    })
    mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "ok")
    })

    log.Printf("task API listening on :%s", port)
    log.Fatal(http.ListenAndServe(":"+port, mux))
}
EOF
```

3. Verify the API works before containerizing it:

```bash
go build -o taskapi . && PORT=8081 ./taskapi &
sleep 1
curl http://localhost:8081/tasks
curl -X POST -d '{"name":"Push image to registry"}' http://localhost:8081/tasks
curl http://localhost:8081/tasks
kill %1
```

**Expected Output (real, verified with `go build`/`go run` in this environment):**

```
[{"id":1,"name":"Write Dockerfile","done":false},{"id":2,"name":"Build image","done":false}]
{"id":3,"name":"Push image to registry","done":false}
[{"id":1,"name":"Write Dockerfile","done":false},{"id":2,"name":"Build image","done":false},{"id":3,"name":"Push image to registry","done":false}]
```

4. Write a `.dockerignore`:

```bash
cat > .dockerignore << 'EOF'
.git
.gitignore
*.md
.env
.vscode/
.idea/
taskapi
EOF
```

5. Write a production-quality multi-stage Dockerfile: cache-friendly layer order, `CGO_ENABLED=0`, a minimal runtime base, a non-root `USER`, and a `HEALTHCHECK`:

```bash
cat > Dockerfile << 'EOF'
# ---- Stage 1: build ----
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o taskapi .

# ---- Stage 2: runtime ----
FROM alpine:3.20

RUN adduser -D -u 10001 appuser
WORKDIR /app

COPY --from=builder /app/taskapi .

ENV PORT=8080
EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=2s --start-period=5s --retries=3 \
    CMD wget -q -O- http://localhost:8080/health || exit 1

USER appuser

CMD ["./taskapi"]
EOF
```

6. Build and run:

```bash
docker build -t level35-taskapi:1.0 .
docker run -d --name level35-taskapi -p 8080:8080 level35-taskapi:1.0
curl http://localhost:8080/tasks
docker stop level35-taskapi && docker rm level35-taskapi
```

**Illustrative Output (unverified in this environment):**

```
$ docker build -t level35-taskapi:1.0 .
[+] Building 11.3s (15/15) FINISHED
 ...
 => => naming to docker.io/library/level35-taskapi:1.0

$ docker images level35-taskapi
REPOSITORY        TAG    IMAGE ID       CREATED         SIZE
level35-taskapi   1.0    f7e8d9c0b1a2   3 seconds ago   13.8MB

$ curl http://localhost:8080/tasks
[{"id":1,"name":"Write Dockerfile","done":false},{"id":2,"name":"Build image","done":false}]
```

**Learning Objectives:**
- ✅ Containerize a full REST API end-to-end, from Go source through a production-quality Dockerfile
- ✅ Combine multi-stage builds, layer caching order, `.dockerignore`, `HEALTHCHECK`, and a non-root `USER` in one Dockerfile
- ✅ Distinguish what was genuinely verified (the Go API's behavior) from what is illustrative (the container build/run output)

---

## Bonus Challenges

### Challenge 1: A Non-Root USER Setup From Scratch

Take any Dockerfile from this level and, without looking at Exercise 10's answer, add a dedicated non-root user to the runtime stage and confirm (via `docker exec <container> whoami` if you're on `alpine`, or by reading the Dockerfile alone if you're on `scratch`/`distroless`) that the container does not run as root.

**Hints:**
- On Alpine: `RUN adduser -D -u 10001 appuser` then `USER appuser` near the end of the Dockerfile, after any `COPY`/`RUN` steps that need root (like installing packages)
- On `distroless/static`, there's already a built-in non-root `nonroot` user/UID you can select with `USER nonroot` - no `adduser` needed
- Order matters: `USER` should come after everything that needs root privileges (file copies owned by root are still readable; the point is the *process* doesn't run as root)

### Challenge 2: Multi-Architecture Build Awareness

Go can cross-compile for different OS/architecture combinations using `GOOS`/`GOARCH` build environment variables, and Docker's `buildx` can build multi-architecture images (e.g., for both `amd64` and `arm64`) from one Dockerfile.

**Hints:**
- Research `docker buildx build --platform linux/amd64,linux/arm64 -t myimage .`
- In the build stage, `ARG TARGETOS TARGETARCH` combined with `RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build ...` lets one Dockerfile produce the right binary for each requested platform
- This matters because a Go binary built for `arm64` won't run in an `amd64` container and vice versa - unlike an interpreted language, there is no "just works everywhere" fallback at the binary level

### Challenge 3: A Deliberately Broken .dockerignore Audit

Write a `.dockerignore` with an intentional gap (e.g., forget to exclude `.env`), then reason through - without running anything - exactly what would go wrong if this image were built and pushed to a shared registry.

**Hints:**
- Trace the path: build context → `COPY . .` → image layer → pushed image → anyone who can `docker pull` it
- Consider `docker history <image>` as a real tool for inspecting what previous layers contain
- Cross-reference with [Common Mistakes](./README.md#common-mistakes) Mistake 3 (hardcoding secrets) - a leaked `.env` inside an image is functionally the same failure

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Write a basic single-stage Dockerfile and explain why it's inefficient
✅ Convert it into a multi-stage build and explain why the image shrinks dramatically
✅ Write a `.dockerignore` that keeps secrets and clutter out of the build context
✅ Order Dockerfile instructions for layer-caching efficiency
✅ Run containers with mapped ports and injected environment variables
✅ Set `ENV` defaults and override them at runtime, connecting to Level 34's configuration patterns
✅ Write a `docker-compose.yml` combining an app and a database service
✅ Add a `HEALTHCHECK` and a non-root `USER` to a production-quality Dockerfile
✅ Fully containerize a REST API end-to-end

---

## Next Level

Level 36: Kubernetes
- Running containers at scale across a cluster
- Deployments, Services, and Pods
- Configuring containers with ConfigMaps and Secrets, building on this level's environment variables

Great work! You've learned to package Go programs the way they actually ship in the real world! 🚀
