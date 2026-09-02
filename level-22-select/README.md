# Level 22: Select - Complete Guide

## Introduction

Welcome to Level 22! Channels (Level 21) taught you how one goroutine hands a value to another safely. But real programs rarely wait on just one channel - a server might be waiting for a request, a shutdown signal, and a timeout all at once. What do you do when you need to wait on **several** channel operations at the same time, and proceed with whichever one happens first?

Go's answer is the `select` statement. It looks like a `switch`, but it behaves nothing like one: instead of comparing a value against cases, `select` waits on multiple channel operations simultaneously and runs whichever one becomes ready first. This level is marked **Advanced** because `select` is where several of Go's concurrency primitives - channels, goroutines, timers - come together, and because a few of its rules (random tie-breaking, `default`, the timer-leak gotcha) trip up even experienced Go programmers who assume it behaves like ordinary control flow.

---

## Table of Contents

1. [What select Does](#what-select-does)
2. [Random Tie-Breaking](#random-tie-breaking)
3. [The default Case](#the-default-case)
4. [Timeouts With time.After](#timeouts-with-timeafter)
5. [The done-Channel Cancellation Pattern](#the-done-channel-cancellation-pattern)
6. [select With Send Cases](#select-with-send-cases)
7. [select on a Closed Channel](#select-on-a-closed-channel)
8. [The time.After Leak Gotcha](#the-timeafter-leak-gotcha)
9. [Combining select With a for Loop: an Event Multiplexer](#combining-select-with-a-for-loop-an-event-multiplexer)
10. [Best Practices](#best-practices)
11. [Common Mistakes](#common-mistakes)

---

## What select Does

`select` waits on multiple channel operations at once and proceeds with whichever one is ready first. Each `case` in a `select` is a send or a receive on a channel - not a value comparison like `switch`.

```go
select {
case msg1 := <-ch1:
    fmt.Println("received:", msg1)
case msg2 := <-ch2:
    fmt.Println("received:", msg2)
}
```

`select` blocks until **at least one** of its cases can proceed, then runs that case's body. If two goroutines are racing to send on `ch1` and `ch2`, whichever one actually delivers a value first is the case `select` runs - `select` doesn't know or care which case is "first" in the source code.

```go
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
```

```
received: from ch1
received: from ch2
done
```

`ch1`'s goroutine sleeps for less time, so it's ready first - `select` picks it up on the first loop iteration, then picks up `ch2` on the second. Looping around a `select` is how you consume more than one ready value: each pass through the loop is a fresh wait on all the cases.

---

## Random Tie-Breaking

Here's the rule that makes `select` fundamentally different from a `switch` or an `if`/`else if` chain: **when multiple cases are ready at the same time, Go picks one of them at random.** It is not top-to-bottom like a `switch`, and it is not first-declared-wins. Every ready case has a roughly equal chance of being chosen on any given `select`.

Why? If `select` always favored the first ready case in source order, programs would start silently depending on case *order* for correctness - a fragile, invisible coupling. Random selection forces you to treat every ready case as equally likely, which is what actually happens in a truly concurrent system.

Let's verify this empirically instead of taking it on faith. Run a `select` with two channels that are **both already ready** (pre-filled buffered channels), many times in a loop, and tally which case fires:

```go
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
```

Captured from a real run:

```
chA fired: 49986 times (50.0%)
chB fired: 50014 times (50.0%)
(exact numbers vary between runs, but both should land close to 50%)
```

Both channels are filled before the `select` ever runs, so both cases are ready on every single iteration - and yet neither one dominates. Three separate runs during verification landed at 49.8%/50.2%, 49.9%/50.1%, and 50.0%/50.0%: the exact split moves around slightly every time you run it (it's a real coin flip made 100,000 times), but it always stays close to an even 50/50 split. If `select` favored `chA` because it's listed first, you'd see something close to 100%/0% instead - that never happens.

---

## The default Case

Add a `default` case to make `select` **non-blocking**: if no other case is ready right away, `default` runs immediately instead of waiting.

```go
select {
case v := <-ch:
    fmt.Println("got", v)
default:
    fmt.Println("nothing ready")
}
```

Without `default`, a `select` with no ready case blocks until one becomes ready - potentially forever. With `default`, it never blocks at all: it either takes a ready case or falls through to `default` on the spot.

```go
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
```

```
=== Non-blocking receive from an empty channel ===
no value ready, default case ran

=== Non-blocking send into a full buffered channel ===
channel full, default case ran

=== Non-blocking send into a channel with room ===
sent 42, no default needed
```

The first `ch` has no sender and nothing buffered, so the receive case can't proceed - `default` runs instantly rather than blocking forever. The second channel is a full buffer, so the send case can't proceed either - same result. Only the third `select`, where the channel genuinely has room, takes its non-`default` case.

---

## Timeouts With time.After

`time.After(d)` returns a channel that receives a single value after duration `d` has elapsed. Used as a `select` case, it turns "wait for this channel, but give up after a while" into a few lines of code - one of the most common `select` patterns in real Go programs.

```go
select {
case result := <-ch:
    fmt.Println("got result:", result)
case <-time.After(2 * time.Second):
    fmt.Println("timed out")
}
```

```go
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
```

Captured from a real run:

```
=== Case 1: the operation finishes before the timeout ===
received "work done" before the timeout

=== Case 2: the timeout wins ===
timed out (elapsed roughly 100ms, bounded by the 100ms time.After)
```

In Case 1, the worker finishes in 50ms - well inside the 200ms budget - so the real result wins. In Case 2, the worker takes 300ms but the timeout is only 100ms, so `time.After`'s case fires first and the (still-running) worker's result is simply never received. The elapsed time printed is bounded by the 100ms timeout by construction; the exact number can drift a few milliseconds either way depending on your machine and system load, though the rounding above kept it at exactly 100ms across three separate runs during verification.

---

## The done-Channel Cancellation Pattern

A long-running goroutine needs a way to be told "stop now" from the outside. The idiomatic pre-`context` solution is a **done channel**: a `chan struct{}` (or `chan bool`) that a goroutine watches inside a `select`, alongside whatever real work it's doing. Closing the done channel is the signal - it makes every future receive on it immediately ready (see [select on a Closed Channel](#select-on-a-closed-channel)), so *every* goroutine watching it notices at once, no matter how many there are.

```go
select {
case <-done:
    return // told to stop
case w := <-work:
    // do something with w
}
```

```go
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
```

```
worker: processing 1
worker: processing 2
worker: processing 3
worker: received done signal, stopping
main: worker confirmed stopped
```

The worker's `for { select { ... } }` loop checks `done` on every single pass. As long as `main` keeps sending, the `work` case wins; the moment `main` calls `close(done)`, the worker's very next `select` sees `done` ready and returns. `main` waits on `finished` to know for certain the worker has actually exited before printing its own final line - without that synchronization, `main` could print "confirmed stopped" before the worker goroutine had actually finished.

This exact shape - a goroutine `select`ing between real work and a cancellation signal - is the direct precursor to Go's `context` package (covered in a later level, Level 24), which formalizes cancellation, deadlines, and timeouts into a single reusable `context.Context` value instead of a hand-rolled `done` channel. Everything you just did by hand here, `context.WithCancel` and `context.WithTimeout` do for you - but understanding the raw channel version first is what makes `context` make sense later.

---

## select With Send Cases

A `select` case isn't limited to receiving - a **send** can be a case too: `case ch <- v:`. The case is only "ready" if that send could complete right now (an unbuffered channel with a receiver waiting, or a buffered channel with room). This lets a goroutine try to hand off a value *or* bail out via some other channel, whichever is possible first.

```go
select {
case ch <- v:
    // the send succeeded - someone (or something) was ready to receive
case <-done:
    // told to stop before we ever got to send
    return
}
```

```go
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
```

```
=== Send case is NOT ready: channel full, done wins ===
done fired instead, since the send case was not ready
channel still holds: 100

=== Send case IS ready: channel has room, send wins ===
sent 42 - the channel had room, so the send case fired
received back: 42
```

In the first block, `ch` already holds its one buffered value, so `ch <- 200` cannot complete right now - it isn't a candidate at all, and the already-closed `done` case is the only one that can fire. In the second block, `ch2` has room and `done2` is never closed, so the send goes through. A send case behaves exactly like a receive case in every other way: it's just one more channel operation `select` can wait on and choose between.

---

## select on a Closed Channel

A receive case on a **closed** channel is a special kind of "always ready": it never blocks, and it immediately returns the channel's zero value (with `ok == false` in the comma-ok form, once any buffered values are drained). That makes it behave exactly like any other ready case as far as `select` is concerned - `select` will happily pick it, over and over, as fast as it's asked to.

```go
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
```

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

This is exactly why the [done-channel pattern](#the-done-channel-cancellation-pattern) works: closing `done` makes its receive case permanently ready, so every goroutine `select`ing on it notices - immediately, and forever after. But it's a double-edged sword: a `select` case reading a channel that's already closed and drained is ready on **every single call**, with zero waiting. A loop that keeps hitting that case without ever checking `ok` and exiting will spin at full CPU, doing no useful work - this is the busy-loop risk to watch for whenever a channel in a `select` might end up closed.

---

## The time.After Leak Gotcha

`time.After(d)` is convenient, but calling it **inside a loop** is a real, well-known gotcha: every call creates a brand-new `time.Timer` under the hood, and that timer isn't eligible for garbage collection until it actually fires - even though only the *current* iteration's timer is ever useful. In a loop that runs many times before its exit condition is met, or a loop with a short body and a long timeout, this steadily piles up live timers that just sit there, doing nothing, waiting to expire.

```go
// ANTI-PATTERN: a new timer is created and abandoned every iteration
for {
    select {
    case v := <-ch:
        process(v)
    case <-time.After(200 * time.Millisecond):
        fmt.Println("no value in time")
    }
}
```

The fix: create **one** `time.NewTimer` outside the loop, and `Stop()` + `Reset()` it each iteration instead of calling `time.After` again:

```go
timer := time.NewTimer(200 * time.Millisecond)
defer timer.Stop()

for {
    if !timer.Stop() {
        // drain a pending fire so Reset starts from a clean slate
        select {
        case <-timer.C:
        default:
        }
    }
    timer.Reset(200 * time.Millisecond)

    select {
    case v := <-ch:
        process(v)
    case <-timer.C:
        fmt.Println("no value in time")
    }
}
```

Verified side by side - both loops receive the exact same three values and see the exact same output, but only the second one reuses a single timer:

```go
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
}
```

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
```

The visible output is identical - that's the point. Nothing about the *behavior* changes; what changes is that the naive version silently created and abandoned three separate timers over the life of the loop (one per iteration), while the better version created exactly one and reused it. For a loop that only runs three times, the difference is trivial - but the same naive code in a server's hot request-handling loop, running thousands or millions of times, quietly accumulates thousands or millions of abandoned timers instead of one. `time.NewTimer` + `Stop()`/`Reset()` is the idiomatic fix, and it costs only a few extra lines.

---

## Combining select With a for Loop: an Event Multiplexer

Wrapping a `select` in a `for` loop turns a one-shot "wait for whichever is ready" into a standing **multiplexer**: a single point that keeps servicing several channels for as long as the program needs it to, and stops cleanly when told to.

One extra trick makes multiplexing multiple *sources* clean: once a source channel is exhausted (closed and drained), set your local variable for it to `nil` instead of continuing to select on the closed channel. A `nil` channel case in a `select` is never ready - it's permanently disabled, not "always ready" like a closed one - so this removes that source from consideration without an `if` check inside every case, and sidesteps exactly the busy-loop risk described in the [previous section](#select-on-a-closed-channel).

```go
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
```

```
from A: 1
from B: 100
channel B exhausted, disabling its case (setting to nil)
from A: 2
channel A exhausted, disabling its case (setting to nil)
both channels exhausted and disabled; total values received: 3
```

The loop condition `a != nil || b != nil` is the multiplexer's stopping condition: once both sources are disabled, there's nothing left to wait for, so the loop ends on its own. Note that this same `nil`-to-disable trick works for **any** case you want to turn off permanently at runtime - not just an exhausted source, but a `done` channel you only want to check once, or an output channel you only want open during part of the multiplexer's life. This confirms exactly what Level 21 previewed: a `nil` channel isn't broken, it's a channel `select` will simply never choose - which is precisely how you switch a case off without restructuring the whole `select`.

---

## Best Practices

### 1. Always Give an Unbounded select Loop a Way Out

If a `for { select { ... } }` loop is meant to run indefinitely, make sure at least one case is a stopping signal (a `done` channel, a counter check, a timeout) - otherwise the only way out is killing the process.

```go
// ✅ Good - has a real exit path
for {
    select {
    case v := <-work:
        process(v)
    case <-done:
        return
    }
}
```

### 2. Never Assume Case Order Means Priority

Don't write a `select` assuming the first case "wins" when multiple are ready - it doesn't. If you genuinely need priority between two channels, check the higher-priority one first with its own non-blocking `select`/`default`, then fall through to a blocking `select` on the rest.

```go
// ✅ Good - explicitly prioritizes highPriority over the rest
select {
case v := <-highPriority:
    handle(v)
default:
    select {
    case v := <-highPriority:
        handle(v)
    case v := <-normal:
        handle(v)
    }
}
```

### 3. Reuse a time.NewTimer Inside Loops, Not time.After

As shown above, `time.After` inside a loop is fine for a handful of iterations but leaks timers under sustained use. Reach for `time.NewTimer` with `Stop()`/`Reset()` in any loop that runs more than a few times or runs for a long time.

### 4. Use a nil Channel to Cleanly Disable a Case

Rather than adding `if` guards inside a case body to ignore values from an exhausted or not-yet-ready source, set that channel variable to `nil`. `select` will simply never choose a `nil` case again - it's the idiomatic way to "turn off" part of a multiplexer at runtime.

### 5. Prefer a Send-Case + done Over a Bare Blocking Send

If a goroutine's send might have no one left to receive it (a consumer that already gave up, for example), pair the send with a `done` case so the goroutine can bail out instead of blocking forever.

```go
select {
case out <- result:
    // delivered
case <-done:
    return // nobody's listening anymore
}
```

---

## Common Mistakes

### Mistake 1: Forgetting default Makes select Block

```go
// ❌ WRONG ASSUMPTION - without default, this BLOCKS if ch has nothing ready
select {
case v := <-ch:
    fmt.Println(v)
}

// ✅ RIGHT - add default for a non-blocking check
select {
case v := <-ch:
    fmt.Println(v)
default:
    fmt.Println("nothing ready")
}
```

### Mistake 2: Misunderstanding Tie-Breaking as Ordered

```go
// ❌ WRONG ASSUMPTION - "ch1 is listed first, so it always wins if both are ready"
select {
case v := <-ch1:
    fmt.Println("ch1:", v)
case v := <-ch2:
    fmt.Println("ch2:", v)
}
// If both are ready at the same instant, Go picks ONE AT RANDOM -
// case order in the source has no effect on which one fires.
```

### Mistake 3: Leaking Timers in a Loop

```go
// ❌ WRONG - a brand-new timer is created and abandoned every iteration
for {
    select {
    case v := <-ch:
        process(v)
    case <-time.After(time.Second):
        fmt.Println("no value")
    }
}

// ✅ RIGHT - one timer, reused with Stop()/Reset()
timer := time.NewTimer(time.Second)
defer timer.Stop()
for {
    if !timer.Stop() {
        select {
        case <-timer.C:
        default:
        }
    }
    timer.Reset(time.Second)
    select {
    case v := <-ch:
        process(v)
    case <-timer.C:
        fmt.Println("no value")
    }
}
```

### Mistake 4: Forgetting a Way to Ever Stop the Loop

```go
// ❌ WRONG - no done case, no counter, nothing - runs forever with no exit
for {
    select {
    case v := <-work:
        process(v)
    }
}

// ✅ RIGHT - always include a real stopping condition
for {
    select {
    case v := <-work:
        process(v)
    case <-done:
        return
    }
}
```

### Mistake 5: Busy-Looping on an Already-Closed Channel

```go
// ❌ WRONG - once ch is closed, this case is ALWAYS ready, spinning at full CPU
for {
    select {
    case v := <-ch:
        fmt.Println(v) // keeps printing the zero value forever, as fast as possible
    }
}

// ✅ RIGHT - check ok and stop (or nil out the channel) once it's closed
for {
    select {
    case v, ok := <-ch:
        if !ok {
            return // or: ch = nil to disable this case instead
        }
        fmt.Println(v)
    }
}
```

---

## Summary

**What select Does:**
- Waits on multiple channel sends/receives at once, proceeds with whichever is ready first
- With no ready case and no `default`, `select` blocks - just like a single channel op would

**Tie-Breaking:**
- Multiple ready cases → Go picks one **at random**, verified empirically at a roughly even split
- Case order in the source has zero effect on which one is chosen

**default:**
- Makes `select` non-blocking - runs immediately if nothing else is ready

**Timeouts:**
- `case <-time.After(d):` as one arm of a `select` is the standard timeout pattern
- Inside a loop, prefer `time.NewTimer` + `Stop()`/`Reset()` over repeated `time.After` calls

**Cancellation:**
- A `done` channel, closed once, makes every `select` watching it fire immediately - the precursor to `context`

**Closed Channels:**
- A closed channel's receive case is permanently ready - useful for `done`, dangerous as an unguarded busy-loop source

**nil Channels:**
- A `nil` channel's case is never ready - the idiomatic way to disable one branch of a `select` at runtime

**Patterns:**
- `select` + `for` = an event multiplexer that services several channels until told to stop

---

## Next Steps

You now understand:
- ✅ What `select` does and how it differs from `switch`
- ✅ Random tie-breaking among ready cases, confirmed with a real measured tally
- ✅ Non-blocking channel operations with `default`
- ✅ Timeouts with `time.After`, and the timer-leak gotcha with the `time.NewTimer` fix
- ✅ The done-channel cancellation pattern, and how it foreshadows `context`
- ✅ Send cases, closed-channel cases, and `nil` channels for disabling a case
- ✅ Building a multiplexer by combining `select` with a `for` loop

**Next level:** Level 23 - Mutex, WaitGroup & Atomic
- `sync.Mutex` and `sync.RWMutex` for protecting shared state directly
- `sync.WaitGroup` in depth for waiting on many goroutines
- The `sync/atomic` package for lock-free counters and flags
- Choosing between channels, mutexes, and atomics for a given problem

You can now listen to several conversations at once. Keep going! 🚀
