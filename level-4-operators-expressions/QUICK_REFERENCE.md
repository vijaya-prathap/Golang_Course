# Level 4: Quick Reference Card

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
    a, b := 10, 3

    fmt.Println(a+b, a-b, a*b, a/b, a%b) // 13 7 30 3 1
    fmt.Println(a > b && b > 0)          // true
    fmt.Println(a&b, a|b, a^b)           // 2 11 9

    a += 5
    fmt.Println(a) // 15
}
EOF

# Run
go run main.go
```

---

## ➕ Arithmetic

```go
a + b   // addition
a - b   // subtraction
a * b   // multiplication
a / b   // division (int / int truncates!)
a % b   // modulo (remainder)
-a      // unary negation
+a      // unary affirmation

float64(a) / b   // use this to get a decimal result
```

---

## 🔀 Comparison (all return bool)

```go
a == b   // equal
a != b   // not equal
a > b    // greater than
a >= b   // greater than or equal
a < b    // less than
a <= b   // less than or equal

"abc" < "xyz"   // lexicographic string comparison
```

---

## 🔗 Logical (work on bool)

```go
a && b   // AND — both must be true
a || b   // OR — at least one must be true
!a       // NOT — inverts

// Short-circuit:
// a && b -> if a is false, b is never evaluated
// a || b -> if a is true, b is never evaluated
```

---

## 🔢 Bitwise (work on integers)

```go
a & b    // AND (bit is 1 if both are 1)
a | b    // OR  (bit is 1 if either is 1)
a ^ b    // XOR (bit is 1 if bits differ)
^a       // NOT (unary, flips all bits)
a &^ b   // AND NOT / bit clear (clears bits set in b)

a << n   // left shift  (multiply by 2^n)
a >> n   // right shift (divide by 2^n, positive numbers)
```

### Bit Flags Pattern

```go
const (
    Read   = 1 << iota // 1
    Write              // 2
    Execute            // 4
)

perms := Read | Write         // grant multiple
perms &^= Write                // revoke one
hasRead := perms&Read == Read  // check one
```

---

## 📝 Assignment

```go
x = 5     // assign
x += 3    // x = x + 3
x -= 3    // x = x - 3
x *= 3    // x = x * 3
x /= 3    // x = x / 3
x %= 3    // x = x % 3
x &= 3    // x = x & 3
x |= 3    // x = x | 3
x ^= 3    // x = x ^ 3
x <<= 1   // x = x << 1
x >>= 1   // x = x >> 1

x++       // x = x + 1  (statement only, not an expression)
x--       // x = x - 1  (statement only, not an expression)
```

---

## 📊 Operator Precedence (highest → lowest)

```
1. ()
2. ++  --
3. *  /  %  &  <<  >>
4. +  -  |  ^
5. ==  !=  <  >  <=  >=
6. &&
7. ||
```

```go
2 + 3 * 4      // 14, not 20
(2 + 3) * 4    // 20

true || false && false    // true  (&& first)
(true || false) && false  // false (parens override)
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Precedence surprise | `2 + 3 * 4 // expecting 20` | `(2 + 3) * 4` |
| Integer division loss | `pct := 1 / 3 * 100 // 0` | `pct := int(float64(1) / 3 * 100)` |
| `=` vs `==` in condition | `if x = 5 {}` | `if x == 5 {}` |
| `++` as expression | `y := x++` | `x++; y := x` |
| Unsafe short-circuit order | `if a/b > 10 && b != 0 {}` | `if b != 0 && a/b > 10 {}` |
| Mixing numeric types | `var f float64 = i` | `var f float64 = float64(i)` |

---

## 🎓 Before Next Level

Can you:
- [ ] Use all arithmetic operators, including modulo?
- [ ] Explain why `10 / 3` is `3`, not `3.33`?
- [ ] Combine conditions with &&, ||, !?
- [ ] Build a bit-flag system with &, |, ^, <<?
- [ ] Use compound assignment operators?
- [ ] Predict precedence without running the code?
- [ ] Explain short-circuit evaluation?

If YES → You're ready for Level 5!

---

## 📚 Next Level

Level 5: Conditions (if/else/switch)
- `if` / `else if` / `else`
- `switch` statements
- Making decisions with expressions

You've got operators down! 💪
