# Level 20: Quick Reference Card

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
)

func main() {
    var wg sync.WaitGroup

    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            fmt.Println("working")
        }()
    }

    wg.Wait()
    fmt.Println("all done")
}
EOF

# Run
go run main.go

# Run with the race detector
go run -race main.go
```

---

## 📋 Launching a Goroutine

```go
go someFunc()          // launches, returns immediately - doesn't wait
go func() { ... }()    // anonymous goroutine
```

⚠️ **`main()` returning kills every goroutine still running.** Never assume a background goroutine "had enough time" - always synchronize explicitly.

---

## 🔁 sync.WaitGroup

```go
var wg sync.WaitGroup

wg.Add(1)               // BEFORE launching - increments the counter
go func() {
    defer wg.Done()      // decrements the counter when this goroutine ends
    // work
}()

wg.Wait()                // blocks until the counter reaches 0
```

| Method | When to call it |
|--------|------------------|
| `Add(n)` | Before launching the goroutine(s) it counts |
| `Done()` | Inside the goroutine, via `defer`, when it finishes |
| `Wait()` | Wherever you need to block until all are done |

Always pass `*sync.WaitGroup` (a pointer) - never copy a `WaitGroup` after first use.

---

## 🔒 sync.Mutex (Preview - full coverage: Level 23)

```go
var mu sync.Mutex

mu.Lock()
sharedValue++    // critical section - only one goroutine at a time
mu.Unlock()
```

Any variable read **and** written by more than one goroutine needs a mutex around **every** access, not just some.

---

## 🎁 Loop Variable Capture (Go 1.22+, this project: go1.26.5)

```go
for i := 0; i < 5; i++ {
    go func() {
        fmt.Println(i)   // each goroutine gets its OWN i - safe on Go 1.22+
    }()
}
```

Pre-1.22 code you find online may show `go func(i int) { ... }(i)` to force a copy - harmless but unnecessary on this project's Go version. See Level 11's closures-in-loops section for the full history.

---

## 🚨 Data Races

```go
counter++                     // ❌ NOT atomic: read, add, write - a race
                               //    if done from multiple goroutines

mu.Lock(); counter++; mu.Unlock()   // ✅ fixed
```

```bash
go run -race main.go     # catches races in the act; requires cgo
```

A data race = concurrent unsynchronized access to shared memory, with at least one write. Output from a racy program is **genuinely non-deterministic** - don't expect the same wrong answer twice.

---

## 💀 Goroutine Leaks

```go
// ❌ blocks forever if nobody ever reads ch
go func() { ch <- result }()

// ✅ always give it a way out
go func() {
    select {
    case ch <- result:
    case <-time.After(timeout):
    }
}()
```

A leaked goroutine is never garbage collected - it just sits blocked, consuming memory, for the life of the process.

---

## 🧩 Concurrent-Safe Result Collection

```go
var wg sync.WaitGroup
var mu sync.Mutex
var results []int

for _, n := range inputs {
    wg.Add(1)
    go func(n int) {
        defer wg.Done()
        r := compute(n)
        mu.Lock()
        results = append(results, r)
        mu.Unlock()
    }(n)
}
wg.Wait()
sort.Ints(results) // make the CHECKED result deterministic
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Forgetting to wait | `go fmt.Println("hi")` then return | `wg.Add`/`Done`/`Wait` around it |
| Guessing at timing | `time.Sleep(100ms)` and hope | `sync.WaitGroup` |
| Assuming pre-1.22 sharing | `i := i` workaround "just in case" | Trust per-iteration capture on Go 1.22+ |
| Unsynchronized shared state | `counter++` from many goroutines | `mu.Lock(); counter++; mu.Unlock()` |
| Goroutine with no exit | `ch <- v` with no reader, ever | Timeout / done channel / context |
| Goroutine for trivial work | `go func() { total += a + b }()` | Just `total += a + b` |

---

## 🎓 Before Next Level

Can you:
- [ ] Explain why a goroutine's output isn't guaranteed without synchronization?
- [ ] Use `sync.WaitGroup` to wait for N goroutines?
- [ ] Explain why loop variable capture is safe on this project's Go version (1.26.5)?
- [ ] Reproduce a data race and explain what `-race` reports?
- [ ] Fix a race with a minimal `sync.Mutex`?
- [ ] Explain a goroutine leak and why it's never garbage collected?
- [ ] Say when a goroutine is (and isn't) worth using?

If YES → You're ready for Level 21!

---

## 📚 Next Level

Level 21: Channels
- Sending and receiving values between goroutines
- Buffered vs. unbuffered channels
- `select` for waiting on multiple channels

You've got goroutines down! 💪
