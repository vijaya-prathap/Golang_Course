# Level 6: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: Loops
```
Day 1:  Three-component for loop
Day 2:  while-style and infinite for loops
Day 3:  break and continue
Day 4:  range over slices and arrays
Day 5:  range over strings (runes) and maps
Day 6:  Labeled break/continue
Day 7:  Nested loops and common patterns
```

### Week 2: Practice & Application
```
Day 1:  Exercises 1-3 (for forms, break)
Day 2:  Exercises 4-6 (continue, range slices, range strings)
Day 3:  Exercises 7-8 (range maps, labeled loops)
Day 4:  Exercises 9-10 (patterns, comprehensive practice)
Day 5:  Bonus challenges
Day 6-7: Practice & review
```

---

## 🎯 The Three for Loop Shapes

| Shape | Syntax | Like... |
|-------|--------|---------|
| **Three-component** | `for i := 0; i < n; i++ { }` | classic C `for` |
| **while-style** | `for condition { }` | `while` in other languages |
| **infinite** | `for { }` | `while (true)` / `loop` |

```
for <init>; <condition>; <post> { }
    ├─ init:      runs ONCE, before anything else
    ├─ condition: checked BEFORE every iteration
    └─ post:      runs AFTER every iteration

for <condition> { }
    └─ just the condition - Go's while loop

for { }
    └─ no clauses - runs forever until break/return
```

---

## 🔁 Loop Control Flow

```
break     → exits the loop immediately, skips everything after it
continue  → skips the REST of this iteration, jumps to the next one

for i := 0; i < 5; i++ {
    if i == 3 {
        break        // stops here: 0, 1, 2
    }
    fmt.Println(i)
}

for i := 0; i < 5; i++ {
    if i == 3 {
        continue     // skips 3: 0, 1, 2, 4
    }
    fmt.Println(i)
}
```

---

## 🎁 range Cheat Sheet

| Collection | `for a, b := range x` gives you |
|------------|----------------------------------|
| `[]T` (slice) | index (`int`), value (`T`) |
| `[N]T` (array) | index (`int`), value (`T`) |
| `string` | byte-offset index (`int`), rune (`rune`) |
| `map[K]V` | key (`K`), value (`V`) — **order is randomized** |
| `chan T` | value only (`T`) — single-variable form |

```go
for i, v := range slice { }     // both index and value
for _, v := range slice { }     // value only
for i := range slice { }        // index only (also works to just count)
for range slice { }             // neither - just loop len(slice) times
```

---

## 🔤 Strings: Rune Iteration Explained

```
s := "héllo"

Bytes:  h  é(2 bytes)  l  l  o     → len(s) == 6
Runes:  h  é            l  l  o     → 5 characters

for i, r := range s {
    // i is the BYTE offset where rune r starts
    // r is the decoded Unicode code point (rune / int32)
}

Result:
  index=0 rune=h
  index=1 rune=é   (occupies bytes 1-2)
  index=3 rune=l   (index jumped from 1 to 3!)
  index=4 rune=l
  index=5 rune=o
```

**Rule:** if you need raw bytes, index the string directly (`s[i]`). If you need characters, use `range`.

---

## 🗺️ Maps: Order Is Randomized

```
Go DELIBERATELY randomizes map iteration order between runs.
This prevents code from accidentally depending on an order
that was never guaranteed.

Need a predictable order? Sort the keys first:

    keys := make([]string, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    for _, k := range keys {
        fmt.Println(k, m[k])
    }
```

---

## 🏷️ Labeled break / continue

```
Without a label, break/continue only affect the INNERMOST loop.

outer:
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if j == 1 {
                break outer     // exits BOTH loops
                // continue outer would instead skip to the next i
            }
        }
    }
```

### Decision Tree

```
Need to stop/skip from inside a nested loop?
├─ Only need to affect the innermost loop?
│  └─ plain break / continue
└─ Need to affect an OUTER loop?
   └─ label the outer loop, use break <label> / continue <label>
```

---

## 🧩 Common Loop Patterns

```go
// Accumulate
sum := 0
for _, n := range nums { sum += n }

// Find (stop early once found)
found := false
for _, v := range items {
    if v == target {
        found = true
        break
    }
}

// Filter (build a new slice)
var result []int
for _, n := range nums {
    if n%2 == 0 {
        result = append(result, n)
    }
}

// Count (frequency map)
counts := make(map[string]int)
for _, w := range words {
    counts[w]++
}
```

---

## 🚨 Common Mistakes

### Mistake 1: Off-by-One

```go
// ❌ WRONG - misses the last valid index
for i := 0; i < len(items)-1; i++ { }

// ✅ RIGHT
for i := 0; i < len(items); i++ { }
```

### Mistake 2: Forgetting to Update the Loop Variable

```go
// ❌ WRONG - infinite loop, i never changes
for i := 0; i < 5; {
    fmt.Println(i)
}

// ✅ RIGHT
for i := 0; i < 5; i++ { }
```

### Mistake 3: Assuming Map Order Is Stable

```go
// ❌ WRONG ASSUMPTION
for k := range m {
    fmt.Println(k) // order varies run to run
}
```

### Mistake 4: Unlabeled break Inside Nested Loops

```go
// ❌ WRONG - only exits the inner loop
for i := 0; i < 3; i++ {
    for j := 0; j < 3; j++ {
        if j == 1 { break }
    }
}

// ✅ RIGHT - label it if you mean to exit both
outer:
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if j == 1 { break outer }
        }
    }
```

### Mistake 5: Treating String Indices as Character Positions

```go
s := "héllo"
// ❌ WRONG ASSUMPTION - index 2 is NOT the 3rd character here
// (because é occupies 2 bytes)
```

---

## 📈 Progression Summary

### Understanding Level 6

Level 6 teaches how to repeat work using conditions (Level 5) as guards:

1. **Three loop shapes** — counted, conditional, infinite
2. **break/continue** — early exit and skip
3. **range** — the idiomatic way to iterate collections
4. **Labels** — controlling outer loops from nested ones
5. **Patterns** — accumulate, find, filter, count

### Prerequisites for Level 7

Before moving to Level 7 (Arrays), you need:

- ✅ Comfortable with all three for loop forms
- ✅ Comfortable with break, continue, and labels
- ✅ Can range over slices, strings, and maps correctly
- ✅ Understand rune vs byte iteration for strings
- ✅ Understand map iteration order is randomized

### Ready for Level 7?

Level 7 teaches the fixed-size collection type you've already been ranging over informally:
- Arrays in depth (fixed size, value semantics)
- Array vs slice differences
- Multi-dimensional arrays

---

## ✅ Checklist Before Level 7

- [ ] Can write all three for loop forms
- [ ] Can use break and continue correctly
- [ ] Can range over a slice, ignoring index or value as needed
- [ ] Understand why range over a string gives runes, not bytes
- [ ] Know map iteration order is randomized and how to sort keys
- [ ] Can use a labeled break/continue in nested loops
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The One Keyword
`for` is Go's only loop keyword — three shapes, one keyword.

### The Control
`break` exits, `continue` skips — both target only the innermost loop unless labeled.

### The Iteration
`range` is idiomatic for slices, arrays, strings, and maps — but strings give runes, and maps give randomized order.

### The Patterns
Accumulate, find, filter, count — nearly every loop you write is one of these four patterns (or a combination).

---

## 📚 Next Level

Level 7: Arrays
- Fixed-size arrays in depth
- Array vs slice semantics
- Multi-dimensional arrays
- When to use arrays vs slices

You've got loops down! Keep going! 🚀
