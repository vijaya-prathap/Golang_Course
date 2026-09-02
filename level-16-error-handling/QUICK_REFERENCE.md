# Level 16: Quick Reference Card

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir myapp && cd myapp
go mod init myapp

# Create main.go
cat > main.go << 'EOF'
package main

import (
    "errors"
    "fmt"
)

var ErrNotFound = errors.New("not found")

func find(id int) error {
    if id != 1 {
        return fmt.Errorf("looking up id %d: %w", id, ErrNotFound)
    }
    return nil
}

func main() {
    err := find(42)
    if errors.Is(err, ErrNotFound) {
        fmt.Println("not found:", err)
    }
}
EOF

# Run
go run main.go
```

---

## 📋 The error Interface

```go
type error interface {
    Error() string
}
```

Any type with `Error() string` satisfies it - implicitly, like every interface (Level 15).

---

## ✍️ Creating Errors

```go
errors.New("division by zero")            // fixed message
fmt.Errorf("user %d not found", id)        // formatted message
fmt.Errorf("loading config: %w", err)      // WRAPPED - keeps err in the chain
```

**Message style:** lowercase, no trailing period, include relevant values.

---

## 🔁 The Fundamental Flow

```go
result, err := doSomething()
if err != nil {
    return fmt.Errorf("doing something: %w", err)  // add context, return early
}
// happy path continues here, un-indented
```

---

## 🔗 Wrapping and Unwrapping

```go
wrapped := fmt.Errorf("context: %w", err)  // %w links err into the chain
errors.Unwrap(wrapped)                     // -> err (one layer)
errors.Unwrap(err)                         // -> nil if err has no wrapped cause

fmt.Errorf("context: %v", err)             // %v copies TEXT ONLY - breaks the chain!
```

---

## 🎯 errors.Is - Sentinel Matching

```go
var ErrNotFound = errors.New("not found")   // define once, package-level

err := fmt.Errorf("loading invoice: %w", fmt.Errorf("order 42: %w", ErrNotFound))

errors.Is(err, ErrNotFound)   // true  - finds it ANYWHERE in the chain
err == ErrNotFound            // false - only checks the outermost error
```

**Rule: use `errors.Is`, never `==`, once wrapping is possible.**

---

## 🏷️ errors.As - Extracting a Typed Error

```go
type ValidationError struct {
    Field, Reason string
}
func (e *ValidationError) Error() string {
    return fmt.Sprintf("invalid %s: %s", e.Field, e.Reason)
}

var vErr *ValidationError
if errors.As(err, &vErr) {           // pass a POINTER to the target
    fmt.Println(vErr.Field, vErr.Reason)
}
```

---

## 🧭 Sentinel vs Custom Type vs Ad-hoc

```go
var ErrX = errors.New("...")     // caller only needs to recognize the KIND
type XError struct{ ... }        // caller needs DATA about the failure
fmt.Errorf("...")                // caller will only ever print it
```

---

## 🚨 The Typed-Nil Trap

```go
func broken() *MyError { return nil }  // ❌ concrete pointer return type
var err error = broken()
err == nil   // false! interface holds (*MyError, nil) - non-nil interface

func fixed() error { return nil }      // ✅ interface return type
var err error = fixed()
err == nil   // true - untyped nil, interface is (nil, nil)
```

**Always declare `error`, never a concrete error pointer, as a return type.**

---

## 💥 panic and recover

```go
panic("message")   // stops execution, runs defers, unwinds, crashes if uncaught

func safe() (err error) {
    defer func() {
        if r := recover(); r != nil {      // ONLY works directly inside a defer
            err = fmt.Errorf("recovered: %v", r)
        }
    }()
    riskyCall()
    return nil
}
```

Uncaught panic shape:

```
panic: <message>

goroutine 1 [running]:
main.main()
        /path/to/main.go:N +0xNN
exit status 2
```

**panic for:** bugs, impossible states, unrecoverable startup failures.
**error for:** everything else (95%+ of failures).

---

## 🧹 defer for Cleanup

```go
r, err := open("data.txt")
if err != nil {
    return fmt.Errorf("opening resource: %w", err)  // check error FIRST
}
defer r.Close()   // defer immediately after a successful acquire
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Ignoring errors | `value, _ := strconv.Atoi(s)` | check and handle `err` |
| Shadowing err | `data, err := f()` inside an `if` | plain `=`, not `:=` |
| Typed nil | `func f() *MyError` | `func f() error` |
| == on wrapped errors | `err == ErrX` | `errors.Is(err, ErrX)` |
| Panic for normal failures | `panic("not found")` | `return fmt.Errorf(...)` |
| %v breaks the chain | `fmt.Errorf("...: %v", err)` | `fmt.Errorf("...: %w", err)` |
| Log AND return | log the error, then also return it | return with context OR handle - not both |

---

## 🎓 Before Next Level

Can you:
- [ ] Explain the error interface and why Go returns errors instead of throwing?
- [ ] Wrap an error with %w and unwrap it with errors.Unwrap?
- [ ] Use errors.Is to match a sentinel through a wrapped chain?
- [ ] Use errors.As to extract a custom error type and read its fields?
- [ ] Reproduce and fix the typed-nil error trap?
- [ ] Use recover() inside a deferred function to convert a panic into an error?

If YES → You're ready for Level 17!

---

## 📚 Next Level

Level 17: Packages & Modules
- Organizing code into multiple packages
- Exported vs unexported identifiers at package scale
- go.mod, module paths, and third-party dependencies

You've got error handling down - and finished OOP CONCEPTS! 💪
