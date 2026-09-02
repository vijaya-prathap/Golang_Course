# Level 6: Loops - INDEX

Welcome to **Level 6: Loops**! This is where your programs start repeating work instead of executing top-to-bottom just once.

---

## 📖 What You'll Learn

- ✅ The three-component for loop
- ✅ The while-style and infinite for loop forms
- ✅ break and continue
- ✅ range over slices, arrays, strings, and maps
- ✅ Labeled break and continue for nested loops
- ✅ Common loop patterns: accumulate, find, filter, count

---

## 🗂️ Level 6 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview and welcome message. Start here first!

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- All three for loop shapes
- break, continue, and labeled loop control
- range over every collection type
- Nested loops and common patterns
- Best practices
- Common mistakes

**Read Time:** 40-55 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 detailed, step-by-step exercises:
1. Three-component for loop
2. while-style for loop
3. Infinite loop with break
4. continue
5. range over slices
6. range over strings (runes)
7. range over maps
8. Labeled break and continue
9. Common loop patterns
10. Comprehensive practice (prime number finder)

**Time Commitment:** 4-5 hours
**When:** After reading README

### 4. **STUDY_GUIDE.md** (Visual Learning)
- The three for loop shapes compared
- range cheat sheet by collection type
- Rune iteration diagrams
- Labeled break/continue decision tree
- Common patterns and mistakes guide

**Read Time:** 30-40 minutes
**When:** Use while doing exercises

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- Essential syntax reference
- range syntax by collection type
- Common patterns quick guide
- Common mistakes table

**Read Time:** 10-15 minutes
**When:** Quick lookup during practice

---

## 🎯 Recommended Learning Path

### Day 1: The for Loop Shapes (2 hours)
1. Read **START_HERE.txt** (10 min)
2. Read **README.md** sections 1-3 (30 min)
3. Complete Exercises 1-3 (1 hour 20 min)

### Day 2: break, continue & range Basics (2 hours)
1. Read **README.md** sections 4-5 (25 min)
2. Complete Exercises 4-6 (1.5 hours)

### Day 3: range Over Maps & Labeled Loops (1.5 hours)
1. Read **README.md** sections 6-7 (25 min)
2. Complete Exercises 7-8 (1 hour 5 min)

### Day 4: Patterns & Real-World Practice (1.5 hours)
1. Read **README.md** sections 8-10 (20 min)
2. Complete Exercises 9-10 (1 hour 10 min)

### Day 5: Consolidation (1 hour)
1. Try bonus challenges
2. Review with **QUICK_REFERENCE.md**
3. Make sure range's behavior for strings and maps feels automatic

---

## 💡 Key Concepts At A Glance

### Three for Loop Shapes
```go
for i := 0; i < 5; i++ { }   // three-component
for count > 0 { }             // while-style
for { }                       // infinite (needs break)
```

### break and continue
```go
if done { break }       // exit the loop
if skip { continue }    // skip to next iteration
```

### range Over Collections
```go
for i, v := range slice { }    // index + value
for i, r := range someString { } // byte offset + rune
for k, v := range someMap { }  // key + value (order randomized)
```

### Labeled break/continue
```go
outer:
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if j == 1 { break outer }
        }
    }
```

---

## ✅ Prerequisites

Make sure you've completed **Level 5: Conditions (if/else/switch)**

You need:
- ✅ Comfort with if/else and switch statements
- ✅ Comfort with comparison and logical operators
- ✅ Basic familiarity with slices and maps (introduced in Level 3)

---

## 🎓 Learning Objectives

By the end of Level 6, you'll be able to:

- ✅ Write all three for loop forms
- ✅ Use break and continue correctly
- ✅ Range over slices, arrays, strings, and maps
- ✅ Explain why range over a string gives runes, not bytes
- ✅ Explain why map iteration order is randomized, and sort keys when needed
- ✅ Use labeled break/continue to control nested loops
- ✅ Apply accumulate, find, filter, and count patterns fluently

---

## 📊 Statistics

- **Main Theory:** README.md covering all loop forms, range, and common patterns
- **Exercises:** 10 detailed exercises + 3 bonus challenges
- **Study Guide:** range cheat sheet, rune diagrams, pattern library
- **Quick Reference:** One-page cheat sheet
- **Practice Time:** 4-5 hours
- **Difficulty:** ⭐⭐ (Intermediate)

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for comprehensive theory
2. Use STUDY_GUIDE.md alongside for the range cheat sheet and diagrams
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
3. Practice range's behavior for strings/maps until it's automatic

---

## 🆘 Common Questions

**Q: Does Go have a while loop?**
A: Not as a separate keyword — `for condition { }` (dropping the init/post clauses) IS Go's while loop.

**Q: Why does ranging over a string give weird jumping indices?**
A: `range` decodes UTF-8 into runes (Unicode code points). Multi-byte characters (like `é`) make the byte-offset index skip ahead by more than 1.

**Q: Why is my map printing in a different order every time I run it?**
A: Go deliberately randomizes map iteration order to stop code from silently depending on an order that was never guaranteed. Sort the keys if you need a stable order.

**Q: What's the difference between break and a labeled break?**
A: A plain `break` only exits the innermost loop. A labeled `break <label>` exits the labeled (usually outer) loop, even from inside nested loops.

**Q: When should I use range vs a manual index loop?**
A: Use `range` almost always — it's more idiomatic and less error-prone. Reach for a manual `for i := 0; i < n; i++` only when you need fine control (e.g., skipping indices, moving backward, or comparing two positions).

---

## 🎯 Before Moving to Level 7

Make sure you can answer these questions:

- [ ] What are the three for loop shapes, and when would you use each?
- [ ] What's the difference between break and continue?
- [ ] Why does range over a string give runes instead of bytes?
- [ ] Why is map iteration order randomized, and how do you get a stable order?
- [ ] How does a labeled break differ from a plain break in nested loops?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (10 min)
   - Gets you oriented

2. **README.md** Sections 1-4 (25 min)
   - The three for shapes, break/continue

3. **README.md** Sections 5-8 (30 min)
   - range, labels, nesting, patterns

4. **EXERCISES.md** Exercises 1-5 (2 hours)
   - Get hands-on with loop shapes and slice iteration

5. **STUDY_GUIDE.md** (30 min)
   - Study the range cheat sheet and rune diagrams

6. **EXERCISES.md** Exercises 6-10 (2+ hours)
   - Deep practice with strings, maps, labels, and real programs

7. **QUICK_REFERENCE.md** (10 min)
   - Create your mental cheat sheet

---

## 🏆 Success Indicators

You've mastered Level 6 when:

- ✅ You reach for range by default and only use manual indices when you need to
- ✅ You can explain rune vs byte iteration for strings without hesitation
- ✅ You never assume map iteration order is stable
- ✅ You can use a labeled break/continue to control nested loops
- ✅ You recognize accumulate/find/filter/count patterns instantly
- ✅ You've completed 8+ exercises
- ✅ You can explain loops to someone else

---

## 🚀 What's Next?

After Level 6, you're ready for:

**Level 7: Arrays**
- Fixed-size arrays in depth
- Array vs slice semantics (value vs reference-like behavior)
- Multi-dimensional arrays
- When to use arrays vs slices

---

## 💬 Key Takeaway

> **`for` is Go's only loop keyword, but range makes most loops write themselves — just watch out for rune iteration and map ordering.**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Exercises
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Patterns
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Happy learning! Level 6 gives your programs the power to repeat and process data at scale! 🎉

*Estimated time to complete Level 6: 4-5 hours*
*Difficulty: ⭐⭐ (Intermediate)*
*Next Level: Level 7 - Arrays*
