# Level 6: Quick Reference Card

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
    for i := 0; i < 5; i++ {
        fmt.Println(i)
    }

    fruits := []string{"apple", "banana", "cherry"}
    for i, f := range fruits {
        fmt.Println(i, f)
    }
}
EOF

# Run
go run main.go
```

---

## 📋 The Three for Loop Shapes

```go
// Three-component (classic for)
for i := 0; i < 5; i++ {
    // ...
}

// while-style (condition only)
for count > 0 {
    // ...
}

// infinite (must break to exit)
for {
    if done {
        break
    }
}
```

---

## 🔁 break and continue

```go
break      // exit the loop immediately
continue   // skip to the next iteration
```

```go
for i := 0; i < 10; i++ {
    if i == 5 {
        break
    }
    if i%2 != 0 {
        continue
    }
    fmt.Println(i)
}
```

---

## 🎁 range Syntax

```go
for i, v := range slice { }   // index + value
for _, v := range slice { }   // value only
for i := range slice { }      // index only
for range slice { }           // neither, just iterate

for i, r := range someString { }   // byte-offset index + rune

for k, v := range someMap { }      // key + value (order randomized!)
for k := range someMap { }         // keys only
```

---

## 🏷️ Labeled break / continue

```go
outer:
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if j == 1 {
                break outer       // exits BOTH loops
                // continue outer would skip to next i instead
            }
        }
    }
```

---

## 🧩 Common Patterns

```go
// Accumulate
sum := 0
for _, n := range nums { sum += n }

// Find
for _, v := range items {
    if v == target { break }
}

// Filter
var result []int
for _, n := range nums {
    if n%2 == 0 { result = append(result, n) }
}

// Count
counts := make(map[string]int)
for _, w := range words { counts[w]++ }

// Sorted map iteration
keys := make([]string, 0, len(m))
for k := range m { keys = append(keys, k) }
sort.Strings(keys)
for _, k := range keys { fmt.Println(k, m[k]) }
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Off-by-one | `i < len(items)-1` | `i < len(items)` |
| Forgotten post | `for i:=0; i<5; { }` | `for i:=0; i<5; i++ { }` |
| Assuming map order | relying on iteration order | sort keys first |
| Unlabeled break in nested loop | `break` (only exits inner) | `break outer` (labeled) |
| Treating string index as char position | assuming 1 index = 1 char | remember multi-byte runes shift indices |

---

## 🎓 Before Next Level

Can you:
- [ ] Write all three for loop forms from memory?
- [ ] Explain the difference between break and continue?
- [ ] Range over a slice, string, and map correctly?
- [ ] Explain why range over a string gives runes, not bytes?
- [ ] Use a labeled break to exit nested loops?
- [ ] Write accumulate/find/filter/count patterns without hesitation?

If YES → You're ready for Level 7!

---

## 📚 Next Level

Level 7: Arrays
- Fixed-size arrays in depth
- Array vs slice semantics
- Multi-dimensional arrays

You've got loops down! 💪
