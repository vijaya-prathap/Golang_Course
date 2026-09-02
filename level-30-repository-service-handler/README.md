# Level 30: Repository-Service-Handler - Complete Guide

## Introduction

Welcome to Level 30! You've built HTTP APIs (Level 27), a REST framework workflow (Level 28), and talked to a real database (Level 29). Every one of those levels could still be written as a single `main.go` file with everything jammed together: SQL queries, validation logic, and `http.ResponseWriter` calls all in the same function. It would work. It would also become unmaintainable and untestable the moment the program grew past a toy size.

**Repository-Service-Handler** is the architectural pattern that fixes this by giving every piece of code exactly **one job**:

- **Repository** - talks to storage. Nothing else.
- **Service** - enforces business rules. Nothing else.
- **Handler** - translates HTTP in and out. Nothing else.

This level is about the *shape* of the pattern, not about any specific storage technology. Every example here uses an **in-memory repository** - a map or a slice guarded by a mutex - instead of Level 29's SQLite setup. The whole point of the pattern is that the Repository is just an interface: it could be backed by SQL (Level 29), an in-memory store (this level), a file, a remote API, or a test fake, and the Service and Handler layers would never know the difference. Using in-memory storage here keeps every exercise 100% offline and fast, matching every level in this course except Levels 28-29.

---

## Table of Contents

1. [What the Pattern Is and Why](#what-the-pattern-is-and-why)
2. [The Repository Layer](#the-repository-layer)
3. [The Service Layer](#the-service-layer)
4. [The Handler Layer](#the-handler-layer)
5. [Dependency Direction](#dependency-direction)
6. [Wiring It Together in main()](#wiring-it-together-in-main)
7. [Why This Enables Easy Testing](#why-this-enables-easy-testing)
8. [Error Handling Across Layers](#error-handling-across-layers)
9. [When This Is Overkill](#when-this-is-overkill)
10. [Best Practices](#best-practices)
11. [Common Mistakes](#common-mistakes)

---

## What the Pattern Is and Why

A program with no architecture at all tends to grow like this: an HTTP handler function that reads the request body, runs a SQL query (or touches a map) directly, checks a few `if` statements for validity somewhere in the middle, and writes the response - all in one function, 150 lines long.

That function has at least three separate responsibilities tangled together:

1. Understanding HTTP (reading a request body, writing a status code)
2. Understanding business rules (is this input valid? is this allowed?)
3. Understanding storage (how records are actually stored and retrieved)

Repository-Service-Handler pulls those three responsibilities into three layers:

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Handler   │ --> │   Service   │ --> │ Repository  │
│ (transport) │     │ (business   │     │  (storage)  │
│             │     │   rules)    │     │             │
└─────────────┘     └─────────────┘     └─────────────┘
   net/http           validation,          in-memory
   JSON in/out        uniqueness,          map or slice
                       computed fields      (this level)
```

Each layer has one job, which means each layer can be:

- **Tested independently** - test business rules without spinning up an HTTP server or a real database
- **Changed independently** - swap the storage technology without touching a single validation rule
- **Understood independently** - a new team member reading the Handler doesn't need to understand how records are stored

This is not a Go-specific idea (it appears across languages as "layered architecture," "clean architecture," "hexagonal architecture," and half a dozen other names), but Go's implicit interface satisfaction (Level 15) makes it unusually easy to wire up without frameworks, base classes, or annotations.

---

## The Repository Layer

The Repository layer's entire job is: **given an interface describing storage operations, provide a concrete implementation.** It knows nothing about HTTP, nothing about validation - just how to save, fetch, update, and delete records.

```go
// User is the data record this level's examples work with.
type User struct {
    ID    string
    Name  string
    Email string
}

var ErrNotFound = errors.New("user not found")

// UserRepository is the Repository layer's contract: pure data access,
// no business rules, no HTTP. Any type with these four methods satisfies
// it implicitly (Level 15) - a map-backed store here, a SQL-backed store
// in Level 29, anything else tomorrow.
type UserRepository interface {
    Create(u User) error
    Get(id string) (User, error)
    Update(u User) error
    Delete(id string) error
}
```

And one concrete implementation, backed by a map guarded by a mutex (Level 23) so it stays safe once concurrent HTTP requests start hitting it:

```go
type InMemoryUserRepository struct {
    mu    sync.Mutex
    users map[string]User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
    return &InMemoryUserRepository{users: make(map[string]User)}
}

func (r *InMemoryUserRepository) Create(u User) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.users[u.ID] = u
    return nil
}

func (r *InMemoryUserRepository) Get(id string) (User, error) {
    r.mu.Lock()
    defer r.mu.Unlock()
    u, ok := r.users[id]
    if !ok {
        return User{}, fmt.Errorf("get user %s: %w", id, ErrNotFound)
    }
    return u, nil
}

// Update and Delete follow the same shape.
```

`InMemoryUserRepository` satisfies `UserRepository` **implicitly** - nothing declares "this type implements that interface." This is Level 15's core idea revisited: because satisfaction is structural, you can verify it compiles with a single line that does nothing at runtime except fail to *compile* if the methods ever drift:

```go
var _ UserRepository = (*InMemoryUserRepository)(nil)
```

Running a small program against this repository (real output, captured directly):

```
=== Create ===
created users 1 and 2

=== Get ===
found: {ID:1 Name:Alice Email:alice@example.com}

=== Update ===
after update: {ID:1 Name:Alice Updated Email:alice@example.com}

=== Delete ===
deleted user 2

=== Get a Missing User ===
error: get user 2: user not found
errors.Is(err, ErrNotFound): true
```

Notice `ErrNotFound` is wrapped with `%w` (Level 16), not just embedded as text - callers use `errors.Is` to check for it, which matters once this error needs to cross the Service and Handler layers later in this level.

---

## The Service Layer

The Service layer holds **business rules the repository shouldn't know about**. A map or a SQL table will happily store an empty name or a malformed email - it has no concept of "valid." Deciding what's valid is the Service's job.

Critically, the Service depends on the `UserRepository` **interface**, never on `InMemoryUserRepository` directly:

```go
var (
    ErrInvalidName  = errors.New("name must not be empty")
    ErrInvalidEmail = errors.New("email must contain @")
)

// UserService holds the business rules and depends on the UserRepository
// INTERFACE, never on InMemoryUserRepository directly. That's what lets
// main() and tests hand it completely different concrete repositories
// without changing a line of this type.
type UserService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

// CreateUser validates BEFORE ever calling the repository. An invalid
// user never reaches storage - the repository is not even asked.
func (s *UserService) CreateUser(u User) error {
    if strings.TrimSpace(u.Name) == "" {
        return fmt.Errorf("create user: %w", ErrInvalidName)
    }
    if !strings.Contains(u.Email, "@") {
        return fmt.Errorf("create user: %w", ErrInvalidEmail)
    }
    return s.repo.Create(u)
}
```

### Proving the Repository Is Never Touched

"Validates before touching the repository" is a claim that needs proof, not just a comment. A **poison repository** - a fake whose every method panics - proves it directly: if the Service ever reached storage with invalid input, the program would crash instead of printing a clean error.

```go
type poisonRepository struct{}

func (poisonRepository) Create(User) error {
    panic("poisonRepository.Create called - service should have rejected this input first")
}
// Get, Update, Delete panic the same way.

var _ UserRepository = poisonRepository{}
```

Real captured output:

```
=== Valid User: Reaches the Real Repository ===
error: <nil>

=== Invalid Name: Rejected Before Touching the Repository ===
error: create user: name must not be empty
errors.Is(err, ErrInvalidName): true
(no panic above - the poison repository was never called)

=== Invalid Email: Rejected Before Touching the Repository ===
error: create user: email must contain @
errors.Is(err, ErrInvalidEmail): true
(no panic above - the poison repository was never called)
```

### Rules That Need the Repository: Uniqueness

Not every business rule can be checked in isolation. "This email must be unique" requires looking at *existing* records - a rule the repository still shouldn't own, because "what counts as a duplicate" is a business decision, not a storage concern:

```go
var ErrEmailTaken = errors.New("email already registered")

func (s *UserService) CreateUser(u User) error {
    if strings.TrimSpace(u.Name) == "" {
        return fmt.Errorf("create user: %w", ErrInvalidName)
    }
    if !strings.Contains(u.Email, "@") {
        return fmt.Errorf("create user: %w", ErrInvalidEmail)
    }
    for _, existing := range s.repo.List() {
        if existing.Email == u.Email {
            return fmt.Errorf("create user: %w", ErrEmailTaken)
        }
    }
    return s.repo.Create(u)
}
```

Real captured output:

```
=== Create Alice ===
error: <nil>

=== Create Bob With Alice's Email (Rejected) ===
error: create user: email already registered
errors.Is(err, ErrEmailTaken): true

=== Create Bob With His Own Email (Accepted) ===
error: <nil>
```

The repository's `List()` method just lists what's there - it has no idea "uniqueness" is being checked against its results. That decision lives entirely in the Service.

---

## The Handler Layer

The Handler layer is deliberately **thin**: decode the request, call the Service, encode the response. No validation, no storage logic - both already live in the layers below it. This level uses Level 27's stdlib `net/http` patterns, including the method-and-path-parameter routing built into `http.ServeMux` since Go 1.22 (`"GET /users/{id}"`).

```go
type UserHandler struct {
    service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
    return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var u User
    if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
        http.Error(w, "invalid JSON body", http.StatusBadRequest)
        return
    }

    if err := h.service.CreateUser(u); err != nil {
        switch {
        case errors.Is(err, ErrInvalidName), errors.Is(err, ErrInvalidEmail):
            http.Error(w, err.Error(), http.StatusBadRequest)
        default:
            http.Error(w, "internal server error", http.StatusInternalServerError)
        }
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(u)
}

func NewMux(h *UserHandler) *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("POST /users", h.CreateUser)
    mux.HandleFunc("GET /users/{id}", h.GetUser)
    return mux
}
```

Because the Handler is just a function that takes an `http.ResponseWriter` and an `*http.Request`, it can be tested with `httptest` (Level 26's preview, paid off in full here) - no real network socket, no real server process:

```go
func TestCreateUserHandler_Success(t *testing.T) {
    mux := newTestMux()

    body, _ := json.Marshal(User{ID: "1", Name: "Alice", Email: "alice@example.com"})
    req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
    rec := httptest.NewRecorder()

    mux.ServeHTTP(rec, req)

    if rec.Code != http.StatusCreated {
        t.Fatalf("status = %d; want %d, body = %s", rec.Code, http.StatusCreated, rec.Body.String())
    }
    // ...decode rec.Body and assert on the fields
}
```

Real captured output:

```
=== RUN   TestCreateUserHandler_Success
--- PASS: TestCreateUserHandler_Success (0.00s)
=== RUN   TestCreateUserHandler_InvalidName
--- PASS: TestCreateUserHandler_InvalidName (0.00s)
=== RUN   TestGetUserHandler_Success
--- PASS: TestGetUserHandler_Success (0.00s)
=== RUN   TestGetUserHandler_NotFound
--- PASS: TestGetUserHandler_NotFound (0.00s)
PASS
ok  	level30.example/exercise4	1.049s
```

---

## Dependency Direction

The three layers only make sense if dependencies flow **one direction**:

```
Handler  ──depends on──>  Service  ──depends on──>  Repository (interface)
   │                          │                             ▲
   │                          │                             │
   ▼                          ▼                             │
uses *UserService       uses UserRepository       implemented by
  (concrete)                (interface)         InMemoryUserRepository
                                                      (concrete)
```

- The **Handler** knows about the **Service** (a concrete type it calls methods on).
- The **Service** knows about the **Repository interface** (never the concrete implementation).
- The **Repository implementation** knows about neither the Service nor the Handler - it just implements the interface.

Interfaces flow *against* the dependency arrow: the Repository's interface is defined by (and for) the Service that consumes it, then an implementation satisfies it from the other direction. Nothing above the Handler depends on anything below the Repository, and nothing below the Repository knows layers above it exist. If you ever find a Repository importing an HTTP package, or a Service importing `net/http`, the direction has broken.

---

## Wiring It Together in main()

Somewhere, something has to actually construct the concrete pieces and hand them to each other. That "somewhere" is `main()`:

```go
// This is the ONLY place that knows about the concrete repository type.
// Every layer above only ever sees an interface. Constructing the real
// dependencies here and handing them down - Repository into Service,
// Service into Handler - is manual dependency injection: Level 31 gives
// this pattern a name and shows how to scale it beyond "wire it by hand".
func main() {
    repo := NewInMemoryUserRepository()  // concrete repository
    service := NewUserService(repo)      // service depends on the INTERFACE
    handler := NewUserHandler(service)    // handler depends on the service
    mux := NewMux(handler)

    listener, _ := net.Listen("tcp", "127.0.0.1:0")
    server := &http.Server{Handler: mux}
    go server.Serve(listener)
    defer server.Shutdown(context.Background())

    // ...issue real HTTP requests against listener.Addr()
}
```

Real captured output from running this exact program (the port number will differ on your machine - that's expected, it's chosen by the OS):

```
=== Wiring: Repository -> Service -> Handler -> http.Server ===
server listening on 127.0.0.1:49297

=== POST /users (real HTTP request over loopback) ===
status: 201
body: {"id":"1","name":"Alice","email":"alice@example.com"}

=== GET /users/1 (real HTTP request over loopback) ===
status: 200
body: {"id":"1","name":"Alice","email":"alice@example.com"}

=== GET /users/999 (not created - 404) ===
status: 404
body: get user 999: user not found
```

Passing a fully-constructed dependency into a constructor (`NewUserService(repo)`, `NewUserHandler(service)`) instead of having each type construct its own dependencies internally **is** dependency injection - you're already doing it. **Level 31: Dependency Injection** picks this exact pattern up and shows how to manage it as the dependency graph grows past three layers and a handful of constructors.

---

## Why This Enables Easy Testing

Because `UserService` depends on the `UserRepository` interface, tests can hand it a repository that does something a real one never would - like failing every single call - without touching a database, a file, or even the real `InMemoryUserRepository`.

```go
// ErrStorageUnavailable simulates a real-world failure: a database that's
// down, a disk that's full, a network partition - anything the SERVICE
// can't control but must still handle gracefully.
var ErrStorageUnavailable = errors.New("storage backend unavailable")

// alwaysErrorRepository is a fake UserRepository (Level 15/26's mocking
// pattern) that fails every call.
type alwaysErrorRepository struct{}

func (alwaysErrorRepository) Create(User) error       { return ErrStorageUnavailable }
func (alwaysErrorRepository) Get(string) (User, error) { return User{}, ErrStorageUnavailable }

var _ UserRepository = alwaysErrorRepository{}

func TestUserService_CreateUser_RepositoryFailureSurfaces(t *testing.T) {
    service := NewUserService(alwaysErrorRepository{})

    err := service.CreateUser(User{ID: "1", Name: "Alice", Email: "alice@example.com"})

    if !errors.Is(err, ErrStorageUnavailable) {
        t.Errorf("error = %v; want it to wrap ErrStorageUnavailable", err)
    }
}
```

Real captured output:

```
=== RUN   TestUserService_CreateUser_Success
--- PASS: TestUserService_CreateUser_Success (0.00s)
=== RUN   TestUserService_CreateUser_RepositoryFailureSurfaces
--- PASS: TestUserService_CreateUser_RepositoryFailureSurfaces (0.00s)
PASS
ok  	level30.example/exercise6	0.781s
```

No real storage was ever involved in proving the Service handles a storage failure correctly. This is the entire payoff of the pattern: business logic gets tested as business logic, independent of whatever happens to be storing the data today.

---

## Error Handling Across Layers

Errors need to be translated as they cross layer boundaries (bridging Level 16):

1. The **Repository** returns a storage-flavored error: `ErrNotFound`, wrapped with context (`fmt.Errorf("get user %s: %w", id, ErrNotFound)`).
2. The **Service** either passes it through, wraps it further, or turns it into its own business error (`ErrInvalidName`, `ErrEmailTaken`).
3. The **Handler** is the *only* layer that knows what an HTTP status code is - it inspects the error with `errors.Is` and picks a status.

```go
// statusFor is the single place that translates a layered error into an
// HTTP status code. It checks the MOST SPECIFIC errors first (validation,
// conflict, not-found) and falls back to 500 for anything it doesn't
// recognize - an unrecognized error must never leak as a 200.
func statusFor(err error) int {
    switch {
    case errors.Is(err, ErrInvalidName), errors.Is(err, ErrInvalidEmail):
        return http.StatusBadRequest // 400
    case errors.Is(err, ErrEmailTaken):
        return http.StatusConflict // 409
    case errors.Is(err, ErrNotFound):
        return http.StatusNotFound // 404
    default:
        return http.StatusInternalServerError // 500 - unexpected, e.g. a storage failure
    }
}
```

| Origin | Error | HTTP Status |
|--------|-------|--------------|
| Service (validation) | `ErrInvalidName` / `ErrInvalidEmail` | 400 Bad Request |
| Service (business rule) | `ErrEmailTaken` | 409 Conflict |
| Repository (via Service) | `ErrNotFound` | 404 Not Found |
| Anything else | unrecognized error | 500 Internal Server Error |

A table-driven `httptest` suite proves this end-to-end, including an injected repository failure that must fall through to 500 (real captured output):

```
=== RUN   TestErrorTranslationAcrossLayers
=== RUN   TestErrorTranslationAcrossLayers/valid_create_->_201
=== RUN   TestErrorTranslationAcrossLayers/invalid_name_->_400
=== RUN   TestErrorTranslationAcrossLayers/duplicate_email_->_409
=== RUN   TestErrorTranslationAcrossLayers/unexpected_storage_failure_->_500
=== RUN   TestErrorTranslationAcrossLayers/get_existing_->_200
=== RUN   TestErrorTranslationAcrossLayers/get_missing_->_404
--- PASS: TestErrorTranslationAcrossLayers (0.00s)
    --- PASS: TestErrorTranslationAcrossLayers/valid_create_->_201 (0.00s)
    --- PASS: TestErrorTranslationAcrossLayers/invalid_name_->_400 (0.00s)
    --- PASS: TestErrorTranslationAcrossLayers/duplicate_email_->_409 (0.00s)
    --- PASS: TestErrorTranslationAcrossLayers/unexpected_storage_failure_->_500 (0.00s)
    --- PASS: TestErrorTranslationAcrossLayers/get_existing_->_200 (0.00s)
    --- PASS: TestErrorTranslationAcrossLayers/get_missing_->_404 (0.00s)
PASS
ok  	level30.example/exercise7	0.658s
```

The key discipline: **`errors.Is` at every boundary**, never string-matching an error message. A repository can add more context to its error message freely, as long as it keeps wrapping the same sentinel error, and every layer above it keeps working.

---

## When This Is Overkill

Three layers, two interfaces, and a wiring step in `main()` is real ceremony. For a 40-line script that reads a CSV file once and prints a total, this pattern adds files and indirection without buying anything back - there's no business logic worth isolating, no storage technology that will ever change, and no test suite complex enough to need fakes.

The pattern earns its keep once **any** of these become true:

- **Business logic exists and matters.** Validation, uniqueness checks, computed fields, authorization rules - anything beyond "store what I was given."
- **Testability matters.** You want to test business rules without a real database, a real file system, or a real network call.
- **Storage might change or vary.** Today it's a map; tomorrow it's SQL (Level 29); in tests it's a fake. The Service shouldn't have to change for any of these.

A single `main.go` with one struct and no rules doesn't need a Repository, a Service, and a Handler - it needs a function. Reach for this pattern when the program has grown enough that "which part of this file is responsible for X?" stops having an obvious answer.

---

## Best Practices

### 1. Define Repository Interfaces in the SERVICE Package, Not the Repository Package

Per Level 15's Rule 3 ("define interfaces where they're used"), the interface the Service depends on belongs with the Service, not exported by the repository implementation:

```go
// ✅ Good - service.go defines exactly what IT needs
package service

type NoteRepository interface {
    Create(n Note) error
    Get(id string) (Note, error)
}

// ❌ Provider-side (not idiomatic Go) - repository package dictates the
// shape every consumer must use, whether they need all of it or not
package repository

type Repository interface { /* ... */ }
```

### 2. Keep Handlers Free of Business Logic

A Handler that validates input itself, checks uniqueness itself, or computes a derived field itself has quietly absorbed the Service's job. If you can imagine calling the same logic from a CLI command or a message queue consumer instead of HTTP, it belongs in the Service.

### 3. Keep Services Free of HTTP and Storage Details

A Service that imports `net/http` or knows what a `sql.Tx` is has leaked a detail from the layer above or below it. The Service should be testable with nothing but plain Go values - no request, no response writer, no database handle.

### 4. Accept Interfaces, Return Concrete Types (Level 15)

`NewUserService(repo UserRepository)` accepts an interface; it returns `*UserService`, a concrete type. Callers get full access to everything `UserService` offers, while the constructor stays flexible about what it's handed.

### 5. Verify Implicit Satisfaction With a Compile-Time Check

`var _ UserRepository = (*InMemoryUserRepository)(nil)` costs nothing at runtime and turns a drifted method signature into a compile error instead of a runtime surprise.

---

## Common Mistakes

### Mistake 1: Services Depending on Concrete Repository Types

```go
// ❌ WRONG - welded to one repository forever
type UserService struct {
    repo *InMemoryUserRepository
}

// ✅ RIGHT - works with any UserRepository implementation
type UserService struct {
    repo UserRepository
}
```

A service built against a concrete type cannot accept a different repository implementation, a decorator (like a caching wrapper), or a test fake without editing the service itself - defeating the entire purpose of the pattern.

### Mistake 2: Business Logic Leaking Into Handlers

```go
// ❌ WRONG - validation belongs in the Service, not the Handler
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var u User
    json.NewDecoder(r.Body).Decode(&u)
    if strings.TrimSpace(u.Name) == "" {
        http.Error(w, "name required", http.StatusBadRequest)
        return
    }
    h.repo.Create(u) // also wrong - skips the Service entirely
}

// ✅ RIGHT - Handler decodes and delegates, Service decides what's valid
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var u User
    json.NewDecoder(r.Body).Decode(&u)
    if err := h.service.CreateUser(u); err != nil {
        http.Error(w, err.Error(), statusFor(err))
        return
    }
    // ...
}
```

### Mistake 3: Repositories That Know About HTTP Status Codes

```go
// ❌ WRONG - the Repository layer has no business knowing what a 404 is
func (r *InMemoryUserRepository) Get(id string) (User, int, error) {
    // ...
    return User{}, http.StatusNotFound, ErrNotFound
}

// ✅ RIGHT - the Repository returns a plain error; the HANDLER decides
// what status code that error maps to
func (r *InMemoryUserRepository) Get(id string) (User, error) {
    // ...
    return User{}, ErrNotFound
}
```

A Repository that returns HTTP concepts can never be reused by a CLI tool, a gRPC service, or a background job - only by something that already speaks HTTP.

### Mistake 4: Skipping the Service and Calling the Repository Directly From Handlers

```go
// ❌ WRONG - no validation, no uniqueness check, nowhere for business
// rules to live; the Handler is now also the Service
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var u User
    json.NewDecoder(r.Body).Decode(&u)
    h.repo.Create(u)
}

// ✅ RIGHT - the Handler only ever talks to the Service
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var u User
    json.NewDecoder(r.Body).Decode(&u)
    if err := h.service.CreateUser(u); err != nil {
        http.Error(w, err.Error(), statusFor(err))
        return
    }
    // ...
}
```

Even a Service that starts out as a thin pass-through is worth keeping, because it's the one place business rules can be added later without touching the Handler.

---

## Summary

**The Three Layers:**
- **Repository** - data access only (an interface + a concrete implementation, in-memory in this level)
- **Service** - business rules only (validation, uniqueness, computed fields) - depends on the Repository *interface*
- **Handler** - HTTP translation only - depends on the Service, decodes requests, encodes responses

**Dependency Direction:**
- Handler → Service → Repository (interface); implementations satisfy the interface from the other direction
- Nothing below a layer knows that layer, or anything above it, exists

**Testing Payoff:**
- Swap the real repository for a fake (a poison repository, an always-erroring repository) to test the Service in complete isolation
- Test the Handler with `httptest` against a real or fake Service, no network required

**Error Translation:**
- Repository errors (wrapped sentinels) → Service errors (passed through or added to) → HTTP status codes (Handler, via `errors.Is`)

**When to Use It:**
- Once business logic, testability, or swappable storage genuinely matter - not for a 40-line script

---

## Next Steps

You now understand:
- ✅ Why splitting Repository, Service, and Handler into three layers with one job each pays off
- ✅ How to define a small Repository interface and satisfy it with an in-memory implementation
- ✅ How to write Service logic that depends on interfaces, not concrete types, and enforces rules the repository can't
- ✅ How to keep Handlers thin and testable with `httptest`
- ✅ How to wire all three layers together by hand in `main()` - manual dependency injection
- ✅ How to translate errors across layer boundaries into the right HTTP status codes
- ✅ When this pattern is worth its ceremony, and when it isn't

**Next level:** Level 31 - Dependency Injection
- Naming and formalizing the "construct dependencies and hand them down" pattern you just did by hand in `main()`
- Managing a dependency graph as it grows past three layers
- Constructor injection, wiring functions, and where DI frameworks fit into Go (spoiler: rarely)

You're now structuring programs the way real production Go services are built. Keep going! 🚀
