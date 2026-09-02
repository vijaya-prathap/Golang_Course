# Level 7: Arrays - INDEX

Welcome to **Level 7: Arrays**! This is where you learn Go's fixed-size collection type - and the value-semantics rule that explains a huge amount of "why does Go behave this way" once you reach slices.

---

## 📖 What You'll Learn

- ✅ What makes arrays fixed-size and homogeneous
- ✅ Every way to declare an array: var, literals, [...], and sparse index literals
- ✅ Indexing, len(), zero values, and bounds checking
- ✅ Why arrays are value types - assignment and function calls both COPY them
- ✅ Comparing arrays with == and using them as map keys
- ✅ Multi-dimensional arrays for grid-shaped data
- ✅ When to actually choose an array over a slice (it's rarer than you'd think)

---

## 🗂️ Level 7 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview and welcome message. Start here first!

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- What arrays are and how they differ from slices
- Every array declaration form
- Indexing, len(), and zero values
- Value semantics (the big idea) with verified, real output
- Array comparison and multi-dimensional arrays
- range over arrays
- Array vs slice decision-making
- Best practices
- Common mistakes

**Read Time:** 40-55 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 detailed, step-by-step exercises:
1. Declaring arrays
2. Indexing, len(), and zero values
3. Arrays are value types - assignment copies
4. Arrays are value types - function parameters
5. Comparing arrays with ==
6. Multi-dimensional arrays
7. Iterating arrays with range
8. Array vs slice - converting and sharing memory
9. Tic-tac-toe board (2D array project)
10. Comprehensive practice (Sudoku row validator & scoreboard)

**Time Commitment:** 4-5 hours
**When:** After reading README

### 4. **STUDY_GUIDE.md** (Visual Learning)
- Array declaration forms table
- The array copy diagram (two separate memory blocks, contrasted with slices)
- Multi-dimensional array indexing diagram
- Array vs slice decision tree
- Common mistakes guide

**Read Time:** 30-40 minutes
**When:** Use while doing exercises

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- Essential syntax reference
- The one rule that matters: arrays copy, slices share
- Common mistakes table

**Read Time:** 10-15 minutes
**When:** Quick lookup during practice

---

## 🎯 Recommended Learning Path

### Day 1: Declaring and Indexing Arrays (2 hours)
1. Read **START_HERE.txt** (10 min)
2. Read **README.md** sections 1-3 (30 min)
3. Complete Exercises 1-2 (1 hour 20 min)

### Day 2: Value Semantics - The Big Idea (2 hours)
1. Read **README.md** section 4 carefully (25 min)
2. Complete Exercises 3-4 (1.5 hours) - run the code, don't just read it

### Day 3: Comparison & Multi-Dimensional Arrays (1.5 hours)
1. Read **README.md** sections 5-6 (25 min)
2. Complete Exercises 5-6 (1 hour 5 min)

### Day 4: range, Array vs Slice, & Real Projects (1.5 hours)
1. Read **README.md** sections 7-8 (20 min)
2. Complete Exercises 7-10 (1 hour 10 min)

### Day 5: Consolidation (1 hour)
1. Try bonus challenges
2. Review with **QUICK_REFERENCE.md**
3. Make sure you can PROVE (not just recite) that arrays copy on assignment and function calls

---

## 💡 Key Concepts At A Glance

### Declaration Forms
```go
var arr [5]int                  // zero-valued
lit := [3]string{"a", "b", "c"} // explicit size
inf := [...]int{1, 2, 3}        // size inferred
sparse := [5]int{0: 1, 4: 5}    // sparse indices
```

### The Big Idea: Arrays Copy
```go
a := [3]int{1, 2, 3}
b := a          // COPY
b[0] = 99
// a is [1 2 3] - completely unaffected
```

### Comparison
```go
[3]int{1,2,3} == [3]int{1,2,3}   // true
[3]int{1,2,3} == [4]int{1,2,3,4} // compile error - different types!
```

### Multi-Dimensional
```go
var grid [3][3]int
grid[1][2] = 7   // grid[row][col]
```

---

## ✅ Prerequisites

Make sure you've completed **Level 6: Loops**

You need:
- ✅ Comfort with all three for loop forms
- ✅ Comfort with range over slices, strings, and maps
- ✅ Basic familiarity with slices (introduced informally in earlier levels)

---

## 🎓 Learning Objectives

By the end of Level 7, you'll be able to:

- ✅ Declare arrays using every available syntax
- ✅ Index, assign, and use len() correctly, including on multi-dimensional arrays
- ✅ Explain AND demonstrate why arrays are value types
- ✅ Predict exactly what happens when an array is assigned or passed to a function
- ✅ Compare arrays with == and use them as map keys
- ✅ Build and iterate multi-dimensional arrays
- ✅ Decide confidently when an array is the right tool instead of a slice

---

## 📊 Statistics

- **Main Theory:** README.md covering declaration, value semantics, comparison, and multi-dimensional arrays
- **Exercises:** 10 detailed exercises + 3 bonus challenges
- **Study Guide:** copy-semantics diagram, indexing diagram, decision tree
- **Quick Reference:** One-page cheat sheet
- **Practice Time:** 4-5 hours
- **Difficulty:** ⭐⭐ (Intermediate)

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for comprehensive theory
2. Use STUDY_GUIDE.md alongside for the copy-semantics diagram
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
3. Practice explaining the array-copy behavior out loud until it's automatic

---

## 🆘 Common Questions

**Q: Why does Go even have arrays if slices are more flexible?**
A: Sometimes a fixed size IS the correct guarantee - a chess board is always 8x8, a hash is always a fixed number of bytes. Arrays also being comparable and usable as map keys is something slices can never do. Plus, slices are literally built on top of arrays.

**Q: If I assign one array to another, does it create a reference like a slice?**
A: No. Array assignment always copies every element. This is true for `:=`/`=` assignment AND for passing an array as a function argument. Slices are the ones that share memory.

**Q: Are `[3]int` and `[4]int` really different types?**
A: Yes - the size is part of the type. You can't assign one to the other, compare them with ==, or pass one where the other is expected. This is different from, say, a slice, where `len()` can change freely.

**Q: Why did my out-of-range access get caught at compile time instead of panicking?**
A: If the index is a literal constant and the array's size is known, the Go compiler checks the bounds itself. Only when the index comes from a variable (not known until runtime) does it become a runtime panic instead.

**Q: Should I use arrays or slices for my function parameters?**
A: Slices, almost always. Arrays as parameters get fully copied on every call and can't be resized - that's rarely what you want for general-purpose code.

---

## 🎯 Before Moving to Level 8

Make sure you can answer these questions:

- [ ] What are all the ways to declare an array?
- [ ] What happens when you assign one array to another? Can you prove it with real code?
- [ ] What happens when you pass an array to a function? How is that different from passing a slice?
- [ ] When can two arrays be compared with ==, and when is it a compile error?
- [ ] How do you declare, fill, and index a multi-dimensional array?
- [ ] When would you actually choose an array over a slice?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (10 min)
   - Gets you oriented

2. **README.md** Sections 1-4 (30 min)
   - Declaration forms and the value-semantics rule

3. **README.md** Sections 5-8 (30 min)
   - Comparison, multi-dimensional arrays, range, array vs slice

4. **EXERCISES.md** Exercises 1-4 (2 hours)
   - Get hands-on with declaration and prove value semantics yourself

5. **STUDY_GUIDE.md** (30 min)
   - Study the copy-semantics diagram and indexing diagram

6. **EXERCISES.md** Exercises 5-10 (2+ hours)
   - Deep practice with comparison, multi-dimensional arrays, and real projects

7. **QUICK_REFERENCE.md** (10 min)
   - Create your mental cheat sheet

---

## 🏆 Success Indicators

You've mastered Level 7 when:

- ✅ You can declare an array five different ways without hesitation
- ✅ You can explain - and demonstrate with running code - why arrays copy on assignment and function calls
- ✅ You never assume a function can mutate an array parameter in place
- ✅ You can compare arrays with == and explain the size/type rule
- ✅ You can build and index a multi-dimensional array
- ✅ You've completed 8+ exercises
- ✅ You can explain arrays to someone else, including why Go still has slices

---

## 🚀 What's Next?

After Level 7, you're ready for:

**Level 8: Slices**
- Slices in depth: the pointer/length/capacity header
- append, growth, and re-slicing
- Why slices share memory - and when that surprises you

---

## 💬 Key Takeaway

> **Arrays are values: assigning one, or passing one to a function, always copies every element. Slices share memory instead - and that single difference explains most of what feels surprising about Go's collections.**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Exercises
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Patterns
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Happy learning! Level 7 gives you the foundation that makes slices in Level 8 make sense! 🎉

*Estimated time to complete Level 7: 4-5 hours*
*Difficulty: ⭐⭐ (Intermediate)*
*Next Level: Level 8 - Slices*
