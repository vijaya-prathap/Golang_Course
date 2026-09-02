# Level 22: Quick Reference Card

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir myapp && cd myapp
go mod init myapp

# Create main.go
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "time"
)

func main() {
    ch := make(chan string)
    go func() {
        time.Sleep(50 * time.Millisecond)
        ch <- "hello"
    }()

    select {
    case msg := <-ch:
        fmt.Println("received:", msg)
    case <-time.After(200 * time.Millisecond):
        fmt.Println("timed out")
    }
}
EOF

# Run
go run main.go
```

---

## 📋 Basic select

```go
select {
case v := <-ch1:
    // ch1 was ready first
case v2 := <-ch2:
    // ch2 was ready first
}
// blocks until at least one case is ready
// if MULTIPLE are ready at once, Go picks ONE AT RANDOM
```

---

## 🎲 Random Tie-Breaking

```go
// If ch1 and ch2 are BOTH ready, this does NOT always pick ch1.
// Go picks uniformly at random among every ready case, every time.
select {
case <-ch1:
case <-ch2:
}
```

---

## 🚦 default: Non-Blocking select

```go
select {
case v := <-ch:
    use(v)
default:
    // runs immediately if ch had nothing ready
}

select {
case ch <- v:
    // sent - there was room / a receiver
default:
    // channel was full / had no receiver
}
```

---

## ⏱️ Timeout With time.After

```go
select {
case result := <-work:
    use(result)
case <-time.After(2 * time.Second):
    fmt.Println("timed out")
}
```

---

## 🛑 done-Channel Cancellation

```go
func worker(done <-chan struct{}, work <-chan int) {
    for {
        select {
        case <-done:
            return
        case w := <-work:
            process(w)
        }
    }
}
// close(done) tells every goroutine watching it to stop -
// precursor to the context package (a later level)
```

---

## ➡️ select With a Send Case

```go
select {
case out <- result:
    // delivered
case <-done:
    return // nobody's listening anymore
}
```

---

## 🔒 select on a Closed Channel

```go
close(ch)

select {
case v, ok := <-ch:
    // ok == false -> ALWAYS ready, returns zero value
    // WITHOUT an exit check, this becomes a CPU busy-loop
}
```

---

## ⏲️ time.After Leak: The Fix

```go
// ❌ Leaky - new timer every iteration
for {
    select {
    case v := <-ch:
        use(v)
    case <-time.After(d):
        onTimeout()
    }
}

// ✅ Reused - one timer for the whole loop
timer := time.NewTimer(d)
defer timer.Stop()
for {
    if !timer.Stop() {
        select {
        case <-timer.C:
        default:
        }
    }
    timer.Reset(d)
    select {
    case v := <-ch:
        use(v)
    case <-timer.C:
        onTimeout()
    }
}
```

---

## 🕳️ nil Channel: Disable a Case

```go
var a chan int = someChannel
// ...
select {
case v, ok := <-a:
    if !ok {
        a = nil // this case is never chosen again
        continue
    }
    use(v)
case v := <-b:
    use(v)
}
```

---

## 🧩 select + for = Multiplexer

```go
for a != nil || b != nil {
    select {
    case v, ok := <-a:
        if !ok { a = nil; continue }
        handle(v)
    case v, ok := <-b:
        if !ok { b = nil; continue }
        handle(v)
    }
}
// stops once every source is exhausted and disabled
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Forgetting default blocks | `select { case v := <-ch: ... }` alone | add `default:` for a non-blocking check |
| Assuming order = priority | believing the first case wins ties | ties are broken randomly, order is irrelevant |
| Leaking timers in a loop | `case <-time.After(d):` every iteration | `time.NewTimer(d)` once, `Stop()`/`Reset()` per pass |
| No way to stop the loop | `for { select { ... } }` with no exit case | always include a `done`, counter, or timeout case |
| Busy-loop on closed channel | reading a closed channel's case without checking `ok` | check `ok`, then `return`/`break`/set to `nil` |

---

## 🎓 Before Next Level

Can you:
- [ ] Write a select with multiple receive and send cases?
- [ ] Explain why select's tie-breaking is random, not ordered?
- [ ] Use default to make a select non-blocking?
- [ ] Implement a timeout with time.After, and explain the timer-leak gotcha?
- [ ] Build a done-channel cancellation pattern from scratch?
- [ ] Use a nil channel to disable a select case at runtime?

If YES → You're ready for Level 23!

---

## 📚 Next Level

Level 23: Mutex, WaitGroup & Atomic
- sync.Mutex and sync.RWMutex for protecting shared state directly
- sync.WaitGroup in depth for waiting on many goroutines
- The sync/atomic package for lock-free counters and flags

You've got select down! 💪
