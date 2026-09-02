# Level 33: Logging - INDEX

Welcome to **Level 33: Logging**! This is where your programs stop shouting into a terminal only you can see, and start producing a structured, queryable record of what they're doing.

---

## 📖 What You'll Learn

- ✅ Why structured logging beats scattered `fmt.Println` calls
- ✅ The standard `log` package - what you've used incidentally, and its limits
- ✅ `log/slog`'s levels, handlers, and structured key-value attributes
- ✅ `TextHandler` vs `JSONHandler` - human-readable vs machine-parseable output
- ✅ Filtering noise with a minimum log level
- ✅ Contextual loggers with `.With()` and carrying them through `context.Context`
- ✅ Writing and testing a request-logging HTTP middleware

---

## 🗂️ Level 33 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview and welcome message. Start here first!

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- Why structured logging matters, and the standard `log` package's limits
- `log/slog` basics, handlers, levels, and structured attributes
- Contextual loggers, `context.Context` integration, and HTTP middleware
- `slog.SetDefault()` vs explicit loggers
- Best practices
- Common mistakes

**Read Time:** 45-60 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 detailed, step-by-step exercises + 3 bonus challenges:
1. Recap of the standard log package
2. Basic slog.Info/Warn/Error with attributes
3. TextHandler vs JSONHandler side by side
4. Filtering by minimum log level
5. Grouped attributes and typed helpers
6. A contextual logger with slog.With()
7. Stashing a logger in context.Context
8. A request-logging HTTP middleware (captured with httptest)
9. The log-and-return anti-pattern, demonstrated and fixed
10. Comprehensive practice - a small HTTP service with request-ID logging

**Time Commitment:** 4-5 hours
**When:** After reading README

### 4. **STUDY_GUIDE.md** (Visual Learning)
- log vs log/slog comparison table
- TextHandler vs JSONHandler output comparison
- Log-level filtering diagram
- The context-carried-logger pattern diagram
- Request-logging-middleware flow diagram
- Common mistakes guide

**Read Time:** 30-40 minutes
**When:** Use while doing exercises

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- Essential syntax reference
- Structured attributes and grouping quick guide
- Context-carried-logger pattern
- Common mistakes table

**Read Time:** 10-15 minutes
**When:** Quick lookup during practice

---

## 🎯 Recommended Learning Path

### Day 1: Why Logging Matters & the log Package (2 hours)
1. Read **START_HERE.txt** (10 min)
2. Read **README.md** sections 1-2 (20 min)
3. Complete Exercise 1 (30 min)

### Day 2: log/slog Basics & Handlers (2 hours)
1. Read **README.md** sections 3-5 (30 min)
2. Complete Exercises 2-4 (1.5 hours)

### Day 3: Structure & Contextual Loggers (2 hours)
1. Read **README.md** sections 6-8 (30 min)
2. Complete Exercises 5-7 (1.5 hours)

### Day 4: Middleware & Real-World Practice (2 hours)
1. Read **README.md** sections 9-10 (20 min)
2. Complete Exercises 8-10 (1.5 hours)

### Day 5: Consolidation (1 hour)
1. Try bonus challenges
2. Review with **QUICK_REFERENCE.md**
3. Make sure the context-carried-logger pattern feels automatic

---

## 💡 Key Concepts At A Glance

### log/slog Basics
```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
logger.Info("user logged in", "user_id", 42, "method", "password")
```

### Handlers
```go
slog.NewTextHandler(w, opts)   // human-readable key=value
slog.NewJSONHandler(w, opts)   // machine-parseable JSON
```

### Levels
```go
opts := &slog.HandlerOptions{Level: slog.LevelWarn}
// Debug and Info calls now produce no output
```

### Contextual Loggers
```go
reqLogger := base.With("request_id", id)
ctx := withLogger(r.Context(), reqLogger)
// any function reachable from ctx can now log with request_id attached
```

---

## ✅ Prerequisites

Make sure you've completed **Level 32: Authentication & JWT**

You need:
- ✅ Comfort with HTTP handlers and middleware (Level 27, Level 28)
- ✅ Comfort with `context.Context` and typed, unexported context keys (Level 24)
- ✅ Comfort with error wrapping and the log-and-return anti-pattern (Level 16)
- ✅ Basic familiarity with `httptest` (Level 26)

---

## 🎓 Learning Objectives

By the end of Level 33, you'll be able to:

- ✅ Explain why structured logging beats scattered `fmt.Println` calls
- ✅ Recap the standard `log` package and name its limitations
- ✅ Create loggers with `slog.NewTextHandler` and `slog.NewJSONHandler`
- ✅ Set and verify a minimum log level with `HandlerOptions`
- ✅ Use typed attribute helpers and `slog.Group`
- ✅ Build contextual loggers with `.With()` and carry them through `context.Context`
- ✅ Write and test a request-logging HTTP middleware with real captured output
- ✅ Identify and fix the log-and-return anti-pattern

---

## 📊 Statistics

- **Main Theory:** README.md covering the log package, log/slog, handlers, levels, structure, context, and middleware
- **Exercises:** 10 detailed exercises + 3 bonus challenges
- **Study Guide:** comparison tables, filtering diagram, context-logger diagram, middleware flow diagram
- **Quick Reference:** One-page cheat sheet
- **Practice Time:** 4-5 hours
- **Difficulty:** ⭐⭐⭐ (Intermediate-Advanced)

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for comprehensive theory
2. Use STUDY_GUIDE.md alongside for the comparison tables and diagrams
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
3. Practice the context-carried-logger pattern until it's automatic

---

## 🆘 Common Questions

**Q: Do I still need the `log` package after this level?**
A: For quick scripts and CLI tools, sure. For anything that runs as a service, prefer `log/slog` - it gives you levels and structure `log` never will.

**Q: Should I use TextHandler or JSONHandler in production?**
A: JSONHandler, almost always - it's what log aggregators expect. TextHandler is nicer for a human watching a local terminal during development.

**Q: Why store a logger in context.Context instead of just passing it as a parameter?**
A: Either works. Level 24 explains the tradeoff: `context.Value` suits cross-cutting, request-scoped data that many unrelated layers need (like a request ID or a request-scoped logger) without threading it through every signature by hand. If only one or two functions need it, an explicit parameter is more discoverable.

**Q: What's wrong with logging an error and then also returning it?**
A: The same failure ends up logged multiple times as it bubbles up through layers, making it impossible to tell from the logs whether one incident happened or several. Log exactly once, at the boundary - Level 16 introduced this rule for plain errors, and it applies identically to `slog`.

**Q: Should I call `slog.SetDefault()` everywhere?**
A: No - reserve it for `main()`. Pass explicit `*slog.Logger` instances into constructors elsewhere; it's more testable and lets different components use different loggers.

---

## 🎯 Before Moving to Level 34

Make sure you can answer these questions:

- [ ] What does `log/slog` give you that the standard `log` package doesn't?
- [ ] What's the difference between `TextHandler` and `JSONHandler`, and when do you use each?
- [ ] How does `HandlerOptions.Level` filter log output, and why does that matter in production?
- [ ] How do you carry a request-scoped logger through several function calls without a `*slog.Logger` parameter on every signature?
- [ ] Why is logging an error and then returning it considered a mistake?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (10 min)
   - Gets you oriented

2. **README.md** Sections 1-5 (30 min)
   - Why logging matters, the log package recap, log/slog basics, handlers

3. **README.md** Sections 6-10 (30 min)
   - Structured attributes, contextual loggers, context.Context, middleware, defaults

4. **EXERCISES.md** Exercises 1-5 (2 hours)
   - Get hands-on with the log package, slog basics, handlers, and levels

5. **STUDY_GUIDE.md** (30 min)
   - Study the comparison tables and the context-carried-logger diagram

6. **EXERCISES.md** Exercises 6-10 (2+ hours)
   - Deep practice with contextual loggers, middleware, and a full small service

7. **QUICK_REFERENCE.md** (10 min)
   - Create your mental cheat sheet

---

## 🏆 Success Indicators

You've mastered Level 33 when:

- ✅ You reach for `log/slog` by default for anything beyond a quick script
- ✅ You can explain the TextHandler/JSONHandler tradeoff without hesitation
- ✅ You never leave a service without an explicit minimum log level
- ✅ You can build and thread a context-carried logger through multiple layers
- ✅ You can write and test a logging middleware with deterministic, captured output
- ✅ You instinctively recognize and avoid the log-and-return anti-pattern
- ✅ You've completed 8+ exercises
- ✅ You can explain structured logging to someone else

---

## 🚀 What's Next?

After Level 33, you're ready for:

**Level 34: Configuration**
- Loading configuration from environment variables and files
- Setting the log level (and everything else) from external config instead of hardcoding it
- Validating configuration at startup instead of failing deep inside a request

---

## 💬 Key Takeaway

> **`log/slog` turns "print statements a human might glance at" into structured, leveled, queryable data a machine can index and a human can still read - carry it through `context.Context` so every layer of a request logs consistently, and log each failure exactly once.**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Exercises
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Patterns
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Happy learning! Level 33 gives your programs a voice that production systems can actually listen to! 🎉

*Estimated time to complete Level 33: 4-5 hours*
*Difficulty: ⭐⭐⭐ (Intermediate-Advanced)*
*Next Level: Level 34 - Configuration*
