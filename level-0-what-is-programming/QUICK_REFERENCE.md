# Level 0: Quick Reference Card

## 📌 Print This Out!

This is a handy reference you can print or keep on your phone.

---

## 🎯 What is Programming?

**Definition:** A sequence of instructions that tells a computer what to do.

**Goal:** Solve problems and automate tasks.

**Key Principle:** Input → Process → Output

---

## 🧩 Four Core Concepts

### 1. SEQUENCE
```
Do Step 1
Do Step 2
Do Step 3
```
Execute one instruction after another in order.

### 2. SELECTION
```
IF condition is true
    Do Action A
ELSE
    Do Action B
```
Choose which action based on a condition.

### 3. REPETITION
```
FOR i = 1 TO n
    Do Action
REPEAT
```
Perform the same action multiple times.

### 4. FUNCTIONS
```
FUNCTION doSomething(input):
    RETURN output
```
Group actions into reusable blocks.

---

## 💾 Data Types Cheat Sheet

| Type | Use For | Examples |
|------|---------|----------|
| Integer | Whole numbers | -5, 0, 42 |
| Float | Decimals | 3.14, -2.5 |
| String | Text | "Hello", "Golang" |
| Boolean | True/False | true, false |

---

## 📊 Operators Quick Reference

### Arithmetic
```
+  Add      5 + 3 = 8
-  Subtract 5 - 3 = 2
*  Multiply 5 * 3 = 15
/  Divide   6 / 3 = 2
%  Modulo   7 % 3 = 1 (remainder)
```

### Comparison
```
=  Equal     5 = 5 (true)
≠  Not equal 5 ≠ 3 (true)
>  Greater   5 > 3 (true)
<  Less      5 < 3 (false)
≥  Greater/equal 5 ≥ 5 (true)
≤  Less/equal    3 ≤ 5 (true)
```

### Logical
```
AND  Both true      (5>3) AND (3>1) = true
OR   At least one   (5>3) OR (1>3) = true
NOT  Reverse        NOT (5>3) = false
```

---

## 🔄 Control Structures

### IF-ELSE
```
IF condition THEN
    Action 1
ELSE IF another_condition THEN
    Action 2
ELSE
    Action 3
END IF
```

### FOR LOOP
```
FOR i = start TO end STEP increment DO
    Action (repeated)
END FOR
```

### WHILE LOOP
```
WHILE condition DO
    Action
END WHILE
```

### SWITCH
```
SWITCH value
    CASE option1: Action1
    CASE option2: Action2
    DEFAULT: Action3
END SWITCH
```

---

## 🎯 Problem-Solving Template

### Step 1: Understand
- Input: What information is given?
- Output: What should the result be?
- Constraints: Any limitations?

### Step 2: Plan
```
Write pseudocode:
1. First action
2. Second action
3. Third action
```

### Step 3: Verify
- Does it handle edge cases?
- Will it always terminate?
- Is the logic correct?

### Step 4: Test
- Try with normal inputs
- Try with edge cases
- Try with invalid inputs

---

## ⚡ Algorithm Patterns

### Pattern: Find Maximum
```
max = first
FOR each element:
    IF element > max:
        max = element
RETURN max
```

### Pattern: Count Items
```
count = 0
FOR each item:
    IF matches:
        count = count + 1
RETURN count
```

### Pattern: Sum All
```
sum = 0
FOR each element:
    sum = sum + element
RETURN sum
```

### Pattern: Check All
```
FOR each item:
    IF not valid:
        RETURN false
RETURN true
```

---

## 🐛 Debugging Checklist

- [ ] Does the algorithm handle empty input?
- [ ] Does it handle single item?
- [ ] Does it handle boundary values?
- [ ] Are all variables initialized?
- [ ] Will the loop terminate?
- [ ] Are the conditions correct?
- [ ] Does it handle errors gracefully?

---

## 💡 Best Practices

✅ DO:
- Use descriptive variable names
- Add comments explaining WHY
- Test with multiple inputs
- Break problems into smaller pieces
- Trace through with sample data

❌ DON'T:
- Use cryptic names (x, temp, a1)
- Write one huge function
- Ignore edge cases
- Skip testing
- Over-complicate solutions

---

## 🎓 Common Mistakes

| Mistake | Fix |
|---------|-----|
| Infinite loop | Make sure condition eventually becomes false |
| Off-by-one | Check array boundaries carefully |
| Wrong output | Trace through step-by-step manually |
| Logic error | Test with multiple inputs |
| Uninitialized var | Always assign initial value |
| Division by zero | Check denominator first |

---

## 📚 Pseudocode Template

```
ALGORITHM AlgorithmName(input1, input2):
    // Description of what this does
    
    // Initialize variables
    variable1 = initial_value
    variable2 = initial_value
    
    // Main logic
    IF condition THEN
        action1
    ELSE
        action2
    END IF
    
    FOR i = 1 TO n DO
        action3
    END FOR
    
    // Return result
    RETURN result
    
END ALGORITHM
```

---

## 🚀 Performance Notes

Faster algorithms win when dealing with large amounts of data:

```
O(1)     - Instant           ⚡⚡⚡ (Best)
O(log n) - Very fast          ⚡⚡
O(n)     - Fast               ⚡
O(n²)    - Slow              🐢
O(2^n)   - Very slow         🐢🐢 (Avoid)
```

---

## 📋 Exercise Quick Links

1. **Simple Calculation** - Math operations, error handling
2. **Number Classification** - If-else conditions
3. **Grade Calculator** - Switch or nested if-else
4. **Fibonacci** - Loops and pattern recognition
5. **Prime Checker** - Logic optimization (sqrt trick)
6. **Reverse String** - String manipulation
7. **Palindrome** - Two-pointer technique
8. **Count Occurrences** - Loop and condition

---

## ✅ Level 0 Mastery Checklist

- [ ] Explain what a program is
- [ ] Draw flowcharts for algorithms
- [ ] Write pseudocode without syntax errors
- [ ] Trace through algorithms manually
- [ ] Identify and fix logic bugs
- [ ] Choose correct control structures
- [ ] Design efficient algorithms
- [ ] Handle edge cases properly
- [ ] Complete all 8 exercises
- [ ] Explain solution approaches

---

## 🎯 Key Takeaway

**"The art of programming is about clear thinking, not complex syntax."**

Master the logic first. The syntax comes later.

---

## 📞 Need Help?

When stuck:
1. Re-read the problem carefully
2. Break into smaller steps
3. Write pseudocode first
4. Trace through with sample data
5. Check for edge cases
6. Compare with the solution

**Never copy without understanding!**

---

## 🎉 Next Steps

You're ready for **Level 1: Go Introduction & Environment** when:
- ✅ You understand all core concepts
- ✅ You can write pseudocode fluently
- ✅ You've completed all exercises

See you in Level 1! 🚀
