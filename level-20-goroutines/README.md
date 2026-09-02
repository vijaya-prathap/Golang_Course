# Level 20: Goroutines - Complete Guide

## Introduction

Welcome to Level 20! You've learned how to read and write JSON (Level 19) - structuring data as it moves in and out of your programs. Now it's time to learn how Go runs **multiple things at once**: goroutines.

This is where Go starts to feel genuinely different from most languages you may have used before. Concurrency is built into the language itself, not bolted on as a library - and it's one of the biggest reasons Go became popular for servers, network tools, and anything that has to juggle many things happening independently.

This level is harder than what came before. Concurrency bugs don't always show up the same way twice, and that's exactly the point - you need to build the right mental model now, before you're debugging a flaky production issue at 2 AM.

---

## Table of Contents

1. [What Is a Goroutine?](#what-is-a-goroutine)
2. [The Classic Beginner Bug: main() Exits Too Early](#the-classic-beginner-bug-main-exits-too-early)
3. [sync.WaitGroup: Waiting for Goroutines to Finish](#syncwaitgroup-waiting-for-goroutines-to-finish)
4. [Closures and Goroutines: The Loop Variable Question](#closures-and-goroutines-the-loop-variable-question)
5. [Race Conditions](#race-conditions)
6. [Fixing Races: A Preview of sync.Mutex](#fixing-races-a-preview-of-syncmutex)
7. [Goroutine Leaks](#goroutine-leaks)
8. [When to Use Goroutines (and When Not To)](#when-to-use-goroutines-and-when-not-to)
9. [Best Practices](#best-practices)
10. [Common Mistakes](#common-mistakes)

---

## What Is a Goroutine?

A **goroutine** is a lightweight, independently executing function managed by the Go runtime - often described as a "green thread." You don't create operating system threads directly in Go; instead the Go runtime multiplexes potentially thousands of goroutines onto a much smaller number of real OS threads, scheduling them itself.

You launch one with the `go` keyword in front of a function call:

```go
go doSomething()
```

That's it. `doSomething()` now runs concurrently with the rest of `main()` - the `go` statement itself returns immediately, without waiting for `doSomething` to finish.

Goroutines are cheap compared to OS threads: a goroutine starts with a tiny stack (a few kilobytes) that grows as needed, whereas an OS thread typically reserves megabytes. This is why Go programs can comfortably run tens of thousands of goroutines where a thread-per-task design in other languages would fall over.

### Concurrency vs. Parallelism

These two words get used interchangeably in casual conversation, but they mean different things - and the distinction matters for understanding goroutines correctly.

- **Concurrency** is about *structure*: dealing with multiple independent tasks by interleaving their progress, whether or not they ever literally execute at the same instant. A single CPU core can be concurrent by rapidly switching between tasks.
- **Parallelism** is about *execution*: multiple tasks physically running at the same time, which requires multiple CPU cores.

> Concurrency is about **dealing with** lots of things at once. Parallelism is about **doing** lots of things at once.

Goroutines give you concurrency for free - the Go runtime schedules them, switching between them as needed. Whether that concurrency also becomes parallelism depends on `GOMAXPROCS` (the number of OS threads the runtime is allowed to run Go code on simultaneously, defaulting to the number of CPU cores) and how many CPU cores are actually available. On a multi-core machine, goroutines that aren't blocked on I/O or synchronization genuinely can run in parallel; on a single core, they're concurrent but not parallel - the runtime just switches between them very quickly.

You write concurrent code (independent goroutines) and let the Go runtime decide how much of it becomes parallel, based on the hardware it's running on.

---

## The Classic Beginner Bug: main() Exits Too Early

Here is the single most common mistake every Go beginner makes with goroutines: launching one and expecting its output to appear, when `main()` doesn't wait around for it.

```go
package main

import "fmt"

func sayHello() {
    fmt.Println("Hello from the goroutine!")
}

func main() {
    fmt.Println("main: launching goroutine")
    go sayHello()
    fmt.Println("main: returning immediately (no wait!)")
}
```

When `main()` returns, the entire program exits immediately - **every goroutine still running is simply killed**, no matter what it was doing. Go does not wait for background goroutines before shutting down.

**Actually run 8 times in a row on this project's environment:**

```
--run 1--
main: launching goroutine
main: returning immediately (no wait!)
--run 2--
main: launching goroutine
main: returning immediately (no wait!)
--run 3--
main: launching goroutine
main: returning immediately (no wait!)
--run 4--
main: launching goroutine
main: returning immediately (no wait!)
--run 5--
main: launching goroutine
main: returning immediately (no wait!)
Hello from the goroutine!
--run 6--
main: launching goroutine
main: returning immediately (no wait!)
--run 7--
main: launching goroutine
main: returning immediately (no wait!)
--run 8--
main: launching goroutine
main: returning immediately (no wait!)
```

**This is real, captured output - not a hypothetical.** Seven of eight runs never printed the goroutine's message at all: `main()` reached its last `fmt.Println` and the process exited before the scheduler ever got around to running `sayHello`. On run 5, the goroutine got lucky and squeezed its print in before `main()` returned. **The exact mix you get will differ run to run** - that unpredictability is the whole lesson. You cannot rely on a goroutine "probably having enough time."

The fix is never "make main sleep a little" (that's a band-aid that can still fail under load, and wastes time when the goroutine finishes early) - the real fix is to **synchronize explicitly**, which is exactly what `sync.WaitGroup` is for.

---

## sync.WaitGroup: Waiting for Goroutines to Finish

`sync.WaitGroup` lets `main()` (or any goroutine) block until a known number of other goroutines have finished. It has three methods:

- `Add(n)` - increment the counter by `n` (call this *before* launching the goroutines it counts)
- `Done()` - decrement the counter by 1 (call this when a goroutine finishes - almost always via `defer`)
- `Wait()` - block until the counter reaches zero

```go
package main

import (
    "fmt"
    "sync"
)

func sayHello(wg *sync.WaitGroup) {
    defer wg.Done()
    fmt.Println("Hello from the goroutine!")
}

func main() {
    var wg sync.WaitGroup

    fmt.Println("main: launching goroutine")
    wg.Add(1)
    go sayHello(&wg)

    wg.Wait()
    fmt.Println("main: goroutine finished, safe to return")
}
```

**Verified output (stable across 5 runs):**

```
main: launching goroutine
Hello from the goroutine!
main: goroutine finished, safe to return
```

Unlike the previous example, this output is **deterministic** - not because goroutine scheduling became predictable, but because `wg.Wait()` blocks `main()` until `wg.Done()` has actually run. The *order* of internal scheduling doesn't matter anymore; the *observable* result (the message prints before `main` continues) is now guaranteed.

A `sync.WaitGroup` must never be copied after first use - always pass it around by pointer (`*sync.WaitGroup`), as shown above.

### Waiting for Multiple Goroutines

`Add` accepts any count, and each goroutine calls `Done()` independently:

```go
var wg sync.WaitGroup
const numWorkers = 5

for i := 0; i < numWorkers; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        // work
    }()
}

wg.Wait() // blocks until all 5 have called Done()
```

Calling `wg.Add(1)` inside the loop, once per goroutine, is the idiomatic pattern - it keeps the "how many am I waiting for" logic next to the "launch one" logic.

---

## Closures and Goroutines: The Loop Variable Question

Launching goroutines inside a loop raises a question every Go programmer eventually runs into: does each goroutine see its own value of the loop variable, or do they all share one?

This is directly related to the closure-capture behavior covered in [Level 11's section on closures in loops](../level-11-functions/README.md#closures-in-loops-the-loop-variable-gotcha) - the short version is repeated here because it matters just as much for goroutines as it does for stored closures.

**Before Go 1.22**, a `for` loop reused a single loop variable across every iteration. A closure (including a goroutine launched with `go func() { ... }()`) that referenced the loop variable captured that *one shared variable* - so by the time the goroutine actually ran, the loop might have already moved on (or finished), and every goroutine could observe the *same*, often-final, value. This was one of the most notorious bugs in Go's history, and produced the well-known workaround of shadowing the variable inside the loop body (`i := i`) to give each iteration its own copy.

**As of Go 1.22**, `for` loops create a fresh instance of each loop variable on every iteration. This project's `go.mod` specifies `go 1.26.5` (confirmed by reading `/Users/macbookpro/go/Golang/go.mod`), so the old bug simply does not apply here.

```go
package main

import (
    "fmt"
    "sort"
    "sync"
)

func main() {
    var wg sync.WaitGroup
    var mu sync.Mutex
    var results []int

    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            mu.Lock()
            results = append(results, i) // captures THIS iteration's i
            mu.Unlock()
        }()
    }

    wg.Wait()
    sort.Ints(results)
    fmt.Println("Captured values (sorted):", results)
}
```

**Verified output on this project's Go version (go1.26.5), stable across 5 runs:**

```
Captured values (sorted): [0 1 2 3 4]
```

Every goroutine captured its own independent `i`, so all five values (0 through 4) show up exactly once. On Go 1.21 or earlier, this program could instead have printed something like `[5 5 5 5 5]` (or any repeated/partial mix, since it was also compounded by the data race of appending without a mutex).

If you read older tutorials, blog posts, or Stack Overflow answers warning you to write `go func(i int) { ... }(i)` (passing the loop variable in as a parameter to force a copy), that advice was correct for Go 1.21 and earlier. It's unnecessary on Go 1.22+, though still harmless if you see it in existing code - passing an already-independent variable as a parameter is a redundant copy, not a bug.

Note that the mutex (`sync.Mutex`) protecting `results` in the example above is doing separate, still-necessary work: even with correct per-iteration capture, multiple goroutines appending to the *same slice* at the same time is its own data race (covered next), unrelated to loop variable semantics.

---

## Race Conditions

A **data race** occurs when two or more goroutines access the same memory location concurrently, at least one of the accesses is a write, and there is no synchronization ordering the accesses. Data races are a category of bug distinct from a general "race condition" (a broader term for any outcome that depends on timing) - a data race specifically means unsynchronized concurrent memory access, and Go's tooling can detect it precisely.

Here's a shared counter incremented by 1,000 goroutines with no synchronization at all:

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var wg sync.WaitGroup
    counter := 0

    const numGoroutines = 1000
    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter++ // UNSYNCHRONIZED read-modify-write - a data race!
        }()
    }

    wg.Wait()
    fmt.Println("Expected:", numGoroutines, "Got:", counter)
}
```

`counter++` looks like one operation, but it's really three: read `counter`, add 1, write it back. When two goroutines interleave those three steps, one goroutine's increment can be silently lost.

**Actual output from 5 plain runs (no race detector) in this sandbox:**

```
Expected: 1000 Got: 970
Expected: 1000 Got: 987
Expected: 1000 Got: 958
Expected: 1000 Got: 943
Expected: 1000 Got: 943
```

Every single run lost increments, and the exact count lost is different every time - this is genuinely non-deterministic output, not a fixed "wrong" number. **Do not treat any specific "Got:" value as the expected output of this program** - the only correct expectation is "Got: 1000 is *not* guaranteed, and will usually be lower."

### Catching It With the Race Detector

Go ships a built-in race detector, enabled with `go run -race`. **This was actually attempted in this sandbox, and it worked** - it requires cgo, and cgo is available here. Real captured output:

```
==================
WARNING: DATA RACE
Read at 0x00c000012158 by goroutine 8:
  main.main.func1()
      /private/tmp/level20-verify/ex5/main.go:17 +0x68

Previous write at 0x00c000012158 by goroutine 7:
  main.main.func1()
      /private/tmp/level20-verify/ex5/main.go:17 +0x78

Goroutine 8 (running) created at:
  main.main()
      /private/tmp/level20-verify/ex5/main.go:15 +0x70

Goroutine 7 (finished) created at:
  main.main()
      /private/tmp/level20-verify/ex5/main.go:15 +0x70
==================
Expected: 1000 Got: 869
Found 2 data race(s)
exit status 1
```

The race detector points at the exact line (`counter++`) and even names the two goroutines whose unsynchronized read and write collided - this is dramatically more useful than just noticing the final count is wrong. `go run -race` exits with a non-zero status when it finds a race, which makes it suitable for CI pipelines. It does have real costs: instrumented programs run several times slower and use significantly more memory, so `-race` is a development/testing tool, not something you ship to production.

> **If `-race` isn't available in your environment:** it depends on cgo, which isn't available on every platform/configuration (some minimal or cross-compiled setups disable cgo). If `go run -race` fails for you with a cgo-related error, you can still fully understand the bug from the plain-run output above - the varying, too-low "Got:" counts already prove the race exists; the detector just makes it easier to *locate*.

---

## Fixing Races: A Preview of sync.Mutex

`sync.Mutex` ("mutual exclusion lock") ensures only one goroutine at a time can execute the code between `Lock()` and `Unlock()`. This is just enough to fix the race above - full coverage of mutexes, `RWMutex`, and the `sync/atomic` package is Level 23's job.

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var wg sync.WaitGroup
    var mu sync.Mutex
    counter := 0

    const numGoroutines = 1000
    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            mu.Lock()
            counter++
            mu.Unlock()
        }()
    }

    wg.Wait()
    fmt.Println("Expected:", numGoroutines, "Got:", counter)
}
```

**Verified output, stable across 5 plain runs AND a `-race` run (0 races found):**

```
Expected: 1000 Got: 1000
```

`mu.Lock()` and `mu.Unlock()` bracket the critical section (`counter++`) so that no two goroutines can be inside it at the same time - the read-modify-write is now effectively atomic from every other goroutine's point of view. Always `Unlock()` what you `Lock()`, ideally with `defer mu.Unlock()` immediately after locking when the critical section is the rest of the function, so a panic or an early return can't leave the mutex held forever.

---

## Goroutine Leaks

A **goroutine leak** happens when a goroutine blocks forever and can never finish - unlike a leaked object in a garbage-collected language, a leaked goroutine is never cleaned up by the garbage collector, because from the runtime's point of view it's still "doing work" (even though that work will never complete). Leaked goroutines accumulate over the lifetime of a long-running program (like a server) and quietly consume memory until something breaks.

The most common cause is sending to (or receiving from) a channel when nobody on the other end will ever show up:

```go
package main

import (
    "fmt"
    "time"
)

func leaky(ch chan int) {
    val := 42
    fmt.Println("leaky goroutine: about to send, but nobody is receiving...")
    ch <- val // blocks forever - nobody ever reads from ch
    fmt.Println("leaky goroutine: this line never runs")
}

func main() {
    ch := make(chan int) // unbuffered: a send blocks until a receiver is ready

    fmt.Println("main: launching a goroutine that will block forever")
    go leaky(ch)

    select {
    case <-time.After(200 * time.Millisecond):
        fmt.Println("main: gave up waiting - the goroutine above is now leaked")
    }

    fmt.Println("main: returning; the leaked goroutine dies with the process")
}
```

**Verified output, stable across 3 runs:**

```
main: launching a goroutine that will block forever
leaky goroutine: about to send, but nobody is receiving...
main: gave up waiting - the goroutine above is now leaked
main: returning; the leaked goroutine dies with the process
```

`leaky` is permanently stuck on `ch <- val` because `ch` is unbuffered and nothing ever calls `<-ch`. `main` doesn't wait for it - it gives up after a 200ms timeout (using `time.After`, a small preview of the channel operations Level 21 covers in depth) and moves on. In this exercise the leaked goroutine is harmless because the whole process exits right after, which kills every goroutine along with it. In a real long-running program (a server that never exits), a leak like this would sit there forever, one instance accumulating per request that hit this code path, gradually consuming memory - that's the actual danger of a goroutine leak.

The fix is always the same idea: know how every goroutine you start is going to *stop*. Common tools: a `context.Context` with cancellation (Level 24), a done channel, or a `select` with a timeout like the one above on the *sending* side too, so it gives up instead of blocking forever.

---

## When to Use Goroutines (and When Not To)

**Good fits for goroutines:**
- **I/O-bound work** - network requests, file reads, database queries - where a goroutine spends most of its time waiting, not computing. While one goroutine waits on a slow network call, the runtime can run others.
- **Fan-out patterns** - splitting one big task into independent pieces (e.g., processing chunks of a slice, or querying several downstream services) that can run concurrently and then have their results collected.
- **Background work** that doesn't need to block the caller - logging, metrics, cache warming.

**Poor fits for goroutines:**
- **Tiny, trivial work.** Launching a goroutine has real (if small) overhead - stack allocation, scheduler bookkeeping - and every goroutine that touches shared state needs synchronization, which has its own cost and complexity. For something as cheap as `n * n` on a handful of numbers, a goroutine is very likely to be *slower* than just doing it inline, once you count the overhead. Concurrency should make a real workload faster or better structured - it shouldn't be reached for by default.
- **Work with no independence.** If step B always needs step A's result, there's nothing to run concurrently - a goroutine wouldn't change anything except add synchronization overhead.
- **Anywhere the added complexity isn't paying for itself.** Every goroutine you launch is a piece of program state (is it running? blocked? did it finish? did it error?) that you now have to reason about. Sequential code you can read top-to-bottom is easier to get right - only pay the concurrency complexity tax when there's a real payoff (usually: waiting on something slow, or genuinely independent parallel work).

The bonus challenges include a benchmark comparing sequential vs. goroutine-based versions of the same work, so you can see this tradeoff with real numbers instead of just taking it on faith.

---

## Best Practices

### 1. Always Know How a Goroutine Will Stop

Before writing `go someFunc()`, be able to answer: "under what condition does this goroutine's function return?" If the answer is "it doesn't, unless X happens" - make sure X is guaranteed to eventually happen.

### 2. Use WaitGroup or Channels to Synchronize - Never time.Sleep

```go
// ❌ Fragile - guessing at timing
go doWork()
time.Sleep(100 * time.Millisecond)
fmt.Println("probably done?")

// ✅ Correct - explicit synchronization
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    doWork()
}()
wg.Wait()
fmt.Println("definitely done")
```

### 3. Avoid Unbounded Goroutine Creation

Launching one goroutine per item in a small, fixed-size slice is fine. Launching one goroutine per incoming request on a server with no limit, or one per row in an unbounded database query, can exhaust memory under load. Bound concurrency with a pattern like the semaphore shown in the bonus challenges, or a fixed-size worker pool (a full treatment arrives with buffered channels in Level 21).

### 4. Protect Every Piece of Shared, Mutable State

If more than one goroutine can read *and* at least one can write the same variable, slice, or map, it needs a `sync.Mutex` (or another synchronization primitive) around every access - not just some of them. A single unprotected access is enough to make the whole thing a data race.

### 5. Run `go run -race` (or `go test -race`) During Development

The race detector doesn't catch every possible race (it can only report races it actually observes during that particular run), but it's extremely effective at catching the interleavings that do happen to occur. Make it part of your normal testing habit once your code touches goroutines.

---

## Common Mistakes

### Mistake 1: Forgetting to Wait for Goroutines

```go
// ❌ WRONG - main returns before the goroutine can print anything
func main() {
    go fmt.Println("hello")
}

// ✅ RIGHT - wait for it explicitly
func main() {
    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        defer wg.Done()
        fmt.Println("hello")
    }()
    wg.Wait()
}
```

### Mistake 2: Assuming the Pre-1.22 Loop Variable Bug Still Applies

```go
// This is CORRECT on Go 1.22+ (including this project's go 1.26.5) -
// each goroutine gets its own i automatically.
for i := 0; i < 5; i++ {
    go func() {
        fmt.Println(i)
    }()
}

// This extra parameter is now unnecessary (but harmless) on Go 1.22+:
for i := 0; i < 5; i++ {
    go func(i int) {
        fmt.Println(i)
    }(i)
}
```

If you're reading material written before early 2024, assume it's describing the *old* behavior unless it says otherwise. Always check your own `go.mod` version before trusting inherited assumptions about loop variable capture.

### Mistake 3: Data Races From Unsynchronized Shared State

```go
// ❌ WRONG - counter++ is not atomic; concurrent goroutines lose updates
counter := 0
for i := 0; i < 1000; i++ {
    go func() {
        defer wg.Done()
        counter++
    }()
}

// ✅ RIGHT - protect the shared variable
var mu sync.Mutex
for i := 0; i < 1000; i++ {
    go func() {
        defer wg.Done()
        mu.Lock()
        counter++
        mu.Unlock()
    }()
}
```

Remember: a data race can exist even if a particular run happens to print the "right" answer. Absence of a visibly wrong result on one run is not proof of correctness - run `-race` to be sure.

### Mistake 4: Goroutine Leaks From Blocking Sends/Receives With No Escape

```go
// ❌ WRONG - if nobody ever reads from ch, this goroutine blocks forever
go func() {
    ch <- computeResult()
}()

// ✅ RIGHT - give it a way out (a timeout, a done channel, or a context)
go func() {
    select {
    case ch <- computeResult():
    case <-time.After(time.Second):
        return // give up instead of blocking forever
    }
}()
```

### Mistake 5: Reaching for Goroutines for Trivially Small Work

```go
// ❌ WRONG - the goroutine + mutex overhead likely costs more than
// the tiny amount of work it's parallelizing
go func() {
    mu.Lock()
    total += a + b
    mu.Unlock()
}()

// ✅ RIGHT - just do it inline
total += a + b
```

---

## Summary

**What a Goroutine Is:**
- A lightweight, runtime-managed "green thread" launched with `go someFunc()`
- Concurrency (structuring independent tasks) is not the same as parallelism (running them simultaneously on multiple cores) - goroutines give you the former, and the runtime turns it into the latter when hardware allows

**Synchronization:**
- `main()` returning kills every goroutine still running - never assume a background goroutine "had enough time"
- `sync.WaitGroup` (`Add`/`Done`/`Wait`) is the standard way to wait for a known number of goroutines to finish

**Closures in Loops:**
- On Go 1.22+ (this project uses go1.26.5), each loop iteration gets its own copy of the loop variable - goroutines launched in a loop capture their own value correctly, with no extra workaround needed

**Data Races:**
- Unsynchronized concurrent access to shared memory, with at least one write, is a data race - and produces genuinely non-deterministic, often silently wrong, results
- `go run -race` catches races in the act and points at the exact lines involved
- `sync.Mutex` (`Lock`/`Unlock`) fixes a race by ensuring only one goroutine touches the shared state at a time (full coverage in Level 23)

**Goroutine Leaks:**
- A goroutine blocked forever (e.g., on a channel send nobody will ever receive) is never garbage collected - always know how a goroutine will stop

**When to Use Them:**
- Good fits: I/O-bound work, fan-out over independent tasks, background work
- Poor fits: trivial work, strictly sequential dependencies, anything where the complexity outweighs the benefit

---

## Next Steps

You now understand:
- ✅ What a goroutine is, and concurrency vs. parallelism
- ✅ Why `main()` can exit before a goroutine's output ever appears
- ✅ How to synchronize with `sync.WaitGroup`
- ✅ Why loop variable capture in goroutines is safe on this project's Go version
- ✅ What a data race is, how to see one with `-race`, and how to fix one with `sync.Mutex`
- ✅ What a goroutine leak is and why it's never garbage collected
- ✅ When goroutines are worth the complexity, and when they aren't

**Next level:** Level 21 - Channels
- Sending and receiving values between goroutines
- Buffered vs. unbuffered channels
- `select` for waiting on multiple channels
- Building on the WaitGroup and mutex patterns from this level with a more idiomatic, channel-based communication style

You've just started writing genuinely concurrent programs - one of Go's biggest strengths. Keep going! 🚀
