# Level 33: Quick Reference Card

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir myapp && cd myapp
go mod init myapp

# Create main.go
cat > main.go << 'EOF'
package main

import (
    "log/slog"
    "os"
)

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

    logger.Info("server started", "port", 8080)
    logger.Warn("cache nearly full", "used_pct", 92)
    logger.Error("db connection failed", "err", "timeout")
}
EOF

# Run
go run main.go
```

---

## 📋 Creating a Logger

```go
// Text output (human-readable)
logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

// JSON output (machine-parseable)
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

// Writing into a buffer instead of stdout (deterministic, testable)
var buf bytes.Buffer
logger := slog.New(slog.NewJSONHandler(&buf, nil))
```

---

## 🔊 Logging Calls

```go
logger.Debug("msg", "key", value)   // fine-grained detail
logger.Info("msg", "key", value)    // normal operational events
logger.Warn("msg", "key", value)    // unexpected, but still succeeded
logger.Error("msg", "key", value)   // an operation failed
```

---

## 🚦 Minimum Level Filtering

```go
opts := &slog.HandlerOptions{Level: slog.LevelWarn}
logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

logger.Debug("dropped")   // no output
logger.Info("dropped")    // no output
logger.Warn("printed")    // printed
logger.Error("printed")   // printed
```

| Level | Value |
|-------|-------|
| `slog.LevelDebug` | -4 |
| `slog.LevelInfo` | 0 |
| `slog.LevelWarn` | 4 |
| `slog.LevelError` | 8 |

---

## 🧩 Structured Attributes

```go
logger.Info("msg",
    slog.String("username", "alice"),
    slog.Int("attempt", 3),
    slog.Bool("success", true),
    slog.Any("roles", []string{"admin", "editor"}),
)

// Grouping related attributes
logger.Info("request handled",
    slog.Group("http",
        slog.Int("status", 201),
        slog.Duration("latency", 42*time.Millisecond),
    ),
)
```

---

## 🎒 Contextual Loggers (.With())

```go
base := slog.New(slog.NewTextHandler(os.Stdout, nil))
reqLogger := base.With("request_id", "req-7f3a", "user_id", 42)

reqLogger.Info("request started")   // carries request_id + user_id automatically
reqLogger.Info("request completed")

// base itself is untouched - .With() returns a NEW logger
```

---

## 🗂️ Logger in context.Context

```go
type loggerKey struct{}   // typed, unexported key - Level 24 pattern

func withLogger(ctx context.Context, l *slog.Logger) context.Context {
    return context.WithValue(ctx, loggerKey{}, l)
}

func loggerFrom(ctx context.Context) *slog.Logger {
    if l, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok {
        return l
    }
    return slog.Default()   // safe fallback
}

// deep inside a call chain, with no logger parameter:
func someHandler(ctx context.Context) {
    logger := loggerFrom(ctx)
    logger.InfoContext(ctx, "did something")
}
```

---

## 🌐 Request-Logging Middleware

```go
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
            "method", r.Method, "path", r.URL.Path,
            "status", rec.status,
            "duration_ms", time.Since(start).Milliseconds(),
        )
    })
}
```

---

## ⚖️ SetDefault vs Explicit

```go
// Package-level default - use in main() only
slog.SetDefault(slog.New(handler))
slog.Info("uses the default")

// Explicit instance - prefer this everywhere else (more testable)
logger := slog.New(handler)
logger.Info("uses an explicit logger")
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Log-and-return | log AND return the same error at every layer | return only; log once at the boundary |
| Noisy hot paths | unconditional `Debug` in a tight loop | sample it, or guard behind a check |
| String concatenation | `logger.Info("user=" + u)` | `logger.Info("msg", "user", u)` |
| No minimum level in prod | `HandlerOptions{}` left at default everywhere | set `Level` explicitly from config |
| Logging secrets | `logger.Info("login", "password", pw)` | `logger.Info("login", "password_provided", pw != "")` |

---

## 🎓 Before Next Level

Can you:
- [ ] Create a logger with `TextHandler` and `JSONHandler`, and explain when to use each?
- [ ] Set a minimum log level and verify filtering with real output?
- [ ] Use `slog.Group` and typed attribute helpers?
- [ ] Build a contextual logger with `.With()`?
- [ ] Stash and retrieve a logger from `context.Context`?
- [ ] Write and test a request-logging HTTP middleware?
- [ ] Explain why log-and-return is a mistake?

If YES → You're ready for Level 34!

---

## 📚 Next Level

Level 34: Configuration
- Loading configuration from environment variables and files
- Setting the log level (and everything else) from external config

You've got structured logging down! 💪
