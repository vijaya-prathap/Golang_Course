# Level 24: Context - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

**A note on timing and determinism:** every exercise below uses short, real `time.Sleep`/`time.After` durations so each program terminates in well under a second. Exact elapsed times will vary slightly between runs and machines (scheduling always adds a little overhead) - where that matters, it's called out explicitly, and every expected output below is a real, captured run, not a hand-computed guess.

---

## Exercise 1: context.Background() and context.TODO()

**Objective:** Create both root contexts and confirm they behave identically (no error, no deadline, no value)

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level24-exercise1
cd ~/projects/level24-exercise1
go mod init level24.example/exercise1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "context"
    "fmt"
)

func main() {
    bg := context.Background()
    todo := context.TODO()

    fmt.Println("=== context.Background() ===")
    fmt.Println("Err():", bg.Err())
    fmt.Println("Done():", bg.Done())
    fmt.Println("Value(\"x\"):", bg.Value("x"))
    fmt.Printf("Type: %T\n", bg)

    fmt.Println("\n=== context.TODO() ===")
    fmt.Println("Err():", todo.Err())
    fmt.Println("Done():", todo.Done())
    fmt.Println("Value(\"x\"):", todo.Value("x"))
    fmt.Printf("Type: %T\n", todo)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== context.Background() ===
Err(): <nil>
Done(): <nil>
Value("x"): <nil>
Type: context.backgroundCtx

=== context.TODO() ===
Err(): <nil>
Done(): <nil>
Value("x"): <nil>
Type: context.todoCtx
```

**Learning Objectives:**
- ✅ Create both root contexts with no parent
- ✅ Confirm neither is ever canceled, deadlined, or carries a value
- ✅ Understand Background() marks a true root; TODO() marks unfinished context plumbing

---

## Exercise 2: context.WithCancel Stopping a Real Goroutine

**Objective:** Cancel a running worker goroutine manually and verify it actually stops

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level24-exercise2
cd ~/projects/level24-exercise2
go mod init level24.example/exercise2
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "context"
    "fmt"
    "time"
)

func worker(ctx context.Context, stopped chan<- struct{}) {
    ticks := 0
    for {
        select {
        case <-ctx.Done():
            fmt.Println("worker: received cancellation, stopping. ticks completed:", ticks)
            stopped <- struct{}{}
            return
        case <-time.After(20 * time.Millisecond):
            ticks++
            fmt.Println("worker: tick", ticks)
        }
    }
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    stopped := make(chan struct{})

    go worker(ctx, stopped)

    time.Sleep(100 * time.Millisecond) // let the worker tick a few times
    fmt.Println("main: calling cancel()")
    cancel()

    <-stopped
    fmt.Println("main: confirmed worker actually stopped")
    fmt.Println("ctx.Err():", ctx.Err())
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
worker: tick 1
worker: tick 2
worker: tick 3
worker: tick 4
main: calling cancel()
worker: received cancellation, stopping. ticks completed: 4
main: confirmed worker actually stopped
ctx.Err(): context canceled
```

The exact tick count may vary by a tick or two between runs/machines, but `main` genuinely blocks on `<-stopped` until the worker's own `select` observes `ctx.Done()` - this is a real, verified stop, not a hope.

**Learning Objectives:**
- ✅ Create a cancelable context with context.WithCancel
- ✅ Stop a goroutine on demand by calling the returned cancel()
- ✅ Verify the stop for real via a channel, not just by calling cancel and moving on

---

## Exercise 3: context.WithTimeout - Real Elapsed-Time Verification

**Objective:** Let a context cancel itself automatically after a duration, and measure the real elapsed time

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level24-exercise3
cd ~/projects/level24-exercise3
go mod init level24.example/exercise3
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
    defer cancel()

    start := time.Now()
    <-ctx.Done()
    elapsed := time.Since(start)

    fmt.Println("ctx.Done() fired")
    fmt.Println("ctx.Err():", ctx.Err())
    fmt.Printf("elapsed: %v (requested timeout: 200ms)\n", elapsed.Round(time.Millisecond))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output (one real captured run - your exact elapsed time will vary slightly by machine):**

```
ctx.Done() fired
ctx.Err(): context deadline exceeded
elapsed: 201ms (requested timeout: 200ms)
```

**Learning Objectives:**
- ✅ Create a self-canceling context with context.WithTimeout
- ✅ Confirm ctx.Done() fires automatically with no manual cancel() call
- ✅ Measure and understand real elapsed time vs. requested duration

---

## Exercise 4: context.WithDeadline - An Absolute Cutoff Time

**Objective:** Cancel a context at an absolute point in time rather than a duration

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level24-exercise4
cd ~/projects/level24-exercise4
go mod init level24.example/exercise4
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    deadline := time.Now().Add(150 * time.Millisecond)
    ctx, cancel := context.WithDeadline(context.Background(), deadline)
    defer cancel()

    start := time.Now()
    <-ctx.Done()
    elapsed := time.Since(start)

    fmt.Println("ctx.Done() fired")
    fmt.Println("ctx.Err():", ctx.Err())
    fmt.Printf("elapsed: %v (deadline was 150ms from start)\n", elapsed.Round(time.Millisecond))

    d, ok := ctx.Deadline()
    fmt.Println("Deadline() ok:", ok)
    fmt.Println("Deadline() time equals requested:", d.Equal(deadline))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output (one real captured run - exact elapsed time varies by machine):**

```
ctx.Done() fired
ctx.Err(): context deadline exceeded
elapsed: 151ms (deadline was 150ms from start)
Deadline() ok: true
Deadline() time equals requested: true
```

**Learning Objectives:**
- ✅ Create a context tied to an absolute time.Time with context.WithDeadline
- ✅ Retrieve the deadline back out with ctx.Deadline()
- ✅ Understand WithTimeout is WithDeadline expressed as a relative duration

---

## Exercise 5: Distinguishing Canceled From DeadlineExceeded

**Objective:** Produce both context.Canceled and context.DeadlineExceeded for real, in the same run

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level24-exercise5
cd ~/projects/level24-exercise5
go mod init level24.example/exercise5
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "context"
    "errors"
    "fmt"
    "time"
)

func main() {
    fmt.Println("=== Manual Cancellation ===")
    ctx1, cancel1 := context.WithCancel(context.Background())
    cancel1()
    <-ctx1.Done()
    fmt.Println("ctx1.Err():", ctx1.Err())
    fmt.Println("errors.Is(ctx1.Err(), context.Canceled):", errors.Is(ctx1.Err(), context.Canceled))

    fmt.Println("\n=== Timeout Expiry ===")
    ctx2, cancel2 := context.WithTimeout(context.Background(), 50*time.Millisecond)
    defer cancel2()
    <-ctx2.Done()
    fmt.Println("ctx2.Err():", ctx2.Err())
    fmt.Println("errors.Is(ctx2.Err(), context.DeadlineExceeded):", errors.Is(ctx2.Err(), context.DeadlineExceeded))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Manual Cancellation ===
ctx1.Err(): context canceled
errors.Is(ctx1.Err(), context.Canceled): true

=== Timeout Expiry ===
ctx2.Err(): context deadline exceeded
errors.Is(ctx2.Err(), context.DeadlineExceeded): true
```

**Learning Objectives:**
- ✅ Produce a real context.Canceled error via manual cancel()
- ✅ Produce a real context.DeadlineExceeded error via an expired timeout
- ✅ Use errors.Is to check against the sentinel errors idiomatically

---

## Exercise 6: The ctx.Done() select Pattern - Work vs. Cancellation

**Objective:** Race real work against context cancellation using select, and observe both outcomes

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level24-exercise6
cd ~/projects/level24-exercise6
go mod init level24.example/exercise6
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "context"
    "fmt"
    "time"
)

func doWork(ctx context.Context, workDuration time.Duration, workCh chan<- string) {
    time.Sleep(workDuration)
    select {
    case workCh <- "work result":
    case <-ctx.Done():
    }
}

func runWithSelect(label string, ctx context.Context, workCh <-chan string) {
    select {
    case <-ctx.Done():
        fmt.Printf("[%s] canceled before work finished: %v\n", label, ctx.Err())
    case result := <-workCh:
        fmt.Printf("[%s] work finished first: %s\n", label, result)
    }
}

func main() {
    fmt.Println("=== Case 1: Work finishes before timeout ===")
    ctx1, cancel1 := context.WithTimeout(context.Background(), 200*time.Millisecond)
    defer cancel1()
    workCh1 := make(chan string)
    go doWork(ctx1, 20*time.Millisecond, workCh1)
    runWithSelect("fast work", ctx1, workCh1)

    fmt.Println("\n=== Case 2: Timeout fires before work finishes ===")
    ctx2, cancel2 := context.WithTimeout(context.Background(), 50*time.Millisecond)
    defer cancel2()
    workCh2 := make(chan string)
    go doWork(ctx2, 300*time.Millisecond, workCh2)
    runWithSelect("slow work", ctx2, workCh2)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Case 1: Work finishes before timeout ===
[fast work] work finished first: work result

=== Case 2: Timeout fires before work finishes ===
[slow work] canceled before work finished: context deadline exceeded
```

Both outcomes are exactly deterministic here because the work duration and the timeout are far enough apart (20ms vs. 200ms, and 300ms vs. 50ms) that scheduling jitter can't flip the winner.

**Learning Objectives:**
- ✅ Write the standard `select { case <-ctx.Done(): ...; case result := <-workCh: ... }` pattern
- ✅ See real work win the race when it's fast enough
- ✅ See cancellation win the race when work is too slow
- ✅ Understand why doWork itself must also select on ctx.Done() when sending, to avoid leaking

---

## Exercise 7: context.WithValue - Typed Keys and the Collision Risk

**Objective:** Use a correctly-typed unexported key, then prove why a plain string key is dangerous

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level24-exercise7
cd ~/projects/level24-exercise7
go mod init level24.example/exercise7
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "context"
    "fmt"
)

type requestIDKey struct{} // unexported, unique type - avoids collisions with any other package's keys

func withRequestID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, requestIDKey{}, id)
}

func requestIDFrom(ctx context.Context) (string, bool) {
    id, ok := ctx.Value(requestIDKey{}).(string)
    return id, ok
}

func main() {
    fmt.Println("=== Correct Pattern: Typed Unexported Key ===")
    ctx := withRequestID(context.Background(), "req-12345")
    id, ok := requestIDFrom(ctx)
    fmt.Println("requestID:", id, "found:", ok)

    fmt.Println("\n=== The Collision Risk of Plain String Keys ===")
    // Two unrelated pieces of code both decide to use the string "id" as a context key.
    ctxA := context.WithValue(context.Background(), "id", "package-A-value")
    ctxB := context.WithValue(ctxA, "id", "package-B-value") // silently overwrites package A's value!

    fmt.Println("Value for string key \"id\":", ctxB.Value("id"))
    fmt.Println("package A's original value is now unreachable - overwritten by package B's use of the same string key")

    fmt.Println("\n=== Distinct Typed Keys Never Collide ===")
    type keyA struct{}
    type keyB struct{}
    ctxTyped := context.WithValue(context.Background(), keyA{}, "package-A-value")
    ctxTyped = context.WithValue(ctxTyped, keyB{}, "package-B-value")
    fmt.Println("keyA value:", ctxTyped.Value(keyA{}))
    fmt.Println("keyB value:", ctxTyped.Value(keyB{}))
    fmt.Println("both values survive - distinct unexported types can never collide, even carrying identical string data")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Correct Pattern: Typed Unexported Key ===
requestID: req-12345 found: true

=== The Collision Risk of Plain String Keys ===
Value for string key "id": package-B-value
package A's original value is now unreachable - overwritten by package B's use of the same string key

=== Distinct Typed Keys Never Collide ===
keyA value: package-A-value
keyB value: package-B-value
both values survive - distinct unexported types can never collide, even carrying identical string data
```

**Learning Objectives:**
- ✅ Store and retrieve a value using a typed, unexported context key
- ✅ Watch a real, silent collision happen when two callers reuse the same string key
- ✅ Confirm distinct typed keys never collide, even with identical string contents

---

## Exercise 8: Always defer cancel() - What Leaks If You Skip It

**Objective:** Compare a discarded cancel() against a deferred one, and capture the real go vet warning for the mistake

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level24-exercise8
cd ~/projects/level24-exercise8
go mod init level24.example/exercise8
```

2. Create the "bad" version, `bad.go`, that discards the cancel function:

```bash
cat > bad.go << 'EOF'
package main

import (
    "context"
    "fmt"
    "time"
)

func doSomething(ctx context.Context) {
    select {
    case <-ctx.Done():
        fmt.Println("bad: ctx already canceled/expired:", ctx.Err())
    case <-time.After(10 * time.Millisecond):
        fmt.Println("bad: finished a quick unit of work")
    }
}

func main() {
    ctx, _ := context.WithTimeout(context.Background(), 5*time.Second) // cancel discarded - LEAK
    doSomething(ctx)
    fmt.Println("main: done (the 5-second internal timer is still alive until it fires)")
}
EOF
```

3. Run it, then vet it:

```bash
go run bad.go
go vet bad.go
```

**Expected Output (run):**

```
bad: finished a quick unit of work
main: done (the 5-second internal timer is still alive until it fires)
```

**Expected Output (go vet - a real, captured warning):**

```
bad.go:19:10: the cancel function returned by context.WithTimeout should be called, not discarded, to avoid a context leak
```

4. Now replace it with the fixed version, `bad.go` (same file, corrected):

```bash
cat > bad.go << 'EOF'
package main

import (
    "context"
    "fmt"
    "time"
)

func doSomething(ctx context.Context) {
    select {
    case <-ctx.Done():
        fmt.Println("good: ctx already canceled/expired:", ctx.Err())
    case <-time.After(10 * time.Millisecond):
        fmt.Println("good: finished a quick unit of work")
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel() // released the instant main() returns - the timer is stopped immediately, not in 5s

    doSomething(ctx)
    fmt.Println("main: done (deferred cancel will release the timer right now, not in 5 seconds)")
}
EOF
```

5. Run and vet the fixed version:

```bash
go run bad.go
go vet bad.go
```

**Expected Output (run):**

```
good: finished a quick unit of work
main: done (deferred cancel will release the timer right now, not in 5 seconds)
```

**Expected Output (go vet):** no output at all - the command exits cleanly, confirming the leak is fixed.

**What actually leaks if you skip `defer cancel()`:** the context object and its internal timer stay alive and registered for the *entire* requested timeout/deadline, even after your real work is long done. In a short CLI program the process exits and cleans everything up anyway - but in a long-running server handling many requests, every forgotten `cancel()` leaves one more context and timer alive until its own timeout eventually fires, and that adds up to real, accumulating memory and scheduler overhead across thousands of requests.

**Learning Objectives:**
- ✅ See the real go vet warning for a discarded WithTimeout cancel function
- ✅ Fix it with defer cancel() and confirm go vet goes silent
- ✅ Explain precisely what resource stays alive when cancel() is skipped, and why it matters more in long-running programs than short ones

---

## Exercise 9: Propagating One Context Through Fanned-Out Goroutines

**Objective:** Cancel a single parent context and confirm every one of several worker goroutines stops

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level24-exercise9
cd ~/projects/level24-exercise9
go mod init level24.example/exercise9
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "context"
    "fmt"
    "sync"
    "time"
)

func worker(ctx context.Context, id int, wg *sync.WaitGroup, stopped chan<- int) {
    defer wg.Done()
    for {
        select {
        case <-ctx.Done():
            stopped <- id
            return
        case <-time.After(15 * time.Millisecond):
            // simulate doing a small unit of work
        }
    }
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    const numWorkers = 5

    var wg sync.WaitGroup
    stopped := make(chan int, numWorkers)

    for i := 1; i <= numWorkers; i++ {
        wg.Add(1)
        go worker(ctx, i, &wg, stopped)
    }

    time.Sleep(80 * time.Millisecond)
    fmt.Println("main: canceling parent context - all workers should stop")
    cancel()

    wg.Wait()
    close(stopped)

    count := 0
    for range stopped {
        count++
    }
    fmt.Printf("all %d workers confirmed stopped after cancellation\n", count)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
main: canceling parent context - all workers should stop
all 5 workers confirmed stopped after cancellation
```

Which worker finishes its current 15ms tick first is non-deterministic, but the final count is always exactly 5 - `wg.Wait()` doesn't return until every single one has confirmed via `stopped`.

**Learning Objectives:**
- ✅ Share one context across multiple fanned-out goroutines
- ✅ Confirm a single cancel() stops every worker, verified with sync.WaitGroup
- ✅ Understand a context.Context is a cheap-to-copy interface value, safely shared across goroutines

---

## Exercise 10: Comprehensive Practice - A Cancelable Pipeline Job

**Objective:** Run a multi-step "pipeline" that can be canceled partway through, cleanly stopping and reporting how far it got

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level24-exercise10
cd ~/projects/level24-exercise10
go mod init level24.example/exercise10
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "context"
    "fmt"
    "time"
)

type step struct {
    name     string
    duration time.Duration
}

func runPipeline(ctx context.Context, steps []step) (completed []string, err error) {
    for _, s := range steps {
        select {
        case <-ctx.Done():
            return completed, ctx.Err()
        case <-time.After(s.duration):
            completed = append(completed, s.name)
        }
    }
    return completed, nil
}

func main() {
    steps := []step{
        {"validate", 30 * time.Millisecond},
        {"fetch", 30 * time.Millisecond},
        {"transform", 30 * time.Millisecond},
        {"persist", 30 * time.Millisecond},
        {"notify", 30 * time.Millisecond},
    }

    fmt.Println("=== Run 1: Enough time for the whole pipeline ===")
    ctx1, cancel1 := context.WithTimeout(context.Background(), 500*time.Millisecond)
    defer cancel1()
    completed1, err1 := runPipeline(ctx1, steps)
    fmt.Println("completed steps:", completed1)
    fmt.Println("error:", err1)

    fmt.Println("\n=== Run 2: Canceled partway through ===")
    ctx2, cancel2 := context.WithTimeout(context.Background(), 70*time.Millisecond)
    defer cancel2()
    completed2, err2 := runPipeline(ctx2, steps)
    fmt.Println("completed steps:", completed2)
    fmt.Println("error:", err2)
    fmt.Printf("pipeline stopped after %d of %d steps\n", len(completed2), len(steps))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Run 1: Enough time for the whole pipeline ===
completed steps: [validate fetch transform persist notify]
error: <nil>

=== Run 2: Canceled partway through ===
completed steps: [validate fetch]
error: context deadline exceeded
pipeline stopped after 2 of 5 steps
```

**Learning Objectives:**
- ✅ Combine multiple context concepts into one realistic, cancelable job
- ✅ Report partial progress cleanly instead of losing all information on cancellation
- ✅ See the exact same pipeline function behave correctly whether it completes or is cut off

---

## Bonus Challenges

### Challenge 1: Context-Aware Retry With Backoff

Write `retryWithBackoff(ctx context.Context, maxAttempts int, base time.Duration, op func(int) error) error` that retries `op`, waiting `base * attempt` between attempts, and gives up early - returning `ctx.Err()` - if the context is canceled while waiting.

```bash
mkdir -p ~/projects/level24-bonus1
cd ~/projects/level24-bonus1
go mod init level24.example/bonus1
```

**Hints:**
- Between attempts, `select` on `<-ctx.Done()` and `<-time.After(backoff)` - whichever fires first decides whether you retry or give up
- Verify it two ways: a context with plenty of time (the retry eventually succeeds) and a very short-timeout context (it gives up mid-backoff with `context.DeadlineExceeded`)
- A real captured run of this exact pattern: attempts 1-3 fail, backoff grows each time, and a short-timeout context stops it after attempt 2 with `retry: context canceled during backoff, giving up early` followed by `final error: context deadline exceeded`

### Challenge 2: A Cancellable Worker Pool

Recall the Level 21 worker pool pattern (fixed workers pulling from a `jobs` channel, pushing to a `results` channel) and make it cancellable: every worker, and the goroutine feeding `jobs`, should select on `ctx.Done()` so the whole pool stops cleanly the moment a `context.WithTimeout` expires - not just when `jobs` runs out.

```bash
mkdir -p ~/projects/level24-bonus2
cd ~/projects/level24-bonus2
go mod init level24.example/bonus2
```

**Hints:**
- Every place a worker or the job-feeder would otherwise block on a channel operation needs a second `select` case on `<-ctx.Done()`
- Use an infinite job generator (`for j := 1; ; j++`) so the *only* thing that ever stops the pool is the context timing out
- Confirm with `wg.Wait()` returning promptly after the timeout, and `ctx.Err()` reporting `context deadline exceeded`

### Challenge 3: First Successful Result Wins

Launch several goroutines that each simulate fetching from a different "mirror" with a different delay. Take the first one to respond, then `cancel()` the shared context so every other in-flight goroutine stops immediately instead of continuing to run for no reason.

```bash
mkdir -p ~/projects/level24-bonus3
cd ~/projects/level24-bonus3
go mod init level24.example/bonus3
```

**Hints:**
- One shared `context.WithCancel`, passed to every goroutine; a single shared, buffered `results` channel sized so a "losing" goroutine's send never blocks
- Each fetch goroutine should itself `select` on `ctx.Done()` so it stops waiting the instant a winner has already been chosen
- Use a `sync.WaitGroup` so `main` can confirm every loser actually observed the cancellation before the program exits, rather than guessing with a sleep
- The *order* the losers print "canceled" in in is not guaranteed - only that there is exactly one winner and the rest all eventually stop

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Create and distinguish context.Background() from context.TODO()
✅ Cancel a real running goroutine manually with context.WithCancel and verify the stop
✅ Use context.WithTimeout and context.WithDeadline with real, measured elapsed time
✅ Produce and distinguish context.Canceled from context.DeadlineExceeded via ctx.Err()
✅ Write the standard ctx.Done()-select pattern for racing real work against cancellation
✅ Use context.WithValue safely with a typed key, and explain the string-key collision risk it avoids
✅ Explain exactly what leaks when you skip defer cancel(), backed by a real go vet warning
✅ Propagate one context through a fan-out of goroutines that all stop together
✅ Build a cancelable multi-step pipeline that reports real partial progress

---

## Next Level

Level 25: Generics
- Type parameters for functions and types
- Type constraints and the comparable/any built-in constraints
- Writing generic data structures and algorithms
- When generics help vs. when a plain interface is simpler

Great work! You can now write Go concurrency that stops exactly when it's told to. 🚀
