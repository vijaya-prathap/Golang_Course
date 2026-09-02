# Level 22: Select - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

**A note on timing and randomness:** several exercises below measure real elapsed time or real random tie-breaking. Every one of them is designed to terminate on its own - via a fixed iteration count, a bounded sleep, or a real timeout - so none of them can hang. Where output depends on real time or randomness, the exact numbers will vary slightly between runs and machines; what's called out as "expected" is the shape of the result (which case fired, roughly what percentage, roughly how much time elapsed), not a byte-for-byte guarantee.

---

## Exercise 1: Basic select Between Two Channels

**Objective:** Wait on two channels at once and see select proceed with whichever is ready first

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level22-exercise1
cd ~/projects/level22-exercise1
go mod init level22.example/exercise1
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
    ch1 := make(chan string)
    ch2 := make(chan string)

    go func() {
        time.Sleep(50 * time.Millisecond)
        ch1 <- "from ch1"
    }()
    go func() {
        time.Sleep(150 * time.Millisecond)
        ch2 <- "from ch2"
    }()

    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-ch1:
            fmt.Println("received:", msg1)
        case msg2 := <-ch2:
            fmt.Println("received:", msg2)
        }
    }
    fmt.Println("done")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
received: from ch1
received: from ch2
done
```

`ch1`'s goroutine sleeps for less time, so its send is ready first and `select` picks it up on the first loop pass; `ch2`'s value is picked up on the second pass. Looping around a `select` twice is what lets you receive from both channels.

**Learning Objectives:**
- ✅ Write a `select` with two receive cases
- ✅ Understand that `select` proceeds with whichever case is ready first
- ✅ Loop around a `select` to consume more than one ready value

---

## Exercise 2: Empirically Measuring Random Tie-Breaking

**Objective:** Prove, with real tallied counts, that select picks randomly among simultaneously-ready cases

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level22-exercise2
cd ~/projects/level22-exercise2
go mod init level22.example/exercise2
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    const iterations = 100000
    countA := 0
    countB := 0

    for i := 0; i < iterations; i++ {
        chA := make(chan int, 1)
        chB := make(chan int, 1)
        chA <- 1
        chB <- 1

        select {
        case <-chA:
            countA++
        case <-chB:
            countB++
        }
    }

    fmt.Printf("chA fired: %d times (%.1f%%)\n", countA, float64(countA)/iterations*100)
    fmt.Printf("chB fired: %d times (%.1f%%)\n", countB, float64(countB)/iterations*100)
    fmt.Println("(exact numbers vary between runs, but both should land close to 50%)")
}
EOF
```

3. Run the program (try running it 2-3 times to see the numbers shift slightly):

```bash
go run main.go
```

**Expected Output (one real captured run - your exact numbers will differ):**

```
chA fired: 49986 times (50.0%)
chB fired: 50014 times (50.0%)
(exact numbers vary between runs, but both should land close to 50%)
```

Both channels are pre-filled before the `select` ever runs, so both cases are ready on **every single iteration** - and yet the tally still lands close to a 50/50 split across 100,000 trials. If `select` favored whichever case is written first, you'd see something close to 100%/0% instead of a roughly even split.

**Learning Objectives:**
- ✅ Build two simultaneously-ready channel cases on purpose
- ✅ Tally real results across many iterations instead of assuming behavior
- ✅ Confirm empirically that select's tie-breaking is random, not ordered

---

## Exercise 3: Non-Blocking Operations With default

**Objective:** Use default to make both a receive and a send non-blocking

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level22-exercise3
cd ~/projects/level22-exercise3
go mod init level22.example/exercise3
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== Non-blocking receive from an empty channel ===")
    ch := make(chan int)
    select {
    case v := <-ch:
        fmt.Println("received:", v)
    default:
        fmt.Println("no value ready, default case ran")
    }

    fmt.Println("\n=== Non-blocking send into a full buffered channel ===")
    full := make(chan int, 1)
    full <- 100 // fill the only slot
    select {
    case full <- 200:
        fmt.Println("sent 200")
    default:
        fmt.Println("channel full, default case ran")
    }

    fmt.Println("\n=== Non-blocking send into a channel with room ===")
    room := make(chan int, 1)
    select {
    case room <- 42:
        fmt.Println("sent 42, no default needed")
    default:
        fmt.Println("this should not run")
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
=== Non-blocking receive from an empty channel ===
no value ready, default case ran

=== Non-blocking send into a full buffered channel ===
channel full, default case ran

=== Non-blocking send into a channel with room ===
sent 42, no default needed
```

**Learning Objectives:**
- ✅ Use `default` to make a receive non-blocking
- ✅ Use `default` to make a send non-blocking
- ✅ See that `default` only runs when every other case is genuinely not ready

---

## Exercise 4: Timeouts With time.After

**Objective:** Bound how long you wait for a channel using time.After, and observe both outcomes

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level22-exercise4
cd ~/projects/level22-exercise4
go mod init level22.example/exercise4
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
    fmt.Println("=== Case 1: the operation finishes before the timeout ===")
    ch := make(chan string)
    go func() {
        time.Sleep(50 * time.Millisecond)
        ch <- "work done"
    }()

    select {
    case msg := <-ch:
        fmt.Printf("received %q before the timeout\n", msg)
    case <-time.After(200 * time.Millisecond):
        fmt.Println("timed out")
    }

    fmt.Println("\n=== Case 2: the timeout wins ===")
    ch2 := make(chan string)
    go func() {
        time.Sleep(300 * time.Millisecond)
        ch2 <- "work done (too late)"
    }()

    start := time.Now()
    select {
    case msg := <-ch2:
        fmt.Println("received:", msg)
    case <-time.After(100 * time.Millisecond):
        elapsed := time.Since(start)
        fmt.Printf("timed out (elapsed roughly %v, bounded by the 100ms time.After)\n", elapsed.Round(10*time.Millisecond))
    }
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output (elapsed time will vary slightly by machine; this is a real captured run):**

```
=== Case 1: the operation finishes before the timeout ===
received "work done" before the timeout

=== Case 2: the timeout wins ===
timed out (elapsed roughly 100ms, bounded by the 100ms time.After)
```

The worker in Case 1 finishes in 50ms, well inside its 200ms budget, so the real result wins. The worker in Case 2 takes 300ms but only gets a 100ms budget, so the timeout fires first - the elapsed time is bounded by (approximately equal to) the `time.After` duration, not the worker's sleep.

**Learning Objectives:**
- ✅ Use `case <-time.After(d):` as a timeout arm of a select
- ✅ See a real operation win a race against its timeout
- ✅ See a real timeout win against a slower operation, with a measured elapsed time

---

## Exercise 5: The done-Channel Cancellation Pattern

**Objective:** Stop a running goroutine cleanly from the outside using a done channel

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level22-exercise5
cd ~/projects/level22-exercise5
go mod init level22.example/exercise5
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "time"
)

func worker(done <-chan struct{}, work <-chan int) {
    for {
        select {
        case <-done:
            fmt.Println("worker: received done signal, stopping")
            return
        case w, ok := <-work:
            if !ok {
                fmt.Println("worker: work channel closed, stopping")
                return
            }
            fmt.Println("worker: processing", w)
        }
    }
}

func main() {
    done := make(chan struct{})
    work := make(chan int)

    finished := make(chan struct{})
    go func() {
        worker(done, work)
        close(finished)
    }()

    for i := 1; i <= 3; i++ {
        work <- i
        time.Sleep(20 * time.Millisecond)
    }

    close(done)
    <-finished
    fmt.Println("main: worker confirmed stopped")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
worker: processing 1
worker: processing 2
worker: processing 3
worker: received done signal, stopping
main: worker confirmed stopped
```

**Learning Objectives:**
- ✅ Build a goroutine that selects between real work and a cancellation signal
- ✅ Use `close(done)` to signal cancellation to a running goroutine
- ✅ Use a second channel (`finished`) to confirm the goroutine actually stopped before moving on

---

## Exercise 6: select With a Send Case

**Objective:** Use a channel send as a select case, and see it lose to another case when it isn't ready

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level22-exercise6
cd ~/projects/level22-exercise6
go mod init level22.example/exercise6
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== Send case is NOT ready: channel full, done wins ===")
    ch := make(chan int, 1)
    ch <- 100 // fill it so a further send would block

    done := make(chan struct{})
    close(done) // already closed, so this case is immediately ready

    select {
    case ch <- 200:
        fmt.Println("sent 200 (unexpected - channel should be full)")
    case <-done:
        fmt.Println("done fired instead, since the send case was not ready")
    }
    fmt.Println("channel still holds:", <-ch)

    fmt.Println("\n=== Send case IS ready: channel has room, send wins ===")
    ch2 := make(chan int, 1)
    done2 := make(chan struct{}) // never closed - will never be ready

    select {
    case ch2 <- 42:
        fmt.Println("sent 42 - the channel had room, so the send case fired")
    case <-done2:
        fmt.Println("this should never happen")
    }
    fmt.Println("received back:", <-ch2)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Send case is NOT ready: channel full, done wins ===
done fired instead, since the send case was not ready
channel still holds: 100

=== Send case IS ready: channel has room, send wins ===
sent 42 - the channel had room, so the send case fired
received back: 42
```

**Learning Objectives:**
- ✅ Write a select case that sends (`case ch <- v:`) instead of receives
- ✅ Understand a send case is only "ready" if that send could complete immediately
- ✅ See a send case both win and lose against a competing receive case

---

## Exercise 7: select on a Closed Channel

**Objective:** Verify that a closed channel's case is always immediately ready

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level22-exercise7
cd ~/projects/level22-exercise7
go mod init level22.example/exercise7
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    ch := make(chan int, 3)
    ch <- 1
    ch <- 2
    close(ch)

    fmt.Println("=== A closed channel case is ready immediately, even with buffered values remaining ===")
    for i := 0; i < 2; i++ {
        select {
        case v, ok := <-ch:
            fmt.Println("received:", v, "ok:", ok)
        }
    }

    fmt.Println("\n=== Once drained, a closed channel keeps returning the zero value + false, forever ===")
    for i := 0; i < 3; i++ {
        select {
        case v, ok := <-ch:
            fmt.Println("received:", v, "ok:", ok)
        }
    }

    fmt.Println("\nNote: a select case reading a closed, drained channel is ALWAYS immediately ready.")
    fmt.Println("If that were the only case in an unbounded 'for { select {...} }' loop, the loop")
    fmt.Println("would busy-loop, spinning as fast as the CPU allows. Always pair a case like this")
    fmt.Println("with a way to stop (a counter, an ok-check that returns/breaks, or another exit case).")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== A closed channel case is ready immediately, even with buffered values remaining ===
received: 1 ok: true
received: 2 ok: true

=== Once drained, a closed channel keeps returning the zero value + false, forever ===
received: 0 ok: false
received: 0 ok: false
received: 0 ok: false

Note: a select case reading a closed, drained channel is ALWAYS immediately ready.
If that were the only case in an unbounded 'for { select {...} }' loop, the loop
would busy-loop, spinning as fast as the CPU allows. Always pair a case like this
with a way to stop (a counter, an ok-check that returns/breaks, or another exit case).
```

This program deliberately caps each demonstration loop at a fixed count (2, then 3) instead of looping forever - reading a closed channel in a real unbounded loop without checking `ok` is exactly the busy-loop risk being demonstrated, not something to actually run unbounded.

**Learning Objectives:**
- ✅ See a closed channel's case fire immediately, draining any buffered values first
- ✅ See a closed, drained channel keep returning the zero value and `ok == false` forever
- ✅ Understand why an unguarded closed-channel case risks a CPU busy-loop

---

## Exercise 8: Replacing a Naive time.After Loop With time.NewTimer

**Objective:** Compare the leaky time.After-in-a-loop pattern with the correct time.NewTimer pattern

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level22-exercise8
cd ~/projects/level22-exercise8
go mod init level22.example/exercise8
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "time"
)

func naiveLoop() {
    fmt.Println("=== ANTI-PATTERN: time.After inside a loop (works, but leaks a timer each iteration) ===")
    ch := make(chan int)
    go func() {
        for i := 1; i <= 3; i++ {
            time.Sleep(30 * time.Millisecond)
            ch <- i
        }
        close(ch)
    }()

    for {
        select {
        case v, ok := <-ch:
            if !ok {
                fmt.Println("channel closed, anti-pattern loop done")
                return
            }
            fmt.Println("naive loop received:", v)
        case <-time.After(200 * time.Millisecond):
            fmt.Println("naive loop: unexpected timeout")
        }
    }
}

func betterLoop() {
    fmt.Println("\n=== BETTER: reuse a single time.NewTimer with Stop()/Reset() ===")
    ch := make(chan int)
    go func() {
        for i := 1; i <= 3; i++ {
            time.Sleep(30 * time.Millisecond)
            ch <- i
        }
        close(ch)
    }()

    timer := time.NewTimer(200 * time.Millisecond)
    defer timer.Stop()

    for {
        if !timer.Stop() {
            select {
            case <-timer.C:
            default:
            }
        }
        timer.Reset(200 * time.Millisecond)

        select {
        case v, ok := <-ch:
            if !ok {
                fmt.Println("channel closed, better loop done")
                return
            }
            fmt.Println("better loop received:", v)
        case <-timer.C:
            fmt.Println("better loop: unexpected timeout")
        }
    }
}

func main() {
    naiveLoop()
    betterLoop()
    fmt.Println("\nBoth loops print the same visible output. The difference is under the hood:")
    fmt.Println("the naive version creates and abandons a brand-new timer on every single")
    fmt.Println("iteration (3 of them here - far worse in a hot loop with thousands of")
    fmt.Println("iterations), while the better version reuses exactly one timer for the whole loop.")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== ANTI-PATTERN: time.After inside a loop (works, but leaks a timer each iteration) ===
naive loop received: 1
naive loop received: 2
naive loop received: 3
channel closed, anti-pattern loop done

=== BETTER: reuse a single time.NewTimer with Stop()/Reset() ===
better loop received: 1
better loop received: 2
better loop received: 3
channel closed, better loop done

Both loops print the same visible output. The difference is under the hood:
the naive version creates and abandons a brand-new timer on every single
iteration (3 of them here - far worse in a hot loop with thousands of
iterations), while the better version reuses exactly one timer for the whole loop.
```

**Learning Objectives:**
- ✅ See the naive `time.After`-in-a-loop pattern produce correct but wasteful behavior
- ✅ Replace it with a single `time.NewTimer`, using `Stop()`/`Reset()` correctly each iteration
- ✅ Understand that both patterns are functionally identical from the outside - the cost is purely in abandoned timers

---

## Exercise 9: A Multiplexer That Disables Sources With nil Channels

**Objective:** Multiplex two channels with select+for, disabling each source with a nil channel once it's exhausted

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level22-exercise9
cd ~/projects/level22-exercise9
go mod init level22.example/exercise9
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
    chA := make(chan int)
    chB := make(chan int)

    go func() {
        chA <- 1
        time.Sleep(30 * time.Millisecond)
        chA <- 2
        close(chA)
    }()
    go func() {
        time.Sleep(15 * time.Millisecond)
        chB <- 100
        close(chB)
    }()

    a, b := chA, chB
    received := 0

    for a != nil || b != nil {
        select {
        case v, ok := <-a:
            if !ok {
                fmt.Println("channel A exhausted, disabling its case (setting to nil)")
                a = nil
                continue
            }
            fmt.Println("from A:", v)
            received++
        case v, ok := <-b:
            if !ok {
                fmt.Println("channel B exhausted, disabling its case (setting to nil)")
                b = nil
                continue
            }
            fmt.Println("from B:", v)
            received++
        }
    }
    fmt.Println("both channels exhausted and disabled; total values received:", received)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
from A: 1
from B: 100
channel B exhausted, disabling its case (setting to nil)
from A: 2
channel A exhausted, disabling its case (setting to nil)
both channels exhausted and disabled; total values received: 3
```

The sleeps are spaced far enough apart (15ms and 30ms) that the arrival order is reliably deterministic on a normal machine: `A`'s first value, then `B`'s only value, then `B` closing, then `A`'s second value, then `A` closing. Once a source's case reports `ok == false`, setting that local variable to `nil` removes it from the `select` for good - the loop's `a != nil || b != nil` condition is what lets it stop once both are gone.

**Learning Objectives:**
- ✅ Multiplex two channels inside a single `select` + `for` loop
- ✅ Use a `nil` channel to permanently disable a case once its source is exhausted
- ✅ Avoid the closed-channel busy-loop risk from Exercise 7 by disabling, not just detecting, an exhausted source

---

## Exercise 10: Comprehensive Practice - An Event Dispatcher

**Objective:** Build a dispatcher that multiplexes an events channel and a done channel, stopping after N events or a timeout, whichever comes first

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level22-exercise10
cd ~/projects/level22-exercise10
go mod init level22.example/exercise10
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "time"
)

func dispatcher(events <-chan string, done <-chan struct{}, maxEvents int, timeout time.Duration) {
    count := 0
    for {
        select {
        case ev, ok := <-events:
            if !ok {
                fmt.Println("dispatcher: events channel closed, stopping")
                return
            }
            count++
            fmt.Printf("dispatcher: handled event %d: %s\n", count, ev)
            if count >= maxEvents {
                fmt.Println("dispatcher: reached max events, stopping")
                return
            }
        case <-done:
            fmt.Println("dispatcher: done signal received, stopping early")
            return
        case <-time.After(timeout):
            fmt.Println("dispatcher: timed out waiting for an event, stopping")
            return
        }
    }
}

func main() {
    fmt.Println("=== Scenario 1: dispatcher stops after reaching max events ===")
    events := make(chan string)
    done := make(chan struct{})
    go func() {
        for _, e := range []string{"login", "click", "purchase"} {
            events <- e
        }
    }()
    dispatcher(events, done, 3, 300*time.Millisecond)

    fmt.Println("\n=== Scenario 2: dispatcher stops early via a done signal ===")
    events2 := make(chan string)
    done2 := make(chan struct{})
    go func() {
        events2 <- "start"
        close(done2) // tell the dispatcher to stop before any more events arrive
    }()
    dispatcher(events2, done2, 10, 300*time.Millisecond)

    fmt.Println("\n=== Scenario 3: dispatcher stops via timeout (no events ever arrive) ===")
    events3 := make(chan string)
    done3 := make(chan struct{})
    dispatcher(events3, done3, 10, 150*time.Millisecond)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Scenario 1: dispatcher stops after reaching max events ===
dispatcher: handled event 1: login
dispatcher: handled event 2: click
dispatcher: handled event 3: purchase
dispatcher: reached max events, stopping

=== Scenario 2: dispatcher stops early via a done signal ===
dispatcher: handled event 1: start
dispatcher: done signal received, stopping early

=== Scenario 3: dispatcher stops via timeout (no events ever arrive) ===
dispatcher: timed out waiting for an event, stopping
```

Scenario 3 genuinely waits out its full 150ms timeout (nothing is ever sent on `events3` or `done3`) before printing - a real, bounded wait, not an instant return. All three scenarios are guaranteed to terminate: by count, by signal, or by timeout.

**Learning Objectives:**
- ✅ Combine all three `select` exit conditions (data received, cancellation, timeout) in one dispatcher
- ✅ Design a function that reliably terminates via whichever condition happens first
- ✅ See the same dispatcher function handle three completely different real-world shutdown reasons correctly

---

## Bonus Challenges

### Challenge 1: Two-Way Ping-Pong

Write a program where two goroutines pass a "ball" back and forth over two channels a fixed number of times using `select`, then stop cleanly.

```bash
mkdir -p ~/projects/level22-bonus1
cd ~/projects/level22-bonus1
go mod init level22.example/bonus1
```

**Hints:**
- Use two channels, `ping` and `pong`; goroutine A sends on `ping` and receives on `pong`, goroutine B does the reverse
- Give each goroutine a `select` with a `done` case (or just a bounded loop counter) so the exchange stops after, say, 5 rounds instead of continuing forever
- Have `main` kick off the very first send so both goroutines have something to react to

### Challenge 2: Fan-In With select

Write a `merge(a, b <-chan int) <-chan int` that uses a single `select` in a loop (instead of two separate forwarding goroutines) to pull from whichever of `a` or `b` is ready and forward it to one output channel.

```bash
mkdir -p ~/projects/level22-bonus2
cd ~/projects/level22-bonus2
go mod init level22.example/bonus2
```

**Hints:**
- Use the `nil`-channel trick from Exercise 9: when one input closes, set its local variable to `nil` so the `select` stops considering it
- Close the output channel only once **both** input variables are `nil`
- The interleaving of values from `a` and `b` is genuinely non-deterministic - only the final set received is guaranteed

### Challenge 3: A Rate-Limited Worker With a Ticker

Use `time.NewTicker` (or `time.Tick`) inside a `select` to process at most one unit of work per tick, alongside a `done` channel that stops the worker early.

```bash
mkdir -p ~/projects/level22-bonus3
cd ~/projects/level22-bonus3
go mod init level22.example/bonus3
```

**Hints:**
- `ticker := time.NewTicker(100 * time.Millisecond)` produces a value on `ticker.C` every interval - `select` on `ticker.C` alongside `work` and `done`
- Always `defer ticker.Stop()` - unlike a one-shot `time.After`, a ticker keeps firing forever until stopped, so an un-stopped ticker is a real, ongoing leak, not just a one-time allocation
- Bound the exercise with a fixed number of ticks or a fixed number of processed items so the program still terminates on its own

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Write a select that waits on multiple channels and proceeds with whichever is ready first
✅ Prove select's random tie-breaking empirically instead of assuming it
✅ Use default to make both sends and receives non-blocking
✅ Bound a wait with time.After and measure the real elapsed time on both sides of a race
✅ Build a done-channel cancellation pattern that a select loop checks alongside real work
✅ Write and reason about select cases that send rather than receive
✅ Recognize a closed channel's case as permanently ready, and the busy-loop risk that implies
✅ Replace a leaky time.After-in-a-loop with a reused time.NewTimer
✅ Use nil channels to cleanly disable a select case at runtime
✅ Combine select with a for loop to build a multiplexer or dispatcher with a guaranteed termination path

---

## Next Level

Level 23: Mutex, WaitGroup & Atomic
- sync.Mutex and sync.RWMutex for protecting shared state directly
- sync.WaitGroup in depth for waiting on many goroutines
- The sync/atomic package for lock-free counters and flags
- Choosing between channels, mutexes, and atomics for a given problem

Great work! You're mastering concurrent Go! 🚀
