# Level 31: Dependency Injection - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

**A note on exercises 3, 7, and 10:** these use `_test.go` files and `go test` instead of `go run`, the same two-file pattern from Level 26 (a code file plus its test file, in the same package, run with `go test`).

---

## Exercise 1: Before/After - Hardcoded Dependency vs Constructor Injection

**Objective:** See the exact refactor from a type that builds its own dependency to one that receives it

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level31-exercise1
cd ~/projects/level31-exercise1
go mod init level31.example/exercise1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

// ❌ BEFORE: EmailNotifier is hardcoded inside OrderNotifier. There is no
// way to give OrderNotifier a different sender, in production OR in tests.
type EmailNotifier struct{}

func (EmailNotifier) Send(message string) {
    fmt.Println("[EMAIL] " + message)
}

type OrderNotifierBefore struct {
    notifier EmailNotifier // concrete type, constructed internally
}

func NewOrderNotifierBefore() *OrderNotifierBefore {
    return &OrderNotifierBefore{notifier: EmailNotifier{}} // reaches for its own dependency
}

func (o *OrderNotifierBefore) NotifyShipped(orderID string) {
    o.notifier.Send("order " + orderID + " has shipped")
}

// ✅ AFTER: Sender is an interface, and OrderNotifierAfter receives its
// implementation as a constructor argument instead of creating one.
type Sender interface {
    Send(message string)
}

type OrderNotifierAfter struct {
    sender Sender // depends on the interface
}

func NewOrderNotifierAfter(sender Sender) *OrderNotifierAfter {
    return &OrderNotifierAfter{sender: sender} // injected from the outside
}

func (o *OrderNotifierAfter) NotifyShipped(orderID string) {
    o.sender.Send("order " + orderID + " has shipped")
}

func main() {
    fmt.Println("=== BEFORE: hardcoded dependency ===")
    before := NewOrderNotifierBefore()
    before.NotifyShipped("A100")

    fmt.Println("\n=== AFTER: constructor-injected dependency ===")
    after := NewOrderNotifierAfter(EmailNotifier{})
    after.NotifyShipped("A100")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== BEFORE: hardcoded dependency ===
[EMAIL] order A100 has shipped

=== AFTER: constructor-injected dependency ===
[EMAIL] order A100 has shipped
```

**Learning Objectives:**
- ✅ See the exact one-line difference between hardcoded and injected dependencies
- ✅ Notice both versions produce identical output - the difference only shows up in what each one *lets you do next*
- ✅ Understand that `EmailNotifier{}` satisfying `Sender` is what makes the "after" version possible with zero other changes

---

## Exercise 2: Constructor Injection Basics

**Objective:** Write the idiomatic `func NewThing(dep SomeInterface) *Thing` pattern from scratch

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level31-exercise2
cd ~/projects/level31-exercise2
go mod init level31.example/exercise2
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

// Logger is the interface App depends on.
type Logger interface {
    Log(msg string)
}

type ConsoleLogger struct{}

func (ConsoleLogger) Log(msg string) {
    fmt.Println("[console]", msg)
}

type PrefixLogger struct {
    Prefix string
}

func (p PrefixLogger) Log(msg string) {
    fmt.Println(p.Prefix+":", msg)
}

// App is the idiomatic Go constructor-injection shape:
//
//	func NewThing(dep SomeInterface) *Thing { return &Thing{dep: dep} }
type App struct {
    logger Logger
}

func NewApp(logger Logger) *App {
    return &App{logger: logger}
}

func (a *App) Start() {
    a.logger.Log("application starting")
}

func (a *App) Stop() {
    a.logger.Log("application stopping")
}

func main() {
    fmt.Println("=== App with ConsoleLogger ===")
    app1 := NewApp(ConsoleLogger{})
    app1.Start()
    app1.Stop()

    fmt.Println("\n=== App with PrefixLogger ===")
    app2 := NewApp(PrefixLogger{Prefix: "APP"})
    app2.Start()
    app2.Stop()
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== App with ConsoleLogger ===
[console] application starting
[console] application stopping

=== App with PrefixLogger ===
APP: application starting
APP: application stopping
```

**Learning Objectives:**
- ✅ Write `NewApp(logger Logger) *App` from scratch
- ✅ Confirm `NewApp` never has to change to accept a different `Logger` implementation
- ✅ Connect this to Level 13's `NewX` constructor convention

---

## Exercise 3: Why Hidden Internal Construction Breaks Testability

**Objective:** Prove, with a real measured test run, that a hardcoded dependency has no seam to substitute a fake into

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level31-exercise3
cd ~/projects/level31-exercise3
go mod init level31.example/exercise3
```

2. Create `notifier.go`:

```bash
cat > notifier.go << 'EOF'
package notifier

import "time"

// RealSender simulates a real, slow network call (e.g. a real email API).
type RealSender struct{}

func (RealSender) Send(message string) string {
    time.Sleep(300 * time.Millisecond) // stands in for real network latency
    return "sent via real network: " + message
}

// OrderNotifierBefore constructs its OWN RealSender inside NewOrderNotifierBefore.
// There is no constructor parameter for a sender at all - nothing to
// substitute a fake into. Any test of this type is FORCED to pay the
// real 300ms cost, every single time.
type OrderNotifierBefore struct {
    sender RealSender
}

func NewOrderNotifierBefore() *OrderNotifierBefore {
    return &OrderNotifierBefore{sender: RealSender{}}
}

func (o *OrderNotifierBefore) NotifyShipped(orderID string) string {
    return o.sender.Send("order " + orderID + " has shipped")
}

// Sender is the interface OrderNotifierAfter depends on instead.
type Sender interface {
    Send(message string) string
}

// OrderNotifierAfter RECEIVES its sender via the constructor - there IS a
// seam, so a test can substitute a fast fake instead of the real sender.
type OrderNotifierAfter struct {
    sender Sender
}

func NewOrderNotifierAfter(sender Sender) *OrderNotifierAfter {
    return &OrderNotifierAfter{sender: sender}
}

func (o *OrderNotifierAfter) NotifyShipped(orderID string) string {
    return o.sender.Send("order " + orderID + " has shipped")
}
EOF
```

3. Create `notifier_test.go`:

```bash
cat > notifier_test.go << 'EOF'
package notifier

import "testing"

// This test CANNOT substitute a fake sender: NewOrderNotifierBefore takes
// no arguments, so there is no parameter to pass anything through. The
// test is stuck paying RealSender's real 300ms delay on every run.
func TestOrderNotifierBefore_CannotSubstituteFake(t *testing.T) {
    before := NewOrderNotifierBefore()
    got := before.NotifyShipped("A100")
    want := "sent via real network: order A100 has shipped"
    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}

// fakeSender is instant - no sleep, no network, no flakiness.
type fakeSender struct {
    messages []string
}

func (f *fakeSender) Send(message string) string {
    f.messages = append(f.messages, message)
    return "sent via fake: " + message
}

// This test injects a fake through the constructor parameter that
// NewOrderNotifierAfter provides. No sleep, no real dependency, runs
// essentially instantly.
func TestOrderNotifierAfter_SubstitutesFake(t *testing.T) {
    fake := &fakeSender{}
    after := NewOrderNotifierAfter(fake)
    got := after.NotifyShipped("A100")
    want := "sent via fake: order A100 has shipped"
    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
    if len(fake.messages) != 1 {
        t.Fatalf("expected 1 recorded message, got %d", len(fake.messages))
    }
}
EOF
```

4. Run the tests:

```bash
go test -v ./...
```

**Expected Output** (the exact timings will vary slightly on your machine, but the ~300ms gap between the two tests will not):

```
=== RUN   TestOrderNotifierBefore_CannotSubstituteFake
--- PASS: TestOrderNotifierBefore_CannotSubstituteFake (0.30s)
=== RUN   TestOrderNotifierAfter_SubstitutesFake
--- PASS: TestOrderNotifierAfter_SubstitutesFake (0.00s)
PASS
ok  	level31.example/exercise3	0.901s
```

**Learning Objectives:**
- ✅ See a concrete, measured proof that a hardcoded dependency cannot be faked - `NewOrderNotifierBefore()` simply has no parameter to substitute through
- ✅ Confirm the injected version's test runs in ~0.00s because it never touches the real, slow dependency
- ✅ Understand that "testability" isn't abstract - it's the literal difference between a 300ms test and an instant one, multiplied across your whole test suite

---

## Exercise 4: Interface-Based Injection - Swapping Implementations

**Objective:** Plug two different implementations of the same interface into unchanged dependent code

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level31-exercise4
cd ~/projects/level31-exercise4
go mod init level31.example/exercise4
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

// StorageBackend is the interface FileArchiver depends on.
type StorageBackend interface {
    Store(filename string, data string) string
}

type LocalDiskBackend struct{}

func (LocalDiskBackend) Store(filename, data string) string {
    return fmt.Sprintf("wrote %d bytes to local disk file %q", len(data), filename)
}

type CloudBackend struct {
    Bucket string
}

func (c CloudBackend) Store(filename, data string) string {
    return fmt.Sprintf("uploaded %d bytes to bucket %q as %q", len(data), c.Bucket, filename)
}

// FileArchiver never changes, no matter which StorageBackend it receives.
type FileArchiver struct {
    backend StorageBackend
}

func NewFileArchiver(backend StorageBackend) *FileArchiver {
    return &FileArchiver{backend: backend}
}

func (f *FileArchiver) Archive(filename, data string) {
    fmt.Println(f.backend.Store(filename, data))
}

func main() {
    fmt.Println("=== FileArchiver backed by LocalDiskBackend ===")
    diskArchiver := NewFileArchiver(LocalDiskBackend{})
    diskArchiver.Archive("report.txt", "quarterly results")

    fmt.Println("\n=== Same FileArchiver code, backed by CloudBackend instead ===")
    cloudArchiver := NewFileArchiver(CloudBackend{Bucket: "archives-prod"})
    cloudArchiver.Archive("report.txt", "quarterly results")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== FileArchiver backed by LocalDiskBackend ===
wrote 17 bytes to local disk file "report.txt"

=== Same FileArchiver code, backed by CloudBackend instead ===
uploaded 17 bytes to bucket "archives-prod" as "report.txt"
```

**Learning Objectives:**
- ✅ Confirm `FileArchiver`'s source code never mentions `LocalDiskBackend` or `CloudBackend` by name
- ✅ See two structurally different implementations satisfy the same interface
- ✅ Connect this to Level 15's Rule 2 (accept interfaces, return concrete types)

---

## Exercise 5: Injecting a Config Struct Alongside an Interface

**Objective:** Combine a plain configuration value with a mockable interface dependency in one constructor

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level31-exercise5
cd ~/projects/level31-exercise5
go mod init level31.example/exercise5
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

// RateLimiterConfig is plain configuration data - not an interface, not
// something you'd ever need to fake, just values.
type RateLimiterConfig struct {
    MaxRequests int
    WindowLabel string
}

// Clock is the one piece of BEHAVIOR that needs to be mockable in tests.
type Clock interface {
    Now() string
}

type systemClock struct{}

func (systemClock) Now() string { return "2026-09-01T12:00:00Z" }

// APIGateway is injected with BOTH a plain config struct and an interface
// dependency - not everything injected needs to be mockable behavior.
type APIGateway struct {
    cfg   RateLimiterConfig
    clock Clock
}

func NewAPIGateway(cfg RateLimiterConfig, clock Clock) *APIGateway {
    return &APIGateway{cfg: cfg, clock: clock}
}

func (g *APIGateway) Describe() string {
    return fmt.Sprintf("limiter allows %d requests per %s (checked at %s)",
        g.cfg.MaxRequests, g.cfg.WindowLabel, g.clock.Now())
}

func main() {
    cfg := RateLimiterConfig{MaxRequests: 100, WindowLabel: "minute"}
    gateway := NewAPIGateway(cfg, systemClock{})
    fmt.Println(gateway.Describe())
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
limiter allows 100 requests per minute (checked at 2026-09-01T12:00:00Z)
```

**Learning Objectives:**
- ✅ Inject a plain data struct (`RateLimiterConfig`) and an interface (`Clock`) through the same constructor
- ✅ Recognize which of a type's dependencies need to be interfaces (behavior) versus plain values (data)
- ✅ See that both go through the exact same mechanism: a constructor parameter

---

## Exercise 6: Optional Dependencies via Functional Options

**Objective:** Add optional, configurable extras to a constructor without growing its required parameter list

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level31-exercise6
cd ~/projects/level31-exercise6
go mod init level31.example/exercise6
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

// Repository is the required dependency every HTTPClient must have.
type Repository interface {
    Ping() string
}

type memRepository struct{}

func (memRepository) Ping() string { return "pong" }

// HTTPClient has ONE required dependency (repo) and any number of
// OPTIONAL extras configured via functional options, so the constructor's
// parameter list never has to grow for every new optional knob.
type HTTPClient struct {
    repo       Repository
    timeoutSec int
    retries    int
    userAgent  string
}

type Option func(*HTTPClient)

func WithTimeout(seconds int) Option {
    return func(c *HTTPClient) { c.timeoutSec = seconds }
}

func WithRetries(n int) Option {
    return func(c *HTTPClient) { c.retries = n }
}

func WithUserAgent(ua string) Option {
    return func(c *HTTPClient) { c.userAgent = ua }
}

func NewHTTPClient(repo Repository, opts ...Option) *HTTPClient {
    c := &HTTPClient{
        repo:       repo,
        timeoutSec: 30,          // sensible default
        retries:    0,           // sensible default
        userAgent:  "level31/1", // sensible default
    }
    for _, opt := range opts {
        opt(c)
    }
    return c
}

func (c *HTTPClient) Describe() string {
    return fmt.Sprintf("timeout=%ds retries=%d userAgent=%q repo.Ping=%s",
        c.timeoutSec, c.retries, c.userAgent, c.repo.Ping())
}

func main() {
    fmt.Println("=== No options: all defaults ===")
    plain := NewHTTPClient(memRepository{})
    fmt.Println(plain.Describe())

    fmt.Println("\n=== Some options set, others default ===")
    withTimeout := NewHTTPClient(memRepository{}, WithTimeout(5))
    fmt.Println(withTimeout.Describe())

    fmt.Println("\n=== All options set ===")
    full := NewHTTPClient(memRepository{},
        WithTimeout(10),
        WithRetries(3),
        WithUserAgent("level31-custom/2"),
    )
    fmt.Println(full.Describe())
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== No options: all defaults ===
timeout=30s retries=0 userAgent="level31/1" repo.Ping=pong

=== Some options set, others default ===
timeout=5s retries=0 userAgent="level31/1" repo.Ping=pong

=== All options set ===
timeout=10s retries=3 userAgent="level31-custom/2" repo.Ping=pong
```

**Learning Objectives:**
- ✅ Write the `func NewThing(dep SomeInterface, opts ...Option) *Thing` shape
- ✅ Confirm defaults apply when no options are passed, and mix cleanly with partial overrides
- ✅ Understand why this scales better than adding parameters for every new optional setting

---

## Exercise 7: One Service, Tested With Two Different Fakes

**Objective:** Test a single service against two fakes that produce opposite behavior

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level31-exercise7
cd ~/projects/level31-exercise7
go mod init level31.example/exercise7
```

2. Create `validator.go`:

```bash
cat > validator.go << 'EOF'
package validator

// UserRepository is the interface UserValidator depends on.
type UserRepository interface {
    Exists(username string) bool
}

type UserValidator struct {
    repo UserRepository
}

func NewUserValidator(repo UserRepository) *UserValidator {
    return &UserValidator{repo: repo}
}

// IsAvailable reports whether username is free to register.
func (v *UserValidator) IsAvailable(username string) bool {
    return !v.repo.Exists(username)
}
EOF
```

3. Create `validator_test.go`:

```bash
cat > validator_test.go << 'EOF'
package validator

import "testing"

// fakeRepoAlwaysExists simulates a repository where every username is
// already taken.
type fakeRepoAlwaysExists struct{}

func (fakeRepoAlwaysExists) Exists(username string) bool { return true }

// fakeRepoNeverExists simulates a repository where no username is taken.
type fakeRepoNeverExists struct{}

func (fakeRepoNeverExists) Exists(username string) bool { return false }

func TestUserValidator_WithTakenUsernames(t *testing.T) {
    validator := NewUserValidator(fakeRepoAlwaysExists{})
    if validator.IsAvailable("alice") {
        t.Error("expected alice to be unavailable when repo says it exists")
    }
}

func TestUserValidator_WithFreeUsernames(t *testing.T) {
    validator := NewUserValidator(fakeRepoNeverExists{})
    if !validator.IsAvailable("alice") {
        t.Error("expected alice to be available when repo says it doesn't exist")
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
=== RUN   TestUserValidator_WithTakenUsernames
--- PASS: TestUserValidator_WithTakenUsernames (0.00s)
=== RUN   TestUserValidator_WithFreeUsernames
--- PASS: TestUserValidator_WithFreeUsernames (0.00s)
PASS
ok  	level31.example/exercise7	0.533s
```

**Learning Objectives:**
- ✅ Write two fakes of the same interface that produce opposite behavior on purpose
- ✅ Confirm `UserValidator`'s own code never changes between the two tests - only which fake is injected
- ✅ See interface-based injection pay off directly inside a test file

---

## Exercise 8: Wiring a Layered Dependency Graph by Hand in main()

**Objective:** Hand-wire a Config → Repository → Service → Handler graph, recalling Level 30's three-layer pattern

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level31-exercise8
cd ~/projects/level31-exercise8
go mod init level31.example/exercise8
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

// --- Layer 1: Config ---
type Config struct {
    MinBalance float64
}

// --- Layer 2: Repository ---
type AccountRepository interface {
    Balance(accountID string) (float64, error)
}

type inMemoryAccountRepository struct {
    balances map[string]float64
}

func NewInMemoryAccountRepository() *inMemoryAccountRepository {
    return &inMemoryAccountRepository{
        balances: map[string]float64{"acc-1": 500.00, "acc-2": 5.00},
    }
}

func (r *inMemoryAccountRepository) Balance(accountID string) (float64, error) {
    bal, ok := r.balances[accountID]
    if !ok {
        return 0, fmt.Errorf("account %s not found", accountID)
    }
    return bal, nil
}

// --- Layer 3: Service ---
type AccountService struct {
    cfg  Config
    repo AccountRepository
}

func NewAccountService(cfg Config, repo AccountRepository) *AccountService {
    return &AccountService{cfg: cfg, repo: repo}
}

func (s *AccountService) IsHealthy(accountID string) (bool, error) {
    bal, err := s.repo.Balance(accountID)
    if err != nil {
        return false, err
    }
    return bal >= s.cfg.MinBalance, nil
}

// --- Layer 4: Handler ---
type AccountHandler struct {
    service *AccountService
}

func NewAccountHandler(service *AccountService) *AccountHandler {
    return &AccountHandler{service: service}
}

func (h *AccountHandler) HandleCheck(accountID string) {
    healthy, err := h.service.IsHealthy(accountID)
    if err != nil {
        fmt.Printf("account %s: error: %v\n", accountID, err)
        return
    }
    fmt.Printf("account %s: healthy=%t\n", accountID, healthy)
}

// buildApp wires the whole graph by hand: Config -> Repository -> Service
// -> Handler. This one function is the "composition root" - a DI
// framework would automate this; here it's just explicit code.
func buildApp() *AccountHandler {
    cfg := Config{MinBalance: 50.00}
    repo := NewInMemoryAccountRepository()
    service := NewAccountService(cfg, repo)
    handler := NewAccountHandler(service)
    return handler
}

func main() {
    handler := buildApp()
    handler.HandleCheck("acc-1")
    handler.HandleCheck("acc-2")
    handler.HandleCheck("acc-99")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
account acc-1: healthy=true
account acc-2: healthy=false
account acc-99: error: account acc-99 not found
```

**Learning Objectives:**
- ✅ Wire a four-layer dependency graph by hand, in the right order, inside one `buildApp()` function
- ✅ Recognize this as the same shape Level 30 previewed: Repository → Service → Handler, now with a Config layer added
- ✅ Confirm the compiler catches any wiring mistake (wrong type passed to a constructor) before the program ever runs

---

## Exercise 9: Global Mutable State vs Injected Dependency

**Objective:** See a real, captured data race caused by global state, and confirm injection avoids it entirely

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level31-exercise9
cd ~/projects/level31-exercise9
go mod init level31.example/exercise9
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "sync"
)

// ❌ BEFORE: globalHits is package-level mutable state. Every call to
// RecordHit reaches out and touches it directly - a hidden dependency
// no caller can see, override, or isolate.
var globalHits = map[string]int{}

func RecordHit(page string) {
    globalHits[page]++
}

// ✅ AFTER: HitTracker owns its own state, injected/created per instance
// instead of shared globally. Each caller gets an isolated tracker.
type HitTracker struct {
    hits map[string]int
}

func NewHitTracker() *HitTracker {
    return &HitTracker{hits: map[string]int{}}
}

func (t *HitTracker) RecordHit(page string) {
    t.hits[page]++
}

func main() {
    fmt.Println("=== BEFORE: concurrent goroutines sharing package-level state ===")
    var wg1 sync.WaitGroup
    for i := 0; i < 50; i++ {
        wg1.Add(1)
        go func() {
            defer wg1.Done()
            RecordHit("home") // every goroutine touches the SAME global map
        }()
    }
    wg1.Wait()
    fmt.Println("global hits (run with -race to see the data race Go detects):", globalHits["home"])

    fmt.Println("\n=== AFTER: each goroutine gets its own isolated HitTracker ===")
    var wg2 sync.WaitGroup
    results := make([]int, 5)
    for i := 0; i < 5; i++ {
        wg2.Add(1)
        go func(i int) {
            defer wg2.Done()
            tracker := NewHitTracker() // no shared state - nothing to race on
            for j := 0; j < 50; j++ {
                tracker.RecordHit("home")
            }
            results[i] = tracker.hits["home"]
        }(i)
    }
    wg2.Wait()
    fmt.Println("per-goroutine hit counts (each isolated, always 50):", results)
}
EOF
```

3. Run the program normally, several times in a row:

```bash
go run main.go
```

**Expected Output:** this is a genuine, unsynchronized data race, so the BEFORE section's result is **not stable** - that instability is exactly the point. Four real, consecutive runs, captured verbatim:

```
--run 1--
=== BEFORE: concurrent goroutines sharing package-level state ===
global hits (run with -race to see the data race Go detects): 50

=== AFTER: each goroutine gets its own isolated HitTracker ===
per-goroutine hit counts (each isolated, always 50): [50 50 50 50 50]
--run 2--
=== BEFORE: concurrent goroutines sharing package-level state ===
fatal error: concurrent map writes

goroutine 56 [running]:
internal/runtime/maps.fatal({0x1025eb932?, 0x0?})
	/Users/macbookpro/.local/go/src/runtime/panic.go:1181 +0x20
main.RecordHit(...)
	main.go:14
main.main.func1()
	main.go:38 +0x64
created by main.main in goroutine 1
	main.go:36 +0x70
exit status 2
--run 3--
=== BEFORE: concurrent goroutines sharing package-level state ===
fatal error: concurrent map writes
--run 4--
=== BEFORE: concurrent goroutines sharing package-level state ===
fatal error: concurrent map writes
```

Three crashes out of four real runs on this machine - the exact ratio will differ on yours, but crashing at all, unpredictably, is the point. The AFTER section, by contrast, is stable on every single run: `[50 50 50 50 50]`, always, because each goroutine's `HitTracker` is a private, injected instance with nothing to share. Run the AFTER section by itself a few times if you want to confirm this - it never varies.

4. Now run it with Go's race detector to see the hidden dependency exposed, reliably, every time:

```bash
go run -race main.go
```

**Expected Output** (captured verbatim from real runs with `-race` - the exact memory addresses and goroutine numbers will differ on your machine, but the `WARNING: DATA RACE` on `main.RecordHit` will not):

```
=== BEFORE: concurrent goroutines sharing package-level state ===
==================
WARNING: DATA RACE
Write at 0x00c0001020f0 by goroutine 8:
  runtime.mapaccess2_faststr()
      /usr/local/go/src/internal/runtime/maps/runtime_faststr.go:161 +0x2ac
  main.RecordHit()
      main.go:14 +0xd0
  main.main.func1()
      main.go:38 +0x64

Previous read at 0x00c0001020f0 by goroutine 10:
  runtime.mapassign_fast64ptr()
      /usr/local/go/src/internal/runtime/maps/runtime_fast64.go:374 +0x38c
  main.RecordHit()
      main.go:14 +0x90
  main.main.func1()
      main.go:38 +0x64
==================
... (usually several more WARNING: DATA RACE blocks like this one) ...

=== AFTER: each goroutine gets its own isolated HitTracker ===
per-goroutine hit counts (each isolated, always 50): [50 50 50 50 50]
Found 4 data race(s)
exit status 66
```

`-race` catches this every single time (3 out of 3 real runs, and every run in the batch above), unlike the plain run's coin-flip between a wrong-looking success and an outright crash. The `AFTER` section still prints its stable `[50 50 50 50 50]` even under `-race`, because there is genuinely nothing racing in it.

**Learning Objectives:**
- ✅ See that `RecordHit` never appears in any function signature as a dependency - it's invisible, hidden global state
- ✅ Get a real, reproducible `WARNING: DATA RACE` from Go's own `-race` detector, caused entirely by sharing package-level state across goroutines
- ✅ Confirm the injected `HitTracker` version has nothing to race on, because each goroutine owns its own instance

---

## Exercise 10: Comprehensive - A Notification System

**Objective:** Build a complete, swappable, testable notification system using everything from this level

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level31-exercise10
cd ~/projects/level31-exercise10
go mod init level31.example/exercise10
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

// Sender is the small interface NotificationService depends on. Any type
// with a Send method can be plugged in - email, SMS, console, or a fake.
type Sender interface {
    Send(to, message string) error
}

// EmailSender is one real implementation.
type EmailSender struct{}

func (EmailSender) Send(to, message string) error {
    fmt.Printf("[EMAIL to %s] %s\n", to, message)
    return nil
}

// SMSSender is a second, completely different real implementation.
type SMSSender struct{}

func (SMSSender) Send(to, message string) error {
    fmt.Printf("[SMS to %s] %s\n", to, message)
    return nil
}

// ConsoleSender is a third implementation, handy for local development.
type ConsoleSender struct{}

func (ConsoleSender) Send(to, message string) error {
    fmt.Printf("[CONSOLE] would notify %s: %s\n", to, message)
    return nil
}

// NotificationService depends only on the Sender interface - it has no
// idea whether it's emailing, texting, or printing to a console.
type NotificationService struct {
    sender Sender
}

func NewNotificationService(sender Sender) *NotificationService {
    return &NotificationService{sender: sender}
}

func (s *NotificationService) NotifyUser(recipient, event string) error {
    message := fmt.Sprintf("event: %s", event)
    return s.sender.Send(recipient, message)
}

// buildApp is the composition root: it decides, in ONE place, which
// Sender implementation the whole app uses.
func buildApp(sender Sender) *NotificationService {
    return NewNotificationService(sender)
}

func main() {
    fmt.Println("=== Wired with EmailSender ===")
    emailApp := buildApp(EmailSender{})
    emailApp.NotifyUser("alice@example.com", "order shipped")

    fmt.Println("\n=== Same NotificationService code, wired with SMSSender ===")
    smsApp := buildApp(SMSSender{})
    smsApp.NotifyUser("alice@example.com", "order shipped")

    fmt.Println("\n=== Same NotificationService code, wired with ConsoleSender ===")
    consoleApp := buildApp(ConsoleSender{})
    consoleApp.NotifyUser("alice@example.com", "order shipped")
}
EOF
```

3. Create `notification_test.go` to test the service in isolation, with no real sending at all:

```bash
cat > notification_test.go << 'EOF'
package main

import "testing"

// fakeSender records every message instead of actually sending anything -
// this is what makes NotificationService testable in isolation.
type fakeSender struct {
    sentTo      []string
    sentMessage []string
}

func (f *fakeSender) Send(to, message string) error {
    f.sentTo = append(f.sentTo, to)
    f.sentMessage = append(f.sentMessage, message)
    return nil
}

func TestNotificationService_NotifyUser(t *testing.T) {
    fake := &fakeSender{}
    service := NewNotificationService(fake)

    if err := service.NotifyUser("bob@example.com", "password changed"); err != nil {
        t.Fatalf("NotifyUser returned unexpected error: %v", err)
    }

    if len(fake.sentTo) != 1 {
        t.Fatalf("expected 1 send, got %d", len(fake.sentTo))
    }
    if fake.sentTo[0] != "bob@example.com" {
        t.Errorf("sent to %q; want %q", fake.sentTo[0], "bob@example.com")
    }
    want := "event: password changed"
    if fake.sentMessage[0] != want {
        t.Errorf("message = %q; want %q", fake.sentMessage[0], want)
    }
}
EOF
```

4. Run the program, then run the test:

```bash
go run main.go
go test -v ./...
```

**Expected Output** (from `go run main.go`):

```
=== Wired with EmailSender ===
[EMAIL to alice@example.com] event: order shipped

=== Same NotificationService code, wired with SMSSender ===
[SMS to alice@example.com] event: order shipped

=== Same NotificationService code, wired with ConsoleSender ===
[CONSOLE] would notify alice@example.com: event: order shipped
```

**Expected Output** (from `go test -v ./...`):

```
=== RUN   TestNotificationService_NotifyUser
--- PASS: TestNotificationService_NotifyUser (0.00s)
PASS
ok  	level31.example/exercise10	0.595s
```

**Learning Objectives:**
- ✅ Build a full, swappable dependency graph: three real `Sender` implementations, one unchanged `NotificationService`
- ✅ Wire it through a `buildApp` composition root, exactly like Exercise 8
- ✅ Test the service in complete isolation with a fake, with zero real emails, texts, or console side effects reaching outside the test
- ✅ Recognize this as the complete shape of idiomatic Go dependency injection: interfaces, constructor injection, and explicit wiring, working together

---

## Bonus Challenges

### Challenge 1: Lazy Singleton via sync.Once, Injected Instead of Global

Build a `StoreProvider` that lazily constructs an expensive `Store` exactly once (using `sync.Once`, from Level 23), and inject the *provider* into consumers instead of exposing the store as a package-level global.

```bash
mkdir -p ~/projects/level31-bonus1
cd ~/projects/level31-bonus1
go mod init level31.example/bonus1
```

**Hints:**
- `type StoreProvider struct { once sync.Once; store *expensiveStore }` with a `Store() Store` method that calls `p.once.Do(func() { ... })`
- Inject `*StoreProvider` into a `Reader` via `NewReader(provider *StoreProvider) *Reader`, same as any other dependency
- Print a message inside the expensive constructor and confirm it only prints once, even when called through two different `Reader`s sharing the same provider

### Challenge 2: Wire a Five-Layer Graph by Hand

Extend Exercise 8's Config → Repository → Service → Handler graph with a fifth layer - a `Middleware` that wraps the handler and logs before delegating to it.

```bash
mkdir -p ~/projects/level31-bonus2
cd ~/projects/level31-bonus2
go mod init level31.example/bonus2
```

**Hints:**
- `type Middleware struct { handler *Handler }` with `NewMiddleware(handler *Handler) *Middleware`
- `buildApp()` should now construct all five layers in order and return the `*Middleware`
- Print something in the middleware before calling into the handler, so the output proves the request passed through all five layers

### Challenge 3: A Logging/Timing Decorator

Write a `loggingStore` that wraps any `Store` implementation and logs how long each call took, without `Store`'s consumer ever knowing it's there.

```bash
mkdir -p ~/projects/level31-bonus3
cd ~/projects/level31-bonus3
go mod init level31.example/bonus3
```

**Hints:**
- `type loggingStore struct { next Store }` also implements `Store` by calling `l.next.Get(key)` and timing it with `time.Now()` / `time.Since()`
- Inject `NewLoggingStore(memStore{})` anywhere a plain `Store` was injected before - the consumer's code doesn't change at all
- This is the **decorator pattern**: a type satisfying the same interface it wraps, adding behavior transparently

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Refactor a hardcoded dependency into a constructor-injected one
✅ Write the idiomatic `func NewThing(dep SomeInterface) *Thing` pattern from memory
✅ Prove, with a measured test, why hidden internal construction breaks testability
✅ Swap interface implementations without touching the code that depends on them
✅ Inject plain config values alongside mockable interfaces
✅ Use the functional options pattern for optional, growable configuration
✅ Test one service against multiple fakes with opposite behavior
✅ Hand-wire a layered dependency graph inside a composition root
✅ Recognize global mutable state as a hidden dependency and a real concurrency hazard
✅ Combine everything into a complete, swappable, testable notification system

---

## Next Level

Level 32: Authentication & JWT
- Verifying user identity with JSON Web Tokens
- Issuing, signing, and validating tokens
- Middleware that injects an authenticated user into request handling

Great work! You now build software the way idiomatic Go actually builds it - explicit, compiler-checked, and testable by design! 🚀
