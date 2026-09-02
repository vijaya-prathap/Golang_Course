# Level 26: Testing - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

**A note on this level's format:** Go tests live in `_test.go` files alongside the code they test, not in a single `main.go`. So each exercise below creates **two** files - the code file and its test file - then runs `go test` (instead of `go run main.go`) and shows you the real captured test output.

---

## Exercise 1: A Basic Test Function

**Objective:** Write your first `TestXxx` function against a simple pure function

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level26-exercise1
cd ~/projects/level26-exercise1
go mod init level26.example/exercise1
```

2. Create `math.go`:

```bash
cat > math.go << 'EOF'
package mathutil

// Add returns the sum of a and b.
func Add(a, b int) int {
    return a + b
}
EOF
```

3. Create `math_test.go`:

```bash
cat > math_test.go << 'EOF'
package mathutil

import "testing"

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Add(2, 3) = %d; want 5", result)
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
=== RUN   TestAdd
--- PASS: TestAdd (0.00s)
PASS
ok  	level26.example/exercise1	0.547s
```

*(The `0.547s` total is machine-dependent - yours will differ slightly.)*

**Learning Objectives:**
- ✅ Name a test function so `go test` recognizes it (`Test` + capital letter, `*testing.T` parameter)
- ✅ Put test code in a `_test.go` file, same package, same directory as the code under test
- ✅ Read a passing `go test -v` result

---

## Exercise 2: t.Error vs t.Fatal

**Objective:** See, with real output, exactly how t.Error differs from t.Fatal

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level26-exercise2
cd ~/projects/level26-exercise2
go mod init level26.example/exercise2
```

2. Create `divide.go`:

```bash
cat > divide.go << 'EOF'
package divide

import "errors"

// Divide returns a / b, or an error if b is zero.
func Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}
EOF
```

3. Create `divide_test.go` - both tests below are **deliberately seeded to fail**, so you can see the contrast in behavior:

```bash
cat > divide_test.go << 'EOF'
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
EOF
```

4. Run the tests:

```bash
go test -v ./...
```

**Expected Output:**

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

*(`go test` exits with a non-zero status here - that's correct, since both tests genuinely fail on purpose to demonstrate the difference.)*

**Learning Objectives:**
- ✅ See `t.Error`/`t.Errorf` report a failure and keep executing the rest of the test
- ✅ See `t.Fatal`/`t.Fatalf` report a failure and stop the test immediately (the trailing `t.Log` never runs)
- ✅ Understand why you'd choose one over the other (continuing safely vs. avoiding invalid state)

---

## Exercise 3: Table-Driven Tests Catch a Real Bug

**Objective:** Use a table-driven test to catch a genuine, deliberately-seeded bug

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level26-exercise3
cd ~/projects/level26-exercise3
go mod init level26.example/exercise3
```

2. Create `maxfinder.go` - **this version has a real bug**:

```bash
cat > maxfinder.go << 'EOF'
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
EOF
```

3. Create `maxfinder_test.go`:

```bash
cat > maxfinder_test.go << 'EOF'
package maxfinder

import "testing"

func TestMax(t *testing.T) {
    tests := []struct {
        name  string
        input []int
        want  int
    }{
        {"all positive", []int{3, 1, 4, 1, 5, 9}, 9},
        {"single element", []int{7}, 7},
        {"all negative", []int{-5, -2, -9, -1}, -1},
        {"mixed with zero", []int{-3, 0, -7}, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Max(tt.input)
            if got != tt.want {
                t.Errorf("Max(%v) = %d; want %d", tt.input, got, tt.want)
            }
        })
    }
}
EOF
```

4. Run the tests against the buggy code:

```bash
go test -v ./...
```

**Expected Output (buggy version - the bug is real, so this is a genuine failure):**

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

5. Now fix the bug - update `maxfinder.go`:

```bash
cat > maxfinder.go << 'EOF'
package maxfinder

// Max returns the largest value in nums.
func Max(nums []int) int {
    max := nums[0] // FIXED: start from the first element, not 0
    for _, n := range nums {
        if n > max {
            max = n
        }
    }
    return max
}
EOF
```

6. Run the tests again:

```bash
go test -v ./...
```

**Expected Output (fixed version):**

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

**Learning Objectives:**
- ✅ Write a table of `{name, input, want}` cases and loop over it
- ✅ See a table-driven test catch a real edge-case bug (`max := 0` fails for all-negative input)
- ✅ Fix the bug and confirm the same table now passes

---

## Exercise 4: Subtests and Running One With -run

**Objective:** Use t.Run to name each case, then run a single subtest by name

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level26-exercise4
cd ~/projects/level26-exercise4
go mod init level26.example/exercise4
```

2. Create `validate.go`:

```bash
cat > validate.go << 'EOF'
package validate

import "strings"

// IsValidUsername reports whether name is a valid username: 3-12 characters,
// no spaces, and not empty.
func IsValidUsername(name string) bool {
    if len(name) < 3 || len(name) > 12 {
        return false
    }
    if strings.Contains(name, " ") {
        return false
    }
    return true
}
EOF
```

3. Create `validate_test.go`:

```bash
cat > validate_test.go << 'EOF'
package validate

import "testing"

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
EOF
```

4. Run all subtests:

```bash
go test -v ./...
```

**Expected Output:**

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

5. Now run just the `empty_string` case:

```bash
go test -v -run 'TestIsValidUsername/empty_string' ./...
```

**Expected Output:**

```
=== RUN   TestIsValidUsername
=== RUN   TestIsValidUsername/empty_string
--- PASS: TestIsValidUsername (0.00s)
    --- PASS: TestIsValidUsername/empty_string (0.00s)
PASS
ok  	level26.example/exercise4	0.248s
```

**Learning Objectives:**
- ✅ See how `t.Run` names show up in `-v` output as `TestName/subtest_name`
- ✅ Use `-run 'TestName/subtest_name'` to execute exactly one subtest
- ✅ Understand why this matters when debugging one failing case in a large suite

---

## Exercise 5: Measuring Test Coverage

**Objective:** Use go test -cover and go tool cover to see what your tests actually exercise

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level26-exercise5
cd ~/projects/level26-exercise5
go mod init level26.example/exercise5
```

2. Create `grade.go`:

```bash
cat > grade.go << 'EOF'
package grade

// Letter converts a numeric score (0-100) into a letter grade.
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
EOF
```

3. Create `grade_test.go` - **deliberately covering only some branches** (`A`, `B`, and `F`, skipping `C` and `D`):

```bash
cat > grade_test.go << 'EOF'
package grade

import "testing"

func TestLetter(t *testing.T) {
    tests := []struct {
        name  string
        score int
        want  string
    }{
        {"A grade", 95, "A"},
        {"B grade", 85, "B"},
        {"F grade", 40, "F"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Letter(tt.score)
            if got != tt.want {
                t.Errorf("Letter(%d) = %q; want %q", tt.score, got, tt.want)
            }
        })
    }
}
EOF
```

4. Run with coverage:

```bash
go test -cover ./...
```

**Expected Output:**

```
ok  	level26.example/exercise5	0.896s	coverage: 66.7% of statements
```

5. Generate a coverage profile and inspect it per-function:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

**Expected Output:**

```
level26.example/exercise5/grade.go:4:	Letter		66.7%
total:					(statements)	66.7%
```

6. Generate the HTML report (open `coverage.html` in a browser to see covered/uncovered lines highlighted):

```bash
go tool cover -html=coverage.out -o coverage.html
```

**Learning Objectives:**
- ✅ Read a `go test -cover` percentage and understand what "statements" means
- ✅ Use `go tool cover -func` to see coverage broken down per function
- ✅ Generate an HTML coverage report with `go tool cover -html`
- ✅ Recognize that 66.7% here means the `C` and `D` branches were never exercised

---

## Exercise 6: Writing and Running a Benchmark

**Objective:** Write a BenchmarkXxx function and read real ns/op output

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level26-exercise6
cd ~/projects/level26-exercise6
go mod init level26.example/exercise6
```

2. Create `reverse.go`:

```bash
cat > reverse.go << 'EOF'
package strutil

// Reverse returns s with its runes in reverse order.
func Reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}
EOF
```

3. Create `reverse_test.go`:

```bash
cat > reverse_test.go << 'EOF'
package strutil

import "testing"

func TestReverse(t *testing.T) {
    got := Reverse("hello")
    want := "olleh"
    if got != want {
        t.Errorf("Reverse(\"hello\") = %q; want %q", got, want)
    }
}

func BenchmarkReverse(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Reverse("The quick brown fox jumps over the lazy dog")
    }
}
EOF
```

4. Run the regular test first:

```bash
go test -v ./...
```

**Expected Output:**

```
=== RUN   TestReverse
--- PASS: TestReverse (0.00s)
PASS
ok  	level26.example/exercise6	0.713s
```

5. Run the benchmark (`-run=^$` skips regular tests so only the benchmark runs):

```bash
go test -bench=. -run=^$ ./...
```

**Expected Output** *(your `ns/op` and iteration count WILL differ - these numbers are entirely dependent on your CPU, load, and Go version; the values below were captured on one specific machine)*:

```
goos: darwin
goarch: arm64
pkg: level26.example/exercise6
cpu: Apple M1 Pro
BenchmarkReverse-10    	 3444912	       332.6 ns/op
PASS
ok  	level26.example/exercise6	1.746s
```

6. Add `-benchmem` to also see allocations:

```bash
go test -bench=. -run=^$ -benchmem ./...
```

**Expected Output** *(again, machine-dependent)*:

```
goos: darwin
goarch: arm64
pkg: level26.example/exercise6
cpu: Apple M1 Pro
BenchmarkReverse-10    	 3508375	       333.8 ns/op	     224 B/op	       2 allocs/op
PASS
ok  	level26.example/exercise6	1.757s
```

**Learning Objectives:**
- ✅ Write a `func BenchmarkXxx(b *testing.B)` using the `b.N` loop
- ✅ Run benchmarks with `go test -bench=.` and read `ns/op`
- ✅ Understand that benchmark numbers vary by machine and are for relative, not absolute, comparison
- ✅ Use `-benchmem` to see bytes and allocations per operation

---

## Exercise 7: A Test Helper With t.Helper()

**Objective:** Extract a test helper and see t.Helper() fix failure line-number reporting

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level26-exercise7
cd ~/projects/level26-exercise7
go mod init level26.example/exercise7
```

2. Create `account.go`:

```bash
cat > account.go << 'EOF'
package account

// Balance represents a simple account balance in cents.
type Balance struct {
    cents int
}

// Deposit adds amount cents to the balance.
func (b *Balance) Deposit(amount int) {
    b.cents += amount
}

// Cents returns the current balance in cents.
func (b *Balance) Cents() int {
    return b.cents
}
EOF
```

3. Create `account_test.go` - the second test has a **deliberately wrong expectation**, so you can see exactly where the failure gets reported:

```bash
cat > account_test.go << 'EOF'
package account

import "testing"

// assertCents is a test helper. t.Helper() tells the testing package that
// failures here should be blamed on the CALLER's line, not this line.
func assertCents(t *testing.T, b *Balance, want int) {
    t.Helper()
    if got := b.Cents(); got != want {
        t.Errorf("balance = %d cents; want %d cents", got, want)
    }
}

func TestDeposit(t *testing.T) {
    b := &Balance{}
    b.Deposit(500)
    assertCents(t, b, 500)
}

func TestDepositWrongExpectation(t *testing.T) {
    b := &Balance{}
    b.Deposit(500)
    assertCents(t, b, 999) // deliberately wrong, to see WHERE the failure is reported
}
EOF
```

4. Run the tests:

```bash
go test -v ./...
```

**Expected Output** (with `t.Helper()` in place - the failure points at line 23, the call site inside `TestDepositWrongExpectation`):

```
=== RUN   TestDeposit
--- PASS: TestDeposit (0.00s)
=== RUN   TestDepositWrongExpectation
    account_test.go:23: balance = 500 cents; want 999 cents
--- FAIL: TestDepositWrongExpectation (0.00s)
FAIL
FAIL	level26.example/exercise7	0.538s
FAIL
```

5. Now comment out (or delete) the `t.Helper()` line inside `assertCents` and run the tests again. **Expected Output** (same bug, but now blamed on line 8 - inside the helper itself, which is far less useful when many tests share this helper):

```
=== RUN   TestDeposit
--- PASS: TestDeposit (0.00s)
=== RUN   TestDepositWrongExpectation
    account_test.go:8: balance = 500 cents; want 999 cents
--- FAIL: TestDepositWrongExpectation (0.00s)
FAIL
FAIL	level26.example/exercise7	0.649s
FAIL
```

**Learning Objectives:**
- ✅ Extract a repeated assertion into a helper function taking `*testing.T`
- ✅ Call `t.Helper()` as the helper's first line
- ✅ See, with a real before/after comparison, exactly which source line changes in the failure output

---

## Exercise 8: Setup and Teardown (TestMain and t.Cleanup)

**Objective:** Use TestMain for package-level setup/teardown, and t.Cleanup for per-test cleanup

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level26-exercise8
cd ~/projects/level26-exercise8
go mod init level26.example/exercise8
```

2. Create `store.go`:

```bash
cat > store.go << 'EOF'
package store

// Store is a tiny in-memory key-value store.
type Store struct {
    data map[string]int
}

// NewStore creates an empty Store.
func NewStore() *Store {
    return &Store{data: make(map[string]int)}
}

// Set stores value under key.
func (s *Store) Set(key string, value int) {
    s.data[key] = value
}

// Get returns the value stored under key, and whether it was found.
func (s *Store) Get(key string) (int, bool) {
    v, ok := s.data[key]
    return v, ok
}
EOF
```

3. Create `store_test.go`:

```bash
cat > store_test.go << 'EOF'
package store

import (
    "fmt"
    "os"
    "testing"
)

// sharedStore is set up once, before ANY test in this package runs, and
// torn down once, after ALL of them finish.
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

// newTempStore demonstrates PER-TEST setup/teardown with t.Cleanup: each
// test gets its own fresh Store, and the cleanup runs when that single
// test finishes - not once for the whole package.
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
EOF
```

4. Run the tests:

```bash
go test -v ./...
```

**Expected Output:**

```
=== SETUP: seeding shared store ===
=== RUN   TestLookupAlice
--- PASS: TestLookupAlice (0.00s)
=== RUN   TestLookupBob
--- PASS: TestLookupBob (0.00s)
=== RUN   TestTempStoreIsIsolated
    store_test.go:48: cleanup: discarding temp store for TestTempStoreIsIsolated
--- PASS: TestTempStoreIsIsolated (0.00s)
PASS
=== TEARDOWN: clearing shared store ===
ok  	level26.example/exercise8	0.536s
```

**Learning Objectives:**
- ✅ Use `TestMain(m *testing.M)` to run setup/teardown exactly once per package
- ✅ See that `SETUP`/`TEARDOWN` print once, regardless of how many tests run
- ✅ Register a per-test cleanup with `t.Cleanup` from inside a helper
- ✅ Understand `t.Cleanup` runs automatically when that specific test finishes

---

## Exercise 9: Testing an HTTP Handler With httptest

**Objective:** Test an HTTP handler without starting a real server

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level26-exercise9
cd ~/projects/level26-exercise9
go mod init level26.example/exercise9
```

2. Create `handler.go`:

```bash
cat > handler.go << 'EOF'
package api

import (
    "fmt"
    "net/http"
)

// GreetHandler writes a greeting for the "name" query parameter, defaulting
// to "World" when it's missing.
func GreetHandler(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    if name == "" {
        name = "World"
    }
    w.Header().Set("Content-Type", "text/plain")
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Hello, %s!", name)
}
EOF
```

3. Create `handler_test.go`:

```bash
cat > handler_test.go << 'EOF'
package api

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

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
EOF
```

4. Run the tests:

```bash
go test -v ./...
```

**Expected Output:**

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

**Learning Objectives:**
- ✅ Build a fake request with `httptest.NewRequest`
- ✅ Capture a handler's response with `httptest.NewRecorder`
- ✅ Assert on `rec.Code` and `rec.Body` without a real network connection
- ✅ Combine table-driven tests with HTTP handler testing

---

## Exercise 10: Comprehensive Practice — A Mockable Discount Service

**Objective:** Combine interfaces, fakes/mocks, and table-driven tests in one realistic example

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level26-exercise10
cd ~/projects/level26-exercise10
go mod init level26.example/exercise10
```

2. Create `clock.go`:

```bash
cat > clock.go << 'EOF'
package pricing

import "time"

// Clock is the small interface the pricing logic depends on instead of
// calling time.Now() directly - this is what makes it testable.
type Clock interface {
    Now() time.Time
}

// RealClock is the production implementation: it defers to time.Now().
type RealClock struct{}

func (RealClock) Now() time.Time {
    return time.Now()
}

// FakeClock is a test double: it always returns a fixed time, so tests
// never depend on what day it happens to be when they run.
type FakeClock struct {
    FixedTime time.Time
}

func (f FakeClock) Now() time.Time {
    return f.FixedTime
}
EOF
```

3. Create `pricing.go`:

```bash
cat > pricing.go << 'EOF'
package pricing

import "time"

// DiscountService calculates discounts. It depends on the Clock interface,
// not on time.Now() directly, so tests can inject a FakeClock.
type DiscountService struct {
    clock Clock
}

// NewDiscountService builds a DiscountService using the given Clock.
func NewDiscountService(clock Clock) *DiscountService {
    return &DiscountService{clock: clock}
}

// CalculateDiscount returns the discount percentage for a purchase of the
// given amount:
//   - amount >= 100:  10% base discount
//   - amount >= 50:   5% base discount
//   - otherwise:      0% base discount
//   - weekends (Saturday/Sunday) add an extra 5% on top of the base
func (s *DiscountService) CalculateDiscount(amount float64) int {
    discount := 0
    switch {
    case amount >= 100:
        discount = 10
    case amount >= 50:
        discount = 5
    }

    weekday := s.clock.Now().Weekday()
    if weekday == time.Saturday || weekday == time.Sunday {
        discount += 5
    }
    return discount
}
EOF
```

4. Create `pricing_test.go`:

```bash
cat > pricing_test.go << 'EOF'
package pricing

import (
    "testing"
    "time"
)

// mustDate is a small test helper for building fixed dates. t.Helper()
// keeps failures pointing at the test case, not this helper.
func mustDate(t *testing.T, layout, value string) time.Time {
    t.Helper()
    d, err := time.Parse(layout, value)
    if err != nil {
        t.Fatalf("mustDate(%q): %v", value, err)
    }
    return d
}

func TestCalculateDiscount(t *testing.T) {
    tests := []struct {
        name   string
        amount float64
        date   string // RFC3339 date
        want   int
    }{
        {"small order on a weekday", 20, "2026-08-31T10:00:00Z", 0},   // Monday
        {"mid order on a weekday", 60, "2026-08-31T10:00:00Z", 5},     // Monday
        {"large order on a weekday", 150, "2026-08-31T10:00:00Z", 10}, // Monday
        {"small order on a weekend", 20, "2026-09-05T10:00:00Z", 5},   // Saturday
        {"large order on a weekend", 150, "2026-09-06T10:00:00Z", 15}, // Sunday
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            fake := FakeClock{FixedTime: mustDate(t, time.RFC3339, tt.date)}
            svc := NewDiscountService(fake)

            got := svc.CalculateDiscount(tt.amount)
            if got != tt.want {
                t.Errorf("CalculateDiscount(%.2f) on %s = %d; want %d", tt.amount, tt.date, got, tt.want)
            }
        })
    }
}

// TestRealClockIsWired confirms RealClock actually satisfies Clock and
// returns a sane (non-zero) time - it does NOT assert on a specific
// discount, since "today" changes every day this test runs.
func TestRealClockIsWired(t *testing.T) {
    svc := NewDiscountService(RealClock{})
    discount := svc.CalculateDiscount(150)
    if discount < 10 || discount > 15 {
        t.Errorf("CalculateDiscount(150) with RealClock = %d; want between 10 and 15", discount)
    }
}
EOF
```

5. Run the tests:

```bash
go test -v ./...
```

**Expected Output:**

```
=== RUN   TestCalculateDiscount
=== RUN   TestCalculateDiscount/small_order_on_a_weekday
=== RUN   TestCalculateDiscount/mid_order_on_a_weekday
=== RUN   TestCalculateDiscount/large_order_on_a_weekday
=== RUN   TestCalculateDiscount/small_order_on_a_weekend
=== RUN   TestCalculateDiscount/large_order_on_a_weekend
--- PASS: TestCalculateDiscount (0.00s)
    --- PASS: TestCalculateDiscount/small_order_on_a_weekday (0.00s)
    --- PASS: TestCalculateDiscount/mid_order_on_a_weekday (0.00s)
    --- PASS: TestCalculateDiscount/large_order_on_a_weekday (0.00s)
    --- PASS: TestCalculateDiscount/small_order_on_a_weekend (0.00s)
    --- PASS: TestCalculateDiscount/large_order_on_a_weekend (0.00s)
=== RUN   TestRealClockIsWired
--- PASS: TestRealClockIsWired (0.00s)
PASS
ok  	level26.example/exercise10	0.561s
```

6. Confirm full coverage:

```bash
go test -cover ./...
```

**Expected Output:**

```
ok  	level26.example/exercise10	0.554s	coverage: 100.0% of statements
```

**Learning Objectives:**
- ✅ Define a small interface (`Clock`) at the point where it's consumed
- ✅ Write both a real implementation (`RealClock`) and a fake/test implementation (`FakeClock`)
- ✅ Inject the fake into the service under test to make date-dependent logic fully deterministic
- ✅ Combine everything from this level: table-driven tests, subtests, a test helper with `t.Helper()`, and mocking via interfaces

---

## Bonus Challenges

### Challenge 1: Testing a Level 25 Generic Function

Write table-driven tests for generic `Map` and `Filter` functions (a preview tying back to Level 25: Generics).

```bash
mkdir -p ~/projects/level26-bonus1
cd ~/projects/level26-bonus1
go mod init level26.example/bonus1
```

```bash
cat > generic.go << 'EOF'
package slices2

// Map applies f to every element of s and returns a new slice of the results.
// (A preview of Level 25: Generics - T and U can be any types.)
func Map[T, U any](s []T, f func(T) U) []U {
    result := make([]U, len(s))
    for i, v := range s {
        result[i] = f(v)
    }
    return result
}

// Filter returns a new slice containing only the elements of s for which
// keep returns true.
func Filter[T any](s []T, keep func(T) bool) []T {
    var result []T
    for _, v := range s {
        if keep(v) {
            result = append(result, v)
        }
    }
    return result
}
EOF
```

```bash
cat > generic_test.go << 'EOF'
package slices2

import (
    "reflect"
    "testing"
)

func TestMap(t *testing.T) {
    nums := []int{1, 2, 3, 4}
    got := Map(nums, func(n int) int { return n * n })
    want := []int{1, 4, 9, 16}
    if !reflect.DeepEqual(got, want) {
        t.Errorf("Map(squares) = %v; want %v", got, want)
    }
}

func TestMapDifferentTypes(t *testing.T) {
    nums := []int{1, 2, 3}
    got := Map(nums, func(n int) string {
        if n%2 == 0 {
            return "even"
        }
        return "odd"
    })
    want := []string{"odd", "even", "odd"}
    if !reflect.DeepEqual(got, want) {
        t.Errorf("Map(parity) = %v; want %v", got, want)
    }
}

func TestFilter(t *testing.T) {
    nums := []int{1, 2, 3, 4, 5, 6}
    got := Filter(nums, func(n int) bool { return n%2 == 0 })
    want := []int{2, 4, 6}
    if !reflect.DeepEqual(got, want) {
        t.Errorf("Filter(evens) = %v; want %v", got, want)
    }
}
EOF
```

```bash
go test -v ./...
```

**Real captured output:**

```
=== RUN   TestMap
--- PASS: TestMap (0.00s)
=== RUN   TestMapDifferentTypes
--- PASS: TestMapDifferentTypes (0.00s)
=== RUN   TestFilter
--- PASS: TestFilter (0.00s)
PASS
ok  	level26.example/bonus1	0.534s
```

**Note on generic test functions:** `Map[T, U any]` has two independently-inferred type parameters, so `TestMap` (int → int) and `TestMapDifferentTypes` (int → string) both compile and pass against the *same* generic function - no duplicated logic per type pair.

### Challenge 2: Fuzz Testing

Write `func FuzzXxx(f *testing.F)` for a `Reverse` function, checking the property "reversing twice returns the original string" - and see what the fuzzer actually finds.

```bash
mkdir -p ~/projects/level26-bonus2
cd ~/projects/level26-bonus2
go mod init level26.example/bonus2
```

```bash
cat > reverse.go << 'EOF'
package strutil

// Reverse returns s with its runes in reverse order.
func Reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}
EOF
```

```bash
cat > reverse_test.go << 'EOF'
package strutil

import "testing"

func TestReverse(t *testing.T) {
    got := Reverse("hello")
    want := "olleh"
    if got != want {
        t.Errorf("Reverse(\"hello\") = %q; want %q", got, want)
    }
}

// FuzzReverse checks a PROPERTY instead of a fixed input/output pair:
// reversing a string twice must always return the original string, for
// any input the fuzzer can generate.
func FuzzReverse(f *testing.F) {
    // Seed corpus: known starting inputs the fuzzer mutates from.
    f.Add("hello")
    f.Add("")
    f.Add("héllo")

    f.Fuzz(func(t *testing.T, s string) {
        once := Reverse(s)
        twice := Reverse(once)
        if twice != s {
            t.Errorf("Reverse(Reverse(%q)) = %q; want %q", s, twice, s)
        }
    })
}
EOF
```

First, a normal `go test` only runs the seed corpus (`f.Add` calls) - it does NOT fuzz automatically:

```bash
go test -v ./...
```

**Real captured output:**

```
=== RUN   TestReverse
--- PASS: TestReverse (0.00s)
=== RUN   FuzzReverse
=== RUN   FuzzReverse/seed#0
=== RUN   FuzzReverse/seed#1
=== RUN   FuzzReverse/seed#2
--- PASS: FuzzReverse (0.00s)
    --- PASS: FuzzReverse/seed#0 (0.00s)
    --- PASS: FuzzReverse/seed#1 (0.00s)
    --- PASS: FuzzReverse/seed#2 (0.00s)
PASS
ok  	level26.example/bonus2	0.627s
```

Now actually fuzz it for a few seconds with `-fuzz`:

```bash
go test -fuzz=FuzzReverse -fuzztime=5s ./...
```

**Real captured output - fuzzing genuinely found a bug within a fraction of a second:**

```
fuzz: elapsed: 0s, gathering baseline coverage: 0/3 completed
fuzz: elapsed: 0s, gathering baseline coverage: 3/3 completed, now fuzzing with 10 workers
fuzz: minimizing 36-byte failing input file
fuzz: elapsed: 0s, minimizing
--- FAIL: FuzzReverse (0.03s)
    --- FAIL: FuzzReverse (0.00s)
        reverse_test.go:26: Reverse(Reverse("\xdb")) = "�"; want "\xdb"
    
    Failing input written to testdata/fuzz/FuzzReverse/a2f306083e514b06
    To re-run:
    go test -run=FuzzReverse/a2f306083e514b06
FAIL
exit status 1
FAIL	level26.example/bonus2	0.567s
```

**What actually happened:** this is verified, genuine fuzzing behavior on Go 1.26 - not a hypothetical. `\xdb` alone is not valid UTF-8. Converting it to `[]rune` decodes it as the Unicode replacement character `U+FFFD` (`�`), which is NOT reversible back to the original byte `\xdb`. The fuzzer found this in well under a second because a lone invalid byte is easy to stumble onto randomly. This is precisely what fuzz testing is for: it explores inputs no one thought to write by hand, and it found a real, subtle bug - `Reverse` silently assumes its input is valid UTF-8, which Go strings never guarantee.

The failing input is saved as a permanent regression test:

```bash
cat testdata/fuzz/FuzzReverse/a2f306083e514b06
```

```
go test fuzz v1
string("\xdb")
```

Re-running just that failing case (without fuzzing again):

```bash
go test -run=FuzzReverse/a2f306083e514b06 -v ./...
```

```
=== RUN   FuzzReverse
=== RUN   FuzzReverse/a2f306083e514b06
    reverse_test.go:26: Reverse(Reverse("\xdb")) = "�"; want "\xdb"
--- FAIL: FuzzReverse (0.00s)
    --- FAIL: FuzzReverse/a2f306083e514b06 (0.00s)
FAIL
FAIL	level26.example/bonus2	0.277s
FAIL
```

This file now lives under `testdata/fuzz/` and will run automatically as a regular test case on every future `go test`, even without `-fuzz` - the fuzzer's discovery becomes a permanent part of your regression suite.

### Challenge 3: Golden File Testing

Compare a function's output against a saved "golden" file, with an `-update` flag to regenerate it deliberately.

```bash
mkdir -p ~/projects/level26-bonus3/testdata
cd ~/projects/level26-bonus3
go mod init level26.example/bonus3
```

```bash
cat > report.go << 'EOF'
package report

import (
    "fmt"
    "strings"
)

// Item is a single line item on an invoice.
type Item struct {
    Name  string
    Price float64
}

// Render builds a plain-text invoice report for the given items.
func Render(items []Item) string {
    var b strings.Builder
    b.WriteString("INVOICE\n")
    b.WriteString("=======\n")
    total := 0.0
    for _, it := range items {
        fmt.Fprintf(&b, "%-20s $%.2f\n", it.Name, it.Price)
        total += it.Price
    }
    b.WriteString("-------\n")
    fmt.Fprintf(&b, "%-20s $%.2f\n", "TOTAL", total)
    return b.String()
}
EOF
```

```bash
cat > report_test.go << 'EOF'
package report

import (
    "flag"
    "os"
    "testing"
)

// -update regenerates the golden file instead of comparing against it.
// Run: go test -update
var update = flag.Bool("update", false, "update golden files")

func TestRender(t *testing.T) {
    items := []Item{
        {"Widget", 9.99},
        {"Gadget", 19.99},
        {"Gizmo", 4.50},
    }
    got := Render(items)

    golden := "testdata/invoice.golden"

    if *update {
        if err := os.WriteFile(golden, []byte(got), 0644); err != nil {
            t.Fatalf("failed to update golden file: %v", err)
        }
    }

    want, err := os.ReadFile(golden)
    if err != nil {
        t.Fatalf("failed to read golden file: %v", err)
    }

    if got != string(want) {
        t.Errorf("Render() output does not match golden file %s\n--- got ---\n%s\n--- want ---\n%s", golden, got, string(want))
    }
}
EOF
```

First, generate the golden file:

```bash
go test -v -update ./...
```

**Real captured output:**

```
=== RUN   TestRender
--- PASS: TestRender (0.00s)
PASS
ok  	level26.example/bonus3	0.559s
```

The generated `testdata/invoice.golden`:

```
INVOICE
=======
Widget               $9.99
Gadget               $19.99
Gizmo                $4.50
-------
TOTAL                $34.48
```

Now run normally (no `-update`) - it compares against the saved file instead of regenerating it:

```bash
go test -v ./...
```

**Real captured output:**

```
=== RUN   TestRender
--- PASS: TestRender (0.00s)
PASS
ok  	level26.example/bonus3	0.251s
```

**Learning Objectives (Bonus Challenges):**
- ✅ Test a generic function against multiple concrete type instantiations
- ✅ Write a fuzz test and confirm `-fuzz` genuinely explores random inputs (verified working on Go 1.26.5)
- ✅ See a fuzzer find a real bug (invalid UTF-8 breaks the reverse-twice property) and understand why
- ✅ Implement the golden-file pattern with an `-update` flag to regenerate expected output deliberately

---

## What You've Learned

After completing these 10 exercises and 3 bonus challenges, you can:

✅ Write and name `TestXxx` functions correctly, and run them with `go test`
✅ Choose between `t.Error`/`t.Errorf` and `t.Fatal`/`t.Fatalf` deliberately
✅ Write table-driven tests that catch real bugs, including edge cases
✅ Use `t.Run` subtests and run a single one with `-run`
✅ Measure and interpret test coverage with `go test -cover` and `go tool cover`
✅ Write and read benchmarks, understanding that raw numbers are machine-dependent
✅ Extract test helpers correctly using `t.Helper()`
✅ Set up and tear down test state with `TestMain` and `t.Cleanup`
✅ Test HTTP handlers with `httptest.NewRecorder`/`NewRequest`
✅ Mock dependencies via small interfaces and fake implementations
✅ (Bonus) Test generics, write a fuzz test, and use golden files

---

## Next Level

Level 27: HTTP & REST APIs
- Building real HTTP servers and handlers
- Routing, middleware, and request/response handling
- Testing HTTP services in depth, building on this level's httptest preview

Great work! You can now prove your code works! 🚀
