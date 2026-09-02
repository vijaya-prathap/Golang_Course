# Level 16: Error Handling - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

---

## Exercise 1: Errors Are Values - The if err != nil Flow

**Objective:** Create errors with errors.New and fmt.Errorf and handle them with the fundamental pattern

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level16-exercise1
cd ~/projects/level16-exercise1
go mod init level16.example/exercise1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "errors"
    "fmt"
)

func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

func findUser(id int) (string, error) {
    users := map[int]string{1: "Alice", 2: "Bob"}
    name, ok := users[id]
    if !ok {
        return "", fmt.Errorf("user %d not found", id)
    }
    return name, nil
}

func main() {
    fmt.Println("=== The if err != nil Flow ===")
    result, err := divide(10, 2)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("10 / 2 =", result)
    }

    result, err = divide(10, 0)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("10 / 0 =", result)
    }

    fmt.Println("\n=== fmt.Errorf Builds Formatted Errors ===")
    for _, id := range []int{1, 42} {
        name, err := findUser(id)
        if err != nil {
            fmt.Println("Error:", err)
            continue
        }
        fmt.Println("Found:", name)
    }

    fmt.Println("\n=== An error Is Just a Value ===")
    e := errors.New("something went wrong")
    fmt.Printf("Type: %T\n", e)
    fmt.Println("Message:", e.Error())
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== The if err != nil Flow ===
10 / 2 = 5
Error: division by zero

=== fmt.Errorf Builds Formatted Errors ===
Found: Alice
Error: user 42 not found

=== An error Is Just a Value ===
Type: *errors.errorString
Message: something went wrong
```

**Learning Objectives:**
- ✅ Create errors with errors.New and fmt.Errorf
- ✅ Handle failure and success paths with if err != nil
- ✅ See that an error is an ordinary value with a concrete type

---

## Exercise 2: Adding Context Through Two Layers

**Objective:** Build an error message that tells the whole story, layer by layer

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level16-exercise2
cd ~/projects/level16-exercise2
go mod init level16.example/exercise2
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "strconv"
)

func parsePort(raw string) (int, error) {
    port, err := strconv.Atoi(raw)
    if err != nil {
        return 0, fmt.Errorf("parsing port %q: %w", raw, err)
    }
    if port < 1 || port > 65535 {
        return 0, fmt.Errorf("port %d out of range (1-65535)", port)
    }
    return port, nil
}

func loadConfig(rawPort string) (int, error) {
    port, err := parsePort(rawPort)
    if err != nil {
        return 0, fmt.Errorf("loading config: %w", err)
    }
    return port, nil
}

func main() {
    fmt.Println("=== Success Path ===")
    port, err := loadConfig("8080")
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Server port:", port)
    }

    fmt.Println("\n=== Bad Number: Two Layers of Context ===")
    _, err = loadConfig("abc")
    if err != nil {
        fmt.Println("Error:", err)
    }

    fmt.Println("\n=== Out of Range: Context From One Layer ===")
    _, err = loadConfig("99999")
    if err != nil {
        fmt.Println("Error:", err)
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
=== Success Path ===
Server port: 8080

=== Bad Number: Two Layers of Context ===
Error: loading config: parsing port "abc": strconv.Atoi: parsing "abc": invalid syntax

=== Out of Range: Context From One Layer ===
Error: loading config: port 99999 out of range (1-65535)
```

**Learning Objectives:**
- ✅ Return early on error at every layer
- ✅ Add one clause of context per layer with fmt.Errorf
- ✅ Read a layered error message from outermost operation to root cause

---

## Exercise 3: Wrapping With %w and errors.Unwrap

**Objective:** Build an error chain and walk it layer by layer

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level16-exercise3
cd ~/projects/level16-exercise3
go mod init level16.example/exercise3
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "errors"
    "fmt"
)

func main() {
    base := errors.New("connection refused")
    mid := fmt.Errorf("dialing database: %w", base)
    top := fmt.Errorf("starting server: %w", mid)

    fmt.Println("=== The Three Errors ===")
    fmt.Println("base:", base)
    fmt.Println("mid: ", mid)
    fmt.Println("top: ", top)

    fmt.Println("\n=== Walking the Chain With errors.Unwrap ===")
    for err := error(top); err != nil; err = errors.Unwrap(err) {
        fmt.Printf("-> %v\n", err)
    }

    fmt.Println("\n=== Unwrapping Past the End Gives nil ===")
    fmt.Println(errors.Unwrap(base))

    fmt.Println("\n=== %v vs %w: Only %w Builds a Chain ===")
    plain := fmt.Errorf("dialing database: %v", base)
    fmt.Println("wrapped with %v, Unwrap gives:", errors.Unwrap(plain))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== The Three Errors ===
base: connection refused
mid:  dialing database: connection refused
top:  starting server: dialing database: connection refused

=== Walking the Chain With errors.Unwrap ===
-> starting server: dialing database: connection refused
-> dialing database: connection refused
-> connection refused

=== Unwrapping Past the End Gives nil ===
<nil>

=== %v vs %w: Only %w Builds a Chain ===
wrapped with %v, Unwrap gives: <nil>
```

**Learning Objectives:**
- ✅ Build a multi-layer error chain with %w
- ✅ Walk a chain with errors.Unwrap until it returns nil
- ✅ Prove that %v copies text but does NOT build a chain

---

## Exercise 4: errors.Is - A Sentinel Through a Wrapped Chain

**Objective:** Define a sentinel error and match it through two layers of wrapping

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level16-exercise4
cd ~/projects/level16-exercise4
go mod init level16.example/exercise4
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "errors"
    "fmt"
)

var ErrNotFound = errors.New("not found")

func findOrder(id int) error {
    if id != 1 {
        return fmt.Errorf("order %d: %w", id, ErrNotFound)
    }
    return nil
}

func loadInvoice(orderID int) error {
    if err := findOrder(orderID); err != nil {
        return fmt.Errorf("loading invoice: %w", err)
    }
    return nil
}

func main() {
    fmt.Println("=== errors.Is Sees Through the Chain ===")
    err := loadInvoice(42)
    fmt.Println("full error:", err)
    fmt.Println("errors.Is(err, ErrNotFound):", errors.Is(err, ErrNotFound))

    fmt.Println("\n=== == Only Matches the Outermost Error ===")
    fmt.Println("err == ErrNotFound:", err == ErrNotFound)

    fmt.Println("\n=== No Match for a Different Sentinel ===")
    ErrTimeout := errors.New("timeout")
    fmt.Println("errors.Is(err, ErrTimeout):", errors.Is(err, ErrTimeout))

    fmt.Println("\n=== The Success Path Has No Error ===")
    err = loadInvoice(1)
    fmt.Println("err:", err)
    fmt.Println("errors.Is(nil, ErrNotFound):", errors.Is(err, ErrNotFound))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== errors.Is Sees Through the Chain ===
full error: loading invoice: order 42: not found
errors.Is(err, ErrNotFound): true

=== == Only Matches the Outermost Error ===
err == ErrNotFound: false

=== No Match for a Different Sentinel ===
errors.Is(err, ErrTimeout): false

=== The Success Path Has No Error ===
err: <nil>
errors.Is(nil, ErrNotFound): false
```

**Learning Objectives:**
- ✅ Define a sentinel error with the Err- naming convention
- ✅ Match a sentinel through two layers of wrapping with errors.Is
- ✅ See why == fails on wrapped errors

---

## Exercise 5: errors.As - Extracting a Custom Error Type

**Objective:** Pull a typed error out of a chain and read its fields

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level16-exercise5
cd ~/projects/level16-exercise5
go mod init level16.example/exercise5
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "errors"
    "fmt"
)

type HTTPError struct {
    StatusCode int
    URL        string
}

func (e *HTTPError) Error() string {
    return fmt.Sprintf("http %d fetching %s", e.StatusCode, e.URL)
}

func fetch(url string) error {
    if url == "https://example.com/missing" {
        return &HTTPError{StatusCode: 404, URL: url}
    }
    return nil
}

func loadProfile(url string) error {
    if err := fetch(url); err != nil {
        return fmt.Errorf("loading profile: %w", err)
    }
    return nil
}

func main() {
    fmt.Println("=== errors.As Extracts the Typed Error ===")
    err := loadProfile("https://example.com/missing")
    fmt.Println("full error:", err)

    var httpErr *HTTPError
    if errors.As(err, &httpErr) {
        fmt.Println("status code:", httpErr.StatusCode)
        fmt.Println("url:       ", httpErr.URL)
        if httpErr.StatusCode == 404 {
            fmt.Println("decision:   show the 'page not found' screen")
        }
    }

    fmt.Println("\n=== errors.As Reports false When the Type Is Absent ===")
    plain := errors.New("no HTTPError anywhere in here")
    var target *HTTPError
    fmt.Println("errors.As(plain, &target):", errors.As(plain, &target))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== errors.As Extracts the Typed Error ===
full error: loading profile: http 404 fetching https://example.com/missing
status code: 404
url:        https://example.com/missing
decision:   show the 'page not found' screen

=== errors.As Reports false When the Type Is Absent ===
errors.As(plain, &target): false
```

**Learning Objectives:**
- ✅ Extract a custom error type from a wrapped chain with errors.As
- ✅ Pass a pointer to a variable of the error's type as the target
- ✅ Make a decision based on the extracted error's fields

---

## Exercise 6: A Custom ValidationError With Structured Fields

**Objective:** Design an error type that carries data, not just a message

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level16-exercise6
cd ~/projects/level16-exercise6
go mod init level16.example/exercise6
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "errors"
    "fmt"
)

type ValidationError struct {
    Field  string
    Reason string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("invalid %s: %s", e.Field, e.Reason)
}

func validateForm(username, email string) error {
    if len(username) < 3 {
        return &ValidationError{Field: "username", Reason: "must be at least 3 characters"}
    }
    if email == "" {
        return &ValidationError{Field: "email", Reason: "must not be empty"}
    }
    return nil
}

func main() {
    fmt.Println("=== A Struct Can Be an Error ===")
    submissions := []struct{ username, email string }{
        {"al", "al@example.com"},
        {"alice", ""},
        {"alice", "alice@example.com"},
    }

    for _, s := range submissions {
        err := validateForm(s.username, s.email)
        if err == nil {
            fmt.Printf("%q / %q -> ok\n", s.username, s.email)
            continue
        }
        fmt.Printf("%q / %q -> %v\n", s.username, s.email, err)

        var vErr *ValidationError
        if errors.As(err, &vErr) {
            fmt.Printf("    field=%q reason=%q\n", vErr.Field, vErr.Reason)
        }
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
=== A Struct Can Be an Error ===
"al" / "al@example.com" -> invalid username: must be at least 3 characters
    field="username" reason="must be at least 3 characters"
"alice" / "" -> invalid email: must not be empty
    field="email" reason="must not be empty"
"alice" / "alice@example.com" -> ok
```

**Learning Objectives:**
- ✅ Implement the error interface on your own struct
- ✅ Carry structured data (field, reason) inside an error
- ✅ Let callers read fields directly instead of parsing message strings

---

## Exercise 7: The Typed-Nil Error Trap (and the Fix)

**Objective:** Reproduce the most famous error-handling bug in Go, then fix it

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level16-exercise7
cd ~/projects/level16-exercise7
go mod init level16.example/exercise7
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type ConfigError struct {
    Path string
}

func (e *ConfigError) Error() string {
    return "bad config file: " + e.Path
}

// ❌ BROKEN: declares a *ConfigError return type
func loadBroken(path string) *ConfigError {
    if path == "" {
        return &ConfigError{Path: "(empty path)"}
    }
    return nil // typed nil: a nil *ConfigError, not a nil error
}

// ✅ FIXED: declares the error INTERFACE as the return type
func loadFixed(path string) error {
    if path == "" {
        return &ConfigError{Path: "(empty path)"}
    }
    return nil // untyped nil: a truly nil error
}

func main() {
    fmt.Println("=== The Trap ===")
    var err error = loadBroken("app.conf") // load succeeded...
    fmt.Println("err:", err)
    fmt.Printf("err's type: %T\n", err)
    fmt.Println("err == nil:", err == nil)
    if err != nil {
        fmt.Println("BUG: we report a failure that never happened!")
    }

    fmt.Println("\n=== The Fix ===")
    err = loadFixed("app.conf")
    fmt.Println("err:", err)
    fmt.Printf("err's type: %T\n", err)
    fmt.Println("err == nil:", err == nil)
    if err == nil {
        fmt.Println("correct: success is reported as success")
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
=== The Trap ===
err: <nil>
err's type: *main.ConfigError
err == nil: false
BUG: we report a failure that never happened!

=== The Fix ===
err: <nil>
err's type: <nil>
err == nil: true
correct: success is reported as success
```

Look closely at the trap: the broken version even *prints* as `<nil>`, but `err == nil` is `false` - the interface holds a `(*ConfigError, nil)` pair, and that's not a nil interface. This is exactly the typed-nil gotcha from Level 15, in its most dangerous real-world form.

**Learning Objectives:**
- ✅ Reproduce the typed-nil error trap with a concrete error return type
- ✅ Explain it via the interface's (type, value) pair from Level 15
- ✅ Fix it by always declaring `error` as the return type

---

## Exercise 8: panic and recover

**Objective:** Watch a real panic crash a program, then catch one with recover

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level16-exercise8
cd ~/projects/level16-exercise8
go mod init level16.example/exercise8
```

2. First, create a `main.go` that panics with NO recovery:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func mustPositive(n int) int {
    if n <= 0 {
        panic(fmt.Sprintf("mustPositive called with %d", n))
    }
    return n
}

func main() {
    defer fmt.Println("deferred: I still run while the panic unwinds")
    fmt.Println("before the panic")
    mustPositive(-5)
    fmt.Println("after the panic - never printed")
}
EOF
```

3. Run it and watch it crash:

```bash
go run main.go
```

**Expected Output (the crash):**

```
before the panic
deferred: I still run while the panic unwinds
panic: mustPositive called with -5

goroutine 1 [running]:
main.mustPositive(...)
        /path/to/main.go:7
main.main()
        /path/to/main.go:15 +0xd0
exit status 2
```

The exact file path and offset will differ on your machine, but everything else - including the deferred line printing BEFORE the crash - is exactly what you'll see. Note the panic message, the goroutine stack trace, and `exit status 2`.

4. Now replace `main.go` with a version that recovers:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func mustPositive(n int) int {
    if n <= 0 {
        panic(fmt.Sprintf("mustPositive called with %d", n))
    }
    return n
}

func safeCall(n int) (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("recovered from panic: %v", r)
        }
    }()
    return mustPositive(n), nil
}

func main() {
    fmt.Println("=== A Call That Does Not Panic ===")
    result, err := safeCall(10)
    fmt.Println("result:", result, "err:", err)

    fmt.Println("\n=== A Call That Panics - and Is Recovered ===")
    result, err = safeCall(-5)
    fmt.Println("result:", result, "err:", err)

    fmt.Println("\n=== The Program Is Still Alive ===")
    fmt.Println("main keeps running after the recovered panic")
}
EOF
```

5. Run it again:

```bash
go run main.go
```

**Expected Output (recovered):**

```
=== A Call That Does Not Panic ===
result: 10 err: <nil>

=== A Call That Panics - and Is Recovered ===
result: 0 err: recovered from panic: mustPositive called with -5

=== The Program Is Still Alive ===
main keeps running after the recovered panic
```

**Learning Objectives:**
- ✅ See a real uncaught panic: defers run, then stack trace, then exit status 2
- ✅ Use recover() inside a deferred function to stop a panic
- ✅ Convert a panic into an ordinary error using a named return value

---

## Exercise 9: Handle Errors Exactly Once (Refactoring)

**Objective:** Start with the log-and-return anti-pattern, then refactor to handle the error once

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level16-exercise9
cd ~/projects/level16-exercise9
go mod init level16.example/exercise9
```

2. First, create the ANTI-PATTERN version - every layer logs AND returns:

```bash
cat > main.go << 'EOF'
package main

import (
    "errors"
    "fmt"
)

// ❌ ANTI-PATTERN: every layer logs AND returns the same error

func readSettings(name string) (string, error) {
    if name != "settings.json" {
        err := errors.New("file does not exist")
        fmt.Println("ERROR in readSettings:", err) // logged here...
        return "", err                             // ...AND returned
    }
    return `{"debug": true}`, nil
}

func loadApp(name string) (string, error) {
    data, err := readSettings(name)
    if err != nil {
        fmt.Println("ERROR in loadApp:", err) // logged AGAIN...
        return "", err                        // ...and returned again
    }
    return data, nil
}

func main() {
    _, err := loadApp("missing.json")
    if err != nil {
        fmt.Println("ERROR in main:", err) // third log of ONE failure!
    }
}
EOF
```

3. Run it and count how many times ONE failure gets reported:

```bash
go run main.go
```

**Expected Output (the anti-pattern):**

```
ERROR in readSettings: file does not exist
ERROR in loadApp: file does not exist
ERROR in main: file does not exist
```

One missing file produced three error reports - and none of them says *what the program was trying to do*. In a real log, nobody can tell if this was one incident or three.

4. Now refactor: lower layers only add context and return; main handles the error exactly once:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

// ✅ FIXED: lower layers only ADD CONTEXT and return;
// exactly one place (main) handles the error

func readSettings(name string) (string, error) {
    if name != "settings.json" {
        return "", fmt.Errorf("settings file %q does not exist", name)
    }
    return `{"debug": true}`, nil
}

func loadApp(name string) (string, error) {
    data, err := readSettings(name)
    if err != nil {
        return "", fmt.Errorf("loading app: %w", err)
    }
    return data, nil
}

func main() {
    _, err := loadApp("missing.json")
    if err != nil {
        fmt.Println("ERROR:", err) // handled exactly once, full story intact
    }
}
EOF
```

5. Run it again:

```bash
go run main.go
```

**Expected Output (refactored):**

```
ERROR: loading app: settings file "missing.json" does not exist
```

One failure, one report - and thanks to the context added at each layer, this single line tells a *better* story than the three noisy lines did.

**Learning Objectives:**
- ✅ Recognize the log-and-return anti-pattern and the duplicate reports it causes
- ✅ Refactor so each layer either adds context and returns, or handles - never both
- ✅ See that one wrapped error line beats several raw log lines

---

## Exercise 10: Comprehensive Practice - User Registration Validator

**Objective:** Combine sentinels, custom error types, wrapping, errors.Is, and errors.As in one real program

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level16-exercise10
cd ~/projects/level16-exercise10
go mod init level16.example/exercise10
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "errors"
    "fmt"
    "strconv"
)

var ErrDuplicateUsername = errors.New("username already taken")

type ValidationError struct {
    Field  string
    Reason string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("invalid %s: %s", e.Field, e.Reason)
}

var existingUsers = map[string]bool{"alice": true, "bob": true}

func validateUsername(username string) error {
    if len(username) < 3 {
        return &ValidationError{Field: "username", Reason: "must be at least 3 characters"}
    }
    if existingUsers[username] {
        return fmt.Errorf("checking availability: %w", ErrDuplicateUsername)
    }
    return nil
}

func validateAge(rawAge string) error {
    age, err := strconv.Atoi(rawAge)
    if err != nil {
        return fmt.Errorf("parsing age: %w", err)
    }
    if age < 13 || age > 120 {
        return &ValidationError{Field: "age", Reason: fmt.Sprintf("must be between 13 and 120, got %d", age)}
    }
    return nil
}

func register(username, rawAge string) error {
    if err := validateUsername(username); err != nil {
        return fmt.Errorf("registering %q: %w", username, err)
    }
    if err := validateAge(rawAge); err != nil {
        return fmt.Errorf("registering %q: %w", username, err)
    }
    return nil
}

func main() {
    candidates := []struct{ username, age string }{
        {"charlie", "30"},
        {"alice", "25"},
        {"al", "40"},
        {"diana", "abc"},
        {"erin", "7"},
    }

    for _, c := range candidates {
        fmt.Printf("--- %s / %s ---\n", c.username, c.age)
        err := register(c.username, c.age)
        if err == nil {
            fmt.Println("registered - welcome aboard!")
            continue
        }

        var vErr *ValidationError
        switch {
        case errors.Is(err, ErrDuplicateUsername):
            fmt.Println("reaction: that name is taken - pick another username")
        case errors.As(err, &vErr):
            fmt.Printf("reaction: fix the %s field (%s)\n", vErr.Field, vErr.Reason)
        case errors.Is(err, strconv.ErrSyntax):
            fmt.Println("reaction: age must be a whole number")
        default:
            fmt.Println("reaction: unexpected error")
        }
        fmt.Println("full error:", err)
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
--- charlie / 30 ---
registered - welcome aboard!
--- alice / 25 ---
reaction: that name is taken - pick another username
full error: registering "alice": checking availability: username already taken
--- al / 40 ---
reaction: fix the username field (must be at least 3 characters)
full error: registering "al": invalid username: must be at least 3 characters
--- diana / abc ---
reaction: age must be a whole number
full error: registering "diana": parsing age: strconv.Atoi: parsing "abc": invalid syntax
--- erin / 7 ---
reaction: fix the age field (must be between 13 and 120, got 7)
full error: registering "erin": invalid age: must be between 13 and 120, got 7
```

Notice `diana`'s case: `strconv.Atoi` returns an error that itself wraps the standard library's `strconv.ErrSyntax` sentinel - so `errors.Is` matches it even through your own `parsing age:` and `registering "diana":` wrappers. Three layers of chain, one clean check.

**Learning Objectives:**
- ✅ Combine a sentinel, a custom type, and wrapped ad-hoc errors in one design
- ✅ React differently per error kind with an errors.Is / errors.As switch
- ✅ Match a standard-library sentinel (strconv.ErrSyntax) through your own wrappers
- ✅ Keep the full wrapped story available for logging while reacting on the classification

---

## Bonus Challenges

### Challenge 1: Retry Helper That Gives Up on Permanent Errors

Write a `retry(attempts int, fn func() error) error` helper that calls `fn` up to `attempts` times - but gives up immediately if the error is permanent (no point retrying a "not authorized" failure).

```bash
mkdir -p ~/projects/level16-bonus1
cd ~/projects/level16-bonus1
go mod init level16.example/bonus1
```

**Hints:**
- Define `var ErrPermanent = errors.New("permanent failure")`
- Inside the loop: if `fn()` returns nil, return nil; if `errors.Is(err, ErrPermanent)`, return the error immediately
- Have test functions wrap the sentinel (`fmt.Errorf("auth rejected: %w", ErrPermanent)`) to prove `errors.Is` still catches it
- Simulate a flaky operation with a counter closure (Level 11) that fails twice, then succeeds
- After the loop, return `fmt.Errorf("giving up after %d attempts: %w", attempts, err)`

### Challenge 2: Multi-Error Collector With errors.Join

A form validator that stops at the FIRST bad field is annoying - collect ALL the problems and report them together using `errors.Join` (standard library, Go 1.20+).

```bash
mkdir -p ~/projects/level16-bonus2
cd ~/projects/level16-bonus2
go mod init level16.example/bonus2
```

**Hints:**
- Collect each failed check into a `[]error` slice, then `return errors.Join(errs...)`
- Verified on this toolchain: printing a joined error shows each message on its own line (they're separated by `\n`)
- `errors.Join` with all-nil arguments returns nil - so an empty/all-nil slice means "valid form" automatically
- `errors.Is` and `errors.As` search EVERY branch of a joined error - verify that a sentinel joined in among other errors is still found
- Validate at least three fields (username, email, age) and submit a form that fails all three

### Challenge 3: Implement Unwrap() on Your Own Wrapper Type

`fmt.Errorf` with `%w` isn't magic - any type can join the wrapping ecosystem by implementing `Unwrap() error`. Build a `TimestampError` that wraps another error and adds when it happened.

```bash
mkdir -p ~/projects/level16-bonus3
cd ~/projects/level16-bonus3
go mod init level16.example/bonus3
```

**Hints:**
- Struct fields: `When time.Time` and `Err error`
- Implement BOTH `Error() string` (include `e.When.Format(time.RFC3339)` and `e.Err`) and `Unwrap() error` (just `return e.Err`)
- Wrap a sentinel like `ErrNotFound` inside a `&TimestampError{...}` and prove `errors.Is(err, ErrNotFound)` returns true - errors.Is finds it by calling YOUR Unwrap method
- Comment out the `Unwrap` method and rerun: `errors.Is` now returns false. That's the whole mechanism revealed
- Extra credit: scavenge the `errors` package docs (`go doc errors`) and list which functions rely on `Unwrap`

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Create errors with errors.New and fmt.Errorf and handle them with if err != nil
✅ Add context at every layer and return early, keeping the happy path flat
✅ Build and walk error chains with %w and errors.Unwrap
✅ Match sentinels through wrapped chains with errors.Is
✅ Extract typed errors and their fields with errors.As
✅ Design custom error types with structured data
✅ Avoid the typed-nil error trap by returning the error interface
✅ Read a real panic stack trace and recover panics into ordinary errors
✅ Handle each error exactly once - context on the way up, one handler at the top
✅ Build a full validation flow where callers react differently per error kind

---

## Next Level

Level 17: Packages & Modules
- Organizing code into multiple packages
- Exported vs unexported identifiers at package scale
- go.mod, module paths, and third-party dependencies
- Designing clean package APIs (including which errors to export)

Great work! You've completed the OOP CONCEPTS tier! 🚀
