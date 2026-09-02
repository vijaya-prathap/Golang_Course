# Level 4: Operators & Expressions - Complete Guide

## Introduction

Welcome to Level 4! You've mastered variables and data types (Level 3). Now it's time to learn about **operators and expressions** - how to work WITH those values.

An **operator** performs an action on one or more values. An **expression** is code that produces a value.

For example:
- `5 + 3` is an expression using the `+` operator
- `x > 10` is an expression using the `>` operator
- `a && b` is an expression using the `&&` operator

---

## Table of Contents

1. [Arithmetic Operators](#arithmetic-operators)
2. [Comparison Operators](#comparison-operators)
3. [Logical Operators](#logical-operators)
4. [Bitwise Operators](#bitwise-operators)
5. [Assignment Operators](#assignment-operators)
6. [Operator Precedence](#operator-precedence)
7. [Expressions](#expressions)
8. [Short-Circuit Evaluation](#short-circuit-evaluation)
9. [Best Practices](#best-practices)
10. [Common Mistakes](#common-mistakes)

---

## Arithmetic Operators

Arithmetic operators work with numbers (int, float).

### Basic Arithmetic

```go
a := 10
b := 3

a + b  // 13 (addition)
a - b  // 7  (subtraction)
a * b  // 30 (multiplication)
a / b  // 3  (division - integer division!)
a % b  // 1  (modulo - remainder)
```

### Important: Integer Division

When both operands are integers, division truncates (loses decimal):

```go
10 / 3   // 3 (not 3.333)
10 / 4   // 2 (not 2.5)

// To get decimal result, convert to float
float64(10) / 3  // 3.333...
```

### Modulo (Remainder)

```go
10 % 3   // 1 (10 = 3*3 + 1)
15 % 4   // 3 (15 = 4*3 + 3)
20 % 5   // 0 (20 = 5*4 + 0)

// Common use: check if number is even or odd
if x % 2 == 0 {
    fmt.Println("Even")
} else {
    fmt.Println("Odd")
}
```

### Unary Operators

```go
x := 5
-x   // -5 (negation)
+x   // 5  (affirmation - rarely used)
```

---

## Comparison Operators

Comparison operators compare two values and return a boolean (true/false).

### All Comparison Operators

```go
a := 10
b := 3

a == b  // false (equal)
a != b  // true  (not equal)
a > b   // true  (greater than)
a >= b  // true  (greater than or equal)
a < b   // false (less than)
a <= b  // false (less than or equal)
```

### With Different Types

```go
5 > 3         // true
5.5 > 3.2     // true
"apple" == "apple"  // true
"abc" < "xyz" // true (lexicographic comparison)

// Can't compare different types!
// 5 == "5"   // ERROR!
```

### String Comparison

```go
"abc" == "abc"   // true
"abc" != "xyz"   // true
"abc" < "xyz"    // true (alphabetical)
"apple" > "apricot"  // false
```

---

## Logical Operators

Logical operators work with boolean values.

### AND Operator (&&)

Both conditions must be true:

```go
a := true
b := false

a && b   // false (both must be true)
a && a   // true
b && b   // false
```

### OR Operator (||)

At least one condition must be true:

```go
a := true
b := false

a || b   // true (at least one is true)
a || a   // true
b || b   // false
```

### NOT Operator (!)

Inverts the boolean:

```go
a := true
b := false

!a   // false
!b   // true
!!a  // true (not not true = true)
```

### Combining Logical Operators

```go
x := 5
y := 10

(x > 0) && (y > 0)           // true (both positive)
(x < 0) || (y < 0)           // false (neither negative)
!(x == y)                    // true (not equal)

// Common pattern: validation
if age >= 18 && hasLicense {
    fmt.Println("Can drive")
}

if score >= 90 || hasExtraCredit {
    fmt.Println("Passed")
}
```

---

## Bitwise Operators

Bitwise operators work on individual bits of integers.

### Bitwise AND (&)

Both bits must be 1:

```go
5 & 3  // 1
// 5 = 0101
// 3 = 0011
//     0001 = 1
```

### Bitwise OR (|)

At least one bit must be 1:

```go
5 | 3  // 7
// 5 = 0101
// 3 = 0011
//     0111 = 7
```

### Bitwise XOR (^)

Bits must be different:

```go
5 ^ 3  // 6
// 5 = 0101
// 3 = 0011
//     0110 = 6
```

### Bitwise NOT (^) - Unary

Inverts all bits:

```go
^5  // -6 (in two's complement)
```

### Left Shift (<<)

Shifts bits left, fills with 0s:

```go
5 << 1   // 10 (101 becomes 1010)
5 << 2   // 20 (101 becomes 10100)

// Multiply by 2^n
x << 1   // Same as x * 2
x << 3   // Same as x * 8
```

### Right Shift (>>)

Shifts bits right, fills with sign bit:

```go
5 >> 1   // 2 (101 becomes 10)
5 >> 2   // 1 (101 becomes 1)

// Divide by 2^n (for positive numbers)
x >> 1   // Same as x / 2
x >> 3   // Same as x / 8
```

### Bitwise Operators in Practice

Flag checking (from Level 3):

```go
const (
    Read   = 1 << iota  // 1
    Write              // 2
    Execute            // 4
)

perms := Read | Write  // 3 (binary: 011)

// Check if Read permission set
if (perms & Read) == Read {
    fmt.Println("Can read")
}

// Check if Execute permission set
if (perms & Execute) == Execute {
    fmt.Println("Can execute")
} else {
    fmt.Println("Cannot execute")
}
```

---

## Assignment Operators

Assignment operators assign values to variables.

### Basic Assignment

```go
x := 5      // Assign 5 to x
x = 10      // Reassign 10 to x
```

### Compound Assignment

```go
x := 5

x += 3   // x = x + 3  (x is now 8)
x -= 2   // x = x - 2  (x is now 6)
x *= 2   // x = x * 2  (x is now 12)
x /= 3   // x = x / 3  (x is now 4)
x %= 3   // x = x % 3  (x is now 1)

x &= 7   // x = x & 7  (bitwise AND)
x |= 2   // x = x | 2  (bitwise OR)
x ^= 1   // x = x ^ 1  (bitwise XOR)

x <<= 1  // x = x << 1 (left shift)
x >>= 2  // x = x >> 2 (right shift)
```

### Increment and Decrement

```go
x := 5

x++   // x = x + 1 (x is now 6)
x--   // x = x - 1 (x is now 5)

// Note: In Go, ++ and -- are statements, not expressions!
// This is WRONG:
// y := x++     // Syntax error!

// This is RIGHT:
x++
y := x
```

---

## Operator Precedence

Operators have a priority order. Higher precedence operators execute first.

### Precedence Order (Highest to Lowest)

```
1. ()           Parentheses
2. ++, --       Increment/decrement
3. *, /, %, &, << >>  Multiplication, division, bitwise
4. +, -, |, ^   Addition, subtraction, bitwise
5. ==, !=, <, >, <=, >=  Comparison
6. &&           Logical AND
7. ||           Logical OR
```

### Examples

```go
2 + 3 * 4      // 14 (multiply first, then add)
(2 + 3) * 4    // 20 (parentheses force addition first)

5 > 3 && 10 < 20     // true (comparisons first, then AND)
5 > 3 || 10 > 20     // true (comparisons first, then OR)

true || false && false   // true (AND before OR)
(true || false) && false // false (parentheses override)

2 + 3 * 4 - 5      // 9 (mult: 12, then 2+12-5)
```

### Use Parentheses for Clarity

Even when not required, parentheses make intent clear:

```go
// Clear intent
if (age >= 18) && (hasLicense) {
    fmt.Println("Can drive")
}

// Less clear
if age >= 18 && hasLicense {
    fmt.Println("Can drive")
}
```

---

## Expressions

An **expression** is a piece of code that produces a value.

### Simple Expressions

```go
5 + 3           // Expression: evaluates to 8
x > 10          // Expression: evaluates to true or false
true && false   // Expression: evaluates to false
```

### Expressions in Assignments

```go
result := 5 + 3 * 2   // result gets 11
isAdult := age >= 18
canDrive := isAdult && hasLicense
```

### Expressions in Conditions

```go
if 5 > 3 {
    // This block runs
}

if x >= 10 && x <= 20 {
    // x is between 10 and 20
}
```

### Type of Expression

Expressions have types too:

```go
5 + 3           // int
5.0 + 3.0       // float64
5 > 3           // bool
"hello" + " world"  // string
```

---

## Short-Circuit Evaluation

Logical operators don't always evaluate all parts.

### AND Short-Circuit

If first operand is false, second isn't evaluated:

```go
false && (x / 0)   // No error! (x/0) isn't evaluated
true && (x / 0)    // ERROR! (x/0) is evaluated
```

### OR Short-Circuit

If first operand is true, second isn't evaluated:

```go
true || (x / 0)    // No error! (x/0) isn't evaluated
false || (x / 0)   // ERROR! (x/0) is evaluated
```

### Practical Use

```go
// Safe: doesn't divide by zero
if b != 0 && a/b > 10 {
    // ...
}

// Not safe: might divide by zero
if a/b > 10 && b != 0 {  // Check order matters!
    // ...
}
```

---

## Best Practices

### 1. Use Parentheses for Clarity

```go
// ✅ Good
if (age >= 18) && (hasLicense) {
    canDrive = true
}

// ❌ Unclear
if age >= 18 && hasLicense {
    canDrive = true
}
```

### 2. Order Conditions Efficiently

Put cheaper checks first:

```go
// ✅ Good (variable check is cheaper than function call)
if found && isValid(data) {
    // ...
}

// ❌ Less efficient
if isValid(data) && found {
    // ...
}
```

### 3. Avoid Complex Expressions

```go
// ✅ Good - readable
isEligible := age >= 18
hasRequirements := experience > 5 && score > 80
if isEligible && hasRequirements {
    hire()
}

// ❌ Hard to read
if age >= 18 && experience > 5 && score > 80 {
    hire()
}
```

### 4. Use Meaningful Variable Names

```go
// ✅ Good
isValid := x > 0 && x < 100
hasPermission := user.Role == "admin"
canProceed := isValid && hasPermission

// ❌ Unclear
ok := x > 0 && x < 100 && u == "a"
```

### 5. Handle Integer Division Carefully

```go
// ✅ Good - awareness of integer division
percentage := int(float64(correct) / float64(total) * 100)

// ❌ Wrong - loses decimal places
percentage := correct / total * 100
```

---

## Common Mistakes

### Mistake 1: Forgetting Operator Precedence

```go
// ❌ WRONG - expects 15, gets 14
result := 2 + 3 * 4   // 14, not 20

// ✅ RIGHT
result := (2 + 3) * 4  // 20
```

### Mistake 2: Integer Division Loss

```go
// ❌ WRONG - gets 0
percentage := 1 / 3 * 100   // 0

// ✅ RIGHT
percentage := int(float64(1) / 3 * 100)  // 33
```

### Mistake 3: Using = Instead of ==

```go
// ❌ WRONG - assignment, not comparison
if x = 5 {
    // Error!
}

// ✅ RIGHT
if x == 5 {
    // Code
}
```

### Mistake 4: Comparing Different Types

```go
// ❌ WRONG
if x == "5" {  // x is int, "5" is string
    // Compiler error!
}

// ✅ RIGHT
if x == 5 {
    // OK
}
```

### Mistake 5: Complex Conditions Without Parentheses

```go
// ❌ WRONG - confusing precedence
if x > 5 && y > 10 || z < 3 {
    // Hard to read!
}

// ✅ RIGHT
if (x > 5 && y > 10) || (z < 3) {
    // Clear intent!
}
```

### Mistake 6: Logical Operator Confusion

```go
// ❌ WRONG - will always be true if either value is positive
if x > 0 || x < 0 {
    // This doesn't check for both positive AND negative
}

// ✅ RIGHT - for AND, both must be true
if x > 0 && y > 0 {
    // Both are positive
}
```

---

## Summary

**Arithmetic Operators:**
- `+`, `-`, `*`, `/`, `%` for numbers
- `/` with integers truncates
- `%` gives remainder

**Comparison Operators:**
- `==`, `!=`, `<`, `>`, `<=`, `>=`
- Return boolean (true/false)

**Logical Operators:**
- `&&` (AND), `||` (OR), `!` (NOT)
- Work with booleans
- Short-circuit evaluation

**Bitwise Operators:**
- `&`, `|`, `^`, `~` for bit manipulation
- `<<`, `>>` for shifting
- Used for flags and low-level work

**Operator Precedence:**
- Use parentheses for clarity
- Precedence: (), ++/--, *, /, +, -, ==, <, &&, ||

**Expressions:**
- Code that produces a value
- Can be assigned to variables
- Used in conditions

---

## Next Steps

You now understand:
- ✅ All basic operators
- ✅ Operator precedence
- ✅ Expression evaluation
- ✅ Short-circuit evaluation
- ✅ Best practices

**Next level:** Level 5 - Conditions (if/else/switch)
- Control flow with if statements
- Switch statements
- Making decisions in code

You're building control over Go programs! Keep going! 🚀
