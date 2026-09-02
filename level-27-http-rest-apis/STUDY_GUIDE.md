# Level 27: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: Fundamentals, Routing, and Data
```
Day 1:  http.Handler / http.HandlerFunc, ResponseWriter, Request
Day 2:  http.ServeMux routing - method patterns and {wildcards}
Day 3:  Reading requests - path values, query params, JSON bodies
Day 4:  Writing responses - Content-Type, JSON encoding, status codes
Day 5:  Building the in-memory CRUD resource (list/create/get/delete)
Day 6:  The middleware pattern - func(http.Handler) http.Handler
Day 7:  Review and consolidate Sections 1-6
```

### Week 2: Testing, Context, Errors, and Practice
```
Day 1:  httptest.NewRecorder vs httptest.NewServer
Day 2:  Exercises 1-4 (hello handler, routing, query params, JSON)
Day 3:  Exercises 5-7 (CRUD resource, middleware, round-trip test)
Day 4:  context.WithTimeout inside a handler
Day 5:  Exercises 8-9 (context-aware handler, JSON error mapping)
Day 6:  Exercise 10 (comprehensive task manager API)
Day 7:  Bonus challenges & review
```

---

## 🎯 The http.Handler / http.HandlerFunc Relationship

```
                    ┌─────────────────────────────┐
                    │   type Handler interface {  │
                    │       ServeHTTP(w, r)       │
                    │   }                          │
                    └──────────────┬───────────────┘
                                   │ satisfied by
              ┌────────────────────┼────────────────────┐
              │                    │                    │
   ┌──────────▼─────────┐ ┌────────▼─────────┐ ┌────────▼─────────┐
   │   *http.ServeMux    │ │  http.HandlerFunc │ │  your middleware  │
   │  (a router IS a     │ │  (adapts a plain  │ │  (wraps a Handler  │
   │   Handler)           │ │   func into one)  │ │   and returns one) │
   └──────────────────────┘ └───────────────────┘ └─────────────────┘
```

```go
type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) { f(w, r) }
```

**The trick:** `HandlerFunc` is a function type with a `ServeHTTP` method that just calls itself. That's the entire adapter - it's what lets `mux.HandleFunc(pattern, myPlainFunc)` accept an ordinary function anywhere a `Handler` is required.

---

## 🔄 Request Lifecycle Diagram

```
 client                 net/http server              your code
   │                          │                            │
   │  HTTP request over TCP   │                            │
   ├─────────────────────────►│                            │
   │                          │  parse into *http.Request  │
   │                          │────────────┐               │
   │                          │◄───────────┘               │
   │                          │                            │
   │                          │  ServeMux.ServeHTTP(w, r)   │
   │                          │  → match pattern, extract   │
   │                          │    path values               │
   │                          ├───────────────────────────►│
   │                          │        Middleware N          │
   │                          │◄─ wraps ──┐                  │
   │                          │           │  Middleware 1     │
   │                          │           │◄─ wraps ──┐       │
   │                          │           │           │        │
   │                          │           │      Handler        │
   │                          │           │  (reads r, writes w)│
   │                          │           │           │        │
   │                          │           └──── after ─┘        │
   │                          │◄──── after ───┘                  │
   │  HTTP response over TCP  │                                  │
   │◄─────────────────────────┤                                  │
```

Each middleware layer runs code **before** calling `next.ServeHTTP` (e.g. start a timer, generate a request ID) and optionally **after** it returns (e.g. log the final status code) - the handler itself sits at the center, innermost.

---

## 🗺️ ServeMux Pattern Syntax (Go 1.22+, confirmed on this project's Go 1.26.5)

| Pattern | Matches |
|---|---|
| `"/"` | Any path, any method (catch-all) |
| `"/items"` | Exact path `/items`, any method |
| `"GET /items"` | Path `/items`, only `GET` (and `HEAD`) |
| `"POST /items"` | Path `/items`, only `POST` |
| `"GET /items/{id}"` | `/items/<anything-one-segment>`, captured as `r.PathValue("id")` |
| `"GET /files/{path...}"` | `/files/<rest-of-path-including-slashes>`, captured as `r.PathValue("path")` |
| `"example.com/"` | Any path, but only when the `Host` header is `example.com` |
| `"/items/{$}"` | Matches `/items/` exactly - `{$}` means "end of path here" |

**Verified, real behavior:** registering `GET /items` and `POST /items` but sending `PUT /items` returns `405 Method Not Allowed` automatically - `ServeMux` does this without any code from you. Longest/most-specific pattern wins when more than one could match.

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /books", listBooks)
mux.HandleFunc("GET /books/{id}", getBook)     // r.PathValue("id")
mux.HandleFunc("POST /books", createBook)
mux.HandleFunc("DELETE /books/{id}", deleteBook)
```

---

## 🧪 httptest.NewRecorder vs httptest.NewServer

| | `httptest.NewRecorder` | `httptest.NewServer` |
|---|---|---|
| **What it is** | A fake `http.ResponseWriter` that captures what's written | A real `net.Listener` on `127.0.0.1` + an OS-assigned port |
| **Network involved** | None - fully in-process | Real TCP, loopback-only |
| **How you call the handler** | `handler.ServeHTTP(rec, req)` directly | `server.Client().Get(server.URL + path)` |
| **Request comes from** | `httptest.NewRequest(method, target, body)` | A real `*http.Client` |
| **Speed** | Fastest - no I/O | Fast, but real I/O (still local) |
| **Best for** | Unit-testing one handler or middleware chain in isolation | Proving the full stack (routing + middleware + handler) works end to end |
| **Cleanup** | None needed | `defer server.Close()` |

```go
// NewRecorder: fast, in-process
req := httptest.NewRequest(http.MethodGet, "/hello", nil)
rec := httptest.NewRecorder()
HelloHandler(rec, req)
// assert on rec.Code, rec.Body, rec.Header()

// NewServer: real (loopback) round trip
server := httptest.NewServer(mux)
defer server.Close()
resp, err := server.Client().Get(server.URL + "/hello")
// assert on resp.StatusCode, decode resp.Body
```

**Rule of thumb:** reach for `NewRecorder` by default - it's what you'll use for nearly every handler test. Reach for `NewServer` when you specifically need to prove routing and the whole `ServeMux` + middleware chain behave correctly together, or when testing code that itself makes real HTTP calls (like an API client).

---

## 📊 Status Code Cheat Table

| Code | Constant | Use it when... |
|---|---|---|
| 200 | `http.StatusOK` | The request succeeded and there's a body to return |
| 201 | `http.StatusCreated` | A `POST` successfully created a new resource |
| 204 | `http.StatusNoContent` | The request succeeded but there's nothing to return (e.g. a successful `DELETE`) |
| 400 | `http.StatusBadRequest` | The client sent malformed JSON, a missing required field, or an invalid value |
| 401 | `http.StatusUnauthorized` | The request has no (or invalid) authentication |
| 403 | `http.StatusForbidden` | The client is authenticated but not allowed to do this |
| 404 | `http.StatusNotFound` | The requested resource (or route) doesn't exist |
| 405 | `http.StatusMethodNotAllowed` | The path exists but not for this HTTP method (`ServeMux` sends this automatically) |
| 409 | `http.StatusConflict` | The request conflicts with the resource's current state (e.g. duplicate create) |
| 429 | `http.StatusTooManyRequests` | The client is being rate-limited |
| 500 | `http.StatusInternalServerError` | An unexpected/unmapped error occurred - the correct default fallback |
| 504 | `http.StatusGatewayTimeout` | A context deadline was exceeded waiting on slow work |

---

## 🧩 The Middleware Pattern at a Glance

```go
type Middleware func(http.Handler) http.Handler

func Logging(logger *log.Logger, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // ... before ...
        next.ServeHTTP(w, r)
        // ... after ...
    })
}

// Composing several layers - innermost (mux) to outermost (Logging):
handler := Logging(logger, RequestID(RateLimit(mux)))
```

Because every layer's input and output are both `http.Handler`, layers stack in any order you choose, and each one only needs to know about the single handler it wraps - not the whole chain.

---

## 🚨 Common Mistakes

### Mistake 1: WriteHeader After Write (or Never Calling It)

```go
// ❌ WRONG - Write() implicitly sends 200 OK; this WriteHeader is too late
w.Write([]byte("created"))
w.WriteHeader(http.StatusCreated)

// ✅ RIGHT - status BEFORE body
w.WriteHeader(http.StatusCreated)
w.Write([]byte("created"))
```

### Mistake 2: Ignoring the JSON Decode Error

```go
// ❌ WRONG - malformed body silently leaves req at its zero value
json.NewDecoder(r.Body).Decode(&req)

// ✅ RIGHT
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    http.Error(w, "invalid JSON body", http.StatusBadRequest)
    return
}
```

### Mistake 3: A Goroutine That Ignores Context Cancellation

```go
// ❌ WRONG - keeps running even after the handler (and client) gave up
go func() {
    time.Sleep(10 * time.Second)
    result <- "done"
}()

// ✅ RIGHT - selects on ctx.Done() too, so it actually stops
go func() {
    select {
    case <-time.After(10 * time.Second):
        result <- "done"
    case <-ctx.Done():
        return
    }
}()
```

### Mistake 4: Not Closing an Outgoing Response Body

```go
// ❌ WRONG - leaks the connection back to the pool
resp, _ := http.Get(someURL)
data, _ := io.ReadAll(resp.Body)

// ✅ RIGHT
resp, err := http.Get(someURL)
if err != nil { /* handle */ }
defer resp.Body.Close()
```

---

## 📈 Progression Summary

### Understanding Level 27

Level 27 teaches how everything from Levels 16-26 comes together behind an HTTP boundary:

1. **The interface** - `http.Handler`/`http.HandlerFunc`, the whole foundation
2. **Routing** - `http.ServeMux`'s real, modern method+wildcard syntax
3. **I/O** - reading path values/query params/JSON in, writing JSON+status codes out
4. **A real resource** - CRUD, backed by a mutex-protected store (Level 23)
5. **Middleware** - cross-cutting behavior via `func(http.Handler) http.Handler`
6. **Testing** - `httptest.NewRecorder` (unit) and `httptest.NewServer` (full stack)
7. **Context** - bounding slow work with `context.WithTimeout` (Level 24)
8. **Errors** - mapping Go errors (Level 16) to consistent HTTP responses

### Prerequisites for Level 28

Before moving to Level 28 (Gin Framework), you need:

- ✅ Comfortable with `http.Handler`, `http.HandlerFunc`, and `http.ServeMux` routing
- ✅ Can read path values, query parameters, and JSON bodies from a request
- ✅ Can write JSON responses with correct status codes
- ✅ Can build and test a small CRUD REST resource
- ✅ Understand the middleware pattern and can write one
- ✅ Comfortable testing handlers with both `httptest.NewRecorder` and `httptest.NewServer`
- ✅ Can use `context.WithTimeout` inside a handler
- ✅ Can map application errors to consistent JSON error responses

### Ready for Level 28?

Level 28 introduces the Gin framework - everything you just hand-built (routing, middleware chains, JSON binding) gets a purpose-built API on top of the same underlying `net/http` concepts:
- Gin's `*gin.Context` (a richer `ResponseWriter`+`Request` bundle)
- Route groups and chained middleware
- Struct-tag-based binding and validation
- Building the same REST API with far less boilerplate

---

## ✅ Checklist Before Level 28

- [ ] Can explain why `http.HandlerFunc` lets a plain function satisfy `http.Handler`
- [ ] Can register method-specific and wildcard `ServeMux` patterns from memory
- [ ] Can read `r.PathValue`, `r.URL.Query()`, and decode a JSON body
- [ ] Always sets `Content-Type` and calls `WriteHeader` before writing a body
- [ ] Can write a `func(http.Handler) http.Handler` middleware
- [ ] Knows when to reach for `httptest.NewRecorder` vs `httptest.NewServer`
- [ ] Can derive `context.WithTimeout` from `r.Context()` and check `ctx.Err()`
- [ ] Returns one consistent JSON error shape, mapped from Go errors
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The One Interface
Everything in `net/http` - servers, routers, middleware, handlers - is just something with a `ServeHTTP(w, r)` method.

### The Modern Router
`http.ServeMux` on Go 1.22+ handles method-specific routing and path wildcards natively - no third-party router required for most APIs.

### The Testing Split
`httptest.NewRecorder` for fast, in-process unit tests of one handler; `httptest.NewServer` for real (loopback) round trips proving the whole stack.

### The Boundary
A handler's job is translation - HTTP in, plain Go function calls, HTTP out. Keep it thin; push logic into testable, HTTP-unaware types.

---

## 📚 Next Level

Level 28: Gin Framework
- What a routing framework buys you over hand-rolled `http.ServeMux`
- Gin's context object, route groups, and middleware chains
- Binding and validating JSON with struct tags
- Building the same kind of REST API faster, with less boilerplate

You've built REST APIs the hard way - the way that makes every framework afterward make sense! 🚀
