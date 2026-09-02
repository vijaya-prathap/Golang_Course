# Level 27: HTTP & REST APIs - Complete Guide

## Introduction

Welcome to Level 27! You've mastered testing (Level 26) - you can write `TestXxx` functions, table-driven tests, and you got a preview of `httptest.NewRecorder` for testing HTTP handlers without a real network connection. Now it's time to build the thing those handlers actually serve: real HTTP servers and REST APIs, using nothing but Go's standard library.

This level is marked **Advanced** because it's where everything comes together. Structs and JSON tags (Level 19) become request/response bodies. Errors (Level 16) become HTTP status codes. Goroutines and mutexes (Levels 20, 23) protect shared in-memory state under concurrent requests. Context (Level 24) bounds how long a handler will wait for slow work. And testing (Level 26) is how you prove all of it actually works - not by eyeballing curl output, but with real, repeatable `go test` runs.

Every example and exercise in this level uses **only** `net/http` from the standard library - no routing frameworks, no third-party middleware. `Level 28: Gin Framework` picks up right where this level leaves off, showing you what a framework buys you once you've felt what it's like to build routing, middleware, and JSON handling by hand.

---

## Table of Contents

1. [net/http Fundamentals](#nethttp-fundamentals)
2. [Routing With http.ServeMux](#routing-with-httpservemux)
3. [Reading Requests](#reading-requests)
4. [Writing Responses](#writing-responses)
5. [A Complete Small REST Resource](#a-complete-small-rest-resource)
6. [The Middleware Pattern](#the-middleware-pattern)
7. [Testing HTTP Handlers](#testing-http-handlers)
8. [Using Context in Handlers](#using-context-in-handlers)
9. [Error Responses](#error-responses)
10. [Best Practices](#best-practices)
11. [Common Mistakes](#common-mistakes)

---

## net/http Fundamentals

Everything in Go's HTTP server story rests on one interface:

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

Anything with a `ServeHTTP(w http.ResponseWriter, r *http.Request)` method **is** a handler - a server, a router, a piece of middleware, all of them. There's no base class to extend, no framework object to instantiate. Just one method.

### http.HandlerFunc: Turning a Plain Function Into a Handler

Writing a full type with a `ServeHTTP` method every time would be tedious for simple cases, so the standard library provides an adapter:

```go
type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}
```

`HandlerFunc` is a function type that also satisfies the `Handler` interface - its `ServeHTTP` method just calls itself. This means any ordinary function with the right signature can be used as a `Handler` by wrapping it in `http.HandlerFunc(...)`:

```go
func hello(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello, World!")
}

var h http.Handler = http.HandlerFunc(hello) // now hello satisfies Handler
```

In practice, `mux.HandleFunc(pattern, hello)` does this wrapping for you (see the next section) - but understanding the adapter is what makes middleware (Section 6) make sense: middleware is just a function that takes a `Handler` and returns a new one.

### http.ResponseWriter and *http.Request

Every handler receives exactly two things:

- **`http.ResponseWriter`** - an interface you write the response *to*: `Header()` for response headers, `WriteHeader(statusCode)` for the status line, `Write([]byte)` for the body.
- **`*http.Request`** - a struct you read the incoming request *from*: `r.Method`, `r.URL`, `r.Header`, `r.Body`, and (Go 1.22+) `r.PathValue(name)`.

```go
func PingHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "pong")
}
```

Order matters: set headers, call `WriteHeader`, *then* write the body. Once you call `Write` (or `WriteHeader`), the status line and headers are sent - changing them afterward has no effect (see Common Mistake 1).

### Starting a Real Server

`http.ListenAndServe` blocks forever, serving a `Handler` on a TCP address. For a fully offline, verifiable demo, bind explicitly to loopback with an OS-assigned port instead of a fixed one:

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /ping", PingHandler)

listener, err := net.Listen("tcp", "127.0.0.1:0") // loopback, OS picks a free port
if err != nil {
    panic(err)
}
fmt.Println("listening on", listener.Addr())

server := &http.Server{Handler: mux}
go server.Serve(listener)
```

Real captured output from this exact snippet (the port number varies by run - it's OS-assigned):

```
listening on 127.0.0.1:49162
GET /ping -> 200 pong
```

For tests, you rarely do this by hand - `httptest.NewServer` (Section 7) does exactly this setup for you.

---

## Routing With http.ServeMux

This project's `go.mod` pins `go 1.26.5`, well past the Go 1.22 release that gave `http.ServeMux` **method-specific and wildcard routing** built into the standard library - no router package required. This level uses that real, modern syntax throughout (confirmed by running the code in this file's exercises on this exact toolchain).

### Pattern Syntax

A `ServeMux` pattern has the shape `[METHOD ][HOST]/[PATH]` - every part is optional except `/`:

```go
mux := http.NewServeMux()

mux.HandleFunc("GET /items", listItems)          // only matches GET
mux.HandleFunc("POST /items", createItem)        // only matches POST
mux.HandleFunc("GET /items/{id}", getItem)        // {id} is a wildcard path segment
mux.HandleFunc("DELETE /items/{id}", deleteItem)
```

- A pattern with **no method** matches every method.
- A pattern with method `GET` matches both `GET` and `HEAD` requests.
- `{name}` matches exactly one path segment; `{name...}` matches the rest of the path (must be last).
- Read a wildcard's matched value with `r.PathValue("name")`.

### Verified Behavior: Unmatched Methods Get an Automatic 405

Register `GET /items` and `POST /items` but not `PUT /items`, and the `ServeMux` itself - with no code from you - responds `405 Method Not Allowed` for `PUT /items`. Real captured output from `go run` against an `httptest.NewServer`-backed mux with exactly these routes:

```
server listening on: http://127.0.0.1:65059
GET /items -> 200 "list of items\n"
POST /items -> 200 "created an item\n"
GET /items/42 -> 200 "item id=42\n"
DELETE /items/42 -> 200 "deleted item id=42\n"
PUT /items -> 405 "Method Not Allowed\n"
```

(The port number is OS-assigned and will differ on your machine each run - everything else is deterministic.)

### Longest Pattern Wins

If two patterns could both match a request, `ServeMux` picks the more specific one (the "longest" match), so you can register `/items` and `/items/{id}` without one shadowing the other.

---

## Reading Requests

### Path Values (Go 1.22+)

```go
mux.HandleFunc("GET /books/{id}", func(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id") // the string matched by {id}
    fmt.Fprintf(w, "book id=%s\n", id)
})
```

`r.PathValue` only returns values for wildcards declared in the pattern that matched - there's no reflection or magic, it's a direct lookup by name.

### Query Parameters

```go
func SearchHandler(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query()      // url.Values, a map[string][]string
    q := query.Get("q")         // first value for "q", or "" if absent
    limit := query.Get("limit")
    // ...
}
```

`query.Get(key)` always returns a string (empty if the key is missing) - there's no error to check for a missing key, only for parsing what you got (e.g. `strconv.Atoi` on `limit`).

### JSON Request Bodies

This is exactly the `json.Decoder` from Level 19, applied to `r.Body` (which implements `io.Reader`):

```go
type CreateUserRequest struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid JSON body: "+err.Error(), http.StatusBadRequest)
        return
    }
    // req.Name, req.Age are now populated
}
```

Always check the error from `Decode` (see Common Mistake 2) - a client can send malformed JSON, an empty body, or the wrong shape entirely, and your handler must not trust it blindly.

---

## Writing Responses

### Content-Type and Status Codes

```go
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusCreated) // 201
```

`http.StatusXxx` constants (`http.StatusOK`, `http.StatusNotFound`, `http.StatusBadRequest`, `http.StatusCreated`, `http.StatusNoContent`, `http.StatusInternalServerError`, ...) spell out status codes so handlers never hardcode magic numbers like `404`.

### Encoding JSON Responses

The mirror image of `json.NewDecoder(r.Body).Decode(&v)`:

```go
type CreateUserResponse struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(v)
}
```

`json.NewEncoder(w).Encode(v)` streams JSON straight to the response body - no intermediate `[]byte` needed - and appends a trailing newline. Real captured output from Exercise 4's handler:

```go
writeJSON(w, http.StatusCreated, CreateUserResponse{ID: 1, Name: "Alice", Age: 30})
```

```
Status: 201
Body:   {"id":1,"name":"Alice","age":30}
```

**Order matters here too:** set `Content-Type`, call `WriteHeader`, *then* encode/write the body - once bytes are written, the status code is locked in.

---

## A Complete Small REST Resource

Putting routing, request reading, and response writing together, here is a full in-memory "books" resource: `GET /books` (list), `POST /books` (create), `GET /books/{id}` (get one), `DELETE /books/{id}` (delete). The store uses a `sync.Mutex` (Level 23) because multiple requests can arrive concurrently.

```go
type Book struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
}

type BookStore struct {
    mu     sync.Mutex
    nextID int
    books  map[int]Book
}

func NewBookStore() *BookStore {
    return &BookStore{nextID: 1, books: make(map[int]Book)}
}

func (s *BookStore) List() []Book {
    s.mu.Lock()
    defer s.mu.Unlock()
    result := make([]Book, 0, len(s.books))
    for id := 1; id < s.nextID; id++ {
        if b, ok := s.books[id]; ok {
            result = append(result, b)
        }
    }
    return result
}

func (s *BookStore) Create(title, author string) Book {
    s.mu.Lock()
    defer s.mu.Unlock()
    b := Book{ID: s.nextID, Title: title, Author: author}
    s.books[b.ID] = b
    s.nextID++
    return b
}

func (s *BookStore) Get(id int) (Book, bool) {
    s.mu.Lock()
    defer s.mu.Unlock()
    b, ok := s.books[id]
    return b, ok
}

func (s *BookStore) Delete(id int) bool {
    s.mu.Lock()
    defer s.mu.Unlock()
    if _, ok := s.books[id]; !ok {
        return false
    }
    delete(s.books, id)
    return true
}
```

```go
func newBookMux(store *BookStore) *http.ServeMux {
    mux := http.NewServeMux()

    mux.HandleFunc("GET /books", func(w http.ResponseWriter, r *http.Request) {
        writeJSON(w, http.StatusOK, store.List())
    })

    mux.HandleFunc("POST /books", func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Title  string `json:"title"`
            Author string `json:"author"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
            return
        }
        if req.Title == "" {
            writeJSON(w, http.StatusBadRequest, map[string]string{"error": "title is required"})
            return
        }
        writeJSON(w, http.StatusCreated, store.Create(req.Title, req.Author))
    })

    mux.HandleFunc("GET /books/{id}", func(w http.ResponseWriter, r *http.Request) {
        id, err := strconv.Atoi(r.PathValue("id"))
        if err != nil {
            writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid book id"})
            return
        }
        book, ok := store.Get(id)
        if !ok {
            writeJSON(w, http.StatusNotFound, map[string]string{"error": "book not found"})
            return
        }
        writeJSON(w, http.StatusOK, book)
    })

    mux.HandleFunc("DELETE /books/{id}", func(w http.ResponseWriter, r *http.Request) {
        id, err := strconv.Atoi(r.PathValue("id"))
        if err != nil {
            writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid book id"})
            return
        }
        if !store.Delete(id) {
            writeJSON(w, http.StatusNotFound, map[string]string{"error": "book not found"})
            return
        }
        w.WriteHeader(http.StatusNoContent)
    })

    return mux
}
```

Real captured output, driving this exact resource end-to-end over `httptest.NewServer` (a real loopback TCP round trip):

```
GET /books -> 200 []
POST /books {"title":"The Go Programming Language","author":"Donovan and Kernighan"} -> 201 {"id":1,"title":"The Go Programming Language","author":"Donovan and Kernighan"}
POST /books {"title":"The Pragmatic Programmer","author":"Hunt and Thomas"} -> 201 {"id":2,"title":"The Pragmatic Programmer","author":"Hunt and Thomas"}
GET /books -> 200 [{"id":1,...},{"id":2,...}]
GET /books/1 -> 200 {"id":1,"title":"The Go Programming Language","author":"Donovan and Kernighan"}
GET /books/99 -> 404 {"error":"book not found"}
DELETE /books/1 -> 204 ""
GET /books -> 200 [{"id":2,"title":"The Pragmatic Programmer","author":"Hunt and Thomas"}]
DELETE /books/1 -> 404 {"error":"book not found"}
```

Notice each handler stays thin: decode/parse, call one store method, write one response. All the actual logic (locking, ID assignment, lookup) lives in `BookStore` - a plain Go type this API happens to expose over HTTP (see Best Practice 4).

---

## The Middleware Pattern

Middleware wraps a `Handler` with a new `Handler` that adds behavior before and/or after calling the original:

```go
func(next http.Handler) http.Handler
```

Because both the input and output are `http.Handler`, middleware composes: `Logging(RequestID(RateLimit(mux)))` wraps three layers around one mux, and each layer only knows about the one directly inside it.

### A Logging Middleware

```go
type statusRecorder struct {
    http.ResponseWriter
    status int
}

func (r *statusRecorder) WriteHeader(code int) {
    r.status = code
    r.ResponseWriter.WriteHeader(code)
}

func Logging(logger *log.Logger, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

        next.ServeHTTP(rec, r) // call the wrapped handler

        logger.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, rec.status, time.Since(start).Round(time.Microsecond))
    })
}
```

`statusRecorder` embeds `http.ResponseWriter` and overrides only `WriteHeader`, so it can observe the status code the inner handler chose without changing any other behavior - a small, common trick for HTTP middleware.

Real captured output (a `log.Logger` writing into a `bytes.Buffer` instead of stdout, so the exact text can be asserted on in a test):

```
GET /brew -> 418 (1µs)
```

(The duration is machine-dependent; the method, path, and status are exact and deterministic.)

---

## Testing HTTP Handlers

Level 26 previewed this; here's the full picture, using both tools built into `net/http/httptest` for different jobs.

### httptest.NewRecorder: Unit-Test a Handler Directly

No network, no goroutine, no server - just call the handler function with a fake request and a fake response writer, then inspect what it wrote:

```go
func TestHelloHandler(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/hello", nil)
    rec := httptest.NewRecorder()

    HelloHandler(rec, req)

    if rec.Code != http.StatusOK {
        t.Errorf("status = %d; want %d", rec.Code, http.StatusOK)
    }
    if got := rec.Body.String(); got != "Hello, World!\n" {
        t.Errorf("body = %q; want %q", got, "Hello, World!\n")
    }
}
```

Real captured output:

```
=== RUN   TestHelloHandler
--- PASS: TestHelloHandler (0.00s)
PASS
ok  	level27.example/exercise1	0.826s
```

Use this for the vast majority of handler tests - it's fast (no real I/O) and works for any single `http.HandlerFunc`, including one wrapped in middleware (just call the wrapped handler's `ServeHTTP`).

### httptest.NewServer: A Full Round Trip Over Real HTTP

When you need to prove the *entire* stack works together - real TCP, a real `*http.Client`, routing through a real `ServeMux` - `httptest.NewServer` starts an actual server on `127.0.0.1` with an OS-assigned port, entirely on loopback:

```go
func TestBookAPIRoundTrip(t *testing.T) {
    store := NewBookStore()
    server := httptest.NewServer(newBookMux(store))
    defer server.Close()

    client := server.Client()

    resp, err := client.Get(server.URL + "/books")
    if err != nil {
        t.Fatalf("GET /books: %v", err)
    }
    defer resp.Body.Close()
    // ... assert on resp.StatusCode, decode resp.Body ...
}
```

Real captured output, a table of subtests covering list/create/get/404:

```
=== RUN   TestBookAPIRoundTrip
=== RUN   TestBookAPIRoundTrip/list_starts_empty
=== RUN   TestBookAPIRoundTrip/create_a_book
=== RUN   TestBookAPIRoundTrip/get_the_created_book_by_id
=== RUN   TestBookAPIRoundTrip/get_a_nonexistent_book_returns_404
--- PASS: TestBookAPIRoundTrip (0.00s)
    --- PASS: TestBookAPIRoundTrip/list_starts_empty (0.00s)
    --- PASS: TestBookAPIRoundTrip/create_a_book (0.00s)
    --- PASS: TestBookAPIRoundTrip/get_the_created_book_by_id (0.00s)
    --- PASS: TestBookAPIRoundTrip/get_a_nonexistent_book_returns_404 (0.00s)
PASS
ok  	level27.example/exercise7	0.576s
```

`server.Close()` (always `defer`red) shuts the listener down cleanly at the end of the test - since it's bound to loopback with an OS-assigned port, it never touches or requires real network access.

| | `httptest.NewRecorder` | `httptest.NewServer` |
|---|---|---|
| What runs | Your handler function directly, in-process | A real `net.Listener` + goroutine serving your `Handler` |
| Network involved | None - it's a fake `ResponseWriter` | Real TCP, on `127.0.0.1` only |
| Speed | Very fast | Slightly slower (real I/O, still local) |
| Good for | Testing one handler/middleware chain in isolation | Proving routing + the whole stack works end-to-end |
| Client | You call the handler's `ServeHTTP` yourself | A real `*http.Client` (`server.Client()`) |

---

## Using Context in Handlers

Every `*http.Request` carries a `context.Context` (Level 24), accessible via `r.Context()`. The Go HTTP server cancels it automatically if the client disconnects or the request's connection closes - and you can derive a *tighter* deadline from it for a specific slow operation:

```go
func slowOperation(ctx context.Context, delay time.Duration) error {
    select {
    case <-time.After(delay):
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

func ReportHandler(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 50*time.Millisecond)
    defer cancel()

    err := slowOperation(ctx, 200*time.Millisecond) // slower than the timeout, on purpose

    w.Header().Set("Content-Type", "application/json")
    if err != nil {
        w.WriteHeader(http.StatusGatewayTimeout) // 504
        json.NewEncoder(w).Encode(map[string]string{"error": "report generation timed out", "cause": err.Error()})
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "report ready"})
}
```

Real captured output from a table-driven test covering both the slow (times out) and fast (succeeds) cases:

```
=== RUN   TestReportHandler_TimesOut
--- PASS: TestReportHandler_TimesOut (0.05s)
=== RUN   TestFastReportHandler_Succeeds
--- PASS: TestFastReportHandler_Succeeds (0.01s)
PASS
ok  	level27.example/exercise8	0.646s
```

`ctx.Err()` after a timeout returns `context.DeadlineExceeded` - checking it (or using `errors.Is`) lets a handler distinguish "the client canceled" from "we gave up waiting," and choose an accurate status code (`499`-style client-closed vs. `504 Gateway Timeout`) instead of a generic `500`.

Deriving `context.WithTimeout(r.Context(), ...)` - rather than `context.WithTimeout(context.Background(), ...)` - matters: if the *client* disconnects, `r.Context()` is canceled too, and that cancellation propagates into your derived context immediately, stopping the slow operation even before your own timeout would have.

---

## Error Responses

A REST API should return errors in **one consistent JSON shape**, no matter which internal error caused them - callers should never have to guess the response format based on which endpoint failed.

```go
type errorResponse struct {
    Error string `json:"error"`
}

var (
    ErrNotFound     = errors.New("resource not found")
    ErrInvalidInput = errors.New("invalid input")
    ErrUnauthorized = errors.New("unauthorized")
)

// writeError maps an application error (Level 16) to an HTTP status code,
// using errors.Is so wrapped errors still match their sentinel.
func writeError(w http.ResponseWriter, err error) {
    status := http.StatusInternalServerError
    switch {
    case errors.Is(err, ErrNotFound):
        status = http.StatusNotFound
    case errors.Is(err, ErrInvalidInput):
        status = http.StatusBadRequest
    case errors.Is(err, ErrUnauthorized):
        status = http.StatusUnauthorized
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
}
```

A wrapped error (Level 16's `Unwrap() error`) still maps correctly, because `errors.Is` walks the chain:

```go
type ValidationError struct {
    Field, Message string
}

func (e *ValidationError) Error() string { return "validation failed on " + e.Field + ": " + e.Message }
func (e *ValidationError) Unwrap() error { return ErrInvalidInput } // still matches ErrInvalidInput
```

Real captured output, one table-driven test covering a wrapped validation error, two sentinel errors, and one *unmapped* error (correctly falling through to `500`):

```
=== RUN   TestUserHandler_ErrorMapping
=== RUN   TestUserHandler_ErrorMapping/empty_id_->_validation_error_->_400
=== RUN   TestUserHandler_ErrorMapping/404_id_->_not_found_->_404
=== RUN   TestUserHandler_ErrorMapping/401_id_->_unauthorized_->_401
=== RUN   TestUserHandler_ErrorMapping/boom_id_->_unmapped_error_->_500
=== RUN   TestUserHandler_ErrorMapping/valid_id_->_200
--- PASS: TestUserHandler_ErrorMapping (0.00s)
    --- PASS: TestUserHandler_ErrorMapping/empty_id_->_validation_error_->_400 (0.00s)
    --- PASS: TestUserHandler_ErrorMapping/404_id_->_not_found_->_404 (0.00s)
    --- PASS: TestUserHandler_ErrorMapping/401_id_->_unauthorized_->_401 (0.00s)
    --- PASS: TestUserHandler_ErrorMapping/boom_id_->_unmapped_error_->_500 (0.00s)
    --- PASS: TestUserHandler_ErrorMapping/valid_id_->_200 (0.00s)
PASS
ok  	level27.example/exercise9	0.620s
```

An unmapped error deliberately falling through to `500 Internal Server Error` is the *correct* default - it means "something the API author didn't anticipate," which is exactly what a 500 should signal, without leaking stack traces or internal detail into `err.Error()` in a real production system (this level's examples keep messages simple for learning; Level 33: Logging and Level 38: Production Debugging cover safely separating internal detail from client-facing messages).

---

## Best Practices

### 1. Always Set Content-Type

Every response - success or error - should declare what it is. Clients (and tests) that parse the body depend on it:

```go
// ✅ Good
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
json.NewEncoder(w).Encode(v)

// ❌ Risky - client has to guess/sniff the content type
w.WriteHeader(http.StatusOK)
json.NewEncoder(w).Encode(v)
```

### 2. Validate and Sanitize Input

Never trust `r.Body`, `r.URL.Query()`, or `r.PathValue()` blindly - check decode errors, required fields, and value ranges before acting on them (Section 9's error-mapping pattern is how you report the rejection consistently).

### 3. Use Context for Cancellation

Derive timeouts from `r.Context()`, not `context.Background()`, so a slow handler both respects its own deadline *and* stops early if the client goes away:

```go
// ✅ Good - inherits client disconnection AND adds a tighter deadline
ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
defer cancel()

// ❌ Ignores the request entirely - keeps working even after the client is gone
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()
```

### 4. Keep Handlers Thin

A handler's job is translation: HTTP in, a plain Go function call, HTTP out. Push actual logic into separate functions/types (like `BookStore` in Section 5) that don't know anything about `http.ResponseWriter` or `*http.Request` - they're then trivially testable and reusable outside HTTP entirely.

```go
// ✅ Good - handler is a thin adapter
mux.HandleFunc("POST /books", func(w http.ResponseWriter, r *http.Request) {
    var req createBookRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil { /* ... */ }
    writeJSON(w, http.StatusCreated, store.Create(req.Title, req.Author))
})

// ❌ Business logic tangled directly into the handler makes it hard to
// test or reuse without spinning up an HTTP request
```

### 5. Return Consistent Error Shapes

One `errorResponse{Error string}` (or similar) shape for every failure path in the whole API - never a different JSON structure per endpoint (Section 9).

---

## Common Mistakes

### Mistake 1: Calling w.Write Before w.WriteHeader (or Not At All)

```go
// ❌ WRONG - Write() before WriteHeader() implicitly sends 200 OK first;
// this WriteHeader call is silently ignored and Go logs a warning
w.Write([]byte("created"))
w.WriteHeader(http.StatusCreated) // too late - status is already 200

// ✅ RIGHT - set the status BEFORE writing any body bytes
w.WriteHeader(http.StatusCreated)
w.Write([]byte("created"))
```

If you never call `WriteHeader` at all, the first `Write` implicitly sends `200 OK` - which is fine for simple success cases, but means you cannot "change your mind" about the status code after any bytes have gone out.

### Mistake 2: Not Checking the JSON Decode Error

```go
// ❌ WRONG - ignores a malformed/empty body; req keeps its zero values
// and the handler proceeds as if everything was fine
json.NewDecoder(r.Body).Decode(&req)

// ✅ RIGHT
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    http.Error(w, "invalid JSON body", http.StatusBadRequest)
    return
}
```

### Mistake 3: Leaking a Goroutine That Ignores Context Cancellation

```go
// ❌ WRONG - if ctx times out, ReportHandler returns, but this goroutine
// keeps running in the background forever, doing wasted work with no
// way for anyone to stop it
func slowWork(ctx context.Context) <-chan string {
    result := make(chan string)
    go func() {
        time.Sleep(10 * time.Second) // ignores ctx entirely
        result <- "done"
    }()
    return result
}

// ✅ RIGHT - the goroutine itself selects on ctx.Done() (or the work it
// calls does), so it actually stops instead of running to completion
// after nobody is listening anymore
func slowWork(ctx context.Context) <-chan string {
    result := make(chan string, 1)
    go func() {
        select {
        case <-time.After(10 * time.Second):
            result <- "done"
        case <-ctx.Done():
            return // stop - no one is waiting for `result` anymore
        }
    }()
    return result
}
```

### Mistake 4: Not Closing (or Draining) r.Body When It Matters

For handlers that only `Decode()` the body once and return, the Go HTTP server manages `r.Body`'s lifecycle for you - but if you make an **outgoing** HTTP request from within a handler (proxying, calling another service) via `http.Client`, its response body must be closed explicitly:

```go
// ❌ WRONG - leaks the underlying connection back to the pool
resp, _ := http.Get(someURL)
data, _ := io.ReadAll(resp.Body)

// ✅ RIGHT
resp, err := http.Get(someURL)
if err != nil {
    // handle err
}
defer resp.Body.Close()
data, _ := io.ReadAll(resp.Body)
```

### Mistake 5: Hardcoding Magic Status Code Numbers

```go
// ❌ WRONG - what is 404? readable, but not self-documenting, and easy
// to typo into a wrong number
w.WriteHeader(404)

// ✅ RIGHT
w.WriteHeader(http.StatusNotFound)
```

---

## Summary

**The Core Interface:**
- `http.Handler` - one method, `ServeHTTP(ResponseWriter, *Request)`
- `http.HandlerFunc` - adapts a plain function into a `Handler`

**Routing (Go 1.22+, confirmed on this project's Go 1.26.5 toolchain):**
- `mux.HandleFunc("METHOD /path/{wildcard}", handler)` - method-specific, wildcard-capable, built into `net/http`
- `r.PathValue("wildcard")` reads a matched segment
- Unmatched methods on a registered path get an automatic `405`

**Reading Requests:**
- `r.PathValue(name)`, `r.URL.Query().Get(key)`, `json.NewDecoder(r.Body).Decode(&v)`

**Writing Responses:**
- `w.Header().Set(...)` → `w.WriteHeader(status)` → `w.Write(...)` / `json.NewEncoder(w).Encode(v)`, strictly in that order

**Middleware:**
- `func(http.Handler) http.Handler` - wraps cross-cutting behavior (logging, auth, rate limiting) around any handler, and composes

**Testing:**
- `httptest.NewRecorder` - fast, in-process, unit-test one handler
- `httptest.NewServer` - real (loopback-only) TCP round trip, proves the whole stack

**Context:**
- `r.Context()` is canceled automatically on client disconnect; derive `context.WithTimeout` from it, never from `context.Background()`, inside a handler

**Errors:**
- One consistent JSON error shape; map application errors (Level 16) to status codes with `errors.Is`

---

## Next Steps

You now understand:
- ✅ The `http.Handler`/`http.HandlerFunc` relationship
- ✅ Modern `http.ServeMux` routing with methods and wildcards
- ✅ Reading path values, query parameters, and JSON bodies
- ✅ Writing JSON responses with correct status codes
- ✅ Building a real in-memory REST resource end to end
- ✅ The middleware pattern, verified with real logging output
- ✅ Testing handlers with both `httptest.NewRecorder` and `httptest.NewServer`
- ✅ Bounding slow work with `context.WithTimeout` inside a handler
- ✅ Consistent, error-mapped JSON error responses

**Next level:** Level 28 - Gin Framework
- What a routing framework buys you over hand-rolled `http.ServeMux`
- Gin's context object, route groups, and middleware chains
- Binding and validating JSON with struct tags
- Building the same kind of REST API faster, with less boilerplate

You just built REST APIs the way the standard library actually works underneath every framework - keep going! 🚀
