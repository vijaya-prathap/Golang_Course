# Level 0: Study Guide & Quick Reference

## 🎓 Learning Path

```
Week 1: Fundamentals
├─ What is Programming?
├─ Languages & Compilation
├─ Algorithm Basics
└─ Problem-Solving Methodology

Week 2: Core Concepts
├─ Variables & Data Types
├─ Control Flow (If/Else)
├─ Loops (For/While)
└─ Functions

Week 3: Practice & Application
├─ Solve Exercise 1-4
├─ Solve Exercise 5-8
├─ Challenge Exercises
└─ Pseudocode Mastery
```

---

## 📊 Control Flow Diagrams

### If-Else Statement
```
        ┌─────────────────┐
        │  START          │
        └────────┬────────┘
                 │
        ┌────────▼────────┐
        │ Condition True? │
        └────┬────────┬───┘
             │        │
            Yes       No
             │        │
      ┌──────▼──┐  ┌──▼──────┐
      │ Action A│  │ Action B │
      └──────┬──┘  └──┬──────┘
             │        │
        ┌────▼────────▼──┐
        │      END       │
        └────────────────┘
```

### While Loop
```
        ┌─────────────────┐
        │  START          │
        └────────┬────────┘
                 │
        ┌────────▼────────────┐
        │ Condition True?     │
        └────┬────────────┬───┘
            Yes           No
             │            │
      ┌──────▼──────┐  ┌──▼──────┐
      │ Do Action   │  │    END   │
      │ Update Loop │  └──────────┘
      │ Variable    │
      └──────┬──────┘
             │
        ┌────▼────────────────┐
        │ Go back to condition│
        └────────────────────┘
```

### For Loop
```
        ┌─────────────────┐
        │  START          │
        └────────┬────────┘
                 │
        ┌────────▼────────┐
        │ Initialize i=1  │
        └────────┬────────┘
                 │
        ┌────────▼────────────┐
        │ i <= n?             │
        └────┬────────────┬───┘
            Yes           No
             │            │
      ┌──────▼──────┐  ┌──▼──────┐
      │ Do Action   │  │    END   │
      │ i = i + 1   │  └──────────┘
      └──────┬──────┘
             │
        ┌────▼────────────────┐
        │ Go back to condition│
        └────────────────────┘
```

### Switch Statement
```
        ┌──────────────┐
        │   START      │
        └────────┬─────┘
                 │
        ┌────────▼──────────┐
        │   Switch(value)   │
        └────┬──┬──┬──┬─────┘
             │  │  │  │
        ┌────┴──┴──┴──┴────┐
        │                  │
      ┌─▼─┐ ┌─▼─┐ ┌─▼──┐ ┌─▼──┐
      │A  │ │B  │ │C   │ │DEF │
      │   │ │   │ │    │ │    │
      └─┬─┘ └─┬─┘ └─┬──┘ └─┬──┘
        │     │     │      │
        └─────┴─────┴──────┴───┐
                              │
                        ┌─────▼──┐
                        │   END   │
                        └─────────┘
```

---

## 📝 Data Type Quick Reference

```
┌─────────────────────────────────────────┐
│          PRIMITIVE DATA TYPES            │
├─────────────────────────────────────────┤
│ Integer (Whole numbers)                 │
│  Example: -5, 0, 42, 1000              │
│  Operations: +, -, ×, ÷, MOD (%)       │
│                                         │
│ Float (Decimal numbers)                 │
│  Example: 3.14, -2.5, 0.001            │
│  Operations: +, -, ×, ÷                │
│                                         │
│ String (Text)                           │
│  Example: "Hello", "Go", "123"         │
│  Operations: Concatenation, length     │
│                                         │
│ Boolean (True/False)                    │
│  Example: true, false                  │
│  Operations: AND, OR, NOT              │
└─────────────────────────────────────────┘
```

---

## 🔄 Variable Lifecycle

```
Declaration  → Initialization  → Usage  → Modification  → Return/End
   var x         x = 5         print(x)    x = x + 1      cleanup
```

**Example:**
```
// Declaration: Create variable
age

// Initialization: Give it a value
age = 25

// Usage: Use the value
print(age)

// Modification: Change the value
age = age + 1

// Now age = 26
```

---

## 🧮 Operator Reference

### Arithmetic Operators
```
+  Addition       → 5 + 3 = 8
-  Subtraction    → 5 - 3 = 2
×  Multiplication → 5 × 3 = 15
÷  Division       → 6 ÷ 3 = 2
%  Modulo (Remainder) → 7 % 3 = 1
```

### Comparison Operators
```
=   Equal to                → 5 = 5 (true)
≠   Not equal to            → 5 ≠ 3 (true)
>   Greater than            → 5 > 3 (true)
<   Less than               → 5 < 3 (false)
≥   Greater than or equal   → 5 ≥ 5 (true)
≤   Less than or equal      → 5 ≤ 3 (false)
```

### Logical Operators
```
AND  Both conditions true   → (5>3) AND (3>1) = true
OR   At least one true      → (5>3) OR (1>3) = true
NOT  Reverse the value      → NOT (5>3) = false
```

---

## 🎯 Problem-Solving Flowchart

```
         START
           │
           ▼
    ┌──────────────────┐
    │ Understand       │
    │ Problem?         │
    └────┬─────┬──────┘
        Yes    No
         │      └──→ Re-read problem
         │           │
         ▼           │
    ┌──────────────────┐
    │ Plan Solution    │◄──┘
    │ (Algorithm)      │
    └────┬─────┬──────┘
        Done   Not Clear
         │      └──→ Break into smaller parts
         │           │
         ▼           │
    ┌──────────────────┐
    │ Implement Code   │◄──┘
    └────┬──────┬─────┘
       Done    Issues
        │       └──→ Debug
        │           │
        ▼           │
    ┌──────────────────┐
    │ Test Code        │
    │ (Edge Cases)     │◄──┘
    └────┬──────┬─────┘
      Pass   Fail
       │      └──→ Fix bugs
       │          │
       ▼          │
    ┌──────────────────┐
    │ Optimize?       │◄──┘
    │ (if needed)      │
    └────┬──────┬─────┘
      Done   Yes
       │      └──→ Make faster/cleaner
       │          │
       ▼          │
    COMPLETE  ◄──┘
```

---

## 💡 Common Patterns

### Pattern 1: Finding Maximum
```
max = array[0]
FOR each element in array:
    IF element > max:
        max = element
RETURN max
```

### Pattern 2: Counting Occurrences
```
count = 0
FOR each element in array:
    IF element == target:
        count = count + 1
RETURN count
```

### Pattern 3: Summing Values
```
sum = 0
FOR each element in array:
    sum = sum + element
RETURN sum
```

### Pattern 4: Filtering
```
result = []
FOR each element in array:
    IF element meets criteria:
        ADD element to result
RETURN result
```

### Pattern 5: Validation
```
FOR each character in input:
    IF character is invalid:
        RETURN false
RETURN true
```

---

## ⚡ Quick Tips

1. **Always Test Edge Cases:**
   - Empty input
   - Single element
   - Boundary values
   - Negative numbers

2. **Use Meaningful Names:**
   - `userAge` not `ua`
   - `isValid` not `v`
   - `calculateTotal` not `ct`

3. **Break Down Complex Problems:**
   - Don't try to solve everything at once
   - Solve smaller sub-problems first
   - Combine solutions

4. **Document Your Logic:**
   - Add comments explaining WHY, not WHAT
   - Code shows WHAT, comments explain WHY

5. **Test Before Moving On:**
   - Always verify your algorithm works
   - Test with multiple inputs
   - Think about what could break it

---

## 🎮 Practice Approach

### Day 1-2: Learn Concepts
- Read through README.md
- Understand control flow diagrams
- Study examples

### Day 3-4: Pseudocode
- Write pseudocode for exercises 1-4
- Review and refine
- Compare with solutions

### Day 5-6: More Exercises
- Write pseudocode for exercises 5-8
- Trace through with sample inputs
- Verify correctness

### Day 7: Challenge
- Attempt challenge exercises
- Combine multiple concepts
- Test thoroughly

---

## ✅ Self-Assessment

Ask yourself these questions:

- [ ] Can I explain what a program is?
- [ ] Can I write pseudocode for a given problem?
- [ ] Do I understand if-else, loops, and functions?
- [ ] Can I trace through an algorithm step-by-step?
- [ ] Can I identify and fix bugs in pseudocode?
- [ ] Do I know how to approach a new problem?
- [ ] Can I write pseudocode without looking at examples?
- [ ] Can I optimize an algorithm for performance?

If you can answer YES to all, you're ready for Level 1! 🚀
