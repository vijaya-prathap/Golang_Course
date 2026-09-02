# Level 21: Quick Reference Card

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir myapp && cd myapp
go mod init myapp

# Create main.go
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    ch := make(chan string)

    go func() {
        ch <- "hello from a goroutine"
    }()

    msg := <-ch
    fmt.Println(msg)
}
EOF

# Run
go run main.go
```

---

## 📋 Creating Channels

```go
ch := make(chan int)      // unbuffered - send/receive both block until the other side is ready
ch := make(chan int, 5)   // buffered, capacity 5 - send blocks only once 5 items are waiting
```

---

## 🔁 Send, Receive, Close

```go
ch <- v          // send v (blocks per the rules above)
v := <-ch        // receive (single-value form)
v, ok := <-ch    // receive with comma-ok
close(ch)        // signal "no more sends" - never send after this
```

| Operation | On a closed, drained channel |
|---|---|
| `v := <-ch` | returns the zero value immediately, no block |
| `v, ok := <-ch` | `v` = zero value, `ok` = `false` |
| `ch <- x` | **panics**: `send on closed channel` |
| `close(ch)` again | **panics**: `close of closed channel` |

---

## 🎁 range Over a Channel

```go
for v := range ch {
    // runs once per value, in send order
    // exits automatically once ch is closed AND drained
}
```

---

## ➡️⬅️ Direction Types

```go
func sendOnly(ch chan<- int)   { ch <- 1 }        // may only send
func receiveOnly(ch <-chan int) int { return <-ch } // may only receive

// misuse is a COMPILE ERROR:
// ./main.go:4:12: invalid operation: cannot receive from
// send-only channel chan<- int ch (variable of type chan<- int)
```

---

## 💥 Deadlock (real, captured output)

```go
ch := make(chan int)
ch <- 1 // nobody will ever receive - deadlock
```

```
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan send]:
main.main()
        /path/to/main.go:8 +0x68
exit status 2
```

Terminates immediately - never hangs.

---

## 💣 Panics (real, captured output)

```go
ch := make(chan int)
close(ch)
ch <- 1        // panic: send on closed channel
```

```go
ch := make(chan int)
close(ch)
close(ch)      // panic: close of closed channel
```

---

## 🏭 Worker Pool Skeleton

```go
func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()
    for j := range jobs {
        results <- j * j
    }
}

jobs := make(chan int, numJobs)
results := make(chan int, numJobs)

var wg sync.WaitGroup
for w := 1; w <= numWorkers; w++ {
    wg.Add(1)
    go worker(w, jobs, results, &wg)
}
for j := 1; j <= numJobs; j++ { jobs <- j }
close(jobs)

wg.Wait()
close(results) // safe only after every worker is confirmed done
```

---

## 🧵 Pipeline Skeleton

```go
func stage(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for v := range in {
            out <- transform(v)
        }
    }()
    return out
}

result := stage3(stage2(stage1(source)))
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Deadlock | send/receive with no possible counterpart, ever | make sure a concurrent goroutine (or buffer) is on the other side |
| Send on closed | `close(ch); ch <- 1` | only the sender closes, only when truly done |
| Double close | `close(ch); close(ch)` | close exactly once (guard with `sync.Once` if unsure) |
| Forgotten close | producer never closes; `range` consumer hangs forever | `defer close(ch)` in the producer |
| Wrong side closes | receiver calls `close` | only the sender closes |
| nil channel | assuming it behaves like a closed one | nil channel blocks forever on send AND receive |

---

## 🎓 Before Next Level

Can you:
- [ ] Explain exactly when an unbuffered vs. buffered send blocks?
- [ ] Use comma-ok to tell a real zero value from a closed channel?
- [ ] Write a producer/consumer pair using range and close correctly?
- [ ] Write a function with a `chan<-` or `<-chan` parameter?
- [ ] Explain why Go crashes on a deadlock instead of hanging?
- [ ] Build a worker pool with a fixed number of workers?

If YES → You're ready for Level 22!

---

## 📚 Next Level

Level 22: Select
- Waiting on multiple channels at once with `select`
- `default` cases for non-blocking operations
- Timeouts with `time.After`

You've got channels down! 💪
