# Level 20: Goroutines - INDEX

Welcome to **Level 20: Goroutines**! This is where your programs stop running strictly top-to-bottom and start doing multiple independent things at once.

---

## 📖 What You'll Learn

- ✅ What a goroutine is, and concurrency vs. parallelism
- ✅ The classic beginner bug: `main()` exiting before a goroutine finishes
- ✅ `sync.WaitGroup` to wait for goroutines to complete
- ✅ Closures and loop variable capture in goroutines (verified on this project's Go version)
- ✅ Data races: what they are, seeing one for real, and catching one with `-race`
- ✅ Fixing races with a minimal `sync.Mutex` (a preview of Level 23)
- ✅ Goroutine leaks and why they're never garbage collected
- ✅ When goroutines are (and aren't) worth using

---

## 🗂️ Level 20 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview and welcome message. Start here first!

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- Goroutines, concurrency vs. parallelism
- The main-exits-too-early bug, demonstrated for real
- `sync.WaitGroup`, closures in loops, data races, `sync.Mutex`
- Goroutine leaks and when (not) to use goroutines
- Best practices
- Common mistakes

**Read Time:** 45-60 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 detailed, step-by-step exercises:
1. The main() exits too early problem
2. Fixing it with sync.WaitGroup
3. Waiting for multiple goroutines
4. Closures and loop variable capture
5. Data race - shared counter, no synchronization
6. Fixing the race with sync.Mutex
7. Goroutine leaks
8. Fan-out worker pattern - parallel partial sums
9. Mutex-protected results collection
10. Comprehensive practice - concurrent URL fetch simulation

**Time Commitment:** 5-6 hours
**When:** After reading README

### 4. **STUDY_GUIDE.md** (Visual Learning)
- Goroutine lifecycle diagram
- Concurrency vs. parallelism diagram
- sync.WaitGroup flow diagram
- Race condition timeline (two goroutines interleaving)
- Goroutine leak diagram
- Common mistakes guide

**Read Time:** 35-45 minutes
**When:** Use while doing exercises

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- Essential syntax reference
- WaitGroup and Mutex quick guide
- Common patterns and mistakes table

**Read Time:** 10-15 minutes
**When:** Quick lookup during practice

---

## 🎯 Recommended Learning Path

### Day 1: Goroutine Basics (2 hours)
1. Read **START_HERE.txt** (10 min)
2. Read **README.md** sections 1-2 (25 min)
3. Complete Exercise 1 (45 min)

### Day 2: WaitGroup & Loop Capture (2 hours)
1. Read **README.md** sections 3-4 (25 min)
2. Complete Exercises 2-4 (1.5 hours)

### Day 3: Races & Mutex (2 hours)
1. Read **README.md** sections 5-6 (30 min)
2. Complete Exercises 5-6 (1.5 hours)

### Day 4: Leaks & Fan-Out Patterns (2 hours)
1. Read **README.md** sections 7-8 (25 min)
2. Complete Exercises 7-9 (1.5 hours)

### Day 5: Comprehensive Practice & Consolidation (1.5 hours)
1. Complete Exercise 10 (45 min)
2. Try bonus challenges
3. Review with **QUICK_REFERENCE.md**

---

## 💡 Key Concepts At A Glance

### Launching a Goroutine
```go
go someFunc()   // returns immediately, runs concurrently
```

### sync.WaitGroup
```go
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    // work
}()
wg.Wait()
```

### Data Race → Fix With Mutex
```go
mu.Lock()
counter++   // safe: only one goroutine at a time
mu.Unlock()
```

### Goroutine Leak
```go
// blocks forever if nobody ever reads ch - a leak
go func() { ch <- result }()
```

---

## ✅ Prerequisites

Make sure you've completed **Level 19: JSON**

You need:
- ✅ Comfort with functions, closures, and Go's basic types
- ✅ Comfort with pointers (Level 12) - goroutines are frequently coordinated by pointer
- ✅ Familiarity with structuring and processing data (Levels 17-19)

---

## 🎓 Learning Objectives

By the end of Level 20, you'll be able to:

- ✅ Explain what a goroutine is and how it differs from an OS thread
- ✅ Distinguish concurrency from parallelism clearly
- ✅ Use `sync.WaitGroup` to synchronize a known number of goroutines
- ✅ Explain why loop variable capture in goroutines is safe on this project's Go version (1.26.5)
- ✅ Reproduce a data race and interpret `-race` detector output
- ✅ Fix a race with a minimal `sync.Mutex`
- ✅ Explain a goroutine leak and why it's never garbage collected
- ✅ Judge when goroutines are (and aren't) the right tool

---

## 📊 Statistics

- **Main Theory:** README.md covering goroutines, WaitGroup, races, mutexes, and leaks
- **Exercises:** 10 detailed exercises + 3 bonus challenges
- **Study Guide:** lifecycle, concurrency/parallelism, WaitGroup flow, and race-timeline diagrams
- **Quick Reference:** One-page cheat sheet
- **Practice Time:** 5-6 hours
- **Difficulty:** ⭐⭐⭐⭐ (Advanced)

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for comprehensive theory
2. Use STUDY_GUIDE.md alongside for the lifecycle and race-timeline diagrams
3. Reference QUICK_REFERENCE.md for fast lookup

### For Practice
1. Read exercise description
2. Copy provided bash commands (or type them)
3. Follow step-by-step instructions
4. Run each program multiple times where the exercise calls for it
5. Verify output matches the Expected Output section - and pay attention to which exercises are explicitly non-deterministic

### For Review
1. Use QUICK_REFERENCE.md for quick recap
2. Refer to the common mistakes section
3. Make sure you can explain WaitGroup, Mutex, and goroutine leaks to someone else

---

## 🆘 Common Questions

**Q: Why did my goroutine's `fmt.Println` never print anything?**
A: `main()` almost certainly returned before the goroutine got scheduled. This is Exercise 1's entire point - fix it with `sync.WaitGroup` (Exercise 2).

**Q: Is a goroutine the same as an OS thread?**
A: No. A goroutine is a lightweight, runtime-managed unit of execution; the Go runtime multiplexes many goroutines onto a smaller number of real OS threads.

**Q: Does `go run -race` always work?**
A: It requires cgo. It worked in this course's sandbox (verified with real captured output in Exercise 5) - but if cgo isn't available in your environment, you can still confirm a race exists from the varying, too-low output of a plain run.

**Q: Do I need to write `i := i` inside my loops before launching goroutines?**
A: Not on this project's Go version (1.26.5) or any Go 1.22+. Loops create a fresh copy of the loop variable every iteration now. See Exercise 4 and Level 11's closures-in-loops section.

**Q: What's the difference between a data race and a "race condition"?**
A: A data race specifically means unsynchronized concurrent memory access with at least one write - Go's tooling can detect it precisely. "Race condition" is a broader, informal term for any bug caused by timing.

---

## 🎯 Before Moving to Level 21

Make sure you can answer these questions:

- [ ] What is a goroutine, and how does it differ from an OS thread?
- [ ] What's the difference between concurrency and parallelism?
- [ ] Why does `sync.WaitGroup` make goroutine completion deterministic?
- [ ] Why is loop variable capture safe in goroutines on this project's Go version?
- [ ] What is a data race, and how does `sync.Mutex` fix one?
- [ ] What is a goroutine leak, and why isn't it garbage collected?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (10 min)
   - Gets you oriented

2. **README.md** Sections 1-4 (35 min)
   - Goroutines, the main-exits-early bug, WaitGroup, loop capture

3. **README.md** Sections 5-8 (35 min)
   - Races, Mutex preview, leaks, when to use goroutines

4. **EXERCISES.md** Exercises 1-5 (2.5 hours)
   - Get hands-on with the core bug, its fix, and races

5. **STUDY_GUIDE.md** (35 min)
   - Study the lifecycle and race-timeline diagrams

6. **EXERCISES.md** Exercises 6-10 (2.5+ hours)
   - Fix the race, explore leaks, and build the fan-out and fetch-simulation patterns

7. **QUICK_REFERENCE.md** (10 min)
   - Create your mental cheat sheet

---

## 🏆 Success Indicators

You've mastered Level 20 when:

- ✅ You never assume a background goroutine "had enough time" without explicit synchronization
- ✅ You reach for `sync.WaitGroup` automatically when launching a known number of goroutines
- ✅ You can explain why this project's Go version makes loop variable capture in goroutines safe
- ✅ You can reproduce a data race and read `-race` output (or reason about the bug without it)
- ✅ You know `sync.Mutex` is the minimal fix for a race, with full coverage still ahead in Level 23
- ✅ You can explain a goroutine leak without hesitation
- ✅ You can judge when a goroutine is worth its overhead and complexity, and when it isn't
- ✅ You've completed 8+ exercises

---

## 🚀 What's Next?

After Level 20, you're ready for:

**Level 21: Channels**
- Sending and receiving values between goroutines
- Buffered vs. unbuffered channels
- `select` for waiting on multiple channels
- Replacing this level's mutex-protected shared state with idiomatic channel communication

---

## 💬 Key Takeaway

> **A goroutine gives you concurrency for free with the `go` keyword - but every goroutine you launch needs an answer to two questions: "how do I know when it's done?" and "how do I know it will ever stop?"**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Exercises
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Patterns
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Happy learning! Level 20 gives your programs the power to do real work concurrently! 🎉

*Estimated time to complete Level 20: 5-6 hours*
*Difficulty: ⭐⭐⭐⭐ (Advanced)*
*Next Level: Level 21 - Channels*
