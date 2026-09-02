# Level 4: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: Operators
```
Day 1:  Arithmetic operators & integer division
Day 2:  Comparison operators
Day 3:  Logical operators
Day 4:  Bitwise operators
Day 5:  Assignment operators
Day 6:  Operator precedence
Day 7:  Short-circuit evaluation
```

### Week 2: Practice & Application
```
Day 1:  Exercises 1-3 (Arithmetic, division, comparison)
Day 2:  Exercises 4-6 (Logical, bitwise, assignment)
Day 3:  Exercises 7-8 (Precedence, short-circuit)
Day 4:  Exercises 9-10 (Expressions, comprehensive practice)
Day 5:  Bonus challenges
Day 6-7: Practice & review
```

---

## 🎯 Operator Categories At A Glance

| Category | Operators | Returns |
|----------|-----------|---------|
| **Arithmetic** | `+ - * / %` | Same numeric type |
| **Comparison** | `== != < > <= >=` | `bool` |
| **Logical** | `&& \|\| !` | `bool` |
| **Bitwise** | `& \| ^ ^(unary) << >>` | Same integer type |
| **Assignment** | `= += -= *= /= %= &= \|= ^= <<= >>= ++ --` | (no value; statement) |

---

## ➕ Arithmetic Operators

```
a := 10
b := 3

a + b  → 13   addition
a - b  → 7    subtraction
a * b  → 30   multiplication
a / b  → 3    division (integer division truncates!)
a % b  → 1    modulo (remainder)

-a     → -10  unary negation
+a     → 10   unary affirmation (rarely used)
```

### Integer Division Truncation

```
int / int  → int (fraction is DISCARDED, not rounded)

10 / 3   → 3     (not 3.333)
10 / 4   → 2     (not 2.5)
-7 / 2   → -3    (truncates toward zero)

To get a decimal result, convert BEFORE dividing:
float64(10) / 3   → 3.3333...
```

### Modulo Use Cases

```
Even/odd check:      n % 2 == 0
Cyclic wrap (clock):  hour % 24
Grouping into buckets: index % bucketCount
```

---

## 🔀 Comparison Operators

```
a := 10
b := 3

a == b  → false   equal
a != b  → true    not equal
a > b   → true    greater than
a >= b  → true    greater than or equal
a < b   → false   less than
a <= b  → false   less than or equal
```

### Rules

```
✅ Both sides must be the same type
   var i int = 5
   var f float64 = 5.0
   i == f              ❌ ERROR (can't compare int and float64 directly)
   float64(i) == f     ✅ RIGHT (convert first)

✅ Strings compare lexicographically   "abc" < "xyz"  → true
✅ String comparison is case-sensitive "Apple" != "apple"
```

---

## 🔗 Logical Operators

### Truth Tables

**AND (`&&`)** — both must be true

| a | b | a && b |
|---|---|--------|
| true | true | **true** |
| true | false | false |
| false | true | false |
| false | false | false |

**OR (`\|\|`)** — at least one must be true

| a | b | a \|\| b |
|---|---|----------|
| true | true | **true** |
| true | false | **true** |
| false | true | **true** |
| false | false | false |

**NOT (`!`)** — inverts

| a | !a |
|---|----|
| true | false |
| false | true |

---

## 🔢 Bitwise Operators

### Truth Tables (per bit)

**AND (`&`)** | **OR (`\|`)** | **XOR (`^`)**

| a | b | a&b | a\|b | a^b |
|---|---|-----|------|-----|
| 0 | 0 | 0 | 0 | 0 |
| 0 | 1 | 0 | 1 | 1 |
| 1 | 0 | 0 | 1 | 1 |
| 1 | 1 | 1 | 1 | 0 |

### Worked Example: 5 & 3, 5 | 3, 5 ^ 3

```
  5 = 0101
  3 = 0011
  ------
5&3 = 0001 = 1

  5 = 0101
  3 = 0011
  ------
5|3 = 0111 = 7

  5 = 0101
  3 = 0011
  ------
5^3 = 0110 = 6
```

### Shifts

```
5 << 1   0101 → 1010   = 10   (multiply by 2)
5 << 2   0101 → 10100  = 20   (multiply by 4)
5 >> 1   0101 → 0010   = 2    (divide by 2)
5 >> 2   0101 → 0001   = 1    (divide by 4)

x << n  is x * 2^n
x >> n  is x / 2^n  (for positive numbers)
```

### Bit Flags Pattern (from Level 3's `iota`)

```go
const (
    Read   = 1 << iota  // 001 = 1
    Write              // 010 = 2
    Execute            // 100 = 4
)

perms := Read | Write        // combine:  011 = 3
hasRead := perms&Read == Read // check:   001 & 011 = 001 (true)
perms |= Execute              // grant:    011 | 100 = 111
perms &^= Write                // revoke:  111 &^ 010 = 101
```

---

## 📝 Assignment Operators

### Compound Assignment Table

| Operator | Meaning | Example (`x` starts at 10) |
|----------|---------|------------------------------|
| `+=` | `x = x + n` | `x += 5` → 15 |
| `-=` | `x = x - n` | `x -= 5` → 5 |
| `*=` | `x = x * n` | `x *= 2` → 20 |
| `/=` | `x = x / n` | `x /= 2` → 5 |
| `%=` | `x = x % n` | `x %= 3` → 1 |
| `&=` | `x = x & n` | bitwise AND-assign |
| `\|=` | `x = x \| n` | bitwise OR-assign |
| `^=` | `x = x ^ n` | bitwise XOR-assign |
| `<<=` | `x = x << n` | left-shift-assign |
| `>>=` | `x = x >> n` | right-shift-assign |

### Increment / Decrement

```go
x++   // x = x + 1   (STATEMENT, not expression)
x--   // x = x - 1   (STATEMENT, not expression)

y := x++   // ❌ COMPILE ERROR — can't use ++ as a value
x++
y := x     // ✅ RIGHT
```

---

## 📊 Operator Precedence Table

From **highest** to **lowest** — higher rows evaluate first:

| Precedence | Operators |
|------------|-----------|
| 1 (highest) | `()` parentheses |
| 2 | `++` `--` |
| 3 | `*` `/` `%` `&` `<<` `>>` |
| 4 | `+` `-` `\|` `^` |
| 5 | `==` `!=` `<` `>` `<=` `>=` |
| 6 | `&&` |
| 7 (lowest) | `\|\|` |

### Decision Tree

```
Unsure how an expression will evaluate?
├─ Does it mix arithmetic operators?
│  └─ *, /, % run before +, -
├─ Does it mix comparison and logical operators?
│  └─ Comparisons run before && and ||
├─ Does it mix && and ||?
│  └─ && binds tighter than || (evaluates first)
└─ Still unsure or want to be safe?
   └─ Add parentheses — always allowed, never wrong
```

### Worked Examples

```
2 + 3 * 4        → 14   (3*4 first, then +2)
(2 + 3) * 4      → 20   (parens force + first)

5 > 3 && 10 < 20         → true   (both comparisons run first)
true || false && false   → true   (&& binds tighter: false && false = false, then true || false)
(true || false) && false → false  (parens override)
```

---

## ⚡ Short-Circuit Evaluation

```
a && b
   ├─ a is false → b is NEVER evaluated, result is false
   └─ a is true  → b IS evaluated, result is b

a || b
   ├─ a is true  → b is NEVER evaluated, result is true
   └─ a is false → b IS evaluated, result is b
```

### Why It Matters: Safe Guards

```go
// ✅ SAFE — b != 0 is checked first; if false, a/b never runs
if b != 0 && a/b > 10 {
    // ...
}

// ❌ UNSAFE — a/b runs even when b might be 0, order matters!
if a/b > 10 && b != 0 {
    // panics if b == 0
}
```

---

## 🚨 Common Mistakes

### Mistake 1: Forgetting Precedence

```go
// ❌ Expects 20, gets 14
result := 2 + 3 * 4

// ✅ RIGHT
result := (2 + 3) * 4
```

### Mistake 2: Integer Division Loses Precision

```go
// ❌ WRONG — gets 0
percentage := 1 / 3 * 100

// ✅ RIGHT
percentage := int(float64(1) / 3 * 100)
```

### Mistake 3: `=` Instead of `==`

```go
// ❌ WRONG — Go won't even compile this as a condition
if x = 5 { }

// ✅ RIGHT
if x == 5 { }
```

### Mistake 4: Using `++`/`--` As An Expression

```go
// ❌ WRONG — compile error
y := x++

// ✅ RIGHT
x++
y := x
```

### Mistake 5: Unsafe Guard Ordering

```go
// ❌ WRONG — can panic
if a/b > 10 && b != 0 { }

// ✅ RIGHT — short-circuit protects the division
if b != 0 && a/b > 10 { }
```

---

## 📈 Progression Summary

### Understanding Level 4

Level 4 teaches how to act on the values from Level 3:

1. **Arithmetic** — compute with numbers
2. **Comparison** — produce booleans from values
3. **Logical** — combine booleans
4. **Bitwise** — manipulate individual bits
5. **Assignment** — update variables in place
6. **Precedence** — know what runs first
7. **Short-circuit** — write safe, efficient conditions

### Prerequisites for Level 5

Before moving to Level 5 (Conditions), you need:

- ✅ Comfortable with all arithmetic/comparison/logical operators
- ✅ Understand precedence well enough to predict results
- ✅ Understand short-circuit evaluation
- ✅ Can write and evaluate compound expressions

### Ready for Level 5?

Level 5 teaches control flow that *uses* these expressions as conditions:
- `if` / `else if` / `else`
- `switch` statements
- Making decisions in code

---

## ✅ Checklist Before Level 5

- [ ] Can use all arithmetic operators, including modulo
- [ ] Understand integer division truncation
- [ ] Can use all comparison operators
- [ ] Can combine conditions with &&, ||, !
- [ ] Can use bitwise operators for flags
- [ ] Can use compound assignment operators
- [ ] Can predict operator precedence
- [ ] Understand and can prove short-circuit evaluation
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The Categories
Arithmetic, Comparison, Logical, Bitwise, Assignment — five families of operators.

### The Rule
`/` on two ints truncates. Convert to `float64` first if you need decimals.

### The Precedence
`()` → `*/％&<<>>` → `+-|^` → comparisons → `&&` → `||`

### The Safety Net
Short-circuit evaluation lets `&&`/`||` skip unnecessary (or unsafe) work — order your conditions accordingly.

---

## 📚 Next Level

Level 5: Conditions (if/else/switch)
- `if` / `else if` / `else` statements
- `switch` statements
- Making decisions with the expressions you just learned

You've got operators down! Keep going! 🚀
