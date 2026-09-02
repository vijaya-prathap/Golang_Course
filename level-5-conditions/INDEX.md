# Level 5: Conditions (if/else/switch) - INDEX

Welcome to **Level 5: Conditions**! This is where your programs start making decisions using the operators and expressions you mastered in Level 4.

---

## 📖 What You'll Learn

- ✅ if, if/else, and if/else if chains
- ✅ The if-with-short-init pattern (and the comma-ok idiom)
- ✅ switch statements, including multi-value cases
- ✅ Conditionless switch as a clean if/else-if replacement
- ✅ fallthrough, and why Go's switch doesn't fall through by default
- ✅ How to flatten nested conditions with early returns

---

## 🗂️ Level 5 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview and welcome message. Start here first!

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- if, else, and else-if chains
- Short-init form for if and switch
- switch statements, conditionless switch, and fallthrough
- Nested conditions and early returns
- Best practices
- Common mistakes

**Read Time:** 40-55 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 detailed, step-by-step exercises:
1. Basic if statement
2. if/else
3. if/else if/else chains
4. if with short init (comma-ok idiom)
5. Basic switch
6. switch with short init
7. Conditionless switch
8. fallthrough
9. Avoiding deep nesting with early returns
10. Comprehensive practice (ticket pricing system)

**Time Commitment:** 4-5 hours
**When:** After reading README

### 4. **STUDY_GUIDE.md** (Visual Learning)
- if vs switch decision guide
- Short-init pattern breakdown
- switch anatomy and the four switch shapes
- fallthrough flow diagrams
- Nested-vs-flattened code comparisons
- Common mistakes guide

**Read Time:** 30-40 minutes
**When:** Use while doing exercises

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- Essential syntax reference
- if vs switch decision guide
- Common mistakes table

**Read Time:** 10-15 minutes
**When:** Quick lookup during practice

---

## 🎯 Recommended Learning Path

### Day 1: if, else, else-if (2 hours)
1. Read **START_HERE.txt** (10 min)
2. Read **README.md** sections 1-3 (30 min)
3. Complete Exercises 1-3 (1 hour 20 min)

### Day 2: Short Init & switch Basics (2 hours)
1. Read **README.md** sections 4-6 (30 min)
2. Complete Exercises 4-6 (1.5 hours)

### Day 3: Conditionless switch & fallthrough (1.5 hours)
1. Read **README.md** sections 7-8 (20 min)
2. Complete Exercises 7-8 (1 hour 10 min)

### Day 4: Nesting & Real-World Practice (1.5 hours)
1. Read **README.md** sections 9-10 (20 min)
2. Complete Exercises 9-10 (1 hour 10 min)

### Day 5: Consolidation (1 hour)
1. Try bonus challenges
2. Review with **QUICK_REFERENCE.md**
3. Make sure the if-vs-switch decision feels automatic

---

## 💡 Key Concepts At A Glance

### if / else if / else
```go
if score >= 90 {
    grade = "A"
} else if score >= 80 {
    grade = "B"
} else {
    grade = "F"
}
```

### Short Init
```go
if value, err := strconv.Atoi(s); err == nil {
    // use value
} else {
    // handle err
}
```

### switch (No Fallthrough By Default)
```go
switch day {
case "Saturday", "Sunday":
    fmt.Println("Weekend")
default:
    fmt.Println("Weekday")
}
```

### Conditionless switch
```go
switch {
case score >= 90:
    grade = "A"
case score >= 80:
    grade = "B"
default:
    grade = "F"
}
```

### Early Returns Over Nesting
```go
if !active { return "Inactive" }
if !verified { return "Not verified" }
return "OK"
```

---

## ✅ Prerequisites

Make sure you've completed **Level 4: Operators & Expressions**

You need:
- ✅ Comfort with comparison operators (==, !=, <, >, <=, >=)
- ✅ Comfort with logical operators (&&, ||, !) and short-circuit evaluation
- ✅ Understanding of operator precedence

---

## 🎓 Learning Objectives

By the end of Level 5, you'll be able to:

- ✅ Write correctly ordered if/else-if chains
- ✅ Use the short-init form to scope helper variables
- ✅ Apply the comma-ok idiom for error checks and map lookups
- ✅ Write both value-based and conditionless switch statements
- ✅ Explain why Go's switch doesn't fall through by default
- ✅ Use fallthrough deliberately, when appropriate
- ✅ Flatten deeply nested conditions with early returns

---

## 📊 Statistics

- **Main Theory:** README.md covering if/else, switch, fallthrough, and nesting
- **Exercises:** 10 detailed exercises + 3 bonus challenges
- **Study Guide:** Decision guides, pattern breakdowns, flow diagrams
- **Quick Reference:** One-page cheat sheet
- **Practice Time:** 4-5 hours
- **Difficulty:** ⭐⭐ (Intermediate)

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for comprehensive theory
2. Use STUDY_GUIDE.md alongside for decision trees and flow diagrams
3. Reference QUICK_REFERENCE.md for fast lookup

### For Practice
1. Read exercise description
2. Copy provided bash commands (or type them)
3. Follow step-by-step instructions
4. Verify output matches expected results
5. Understand what each step does

### For Review
1. Use QUICK_REFERENCE.md for quick recap
2. Refer to the common mistakes section
3. Practice choosing if vs switch until it's automatic

---

## 🆘 Common Questions

**Q: When should I use if vs switch?**
A: Use `if`/`else` for one or two conditions. Use `switch` when comparing one value against many cases, or a conditionless `switch` when you have many independent conditions/ranges.

**Q: Does Go's switch fall through like C or Java?**
A: No. Each case automatically breaks after its body runs. Use the `fallthrough` keyword to explicitly opt into the next case.

**Q: What's the short-init form actually for?**
A: Scoping a helper variable (like an error or a map lookup result) so it only exists where it's needed — it disappears once the if/else or switch block ends.

**Q: Why does order matter in if/else-if chains?**
A: Go runs the first true condition and skips the rest. For overlapping ranges (like grade thresholds), you must check the highest/most specific condition first.

**Q: Is nesting if statements bad?**
A: Not inherently, but two or three levels deep it usually gets harder to read. Prefer early returns (guard clauses) when every branch would otherwise return/exit anyway.

---

## 🎯 Before Moving to Level 6

Make sure you can answer these questions:

- [ ] Why must if/else-if thresholds be ordered from most to least specific?
- [ ] What is the short-init pattern, and where is its variable scoped?
- [ ] How does Go's switch differ from C/Java's switch regarding fallthrough?
- [ ] What does a conditionless switch look like, and when would you use one?
- [ ] How do you flatten nested conditions using early returns?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (10 min)
   - Gets you oriented

2. **README.md** Sections 1-4 (25 min)
   - if, else, else-if, short init

3. **README.md** Sections 5-8 (25 min)
   - switch, conditionless switch, fallthrough, nesting

4. **EXERCISES.md** Exercises 1-5 (2 hours)
   - Get hands-on with if/else and basic switch

5. **STUDY_GUIDE.md** (30 min)
   - Study the decision guides and flow diagrams

6. **EXERCISES.md** Exercises 6-10 (2+ hours)
   - Deep practice with switch variants and real programs

7. **QUICK_REFERENCE.md** (10 min)
   - Create your mental cheat sheet

---

## 🏆 Success Indicators

You've mastered Level 5 when:

- ✅ You order if/else-if thresholds correctly without thinking twice
- ✅ You reach for the short-init form naturally for errors and lookups
- ✅ You choose between if and switch based on the shape of the decision
- ✅ You know Go's switch doesn't fall through, and why that's safer
- ✅ You instinctively flatten nested conditions with early returns
- ✅ You've completed 8+ exercises
- ✅ You can explain conditions to someone else

---

## 🚀 What's Next?

After Level 5, you're ready for:

**Level 6: Loops**
- `for` loops (Go's only loop keyword)
- while-style and infinite loops
- `break` and `continue`
- Looping over arrays, slices, and maps

---

## 💬 Key Takeaway

> **Conditions turn expressions into decisions. Order matters, scope matters, and flat code beats nested code.**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Exercises
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Patterns
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Happy learning! Level 5 gives your programs the power to decide! 🎉

*Estimated time to complete Level 5: 4-5 hours*
*Difficulty: ⭐⭐ (Intermediate)*
*Next Level: Level 6 - Loops*
