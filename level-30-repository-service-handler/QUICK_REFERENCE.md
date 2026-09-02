# Level 30: Quick Reference Card

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir myapp && cd myapp
go mod init myapp

# Create main.go
cat > main.go << 'EOF'
package main

import (
    "errors"
    "fmt"
    "sync"
)

type User struct {
    ID, Name, Email string
}

var ErrNotFound = errors.New("user not found")

// Repository: interface + in-memory implementation
type UserRepository interface {
    Create(u User) error
    Get(id string) (User, error)
}

type InMemoryUserRepository struct {
    mu    sync.Mutex
    users map[string]User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
    return &InMemoryUserRepository{users: make(map[string]User)}
}

func (r *InMemoryUserRepository) Create(u User) error {
    r.mu.Lock(); defer r.mu.Unlock()
    r.users[u.ID] = u
    return nil
}

func (r *InMemoryUserRepository) Get(id string) (User, error) {
    r.mu.Lock(); defer r.mu.Unlock()
    u, ok := r.users[id]
    if !ok {
        return User{}, fmt.Errorf("get user %s: %w", id, ErrNotFound)
    }
    return u, nil
}

// Service: depends on the INTERFACE, enforces rules
type UserService struct{ repo UserRepository }

func NewUserService(repo UserRepository) *UserService { return &UserService{repo: repo} }

func (s *UserService) CreateUser(u User) error {
    if u.Name == "" {
        return errors.New("name required")
    }
    return s.repo.Create(u)
}

func main() {
    service := NewUserService(NewInMemoryUserRepository())
    fmt.Println(service.CreateUser(User{ID: "1", Name: "Alice"}))
}
EOF

# Run
go run main.go
```

---

## 📋 The Three Layers

```go
// Repository - data access ONLY
type UserRepository interface {
    Create(u User) error
    Get(id string) (User, error)
    Update(u User) error
    Delete(id string) error
}

// Service - business rules ONLY, depends on the INTERFACE
type UserService struct {
    repo UserRepository // never the concrete type
}

// Handler - HTTP translation ONLY, depends on the Service
type UserHandler struct {
    service *UserService
}
```

---

## 🔀 Dependency Direction

```
Handler  ──depends on──>  Service  ──depends on──>  Repository (interface)
                                                             ▲
                                              implemented by │
                                                InMemoryUserRepository
```

Interface defined where it's USED (in the Service's package), not where the implementation lives.

---

## ✅ Compile-Time Interface Check

```go
var _ UserRepository = (*InMemoryUserRepository)(nil)
```

---

## 🧪 Testing With Fakes

```go
// Fails every call - proves the Service surfaces errors correctly
type alwaysErrorRepository struct{}
func (alwaysErrorRepository) Create(User) error        { return errStorageDown }
func (alwaysErrorRepository) Get(string) (User, error) { return User{}, errStorageDown }

// Panics on every call - proves validation runs BEFORE the repository
type poisonRepository struct{}
func (poisonRepository) Create(User) error { panic("should not be called") }

service := NewUserService(alwaysErrorRepository{}) // or poisonRepository{}
```

---

## 🌐 Handler + httptest

```go
mux := NewMux(NewUserHandler(NewUserService(NewInMemoryUserRepository())))

req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
rec := httptest.NewRecorder()
mux.ServeHTTP(rec, req)

// rec.Code, rec.Body.String()
```

```go
// Go 1.22+ ServeMux: method + path parameter routing, no third-party router
mux.HandleFunc("POST /users", h.CreateUser)
mux.HandleFunc("GET /users/{id}", h.GetUser)
id := r.PathValue("id")
```

---

## 🔀 Error Translation Table

```go
func statusFor(err error) int {
    switch {
    case errors.Is(err, ErrInvalidName), errors.Is(err, ErrInvalidEmail):
        return http.StatusBadRequest // 400
    case errors.Is(err, ErrEmailTaken):
        return http.StatusConflict // 409
    case errors.Is(err, ErrNotFound):
        return http.StatusNotFound // 404
    default:
        return http.StatusInternalServerError // 500
    }
}
```

| Layer | Error | Status |
|-------|-------|--------|
| Service (validation) | `ErrInvalidName`/`ErrInvalidEmail` | 400 |
| Service (business rule) | `ErrEmailTaken` | 409 |
| Repository (via Service) | `ErrNotFound` | 404 |
| Anything unrecognized | - | 500 |

---

## 🔌 Wiring in main() (Level 31 Preview)

```go
func main() {
    repo := NewInMemoryUserRepository()   // concrete
    service := NewUserService(repo)       // interface injected
    handler := NewUserHandler(service)     // service injected
    mux := NewMux(handler)

    http.ListenAndServe(":8080", mux)
}
```

---

## 📦 Multi-Package Layout (Comprehensive Exercise)

```
myapp/
├── go.mod
├── main.go            (wires repository -> service -> handler)
├── repository/
│   └── repository.go  (Note, ErrNotFound, InMemoryNoteRepository)
├── service/
│   ├── service.go      (NoteRepository interface, NoteService)
│   └── service_test.go (tests Service with a fake repository)
└── handler/
    ├── handler.go       (NoteService interface, NoteHandler, NewMux)
    └── handler_test.go  (tests Handler with a fake service)
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Concrete repository dependency | `repo *InMemoryUserRepository` | `repo UserRepository` |
| Validation in the Handler | checking `u.Name == ""` in the handler func | delegate to `service.CreateUser(u)` |
| Repository knows HTTP | `Get(id) (User, int, error)` | `Get(id) (User, error)` |
| Handler skips the Service | `h.repo.Create(u)` from the handler | `h.service.CreateUser(u)` |
| String-matching errors | `err.Error() == "not found"` | `errors.Is(err, ErrNotFound)` |

---

## 🎓 Before Next Level

Can you:
- [ ] Define a Repository interface and satisfy it with an in-memory implementation?
- [ ] Write a Service that rejects invalid input before touching the repository?
- [ ] Write a thin Handler and test it with `httptest`?
- [ ] Explain the Handler → Service → Repository dependency direction?
- [ ] Wire all three layers together by hand in `main()`?
- [ ] Swap in a fake repository to test a Service in isolation?
- [ ] Translate a not-found error into a 404 across two layers?

If YES → You're ready for Level 31!

---

## 📚 Next Level

Level 31: Dependency Injection
- Naming and scaling the wiring pattern from this level
- Managing a dependency graph as it grows
- Constructor injection and wiring functions

You've got layered architecture down! 💪
