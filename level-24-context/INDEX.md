# Level 24: Context - INDEX

Welcome to **Level 24: Context**! This is where Level 22's "done channel" cancellation pattern gets formalized into Go's standard tool for propagating cancellation, deadlines, and request-scoped data through an entire call graph.

---

## 📖 What You'll Learn

- ✅ What context.Context carries, and why it formalizes the done-channel pattern
- ✅ context.Background() vs. context.TODO() - the two root contexts
- ✅ context.WithCancel for manual cancellation, verified with a real goroutine
- ✅ context.WithTimeout and context.WithDeadline, with real elapsed-time verification
- ✅ Distinguishing context.Canceled from context.DeadlineExceeded via ctx.Err()
- ✅ The ctx.Done()-select pattern for racing real work against cancellation
- ✅ context.WithValue with a typed key - and the string-key collision risk it avoids
- ✅ Why defer cancel() matters, backed by a real go vet warning
- ✅ Propagating one context through a call chain and a fan-out of goroutines
- ✅ A simulated HTTP-handler-like pattern previewing Level 27

---

## 🗂️ Level 24 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview and welcome message. Start here first!

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- What context.Context is for and why it exists
- Background()/TODO(), WithCancel, WithTimeout, WithDeadline
- ctx.Err(), ctx.Done(), and the standard select pattern
- WithValue and its correct, narrow use
- Always calling defer cancel()
- Propagation through call chains and fan-out goroutines
- A simulated HTTP-handler-like preview of Level 27
- Best practices
- Common mistakes

**Read Time:** 50-65 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 detailed, step-by-step exercises:
1. context.Background() and context.TODO()
2. context.WithCancel stopping a real goroutine
3. context.WithTimeout - real elapsed-time verification
4. context.WithDeadline - an absolute cutoff time
5. Distinguishing Canceled from DeadlineExceeded
6. The ctx.Done() select pattern - work vs. cancellation
7. context.WithValue - typed keys and the collision risk
8. Always defer cancel() - what leaks if you skip it (with a real go vet warning)
9. Propagating one context through fanned-out goroutines
10. Comprehensive practice - a cancelable pipeline job

Plus 3 bonus challenges: context-aware retry with backoff, a cancellable worker pool, and first-successful-result-wins.

**Time Commitment:** 5-6 hours
**When:** After reading README

### 4. **STUDY_GUIDE.md** (Visual Learning)
- The context tree / propagation diagram
- WithCancel vs. WithTimeout vs. WithDeadline comparison table
- The ctx.Done()-select pattern diagram
- context.Value key-collision-avoidance diagram (string vs. typed key)
- The defer cancel() leak, visualized
- Fan-out propagation diagram
- Common mistakes guide

**Read Time:** 35-45 minutes
**When:** Use while doing exercises

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- Essential syntax reference
- Real captured go vet warning for a discarded cancel
- Fan-out skeleton
- Common mistakes table

**Read Time:** 10-15 minutes
**When:** Quick lookup during practice

---

## 🎯 Recommended Learning Path

### Day 1: Root Contexts and Manual Cancellation (2 hours)
1. Read **START_HERE.txt** (10 min)
2. Read **README.md** sections 1-3 (30 min)
3. Complete Exercises 1-2 (1 hour 20 min)

### Day 2: Timeouts, Deadlines, and Err() (2 hours)
1. Read **README.md** sections 4-5 (25 min)
2. Complete Exercises 3-5 (1.5 hours)

### Day 3: The select Pattern and WithValue (1.5 hours)
1. Read **README.md** sections 6-7 (25 min)
2. Complete Exercises 6-7 (1 hour 5 min)

### Day 4: Discipline, Propagation, and the Handler Preview (2 hours)
1. Read **README.md** sections 8-10 (25 min)
2. Complete Exercises 8-10 (1.5+ hours)

### Day 5: Consolidation (1 hour)
1. Try bonus challenges
2. Review with **QUICK_REFERENCE.md**
3. Make sure the ctx.Done()-select pattern and the typed-key rule feel automatic

---

## 💡 Key Concepts At A Glance

### Root Contexts
```go
ctx := context.Background()   // real root
ctx := context.TODO()         // placeholder - plumbing still needed
```

### Deriving Cancelable Contexts
```go
ctx, cancel := context.WithCancel(parent)
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
ctx, cancel := context.WithDeadline(parent, someTime)
defer cancel() // always
```

### The select Pattern
```go
select {
case <-ctx.Done():
    return ctx.Err()
case result := <-workCh:
    // use result
}
```

### Typed WithValue Keys
```go
type requestIDKey struct{}
ctx = context.WithValue(ctx, requestIDKey{}, "req-123")
```

---

## ✅ Prerequisites

Make sure you've completed **Level 23: Mutex, WaitGroup & Atomic**

You need:
- ✅ Comfort protecting shared state with sync.Mutex
- ✅ Comfort using sync.WaitGroup to wait for goroutines to finish
- ✅ Familiarity with atomic operations for simple counters/flags
- ✅ Comfort with channels and select (Levels 21-22), including the "done channel" cancellation pattern

---

## 🎓 Learning Objectives

By the end of Level 24, you'll be able to:

- ✅ Explain what context.Context is for and why it formalizes the done-channel pattern
- ✅ Choose correctly between context.Background() and context.TODO()
- ✅ Cancel work manually with WithCancel and automatically with WithTimeout/WithDeadline
- ✅ Distinguish context.Canceled from context.DeadlineExceeded via ctx.Err()
- ✅ Write the standard ctx.Done()-select pattern for any cancellable operation
- ✅ Use context.WithValue correctly and explain the collision risk of string keys
- ✅ Always call the returned cancel function and explain precisely what leaks if you don't
- ✅ Propagate one context through a call chain and a fan-out of goroutines

---

## 📊 Statistics

- **Main Theory:** README.md covering root contexts, cancellation, timeouts, values, and propagation
- **Exercises:** 10 detailed exercises + 3 bonus challenges
- **Study Guide:** context tree diagram, comparison tables, select pattern and leak diagrams
- **Quick Reference:** One-page cheat sheet
- **Practice Time:** 5-6 hours
- **Difficulty:** ⭐⭐⭐⭐ (Advanced)

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for comprehensive theory
2. Use STUDY_GUIDE.md alongside for the context tree and select-pattern diagrams
3. Reference QUICK_REFERENCE.md for fast lookup

### For Practice
1. Read exercise description
2. Copy provided bash commands (or type them)
3. Follow step-by-step instructions
4. Verify output matches expected results - including the real go vet warning in Exercise 8
5. Understand what each step does

### For Review
1. Use QUICK_REFERENCE.md for quick recap
2. Refer to the common mistakes section
3. Practice writing the ctx.Done()-select pattern until it's automatic

---

## 🆘 Common Questions

**Q: What's the difference between context.Background() and context.TODO()?**
A: They behave identically at runtime - neither is ever canceled or has a deadline. The difference is purely documentation: Background() marks a genuine root; TODO() marks a spot where real context plumbing still needs to be wired up.

**Q: Do I really need defer cancel() if I'm using WithTimeout, which cancels itself anyway?**
A: Yes. Without it, the context and its internal timer stay alive for the *entire* timeout duration, even after your work finishes early. go vet actually flags a discarded WithTimeout cancel function as a real leak.

**Q: Why can't I just use a string as a context key?**
A: Two unrelated pieces of code that happen to pick the same string key will silently overwrite each other's values with no error, at runtime. A typed, unexported key type makes that collision impossible, since only your package can construct it.

**Q: How is this different from the "done channel" pattern from Level 22's select?**
A: It's the same idea, standardized. Instead of every function inventing its own `chan struct{}`, context.Context gives every function - yours, the standard library's, third-party packages' - the exact same Done()/Err()/Deadline() shape, so cancellation composes across an entire program.

**Q: Should I ever store a context in a struct field?**
A: No. A context describes one operation's lifetime, which rarely matches a struct's lifetime. Pass ctx as a parameter to each method call that needs it instead.

---

## 🎯 Before Moving to Level 25

Make sure you can answer these questions:

- [ ] What's the practical difference between context.Background() and context.TODO()?
- [ ] What does ctx.Err() return before and after cancellation, and how do the two sentinel errors differ?
- [ ] What's the standard select pattern for racing real work against ctx.Done()?
- [ ] Why must a context.Value key be a typed, unexported type rather than a string?
- [ ] What resource actually leaks if you skip defer cancel(), and why does it matter more in a server than a CLI tool?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (10 min)
   - Gets you oriented

2. **README.md** Sections 1-5 (35 min)
   - What context is for, root contexts, WithCancel, WithTimeout/WithDeadline, Err()

3. **README.md** Sections 6-10 (30 min)
   - The select pattern, WithValue, defer cancel(), propagation, the handler preview

4. **EXERCISES.md** Exercises 1-5 (2.5 hours)
   - Get hands-on with root contexts, cancellation, and timeouts

5. **STUDY_GUIDE.md** (40 min)
   - Study the context tree diagram and the select pattern diagram

6. **EXERCISES.md** Exercises 6-10 (2.5+ hours)
   - Deep practice with WithValue, defer cancel(), propagation, and the pipeline

7. **QUICK_REFERENCE.md** (10 min)
   - Create your mental cheat sheet

---

## 🏆 Success Indicators

You've mastered Level 24 when:

- ✅ You reach for context.Context as the first parameter of any function that does I/O or might run long
- ✅ You can explain Background() vs. TODO() without hesitation
- ✅ You've personally triggered and read the real go vet "context leak" warning
- ✅ You never store a context in a struct, and never use context.Value for business parameters
- ✅ You can build a fan-out of goroutines that all stop together on one cancel() call
- ✅ You've completed 8+ exercises
- ✅ You can explain context to someone else

---

## 🚀 What's Next?

After Level 24, you're ready for:

**Level 25: Generics**
- Type parameters for functions and types
- Type constraints and the comparable/any built-in constraints
- Writing generic data structures and algorithms
- When generics help vs. when a plain interface is simpler

---

## 💬 Key Takeaway

> **context.Context turns the "done channel" you built by hand in Level 22 into a standard, composable signal for cancellation, deadlines, and (sparingly) request-scoped values - shared cheaply across every function and goroutine working on one operation.**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Exercises
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Patterns
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Happy learning! Level 24 gives your concurrent programs a standard way to stop exactly when they're told to! 🎉

*Estimated time to complete Level 24: 5-6 hours*
*Difficulty: ⭐⭐⭐⭐ (Advanced)*
*Next Level: Level 25 - Generics*
