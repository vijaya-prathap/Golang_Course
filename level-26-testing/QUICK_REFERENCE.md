# Level 26: Quick Reference Card

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir myapp && cd myapp
go mod init myapp

# Create math.go
cat > math.go << 'EOF'
package mathutil

func Add(a, b int) int {
    return a + b
}
EOF

# Create math_test.go
cat > math_test.go << 'EOF'
package mathutil

import "testing"

func TestAdd(t *testing.T) {
    if got := Add(2, 3); got != 5 {
        t.Errorf("Add(2, 3) = %d; want 5", got)
    }
}
EOF

# Run
go test -v ./...
```

---

## 📋 Test Function Rules

```go
func TestXxx(t *testing.T) { }

✅ Name starts with "Test"
✅ Capital letter (or nothing) right after "Test"
✅ Exactly one parameter: *testing.T
✅ Lives in a file ending in _test.go
```

---

## ⚖️ t.Error vs t.Fatal

```go
t.Error(args...)             // fail, KEEP running
t.Errorf(format, args...)    // fail, KEEP running
t.Fatal(args...)             // fail, STOP this test now
t.Fatalf(format, args...)    // fail, STOP this test now
```

```go
// Use Fatal when continuing would be unsafe/meaningless:
cfg, err := ParseConfig(input)
if err != nil {
    t.Fatalf("ParseConfig: %v", err) // cfg may be nil - must stop
}
if cfg.Timeout != 30 {
    t.Errorf("Timeout = %d; want 30", cfg.Timeout)
}
```

---

## 🧱 Table-Driven Tests + Subtests

```go
func TestMax(t *testing.T) {
    tests := []struct {
        name  string
        input []int
        want  int
    }{
        {"positive", []int{1, 5, 3}, 5},
        {"negative", []int{-5, -2, -9}, -2},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := Max(tt.input); got != tt.want {
                t.Errorf("Max(%v) = %d; want %d", tt.input, got, tt.want)
            }
        })
    }
}
```

```bash
go test -v                              # everything
go test -v -run 'TestMax/negative'      # one subtest only
go test -run Max                        # regex match on test names
```

---

## 📊 Coverage

```bash
go test -cover                              # quick percentage
go test -coverprofile=coverage.out          # write profile
go tool cover -func=coverage.out            # per-function %
go tool cover -html=coverage.out -o out.html  # visual HTML report
```

---

## 🏎️ Benchmarks

```go
func BenchmarkReverse(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Reverse("some input")
    }
}
```

```bash
go test -bench=.              # run all benchmarks (also runs regular tests)
go test -bench=. -run=^$      # run ONLY benchmarks, skip regular tests
go test -bench=. -benchmem    # also show B/op and allocs/op
```

`BenchmarkReverse-10   3508375   333.8 ns/op   224 B/op   2 allocs/op`
→ name-GOMAXPROCS, iterations, ns/call, bytes/call, allocs/call — **all machine-dependent**.

---

## 🧰 Test Helpers

```go
func assertCents(t *testing.T, b *Balance, want int) {
    t.Helper() // failures blame the CALLER's line, not this one
    if got := b.Cents(); got != want {
        t.Errorf("balance = %d; want %d", got, want)
    }
}
```

---

## 🏗️ Setup / Teardown

```go
// Package-level, runs ONCE for the whole file's package:
func TestMain(m *testing.M) {
    setup()
    code := m.Run()
    teardown()
    os.Exit(code)
}

// Per-test, via a helper:
func newTempStore(t *testing.T) *Store {
    t.Helper()
    s := NewStore()
    t.Cleanup(func() { /* runs when THIS test finishes */ })
    return s
}
```

---

## 🌐 httptest Preview

```go
req := httptest.NewRequest(http.MethodGet, "/greet?name=Gopher", nil)
rec := httptest.NewRecorder()

GreetHandler(rec, req)

if rec.Code != http.StatusOK { /* ... */ }
if rec.Body.String() != "Hello, Gopher!" { /* ... */ }
```

Full depth (routers, middleware, `httptest.NewServer`) → Level 27.

---

## 🎭 Mocking via Interfaces

```go
type Notifier interface {          // small, consumer-defined interface
    Notify(message string) error
}

type fakeNotifier struct{ messages []string }

func (f *fakeNotifier) Notify(m string) error {
    f.messages = append(f.messages, m)
    return nil
}

svc := NewOrderService(&fakeNotifier{}) // inject the fake in tests
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Ignoring errors in tests | `cfg, _ := ParseConfig(x)` | `cfg, err := ...; if err != nil { t.Fatalf(...) }` |
| t.Parallel() + shared state | mutating a shared var across parallel tests | isolate state per test |
| Testing private internals | testing an unexported helper directly | test the exported behavior it supports |
| Table without subtests | one `t.Errorf` per loop, no `t.Run` | wrap each case in `t.Run(tt.name, ...)` |

---

## 🎓 Before Next Level

Can you:
- [ ] Name a test function correctly from memory?
- [ ] Explain t.Error vs t.Fatal without hesitation?
- [ ] Write a table-driven test with subtests?
- [ ] Run a single subtest with `-run`?
- [ ] Read a `go test -cover` and a benchmark's `ns/op` line?
- [ ] Use `t.Helper()` and `t.Cleanup()` correctly?
- [ ] Mock a dependency with a small interface + fake struct?

If YES → You're ready for Level 27!

---

## 📚 Next Level

Level 27: HTTP & REST APIs
- Building HTTP servers and handlers
- Routing, middleware, and request/response handling
- Testing HTTP services in depth

You've got testing down! 💪
