# Level 4: Operators & Expressions - INDEX

Welcome to **Level 4: Operators & Expressions**! This is where you learn to act on the variables and types you mastered in Level 3.

---

## 📖 What You'll Learn

- ✅ Arithmetic operators (+, -, *, /, %) and integer division truncation
- ✅ Comparison operators (==, !=, <, >, <=, >=)
- ✅ Logical operators (&&, ||, !) and short-circuit evaluation
- ✅ Bitwise operators (&, |, ^, <<, >>) and bit-flag systems
- ✅ Assignment operators (+=, -=, *=, /=, ++, --, etc.)
- ✅ Operator precedence
- ✅ How to build and evaluate expressions

---

## 🗂️ Level 4 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview and welcome message. Start here first!

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- Arithmetic, comparison, logical, and bitwise operators
- Assignment operators
- Operator precedence
- Expressions and short-circuit evaluation
- Best practices
- Common mistakes

**Read Time:** 45-60 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 detailed, step-by-step exercises:
1. Arithmetic operators
2. Integer division & modulo
3. Comparison operators
4. Logical operators
5. Bitwise operators (permission flags)
6. Assignment operators
7. Operator precedence
8. Short-circuit evaluation
9. Building expressions (mini calculator)
10. Comprehensive practice (BMI calculator)

**Time Commitment:** 4-5 hours
**When:** After reading README

### 4. **STUDY_GUIDE.md** (Visual Learning)
- Operator category tables
- Truth tables (logical and bitwise)
- Bit-diagram walkthroughs
- Precedence table and decision tree
- Common mistakes guide

**Read Time:** 30-40 minutes
**When:** Use while doing exercises

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- Essential syntax reference
- Precedence table
- Common operations quick guide
- Common mistakes table

**Read Time:** 10-15 minutes
**When:** Quick lookup during practice

---

## 🎯 Recommended Learning Path

### Day 1: Arithmetic & Comparison (2 hours)
1. Read **START_HERE.txt** (10 min)
2. Read **README.md** sections 1-2 (30 min)
3. Look at **STUDY_GUIDE.md** arithmetic/comparison sections (20 min)
4. Complete Exercises 1-3 (1 hour)

### Day 2: Logical & Bitwise (2 hours)
1. Read **README.md** sections 3-4 (30 min)
2. Complete Exercises 4-5 (1 hour 30 min)

### Day 3: Assignment, Precedence & Short-Circuit (2 hours)
1. Read **README.md** sections 5-8 (40 min)
2. Complete Exercises 6-8 (1 hour 20 min)

### Day 4: Expressions & Real-World Practice (1.5 hours)
1. Complete Exercises 9-10 (1.5 hours)
2. Try bonus challenges if time allows

### Day 5: Consolidation (1 hour)
1. Re-read tricky sections (precedence, short-circuit)
2. Do bonus challenges
3. Review with **QUICK_REFERENCE.md**

---

## 💡 Key Concepts At A Glance

### The Five Operator Families
```go
// Arithmetic
a + b, a - b, a * b, a / b, a % b

// Comparison (returns bool)
a == b, a != b, a < b, a > b, a <= b, a >= b

// Logical (works on bool)
a && b, a || b, !a

// Bitwise (works on integers)
a & b, a | b, a ^ b, ^a, a << n, a >> n

// Assignment
a += b, a -= b, a *= b, a /= b, a++, a--
```

### Integer Division Truncates
```go
10 / 3                 // 3, not 3.333
float64(10) / 3        // 3.3333...
```

### Precedence (Highest → Lowest)
```go
()  →  ++/--  →  * / % & << >>  →  + - | ^  →  comparisons  →  &&  →  ||
```

### Short-Circuit Evaluation
```go
if b != 0 && a/b > 10 {   // b != 0 checked FIRST; a/b never runs if b is 0
    // safe
}
```

---

## ✅ Prerequisites

Make sure you've completed **Level 3: Variables & Data Types**

You need:
- ✅ Understanding of Go's basic types (int, float64, string, bool)
- ✅ Ability to declare variables with var and :=
- ✅ Understanding of type conversion
- ✅ Familiarity with constants and iota

---

## 🎓 Learning Objectives

By the end of Level 4, you'll be able to:

- ✅ Use all arithmetic operators, including modulo
- ✅ Explain and avoid integer division truncation surprises
- ✅ Use all comparison operators on numbers and strings
- ✅ Combine conditions with logical operators
- ✅ Manipulate bits and build a permission-flag system
- ✅ Use compound assignment and increment/decrement operators
- ✅ Predict how operator precedence affects a result
- ✅ Explain and rely on short-circuit evaluation
- ✅ Build and evaluate expressions in realistic programs

---

## 📊 Statistics

- **Main Theory:** ~14,000 characters (README.md, 670 lines)
- **Exercises:** 10 detailed exercises + 3 bonus challenges
- **Study Guide:** Truth tables, bit diagrams, precedence table
- **Quick Reference:** One-page cheat sheet
- **Practice Time:** 4-5 hours
- **Difficulty:** ⭐⭐ (Intermediate)

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for comprehensive theory
2. Use STUDY_GUIDE.md alongside for truth tables and bit diagrams
3. Reference QUICK_REFERENCE.md for fast lookup

### For Practice
1. Read exercise description
2. Copy provided bash commands (or type them)
3. Follow step-by-step instructions
4. Verify output matches expected results
5. Understand what each step does

### For Review
1. Use QUICK_REFERENCE.md for quick recap
2. Refer to common mistakes section
3. Review the precedence table until it's automatic

---

## 🆘 Common Questions

**Q: Why does `10 / 3` give `3` instead of `3.33`?**
A: When both operands are integers, Go performs integer division and truncates the decimal part. Convert to `float64` first if you need a precise result.

**Q: Do I need to memorize the precedence table?**
A: Not fully — but know that `*`, `/`, `%` bind tighter than `+`, `-`, and that `&&` binds tighter than `||`. When in doubt, use parentheses.

**Q: What's short-circuit evaluation actually useful for?**
A: Writing safe guard conditions, like checking `b != 0` before dividing by `b` in the same `if`.

**Q: Why would I use bitwise operators in normal application code?**
A: Permission/flag systems (like Read/Write/Execute), efficient set membership, and low-level protocol or hardware work.

**Q: Can I use `++` inside an expression like `y := x++`?**
A: No. In Go, `++` and `--` are statements, not expressions. Increment first, then use the variable.

---

## 🎯 Before Moving to Level 5

Make sure you can answer these questions:

- [ ] Why does integer division truncate, and how do you avoid it?
- [ ] What's the difference between `&&` and `||`?
- [ ] How do bit-flags work with `|`, `&`, and `<<`?
- [ ] What are the compound assignment operators?
- [ ] What is operator precedence, and how do parentheses override it?
- [ ] What is short-circuit evaluation, and why does operand order matter?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (10 min)
   - Gets you oriented

2. **README.md** Sections 1-4 (30 min)
   - Arithmetic, comparison, logical, bitwise operators

3. **README.md** Sections 5-8 (30 min)
   - Assignment, precedence, expressions, short-circuit

4. **EXERCISES.md** Exercises 1-5 (2 hours)
   - Get hands-on with each operator family

5. **STUDY_GUIDE.md** (30 min)
   - Study the truth tables and bit diagrams

6. **EXERCISES.md** Exercises 6-10 (2+ hours)
   - Deep practice, precedence, short-circuit, real programs

7. **QUICK_REFERENCE.md** (10 min)
   - Create your mental cheat sheet

---

## 🏆 Success Indicators

You've mastered Level 4 when:

- ✅ You can use all five operator families without hesitation
- ✅ You instinctively convert to float64 before dividing when precision matters
- ✅ You can predict precedence for mixed expressions
- ✅ You can build a bit-flag system from scratch
- ✅ You order && conditions to guard unsafe operations
- ✅ You've completed 8+ exercises
- ✅ You can explain operators and expressions to someone else

---

## 🚀 What's Next?

After Level 4, you're ready for:

**Level 5: Conditions (if/else/switch)**
- `if` / `else if` / `else` statements
- `switch` statements
- Making decisions in code using the expressions you just learned

---

## 💬 Key Takeaway

> **Operators turn values into decisions. Precedence and short-circuiting determine what actually runs.**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Exercises
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Patterns
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Happy learning! Level 4 gives you the tools to build every condition and calculation ahead! 🎉

*Estimated time to complete Level 4: 4-5 hours*
*Difficulty: ⭐⭐ (Intermediate)*
*Next Level: Level 5 - Conditions (if/else/switch)*
