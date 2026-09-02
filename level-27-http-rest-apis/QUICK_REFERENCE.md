# Level 27: Quick Reference Card

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir myapi && cd myapi
go mod init myapi

# Create main.go
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "net/http"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello, World!")
    })
    mux.HandleFunc("GET /items/{id}", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "item id=%s\n", r.PathValue("id"))
    })

    fmt.Println("listening on :8080")
    http.ListenAndServe(":8080", mux)
}
EOF

# Run
go run main.go
```

---

## 📋 The Handler Interface

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}

type HandlerFunc func(ResponseWriter, *Request) // adapts a func into a Handler

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "hi")
}
var h http.Handler = http.HandlerFunc(handler)
```

---

## 🗺️ ServeMux Routing (Go 1.22+)

```go
mux := http.NewServeMux()

mux.HandleFunc("GET /items", listItems)          // method-specific
mux.HandleFunc("POST /items", createItem)
mux.HandleFunc("GET /items/{id}", getItem)         // wildcard segment
mux.HandleFunc("DELETE /items/{id}", deleteItem)
mux.HandleFunc("GET /files/{path...}", serveFile)  // rest-of-path wildcard

id := r.PathValue("id") // read a matched wildcard
```

Unmatched method on a registered path → automatic `405 Method Not Allowed`.

---

## 📥 Reading Requests

```go
id := r.PathValue("id")              // path wildcard (Go 1.22+)
q := r.URL.Query().Get("q")          // query parameter (always a string, "" if absent)

var req MyStruct
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    http.Error(w, "invalid JSON body", http.StatusBadRequest)
    return
}
```

---

## 📤 Writing Responses

```go
w.Header().Set("Content-Type", "application/json") // 1. headers
w.WriteHeader(http.StatusCreated)                    // 2. status
json.NewEncoder(w).Encode(myStruct)                  // 3. body

http.Error(w, "message", http.StatusBadRequest) // shortcut for a plain-text error body
```

**Order matters:** headers → `WriteHeader` → body. Once bytes are written, the status is locked in.

---

## 🧩 Middleware Pattern

```go
func Logging(logger *log.Logger, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        logger.Printf("%s %s (%s)", r.Method, r.URL.Path, time.Since(start))
    })
}

handler := Logging(logger, mux) // wrap once, compose as many layers as needed
```

---

## 🧪 Testing Handlers

```go
// httptest.NewRecorder - fast, in-process, unit-test ONE handler
req := httptest.NewRequest(http.MethodGet, "/hello", nil)
rec := httptest.NewRecorder()
HelloHandler(rec, req)
// rec.Code, rec.Body.String(), rec.Header()

// httptest.NewServer - real (loopback) round trip, tests the WHOLE stack
server := httptest.NewServer(mux)
defer server.Close()
resp, _ := server.Client().Get(server.URL + "/hello")
// resp.StatusCode, json.NewDecoder(resp.Body).Decode(&v)
```

---

## ⏱️ Context in Handlers

```go
func SlowHandler(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
    defer cancel()

    select {
    case result := <-doWork(ctx):
        writeJSON(w, http.StatusOK, result)
    case <-ctx.Done():
        writeJSON(w, http.StatusGatewayTimeout, map[string]string{"error": ctx.Err().Error()})
    }
}
```

Always derive from `r.Context()`, never `context.Background()`, inside a handler.

---

## 🚨 Consistent JSON Errors

```go
type errorResponse struct {
    Error string `json:"error"`
}

func writeError(w http.ResponseWriter, err error) {
    status := http.StatusInternalServerError
    switch {
    case errors.Is(err, ErrNotFound):
        status = http.StatusNotFound
    case errors.Is(err, ErrInvalidInput):
        status = http.StatusBadRequest
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
}
```

---

## 📊 Status Codes at a Glance

| Code | Constant | Meaning |
|---|---|---|
| 200 | `StatusOK` | Success, body attached |
| 201 | `StatusCreated` | POST created a resource |
| 204 | `StatusNoContent` | Success, no body (e.g. DELETE) |
| 400 | `StatusBadRequest` | Malformed/invalid client input |
| 401 | `StatusUnauthorized` | Not authenticated |
| 404 | `StatusNotFound` | Resource/route doesn't exist |
| 405 | `StatusMethodNotAllowed` | Wrong method for this path (automatic) |
| 429 | `StatusTooManyRequests` | Rate-limited |
| 500 | `StatusInternalServerError` | Unexpected/unmapped error |
| 504 | `StatusGatewayTimeout` | Context deadline exceeded |

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| WriteHeader too late | `w.Write(...)` then `w.WriteHeader(...)` | `w.WriteHeader(...)` then `w.Write(...)` |
| Ignoring decode errors | `json.NewDecoder(r.Body).Decode(&v)` alone | check the returned `error`, respond 400 |
| Goroutine ignores context | `go func(){ time.Sleep(...); ... }()` | `select` on `ctx.Done()` inside the goroutine |
| Leaking response bodies | `resp, _ := http.Get(url)` then read | `defer resp.Body.Close()` right after the error check |
| Magic status numbers | `w.WriteHeader(404)` | `w.WriteHeader(http.StatusNotFound)` |

---

## 🎓 Before Next Level

Can you:
- [ ] Explain the `http.Handler`/`http.HandlerFunc` relationship from memory?
- [ ] Register a method-specific, wildcard `ServeMux` route without looking it up?
- [ ] Read path values, query parameters, and a JSON body?
- [ ] Write a JSON response with the correct `Content-Type` and status code?
- [ ] Write a `func(http.Handler) http.Handler` middleware?
- [ ] Choose between `httptest.NewRecorder` and `httptest.NewServer` correctly?
- [ ] Bound a slow handler with `context.WithTimeout(r.Context(), ...)`?
- [ ] Map application errors to one consistent JSON error shape?

If YES → You're ready for Level 28!

---

## 📚 Next Level

Level 28: Gin Framework
- What a routing framework buys you over hand-rolled `http.ServeMux`
- Gin's context object, route groups, and middleware chains
- Binding and validating JSON with struct tags

You've got HTTP & REST APIs down! 💪
