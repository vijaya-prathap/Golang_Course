# Level 40: Real-World Projects - INDEX

Welcome to **Level 40: Real-World Projects** - the capstone! This is where every skill from Levels 1-39 stops being practiced in isolation and gets combined into one real, working, tested system.

---

## 📖 What You'll Learn

- ✅ How to design an API's endpoints and data model before writing any code
- ✅ How to build a complete Repository-Service-Handler stack from scratch (Level 30)
- ✅ How to wire the whole dependency graph by hand in a composition root (Level 31)
- ✅ How to protect specific endpoints with JWT auth while leaving others public (Level 32)
- ✅ How to add structured logging and fail-fast configuration as real middleware (Levels 33-34)
- ✅ How to bridge Level 39's rate-limiting theory into real, tested token-bucket middleware
- ✅ How to prove an entire system works with a full end-to-end test (Level 26)
- ✅ How to carry Levels 35-36's Docker/Kubernetes honesty pattern into your own projects

---

## 🗂️ Level 40 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview and welcome message. Start here first!

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- Why this level exists and what a "capstone" actually means
- Three capstone project options, and why this level commits to a URL Shortener
- The full API design and data model, decided before any code
- The complete architecture: Repository → Service → Handler → Middleware
- Building every layer, in order, with real captured output at each step
- Adding auth, logging, configuration, and a rate limiter
- Testing the whole system end to end
- Containerizing and sketching Kubernetes (following Levels 35-36's honesty pattern)
- Best practices and common mistakes for capstone-scale projects

**Read Time:** 60-75 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 sequential milestones building ONE continuous project:
1. Project setup, domain models, Link repository
2. User repository
3. LinkService - validation, per-user limit, short codes
4. UserService - bcrypt hashing, JWT issue/parse
5. Handler layer - routes, isolated handler tests
6. Real JWT auth middleware
7. Structured logging + configuration (catches a real bug!)
8. Token-bucket rate limiter
9. main.go wiring + full end-to-end test (catches a second real bug!)
10. Dockerfile + Kubernetes manifests (illustrative in this environment)

Plus 3 bonus challenges: click analytics, link expiration, custom aliases.

**Time Commitment:** 8-12 hours
**When:** After reading README.md's Sections 1-4

### 4. **STUDY_GUIDE.md** (Visual Learning)
- The complete architecture diagram, all layers shown together
- The full API endpoint reference table
- Two full request-flow diagrams (the JWT-protected path and the rate-limited path)
- A milestone roadmap diagram
- Common mistakes at a glance
- Prerequisites, readiness checklist, and pre-Level-41 checklist

**Read Time:** 30-40 minutes
**When:** Use while doing the milestones

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- Every endpoint, every layer's core code shape, one page
- The two real bugs this level caught, and their fixes
- Quick Docker/Kubernetes honesty-pattern reminder

**Read Time:** 10-15 minutes
**When:** Quick lookup while building

---

## 🎯 Recommended Learning Path

### Week 1: Foundations (Milestones 1-4)
1. Read **START_HERE.txt** (10 min)
2. Read **README.md** Sections 1-6 (35 min)
3. Complete Milestones 1-4 (3-4 hours)

### Week 2: HTTP and Security (Milestones 5-6)
1. Read **README.md** Sections 7-8 (15 min)
2. Complete Milestones 5-6 (2-3 hours)

### Week 3: Observability, Resilience, Proof (Milestones 7-9)
1. Read **README.md** Sections 9-11 (20 min)
2. Complete Milestones 7-9 (3-4 hours) - budget extra time here, both real bugs live in this stretch

### Week 4: Shipping and Consolidation (Milestone 10)
1. Read **README.md** Sections 12-15 (15 min)
2. Complete Milestone 10 and the bonus challenges
3. Review with **STUDY_GUIDE.md** and **QUICK_REFERENCE.md**

---

## 💡 Key Concepts At A Glance

### The Architecture
```go
Handler  →  Service  →  Repository (interface)
(HTTP)      (rules)      (in-memory, Level 30)
```

### The Composition Root
```go
func buildServer(cfg config.Config, logger *slog.Logger) http.Handler {
    linkRepo := repository.NewInMemoryLinkRepository()
    linkService := service.NewLinkService(linkRepo, cfg.CodeLength)
    linkHandler := handler.NewLinkHandler(linkService)
    // ... the ONLY place that knows any concrete type (Level 31)
}
```

### The Middleware Pipeline
```
Logging (every request) → RequireAuth (POST /shorten only) → Handler
Logging (every request) → RateLimit (GET /{code} only)      → Handler
```

---

## ✅ Prerequisites

Make sure you've completed **Level 39: System Design** and are comfortable with every level it builds on:

You need:
- ✅ Comfort with Level 30's Repository-Service-Handler layering
- ✅ Comfort with Level 31's constructor injection and composition roots
- ✅ Comfort with Level 32's JWT auth middleware and bcrypt password hashing
- ✅ Comfort with Level 33's structured logging (`log/slog`) and Level 34's configuration patterns
- ✅ Comfort with Level 26's table-driven tests and `httptest`
- ✅ Familiarity with Level 39's rate-limiting concepts (token bucket, in particular)

---

## 🎓 Learning Objectives

By the end of Level 40, you'll be able to:

- ✅ Design an API's endpoints, data model, and business rules before writing code
- ✅ Build and verify a Repository-Service-Handler stack, one layer at a time
- ✅ Wire an entire application's dependency graph by hand, readably
- ✅ Protect one endpoint with JWT auth while leaving others intentionally public
- ✅ Add structured logging and environment-driven configuration as real middleware
- ✅ Implement and test a token-bucket rate limiter from scratch
- ✅ Prove a whole system works with one end-to-end test, on top of per-layer unit tests
- ✅ Recognize the Docker/Kubernetes verification pattern this course has used since Level 35

---

## 📊 Statistics

- **Main Theory:** README.md covering project design, full architecture, and every layer's build
- **Exercises:** 10 sequential milestones (one evolving project) + 3 bonus challenges
- **Study Guide:** full architecture diagram, endpoint table, two request-flow diagrams, milestone roadmap
- **Quick Reference:** one-page cheat sheet covering every layer and both real bugs found during authoring
- **Practice Time:** 8-12 hours
- **Difficulty:** ⭐⭐⭐⭐⭐ (Capstone)

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for the full design rationale and architecture
2. Use STUDY_GUIDE.md alongside for diagrams while you build
3. Reference QUICK_REFERENCE.md for fast lookup mid-milestone

### For Practice
1. Read each milestone's Objective and Instructions in EXERCISES.md
2. Create the files exactly as shown, in the SAME evolving project directory
3. Run the verification commands and compare against Expected Output
4. When a milestone shows a real bug (Milestones 7 and 9), reproduce the failure before applying the fix - that's the point

### For Review
1. Use QUICK_REFERENCE.md for a fast recap of every layer
2. Walk through STUDY_GUIDE.md's request-flow diagrams from memory
3. Re-read README.md's Best Practices and Common Mistakes sections

---

## 🆘 Common Questions

**Q: Why in-memory storage instead of Level 29's SQLite approach?**
A: Network access for `modernc.org/sqlite` was verified working, but keeping the capstone's storage layer 100% offline (matching every level except 28, 29, 32) means the WHOLE build works with zero network dependency for storage - and Level 30's own argument is that swapping in SQLite later is a pure drop-in behind the same repository interface.

**Q: Why does this level use bcrypt and JWT if the rest of the course avoids network dependencies?**
A: Level 32 already established this exact, clearly-flagged exception - both packages are pinned versions, fetched once, then cached. This level follows that precedent rather than re-deciding it.

**Q: Why is `GET /{code}` rate-limited but `GET /stats/{code}` isn't?**
A: `GET /{code}` is the one endpoint the entire internet can hit without any credentials - it's the realistic abuse target. `GET /stats/{code}` matters less for abuse and more for correctness (never incrementing a click count by accident), which is why it's a *separate* Service method instead.

**Q: What were the two real bugs this level's authoring caught?**
A: `os.LookupEnv` correctly treating an explicitly-empty environment variable as "present" (Milestone 7), and an `http.Client` never reusing a connection - and therefore never keeping a consistent rate-limit key - when a response body isn't drained before `Close()` (Milestone 9). Both are walked through in full in EXERCISES.md, failure and fix.

**Q: Why aren't the Docker/Kubernetes steps verified the same way the Go code is?**
A: Because they can't be, honestly, in this environment - `docker version` shows the CLI installed with no reachable daemon, and `kubectl` has no reachable cluster. This is the exact situation Levels 35 and 36 already documented; this level follows their established pattern rather than pretending otherwise.

---

## 🎯 Before Moving to Level 41

Make sure you can answer these questions:

- [ ] What are the three layers of this project's architecture, and what does each one own?
- [ ] Why does `main.go`'s `buildServer` function get to know concrete types when nothing else in the project does?
- [ ] How does the JWT auth middleware get the authenticated user's ID to the handler?
- [ ] Why is the rate limiter keyed by `RemoteAddr`, and what real bug did that expose during testing?
- [ ] What's the difference between `LinkService.Resolve` and `LinkService.Stats`, and why does it matter?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (10 min)
   - Gets you oriented on the capstone's scope

2. **README.md** Sections 1-4 (30 min)
   - Why this level exists, project options, API design, architecture

3. **README.md** Sections 5-11 (35 min)
   - Building every layer, auth, observability, rate limiting, testing

4. **EXERCISES.md** Milestones 1-6 (4-5 hours)
   - Build the repository, service, handler, and auth middleware

5. **STUDY_GUIDE.md** (30 min)
   - Study the architecture and request-flow diagrams

6. **EXERCISES.md** Milestones 7-10 (4-5 hours)
   - Logging, config, rate limiting, end-to-end testing, and packaging - including both real bugs

7. **QUICK_REFERENCE.md** (10 min)
   - Consolidate into your own mental cheat sheet

---

## 🏆 Success Indicators

You've mastered Level 40 when:

- ✅ You can build a Repository-Service-Handler stack from a blank directory without referring back to notes
- ✅ You instinctively reach for an interface at every layer boundary, and a concrete type only in the composition root
- ✅ You can explain JWT middleware, logging middleware, and rate-limiting middleware as three instances of the exact same `func(http.Handler) http.Handler` shape
- ✅ You've run `go test ./...` and watched every package pass, including one full end-to-end test
- ✅ You can explain both real bugs this level's authoring hit, in your own words, to someone else
- ✅ You've completed all 10 milestones and at least one bonus challenge

---

## 🚀 What's Next?

After Level 40, you're ready for:

**Level 41: Go Interview Preparation**
- Turning this capstone's architecture into confident, spoken interview answers
- Common Go interview topics: concurrency, interfaces, error handling, the memory model
- Practicing system-design questions using this exact project as your worked example

---

## 💬 Key Takeaway

> **Every level before this one taught a tool. This level is the first time you built something real with all of them at once - and that, not any single syntax feature, is what "knowing Go" actually means.**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Milestones
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Architecture & Diagrams
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Congratulations on reaching the capstone! Level 40 is where everything you've learned finally works together as one real system. 🎉

*Estimated time to complete Level 40: 8-12 hours*
*Difficulty: ⭐⭐⭐⭐⭐ (Capstone)*
*Prerequisite: Level 39 - System Design*
*Next Level: Level 41 - Go Interview Preparation*
