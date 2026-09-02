# Level 14: Methods - INDEX

Welcome to **Level 14: Methods**! This is where your structs (Level 13) stop being passive bags of data and start carrying their own behavior.

---

## 📖 What You'll Learn

- ✅ What a method is: a function with a receiver
- ✅ Value receivers (copy) vs pointer receivers (original)
- ✅ Choosing receivers + the consistency rule
- ✅ Auto-address sugar and the non-addressable compile error
- ✅ Method sets - the foundation for Level 15 interfaces
- ✅ Methods on named non-struct types (Celsius, Path)
- ✅ Method promotion through embedding, and overriding
- ✅ Method values and method expressions

---

## 🗂️ Level 14 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview and welcome message. Start here first!

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- Methods vs plain functions; receiver conventions
- Value vs pointer receivers, proven with runnable programs
- Addressability and the auto-& sugar's limits
- Method sets (T vs *T) and why interfaces will care
- Named non-struct types, promotion, method values
- Best practices
- Common mistakes

**Read Time:** 45-60 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 detailed, step-by-step exercises:
1. Your first methods (vs plain functions)
2. Value receivers - proof of the copy
3. Pointer receivers - mutating the original
4. Auto-address and auto-dereference
5. The non-addressable compile error (and both fixes)
6. The consistency rule refactor
7. Methods on named non-struct types
8. Method promotion and overriding
9. Method values and method expressions
10. Comprehensive practice (ShoppingCart)

**Time Commitment:** 4-6 hours
**When:** After reading README

### 4. **STUDY_GUIDE.md** (Visual Learning)
- Method anatomy diagram
- Value vs pointer receiver comparison table + decision tree
- Method-set table (T vs *T)
- Addressability map and promotion diagram
- Common mistakes guide

**Read Time:** 30-40 minutes
**When:** Use while doing exercises

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- Method declaration syntax
- Receiver decision rules
- Method-set table
- Common mistakes table

**Read Time:** 10-15 minutes
**When:** Quick lookup during practice

---

## 🎯 Recommended Learning Path

### Day 1: Methods & Value Receivers (2 hours)
1. Read **START_HERE.txt** (10 min)
2. Read **README.md** sections 1-3 (30 min)
3. Complete Exercises 1-2 (1 hour 20 min)

### Day 2: Pointer Receivers & Choosing (2 hours)
1. Read **README.md** sections 4-5 (25 min)
2. Complete Exercises 3-4 (1.5 hours)

### Day 3: Addressability & Method Sets (2 hours)
1. Read **README.md** sections 6-7 (30 min)
2. Complete Exercises 5-6 (1.5 hours)

### Day 4: Named Types, Promotion & Method Values (2 hours)
1. Read **README.md** sections 8-10 (30 min)
2. Complete Exercises 7-9 (1.5 hours)

### Day 5: Consolidation (1.5 hours)
1. Complete Exercise 10 (ShoppingCart)
2. Try bonus challenges
3. Review with **QUICK_REFERENCE.md** - make sure the method-set table is automatic

---

## 💡 Key Concepts At A Glance

### A Method Is a Function With a Receiver
```go
func (r Rectangle) Area() float64 { return r.Width * r.Height }
rect.Area() // 12
```

### Value vs Pointer Receivers
```go
func (c Counter) IncrementCopy()  { c.Count++ } // mutation LOST
func (c *Counter) Increment()     { c.Count++ } // mutation STICKS
```

### The Sugar and Its Limit
```go
acct.Deposit(50)              // OK: sugar for (&acct).Deposit(50)
players["alice"].AddPoints(5) // ❌ compile error - not addressable
```

### Method Sets
```go
// T:  value-receiver methods only
// *T: value-receiver AND pointer-receiver methods
```

### Promotion
```go
type Dog struct{ Animal }
d.Describe()       // Animal's method, promoted to Dog
```

---

## ✅ Prerequisites

Make sure you've completed **Level 13: Structs**

You need:
- ✅ Comfort defining structs and struct literals
- ✅ Understanding of embedding and field promotion
- ✅ Constructors (NewXxx functions) from Level 13
- ✅ Pointers and & / * from Level 12

---

## 🎓 Learning Objectives

By the end of Level 14, you'll be able to:

- ✅ Declare methods with value and pointer receivers
- ✅ Predict whether a mutation sticks before running the code
- ✅ Apply the consistency rule when designing a type's methods
- ✅ Explain the non-addressable compile error and fix it two ways
- ✅ State the method-set rule for T vs *T from memory
- ✅ Attach methods to named non-struct types
- ✅ Use method promotion and overriding through embedding
- ✅ Bind method values and use method expressions

---

## 📊 Statistics

- **Main Theory:** README.md covering receivers, method sets, promotion, and method values
- **Exercises:** 10 detailed exercises + 3 bonus challenges
- **Study Guide:** Receiver decision tree, method-set table, promotion diagram
- **Quick Reference:** One-page cheat sheet
- **Practice Time:** 4-6 hours
- **Difficulty:** ⭐⭐⭐ (Intermediate)

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for comprehensive theory
2. Use STUDY_GUIDE.md alongside for the receiver decision tree and method-set table
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
3. Practice the method-set table until it's automatic

---

## 🆘 Common Questions

**Q: Are methods just functions with different syntax?**
A: Essentially yes - same power, different organization. But only methods count toward a type's method set, which is what satisfies interfaces in Level 15. That's why methods matter.

**Q: Why did my method's change to the struct disappear?**
A: Your receiver is a value receiver, so the method mutated a copy. Switch to a pointer receiver: `func (c *Counter) Increment()`.

**Q: Why does `players["alice"].AddPoints(5)` fail to compile?**
A: `AddPoints` needs a `*Player`, and map elements aren't addressable - Go can't take `&players["alice"]`. Copy the value out, modify it, store it back - or store `*Player` in the map.

**Q: If Go auto-takes the address anyway, why do method sets exclude pointer methods from T?**
A: The sugar is a compile-time rewrite that needs an addressable variable at the call site. An interface value stores a copy with no stable address, so there's nothing to auto-address - hence the stricter method-set rule.

**Q: Should I name my receiver this or self?**
A: No - Go convention is a short name derived from the type, like `c` for `*Counter`. Using this/self marks code as non-idiomatic.

**Q: Is method promotion the same as inheritance?**
A: No - it's composition with convenient syntax. `d.Describe()` is shorthand for `d.Animal.Describe()`; there's no virtual dispatch, and an outer same-name method simply shadows the promoted one.

---

## 🎯 Before Moving to Level 15

Make sure you can answer these questions:

- [ ] What exactly does a value receiver receive, and what happens to mutations?
- [ ] What are the three reasons to choose a pointer receiver?
- [ ] Why does calling a pointer method on a map element fail to compile?
- [ ] Which methods are in the method set of T? Of *T?
- [ ] How do you override a promoted method, and how do you still reach the embedded version?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (10 min)
   - Gets you oriented

2. **README.md** Sections 1-5 (35 min)
   - Methods, receivers, and choosing between them

3. **README.md** Sections 6-7 (25 min)
   - Addressability and method sets - the tricky heart of the level

4. **EXERCISES.md** Exercises 1-5 (2.5 hours)
   - Get hands-on, including the intentional compile error

5. **STUDY_GUIDE.md** (30 min)
   - Study the decision tree and method-set table

6. **EXERCISES.md** Exercises 6-10 (2.5 hours)
   - Refactoring, named types, promotion, and the ShoppingCart

7. **QUICK_REFERENCE.md** (10 min)
   - Create your mental cheat sheet

---

## 🏆 Success Indicators

You've mastered Level 14 when:

- ✅ You choose value or pointer receivers deliberately, not by habit
- ✅ You never mix receivers on one type without a reason you can defend
- ✅ You can predict the non-addressable compile error before the compiler tells you
- ✅ You can recite the method-set rule and explain why interfaces will need it
- ✅ You can use promotion and overriding without confusing them with inheritance
- ✅ You've completed 8+ exercises
- ✅ You can explain methods to someone else

---

## 🚀 What's Next?

After Level 14, you're ready for:

**Level 15: Interfaces**
- Defining interfaces and implicit satisfaction
- Method sets in action: which types satisfy which interfaces
- The empty interface and type assertions
- fmt.Stringer and the standard library's small-interface style

---

## 💬 Key Takeaway

> **A method is just a function with a receiver - but the receiver's kind (value or pointer) decides whether mutations stick, and the method set it builds decides which interfaces your type will satisfy in Level 15.**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Exercises
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Patterns
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Happy learning! Level 14 turns your data into types with real behavior! 🎉

*Estimated time to complete Level 14: 4-6 hours*
*Difficulty: ⭐⭐⭐ (Intermediate)*
*Next Level: Level 15 - Interfaces*
