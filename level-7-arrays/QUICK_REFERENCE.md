# Level 7: Quick Reference Card

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir myapp && cd myapp
go mod init myapp

# Create main.go
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    var scores [5]int
    scores[0] = 90
    scores[1] = 85

    colors := [3]string{"Red", "Green", "Blue"}
    primes := [...]int{2, 3, 5, 7, 11}

    fmt.Println(scores)
    fmt.Println(colors)
    fmt.Println(primes, "len:", len(primes))
}
EOF

# Run
go run main.go
```

---

## 📋 Declaration Forms

```go
var arr [5]int                    // zero-valued: [0 0 0 0 0]
lit := [3]string{"a", "b", "c"}   // explicit size, must match count
inf := [...]int{1, 2, 3}          // size inferred: [3]int
sparse := [5]int{0: 1, 4: 5}      // [1 0 0 0 5], rest are zero
var g [4]float64                  // declare, then assign by index
g[0] = 91.5
```

---

## 🔑 The One Rule That Matters

```go
// ARRAYS COPY. SLICES SHARE.

a := [3]int{1, 2, 3}
b := a          // b is an INDEPENDENT copy
b[0] = 99       // a is untouched: a == [1 2 3]

sa := []int{1, 2, 3}
sb := sa        // sb shares sa's underlying array
sb[0] = 99      // sa IS changed: sa == [99 2 3]
```

Same rule for function calls:

```go
func f(arr [3]int) { arr[0] = 0 }   // arr is a copy - caller unaffected
func f(s []int)    { s[0] = 0 }     // s shares data - caller's slice IS affected
```

---

## ⚖️ Comparing Arrays

```go
a := [3]int{1, 2, 3}
b := [3]int{1, 2, 3}
a == b                       // true - same type, same elements

d := [4]int{1, 2, 3, 4}
// a == d                    // compile error: mismatched types [3]int and [4]int

seen := map[[2]int]string{}  // arrays CAN be map keys (slices can't)
seen[[2]int{0, 0}] = "origin"
```

---

## 🧊 Multi-Dimensional Arrays

```go
var grid [3][3]int
grid[1][2] = 7            // grid[row][col]

board := [2][3]string{
    {"X", "O", "X"},
    {"O", "X", "O"},
}

len(grid)      // 3 - number of rows
len(grid[0])   // 3 - number of columns

for r, row := range grid {
    for c, v := range row {
        fmt.Println(r, c, v)
    }
}
```

---

## 🎁 range Over Arrays

```go
for i, v := range arr { }   // index + value (both COPIES)
for _, v := range arr { }   // value only
for i := range arr { }      // index only

// range copies each element - mutating v does NOT change arr:
for _, v := range arr {
    v *= 100   // no-op on arr
}
// to actually mutate, index back in:
for i := range arr {
    arr[i] *= 100
}
```

---

## 🧮 Array ↔ Slice Conversion

```go
arr := [5]int{10, 20, 30, 40, 50}
s := arr[:]        // whole array as a slice - SHARES memory with arr
p := arr[1:3]      // partial slice - also SHARES memory

s[0] = 999
fmt.Println(arr)   // [999 20 30 40 50] - arr changed too!
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Expecting reference semantics | function mutates `[N]T` param, expects caller to see it | pass a slice, or `*[N]T`, if mutation must be visible |
| Mixing array sizes | `var b [4]int = a` where `a` is `[3]int` | sizes must match exactly - use a slice if size varies |
| Confusing `[...]` with dynamic size | `nums := [...]int{1,2,3}; append(nums, 4)` | arrays have no `append` - use `[]int` for growable data |
| Forgetting range copies | `for _, v := range arr { v *= 2 }` (no effect) | `for i := range arr { arr[i] *= 2 }` |
| Assuming bounds checks are always at runtime | expects `arr[3]` (constant, out of range) to panic | constant out-of-range index is a COMPILE error; only a variable index panics at runtime |

---

## 🎓 Before Next Level

Can you:
- [ ] Declare an array five different ways?
- [ ] Prove with real output that arrays copy on assignment and function calls?
- [ ] Explain why slices behave differently with the exact same syntax?
- [ ] Compare two arrays with == and explain the type/size rule?
- [ ] Build and index a 2D array?
- [ ] Explain when to actually reach for an array instead of a slice?

If YES → You're ready for Level 8!

---

## 📚 Next Level

Level 8: Slices
- Slices in depth: pointer/length/capacity
- append, growth, and re-slicing
- Why slices share memory

You've got arrays down! 💪
