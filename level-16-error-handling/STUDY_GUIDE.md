# Level 16: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: Error Handling
```
Day 1:  The error interface, errors.New, fmt.Errorf
Day 2:  Return-early flow, adding context at each layer
Day 3:  Wrapping with %w, errors.Unwrap
Day 4:  errors.Is (sentinels) and errors.As (custom types)
Day 5:  Custom error types, the typed-nil trap
Day 6:  panic and recover
Day 7:  defer for cleanup, best practices, common mistakes
```

### Week 2: Practice & Application
```
Day 1:  Exercises 1-3 (values, context, wrapping)
Day 2:  Exercises 4-6 (Is, As, custom types)
Day 3:  Exercises 7-8 (typed-nil trap, panic/recover)
Day 4:  Exercises 9-10 (refactoring, comprehensive practice)
Day 5:  Bonus challenges
Day 6-7: Practice & review
```

---

## 🎯 The error Interface

```
type error interface {
    Error() string
}

Any type with an Error() string method IS an error - implicitly,
just like any other interface (Level 15).

    errors.New("msg")        -> *errors.errorString  (has Error() string)
    fmt.Errorf("msg %d", n)  -> *fmt.wrapError or *errors.errorString
    &MyError{...}            -> whatever YOU define, as long as it has Error()
```

```
func doSomething() (Result, error)
                              │
                              └─ nil = success, non-nil = failure

result, err := doSomething()
if err != nil {
    // handle it
}
// use result safely here
```

---

## 🏗️ Error Chain: Wrapping and Traversal

```
base := errors.New("connection refused")
mid  := fmt.Errorf("dialing database: %w", base)
top  := fmt.Errorf("starting server: %w", mid)

    top ───Unwrap()──▶ mid ───Unwrap()──▶ base ───Unwrap()──▶ nil
    │                   │                  │
    "starting server:   "dialing           "connection
     dialing database:   database:          refused"
     connection refused"  connection
                          refused"

Printing top shows the FULL chain as one string:
    starting server: dialing database: connection refused
```

### %w vs %v

```
fmt.Errorf("context: %w", err)   -> LINKS err into the chain
                                     (errors.Is / errors.As can find it)

fmt.Errorf("context: %v", err)   -> COPIES err's text only
                                     (errors.Unwrap returns nil - chain is broken)
```

### How errors.Is Walks the Chain

```
errors.Is(err, target)

    err ──▶ is err == target? ──▶ NO ──▶ Unwrap() ──▶ is THIS == target? ──▶ ...
                  │                                           │
                  YES                                        YES
                  │                                           │
                  ▼                                           ▼
                 true                                        true

Repeats until Unwrap() returns nil (chain exhausted -> false).

Verified (Exercise 4):
    err := loading invoice: order 42: not found   (wraps ErrNotFound twice)
    errors.Is(err, ErrNotFound)  -> true
    err == ErrNotFound           -> false   (only compares the OUTERMOST error)
```

### How errors.As Walks the Chain

```
errors.As(err, &target)

    Walks the SAME chain as errors.Is, but instead of comparing VALUES,
    checks each error's TYPE against target's type.

    First match found -> copies it into target, returns true
    No match anywhere  -> target untouched, returns false

Verified (Exercise 5):
    var httpErr *HTTPError
    errors.As(err, &httpErr)  -> true
    httpErr.StatusCode        -> 404   (fields are now readable)
```

---

## 🧭 Decision Tree: Sentinel vs Custom Type vs Ad-hoc

```
Does the caller need to REACT differently to this failure?
│
├─ NO - it will only ever be printed/logged
│   └─ Plain fmt.Errorf("...: %w", err)   (ad-hoc, no special type)
│
└─ YES - the caller needs to branch on it
    │
    ├─ Caller only needs to know WHICH KIND of failure occurred
    │   └─ Sentinel error:  var ErrNotFound = errors.New("not found")
    │      Caller checks:   errors.Is(err, ErrNotFound)
    │
    └─ Caller needs DATA about the failure (field name, status code, position)
        └─ Custom error type:  type ValidationError struct { Field, Reason string }
           Caller checks:      errors.As(err, &vErr); read vErr.Field
```

---

## ⚡ Decision Tree: panic vs Return an Error

```
Something has gone wrong. Now what?
│
├─ Is this a NORMAL, EXPECTED failure? (bad input, missing file,
│  network down, duplicate record, out-of-range value)
│   └─ RETURN AN ERROR
│      - This is 95%+ of all failures in real programs
│      - Caller decides what to do: retry, report, fall back
│
└─ Is this an "IMPOSSIBLE" state - a bug, a broken invariant,
   something that should never happen if the code is correct?
    └─ PANIC
       - Programmer errors: nil map write, index out of range,
         a switch missing a case that "can't" occur
       - Startup failures the program truly cannot run without
       - Library internals MAY panic internally, but should
         recover() at the public API boundary and return an error
```

```
Rule of thumb: if a user could trigger it by typing something silly,
it's an ERROR, not a panic.
```

---

## 🔁 panic / defer / recover Flow

```
func safeDivide(a, b int) (result int, err error) {
    defer func() {              ← ❶ registered BEFORE the risky call
        if r := recover(); r != nil {
            err = fmt.Errorf("recovered: %v", r)   ← ❹ converts panic to error
        }
    }()
    return a / b, nil            ← ❷ this PANICS if b == 0
}                                    ❸ panic unwinds straight to the deferred func

Timeline for safeDivide(10, 0):

  call safeDivide
      │
      ▼
  defer registered  ──────────────────────────────┐
      │                                            │
      ▼                                            │
  a / b panics: "integer divide by zero"           │
      │                                            │
      ▼ (unwind begins - stack pops toward defers)  │
  deferred closure runs ◀──────────────────────────┘
      │
      ▼
  recover() catches the panic value, sets err
      │
      ▼
  safeDivide returns NORMALLY: result=0, err="recovered: ..."
      │
      ▼
  caller's code continues - the program never crashed
```

**Key rule:** `recover()` only stops a panic when called **directly inside a deferred function**. Anywhere else, it's a no-op that returns `nil`.

### Uncaught Panic (no recover) - Verified Shape

```
before the panic
deferred: I still run while the panic unwinds     ← defers STILL run
panic: mustPositive called with -5                ← panic: <value>

goroutine 1 [running]:                            ← always this header
main.mustPositive(...)
        /path/to/main.go:7
main.main()
        /path/to/main.go:15 +0xd0
exit status 2                                     ← process exits with code 2
```

---

## 🚨 The Typed-Nil Error Trap

```
type MyError struct{}
func (e *MyError) Error() string { return "boom" }

func broken() *MyError { return nil }     ← returns a CONCRETE pointer type

var err error = broken()

    error interface value:
    ┌─────────────┬───────────┐
    │  type       │  value    │
    ├─────────────┼───────────┤
    │  *MyError   │  nil      │   ← non-nil TYPE, nil VALUE
    └─────────────┴───────────┘

    err == nil  →  FALSE   (only true when BOTH type and value are nil)

THE FIX: declare the function's return type as the error INTERFACE,
never a concrete pointer type.

    func broken() *MyError   ❌  →  func broken() error   ✅

    With `error` as the signature, `return nil` is an UNTYPED nil,
    so the interface value becomes (nil, nil) - and err == nil is TRUE.
```

Verified:

```
=== The Trap ===             === The Fix ===
err's type: *main.ConfigError    err's type: <nil>
err == nil: false                err == nil: true
```

---

## 🧩 defer + Error Handling for Cleanup

```
func process() error {
    r, err := open("data.txt")
    if err != nil {                  ← ❶ check error FIRST
        return fmt.Errorf("opening resource: %w", err)
    }
    defer r.Close()                  ← ❷ defer AFTER a successful acquire

    // ... work that might return an error from several places ...
    return nil                       ← ❸ Close() still runs, no matter which
                                          return statement (or panic) fires
}
```

Rule: **check-then-defer**. Deferring cleanup before checking the error risks calling a method on something that was never successfully created.

---

## 🚨 Common Mistakes

### Mistake 1: Ignoring Errors With _

```go
value, _ := strconv.Atoi(input)   // ❌ failure vanishes silently
```

### Mistake 2: Shadowing err Inside if

```go
var err error
if cond {
    data, err := readCache()   // ❌ := makes a NEW err, shadows the outer one
}
if err != nil { ... }          // outer err is still nil!
```

### Mistake 3: Returning Typed Nils

```go
func validate() *ValidationError { return nil }  // ❌ non-nil error interface
func validate() error            { return nil }  // ✅ genuinely nil
```

### Mistake 4: Using == Instead of errors.Is on Wrapped Errors

```go
if err == ErrNotFound { ... }          // ❌ false once wrapped
if errors.Is(err, ErrNotFound) { ... } // ✅ true through any depth
```

### Mistake 5: Panicking for Normal Failures

```go
if !ok { panic("user not found") }              // ❌ a lookup miss is not a bug
if !ok { return fmt.Errorf("user not found") }  // ✅ expected failure -> error
```

---

## 📈 Progression Summary

### Understanding Level 16

Level 16 makes explicit what you've been doing informally since Level 3:

1. **The interface** - `error` is one method, implemented implicitly (built on Level 15)
2. **Creating** - `errors.New` for fixed text, `fmt.Errorf` for formatted text
3. **Flow** - return early, add context, keep the happy path flat
4. **Wrapping** - `%w` builds a chain; `errors.Unwrap` walks it one link at a time
5. **Inspecting** - `errors.Is` for identity (sentinels), `errors.As` for type + data (custom types)
6. **Designing** - sentinel vs custom type vs ad-hoc, chosen by what the caller needs
7. **The trap** - typed nils, and why `error` must be the declared return type
8. **The escape hatch** - `panic`/`recover`, used rarely, mostly for programmer bugs

### Prerequisites for Level 17

Before moving to Level 17 (Packages & Modules), you need:

- ✅ Comfortable creating and checking errors with `if err != nil`
- ✅ Can wrap errors with `%w` and explain what `errors.Unwrap` does
- ✅ Can use `errors.Is` for sentinels and `errors.As` for custom types
- ✅ Can explain the typed-nil error trap and always return `error`, not a concrete type
- ✅ Understand when to `panic`/`recover` versus when to return an error
- ✅ Completed 8+ exercises

### Ready for Level 17?

Level 17 teaches how to organize the code you've been writing into real, multi-file, multi-package projects:
- Splitting code across packages
- Exported vs unexported identifiers at package scale
- `go.mod`, module paths, and third-party dependencies
- Designing clean package APIs - including which sentinel errors and custom error types to export

---

## ✅ Checklist Before Level 17

- [ ] Can explain the `error` interface in one sentence
- [ ] Can create errors with both `errors.New` and `fmt.Errorf`
- [ ] Can wrap an error with `%w` and unwrap it with `errors.Unwrap`
- [ ] Can define a sentinel error and match it through a wrapped chain with `errors.Is`
- [ ] Can define a custom error type and extract it with `errors.As`
- [ ] Can reproduce the typed-nil error trap and explain the fix
- [ ] Can use `recover()` inside a deferred function to convert a panic into an error
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The Foundation
`error` is just an interface with one method - `Error() string`. Nothing more.

### The Flow
Return early, add context at each layer with `fmt.Errorf("...: %w", err)`, keep the happy path flat.

### The Chain
`%w` links errors together; `errors.Is` finds a sentinel anywhere in the chain, `errors.As` extracts a typed error anywhere in the chain. `==` and plain type assertions only see the outermost link.

### The Trap
Never declare a concrete error pointer as a return type - only `error`. A typed nil inside an interface is not a nil interface.

### The Escape Hatch
`panic` is for bugs and impossible states, not expected failures. `recover()` only works directly inside a deferred function.

---

## 📚 Next Level

Level 17: Packages & Modules
- Organizing code into multiple packages
- Exported vs unexported identifiers at package scale
- go.mod, module paths, and third-party dependencies
- Designing clean package APIs

You've mastered error handling - and completed the whole OOP CONCEPTS tier! Keep going! 🚀
