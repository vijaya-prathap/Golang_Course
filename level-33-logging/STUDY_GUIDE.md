# Level 33: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: log/slog Fundamentals
```
Day 1:  Why structured logging beats fmt.Println; recap of the standard log package
Day 2:  log/slog basics - Debug/Info/Warn/Error with attributes
Day 3:  TextHandler vs JSONHandler
Day 4:  Log levels and HandlerOptions filtering
Day 5:  Typed attributes and slog.Group
Day 6:  Contextual loggers with slog.With()
Day 7:  Stashing a logger in context.Context
```

### Week 2: Middleware & Practice
```
Day 1:  Exercises 1-3 (log recap, basic slog, handlers)
Day 2:  Exercises 4-6 (levels, grouping, With())
Day 3:  Exercises 7-8 (context-carried logger, HTTP middleware)
Day 4:  Exercises 9-10 (log-and-return fix, comprehensive service)
Day 5:  Bonus challenges
Day 6-7: Practice & review
```

---

## 🎯 log vs log/slog

| | `log` | `log/slog` |
|---|-------|-----------|
| **Structure** | One opaque string per line | Key-value attributes on every line |
| **Levels** | None - everything is the same severity | `Debug`/`Info`/`Warn`/`Error`, filterable by minimum level |
| **Output formats** | Plain text only | `TextHandler` (text) or `JSONHandler` (JSON) - swappable |
| **Machine-parseable** | No - requires regex/guessing | Yes - JSON output indexes cleanly |
| **Context-aware calls** | No | `InfoContext`/`WarnContext`/`ErrorContext`/`DebugContext` |
| **Child loggers with baseline fields** | Not built in | `logger.With(...)` |
| **Standard library since** | Always | Go 1.21+ |
| **Still appropriate for** | Quick CLI tools, one-off scripts | Anything that runs as a long-lived service |

```go
// log: one string, no structure
log.Printf("connected users: %d", 12)

// slog: structured key-value data
logger.Info("connection count", "users", 12)
```

---

## 🔁 TextHandler vs JSONHandler - Same Call, Two Outputs

Real captured output from the identical `logger.Info("order placed", "order_id", "ord-9001", "total", 129.50, "items", 3)` call:

```
TextHandler (slog.NewTextHandler):
  time=2026-09-01T22:42:23.863+05:30 level=INFO msg="order placed" order_id=ord-9001 total=129.5 items=3

JSONHandler (slog.NewJSONHandler):
  {"time":"2026-09-01T22:42:23.864295+05:30","level":"INFO","msg":"order placed","order_id":"ord-9001","total":129.5,"items":3}
```

```
Same logger.Info(...) call
        │
        ▼
   *slog.Logger  ──────────────────────────────┐
        │                                       │
        ▼                                       ▼
  TextHandler                             JSONHandler
  (key=value, one line,                   (one JSON object per line,
   human-readable)                         machine-parseable)
        │                                       │
        ▼                                       ▼
  good for: local dev,                    good for: log aggregators
  a human watching a terminal             (Loki, ELK, CloudWatch, Datadog)
```

The logger's calls never change - only the handler passed to `slog.New(...)` does.

---

## 🚦 Log-Level Filtering

`HandlerOptions.Level` sets the **minimum** severity a handler will emit. Anything below it is dropped before it ever reaches the writer - it costs nothing downstream.

```
Severity increases  →

  Debug (-4)  Info (0)  Warn (4)  Error (8)

  HandlerOptions{Level: slog.LevelWarn}
                              │
        ┌─────────────────────┴─────────────────────┐
        │ DROPPED (below minimum)   │  EMITTED (at/above minimum)
        │  Debug                    │   Warn
        │  Info                     │   Error
        └────────────────────────────┴────────────────┘
```

Real captured output with `Level: slog.LevelWarn` - only 2 of 4 calls survive:

```
time=... level=WARN  msg="this one appears"     level_min=WARN
time=... level=ERROR msg="this one appears too" level_min=WARN
```

(The `Debug` and `Info` calls that preceded them produced zero output.)

| Level | Value | Typical use |
|-------|-------|-------------|
| `slog.LevelDebug` | -4 | Fine-grained detail, active debugging only |
| `slog.LevelInfo` | 0 | Normal operational events |
| `slog.LevelWarn` | 4 | Unexpected, but the operation still succeeded |
| `slog.LevelError` | 8 | An operation failed |

---

## 🗂️ The Context-Carried-Logger Pattern

Problem: a request flows through several functions (handler → service → repository, Level 30's layers). Passing `*slog.Logger` as an explicit parameter through every signature works, but `context.Context` (Level 24) is the idiomatic home for request-scoped, cross-cutting data every layer needs.

```
main() / middleware
   │
   │  reqLogger := base.With("request_id", "req-a1b2")
   │  ctx := withLogger(context.Background(), reqLogger)
   ▼
processOrder(ctx, ...)  ── loggerFrom(ctx) ──▶ logs with request_id=req-a1b2
   │
   ▼
chargeCard(ctx, ...)    ── loggerFrom(ctx) ──▶ logs with request_id=req-a1b2
   (2 calls deep - *slog.Logger never appears as a parameter)
```

```go
type loggerKey struct{}                         // typed, unexported - no collisions (Level 24)

func withLogger(ctx context.Context, l *slog.Logger) context.Context {
    return context.WithValue(ctx, loggerKey{}, l)
}

func loggerFrom(ctx context.Context) *slog.Logger {
    if l, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok {
        return l
    }
    return slog.Default()                       // safe fallback, never nil
}
```

Real captured output - `request_id=req-a1b2` on both lines, even though `chargeCard` never takes a logger parameter:

```
time=... level=INFO msg="processing order" request_id=req-a1b2 order_id=ord-500
time=... level=INFO msg="card charged"     request_id=req-a1b2 order_id=ord-500 amount=59.99
```

---

## 🌐 Request-Logging Middleware Flow

```
   incoming HTTP request
          │
          ▼
   ┌──────────────────────────────────────────────┐
   │  RequestID middleware                         │
   │  - generates a request ID                     │
   │  - reqLogger := base.With("request_id", id)   │
   │  - ctx := withLogger(r.Context(), reqLogger)  │
   │  - next.ServeHTTP(w, r.WithContext(ctx))       │
   └──────────────────────────────────────────────┘
          │
          ▼
   ┌──────────────────────────────────────────────┐
   │  RequestLogging middleware                    │
   │  - start := time.Now()                        │
   │  - rec := &statusRecorder{w, http.StatusOK}   │
   │  - next.ServeHTTP(rec, r)  ─────────────┐     │
   │  - logger := loggerFrom(r.Context())    │     │
   │  - logger.InfoContext(..., method,      │     │
   │      path, status, duration_ms)         │     │
   └──────────────────────────────────────────┼────┘
          │                                    │
          ▼                                    ▼
   ┌────────────────────┐          handler + service layers
   │ mux → ordersHandler │ ◀───────  log through the SAME
   │  → createOrder      │           context-carried logger
   │    → chargePayment  │
   └────────────────────┘
```

**Critical ordering rule:** `RequestID` must wrap *outside* `RequestLogging` (`RequestID(base, RequestLogging(mux))`) so the logger is already in the context by the time `RequestLogging` reads it back out after `next.ServeHTTP` returns. Reversing the order means `RequestLogging` reads from the original, unenriched request context - silently breaking the correlation ID chain. Exercise 10 demonstrates the fixed, correct order and confirms all 4 log lines share one `request_id` in real captured JSON output.

---

## 🚨 Common Mistakes

### Mistake 1: Log-and-Return

```go
// ❌ WRONG - same failure logged at every layer
func fetchUser(logger *slog.Logger, id int) error {
    if err != nil {
        logger.Error("fetch failed", "err", err)   // logged here...
        return err
    }
}
func handleRequest(logger *slog.Logger, id int) {
    if err := fetchUser(logger, id); err != nil {
        logger.Error("request failed", "err", err) // ...and logged again here
    }
}

// ✅ RIGHT - return only, log once at the boundary
func fetchUser(id int) error {
    if err != nil {
        return err   // no logging in the lower layer
    }
}
func handleRequest(logger *slog.Logger, id int) {
    if err := fetchUser(id); err != nil {
        logger.Error("request failed", "err", err) // exactly one ERROR line
    }
}
```

Real captured output shows the difference directly: the bad version produces 2 `ERROR` lines for one failure; the fixed version produces exactly 1.

### Mistake 2: Excessive Debug Logging in Hot Paths

```go
// ❌ RISKY - allocates args and checks Enabled() on every iteration of a hot loop
for _, item := range millionItems {
    logger.Debug("processing item", "item", item)
}

// ✅ BETTER - sample it (see Bonus Challenge 3), or guard behind an explicit check
```

### Mistake 3: Unstructured String Concatenation

```go
// ❌ WRONG - defeats the purpose of switching to slog at all
logger.Info("user=" + username + " action=" + action)

// ✅ RIGHT - each value stays its own filterable field
logger.Info("user action", "user", username, "action", action)
```

### Mistake 4: Forgetting to Set a Minimum Level in Production

```go
// ❌ WRONG - nil options default to Info, but leftover Debug calls
// still cost nothing to filter IF this is actually wired up in prod
handler := slog.NewJSONHandler(os.Stdout, nil)

// ✅ RIGHT - read the minimum level from config (Level 34) explicitly
handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: levelFromConfig(),
})
```

---

## 📈 Progression Summary

### Understanding Level 33

Level 33 replaces ad-hoc `fmt.Println` debugging with structured, leveled, queryable logs:

1. **Why it matters** - queryable key-value data, consistent format, filterable severity
2. **log/slog basics** - `Debug`/`Info`/`Warn`/`Error`, key-value attributes
3. **Handlers** - `TextHandler` (human) vs `JSONHandler` (machine), same call sites
4. **Levels** - `HandlerOptions.Level` filters noise for free
5. **Structure** - typed attributes, `slog.Group`, `.With()` for baseline fields
6. **Context** - a request-scoped logger carried through `context.Context` (Level 24)
7. **Middleware** - a real, tested request-logging middleware (Level 27's pattern, upgraded)

### Prerequisites for Level 34

Before moving to Level 34 (Configuration), you need:

- ✅ Comfortable creating `*slog.Logger` instances with `TextHandler`/`JSONHandler`
- ✅ Comfortable setting and reasoning about a minimum log level
- ✅ Can build and use a contextual logger with `.With()`
- ✅ Can stash and retrieve a logger from `context.Context` across nested calls
- ✅ Can write and test a logging HTTP middleware with `httptest`
- ✅ Understand why log-and-return is a mistake, and how to avoid it

### Ready for Level 34?

Level 34 teaches you where that minimum log level - and everything else you've been hardcoding - should actually come from:
- Loading configuration from environment variables and files
- Passing the configured log level into `HandlerOptions.Level` instead of a hardcoded constant
- Validating configuration at startup instead of failing deep inside a request

---

## ✅ Checklist Before Level 34

- [ ] Can explain what `log` lacks that `log/slog` provides
- [ ] Can create a logger with `slog.NewTextHandler` and `slog.NewJSONHandler`
- [ ] Can set a minimum level with `HandlerOptions.Level` and verify filtering with real output
- [ ] Can use `slog.String`/`Int`/`Bool`/`Any` and `slog.Group`
- [ ] Can build a child logger with `.With()` that carries baseline attributes
- [ ] Can stash a logger in `context.Context` with a typed, unexported key and retrieve it in a nested function
- [ ] Can write a request-logging middleware and verify its output via a `bytes.Buffer` and `httptest`
- [ ] Can identify and fix the log-and-return anti-pattern
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The Shift
`fmt.Println` is for a human watching a terminal. `log/slog` is for a machine (and a human) reading logs from a service that's already running somewhere else.

### The Separation
A `*slog.Logger` is a thin wrapper around a `Handler`. Swapping `TextHandler` for `JSONHandler` never touches a single `logger.Info(...)` call site.

### The Filter
`HandlerOptions.Level` drops anything below the minimum severity for free - that's how production turns down noise without redeploying.

### The Carry
`.With()` builds a child logger with baseline attributes for one function's scope; stashing that logger in `context.Context` (Level 24's pattern) carries it through an entire call chain without threading it through every signature by hand.

### The Rule
Log exactly once per failure, at the boundary - never log and return the same error up through every layer it passes.

---

## 📚 Next Level

Level 34: Configuration
- Loading configuration from environment variables and files
- Setting the log level (and everything else) from external config instead of hardcoding it
- Validating configuration at startup instead of failing deep inside a request

You've got structured logging down! Keep going! 🚀
