# Level 30: Repository-Service-Handler - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge. Every exercise here runs 100% offline - all storage is in-memory, and any HTTP traffic stays on `127.0.0.1` (your own machine talking to itself).

---

## Exercise 1: The Repository Layer - Interface + In-Memory Implementation

**Objective:** Define a small Repository interface and satisfy it with an in-memory implementation

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level30-exercise1
cd ~/projects/level30-exercise1
go mod init level30.example/exercise1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "errors"
    "fmt"
    "sync"
)

// User is the record every example in this level works with.
type User struct {
    ID    string
    Name  string
    Email string
}

// ErrNotFound is returned by a repository when no record matches the given ID.
var ErrNotFound = errors.New("user not found")

// UserRepository is the Repository layer's contract: pure data access,
// nothing else. It knows nothing about validation, HTTP, or business
// rules - just how to store and retrieve a User by ID.
//
// Any type with these four methods satisfies this interface implicitly
// (Level 15) - this level uses an in-memory map, but a SQL-backed store
// (Level 29) or anything else could satisfy it identically.
type UserRepository interface {
    Create(u User) error
    Get(id string) (User, error)
    Update(u User) error
    Delete(id string) error
}

// InMemoryUserRepository is one concrete implementation of UserRepository,
// backed by a map guarded by a mutex (Level 23) so it stays safe to use
// concurrently once HTTP requests start hitting it in later exercises.
type InMemoryUserRepository struct {
    mu    sync.Mutex
    users map[string]User
}

// NewInMemoryUserRepository builds an empty, ready-to-use repository.
func NewInMemoryUserRepository() *InMemoryUserRepository {
    return &InMemoryUserRepository{users: make(map[string]User)}
}

func (r *InMemoryUserRepository) Create(u User) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, exists := r.users[u.ID]; exists {
        return fmt.Errorf("create user %s: already exists", u.ID)
    }
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

func (r *InMemoryUserRepository) Update(u User) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, ok := r.users[u.ID]; !ok {
        return fmt.Errorf("update user %s: %w", u.ID, ErrNotFound)
    }
    r.users[u.ID] = u
    return nil
}

func (r *InMemoryUserRepository) Delete(id string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, ok := r.users[id]; !ok {
        return fmt.Errorf("delete user %s: %w", id, ErrNotFound)
    }
    delete(r.users, id)
    return nil
}

// Compile-time check: InMemoryUserRepository must satisfy UserRepository.
// This line does nothing at runtime - if the methods above ever drift
// from the interface, this fails to COMPILE instead of failing silently.
var _ UserRepository = (*InMemoryUserRepository)(nil)

func main() {
    var repo UserRepository = NewInMemoryUserRepository()

    fmt.Println("=== Create ===")
    if err := repo.Create(User{ID: "1", Name: "Alice", Email: "alice@example.com"}); err != nil {
        fmt.Println("unexpected error:", err)
    }
    if err := repo.Create(User{ID: "2", Name: "Bob", Email: "bob@example.com"}); err != nil {
        fmt.Println("unexpected error:", err)
    }
    fmt.Println("created users 1 and 2")

    fmt.Println("\n=== Get ===")
    u, err := repo.Get("1")
    if err != nil {
        fmt.Println("unexpected error:", err)
    }
    fmt.Printf("found: %+v\n", u)

    fmt.Println("\n=== Update ===")
    u.Name = "Alice Updated"
    if err := repo.Update(u); err != nil {
        fmt.Println("unexpected error:", err)
    }
    updated, _ := repo.Get("1")
    fmt.Printf("after update: %+v\n", updated)

    fmt.Println("\n=== Delete ===")
    if err := repo.Delete("2"); err != nil {
        fmt.Println("unexpected error:", err)
    }
    fmt.Println("deleted user 2")

    fmt.Println("\n=== Get a Missing User ===")
    _, err = repo.Get("2")
    if err != nil {
        fmt.Println("error:", err)
        fmt.Println("errors.Is(err, ErrNotFound):", errors.Is(err, ErrNotFound))
    }
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

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

**Learning Objectives:**
- ✅ Define a small Repository interface (Get/Create/Update/Delete)
- ✅ Satisfy it implicitly with an in-memory implementation (Level 15 revisited)
- ✅ Verify satisfaction with a compile-time `var _ Interface = (*Type)(nil)` check
- ✅ Wrap a not-found error with `%w` so `errors.Is` can detect it later (Level 16)

---

## Exercise 2: The Service Layer - Validation Before Storage

**Objective:** Write a Service that enforces a rule the repository doesn't know about, and prove it never touches the repository for invalid input

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level30-exercise2
cd ~/projects/level30-exercise2
go mod init level30.example/exercise2
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "errors"
    "fmt"
    "strings"
    "sync"
)

// --- Repository layer (same as Exercise 1) ---

type User struct {
    ID    string
    Name  string
    Email string
}

var ErrNotFound = errors.New("user not found")

type UserRepository interface {
    Create(u User) error
    Get(id string) (User, error)
    Update(u User) error
    Delete(id string) error
}

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

func (r *InMemoryUserRepository) Update(u User) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.users[u.ID] = u
    return nil
}

func (r *InMemoryUserRepository) Delete(id string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    delete(r.users, id)
    return nil
}

var _ UserRepository = (*InMemoryUserRepository)(nil)

// --- Service layer ---

// Sentinel errors for rules the REPOSITORY has no concept of. A map or a
// SQL table doesn't know what a "valid" name or email looks like - that's
// a business rule, so it belongs here, not in the Repository layer.
var (
    ErrInvalidName  = errors.New("name must not be empty")
    ErrInvalidEmail = errors.New("email must contain @")
)

// UserService holds the business rules and depends on the UserRepository
// INTERFACE, never on InMemoryUserRepository directly.
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

// poisonRepository proves validation runs first: every method panics, so
// if the service ever reached the repository for invalid input, this
// program would crash instead of printing a clean error.
type poisonRepository struct{}

func (poisonRepository) Create(User) error {
    panic("poisonRepository.Create called - service should have rejected this input first")
}
func (poisonRepository) Get(string) (User, error) {
    panic("poisonRepository.Get called")
}
func (poisonRepository) Update(User) error {
    panic("poisonRepository.Update called")
}
func (poisonRepository) Delete(string) error {
    panic("poisonRepository.Delete called")
}

var _ UserRepository = poisonRepository{}

func main() {
    fmt.Println("=== Valid User: Reaches the Real Repository ===")
    realService := NewUserService(NewInMemoryUserRepository())
    err := realService.CreateUser(User{ID: "1", Name: "Alice", Email: "alice@example.com"})
    fmt.Println("error:", err)

    fmt.Println("\n=== Invalid Name: Rejected Before Touching the Repository ===")
    poisonService := NewUserService(poisonRepository{})
    err = poisonService.CreateUser(User{ID: "2", Name: "   ", Email: "bob@example.com"})
    fmt.Println("error:", err)
    fmt.Println("errors.Is(err, ErrInvalidName):", errors.Is(err, ErrInvalidName))
    fmt.Println("(no panic above - the poison repository was never called)")

    fmt.Println("\n=== Invalid Email: Rejected Before Touching the Repository ===")
    err = poisonService.CreateUser(User{ID: "3", Name: "Carol", Email: "not-an-email"})
    fmt.Println("error:", err)
    fmt.Println("errors.Is(err, ErrInvalidEmail):", errors.Is(err, ErrInvalidEmail))
    fmt.Println("(no panic above - the poison repository was never called)")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

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

**Learning Objectives:**
- ✅ Write a Service that depends on the Repository interface, not a concrete type
- ✅ Enforce a validation rule the repository has no concept of
- ✅ Prove rejection happens BEFORE storage, using a repository that panics if ever called

---

## Exercise 3: The Service Layer - A Rule That Needs the Repository (Uniqueness)

**Objective:** Enforce a business rule that requires querying existing data

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level30-exercise3
cd ~/projects/level30-exercise3
go mod init level30.example/exercise3
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "errors"
    "fmt"
    "strings"
    "sync"
)

// --- Repository layer ---

type User struct {
    ID    string
    Name  string
    Email string
}

var ErrNotFound = errors.New("user not found")

// UserRepository now also exposes List, which the Service layer needs to
// enforce email uniqueness - a rule the repository itself knows nothing
// about. The repository just answers "what's in storage?"
type UserRepository interface {
    Create(u User) error
    Get(id string) (User, error)
    Update(u User) error
    Delete(id string) error
    List() []User
}

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

func (r *InMemoryUserRepository) Update(u User) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.users[u.ID] = u
    return nil
}

func (r *InMemoryUserRepository) Delete(id string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    delete(r.users, id)
    return nil
}

func (r *InMemoryUserRepository) List() []User {
    r.mu.Lock()
    defer r.mu.Unlock()
    out := make([]User, 0, len(r.users))
    for _, u := range r.users {
        out = append(out, u)
    }
    return out
}

var _ UserRepository = (*InMemoryUserRepository)(nil)

// --- Service layer ---

var (
    ErrInvalidName  = errors.New("name must not be empty")
    ErrInvalidEmail = errors.New("email must contain @")
    ErrEmailTaken   = errors.New("email already registered")
)

type UserService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

// CreateUser validates the input, THEN checks uniqueness against
// existing records, and only then asks the repository to store it.
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

func main() {
    service := NewUserService(NewInMemoryUserRepository())

    fmt.Println("=== Create Alice ===")
    err := service.CreateUser(User{ID: "1", Name: "Alice", Email: "alice@example.com"})
    fmt.Println("error:", err)

    fmt.Println("\n=== Create Bob With Alice's Email (Rejected) ===")
    err = service.CreateUser(User{ID: "2", Name: "Bob", Email: "alice@example.com"})
    fmt.Println("error:", err)
    fmt.Println("errors.Is(err, ErrEmailTaken):", errors.Is(err, ErrEmailTaken))

    fmt.Println("\n=== Create Bob With His Own Email (Accepted) ===")
    err = service.CreateUser(User{ID: "2", Name: "Bob", Email: "bob@example.com"})
    fmt.Println("error:", err)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Create Alice ===
error: <nil>

=== Create Bob With Alice's Email (Rejected) ===
error: create user: email already registered
errors.Is(err, ErrEmailTaken): true

=== Create Bob With His Own Email (Accepted) ===
error: <nil>
```

**Learning Objectives:**
- ✅ Extend a Repository interface with a query method (`List`) a business rule needs
- ✅ Implement a uniqueness rule in the Service, using repository data without owning the rule
- ✅ Recognize that "what counts as a duplicate" is a business decision, not a storage one

---

## Exercise 4: The Handler Layer - Translating HTTP, Tested With httptest

**Objective:** Write a thin HTTP handler that delegates to the Service, and test it without a real server

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level30-exercise4
cd ~/projects/level30-exercise4
go mod init level30.example/exercise4
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "log"
    "net/http"
    "strings"
    "sync"
)

// --- Repository layer ---

type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

var ErrNotFound = errors.New("user not found")

type UserRepository interface {
    Create(u User) error
    Get(id string) (User, error)
    Update(u User) error
    Delete(id string) error
}

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

func (r *InMemoryUserRepository) Update(u User) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.users[u.ID] = u
    return nil
}

func (r *InMemoryUserRepository) Delete(id string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    delete(r.users, id)
    return nil
}

var _ UserRepository = (*InMemoryUserRepository)(nil)

// --- Service layer ---

var (
    ErrInvalidName  = errors.New("name must not be empty")
    ErrInvalidEmail = errors.New("email must contain @")
)

type UserService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) CreateUser(u User) error {
    if strings.TrimSpace(u.Name) == "" {
        return fmt.Errorf("create user: %w", ErrInvalidName)
    }
    if !strings.Contains(u.Email, "@") {
        return fmt.Errorf("create user: %w", ErrInvalidEmail)
    }
    return s.repo.Create(u)
}

func (s *UserService) GetUser(id string) (User, error) {
    return s.repo.Get(id)
}

// --- Handler layer ---
//
// UserHandler is deliberately thin: it only decodes requests, calls the
// Service, and encodes responses. It contains no validation and no
// storage logic - both already live in layers below it.
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

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")

    u, err := h.service.GetUser(id)
    if err != nil {
        switch {
        case errors.Is(err, ErrNotFound):
            http.Error(w, err.Error(), http.StatusNotFound)
        default:
            http.Error(w, "internal server error", http.StatusInternalServerError)
        }
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(u)
}

// NewMux wires the routes to handler methods - net/http's built-in
// ServeMux (Go 1.22+) supports method + path-parameter patterns directly,
// no third-party router required (Level 27's stdlib approach).
func NewMux(h *UserHandler) *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("POST /users", h.CreateUser)
    mux.HandleFunc("GET /users/{id}", h.GetUser)
    return mux
}

func main() {
    repo := NewInMemoryUserRepository()
    service := NewUserService(repo)
    handler := NewUserHandler(service)
    mux := NewMux(handler)

    log.Println("listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
EOF
```

3. Create `main_test.go`:

```bash
cat > main_test.go << 'EOF'
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

func newTestMux() *http.ServeMux {
    repo := NewInMemoryUserRepository()
    service := NewUserService(repo)
    handler := NewUserHandler(service)
    return NewMux(handler)
}

func TestCreateUserHandler_Success(t *testing.T) {
    mux := newTestMux()

    body, _ := json.Marshal(User{ID: "1", Name: "Alice", Email: "alice@example.com"})
    req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
    rec := httptest.NewRecorder()

    mux.ServeHTTP(rec, req)

    if rec.Code != http.StatusCreated {
        t.Fatalf("status = %d; want %d, body = %s", rec.Code, http.StatusCreated, rec.Body.String())
    }

    var got User
    if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
        t.Fatalf("decoding response body: %v", err)
    }
    if got.ID != "1" || got.Name != "Alice" {
        t.Errorf("got %+v; want ID=1 Name=Alice", got)
    }
}

func TestCreateUserHandler_InvalidName(t *testing.T) {
    mux := newTestMux()

    body, _ := json.Marshal(User{ID: "2", Name: "   ", Email: "bob@example.com"})
    req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
    rec := httptest.NewRecorder()

    mux.ServeHTTP(rec, req)

    if rec.Code != http.StatusBadRequest {
        t.Fatalf("status = %d; want %d", rec.Code, http.StatusBadRequest)
    }
}

func TestGetUserHandler_Success(t *testing.T) {
    repo := NewInMemoryUserRepository()
    service := NewUserService(repo)
    handler := NewUserHandler(service)
    mux := NewMux(handler)

    if err := service.CreateUser(User{ID: "1", Name: "Alice", Email: "alice@example.com"}); err != nil {
        t.Fatalf("seeding user: %v", err)
    }

    req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
    rec := httptest.NewRecorder()

    mux.ServeHTTP(rec, req)

    if rec.Code != http.StatusOK {
        t.Fatalf("status = %d; want %d, body = %s", rec.Code, http.StatusOK, rec.Body.String())
    }

    var got User
    if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
        t.Fatalf("decoding response body: %v", err)
    }
    if got.Name != "Alice" {
        t.Errorf("got Name = %q; want Alice", got.Name)
    }
}

func TestGetUserHandler_NotFound(t *testing.T) {
    mux := newTestMux()

    req := httptest.NewRequest(http.MethodGet, "/users/missing", nil)
    rec := httptest.NewRecorder()

    mux.ServeHTTP(rec, req)

    if rec.Code != http.StatusNotFound {
        t.Fatalf("status = %d; want %d", rec.Code, http.StatusNotFound)
    }
}
EOF
```

4. Run the tests:

```bash
go test -v ./...
```

**Expected Output:**

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

**Learning Objectives:**
- ✅ Write a Handler that only decodes/encodes and delegates to the Service
- ✅ Route with the stdlib `http.ServeMux`'s method + path-parameter syntax
- ✅ Test HTTP handlers with `httptest.NewRequest`/`NewRecorder`, no real server needed

---

## Exercise 5: Wiring All Three Layers Together, End-to-End

**Objective:** Construct the concrete repository, inject it into the service, inject the service into the handler, and prove the whole stack works over a real HTTP connection

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level30-exercise5
cd ~/projects/level30-exercise5
go mod init level30.example/exercise5
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "bytes"
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net"
    "net/http"
    "strings"
    "sync"
)

// --- Repository layer ---

type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

var ErrNotFound = errors.New("user not found")

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

var _ UserRepository = (*InMemoryUserRepository)(nil)

// --- Service layer ---

var (
    ErrInvalidName  = errors.New("name must not be empty")
    ErrInvalidEmail = errors.New("email must contain @")
)

type UserService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) CreateUser(u User) error {
    if strings.TrimSpace(u.Name) == "" {
        return fmt.Errorf("create user: %w", ErrInvalidName)
    }
    if !strings.Contains(u.Email, "@") {
        return fmt.Errorf("create user: %w", ErrInvalidEmail)
    }
    return s.repo.Create(u)
}

func (s *UserService) GetUser(id string) (User, error) {
    return s.repo.Get(id)
}

// --- Handler layer ---

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
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(u)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    u, err := h.service.GetUser(r.PathValue("id"))
    if err != nil {
        if errors.Is(err, ErrNotFound) {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }
        http.Error(w, "internal server error", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(u)
}

func NewMux(h *UserHandler) *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("POST /users", h.CreateUser)
    mux.HandleFunc("GET /users/{id}", h.GetUser)
    return mux
}

// --- Wiring (main) ---
//
// This is the ONLY place that knows about the concrete repository type.
// Every layer above only ever sees an interface. Constructing the real
// dependencies here and handing them down - Repository into Service,
// Service into Handler - is manual dependency injection: Level 31 gives
// this pattern a name and shows how to scale it beyond "wire it by hand".
func main() {
    repo := NewInMemoryUserRepository()          // concrete repository
    service := NewUserService(repo)              // service depends on the INTERFACE
    handler := NewUserHandler(service)            // handler depends on the service
    mux := NewMux(handler)

    listener, err := net.Listen("tcp", "127.0.0.1:0")
    if err != nil {
        panic(err)
    }
    addr := listener.Addr().String()

    server := &http.Server{Handler: mux}
    go server.Serve(listener)
    defer server.Shutdown(context.Background())

    fmt.Println("=== Wiring: Repository -> Service -> Handler -> http.Server ===")
    fmt.Println("server listening on", addr)
    baseURL := "http://" + addr

    fmt.Println("\n=== POST /users (real HTTP request over loopback) ===")
    body, _ := json.Marshal(User{ID: "1", Name: "Alice", Email: "alice@example.com"})
    resp, err := http.Post(baseURL+"/users", "application/json", bytes.NewReader(body))
    if err != nil {
        panic(err)
    }
    printResponse(resp)

    fmt.Println("\n=== GET /users/1 (real HTTP request over loopback) ===")
    resp, err = http.Get(baseURL + "/users/1")
    if err != nil {
        panic(err)
    }
    printResponse(resp)

    fmt.Println("\n=== GET /users/999 (not created - 404) ===")
    resp, err = http.Get(baseURL + "/users/999")
    if err != nil {
        panic(err)
    }
    printResponse(resp)
}

func printResponse(resp *http.Response) {
    defer resp.Body.Close()
    data, _ := io.ReadAll(resp.Body)
    fmt.Println("status:", resp.StatusCode)
    fmt.Println("body:", strings.TrimSpace(string(data)))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output** (the port number is chosen by your OS and will differ - everything else will match exactly):

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

**Learning Objectives:**
- ✅ Wire Repository, Service, and Handler together by hand in `main()`
- ✅ Recognize this manual wiring AS dependency injection (a direct Level 31 preview)
- ✅ Run a real `http.Server` on an OS-assigned loopback port, entirely offline
- ✅ Watch a request travel through all three layers and back over an actual TCP connection

---

## Exercise 6: Testing the Service in Isolation With a Fake Repository

**Objective:** Swap in a repository that always fails, and prove the Service surfaces the failure correctly

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level30-exercise6
cd ~/projects/level30-exercise6
go mod init level30.example/exercise6
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "errors"
    "fmt"
    "strings"
    "sync"
)

// --- Repository layer ---

type User struct {
    ID    string
    Name  string
    Email string
}

var ErrNotFound = errors.New("user not found")

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

var _ UserRepository = (*InMemoryUserRepository)(nil)

// --- Service layer ---

var (
    ErrInvalidName  = errors.New("name must not be empty")
    ErrInvalidEmail = errors.New("email must contain @")
)

type UserService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) CreateUser(u User) error {
    if strings.TrimSpace(u.Name) == "" {
        return fmt.Errorf("create user: %w", ErrInvalidName)
    }
    if !strings.Contains(u.Email, "@") {
        return fmt.Errorf("create user: %w", ErrInvalidEmail)
    }
    if err := s.repo.Create(u); err != nil {
        return fmt.Errorf("create user: %w", err)
    }
    return nil
}

func main() {
    fmt.Println("See main_test.go for this exercise's real work: testing")
    fmt.Println("UserService in isolation with a fake repository.")
}
EOF
```

3. Create `main_test.go`:

```bash
cat > main_test.go << 'EOF'
package main

import (
    "errors"
    "testing"
)

// ErrStorageUnavailable simulates a real-world failure: a database that's
// down, a disk that's full, a network partition - anything the SERVICE
// can't control but must still handle gracefully.
var ErrStorageUnavailable = errors.New("storage backend unavailable")

// alwaysErrorRepository is a fake UserRepository (Level 15/26's mocking
// pattern) that fails every call. It exists purely to prove the Service
// surfaces a repository failure correctly instead of hiding it, panicking,
// or returning a false success.
type alwaysErrorRepository struct{}

func (alwaysErrorRepository) Create(User) error {
    return ErrStorageUnavailable
}

func (alwaysErrorRepository) Get(string) (User, error) {
    return User{}, ErrStorageUnavailable
}

var _ UserRepository = alwaysErrorRepository{}

func TestUserService_CreateUser_Success(t *testing.T) {
    service := NewUserService(NewInMemoryUserRepository())

    err := service.CreateUser(User{ID: "1", Name: "Alice", Email: "alice@example.com"})
    if err != nil {
        t.Fatalf("CreateUser with a working repository returned error: %v", err)
    }
}

func TestUserService_CreateUser_RepositoryFailureSurfaces(t *testing.T) {
    service := NewUserService(alwaysErrorRepository{})

    err := service.CreateUser(User{ID: "1", Name: "Alice", Email: "alice@example.com"})

    if err == nil {
        t.Fatal("CreateUser returned nil error; want the repository failure to surface")
    }
    if !errors.Is(err, ErrStorageUnavailable) {
        t.Errorf("error = %v; want it to wrap ErrStorageUnavailable", err)
    }
    // The input itself was valid, so this must NOT look like a validation
    // failure - it must be attributable to the repository.
    if errors.Is(err, ErrInvalidName) || errors.Is(err, ErrInvalidEmail) {
        t.Errorf("error = %v; a valid input must not be reported as a validation error", err)
    }
}
EOF
```

4. Run the tests:

```bash
go test -v ./...
```

**Expected Output:**

```
=== RUN   TestUserService_CreateUser_Success
--- PASS: TestUserService_CreateUser_Success (0.00s)
=== RUN   TestUserService_CreateUser_RepositoryFailureSurfaces
--- PASS: TestUserService_CreateUser_RepositoryFailureSurfaces (0.00s)
PASS
ok  	level30.example/exercise6	0.781s
```

**Learning Objectives:**
- ✅ Write a fake repository that always fails, matching the real interface
- ✅ Test the Service completely in isolation - no real storage involved
- ✅ Prove a repository failure surfaces as an error, not a panic or a false success
- ✅ Distinguish a repository failure from a validation failure in the returned error

---

## Exercise 7: Mapping a Not-Found Repository Error to a 404 Response

**Objective:** Follow a "not found" error from the repository, through the service, into the correct HTTP status code

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level30-exercise7
cd ~/projects/level30-exercise7
go mod init level30.example/exercise7
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "strings"
    "sync"
)

// --- Repository layer ---

type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

var ErrNotFound = errors.New("user not found")

type UserRepository interface {
    Create(u User) error
    Get(id string) (User, error)
    List() []User
}

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

func (r *InMemoryUserRepository) List() []User {
    r.mu.Lock()
    defer r.mu.Unlock()
    out := make([]User, 0, len(r.users))
    for _, u := range r.users {
        out = append(out, u)
    }
    return out
}

var _ UserRepository = (*InMemoryUserRepository)(nil)

// --- Service layer ---

var (
    ErrInvalidName  = errors.New("name must not be empty")
    ErrInvalidEmail = errors.New("email must contain @")
    ErrEmailTaken   = errors.New("email already registered")
)

type UserService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

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
    if err := s.repo.Create(u); err != nil {
        return fmt.Errorf("create user: %w", err)
    }
    return nil
}

func (s *UserService) GetUser(id string) (User, error) {
    u, err := s.repo.Get(id)
    if err != nil {
        return User{}, fmt.Errorf("get user: %w", err)
    }
    return u, nil
}

// --- Handler layer ---
//
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
        http.Error(w, err.Error(), statusFor(err))
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(u)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    u, err := h.service.GetUser(r.PathValue("id"))
    if err != nil {
        http.Error(w, err.Error(), statusFor(err))
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(u)
}

func NewMux(h *UserHandler) *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("POST /users", h.CreateUser)
    mux.HandleFunc("GET /users/{id}", h.GetUser)
    return mux
}

func main() {
    fmt.Println("See main_test.go for this exercise's real work: an end-to-end")
    fmt.Println("error-translation table covering 400, 404, 409, and 500.")
}
EOF
```

3. Create `main_test.go`:

```bash
cat > main_test.go << 'EOF'
package main

import (
    "bytes"
    "encoding/json"
    "errors"
    "net/http"
    "net/http/httptest"
    "testing"
)

// flakyRepository wraps a real InMemoryUserRepository but deliberately
// fails Create for a magic ID, simulating an unexpected storage failure
// (a full disk, a dropped connection) that is NOT one of the Service's
// known sentinel errors.
type flakyRepository struct {
    *InMemoryUserRepository
}

func (f *flakyRepository) Create(u User) error {
    if u.ID == "boom" {
        return errors.New("disk full")
    }
    return f.InMemoryUserRepository.Create(u)
}

var _ UserRepository = (*flakyRepository)(nil)

func TestErrorTranslationAcrossLayers(t *testing.T) {
    tests := []struct {
        name       string
        setup      func(mux *http.ServeMux, service *UserService)
        method     string
        path       string
        body       any
        wantStatus int
    }{
        {
            name:       "valid create -> 201",
            method:     http.MethodPost,
            path:       "/users",
            body:       User{ID: "1", Name: "Alice", Email: "alice@example.com"},
            wantStatus: http.StatusCreated,
        },
        {
            name:       "invalid name -> 400",
            method:     http.MethodPost,
            path:       "/users",
            body:       User{ID: "2", Name: "   ", Email: "bob@example.com"},
            wantStatus: http.StatusBadRequest,
        },
        {
            name: "duplicate email -> 409",
            setup: func(mux *http.ServeMux, service *UserService) {
                _ = service.CreateUser(User{ID: "3", Name: "Carol", Email: "carol@example.com"})
            },
            method:     http.MethodPost,
            path:       "/users",
            body:       User{ID: "4", Name: "Dave", Email: "carol@example.com"},
            wantStatus: http.StatusConflict,
        },
        {
            name:       "unexpected storage failure -> 500",
            method:     http.MethodPost,
            path:       "/users",
            body:       User{ID: "boom", Name: "Eve", Email: "eve@example.com"},
            wantStatus: http.StatusInternalServerError,
        },
        {
            name: "get existing -> 200",
            setup: func(mux *http.ServeMux, service *UserService) {
                _ = service.CreateUser(User{ID: "5", Name: "Frank", Email: "frank@example.com"})
            },
            method:     http.MethodGet,
            path:       "/users/5",
            wantStatus: http.StatusOK,
        },
        {
            name:       "get missing -> 404",
            method:     http.MethodGet,
            path:       "/users/missing",
            wantStatus: http.StatusNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := &flakyRepository{InMemoryUserRepository: NewInMemoryUserRepository()}
            service := NewUserService(repo)
            handler := NewUserHandler(service)
            mux := NewMux(handler)

            if tt.setup != nil {
                tt.setup(mux, service)
            }

            var reqBody *bytes.Reader
            if tt.body != nil {
                data, _ := json.Marshal(tt.body)
                reqBody = bytes.NewReader(data)
            } else {
                reqBody = bytes.NewReader(nil)
            }

            req := httptest.NewRequest(tt.method, tt.path, reqBody)
            rec := httptest.NewRecorder()
            mux.ServeHTTP(rec, req)

            if rec.Code != tt.wantStatus {
                t.Errorf("status = %d; want %d, body = %s", rec.Code, tt.wantStatus, rec.Body.String())
            }
        })
    }
}
EOF
```

4. Run the tests:

```bash
go test -v ./...
```

**Expected Output:**

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

**Learning Objectives:**
- ✅ Build a single `statusFor(err error) int` translation point in the Handler
- ✅ Cover the full table: validation (400), conflict (409), not-found (404), unexpected (500)
- ✅ Prove the not-found error crosses Repository → Service → Handler correctly with a real `httptest` request
- ✅ Inject an unexpected repository failure and prove it falls through to 500, never a false success

---

## Exercise 8: Dependency Direction - Concrete Type vs Interface

**Objective:** Compare a Service coupled to a concrete repository type against one that depends on the interface, and prove only the second can be freely rewired

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level30-exercise8
cd ~/projects/level30-exercise8
go mod init level30.example/exercise8
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "errors"
    "fmt"
    "sync"
)

// --- Repository layer: TWO different implementations of the same contract ---

type User struct {
    ID    string
    Name  string
    Email string
}

var ErrNotFound = errors.New("user not found")

// UserRepository is defined here, at the point it's USED by the service
// layer below (Level 15's Rule 3) - not exported by the repository
// implementations themselves.
type UserRepository interface {
    Create(u User) error
    Get(id string) (User, error)
}

// InMemoryUserRepository: map-backed.
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

var _ UserRepository = (*InMemoryUserRepository)(nil)

// SliceUserRepository: slice-backed - a completely different internal
// shape, but it satisfies the exact same UserRepository interface.
type SliceUserRepository struct {
    mu    sync.Mutex
    users []User
}

func NewSliceUserRepository() *SliceUserRepository {
    return &SliceUserRepository{}
}

func (r *SliceUserRepository) Create(u User) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.users = append(r.users, u)
    return nil
}

func (r *SliceUserRepository) Get(id string) (User, error) {
    r.mu.Lock()
    defer r.mu.Unlock()
    for _, u := range r.users {
        if u.ID == id {
            return u, nil
        }
    }
    return User{}, fmt.Errorf("get user %s: %w", id, ErrNotFound)
}

var _ UserRepository = (*SliceUserRepository)(nil)

// --- BadUserService: depends on a CONCRETE repository type ---
//
// This compiles and works, but it is welded to InMemoryUserRepository.
// It could never accept a SliceUserRepository, a SQL-backed repository
// (Level 29), or a test fake without editing this struct's field type.
type BadUserService struct {
    repo *InMemoryUserRepository // concrete type, not an interface
}

func NewBadUserService(repo *InMemoryUserRepository) *BadUserService {
    return &BadUserService{repo: repo}
}

func (s *BadUserService) CreateUser(u User) error          { return s.repo.Create(u) }
func (s *BadUserService) GetUser(id string) (User, error) { return s.repo.Get(id) }

// --- GoodUserService: depends on the UserRepository INTERFACE ---
//
// Identical business logic, one type change. Because the dependency is
// an interface, ANY UserRepository implementation works: today's
// in-memory store, a future SQL store, or a test-only fake - with zero
// changes to this type.
type GoodUserService struct {
    repo UserRepository // interface
}

func NewGoodUserService(repo UserRepository) *GoodUserService {
    return &GoodUserService{repo: repo}
}

func (s *GoodUserService) CreateUser(u User) error          { return s.repo.Create(u) }
func (s *GoodUserService) GetUser(id string) (User, error) { return s.repo.Get(id) }

func main() {
    fmt.Println("=== BadUserService: Locked to InMemoryUserRepository ===")
    bad := NewBadUserService(NewInMemoryUserRepository())
    _ = bad.CreateUser(User{ID: "1", Name: "Alice", Email: "alice@example.com"})
    u, _ := bad.GetUser("1")
    fmt.Printf("BadUserService result: %+v\n", u)
    fmt.Println("NewBadUserService only accepts *InMemoryUserRepository - passing a")
    fmt.Println("*SliceUserRepository here would be a COMPILE ERROR, not a design choice.")

    fmt.Println("\n=== GoodUserService: Works With InMemoryUserRepository ===")
    goodMap := NewGoodUserService(NewInMemoryUserRepository())
    _ = goodMap.CreateUser(User{ID: "1", Name: "Alice", Email: "alice@example.com"})
    u, _ = goodMap.GetUser("1")
    fmt.Printf("GoodUserService(map-backed) result: %+v\n", u)

    fmt.Println("\n=== GoodUserService: SAME Code, Swapped to SliceUserRepository ===")
    goodSlice := NewGoodUserService(NewSliceUserRepository())
    _ = goodSlice.CreateUser(User{ID: "1", Name: "Alice", Email: "alice@example.com"})
    u, _ = goodSlice.GetUser("1")
    fmt.Printf("GoodUserService(slice-backed) result: %+v\n", u)

    fmt.Println("\nBoth GoodUserService results are identical - the service never")
    fmt.Println("noticed the repository's internals changed from a map to a slice.")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== BadUserService: Locked to InMemoryUserRepository ===
BadUserService result: {ID:1 Name:Alice Email:alice@example.com}
NewBadUserService only accepts *InMemoryUserRepository - passing a
*SliceUserRepository here would be a COMPILE ERROR, not a design choice.

=== GoodUserService: Works With InMemoryUserRepository ===
GoodUserService(map-backed) result: {ID:1 Name:Alice Email:alice@example.com}

=== GoodUserService: SAME Code, Swapped to SliceUserRepository ===
GoodUserService(slice-backed) result: {ID:1 Name:Alice Email:alice@example.com}

Both GoodUserService results are identical - the service never
noticed the repository's internals changed from a map to a slice.
```

**Learning Objectives:**
- ✅ See, in real running code, why a concrete-type dependency locks a Service to one implementation
- ✅ Prove an interface-based Service works unchanged across two different repository implementations
- ✅ Internalize the dependency direction: Service → Repository interface, never Service → Repository struct

---

## Exercise 9: A Full CRUD HTTP API, Table-Driven and Lifecycle-Tested

**Objective:** Complete the REST API with Update and Delete, and test a realistic Create → Read → Update → Read → Delete → Read lifecycle

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level30-exercise9
cd ~/projects/level30-exercise9
go mod init level30.example/exercise9
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "strings"
    "sync"
)

// --- Repository layer ---

type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

var ErrNotFound = errors.New("user not found")

type UserRepository interface {
    Create(u User) error
    Get(id string) (User, error)
    Update(u User) error
    Delete(id string) error
}

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

func (r *InMemoryUserRepository) Update(u User) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, ok := r.users[u.ID]; !ok {
        return fmt.Errorf("update user %s: %w", u.ID, ErrNotFound)
    }
    r.users[u.ID] = u
    return nil
}

func (r *InMemoryUserRepository) Delete(id string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, ok := r.users[id]; !ok {
        return fmt.Errorf("delete user %s: %w", id, ErrNotFound)
    }
    delete(r.users, id)
    return nil
}

var _ UserRepository = (*InMemoryUserRepository)(nil)

// --- Service layer ---

var (
    ErrInvalidName  = errors.New("name must not be empty")
    ErrInvalidEmail = errors.New("email must contain @")
)

type UserService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

func validate(u User) error {
    if strings.TrimSpace(u.Name) == "" {
        return ErrInvalidName
    }
    if !strings.Contains(u.Email, "@") {
        return ErrInvalidEmail
    }
    return nil
}

func (s *UserService) CreateUser(u User) error {
    if err := validate(u); err != nil {
        return fmt.Errorf("create user: %w", err)
    }
    return s.repo.Create(u)
}

func (s *UserService) GetUser(id string) (User, error) {
    u, err := s.repo.Get(id)
    if err != nil {
        return User{}, fmt.Errorf("get user: %w", err)
    }
    return u, nil
}

func (s *UserService) UpdateUser(u User) error {
    if err := validate(u); err != nil {
        return fmt.Errorf("update user: %w", err)
    }
    if err := s.repo.Update(u); err != nil {
        return fmt.Errorf("update user: %w", err)
    }
    return nil
}

func (s *UserService) DeleteUser(id string) error {
    if err := s.repo.Delete(id); err != nil {
        return fmt.Errorf("delete user: %w", err)
    }
    return nil
}

// --- Handler layer: full CRUD ---

func statusFor(err error) int {
    switch {
    case errors.Is(err, ErrInvalidName), errors.Is(err, ErrInvalidEmail):
        return http.StatusBadRequest
    case errors.Is(err, ErrNotFound):
        return http.StatusNotFound
    default:
        return http.StatusInternalServerError
    }
}

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
        http.Error(w, err.Error(), statusFor(err))
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(u)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    u, err := h.service.GetUser(r.PathValue("id"))
    if err != nil {
        http.Error(w, err.Error(), statusFor(err))
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(u)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
    var u User
    if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
        http.Error(w, "invalid JSON body", http.StatusBadRequest)
        return
    }
    u.ID = r.PathValue("id")
    if err := h.service.UpdateUser(u); err != nil {
        http.Error(w, err.Error(), statusFor(err))
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(u)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
    if err := h.service.DeleteUser(r.PathValue("id")); err != nil {
        http.Error(w, err.Error(), statusFor(err))
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

func NewMux(h *UserHandler) *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("POST /users", h.CreateUser)
    mux.HandleFunc("GET /users/{id}", h.GetUser)
    mux.HandleFunc("PUT /users/{id}", h.UpdateUser)
    mux.HandleFunc("DELETE /users/{id}", h.DeleteUser)
    return mux
}

func main() {
    fmt.Println("See main_test.go for this exercise's real work: a full")
    fmt.Println("Create/Read/Update/Delete HTTP API, table-driven and httptest-verified.")
}
EOF
```

3. Create `main_test.go`:

```bash
cat > main_test.go << 'EOF'
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

func doRequest(t *testing.T, mux *http.ServeMux, method, path string, body any) *httptest.ResponseRecorder {
    t.Helper()
    var reqBody *bytes.Reader
    if body != nil {
        data, err := json.Marshal(body)
        if err != nil {
            t.Fatalf("marshaling request body: %v", err)
        }
        reqBody = bytes.NewReader(data)
    } else {
        reqBody = bytes.NewReader(nil)
    }
    req := httptest.NewRequest(method, path, reqBody)
    rec := httptest.NewRecorder()
    mux.ServeHTTP(rec, req)
    return rec
}

// TestFullCRUDLifecycle drives one user through Create -> Read -> Update ->
// Read -> Delete -> Read, all through the real HTTP handlers, proving the
// three layers cooperate correctly across a realistic sequence.
func TestFullCRUDLifecycle(t *testing.T) {
    repo := NewInMemoryUserRepository()
    service := NewUserService(repo)
    mux := NewMux(NewUserHandler(service))

    t.Run("create", func(t *testing.T) {
        rec := doRequest(t, mux, http.MethodPost, "/users", User{ID: "1", Name: "Alice", Email: "alice@example.com"})
        if rec.Code != http.StatusCreated {
            t.Fatalf("status = %d; want %d, body = %s", rec.Code, http.StatusCreated, rec.Body.String())
        }
    })

    t.Run("read after create", func(t *testing.T) {
        rec := doRequest(t, mux, http.MethodGet, "/users/1", nil)
        if rec.Code != http.StatusOK {
            t.Fatalf("status = %d; want %d", rec.Code, http.StatusOK)
        }
        var got User
        if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
            t.Fatalf("decoding body: %v", err)
        }
        if got.Name != "Alice" {
            t.Errorf("Name = %q; want Alice", got.Name)
        }
    })

    t.Run("update", func(t *testing.T) {
        rec := doRequest(t, mux, http.MethodPut, "/users/1", User{Name: "Alice Updated", Email: "alice@example.com"})
        if rec.Code != http.StatusOK {
            t.Fatalf("status = %d; want %d, body = %s", rec.Code, http.StatusOK, rec.Body.String())
        }
    })

    t.Run("read after update", func(t *testing.T) {
        rec := doRequest(t, mux, http.MethodGet, "/users/1", nil)
        var got User
        if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
            t.Fatalf("decoding body: %v", err)
        }
        if got.Name != "Alice Updated" {
            t.Errorf("Name = %q; want Alice Updated", got.Name)
        }
    })

    t.Run("delete", func(t *testing.T) {
        rec := doRequest(t, mux, http.MethodDelete, "/users/1", nil)
        if rec.Code != http.StatusNoContent {
            t.Fatalf("status = %d; want %d", rec.Code, http.StatusNoContent)
        }
    })

    t.Run("read after delete -> 404", func(t *testing.T) {
        rec := doRequest(t, mux, http.MethodGet, "/users/1", nil)
        if rec.Code != http.StatusNotFound {
            t.Fatalf("status = %d; want %d", rec.Code, http.StatusNotFound)
        }
    })
}

// TestCRUDErrorCases covers the failure branches a happy-path lifecycle
// test never reaches: updating/deleting something that was never created,
// and updating with invalid data.
func TestCRUDErrorCases(t *testing.T) {
    tests := []struct {
        name       string
        method     string
        path       string
        body       any
        wantStatus int
    }{
        {"update missing user -> 404", http.MethodPut, "/users/missing", User{Name: "X", Email: "x@example.com"}, http.StatusNotFound},
        {"delete missing user -> 404", http.MethodDelete, "/users/missing", nil, http.StatusNotFound},
        {"update with invalid email -> 400", http.MethodPut, "/users/1", User{Name: "X", Email: "not-an-email"}, http.StatusBadRequest},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := NewInMemoryUserRepository()
            service := NewUserService(repo)
            mux := NewMux(NewUserHandler(service))
            // Seed user "1" so the "update with invalid email" case has
            // something to (fail to) update.
            _ = service.CreateUser(User{ID: "1", Name: "Alice", Email: "alice@example.com"})

            rec := doRequest(t, mux, tt.method, tt.path, tt.body)
            if rec.Code != tt.wantStatus {
                t.Errorf("status = %d; want %d, body = %s", rec.Code, tt.wantStatus, rec.Body.String())
            }
        })
    }
}
EOF
```

4. Run the tests:

```bash
go test -v ./...
```

**Expected Output:**

```
=== RUN   TestFullCRUDLifecycle
=== RUN   TestFullCRUDLifecycle/create
=== RUN   TestFullCRUDLifecycle/read_after_create
=== RUN   TestFullCRUDLifecycle/update
=== RUN   TestFullCRUDLifecycle/read_after_update
=== RUN   TestFullCRUDLifecycle/delete
=== RUN   TestFullCRUDLifecycle/read_after_delete_->_404
--- PASS: TestFullCRUDLifecycle (0.00s)
    --- PASS: TestFullCRUDLifecycle/create (0.00s)
    --- PASS: TestFullCRUDLifecycle/read_after_create (0.00s)
    --- PASS: TestFullCRUDLifecycle/update (0.00s)
    --- PASS: TestFullCRUDLifecycle/read_after_update (0.00s)
    --- PASS: TestFullCRUDLifecycle/delete (0.00s)
    --- PASS: TestFullCRUDLifecycle/read_after_delete_->_404 (0.00s)
=== RUN   TestCRUDErrorCases
=== RUN   TestCRUDErrorCases/update_missing_user_->_404
=== RUN   TestCRUDErrorCases/delete_missing_user_->_404
=== RUN   TestCRUDErrorCases/update_with_invalid_email_->_400
--- PASS: TestCRUDErrorCases (0.00s)
    --- PASS: TestCRUDErrorCases/update_missing_user_->_404 (0.00s)
    --- PASS: TestCRUDErrorCases/delete_missing_user_->_404 (0.00s)
    --- PASS: TestCRUDErrorCases/update_with_invalid_email_->_400 (0.00s)
PASS
ok  	level30.example/exercise9	0.625s
```

**Learning Objectives:**
- ✅ Complete a REST API's four core operations across all three layers
- ✅ Write a sequential lifecycle test that mirrors how a real client would use the API
- ✅ Cover the error branches a happy-path test never reaches (missing IDs, invalid updates)

---

## Exercise 10: Comprehensive Practice - A Complete Notes API

**Objective:** Build a small Notes API as REAL, separate packages - repository, service, and handler - wired together in `main()`, and tested at every layer independently plus one full end-to-end test

**Instructions:**

1. Create working directory and the package layout:

```bash
mkdir -p ~/projects/level30-exercise10/repository
mkdir -p ~/projects/level30-exercise10/service
mkdir -p ~/projects/level30-exercise10/handler
cd ~/projects/level30-exercise10
go mod init level30.example/exercise10
```

2. Create `repository/repository.go`:

```bash
cat > repository/repository.go << 'EOF'
// Package repository is the Repository layer: it knows how to store and
// retrieve Notes, and nothing else. No validation, no HTTP, no business
// rules - just data access.
package repository

import (
    "errors"
    "fmt"
    "sync"
)

// Note is the data record this level's comprehensive exercise works with.
type Note struct {
    ID    string `json:"id"`
    Title string `json:"title"`
    Body  string `json:"body"`
}

// ErrNotFound is returned when no Note matches the given ID.
var ErrNotFound = errors.New("note not found")

// InMemoryNoteRepository stores Notes in a map guarded by a mutex. It is
// ONE possible implementation of the NoteRepository interface the service
// package defines - a SQL-backed store (Level 29) could implement the
// exact same shape without this package changing at all.
type InMemoryNoteRepository struct {
    mu    sync.Mutex
    notes map[string]Note
}

// NewInMemoryNoteRepository builds an empty, ready-to-use repository.
func NewInMemoryNoteRepository() *InMemoryNoteRepository {
    return &InMemoryNoteRepository{notes: make(map[string]Note)}
}

func (r *InMemoryNoteRepository) Create(n Note) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.notes[n.ID] = n
    return nil
}

func (r *InMemoryNoteRepository) Get(id string) (Note, error) {
    r.mu.Lock()
    defer r.mu.Unlock()
    n, ok := r.notes[id]
    if !ok {
        return Note{}, fmt.Errorf("get note %s: %w", id, ErrNotFound)
    }
    return n, nil
}

func (r *InMemoryNoteRepository) Update(n Note) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, ok := r.notes[n.ID]; !ok {
        return fmt.Errorf("update note %s: %w", n.ID, ErrNotFound)
    }
    r.notes[n.ID] = n
    return nil
}

func (r *InMemoryNoteRepository) Delete(id string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, ok := r.notes[id]; !ok {
        return fmt.Errorf("delete note %s: %w", id, ErrNotFound)
    }
    delete(r.notes, id)
    return nil
}

func (r *InMemoryNoteRepository) List() []Note {
    r.mu.Lock()
    defer r.mu.Unlock()
    out := make([]Note, 0, len(r.notes))
    for _, n := range r.notes {
        out = append(out, n)
    }
    return out
}
EOF
```

3. Create `service/service.go`:

```bash
cat > service/service.go << 'EOF'
// Package service is the Service layer: business rules that sit between
// the Handler (transport) and the Repository (storage). It depends on the
// repository package only for the Note type and its sentinel error - the
// STORAGE CONTRACT it actually calls through is an interface defined
// right here, at the point of use (Level 15's Rule 3), not exported by
// the repository package itself.
package service

import (
    "errors"
    "fmt"
    "strings"

    "level30.example/exercise10/repository"
)

// Note is re-exported so callers of this package (handler, main, tests)
// never need to import the repository package just to spell the type.
type Note = repository.Note

// ErrNotFound is re-exported for the same reason - callers translate on
// this, never on repository.ErrNotFound directly.
var ErrNotFound = repository.ErrNotFound

// ErrInvalidTitle is a business rule the repository has no concept of: a
// map or a SQL table will happily store an empty string.
var ErrInvalidTitle = errors.New("title must not be empty")

// NoteRepository is the CONSUMER-SIDE interface: exactly the methods
// NoteService calls, nothing more. Any type satisfying it - the in-memory
// store in repository.go, a fake in a test, a future SQL store - works
// here without this file changing.
type NoteRepository interface {
    Create(n Note) error
    Get(id string) (Note, error)
    Update(n Note) error
    Delete(id string) error
    List() []Note
}

// NoteService holds the business rules and depends on NoteRepository the
// INTERFACE, never on *repository.InMemoryNoteRepository directly.
type NoteService struct {
    repo NoteRepository
}

// NewNoteService injects the repository dependency (Level 31 preview).
func NewNoteService(repo NoteRepository) *NoteService {
    return &NoteService{repo: repo}
}

// CreateNote validates BEFORE ever calling the repository.
func (s *NoteService) CreateNote(n Note) error {
    if strings.TrimSpace(n.Title) == "" {
        return fmt.Errorf("create note: %w", ErrInvalidTitle)
    }
    if err := s.repo.Create(n); err != nil {
        return fmt.Errorf("create note: %w", err)
    }
    return nil
}

func (s *NoteService) GetNote(id string) (Note, error) {
    n, err := s.repo.Get(id)
    if err != nil {
        return Note{}, fmt.Errorf("get note: %w", err)
    }
    return n, nil
}

func (s *NoteService) UpdateNote(n Note) error {
    if strings.TrimSpace(n.Title) == "" {
        return fmt.Errorf("update note: %w", ErrInvalidTitle)
    }
    if err := s.repo.Update(n); err != nil {
        return fmt.Errorf("update note: %w", err)
    }
    return nil
}

func (s *NoteService) DeleteNote(id string) error {
    if err := s.repo.Delete(id); err != nil {
        return fmt.Errorf("delete note: %w", err)
    }
    return nil
}

func (s *NoteService) ListNotes() []Note {
    return s.repo.List()
}
EOF
```

4. Create `service/service_test.go` - tests the Service completely in isolation:

```bash
cat > service/service_test.go << 'EOF'
package service

import (
    "errors"
    "testing"
)

// poisonRepository panics on every call - used to prove validation
// happens before the repository is ever touched.
type poisonRepository struct{}

func (poisonRepository) Create(Note) error        { panic("Create should not be called") }
func (poisonRepository) Get(string) (Note, error) { panic("Get should not be called") }
func (poisonRepository) Update(Note) error        { panic("Update should not be called") }
func (poisonRepository) Delete(string) error      { panic("Delete should not be called") }
func (poisonRepository) List() []Note             { panic("List should not be called") }

var _ NoteRepository = poisonRepository{}

// fakeRepository is a small in-memory test double distinct from the real
// repository package, used to test NoteService entirely in isolation.
type fakeRepository struct {
    notes map[string]Note
}

func newFakeRepository() *fakeRepository {
    return &fakeRepository{notes: make(map[string]Note)}
}

func (f *fakeRepository) Create(n Note) error {
    f.notes[n.ID] = n
    return nil
}

func (f *fakeRepository) Get(id string) (Note, error) {
    n, ok := f.notes[id]
    if !ok {
        return Note{}, ErrNotFound
    }
    return n, nil
}

func (f *fakeRepository) Update(n Note) error {
    f.notes[n.ID] = n
    return nil
}

func (f *fakeRepository) Delete(id string) error {
    delete(f.notes, id)
    return nil
}

func (f *fakeRepository) List() []Note {
    out := make([]Note, 0, len(f.notes))
    for _, n := range f.notes {
        out = append(out, n)
    }
    return out
}

var _ NoteRepository = (*fakeRepository)(nil)

func TestNoteService_CreateNote_ValidationRejectsBeforeRepository(t *testing.T) {
    svc := NewNoteService(poisonRepository{})

    err := svc.CreateNote(Note{ID: "1", Title: "   ", Body: "irrelevant"})

    if err == nil {
        t.Fatal("CreateNote with a blank title returned nil error")
    }
    if !errors.Is(err, ErrInvalidTitle) {
        t.Errorf("error = %v; want it to wrap ErrInvalidTitle", err)
    }
    // Reaching this line at all proves poisonRepository never panicked.
}

func TestNoteService_CreateNote_Success(t *testing.T) {
    repo := newFakeRepository()
    svc := NewNoteService(repo)

    err := svc.CreateNote(Note{ID: "1", Title: "Groceries", Body: "milk, eggs"})
    if err != nil {
        t.Fatalf("CreateNote returned unexpected error: %v", err)
    }

    got, err := svc.GetNote("1")
    if err != nil {
        t.Fatalf("GetNote returned unexpected error: %v", err)
    }
    if got.Title != "Groceries" {
        t.Errorf("Title = %q; want Groceries", got.Title)
    }
}

func TestNoteService_GetNote_NotFoundSurfaces(t *testing.T) {
    svc := NewNoteService(newFakeRepository())

    _, err := svc.GetNote("missing")

    if !errors.Is(err, ErrNotFound) {
        t.Errorf("error = %v; want it to wrap ErrNotFound", err)
    }
}
EOF
```

5. Create `handler/handler.go`:

```bash
cat > handler/handler.go << 'EOF'
// Package handler is the Handler layer: it translates HTTP requests into
// service calls and service results into HTTP responses. It contains no
// validation and no storage logic - both live below it.
package handler

import (
    "encoding/json"
    "errors"
    "net/http"

    "level30.example/exercise10/service"
)

// Note is re-exported purely for caller convenience.
type Note = service.Note

// NoteService is the CONSUMER-SIDE interface this package depends on -
// exactly the methods the handlers call. *service.NoteService satisfies
// it implicitly; a fake satisfies it just as easily in handler tests.
type NoteService interface {
    CreateNote(n Note) error
    GetNote(id string) (Note, error)
    UpdateNote(n Note) error
    DeleteNote(id string) error
    ListNotes() []Note
}

// statusFor translates a layered error into an HTTP status code - the
// ONE place in the program that knows about both errors and status codes.
func statusFor(err error) int {
    switch {
    case errors.Is(err, service.ErrInvalidTitle):
        return http.StatusBadRequest
    case errors.Is(err, service.ErrNotFound):
        return http.StatusNotFound
    default:
        return http.StatusInternalServerError
    }
}

type NoteHandler struct {
    service NoteService
}

func NewNoteHandler(s NoteService) *NoteHandler {
    return &NoteHandler{service: s}
}

func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
    var n Note
    if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
        http.Error(w, "invalid JSON body", http.StatusBadRequest)
        return
    }
    if err := h.service.CreateNote(n); err != nil {
        http.Error(w, err.Error(), statusFor(err))
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(n)
}

func (h *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
    n, err := h.service.GetNote(r.PathValue("id"))
    if err != nil {
        http.Error(w, err.Error(), statusFor(err))
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(n)
}

func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
    var n Note
    if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
        http.Error(w, "invalid JSON body", http.StatusBadRequest)
        return
    }
    n.ID = r.PathValue("id")
    if err := h.service.UpdateNote(n); err != nil {
        http.Error(w, err.Error(), statusFor(err))
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(n)
}

func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
    if err := h.service.DeleteNote(r.PathValue("id")); err != nil {
        http.Error(w, err.Error(), statusFor(err))
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

func (h *NoteHandler) ListNotes(w http.ResponseWriter, r *http.Request) {
    notes := h.service.ListNotes()
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(notes)
}

// NewMux wires routes to handler methods using the stdlib ServeMux
// (Level 27's approach) - no third-party router needed.
func NewMux(h *NoteHandler) *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("POST /notes", h.CreateNote)
    mux.HandleFunc("GET /notes", h.ListNotes)
    mux.HandleFunc("GET /notes/{id}", h.GetNote)
    mux.HandleFunc("PUT /notes/{id}", h.UpdateNote)
    mux.HandleFunc("DELETE /notes/{id}", h.DeleteNote)
    return mux
}
EOF
```

6. Create `handler/handler_test.go` - tests the Handler completely in isolation from the real Service:

```bash
cat > handler/handler_test.go << 'EOF'
package handler

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "level30.example/exercise10/service"
)

// fakeNoteService is a test double satisfying the handler's OWN
// NoteService interface - it never touches a real repository, so these
// tests check ONLY the handler's HTTP-translation logic in isolation.
type fakeNoteService struct {
    createErr error
    getErr    error
    note      Note
}

func (f *fakeNoteService) CreateNote(n Note) error         { return f.createErr }
func (f *fakeNoteService) GetNote(id string) (Note, error) { return f.note, f.getErr }
func (f *fakeNoteService) UpdateNote(n Note) error         { return nil }
func (f *fakeNoteService) DeleteNote(id string) error      { return nil }
func (f *fakeNoteService) ListNotes() []Note                { return []Note{f.note} }

var _ NoteService = (*fakeNoteService)(nil)

func TestNoteHandler_CreateNote_Success(t *testing.T) {
    fake := &fakeNoteService{}
    mux := NewMux(NewNoteHandler(fake))

    body, _ := json.Marshal(Note{ID: "1", Title: "Groceries"})
    req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewReader(body))
    rec := httptest.NewRecorder()
    mux.ServeHTTP(rec, req)

    if rec.Code != http.StatusCreated {
        t.Fatalf("status = %d; want %d, body = %s", rec.Code, http.StatusCreated, rec.Body.String())
    }
}

func TestNoteHandler_CreateNote_ValidationErrorMapsTo400(t *testing.T) {
    fake := &fakeNoteService{createErr: service.ErrInvalidTitle}
    mux := NewMux(NewNoteHandler(fake))

    body, _ := json.Marshal(Note{ID: "1", Title: ""})
    req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewReader(body))
    rec := httptest.NewRecorder()
    mux.ServeHTTP(rec, req)

    if rec.Code != http.StatusBadRequest {
        t.Fatalf("status = %d; want %d", rec.Code, http.StatusBadRequest)
    }
}

func TestNoteHandler_GetNote_NotFoundMapsTo404(t *testing.T) {
    fake := &fakeNoteService{getErr: service.ErrNotFound}
    mux := NewMux(NewNoteHandler(fake))

    req := httptest.NewRequest(http.MethodGet, "/notes/missing", nil)
    rec := httptest.NewRecorder()
    mux.ServeHTTP(rec, req)

    if rec.Code != http.StatusNotFound {
        t.Fatalf("status = %d; want %d", rec.Code, http.StatusNotFound)
    }
}

func TestNoteHandler_GetNote_Success(t *testing.T) {
    fake := &fakeNoteService{note: Note{ID: "1", Title: "Groceries", Body: "milk"}}
    mux := NewMux(NewNoteHandler(fake))

    req := httptest.NewRequest(http.MethodGet, "/notes/1", nil)
    rec := httptest.NewRecorder()
    mux.ServeHTTP(rec, req)

    if rec.Code != http.StatusOK {
        t.Fatalf("status = %d; want %d", rec.Code, http.StatusOK)
    }
    var got Note
    if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
        t.Fatalf("decoding body: %v", err)
    }
    if got.Title != "Groceries" {
        t.Errorf("Title = %q; want Groceries", got.Title)
    }
}
EOF
```

7. Create the top-level `main.go` that wires all three real packages together:

```bash
cat > main.go << 'EOF'
// Command exercise10 wires the three layers together: a concrete
// repository is injected into the service, and the service is injected
// into the handler. This is the whole level in one program, and a direct
// preview of Level 31 (Dependency Injection).
package main

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net"
    "net/http"
    "strings"

    "level30.example/exercise10/handler"
    "level30.example/exercise10/repository"
    "level30.example/exercise10/service"
)

func main() {
    repo := repository.NewInMemoryNoteRepository() // Repository: concrete
    svc := service.NewNoteService(repo)             // Service: depends on the interface
    h := handler.NewNoteHandler(svc)                 // Handler: depends on the service
    mux := handler.NewMux(h)

    listener, err := net.Listen("tcp", "127.0.0.1:0")
    if err != nil {
        panic(err)
    }
    server := &http.Server{Handler: mux}
    go server.Serve(listener)
    defer server.Shutdown(context.Background())

    baseURL := "http://" + listener.Addr().String()
    fmt.Println("=== Notes API: Repository -> Service -> Handler, wired in main() ===")

    fmt.Println("\n=== POST /notes ===")
    body, _ := json.Marshal(map[string]string{"id": "1", "title": "Groceries", "body": "milk, eggs"})
    resp, _ := http.Post(baseURL+"/notes", "application/json", bytes.NewReader(body))
    printResponse(resp)

    fmt.Println("\n=== GET /notes/1 ===")
    resp, _ = http.Get(baseURL + "/notes/1")
    printResponse(resp)

    fmt.Println("\n=== POST /notes with an empty title (rejected by the Service) ===")
    body, _ = json.Marshal(map[string]string{"id": "2", "title": ""})
    resp, _ = http.Post(baseURL+"/notes", "application/json", bytes.NewReader(body))
    printResponse(resp)

    fmt.Println("\n=== GET /notes/999 (never created) ===")
    resp, _ = http.Get(baseURL + "/notes/999")
    printResponse(resp)
}

func printResponse(resp *http.Response) {
    defer resp.Body.Close()
    data, _ := io.ReadAll(resp.Body)
    fmt.Println("status:", resp.StatusCode)
    fmt.Println("body:", strings.TrimSpace(string(data)))
}
EOF
```

8. Create the top-level `main_test.go` - the ONE full end-to-end test wiring all three REAL packages together:

```bash
cat > main_test.go << 'EOF'
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "level30.example/exercise10/handler"
    "level30.example/exercise10/repository"
    "level30.example/exercise10/service"
)

// TestNotesAPI_EndToEnd wires the REAL repository, service, and handler
// together - exactly as main() does - and drives a full lifecycle through
// the HTTP layer with httptest. This is the one test in this exercise
// that exercises all three layers cooperating at once; service_test.go
// and handler/handler_test.go already proved each layer in isolation.
func TestNotesAPI_EndToEnd(t *testing.T) {
    repo := repository.NewInMemoryNoteRepository()
    svc := service.NewNoteService(repo)
    mux := handler.NewMux(handler.NewNoteHandler(svc))

    // Create
    body, _ := json.Marshal(map[string]string{"id": "1", "title": "Groceries", "body": "milk"})
    req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewReader(body))
    rec := httptest.NewRecorder()
    mux.ServeHTTP(rec, req)
    if rec.Code != http.StatusCreated {
        t.Fatalf("create: status = %d; want %d, body = %s", rec.Code, http.StatusCreated, rec.Body.String())
    }

    // Read
    req = httptest.NewRequest(http.MethodGet, "/notes/1", nil)
    rec = httptest.NewRecorder()
    mux.ServeHTTP(rec, req)
    if rec.Code != http.StatusOK {
        t.Fatalf("get: status = %d; want %d", rec.Code, http.StatusOK)
    }

    // Update
    body, _ = json.Marshal(map[string]string{"title": "Groceries v2", "body": "milk, eggs"})
    req = httptest.NewRequest(http.MethodPut, "/notes/1", bytes.NewReader(body))
    rec = httptest.NewRecorder()
    mux.ServeHTTP(rec, req)
    if rec.Code != http.StatusOK {
        t.Fatalf("update: status = %d; want %d, body = %s", rec.Code, http.StatusOK, rec.Body.String())
    }

    // Invalid create is rejected by the Service and mapped to 400 by the Handler
    body, _ = json.Marshal(map[string]string{"id": "2", "title": ""})
    req = httptest.NewRequest(http.MethodPost, "/notes", bytes.NewReader(body))
    rec = httptest.NewRecorder()
    mux.ServeHTTP(rec, req)
    if rec.Code != http.StatusBadRequest {
        t.Fatalf("invalid create: status = %d; want %d", rec.Code, http.StatusBadRequest)
    }

    // Delete
    req = httptest.NewRequest(http.MethodDelete, "/notes/1", nil)
    rec = httptest.NewRecorder()
    mux.ServeHTTP(rec, req)
    if rec.Code != http.StatusNoContent {
        t.Fatalf("delete: status = %d; want %d", rec.Code, http.StatusNoContent)
    }

    // Read after delete -> 404, proving the not-found error crossed all
    // three layers correctly.
    req = httptest.NewRequest(http.MethodGet, "/notes/1", nil)
    rec = httptest.NewRecorder()
    mux.ServeHTTP(rec, req)
    if rec.Code != http.StatusNotFound {
        t.Fatalf("get after delete: status = %d; want %d", rec.Code, http.StatusNotFound)
    }
}
EOF
```

9. Run every layer's tests:

```bash
go test -v ./...
```

**Expected Output:**

```
=== RUN   TestNotesAPI_EndToEnd
--- PASS: TestNotesAPI_EndToEnd (0.00s)
PASS
ok  	level30.example/exercise10	0.273s
=== RUN   TestNoteHandler_CreateNote_Success
--- PASS: TestNoteHandler_CreateNote_Success (0.00s)
=== RUN   TestNoteHandler_CreateNote_ValidationErrorMapsTo400
--- PASS: TestNoteHandler_CreateNote_ValidationErrorMapsTo400 (0.00s)
=== RUN   TestNoteHandler_GetNote_NotFoundMapsTo404
--- PASS: TestNoteHandler_GetNote_NotFoundMapsTo404 (0.00s)
=== RUN   TestNoteHandler_GetNote_Success
--- PASS: TestNoteHandler_GetNote_Success (0.00s)
PASS
ok  	level30.example/exercise10/handler	0.762s
?   	level30.example/exercise10/repository	[no test files]
=== RUN   TestNoteService_CreateNote_ValidationRejectsBeforeRepository
--- PASS: TestNoteService_CreateNote_ValidationRejectsBeforeRepository (0.00s)
=== RUN   TestNoteService_CreateNote_Success
--- PASS: TestNoteService_CreateNote_Success (0.00s)
=== RUN   TestNoteService_GetNote_NotFoundSurfaces
--- PASS: TestNoteService_GetNote_NotFoundSurfaces (0.00s)
PASS
ok  	level30.example/exercise10/service	0.536s
```

10. Run the wired program itself:

```bash
go run .
```

**Expected Output:**

```
=== Notes API: Repository -> Service -> Handler, wired in main() ===

=== POST /notes ===
status: 201
body: {"id":"1","title":"Groceries","body":"milk, eggs"}

=== GET /notes/1 ===
status: 200
body: {"id":"1","title":"Groceries","body":"milk, eggs"}

=== POST /notes with an empty title (rejected by the Service) ===
status: 400
body: create note: title must not be empty

=== GET /notes/999 (never created) ===
status: 404
body: get note: get note 999: note not found
```

**Learning Objectives:**
- ✅ Split Repository, Service, and Handler into REAL separate Go packages, not just sections of one file
- ✅ Define the `NoteRepository` interface in the service package, and the `NoteService` interface in the handler package - each consumer owning its own contract (Level 15's Rule 3, applied twice)
- ✅ Test the Service in isolation with a fake repository, and the Handler in isolation with a fake service
- ✅ Write one true end-to-end test wiring all three real packages together
- ✅ Wire everything by hand in `main()` and run it as a real (offline, loopback-only) HTTP program

---

## Bonus Challenges

### Challenge 1: A Second Repository Implementation

Add a second `UserRepository` implementation with a different internal shape (for example, a slice-backed store instead of a map-backed one, or a store keyed by email instead of ID), and prove the same `UserService` works unchanged against either one.

```bash
mkdir -p ~/projects/level30-bonus1
cd ~/projects/level30-bonus1
go mod init level30.example/bonus1
```

**Hints:**
- Reuse Exercise 8's pattern: two repository implementations, one interface, one service
- Run the exact same sequence of service calls against both repositories and diff the results
- The service code should not need a single `if` statement added to support the second implementation

### Challenge 2: A Caching Decorator

Write a `CachingUserRepository` that wraps any `UserRepository` and satisfies the same interface, caching `Get` results in memory so repeated lookups for the same ID skip the wrapped repository.

```bash
mkdir -p ~/projects/level30-bonus2
cd ~/projects/level30-bonus2
go mod init level30.example/bonus2
```

**Hints:**
- `type CachingUserRepository struct { inner UserRepository; cache map[string]User }`
- `Get` checks the cache first; on a miss, calls `inner.Get`, stores the result, then returns it
- `Create`, `Update`, and `Delete` should invalidate (or update) the cache entry for that ID before delegating to `inner`
- Prove it works by wrapping a repository that counts its own calls, and showing the count doesn't increase on a cached `Get`

### Challenge 3: Service-Level Pagination and Filtering

Add a `ListUsers(limit, offset int) ([]User, error)` method to the Service that returns a page of results, and extend the `UserRepository` interface with whatever the repository needs to support it efficiently.

```bash
mkdir -p ~/projects/level30-bonus3
cd ~/projects/level30-bonus3
go mod init level30.example/bonus3
```

**Hints:**
- The simplest correct approach: `List() []User` in the repository, slicing and sorting (for a stable order) in the Service
- A more realistic approach: change the repository interface itself to `List(limit, offset int) ([]User, error)`, pushing pagination down - decide which layer should own "what does page 2 mean?"
- Either way, write a test with more records than one page holds, and assert the returned slice has exactly `limit` items starting at the right offset

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Define a small Repository interface and satisfy it with an in-memory implementation
✅ Write Service logic that enforces validation and uniqueness rules the repository can't
✅ Prove a Service rejects invalid input before ever touching the repository
✅ Write thin Handlers that translate HTTP to Service calls, tested with `httptest`
✅ Wire Repository, Service, and Handler together by hand in `main()` - manual dependency injection
✅ Test the Service in isolation with fake repositories, including one that always fails
✅ Translate errors across all three layers into the correct HTTP status codes (400/404/409/500)
✅ Explain why a Service depending on a concrete type breaks the pattern, with real running proof
✅ Build a complete CRUD HTTP API across all three layers
✅ Structure a real multi-package Go program along Repository/Service/Handler lines, with each layer's interface owned by its consumer

---

## Next Level

Level 31: Dependency Injection
- Naming and formalizing the "construct dependencies and hand them down" pattern from Exercises 5 and 10
- Managing a dependency graph as it grows past three layers
- Constructor injection, wiring functions, and where DI frameworks fit into Go

Great work! You're structuring programs the way real production Go services are built! 🚀
