# Level 34: Configuration - INDEX

Welcome to **Level 34: Configuration**! This is where your programs stop hardcoding settings and start adapting to whatever environment they're run in - dev, staging, or production - using the exact same compiled binary.

---

## 📖 What You'll Learn

- ✅ Why the same binary needs different settings across environments without recompiling
- ✅ `os.Getenv` vs `os.LookupEnv` - distinguishing "unset" from "set to empty"
- ✅ The `flag` package: defaults, `flag.Parse()`, and command-line overrides
- ✅ Designing a typed `Config` struct with a single `LoadConfig()` function
- ✅ Config precedence: flags override environment variables, which override defaults
- ✅ Fail-fast validation at startup instead of discovering problems in production
- ✅ Type conversion (`strconv`, `time.ParseDuration`) and error handling for env vars
- ✅ Hand-rolling a `.env` file parser for local-dev convenience

---

## 🗂️ Level 34 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview and welcome message. Start here first!

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- Why configuration matters
- `os.Getenv`/`os.LookupEnv` and the `flag` package
- Designing a `Config` struct and `LoadConfig()`
- Precedence, fail-fast validation, and type conversion
- Hand-rolled `.env` file support
- Third-party options as further reading only
- Best practices
- Common mistakes

**Read Time:** 45-60 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 detailed, step-by-step exercises:
1. os.Getenv vs os.LookupEnv
2. The flag package with defaults
3. A Config struct populated from environment variables
4. Config precedence (flags > env > defaults)
5. Fail-fast validation at startup
6. Type-converting numeric and boolean env vars
7. Hand-rolled .env file parser
8. Loading a typed Config struct directly from a .env file
9. .env files must never override real environment variables
10. Comprehensive practice - full service config loading

**Time Commitment:** 4-5 hours
**When:** After reading README

### 4. **STUDY_GUIDE.md** (Visual Learning)
- os.Getenv vs os.LookupEnv comparison table
- flag package lifecycle diagram
- Config precedence flow diagram
- Fail-fast validation flow diagram
- Type conversion cheat sheet
- .env file parsing diagram
- Common mistakes guide

**Read Time:** 30-40 minutes
**When:** Use while doing exercises

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- Essential syntax reference
- Precedence resolver pattern
- .env parser snippet
- Common mistakes table

**Read Time:** 10-15 minutes
**When:** Quick lookup during practice

---

## 🎯 Recommended Learning Path

### Day 1: Reading Config (2 hours)
1. Read **START_HERE.txt** (10 min)
2. Read **README.md** sections 1-3 (30 min)
3. Complete Exercises 1-2 (1 hour 20 min)

### Day 2: Structuring & Prioritizing Config (2 hours)
1. Read **README.md** sections 4-5 (25 min)
2. Complete Exercises 3-4 (1.5 hours)

### Day 3: Validation & Type Safety (1.5 hours)
1. Read **README.md** sections 6-7 (25 min)
2. Complete Exercises 5-6 (1 hour 5 min)

### Day 4: .env Files & Real-World Practice (2 hours)
1. Read **README.md** sections 8-9 (20 min)
2. Complete Exercises 7-10 (1.5+ hours)

### Day 5: Consolidation (1 hour)
1. Try bonus challenges
2. Review with **QUICK_REFERENCE.md**
3. Make sure the flags > env > defaults precedence feels automatic

---

## 💡 Key Concepts At A Glance

### os.Getenv vs os.LookupEnv
```go
value := os.Getenv("KEY")          // "" for unset OR empty
value, ok := os.LookupEnv("KEY")   // ok distinguishes the two
```

### flag Package
```go
port := flag.Int("port", 8080, "server port")
flag.Parse()
```

### Config Struct
```go
type Config struct { Port string }
func LoadConfig() Config { /* ... */ }
```

### Precedence
```go
flags > environment variables > hardcoded defaults
```

### Fail-Fast Validation
```go
if err := cfg.Validate(); err != nil {
    log.Fatal(err)
}
```

---

## ✅ Prerequisites

Make sure you've completed **Level 33: Logging**

You need:
- ✅ Comfort with `log.Fatal` and structured error reporting (Level 33)
- ✅ Comfort with file I/O basics: `os.Open`, `bufio.Scanner` (Level 18)
- ✅ Comfort with `strconv` string/number conversions (Level 3)
- ✅ Basic familiarity with structs and methods (Levels 13-14)

---

## 🎓 Learning Objectives

By the end of Level 34, you'll be able to:

- ✅ Explain why one binary needs different settings per environment
- ✅ Distinguish `os.Getenv` from `os.LookupEnv` and know when the difference matters
- ✅ Define and parse command-line flags with sensible defaults
- ✅ Design a typed `Config` struct with a single `LoadConfig()` entry point
- ✅ Implement and test flags > env > defaults precedence
- ✅ Validate configuration at startup and fail fast on missing/invalid settings
- ✅ Convert environment variable strings into typed values with proper error handling
- ✅ Hand-roll a `.env` file parser, including the "don't clobber real env vars" convention

---

## 📊 Statistics

- **Main Theory:** README.md covering env vars, flags, Config structs, precedence, validation, and .env parsing
- **Exercises:** 10 detailed exercises + 3 bonus challenges
- **Study Guide:** comparison tables, precedence diagrams, .env parsing diagram
- **Quick Reference:** One-page cheat sheet
- **Practice Time:** 4-5 hours
- **Difficulty:** ⭐⭐⭐ (Intermediate-Advanced)

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for comprehensive theory
2. Use STUDY_GUIDE.md alongside for the precedence diagram and comparison tables
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
3. Practice the precedence resolver and .env parser until they're automatic

---

## 🆘 Common Questions

**Q: Why not just use os.Getenv everywhere?**
A: You can for simple cases, but scattered `os.Getenv` calls mean there's no single place to see every setting a program depends on, no way to enforce defaults consistently, and no way to distinguish "unset" from "empty" where that matters. A `Config` struct plus `LoadConfig()` fixes all three.

**Q: Do I need a config library like Viper for this level?**
A: No - every exercise in this level uses only the standard library (`os`, `flag`, `strconv`, `time`, `bufio`) and runs fully offline. Third-party libraries are mentioned only as further reading for when a project's config genuinely outgrows hand-rolled loading.

**Q: Why do flags override environment variables, and not the other way around?**
A: Environment variables are typically set once per deployment (by Docker, Kubernetes, or a shell profile) and apply broadly. Flags are passed explicitly for one specific run, which makes them the most deliberate, most recent instruction - so they win.

**Q: Why shouldn't a .env file override a real environment variable?**
A: Because a real, exported environment variable (say, one set by your deployment platform) was set deliberately and should never be silently overridden by a leftover `.env` file sitting in the working directory. `.env` is only meant to fill gaps for local development convenience.

**Q: What does "fail fast" actually save me from?**
A: A missing `DATABASE_URL` caught by `Validate()` at startup is a one-line, obvious error message before the server ever binds a port. The same missing value, uncaught, might not surface until the first real request tries to use it - possibly in production, possibly much later.

---

## 🎯 Before Moving to Level 35

Make sure you can answer these questions:

- [ ] What's the difference between `os.Getenv` and `os.LookupEnv`, and when does it matter?
- [ ] How do you define a flag with a default value, and when is that default actually available?
- [ ] What's the precedence order for flags, environment variables, and defaults?
- [ ] Why should configuration be validated at startup instead of lazily?
- [ ] Why should a `.env` file never override an already-set real environment variable?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (10 min)
   - Gets you oriented

2. **README.md** Sections 1-4 (25 min)
   - Why configuration matters, env vars, flags, Config structs

3. **README.md** Sections 5-9 (35 min)
   - Precedence, validation, type conversion, .env files, third-party options

4. **EXERCISES.md** Exercises 1-5 (2 hours)
   - Get hands-on with env vars, flags, Config structs, precedence, and validation

5. **STUDY_GUIDE.md** (30 min)
   - Study the precedence diagram and comparison tables

6. **EXERCISES.md** Exercises 6-10 (2+ hours)
   - Deep practice with type conversion, .env parsing, and the comprehensive final exercise

7. **QUICK_REFERENCE.md** (10 min)
   - Create your mental cheat sheet

---

## 🏆 Success Indicators

You've mastered Level 34 when:

- ✅ You reach for `LookupEnv` by default whenever "unset vs empty" could matter
- ✅ You can explain the flags > env > defaults precedence order without hesitation
- ✅ You never let a program start successfully with a missing required setting
- ✅ You can hand-roll a `.env` parser from memory, including comments and blank lines
- ✅ You instinctively keep config-loading code separate from business logic
- ✅ You've completed 8+ exercises
- ✅ You can explain configuration design to someone else

---

## 🚀 What's Next?

After Level 34, you're ready for:

**Level 35: Docker**
- Packaging your configured Go binary into a container image
- Passing environment variables into a container at `docker run` time
- Why the configuration habits from this level make a Go binary "container-friendly"

---

## 💬 Key Takeaway

> **One compiled binary, deployed unchanged everywhere - configuration is what lets it behave correctly in dev, staging, and production without a single recompile.**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Exercises
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Patterns
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Happy learning! Level 34 gives your programs the power to adapt to any environment without ever being recompiled! 🎉

*Estimated time to complete Level 34: 4-5 hours*
*Difficulty: ⭐⭐⭐ (Intermediate-Advanced)*
*Next Level: Level 35 - Docker*
