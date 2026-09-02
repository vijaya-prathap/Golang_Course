# Level 21: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: Channel Fundamentals
```
Day 1:  What a channel is, unbuffered channels, blocking send/receive
Day 2:  Buffered channels and the blocking-behavior contrast
Day 3:  Closing channels, the comma-ok zero-value distinction
Day 4:  Ranging over a channel until closed
Day 5:  Channel direction types and the compile-time guarantee
Day 6:  Deadlocks - causing and reading the real fatal error
Day 7:  Panics - send on closed, close of closed channel
```

### Week 2: Patterns & Practice
```
Day 1:  Exercises 1-3 (hand-off, buffered contrast, comma-ok)
Day 2:  Exercises 4-5 (range until closed, direction types)
Day 3:  Exercises 6-7 (deadlock, panics - captured verbatim)
Day 4:  Exercise 8 (worker pool)
Day 5:  Exercises 9-10 (pipelines)
Day 6:  Bonus challenges (fan-in, rate limiter, pub/sub)
Day 7:  Review & consolidation
```

---

## 🔌 Channel-as-Pipe: Unbuffered vs. Buffered

```
UNBUFFERED  make(chan int)            capacity 0

   sender                 receiver
     |                        |
     |---- ch <- v -----X-----|      X = handshake point
     |                        |
   BLOCKS until a receiver  BLOCKS until a sender
   is ready to take v RIGHT  has a value RIGHT NOW
   NOW

   Think: a phone call. Both parties must be present
   at the same instant for anything to happen.


BUFFERED    make(chan int, 3)         capacity 3

   sender            [ □ □ □ ]            receiver
     |                 buffer                |
     |--- ch<-1 -->[ 1 ]                     |
     |--- ch<-2 -->[ 1 2 ]                   |
     |--- ch<-3 -->[ 1 2 3 ]  (now FULL)     |
     |--- ch<-4 -->  BLOCKS until buffer     |
     |               has room again          |
     |                                       |<-- <-ch  (takes 1, buffer has room)

   Think: a mailbox. The sender can drop off mail (up to
   the box's size) and walk away - no one has to be
   standing there. Only once the box is full does
   dropping off another letter require waiting.
```

| | Unbuffered | Buffered (capacity `n`) |
|---|---|---|
| `make(...)` | `make(chan int)` | `make(chan int, n)` |
| Send blocks until | a receiver is ready | fewer than `n` items are waiting |
| Synchronization guarantee | sender and receiver rendezvous at the same instant | none beyond "buffer isn't over capacity" |
| `cap(ch)` | `0` | `n` |

---

## ✅❌ The comma-ok Truth Table

```go
v, ok := <-ch
```

| Channel state | `v` | `ok` |
|---|---|---|
| Open, a real value was sent | that real value | `true` |
| Open, but nothing sent yet (blocks until one of these becomes true) | — | — |
| **Closed**, values still buffered | the next buffered value (still real!) | `true` |
| **Closed AND fully drained** | the zero value of the channel's type | `false` |

```
ch := make(chan int, 2)
ch <- 10
ch <- 20
close(ch)

<-ch  ->  10, true    (a real, sent value - close doesn't erase what's already buffered)
<-ch  ->  20, true    (still real)
<-ch  ->   0, false   (drained AND closed - this is the "no more values, ever" signal)
<-ch  ->   0, false   (stays this way forever)
```

**Rule of thumb:** `ok == false` is the ONLY reliable way to know "this channel is done." A bare `v := <-ch` can't tell a real zero-valued send apart from a closed channel - always use comma-ok (or `range`) when that distinction matters.

---

## ➡️⬅️ Channel Direction Types

```
   chan int         bidirectional - can send AND receive
   chan<- int       send-only     - can ONLY do ch <- v
   <-chan int       receive-only  - can ONLY do v := <-ch

   The arrow points the way the data flows relative to
   the channel:

       chan<- int      data flows INTO the channel  (send)
       <-chan int      data flows OUT of the channel (receive)
```

```
func produce(out chan<- int)     // caller knows: this function only ever SENDS
func consume(in <-chan int)      // caller knows: this function only ever RECEIVES

A plain chan int passed into either parameter is automatically
narrowed to that direction for the call - no cast needed.

Violating it is a COMPILE ERROR, not a runtime bug:

  func sendOnly(ch chan<- int) {
      v := <-ch   // ❌ won't compile
  }

  ./main.go:4:12: invalid operation: cannot receive from send-only
  channel chan<- int ch (variable of type chan<- int)
```

---

## 🏭 Worker Pool Architecture

```
                     ┌─────────┐
   jobs channel ────>│worker 1 │──┐
   (buffered,        └─────────┘  │
    numJobs slots)   ┌─────────┐  │      results channel
              ┌──────>│worker 2 │──┼────> (buffered,
              │      └─────────┘  │       numJobs slots)
   main sends │      ┌─────────┐  │
   1..9, then │      │worker 3 │──┘
   close(jobs)└──────└─────────┘
                       (3 fixed workers, racing to
                        pull from the SAME jobs channel)

  1. main fills `jobs`, then close(jobs)
  2. each worker: for j := range jobs { results <- j*j }
     - range exits on its own once jobs is closed+drained
  3. main: wg.Wait() blocks until ALL workers call wg.Done()
  4. ONLY THEN: close(results) - safe, because no worker
     can possibly still be sending
  5. main drains results, sorts, sums

  Non-deterministic: which worker gets which job, and the
  order results arrive in.
  Deterministic: the final SET of results (every job run
  exactly once) - sort before printing to get one exact,
  checkable answer every time.
```

---

## 💥 Deadlock Scenario

```
   main goroutine                    (no other goroutine exists)

   ch := make(chan int)
   ch <- 42   ─────────────>  BLOCKED - waiting for a receiver
                                that will NEVER come, because
                                there is no other goroutine
                                left to run

   Go's runtime notices: EVERY goroutine in the program is
   asleep, and none of them can ever be woken by anything
   else in the program (a timer, another goroutine, I/O).
   That's an unrecoverable deadlock, so the runtime kills
   the whole program immediately:

   fatal error: all goroutines are asleep - deadlock!

   goroutine 1 [chan send]:
   main.main()
           /path/to/main.go:8 +0x68
   exit status 2

   This is NOT a hang - the crash is immediate (milliseconds).
   The same thing happens in reverse (receive with no possible
   sender), reported as [chan receive] instead of [chan send].
```

---

## 🚨 Common Mistakes

| Mistake | Wrong | Right |
|---|---|---|
| Deadlock on unbuffered send | `ch := make(chan int); ch <- 1` (nothing ever receives) | launch a receiving goroutine first, or use a buffered channel sized for the load |
| Send on closed channel | `close(ch); ch <- 1` → `panic: send on closed channel` | only the sender closes, and only once it's truly done sending |
| Close a channel twice | `close(ch); close(ch)` → `panic: close of closed channel` | close exactly once, guarded by `sync.Once` if multiple goroutines might try |
| Forgetting to close before a `range` | producer never calls `close`; consumer's `for v := range ch` blocks forever | always `defer close(ch)` in the producer once it's done sending |
| Closing from the receiving side | receiver calls `close(ch)` after its own loop ends | never close a channel you don't own the sending side of |
| Assuming `nil` chan behaves like closed | `var ch chan int; <-ch` | a `nil` channel blocks **forever** on send or receive - it is not the same as closed |

---

## 📈 Progression Summary

### Understanding Level 21

Level 21 gives goroutines (Level 20) a safe way to talk to each other:

1. **Channels** - typed conduits, "share memory by communicating"
2. **Unbuffered vs. buffered** - two different blocking contracts
3. **Closing** - the "no more values, ever" signal, and its comma-ok truth table
4. **range** - the idiomatic consumer loop, stops itself on close
5. **Direction types** - compiler-enforced API contracts
6. **Failure modes** - deadlocks and panics, captured and understood, not feared
7. **Patterns** - worker pools and multi-stage pipelines

### Prerequisites for Level 22

Before moving to Level 22 (Select), you need:

- ✅ Comfortable creating and using both unbuffered and buffered channels
- ✅ Can explain the comma-ok zero-value distinction from memory
- ✅ Can range over a channel and know why the producer must close it
- ✅ Understand send-only/receive-only channel parameter types
- ✅ Have seen (and can explain) a real deadlock and both "closed channel" panics
- ✅ Can build a worker pool and a simple channel pipeline

### Ready for Level 22?

Level 22 teaches you to wait on *several* channels at once:
- `select` - pick whichever channel is ready first
- `default` - a non-blocking channel operation
- `time.After` - implementing a timeout with a channel
- Using `nil` channels (previewed briefly in this level) to disable a `select` case at runtime

---

## ✅ Checklist Before Level 22

- [ ] Can write unbuffered and buffered channel code from memory
- [ ] Can explain exactly when a buffered send blocks vs. an unbuffered send
- [ ] Can use `v, ok := <-ch` and explain every row of the truth table above
- [ ] Can write a producer that closes its channel and a consumer that ranges over it
- [ ] Can write a function with a `chan<-` or `<-chan` parameter and explain why
- [ ] Have personally triggered and read a real deadlock and both closed-channel panics
- [ ] Can build a worker pool with a fixed number of workers
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The Motto
"Don't communicate by sharing memory; share memory by communicating" - channels move ownership of a value, not just its bits.

### The Contract
Unbuffered channels synchronize (rendezvous); buffered channels decouple (up to capacity) - know which guarantee you actually need.

### The Signal
`close()` means "no more sends, ever" - comma-ok and `range` are how a receiver learns that reliably.

### The Safety Net
Deadlocks and closed-channel misuse crash loudly and immediately in Go - that's a feature, not a flaw. A silent hang would be far worse.

### The Patterns
Worker pools bound concurrency; pipelines chain stages together - both are built from nothing but channels, goroutines, and (for closing/waiting) `sync.WaitGroup`.

---

## 📚 Next Level

Level 22: Select
- Waiting on multiple channels at once with `select`
- `default` cases for non-blocking channel operations
- Timeouts with `time.After`
- Using `nil` channels to disable a `select` case at runtime

You've got channels down! Keep going! 🚀
