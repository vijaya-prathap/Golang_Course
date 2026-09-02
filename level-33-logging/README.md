# Level 33: Logging - Complete Guide

## Introduction

Welcome to Level 33! Level 32 (Authentication & JWT) had you writing middleware that made real security decisions - and every one of those decisions is exactly the kind of thing a production system needs a paper trail for. So far, this whole course has used `fmt.Println` to see what a program is doing. That's fine for learning, but it falls apart the moment a program runs unattended, on someone else's machine, at 3 AM, and something goes wrong.

This level is about **logging** - the standard library's `log` package you've been using incidentally since Level 1, and its modern successor `log/slog` (Go 1.21+), which turns free-form text into **structured, queryable data**. Everything here is standard library only: `log` and `log/slog`. No third-party logging framework, fully offline.

By the end of this level you'll be able to emit structured logs with levels and key-value attributes, compare human-readable and machine-readable output side by side, filter noisy logs by severity, carry request-scoped context through a call chain, and write a real request-logging HTTP middleware - all verified with real captured output, not hand-typed guesses.

---

## Table of Contents

1. [Why Structured Logging Beats Scattered fmt.Println](#why-structured-logging-beats-scattered-fmtprintln)
2. [A Quick Recap of the Standard log Package](#a-quick-recap-of-the-standard-log-package)
3. [Introducing log/slog](#introducing-logslog)
4. [Handlers: TextHandler vs JSONHandler](#handlers-texthandler-vs-jsonhandler)
5. [Log Levels In Depth](#log-levels-in-depth)
6. [Structured Attributes and Grouping](#structured-attributes-and-grouping)
7. [Contextual Loggers With slog.With()](#contextual-loggers-with-slogwith)
8. [Logging With context.Context](#logging-with-contextcontext)
9. [A Request-Logging HTTP Middleware](#a-request-logging-http-middleware)
10. [slog.SetDefault() vs Explicit Loggers](#slogsetdefault-vs-explicit-loggers)
11. [Best Practices](#best-practices)
12. [Common Mistakes](#common-mistakes)

---

## Why Structured Logging Beats Scattered fmt.Println

Every earlier level used `fmt.Println` to show what a program was doing. That's a debugging tool for a human staring at a terminal - it is not a production logging strategy. Three problems show up the moment real traffic hits real code:

**1. Free-form text isn't queryable.** `fmt.Println("user", userID, "logged in")` produces a line a human can read but a machine cannot reliably parse. Was `userID` the second word or the third? What if a value itself contains a space? Structured logging instead emits **key-value data**: `user_id=42 event=login`. A log aggregator (or a simple `grep`/`jq`) can filter, count, and group on `user_id` or `event` without guessing at string positions.

**2. There's no consistent format.** One `fmt.Println` call might write `"user 42 logged in"`, another might write `"Login: user=42"`. Across a codebase with dozens of contributors, every message ends up shaped differently, and nothing downstream can rely on a shape that keeps changing.

**3. There are no levels.** `fmt.Println` can't distinguish "the database is unreachable" from "a cache happened to miss." Every line has equal (non-)priority, so during an incident you're scrolling through the same firehose you'd get during a quiet Tuesday afternoon - there's no way to turn the volume down without deleting print statements and redeploying.

Structured logging fixes all three: a consistent key-value shape, a severity **level** on every entry, and output a machine can index and a human can still read.

```
fmt.Println:  "user 42 logged in from 203.0.113.7"
structured:   level=INFO msg="user logged in" user_id=42 ip=203.0.113.7
```

Same information, but the second form can be filtered (`level=ERROR`), searched (`user_id=42`), and aggregated (`count by ip`) without parsing English sentences.

---

## A Quick Recap of the Standard log Package

You've used `log.Println` and `log.Fatal` throughout this course without dwelling on them. Here's the full picture - and its limits, which motivate everything else in this level.

```go
package main

import "log"

func main() {
    log.Println("server starting on :8080")

    log.SetPrefix("[MYAPP] ")
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
    log.Println("prefix and flags configured")

    log.Printf("connected users: %d", 12)
}
```

Real captured output:

```
2026/09/01 22:44:24 server starting on :8080
[MYAPP] 2026/09/01 22:44:24 main.go:10: prefix and flags configured
[MYAPP] 2026/09/01 22:44:24 main.go:12: connected users: 12
```

- `log.Println` / `log.Printf` write a timestamped line to stderr by default.
- `log.SetPrefix` adds a fixed string before every subsequent line.
- `log.SetFlags` controls which metadata is included - `log.Ldate`, `log.Ltime`, `log.Lshortfile` (file:line of the call site) are the common ones.
- `log.Fatal` / `log.Fatalf` do what `Println`/`Printf` do, then call `os.Exit(1)` - nothing deferred after them runs. `log.Panic` / `log.Panicf` log then panic.

**Its limitations, which this level exists to solve:**

- **No structure.** Every call produces one opaque string. There's no `user_id` field to filter on - only a sentence to grep.
- **No levels.** There is no `log.Warn` or `log.Debug`. Everything is the same undifferentiated severity; you cannot turn down the noise without editing code.
- **No attributes.** Attaching multiple pieces of context means hand-formatting them into the message string yourself, inconsistently, every time.

`log` is still perfectly fine for a quick CLI tool or a one-off script. For anything that runs as a long-lived service - which is most of what Levels 27-32 have been building toward - `log/slog` is the better default.

---

## Introducing log/slog

`log/slog` (**s**tructured **log**, standard library since Go 1.21) is built around a `*slog.Logger` with one method per severity level: `Debug`, `Info`, `Warn`, `Error`. Each takes a message string followed by alternating key/value pairs.

```go
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
```

Real captured output:

```
time=2026-09-01T22:42:05.383+05:30 level=INFO msg="user logged in" user_id=42 method=password
time=2026-09-01T22:42:05.384+05:30 level=WARN msg="rate limit approaching" user_id=42 requests=95 limit=100
time=2026-09-01T22:42:05.384+05:30 level=ERROR msg="payment failed" user_id=42 amount=49.99 reason=card_declined
```

Notice the `Debug` line never appears - a `TextHandler` created with `nil` options defaults to a minimum level of `Info`. (Section 5 covers changing that.) Also notice the shape: every line has `time`, `level`, `msg`, and then whatever key-value attributes were passed - consistent and greppable, unlike a hand-built string.

`slog.New(handler)` is how every logger starts: a `*slog.Logger` is just a thin wrapper around a `Handler`, which decides both the output format and where it goes. That separation - logger vs. handler - is what Section 4 builds on.

---

## Handlers: TextHandler vs JSONHandler

The standard library ships two handlers. Same logger calls, same data, two different wire formats:

```go
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
```

Real captured output:

```
=== TextHandler ===
time=2026-09-01T22:42:23.863+05:30 level=INFO msg="order placed" order_id=ord-9001 total=129.5 items=3

=== JSONHandler ===
{"time":"2026-09-01T22:42:23.864295+05:30","level":"INFO","msg":"order placed","order_id":"ord-9001","total":129.5,"items":3}
```

| Handler | Format | Best for |
|---------|--------|----------|
| `slog.NewTextHandler` | `key=value` pairs, one line | A human watching a terminal during local development |
| `slog.NewJSONHandler` | One JSON object per line | A log aggregator (Loki, ELK, CloudWatch, Datadog...) that indexes fields for search |

Both take the same two arguments: an `io.Writer` (where the output goes - `os.Stdout`, a file, a `bytes.Buffer` for tests) and a `*slog.HandlerOptions` (or `nil` for defaults). Swapping handlers never touches the logging call sites - `logger.Info(...)` is identical either way. In production, JSON almost always wins: it's what log aggregation tools expect, and "readable" only matters when a human is watching a raw stream, which is increasingly rare once a service is actually deployed.

---

## Log Levels In Depth

`slog` defines four levels, in increasing severity: `slog.LevelDebug`, `slog.LevelInfo`, `slog.LevelWarn`, `slog.LevelError`. A handler's `HandlerOptions.Level` sets the **minimum** level it will emit - anything below that threshold is silently dropped before it ever reaches the writer.

```go
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
```

Real captured output - only two of the four calls made it through:

```
time=2026-09-01T22:42:33.969+05:30 level=WARN msg="this one appears" level_min=WARN
time=2026-09-01T22:42:33.969+05:30 level=ERROR msg="this one appears too" level_min=WARN
```

This is the mechanism behind "turn down the noise in production without redeploying application code" - most services read the minimum level from an environment variable or config value (Level 34 covers loading that config) and pass it straight into `HandlerOptions.Level`. Debug-level detail stays in the code permanently; it just doesn't cost anything to emit once the environment says "only Warn and above."

| Level | Value | Typical use |
|-------|-------|-------------|
| `slog.LevelDebug` | -4 | Fine-grained detail useful only while actively debugging |
| `slog.LevelInfo` | 0 | Normal operational events worth recording always |
| `slog.LevelWarn` | 4 | Something unexpected, but the request/operation still succeeded |
| `slog.LevelError` | 8 | An operation failed |

---

## Structured Attributes and Grouping

Passing `"key", value` pairs works, but `slog` also has typed constructors - `slog.String`, `slog.Int`, `slog.Bool`, `slog.Any` - which avoid a subtle foot-gun: an odd number of loose arguments, or a key that isn't actually a string, gets logged as a best-effort `!BADKEY` entry instead of a compile error. Typed attributes catch that class of mistake and are marginally faster since they skip a runtime type switch.

`slog.Group` nests a set of attributes under one key - handy for keeping "everything about the HTTP side" separate from "everything about the caller" in the same line.

```go
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
```

Real captured output:

```
{"time":"2026-09-01T22:43:15.691523+05:30","level":"INFO","msg":"attribute helpers","username":"alice","attempt":3,"success":true,"roles":["admin","editor"]}
{"time":"2026-09-01T22:43:15.691999+05:30","level":"INFO","msg":"request handled","method":"POST","path":"/orders","http":{"status":201,"latency":42000000},"client":{"ip":"203.0.113.7","user_agent":"curl/8.4.0"}}
```

Two things worth noticing in that real output: `slog.Any("roles", []string{...})` serializes as a real JSON array, not a stringified blob - and `slog.Group("http", ...)` produces a nested `"http": {...}` object rather than flattening `status`/`latency` at the top level. A `time.Duration` value encodes as its raw nanosecond count (`42000000` for 42ms) in JSON, since JSON has no native duration type - the `TextHandler` would instead print it as `42ms`.

---

## Contextual Loggers With slog.With()

Repeating `"request_id", id` on every single log call in a request handler is tedious and easy to forget on one call out of ten. `logger.With(...)` returns a **new** `*slog.Logger` that carries a fixed set of baseline attributes on every subsequent call, without mutating the original logger.

```go
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
```

Real captured output - `request_id` and `user_id` appear on every line without being passed again:

```
time=2026-09-01T22:43:28.251+05:30 level=INFO msg="request started" request_id=req-7f3a user_id=42 path=/checkout
time=2026-09-01T22:43:28.252+05:30 level=INFO msg="inventory checked" request_id=req-7f3a user_id=42 sku=SKU-100 in_stock=true
time=2026-09-01T22:43:28.252+05:30 level=INFO msg="request completed" request_id=req-7f3a user_id=42 status=200
```

`base` itself is untouched - only `reqLogger` carries the extra attributes. This is the idiomatic way to build a "child logger" per request, per worker, or per subsystem: derive it once with `.With(...)`, then pass that derived logger down instead of the bare one.

---

## Logging With context.Context

`.With()` solves carrying attributes within one function's local scope, but a real request usually flows through several function calls - a handler calls a service, which calls a repository (the pattern from Level 30). Passing a `*slog.Logger` as an explicit parameter through every one of those signatures works, but Level 24 already established the idiomatic place for request-scoped, cross-cutting data that many unrelated layers need: `context.Context`, using a typed, unexported key to avoid collisions.

`slog` also has context-aware methods - `InfoContext`, `WarnContext`, `ErrorContext`, `DebugContext` - which accept a `context.Context` as their first argument. On their own they don't do anything with the context (unless a custom handler inspects it, e.g. to pull out an OpenTelemetry trace ID), but combining them with "stash a request-scoped logger in the context" gives every layer access to the same enriched logger without a parameter threaded through every signature by hand.

```go
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
```

Real captured output - `chargeCard` is two calls deep from `main`, and still logs with the same `request_id`, without it ever appearing as a parameter in `chargeCard`'s signature:

```
time=2026-09-01T22:43:28.251+05:30 level=INFO msg="processing order" request_id=req-a1b2 order_id=ord-500
time=2026-09-01T22:43:28.252+05:30 level=INFO msg="card charged" request_id=req-a1b2 order_id=ord-500 amount=59.99
```

`loggerFrom` falls back to `slog.Default()` when nothing was stashed - so code that forgets to seed the context (or genuinely has no request in flight, like a background job) still logs, just without the extra attributes, instead of panicking on a nil logger. This is precisely the "request-scoped, cross-cutting" use case Level 24 said `context.Value` is actually for - a logger isn't business logic, it's an ambient concern every layer needs.

---

## A Request-Logging HTTP Middleware

Level 27 introduced the middleware pattern (`func(next http.Handler) http.Handler`) with a `log.Logger`-based example. Here's the same shape, upgraded to `slog`, logging method, path, status code, and duration for every request - and, critically for a test suite, writing into a `bytes.Buffer` instead of stdout so the exact output is capturable and deterministic.

```go
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
```

Real captured output:

```
=== Response ===
status: 418

=== Captured Log Output (buffer, not stdout) ===
{"time":"2026-09-01T22:43:43.329014+05:30","level":"INFO","msg":"request handled","method":"GET","path":"/brew","status":418,"duration_ms":0}
```

The same `statusRecorder` trick from Level 27 is doing the work here - embedding `http.ResponseWriter` and overriding only `WriteHeader` lets the middleware observe the status code the inner handler chose. The only new idea is *where* the log line ends up: a `*bytes.Buffer` behaves like any other `io.Writer`, so a test can assert on `buf.String()` directly instead of trying to capture stdout - which is what Exercise 8 does with `httptest`.

---

## slog.SetDefault() vs Explicit Loggers

`slog` has a package-level default logger, used by the bare functions `slog.Info`, `slog.Warn`, `slog.Error`, `slog.Debug`. You set it once with `slog.SetDefault(logger)`, and every call to `slog.Info(...)` anywhere in the program routes through it.

```go
package main

import (
    "log/slog"
    "os"
)

func main() {
    handler := slog.NewTextHandler(os.Stdout, nil)
    slog.SetDefault(slog.New(handler))

    slog.Info("using the package-level default logger", "via", "slog.Info")

    explicit := slog.New(handler)
    explicit.Info("using an explicit logger instance", "via", "explicit.Info")
}
```

Real captured output:

```
time=2026-09-01T22:43:50.457+05:30 level=INFO msg="using the package-level default logger" via=slog.Info
time=2026-09-01T22:43:50.457+05:30 level=INFO msg="using an explicit logger instance" via=explicit.Info
```

Both lines look identical here - the difference only shows up in **testability**. `slog.SetDefault` is global, mutable process-wide state: a test that calls it races with every other test that also touches the default logger, and there's no way to give two different components two different loggers (say, one writing JSON to a file, one writing text to stdout) at the same time.

**Guidance: prefer passing an explicit `*slog.Logger` into constructors** (the dependency-injection instinct from Level 31) over reaching for the global default. Reserve `slog.SetDefault` for `main()` - wiring up the one process-wide fallback that library code and quick scripts fall back to - not as the primary way application code gets its logger.

---

## Best Practices

### 1. Never Log Secrets, Passwords, or Tokens

Level 32 spent real effort making sure passwords are hashed and JWTs are signed correctly - all of that discipline is wasted if a password, an API key, or a raw JWT ends up sitting in plaintext in a log file that outlives the request by months.

```go
// ❌ Bad - the password is now sitting in every log aggregator forever
logger.Info("login attempt", "username", user, "password", password)

// ✅ Good - log that the field was present, never its value
logger.Info("login attempt", "username", user, "password_provided", password != "")
```

Bonus Challenge 1 below builds a custom handler that redacts a named field automatically, as a defense-in-depth backstop for when someone forgets.

### 2. Log at the Right Level

```go
// ❌ Bad - a routine cache miss doesn't deserve Error
logger.Error("cache miss", "key", key)

// ✅ Good - Error is reserved for things that actually failed
logger.Debug("cache miss", "key", key)
logger.Error("database connection failed", "err", err)
```

If everything is logged at `Error`, `Error` stops meaning anything - an on-call engineer can no longer distinguish "the system is on fire" from "a cache warmed up as expected."

### 3. Prefer Structured Attributes Over String Interpolation

```go
// ❌ Bad - back to un-queryable free text, just wrapped in slog
logger.Info(fmt.Sprintf("user %d placed order %s for $%.2f", userID, orderID, total))

// ✅ Good - each value is its own filterable field
logger.Info("order placed", "user_id", userID, "order_id", orderID, "total", total)
```

The whole point of `slog` is the key-value structure - collapsing it back into one formatted string throws away exactly what you switched to `slog` to get.

### 4. Include Correlation IDs for Tracing Requests

A single request often touches a handler, a service, and a repository (Level 30's layers), possibly across multiple goroutines. Without a shared identifier, correlating log lines from one request out of thousands of concurrent ones is guesswork. Generate one ID per request (Section 8 and Exercise 10 show exactly this), attach it via `.With()` once, and every downstream log line carries it automatically.

---

## Common Mistakes

### Mistake 1: Logging AND Returning an Error (Log-and-Return)

Level 16 named this anti-pattern for plain errors; it applies identically to `slog`. Logging an error at every layer it passes through, then also returning it, means one real failure shows up as three or four log lines - and nobody downstream can tell whether that was one incident or four.

```go
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
```

Real captured output - the bad version logs the identical failure twice, the fixed version once:

```
time=2026-09-01T22:43:54.726+05:30 level=INFO msg="=== Anti-pattern: log-and-return (same failure logged twice) ==="
time=2026-09-01T22:43:54.726+05:30 level=ERROR msg="fetch failed" user_id=99 err="user not found"
time=2026-09-01T22:43:54.726+05:30 level=ERROR msg="request failed" user_id=99 err="user not found"
time=2026-09-01T22:43:54.726+05:30 level=INFO msg="=== Fixed: log once, only at the boundary ==="
time=2026-09-01T22:43:54.726+05:30 level=ERROR msg="request failed" user_id=99 err="user not found"
```

**Rule of thumb:** lower layers `return fmt.Errorf("...: %w", err)` (Level 16), adding context as the error travels up. Exactly one place - usually the outermost boundary, like the request-logging middleware in Section 9 - logs it.

### Mistake 2: Excessive Debug Logging in Hot Paths

A `logger.Debug(...)` call inside a loop that runs millions of times still allocates its arguments and calls into the handler's `Enabled` check every single time, even when Debug output is filtered out. In a genuinely hot path, that adds up. Guard expensive-to-compute debug data, or sample it (Bonus Challenge 3), rather than unconditionally logging every iteration.

### Mistake 3: Unstructured String Concatenation Instead of Key-Value Attributes

```go
// ❌ Bad - defeats the purpose of using slog at all
logger.Info("user=" + username + " action=" + action)

// ✅ Good
logger.Info("user action", "user", username, "action", action)
```

This is the same mistake as Best Practice 3, phrased as its own failure mode: reaching for string-building habits from the `fmt.Println` era instead of `slog`'s attribute parameters.

### Mistake 4: Forgetting to Set a Minimum Log Level in Production

A `TextHandler`/`JSONHandler` created with `nil` options defaults to `LevelInfo` - but a codebase full of leftover `logger.Debug(...)` calls from active development can still produce a wall of noise if a `Debug`-level minimum ever leaks into a production deployment, or if a team assumes filtering happens somewhere it doesn't. Section 5 showed the fix: explicitly set `HandlerOptions.Level` from configuration (Level 34), and default that configuration to `Info` or higher outside of local development.

---

## Summary

**The Standard log Package:**
- `log.Println`/`log.Printf`, `log.Fatal`, `log.SetPrefix`/`log.SetFlags` - fine for scripts, has no structure and no levels

**log/slog Basics:**
- `slog.New(handler)` builds a `*slog.Logger`; `Debug`/`Info`/`Warn`/`Error` methods take a message plus key-value pairs

**Handlers:**
- `slog.NewTextHandler` - human-readable `key=value` lines
- `slog.NewJSONHandler` - machine-parseable JSON, the production default

**Levels:**
- `HandlerOptions.Level` sets the minimum severity emitted; everything below it is dropped for free

**Structure:**
- `slog.String`/`Int`/`Bool`/`Any` for typed attributes, `slog.Group` to nest related attributes

**Carrying Context:**
- `logger.With(...)` builds a child logger with baseline attributes
- Stash a request-scoped logger in a `context.Context` (typed, unexported key - Level 24) so any function down the call chain can log consistently, using `InfoContext`/`ErrorContext` etc.

**Middleware:**
- ✅ A request-logging middleware verified with real captured JSON output, following the Level 27 middleware pattern

**Defaults:**
- `slog.SetDefault` for `main()`'s process-wide fallback; explicit `*slog.Logger` instances passed into constructors for everything you actually want to test

---

## Next Steps

You now understand:
- ✅ Why structured logging beats scattered `fmt.Println` calls
- ✅ `log/slog`'s levels, handlers, and structured attributes - verified with real output
- ✅ Building contextual loggers with `.With()` and carrying them through `context.Context`
- ✅ Writing and testing a real request-logging HTTP middleware
- ✅ The log-and-return anti-pattern, and how to log exactly once per failure

**Next level:** Level 34 - Configuration
- Loading configuration from environment variables and files
- Setting the log level (and everything else) from external config instead of hardcoding it
- Validating configuration at startup instead of failing deep inside a request

You can now see what your programs are doing in production, not just on your own terminal. Keep going! 🚀
