# Level 35: Quick Reference Card

> ⚠️ **Docker was not available in the environment that authored this course** (daemon unreachable). Go source shown was verified with `go build`/`go run`; Docker command output in this level is illustrative, per documented behavior, not captured from a real run. Trust your own output once you run these commands.

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir myapp && cd myapp
go mod init myapp

# main.go
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "net/http"
    "os"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello from Docker!")
    })
    http.ListenAndServe(":"+port, nil)
}
EOF

# Multi-stage Dockerfile
cat > Dockerfile << 'EOF'
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server .

FROM alpine:3.20
COPY --from=builder /app/server /server
EXPOSE 8080
CMD ["/server"]
EOF

# Build and run
docker build -t myapp .
docker run -p 8080:8080 myapp
```

---

## 📋 Core Dockerfile Instructions

```dockerfile
FROM golang:1.23 AS builder   # base image, optional stage name
WORKDIR /app                  # set/create working directory
COPY go.mod go.sum ./         # copy files from build context
RUN go mod download           # execute a command at build time
ENV PORT=8080                 # set a default env var (overridable at runtime)
EXPOSE 8080                   # document the listening port (metadata only)
USER appuser                  # run as a non-root user
HEALTHCHECK CMD wget ...      # periodic container health check
CMD ["./server"]              # default command (overridable)
ENTRYPOINT ["./server"]       # fixed command (args appended, not replaced)
```

---

## 🏗️ Multi-Stage Build Skeleton

```dockerfile
# ---- build ----
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# ---- runtime ----
FROM alpine:3.20
RUN adduser -D -u 10001 appuser
WORKDIR /app
COPY --from=builder /app/server .
USER appuser
EXPOSE 8080
CMD ["./server"]
```

Runtime base choices: `scratch` (nothing, smallest) · `gcr.io/distroless/static` (+ CA certs, non-root user, no shell) · `alpine` (+ shell, package manager).

---

## 🚫 .dockerignore Essentials

```dockerignore
.git
.gitignore
*.md
.env
.env.*
.vscode/
.idea/
bin/
```

---

## 🔁 Layer Caching Order

```dockerfile
# ✅ Dependency manifests BEFORE source code
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server .
```

Rule: least-frequently-changing instructions first, most-frequently-changing last.

---

## 🐳 Build & Run Commands

```bash
docker build -t myapp:1.0 .              # build and tag an image
docker images                             # list images (and their sizes)
docker run -p 8080:8080 myapp:1.0         # map host:container port
docker run -e KEY=VALUE myapp:1.0         # inject an env var at runtime
docker run -d --name myapp myapp:1.0      # run detached, name the container
docker ps                                 # list running containers
docker logs myapp                         # view a container's stdout/stderr
docker stop myapp && docker rm myapp      # stop and remove
docker rmi myapp:1.0                      # remove the image
```

---

## 🌱 Environment Variables

```dockerfile
ENV PORT=8080                             # image default
```

```bash
docker run -e PORT=9090 myapp:1.0         # runtime override, no rebuild
```

Same pattern Level 34's configuration code already reads via `os.Getenv` — Docker (and later, Kubernetes in Level 36) is just what sets those variables in production.

---

## 🧩 docker-compose.yml Skeleton

```yaml
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: "postgres://user:pass@db:5432/appdb?sslmode=disable"
    depends_on:
      - db

  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: pass
      POSTGRES_DB: appdb
    volumes:
      - pgdata:/var/lib/postgresql/data

volumes:
  pgdata:
```

```bash
docker compose up --build     # build and start all services
docker compose down           # stop and remove containers
docker compose down -v        # also delete named volumes
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Skipping multi-stage builds | Single `FROM golang` stage ships the whole toolchain | Multi-stage, copy only the compiled binary |
| Bloated build context | `COPY . .` with no `.dockerignore` | Write `.dockerignore` before the first `COPY` |
| Hardcoded secrets | `ENV DB_PASSWORD=hunter2` | `docker run -e DB_PASSWORD=...` at runtime |
| Unpinned tags | `FROM golang:latest` | `FROM golang:1.23-alpine` |
| Running as root | No `USER` instruction | `USER appuser` in the final stage |

---

## 🎓 Before Next Level

Can you:
- [ ] Write a multi-stage Dockerfile from memory?
- [ ] Explain why the runtime stage doesn't need the Go toolchain?
- [ ] Order instructions for layer-caching efficiency?
- [ ] Run a container with a mapped port and an injected env var?
- [ ] Write a `docker-compose.yml` with an app and a database?

If YES → You're ready for Level 36!

---

## 📚 Next Level

Level 36: Kubernetes
- Deployments, Pods, and Services
- Running containers at scale across a cluster
- ConfigMaps/Secrets building on this level's environment variables

You've got containerization down! 💪
