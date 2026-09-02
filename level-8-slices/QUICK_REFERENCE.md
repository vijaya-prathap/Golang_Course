# Level 8: Quick Reference Card

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
    nums := []int{1, 2, 3}
    nums = append(nums, 4, 5)
    fmt.Println(nums, "len:", len(nums), "cap:", cap(nums))

    sub := nums[1:3]
    sub[0] = 999
    fmt.Println(nums) // sub aliases nums - this mutates nums too!
}
EOF

# Run
go run main.go
```

---

## 📋 Creating Slices

```go
nums := []int{1, 2, 3}      // literal - len=3 cap=3
zeros := make([]int, 4)     // make with length - len=4 cap=4
buf := make([]int, 2, 6)    // make with length AND capacity - len=2 cap=6
sub := arr[1:4]             // slicing an array or slice
```

---

## 📏 len() vs cap()

```go
len(s) // elements currently in the slice
cap(s) // room before the NEXT append must reallocate,
       // measured from where s starts to the end of its underlying array
```

---

## ➕ append() Semantics

```go
len(s) < cap(s)   // append WRITES INTO the same underlying array
len(s) == cap(s)  // append ALLOCATES A NEW array and copies everything

// Never assume which one happened without checking -
// always use the slice append() RETURNS, never the old variable
s = append(s, v)
```

---

## 🔗 Aliasing (Slicing Shares Memory)

```go
sub := original[1:4]
sub[0] = 999          // ALSO changes original[1] - same backing array!
```

```go
s[low:high:max]       // three-index (full) slice expression
                       // len = high-low, cap = max-low
                       // caps capacity so the next append() can't leak
                       // into the original array
```

---

## 📋 copy() - True Independence

```go
dst := make([]int, len(src))
n := copy(dst, src)   // n = number of elements actually copied
                       // dst is now a SEPARATE array - no aliasing
```

---

## 🆚 nil vs Empty Slice

```go
var nilSlice []int     // nilSlice == nil is TRUE
emptySlice := []int{}  // emptySlice == nil is FALSE

// Both: len == 0, cap == 0, work fine with append/range/len
// Check emptiness with len(s) == 0, not s == nil
```

---

## 🧱 2D Slices

```go
// ✅ RIGHT - allocate each row separately
grid := make([][]int, rows)
for i := range grid {
    grid[i] = make([]int, cols)
}

// ❌ WRONG - every row is the SAME slice (shared-row bug)
row := make([]int, cols)
for i := range grid {
    grid[i] = row
}
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Confusing slices with arrays | assuming `[]int{1,2}` behaves like `[2]int{1,2}` | remember: no number = slice, has a number = array |
| Assuming append() never mutates | ignoring that `len < cap` writes in place | check `len` vs `cap` before relying on either behavior |
| Assuming append() always mutates | ignoring reallocation when `len == cap` | always use the slice append() returns |
| Forgetting sub-slices alias their source | mutating a sub-slice and being surprised the original changed | expect aliasing; use `copy()` when you need independence |
| Shared-row bug | `row := make(...); grid[i] = row` in a loop | `grid[i] = make([]int, cols)` inside the loop |
| Checking emptiness with `== nil` | `if s == nil` (misses non-nil empty slices) | `if len(s) == 0` |

---

## 🎓 Before Next Level

Can you:
- [ ] Create a slice all four ways (literal, make with len, make with len+cap, slicing)?
- [ ] Explain why cap() can exceed len()?
- [ ] Predict whether a given append() call will mutate in place or reallocate?
- [ ] Demonstrate that sub-slices alias their underlying array?
- [ ] Use `s[low:high:max]` and `copy()` to control or remove aliasing?
- [ ] Explain the nil-vs-empty-slice distinction?
- [ ] Build a correct 2D slice grid without the shared-row bug?

If YES → You're ready for Level 9!

---

## 📚 Next Level

Level 9: Strings & Runes
- Strings as immutable byte sequences
- Runes and Unicode code points
- Converting between strings, []byte, and []rune

You've got slices down! 💪
