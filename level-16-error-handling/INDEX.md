# Level 16: Error Handling - INDEX

Welcome to **Level 16: Error Handling**! This is where the `if err != nil` pattern you've been typing since Level 3 finally becomes the explicit subject - and where Level 15's interfaces pay off in the most common interface Go code ever touches.

---

## 📖 What You'll Learn

- ✅ The `error` interface and why Go returns errors instead of throwing exceptions
- ✅ Creating errors with `errors.New` and `fmt.Errorf`
- ✅ The return-early flow and adding context at every layer
- ✅ Wrapping errors with `%w` and walking chains with `errors.Unwrap`
- ✅ Matching sentinel errors anywhere in a chain with `errors.Is`
- ✅ Extracting custom error types and their fields with `errors.As`
- ✅ Designing custom error types with structured fields
- ✅ The typed-nil error trap - and why you must return `error`, never a concrete type
- ✅ `panic` and `recover`, and when panicking is (rarely) the right call
- ✅ `defer` for cleanup alongside error returns

---

## 🗂️ Level 16 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview and welcome message. Start here first!

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- The error interface, errors.New, fmt.Errorf
- Return-early flow and layered context
- Wrapping with %w, errors.Unwrap, errors.Is, errors.As
- Custom error types and the typed-nil trap
- panic, recover, and defer for cleanup
- Best practices and common mistakes

**Read Time:** 50-65 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 detailed, step-by-step exercises:
1. Errors are values - the if err != nil flow
2. Adding context through two layers
3. Wrapping with %w and errors.Unwrap
4. errors.Is - a sentinel through a wrapped chain
5. errors.As - extracting a custom error type
6. A custom ValidationError with structured fields
7. The typed-nil error trap (and the fix)
8. panic and recover
9. Handle errors exactly once (refactoring)
10. Comprehensive practice - user registration validator

Plus 3 bonus challenges: a retry helper, a multi-error collector with errors.Join, and implementing your own Unwrap().

**Time Commitment:** 5-6 hours
**When:** After reading README

### 4. **STUDY_GUIDE.md** (Visual Learning)
- The error interface diagram
- Error-chain diagrams showing Is/As traversal
- Sentinel vs custom-type vs ad-hoc decision tree
- panic-vs-error decision tree
- defer/recover flow diagram
- Common mistakes guide

**Read Time:** 35-45 minutes
**When:** Use while doing exercises

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- Essential syntax reference
- errors.Is / errors.As patterns
- The typed-nil trap, at a glance
- Common mistakes table

**Read Time:** 10-15 minutes
**When:** Quick lookup during practice

---

## 🎯 Recommended Learning Path

### Day 1: The Foundation (2 hours)
1. Read **START_HERE.txt** (10 min)
2. Read **README.md** sections 1-3 (35 min)
3. Complete Exercises 1-2 (1 hour 15 min)

### Day 2: Wrapping and Inspecting Chains (2 hours)
1. Read **README.md** sections 4-6 (30 min)
2. Complete Exercises 3-5 (1.5 hours)

### Day 3: Custom Types and the Typed-Nil Trap (1.5 hours)
1. Read **README.md** sections 7-8 (20 min)
2. Complete Exercises 6-7 (1 hour 10 min)

### Day 4: panic, recover, and Cleanup (1.5 hours)
1. Read **README.md** sections 9-10 (20 min)
2. Complete Exercise 8 (1 hour 10 min)

### Day 5: Best Practices in Practice (1.5 hours)
1. Read **README.md** sections 11-12 (15 min)
2. Complete Exercises 9-10 (1 hour 15 min)

### Day 6: Consolidation (1 hour)
1. Try bonus challenges
2. Review with **QUICK_REFERENCE.md**
3. Make sure errors.Is vs errors.As feels automatic

---

## 💡 Key Concepts At A Glance

### The error Interface
```go
type error interface {
    Error() string
}
```

### The Fundamental Flow
```go
result, err := doSomething()
if err != nil {
    return fmt.Errorf("doing something: %w", err)
}
```

### Sentinels and Custom Types
```go
var ErrNotFound = errors.New("not found")   // sentinel + errors.Is
type ValidationError struct{ Field, Reason string }  // custom type + errors.As
```

### The Typed-Nil Trap
```go
func f() *MyError { return nil }  // ❌ non-nil error interface
func f() error    { return nil }  // ✅ genuinely nil error
```

### panic and recover
```go
defer func() {
    if r := recover(); r != nil {
        err = fmt.Errorf("recovered: %v", r)
    }
}()
```

---

## ✅ Prerequisites

Make sure you've completed **Level 15: Interfaces**

You need:
- ✅ Comfort with implicit interface satisfaction
- ✅ Comfort with type assertions and type switches
- ✅ Understanding of the typed-nil interface gotcha (this level puts it to direct use)
- ✅ Comfort with pointer receivers vs value receivers (Level 14)
- ✅ Comfort with `defer` and its LIFO ordering (Level 11)

---

## 🎓 Learning Objectives

By the end of Level 16, you'll be able to:

- ✅ Explain the error interface and why Go treats errors as values
- ✅ Create errors with errors.New and fmt.Errorf, following message style conventions
- ✅ Return early and add context at every layer without breaking the chain
- ✅ Wrap errors with %w and inspect chains with errors.Is and errors.As
- ✅ Design sentinel errors and custom error types, and know when to use each
- ✅ Reproduce and avoid the typed-nil error trap
- ✅ Use panic and recover appropriately, and explain why panicking is rare
- ✅ Pair defer with error returns for reliable cleanup

---

## 📊 Statistics

- **Main Theory:** README.md covering the error interface, wrapping, Is/As, custom types, and panic/recover
- **Exercises:** 10 detailed exercises + 3 bonus challenges
- **Study Guide:** error-chain diagrams, decision trees, panic/recover flow
- **Quick Reference:** One-page cheat sheet
- **Practice Time:** 5-6 hours
- **Difficulty:** ⭐⭐⭐ (Intermediate-Advanced)

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for comprehensive theory
2. Use STUDY_GUIDE.md alongside for the chain diagrams and decision trees
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
3. Practice errors.Is vs errors.As until choosing between them is automatic

---

## 🆘 Common Questions

**Q: Why doesn't Go have try/catch?**
A: Errors are ordinary return values, not a separate control-flow mechanism. A function's signature tells you it can fail, and `if err != nil` handles it right where the call happens - no hidden jumps to a distant catch block.

**Q: When should I use errors.Is vs errors.As?**
A: `errors.Is` answers "is this specific sentinel value anywhere in the chain?" `errors.As` answers "is there an error of this type anywhere in the chain, and can I have it to read its fields?" Use Is for identity checks, As when you need data out of the error.

**Q: Why does err != nil sometimes lie?**
A: The typed-nil trap. If a function's signature returns a concrete pointer type (like `*MyError`) instead of `error`, a nil `*MyError` becomes a non-nil `error` interface value once assigned - because the interface's type half is non-nil even though the value half is nil. Always declare `error` as the return type.

**Q: When is it okay to panic?**
A: Almost never in application code. Reserve it for programmer bugs and impossible states - not for expected failures like bad input or a missing file. Library code may panic internally but should recover at its public boundary and return an error instead.

**Q: What does %w actually do that %v doesn't?**
A: Both format the same way. `%w` additionally records the wrapped error so `errors.Unwrap`, `errors.Is`, and `errors.As` can find it later. `%v` only copies the text - the original error's identity and type are lost.

---

## 🎯 Before Moving to Level 17

Make sure you can answer these questions:

- [ ] What is the error interface, and why does implementing Error() string make a type an error?
- [ ] What's the difference between errors.Is and errors.As, and when do you use each?
- [ ] Why does the typed-nil error trap happen, and how do you avoid it?
- [ ] How does %w differ from %v inside fmt.Errorf?
- [ ] When is panic appropriate, and how does recover work?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (10 min)
   - Gets you oriented

2. **README.md** Sections 1-4 (30 min)
   - The error interface, creating errors, the standard flow, wrapping

3. **README.md** Sections 5-8 (35 min)
   - errors.Is, errors.As, custom types, the typed-nil trap

4. **EXERCISES.md** Exercises 1-6 (2.5 hours)
   - Get hands-on with values, context, wrapping, Is, As, and custom types

5. **STUDY_GUIDE.md** (35 min)
   - Study the chain diagrams and decision trees

6. **README.md** Sections 9-12 (25 min)
   - panic, recover, defer, best practices, common mistakes

7. **EXERCISES.md** Exercises 7-10 (2.5 hours)
   - Deep practice with the typed-nil trap, panic/recover, refactoring, and the comprehensive validator

8. **QUICK_REFERENCE.md** (10 min)
   - Create your mental cheat sheet

---

## 🏆 Success Indicators

You've mastered Level 16 when:

- ✅ You reach for `if err != nil` and return-with-context automatically
- ✅ You never compare a possibly-wrapped error with `==`
- ✅ You can explain the typed-nil error trap to someone else without hesitation
- ✅ You know instantly whether a new error should be a sentinel, a custom type, or ad-hoc
- ✅ You reserve panic for bugs, not for expected failures
- ✅ You've completed 8+ exercises
- ✅ You can explain error handling to someone else

---

## 🚀 What's Next?

After Level 16, you're ready for:

**Level 17: Packages & Modules**
- Organizing code into multiple packages
- Exported vs unexported identifiers at package scale
- go.mod, module paths, and third-party dependencies
- Designing clean package APIs (including which errors to export)

---

## 💬 Key Takeaway

> **An error is just a value with an Error() string method - handle it once, add context on the way up, and reach for panic only when something has truly gone impossible.**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Exercises
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Patterns
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Happy learning! Level 16 gives your programs the discipline to fail gracefully and predictably! 🎉

*Estimated time to complete Level 16: 5-6 hours*
*Difficulty: ⭐⭐⭐ (Intermediate-Advanced)*
*Next Level: Level 17 - Packages & Modules*
