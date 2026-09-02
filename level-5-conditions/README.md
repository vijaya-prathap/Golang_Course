# Level 5: Conditions (if/else/switch) - Complete Guide

## Introduction

Welcome to Level 5! You've mastered operators and expressions (Level 4). Now it's time to learn **conditions** - how to make your programs choose what to do based on the values and expressions you can now build.

A **condition** is a boolean expression Go evaluates to decide which code runs. `if`, `else`, and `switch` are Go's tools for branching based on conditions.

For example:
- `if age >= 18 { ... }` runs the block only when the condition is true
- `switch day { case "Mon": ... }` picks one branch out of many

---

## Table of Contents

1. [The if Statement](#the-if-statement)
2. [if/else](#ifelse)
3. [if/else if/else Chains](#ifelse-ifelse-chains)
4. [The if Statement With a Short Init](#the-if-statement-with-a-short-init)
5. [The switch Statement](#the-switch-statement)
6. [switch Without a Condition](#switch-without-a-condition)
7. [Fallthrough](#fallthrough)
8. [Nested Conditions](#nested-conditions)
9. [Best Practices](#best-practices)
10. [Common Mistakes](#common-mistakes)

---

## The if Statement

An `if` statement runs a block only when its condition is `true`.

### Basic Syntax

```go
if condition {
    // runs only if condition is true
}
```

### Example

```go
age := 20

if age >= 18 {
    fmt.Println("You are an adult")
}
```

### No Parentheses Required

Unlike C, Java, or JavaScript, Go does NOT require parentheses around the condition — and the opening brace `{` MUST be on the same line as `if`:

```go
// ✅ Go style
if age >= 18 {
    fmt.Println("Adult")
}

// ❌ Compile error - brace on its own line
if age >= 18
{
    fmt.Println("Adult")
}
```

---

## if/else

`else` runs when the `if` condition is false.

```go
age := 15

if age >= 18 {
    fmt.Println("You are an adult")
} else {
    fmt.Println("You are a minor")
}
```

Exactly one of the two blocks always runs.

---

## if/else if/else Chains

Chain multiple conditions with `else if`. Go checks each condition top-to-bottom and runs the **first** one that's true.

```go
score := 82

if score >= 90 {
    fmt.Println("Grade: A")
} else if score >= 80 {
    fmt.Println("Grade: B")
} else if score >= 70 {
    fmt.Println("Grade: C")
} else {
    fmt.Println("Grade: F")
}
```

Output: `Grade: B` (the first true condition wins; later ones are never checked)

### Order Matters

```go
// ❌ WRONG - score >= 70 is always true when score >= 90 is also true,
// but the FIRST matching branch runs, so order still needs to go
// from most specific/highest to least specific/lowest for range checks
score := 95

if score >= 70 {
    fmt.Println("Grade: C") // Wrong! prints C, never reaches the A/B checks
} else if score >= 90 {
    fmt.Println("Grade: A")
}

// ✅ RIGHT - check the highest threshold first
if score >= 90 {
    fmt.Println("Grade: A")
} else if score >= 70 {
    fmt.Println("Grade: C")
}
```

---

## The if Statement With a Short Init

`if` can run a short statement before the condition, scoped only to the `if`/`else` blocks. This is one of Go's most distinctive features.

```go
if x := compute(); x > 10 {
    fmt.Println("x is large:", x)
} else {
    fmt.Println("x is small:", x)
}
// x does NOT exist here - it's scoped to the if/else
```

### Common Pattern: Error Checking

```go
if value, err := strconv.Atoi("42"); err != nil {
    fmt.Println("Conversion failed:", err)
} else {
    fmt.Println("Converted value:", value)
}
```

### Common Pattern: Map Lookups

```go
scores := map[string]int{"Alice": 95, "Bob": 87}

if score, ok := scores["Alice"]; ok {
    fmt.Println("Alice's score:", score)
} else {
    fmt.Println("Alice not found")
}
```

---

## The switch Statement

`switch` compares one value against multiple cases. Unlike C/Java, Go's `switch` does **not** fall through by default — each case automatically `break`s after its block runs.

### Basic Syntax

```go
day := "Tuesday"

switch day {
case "Monday":
    fmt.Println("Start of the week")
case "Tuesday", "Wednesday", "Thursday":
    fmt.Println("Midweek")
case "Friday":
    fmt.Println("Almost the weekend")
default:
    fmt.Println("Weekend")
}
```

Notice a single `case` can list multiple comma-separated values.

### switch With a Short Init

Just like `if`, `switch` supports a short statement before the value being switched on:

```go
switch hour := time.Now().Hour(); {
case hour < 12:
    fmt.Println("Good morning")
case hour < 18:
    fmt.Println("Good afternoon")
default:
    fmt.Println("Good evening")
}
```

---

## switch Without a Condition

A `switch` with no value after the keyword acts like a cleaner `if/else if` chain — each `case` is a boolean expression, and the first true one runs.

```go
score := 82

switch {
case score >= 90:
    fmt.Println("Grade: A")
case score >= 80:
    fmt.Println("Grade: B")
case score >= 70:
    fmt.Println("Grade: C")
default:
    fmt.Println("Grade: F")
}
```

This is often more readable than a long `if/else if` chain.

---

## Fallthrough

Go's `switch` does not fall through by default. Use the `fallthrough` keyword to explicitly run the next case too.

```go
switch n := 2; n {
case 1:
    fmt.Println("one")
    fallthrough
case 2:
    fmt.Println("two")
    fallthrough
case 3:
    fmt.Println("three")
case 4:
    fmt.Println("four")
}
```

Output:
```
two
three
```

`fallthrough` unconditionally runs the next case's body, even if that case's own condition wouldn't have matched. Use it sparingly — it's rarely needed in idiomatic Go.

---

## Nested Conditions

`if` and `switch` can be nested inside each other, but deep nesting hurts readability.

```go
func classify(age int, hasLicense bool) string {
    if age >= 18 {
        if hasLicense {
            return "Adult driver"
        }
        return "Adult, no license"
    }
    return "Minor"
}
```

### Prefer Early Returns Over Deep Nesting

```go
// ❌ Deeply nested
func process(user User) string {
    if user.Active {
        if user.Verified {
            if user.HasPermission {
                return "OK"
            } else {
                return "No permission"
            }
        } else {
            return "Not verified"
        }
    } else {
        return "Inactive"
    }
}

// ✅ Early returns - flatter and easier to follow
func process(user User) string {
    if !user.Active {
        return "Inactive"
    }
    if !user.Verified {
        return "Not verified"
    }
    if !user.HasPermission {
        return "No permission"
    }
    return "OK"
}
```

---

## Best Practices

### 1. Prefer switch Over Long if/else if Chains

```go
// ✅ Good - clear and scannable
switch {
case score >= 90:
    grade = "A"
case score >= 80:
    grade = "B"
default:
    grade = "F"
}

// ❌ Harder to scan
if score >= 90 {
    grade = "A"
} else if score >= 80 {
    grade = "B"
} else {
    grade = "F"
}
```

### 2. Use the Short Init Form to Limit Scope

```go
// ✅ Good - err only exists where it's needed
if err := doSomething(); err != nil {
    return err
}

// ❌ Less good - err leaks into surrounding scope
err := doSomething()
if err != nil {
    return err
}
```

### 3. Use Early Returns to Avoid Nesting

Return, `continue`, or `break` as soon as you know the answer, instead of nesting the "happy path" deeper and deeper.

### 4. Order if/else if From Most to Least Specific

When ranges overlap (like grade thresholds), check the highest/most specific condition first.

### 5. Avoid fallthrough Unless Truly Needed

Most switch logic reads more clearly without it. Reach for `fallthrough` only when multiple cases genuinely share trailing behavior.

---

## Common Mistakes

### Mistake 1: Wrong Order in if/else if Chains

```go
// ❌ WRONG - score=95 incorrectly prints "C"
if score >= 70 {
    fmt.Println("C")
} else if score >= 90 {
    fmt.Println("A")
}

// ✅ RIGHT - highest threshold first
if score >= 90 {
    fmt.Println("A")
} else if score >= 70 {
    fmt.Println("C")
}
```

### Mistake 2: Expecting switch to Fall Through

```go
// ❌ WRONG ASSUMPTION - this does NOT print "two" then "three"
switch 2 {
case 1:
    fmt.Println("one")
case 2:
    fmt.Println("two") // only this prints
case 3:
    fmt.Println("three")
}

// ✅ Use fallthrough explicitly if you want that behavior
switch 2 {
case 1:
    fmt.Println("one")
case 2:
    fmt.Println("two")
    fallthrough
case 3:
    fmt.Println("three") // now this also prints
}
```

### Mistake 3: Brace on the Wrong Line

```go
// ❌ WRONG - compile error
if x > 5
{
    fmt.Println("big")
}

// ✅ RIGHT - opening brace on the same line
if x > 5 {
    fmt.Println("big")
}
```

### Mistake 4: Using = Instead of ==

```go
// ❌ WRONG - assignment isn't allowed as a condition, won't compile
if x = 5 {
}

// ✅ RIGHT
if x == 5 {
}
```

### Mistake 5: Forgetting Short-Init Variables Are Scoped

```go
if value, ok := lookup(); ok {
    fmt.Println(value)
}
fmt.Println(value) // ❌ ERROR - value doesn't exist here
```

### Mistake 6: Deep Nesting Instead of Early Returns

```go
// ❌ Hard to follow after a few more conditions
if a {
    if b {
        if c {
            doWork()
        }
    }
}

// ✅ Flatten with early returns
if !a {
    return
}
if !b {
    return
}
if !c {
    return
}
doWork()
```

---

## Summary

**if Statement:**
- Runs a block only when its condition is true
- Opening brace must be on the same line as `if`
- No parentheses required around the condition

**if/else and else if:**
- `else` runs when `if` is false
- `else if` chains check conditions top-to-bottom; first true one wins
- Order matters for overlapping ranges - check most specific first

**Short Init:**
- `if x := f(); x > 0 { ... }` scopes `x` to the if/else blocks only
- Common for error checking and map lookups (`value, ok := m[key]`)

**switch Statement:**
- Compares one value against multiple cases
- Automatically breaks after each case (no fallthrough by default)
- A single case can list multiple comma-separated values
- `switch {}` with no value acts like a clean if/else if chain

**fallthrough:**
- Explicitly runs the next case's body
- Rarely needed in idiomatic Go

**Best Practices:**
- Prefer `switch` for many branches, `if/else` for one or two
- Use early returns to avoid deep nesting
- Use short-init to keep variables scoped tightly

---

## Next Steps

You now understand:
- ✅ if, else, and else if
- ✅ The short-init form
- ✅ switch, including the conditionless form
- ✅ fallthrough
- ✅ How to avoid deep nesting

**Next level:** Level 6 - Loops
- for loops (Go's only loop keyword)
- while-style and infinite loops
- break and continue
- Looping over collections

You're building real decision-making logic! Keep going! 🚀
