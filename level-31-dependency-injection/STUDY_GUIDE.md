# Level 31: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: Dependency Injection Concepts
```
Day 1:  What DI means - before/after hardcoded vs injected
Day 2:  Constructor injection pattern (NewThing(dep))
Day 3:  Why Go skips DI frameworks - main() as composition root
Day 4:  Interface-based injection - swapping implementations
Day 5:  Injecting config alongside interfaces
Day 6:  Functional options for optional dependencies
Day 7:  Hand-rolled containers and when frameworks help
```

### Week 2: Practice & Application
```
Day 1:  Exercises 1-3 (before/after, constructor basics, testability proof)
Day 2:  Exercises 4-6 (interface swap, config injection, functional options)
Day 3:  Exercises 7-8 (two fakes, layered graph wiring)
Day 4:  Exercises 9-10 (global state hazard, comprehensive notification system)
Day 5:  Bonus challenges & review
```

---

## 🎯 Hardcoded Dependency vs Injected Dependency

```
❌ HARDCODED (the type builds its own dependency)

    ┌─────────────────────────┐
    │   ReportGenerator        │
    │                          │
    │   logger := &ConsoleLogger{}   ← built INSIDE, every time
    │                          │
    └─────────────────────────┘
    Locked to *ConsoleLogger forever. No test can intercept it.


✅ INJECTED (the type receives its dependency)

    caller ──── passes in ────▶  ┌─────────────────────────┐
    (real Logger OR test fake)   │   ReportGenerator        │
                                  │                          │
                                  │   logger Logger  ← FIELD, set by constructor
                                  │                          │
                                  └─────────────────────────┘
    Works with ANY type satisfying Logger - real or fake.
```

```go
// ❌ Before
func NewReportGenerator() *ReportGenerator {
    return &ReportGenerator{logger: &ConsoleLogger{}} // reaches for it
}

// ✅ After
func NewReportGenerator(logger Logger) *ReportGenerator {
    return &ReportGenerator{logger: logger} // receives it
}
```

---

## 🧩 Constructor Injection Pattern Anatomy

```go
func NewThing(dep SomeInterface) *Thing {
    //  ▲         ▲       ▲              ▲
    //  │         │       │              └─ pointer to a ready-to-use struct
    //  │         │       └─ the TYPE is an interface, not a concrete struct
    //  │         └─ the dependency arrives as a PARAMETER
    //  └─ Level 13's NewX convention
    return &Thing{dep: dep}
    //            ▲
    //            └─ stored, never constructed here
}
```

| Piece | Role |
|---|---|
| `NewThing(...)` | Level 13's constructor convention - the obvious way to build the type |
| `dep SomeInterface` | the dependency, typed as an interface so any implementation fits |
| `&Thing{dep: dep}` | dependency is **stored**, never built inside the constructor |

This shape is basically all of DI in idiomatic Go. Everything else in this level (config values, optional extras, layered graphs) is this same shape applied more than once.

---

## 🔀 Interface Swap: Same Dependent Code, Two Implementations

```
                     ┌───────────────────┐
     StripeGateway ──▶                   │
                     │  PaymentProcessor  │──▶ Pay(cents)
     PayPalGateway ──▶  (unchanged code)  │
                     │                   │
     fakeGateway   ──▶                   │  (in tests)
                     └───────────────────┘

PaymentProcessor's source mentions ONLY:
    type PaymentGateway interface { Charge(cents int) error }

It never names StripeGateway, PayPalGateway, or any fake directly.
```

```go
stripeProcessor := NewPaymentProcessor(StripeGateway{})
paypalProcessor := NewPaymentProcessor(PayPalGateway{})
// same NewPaymentProcessor, same PaymentProcessor.Pay - different gateway
```

---

## 🏗️ Hand-Wired Dependency Graph: main() as the Composition Root

```
buildApp()  ─── the ONLY function that knows every concrete type
    │
    ├─▶ cfg     := Config{...}                        (plain data)
    │
    ├─▶ repo    := NewInMemoryAccountRepository()      (no dependencies)
    │
    ├─▶ service := NewAccountService(cfg, repo)         (needs cfg + repo)
    │
    └─▶ handler := NewAccountHandler(service)            (needs service)
            │
            ▼
        returned to main(), ready to run
```

```
Config           →  no dependencies
Repository       →  no dependencies
Service          →  needs Config + Repository
Handler          →  needs Service
```

Read top to bottom, `buildApp()` **is** the dependency graph. No registration step, no reflection, no framework - just function calls in the right order, checked by the compiler.

---

## 🎛️ Functional Options Pattern

```go
func NewCache(store Store, opts ...Option) *Cache
                 ▲            ▲
                 │            └─ zero or more optional extras
                 └─ the one REQUIRED dependency, always explicit
```

```
NewCache(store)                              → all defaults
NewCache(store, WithTTL(30*time.Second))     → one override, rest default
NewCache(store, WithTTL(x), WithMaxSize(y))  → both overridden, any order
```

```
type Option func(*Cache)              ← an Option just mutates the built value

func WithTTL(ttl time.Duration) Option {
    return func(c *Cache) { c.ttl = ttl }
}
```

**Why not just more constructor parameters?** Every new optional setting would force every existing call site to be rewritten (or force you to pass zero values everywhere). Functional options let old call sites keep compiling unchanged forever.

---

## 🚨 Common Mistakes

### Mistake 1: Hidden Construction Instead of Injection

```go
// ❌ Type builds its own dependency - hidden coupling
func NewReportGenerator() *ReportGenerator {
    return &ReportGenerator{logger: &ConsoleLogger{}}
}
```

### Mistake 2: Depending on a Concrete Type

```go
// ❌ Can't substitute a fake in tests
type Service struct {
    repo *PostgresRepository
}
```

### Mistake 3: Global Mutable State

```go
// ❌ Hidden dependency AND a real data race under concurrency
var globalHits = map[string]int{}
func RecordHit(page string) { globalHits[page]++ }
```

### Mistake 4: Interfaces "Just in Case"

```go
// ❌ One implementation, one consumer - pure indirection
type Adder interface{ Add(a, b int) int }
```

---

## 📈 Progression Summary

### Understanding Level 31

Level 31 turns Level 30's preview (wiring Repository → Service → Handler in `main()`) into a named, deliberate practice:

1. **The core shift** — receive dependencies, don't construct them
2. **Constructor injection** — `NewThing(dep)`, nearly all of Go DI
3. **No framework needed** — `main()` is the composition root
4. **Interfaces make it swappable** — same code, different implementations
5. **Config vs behavior** — plain values injected too, not just interfaces
6. **Functional options** — optional extras without parameter-list explosion

### Prerequisites for Level 32

Before moving to Level 32 (Authentication & JWT), you need:

- ✅ Comfortable writing `func NewThing(dep SomeInterface) *Thing`
- ✅ Can explain why a hardcoded dependency breaks testability
- ✅ Can hand-wire a layered dependency graph inside a `buildApp()`/`main()` composition root
- ✅ Understand the functional options pattern for optional dependencies
- ✅ Recognize global mutable state as a hidden dependency and concurrency hazard

### Ready for Level 32?

Level 32 builds directly on this level's habits:
- Middleware that authenticates a request and injects the resulting user into the next handler - constructor-injection thinking applied to request-scoped values
- Services depending on a `TokenValidator` interface, swappable exactly like this level's `PaymentGateway` or `Sender`
- The same composition-root wiring, now with an auth layer added to the graph

---

## ✅ Checklist Before Level 32

- [ ] Can refactor a hardcoded dependency into a constructor-injected one
- [ ] Can write the `NewThing(dep SomeInterface) *Thing` pattern from memory
- [ ] Can explain testability, flexibility, and decoupling as DI's three payoffs
- [ ] Can inject a plain config struct alongside an interface dependency
- [ ] Can write a functional-options constructor for optional extras
- [ ] Can hand-wire a 3+ layer dependency graph inside one composition-root function
- [ ] Understand why Go usually skips DI frameworks in favor of explicit `main()` wiring
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The One Sentence
A type should receive its dependencies from the outside instead of constructing or reaching for them itself.

### The Pattern
`func NewThing(dep SomeInterface) *Thing { return &Thing{dep: dep} }` - nearly all of Go DI, in one line.

### The Payoff
Testability (fakes in, real dependencies out), flexibility (swap implementations), decoupling (depend on interfaces, not concretions).

### The Non-Framework
`main()` (or a small `buildApp()`) is a hand-rolled composition root - explicit, compiler-checked, and usually all a Go program needs.

---

## 📚 Next Level

Level 32: Authentication & JWT
- Verifying user identity with JSON Web Tokens
- Issuing, signing, and validating tokens
- Middleware that injects an authenticated user into request handling

You've mastered dependency injection the Go way - explicit, simple, and framework-free! 🚀
