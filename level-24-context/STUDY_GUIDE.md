# Level 24: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: Context Fundamentals
```
Day 1:  What context.Context is for; Background() vs TODO()
Day 2:  WithCancel - manual cancellation, verified with a real goroutine
Day 3:  WithTimeout and WithDeadline - real elapsed-time verification
Day 4:  ctx.Err() - Canceled vs DeadlineExceeded
Day 5:  ctx.Done() and the select pattern
Day 6:  WithValue - typed keys and the collision risk
Day 7:  defer cancel(), propagation, and the pipeline preview
```

### Week 2: Practice & Application
```
Day 1:  Exercises 1-3 (root contexts, WithCancel, WithTimeout)
Day 2:  Exercises 4-6 (WithDeadline, Err() distinction, select pattern)
Day 3:  Exercises 7-8 (WithValue, defer cancel() and go vet)
Day 4:  Exercises 9-10 (fan-out propagation, comprehensive pipeline)
Day 5:  Bonus challenges (retry-with-backoff, cancellable pool, first-wins)
Day 6-7: Review & consolidation
```

---

## 🌳 The Context Tree

```
                     context.Background()
                       (never canceled,
                        no deadline, no value)
                             │
              ┌──────────────┼──────────────┐
              │                              │
      context.WithCancel(parent)    context.WithTimeout(parent, d)
      returns (ctx, cancel)          returns (ctx, cancel)
              │                              │
      ┌───────┴───────┐              ┌───────┴───────┐
      │               │              │               │
  worker 1        worker 2       WithValue(ctx, k, v)  ... more children

  Canceling a context cancels EVERY context derived from it -
  cancellation flows DOWN the tree, never up, never sideways.
  Canceling a child affects only that child and ITS children,
  never its parent or siblings.
```

Every `With*` function takes a parent `context.Context` and returns a new **child** context. A child's `Done()` channel closes when *either* it is directly canceled/times out, *or* any of its ancestors is. That's the entire mechanism behind fan-out: cancel one context near the root, and every descendant - however many goroutines are holding a reference to it - finds out.

---

## 🎛️ WithCancel vs. WithTimeout vs. WithDeadline

| | `WithCancel(parent)` | `WithTimeout(parent, d)` | `WithDeadline(parent, t)` |
|---|---|---|---|
| Cancels when | you call the returned `cancel()` | `d` elapses, or `cancel()` is called early | wall-clock reaches `t`, or `cancel()` is called early |
| Input | nothing beyond the parent | a `time.Duration` (relative) | a `time.Time` (absolute) |
| `ctx.Err()` after firing | `context.Canceled` | `context.DeadlineExceeded` (or `Canceled` if you canceled it manually first) | `context.DeadlineExceeded` (or `Canceled` if canceled manually first) |
| Typical use | you decide, at runtime, when to stop (e.g., after a fan-out winner is found) | "this must finish within N seconds" | "this must finish by 5:00pm" / a computed absolute cutoff |
| Must still `defer cancel()`? | Yes - it's your only way to release it | **Yes**, even though it self-cancels eventually | **Yes**, even though it self-cancels eventually |

`WithTimeout(parent, d)` is implemented as `WithDeadline(parent, time.Now().Add(d))` internally - they are the same mechanism, expressed two different ways.

---

## 🔀 The ctx.Done()-select Pattern

```
   select {
   case <-ctx.Done():
       // parent canceled OR timed out - stop now
       return ctx.Err()
   case result := <-workCh:
       // real work finished first - use it
   }

        ┌─────────────┐          ┌─────────────┐
        │  ctx.Done()  │          │   workCh     │
        │ (closes on   │          │ (real result │
        │  cancel/     │          │  arrives)    │
        │  timeout)    │          │              │
        └──────┬───────┘          └──────┬───────┘
               │                          │
               └────────────┬─────────────┘
                             │
                        select picks
                      WHICHEVER FIRES
                            FIRST
```

This is Level 22's done-channel `select` pattern, standardized: instead of a bespoke `chan struct{}` per operation, every function that takes a `context.Context` shares the exact same `Done()` shape - so any function, from any package, composes with any other.

**Rule:** any goroutine that blocks on a channel operation while doing work on behalf of a context should have a `<-ctx.Done()` case at *every* one of those blocking points - not just its outermost loop - or a canceled caller can still leave it stuck forever.

---

## 🗝️ context.Value: String Key vs. Typed Key

```
STRING KEY (❌ collision risk)

  package A:  ctx = context.WithValue(ctx, "id", "A's value")
  package B:  ctx = context.WithValue(ctx, "id", "B's value")
                                        ^^^^
                            SAME KEY - B silently overwrites A!

  ctx.Value("id")  →  "B's value"     (A's value is gone, no error, no warning)


TYPED KEY (✅ collision-proof)

  package A:  type keyA struct{}
              ctx = context.WithValue(ctx, keyA{}, "A's value")
  package B:  type keyB struct{}
              ctx = context.WithValue(ctx, keyB{}, "B's value")

  ctx.Value(keyA{})  →  "A's value"   (each type is its own namespace)
  ctx.Value(keyB{})  →  "B's value"   (both survive - impossible to collide)

  An unexported struct type can only be constructed inside the
  package that defines it - no outside code can ever forge a
  colliding key, even by accident.
```

**Rule of thumb:** if you can imagine two different packages picking the exact same variable/field name for a context key, a plain `string` key WILL eventually collide in a large enough codebase. A typed unexported key makes that mathematically impossible.

---

## 💧 The defer cancel() Leak, Visualized

```
❌ WITHOUT defer cancel()

  ctx, _ := context.WithTimeout(parent, 5*time.Second)
  doQuickWork(ctx)   // finishes in 10ms
  // ... 4.99 seconds of nothing ...
  // the internal timer is STILL alive, still registered,
  // doing nothing useful, until it finally self-cancels at 5s

  go vet catches this:
    "the cancel function returned by context.WithTimeout
     should be called, not discarded, to avoid a context leak"


✅ WITH defer cancel()

  ctx, cancel := context.WithTimeout(parent, 5*time.Second)
  defer cancel()
  doQuickWork(ctx)   // finishes in 10ms
  // main() returns → deferred cancel() runs IMMEDIATELY →
  // timer released right now, not in 5 seconds

  go vet: silent, no warning
```

In a single short-lived program, the leaked timer is harmless - the process exits and the OS reclaims everything anyway. In a **long-running server** handling many requests, every forgotten `cancel()` is one more context + timer that outlives the work it was for - multiply that by thousands of requests and it becomes real, measurable memory and scheduler overhead.

---

## 🌐 Fan-Out Propagation

```
                context.WithCancel(parent)
                    ctx, cancel
                         │
        ┌────────────────┼────────────────┬────────────────┐
        │                │                │                │
   worker 1          worker 2         worker 3         worker 4
   select {          select {         select {         select {
   <-ctx.Done()      <-ctx.Done()     <-ctx.Done()     <-ctx.Done()
   <-time.After()    <-time.After()   <-time.After()   <-time.After()
   }                 }                }                }

   main calls cancel() ONCE
        │
        ▼
   ctx.Done() closes ONCE
        │
        ▼
   EVERY worker's select wakes up on its <-ctx.Done() case,
   simultaneously, because they all share the SAME ctx value
   (a context.Context is a small interface value - cheap to
   copy, safe to share across any number of goroutines)
```

---

## 🚨 Common Mistakes

| Mistake | Wrong | Right |
|---|---|---|
| Forgetting defer cancel() | `ctx, _ := context.WithTimeout(...)` | `ctx, cancel := context.WithTimeout(...); defer cancel()` |
| Overusing context.Value | hiding required business parameters (`userID`, `discount`) behind `ctx.Value(...)` | pass them as explicit function arguments |
| Using Background() mid-chain | `ctx := context.Background()` inside a function that received a real `ctx` parameter | accept and forward the caller's `ctx` |
| Not checking Done()/Err() | a function takes `ctx` but never reads it in its loop | check `<-ctx.Done()` at every blocking step of long-running work |
| String context keys | `context.WithValue(ctx, "id", v)` | `context.WithValue(ctx, myKeyType{}, v)` with an unexported type |
| Storing ctx in a struct | `type Client struct { ctx context.Context }` | pass `ctx` into each method call instead |

---

## 📈 Progression Summary

### Understanding Level 24

Level 24 formalizes Level 22's done-channel pattern into a standard, composable type:

1. **Root contexts** - Background() for real roots, TODO() as an honest placeholder
2. **Derivation** - WithCancel, WithTimeout, WithDeadline all return `(ctx, cancel)`
3. **Inspection** - Done() to wait, Err() to learn why, Deadline() to check the cutoff
4. **Values** - a narrow, typed-key-only escape hatch for cross-cutting request data
5. **Discipline** - always defer cancel(), always check Done()/Err() in long-running work
6. **Propagation** - one context, shared by value, fans out to any number of goroutines

### Prerequisites for Level 25

Before moving to Level 25 (Generics), you need:

- ✅ Comfortable creating and canceling contexts with WithCancel
- ✅ Comfortable using WithTimeout/WithDeadline and reading ctx.Err()
- ✅ Can write the ctx.Done()-select pattern from memory
- ✅ Understand the context.Value collision risk and the typed-key fix
- ✅ Can explain exactly what defer cancel() prevents
- ✅ Can propagate one context through a fan-out of goroutines

### Ready for Level 25?

Level 25 teaches a completely different kind of flexibility - writing code once that works safely across many types:
- Type parameters for functions and types
- Type constraints, `comparable`, and `any`
- Generic data structures and algorithms
- When generics help vs. when a plain interface is simpler

---

## ✅ Checklist Before Level 25

- [ ] Can explain the difference between Background() and TODO()
- [ ] Can cancel a running goroutine with WithCancel and verify the stop for real
- [ ] Can use WithTimeout/WithDeadline and read the real elapsed time
- [ ] Can distinguish context.Canceled from context.DeadlineExceeded via ctx.Err()
- [ ] Can write the ctx.Done()-select pattern without hesitation
- [ ] Can explain why context.Value needs a typed, unexported key
- [ ] Always writes defer cancel() immediately after every With* call
- [ ] Can propagate one context through several fanned-out goroutines
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The Formalization
`context.Context` takes Level 22's done-channel idea and gives it a standard shape every function, package, and goroutine can share.

### The Discipline
`defer cancel()` immediately after every `With*` call - always, even for self-canceling timeouts and deadlines.

### The Narrow Escape Hatch
`context.Value` is for cross-cutting, request-scoped metadata only - never for parameters your function actually needs to do its job.

### The Propagation
A context is a small, cheap interface value - pass the real one down, never manufacture a fresh `Background()` partway through a call chain.

---

## 📚 Next Level

Level 25: Generics
- Type parameters for functions and types
- Type constraints and the comparable/any built-in constraints
- Writing generic data structures and algorithms
- When generics help vs. when a plain interface is simpler

You've got context down! Keep going! 🚀
