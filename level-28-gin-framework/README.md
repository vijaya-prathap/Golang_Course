# Level 28: Gin Framework - Complete Guide

## Introduction

Welcome to Level 28! You've mastered `Level 27: HTTP & REST APIs` - routing with `net/http`'s `ServeMux`, middleware as function-wrapping, `httptest`, and passing `context.Context` through handlers, all using only the standard library. Now it's time to meet **Gin** (`github.com/gin-gonic/gin`), the most widely used web framework in the Go ecosystem, and see what a framework buys you once your API grows past a handful of routes.

> ⚠️ **This level is different from every other level in this course.** Levels 0-27 work **100% offline** - every exercise runs with nothing but the Go standard library. Gin is a **third-party package**, so this level **requires internet access** the first time you set up each exercise, to run `go get github.com/gin-gonic/gin@v1.12.0` and download it (and its own dependencies) from the Go module proxy. Once downloaded, Go caches it locally (`$GOPATH/pkg/mod`) so you won't need the network again for the same version. Every exercise below was written, run with `go run` / `go test`, and its output captured for real using Gin `v1.12.0` - nothing here is hand-computed. If you're offline, come back to this level once you have a connection; every other level in the course does not have this requirement.

Go's standard library gives you everything you need to build a correct, production-grade HTTP API - Level 27 proved that. So why would you reach for a framework at all?

### Before: A Level-27-Style Handler

```go
// Raw net/http (Level 27 style): manual method check, manual path parsing,
// manual JSON encoding, manual error handling for a malformed body.
func getUserHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }
    id := r.PathValue("id") // Go 1.22+ ServeMux path values
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"user_id": id})
}

mux := http.NewServeMux()
mux.HandleFunc("GET /users/{id}", getUserHandler)
```

### After: The Gin Equivalent

```go
// Gin: the method is part of router.GET, the path parameter is decoded for
// you, and c.JSON handles both the header and the encoding in one call.
func getUserHandler(c *gin.Context) {
    id := c.Param("id")
    c.JSON(http.StatusOK, gin.H{"user_id": id})
}

router := gin.Default()
router.GET("/users/:id", getUserHandler)
```

Both versions do the same thing. The Gin version is shorter, but the real payoff shows up as an API grows: dozens of routes, nested versioned groups, JSON bodies that need validation, and middleware stacked several layers deep. That's what this level is about.

---

## Table of Contents

1. [Why a Framework?](#why-a-framework)
2. [Setup: gin.Default() vs gin.New()](#setup-gindefault-vs-ginnew)
3. [Routing: Methods, Path Parameters, and Groups](#routing-methods-path-parameters-and-groups)
4. [Reading Input: Query Params, Path Params, and JSON Bodies](#reading-input-query-params-path-params-and-json-bodies)
5. [Writing Responses](#writing-responses)
6. [Middleware](#middleware)
7. [Built-in Recovery Middleware](#built-in-recovery-middleware)
8. [A Complete REST Resource in Gin](#a-complete-rest-resource-in-gin)
9. [Testing Gin Handlers](#testing-gin-handlers)
10. [Gin vs Plain net/http: When to Use Which](#gin-vs-plain-nethttp-when-to-use-which)
11. [Best Practices](#best-practices)
12. [Common Mistakes](#common-mistakes)

---

## Why a Framework?

`net/http`'s `ServeMux` (since Go 1.22) can match methods and path parameters, but it stops there - everything else (JSON validation, structured error responses, composable middleware groups) is code you write and re-write on every project. A framework packages up those repeated decisions:

- **Faster, richer routing.** Gin's router is a **radix tree** (a compressed prefix tree) instead of `ServeMux`'s linear-ish pattern matching. Lookups stay fast even with thousands of routes, and it natively supports parameters (`:id`) and catch-alls (`*filepath`) in the same tree.
- **Built-in JSON binding and validation.** `c.ShouldBindJSON` decodes a request body into a struct *and* runs validation rules declared as struct tags (`binding:"required"`), in one call.
- **Composable middleware.** `engine.Use()` for global middleware, per-route middleware, and per-group middleware, all sharing one `c.Next()`-based chain.
- **Less boilerplate per handler.** No manual `w.Header().Set` + `json.NewEncoder(w).Encode(...)` pairs, no manual method dispatch, no manual path-parameter parsing.

None of this makes Gin *capable* of something `net/http` can't do - Level 27 already builds real REST APIs with the standard library alone. Gin just means writing less repetitive code to get there, at the cost of one external dependency.

---

## Setup: gin.Default() vs gin.New()

Every Gin app starts with an **engine** (`*gin.Engine`), which is also an `http.Handler` - meaning you can still hand it to `http.ListenAndServe` if you ever need to mix it with standard-library code.

```go
router := gin.Default() // gin.New() + Logger() + Recovery() middleware attached
```

```go
router := gin.New() // a bare engine - no middleware attached at all
```

`gin.Default()` is not magic - it is *defined* as `gin.New()` with two middleware functions attached via `Use()`. You can verify this by rebuilding it by hand:

```go
handRolled := gin.New()
handRolled.Use(gin.Logger(), gin.Recovery())
// handRolled now behaves exactly like gin.Default()
```

| | `gin.New()` | `gin.Default()` |
|---|---|---|
| Request logging | ❌ none | ✅ `gin.Logger()` |
| Panic recovery | ❌ none | ✅ `gin.Recovery()` |
| Use when | You want full control over middleware (e.g. a custom logger) | The common case - most real projects start here |

Once you have an engine, register routes on it and call `router.Run(":8080")` to start listening (internally, `Run` just calls `http.ListenAndServe(addr, router)`). The exercises in this level mostly skip `Run` and instead simulate requests directly against the engine with `httptest`, so the output is captured and deterministic - the same technique Level 26 previewed and Level 27 used in full.

---

## Routing: Methods, Path Parameters, and Groups

### HTTP Methods

Gin has one method on the engine per HTTP verb - no manual `if r.Method == ...` needed:

```go
router.GET("/books", listBooks)
router.POST("/books", createBook)
router.PUT("/books/:id", updateBook)
router.DELETE("/books/:id", deleteBook)
```

### Path Parameters

A segment prefixed with `:` is a named parameter, read with `c.Param`:

```go
router.GET("/users/:id", func(c *gin.Context) {
    id := c.Param("id") // string - convert yourself if you need an int
    c.JSON(200, gin.H{"user_id": id})
})
```

Multiple parameters work the same way:

```go
router.GET("/users/:id/posts/:postId", func(c *gin.Context) {
    c.JSON(200, gin.H{
        "user_id": c.Param("id"),
        "post_id": c.Param("postId"),
    })
})
```

A segment prefixed with `*` is a catch-all wildcard that matches the rest of the path, **including the leading slash**:

```go
router.GET("/files/*filepath", func(c *gin.Context) {
    c.JSON(200, gin.H{"filepath": c.Param("filepath")})
})
// GET /files/images/2024/photo.png -> filepath = "/images/2024/photo.png"
```

### Route Groups

A `Group` shares a path prefix (and, as you'll see in the middleware section, a middleware stack) across many routes - ideal for API versioning:

```go
v1 := router.Group("/api/v1")
{
    v1.GET("/users", listUsersV1)
    v1.GET("/users/:id", getUserV1)
}

v2 := router.Group("/api/v2")
{
    v2.GET("/users", listUsersV2) // same path, different version, no collision
}
```

Groups nest, too - `admin.Group("/reports")` builds `/admin/reports/...` on top of an existing `/admin` group. The `{ }` blocks above are just a Go convention for visually scoping which routes belong to which group - they carry no special meaning to the compiler.

---

## Reading Input: Query Params, Path Params, and JSON Bodies

### Query Parameters

```go
router.GET("/search", func(c *gin.Context) {
    q := c.Query("q")                     // "" if missing
    limit := c.DefaultQuery("limit", "10") // "10" if missing
    c.JSON(200, gin.H{"query": q, "limit": limit})
})
```

`GET /search?q=golang&limit=5` -> `{"limit":"5","query":"golang"}`
`GET /search?q=golang` (no `limit`) -> `{"limit":"10","query":"golang"}`

If you need to distinguish "missing" from "present but empty", use `c.GetQuery`, which also returns a `bool`:

```go
q, wasProvided := c.GetQuery("q")
```

### Path Parameters

Already covered above - `c.Param("name")` always returns a `string`.

### JSON Request Bodies

`c.ShouldBindJSON` decodes the body into a struct and enforces any `binding` tags on that struct in one step:

```go
type CreateUserInput struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age" binding:"required,gte=0,lte=130"`
}

router.POST("/users", func(c *gin.Context) {
    var input CreateUserInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    c.JSON(201, gin.H{"created": input})
})
```

`binding` tags come from the [`go-playground/validator`](https://github.com/go-playground/validator) package, which Gin uses internally. Common tags: `required`, `email`, `gt=0` / `gte=0` / `lt=100` / `lte=100`, `min=3` / `max=50`, `oneof=admin user guest`. A validation failure returns an error whose `.Error()` message names the offending field and rule, e.g.:

```
Key: 'CreateUserInput.Email' Error:Field validation for 'Email' failed on the 'email' tag
```

`ShouldBindJSON` vs `BindJSON`: `BindJSON` automatically calls `c.AbortWithError(400, err)` on failure (it decides the response for you); `ShouldBindJSON` just returns the `error` and lets *you* decide the status code and response shape. Prefer `ShouldBindJSON` when you want control over the error response format - which is almost always, in a real API.

---

## Writing Responses

```go
c.JSON(http.StatusOK, gin.H{"message": "pong"})   // JSON body + Content-Type header
c.String(http.StatusOK, "plain text, %s!", "here") // text/plain
c.Status(http.StatusNoContent)                     // status only, no body
```

`gin.H` is nothing exotic - it's just `type H map[string]any`, a shorthand for building an ad-hoc JSON object without declaring a struct. Because it's a Go map, `encoding/json` serializes its keys in **sorted alphabetical order** - that's why `gin.H{"user_id": id, "post_id": postID}` comes out as `{"post_id":...,"user_id":...}` in the exercises below. A named struct, by contrast, serializes fields in **declaration order** - reach for a struct instead of `gin.H` whenever field order in the response matters to you or to API consumers.

---

## Middleware

A Gin middleware is a `gin.HandlerFunc` - the exact same type as a regular handler (`func(c *gin.Context)`). What makes it middleware is that it calls `c.Next()` to hand control to whatever comes next in the chain, then (optionally) runs more code afterward:

```go
func timingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Println("before handler")
        c.Next() // runs every remaining middleware AND the final handler
        fmt.Println("after handler, status:", c.Writer.Status())
    }
}
```

Code before `c.Next()` runs on the way **in**; code after it runs on the way **out**, once everything downstream has finished - the same "wrap the next thing" shape Level 27 built by hand with plain functions, just given a name (`c.Next()`) and a shared `*gin.Context` to pass state through.

### Global Middleware

```go
router.Use(timingMiddleware()) // runs for every route on this engine
```

### Per-Route Middleware

Pass it as an extra argument before the handler - Gin accepts any number of `gin.HandlerFunc` per route:

```go
router.GET("/private", requireHeaderMiddleware("X-API-Key"), privateHandler)
```

### Per-Group Middleware

```go
protected := router.Group("/api/v1")
protected.Use(authMiddleware(apiKey)) // applies to every route added to this group
```

### Stopping the Chain: c.Abort()

If a middleware decides the request shouldn't continue (e.g. missing auth), it must call one of the `Abort*` methods instead of just returning - otherwise Gin has no way to know the chain should stop:

```go
func requireHeaderMiddleware(header string) gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.GetHeader(header) == "" {
            c.AbortWithStatusJSON(401, gin.H{"error": "missing " + header})
            return // still return - Abort* doesn't stop the current function
        }
        c.Next()
    }
}
```

`c.AbortWithStatusJSON` sets the status, writes the JSON body, and marks the context aborted (so no later middleware or handler runs) - all in one call.

---

## Built-in Recovery Middleware

Level 16 covered `panic`/`recover`: an unrecovered panic unwinds the stack, prints a trace, and crashes the program with exit status 2. In a raw `net/http` server (Level 27), Go's own HTTP server recovers a panicking handler *per-connection* so one bad request doesn't take down the whole process - but by default it just logs the panic to stderr and closes the connection; the client gets a bare connection reset, not a clean response.

Gin's `Recovery()` middleware (included automatically by `gin.Default()`) does this more gracefully: it recovers the panic, logs it, and turns it into a proper `500 Internal Server Error` response - so callers always get a real HTTP response, and the server keeps serving later requests without missing a beat.

```go
router := gin.Default() // includes Recovery()

router.GET("/divide/:a/:b", func(c *gin.Context) {
    var a, b int
    fmt.Sscanf(c.Param("a"), "%d", &a)
    fmt.Sscanf(c.Param("b"), "%d", &b)
    result := a / b // panics if b is 0
    c.JSON(200, gin.H{"result": result})
})
```

`GET /divide/10/0` triggers a real `panic` (integer division by zero) inside the handler. With `Recovery()` in the chain, the response is `500` with an empty body - not a crash - and the very next request to `/healthy` succeeds normally, proving the server survived. `Recovery()` also logs the recovered panic and its stack trace to `gin.DefaultErrorWriter` (`os.Stderr` by default) so you don't lose visibility into what went wrong, even though the client just sees a clean 500.

If you want a custom JSON shape for recovered panics instead of an empty 500 body, use `gin.CustomRecoveryWithWriter` (or `gin.CustomRecovery` on modern Gin) instead of the default `gin.Recovery()`.

---

## A Complete REST Resource in Gin

Here's a full CRUD "todos" resource - the same shape of resource Level 27 built with raw `net/http`, rewritten with Gin's routing, binding, and grouping:

```go
type Todo struct {
    ID    int    `json:"id"`
    Title string `json:"title"`
    Done  bool   `json:"done"`
}

type CreateTodoInput struct {
    Title string `json:"title" binding:"required"`
}

var (
    todos  = map[int]Todo{}
    nextID = 1
)

func setupRouter() *gin.Engine {
    router := gin.New()
    router.Use(gin.Recovery())

    api := router.Group("/api/v1")
    {
        api.GET("/todos", func(c *gin.Context) {
            result := make([]Todo, 0, len(todos))
            for id := 1; id < nextID; id++ {
                if t, ok := todos[id]; ok {
                    result = append(result, t)
                }
            }
            c.JSON(200, gin.H{"todos": result})
        })

        api.POST("/todos", func(c *gin.Context) {
            var input CreateTodoInput
            if err := c.ShouldBindJSON(&input); err != nil {
                c.JSON(400, gin.H{"error": err.Error()})
                return
            }
            t := Todo{ID: nextID, Title: input.Title}
            todos[t.ID] = t
            nextID++
            c.JSON(201, t)
        })
        // ...PUT and DELETE follow the same shape
    }
    return router
}
```

Compare this to a `net/http` version: no manual method switch, no manual JSON decode-and-check-error boilerplate for the request body, and versioning is a `Group` call instead of string-prefix path handling. Exercise 8 builds this resource in full, with every HTTP verb, and Exercise 10 goes further with grouped auth middleware and validated bodies.

---

## Testing Gin Handlers

Level 26 previewed `httptest.NewRequest` + `httptest.NewRecorder`; Level 27 used them against raw `net/http` handlers. Testing a Gin router uses the exact same two building blocks - Gin's `*gin.Engine` implements `http.Handler`, so `router.ServeHTTP(recorder, request)` works exactly like calling any other handler's `ServeHTTP`:

```go
func TestPingRoute(t *testing.T) {
    gin.SetMode(gin.TestMode) // silences Gin's startup debug output during tests

    router := SetupRouter()

    req := httptest.NewRequest(http.MethodGet, "/ping", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("status = %d; want %d", w.Code, http.StatusOK)
    }
    want := `{"message":"pong"}`
    if got := w.Body.String(); got != want {
        t.Errorf("body = %q; want %q", got, want)
    }
}
```

`gin.SetMode(gin.TestMode)` is the third mode alongside `gin.DebugMode` (the default - prints a route-registration banner) and `gin.ReleaseMode` (silences it for production). Setting it in tests keeps `go test -v` output focused on your own `t.Run` lines instead of Gin's internal debug banner.

The pattern that keeps handlers testable: put your routing setup in an exported function (`func SetupRouter() *gin.Engine`) instead of directly inside `main()`. Both `main()` and every test then call the same function, so tests always exercise the exact routes and middleware your real server runs.

---

## Gin vs Plain net/http: When to Use Which

| Use Gin when... | Use plain net/http when... |
|---|---|
| The API has many routes, especially with path/query parameters | The service is small - a handful of endpoints |
| Request bodies need struct-tag validation | You want zero non-standard-library dependencies |
| You want composable, groupable middleware without writing the chaining logic yourself | You're building a library, not a service, and don't want to force a framework on consumers |
| The team already knows Gin, or the project will grow | You need the absolute smallest binary/dependency footprint |

Neither is "more correct" - `net/http` (as of Go 1.22+) can express most of what Gin does, just with more code per route. Gin trades one dependency for less repeated boilerplate as an API's surface area grows. Many production Go services use `net/http` directly and are perfectly happy; many others use Gin (or a sibling like Echo or Chi) and are equally happy. Level 27's skills transfer directly either way - a Gin handler's `*gin.Context` wraps the same `*http.Request` and `http.ResponseWriter` you already know.

---

## Best Practices

### 1. Prefer `ShouldBindJSON` Over `BindJSON`

```go
// ✅ Good - you control the error response shape
if err := c.ShouldBindJSON(&input); err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
}

// ⚠️ Works, but Gin decides the response format for you
if err := c.BindJSON(&input); err != nil {
    return // BindJSON already wrote a 400 response
}
```

### 2. Build the Router in One Function, Shared by main() and Tests

```go
// ✅ Good
func SetupRouter() *gin.Engine { /* ... */ }

func main() {
    SetupRouter().Run(":8080")
}
```

Tests then call `SetupRouter()` directly, guaranteeing they exercise real routes and real middleware - not a hand-copied subset.

### 3. Use Route Groups for Versioning and Shared Middleware, Not String Concatenation

```go
// ✅ Good
v1 := router.Group("/api/v1")
v1.Use(authMiddleware())

// ❌ Repetitive and error-prone
router.GET("/api/v1/users", authMiddleware(), listUsers)
router.GET("/api/v1/posts", authMiddleware(), listPosts)
```

### 4. Always Check the Error From ShouldBindJSON

Never assume a request body matches your struct - see Common Mistakes below for what happens if you don't.

### 5. Use gin.TestMode in Tests, gin.ReleaseMode in Production

```go
gin.SetMode(gin.TestMode)    // in tests
gin.SetMode(gin.ReleaseMode) // in production main() - silences the debug banner
```

### 6. Give Every Struct Field That Must Be Present a `binding:"required"` Tag

A field with no `binding` tag silently accepts the zero value (`""`, `0`, `false`) if the client omits it - it will never produce a validation error on its own.

---

## Common Mistakes

### Mistake 1: Forgetting `binding:"required"` - Silent Zero Values

```go
// ❌ WRONG - Age has no binding tag
type CreateUserInput struct {
    Name string `json:"name" binding:"required"`
    Age  int    `json:"age"` // missing "age" in the request body silently becomes 0
}
```

```go
// ✅ RIGHT - now a missing age is caught explicitly
type CreateUserInput struct {
    Name string `json:"name" binding:"required"`
    Age  int    `json:"age" binding:"required"`
}
```

### Mistake 2: Calling BindJSON/ShouldBindJSON Twice on the Same Request

The request body is an `io.Reader` - once it's been read, it's drained. A second bind call on the same request gets `EOF` or an empty result, not the same data again.

```go
// ❌ WRONG - the second call sees an already-drained body
var a, b Input
c.ShouldBindJSON(&a)
c.ShouldBindJSON(&b) // fails or binds zero values

// ✅ RIGHT - bind once into one struct (or use c.ShouldBindBodyWith to buffer it)
var input Input
c.ShouldBindJSON(&input)
```

### Mistake 3: Ignoring the Error From ShouldBindJSON

```go
// ❌ WRONG - input may be partially/zero-valued and you'll never know
c.ShouldBindJSON(&input)
saveUser(input)

// ✅ RIGHT
if err := c.ShouldBindJSON(&input); err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
}
saveUser(input)
```

### Mistake 4: Blocking Group Middleware Patterns

Attaching middleware with `Use()` **after** routes have already been registered on that group does not retroactively apply it to those earlier routes:

```go
// ❌ WRONG - /books was already registered before Use() ran
v1 := router.Group("/api/v1")
v1.GET("/books", listBooks)
v1.Use(authMiddleware()) // too late for /books

// ✅ RIGHT - attach middleware before registering routes that need it
v1 := router.Group("/api/v1")
v1.Use(authMiddleware())
v1.GET("/books", listBooks)
```

### Mistake 5: Forgetting c.Next() Silently Skips the Rest of the Chain

```go
// ❌ WRONG - handler never runs, request just hangs with no response written
func brokenMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Println("middleware ran")
        // forgot c.Next() - nothing downstream executes
    }
}

// ✅ RIGHT
func fixedMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Println("middleware ran")
        c.Next()
    }
}
```

### Mistake 6: Returning From a Middleware Without Calling Abort

```go
// ❌ WRONG - returning stops THIS function, but doesn't tell Gin the chain
// should stop, and doesn't write any response - the client hangs waiting
func brokenAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.GetHeader("X-API-Key") == "" {
            return // handler still runs next!
        }
        c.Next()
    }
}

// ✅ RIGHT
func fixedAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.GetHeader("X-API-Key") == "" {
            c.AbortWithStatusJSON(401, gin.H{"error": "missing API key"})
            return
        }
        c.Next()
    }
}
```

---

## Summary

**Setup:**
- `gin.Default()` = `gin.New()` + `Logger()` + `Recovery()` middleware
- `gin.New()` = a bare engine, add your own middleware

**Routing:**
- `router.GET/POST/PUT/DELETE(path, handlers...)`, `:param` path params, `*wildcard` catch-alls
- `router.Group(prefix)` for versioning and shared middleware

**Input:**
- `c.Query` / `c.DefaultQuery` / `c.GetQuery` for query strings
- `c.Param` for path parameters
- `c.ShouldBindJSON(&struct)` for validated request bodies (`binding:"required"` etc.)

**Output:**
- `c.JSON`, `c.String`, `c.Status`

**Middleware:**
- A `gin.HandlerFunc` that calls `c.Next()` to continue the chain
- Global (`engine.Use`), per-route, or per-group (`group.Use`)
- `c.Abort...` to stop the chain and write an error response

**Recovery:**
- `gin.Default()`'s `Recovery()` middleware turns a panicking handler into a clean `500`, not a crash

**Testing:**
- `httptest.NewRequest` + `httptest.NewRecorder` + `router.ServeHTTP(w, req)`, same as testing any `http.Handler`

---

## Next Steps

You now understand:
- ✅ Why a framework exists on top of `net/http`, and what Gin specifically adds
- ✅ Setting up `gin.Default()` vs `gin.New()`
- ✅ Routing with path parameters and route groups
- ✅ Reading query params, path params, and validated JSON bodies
- ✅ Writing JSON/string/status-only responses
- ✅ Writing and attaching custom middleware, global and scoped
- ✅ How Gin's Recovery middleware turns panics into clean 500 responses
- ✅ Testing a Gin router with `httptest`, the same way you'd test any `http.Handler`

**Next level:** Level 29 - Database/SQL
- Connecting to a real database from Go
- `database/sql`, connection pools, and drivers
- Queries, scanning rows into structs, and prepared statements
- Wiring a database layer underneath the API you just built with Gin

You're now equipped to build real, production-shaped APIs quickly. Keep going! 🚀
