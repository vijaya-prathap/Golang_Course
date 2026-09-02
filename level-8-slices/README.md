# Level 8: Slices - Complete Guide

## Introduction

Welcome to Level 8! You've mastered arrays (Level 7) - fixed-size, value-semantic collections. Now it's time to learn **slices**, the collection type you'll actually use almost everywhere in real Go code.

A slice is **not** a fixed-size container like an array. A slice is a small, lightweight **view** into an underlying array: a pointer, a length, and a capacity. Understanding that one sentence is the key to everything in this level - growth, aliasing, and every "gotcha" Go programmers eventually hit.

---

## Table of Contents

1. [What Is a Slice?](#what-is-a-slice)
2. [Creating Slices](#creating-slices)
3. [len() vs cap()](#len-vs-cap)
4. [append() and Growth Semantics](#append-and-growth-semantics)
5. [Slicing Shares the Underlying Array](#slicing-shares-the-underlying-array)
6. [The Three-Index (Full) Slice Expression](#the-three-index-full-slice-expression)
7. [copy() - Making Independent Copies](#copy---making-independent-copies)
8. [nil Slice vs Empty Slice](#nil-slice-vs-empty-slice)
9. [Two-Dimensional Slices ([][]T)](#two-dimensional-slices-t)
10. [Best Practices](#best-practices)
11. [Common Mistakes](#common-mistakes)

---

## What Is a Slice?

An array's size is part of its type (`[5]int` and `[10]int` are different types) and copying an array copies every element. A slice fixes both problems by being a **view**, not a container.

Under the hood, a slice header is three fields:

```
type sliceHeader struct {
    ptr *T   // pointer to the first element the slice can see
    len int  // number of elements currently in the slice
    cap int  // number of elements available before a reallocation is needed
}
```

You never write that struct yourself - it's the Go runtime's internal representation - but every slice you ever use is exactly this: a pointer into *some* underlying array, plus a length and a capacity.

```go
arr := [5]int{10, 20, 30, 40, 50}
s := arr[1:4] // s is a VIEW into arr, not a new array
```

**The most important idea in this level:** a slice does not own its data the way an array owns its elements. It borrows a window into an array that lives somewhere else. That's why slices are cheap to pass around (just copy the header - 3 machine words) and why mutating through one slice can be visible through another.

---

## Creating Slices

There are four common ways to create a slice.

### 1. Slice Literal

```go
nums := []int{10, 20, 30}
// len=3 cap=3 - Go allocates a backing array sized exactly to the literal
```

Notice there's no length between the brackets (`[]int`, not `[3]int`) - that's what makes it a slice type instead of an array type. It's easy to typo one as the other, so watch the brackets carefully.

### 2. make() With Length Only

```go
zeros := make([]int, 4)
// [0 0 0 0]  len=4 cap=4 - zero-valued elements, cap equals len
```

Use this when you know how many elements you need up front and will fill them by index.

### 3. make() With Length AND Capacity

```go
buf := make([]int, 2, 6)
// [0 0]  len=2 cap=6 - room to grow to 6 elements before reallocating
```

This is the pattern to reach for when you know roughly how big a slice will eventually get (e.g., building up a result with `append` in a loop) - it avoids repeated reallocation.

### 4. Slicing an Existing Array or Slice

```go
arr := [6]int{10, 20, 30, 40, 50, 60}
s := arr[1:4]
// s == [20 30 40], len=3, cap=5 (see the len vs cap section for why cap is 5)
```

`s[low:high]` means "give me a view starting at index `low`, up to but not including index `high`." Both bounds are optional: `s[:3]` starts at 0, `s[3:]` runs to the end, `s[:]` is the whole thing.

```go
nums := []int{1, 2, 3, 4, 5, 6, 7, 8}
nums[:3]  // [1 2 3]
nums[3:]  // [4 5 6 7 8]
nums[2:5] // [3 4 5]
nums[:]   // [1 2 3 4 5 6 7 8]
```

---

## len() vs cap()

- `len(s)` - how many elements are currently in the slice (what `for range` walks over, what `s[i]` lets you index up to).
- `cap(s)` - how many elements the underlying array has room for *starting from the slice's own pointer* before Go must allocate a new array.

These are **not** always equal, and the difference is the whole reason `append` sometimes mutates in place and sometimes doesn't.

```
Underlying array (length 6):   index:   0    1    2    3    4    5
                                value:  10   20   30   40   50   60

s := arr[1:4]

                pointer──┐
                         ▼
                        20   30   40 │  50   60
                        └─────┬─────┘  └──┬──┘
                          len(s) = 3   (unused capacity)

                cap(s) = 5   (from index 1 through the end of the array: 20,30,40,50,60)
```

Verified with `go run` on this Go toolchain:

```go
arr := [6]int{10, 20, 30, 40, 50, 60}
s := arr[1:4]
fmt.Println(s, len(s), cap(s))
```

```
[20 30 40] 3 5
```

`cap` is measured from where the slice *starts*, not from the start of the original array. A slice that starts later into the same array has a smaller capacity:

```go
nums := []int{1, 2, 3, 4, 5, 6, 7, 8} // len=8 cap=8
mid := nums[3:5]                      // len=2 cap=5 (index 3 through index 7)
```

---

## append() and Growth Semantics

`append(s, v)` does one of two things, and which one happens is invisible unless you check:

1. **If `len(s) < cap(s)`:** there's room. Go writes the new element right after the last one, **in the same underlying array**, and returns a slice with `len+1`. No allocation. Any other slice that shares that array can observe the write.
2. **If `len(s) == cap(s)`:** there's no room. Go allocates a **new**, bigger underlying array, copies every existing element over, writes the new element, and returns a slice pointing at the *new* array. The old array is now completely disconnected from the returned slice.

This is the single most common source of surprising bugs for people coming from other languages: **append sometimes returns a slice backed by the same array you started with, and sometimes doesn't**, and you cannot tell which happened just by looking at the call.

### Proof: Case 1 - append() Within Capacity Mutates the Shared Array

```go
base := make([]int, 3, 5)
base[0], base[1], base[2] = 1, 2, 3
view := base[:2]
fmt.Printf("base: %v ptr=%p\n", base, base)
fmt.Printf("view: %v ptr=%p\n", view, view)

view = append(view, 999)
fmt.Printf("base: %v ptr=%p\n", base, base)
fmt.Printf("view: %v ptr=%p\n", view, view)
```

**Actual output (addresses will differ on your machine, but will match each other exactly as shown):**

```
base: [1 2 3] ptr=0x671f46834030
view: [1 2] ptr=0x671f46834030

base: [1 2 999] ptr=0x671f46834030
view: [1 2 999] ptr=0x671f46834030
```

`view` had `cap == 5` and `len == 2`, so `append` had room. It silently overwrote `base[2]` - a variable `view` never even mentioned - because both slices pointed at the same array.

### Proof: Case 2 - append() Beyond Capacity Reallocates

```go
small := make([]int, 2, 2)
small[0], small[1] = 10, 20
grown := append(small, 30)
grown[0] = 111
fmt.Println(small) // unaffected
fmt.Println(grown)
```

**Actual output:**

```
[10 20]
[111 20 30]
```

`small` had `len == cap == 2`, so `append` had to grow. Go allocated a brand-new array, copied `10, 20` into it, appended `30`, and returned that as `grown`. Mutating `grown[0]` afterward has zero effect on `small` because they no longer share any memory.

**Why this matters:** never assume `append` did (or didn't) reallocate. If you need the old slice to observe changes, make sure it shares capacity on purpose (or just don't rely on it - see Best Practices). If you need the old slice to be *protected* from changes, either make sure a reallocation happened, or copy explicitly.

---

## Slicing Shares the Underlying Array

Taking a sub-slice (`s[low:high]`) does not copy data - it creates a second header pointing at the same array. Mutating an element through either slice is visible through the other.

```go
original := []int{1, 2, 3, 4, 5}
sub := original[1:4]
sub[0] = 999
fmt.Println(sub)      // [999 3 4]
fmt.Println(original) // [1 999 3 4 5] - index 1 changed too!
```

**Actual verified output:**

```
sub: [999 3 4]
original: [1 999 3 4 5] (index 1 changed too!)
```

Two slices can even overlap partially:

```go
base := []int{100, 200, 300, 400, 500}
left := base[0:3]  // [100 200 300]
right := base[2:5] // [300 400 500]

left[2] = 9999
fmt.Println(left)  // [100 200 9999]
fmt.Println(right) // [9999 400 500] - right[0] IS base[2] IS left[2]
```

`left` and `right` overlap at `base[2]`; writing through one is visible through the other because there is only ever **one** array underneath both of them.

---

## The Three-Index (Full) Slice Expression

The two-index form `s[low:high]` sets `len` but leaves `cap` extending all the way to the end of the underlying array - which is exactly what caused the aliasing surprise above. The **three-index** (full) slice expression lets you cap the capacity too:

```go
s[low:high:max]
// len = high - low
// cap = max  - low
```

### Proof: Regular Slicing Leaks Capacity

```go
source := []int{1, 2, 3, 4, 5, 6, 7, 8}
leaky := source[1:3] // len=2, cap=7 (leaks all the way to source's end)
leaky = append(leaky, 999)
fmt.Println(leaky)  // [2 3 999]
fmt.Println(source) // [1 2 3 999 5 6 7 8] - source[3] got clobbered!
```

### Proof: Three-Index Slicing Prevents It

```go
source2 := []int{1, 2, 3, 4, 5, 6, 7, 8}
safe := source2[1:3:3] // len=2, cap=2 - capacity capped at "3"
safe = append(safe, 999)
fmt.Println(safe)    // [2 3 999]
fmt.Println(source2) // [1 2 3 4 5 6 7 8] - unchanged!
```

**Actual verified output:**

```
safe := source2[1:3:3]: [2 3] len: 2 cap: 2

--- after safe = append(safe, 999) ---
safe: [2 3 999]
source2: [1 2 3 4 5 6 7 8] (unchanged - append forced a reallocation)
```

Because `cap(safe)` was forced down to `2` (equal to its own `len`), the moment you `append` to it, `len == cap` triggers a reallocation instead of writing into `source2`'s backing array. Reach for `s[low:high:max]` whenever you hand a sub-slice to code you don't control and want to guarantee it can never silently corrupt your original array.

---

## copy() - Making Independent Copies

`copy(dst, src)` copies `min(len(dst), len(src))` elements from `src` into `dst` and returns the number of elements copied. Unlike slicing, it duplicates the *values* into `dst`'s own backing array - there's no aliasing afterward.

```go
src := []int{1, 2, 3, 4, 5}
dst := make([]int, len(src))
n := copy(dst, src)

dst[0] = 999
fmt.Println(src) // [1 2 3 4 5]   - unaffected
fmt.Println(dst) // [999 2 3 4 5]
fmt.Println(n)   // 5
```

If `dst` is shorter than `src`, only `len(dst)` elements are copied - `copy` never grows `dst` for you:

```go
small := make([]int, 3)
copied := copy(small, src) // copies only the first 3 elements
fmt.Println(small, copied) // [1 2 3] 3
```

**The idiomatic "give me my own copy" pattern:**

```go
clone := make([]T, len(original))
copy(clone, original)
```

Use `copy()` whenever you need to hand a slice to code that might mutate it, and you want your own version to stay untouched.

---

## nil Slice vs Empty Slice

```go
var nilSlice []int      // zero value of a slice type - the pointer is nil
emptySlice := []int{}   // a real (empty) backing array - the pointer is non-nil
```

Both have `len == 0` and `cap == 0`, and both print as `[]` with `%v`. But they are **not** the same thing:

```go
fmt.Println(nilSlice == nil)   // true
fmt.Println(emptySlice == nil) // false
```

**Actual verified output:**

```
nilSlice == nil: true
emptySlice == nil: false

nilSlice:   []   []int(nil)
emptySlice: []   []int{}
```

(`%#v` reveals the difference that `%v` hides.)

Despite that difference, **every common operation treats them identically**:

```go
nilSlice = append(nilSlice, 1)     // works fine, [1]
emptySlice = append(emptySlice, 1) // works fine, [1]

for range nilSlice { }  // zero iterations, no panic
len(nilSlice)           // 0, no panic
```

**Rule of thumb:** don't write `if mySlice == nil` to mean "is it empty?" unless you specifically care about the nil-vs-non-nil distinction (e.g., JSON encodes a nil slice as `null` but an empty slice as `[]`). To check emptiness, use `len(mySlice) == 0` - it's correct for both.

---

## Two-Dimensional Slices ([][]T)

A `[][]int` is a slice of slices - each "row" is its own independent slice with its own backing array. There is no built-in 2D array type the way some languages have; you build the grid yourself, one row at a time.

### The Shared-Row Bug

```go
func buildGridWrong(rows, cols int) [][]int {
    row := make([]int, cols)
    grid := make([][]int, rows)
    for i := range grid {
        grid[i] = row // BUG: every row IS the same slice
    }
    return grid
}
```

```go
wrong := buildGridWrong(3, 3)
wrong[0][0] = 1
// [1 0 0]
// [1 0 0]
// [1 0 0]     <- every row changed!
```

**Actual verified output confirms the bug:**

```
After setting wrong[0][0] = 1:
[1 0 0]
[1 0 0]
[1 0 0]
(every row shows the change - they all share ONE underlying array!)
```

`row` was created **once**, outside the loop, and every entry in `grid` was assigned that same slice header. All three "rows" point at the exact same backing array.

### The Correct Way

```go
func buildGridRight(rows, cols int) [][]int {
    grid := make([][]int, rows)
    for i := range grid {
        grid[i] = make([]int, cols) // a FRESH backing array every iteration
    }
    return grid
}
```

```go
right := buildGridRight(3, 3)
right[0][0] = 1
right[1][1] = 2
right[2][2] = 3
// [1 0 0]
// [0 2 0]
// [0 0 3]     <- each row is independent
```

### 2D Slice Literal

For fixed, known-upfront data, you can also write the whole grid as a literal:

```go
grid := [][]int{
    {1, 2, 3},
    {4, 5, 6},
    {7, 8, 9},
}
```

Each `{...}` row literal already gets its own backing array, so the shared-row bug doesn't apply to this form - it only bites you when you build rows programmatically by reusing one slice variable.

---

## Best Practices

### 1. Preallocate Capacity When the Size Is Known

```go
// ✅ Good - one allocation
result := make([]int, 0, len(input))
for _, v := range input {
    result = append(result, transform(v))
}

// ❌ Wasteful - repeated reallocation as the slice grows
var result []int
for _, v := range input {
    result = append(result, transform(v))
}
```

### 2. Use the Three-Index Slice When Handing Out Sub-Slices

```go
// ✅ Good - caller's append() can't corrupt your array
func firstN(s []int, n int) []int {
    return s[:n:n]
}

// ❌ Risky - caller's append() can silently overwrite s[n], s[n+1], ...
func firstN(s []int, n int) []int {
    return s[:n]
}
```

### 3. Use copy() When You Need Independence, Not Just a Slice Expression

```go
// ✅ Good - clone truly independent from original
clone := make([]int, len(original))
copy(clone, original)

// ❌ Still aliased - this is a view, not a copy
alias := original[:]
```

### 4. Check Emptiness With len(), Not == nil

```go
// ✅ Good - correct for both nil and non-nil empty slices
if len(tasks) == 0 { }

// ❌ Fragile - misses the case of a non-nil empty slice
if tasks == nil { }
```

### 5. Don't Build a New Slice of Sub-Slices From One Reused Buffer

```go
// ✅ Good - each window gets copied into its own array
var windows [][]int
for i := 0; i < len(data)-1; i++ {
    w := make([]int, 2)
    copy(w, data[i:i+2])
    windows = append(windows, w)
}

// ❌ Every window still aliases data - a later data[i] = x can retroactively
//    "change" windows you already stored
var windows [][]int
for i := 0; i < len(data)-1; i++ {
    windows = append(windows, data[i:i+2])
}
```

---

## Common Mistakes

### Mistake 1: Thinking Slice Literal Syntax Means Slices Are Arrays

```go
// ❌ WRONG ASSUMPTION - "[]int{1,2,3} looks like [3]int{1,2,3}, so they must behave the same"
var a [3]int = [3]int{1, 2, 3}   // array: fixed size, copies by value
var s []int = []int{1, 2, 3}     // slice: dynamic view, copies by reference-like header

// ✅ RIGHT - remember the missing number between the brackets is the whole story:
//    [3]int is an array type. []int is a slice type. Different types entirely.
```

### Mistake 2: Assuming append() Always Mutates in Place

```go
// ❌ WRONG ASSUMPTION
s1 := make([]int, 3, 3) // len == cap, no room
s2 := append(s1, 4)     // forces a reallocation - s2 is now a DIFFERENT array
s2[0] = 999
fmt.Println(s1[0]) // still the old value, NOT 999!

// ✅ RIGHT - always use the slice append() returns, and never assume the old
//    slice variable will (or won't) see later mutations through the new one
```

### Mistake 3: Assuming append() Never Mutates the Original

```go
// ❌ WRONG ASSUMPTION
s1 := make([]int, 2, 5) // len < cap, room to spare
s2 := append(s1, 99)    // NO reallocation - s2 shares s1's array
s2[0] = 111
fmt.Println(s1[0]) // 111 - s1 changed too, even though you never touched s1 directly!

// ✅ RIGHT - if you need s1 protected, slice it with a capped capacity first:
//    s1 := arr[low:high:high]
```

### Mistake 4: The Shared-Row Bug in 2D Slices

```go
// ❌ WRONG - every row is the same slice
row := make([]int, cols)
grid := make([][]int, rows)
for i := range grid {
    grid[i] = row
}

// ✅ RIGHT - allocate a new row each iteration
grid := make([][]int, rows)
for i := range grid {
    grid[i] = make([]int, cols)
}
```

### Mistake 5: Using == nil to Check for "Empty"

```go
// ❌ WRONG - emptySlice := []int{} is NOT nil, so this misses it
if tasks == nil {
    fmt.Println("no tasks")
}

// ✅ RIGHT - len() is correct for both nil and non-nil empty slices
if len(tasks) == 0 {
    fmt.Println("no tasks")
}
```

### Mistake 6: Sub-Slicing in a Loop Without Copying

```go
// ❌ RISKY - every stored slice still aliases data's backing array
var chunks [][]int
for i := 0; i < len(data); i += 2 {
    chunks = append(chunks, data[i:i+2]) // no copy!
}
// mutating data later silently changes chunks too

// ✅ RIGHT - copy each chunk into its own backing array
var chunks [][]int
for i := 0; i < len(data); i += 2 {
    c := make([]int, 2)
    copy(c, data[i:i+2])
    chunks = append(chunks, c)
}
```

---

## Summary

**What a Slice Is:**
- A header of `{pointer, len, cap}` pointing into an underlying array - NOT a container that owns fixed-size storage like an array

**Creating Slices:**
- Literal: `[]int{1, 2, 3}`
- `make([]int, len)` or `make([]int, len, cap)`
- Slicing: `arr[low:high]`, `s[low:high]`

**len() vs cap():**
- `len` = elements currently in the slice; `cap` = room before a reallocation is needed, measured from where the slice starts

**append():**
- Mutates the shared array in place when `len < cap`
- Reallocates to a new array when `len == cap` - breaking the connection to the old array
- Never assume which one happened without checking

**Aliasing:**
- Sub-slices share the underlying array; mutating one is visible through any other slice that overlaps it
- `s[low:high:max]` caps capacity to prevent a later `append` from leaking into shared memory

**copy():**
- The only way to get a truly independent duplicate of a slice's data

**nil vs Empty:**
- `var s []int` is nil; `s := []int{}` is not nil - but both behave identically under `len`, `range`, and `append`

**2D Slices:**
- Always allocate each row's own backing array - reusing one row slice for every row is the classic bug

---

## Next Steps

You now understand:
- ✅ What a slice really is under the hood (pointer + len + cap)
- ✅ How to create slices four different ways
- ✅ The difference between len() and cap(), and why it matters
- ✅ Exactly when append() mutates in place vs. reallocates
- ✅ How sub-slices alias their underlying array
- ✅ How to use the three-index slice expression to prevent aliasing surprises
- ✅ How copy() gives you true independence from the original data
- ✅ The nil-vs-empty-slice distinction
- ✅ How to correctly build 2D slices without the shared-row bug

**Next level:** Level 9 - Strings & Runes
- Strings as immutable byte sequences
- Runes and Unicode code points
- Converting between strings, []byte, and []rune
- Common string-processing patterns

Slices are the workhorse of Go - you'll use them in nearly every program from here on. Keep going! 🚀
