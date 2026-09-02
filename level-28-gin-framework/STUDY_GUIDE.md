# Level 28: Study Guide & Visual Reference

> ⚠️ **Reminder:** This level requires internet access (`go get github.com/gin-gonic/gin@v1.12.0`) - the one exception in an otherwise fully-offline course. Every example below was captured from real, running Gin `v1.12.0` code.

## 📚 Learning Path

### Week 1: Gin Fundamentals
```
Day 1:  Why a framework? gin.Default() vs gin.New()
Day 2:  Routing - methods, path params, wildcards, groups
Day 3:  Reading input - query params, path params, JSON binding
Day 4:  Writing responses, custom middleware
Day 5:  Recovery middleware and panic safety
Day 6:  A complete REST resource (todos rewrite)
Day 7:  Testing Gin handlers with httptest
```

### Week 2: Practice & Application
```
Day 1:  Exercises 1-3 (setup, path params, groups)
Day 2:  Exercises 4-5 (query params, JSON binding/validation)
Day 3:  Exercises 6-7 (middleware, panic recovery)
Day 4:  Exercises 8-9 (todos API rewrite, testing)
Day 5:  Exercise 10 (comprehensive bookstore API)
Day 6:  Bonus challenges
Day 7:  Review & consolidation
```

---

## 🎯 Gin vs Plain net/http (Level 27) Comparison

| Concern | net/http (Level 27) | Gin |
|---|---|---|
| Router | `http.ServeMux` (Go 1.22+: methods + `{param}`) | Radix-tree router, `:param` and `*wildcard` |
| Handler signature | `func(w http.ResponseWriter, r *http.Request)` | `func(c *gin.Context)` |
| Path parameter | `r.PathValue("id")` | `c.Param("id")` |
| Query parameter | `r.URL.Query().Get("q")` | `c.Query("q")` / `c.DefaultQuery("q", def)` |
| JSON response | `json.NewEncoder(w).Encode(v)` + manual header | `c.JSON(status, v)` |
| JSON request + validation | Manual `json.Decode` + manual `if` checks | `c.ShouldBindJSON(&v)` + `binding` struct tags |
| Middleware | Wrap `http.Handler` in another `http.Handler` | `gin.HandlerFunc` + `c.Next()` |
| Route grouping | Manual prefix string handling | `router.Group(prefix)` |
| Panic safety | Per-connection recover in `net/http` server (closes connection, no clean response by default) | `gin.Recovery()` returns a clean `500` |
| Testing | `httptest.NewRequest` + `httptest.NewRecorder` + call the handler | Identical - `router.ServeHTTP(w, req)` since `*gin.Engine` is an `http.Handler` |
| Dependencies | None (standard library only) | One external module (`gin-gonic/gin` + its own deps) |

---

## 🗺️ gin.Context Method Reference

| Method | Purpose |
|---|---|
| `c.Param(name)` | Read a `:name` path parameter (always a `string`) |
| `c.Query(name)` | Read a query-string value, `""` if missing |
| `c.DefaultQuery(name, def)` | Read a query-string value with a fallback |
| `c.GetQuery(name)` | Read a query-string value + whether it was present (`(string, bool)`) |
| `c.GetHeader(name)` | Read a request header |
| `c.Header(name, value)` | Set a response header |
| `c.ShouldBindJSON(&v)` | Decode + validate a JSON body, return the `error` yourself |
| `c.BindJSON(&v)` | Same, but auto-writes a 400 response on failure |
| `c.JSON(status, v)` | Write a JSON response with the right `Content-Type` |
| `c.String(status, format, args...)` | Write a plain-text response |
| `c.Status(status)` | Write only a status code, no body |
| `c.Set(key, value)` / `c.Get(key)` / `c.MustGet(key)` | Store/read per-request values, shared across middleware and the handler |
| `c.Next()` | Continue to the next middleware/handler in the chain |
| `c.Abort()` | Mark the chain as stopped (no later middleware/handler runs), write nothing |
| `c.AbortWithStatus(code)` | Abort + set only the status code |
| `c.AbortWithStatusJSON(code, v)` | Abort + write a JSON error body |
| `c.Writer.Status()` | Read back the status code that's about to be (or was) written |
| `c.Request` | The underlying `*http.Request` - everything from Level 27 still applies |

---

## 🔁 Middleware Chain Flow

```
router.Use(A)
router.GET("/x", B, C)     // per-route middleware B, then final handler C

Request for GET /x:

    A: before c.Next() ─────────────────┐
                                         ▼
        B: before c.Next() ─────────────┐
                                         ▼
            C (final handler) - writes the response
                                         │
        B: after c.Next() ◄─────────────┘
                                         │
    A: after c.Next() ◄─────────────────┘

    Response sent to the client
```

Every middleware in the chain gets a "before" phase (before it calls `c.Next()`) and an "after" phase (after `c.Next()` returns) - and they nest, like `defer` statements from Level 11, but running around `c.Next()` instead of around a `return`. If a middleware never calls `c.Next()`, nothing downstream ever runs and no response gets written unless that middleware wrote one itself.

**Aborting stops the chain early:**

```
router.Use(A)
router.GET("/x", authCheck, C)

Request with a missing auth header:

    A: before c.Next() ──┐
                          ▼
        authCheck: header missing -> c.AbortWithStatusJSON(401, ...); return
                          │
                          ✗  (C never runs - the chain stopped here)
                          │
    A: after c.Next() ◄──┘   (A still runs its "after" phase - c.Next() always returns)
```

---

## ✅ Binding & Validation Tag Reference

| Tag | Meaning |
|---|---|
| `binding:"required"` | Field must be present and non-zero-valued |
| `binding:"email"` | Must be a syntactically valid email address |
| `binding:"gt=0"` / `gte=0"` | Greater than / greater-than-or-equal to a number |
| `binding:"lt=100"` / `lte=100"` | Less than / less-than-or-equal to a number |
| `binding:"min=3"` / `max=50"` | Minimum/maximum length (strings) or value (numbers) |
| `binding:"oneof=admin user guest"` | Must be one of a fixed set of values |
| `binding:"required,email"` | Tags combine with commas - all must pass |

```go
type CreateUserInput struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age" binding:"required,gte=0,lte=130"`
}
```

A failed rule produces an error like:
```
Key: 'CreateUserInput.Email' Error:Field validation for 'Email' failed on the 'email' tag
```

**No `binding` tag = no validation.** A field with only a `json` tag silently accepts the zero value (`""`, `0`, `false`) if the client omits it entirely - see Common Mistakes below.

---

## 🚨 Common Mistakes

### Mistake 1: Missing `binding:"required"` - Silent Zero Values

```go
// ❌ WRONG - Age has no binding tag; a missing "age" silently becomes 0
type Input struct {
    Age int `json:"age"`
}

// ✅ RIGHT
type Input struct {
    Age int `json:"age" binding:"required"`
}
```

### Mistake 2: Binding the Same Request Body Twice

```go
// ❌ WRONG - the request body is a stream; the second bind gets nothing
c.ShouldBindJSON(&a)
c.ShouldBindJSON(&b) // EOF or zero values

// ✅ RIGHT - bind once, or use c.ShouldBindBodyWith to buffer it for reuse
```

### Mistake 3: Ignoring the Bind Error

```go
// ❌ WRONG
c.ShouldBindJSON(&input)
save(input) // input might be zero-valued and invalid

// ✅ RIGHT
if err := c.ShouldBindJSON(&input); err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
}
```

### Mistake 4: Attaching Group Middleware After Registering Routes

```go
// ❌ WRONG - too late for routes already added
v1 := router.Group("/api/v1")
v1.GET("/books", listBooks)
v1.Use(authMiddleware())

// ✅ RIGHT
v1 := router.Group("/api/v1")
v1.Use(authMiddleware())
v1.GET("/books", listBooks)
```

### Mistake 5: Forgetting c.Next() in Custom Middleware

```go
// ❌ WRONG - nothing downstream ever runs
func broken() gin.HandlerFunc {
    return func(c *gin.Context) { fmt.Println("ran") }
}

// ✅ RIGHT
func fixed() gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Println("ran")
        c.Next()
    }
}
```

### Mistake 6: Returning Instead of Aborting to Stop the Chain

```go
// ❌ WRONG - the handler still runs after this middleware "returns"
if unauthorized {
    return
}

// ✅ RIGHT
if unauthorized {
    c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
    return
}
```

---

## 📈 Progression Summary

### Understanding Level 28

Level 28 layers a framework on top of everything Level 27 already taught:

1. **Setup** — `gin.Default()` (with Logger + Recovery) vs `gin.New()` (bare)
2. **Routing** — methods, `:param`, `*wildcard`, `Group` for versioning
3. **Input** — query params, path params, `ShouldBindJSON` + validation tags
4. **Output** — `c.JSON`, `c.String`, `c.Status`
5. **Middleware** — `c.Next()`-based chains, global/per-route/per-group
6. **Recovery** — panics become clean 500s instead of crashes
7. **Testing** — identical `httptest` pattern as any other `http.Handler`

### Prerequisites for Level 29

Before moving to Level 29 (Database/SQL), you need:

- ✅ Comfortable setting up a Gin engine and registering routes
- ✅ Comfortable reading path params, query params, and binding JSON bodies with validation
- ✅ Can write and attach custom middleware, and explain `c.Next()` / `c.Abort()`
- ✅ Understand why `Recovery()` prevents one bad request from crashing the server
- ✅ Can test a Gin router with `httptest`, sharing a `SetupRouter` function

### Ready for Level 29?

Level 29 replaces the in-memory `map` storage used in this level's exercises with a real database:
- `database/sql`, connection pools, and drivers
- Real persistence for the todos/bookstore-style resources you just built
- Scanning rows into structs, prepared statements, and transactions

---

## ✅ Checklist Before Level 29

- [ ] Can explain what `gin.Default()` adds over `gin.New()`
- [ ] Can register a route with a path parameter and read it with `c.Param`
- [ ] Can build a versioned API with `router.Group`
- [ ] Can validate a JSON request body with `binding` struct tags
- [ ] Can write a custom middleware using `c.Next()`, both global and per-route
- [ ] Can explain how `Recovery()` turns a panic into a `500` instead of a crash
- [ ] Can test a Gin router with `httptest.NewRequest` + `httptest.NewRecorder`
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The Trade
Gin doesn't do anything `net/http` can't - it just removes repeated boilerplate (method dispatch, path parsing, JSON encode/decode, middleware chaining) at the cost of one dependency.

### The Context
`*gin.Context` wraps the same `*http.Request` / `http.ResponseWriter` you already know from Level 27 - `c.Request` and `c.Writer` are right there if you ever need the raw stdlib types.

### The Chain
Middleware is just a handler that calls `c.Next()`. Code before `c.Next()` runs going in; code after runs coming out - the same nesting shape as `defer`, but built around one function call instead of a `return`.

### The Safety Net
`gin.Default()`'s `Recovery()` middleware means a single panicking handler returns a clean `500` and never takes the whole server down - directly building on Level 16's `panic`/`recover`.

---

## 📚 Next Level

Level 29: Database/SQL
- `database/sql`, connection pools, and drivers
- Queries, scanning rows into structs, prepared statements
- Wiring real persistence underneath the Gin API you just built

You've got Gin down! Keep going! 🚀
