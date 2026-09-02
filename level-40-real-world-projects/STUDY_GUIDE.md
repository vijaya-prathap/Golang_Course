# Level 40: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: Foundations (Milestones 1-4)
```
Day 1:  Milestone 1 - Project setup, domain models, Link repository
Day 2:  Milestone 2 - User repository
Day 3:  Milestone 3 - LinkService: validation, per-user limit, short codes
Day 4:  Milestone 4 - UserService: bcrypt hashing, JWT issue/parse
```

### Week 2: HTTP and Security (Milestones 5-6)
```
Day 1-2: Milestone 5 - Handler layer, NewMux, isolated handler tests
Day 3-4: Milestone 6 - Real JWT auth middleware, replacing the test stub
```

### Week 3: Observability, Resilience, Proof (Milestones 7-9)
```
Day 1-2: Milestone 7 - Structured logging + configuration (and a real bug!)
Day 3:   Milestone 8 - Token-bucket rate limiter
Day 4-5: Milestone 9 - main.go wiring + full end-to-end test (a second real bug!)
```

### Week 4: Shipping (Milestone 10 + Review)
```
Day 1: Milestone 10 - Dockerfile + Kubernetes manifests (illustrative)
Day 2: Bonus challenges
Day 3-4: Review architecture end to end; re-read README.md Sections 4, 11-12
```

---

## 🏗️ The Complete Architecture

Every layer this course built, in one running system:

```
                                    ┌──────────────────────┐
                                    │        main()          │
                                    │  composition root       │
                                    │     (Level 31)           │
                                    └───────────┬──────────────┘
                                                │ builds & wires everything
                                                ▼
┌───────────────────────────────────────────────────────────────────────────┐
│                          MIDDLEWARE (wraps handlers)                       │
│                                                                             │
│   Logging (Level 33)  ─── every request                                   │
│   RequireAuth (Level 32) ─── POST /shorten only                           │
│   RateLimit (bridges 39) ─── GET /{code} only                             │
└───────────────────────────────────────────────────────────────────────────┘
                                                │
                                                ▼
┌───────────────────────────────────────────────────────────────────────────┐
│                         HANDLER  (Level 27, Level 30)                     │
│                                                                             │
│   LinkHandler: CreateLink, Redirect, Stats                                 │
│   AuthHandler: Register, Login                                             │
│   -> decodes JSON, calls the Service, encodes the response                 │
│   -> depends on LinkService / AuthService INTERFACES only                  │
└───────────────────────────────────────────────────────────────────────────┘
                                                │ interface
                                                ▼
┌───────────────────────────────────────────────────────────────────────────┐
│                         SERVICE  (Level 30, 31, 32)                      │
│                                                                             │
│   LinkService: validateURL, MaxLinksPerUser, GenerateCode (crypto/rand)    │
│   UserService: bcrypt.GenerateFromPassword/CompareHashAndPassword,         │
│                jwt.NewWithClaims / jwt.ParseWithClaims (HS256 only)        │
│   -> ALL business rules live here, nowhere else                           │
│   -> depends on LinkRepository / UserRepository INTERFACES only           │
└───────────────────────────────────────────────────────────────────────────┘
                                                │ interface
                                                ▼
┌───────────────────────────────────────────────────────────────────────────┐
│                       REPOSITORY  (Level 30, Level 23)                   │
│                                                                             │
│   InMemoryLinkRepository:  map[string]domain.Link + sync.Mutex            │
│   InMemoryUserRepository:  map[string]domain.User + sync.Mutex            │
│   -> pure storage, no validation, no HTTP, no JWTs                        │
└───────────────────────────────────────────────────────────────────────────┘
```

**Dependency direction, unchanged from Level 30 at any scale:** Handler → Service → Repository (interface). `middleware` depends on `handler` (for the context-key helpers); `handler` never depends on `middleware`. Nothing below the Repository knows anything above it exists.

---

## 🗺️ API Endpoint Reference

| Method | Path | Auth? | Rate-Limited? | Service Method | Success | Failure Modes |
|---|---|:---:|:---:|---|---|---|
| `POST` | `/auth/register` | No | No | `UserService.Register` | `201` | `400` invalid input, `409` username taken |
| `POST` | `/auth/login` | No | No | `UserService.Login` | `200` `{token}` | `401` invalid credentials |
| `POST` | `/shorten` | **Yes** | No | `LinkService.CreateLink` | `201` Link | `400` invalid URL, `401` no/bad token, `429` per-user limit |
| `GET` | `/{code}` | No | **Yes** | `LinkService.Resolve` | `302` redirect | `404` not found, `429` rate-limited |
| `GET` | `/stats/{code}` | No | No | `LinkService.Stats` | `200` Link | `404` not found |
| `GET` | `/healthz` | No | No | (none) | `200` `ok` | - |

---

## 🔀 Request Flow: One Full Journey

The middleware pipeline positions, shown together for reference. In practice, `RequireAuth` only wraps `POST /shorten` and `RateLimit` only wraps `GET /{code}` (see the table above) - no single real request passes through both. Here is what each of those two protected paths looks like end to end:

### Path A: `POST /shorten` (protected by JWT)

```
 client
   │  POST /shorten
   │  Authorization: Bearer <jwt>
   │  {"long_url": "https://..."}
   ▼
┌─────────────┐
│  Logging     │  starts a timer, notes method+path
└──────┬──────┘
       ▼
┌─────────────┐
│ RequireAuth  │  extracts Bearer token -> UserService.ParseUserID
│ (Level 32)   │  invalid/missing? -> 401, STOP HERE
└──────┬──────┘  valid -> injects userID into context (Level 24)
       ▼
┌─────────────┐
│ LinkHandler  │  decodes {"long_url": "..."} from the body
│ .CreateLink  │  reads userID from context (never from the body)
└──────┬──────┘
       ▼
┌─────────────┐
│ LinkService  │  validateURL -> CountByOwner -> compare to MaxLinksPerUser
│ .CreateLink  │  -> GenerateCode (crypto/rand) -> repo.Create
└──────┬──────┘
       ▼
┌─────────────┐
│ LinkRepo-    │  map write, guarded by mutex
│ sitory       │  returns ErrCodeTaken on the rare collision (retried above)
└──────┬──────┘
       ▼
   201 Created + Link JSON
       │
       ▼
┌─────────────┐
│  Logging     │  logs method, path, status=201, latency
└─────────────┘
```

### Path B: `GET /{code}` (public, rate-limited)

```
 client
   │  GET /{code}
   ▼
┌─────────────┐
│  Logging     │  starts a timer
└──────┬──────┘
       ▼
┌─────────────┐
│ RateLimit    │  Limiter.Allow(r.RemoteAddr)
│ (bridges 39) │  no tokens left? -> 429, STOP HERE
└──────┬──────┘  token available -> consume one, continue
       ▼
┌─────────────┐
│ LinkHandler  │  reads {code} path value
│ .Redirect    │
└──────┬──────┘
       ▼
┌─────────────┐
│ LinkService  │  repo.Get(code) -> repo.IncrementClicks(code)
│ .Resolve     │  (a WRITE - never called by Stats)
└──────┬──────┘
       ▼
┌─────────────┐
│ LinkRepo-    │  map read + map write, guarded by mutex
│ sitory       │
└──────┬──────┘
       ▼
   302 Found, Location: <long url>
       │
       ▼
┌─────────────┐
│  Logging     │  logs method, path, status=302, latency
└─────────────┘
```

---

## 🧭 Milestone Roadmap

```
 1. Domain + Link Repository ─────────┐
 2. User Repository ───────────────────┼── REPOSITORY layer complete
                                       │
 3. LinkService (validate, limit,     │
    generate codes) ──────────────────┤
 4. UserService (bcrypt, JWT) ─────────┴── SERVICE layer complete
                                       │
 5. Handler layer + NewMux ────────────┴── HANDLER layer complete
                                       │
 6. JWT auth middleware ────────────────┐
 7. Logging + Config ───────────────────┼── MIDDLEWARE + config complete
 8. Token-bucket rate limiter ──────────┘   (real bug caught in #7!)
                                       │
 9. main.go wiring + end-to-end test ──┴── WHOLE SYSTEM proven together
                                             (a second real bug caught here!)
10. Dockerfile + Kubernetes manifests ───── SHIPPED (illustrative in this env)
```

Each arrow into a milestone assumes every file from every earlier milestone already compiles and passes its tests - exactly the "build bottom-up, test as you go" discipline Section 14 of `README.md` names explicitly.

---

## 🚨 Common Mistakes (Quick Reference)

| Mistake | Wrong | Right |
|---|---|---|
| Designing while coding | Inventing the request shape inside the first handler | Write the endpoint table and structs FIRST (Section 3) |
| Business logic in the handler | `if !strings.HasPrefix(url, "http")` inside a handler func | Validation lives in the Service (`validateURL`) |
| Skipping the Service layer | Handler calls `repo.Create` directly | Handler calls `service.CreateLink`, always |
| `os.LookupEnv` empty-string trap | `t.Setenv("PORT", "")` expecting the default to apply | `os.Unsetenv("PORT")` - empty is PRESENT, not unset |
| Not draining a response body | `resp.Body.Close()` without reading first | `io.Copy(io.Discard, resp.Body)` before `Close()` - required for keep-alive reuse |
| Testing only at the end | Writing all 10 milestones, then running `go test ./...` once | Run tests after every milestone - this level's own two real bugs were caught this way |
| Global mutable state | A package-level `var db *sql.DB` | Constructor-injected dependencies everywhere (Level 31) |
| Over-scoping the capstone | Planning an admin dashboard, OAuth, and analytics before writing a line | One resource, fully built, fully tested, fully shippable |

---

## ✅ Prerequisites for Level 41

Before moving to Level 41 (Go Interview Preparation), you need:

- ✅ Comfortable explaining Repository-Service-Handler layering out loud, not just recognizing it in code
- ✅ Comfortable explaining WHY dependencies are injected as interfaces, not concrete types
- ✅ Comfortable explaining how JWT auth middleware and a rate limiter both use the same `func(http.Handler) http.Handler` shape
- ✅ Comfortable explaining the difference between `os.Getenv` and `os.LookupEnv`, from having been bitten by it
- ✅ Comfortable explaining why HTTP keep-alive connection reuse requires draining a response body
- ✅ Completed all 10 milestones with real, passing `go test` output

---

## 🎯 Ready for Level 41?

Level 41 turns everything you just built into spoken, interview-ready answers. Before moving on, make sure you can answer each of these OUT LOUD, in under two minutes, without looking at your code:

- [ ] "Walk me through what happens, layer by layer, when a client calls `POST /shorten`."
- [ ] "Why does `LinkService` depend on a `LinkRepository` interface instead of `*InMemoryLinkRepository` directly?"
- [ ] "How would you swap this project's storage from in-memory to SQLite without touching the Service or Handler?"
- [ ] "Why is `GET /stats/{code}` a separate service method from `GET /{code}`'s resolve-and-redirect, instead of one method with a flag?"
- [ ] "What does your rate limiter do if two requests from the same client arrive at literally the same instant?" (Hint: think about the mutex.)
- [ ] "What are the two real bugs this level's own test suite caught, and what did each one teach you?"

---

## 📋 Checklist Before Level 41

- [ ] All 10 milestones completed with real, captured `go test` / `go run` output
- [ ] `go test ./...` passes cleanly across every package
- [ ] Can explain the full request flow diagrams above without referring back to them
- [ ] Completed at least one bonus challenge (click analytics, link expiration, or custom aliases)
- [ ] Comfortable narrating this project's architecture as if explaining it to an interviewer

---

## 💡 Key Takeaways

### The Payoff
Every level from 1-39 taught one tool. This level is the first time all of them had to work together, under real HTTP traffic, with real tests proving it - which is what "software engineering," as opposed to "learning Go syntax," actually means day to day.

### The Two Real Bugs
Both bugs this level's own authoring hit - `os.LookupEnv`'s empty-vs-unset distinction, and HTTP keep-alive's dependence on draining a response body - were caught within seconds by a test, not discovered in production. That is Best Practice #3 (test as you build) working exactly as intended.

### The Shape Generalizes
Repository-Service-Handler, constructor injection, JWT middleware, structured logging middleware, and a rate limiter are not URL-shortener-specific. The exact same shape builds a Todo API, a Notes API, or almost any other small backend service (Section 2 of `README.md`).

---

## 📚 Next Level

Level 41: Go Interview Preparation
- Turning this capstone's architecture into confident, spoken interview answers
- Common Go interview topics: concurrency, interfaces, error handling, the memory model
- Practicing system-design questions using this exact project as your worked example

You built a complete, working, tested system. Now let's make sure you can talk about it as confidently as you built it! 🚀
