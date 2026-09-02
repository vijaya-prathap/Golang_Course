# Level 33: Logging - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

---

## Exercise 1: Recap of the Standard log Package

**Objective:** Use `log.Println`, `log.SetPrefix`, and `log.SetFlags`, and see why the output has no structure

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level33-exercise1
cd ~/projects/level33-exercise1
go mod init level33.example/exercise1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "log"

func main() {
    log.Println("server starting on :8080")

    log.SetPrefix("[MYAPP] ")
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
    log.Println("prefix and flags configured")

    log.Printf("connected users: %d", 12)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
2026/09/01 22:44:24 server starting on :8080
[MYAPP] 2026/09/01 22:44:24 main.go:10: prefix and flags configured
[MYAPP] 2026/09/01 22:44:24 main.go:12: connected users: 12
```

(The date and time will reflect when you actually run it - only the shape and the `main.go:10`/`main.go:12` line numbers are guaranteed.)

**Learning Objectives:**
- ✅ Use `log.Println`/`log.Printf`, `log.SetPrefix`, and `log.SetFlags`
- ✅ See that every line is one opaque string - there's no `user_id=42`-style field to filter on
- ✅ Recognize why this package alone isn't enough for a production service

---

## Exercise 2: Basic slog.Info/Warn/Error With Attributes

**Objective:** Emit structured log lines at different severities with key-value attributes

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level33-exercise2
cd ~/projects/level33-exercise2
go mod init level33.example/exercise2
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "log/slog"
    "os"
)

func main() {
    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

    logger.Debug("cache miss", "key", "user:42")
    logger.Info("user logged in", "user_id", 42, "method", "password")
    logger.Warn("rate limit approaching", "user_id", 42, "requests", 95, "limit", 100)
    logger.Error("payment failed", "user_id", 42, "amount", 49.99, "reason", "card_declined")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
time=2026-09-01T22:42:05.383+05:30 level=INFO msg="user logged in" user_id=42 method=password
time=2026-09-01T22:42:05.384+05:30 level=WARN msg="rate limit approaching" user_id=42 requests=95 limit=100
time=2026-09-01T22:42:05.384+05:30 level=ERROR msg="payment failed" user_id=42 amount=49.99 reason=card_declined
```

(The `time=` value will differ on your machine. The `Debug` line does not appear - see Exercise 4 for why.)

**Learning Objectives:**
- ✅ Create a `*slog.Logger` with `slog.New` and `slog.NewTextHandler`
- ✅ Call `Debug`/`Info`/`Warn`/`Error` with alternating key-value attributes
- ✅ Notice `Debug` is filtered out by the default minimum level

---

## Exercise 3: TextHandler vs JSONHandler Side By Side

**Objective:** Compare the two standard-library handlers on the identical log call

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level33-exercise3
cd ~/projects/level33-exercise3
go mod init level33.example/exercise3
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "log/slog"
    "os"
)

func main() {
    fmt.Println("=== TextHandler ===")
    textLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
    textLogger.Info("order placed", "order_id", "ord-9001", "total", 129.50, "items", 3)

    fmt.Println("\n=== JSONHandler ===")
    jsonLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    jsonLogger.Info("order placed", "order_id", "ord-9001", "total", 129.50, "items", 3)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== TextHandler ===
time=2026-09-01T22:42:23.863+05:30 level=INFO msg="order placed" order_id=ord-9001 total=129.5 items=3

=== JSONHandler ===
{"time":"2026-09-01T22:42:23.864295+05:30","level":"INFO","msg":"order placed","order_id":"ord-9001","total":129.5,"items":3}
```

**Learning Objectives:**
- ✅ Swap handlers without touching any `logger.Info(...)` call sites
- ✅ Compare human-readable text output against machine-parseable JSON output
- ✅ Recognize which format a log aggregator expects versus a terminal

---

## Exercise 4: Filtering by Minimum Log Level

**Objective:** Set a minimum level with `HandlerOptions` and verify lower-severity messages are dropped

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level33-exercise4
cd ~/projects/level33-exercise4
go mod init level33.example/exercise4
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "log/slog"
    "os"
)

func main() {
    opts := &slog.HandlerOptions{Level: slog.LevelWarn}
    logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

    logger.Debug("this is filtered out", "reason", "below minimum level")
    logger.Info("this is also filtered out", "reason", "below minimum level")
    logger.Warn("this one appears", "level_min", "WARN")
    logger.Error("this one appears too", "level_min", "WARN")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
time=2026-09-01T22:42:33.969+05:30 level=WARN msg="this one appears" level_min=WARN
time=2026-09-01T22:42:33.969+05:30 level=ERROR msg="this one appears too" level_min=WARN
```

Only two of the four log calls printed anything - `Debug` and `Info` were silently dropped because the handler's minimum level is `slog.LevelWarn`.

**Learning Objectives:**
- ✅ Set a minimum severity with `slog.HandlerOptions{Level: ...}`
- ✅ Verify, with real output, that lower-severity calls produce no output at all
- ✅ Understand this is how production noise gets turned down without touching call sites

---

## Exercise 5: Grouped Attributes and Typed Helpers

**Objective:** Use `slog.String`/`Int`/`Bool`/`Any` and nest related attributes with `slog.Group`

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level33-exercise5
cd ~/projects/level33-exercise5
go mod init level33.example/exercise5
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "log/slog"
    "os"
    "time"
)

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

    logger.Info("attribute helpers",
        slog.String("username", "alice"),
        slog.Int("attempt", 3),
        slog.Bool("success", true),
        slog.Any("roles", []string{"admin", "editor"}),
    )

    logger.Info("request handled",
        slog.String("method", "POST"),
        slog.String("path", "/orders"),
        slog.Group("http",
            slog.Int("status", 201),
            slog.Duration("latency", 42*time.Millisecond),
        ),
        slog.Group("client",
            slog.String("ip", "203.0.113.7"),
            slog.String("user_agent", "curl/8.4.0"),
        ),
    )
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
{"time":"2026-09-01T22:43:15.691523+05:30","level":"INFO","msg":"attribute helpers","username":"alice","attempt":3,"success":true,"roles":["admin","editor"]}
{"time":"2026-09-01T22:43:15.691999+05:30","level":"INFO","msg":"request handled","method":"POST","path":"/orders","http":{"status":201,"latency":42000000},"client":{"ip":"203.0.113.7","user_agent":"curl/8.4.0"}}
```

**Learning Objectives:**
- ✅ Use `slog.String`/`Int`/`Bool`/`Any` typed attribute constructors
- ✅ Nest related attributes under one key with `slog.Group`
- ✅ Observe that a `time.Duration` serializes as raw nanoseconds in JSON

---

## Exercise 6: A Contextual Logger With slog.With()

**Objective:** Build a child logger that carries a request ID across multiple calls without repeating it

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level33-exercise6
cd ~/projects/level33-exercise6
go mod init level33.example/exercise6
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "log/slog"
    "os"
)

func main() {
    base := slog.New(slog.NewTextHandler(os.Stdout, nil))

    reqLogger := base.With("request_id", "req-7f3a", "user_id", 42)

    reqLogger.Info("request started", "path", "/checkout")
    reqLogger.Info("inventory checked", "sku", "SKU-100", "in_stock", true)
    reqLogger.Info("request completed", "status", 200)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
time=2026-09-01T22:43:28.251+05:30 level=INFO msg="request started" request_id=req-7f3a user_id=42 path=/checkout
time=2026-09-01T22:43:28.252+05:30 level=INFO msg="inventory checked" request_id=req-7f3a user_id=42 sku=SKU-100 in_stock=true
time=2026-09-01T22:43:28.252+05:30 level=INFO msg="request completed" request_id=req-7f3a user_id=42 status=200
```

**Learning Objectives:**
- ✅ Use `logger.With(...)` to build a derived logger carrying baseline attributes
- ✅ Confirm `request_id` and `user_id` appear on every call without repeating them
- ✅ Confirm the original `base` logger is unaffected by `.With()`

---

## Exercise 7: Stashing a Logger in context.Context

**Objective:** Carry a request-scoped logger through nested function calls via `context.Context`

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level33-exercise7
cd ~/projects/level33-exercise7
go mod init level33.example/exercise7
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "context"
    "log/slog"
    "os"
)

type loggerKey struct{}

func withLogger(ctx context.Context, logger *slog.Logger) context.Context {
    return context.WithValue(ctx, loggerKey{}, logger)
}

func loggerFrom(ctx context.Context) *slog.Logger {
    if logger, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok {
        return logger
    }
    return slog.Default()
}

func processOrder(ctx context.Context, orderID string) {
    logger := loggerFrom(ctx)
    logger.InfoContext(ctx, "processing order", "order_id", orderID)
    chargeCard(ctx, orderID)
}

func chargeCard(ctx context.Context, orderID string) {
    logger := loggerFrom(ctx)
    logger.InfoContext(ctx, "card charged", "order_id", orderID, "amount", 59.99)
}

func main() {
    base := slog.New(slog.NewTextHandler(os.Stdout, nil))
    reqLogger := base.With("request_id", "req-a1b2")

    ctx := withLogger(context.Background(), reqLogger)
    processOrder(ctx, "ord-500")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
time=2026-09-01T22:43:28.251+05:30 level=INFO msg="processing order" request_id=req-a1b2 order_id=ord-500
time=2026-09-01T22:43:28.252+05:30 level=INFO msg="card charged" request_id=req-a1b2 order_id=ord-500 amount=59.99
```

`chargeCard` is two calls deep from `main` and never receives `*slog.Logger` as an explicit parameter - it pulls the request-scoped logger straight out of `ctx`.

**Learning Objectives:**
- ✅ Store a `*slog.Logger` in a `context.Context` using a typed, unexported key (the Level 24 pattern)
- ✅ Retrieve it in a function nested two calls deep and log with `InfoContext`
- ✅ Provide a safe fallback (`slog.Default()`) for a context that was never seeded

---

## Exercise 8: A Request-Logging HTTP Middleware (Captured With httptest)

**Objective:** Write and test a middleware that logs method, path, status, and duration for every request

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level33-exercise8
cd ~/projects/level33-exercise8
go mod init level33.example/exercise8
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "bytes"
    "fmt"
    "log/slog"
    "net/http"
    "net/http/httptest"
    "time"
)

type statusRecorder struct {
    http.ResponseWriter
    status int
}

func (r *statusRecorder) WriteHeader(code int) {
    r.status = code
    r.ResponseWriter.WriteHeader(code)
}

func RequestLogger(logger *slog.Logger, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

        next.ServeHTTP(rec, r)

        logger.Info("request handled",
            "method", r.Method,
            "path", r.URL.Path,
            "status", rec.status,
            "duration_ms", time.Since(start).Milliseconds(),
        )
    })
}

func main() {
    var buf bytes.Buffer
    logger := slog.New(slog.NewJSONHandler(&buf, nil))

    mux := http.NewServeMux()
    mux.HandleFunc("/brew", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusTeapot)
    })

    handler := RequestLogger(logger, mux)

    req := httptest.NewRequest(http.MethodGet, "/brew", nil)
    rec := httptest.NewRecorder()
    handler.ServeHTTP(rec, req)

    fmt.Println("=== Response ===")
    fmt.Println("status:", rec.Code)

    fmt.Println("\n=== Captured Log Output (buffer, not stdout) ===")
    fmt.Print(buf.String())
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Response ===
status: 418

=== Captured Log Output (buffer, not stdout) ===
{"time":"2026-09-01T22:43:43.329014+05:30","level":"INFO","msg":"request handled","method":"GET","path":"/brew","status":418,"duration_ms":0}
```

**Learning Objectives:**
- ✅ Write a middleware that wraps an `http.Handler` and logs after calling it
- ✅ Drive it with `httptest.NewRequest`/`httptest.NewRecorder` instead of a real network listener
- ✅ Capture log output into a `bytes.Buffer` for deterministic, assertable results instead of relying on stdout

---

## Exercise 9: The Log-and-Return Anti-Pattern, Demonstrated and Fixed

**Objective:** See the same failure logged twice, then fix it so it's logged exactly once

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level33-exercise9
cd ~/projects/level33-exercise9
go mod init level33.example/exercise9
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "errors"
    "log/slog"
    "os"
)

var errNotFound = errors.New("user not found")

func fetchUserBad(logger *slog.Logger, id int) error {
    if id != 1 {
        logger.Error("fetch failed", "user_id", id, "err", errNotFound.Error())
        return errNotFound
    }
    return nil
}

func handleRequestBad(logger *slog.Logger, id int) {
    if err := fetchUserBad(logger, id); err != nil {
        logger.Error("request failed", "user_id", id, "err", err.Error())
        return
    }
}

func fetchUserGood(id int) error {
    if id != 1 {
        return errNotFound
    }
    return nil
}

func handleRequestGood(logger *slog.Logger, id int) {
    if err := fetchUserGood(id); err != nil {
        logger.Error("request failed", "user_id", id, "err", err.Error())
        return
    }
}

func main() {
    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

    logger.Info("=== Anti-pattern: log-and-return (same failure logged twice) ===")
    handleRequestBad(logger, 99)

    logger.Info("=== Fixed: log once, only at the boundary ===")
    handleRequestGood(logger, 99)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
time=2026-09-01T22:43:54.726+05:30 level=INFO msg="=== Anti-pattern: log-and-return (same failure logged twice) ==="
time=2026-09-01T22:43:54.726+05:30 level=ERROR msg="fetch failed" user_id=99 err="user not found"
time=2026-09-01T22:43:54.726+05:30 level=ERROR msg="request failed" user_id=99 err="user not found"
time=2026-09-01T22:43:54.726+05:30 level=INFO msg="=== Fixed: log once, only at the boundary ==="
time=2026-09-01T22:43:54.726+05:30 level=ERROR msg="request failed" user_id=99 err="user not found"
```

Count the `ERROR` lines for each version: the bad path produces two for one real failure, the fixed path produces exactly one.

**Learning Objectives:**
- ✅ Reproduce the log-and-return anti-pattern from Level 16, this time with `slog`
- ✅ Fix it by having the lower layer only return the error, never log it
- ✅ Confirm with real output that the failure is now logged exactly once, at the boundary

---

## Exercise 10: Comprehensive Practice — A Small HTTP Service With Request-ID Logging

**Objective:** Combine request-ID injection, a context-carried logger, and multi-layer logging into one small service

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level33-exercise10
cd ~/projects/level33-exercise10
go mod init level33.example/exercise10
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "bytes"
    "context"
    "fmt"
    "log/slog"
    "net/http"
    "net/http/httptest"
    "time"
)

type loggerKey struct{}

func withLogger(ctx context.Context, logger *slog.Logger) context.Context {
    return context.WithValue(ctx, loggerKey{}, logger)
}

func loggerFrom(ctx context.Context) *slog.Logger {
    if logger, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok {
        return logger
    }
    return slog.Default()
}

type statusRecorder struct {
    http.ResponseWriter
    status int
}

func (r *statusRecorder) WriteHeader(code int) {
    r.status = code
    r.ResponseWriter.WriteHeader(code)
}

// RequestID injects a unique per-request ID and a request-scoped logger into the context.
func RequestID(base *slog.Logger, next http.Handler) http.Handler {
    var counter int
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        counter++
        id := fmt.Sprintf("req-%03d", counter)
        reqLogger := base.With("request_id", id)
        ctx := withLogger(r.Context(), reqLogger)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// RequestLogging logs method, path, status, and duration for every request.
func RequestLogging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

        next.ServeHTTP(rec, r)

        logger := loggerFrom(r.Context())
        logger.InfoContext(r.Context(), "request handled",
            "method", r.Method,
            "path", r.URL.Path,
            "status", rec.status,
            "duration_ms", time.Since(start).Milliseconds(),
        )
    })
}

// --- Layer 2: service layer ---
func createOrder(ctx context.Context, sku string) string {
    logger := loggerFrom(ctx)
    logger.InfoContext(ctx, "creating order", "sku", sku)
    orderID := chargePayment(ctx, sku)
    logger.InfoContext(ctx, "order created", "order_id", orderID)
    return orderID
}

// --- Layer 3: deeper helper ---
func chargePayment(ctx context.Context, sku string) string {
    logger := loggerFrom(ctx)
    logger.InfoContext(ctx, "charging payment", "sku", sku)
    return "ord-777"
}

func ordersHandler(w http.ResponseWriter, r *http.Request) {
    orderID := createOrder(r.Context(), "SKU-42")
    w.WriteHeader(http.StatusCreated)
    fmt.Fprintf(w, "created %s", orderID)
}

func main() {
    var buf bytes.Buffer
    baseLogger := slog.New(slog.NewJSONHandler(&buf, nil))

    mux := http.NewServeMux()
    mux.HandleFunc("/orders", ordersHandler)

    handler := RequestID(baseLogger, RequestLogging(mux))

    req := httptest.NewRequest(http.MethodPost, "/orders", nil)
    rec := httptest.NewRecorder()
    handler.ServeHTTP(rec, req)

    fmt.Println("=== HTTP Response ===")
    fmt.Println("status:", rec.Code)
    fmt.Println("body:", rec.Body.String())

    fmt.Println("\n=== Captured Log Output (every line carries the same request_id) ===")
    fmt.Print(buf.String())
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== HTTP Response ===
status: 201
body: created ord-777

=== Captured Log Output (every line carries the same request_id) ===
{"time":"2026-09-01T22:45:17.522409+05:30","level":"INFO","msg":"creating order","request_id":"req-001","sku":"SKU-42"}
{"time":"2026-09-01T22:45:17.52266+05:30","level":"INFO","msg":"charging payment","request_id":"req-001","sku":"SKU-42"}
{"time":"2026-09-01T22:45:17.522663+05:30","level":"INFO","msg":"order created","request_id":"req-001","order_id":"ord-777"}
{"time":"2026-09-01T22:45:17.522669+05:30","level":"INFO","msg":"request handled","request_id":"req-001","method":"POST","path":"/orders","status":201,"duration_ms":0}
```

All four log lines - across the middleware and two layers of business logic three calls deep - carry the identical `request_id`, injected exactly once. **Note the middleware order**: `RequestID` must wrap *outside* `RequestLogging` so the logger is stashed in the context before `RequestLogging` ever reads it back out. Reversing the order silently breaks the chain - the outer middleware would read from a context that hasn't been enriched yet.

**Learning Objectives:**
- ✅ Inject a unique request ID and a derived logger via middleware
- ✅ Thread the context-carried logger through a handler, a service layer, and a deeper helper
- ✅ Understand why middleware ordering matters when one middleware depends on context another middleware sets up
- ✅ Verify, with real captured JSON, that every log line for one request shares the same correlation ID

---

## Bonus Challenges

### Challenge 1: A Redacting Custom slog.Handler

Write a `slog.Handler` that wraps another handler and replaces the value of any attribute named `"password"` with `"[REDACTED]"` before passing the record on - a defense-in-depth backstop for Best Practice 1 (never log secrets).

```bash
mkdir -p ~/projects/level33-bonus1
cd ~/projects/level33-bonus1
go mod init level33.example/bonus1
```

**Hints:**
- Implement all four methods of the `slog.Handler` interface: `Enabled`, `Handle`, `WithAttrs`, `WithGroup`
- Inside `Handle`, build a new `slog.Record` with `slog.NewRecord(r.Time, r.Level, r.Message, r.PC)`, then use `r.Attrs(func(a slog.Attr) bool { ... })` to copy each attribute across, swapping the value when the key matches
- Delegate `WithAttrs`/`WithGroup` to the wrapped handler and re-wrap the result, so the redaction survives `.With(...)` calls too

### Challenge 2: TextHandler vs JSONHandler Overhead

Time 100,000 log calls through a `TextHandler` and through a `JSONHandler`, both writing to `io.Discard`, and print the per-call overhead in nanoseconds for each.

```bash
mkdir -p ~/projects/level33-bonus2
cd ~/projects/level33-bonus2
go mod init level33.example/bonus2
```

**Hints:**
- Use `io.Discard` as the writer so you're timing encoding, not I/O
- Wrap each loop in `time.Now()`/`time.Since(start)`
- Print a note that exact numbers vary by machine and load - for reliable numbers, use `go test -bench` instead of a hand-rolled timing loop

### Challenge 3: A Log-Sampling Helper

Write a `Sampler` type that only allows logging 1 in every N occurrences of a noisy, high-frequency event (e.g. a health check hit thousands of times a minute), to address the "excessive debug logging in hot paths" mistake without losing visibility entirely.

```bash
mkdir -p ~/projects/level33-bonus3
cd ~/projects/level33-bonus3
go mod init level33.example/bonus3
```

**Hints:**
- Keep an internal counter, incremented with `sync/atomic` so it's safe if called from multiple goroutines
- `Should()` returns `true` only when the counter modulo N lands on a fixed remainder (e.g. `count%n == 1`), so it fires on the 1st, (N+1)th, (2N+1)th call, and so on
- Log a summary line at the end showing how many of the total events were actually logged, to confirm the sampling ratio

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Recognize what the standard `log` package can't do, and why `log/slog` exists
✅ Emit structured log lines at Debug/Info/Warn/Error with key-value attributes
✅ Compare `TextHandler` and `JSONHandler` output on identical calls
✅ Set and verify a minimum log level that filters out noisy lower-severity messages
✅ Use typed attribute helpers and `slog.Group` to nest related data
✅ Build contextual loggers with `.With()` and carry them through `context.Context`
✅ Write and test a request-logging HTTP middleware with deterministic, buffer-captured output
✅ Identify and fix the log-and-return anti-pattern
✅ Combine request-ID injection, middleware, and multi-layer logging into one small, verifiable HTTP service

---

## Next Level

Level 34: Configuration
- Loading configuration from environment variables and files
- Setting the log level (and everything else) from external config instead of hardcoding it
- Validating configuration at startup instead of failing deep inside a request

Great work! Your programs can now explain themselves in production! 🚀
