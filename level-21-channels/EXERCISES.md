# Level 21: Channels - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

**A note on ordering:** channel-based programs involve multiple goroutines, so in general their output order can vary between runs. Every exercise below is deliberately designed so its output is deterministic and checkable - either because only one goroutine ever prints, or because a channel operation forces one side to wait for the other before continuing. Where an exercise's *internal* execution order (e.g. which worker picks up which job) is genuinely non-deterministic, it's called out explicitly, and the exercise is designed so the final printed result is still exactly checkable.

---

## Exercise 1: Unbuffered Channel Hand-off

**Objective:** Send a value between two goroutines over an unbuffered channel and observe real blocking

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level21-exercise1
cd ~/projects/level21-exercise1
go mod init level21.example/exercise1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "time"
)

func main() {
    ch := make(chan string)

    go func() {
        time.Sleep(100 * time.Millisecond) // simulate the worker doing some work first
        ch <- "hello from the worker goroutine"
    }()

    fmt.Println("main: blocked, waiting to receive from the unbuffered channel...")
    msg := <-ch
    fmt.Println("main: received:", msg)
    fmt.Println("main: done")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
main: blocked, waiting to receive from the unbuffered channel...
main: received: hello from the worker goroutine
main: done
```

This output is exactly deterministic: `main` prints its own "blocked, waiting" line before it ever touches the channel, and it cannot print anything after that until the goroutine's send completes 100ms later.

**Learning Objectives:**
- ✅ Create an unbuffered channel with `make(chan string)`
- ✅ Send with `ch <- v` and receive with `v := <-ch`
- ✅ Observe that an unbuffered receive genuinely blocks until a sender is ready

---

## Exercise 2: Buffered vs. Unbuffered Blocking Behavior

**Objective:** Contrast a buffered channel's non-blocking sends with an unbuffered channel's blocking send

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level21-exercise2
cd ~/projects/level21-exercise2
go mod init level21.example/exercise2
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== Buffered Channel (capacity 3) ===")
    buffered := make(chan int, 3)

    fmt.Println("sending 1 (buffer has room, does not block)")
    buffered <- 1
    fmt.Println("sending 2 (buffer has room, does not block)")
    buffered <- 2
    fmt.Println("sending 3 (buffer has room, does not block)")
    buffered <- 3
    fmt.Printf("all 3 sends completed with no receiver at all. len=%d cap=%d\n", len(buffered), cap(buffered))

    fmt.Println("\ndraining the buffer:")
    fmt.Println(<-buffered)
    fmt.Println(<-buffered)
    fmt.Println(<-buffered)

    fmt.Println("\n=== Unbuffered Channel (capacity 0) ===")
    unbuffered := make(chan int)
    done := make(chan bool)

    fmt.Println("sender: about to send on the unbuffered channel (this blocks until a receiver is ready)")
    go func() {
        v := <-unbuffered
        fmt.Println("receiver: got", v)
        done <- true
    }()
    unbuffered <- 42
    <-done
    fmt.Println("sender: confirmed - the receiver above had to be ready before this line could run")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Buffered Channel (capacity 3) ===
sending 1 (buffer has room, does not block)
sending 2 (buffer has room, does not block)
sending 3 (buffer has room, does not block)
all 3 sends completed with no receiver at all. len=3 cap=3

draining the buffer:
1
2
3

=== Unbuffered Channel (capacity 0) ===
sender: about to send on the unbuffered channel (this blocks until a receiver is ready)
receiver: got 42
sender: confirmed - the receiver above had to be ready before this line could run
```

All three buffered sends complete instantly with **nothing** receiving yet - proof the buffer, not a receiver, is absorbing them. The unbuffered send, by contrast, cannot complete until the goroutine's receive is ready; main's final print is gated behind the `done` channel, so it can only run after the receiver's own print - this is fully deterministic every run.

**Learning Objectives:**
- ✅ Create a buffered channel with `make(chan int, n)`
- ✅ See that buffered sends don't block until the buffer is full
- ✅ Confirm an unbuffered send still blocks until a receiver is ready
- ✅ Use `len()`/`cap()` to inspect a channel's buffer

---

## Exercise 3: Closing Channels and the comma-ok Zero Value

**Objective:** Distinguish a real zero value from "the channel is closed" using comma-ok

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level21-exercise3
cd ~/projects/level21-exercise3
go mod init level21.example/exercise3
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    ch := make(chan int, 2)
    ch <- 10
    ch <- 20
    close(ch)

    v1, ok1 := <-ch
    fmt.Println("v1:", v1, "ok1:", ok1)

    v2, ok2 := <-ch
    fmt.Println("v2:", v2, "ok2:", ok2)

    v3, ok3 := <-ch
    fmt.Println("v3 (channel is drained AND closed):", v3, "ok3:", ok3)

    v4 := <-ch
    fmt.Println("v4 (single-value form on a closed, drained channel):", v4)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
v1: 10 ok1: true
v2: 20 ok2: true
v3 (channel is drained AND closed): 0 ok3: false
v4 (single-value form on a closed, drained channel): 0
```

**Learning Objectives:**
- ✅ Close a channel with `close(ch)`
- ✅ Use `v, ok := <-ch` to distinguish a real value from a closed/drained channel
- ✅ See that a closed, drained channel keeps returning the zero value forever

---

## Exercise 4: Ranging Over a Channel Until Closed

**Objective:** Use `for v := range ch` to consume values until a producer closes the channel

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level21-exercise4
cd ~/projects/level21-exercise4
go mod init level21.example/exercise4
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func produce(ch chan<- int, n int) {
    for i := 1; i <= n; i++ {
        ch <- i
    }
    close(ch)
}

func main() {
    ch := make(chan int)
    go produce(ch, 5)

    for v := range ch {
        fmt.Println("received:", v)
    }
    fmt.Println("channel closed, range loop exited")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
received: 1
received: 2
received: 3
received: 4
received: 5
channel closed, range loop exited
```

With a single producer and a single consumer, values always arrive in send order - this output is exactly deterministic every run.

**Learning Objectives:**
- ✅ Use `for v := range ch` to consume a channel
- ✅ See the range loop terminate automatically when the producer calls `close`
- ✅ Confirm send order is preserved for a single producer/single consumer

---

## Exercise 5: Channel Direction Types and a Real Compile Error

**Objective:** Restrict channel parameters to send-only/receive-only, then trigger the real compile error from misusing one

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level21-exercise5
cd ~/projects/level21-exercise5
go mod init level21.example/exercise5
```

2. Create the working version, `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func sendOnly(ch chan<- int) {
    ch <- 100
}

func receiveOnly(ch <-chan int) int {
    return <-ch
}

func main() {
    ch := make(chan int, 1)
    sendOnly(ch)
    fmt.Println("received:", receiveOnly(ch))
}
EOF
```

3. Run it:

```bash
go run main.go
```

**Expected Output:**

```
received: 100
```

4. Now see the compiler reject a direction violation. Replace `main.go`:

```bash
cat > main.go << 'EOF'
package main

func sendOnly(ch chan<- int) {
    v := <-ch // trying to receive from a send-only channel
    _ = v
}

func main() {
    ch := make(chan int, 1)
    sendOnly(ch)
}
EOF
```

5. Try to build it:

```bash
go build main.go
```

**Expected Output (a real compile error, captured verbatim):**

```
# command-line-arguments
./main.go:4:12: invalid operation: cannot receive from send-only channel chan<- int ch (variable of type chan<- int)
```

This is caught before the program ever runs. `go run main.go` produces the identical error and also never executes.

**Learning Objectives:**
- ✅ Declare send-only (`chan<-`) and receive-only (`<-chan`) function parameters
- ✅ Understand that a bidirectional channel is automatically restricted when passed into such a parameter
- ✅ See the compiler reject a direction violation at build time, not runtime

---

## Exercise 6: A Real Deadlock

**Objective:** Trigger and capture the exact fatal error Go produces for a deadlock

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level21-exercise6
cd ~/projects/level21-exercise6
go mod init level21.example/exercise6
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("about to send on an unbuffered channel with no receiver...")
    ch := make(chan int)
    ch <- 42 // nobody will ever receive this - deadlock
    fmt.Println("this line never runs")
}
EOF
```

3. Run it:

```bash
go run main.go
```

**Expected Output (captured verbatim from a real run):**

```
about to send on an unbuffered channel with no receiver...
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan send]:
main.main()
        /path/to/main.go:8 +0x68
exit status 2
```

The exact file path and offset (`+0x68`) will differ on your machine, but the `fatal error: all goroutines are asleep - deadlock!` line, the `[chan send]` state, and `exit status 2` are exactly what you'll see. Notice the program **terminates immediately** - a Go deadlock is a fast, safe crash, not an infinite hang.

**Learning Objectives:**
- ✅ Understand exactly what causes a deadlock: a send (or receive) with no possible counterpart, ever
- ✅ Recognize the real `fatal error: all goroutines are asleep - deadlock!` output
- ✅ Confirm Go's deadlock detector terminates the program rather than hanging it

---

## Exercise 7: Two Real Panics - Send on Closed, Close Twice

**Objective:** Trigger and capture the exact panics from sending on a closed channel and from closing a channel twice

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level21-exercise7
cd ~/projects/level21-exercise7
go mod init level21.example/exercise7
```

2. First, create `main.go` for the "send on closed channel" panic:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    ch := make(chan int)
    close(ch)
    fmt.Println("channel closed successfully")
    ch <- 1 // panic: send on closed channel
    fmt.Println("this line never runs")
}
EOF
```

3. Run it:

```bash
go run main.go
```

**Expected Output (captured verbatim):**

```
channel closed successfully
panic: send on closed channel

goroutine 1 [running]:
main.main()
        /path/to/main.go:9 +0x74
exit status 2
```

4. Now replace `main.go` to trigger the "close of closed channel" panic instead:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    ch := make(chan int)
    close(ch)
    fmt.Println("channel closed successfully (first time)")
    close(ch) // panic: close of closed channel
    fmt.Println("this line never runs")
}
EOF
```

5. Run it again:

```bash
go run main.go
```

**Expected Output (captured verbatim):**

```
channel closed successfully (first time)
panic: close of closed channel

goroutine 1 [running]:
main.main()
        /path/to/main.go:9 +0x6c
exit status 2
```

As with the deadlock, the file path and offset will differ on your machine, but the `panic:` line, goroutine trace shape, and `exit status 2` are exactly what you'll see in both cases.

**Learning Objectives:**
- ✅ Reproduce the real `panic: send on closed channel`
- ✅ Reproduce the real `panic: close of closed channel`
- ✅ Understand why both are panics (loud, immediate) rather than silent no-ops

---

## Exercise 8: The Worker Pool Pattern

**Objective:** Build a fixed-size pool of workers that consume jobs and produce results

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level21-exercise8
cd ~/projects/level21-exercise8
go mod init level21.example/exercise8
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "sort"
    "sync"
)

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()
    for j := range jobs {
        results <- j * j
    }
    _ = id
}

func main() {
    const numWorkers = 3
    const numJobs = 9

    jobs := make(chan int, numJobs)
    results := make(chan int, numJobs)

    var wg sync.WaitGroup
    for w := 1; w <= numWorkers; w++ {
        wg.Add(1)
        go worker(w, jobs, results, &wg)
    }

    for j := 1; j <= numJobs; j++ {
        jobs <- j
    }
    close(jobs)

    wg.Wait()
    close(results)

    var squares []int
    for r := range results {
        squares = append(squares, r)
    }
    sort.Ints(squares)
    fmt.Println("squares:", squares)

    total := 0
    for _, s := range squares {
        total += s
    }
    fmt.Println("total:", total)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
squares: [1 4 9 16 25 36 49 64 81]
total: 285
```

**Note on non-determinism:** which of the 3 workers processes which job, and the order they finish in, is genuinely random - they're racing to read from the same `jobs` channel. That's fine, because every job is still processed exactly once. Sorting the collected results before printing (and summing them) turns that race into an exactly checkable, deterministic final output - you'll get `squares: [1 4 9 16 25 36 49 64 81]` and `total: 285` on every run, regardless of which worker did what.

**Learning Objectives:**
- ✅ Build a fixed-size worker pool with a shared `jobs` channel and `results` channel
- ✅ Use `sync.WaitGroup` to know when every worker has finished before closing `results`
- ✅ Recognize which parts of a concurrent program are non-deterministic vs. which can still be made checkable

---

## Exercise 9: A Two-Stage Pipeline

**Objective:** Chain two channel-producing functions together, one stage's output feeding the next stage's input

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level21-exercise9
cd ~/projects/level21-exercise9
go mod init level21.example/exercise9
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func generator(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            out <- n
        }
    }()
    return out
}

func double(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            out <- n * 2
        }
    }()
    return out
}

func main() {
    gen := generator(1, 2, 3, 4, 5)
    doubled := double(gen)

    for v := range doubled {
        fmt.Println(v)
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
2
4
6
8
10
```

Each stage is a single sender feeding a single receiver, so order is preserved end to end - deterministic every run.

**Learning Objectives:**
- ✅ Write a function that returns a receive-only channel (`<-chan int`)
- ✅ Chain pipeline stages by feeding one stage's output channel into the next stage's input
- ✅ See each stage close its own output channel when its input is exhausted

---

## Exercise 10: Comprehensive Practice - Three-Stage Pipeline

**Objective:** Combine a generator, a squaring stage, and a summing stage into a full pipeline with one deterministic final answer

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level21-exercise10
cd ~/projects/level21-exercise10
go mod init level21.example/exercise10
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func generate(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            out <- n
        }
    }()
    return out
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            out <- n * n
        }
    }()
    return out
}

func sum(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        total := 0
        for n := range in {
            total += n
        }
        out <- total
    }()
    return out
}

func main() {
    nums := generate(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
    squares := square(nums)
    total := <-sum(squares)

    fmt.Println("sum of squares 1..10:", total)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
sum of squares 1..10: 385
```

Three stages, three goroutines, one channel between each pair - but because each stage has exactly one sender and one receiver, the whole pipeline is fully deterministic: `1²+2²+...+10² = 385`, every single run.

**Learning Objectives:**
- ✅ Compose three pipeline stages (generate → square → sum) using channels alone
- ✅ Understand each stage runs concurrently as its own goroutine, connected only by channels
- ✅ See that a well-designed channel pipeline yields exactly one deterministic, checkable answer

---

## Bonus Challenges

### Challenge 1: Fan-In

Write a `merge(a, b <-chan int) <-chan int` function that reads from two input channels concurrently and forwards every value onto a single output channel, closing the output only after both inputs are exhausted.

```bash
mkdir -p ~/projects/level21-bonus1
cd ~/projects/level21-bonus1
go mod init level21.example/bonus1
```

**Hints:**
- Launch one goroutine per input channel, each doing `for v := range in { out <- v }`
- Use a `sync.WaitGroup` with `Add(2)` so a third goroutine can `wg.Wait()` and then `close(out)` once both forwarders are done
- The *interleaving* of values from `a` and `b` on the merged channel is genuinely non-deterministic - only the final *set* of values received is guaranteed

### Challenge 2: Rate Limiter With a Token Bucket

Use a buffered channel as a "token bucket" to limit how many operations can run without waiting: pre-fill a buffered channel with a fixed number of tokens, have each operation receive a token before proceeding and (optionally) return it later.

```bash
mkdir -p ~/projects/level21-bonus2
cd ~/projects/level21-bonus2
go mod init level21.example/bonus2
```

**Hints:**
- `tokens := make(chan struct{}, 3)` gives you 3 "slots"; fill it with 3 sends before starting
- Each unit of work does `<-tokens` before running and `tokens <- struct{}{}` when finished, to return the slot
- `struct{}` is a zero-size type - perfect for a channel used purely as a signal, since the value itself carries no information

### Challenge 3: Simple Pub/Sub

Build a tiny publish/subscribe system: one publisher, multiple subscriber channels. Each message the publisher sends should reach every current subscriber.

```bash
mkdir -p ~/projects/level21-bonus3
cd ~/projects/level21-bonus3
go mod init level21.example/bonus3
```

**Hints:**
- Keep a slice of subscriber channels (`[]chan string`) inside a small `Broker` struct
- A `Subscribe()` method creates a new channel, appends it to the slice, and returns it to the caller
- A `Publish(msg string)` method loops over every subscriber channel and sends the message to each one
- Protect the subscriber slice with a `sync.Mutex` if `Subscribe` and `Publish` can be called concurrently - this is a case where a mutex genuinely earns its place alongside channels (see Best Practices in the README)

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Send and receive over unbuffered channels and observe genuine blocking
✅ Contrast buffered and unbuffered blocking behavior with real, verified evidence
✅ Close channels and distinguish a real zero value from a closed channel with comma-ok
✅ Range over a channel until it closes, and know why the producer must call `close`
✅ Use channel direction types and read a real compiler error for misusing one
✅ Recognize (and safely reproduce) a real deadlock and two real "closed channel" panics
✅ Build a worker pool with a bounded number of workers
✅ Compose multi-stage pipelines out of nothing but channels and goroutines

---

## Next Level

Level 22: Select
- Waiting on multiple channels at once with `select`
- `default` cases for non-blocking channel operations
- Timeouts with `time.After`
- Using `nil` channels to disable a `select` case at runtime

Great work! You're mastering concurrent Go! 🚀
