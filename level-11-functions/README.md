# Level 11: Functions - Complete Guide

## Introduction

Welcome to Level 11! You've mastered maps (Level 10) - Go's key-value collection type. Now it's time to make **functions** the explicit subject, instead of the background tool you've been using since Level 3.

You've already been writing functions without much ceremony: `divmod` in Level 3 returned two values, and `isEven`, `isPrime`, `grade`, and `calculate` all showed up as helpers in Levels 4-6. This level formalizes everything about functions in Go - multiple return values, named returns, variadic parameters, functions as first-class values, closures, recursion, and `defer` - and goes well beyond what you've seen incidentally so far.

Functions are the last topic in the **FUNDAMENTALS** tier of this course (Levels 0-11). Once you're comfortable here, Level 12 begins the **OOP CONCEPTS** tier (Pointers, Structs, Methods, Interfaces, Error Handling) - and nearly everything in that tier is built on top of the function skills you'll practice here.

---

## Table of Contents

1. [Function Declaration Syntax](#function-declaration-syntax)
2. [Multiple Return Values](#multiple-return-values)
3. [Named Return Values and Naked Returns](#named-return-values-and-naked-returns)
4. [Variadic Functions](#variadic-functions)
5. [Functions as Values](#functions-as-values)
6. [Anonymous Functions (Function Literals)](#anonymous-functions-function-literals)
7. [Closures](#closures)
8. [Closures in Loops: The Loop Variable Gotcha](#closures-in-loops-the-loop-variable-gotcha)
9. [Recursion](#recursion)
10. [defer](#defer)
11. [Best Practices](#best-practices)
12. [Common Mistakes](#common-mistakes)

---

## Function Declaration Syntax

Every Go function follows the same basic shape:

```go
func name(parameterName parameterType, ...) returnType {
    // body
    return value
}
```

### Anatomy

```go
func add(a int, b int) int {
    return a + b
}

func greet(name string) string {
    return "Hello, " + name + "!"
}

func main() {
    sum := add(3, 4)
    fmt.Println(sum) // 7
    fmt.Println(greet("Gopher")) // Hello, Gopher!
}
```

```
func <name>(<param> <type>, <param> <type>) <returnType> {
    ├─ name:       identifier used to call the function
    ├─ parameters: zero or more name+type pairs (comma separated)
    ├─ returnType: the type of value the function sends back
    └─ body:       the statements that run when the function is called
}
```

### Parameters of the Same Type

When consecutive parameters share a type, you can drop the type on all but the last one:

```go
// These two declarations are identical:
func add(a int, b int) int { return a + b }
func add(a, b int) int     { return a + b }
```

### No Return Value

A function that doesn't return anything simply omits the return type:

```go
func logMessage(msg string) {
    fmt.Println("[LOG]", msg)
    // no return statement needed
}
```

---

## Multiple Return Values

Unlike many languages, Go functions can return more than one value directly - no wrapping them in a struct, tuple, or array required. You already saw this with `divmod` back in Level 3; here it is formalized.

```go
func divmod(a, b int) (int, int) {
    return a / b, a % b
}

func minMax(nums []int) (int, int) {
    min, max := nums[0], nums[0]
    for _, n := range nums {
        if n < min {
            min = n
        }
        if n > max {
            max = n
        }
    }
    return min, max
}

func main() {
    q, r := divmod(17, 5)
    fmt.Printf("17 / 5 = %d remainder %d\n", q, r)
    // 17 / 5 = 3 remainder 2

    min, max := minMax([]int{4, 8, 15, 16, 23, 42})
    fmt.Println("min:", min, "max:", max)
    // min: 4 max: 42

    // Ignoring one of the return values with _
    quotient, _ := divmod(20, 3)
    fmt.Println("quotient only:", quotient)
    // quotient only: 6
}
```

### Multiple-Value Assignment

Any time a function returns multiple values, you must either capture all of them (using `_` for the ones you don't need) or capture none of them (calling the function only for its side effects, if any):

```go
result, err := someFunc()  // capture both
_, err := someFunc()       // capture only the error
someFunc()                 // capture neither (only valid if you don't need the values)
```

The most common use of multiple return values in idiomatic Go is `(value, error)` - a pattern you've already seen with `strconv.Atoi` in Level 3, and one that Level 16 (Error Handling) will cover in depth.

---

## Named Return Values and Naked Returns

You can name the return values in a function's signature. Named return values are declared like parameters, start at their type's zero value, and can be returned with a bare `return` statement - called a **naked return**.

```go
func minMaxNamed(nums []int) (min int, max int) {
    min, max = nums[0], nums[0]
    for _, n := range nums {
        if n < min {
            min = n
        }
        if n > max {
            max = n
        }
    }
    return // naked return - sends back current min, max
}

func main() {
    min, max := minMaxNamed([]int{9, 3, 7, 1, 5})
    fmt.Println("min:", min, "max:", max)
    // min: 1 max: 9
}
```

Named returns are also useful purely as **documentation** - the signature tells the reader what each returned value means, even if you still write an explicit `return`:

```go
func splitFullName(full string) (first string, last string) {
    for i := 0; i < len(full); i++ {
        if full[i] == ' ' {
            return full[:i], full[i+1:]
        }
    }
    return full, ""
}
```

### Best Practice: Don't Overuse Naked Returns

Naked returns work fine in short functions, but they hurt readability once a function gets longer - by the time you reach `return`, the reader has to scroll back up to remember what each named value currently holds:

```go
// ⚠️ Naked return in a longer function - unclear at a glance what's returned
func processOrder(id int, quantity int, price float64) (total float64, err error) {
    if quantity <= 0 {
        err = fmt.Errorf("invalid quantity: %d", quantity)
        return
    }
    if price < 0 {
        err = fmt.Errorf("invalid price: %.2f", price)
        return
    }
    total = float64(quantity) * price
    if total > 10000 {
        total = total * 0.9 // bulk discount
    }
    return // <- what is being returned here? you have to scroll up to check
}

// ✅ Same logic with explicit returns - each exit point is self-contained
func processOrderExplicit(id int, quantity int, price float64) (float64, error) {
    if quantity <= 0 {
        return 0, fmt.Errorf("invalid quantity: %d", quantity)
    }
    if price < 0 {
        return 0, fmt.Errorf("invalid price: %.2f", price)
    }
    total := float64(quantity) * price
    if total > 10000 {
        total = total * 0.9
    }
    return total, nil
}
```

**Rule of thumb:** named returns are great for documentation and for short functions (a handful of lines). Once a function has multiple exit points or grows past 10-15 lines, prefer explicit `return value, err` statements at each exit point.

---

## Variadic Functions

A variadic parameter accepts zero or more arguments of the same type. Mark it with `...` before the type, and it becomes an ordinary slice inside the function body.

```go
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

func main() {
    fmt.Println(sum())           // 0  - no arguments at all
    fmt.Println(sum(5))          // 5  - one argument
    fmt.Println(sum(1, 2, 3, 4)) // 10 - several arguments

    values := []int{10, 20, 30}
    fmt.Println(sum(values...)) // 60 - spread a slice into the variadic call
}
```

### Spreading a Slice With `...`

If you already have a slice and want to pass its elements as individual variadic arguments, append `...` to the slice at the call site:

```go
values := []int{10, 20, 30}
sum(values...)   // equivalent to sum(10, 20, 30)
```

Without the trailing `...`, `sum(values)` would be a compile error - `values` is a `[]int`, not an `int`.

### Rules

- A variadic parameter must be the **last** parameter in the list.
- A function can have regular parameters before the variadic one: `func labelAndSum(label string, nums ...int) string`.
- A function can only have **one** variadic parameter.
- Inside the function, the variadic parameter behaves exactly like a slice - `len()`, `range`, indexing, all work normally.

---

## Functions as Values

In Go, functions are **first-class values**: you can store them in variables, pass them as arguments, and return them from other functions. A function's type is written as `func(paramTypes) returnType`.

```go
func apply(nums []int, f func(int) int) []int {
    result := make([]int, len(nums))
    for i, n := range nums {
        result[i] = f(n)
    }
    return result
}

func filter(nums []int, keep func(int) bool) []int {
    var result []int
    for _, n := range nums {
        if keep(n) {
            result = append(result, n)
        }
    }
    return result
}

func double(n int) int  { return n * 2 }
func isEven(n int) bool { return n%2 == 0 }

func main() {
    nums := []int{1, 2, 3, 4, 5}

    // Assigning a function to a variable - fn has type func(int) int
    var fn func(int) int = double
    fmt.Println("fn(10):", fn(10)) // 20

    doubled := apply(nums, double)
    fmt.Println("doubled:", doubled) // [2 4 6 8 10]

    evens := filter(nums, isEven)
    fmt.Println("evens:", evens) // [2 4]
}
```

A function that takes another function as a parameter (like `apply` and `filter` above) - or returns one - is called a **higher-order function**. `apply` is a "map"-style helper (transform every element), and `filter` keeps only the elements that pass a test. You'll build both, plus a `reduce` helper, in the exercises.

---

## Anonymous Functions (Function Literals)

A **function literal** is a function without a name. You can call it immediately, assign it to a variable, or pass it directly as an argument - useful for small, one-off pieces of logic that don't deserve a top-level name.

```go
func main() {
    // Declared and called immediately (IIFE-style)
    result := func(a, b int) int {
        return a + b
    }(3, 4)
    fmt.Println("immediate call result:", result) // 7

    // Assigned to a variable, called later
    square := func(n int) int {
        return n * n
    }
    fmt.Println("square(5):", square(5)) // 25

    // Passed directly as an argument (no separate named function needed)
    nums := []int{1, 2, 3}
    tripled := apply(nums, func(n int) int {
        return n * 3
    })
    fmt.Println("tripled:", tripled) // [3 6 9]
}
```

Anonymous functions passed inline are extremely common in idiomatic Go - `sort.Slice(names, func(i, j int) bool { ... })` is a pattern you'll see constantly once you start using the standard library more.

---

## Closures

A **closure** is a function literal that references variables from outside its own body. The function "closes over" those variables - it keeps a live reference to them, not a snapshot, so it can both read and mutate them even after the enclosing function has returned.

```go
// makeCounter returns a function literal that closes over "count".
// Each call to the returned function increments the SAME "count"
// variable - it is not reset between calls.
func makeCounter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

func main() {
    counter := makeCounter()
    fmt.Println(counter()) // 1
    fmt.Println(counter()) // 2
    fmt.Println(counter()) // 3

    // A second counter has its own independent "count"
    counter2 := makeCounter()
    fmt.Println(counter2()) // 1
    fmt.Println(counter())  // 4 - counter's state is unaffected by counter2
}
```

**Verified output** (from `go run`):

```
1
2
3
1
4
```

This proves two things at once: `count` is genuinely **shared, mutable state** across calls to the same closure (1, 2, 3, 4 - it keeps climbing), and each call to `makeCounter()` creates a **brand-new, independent** `count` (`counter2` starts back at 1 without disturbing `counter`).

### Why This Matters

A closure captures the *variable*, not the *value it held at the time of capture*. This is the mental model you need for the loop-variable gotcha in the next section - and it's also the foundation for patterns like the memoized-fibonacci and retry-with-backoff exercises later in this level.

---

## Closures in Loops: The Loop Variable Gotcha

This is one of the most famous gotchas in Go's history - and it's **no longer a gotcha** in the version of Go this course uses. Read this section carefully, because a lot of blog posts, Stack Overflow answers, and older books will confidently tell you the *old* behavior as if it still applies.

### The Historical Problem (Go 1.21 and Earlier)

Before Go 1.22, a `for` loop declared its index/range variables **once**, and every iteration reused the same variable. A closure created inside the loop body captured that one shared variable - so if you saved the closures and called them all *after* the loop finished, they would all see whatever value the variable held at the end of the loop, not the value it had "at capture time." This tripped up huge numbers of Go programmers and produced a well-known idiom to work around it: shadow the loop variable with `i := i` inside the loop body to give each iteration its own copy.

### The Current Behavior (Go 1.22+)

**This project's `go.mod` specifies `go 1.26.5`** (confirmed by reading `/Users/macbookpro/go/Golang/go.mod`), which is far past the Go 1.22 release that changed loop variable semantics. As of Go 1.22, `for` loops create a **new instance of each loop variable on every iteration**. Each closure captures its own, independent variable - the old bug simply cannot happen anymore.

```go
func main() {
    var funcs []func()
    for i := 0; i < 3; i++ {
        funcs = append(funcs, func() {
            fmt.Println(i)
        })
    }
    for _, f := range funcs {
        f()
    }
}
```

**Verified output on this project's Go version (go1.26.5):**

```
0
1
2
```

Each closure printed the value `i` held *during its own iteration* - `0`, `1`, `2` - not three copies of the final value `3`. On Go 1.21 or earlier, this exact program would have printed `3` three times.

### Why You Might Still See Warnings About This

If you read a Go blog post, book, or Stack Overflow answer from before early 2024 (when Go 1.22 shipped), it will very likely warn you about "the loop variable capture bug" and recommend the `i := i` shadowing workaround. That advice was correct for the Go version it was written against - but on Go 1.22 and later (including the `go 1.26.5` this course targets), it's unnecessary. The workaround is still **harmless** if you see it in old code (shadowing a variable with an identical copy of itself is a no-op under the new semantics), so you don't need to "fix" it if you encounter it - just know you don't need to *add* it in new code.

```go
// Still works, but is now redundant on Go 1.22+:
for i := 0; i < 3; i++ {
    i := i // pre-1.22 idiom: no longer changes behavior on this Go version
    funcs = append(funcs, func() { fmt.Println(i) })
}
```

This same fixed semantics applies to `range` loops too - a closure capturing a `range` variable (index or value) also gets its own copy per iteration on Go 1.22+.

---

## Recursion

A **recursive function** calls itself, working toward a **base case** that stops the recursion.

```go
func factorial(n int) int {
    if n <= 1 {
        return 1 // base case
    }
    return n * factorial(n-1) // recursive case
}
```

Every recursive function needs both parts: a base case that returns without recursing, and a recursive case that makes progress toward the base case. Miss the base case (or never make progress toward it) and you get infinite recursion - which in Go ends in a **stack overflow** crash, not an infinite loop that silently hangs.

### Fibonacci: Naive Recursion vs. Iteration

A classic example of recursion that "works" but scales badly:

```go
func fibNaive(n int) int {
    if n < 2 {
        return n
    }
    return fibNaive(n-1) + fibNaive(n-2)
}

func fibIterative(n int) int {
    if n < 2 {
        return n
    }
    a, b := 0, 1
    for i := 2; i <= n; i++ {
        a, b = b, a+b
    }
    return b
}
```

`fibNaive` recomputes the same sub-problems over and over - `fibNaive(5)` calls `fibNaive(4)` and `fibNaive(3)`, but `fibNaive(4)` *also* calls `fibNaive(3)` internally, duplicating work. Its running time grows **exponentially** with `n` (roughly O(2ⁿ)). `fibIterative` solves the identical problem in a single pass with O(n) time and no extra stack frames.

```
fibNaive(0..9):     0 1 1 2 3 5 8 13 21 34
fibIterative(0..9): 0 1 1 2 3 5 8 13 21 34
```

Both produce identical results - the difference is purely how much work it takes to get there. You'll see this gap widen dramatically at `fib(40)` in the exercises (naive recursion becomes noticeably slow; an iterative or memoized version stays instant).

### Stack Depth Awareness

Every recursive call adds a new frame to the call stack. Go's goroutine stacks grow dynamically (starting small and expanding as needed), which is much more forgiving than languages with a fixed stack size - but a recursive function with no base case, or one that recurses far too deep (e.g., processing a list with a million elements one recursive call at a time), can still exhaust available memory and crash with `runtime: goroutine stack exceeds ... - fatal error: stack overflow`.

### When Recursion Is (and Isn't) Idiomatic in Go

Recursion is a natural fit when the *data* is naturally recursive - trees, nested JSON/graph structures, divide-and-conquer algorithms (like quicksort or binary search on a tree). It's less idiomatic in Go for simple linear iteration over a fixed range or a slice, where a plain `for` loop (Level 6) is clearer, avoids stack growth, and is usually faster. When a recursive solution's naive form has overlapping sub-problems (like `fibNaive`), prefer an iterative version or add memoization (caching previously computed results) - you'll build a memoized Fibonacci using a closure in the bonus challenges.

---

## defer

`defer` schedules a function call to run **after** the surrounding function returns - useful for cleanup work (closing a file, unlocking a mutex, closing a database connection) that needs to happen no matter which path the function takes to return. This level introduces `defer` as a function-level concept; Level 18 (File Handling) and Level 16 (Error Handling) will deepen the cleanup patterns built on top of it.

```go
func greet() {
    defer fmt.Println("2: deferred - runs when greet() returns")
    fmt.Println("1: runs immediately")
}

func main() {
    greet()
    fmt.Println("3: back in main")
}
```

**Verified output:**

```
1: runs immediately
2: deferred - runs when greet() returns
3: back in main
```

The deferred call is registered the moment `defer` runs, but its execution is delayed until `greet()` is about to return - after every other statement in `greet()` has already run.

### Multiple defers Run in LIFO Order

If a function has more than one `defer` statement, they run in **Last In, First Out** order - like a stack. The most recently deferred call runs first.

```go
func main() {
    fmt.Println("start")
    defer fmt.Println("deferred 1 (runs LAST)")
    defer fmt.Println("deferred 2")
    defer fmt.Println("deferred 3 (runs FIRST of the defers)")
    fmt.Println("end")
}
```

**Verified output:**

```
start
end
deferred 3 (runs FIRST of the defers)
deferred 2
deferred 1 (runs LAST)
```

### defer Arguments Are Evaluated Immediately

This is a precise and frequently-misremembered rule: when a `defer` statement runs, **the arguments to the deferred call are evaluated right then** - only the *call itself* is postponed. Changing a variable after the `defer` statement does not change what was already captured.

```go
func main() {
    i := 0
    // The argument to fmt.Println (the value of i) is evaluated RIGHT NOW,
    // when the defer statement runs - only the CALL itself is delayed.
    defer fmt.Println("deferred: i was", i)
    i++
    fmt.Println("current: i is now", i)
}
```

**Verified output:**

```
current: i is now 1
deferred: i was 0
```

Even though `i` became `1` before `main()` returned, the deferred call already had `0` locked in as its argument at the moment `defer` executed. (If the deferred call were instead a closure like `defer func() { fmt.Println(i) }()`, it *would* see the latest value of `i` when it finally runs, because then it's the variable being captured, not a value being passed as an argument - the same capture-by-reference behavior you saw with closures above.)

---

## Best Practices

### 1. Keep Functions Small and Focused

```go
// ✅ Good - one clear responsibility
func isValidEmail(email string) bool {
    return strings.Contains(email, "@")
}

// ❌ Doing too much - validating, formatting, AND logging
func processEmail(email string) string {
    if !strings.Contains(email, "@") {
        log.Println("invalid email")
    }
    return strings.ToLower(email)
}
```

### 2. Prefer Explicit Returns Over Naked Returns in Longer Functions

Naked returns are fine for a handful of lines; beyond that, explicit `return value, err` at each exit point is easier to audit (see the [naked returns section](#named-return-values-and-naked-returns) above).

### 3. Put the Variadic Parameter Last, and Only Use One

```go
// ✅ Good
func logEvent(level string, args ...interface{}) { }

// ❌ Won't compile - variadic must be last, and only one is allowed
func broken(a ...int, b ...string) { }
```

### 4. Use Higher-Order Functions to Avoid Repeating Loop Logic

```go
// ✅ Good - the looping logic lives in one place (apply/filter)
evens := filter(nums, isEven)
doubled := apply(nums, double)

// ❌ Repeats the same loop shape everywhere it's needed
var evens []int
for _, n := range nums {
    if n%2 == 0 {
        evens = append(evens, n)
    }
}
```

### 5. Don't Assume Old Loop-Closure Warnings Still Apply

On this course's Go version (`go 1.26.5`, Go 1.22+ semantics), each loop iteration already gets its own copy of the loop variable. Don't reflexively add `i := i` shadowing out of habit from older tutorials - it's harmless, but unnecessary noise.

### 6. Use defer Right After Acquiring a Resource

```go
// ✅ Good - defer sits right next to the acquisition, impossible to forget
file, err := os.Open("data.txt")
if err != nil {
    return err
}
defer file.Close()

// ❌ Risky - easy to add an early return between Open and Close and leak the file
file, err := os.Open("data.txt")
if err != nil {
    return err
}
// ... lots of code ...
file.Close() // if an early return happens above, this never runs
```

### 7. Prefer Iteration (or Memoization) Over Naive Recursion for Overlapping Sub-Problems

If a recursive function recomputes the same inputs repeatedly (like `fibNaive`), switch to an iterative version or cache results with memoization - both keep the same result while cutting the work from exponential to linear.

---

## Common Mistakes

### Mistake 1: Forgetting a Variadic Parameter Must Be Last

```go
// ❌ WRONG - compile error, variadic parameter isn't last
func broken(nums ...int, label string) { }

// ✅ RIGHT
func fixed(label string, nums ...int) { }
```

### Mistake 2: Passing a Slice to a Variadic Parameter Without `...`

```go
values := []int{1, 2, 3}

// ❌ WRONG - compile error: cannot use values (type []int) as type int
sum(values)

// ✅ RIGHT - spread the slice
sum(values...)
```

### Mistake 3: Assuming a Deferred Closure's Argument Is Evaluated Immediately

```go
// ❌ SURPRISING - i is captured BY THE CLOSURE, not evaluated until the
// closure actually runs, so this prints the FINAL value of i, not each one
i := 0
defer func() {
    fmt.Println(i) // whatever i is when this finally executes
}()
i = 99
// prints 99, not 0!

// ✅ If you want the value AT THE defer STATEMENT, pass it as an argument
i := 0
defer func(captured int) {
    fmt.Println(captured)
}(i)
i = 99
// prints 0 - the argument was evaluated when defer ran
```

### Mistake 4: Expecting the Pre-1.22 Loop Variable Sharing Bug

```go
// On Go 1.22+ (this project uses go 1.26.5), this prints 0, 1, 2 -
// NOT three copies of 3 like it would have on Go 1.21 and earlier.
var funcs []func()
for i := 0; i < 3; i++ {
    funcs = append(funcs, func() { fmt.Println(i) })
}
for _, f := range funcs {
    f()
}
```

### Mistake 5: Naive Recursion on a Problem With Overlapping Sub-Problems

```go
// ❌ WRONG (not really wrong, but painfully slow) - exponential time
func fibNaive(n int) int {
    if n < 2 {
        return n
    }
    return fibNaive(n-1) + fibNaive(n-2)
}
fibNaive(40) // takes a noticeable amount of time

// ✅ RIGHT - iterative, or memoized with a closure over a cache
func fibIterative(n int) int {
    if n < 2 {
        return n
    }
    a, b := 0, 1
    for i := 2; i <= n; i++ {
        a, b = b, a+b
    }
    return b
}
fibIterative(40) // instant
```

### Mistake 6: Forgetting That Named Returns Start at Their Zero Value

```go
// ❌ WRONG ASSUMPTION - if the "if" is never true, err is nil (its zero
// value), not some "unset" state - so make sure every path sets what it needs
func validate(age int) (err error) {
    if age < 0 {
        err = fmt.Errorf("age cannot be negative")
    }
    return // if age >= 0, this returns err == nil, which is CORRECT here,
           // but it's easy to forget a case and silently return a zero value
}
```

### Mistake 7: Overusing Naked Returns in Long Functions

See the [naked returns best-practice note](#named-return-values-and-naked-returns) - a naked `return` several branches and a dozen lines away from the named return declaration forces the reader to scroll back up to know what's being sent back.

---

## Summary

**Function Basics:**
- `func name(params) returnType { }` - the universal shape of a Go function
- Functions can return multiple values directly - no wrapping needed
- Named return values pre-declare and zero-initialize the results; `return` alone (a naked return) sends back their current values

**Variadic and First-Class Functions:**
- `func f(nums ...int)` accepts zero or more arguments, available as a slice inside the function
- Spread a slice into a variadic call with `slice...`
- Functions are values: they have types (`func(int) int`), can be stored in variables, passed as arguments, and returned from other functions

**Anonymous Functions and Closures:**
- A function literal can be called immediately, assigned to a variable, or passed inline
- A closure captures the enclosing scope's variables by reference, not by value - shared, mutable state survives across calls
- On this course's Go version (1.22+ semantics, `go 1.26.5`), each loop iteration gets its own copy of the loop variable, so closures created in a loop capture distinct values - the historic "loop variable capture bug" is fixed

**Recursion and defer:**
- A recursive function needs a base case and a recursive case that makes progress toward it
- Naive recursion on overlapping sub-problems (like Fibonacci) is exponential; prefer iteration or memoization
- `defer` delays a *call*, not its *argument evaluation* - arguments are captured immediately
- Multiple `defer` statements in one function run in LIFO order

---

## Next Steps

You now understand:
- ✅ Function declaration, parameters, and return values
- ✅ Multiple return values and named returns (including naked-return trade-offs)
- ✅ Variadic functions and spreading slices into them
- ✅ Functions as values and higher-order functions (map/filter-style helpers)
- ✅ Anonymous functions and closures, including the fixed Go 1.22+ loop-variable semantics
- ✅ Recursion, its costs, and when to prefer iteration
- ✅ `defer`, its LIFO ordering, and its immediate argument evaluation

**Next level:** Level 12 - Pointers
- What a pointer is and why Go has them
- The `&` (address-of) and `*` (dereference) operators
- Pointer receivers vs. value receivers (foundation for Level 14: Methods)
- When to pass a pointer instead of a value

This closes out the FUNDAMENTALS tier of the course - congratulations! Level 12 begins OOP CONCEPTS, where structs, methods, and interfaces build directly on the functions you just mastered. Keep going! 🚀
