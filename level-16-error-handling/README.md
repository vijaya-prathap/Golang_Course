# Level 16: Error Handling - Complete Guide

## Introduction

Welcome to Level 16! You've mastered interfaces (Level 15). Now it's time to make **errors** the explicit subject - the pattern you've been typing since Level 3 (`value, err := strconv.Atoi(...)`) finally gets fully explained.

Here's the punchline up front: **`error` is just an interface**. Everything you learned in Level 15 - implicit satisfaction, type assertions, the typed-nil gotcha - applies directly to error handling. That's why this level closes out the **OOP CONCEPTS** tier (Levels 12-16): errors are where pointers, structs, methods, and interfaces all come together in everyday Go code.

Go doesn't throw exceptions. Errors are ordinary **values** that functions return, and you handle them with ordinary code: `if err != nil`. This level teaches you to create, wrap, inspect, and design errors like a Go professional - plus the escape hatch (`panic`/`recover`) and when it's actually appropriate.

---

## Table of Contents

1. [Errors Are Values: The error Interface](#errors-are-values-the-error-interface)
2. [Creating Errors: errors.New and fmt.Errorf](#creating-errors-errorsnew-and-fmterrorf)
3. [The Standard Flow: Return Early, Add Context](#the-standard-flow-return-early-add-context)
4. [Error Wrapping With %w](#error-wrapping-with-w)
5. [errors.Is: Matching Sentinel Errors](#errorsis-matching-sentinel-errors)
6. [errors.As: Extracting Custom Error Types](#errorsas-extracting-custom-error-types)
7. [Custom Error Types](#custom-error-types)
8. [The Typed-Nil Error Trap](#the-typed-nil-error-trap)
9. [panic and recover](#panic-and-recover)
10. [defer for Cleanup Alongside Errors](#defer-for-cleanup-alongside-errors)
11. [Best Practices](#best-practices)
12. [Common Mistakes](#common-mistakes)

---

## Errors Are Values: The error Interface

The entire error system rests on one tiny interface built into the language:

```go
type error interface {
    Error() string
}
```

That's it. Any type with an `Error() string` method **is** an error - implicitly, exactly like every other interface you met in Level 15. No special keywords, no exception classes, no `throw`.

### Why Values Instead of Exceptions?

In exception-based languages, any call might secretly jump to a distant `catch` block. In Go, a function that can fail **says so in its signature**:

```go
func divide(a, b float64) (float64, error)
```

The error is part of the return values, so:

- You can't call the function without *seeing* that it can fail
- Handling failure is normal `if` code, not a separate control-flow system
- The happy path and the failure path are both visible right where the call happens

### The Fundamental Pattern

```go
result, err := divide(10, 2)
if err != nil {
    // failure path: handle it, then usually return or continue
    fmt.Println("Error:", err)
    return
}
// success path: err was nil, result is safe to use
fmt.Println("10 / 2 =", result)
```

`nil` means "no error". Anything non-nil means "something went wrong, and calling `.Error()` on me tells you what."

### A Complete Example

```go
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

func main() {
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
}
```

Verified output:

```
10 / 2 = 5
Error: division by zero
```

### Under the Hood

An error created with `errors.New` is a pointer to a small unexported struct in the `errors` package:

```go
e := errors.New("something went wrong")
fmt.Printf("Type: %T\n", e)
fmt.Println("Message:", e.Error())
```

Verified output:

```
Type: *errors.errorString
Message: something went wrong
```

Printing an error with `fmt.Println(err)` automatically calls its `Error()` method - you almost never call `.Error()` yourself.

---

## Creating Errors: errors.New and fmt.Errorf

### errors.New - Fixed Messages

```go
err := errors.New("division by zero")
```

Use it when the message needs no dynamic values.

### fmt.Errorf - Formatted Messages

```go
err := fmt.Errorf("user %d not found", id)
```

`fmt.Errorf` works exactly like `fmt.Sprintf` but produces an `error`. Use it whenever the message should include values from the situation.

```go
func findUser(id int) (string, error) {
    users := map[int]string{1: "Alice", 2: "Bob"}
    name, ok := users[id]
    if !ok {
        return "", fmt.Errorf("user %d not found", id)
    }
    return name, nil
}
```

Verified output (calling with id 1, then 42):

```
Found: Alice
Error: user 42 not found
```

### Error Message Style Conventions

Go has firm conventions for error strings, because messages get **joined into chains** (you'll see how in the wrapping section):

| Convention | ❌ Bad | ✅ Good |
|-----------|--------|---------|
| Lowercase first letter | `"File not found"` | `"file not found"` |
| No trailing punctuation | `"file not found."` | `"file not found"` |
| Include relevant values | `"bad port"` | `"port 99999 out of range (1-65535)"` |
| Describe the operation, not blame | `"you passed a bad id"` | `"user 42 not found"` |

Why lowercase and no period? Because errors get embedded mid-sentence when wrapped:

```
loading config: parsing port "abc": invalid syntax
```

A capitalized `Parsing Port "abc".` would look broken inside that chain.

---

## The Standard Flow: Return Early, Add Context

Idiomatic Go error handling has a rhythm you'll see in every real codebase:

1. Call something that can fail
2. `if err != nil` → **return early**, adding context
3. Continue with the happy path, un-indented

```go
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
```

Each layer **adds one clause of context** describing what *it* was doing. The result reads like a story, from the outermost operation down to the root cause:

Verified output for `loadConfig("abc")`:

```
Error: loading config: parsing port "abc": strconv.Atoi: parsing "abc": invalid syntax
```

Verified output for `loadConfig("99999")`:

```
Error: loading config: port 99999 out of range (1-65535)
```

### Why Return Early?

```go
// ✅ Good - failure paths exit immediately, happy path stays flat
func process() error {
    a, err := stepOne()
    if err != nil {
        return fmt.Errorf("step one: %w", err)
    }
    b, err := stepTwo(a)
    if err != nil {
        return fmt.Errorf("step two: %w", err)
    }
    return finish(b)
}

// ❌ Avoid - nesting the happy path inside else blocks
func process() error {
    a, err := stepOne()
    if err == nil {
        b, err := stepTwo(a)
        if err == nil {
            return finish(b)
        } else {
            return err
        }
    } else {
        return err
    }
}
```

The rule of thumb: **the happy path flows down the left edge; errors indent and leave.**

---

## Error Wrapping With %w

`fmt.Errorf` has a special verb: `%w` (w for **wrap**). It formats exactly like `%v`, but it also **records the original error inside the new one**, building a chain that tools can walk later.

```go
base := errors.New("connection refused")
wrapped := fmt.Errorf("dialing database: %w", base)
top := fmt.Errorf("starting server: %w", wrapped)
```

Now `top` isn't just a longer string - it's a linked chain: `top → wrapped → base`.

### errors.Unwrap Walks the Chain

```go
fmt.Println(top)
fmt.Println(errors.Unwrap(top))
fmt.Println(errors.Unwrap(errors.Unwrap(top)))
fmt.Println(errors.Unwrap(base))
```

Verified output:

```
starting server: dialing database: connection refused
dialing database: connection refused
connection refused
<nil>
```

Each `Unwrap` peels off one layer; unwrapping past the end gives `nil`.

### %w vs %v - Only %w Builds a Chain

```go
base := errors.New("connection refused")
plain := fmt.Errorf("dialing database: %v", base)

fmt.Println("wrapped with %v, Unwrap gives:", errors.Unwrap(plain))
```

Verified output:

```
wrapped with %v, Unwrap gives: <nil>
```

With `%v` the original error's *text* is copied in, but its *identity* is lost - `errors.Is` and `errors.As` (next two sections) can no longer find it. **Default to `%w` when adding context.** Use `%v` only when you deliberately want to hide the underlying error from callers (for example, at an API boundary where the cause is an internal detail).

You rarely call `errors.Unwrap` yourself. Its real job is powering the two functions you *will* use constantly: `errors.Is` and `errors.As`.

---

## errors.Is: Matching Sentinel Errors

A **sentinel error** is a package-level error value that callers can compare against - a named, exported "kind of failure":

```go
var ErrNotFound = errors.New("not found")
```

The convention: sentinel names start with `Err`. The standard library is full of them (`io.EOF`, `sql.ErrNoRows`, `os.ErrNotExist`).

### The Problem With ==

Once errors get wrapped, a direct comparison fails - the outermost error is a wrapper, not the sentinel itself. `errors.Is(err, target)` fixes this: it reports whether `target` appears **anywhere in err's chain**.

### Verified Demo: A Sentinel Through Two Layers

```go
var ErrNotFound = errors.New("not found")

func fetchUser(id int) error {
    return fmt.Errorf("user %d: %w", id, ErrNotFound)
}

func renderProfile(id int) error {
    if err := fetchUser(id); err != nil {
        return fmt.Errorf("rendering profile: %w", err)
    }
    return nil
}

func main() {
    err := renderProfile(7)
    fmt.Println(err)
    fmt.Println("errors.Is:", errors.Is(err, ErrNotFound))
    fmt.Println("== compare:", err == ErrNotFound)
}
```

Verified output:

```
rendering profile: user 7: not found
errors.Is: true
== compare: false
```

The sentinel was wrapped twice (`fetchUser` wrapped it, `renderProfile` wrapped it again), yet `errors.Is` still finds it. `==` sees only the outermost wrapper and says `false`.

**Rule: whenever wrapping *might* have happened, compare with `errors.Is`, never `==`.**

`errors.Is(nil, ErrNotFound)` is simply `false` - safe to call on the success path.

---

## errors.As: Extracting Custom Error Types

`errors.Is` answers "is this **that specific value**?" `errors.As` answers a different question: "is there an error of this **type** in the chain - and if so, give it to me so I can read its fields."

```go
type QueryError struct {
    Query string
    Code  int
}

func (e *QueryError) Error() string {
    return fmt.Sprintf("query failed (code %d): %s", e.Code, e.Query)
}

func main() {
    err := fmt.Errorf("generating report: %w", &QueryError{Query: "SELECT *", Code: 1064})

    var qErr *QueryError
    if errors.As(err, &qErr) {
        fmt.Println("full error:", err)
        fmt.Println("code:      ", qErr.Code)
        fmt.Println("query:     ", qErr.Query)
    }
}
```

Verified output:

```
full error: generating report: query failed (code 1064): SELECT *
code:       1064
query:      SELECT *
```

### How to Read That Call

```go
var qErr *QueryError      // 1. declare a variable of the error's type
if errors.As(err, &qErr)  // 2. pass a POINTER to it
```

`errors.As` walks the chain; if it finds a `*QueryError`, it copies it into `qErr` and returns `true`. It's the error-chain-aware cousin of the type assertion `err.(*QueryError)` you learned in Level 15 - a plain assertion only checks the outermost error, `errors.As` checks the whole chain.

If no error of that type exists anywhere in the chain, `errors.As` returns `false` and leaves the target untouched.

---

## Custom Error Types

Sentinels carry only an identity. When a failure has **structured data** callers need - which field failed, what status code came back, which line of a file was bad - define a custom error type: a struct with an `Error() string` method.

```go
type ValidationError struct {
    Field  string
    Reason string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("invalid %s: %s", e.Field, e.Reason)
}
```

That pointer-receiver method makes `*ValidationError` satisfy the `error` interface (Level 14 receivers + Level 15 implicit satisfaction, working together).

### Using It

```go
func validateForm(username, email string) error {
    if len(username) < 3 {
        return &ValidationError{Field: "username", Reason: "must be at least 3 characters"}
    }
    if email == "" {
        return &ValidationError{Field: "email", Reason: "must not be empty"}
    }
    return nil
}
```

Verified output (submitting `"al" / "al@example.com"`, then `"alice" / ""`):

```
"al" / "al@example.com" -> invalid username: must be at least 3 characters
"alice" / "" -> invalid email: must not be empty
```

Callers who only print the error get a clean message. Callers who need details use `errors.As` and read `.Field` and `.Reason` directly - no string parsing.

### Sentinel vs Custom Type - When to Use Which

| Situation | Use |
|-----------|-----|
| Callers only need to recognize the *kind* of failure | Sentinel + `errors.Is` |
| Callers need *data* about the failure (field, code, position) | Custom type + `errors.As` |
| Callers will only ever print the message | Plain `fmt.Errorf` (ad-hoc) |

Don't over-engineer: most errors in real code are ad-hoc `fmt.Errorf` calls. Reach for sentinels and custom types only when a caller genuinely needs to *react differently* to that failure.

---

## The Typed-Nil Error Trap

This is the most famous error-handling bug in Go, and it's exactly the typed-nil interface gotcha you saw in Level 15 - now wearing its most dangerous costume.

### The Reproduction

```go
type MyError struct{}

func (e *MyError) Error() string { return "boom" }

func broken() *MyError {
    return nil // typed nil: a nil *MyError
}

func main() {
    var err error = broken()
    fmt.Println("err == nil:", err == nil)
    fmt.Printf("err's type: %T\n", err)
}
```

Verified output:

```
err == nil: false
err's type: *main.MyError
```

The function succeeded and returned `nil` - yet `err != nil` is **true**, so every caller thinks it failed!

### Why It Happens

As you saw in Level 15, an interface value holds a (type, value) pair. Assigning a `nil *MyError` to an `error` variable produces an interface holding `(*MyError, nil)` - a non-nil interface wrapping a nil pointer. `err == nil` only succeeds when **both** halves are nil.

### The Fix

**Declare error returns as the `error` interface type, never as a concrete error pointer type:**

```go
// ❌ BROKEN - concrete error type in the signature
func loadBroken(path string) *ConfigError {
    // ...
    return nil // becomes a non-nil error at the call site!
}

// ✅ FIXED - the interface type in the signature
func loadFixed(path string) error {
    if path == "" {
        return &ConfigError{Path: "(empty path)"}
    }
    return nil // a genuinely nil error
}
```

Verified output comparing both versions on the success path:

```
=== The Trap ===
err's type: *main.ConfigError
err == nil: false

=== The Fix ===
err's type: <nil>
err == nil: true
```

With `error` as the declared return type, the literal `nil` in the return statement is an untyped nil - the interface stays `(nil, nil)` and `err == nil` behaves correctly.

---

## panic and recover

Errors are for failures you **expect** and handle. `panic` is for situations where continuing is meaningless - and it should be rare.

### What panic Does

When a function panics:

1. Normal execution stops immediately
2. All **deferred calls in the current function run** (Level 11's `defer`, doing its job to the last)
3. The panic unwinds up the call stack, running each caller's defers too
4. If nothing recovers it, the program crashes, printing the panic value and a stack trace, and exits with status 2

```go
func main() {
    fmt.Println("before the panic")
    panic("something went terribly wrong")
}
```

Verified output:

```
before the panic
panic: something went terribly wrong

goroutine 1 [running]:
main.main()
        /path/to/main.go:7 +0x60
exit status 2
```

The exact file path and offset will differ on your machine, but the shape - `panic:` line, `goroutine 1 [running]:`, the stack trace, `exit status 2` - is always the same.

### Defers Run During the Unwind

```go
func main() {
    defer fmt.Println("deferred 1 (registered first, runs last)")
    defer fmt.Println("deferred 2 (registered last, runs first)")
    panic("boom")
}
```

Verified output:

```
deferred 2 (registered last, runs first)
deferred 1 (registered first, runs last)
panic: boom

goroutine 1 [running]:
main.main()
        /path/to/main.go:8 +0x8c
exit status 2
```

Both defers ran - in LIFO order, as always - *before* the crash. This is why `defer`-based cleanup is reliable even in the face of panics.

### recover: Catching a Panic

`recover()` stops an in-flight panic - but **only when called directly inside a deferred function**. Called anywhere else, it returns `nil` and does nothing.

```go
func safeDivide(a, b int) (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("recovered: %v", r)
        }
    }()
    return a / b, nil
}

func main() {
    result, err := safeDivide(10, 2)
    fmt.Println("result:", result, "err:", err)

    result, err = safeDivide(10, 0)
    fmt.Println("result:", result, "err:", err)

    fmt.Println("still running!")
}
```

Verified output:

```
result: 5 err: <nil>
result: 0 err: recovered: runtime error: integer divide by zero
still running!
```

Notice the idiom: the deferred closure converts the panic into an ordinary `error` through the **named return value** `err` (Level 11) - so callers see a normal Go error, not a crash. Dividing by zero is a real runtime panic, and `recover` caught it.

### When to panic vs When to Return an Error

| Situation | Do this |
|-----------|---------|
| File missing, bad user input, network down, invalid data | **Return an error** - expected failures |
| A bug: "impossible" state, broken invariant, nil that can't be nil | `panic` - programmer error |
| Initialization that the program cannot run without | `panic` (or `log.Fatal`) at startup |
| Library code hitting an internal bug | `panic` internally, `recover` at the API boundary, return an error |

The overwhelming default is **return an error**. If you're reaching for `panic` because typing `if err != nil` feels tedious - stop. That tedium is Go working as designed: every failure path visible and handled.

---

## defer for Cleanup Alongside Errors

Real functions acquire resources (files, connections, locks) that must be released **no matter which of the many `return err` paths runs**. `defer` (Level 11) plus early returns (this level) is the pattern:

```go
func process() error {
    r, err := open("data.txt")
    if err != nil {
        return fmt.Errorf("opening resource: %w", err)
    }
    defer r.Close() // runs no matter which return path we take

    fmt.Println("working with", r.name)
    return nil
}
```

Verified output (with `open` and `Close` printing trace lines):

```
opening data.txt
working with data.txt
closing data.txt
done
```

The two rules that make this pattern safe:

1. **Check the error BEFORE deferring the cleanup.** If `open` failed there's nothing to close - deferring `r.Close()` first would call a method on a nil resource.
2. **Defer immediately after a successful acquisition.** Every later `return err` - and even a panic - still triggers the cleanup.

This is a sketch with a toy `resource` type; in Level 18 (File Handling) you'll use exactly this shape with real `os.File` values.

---

## Best Practices

### 1. Handle Each Error Exactly Once

An error should be either **returned** (with context) or **handled** (logged, retried, converted into a response) - never both.

```go
// ❌ Bad - logs AND returns; the same failure gets reported at every layer
if err != nil {
    fmt.Println("ERROR in readSettings:", err)
    return "", err
}

// ✅ Good - lower layers add context and return; ONE place at the top handles it
if err != nil {
    return "", fmt.Errorf("reading settings: %w", err)
}
```

Log-and-return produces logs where one failure appears three or four times, and nobody can tell whether that was one incident or four.

### 2. Wrap With Context at Each Layer

```go
// ✅ Good - each layer says what IT was doing
return fmt.Errorf("loading config: %w", err)

// ❌ Bad - passing the error up raw; the caller learns nothing about where/why
return err
```

### 3. Keep Messages Lowercase, No Trailing Punctuation

```go
// ✅ Good
errors.New("connection refused")

// ❌ Bad - breaks mid-chain formatting
errors.New("Connection refused.")
```

### 4. Compare With errors.Is, Not ==, When Chains May Exist

```go
// ✅ Good - finds the sentinel through any number of wraps
if errors.Is(err, ErrNotFound) { ... }

// ❌ Fragile - breaks the moment anyone adds a %w wrapper
if err == ErrNotFound { ... }
```

### 5. Return the error Interface, Never a Concrete Error Pointer

```go
// ✅ Good
func load() error

// ❌ Bad - invites the typed-nil trap
func load() *ConfigError
```

### 6. Reserve panic for the Truly Unrecoverable

Expected failures - bad input, missing files, network problems - are errors. `panic` is for bugs and impossible states. If a user can trigger it with bad input, it must be an error, not a panic.

---

## Common Mistakes

### Mistake 1: Ignoring Errors With _

```go
// ❌ WRONG - the failure vanishes silently
value, _ := strconv.Atoi(input)

// ✅ RIGHT
value, err := strconv.Atoi(input)
if err != nil {
    return fmt.Errorf("parsing input: %w", err)
}
```

Every `_, _ =` discarding an error is a landmine. If you truly don't care, that's rare enough to deserve a comment saying why.

### Mistake 2: Shadowing err Inside if

```go
// ❌ WRONG - := inside the block declares NEW variables
var data string
var err error
if useCache {
    data, err := readCache() // shadows BOTH outer variables!
    _ = data
}
if err != nil { // outer err is still nil - failure silently missed
    return err
}

// ✅ RIGHT - plain = assigns to the existing variables
if useCache {
    data, err = readCache()
}
```

`:=` creates fresh variables scoped to the block (Level 5 scoping rules). The compiler often can't warn you, because the inner variables *are* used.

### Mistake 3: Returning Typed Nils

```go
// ❌ WRONG - a nil *ValidationError becomes a NON-nil error
func validate() *ValidationError {
    return nil
}

// ✅ RIGHT - declare error, return untyped nil
func validate() error {
    return nil
}
```

(See [The Typed-Nil Error Trap](#the-typed-nil-error-trap) for the verified reproduction.)

### Mistake 4: Using == on Wrapped Errors

```go
err := fmt.Errorf("loading user: %w", ErrNotFound)

// ❌ WRONG - false, the outermost error is the wrapper
if err == ErrNotFound { ... }

// ✅ RIGHT - true, errors.Is walks the chain
if errors.Is(err, ErrNotFound) { ... }
```

### Mistake 5: Panicking for Normal Failures

```go
// ❌ WRONG - a missing user is not a program bug
func findUser(id int) string {
    name, ok := users[id]
    if !ok {
        panic("user not found")
    }
    return name
}

// ✅ RIGHT - expected failures are errors
func findUser(id int) (string, error) {
    name, ok := users[id]
    if !ok {
        return "", fmt.Errorf("user %d not found", id)
    }
    return name, nil
}
```

### Mistake 6: Losing the Chain With %v

```go
// ❌ WRONG - %v copies the text but breaks the chain
return fmt.Errorf("loading config: %v", err)
// ...errors.Is/errors.As can no longer see the original error

// ✅ RIGHT - %w preserves the chain
return fmt.Errorf("loading config: %w", err)
```

---

## Summary

**The Foundation:**
- `error` is a one-method interface: `Error() string` - any type with that method is an error
- Errors are returned values, not thrown exceptions; `nil` means success
- The rhythm: call, `if err != nil`, return early with context, happy path stays flat

**Creating and Wrapping:**
- `errors.New` for fixed messages, `fmt.Errorf` for formatted ones
- Messages: lowercase, no trailing punctuation, include relevant values
- `%w` wraps the original error into a chain; `%v` copies only its text
- `errors.Unwrap` peels one layer; you rarely call it directly

**Inspecting Chains:**
- Sentinel errors (`var ErrNotFound = errors.New(...)`) + `errors.Is` = "is this kind of failure anywhere in the chain?"
- Custom error types (struct + `Error()` method) + `errors.As` = "extract the typed error so I can read its fields"
- `==` and plain type assertions only see the outermost error - use `Is`/`As`

**The Traps and the Escape Hatch:**
- Never declare a concrete error pointer as a return type - typed nil makes `err != nil` lie
- `panic` unwinds the stack, running defers as it goes; unrecovered panics crash with a stack trace
- `recover()` works only directly inside a deferred function; the idiom converts panics to errors via named returns
- Errors for expected failures, panic for bugs and impossible states

---

## Next Steps

You now understand:
- ✅ The `error` interface and why Go treats errors as values
- ✅ Creating errors with `errors.New` and `fmt.Errorf`, with idiomatic messages
- ✅ The return-early flow and adding context at every layer
- ✅ Wrapping with `%w` and inspecting chains with `errors.Is` and `errors.As`
- ✅ Sentinel errors, custom error types, and when to use each
- ✅ The typed-nil error trap and how to avoid it
- ✅ `panic`, `recover`, and the discipline of using them rarely

**Next level:** Level 17 - Packages & Modules
- Organizing code into multiple packages
- Exported vs unexported identifiers at package scale
- go.mod, module paths, and third-party dependencies
- Designing clean package APIs (including which errors to export!)

With Level 16 done, you've completed the entire **OOP CONCEPTS** tier (Levels 12-16): pointers, structs, methods, interfaces, and error handling. Level 17 opens the **Packages & Modules** tier, where you'll organize everything you've built into real, multi-file projects. Congratulations - keep going! 🚀
