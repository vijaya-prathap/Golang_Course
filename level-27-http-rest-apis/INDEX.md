# Level 27: HTTP & REST APIs - INDEX

Welcome to **Level 27: HTTP & REST APIs**! This is where your programs stop being scripts that print to a terminal and start being services that answer requests over the network.

---

## 📖 What You'll Learn

- ✅ The `http.Handler`/`http.HandlerFunc` interface underneath every Go HTTP server
- ✅ Routing with `http.ServeMux`'s real, modern method-specific and wildcard syntax
- ✅ Reading path values, query parameters, and JSON request bodies
- ✅ Writing JSON responses with correct headers and status codes
- ✅ Building a complete in-memory CRUD REST resource
- ✅ The middleware pattern for cross-cutting behavior like logging
- ✅ Testing handlers with `httptest.NewRecorder` and `httptest.NewServer`
- ✅ Using `context.WithTimeout` inside a handler to bound slow work
- ✅ Consistent JSON error responses mapped from Go errors

---

## 🗂️ Level 27 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview and welcome message. Start here first!

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- `net/http` fundamentals (`Handler`, `HandlerFunc`, `ResponseWriter`, `*Request`)
- `http.ServeMux` routing, verified against this project's actual Go version
- Reading and writing requests/responses
- A complete in-memory REST resource
- The middleware pattern
- Testing HTTP handlers in depth
- Context in handlers
- Consistent error responses
- Best practices and common mistakes

**Read Time:** 55-70 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 detailed, step-by-step exercises:
1. A "Hello, World" handler tested with `httptest.NewRecorder`
2. `ServeMux` routing with method and wildcard patterns
3. Reading query parameters
4. Decoding JSON requests and encoding JSON responses
5. An in-memory CRUD REST resource
6. A logging middleware, verified with captured output
7. A full `httptest.NewServer` round-trip test
8. A context-timeout-aware handler
9. Consistent JSON error responses
10. Comprehensive practice - a Task Manager REST API

Plus 3 bonus challenges (rate limiting, request IDs, pagination).

**Time Commitment:** 6-8 hours
**When:** After reading README

### 4. **STUDY_GUIDE.md** (Visual Learning)
- The `http.Handler`/`http.HandlerFunc` relationship diagram
- Request-lifecycle diagram (request in → routing → middleware → handler → response out)
- `ServeMux` pattern syntax table
- `httptest.NewRecorder` vs `httptest.NewServer` comparison
- Status code cheat table
- Common mistakes guide

**Read Time:** 35-45 minutes
**When:** Use while doing exercises

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- Essential syntax reference
- Routing, request/response, middleware, testing, and context snippets
- Status code table
- Common mistakes table

**Read Time:** 10-15 minutes
**When:** Quick lookup during practice

---

## 🎯 Recommended Learning Path

### Day 1: Fundamentals & Routing (2.5 hours)
1. Read **START_HERE.txt** (10 min)
2. Read **README.md** sections 1-2 (35 min)
3. Complete Exercises 1-2 (1 hour 45 min)

### Day 2: Reading & Writing HTTP (2 hours)
1. Read **README.md** sections 3-4 (25 min)
2. Complete Exercises 3-4 (1.5 hours)

### Day 3: A Real REST Resource (2 hours)
1. Read **README.md** section 5 (20 min)
2. Complete Exercise 5 (1.5 hours)

### Day 4: Middleware & Testing (2.5 hours)
1. Read **README.md** sections 6-7 (30 min)
2. Complete Exercises 6-7 (2 hours)

### Day 5: Context & Errors (2 hours)
1. Read **README.md** sections 8-9 (25 min)
2. Complete Exercises 8-9 (1.5 hours)

### Day 6: Comprehensive Practice (2 hours)
1. Complete Exercise 10 (1.5 hours)
2. Try bonus challenges

### Day 7: Consolidation (1 hour)
1. Review with **QUICK_REFERENCE.md** and **STUDY_GUIDE.md**
2. Make sure `ServeMux` routing syntax and the middleware pattern feel automatic

---

## 💡 Key Concepts At A Glance

### The One Interface
```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

### Modern Routing (Go 1.22+)
```go
mux.HandleFunc("GET /items/{id}", handler)
id := r.PathValue("id")
```

### The Middleware Pattern
```go
func Logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        next.ServeHTTP(w, r)
    })
}
```

### Testing, Two Ways
```go
rec := httptest.NewRecorder()            // fast, in-process
server := httptest.NewServer(mux)        // real (loopback) round trip
```

---

## ✅ Prerequisites

Make sure you've completed **Level 26: Testing**

You need:
- ✅ Comfort writing `TestXxx` functions and table-driven tests
- ✅ The `httptest.NewRecorder` preview from Level 26
- ✅ Comfort with JSON structs and tags (Level 19)
- ✅ Comfort with errors and `errors.Is`/wrapping (Level 16)
- ✅ Comfort with `sync.Mutex` (Level 23) and `context.Context` (Level 24)

---

## 🎓 Learning Objectives

By the end of Level 27, you'll be able to:

- ✅ Explain the `http.Handler`/`http.HandlerFunc` relationship
- ✅ Route requests with `http.ServeMux`'s real method-specific and wildcard syntax
- ✅ Read path values, query parameters, and JSON bodies from a request
- ✅ Write JSON responses with correct `Content-Type` and status codes
- ✅ Build and test a complete in-memory CRUD REST resource
- ✅ Write and verify `func(http.Handler) http.Handler` middleware
- ✅ Choose correctly between `httptest.NewRecorder` and `httptest.NewServer`
- ✅ Bound slow work in a handler with `context.WithTimeout`
- ✅ Return consistent JSON error responses mapped from Go errors

---

## 📊 Statistics

- **Main Theory:** README.md covering fundamentals, routing, CRUD, middleware, testing, context, and errors
- **Exercises:** 10 detailed exercises + 3 bonus challenges
- **Study Guide:** diagrams, pattern tables, status code cheat sheet
- **Quick Reference:** One-page cheat sheet
- **Practice Time:** 6-8 hours
- **Difficulty:** ⭐⭐⭐⭐ (Advanced)

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for comprehensive theory
2. Use STUDY_GUIDE.md alongside for diagrams and the status code table
3. Reference QUICK_REFERENCE.md for fast lookup

### For Practice
1. Read exercise description
2. Copy provided bash commands (or type them)
3. Follow step-by-step instructions
4. Run `go test -v` (or `go run`) and verify output matches expected results
5. Understand what each step does

### For Review
1. Use QUICK_REFERENCE.md for quick recap
2. Refer to the common mistakes section
3. Practice writing a `ServeMux` route and a middleware function from memory

---

## 🆘 Common Questions

**Q: Do I need a routing framework to build a real REST API in Go?**
A: Not necessarily - since Go 1.22, `http.ServeMux` handles method-specific routes and path wildcards natively. Frameworks like Gin (Level 28) add convenience (binding, validation, route groups) but the underlying model is exactly what this level teaches.

**Q: What's the difference between `httptest.NewRecorder` and `httptest.NewServer`?**
A: `NewRecorder` is an in-process fake `ResponseWriter` - fast, no network, best for unit-testing one handler. `NewServer` starts a real (loopback-only) TCP server - use it to prove routing and middleware work together end to end.

**Q: Why derive `context.WithTimeout` from `r.Context()` instead of `context.Background()`?**
A: `r.Context()` is already canceled automatically if the client disconnects. Deriving from it means your handler both respects your own timeout AND stops early if the client gives up first.

**Q: Why does an unmapped error fall through to 500 instead of some other code?**
A: `500 Internal Server Error` is the correct default for "something the API author didn't anticipate" - it's a deliberate fallback, not an oversight.

**Q: Is it safe to test all of this without a real network connection?**
A: Yes - every exercise in this level uses `httptest.NewRecorder` (no network at all) or `httptest.NewServer` (a real TCP listener, but bound to `127.0.0.1` only, never reachable from outside the machine).

---

## 🎯 Before Moving to Level 28

Make sure you can answer these questions:

- [ ] What single method must a type have to satisfy `http.Handler`?
- [ ] How do you register a route that only matches `DELETE` requests with a path wildcard?
- [ ] What's the correct order of operations when writing a response (headers, status, body)?
- [ ] What does `func(http.Handler) http.Handler` middleware actually wrap and return?
- [ ] When would you reach for `httptest.NewServer` instead of `httptest.NewRecorder`?
- [ ] Why should a handler's `context.WithTimeout` be derived from `r.Context()`?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (10 min)
   - Gets you oriented

2. **README.md** Sections 1-4 (35 min)
   - Fundamentals, routing, reading and writing requests

3. **EXERCISES.md** Exercises 1-4 (2.5 hours)
   - Get hands-on with handlers, routing, and JSON

4. **README.md** Sections 5-7 (35 min)
   - The CRUD resource, middleware, and testing in depth

5. **EXERCISES.md** Exercises 5-7 (3 hours)
   - Build the resource, write middleware, verify with a real round trip

6. **README.md** Sections 8-11 (25 min)
   - Context, errors, best practices, and mistakes

7. **EXERCISES.md** Exercises 8-10 + bonus (2.5+ hours)
   - Context-aware handlers, error mapping, and the comprehensive task manager

8. **STUDY_GUIDE.md** and **QUICK_REFERENCE.md** (as needed)
   - Diagrams, tables, and quick lookup throughout

---

## 🏆 Success Indicators

You've mastered Level 27 when:

- ✅ You can write a working `http.ServeMux`-routed handler from a blank file, no reference needed
- ✅ You reach for `httptest.NewRecorder` automatically when testing a new handler
- ✅ You can explain why `context.WithTimeout(r.Context(), ...)` beats `context.WithTimeout(context.Background(), ...)` in a handler
- ✅ You never guess a status code - you reach for the `http.StatusXxx` constant
- ✅ You can write a middleware function without looking up the pattern
- ✅ You've completed 8+ exercises
- ✅ You can explain the request lifecycle (routing → middleware → handler → response) to someone else

---

## 🚀 What's Next?

After Level 27, you're ready for:

**Level 28: Gin Framework**
- What a routing framework buys you over hand-rolled `http.ServeMux`
- Gin's context object, route groups, and middleware chains
- Binding and validating JSON with struct tags
- Building the same kind of REST API faster, with less boilerplate

---

## 💬 Key Takeaway

> **Every Go HTTP server, framework, and router is built from one interface - `http.Handler` - and one method - `ServeHTTP`. Learn that, and every framework you touch afterward is just convenience on top of what you already understand.**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Exercises
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Patterns
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Happy learning! Level 27 turns everything you've learned so far into a real, working web service! 🎉

*Estimated time to complete Level 27: 6-8 hours*
*Difficulty: ⭐⭐⭐⭐ (Advanced)*
*Next Level: Level 28 - Gin Framework*
