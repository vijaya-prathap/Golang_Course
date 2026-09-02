# Level 6: Loops - Complete Guide

## Introduction

Welcome to Level 6! You've mastered conditions (Level 5). Now it's time to learn **loops** - how to repeat code without copy-pasting it.

Go has exactly **one** loop keyword: `for`. Unlike most languages, there's no separate `while` or `do-while` - `for` covers all of them through different forms.

---

## Table of Contents

1. [The Three-Component for Loop](#the-three-component-for-loop)
2. [The while-Style for Loop](#the-while-style-for-loop)
3. [The Infinite for Loop](#the-infinite-for-loop)
4. [break and continue](#break-and-continue)
5. [Looping Over Collections With range](#looping-over-collections-with-range)
6. [Labeled break and continue](#labeled-break-and-continue)
7. [Nested Loops](#nested-loops)
8. [Common Loop Patterns](#common-loop-patterns)
9. [Best Practices](#best-practices)
10. [Common Mistakes](#common-mistakes)

---

## The Three-Component for Loop

The classic C-style loop: init, condition, post.

```go
for i := 0; i < 5; i++ {
    fmt.Println(i)
}
// prints 0, 1, 2, 3, 4
```

### Anatomy

```go
for <init>; <condition>; <post> {
    // body
}

init:      runs once, before the loop starts (i := 0)
condition: checked before EVERY iteration; loop stops when false (i < 5)
post:      runs after EVERY iteration, before the next check (i++)
```

Just like `if`, the opening brace must be on the same line, and no parentheses are needed around the clauses.

---

## The while-Style for Loop

Drop the init and post clauses, keep only the condition - this is Go's `while`.

```go
count := 0
for count < 5 {
    fmt.Println(count)
    count++
}
// prints 0, 1, 2, 3, 4
```

There is no separate `while` keyword in Go - this IS the while loop.

---

## The Infinite for Loop

Drop all three clauses for an infinite loop. You must `break` (or `return`) to exit.

```go
count := 0
for {
    if count >= 5 {
        break
    }
    fmt.Println(count)
    count++
}
```

Common uses: server request-handling loops, retry loops, game loops - anywhere the exit condition is checked inside the body rather than at the top.

---

## break and continue

`break` exits the loop immediately. `continue` skips to the next iteration.

```go
for i := 0; i < 10; i++ {
    if i == 5 {
        break // stop the loop entirely
    }
    fmt.Println(i)
}
// prints 0, 1, 2, 3, 4
```

```go
for i := 0; i < 10; i++ {
    if i%2 != 0 {
        continue // skip odd numbers
    }
    fmt.Println(i)
}
// prints 0, 2, 4, 6, 8
```

---

## Looping Over Collections With range

`range` iterates over arrays, slices, strings, maps, and channels.

### Slices and Arrays

```go
fruits := []string{"apple", "banana", "cherry"}

for index, fruit := range fruits {
    fmt.Println(index, fruit)
}
// 0 apple
// 1 banana
// 2 cherry
```

### Ignoring the Index or Value

```go
// Only need the value? Ignore the index with _
for _, fruit := range fruits {
    fmt.Println(fruit)
}

// Only need the index?
for i := range fruits {
    fmt.Println(i)
}
```

### Strings (range gives you runes, not bytes!)

```go
for i, r := range "héllo" {
    fmt.Printf("index=%d rune=%c\n", i, r)
}
// index=0 rune=h
// index=1 rune=é     (takes 2 bytes, so next index jumps to 3)
// index=3 rune=l
// index=4 rune=l
// index=5 rune=o
```

### Maps (order is NOT guaranteed!)

```go
scores := map[string]int{"Alice": 95, "Bob": 87}

for name, score := range scores {
    fmt.Println(name, score)
}
// Order varies between runs - Go randomizes map iteration order on purpose!
```

---

## Labeled break and continue

By default, `break`/`continue` only affect the innermost loop. Use a **label** to target an outer loop.

```go
outer:
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if j == 1 {
                continue outer // skip to the next i, not just the next j
            }
            fmt.Println(i, j)
        }
    }
```

```go
search:
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if i == 1 && j == 1 {
                break search // exits BOTH loops immediately
            }
            fmt.Println(i, j)
        }
    }
```

---

## Nested Loops

Loops can nest inside each other, commonly for grid/table-shaped data.

```go
for row := 1; row <= 3; row++ {
    for col := 1; col <= 3; col++ {
        fmt.Printf("%d ", row*col)
    }
    fmt.Println()
}
// 1 2 3
// 2 4 6
// 3 6 9
```

---

## Common Loop Patterns

### Accumulating a Total

```go
sum := 0
for _, n := range []int{1, 2, 3, 4, 5} {
    sum += n
}
```

### Finding an Item

```go
found := false
for _, name := range names {
    if name == "Alice" {
        found = true
        break
    }
}
```

### Building a New Slice

```go
var evens []int
for _, n := range numbers {
    if n%2 == 0 {
        evens = append(evens, n)
    }
}
```

### Counting Occurrences

```go
counts := make(map[string]int)
for _, word := range words {
    counts[word]++
}
```

---

## Best Practices

### 1. Prefer range Over Manual Indexing

```go
// ✅ Good
for i, v := range items {
    fmt.Println(i, v)
}

// ❌ Unnecessary
for i := 0; i < len(items); i++ {
    fmt.Println(i, items[i])
}
```

### 2. Use _ to Ignore Unused Loop Variables

```go
// ✅ Good - index unused, explicitly ignored
for _, v := range items {
    process(v)
}

// ❌ Compile error - unused variable
for i, v := range items {
    process(v)
}
```

### 3. Don't Rely on Map Iteration Order

If you need a predictable order, sort the keys first:

```go
keys := make([]string, 0, len(myMap))
for k := range myMap {
    keys = append(keys, k)
}
sort.Strings(keys)
for _, k := range keys {
    fmt.Println(k, myMap[k])
}
```

### 4. Use Labeled Breaks Sparingly

They're useful for exiting nested loops early, but overuse hurts readability. Consider extracting the nested loop into its own function with a normal `return` instead.

### 5. Prefer the while-Style Loop Over a Fake do-while

Go has no `do-while`. Don't fake it awkwardly - restructure with a plain condition or an infinite loop with an internal check.

---

## Common Mistakes

### Mistake 1: Off-by-One Errors

```go
// ❌ WRONG - misses the last index (skips len(items)-1)
for i := 0; i < len(items)-1; i++ { }

// ✅ RIGHT
for i := 0; i < len(items); i++ { }
```

### Mistake 2: Infinite Loop From a Forgotten Post/Condition

```go
// ❌ WRONG - i never changes, loops forever
for i := 0; i < 5; {
    fmt.Println(i)
}

// ✅ RIGHT
for i := 0; i < 5; i++ {
    fmt.Println(i)
}
```

### Mistake 3: Modifying a Slice While Ranging Over It

```go
// ❌ RISKY - appending during range can behave unexpectedly
for _, v := range items {
    if v == 0 {
        items = append(items, v) // don't do this
    }
}

// ✅ RIGHT - build a new slice instead
var result []int
for _, v := range items {
    result = append(result, v)
}
```

### Mistake 4: Expecting Map Iteration to Be Ordered

```go
// ❌ WRONG ASSUMPTION - order is randomized on each run
for k := range myMap {
    fmt.Println(k) // don't assume any particular order
}

// ✅ RIGHT - sort keys first if order matters
```

### Mistake 5: Unlabeled break Inside Nested Loops

```go
// ❌ WRONG - only breaks the inner loop, outer loop continues
for i := 0; i < 3; i++ {
    for j := 0; j < 3; j++ {
        if j == 1 {
            break // only exits the j loop!
        }
    }
}

// ✅ RIGHT - use a label to break the outer loop too
outer:
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if j == 1 {
                break outer
            }
        }
    }
```

### Mistake 6: Assuming range Gives Bytes for Strings

```go
s := "héllo"
// ❌ WRONG ASSUMPTION - indices skip around because é is 2 bytes
for i, r := range s {
    fmt.Println(i, r) // index isn't 0,1,2,3,4 - it's 0,1,3,4,5
}

// If you need byte-by-byte, index the string directly: s[i]
```

---

## Summary

**Three for Loop Forms:**
- `for init; cond; post { }` - classic three-component loop
- `for cond { }` - while-style loop
- `for { }` - infinite loop (must break to exit)

**Loop Control:**
- `break` - exit the loop immediately
- `continue` - skip to the next iteration
- Labels (`outer:`) let break/continue target an outer loop

**range:**
- Iterates over arrays, slices, strings, maps, and channels
- Strings: range gives you runes (and byte-offset indices), not raw bytes
- Maps: iteration order is randomized on purpose - sort keys if order matters
- Use `_` to explicitly ignore the index or value you don't need

**Best Practices:**
- Prefer `range` over manual indexing when you don't need low-level control
- Don't rely on map iteration order
- Use labels sparingly; consider extracting nested loops into functions

---

## Next Steps

You now understand:
- ✅ All three for loop forms
- ✅ break, continue, and labeled loop control
- ✅ range over slices, arrays, strings, and maps
- ✅ Common accumulation/filtering/counting patterns

**Next level:** Level 7 - Arrays
- Fixed-size arrays in depth
- Array vs slice semantics
- Multi-dimensional arrays
- When to use arrays vs slices

You're building the tools to process real data! Keep going! 🚀
