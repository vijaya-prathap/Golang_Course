# Level 13: Structs - INDEX

Welcome to **Level 13: Structs**! This is where your programs start grouping related data into named types instead of juggling separate variables - Go's answer to "objects" (the data half, at least).

---

## 📖 What You'll Learn

- ✅ Defining structs, field types, and zero-value structs
- ✅ Struct literals: keyed (preferred) vs positional
- ✅ Struct value semantics: copy on assignment vs pointer sharing
- ✅ Comparing structs with == (and when it won't compile)
- ✅ Nested structs and embedded structs with field promotion
- ✅ Anonymous structs, struct tags, and the NewX constructor convention

---

## 🗂️ Level 13 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview and welcome message. Start here first!

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- What a struct is and how to define one
- Zero values, keyed vs positional literals
- Value semantics: copying vs pointer-based sharing
- Struct comparison rules and the uncomparable-field compile error
- Nested structs, embedding, field promotion, and shadowing
- Anonymous structs, struct tags, and NewX constructors
- Best practices and common mistakes

**Read Time:** 45-60 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 detailed, step-by-step exercises:
1. Defining and initializing structs (all literal forms)
2. Exploring the zero-value struct
3. Structs are value types (copy vs pointer)
4. Comparing structs (and when you can't - captured compile error)
5. Nested structs
6. Embedded structs and field promotion
7. Anonymous structs
8. Struct tags
9. The constructor convention (NewX)
10. Comprehensive practice - library inventory

Plus 3 bonus challenges (2D vector type, config struct with defaults, company org model).

**Time Commitment:** 4-5 hours
**When:** After reading README

### 4. **STUDY_GUIDE.md** (Visual Learning)
- Struct memory layout diagram
- Literal forms comparison table
- Value-copy vs pointer-sharing diagram
- Embedding/promotion and shadowing diagrams
- Comparable vs uncomparable field types table
- Common mistakes guide

**Read Time:** 30-40 minutes
**When:** Use while doing exercises

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- Essential syntax reference
- Comparison and embedding syntax at a glance
- Common mistakes table

**Read Time:** 10-15 minutes
**When:** Quick lookup during practice

---

## 🎯 Recommended Learning Path

### Day 1: Defining Structs & Zero Values (2 hours)
1. Read **START_HERE.txt** (10 min)
2. Read **README.md** sections 1-2 (25 min)
3. Complete Exercises 1-2 (1 hour 25 min)

### Day 2: Value Semantics & Comparison (2 hours)
1. Read **README.md** sections 3-5 (30 min)
2. Complete Exercises 3-4 (1.5 hours)

### Day 3: Nesting & Embedding (2 hours)
1. Read **README.md** sections 6-7 (25 min)
2. Complete Exercises 5-6 (1 hour 35 min)

### Day 4: Anonymous Structs, Tags & Constructors (1.5 hours)
1. Read **README.md** sections 8-10 (20 min)
2. Complete Exercises 7-9 (1 hour 10 min)

### Day 5: Comprehensive Practice & Consolidation (1.5 hours)
1. Complete Exercise 10 (45 min)
2. Try bonus challenges
3. Review with **QUICK_REFERENCE.md**

---

## 💡 Key Concepts At A Glance

### Defining and Zero Values
```go
type Person struct {
    Name string
    Age  int
}
var p Person // {Name:"" Age:0} - every field zero, immediately usable
```

### Keyed vs Positional Literals
```go
alice := Person{Name: "Alice", Age: 30} // ✅ preferred - order-independent
bob := Person{"Bob", 25}                 // ❌ fragile - avoid
```

### Value Semantics
```go
b := a       // COPIES a's fields - independent
ptr := &a
ptr.Field = x // shares a - mutates through the pointer
```

### Embedding and Promotion
```go
type Dog struct {
    Animal // embedded - fields promoted
    Breed string
}
d.Name // same as d.Animal.Name
```

---

## ✅ Prerequisites

Make sure you've completed **Level 12: Pointers**

You need:
- ✅ Comfort with `&` and `*`, and nil pointers
- ✅ Understand value vs pointer parameters
- ✅ Comfortable creating values with `&T{}`

---

## 🎓 Learning Objectives

By the end of Level 13, you'll be able to:

- ✅ Define struct types and create them with every literal form
- ✅ Predict every field's zero value in a zero-value struct
- ✅ Prove copy-on-assignment and choose pointers when sharing is needed
- ✅ Explain exactly when struct comparison compiles and when it doesn't
- ✅ Build nested structs and traverse deep fields
- ✅ Use embedding, field promotion, and explain shadowing
- ✅ Reach for anonymous structs and recognize struct tag syntax
- ✅ Write NewX constructor functions

---

## 📊 Statistics

- **Main Theory:** README.md covering struct definition, semantics, composition, and conventions
- **Exercises:** 10 detailed exercises + 3 bonus challenges
- **Study Guide:** memory layout diagrams, literal/comparison tables, embedding diagrams
- **Quick Reference:** One-page cheat sheet
- **Practice Time:** 4-5 hours
- **Difficulty:** ⭐⭐⭐ (Intermediate-Advanced)

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for comprehensive theory
2. Use STUDY_GUIDE.md alongside for memory-layout and promotion diagrams
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
3. Practice explaining copy-vs-pointer and promotion-vs-shadowing until they're automatic

---

## 🆘 Common Questions

**Q: Are Go structs the same as classes?**
A: No. A struct is only the data half. There are no constructors, no inheritance, and no access modifiers built into the type. Level 14 (Methods) adds behavior, and Level 15 (Interfaces) adds polymorphism - together they're Go's composition-based alternative to classes.

**Q: Why does assigning a struct copy it instead of sharing it?**
A: A struct is a value type, exactly like an array (Level 7). Its fields live inline in one block of memory - assignment copies that whole block. Use a pointer (`*T`) when you want changes to be visible elsewhere.

**Q: Why can't I compare two structs that have a slice field?**
A: Slices (and maps and functions) don't support `==` - Go can't define what "equal" means for them cheaply. Since struct comparison needs every field to be comparable, one uncomparable field disables `==` for the whole struct.

**Q: What's the difference between a nested struct and an embedded struct?**
A: A nested struct has a field name (`Address Address`) and you always use the full path (`p.Address.City`). An embedded struct has no field name - just the type - and its fields are *promoted* so you can access them directly on the outer struct (`d.Name` instead of `d.Animal.Name`).

**Q: Do I need struct tags right now?**
A: Just recognize the syntax (`` `json:"name"` ``) and know libraries read them via reflection. You won't use them for real until Level 19 (JSON).

---

## 🎯 Before Moving to Level 14

Make sure you can answer these questions:

- [ ] What's the difference between a keyed and a positional struct literal, and why prefer keyed?
- [ ] What is every field's value in a zero-value struct you just declared with `var`?
- [ ] Does assigning a struct to a new variable copy it or share it? What about a pointer?
- [ ] When does `==` fail to compile on a struct, and why?
- [ ] What does field promotion mean, and what happens when a field name is shadowed?
- [ ] Why does Go use a NewX function instead of a constructor?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (10 min)
   - Gets you oriented

2. **README.md** Sections 1-5 (30 min)
   - Defining structs, literals, value semantics, comparison

3. **README.md** Sections 6-10 (30 min)
   - Nesting, embedding, anonymous structs, tags, constructors

4. **EXERCISES.md** Exercises 1-5 (2 hours)
   - Get hands-on with literals, zero values, copy semantics, and comparison

5. **STUDY_GUIDE.md** (30 min)
   - Study the memory layout and embedding/promotion diagrams

6. **EXERCISES.md** Exercises 6-10 (2+ hours)
   - Deep practice with embedding, anonymous structs, tags, and a full mini-model

7. **QUICK_REFERENCE.md** (10 min)
   - Create your mental cheat sheet

---

## 🏆 Success Indicators

You've mastered Level 13 when:

- ✅ You reach for keyed literals by default and never write positional ones for real structs
- ✅ You can explain and demonstrate copy vs pointer sharing without hesitation
- ✅ You never assume `==` works on a struct without checking its field types
- ✅ You can build nested and embedded structs and explain promotion and shadowing
- ✅ You recognize when a NewX constructor is worth writing, and when it isn't
- ✅ You've completed 8+ exercises
- ✅ You can explain structs to someone else

---

## 🚀 What's Next?

After Level 13, you're ready for:

**Level 14: Methods**
- Attaching behavior to your struct types
- Value receivers vs pointer receivers
- Method promotion through embedding
- Method sets - the doorway to interfaces (Level 15)

---

## 💬 Key Takeaway

> **A struct groups fields into one named value that copies as a whole - embedding lets one struct borrow another's fields (and soon, methods) without inheritance.**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Exercises
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Patterns
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Happy learning! Level 13 gives your programs the power to model real things as structured data! 🎉

*Estimated time to complete Level 13: 4-5 hours*
*Difficulty: ⭐⭐⭐ (Intermediate-Advanced)*
*Next Level: Level 14 - Methods*
