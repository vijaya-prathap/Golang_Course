# Level 5: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: Conditions
```
Day 1:  if statement basics
Day 2:  if/else and else if chains
Day 3:  if with short init
Day 4:  switch statement basics
Day 5:  switch with short init & conditionless switch
Day 6:  fallthrough
Day 7:  Nested conditions & early returns
```

### Week 2: Practice & Application
```
Day 1:  Exercises 1-3 (if, if/else, else-if chains)
Day 2:  Exercises 4-6 (short init, switch, switch+init)
Day 3:  Exercises 7-8 (conditionless switch, fallthrough)
Day 4:  Exercises 9-10 (early returns, comprehensive practice)
Day 5:  Bonus challenges
Day 6-7: Practice & review
```

---

## 🎯 if vs switch At A Glance

| Situation | Prefer |
|-----------|--------|
| One or two conditions | `if` / `if-else` |
| Many conditions on the same value | `switch value { case ... }` |
| Many conditions, no single value (ranges) | `switch { case cond1: ... }` |
| Need a scoped helper variable | short-init form of either |

---

## 🌳 if / else if / else Decision Flow

```
if condition1 {
    // runs if condition1 is true
} else if condition2 {
    // runs if condition1 is false AND condition2 is true
} else {
    // runs if NEITHER condition1 NOR condition2 is true
}

Evaluation order: top to bottom.
The FIRST true condition's block runs; the rest are skipped entirely.
```

### Ordering Rule For Overlapping Ranges

```
Thresholds must be checked from MOST specific (highest) to LEAST specific (lowest):

score := 95

✅ RIGHT                          ❌ WRONG
if score >= 90 { "A" }            if score >= 70 { "C" }   ← matches first, wrong!
else if score >= 70 { "C" }       else if score >= 90 { "A" }
```

---

## 🔑 The Short-Init Pattern

```go
if <init-statement>; <condition> {
    // init-statement's variables are visible here
} else {
    // ...and here
}
// but NOT here - scope ends with the if/else
```

### Two Common Uses

**1. Error checking**
```go
if value, err := strconv.Atoi(s); err != nil {
    // handle error
} else {
    // use value
}
```

**2. The "comma-ok" idiom (map/type lookups)**
```go
if value, ok := myMap[key]; ok {
    // key was present, value is valid
} else {
    // key was absent
}
```

---

## 🔀 switch Anatomy

```go
switch <optional-init>; <optional-value> {
case val1, val2:      // comma-separated = OR
    // ...
case val3:
    // ...
default:
    // ...
}
```

### Four switch Shapes

```
1. switch value { case a: ... }              — basic
2. switch init; value { case a: ... }        — with short init
3. switch { case cond1: ... }                — conditionless (like if/else if)
4. switch init; { case cond1: ... }          — conditionless + short init
```

### Key Behavior: No Automatic Fallthrough

```
Go switch:  each case implicitly breaks after its body runs.
C/Java switch: falls through to the next case unless you break.

Go is the SAFER default - most bugs in C-style switches come from
forgetting a break. Go requires you to opt IN with `fallthrough`.
```

---

## ➡️ fallthrough Flow

```go
switch n {
case 1:
    A()
    fallthrough   // unconditionally runs case 2's body next
case 2:
    B()
    fallthrough
case 3:
    C()
case 4:
    D()
}
```

```
n == 1  →  runs A(), B(), C()          (falls through twice)
n == 2  →  runs B(), C()               (falls through once)
n == 3  →  runs C()                    (no fallthrough)
n == 4  →  runs D()                    (never reached from above)
```

⚠️ fallthrough ignores the next case's own condition — it always runs, even if that case wouldn't have matched on its own.

---

## 🪆 Flattening Nested Conditions

### Before: Nested (grows sideways with every new rule)

```go
func process(u user) string {
    if u.active {
        if u.verified {
            if u.hasPermission {
                return "OK"
            }
            return "No permission"
        }
        return "Not verified"
    }
    return "Inactive"
}
```

### After: Early Returns (grows downward, each rule is independent)

```go
func process(u user) string {
    if !u.active {
        return "Inactive"
    }
    if !u.verified {
        return "Not verified"
    }
    if !u.hasPermission {
        return "No permission"
    }
    return "OK"
}
```

### Decision Tree: When To Flatten

```
Are you nesting if inside if inside if?
├─ Each branch returns/exits immediately?
│  └─ Flatten with early returns (guard clauses)
├─ Branches need to do more work after the nested check?
│  └─ Consider extracting the nested logic into its own function
└─ Only 2 levels and both branches are short?
   └─ Nesting is fine - don't over-engineer
```

---

## 🚨 Common Mistakes

### Mistake 1: Wrong Threshold Order

```go
// ❌ WRONG
if score >= 70 { grade = "C" } else if score >= 90 { grade = "A" }

// ✅ RIGHT
if score >= 90 { grade = "A" } else if score >= 70 { grade = "C" }
```

### Mistake 2: Expecting C-Style Fallthrough

```go
// ❌ WRONG ASSUMPTION - only "two" prints, "three" does NOT
switch 2 {
case 2:
    fmt.Println("two")
case 3:
    fmt.Println("three")
}

// ✅ Use fallthrough explicitly if that's what you want
```

### Mistake 3: Brace On The Wrong Line

```go
// ❌ WRONG - compile error
if x > 5
{
}

// ✅ RIGHT
if x > 5 {
}
```

### Mistake 4: Using Short-Init Variable Outside Its Scope

```go
if value, ok := lookup(); ok {
    fmt.Println(value)
}
fmt.Println(value)  // ❌ ERROR - out of scope
```

### Mistake 5: Deep Nesting Instead Of Early Returns

```go
// ❌ Gets harder to follow with every new rule
if a {
    if b {
        if c {
            doWork()
        }
    }
}

// ✅ Flatten
if !a { return }
if !b { return }
if !c { return }
doWork()
```

---

## 📈 Progression Summary

### Understanding Level 5

Level 5 teaches how to branch based on the expressions from Level 4:

1. **if / else** — binary decisions
2. **else if chains** — multiple ordered conditions
3. **short init** — scoped helper variables
4. **switch** — clean multi-way branching on one value
5. **conditionless switch** — clean multi-way branching on multiple conditions
6. **fallthrough** — explicit case chaining (rare)
7. **early returns** — flattening nested logic

### Prerequisites for Level 6

Before moving to Level 6 (Loops), you need:

- ✅ Comfortable writing if/else and else-if chains in the right order
- ✅ Comfortable with the short-init pattern
- ✅ Can choose between if and switch appropriately
- ✅ Understand why Go's switch doesn't fall through by default
- ✅ Can flatten nested conditions with early returns

### Ready for Level 6?

Level 6 teaches how to repeat code using conditions as loop guards:
- `for` loops (Go's only loop keyword)
- while-style and infinite loops
- `break` and `continue`

---

## ✅ Checklist Before Level 6

- [ ] Can write if / else / else-if chains in correct order
- [ ] Can use the short-init form for if and switch
- [ ] Can write a basic switch with multiple cases
- [ ] Can write a conditionless switch
- [ ] Understand fallthrough and know it's rarely needed
- [ ] Can flatten nested conditions using early returns
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The Rule
The FIRST true condition wins in both if/else-if chains and conditionless switches — order matters.

### The Default
Go's switch does NOT fall through. Each case implicitly breaks.

### The Scope
Short-init variables (`if x := f(); ...`) live only inside that if/else or switch.

### The Style
Prefer flat, early-return code over deeply nested conditions.

---

## 📚 Next Level

Level 6: Loops
- for loops (Go's only loop keyword)
- while-style and infinite loops
- break and continue
- Looping over collections

You've got decision-making down! Keep going! 🚀
