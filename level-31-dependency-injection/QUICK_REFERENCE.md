# Level 31: Quick Reference Card

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir myapp && cd myapp
go mod init myapp

# Create main.go
cat > main.go << 'EOF'
package main

import "fmt"

type Logger interface {
    Log(msg string)
}

type ConsoleLogger struct{}

func (ConsoleLogger) Log(msg string) {
    fmt.Println("[LOG]", msg)
}

// NewThing(dep SomeInterface) *Thing - the idiomatic Go DI pattern
type App struct {
    logger Logger
}

func NewApp(logger Logger) *App {
    return &App{logger: logger} // RECEIVED, not constructed
}

func (a *App) Start() {
    a.logger.Log("starting")
}

func main() {
    app := NewApp(ConsoleLogger{})
    app.Start()
}
EOF

# Run
go run main.go
```

---

## 📋 The Core Idea

```
DI = a type receives its dependencies from the OUTSIDE
     (as constructor arguments) instead of building or
     reaching for them itself.
```

```go
// ❌ Before - builds its own dependency
func NewThing() *Thing {
    return &Thing{dep: &ConcreteDep{}}
}

// ✅ After - receives its dependency
func NewThing(dep SomeInterface) *Thing {
    return &Thing{dep: dep}
}
```

---

## 🏗️ Constructor Injection

```go
func NewThing(dep SomeInterface) *Thing {
    return &Thing{dep: dep}
}
```

This is essentially ALL of DI in idiomatic Go. Depend on the interface, not the concrete type, and pass it in through the constructor (Level 13's `NewX` convention).

---

## 🔀 Interface-Based Injection (Swap Implementations)

```go
type Gateway interface {
    Charge(cents int) error
}

// Same NewProcessor, same Processor code - any Gateway works
proc1 := NewProcessor(StripeGateway{})
proc2 := NewProcessor(PayPalGateway{})
proc3 := NewProcessor(fakeGateway{}) // in tests
```

---

## ⚙️ Injecting Config Alongside an Interface

```go
type Config struct {   // plain data - no interface needed
    Host string
    Port int
}

func NewServer(cfg Config, store Store) *Server {
    return &Server{cfg: cfg, store: store}
}
```

Interface = behavior you might fake/swap. Plain struct = just data.

---

## 🎛️ Functional Options (Optional Dependencies)

```go
type Option func(*Cache)

func WithTTL(ttl time.Duration) Option {
    return func(c *Cache) { c.ttl = ttl }
}

func NewCache(store Store, opts ...Option) *Cache {
    c := &Cache{store: store, ttl: 5 * time.Minute} // defaults
    for _, opt := range opts {
        opt(c)
    }
    return c
}

// Usage:
NewCache(store)                          // defaults
NewCache(store, WithTTL(30*time.Second)) // one override
```

---

## 🧱 Hand-Wiring a Dependency Graph (main() as Composition Root)

```go
func buildApp() *Handler {
    cfg := Config{...}
    repo := NewRepository()
    service := NewService(cfg, repo)
    handler := NewHandler(service)
    return handler
}
```

No framework needed - explicit function calls, checked by the compiler, readable top to bottom. This IS a "poor man's DI container."

---

## 🧪 Testing With Fakes

```go
type fakeSender struct {
    messages []string
}

func (f *fakeSender) Send(msg string) error {
    f.messages = append(f.messages, msg)
    return nil
}

func TestService(t *testing.T) {
    fake := &fakeSender{}
    svc := NewService(fake) // injected fake, no real dependency touched
    ...
}
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Hidden construction | `&Thing{dep: &ConcreteDep{}}` inside `NewThing()` | `NewThing(dep SomeInterface)` |
| Concrete dependency | `repo *PostgresRepository` | `repo Repository` (interface) |
| Global mutable state | `var globalStore = map[string]int{}` | inject a `*Store` instance instead |
| Interfaces "just in case" | wrapping a pure `Add(a, b int) int` in an interface | keep it a plain function until a real need appears |

---

## 🎓 Before Next Level

Can you:
- [ ] Write `func NewThing(dep SomeInterface) *Thing` from memory?
- [ ] Explain why a hardcoded dependency can't be faked in a test?
- [ ] Swap two interface implementations without touching the dependent code?
- [ ] Inject a config struct alongside an interface dependency?
- [ ] Write a functional-options constructor for optional extras?
- [ ] Hand-wire a layered graph (Config → Repository → Service → Handler) in one function?

If YES → You're ready for Level 32!

---

## 📚 Next Level

Level 32: Authentication & JWT
- Verifying user identity with JSON Web Tokens
- Issuing, signing, and validating tokens
- Middleware that injects an authenticated user into request handling

You've got dependency injection down - the Go way! 💪
