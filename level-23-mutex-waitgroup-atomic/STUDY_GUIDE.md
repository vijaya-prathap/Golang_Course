# Level 23: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: Locking Fundamentals
```
Day 1:  sync.Mutex deep dive - the racy counter, fixed and -race-verified
Day 2:  The "always unlock" pattern (mu.Lock(); defer mu.Unlock())
Day 3:  sync.RWMutex - concurrent readers, exclusive writers
Day 4:  sync.WaitGroup deep dive - Add() timing, the by-value bug
Day 5:  sync/atomic - typed counters, lock-free
Day 6:  Atomics vs mutex - choosing the right tool
Day 7:  sync.Once and channels-vs-mutex philosophy
```

### Week 2: Practice & the Real Deadlock
```
Day 1:  Exercises 1-3 (mutex fix, defer pattern, RWMutex)
Day 2:  Exercises 4-5 (WaitGroup bugs, including the by-value deadlock)
Day 3:  Exercises 6-8 (atomics, atomic-vs-mutex, sync.Once)
Day 4:  Exercise 9 (lock-ordering deadlock, real detector output)
Day 5:  Exercise 10 (comprehensive cache/counter service)
Day 6:  Bonus challenges
Day 7:  Practice & review
```

---

## 🎯 sync.Mutex: The Critical Section

```
Without a mutex:                     With a mutex:

goroutine A: counter++  ─┐           goroutine A: mu.Lock() ──┐
goroutine B: counter++  ─┤ INTERLEAVE        counter++        │ EXCLUSIVE
goroutine C: counter++  ─┘  (lost updates)   mu.Unlock() ─────┘
                                      goroutine B: mu.Lock() ──┐
                                             counter++         │ waits, then runs
                                             mu.Unlock() ──────┘
                                      goroutine C: ... (same pattern)

Result: some increments silently     Result: EVERY increment counted -
        lost - "Got: 947" not 1000          "Got: 1000", every time
```

```go
mu.Lock()
// ---- critical section: only ONE goroutine at a time is ever in here ----
counter++
// --------------------------------------------------------------------
mu.Unlock()
```

**Verified:** the racy version lost increments on every one of 5 runs (976, 962, 947, 958, 945). The mutex-protected version hit exactly 1000 on every run, plain AND under `-race` (0 races found).

---

## 🔒 The "Always Unlock" Pattern

```
mu.Lock()
defer mu.Unlock()   ← placed IMMEDIATELY after Lock()
                       guaranteed to run on:
                         - normal return
                         - an early return (guard clause)
                         - a panic unwinding the stack
```

```
Decision: where does the critical section end?

"The rest of this function" ──► defer mu.Unlock() right after Lock()  (default choice)
"Just a few lines in the middle" ──► explicit Unlock() right after those lines,
                                       to release the lock earlier
```

---

## 🗂️ sync.RWMutex: Read/Write Access

```
Readers (RLock/RUnlock)              Writer (Lock/Unlock)

  R1 ──RLock──┐                        W1 ──Lock──┐
  R2 ──RLock──┼── all run CONCURRENTLY  (excludes  │  EXCLUSIVE -
  R3 ──RLock──┘   (no writer active)    everyone)  │  no readers, no other writers
                                        ────Unlock──┘

If a writer is waiting or active, new readers wait too - writers are
never starved indefinitely by a constant stream of readers.
```

| Situation | Use |
|-----------|-----|
| Reads vastly outnumber writes AND each read is substantial | `sync.RWMutex` |
| Critical section is trivially cheap either way | plain `sync.Mutex` may be simpler *and* faster |
| Unsure / haven't measured | start with `sync.Mutex`, switch only if profiling shows read contention |

**Verified (this project's own hardware, 50 readers × 2,000 reads each):**
```
RWMutex (read-heavy): 20.04475ms for 100000 reads
Mutex   (read-heavy): 13.599458ms for 100000 reads
```
Plain `Mutex` won here - the lookups were too cheap for `RWMutex`'s extra bookkeeping to pay off. Exact numbers vary by machine; the *lesson* (measure before assuming) doesn't.

---

## 🔁 sync.WaitGroup Lifecycle

```
        Add(n)                Done()  x n times           Wait()
          │                       │                          │
          ▼                       ▼                          ▼
   counter = n    goroutines finish, each  counter reaches 0 ─► Wait() returns
                  calls Done() (-1 each)

RULE: Add() must happen BEFORE the `go` statement that launches
      the goroutine which will eventually call Done().
```

```
❌ WRONG (Add inside the goroutine - races with Wait):

  for i := 0; i < 5; i++ {
      go func() {
          wg.Add(1)     ◄── Wait() below might check the counter
          defer wg.Done()    before this line ever executes!
      }()
  }
  wg.Wait()             ◄── can return with counter still at 0

✅ RIGHT (Add before `go`):

  for i := 0; i < 5; i++ {
      wg.Add(1)         ◄── happens synchronously, before launch
      go func() {
          defer wg.Done()
      }()
  }
  wg.Wait()             ◄── guaranteed to wait for all 5
```

**Verified:** the buggy version printed `completed when Wait() returned: 0` in all 13 test runs in this sandbox - `Wait()` kept winning the race against the goroutines' own `Add()` calls.

### The By-Value Trap

```
func worker(wg sync.WaitGroup) {   ◄── VALUE parameter = a COPY
    defer wg.Done()                     Done() decrements the COPY
}                                       the ORIGINAL never moves

main's wg.Add(1) ──► counter = 1
go worker(wg)         (copy made HERE, at the call)
main's wg.Wait() ──► blocks FOREVER (original counter stuck at 1)
                      ▼
     fatal error: all goroutines are asleep - deadlock!
```

**Verified verbatim, 3 runs, identical every time (memory addresses aside):**
```
processing: invoice-42
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [sync.WaitGroup.Wait]:
sync.runtime_SemacquireWaitGroup(...)
sync.(*WaitGroup).Wait(...)
main.main()
exit status 2
```
`go vet` catches this WITHOUT running anything:
```
main.go:11:21: processItem passes lock by value: sync.WaitGroup contains sync.noCopy
```

---

## ⚛️ Atomic vs Mutex Decision Tree

```
Do you need to change/read ONE independent variable
(a counter, a flag, a pointer)?
│
├─ YES ──► Is it JUST that one variable, with nothing else
│          that must stay consistent with it?
│          │
│          ├─ YES ──► sync/atomic (atomic.Int64, atomic.Bool, ...)
│          │           Lock-free, simplest, often fastest.
│          │
│          └─ NO, other state must move with it ──► sync.Mutex
│
└─ NO, it's a map/struct/multiple fields that must be
   read or updated together as one consistent unit ──► sync.Mutex
```

**Verified example:** `Stats{requests, errors}` (independent atomics: 200 / 20) vs. `Ledger{balance, totalFeesPaid}` (one mutex, both fields always consistent: 9800 / 200, sum always 10000).

### Typed Atomics Available (Go 1.19+, confirmed on this project's go1.26.5)

```go
var i  atomic.Int64      // i.Add(1), i.Load(), i.Store(n), i.CompareAndSwap(old, new)
var i32 atomic.Int32
var u  atomic.Uint64
var b  atomic.Bool        // b.Load(), b.Store(true), b.CompareAndSwap(false, true)
var p  atomic.Pointer[T]  // p.Load(), p.Store(&v), p.CompareAndSwap(old, new)
```

---

## 🔂 sync.Once Flow

```
goroutine 1: once.Do(fn) ──► fn hasn't run ──► RUNS fn ──► fn completes
goroutine 2: once.Do(fn) ──► BLOCKS until goroutine 1's fn finishes ──► returns, fn NOT run again
goroutine 3: once.Do(fn) ──► same as goroutine 2
   ...
goroutine 100: once.Do(fn) ──► same

Result: fn runs EXACTLY ONCE, no matter how many goroutines
        call once.Do(fn), no matter how concurrently they call it.
```

**Verified:** 100 concurrent callers, `"loading configuration..."` printed exactly once, `initCalls == 1`, all 100 goroutines received the identical `*Config` pointer - every run, plain and under `-race`.

---

## 💀 Lock-Ordering Deadlock

```
transferAtoB:  Lock(mu1) ──► [sleep] ──► Lock(mu2)  (waiting...)
transferBtoA:  Lock(mu2) ──► [sleep] ──► Lock(mu1)  (waiting...)

                    mu1 held by A, wanted by B
                    mu2 held by B, wanted by A
                         ╲           ╱
                          ╲  DEADLOCK ╱
                           ╲_________╱
              Neither goroutine can EVER proceed.
              main() is also stuck in wg.Wait().
              ALL goroutines asleep ──► Go's runtime detects this
              and crashes the process:

     fatal error: all goroutines are asleep - deadlock!
```

**Verified verbatim** (3 runs, always deadlocked) - the runtime dump names exactly which goroutine is stuck on which primitive:
```
goroutine 1 [sync.WaitGroup.Wait]: ... main.main()
goroutine 7 [sync.Mutex.Lock]:     ... main.transferAtoB(...)
goroutine 8 [sync.Mutex.Lock]:     ... main.transferBtoA(...)
```

**The fix:** pick ONE global lock order and never deviate from it.
```
❌ transferAtoB: mu1 → mu2        transferBtoA: mu2 → mu1   (opposite - deadlock risk)
✅ transferAtoB: mu1 → mu2        transferBtoA: mu1 → mu2   (same order - safe)
```

---

## 🚨 Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Forgetting to unlock | manual `Unlock()` that an early `return` skips | `mu.Lock(); defer mu.Unlock()` |
| Copying a Mutex/WaitGroup | `func f(wg sync.WaitGroup)` | `func f(wg *sync.WaitGroup)` |
| Add() racing with Wait() | `go func(){ wg.Add(1); ... }()` | `wg.Add(1); go func(){ ... }()` |
| Lock-ordering deadlock | two goroutines, opposite lock order | one consistent global lock order everywhere |
| Mutex where a channel fits better | polling a mutex-guarded "next job" variable | `jobs := make(chan *Job)` |

---

## 📈 Progression Summary

### Understanding Level 23

Level 23 turns Level 20's brief mutex/WaitGroup preview into full mastery, and adds the tools that preview didn't cover:

1. **sync.Mutex** — the critical section, verified `-race`-clean
2. **defer mu.Unlock()** — the idiomatic unlock pattern
3. **sync.RWMutex** — concurrent readers, exclusive writers
4. **sync.WaitGroup** — real rules: Add() timing, never by value
5. **sync/atomic** — lock-free typed counters
6. **Atomics vs mutex** — one variable vs. multiple related fields
7. **sync.Once** — guaranteed one-time initialization
8. **Channels vs mutex** — two valid philosophies, used together
9. **Lock-ordering deadlocks** — a real, captured `fatal error`

### Prerequisites for Level 24

Before moving to Level 24 (Context), you need:

- ✅ Comfortable locking and unlocking a `sync.Mutex` around a critical section
- ✅ Comfortable with `sync.RWMutex` for read-heavy shared state
- ✅ Know the WaitGroup rules: `Add()` before `go`, always by pointer
- ✅ Comfortable choosing between `sync/atomic` and `sync.Mutex`
- ✅ Understand `sync.Once` for one-time initialization
- ✅ Can read a real Go deadlock dump and explain what it means
- ✅ Completed Level 22 (Select)

### Ready for Level 24?

Level 24 teaches the tool Go programs use to control how long a goroutine runs and when it should stop, building directly on this level's synchronization primitives:
- `context.Context` for cancellation and deadlines
- Passing request-scoped values safely
- Combining context cancellation with the mutexes, WaitGroups, and channels you already know

---

## ✅ Checklist Before Level 24

- [ ] Can fix a racy counter with `sync.Mutex` and confirm it with `-race`
- [ ] Always write `mu.Lock(); defer mu.Unlock()` immediately after locking
- [ ] Can explain when `sync.RWMutex` helps vs. when plain `Mutex` is simpler/faster
- [ ] Know never to call `Add()` from inside the goroutine it's counting
- [ ] Know never to pass a `sync.Mutex` or `sync.WaitGroup` by value
- [ ] Can use a typed atomic (`atomic.Int64`, etc.) for a simple lock-free counter
- [ ] Can explain why multiple related fields need a mutex, not separate atomics
- [ ] Can use `sync.Once` for one-time initialization under concurrency
- [ ] Can explain a lock-ordering deadlock and how to prevent it
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The Critical Section
Everything between `Lock()` and `Unlock()` runs as if no other goroutine exists - that's the entire guarantee, and it's why keeping it small matters.

### The Unlock Discipline
`defer mu.Unlock()` right after `Lock()` is not a style preference - it's what makes forgetting to unlock structurally difficult.

### The WaitGroup Rules
`Add()` before `go`, never by value. Both mistakes were verified here to produce real, observable bugs - one silent (Wait returns early), one loud (a genuine runtime deadlock).

### Atomic or Mutex
One independent variable → atomic. Several fields that must move together → mutex. There is no in-between.

### The Deadlock Detector
Go's runtime can prove when literally nothing in the process can make progress - that's what `fatal error: all goroutines are asleep - deadlock!` means, and this level captured it twice, from two different root causes.

---

## 📚 Next Level

Level 24: Context
- `context.Context` for cancellation, deadlines, and timeouts
- Passing request-scoped values safely
- Controlling goroutine lifetimes on top of the primitives from this level

You've got Go's classic concurrency toolkit down! Keep going! 🚀
