# Level 35: Docker - INDEX

Welcome to **Level 35: Docker**! This is where your Go programs stop being "something that runs on my machine" and become portable images that run identically anywhere.

> ⚠️ **Note:** Docker's daemon was not reachable in the environment that authored this course (checked with `docker version` / `docker info`). The Go source code in this level was genuinely compiled and run; the Docker build/run output shown throughout is illustrative and correct per documented behavior, but not captured from a real execution. See EXERCISES.md's banner for full details.

---

## 📖 What You'll Learn

- ✅ Why containerize a Go app, and why Go's static binaries are unusually well-suited to it
- ✅ Dockerfile basics: FROM, WORKDIR, COPY, RUN, EXPOSE, CMD/ENTRYPOINT
- ✅ Multi-stage builds: compiling in one stage, shipping only the binary in another
- ✅ .dockerignore for a clean, secret-free build context
- ✅ Layer caching and instruction ordering for fast repeated builds
- ✅ Building and running containers with `docker build`/`docker run`
- ✅ Environment variables in containers, connecting to Level 34's configuration patterns
- ✅ docker-compose.yml for multi-service apps (app + database)

---

## 🗂️ Level 35 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview and welcome message. Start here first!

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- Why containerize a Go app
- Dockerfile basics and a first simple Dockerfile
- Multi-stage builds and choosing a runtime base image
- .dockerignore, build context, and layer caching
- Building, running, and injecting environment variables
- A docker-compose.yml example
- Best practices
- Common mistakes

**Read Time:** 40-55 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 detailed, step-by-step exercises:
1. A basic single-stage Dockerfile
2. Converting to a multi-stage build
3. Writing a .dockerignore file
4. Layer caching with COPY go.mod before COPY . .
5. Running with a mapped port and injected environment variable
6. Adding an ENV default in the Dockerfile
7. A .dockerignore audit
8. A docker-compose.yml with an app and a database
9. HEALTHCHECK and graceful container lifecycle
10. Comprehensive practice - fully containerizing a REST API

**Time Commitment:** 4-5 hours
**When:** After reading README

### 4. **STUDY_GUIDE.md** (Visual Learning)
- Multi-stage build diagram (build stage → copy binary → runtime stage)
- Conceptual image-size comparison table
- Layer-caching order diagram
- Dockerfile instruction reference table
- Common mistakes guide

**Read Time:** 30-40 minutes
**When:** Use while doing exercises

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- Essential Dockerfile syntax
- Build/run command reference
- docker-compose.yml skeleton
- Common mistakes table

**Read Time:** 10-15 minutes
**When:** Quick lookup during practice

---

## 🎯 Recommended Learning Path

### Day 1: Dockerfile Basics (2 hours)
1. Read **START_HERE.txt** (10 min)
2. Read **README.md** sections 1-2 (25 min)
3. Complete Exercise 1 (1 hour 25 min)

### Day 2: Multi-Stage Builds & .dockerignore (2 hours)
1. Read **README.md** sections 3-4 (25 min)
2. Complete Exercises 2-3 (1.5 hours)

### Day 3: Caching, Running & Environment Variables (2 hours)
1. Read **README.md** sections 5-7 (30 min)
2. Complete Exercises 4-6 (1.5 hours)

### Day 4: Compose & Real-World Practice (2 hours)
1. Read **README.md** sections 8-10 (25 min)
2. Complete Exercises 7-10 (1.5 hours)

### Day 5: Consolidation (1 hour)
1. Try bonus challenges
2. Review with **QUICK_REFERENCE.md**
3. Make sure the multi-stage build pattern feels automatic

---

## 💡 Key Concepts At A Glance

### A Minimal Multi-Stage Dockerfile
```dockerfile
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server .

FROM alpine:3.20
COPY --from=builder /app/server /server
CMD ["/server"]
```

### Build and Run
```bash
docker build -t myapp:1.0 .
docker run -p 8080:8080 -e PORT=8080 myapp:1.0
```

### docker-compose.yml
```yaml
services:
  app:
    build: .
    ports: ["8080:8080"]
    depends_on: [db]
  db:
    image: postgres:16-alpine
```

---

## ✅ Prerequisites

Make sure you've completed **Level 34: Configuration**

You need:
- ✅ Comfort reading settings from environment variables (`os.Getenv`) and understanding why hardcoded config is a problem
- ✅ Basic familiarity with building and running Go binaries (`go build`)
- ✅ Comfort with HTTP servers from `Level 27: HTTP & REST APIs` (used in the exercises' example programs)

You also need **Docker Desktop (or an equivalent Docker Engine) installed and running** on your own machine - this is one of the few levels in this course that requires a tool beyond the Go toolchain itself.

---

## 🎓 Learning Objectives

By the end of Level 35, you'll be able to:

- ✅ Explain why Go's statically-linked binaries are well-suited to minimal containers
- ✅ Write a multi-stage Dockerfile from scratch
- ✅ Choose an appropriate runtime base image (scratch, distroless, or alpine)
- ✅ Write a .dockerignore that excludes secrets and clutter
- ✅ Order Dockerfile instructions for layer-caching efficiency
- ✅ Run containers with mapped ports and injected environment variables
- ✅ Connect Docker's environment-variable injection to Level 34's configuration patterns
- ✅ Write a docker-compose.yml combining an app and a database service

---

## 📊 Statistics

- **Main Theory:** README.md covering Dockerfile basics through docker-compose.yml
- **Exercises:** 10 detailed exercises + 3 bonus challenges
- **Study Guide:** multi-stage build diagram, size comparison, caching diagram, instruction reference
- **Quick Reference:** One-page cheat sheet
- **Practice Time:** 4-5 hours
- **Difficulty:** ⭐⭐⭐ (Advanced)

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for comprehensive theory
2. Use STUDY_GUIDE.md alongside for the multi-stage diagram and instruction reference table
3. Reference QUICK_REFERENCE.md for fast lookup

### For Practice
1. Read exercise description
2. Copy provided bash/Dockerfile/YAML commands (or type them)
3. Follow step-by-step instructions
4. For the Go portions: verify output matches the real, captured results
5. For the Docker portions: run them yourself and trust your own output over the illustrative blocks

### For Review
1. Use QUICK_REFERENCE.md for quick recap
2. Refer to the common mistakes section
3. Practice writing a multi-stage Dockerfile from memory until it's automatic

---

## 🆘 Common Questions

**Q: Why does this level need a real Docker installation when no other level needed extra tools?**
A: Docker is a separate piece of software with its own daemon - unlike Gin (Level 28) or a SQL driver (Level 29), which are just Go packages fetched with `go get`, Docker isn't part of the Go toolchain at all. You need Docker Desktop (or an equivalent engine) installed and running to actually build/run the images in this level.

**Q: Why is a multi-stage build so much smaller than a single-stage one?**
A: The build stage needs the full Go compiler, standard library, and your source code to produce a binary. The runtime stage only needs to *run* that already-compiled binary - so it starts fresh from a minimal base and copies in just the one file, discarding the compiler and source entirely.

**Q: Why copy go.mod before the rest of the source code?**
A: Docker caches each instruction as a layer. If `go.mod`/`go.sum` haven't changed, `RUN go mod download` can reuse its cached layer even when your `.go` files have changed - so ordinary source edits don't force dependencies to re-download.

**Q: Why not just use `ENV` for secrets?**
A: `ENV` bakes its value into the image itself, permanently - anyone with the image can extract it, even from a layer a later instruction appears to overwrite. Secrets belong in `docker run -e`, mounted secret files, or a secrets manager, injected at runtime, never at build time.

**Q: How does this connect to Level 34's configuration?**
A: Level 34 teaches your Go program to read its settings from environment variables instead of hardcoding them. This level is what actually *sets* those environment variables in a deployed context - `ENV` in the Dockerfile for defaults, `docker run -e` (and later, Kubernetes manifests in Level 36) to override them per environment.

---

## 🎯 Before Moving to Level 36

Make sure you can answer these questions:

- [ ] Why are Go binaries well-suited to minimal containers, compared to Python/Node/Java?
- [ ] What's the difference between CMD and ENTRYPOINT?
- [ ] Why does a multi-stage build produce a dramatically smaller image?
- [ ] Why does .dockerignore matter even before any COPY instruction runs?
- [ ] Why should go.mod/go.sum be copied (and go mod download run) before the rest of the source?
- [ ] How do you inject an environment variable into a running container without rebuilding the image?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (10 min)
   - Gets you oriented, including the Docker-availability caveat

2. **README.md** Sections 1-3 (30 min)
   - Why containerize, Dockerfile basics, multi-stage builds

3. **README.md** Sections 4-7 (30 min)
   - .dockerignore, layer caching, building/running, environment variables

4. **EXERCISES.md** Exercises 1-5 (2 hours)
   - Get hands-on with single-stage, multi-stage, and .dockerignore

5. **STUDY_GUIDE.md** (30 min)
   - Study the multi-stage diagram and instruction reference table

6. **EXERCISES.md** Exercises 6-10 (2+ hours)
   - Deep practice with ENV, compose, HEALTHCHECK, and a full REST API

7. **QUICK_REFERENCE.md** (10 min)
   - Create your mental cheat sheet

---

## 🏆 Success Indicators

You've mastered Level 35 when:

- ✅ You reach for a multi-stage build by default, never a single-stage one
- ✅ You can explain rune-deep *why* the runtime image is so much smaller, not just that it is
- ✅ You never bake a secret into an image with ENV
- ✅ You always write .dockerignore before your first COPY . .
- ✅ You can write a docker-compose.yml combining an app and a database from memory
- ✅ You've completed 8+ exercises
- ✅ You can explain Docker to someone else

---

## 🚀 What's Next?

After Level 35, you're ready for:

**Level 36: Kubernetes**
- Running containers at scale across a cluster
- Deployments, Pods, and Services
- ConfigMaps and Secrets, building directly on this level's environment variables
- Scaling and self-healing beyond a single `docker run`

---

## 💬 Key Takeaway

> **A container ships the environment along with the code - and Go's statically-linked binaries mean that environment can be almost nothing at all.**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Exercises
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Patterns
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Happy learning! Level 35 gives your Go programs a passport that works in any environment! 🎉

*Estimated time to complete Level 35: 4-5 hours*
*Difficulty: ⭐⭐⭐ (Advanced)*
*Next Level: Level 36 - Kubernetes*
