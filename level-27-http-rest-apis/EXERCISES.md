# Level 27: HTTP & REST APIs - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

**A note on this level's format:** most exercises here follow Level 26's pattern - a code file plus a `_test.go` file, verified with `go test -v`, since that's how you actually prove an HTTP handler works without touching a real network. A couple of exercises use `go run main.go` instead, when the point is to watch a real (loopback-only) request/response round trip happen.

---

## Exercise 1: A "Hello, World" Handler Tested With httptest.NewRecorder

**Objective:** Write the smallest possible HTTP handler and unit-test it without a network connection

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level27-exercise1
cd ~/projects/level27-exercise1
go mod init level27.example/exercise1
```

2. Create `handler.go`:

```bash
cat > handler.go << 'EOF'
package api

import (
    "fmt"
    "net/http"
)

// HelloHandler writes a simple plain-text greeting. It's the smallest
// possible http.HandlerFunc: read nothing from the request, write a
// status code and a body.
func HelloHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "Hello, World!")
}
EOF
```

3. Create `handler_test.go`:

```bash
cat > handler_test.go << 'EOF'
package api

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHelloHandler(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/hello", nil)
    rec := httptest.NewRecorder()

    HelloHandler(rec, req)

    if rec.Code != http.StatusOK {
        t.Errorf("status = %d; want %d", rec.Code, http.StatusOK)
    }
    wantBody := "Hello, World!\n"
    if got := rec.Body.String(); got != wantBody {
        t.Errorf("body = %q; want %q", got, wantBody)
    }
    wantType := "text/plain; charset=utf-8"
    if got := rec.Header().Get("Content-Type"); got != wantType {
        t.Errorf("Content-Type = %q; want %q", got, wantType)
    }
}
EOF
```

4. Run the tests:

```bash
go test -v ./...
```

**Expected Output:**

```
=== RUN   TestHelloHandler
--- PASS: TestHelloHandler (0.00s)
PASS
ok  	level27.example/exercise1	0.826s
```

**Learning Objectives:**
- ✅ Write a plain function matching the `func(http.ResponseWriter, *http.Request)` signature
- ✅ Build a fake request with `httptest.NewRequest` and capture output with `httptest.NewRecorder`
- ✅ Assert on `rec.Code`, `rec.Body`, and response headers with zero real network I/O

---

## Exercise 2: ServeMux Routing With Method and Wildcard Patterns

**Objective:** Register method-specific routes with path wildcards using the real syntax this Go version supports, and watch real routing decisions happen over a loopback server

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level27-exercise2
cd ~/projects/level27-exercise2
go mod init level27.example/exercise2
```

2. Create `main.go`. This project's `go.mod` targets `go 1.26.5`, well past Go 1.22, so `http.ServeMux` supports method-specific and wildcard patterns natively - no router library needed:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "io"
    "net/http"
    "net/http/httptest"
)

func itemsHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "list of items")
}

func itemByIDHandler(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    fmt.Fprintf(w, "item id=%s\n", id)
}

func createItemHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "created an item")
}

func deleteItemHandler(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    fmt.Fprintf(w, "deleted item id=%s\n", id)
}

func newMux() *http.ServeMux {
    mux := http.NewServeMux()
    // Method-specific and wildcard patterns, native to net/http since Go 1.22.
    mux.HandleFunc("GET /items", itemsHandler)
    mux.HandleFunc("POST /items", createItemHandler)
    mux.HandleFunc("GET /items/{id}", itemByIDHandler)
    mux.HandleFunc("DELETE /items/{id}", deleteItemHandler)
    return mux
}

func main() {
    // httptest.NewServer binds to 127.0.0.1 on an OS-assigned port -
    // a real TCP round trip, entirely on loopback, no outbound network.
    server := httptest.NewServer(newMux())
    defer server.Close()

    fmt.Println("server listening on:", server.URL)

    requests := []struct {
        method string
        path   string
    }{
        {http.MethodGet, "/items"},
        {http.MethodPost, "/items"},
        {http.MethodGet, "/items/42"},
        {http.MethodDelete, "/items/42"},
        {http.MethodPut, "/items"}, // no PUT /items pattern registered
    }

    for _, req := range requests {
        httpReq, err := http.NewRequest(req.method, server.URL+req.path, nil)
        if err != nil {
            fmt.Println("build request error:", err)
            continue
        }
        resp, err := http.DefaultClient.Do(httpReq)
        if err != nil {
            fmt.Println("request error:", err)
            continue
        }
        body, _ := io.ReadAll(resp.Body)
        resp.Body.Close()
        fmt.Printf("%s %s -> %d %q\n", req.method, req.path, resp.StatusCode, string(body))
    }
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output** *(the port number is OS-assigned and will differ on your machine - everything else is exact)*:

```
server listening on: http://127.0.0.1:65059
GET /items -> 200 "list of items\n"
POST /items -> 200 "created an item\n"
GET /items/42 -> 200 "item id=42\n"
DELETE /items/42 -> 200 "deleted item id=42\n"
PUT /items -> 405 "Method Not Allowed\n"
```

**Learning Objectives:**
- ✅ Register method-specific patterns (`"GET /items"`) and wildcard path segments (`"{id}"`) on `http.ServeMux`
- ✅ Read a matched wildcard with `r.PathValue`
- ✅ See, with real output, that a registered path with no matching method automatically returns `405 Method Not Allowed` - you didn't write that behavior, `ServeMux` did

---

## Exercise 3: Reading Query Parameters

**Objective:** Read and validate optional and required query parameters

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level27-exercise3
cd ~/projects/level27-exercise3
go mod init level27.example/exercise3
```

2. Create `handler.go`:

```bash
cat > handler.go << 'EOF'
package api

import (
    "fmt"
    "net/http"
    "strconv"
)

// SearchHandler reads query parameters from the request URL:
//   ?q=<term>          required search term
//   ?limit=<n>         optional, defaults to 10
func SearchHandler(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query()

    q := query.Get("q")
    if q == "" {
        http.Error(w, "missing required query parameter: q", http.StatusBadRequest)
        return
    }

    limit := 10
    if raw := query.Get("limit"); raw != "" {
        n, err := strconv.Atoi(raw)
        if err != nil || n < 1 {
            http.Error(w, "invalid limit: must be a positive integer", http.StatusBadRequest)
            return
        }
        limit = n
    }

    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    fmt.Fprintf(w, "searching for %q, limit=%d\n", q, limit)
}
EOF
```

3. Create `handler_test.go`:

```bash
cat > handler_test.go << 'EOF'
package api

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestSearchHandler(t *testing.T) {
    tests := []struct {
        name       string
        url        string
        wantStatus int
        wantBody   string
    }{
        {"missing q", "/search", http.StatusBadRequest, "missing required query parameter: q\n"},
        {"q only, default limit", "/search?q=golang", http.StatusOK, `searching for "golang", limit=10` + "\n"},
        {"q and limit", "/search?q=golang&limit=5", http.StatusOK, `searching for "golang", limit=5` + "\n"},
        {"invalid limit", "/search?q=golang&limit=abc", http.StatusBadRequest, "invalid limit: must be a positive integer\n"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest(http.MethodGet, tt.url, nil)
            rec := httptest.NewRecorder()

            SearchHandler(rec, req)

            if rec.Code != tt.wantStatus {
                t.Errorf("status = %d; want %d", rec.Code, tt.wantStatus)
            }
            if got := rec.Body.String(); got != tt.wantBody {
                t.Errorf("body = %q; want %q", got, tt.wantBody)
            }
        })
    }
}
EOF
```

4. Run the tests:

```bash
go test -v ./...
```

**Expected Output:**

```
=== RUN   TestSearchHandler
=== RUN   TestSearchHandler/missing_q
=== RUN   TestSearchHandler/q_only,_default_limit
=== RUN   TestSearchHandler/q_and_limit
=== RUN   TestSearchHandler/invalid_limit
--- PASS: TestSearchHandler (0.00s)
    --- PASS: TestSearchHandler/missing_q (0.00s)
    --- PASS: TestSearchHandler/q_only,_default_limit (0.00s)
    --- PASS: TestSearchHandler/q_and_limit (0.00s)
    --- PASS: TestSearchHandler/invalid_limit (0.00s)
PASS
ok  	level27.example/exercise3	0.692s
```

**Learning Objectives:**
- ✅ Read query parameters with `r.URL.Query().Get(key)`
- ✅ Distinguish a missing required parameter from an invalid optional one, with different messages
- ✅ Table-drive a handler test across several query-string variations

---

## Exercise 4: Decoding JSON Requests and Encoding JSON Responses

**Objective:** Decode a JSON request body, validate it, and encode a JSON response with an explicit status code

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level27-exercise4
cd ~/projects/level27-exercise4
go mod init level27.example/exercise4
```

2. Create `handler.go`:

```bash
cat > handler.go << 'EOF'
package api

import (
    "encoding/json"
    "net/http"
)

// CreateUserRequest is the shape of the JSON body clients send.
type CreateUserRequest struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

// CreateUserResponse is the shape of the JSON body the handler sends back.
type CreateUserResponse struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Age  int    `json:"age"`
}

// CreateUserHandler decodes a JSON request body (bridging to Level 19),
// validates it, and encodes a JSON response with an explicit status code.
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON body: " + err.Error()})
        return
    }

    if req.Name == "" {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "name is required"})
        return
    }

    resp := CreateUserResponse{ID: 1, Name: req.Name, Age: req.Age}

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(resp)
}
EOF
```

3. Create `handler_test.go`:

```bash
cat > handler_test.go << 'EOF'
package api

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

func TestCreateUserHandler(t *testing.T) {
    tests := []struct {
        name       string
        body       string
        wantStatus int
        wantBody   string
    }{
        {
            name:       "valid request",
            body:       `{"name":"Alice","age":30}`,
            wantStatus: http.StatusCreated,
            wantBody:   `{"id":1,"name":"Alice","age":30}` + "\n",
        },
        {
            name:       "missing name",
            body:       `{"age":30}`,
            wantStatus: http.StatusBadRequest,
            wantBody:   `{"error":"name is required"}` + "\n",
        },
        {
            name:       "malformed JSON",
            body:       `{"name":`,
            wantStatus: http.StatusBadRequest,
            wantBody:   "", // checked separately below - error text includes decoder detail
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(tt.body))
            rec := httptest.NewRecorder()

            CreateUserHandler(rec, req)

            if rec.Code != tt.wantStatus {
                t.Errorf("status = %d; want %d", rec.Code, tt.wantStatus)
            }
            if got := rec.Header().Get("Content-Type"); got != "application/json" {
                t.Errorf("Content-Type = %q; want %q", got, "application/json")
            }
            if tt.wantBody != "" {
                if got := rec.Body.String(); got != tt.wantBody {
                    t.Errorf("body = %q; want %q", got, tt.wantBody)
                }
            } else if !strings.Contains(rec.Body.String(), "invalid JSON body") {
                t.Errorf("body = %q; want it to contain %q", rec.Body.String(), "invalid JSON body")
            }
        })
    }
}
EOF
```

4. Run the tests:

```bash
go test -v ./...
```

**Expected Output:**

```
=== RUN   TestCreateUserHandler
=== RUN   TestCreateUserHandler/valid_request
=== RUN   TestCreateUserHandler/missing_name
=== RUN   TestCreateUserHandler/malformed_JSON
--- PASS: TestCreateUserHandler (0.00s)
    --- PASS: TestCreateUserHandler/valid_request (0.00s)
    --- PASS: TestCreateUserHandler/missing_name (0.00s)
    --- PASS: TestCreateUserHandler/malformed_JSON (0.00s)
PASS
ok  	level27.example/exercise4	0.560s
```

**Learning Objectives:**
- ✅ Decode a JSON request body with `json.NewDecoder(r.Body).Decode(&v)` and always check the error
- ✅ Encode a JSON response with `json.NewEncoder(w).Encode(v)` after setting `Content-Type` and the status code
- ✅ Return a `400 Bad Request` for both malformed JSON and a failed validation rule

---

## Exercise 5: An In-Memory CRUD REST Resource

**Objective:** Build list, create, get-by-id, and delete handlers backed by a thread-safe in-memory store

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level27-exercise5
cd ~/projects/level27-exercise5
go mod init level27.example/exercise5
```

2. Create `store.go`:

```bash
cat > store.go << 'EOF'
package main

import "sync"

// Book is the resource this API manages.
type Book struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
}

// BookStore is a tiny thread-safe in-memory store (Level 23: Mutex).
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
EOF
```

3. Create `handlers.go`:

```bash
cat > handlers.go << 'EOF'
package main

import (
    "encoding/json"
    "net/http"
    "strconv"
)

type createBookRequest struct {
    Title  string `json:"title"`
    Author string `json:"author"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(v)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
    writeJSON(w, status, map[string]string{"error": message})
}

func newBookMux(store *BookStore) *http.ServeMux {
    mux := http.NewServeMux()

    mux.HandleFunc("GET /books", func(w http.ResponseWriter, r *http.Request) {
        writeJSON(w, http.StatusOK, store.List())
    })

    mux.HandleFunc("POST /books", func(w http.ResponseWriter, r *http.Request) {
        var req createBookRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
            return
        }
        if req.Title == "" {
            writeJSONError(w, http.StatusBadRequest, "title is required")
            return
        }
        book := store.Create(req.Title, req.Author)
        writeJSON(w, http.StatusCreated, book)
    })

    mux.HandleFunc("GET /books/{id}", func(w http.ResponseWriter, r *http.Request) {
        id, err := strconv.Atoi(r.PathValue("id"))
        if err != nil {
            writeJSONError(w, http.StatusBadRequest, "invalid book id")
            return
        }
        book, ok := store.Get(id)
        if !ok {
            writeJSONError(w, http.StatusNotFound, "book not found")
            return
        }
        writeJSON(w, http.StatusOK, book)
    })

    mux.HandleFunc("DELETE /books/{id}", func(w http.ResponseWriter, r *http.Request) {
        id, err := strconv.Atoi(r.PathValue("id"))
        if err != nil {
            writeJSONError(w, http.StatusBadRequest, "invalid book id")
            return
        }
        if !store.Delete(id) {
            writeJSONError(w, http.StatusNotFound, "book not found")
            return
        }
        w.WriteHeader(http.StatusNoContent)
    })

    return mux
}
EOF
```

4. Create `main.go` to drive it end to end over a real (loopback-only) server:

```bash
cat > main.go << 'EOF'
package main

import (
    "bytes"
    "fmt"
    "io"
    "net/http"
    "net/http/httptest"
)

func main() {
    store := NewBookStore()
    server := httptest.NewServer(newBookMux(store))
    defer server.Close()

    get := func(path string) {
        resp, err := http.Get(server.URL + path)
        if err != nil {
            fmt.Println("GET error:", err)
            return
        }
        body, _ := io.ReadAll(resp.Body)
        resp.Body.Close()
        fmt.Printf("GET %s -> %d %s\n", path, resp.StatusCode, body)
    }

    post := func(path, jsonBody string) {
        resp, err := http.Post(server.URL+path, "application/json", bytes.NewBufferString(jsonBody))
        if err != nil {
            fmt.Println("POST error:", err)
            return
        }
        body, _ := io.ReadAll(resp.Body)
        resp.Body.Close()
        fmt.Printf("POST %s %s -> %d %s\n", path, jsonBody, resp.StatusCode, body)
    }

    del := func(path string) {
        req, _ := http.NewRequest(http.MethodDelete, server.URL+path, nil)
        resp, err := http.DefaultClient.Do(req)
        if err != nil {
            fmt.Println("DELETE error:", err)
            return
        }
        body, _ := io.ReadAll(resp.Body)
        resp.Body.Close()
        fmt.Printf("DELETE %s -> %d %q\n", path, resp.StatusCode, string(body))
    }

    get("/books")
    post("/books", `{"title":"The Go Programming Language","author":"Donovan and Kernighan"}`)
    post("/books", `{"title":"The Pragmatic Programmer","author":"Hunt and Thomas"}`)
    get("/books")
    get("/books/1")
    get("/books/99")
    del("/books/1")
    get("/books")
    del("/books/1")
}
EOF
```

5. Run it:

```bash
go run .
```

**Expected Output:**

```
GET /books -> 200 []

POST /books {"title":"The Go Programming Language","author":"Donovan and Kernighan"} -> 201 {"id":1,"title":"The Go Programming Language","author":"Donovan and Kernighan"}

POST /books {"title":"The Pragmatic Programmer","author":"Hunt and Thomas"} -> 201 {"id":2,"title":"The Pragmatic Programmer","author":"Hunt and Thomas"}

GET /books -> 200 [{"id":1,"title":"The Go Programming Language","author":"Donovan and Kernighan"},{"id":2,"title":"The Pragmatic Programmer","author":"Hunt and Thomas"}]

GET /books/1 -> 200 {"id":1,"title":"The Go Programming Language","author":"Donovan and Kernighan"}

GET /books/99 -> 404 {"error":"book not found"}

DELETE /books/1 -> 204 ""
GET /books -> 200 [{"id":2,"title":"The Pragmatic Programmer","author":"Hunt and Thomas"}]

DELETE /books/1 -> 404 "{\"error\":\"book not found\"}\n"
```

**Learning Objectives:**
- ✅ Combine routing, JSON decode/encode, and a mutex-protected store into one working REST resource
- ✅ Implement list, create, get-by-id, and delete over `http.ServeMux`
- ✅ See `204 No Content` (delete success, empty body) and `404` (delete/get on a missing id) used correctly

---

## Exercise 6: A Logging Middleware, Verified With Captured Output

**Objective:** Write `func(http.Handler) http.Handler` middleware and prove it logs real request data

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level27-exercise6
cd ~/projects/level27-exercise6
go mod init level27.example/exercise6
```

2. Create `middleware.go`:

```bash
cat > middleware.go << 'EOF'
package main

import (
    "log"
    "net/http"
    "time"
)

// statusRecorder wraps http.ResponseWriter to capture the status code that
// gets written, since the standard ResponseWriter doesn't expose it back.
type statusRecorder struct {
    http.ResponseWriter
    status int
}

func (r *statusRecorder) WriteHeader(code int) {
    r.status = code
    r.ResponseWriter.WriteHeader(code)
}

// Logging is middleware: it takes a handler and returns a new handler that
// wraps it with cross-cutting behavior (request logging), then delegates
// to the original. This is the func(http.Handler) http.Handler pattern.
func Logging(logger *log.Logger, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

        next.ServeHTTP(rec, r)

        logger.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, rec.status, time.Since(start).Round(time.Microsecond))
    })
}
EOF
```

3. Create `middleware_test.go`:

```bash
cat > middleware_test.go << 'EOF'
package main

import (
    "bytes"
    "log"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

func TestLoggingMiddleware(t *testing.T) {
    var buf bytes.Buffer
    logger := log.New(&buf, "", 0) // no timestamp prefix, for a deterministic assertion

    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusTeapot)
        w.Write([]byte("short and stout"))
    })

    wrapped := Logging(logger, handler)

    req := httptest.NewRequest(http.MethodGet, "/brew", nil)
    rec := httptest.NewRecorder()

    wrapped.ServeHTTP(rec, req)

    if rec.Code != http.StatusTeapot {
        t.Errorf("status = %d; want %d", rec.Code, http.StatusTeapot)
    }

    logged := buf.String()
    if !strings.Contains(logged, "GET /brew -> 418") {
        t.Errorf("log output = %q; want it to contain %q", logged, "GET /brew -> 418")
    }
    t.Logf("captured log line: %s", strings.TrimSpace(logged))
}
EOF
```

4. Run the tests:

```bash
go test -v ./...
```

**Expected Output** *(the exact microsecond duration is machine-dependent; the method, path, and status are exact)*:

```
=== RUN   TestLoggingMiddleware
    middleware_test.go:36: captured log line: GET /brew -> 418 (1µs)
--- PASS: TestLoggingMiddleware (0.00s)
PASS
ok  	level27.example/exercise6	0.536s
```

**Learning Objectives:**
- ✅ Write middleware matching `func(http.Handler) http.Handler`, wrapping cross-cutting behavior around any handler
- ✅ Capture the response status code from inside middleware by embedding and overriding `WriteHeader`
- ✅ Verify middleware behavior with a real captured log line, not just "it compiled"

---

## Exercise 7: A Full httptest.NewServer Round-Trip Test

**Objective:** Test routing, JSON, and status codes together over a real (loopback) TCP connection

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level27-exercise7
cd ~/projects/level27-exercise7
go mod init level27.example/exercise7
```

2. Create `store.go`:

```bash
cat > store.go << 'EOF'
package main

import "sync"

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
EOF
```

3. Create `handlers.go`:

```bash
cat > handlers.go << 'EOF'
package main

import (
    "encoding/json"
    "net/http"
    "strconv"
)

type createBookRequest struct {
    Title  string `json:"title"`
    Author string `json:"author"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(v)
}

func newBookMux(store *BookStore) *http.ServeMux {
    mux := http.NewServeMux()

    mux.HandleFunc("GET /books", func(w http.ResponseWriter, r *http.Request) {
        writeJSON(w, http.StatusOK, store.List())
    })

    mux.HandleFunc("POST /books", func(w http.ResponseWriter, r *http.Request) {
        var req createBookRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
            return
        }
        book := store.Create(req.Title, req.Author)
        writeJSON(w, http.StatusCreated, book)
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

    return mux
}
EOF
```

4. Create `server_test.go`:

```bash
cat > server_test.go << 'EOF'
package main

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

// TestBookAPIRoundTrip spins up a REAL httptest.NewServer - a real TCP
// listener on 127.0.0.1 with an OS-assigned port - and drives it with a
// real *http.Client, unlike httptest.NewRecorder's in-process fake. This
// exercises the full stack: routing, JSON encoding, and status codes,
// over an actual (loopback-only) network round trip.
func TestBookAPIRoundTrip(t *testing.T) {
    store := NewBookStore()
    server := httptest.NewServer(newBookMux(store))
    defer server.Close()

    client := server.Client()

    t.Run("list starts empty", func(t *testing.T) {
        resp, err := client.Get(server.URL + "/books")
        if err != nil {
            t.Fatalf("GET /books: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            t.Errorf("status = %d; want %d", resp.StatusCode, http.StatusOK)
        }
        var books []Book
        if err := json.NewDecoder(resp.Body).Decode(&books); err != nil {
            t.Fatalf("decode: %v", err)
        }
        if len(books) != 0 {
            t.Errorf("len(books) = %d; want 0", len(books))
        }
    })

    var createdID int
    t.Run("create a book", func(t *testing.T) {
        body := strings.NewReader(`{"title":"Clean Code","author":"Robert C. Martin"}`)
        resp, err := client.Post(server.URL+"/books", "application/json", body)
        if err != nil {
            t.Fatalf("POST /books: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusCreated {
            t.Errorf("status = %d; want %d", resp.StatusCode, http.StatusCreated)
        }
        var created Book
        if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
            t.Fatalf("decode: %v", err)
        }
        if created.Title != "Clean Code" {
            t.Errorf("title = %q; want %q", created.Title, "Clean Code")
        }
        createdID = created.ID
    })

    t.Run("get the created book by id", func(t *testing.T) {
        resp, err := client.Get(server.URL + "/books/1")
        if err != nil {
            t.Fatalf("GET /books/1: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            t.Errorf("status = %d; want %d", resp.StatusCode, http.StatusOK)
        }
        var got Book
        if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
            t.Fatalf("decode: %v", err)
        }
        if got.ID != createdID {
            t.Errorf("id = %d; want %d", got.ID, createdID)
        }
    })

    t.Run("get a nonexistent book returns 404", func(t *testing.T) {
        resp, err := client.Get(server.URL + "/books/999")
        if err != nil {
            t.Fatalf("GET /books/999: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusNotFound {
            t.Errorf("status = %d; want %d", resp.StatusCode, http.StatusNotFound)
        }
    })
}
EOF
```

5. Run the tests:

```bash
go test -v ./...
```

**Expected Output:**

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

**Learning Objectives:**
- ✅ Start a real (loopback-only) server for a test with `httptest.NewServer`
- ✅ Drive it with `server.Client()`, a real `*http.Client`
- ✅ Structure a multi-step round-trip test as subtests that share state (the created book's id) across steps

---

## Exercise 8: A Context-Timeout-Aware Handler

**Objective:** Bound simulated slow work with context.WithTimeout and map the outcome to the right status code

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level27-exercise8
cd ~/projects/level27-exercise8
go mod init level27.example/exercise8
```

2. Create `handler.go`:

```bash
cat > handler.go << 'EOF'
package main

import (
    "context"
    "encoding/json"
    "net/http"
    "time"
)

// slowOperation simulates work that takes `delay` to complete, but
// respects context cancellation - it stops early if ctx is done first.
func slowOperation(ctx context.Context, delay time.Duration) error {
    select {
    case <-time.After(delay):
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

// ReportHandler wraps the incoming request's context with a timeout
// (Level 24), so a slow downstream operation can't hang the handler
// forever. It checks ctx.Err() to decide which status code to send.
func ReportHandler(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 50*time.Millisecond)
    defer cancel()

    err := slowOperation(ctx, 200*time.Millisecond) // deliberately slower than the timeout

    w.Header().Set("Content-Type", "application/json")
    if err != nil {
        w.WriteHeader(http.StatusGatewayTimeout)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "report generation timed out",
            "cause": err.Error(),
        })
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "report ready"})
}

// FastReportHandler is the same shape but with a generous timeout, so the
// simulated work finishes in time - contrast case for the test below.
func FastReportHandler(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
    defer cancel()

    err := slowOperation(ctx, 10*time.Millisecond)

    w.Header().Set("Content-Type", "application/json")
    if err != nil {
        w.WriteHeader(http.StatusGatewayTimeout)
        json.NewEncoder(w).Encode(map[string]string{"error": "report generation timed out"})
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "report ready"})
}
EOF
```

3. Create `handler_test.go`:

```bash
cat > handler_test.go << 'EOF'
package main

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

func TestReportHandler_TimesOut(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/report", nil)
    rec := httptest.NewRecorder()

    ReportHandler(rec, req)

    if rec.Code != http.StatusGatewayTimeout {
        t.Errorf("status = %d; want %d", rec.Code, http.StatusGatewayTimeout)
    }
    if body := rec.Body.String(); !strings.Contains(body, "context deadline exceeded") {
        t.Errorf("body = %q; want it to contain %q", body, "context deadline exceeded")
    }
}

func TestFastReportHandler_Succeeds(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/report-fast", nil)
    rec := httptest.NewRecorder()

    FastReportHandler(rec, req)

    if rec.Code != http.StatusOK {
        t.Errorf("status = %d; want %d", rec.Code, http.StatusOK)
    }
    if body := rec.Body.String(); !strings.Contains(body, "report ready") {
        t.Errorf("body = %q; want it to contain %q", body, "report ready")
    }
}
EOF
```

4. Run the tests:

```bash
go test -v ./...
```

**Expected Output:**

```
=== RUN   TestReportHandler_TimesOut
--- PASS: TestReportHandler_TimesOut (0.05s)
=== RUN   TestFastReportHandler_Succeeds
--- PASS: TestFastReportHandler_Succeeds (0.01s)
PASS
ok  	level27.example/exercise8	0.646s
```

**Learning Objectives:**
- ✅ Derive `context.WithTimeout` from `r.Context()` inside a handler
- ✅ Use a `select` between the real work finishing and `ctx.Done()` (Level 24's pattern, applied to HTTP)
- ✅ Map a timed-out context to `504 Gateway Timeout` and a successful one to `200 OK`, in the same handler shape

---

## Exercise 9: Consistent JSON Error Responses

**Objective:** Map application errors (Level 16) to HTTP status codes through one shared error-writing function

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level27-exercise9
cd ~/projects/level27-exercise9
go mod init level27.example/exercise9
```

2. Create `errors.go`:

```bash
cat > errors.go << 'EOF'
package main

import "errors"

// Sentinel application errors (Level 16) - handlers map these to HTTP
// status codes instead of leaking implementation details to clients.
var (
    ErrNotFound     = errors.New("resource not found")
    ErrInvalidInput = errors.New("invalid input")
    ErrUnauthorized = errors.New("unauthorized")
)

// ValidationError wraps ErrInvalidInput with a field-specific message,
// demonstrating errors.Is/errors.As-friendly wrapping (Level 16).
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return "validation failed on " + e.Field + ": " + e.Message
}

func (e *ValidationError) Unwrap() error {
    return ErrInvalidInput
}
EOF
```

3. Create `handler.go`:

```bash
cat > handler.go << 'EOF'
package main

import (
    "encoding/json"
    "errors"
    "net/http"
)

// errorResponse is the ONE consistent JSON shape every error path returns,
// regardless of which internal error caused it.
type errorResponse struct {
    Error string `json:"error"`
}

// writeError maps an application error (Level 16) to an HTTP status code
// and writes it as the standard error JSON shape. This is the single place
// that translates Go errors into HTTP semantics - handlers never
// hand-roll status codes for errors themselves.
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

// lookupUser simulates a lookup that fails in different ways depending on id.
func lookupUser(id string) error {
    switch id {
    case "":
        return &ValidationError{Field: "id", Message: "must not be empty"}
    case "404":
        return ErrNotFound
    case "401":
        return ErrUnauthorized
    case "boom":
        return errors.New("unexpected database failure") // unmapped -> 500
    default:
        return nil
    }
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    if err := lookupUser(id); err != nil {
        writeError(w, err)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"id": id, "status": "found"})
}
EOF
```

4. Create `handler_test.go`:

```bash
cat > handler_test.go << 'EOF'
package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestUserHandler_ErrorMapping(t *testing.T) {
    tests := []struct {
        name       string
        id         string
        wantStatus int
        wantBody   string
    }{
        {"empty id -> validation error -> 400", "", http.StatusBadRequest, `{"error":"validation failed on id: must not be empty"}` + "\n"},
        {"404 id -> not found -> 404", "404", http.StatusNotFound, `{"error":"resource not found"}` + "\n"},
        {"401 id -> unauthorized -> 401", "401", http.StatusUnauthorized, `{"error":"unauthorized"}` + "\n"},
        {"boom id -> unmapped error -> 500", "boom", http.StatusInternalServerError, `{"error":"unexpected database failure"}` + "\n"},
        {"valid id -> 200", "42", http.StatusOK, `{"id":"42","status":"found"}` + "\n"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest(http.MethodGet, "/user?id="+tt.id, nil)
            rec := httptest.NewRecorder()

            UserHandler(rec, req)

            if rec.Code != tt.wantStatus {
                t.Errorf("status = %d; want %d", rec.Code, tt.wantStatus)
            }
            if got := rec.Body.String(); got != tt.wantBody {
                t.Errorf("body = %q; want %q", got, tt.wantBody)
            }
        })
    }
}
EOF
```

5. Run the tests:

```bash
go test -v ./...
```

**Expected Output:**

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

**Learning Objectives:**
- ✅ Map sentinel and wrapped application errors (Level 16) to HTTP status codes using `errors.Is`
- ✅ Return one consistent `{"error": "..."}` JSON shape across every failure path
- ✅ Confirm an unmapped/unexpected error correctly falls through to `500 Internal Server Error`

---

## Exercise 10: Comprehensive Practice — A Task Manager REST API

**Objective:** Combine routing, JSON, middleware, error mapping, and context into one fully tested API

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level27-exercise10
cd ~/projects/level27-exercise10
go mod init level27.example/exercise10
```

2. Create `store.go`:

```bash
cat > store.go << 'EOF'
package main

import (
    "errors"
    "sync"
)

// ErrTaskNotFound is the sentinel error (Level 16) the store returns when
// an id doesn't exist - handlers map it to a 404.
var ErrTaskNotFound = errors.New("task not found")

type Task struct {
    ID    int    `json:"id"`
    Title string `json:"title"`
    Done  bool   `json:"done"`
}

// TaskStore is a thread-safe (Level 23) in-memory task list.
type TaskStore struct {
    mu     sync.Mutex
    nextID int
    tasks  map[int]Task
}

func NewTaskStore() *TaskStore {
    return &TaskStore{nextID: 1, tasks: make(map[int]Task)}
}

func (s *TaskStore) List() []Task {
    s.mu.Lock()
    defer s.mu.Unlock()
    result := make([]Task, 0, len(s.tasks))
    for id := 1; id < s.nextID; id++ {
        if t, ok := s.tasks[id]; ok {
            result = append(result, t)
        }
    }
    return result
}

func (s *TaskStore) Create(title string) Task {
    s.mu.Lock()
    defer s.mu.Unlock()
    t := Task{ID: s.nextID, Title: title, Done: false}
    s.tasks[t.ID] = t
    s.nextID++
    return t
}

func (s *TaskStore) Get(id int) (Task, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    t, ok := s.tasks[id]
    if !ok {
        return Task{}, ErrTaskNotFound
    }
    return t, nil
}

func (s *TaskStore) MarkDone(id int) (Task, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    t, ok := s.tasks[id]
    if !ok {
        return Task{}, ErrTaskNotFound
    }
    t.Done = true
    s.tasks[id] = t
    return t, nil
}

func (s *TaskStore) Delete(id int) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    if _, ok := s.tasks[id]; !ok {
        return ErrTaskNotFound
    }
    delete(s.tasks, id)
    return nil
}
EOF
```

3. Create `middleware.go`:

```bash
cat > middleware.go << 'EOF'
package main

import (
    "log"
    "net/http"
    "time"
)

type statusRecorder struct {
    http.ResponseWriter
    status int
}

func (r *statusRecorder) WriteHeader(code int) {
    r.status = code
    r.ResponseWriter.WriteHeader(code)
}

// Logging is cross-cutting middleware: func(http.Handler) http.Handler.
func Logging(logger *log.Logger, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
        next.ServeHTTP(rec, r)
        logger.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, rec.status, time.Since(start).Round(time.Microsecond))
    })
}
EOF
```

4. Create `errors.go`:

```bash
cat > errors.go << 'EOF'
package main

import (
    "encoding/json"
    "errors"
    "net/http"
)

type errorResponse struct {
    Error string `json:"error"`
}

// writeError maps application errors (Level 16) to HTTP status codes and
// writes the ONE consistent JSON error shape used everywhere in this API.
func writeError(w http.ResponseWriter, err error) {
    status := http.StatusInternalServerError
    switch {
    case errors.Is(err, ErrTaskNotFound):
        status = http.StatusNotFound
    case errors.Is(err, ErrInvalidTask):
        status = http.StatusBadRequest
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(v)
}
EOF
```

5. Create `handlers.go` - the API's routing table, plus a context-bounded `/stats` endpoint:

```bash
cat > handlers.go << 'EOF'
package main

import (
    "context"
    "encoding/json"
    "errors"
    "net/http"
    "strconv"
    "time"
)

// ErrInvalidTask is returned when a create/update request body is invalid.
var ErrInvalidTask = errors.New("invalid task")

type createTaskRequest struct {
    Title string `json:"title"`
}

// slowStatsQuery simulates a slow aggregate computation over the task
// store, respecting context cancellation (Level 24).
func slowStatsQuery(ctx context.Context, delay time.Duration, compute func() map[string]int) (map[string]int, error) {
    resultCh := make(chan map[string]int, 1)
    go func() {
        time.Sleep(delay)
        resultCh <- compute()
    }()

    select {
    case result := <-resultCh:
        return result, nil
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

func newTaskMux(store *TaskStore) *http.ServeMux {
    mux := http.NewServeMux()

    mux.HandleFunc("GET /tasks", func(w http.ResponseWriter, r *http.Request) {
        writeJSON(w, http.StatusOK, store.List())
    })

    mux.HandleFunc("POST /tasks", func(w http.ResponseWriter, r *http.Request) {
        var req createTaskRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" {
            writeError(w, ErrInvalidTask)
            return
        }
        task := store.Create(req.Title)
        writeJSON(w, http.StatusCreated, task)
    })

    mux.HandleFunc("GET /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
        id, err := strconv.Atoi(r.PathValue("id"))
        if err != nil {
            writeError(w, ErrInvalidTask)
            return
        }
        task, err := store.Get(id)
        if err != nil {
            writeError(w, err)
            return
        }
        writeJSON(w, http.StatusOK, task)
    })

    mux.HandleFunc("PATCH /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
        id, err := strconv.Atoi(r.PathValue("id"))
        if err != nil {
            writeError(w, ErrInvalidTask)
            return
        }
        task, err := store.MarkDone(id)
        if err != nil {
            writeError(w, err)
            return
        }
        writeJSON(w, http.StatusOK, task)
    })

    mux.HandleFunc("DELETE /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
        id, err := strconv.Atoi(r.PathValue("id"))
        if err != nil {
            writeError(w, ErrInvalidTask)
            return
        }
        if err := store.Delete(id); err != nil {
            writeError(w, err)
            return
        }
        w.WriteHeader(http.StatusNoContent)
    })

    // GET /stats simulates a slow aggregate query, bounded by a context
    // timeout, so a slow backend can never hang this handler forever.
    mux.HandleFunc("GET /stats", func(w http.ResponseWriter, r *http.Request) {
        ctx, cancel := context.WithTimeout(r.Context(), 30*time.Millisecond)
        defer cancel()

        result, err := slowStatsQuery(ctx, 100*time.Millisecond, func() map[string]int {
            tasks := store.List()
            done := 0
            for _, t := range tasks {
                if t.Done {
                    done++
                }
            }
            return map[string]int{"total": len(tasks), "done": done}
        })
        if err != nil {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusGatewayTimeout)
            json.NewEncoder(w).Encode(errorResponse{Error: "stats query timed out: " + err.Error()})
            return
        }
        writeJSON(w, http.StatusOK, result)
    })

    return mux
}
EOF
```

6. Create `task_test.go` - a table-driven suite that drives the full middleware-wrapped mux through `httptest.NewRecorder`, plus a content-type check:

```bash
cat > task_test.go << 'EOF'
package main

import (
    "bytes"
    "encoding/json"
    "log"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

// newTestServer builds the full stack under test: store -> mux -> logging
// middleware - exactly what a real deployment would run, minus the actual
// TCP listener (httptest.NewRecorder drives it in-process instead).
func newTestServer() (http.Handler, *bytes.Buffer) {
    var logBuf bytes.Buffer
    logger := log.New(&logBuf, "", 0)
    store := NewTaskStore()
    return Logging(logger, newTaskMux(store)), &logBuf
}

func TestTaskAPI_TableDriven(t *testing.T) {
    handler, logBuf := newTestServer()

    type step struct {
        name       string
        method     string
        path       string
        body       string
        wantStatus int
        wantBody   string // exact match when non-empty
        wantSubstr string // substring match when non-empty (for varying bodies)
    }

    steps := []step{
        {name: "list is empty", method: http.MethodGet, path: "/tasks", wantStatus: http.StatusOK, wantBody: "[]\n"},
        {name: "reject empty title", method: http.MethodPost, path: "/tasks", body: `{"title":""}`, wantStatus: http.StatusBadRequest, wantBody: `{"error":"invalid task"}` + "\n"},
        {name: "create a task", method: http.MethodPost, path: "/tasks", body: `{"title":"Write tests"}`, wantStatus: http.StatusCreated, wantBody: `{"id":1,"title":"Write tests","done":false}` + "\n"},
        {name: "create a second task", method: http.MethodPost, path: "/tasks", body: `{"title":"Ship it"}`, wantStatus: http.StatusCreated, wantBody: `{"id":2,"title":"Ship it","done":false}` + "\n"},
        {name: "list has two tasks", method: http.MethodGet, path: "/tasks", wantStatus: http.StatusOK, wantBody: `[{"id":1,"title":"Write tests","done":false},{"id":2,"title":"Ship it","done":false}]` + "\n"},
        {name: "get task 1", method: http.MethodGet, path: "/tasks/1", wantStatus: http.StatusOK, wantBody: `{"id":1,"title":"Write tests","done":false}` + "\n"},
        {name: "get missing task", method: http.MethodGet, path: "/tasks/99", wantStatus: http.StatusNotFound, wantBody: `{"error":"task not found"}` + "\n"},
        {name: "mark task 1 done", method: http.MethodPatch, path: "/tasks/1", wantStatus: http.StatusOK, wantBody: `{"id":1,"title":"Write tests","done":true}` + "\n"},
        {name: "stats after one done", method: http.MethodGet, path: "/stats", wantStatus: http.StatusGatewayTimeout, wantSubstr: "stats query timed out"},
        {name: "delete task 1", method: http.MethodDelete, path: "/tasks/1", wantStatus: http.StatusNoContent, wantBody: ""},
        {name: "delete task 1 again is 404", method: http.MethodDelete, path: "/tasks/1", wantStatus: http.StatusNotFound, wantBody: `{"error":"task not found"}` + "\n"},
        {name: "invalid id is 400", method: http.MethodGet, path: "/tasks/abc", wantStatus: http.StatusBadRequest, wantBody: `{"error":"invalid task"}` + "\n"},
    }

    for _, s := range steps {
        t.Run(s.name, func(t *testing.T) {
            var body *bytes.Reader
            if s.body != "" {
                body = bytes.NewReader([]byte(s.body))
            } else {
                body = bytes.NewReader(nil)
            }
            req := httptest.NewRequest(s.method, s.path, body)
            rec := httptest.NewRecorder()

            handler.ServeHTTP(rec, req)

            if rec.Code != s.wantStatus {
                t.Errorf("status = %d; want %d (body=%s)", rec.Code, s.wantStatus, rec.Body.String())
            }
            got := rec.Body.String()
            if s.wantSubstr != "" {
                if !strings.Contains(got, s.wantSubstr) {
                    t.Errorf("body = %q; want it to contain %q", got, s.wantSubstr)
                }
            } else if got != s.wantBody {
                t.Errorf("body = %q; want %q", got, s.wantBody)
            }
        })
    }

    t.Run("middleware logged every request", func(t *testing.T) {
        logged := logBuf.String()
        lines := strings.Count(logged, "\n")
        if lines != len(steps) {
            t.Errorf("logged %d lines; want %d (one per request)", lines, len(steps))
        }
        if !strings.Contains(logged, "GET /tasks -> 200") {
            t.Errorf("log output missing expected line, got:\n%s", logged)
        }
    })
}

func TestTaskAPI_ContentTypeIsAlwaysJSON(t *testing.T) {
    handler, _ := newTestServer()

    req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
    rec := httptest.NewRecorder()
    handler.ServeHTTP(rec, req)

    if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
        t.Errorf("Content-Type = %q; want %q", ct, "application/json")
    }

    var tasks []Task
    if err := json.NewDecoder(rec.Body).Decode(&tasks); err != nil {
        t.Fatalf("response body did not decode as JSON: %v", err)
    }
}
EOF
```

7. Run the tests:

```bash
go test -v ./...
```

**Expected Output:**

```
=== RUN   TestTaskAPI_TableDriven
=== RUN   TestTaskAPI_TableDriven/list_is_empty
=== RUN   TestTaskAPI_TableDriven/reject_empty_title
=== RUN   TestTaskAPI_TableDriven/create_a_task
=== RUN   TestTaskAPI_TableDriven/create_a_second_task
=== RUN   TestTaskAPI_TableDriven/list_has_two_tasks
=== RUN   TestTaskAPI_TableDriven/get_task_1
=== RUN   TestTaskAPI_TableDriven/get_missing_task
=== RUN   TestTaskAPI_TableDriven/mark_task_1_done
=== RUN   TestTaskAPI_TableDriven/stats_after_one_done
=== RUN   TestTaskAPI_TableDriven/delete_task_1
=== RUN   TestTaskAPI_TableDriven/delete_task_1_again_is_404
=== RUN   TestTaskAPI_TableDriven/invalid_id_is_400
=== RUN   TestTaskAPI_TableDriven/middleware_logged_every_request
--- PASS: TestTaskAPI_TableDriven (0.03s)
    --- PASS: TestTaskAPI_TableDriven/list_is_empty (0.00s)
    --- PASS: TestTaskAPI_TableDriven/reject_empty_title (0.00s)
    --- PASS: TestTaskAPI_TableDriven/create_a_task (0.00s)
    --- PASS: TestTaskAPI_TableDriven/create_a_second_task (0.00s)
    --- PASS: TestTaskAPI_TableDriven/list_has_two_tasks (0.00s)
    --- PASS: TestTaskAPI_TableDriven/get_task_1 (0.00s)
    --- PASS: TestTaskAPI_TableDriven/get_missing_task (0.00s)
    --- PASS: TestTaskAPI_TableDriven/mark_task_1_done (0.00s)
    --- PASS: TestTaskAPI_TableDriven/stats_after_one_done (0.03s)
    --- PASS: TestTaskAPI_TableDriven/delete_task_1 (0.00s)
    --- PASS: TestTaskAPI_TableDriven/delete_task_1_again_is_404 (0.00s)
    --- PASS: TestTaskAPI_TableDriven/invalid_id_is_400 (0.00s)
    --- PASS: TestTaskAPI_TableDriven/middleware_logged_every_request (0.00s)
=== RUN   TestTaskAPI_ContentTypeIsAlwaysJSON
--- PASS: TestTaskAPI_ContentTypeIsAlwaysJSON (0.00s)
PASS
ok  	level27.example/exercise10	1.901s
```

**Learning Objectives:**
- ✅ Combine `http.ServeMux` routing, JSON decode/encode, logging middleware, and a context-bounded endpoint in one API
- ✅ Drive a middleware-wrapped handler through a table of steps that share state (created task ids) across requests
- ✅ Confirm cross-cutting behavior (every request gets logged) as part of the same test suite that checks business logic
- ✅ Recognize this as the shape of a real, small production service - built entirely with the standard library

---

## Bonus Challenges

### Challenge 1: Token-Bucket Rate-Limiting Middleware

Write `func(http.Handler) http.Handler` middleware that allows only N requests per second per client, using a token bucket refilled by a `time.Ticker` (recall Level 22's ticker pattern). Requests over the limit get `429 Too Many Requests`.

```bash
mkdir -p ~/projects/level27-bonus1
cd ~/projects/level27-bonus1
go mod init level27.example/bonus1
```

**Hints:**
- A struct with a buffered `chan struct{}` of capacity N works as the bucket: `Allow()` tries a non-blocking receive (`select` with `default`)
- A background goroutine with a `time.Ticker` refills one token at a time with a non-blocking send
- Wrap the check at the top of the middleware's returned `http.HandlerFunc`, before calling `next.ServeHTTP`
- Test it with `httptest.NewRecorder` by firing more requests than the bucket capacity in a tight loop and asserting some get `429`

### Challenge 2: Request-ID Middleware

Write middleware that generates a unique ID for each incoming request, stores it in the request's context with `context.WithValue`, and includes it in both the log line and an `X-Request-ID` response header.

```bash
mkdir -p ~/projects/level27-bonus2
cd ~/projects/level27-bonus2
go mod init level27.example/bonus2
```

**Hints:**
- A simple counter (`atomic.Int64`, Level 23) or `crypto/rand` bytes hex-encoded both make a fine "UUID-like" ID for this exercise - a real UUID library isn't required
- Define an unexported context key type (`type contextKey string`) rather than a bare string, per Level 24's `context.WithValue` guidance
- `r = r.WithContext(context.WithValue(r.Context(), requestIDKey, id))` before calling `next.ServeHTTP(w, r)`
- Retrieve it later with a small helper: `func RequestIDFromContext(ctx context.Context) (string, bool)`

### Challenge 3: Pagination Query Parameters

Add `?page=` and `?limit=` support to a list handler, returning only the requested slice of results along with metadata about the total count and whether more pages exist.

```bash
mkdir -p ~/projects/level27-bonus3
cd ~/projects/level27-bonus3
go mod init level27.example/bonus3
```

**Hints:**
- Default `page=1`, `limit=10` when absent; reject non-positive values with `400`
- Compute `start := (page-1)*limit` and `end := min(start+limit, len(items))`; if `start >= len(items)`, return an empty slice, not an error
- Wrap the response in an envelope: `{"items": [...], "page": 1, "limit": 10, "total": 42, "has_more": true}`
- Table-drive the test across page 1, a middle page, the last (partial) page, and a page past the end

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Write handlers matching `func(http.ResponseWriter, *http.Request)` and unit-test them with `httptest.NewRecorder`
✅ Route requests with `http.ServeMux`'s real method-specific and wildcard pattern syntax
✅ Read path values, query parameters, and JSON request bodies
✅ Encode JSON responses with correct `Content-Type` headers and status codes
✅ Build a complete in-memory CRUD REST resource, thread-safe under concurrent requests
✅ Write and verify `func(http.Handler) http.Handler` middleware with real captured log output
✅ Test the full stack over a real (loopback-only) connection with `httptest.NewServer`
✅ Bound slow work in a handler with `context.WithTimeout` derived from `r.Context()`
✅ Map application errors to HTTP status codes through one consistent JSON error shape
✅ Combine all of the above into a fully tested, table-driven REST API

---

## Next Level

Level 28: Gin Framework
- What a routing framework buys you over hand-rolled `http.ServeMux`
- Gin's context object, route groups, and middleware chains
- Binding and validating JSON with struct tags
- Building the same kind of REST API faster, with less boilerplate

Great work! You just built REST APIs the way the standard library actually works underneath every framework! 🚀
