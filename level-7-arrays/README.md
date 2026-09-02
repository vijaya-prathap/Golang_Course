# Level 7: Arrays - Complete Guide

## Introduction

Welcome to Level 7! You've mastered loops (Level 6) and can now iterate over any collection. It's time to look closely at the simplest collection type in Go: the **array**.

You've actually already glimpsed arrays back in Level 3's composite types preview and used `range` over them in Level 6. This level goes deep: how arrays are declared, why their **fixed size** is part of their type, and - most importantly - why arrays behave completely differently from slices when you assign or pass them around.

That last point is the single most important idea in this level. Get comfortable with it now, because it explains a lot of "weird" Go behavior once you get to Level 8 (Slices).

---

## Table of Contents

1. [What Arrays Are](#what-arrays-are)
2. [Declaring Arrays](#declaring-arrays)
3. [Indexing, len(), and Zero Values](#indexing-len-and-zero-values)
4. [Arrays Are Value Types](#arrays-are-value-types)
5. [Comparing Arrays with ==](#comparing-arrays-with-)
6. [Multi-Dimensional Arrays](#multi-dimensional-arrays)
7. [Iterating Arrays with range](#iterating-arrays-with-range)
8. [Array vs Slice: When to Use Which](#array-vs-slice-when-to-use-which)
9. [Best Practices](#best-practices)
10. [Common Mistakes](#common-mistakes)

---

## What Arrays Are

An array is a **fixed-size**, **homogeneous** collection: a set number of elements, all of the same type, stored contiguously in memory.

```go
var scores [5]int // exactly 5 ints, always 5 ints
```

Two words matter enormously here:

- **Fixed-size** - the length is part of the array's *type*. `[5]int` and `[10]int` are two different types, the same way `int` and `string` are different types. You cannot grow or shrink an array.
- **Homogeneous** - every element must be the same type. `[5]int` holds only `int`s.

This is fundamentally different from a **slice** (which you'll meet properly in Level 8), a dynamically-sized view over an array that can grow with `append`. Every slice is backed by an array under the hood, but a slice itself is *not* an array - it's a small header (pointer, length, capacity) pointing at one.

**Why does Go even have arrays if slices are more flexible?** Because "fixed-size" is sometimes exactly the guarantee you want: a chess board is always 8x8, an IPv4 address is always 4 bytes, a SHA-256 hash is always 32 bytes. Arrays let the type system enforce that size for you. They're also the foundation slices are built on - understanding arrays makes slices make sense.

---

## Declaring Arrays

There are several ways to declare an array, and they all end up producing a value of type `[N]T`.

### 1. var With a Fixed Size (Zero-Valued)

```go
var arr [5]int
fmt.Println(arr) // [0 0 0 0 0]
```

Every element starts at that type's zero value (see [Section 3](#indexing-len-and-zero-values)).

### 2. Array Literal With an Explicit Size

```go
colors := [3]string{"a", "b", "c"}
fmt.Println(colors) // [a b c]
```

The number of elements in the literal must exactly match the declared size, or the compiler rejects it.

### 3. Array Literal With Size Inference: `[...]T`

Let Go count the elements for you:

```go
primes := [...]int{2, 3, 5, 7, 11}
fmt.Println(primes)        // [2 3 5 7 11]
fmt.Println(len(primes))   // 5
```

`[...]int` is not "any size" - the size is fixed to exactly the number of elements you list, and it's determined once, at compile time. `primes` here has the concrete type `[5]int`, identical to writing `[5]int{2, 3, 5, 7, 11}` by hand.

### 4. Sparse Index Literals

You can set only specific indices and let everything else default to the zero value:

```go
sparse := [5]int{0: 1, 4: 5}
fmt.Println(sparse) // [1 0 0 0 5]
```

Index `0` gets `1`, index `4` gets `5`, and indices `1`, `2`, `3` fall back to the zero value `0`.

This also combines with `[...]`, where the size is inferred from the *highest* index used, not the number of elements written:

```go
sparseInferred := [...]int{2: 100, 5: 200}
fmt.Println(sparseInferred)      // [0 0 100 0 0 200]
fmt.Println(len(sparseInferred)) // 6 (highest index 5, so length is 6)
```

### 5. Declare Then Assign

```go
var grades [4]float64
grades[0] = 91.5
grades[3] = 88.0
fmt.Println(grades) // [91.5 0 0 88]
```

---

## Indexing, len(), and Zero Values

Arrays index like you'd expect, starting at `0`:

```go
fruits := [3]string{"apple", "banana", "cherry"}
fmt.Println(fruits[0]) // apple
fmt.Println(fruits[1]) // banana
fmt.Println(fruits[2]) // cherry
```

`len()` always returns the array's fixed size:

```go
fmt.Println(len(fruits)) // 3
```

Elements are assignable through their index:

```go
fruits[1] = "blueberry"
fmt.Println(fruits) // [apple blueberry cherry]
```

### Zero Values

An array's zero value is an array of the same length, with every element set to *that type's* zero value:

```go
var ints [3]int       // [0 0 0]
var floats [3]float64 // [0 0 0]
var strs [3]string    // [  ]   (three empty strings)
var bools [3]bool     // [false false false]
```

Unlike slices and maps, an array's zero value is immediately usable - there's no `nil` array. `var arr [5]int` is a complete, ready-to-use array of five zeros, not a "not yet allocated" placeholder.

### Out-of-Range Access

Go checks array bounds. If the index is a **compile-time constant**, an out-of-range access is caught while compiling:

```go
arr := [3]int{1, 2, 3}
fmt.Println(arr[3]) // compile error:
// invalid argument: index 3 out of bounds [0:3]
```

If the index is only known at **runtime** (e.g., it comes from a variable), the same mistake instead panics while the program is running:

```go
i := 3
fmt.Println(arr[i])
// panic: runtime error: index out of range [3] with length 3
```

Either way, Go never silently reads or writes past the end of an array.

---

## Arrays Are Value Types

This is the most important - and most commonly misunderstood - fact about arrays in Go: **an array is a value, not a reference**. Copying an array copies every element. There is no shared backing storage the way there is with slices, maps, or pointers.

### Assignment Copies the Whole Array

```go
a := [3]int{1, 2, 3}
b := a    // COPY - b gets its own independent [3]int
b[0] = 99

fmt.Println(a) // [1 2 3]   <- unchanged!
fmt.Println(b) // [99 2 3]
```

Contrast that with a slice, where assignment copies the slice *header* (pointer + length + capacity), not the underlying data - so both variables still point at the same array:

```go
sa := []int{1, 2, 3}
sb := sa  // sb shares sa's underlying array
sb[0] = 99

fmt.Println(sa) // [99 2 3]  <- changed too!
fmt.Println(sb) // [99 2 3]
```

Same syntax (`b := a`), opposite behavior - because one is an array (value type) and the other is a slice (reference-like type).

### Passing an Array to a Function Passes a COPY

The exact same rule applies to function calls, since passing an argument is really just "assigning it to the parameter." A function that receives an array parameter gets its own private copy - mutating it inside the function has **zero effect** on the caller's array.

```go
func zeroOutArray(arr [3]int) {
    for i := range arr {
        arr[i] = 0
    }
    fmt.Println("  inside:", arr) // [0 0 0]
}

func main() {
    arr := [3]int{1, 2, 3}
    zeroOutArray(arr)
    fmt.Println(arr) // [1 2 3]  <- unchanged! the function only zeroed its own copy
}
```

Now compare a slice parameter, which shares the underlying array with the caller:

```go
func zeroOutSlice(s []int) {
    for i := range s {
        s[i] = 0
    }
}

func main() {
    s := []int{1, 2, 3}
    zeroOutSlice(s)
    fmt.Println(s) // [0 0 0]  <- changed! slices share the underlying array
}
```

**This is the number one reason Go code almost always uses slices instead of arrays as function parameters.** If you ever find yourself passing large arrays around expecting a function to modify them in place, that's the value-semantics gotcha catching you - either return the modified array, use a slice, or pass a pointer to the array (`*[3]int`).

---

## Comparing Arrays with ==

Because arrays are values (not references), and because Go can compare them element by element, arrays of **comparable element type and identical size** support `==` and `!=` directly - something slices cannot do at all.

```go
a := [3]int{1, 2, 3}
b := [3]int{1, 2, 3}
c := [3]int{1, 2, 99}

fmt.Println(a == b) // true  - same length, same type, same elements
fmt.Println(a == c) // false - element at index 2 differs
```

This also works for arrays of strings, structs, or any other comparable type:

```go
names1 := [2]string{"Alice", "Bob"}
names2 := [2]string{"Alice", "Bob"}
fmt.Println(names1 == names2) // true
```

Because arrays are comparable, they can also be used as **map keys** - something slices can never do:

```go
seen := map[[2]int]string{}
seen[[2]int{0, 0}] = "origin"
seen[[2]int{1, 1}] = "diagonal"
fmt.Println(seen[[2]int{0, 0}]) // origin
```

### Size Is Part of the Type

`[3]int` and `[4]int` are **different types**, not just different lengths of the same type. You cannot compare them, assign one to the other, or pass one where the other is expected:

```go
a := [3]int{1, 2, 3}
d := [4]int{1, 2, 3, 4}

fmt.Println(a == d)
// compile error: invalid operation: a == d (mismatched types [3]int and [4]int)
```

---

## Multi-Dimensional Arrays

Arrays can hold other arrays, giving you grid-shaped data - handy for boards, matrices, and tables.

### Declaring a 2D Array

```go
var grid [3][3]int
fmt.Println(grid) // [[0 0 0] [0 0 0] [0 0 0]]
```

`grid` is "an array of 3 elements, where each element is itself an array of 3 ints." Read `[3][3]int` as `[3]([3]int)`.

### Indexing a 2D Array

Index each dimension in turn: `grid[row][col]`.

```go
for row := 0; row < 3; row++ {
    for col := 0; col < 3; col++ {
        grid[row][col] = row*3 + col
    }
}
fmt.Println(grid) // [[0 1 2] [3 4 5] [6 7 8]]
```

### 2D Array Literals

```go
board := [2][3]string{
    {"X", "O", "X"},
    {"O", "X", "O"},
}
fmt.Println(board[0][1]) // O
fmt.Println(board[1][2]) // O
```

### len() on Each Dimension

```go
fmt.Println(len(board))    // 2 - number of rows
fmt.Println(len(board[0])) // 3 - number of columns in row 0
```

### Value Semantics Apply Here Too

A `[3][3]int` is still just an array (of arrays) - so everything from [Section 4](#arrays-are-value-types) applies. Copying a 2D array, or passing it to a function, copies **every** element in **every** row:

```go
boardCopy := board
boardCopy[0][0] = "changed"
// board[0][0] is untouched - boardCopy is a fully independent copy
```

---

## Iterating Arrays with range

`range` works on arrays exactly the way it works on slices (as you saw in Level 6):

```go
scores := [4]int{85, 92, 78, 95}
for i, score := range scores {
    fmt.Printf("index=%d score=%d\n", i, score)
}
```

Ignore the index when you only need values:

```go
sum := 0
for _, score := range scores {
    sum += score
}
```

### Ranging Over a 2D Array

Nest the loops, just like Level 6's nested-loop patterns:

```go
grid := [2][2]int{{1, 2}, {3, 4}}
for rowIndex, row := range grid {
    for colIndex, value := range row {
        fmt.Printf("grid[%d][%d] = %d\n", rowIndex, colIndex, value)
    }
}
```

### Reminder: range Gives You a COPY of Each Element

Because arrays are value types, the value `range` hands you each iteration is a **copy**. Modifying it does nothing to the original array:

```go
nums := [3]int{1, 2, 3}
for _, v := range nums {
    v *= 100 // modifies the local copy `v` only
}
fmt.Println(nums) // [1 2 3] - completely unchanged
```

If you need to actually mutate the array while ranging, index back into it with `nums[i] = ...` inside a `for i := range nums` loop instead.

---

## Array vs Slice: When to Use Which

**Short answer: reach for a slice almost every time.** Slices are Go's everyday, idiomatic collection type - you'll use them constantly starting in Level 8. Arrays are the exception, not the rule.

### When Arrays Are the Right Choice

- **The size is fixed by the domain, not by how much data happens to show up.** A tic-tac-toe board is always 3x3. A chess board is always 8x8. A SHA-256 hash is always 32 bytes. An IPv4 address is always 4 bytes.
- **You want value semantics on purpose.** Sometimes you *want* a copy - passing a small fixed-size record around without worrying that a callee might mutate your data.
- **You're modeling something with a genuinely fixed protocol shape**, like a fixed-size header in a binary format, where "exactly N bytes" is a hard requirement you want the type system to enforce.
- **Arrays can be map keys or compared with `==`** - if you need either of those and your data has a fixed size, an array is the natural (sometimes only) fit.

### When Slices Are Idiomatic (Almost Always)

- The number of elements can grow, shrink, or isn't known until runtime (reading lines from a file, results from a query, items a user adds).
- You want to pass the collection to a function and have mutations visible to the caller without extra ceremony.
- You want to use `append`, re-slice (`s[1:3]`), or pass sub-ranges around cheaply.
- You're accepting a collection as a function parameter - even if you happen to have an array, converting it to a slice (`arr[:]`) at the call site is the idiomatic move.

### Converting Between Them

An array can be turned into a slice with a full slice expression, and that slice **shares the array's underlying memory** (it is not a copy):

```go
arr := [5]int{10, 20, 30, 40, 50}
s := arr[:]       // slice over the whole array
s[0] = 999
fmt.Println(arr)  // [999 20 30 40 50] - arr changed too!
```

This is the one place where slicing an array does *not* behave like the value-copy rules above - a slice, by design, is a view, not a copy.

---

## Best Practices

### 1. Default to Slices, Reach for Arrays Deliberately

```go
// ✅ Good - idiomatic default
func sum(nums []int) int { ... }

// ⚠️ Only if the size is genuinely fixed and meaningful
func sum(nums [10]int) int { ... }
```

### 2. Use [...] When the Element Count Speaks for Itself

```go
// ✅ Good - reader doesn't have to count and cross-check
weekdays := [...]string{"Mon", "Tue", "Wed", "Thu", "Fri"}

// ⚠️ Fragile - easy to miscount when editing later
weekdays := [5]string{"Mon", "Tue", "Wed", "Thu", "Fri"}
```

### 3. Pass Large Arrays by Pointer If You Must Mutate Them

```go
// ✅ Good - shares the caller's array, avoids a large copy
func reset(arr *[100]int) {
    for i := range arr {
        arr[i] = 0
    }
}

// ⚠️ Copies all 100 ints on every call, and can't mutate the caller's array
func reset(arr [100]int) { ... }
```

### 4. Use Named Constants for Array Sizes

```go
// ✅ Good - the size has a name and a reason
const BoardSize = 3
var board [BoardSize][BoardSize]string

// ⚠️ Magic number, unclear if 3 is meaningful or arbitrary
var board [3][3]string
```

### 5. Prefer range Over Manual Indexing (Callback to Level 6)

```go
// ✅ Good
for i, v := range arr {
    fmt.Println(i, v)
}

// ❌ Unnecessary
for i := 0; i < len(arr); i++ {
    fmt.Println(i, arr[i])
}
```

---

## Common Mistakes

### Mistake 1: Expecting Reference Semantics

```go
// ❌ WRONG ASSUMPTION - arr inside the function is a COPY
func addOne(arr [3]int) {
    arr[0]++
}

nums := [3]int{1, 2, 3}
addOne(nums)
fmt.Println(nums) // still [1 2 3] - the caller's array never changed!

// ✅ RIGHT - use a slice, or pass a pointer, if you need mutation to be visible
func addOne(nums []int) {
    nums[0]++
}
```

### Mistake 2: Treating [3]int and [4]int as the Same Type

```go
// ❌ WRONG - these are different types; this is a compile error
var a [3]int
var b [4]int = a // cannot use a (variable of type [3]int) as [4]int value

// ✅ RIGHT - sizes must match, or use a slice if the size varies
var a [3]int
var b [3]int = a
```

### Mistake 3: Confusing [...] With "Dynamic Size"

```go
// ❌ WRONG ASSUMPTION - this is NOT a growable collection
nums := [...]int{1, 2, 3}
// nums = append(nums, 4) // compile error - arrays have no append!

// ✅ RIGHT - use a slice if the size needs to change later
nums := []int{1, 2, 3}
nums = append(nums, 4) // fine
```

### Mistake 4: Forgetting range Gives Copies

```go
// ❌ WRONG - this loop does nothing to the original array
prices := [3]float64{9.99, 19.99, 29.99}
for _, p := range prices {
    p *= 1.1 // only changes the local copy `p`
}

// ✅ RIGHT - index back into the array to mutate it
for i := range prices {
    prices[i] *= 1.1
}
```

### Mistake 5: Assuming Out-of-Range Indexing Is Always Caught the Same Way

```go
arr := [3]int{1, 2, 3}

// A CONSTANT index out of range is a COMPILE error:
// fmt.Println(arr[3])
// invalid argument: index 3 out of bounds [0:3]

// A VARIABLE index out of range is a RUNTIME panic instead:
i := 3
fmt.Println(arr[i])
// panic: runtime error: index out of range [3] with length 3
```

### Mistake 6: Reaching for Arrays by Default

```go
// ⚠️ Usually wrong instinct - arrays as general-purpose containers
func processUsers(users [100]User) { ... } // what if there are 3 users? or 500?

// ✅ RIGHT - slices handle "however many there happen to be"
func processUsers(users []User) { ... }
```

---

## Summary

**What Arrays Are:**
- Fixed-size, homogeneous collections - the length is part of the type (`[5]int` != `[10]int`)
- The zero value is immediately usable - every element set to its type's zero value

**Declaring:**
- `var arr [5]int`, `[3]string{"a","b","c"}`, `[...]int{1,2,3}`, sparse `[5]int{0: 1, 4: 5}`

**Value Semantics (the big one):**
- Assigning an array copies it; passing an array to a function passes a copy
- Slices, by contrast, share their underlying array on assignment/passing
- `range` also hands you copies of each element

**Comparison:**
- `==` and `!=` work between arrays of the same type (same element type AND same size)
- This makes arrays usable as map keys, unlike slices

**Multi-Dimensional:**
- `[3][3]int` is an array of arrays - index with `grid[row][col]`

**Array vs Slice:**
- Slices are the idiomatic default almost everywhere
- Arrays shine when the size is fixed by the domain, or when you want value semantics or comparability on purpose

---

## Next Steps

You now understand:
- ✅ What makes arrays fixed-size and homogeneous
- ✅ Every way to declare an array, including size inference and sparse literals
- ✅ Why arrays are value types - and how that differs from slices
- ✅ How array comparison and multi-dimensional arrays work
- ✅ When to actually reach for an array instead of a slice

**Next level:** Level 8 - Slices
- Slices in depth: the pointer/length/capacity header
- `append`, growth, and re-slicing
- Why slices share memory - and when that surprises you

You now understand the foundation slices are built on. Keep going! 🚀
