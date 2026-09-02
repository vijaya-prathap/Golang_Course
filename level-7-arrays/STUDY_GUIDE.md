# Level 7: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: Arrays
```
Day 1:  What arrays are, declaring arrays (var, literals, [...], sparse)
Day 2:  Indexing, len(), zero values, bounds checking
Day 3:  Value semantics - assignment copies (THE big idea)
Day 4:  Value semantics - function parameters copy too
Day 5:  Comparing arrays with ==
Day 6:  Multi-dimensional arrays
Day 7:  range over arrays, array vs slice decision-making
```

### Week 2: Practice & Application
```
Day 1:  Exercises 1-3 (declaring, indexing, assignment copy)
Day 2:  Exercises 4-5 (function copy, comparison)
Day 3:  Exercises 6-7 (multi-dimensional, range)
Day 4:  Exercises 8-10 (array/slice conversion, tic-tac-toe, comprehensive)
Day 5:  Bonus challenges
Day 6-7: Practice & review
```

---

## 🎯 Array Declaration Forms

| Form | Syntax | Length |
|------|--------|--------|
| **var, zero-valued** | `var arr [5]int` | 5 (fixed, elements are zero values) |
| **literal, explicit size** | `[3]string{"a","b","c"}` | 3 (must match element count) |
| **literal, inferred size** | `[...]int{1,2,3}` | 3 (counted at compile time) |
| **sparse, explicit size** | `[5]int{0: 1, 4: 5}` | 5 (unset indices are zero) |
| **sparse, inferred size** | `[...]int{2: 100, 5: 200}` | 6 (highest index + 1) |
| **declare then assign** | `var g [4]float64; g[0] = 1` | 4 (fixed, filled by index) |

```
[N]T
 │└─ element type (must be uniform)
 └── fixed length - part of the TYPE itself

[5]int  and  [10]int  are as different as  int  and  string
```

---

## 🧊 The Array Copy Diagram

This is the single most important picture in this level. Two SEPARATE memory
blocks exist after assignment - not one block with two names pointing at it.

```go
a := [3]int{1, 2, 3}
b := a       // COPY - a brand new [3]int, independent of a
b[0] = 99
```

```
BEFORE b[0] = 99:

   a  ──────────────┐
                     ▼
                ┌───┬───┬───┐
                │ 1 │ 2 │ 3 │   (a's own memory block)
                └───┴───┴───┘

   b  ──────────────┐
                     ▼
                ┌───┬───┬───┐
                │ 1 │ 2 │ 3 │   (b's own SEPARATE memory block)
                └───┴───┴───┘

AFTER b[0] = 99:

   a  ──────────────┐
                     ▼
                ┌───┬───┬───┐
                │ 1 │ 2 │ 3 │   ← a is COMPLETELY UNTOUCHED
                └───┴───┴───┘

   b  ──────────────┐
                     ▼
                ┌───┬───┬───┐
                │99 │ 2 │ 3 │   ← only b's block changed
                └───┴───┴───┘

Result:  a == [1 2 3]   b == [99 2 3]
```

### Contrast: What a Slice Would Do

A slice is a small header (pointer, length, capacity). Assigning a slice
copies the HEADER, not the data it points to - so both headers still point
at the same block.

```go
sa := []int{1, 2, 3}
sb := sa     // copies the HEADER only
sb[0] = 99
```

```
   sa header: {ptr ──┐, len:3, cap:3}
                      │
                      ▼
                 ┌───┬───┬───┐
                 │99 │ 2 │ 3 │   ← ONE shared memory block
                 └───┴───┴───┘
                      ▲
                      │
   sb header: {ptr ──┘, len:3, cap:3}

Result:  sa == [99 2 3]   sb == [99 2 3]   (BOTH changed - same backing array)
```

**The same `x := y` syntax means "copy the data" for arrays and "copy the header
(share the data)" for slices.** This is the root of nearly every array/slice
surprise you'll run into.

### The Same Rule Applies to Function Calls

Passing an argument to a function is just "assigning it to the parameter" -
so everything above applies identically:

```
func f(arr [3]int)   → arr is a COPY.   Mutating it never touches the caller's array.
func f(s []int)      → s shares the underlying array. Mutating s[i] DOES touch the caller's data.
```

---

## 🗺️ Multi-Dimensional Array Indexing

```go
var grid [3][3]int
grid[1][2] = 7
```

```
grid is "an array of 3 rows, each row is an array of 3 ints"

        col:  0    1    2
              │    │    │
row 0  ─────┬───┬───┬───┐
        │ 0 │ 0 │ 0 │
        ├───┼───┼───┤
row 1  ─────┤ 0 │ 0 │ 7 │  ← grid[1][2] = 7
        ├───┼───┼───┤
row 2  ─────┤ 0 │ 0 │ 0 │
        └───┴───┴───┘

grid[row][col]
     └─┬─┘  └─┬─┘
    outer   inner
    index   index
   (which  (position
    row)    in that row)

len(grid)     → 3   (number of rows - the outer dimension)
len(grid[0])  → 3   (number of columns - the inner dimension)
```

### Iterating a 2D Array

```go
for rowIndex, row := range grid {       // row is a COPY of that [3]int row
    for colIndex, value := range row {  // value is a COPY of that int
        fmt.Println(rowIndex, colIndex, value)
    }
}
```

---

## ⚖️ Array Comparison Rules

```
Comparable with ==  when:
  ✅ same element type
  ✅ same length
  ✅ element type is itself comparable (numbers, strings, bools, structs
     of comparable fields, other arrays of comparable elements)

NOT comparable when:
  ❌ different lengths   → [3]int vs [4]int is a COMPILE ERROR
  ❌ different element types
  ❌ elements aren't comparable (e.g. an array of slices, maps, or funcs)
```

```go
[3]int{1,2,3} == [3]int{1,2,3}   // true  - compiles, same type
[3]int{1,2,3} == [3]int{1,2,4}   // false - compiles, same type
[3]int{1,2,3} == [4]int{1,2,3,4} // compile error - different types entirely
```

Because arrays are comparable, they can be **map keys** - slices never can:

```go
seen := map[[2]int]string{}
seen[[2]int{0, 0}] = "origin"
```

---

## 🧩 Array vs Slice Decision Tree

```
Need a collection?
├─ Size is fixed by the DOMAIN (board size, hash length, header format)?
│  ├─ Want value semantics (copy) or need it as a map key / to compare with ==?
│  │  └─ ✅ Use an ARRAY: [N]T
│  └─ Just happens to be a known constant but no comparison/copy need?
│     └─ Consider a slice anyway - it's more flexible with little cost
└─ Size can grow/shrink, or isn't known until runtime?
   └─ ✅ Use a SLICE: []T   (this is almost always the answer)
```

| Question | Array | Slice |
|----------|-------|-------|
| Can it grow with `append`? | ❌ No | ✅ Yes |
| Copies on assignment/pass? | ✅ Yes (full copy) | ❌ No (shares backing array) |
| Comparable with `==`? | ✅ Yes (same type) | ❌ No (compile error) |
| Can be a map key? | ✅ Yes | ❌ No |
| Size is part of the type? | ✅ Yes | ❌ No |
| Idiomatic default in Go? | ❌ Rare | ✅ Almost always |

---

## 🚨 Common Mistakes

### Mistake 1: Expecting Reference Semantics

```go
// ❌ WRONG - arr inside the function is a COPY
func addOne(arr [3]int) { arr[0]++ }

nums := [3]int{1, 2, 3}
addOne(nums)
fmt.Println(nums) // still [1 2 3]!
```

### Mistake 2: Treating [3]int and [4]int as the Same Type

```go
// ❌ WRONG - compile error, mismatched types
var a [3]int
var b [4]int = a
```

### Mistake 3: Confusing [...] With "Dynamic Size"

```go
// ❌ WRONG - arrays have no append; this is still fixed-size
nums := [...]int{1, 2, 3}
// nums = append(nums, 4) // compile error
```

### Mistake 4: Forgetting range Gives Copies

```go
// ❌ WRONG - does nothing to the original array
for _, p := range prices {
    p *= 1.1 // only changes the local copy
}

// ✅ RIGHT
for i := range prices {
    prices[i] *= 1.1
}
```

### Mistake 5: Assuming All Out-of-Range Access Panics at Runtime

```go
arr := [3]int{1, 2, 3}
// fmt.Println(arr[3])   // CONSTANT index → COMPILE error

i := 3
fmt.Println(arr[i])      // VARIABLE index → RUNTIME panic
```

---

## 📈 Progression Summary

### Understanding Level 7

Level 7 teaches the fixed-size collection type that underlies slices:

1. **Declaration** - every syntax for building an array
2. **Value semantics** - the defining trait: copy on assignment and on call
3. **Comparison** - arrays can be compared and used as map keys
4. **Multi-dimensional** - arrays of arrays for grid-shaped data
5. **When to choose arrays** - rare, but real: fixed protocols, comparability, deliberate copying

### Prerequisites for Level 8

Before moving to Level 8 (Slices), you need:

- ✅ Comfortable declaring arrays with all forms (var, literal, [...], sparse)
- ✅ Can explain and DEMONSTRATE that arrays copy on assignment and function calls
- ✅ Understand why [3]int and [4]int are different types
- ✅ Can compare arrays with == and know when it's legal
- ✅ Can build and iterate multi-dimensional arrays
- ✅ Know when an array is the right choice instead of a slice (and that this is rare)

### Ready for Level 8?

Level 8 teaches the collection type you'll use nearly everywhere in Go:
- Slices as a header over an array (pointer, length, capacity)
- append, growth strategy, and re-slicing
- Why slices sharing memory is a feature - and a footgun if you forget it

---

## ✅ Checklist Before Level 8

- [ ] Can declare an array five different ways
- [ ] Can prove (not just state) that array assignment copies the data
- [ ] Can prove that passing an array to a function passes a copy
- [ ] Know why the same slice code shares data instead of copying
- [ ] Can compare two arrays with == and explain the size/type rule
- [ ] Can build and index a 2D array (`grid[row][col]`)
- [ ] Can explain when an array is the right choice over a slice
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The Type
An array's length is part of its type - `[3]int` and `[4]int` are as different as `int` and `string`.

### The Copy
Arrays are value types. Assignment and function calls copy the entire array - always. This is the #1 thing to internalize before Level 8.

### The Comparison
Arrays (of comparable element type and equal size) support `==`, `!=`, and can be map keys - slices can do none of that.

### The Default
Reach for a slice almost every time. Arrays are for genuinely fixed-size domains, or when you want copy semantics or comparability on purpose.

---

## 📚 Next Level

Level 8: Slices
- Slices in depth: the pointer/length/capacity header
- append, growth, and re-slicing
- Why slices share memory - and when that surprises you

You've got arrays down! Keep going! 🚀
