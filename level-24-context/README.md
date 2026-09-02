# Level 24: Context - Complete Guide

## Introduction

Welcome to Level 24! You've mastered `sync.Mutex`, `sync.WaitGroup`, and atomic operations (Level 23) - you can protect shared state and wait for goroutines to finish. Level 22 (Select) showed you the "done channel" pattern: a goroutine that watches a `chan struct{}` alongside its real work, and stops the moment that channel closes. That pattern is powerful, but every team that used it ended up reinventing it slightly differently - one done channel per operation, no standard way to attach a deadline, no standard way to say *why* something stopped.

`context.Context` is Go's answer: it takes the done-channel idea and formalizes it into a standard type that carries **cancellation signals**, **deadlines**, and (sparingly) **request-scoped values** through an entire call graph - across function calls and across goroutines. This level is marked **Advanced** because context threads through nearly every real Go program you'll write from here on: HTTP handlers (Level 27), database calls (Level 29), and any goroutine that needs to know "should I stop now?"

---

## Table of Contents

1. [What context.Context Is For](#what-contextcontext-is-for)
2. [context.Background() and context.TODO()](#contextbackground-and-contexttodo)
3. [context.WithCancel: Manual Cancellation](#contextwithcancel-manual-cancellation)
4. [context.WithTimeout and context.WithDeadline](#contextwithtimeout-and-contextwithdeadline)
5. [ctx.Err(): Canceled vs. DeadlineExceeded](#ctxerr-canceled-vs-deadlineexceeded)
6. [ctx.Done() and the select Pattern](#ctxdone-and-the-select-pattern)
7. [context.WithValue: Request-Scoped Data](#contextwithvalue-request-scoped-data)
8. [Always Call the Returned cancel Function](#always-call-the-returned-cancel-function)
9. [Propagating Context Through a Call Chain and Goroutines](#propagating-context-through-a-call-chain-and-goroutines)
10. [A Simulated HTTP-Handler-Like Pattern](#a-simulated-http-handler-like-pattern)
11. [Best Practices](#best-practices)
12. [Common Mistakes](#common-mistakes)

---

## What context.Context Is For

`context.Context` is Go's standard way to propagate three things through a call graph - a chain of function calls and the goroutines they launch:

1. **Cancellation signals** - "stop what you're doing, right now"
2. **Deadlines** - "stop what you're doing if you're not done by this time"
3. **Request-scoped values** (sparingly) - data that belongs to one request/operation, like a trace ID

```go
type Context interface {
    Done() <-chan struct{}
    Err() error
    Deadline() (deadline time.Time, ok bool)
    Value(key any) any
}
```

Level 22's `select { case <-done: ...; case <-ch: ... }` pattern already showed you the shape of this: a channel that closes to mean "stop." `context.Context` standardizes that channel (`Done()`), adds a reason you can inspect after the fact (`Err()`), adds an optional deadline (`Deadline()`), and adds a narrow escape hatch for request-scoped data (`Value()`). Because it's a single interface, every function that takes a `context.Context` - whether it's yours, the standard library's, or a third-party package's - can be canceled and timed out the exact same way.

The core idea: a context is created once (usually at the top of a request or operation), and passed down through every function and goroutine that does work on behalf of that operation. If the top-level context is canceled or times out, every downstream piece of work that's paying attention finds out - without any of them needing to know about each other.

---

## context.Background() and context.TODO()

Every context tree starts from a **root context** with no parent. There are exactly two ways to create one:

```go
ctx := context.Background()
ctx := context.TODO()
```

Both behave identically at runtime - neither is ever canceled, has a deadline, or carries a value:

```go
bg := context.Background()
todo := context.TODO()

fmt.Println(bg.Err())          // <nil>
fmt.Println(bg.Done())         // <nil>
fmt.Println(bg.Value("x"))     // <nil>
```

Real output, captured verbatim from a run (Go 1.21+ gives them distinct internal types purely to make debugging clearer - `backgroundCtx` and `todoCtx` - but their observable behavior is identical):

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

The difference between them is purely **documentation of intent**, and it matters for code review:

| | `context.Background()` | `context.TODO()` |
|---|---|---|
| Use when | You're at the actual top of a call graph: `main()`, a test, or an incoming request's real entry point | You know a `context.Context` should be threaded through eventually, but you haven't wired it up yet, or it's genuinely unclear which context should be used here |
| Signals to reviewers | "This is deliberately the root - nothing to propagate" | "This is a placeholder - someone should replace this with a real, propagated context" |

`go vet` and static analysis tools can flag lingering `context.TODO()` calls as a to-do list of places that still need proper context plumbing. In finished, production code, `context.TODO()` should be rare - most contexts are either `context.Background()` at a genuine root, or derived from a context passed into your function.

---

## context.WithCancel: Manual Cancellation

`context.WithCancel` derives a **child context** from a parent, along with a `cancel` function you call to cancel it manually:

```go
ctx, cancel := context.WithCancel(context.Background())
```

Calling `cancel()` closes `ctx.Done()` and sets `ctx.Err()` to `context.Canceled`. Any goroutine selecting on `ctx.Done()` wakes up immediately.

```go
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
```

Real captured output:

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

The exact tick count before cancellation depends on scheduler timing and can vary slightly between runs or machines (you might see 4 or 5), but the shape is guaranteed: `main` blocks on `<-stopped` until the worker's `select` actually observes `ctx.Done()` and returns - this is a **real, verified stop**, not just "we told it to stop and hoped."

---

## context.WithTimeout and context.WithDeadline

`context.WithTimeout` cancels a context automatically after a **duration** elapses. `context.WithDeadline` does the same for an **absolute time**. Both return the same `(ctx, cancel)` shape as `WithCancel`, and both still let you call `cancel()` manually to stop early.

```go
ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
defer cancel()

start := time.Now()
<-ctx.Done()
elapsed := time.Since(start)

fmt.Println("ctx.Done() fired")
fmt.Println("ctx.Err():", ctx.Err())
fmt.Printf("elapsed: %v (requested timeout: 200ms)\n", elapsed.Round(time.Millisecond))
```

Real captured output (exact elapsed time varies by machine - scheduling always adds a small overhead above the requested duration; this is one real run):

```
ctx.Done() fired
ctx.Err(): context deadline exceeded
elapsed: 201ms (requested timeout: 200ms)
```

`context.WithDeadline` takes an absolute `time.Time` instead of a duration - useful when you're computing "must finish by X" rather than "must finish within X":

```go
deadline := time.Now().Add(150 * time.Millisecond)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()

<-ctx.Done()
fmt.Println("ctx.Err():", ctx.Err())

d, ok := ctx.Deadline()
fmt.Println("Deadline() ok:", ok)
fmt.Println("Deadline() time equals requested:", d.Equal(deadline))
```

Real captured output:

```
ctx.Done() fired
ctx.Err(): context deadline exceeded
elapsed: 151ms (deadline was 150ms from start)
Deadline() ok: true
Deadline() time equals requested: true
```

`context.WithTimeout(parent, d)` is implemented in terms of `context.WithDeadline(parent, time.Now().Add(d))` - they're the same mechanism, just expressed as a relative duration vs. an absolute instant. Use whichever reads more naturally at the call site.

---

## ctx.Err(): Canceled vs. DeadlineExceeded

`ctx.Err()` returns `nil` while a context is still active, and one of exactly two sentinel errors once it's done:

- `context.Canceled` - a `cancel()` function was called manually
- `context.DeadlineExceeded` - a timeout or deadline elapsed on its own

```go
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
```

Real captured output - both sentinel errors produced for real in the same run:

```
=== Manual Cancellation ===
ctx1.Err(): context canceled
errors.Is(ctx1.Err(), context.Canceled): true

=== Timeout Expiry ===
ctx2.Err(): context deadline exceeded
errors.Is(ctx2.Err(), context.DeadlineExceeded): true
```

Always compare with `errors.Is(ctx.Err(), context.Canceled)` rather than `ctx.Err() == context.Canceled` - it's the idiomatic form and remains correct even if a context implementation wraps the sentinel error. Distinguishing the two matters in real code: a `context.DeadlineExceeded` on an HTTP client call usually means "retry, this was slow," while a `context.Canceled` usually means "the caller gave up, don't retry."

---

## ctx.Done() and the select Pattern

The standard shape for any goroutine doing real work while respecting cancellation is a `select` between the context's `Done()` channel and whatever channel carries the actual result - exactly the "done channel" idea from Level 22's `select`, now standardized:

```go
select {
case <-ctx.Done():
    // canceled or timed out - stop, clean up, return ctx.Err()
case result := <-workCh:
    // real work finished first - use the result
}
```

```go
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
```

Run once with the work fast enough to win, and once with a timeout short enough that it wins instead - both real outcomes, captured verbatim:

```
=== Case 1: Work finishes before timeout ===
[fast work] work finished first: work result

=== Case 2: Timeout fires before work finishes ===
[slow work] canceled before work finished: context deadline exceeded
```

Notice `doWork` also selects on `ctx.Done()` when it tries to *send* the result. Without that inner `select`, a canceled caller that's no longer receiving would leave `doWork`'s goroutine blocked forever trying to send on an unbuffered channel nobody's reading anymore - a goroutine leak. Respecting `ctx.Done()` on every blocking operation, not just the outermost one, is what makes cancellation actually propagate all the way down.

---

## context.WithValue: Request-Scoped Data

`context.WithValue` attaches a key/value pair to a context, retrievable from any context derived from it:

```go
type requestIDKey struct{} // unexported, unique type - avoids collisions with any other package's keys

func withRequestID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, requestIDKey{}, id)
}

func requestIDFrom(ctx context.Context) (string, bool) {
    id, ok := ctx.Value(requestIDKey{}).(string)
    return id, ok
}
```

```go
ctx := withRequestID(context.Background(), "req-12345")
id, ok := requestIDFrom(ctx)
fmt.Println("requestID:", id, "found:", ok)
```

Real captured output:

```
=== Correct Pattern: Typed Unexported Key ===
requestID: req-12345 found: true
```

### 🚨 A Strong Warning: Don't Overuse context.Value

`context.Value` is tempting to reach for because it lets you skip changing a function signature - but that's exactly the problem. Values stored in a context are **untyped, unchecked by the compiler, and invisible in a function's signature**. A function that secretly needs `ctx.Value(someKey{})` to work correctly is hiding a real dependency that a caller has no way to discover except by reading its source.

**Use `context.Value` only for data that is genuinely request-scoped and cross-cutting** - things like a trace ID, a request ID, or an authenticated user identity that many unrelated layers need without explicitly threading it through every signature. **Never use it for optional parameters, configuration, or anything that changes a function's business behavior** - those belong as explicit function arguments, where the compiler and every reader of the signature can see them.

**Always use a typed, unexported key type - never a plain string or other exported/basic type.** Here's exactly why: two unrelated pieces of code that both happen to use the string `"id"` as a key will silently collide - one overwrites the other with no warning, at runtime, in production.

```go
// Two unrelated pieces of code both decide to use the string "id" as a context key.
ctxA := context.WithValue(context.Background(), "id", "package-A-value")
ctxB := context.WithValue(ctxA, "id", "package-B-value") // silently overwrites package A's value!

fmt.Println("Value for string key \"id\":", ctxB.Value("id"))
```

Real captured output - the collision actually happens:

```
=== The Collision Risk of Plain String Keys ===
Value for string key "id": package-B-value
package A's original value is now unreachable - overwritten by package B's use of the same string key
```

Now the same scenario with distinct, unexported key **types** instead of string values - no collision is possible, even though the underlying data is otherwise identical:

```go
type keyA struct{}
type keyB struct{}
ctxTyped := context.WithValue(context.Background(), keyA{}, "package-A-value")
ctxTyped = context.WithValue(ctxTyped, keyB{}, "package-B-value")
fmt.Println("keyA value:", ctxTyped.Value(keyA{}))
fmt.Println("keyB value:", ctxTyped.Value(keyB{}))
```

Real captured output:

```
=== Distinct Typed Keys Never Collide ===
keyA value: package-A-value
keyB value: package-B-value
both values survive - distinct unexported types can never collide, even carrying identical string data
```

An unexported struct type (`type requestIDKey struct{}`) can only ever be constructed inside the package that defines it, so no outside package - not even one that guesses the same "field name" or string - can ever produce a colliding key by accident.

---

## Always Call the Returned cancel Function

Every `With*` function that returns a `cancel` - `WithCancel`, `WithTimeout`, `WithDeadline` - **must have that `cancel` called**, typically via `defer cancel()` immediately after creating the context:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel() // always, even though the timeout will eventually fire on its own
```

This is true **even for `WithTimeout`/`WithDeadline`, which cancel themselves automatically once their time is up.** The reason to call `cancel` anyway is releasing resources **the instant you're actually done**, rather than leaving them alive until a timer that may be much longer than your real work finally fires.

Here's the difference made concrete. First, the version that discards `cancel` - a real, common mistake:

```go
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
```

Real captured output:

```
bad: finished a quick unit of work
main: done (the 5-second internal timer is still alive until it fires)
```

The program did its real work in 10 milliseconds, but that discarded `cancel` left the context's internal timer registered and alive for the full 5 seconds - doing nothing useful, just waiting to fire so it can clean itself up. Go's own `go vet` tool catches exactly this mistake, for real:

```
$ go vet ./...
main.go:19:10: the cancel function returned by context.WithTimeout should be called, not discarded, to avoid a context leak
```

Now the fixed version, with `defer cancel()`:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel() // released the instant main() returns - the timer is stopped immediately, not in 5s

doSomething(ctx)
fmt.Println("main: done (deferred cancel will release the timer right now, not in 5 seconds)")
```

Real captured output:

```
good: finished a quick unit of work
main: done (deferred cancel will release the timer right now, not in 5 seconds)
```

Running `go vet ./...` on this version reports nothing - it exits cleanly with no warnings, confirming the fix.

**What actually leaks if you skip it:** the context object and its internal timer stay registered in the runtime for as long as the timeout/deadline you specified, holding memory and doing scheduler bookkeeping, even after you're finished with the operation. For one short CLI program that exits in milliseconds, the process death cleans everything up anyway - the leak is harmless. But in a **long-running server** handling thousands of requests, each one deriving a `WithTimeout` context and forgetting to cancel means every single request leaves its context and timer alive until its own timeout eventually fires - across enough concurrent requests, that's real, accumulating memory and CPU overhead for work nobody needs anymore. `defer cancel()` costs nothing and prevents all of it.

---

## Propagating Context Through a Call Chain and Goroutines

A context created once at the top of an operation should be **passed down**, unchanged in identity, through every function and goroutine working on that operation's behalf - a fan-out pattern where canceling the one parent context stops every worker at once:

```go
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
```

Real captured output:

```
main: canceling parent context - all workers should stop
all 5 workers confirmed stopped after cancellation
```

Every one of the 5 workers received the **exact same `ctx`** - there's only one context object, shared by value (a `context.Context` is an interface value, cheap to copy and pass around). One `cancel()` call closes one `Done()` channel, and every `select` watching that channel across every goroutine wakes up at once. This is the payoff of formalizing the done-channel pattern: fan-out to any number of workers, and stopping all of them is always exactly one function call.

---

## A Simulated HTTP-Handler-Like Pattern

Real HTTP handlers (Level 27) receive a `context.Context` tied to the incoming request - it's canceled automatically if the client disconnects or the request times out. You can preview that exact shape today with a function that simulates multi-step work, checking `ctx` between steps and stopping cleanly and reporting how far it got:

```go
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
```

Run it once with plenty of time, and once with a timeout that cuts it off partway through:

```go
steps := []step{
    {"validate", 30 * time.Millisecond},
    {"fetch", 30 * time.Millisecond},
    {"transform", 30 * time.Millisecond},
    {"persist", 30 * time.Millisecond},
    {"notify", 30 * time.Millisecond},
}

ctx1, cancel1 := context.WithTimeout(context.Background(), 500*time.Millisecond)
defer cancel1()
completed1, err1 := runPipeline(ctx1, steps)

ctx2, cancel2 := context.WithTimeout(context.Background(), 70*time.Millisecond)
defer cancel2()
completed2, err2 := runPipeline(ctx2, steps)
```

Real captured output:

```
=== Run 1: Enough time for the whole pipeline ===
completed steps: [validate fetch transform persist notify]
error: <nil>

=== Run 2: Canceled partway through ===
completed steps: [validate fetch]
error: context deadline exceeded
pipeline stopped after 2 of 5 steps
```

This is precisely the shape a real HTTP handler takes: `r.Context()` gives you a context that's canceled the moment the client goes away, and every downstream step - a database query, a call to another service - should take that same `ctx` and stop as soon as it's done, rather than continuing to burn CPU and I/O on a response nobody will ever receive.

---

## Best Practices

### 1. Context Is Always the First Parameter, Named ctx

```go
// ✅ Good
func FetchUser(ctx context.Context, id string) (*User, error) { }

// ❌ Wrong - not first, not conventionally named
func FetchUser(id string, c context.Context) (*User, error) { }
```

This is such a strong Go convention that `go vet` and linters will flag violations. Every reader immediately knows where to look for the cancellation/deadline plumbing.

### 2. Never Store a Context in a Struct

```go
// ❌ WRONG - a stored context can outlive the request it belonged to,
// and different calls to the same struct's methods may need different contexts
type Client struct {
    ctx context.Context
}

// ✅ RIGHT - pass ctx into each method call that needs it
type Client struct{}
func (c *Client) DoRequest(ctx context.Context) error { }
```

A context describes one specific operation's lifetime. A struct's lifetime and a single operation's lifetime are different things - baking one into the other means either the context goes stale or the struct becomes unusable across unrelated calls.

### 3. Don't Use context.Value for Things That Should Be Explicit Parameters

If a piece of data changes a function's business logic, or is required for it to work correctly, put it in the function signature - not in the context. Reserve `context.Value` for cross-cutting, request-scoped metadata (trace IDs, request IDs) that truly doesn't belong to any single layer's explicit API. See the strong warning in [context.WithValue](#contextwithvalue-request-scoped-data) above.

### 4. Derive, Don't Recreate

Always derive a new context from the one you were given (`context.WithTimeout(ctx, ...)`), never discard it and start over from `context.Background()` partway through a call chain - see [Common Mistakes](#common-mistakes) below for what that breaks.

### 5. Check ctx.Done()/ctx.Err() in Any Loop or Long-Running Operation

Any function that loops, retries, or otherwise might run for a while should have cancellation baked into every blocking step - not just its own entry point. See [ctx.Done() and the select Pattern](#ctxdone-and-the-select-pattern).

---

## Common Mistakes

### Mistake 1: Forgetting defer cancel()

```go
// ❌ WRONG - the cancel func returned by context.WithTimeout should be called,
// not discarded, to avoid a context leak (this is a real go vet warning)
ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

// ✅ RIGHT
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

### Mistake 2: Overusing context.Value

```go
// ❌ WRONG - business parameters hidden in the context instead of the signature
func ProcessOrder(ctx context.Context) error {
    userID := ctx.Value("userID").(string) // hidden dependency, plain string key
    discount := ctx.Value("discount").(float64) // this should just be an argument!
    // ...
}

// ✅ RIGHT - explicit parameters the compiler and every reader can see
func ProcessOrder(ctx context.Context, userID string, discount float64) error {
    // ...
}
```

### Mistake 3: Using context.Background() Everywhere Instead of Propagating

```go
// ❌ WRONG - creates a brand-new, never-cancelable context, breaking the
// cancellation chain from every caller above this function
func fetchData() error {
    ctx := context.Background()
    return db.QueryContext(ctx, "...")
}

// ✅ RIGHT - accept and forward the real parent context
func fetchData(ctx context.Context) error {
    return db.QueryContext(ctx, "...")
}
```

Once you call `context.Background()` partway down a call chain, canceling the real top-level context stops having any effect on everything below that point - the chain is broken.

### Mistake 4: Not Checking ctx.Err()/ctx.Done() in Long-Running Work

```go
// ❌ WRONG - takes a context but never actually looks at it; cancellation
// is accepted as a parameter and then completely ignored
func process(ctx context.Context, items []string) {
    for _, item := range items {
        expensiveWork(item) // keeps going even after the caller gave up
    }
}

// ✅ RIGHT - checked every iteration, so cancellation actually stops the work
func process(ctx context.Context, items []string) {
    for _, item := range items {
        select {
        case <-ctx.Done():
            return
        default:
        }
        expensiveWork(item)
    }
}
```

A function that accepts a `context.Context` parameter but never reads `Done()` or `Err()` is worse than one that doesn't take a context at all - it looks cancellable to every caller, but isn't.

---

## Summary

**The Two Root Contexts:**
- `context.Background()` - the real root, used at genuine top-level entry points
- `context.TODO()` - a placeholder signaling "proper context plumbing still needed here"

**Derivation:**
- `context.WithCancel(parent)` - manual cancellation via the returned `cancel()`
- `context.WithTimeout(parent, d)` / `context.WithDeadline(parent, t)` - automatic cancellation after a duration/absolute time
- Always `defer cancel()`, even for the self-canceling `With*` functions

**Inspecting State:**
- `ctx.Done()` - a channel that closes when the context is canceled or times out
- `ctx.Err()` - `nil` while active; `context.Canceled` or `context.DeadlineExceeded` once done
- `ctx.Value(key)` - request-scoped data, retrieved with a typed unexported key

**The Pattern:**
```go
select {
case <-ctx.Done():
    return ctx.Err()
case result := <-workCh:
    // use result
}
```

**Propagation:**
- One context, shared by value, fans out to any number of goroutines
- Cancel it once, and every `select` watching `Done()` wakes up together

---

## Next Steps

You now understand:
- ✅ What context.Context carries, and why it formalizes Level 22's done-channel pattern
- ✅ context.Background() vs. context.TODO(), and when to use each
- ✅ context.WithCancel, context.WithTimeout, and context.WithDeadline
- ✅ Distinguishing context.Canceled from context.DeadlineExceeded via ctx.Err()
- ✅ The ctx.Done()-select pattern for real, cancellable work
- ✅ context.WithValue's narrow, typed-key-only correct use - and its overuse risk
- ✅ Why defer cancel() matters even for self-canceling contexts
- ✅ Propagating one context through a call chain and a fan-out of goroutines

**Next level:** Level 25 - Generics
- Type parameters for functions and types
- Type constraints and the `comparable`/`any` built-in constraints
- Writing generic data structures and algorithms
- When generics help vs. when a plain interface is simpler

You can now write concurrent Go that stops cleanly on command instead of running forever. Generics are next - a completely different kind of flexibility: writing code once that works safely across many types. Keep going! 🚀
