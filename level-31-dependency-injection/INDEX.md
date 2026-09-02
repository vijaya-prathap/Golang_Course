# Level 31: Dependency Injection - INDEX

Welcome to **Level 31: Dependency Injection**! This is where Level 30's preview - wiring a repository into a service, and a service into a handler, inside `main()` - gets a name, a full explanation, and a lot more depth.

---

## 📖 What You'll Learn

- ✅ What dependency injection actually means, with no framework mystique
- ✅ Constructor injection: `func NewThing(dep SomeInterface) *Thing` - nearly all of Go DI
- ✅ Why Go usually skips DI frameworks/containers in favor of explicit `main()` wiring
- ✅ Interface-based injection: swapping implementations without touching dependent code
- ✅ Injecting configuration alongside interfaces
- ✅ The functional options pattern for optional dependencies
- ✅ Hand-wiring a full layered dependency graph in a composition root

---

## 🗂️ Level 31 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview and welcome message. Start here first!

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- What DI means, verified with a before/after refactor
- Why it matters: testability, flexibility, decoupling
- Constructor injection as the idiomatic Go pattern
- Why Go doesn't need a DI framework - `main()` as composition root
- Interface-based injection, config injection, functional options
- Best practices
- Common mistakes

**Read Time:** 45-60 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 detailed, step-by-step exercises:
1. Before/After: hardcoded dependency vs constructor injection
2. Constructor injection basics
3. Why hidden internal construction breaks testability (measured proof)
4. Interface-based injection: swapping implementations
5. Injecting a config struct alongside an interface
6. Optional dependencies via functional options
7. One service, tested with two different fakes
8. Wiring a layered dependency graph by hand in main()
9. Global mutable state vs injected dependency (a real data race)
10. Comprehensive: a notification system

**Time Commitment:** 5-6 hours
**When:** After reading README

### 4. **STUDY_GUIDE.md** (Visual Learning)
- Hardcoded vs injected dependency diagram
- Constructor injection pattern anatomy
- Interface-swap diagram
- Hand-wired dependency graph diagram
- Functional options pattern diagram
- Common mistakes guide

**Read Time:** 30-40 minutes
**When:** Use while doing exercises

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- Essential syntax reference
- Constructor injection, interface swap, config injection, functional options
- Common mistakes table

**Read Time:** 10-15 minutes
**When:** Quick lookup during practice

---

## 🎯 Recommended Learning Path

### Day 1: The Core Idea (2 hours)
1. Read **START_HERE.txt** (10 min)
2. Read **README.md** sections 1-3 (35 min)
3. Complete Exercises 1-3 (1 hour 15 min)

### Day 2: Interfaces, Config, and Options (2 hours)
1. Read **README.md** sections 4-7 (35 min)
2. Complete Exercises 4-6 (1.5 hours)

### Day 3: Testing and Wiring (2 hours)
1. Read **README.md** section 8 (10 min)
2. Complete Exercises 7-8 (1 hour 50 min)

### Day 4: Real-World Practice (1.5 hours)
1. Read **README.md** sections 9-11 (20 min)
2. Complete Exercises 9-10 (1 hour 10 min)

### Day 5: Consolidation (1 hour)
1. Try bonus challenges
2. Review with **QUICK_REFERENCE.md**
3. Make sure you can write `NewThing(dep SomeInterface)` and a composition root from memory

---

## 💡 Key Concepts At A Glance

### The Core Idea
```go
// ❌ Before - builds its own dependency
func NewThing() *Thing {
    return &Thing{dep: &ConcreteDep{}}
}

// ✅ After - receives its dependency
func NewThing(dep SomeInterface) *Thing {
    return &Thing{dep: dep}
}
```

### Interface-Based Injection
```go
proc1 := NewProcessor(StripeGateway{}) // same NewProcessor,
proc2 := NewProcessor(PayPalGateway{}) // different implementation
```

### Functional Options
```go
func NewCache(store Store, opts ...Option) *Cache { ... }

NewCache(store)                          // all defaults
NewCache(store, WithTTL(30*time.Second)) // one override
```

### Composition Root
```go
func buildApp() *Handler {
    cfg := Config{...}
    repo := NewRepository()
    service := NewService(cfg, repo)
    return NewHandler(service)
}
```

---

## ✅ Prerequisites

Make sure you've completed **Level 30: Repository-Service-Handler**

You need:
- ✅ Comfort with the repository → service → handler three-layer pattern
- ✅ Comfort with interfaces (Level 15) and the `NewX` constructor convention (Level 13)
- ✅ Comfort testing with fakes (Level 26)

---

## 🎓 Learning Objectives

By the end of Level 31, you'll be able to:

- ✅ Explain dependency injection without any framework jargon
- ✅ Refactor a hardcoded dependency into a constructor-injected one
- ✅ Write `func NewThing(dep SomeInterface) *Thing` from memory
- ✅ Explain why Go usually doesn't need a DI framework
- ✅ Swap interface implementations without touching dependent code
- ✅ Inject config values alongside interfaces
- ✅ Use the functional options pattern for optional dependencies
- ✅ Hand-wire a layered dependency graph inside a composition root

---

## 📊 Statistics

- **Main Theory:** README.md covering DI fundamentals, constructor injection, and hand-wired graphs
- **Exercises:** 10 detailed exercises + 3 bonus challenges
- **Study Guide:** diagrams for hardcoded vs injected, constructor anatomy, interface swap, graph wiring, functional options
- **Quick Reference:** One-page cheat sheet
- **Practice Time:** 5-6 hours
- **Difficulty:** ⭐⭐⭐⭐ (Advanced)

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for comprehensive theory
2. Use STUDY_GUIDE.md alongside for the diagrams
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
3. Practice writing a composition root from memory until it's automatic

---

## 🆘 Common Questions

**Q: Isn't dependency injection a framework thing?**
A: Not in Go. It's just a type receiving its dependencies as constructor arguments. No framework, container, or annotation is required - the whole idea fits in one sentence.

**Q: Why does Go's standard library not ship a DI container?**
A: Explicit constructor calls in `main()` are simple, readable, and compiler-checked. A "poor man's DI container" is just a well-organized `main()` function - most Go programs never outgrow that.

**Q: When should I inject an interface vs a plain struct?**
A: Inject an interface when the dependency has *behavior* worth swapping or faking (a database, a network client). Inject a plain struct when it's just *data* (configuration values).

**Q: What's the functional options pattern for?**
A: Adding optional, growable configuration to a constructor without forcing every call site to pass every possible setting, or rewriting every call site each time a new option is added.

**Q: Are third-party DI frameworks (wire, fx, dig) ever worth using?**
A: For very large, deep dependency graphs, yes - some teams use them. For most Go programs, including everything in this level, hand-wiring in `main()` is simpler and more readable.

---

## 🎯 Before Moving to Level 32

Make sure you can answer these questions:

- [ ] What does "dependency injection" mean, stripped of framework language?
- [ ] Why does a hardcoded dependency break testability?
- [ ] What does `func NewThing(dep SomeInterface) *Thing` buy you over building the dependency inside?
- [ ] Why does Go usually not need a DI framework?
- [ ] When should a dependency be an interface vs a plain config struct?
- [ ] What problem does the functional options pattern solve?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (10 min)
   - Gets you oriented

2. **README.md** Sections 1-4 (35 min)
   - The core idea, why it matters, constructor injection, no framework needed

3. **README.md** Sections 5-8 (30 min)
   - Interfaces, config, functional options, hand-wired graphs

4. **EXERCISES.md** Exercises 1-5 (2.5 hours)
   - Get hands-on with the before/after refactor and constructor injection

5. **STUDY_GUIDE.md** (30 min)
   - Study the diagrams for constructor anatomy and graph wiring

6. **EXERCISES.md** Exercises 6-10 (3+ hours)
   - Deep practice with options, testing with fakes, and the comprehensive notification system

7. **QUICK_REFERENCE.md** (10 min)
   - Create your mental cheat sheet

---

## 🏆 Success Indicators

You've mastered Level 31 when:

- ✅ You reach for constructor injection by default, without thinking about it
- ✅ You can explain testability, flexibility, and decoupling as DI's three payoffs
- ✅ You never let a type construct its own dependency when it could receive one instead
- ✅ You can hand-wire a layered dependency graph inside a composition root confidently
- ✅ You recognize global mutable state as a hidden dependency the moment you see it
- ✅ You've completed 8+ exercises
- ✅ You can explain dependency injection to someone else without using the word "framework"

---

## 🚀 What's Next?

After Level 31, you're ready for:

**Level 32: Authentication & JWT**
- Verifying user identity with JSON Web Tokens
- Issuing, signing, and validating tokens
- Middleware that injects an authenticated user into request handling

---

## 💬 Key Takeaway

> **Dependency injection in Go is just a type receiving what it needs as constructor arguments - `main()` is your DI container, and the compiler is your framework.**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Exercises
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Patterns
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Happy learning! Level 31 gives your programs the structure to grow, swap parts, and stay testable as they scale! 🎉

*Estimated time to complete Level 31: 5-6 hours*
*Difficulty: ⭐⭐⭐⭐ (Advanced)*
*Next Level: Level 32 - Authentication & JWT*
