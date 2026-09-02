# Level 26: Testing - Complete Guide

## Introduction

Welcome to Level 26! You've mastered generics (Level 25) - writing functions and types that work across many types safely. Now it's time for **testing**, the practice that separates "code that seems to work" from "code you can prove works, and keep proving as it changes."

Testing isn't an optional extra bolted onto Go - it's baked into the toolchain. Every Go installation ships with the `testing` package and the `go test` command, no third-party framework required. This is one of the reasons Go code tends to be so heavily tested in practice: the tools are right there, they're fast, and they're the same tools everyone else uses, so any Go developer can pick up any Go project and immediately know how to run its tests.

This level is overdue and important. Everything you've learned so far - functions (Level 11), error handling (Level 16), interfaces (Level 15), goroutines and channels (Levels 20-21), generics (Level 25) - becomes something you can verify automatically, on every change, forever. Testing is also what makes refactoring safe: a good test suite tells you the instant you break something, instead of leaving you to discover it in production.

---

## Table of Contents

1. [The testing Package and go test](#the-testing-package-and-go-test)
2. [t.Error vs t.Fatal](#terror-vs-tfatal)
3. [Table-Driven Tests](#table-driven-tests)
4. [Subtests with t.Run](#subtests-with-trun)
5. [Test Coverage](#test-coverage)
6. [Benchmarks](#benchmarks)
7. [Test Helpers and t.Helper()](#test-helpers-and-thelper)
8. [Setup and Teardown](#setup-and-teardown)
9. [httptest Package Preview](#httptest-package-preview)
10. [Mocking via Interfaces](#mocking-via-interfaces)
11. [Best Practices](#best-practices)
12. [Common Mistakes](#common-mistakes)

---

## The testing Package and go test

Go's test tooling follows one strict convention: put tests in a file ending in `_test.go`, right next to the code they test, in the same package.

```
myapp/
├── math.go          // the code
└── math_test.go     // its tests
```

`_test.go` files are excluded from normal builds (`go build`, `go run`) - only `go test` compiles and runs them. This means test code can live forever alongside production code with zero runtime cost.

### The Test Function Signature

A test function has three hard requirements:

1. Its name must **start with `Test`**
2. The character immediately after `Test` must be **uppercase** (or nothing) - `TestAdd` and `TestHTTPClient` are valid, `Testadd` is not recognized as a test
3. It must take exactly one parameter: **`t *testing.T`**

```go
package mathutil

import "testing"

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Add(2, 3) = %d; want 5", result)
    }
}
```

There's no assertion library built in - a test is just a function that calls the code under test and reports a failure via `t.Error`/`t.Errorf`/`t.Fatal`/`t.Fatalf` (next section) when the result is wrong. If none of those are called, the test passes.

### Running Tests

```bash
go test              # run tests in the current package
go test ./...        # run tests in the current package AND all subpackages
go test -v           # verbose - show every test name and PASS/FAIL
go test -run Add     # only run tests whose name matches the regex "Add"
```

Real output, run against the `TestAdd` example above:

```
=== RUN   TestAdd
--- PASS: TestAdd (0.00s)
PASS
ok  	level26.example/exercise1	0.547s
```

Each line matters:
- `=== RUN   TestAdd` - only appears with `-v`, announces the test is starting
- `--- PASS: TestAdd (0.00s)` - the test's own result and how long it took
- `PASS` - the whole package's result
- `ok  	level26.example/exercise1	0.547s` - package name and total time (this number is **machine-dependent** - yours will differ)

---

## t.Error vs t.Fatal

`*testing.T` gives you two families of failure-reporting methods, and picking the right one matters:

| Method | Marks test failed? | Stops the test immediately? | Use when... |
|--------|--------------------|------------------------------|-------------|
| `t.Error(args...)` | ✅ Yes | ❌ No - keeps running | You want to see every problem, not just the first |
| `t.Errorf(format, args...)` | ✅ Yes | ❌ No - keeps running | Same as `t.Error`, with `Printf`-style formatting |
| `t.Fatal(args...)` | ✅ Yes | ✅ Yes - stops THIS test | Continuing would panic or produce meaningless results |
| `t.Fatalf(format, args...)` | ✅ Yes | ✅ Yes - stops THIS test | Same as `t.Fatal`, with `Printf`-style formatting |

Both families mark the test as failed - the difference is purely about whether execution continues. `t.Fatal`/`t.Fatalf` work by calling `runtime.Goexit()`, which only stops the **current goroutine** (usually the test function itself) - it doesn't stop other tests, and it must be called from the test's own goroutine, not from a helper goroutine you spawned.

### Seeing Both in Action

```go
package divide

import "testing"

// t.Error/t.Errorf marks the test as failed but lets it keep running,
// so every problem in the test shows up in one run instead of one-at-a-time.
func TestErrorReportsEveryProblem(t *testing.T) {
    result, err := Divide(10, 4)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    if result != 3 { // deliberately wrong expectation, to prove execution continues
        t.Errorf("check 1: Divide(10, 4) = %v; want %v", result, 3.0)
    }
    if result != 2 { // also deliberately wrong
        t.Errorf("check 2: Divide(10, 4) = %v; want %v", result, 2.0)
    }
    t.Log("reached the end of the test despite two failed checks above")
}

// t.Fatal/t.Fatalf marks the test as failed and stops THIS test immediately -
// nothing after it in the same test function runs.
func TestFatalStopsImmediately(t *testing.T) {
    _, err := Divide(10, 4)
    if err == nil {
        t.Fatalf("expected an error for this deliberately wrong assertion, got nil")
    }
    t.Log("this line never runs - t.Fatalf above stopped the test")
}
```

Real captured output (both tests are deliberately seeded to fail, to make the contrast visible):

```
=== RUN   TestErrorReportsEveryProblem
    divide_test.go:13: check 1: Divide(10, 4) = 2.5; want 3
    divide_test.go:16: check 2: Divide(10, 4) = 2.5; want 2
    divide_test.go:18: reached the end of the test despite two failed checks above
--- FAIL: TestErrorReportsEveryProblem (0.00s)
=== RUN   TestFatalStopsImmediately
    divide_test.go:26: expected an error for this deliberately wrong assertion, got nil
--- FAIL: TestFatalStopsImmediately (0.00s)
FAIL
FAIL	level26.example/exercise2	0.543s
FAIL
```

Notice:
- `TestErrorReportsEveryProblem` reports **both** wrong-expectation lines AND the final `t.Log` line - execution ran all the way to the end even though it already failed twice.
- `TestFatalStopsImmediately` reports only the `t.Fatalf` line - the trailing `t.Log` never executes, because `t.Fatalf` stopped the test right there.

💡 A common real-world reason to reach for `t.Fatal`: you got an `err` back and need to stop before dereferencing a `nil` result that the error implies.

```go
cfg, err := ParseConfig(input)
if err != nil {
    t.Fatalf("ParseConfig returned error: %v", err) // must stop - cfg may be nil
}
if cfg.Timeout != 30 { // safe: we already know err was nil
    t.Errorf("Timeout = %d; want 30", cfg.Timeout)
}
```

---

## Table-Driven Tests

The single most idiomatic pattern in Go testing: define your test cases as a **slice of structs**, then loop over them and run the same test logic against each one. This scales from 2 cases to 200 without duplicating test code.

### Anatomy

```go
tests := []struct {
    name  string // human-readable description of the case
    input []int  // input(s) to the function under test
    want  int     // expected result
}{
    {"all positive", []int{3, 1, 4, 1, 5, 9}, 9},
    {"single element", []int{7}, 7},
    {"all negative", []int{-5, -2, -9, -1}, -1},
    {"mixed with zero", []int{-3, 0, -7}, 0},
}

for _, tt := range tests {
    got := Max(tt.input)
    if got != tt.want {
        t.Errorf("Max(%v) = %d; want %d", tt.input, got, tt.want)
    }
}
```

### Catching a Real Bug

Table-driven tests earn their keep by covering edge cases you might not think to write a whole separate test for. Here's a deliberately-buggy `Max` function:

```go
package maxfinder

// Max returns the largest value in nums.
func Max(nums []int) int {
    max := 0 // BUG: should start from nums[0], not 0
    for _, n := range nums {
        if n > max {
            max = n
        }
    }
    return max
}
```

With the table above (including the crucial `"all negative"` case), `go test -v` catches it immediately:

```
=== RUN   TestMax
=== RUN   TestMax/all_positive
=== RUN   TestMax/single_element
=== RUN   TestMax/all_negative
    maxfinder_test.go:21: Max([-5 -2 -9 -1]) = 0; want -1
=== RUN   TestMax/mixed_with_zero
--- FAIL: TestMax (0.00s)
    --- PASS: TestMax/all_positive (0.00s)
    --- PASS: TestMax/single_element (0.00s)
    --- FAIL: TestMax/all_negative (0.00s)
    --- PASS: TestMax/mixed_with_zero (0.00s)
FAIL
FAIL	level26.example/exercise3	3.216s
FAIL
```

The bug: initializing `max := 0` means a slice of all-negative numbers can never beat the starting value, so `Max` wrongly returns `0`. The fix is to seed `max` from the data itself:

```go
func Max(nums []int) int {
    max := nums[0] // FIXED: start from the first element, not 0
    for _, n := range nums {
        if n > max {
            max = n
        }
    }
    return max
}
```

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
ok  	level26.example/exercise3	0.539s
```

This is exactly why table-driven tests are the default in Go: a hand-written single-case test for `Max` might never have included an all-negative slice. A table invites you to enumerate the edge cases as data, which makes them cheap to add and easy to review.

---

## Subtests with t.Run

`t.Run(name, func(t *testing.T) {...})` turns each table row into its own **named subtest** - notice it already appeared in the output above. Subtests give you three things a bare loop doesn't:

1. **Named failures** - you instantly see *which* case failed (`TestMax/all_negative`), not just that the table loop failed somewhere
2. **Independent pass/fail** - one failing case doesn't stop the others from running
3. **Selective execution** - you can run exactly one subtest with `-run`

```go
func TestIsValidUsername(t *testing.T) {
    tests := []struct {
        name  string
        input string
        want  bool
    }{
        {"valid name", "alice99", true},
        {"too short", "ab", false},
        {"too long", "waaaaytoolongname", false},
        {"contains space", "al ice", false},
        {"empty string", "", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := IsValidUsername(tt.input)
            if got != tt.want {
                t.Errorf("IsValidUsername(%q) = %v; want %v", tt.input, got, tt.want)
            }
        })
    }
}
```

Go turns each `name` into a slash-separated path by replacing spaces with underscores: `TestIsValidUsername/empty_string`.

```
=== RUN   TestIsValidUsername
=== RUN   TestIsValidUsername/valid_name
=== RUN   TestIsValidUsername/too_short
=== RUN   TestIsValidUsername/too_long
=== RUN   TestIsValidUsername/contains_space
=== RUN   TestIsValidUsername/empty_string
--- PASS: TestIsValidUsername (0.00s)
    --- PASS: TestIsValidUsername/valid_name (0.00s)
    --- PASS: TestIsValidUsername/too_short (0.00s)
    --- PASS: TestIsValidUsername/too_long (0.00s)
    --- PASS: TestIsValidUsername/contains_space (0.00s)
    --- PASS: TestIsValidUsername/empty_string (0.00s)
PASS
ok  	level26.example/exercise4	0.958s
```

### Running One Subtest

`-run` takes a regular expression matched against `Test.../Subtest...`. To run just `empty_string`:

```bash
go test -v -run 'TestIsValidUsername/empty_string'
```

```
=== RUN   TestIsValidUsername
=== RUN   TestIsValidUsername/empty_string
--- PASS: TestIsValidUsername (0.00s)
    --- PASS: TestIsValidUsername/empty_string (0.00s)
PASS
ok  	level26.example/exercise4	0.248s
```

This is invaluable while debugging: instead of re-running an entire (possibly slow) test file, isolate the one case you're chasing.

---

## Test Coverage

`go test -cover` reports what percentage of your code's **statements** were executed by the tests that just ran. It's a floor, not a target - 100% coverage doesn't mean your tests check the right *behavior*, only that every line ran at least once.

```go
func Letter(score int) string {
    switch {
    case score >= 90:
        return "A"
    case score >= 80:
        return "B"
    case score >= 70:
        return "C"
    case score >= 60:
        return "D"
    default:
        return "F"
    }
}
```

Testing only the `A`, `B`, and `F` branches:

```bash
go test -cover
```

```
ok  	level26.example/exercise5	0.896s	coverage: 66.7% of statements
```

66.7% because the `C` and `D` branches never ran. To see exactly *which* function (and, with more detail, which lines) are undertested, generate a profile and inspect it:

```bash
go test -coverprofile=coverage.out
go tool cover -func=coverage.out
```

```
level26.example/exercise5/grade.go:4:	Letter		66.7%
total:					(statements)	66.7%
```

For a visual, line-by-line view (green = covered, red = not covered) in your browser:

```bash
go tool cover -html=coverage.out -o coverage.html
```

This writes an HTML report you can open directly; CI systems often generate this artifact automatically so reviewers can see coverage changes on a pull request.

🚨 Coverage percentage is a diagnostic, not a goal to chase blindly - see [Common Mistakes](#common-mistakes).

---

## Benchmarks

A benchmark measures how long code takes to run and how much it allocates. It lives in the same `_test.go` file, named `BenchmarkXxx`, and takes `*testing.B` instead of `*testing.T`.

```go
func BenchmarkReverse(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Reverse("The quick brown fox jumps over the lazy dog")
    }
}
```

`b.N` is chosen automatically by the testing framework: it starts small and increases until the benchmark has run long enough to produce a stable measurement. You never set `b.N` yourself - just put the code you want timed inside the loop.

Benchmarks don't run during a normal `go test` - you opt in with `-bench` (which takes a regex, `.` meaning "all benchmarks"), and it's conventional to pair it with `-run=^$` so no regular tests also run:

```bash
go test -bench=. -run=^$
```

```
goos: darwin
goarch: arm64
pkg: level26.example/exercise6
cpu: Apple M1 Pro
BenchmarkReverse-10    	 3444912	       332.6 ns/op
PASS
ok  	level26.example/exercise6	1.746s
```

Reading the line `BenchmarkReverse-10    3444912    332.6 ns/op`:
- `BenchmarkReverse-10` - the benchmark name, with `-10` meaning `GOMAXPROCS=10` (the number of logical CPUs used)
- `3444912` - how many iterations `b.N` reached
- `332.6 ns/op` - nanoseconds per single call to `Reverse`, averaged over all iterations

🚨 **These exact numbers are machine- and load-dependent.** Re-running the same benchmark on the same machine will already vary slightly; running it on a different CPU, a laptop on battery power, or a busy CI runner can produce very different numbers. Never hard-code an expected `ns/op` value into a test assertion - use benchmarks to compare *relative* performance (before vs. after a change, one approach vs. another), not to assert an absolute number.

Add `-benchmem` to also see memory allocations, useful for spotting unnecessary `alloc`s:

```
BenchmarkReverse-10    	 3508375	       333.8 ns/op	     224 B/op	       2 allocs/op
```

---

## Test Helpers and t.Helper()

Real test suites accumulate repeated setup and assertion logic. Extracting it into a helper function keeps tests short - but without `t.Helper()`, a failure inside the helper gets blamed on the **helper's** line, not the test case that actually triggered it, which makes large suites painful to debug.

```go
// assertCents is a test helper. t.Helper() tells the testing package that
// failures here should be blamed on the CALLER's line, not this line.
func assertCents(t *testing.T, b *Balance, want int) {
    t.Helper()
    if got := b.Cents(); got != want {
        t.Errorf("balance = %d cents; want %d cents", got, want)
    }
}

func TestDepositWrongExpectation(t *testing.T) {
    b := &Balance{}
    b.Deposit(500)
    assertCents(t, b, 999) // deliberately wrong, to see WHERE the failure is reported
}
```

**With `t.Helper()`** (as written above), the failure points at the call site inside the test function:

```
=== RUN   TestDepositWrongExpectation
    account_test.go:23: balance = 500 cents; want 999 cents
--- FAIL: TestDepositWrongExpectation (0.00s)
```

Line 23 is `assertCents(t, b, 999)` inside `TestDepositWrongExpectation` itself.

**Without `t.Helper()`** (removing that one line, otherwise identical code), the exact same failure instead points *inside the helper*:

```
=== RUN   TestDepositWrongExpectation
    account_test.go:8: balance = 500 cents; want 999 cents
--- FAIL: TestDepositWrongExpectation (0.00s)
```

Line 8 is the `t.Errorf(...)` call inside `assertCents` - useless information when you have dozens of tests calling the same helper, since every failure would report the identical line no matter which test caused it.

**Rule of thumb:** any function that takes `*testing.T` and calls `t.Error*`/`t.Fatal*` on behalf of a caller should call `t.Helper()` as its first line.

---

## Setup and Teardown

Sometimes a whole file (or package) of tests shares expensive setup - seeding a database, starting a server, building test fixtures. Go gives you two mechanisms depending on the scope.

### Package-Level: TestMain

`func TestMain(m *testing.M)` replaces the default test runner for the whole file's package. It runs **once**, before any test, and you decide when the tests actually run by calling `m.Run()`:

```go
var sharedStore *Store

func TestMain(m *testing.M) {
    fmt.Println("=== SETUP: seeding shared store ===")
    sharedStore = NewStore()
    sharedStore.Set("alice", 100)
    sharedStore.Set("bob", 200)

    code := m.Run()

    fmt.Println("=== TEARDOWN: clearing shared store ===")
    sharedStore = nil

    os.Exit(code)
}

func TestLookupAlice(t *testing.T) {
    got, ok := sharedStore.Get("alice")
    if !ok || got != 100 {
        t.Errorf("Get(\"alice\") = %d, %v; want 100, true", got, ok)
    }
}

func TestLookupBob(t *testing.T) {
    got, ok := sharedStore.Get("bob")
    if !ok || got != 200 {
        t.Errorf("Get(\"bob\") = %d, %v; want 200, true", got, ok)
    }
}
```

```
=== SETUP: seeding shared store ===
=== RUN   TestLookupAlice
--- PASS: TestLookupAlice (0.00s)
=== RUN   TestLookupBob
--- PASS: TestLookupBob (0.00s)
PASS
=== TEARDOWN: clearing shared store ===
ok  	level26.example/exercise8	0.536s
```

Notice `SETUP` prints exactly once (before both tests) and `TEARDOWN` prints exactly once (after both) - not once per test. You must call `os.Exit(code)` yourself; forgetting it means `go test`'s exit status won't reflect whether the tests actually passed.

### Per-Test: A Helper Returning Cleanup, Used With t.Cleanup

When each test needs its *own* fresh setup/teardown (not shared across the whole package), write a helper that registers a cleanup function with `t.Cleanup`. Go runs cleanups in LIFO order, automatically, when the test (or subtest) that registered them finishes - including if it fails or panics.

```go
func newTempStore(t *testing.T) *Store {
    t.Helper()
    s := NewStore()
    t.Cleanup(func() {
        t.Log("cleanup: discarding temp store for " + t.Name())
    })
    return s
}

func TestTempStoreIsIsolated(t *testing.T) {
    s := newTempStore(t)
    s.Set("temp", 1)
    got, ok := s.Get("temp")
    if !ok || got != 1 {
        t.Errorf("Get(\"temp\") = %d, %v; want 1, true", got, ok)
    }
}
```

```
=== RUN   TestTempStoreIsIsolated
    store_test.go:48: cleanup: discarding temp store for TestTempStoreIsIsolated
--- PASS: TestTempStoreIsIsolated (0.00s)
```

`t.Cleanup` is generally preferred over a manual `defer` in the test body when the setup itself lives in a helper - the helper can register its own cleanup without the test author needing to remember to `defer helperCleanup()` themselves.

---

## httptest Package Preview

Testing HTTP handlers doesn't require starting a real server on a real port. The standard library's `net/http/httptest` package gives you an in-memory request and a recorder that captures whatever a handler writes - `Level 27: HTTP & REST APIs` covers this in full depth; this is a preview of the essential pattern.

```go
func GreetHandler(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    if name == "" {
        name = "World"
    }
    w.Header().Set("Content-Type", "text/plain")
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Hello, %s!", name)
}
```

```go
func TestGreetHandler(t *testing.T) {
    tests := []struct {
        name string
        url  string
        want string
    }{
        {"default name", "/greet", "Hello, World!"},
        {"custom name", "/greet?name=Gopher", "Hello, Gopher!"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest(http.MethodGet, tt.url, nil)
            rec := httptest.NewRecorder()

            GreetHandler(rec, req)

            if rec.Code != http.StatusOK {
                t.Errorf("status = %d; want %d", rec.Code, http.StatusOK)
            }
            if got := rec.Body.String(); got != tt.want {
                t.Errorf("body = %q; want %q", got, tt.want)
            }
        })
    }
}
```

```
=== RUN   TestGreetHandler
=== RUN   TestGreetHandler/default_name
=== RUN   TestGreetHandler/custom_name
--- PASS: TestGreetHandler (0.00s)
    --- PASS: TestGreetHandler/default_name (0.00s)
    --- PASS: TestGreetHandler/custom_name (0.00s)
PASS
ok  	level26.example/exercise9	0.634s
```

Two building blocks did all the work:
- `httptest.NewRequest(method, target, body)` builds an `*http.Request` without touching the network
- `httptest.NewRecorder()` builds a `*httptest.ResponseRecorder` that implements `http.ResponseWriter`, capturing the status code and body into fields you can assert on directly (`rec.Code`, `rec.Body`)

For testing a full handler chain (middleware, routers) against a real TCP connection, `httptest.NewServer` spins up a real local server backed by your handler and hands you its URL - `Level 27` will use this extensively.

---

## Mocking via Interfaces

Level 15 established the core idea: an interface describes *behavior*, and any type with the right methods satisfies it - including a type you wrote purely for tests. This is Go's whole approach to "mocking": no framework, no code generation required (though generators exist) - just a small interface and a fake implementation.

Suppose `OrderService` needs to notify someone when an order ships. If it depended directly on a concrete email client, testing it would mean sending real emails. Instead, it depends on a small interface:

```go
// Notifier is the small interface OrderService depends on instead of an
// email/SMS client directly - this is what makes it testable.
type Notifier interface {
    Notify(message string) error
}

type OrderService struct {
    notifier Notifier
}

func NewOrderService(n Notifier) *OrderService {
    return &OrderService{notifier: n}
}

func (s *OrderService) Ship(orderID string) error {
    return s.notifier.Notify("order " + orderID + " has shipped")
}
```

Production code injects a real notifier (an email client, a Slack webhook, whatever satisfies `Notifier`). Tests inject a **fake** that just records what happened:

```go
// fakeNotifier is a test double: instead of sending a real notification,
// it just records what it was asked to send.
type fakeNotifier struct {
    messages []string
}

func (f *fakeNotifier) Notify(message string) error {
    f.messages = append(f.messages, message)
    return nil
}

func TestOrderServiceShip(t *testing.T) {
    fake := &fakeNotifier{}
    svc := NewOrderService(fake)

    if err := svc.Ship("A100"); err != nil {
        t.Fatalf("Ship returned unexpected error: %v", err)
    }

    if len(fake.messages) != 1 {
        t.Fatalf("expected 1 notification, got %d", len(fake.messages))
    }
    want := "order A100 has shipped"
    if fake.messages[0] != want {
        t.Errorf("notification = %q; want %q", fake.messages[0], want)
    }
}
```

```
=== RUN   TestOrderServiceShip
--- PASS: TestOrderServiceShip (0.00s)
PASS
ok  	level26.example/readmemock	0.581s
```

No network call, no real notification, no flakiness from an external system being down - just a plain Go value asserting on what it recorded. This same pattern (define the interface where it's *used*, per Level 15's Rule 3) is how Go code tests database access, clock-dependent logic, external APIs, and file systems, without ever touching the real thing. The comprehensive exercise at the end of this level's `EXERCISES.md` builds a second example of this with a `Clock` interface.

---

## Best Practices

### 1. Table-Driven Tests by Default

Reach for a table + `t.Run` even for what looks like a "simple" function - it costs almost nothing when you have one case, and it scales cleanly the moment you think of a second.

```go
// ✅ Good - scales to any number of cases, each named and independent
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) { ... })
}

// ❌ Avoid - a wall of separate, near-identical test functions
func TestAddPositive(t *testing.T) { ... }
func TestAddNegative(t *testing.T) { ... }
func TestAddZero(t *testing.T) { ... }
```

### 2. Test Behavior, Not Implementation Details

Assert on what a function *returns* or *does* (its public contract), not on internal fields or private helper functions it happens to call today. Implementation-detail tests break every time you refactor, even when behavior hasn't changed.

### 3. One Logical Assertion Focus Per Test

A test named `TestParseConfig` should be about parsing a config - if it also asserts on logging behavior and network retries, split those into their own tests. A failing test name should already tell you roughly what broke.

### 4. Descriptive Test and Subtest Names

`TestIsValidUsername/contains_space` tells you exactly what failed at a glance. `TestIsValidUsername/case_3` does not - you'd have to open the file and count.

### 5. Call t.Helper() in Every Test Helper

Any function taking `*testing.T` that can fail the test should start with `t.Helper()` - see [Test Helpers](#test-helpers-and-thelper) for what goes wrong without it.

---

## Common Mistakes

### Mistake 1: Not Checking Errors in Tests

```go
// ❌ WRONG - if ParseConfig fails, cfg is likely nil/zero and the next
// line either panics or silently passes against garbage data
cfg, _ := ParseConfig(input)
if cfg.Timeout != 30 {
    t.Errorf("Timeout = %d; want 30", cfg.Timeout)
}

// ✅ RIGHT - fail loudly and stop before using an invalid result
cfg, err := ParseConfig(input)
if err != nil {
    t.Fatalf("ParseConfig returned error: %v", err)
}
if cfg.Timeout != 30 {
    t.Errorf("Timeout = %d; want 30", cfg.Timeout)
}
```

### Mistake 2: Forgetting What t.Parallel() Implies About Shared State

Calling `t.Parallel()` inside a test tells Go it may run concurrently with other parallel tests in the same package. That's a problem if those tests share mutable state (a package-level variable, a shared file, a shared map) without synchronization - you get a data race instead of a faster test suite.

```go
// ❌ RISKY - if other tests also call t.Parallel() and touch sharedCounter,
// this is now a data race
func TestIncrement(t *testing.T) {
    t.Parallel()
    sharedCounter++
    if sharedCounter != 1 {
        t.Errorf("sharedCounter = %d; want 1", sharedCounter)
    }
}

// ✅ RIGHT - give each parallel test its own isolated state
func TestIncrement(t *testing.T) {
    t.Parallel()
    counter := &Counter{}
    counter.Increment()
    if counter.Value() != 1 {
        t.Errorf("counter = %d; want 1", counter.Value())
    }
}
```

Also remember: inside a table-driven loop, capture the loop variable per-iteration before calling `t.Parallel()` in the subtest closure, or every parallel subtest can end up sharing the same final `tt` value (this was a common footgun before Go 1.22 changed loop variable semantics; it's worth understanding even though newer Go versions no longer require the manual `tt := tt` workaround).

### Mistake 3: Testing Private Implementation Details Instead of Public Behavior

```go
// ❌ WRONG - couples the test to an internal helper that could be renamed,
// inlined, or removed without changing behavior at all
func TestInternalNormalizeHelper(t *testing.T) {
    got := normalizeWhitespace("  a  b ") // unexported helper
    if got != "a b" {
        t.Errorf("got %q", got)
    }
}

// ✅ RIGHT - test the exported function that uses it; internals can change freely
func TestFormatMessage(t *testing.T) {
    got := FormatMessage("  a  b ")
    if got != "a b" {
        t.Errorf("got %q", got)
    }
}
```

### Mistake 4: Table-Driven Tests Without Subtests

```go
// ❌ WRONG - if case 3 of 10 fails, the failure message is your only clue;
// there's no name, and a t.Fatalf here would skip cases 4-10 entirely
for _, tt := range tests {
    got := Max(tt.input)
    if got != tt.want {
        t.Errorf("Max(%v) = %d; want %d", tt.input, got, tt.want)
    }
}

// ✅ RIGHT - t.Run gives every case a name and independent pass/fail
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        got := Max(tt.input)
        if got != tt.want {
            t.Errorf("Max(%v) = %d; want %d", tt.input, got, tt.want)
        }
    })
}
```

---

## Summary

**The Basics:**
- Tests live in `_test.go` files, next to the code they test, same package
- `func TestXxx(t *testing.T)` - must start with `Test`, capital after it, takes `*testing.T`
- `go test`, `go test -v`, `go test -run <pattern>`

**Failure Reporting:**
- `t.Error`/`t.Errorf` - fail, keep running (see every problem)
- `t.Fatal`/`t.Fatalf` - fail, stop this test now (avoid using invalid state)

**The Idiomatic Shape:**
- Table-driven tests (a slice of `{name, input, want}` structs) run through `t.Run` subtests
- `-run 'TestName/subtest_name'` runs exactly one case

**Measuring Quality and Speed:**
- `go test -cover`, `go tool cover -func`/`-html` - which statements ran
- `func BenchmarkXxx(b *testing.B)`, `go test -bench=.` - `ns/op`, machine-dependent

**Keeping Suites Maintainable:**
- `t.Helper()` in every helper - correct failure line numbers
- `TestMain(m *testing.M)` for package-wide setup/teardown; `t.Cleanup()` for per-test
- `httptest.NewRecorder`/`NewRequest` - test handlers without a real server
- Small interfaces + fake implementations - test without real dependencies

---

## Next Steps

You now understand:
- ✅ How to write, name, and run Go tests correctly
- ✅ When to use t.Error vs t.Fatal
- ✅ Table-driven tests and subtests - the idiomatic Go testing shape
- ✅ How to measure coverage and benchmark performance
- ✅ How to keep large test suites maintainable with helpers, setup/teardown, and mocking

**Next level:** Level 27 - HTTP & REST APIs
- Building real HTTP servers and handlers
- Routing, middleware, and request/response handling
- Testing HTTP services in depth (building on this level's httptest preview)

You can now prove your code works - and keep proving it as it grows! 🚀
