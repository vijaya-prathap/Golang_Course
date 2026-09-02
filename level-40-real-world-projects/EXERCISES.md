# Level 40: Real-World Projects - Exercises

This level is **one continuous guided build**, not ten independent exercises. Every milestone adds to the SAME project - a URL Shortener API - in the order a real project actually gets built: storage, then business rules, then HTTP, then security, then observability, then resilience, then proof, then packaging.

**Project layout across all 10 milestones** (create this once, in Milestone 1, and keep adding to it):

```
~/projects/level40-capstone/urlshortener/
├── go.mod
├── main.go                    # added in Milestone 9
├── domain/
│   ├── link.go                # Milestone 1
│   └── user.go                # Milestone 2
├── repository/
│   ├── link_repository.go     # Milestone 1
│   ├── link_repository_test.go
│   ├── user_repository.go     # Milestone 2
│   └── user_repository_test.go
├── service/
│   ├── shortcode.go           # Milestone 3
│   ├── shortcode_test.go
│   ├── link_service.go        # Milestone 3
│   ├── link_service_test.go
│   ├── user_service.go        # Milestone 4
│   └── user_service_test.go
├── handler/
│   ├── context.go             # Milestone 5
│   ├── link_handler.go        # Milestone 5
│   ├── auth_handler.go        # Milestone 5
│   ├── routes.go              # Milestone 5
│   └── handler_test.go
├── middleware/
│   ├── auth.go                # Milestone 6
│   ├── auth_test.go
│   ├── logging.go             # Milestone 7
│   └── logging_test.go
├── config/
│   ├── config.go              # Milestone 7
│   └── config_test.go
├── ratelimit/
│   ├── tokenbucket.go         # Milestone 8
│   └── tokenbucket_test.go
├── main_test.go                # Milestone 9
├── Dockerfile                  # Milestone 10
└── k8s/
    ├── deployment.yaml         # Milestone 10
    └── service.yaml            # Milestone 10
```

Complete milestones in order - each one assumes every file from every earlier milestone already exists.

---

## Milestone 1: Project Setup, Domain Models, and the Link Repository

**Objective:** Set up the module and build the Repository layer's Link half - an interface plus an in-memory, mutex-guarded implementation (Level 30, Level 23) - with its own passing tests before anything else exists.

**Instructions:**

1. Create the project and module:

```bash
mkdir -p ~/projects/level40-capstone/urlshortener
cd ~/projects/level40-capstone/urlshortener
go mod init level40.example/urlshortener
mkdir -p domain repository service handler middleware ratelimit config
```

2. Create `domain/link.go`:

```bash
cat > domain/link.go << 'EOF'
// Package domain holds the plain data records shared across every layer -
// no HTTP, no storage, no business rules. Just the shapes of the data.
package domain

import "time"

// Link is a shortened URL: a generated Code that redirects to a
// LongURL, owned by the user who created it, with a running click count.
type Link struct {
    Code      string    `json:"code"`
    LongURL   string    `json:"long_url"`
    OwnerID   string    `json:"owner_id"`
    Clicks    int       `json:"clicks"`
    CreatedAt time.Time `json:"created_at"`
}
EOF
```

3. Create `repository/link_repository.go`:

```bash
cat > repository/link_repository.go << 'EOF'
// Package repository is the Repository layer (Level 30): it knows how to
// store and retrieve Links and Users, and nothing else. No validation, no
// short-code generation, no HTTP - just data access, backed by an
// in-memory map guarded by a mutex (Level 23) so it stays safe once
// concurrent HTTP requests start hitting it.
package repository

import (
    "errors"
    "fmt"
    "sync"

    "level40.example/urlshortener/domain"
)

// ErrLinkNotFound is returned when no Link matches the given code.
var ErrLinkNotFound = errors.New("link not found")

// ErrCodeTaken is returned when a code is already in use - the repository
// enforces storage-level uniqueness; the SERVICE layer (Level 30) decides
// what to do about it (generate a different code, retry, etc).
var ErrCodeTaken = errors.New("short code already exists")

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

// InMemoryLinkRepository stores Links in a map guarded by a mutex.
type InMemoryLinkRepository struct {
    mu    sync.Mutex
    links map[string]domain.Link
}

// NewInMemoryLinkRepository builds an empty, ready-to-use repository.
func NewInMemoryLinkRepository() *InMemoryLinkRepository {
    return &InMemoryLinkRepository{links: make(map[string]domain.Link)}
}

func (r *InMemoryLinkRepository) Create(l domain.Link) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, exists := r.links[l.Code]; exists {
        return fmt.Errorf("create link %s: %w", l.Code, ErrCodeTaken)
    }
    r.links[l.Code] = l
    return nil
}

func (r *InMemoryLinkRepository) Get(code string) (domain.Link, error) {
    r.mu.Lock()
    defer r.mu.Unlock()
    l, ok := r.links[code]
    if !ok {
        return domain.Link{}, fmt.Errorf("get link %s: %w", code, ErrLinkNotFound)
    }
    return l, nil
}

// IncrementClicks bumps a Link's click counter by one. It is its own
// method (rather than folded into Get) because incrementing on every
// redirect is a WRITE, and callers that just want to read stats (GET
// /stats/:code) must never trigger a side effect.
func (r *InMemoryLinkRepository) IncrementClicks(code string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    l, ok := r.links[code]
    if !ok {
        return fmt.Errorf("increment clicks %s: %w", code, ErrLinkNotFound)
    }
    l.Clicks++
    r.links[code] = l
    return nil
}

func (r *InMemoryLinkRepository) CountByOwner(ownerID string) (int, error) {
    r.mu.Lock()
    defer r.mu.Unlock()
    count := 0
    for _, l := range r.links {
        if l.OwnerID == ownerID {
            count++
        }
    }
    return count, nil
}

func (r *InMemoryLinkRepository) List() []domain.Link {
    r.mu.Lock()
    defer r.mu.Unlock()
    out := make([]domain.Link, 0, len(r.links))
    for _, l := range r.links {
        out = append(out, l)
    }
    return out
}

// Compile-time check: a drifted method signature becomes a build error,
// not a runtime surprise (Level 30's Best Practice #5).
var _ LinkRepository = (*InMemoryLinkRepository)(nil)
EOF
```

4. Create `repository/link_repository_test.go`:

```bash
cat > repository/link_repository_test.go << 'EOF'
package repository

import (
    "errors"
    "testing"
    "time"

    "level40.example/urlshortener/domain"
)

func TestInMemoryLinkRepository_CreateAndGet(t *testing.T) {
    repo := NewInMemoryLinkRepository()
    link := domain.Link{Code: "abc123", LongURL: "https://example.com", OwnerID: "u1", CreatedAt: time.Now()}

    if err := repo.Create(link); err != nil {
        t.Fatalf("Create returned unexpected error: %v", err)
    }

    got, err := repo.Get("abc123")
    if err != nil {
        t.Fatalf("Get returned unexpected error: %v", err)
    }
    if got.LongURL != link.LongURL {
        t.Errorf("LongURL = %q; want %q", got.LongURL, link.LongURL)
    }
}

func TestInMemoryLinkRepository_CreateDuplicateCodeRejected(t *testing.T) {
    repo := NewInMemoryLinkRepository()
    link := domain.Link{Code: "abc123", LongURL: "https://example.com", OwnerID: "u1"}
    _ = repo.Create(link)

    err := repo.Create(link)

    if !errors.Is(err, ErrCodeTaken) {
        t.Errorf("error = %v; want it to wrap ErrCodeTaken", err)
    }
}

func TestInMemoryLinkRepository_GetMissing(t *testing.T) {
    repo := NewInMemoryLinkRepository()

    _, err := repo.Get("missing")

    if !errors.Is(err, ErrLinkNotFound) {
        t.Errorf("error = %v; want it to wrap ErrLinkNotFound", err)
    }
}

func TestInMemoryLinkRepository_IncrementClicks(t *testing.T) {
    repo := NewInMemoryLinkRepository()
    _ = repo.Create(domain.Link{Code: "abc123", LongURL: "https://example.com", OwnerID: "u1"})

    _ = repo.IncrementClicks("abc123")
    _ = repo.IncrementClicks("abc123")
    got, _ := repo.Get("abc123")

    if got.Clicks != 2 {
        t.Errorf("Clicks = %d; want 2", got.Clicks)
    }
}

func TestInMemoryLinkRepository_CountByOwner(t *testing.T) {
    repo := NewInMemoryLinkRepository()
    _ = repo.Create(domain.Link{Code: "a1", OwnerID: "u1"})
    _ = repo.Create(domain.Link{Code: "a2", OwnerID: "u1"})
    _ = repo.Create(domain.Link{Code: "a3", OwnerID: "u2"})

    count, err := repo.CountByOwner("u1")
    if err != nil {
        t.Fatalf("CountByOwner returned unexpected error: %v", err)
    }
    if count != 2 {
        t.Errorf("count = %d; want 2", count)
    }
}
EOF
```

5. Run the tests:

```bash
go vet ./...
go test ./repository/... -v
```

**Expected Output:**

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
PASS
ok  	level40.example/urlshortener/repository	0.542s
```

**Learning Objectives:**
- ✅ Set up a multi-package module the way a real project starts (Level 17)
- ✅ Build a Repository interface plus a mutex-guarded in-memory implementation (Level 23, Level 30)
- ✅ Separate a read (`Get`) from a write (`IncrementClicks`) at the interface level, on purpose

---

## Milestone 2: The User Repository

**Objective:** Add the Repository layer's other half - Users, with username uniqueness enforced at the storage layer.

**Instructions:**

1. Create `domain/user.go`:

```bash
cat > domain/user.go << 'EOF'
package domain

import "time"

// User is a registered account. PasswordHash is a bcrypt hash - the
// plaintext password is never stored (Level 32).
type User struct {
    ID           string    `json:"id"`
    Username     string    `json:"username"`
    PasswordHash string    `json:"-"`
    CreatedAt    time.Time `json:"created_at"`
}
EOF
```

2. Create `repository/user_repository.go`:

```bash
cat > repository/user_repository.go << 'EOF'
package repository

import (
    "errors"
    "fmt"
    "sync"

    "level40.example/urlshortener/domain"
)

// ErrUserNotFound is returned when no User matches the given lookup.
var ErrUserNotFound = errors.New("user not found")

// ErrUsernameTaken is returned when a username is already registered.
var ErrUsernameTaken = errors.New("username already registered")

// UserRepository is the Repository layer's contract for Users.
type UserRepository interface {
    Create(u domain.User) error
    GetByUsername(username string) (domain.User, error)
    GetByID(id string) (domain.User, error)
}

// InMemoryUserRepository stores Users in a map guarded by a mutex.
type InMemoryUserRepository struct {
    mu    sync.Mutex
    users map[string]domain.User // keyed by ID
}

// NewInMemoryUserRepository builds an empty, ready-to-use repository.
func NewInMemoryUserRepository() *InMemoryUserRepository {
    return &InMemoryUserRepository{users: make(map[string]domain.User)}
}

func (r *InMemoryUserRepository) Create(u domain.User) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    for _, existing := range r.users {
        if existing.Username == u.Username {
            return fmt.Errorf("create user %s: %w", u.Username, ErrUsernameTaken)
        }
    }
    r.users[u.ID] = u
    return nil
}

func (r *InMemoryUserRepository) GetByUsername(username string) (domain.User, error) {
    r.mu.Lock()
    defer r.mu.Unlock()
    for _, u := range r.users {
        if u.Username == username {
            return u, nil
        }
    }
    return domain.User{}, fmt.Errorf("get user %s: %w", username, ErrUserNotFound)
}

func (r *InMemoryUserRepository) GetByID(id string) (domain.User, error) {
    r.mu.Lock()
    defer r.mu.Unlock()
    u, ok := r.users[id]
    if !ok {
        return domain.User{}, fmt.Errorf("get user %s: %w", id, ErrUserNotFound)
    }
    return u, nil
}

var _ UserRepository = (*InMemoryUserRepository)(nil)
EOF
```

3. Create `repository/user_repository_test.go`:

```bash
cat > repository/user_repository_test.go << 'EOF'
package repository

import (
    "errors"
    "testing"

    "level40.example/urlshortener/domain"
)

func TestInMemoryUserRepository_CreateAndGet(t *testing.T) {
    repo := NewInMemoryUserRepository()
    user := domain.User{ID: "u1", Username: "alice", PasswordHash: "hash"}

    if err := repo.Create(user); err != nil {
        t.Fatalf("Create returned unexpected error: %v", err)
    }

    byID, err := repo.GetByID("u1")
    if err != nil || byID.Username != "alice" {
        t.Errorf("GetByID = %+v, err=%v; want alice, nil", byID, err)
    }

    byName, err := repo.GetByUsername("alice")
    if err != nil || byName.ID != "u1" {
        t.Errorf("GetByUsername = %+v, err=%v; want u1, nil", byName, err)
    }
}

func TestInMemoryUserRepository_DuplicateUsernameRejected(t *testing.T) {
    repo := NewInMemoryUserRepository()
    _ = repo.Create(domain.User{ID: "u1", Username: "alice"})

    err := repo.Create(domain.User{ID: "u2", Username: "alice"})

    if !errors.Is(err, ErrUsernameTaken) {
        t.Errorf("error = %v; want it to wrap ErrUsernameTaken", err)
    }
}

func TestInMemoryUserRepository_GetMissing(t *testing.T) {
    repo := NewInMemoryUserRepository()

    _, err := repo.GetByUsername("ghost")

    if !errors.Is(err, ErrUserNotFound) {
        t.Errorf("error = %v; want it to wrap ErrUserNotFound", err)
    }
}
EOF
```

4. Run the full repository suite:

```bash
go vet ./...
go test ./repository/... -v
```

**Expected Output:**

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
ok  	level40.example/urlshortener/repository	0.548s
```

**Learning Objectives:**
- ✅ Enforce a uniqueness rule (username) at the storage layer while leaving password rules to the Service
- ✅ Two lookup methods (`GetByUsername`, `GetByID`) on the same repository, each serving a different caller (login vs. token validation later)

---

## Milestone 3: The Service Layer - Short Codes and Link Business Rules

**Objective:** Build the Service layer's Link half: URL validation, a per-user link limit, and `crypto/rand`-based short code generation - all tested against a **fake** repository, never the real one (Level 30's testing payoff).

**Instructions:**

1. Create `service/shortcode.go`:

```bash
cat > service/shortcode.go << 'EOF'
// Package service is the Service layer (Level 30): business rules that sit
// between the Handler (transport) and the Repository (storage). It depends
// on repository interfaces defined AT THE POINT OF USE, right here
// (Level 15's Rule 3, Level 30's Best Practice #1) - never on the concrete
// in-memory implementations directly.
package service

import (
    "crypto/rand"
    "math/big"
)

const codeAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateCode returns a random, URL-safe short code of the given length,
// drawn from crypto/rand rather than math/rand - a short code that an
// attacker could predict would let them enumerate other users' links.
func GenerateCode(length int) (string, error) {
    code := make([]byte, length)
    for i := range code {
        n, err := rand.Int(rand.Reader, big.NewInt(int64(len(codeAlphabet))))
        if err != nil {
            return "", err
        }
        code[i] = codeAlphabet[n.Int64()]
    }
    return string(code), nil
}
EOF
```

2. Create `service/shortcode_test.go`:

```bash
cat > service/shortcode_test.go << 'EOF'
package service

import "testing"

func TestGenerateCode_Length(t *testing.T) {
    code, err := GenerateCode(7)
    if err != nil {
        t.Fatalf("GenerateCode returned unexpected error: %v", err)
    }
    if len(code) != 7 {
        t.Errorf("len(code) = %d; want 7", len(code))
    }
}

func TestGenerateCode_Uniqueness(t *testing.T) {
    seen := make(map[string]bool)
    for i := 0; i < 1000; i++ {
        code, err := GenerateCode(7)
        if err != nil {
            t.Fatalf("GenerateCode returned unexpected error: %v", err)
        }
        if seen[code] {
            t.Fatalf("GenerateCode produced a duplicate: %s", code)
        }
        seen[code] = true
    }
}
EOF
```

3. Create `service/link_service.go`:

```bash
cat > service/link_service.go << 'EOF'
package service

import (
    "errors"
    "fmt"
    "net/url"
    "strings"
    "time"

    "level40.example/urlshortener/domain"
    "level40.example/urlshortener/repository"
)

// Re-exported so callers (handler, main, tests) never need to import the
// repository package just to spell the type or check a sentinel error -
// the same re-export convention Level 30 uses throughout.
type Link = domain.Link

var (
    ErrInvalidURL     = errors.New("long_url must be an absolute http(s) URL")
    ErrLinkLimit      = errors.New("link limit reached for this user")
    ErrLinkNotFound   = repository.ErrLinkNotFound
    ErrCodeGeneration = errors.New("could not generate a unique short code")
)

// MaxLinksPerUser is the per-user limit the service enforces - a business
// rule the repository has no concept of.
const MaxLinksPerUser = 5

// maxCodeAttempts bounds how many times CreateLink retries GenerateCode
// after a collision before giving up.
const maxCodeAttempts = 5

// LinkRepository is the CONSUMER-SIDE interface: exactly what LinkService
// calls, nothing more.
type LinkRepository interface {
    Create(l domain.Link) error
    Get(code string) (domain.Link, error)
    IncrementClicks(code string) error
    CountByOwner(ownerID string) (int, error)
    List() []domain.Link
}

// LinkService holds the business rules for shortening and resolving URLs.
// It depends on the LinkRepository INTERFACE, never a concrete type.
type LinkService struct {
    repo    LinkRepository
    codeLen int
    newCode func(int) (string, error)
}

// NewLinkService injects the repository dependency (Level 31's constructor
// injection). codeLen controls how many characters a generated short code
// has.
func NewLinkService(repo LinkRepository, codeLen int) *LinkService {
    return &LinkService{repo: repo, codeLen: codeLen, newCode: GenerateCode}
}

// validateURL rejects anything that is not an absolute http/https URL -
// the repository would happily store "not a url" forever.
func validateURL(raw string) error {
    raw = strings.TrimSpace(raw)
    if raw == "" {
        return ErrInvalidURL
    }
    u, err := url.ParseRequestURI(raw)
    if err != nil {
        return ErrInvalidURL
    }
    if u.Scheme != "http" && u.Scheme != "https" {
        return ErrInvalidURL
    }
    if u.Host == "" {
        return ErrInvalidURL
    }
    return nil
}

// CreateLink validates the URL and enforces the per-user link limit BEFORE
// ever calling the repository, then generates a unique short code,
// retrying on the rare collision.
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
        if err != nil {
            return domain.Link{}, fmt.Errorf("create link: %w", err)
        }
        link := domain.Link{Code: code, LongURL: strings.TrimSpace(longURL), OwnerID: ownerID, CreatedAt: time.Now()}
        err = s.repo.Create(link)
        if err == nil {
            return link, nil
        }
        if errors.Is(err, repository.ErrCodeTaken) {
            continue // rare collision - try another random code
        }
        return domain.Link{}, fmt.Errorf("create link: %w", err)
    }
    return domain.Link{}, fmt.Errorf("create link: %w", ErrCodeGeneration)
}

// Resolve looks up a Link by code and records a click. Read (Get) and
// write (IncrementClicks) are separate repository calls on purpose: a
// caller that only wants stats (Stats, below) must never trigger a click.
func (s *LinkService) Resolve(code string) (domain.Link, error) {
    link, err := s.repo.Get(code)
    if err != nil {
        return domain.Link{}, fmt.Errorf("resolve link: %w", ErrLinkNotFound)
    }
    if err := s.repo.IncrementClicks(code); err != nil {
        return domain.Link{}, fmt.Errorf("resolve link: %w", err)
    }
    link.Clicks++
    return link, nil
}

// Stats returns a Link's current data WITHOUT incrementing its click count.
func (s *LinkService) Stats(code string) (domain.Link, error) {
    link, err := s.repo.Get(code)
    if err != nil {
        return domain.Link{}, fmt.Errorf("stats: %w", ErrLinkNotFound)
    }
    return link, nil
}
EOF
```

4. Create `service/link_service_test.go`:

```bash
cat > service/link_service_test.go << 'EOF'
package service

import (
    "errors"
    "testing"

    "level40.example/urlshortener/domain"
)

// fakeLinkRepository is a small in-memory test double distinct from the
// real repository package, used to test LinkService entirely in isolation
// (Level 30's testing payoff).
type fakeLinkRepository struct {
    links map[string]domain.Link
}

func newFakeLinkRepository() *fakeLinkRepository {
    return &fakeLinkRepository{links: make(map[string]domain.Link)}
}

func (f *fakeLinkRepository) Create(l domain.Link) error {
    if _, exists := f.links[l.Code]; exists {
        return errors.New("code taken")
    }
    f.links[l.Code] = l
    return nil
}

func (f *fakeLinkRepository) Get(code string) (domain.Link, error) {
    l, ok := f.links[code]
    if !ok {
        return domain.Link{}, ErrLinkNotFound
    }
    return l, nil
}

func (f *fakeLinkRepository) IncrementClicks(code string) error {
    l, ok := f.links[code]
    if !ok {
        return ErrLinkNotFound
    }
    l.Clicks++
    f.links[code] = l
    return nil
}

func (f *fakeLinkRepository) CountByOwner(ownerID string) (int, error) {
    count := 0
    for _, l := range f.links {
        if l.OwnerID == ownerID {
            count++
        }
    }
    return count, nil
}

func (f *fakeLinkRepository) List() []domain.Link {
    out := make([]domain.Link, 0, len(f.links))
    for _, l := range f.links {
        out = append(out, l)
    }
    return out
}

var _ LinkRepository = (*fakeLinkRepository)(nil)

// poisonLinkRepository panics on every call - proves validation happens
// before the repository is ever touched (Level 30's poison-repository
// technique).
type poisonLinkRepository struct{}

func (poisonLinkRepository) Create(domain.Link) error        { panic("Create should not be called") }
func (poisonLinkRepository) Get(string) (domain.Link, error) { panic("Get should not be called") }
func (poisonLinkRepository) IncrementClicks(string) error {
    panic("IncrementClicks should not be called")
}
func (poisonLinkRepository) CountByOwner(string) (int, error) {
    panic("CountByOwner should not be called")
}
func (poisonLinkRepository) List() []domain.Link { panic("List should not be called") }

var _ LinkRepository = poisonLinkRepository{}

func TestLinkService_CreateLink_Success(t *testing.T) {
    svc := NewLinkService(newFakeLinkRepository(), 7)

    link, err := svc.CreateLink("u1", "https://example.com/very/long/path")
    if err != nil {
        t.Fatalf("CreateLink returned unexpected error: %v", err)
    }
    if len(link.Code) != 7 {
        t.Errorf("len(Code) = %d; want 7", len(link.Code))
    }
    if link.OwnerID != "u1" {
        t.Errorf("OwnerID = %q; want u1", link.OwnerID)
    }
}

func TestLinkService_CreateLink_InvalidURLRejectedBeforeRepository(t *testing.T) {
    svc := NewLinkService(poisonLinkRepository{}, 7)

    _, err := svc.CreateLink("u1", "not-a-url")

    if !errors.Is(err, ErrInvalidURL) {
        t.Errorf("error = %v; want it to wrap ErrInvalidURL", err)
    }
    // Reaching this line at all proves poisonLinkRepository never panicked.
}

func TestLinkService_CreateLink_EnforcesPerUserLimit(t *testing.T) {
    repo := newFakeLinkRepository()
    svc := NewLinkService(repo, 7)

    for i := 0; i < MaxLinksPerUser; i++ {
        if _, err := svc.CreateLink("u1", "https://example.com"); err != nil {
            t.Fatalf("CreateLink #%d returned unexpected error: %v", i, err)
        }
    }

    _, err := svc.CreateLink("u1", "https://example.com")

    if !errors.Is(err, ErrLinkLimit) {
        t.Errorf("error = %v; want it to wrap ErrLinkLimit", err)
    }
}

func TestLinkService_Resolve_IncrementsClicksAndStatsDoesNot(t *testing.T) {
    repo := newFakeLinkRepository()
    svc := NewLinkService(repo, 7)
    link, _ := svc.CreateLink("u1", "https://example.com")

    if _, err := svc.Stats(link.Code); err != nil {
        t.Fatalf("Stats returned unexpected error: %v", err)
    }
    statsAfter, _ := svc.Stats(link.Code)
    if statsAfter.Clicks != 0 {
        t.Errorf("Clicks after Stats = %d; want 0 (Stats must not increment)", statsAfter.Clicks)
    }

    resolved, err := svc.Resolve(link.Code)
    if err != nil {
        t.Fatalf("Resolve returned unexpected error: %v", err)
    }
    if resolved.Clicks != 1 {
        t.Errorf("Clicks after Resolve = %d; want 1", resolved.Clicks)
    }
}

func TestLinkService_Resolve_NotFound(t *testing.T) {
    svc := NewLinkService(newFakeLinkRepository(), 7)

    _, err := svc.Resolve("missing")

    if !errors.Is(err, ErrLinkNotFound) {
        t.Errorf("error = %v; want it to wrap ErrLinkNotFound", err)
    }
}
EOF
```

5. Run the service tests so far:

```bash
gofmt -l .
go vet ./...
go test ./service/... -v
```

**Expected Output:**

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
PASS
ok  	level40.example/urlshortener/service	0.571s
```

**Learning Objectives:**
- ✅ Generate unpredictable identifiers with `crypto/rand`, not `math/rand`
- ✅ Enforce a business rule (`MaxLinksPerUser`) that has nothing to do with storage
- ✅ Use a poison repository to PROVE validation runs before storage is ever touched
- ✅ Split "resolve and count a click" from "read stats without counting" at the Service layer

---

## Milestone 4: The Service Layer - Users, bcrypt, and JWTs

**Objective:** Add `UserService`: bcrypt password hashing on registration (Level 32) and JWT issuing/parsing on login (Level 32), depending only on a small repository interface.

**Instructions:**

1. Fetch the two third-party packages this milestone needs - like Level 32, this is a deliberate, flagged exception to the rest of this level's offline-by-default storage layer:

```bash
go get golang.org/x/crypto/bcrypt@v0.55.0
go get github.com/golang-jwt/jwt/v5@v5.3.1
```

2. Create `service/user_service.go`:

```bash
cat > service/user_service.go << 'EOF'
package service

import (
    "errors"
    "fmt"
    "strings"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"

    "level40.example/urlshortener/domain"
    "level40.example/urlshortener/repository"
)

// Re-exported for the same reason as Link, above.
type User = domain.User

var (
    ErrInvalidUsername    = errors.New("username must not be empty")
    ErrInvalidPassword    = errors.New("password must be at least 8 characters")
    ErrUsernameTaken      = repository.ErrUsernameTaken
    ErrInvalidCredentials = errors.New("invalid username or password")
)

// UserRepository is the CONSUMER-SIDE interface UserService depends on.
type UserRepository interface {
    Create(u domain.User) error
    GetByUsername(username string) (domain.User, error)
    GetByID(id string) (domain.User, error)
}

// UserService holds authentication business rules: hashing passwords
// (Level 32's bcrypt), issuing JWTs, and validating credentials. It never
// stores a plaintext password, and it depends only on the UserRepository
// interface.
type UserService struct {
    repo      UserRepository
    jwtSecret []byte
    tokenTTL  time.Duration
    newID     func() string
}

// NewUserService injects the repository and the JWT signing secret
// (Level 34: configuration - the secret comes from environment config,
// never a hardcoded literal, in main()).
func NewUserService(repo UserRepository, jwtSecret string, tokenTTL time.Duration) *UserService {
    return &UserService{
        repo:      repo,
        jwtSecret: []byte(jwtSecret),
        tokenTTL:  tokenTTL,
        newID:     func() string { return fmt.Sprintf("u_%d", time.Now().UnixNano()) },
    }
}

// Register validates the input, hashes the password with bcrypt (Level 32
// - never store or log a plaintext password), and creates the user.
func (s *UserService) Register(username, password string) (domain.User, error) {
    username = strings.TrimSpace(username)
    if username == "" {
        return domain.User{}, fmt.Errorf("register: %w", ErrInvalidUsername)
    }
    if len(password) < 8 {
        return domain.User{}, fmt.Errorf("register: %w", ErrInvalidPassword)
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return domain.User{}, fmt.Errorf("register: %w", err)
    }

    user := domain.User{
        ID:           s.newID(),
        Username:     username,
        PasswordHash: string(hash),
        CreatedAt:    time.Now(),
    }
    if err := s.repo.Create(user); err != nil {
        return domain.User{}, fmt.Errorf("register: %w", err)
    }
    return user, nil
}

// Login verifies credentials with bcrypt.CompareHashAndPassword (Level 32)
// and, on success, issues a signed JWT (HS256, Level 32) carrying the
// user's ID as the subject claim.
func (s *UserService) Login(username, password string) (string, error) {
    user, err := s.repo.GetByUsername(username)
    if err != nil {
        // Deliberately the SAME error for "no such user" and "wrong
        // password" below - never let a caller distinguish the two,
        // or they can enumerate valid usernames.
        return "", fmt.Errorf("login: %w", ErrInvalidCredentials)
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
        return "", fmt.Errorf("login: %w", ErrInvalidCredentials)
    }

    claims := jwt.RegisteredClaims{
        Subject:   user.ID,
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenTTL)),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signed, err := token.SignedString(s.jwtSecret)
    if err != nil {
        return "", fmt.Errorf("login: %w", err)
    }
    return signed, nil
}

// ParseUserID validates a signed JWT and returns the user ID it carries.
// It hard-codes acceptance of HMAC signing only (Level 32's alg:none
// defense) regardless of what the token's own header claims.
func (s *UserService) ParseUserID(tokenString string) (string, error) {
    token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
        }
        return s.jwtSecret, nil
    })
    if err != nil || !token.Valid {
        return "", fmt.Errorf("parse token: %w", ErrInvalidCredentials)
    }
    claims, ok := token.Claims.(*jwt.RegisteredClaims)
    if !ok {
        return "", fmt.Errorf("parse token: %w", ErrInvalidCredentials)
    }
    return claims.Subject, nil
}
EOF
```

3. Create `service/user_service_test.go`:

```bash
cat > service/user_service_test.go << 'EOF'
package service

import (
    "errors"
    "testing"
    "time"

    "level40.example/urlshortener/domain"
)

type fakeUserRepository struct {
    byID       map[string]domain.User
    byUsername map[string]domain.User
}

func newFakeUserRepository() *fakeUserRepository {
    return &fakeUserRepository{byID: map[string]domain.User{}, byUsername: map[string]domain.User{}}
}

func (f *fakeUserRepository) Create(u domain.User) error {
    if _, exists := f.byUsername[u.Username]; exists {
        return ErrUsernameTaken
    }
    f.byID[u.ID] = u
    f.byUsername[u.Username] = u
    return nil
}

func (f *fakeUserRepository) GetByUsername(username string) (domain.User, error) {
    u, ok := f.byUsername[username]
    if !ok {
        return domain.User{}, errors.New("not found")
    }
    return u, nil
}

func (f *fakeUserRepository) GetByID(id string) (domain.User, error) {
    u, ok := f.byID[id]
    if !ok {
        return domain.User{}, errors.New("not found")
    }
    return u, nil
}

var _ UserRepository = (*fakeUserRepository)(nil)

func TestUserService_RegisterAndLogin_Success(t *testing.T) {
    svc := NewUserService(newFakeUserRepository(), "test-secret", time.Hour)

    user, err := svc.Register("alice", "supersecret")
    if err != nil {
        t.Fatalf("Register returned unexpected error: %v", err)
    }
    if user.PasswordHash == "supersecret" {
        t.Fatal("PasswordHash equals the plaintext password - it must be hashed")
    }

    token, err := svc.Login("alice", "supersecret")
    if err != nil {
        t.Fatalf("Login returned unexpected error: %v", err)
    }
    if token == "" {
        t.Fatal("Login returned an empty token")
    }

    userID, err := svc.ParseUserID(token)
    if err != nil {
        t.Fatalf("ParseUserID returned unexpected error: %v", err)
    }
    if userID != user.ID {
        t.Errorf("ParseUserID = %q; want %q", userID, user.ID)
    }
}

func TestUserService_Register_ShortPasswordRejected(t *testing.T) {
    svc := NewUserService(newFakeUserRepository(), "test-secret", time.Hour)

    _, err := svc.Register("bob", "short")

    if !errors.Is(err, ErrInvalidPassword) {
        t.Errorf("error = %v; want it to wrap ErrInvalidPassword", err)
    }
}

func TestUserService_Login_WrongPasswordRejected(t *testing.T) {
    svc := NewUserService(newFakeUserRepository(), "test-secret", time.Hour)
    _, _ = svc.Register("alice", "supersecret")

    _, err := svc.Login("alice", "wrongpassword")

    if !errors.Is(err, ErrInvalidCredentials) {
        t.Errorf("error = %v; want it to wrap ErrInvalidCredentials", err)
    }
}

func TestUserService_Login_UnknownUserGetsSameError(t *testing.T) {
    svc := NewUserService(newFakeUserRepository(), "test-secret", time.Hour)

    _, err := svc.Login("ghost", "whatever1")

    if !errors.Is(err, ErrInvalidCredentials) {
        t.Errorf("error = %v; want it to wrap ErrInvalidCredentials (same as wrong password)", err)
    }
}

func TestUserService_ParseUserID_RejectsTokenSignedWithDifferentSecret(t *testing.T) {
    svcA := NewUserService(newFakeUserRepository(), "secret-a", time.Hour)
    repo := newFakeUserRepository()
    svcB := NewUserService(repo, "secret-b", time.Hour)
    _, _ = svcB.Register("alice", "supersecret")
    tokenFromB, _ := svcB.Login("alice", "supersecret")

    _, err := svcA.ParseUserID(tokenFromB)

    if !errors.Is(err, ErrInvalidCredentials) {
        t.Errorf("error = %v; want a token signed with a different secret to be rejected", err)
    }
}
EOF
```

4. Run the whole service package:

```bash
gofmt -l .
go vet ./...
go test ./service/... -v
```

**Expected Output:**

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

(The three `~0.15s` tests are bcrypt's deliberate slowness at work, Level 32's whole point.)

**Learning Objectives:**
- ✅ Hash passwords with bcrypt, never store or compare plaintext (Level 32)
- ✅ Issue and parse HS256 JWTs, hard-coding the accepted signing method (Level 32's `alg:none` defense)
- ✅ Return the SAME error for "no such user" and "wrong password" to prevent username enumeration
- ✅ Prove a token signed with a different secret is rejected - not just that a valid one is accepted

---

## Milestone 5: The Handler Layer

**Objective:** Build the Handler layer - `LinkHandler`, `AuthHandler`, and `NewMux` - using the stdlib `http.ServeMux` (Level 27). Test every handler against fakes with no-op middleware, before any real middleware exists.

**Instructions:**

1. Create `handler/context.go`:

```bash
cat > handler/context.go << 'EOF'
package handler

import "context"

// contextKey is an unexported type so keys from this package can never
// collide with a context key from another package (Level 24's rule).
type contextKey string

const userIDContextKey contextKey = "userID"

// WithUserID returns a new context carrying the authenticated user's ID.
// Called by the JWT auth middleware (middleware/auth.go) once a token has
// been validated.
func WithUserID(ctx context.Context, userID string) context.Context {
    return context.WithValue(ctx, userIDContextKey, userID)
}

// UserIDFromContext retrieves the authenticated user's ID stashed by the
// auth middleware. ok is false if no middleware ran (or the route is
// public).
func UserIDFromContext(ctx context.Context) (string, bool) {
    userID, ok := ctx.Value(userIDContextKey).(string)
    return userID, ok
}
EOF
```

2. Create `handler/link_handler.go`:

```bash
cat > handler/link_handler.go << 'EOF'
// Package handler is the Handler layer (Level 30): it translates HTTP
// requests into service calls and service results into HTTP responses.
// No validation, no storage logic - both live in the service and
// repository layers below it. Routing uses the stdlib http.ServeMux
// method-and-wildcard syntax (Level 27), no third-party router.
package handler

import (
    "encoding/json"
    "errors"
    "net/http"

    "level40.example/urlshortener/service"
)

// Link is re-exported purely for caller convenience.
type Link = service.Link

// LinkService is the CONSUMER-SIDE interface this package depends on -
// exactly the methods LinkHandler calls.
type LinkService interface {
    CreateLink(ownerID, longURL string) (Link, error)
    Resolve(code string) (Link, error)
    Stats(code string) (Link, error)
}

// linkStatusFor translates a Link-related error into an HTTP status code -
// the ONE place that knows about both errors and status codes (Level 30).
func linkStatusFor(err error) int {
    switch {
    case errors.Is(err, service.ErrInvalidURL):
        return http.StatusBadRequest
    case errors.Is(err, service.ErrLinkLimit):
        return http.StatusTooManyRequests
    case errors.Is(err, service.ErrLinkNotFound):
        return http.StatusNotFound
    default:
        return http.StatusInternalServerError
    }
}

type LinkHandler struct {
    service LinkService
}

func NewLinkHandler(s LinkService) *LinkHandler {
    return &LinkHandler{service: s}
}

type createLinkRequest struct {
    LongURL string `json:"long_url"`
}

// CreateLink handles POST /shorten. It is protected by the JWT auth
// middleware (Level 32) - by the time this runs, r.Context() already
// carries the authenticated user ID (see middleware/auth.go).
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

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(link)
}

// Redirect handles GET /{code} - PUBLIC, no auth required, but protected
// by the rate-limiter middleware (Level 39's token bucket) since it's the
// one endpoint the whole internet can hit.
func (h *LinkHandler) Redirect(w http.ResponseWriter, r *http.Request) {
    code := r.PathValue("code")
    link, err := h.service.Resolve(code)
    if err != nil {
        http.Error(w, err.Error(), linkStatusFor(err))
        return
    }
    http.Redirect(w, r, link.LongURL, http.StatusFound)
}

// Stats handles GET /stats/{code} - PUBLIC, read-only, never increments
// the click counter (that's the whole reason Resolve and Stats are
// separate service methods).
func (h *LinkHandler) Stats(w http.ResponseWriter, r *http.Request) {
    code := r.PathValue("code")
    link, err := h.service.Stats(code)
    if err != nil {
        http.Error(w, err.Error(), linkStatusFor(err))
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(link)
}
EOF
```

3. Create `handler/auth_handler.go`:

```bash
cat > handler/auth_handler.go << 'EOF'
package handler

import (
    "encoding/json"
    "errors"
    "net/http"

    "level40.example/urlshortener/service"
)

// User is re-exported purely for caller convenience.
type User = service.User

// AuthService is the CONSUMER-SIDE interface this package depends on.
type AuthService interface {
    Register(username, password string) (User, error)
    Login(username, password string) (string, error)
}

func authStatusFor(err error) int {
    switch {
    case errors.Is(err, service.ErrInvalidUsername), errors.Is(err, service.ErrInvalidPassword):
        return http.StatusBadRequest
    case errors.Is(err, service.ErrUsernameTaken):
        return http.StatusConflict
    case errors.Is(err, service.ErrInvalidCredentials):
        return http.StatusUnauthorized
    default:
        return http.StatusInternalServerError
    }
}

type AuthHandler struct {
    service AuthService
}

func NewAuthHandler(s AuthService) *AuthHandler {
    return &AuthHandler{service: s}
}

type credentialsRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

// Register handles POST /auth/register - PUBLIC.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req credentialsRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid JSON body", http.StatusBadRequest)
        return
    }

    user, err := h.service.Register(req.Username, req.Password)
    if err != nil {
        http.Error(w, err.Error(), authStatusFor(err))
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

type loginResponse struct {
    Token string `json:"token"`
}

// Login handles POST /auth/login - PUBLIC. Returns a signed JWT to use as
// a Bearer token on POST /shorten.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req credentialsRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid JSON body", http.StatusBadRequest)
        return
    }

    token, err := h.service.Login(req.Username, req.Password)
    if err != nil {
        http.Error(w, err.Error(), authStatusFor(err))
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(loginResponse{Token: token})
}
EOF
```

4. Create `handler/routes.go`:

```bash
cat > handler/routes.go << 'EOF'
package handler

import "net/http"

// Middleware is the standard func(http.Handler) http.Handler wrapping
// shape used throughout the course since Level 27. Accepting middleware
// as PARAMETERS here - rather than importing the middleware package
// directly - keeps the dependency direction one-way: middleware may
// depend on handler (it needs WithUserID/UserIDFromContext), but handler
// never depends on middleware (Level 30's Dependency Direction rule).
type Middleware func(http.Handler) http.Handler

// chain applies middlewares in order, so the FIRST one in the list runs
// OUTERMOST (first on the way in, last on the way out).
func chain(h http.Handler, mws ...Middleware) http.Handler {
    for i := len(mws) - 1; i >= 0; i-- {
        h = mws[i](h)
    }
    return h
}

// NewMux wires every route to its handler method, with the right
// middleware applied to each: logging on everything, JWT auth ONLY on
// POST /shorten, and the rate limiter ONLY on the public redirect.
func NewMux(linkH *LinkHandler, authH *AuthHandler, logging, requireAuth, rateLimit Middleware) *http.ServeMux {
    mux := http.NewServeMux()

    mux.Handle("POST /auth/register", chain(http.HandlerFunc(authH.Register), logging))
    mux.Handle("POST /auth/login", chain(http.HandlerFunc(authH.Login), logging))

    mux.Handle("POST /shorten", chain(http.HandlerFunc(linkH.CreateLink), logging, requireAuth))
    mux.Handle("GET /{code}", chain(http.HandlerFunc(linkH.Redirect), logging, rateLimit))
    mux.Handle("GET /stats/{code}", chain(http.HandlerFunc(linkH.Stats), logging))

    mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("ok"))
    })

    return mux
}
EOF
```

5. Create `handler/handler_test.go` - fakes for both services, plus identity/stub middleware since the REAL auth middleware doesn't exist until Milestone 6:

```bash
cat > handler/handler_test.go << 'EOF'
package handler

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "level40.example/urlshortener/service"
)

type fakeLinkService struct {
    createLink func(ownerID, longURL string) (Link, error)
    resolve    func(code string) (Link, error)
    stats      func(code string) (Link, error)
}

func (f *fakeLinkService) CreateLink(ownerID, longURL string) (Link, error) {
    return f.createLink(ownerID, longURL)
}
func (f *fakeLinkService) Resolve(code string) (Link, error) { return f.resolve(code) }
func (f *fakeLinkService) Stats(code string) (Link, error)   { return f.stats(code) }

var _ LinkService = (*fakeLinkService)(nil)

type fakeAuthService struct {
    register func(username, password string) (User, error)
    login    func(username, password string) (string, error)
}

func (f *fakeAuthService) Register(username, password string) (User, error) {
    return f.register(username, password)
}
func (f *fakeAuthService) Login(username, password string) (string, error) {
    return f.login(username, password)
}

var _ AuthService = (*fakeAuthService)(nil)

// identity is a no-op middleware, used where a test doesn't care about
// logging/auth/rate-limiting behavior - just the handler's own logic.
func identity(h http.Handler) http.Handler { return h }

// withFakeUser simulates what the REAL auth middleware would do once a
// valid token is presented: stash a user ID in the context.
func withFakeUser(userID string) Middleware {
    return func(h http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            h.ServeHTTP(w, r.WithContext(WithUserID(r.Context(), userID)))
        })
    }
}

func TestCreateLinkHandler_Success(t *testing.T) {
    linkSvc := &fakeLinkService{
        createLink: func(ownerID, longURL string) (Link, error) {
            return Link{Code: "abc1234", LongURL: longURL, OwnerID: ownerID}, nil
        },
    }
    mux := NewMux(NewLinkHandler(linkSvc), NewAuthHandler(&fakeAuthService{}), identity, withFakeUser("u1"), identity)

    body, _ := json.Marshal(map[string]string{"long_url": "https://example.com"})
    req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(body))
    rec := httptest.NewRecorder()

    mux.ServeHTTP(rec, req)

    if rec.Code != http.StatusCreated {
        t.Fatalf("status = %d; want %d, body = %s", rec.Code, http.StatusCreated, rec.Body.String())
    }
    var got Link
    if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
        t.Fatalf("could not decode response body: %v", err)
    }
    if got.Code != "abc1234" {
        t.Errorf("Code = %q; want abc1234", got.Code)
    }
}

func TestCreateLinkHandler_InvalidURL(t *testing.T) {
    linkSvc := &fakeLinkService{
        createLink: func(ownerID, longURL string) (Link, error) {
            return Link{}, service.ErrInvalidURL
        },
    }
    mux := NewMux(NewLinkHandler(linkSvc), NewAuthHandler(&fakeAuthService{}), identity, withFakeUser("u1"), identity)

    body, _ := json.Marshal(map[string]string{"long_url": "not-a-url"})
    req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(body))
    rec := httptest.NewRecorder()

    mux.ServeHTTP(rec, req)

    if rec.Code != http.StatusBadRequest {
        t.Errorf("status = %d; want %d", rec.Code, http.StatusBadRequest)
    }
}

func TestRedirectHandler_Success(t *testing.T) {
    linkSvc := &fakeLinkService{
        resolve: func(code string) (Link, error) {
            return Link{Code: code, LongURL: "https://example.com/target"}, nil
        },
    }
    mux := NewMux(NewLinkHandler(linkSvc), NewAuthHandler(&fakeAuthService{}), identity, identity, identity)

    req := httptest.NewRequest(http.MethodGet, "/abc1234", nil)
    rec := httptest.NewRecorder()

    mux.ServeHTTP(rec, req)

    if rec.Code != http.StatusFound {
        t.Fatalf("status = %d; want %d", rec.Code, http.StatusFound)
    }
    if loc := rec.Header().Get("Location"); loc != "https://example.com/target" {
        t.Errorf("Location = %q; want https://example.com/target", loc)
    }
}

func TestRedirectHandler_NotFound(t *testing.T) {
    linkSvc := &fakeLinkService{
        resolve: func(code string) (Link, error) { return Link{}, service.ErrLinkNotFound },
    }
    mux := NewMux(NewLinkHandler(linkSvc), NewAuthHandler(&fakeAuthService{}), identity, identity, identity)

    req := httptest.NewRequest(http.MethodGet, "/missing", nil)
    rec := httptest.NewRecorder()

    mux.ServeHTTP(rec, req)

    if rec.Code != http.StatusNotFound {
        t.Errorf("status = %d; want %d", rec.Code, http.StatusNotFound)
    }
}

func TestStatsHandler_DoesNotRequireAuth(t *testing.T) {
    linkSvc := &fakeLinkService{
        stats: func(code string) (Link, error) { return Link{Code: code, Clicks: 5}, nil },
    }
    mux := NewMux(NewLinkHandler(linkSvc), NewAuthHandler(&fakeAuthService{}), identity, identity, identity)

    req := httptest.NewRequest(http.MethodGet, "/stats/abc1234", nil)
    rec := httptest.NewRecorder()

    mux.ServeHTTP(rec, req)

    if rec.Code != http.StatusOK {
        t.Fatalf("status = %d; want %d, body = %s", rec.Code, http.StatusOK, rec.Body.String())
    }
    var got Link
    json.Unmarshal(rec.Body.Bytes(), &got)
    if got.Clicks != 5 {
        t.Errorf("Clicks = %d; want 5", got.Clicks)
    }
}

func TestRegisterHandler_UsernameTaken(t *testing.T) {
    authSvc := &fakeAuthService{
        register: func(username, password string) (User, error) { return User{}, service.ErrUsernameTaken },
    }
    mux := NewMux(NewLinkHandler(&fakeLinkService{}), NewAuthHandler(authSvc), identity, identity, identity)

    body, _ := json.Marshal(map[string]string{"username": "alice", "password": "supersecret"})
    req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
    rec := httptest.NewRecorder()

    mux.ServeHTTP(rec, req)

    if rec.Code != http.StatusConflict {
        t.Errorf("status = %d; want %d", rec.Code, http.StatusConflict)
    }
}

func TestLoginHandler_Success(t *testing.T) {
    authSvc := &fakeAuthService{
        login: func(username, password string) (string, error) { return "fake.jwt.token", nil },
    }
    mux := NewMux(NewLinkHandler(&fakeLinkService{}), NewAuthHandler(authSvc), identity, identity, identity)

    body, _ := json.Marshal(map[string]string{"username": "alice", "password": "supersecret"})
    req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
    rec := httptest.NewRecorder()

    mux.ServeHTTP(rec, req)

    if rec.Code != http.StatusOK {
        t.Fatalf("status = %d; want %d, body = %s", rec.Code, http.StatusOK, rec.Body.String())
    }
    var got loginResponse
    json.Unmarshal(rec.Body.Bytes(), &got)
    if got.Token != "fake.jwt.token" {
        t.Errorf("Token = %q; want fake.jwt.token", got.Token)
    }
}
EOF
```

6. Run the handler tests:

```bash
gofmt -l .
go vet ./...
go test ./handler/... -v
```

**Expected Output:**

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
ok  	level40.example/urlshortener/handler	0.663s
```

**Learning Objectives:**
- ✅ Build a thin Handler layer over `http.ServeMux`'s method-and-wildcard routing (Level 27)
- ✅ Accept middleware as constructor/function PARAMETERS to keep the dependency direction one-way
- ✅ Test handlers against fakes and stub middleware, entirely independent of any real service or auth logic

---

## Milestone 6: JWT Auth Middleware

**Objective:** Replace the stub `withFakeUser` middleware from Milestone 5 with the real thing: a `RequireAuth` middleware that verifies a Bearer JWT and injects the user ID (Level 32, Level 24).

**Instructions:**

1. Create `middleware/auth.go`:

```bash
cat > middleware/auth.go << 'EOF'
// Package middleware holds cross-cutting HTTP concerns - authentication,
// logging, rate limiting - each shaped as the standard
// func(http.Handler) http.Handler wrapper from Level 27, brought together
// here with Level 32's JWT verification and Level 24's context values.
package middleware

import (
    "net/http"
    "strings"

    "level40.example/urlshortener/handler"
)

// TokenParser is the CONSUMER-SIDE interface this middleware depends on -
// exactly the one method it calls. *service.UserService satisfies it
// implicitly.
type TokenParser interface {
    ParseUserID(token string) (string, error)
}

// RequireAuth extracts a Bearer token from the Authorization header,
// validates it, and injects the resulting user ID into the request
// context (handler.WithUserID) for every handler downstream. Missing or
// invalid tokens short-circuit with 401 before the wrapped handler ever
// runs.
func RequireAuth(parser TokenParser) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            const prefix = "Bearer "
            if !strings.HasPrefix(authHeader, prefix) {
                http.Error(w, "missing or malformed Authorization header", http.StatusUnauthorized)
                return
            }
            tokenString := strings.TrimPrefix(authHeader, prefix)

            userID, err := parser.ParseUserID(tokenString)
            if err != nil {
                http.Error(w, "invalid or expired token", http.StatusUnauthorized)
                return
            }

            ctx := handler.WithUserID(r.Context(), userID)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
EOF
```

2. Create `middleware/auth_test.go`:

```bash
cat > middleware/auth_test.go << 'EOF'
package middleware

import (
    "errors"
    "net/http"
    "net/http/httptest"
    "testing"

    "level40.example/urlshortener/handler"
)

type fakeTokenParser struct {
    userID string
    err    error
}

func (f *fakeTokenParser) ParseUserID(token string) (string, error) {
    if f.err != nil {
        return "", f.err
    }
    return f.userID, nil
}

func TestRequireAuth_ValidTokenInjectsUserID(t *testing.T) {
    var gotUserID string
    next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        gotUserID, _ = handler.UserIDFromContext(r.Context())
        w.WriteHeader(http.StatusOK)
    })

    wrapped := RequireAuth(&fakeTokenParser{userID: "u1"})(next)

    req := httptest.NewRequest(http.MethodPost, "/shorten", nil)
    req.Header.Set("Authorization", "Bearer valid.token.here")
    rec := httptest.NewRecorder()

    wrapped.ServeHTTP(rec, req)

    if rec.Code != http.StatusOK {
        t.Fatalf("status = %d; want %d", rec.Code, http.StatusOK)
    }
    if gotUserID != "u1" {
        t.Errorf("userID injected into context = %q; want u1", gotUserID)
    }
}

func TestRequireAuth_MissingHeaderRejected(t *testing.T) {
    next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        t.Fatal("next handler must not run without a valid token")
    })
    wrapped := RequireAuth(&fakeTokenParser{userID: "u1"})(next)

    req := httptest.NewRequest(http.MethodPost, "/shorten", nil)
    rec := httptest.NewRecorder()

    wrapped.ServeHTTP(rec, req)

    if rec.Code != http.StatusUnauthorized {
        t.Errorf("status = %d; want %d", rec.Code, http.StatusUnauthorized)
    }
}

func TestRequireAuth_InvalidTokenRejected(t *testing.T) {
    next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        t.Fatal("next handler must not run with an invalid token")
    })
    wrapped := RequireAuth(&fakeTokenParser{err: errors.New("bad token")})(next)

    req := httptest.NewRequest(http.MethodPost, "/shorten", nil)
    req.Header.Set("Authorization", "Bearer garbage")
    rec := httptest.NewRecorder()

    wrapped.ServeHTTP(rec, req)

    if rec.Code != http.StatusUnauthorized {
        t.Errorf("status = %d; want %d", rec.Code, http.StatusUnauthorized)
    }
}
EOF
```

3. Run the middleware tests:

```bash
gofmt -l .
go vet ./...
go test ./middleware/... -v
```

**Expected Output:**

```
=== RUN   TestRequireAuth_ValidTokenInjectsUserID
--- PASS: TestRequireAuth_ValidTokenInjectsUserID (0.00s)
=== RUN   TestRequireAuth_MissingHeaderRejected
--- PASS: TestRequireAuth_MissingHeaderRejected (0.00s)
=== RUN   TestRequireAuth_InvalidTokenRejected
--- PASS: TestRequireAuth_InvalidTokenRejected (0.00s)
PASS
ok  	level40.example/urlshortener/middleware	0.557s
```

**Learning Objectives:**
- ✅ Extract and verify a Bearer token in middleware, short-circuiting with 401 before the handler runs
- ✅ Inject an authenticated identity into `context.Context`, never into a global or a struct field
- ✅ Depend on a one-method consumer-side interface (`TokenParser`) so the real `UserService` and a test fake are interchangeable

---

## Milestone 7: Structured Logging and Configuration

**Objective:** Add a logging middleware (Level 33's `log/slog`) and a `config` package (Level 34) that fails fast on a missing `JWT_SECRET`.

**Instructions:**

1. Create `middleware/logging.go`:

```bash
cat > middleware/logging.go << 'EOF'
package middleware

import (
    "log/slog"
    "net/http"
    "time"
)

// statusRecorder wraps http.ResponseWriter to capture the status code a
// handler wrote, since http.ResponseWriter has no getter for it.
type statusRecorder struct {
    http.ResponseWriter
    status int
}

func (r *statusRecorder) WriteHeader(code int) {
    r.status = code
    r.ResponseWriter.WriteHeader(code)
}

// Logging returns a middleware that logs one structured line per request
// (Level 33's log/slog) - method, path, status, and latency - using the
// injected logger rather than a package-level global (Level 31's
// dependency-injection habit carried into middleware).
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
EOF
```

2. Create `middleware/logging_test.go`:

```bash
cat > middleware/logging_test.go << 'EOF'
package middleware

import (
    "bytes"
    "encoding/json"
    "log/slog"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestLogging_EmitsOneStructuredLineWithStatusAndPath(t *testing.T) {
    var buf bytes.Buffer
    logger := slog.New(slog.NewJSONHandler(&buf, nil))

    next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusTeapot)
    })
    wrapped := Logging(logger)(next)

    req := httptest.NewRequest(http.MethodGet, "/brew", nil)
    rec := httptest.NewRecorder()
    wrapped.ServeHTTP(rec, req)

    var entry map[string]any
    if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
        t.Fatalf("could not decode log line as JSON: %v\nraw: %s", err, buf.String())
    }
    if entry["msg"] != "http request" {
        t.Errorf("msg = %v; want %q", entry["msg"], "http request")
    }
    if entry["path"] != "/brew" {
        t.Errorf("path = %v; want /brew", entry["path"])
    }
    if entry["status"] != float64(http.StatusTeapot) {
        t.Errorf("status = %v; want %d", entry["status"], http.StatusTeapot)
    }
}
EOF
```

3. Create `config/config.go`:

```bash
cat > config/config.go << 'EOF'
// Package config collects every environment-driven setting into one typed
// struct, populated by a single LoadConfig function (Level 34's pattern).
// The rest of the program depends only on Config - never on os.Getenv
// directly.
package config

import (
    "fmt"
    "os"
    "strconv"
    "time"
)

// Config holds every setting the service needs at startup.
type Config struct {
    Port            string
    JWTSecret       string
    TokenTTL        time.Duration
    CodeLength      int
    RateLimitPerSec float64
    RateLimitBurst  int
}

// LoadConfig reads settings from environment variables, applying sensible
// defaults for local development and FAILING FAST (Level 34's Best
// Practice #3) if a required, security-sensitive setting like
// JWT_SECRET is missing outside of explicit local-dev opt-in.
func LoadConfig() (Config, error) {
    cfg := Config{
        Port:            getEnv("PORT", "8080"),
        JWTSecret:       os.Getenv("JWT_SECRET"),
        TokenTTL:        time.Hour,
        CodeLength:      7,
        RateLimitPerSec: 5,
        RateLimitBurst:  10,
    }

    if v, ok := os.LookupEnv("TOKEN_TTL_MINUTES"); ok {
        minutes, err := strconv.Atoi(v)
        if err != nil {
            return Config{}, fmt.Errorf("load config: invalid TOKEN_TTL_MINUTES %q: %w", v, err)
        }
        cfg.TokenTTL = time.Duration(minutes) * time.Minute
    }

    if v, ok := os.LookupEnv("SHORT_CODE_LENGTH"); ok {
        n, err := strconv.Atoi(v)
        if err != nil {
            return Config{}, fmt.Errorf("load config: invalid SHORT_CODE_LENGTH %q: %w", v, err)
        }
        cfg.CodeLength = n
    }

    if v, ok := os.LookupEnv("RATE_LIMIT_PER_SEC"); ok {
        f, err := strconv.ParseFloat(v, 64)
        if err != nil {
            return Config{}, fmt.Errorf("load config: invalid RATE_LIMIT_PER_SEC %q: %w", v, err)
        }
        cfg.RateLimitPerSec = f
    }

    if v, ok := os.LookupEnv("RATE_LIMIT_BURST"); ok {
        n, err := strconv.Atoi(v)
        if err != nil {
            return Config{}, fmt.Errorf("load config: invalid RATE_LIMIT_BURST %q: %w", v, err)
        }
        cfg.RateLimitBurst = n
    }

    if cfg.JWTSecret == "" {
        if os.Getenv("ALLOW_DEV_JWT_SECRET") == "true" {
            cfg.JWTSecret = "dev-only-insecure-secret"
        } else {
            return Config{}, fmt.Errorf("load config: JWT_SECRET is required (set ALLOW_DEV_JWT_SECRET=true for local development only)")
        }
    }
    if cfg.CodeLength < 4 {
        return Config{}, fmt.Errorf("load config: SHORT_CODE_LENGTH must be at least 4, got %d", cfg.CodeLength)
    }

    return cfg, nil
}

func getEnv(key, fallback string) string {
    if v, ok := os.LookupEnv(key); ok {
        return v
    }
    return fallback
}
EOF
```

4. Create `config/config_test.go` - **including a test that catches a real bug**, on purpose, kept exactly as first written to show what the failure looked like before the fix below:

```bash
cat > config/config_test.go << 'EOF'
package config

import (
    "os"
    "testing"
)

func TestLoadConfig_FailsFastWithoutJWTSecret(t *testing.T) {
    t.Setenv("JWT_SECRET", "")
    t.Setenv("ALLOW_DEV_JWT_SECRET", "")

    _, err := LoadConfig()

    if err == nil {
        t.Fatal("LoadConfig succeeded with no JWT_SECRET and no dev opt-in; want an error")
    }
}

func TestLoadConfig_DevOptInAllowsMissingSecret(t *testing.T) {
    t.Setenv("JWT_SECRET", "")
    t.Setenv("ALLOW_DEV_JWT_SECRET", "true")

    cfg, err := LoadConfig()

    if err != nil {
        t.Fatalf("LoadConfig returned unexpected error: %v", err)
    }
    if cfg.JWTSecret == "" {
        t.Error("JWTSecret is empty even with ALLOW_DEV_JWT_SECRET=true")
    }
}

func TestLoadConfig_DefaultsApplied(t *testing.T) {
    t.Setenv("JWT_SECRET", "real-secret")
    t.Setenv("PORT", "")

    cfg, err := LoadConfig()

    if err != nil {
        t.Fatalf("LoadConfig returned unexpected error: %v", err)
    }
    if cfg.Port != "8080" {
        t.Errorf("Port = %q; want default 8080", cfg.Port)
    }
    if cfg.CodeLength != 7 {
        t.Errorf("CodeLength = %d; want default 7", cfg.CodeLength)
    }
}

func TestLoadConfig_InvalidShortCodeLengthRejected(t *testing.T) {
    t.Setenv("JWT_SECRET", "real-secret")
    t.Setenv("SHORT_CODE_LENGTH", "2")

    _, err := LoadConfig()

    if err == nil {
        t.Fatal("LoadConfig succeeded with SHORT_CODE_LENGTH=2; want an error (below minimum of 4)")
    }
}
EOF
```

5. Run it - this fails, for real, and the failure is instructive:

```bash
go test ./config/... -v
```

**Actual Output (a genuine bug, caught by the test as written above):**

```
=== RUN   TestLoadConfig_FailsFastWithoutJWTSecret
--- PASS: TestLoadConfig_FailsFastWithoutJWTSecret (0.00s)
=== RUN   TestLoadConfig_DevOptInAllowsMissingSecret
--- PASS: TestLoadConfig_DevOptInAllowsMissingSecret (0.00s)
=== RUN   TestLoadConfig_DefaultsApplied
    config_test.go:40: Port = ""; want default 8080
--- FAIL: TestLoadConfig_DefaultsApplied (0.00s)
=== RUN   TestLoadConfig_InvalidShortCodeLengthRejected
--- PASS: TestLoadConfig_InvalidShortCodeLengthRejected (0.00s)
FAIL
FAIL	level40.example/urlshortener/config	1.119s
```

**What went wrong:** `t.Setenv("PORT", "")` sets `PORT` to an explicitly empty string - which `os.LookupEnv` correctly reports as **present** (`ok == true`), not unset. That's exactly Level 34's Section 2 lesson (`os.Getenv` vs. `os.LookupEnv`, "unset" vs. "empty"): the test author, not `LoadConfig`, confused the two. `getEnv`'s fallback only ever triggers when the variable is truly unset. Fix the TEST, not the code:

```bash
```

Open `config/config_test.go` and replace the `TestLoadConfig_DefaultsApplied` body's second line:

```go
// ❌ WRONG - sets PORT to a PRESENT empty string, which os.LookupEnv
// correctly reports as ok=true, so the default is never applied
t.Setenv("PORT", "")

// ✅ RIGHT - PORT must be genuinely UNSET for the default to kick in
os.Unsetenv("PORT")
```

Re-run:

```bash
go test ./config/... -v
```

**Expected Output (after the fix):**

```
=== RUN   TestLoadConfig_FailsFastWithoutJWTSecret
--- PASS: TestLoadConfig_FailsFastWithoutJWTSecret (0.00s)
=== RUN   TestLoadConfig_DevOptInAllowsMissingSecret
--- PASS: TestLoadConfig_DevOptInAllowsMissingSecret (0.00s)
=== RUN   TestLoadConfig_DefaultsApplied
--- PASS: TestLoadConfig_DefaultsApplied (0.00s)
=== RUN   TestLoadConfig_InvalidShortCodeLengthRejected
--- PASS: TestLoadConfig_InvalidShortCodeLengthRejected (0.00s)
PASS
ok  	level40.example/urlshortener/config	0.561s
```

6. Run the logging middleware test too:

```bash
go test ./middleware/... -v
```

**Expected Output:**

```
=== RUN   TestRequireAuth_ValidTokenInjectsUserID
--- PASS: TestRequireAuth_ValidTokenInjectsUserID (0.00s)
=== RUN   TestRequireAuth_MissingHeaderRejected
--- PASS: TestRequireAuth_MissingHeaderRejected (0.00s)
=== RUN   TestRequireAuth_InvalidTokenRejected
--- PASS: TestRequireAuth_InvalidTokenRejected (0.00s)
=== RUN   TestLogging_EmitsOneStructuredLineWithStatusAndPath
--- PASS: TestLogging_EmitsOneStructuredLineWithStatusAndPath (0.00s)
PASS
ok  	level40.example/urlshortener/middleware	0.609s
```

**Learning Objectives:**
- ✅ Emit one structured JSON log line per request via an INJECTED `*slog.Logger` (Level 33, Level 31)
- ✅ Fail fast at startup on a missing security-sensitive setting, with an explicit dev opt-in (Level 34)
- ✅ Experience firsthand why `os.LookupEnv("")` reports "present," not "unset" - by writing a test that got it wrong, then fixing the test

---

## Milestone 8: Token-Bucket Rate Limiter

**Objective:** Bridge Level 39's rate-limiting coverage into real, runnable middleware protecting the public redirect endpoint.

**Instructions:**

1. Create `ratelimit/tokenbucket.go`:

```bash
cat > ratelimit/tokenbucket.go << 'EOF'
// Package ratelimit implements a token-bucket rate limiter (bridging
// Level 39: System Design's coverage of rate-limiting algorithms) as
// reusable, dependency-free HTTP middleware.
package ratelimit

import (
    "net/http"
    "sync"
    "time"
)

// bucket tracks one client's available tokens and when it was last
// refilled. Guarded by the Limiter's mutex (Level 23), since concurrent
// requests from the same client hit the same bucket.
type bucket struct {
    tokens     float64
    lastRefill time.Time
}

// Limiter is a token-bucket rate limiter keyed by an arbitrary string
// (typically the client's IP address). ratePerSecond tokens are added
// back per second, up to burst - the classic token-bucket shape: it
// allows short bursts up to `burst` while enforcing a steady-state
// average of `ratePerSecond`.
type Limiter struct {
    mu            sync.Mutex
    buckets       map[string]*bucket
    ratePerSecond float64
    burst         int
    now           func() time.Time // injected for deterministic tests
}

// NewLimiter builds a Limiter. burst is both the bucket's maximum
// capacity and its starting balance, so the very first request from any
// client is never throttled.
func NewLimiter(ratePerSecond float64, burst int) *Limiter {
    return &Limiter{
        buckets:       make(map[string]*bucket),
        ratePerSecond: ratePerSecond,
        burst:         burst,
        now:           time.Now,
    }
}

// Allow reports whether the given key has a token available, consuming
// one if so. Safe for concurrent use.
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
    b.tokens += elapsed * l.ratePerSecond
    if b.tokens > float64(l.burst) {
        b.tokens = float64(l.burst)
    }
    b.lastRefill = now

    if b.tokens < 1 {
        return false
    }
    b.tokens--
    return true
}

// Middleware wraps an http.Handler, rejecting requests over the limit
// with 429 Too Many Requests before the wrapped handler ever runs. keyFunc
// extracts the rate-limit key from a request - typically the client's
// remote address - kept as a parameter (constructor-injection-style,
// Level 31) so it's easy to swap or fake in tests.
func (l *Limiter) Middleware(keyFunc func(*http.Request) string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            key := keyFunc(r)
            if !l.Allow(key) {
                http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

// KeyByRemoteAddr is the default keyFunc: rate-limit by client IP.
func KeyByRemoteAddr(r *http.Request) string {
    return r.RemoteAddr
}
EOF
```

2. Create `ratelimit/tokenbucket_test.go`:

```bash
cat > ratelimit/tokenbucket_test.go << 'EOF'
package ratelimit

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
)

func TestLimiter_AllowsUpToBurstThenBlocks(t *testing.T) {
    l := NewLimiter(1, 3) // 1 token/sec refill, burst of 3

    for i := 0; i < 3; i++ {
        if !l.Allow("client-a") {
            t.Fatalf("request %d should be allowed within burst of 3", i+1)
        }
    }
    if l.Allow("client-a") {
        t.Fatal("4th immediate request should be blocked - burst exhausted")
    }
}

func TestLimiter_RefillsOverTime(t *testing.T) {
    l := NewLimiter(1, 1) // 1 token/sec, burst of 1
    fakeNow := time.Now()
    l.now = func() time.Time { return fakeNow }

    if !l.Allow("client-a") {
        t.Fatal("first request should be allowed")
    }
    if l.Allow("client-a") {
        t.Fatal("second immediate request should be blocked")
    }

    fakeNow = fakeNow.Add(1100 * time.Millisecond) // advance past one refill interval
    if !l.Allow("client-a") {
        t.Fatal("request after refill interval should be allowed")
    }
}

func TestLimiter_TracksKeysIndependently(t *testing.T) {
    l := NewLimiter(1, 1)

    if !l.Allow("client-a") {
        t.Fatal("client-a's first request should be allowed")
    }
    if !l.Allow("client-b") {
        t.Fatal("client-b's first request should be allowed independently of client-a")
    }
}

func TestMiddleware_RejectsOverLimitWith429(t *testing.T) {
    l := NewLimiter(1, 1)
    next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })
    wrapped := l.Middleware(KeyByRemoteAddr)(next)

    req := httptest.NewRequest(http.MethodGet, "/abc1234", nil)
    req.RemoteAddr = "203.0.113.5:1234"

    rec1 := httptest.NewRecorder()
    wrapped.ServeHTTP(rec1, req)
    if rec1.Code != http.StatusOK {
        t.Fatalf("first request status = %d; want %d", rec1.Code, http.StatusOK)
    }

    rec2 := httptest.NewRecorder()
    wrapped.ServeHTTP(rec2, req)
    if rec2.Code != http.StatusTooManyRequests {
        t.Fatalf("second request status = %d; want %d", rec2.Code, http.StatusTooManyRequests)
    }
}
EOF
```

3. Run the rate limiter tests:

```bash
gofmt -l .
go vet ./...
go test ./ratelimit/... -v
```

**Expected Output:**

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
ok  	level40.example/urlshortener/ratelimit	0.663s
```

**Learning Objectives:**
- ✅ Implement a token-bucket rate limiter from scratch: refill rate, burst capacity, per-key buckets
- ✅ Inject the clock (`now func() time.Time`) so a refill can be tested deterministically instead of with a real `time.Sleep`
- ✅ Wrap it as `func(http.Handler) http.Handler` middleware consistent with every other middleware in this project

---

## Milestone 9: Wiring It All Together and the Full End-to-End Test

**Objective:** Build the composition root (`main.go`) that constructs the real dependency graph, then prove the ENTIRE system works together with one `httptest.NewServer`-based end-to-end test - not just each layer in isolation.

**Instructions:**

1. Create `main.go`:

```bash
cat > main.go << 'EOF'
// Command urlshortener is the capstone project: a URL shortener API that
// wires together every layer built across this course - Repository (30),
// Service, Handler, Dependency Injection (31), JWT auth (32), structured
// logging (33), configuration (34), and a token-bucket rate limiter
// (bridging 39) - into one working service.
package main

import (
    "log/slog"
    "net/http"
    "os"

    "level40.example/urlshortener/config"
    "level40.example/urlshortener/handler"
    "level40.example/urlshortener/middleware"
    "level40.example/urlshortener/ratelimit"
    "level40.example/urlshortener/repository"
    "level40.example/urlshortener/service"
)

// buildServer is the composition root (Level 31): the ONE place that
// knows how the entire dependency graph fits together. Everything above
// this function only ever sees interfaces and plain config values.
func buildServer(cfg config.Config, logger *slog.Logger) http.Handler {
    // Repository layer - concrete, in-memory (Level 30).
    linkRepo := repository.NewInMemoryLinkRepository()
    userRepo := repository.NewInMemoryUserRepository()

    // Service layer - depends on repository INTERFACES only.
    linkService := service.NewLinkService(linkRepo, cfg.CodeLength)
    userService := service.NewUserService(userRepo, cfg.JWTSecret, cfg.TokenTTL)

    // Handler layer - depends on service INTERFACES only.
    linkHandler := handler.NewLinkHandler(linkService)
    authHandler := handler.NewAuthHandler(userService)

    // Middleware - logging on everything, JWT auth on POST /shorten,
    // token-bucket rate limiting on the public redirect.
    logging := middleware.Logging(logger)
    requireAuth := middleware.RequireAuth(userService)
    limiter := ratelimit.NewLimiter(cfg.RateLimitPerSec, cfg.RateLimitBurst)
    rateLimit := limiter.Middleware(ratelimit.KeyByRemoteAddr)

    return handler.NewMux(linkHandler, authHandler, logging, requireAuth, rateLimit)
}

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

    cfg, err := config.LoadConfig()
    if err != nil {
        logger.Error("failed to load config", slog.String("error", err.Error()))
        os.Exit(1)
    }

    mux := buildServer(cfg, logger)

    logger.Info("starting server", slog.String("port", cfg.Port))
    if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
        logger.Error("server stopped", slog.String("error", err.Error()))
        os.Exit(1)
    }
}
EOF
```

2. Confirm it builds:

```bash
go build ./...
echo BUILD_OK
```

**Expected Output:**

```
BUILD_OK
```

3. Create `main_test.go` - the full end-to-end test:

```bash
cat > main_test.go << 'EOF'
package main

import (
    "bytes"
    "encoding/json"
    "io"
    "log/slog"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "level40.example/urlshortener/config"
)

// newTestServer builds the FULL real dependency graph (repository ->
// service -> handler -> middleware) exactly as main() does, then wraps it
// in httptest.NewServer for a real TCP round trip over loopback.
func newTestServer(t *testing.T) *httptest.Server {
    t.Helper()
    cfg := config.Config{
        Port:            "0",
        JWTSecret:       "test-secret",
        TokenTTL:        time.Hour,
        CodeLength:      7,
        RateLimitPerSec: 2,
        RateLimitBurst:  2,
    }
    logger := slog.New(slog.NewTextHandler(io.Discard, nil))
    mux := buildServer(cfg, logger)
    srv := httptest.NewServer(mux)
    t.Cleanup(srv.Close)
    return srv
}

func postJSON(t *testing.T, url string, payload any, token string) *http.Response {
    t.Helper()
    body, err := json.Marshal(payload)
    if err != nil {
        t.Fatalf("json.Marshal: %v", err)
    }
    req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
    if err != nil {
        t.Fatalf("http.NewRequest: %v", err)
    }
    req.Header.Set("Content-Type", "application/json")
    if token != "" {
        req.Header.Set("Authorization", "Bearer "+token)
    }
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        t.Fatalf("http.Do: %v", err)
    }
    return resp
}

// TestEndToEnd_FullUserJourney exercises the REAL HTTP flow a real client
// would follow: register, log in, create a short link (protected by JWT),
// follow the public redirect, and read public stats.
func TestEndToEnd_FullUserJourney(t *testing.T) {
    srv := newTestServer(t)

    registerResp := postJSON(t, srv.URL+"/auth/register", map[string]string{
        "username": "alice", "password": "supersecret1",
    }, "")
    if registerResp.StatusCode != http.StatusCreated {
        t.Fatalf("register status = %d; want %d", registerResp.StatusCode, http.StatusCreated)
    }
    registerResp.Body.Close()

    loginResp := postJSON(t, srv.URL+"/auth/login", map[string]string{
        "username": "alice", "password": "supersecret1",
    }, "")
    if loginResp.StatusCode != http.StatusOK {
        t.Fatalf("login status = %d; want %d", loginResp.StatusCode, http.StatusOK)
    }
    var loginBody struct {
        Token string `json:"token"`
    }
    json.NewDecoder(loginResp.Body).Decode(&loginBody)
    loginResp.Body.Close()
    if loginBody.Token == "" {
        t.Fatal("login returned an empty token")
    }

    shortenResp := postJSON(t, srv.URL+"/shorten", map[string]string{
        "long_url": "https://go.dev/doc/effective_go",
    }, loginBody.Token)
    if shortenResp.StatusCode != http.StatusCreated {
        body, _ := io.ReadAll(shortenResp.Body)
        t.Fatalf("shorten status = %d; want %d, body = %s", shortenResp.StatusCode, http.StatusCreated, body)
    }
    var link struct {
        Code    string `json:"code"`
        LongURL string `json:"long_url"`
    }
    json.NewDecoder(shortenResp.Body).Decode(&link)
    shortenResp.Body.Close()
    if len(link.Code) != 7 {
        t.Fatalf("len(Code) = %d; want 7", len(link.Code))
    }

    noAuthResp := postJSON(t, srv.URL+"/shorten", map[string]string{"long_url": "https://example.com"}, "")
    if noAuthResp.StatusCode != http.StatusUnauthorized {
        t.Errorf("unauthenticated shorten status = %d; want %d", noAuthResp.StatusCode, http.StatusUnauthorized)
    }
    noAuthResp.Body.Close()

    noRedirectClient := &http.Client{
        CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
    }
    redirectResp, err := noRedirectClient.Get(srv.URL + "/" + link.Code)
    if err != nil {
        t.Fatalf("GET /%s: %v", link.Code, err)
    }
    if redirectResp.StatusCode != http.StatusFound {
        t.Fatalf("redirect status = %d; want %d", redirectResp.StatusCode, http.StatusFound)
    }
    if loc := redirectResp.Header.Get("Location"); loc != link.LongURL {
        t.Errorf("Location = %q; want %q", loc, link.LongURL)
    }
    redirectResp.Body.Close()

    statsResp, err := http.Get(srv.URL + "/stats/" + link.Code)
    if err != nil {
        t.Fatalf("GET /stats/%s: %v", link.Code, err)
    }
    var stats struct {
        Clicks int `json:"clicks"`
    }
    json.NewDecoder(statsResp.Body).Decode(&stats)
    statsResp.Body.Close()
    if stats.Clicks != 1 {
        t.Errorf("Clicks = %d; want 1 (one redirect happened above)", stats.Clicks)
    }
}

// TestEndToEnd_RateLimiterProtectsRedirect proves the public redirect
// endpoint really is rate-limited: with a burst of 2, the 3rd immediate
// request from the same client must be rejected with 429.
func TestEndToEnd_RateLimiterProtectsRedirect(t *testing.T) {
    srv := newTestServer(t)

    registerResp := postJSON(t, srv.URL+"/auth/register", map[string]string{"username": "bob", "password": "supersecret1"}, "")
    registerResp.Body.Close()
    loginResp := postJSON(t, srv.URL+"/auth/login", map[string]string{"username": "bob", "password": "supersecret1"}, "")
    var loginBody struct {
        Token string `json:"token"`
    }
    json.NewDecoder(loginResp.Body).Decode(&loginBody)
    loginResp.Body.Close()

    shortenResp := postJSON(t, srv.URL+"/shorten", map[string]string{"long_url": "https://example.com"}, loginBody.Token)
    var link struct {
        Code string `json:"code"`
    }
    json.NewDecoder(shortenResp.Body).Decode(&link)
    shortenResp.Body.Close()

    client := &http.Client{
        CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
    }

    var statuses []int
    for i := 0; i < 3; i++ {
        resp, err := client.Get(srv.URL + "/" + link.Code)
        if err != nil {
            t.Fatalf("request %d: %v", i+1, err)
        }
        statuses = append(statuses, resp.StatusCode)
        resp.Body.Close()
    }

    if statuses[0] != http.StatusFound || statuses[1] != http.StatusFound {
        t.Errorf("first two requests = %v; want both %d (within burst of 2)", statuses[:2], http.StatusFound)
    }
    if statuses[2] != http.StatusTooManyRequests {
        t.Errorf("third request = %d; want %d (burst exhausted)", statuses[2], http.StatusTooManyRequests)
    }
}
EOF
```

4. Run it - **this fails on the second test**, and the failure teaches an HTTP fundamental:

```bash
go test . -run TestEndToEnd -v -count=1
```

**Actual Output (a second genuine bug):**

```
=== RUN   TestEndToEnd_FullUserJourney
--- PASS: TestEndToEnd_FullUserJourney (0.16s)
=== RUN   TestEndToEnd_RateLimiterProtectsRedirect
    main_test.go:201: third request = 302; want 429 (burst exhausted)
--- FAIL: TestEndToEnd_RateLimiterProtectsRedirect (0.15s)
FAIL
FAIL	level40.example/urlshortener	0.604s
```

**What went wrong:** the loop closed each response body without first reading it to EOF. Go's `http.Client` can only reuse a persistent (keep-alive) TCP connection when the previous response body has been fully drained before `Close()` - otherwise it opens a **new** connection (and a new ephemeral client port) for the next request. Since the rate limiter's key is `r.RemoteAddr`, a new connection means a new key, which means a brand-new, full bucket every single request - so the limiter never actually saw three requests from "the same client." Fix: drain the body before closing.

```go
// ❌ WRONG - skips draining, so the client opens a NEW connection (and a
// NEW RemoteAddr) for every request, giving each one a fresh rate-limit
// bucket
resp, err := client.Get(srv.URL + "/" + link.Code)
statuses = append(statuses, resp.StatusCode)
resp.Body.Close()

// ✅ RIGHT - draining the body lets the client REUSE the connection, so
// RemoteAddr (and therefore the rate-limit bucket) stays the SAME across
// requests, the way it would for a real repeat client
resp, err := client.Get(srv.URL + "/" + link.Code)
statuses = append(statuses, resp.StatusCode)
io.Copy(io.Discard, resp.Body)
resp.Body.Close()
```

Apply that fix inside the loop in `TestEndToEnd_RateLimiterProtectsRedirect`, then re-run:

```bash
go test . -run TestEndToEnd -v -count=1
```

**Expected Output (after the fix):**

```
=== RUN   TestEndToEnd_FullUserJourney
--- PASS: TestEndToEnd_FullUserJourney (0.16s)
=== RUN   TestEndToEnd_RateLimiterProtectsRedirect
--- PASS: TestEndToEnd_RateLimiterProtectsRedirect (0.15s)
PASS
ok  	level40.example/urlshortener	0.999s
```

5. Run the ENTIRE module's test suite - every package, every layer, one command:

```bash
gofmt -l .
go vet ./...
go build ./...
go test ./... -count=1
```

**Expected Output:**

```
ok  	level40.example/urlshortener	2.492s
ok  	level40.example/urlshortener/config	0.257s
?   	level40.example/urlshortener/domain	[no test files]
ok  	level40.example/urlshortener/handler	1.628s
ok  	level40.example/urlshortener/middleware	1.090s
ok  	level40.example/urlshortener/ratelimit	0.547s
ok  	level40.example/urlshortener/repository	2.880s
ok  	level40.example/urlshortener/service	3.156s
```

6. Finally, run the REAL server and hit it with real `curl` commands - the ultimate proof that this is a working service, not just passing tests:

```bash
go build -o urlshortener-bin .
JWT_SECRET="demo-secret-for-live-run" PORT=8099 RATE_LIMIT_PER_SEC=5 RATE_LIMIT_BURST=10 \
  ./urlshortener-bin &
SERVER_PID=$!
sleep 1

curl -s -w "\nHTTP %{http_code}\n" -X POST http://localhost:8099/auth/register \
  -d '{"username":"alice","password":"supersecret1"}'

LOGIN_JSON=$(curl -s -X POST http://localhost:8099/auth/login \
  -d '{"username":"alice","password":"supersecret1"}')
TOKEN=$(echo "$LOGIN_JSON" | sed -E 's/.*"token":"([^"]+)".*/\1/')

curl -s -w "\nHTTP %{http_code}\n" -X POST http://localhost:8099/shorten \
  -H "Authorization: Bearer $TOKEN" -d '{"long_url":"https://go.dev/doc/effective_go"}'

curl -s -w "\nHTTP %{http_code}\n" -X POST http://localhost:8099/shorten \
  -d '{"long_url":"https://example.com"}'

kill $SERVER_PID
```

**Expected Output (real, captured):**

```
{"id":"u_1788285897194322000","username":"alice","created_at":"2026-09-01T23:34:57.194323+05:30"}

HTTP 201
{"code":"YXD97v0","long_url":"https://go.dev/doc/effective_go","owner_id":"u_1788285897194322000","clicks":0,"created_at":"2026-09-01T23:34:57.292782+05:30"}

HTTP 201
missing or malformed Authorization header

HTTP 401
```

**Learning Objectives:**
- ✅ Build a composition root (`buildServer`) that wires the entire dependency graph by hand (Level 31, at application scale)
- ✅ Write one end-to-end test that drives a real HTTP server through the real flow a client would follow (Level 26, Level 30's Section 6 technique)
- ✅ Learn - by hitting it directly - why HTTP keep-alive connection reuse depends on fully draining a response body, and why that matters for anything keyed by `RemoteAddr`
- ✅ Confirm the whole system with real `go test` output AND a real running binary hit with real `curl`

---

## Milestone 10: Containerizing and Sketching Kubernetes

> ⚠️ **Same honesty pattern as Levels 35 and 36.** `docker version` shows the Docker CLI installed but its daemon unreachable in this environment (`Cannot connect to the Docker daemon...`); `kubectl` is installed but no cluster is reachable, so even `kubectl apply --dry-run=client` fails (client-side dry-run still needs a live API server for its OpenAPI schema). Everything below is written correctly per Docker's and Kubernetes' documented behavior and checked for syntax, but is **illustrative, not verified end-to-end** here - exactly like the corresponding sections of Levels 35 and 36. Run it yourself with a real daemon and cluster.

**Objective:** Package the finished service as a container image and sketch its Kubernetes deployment.

**Instructions:**

1. Create `Dockerfile`:

```bash
cat > Dockerfile << 'EOF'
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
EOF
```

2. Attempt to build it:

```bash
docker build -t urlshortener:1.0.0 .
```

**Actual Output (in this environment):**

```
ERROR: Cannot connect to the Docker daemon at unix:///Users/macbookpro/.docker/run/docker.sock. Is the docker daemon running?
```

> ⚠️ **Illustrative output (unverified in this environment)** - what a successful build is documented to look like:
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

3. Create the Kubernetes manifests:

```bash
mkdir -p k8s
cat > k8s/deployment.yaml << 'EOF'
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
EOF

cat > k8s/service.yaml << 'EOF'
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
EOF
```

4. Attempt to validate them:

```bash
kubectl apply --dry-run=client -f k8s/deployment.yaml
```

**Actual Output (in this environment):**

```
error: error validating "k8s/deployment.yaml": error validating data: failed to download openapi: Get "http://localhost:8080/openapi/v2?timeout=32s": dial tcp [::1]:8080: connect: connection refused; if you choose to ignore these errors, turn validation off with --validate=false
```

This is the identical situation Level 36 documented: `kubectl` client-side dry-run still needs to reach a live API server for schema validation, so even the "safe," non-mutating command fails with no cluster present. The YAML above was hand-checked for well-formed syntax and matches Level 36's `Deployment`/`Service` shape exactly, but has not been applied to a real cluster in this environment.

**Learning Objectives:**
- ✅ Write a multi-stage `Dockerfile` that builds in a full Go image and ships from a minimal one (Level 35)
- ✅ Write a `Deployment` + `ClusterIP` `Service` pair, injecting a secret via `secretKeyRef` rather than baking it into the image (Level 34, Level 36)
- ✅ Recognize the exact same "tool installed, daemon/cluster unreachable" situation this course has already documented twice, and know what real verification would require

---

## Bonus Challenges

### Challenge 1: Click Analytics

Extend `Link` with a slice of click timestamps (or a separate `ClickEvent` type keyed by code) instead of a bare `int` counter, and add a `GET /stats/{code}/timeline` endpoint that buckets clicks by hour.

**Hints:**
- Keep `Clicks int` for backward compatibility and add `ClickEvents []time.Time` alongside it, or introduce a second repository (`ClickRepository`) if the event volume would ever matter for a real system
- The Service layer, not the Repository, should own "bucket by hour" - the Repository just needs to return raw events

### Challenge 2: Link Expiration

Add an optional `ExpiresAt *time.Time` to `Link`, settable at creation time, after which `Resolve` and `Stats` both return a new `ErrLinkExpired` sentinel instead of the link.

**Hints:**
- Check expiration in the Service layer (`Resolve`/`Stats`), not the Repository - "expired" is a business rule, not a storage fact
- Add `linkStatusFor` a new case mapping `ErrLinkExpired` to `410 Gone`, a status code built exactly for this

### Challenge 3: Custom Aliases

Let a user optionally request their own short code (e.g. `{"long_url": "...", "alias": "my-link"}`) instead of always generating a random one, falling back to `ErrCodeGeneration`'s existing collision handling if the alias is already taken.

**Hints:**
- Validate the alias in the Service (allowed characters, length) before ever calling `repo.Create`
- Reuse `repository.ErrCodeTaken` for "alias already in use" - the Repository doesn't need to know an alias is different from a generated code

---

## What You've Learned

After completing all 10 milestones, you can:

✅ Design an API's endpoints, data model, and business rules before writing any code
✅ Build a complete Repository-Service-Handler stack from scratch, layer by layer, with tests at every step
✅ Protect specific endpoints with JWT auth while leaving others public, using request context correctly
✅ Add structured logging and fail-fast, environment-driven configuration as real middleware and startup logic
✅ Implement a token-bucket rate limiter from scratch and wire it in as middleware
✅ Prove an entire system works with a full end-to-end `httptest`-based test, not just isolated unit tests
✅ Recognize and fix two classes of real bugs this level's own authoring hit: `os.LookupEnv`'s "empty vs. unset" distinction, and HTTP keep-alive's dependency on draining response bodies
✅ Follow the established honesty pattern for Docker/Kubernetes verification when no daemon or cluster is available

---

## Next Level

Level 41: Go Interview Preparation
- Turning this capstone's architecture into confident, spoken answers about concurrency, interfaces, error handling, and system design
- Common Go interview questions and how to reason through them out loud
- Practice explaining exactly why this project is shaped the way it is

You just built, tested, and (as far as this environment allows) packaged a complete backend service using nothing but what this course already taught you. Onward to Level 41! 🚀
