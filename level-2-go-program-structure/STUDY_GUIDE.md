# Level 2: Study Guide & Quick Reference

## 📚 Learning Path

### Week 1: Packages & Organization
```
Day 1:  Why packages? Package structure
Day 2:  Creating packages
Day 3:  Import statements
Day 4:  Exported vs unexported
Day 5:  Package-level variables
Day 6:  Init functions
Day 7:  Multi-file projects
```

### Week 2: Practice
```
Day 1:  Exercises 1-3 (Basic packages)
Day 2:  Exercises 4-6 (Variables & imports)
Day 3:  Exercises 7-8 (Real projects)
Day 4:  Exercise 9-10 (Build & docs)
Day 5:  Bonus challenges
Day 6-7: Practice & review
```

---

## 🎯 Package Architecture

### Project Structure (Typical)

```
myapp/
├── go.mod
├── main.go              (package main)
├── models/              (package models)
│   ├── user.go
│   └── product.go
├── services/            (package services)
│   ├── user_service.go
│   └── product_service.go
├── handlers/            (package handlers)
│   └── http_handlers.go
├── repository/          (package repository)
│   └── db.go
└── utils/               (package utils)
    └── helpers.go
```

### Directory = Package

```
One directory can have only ONE package
But one package can span MULTIPLE files
```

### Example

```
math/          ← Directory (one package)
├── add.go      (package math)
├── sub.go      (package math)
├── mul.go      (package math)
└── div.go      (package math)
```

All files declare: `package math`

---

## 🔍 Visibility Rules

### The Capitalization Rule

```
UPPERCASE = Exported (visible outside package)
lowercase = Unexported (hidden from other packages)
```

### At Package Level

```go
// ✅ EXPORTED
const MaxSize = 100
var GlobalCount int
type Config struct { }
func DoSomething() { }

// ❌ UNEXPORTED
const minSize = 10
var internalValue int
type config struct { }
func doSomething() { }
```

### Inside Structs

```go
type User struct {
    // ✅ Exported field
    Name string
    
    // ❌ Unexported field
    password string
}

// ✅ Exported method
func (u User) GetName() string { }

// ❌ Unexported method
func (u User) getPassword() string { }
```

---

## 📊 Import Patterns

### Standard Patterns

```go
// Single import
import "fmt"

// Multiple imports
import (
    "fmt"
    "os"
    "time"
)

// Local packages
import (
    "myapp/models"
    "myapp/services"
)

// With aliases
import (
    "fmt"
    m "myapp/models"
    s "myapp/services"
)
```

### Usage

```go
// Direct use
fmt.Println()

// With alias
m.User
s.CreateUser()
```

---

## 🚀 Quick Start: Create a Package

```bash
# 1. Setup
mkdir myapp && cd myapp
go mod init myapp.com/project

# 2. Create package directory
mkdir helpers

# 3. Create package file
cat > helpers/math.go << 'EOF'
package helpers

// Exported function
func Add(a, b int) int {
    return a + b
}

// Unexported function
func validate(x int) bool {
    return x > 0
}
EOF

# 4. Use in main
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "myapp.com/project/helpers"
)

func main() {
    result := helpers.Add(5, 3)
    fmt.Println(result)
}
EOF

# 5. Run
go run .
```

---

## 📋 Initialization Order

### When Programs Start

```
1. Import packages
2. Global constants initialized
3. Global variables initialized
4. init() functions run (all files, in order)
5. main() function runs
```

### Example

```go
package main

import "fmt"

var count = 0  // Step 2

func init() {  // Step 4
    fmt.Println("Init running")
    count = 10
}

func main() {  // Step 5
    fmt.Println(count)  // Prints: 10
}
```

---

## 🔄 Package Lifecycle

```
Import Package
    ↓
Initialize global constants
    ↓
Initialize global variables
    ↓
Run all init() functions
    ↓
Package ready to use
```

---

## 📁 File Organization

### How Many Files Per Package?

**Good reasons to split:**
- Different responsibilities
- Large file (> 500 lines)
- Related logic grouped
- Testing convenience

**Files in a package can:**
- ✅ Access each other's unexported functions
- ✅ Share package-level variables
- ✅ Call package-level functions

**Example:**

```
models/
├── user.go       (User, unexported helpers)
├── product.go    (Product, unexported helpers)
└── shared.go     (Shared constants, helpers)

main.go imports models package:
models.User
models.Product
```

---

## 🎯 Common Patterns

### Pattern 1: Service Layer

```
services/
├── user_service.go     (Business logic)
└── product_service.go

// Usage:
services.CreateUser()
services.GetProduct()
```

### Pattern 2: Model Layer

```
models/
├── user.go       (User struct)
├── product.go    (Product struct)
└── order.go      (Order struct)

// Usage:
models.NewUser()
models.NewProduct()
```

### Pattern 3: Handler Layer

```
handlers/
├── user_handlers.go       (User endpoints)
├── product_handlers.go    (Product endpoints)

// Usage:
handlers.GetUserHandler()
handlers.CreateProductHandler()
```

### Pattern 4: Repository Pattern

```
repository/
├── user_repository.go       (Database operations)

// Usage:
repository.SaveUser()
repository.GetUser()
```

---

## ⚠️ Common Mistakes

### Mistake 1: Wrong Package Name

```go
// ❌ WRONG
// Directory: helpers/
package utilities  // Doesn't match!

// ✅ RIGHT
// Directory: helpers/
package helpers    // Matches directory
```

### Mistake 2: Case Sensitivity

```go
// ❌ WRONG - will be unexported
func doSomething() { }

// ✅ RIGHT - will be exported
func DoSomething() { }
```

### Mistake 3: Importing Non-Existent Package

```go
// ❌ WRONG
import "myapp/models"    // Package doesn't exist

// ✅ RIGHT
import "myapp/model"     // Correct path
```

### Mistake 4: Circular Imports

```
models → imports database
database → imports models
❌ ERROR: Circular import
```

**Fix:** Use interfaces or extract common types

### Mistake 5: Access Unexported from Outside

```go
// In models/user.go
func (u User) getPassword() string { }  // lowercase = unexported

// In main.go
user.getPassword()  // ❌ Compiler error!
```

---

## 📊 Multiple Files Same Package

### How It Works

```
math/
├── add.go       (package math)
├── sub.go       (package math)
└── mul.go       (package math)

// All can access each other!
```

### Example

**add.go**
```go
package math

func Add(a, b int) int {
    return validateThenAdd(a, b)  // Call from sub.go!
}
```

**sub.go**
```go
package math

func validateThenAdd(a, b int) int {
    if validate(a, b) {  // Call from mul.go!
        return a + b
    }
    return 0
}
```

**mul.go**
```go
package math

func validate(a, b int) bool {
    return a > 0 && b > 0
}
```

### Accessing from main

```go
package main

import "myapp/math"

func main() {
    math.Add(5, 3)        // ✅ Exported
    // math.validateThenAdd(5, 3)  // ❌ Unexported
}
```

---

## 🔗 Import Paths

### Relative Paths (Wrong)

```go
// ❌ DON'T DO THIS
import "./helpers"
import "../utils"
```

### Absolute Paths (Right)

```go
// ✅ CORRECT
import "myapp/helpers"
import "myapp/utils"
```

### Module Name Matters

```
go.mod declares: module myapp

So import paths start with: myapp/
```

---

## 📚 Documentation Patterns

### Package Documentation

```go
// Package models contains data structures.
package models

// User represents a user in the system.
type User struct {
    ID   int    // User's unique identifier
    Name string // User's full name
}

// NewUser creates a new user with the given ID and name.
func NewUser(id int, name string) *User {
    return &User{ID: id, Name: name}
}
```

### View Docs

```bash
go doc models
go doc models.User
go doc models.NewUser
```

---

## 🎯 Checklist Before Level 3

- [ ] Understand package structure
- [ ] Can create packages
- [ ] Know exported vs unexported
- [ ] Can organize projects
- [ ] Understand init() functions
- [ ] Can avoid circular imports
- [ ] Can write package documentation
- [ ] Completed 8+ exercises
- [ ] Can build multi-package projects

---

## 💡 Key Takeaways

1. **One directory = One package**
2. **Uppercase = Exported**
3. **Lowercase = Unexported**
4. **init() runs automatically**
5. **Avoid circular imports**
6. **Document your packages**
7. **Organize by responsibility**

---

## 🚀 Next Level

Level 3: Variables & Data Types
- Variable declaration
- Basic types
- Type conversion
- Constants
- Zero values

You're almost there! Keep going! 🎉
