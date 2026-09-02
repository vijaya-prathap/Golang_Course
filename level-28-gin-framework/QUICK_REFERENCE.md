# Level 28: Quick Reference Card

> ⚠️ **Requires internet access** - `go get github.com/gin-gonic/gin@v1.12.0` (the version this course was verified against). This is the one level in the course that isn't fully offline.

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir myapp && cd myapp
go mod init myapp
go get github.com/gin-gonic/gin@v1.12.0

# Create main.go
cat > main.go << 'EOF'
package main

import "github.com/gin-gonic/gin"

func main() {
    router := gin.Default()

    router.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })

    router.GET("/users/:id", func(c *gin.Context) {
        c.JSON(200, gin.H{"user_id": c.Param("id")})
    })

    router.Run(":8080") // blocks, serving on http://localhost:8080
}
EOF

# Run (then curl http://localhost:8080/ping from another terminal)
go run main.go
```

---

## 📋 Setup

```go
router := gin.Default() // gin.New() + Logger() + Recovery()
router := gin.New()     // bare engine, no middleware
```

---

## 🗺️ Routing

```go
router.GET("/path", handler)
router.POST("/path", handler)
router.PUT("/path", handler)
router.DELETE("/path", handler)

router.GET("/users/:id", handler)             // path param -> c.Param("id")
router.GET("/users/:id/posts/:postId", h)     // multiple params
router.GET("/files/*filepath", handler)       // wildcard, includes leading "/"

v1 := router.Group("/api/v1")
v1.GET("/users", handler)                     // -> /api/v1/users
```

---

## 📥 Reading Input

```go
c.Param("id")                    // path parameter, always a string
c.Query("q")                     // query param, "" if missing
c.DefaultQuery("limit", "10")    // query param with fallback
q, ok := c.GetQuery("q")         // value + whether it was present
c.GetHeader("X-API-Key")         // request header

type Input struct {
    Name string `json:"name" binding:"required"`
    Age  int    `json:"age" binding:"required,gte=0"`
}
var input Input
if err := c.ShouldBindJSON(&input); err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
}
```

---

## 📤 Writing Responses

```go
c.JSON(200, gin.H{"key": "value"}) // JSON body
c.String(200, "plain text")        // text/plain
c.Status(204)                      // status only, no body
```

`gin.H` is `map[string]any` - keys serialize alphabetically. Use a struct instead when field order matters.

---

## 🔗 Middleware

```go
func myMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // ...before...
        c.Next()
        // ...after...
    }
}

router.Use(myMiddleware())                       // global
router.GET("/x", myMiddleware(), handler)         // per-route
group.Use(myMiddleware())                         // per-group (before registering routes!)

c.AbortWithStatusJSON(401, gin.H{"error": "..."}) // stop the chain + write an error
```

---

## 🛟 Recovery

```go
router := gin.Default() // Recovery() included - a panicking handler
                         // returns 500 instead of crashing the process
```

---

## 🧪 Testing

```go
gin.SetMode(gin.TestMode)
router := SetupRouter()

req := httptest.NewRequest(http.MethodGet, "/ping", nil)
w := httptest.NewRecorder()
router.ServeHTTP(w, req)

// assert on w.Code and w.Body.String()
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Missing `binding:"required"` | field silently zero-valued | tag every required field |
| Double-binding a body | `ShouldBindJSON` called twice | bind once (or buffer with `ShouldBindBodyWith`) |
| Ignoring bind errors | `c.ShouldBindJSON(&x)` then using `x` | check the returned `error` first |
| Group `Use()` too late | `Use()` after routes registered | `Use()` before registering routes |
| Forgetting `c.Next()` | middleware never continues the chain | always call `c.Next()` unless intentionally aborting |
| `return` instead of `Abort` | handler still runs after "returning" | call `c.AbortWithStatus...` then `return` |

---

## 🎓 Before Next Level

Can you:
- [ ] Explain `gin.Default()` vs `gin.New()`?
- [ ] Register a route with a path parameter and read it?
- [ ] Validate a JSON body with `binding` struct tags?
- [ ] Write a custom middleware using `c.Next()`?
- [ ] Explain how `Recovery()` prevents a crash?
- [ ] Test a Gin router with `httptest`?

If YES → You're ready for Level 29!

---

## 📚 Next Level

Level 29: Database/SQL
- `database/sql`, connection pools, and drivers
- Queries, scanning rows into structs, prepared statements

You've got Gin down! 💪
