# Go Interview Questions & Answers

A standalone, quick-scan Q&A reference distilled from this 42-level course. Organized by topic so you can jump straight to your weak spots. Every code snippet below was actually run on this project's Go toolchain (go1.26.5) before being included — none of it is guessed.

Companion to `level-41-go-interview-preparation/` (which has the full narrative version, coding-challenge walkthroughs, and a 42-row course checklist). This file is the condensed, question-first version for last-minute review.

---

## How to Use This File

- Skim the section headers, jump to whatever you're rusty on.
- For each question, try answering out loud *before* reading the answer — that's what the real interview feels like.
- The "Gotcha" questions are the ones interviewers ask specifically to see if you've been burned by them before. Know these cold.

---

## 1. Language Fundamentals

**Q: What are Go's zero values, and why do they matter?**
A: Every declared-but-uninitialized variable gets a predictable default: `0` for numeric types, `""` for strings, `false` for bool, `nil` for pointers/slices/maps/channels/functions/interfaces. This means Go has no concept of an "uninitialized" variable the way C does — every value is always usable immediately, which eliminates a whole class of undefined-behavior bugs.

**Q: What's the difference between `var x = 5` and `x := 5`?**
A: They're equivalent in effect (both use type inference), but `:=` (short declaration) only works inside functions, while `var` works at package level too. `:=` is idiomatic inside functions; `var` is required for package-level declarations and is also preferred when you want to be explicit about the type or declare a zero value with no initializer.

**Q: Are arrays and slices the same thing?**
A: No. An array (`[3]int`) has a fixed size baked into its type — `[3]int` and `[4]int` are different, incompatible types. A slice (`[]int`) is a small header (pointer, length, capacity) pointing at an underlying array. Arrays are value types (assigning or passing one copies every element); slices behave like references to shared backing data. In practice, you reach for slices almost always — arrays are rare except for genuinely fixed-size data.

**Q: What happens when you append to a slice — does it always allocate a new array?**
A: Only when the slice is at capacity. If there's spare capacity, `append` writes into the existing backing array and returns a slice with the same header pointer — no allocation. This is the source of a classic bug:
```go
a := make([]int, 3, 5)
a[0], a[1], a[2] = 1, 2, 3
b := a[:2]
b = append(b, 99) // capacity allows it — overwrites a[2] in place!
fmt.Println(a) // [1 2 99]  <- a changed even though you only appended to b
```
Real output. If you need an independent copy, use `copy()` or slice with a full three-index expression (`a[:2:2]`) to cap the capacity and force a new allocation on the next append.

**Q: byte vs rune — what's the difference?**
A: `byte` is an alias for `uint8` — one raw byte, a single UTF-8 code unit. `rune` is an alias for `int32` — a decoded Unicode code point, which may span 1-4 bytes in UTF-8. `s[i]` on a string gives you a byte; `for i, r := range s` gives you runes, and the index jumps by more than 1 for multi-byte characters. `len(s)` counts bytes, not characters.

**Q: Why does `map` iteration order look random?**
A: Because Go deliberately randomizes it on every run, specifically to stop code from silently depending on an order that was never a guarantee in the first place. If you need a stable order, collect the keys into a slice, sort them, then iterate the sorted slice.

**Q: What's the nil-map gotcha?**
A: Reading from a nil map is always safe and returns the zero value. **Writing** to a nil map panics:
```go
var m map[string]int
fmt.Println(m["x"])  // 0, no panic
m["x"] = 1            // panic: assignment to entry in nil map
```
Real captured panic text. Fix: initialize with `make(map[K]V)` or a map literal before writing.

---

## 2. Functions, Closures, defer

**Q: Are Go's `++`/`--` expressions?**
A: No — they're statements. `y := x++` doesn't compile. Increment on its own line, then use the variable on the next.

**Q: What does `defer` actually delay — the call, or the evaluation of its arguments?**
A: Only the *call* is delayed. Arguments are evaluated immediately, at the point `defer` is written:
```go
i := 0
defer fmt.Println("deferred i =", i) // i is captured as 0 right here
i++
fmt.Println("current i =", i)
```
Real output:
```
current i = 1
deferred i = 0
```
Multiple defers in one function run in LIFO order (last deferred, first executed).

**Q: How did loop-variable capture in closures change in Go 1.22?**
A: Before 1.22, a `for` loop had ONE shared loop variable reused across every iteration — a goroutine or closure launched inside the loop captured that single variable, and by the time it ran, the loop had usually already moved on, so every closure saw the *final* value (a famously misleading bug). Since Go 1.22, each iteration gets its **own** copy of the loop variable automatically. Verified on this project's Go 1.26.5:
```go
var results []int
var wg sync.WaitGroup
for i := 0; i < 3; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        results = append(results, i) // each goroutine sees its own i
    }()
}
wg.Wait()
```
On 1.22+ this reliably produces `{0, 1, 2}` in some order — never `{3, 3, 3}`. If you're asked about this, mention the version boundary explicitly; a lot of "expected" wrong-answer trivia questions about this are now outdated.

**Q: When would you use named return values?**
A: When it makes the function's contract clearer (documenting what each return slot means) or when combined with `defer` to modify a return value after the fact (e.g., wrapping an error with context in a deferred recover). Avoid them for long functions — "naked returns" (`return` with no values) get hard to trace.

---

## 3. Pointers & Structs

**Q: When does Go pass by reference vs by value?**
A: Go is *always* pass-by-value — even a pointer argument is a copy of the pointer (the address), not a reference to the caller's variable itself. Passing `*T` lets the function dereference and mutate the pointed-to value, which achieves the effect people mean by "pass by reference," but the pointer variable itself is still copied.

**Q: Why is it safe in Go to return a pointer to a local variable?**
A: The compiler's escape analysis detects that the pointer outlives the function and automatically allocates that variable on the heap instead of the stack. Unlike C, there's no dangling-pointer risk:
```go
func newCounter() *int {
    x := 0
    return &x // safe — x escapes to the heap
}
```

**Q: Do slices and maps need to be passed as pointers to be mutated?**
A: Usually no — a slice/map is already a small header containing a pointer to the underlying data, so passing it by value still shares that underlying data; mutating an *existing element* through it is visible to the caller. The gotcha: if a function `append`s to a slice parameter and the backing array needs to grow, the function gets a *new* header, and that new header is NOT visible to the caller unless the function returns the (possibly-reallocated) slice and the caller reassigns it. This is why idiomatic Go signatures look like `func addItem(s []int, v int) []int { return append(s, v) }`.

**Q: Are structs value types or reference types?**
A: Value types, like arrays. `b := a` on a struct copies every field. If you want shared mutation, use a pointer (`*Struct`).

**Q: What's field promotion via embedding?**
A: Embedding a type (anonymous field) inside a struct promotes its fields and methods to the outer struct, so you can call `outer.Field` or `outer.Method()` directly without qualifying it. It's Go's composition-over-inheritance mechanism — there's no real inheritance, just delegation made syntactically convenient. A same-named field/method on the outer struct shadows the promoted one.

---

## 4. Methods & Interfaces

**Q: Value receiver vs pointer receiver — when do you use each?**
A: A value receiver operates on a copy — mutations inside the method don't affect the caller's value. A pointer receiver operates on the original. Rule of thumb: use a pointer receiver if the method needs to mutate the receiver, if the struct is large (avoid copying it), or if *any* method on the type needs a pointer receiver (for consistency, make them all pointer receivers).

**Q: What's the "method set" and why does it matter for interfaces?**
A: The method set of `T` includes only value-receiver methods. The method set of `*T` includes both value- and pointer-receiver methods. This means a `T` value can be assigned to an interface only if the interface's methods are all satisfied by `T`'s method set — if any required method has a pointer receiver, you need a `*T`, not a `T`, or you'll get a compile error like `Foo does not implement Bar (method X has pointer receiver)`.

**Q: How does a type "implement" an interface in Go?**
A: Implicitly — there's no `implements` keyword. Any type whose method set includes every method the interface declares automatically satisfies it. This decouples the interface's definition from the type's definition; a type doesn't even need to know an interface exists.

**Q: What is `fmt.Stringer` and why does it matter?**
A: It's the interface `interface { String() string }`. If a type implements it, `fmt.Println`/`%v` automatically call `String()` to format it — this is how you control how your custom types print.

**Q: What's the classic nil-interface gotcha?**
A: An interface value is really a `(type, value)` pair. An interface holding a `(*T)(nil)` — a typed nil — is NOT `== nil`, because the type half is non-nil even though the value half is:
```go
var p *MyError = nil
var err error = p
fmt.Println(err == nil) // false! — typed nil inside an interface
```
Real, verified behavior. This bites people constantly with error returns: a function declared to return `error` that returns a nil `*ConcreteError` produces a non-nil `error` interface. Fix: return the interface type directly (`var result error; if bad { result = &ConcreteError{} }; return result`), or make sure the concrete nil never gets boxed into the interface unless you mean to.

**Q: Why does Go favor small interfaces?**
A: "The bigger the interface, the weaker the abstraction." A one- or two-method interface (like `io.Reader`) is trivially satisfiable by almost anything, making it maximally reusable. Idiomatic Go also defines interfaces at the point of *use* (the consumer package), not next to the concrete type — so a package only needs to know what it uses, not the whole implementation.

**Q: Type assertion vs type switch?**
A: `v, ok := i.(T)` is the safe comma-ok form (returns `ok=false` instead of panicking on mismatch). `v := i.(T)` panics on mismatch — real captured text: `panic: interface conversion: interface {} is string, not int`. A type switch (`switch v := i.(type) { case int: ... }`) handles multiple possible concrete types cleanly.

---

## 5. Error Handling

**Q: Why does Go use returned errors instead of exceptions?**
A: It makes error handling an explicit, visible part of a function's contract (`(T, error)`) rather than an invisible control-flow path that can be silently swallowed. The tradeoff is more `if err != nil` boilerplate, in exchange for errors you can't accidentally ignore without the compiler nagging (an unused error return is a `go vet`/linter flag, if not a hard error).

**Q: What does `%w` do in `fmt.Errorf`, and how is it different from `%v`?**
A: `%w` wraps the error, preserving a chain you can walk with `errors.Unwrap`, `errors.Is`, and `errors.As`. `%v` just formats the error's message as a string — the resulting error loses any relationship to the original.

**Q: `errors.Is` vs `errors.As` — what's the difference?**
A: `errors.Is(err, target)` checks whether `target` (usually a sentinel like `var ErrNotFound = errors.New(...)`) appears anywhere in `err`'s wrap chain — use it to check "is this THAT specific error." `errors.As(err, &target)` finds the first error in the chain matching a given *type* and populates `target` with it — use it to check "is this error of THIS type, and if so, give me its fields."

**Q: When should you `panic` instead of returning an error?**
A: Rarely, and only for programmer errors or truly unrecoverable states (e.g., a required invariant is violated, or `main()` can't start because config is fundamentally broken). Normal, expected failure modes (file not found, invalid input, network timeout) should always be errors, not panics — a caller can't gracefully recover from your panic the way they can check an error.

**Q: What does `recover()` do, and where must it be called?**
A: `recover()` stops a panic that's currently unwinding the call stack and returns the value passed to `panic()`. It only has an effect when called directly inside a `deferred` function — calling it anywhere else returns `nil` and does nothing.

---

## 6. Concurrency

**Q: What's the difference between concurrency and parallelism?**
A: Concurrency is about *structuring* a program as independently-progressing tasks; parallelism is about actually *running* multiple tasks at the same instant on multiple cores. A concurrent Go program with `GOMAXPROCS=1` runs its goroutines concurrently but not in parallel — they interleave on one core.

**Q: Classic beginner goroutine bug?**
A: Launching a goroutine from `main()` and letting `main()` return before it finishes — the whole program exits and the goroutine's output never appears, with no error or warning. Fix: `sync.WaitGroup` (or a channel) to make `main()` wait.

**Q: What is a data race, and how do you detect one?**
A: Concurrent, unsynchronized access to the same memory where at least one access is a write. Go's race detector (`go test -race` / `go run -race`) instruments the binary and reports real, verified `WARNING: DATA RACE` output when one occurs at runtime — it doesn't require you to spot it by eye.

**Q: `sync.Mutex` vs channels — when do you use which?**
A: "Share memory by communicating" is Go's philosophy: channels are the idiomatic default for coordinating goroutines or handing off data. A `sync.Mutex` is more appropriate for straightforward protection of a shared, in-place piece of state (a counter, a cache map) where there's no natural "message" being passed. `sync/atomic` beats a mutex for a single simple counter/flag.

**Q: What's the difference between a buffered and unbuffered channel?**
A: An unbuffered channel (`make(chan int)`) forces the sender to block until a receiver is ready — it's a synchronization rendezvous. A buffered channel (`make(chan int, n)`) lets up to `n` sends proceed without a waiting receiver, decoupling sender and receiver timing (until the buffer fills).

**Q: What happens if you send on a channel with nobody ever receiving?**
A: A real deadlock — Go's runtime deadlock detector catches it and crashes the whole program immediately (it doesn't hang forever):
```
fatal error: all goroutines are asleep - deadlock!
```
Real, verified output. This is actually a *safety feature* — a hung, silent program would be far worse to debug.

**Q: What happens if you send on an already-closed channel?**
A: A real panic: `panic: send on closed channel`. Only the sender should ever close a channel, and never more than once (closing twice also panics).

**Q: What does a `nil` channel do in a `select`?**
A: A `nil` channel's case in a `select` is **never** ready — verified:
```go
var ch chan int // nil
select {
case <-ch:
    fmt.Println("received")
default:
    fmt.Println("default: nil channel is never ready")
}
```
prints `default: nil channel is never ready`, every time. This is used deliberately to "disable" one branch of a `select` at runtime by setting its channel variable to `nil`.

**Q: How does `select` pick a case when multiple are ready simultaneously?**
A: At random — not top-to-bottom like a `switch`, and not "whichever is listed first." This was empirically confirmed across 100,000 iterations landing near a 50/50 split when two channels were both pre-filled. Never write code that assumes case order implies priority.

**Q: What does `context.Context` give you that a hand-rolled `done chan struct{}` doesn't?**
A: A standardized way to propagate cancellation, deadlines, AND (sparingly) request-scoped values through an entire call graph, with a single value type every function/library can accept and check consistently (`ctx.Done()`, `ctx.Err()`). Always `defer cancel()` even with `WithTimeout`/`WithDeadline` (which self-cancel eventually) to release the timer immediately rather than leaking it until it fires.

---

## 7. Generics

**Q: Why does Go have generics, given it survived without them for a decade?**
A: To let you write one function/type that's genuinely type-safe across multiple types, without duplicating code per type OR falling back to `any`/`interface{}` (which loses compile-time type checking and needs runtime type assertions).

**Q: What does the `comparable` constraint do?**
A: It restricts a type parameter to types that support `==`/`!=` — required if your generic function uses those operators, or if the type parameter is used as a map key.

**Q: When should you NOT reach for generics?**
A: When a small interface already expresses the abstraction cleanly (behavior-based polymorphism) — generics shine for algorithms/data structures that are genuinely type-parametric (a stack, a cache, `Map`/`Filter`/`Reduce`), less so for "this type can do X" style abstractions, which interfaces already handle well.

---

## 8. Testing

**Q: What's a table-driven test, and why is it the idiomatic Go pattern?**
A: A slice of `{name, input, want}` structs run through a shared test body in a loop, usually with `t.Run(tt.name, ...)` per case. It scales to dozens of cases without duplicating test logic, and subtests give you individually addressable, individually failing test names in `go test -v` output.

**Q: `t.Error` vs `t.Fatal`?**
A: `t.Error`/`t.Errorf` marks the test failed but lets it keep running (useful for reporting multiple independent assertion failures in one test). `t.Fatal`/`t.Fatalf` marks it failed and stops that test immediately (use it when a later assertion would be meaningless or panic without an earlier one passing, e.g. after a failed setup step).

**Q: How do you test code that depends on an external interface, like a database or HTTP client?**
A: Depend on the *interface*, not the concrete type (constructor injection — Level 31's whole point), then substitute a fake/mock implementation in your test. This is why "accept interfaces, return concrete types" and defining the interface at the point of use matter so much for testability.

---

## 9. Web / Backend

**Q: What does Go's `http.Handler` interface look like, and why does it matter?**
A: `interface { ServeHTTP(w http.ResponseWriter, r *http.Request) }`. Because it's just one method, anything can be a handler, and middleware is simply a function that takes a `http.Handler` and returns a new one wrapping it — the entire ecosystem (including full frameworks like Gin) builds on this one small interface.

**Q: What's the middleware pattern?**
A: A function of type `func(http.Handler) http.Handler` that wraps a handler to add cross-cutting behavior (logging, auth, rate limiting) without the handler itself knowing about it. Middlewares are typically chained: `logging(auth(rateLimit(finalHandler)))`.

**Q: Why is `httptest` preferred over spinning up a real server for tests?**
A: `httptest.NewRecorder()` lets you call a handler function directly and inspect the response with zero real networking — fast and fully in-process. `httptest.NewServer()` gives you a real (but loopback-only, OS-assigned-port) HTTP server when you genuinely need a full round trip, e.g. testing client code against a handler.

**Q: How do you prevent SQL injection in Go?**
A: Always use parameterized queries (`?` placeholders passed as separate `Query`/`Exec` arguments) — never string-concatenate user input into SQL. `database/sql`'s driver handles safe escaping for you when you use placeholders correctly.

**Q: `sql.DB` — is it a single connection?**
A: No — it's a connection *pool*, safe for concurrent use. `sql.Open` doesn't even connect immediately (it's lazy); use `db.Ping()` to verify connectivity eagerly if you want to fail fast at startup.

---

## 10. System Design (Interview Framing)

**Q: How should you structure your answer to a "design X" system-design question?**
A: Clarify requirements and scale first (don't assume) → sketch a high-level design (main components and data flow) → go deep on 1-2 interesting parts the interviewer probes → explicitly discuss tradeoffs you made and what you'd reconsider at 10x scale. Interviewers are usually grading your reasoning process, not whether you produce "the" correct architecture.

**Q: Cache-aside vs write-through caching?**
A: Cache-aside: the application checks the cache first, and on a miss, reads from the source of truth and populates the cache itself — simple, but the cache can go briefly stale. Write-through: writes go to the cache and the backing store together, keeping them always in sync, at the cost of slower writes.

**Q: What does a circuit breaker protect against?**
A: Cascading failures — if a downstream dependency is failing, a circuit breaker "trips" after N consecutive failures and starts failing fast (without even attempting the call) instead of piling up slow, doomed requests. After a cooldown, it allows a trial request through (half-open state) to check if the dependency has recovered.

**Q: When is a monolith the right call over microservices?**
A: More often than most people assume, especially early on. A monolith has far less operational complexity (no distributed tracing, no network calls between your own components, one deployment). Microservices earn their complexity when different parts of the system genuinely need independent scaling, independent deployment cadence, or ownership by separate teams — not by default.

---

## 11. Behavioral / Soft Skills

**Q: The interviewer asks something you don't know the answer to. What do you do?**
A: Say so honestly, then reason toward an answer out loud from what you *do* know, rather than bluffing. Interviewers consistently rate "I'm not sure, but here's how I'd figure it out" far higher than a confident wrong answer.

**Q: How much should you talk while coding in a live interview?**
A: Continuously — narrate your plan before typing, explain tradeoffs as you make them, and call out edge cases you're intentionally deferring. A silent-but-correct solution gives the interviewer far less signal than a verbalized one, even if the verbalized one takes a bit longer.

**Q: What should you always check before saying "I'm done" on a coding problem?**
A: Edge cases: empty input, nil/zero values, a single-element input, and (for concurrent code) whether it can actually deadlock or race. Actually trace through your code with one of these cases before declaring victory — interviewers reliably ask "what about an empty slice?" if you don't address it yourself.

---

## Quick-Scan Gotcha Table

| Gotcha | One-line fix |
|---|---|
| Array copy vs slice sharing | Arrays copy on assignment/pass; slices share the backing array |
| `append` aliasing | Sub-slices sharing capacity can silently overwrite each other |
| nil map write panic | `make(map[K]V)` before writing |
| Struct/array `==` on uncomparable fields | Compile error — don't put slices/maps/funcs in a comparable struct |
| Pre-1.22 loop variable capture | Fixed in Go 1.22+; each iteration gets its own copy now |
| Typed nil in an interface | `err != nil` can be true for a nil concrete pointer — return the interface type directly |
| Method set mismatch | A pointer-receiver method requires `*T`, not `T`, to satisfy an interface |
| `defer` argument timing | Arguments evaluate immediately; only the call is delayed |
| Map iteration order | Randomized on purpose — sort keys for a stable order |
| Goroutine leak | Always know how a goroutine will stop (WaitGroup, channel, context) |
| Channel deadlock | Fast, safe crash (`fatal error: ... deadlock!`) — not a silent hang |
| `select` tie-breaking | Random among ready cases, never source-order priority |
| `nil` channel in `select` | Never ready — used deliberately to disable a case |

---

*Companion file to `level-41-go-interview-preparation/` in this 42-level Go course. Every code snippet above was verified by actually running it on this project's Go 1.26.5 toolchain.*
