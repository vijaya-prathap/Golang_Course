# Level 5: Quick Reference Card

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
    age := 20

    if age >= 18 {
        fmt.Println("Adult")
    } else {
        fmt.Println("Minor")
    }

    switch {
    case age < 13:
        fmt.Println("Child")
    case age < 20:
        fmt.Println("Teen")
    default:
        fmt.Println("Adult")
    }
}
EOF

# Run
go run main.go
```

---

## 📋 if / else Syntax

```go
if condition {
    // ...
}

if condition {
    // ...
} else {
    // ...
}

if condition1 {
    // ...
} else if condition2 {
    // ...
} else {
    // ...
}

// Short init - scoped to if/else only
if x := compute(); x > 0 {
    // use x
} else {
    // use x
}
// x not visible here
```

### Rule
Opening brace `{` MUST be on the same line as `if` / `else if` / `else`.

---

## 🔀 switch Syntax

```go
// Basic
switch value {
case a:
    // ...
case b, c:      // OR: matches b or c
    // ...
default:
    // ...
}

// With short init
switch x := compute(); x {
case 1:
    // ...
}

// Conditionless (like if/else if)
switch {
case x < 10:
    // ...
case x < 20:
    // ...
default:
    // ...
}

// fallthrough - explicitly run next case too
switch n {
case 1:
    fmt.Println("one")
    fallthrough
case 2:
    fmt.Println("two")
}
```

### Rule
Go's switch does **not** fall through by default — each case auto-breaks. Use `fallthrough` to opt in.

---

## 🎯 Comma-OK Idiom

```go
// Map lookup
if value, ok := myMap[key]; ok {
    // key exists, value is valid
}

// Type assertion
if s, ok := x.(string); ok {
    // x was a string, s holds it
}

// Error checking
if value, err := strconv.Atoi(s); err == nil {
    // parsed successfully
}
```

---

## 🪆 Flattening Nested Conditions

```go
// ❌ Nested
if a {
    if b {
        doWork()
    }
}

// ✅ Early returns (guard clauses)
if !a {
    return
}
if !b {
    return
}
doWork()
```

---

## 🌳 if vs switch Decision Guide

```
One or two conditions?
├─ YES → use if / if-else
└─ NO  → many conditions?
    ├─ All checking the SAME value?
    │  └─ use switch value { case a: ... }
    └─ Checking DIFFERENT conditions/ranges?
       └─ use switch { case cond: ... }
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Wrong threshold order | `if s>=70{"C"} else if s>=90{"A"}` | `if s>=90{"A"} else if s>=70{"C"}` |
| Expecting fallthrough | relying on switch to fall through | add `fallthrough` explicitly |
| Brace on new line | `if x > 5\n{ }` | `if x > 5 { }` |
| Using = instead of == | `if x = 5 {}` | `if x == 5 {}` |
| Using short-init var outside scope | `if v,ok:=m[k];ok{}; fmt.Println(v)` | keep usage inside the if/else |
| Deep nesting | `if a { if b { if c {} } }` | early returns / guard clauses |

---

## 🎓 Before Next Level

Can you:
- [ ] Write an if/else-if chain with correctly ordered thresholds?
- [ ] Use the short-init form for if and switch?
- [ ] Write both a value-based and conditionless switch?
- [ ] Explain why Go's switch doesn't fall through by default?
- [ ] Use fallthrough correctly when you actually need it?
- [ ] Flatten nested conditions with early returns?

If YES → You're ready for Level 6!

---

## 📚 Next Level

Level 6: Loops
- `for` loops (Go's only loop keyword)
- while-style and infinite loops
- `break` and `continue`
- Looping over collections

You've got conditions down! 💪
