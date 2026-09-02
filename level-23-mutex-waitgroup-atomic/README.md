# Level 23: Mutex, WaitGroup & Atomic - Complete Guide

## Introduction

Welcome to Level 23! You've learned to coordinate goroutines with channels and `select` (Level 22) - Go's "share memory by communicating" style of concurrency. Now it's time to go deep on the *other* half of Go's concurrency toolkit: classic locking and atomic operations from the `sync` and `sync/atomic` packages.

Level 20 gave you a brief preview of `sync.WaitGroup` and a glimpse of `sync.Mutex` fixing a racy counter. This level is where that preview becomes mastery: `sync.Mutex` and `sync.RWMutex` in full depth, `sync.WaitGroup`'s real rules (and its sharpest gotchas), the modern typed `sync/atomic` API, `sync.Once`, and a real, captured deadlock from getting lock ordering wrong. Everything here is verified by actually running the code in this project's sandbox - including a genuine `fatal error: all goroutines are asleep - deadlock!` from Go's own runtime detector.

---

## Table of Contents

1. [sync.Mutex Deep Dive](#syncmutex-deep-dive)
2. [The "Always Unlock" Pattern](#the-always-unlock-pattern)
3. [sync.RWMutex: Readers and Writers](#syncrwmutex-readers-and-writers)
4. [sync.WaitGroup Deep Dive](#syncwaitgroup-deep-dive)
5. [sync/atomic: Lock-Free Counters](#syncatomic-lock-free-counters)
6. [Atomics vs Mutex: Which One?](#atomics-vs-mutex-which-one)
7. [sync.Once: Guaranteed One-Time Initialization](#synconce-guaranteed-one-time-initialization)
8. [Channels vs Mutex: Two Philosophies](#channels-vs-mutex-two-philosophies)
9. [A Real Deadlock: Lock Ordering Gone Wrong](#a-real-deadlock-lock-ordering-gone-wrong)
10. [Best Practices](#best-practices)
11. [Common Mistakes](#common-mistakes)

---

## sync.Mutex Deep Dive

A `sync.Mutex` ("mutual exclusion lock") guarantees that only one goroutine at a time can be inside the code between `Lock()` and `Unlock()` - the **critical section**. Every other goroutine that calls `Lock()` while it's held simply blocks until the current holder calls `Unlock()`.

Level 20 previewed this fixing a racy counter. Here it is again, verified fresh in this level's own sandbox, first broken and then fixed, so you can see the actual before/after side by side.

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    fmt.Println("=== Racy Version (No Synchronization) ===")
    racyCounter := 0
    var wg1 sync.WaitGroup
    const n = 1000
    for i := 0; i < n; i++ {
        wg1.Add(1)
        go func() {
            defer wg1.Done()
            racyCounter++ // data race: read-modify-write with no protection
        }()
    }
    wg1.Wait()
    fmt.Println("Expected:", n, "Got:", racyCounter)

    fmt.Println("\n=== Fixed Version (sync.Mutex) ===")
    var mu sync.Mutex
    fixedCounter := 0
    var wg2 sync.WaitGroup
    for i := 0; i < n; i++ {
        wg2.Add(1)
        go func() {
            defer wg2.Done()
            mu.Lock()
            fixedCounter++
            mu.Unlock()
        }()
    }
    wg2.Wait()
    fmt.Println("Expected:", n, "Got:", fixedCounter)
}
```

**Actual output from 5 plain runs (no race detector):**

```
--run 1--
=== Racy Version (No Synchronization) ===
Expected: 1000 Got: 976

=== Fixed Version (sync.Mutex) ===
Expected: 1000 Got: 1000
--run 2--
=== Racy Version (No Synchronization) ===
Expected: 1000 Got: 962

=== Fixed Version (sync.Mutex) ===
Expected: 1000 Got: 1000
--run 3--
=== Racy Version (No Synchronization) ===
Expected: 1000 Got: 947

=== Fixed Version (sync.Mutex) ===
Expected: 1000 Got: 1000
--run 4--
=== Racy Version (No Synchronization) ===
Expected: 1000 Got: 958

=== Fixed Version (sync.Mutex) ===
Expected: 1000 Got: 1000
--run 5--
=== Racy Version (No Synchronization) ===
Expected: 1000 Got: 945

=== Fixed Version (sync.Mutex) ===
Expected: 1000 Got: 1000
```

The racy version lost a different number of increments every single run (976, 962, 947, 958, 945) - never once hitting 1000. The mutex-protected version hit exactly 1000 every time, because `mu.Lock()`/`mu.Unlock()` around `fixedCounter++` means no two goroutines can ever be inside that increment at the same time.

### This Time, With `-race` - and It's Actually Clean

**This was attempted for real in this sandbox with `go run -race`, and `-race` is available here** (it depends on cgo, and cgo works in this environment). Running the racy version under `-race` immediately reports real data races:

```
=== Racy Version (No Synchronization) ===
==================
WARNING: DATA RACE
Read at 0x00c00019c038 by goroutine 13:
  main.main.func1()
      /private/tmp/level23-verify/ex1/main.go:17 +0x68

Previous write at 0x00c00019c038 by goroutine 8:
  main.main.func1()
      /private/tmp/level23-verify/ex1/main.go:17 +0x78
...
==================
Expected: 1000 Got: 915
```

...and then, immediately after, the **same run** shows the mutex-protected version reporting **zero races** and the exact count:

```
=== Fixed Version (sync.Mutex) ===
Expected: 1000 Got: 1000
```

`Found 2 data race(s)` and a non-zero exit status came only from the racy half of the program. This is the concrete meaning of "`-race`-clean": run the mutex-protected code under `go run -race`, and the race detector finds nothing to report, on top of the count already being exactly right.

---

## The "Always Unlock" Pattern

`mu.Lock(); defer mu.Unlock()` immediately after the lock is the idiomatic Go pattern for any critical section that spans "the rest of the function." This directly reuses `defer`'s core guarantee from [Level 11's defer section](../level-11-functions/README.md#defer) and its use for cleanup in [Level 16's error handling](../level-16-error-handling/README.md): a deferred call runs when the function returns, **no matter how it returns** - a normal `return`, an early `return` from a guard clause, or even a `panic` unwinding the stack.

```go
type BankAccount struct {
    mu      sync.Mutex
    balance int
    history []string
}

func (a *BankAccount) Withdraw(amount int) bool {
    a.mu.Lock()
    defer a.mu.Unlock() // guaranteed to run, even on the early return below
    if amount > a.balance {
        a.history = append(a.history, fmt.Sprintf("rejected withdrawal %d (insufficient funds)", amount))
        return false // WITHOUT defer, a naive `a.mu.Unlock()` here would be easy to forget
    }
    a.balance -= amount
    a.history = append(a.history, fmt.Sprintf("withdraw %d", amount))
    return true
}
```

Without `defer`, every early `return` (and there's often more than one, as above) needs its own manual `Unlock()` call - miss just one, and every future caller blocks on `Lock()` forever. `defer mu.Unlock()` written on the line right after `mu.Lock()` makes it structurally impossible to forget, and makes the paired Lock/Unlock visually obvious at a glance.

**Verified: 50 concurrent deposits of 10 and 50 concurrent withdrawals of 5 against a starting balance of 1000, stable across 3 plain runs and a `-race` run (0 races found):**

```
Final balance: 1250
Expected balance: 1250
History entries recorded: 100
Expected history entries: 100
```

Both `balance` and `history` moved together consistently on every single one of the 100 concurrent operations - exactly the guarantee a mutex-protected critical section provides.

The one time `defer mu.Unlock()` is *not* automatically the right call: when the critical section is only a few lines in the *middle* of a longer function, and you want the lock released well before the function returns. In that case, an explicit `mu.Unlock()` right after the critical section (not deferred) avoids holding the lock longer than necessary. Deferring is the default; narrowing the critical section further is an optimization for when the section is genuinely small relative to the rest of the function.

---

## sync.RWMutex: Readers and Writers

A plain `sync.Mutex` allows exactly one goroutine in the critical section at a time - even if every single goroutine only wants to *read*. `sync.RWMutex` adds a second mode: `RLock()`/`RUnlock()` for readers, which allows **any number of readers to hold the lock simultaneously**, as long as no writer holds it. `Lock()`/`Unlock()` still work as before for writers, and a writer excludes everyone (readers and other writers) while it holds the lock.

```go
type Cache struct {
    mu   sync.RWMutex
    data map[string]int
}

func (c *Cache) Get(key string) (int, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    v, ok := c.data[key]
    return v, ok
}

func (c *Cache) Set(key string, value int) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = value
}
```

**Verified with 20 concurrent readers and 3 occasional writers, stable across 3 plain runs and a `-race` run (0 races found):**

```
Reads completed: 20 (expected 20)
Writes performed: 3 (expected 3)
Final value of x is one of {2,3,4}: true
```

All 20 reads succeeded concurrently with each other, and the 3 writes each took exclusive access in turn - the final value is whichever of the 3 writers happened to run last, which is expected and correct.

### When RWMutex Helps vs. When Plain Mutex Is Simpler

`RWMutex` is a genuine win when reads vastly outnumber writes **and** each read does enough work that letting readers run concurrently actually matters (e.g., reading a moderately expensive-to-copy value, or holding the lock for more than a trivial number of instructions). It is *not* automatically faster - `RWMutex` has more internal bookkeeping than a plain `Mutex`, so for very cheap, very short critical sections, that extra bookkeeping can outweigh the benefit of concurrent readers.

This is not a hand-wavy claim - the bonus challenges include an informal benchmark that measured this directly on this project's own hardware:

```
NOTE: exact timings vary by machine and load - the point is the shape, not the numbers.
RWMutex (read-heavy): 20.04475ms for 100000 reads
Mutex   (read-heavy): 13.599458ms for 100000 reads
```

On this run, the plain `Mutex` was actually *faster* than the `RWMutex` for 100,000 trivial map-lookup reads with 50 concurrent readers - the lookup itself is so cheap that `RWMutex`'s extra reader-tracking overhead lost out to `Mutex`'s simplicity. **Don't reach for `RWMutex` by default just because a workload is "read-heavy" in name** - reach for it when reads are both frequent *and* substantial enough that concurrent read access is worth the added bookkeeping. When in doubt, start with a plain `Mutex` (it's simpler to reason about) and switch to `RWMutex` only if profiling shows read contention is a real bottleneck.

---

## sync.WaitGroup Deep Dive

Level 20 introduced the three methods - `Add(n)`, `Done()`, `Wait()` - as a way to make `main()` wait for goroutines. Here's the full picture, including the two mistakes that bite hardest in real code.

### Mistake: Calling Add() From Inside the Goroutine

`Add()` must be called **before** the goroutine that will call `Done()` is launched - not from inside that goroutine. Calling it from inside creates a race between the `Add()` call and `Wait()`: if `Wait()` happens to run before any goroutine's `Add()` has executed, `Wait()` can see a counter that is still zero and **return immediately**, without waiting for anything.

```go
fmt.Println("=== Correct: Add() called BEFORE launching each goroutine ===")
var wg sync.WaitGroup
for i := 0; i < 5; i++ {
    wg.Add(1) // called from the loop, before `go`
    go func(id int) {
        defer wg.Done()
        // ... work ...
    }(i)
}
wg.Wait()

fmt.Println("=== Buggy: Add() called INSIDE the goroutine (races with Wait) ===")
var wg2 sync.WaitGroup
for i := 0; i < 5; i++ {
    go func(id int) {
        wg2.Add(1) // WRONG: main's Wait() below might already be checking the counter
        defer wg2.Done()
        // ... work ...
    }(i)
}
wg2.Wait()
```

**Verified output, stable across 13 total runs of the buggy version in this sandbox:**

```
=== Correct: Add() called BEFORE launching each goroutine ===
completed: 5 (always 5 - deterministic)

=== Buggy: Add() called INSIDE the goroutine (races with Wait) ===
completed when Wait() returned: 0 (NOT reliably 5 - this is the bug)
```

Every one of the 13 runs against this sandbox's scheduler had `Wait()` return before any of the 5 goroutines got around to calling `wg2.Add(1)` - `main` moved on immediately, believing (wrongly) that there was nothing to wait for. **The exact number you'll see is scheduler-dependent** - a separate test of the same underlying bug pattern, run 5 times, printed `0` four times and `5` once, when a goroutine happened to win the race to `Add()` before `Wait()` checked. The unreliability itself - not any specific printed number - is the actual bug.

The fix is always the pattern shown in the "correct" half above: call `Add()` synchronously, in the same goroutine that will do the `go` statement, before it.

### Mistake: Passing a WaitGroup By Value

A `sync.WaitGroup` must never be copied after it's been used - it contains an internal counter and semaphore state that a copy separates from the original. Passing one **by value** into a function is exactly this mistake, and it's easy to write by accident:

```go
// BUG: wg is a sync.WaitGroup VALUE parameter, so every call gets its own COPY.
// Done() decrements the copy inside the function, never the caller's original -
// so the caller's wg.Wait() below can never see the counter reach zero.
func processItem(wg sync.WaitGroup, item string) {
    defer wg.Done()
    fmt.Println("processing:", item)
}

func main() {
    var wg sync.WaitGroup
    wg.Add(1)
    go processItem(wg, "invoice-42") // wg is COPIED into processItem here
    wg.Wait()                        // blocks forever - the ORIGINAL counter is still 1
    fmt.Println("all done")
}
```

**`go vet` catches this automatically, without even running the program:**

```
main.go:11:21: processItem passes lock by value: sync.WaitGroup contains sync.noCopy
main.go:19:17: call of processItem copies lock value: sync.WaitGroup contains sync.noCopy
```

`sync.WaitGroup` embeds a `noCopy` marker specifically so `go vet` can detect exactly this mistake at analysis time. `go vet` is *not* run automatically as part of `go run` for this particular check, so the program still compiles and runs - which is why it's worth running `go vet ./...` as a habit, not just relying on `go build`/`go run` to catch everything.

**And here is what actually happens if you run it anyway - captured verbatim from 3 real runs, every one identical (only the memory addresses differ):**

```
processing: invoice-42
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [sync.WaitGroup.Wait]:
sync.runtime_SemacquireWaitGroup(0x1049cbad0?, 0x18?)
	/Users/macbookpro/.local/go/src/runtime/sema.go:114 +0x38
sync.(*WaitGroup).Wait(0x6a2e65b88020)
	/Users/macbookpro/.local/go/src/sync/waitgroup.go:206 +0xa8
main.main()
	/private/tmp/level23-verify/ex5/main.go:20 +0x74
exit status 2
```

`processItem` printed its line and returned normally - its `defer wg.Done()` decremented its own *copy* of the WaitGroup down to zero, harmlessly, in a copy nobody else can see. Meanwhile the original `wg` in `main` never moved off `1`, so `wg.Wait()` blocked forever. Once the goroutine finished, only `main` was left, permanently blocked - the Go runtime's deadlock detector noticed **all** goroutines were asleep with no possibility of progress and crashed the process itself with `fatal error: all goroutines are asleep - deadlock!`, rather than hanging silently forever.

The fix: always pass a `*sync.WaitGroup` (a pointer), never a `sync.WaitGroup` by value - exactly as Level 20 already showed for the correct usage.

---

## sync/atomic: Lock-Free Counters

`sync/atomic` provides operations that are indivisible at the hardware level - no goroutine can ever observe a torn, half-updated value, and no `Lock()`/`Unlock()` is needed at all. This project's `go.mod` specifies `go 1.26.5` (confirmed by reading `/Users/macbookpro/go/Golang/go.mod`), so the modern **typed atomic types** added in Go 1.19 are fully available and are the idiomatic choice here: `atomic.Int32`, `atomic.Int64`, `atomic.Uint32`, `atomic.Uint64`, `atomic.Bool`, and `atomic.Pointer[T]` (a generic pointer type, usable because this project's Go version also has full generics support from Level 25's era). Each typed atomic wraps a raw memory location and exposes methods directly on it - no more manually passing `*int64` addresses to free functions like the older `atomic.AddInt64(&x, 1)` style required.

```go
var counter atomic.Int64

counter.Add(1)          // atomically add 1
counter.Load()           // atomically read the current value
counter.Store(0)         // atomically set to an exact value
counter.CompareAndSwap(0, 5) // atomically: if current value == 0, set it to 5; returns whether it swapped
```

```go
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

func main() {
    var wg sync.WaitGroup
    var counter atomic.Int64

    const numGoroutines = 1000
    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter.Add(1)
        }()
    }

    wg.Wait()
    fmt.Println("Expected:", numGoroutines, "Got:", counter.Load())
}
```

**Verified output, stable across 5 plain runs AND a `-race` run (0 races found):**

```
Expected: 1000 Got: 1000
```

Every single run hit exactly 1000 - no mutex was needed anywhere in this program. Compare this to the *un-guarded* `counter++` from the very first example in this level, which lost increments on every run. `atomic.Int64.Add` performs the entire read-modify-write as one indivisible hardware operation, so there's no window where two goroutines can interleave and stomp on each other's update.

A second exercise pushed this further with 5,000 concurrent increments and got the identical result every time:

```
Expected: 5000 Got: 5000
MATCH: atomic counter is exact under concurrent increments
```

---

## Atomics vs Mutex: Which One?

Atomics and mutexes solve overlapping problems, but they aren't interchangeable in general - the deciding question is **how many pieces of state need to change together, as one consistent unit**.

- **A single, independent variable** (a counter, a flag, a pointer being swapped) is exactly what atomics are for. There's nothing else that needs to be consistent with it, so an indivisible hardware operation is enough on its own.
- **Multiple related fields that must always be seen/updated together** (a balance and a transaction log, a map plus a size field, a struct with several fields that must stay mutually consistent) need a mutex. You cannot make two *separate* atomic operations happen "at the same instant" - between them, another goroutine can observe a half-updated, inconsistent state. A mutex's critical section is the tool for "these N operations happen as one atomic unit," no matter how many variables are involved.

```go
// Good fit for atomic: one independent counter, nothing else to keep in sync.
type Stats struct {
    requests atomic.Int64
    errors   atomic.Int64
}

// Good fit for mutex: two fields that must always be updated TOGETHER.
type Ledger struct {
    mu            sync.Mutex
    balance       int
    totalFeesPaid int
}

func (l *Ledger) ChargeFee(fee int) {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.balance -= fee
    l.totalFeesPaid += fee // both fields move together - can't be two separate atomics safely
}
```

**Verified with 200 concurrent goroutines hitting each type, stable across 3 plain runs and a `-race` run (0 races found):**

```
=== atomic: independent counters ===
requests: 200 errors: 20

=== mutex: fields that must move together ===
balance: 9800 (expected 9800)
totalFeesPaid: 200 (expected 200)
balance + totalFeesPaid + net change matches original: true
```

If `balance` and `totalFeesPaid` were each their own separate `atomic.Int64` instead of being protected together by one mutex, a goroutine could observe `balance` already decremented but `totalFeesPaid` not yet incremented (or vice versa) - a momentarily inconsistent snapshot. The mutex's critical section is what makes "subtract the fee AND record it" happen as one atomic-from-the-outside unit, which two independent atomics cannot express.

**Rule of thumb:** one counter/flag → atomic. Two or more things that must never be observed out of sync with each other → mutex.

---

## sync.Once: Guaranteed One-Time Initialization

`sync.Once` guarantees that a function runs **exactly once**, no matter how many goroutines call `once.Do(fn)` concurrently - every caller blocks until the first call to actually run `fn` has completed, and every caller after that returns immediately without running `fn` again.

```go
var (
    once      sync.Once
    cfg       *Config
    initCalls atomic.Int32
)

func loadConfig() *Config {
    once.Do(func() {
        initCalls.Add(1)
        fmt.Println("loading configuration from disk (expensive, should happen once)...")
        cfg = &Config{value: "production"}
    })
    return cfg
}
```

**Verified with 100 concurrent goroutines all calling `loadConfig()`, stable across 5 plain runs and a `-race` run (0 races found):**

```
loading configuration from disk (expensive, should happen once)...
initCalls (actual runs of the init function): 1 (expected 1)
all 100 goroutines got the SAME *Config instance: true
```

The line `loading configuration from disk...` printed exactly once per run, every run - not once-ish, not "usually once." All 100 goroutines received the identical `*Config` pointer, proving they all saw the result of the *same* initialization, not 100 separate ones racing each other. `sync.Once` is the standard idiom for lazy singletons, one-time setup (opening a shared resource, compiling a regex, seeding a cache) that might be reached from multiple goroutines before you can guarantee it's already been done.

---

## Channels vs Mutex: Two Philosophies

Level 21 introduced Go's guiding proverb:

> Don't communicate by sharing memory; share memory by communicating.

That's a philosophy, not a law - real Go code uses both styles, often side by side, and this level's own examples reach for a mutex constantly. The practical guidance:

- **Reach for a channel** when the real problem is *handing off* a value or a unit of work between goroutines, or coordinating a sequence of events (a pipeline stage finishing, a worker picking up the next job, a signal that something happened). Channels make the *flow* of data and control between goroutines part of the type system and the code's shape.
- **Reach for a mutex (or atomic)** when the real problem is *protecting a piece of shared state that several goroutines need to read and write directly* - a cache, a counter, an in-memory map, a struct multiple goroutines all hold a reference to. There's no "message" being passed here; it's genuinely one piece of state with multiple concurrent accessors, which is exactly what `sync.Mutex` was built for.

Trying to force the "wrong" tool onto a problem usually shows up as awkward code: passing every single field access through a channel to a dedicated "owner" goroutine is possible but often more ceremony than a mutex, for state that's naturally just... shared. Conversely, orchestrating a multi-stage worker pipeline entirely with a shared mutex-protected struct (instead of channels connecting stages) tends to produce code that's harder to read than the equivalent channel pipeline. This level's final exercise - a cache/counter service - is a case where a mutex protecting the map is simply the more natural and idiomatic fit than routing every `Get`/`Set` through a channel to an owner goroutine.

---

## A Real Deadlock: Lock Ordering Gone Wrong

A **lock-ordering deadlock** happens when two goroutines each need two locks, but acquire them in opposite order. If goroutine A grabs lock 1 and then waits for lock 2, while goroutine B has already grabbed lock 2 and is waiting for lock 1, neither can ever proceed - each is holding what the other needs.

```go
var mu1, mu2 sync.Mutex

// transferAtoB locks mu1 then mu2.
func transferAtoB(wg *sync.WaitGroup) {
    defer wg.Done()
    mu1.Lock()
    fmt.Println("transferAtoB: locked mu1")
    time.Sleep(100 * time.Millisecond) // widens the window so BOTH goroutines commit to their first lock
    fmt.Println("transferAtoB: waiting for mu2")
    mu2.Lock()
    fmt.Println("transferAtoB: locked mu2")
    mu2.Unlock()
    mu1.Unlock()
}

// transferBtoA locks mu2 then mu1 - OPPOSITE ORDER. This is the bug.
func transferBtoA(wg *sync.WaitGroup) {
    defer wg.Done()
    mu2.Lock()
    fmt.Println("transferBtoA: locked mu2")
    time.Sleep(100 * time.Millisecond)
    fmt.Println("transferBtoA: waiting for mu1")
    mu1.Lock()
    fmt.Println("transferBtoA: locked mu1")
    mu1.Unlock()
    mu2.Unlock()
}

func main() {
    var wg sync.WaitGroup
    wg.Add(2)
    go transferAtoB(&wg)
    go transferBtoA(&wg)
    wg.Wait()
    fmt.Println("both transfers completed")
}
```

**This was actually run 3 times in this sandbox. Every run deadlocked, and the runtime caught it every time - here is one full run, captured verbatim:**

```
transferBtoA: locked mu2
transferAtoB: locked mu1
transferAtoB: waiting for mu2
transferBtoA: waiting for mu1
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [sync.WaitGroup.Wait]:
sync.runtime_SemacquireWaitGroup(0x103033ad0?, 0x80?)
	/Users/macbookpro/.local/go/src/runtime/sema.go:114 +0x38
sync.(*WaitGroup).Wait(0x3ef5046f20f0)
	/Users/macbookpro/.local/go/src/sync/waitgroup.go:206 +0xa8
main.main()
	/private/tmp/level23-verify/ex9/main.go:42 +0xb8

goroutine 7 [sync.Mutex.Lock]:
internal/sync.runtime_SemacquireMutex(0x3ef5046e2040?, 0x0?, 0x1e?)
	/Users/macbookpro/.local/go/src/runtime/sema.go:95 +0x28
internal/sync.(*Mutex).lockSlow(0x103194f58)
	/Users/macbookpro/.local/go/src/internal/sync/mutex.go:149 +0x170
internal/sync.(*Mutex).Lock(...)
	/Users/macbookpro/.local/go/src/internal/sync/mutex.go:70
sync.(*Mutex).Lock(...)
	/Users/macbookpro/.local/go/src/sync/mutex.go:46
main.transferAtoB(0x0?)
	/private/tmp/level23-verify/ex9/main.go:18 +0x184
created by main.main in goroutine 1
	/private/tmp/level23-verify/ex9/main.go:40 +0x70

goroutine 8 [sync.Mutex.Lock]:
internal/sync.runtime_SemacquireMutex(0x3ef5046e2040?, 0xa0?, 0x1e?)
	/Users/macbookpro/.local/go/src/runtime/sema.go:95 +0x28
internal/sync.(*Mutex).lockSlow(0x103194f50)
	/Users/macbookpro/.local/go/src/internal/sync/mutex.go:149 +0x170
internal/sync.(*Mutex).Lock(...)
	/Users/macbookpro/.local/go/src/internal/sync/mutex.go:70
sync.(*Mutex).Lock(...)
	/Users/macbookpro/.local/go/src/sync/mutex.go:46
main.transferBtoA(0x0?)
	/private/tmp/level23-verify/ex9/main.go:31 +0x184
created by main.main in goroutine 1
	/private/tmp/level23-verify/ex9/main.go:41 +0xb0
exit status 2
```

Read the goroutine dump like a diagnosis: goroutine 1 (`main`) is asleep in `wg.Wait()`, goroutine 7 (`transferAtoB`) is asleep trying to `Lock()` a mutex it can't get, and goroutine 8 (`transferBtoA`) is asleep trying to `Lock()` a *different* mutex it can't get. **Every** goroutine in the program is permanently blocked, with no external event (no timer, no other running goroutine) that could ever unblock any of them - that's precisely what "all goroutines are asleep - deadlock!" means, and it's why Go's runtime can detect and crash on it reliably rather than hanging forever: it can prove that literally nothing in the process is runnable.

This is a genuinely different situation from the WaitGroup-by-value deadlock earlier - here, *two* mutexes and *two* worker goroutines are all stuck waiting on each other, not one goroutine stuck on a WaitGroup whose counter can never move. Both count as real, detector-caught deadlocks; the mechanism differs.

**The fix:** always acquire multiple locks in a single, globally consistent order everywhere in the codebase (e.g., always lock the lower-numbered account ID first, or always lock in a fixed, documented field order for a struct with multiple mutexes). If every goroutine agrees on the order, this specific deadlock pattern becomes structurally impossible.

---

## Best Practices

### 1. Keep Critical Sections Small

Hold a lock for the minimum amount of work necessary - don't do I/O, logging, or unrelated computation while holding a mutex. A smaller critical section means less contention and fewer chances for another slow operation to block everyone else waiting on the same lock.

```go
// ❌ Doing unrelated slow work while holding the lock
mu.Lock()
data[key] = value
time.Sleep(time.Second) // nothing to do with the critical section - blocks every other goroutine needlessly
mu.Unlock()

// ✅ Lock only what needs protecting
mu.Lock()
data[key] = value
mu.Unlock()
time.Sleep(time.Second) // do unrelated work outside the lock
```

### 2. Always Unlock via defer

`mu.Lock(); defer mu.Unlock()` immediately after acquiring the lock, as covered above, makes forgetting to unlock structurally difficult - even across early returns and panics.

### 3. Prefer Atomics for Simple Single-Variable Counters and Flags

If the state in question is one independent counter, one flag, or one pointer being swapped, `sync/atomic`'s typed types (`atomic.Int64`, `atomic.Bool`, etc.) are simpler and typically faster than a mutex protecting the same single variable - there's no lock to forget to release, because there's no lock at all.

### 4. Be Consistent About Lock Ordering

If any part of the program ever needs to hold two or more locks at once, fix a single global order for acquiring them (by variable name, by an ID, by field order in a struct) and follow it *everywhere* that combination of locks is taken. Inconsistent ordering is the root cause of every lock-ordering deadlock, including the one demonstrated above.

### 5. Never Copy a Mutex or WaitGroup After First Use

Always store them as struct fields or pass them by pointer (`*sync.Mutex`, `*sync.WaitGroup`). `go vet` will flag most accidental copies (as shown above) - don't ignore that warning.

### 6. Run `go vet` in Addition to `go build`/`go run`

`go build` and `go run` only run a small, automatic subset of vet checks. The full `go vet ./...` catches additional problems (like passing a lock by value into a function) that the automatic subset does not always block on.

---

## Common Mistakes

### Mistake 1: Forgetting to Unlock

```go
// ❌ WRONG - if an early return is added later, Unlock() is easy to skip
func (a *BankAccount) Withdraw(amount int) bool {
    a.mu.Lock()
    if amount > a.balance {
        return false // BUG: mutex is never unlocked - every future caller deadlocks
    }
    a.balance -= amount
    a.mu.Unlock()
    return true
}

// ✅ RIGHT - defer guarantees the unlock on every path
func (a *BankAccount) Withdraw(amount int) bool {
    a.mu.Lock()
    defer a.mu.Unlock()
    if amount > a.balance {
        return false
    }
    a.balance -= amount
    return true
}
```

### Mistake 2: Copying a Mutex or WaitGroup by Value

```go
// ❌ WRONG - wg is copied; Done() inside never affects the caller's original
func processItem(wg sync.WaitGroup, item string) {
    defer wg.Done()
}

// ✅ RIGHT - pass a pointer
func processItem(wg *sync.WaitGroup, item string) {
    defer wg.Done()
}
```

As verified above, the by-value version doesn't just misbehave subtly - it causes a genuine `fatal error: all goroutines are asleep - deadlock!` once the caller's `Wait()` can never see the counter reach zero.

### Mistake 3: Calling Add() Concurrently With Wait()

```go
// ❌ WRONG - Add() races with Wait(); Wait() can return before any Add() ran
for i := 0; i < 5; i++ {
    go func() {
        wg.Add(1) // too late - Wait() might already be checking
        defer wg.Done()
    }()
}
wg.Wait()

// ✅ RIGHT - Add() happens before the goroutine is even launched
for i := 0; i < 5; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
    }()
}
wg.Wait()
```

### Mistake 4: Lock-Ordering Deadlocks

```go
// ❌ WRONG - opposite acquisition order between the two goroutines
func transferAtoB() { mu1.Lock(); mu2.Lock(); /* ... */ }
func transferBtoA() { mu2.Lock(); mu1.Lock(); /* ... */ } // deadlock risk

// ✅ RIGHT - always acquire in the same global order
func transferAtoB() { mu1.Lock(); mu2.Lock(); /* ... */ }
func transferBtoA() { mu1.Lock(); mu2.Lock(); /* ... */ } // same order, both times
```

### Mistake 5: Using a Mutex Where a Channel Would Say It Better

```go
// ❌ AWKWARD - a shared "next job" variable, polled under a mutex
var mu sync.Mutex
var nextJob *Job
// producer: mu.Lock(); nextJob = job; mu.Unlock()
// consumer: loops, locking/checking/unlocking, polling for nextJob != nil

// ✅ CLEARER - a channel expresses "hand this job to whoever's ready" directly
jobs := make(chan *Job)
// producer: jobs <- job
// consumer: job := <-jobs
```

If the real shape of the problem is "one goroutine hands something to another," a channel usually reads better than a mutex-protected variable being polled or checked.

---

## Summary

**sync.Mutex:**
- `Lock()`/`Unlock()` bracket a critical section so only one goroutine at a time can execute it
- The Level 20 racy counter, fixed with a mutex, is verified `-race`-clean in this level: 1000/1000, every run, zero races reported

**The Unlock Pattern:**
- `mu.Lock(); defer mu.Unlock()` immediately after locking guarantees unlock on every code path, reusing `defer`'s guarantee from Level 11/16

**sync.RWMutex:**
- `RLock()`/`RUnlock()` allow concurrent readers; `Lock()`/`Unlock()` remain exclusive for writers
- Helps when reads are both frequent and substantial; a plain `Mutex` can be simpler *and* faster for very cheap critical sections, as this level's own benchmark showed

**sync.WaitGroup:**
- `Add()` must happen before the `go` statement, never from inside the goroutine (racing with `Wait()` otherwise)
- Never pass a `WaitGroup` by value - `go vet` can catch it, and doing it anyway produces a real, verified `fatal error: all goroutines are asleep - deadlock!`

**sync/atomic:**
- Typed atomics (`atomic.Int64`, `atomic.Bool`, etc.) give lock-free, indivisible operations for single variables - verified exact (1000/1000, then 5000/5000) under concurrent increments where the unguarded version was never exact

**Atomics vs Mutex:**
- One independent counter/flag → atomic. Multiple fields that must move together as one unit → mutex.

**sync.Once:**
- Guarantees a function runs exactly once even under concurrent calls - verified: 1 init call, 1 shared instance, across 100 concurrent callers

**Channels vs Mutex:**
- Channels for handing off values/coordinating flow; mutexes for protecting shared state accessed directly by multiple goroutines - both are idiomatic Go, used side by side

**Deadlocks:**
- A lock-ordering deadlock (two goroutines, two locks, opposite order) is real, reproducible, and caught by Go's runtime deadlock detector - verified with the exact `fatal error: all goroutines are asleep - deadlock!` text and full goroutine dump

---

## Next Steps

You now understand:
- ✅ `sync.Mutex` and the critical section it protects, with a verified `-race`-clean fix of Level 20's racy counter
- ✅ Why `defer mu.Unlock()` is the idiomatic pattern, tying back to Level 11/16's `defer` coverage
- ✅ `sync.RWMutex` and when it helps vs. when a plain `Mutex` is simpler and just as fast
- ✅ `sync.WaitGroup`'s real rules: `Add()` before `go`, never by value - both verified with real bugs, including a genuine captured deadlock
- ✅ The modern typed `sync/atomic` API for lock-free counters and flags
- ✅ When to reach for an atomic vs. a mutex
- ✅ `sync.Once` for guaranteed one-time initialization under concurrency
- ✅ The "share memory by communicating" philosophy vs. classic locking, and when each fits better
- ✅ A real lock-ordering deadlock, captured verbatim from Go's own runtime detector

**Next level:** Level 24 - Context
- `context.Context` for cancellation, deadlines, and timeouts
- Passing request-scoped values safely
- Building on this level's synchronization primitives to control goroutine lifetimes explicitly

You now have the full classic concurrency toolkit - channels, mutexes, atomics, and the coordination primitives that hold them together. Keep going! 🚀
