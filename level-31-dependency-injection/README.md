# Level 31: Dependency Injection - Complete Guide

## Introduction

Welcome to Level 31! Level 30 (Repository-Service-Handler) gave you a first taste of this level's whole subject: you wired a repository into a service, and a service into a handler, all inside `main()`. That wiring - passing dependencies in from the outside instead of letting each layer build its own - already *was* dependency injection. You just hadn't named it yet.

This level names it, and goes deep. Dependency Injection (DI) has a reputation in other languages for meaning heavyweight frameworks, XML wiring files, and `@Autowired` annotations that feel like magic. In Go, it means almost none of that. Idiomatic Go DI is just: **a type receives its dependencies as constructor arguments instead of constructing or reaching for them itself.** That's it. No framework, no container library, no reflection-based magic - just function parameters and interfaces, both of which you already know cold.

By the end of this level you'll understand why Go deliberately has no dominant DI framework (unlike Java's Spring or C#'s .NET), how to hand-wire an entire application's dependency graph in a single, readable `main()` function, and how to use interfaces (Level 15), constructors (Level 13), and tests with fakes (Level 26) together to build software that's flexible and provably testable.

---

## Table of Contents

1. [What Dependency Injection Actually Means](#what-dependency-injection-actually-means)
2. [Why It Matters: Testability, Flexibility, Decoupling](#why-it-matters-testability-flexibility-decoupling)
3. [Constructor Injection: The Idiomatic Go Pattern](#constructor-injection-the-idiomatic-go-pattern)
4. [Why Go Usually Doesn't Need a DI Framework](#why-go-usually-doesnt-need-a-di-framework)
5. [Interface-Based Injection: Swapping Implementations](#interface-based-injection-swapping-implementations)
6. [Injecting Configuration Alongside Interfaces](#injecting-configuration-alongside-interfaces)
7. [Optional Dependencies: The Functional Options Pattern](#optional-dependencies-the-functional-options-pattern)
8. [A Hand-Rolled Container: Wiring a Dependency Graph in main()](#a-hand-rolled-container-wiring-a-dependency-graph-in-main)
9. [When a DI Framework Might Actually Help](#when-a-di-framework-might-actually-help)
10. [Best Practices](#best-practices)
11. [Common Mistakes](#common-mistakes)

---

## What Dependency Injection Actually Means

Strip away the framework mystique and DI is a one-sentence idea:

> **A type should receive the things it depends on from the outside, instead of constructing them or reaching for them itself.**

"The things it depends on" are its **dependencies** - other types, values, or services it needs to do its job. "From the outside" usually means: as arguments to its constructor.

### Before: A Type That Builds Its Own Dependency

```go
package main

import "fmt"

// ❌ BEFORE - ReportGenerator constructs its own logger internally.
// It is tightly coupled to *ConsoleLogger and can never use anything else.
type ConsoleLogger struct{}

func (c *ConsoleLogger) Log(msg string) {
    fmt.Println("[LOG]", msg)
}

type ReportGenerator struct {
    logger *ConsoleLogger // hardcoded concrete type, built inside the struct
}

func NewReportGenerator() *ReportGenerator {
    return &ReportGenerator{
        logger: &ConsoleLogger{}, // reaches for its own dependency
    }
}

func (r *ReportGenerator) Generate(name string) {
    r.logger.Log("generating report: " + name)
}

func main() {
    fmt.Println("=== BEFORE: hardcoded dependency ===")
    before := NewReportGenerator()
    before.Generate("Q3 Sales")
}
```

Real output:

```
=== BEFORE: hardcoded dependency ===
[LOG] generating report: Q3 Sales
```

It works. But `ReportGenerator` is welded to `*ConsoleLogger` forever - there's no way to log somewhere else, and no way to test `Generate` without also running real `ConsoleLogger` code every single time.

### After: A Type That Receives Its Dependency

```go
package main

import "fmt"

// ✅ AFTER - Logger is an interface, and ReportGenerator RECEIVES its
// implementation from the outside via the constructor. ReportGenerator
// never constructs a logger itself.
type Logger interface {
    Log(msg string)
}

type ConsoleLogger struct{}

func (c *ConsoleLogger) Log(msg string) {
    fmt.Println("[LOG]", msg)
}

type ReportGenerator struct {
    logger Logger // depends on the interface, not a concrete type
}

func NewReportGenerator(logger Logger) *ReportGenerator {
    return &ReportGenerator{logger: logger} // injected, not created
}

func (r *ReportGenerator) Generate(name string) {
    r.logger.Log("generating report: " + name)
}

// silentLogger is a stand-in that could just as easily be a test fake.
type silentLogger struct {
    messages []string
}

func (s *silentLogger) Log(msg string) {
    s.messages = append(s.messages, msg)
}

func main() {
    fmt.Println("=== AFTER: constructor-injected dependency ===")
    real := NewReportGenerator(&ConsoleLogger{})
    real.Generate("Q3 Sales")

    fmt.Println("\n=== AFTER: swapped for a silent fake, zero changes to ReportGenerator ===")
    fake := &silentLogger{}
    withFake := NewReportGenerator(fake)
    withFake.Generate("Q4 Sales")
    fmt.Println("captured message:", fake.messages[0])
}
```

Real output:

```
=== AFTER: constructor-injected dependency ===
[LOG] generating report: Q3 Sales

=== AFTER: swapped for a silent fake, zero changes to ReportGenerator ===
captured message: generating report: Q4 Sales
```

Nothing about `ReportGenerator`'s own code changed between the two calls to `NewReportGenerator` in the second example - only what was passed in. That's the entire idea. Everything else in this level is elaboration on this one shift in who's responsible for construction.

---

## Why It Matters: Testability, Flexibility, Decoupling

Three concrete payoffs fall out of "receive, don't construct":

### 1. Testability

When a dependency is injected, a test can substitute a **fake** (Level 26) for the real thing - no real database, no real network call, no real email sent. The "before" `ReportGenerator` above can only ever use a real `*ConsoleLogger`; the "after" version can be tested with `silentLogger` in complete isolation. Section 3 of this level makes this concrete with a directly measured before/after.

### 2. Flexibility

Swap implementations without touching the dependent code at all. A `PaymentProcessor` that depends on a `PaymentGateway` interface doesn't care whether it's charging through Stripe or PayPal - both work through the exact same `Pay` method with zero changes to `PaymentProcessor` itself (Section 5 shows this side by side).

### 3. Decoupling

The dependent code only needs to know the dependency's **interface** - the small contract from Level 15 - not its implementation details, its own dependencies, or how it was built. `ReportGenerator` knows `Logger` has a `Log(msg string)` method; it knows nothing else about `ConsoleLogger`, and it shouldn't.

These three benefits are really one benefit seen from different angles: giving up control over *how* a dependency is built buys you the freedom to change, fake, or replace it later without touching the code that uses it.

---

## Constructor Injection: The Idiomatic Go Pattern

Nearly all of dependency injection in idiomatic Go is one recurring shape:

```go
func NewThing(dep SomeInterface) *Thing {
    return &Thing{dep: dep}
}
```

This is **constructor injection** - the dependency arrives as a parameter to the `NewX` constructor (Level 13's convention) and gets stored on the struct for later use. Compare it to the two alternatives DI theory usually lists:

| Injection style | How it looks | Used in idiomatic Go? |
|---|---|---|
| **Constructor injection** | `NewThing(dep)` stores it on the struct | ✅ Yes - the default, nearly universal choice |
| **Setter injection** | `thing.SetDep(dep)` after construction | ⚠️ Rare - leaves a window where `thing` exists with no dependency set |
| **Field injection** | Reflection sets private fields directly (common in Java/Spring) | ❌ No standard mechanism, and would fight Go's lack of reflection-heavy conventions |

Constructor injection wins in Go because it lines up with everything you already know:

- It reuses the `NewX(...) *X` convention from Level 13 - no new syntax to learn.
- The struct is **never** in a half-built state: by the time `NewThing` returns, every dependency is already set.
- The compiler enforces it. Forget to pass a required dependency, and `NewThing(dep)` simply won't compile - there's no equivalent of a `NullPointerException` from a dependency nobody set.

Here it is applied to a small `App` type, with two interchangeable `Logger` implementations proving the constructor doesn't care which one it gets:

```go
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
```

Running it with each logger produces exactly what you'd expect - `NewApp` never had to change:

```
=== App with ConsoleLogger ===
[console] application starting
[console] application stopping

=== App with PrefixLogger ===
APP: application starting
APP: application stopping
```

There's rarely a need for anything fancier than this pattern. Sections 5-8 build on it (interfaces, config values, optional extras, and full application graphs), but the core mechanism never changes: **pass what a type needs into its constructor.**

---

## Why Go Usually Doesn't Need a DI Framework

In Java (Spring) or C# (.NET's built-in DI container), a "DI framework" typically means: you annotate classes, register them with a container at startup, and the container uses reflection to figure out the dependency graph and construct everything for you at runtime.

Go's standard answer is simpler: **write the constructor calls yourself, in `main()`.** This works because of a few things Go has that those ecosystems often don't lean on as heavily:

- **Explicit imports and function calls are already fast to write and read.** `NewService(NewRepository(db))` is not meaningfully harder to write than a container registration, and it's far easier to read - you can see the whole graph in one place instead of tracing annotations across files.
- **The compiler checks everything at compile time.** If `NewHandler` needs a `*Service` and you pass a `*Repository` by mistake, that's a compile error, not a runtime container failure discovered when a request comes in.
- **No reflection, no runtime surprises.** A container that builds your graph via reflection at startup can fail in ways that are hard to trace back to a specific line of source code. A plain function call fails exactly where you'd expect.

Because of this, people say **"a poor man's DI container is just a well-organized `main()` function."** It's not a joke or a workaround - it's the idiomatic Go answer. Here's that "container" for a small, layered app - `Config` → `Repository` → `Service` → `Handler` - built by hand:

```go
package main

import "fmt"

// --- Config layer ---
type Config struct {
    AppName string
}

// --- Repository layer ---
type UserRepository interface {
    FindByID(id int) (string, error)
}

type inMemoryUserRepository struct {
    users map[int]string
}

func NewInMemoryUserRepository() *inMemoryUserRepository {
    return &inMemoryUserRepository{
        users: map[int]string{1: "Alice", 2: "Bob"},
    }
}

func (r *inMemoryUserRepository) FindByID(id int) (string, error) {
    name, ok := r.users[id]
    if !ok {
        return "", fmt.Errorf("user %d not found", id)
    }
    return name, nil
}

// --- Service layer ---
type UserService struct {
    cfg  Config
    repo UserRepository
}

func NewUserService(cfg Config, repo UserRepository) *UserService {
    return &UserService{cfg: cfg, repo: repo}
}

func (s *UserService) Greet(id int) (string, error) {
    name, err := s.repo.FindByID(id)
    if err != nil {
        return "", err
    }
    return fmt.Sprintf("[%s] Hello, %s!", s.cfg.AppName, name), nil
}

// --- Handler layer ---
type UserHandler struct {
    service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
    return &UserHandler{service: service}
}

func (h *UserHandler) HandleGreet(id int) {
    greeting, err := h.service.Greet(id)
    if err != nil {
        fmt.Println("error:", err)
        return
    }
    fmt.Println(greeting)
}

// buildApp is the "composition root" - a hand-rolled container. This is
// the ONLY place that knows how the whole dependency graph fits together.
func buildApp() *UserHandler {
    cfg := Config{AppName: "GreeterApp"}
    repo := NewInMemoryUserRepository()
    service := NewUserService(cfg, repo)
    handler := NewUserHandler(service)
    return handler
}

func main() {
    handler := buildApp()
    handler.HandleGreet(1)
    handler.HandleGreet(2)
    handler.HandleGreet(99)
}
```

Real output:

```
[GreeterApp] Hello, Alice!
[GreeterApp] Hello, Bob!
error: user 99 not found
```

Every dependency in that graph was constructed exactly once, in exactly the right order, by plain function calls anyone can read top to bottom. No annotations, no registration step, no runtime resolution - `buildApp` **is** the dependency graph.

---

## Interface-Based Injection: Swapping Implementations

Constructor injection only earns its keep when the parameter type is an **interface** (Level 15), not a concrete struct. Depending on a concrete type still lets you inject *a* value, but only ever one kind of value - you lose the swappability that makes DI worth doing.

```go
package main

import "fmt"

// PaymentGateway is the small interface PaymentProcessor depends on.
type PaymentGateway interface {
    Charge(cents int) error
}

type StripeGateway struct{}

func (StripeGateway) Charge(cents int) error {
    fmt.Printf("stripe: charged %d cents\n", cents)
    return nil
}

type PayPalGateway struct{}

func (PayPalGateway) Charge(cents int) error {
    fmt.Printf("paypal: charged %d cents\n", cents)
    return nil
}

// PaymentProcessor depends only on the PaymentGateway interface - it has
// no idea whether it's talking to Stripe, PayPal, or a test fake.
type PaymentProcessor struct {
    gateway PaymentGateway
}

func NewPaymentProcessor(gateway PaymentGateway) *PaymentProcessor {
    return &PaymentProcessor{gateway: gateway}
}

func (p *PaymentProcessor) Pay(cents int) error {
    return p.gateway.Charge(cents)
}

func main() {
    fmt.Println("=== Same PaymentProcessor code, two different gateways ===")

    stripeProcessor := NewPaymentProcessor(StripeGateway{})
    stripeProcessor.Pay(2500)

    paypalProcessor := NewPaymentProcessor(PayPalGateway{})
    paypalProcessor.Pay(2500)
}
```

Real output:

```
=== Same PaymentProcessor code, two different gateways ===
stripe: charged 2500 cents
paypal: charged 2500 cents
```

`PaymentProcessor`'s source code appears exactly once. It never mentions `StripeGateway` or `PayPalGateway` by name - it only knows about `Charge(cents int) error`. That's Level 15's Rule 2 ("accept interfaces, return concrete types") doing real work: `NewPaymentProcessor` *accepts* the `PaymentGateway` interface, so any current or future implementation - including a test fake - plugs in without a single line of `PaymentProcessor` changing.

---

## Injecting Configuration Alongside Interfaces

Not everything a type depends on needs to be an interface. Plain configuration - a `Host string`, a `MaxRetries int`, a `TTL time.Duration` - has no behavior to mock. It's just data. It still gets **injected** the same way (as a constructor argument), it just doesn't need an interface wrapped around it.

```go
package main

import "fmt"

// ServerConfig is a plain value - not an interface, not mockable behavior,
// just configuration data. It gets injected too.
type ServerConfig struct {
    Host string
    Port int
}

type Store interface {
    Get(key string) (string, bool)
}

type memStore struct{ data map[string]string }

func (m *memStore) Get(key string) (string, bool) {
    v, ok := m.data[key]
    return v, ok
}

// Server depends on BOTH an interface (Store, mockable) and a plain
// config struct (not mockable - just data).
type Server struct {
    cfg   ServerConfig
    store Store
}

func NewServer(cfg ServerConfig, store Store) *Server {
    return &Server{cfg: cfg, store: store}
}

func (s *Server) Address() string {
    return fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
}

func main() {
    cfg := ServerConfig{Host: "localhost", Port: 8080}
    store := &memStore{data: map[string]string{"key1": "value1"}}

    srv := NewServer(cfg, store)
    fmt.Println("Address:", srv.Address())

    if v, ok := srv.store.Get("key1"); ok {
        fmt.Println("Store lookup:", v)
    }
}
```

Real output:

```
Address: localhost:8080
Store lookup: value1
```

The rule of thumb: inject an **interface** when a dependency has *behavior* you might want to fake or swap. Inject a **plain struct or value** when a dependency is just *data* - configuration, options, constants read from the environment. Both go through the same constructor-argument mechanism; only the type differs.

---

## Optional Dependencies: The Functional Options Pattern

Constructor injection handles *required* dependencies cleanly, but what about dependencies (or settings) that are optional, or that keep growing over a type's lifetime? Adding a new required parameter to every constructor call every time you add a knob doesn't scale. Go's idiomatic answer is the **functional options pattern**:

```go
func NewThing(dep SomeInterface, opts ...Option) *Thing
```

`Option` is a function that mutates the thing being built. Callers who don't care about the extras just omit `opts` entirely; callers who do care pass however many `With...` functions they need, in any order.

```go
package main

import (
    "fmt"
    "time"
)

type Store interface {
    Get(key string) (string, bool)
}

type memStore struct{}

func (memStore) Get(key string) (string, bool) { return "value", true }

// Cache depends on the required Store interface via a normal constructor
// parameter, and accepts any number of OPTIONAL extras through the
// functional options pattern - no exploding parameter list.
type Cache struct {
    store   Store
    ttl     time.Duration
    maxSize int
}

type Option func(*Cache)

func WithTTL(ttl time.Duration) Option {
    return func(c *Cache) { c.ttl = ttl }
}

func WithMaxSize(size int) Option {
    return func(c *Cache) { c.maxSize = size }
}

func NewCache(store Store, opts ...Option) *Cache {
    c := &Cache{
        store:   store,
        ttl:     5 * time.Minute, // sensible defaults
        maxSize: 1000,
    }
    for _, opt := range opts {
        opt(c)
    }
    return c
}

func main() {
    fmt.Println("=== Default options ===")
    defaultCache := NewCache(memStore{})
    fmt.Printf("ttl=%s maxSize=%d\n", defaultCache.ttl, defaultCache.maxSize)

    fmt.Println("\n=== Customized via functional options ===")
    customCache := NewCache(memStore{}, WithTTL(30*time.Second), WithMaxSize(50))
    fmt.Printf("ttl=%s maxSize=%d\n", customCache.ttl, customCache.maxSize)
}
```

Real output:

```
=== Default options ===
ttl=5m0s maxSize=1000

=== Customized via functional options ===
ttl=30s maxSize=50
```

Notice `NewCache`'s signature never has to change to support a new optional setting - just add another `WithX` function. This is the standard library's own pattern (`grpc.Dial(target, opts ...grpc.DialOption)`, `http.NewRequestWithContext` variants elsewhere) and it composes perfectly with constructor injection: `dep` stays a required, explicit argument; everything optional rides in through `opts`.

---

## A Hand-Rolled "Container" in main()

Section 4 showed a small three-layer graph. Real applications - like the Repository-Service-Handler apps from Level 30 - often need a **container**: one function whose entire job is constructing the full dependency graph, in the right order, and handing back the top-level piece (usually a handler or a server) ready to run.

A DI framework in another language would build this graph automatically via reflection and registration. In Go, you write the function:

```go
func buildApp() *AccountHandler {
    cfg := Config{MinBalance: 50.00}
    repo := NewInMemoryAccountRepository()
    service := NewAccountService(cfg, repo)
    handler := NewAccountHandler(service)
    return handler
}
```

Read top to bottom, this **is** the dependency graph: `Config` needs nothing, `Repository` needs nothing, `Service` needs `Config` and `Repository`, `Handler` needs `Service`. Each line constructs exactly one thing and hands it to the next constructor. There's no way to get the order wrong without the compiler telling you immediately - you can't pass `repo` to `NewAccountHandler` when it wants a `*AccountService`.

This function is sometimes called the **composition root**: the single place in the whole program where concrete types get chosen and wired together. Every other function in the program - `AccountService`, `AccountHandler`, and everything below them - only ever sees interfaces and plain data, never `buildApp` itself. Section 8 of `EXERCISES.md` builds this exact graph end to end and runs it for real.

---

## When a DI Framework Might Actually Help

None of this means third-party DI frameworks (e.g. `google/wire`, `uber-go/fx`, `uber-go/dig`) are never useful in Go - they exist and see real production use. The honest guidance:

**A framework might help when:**
- The dependency graph is genuinely large and deep (dozens of services, several environments' worth of swapped implementations) and hand-wiring `main()` has become hundreds of lines that are hard to review.
- Multiple teams need to add new dependencies to a shared graph regularly, and a compile-time code generator (like `wire`, which generates the same explicit `NewX(NewY(...))` code you'd write by hand) keeps things consistent without runtime reflection cost.

**It's needless complexity when:**
- Your program has a handful of services, like every exercise in this level and most real Go services below a certain size. Hand-wiring is more readable, not less, at this scale.
- You reach for a container "because that's what you do in Java/C#" rather than because `main()` has actually become unmanageable.

Notice that even the popular Go option (`wire`) generates plain constructor-call code at compile time rather than resolving dependencies via reflection at runtime - it automates *writing* the composition root, not replacing the idea of one. That's a sign of how deeply "explicit constructor calls in `main()`" is baked into idiomatic Go, even in tools built to reduce its tedium at scale.

---

## Best Practices

### 1. Depend on Interfaces, Not Concretions

```go
// ✅ Good - can swap in any implementation, including a test fake
func NewService(repo Repository) *Service

// ❌ Avoid - locked to one concrete type forever
func NewService(repo *PostgresRepository) *Service
```

### 2. Construct Dependencies at the Top, Pass Them Down

Build concrete types once, near `main()` (the composition root), and pass the results down through constructors. Don't let a type three layers deep reach back up and construct its own dependency - that hides coupling exactly where Section 1's "before" example did.

### 3. Keep Constructors Simple and Explicit

A constructor's job is to store what it's given and maybe validate it (Level 13) - not to do real work, spawn goroutines, or make network calls. If building a dependency is expensive, build it once in the composition root and inject the *result*.

### 4. Avoid Global/Package-Level Singletons That Hide Dependencies

```go
// ❌ Avoid - every caller silently depends on this global, invisibly
var DB *sql.DB

func GetUser(id int) (*User, error) {
    return DB.Query(...) // hidden dependency, hard to fake, hard to run in parallel
}

// ✅ Good - the dependency is visible in the signature and swappable
type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}
```

A function signature that lists everything it needs is honest about its dependencies. A function that quietly reaches for a package-level global is not - and Section 9 of `EXERCISES.md` shows exactly how that dishonesty turns into broken tests under real concurrency.

---

## Common Mistakes

### Mistake 1: A Type Constructs Its Own Dependencies Internally

```go
// ❌ WRONG - hidden coupling; ReportGenerator can never use anything but
// *ConsoleLogger, and no test can intercept what it logs
type ReportGenerator struct {
    logger *ConsoleLogger
}

func NewReportGenerator() *ReportGenerator {
    return &ReportGenerator{logger: &ConsoleLogger{}}
}

// ✅ RIGHT - the dependency is received, not constructed
type ReportGenerator struct {
    logger Logger
}

func NewReportGenerator(logger Logger) *ReportGenerator {
    return &ReportGenerator{logger: logger}
}
```

### Mistake 2: Depending on a Concrete Struct Instead of an Interface

```go
// ❌ WRONG - can't substitute a fake in tests; every test pays the real cost
type Service struct {
    repo *PostgresRepository
}

// ✅ RIGHT - any type satisfying the interface works, including test fakes
type Service struct {
    repo Repository // an interface
}
```

### Mistake 3: Global Mutable State Instead of Injection

```go
// ❌ WRONG - hidden dependency, and a genuine data race if two goroutines
// (or two parallel tests) touch it at once
var globalHits = map[string]int{}

func RecordHit(page string) {
    globalHits[page]++
}

// ✅ RIGHT - each caller gets its own isolated, injected instance
type HitTracker struct {
    hits map[string]int
}

func NewHitTracker() *HitTracker {
    return &HitTracker{hits: map[string]int{}}
}
```

Package-level mutable state is invisible in every function signature that touches it - you can't tell a function depends on it just by reading its parameters, and two tests (or two goroutines) touching it concurrently can corrupt it or trigger Go's own `-race` detector. Exercise 9 in this level's `EXERCISES.md` captures this as a real, reproducible data race.

### Mistake 4: Over-Engineering With Unnecessary Interfaces "Just in Case"

```go
// ❌ WRONG - one implementation, one consumer, pure indirection
type Adder interface {
    Add(a, b int) int
}

type realAdder struct{}

func (realAdder) Add(a, b int) int { return a + b }

func NewCalculator(adder Adder) *Calculator { ... }

// ✅ RIGHT - no interface needed until a second implementation or a real
// testing need actually shows up (Level 15's Rule 4)
func Add(a, b int) int { return a + b }
```

Every interface is a layer of indirection. DI is valuable when a dependency has real behavior worth swapping or faking - a database, a network client, a clock. Wrapping a plain, pure function like `Add` in an interface "for flexibility" adds ceremony with no payoff; nothing about it will ever need a second implementation.

---

## Summary

**The Core Idea:**
- DI = a type receives its dependencies from the outside (constructor arguments) instead of building or reaching for them itself

**Why It Matters:**
- Testability (swap in a fake), flexibility (swap implementations), decoupling (depend on interfaces, not concretions)

**The Pattern:**
- `func NewThing(dep SomeInterface) *Thing { return &Thing{dep: dep} }` is essentially all of Go DI
- Plain config values get injected the same way, without needing an interface
- `opts ...Option` (functional options) handles optional extras without exploding the parameter list

**No Framework Needed:**
- `main()` (or a small `buildApp()` function) is a hand-rolled composition root - explicit, compiler-checked, readable top to bottom
- Frameworks exist for very large graphs, but most Go programs never need one

---

## Next Steps

You now understand:
- ✅ What dependency injection means, with no framework mystique attached
- ✅ Constructor injection as the idiomatic Go pattern - nearly all of DI in practice
- ✅ Why explicit wiring in `main()` beats a reflection-based container for most Go programs
- ✅ Interface-based injection, config injection, and the functional options pattern
- ✅ How to hand-wire a layered dependency graph, and when (rarely) a framework might help

**Next level:** Level 32 - Authentication & JWT
- Verifying user identity with JSON Web Tokens
- Issuing, signing, and validating tokens
- Middleware that injects an authenticated user into request handling - this level's constructor-injection habits carry straight over

You've turned "poor man's DI container" from a joke into a real, hand-built skill. Level 32's auth middleware is going to lean on exactly the constructor-injection instincts you just built! 🚀
