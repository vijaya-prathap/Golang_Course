# Level 23: Mutex, WaitGroup & Atomic - INDEX

Welcome to **Level 23: Mutex, WaitGroup & Atomic**! This is where Level 20's brief preview of locking becomes full mastery of Go's classic concurrency toolkit - `sync.Mutex`, `sync.RWMutex`, `sync.WaitGroup`, `sync/atomic`, and `sync.Once`.

---

## 📖 What You'll Learn

- ✅ `sync.Mutex` deep dive, with a fully re-verified `-race`-clean fix of Level 20's racy counter
- ✅ The `mu.Lock(); defer mu.Unlock()` pattern and why `defer` is idiomatic here
- ✅ `sync.RWMutex` for concurrent readers and exclusive writers
- ✅ `sync.WaitGroup`'s real rules: `Add()` timing and the by-value copying bug
- ✅ The modern typed `sync/atomic` API for lock-free counters
- ✅ Choosing between atomics and a mutex
- ✅ `sync.Once` for guaranteed one-time initialization
- ✅ Channels vs. mutexes - two valid concurrency philosophies
- ✅ A real, captured lock-ordering deadlock from Go's own runtime detector

---

## 🗂️ Level 23 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview and welcome message. Start here first!

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- `sync.Mutex`, `sync.RWMutex`, `sync.WaitGroup`, `sync/atomic`, `sync.Once`
- The WaitGroup Add()-timing bug and the by-value bug (with a real captured deadlock)
- A real lock-ordering deadlock, captured verbatim
- Best practices and common mistakes

**Read Time:** 55-70 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 detailed, step-by-step exercises:
1. Fixing the racy counter with sync.Mutex
2. The "always unlock" pattern with defer
3. sync.RWMutex with concurrent readers and writers
4. WaitGroup - Add() timing matters
5. The WaitGroup-passed-by-value bug
6. A typed atomic counter
7. Atomics vs mutex - choosing the right tool
8. sync.Once under concurrent calls
9. A real lock-ordering deadlock
10. Comprehensive practice - thread-safe cache/counter service

**Time Commitment:** 5-6 hours
**When:** After reading README

### 4. **STUDY_GUIDE.md** (Visual Learning)
- Mutex critical-section diagram
- RWMutex read/write access diagram
- WaitGroup Add/Done/Wait lifecycle
- Atomic-vs-mutex decision tree
- sync.Once flow diagram
- Lock-ordering deadlock diagram
- Common mistakes table

**Read Time:** 35-45 minutes
**When:** Use while doing exercises

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- Essential syntax reference for every primitive in this level
- Atomic vs mutex quick table
- Common mistakes table

**Read Time:** 10-15 minutes
**When:** Quick lookup during practice

---

## 🎯 Recommended Learning Path

### Day 1: Mutex Fundamentals (2 hours)
1. Read **START_HERE.txt** (10 min)
2. Read **README.md** sections 1-2 (25 min)
3. Complete Exercises 1-2 (1 hour 25 min)

### Day 2: RWMutex & WaitGroup Deep Dive (2.5 hours)
1. Read **README.md** sections 3-4 (30 min)
2. Complete Exercises 3-5 (2 hours)

### Day 3: Atomics & sync.Once (2 hours)
1. Read **README.md** sections 5-7 (30 min)
2. Complete Exercises 6-8 (1.5 hours)

### Day 4: Philosophy & the Real Deadlock (2 hours)
1. Read **README.md** sections 8-9 (25 min)
2. Complete Exercise 9 (45 min)
3. Complete Exercise 10 (50 min)

### Day 5: Consolidation (1 hour)
1. Try bonus challenges
2. Review with **QUICK_REFERENCE.md**
3. Make sure the atomic-vs-mutex decision feels automatic

---

## 💡 Key Concepts At A Glance

### sync.Mutex
```go
mu.Lock()
defer mu.Unlock()
// critical section
```

### sync.RWMutex
```go
mu.RLock()  // many concurrent readers
mu.Lock()   // one exclusive writer
```

### sync.WaitGroup
```go
wg.Add(1)     // BEFORE `go`, never inside the goroutine
go func() { defer wg.Done(); /* work */ }()
wg.Wait()
```

### sync/atomic
```go
var c atomic.Int64
c.Add(1); c.Load(); c.Store(0); c.CompareAndSwap(0, 5)
```

### sync.Once
```go
once.Do(func() { /* runs exactly once */ })
```

---

## ✅ Prerequisites

Make sure you've completed **Level 22: Select**

You need:
- ✅ Comfort with goroutines, `sync.WaitGroup`, and basic `sync.Mutex` from Level 20
- ✅ Comfort with channels (buffered/unbuffered, closing, ranging) from Level 21
- ✅ Comfort with `select` for waiting on multiple channels from Level 22
- ✅ Comfort with `defer` from Level 11/16

---

## 🎓 Learning Objectives

By the end of Level 23, you'll be able to:

- ✅ Fix a data race with `sync.Mutex` and confirm it with `go run -race`
- ✅ Apply the `mu.Lock(); defer mu.Unlock()` pattern by default
- ✅ Use `sync.RWMutex` where it genuinely helps, and recognize where it doesn't
- ✅ Avoid both major `sync.WaitGroup` mistakes: Add() timing and pass-by-value
- ✅ Use the typed `sync/atomic` API for lock-free counters and flags
- ✅ Choose correctly between an atomic and a mutex for any given piece of state
- ✅ Use `sync.Once` for safe one-time initialization
- ✅ Explain when to reach for a channel vs. a mutex
- ✅ Recognize and prevent a lock-ordering deadlock

---

## 📊 Statistics

- **Main Theory:** README.md covering all 5 synchronization primitives plus deadlocks
- **Exercises:** 10 detailed exercises + 3 bonus challenges
- **Study Guide:** critical-section diagrams, decision trees, and a captured deadlock dump
- **Quick Reference:** One-page cheat sheet
- **Practice Time:** 5-6 hours
- **Difficulty:** ⭐⭐⭐⭐ (Advanced)

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for comprehensive theory
2. Use STUDY_GUIDE.md alongside for the diagrams and decision trees
3. Reference QUICK_REFERENCE.md for fast lookup

### For Practice
1. Read exercise description
2. Copy provided bash commands (or type them)
3. Follow step-by-step instructions
4. Verify output matches expected results
5. Understand what each step does - especially the two exercises that deliberately deadlock

### For Review
1. Use QUICK_REFERENCE.md for quick recap
2. Refer to the common mistakes section
3. Practice the atomic-vs-mutex decision until it's automatic

---

## 🆘 Common Questions

**Q: Why does `-race` matter if my program "looks" correct?**
A: A data race can exist even when a particular run happens to print the "right" answer - absence of a visibly wrong result on one run is not proof of correctness. `-race` was confirmed working in this project's sandbox and catches races other than by luck.

**Q: Is it ever safe to call `Add()` from inside a goroutine?**
A: Only if you can guarantee `Wait()` won't be called until after that `Add()` has definitely executed - in practice this is fragile and easy to get wrong, so the safe rule is: always call `Add()` before the `go` statement.

**Q: Why does passing a `sync.WaitGroup` by value cause a deadlock instead of just a wrong answer?**
A: Because the function's `Done()` only decrements its own local copy - the original counter the caller's `Wait()` is watching never moves. If nothing else in the program is running, that permanent block is detected by Go's runtime as "all goroutines are asleep."

**Q: When should I reach for `sync/atomic` instead of `sync.Mutex`?**
A: When you have exactly one independent variable (a counter, a flag, a pointer) with nothing else that must stay consistent with it. The moment two or more related fields need to move together, use a mutex.

**Q: Is `sync.RWMutex` always faster than `sync.Mutex` for read-heavy code?**
A: No - this level's own informal benchmark measured plain `Mutex` as *faster* than `RWMutex` for very cheap reads. `RWMutex` pays off when reads are both frequent and substantial enough to offset its extra bookkeeping.

---

## 🎯 Before Moving to Level 24

Make sure you can answer these questions:

- [ ] What exactly does a mutex's critical section guarantee?
- [ ] Why does `defer mu.Unlock()` go immediately after `Lock()`?
- [ ] When does `RWMutex` actually help, versus just adding overhead?
- [ ] What are the two classic `sync.WaitGroup` mistakes, and what does each one actually cause?
- [ ] How do you decide between an atomic and a mutex for a piece of shared state?
- [ ] What does `sync.Once` guarantee, and how did this level verify it?
- [ ] What makes a lock-ordering deadlock happen, and how do you prevent it?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (10 min)
   - Gets you oriented

2. **README.md** Sections 1-4 (35 min)
   - Mutex, defer pattern, RWMutex, WaitGroup deep dive

3. **README.md** Sections 5-9 (35 min)
   - Atomics, atomic-vs-mutex, sync.Once, channels vs mutex, the real deadlock

4. **EXERCISES.md** Exercises 1-5 (2.5 hours)
   - Get hands-on with mutexes, RWMutex, and both WaitGroup bugs

5. **STUDY_GUIDE.md** (40 min)
   - Study the diagrams and decision trees

6. **EXERCISES.md** Exercises 6-10 (2.5+ hours)
   - Deep practice with atomics, sync.Once, the real deadlock, and the comprehensive service

7. **QUICK_REFERENCE.md** (10 min)
   - Create your mental cheat sheet

---

## 🏆 Success Indicators

You've mastered Level 23 when:

- ✅ You reach for `defer mu.Unlock()` automatically, every time
- ✅ You can explain why RWMutex isn't automatically faster than Mutex
- ✅ You never call `Add()` from inside a goroutine, and never pass a WaitGroup by value
- ✅ You instinctively know whether a piece of state needs an atomic or a mutex
- ✅ You can read a real Go deadlock dump and identify exactly which goroutine is stuck where
- ✅ You've completed 8+ exercises
- ✅ You can explain the difference between channels and mutexes to someone else

---

## 🚀 What's Next?

After Level 23, you're ready for:

**Level 24: Context**
- `context.Context` for cancellation, deadlines, and timeouts
- Passing request-scoped values safely
- Controlling goroutine lifetimes on top of the primitives from this level

---

## 💬 Key Takeaway

> **A mutex protects a critical section; an atomic protects a single variable; a WaitGroup waits for a known count of goroutines - mix them up (copy one by value, or grab two locks in different orders across goroutines) and Go's own runtime will tell you exactly how, with a real deadlock dump.**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Exercises
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Patterns
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Happy learning! Level 23 completes your classic Go concurrency toolkit! 🎉

*Estimated time to complete Level 23: 5-6 hours*
*Difficulty: ⭐⭐⭐⭐ (Advanced)*
*Next Level: Level 24 - Context*
