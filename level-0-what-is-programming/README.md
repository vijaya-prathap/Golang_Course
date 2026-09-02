# Level 0: What is Programming?

## Overview
This foundational level covers the fundamental concepts and principles of programming. Understanding these concepts is essential before diving into any programming language.

---

## 📚 Core Concepts

### 1. What is a Program?
A program is a set of instructions written in a specific order that a computer follows to perform a task. Programs solve problems and automate processes.

**Key Points:**
- Programs are **deterministic** - same input produces same output
- Programs follow a **sequence of steps**
- Programs process **input** and produce **output**

**Real-world Example:**
```
User Input: "5"
↓
Program: Multiply by 2
↓
Output: "10"
```

---

### 2. What is a Programming Language?
A programming language is a formal language with a set of rules (syntax) and meanings (semantics) used to communicate instructions to a computer.

**Types of Programming Languages:**

#### High-Level vs Low-Level
| Aspect | High-Level (Go, Python, Java) | Low-Level (Assembly, Machine Code) |
|--------|------|------|
| Human Readable | ✅ Yes | ❌ No |
| Abstraction | ✅ High | ❌ Low |
| Learning Curve | ✅ Easy | ❌ Steep |
| Performance | ⚠️ Slower | ✅ Faster |

#### Compiled vs Interpreted
| Aspect | Compiled (Go, C, C++) | Interpreted (Python, JavaScript) |
|--------|------|------|
| Compilation | Required | Not Required |
| Execution | Fast | Slower |
| Debugging | Harder | Easier |
| Portability | Platform-specific | Platform-independent |

**Go's Position:** Go is a compiled, statically-typed, high-level language.

---

### 3. Basic Programming Logic

#### Sequence
Instructions executed in order from top to bottom.
```
1. Start
2. Read input
3. Process data
4. Display output
5. End
```

#### Selection (Conditional Logic)
Making decisions based on conditions.
```
IF condition is true
   DO action A
ELSE
   DO action B
```

#### Repetition (Loops)
Repeating a block of code multiple times.
```
WHILE condition is true
   DO action
```

#### Modularization (Functions)
Breaking code into reusable blocks.
```
FUNCTION doSomething(input):
   RETURN result
```

---

### 4. Algorithms
An algorithm is a step-by-step procedure to solve a problem.

**Example: Finding the largest number**
```
1. Read first number, store as largest
2. Read next number
3. If new number > largest, update largest
4. Repeat steps 2-3 until all numbers read
5. Display largest
```

**Algorithm Characteristics:**
- **Finite:** Must have a definite end
- **Definite:** Each step is clear and unambiguous
- **Effective:** Steps must be executable
- **Input:** May have zero or more inputs
- **Output:** Must produce at least one output

---

### 5. Data & Variables
Data is information (numbers, text, true/false). Variables are containers that store data.

**Variable Naming Convention:**
- Use descriptive names: `age`, `totalPrice`, `isActive`
- Avoid cryptic names: `x`, `temp`, `a1`

**Data Types:**
| Type | Purpose | Examples |
|------|---------|----------|
| Integer | Whole numbers | -5, 0, 42 |
| Float | Decimal numbers | 3.14, -2.5 |
| String | Text | "Hello", "John" |
| Boolean | True/False | true, false |

---

### 6. Control Flow

#### If-Else Statement
```
IF age >= 18 THEN
    PRINT "You are an adult"
ELSE
    PRINT "You are a minor"
END IF
```

#### Switch Statement
```
SWITCH day:
    CASE Monday: PRINT "Start of work week"
    CASE Friday: PRINT "Almost weekend"
    DEFAULT: PRINT "Regular day"
END SWITCH
```

#### Loops
```
FOR i = 1 TO 10 DO
    PRINT i
END FOR

WHILE condition DO
    action
END WHILE
```

---

### 7. Functions/Procedures
Functions are reusable blocks of code that perform specific tasks.

**Benefits:**
- **Reusability:** Write once, use many times
- **Modularity:** Break complex problems into smaller pieces
- **Maintainability:** Easier to update and debug
- **Readability:** Makes code clearer

**Anatomy:**
```
FUNCTION add(a, b):
    result = a + b
    RETURN result
END FUNCTION

sum = add(5, 3)  // Calling the function
PRINT sum        // Output: 8
```

---

### 8. Problem-Solving Steps

### Step 1: Understand the Problem
- What is the input?
- What is the desired output?
- What are the constraints?

### Step 2: Plan the Solution
- Break into smaller steps
- Think about edge cases
- Design the algorithm

### Step 3: Code the Solution
- Write the instructions
- Follow language syntax
- Test as you go

### Step 4: Test & Debug
- Test with multiple inputs
- Handle edge cases
- Fix bugs

---

## 🧠 Key Principles

### 1. DRY (Don't Repeat Yourself)
Avoid writing the same code multiple times. Use functions and loops.

### 2. KISS (Keep It Simple, Stupid)
Write simple, readable code. Avoid unnecessary complexity.

### 3. Single Responsibility
Each function/module should do one thing well.

### 4. Comments & Documentation
Explain "why" not "what" (code shows what).

### 5. Testing
Test your code with various inputs to ensure correctness.

---

## 📊 Algorithm Complexity (Introduction)

### Time Complexity
How long an algorithm takes (in terms of operations, not seconds).

| Notation | Name | Example |
|----------|------|---------|
| O(1) | Constant | Direct array access |
| O(n) | Linear | Scanning through array |
| O(n²) | Quadratic | Nested loops |
| O(log n) | Logarithmic | Binary search |
| O(n log n) | Linearithmic | Efficient sorting |

### Space Complexity
How much memory an algorithm uses.

---

## 💡 Pseudocode Examples

### Example 1: Sum of Numbers
```
ALGORITHM SumNumbers(n):
    sum = 0
    FOR i = 1 TO n DO
        sum = sum + i
    END FOR
    RETURN sum
END ALGORITHM

Result: SumNumbers(5) = 1+2+3+4+5 = 15
```

### Example 2: Find Maximum
```
ALGORITHM FindMax(array):
    max = array[0]
    FOR i = 1 TO length(array)-1 DO
        IF array[i] > max THEN
            max = array[i]
        END IF
    END FOR
    RETURN max
END ALGORITHM
```

### Example 3: Factorial
```
ALGORITHM Factorial(n):
    IF n = 0 OR n = 1 THEN
        RETURN 1
    END IF
    RETURN n × Factorial(n-1)
END ALGORITHM

Result: Factorial(5) = 5 × 4 × 3 × 2 × 1 = 120
```

---

## 🎯 Common Mistakes to Avoid

1. **Infinite Loops:** Make sure loop conditions eventually become false
2. **Off-by-One Errors:** Carefully handle array/list boundaries
3. **Uninitialized Variables:** Always set initial values
4. **Logic Errors:** Test edge cases (empty input, negative numbers, etc.)
5. **Division by Zero:** Always check denominators
6. **Type Mismatches:** Use correct data types for operations

---

## 📝 Exercises

### Exercise 1: Simple Calculation
**Problem:** Write an algorithm that takes two numbers and returns:
- Sum
- Difference
- Product
- Quotient

**Challenge:** Handle division by zero

### Exercise 2: Number Classification
**Problem:** Write an algorithm that classifies a number as:
- Positive
- Negative
- Zero

### Exercise 3: Grade Calculator
**Problem:** Write an algorithm that takes a score (0-100) and returns:
- A: 90-100
- B: 80-89
- C: 70-79
- D: 60-69
- F: Below 60

### Exercise 4: Fibonacci Sequence
**Problem:** Write an algorithm that generates the first n numbers in the Fibonacci sequence.
```
Fibonacci: 0, 1, 1, 2, 3, 5, 8, 13, ...
Each number = previous two numbers added
```

### Exercise 5: Prime Number Checker
**Problem:** Write an algorithm that checks if a number is prime.
```
Prime: Only divisible by 1 and itself
Examples: 2, 3, 5, 7, 11, 13
```

### Exercise 6: Reverse a String
**Problem:** Write an algorithm that reverses a string.
```
Input: "Hello"
Output: "olleH"
```

### Exercise 7: Palindrome Checker
**Problem:** Write an algorithm that checks if a string is a palindrome.
```
Palindrome: reads same forwards and backwards
Examples: "racecar", "madam", "noon"
```

### Exercise 8: Count Occurrences
**Problem:** Write an algorithm that counts how many times a character appears in a string.
```
Input: "programming", character 'r'
Output: 2
```

---

## ✅ Level 0 Checklist

- [ ] Understand what a program is
- [ ] Know the difference between compiled and interpreted languages
- [ ] Understand sequence, selection, repetition, modularization
- [ ] Can write simple pseudocode
- [ ] Understand variables and data types
- [ ] Know how to structure control flow
- [ ] Understand functions and their benefits
- [ ] Can follow problem-solving steps
- [ ] Completed all exercises
- [ ] Can explain programming concepts in simple terms

---

## 🎓 Next Steps
Once you've mastered these foundational concepts, you'll be ready for **Level 1: Go Introduction & Environment**, where we'll set up Go and write your first program!

---

## 📚 Additional Resources
- Pseudocode is language-independent
- These concepts apply to ANY programming language
- Practice thinking algorithmically before coding
- Focus on logic, not syntax at this stage
