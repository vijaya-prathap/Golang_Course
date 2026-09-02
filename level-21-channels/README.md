# Level 21: Channels - Complete Guide

## Introduction

Welcome to Level 21! You've mastered goroutines (Level 20) - you can launch concurrent work with the `go` keyword, wait for it to finish with `sync.WaitGroup`, and you got a brief preview of `sync.Mutex` for protecting shared state. But goroutines on their own have a problem: how does one goroutine hand a *value* to another safely, without both of them poking at the same shared variable at the same time?

Go's answer is captured in one of its most famous mottos:

> **Don't communicate by sharing memory; share memory by communicating.**

Instead of two goroutines fighting over a shared variable (guarded by a mutex), you pass the value itself through a **channel** - a typed conduit built into the language. Channels are how goroutines talk to each other, synchronize their timing, and hand off ownership of data cleanly. This level is marked **Advanced** for a reason: channels are simple to describe but easy to misuse, and misusing them produces some of Go's most dramatic runtime crashes - which is exactly why we'll trigger a few of them on purpose.

---

## Table of Contents

1. [What a Channel Is](#what-a-channel-is)
2. [Creating and Using Unbuffered Channels](#creating-and-using-unbuffered-channels)
3. [Buffered Channels](#buffered-channels)
4. [Closing Channels](#closing-channels)
5. [Ranging Over a Channel](#ranging-over-a-channel)
6. [Channel Direction Types](#channel-direction-types)
7. [Deadlocks](#deadlocks)
8. [The Worker Pool Pattern](#the-worker-pool-pattern)
9. [nil Channels](#nil-channels)
10. [Best Practices](#best-practices)
11. [Common Mistakes](#common-mistakes)

---

## What a Channel Is

A channel is a **typed conduit** for sending and receiving values between goroutines. Think of it as a pipe: one goroutine puts a value in one end, another goroutine takes it out the other end. The type in `chan int` or `chan string` says exactly what kind of value can travel through it - the compiler enforces this, just like any other type in Go.

```go
ch := make(chan int)     // a channel that carries int values
ch2 := make(chan string) // a channel that carries string values
```

Channels exist to solve the coordination problem that Level 20 left open. A `sync.Mutex` protects a shared variable so two goroutines don't corrupt it - but you still have to remember to lock and unlock correctly every time you touch that variable. A channel instead moves the *value* from one goroutine to another. Once you've sent it, you no longer touch it from the sending side - ownership has transferred. There's nothing shared left to protect.

That's the philosophy behind "share memory by communicating": rather than two goroutines both holding a reference to the same data and taking turns locking it, one goroutine owns the data, does something with it, and then hands it to the next goroutine over a channel. At any instant, exactly one goroutine has it.

---

## Creating and Using Unbuffered Channels

Create a channel with `make`. With no second argument, you get an **unbuffered** channel:

```go
ch := make(chan int) // unbuffered channel of int
```

Three operations work on a channel:

```go
ch <- v      // send: put value v onto the channel
v := <-ch    // receive: take a value off the channel, assign it to v
close(ch)    // close: signal that no more values will be sent (see below)
```

The defining trait of an **unbuffered** channel: **a send blocks until a receiver is ready to take the value, and a receive blocks until a sender has a value to give.** Neither side can proceed alone - they must "rendezvous" at the same instant. This is what makes unbuffered channels a synchronization tool, not just a data-passing one.

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch := make(chan string)

    go func() {
        time.Sleep(100 * time.Millisecond) // simulate the worker doing some work
        ch <- "hello from the worker goroutine"
    }()

    fmt.Println("main: blocked, waiting to receive from the unbuffered channel...")
    msg := <-ch
    fmt.Println("main: received:", msg)
}
```

```
main: blocked, waiting to receive from the unbuffered channel...
main: received: hello from the worker goroutine
```

Notice `main` genuinely **waits** for the goroutine - even though the goroutine sleeps for 100ms first, `main`'s `<-ch` doesn't proceed until the send happens. That blocking is the synchronization: the channel guarantees "the worker got this far" before `main` can continue. If the goroutine crashed before ever sending, `main` would block on `<-ch` forever (we'll deliberately trigger this failure mode in the [Deadlocks](#deadlocks) section).

---

## Buffered Channels

Add a capacity as the second argument to `make` to get a **buffered** channel:

```go
ch := make(chan int, 3) // buffered channel, capacity 3
```

A buffered channel behaves differently: **a send only blocks once the buffer is full.** Up to capacity, values can pile up in the channel without any goroutine standing by to receive them yet.

```go
buffered := make(chan int, 3)

buffered <- 1 // does not block - buffer has room
buffered <- 2 // does not block - buffer has room
buffered <- 3 // does not block - buffer is now full

fmt.Println(len(buffered), cap(buffered)) // 3 3
```

A **fourth** send on that same channel, with nothing draining it, would block - the buffer is full and there's no room until something receives. This is the crucial contrast with unbuffered channels:

| | Unbuffered (`make(chan int)`) | Buffered (`make(chan int, n)`) |
|---|---|---|
| Send blocks until... | a receiver is ready | the buffer has room (i.e. fewer than `n` items waiting) |
| Guarantees | sender and receiver rendezvous at the same instant | only that the buffer isn't over capacity - no timing guarantee |
| Good for | synchronization, "did the other side actually get this?" | decoupling producer/consumer speed, avoiding sender blocking on every single item |

`len(ch)` tells you how many values are currently buffered; `cap(ch)` tells you the buffer's capacity (`0` for an unbuffered channel).

---

## Closing Channels

`close(ch)` marks a channel as finished: **no more values will ever be sent on it.** It's a signal, not a value - you close a channel to tell receivers "that's everything."

Two rules follow directly from that:

1. **Receiving from a closed channel never blocks.** If there are values still buffered, you get them first, in order. Once it's drained, every further receive returns immediately with the **zero value** of the channel's type.
2. **You can tell a real zero value apart from "the channel is closed" using the comma-ok form:**

```go
v, ok := <-ch
// ok == true  -> v is a real value that was sent
// ok == false -> the channel is closed AND drained; v is the zero value
```

```go
ch := make(chan int, 2)
ch <- 10
ch <- 20
close(ch)

v1, ok1 := <-ch // 10, true  (a real value, sent before close)
v2, ok2 := <-ch // 20, true  (still a real value)
v3, ok3 := <-ch // 0,  false (drained AND closed)
```

The single-value form `v := <-ch` still works on a closed channel - it just silently gives you the zero value once drained, with no way to tell that apart from a legitimately-sent zero. That's exactly why the comma-ok form exists.

**Sending on a closed channel panics - it does not block, and it does not silently fail.** This is the language telling you loudly that closing means "done," permanently:

```go
ch := make(chan int)
close(ch)
ch <- 1 // panic: send on closed channel
```

We'll capture that exact panic text in the exercises.

---

## Ranging Over a Channel

`for v := range ch` receives values from a channel repeatedly and **automatically stops when the channel is closed and drained** - no manual comma-ok check needed:

```go
func produce(ch chan<- int, n int) {
    for i := 1; i <= n; i++ {
        ch <- i
    }
    close(ch) // this is what lets the range loop below terminate
}

func main() {
    ch := make(chan int)
    go produce(ch, 5)

    for v := range ch {
        fmt.Println("received:", v)
    }
    fmt.Println("channel closed, range loop exited")
}
```

```
received: 1
received: 2
received: 3
received: 4
received: 5
channel closed, range loop exited
```

With a single producer sending in order and a single consumer ranging, the values arrive in the exact order they were sent - a channel is FIFO for one sender and one receiver. If the producer never calls `close`, the `range` loop never sees "no more values" and blocks forever waiting for the next one - a real and common bug (see [Common Mistakes](#common-mistakes)).

---

## Channel Direction Types

A function parameter can restrict a channel to **send-only** or **receive-only**, using arrow syntax in the type itself:

```go
chan<- int  // send-only: you may only do ch <- v
<-chan int  // receive-only: you may only do v := <-ch
```

```go
func sendOnly(ch chan<- int) {
    ch <- 100
    // v := <-ch would NOT compile here - ch is send-only
}

func receiveOnly(ch <-chan int) int {
    return <-ch
    // ch <- 5 would NOT compile here - ch is receive-only
}
```

A plain bidirectional `chan int` passed into either function is automatically restricted to that direction for the duration of the call - Go converts it for you. This is good API design: the function signature documents, and the **compiler enforces**, exactly how the function is allowed to use the channel. A caller reading `func produce(out chan<- int)` knows immediately that `produce` only ever sends - it can't accidentally receive and drain the very channel it's supposed to be filling.

Trying to violate the direction is a genuine **compile error**, not a runtime issue:

```go
func sendOnly(ch chan<- int) {
    v := <-ch // trying to receive from a send-only channel
    _ = v
}
```

```
./main.go:4:12: invalid operation: cannot receive from send-only channel chan<- int ch (variable of type chan<- int)
```

That's caught before the program ever runs - one of the strongest arguments for using direction types on any function that takes a channel.

---

## Deadlocks

Because an unbuffered send blocks until a receiver is ready (and vice versa), it's easy to write a program where that "someone will receive" promise is never kept. When **every** goroutine in the program is blocked and none of them can ever be woken up, Go's runtime detects it and crashes the whole program immediately with a `fatal error` - it does not hang forever:

```go
package main

import "fmt"

func main() {
    fmt.Println("about to send on an unbuffered channel with no receiver...")
    ch := make(chan int)
    ch <- 42 // nobody will ever receive this - deadlock
    fmt.Println("this line never runs")
}
```

Running this produces (captured verbatim from a real run):

```
about to send on an unbuffered channel with no receiver...
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan send]:
main.main()
        /path/to/main.go:8 +0x68
exit status 2
```

This is a genuinely useful safety net: a deadlock in Go is loud and immediate, not a silent hang. The exact file path and offset (`+0x68`) will differ on your machine, but the `fatal error: all goroutines are asleep - deadlock!` line, the `[chan send]` goroutine state, and `exit status 2` are exactly what you'll see. The same failure happens in reverse - receiving from an unbuffered channel that nothing will ever send on deadlocks identically, just reported as `[chan receive]`.

---

## The Worker Pool Pattern

One of the most common uses for channels: a fixed number of **worker** goroutines pull jobs off a shared `jobs` channel and push results onto a shared `results` channel. The number of workers is bounded (unlike spawning one goroutine per job), and the channels do all the coordination - no manual locking required.

```go
func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()
    for j := range jobs { // keeps pulling jobs until jobs is closed and drained
        results <- j * j
    }
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
    close(jobs) // no more jobs - lets each worker's `range jobs` loop end

    wg.Wait()     // wait for all workers to finish draining jobs
    close(results) // safe to close now - nothing sends to results anymore
}
```

The shape to notice: `jobs` is closed by the goroutine that's *sending* jobs (`main`), which lets every worker's `for j := range jobs` loop end on its own once the jobs are drained. `results` is closed only **after** `wg.Wait()` confirms every worker is done sending to it - closing it any earlier would risk a "send on closed channel" panic from a worker still in flight. This ordering - sender closes, and only after every other sender is guaranteed finished - is the core discipline of the worker pool pattern (and of channel closing in general; see [Best Practices](#best-practices)).

The exact job each worker picks up, and the order workers finish, is **not deterministic** - three workers are racing to pull from the same `jobs` channel. What *is* deterministic is the final **set** of results, since every job is processed exactly once. Sorting the collected results before printing (or summing them) turns a nondeterministic race into a checkable, deterministic outcome - exactly what the exercises do.

---

## nil Channels

A channel variable that was never assigned with `make` is `nil` - its zero value. **Both sending and receiving on a nil channel block forever**, with no exception and no deadlock detection escape hatch of their own (a program that only ever touches a nil channel, with nothing else running, is a deadlock like any other):

```go
var ch chan int // ch is nil
<-ch            // blocks forever
```

That sounds useless, but it's actually a deliberate tool - a `nil` channel is a way to **disable a case in a `select` statement** (Level 22's full subject) at runtime. Since a `select` never picks a case that would block forever, setting one of its channels to `nil` effectively removes that case from consideration until you reassign it to a real channel. You'll use this trick to build things like "stop reading from this channel after the first value" once `select` is in your toolkit. For now, just remember: a nil channel isn't broken, it's a channel that will never be ready - and `select` is what makes that useful instead of dangerous.

---

## Best Practices

### 1. The Sender Should Usually Close the Channel

Only the goroutine that sends values should call `close` on that channel. A receiver has no way to know whether another sender is still coming - if a receiver closes a channel a sender is still using, that sender's next send panics.

```go
// ✅ Good - the producer (sender) closes when it's done
func produce(ch chan<- int) {
    defer close(ch)
    for i := 0; i < 5; i++ {
        ch <- i
    }
}
```

### 2. Never Close a Channel From the Receiving Side

```go
// ❌ WRONG - the receiver doesn't know if more sends are coming
func consume(ch chan int) {
    v := <-ch
    fmt.Println(v)
    close(ch) // dangerous - some other goroutine might still be sending
}
```

### 3. Never Close the Same Channel Twice

Closing an already-closed channel panics, just like sending on one:

```go
ch := make(chan int)
close(ch)
close(ch) // panic: close of closed channel
```

```
panic: close of closed channel

goroutine 1 [running]:
main.main()
        /path/to/main.go:9 +0x6c
exit status 2
```

If multiple goroutines might all try to close the same channel, coordinate with a `sync.Once` or make sure only one designated goroutine ever calls `close`.

### 4. Prefer Channels for Coordination, a Mutex for Pure Shared-State Protection

Channels shine when you're **passing ownership of a value** or **signaling that something happened** (a job is ready, a stage is done). If you genuinely just need to protect a shared counter or map that many goroutines read and write in place - with nothing to "hand off" - a `sync.Mutex` (previewed in Level 20) is often simpler and cheaper. Don't force every synchronization problem through a channel just because channels feel more idiomatic; use the tool that matches the shape of the problem.

### 5. Design Function Signatures With Directional Channel Types

If a function only ever sends, declare its parameter `chan<-`. If it only ever receives, declare it `<-chan`. This documents intent and lets the compiler catch misuse - see [Channel Direction Types](#channel-direction-types).

---

## Common Mistakes

### Mistake 1: Deadlocking on an Unbuffered Channel

```go
// ❌ WRONG - nothing will ever receive this; the whole program crashes
func main() {
    ch := make(chan int)
    ch <- 1
}
```

```
fatal error: all goroutines are asleep - deadlock!
```

```go
// ✅ RIGHT - a receiver must be running concurrently
func main() {
    ch := make(chan int)
    go func() { fmt.Println(<-ch) }()
    ch <- 1
}
```

### Mistake 2: Sending on a Closed Channel

```go
// ❌ WRONG - panics immediately
ch := make(chan int)
close(ch)
ch <- 1 // panic: send on closed channel
```

```go
// ✅ RIGHT - only the sender closes, and only after it's truly done sending
ch := make(chan int)
go func() {
    defer close(ch)
    ch <- 1
}()
```

### Mistake 3: Closing a Channel Twice

```go
// ❌ WRONG - panics on the second close
ch := make(chan int)
close(ch)
close(ch) // panic: close of closed channel
```

```go
// ✅ RIGHT - close exactly once, from exactly one place
var once sync.Once
closeCh := func() { once.Do(func() { close(ch) }) }
```

### Mistake 4: Forgetting to Close a Channel a range Loop Is Waiting On

```go
// ❌ WRONG - producer never calls close; the range loop below blocks forever
func produce(ch chan<- int) {
    for i := 0; i < 5; i++ {
        ch <- i
    }
    // missing: close(ch)
}

func main() {
    ch := make(chan int)
    go produce(ch)
    for v := range ch { // never sees "closed" - hangs after the 5th value
        fmt.Println(v)
    }
}
```

```go
// ✅ RIGHT - always close when the producer is done sending
func produce(ch chan<- int) {
    defer close(ch)
    for i := 0; i < 5; i++ {
        ch <- i
    }
}
```

### Mistake 5: Closing a Channel From the Wrong Side

```go
// ❌ WRONG - the receiver closes a channel it doesn't own
func consume(ch chan int) {
    for v := range ch {
        fmt.Println(v)
    }
    close(ch) // too late to matter here, but dangerous as a habit -
              // a sender elsewhere may still be trying to send
}
```

### Mistake 6: Assuming a nil Channel "Just Works Like an Empty One"

```go
// ❌ WRONG ASSUMPTION - a nil channel doesn't return a zero value like a
// closed one does; it blocks FOREVER
var ch chan int
<-ch // blocks forever - this is not the same as a closed channel
```

---

## Summary

**Channel Basics:**
- `make(chan T)` - unbuffered: send and receive block until both sides are ready
- `make(chan T, n)` - buffered: send only blocks once `n` items are already waiting
- `ch <- v` sends, `v := <-ch` receives, `close(ch)` signals "no more sends"

**Closing:**
- Receiving from a closed, drained channel returns the zero value immediately
- `v, ok := <-ch` - `ok` is `false` only when the channel is closed AND drained
- Sending on a closed channel panics; closing a closed channel panics

**range:**
- `for v := range ch` receives until the channel is closed, then stops automatically
- Preserves send order for a single sender

**Direction Types:**
- `chan<- T` (send-only), `<-chan T` (receive-only) as parameter types
- The compiler enforces correct usage - misuse is a compile error, not a bug you find at runtime

**Failure Modes (all real, all captured in the exercises):**
- Deadlock: `fatal error: all goroutines are asleep - deadlock!` (program crashes immediately, doesn't hang)
- `panic: send on closed channel`
- `panic: close of closed channel`

**Patterns:**
- Worker pool: fixed workers pull from a shared `jobs` channel, push to a shared `results` channel
- Pipeline: chain stages together, each stage's output channel is the next stage's input

---

## Next Steps

You now understand:
- ✅ Unbuffered vs. buffered channel blocking behavior
- ✅ Closing channels, the comma-ok zero-value distinction, and ranging until closed
- ✅ Channel direction types and why the compiler enforcing them matters
- ✅ What a real deadlock and a real "send/close on closed channel" panic look like
- ✅ The worker pool pattern and multi-stage pipelines

**Next level:** Level 22 - Select
- Waiting on multiple channels at once with `select`
- `default` cases for non-blocking channel operations
- Timeouts with `time.After`
- Using `nil` channels to disable a `select` case at runtime

You can now make goroutines talk to each other safely and predictably. `select` is what lets you listen to several conversations at once. Keep going! 🚀
