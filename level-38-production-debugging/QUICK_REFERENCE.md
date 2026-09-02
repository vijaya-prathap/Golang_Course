# Level 38: Quick Reference Card

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir myapp && cd myapp
go mod init myapp

# Create main.go
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "net/http"
    _ "net/http/pprof" // side-effect import - registers on DefaultServeMux
)

func main() {
    http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "hello")
    })
    http.ListenAndServe("127.0.0.1:6060", nil) // loopback only!
}
EOF

# Run, then in another terminal:
go run main.go &
curl http://127.0.0.1:6060/debug/pprof/
```

---

## 📋 The pprof Endpoints

```
/debug/pprof/                      index page (HTML)
/debug/pprof/profile?seconds=30    CPU profile (BLOCKS for N seconds)
/debug/pprof/heap                  live heap snapshot
/debug/pprof/heap?gc=1             ...same, but forces GC first
/debug/pprof/allocs                all-time allocation profile
/debug/pprof/goroutine             goroutine profile (pprof format)
/debug/pprof/goroutine?debug=2     full human-readable stack dump
/debug/pprof/trace?seconds=5       execution trace (for go tool trace)
```

⚠️ Only registers on `http.DefaultServeMux`. Custom mux? Register `pprof.Index`, `pprof.Profile`, `pprof.Symbol`, `pprof.Trace`, `pprof.Cmdline` yourself.

---

## 🔬 go tool pprof Cheat Sheet

```bash
# Capture a CPU profile from a running server
curl -o cpu.pprof "http://127.0.0.1:6060/debug/pprof/profile?seconds=30"

# Capture a heap profile
curl -o heap.pprof "http://127.0.0.1:6060/debug/pprof/heap?gc=1"

# Top functions by flat time/space
go tool pprof -top -nodecount=10 ./binary cpu.pprof
go tool pprof -top -sample_index=inuse_space ./binary heap.pprof

# Compare two heap snapshots (find what GREW)
go tool pprof -top -sample_index=inuse_space -diff_base=old.pprof ./binary new.pprof

# Interactive mode (then type: top, list <func>, web, quit)
go tool pprof ./binary cpu.pprof
```

---

## 🧵 Lightweight Diagnostics

```go
runtime.NumGoroutine()   // int - current goroutine count, instant

var m runtime.MemStats
runtime.ReadMemStats(&m)
m.Alloc        // currently live heap bytes
m.TotalAlloc   // all-time allocated bytes (never decreases)
m.Sys          // total memory obtained from the OS
m.NumGC        // number of completed GC cycles
```

---

## 🚨 Issue Signatures

| Symptom | Real signature |
|---------|-----------------|
| Goroutine leak | `NumGoroutine()` rises, never falls; dump shows one function repeating in `[chan receive]` |
| Memory leak | `Alloc` rises, GC doesn't reclaim it; `-diff_base` shows one function's share growing |
| Deadlock | `fatal error: all goroutines are asleep - deadlock!` (crashes immediately, doesn't hang) |
| CPU-bound slowness | One function dominates `flat%` in a CPU profile taken under real load |

---

## 🐛 Delve (dlv) Essentials

```bash
dlv version                                    # confirm it's installed
go build -gcflags="all=-N -l" -o app .         # disable inlining/optimization first
dlv exec ./app                                 # debug from process start
dlv attach <pid>                               # attach to an ALREADY-RUNNING process
dlv test ./...                                 # debug inside a test
```

```
(dlv) break main.someFunc     # set a breakpoint
(dlv) continue                # run until breakpoint/exit
(dlv) print someVar           # inspect a variable
(dlv) next                    # step over one line
(dlv) bt                      # backtrace
(dlv) quit                    # exit (answers a kill/detach prompt)
```

---

## 🏥 Real Health Check Pattern

```go
func healthzHandler(w http.ResponseWriter, r *http.Request) {
    healthy := dbUp.Load() && runtime.NumGoroutine() < maxHealthyGoroutines
    if healthy {
        w.WriteHeader(http.StatusOK)          // 200
    } else {
        w.WriteHeader(http.StatusServiceUnavailable) // 503
    }
    json.NewEncoder(w).Encode(status)
}
```

Check real dependencies (DB, goroutine count, downstream services) - never return a hardcoded `200`.

---

## 📝 Correlation ID Pattern

```go
func withRequestLogger(logger *slog.Logger, next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        reqID := r.Header.Get("X-Request-ID")
        if reqID == "" {
            reqID = newRequestID()
        }
        l := logger.With("request_id", reqID)
        l.Info("request started", "path", r.URL.Path)
        next(w, r.WithContext(context.WithValue(r.Context(), requestIDKey, reqID)))
        l.Info("request completed", "path", r.URL.Path)
    }
}
```

`grep '"request_id":"<id>"'` across logs reconstructs one request's full life.

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Public pprof | binding to `0.0.0.0` | bind to `127.0.0.1`, or gate behind auth |
| Lost pprof on custom mux | `http.NewServeMux()` + expecting `/debug/pprof/*` to work | use `http.DefaultServeMux`, or register pprof handlers explicitly |
| Idle profiling | capturing a profile with no load | profile while real/simulated load is running |
| Ignoring trends | one high goroutine-count reading, ignored | alert on the trend (rising over time), not a single point |
| No correlation ID | untagged log lines during an incident | tag every line with the same `request_id` |

---

## 🎓 Before Next Level

Can you:
- [ ] Wire up `net/http/pprof` and explain the `DefaultServeMux` gotcha?
- [ ] Capture and read a real CPU profile with `go tool pprof -top`?
- [ ] Capture and read a real heap profile, including comparing two with `-diff_base`?
- [ ] Capture a goroutine dump and recognize a leak's signature?
- [ ] Use `runtime.NumGoroutine()`/`ReadMemStats()` confidently?
- [ ] Recognize a real deadlock crash immediately?
- [ ] Set a breakpoint and inspect a variable with Delve (or read a real session)?
- [ ] Design a `/healthz` endpoint around real dependency checks?
- [ ] Explain why correlation IDs matter for incident logs?

If YES → You're ready for Level 39!

---

## 📚 Next Level

Level 39: System Design
- Designing systems at a higher level than a single service
- Applying observability, deployment, and debugging practices at architectural scale

You've got production debugging down! 💪
