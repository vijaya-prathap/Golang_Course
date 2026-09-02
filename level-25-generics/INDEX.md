# Level 25: Generics - INDEX

Welcome to **Level 25: Generics**! This is where you learn to write one function or type that works safely across many types - no more copy-pasting `SumInts`/`SumFloats`, and no more falling back to `any` and losing the compiler's help.

---

## 📖 What You'll Learn

- ✅ Why generics exist: duplication vs. `any`'s lost type safety
- ✅ Type parameter syntax, calling with inference and with explicit type arguments
- ✅ The built-in `any` and `comparable` constraints
- ✅ Defining your own constraints locally with type sets (no external packages)
- ✅ Generic functions: `Map`/`Filter`/`Reduce`, revisited from Level 11
- ✅ Generic types: `Stack[T]`, `Pair[K, V]`, and a `Cache[K, V]`
- ✅ Type inference rules - when it works, and when you must specify explicitly
- ✅ When a small interface (Level 15) is the better tool than a generic type parameter

---

## 🗂️ Level 25 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview and welcome message. Start here first!

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- The duplication-vs-`any` problem generics solve
- Type parameter syntax and inference
- `any`, `comparable`, and custom type-set constraints
- Generic functions and generic types with methods
- When NOT to use generics (revisiting Level 15)
- Best practices
- Common mistakes (with real compiler errors)

**Read Time:** 45-60 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 detailed, step-by-step exercises:
1. Before/After - duplication vs. any vs. generics
2. Type parameter syntax and inference
3. The any constraint
4. The comparable constraint (+ real compile error)
5. A custom Number constraint (+ real compile error)
6. Map, Filter, Reduce - revisited from Level 11
7. A generic Stack[T]
8. A generic Pair[K, V]
9. Type inference - when it works and when it fails (+ real compile error)
10. Comprehensive practice - a generic Cache[K, V]

Plus 3 bonus challenges: generic binary search, a generic linked list, and generic Zip.

**Time Commitment:** 5-6 hours
**When:** After reading README

### 4. **STUDY_GUIDE.md** (Visual Learning)
- Type parameter syntax anatomy diagram
- Constraint hierarchy: any → comparable → custom
- Generics-vs-interfaces decision tree (ties back to Level 15)
- Generic Stack[T] memory/operation diagram
- Common mistakes with real compiler errors

**Read Time:** 35-45 minutes
**When:** Use while doing exercises

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- Essential syntax reference
- Constraint quick guide
- Common mistakes table with real compile errors

**Read Time:** 10-15 minutes
**When:** Quick lookup during practice

---

## 🎯 Recommended Learning Path

### Day 1: The Problem and the Syntax (2 hours)
1. Read **START_HERE.txt** (10 min)
2. Read **README.md** sections 1-2 (30 min)
3. Complete Exercises 1-2 (1 hour 20 min)

### Day 2: Constraints (2.5 hours)
1. Read **README.md** sections 3-5 (35 min)
2. Complete Exercises 3-5, including both real compile errors (2 hours)

### Day 3: Generic Functions & Types (2 hours)
1. Read **README.md** sections 6-7 (25 min)
2. Complete Exercises 6-8 (1.5 hours)

### Day 4: Inference & Judgment (1.5 hours)
1. Read **README.md** sections 8-9 (20 min)
2. Complete Exercise 9, including the real compile error (1 hour 10 min)

### Day 5: Comprehensive Practice & Consolidation (1.5 hours)
1. Complete Exercise 10 (45 min)
2. Try bonus challenges
3. Review with **QUICK_REFERENCE.md**

---

## 💡 Key Concepts At A Glance

### Type Parameter Syntax
```go
func Map[T, U any](s []T, f func(T) U) []U { }
```

### Constraints
```go
any                                          // anything
comparable                                   // == / != / map keys
type Number interface { int | int64 | float64 }  // your own type set
```

### Generic Types
```go
type Stack[T any] struct { items []T }
func (s *Stack[T]) Push(v T) { }
```

### Inference
```go
Map(nums, toString)          // inferred
Map[int, string](nums, toString)  // explicit (always valid)
zero := Zero[int]()          // required when T is only in the return type
```

---

## ✅ Prerequisites

Make sure you've completed **Level 24: Context**

You need:
- ✅ Comfort with functions as values and higher-order functions (Level 11)
- ✅ Comfort with interfaces and small-interface design (Level 15)
- ✅ Comfort with structs and methods (Levels 13-14)
- ✅ Basic familiarity with maps and slices (Levels 8, 10)

---

## 🎓 Learning Objectives

By the end of Level 25, you'll be able to:

- ✅ Explain the exact problem generics solve (duplication vs. lost type safety)
- ✅ Declare and call generic functions with and without explicit type arguments
- ✅ Use `any` and `comparable` correctly, and know which one a function needs
- ✅ Define your own constraints locally with type sets, fully offline
- ✅ Rewrite Level 11's `Map`/`Filter`/`Reduce` as truly generic functions
- ✅ Build generic types (`Stack[T]`, `Pair[K, V]`, `Cache[K, V]`) with methods
- ✅ Predict when type inference works and when it fails
- ✅ Decide between a generic type parameter and a small interface (Level 15)

---

## 📊 Statistics

- **Main Theory:** README.md covering the why, syntax, constraints, generic functions/types, and design judgment
- **Exercises:** 10 detailed exercises + 3 bonus challenges
- **Study Guide:** syntax anatomy, constraint hierarchy, decision tree, memory diagram
- **Quick Reference:** One-page cheat sheet
- **Practice Time:** 5-6 hours
- **Difficulty:** ⭐⭐⭐⭐ (Advanced)

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for comprehensive theory
2. Use STUDY_GUIDE.md alongside for the constraint hierarchy and decision tree
3. Reference QUICK_REFERENCE.md for fast lookup

### For Practice
1. Read exercise description
2. Copy provided bash commands (or type them)
3. Follow step-by-step instructions
4. Verify output matches expected results - including the real compiler errors in Exercises 4, 5, and 9
5. Understand what each step does

### For Review
1. Use QUICK_REFERENCE.md for quick recap
2. Refer to the common mistakes section
3. Practice the generics-vs-interfaces decision until it's automatic

---

## 🆘 Common Questions

**Q: Is `any` the same thing as `interface{}`?**
A: Yes - `any` is a built-in alias for the empty interface `interface{}`, introduced in Go 1.18 alongside generics. They're interchangeable; `any` is just easier to read.

**Q: Why can't I use `==` inside a generic function constrained with `any`?**
A: `any` only promises "some value exists" - it doesn't promise the value supports `==`. Use `comparable` instead, which specifically guarantees `==`/`!=` work.

**Q: Do I need an external package to write a custom constraint?**
A: No. `type Number interface { int | int64 | float64 }` is just a regular interface with a type set - defined entirely in your own code, no imports required.

**Q: When should I use a generic type parameter instead of an interface?**
A: When the same algorithm or data structure needs to work over arbitrary *element types* that don't need to "do" anything themselves (a stack, a cache, `Sum`/`Map`/`Filter`). Use an interface (Level 15) when different concrete types need to *behave* differently behind a shared contract.

**Q: Why does Go sometimes require me to write `Zero[int]()` instead of just `Zero()`?**
A: Go infers type arguments from the regular arguments you pass to a function. If a type parameter appears only in the return type, there's nothing to infer it from, so you must specify it explicitly.

---

## 🎯 Before Moving to Level 26

Make sure you can answer these questions:

- [ ] What problem do generics solve that neither type-specific functions nor `any` solve well?
- [ ] What's the difference between `any` and `comparable` as constraints?
- [ ] How do you define a custom constraint, and why doesn't it need an external package?
- [ ] Why must a generic type's methods repeat its type parameter (`func (s *Stack[T]) Push(...)`)?
- [ ] When does Go's type inference fail, and how do you work around it?
- [ ] When is a small interface (Level 15) the better choice over a generic type parameter?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (10 min)
   - Gets you oriented

2. **README.md** Sections 1-5 (35 min)
   - The problem, syntax, and both built-in constraints

3. **EXERCISES.md** Exercises 1-5 (2.5 hours)
   - Get hands-on with syntax, `any`, `comparable`, and custom constraints - including two real compile errors

4. **STUDY_GUIDE.md** (40 min)
   - Study the constraint hierarchy and decision tree

5. **README.md** Sections 6-9 (30 min)
   - Generic functions/types, inference, and when not to use generics

6. **EXERCISES.md** Exercises 6-10 (3+ hours)
   - Deep practice with Map/Filter/Reduce, Stack, Pair, inference, and the comprehensive Cache

7. **QUICK_REFERENCE.md** (10 min)
   - Create your mental cheat sheet

---

## 🏆 Success Indicators

You've mastered Level 25 when:

- ✅ You reach for the narrowest constraint that compiles, not the loosest one that "probably works"
- ✅ You can explain why `any` and `comparable` aren't interchangeable
- ✅ You can define a custom type-set constraint without looking it up
- ✅ You never forget to restate a generic type's type parameter on its methods
- ✅ You can predict a type-inference failure before the compiler tells you
- ✅ You default to a small interface for behavior and a type parameter for element types
- ✅ You've completed 8+ exercises, including all three real compile-error exercises
- ✅ You can explain generics to someone else

---

## 🚀 What's Next?

After Level 25, you're ready for:

**Level 26: Testing**
- Table-driven tests with the `testing` package
- `go test`, subtests, and coverage
- Testing generic functions across multiple type instantiations
- Mocking with the small-interface patterns from Level 15

---

## 💬 Key Takeaway

> **Generics give you the reusability of `any` without giving up the compiler's type checking - but they're for arbitrary element types, not arbitrary behavior. When behavior varies, reach for Level 15's small interfaces instead.**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Exercises
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Patterns
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Happy learning! Level 25 gives your programs the power to reuse algorithms and data structures across any type, safely! 🎉

*Estimated time to complete Level 25: 5-6 hours*
*Difficulty: ⭐⭐⭐⭐ (Advanced)*
*Next Level: Level 26 - Testing*
