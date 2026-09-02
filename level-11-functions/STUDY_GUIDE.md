# Level 11: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: Functions
```
Day 1:  Function declaration syntax, multiple return values
Day 2:  Named return values and naked returns
Day 3:  Variadic functions
Day 4:  Functions as values, higher-order functions
Day 5:  Anonymous functions and closures
Day 6:  The loop variable capture gotcha (old vs current Go)
Day 7:  Recursion and defer
```

### Week 2: Practice & Application
```
Day 1:  Exercises 1-3 (multiple returns, named returns, variadic)
Day 2:  Exercises 4-6 (functions as values, anonymous functions, closures)
Day 3:  Exercises 7-8 (loop variable capture, recursion)
Day 4:  Exercises 9-10 (defer, comprehensive dispatcher)
Day 5:  Bonus challenges
Day 6-7: Practice & review
```

---

## 🎯 Function Declaration Anatomy

```
func   name   (  param  type,  param  type  )   returnType   {
 │      │        │       │      │      │            │          │
 │      │        │       │      │      │            │          └─ body starts here
 │      │        │       │      │      │            └─ what the function sends back
 │      │        │       │      │      └─ second parameter's type
 │      │        │       │      └─ second parameter's name
 │      │        │       └─ first parameter's type
 │      │        └─ first parameter's name
 │      └─ the function's identifier (how you call it)
 └─ keyword that starts every function declaration
```

```go
func add(a int, b int) int {
    return a + b
}

// Shared-type shorthand - identical to the above
func add(a, b int) int {
    return a + b
}

// No return value - return type is simply omitted
func logMessage(msg string) {
    fmt.Println("[LOG]", msg)
}
```

---

## 🔁 Multiple Return Values

| Pattern | Example | Notes |
|---------|---------|-------|
| Two values | `func divmod(a, b int) (int, int)` | Types listed in parentheses, comma separated |
| Capture both | `q, r := divmod(17, 5)` | Standard multi-value assignment |
| Ignore one | `q, _ := divmod(17, 5)` | `_` discards a value you don't need |
| `(value, error)` | `n, err := strconv.Atoi("42")` | The most common multi-return shape in Go |

```
divmod(17, 5)
      │
      ├─→ 3   (quotient:  17 / 5)
      └─→ 2   (remainder: 17 % 5)

q, r := divmod(17, 5)   // q = 3, r = 2
```

---

## 🏷️ Named Returns & Naked Returns

```go
func minMaxNamed(nums []int) (min int, max int) {
//                             └──────┬──────┘
//                    declared here, zero-valued at function start
    min, max = nums[0], nums[0]
    for _, n := range nums {
        if n < min { min = n }
        if n > max { max = n }
    }
    return          // ← "naked return": sends back CURRENT min, max
}
```

### Decision Tree: Naked Return or Explicit Return?

```
Is the function short (roughly < 10-15 lines) with ONE exit point?
├─ YES → naked "return" is fine, and documents intent via the signature
└─ NO  → prefer explicit "return value, err" at each exit point
          (a naked return several lines/branches away forces the
           reader to scroll back up to remember what's returned)
```

---

## 🎒 Variadic Functions

```
func sum(nums ...int) int
              │
              └─ inside the function body, nums has type []int

sum()              → nums == []int{}          → 0
sum(5)             → nums == []int{5}         → 5
sum(1, 2, 3)       → nums == []int{1, 2, 3}   → 6

values := []int{10, 20, 30}
sum(values...)     ← spread operator "..." unpacks the slice
                      into individual variadic arguments
sum(values)         ✗ compile error - []int is not int
```

**Rule:** a variadic parameter must be the LAST parameter, and a function may have only ONE.

---

## 🧩 Functions as Values & Higher-Order Functions

```
Function TYPE:     func(int) int
Function VALUE:    double            (a specific func(int) int)

var fn func(int) int = double     // store a function in a variable
fn(10)                            // call it indirectly → 20

apply(nums, double)   // pass a function AS AN ARGUMENT ("higher-order")
                      //   apply is "higher-order" because it takes
                      //   a func(int) int parameter
```

| Helper | Shape | Purpose |
|--------|-------|---------|
| `apply` (map-style) | `func(nums []int, f func(int) int) []int` | Transform every element |
| `filter` | `func(nums []int, keep func(int) bool) []int` | Keep only elements passing a test |
| `reduce` (bonus) | `func(nums []int, initial int, f func(acc, cur int) int) int` | Fold a slice down to one value |

---

## 🔗 Closures Capture Variables by Reference

```
func makeCounter() func() int {
    count := 0              ┐
    return func() int {     │  the returned function CLOSES OVER
        count++             │  "count" - it keeps a live link to
        return count        │  the SAME variable, not a copy
    }                       ┘
}

counter := makeCounter()   // count = 0, owned by this closure instance
counter()   → 1             // count becomes 1
counter()   → 2             // count becomes 2
counter()   → 3             // count becomes 3

counter2 := makeCounter()  // a NEW, independent count = 0
counter2()  → 1             // counter2's own count becomes 1
counter()   → 4             // counter's count is UNAFFECTED - still climbing
```

```
┌─────────────────────┐        ┌─────────────────────┐
│ makeCounter() call 1 │        │ makeCounter() call 2 │
│   count = 4 (shared) │        │   count = 1 (shared) │
└─────────┬────────────┘        └─────────┬────────────┘
          │ closure reads/writes           │ closure reads/writes
          ▼ the SAME variable               ▼ its OWN variable
      counter()                          counter2()
```

**Key idea:** a closure captures the *variable itself*, not a snapshot of its value at the moment the closure was created. This is exactly the mental model needed for the loop-variable topic below.

---

## 🔂 Loop Variable Capture: Before vs. After Go 1.22

This project's `go.mod` specifies **`go 1.26.5`** - well past the Go 1.22 semantics change. The table below is accurate for **this course's Go version**.

| | Go 1.21 and earlier | Go 1.22 and later (this course: go 1.26.5) |
|---|---|---|
| Loop variable lifetime | ONE variable, reused every iteration | A NEW variable created each iteration |
| What a closure captures | The single shared variable | Its own iteration's variable |
| Output of the classic example | `3 3 3` (all closures see the final value) | `0 1 2` (each closure sees its own value) |
| Common workaround | `i := i` shadow inside the loop body | Not needed (harmless if present) |

```go
var funcs []func()
for i := 0; i < 3; i++ {
    funcs = append(funcs, func() { fmt.Println(i) })
}
for _, f := range funcs {
    f()
}
```

```
Go 1.21 and earlier:          Go 1.22+ (verified on go1.26.5 for this course):
┌───────────────────┐         ┌───────────────────┐
│ one shared "i"     │         │ i=0 (iteration 0)  │──→ closure 0 → prints 0
│ ends at 3          │         │ i=1 (iteration 1)  │──→ closure 1 → prints 1
│ ALL closures read  │         │ i=2 (iteration 2)  │──→ closure 2 → prints 2
│ the SAME final "i" │         └───────────────────┘
└───────────────────┘
  prints: 3 3 3                  prints: 0 1 2
```

**Verified output on this project's Go version:**

```
0
1
2
```

If you read older tutorials warning about "the loop variable bug," they were correct for their time (pre-2024 / pre-Go-1.22) - just outdated for the Go version this course uses.

---

## 🌀 Recursion: Call Stack Diagram

```go
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n-1)
}
```

`factorial(3)` unwinding:

```
factorial(3)
 └─ calls factorial(2)
     └─ calls factorial(1)
         └─ base case: return 1
     └─ factorial(2) resumes: return 2 * 1 = 2
 └─ factorial(3) resumes: return 3 * 2 = 6

Call stack grows DOWN as calls are made, then unwinds UP as each
call returns:

  factorial(3)  ─┐
  factorial(2)   │  stack grows
  factorial(1)  ─┘  (deepest point: 3 frames)
  ── base case hit, returns 1 ──
  factorial(2) resumes → 2 * 1 = 2
  factorial(3) resumes → 3 * 2 = 6
```

### Naive vs. Iterative Fibonacci: Why the Difference Matters

```
fibNaive(5)
├─ fibNaive(4)
│  ├─ fibNaive(3)          ← computed here...
│  │  ├─ fibNaive(2)
│  │  └─ fibNaive(1)
│  └─ fibNaive(2)
└─ fibNaive(3)             ← ...and AGAIN here (duplicated work!)
   ├─ fibNaive(2)
   └─ fibNaive(1)

fibNaive(n)     → O(2^n)  time  (exponential - recomputes sub-problems)
fibIterative(n) → O(n)    time  (linear - one pass, no repeated work)
```

---

## ⏱️ defer: LIFO Order Diagram

```go
func main() {
    defer fmt.Println("deferred 1 (runs LAST)")
    defer fmt.Println("deferred 2")
    defer fmt.Println("deferred 3 (runs FIRST of the defers)")
    fmt.Println("end")
}
```

```
Registration order (top to bottom):     Execution order (LIFO - a stack):
┌─────────────────┐                     ┌─────────────────┐
│ defer 1          │  pushed 1st         │ defer 3 runs 1st │  popped last-in-first-out
│ defer 2          │  pushed 2nd         │ defer 2 runs 2nd │
│ defer 3          │  pushed 3rd         │ defer 1 runs 3rd │
└─────────────────┘                     └─────────────────┘

Verified output:
end
deferred 3 (runs FIRST of the defers)
deferred 2
deferred 1 (runs LAST)
```

### defer Argument Evaluation Timing

```go
i := 0
defer fmt.Println("deferred: i was", i)   // "i" (== 0) is READ RIGHT NOW
i++
fmt.Println("current: i is now", i)
```

```
Timeline:
 t0: defer statement runs → argument "i" evaluated NOW → 0 is locked in
 t1: i++                  → i becomes 1 (too late, already captured 0)
 t2: fmt.Println runs     → prints "current: i is now 1"
 t3: main() about to return → deferred call FINALLY executes
                            → prints "deferred: i was 0"

Verified output:
current: i is now 1
deferred: i was 0
```

**Contrast:** a deferred *closure* (`defer func() { fmt.Println(i) }()`) captures the variable `i` itself, so it WOULD see the latest value when it finally runs - only a directly deferred call's arguments are locked in immediately.

---

## 🚨 Common Mistakes

### Mistake 1: Variadic Parameter Not Last

```go
// ❌ WRONG - compile error
func broken(nums ...int, label string) { }

// ✅ RIGHT
func fixed(label string, nums ...int) { }
```

### Mistake 2: Forgetting `...` When Spreading a Slice

```go
values := []int{1, 2, 3}
sum(values)    // ❌ compile error - []int is not int
sum(values...) // ✅ spreads the slice
```

### Mistake 3: Confusing a Deferred Closure With a Deferred Call's Arguments

```go
i := 0
defer func() { fmt.Println(i) }()  // captures the VARIABLE → prints 99
i = 99

i2 := 0
defer fmt.Println(i2)              // captures the VALUE (0) right now → prints 0
i2 = 99
```

### Mistake 4: Expecting the Pre-1.22 Loop-Sharing Behavior

```go
// On go 1.26.5 (this course), prints 0, 1, 2 - NOT three copies of 3.
for i := 0; i < 3; i++ {
    funcs = append(funcs, func() { fmt.Println(i) })
}
```

### Mistake 5: Naive Recursion on Overlapping Sub-Problems

```go
fibNaive(40)      // ❌ noticeably slow - exponential time
fibIterative(40)  // ✅ instant - linear time
```

---

## 📈 Progression Summary

### Understanding Level 11

Level 11 makes functions the explicit subject after 10 levels of using them incidentally:

1. **Declaration** - the universal `func name(params) returnType { }` shape
2. **Multiple returns** - returning more than one value without a wrapper type
3. **Named returns** - self-documenting signatures and naked returns
4. **Variadic** - accepting a variable number of arguments
5. **First-class functions** - storing, passing, and returning functions
6. **Closures** - capturing enclosing-scope variables by reference
7. **Loop variable capture** - verified, current (Go 1.22+) semantics
8. **Recursion** - base cases, recursive cases, and when to prefer iteration
9. **defer** - LIFO ordering and immediate argument evaluation

### Prerequisites for Level 12

Before moving to Level 12 (Pointers), you need:

- ✅ Comfortable writing functions with multiple return values
- ✅ Comfortable with variadic functions and spreading slices with `...`
- ✅ Can write and explain a closure that shares mutable state across calls
- ✅ Understand how this Go version (1.22+) captures loop variables in closures
- ✅ Can write a recursive function and identify its base/recursive cases
- ✅ Can predict `defer`'s LIFO order and its immediate argument evaluation

### Ready for Level 12?

Level 12 teaches the mechanism that lets functions share and modify data across calls without closures:
- What a pointer is, and the `&`/`*` operators
- Passing a pointer instead of a value (and why)
- The foundation for pointer receivers in Level 14 (Methods)

---

## ✅ Checklist Before Level 12

- [ ] Can declare a function with multiple return values and unpack them
- [ ] Can explain when a naked return helps vs. hurts readability
- [ ] Can write a variadic function and spread a slice into it
- [ ] Can assign a function to a variable and pass one as an argument
- [ ] Can write a closure (like `makeCounter`) and explain why its state persists
- [ ] Can explain - correctly, for THIS Go version - how closures capture loop variables
- [ ] Can write a recursive function and know when to prefer an iterative version
- [ ] Can predict the output of a function using multiple `defer` statements
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The Shape
Every function is `func name(params) returnType { }` - parameters, multiple returns, and no return type are all variations on this one shape.

### The Capture
A closure captures variables by reference - shared, mutable, and alive as long as something can still call the closure.

### The Version-Specific Fact
On Go 1.22+ (this course: `go 1.26.5`), each loop iteration gets its own copy of the loop variable - closures created in a loop no longer share one final value.

### The Order
`defer` runs LIFO; a deferred call's arguments are locked in the moment `defer` executes, not when the call finally runs.

---

## 📚 Next Level

Level 12: Pointers
- What a pointer is and why Go has them
- The `&` (address-of) and `*` (dereference) operators
- Pointer receivers vs. value receivers
- When to pass a pointer instead of a value

You've mastered functions - and finished the FUNDAMENTALS tier! Keep going! 🚀
