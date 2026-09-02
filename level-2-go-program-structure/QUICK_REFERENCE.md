# Level 2: Quick Reference Card

## 🚀 Quick Start (10 Minutes)

```bash
# Setup project
mkdir myapp && cd myapp
go mod init myapp.com/app

# Create package
mkdir helpers
cat > helpers/strings.go << 'EOF'
package helpers

func Reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}
EOF

# Use package
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "myapp.com/app/helpers"
)

func main() {
    fmt.Println(helpers.Reverse("hello"))
}
EOF

# Run
go run .
```

---

## 📋 Essential Concepts

### Package Structure
```
One directory = One package
One package = Can have multiple files
Package name ≠ Directory name (but usually same)
```

### Visibility Rules
```
UPPERCASE = Exported (public)
lowercase = Unexported (private)
```

### Initialization Order
```
1. Constants
2. Variables
3. init() functions
4. main() function
```

---

## 🔧 Common Tasks

### Create Package
```bash
mkdir packagename
cat > packagename/file.go << 'EOF'
package packagename

func DoSomething() string {
    return "something"
}
EOF
```

### Import Package
```go
import "myapp/packagename"

// Use:
packagename.DoSomething()
```

### Export vs Unexport
```go
func Public() { }       // ✅ Exported (uppercase)
func private() { }      // ❌ Unexported (lowercase)

type Public struct { }  // ✅ Exported
type private struct { } // ❌ Unexported
```

### Package-Level Variables
```go
var GlobalVar = 100     // Exported
var privateVar = 100    // Unexported

func GetVar() int {
    return privateVar   // Access within package
}
```

### Multiple init() Functions
```go
func init() {
    // Runs first
}

func init() {
    // Runs second
}
```

### Import Alias
```go
import (
    "fmt"
    m "myapp/models"    // Alias
)

// Use:
fmt.Println()
m.NewUser()
```

---

## 📁 Project Structure Template

```
myapp/
├── go.mod
├── main.go              (package main)
├── models/              (package models)
│   └── user.go
├── services/            (package services)
│   └── user_service.go
└── handlers/            (package handlers)
    └── handlers.go
```

---

## ⚡ Key Commands

```bash
# View documentation
go doc mypackage
go doc mypackage.Function

# Format all files
go fmt ./...

# Check for issues
go vet ./...

# Run all files
go run .

# Build
go build

# Build with name
go build -o myapp
```

---

## 🎯 Visibility Patterns

### Struct with Private Fields
```go
type User struct {
    ID   int       // Exported
    Name string    // Exported
    password string // Unexported
}

func (u User) GetPassword() string {
    return u.password  // Getter
}

func (u *User) SetPassword(p string) {
    u.password = p     // Setter
}
```

### Factory Pattern
```go
type config struct {  // Unexported
    debug bool
    timeout int
}

func NewConfig() *config {  // Exported factory
    return &config{
        debug:   false,
        timeout: 30,
    }
}
```

---

## ❌ Common Errors

| Error | Cause | Fix |
|-------|-------|-----|
| `package not found` | Wrong import path | Check module name in go.mod |
| `undefined: Function` | Unexported function | Use capital letter: Function() |
| `circular import` | Packages import each other | Restructure or use interface |
| `cannot refer to unexported` | Access private outside package | Only access exported identifiers |

---

## 📊 Multiple Files Same Package

Files in same package can:
- ✅ Call each other's unexported functions
- ✅ Share package-level variables
- ✅ Use same type definitions

```
math/
├── add.go    (func Add, func validate)
├── sub.go    (func Subtract, uses validate)
└── mul.go    (func Multiply)

All functions can call validate
Even though it's in add.go
```

---

## 🔄 Import Best Practices

### ✅ DO:
```go
import (
    "fmt"
    "myapp/models"
    "myapp/services"
)
```

### ❌ DON'T:
```go
import (
    "./models"     // Relative paths
)

import . "fmt"     // Dot imports
```

---

## 💾 Package-Level Variables

```go
// Access only within package
var internalState = 0

// But create exported getter
func GetState() int {
    return internalState
}

func SetState(s int) {
    internalState = s
}
```

---

## 🎓 Checklist Before Moving On

- [ ] Created 3+ packages
- [ ] Used import statements correctly
- [ ] Understand exported/unexported
- [ ] Know initialization order
- [ ] Used init() function
- [ ] Organized a small project
- [ ] Completed 8 exercises

---

## 📚 File Organization

### Single Responsibility
```
Each package = one responsibility
Each file = related functions/types
```

### Example Structure
```
models/          (Data structures)
├── user.go
└── product.go

services/        (Business logic)
├── user_service.go
└── product_service.go

handlers/        (Output)
└── handlers.go

main.go          (Entry point)
```

---

## 🚀 Next Level

Level 3: Variables & Data Types
- Variable declaration styles
- Basic types (int, float, string, bool)
- Type conversion
- Constants and iota
- Zero values

Keep learning! 🎉

---

## Remember

```
One directory = One package
Capitalization matters
init() runs automatically
Avoid circular imports
Organize by responsibility
```

You're building real Go skills now! 💪
