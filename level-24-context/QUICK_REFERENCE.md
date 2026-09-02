# Level 24: Quick Reference Card

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir myapp && cd myapp
go mod init myapp

# Create main.go
cat > main.go << 'EOF'
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
    defer cancel()

    select {
    case <-ctx.Done():
        fmt.Println("stopped:", ctx.Err())
    case <-time.After(500 * time.Millisecond):
        fmt.Println("this won't print - the timeout wins")
    }
}
EOF

# Run
go run main.go
```

---

## 📋 Creating Contexts

```go
ctx := context.Background()              // real root - use at true top-level entry points
ctx := context.TODO()                    // placeholder - "context plumbing still needed here"

ctx, cancel := context.WithCancel(parent)                    // manual cancellation
ctx, cancel := context.WithTimeout(parent, 5*time.Second)    // auto-cancels after a duration
ctx, cancel := context.WithDeadline(parent, someTime)        // auto-cancels at an absolute time

defer cancel() // ALWAYS - even for WithTimeout/WithDeadline
```

---

## 🔍 Inspecting a Context

```go
<-ctx.Done()          // blocks until canceled/timed out; channel closes, doesn't send a value
ctx.Err()             // nil while active; context.Canceled or context.DeadlineExceeded once done
d, ok := ctx.Deadline() // the absolute deadline, if one was set
ctx.Value(key)        // request-scoped data (see WithValue below)
```

| State | `ctx.Err()` |
|---|---|
| Still active | `nil` |
| `cancel()` was called | `context.Canceled` |
| Timeout/deadline elapsed | `context.DeadlineExceeded` |

```go
errors.Is(ctx.Err(), context.Canceled)          // idiomatic check
errors.Is(ctx.Err(), context.DeadlineExceeded)  // idiomatic check
```

---

## 🔀 The select Pattern

```go
select {
case <-ctx.Done():
    return ctx.Err()
case result := <-workCh:
    // use result
}
```

Apply this at **every** blocking point in long-running work, not just the outermost call - including when a goroutine tries to *send* a result, so it doesn't block forever on a caller who already gave up.

---

## 🗝️ context.WithValue

```go
type requestIDKey struct{} // unexported, unique type

ctx = context.WithValue(ctx, requestIDKey{}, "req-123")
id, ok := ctx.Value(requestIDKey{}).(string)
```

| Key type | Safe? |
|---|---|
| `string` / other basic type | ❌ collides silently with any other code using the same string |
| unexported `struct{}` type | ✅ only your package can ever construct this key |

**Use for:** trace IDs, request IDs, cross-cutting metadata.
**Never use for:** anything your function needs to do its job - those are explicit parameters.

---

## 💧 defer cancel() - Always

```go
// ❌ WRONG - go vet catches this
ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

// ✅ RIGHT
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

Real captured `go vet` warning for the wrong version:
```
main.go:19:10: the cancel function returned by context.WithTimeout should be called, not discarded, to avoid a context leak
```

What leaks: the context object and its internal timer stay alive for the *entire* requested duration/deadline, even after your work is done - harmless in a short CLI run, real accumulating overhead in a long-running server handling many requests.

---

## 🌐 Fan-Out Skeleton

```go
ctx, cancel := context.WithCancel(context.Background())
var wg sync.WaitGroup

for i := 0; i < numWorkers; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        for {
            select {
            case <-ctx.Done():
                return
            case <-time.After(interval):
                // do work
            }
        }
    }()
}

cancel()   // stops EVERY worker at once - they all share the same ctx
wg.Wait()
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Forgotten cancel | `ctx, _ := context.WithTimeout(...)` | `ctx, cancel := ...; defer cancel()` |
| Overusing Value | business params hidden in `ctx.Value(...)` | pass them as real function arguments |
| Background() mid-chain | `ctx := context.Background()` inside a function that got a real `ctx` | forward the caller's `ctx` |
| Ignoring Done()/Err() | a function takes `ctx` but never checks it | select on `<-ctx.Done()` in every blocking step |
| String context keys | `context.WithValue(ctx, "id", v)` | `context.WithValue(ctx, myKeyType{}, v)` |
| Storing ctx in a struct | `type S struct { ctx context.Context }` | pass `ctx` as a method parameter |

---

## 🎓 Before Next Level

Can you:
- [ ] Explain Background() vs. TODO() and when to use each?
- [ ] Cancel a real goroutine with WithCancel and verify it stopped?
- [ ] Use WithTimeout/WithDeadline and read the real elapsed time?
- [ ] Distinguish context.Canceled from context.DeadlineExceeded?
- [ ] Write the ctx.Done()-select pattern from memory?
- [ ] Explain why context.Value needs a typed, unexported key?
- [ ] Explain exactly what defer cancel() prevents?

If YES → You're ready for Level 25!

---

## 📚 Next Level

Level 25: Generics
- Type parameters for functions and types
- Type constraints and the comparable/any built-in constraints
- Writing generic data structures and algorithms

You've got context down! 💪
