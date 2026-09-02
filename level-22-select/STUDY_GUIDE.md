# Level 22: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: select Fundamentals
```
Day 1:  What select does; basic multi-channel waiting
Day 2:  Random tie-breaking, verified with real tallies
Day 3:  The default case and non-blocking operations
Day 4:  Timeouts with time.After
Day 5:  The done-channel cancellation pattern
Day 6:  select with send cases; select on a closed channel
Day 7:  The time.After leak gotcha and time.NewTimer
```

### Week 2: Practice & Application
```
Day 1:  Exercises 1-3 (basic select, tie-breaking, default)
Day 2:  Exercises 4-6 (timeouts, done pattern, send cases)
Day 3:  Exercises 7-8 (closed channel, timer leak fix)
Day 4:  Exercises 9-10 (nil-channel multiplexer, event dispatcher)
Day 5:  Bonus challenges (ping-pong, fan-in, rate limiter)
Day 6-7: Review & consolidation
```

---

## 🎯 select as a Multiplexer

```
                    ┌─────────────────────┐
        ch1 ───────>│                     │
                     │                     │
        ch2 ───────>│      select         │───> runs whichever
                     │                     │     case is ready
      done ───────>│                     │     FIRST
                     │  (default, if      │
                     │   present, runs    │
                     │   only if NOTHING  │
                     │   else is ready)   │
                     └─────────────────────┘

select is NOT a switch:
  - switch compares a value against cases, top to bottom
  - select waits on channel OPERATIONS (send or receive),
    and does not care about source-code order at all
```

---

## 🎲 Random Tie-Breaking, Illustrated

```
Both chA and chB are ready RIGHT NOW:

    select {
    case <-chA:   ← ready
    case <-chB:   ← ALSO ready
    }

    Go does NOT pick chA just because it's listed first.
    Internally, Go picks uniformly at random among every
    ready case, every single time.

    100,000 trials, both channels pre-filled every time:

        chA: ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  49,986  (50.0%)
        chB: ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  50,014  (50.0%)

    Roughly even, every run - never anywhere close to 100/0.
```

---

## 🚦 Blocking vs. Non-Blocking (default)

| | Without `default` | With `default` |
|---|---|---|
| No case ready | **Blocks** - waits until one becomes ready | Runs `default` **immediately** |
| A case becomes ready later | Proceeds with it whenever it happens | Never revisited - `default` already ran |
| Use when... | You genuinely want to wait | You want to "check and move on" |

```go
// Blocks until ch has something (or forever, if it never does)
select {
case v := <-ch:
    use(v)
}

// Never blocks - checks once, falls through instantly if empty
select {
case v := <-ch:
    use(v)
default:
    // ch had nothing ready right now
}
```

---

## ⏱️ Timeout Pattern

```
   select {
   case result := <-work:      ──────┐
   case <-time.After(d):       ──────┤   whichever
   }                                 │   arrives
                                      │   FIRST wins
        work          ╌╌╌╌╌╌╌╌╌╌╌╌╌╌>│
        (real answer, arrives           <-- if work is
         whenever it's ready)              faster than d

        time.After(d) ───────────────>│
        (fires exactly once,             <-- if work takes
         after duration d)                  longer than d

   Both are just channels to select - "timeout" isn't a
   special case, it's an ordinary race between two channels.
```

---

## 🛑 done-Channel Cancellation Flow

```
   main                              worker goroutine
   ────                              ────────────────
                                      for {
   work <- 1  ─────────────────────>   select {
                                        case w := <-work: process(w)
   work <- 2  ─────────────────────>   case <-done:      return
                                        }
   work <- 3  ─────────────────────>  }
                                      (still looping - done not closed yet)

   close(done) ─────────────────────>  <-done fires on the NEXT select
                                        "received done signal, stopping"
                                        return

   <-finished  <────────────────────  (worker goroutine confirms exit)

   Closing done is a BROADCAST: every goroutine select-ing on it
   notices on their very next pass through the loop, however many
   goroutines are watching. This is the direct precursor to the
   context package (a later level formalizes exactly this pattern).
```

---

## ⏲️ time.After vs. time.NewTimer (in a loop)

| | `time.After(d)` called every iteration | `time.NewTimer(d)` created once, `Stop()`/`Reset()` each iteration |
|---|---|---|
| Timers created over N iterations | N (one per call) | 1 (reused) |
| Each abandoned timer | Lives until it fires, unreclaimed until then | N/A - never abandoned |
| Code needed | `case <-time.After(d):` - trivial | A few extra lines: `Stop()` before `Reset()`, drain if needed |
| Best for | A handful of calls, or a one-shot wait outside any loop | Any loop that runs many times or for a long time |

```go
// ❌ Leaky in a loop - a new timer every pass
for {
    select {
    case v := <-ch:
        use(v)
    case <-time.After(d):
        onTimeout()
    }
}

// ✅ Reused - one timer for the whole loop's life
timer := time.NewTimer(d)
defer timer.Stop()
for {
    if !timer.Stop() {
        select {
        case <-timer.C:
        default:
        }
    }
    timer.Reset(d)
    select {
    case v := <-ch:
        use(v)
    case <-timer.C:
        onTimeout()
    }
}
```

---

## 🕳️ nil Channels: Disabling a Case

```
   var ch chan int   // ch is nil

   select {
   case v := <-ch:     // this case is NEVER ready - nil channels
       use(v)           // never send, never receive, ever
   case <-other:
       ...
   }

   Useful for turning OFF a case at runtime without restructuring
   the select: once a source is exhausted, set it to nil instead
   of leaving its (now-closed) case active:

       case v, ok := <-a:
           if !ok {
               a = nil   // this case is permanently skipped from now on
               continue
           }
```

---

## 🚨 Common Mistakes

| Mistake | Wrong | Right |
|---|---|---|
| Forgetting `default` blocks | `select { case v := <-ch: ... }` with nothing ready | Add `default:` if you need a non-blocking check |
| Assuming case order = priority | Believing the first-listed case wins ties | Ties are broken **randomly** - order has no effect |
| Leaking timers in a loop | `case <-time.After(d):` called every iteration | `time.NewTimer(d)` once, `Stop()`/`Reset()` per iteration |
| No way to stop the loop | `for { select { case v := <-ch: ... } }` forever | Always include a `done` case, counter, or timeout |
| Busy-looping on a closed channel | Reading a closed channel's case without checking `ok` | Check `ok`, and `return`/`break`/set the channel to `nil` |

---

## 📈 Progression Summary

### Understanding Level 22

Level 22 gives channels (Level 21) a way to be waited on together:

1. **select** - wait on multiple channel ops, proceed with whichever is ready
2. **Random tie-breaking** - fairness by design, verified empirically
3. **default** - the non-blocking escape hatch
4. **time.After** - timeouts as an ordinary select case
5. **done channels** - the hand-rolled ancestor of `context`
6. **Send cases, closed channels, nil channels** - the full vocabulary of what a case can be
7. **select + for** - the shape of every real multiplexer and dispatcher

### Prerequisites for Level 23

Before moving to Level 23 (Mutex, WaitGroup & Atomic), you need:

- ✅ Comfortable writing a `select` with multiple receive and send cases
- ✅ Can explain why select's tie-breaking is random, not ordered
- ✅ Can use `default` to make a channel operation non-blocking
- ✅ Can implement a timeout with `time.After` and explain the timer-leak gotcha
- ✅ Can build a done-channel cancellation pattern from scratch
- ✅ Can use a `nil` channel to disable a select case at runtime
- ✅ Can combine select with a for loop into a multiplexer that reliably terminates

### Ready for Level 23?

Level 23 moves from channel-based coordination to direct shared-state protection:
- `sync.Mutex` / `sync.RWMutex` for guarding a variable in place
- `sync.WaitGroup` in depth
- `sync/atomic` for lock-free counters and flags
- When to reach for a mutex or atomic instead of a channel

---

## ✅ Checklist Before Level 23

- [ ] Can write a select with 2+ cases from memory
- [ ] Have personally run the random tie-breaking measurement and seen a roughly even split
- [ ] Can explain exactly when `default` runs vs. when a select blocks
- [ ] Can implement a `time.After` timeout and explain the timer-leak gotcha
- [ ] Can build a done-channel cancellation pattern and explain how it foreshadows `context`
- [ ] Can write a select case that sends, not just receives
- [ ] Understand why a closed channel's case is always ready, and the busy-loop risk
- [ ] Can use a nil channel to disable a select case at runtime
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The Wait
`select` waits on multiple channel operations and proceeds with whichever is ready first - it is not a switch, and case order carries no meaning.

### The Fairness
When more than one case is ready, Go picks randomly - verified, not assumed.

### The Escape Hatch
`default` turns any select non-blocking; a `done` case turns any loop cancellable; a timeout case turns any wait bounded.

### The Gotcha
`time.After` inside a loop leaks a timer per iteration - `time.NewTimer` with `Stop()`/`Reset()` is the fix.

### The Toggle
A `nil` channel's case is never ready - the clean way to disable part of a select at runtime.

---

## 📚 Next Level

Level 23: Mutex, WaitGroup & Atomic
- sync.Mutex and sync.RWMutex for protecting shared state directly
- sync.WaitGroup in depth for waiting on many goroutines
- The sync/atomic package for lock-free counters and flags
- Choosing between channels, mutexes, and atomics for a given problem

You've got select down! Keep going! 🚀
