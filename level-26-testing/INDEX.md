# Level 26: Testing - INDEX

Welcome to **Level 26: Testing**! This is where your programs stop being "code that seems to work" and become code you can *prove* works - automatically, every time it changes.

---

## 📖 What You'll Learn

- ✅ The `testing` package, `_test.go` naming, and `go test`
- ✅ t.Error/t.Errorf vs t.Fatal/t.Fatalf
- ✅ Table-driven tests - the idiomatic Go testing pattern
- ✅ Subtests with t.Run, and running one with -run
- ✅ Test coverage with go test -cover and go tool cover
- ✅ Benchmarks with BenchmarkXxx and go test -bench
- ✅ Test helpers and t.Helper()
- ✅ Setup/teardown with TestMain and t.Cleanup()
- ✅ Testing HTTP handlers with httptest (a Level 27 preview)
- ✅ Mocking dependencies via interfaces (revisiting Level 15)

---

## 🗂️ Level 26 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview and welcome message. Start here first!

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- The testing package, naming rules, and go test
- t.Error vs t.Fatal, verified with real output
- Table-driven tests catching a real seeded bug
- Subtests, coverage, benchmarks
- Test helpers, setup/teardown, httptest, mocking via interfaces
- Best practices
- Common mistakes

**Read Time:** 50-65 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 detailed, step-by-step exercises (each creates a code file AND a test file):
1. A basic test function
2. t.Error vs t.Fatal
3. Table-driven tests catch a real bug
4. Subtests and running one with -run
5. Measuring test coverage
6. Writing and running a benchmark
7. A test helper with t.Helper()
8. Setup and teardown (TestMain and t.Cleanup)
9. Testing an HTTP handler with httptest
10. Comprehensive practice - a mockable discount service

Plus 3 bonus challenges: testing a generic function, fuzz testing, and golden-file testing.

**Time Commitment:** 5-6 hours
**When:** After reading README

### 4. **STUDY_GUIDE.md** (Visual Learning)
- Test file naming/structure diagram
- t.Error vs t.Fatal comparison table
- Table-driven test anatomy diagram
- Subtest naming/hierarchy diagram
- Coverage and benchmark output explained
- Mocking-via-interfaces diagram (ties to Level 15)
- Common mistakes guide

**Read Time:** 30-40 minutes
**When:** Use while doing exercises

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- Essential syntax reference
- Table-driven test template
- Coverage/benchmark command reference
- Common mistakes table

**Read Time:** 10-15 minutes
**When:** Quick lookup during practice

---

## 🎯 Recommended Learning Path

### Day 1: The Basics (2 hours)
1. Read **START_HERE.txt** (10 min)
2. Read **README.md** sections 1-2 (25 min)
3. Complete Exercises 1-2 (1 hour 25 min)

### Day 2: The Idiomatic Shape (2 hours)
1. Read **README.md** sections 3-4 (25 min)
2. Complete Exercises 3-4 (1.5 hours)

### Day 3: Measuring Your Tests (1.5 hours)
1. Read **README.md** sections 5-6 (25 min)
2. Complete Exercises 5-6 (1 hour 5 min)

### Day 4: Maintainable Suites (2 hours)
1. Read **README.md** sections 7-10 (35 min)
2. Complete Exercises 7-10 (1.5 hours)

### Day 5: Consolidation (1.5 hours)
1. Try the bonus challenges
2. Review with **QUICK_REFERENCE.md**
3. Make sure table-driven tests + subtests feel automatic

---

## 💡 Key Concepts At A Glance

### Naming and Running
```go
func TestAdd(t *testing.T) { }   // must start with Test + capital
```
```bash
go test -v ./...
```

### Error vs Fatal
```go
t.Errorf(...)   // fail, keep going
t.Fatalf(...)   // fail, stop this test now
```

### The Idiomatic Shape
```go
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) { ... })
}
```

### Mocking via Interfaces
```go
type Notifier interface{ Notify(string) error }
svc := NewOrderService(&fakeNotifier{}) // inject the fake in tests
```

---

## ✅ Prerequisites

Make sure you've completed **Level 25: Generics**

You need:
- ✅ Comfort with functions, methods, and error handling (Levels 11, 14, 16)
- ✅ Comfort with interfaces and implicit satisfaction (Level 15)
- ✅ Basic familiarity with the `net/http` package helps for section 9, but isn't required

---

## 🎓 Learning Objectives

By the end of Level 26, you'll be able to:

- ✅ Write and correctly name `TestXxx` functions
- ✅ Choose between t.Error and t.Fatal deliberately
- ✅ Write table-driven tests with subtests, and run a single one with -run
- ✅ Measure and interpret coverage and benchmark output
- ✅ Extract test helpers using t.Helper() correctly
- ✅ Set up and tear down test state with TestMain and t.Cleanup
- ✅ Test HTTP handlers with httptest
- ✅ Mock dependencies via small, consumer-defined interfaces

---

## 📊 Statistics

- **Main Theory:** README.md covering the testing package, table-driven tests, coverage, benchmarks, helpers, setup/teardown, httptest, and mocking
- **Exercises:** 10 detailed exercises + 3 bonus challenges (generics, fuzzing, golden files)
- **Study Guide:** naming diagrams, Error/Fatal comparison, table anatomy, coverage/benchmark explanations
- **Quick Reference:** One-page cheat sheet
- **Practice Time:** 5-6 hours
- **Difficulty:** ⭐⭐⭐⭐ (Advanced)

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for comprehensive theory
2. Use STUDY_GUIDE.md alongside for diagrams and comparison tables
3. Reference QUICK_REFERENCE.md for fast lookup

### For Practice
1. Read exercise description
2. Copy provided bash commands (or type them) - each exercise creates a code file AND a test file
3. Run `go test -v` (or the coverage/benchmark variant shown) and compare against the captured Expected Output
4. Understand what each step does, especially where output is deliberately shown failing

### For Review
1. Use QUICK_REFERENCE.md for quick recap
2. Refer to the common mistakes section
3. Practice writing a table-driven test with subtests from a blank file until it's automatic

---

## 🆘 Common Questions

**Q: Do I need a testing framework like in other languages?**
A: No - the standard library's `testing` package and `go test` command are the framework. Third-party assertion libraries exist but are optional, not required.

**Q: When should I use t.Fatal instead of t.Error?**
A: When continuing the test after a failure would be unsafe or meaningless - typically right after checking an `error` return, before you'd otherwise dereference a possibly-nil result.

**Q: Why does my table-driven test's failure only show a line number, not which case failed?**
A: You're missing `t.Run`. Wrap each case in a named subtest (`t.Run(tt.name, func(t *testing.T) {...})`) so failures are attributed to a specific, named case.

**Q: Is 100% test coverage the goal?**
A: No. Coverage tells you which statements *ran*, not whether the assertions checking them are meaningful. Aim for coverage of your important logic and edge cases, not a percentage number.

**Q: Do benchmark ns/op numbers mean anything across different machines?**
A: Only in relative terms. Use them to compare two approaches on the same machine in the same run - never treat an absolute `ns/op` value as portable or something to assert on directly.

**Q: Does fuzz testing actually work in this course's Go version?**
A: Yes - verified directly on Go 1.26.5. `go test -fuzz=FuzzXxx -fuzztime=5s` genuinely mutates inputs and can find real bugs (see the bonus challenge, where it found a genuine UTF-8 edge case in a `Reverse` function in well under a second).

---

## 🎯 Before Moving to Level 27

Make sure you can answer these questions:

- [ ] What are the exact naming rules for a test function to be recognized by `go test`?
- [ ] What's the difference between t.Error and t.Fatal, and when do you use each?
- [ ] Why is table-driven testing the idiomatic default in Go?
- [ ] What does `go test -cover` actually measure, and what doesn't it tell you?
- [ ] Why are benchmark numbers machine-dependent, and how should you use them anyway?
- [ ] Why does a test helper need `t.Helper()`?
- [ ] How do you mock a dependency in Go without a mocking framework?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (10 min)
   - Gets you oriented

2. **README.md** Sections 1-4 (30 min)
   - Test basics, Error/Fatal, table-driven tests, subtests

3. **README.md** Sections 5-8 (35 min)
   - Coverage, benchmarks, helpers, setup/teardown

4. **EXERCISES.md** Exercises 1-6 (2.5 hours)
   - Get hands-on with the core mechanics

5. **README.md** Sections 9-12 (25 min)
   - httptest, mocking, best practices, common mistakes

6. **STUDY_GUIDE.md** (30 min)
   - Study the diagrams and comparison tables

7. **EXERCISES.md** Exercises 7-10 + Bonus (2.5+ hours)
   - Deep practice with helpers, setup/teardown, httptest, mocking, and the bonus challenges

8. **QUICK_REFERENCE.md** (10 min)
   - Create your mental cheat sheet

---

## 🏆 Success Indicators

You've mastered Level 26 when:

- ✅ You reach for a table-driven test with subtests by default, even for "simple" functions
- ✅ You can explain t.Error vs t.Fatal without hesitation
- ✅ You know how to measure coverage and read a benchmark's ns/op line
- ✅ You always call t.Helper() in a test helper
- ✅ You can mock any dependency by defining a small interface at the point it's used
- ✅ You've completed 8+ exercises
- ✅ You can explain why Go's built-in testing tools are enough for most projects

---

## 🚀 What's Next?

After Level 26, you're ready for:

**Level 27: HTTP & REST APIs**
- Building real HTTP servers and handlers
- Routing, middleware, and request/response handling
- Testing HTTP services in depth, building directly on this level's httptest preview

---

## 💬 Key Takeaway

> **Go's testing tools ship in the standard library, so there's no excuse not to test: `_test.go` files, table-driven tests with subtests, and small interfaces for mocking will cover the overwhelming majority of what you ever need.**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Exercises
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Patterns
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Happy learning! Level 26 gives you the power to prove your programs work - and keep proving it as they grow! 🎉

*Estimated time to complete Level 26: 5-6 hours*
*Difficulty: ⭐⭐⭐⭐ (Advanced)*
*Next Level: Level 27 - HTTP & REST APIs*
