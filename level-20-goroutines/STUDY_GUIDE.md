# Level 20: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: Goroutine Fundamentals
```
Day 1:  What a goroutine is; concurrency vs parallelism
Day 2:  The main-exits-too-early bug (Exercise 1)
Day 3:  sync.WaitGroup (Exercises 2-3)
Day 4:  Closures and loop variable capture (Exercise 4)
Day 5:  Data races and go run -race (Exercise 5)
Day 6:  Fixing races with sync.Mutex (Exercise 6)
Day 7:  Goroutine leaks (Exercise 7)
```

### Week 2: Patterns & Practice
```
Day 1:  Fan-out worker pattern (Exercise 8)
Day 2:  Mutex-protected collections (Exercise 9)
Day 3-4: Comprehensive fetch simulation (Exercise 10)
Day 5:  Bonus challenges
Day 6-7: Review & consolidation
```

---

## 🎯 Goroutine Lifecycle Diagram

```
                     go someFunc()
                          │
                          ▼
                 ┌─────────────────┐
                 │   RUNNABLE       │  ← queued, waiting for the
                 │  (in scheduler   │    Go runtime to give it
                 │   run queue)     │    a thread to run on
                 └────────┬─────────┘
                          │  scheduler assigns an OS thread
                          ▼
                 ┌─────────────────┐
                 │    RUNNING       │  ← actually executing on
                 │                  │    a CPU core right now
                 └────────┬─────────┘
               ┌──────────┼───────────┐
     blocks on │          │           │ finishes normally
   channel/lock│          │ preempted │ (function returns)
               ▼          ▼           ▼
       ┌──────────┐ ┌───────────┐ ┌──────────┐
       │ WAITING   │ │ RUNNABLE  │ │  DONE    │
       │ (blocked) │ │(re-queued)│ │ (exited) │
       └─────┬─────┘ └───────────┘ └──────────┘
             │ unblocked (e.g. lock acquired,
             │ channel ready, WaitGroup released)
             └──────────────► back to RUNNABLE

A goroutine that never reaches DONE (stuck in WAITING forever)
is a LEAK - see the leak diagram below.
```

---

## 🎯 Concurrency vs. Parallelism

```
CONCURRENCY (structure) - dealing with multiple things at once
  Single core, two goroutines, interleaved:

  Core 1:  [ A ][ B ][ A ][ B ][ A ][ B ]
                time ─────────────────►

  Both A and B make progress, but never run at the EXACT
  same instant. This is concurrency without parallelism.


PARALLELISM (execution) - doing multiple things at once
  Two cores, two goroutines, simultaneous:

  Core 1:  [ A ][ A ][ A ][ A ][ A ][ A ]
  Core 2:  [ B ][ B ][ B ][ B ][ B ][ B ]
                time ─────────────────►

  A and B run at the literal same instant on separate cores.
  This is concurrency THAT HAPPENS TO ALSO be parallel.


  Concurrency is about STRUCTURE (independent tasks, interleavable).
  Parallelism is about EXECUTION (multiple cores, same instant).
  Goroutines always give you the former; the Go runtime + your
  hardware's core count decide how much of the latter you get.
```

---

## 🎯 sync.WaitGroup Flow Diagram

```
                    var wg sync.WaitGroup
                            │
              ┌─────────────┼─────────────┐
              │             │             │
        wg.Add(1)     wg.Add(1)     wg.Add(1)     counter = 3
              │             │             │
              ▼             ▼             ▼
        go worker()   go worker()   go worker()
              │             │             │
       (each runs      (each runs    (each runs
        concurrently,   concurrently,  concurrently,
        eventually      eventually     eventually
        calls)          calls)         calls)
              │             │             │
        defer wg.Done() defer wg.Done() defer wg.Done()
              │             │             │
              ▼             ▼             ▼
         counter=2     counter=1     counter=0
                                          │
                                          ▼
                                   wg.Wait() UNBLOCKS
                                   (was blocking since
                                    the call site, below)

    main goroutine:
        wg.Wait()  ◄──── blocks here until counter reaches 0
        (continues only after ALL three workers called Done())
```

**Rule:** `Add` before you `go`. `Done` (via `defer`) inside the goroutine. `Wait` wherever you need to block until they're all finished.

---

## 🎯 Race Condition Timeline

```
Shared variable: counter = 0

  Goroutine A                    Goroutine B
  -----------                    -----------
  read counter (0)
                                  read counter (0)     ← both read
                                                          BEFORE
  compute 0 + 1 = 1                                      either
                                  compute 0 + 1 = 1       writes!
  write counter = 1
                                  write counter = 1

  Final value: counter == 1  ← WRONG! Should be 2.
  One of the two increments was silently lost because both
  goroutines based their "+1" on the same stale read of 0.

  This is exactly what happens (at scale) with:
      counter++     // read, add 1, write - NOT one atomic step

  across 1000 goroutines - Exercise 5 shows real captured output
  where the final count comes out lower than 1000 every time,
  by a different amount each run.


THE FIX - a mutex serializes the read-modify-write:

  Goroutine A                    Goroutine B
  -----------                    -----------
  mu.Lock()   ◄─ acquires lock
  read counter (0)
  compute 0 + 1 = 1
  write counter = 1
  mu.Unlock() ─► releases lock
                                  mu.Lock()   ◄─ NOW acquires it
                                  read counter (1)
                                  compute 1 + 1 = 2
                                  write counter = 2
                                  mu.Unlock()

  Final value: counter == 2  ← CORRECT. No interleaving possible
  inside the locked section - Exercise 6 verifies this exact fix.
```

---

## 🎯 Goroutine Leak Diagram

```
   main()                              leaky-goroutine
   ------                              ---------------
   ch := make(chan int)
   go leaky(ch)  ─────────────────►    starts running
                                        ch <- val   (send)
                                             │
                                             │  BLOCKS - no
                                             │  receiver ever
   select {                                 │  shows up
   case <-time.After(200ms):                │
       "gave up waiting"                    ▼
   }                                   stuck here FOREVER
   (main moves on and                  (this goroutine's stack,
    eventually returns)                 and anything it was
                                         holding, is never
                                         reclaimed while the
                                         process keeps running)

   In THIS exercise: main() exits shortly after, which kills
   every goroutine (leaked or not) - safe, because it's a
   short-lived demo program.

   In a REAL long-running server: this exact shape (one leaked
   goroutine per request that hits this code path) accumulates
   forever, slowly consuming memory until the process is
   restarted or runs out of resources.
```

---

## 🧩 Deterministic vs. Non-Deterministic Output

```
NON-DETERMINISTIC (this course says so explicitly, every time):
  - Exercise 1: whether the goroutine's print appears at all
  - Exercise 5: the exact "Got:" count from the racy counter

DETERMINISTIC (verified stable across 5+ real runs):
  - Everything synchronized with WaitGroup + Mutex, and where the
    CHECKED result doesn't depend on internal print/finish order
    (Exercises 2, 3, 4, 6, 7, 8, 9, 10)

THE TECHNIQUE that makes concurrent output checkable:
  1. Let goroutines finish in whatever order the scheduler picks
  2. wg.Wait() to know they're ALL done
  3. Only THEN read/sort/print the shared result

  Internal timing: unpredictable.
  Final, checked answer: predictable.
```

---

## 🚨 Common Mistakes

### Mistake 1: Forgetting to Wait

```go
// ❌ WRONG
go fmt.Println("hello")
// main() may exit before this ever runs

// ✅ RIGHT
wg.Add(1)
go func() {
    defer wg.Done()
    fmt.Println("hello")
}()
wg.Wait()
```

### Mistake 2: Assuming Pre-1.22 Loop Variable Sharing

```go
// ✅ CORRECT on Go 1.22+ (this project: go1.26.5) - no workaround needed
for i := 0; i < 5; i++ {
    go func() { fmt.Println(i) }()   // each goroutine gets its own i
}
```

### Mistake 3: Data Races From Unsynchronized Shared State

```go
// ❌ WRONG - counter++ from many goroutines with no lock
counter++

// ✅ RIGHT
mu.Lock()
counter++
mu.Unlock()
```

### Mistake 4: Goroutine Leaks

```go
// ❌ WRONG - blocks forever if nobody reads ch
go func() { ch <- result }()

// ✅ RIGHT - always have an escape hatch
go func() {
    select {
    case ch <- result:
    case <-time.After(timeout):
    }
}()
```

---

## 📈 Progression Summary

### Understanding Level 20

Level 20 teaches how Go runs independent work concurrently:

1. **Goroutines** — lightweight, runtime-scheduled concurrent execution
2. **WaitGroup** — the standard way to know when a batch of goroutines is done
3. **Loop variable capture** — verified safe on Go 1.22+ (this project's go1.26.5)
4. **Data races** — what they are, how to see one, how `-race` finds them
5. **Mutex** — the minimal fix for a race (full coverage: Level 23)
6. **Leaks** — goroutines blocked forever are never garbage collected
7. **Judgment** — knowing when goroutines are worth their complexity

### Prerequisites for Level 21

Before moving to Level 21 (Channels), you need:

- ✅ Comfortable launching goroutines with `go`
- ✅ Comfortable with `sync.WaitGroup`'s Add/Done/Wait
- ✅ Understand why loop variable capture is safe on this project's Go version
- ✅ Can explain what a data race is and how to fix one with `sync.Mutex`
- ✅ Understand what a goroutine leak is and why it matters

### Ready for Level 21?

Level 21 replaces the mutex-protected-slice/map patterns you just practiced with Go's native communication primitive:
- Sending and receiving values through channels
- Buffered vs. unbuffered channels
- `select` for waiting on multiple channels at once
- A more idiomatic way to coordinate goroutines than shared memory + locks

---

## ✅ Checklist Before Level 21

- [ ] Can launch a goroutine and explain why its output isn't guaranteed without synchronization
- [ ] Can use `sync.WaitGroup` to wait for a known number of goroutines
- [ ] Can explain why this project's Go version (1.26.5) makes per-iteration loop variable capture safe
- [ ] Can reproduce a data race and read (or reason about, if `-race` is unavailable) its output
- [ ] Can fix a race with a minimal `sync.Mutex`
- [ ] Can explain a goroutine leak and why it's never garbage collected
- [ ] Can explain when NOT to reach for a goroutine
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The Keyword
`go someFunc()` launches a goroutine - that's the entire syntax. Everything else in this level is about coordinating what happens next.

### The Synchronization
`sync.WaitGroup` answers "are they all done yet?" `sync.Mutex` answers "can I safely touch this shared thing right now?" Neither is optional once goroutines touch shared state.

### The Honesty
Concurrent programs have real non-determinism in *when* things happen. Good concurrent code design makes the *checked, final* result deterministic anyway - that's the skill this level is really teaching.

### The Boundary
Goroutines shine for I/O-bound and genuinely independent work. They cost real overhead and complexity - don't reach for one to add `a + b`.

---

## 📚 Next Level

Level 21: Channels
- Sending and receiving values between goroutines
- Buffered vs. unbuffered channels
- `select` for waiting on multiple channels

You've written real concurrent Go programs - keep going! 🚀
