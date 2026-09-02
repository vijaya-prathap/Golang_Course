# Level 28: Gin Framework - INDEX

Welcome to **Level 28: Gin Framework**! This is where you meet your first third-party package - a real, production-grade web framework used across the Go ecosystem.

> ⚠️ **This level requires internet access.** Unlike every other level in this course, Level 28 needs `go get github.com/gin-gonic/gin@v1.12.0` to download a real dependency from the Go module proxy. Every exercise was verified against that exact pinned version - see START_HERE.txt and EXERCISES.md for details.

---

## 📖 What You'll Learn

- ✅ Why a framework exists on top of `net/http`, and what Gin specifically adds
- ✅ `gin.Default()` vs `gin.New()`
- ✅ Routing with path parameters, wildcards, and route groups
- ✅ Reading query params, path params, and validated JSON bodies
- ✅ Writing JSON, string, and status-only responses
- ✅ Writing and attaching custom middleware, global and scoped
- ✅ How Gin's built-in `Recovery()` middleware turns panics into clean responses
- ✅ Testing a Gin router with `httptest`

---

## 🗂️ Level 28 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview, the internet-access requirement, and welcome message. Start here first!

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- Why a framework, with a before/after comparison to Level 27
- Setup, routing, input, output, middleware, and recovery
- A complete REST resource rewritten in Gin
- Testing Gin handlers with `httptest`
- Gin vs plain `net/http` decision guide
- Best practices and common mistakes

**Read Time:** 55-70 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 detailed, step-by-step exercises:
1. `gin.Default()` vs `gin.New()`
2. Routing with path parameters
3. Route groups
4. Query parameters
5. JSON binding with validation
6. Custom middleware
7. Panic recovery demo
8. Full todos API rewrite
9. Testing Gin handlers with httptest
10. Comprehensive practice - a bookstore API

Plus 3 bonus challenges (request-ID middleware, a route group with its own middleware stack, converting a Level 27 stdlib handler to Gin).

**Time Commitment:** 5-6 hours
**When:** After reading README - **requires internet access** for each exercise's `go get`

### 4. **STUDY_GUIDE.md** (Visual Learning)
- Gin vs plain `net/http` comparison table
- `gin.Context` method reference table
- Middleware chain flow diagram (`c.Next()`)
- Binding/validation tag reference
- Common mistakes guide

**Read Time:** 35-45 minutes
**When:** Use while doing exercises

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- Essential syntax reference
- Routing, input, output, middleware, testing at a glance
- Common mistakes table

**Read Time:** 10-15 minutes
**When:** Quick lookup during practice

---

## 🎯 Recommended Learning Path

### Day 1: Why Gin, Setup, and Routing (2.5 hours)
1. Read **START_HERE.txt** (10 min)
2. Read **README.md** sections 1-3 (40 min)
3. Complete Exercises 1-3 (1 hour 40 min)

### Day 2: Input, Output & Middleware (2.5 hours)
1. Read **README.md** sections 4-6 (35 min)
2. Complete Exercises 4-6 (1 hour 55 min)

### Day 3: Recovery & a Real Resource (2 hours)
1. Read **README.md** sections 7-8 (20 min)
2. Complete Exercises 7-8 (1 hour 40 min)

### Day 4: Testing & Comprehensive Practice (2 hours)
1. Read **README.md** sections 9-10 (20 min)
2. Complete Exercises 9-10 (1 hour 40 min)

### Day 5: Consolidation (1 hour)
1. Try the bonus challenges
2. Review with **QUICK_REFERENCE.md**
3. Make sure the middleware chain (`c.Next()`) flow feels automatic

---

## 💡 Key Concepts At A Glance

### Setup
```go
router := gin.Default() // gin.New() + Logger() + Recovery()
router := gin.New()     // bare engine
```

### Routing
```go
router.GET("/users/:id", handler)
v1 := router.Group("/api/v1")
```

### Input & Output
```go
id := c.Param("id")
c.ShouldBindJSON(&input) // + binding:"required" struct tags
c.JSON(200, gin.H{"key": "value"})
```

### Middleware
```go
func mw() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
    }
}
router.Use(mw())
```

---

## ✅ Prerequisites

Make sure you've completed **Level 27: HTTP & REST APIs**

You need:
- ✅ Comfort with `net/http`, `ServeMux` routing, and middleware as function-wrapping
- ✅ Comfort with `httptest.NewRequest` / `httptest.NewRecorder`
- ✅ Comfort passing `context.Context` through handlers
- ✅ Level 16's `panic`/`recover` model
- ✅ Level 26's testing fundamentals (`go test`, table-driven tests, subtests)

---

## 🎓 Learning Objectives

By the end of Level 28, you'll be able to:

- ✅ Set up a Gin engine and explain `gin.Default()` vs `gin.New()`
- ✅ Route requests with path parameters, wildcards, and versioned groups
- ✅ Read query params, path params, and bind + validate JSON bodies
- ✅ Write JSON, string, and status-only responses
- ✅ Write custom middleware and control the chain with `c.Next()` / `c.Abort()`
- ✅ Explain how `Recovery()` prevents a single panic from crashing the server
- ✅ Build a complete CRUD REST resource with Gin
- ✅ Test a Gin router with `httptest`
- ✅ Decide when Gin is worth the dependency versus plain `net/http`

---

## 📊 Statistics

- **Main Theory:** README.md covering setup, routing, input/output, middleware, recovery, testing, and Gin vs stdlib
- **Exercises:** 10 detailed exercises + 3 bonus challenges
- **Study Guide:** comparison tables, context method reference, middleware flow diagram
- **Quick Reference:** One-page cheat sheet
- **Practice Time:** 5-6 hours
- **Difficulty:** ⭐⭐⭐⭐ (Advanced)
- **Special requirement:** Internet access for `go get github.com/gin-gonic/gin@v1.12.0` - the only level in this course that isn't fully offline

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for comprehensive theory
2. Use STUDY_GUIDE.md alongside for the comparison tables and chain diagram
3. Reference QUICK_REFERENCE.md for fast lookup

### For Practice
1. Read exercise description
2. Run the `go get` setup step (requires internet)
3. Copy provided bash commands (or type them)
4. Follow step-by-step instructions
5. Verify output matches expected results
6. Understand what each step does

### For Review
1. Use QUICK_REFERENCE.md for quick recap
2. Refer to the common mistakes section
3. Practice explaining the middleware chain (`c.Next()`) until it's automatic

---

## 🆘 Common Questions

**Q: Why does this level need internet access when nothing else in the course does?**
A: Gin is a real third-party package, not part of Go's standard library. `go get` has to fetch it (and its own dependencies) from the Go module proxy the first time. Once downloaded, Go caches it locally, so you won't need the network again for the same version.

**Q: Is `gin.Default()` doing something special I can't replicate?**
A: No - it's literally `gin.New()` with `gin.Logger()` and `gin.Recovery()` attached via `Use()`. Exercise 1 rebuilds it by hand to prove this.

**Q: What's the difference between `c.BindJSON` and `c.ShouldBindJSON`?**
A: `BindJSON` automatically writes a 400 response on failure. `ShouldBindJSON` just returns the `error` so you control the response format yourself - prefer this in real APIs.

**Q: Does Gin replace what I learned in Level 27?**
A: No - it builds on it. A `*gin.Context` wraps the exact same `*http.Request` and `http.ResponseWriter` from `net/http`. Everything about HTTP semantics, status codes, and testing with `httptest` still applies.

**Q: Should I always use Gin instead of plain net/http from now on?**
A: No. See README.md section 10 - Gin earns its keep on larger APIs with lots of routes and validation needs; plain `net/http` is often the better choice for smaller services or when minimizing dependencies matters.

---

## 🎯 Before Moving to Level 29

Make sure you can answer these questions:

- [ ] What does `gin.Default()` add over `gin.New()`?
- [ ] How do you read a path parameter? A query parameter? A JSON body?
- [ ] What does `binding:"required"` actually enforce, and what happens without it?
- [ ] How does `c.Next()` control the middleware chain?
- [ ] How does `Recovery()` change what happens when a handler panics?
- [ ] How do you test a Gin router without starting a real server?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (10 min)
   - Gets you oriented, including the internet-access requirement

2. **README.md** Sections 1-5 (35 min)
   - Why a framework, setup, routing, input, output

3. **README.md** Sections 6-9 (35 min)
   - Middleware, recovery, a complete resource, testing

4. **EXERCISES.md** Exercises 1-5 (2.5 hours)
   - Get hands-on with setup, routing, and validated input

5. **STUDY_GUIDE.md** (35 min)
   - Study the comparison tables and middleware chain diagram

6. **EXERCISES.md** Exercises 6-10 (3+ hours)
   - Deep practice with middleware, recovery, a full resource, and testing

7. **QUICK_REFERENCE.md** (10 min)
   - Create your mental cheat sheet

---

## 🏆 Success Indicators

You've mastered Level 28 when:

- ✅ You can set up a Gin engine and explain every middleware `gin.Default()` attaches
- ✅ You reach for route groups instead of repeating path prefixes and middleware per route
- ✅ You never forget a `binding` tag on a field that must be present
- ✅ You can trace a request through a middleware chain, including where `c.Next()` hands off control
- ✅ You can explain why `Recovery()` keeps one bad request from taking down the whole server
- ✅ You can test a Gin router exactly like you'd test any other `http.Handler`
- ✅ You've completed 8+ exercises
- ✅ You can explain when Gin is worth it and when plain `net/http` is the better call

---

## 🚀 What's Next?

After Level 28, you're ready for:

**Level 29: Database/SQL**
- `database/sql`, connection pools, and drivers
- Queries, scanning rows into structs, and prepared statements
- Wiring real, persistent storage underneath the Gin API you just built

---

## 💬 Key Takeaway

> **Gin doesn't replace what you learned in Level 27 - it removes the boilerplate around it. `*gin.Context` is still `*http.Request` and `http.ResponseWriter` underneath, `c.Next()` is still just a function call, and testing a Gin router is exactly like testing any other `http.Handler`.**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Exercises
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Patterns
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Happy learning! Level 28 turns everything you know about HTTP into a faster, less repetitive way to build real APIs! 🎉

*Estimated time to complete Level 28: 5-6 hours*
*Difficulty: ⭐⭐⭐⭐ (Advanced)*
*Next Level: Level 29 - Database/SQL*
