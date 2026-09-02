# Level 38: Production Debugging - Complete Guide

> **Verification note:** Unlike Level 35 (Docker) and Level 36 (Kubernetes), every tool in this level ships in the Go standard toolchain - `net/http/pprof`, `runtime/pprof`, `go tool pprof`, `runtime.NumGoroutine`, `runtime.ReadMemStats` - and needs zero external dependencies. Everything below was actually run in this sandbox: real HTTP servers, real captured `.pprof` files, real `go tool pprof -top` output, real goroutine dumps, real `ReadMemStats` numbers. The Delve debugger (`dlv`) was also present in this sandbox (`dlv version` reported `1.27.0`), so Section 7's debugger walkthrough is real, captured output too - not a conceptual placeholder. Exact numbers (CPU percentages, allocation byte counts, timings) will differ on your machine; that variability is normal and expected for profiling data.

## Introduction

Welcome to Level 38! Level 37 (CI/CD) automated getting your code from a commit to a running deployment. But shipping code is not the end of the story - once it's running in production, serving real traffic, on a machine you probably can't `attach` a local debugger to whenever you feel like it, a whole new category of problem shows up: the program runs slow, or it runs out of memory, or it deadlocks, or a goroutine count keeps climbing until the process falls over at 3 AM. None of your earlier levels' tools - `fmt.Println`, a local debugger breakpoint, `go run` on your laptop - work the same way against a live, remote process serving real users.

This level is about the tools Go gives you to answer "what is this running process actually doing right now, and why" without stopping it, without a debugger attached from the start, and often without being able to reproduce the problem locally at all: CPU and memory profiling with `net/http/pprof`, goroutine diagnostics, lightweight runtime introspection, real health checks, structured incident logging, and the Delve debugger for the cases where you do get to attach one.

Difficulty: ⭐⭐⭐⭐ (Advanced).

---

## Table of Contents

1. [Why Production Debugging Is Different](#why-production-debugging-is-different)
2. [CPU Profiling With net/http/pprof](#cpu-profiling-with-nethttppprof)
3. [Memory Profiling](#memory-profiling)
4. [Goroutine Profiling](#goroutine-profiling)
5. [Lightweight Diagnostics: NumGoroutine and ReadMemStats](#lightweight-diagnostics-numgoroutine-and-readmemstats)
6. [Common Production Issues and How to Spot Them](#common-production-issues-and-how-to-spot-them)
7. [The Delve Debugger](#the-delve-debugger)
8. [Health Check Endpoints](#health-check-endpoints)
9. [Structured Logging for Incident Response](#structured-logging-for-incident-response)
10. [Best Practices](#best-practices)
11. [Common Mistakes](#common-mistakes)

---

## Why Production Debugging Is Different

Every debugging technique you've used so far assumes you can stop the program, or that you started it under a debugger, or at minimum that you can reproduce the failure on your own machine. Production breaks all three assumptions:

- **No debugger attached from the start.** A production server has usually been running for hours or days before anything goes wrong. You cannot go back in time and start it under `dlv debug`. Whatever tools you use have to work against a process that is already running, right now, serving traffic you don't want to interrupt.
- **You often can't reproduce it locally.** The bug might depend on production data volume, real network latency, a specific client's request shape, or load levels your laptop never sees. "Works on my machine" is exactly the failure mode this level exists to avoid depending on.
- **Observability has to be captured in the moment.** If you don't have logs, metrics, or a profile from *while the problem was happening*, the evidence is gone the moment it passes (or the moment the process restarts). Level 33's structured logging and this level's profiling are both about capturing enough of "what happened" in real time that you don't need to reproduce the bug to diagnose it.
- **You need to gather diagnostics from a live process.** Go's answer to this is exactly what makes this level "almost entirely genuinely verified, not illustrative": `net/http/pprof` lets you pull a CPU profile, a heap snapshot, or a full goroutine dump from a process **while it keeps serving requests**, over a plain HTTP request. No restart, no attached debugger, no lost state.

The rest of this level walks through each of those tools for real: profiling a program that's actually doing (and allocating, and blocking on) real work, and reading the real output it produces.

---

## CPU Profiling With net/http/pprof

`net/http/pprof` is a standard-library package that, when imported **for its side effect**, registers a handful of `/debug/pprof/*` HTTP handlers that expose the Go runtime's built-in profiler. You never call anything in the package directly:

```go
import (
    "net/http"
    _ "net/http/pprof" // side-effect import: registers /debug/pprof/* handlers
)
```

### The Gotcha: It Registers on http.DefaultServeMux

This is the single most common real-world mistake with `net/http/pprof`, and it was reproduced for real while writing this level: the side-effect import registers its handlers on **`http.DefaultServeMux`** - the package-level default mux `http.HandleFunc` uses. If your server builds its own mux with `http.NewServeMux()` and passes that to `http.Server{Handler: mux}`, the pprof endpoints are never mounted on it and every request to `/debug/pprof/...` returns a plain `404 page not found`, with no hint that pprof was even imported. The fix is either to use `http.DefaultServeMux` (directly or implicitly via `http.HandleFunc`), or to explicitly register the exported handlers (`pprof.Index`, `pprof.Profile`, `pprof.Symbol`, `pprof.Trace`, `pprof.Cmdline`) on your own mux.

### Exploring the Endpoints

Here is the real index page (`GET /debug/pprof/`) served by a minimal Go program in this sandbox:

```
$ curl -s http://127.0.0.1:6060/debug/pprof/
<html>
<head>
<title>/debug/pprof/</title>
...
</head>
<body>
/debug/pprof/
<br>
<p>Set debug=1 as a query parameter to export in legacy text format</p>
<br>
Types of profiles available:
<table>
<thead><td>Count</td><td>Profile</td></thead>
<tr><td>3</td><td><a href='allocs?debug=1'>allocs</a></td></tr>
<tr><td>0</td><td><a href='block?debug=1'>block</a></td></tr>
<tr><td>0</td><td><a href='cmdline?debug=1'>cmdline</a></td></tr>
<tr><td>4</td><td><a href='goroutine?debug=1'>goroutine</a></td></tr>
<tr><td>3</td><td><a href='heap?debug=1'>heap</a></td></tr>
<tr><td>0</td><td><a href='mutex?debug=1'>mutex</a></td></tr>
<tr><td>0</td><td><a href='profile?debug=1'>profile</a></td></tr>
<tr><td>0</td><td><a href='symbol?debug=1'>symbol</a></td></tr>
<tr><td>5</td><td><a href='threadcreate?debug=1'>threadcreate</a></td></tr>
<tr><td>0</td><td><a href='trace?debug=1'>trace</a></td></tr>
</table>
<a href="goroutine?debug=2">full goroutine stack dump</a>
```

That's the real, unedited HTML this project's own verification server returned. The important endpoints for this level:

| Endpoint | What it gives you |
|----------|--------------------|
| `/debug/pprof/profile?seconds=N` | Blocks for N seconds, sampling CPU usage, then returns a binary CPU profile |
| `/debug/pprof/heap` | A snapshot of the current heap: live allocations by call site |
| `/debug/pprof/goroutine` | The current count and stacks of every goroutine, in pprof format |
| `/debug/pprof/goroutine?debug=2` | The same thing, but as a human-readable full stack dump |
| `/debug/pprof/allocs` | Historical allocation profile (all allocations, not just live ones) |

### Capturing a Real CPU Profile

The setup: a small HTTP server with one CPU-bound handler (`/work`, a naive prime counter up to 200,000) and a load generator that fires concurrent requests at it for several seconds. While the load runs, the program itself issues `GET /debug/pprof/profile?seconds=5` - the standard way to capture a CPU profile, whether you do it from `curl`, a browser, or (as here) from Go code:

```go
resp, err := client.Get("http://127.0.0.1:6061/debug/pprof/profile?seconds=5")
// ... write resp.Body to a file, e.g. cpu.pprof
```

That request **blocks for the full 5 seconds** while the Go runtime's profiler samples the process 100 times per second, then streams back a binary profile. Real output from this run:

```
server listening on 127.0.0.1:6061 (loopback only)
generating load against /work for 8s while a profile is captured...
capturing CPU profile: GET /debug/pprof/profile?seconds=5
wrote 6492 bytes to cpu.pprof
load generation finished; shutting down
```

### Reading It With go tool pprof

```bash
go tool pprof -top -nodecount=10 ./server cpu.pprof
```

Real captured output:

```
File: server
Type: cpu
Time: 2026-09-01 23:23:24 IST
Duration: 5.11s, Total samples = 15.84s (309.81%)
Showing nodes accounting for 15.70s, 99.12% of 15.84s total
Dropped 104 nodes (cum <= 0.08s)
Showing top 10 nodes out of 49
      flat  flat%   sum%        cum   cum%
     9.42s 59.47% 59.47%      9.42s 59.47%  syscall.rawsyscalln
     5.57s 35.16% 94.63%      5.72s 36.11%  main.isPrime (inline)
     0.29s  1.83% 96.46%      6.03s 38.07%  main.countPrimesHandler
     0.16s  1.01% 97.47%      0.16s  1.01%  runtime.asyncPreempt
     0.16s  1.01% 98.48%      0.16s  1.01%  runtime.madvise
     0.09s  0.57% 99.05%      0.09s  0.57%  runtime.kevent
     0.01s 0.063% 99.12%      0.17s  1.07%  runtime.stackpoolalloc
         0     0% 99.12%      9.25s 58.40%  bufio.(*Writer).Flush
         0     0% 99.12%      9.22s 58.21%  internal/poll.(*FD).Write
         0     0% 99.12%      9.25s 58.40%  internal/poll.ignoringEINTRIO (inline)
```

Two things worth reading carefully here. First, `Total samples = 15.84s (309.81%)` over a 5.11s wall-clock window is normal, not a bug - four concurrent workers meant roughly 4x CPU-seconds of work happened during those 5 real seconds (`GOMAXPROCS` permitting). Second, on this machine `syscall.rawsyscalln` (the network I/O writing each response) edged out the actual CPU-bound code as the single hottest *flat* frame - a reminder to always look at what's actually consuming time rather than assuming the function you *meant* to stress will automatically dominate. `main.isPrime`, the function actually doing the CPU-bound work, is unambiguously the second-hottest frame and the real subject of this profile - `go tool pprof -top` telling you exactly that, from real samples, is the entire point of the exercise.

---

## Memory Profiling

The heap profile (`/debug/pprof/heap`) reports live allocations grouped by the call site that made them. To see it reflect something predictable, this level's verification program allocates a fixed amount of memory per request (1,000 chunks of 10KB, retained in a package-level slice so garbage collection can't reclaim them before the profile is taken) and calls that endpoint 5 times before capturing:

```go
var retained [][]byte

func allocateHandler(w http.ResponseWriter, r *http.Request) {
    const chunks, chunkSize = 1000, 10*1024
    for i := 0; i < chunks; i++ {
        b := make([]byte, chunkSize)
        for j := range b {
            b[j] = byte(j) // touch it so it's not optimized away
        }
        retained = append(retained, b)
    }
}
```

Real output from calling `/allocate` five times, then capturing with `/debug/pprof/heap?gc=1` (the `gc=1` query parameter forces a garbage collection cycle first, so the profile reflects genuinely live memory, not garbage waiting to be swept):

```
allocated 1000 more chunks of 10240 bytes, total retained chunks: 1000
allocated 1000 more chunks of 10240 bytes, total retained chunks: 2000
allocated 1000 more chunks of 10240 bytes, total retained chunks: 3000
allocated 1000 more chunks of 10240 bytes, total retained chunks: 4000
allocated 1000 more chunks of 10240 bytes, total retained chunks: 5000
wrote 1329 bytes to heap.pprof
```

Analyzing it (`-sample_index=inuse_space` selects "currently live bytes," as opposed to `alloc_space`, which would include memory already freed):

```bash
go tool pprof -top -nodecount=10 -sample_index=inuse_space ./server heap.pprof
```

```
File: server
Type: inuse_space
Time: 2026-09-01 23:24:01 IST
Showing nodes accounting for 52.50MB, 100% of 52.50MB total
Showing top 10 nodes out of 25
      flat  flat%   sum%        cum   cum%
   50.49MB 96.17% 96.17%    50.49MB 96.17%  main.allocateHandler
    2.01MB  3.83%   100%     2.01MB  3.83%  runtime.mallocgc
         0     0%   100%    50.49MB 96.17%  net/http.(*ServeMux).ServeHTTP
         0     0%   100%    50.49MB 96.17%  net/http.(*conn).serve
```

5,000 chunks x 10KB is exactly ~50MB, and `main.allocateHandler` accounts for 50.49MB (96.17%) of the live heap - the profile reflects the predictable allocation almost exactly, with the small remainder (`runtime.mallocgc`, a couple of megabytes) being ordinary runtime bookkeeping. This is the core skill: when you see a function name accounting for a disproportionate share of `inuse_space`, that's your leak (or your legitimately large cache) - read top-down and go straight to the flat% column.

---

## Goroutine Profiling

`/debug/pprof/goroutine` reports every currently running goroutine and, in the `debug=2` text form, its **full stack trace** - exactly what you want when a goroutine count is climbing and you need to know where the extra goroutines are stuck. This connects directly to [Level 20's goroutine leak content](../level-20-goroutines/README.md#goroutine-leaks): a leak isn't abstract here, it's a real, capturable stack.

The verification program launches 20 goroutines that all block receiving from a channel nobody sends on (`blockForever`) - a real, but **bounded and eventually-cleaned-up** leak, so the program terminates:

```go
func blockForever(unblock <-chan struct{}, wg *sync.WaitGroup) {
    defer wg.Done()
    <-unblock // blocks here until the exercise unblocks it
}
```

Real captured output:

```
server listening on 127.0.0.1:6063 (loopback only)
baseline goroutines: 2
goroutines after launching 20 blocked workers: 22
wrote 10844 bytes to goroutine.dump
goroutines after cleanup: 5
```

And a real excerpt from `goroutine.dump` (`GET /debug/pprof/goroutine?debug=2`), showing exactly where those 20 extra goroutines are stuck:

```
goroutine 6 [chan receive]:
main.blockForever(0x0?, 0x0?)
	/private/tmp/level38-verify/ex3-goroutine/main.go:22 +0x44
created by main.main in goroutine 1
	/private/tmp/level38-verify/ex3-goroutine/main.go:46 +0x25c

goroutine 7 [chan receive]:
main.blockForever(0x0?, 0x0?)
	/private/tmp/level38-verify/ex3-goroutine/main.go:22 +0x44
created by main.main in goroutine 1
	/private/tmp/level38-verify/ex3-goroutine/main.go:46 +0x25c
```

`[chan receive]` next to each goroutine tells you exactly what it's blocked on, and the `created by` line tells you exactly where it was launched - in a real incident, that's the line you go fix. The same profile, fetched instead from `/debug/pprof/goroutine` (binary form) and summarized with `go tool pprof -top`:

```
File: server
Type: goroutine
Showing nodes accounting for 25, 96.15% of 26 total
      flat  flat%   sum%        cum   cum%
        24 92.31% 92.31%         24 92.31%  runtime.gopark
         1  3.85% 96.15%          1  3.85%  runtime.goroutineProfileWithLabels
...
         0     0% 96.15%         20 76.92%  main.blockForever
...
         0     0% 96.15%         20 76.92%  runtime.chanrecv
```

`main.blockForever` accounts for 20 of the 26 total goroutines (76.92%), every one of them parked in `runtime.chanrecv` - this is exactly the signature you're looking for when triaging a suspected goroutine leak: one function name accounting for a growing, disproportionate share of all goroutines.

---

## Lightweight Diagnostics: NumGoroutine and ReadMemStats

Full pprof captures are the right tool when you need to see *where* time or memory is going. Sometimes you just need a number, cheaply, without spinning up an HTTP capture - `runtime.NumGoroutine()` and `runtime.ReadMemStats()` do exactly that, in-process, with no separate tooling.

```go
var m runtime.MemStats
runtime.ReadMemStats(&m)
fmt.Println("Alloc:", m.Alloc, "NumGoroutine:", runtime.NumGoroutine())
```

Real captured run - launching 50 goroutines, then allocating and retaining ~20MB, then releasing the goroutines:

```
--- baseline ---
Alloc = 179 KB
TotalAlloc = 179 KB
Sys = 7698 KB
NumGC = 0
NumGoroutine = 1
--- after launching 50 goroutines ---
Alloc = 217 KB
TotalAlloc = 217 KB
Sys = 8210 KB
NumGC = 0
NumGoroutine = 51
--- after allocating ~20MB and forcing GC ---
Alloc = 20291 KB
TotalAlloc = 20398 KB
Sys = 29858 KB
NumGC = 4
NumGoroutine = 51
retained chunks: 2000
--- after releasing goroutines (retained memory is untouched) ---
Alloc = 20293 KB
TotalAlloc = 20400 KB
Sys = 29858 KB
NumGC = 4
NumGoroutine = 1
```

Every field here is a real, live number: `Alloc` (currently live heap bytes) jumps from ~200KB to ~20MB exactly when the retained allocation happens and stays there even after the goroutines are released (because the memory, unlike the goroutines, is still referenced). `NumGoroutine` tracks the 50 launched goroutines precisely and drops back to 1 the instant they're released. `NumGC` counting up to 4 shows the garbage collector actually ran (triggered here by an explicit `runtime.GC()` call plus normal heap-growth triggers). This is cheap enough to expose on a `/healthz` or `/debug/vars`-style endpoint and check on every request, unlike a full profile capture.

---

## Common Production Issues and How to Spot Them

| Issue | Signature | Where to look |
|-------|-----------|----------------|
| **Goroutine leak** | `NumGoroutine()` (or `/debug/pprof/goroutine`'s count) rises steadily over time and never comes back down, even when load drops | A goroutine dump repeatedly shows the same function name accumulating in `[chan receive]`, `[chan send]`, or `[select]` states |
| **Memory leak** | Heap (`Alloc` / `/debug/pprof/heap`) rises and never drops back after a GC, even under steady-state load | Compare two heap profiles over time (`go tool pprof -diff_base`, covered in the bonus challenges) - the function whose share keeps growing is your leak |
| **Deadlock** | The entire process crashes immediately with `fatal error: all goroutines are asleep - deadlock!`, rather than hanging silently | This is Go's runtime detector proving *every* goroutine is permanently blocked - see Level 21's [Deadlocks](../level-21-channels/README.md#deadlocks) and Level 23's [real lock-ordering deadlock](../level-23-mutex-waitgroup-atomic/README.md#a-real-deadlock-lock-ordering-gone-wrong) for the full mechanism. Reproduced fresh for this level: |

```go
func main() {
    ch := make(chan int)
    <-ch // nobody will ever send - deadlock
}
```

```
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan receive]:
main.main()
	/private/tmp/level38-verify/ex10-deadlock/main.go:5 +0x30
exit status 2
```

A deadlock, unlike a leak, is loud and immediate - you don't need pprof for it, you need the goroutine dump the runtime already printed for you in the crash output.

There's one more diagnostic worth knowing even though it's not pprof: sending `SIGQUIT` to a running Go process (`kill -QUIT <pid>` on Unix) makes the runtime print a full stack dump of every goroutine to stderr and then terminate - a last-resort way to see "what was every goroutine doing" from a process that's hung and isn't even responding to an HTTP request for `/debug/pprof/goroutine`. A real, trimmed excerpt from this sandbox:

```
SIGQUIT: quit
PC=0x181761fc4 m=0 sigcode=0

goroutine 1 gp=0x39075b2521e0 m=nil [sleep]:
runtime.gopark(...)
	/Users/macbookpro/.local/go/src/runtime/proc.go:462 +0xbc
time.Sleep(0x12a05f200)
	/Users/macbookpro/.local/go/src/runtime/time.go:363 +0x150
main.main()
	/private/tmp/level38-verify/ex11-sigquit/main.go:10 +0x5c
```

Because `SIGQUIT` terminates the process, this is a diagnostic of last resort (when the process is unresponsive anyway), not something to send to a healthy server.

---

## The Delve Debugger

`dlv` (Delve) is the Go community's standard debugger - real breakpoints, real variable inspection, real stepping, unlike `fmt.Println`-driven debugging. This sandbox had Delve installed (`dlv version` reported `Delve Debugger, Version: 1.27.0`), so everything in this section is real, captured `dlv` output.

### Debugging a Program From the Start

```bash
go build -gcflags="all=-N -l" -o app .   # disable optimizations/inlining for clean debugging
dlv exec ./app
```

That `-gcflags="all=-N -l"` matters: without it, the compiler may inline small functions or optimize away variables, and Delve will warn `debugging optimized function` and sometimes fail to read locals at all. Real session, debugging this small program:

```go
func computeTotal(prices []int) int {
    total := 0
    for _, p := range prices {
        total += p
    }
    return total
}
```

```
(dlv) break main.computeTotal
Breakpoint 1 set at 0x1050125b0 for main.computeTotal() ./main.go:5
(dlv) continue
> [Breakpoint 1] main.computeTotal() ./main.go:5 (hits goroutine(1):1 total:1)
(dlv) print prices
[]int len: 4, cap: 4, [10,25,40,5]
(dlv) next
> main.computeTotal() ./main.go:6
(dlv) next
> main.computeTotal() ./main.go:7
(dlv) print total
0
(dlv) continue
total: 80
Process 32543 has exited with status 0
```

Every value here is real: `prices` shows the actual slice contents, and `total` correctly reads `0` at that point in execution - line 7 is the `for` loop header, reached right after `total := 0` ran but before the loop has added anything. Stepping further and printing again would show the sum accumulating.

### Attaching to an Already-Running Process

This is the scenario that matters most for production: the process is already running (you didn't start it under a debugger), and you need to inspect it *right now*.

```bash
dlv attach <pid>
```

Real session, attaching to a running process by PID:

```
(dlv) break main.computeTotal
Breakpoint 1 set at 0x1001c27ec for main.computeTotal() ./main.go:8
(dlv) continue
> [Breakpoint 1] main.computeTotal() ./main.go:8 (hits goroutine(1):1 total:1)
(dlv) print prices
[]int len: 4, cap: 4, [10,25,40,5]
(dlv) next
> main.computeTotal() ./main.go:9
(dlv) next
> main.computeTotal() ./main.go:10
(dlv) print total
0
(dlv) quit
Would you like to kill the process? [Y/n]
```

`dlv attach` works against any running process it has permission to trace - no restart, no pre-planning required, which is exactly why it's the tool of last resort when a production process is doing something wrong and you need to see its actual state. Delve pauses every goroutine in the process while you're stopped at a breakpoint, so use it deliberately (never against a process serving live traffic you can't afford to pause) and detach (or answer "no" to the kill prompt) rather than terminate the process when you're done, unless termination is actually what you want.

### Other Real Delve Workflows

- `dlv test ./...` - runs a test binary under the debugger, so you can set breakpoints inside the code a specific test exercises instead of `main()`.
- `dlv connect` / `dlv --headless --listen=:PORT` - runs Delve as a headless server an IDE (VS Code, GoLand) can connect to remotely, useful for attaching to a debugger running inside a container or a remote VM.
- `bt` (backtrace), `locals`, `args`, `step`, `stepout` - round out the same breakpoint/inspect/step workflow shown above.

---

## Health Check Endpoints

A health check that unconditionally returns `200 OK` tells you nothing - it answers "is the HTTP listener accepting connections," not "can this instance actually do its job." A real health check reports **actual internal state**: can it reach its dependencies, is it internally healthy (not leaking goroutines), and so on. This bridges [Level 27](../level-27-http-rest-apis/README.md)'s HTTP handlers and [Level 36](../level-36-kubernetes/README.md#readiness-and-liveness-probes)'s readiness/liveness probes, which call exactly this kind of endpoint.

```go
func healthzHandler(w http.ResponseWriter, r *http.Request) {
    status := healthStatus{
        DBConnected:  dbUp.Load(),
        NumGoroutine: runtime.NumGoroutine(),
    }
    healthy := status.DBConnected && status.NumGoroutine < maxHealthyGoroutines
    if healthy {
        status.Status = "ok"
        w.WriteHeader(http.StatusOK)
    } else {
        status.Status = "unhealthy"
        w.WriteHeader(http.StatusServiceUnavailable)
    }
    json.NewEncoder(w).Encode(status)
}
```

Real captured output, flipping a simulated database-connectivity flag to prove the endpoint reflects real state rather than a hardcoded response:

```
[db up] HTTP 200 body={Status:ok DBConnected:true NumGoroutine:6 GoroutineLimit:100}
simulating a database outage...
[db down] HTTP 503 body={Status:unhealthy DBConnected:false NumGoroutine:6 GoroutineLimit:100}
simulating recovery...
[db recovered] HTTP 200 body={Status:ok DBConnected:true NumGoroutine:6 GoroutineLimit:100}
```

The HTTP status code genuinely changes (`200` → `503` → `200`) as the underlying state changes - this is what makes the endpoint useful to an orchestrator like Kubernetes: a `503` here is a real, actionable signal ("stop routing traffic to this instance"), not noise.

---

## Structured Logging for Incident Response

Level 33 covered `log/slog` in depth. This section is about applying it specifically to incidents: what to log at each severity, and how to trace one request across every log line it produces.

**What to log at each level during an incident:**
- `Debug` - fine-grained internal state, off by default in production, turned on temporarily while actively investigating
- `Info` - normal request lifecycle events ("request started," "order placed") - the backbone of reconstructing what a system did
- `Warn` - degraded-but-still-functioning conditions (inventory low, health check failing, goroutine count above threshold) - things worth a human's attention but not yet an outage
- `Error` - a request or operation actually failed - these are the lines you grep for first when triaging

**Correlation IDs** solve the problem of a request passing through multiple layers (middleware, a handler, downstream calls) and generating multiple log lines - without a shared ID, you can't tell whether five error lines are one incident or five:

```go
func withRequestLogger(logger *slog.Logger, next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        reqID := r.Header.Get("X-Request-ID")
        if reqID == "" {
            reqID = newRequestID()
        }
        reqLogger := logger.With("request_id", reqID)
        reqLogger.Info("request started", "method", r.Method, "path", r.URL.Path)
        next(w, r.WithContext(context.WithValue(r.Context(), requestIDKey, reqID)))
        reqLogger.Info("request completed", "path", r.URL.Path)
    }
}
```

Real captured output - one successful request (generates its own ID) and one failing request (carries a client-supplied ID, e.g. propagated from an upstream service), every line tagged with the same `request_id`:

```
=== Successful request (own correlation ID) ===
{"time":"...","level":"INFO","msg":"request started","request_id":"f7cf3950","method":"GET","path":"/checkout"}
{"time":"...","level":"INFO","msg":"validating cart","request_id":"f7cf3950","item_count":3}
{"time":"...","level":"WARN","msg":"inventory low","request_id":"f7cf3950","sku":"SKU-42","remaining":1}
{"time":"...","level":"INFO","msg":"order placed","request_id":"f7cf3950","order_id":"ord-9001"}
{"time":"...","level":"INFO","msg":"request completed","request_id":"f7cf3950","path":"/checkout"}

=== Failing request (client-supplied correlation ID, traced across every log line) ===
{"time":"...","level":"INFO","msg":"request started","request_id":"req-incident-77","method":"GET","path":"/checkout"}
{"time":"...","level":"INFO","msg":"validating cart","request_id":"req-incident-77","item_count":3}
{"time":"...","level":"WARN","msg":"inventory low","request_id":"req-incident-77","sku":"SKU-42","remaining":1}
{"time":"...","level":"ERROR","msg":"payment gateway timeout","request_id":"req-incident-77","gateway":"stripe","err":"context deadline exceeded"}
{"time":"...","level":"INFO","msg":"request completed","request_id":"req-incident-77","path":"/checkout"}
```

(Timestamps trimmed to `"..."` here for readability - the real output has full RFC3339 timestamps on every line.) During a real incident, `grep '"request_id":"req-incident-77"'` across every log line - across every layer, potentially across every service - reconstructs this request's entire life in order, which is the entire value of a correlation ID.

---

## Best Practices

### 1. Never Expose pprof Endpoints Publicly

`/debug/pprof/*` should be bound to loopback only, put behind auth, or exposed only on an internal-only port/network - never on the same public listener that serves customer traffic. Every server in this level bound explicitly to `127.0.0.1`, never `0.0.0.0` or a public interface, for exactly this reason.

### 2. Profile Under Real Load, Not Idle Conditions

A CPU profile of an idle process just shows the runtime's own housekeeping. Every CPU profile in this level was captured while a load generator was actively hammering the CPU-bound endpoint - profile while the problem is actually happening (or while you're deliberately reproducing production-like load), not before or after.

### 3. Keep Historical Profiles and Metrics for Comparison

A single heap snapshot tells you what's live right now; it can't tell you whether that number is growing. Save profiles periodically (or on each deploy) so you can diff today's against last week's with `go tool pprof -diff_base` (see the bonus challenges) and answer "is this actually a leak, or has it always looked like this."

### 4. Disable Optimizations When Debugging With Delve

`go build -gcflags="all=-N -l"` before `dlv debug`/`dlv exec` - without it, inlined functions and optimized-away variables make breakpoints land in surprising places and `print` fail on locals that technically still exist.

### 5. Make Health Checks Check Real Dependencies

A `/healthz` that always returns `200` is worse than no health check at all - it actively lies to whatever is deciding whether to route traffic to this instance. Check the things that actually determine whether this instance can do its job (database connectivity, goroutine count, downstream service reachability).

---

## Common Mistakes

### Mistake 1: Leaving /debug/pprof/ Publicly Exposed

This is a real security issue, not a style nitpick: pprof endpoints can leak source file paths, memory contents (via heap dumps), and let an attacker trigger a 30-second CPU profile capture (a cheap denial-of-service lever) or a full-process CPU trace on demand. If a public-facing server accidentally serves `http.DefaultServeMux` (see the gotcha in Section 2) on its public listener, importing `net/http/pprof` anywhere in the binary is enough to expose it.

### Mistake 2: Building Your Own Mux and Forgetting pprof Won't Follow

As demonstrated in Section 2: `_ "net/http/pprof"` only registers on `http.DefaultServeMux`. A custom `http.NewServeMux()` silently gets none of it - the first sign is usually a `404` on `/debug/pprof/` in production with no explanation, until someone remembers this rule.

### Mistake 3: Profiling Under Idle or Artificial Load

A profile captured while nothing interesting is happening (or against a synthetic load that doesn't resemble real traffic shape) tells you what the runtime does at rest, not what's actually slow for real users. Always profile against the workload you actually care about.

### Mistake 4: Ignoring Rising Goroutine/Memory Trends Until It's an Outage

A goroutine count or heap size that climbs slowly over hours or days looks harmless on any single check - it only becomes obviously a problem once it causes an outage. Section 6's leak signatures (and the goroutine-leak-detector bonus challenge) exist specifically to catch the *trend*, not just a single bad snapshot.

### Mistake 5: Missing Correlation IDs

Without a shared `request_id` (or trace ID) threading through every log line a request produces, an incident's logs are a pile of individually-true statements with no way to tell which ones belong to the same failure. Section 9's middleware pattern is the fix - apply it at the boundary (the first place a request enters your system), before any handler-specific logging happens.

---

## Summary

**Profiling Endpoints (net/http/pprof):**
- `/debug/pprof/profile?seconds=N` - CPU profile, blocks for N seconds
- `/debug/pprof/heap` - live heap snapshot (`?gc=1` forces a GC first)
- `/debug/pprof/goroutine` / `?debug=2` - goroutine count/stacks (pprof format / human-readable)
- Import for the side effect; remember it only registers on `http.DefaultServeMux`

**Lightweight Diagnostics:**
- `runtime.NumGoroutine()` - current goroutine count, no HTTP round trip needed
- `runtime.ReadMemStats(&m)` - `Alloc`, `TotalAlloc`, `Sys`, `NumGC`, and more, in-process

**Issue Signatures:**
- Goroutine leak: `NumGoroutine()` rises and never falls
- Memory leak: `Alloc` rises and a GC doesn't reclaim it
- Deadlock: `fatal error: all goroutines are asleep - deadlock!` - loud and immediate, not a silent hang

**Delve:**
- `dlv exec` / `dlv debug` - debug from process start
- `dlv attach <pid>` - attach to an already-running process, no restart needed
- Build with `-gcflags="all=-N -l"` for clean breakpoints and variable inspection

**Health and Logging:**
- `/healthz` should report real dependency state, not a hardcoded `200`
- Correlation IDs (`request_id`) let you trace one request across every log line it produces

---

## Next Steps

You now understand:
- ✅ Why production debugging requires different tools than local development
- ✅ Capturing and reading real CPU profiles with `net/http/pprof` and `go tool pprof`
- ✅ Capturing and reading real heap profiles to confirm predictable (and unpredictable) allocations
- ✅ Capturing real goroutine dumps and recognizing a leak's signature
- ✅ Lightweight in-process diagnostics with `NumGoroutine`/`ReadMemStats`
- ✅ Spotting goroutine leaks, memory leaks, and deadlocks from their real signatures
- ✅ Using the Delve debugger to attach to and inspect a running process
- ✅ Building health checks that report real internal state
- ✅ Structured incident logging with correlation IDs

**Next level:** Level 39 - System Design
- Designing systems at a higher level than a single service
- Applying everything from observability (this level) to deployment (Levels 35-37) at architectural scale

You can now see inside a running Go process instead of guessing at it! Keep going! 🚀
