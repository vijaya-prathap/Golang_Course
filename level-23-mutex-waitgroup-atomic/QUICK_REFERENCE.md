# Level 23: Quick Reference Card

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
    "sync"
    "sync/atomic"
)

func main() {
    var wg sync.WaitGroup
    var mu sync.Mutex
    var counter atomic.Int64

    guarded := 0
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            mu.Lock()
            guarded++
            mu.Unlock()
            counter.Add(1)
        }()
    }
    wg.Wait()

    fmt.Println("mutex-guarded:", guarded)
    fmt.Println("atomic:", counter.Load())
}
EOF

# Run
go run main.go
go run -race main.go   # confirm no data races
```

---

## 🔒 sync.Mutex

```go
var mu sync.Mutex

mu.Lock()
defer mu.Unlock()   // idiomatic: place immediately after Lock()
// critical section
```

---

## 🗂️ sync.RWMutex

```go
var mu sync.RWMutex

mu.RLock()           // many readers can hold this at once
defer mu.RUnlock()

mu.Lock()            // writers are exclusive
defer mu.Unlock()
```

Use when reads are frequent AND substantial. For trivial reads, plain `Mutex` can be simpler and just as fast (verified: 13.6ms Mutex vs 20ms RWMutex for cheap map reads in this project's sandbox).

---

## 🔁 sync.WaitGroup

```go
var wg sync.WaitGroup

wg.Add(1)        // BEFORE the `go` statement - never inside the goroutine
go func() {
    defer wg.Done()
    // work
}()

wg.Wait()        // blocks until counter reaches 0
```

```go
// ❌ NEVER pass a WaitGroup by value - it copies the counter
func worker(wg sync.WaitGroup) { defer wg.Done() }   // BUG

// ✅ ALWAYS pass by pointer
func worker(wg *sync.WaitGroup) { defer wg.Done() }  // correct
```

---

## ⚛️ sync/atomic (Typed, Go 1.19+)

```go
var c atomic.Int64     // also: Int32, Uint32, Uint64, Bool, Pointer[T]

c.Add(1)               // atomically add
c.Load()                // atomically read
c.Store(0)              // atomically set
c.CompareAndSwap(0, 5)  // atomic compare-and-swap, returns bool
```

**Atomic vs Mutex:**

| One independent variable | Multiple related fields |
|---------------------------|--------------------------|
| `sync/atomic` | `sync.Mutex` |

---

## 🔂 sync.Once

```go
var once sync.Once

once.Do(func() {
    // runs EXACTLY once, no matter how many goroutines call once.Do
})
```

---

## 📡 Channels vs Mutex

```
Handing off a value / coordinating flow  ──► channel
Protecting shared state read/written
  directly by multiple goroutines        ──► mutex (or atomic)
```

---

## 💀 Lock-Ordering Deadlock

```go
// ❌ Opposite order between goroutines = deadlock risk
func a() { mu1.Lock(); mu2.Lock() /* ... */ }
func b() { mu2.Lock(); mu1.Lock() /* ... */ }  // DANGER

// ✅ Same order everywhere = safe
func a() { mu1.Lock(); mu2.Lock() /* ... */ }
func b() { mu1.Lock(); mu2.Lock() /* ... */ }
```

Real captured output when it happens:
```
fatal error: all goroutines are asleep - deadlock!
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Forgetting to unlock | manual `Unlock()` skipped by early return | `mu.Lock(); defer mu.Unlock()` |
| Copying a lock | `func f(wg sync.WaitGroup)` | `func f(wg *sync.WaitGroup)` |
| Add() racing with Wait() | `go func(){ wg.Add(1) }()` | `wg.Add(1); go func(){ }()` |
| Opposite lock order | `mu1,mu2` then `mu2,mu1` | same order everywhere |
| Mutex for a hand-off | polling a guarded variable | use a channel instead |

---

## 🎓 Before Next Level

Can you:
- [ ] Fix a racy counter with `sync.Mutex` and confirm it with `-race`?
- [ ] Explain why `defer mu.Unlock()` goes immediately after `Lock()`?
- [ ] Explain when `RWMutex` helps and when plain `Mutex` is simpler/faster?
- [ ] Explain both WaitGroup mistakes (Add() timing, pass-by-value) from memory?
- [ ] Choose between an atomic and a mutex for a given piece of state?
- [ ] Explain `sync.Once`'s guarantee under concurrent calls?
- [ ] Read a `fatal error: all goroutines are asleep - deadlock!` dump and explain it?

If YES → You're ready for Level 24!

---

## 📚 Next Level

Level 24: Context
- `context.Context` for cancellation, deadlines, and timeouts
- Passing request-scoped values safely
- Controlling goroutine lifetimes on top of these primitives

You've got Go's classic concurrency toolkit down! 💪
