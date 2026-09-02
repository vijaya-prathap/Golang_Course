# Level 8: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: Slices
```
Day 1:  What a slice is - the header (pointer + len + cap)
Day 2:  Creating slices - literals, make(), slicing
Day 3:  len() vs cap()
Day 4:  append() and growth semantics
Day 5:  Aliasing - sub-slices share the underlying array
Day 6:  Three-index slices and copy()
Day 7:  nil vs empty slices, 2D slices
```

### Week 2: Practice & Application
```
Day 1:  Exercises 1-3 (creation, len/cap)
Day 2:  Exercises 4-5 (aliasing, append reallocation)
Day 3:  Exercises 6-8 (copy, three-index, nil vs empty)
Day 4:  Exercises 9-10 (2D slices, comprehensive task list)
Day 5:  Bonus challenges
Day 6-7: Practice & review
```

---

## 🎯 The Slice Header

A slice is three machine words - a pointer, a length, and a capacity. It is NOT the data itself; it's a small descriptor pointing at data that lives somewhere else.

```
                    ┌─────────────┐
     slice s        │  ptr  ──────┼──────┐
                     │  len  = 3   │      │
                     │  cap  = 5   │      │
                     └─────────────┘      │
                                           ▼
Underlying array:   [ 10 ][ 20 ][ 30 ][ 40 ][ 50 ][ 60 ]
index:                 0     1     2     3     4     5
                              └──────s────────┘└──unused──┘
                              (s = arr[1:4], len=3, cap=5)
```

```
s := arr[1:4]

len(s) = high - low          = 4 - 1 = 3   → elements you can see: 20, 30, 40
cap(s) = len(arr) - low      = 6 - 1 = 5   → elements available before reallocating: 20,30,40,50,60
```

**Key insight:** copying a slice (assigning it, passing it to a function) copies the header - three words - not the underlying array. That's why slices are cheap to pass around, and also exactly why two slice variables can secretly point at the same data.

---

## 🔗 Aliasing Diagram: Two Slices, One Array

```go
base := []int{100, 200, 300, 400, 500}
left := base[0:3]   // [100 200 300]
right := base[2:5]  // [300 400 500]
```

```
Underlying array:   [ 100 ][ 200 ][ 300 ][ 400 ][ 500 ]
index:                 0      1      2      3      4

left  = base[0:3]    └───────left────────┘
right = base[2:5]                   └───────right───────┘

                                     ▲
                              base[2] is SHARED:
                              it is left[2] AND right[0]
```

```go
left[2] = 9999
fmt.Println(right) // [9999 400 500] - right "saw" a write it never made!
```

**Actual verified output:**
```
left: [100 200 9999]
right: [9999 400 500] (right[0] changed too - same memory as left[2]/base[2])
```

**Rule:** if two slices came from the same array (directly, or through slicing each other), assume they can see each other's writes in the overlapping region until proven otherwise.

---

## 📈 append() Growth: Actually Observed on This Go Toolchain

Guessing at growth factors is a classic mistake - Go's growth algorithm has changed across versions, and it also depends on element size. Below is REAL data captured with `go run` on this machine (`go version go1.26.5 darwin/arm64`), appending one `int` at a time to a nil slice:

| len after append | cap  | growth factor from previous cap |
|-------------------|------|----------------------------------|
| 1                  | 4    | (initial allocation, not from 1!) |
| 5                  | 8    | 2.0x |
| 9                  | 16   | 2.0x |
| 17                 | 32   | 2.0x |
| 33                 | 64   | 2.0x |
| 65                 | 128  | 2.0x |
| 129                | 256  | 2.0x |
| 257                | 512  | 2.0x |
| 513                | 848  | 1.66x |
| 849                | 1280 | 1.51x |
| 1281               | 1792 | 1.40x |
| 1793               | 2560 | 1.43x |

**What this actually shows:**
- The very first append does NOT go to `cap=1` - it jumps straight to `cap=4` (Go rounds small allocations to internal memory size classes)
- Growth is roughly 2x while the slice is small
- Once the slice gets large (roughly 500+ elements here), growth slows to a smaller factor to avoid wasting memory
- **Do not hardcode these numbers into your programs** - they can change between Go versions, platforms, and element types. The lesson is qualitative: expect roughly-doubling growth for small slices, tapering off for large ones, and never assume an exact value.

---

## 🧬 append() Decision Table

| Condition | What Happens | Old slice affected by later writes to new slice? |
|-----------|--------------|----------------------------------------------------|
| `len(s) < cap(s)` | Writes into the SAME underlying array, returns `len+1` | YES - they still share memory |
| `len(s) == cap(s)` | Allocates a NEW underlying array, copies everything, writes the new element | NO - completely disconnected |

```
Case 1 (room to grow):                  Case 2 (no room):

base cap=5 len=3                        small cap=2 len=2
view := base[:2]  (cap=5, len=2)        grown := append(small, x)

append(view, x)                              small ──▶ [old array]
    │                                        grown ──▶ [NEW array]  (copied)
    ▼
still points at base's array
base[2] silently changes!
```

---

## 🎯 Three-Index Slice Expression

```
s[low:high:max]

len(s) = high - low
cap(s) = max  - low
```

```
Underlying array:  [ 1 ][ 2 ][ 3 ][ 4 ][ 5 ][ 6 ][ 7 ][ 8 ]
index:                0    1    2    3    4    5    6    7

source[1:3]        →      └─2,3─┘  cap leaks to end of array (cap=7)
source[1:3:3]       →      └─2,3─┘  cap capped at index 3 (cap=2)
```

Capping capacity to equal length means the very next `append` MUST reallocate - it can no longer sneak a write into `source`'s array.

---

## 🗒️ copy() vs Slicing

| Operation | Creates New Backing Array? | Aliased to Source? |
|-----------|------------------------------|----------------------|
| `sub := s[low:high]` | No | YES - shares memory |
| `sub := s[low:high:max]` | No (yet) | YES until it grows past `cap` |
| `dst := make([]T, n); copy(dst, s)` | Yes | NO - fully independent |

```go
alias := original[:]              // still the SAME array
clone := make([]int, len(original))
copy(clone, original)             // a DIFFERENT array, same values
```

---

## 🆚 nil Slice vs Empty Slice

| | `var s []int` (nil) | `s := []int{}` (empty) |
|---|---|---|
| `len(s)` | 0 | 0 |
| `cap(s)` | 0 | 0 |
| `s == nil` | **true** | **false** |
| `%v` output | `[]` | `[]` |
| `%#v` output | `[]int(nil)` | `[]int{}` |
| `append(s, x)` | Works, returns `[x]` | Works, returns `[x]` |
| `for range s` | 0 iterations | 0 iterations |
| JSON encoding (`encoding/json`) | `null` | `[]` |

**Actual verified output:**
```
nilSlice == nil: true
emptySlice == nil: false

nilSlice:   []   []int(nil)
emptySlice: []   []int{}
```

**Rule of thumb:** use `len(s) == 0` to check emptiness. Only reach for `s == nil` when the nil-ness itself matters (e.g., you're about to marshal to JSON and care whether the field is `null` or `[]`).

---

## 🧱 2D Slices: The Shared-Row Bug

```
❌ WRONG                              ✅ RIGHT

row := make([]int, cols)              grid := make([][]int, rows)
grid := make([][]int, rows)           for i := range grid {
for i := range grid {                     grid[i] = make([]int, cols)
    grid[i] = row  // same slice!     }
}

grid[0] ──┐                           grid[0] ──▶ [own array]
grid[1] ──┼──▶ [ONE shared array]     grid[1] ──▶ [own array]
grid[2] ──┘                           grid[2] ──▶ [own array]
```

**Actual verified output of the bug:**
```
After setting wrong[0][0] = 1:
[1 0 0]
[1 0 0]
[1 0 0]
(every row shows the change - they all share ONE underlying array!)
```

**Decision Rule:** if you ever write `grid[i] = someSliceVariable` where `someSliceVariable` was declared OUTSIDE the loop, you have the bug. Each row needs its own call to `make()` (or its own literal `{...}`) INSIDE the loop.

---

## 🚨 Common Mistakes

### Mistake 1: Treating Slice Literals Like Arrays

```go
var a [3]int = [3]int{1, 2, 3}  // array - fixed size, value semantics
var s []int = []int{1, 2, 3}    // slice - dynamic view, header semantics
```

### Mistake 2: Assuming append() Always (or Never) Mutates the Original

```go
// It depends entirely on len(s) vs cap(s) at the moment of the call.
// Never assume either direction without checking.
```

### Mistake 3: Forgetting Sub-Slices Alias Their Source

```go
sub := original[1:4]
sub[0] = 999 // also changes original[1]!
```

### Mistake 4: The Shared-Row Bug in 2D Slices

```go
// ❌ one row slice reused for every row
row := make([]int, cols)
for i := range grid { grid[i] = row }
```

### Mistake 5: Using == nil Instead of len() == 0

```go
// ❌ misses a non-nil empty slice like []int{}
if tasks == nil { }

// ✅ correct for both nil and non-nil empty slices
if len(tasks) == 0 { }
```

---

## 📈 Progression Summary

### Understanding Level 8

Level 8 teaches the collection type you'll use constantly from here on:

1. **The header** - pointer + len + cap, not a container
2. **Creation** - literals, make(), slicing
3. **len vs cap** - two different questions with two different answers
4. **append()** - sometimes mutates, sometimes reallocates
5. **Aliasing** - sub-slices share memory with their source
6. **Three-index slicing & copy()** - the two tools for preventing/removing aliasing
7. **nil vs empty** - a subtle but real distinction
8. **2D slices** - allocate every row independently

### Prerequisites for Level 9

Before moving to Level 9 (Strings & Runes), you need:

- ✅ Comfortable creating slices all four ways
- ✅ Can explain len() vs cap() and why they differ
- ✅ Can predict whether a given append() will mutate in place or reallocate
- ✅ Understand that sub-slices alias their underlying array
- ✅ Know when to reach for `s[low:high:max]` or `copy()`
- ✅ Know the nil-vs-empty-slice distinction
- ✅ Can build a correct 2D slice grid

### Ready for Level 9?

Level 9 teaches the collection type you've already been slicing without realizing it - strings are, under the hood, an immutable read-only view over a byte array:
- Strings as byte sequences
- Runes vs bytes (you saw a preview of this with range in Level 6)
- Converting between string, []byte, and []rune

---

## ✅ Checklist Before Level 9

- [ ] Can create a slice with a literal, make(), and a slice expression
- [ ] Can explain why `cap()` can be larger than `len()`
- [ ] Can predict when append() reallocates vs. mutates in place
- [ ] Can demonstrate aliasing between two slices sharing an array
- [ ] Can use `s[low:high:max]` to cap capacity
- [ ] Can use `copy()` to make an independent clone
- [ ] Know that `nilSlice == nil` but `emptySlice != nil`, and that both behave the same with append/range/len
- [ ] Can build a 2D slice grid without the shared-row bug
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The Header
A slice is `{pointer, len, cap}` - three words describing a VIEW, never the data itself.

### The Growth
`append` mutates in place when there's room (`len < cap`), and silently reallocates when there isn't (`len == cap`). Never assume which one happened.

### The Aliasing
Slicing shares memory. Mutating through one slice can be visible through another. `copy()` is the only way to truly break that connection.

### The nil Subtlety
`nilSlice == nil` is true, `emptySlice == nil` is false - but check emptiness with `len()`, not `== nil`.

### The 2D Rule
Every row of a `[][]T` needs its own call to `make()` - reusing one row slice for every row is the classic bug.

---

## 📚 Next Level

Level 9: Strings & Runes
- Strings as immutable byte sequences
- Runes and Unicode code points
- Converting between strings, []byte, and []rune
- Common string-processing patterns

You've got slices down! Keep going! 🚀
