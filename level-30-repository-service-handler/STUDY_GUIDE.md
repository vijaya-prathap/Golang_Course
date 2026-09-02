# Level 30: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: The Three Layers
```
Day 1:  What the pattern is and why - Repository, Service, Handler
Day 2:  The Repository layer - interface + in-memory implementation
Day 3:  The Service layer - validation rules the repository can't enforce
Day 4:  The Service layer - rules that need repository data (uniqueness)
Day 5:  The Handler layer - thin HTTP translation, tested with httptest
Day 6:  Dependency direction and wiring in main()
Day 7:  Review & practice
```

### Week 2: Testing & Real-World Practice
```
Day 1:  Exercises 1-3 (repository, validation, uniqueness)
Day 2:  Exercises 4-5 (handler + httptest, full wiring end-to-end)
Day 3:  Exercises 6-7 (fake repository, error translation to 404/500)
Day 4:  Exercises 8-9 (dependency direction proof, full CRUD API)
Day 5:  Exercise 10 (comprehensive: multi-package Notes API)
Day 6-7: Bonus challenges & consolidation
```

---

## 🏗️ The Three-Layer Architecture

```
┌───────────────┐        ┌───────────────┐        ┌───────────────────┐
│    Handler     │  ───>  │    Service     │  ───>  │  Repository        │
│  (transport)   │        │ (business      │        │  (interface)       │
│                │        │   rules)       │        │        ▲            │
│  net/http:     │        │                │        │        │            │
│  decode JSON,  │        │ validation,    │        │  implemented by     │
│  call Service, │        │ uniqueness,    │        │        │            │
│  encode JSON   │        │ computed data  │        │ InMemoryUserRepo    │
└───────────────┘        └───────────────┘        └───────────────────┘
       │                         │                          │
       └── depends on ───────────┘                          │
                    concrete *UserService          depends on the
                                                    UserRepository INTERFACE,
                                                    never the concrete type
```

**Read the arrows as "depends on," not "calls data to."** The Handler depends on the Service; the Service depends on the Repository *interface*. The concrete `InMemoryUserRepository` satisfies that interface from the opposite direction - it doesn't know the Service or Handler exist at all.

```
Nothing below a layer knows that layer exists:
  Repository doesn't know Service exists
  Service doesn't know Handler exists

Nothing above a layer knows HOW the layer below it works:
  Handler doesn't know how Service validates
  Service doesn't know if storage is a map, SQL, or a file
```

---

## 🧩 Interface Ownership: Defined Where It's Used

```
┌─────────────────────────────────────────────────────────────────┐
│  package service              (the CONSUMER)                     │
│                                                                     │
│      type UserRepository interface {   ← interface lives HERE,     │
│          Create(u User) error            next to the code that      │
│          Get(id string) (User, error)    actually CALLS it           │
│      }                                                                 │
│                                                                          │
│      type UserService struct { repo UserRepository }                    │
└─────────────────────────────────────────────────────────────────────────┘
                                    ▲
                     satisfies (implicitly, Level 15)
                                    │
┌─────────────────────────────────────────────────────────────────────────┐
│  package repository            (the PROVIDER)                            │
│                                                                             │
│      type InMemoryUserRepository struct { ... }                            │
│      func (r *InMemoryUserRepository) Create(u User) error { ... }          │
│      func (r *InMemoryUserRepository) Get(id string) (User, error) { ... }   │
│      // No interface declared here at all - the repository just HAS         │
│      // the right methods.                                                    │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Why this matters:** the consumer declares exactly the contract it needs. A test fake, a second repository implementation, or a future SQL-backed store can all satisfy `UserRepository` without ever importing anything from a "contracts" package - they just need the right method signatures.

---

## 🧪 Testing With Fakes at Every Layer

```
┌────────────────────────────────────────────────────────────────────┐
│                    type UserRepository interface { ... }              │
│                            ▲                    ▲                       │
│              satisfies     │                    │  satisfies             │
│         ┌──────────────────┘                    └───────────────┐         │
│         │                                                         │         │
│  InMemoryUserRepository (production)              alwaysErrorRepository    │
│  real map-backed storage                          (test) fails every call   │
│                                                                                │
│  UserService depends on the INTERFACE, never on either concrete type:          │
│                                                                                   │
│      NewUserService(NewInMemoryUserRepository())   ← production                  │
│      NewUserService(alwaysErrorRepository{})        ← test: proves error handling │
│      NewUserService(poisonRepository{})             ← test: proves validation runs │
│                                                        BEFORE the repository is used │
└────────────────────────────────────────────────────────────────────────────────────┘
```

```
Same pattern, one layer up:

  type NoteService interface { ... }      ← owned by the HANDLER package
          ▲                    ▲
          │                    │
  *service.NoteService    fakeNoteService
  (production)            (test - returns canned results, no real repository)
```

Every layer's dependency is an interface, so every layer can be tested with a fake standing in for the layer below it - no real storage, no real HTTP server, no framework required (Level 26's mocking-via-interfaces, applied at every boundary).

---

## 🔀 Error Translation Across Layers

```
┌──────────────┐   wraps with %w   ┌──────────────┐   errors.Is()    ┌──────────────┐
│  Repository   │ ────────────────> │   Service     │ ───────────────> │   Handler     │
│               │                    │               │                   │               │
│ ErrNotFound   │                    │ passes through│                   │ statusFor(err) │
│ (storage-     │                    │ OR returns    │                   │ picks the      │
│  flavored)    │                    │ its own error │                   │ HTTP status    │
│               │                    │ (ErrInvalid*, │                   │                │
│               │                    │  ErrEmailTaken)│                   │               │
└──────────────┘                    └──────────────┘                   └──────────────┘
```

```
statusFor(err error) int  - the ONE function that knows about status codes

  ErrInvalidName / ErrInvalidEmail  ──────────────>  400 Bad Request
  ErrEmailTaken                     ──────────────>  409 Conflict
  ErrNotFound                       ──────────────>  404 Not Found
  anything unrecognized             ──────────────>  500 Internal Server Error
                                                       (never leaks as 200!)
```

Real captured output proving all four branches, end-to-end through `httptest`:

```
=== RUN   TestErrorTranslationAcrossLayers
=== RUN   TestErrorTranslationAcrossLayers/valid_create_->_201
=== RUN   TestErrorTranslationAcrossLayers/invalid_name_->_400
=== RUN   TestErrorTranslationAcrossLayers/duplicate_email_->_409
=== RUN   TestErrorTranslationAcrossLayers/unexpected_storage_failure_->_500
=== RUN   TestErrorTranslationAcrossLayers/get_existing_->_200
=== RUN   TestErrorTranslationAcrossLayers/get_missing_->_404
--- PASS: TestErrorTranslationAcrossLayers (0.00s)
PASS
ok  	level30.example/exercise7	0.658s
```

**Rule:** always check with `errors.Is`, never by comparing error message strings. A repository can add more context to its wrapped error (`fmt.Errorf("get user %s: %w", id, ErrNotFound)`) without breaking a single `errors.Is(err, ErrNotFound)` check anywhere above it.

---

## 🎯 Dependency Injection Preview (Level 31)

```
func main() {
    repo := NewInMemoryUserRepository()   // 1. build the concrete dependency
    service := NewUserService(repo)       // 2. inject it into the next layer up
    handler := NewUserHandler(service)     // 3. inject THAT into the next layer up
    mux := NewMux(handler)
}
```

This is manual dependency injection - nothing more than passing constructed values into constructors. Level 31 gives it a name, shows how it scales as the dependency graph grows, and covers where (and whether) a DI framework earns its place in Go.

---

## 🚨 Common Mistakes

### Mistake 1: Service Depends on a Concrete Repository Type

```go
// ❌ WRONG - welded to one repository forever
type UserService struct { repo *InMemoryUserRepository }

// ✅ RIGHT - works with any UserRepository implementation
type UserService struct { repo UserRepository }
```

### Mistake 2: Business Logic Leaking Into the Handler

```go
// ❌ WRONG - validation belongs in the Service
if strings.TrimSpace(u.Name) == "" {
    http.Error(w, "name required", http.StatusBadRequest)
    return
}

// ✅ RIGHT - Handler delegates, Service decides what's valid
if err := h.service.CreateUser(u); err != nil {
    http.Error(w, err.Error(), statusFor(err))
    return
}
```

### Mistake 3: Repository Knows About HTTP Status Codes

```go
// ❌ WRONG - the Repository has no business knowing what a 404 is
func (r *InMemoryUserRepository) Get(id string) (User, int, error)

// ✅ RIGHT - a plain error; the HANDLER decides the status code
func (r *InMemoryUserRepository) Get(id string) (User, error)
```

### Mistake 4: Handler Calls the Repository Directly, Skipping the Service

```go
// ❌ WRONG - no validation, no business rules, ever
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var u User
    json.NewDecoder(r.Body).Decode(&u)
    h.repo.Create(u)
}

// ✅ RIGHT - the Handler only ever talks to the Service
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var u User
    json.NewDecoder(r.Body).Decode(&u)
    h.service.CreateUser(u)
}
```

---

## 📈 Progression Summary

### Understanding Level 30

Level 30 teaches how to split a program into layers with one job each:

1. **Repository** — data access only, an interface + an in-memory implementation
2. **Service** — business rules only, depends on the Repository interface
3. **Handler** — HTTP translation only, depends on the Service
4. **Dependency direction** — Handler → Service → Repository interface, never the reverse
5. **Testing payoff** — fakes at every boundary, no real storage or network required
6. **Error translation** — sentinel errors, wrapped with `%w`, mapped to HTTP status via `errors.Is`

### Prerequisites for Level 31

Before moving to Level 31 (Dependency Injection), you need:

- ✅ Comfortable defining a small interface and satisfying it implicitly (Level 15)
- ✅ Comfortable writing table-driven tests and fakes (Level 26)
- ✅ Comfortable with `net/http` handlers and `httptest` (Level 27)
- ✅ Can explain why Repository, Service, and Handler each have exactly one job
- ✅ Can explain the dependency direction and why interfaces point "backwards" against it
- ✅ Have manually wired a repository → service → handler chain in `main()`

### Ready for Level 31?

Level 31 takes the manual wiring you did by hand in Exercises 5 and 10 and turns it into a repeatable pattern:
- Constructor injection at scale, once there are more than three layers
- Wiring functions that assemble a whole dependency graph in one place
- Where (and whether) a DI framework or code generator earns its place in Go

---

## ✅ Checklist Before Level 31

- [ ] Can define a Repository interface with the methods a Service actually needs
- [ ] Can implement that interface with an in-memory (map or slice) store
- [ ] Can write a Service that validates input before ever calling the repository
- [ ] Can write a Service rule that needs repository data (like uniqueness)
- [ ] Can write a thin Handler and test it with `httptest`, no real server
- [ ] Can explain the Handler → Service → Repository dependency direction from memory
- [ ] Can wire all three layers together by hand in `main()`
- [ ] Can swap in a fake repository to test a Service in isolation
- [ ] Can translate a not-found error into a 404 through two layers
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The Split
Repository (data), Service (rules), Handler (transport) - one job each, so each can be tested and changed independently.

### The Direction
Handler → Service → Repository interface. Interfaces are owned by the layer that USES them, not the layer that implements them.

### The Payoff
Fakes at every boundary mean every layer is testable without a real database, a real file, or a real network call.

### The Discipline
`errors.Is` at every boundary, never string-matching. The Handler is the only layer that knows what an HTTP status code is.

### The Honesty
This pattern is ceremony a tiny script doesn't need. It earns its keep once business logic, testability, or swappable storage genuinely matter.

---

## 📚 Next Level

Level 31: Dependency Injection
- Naming and scaling the "construct dependencies and hand them down" pattern from this level
- Managing a dependency graph as it grows past three layers
- Constructor injection, wiring functions, and where DI frameworks fit into Go

You've got layered architecture down! Keep going! 🚀
