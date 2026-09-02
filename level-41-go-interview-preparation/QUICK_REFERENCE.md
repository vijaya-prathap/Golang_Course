# Level 41: Quick Reference Card

## 🚀 Interview-Day Cheat Sheet (5 Minutes)

```bash
# The night before: skim this whole file once.
# The morning of: re-read the gotchas table and the one-sentence concepts below.
# Right before: take a breath - you know this material, you've verified all of it yourself.
```

---

## 🚨 Go Gotchas (Condensed)

| Gotcha | One-Line Fix |
|---|---|
| Slice/append aliasing | `copy()` for an independent slice |
| Loop variable capture | Go 1.22+ = per-iteration copy (check `go.mod`!) |
| nil interface vs typed-nil | Never return a concrete nil pointer as an interface |
| Map iteration order | Sort keys if order matters |
| Uncomparable struct `==` | Use `reflect.DeepEqual` |
| Goroutine leaks | Always give a goroutine a cancellation path |
| Channel deadlock | Never send/receive with no possible other end |
| nil channel | Blocks forever - useful for disabling a `select` case |

---

## 🧩 Coding-Challenge Checklist

Before the interview, be able to write each of these cold, in under 10 minutes:

- [ ] Reverse a string by runes (`[]rune`, two-pointer swap)
- [ ] Detect a cycle in a linked list (slow/fast pointers)
- [ ] Find the first non-repeating character (frequency map + ordered pass)
- [ ] Implement an LRU cache (`map` + `container/list`)
- [ ] Merge overlapping intervals (sort by start, sweep once)
- [ ] Build a worker pool (`jobs` channel in, `results` channel out, `sync.WaitGroup`)

```go
// Two-pointer skeleton - reused across reversal, palindrome, and merge-style problems
for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
    x[i], x[j] = x[j], x[i]
}
```

```go
// Worker pool skeleton
jobs := make(chan Job, n)
results := make(chan Result, n)
var wg sync.WaitGroup
for w := 0; w < numWorkers; w++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        for j := range jobs {
            results <- process(j)
        }
    }()
}
// ... send jobs, then close(jobs) ...
go func() { wg.Wait(); close(results) }()
for r := range results { /* collect */ }
```

---

## 💬 One-Sentence Answers to Key Concepts

- **Value vs pointer receiver:** pointer when you mutate or the type is large; value for small, copy-cheap types - and be consistent across a type's whole method set.
- **Garbage collector:** concurrent, tri-color mark-and-sweep - runs alongside your program, not fully stop-the-world.
- **Goroutine vs OS thread:** a lightweight, runtime-scheduled function with a tiny growable stack, multiplexed onto far fewer real OS threads.
- **Empty interface (`any`):** satisfied by everything, but every use needs a type assertion - generics (Level 25) are the type-safe alternative in most cases.
- **defer/panic/recover:** `defer` runs at function return (LIFO); `panic` unwinds the stack; `recover` (only inside a deferred call) stops the unwind and lets the function return normally.
- **Interface satisfaction:** structural - any type with the right method set satisfies an interface automatically, no `implements` keyword.
- **Unbuffered vs buffered channel:** unbuffered is a synchronization handoff (send blocks until received); buffered decouples timing up to its capacity, but a successful send no longer proves anyone received it.

---

## 🎯 Interview Structure, In Four Words

```
Clarify → Design → Code → Tradeoffs
```

Ask before you assume. Sketch before you dive. Narrate while you build. Discuss what you'd change before you call it done.

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---|---|---|
| Memorizing without understanding | Reciting "nil interface" rule from memory | Being able to explain the type+value pair underneath it |
| Skipping clarifying questions | Coding immediately on an ambiguous prompt | Asking about scale, edge cases, and scope first |
| Ignoring edge cases | Only handling the happy path | Naming empty input, nil, and zero values out loud |
| Not testing before "done" | "I think this works" | Running it, or tracing it by hand, before declaring it finished |

---

## 🎓 Final Self-Check

Can you:
- [ ] Explain all eight gotchas above without looking at the table?
- [ ] Implement any of the six coding challenges cold, in under 10 minutes?
- [ ] Structure a system design answer as clarify → design → deep dive → tradeoffs?
- [ ] Say "I don't know, but here's how I'd find out" without it feeling like failure?

If YES to all four → you're ready to walk into a Go interview.

---

## 🎉 You've Completed the Course

There is no "Next Level" here - this is the last page of a 42-level course, not a checkpoint on the way to another one.

You started at Level 0 asking what a program even is. You're finishing at Level 41 able to reason about interface internals, fix a goroutine leak, and structure a system design answer under time pressure. That's the whole arc.

Go get the job - and then go keep building. 🚀
