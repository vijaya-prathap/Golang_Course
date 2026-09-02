# Level 11: Functions - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

---

## Exercise 1: Multiple Return Values

**Objective:** Write functions that return more than one value

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level11-exercise1
cd ~/projects/level11-exercise1
go mod init level11.example/exercise1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func divide(a, b int) (int, int) {
    return a / b, a % b
}

func rectangleStats(width, height float64) (float64, float64) {
    area := width * height
    perimeter := 2 * (width + height)
    return area, perimeter
}

func swap(a, b string) (string, string) {
    return b, a
}

func main() {
    fmt.Println("=== Multiple Return Values ===")
    quotient, remainder := divide(17, 5)
    fmt.Printf("17 / 5 = %d remainder %d\n", quotient, remainder)

    fmt.Println("\n=== Rectangle Stats ===")
    area, perimeter := rectangleStats(4.5, 2.0)
    fmt.Printf("Area: %.2f, Perimeter: %.2f\n", area, perimeter)

    fmt.Println("\n=== Swap ===")
    first, second := "Alice", "Bob"
    first, second = swap(first, second)
    fmt.Printf("After swap: first=%s, second=%s\n", first, second)

    fmt.Println("\n=== Ignoring a Return Value ===")
    q, _ := divide(100, 7)
    fmt.Println("quotient only:", q)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Multiple Return Values ===
17 / 5 = 3 remainder 2

=== Rectangle Stats ===
Area: 9.00, Perimeter: 13.00

=== Swap ===
After swap: first=Bob, second=Alice

=== Ignoring a Return Value ===
quotient only: 14
```

**Learning Objectives:**
- ✅ Declare a function with multiple return values
- ✅ Capture multiple values with multiple-value assignment
- ✅ Use `_` to ignore a return value you don't need

---

## Exercise 2: Named Return Values

**Objective:** Use named return values and naked returns

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level11-exercise2
cd ~/projects/level11-exercise2
go mod init level11.example/exercise2
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

// Named return values: area and circumference are declared in the
// signature, start at their zero value, and are filled in by the body.
func circleStats(radius float64) (area float64, circumference float64) {
    area = 3.14159 * radius * radius
    circumference = 2 * 3.14159 * radius
    return // naked return - sends back the current area, circumference
}

// Named returns are also useful for documenting what a function
// returns, even when you still use an explicit return.
func divideSafe(a, b int) (result int, err error) {
    if b == 0 {
        err = fmt.Errorf("cannot divide %d by zero", a)
        return
    }
    result = a / b
    return
}

func main() {
    fmt.Println("=== Named Return Values ===")
    area, circumference := circleStats(5)
    fmt.Printf("Radius 5 -> Area: %.2f, Circumference: %.2f\n", area, circumference)

    fmt.Println("\n=== Named Returns With Error Handling ===")
    result, err := divideSafe(10, 2)
    fmt.Printf("10 / 2 = %d, err = %v\n", result, err)

    result, err = divideSafe(10, 0)
    fmt.Printf("10 / 0 = %d, err = %v\n", result, err)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Named Return Values ===
Radius 5 -> Area: 78.54, Circumference: 31.42

=== Named Returns With Error Handling ===
10 / 2 = 5, err = <nil>
10 / 0 = 0, err = cannot divide 10 by zero
```

**Learning Objectives:**
- ✅ Declare named return values in a function signature
- ✅ Use a naked `return` to send back the current named values
- ✅ See named returns used for self-documenting `(result, err)` signatures

---

## Exercise 3: Variadic Functions

**Objective:** Accept a variable number of arguments, and spread a slice into one

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level11-exercise3
cd ~/projects/level11-exercise3
go mod init level11.example/exercise3
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

func average(nums ...int) float64 {
    if len(nums) == 0 {
        return 0
    }
    return float64(sum(nums...)) / float64(len(nums))
}

// Variadic parameters must come last, but ordinary parameters can come first.
func labelAndSum(label string, nums ...int) string {
    return fmt.Sprintf("%s: %d", label, sum(nums...))
}

func main() {
    fmt.Println("=== Variadic Basics ===")
    fmt.Println("sum():", sum())
    fmt.Println("sum(5):", sum(5))
    fmt.Println("sum(1, 2, 3, 4):", sum(1, 2, 3, 4))

    fmt.Println("\n=== Spreading a Slice ===")
    scores := []int{88, 92, 79, 95}
    fmt.Println("scores:", scores)
    fmt.Println("sum(scores...):", sum(scores...))
    fmt.Printf("average(scores...): %.2f\n", average(scores...))

    fmt.Println("\n=== Mixing Regular and Variadic Params ===")
    fmt.Println(labelAndSum("Total", scores...))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Variadic Basics ===
sum(): 0
sum(5): 5
sum(1, 2, 3, 4): 10

=== Spreading a Slice ===
scores: [88 92 79 95]
sum(scores...): 354
average(scores...): 88.50

=== Mixing Regular and Variadic Params ===
Total: 354
```

**Learning Objectives:**
- ✅ Declare a variadic parameter with `...`
- ✅ Call a variadic function with zero, one, or many arguments
- ✅ Spread an existing slice into a variadic call with `slice...`
- ✅ Combine a regular parameter with a trailing variadic one

---

## Exercise 4: Functions as Values (Higher-Order Functions)

**Objective:** Store functions in variables and pass them as arguments

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level11-exercise4
cd ~/projects/level11-exercise4
go mod init level11.example/exercise4
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

// apply is "map"-style: build a new slice by transforming every element.
func apply(nums []int, f func(int) int) []int {
    result := make([]int, len(nums))
    for i, n := range nums {
        result[i] = f(n)
    }
    return result
}

// filter keeps only elements where keep(n) is true.
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
func square(n int) int  { return n * n }
func isEven(n int) bool { return n%2 == 0 }
func isOdd(n int) bool  { return n%2 != 0 }

func main() {
    nums := []int{1, 2, 3, 4, 5, 6}

    fmt.Println("=== Assigning a Function to a Variable ===")
    var transform func(int) int = double
    fmt.Println("transform(21):", transform(21))

    fmt.Println("\n=== Passing Functions as Arguments (map-style) ===")
    fmt.Println("doubled:", apply(nums, double))
    fmt.Println("squared:", apply(nums, square))

    fmt.Println("\n=== Passing Functions as Arguments (filter-style) ===")
    fmt.Println("evens:", filter(nums, isEven))
    fmt.Println("odds:", filter(nums, isOdd))

    fmt.Println("\n=== Swapping Behavior at Runtime ===")
    operations := []func(int) int{double, square}
    for _, op := range operations {
        fmt.Println("op(4):", op(4))
    }
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Assigning a Function to a Variable ===
transform(21): 42

=== Passing Functions as Arguments (map-style) ===
doubled: [2 4 6 8 10 12]
squared: [1 4 9 16 25 36]

=== Passing Functions as Arguments (filter-style) ===
evens: [2 4 6]
odds: [1 3 5]

=== Swapping Behavior at Runtime ===
op(4): 8
op(4): 16
```

**Learning Objectives:**
- ✅ Understand that functions have types (`func(int) int`)
- ✅ Assign a function to a variable and call it indirectly
- ✅ Write higher-order functions (`apply`, `filter`) that accept a function parameter
- ✅ Store functions in a slice and call them in a loop

---

## Exercise 5: Anonymous Functions

**Objective:** Declare function literals inline - called immediately, stored, or passed directly

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level11-exercise5
cd ~/projects/level11-exercise5
go mod init level11.example/exercise5
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "sort"
)

func apply(nums []int, f func(int) int) []int {
    result := make([]int, len(nums))
    for i, n := range nums {
        result[i] = f(n)
    }
    return result
}

func main() {
    fmt.Println("=== Anonymous Function Called Immediately ===")
    result := func(a, b int) int {
        return a * b
    }(6, 7)
    fmt.Println("result:", result)

    fmt.Println("\n=== Anonymous Function Assigned to a Variable ===")
    greet := func(name string) string {
        return "Hi, " + name + "!"
    }
    fmt.Println(greet("Gopher"))
    fmt.Println(greet("World"))

    fmt.Println("\n=== Anonymous Function Passed Directly as an Argument ===")
    nums := []int{1, 2, 3, 4}
    plusTen := apply(nums, func(n int) int {
        return n + 10
    })
    fmt.Println("plusTen:", plusTen)

    fmt.Println("\n=== Real-World Use: sort.Slice with an Anonymous Function ===")
    names := []string{"Charlie", "Alice", "Bob"}
    sort.Slice(names, func(i, j int) bool {
        return names[i] < names[j]
    })
    fmt.Println("sorted names:", names)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Anonymous Function Called Immediately ===
result: 42

=== Anonymous Function Assigned to a Variable ===
Hi, Gopher!
Hi, World!

=== Anonymous Function Passed Directly as an Argument ===
plusTen: [11 12 13 14]

=== Real-World Use: sort.Slice with an Anonymous Function ===
sorted names: [Alice Bob Charlie]
```

**Learning Objectives:**
- ✅ Declare and immediately call an anonymous function
- ✅ Assign a function literal to a variable
- ✅ Pass a function literal directly as an argument
- ✅ Recognize `sort.Slice` as a real-world use of an inline comparison function

---

## Exercise 6: Closures - Shared Mutable State

**Objective:** Prove that a closure captures a variable by reference, sharing mutable state across calls

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level11-exercise6
cd ~/projects/level11-exercise6
go mod init level11.example/exercise6
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

// makeCounter closes over "count". Every call to the returned function
// shares and mutates the SAME "count" variable.
func makeCounter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

// makeAccount closes over "balance", giving us deposit/withdraw functions
// that share private, mutable state - a simple encapsulation pattern.
func makeAccount(initial float64) (deposit func(float64), withdraw func(float64) bool, balance func() float64) {
    bal := initial
    deposit = func(amount float64) {
        bal += amount
    }
    withdraw = func(amount float64) bool {
        if amount > bal {
            return false
        }
        bal -= amount
        return true
    }
    balance = func() float64 {
        return bal
    }
    return
}

func main() {
    fmt.Println("=== Counter Closure ===")
    counter := makeCounter()
    fmt.Println(counter()) // 1
    fmt.Println(counter()) // 2
    fmt.Println(counter()) // 3

    fmt.Println("\n=== Independent Closures Have Independent State ===")
    counterB := makeCounter()
    fmt.Println("counterB():", counterB()) // 1 - fresh state
    fmt.Println("counter():", counter())   // 4 - unaffected by counterB

    fmt.Println("\n=== Bank Account Closure ===")
    deposit, withdraw, balance := makeAccount(100)
    deposit(50)
    fmt.Println("balance after deposit(50):", balance())
    ok := withdraw(30)
    fmt.Println("withdraw(30) ok:", ok, "balance:", balance())
    ok = withdraw(1000)
    fmt.Println("withdraw(1000) ok:", ok, "balance:", balance())
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Counter Closure ===
1
2
3

=== Independent Closures Have Independent State ===
counterB(): 1
counter(): 4

=== Bank Account Closure ===
balance after deposit(50): 150
withdraw(30) ok: true balance: 120
withdraw(1000) ok: false balance: 120
```

**Learning Objectives:**
- ✅ Write a closure that captures and mutates a variable across calls
- ✅ Prove that each call to a closure-returning function creates independent state
- ✅ Use multiple closures sharing one private variable (`deposit`/`withdraw`/`balance` all share `bal`)

---

## Exercise 7: The Loop Variable Capture Gotcha (Verified on This Project's Go Version)

**Objective:** Verify, with real executed code, how closures capture loop variables on the Go version this course uses

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level11-exercise7
cd ~/projects/level11-exercise7
go mod init level11.example/exercise7
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== Loop Variable Capture (Go 1.22+ semantics) ===")
    // As of Go 1.22, "for" loops create a NEW instance of each loop
    // variable on every iteration. Each closure below captures its OWN
    // "i" - they do NOT all share one final value.
    var funcs []func()
    for i := 0; i < 3; i++ {
        funcs = append(funcs, func() {
            fmt.Println("captured i:", i)
        })
    }
    for _, f := range funcs {
        f()
    }

    fmt.Println("\n=== range Loop Variable Capture ===")
    var rangeFuncs []func()
    fruits := []string{"apple", "banana", "cherry"}
    for _, fruit := range fruits {
        rangeFuncs = append(rangeFuncs, func() {
            fmt.Println("captured fruit:", fruit)
        })
    }
    for _, f := range rangeFuncs {
        f()
    }

    fmt.Println("\n=== The Old Workaround Still Works (but is now unnecessary) ===")
    var oldStyleFuncs []func()
    for i := 0; i < 3; i++ {
        i := i // pre-1.22 idiom: shadow i with a fresh copy per iteration
        oldStyleFuncs = append(oldStyleFuncs, func() {
            fmt.Println("shadowed i:", i)
        })
    }
    for _, f := range oldStyleFuncs {
        f()
    }
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output** (verified against this project's `go.mod`, which specifies `go 1.26.5` - well past the Go 1.22 semantics change):

```
=== Loop Variable Capture (Go 1.22+ semantics) ===
captured i: 0
captured i: 1
captured i: 2

=== range Loop Variable Capture ===
captured fruit: apple
captured fruit: banana
captured fruit: cherry

=== The Old Workaround Still Works (but is now unnecessary) ===
shadowed i: 0
shadowed i: 1
shadowed i: 2
```

**Important:** if you read older Go material (pre-2024, before Go 1.22), it will tell you the first block should print `3` three times, because every closure supposedly shared one `i` that ended at `3`. That was true on Go 1.21 and earlier. It is **not** true here - each iteration gets its own `i`, and all three sections above print 0, 1, 2 (or apple/banana/cherry) exactly once each, in order.

**Learning Objectives:**
- ✅ Verify, by actually running the code, how loop variables are captured by closures on Go 1.26.5
- ✅ Understand why older tutorials warn about a bug that no longer applies
- ✅ Recognize that the pre-1.22 `i := i` shadowing workaround is now redundant, though harmless

---

## Exercise 8: Recursion - Factorial and Fibonacci

**Objective:** Write recursive functions, and compare naive recursion to iteration

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level11-exercise8
cd ~/projects/level11-exercise8
go mod init level11.example/exercise8
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n-1)
}

// fibNaive recomputes the same sub-problems over and over.
// Its runtime grows exponentially with n (roughly O(2^n)).
func fibNaive(n int) int {
    if n < 2 {
        return n
    }
    return fibNaive(n-1) + fibNaive(n-2)
}

// fibIterative solves the same problem in O(n) time and O(1) space -
// no call stack growth, no repeated work.
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

func main() {
    fmt.Println("=== Factorial (Recursive) ===")
    for i := 0; i <= 6; i++ {
        fmt.Printf("factorial(%d) = %d\n", i, factorial(i))
    }

    fmt.Println("\n=== Fibonacci: Naive Recursion vs Iterative ===")
    for i := 0; i <= 10; i++ {
        naive := fibNaive(i)
        iterative := fibIterative(i)
        match := naive == iterative
        fmt.Printf("fib(%2d): naive=%-4d iterative=%-4d match=%v\n", i, naive, iterative, match)
    }
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Factorial (Recursive) ===
factorial(0) = 1
factorial(1) = 1
factorial(2) = 2
factorial(3) = 6
factorial(4) = 24
factorial(5) = 120
factorial(6) = 720

=== Fibonacci: Naive Recursion vs Iterative ===
fib( 0): naive=0    iterative=0    match=true
fib( 1): naive=1    iterative=1    match=true
fib( 2): naive=1    iterative=1    match=true
fib( 3): naive=2    iterative=2    match=true
fib( 4): naive=3    iterative=3    match=true
fib( 5): naive=5    iterative=5    match=true
fib( 6): naive=8    iterative=8    match=true
fib( 7): naive=13   iterative=13   match=true
fib( 8): naive=21   iterative=21   match=true
fib( 9): naive=34   iterative=34   match=true
fib(10): naive=55   iterative=55   match=true
```

**Learning Objectives:**
- ✅ Write a recursive function with a base case and a recursive case
- ✅ Confirm naive recursive Fibonacci and iterative Fibonacci agree on every value
- ✅ Understand why naive recursive Fibonacci is exponential (recomputes overlapping sub-problems), while the iterative version is linear

---

## Exercise 9: defer - LIFO Order and Immediate Argument Evaluation

**Objective:** Verify defer's execution order and argument evaluation timing with real output

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level11-exercise9
cd ~/projects/level11-exercise9
go mod init level11.example/exercise9
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func lifoDemo() {
    fmt.Println("start of lifoDemo")
    defer fmt.Println("deferred #1 (registered first, runs LAST)")
    defer fmt.Println("deferred #2")
    defer fmt.Println("deferred #3 (registered last, runs FIRST)")
    fmt.Println("end of lifoDemo")
}

func argEvalDemo() {
    x := 10
    // x's value is captured (as 10) the moment this defer statement runs.
    // Changing x afterward does not change what gets printed.
    defer fmt.Println("deferred saw x =", x)
    x = 99
    fmt.Println("current x =", x)
}

func multipleArgEvalDemo() {
    for i := 1; i <= 3; i++ {
        // Each iteration's "i" is evaluated immediately and baked into
        // its own deferred call - they'll still print in LIFO order
        // when multipleArgEvalDemo returns.
        defer fmt.Println("deferred with i =", i)
    }
    fmt.Println("loop finished")
}

func main() {
    fmt.Println("=== defer Runs in LIFO Order ===")
    lifoDemo()

    fmt.Println("\n=== defer Arguments Are Evaluated Immediately ===")
    argEvalDemo()

    fmt.Println("\n=== Combining LIFO Order With Immediate Argument Evaluation ===")
    multipleArgEvalDemo()
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== defer Runs in LIFO Order ===
start of lifoDemo
end of lifoDemo
deferred #3 (registered last, runs FIRST)
deferred #2
deferred #1 (registered first, runs LAST)

=== defer Arguments Are Evaluated Immediately ===
current x = 99
deferred saw x = 10

=== Combining LIFO Order With Immediate Argument Evaluation ===
loop finished
deferred with i = 3
deferred with i = 2
deferred with i = 1
```

**Learning Objectives:**
- ✅ Confirm multiple `defer` statements run in Last-In-First-Out order
- ✅ Confirm a deferred call's arguments are evaluated when `defer` runs, not when the call finally executes
- ✅ See both rules combine inside a loop that defers once per iteration

---

## Exercise 10: Comprehensive Practice - Function Dispatcher Calculator

**Objective:** Combine multiple return values, named returns, and functions-as-values in one realistic program

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level11-exercise10
cd ~/projects/level11-exercise10
go mod init level11.example/exercise10
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "os"
)

func add(a, b float64) float64 { return a + b }
func sub(a, b float64) float64 { return a - b }
func mul(a, b float64) float64 { return a * b }

func div(a, b float64) (result float64, err error) {
    if b == 0 {
        err = fmt.Errorf("division by zero")
        return
    }
    result = a / b
    return
}

// calculate is a dispatcher: it looks up an operation by name in a map
// of functions and calls it. Adding a new operation means adding one
// map entry - no new if/else branches.
func calculate(op string, a, b float64) (float64, error) {
    operations := map[string]func(float64, float64) float64{
        "+": add,
        "-": sub,
        "*": mul,
    }

    if op == "/" {
        return div(a, b)
    }

    fn, ok := operations[op]
    if !ok {
        return 0, fmt.Errorf("unknown operator: %q", op)
    }
    return fn(a, b), nil
}

func runCalculation(op string, a, b float64) {
    result, err := calculate(op, a, b)
    if err != nil {
        fmt.Fprintf(os.Stdout, "%.1f %s %.1f -> error: %v\n", a, op, b, err)
        return
    }
    fmt.Printf("%.1f %s %.1f = %.2f\n", a, op, b, result)
}

func main() {
    fmt.Println("=== Function Dispatcher Calculator ===")
    runCalculation("+", 4, 5)
    runCalculation("-", 10, 3)
    runCalculation("*", 6, 7)
    runCalculation("/", 20, 4)
    runCalculation("/", 5, 0)
    runCalculation("%", 5, 2)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Function Dispatcher Calculator ===
4.0 + 5.0 = 9.00
10.0 - 3.0 = 7.00
6.0 * 7.0 = 42.00
20.0 / 4.0 = 5.00
5.0 / 0.0 -> error: division by zero
5.0 % 2.0 -> error: unknown operator: "%"
```

**Learning Objectives:**
- ✅ Combine multiple return values, named returns, and a `map[string]func(...)` dispatcher in one program
- ✅ Use a map of functions to avoid a long if/else chain of operator checks
- ✅ Handle both a "not found" error and a domain error (division by zero) through the same `(result, error)` pattern

---

## Bonus Challenges

### Challenge 1: reduce Helper

Write a generic-free `reduce` function that folds a slice of `int` down to a single value using an accumulator function and a starting value - the classic "reduce"/"fold" pattern that `apply` (map) and `filter` are usually taught alongside.

```bash
mkdir -p ~/projects/level11-bonus1
cd ~/projects/level11-bonus1
go mod init level11.example/bonus1
```

**Hints:**
- Signature: `func reduce(nums []int, initial int, f func(acc, cur int) int) int`
- Loop over `nums`, updating `acc := f(acc, cur)` each time
- Try it with a sum function, a product function, and a max function

**Verified reference solution and output:**

```go
package main

import "fmt"

func reduce(nums []int, initial int, f func(acc, cur int) int) int {
    acc := initial
    for _, n := range nums {
        acc = f(acc, n)
    }
    return acc
}

func main() {
    nums := []int{1, 2, 3, 4, 5}

    sum := reduce(nums, 0, func(acc, cur int) int {
        return acc + cur
    })
    fmt.Println("sum:", sum)

    product := reduce(nums, 1, func(acc, cur int) int {
        return acc * cur
    })
    fmt.Println("product:", product)

    max := reduce(nums, nums[0], func(acc, cur int) int {
        if cur > acc {
            return cur
        }
        return acc
    })
    fmt.Println("max:", max)
}
```

Output:

```
sum: 15
product: 120
max: 5
```

### Challenge 2: Memoized Fibonacci Using a Closure

Write a `makeMemoFib() func(int) int` that returns a Fibonacci function closing over a cache (`map[int]int`), so repeated calls to the *same* returned function reuse previously computed results instead of recomputing them.

```bash
mkdir -p ~/projects/level11-bonus2
cd ~/projects/level11-bonus2
go mod init level11.example/bonus2
```

**Hints:**
- Declare the recursive closure with `var fib func(int) int` first, so it can call itself by name
- Check the cache before recursing; store the result in the cache before returning
- Try `memoFib(40)` - it returns instantly, unlike `fibNaive(40)` from Exercise 8

**Verified reference solution and output:**

```go
package main

import "fmt"

func makeMemoFib() func(int) int {
    cache := make(map[int]int)
    var fib func(int) int
    fib = func(n int) int {
        if n < 2 {
            return n
        }
        if v, ok := cache[n]; ok {
            return v
        }
        result := fib(n-1) + fib(n-2)
        cache[n] = result
        return result
    }
    return fib
}

func main() {
    memoFib := makeMemoFib()
    for i := 0; i <= 15; i++ {
        fmt.Printf("fib(%d) = %d\n", i, memoFib(i))
    }

    fmt.Println("\nA larger value that would be painfully slow with naive recursion:")
    fmt.Println("fib(40) =", memoFib(40))
}
```

Output:

```
fib(0) = 0
fib(1) = 1
fib(2) = 1
fib(3) = 2
fib(4) = 3
fib(5) = 5
fib(6) = 8
fib(7) = 13
fib(8) = 21
fib(9) = 34
fib(10) = 55
fib(11) = 89
fib(12) = 144
fib(13) = 233
fib(14) = 377
fib(15) = 610

A larger value that would be painfully slow with naive recursion:
fib(40) = 102334155
```

### Challenge 3: Retry With Backoff Helper

Write a `retry(maxAttempts int, fn func(attempt int) error) error` helper that calls `fn`, retrying (with a simulated, doubling backoff delay) until `fn` succeeds or `maxAttempts` is reached.

```bash
mkdir -p ~/projects/level11-bonus3
cd ~/projects/level11-bonus3
go mod init level11.example/bonus3
```

**Hints:**
- `fn` takes the current attempt number and returns an `error` (`nil` means success)
- Track and double a `delay` counter each time an attempt fails, just to demonstrate the pattern (no real `time.Sleep` needed for this exercise)
- Return `nil` as soon as `fn` succeeds; otherwise return a wrapped error after the last attempt using `fmt.Errorf("...: %w", lastErr)`

**Verified reference solution and output:**

```go
package main

import (
    "errors"
    "fmt"
)

func retry(maxAttempts int, fn func(attempt int) error) error {
    var lastErr error
    delay := 1
    for attempt := 1; attempt <= maxAttempts; attempt++ {
        err := fn(attempt)
        if err == nil {
            fmt.Printf("attempt %d succeeded\n", attempt)
            return nil
        }
        lastErr = err
        fmt.Printf("attempt %d failed: %v\n", attempt, err)
        if attempt < maxAttempts {
            fmt.Printf("  backing off for %d unit(s) before retrying\n", delay)
            delay *= 2
        }
    }
    return fmt.Errorf("all %d attempts failed: %w", maxAttempts, lastErr)
}

func main() {
    fmt.Println("=== Retry Succeeds on the 3rd Attempt ===")
    err := retry(5, func(attempt int) error {
        if attempt < 3 {
            return errors.New("simulated transient failure")
        }
        return nil
    })
    fmt.Println("final result:", err)

    fmt.Println("\n=== Retry Exhausts All Attempts ===")
    err = retry(3, func(attempt int) error {
        return errors.New("simulated permanent failure")
    })
    fmt.Println("final result:", err)
}
```

Output:

```
=== Retry Succeeds on the 3rd Attempt ===
attempt 1 failed: simulated transient failure
  backing off for 1 unit(s) before retrying
attempt 2 failed: simulated transient failure
  backing off for 2 unit(s) before retrying
attempt 3 succeeded
final result: <nil>

=== Retry Exhausts All Attempts ===
attempt 1 failed: simulated permanent failure
  backing off for 1 unit(s) before retrying
attempt 2 failed: simulated permanent failure
  backing off for 2 unit(s) before retrying
attempt 3 failed: simulated permanent failure
final result: all 3 attempts failed: simulated permanent failure
```

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Write functions with multiple return values and multiple-value assignment
✅ Use named return values and naked returns, and know when naked returns hurt readability
✅ Write variadic functions and spread a slice into one with `...`
✅ Treat functions as values - store them, pass them, and write higher-order functions
✅ Write and use anonymous functions (function literals) inline
✅ Write closures that share and mutate state across calls
✅ Explain (with verified, real output) how this Go version captures loop variables in closures
✅ Write recursive functions, and know when to prefer iteration or memoization
✅ Predict `defer`'s LIFO ordering and its immediate argument evaluation
✅ Combine everything into a realistic function-dispatcher program

---

## Next Level

Level 12: Pointers
- What a pointer is and why Go has them
- The `&` (address-of) and `*` (dereference) operators
- Pointer receivers vs. value receivers
- When to pass a pointer instead of a value

Great work! You've finished the FUNDAMENTALS tier of the course - keep going! 🚀
