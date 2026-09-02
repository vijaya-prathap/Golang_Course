# Level 26: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: Testing Fundamentals
```
Day 1:  The testing package, go test, naming rules
Day 2:  t.Error vs t.Fatal
Day 3:  Table-driven tests
Day 4:  Subtests with t.Run and -run
Day 5:  Coverage: go test -cover, go tool cover
Day 6:  Benchmarks: BenchmarkXxx, go test -bench
Day 7:  Review & practice
```

### Week 2: Real-World Testing Practice
```
Day 1:  Exercises 1-3 (basic tests, Error/Fatal, table-driven bug hunt)
Day 2:  Exercises 4-6 (subtests, coverage, benchmarks)
Day 3:  Exercises 7-8 (helpers with t.Helper, TestMain/t.Cleanup)
Day 4:  Exercises 9-10 (httptest, mocking via interfaces)
Day 5:  Bonus challenges (generics, fuzzing, golden files)
Day 6-7: Consolidation & review
```

---

## 🗂️ Test File Naming and Structure

```
myapp/
├── math.go            ← production code, package mathutil
└── math_test.go       ← tests, SAME package mathutil (or mathutil_test for black-box)

Rules enforced by the toolchain:
  ✅ File MUST end in _test.go             (go build ignores it; go test compiles it)
  ✅ Function name MUST start with "Test"   TestAdd, TestParseConfig ✅   testAdd ❌ (lowercase - ignored)
  ✅ Character after "Test" MUST be         TestHTTPClient ✅   Testadd ❌ (lowercase - not a valid test name)
     uppercase (or nothing)
  ✅ MUST take exactly one parameter:       func TestAdd(t *testing.T) { }
     *testing.T
```

```
┌─────────────────────────────────────────────────────────┐
│  func TestAdd(t *testing.T) {                            │
│       ^^^^                ^                               │
│       │                   └── the ONE required parameter  │
│       └── must start with "Test" + capital letter          │
│                                                             │
│      result := Add(2, 3)     ← call the code under test    │
│      if result != 5 {                                      │
│          t.Errorf(...)       ← report failure if wrong     │
│      }                                                      │
│  }                                                          │
└─────────────────────────────────────────────────────────┘
```

---

## ⚖️ t.Error vs t.Fatal

| | `t.Error` / `t.Errorf` | `t.Fatal` / `t.Fatalf` |
|---|---|---|
| Marks test failed | ✅ Yes | ✅ Yes |
| Continues running the test | ✅ Yes | ❌ No - stops immediately |
| Mechanism | records failure, returns normally | records failure, calls `runtime.Goexit()` |
| Best for | independent checks you want ALL reported | a precondition that must hold before continuing safely |
| Risk if used wrong | none really - just noisier failures | using it for a minor check hides later problems in the same test |

```
TestErrorReportsEveryProblem            TestFatalStopsImmediately
─────────────────────────────           ─────────────────────────
check 1 fails → t.Errorf   ┐            check fails → t.Fatalf
check 2 fails → t.Errorf   ├─ all run    (test stops HERE)
t.Log("reached the end")   ┘            t.Log(...) ← NEVER RUNS
     ↓                                       ↓
  test FAILS,                            test FAILS,
  all 3 lines printed                    only 1 line printed
```

Real captured output showing exactly this contrast:

```
=== RUN   TestErrorReportsEveryProblem
    divide_test.go:13: check 1: Divide(10, 4) = 2.5; want 3
    divide_test.go:16: check 2: Divide(10, 4) = 2.5; want 2
    divide_test.go:18: reached the end of the test despite two failed checks above
--- FAIL: TestErrorReportsEveryProblem (0.00s)
=== RUN   TestFatalStopsImmediately
    divide_test.go:26: expected an error for this deliberately wrong assertion, got nil
--- FAIL: TestFatalStopsImmediately (0.00s)
```

---

## 🧱 Table-Driven Test Anatomy

```
┌──────────────────────────────────────────────────────────────┐
│  tests := []struct {                                           │
│      name  string    ← human-readable case description         │
│      input []int     ← input(s) to the function under test      │
│      want  int       ← expected result                          │
│  }{                                                              │
│      {"all positive", []int{3, 1, 4, 1, 5, 9}, 9},               │
│      {"single element", []int{7}, 7},                            │
│      {"all negative", []int{-5, -2, -9, -1}, -1},  ← the case    │
│                                                        that       │
│                                                     catches a     │
│                                                     seeded bug    │
│  }                                                                │
│                                                                    │
│  for _, tt := range tests {                                       │
│      t.Run(tt.name, func(t *testing.T) {   ← named subtest        │
│          got := Max(tt.input)                                     │
│          if got != tt.want {                                      │
│              t.Errorf(...)                                        │
│          }                                                         │
│      })                                                             │
│  }                                                                   │
└───────────────────────────────────────────────────────────────────┘
```

**Why this is the default in Go:** enumerating edge cases as *data* (rows in a table) makes them cheap to add and easy to scan in a code review - a missing case is a missing row, not a missing function.

---

## 🌳 Subtest Naming & Hierarchy

```
go test -v
─────────────────────────────────────────────
TestMax                              ← parent test
├── TestMax/all_positive             ← subtest (spaces → underscores)
├── TestMax/single_element
├── TestMax/all_negative
└── TestMax/mixed_with_zero

Output shows the tree twice:
  === RUN   lines - announce each one starting, in order
  --- PASS/FAIL lines - indented under the parent, one per subtest
```

Real output:

```
=== RUN   TestMax
=== RUN   TestMax/all_positive
=== RUN   TestMax/single_element
=== RUN   TestMax/all_negative
=== RUN   TestMax/mixed_with_zero
--- PASS: TestMax (0.00s)
    --- PASS: TestMax/all_positive (0.00s)
    --- PASS: TestMax/single_element (0.00s)
    --- PASS: TestMax/all_negative (0.00s)
    --- PASS: TestMax/mixed_with_zero (0.00s)
PASS
```

### Running One Subtest

```bash
go test -v -run 'TestIsValidUsername/empty_string'
```

`-run` is a **regular expression** matched against `TestName/SubtestName`. Use it to isolate exactly the failing case you're debugging instead of re-running the whole suite.

---

## 📊 Coverage Output Explained

```bash
go test -cover
```
```
ok  	level26.example/exercise5	0.896s	coverage: 66.7% of statements
```

```
go test -coverprofile=coverage.out    ← writes a machine-readable profile
go tool cover -func=coverage.out      ← per-function %, from that profile
go tool cover -html=coverage.out -o coverage.html   ← visual, line-by-line HTML report
```

```
go tool cover -func output:
─────────────────────────────────────────────
level26.example/exercise5/grade.go:4:  Letter    66.7%
total:                       (statements)         66.7%
                ↑                  ↑                ↑
             file:line       "statements" unit   percentage covered
```

**66.7% here means:** the `A`, `B`, and `F` branches of a 5-branch switch ran; `C` and `D` never did. Coverage tells you WHERE untested code lives - it does not tell you whether the tested code is *correct*.

---

## 🏎️ Benchmark Output Format

```go
func BenchmarkReverse(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Reverse("The quick brown fox jumps over the lazy dog")
    }
}
```

```bash
go test -bench=. -run=^$ -benchmem
```

```
BenchmarkReverse-10    	 3508375	       333.8 ns/op	     224 B/op	       2 allocs/op
       │            │        │                │                │              │
       │            │        │                │                │              └ allocations per call
       │            │        │                │                └ bytes allocated per call
       │            │        │                └ nanoseconds per call (average over all iterations)
       │            │        └ how many times b.N ran (chosen automatically)
       │            └ GOMAXPROCS used for this run
       └ benchmark name
```

🚨 **`ns/op`, allocation counts, and iteration counts are machine- and load-dependent.** Use benchmarks to compare relative performance (approach A vs. approach B, before vs. after a change) on the *same* machine in the *same* run - never hard-code an absolute expected number into an assertion.

---

## 🎭 Mocking via Interfaces (ties to Level 15)

```
┌─────────────────────────────────────────────────────────────┐
│                    type Notifier interface {                  │
│                        Notify(message string) error            │
│                    }                                             │
│                            ▲                    ▲                 │
│              satisfies     │                    │  satisfies       │
│                            │                    │                   │
│         ┌──────────────────┘                    └───────────────┐   │
│         │                                                         │   │
│  EmailNotifier (production)                          fakeNotifier (test)│
│  Notify() sends a real email                Notify() just records the  │
│                                              message in a slice          │
│                                                                            │
│  OrderService depends on the INTERFACE, never on either concrete type:     │
│                                                                              │
│      type OrderService struct { notifier Notifier }                         │
│                                                                                │
│  Production:  NewOrderService(EmailNotifier{...})                             │
│  Tests:       NewOrderService(&fakeNotifier{})    ← no real email ever sent    │
└──────────────────────────────────────────────────────────────────────────────┘
```

This is Level 15's Rule 3 ("define interfaces where they're used") paying off directly: because `Notifier` is tiny (one method) and owned by the consumer (`OrderService`), a test-only fake can satisfy it with zero framework, zero code generation, and zero coupling to the real implementation.

---

## 🚨 Common Mistakes

### Mistake 1: Not Checking Errors in Tests

```go
// ❌ WRONG - cfg could be nil/zero; next line may panic or pass on garbage
cfg, _ := ParseConfig(input)

// ✅ RIGHT
cfg, err := ParseConfig(input)
if err != nil {
    t.Fatalf("ParseConfig returned error: %v", err)
}
```

### Mistake 2: Forgetting What t.Parallel() Implies About Shared State

```go
// ❌ RISKY - a data race if other parallel tests also touch sharedCounter
func TestIncrement(t *testing.T) {
    t.Parallel()
    sharedCounter++
}

// ✅ RIGHT - isolate state per test
func TestIncrement(t *testing.T) {
    t.Parallel()
    counter := &Counter{}
    counter.Increment()
}
```

### Mistake 3: Testing Private Implementation Details

```go
// ❌ WRONG - couples the test to an internal helper name/shape
func TestInternalNormalizeHelper(t *testing.T) { normalizeWhitespace(...) }

// ✅ RIGHT - test the public behavior; internals can change freely
func TestFormatMessage(t *testing.T) { FormatMessage(...) }
```

### Mistake 4: Table-Driven Tests Without Subtests

```go
// ❌ WRONG - a failure gives you no case name, just a line number
for _, tt := range tests {
    if got := Max(tt.input); got != tt.want { t.Errorf(...) }
}

// ✅ RIGHT - t.Run names each case and isolates its pass/fail
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        if got := Max(tt.input); got != tt.want { t.Errorf(...) }
    })
}
```

---

## 📈 Progression Summary

### Understanding Level 26

Level 26 teaches how to prove your code works, and keep proving it:

1. **The convention** — `_test.go` files, `TestXxx(t *testing.T)`, `go test`
2. **Failure reporting** — `t.Error*` continues, `t.Fatal*` stops
3. **The idiomatic shape** — table-driven tests + `t.Run` subtests
4. **Measuring** — coverage (what ran) and benchmarks (how fast)
5. **Maintainability** — `t.Helper()`, `TestMain`/`t.Cleanup`, mocking via interfaces

### Prerequisites for Level 27

Before moving to Level 27 (HTTP & REST APIs), you need:

- ✅ Comfortable writing and running `TestXxx` functions
- ✅ Know when to use `t.Error*` vs `t.Fatal*`
- ✅ Can write a table-driven test with subtests, and run one with `-run`
- ✅ Can measure coverage and read a benchmark's `ns/op`
- ✅ Can extract a test helper correctly with `t.Helper()`
- ✅ Comfortable with `httptest.NewRecorder`/`NewRequest` for basic handler tests
- ✅ Can mock a dependency via a small interface and a fake implementation

### Ready for Level 27?

Level 27 builds real HTTP servers and REST APIs - and tests them thoroughly using the `httptest` foundation you just built:
- Building HTTP servers and handlers
- Routing and middleware
- Request/response handling
- Testing HTTP services in depth

---

## ✅ Checklist Before Level 27

- [ ] Can name a test function so `go test` recognizes it
- [ ] Can explain the difference between `t.Error` and `t.Fatal` from memory
- [ ] Can write a table-driven test and use `t.Run` for each case
- [ ] Can run a single subtest with `-run 'Test/subtest'`
- [ ] Can run `go test -cover` and explain what the percentage means
- [ ] Can write a `BenchmarkXxx` function and read `go test -bench=.` output
- [ ] Uses `t.Helper()` in every test helper function
- [ ] Can set up package-level state with `TestMain` or per-test state with `t.Cleanup`
- [ ] Can test an HTTP handler with `httptest`
- [ ] Can mock a dependency using a small interface and a fake struct
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The Convention
`_test.go` files, `TestXxx(t *testing.T)`, `go test` - no framework required, built into the toolchain.

### The Choice
`t.Error*` keeps going (see everything); `t.Fatal*` stops now (avoid invalid state).

### The Shape
Table-driven tests + `t.Run` subtests is the default - not an advanced technique, the *normal* way to write a Go test.

### The Measurement
Coverage shows what ran, not what's correct. Benchmarks compare relative speed, not absolute truth - numbers vary by machine.

### The Decoupling
Small interfaces + fake implementations let you test business logic without touching real databases, clocks, networks, or file systems.

---

## 📚 Next Level

Level 27: HTTP & REST APIs
- Building real HTTP servers and handlers
- Routing, middleware, and request/response handling
- Testing HTTP services in depth

You've got testing down! Keep going! 🚀
