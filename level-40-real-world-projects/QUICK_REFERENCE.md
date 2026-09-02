# Level 40: Quick Reference Card

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir -p ~/projects/level40-capstone/urlshortener && cd $_
go mod init level40.example/urlshortener
mkdir -p domain repository service handler middleware ratelimit config

# Build up through the milestones in EXERCISES.md, then run everything:
gofmt -l .
go vet ./...
go build ./...
go test ./... -count=1

# Run the real server:
JWT_SECRET="dev-secret" PORT=8080 go run .
```

---

## 📋 API Endpoints

```
POST /auth/register           { "username": "...", "password": "..." }  -> 201 User
POST /auth/login               { "username": "...", "password": "..." }  -> 200 { "token": "<jwt>" }
POST /shorten   (Bearer <jwt>) { "long_url": "https://..." }             -> 201 Link
GET  /{code}                                                             -> 302 redirect
GET  /stats/{code}                                                       -> 200 Link
GET  /healthz                                                            -> 200 ok
```

Only `POST /shorten` requires auth. Only `GET /{code}` is rate-limited.

---

## 🏗️ The Layers (Level 30 + 31)

```go
// Repository - storage only, an interface + an in-memory implementation
type LinkRepository interface {
    Create(l domain.Link) error
    Get(code string) (domain.Link, error)
    IncrementClicks(code string) error
    CountByOwner(ownerID string) (int, error)
    List() []domain.Link
}

// Service - business rules, depends on the REPOSITORY INTERFACE
type LinkService struct { repo LinkRepository; codeLen int }
func NewLinkService(repo LinkRepository, codeLen int) *LinkService

// Handler - thin HTTP translation, depends on the SERVICE INTERFACE
type LinkHandler struct { service LinkService }
func NewLinkHandler(s LinkService) *LinkHandler

// main.go - the ONLY place that knows the concrete types (Level 31)
func buildServer(cfg config.Config, logger *slog.Logger) http.Handler {
    linkRepo := repository.NewInMemoryLinkRepository()
    linkService := service.NewLinkService(linkRepo, cfg.CodeLength)
    linkHandler := handler.NewLinkHandler(linkService)
    // ...
}
```

---

## 🔐 JWT Auth Middleware (Level 32)

```go
func RequireAuth(parser TokenParser) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
            userID, err := parser.ParseUserID(token)
            if err != nil {
                http.Error(w, "invalid or expired token", http.StatusUnauthorized)
                return
            }
            next.ServeHTTP(w, r.WithContext(handler.WithUserID(r.Context(), userID)))
        })
    }
}
```

```bash
curl -X POST localhost:8080/shorten -H "Authorization: Bearer $TOKEN" -d '{"long_url":"https://example.com"}'
```

---

## 📊 Structured Logging (Level 33)

```go
func Logging(logger *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
            next.ServeHTTP(rec, r)
            logger.Info("http request",
                slog.String("method", r.Method), slog.String("path", r.URL.Path),
                slog.Int("status", rec.status), slog.Duration("latency", time.Since(start)))
        })
    }
}
```

---

## ⚙️ Configuration (Level 34)

```go
type Config struct {
    Port, JWTSecret string
    TokenTTL        time.Duration
    CodeLength      int
    RateLimitPerSec float64
    RateLimitBurst  int
}

func LoadConfig() (Config, error) {
    // os.LookupEnv for anything where "unset" must differ from "empty"
    // FAILS FAST if JWT_SECRET is unset (unless ALLOW_DEV_JWT_SECRET=true)
}
```

```bash
JWT_SECRET=... PORT=8080 RATE_LIMIT_PER_SEC=5 RATE_LIMIT_BURST=10 go run .
```

---

## 🪣 Token-Bucket Rate Limiter (bridges Level 39)

```go
type Limiter struct {
    buckets       map[string]*bucket
    ratePerSecond float64
    burst         int
}

func (l *Limiter) Allow(key string) bool {
    // refill: tokens += elapsed_seconds * ratePerSecond, capped at burst
    // if tokens < 1: deny; else: consume one token, allow
}

func (l *Limiter) Middleware(keyFunc func(*http.Request) string) func(http.Handler) http.Handler
```

```bash
# burst=10 lets 10 requests through immediately, then throttles to 5/sec
```

---

## 🧪 Testing (Level 26)

```go
// Unit test a layer against a FAKE dependency
svc := NewLinkService(newFakeLinkRepository(), 7)

// Unit test a handler against a FAKE service, with stub middleware
mux := NewMux(NewLinkHandler(fakeSvc), NewAuthHandler(fakeAuth), identity, identity, identity)

// Full end-to-end test: the REAL dependency graph, over a real TCP connection
srv := httptest.NewServer(buildServer(cfg, logger))
resp, _ := http.Post(srv.URL+"/shorten", "application/json", body)
```

```bash
go test ./...           # every package
go test . -run TestEndToEnd -v   # just the full-system test
```

---

## ⚠️ Two Real Bugs Worth Remembering

| Symptom | Cause | Fix |
|---|---|---|
| `Port = ""; want default 8080` | `t.Setenv("PORT", "")` sets it PRESENT, not unset | `os.Unsetenv("PORT")` instead |
| Rate limiter never returns 429 in a test | Response body not drained before `Close()` -> no connection reuse -> new `RemoteAddr` every request | `io.Copy(io.Discard, resp.Body)` before `Close()` |

---

## 🐳 Docker / ☸️ Kubernetes (Levels 35-36 honesty pattern)

```bash
docker build -t urlshortener:1.0.0 .      # needs a running Docker daemon
kubectl apply -f k8s/deployment.yaml      # needs a reachable cluster
```

If either tool reports "cannot connect" / "connection refused," that matches this course's own authoring environment exactly - the Dockerfile and manifests are correct per documented syntax, just unverified against a live daemon/cluster here.

---

## 🎓 Before Next Level

Can you:
- [ ] Trace a request through Logging -> RequireAuth -> Handler -> Service -> Repository, out loud?
- [ ] Explain why `LinkService` never imports `net/http`, and `LinkHandler` never imports the repository package?
- [ ] Explain both real bugs this level caught, and why the fix was in the TEST, not the code, for the first one?
- [ ] Run `go test ./...` and get every package passing?

If YES → You're ready for Level 41!

---

## 📚 Next Level

Level 41: Go Interview Preparation
- Turning this exact project into confident, spoken interview answers
- Concurrency, interfaces, error handling, and system design questions
- Practice explaining your own architectural decisions out loud

You built a complete backend service from nothing but what you already knew. 💪
