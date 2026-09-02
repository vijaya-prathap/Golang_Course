# Level 38: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: Profiling
```
Day 1:  Why production debugging is different; net/http/pprof setup
Day 2:  CPU profiling: load generation, go tool pprof -top
Day 3:  Memory profiling: heap snapshots, inuse_space
Day 4:  Goroutine profiling: dumps, leak signatures
Day 5:  runtime.NumGoroutine() / ReadMemStats() lightweight checks
Day 6:  Common production issues: leaks, deadlocks
Day 7:  The Delve debugger
```

### Week 2: Practice & Application
```
Day 1:  Exercises 1-3 (pprof setup, CPU profiling, memory profiling)
Day 2:  Exercises 4-6 (goroutine profiling, lightweight diagnostics, health checks)
Day 3:  Exercises 7-8 (correlation IDs, Delve)
Day 4:  Exercises 9-10 (deadlocks, comprehensive service)
Day 5:  Bonus challenges
Day 6-7: Practice & review
```

---

## 🎯 The pprof Endpoint Reference Table

| Endpoint | Method | What it returns | Blocks? |
|----------|--------|------------------|---------|
| `/debug/pprof/` | GET | HTML index of every available profile | No |
| `/debug/pprof/profile?seconds=N` | GET | Binary CPU profile, sampled over N seconds | **Yes, for N seconds** |
| `/debug/pprof/heap` | GET | Snapshot of current heap allocations | No |
| `/debug/pprof/heap?gc=1` | GET | Same, but forces a GC first (live-only) | No (brief GC pause) |
| `/debug/pprof/allocs` | GET | Historical allocation profile (includes freed memory) | No |
| `/debug/pprof/goroutine` | GET | Binary goroutine profile (for `go tool pprof`) | No |
| `/debug/pprof/goroutine?debug=2` | GET | Full human-readable stack dump, every goroutine | No |
| `/debug/pprof/threadcreate` | GET | OS thread creation profile | No |
| `/debug/pprof/block` | GET | Blocking profile (requires `runtime.SetBlockProfileRate`) | No |
| `/debug/pprof/mutex` | GET | Mutex contention profile (requires `runtime.SetMutexProfileFraction`) | No |
| `/debug/pprof/trace?seconds=N` | GET | Execution trace for `go tool trace` | **Yes, for N seconds** |
| `/debug/pprof/cmdline` | GET | The process's command-line arguments | No |
| `/debug/pprof/symbol` | GET/POST | Resolves program counters to function names | No |

```
IMPORTANT: importing net/http/pprof for its side effect registers ALL
of the above ONLY on http.DefaultServeMux. A custom http.NewServeMux()
gets NONE of them unless you register the exported handlers yourself
(pprof.Index, pprof.Profile, pprof.Symbol, pprof.Trace, pprof.Cmdline).
```

---

## 🧭 CPU vs Memory vs Goroutine Profiling: Decision Diagram

```
Something's wrong in production. Which profile do you reach for?

"The service is slow / using too much CPU"
    │
    ▼
CPU PROFILE  →  GET /debug/pprof/profile?seconds=30
    │             (capture WHILE the slowness is happening -
    │              idle-time profiles tell you nothing useful)
    ▼
go tool pprof -top  →  read flat% top-down, that's your hot function


"Memory keeps growing / OOM-killed"
    │
    ▼
HEAP PROFILE  →  GET /debug/pprof/heap?gc=1
    │             (the ?gc=1 forces garbage collection first, so
    │              you see LIVE memory, not garbage awaiting sweep)
    ▼
go tool pprof -top -sample_index=inuse_space
    │
    ▼
Growing over time? Capture two snapshots, hours/days apart, then:
go tool pprof -diff_base=old.pprof new.pprof   →  isolates the growth


"Goroutine count keeps climbing"
    │
    ▼
GOROUTINE PROFILE  →  GET /debug/pprof/goroutine?debug=2
    │                  (human-readable: shows every stack trace)
    ▼
Look for one function name accumulating in [chan receive],
[chan send], or [select] states - that's your leak's origin


"The whole process just crashed with a fatal error"
    │
    ▼
NOT a profiling problem - read the crash output directly.
"fatal error: all goroutines are asleep - deadlock!" IS the
diagnostic; Go's runtime already printed every blocked goroutine
for you before it exited.
```

---

## 🏥 Health Check Design Diagram

```
BAD (tells you nothing):

    GET /healthz  →  always 200 OK

    This only proves the HTTP listener accepts connections.
    It says NOTHING about whether the service can do its job.


GOOD (checks real dependencies):

    GET /healthz
        │
        ├─ Can we reach the database?          ──── dbUp.Load()
        ├─ Is the goroutine count sane?         ──── runtime.NumGoroutine()
        ├─ (optionally) Can we reach a cache,
        │   a downstream API, a message queue?
        │
        ▼
    ALL checks pass?  →  200 OK,  {"status":"ok", ...real fields...}
    ANY check fails?  →  503 Service Unavailable,  {"status":"unhealthy", ...}

    Real captured behavior from this level's verification:

    [db up]        HTTP 200  {Status:ok        DBConnected:true  NumGoroutine:6}
    [db down]      HTTP 503  {Status:unhealthy DBConnected:false NumGoroutine:6}
    [db recovered] HTTP 200  {Status:ok        DBConnected:true  NumGoroutine:6}
```

```
Why this matters for Level 36 (Kubernetes):

    readinessProbe / livenessProbe  →  hit exactly this endpoint

    A /healthz that always says "ok" means Kubernetes keeps routing
    traffic to (or never restarts) a Pod that is actually broken -
    the probe is only as honest as the endpoint it calls.
```

---

## 🚨 Common Production Issues and Their Profiling Signatures

| Issue | What you observe | What the tool shows |
|-------|-------------------|----------------------|
| **Goroutine leak** | `NumGoroutine()` rises steadily, never falls, even at idle | Goroutine dump: one function repeatedly in `[chan receive]`/`[select]`, count growing run over run |
| **Memory leak** | `Alloc` / RSS rises steadily, a GC cycle doesn't bring it back down | `inuse_space` heap profile: one function's share keeps growing across successive snapshots (`-diff_base` isolates it) |
| **CPU-bound hot path** | High CPU usage, high latency under load | CPU profile: one function dominates `flat%`, matching the code path you'd expect to be slow |
| **I/O-bound latency (not CPU)** | High latency, but CPU usage looks normal/low | CPU profile shows mostly runtime/syscall frames (network, disk) rather than your own code - profile isn't the right tool; check latency/timing logs instead |
| **Deadlock** | Process crashes immediately, doesn't hang | `fatal error: all goroutines are asleep - deadlock!` with a full goroutine dump already printed |
| **Excessive GC pauses** | Latency spikes correlating with GC | `NumGC` climbing fast in `ReadMemStats`; high allocation rate in an `alloc_space` profile |

---

## 🚨 Common Mistakes

### Mistake 1: pprof Endpoints Exposed Publicly

```
❌ WRONG: binding pprof's mux to 0.0.0.0 or a public listener
✅ RIGHT: 127.0.0.1 only, or an internal-only network/auth layer
```

### Mistake 2: A Custom Mux "Losing" pprof

```go
// ❌ WRONG - pprof registered on DefaultServeMux, never mounted here
mux := http.NewServeMux()
mux.HandleFunc("/work", handler)
http.ListenAndServe("127.0.0.1:8080", mux) // /debug/pprof/* all 404

// ✅ RIGHT - use DefaultServeMux (directly or via http.HandleFunc)
http.HandleFunc("/work", handler)
http.ListenAndServe("127.0.0.1:8080", http.DefaultServeMux)
```

### Mistake 3: Profiling an Idle Process

```
❌ WRONG: capturing a CPU profile while nothing is happening
✅ RIGHT: profile WHILE load (real or simulated) is hitting the service
```

### Mistake 4: Ignoring a Slow Upward Trend

```
❌ WRONG: "goroutine count is a little higher than yesterday, probably fine"
✅ RIGHT: alert on the TREND (a leak detector logging when count exceeds
          a threshold), not just a single point-in-time snapshot
```

### Mistake 5: No Correlation ID

```
❌ WRONG: five ERROR log lines during an incident, no way to tell if
          that's one failing request or five
✅ RIGHT: every log line for a request carries the same request_id,
          assigned once at the boundary (middleware) and threaded
          through via context.Context
```

---

## 📈 Progression Summary

### Understanding Level 38

Level 38 teaches how to see inside a running Go process without stopping it:

1. **Profiling endpoints** — CPU, heap, goroutine, all via plain HTTP
2. **go tool pprof** — reading `-top` output, `inuse_space` vs `alloc_space`, `-diff_base`
3. **Lightweight diagnostics** — `NumGoroutine`/`ReadMemStats` for cheap, constant checks
4. **Issue signatures** — recognizing leaks and deadlocks by their real, characteristic shape
5. **Delve** — a real debugger for a real running process, not just `fmt.Println`
6. **Health and logging** — reporting true state and tracing incidents end to end

### Prerequisites for Level 39

Before moving to Level 39 (System Design), you need:

- ✅ Comfortable wiring up `net/http/pprof` and knowing its `DefaultServeMux` gotcha
- ✅ Can capture and read a real CPU profile with `go tool pprof -top`
- ✅ Can capture and read a real heap profile and identify a growing allocator
- ✅ Can capture a goroutine dump and identify a leak's signature
- ✅ Comfortable with `runtime.NumGoroutine()` and `runtime.ReadMemStats()`
- ✅ Can recognize a real deadlock crash and a real leak trend on sight
- ✅ Has used (or read, if `dlv` wasn't available) a real Delve debugging session
- ✅ Can build a health check that reports real dependency state
- ✅ Understands why correlation IDs matter for tracing an incident

### Ready for Level 39?

Level 39 moves up a level of altitude - from debugging one running service to designing systems made of many:
- Trade-offs between architectural styles
- Scaling, reliability, and consistency at a system level
- Applying this course's deployment (Levels 35-37) and observability (this level) foundations to real architecture decisions

---

## ✅ Checklist Before Level 39

- [ ] Can wire up `net/http/pprof` and explain why it must be on `DefaultServeMux`
- [ ] Can capture and interpret a real CPU profile under load
- [ ] Can capture and interpret a real heap profile, including `-diff_base`
- [ ] Can capture a goroutine dump and spot a leak's signature
- [ ] Comfortable using `runtime.NumGoroutine()`/`ReadMemStats()` for lightweight checks
- [ ] Can recognize `fatal error: all goroutines are asleep - deadlock!` immediately
- [ ] Has set a breakpoint and inspected a variable with Delve (or read a real session)
- [ ] Can design a `/healthz` endpoint around real dependency checks
- [ ] Understands correlation IDs and can implement one with middleware + context
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The Access
Production debugging tools work against a process that's already running - no restart, no pre-planted debugger, no reproducing locally required.

### The Endpoints
`net/http/pprof`'s handlers are just HTTP - `curl`, a browser, or Go code can all pull a profile from a live process, as long as it's reachable and not exposed where it shouldn't be.

### The Signatures
A leak isn't abstract - it's a specific function name accumulating in a specific blocked state, visible in a real dump. A deadlock isn't a hang - it's a loud, immediate crash with the evidence already printed.

### The Honesty
A health check, a log line, and a profile are all only useful if they reflect real state - a hardcoded `200`, a log line missing its correlation ID, or a profile captured at idle all fail the same way: they tell you nothing true about what's actually happening.

---

## 📚 Next Level

Level 39: System Design
- Designing systems at a higher level than a single service
- Applying observability, deployment, and debugging practices at architectural scale

You've got production debugging down! Keep going! 🚀
