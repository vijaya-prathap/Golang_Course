# Level 2: Go Program Structure - INDEX

Welcome to **Level 2: Go Program Structure**! This level teaches you how to organize Go code across multiple files and packages, which is essential for building real applications.

---

## 📖 What You'll Learn

- ✅ How to structure Go projects
- ✅ How to create and use packages
- ✅ How to import other packages
- ✅ How to control what's public and private
- ✅ How to initialize packages
- ✅ How to manage dependencies with go.mod
- ✅ How to avoid common structure mistakes

---

## 🗂️ Level 2 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview and welcome message. Start here first to understand the level's goals.

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- Why packages matter
- Package structure and organization
- Import statements and paths
- Exported vs unexported identifiers
- Package-level variables
- Init functions
- Real-world examples
- Troubleshooting
- Best practices

**Read Time:** 45-60 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 detailed, step-by-step exercises:
1. Create your first package
2. Import and use a package
3. Exported vs unexported
4. Package-level variables
5. Multiple files, one package
6. Project with multiple packages
7. Real project structure
8. Circular import problem (and solution!)
9. Init functions
10. Documentation

**Time Commitment:** 3-4 hours
**When:** After reading README

### 4. **STUDY_GUIDE.md** (Visual Learning)
- Visual diagrams and flowcharts
- Architecture patterns
- Common project structures
- Initialization order flowchart
- Package lifecycle diagrams
- Quick pattern reference
- Common mistakes and fixes

**Read Time:** 30-40 minutes
**When:** Use while doing exercises

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- Essential concepts (1-page summary)
- Common tasks and solutions
- Project template
- Key commands
- Visibility patterns
- Common errors table

**Read Time:** 10-15 minutes
**When:** Quick lookup during practice

---

## 🎯 Recommended Learning Path

### Day 1: Understand Concepts
1. Read **START_HERE.txt** (5 min)
2. Read **README.md** sections 1-4 (30 min)
3. Look at **STUDY_GUIDE.md** diagrams (15 min)
4. Start Exercise 1-2 (30 min)

### Day 2: Get Practical
1. Complete Exercises 3-6 (2 hours)
2. Refer to **STUDY_GUIDE.md** for patterns
3. Use **QUICK_REFERENCE.md** for syntax

### Day 3: Real Projects
1. Complete Exercises 7-8 (1 hour)
2. Work on Exercises 9-10 (1 hour)
3. Challenge yourself with bonus exercises

### Day 4: Consolidate
1. Re-read tricky sections of README
2. Try bonus challenges
3. Review your completed exercises
4. Make sure all concepts click

---

## 💡 Key Concepts At A Glance

### Package Structure
```
One directory = One package
One package = Can have multiple files
```

### Visibility (Capitalization Rule)
```
UPPERCASE = Exported (public)
lowercase = Unexported (private)
```

### Import Paths
```
import "myapp/packagename"
```

### Init Functions
```
init() runs automatically before main()
Useful for package setup
```

---

## ✅ Prerequisites

Make sure you've completed **Level 1: Go Introduction & Environment**

You need:
- ✅ Go installed
- ✅ A working Go development environment
- ✅ Ability to run `go run` and `go build`
- ✅ Basic understanding of Go programs

---

## 🎓 Learning Objectives

By the end of Level 2, you'll be able to:

- ✅ Create packages and organize code
- ✅ Import packages correctly
- ✅ Control visibility with capitalization
- ✅ Understand init() function lifecycle
- ✅ Organize real projects into packages
- ✅ Use go.mod effectively
- ✅ Avoid common structure mistakes
- ✅ Write package documentation

---

## 📊 Statistics

- **Total Content:** 55,000+ characters
- **Main Theory:** 20,800+ characters (README.md)
- **Exercises:** 15,000+ characters (10 detailed exercises)
- **Study Guide:** 8,200+ characters
- **Quick Reference:** 5,200+ characters
- **Practice Time:** 4-5 hours
- **Difficulty:** ⭐⭐ (Intermediate)

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for comprehensive theory
2. Use STUDY_GUIDE.md alongside for visual understanding
3. Reference QUICK_REFERENCE.md for fast lookup

### For Practice
1. Read exercise description
2. Copy provided bash commands
3. Follow step-by-step instructions
4. Verify output matches expected results
5. Understand what each step does

### For Review
1. Use QUICK_REFERENCE.md for quick recap
2. Refer to common mistakes section
3. Review project templates

---

## 🆘 Getting Help

### Common Questions

**Q: Why do I need packages?**
A: Packages organize your code and allow reuse. As projects grow, packages help manage complexity.

**Q: Does function case really matter?**
A: Yes! `GetName()` is exported (public). `getName()` is unexported (private). Go uses case to determine visibility.

**Q: What's the difference between "package" and "module"?**
A: Module is your entire project (in go.mod). Package is a directory with related Go files. One module contains many packages.

**Q: What happens if I don't use init()?**
A: Nothing! init() is optional. Use it only for package initialization needs.

---

## 🎯 Before Moving to Level 3

Make sure you can answer:
- [ ] What is a package?
- [ ] How do you create a package?
- [ ] How do you import a package?
- [ ] What's the difference between exported and unexported?
- [ ] When does init() run?
- [ ] What causes circular imports?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (5 min)
   - Gets you oriented
   
2. **README.md** Sections 1-3 (20 min)
   - Understand the "why"
   
3. **README.md** Sections 4-7 (20 min)
   - Learn the concepts
   
4. **EXERCISES.md** Exercises 1-3 (45 min)
   - Get hands-on
   
5. **STUDY_GUIDE.md** (30 min)
   - Visualize patterns
   
6. **EXERCISES.md** Exercises 4-10 (2+ hours)
   - Deep practice
   
7. **QUICK_REFERENCE.md** (10 min)
   - Create your mental cheat sheet

---

## 🏆 Success Indicators

You've mastered Level 2 when:

- ✅ You can create a package without instructions
- ✅ You understand exported vs unexported
- ✅ You can organize a project into multiple packages
- ✅ You know what init() does and when to use it
- ✅ You can avoid circular imports
- ✅ You completed 8+ exercises
- ✅ You can explain packages to someone else

---

## 🚀 What's Next?

After Level 2, you're ready for:

**Level 3: Variables & Data Types**
- Variable declaration
- Basic types (int, float64, string, bool)
- Type conversion
- Constants
- Zero values
- Named types

---

## 💬 Key Takeaway

> **Packages are how Go organizes code. Master packages, and you master Go program structure!**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Exercises
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Patterns
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Happy learning! You're building real Go skills now. 🎉

*Estimated time to complete Level 2: 4-5 hours*
*Difficulty: ⭐⭐ (Intermediate)*
*Next Level: Level 3 - Variables & Data Types*
