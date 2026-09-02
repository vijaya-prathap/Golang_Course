# Level 38: Production Debugging - INDEX

Welcome to **Level 38: Production Debugging**! This is where you learn to see inside a running Go process - profiling CPU, memory, and goroutines, checking real health state, tracing incidents through structured logs, and using a real debugger against a process you didn't start under one.

---

## 📖 What You'll Learn

- ✅ Why production debugging needs different tools than local development
- ✅ CPU profiling with `net/http/pprof` and `go tool pprof`
- ✅ Memory (heap) profiling and confirming predictable allocations
- ✅ Goroutine profiling and recognizing a leak's real signature
- ✅ Lightweight diagnostics with `runtime.NumGoroutine()` and `ReadMemStats()`
- ✅ Common production issues: goroutine leaks, memory leaks, deadlocks
- ✅ The Delve debugger: breakpoints and inspection, even on an already-running process
- ✅ Health check endpoints that report real internal state
- ✅ Structured logging for incident response, with correlation IDs

---

## 🗂️ Level 38 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview and welcome message. Start here first!

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- Why production debugging is different
- CPU, memory, and goroutine profiling with `net/http/pprof` - all real, captured output
- `runtime.NumGoroutine()`/`ReadMemStats()` for lightweight diagnostics
- Recognizing leaks and deadlocks by their real signatures
- The Delve debugger, used for real in this sandbox
- Real health checks and structured incident logging
- Best practices
- Common mistakes

**Read Time:** 55-70 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 detailed, step-by-step exercises:
1. Setting up net/http/pprof
2. Capturing and analyzing a real CPU profile
3. Memory profiling with predictable allocations
4. Goroutine profiling with blocked goroutines
5. Lightweight diagnostics with NumGoroutine and ReadMemStats
6. A real health check endpoint
7. Structured incident logging with correlation IDs
8. Debugging with Delve
9. Diagnosing a deadlock
10. Comprehensive practice — a fully observable service

Plus 3 bonus challenges: a goroutine-leak detector, comparing heap profiles with `-diff_base`, and attaching Delve to a running process.

**Time Commitment:** 6-8 hours
**When:** After reading README

### 4. **STUDY_GUIDE.md** (Visual Learning)
- The pprof endpoint reference table
- CPU vs memory vs goroutine profiling decision diagram
- Health check design diagram
- Common production issues and their profiling signatures table
- Common mistakes guide

**Read Time:** 35-45 minutes
**When:** Use while doing exercises

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- Essential pprof endpoint syntax
- `go tool pprof` command reference
- Delve essentials
- Health check and correlation ID patterns
- Common mistakes table

**Read Time:** 10-15 minutes
**When:** Quick lookup during practice

---

## 🎯 Recommended Learning Path

### Day 1: Why Production Debugging Is Different, and pprof Setup (2 hours)
1. Read **START_HERE.txt** (10 min)
2. Read **README.md** Sections 1-2 (30 min)
3. Complete Exercise 1 (30 min)

### Day 2: CPU and Memory Profiling (2.5 hours)
1. Read **README.md** Sections 2-3 (25 min)
2. Complete Exercises 2-3 (2 hours)

### Day 3: Goroutine Profiling and Lightweight Diagnostics (2 hours)
1. Read **README.md** Sections 4-5 (25 min)
2. Complete Exercises 4-5 (1.5 hours)

### Day 4: Health Checks and Structured Logging (2 hours)
1. Read **README.md** Sections 8-9 (25 min)
2. Complete Exercises 6-7 (1.5 hours)

### Day 5: Delve and Deadlocks (2 hours)
1. Read **README.md** Sections 6-7 (25 min)
2. Complete Exercises 8-9 (1.5 hours)

### Day 6: Comprehensive Practice (2 hours)
1. Complete Exercise 10 (1.5 hours)
2. Try the bonus challenges

### Day 7: Consolidation (1 hour)
1. Review with **STUDY_GUIDE.md** and **QUICK_REFERENCE.md**
2. Make sure you can explain every issue signature from memory

---

## 💡 Key Concepts At A Glance

### Profiling Endpoints
```go
import _ "net/http/pprof" // registers /debug/pprof/* on http.DefaultServeMux
```
```
/debug/pprof/profile?seconds=N   CPU profile
/debug/pprof/heap?gc=1           live heap snapshot
/debug/pprof/goroutine?debug=2   full goroutine stack dump
```

### Lightweight Diagnostics
```go
runtime.NumGoroutine()
runtime.ReadMemStats(&m)   // m.Alloc, m.NumGC, ...
```

### A Real Deadlock
```
fatal error: all goroutines are asleep - deadlock!
```

### Delve
```bash
dlv exec ./app        # debug from the start
dlv attach <pid>       # attach to a running process
```

### Real Health Check
```go
healthy := dbUp.Load() && runtime.NumGoroutine() < maxHealthyGoroutines
```

---

## ✅ Prerequisites

Make sure you've completed **Level 37: CI/CD**

You need:
- ✅ Comfort with goroutines and channels (Level 20-21)
- ✅ Comfort with `sync.Mutex`/`WaitGroup` and real deadlock output (Level 23)
- ✅ Comfort with `log/slog` structured logging (Level 33)
- ✅ Comfort building HTTP servers and handlers (Level 27, Level 28)
- ✅ A general familiarity with what a deployed, running service looks like (Levels 35-37)

---

## 🎓 Learning Objectives

By the end of Level 38, you'll be able to:

- ✅ Explain why production debugging tools have to work against an already-running process
- ✅ Wire up `net/http/pprof` correctly and avoid its `http.DefaultServeMux` gotcha
- ✅ Capture and interpret real CPU, heap, and goroutine profiles
- ✅ Use `runtime.NumGoroutine()`/`ReadMemStats()` for cheap, constant diagnostics
- ✅ Recognize goroutine leaks, memory leaks, and deadlocks by their real signatures
- ✅ Use the Delve debugger against both a fresh and an already-running process
- ✅ Design a health check endpoint that reports genuine internal state
- ✅ Implement correlation IDs to trace one request across an entire incident's logs

---

## 📊 Statistics

- **Main Theory:** README.md covering profiling, diagnostics, debugging, health, and logging - all real, captured output
- **Exercises:** 10 detailed exercises + 3 bonus challenges
- **Study Guide:** pprof reference table, decision diagrams, issue-signature table
- **Quick Reference:** One-page cheat sheet
- **Practice Time:** 6-8 hours
- **Difficulty:** ⭐⭐⭐⭐ (Advanced)

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for comprehensive theory
2. Use STUDY_GUIDE.md alongside for the decision diagrams and signature tables
3. Reference QUICK_REFERENCE.md for fast lookup

### For Practice
1. Read exercise description
2. Copy provided bash commands (or type them)
3. Follow step-by-step instructions
4. Verify output matches expected results (allowing for machine-specific numbers)
5. Understand what each step does

### For Review
1. Use QUICK_REFERENCE.md for quick recap
2. Refer to the common mistakes section
3. Practice explaining each issue signature (leak, deadlock, hot path) from memory

---

## 🆘 Common Questions

**Q: Do I need a live production server to do this level?**
A: No - every exercise runs a real local Go program that profiles, diagnoses, or logs itself. The techniques transfer directly to a real deployed service; only the "how do I reach `/debug/pprof/` on a remote host" part changes (usually via an internal network, SSH tunnel, or port-forward).

**Q: Why did importing net/http/pprof not expose any endpoints for me?**
A: You almost certainly built your own `http.NewServeMux()` and used it as your server's handler. `net/http/pprof`'s side-effect import only registers on `http.DefaultServeMux`. Use `http.HandleFunc`/`http.ListenAndServe(addr, nil)`, or register `pprof.Index`/`pprof.Profile`/etc. on your own mux explicitly.

**Q: Is it safe to leave /debug/pprof/ enabled in production?**
A: Only if it's not reachable from the public internet - bind it to loopback, put it on an internal-only network, or gate it behind authentication. Exposed publicly, it's a real security issue (source paths, memory contents, and a free CPU-profile-triggered denial-of-service lever).

**Q: What if `dlv` isn't installed on my machine?**
A: Install it from https://github.com/go-delve/delve, or read Section 7 and Exercise 8 as a reference - the commands and output shown are real, captured sessions from a sandbox where `dlv` (version 1.27.0) was available.

**Q: What's the difference between a goroutine leak and a deadlock?**
A: A deadlock means *every* goroutine in the process is permanently blocked - Go's runtime detects this and crashes immediately with `fatal error: all goroutines are asleep - deadlock!`. A goroutine leak means *some* goroutines are permanently blocked while the rest of the program keeps running fine - the process doesn't crash, it just slowly accumulates dead weight until something else breaks (usually memory).

---

## 🎯 Before Moving to Level 39

Make sure you can answer these questions:

- [ ] Why can't you just `fmt.Println` your way through a production incident?
- [ ] What's the `http.DefaultServeMux` gotcha with `net/http/pprof`, and how do you avoid it?
- [ ] How do you tell a genuine memory leak from a normal, healthy allocation pattern?
- [ ] What does `fatal error: all goroutines are asleep - deadlock!` actually mean, and why does it crash instead of hang?
- [ ] How would you attach a debugger to a process that's already running?
- [ ] Why is a health check that always returns `200` worse than having no health check at all?
- [ ] What problem does a correlation ID solve that structured logging alone doesn't?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (10 min)
   - Gets you oriented

2. **README.md** Sections 1-2 (25 min)
   - Why production debugging is different, pprof setup

3. **README.md** Sections 3-5 (30 min)
   - Memory, goroutine profiling, lightweight diagnostics

4. **EXERCISES.md** Exercises 1-5 (3 hours)
   - Get hands-on with every profiling technique

5. **STUDY_GUIDE.md** (35 min)
   - Study the decision diagrams and signature tables

6. **README.md** Sections 6-9 (30 min)
   - Issue signatures, Delve, health checks, structured logging

7. **EXERCISES.md** Exercises 6-10 (3+ hours)
   - Deep practice with health checks, logging, Delve, and the comprehensive service

8. **QUICK_REFERENCE.md** (15 min)
   - Create your mental cheat sheet

---

## 🏆 Success Indicators

You've mastered Level 38 when:

- ✅ You reach for `net/http/pprof` by default when a service is slow or leaking, instead of guessing
- ✅ You can read `go tool pprof -top` output and go straight to the real hot function
- ✅ You never assume a health check endpoint that returns `200` is actually telling the truth
- ✅ You recognize a real deadlock crash and a real goroutine-leak trend on sight
- ✅ You know how to attach Delve to a process you didn't start under a debugger
- ✅ You never ship a service without a correlation ID threading through its logs
- ✅ You've completed 8+ exercises

---

## 🚀 What's Next?

After Level 38, you're ready for:

**Level 39: System Design**
- Designing systems at a higher level than a single service
- Applying observability (this level), deployment (Levels 35-37), and architectural trade-offs together

---

## 💬 Key Takeaway

> **A production process doesn't need to be stopped, restarted, or reproduced locally to be understood - `net/http/pprof`, lightweight runtime diagnostics, a real debugger, an honest health check, and a correlation ID are enough to see exactly what it's doing, right now, while it keeps serving traffic.**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Exercises
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Patterns
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Happy learning! Level 38 turns a running Go process from a black box into something you can actually see inside! 🎉

*Estimated time to complete Level 38: 6-8 hours*
*Difficulty: ⭐⭐⭐⭐ (Advanced)*
*Next Level: Level 39 - System Design*
