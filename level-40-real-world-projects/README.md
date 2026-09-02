# Level 40: Real-World Projects - Complete Guide

## Introduction

Welcome to Level 40 - the capstone. Every level from 1 through 39 taught you one skill at a time, in isolation: loops, structs, interfaces, goroutines, JSON, HTTP handlers, SQL, JWTs, structured logging, configuration, Docker, Kubernetes, system design. Each of those levels was deliberately narrow, because learning ten things at once teaches none of them well.

Real engineering doesn't work in isolated slices. A real backend service needs auth AND persistence AND validation AND logging AND rate limiting AND tests AND a way to ship it - all at once, all correctly, all talking to each other through clean boundaries. **This level is the payoff**: you're going to build one complete, working, tested service that genuinely uses the patterns from Level 30 (repository-service-handler), Level 31 (dependency injection), Level 32 (JWT auth), Level 33 (structured logging), Level 34 (configuration), Level 26 (testing), and bridges Level 39's rate-limiting coverage and Levels 35-36's containerization/orchestration - combined into one coherent system instead of ten separate toy programs.

Nothing here is a new Go feature. Every syntax element, every standard library package, every third-party dependency was already introduced in an earlier level. What's new is **composition**: making all of it work together, under real HTTP traffic, with real tests proving it.

---

## Table of Contents

1. [Why This Level Exists](#why-this-level-exists)
2. [Project Options Overview](#project-options-overview)
3. [Requirements and API Design](#requirements-and-api-design)
4. [Architecture](#architecture)
5. [Building the Repository Layer](#building-the-repository-layer)
6. [Building the Service Layer](#building-the-service-layer)
7. [Building the Handler Layer](#building-the-handler-layer)
8. [Adding Authentication](#adding-authentication)
9. [Adding Structured Logging and Configuration](#adding-structured-logging-and-configuration)
10. [Adding a Rate Limiter](#adding-a-rate-limiter)
11. [Testing the Whole Thing](#testing-the-whole-thing)
12. [Containerizing It and Sketching Kubernetes](#containerizing-it-and-sketching-kubernetes)
13. [What's Next](#whats-next)
14. [Best Practices](#best-practices)
15. [Common Mistakes](#common-mistakes)

---

## Why This Level Exists

Picture the levels you've already completed as a toolbox: a hammer (HTTP handlers), a saw (SQL), a level (JWT auth), a tape measure (structured logging). Every one of those tools works. None of them, alone, builds a house.

A junior engineer who has only ever practiced tools in isolation freezes the first time they're handed a real ticket: "add a URL shortener endpoint with rate limiting and auth." Where does validation go? Does the handler touch the database directly? Where does the JWT get checked? What happens when two concerns - like "is this URL valid" and "has this user hit their limit" - both need to reject the same request?

This level answers those questions by making you build the whole thing, in order, the way a real project actually gets built: design the API first, build the storage layer, build the business rules on top of it, expose it over HTTP, lock it down, make it observable, protect it from abuse, prove it works with real tests, and package it to ship. By the end, "combine everything I know into one working system" stops being an abstract idea and becomes something you've physically done once, with your own hands, with real `go test` output proving every step actually worked.

---

## Project Options Overview

Before committing to a single build, it's worth seeing that this pattern - repository, service, handler, auth, tests - generalizes across many small backend projects, not just one. Three realistic capstone options, all roughly the same size and shape:

| Project | Core Entities | What It Exercises |
|---|---|---|
| **URL Shortener API** | `Link` (code → long URL), `User` | Auth, persistence, a generated identifier, a public high-traffic endpoint worth rate-limiting, simple but real per-user business rules |
| **Todo/Task-List API** | `Task`, `User` | Auth, persistence, CRUD across multiple related resources, ownership checks (a user can only see their own tasks) |
| **Minimal Blog/Notes API** | `Post`, `User` | Auth, persistence, public read endpoints vs. authenticated write endpoints, slightly richer validation (title/body/publish state) |

All three fit naturally into Repository-Service-Handler, all three need JWT auth on the write path, and all three benefit from structured logging and configuration. Pick whichever domain interests you more for your own practice - the *shape* of the solution barely changes.

**This level commits to the URL Shortener** for the full guided build, because it's the smallest domain that still touches every capstone requirement meaningfully:

- **Auth** - creating a short link should require a logged-in user; resolving one should not.
- **Persistence** - short codes and their targets need to survive more than one request.
- **A public, high-traffic endpoint** - the redirect endpoint (`GET /{code}`) is the one part of the system the entire internet can hit, which makes it the natural, honest place to demonstrate rate limiting instead of bolting it on artificially.
- **A small but real business rule** - a per-user limit on how many links someone can create - that has to live in the Service layer, not the Handler or the Repository.

A Todo API or a Notes API would work just as well structurally; a URL shortener simply makes the case for *why* rate limiting matters more concretely than either does.

---

## Requirements and API Design

Real projects define the contract before writing a single handler. Here is the URL shortener's complete API surface, decided up front:

| Method | Path | Auth Required? | Purpose |
|---|---|---|---|
| `POST` | `/auth/register` | No | Create a new user account |
| `POST` | `/auth/login` | No | Exchange username/password for a signed JWT |
| `POST` | `/shorten` | **Yes** (Bearer JWT) | Create a short link for a URL, owned by the caller |
| `GET` | `/{code}` | No (rate-limited) | Redirect to the link's original URL; increments its click count |
| `GET` | `/stats/{code}` | No | Return a link's metadata and click count, WITHOUT incrementing it |
| `GET` | `/healthz` | No | Liveness check |

### Data Model

```go
// User is a registered account. PasswordHash is a bcrypt hash - the
// plaintext password is never stored (Level 32).
type User struct {
    ID           string
    Username     string
    PasswordHash string
    CreatedAt    time.Time
}

// Link is a shortened URL: a generated Code that redirects to a LongURL,
// owned by the user who created it, with a running click count.
type Link struct {
    Code      string
    LongURL   string
    OwnerID   string
    Clicks    int
    CreatedAt time.Time
}
```

### Request/Response Shapes

```
POST /auth/register           { "username": "...", "password": "..." }  -> 201 User
POST /auth/login               { "username": "...", "password": "..." }  -> 200 { "token": "<jwt>" }
POST /shorten   (Bearer <jwt>) { "long_url": "https://..." }             -> 201 Link
GET  /{code}                                                             -> 302 redirect
GET  /stats/{code}                                                       -> 200 Link
```

### Business Rules Decided Up Front

- A `long_url` must be an absolute `http://` or `https://` URL - garbage in never reaches storage.
- A user may have at most **5** active links (`MaxLinksPerUser`) - a deliberately small number chosen so the limit is easy to exercise in a test, standing in for whatever a real product's plan-based quota would be.
- Resolving a code (`GET /{code}`) increments its click counter; reading stats (`GET /stats/{code}`) never does - two different service methods, on purpose, so a monitoring dashboard polling stats can't accidentally inflate click counts.
- The same generic "invalid username or password" error covers both "no such user" and "wrong password" at login, so a client can never enumerate valid usernames (Level 32's rule).

Deciding all of this **before** writing a repository or a handler is the whole point of this section - Section 14 revisits "design first" as a named Best Practice, and Section 15 shows what skipping it costs.

---

## Architecture

The system is Level 30's Repository-Service-Handler layering, wired together with Level 31's constructor injection, with two cross-cutting middleware concerns (auth, rate limiting) and one supporting concern (logging) wrapped around the HTTP layer:

```
                          ┌─────────────────────────────────────────┐
                          │                main()                    │
                          │     (composition root - Level 31)        │
                          └───────────────────┬───────────────────────┘
                                              │ constructs & wires
              ┌───────────────────────────────┼───────────────────────────────┐
              ▼                               ▼                               ▼
      ┌───────────────┐             ┌──────────────────┐            ┌──────────────────┐
      │  Repository    │             │     Service       │            │     Handler       │
      │ (Level 30)     │◄────────────│    (Level 30)      │◄───────────│    (Level 30)      │
      │                │  interface  │                    │  interface │                    │
      │ InMemoryLink-  │             │  LinkService:      │            │  LinkHandler,      │
      │ Repository,    │             │  validate URL,     │            │  AuthHandler       │
      │ InMemoryUser-  │             │  generate code,     │            │  (net/http         │
      │ Repository     │             │  enforce per-user   │            │   ServeMux,         │
      │ (map + mutex,  │             │  limit              │            │   Level 27)         │
      │  Level 23)     │             │                    │            │                    │
      │                │             │  UserService:       │            │                    │
      │                │             │  bcrypt hash,        │            │                    │
      │                │             │  issue/parse JWT     │            │                    │
      │                │             │  (Level 32)          │            │                    │
      └───────────────┘             └──────────────────┘            └──────────────────┘
                                                                              ▲
                          ┌───────────────────────────────────────────────────┤
                          │              Middleware (wraps the Handler)        │
                          │  Logging (Level 33) - every request                │
                          │  RequireAuth (Level 32) - POST /shorten only       │
                          │  Rate Limiter (bridges Level 39) - GET /{code} only│
                          └───────────────────────────────────────────────────┘
```

### Package Layout

```
urlshortener/
├── go.mod
├── main.go                    # composition root: builds and wires everything
├── domain/
│   ├── link.go                # Link struct - no behavior, no dependencies
│   └── user.go                # User struct - no behavior, no dependencies
├── repository/
│   ├── link_repository.go     # LinkRepository interface + InMemoryLinkRepository
│   ├── link_repository_test.go
│   ├── user_repository.go     # UserRepository interface + InMemoryUserRepository
│   └── user_repository_test.go
├── service/
│   ├── shortcode.go           # GenerateCode - crypto/rand short code generation
│   ├── shortcode_test.go
│   ├── link_service.go        # LinkService: validation, per-user limit, resolve/stats
│   ├── link_service_test.go
│   ├── user_service.go        # UserService: bcrypt hashing, JWT issue/parse
│   └── user_service_test.go
├── handler/
│   ├── context.go             # typed context key for the authenticated user ID
│   ├── link_handler.go        # POST /shorten, GET /{code}, GET /stats/{code}
│   ├── auth_handler.go        # POST /auth/register, POST /auth/login
│   ├── routes.go              # NewMux - wires routes to handlers + middleware
│   └── handler_test.go
├── middleware/
│   ├── auth.go                # RequireAuth - JWT Bearer token verification
│   ├── auth_test.go
│   ├── logging.go             # Logging - one structured log line per request
│   └── logging_test.go
├── ratelimit/
│   ├── tokenbucket.go         # Limiter - token-bucket rate limiter + middleware
│   └── tokenbucket_test.go
├── config/
│   ├── config.go              # LoadConfig - env-driven Config struct
│   └── config_test.go
├── main_test.go                # full end-to-end httptest.NewServer test
├── Dockerfile
└── k8s/
    ├── deployment.yaml
    └── service.yaml
```

Every arrow in the diagram points the same direction Level 30 established: **Handler → Service → Repository (interface)**. The Repository implementation, the JWT library, and `net/http` itself sit at the edges - nothing in the middle of the graph (the Service layer) imports any of them directly except through the interfaces it defines for itself.

---

## Building the Repository Layer

The Repository layer's whole job, per Level 30, is storage - nothing else. This capstone uses **in-memory storage** (a map guarded by a `sync.Mutex`, exactly Level 23's and Level 30's pattern) rather than Level 29's SQLite setup. That choice is deliberate and explained in full in Section 5 of `EXERCISES.md` - short version: network access for `modernc.org/sqlite` was verified working in this environment, but Level 30 itself makes the case that the *pattern's whole value* is that storage is swappable behind an interface, and keeping the capstone's core 100% offline (matching every level except 28, 29, and 32) means anyone can complete this level with zero network dependency for the storage layer, and can swap in Level 29's SQLite approach later as a pure drop-in.

```go
// LinkRepository is the Repository layer's contract for Links. Any type
// with these five methods satisfies it implicitly (Level 15) - this
// in-memory map today, a SQL-backed store (Level 29's pattern) tomorrow.
type LinkRepository interface {
    Create(l domain.Link) error
    Get(code string) (domain.Link, error)
    IncrementClicks(code string) error
    CountByOwner(ownerID string) (int, error)
    List() []domain.Link
}

type InMemoryLinkRepository struct {
    mu    sync.Mutex
    links map[string]domain.Link
}
```

`IncrementClicks` and `CountByOwner` exist as their own methods rather than being folded into `Get` because they express two facts the interface should say out loud: incrementing a click count is a **write**, and enforcing a per-user limit needs a **count**, not a full list, scan, and count in the Service layer every time.

Real, captured output from the repository's own test suite:

```
=== RUN   TestInMemoryLinkRepository_CreateAndGet
--- PASS: TestInMemoryLinkRepository_CreateAndGet (0.00s)
=== RUN   TestInMemoryLinkRepository_CreateDuplicateCodeRejected
--- PASS: TestInMemoryLinkRepository_CreateDuplicateCodeRejected (0.00s)
=== RUN   TestInMemoryLinkRepository_GetMissing
--- PASS: TestInMemoryLinkRepository_GetMissing (0.00s)
=== RUN   TestInMemoryLinkRepository_IncrementClicks
--- PASS: TestInMemoryLinkRepository_IncrementClicks (0.00s)
=== RUN   TestInMemoryLinkRepository_CountByOwner
--- PASS: TestInMemoryLinkRepository_CountByOwner (0.00s)
=== RUN   TestInMemoryUserRepository_CreateAndGet
--- PASS: TestInMemoryUserRepository_CreateAndGet (0.00s)
=== RUN   TestInMemoryUserRepository_DuplicateUsernameRejected
--- PASS: TestInMemoryUserRepository_DuplicateUsernameRejected (0.00s)
=== RUN   TestInMemoryUserRepository_GetMissing
--- PASS: TestInMemoryUserRepository_GetMissing (0.00s)
PASS
ok  	level40.example/urlshortener/repository	1.399s
```

`UserRepository` follows the identical shape (`Create`, `GetByUsername`, `GetByID`), enforcing username uniqueness at the storage layer while leaving password validation entirely to the Service. `EXERCISES.md` Milestones 1-2 build both from scratch and run these exact tests.

---

## Building the Service Layer

The Service layer is where every business rule from Section 3 actually lives - the Repository has no concept of "valid URL" or "per-user limit," and the Handler shouldn't either.

```go
const MaxLinksPerUser = 5

func (s *LinkService) CreateLink(ownerID, longURL string) (domain.Link, error) {
    if err := validateURL(longURL); err != nil {
        return domain.Link{}, fmt.Errorf("create link: %w", err)
    }

    count, err := s.repo.CountByOwner(ownerID)
    if err != nil {
        return domain.Link{}, fmt.Errorf("create link: %w", err)
    }
    if count >= MaxLinksPerUser {
        return domain.Link{}, fmt.Errorf("create link: %w", ErrLinkLimit)
    }

    for attempt := 0; attempt < maxCodeAttempts; attempt++ {
        code, err := s.newCode(s.codeLen)
        // ...generate, attempt Create, retry only on a genuine code collision
    }
}
```

Short codes are generated with `crypto/rand`, not `math/rand` - a predictable short code would let an attacker enumerate other users' links, the same reasoning Level 32 applied to JWT secrets.

`UserService` carries the authentication rules: bcrypt hashing on `Register` (Level 32 - never a plaintext password touches storage), and JWT issuing/parsing on `Login`/`ParseUserID`, hard-coding HMAC-only signature verification exactly the way Level 32's `alg:none` defense requires. Both services depend only on a small, consumer-side repository interface (`LinkRepository`, `UserRepository`) defined in the `service` package itself, per Level 30's Best Practice #1.

Real, captured output proving the per-user limit and the read/write click-count split both work, entirely without a real repository (a hand-rolled fake, Level 30's testing payoff):

```
=== RUN   TestLinkService_CreateLink_Success
--- PASS: TestLinkService_CreateLink_Success (0.00s)
=== RUN   TestLinkService_CreateLink_InvalidURLRejectedBeforeRepository
--- PASS: TestLinkService_CreateLink_InvalidURLRejectedBeforeRepository (0.00s)
=== RUN   TestLinkService_CreateLink_EnforcesPerUserLimit
--- PASS: TestLinkService_CreateLink_EnforcesPerUserLimit (0.00s)
=== RUN   TestLinkService_Resolve_IncrementsClicksAndStatsDoesNot
--- PASS: TestLinkService_Resolve_IncrementsClicksAndStatsDoesNot (0.00s)
=== RUN   TestLinkService_Resolve_NotFound
--- PASS: TestLinkService_Resolve_NotFound (0.00s)
=== RUN   TestGenerateCode_Length
--- PASS: TestGenerateCode_Length (0.00s)
=== RUN   TestGenerateCode_Uniqueness
--- PASS: TestGenerateCode_Uniqueness (0.00s)
=== RUN   TestUserService_RegisterAndLogin_Success
--- PASS: TestUserService_RegisterAndLogin_Success (0.16s)
=== RUN   TestUserService_Register_ShortPasswordRejected
--- PASS: TestUserService_Register_ShortPasswordRejected (0.00s)
=== RUN   TestUserService_Login_WrongPasswordRejected
--- PASS: TestUserService_Login_WrongPasswordRejected (0.15s)
=== RUN   TestUserService_Login_UnknownUserGetsSameError
--- PASS: TestUserService_Login_UnknownUserGetsSameError (0.00s)
=== RUN   TestUserService_ParseUserID_RejectsTokenSignedWithDifferentSecret
--- PASS: TestUserService_ParseUserID_RejectsTokenSignedWithDifferentSecret (0.15s)
PASS
ok  	level40.example/urlshortener/service	2.593s
```

`TestLinkService_CreateLink_InvalidURLRejectedBeforeRepository` uses a **poison repository** - every method panics - to prove validation genuinely happens before storage is ever touched, the same technique Level 30 introduced.

---

## Building the Handler Layer

The Handler layer stays thin: decode the request, call the Service, translate the result into an HTTP response. It uses the stdlib `http.ServeMux` with method-and-wildcard routing (Level 27), no third-party router.

```go
func (h *LinkHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
    ownerID, ok := UserIDFromContext(r.Context())
    if !ok {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }
    var req createLinkRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid JSON body", http.StatusBadRequest)
        return
    }
    link, err := h.service.CreateLink(ownerID, req.LongURL)
    if err != nil {
        http.Error(w, err.Error(), linkStatusFor(err))
        return
    }
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(link)
}
```

`ownerID` arrives from `r.Context()`, never from the request body - the JWT middleware (Section 8) is the only thing allowed to say who the caller is. `linkStatusFor` is Level 30's single-place error-to-status translation:

| Service Error | HTTP Status |
|---|---|
| `ErrInvalidURL` | 400 Bad Request |
| `ErrLinkLimit` | 429 Too Many Requests |
| `ErrLinkNotFound` | 404 Not Found |
| anything else | 500 Internal Server Error |

`routes.go` wires each endpoint to exactly the middleware it needs - `NewMux` takes the middleware as parameters rather than importing the `middleware` package directly, keeping the dependency arrow one-way (`middleware` depends on `handler`, never the reverse):

```go
func NewMux(linkH *LinkHandler, authH *AuthHandler, logging, requireAuth, rateLimit Middleware) *http.ServeMux {
    mux := http.NewServeMux()
    mux.Handle("POST /auth/register", chain(http.HandlerFunc(authH.Register), logging))
    mux.Handle("POST /auth/login", chain(http.HandlerFunc(authH.Login), logging))
    mux.Handle("POST /shorten", chain(http.HandlerFunc(linkH.CreateLink), logging, requireAuth))
    mux.Handle("GET /{code}", chain(http.HandlerFunc(linkH.Redirect), logging, rateLimit))
    mux.Handle("GET /stats/{code}", chain(http.HandlerFunc(linkH.Stats), logging))
    return mux
}
```

Real, captured `httptest`-based output, testing every handler in isolation from the real Service (fakes satisfying `LinkService`/`AuthService`, Level 30's isolated-handler-testing payoff):

```
=== RUN   TestCreateLinkHandler_Success
--- PASS: TestCreateLinkHandler_Success (0.00s)
=== RUN   TestCreateLinkHandler_InvalidURL
--- PASS: TestCreateLinkHandler_InvalidURL (0.00s)
=== RUN   TestRedirectHandler_Success
--- PASS: TestRedirectHandler_Success (0.00s)
=== RUN   TestRedirectHandler_NotFound
--- PASS: TestRedirectHandler_NotFound (0.00s)
=== RUN   TestStatsHandler_DoesNotRequireAuth
--- PASS: TestStatsHandler_DoesNotRequireAuth (0.00s)
=== RUN   TestRegisterHandler_UsernameTaken
--- PASS: TestRegisterHandler_UsernameTaken (0.00s)
=== RUN   TestLoginHandler_Success
--- PASS: TestLoginHandler_Success (0.00s)
PASS
ok  	level40.example/urlshortener/handler	0.472s
```

---

## Adding Authentication

`POST /shorten` is protected; `GET /{code}` and `GET /stats/{code}` are public, exactly as Section 3's table specified. Level 32's JWT middleware pattern - `func(http.Handler) http.Handler`, extracting a Bearer token, validating it, injecting the result into the request context (Level 24) - applies unchanged:

```go
func RequireAuth(parser TokenParser) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            const prefix = "Bearer "
            if !strings.HasPrefix(authHeader, prefix) {
                http.Error(w, "missing or malformed Authorization header", http.StatusUnauthorized)
                return
            }
            userID, err := parser.ParseUserID(strings.TrimPrefix(authHeader, prefix))
            if err != nil {
                http.Error(w, "invalid or expired token", http.StatusUnauthorized)
                return
            }
            next.ServeHTTP(w, r.WithContext(handler.WithUserID(r.Context(), userID)))
        })
    }
}
```

`TokenParser` is a one-method consumer-side interface (`ParseUserID`) - `*service.UserService` satisfies it implicitly, and a test can fake it without touching real JWTs at all. Real, captured output:

```
=== RUN   TestRequireAuth_ValidTokenInjectsUserID
--- PASS: TestRequireAuth_ValidTokenInjectsUserID (0.00s)
=== RUN   TestRequireAuth_MissingHeaderRejected
--- PASS: TestRequireAuth_MissingHeaderRejected (0.00s)
=== RUN   TestRequireAuth_InvalidTokenRejected
--- PASS: TestRequireAuth_InvalidTokenRejected (0.00s)
```

And a real, live run against the actual running server (`go build` + the compiled binary, real `curl`) proves the whole path end to end - register, log in, shorten with the token, and a rejected attempt without one:

```
=== POST /shorten (Authorization: Bearer <token>) ===
{"code":"YXD97v0","long_url":"https://go.dev/doc/effective_go","owner_id":"u_1788285897194322000","clicks":0,"created_at":"2026-09-01T23:34:57.292782+05:30"}

HTTP 201

=== POST /shorten (no Authorization header) ===
missing or malformed Authorization header

HTTP 401
```

---

## Adding Structured Logging and Configuration

Every request gets one structured log line (Level 33's `log/slog`), emitted by a middleware that depends on an **injected** `*slog.Logger` rather than a package-level global (Level 31's habit, applied to middleware this time):

```go
func Logging(logger *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
            next.ServeHTTP(rec, r)
            logger.Info("http request",
                slog.String("method", r.Method),
                slog.String("path", r.URL.Path),
                slog.Int("status", rec.status),
                slog.Duration("latency", time.Since(start)),
            )
        })
    }
}
```

Real captured JSON log output from the live server run above:

```json
{"time":"2026-09-01T23:34:57.194387+05:30","level":"INFO","msg":"http request","method":"POST","path":"/auth/register","status":201,"latency":77781500}
{"time":"2026-09-01T23:34:57.313757+05:30","level":"INFO","msg":"http request","method":"GET","path":"/YXD97v0","status":302,"latency":21542}
```

Configuration follows Level 34's pattern exactly: one `Config` struct, one `LoadConfig()` function, everything else in the program depends on `Config` and never calls `os.Getenv` directly:

```go
type Config struct {
    Port             string
    JWTSecret        string
    TokenTTL         time.Duration
    CodeLength       int
    RateLimitPerSec  float64
    RateLimitBurst   int
}
```

`JWT_SECRET` fails fast if unset (Level 34's Best Practice #3 - never silently default a security-sensitive setting) unless a developer explicitly opts in with `ALLOW_DEV_JWT_SECRET=true` for local work. Writing this level's own config test caught a real, genuine bug worth calling out: `os.LookupEnv` correctly treats an explicitly-set empty string as **present**, not **unset** (exactly Level 34's Section 2 lesson) - a test that did `t.Setenv("PORT", "")` expecting the default to apply was wrong, not the code. `EXERCISES.md` Milestone 7 shows the failing run and the fix.

---

## Adding a Rate Limiter

The public redirect endpoint (`GET /{code}`) is the one part of this API the whole internet can hit without logging in first, which makes it the honest place to enforce a limit. This bridges Level 39's coverage of rate-limiting algorithms into a real, runnable **token bucket**:

```go
type Limiter struct {
    mu            sync.Mutex
    buckets       map[string]*bucket
    ratePerSecond float64
    burst         int
    now           func() time.Time
}

func (l *Limiter) Allow(key string) bool {
    l.mu.Lock()
    defer l.mu.Unlock()
    now := l.now()
    b, ok := l.buckets[key]
    if !ok {
        b = &bucket{tokens: float64(l.burst), lastRefill: now}
        l.buckets[key] = b
    }
    elapsed := now.Sub(b.lastRefill).Seconds()
    b.tokens = min(b.tokens+elapsed*l.ratePerSecond, float64(l.burst))
    b.lastRefill = now
    if b.tokens < 1 {
        return false
    }
    b.tokens--
    return true
}
```

Each client (keyed by remote address) gets a bucket that starts full (so the very first request is never throttled), refills continuously at `ratePerSecond`, and caps at `burst`. `l.now` is injected (not a bare `time.Now()` call) specifically so a test can fake elapsed time deterministically instead of sleeping for real. Real captured output, including a refill proven with a faked clock:

```
=== RUN   TestLimiter_AllowsUpToBurstThenBlocks
--- PASS: TestLimiter_AllowsUpToBurstThenBlocks (0.00s)
=== RUN   TestLimiter_RefillsOverTime
--- PASS: TestLimiter_RefillsOverTime (0.00s)
=== RUN   TestLimiter_TracksKeysIndependently
--- PASS: TestLimiter_TracksKeysIndependently (0.00s)
=== RUN   TestMiddleware_RejectsOverLimitWith429
--- PASS: TestMiddleware_RejectsOverLimitWith429 (0.00s)
PASS
ok  	level40.example/urlshortener/ratelimit	1.608s
```

---

## Testing the Whole Thing

Every layer above has its own unit tests (Level 26's table-driven, `t.Run`-based style throughout). The capstone adds one more thing none of the isolated levels needed: a **full end-to-end test** that builds the real dependency graph - real repositories, real services, real handlers, real middleware - and drives it through `httptest.NewServer` over an actual TCP connection on loopback, the same technique Level 30 used to prove a whole stack works together:

```go
func TestEndToEnd_FullUserJourney(t *testing.T) {
    srv := newTestServer(t) // the REAL buildServer() from main.go

    // 1. Register -> 2. Log in -> 3. Create a link (protected) ->
    // 3b. The same request without a token must be rejected ->
    // 4. Follow the public redirect -> 5. Confirm stats counted it
}
```

Real, captured `go test` output for the entire module - every package, every layer, the full end-to-end flow, and the rate-limiter proof, all in one run:

```
ok  	level40.example/urlshortener	2.492s
ok  	level40.example/urlshortener/config	0.247s
?   	level40.example/urlshortener/domain	[no test files]
ok  	level40.example/urlshortener/handler	1.628s
ok  	level40.example/urlshortener/middleware	1.090s
ok  	level40.example/urlshortener/ratelimit	0.547s
ok  	level40.example/urlshortener/repository	2.880s
ok  	level40.example/urlshortener/service	3.156s
```

`EXERCISES.md` Milestone 9 walks through building `TestEndToEnd_RateLimiterProtectsRedirect` and hitting a second genuine bug along the way - not draining an HTTP response body before closing it silently defeats connection reuse, which silently defeats a rate limiter keyed by remote address in a test. Both real bugs found during this level's own authoring (the config one above, and this one) are left in the write-up on purpose: real testing finds real bugs, even in code written carefully.

---

## Containerizing It and Sketching Kubernetes

> ⚠️ **Same honesty pattern as Levels 35 and 36.** This level's Go code is written and verified the normal way - `go build`, `go test`, real captured output. The Docker and Kubernetes steps below could not be executed end-to-end in the environment that authored this course: `docker version` shows the Docker CLI is installed, but its daemon was not reachable (`Cannot connect to the Docker daemon at unix:///Users/macbookpro/.docker/run/docker.sock. Is the docker daemon running?`) - the identical situation Level 35 documented. `kubectl` is installed, but with no cluster reachable, even `kubectl apply --dry-run=client` fails, because client-side dry-run still needs to fetch the OpenAPI schema from a live API server - the identical situation Level 36 documented. Every file below is believed correct per Docker's and Kubernetes' documented, stable behavior, and was checked for well-formed syntax, but treat it as **illustrative, not verified**, exactly like the corresponding sections of Levels 35 and 36. Run it yourself once you have a working daemon and cluster - your own output is the one that counts.

A minimal multi-stage `Dockerfile` (Level 35's pattern - build in a full Go image, run from a minimal one):

```dockerfile
# --- build stage ---
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /urlshortener .

# --- run stage ---
FROM gcr.io/distroless/static-debian12
COPY --from=build /urlshortener /urlshortener
EXPOSE 8080
ENV PORT=8080
ENTRYPOINT ["/urlshortener"]
```

> ⚠️ **Illustrative output (unverified in this environment)** - Docker's daemon was not reachable here:
>
> ```
> $ docker build -t urlshortener:1.0.0 .
> [+] Building 14.2s (12/12) FINISHED
>  => [build 1/5] FROM docker.io/library/golang:1.26-alpine
>  => [build 5/5] RUN CGO_ENABLED=0 go build -o /urlshortener .
>  => [stage-1 2/2] COPY --from=build /urlshortener /urlshortener
>  => exporting to image
>  => => naming to docker.io/library/urlshortener:1.0.0
> ```

A minimal Kubernetes `Deployment` + `Service` (Level 36's pattern - a Deployment for self-healing replicas, a `ClusterIP` Service in front of them, config injected via env vars rather than baked into the image):

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: urlshortener
spec:
  replicas: 2
  selector:
    matchLabels:
      app: urlshortener
  template:
    metadata:
      labels:
        app: urlshortener
    spec:
      containers:
        - name: urlshortener
          image: urlshortener:1.0.0
          ports:
            - containerPort: 8080
          env:
            - name: PORT
              value: "8080"
            - name: JWT_SECRET
              valueFrom:
                secretKeyRef:
                  name: urlshortener-secrets
                  key: jwt-secret
---
# k8s/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: urlshortener
spec:
  type: ClusterIP
  selector:
    app: urlshortener
  ports:
    - port: 80
      targetPort: 8080
```

`JWT_SECRET` is injected from a Kubernetes `Secret` (`urlshortener-secrets`), never baked into the image - exactly Level 36's separation of a built image from its runtime configuration, and exactly Level 34's "never hardcode a secret" rule, now enforced at the cluster level too.

---

## What's Next

This capstone is deliberately complete-but-small, per Section 14's Best Practices below. A few natural next steps for a reader who wants to keep extending it:

- **Swap in real persistence.** Replace `InMemoryLinkRepository`/`InMemoryUserRepository` with Level 29's `modernc.org/sqlite` approach - because both are hidden behind the `LinkRepository`/`UserRepository` interfaces, the Service and Handler layers need zero changes. Add a real migration tool (e.g. `golang-migrate`) once schema changes need to be tracked over time.
- **Add metrics.** A `/metrics` endpoint exposing request counts and latencies (Prometheus-style) sits naturally alongside the existing logging middleware - same request lifecycle, different observability signal.
- **Add a frontend.** A tiny static HTML page that calls `/auth/login` and `/shorten` via `fetch` would exercise the API from a real browser client instead of `curl`, and is a natural place to first feel the effects of CORS.

None of these require touching more than one layer at a time - which is exactly the architectural payoff this level set out to demonstrate.

---

## Best Practices

### 1. Design the API Contract Before Writing Any Code

Section 3's endpoint table and data model were written down *before* a single `.go` file existed. Deciding "what are the endpoints, what does each one need, what does each one return" up front turns building into filling in a known shape, instead of discovering the shape by accident three layers deep.

### 2. Build in One Consistent Direction: Repository → Service → Handler

This level built bottom-up - storage first, then business rules, then HTTP - because each layer's tests could lean on the layer below it actually existing and working. Top-down (Handler first, with the Service and Repository stubbed) is equally valid; picking one direction and sticking to it for the whole project is what matters, not which direction.

### 3. Write Tests As You Build, Not At the End

Every milestone in `EXERCISES.md` runs `go test` before moving to the next one. This level's own authoring caught two genuine bugs this way (the `os.LookupEnv("")` config gotcha in Section 9, and the connection-reuse-vs-rate-limiting bug in Section 11) - both would have been far harder to isolate in a "write everything, test once at the end" pass, because there would have been ten new lines of code to suspect instead of one.

### 4. Keep the Composition Root as the Only Place That Knows Concrete Types

`buildServer` in `main.go` is the only function in the entire project that imports `repository.NewInMemoryLinkRepository` by name. Every other function - every Service, every Handler, every middleware - depends on an interface. This is Level 31's constructor injection applied at the scale of a whole application, not just one struct.

### 5. Let Errors Cross Layers as Sentinels, Translate Status Codes in One Place

`linkStatusFor` and `authStatusFor` are the only two functions in the whole project that know what an HTTP status code is. Every error underneath them is a plain Go `error`, checked with `errors.Is` - Level 16's and Level 30's rule, unchanged at capstone scale.

---

## Common Mistakes

### Mistake 1: Skipping the Design Step and Coding Straight Into Handlers

```go
// ❌ WRONG - designing the API by writing whatever the first handler
// happens to need, discovering the real shape three refactors later
func shortenHandler(w http.ResponseWriter, r *http.Request) {
    // wait, what does the request body look like again?
}

// ✅ RIGHT - Section 3's table and struct definitions exist BEFORE this
// function does. Writing the handler becomes translation, not invention.
```

A project designed handler-first tends to grow whatever shape the first `curl` command needed, then breaks every existing caller the moment a second requirement shows up.

### Mistake 2: Mixing Layers - Business Logic Leaking Into Handlers

```go
// ❌ WRONG - the handler is now also the service
func (h *LinkHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
    var req createLinkRequest
    json.NewDecoder(r.Body).Decode(&req)
    if !strings.HasPrefix(req.LongURL, "http") {
        http.Error(w, "bad url", 400)
        return
    }
    h.repo.Create(...) // also wrong - skips the service, and the handler
                        // now imports the repository directly
}

// ✅ RIGHT - the handler decodes and delegates; validation and storage
// both live below it, exactly as Section 6 and Section 7 built them
func (h *LinkHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
    var req createLinkRequest
    json.NewDecoder(r.Body).Decode(&req)
    link, err := h.service.CreateLink(ownerID, req.LongURL)
    if err != nil {
        http.Error(w, err.Error(), linkStatusFor(err))
        return
    }
    json.NewEncoder(w).Encode(link)
}
```

### Mistake 3: Not Testing Until the Whole Thing Is "Done"

Writing all ten milestones' worth of code first and running `go test ./...` only at the very end means a failure could be hiding in the repository, the service, or the handler - and there's no way to tell which without adding tests retroactively anyway. Both real bugs this level's own authoring hit (Sections 9 and 11) were caught within seconds of being introduced, specifically because a test ran immediately after the code that could break it.

### Mistake 4: Over-Scoping the Capstone Instead of Shipping Something Complete But Small

```
❌ WRONG: "My capstone will have URL shortening, AND a full admin dashboard,
AND email notifications, AND an analytics pipeline, AND OAuth..."
   -> six months later, nothing is finished, and nothing is tested

✅ RIGHT: One resource (Link), one auth flow, one rate-limited endpoint,
   fully built, fully tested, fully containerized - actually done.
```

A small, complete, tested, shippable service teaches - and demonstrates - far more than a large, half-built one. Every "What's Next" idea in Section 13 is an extension precisely *because* the base is small enough to extend cleanly.

---

## Summary

**The Capstone:**
- A URL Shortener API - `POST /auth/register`, `POST /auth/login`, `POST /shorten` (JWT-protected), `GET /{code}` (public, rate-limited redirect), `GET /stats/{code}` (public, read-only)

**Architecture:**
- Repository (in-memory, Level 30) → Service (validation, limits, JWT, Level 30/32) → Handler (`net/http`, Level 27) → Middleware (logging, auth, rate limiting)
- Wired entirely by hand in a `buildServer` composition root (Level 31) - the only place that knows any concrete type

**Cross-Cutting Concerns:**
- Structured logging (Level 33, `log/slog`) on every request
- Environment-driven configuration (Level 34) that fails fast on a missing `JWT_SECRET`
- A token-bucket rate limiter (bridging Level 39) protecting the one public, high-traffic endpoint

**Verification:**
- Every layer has its own unit tests; one full `httptest.NewServer`-based end-to-end test drives the real HTTP flow (Level 26)
- Every milestone's code was actually run - `go build`, `go test`, `go run`, and real `curl` against a live server - nothing hand-computed
- Docker and Kubernetes steps follow Levels 35/36's established honesty pattern: written and syntax-checked, explicitly flagged as unverified against a real daemon/cluster in this environment

---

## Next Steps

You now understand:
- ✅ How to design an API's endpoints and data model before writing code
- ✅ How to build a complete Repository-Service-Handler stack from scratch, layer by layer, verifying each one before moving to the next
- ✅ How to protect one endpoint with JWT auth while leaving others public, using the same request context Level 24 introduced
- ✅ How to add structured logging and environment-driven configuration as cross-cutting middleware and startup concerns
- ✅ How to bridge Level 39's rate-limiting theory into a real, tested token-bucket middleware
- ✅ How to prove an entire system works with one full end-to-end test, on top of independent per-layer tests
- ✅ How to carry Levels 35-36's Docker/Kubernetes honesty pattern into your own future projects

**Next level:** Level 41 - Go Interview Preparation
- Turning everything you've built across all 40 levels into confident answers to real interview questions
- Common Go interview topics: concurrency, interfaces, error handling, memory model
- Practice explaining the exact architectural decisions this level made out loud

You just built a real, working, tested backend service using nothing but what you already knew. That's not a metaphor for what real engineering feels like - it *is* what real engineering feels like. Congratulations, and welcome to Level 41! 🚀
