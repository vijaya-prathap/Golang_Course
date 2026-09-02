# Level 30: Repository-Service-Handler - INDEX

Welcome to **Level 30: Repository-Service-Handler**! This is where your programs stop being one tangled file and start being three layers, each with exactly one job - and each independently testable.

---

## 📖 What You'll Learn

- ✅ What the Repository-Service-Handler pattern is, and why each layer has one job
- ✅ The Repository layer: a small interface + an in-memory implementation (revisiting Level 15)
- ✅ The Service layer: business rules that depend on the Repository interface, never the concrete type
- ✅ The Handler layer: a thin `net/http` translation layer, tested with `httptest` (Level 27/26)
- ✅ Dependency direction: Handler → Service → Repository, interfaces flowing the other way
- ✅ Wiring all three layers together by hand in `main()` - a Level 31 preview
- ✅ Why this makes testing easy: fakes at every layer boundary
- ✅ Error handling across layers: repository errors → service errors → HTTP status codes
- ✅ When three layers are overkill, and when they earn their keep

---

## 🗂️ Level 30 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview and welcome message. Start here first!

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- What the pattern is and why it separates concerns
- The Repository layer: interface + in-memory implementation
- The Service layer: validation and uniqueness rules the repository can't enforce
- The Handler layer: thin HTTP translation, verified with httptest
- Dependency direction and interface ownership
- Wiring everything together in main()
- Why this enables easy testing, with a fake repository proof
- Error handling across layers, with a full 400/404/409/500 translation table
- When this pattern is overkill
- Best practices
- Common mistakes

**Read Time:** 50-65 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 detailed, step-by-step exercises:
1. The Repository layer - interface + in-memory implementation
2. The Service layer - validation before storage
3. The Service layer - a rule that needs the repository (uniqueness)
4. The Handler layer - translating HTTP, tested with httptest
5. Wiring all three layers together, end-to-end
6. Testing the Service in isolation with a fake repository
7. Mapping a not-found repository error to a 404 response
8. Dependency direction - concrete type vs interface
9. A full CRUD HTTP API, table-driven and lifecycle-tested
10. Comprehensive practice - a complete Notes API across real packages

Plus 3 bonus challenges: a second repository implementation, a caching decorator, and service-level pagination.

**Time Commitment:** 6-7 hours
**When:** After reading README

### 4. **STUDY_GUIDE.md** (Visual Learning)
- The three-layer architecture diagram with dependency arrows
- Interface-ownership diagram (defined where it's used)
- Testing-with-fakes diagram at every layer boundary
- Error-translation-across-layers diagram
- Common mistakes guide
- Prerequisites/readiness/checklist for Level 31

**Read Time:** 30-40 minutes
**When:** Use while doing exercises

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- Essential syntax reference for all three layers
- Error translation table
- Multi-package layout reference
- Common mistakes table

**Read Time:** 10-15 minutes
**When:** Quick lookup during practice

---

## 🎯 Recommended Learning Path

### Day 1: The Repository Layer (2 hours)
1. Read **START_HERE.txt** (10 min)
2. Read **README.md** sections 1-2 (25 min)
3. Complete Exercise 1 (1 hour 25 min)

### Day 2: The Service Layer (2 hours)
1. Read **README.md** section 3 (20 min)
2. Complete Exercises 2-3 (1.5 hours)

### Day 3: The Handler Layer & Wiring (2 hours)
1. Read **README.md** sections 4-6 (30 min)
2. Complete Exercises 4-5 (1.5 hours)

### Day 4: Testing & Error Translation (2 hours)
1. Read **README.md** sections 7-8 (25 min)
2. Complete Exercises 6-7 (1.5 hours)

### Day 5: Dependency Direction & Full CRUD (2 hours)
1. Read **README.md** section 9 (10 min)
2. Complete Exercises 8-9 (1.75 hours)

### Day 6: The Comprehensive Exercise (2 hours)
1. Read **README.md** sections 10-11 (15 min)
2. Complete Exercise 10 (1.75 hours)

### Day 7: Consolidation (1 hour)
1. Try the bonus challenges
2. Review with **QUICK_REFERENCE.md**
3. Make sure the dependency direction diagram feels automatic

---

## 💡 Key Concepts At A Glance

### The Three Layers
```go
type UserRepository interface { /* data access */ }
type UserService struct { repo UserRepository /* business rules */ }
type UserHandler struct { service *UserService /* HTTP translation */ }
```

### Dependency Direction
```
Handler --> Service --> Repository (interface)
                              ▲
                implemented by InMemoryUserRepository
```

### Testing With Fakes
```go
service := NewUserService(alwaysErrorRepository{}) // no real storage needed
```

### Error Translation
```go
errors.Is(err, ErrNotFound) -> http.StatusNotFound
```

---

## ✅ Prerequisites

Make sure you've completed **Level 29: Database/SQL**

You need:
- ✅ Comfort with interfaces and implicit satisfaction (Level 15)
- ✅ Comfort with error wrapping and `errors.Is` (Level 16)
- ✅ Comfort with table-driven tests and mocking via interfaces (Level 26)
- ✅ Comfort with `net/http` handlers, routing, and `httptest` (Level 27)
- ✅ Familiarity with the idea of a repository/storage layer (Level 29) - this level deliberately replaces it with an in-memory store so every exercise stays offline

---

## 🎓 Learning Objectives

By the end of Level 30, you'll be able to:

- ✅ Define a small Repository interface and satisfy it implicitly with an in-memory implementation
- ✅ Write Service logic that enforces business rules the repository has no concept of
- ✅ Prove a Service rejects invalid input before ever touching the repository
- ✅ Write thin Handlers that only translate HTTP, tested entirely with `httptest`
- ✅ Explain and defend the Handler → Service → Repository dependency direction
- ✅ Wire all three layers together by hand in `main()`
- ✅ Swap in fake repositories (and fake services) to test any layer in isolation
- ✅ Translate errors across layer boundaries into the correct HTTP status codes
- ✅ Recognize when this pattern is worth its ceremony, and when it isn't

---

## 📊 Statistics

- **Main Theory:** README.md covering the three-layer pattern, dependency direction, wiring, testing, and error translation
- **Exercises:** 10 detailed exercises + 3 bonus challenges (second repository, caching decorator, pagination)
- **Study Guide:** architecture diagram, interface-ownership diagram, testing-with-fakes diagram, error-translation diagram
- **Quick Reference:** One-page cheat sheet
- **Practice Time:** 6-7 hours
- **Difficulty:** ⭐⭐⭐⭐ (Advanced)

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for comprehensive theory
2. Use STUDY_GUIDE.md alongside for the architecture and error-translation diagrams
3. Reference QUICK_REFERENCE.md for fast lookup

### For Practice
1. Read exercise description
2. Copy provided bash commands (or type them) - some exercises create multiple files, one creates a real multi-package project
3. Run `go run` or `go test -v` as shown and compare against the captured Expected Output
4. Understand what each step does, especially where a fake repository proves a claim about behavior

### For Review
1. Use QUICK_REFERENCE.md for quick recap
2. Refer to the common mistakes section
3. Practice explaining the dependency direction diagram out loud until it's automatic

---

## 🆘 Common Questions

**Q: Why in-memory storage instead of the SQL setup from Level 29?**
A: This level is about the *architecture*, not the storage technology. In-memory storage keeps every exercise 100% offline and fast, and proves the point of the pattern: the Repository interface could be backed by SQL, memory, a file, or anything else, and the Service and Handler would never know the difference.

**Q: Isn't this just adding unnecessary layers to a simple CRUD app?**
A: For a genuinely tiny program, yes - see "When This Is Overkill" in the README. The pattern earns its keep once business logic, testability, or swappable storage actually matter.

**Q: Where should the Repository interface live - the repository package or the service package?**
A: The service package (or wherever it's consumed), per Level 15's Rule 3 ("define interfaces where they're used"). The repository package just exports concrete types with the right methods.

**Q: How is "wiring in main()" different from what Level 31 teaches?**
A: It isn't fundamentally different - it's the same idea, done by hand. Level 31 gives it a name (dependency injection), shows how to manage it once the dependency graph grows past three layers, and discusses where frameworks or generators fit into Go (rarely, compared to other languages).

**Q: Why does the Handler translate errors instead of the Service returning HTTP status codes directly?**
A: Because the Service shouldn't know HTTP exists - it might be called from a CLI tool, a background job, or a gRPC handler tomorrow. Only the Handler, which already speaks HTTP, should know what a 404 is.

---

## 🎯 Before Moving to Level 31

Make sure you can answer these questions:

- [ ] What is the ONE job each of Repository, Service, and Handler has?
- [ ] Which direction do dependencies flow, and which direction does the interface satisfaction flow?
- [ ] Why does the Repository interface belong in the Service's package, not the Repository's?
- [ ] How do you prove a Service rejects invalid input before touching the repository?
- [ ] How does a "not found" error become a 404, across two layer boundaries?
- [ ] When is this pattern overkill?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (10 min)
   - Gets you oriented

2. **README.md** Sections 1-4 (35 min)
   - The pattern, Repository, Service, Handler

3. **EXERCISES.md** Exercises 1-4 (3 hours)
   - Get hands-on with each layer individually

4. **README.md** Sections 5-8 (30 min)
   - Dependency direction, wiring, testing, error translation

5. **STUDY_GUIDE.md** (30 min)
   - Study the architecture and error-translation diagrams

6. **EXERCISES.md** Exercises 5-9 (3 hours)
   - Wiring, fakes, error mapping, full CRUD

7. **README.md** Sections 9-11 (20 min)
   - When it's overkill, best practices, common mistakes

8. **EXERCISES.md** Exercise 10 + Bonus (2+ hours)
   - The comprehensive multi-package Notes API and bonus challenges

9. **QUICK_REFERENCE.md** (10 min)
   - Create your mental cheat sheet

---

## 🏆 Success Indicators

You've mastered Level 30 when:

- ✅ You instinctively put validation in the Service, never the Handler or Repository
- ✅ You can explain why a Service depending on a concrete repository type breaks the pattern
- ✅ You reach for a fake repository/service to test a layer in isolation without hesitation
- ✅ You can trace a "not found" error from the Repository all the way to a 404 response
- ✅ You know where to put an interface: next to the code that calls it, not the code that implements it
- ✅ You've completed 8+ exercises
- ✅ You can explain the whole pattern to someone else with a simple diagram

---

## 🚀 What's Next?

After Level 30, you're ready for:

**Level 31: Dependency Injection**
- Naming and scaling the "construct dependencies and hand them down" pattern from this level
- Managing a dependency graph as it grows past three layers
- Constructor injection, wiring functions, and where DI frameworks fit into Go

---

## 💬 Key Takeaway

> **Repository, Service, and Handler each get exactly one job - and because every layer depends on an interface owned by its consumer, every layer can be tested and changed without touching the others.**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Exercises
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Patterns
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Happy learning! Level 30 gives your programs the structure real production Go services are built on! 🎉

*Estimated time to complete Level 30: 6-7 hours*
*Difficulty: ⭐⭐⭐⭐ (Advanced)*
*Next Level: Level 31 - Dependency Injection*
