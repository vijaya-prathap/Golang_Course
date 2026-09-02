# Level 2: Go Program Structure

## Overview

Now that you understand Go basics and have it installed, it's time to learn how to structure real-world Go programs. In this level, you'll work with multiple files, create and use packages, manage imports, and understand Go's visibility rules.

---

## 📚 Table of Contents

1. [From Hello World to Real Programs](#from-hello-world-to-real-programs)
2. [Organizing Code with Packages](#organizing-code-with-packages)
3. [Creating Packages](#creating-packages)
4. [Import Statements](#import-statements)
5. [Exported vs Unexported (Public vs Private)](#exported-vs-unexported-public-vs-private)
6. [Package-Level Variables](#package-level-variables)
7. [Init Functions](#init-functions)
8. [The main Package Revisited](#the-main-package-revisited)
9. [Multi-File Projects](#multi-file-projects)
10. [Package Naming Conventions](#package-naming-conventions)
11. [Circular Imports (and Avoiding Them)](#circular-imports-and-avoiding-them)
12. [Working with go.mod](#working-with-gomod)

---

## From Hello World to Real Programs

### The Problem with Single Files

Your first Go program was simple:
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

This works, but real programs are much larger. Writing 10,000 lines in one file causes problems:

**Single File Issues:**
- Hard to find code
- Difficult to organize related functions
- Code duplication
- Hard to test
- Impossible to reuse code in other projects
- Maintenance nightmare

### The Solution: Packages

Go uses **packages** to organize code into logical units:

```
myapp/
├── main.go              (package main)
├── go.mod
├── math/
│   └── math.go          (package math)
├── utils/
│   └── helpers.go       (package utils)
└── database/
    └── db.go            (package database)
```

**Benefits:**
- Organized, easy to navigate
- Reusable code
- Better testing
- Clear dependencies
- Professional structure

---

## Organizing Code with Packages

### What is a Package?

A **package** is a directory of Go files that work together.

**Key Points:**
- One package per directory
- All Go files in a directory must declare the same package
- Package name doesn't need to match directory name (but it usually does)
- Files in same package can access each other's functions

### Package Types

#### 1. Main Package (Executable)
```go
package main

func main() {
    // Program entry point
}
```

**Characteristics:**
- Special package name: `main`
- Must have `func main()`
- Creates an executable
- Only one per project

#### 2. Regular Packages (Libraries)
```go
package helpers      // Directory: helpers/

func DoSomething() {
    // Function available to other packages
}
```

**Characteristics:**
- Any name except `main`
- No `func main()` needed
- Creates a library
- Multiple packages allowed

### File Organization

**Typical Project Structure:**
```
myapp/
├── go.mod                  # Module definition
├── go.sum                  # Dependency lock
├── main.go                 # Entry point (package main)
│
├── models/                 # Package for data structures
│   ├── user.go
│   ├── product.go
│   └── order.go
│
├── utils/                  # Package for utilities
│   ├── helpers.go
│   └── validators.go
│
├── database/               # Package for DB operations
│   ├── connect.go
│   └── queries.go
│
└── api/                    # Package for API handlers
    ├── handlers.go
    └── middleware.go
```

### Why This Structure?

- **Separation of Concerns:** Each package has one responsibility
- **Easy Navigation:** Find what you need quickly
- **Reusability:** Other projects can use your packages
- **Testing:** Test packages independently
- **Maintenance:** Changes are isolated

---

## Creating Packages

### Step 1: Create Directory

```bash
# Create project structure
mkdir -p myapp/helpers
mkdir -p myapp/math
cd myapp
```

### Step 2: Create Package Files

**Directory: helpers/**
**File: helpers/strings.go**
```go
package helpers        // Declare package name

import "strings"

// Helper function for string operations
func CapitalizeWord(word string) string {
    return strings.Title(word)
}

func ReverseString(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}
```

**Directory: math/**
**File: math/calc.go**
```go
package math

// Add two numbers
func Add(a, b int) int {
    return a + b
}

// Subtract two numbers
func Subtract(a, b int) int {
    return a - b
}

// Multiply two numbers
func Multiply(a, b int) int {
    return a * b
}

// Divide two numbers
func Divide(a, b int) int {
    if b == 0 {
        return 0  // Handle division by zero
    }
    return a / b
}
```

### Step 3: Create main.go

```go
package main

import (
    "fmt"
    "myapp/helpers"    // Import helpers package
    "myapp/math"       // Import math package
)

func main() {
    // Use helpers package
    result1 := helpers.CapitalizeWord("hello")
    fmt.Println(result1)  // Output: Hello

    reversed := helpers.ReverseString("golang")
    fmt.Println(reversed)  // Output: gnalог

    // Use math package
    sum := math.Add(5, 3)
    fmt.Println(sum)  // Output: 8

    product := math.Multiply(4, 7)
    fmt.Println(product)  // Output: 28
}
```

### Step 4: Run

```bash
cd myapp
go mod init myapp           # Initialize module
go run .                    # Run all files in directory
# or
go run main.go              # Just main.go (if imports work)
```

### Project Structure Now

```
myapp/
├── go.mod
├── main.go              (package main)
├── helpers/
│   └── strings.go       (package helpers)
└── math/
    └── calc.go          (package math)
```

---

## Import Statements

### Basic Imports

#### Single Import
```go
import "fmt"
```

#### Multiple Imports
```go
import (
    "fmt"
    "strings"
    "os"
)
```

### Import Aliases

#### Renaming Packages
```go
import (
    "fmt"
    math "myapp/math"      // Alias
)

// Use as:
result := math.Add(5, 3)
```

#### Conflicting Names
```go
// If two packages have the same name, use aliases
import (
    stdmath "math"         // Standard library math
    mypkg "myapp/math"     // Your custom math
)

// Use:
stdmath.Sqrt(16)          // Standard library
mypkg.Add(5, 3)           // Your package
```

### Blank Import

```go
import (
    "database/sql"
    _ "github.com/lib/pq"  // Blank import for side effects
)

// The underscore tells Go: "I import this but don't use it directly"
// The package init() function still runs (registering drivers, etc.)
```

### Dot Imports (Use Carefully!)

```go
import (
    . "fmt"    // Import all identifiers from fmt
)

// Now you can use:
Println("Hello")          // Instead of fmt.Println()
```

**Note:** Avoid dot imports! They make code unclear.

### Import Paths

#### Standard Library
```go
import "fmt"
import "os"
import "time"
```

#### Local Packages
```go
import "myapp/helpers"
import "myapp/utils"
```

#### External Packages (GitHub)
```go
import "github.com/gorilla/mux"
import "github.com/go-sql-driver/mysql"
```

---

## Exported vs Unexported (Public vs Private)

### Go's Visibility Rules

Go doesn't have `public`, `private`, or `protected` keywords. Instead, it uses **capitalization**:

- **Uppercase first letter = Exported (Public)**
- **Lowercase first letter = Unexported (Private)**

### Examples

#### In Package: `math/calc.go`

```go
package math

// ✅ EXPORTED - Can be used by other packages
func Add(a, b int) int {
    return a + b
}

// ✅ EXPORTED - Can be used by other packages
const Pi = 3.14159

// ✅ EXPORTED - Can be used by other packages
type Calculator struct {
    LastResult int
}

// ❌ UNEXPORTED - Only used within math package
func validate(a, b int) bool {
    return a > 0 && b > 0
}

// ❌ UNEXPORTED - Only used within math package
var internalCounter = 0

// ❌ UNEXPORTED - Only used within math package
type config struct {
    precision int
}
```

#### In Package: `main.go`

```go
package main

import (
    "fmt"
    "myapp/math"
)

func main() {
    // ✅ Works - Add is exported
    result := math.Add(5, 3)
    fmt.Println(result)

    // ✅ Works - Pi is exported
    fmt.Println(math.Pi)

    // ✅ Works - Calculator is exported
    calc := math.Calculator{}

    // ❌ ERROR - validate is unexported
    // math.validate(5, 3)    // Compiler error!

    // ❌ ERROR - internalCounter is unexported
    // math.internalCounter    // Compiler error!
}
```

### Exported Struct Fields

```go
package models

// ✅ Exported struct - can use in other packages
type User struct {
    // ✅ Exported field - can access: user.Name
    Name string
    
    // ✅ Exported field - can access: user.Email
    Email string
    
    // ❌ Unexported field - cannot access: user.password
    password string
}

// ✅ Exported method - can call: user.GetPassword()
func (u User) GetPassword() string {
    return u.password
}

// ✅ Exported method - can call: user.SetPassword()
func (u *User) SetPassword(p string) {
    u.password = p
}
```

### Exported Functions

```go
package helpers

// ✅ EXPORTED - Starts with capital letter
func ProcessData(input string) string {
    return helpers.validate(input)  // Can call unexported function
}

// ❌ UNEXPORTED - Starts with lowercase
func validate(input string) string {
    // Processing logic
    return input
}
```

### Best Practices

```
✓ Use exported names for public API
✓ Keep internal helpers unexported
✓ Document exported functions with comments
✓ Use unexported for implementation details
✓ Be intentional about your package's interface
```

---

## Package-Level Variables

### What Are Package-Level Variables?

Variables declared outside functions, at package scope.

```go
package models

import "time"

// ✅ Exported package variable
var UserCount = 0

// ❌ Unexported package variable
var lastUpdated time.Time

// ✅ Exported constant
const MaxUsers = 1000

// ❌ Unexported constant
const internalVersion = "1.0"
```

### How They Work

All files in the same package can access package variables:

**File: models/user.go**
```go
package models

var UserCount = 0  // Package-level variable

func CreateUser() {
    UserCount++    // Increment counter
}
```

**File: models/admin.go**
```go
package models

func GetUserCount() int {
    return UserCount   // Access from other file
}
```

**File: main.go**
```go
package main

import "myapp/models"

func main() {
    models.CreateUser()
    count := models.GetUserCount()
    println(count)  // Output: 1
}
```

### Initialization Order

Package-level variables are initialized in:
1. **Const** - Compile-time constants
2. **Var** - Package variables (in order of declaration)
3. **init()** - Initialization functions

---

## Init Functions

### What is init()?

A special function that runs automatically when package is imported.

```go
package config

import "fmt"

// This function runs automatically
func init() {
    fmt.Println("Config package initialized!")
}
```

### Multiple init() Functions

You can have multiple `init()` functions in a package:

```go
package main

import "fmt"

func init() {
    fmt.Println("First init")
}

func init() {
    fmt.Println("Second init")
}

func init() {
    fmt.Println("Third init")
}

func main() {
    fmt.Println("Main function")
}

// Output:
// First init
// Second init
// Third init
// Main function
```

### Execution Order

```
1. Global variables initialized
2. All init() functions (in order)
3. main() function
```

### Real-World Example

**File: database/db.go**
```go
package database

import (
    "database/sql"
    "log"
)

var db *sql.DB

func init() {
    // Initialize database connection
    var err error
    db, err = sql.Open("mysql", "user:pass@/dbname")
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Database initialized")
}

func GetDB() *sql.DB {
    return db
}
```

**File: main.go**
```go
package main

import (
    "fmt"
    "myapp/database"
)

func main() {
    // init() already ran, database is ready
    db := database.GetDB()
    fmt.Println("Using database...")
}
```

---

## The main Package Revisited

### Special Properties

The `main` package has unique properties:

```go
package main

import "fmt"

// ✅ REQUIRED - Every main package needs this
func main() {
    fmt.Println("Program starts here")
}

// ❌ WRONG - main package is for executables only
// func init() { ... }  // Don't use init in main

// ❌ ERROR - main can't be imported
// import "myapp"  // main.go can't be used as a package
```

### main vs other packages

| Aspect | main | Other |
|--------|------|-------|
| Can have `func main()` | ✅ Required | ❌ Not allowed |
| Can be imported | ❌ No | ✅ Yes |
| Creates executable | ✅ Yes | ❌ No |
| Can use `init()` | ⚠️ Rarely | ✅ Yes |

### Typical main Structure

```go
package main

import (
    "fmt"
    "myapp/models"
    "myapp/database"
)

func init() {
    // Optional: Initialize settings
    fmt.Println("Initializing app...")
}

func main() {
    // Your program logic
    fmt.Println("App started!")
    
    user := models.CreateUser("John")
    fmt.Println(user)
}
```

---

## Multi-File Projects

### Organizing Files in a Package

One package can have multiple files:

**Package structure:**
```
math/
├── add.go       (package math)
├── subtract.go  (package math)
├── multiply.go  (package math)
└── divide.go    (package math)
```

### Multi-File Example

**File: math/add.go**
```go
package math

func Add(a, b int) int {
    return a + b
}
```

**File: math/subtract.go**
```go
package math

func Subtract(a, b int) int {
    return a - b
}
```

**File: math/multiply.go**
```go
package math

func Multiply(a, b int) int {
    return a * b
}
```

**File: math/divide.go**
```go
package math

func Divide(a, b int) int {
    if b == 0 {
        return 0
    }
    return a / b
}
```

**File: main.go**
```go
package main

import (
    "fmt"
    "myapp/math"
)

func main() {
    fmt.Println(math.Add(10, 5))      // Output: 15
    fmt.Println(math.Subtract(10, 5)) // Output: 5
    fmt.Println(math.Multiply(10, 5)) // Output: 50
    fmt.Println(math.Divide(10, 5))   // Output: 2
}
```

### Running Multi-File Projects

```bash
# Run all Go files in directory
go run .

# Run specific files
go run main.go

# Build all
go build

# Run executable
./myapp
```

### When to Split Files

**Good reasons:**
- Different concerns (users.go, products.go)
- Large files (> 500 lines)
- Separating interfaces from implementation
- Grouping related functions

**Bad reasons:**
- Just because
- One function per file
- Arbitrary splitting

---

## Package Naming Conventions

### Go Package Naming Rules

Go has conventions for package names:

### 1. Short Names
```go
package net          // Good - short, clear
package http         // Good
package utils        // Good
package myapplication // Bad - too long
```

### 2. Lowercase Only
```go
package math      // ✅ Good
package Math      // ❌ Bad - starts with capital
package myApp     // ❌ Bad - mixed case
```

### 3. Single Word (Usually)
```go
package models        // ✅ Good
package database      // ✅ Good
package user_handlers // ⚠️ Avoid underscores
```

### 4. Avoid These Names
```
// Standard library names (confusing)
package math          // Conflicts with "math"
package http          // Conflicts with "http"
package os            // Conflicts with "os"

// Use your own names:
package mymath        // ✅ Clear
package api           // ✅ Clear
package handlers      // ✅ Clear
```

### Common Naming Patterns

```go
package models      // Data structures (User, Product)
package handlers    // HTTP request handlers
package middleware  // HTTP middleware
package repository  // Database operations
package service     // Business logic
package utils       // Utility functions
package config      // Configuration
package errors      // Error definitions
```

---

## Circular Imports (and Avoiding Them)

### What is a Circular Import?

When packages import each other, creating a cycle:

```
models.go imports database.go
database.go imports models.go
ERROR: Circular import!
```

### Example (DON'T DO THIS!)

**File: models/user.go**
```go
package models

import "myapp/database"  // Imports database

type User struct {
    ID   int
    Name string
}

func GetUser(id int) *User {
    return database.QueryUser(id)  // Uses database
}
```

**File: database/db.go**
```go
package database

import "myapp/models"  // Imports models - CIRCULAR!

func QueryUser(id int) *models.User {
    return &models.User{ID: id}
}
```

**Result:**
```
ERROR: circular import
```

### How to Fix It

#### Solution 1: Dependency Injection

```go
// models/user.go
package models

type User struct {
    ID   int
    Name string
}

// GetUser doesn't call database, caller does
func NewUser(id int, name string) *User {
    return &User{ID: id, Name: name}
}
```

```go
// database/db.go
package database

import "myapp/models"

func QueryUser(id int) *models.User {
    // Database query logic
    return models.NewUser(id, "John")
}
```

```go
// main.go
package main

import (
    "myapp/database"
)

func main() {
    user := database.QueryUser(1)
    println(user.Name)
}
```

#### Solution 2: Extract Common Types

```
common/
├── types.go    (User type - no dependencies)
models/
├── user.go     (imports common)
database/
├── db.go       (imports common)
```

**File: common/types.go**
```go
package common

type User struct {
    ID   int
    Name string
}
```

**File: models/user.go**
```go
package models

import "myapp/common"

func GetUser(user common.User) string {
    return user.Name
}
```

**File: database/db.go**
```go
package database

import "myapp/common"

func SaveUser(user common.User) {
    // Save to database
}
```

#### Solution 3: Use Interfaces

```go
// repository/interface.go
package repository

type UserRepository interface {
    SaveUser(id int, name string) error
    GetUser(id int) (string, error)
}
```

```go
// models/user.go
package models

import "myapp/repository"

type UserService struct {
    repo repository.UserRepository
}

func (s *UserService) CreateUser(name string) error {
    return s.repo.SaveUser(1, name)
}
```

---

## Working with go.mod

### What is go.mod?

A file that declares your module and its dependencies.

### Basic Structure

```go
module myapp

go 1.21

require (
    github.com/gorilla/mux v1.8.0
    github.com/lib/pq v1.10.0
)
```

### Local Imports

When you have local packages:

```bash
$ tree myapp
myapp/
├── go.mod
├── main.go
└── helpers/
    └── help.go
```

**In main.go:**
```go
package main

import "myapp/helpers"  // Local import
```

**go.mod:**
```
module myapp

go 1.21
```

### Checking for Unused Dependencies

```bash
go mod tidy
```

Removes unused imports from go.mod.

### Viewing All Dependencies

```bash
go mod graph
```

Shows dependency tree.

---

## ✅ Level 2 Checklist

By the end of this level, you should:

- [ ] Understand why packages are important
- [ ] Create your own packages
- [ ] Organize projects with multiple files
- [ ] Use import statements correctly
- [ ] Understand exported vs unexported
- [ ] Know Go's visibility rules
- [ ] Use package-level variables
- [ ] Understand init() functions
- [ ] Structure main packages properly
- [ ] Avoid circular imports
- [ ] Know package naming conventions
- [ ] Work with multi-file projects
- [ ] Configure go.mod for local packages

---

## 🎯 Key Concepts Review

### Package Structure
```
One directory = One package
One package = Can have multiple files
```

### Naming Rules
```
Uppercase = Exported (public)
Lowercase = Unexported (private)
```

### Initialization Order
```
1. Global variables
2. init() functions
3. main() function
```

### Circular Imports
```
Avoid by:
- Dependency injection
- Extracting common types
- Using interfaces
```

---

## 🚀 Ready for Level 3?

Once you've completed Level 2:
- ✅ Understand project organization
- ✅ Can create multiple packages
- ✅ Know export/import mechanics
- ✅ Can structure medium-sized projects

**Next Level:** Level 3 - Variables & Data Types
- Variable declaration styles
- Basic types (int, float, string, bool)
- Type conversion
- Constants
- Zero values

---

## 💡 Key Takeaways

1. **Packages organize code** - Use them to structure projects
2. **Capitalization matters** - First letter determines visibility
3. **One directory, one package** - Keep it simple
4. **No circular imports** - Design to avoid them
5. **init() runs automatically** - Use for setup
6. **go.mod manages dependencies** - Keep it clean

---

## 📚 Next Topics

After mastering program structure, you'll learn:
- Variable declaration and types
- Basic operations
- Control flow (if, for, switch)
- Functions and methods
- Structs and interfaces

The foundation is set! Keep building! 🚀
