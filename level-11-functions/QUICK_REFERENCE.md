# Level 11: Quick Reference Card

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir myapp && cd myapp
go mod init myapp

# Create main.go
cat > main.go << 'EOF'
package main

import "fmt"

func divmod(a, b int) (int, int) {
    return a / b, a % b
}

func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

func main() {
    q, r := divmod(17, 5)
    fmt.Println(q, r)

    fmt.Println(sum(1, 2, 3, 4))

    values := []int{10, 20, 30}
    fmt.Println(sum(values...))
}
EOF

# Run
go run main.go
```

---

## 📋 Function Declaration

```go
func name(param type, param type) returnType {
    return value
}

func add(a, b int) int { return a + b }   // shared-type shorthand
func log(msg string)   { fmt.Println(msg) } // no return value
```

---

## 🔁 Multiple Return Values

```go
func divmod(a, b int) (int, int) {
    return a / b, a % b
}

q, r := divmod(17, 5)   // capture both
q, _ := divmod(17, 5)   // ignore remainder
```

---

## 🏷️ Named Returns & Naked Returns

```go
func minMax(nums []int) (min, max int) {
    min, max = nums[0], nums[0]
    for _, n := range nums {
        if n < min { min = n }
        if n > max { max = n }
    }
    return // naked return - sends back current min, max
}
```

Use naked returns only in short functions - prefer explicit `return value, err` once a function has multiple exit points.

---

## 🎒 Variadic Functions

```go
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

sum(1, 2, 3)          // ordinary call
values := []int{1, 2, 3}
sum(values...)        // spread a slice into the call
```

---

## 🧩 Functions as Values

```go
var fn func(int) int = double     // store a function in a variable

func apply(nums []int, f func(int) int) []int {
    result := make([]int, len(nums))
    for i, n := range nums {
        result[i] = f(n)
    }
    return result
}

apply(nums, double)   // pass a function as an argument
```

---

## 🎭 Anonymous Functions

```go
result := func(a, b int) int { return a + b }(3, 4)  // called immediately

square := func(n int) int { return n * n }           // assigned to a variable

apply(nums, func(n int) int { return n * 3 })        // passed inline
```

---

## 🔗 Closures

```go
func makeCounter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

counter := makeCounter()
counter()  // 1
counter()  // 2 - same "count", shared across calls
```

---

## 🔂 Loop Variable Capture (This Course: go 1.26.5, Go 1.22+ semantics)

```go
var funcs []func()
for i := 0; i < 3; i++ {
    funcs = append(funcs, func() { fmt.Println(i) })
}
for _, f := range funcs {
    f()
}
// Prints: 0, 1, 2  (each iteration gets its own "i")
// Pre-Go-1.22 this printed 3, 3, 3 - not on this course's Go version
```

---

## 🌀 Recursion

```go
func factorial(n int) int {
    if n <= 1 {
        return 1 // base case
    }
    return n * factorial(n-1) // recursive case
}
```

Prefer iteration/memoization when sub-problems overlap (e.g. naive Fibonacci is O(2ⁿ); iterative is O(n)).

---

## ⏱️ defer

```go
defer fmt.Println("runs when the function returns")

defer fmt.Println("1st deferred, runs LAST")
defer fmt.Println("2nd deferred, runs FIRST")   // LIFO order

i := 0
defer fmt.Println(i) // argument evaluated NOW (prints 0), not later
i++
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Variadic not last | `func f(nums ...int, s string)` | `func f(s string, nums ...int)` |
| Missing spread | `sum(values)` | `sum(values...)` |
| Deferred closure vs. deferred value | assuming `defer fmt.Println(i)` sees later changes | it captures `i`'s value immediately |
| Old loop-capture assumption | expecting `3 3 3` | this course (go 1.26.5) prints `0 1 2` |
| Naive recursion at scale | `fibNaive(40)` | `fibIterative(40)` or memoize |

---

## 🎓 Before Next Level

Can you:
- [ ] Write a function with multiple return values and unpack them?
- [ ] Explain when a naked return helps vs. hurts readability?
- [ ] Write a variadic function and spread a slice into it?
- [ ] Pass a function as an argument to another function?
- [ ] Write a closure and explain why its state persists across calls?
- [ ] Correctly explain how THIS Go version captures loop variables in closures?
- [ ] Write a recursive function with a clear base case?
- [ ] Predict `defer`'s LIFO order and argument evaluation timing?

If YES → You're ready for Level 12!

---

## 📚 Next Level

Level 12: Pointers
- What a pointer is and why Go has them
- The `&` (address-of) and `*` (dereference) operators
- Pointer receivers vs. value receivers

You've got functions down - and finished FUNDAMENTALS! 💪
