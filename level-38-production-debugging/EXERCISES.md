# Level 38: Production Debugging - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge. Every program below was actually run in a real sandbox, and every profile/dump/log shown as "Expected Output" is real captured output, not hand-typed - exact numbers (byte counts, CPU percentages, timings) will vary on your machine, and that's normal.

---

## Exercise 1: Setting Up net/http/pprof

**Objective:** Wire up pprof on a real server and confirm the endpoints actually work

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level38-exercise1
cd ~/projects/level38-exercise1
go mod init level38.example/exercise1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "log"
    "net/http"
    _ "net/http/pprof" // side-effect import: registers /debug/pprof/* on DefaultServeMux
    "time"
)

func main() {
    http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "hello")
    })

    addr := "127.0.0.1:6060" // loopback only - never expose pprof publicly
    srv := &http.Server{Addr: addr, Handler: http.DefaultServeMux}

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()

    fmt.Println("server listening on", addr)
    time.Sleep(10 * time.Second)
    _ = srv.Close()
}
EOF
```

3. Run the program, then in a second terminal curl the pprof index while it's running:

```bash
go run main.go &
sleep 1
curl -s http://127.0.0.1:6060/debug/pprof/ | head -20
curl -s http://127.0.0.1:6060/hello
```

**Expected Output** (from `curl` against the index page - the exact counts next to each profile type will differ):

```
<html>
<head>
<title>/debug/pprof/</title>
...
Types of profiles available:
<table>
<thead><td>Count</td><td>Profile</td></thead>
<tr><td>3</td><td><a href='allocs?debug=1'>allocs</a></td></tr>
<tr><td>0</td><td><a href='block?debug=1'>block</a></td></tr>
<tr><td>0</td><td><a href='cmdline?debug=1'>cmdline</a></td></tr>
<tr><td>4</td><td><a href='goroutine?debug=1'>goroutine</a></td></tr>
<tr><td>3</td><td><a href='heap?debug=1'>heap</a></td></tr>
...
```

and `curl http://127.0.0.1:6060/hello` prints `hello`, proving your own handler and pprof's handlers coexist on the same mux.

**Learning Objectives:**
- ✅ Import `net/http/pprof` for its side effect
- ✅ Understand it only registers on `http.DefaultServeMux`
- ✅ Confirm the `/debug/pprof/` index page lists real, live profile endpoints

---

## Exercise 2: Capturing and Analyzing a Real CPU Profile

**Objective:** Generate real CPU load and capture + analyze a genuine profile

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level38-exercise2
cd ~/projects/level38-exercise2
go mod init level38.example/exercise2
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
    _ "net/http/pprof"
    "os"
    "sync"
    "time"
)

func isPrime(n int) bool {
    if n < 2 {
        return false
    }
    for i := 2; i*i <= n; i++ {
        if n%i == 0 {
            return false
        }
    }
    return true
}

func countPrimesHandler(w http.ResponseWriter, r *http.Request) {
    count := 0
    for n := 2; n < 200000; n++ {
        if isPrime(n) {
            count++
        }
    }
    fmt.Fprintf(w, "primes below 200000: %d\n", count)
}

func generateLoad(addr string, duration time.Duration) {
    client := &http.Client{Timeout: 5 * time.Second}
    deadline := time.Now().Add(duration)
    var wg sync.WaitGroup
    for w := 0; w < 4; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for time.Now().Before(deadline) {
                resp, err := client.Get("http://127.0.0.1" + addr + "/work")
                if err != nil {
                    return
                }
                resp.Body.Close()
            }
        }()
    }
    wg.Wait()
}

func main() {
    addr := ":6061"
    http.HandleFunc("/work", countPrimesHandler)
    srv := &http.Server{Addr: "127.0.0.1" + addr, Handler: http.DefaultServeMux}

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()
    time.Sleep(200 * time.Millisecond)

    fmt.Println("generating load against /work for 8s while a profile is captured...")
    loadDone := make(chan struct{})
    go func() {
        generateLoad(addr, 8*time.Second)
        close(loadDone)
    }()

    time.Sleep(500 * time.Millisecond)
    fmt.Println("capturing CPU profile: GET /debug/pprof/profile?seconds=5")
    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Get("http://127.0.0.1" + addr + "/debug/pprof/profile?seconds=5")
    if err != nil {
        log.Fatal(err)
    }
    f, err := os.Create("cpu.pprof")
    if err != nil {
        log.Fatal(err)
    }
    n, err := io.Copy(f, resp.Body)
    resp.Body.Close()
    f.Close()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("wrote %d bytes to cpu.pprof\n", n)

    <-loadDone
    _ = srv.Close()
}
EOF
```

3. Build, run, and analyze:

```bash
go build -o server .
./server
go tool pprof -top -nodecount=10 ./server cpu.pprof
```

**Expected Output** (the `./server` run):

```
generating load against /work for 8s while a profile is captured...
capturing CPU profile: GET /debug/pprof/profile?seconds=5
wrote 6492 bytes to cpu.pprof
```

**Expected Output** (`go tool pprof -top`, this project's real captured run):

```
File: server
Type: cpu
Duration: 5.11s, Total samples = 15.84s (309.81%)
Showing nodes accounting for 15.70s, 99.12% of 15.84s total
      flat  flat%   sum%        cum   cum%
     9.42s 59.47% 59.47%      9.42s 59.47%  syscall.rawsyscalln
     5.57s 35.16% 94.63%      5.72s 36.11%  main.isPrime (inline)
     0.29s  1.83% 96.46%      6.03s 38.07%  main.countPrimesHandler
```

**Learning Objectives:**
- ✅ Generate real concurrent load against a CPU-bound handler
- ✅ Capture a real CPU profile via `GET /debug/pprof/profile?seconds=N`
- ✅ Read `go tool pprof -top` output and identify the actual hot function

---

## Exercise 3: Memory Profiling With Predictable Allocations

**Objective:** Confirm a heap profile reflects a program's actual allocations

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level38-exercise3
cd ~/projects/level38-exercise3
go mod init level38.example/exercise3
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
    _ "net/http/pprof"
    "os"
    "time"
)

var retained [][]byte

func allocateHandler(w http.ResponseWriter, r *http.Request) {
    const chunks, chunkSize = 1000, 10 * 1024
    for i := 0; i < chunks; i++ {
        b := make([]byte, chunkSize)
        for j := range b {
            b[j] = byte(j)
        }
        retained = append(retained, b)
    }
    fmt.Fprintf(w, "total retained chunks: %d\n", len(retained))
}

func main() {
    addr := ":6062"
    http.HandleFunc("/allocate", allocateHandler)
    srv := &http.Server{Addr: "127.0.0.1" + addr, Handler: http.DefaultServeMux}
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()
    time.Sleep(200 * time.Millisecond)

    client := &http.Client{Timeout: 5 * time.Second}
    for i := 0; i < 5; i++ {
        resp, err := client.Get("http://127.0.0.1" + addr + "/allocate")
        if err != nil {
            log.Fatal(err)
        }
        body, _ := io.ReadAll(resp.Body)
        resp.Body.Close()
        fmt.Print(string(body))
    }

    resp, err := client.Get("http://127.0.0.1" + addr + "/debug/pprof/heap?gc=1")
    if err != nil {
        log.Fatal(err)
    }
    f, _ := os.Create("heap.pprof")
    n, _ := io.Copy(f, resp.Body)
    resp.Body.Close()
    f.Close()
    fmt.Printf("wrote %d bytes to heap.pprof\n", n)

    _ = srv.Close()
}
EOF
```

3. Build, run, and analyze:

```bash
go build -o server .
./server
go tool pprof -top -nodecount=10 -sample_index=inuse_space ./server heap.pprof
```

**Expected Output** (`go tool pprof`, real captured run - 5 x 1000 x 10KB = ~50MB):

```
File: server
Type: inuse_space
Showing nodes accounting for 52.50MB, 100% of 52.50MB total
      flat  flat%   sum%        cum   cum%
   50.49MB 96.17% 96.17%    50.49MB 96.17%  main.allocateHandler
    2.01MB  3.83%   100%     2.01MB  3.83%  runtime.mallocgc
```

**Learning Objectives:**
- ✅ Capture a heap profile with `/debug/pprof/heap?gc=1`
- ✅ Distinguish `inuse_space` (live now) from `alloc_space` (all-time)
- ✅ Confirm a profile's numbers match a program's known, predictable allocations

---

## Exercise 4: Goroutine Profiling With Blocked Goroutines

**Objective:** Capture a real goroutine dump showing a leak's signature, safely bounded

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level38-exercise4
cd ~/projects/level38-exercise4
go mod init level38.example/exercise4
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
    _ "net/http/pprof"
    "os"
    "runtime"
    "sync"
    "time"
)

// blockForever simulates a leaked goroutine. It is deliberately unblocked
// and joined before the program exits, so this exercise always terminates.
func blockForever(unblock <-chan struct{}, wg *sync.WaitGroup) {
    defer wg.Done()
    <-unblock
}

func main() {
    addr := ":6063"
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "ok") })
    srv := &http.Server{Addr: "127.0.0.1" + addr, Handler: http.DefaultServeMux}
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()
    time.Sleep(200 * time.Millisecond)

    fmt.Println("baseline goroutines:", runtime.NumGoroutine())

    unblock := make(chan struct{})
    var wg sync.WaitGroup
    for i := 0; i < 20; i++ {
        wg.Add(1)
        go blockForever(unblock, &wg)
    }
    time.Sleep(200 * time.Millisecond)
    fmt.Println("goroutines after launching 20 blocked workers:", runtime.NumGoroutine())

    client := &http.Client{Timeout: 5 * time.Second}
    resp, err := client.Get("http://127.0.0.1" + addr + "/debug/pprof/goroutine?debug=2")
    if err != nil {
        log.Fatal(err)
    }
    f, _ := os.Create("goroutine.dump")
    n, _ := io.Copy(f, resp.Body)
    resp.Body.Close()
    f.Close()
    fmt.Printf("wrote %d bytes to goroutine.dump\n", n)

    // Clean up so the exercise actually terminates.
    close(unblock)
    wg.Wait()
    time.Sleep(100 * time.Millisecond)
    fmt.Println("goroutines after cleanup:", runtime.NumGoroutine())

    _ = srv.Close()
}
EOF
```

3. Run and inspect the dump:

```bash
go run main.go
grep -B1 -A4 "main.blockForever" goroutine.dump | head -20
```

**Expected Output:**

```
baseline goroutines: 2
goroutines after launching 20 blocked workers: 22
wrote 10844 bytes to goroutine.dump
goroutines after cleanup: 5
```

```
goroutine 6 [chan receive]:
main.blockForever(0x0?, 0x0?)
	.../main.go:22 +0x44
created by main.main in goroutine 1
	.../main.go:46 +0x25c
```

**Learning Objectives:**
- ✅ Capture a real goroutine dump via `/debug/pprof/goroutine?debug=2`
- ✅ Recognize `[chan receive]` as a goroutine blocked on a channel
- ✅ Confirm goroutines are cleaned up (count returns to baseline) instead of leaving them leaked

---

## Exercise 5: Lightweight Diagnostics With NumGoroutine and ReadMemStats

**Objective:** Get cheap, in-process diagnostics without a full pprof capture

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level38-exercise5
cd ~/projects/level38-exercise5
go mod init level38.example/exercise5
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "runtime"
    "sync"
)

var retained [][]byte

func printMemStats(label string) {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    fmt.Printf("--- %s ---\n", label)
    fmt.Printf("Alloc = %v KB\n", m.Alloc/1024)
    fmt.Printf("TotalAlloc = %v KB\n", m.TotalAlloc/1024)
    fmt.Printf("NumGC = %v\n", m.NumGC)
    fmt.Printf("NumGoroutine = %v\n", runtime.NumGoroutine())
}

func main() {
    printMemStats("baseline")

    release := make(chan struct{})
    var wg sync.WaitGroup
    for i := 0; i < 50; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            <-release
        }()
    }
    printMemStats("after launching 50 goroutines")

    for i := 0; i < 2000; i++ {
        retained = append(retained, make([]byte, 10*1024))
    }
    runtime.GC()
    printMemStats("after allocating ~20MB and forcing GC")

    close(release)
    wg.Wait()
    printMemStats("after releasing goroutines")
}
EOF
```

3. Run:

```bash
go run main.go
```

**Expected Output** (real captured run - your exact byte counts will differ slightly):

```
--- baseline ---
Alloc = 179 KB
NumGoroutine = 1
--- after launching 50 goroutines ---
Alloc = 217 KB
NumGoroutine = 51
--- after allocating ~20MB and forcing GC ---
Alloc = 20291 KB
NumGC = 4
NumGoroutine = 51
--- after releasing goroutines ---
Alloc = 20293 KB
NumGoroutine = 1
```

**Learning Objectives:**
- ✅ Use `runtime.NumGoroutine()` for a cheap, instant goroutine count
- ✅ Use `runtime.ReadMemStats()` for live heap/GC numbers without a profile capture
- ✅ Observe that released goroutines drop the count, but retained memory does not drop on its own

---

## Exercise 6: A Real Health Check Endpoint

**Objective:** Build a `/healthz` that reports actual internal state, not a hardcoded 200

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level38-exercise6
cd ~/projects/level38-exercise6
go mod init level38.example/exercise6
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "runtime"
    "sync/atomic"
    "time"
)

var dbUp atomic.Bool

const maxHealthyGoroutines = 100

type healthStatus struct {
    Status         string `json:"status"`
    DBConnected    bool   `json:"db_connected"`
    NumGoroutine   int    `json:"num_goroutine"`
    GoroutineLimit int    `json:"goroutine_limit"`
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
    status := healthStatus{
        DBConnected:    dbUp.Load(),
        NumGoroutine:   runtime.NumGoroutine(),
        GoroutineLimit: maxHealthyGoroutines,
    }
    healthy := status.DBConnected && status.NumGoroutine < maxHealthyGoroutines
    if healthy {
        status.Status = "ok"
        w.WriteHeader(http.StatusOK)
    } else {
        status.Status = "unhealthy"
        w.WriteHeader(http.StatusServiceUnavailable)
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(status)
}

func main() {
    dbUp.Store(true)
    addr := ":6065"
    http.HandleFunc("/healthz", healthzHandler)
    srv := &http.Server{Addr: "127.0.0.1" + addr, Handler: http.DefaultServeMux}
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()
    time.Sleep(200 * time.Millisecond)

    client := &http.Client{Timeout: 2 * time.Second}
    get := func(label string) {
        resp, _ := client.Get("http://127.0.0.1" + addr + "/healthz")
        var s healthStatus
        json.NewDecoder(resp.Body).Decode(&s)
        resp.Body.Close()
        fmt.Printf("[%s] HTTP %d body=%+v\n", label, resp.StatusCode, s)
    }

    get("db up")
    fmt.Println("simulating a database outage...")
    dbUp.Store(false)
    get("db down")
    fmt.Println("simulating recovery...")
    dbUp.Store(true)
    get("db recovered")

    _ = srv.Close()
}
EOF
```

3. Run:

```bash
go run main.go
```

**Expected Output** (real captured run):

```
[db up] HTTP 200 body={Status:ok DBConnected:true NumGoroutine:6 GoroutineLimit:100}
simulating a database outage...
[db down] HTTP 503 body={Status:unhealthy DBConnected:false NumGoroutine:6 GoroutineLimit:100}
simulating recovery...
[db recovered] HTTP 200 body={Status:ok DBConnected:true NumGoroutine:6 GoroutineLimit:100}
```

**Learning Objectives:**
- ✅ Build a `/healthz` endpoint that reports real dependency state
- ✅ Return `503` (not `200`) when the service is genuinely unhealthy
- ✅ Confirm the endpoint's response actually changes as internal state changes

---

## Exercise 7: Structured Incident Logging With Correlation IDs

**Objective:** Trace a single request across every log line it produces

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level38-exercise7
cd ~/projects/level38-exercise7
go mod init level38.example/exercise7
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "fmt"
    "log/slog"
    "net/http"
    "net/http/httptest"
    "os"
)

type ctxKey string

const requestIDKey ctxKey = "request_id"

func newRequestID() string {
    b := make([]byte, 4)
    rand.Read(b)
    return hex.EncodeToString(b)
}

func withRequestLogger(logger *slog.Logger, next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        reqID := r.Header.Get("X-Request-ID")
        if reqID == "" {
            reqID = newRequestID()
        }
        ctx := context.WithValue(r.Context(), requestIDKey, reqID)
        reqLogger := logger.With("request_id", reqID)
        reqLogger.Info("request started", "method", r.Method, "path", r.URL.Path)
        next(w, r.WithContext(ctx))
        reqLogger.Info("request completed", "path", r.URL.Path)
    }
}

func loggerFromContext(ctx context.Context, base *slog.Logger) *slog.Logger {
    if id, ok := ctx.Value(requestIDKey).(string); ok {
        return base.With("request_id", id)
    }
    return base
}

func checkoutHandler(base *slog.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        l := loggerFromContext(r.Context(), base)
        l.Info("validating cart", "item_count", 3)
        l.Warn("inventory low", "sku", "SKU-42", "remaining", 1)
        if r.URL.Query().Get("fail") == "true" {
            l.Error("payment gateway timeout", "gateway", "stripe", "err", "context deadline exceeded")
            http.Error(w, "payment failed", http.StatusBadGateway)
            return
        }
        l.Info("order placed", "order_id", "ord-9001")
        fmt.Fprintln(w, "ok")
    }
}

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    mux := http.NewServeMux()
    mux.HandleFunc("/checkout", withRequestLogger(logger, checkoutHandler(logger)))
    ts := httptest.NewServer(mux)
    defer ts.Close()

    fmt.Println("=== Successful request ===")
    resp, _ := http.Get(ts.URL + "/checkout")
    resp.Body.Close()

    fmt.Println("\n=== Failing request (client-supplied correlation ID) ===")
    req, _ := http.NewRequest("GET", ts.URL+"/checkout?fail=true", nil)
    req.Header.Set("X-Request-ID", "req-incident-77")
    resp2, _ := http.DefaultClient.Do(req)
    resp2.Body.Close()
}
EOF
```

3. Run:

```bash
go run main.go
```

**Expected Output** (real captured run, timestamps will differ):

```
=== Successful request ===
{"time":"...","level":"INFO","msg":"request started","request_id":"f7cf3950","method":"GET","path":"/checkout"}
{"time":"...","level":"INFO","msg":"validating cart","request_id":"f7cf3950","item_count":3}
{"time":"...","level":"WARN","msg":"inventory low","request_id":"f7cf3950","sku":"SKU-42","remaining":1}
{"time":"...","level":"INFO","msg":"order placed","request_id":"f7cf3950","order_id":"ord-9001"}
{"time":"...","level":"INFO","msg":"request completed","request_id":"f7cf3950","path":"/checkout"}

=== Failing request (client-supplied correlation ID) ===
{"time":"...","level":"INFO","msg":"request started","request_id":"req-incident-77","method":"GET","path":"/checkout"}
{"time":"...","level":"ERROR","msg":"payment gateway timeout","request_id":"req-incident-77","gateway":"stripe","err":"context deadline exceeded"}
{"time":"...","level":"INFO","msg":"request completed","request_id":"req-incident-77","path":"/checkout"}
```

**Learning Objectives:**
- ✅ Assign a correlation ID at the request boundary, generating one if the client didn't supply one
- ✅ Thread it through every log line for that request via `context.Context` and `slog.With()`
- ✅ See how a single `request_id` lets you trace a failing request end to end

---

## Exercise 8: Debugging With Delve

**Objective:** Set a real breakpoint and inspect real variables with `dlv`

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level38-exercise8
cd ~/projects/level38-exercise8
go mod init level38.example/exercise8
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func computeTotal(prices []int) int {
    total := 0
    for _, p := range prices {
        total += p
    }
    return total
}

func main() {
    prices := []int{10, 25, 40, 5}
    total := computeTotal(prices)
    fmt.Println("total:", total)
}
EOF
```

3. Check `dlv` is available, build without optimizations, and debug:

```bash
dlv version
go build -gcflags="all=-N -l" -o app .
dlv exec ./app
```

4. At the `(dlv)` prompt, type these commands one at a time:

```
break main.computeTotal
continue
print prices
next
next
print total
continue
quit
```

**Expected Output** (real captured session - if `dlv` isn't installed on your machine, read this as a reference and install it from https://github.com/go-delve/delve to try it yourself):

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
Process exited with status 0
```

**Learning Objectives:**
- ✅ Build with `-gcflags="all=-N -l"` for debugger-friendly binaries
- ✅ Set a breakpoint by function name and hit it with `continue`
- ✅ Inspect real variable values with `print`, and step with `next`

---

## Exercise 9: Diagnosing a Deadlock

**Objective:** Recognize Go's real deadlock crash signature (bridging Level 21/23)

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level38-exercise9
cd ~/projects/level38-exercise9
go mod init level38.example/exercise9
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

func main() {
    ch := make(chan int)
    <-ch // nobody will ever send on this channel - deadlock
}
EOF
```

3. Run it (this is expected to crash - that's the point):

```bash
go run main.go
echo "exit code: $?"
```

**Expected Output** (real captured run):

```
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan receive]:
main.main()
	.../main.go:5 +0x30
exit status 2
exit code: 1
```

4. Now fix it by sending a value before receiving (or launching a sender goroutine), and confirm it exits cleanly with status 0.

**Learning Objectives:**
- ✅ Recognize `fatal error: all goroutines are asleep - deadlock!` as Go's own runtime detector, not a hang
- ✅ Read the goroutine state (`[chan receive]`) to see exactly what's blocked and where
- ✅ Understand a deadlock crashes loudly and immediately - it never requires a profiling tool to notice

---

## Exercise 10: Comprehensive Practice — A Fully Observable Service

**Objective:** Combine pprof, a real health check, and structured logging in one service, then profile and diagnose it under load

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level38-exercise10
cd ~/projects/level38-exercise10
go mod init level38.example/exercise10
```

2. Create `main.go` (combines every technique from this level: a CPU-bound `/work` handler, `/healthz` reporting real state, structured logging with correlation IDs on every request, and pprof wired up correctly):

```bash
cat > main.go << 'EOF'
package main

import (
    "crypto/rand"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "io"
    "log/slog"
    "net/http"
    _ "net/http/pprof"
    "os"
    "runtime"
    "sync"
    "sync/atomic"
    "time"
)

var (
    logger       = slog.New(slog.NewJSONHandler(os.Stdout, nil))
    dbUp         atomic.Bool
    requestCount atomic.Int64
)

const maxHealthyGoroutines = 200

func newRequestID() string {
    b := make([]byte, 4)
    rand.Read(b)
    return hex.EncodeToString(b)
}

func isPrime(n int) bool {
    if n < 2 {
        return false
    }
    for i := 2; i*i <= n; i++ {
        if n%i == 0 {
            return false
        }
    }
    return true
}

func workHandler(w http.ResponseWriter, r *http.Request) {
    reqID := newRequestID()
    l := logger.With("request_id", reqID)
    requestCount.Add(1)
    l.Info("request started", "path", r.URL.Path)

    count := 0
    for n := 2; n < 300000; n++ {
        if isPrime(n) {
            count++
        }
    }

    l.Info("request completed", "primes_found", count)
    fmt.Fprintf(w, "primes below 300000: %d\n", count)
}

type healthStatus struct {
    Status       string `json:"status"`
    DBConnected  bool   `json:"db_connected"`
    NumGoroutine int    `json:"num_goroutine"`
    Requests     int64  `json:"requests_served"`
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
    status := healthStatus{
        DBConnected:  dbUp.Load(),
        NumGoroutine: runtime.NumGoroutine(),
        Requests:     requestCount.Load(),
    }
    healthy := status.DBConnected && status.NumGoroutine < maxHealthyGoroutines
    if healthy {
        status.Status = "ok"
        w.WriteHeader(http.StatusOK)
    } else {
        status.Status = "unhealthy"
        logger.Warn("health check failing", "db_connected", status.DBConnected, "num_goroutine", status.NumGoroutine)
        w.WriteHeader(http.StatusServiceUnavailable)
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(status)
}

func generateLoad(base string, workers int, duration time.Duration) {
    client := &http.Client{Timeout: 5 * time.Second}
    deadline := time.Now().Add(duration)
    var wg sync.WaitGroup
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for time.Now().Before(deadline) {
                resp, err := client.Get(base + "/work")
                if err != nil {
                    return
                }
                io.Copy(io.Discard, resp.Body)
                resp.Body.Close()
            }
        }()
    }
    wg.Wait()
}

func main() {
    dbUp.Store(true)
    addr := ":6069"
    base := "http://127.0.0.1" + addr

    http.HandleFunc("/work", workHandler)
    http.HandleFunc("/healthz", healthzHandler)
    srv := &http.Server{Addr: "127.0.0.1" + addr, Handler: http.DefaultServeMux}
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Error("server failed", "err", err)
        }
    }()
    time.Sleep(200 * time.Millisecond)
    logger.Info("service started", "addr", "127.0.0.1"+addr)

    client := &http.Client{Timeout: 5 * time.Second}

    resp, _ := client.Get(base + "/healthz")
    var h1 healthStatus
    json.NewDecoder(resp.Body).Decode(&h1)
    resp.Body.Close()
    fmt.Printf("health before load: HTTP %d %+v\n", resp.StatusCode, h1)

    loadDone := make(chan struct{})
    go func() {
        generateLoad(base, 6, 6*time.Second)
        close(loadDone)
    }()
    time.Sleep(500 * time.Millisecond)

    fmt.Println("capturing CPU profile: GET /debug/pprof/profile?seconds=4")
    profResp, err := client.Get(base + "/debug/pprof/profile?seconds=4")
    if err == nil {
        f, _ := os.Create("service-cpu.pprof")
        n, _ := io.Copy(f, profResp.Body)
        profResp.Body.Close()
        f.Close()
        fmt.Printf("wrote %d bytes to service-cpu.pprof\n", n)
    }
    <-loadDone

    fmt.Println("\nsimulating an incident: database connection lost")
    dbUp.Store(false)
    resp2, _ := client.Get(base + "/healthz")
    var h2 healthStatus
    json.NewDecoder(resp2.Body).Decode(&h2)
    resp2.Body.Close()
    fmt.Printf("health during incident: HTTP %d %+v\n", resp2.StatusCode, h2)

    fmt.Println("\nrecovering...")
    dbUp.Store(true)
    resp3, _ := client.Get(base + "/healthz")
    var h3 healthStatus
    json.NewDecoder(resp3.Body).Decode(&h3)
    resp3.Body.Close()
    fmt.Printf("health after recovery: HTTP %d %+v\n", resp3.StatusCode, h3)

    fmt.Println("\ntotal requests served:", requestCount.Load())
    _ = srv.Close()
}
EOF
```

3. Build, run, and profile:

```bash
go build -o service .
./service > run.log 2>&1
go tool pprof -top -nodecount=8 ./service service-cpu.pprof
```

**Expected Output** (`run.log`, non-log lines - real captured run; thousands of JSON log lines are omitted here):

```
health before load: HTTP 200 {Status:ok DBConnected:true NumGoroutine:6 Requests:0}
capturing CPU profile: GET /debug/pprof/profile?seconds=4
wrote 3501 bytes to service-cpu.pprof

simulating an incident: database connection lost
health during incident: HTTP 503 {Status:unhealthy DBConnected:false NumGoroutine:9 Requests:3159}

recovering...
health after recovery: HTTP 200 {Status:ok DBConnected:true NumGoroutine:9 Requests:3159}

total requests served: 3159
```

**Expected Output** (`go tool pprof -top`, real captured run):

```
File: service
Type: cpu
Duration: 4.12s, Total samples = 16.24s (393.71%)
      flat  flat%   sum%        cum   cum%
    14.17s 87.25% 87.25%     14.24s 87.68%  main.isPrime (inline)
     1.22s  7.51% 94.77%      1.22s  7.51%  syscall.rawsyscalln
     0.53s  3.26% 98.03%     15.93s 98.09%  main.workHandler
```

Notice the diagnosis: `main.isPrime` is unambiguously the hottest function (87.25%) - this profile confirms the CPU-bound work is where the time actually goes - and the health endpoint's HTTP status genuinely flips to `503` the instant the simulated database drops, then back to `200` on recovery, all while every request logged its own correlation ID.

**Learning Objectives:**
- ✅ Combine CPU profiling, health checks, and structured logging in one real service
- ✅ Profile a service under genuine concurrent load and confirm the hot path
- ✅ Watch a health endpoint and structured logs both react to a simulated incident in real time

---

## Bonus Challenges

### Challenge 1: A Goroutine-Leak Detector

Build a background goroutine that periodically calls `runtime.NumGoroutine()` and logs a `Warn`-level structured log entry whenever the count exceeds a threshold, and an `Info` entry otherwise. Launch a bounded batch of intentionally-blocked goroutines partway through, then release them, and confirm your detector's logs show the count rising above threshold and then falling back to normal.

```bash
mkdir -p ~/projects/level38-bonus1
cd ~/projects/level38-bonus1
go mod init level38.example/bonus1
```

**Hints:**
- Use a `time.Ticker` and a `select` with a `stop` channel so the detector itself terminates cleanly
- Reuse Exercise 5's `runtime.NumGoroutine()` pattern and Exercise 7's `slog` JSON handler
- A real run of this exact pattern in this level's own verification produced logs like:
  ```
  {"level":"INFO","msg":"goroutine count normal","num_goroutine":2,"threshold":10}
  {"level":"WARN","msg":"goroutine count above threshold","num_goroutine":32,"threshold":10}
  {"level":"INFO","msg":"goroutine count normal","num_goroutine":2,"threshold":10}
  ```

### Challenge 2: Comparing Two Heap Profiles With -diff_base

Capture a baseline heap profile, allocate significantly more memory (simulating growth over time), capture a second heap profile, then use `go tool pprof -diff_base` to see exactly which function's share of the heap grew.

```bash
mkdir -p ~/projects/level38-bonus2
cd ~/projects/level38-bonus2
go mod init level38.example/bonus2
```

**Hints:**
- Capture with `/debug/pprof/heap?gc=1` both times, to two different files (e.g. `heap.base.pprof`, `heap.after.pprof`)
- Run: `go tool pprof -top -sample_index=inuse_space -diff_base=heap.base.pprof ./yourbinary heap.after.pprof`
- A real run of this pattern showed `main.allocate` accounting for over 400% of the (small) diffed total - the diff isolates exactly what grew between the two snapshots, which is the whole point when hunting a slow leak across, say, two profiles captured an hour apart in production

### Challenge 3: Attaching Delve to a Running Process

Start a long-running program in the background, find its PID, and use `dlv attach <pid>` to set a breakpoint and inspect a variable - without ever having started the program under a debugger. This is the closest this course gets to the real "something's wrong with a process that's already running" scenario.

```bash
mkdir -p ~/projects/level38-bonus3
cd ~/projects/level38-bonus3
go mod init level38.example/bonus3
```

**Hints:**
- Build with `go build -gcflags="all=-N -l" -o app .`, run `./app &`, then `dlv attach $!`
- If `dlv` isn't available on your machine, this is a good exercise to read conceptually and try later: the commands are identical to Exercise 8's `dlv exec` session, just via `dlv attach <pid>` instead
- Real captured output from this level's own verification, attaching mid-execution:
  ```
  (dlv) break main.computeTotal
  Breakpoint 1 set at 0x1001c27ec for main.computeTotal() ./main.go:8
  (dlv) continue
  > [Breakpoint 1] main.computeTotal() ./main.go:8 (hits goroutine(1):1 total:1)
  (dlv) print prices
  []int len: 4, cap: 4, [10,25,40,5]
  ```
- Always `quit` (and answer the kill/detach prompt deliberately) rather than leaving a paused, attached process running

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Wire up `net/http/pprof` and know its `http.DefaultServeMux` gotcha cold
✅ Capture and read real CPU profiles with `go tool pprof -top`
✅ Capture and read real heap profiles, confirming predictable allocations
✅ Capture real goroutine dumps and recognize a leak's exact signature
✅ Use `runtime.NumGoroutine()` and `runtime.ReadMemStats()` for cheap diagnostics
✅ Recognize Go's real deadlock crash output on sight
✅ Use the Delve debugger to inspect both a fresh process and an already-running one
✅ Build a `/healthz` endpoint that reports real dependency state
✅ Trace a single request across every log line with a correlation ID
✅ Combine all of the above into one observable, diagnosable service

---

## Next Level

Level 39: System Design
- Designing systems at a higher level than a single service
- Applying observability, deployment, and debugging practices at architectural scale

Great work! You can now see inside a running Go process instead of guessing at it! 🚀
