# Level 28: Gin Framework - Exercises

> ⚠️ **This level requires internet access.** Every exercise below installs a real third-party package - `go get github.com/gin-gonic/gin@v1.12.0` - the exact version this course was verified against. Every command, every program, and every "Expected Output" block on this page was actually run with `go run` / `go test` against that pinned version - nothing here is hand-computed. If `go get` fails, check your network connection before assuming your code is wrong. Every other level in this course works fully offline; this is the one exception.

Complete all exercises in order. Each exercise builds on previous knowledge.

---

## Exercise 1: gin.Default() vs gin.New()

**Objective:** Set up a Gin engine both ways and see that `gin.Default()` is just `gin.New()` plus two middleware calls

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level28-exercise1
cd ~/projects/level28-exercise1
go mod init level28.example/exercise1
go get github.com/gin-gonic/gin@v1.12.0
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "io"
    "net/http/httptest"

    "github.com/gin-gonic/gin"
)

func callPing(router *gin.Engine, label string) {
    w := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/ping", nil)
    router.ServeHTTP(w, req)
    fmt.Printf("[%s] GET /ping -> status=%d body=%s\n", label, w.Code, w.Body.String())
}

func main() {
    gin.SetMode(gin.TestMode)
    gin.DefaultWriter = io.Discard

    fmt.Println("=== gin.New(): a bare engine, no middleware attached ===")
    bare := gin.New()
    bare.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })
    callPing(bare, "bare")

    fmt.Println("\n=== Rebuilding gin.Default() by hand ===")
    handRolled := gin.New()
    handRolled.Use(gin.Logger(), gin.Recovery())
    handRolled.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })
    callPing(handRolled, "hand-rolled")

    fmt.Println("\n=== gin.Default() (same result, one call) ===")
    def := gin.Default()
    def.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })
    callPing(def, "gin.Default()")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== gin.New(): a bare engine, no middleware attached ===
[bare] GET /ping -> status=200 body={"message":"pong"}

=== Rebuilding gin.Default() by hand ===
[hand-rolled] GET /ping -> status=200 body={"message":"pong"}

=== gin.Default() (same result, one call) ===
[gin.Default()] GET /ping -> status=200 body={"message":"pong"}
```

**Learning Objectives:**
- ✅ Create a bare engine with `gin.New()`
- ✅ Understand `gin.Default()` is `gin.New()` + `Logger()` + `Recovery()`
- ✅ Simulate a request against an engine with `httptest`, without starting a real server

---

## Exercise 2: Routing With Path Parameters

**Objective:** Use `:param` and `*wildcard` path segments with `c.Param`

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level28-exercise2
cd ~/projects/level28-exercise2
go mod init level28.example/exercise2
go get github.com/gin-gonic/gin@v1.12.0
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "io"
    "net/http/httptest"

    "github.com/gin-gonic/gin"
)

func call(router *gin.Engine, method, path string) {
    w := httptest.NewRecorder()
    req := httptest.NewRequest(method, path, nil)
    router.ServeHTTP(w, req)
    fmt.Printf("%s %s -> status=%d body=%s\n", method, path, w.Code, w.Body.String())
}

func main() {
    gin.SetMode(gin.TestMode)
    gin.DefaultWriter = io.Discard

    router := gin.New()

    router.GET("/users/:id", func(c *gin.Context) {
        id := c.Param("id")
        c.JSON(200, gin.H{"user_id": id})
    })

    router.GET("/users/:id/posts/:postId", func(c *gin.Context) {
        id := c.Param("id")
        postID := c.Param("postId")
        c.JSON(200, gin.H{"user_id": id, "post_id": postID})
    })

    router.GET("/files/*filepath", func(c *gin.Context) {
        path := c.Param("filepath")
        c.JSON(200, gin.H{"filepath": path})
    })

    fmt.Println("=== Single path parameter ===")
    call(router, "GET", "/users/42")

    fmt.Println("\n=== Multiple path parameters ===")
    call(router, "GET", "/users/42/posts/7")

    fmt.Println("\n=== Wildcard parameter (catch-all) ===")
    call(router, "GET", "/files/images/2024/photo.png")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Single path parameter ===
GET /users/42 -> status=200 body={"user_id":"42"}

=== Multiple path parameters ===
GET /users/42/posts/7 -> status=200 body={"post_id":"7","user_id":"42"}

=== Wildcard parameter (catch-all) ===
GET /files/images/2024/photo.png -> status=200 body={"filepath":"/images/2024/photo.png"}
```

**Learning Objectives:**
- ✅ Read a path parameter with `c.Param`
- ✅ Use multiple path parameters on one route
- ✅ Understand that a `*wildcard` catch-all keeps its leading slash
- ✅ Notice `gin.H` (a map) serializes keys alphabetically, unlike a struct

---

## Exercise 3: Route Groups

**Objective:** Use `router.Group` for API versioning and nested prefixes

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level28-exercise3
cd ~/projects/level28-exercise3
go mod init level28.example/exercise3
go get github.com/gin-gonic/gin@v1.12.0
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "io"
    "net/http/httptest"

    "github.com/gin-gonic/gin"
)

func call(router *gin.Engine, method, path string) {
    w := httptest.NewRecorder()
    req := httptest.NewRequest(method, path, nil)
    router.ServeHTTP(w, req)
    fmt.Printf("%s %s -> status=%d body=%s\n", method, path, w.Code, w.Body.String())
}

func main() {
    gin.SetMode(gin.TestMode)
    gin.DefaultWriter = io.Discard

    router := gin.New()

    v1 := router.Group("/api/v1")
    {
        v1.GET("/users", func(c *gin.Context) {
            c.JSON(200, gin.H{"version": "v1", "users": []string{"alice", "bob"}})
        })
        v1.GET("/users/:id", func(c *gin.Context) {
            c.JSON(200, gin.H{"version": "v1", "user_id": c.Param("id")})
        })
    }

    v2 := router.Group("/api/v2")
    {
        v2.GET("/users", func(c *gin.Context) {
            c.JSON(200, gin.H{"version": "v2", "users": []gin.H{
                {"id": "1", "name": "alice"},
                {"id": "2", "name": "bob"},
            }})
        })
    }

    admin := router.Group("/admin")
    {
        reports := admin.Group("/reports")
        reports.GET("/summary", func(c *gin.Context) {
            c.JSON(200, gin.H{"section": "admin/reports", "report": "summary"})
        })
    }

    fmt.Println("=== v1 group ===")
    call(router, "GET", "/api/v1/users")
    call(router, "GET", "/api/v1/users/7")

    fmt.Println("\n=== v2 group (same resource, different shape) ===")
    call(router, "GET", "/api/v2/users")

    fmt.Println("\n=== Nested groups (admin -> reports) ===")
    call(router, "GET", "/admin/reports/summary")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== v1 group ===
GET /api/v1/users -> status=200 body={"users":["alice","bob"],"version":"v1"}
GET /api/v1/users/7 -> status=200 body={"user_id":"7","version":"v1"}

=== v2 group (same resource, different shape) ===
GET /api/v2/users -> status=200 body={"users":[{"id":"1","name":"alice"},{"id":"2","name":"bob"}],"version":"v2"}

=== Nested groups (admin -> reports) ===
GET /admin/reports/summary -> status=200 body={"report":"summary","section":"admin/reports"}
```

**Learning Objectives:**
- ✅ Create a route group with a shared path prefix
- ✅ Run two versions of the same resource side by side
- ✅ Nest a group inside another group

---

## Exercise 4: Query Parameters

**Objective:** Read optional query-string values with `c.Query`, `c.DefaultQuery`, and `c.GetQuery`

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level28-exercise4
cd ~/projects/level28-exercise4
go mod init level28.example/exercise4
go get github.com/gin-gonic/gin@v1.12.0
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "io"
    "net/http/httptest"

    "github.com/gin-gonic/gin"
)

func call(router *gin.Engine, method, path string) {
    w := httptest.NewRecorder()
    req := httptest.NewRequest(method, path, nil)
    router.ServeHTTP(w, req)
    fmt.Printf("%s %s -> status=%d body=%s\n", method, path, w.Code, w.Body.String())
}

func main() {
    gin.SetMode(gin.TestMode)
    gin.DefaultWriter = io.Discard

    router := gin.New()

    router.GET("/search", func(c *gin.Context) {
        q := c.Query("q")
        limit := c.DefaultQuery("limit", "10")
        c.JSON(200, gin.H{"query": q, "limit": limit})
    })

    router.GET("/search-checked", func(c *gin.Context) {
        q, ok := c.GetQuery("q")
        c.JSON(200, gin.H{"query": q, "provided": ok})
    })

    fmt.Println("=== Query params present ===")
    call(router, "GET", "/search?q=golang&limit=5")

    fmt.Println("\n=== Missing query param falls back to default ===")
    call(router, "GET", "/search?q=golang")

    fmt.Println("\n=== No query params at all ===")
    call(router, "GET", "/search")

    fmt.Println("\n=== GetQuery reports whether the key was present ===")
    call(router, "GET", "/search-checked?q=gin")
    call(router, "GET", "/search-checked")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Query params present ===
GET /search?q=golang&limit=5 -> status=200 body={"limit":"5","query":"golang"}

=== Missing query param falls back to default ===
GET /search?q=golang -> status=200 body={"limit":"10","query":"golang"}

=== No query params at all ===
GET /search -> status=200 body={"limit":"10","query":""}

=== GetQuery reports whether the key was present ===
GET /search-checked?q=gin -> status=200 body={"provided":true,"query":"gin"}
GET /search-checked -> status=200 body={"provided":false,"query":""}
```

**Learning Objectives:**
- ✅ Read a query parameter with `c.Query`
- ✅ Provide a fallback value with `c.DefaultQuery`
- ✅ Distinguish "missing" from "empty" with `c.GetQuery`

---

## Exercise 5: JSON Binding With Validation

**Objective:** Decode and validate a JSON request body in one call with `c.ShouldBindJSON`

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level28-exercise5
cd ~/projects/level28-exercise5
go mod init level28.example/exercise5
go get github.com/gin-gonic/gin@v1.12.0
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "bytes"
    "fmt"
    "io"
    "net/http/httptest"

    "github.com/gin-gonic/gin"
)

type CreateUserInput struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age" binding:"required,gte=0,lte=130"`
}

func postJSON(router *gin.Engine, path string, body string) {
    w := httptest.NewRecorder()
    req := httptest.NewRequest("POST", path, bytes.NewBufferString(body))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)
    fmt.Printf("POST %s body=%s\n  -> status=%d resp=%s\n", path, body, w.Code, w.Body.String())
}

func main() {
    gin.SetMode(gin.TestMode)
    gin.DefaultWriter = io.Discard

    router := gin.New()
    router.POST("/users", func(c *gin.Context) {
        var input CreateUserInput
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }
        c.JSON(201, gin.H{"created": input})
    })

    fmt.Println("=== Valid request body ===")
    postJSON(router, "/users", `{"name":"Alice","email":"alice@example.com","age":30}`)

    fmt.Println("\n=== Missing required field (name) ===")
    postJSON(router, "/users", `{"email":"bob@example.com","age":25}`)

    fmt.Println("\n=== Invalid email format ===")
    postJSON(router, "/users", `{"name":"Carol","email":"not-an-email","age":25}`)

    fmt.Println("\n=== Age out of range ===")
    postJSON(router, "/users", `{"name":"Dave","email":"dave@example.com","age":200}`)

    fmt.Println("\n=== Malformed JSON ===")
    postJSON(router, "/users", `{"name":"Eve",`)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Valid request body ===
POST /users body={"name":"Alice","email":"alice@example.com","age":30}
  -> status=201 resp={"created":{"name":"Alice","email":"alice@example.com","age":30}}

=== Missing required field (name) ===
POST /users body={"email":"bob@example.com","age":25}
  -> status=400 resp={"error":"Key: 'CreateUserInput.Name' Error:Field validation for 'Name' failed on the 'required' tag"}

=== Invalid email format ===
POST /users body={"name":"Carol","email":"not-an-email","age":25}
  -> status=400 resp={"error":"Key: 'CreateUserInput.Email' Error:Field validation for 'Email' failed on the 'email' tag"}

=== Age out of range ===
POST /users body={"name":"Dave","email":"dave@example.com","age":200}
  -> status=400 resp={"error":"Key: 'CreateUserInput.Age' Error:Field validation for 'Age' failed on the 'lte' tag"}

=== Malformed JSON ===
POST /users body={"name":"Eve",
  -> status=400 resp={"error":"unexpected EOF"}
```

**Learning Objectives:**
- ✅ Decode a JSON body into a struct with `c.ShouldBindJSON`
- ✅ Declare validation rules with `binding` struct tags
- ✅ Read the exact validator error message for a failed rule
- ✅ Handle malformed JSON as just another bind error

---

## Exercise 6: Custom Middleware

**Objective:** Write your own `gin.HandlerFunc` middleware, both global and per-route

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level28-exercise6
cd ~/projects/level28-exercise6
go mod init level28.example/exercise6
go get github.com/gin-gonic/gin@v1.12.0
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "io"
    "net/http/httptest"

    "github.com/gin-gonic/gin"
)

// timingMiddleware is a custom global middleware. Code before c.Next()
// runs on the way IN; code after c.Next() runs on the way OUT, once the
// handler (and every middleware after this one) has finished.
func timingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Printf("  [timing] before handler: %s %s\n", c.Request.Method, c.Request.URL.Path)

        c.Next() // hand off to the next middleware/handler in the chain

        fmt.Printf("  [timing] after handler: status=%d\n", c.Writer.Status())
    }
}

// requireHeaderMiddleware only runs on the routes it's attached to.
func requireHeaderMiddleware(header string) gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.GetHeader(header) == "" {
            fmt.Printf("  [auth] missing header %q - aborting chain\n", header)
            c.AbortWithStatusJSON(401, gin.H{"error": "missing " + header})
            return
        }
        fmt.Printf("  [auth] header %q present - continuing\n", header)
        c.Next()
    }
}

func call(router *gin.Engine, method, path, headerKey, headerVal string) {
    w := httptest.NewRecorder()
    req := httptest.NewRequest(method, path, nil)
    if headerKey != "" {
        req.Header.Set(headerKey, headerVal)
    }
    router.ServeHTTP(w, req)
    fmt.Printf("%s %s -> status=%d body=%s\n", method, path, w.Code, w.Body.String())
}

func main() {
    gin.SetMode(gin.TestMode)
    gin.DefaultWriter = io.Discard

    router := gin.New()
    router.Use(timingMiddleware()) // global: runs for every route

    router.GET("/public", func(c *gin.Context) {
        c.JSON(200, gin.H{"area": "public"})
    })

    router.GET("/private", requireHeaderMiddleware("X-API-Key"), func(c *gin.Context) {
        c.JSON(200, gin.H{"area": "private"})
    })

    fmt.Println("=== Global middleware runs on every route ===")
    call(router, "GET", "/public", "", "")

    fmt.Println("\n=== Per-route middleware allows the request through ===")
    call(router, "GET", "/private", "X-API-Key", "secret123")

    fmt.Println("\n=== Per-route middleware blocks the request (c.Abort...) ===")
    call(router, "GET", "/private", "", "")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Global middleware runs on every route ===
  [timing] before handler: GET /public
  [timing] after handler: status=200
GET /public -> status=200 body={"area":"public"}

=== Per-route middleware allows the request through ===
  [timing] before handler: GET /private
  [auth] header "X-API-Key" present - continuing
  [timing] after handler: status=200
GET /private -> status=200 body={"area":"private"}

=== Per-route middleware blocks the request (c.Abort...) ===
  [timing] before handler: GET /private
  [auth] missing header "X-API-Key" - aborting chain
  [timing] after handler: status=401
GET /private -> status=401 body={"error":"missing X-API-Key"}
```

**Learning Objectives:**
- ✅ Write a custom `gin.HandlerFunc` middleware using `c.Next()`
- ✅ Attach middleware globally with `router.Use()`
- ✅ Attach middleware to a single route
- ✅ Stop the chain correctly with `c.AbortWithStatusJSON`

---

## Exercise 7: Panic Recovery Demo

**Objective:** See Gin's built-in `Recovery()` middleware turn a panic into a clean 500, not a crash

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level28-exercise7
cd ~/projects/level28-exercise7
go mod init level28.example/exercise7
go get github.com/gin-gonic/gin@v1.12.0
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "io"
    "net/http/httptest"

    "github.com/gin-gonic/gin"
)

func call(router *gin.Engine, method, path string) {
    w := httptest.NewRecorder()
    req := httptest.NewRequest(method, path, nil)
    router.ServeHTTP(w, req)
    fmt.Printf("%s %s -> status=%d body=%s\n", method, path, w.Code, w.Body.String())
}

func main() {
    gin.SetMode(gin.TestMode)
    gin.DefaultWriter = io.Discard
    gin.DefaultErrorWriter = io.Discard // silence the recovered-panic stack trace for clean output

    // gin.Default() includes Recovery() out of the box.
    router := gin.Default()

    router.GET("/divide/:a/:b", func(c *gin.Context) {
        var a, b int
        fmt.Sscanf(c.Param("a"), "%d", &a)
        fmt.Sscanf(c.Param("b"), "%d", &b)
        result := a / b // panics on division by zero
        c.JSON(200, gin.H{"result": result})
    })

    router.GET("/healthy", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    fmt.Println("=== A normal request works fine ===")
    call(router, "GET", "/divide/10/2")

    fmt.Println("\n=== A panicking handler is recovered, not crashed ===")
    call(router, "GET", "/divide/10/0")

    fmt.Println("\n=== The server is still alive for the next request ===")
    call(router, "GET", "/healthy")

    fmt.Println("\nProgram reached the end normally - no crash, no exit status 2.")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== A normal request works fine ===
GET /divide/10/2 -> status=200 body={"result":5}

=== A panicking handler is recovered, not crashed ===
GET /divide/10/0 -> status=500 body=

=== The server is still alive for the next request ===
GET /healthy -> status=200 body={"status":"ok"}

Program reached the end normally - no crash, no exit status 2.
```

**Learning Objectives:**
- ✅ See a real `panic` (integer division by zero) happen inside a handler
- ✅ Confirm `Recovery()` converts it into a `500` response instead of crashing the process
- ✅ Confirm the server keeps serving later requests after a recovered panic
- ✅ Contrast this with Level 16's unrecovered-panic behavior (crash, exit status 2)

---

## Exercise 8: Full Todos API Rewrite

**Objective:** Rewrite a complete CRUD resource using Gin's routing, binding, and grouping

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level28-exercise8
cd ~/projects/level28-exercise8
go mod init level28.example/exercise8
go get github.com/gin-gonic/gin@v1.12.0
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "bytes"
    "fmt"
    "io"
    "net/http/httptest"
    "strconv"

    "github.com/gin-gonic/gin"
)

type Todo struct {
    ID    int    `json:"id"`
    Title string `json:"title"`
    Done  bool   `json:"done"`
}

type CreateTodoInput struct {
    Title string `json:"title" binding:"required"`
}

type UpdateTodoInput struct {
    Title string `json:"title" binding:"required"`
    Done  bool   `json:"done"`
}

// In-memory store. A real service would use Level 29's database/sql layer instead.
var (
    todos  = map[int]Todo{}
    nextID = 1
)

func listTodos(c *gin.Context) {
    result := make([]Todo, 0, len(todos))
    for id := 1; id < nextID; id++ {
        if t, ok := todos[id]; ok {
            result = append(result, t)
        }
    }
    c.JSON(200, gin.H{"todos": result})
}

func getTodo(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    t, ok := todos[id]
    if !ok {
        c.JSON(404, gin.H{"error": "todo not found"})
        return
    }
    c.JSON(200, t)
}

func createTodo(c *gin.Context) {
    var input CreateTodoInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    t := Todo{ID: nextID, Title: input.Title, Done: false}
    todos[t.ID] = t
    nextID++
    c.JSON(201, t)
}

func updateTodo(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    if _, ok := todos[id]; !ok {
        c.JSON(404, gin.H{"error": "todo not found"})
        return
    }
    var input UpdateTodoInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    t := Todo{ID: id, Title: input.Title, Done: input.Done}
    todos[id] = t
    c.JSON(200, t)
}

func deleteTodo(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    if _, ok := todos[id]; !ok {
        c.JSON(404, gin.H{"error": "todo not found"})
        return
    }
    delete(todos, id)
    c.Status(204)
}

func setupRouter() *gin.Engine {
    router := gin.New()
    router.Use(gin.Recovery())

    api := router.Group("/api/v1")
    {
        api.GET("/todos", listTodos)
        api.GET("/todos/:id", getTodo)
        api.POST("/todos", createTodo)
        api.PUT("/todos/:id", updateTodo)
        api.DELETE("/todos/:id", deleteTodo)
    }
    return router
}

func call(router *gin.Engine, method, path, body string) {
    w := httptest.NewRecorder()
    req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)
    if body != "" {
        fmt.Printf("%s %s body=%s\n  -> status=%d resp=%s\n", method, path, body, w.Code, w.Body.String())
    } else {
        fmt.Printf("%s %s -> status=%d resp=%s\n", method, path, w.Code, w.Body.String())
    }
}

func main() {
    gin.SetMode(gin.TestMode)
    gin.DefaultWriter = io.Discard

    router := setupRouter()

    fmt.Println("=== Create two todos ===")
    call(router, "POST", "/api/v1/todos", `{"title":"Learn Gin"}`)
    call(router, "POST", "/api/v1/todos", `{"title":"Rewrite the todos API"}`)

    fmt.Println("\n=== List all todos ===")
    call(router, "GET", "/api/v1/todos", "")

    fmt.Println("\n=== Get a single todo ===")
    call(router, "GET", "/api/v1/todos/1", "")

    fmt.Println("\n=== Get a todo that doesn't exist ===")
    call(router, "GET", "/api/v1/todos/99", "")

    fmt.Println("\n=== Update a todo ===")
    call(router, "PUT", "/api/v1/todos/1", `{"title":"Learn Gin","done":true}`)

    fmt.Println("\n=== Create with missing required field ===")
    call(router, "POST", "/api/v1/todos", `{}`)

    fmt.Println("\n=== Delete a todo ===")
    call(router, "DELETE", "/api/v1/todos/2", "")

    fmt.Println("\n=== List after delete ===")
    call(router, "GET", "/api/v1/todos", "")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Create two todos ===
POST /api/v1/todos body={"title":"Learn Gin"}
  -> status=201 resp={"id":1,"title":"Learn Gin","done":false}
POST /api/v1/todos body={"title":"Rewrite the todos API"}
  -> status=201 resp={"id":2,"title":"Rewrite the todos API","done":false}

=== List all todos ===
GET /api/v1/todos -> status=200 resp={"todos":[{"id":1,"title":"Learn Gin","done":false},{"id":2,"title":"Rewrite the todos API","done":false}]}

=== Get a single todo ===
GET /api/v1/todos/1 -> status=200 resp={"id":1,"title":"Learn Gin","done":false}

=== Get a todo that doesn't exist ===
GET /api/v1/todos/99 -> status=404 resp={"error":"todo not found"}

=== Update a todo ===
PUT /api/v1/todos/1 body={"title":"Learn Gin","done":true}
  -> status=200 resp={"id":1,"title":"Learn Gin","done":true}

=== Create with missing required field ===
POST /api/v1/todos body={}
  -> status=400 resp={"error":"Key: 'CreateTodoInput.Title' Error:Field validation for 'Title' failed on the 'required' tag"}

=== Delete a todo ===
DELETE /api/v1/todos/2 -> status=204 resp=

=== List after delete ===
GET /api/v1/todos -> status=200 resp={"todos":[{"id":1,"title":"Learn Gin","done":true}]}
```

**Learning Objectives:**
- ✅ Build a full CRUD resource (Create, Read, Update, Delete) with Gin
- ✅ Combine route groups, path parameters, and JSON binding in one resource
- ✅ Return the right status code for each outcome (200, 201, 204, 400, 404)

---

## Exercise 9: Testing Gin Handlers With httptest

**Objective:** Write real `go test` tests against a Gin router, sharing one `SetupRouter` function

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level28-exercise9
cd ~/projects/level28-exercise9
go mod init level28.example/exercise9
go get github.com/gin-gonic/gin@v1.12.0
```

2. Create `router.go`:

```bash
cat > router.go << 'EOF'
package app

import "github.com/gin-gonic/gin"

// SetupRouter builds the Gin engine so both main() and the tests can share
// the exact same routing setup.
func SetupRouter() *gin.Engine {
    router := gin.New()

    router.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })

    router.GET("/greet/:name", func(c *gin.Context) {
        name := c.Param("name")
        c.JSON(200, gin.H{"greeting": "Hello, " + name + "!"})
    })

    return router
}
EOF
```

3. Create `router_test.go`:

```bash
cat > router_test.go << 'EOF'
package app

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
)

func TestPingRoute(t *testing.T) {
    gin.SetMode(gin.TestMode)
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

func TestGreetRoute(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := SetupRouter()

    tests := []struct {
        name string
        path string
        want string
    }{
        {"simple name", "/greet/Alice", `{"greeting":"Hello, Alice!"}`},
        {"another name", "/greet/Gopher", `{"greeting":"Hello, Gopher!"}`},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest(http.MethodGet, tt.path, nil)
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            if w.Code != http.StatusOK {
                t.Errorf("status = %d; want %d", w.Code, http.StatusOK)
            }
            if got := w.Body.String(); got != tt.want {
                t.Errorf("body = %q; want %q", got, tt.want)
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

**Expected Output** *(the final timing, `0.765s`, is machine-dependent - yours will differ; everything above the `ok` line should match)*:

```
=== RUN   TestPingRoute
--- PASS: TestPingRoute (0.00s)
=== RUN   TestGreetRoute
=== RUN   TestGreetRoute/simple_name
=== RUN   TestGreetRoute/another_name
--- PASS: TestGreetRoute (0.00s)
    --- PASS: TestGreetRoute/simple_name (0.00s)
    --- PASS: TestGreetRoute/another_name (0.00s)
PASS
ok  	level28.example/exercise9	0.765s
```

**Learning Objectives:**
- ✅ Put routing setup in an exported `SetupRouter` function shared by `main()` and tests
- ✅ Test a Gin router with `httptest.NewRequest` + `httptest.NewRecorder` + `router.ServeHTTP`
- ✅ Combine table-driven tests (Level 26) with Gin routes

---

## Exercise 10: Comprehensive Practice - A Bookstore API

**Objective:** Combine route groups, JSON validation, custom auth middleware, and httptest-based tests in one realistic API

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level28-exercise10
cd ~/projects/level28-exercise10
go mod init level28.example/exercise10
go get github.com/gin-gonic/gin@v1.12.0
```

2. Create `bookstore.go`:

```bash
cat > bookstore.go << 'EOF'
package bookstore

import (
    "strconv"
    "strings"

    "github.com/gin-gonic/gin"
)

type Book struct {
    ID     int     `json:"id"`
    Title  string  `json:"title"`
    Author string  `json:"author"`
    Price  float64 `json:"price"`
}

type CreateBookInput struct {
    Title  string  `json:"title" binding:"required"`
    Author string  `json:"author" binding:"required"`
    Price  float64 `json:"price" binding:"required,gt=0"`
}

// authMiddleware requires "Authorization: Bearer <apiKey>" on every route
// it's attached to. It's per-group, not global - public routes never see it.
func authMiddleware(apiKey string) gin.HandlerFunc {
    return func(c *gin.Context) {
        want := "Bearer " + apiKey
        if c.GetHeader("Authorization") != want {
            c.AbortWithStatusJSON(401, gin.H{"error": "missing or invalid Authorization header"})
            return
        }
        c.Next()
    }
}

// SetupRouter builds a fresh Gin engine with its own in-memory book store,
// so every call (including every test) starts from a clean slate.
func SetupRouter(apiKey string) *gin.Engine {
    books := map[int]Book{}
    nextID := 1

    router := gin.New()
    router.Use(gin.Recovery())

    v1 := router.Group("/api/v1")
    {
        v1.GET("/books", func(c *gin.Context) {
            result := make([]Book, 0, len(books))
            for id := 1; id < nextID; id++ {
                if b, ok := books[id]; ok {
                    result = append(result, b)
                }
            }
            c.JSON(200, gin.H{"books": result})
        })

        v1.GET("/books/:id", func(c *gin.Context) {
            id, err := strconv.Atoi(c.Param("id"))
            if err != nil {
                c.JSON(400, gin.H{"error": "invalid book id"})
                return
            }
            b, ok := books[id]
            if !ok {
                c.JSON(404, gin.H{"error": "book not found"})
                return
            }
            c.JSON(200, b)
        })
    }

    protected := router.Group("/api/v1")
    protected.Use(authMiddleware(apiKey))
    {
        protected.POST("/books", func(c *gin.Context) {
            var input CreateBookInput
            if err := c.ShouldBindJSON(&input); err != nil {
                c.JSON(400, gin.H{"error": err.Error()})
                return
            }
            b := Book{ID: nextID, Title: input.Title, Author: input.Author, Price: input.Price}
            books[b.ID] = b
            nextID++
            c.JSON(201, b)
        })

        protected.DELETE("/books/:id", func(c *gin.Context) {
            id, err := strconv.Atoi(c.Param("id"))
            if err != nil {
                c.JSON(400, gin.H{"error": "invalid book id"})
                return
            }
            if _, ok := books[id]; !ok {
                c.JSON(404, gin.H{"error": "book not found"})
                return
            }
            delete(books, id)
            c.Status(204)
        })
    }

    router.NoRoute(func(c *gin.Context) {
        c.JSON(404, gin.H{"error": "route not found: " + strings.ToUpper(c.Request.Method) + " " + c.Request.URL.Path})
    })

    return router
}
EOF
```

3. Create `bookstore_test.go`:

```bash
cat > bookstore_test.go << 'EOF'
package bookstore

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/gin-gonic/gin"
)

const testAPIKey = "secret123"

func doRequest(router *gin.Engine, method, path, body, authHeader string) *httptest.ResponseRecorder {
    req := httptest.NewRequest(method, path, strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    if authHeader != "" {
        req.Header.Set("Authorization", authHeader)
    }
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    return w
}

func TestPublicRoutesNoAuthRequired(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := SetupRouter(testAPIKey)

    w := doRequest(router, http.MethodGet, "/api/v1/books", "", "")
    if w.Code != http.StatusOK {
        t.Fatalf("status = %d; want %d", w.Code, http.StatusOK)
    }
    want := `{"books":[]}`
    if got := w.Body.String(); got != want {
        t.Errorf("body = %q; want %q", got, want)
    }
}

func TestCreateBookRequiresAuth(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := SetupRouter(testAPIKey)

    body := `{"title":"The Go Programming Language","author":"Donovan & Kernighan","price":39.99}`

    tests := []struct {
        name       string
        authHeader string
    }{
        {"no auth header at all", ""},
        {"wrong api key", "Bearer wrong-key"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := doRequest(router, http.MethodPost, "/api/v1/books", body, tt.authHeader)
            if w.Code != http.StatusUnauthorized {
                t.Errorf("status = %d; want %d", w.Code, http.StatusUnauthorized)
            }
        })
    }
}

func TestCreateBookValidation(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := SetupRouter(testAPIKey)
    auth := "Bearer " + testAPIKey

    tests := []struct {
        name string
        body string
    }{
        {"missing title", `{"author":"Someone","price":10.0}`},
        {"missing author", `{"title":"Some Book","price":10.0}`},
        {"zero price", `{"title":"Some Book","author":"Someone","price":0}`},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := doRequest(router, http.MethodPost, "/api/v1/books", tt.body, auth)
            if w.Code != http.StatusBadRequest {
                t.Errorf("status = %d; want %d", w.Code, http.StatusBadRequest)
            }
        })
    }
}

func TestFullBookLifecycle(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := SetupRouter(testAPIKey)
    auth := "Bearer " + testAPIKey

    // Create
    createBody := `{"title":"The Go Programming Language","author":"Donovan and Kernighan","price":39.99}`
    w := doRequest(router, http.MethodPost, "/api/v1/books", createBody, auth)
    if w.Code != http.StatusCreated {
        t.Fatalf("create: status = %d; want %d, body=%s", w.Code, http.StatusCreated, w.Body.String())
    }

    // List shows the new book
    w = doRequest(router, http.MethodGet, "/api/v1/books", "", "")
    if w.Code != http.StatusOK {
        t.Fatalf("list: status = %d; want %d", w.Code, http.StatusOK)
    }
    wantList := `{"books":[{"id":1,"title":"The Go Programming Language","author":"Donovan and Kernighan","price":39.99}]}`
    if got := w.Body.String(); got != wantList {
        t.Errorf("list body = %q; want %q", got, wantList)
    }

    // Get by id
    w = doRequest(router, http.MethodGet, "/api/v1/books/1", "", "")
    if w.Code != http.StatusOK {
        t.Fatalf("get: status = %d; want %d", w.Code, http.StatusOK)
    }

    // Delete without auth fails
    w = doRequest(router, http.MethodDelete, "/api/v1/books/1", "", "")
    if w.Code != http.StatusUnauthorized {
        t.Errorf("delete without auth: status = %d; want %d", w.Code, http.StatusUnauthorized)
    }

    // Delete with auth succeeds
    w = doRequest(router, http.MethodDelete, "/api/v1/books/1", "", auth)
    if w.Code != http.StatusNoContent {
        t.Errorf("delete with auth: status = %d; want %d", w.Code, http.StatusNoContent)
    }

    // Now it's gone
    w = doRequest(router, http.MethodGet, "/api/v1/books/1", "", "")
    if w.Code != http.StatusNotFound {
        t.Errorf("get after delete: status = %d; want %d", w.Code, http.StatusNotFound)
    }
}
EOF
```

4. Run the tests:

```bash
go test -v ./...
```

**Expected Output** *(the final timing, `0.816s`, is machine-dependent - yours will differ; everything above the `ok` line should match)*:

```
=== RUN   TestPublicRoutesNoAuthRequired
--- PASS: TestPublicRoutesNoAuthRequired (0.00s)
=== RUN   TestCreateBookRequiresAuth
=== RUN   TestCreateBookRequiresAuth/no_auth_header_at_all
=== RUN   TestCreateBookRequiresAuth/wrong_api_key
--- PASS: TestCreateBookRequiresAuth (0.00s)
    --- PASS: TestCreateBookRequiresAuth/no_auth_header_at_all (0.00s)
    --- PASS: TestCreateBookRequiresAuth/wrong_api_key (0.00s)
=== RUN   TestCreateBookValidation
=== RUN   TestCreateBookValidation/missing_title
=== RUN   TestCreateBookValidation/missing_author
=== RUN   TestCreateBookValidation/zero_price
--- PASS: TestCreateBookValidation (0.00s)
    --- PASS: TestCreateBookValidation/missing_title (0.00s)
    --- PASS: TestCreateBookValidation/missing_author (0.00s)
    --- PASS: TestCreateBookValidation/zero_price (0.00s)
=== RUN   TestFullBookLifecycle
--- PASS: TestFullBookLifecycle (0.00s)
PASS
ok  	level28.example/exercise10	0.816s
```

**Learning Objectives:**
- ✅ Combine public and auth-protected route groups on the same engine
- ✅ Apply custom middleware to only a subset of routes via a group
- ✅ Validate JSON bodies with multiple `binding` rules (`required`, `gt=0`)
- ✅ Write a full request/response lifecycle test, not just single-call tests
- ✅ Give every group its own fresh state so tests never depend on execution order

---

## Bonus Challenges

### Challenge 1: Request-ID Middleware

Write a middleware that generates a short random ID for every incoming request, stores it in the context with `c.Set`, and also sets it as a response header so callers can correlate a request with server-side logs.

```bash
mkdir -p ~/projects/level28-bonus1
cd ~/projects/level28-bonus1
go mod init level28.example/bonus1
go get github.com/gin-gonic/gin@v1.12.0
```

**Hints:**
- Generate the ID with `crypto/rand` + `encoding/hex` (8 random bytes hex-encoded is plenty for a demo)
- Store it with `c.Set("request_id", id)`, then read it later in a handler with `c.MustGet("request_id").(string)`
- Set the response header with `c.Header("X-Request-ID", id)` **before** calling `c.Next()` so it's set even if a later handler panics and gets recovered

### Challenge 2: A Route Group With Its Own Middleware Stack

Build a router with three groups: `/public` (no middleware), `/admin` (a custom `requireRole("admin")` middleware), and `/api/v1` (a request-logging middleware only). Prove all three behave independently.

```bash
mkdir -p ~/projects/level28-bonus2
cd ~/projects/level28-bonus2
go mod init level28.example/bonus2
go get github.com/gin-gonic/gin@v1.12.0
```

**Hints:**
- Each `router.Group(...)` returns its own `*gin.RouterGroup` - calling `.Use()` on one never affects another group or the base engine
- Attach group middleware with `group.Use(...)` **before** registering routes on that group (see Common Mistake 4 in the README)
- Test each group with a request that should pass and one that should be rejected

### Challenge 3: Convert a Level 27 stdlib Handler to Gin

Take a `net/http`-style handler from Level 27 (signature `func(w http.ResponseWriter, r *http.Request)`) and rewrite it as a Gin handler (`func(c *gin.Context)`), side by side, to see the boilerplate Gin removes.

```bash
mkdir -p ~/projects/level28-bonus3
cd ~/projects/level28-bonus3
go mod init level28.example/bonus3
go get github.com/gin-gonic/gin@v1.12.0
```

**Hints:**
- `w.Header().Set("Content-Type", "application/json")` + `json.NewEncoder(w).Encode(v)` becomes one call: `c.JSON(status, v)`
- `r.PathValue("id")` (Go 1.22+ `ServeMux`) becomes `c.Param("id")`
- `r.URL.Query().Get("q")` becomes `c.Query("q")`
- `json.NewDecoder(r.Body).Decode(&input)` plus manual validation becomes `c.ShouldBindJSON(&input)` with `binding` tags
- Write both versions and confirm they produce identical responses for the same request

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Set up a Gin engine with `gin.Default()` or `gin.New()`, and explain the difference
✅ Route requests with path parameters, wildcards, and route groups
✅ Read query parameters, path parameters, and validated JSON request bodies
✅ Write JSON, string, and status-only responses
✅ Write and attach custom middleware, both global and scoped to a route or group
✅ Explain how Gin's `Recovery()` middleware turns a panic into a clean response
✅ Build a complete CRUD REST resource with Gin
✅ Test a Gin router with `httptest`, sharing one `SetupRouter` function between `main()` and tests
✅ Decide when Gin is worth the dependency versus plain `net/http`

---

## Next Level

Level 29: Database/SQL
- Connecting to a real database from Go
- `database/sql`, connection pools, and drivers
- Queries, scanning rows into structs, and prepared statements
- Wiring a database layer underneath the API you just built with Gin

Great work! You're building real, production-shaped APIs! 🚀
