# Level 17: Packages & Modules - INDEX

Welcome to **Level 17: Packages & Modules**! This level starts a new phase of the course - applying Go to real, larger programs built from multiple modules, with deliberately designed package APIs.

---

## 📖 What You'll Learn

- ✅ Every directive in a go.mod file: `module`, `go`, `require`, `replace`, `exclude`
- ✅ How Go modules interpret semantic versioning, including the v2+ import-path rule
- ✅ How to run and read `go mod tidy`, `go list -m all`, `go mod why`, and `go mod graph`
- ✅ How to build a real multi-module project entirely offline with a `replace` directive
- ✅ How `internal/` enforces module-tree-only visibility at compile time
- ✅ How to design cohesive, minimal, focused package APIs
- ✅ How to develop multiple local modules together with `go.work`

---

## 🗂️ Level 17 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview and welcome message. Start here first!

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- go.mod anatomy, deep dive
- Semantic versioning as Go modules interpret it
- Practical module commands, run for real
- Building a local multi-module setup with `replace`
- Internal packages
- Package-level API design
- Multi-module workspaces with `go.work`
- Best practices
- Common mistakes

**Read Time:** 50-65 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 detailed, step-by-step exercises - every command actually run, offline:
1. go.mod anatomy walkthrough
2. A two-module local project with a replace directive
3. Inspecting the build — tidy, why, and graph
4. The internal/ directory — allowed import
5. The internal/ directory — compile error from outside
6. Avoiding circular imports
7. Package API design refactor
8. A multi-module workspace with go.work
9. replace vs. go.work, side by side
10. Comprehensive practice — a multi-module notes system

**Time Commitment:** 5-6 hours
**When:** After reading README

### 4. **STUDY_GUIDE.md** (Visual Learning)
- go.mod anatomy diagram
- Semantic versioning table
- replace vs. go.work comparison diagram
- internal/ visibility diagram
- Package design principles
- Common mistakes guide

**Read Time:** 35-45 minutes
**When:** Use while doing exercises

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- Essential syntax reference
- Module commands quick guide
- internal/ visibility rule
- replace vs. go.work comparison table
- Common mistakes table

**Read Time:** 10-15 minutes
**When:** Quick lookup during practice

---

## 🎯 Recommended Learning Path

### Day 1: go.mod and Local Multi-Module Basics (2.5 hours)
1. Read **START_HERE.txt** (10 min)
2. Read **README.md** sections 1-4 (35 min)
3. Complete Exercises 1-3 (1 hour 45 min)

### Day 2: internal/ and Circular Imports (2 hours)
1. Read **README.md** sections 5-6 (25 min)
2. Complete Exercises 4-6 (1.5 hours)

### Day 3: API Design and go.work (2 hours)
1. Read **README.md** sections 7-8 (25 min)
2. Complete Exercises 7-8 (1.5 hours)

### Day 4: Bringing It Together (1.5 hours)
1. Read **README.md** sections 9-10 (15 min)
2. Complete Exercises 9-10 (1 hour 15 min)

### Day 5: Consolidation (1 hour)
1. Try bonus challenges
2. Review with **QUICK_REFERENCE.md**
3. Make sure `replace` vs `go.work` and `internal/` visibility both feel automatic

---

## 💡 Key Concepts At A Glance

### go.mod Directives
```
module example.com/widgets   // identity + import prefix
go 1.26                      // minimum language version
require example.com/x v1.2.0 // a dependency
replace example.com/x => ../x// local override, main module only
exclude example.com/x v1.0.0 // blacklist a version, main module only
```

### Offline Multi-Module Development
```bash
go mod edit -require=example.com/x@v0.0.0-00010101000000-000000000000
go mod edit -replace=example.com/x=../x
```

### internal/ Visibility
```
example.com/owner/internal/x    → visible only within example.com/owner
```

### go.work
```bash
go work init
go work use ./greetings ./app   # no require/replace needed between them
```

---

## ✅ Prerequisites

Make sure you've completed **Level 16: Error Handling**

You need:
- ✅ Comfort with Go's error handling idioms
- ✅ Comfort with the basics from Level 2: one directory = one package, exported vs. unexported, `go mod init`
- ✅ A working `go` toolchain on your machine

---

## 🎓 Learning Objectives

By the end of Level 17, you'll be able to:

- ✅ Read and write every line of a go.mod file
- ✅ Explain semantic versioning as Go modules interpret it, including the v2+ rule
- ✅ Build a real multi-module project entirely offline with `replace`
- ✅ Run and interpret `go mod tidy`, `go list -m all`, `go mod why`, and `go mod graph`
- ✅ Explain why an `internal/` import fails outside its module tree, and reproduce the real compiler error
- ✅ Design a cohesive, minimally-exported package API and avoid circular imports
- ✅ Set up and use a `go.work` workspace across multiple local modules

---

## 📊 Statistics

- **Main Theory:** README.md covering go.mod anatomy, versioning, module commands, multi-module setups, internal/, API design, and go.work
- **Exercises:** 10 detailed exercises + 3 bonus challenges, every command run for real, fully offline
- **Study Guide:** go.mod diagram, versioning table, replace-vs-go.work diagram, internal/ diagram
- **Quick Reference:** One-page cheat sheet
- **Practice Time:** 5-6 hours
- **Difficulty:** ⭐⭐⭐ (Advanced)

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for comprehensive theory
2. Use STUDY_GUIDE.md alongside for diagrams and comparison tables
3. Reference QUICK_REFERENCE.md for fast lookup

### For Practice
1. Read exercise description
2. Copy provided bash commands (or type them)
3. Follow step-by-step instructions
4. Verify output matches expected results - every expected output in this level was actually run
5. Understand what each step does

### For Review
1. Use QUICK_REFERENCE.md for quick recap
2. Refer to the common mistakes section
3. Practice explaining `replace` vs `go.work` until the tradeoff is automatic

---

## 🆘 Common Questions

**Q: Why does this level avoid `go get` entirely?**
A: So every exercise works with no internet access. A `replace` directive (or a `go.work` file) points a module at a local directory on disk instead of a network registry - the exact same mechanics real projects use to develop unpublished modules together.

**Q: What's the difference between `replace` and `go.work`?**
A: `replace` lives in go.mod, is usually committed, and only affects the main module's own build. `go.work` is a separate, usually-uncommitted file that lets several local modules resolve each other directly, with no `require`/`replace` needed between them - better suited to ongoing multi-module development.

**Q: Why can't I import someone's `internal/` package if I add a `replace` pointing right at it?**
A: `replace`/`go.work` only make a package *findable* on disk - they never change Go's compiler-enforced rule that `internal/` packages are visible only to code rooted at `internal`'s parent directory. That's the whole point of `internal/`.

**Q: Do I need a real v2 module to understand the `/v2` import-path rule?**
A: No - it's explained conceptually in this level since it can't be demonstrated with real published modules offline. Recognize it when you see an import path ending in `/v2`, `/v3`, etc. in the wild.

**Q: How is this different from what Level 2 already taught about packages?**
A: Level 2 covered a single module with several packages. Level 17 covers *multiple modules*, real module commands, module-tree-scoped visibility (`internal/`), deliberate package API design, and developing several modules together with `go.work`.

---

## 🎯 Before Moving to Level 18

Make sure you can answer these questions:

- [ ] What does each line of a go.mod file mean?
- [ ] Why does a v2+ module need a `/vN` suffix in its import path?
- [ ] How do you wire two local, unpublished modules together offline?
- [ ] What does `go mod why` tell you that `go mod graph` doesn't, and vice versa?
- [ ] Why does an `internal/` import fail from outside its module tree, even with `replace`?
- [ ] When would you reach for `go.work` instead of a `replace` directive?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (10 min)
   - Gets you oriented

2. **README.md** Sections 1-4 (35 min)
   - go.mod anatomy, versioning, module commands

3. **README.md** Sections 5-8 (35 min)
   - Multi-module setups, internal/, API design, go.work

4. **EXERCISES.md** Exercises 1-5 (2.5 hours)
   - Get hands-on with go.mod, replace, and internal/

5. **STUDY_GUIDE.md** (35 min)
   - Study the diagrams and comparison tables

6. **EXERCISES.md** Exercises 6-10 (2.5+ hours)
   - Deep practice with API design and go.work

7. **QUICK_REFERENCE.md** (10 min)
   - Create your mental cheat sheet

---

## 🏆 Success Indicators

You've mastered Level 17 when:

- ✅ You can write a go.mod from memory, including `replace` and `exclude`
- ✅ You can wire two local modules together offline without hesitation
- ✅ You reach for `go mod why`/`go mod graph` instinctively when a dependency looks wrong
- ✅ You never try to route around an `internal/` restriction
- ✅ You default to small, cohesive packages instead of a catch-all `utils`
- ✅ You know when to reach for `go.work` instead of `replace`
- ✅ You've completed 8+ exercises
- ✅ You can explain modules vs. packages to someone else

---

## 🚀 What's Next?

After Level 17, you're ready for:

**Level 18: File Handling**
- Reading and writing files
- Working with paths and directories
- Buffered I/O
- Handling file-related errors

---

## 💬 Key Takeaway

> **A module is how Go versions and distributes code; `internal/` and deliberate package design are how you control what that code exposes. Master both, and you can structure software the way real, multi-module Go projects are built - all without ever needing the network.**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Exercises
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Patterns
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Happy learning! Level 17 gives you the tools to structure real, multi-module Go software! 🎉

*Estimated time to complete Level 17: 5-6 hours*
*Difficulty: ⭐⭐⭐ (Advanced)*
*Next Level: Level 18 - File Handling*
