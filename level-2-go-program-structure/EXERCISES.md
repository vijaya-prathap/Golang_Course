# Level 2: Hands-On Exercises & Practice

## 🎯 Practical Exercises

Complete these exercises to master Go program structure, packages, and file organization.

---

## Exercise 1: Create Your First Package

### Objective
Create a simple package and use it from main.

### Steps

```bash
# 1. Create project directory
mkdir level2-project && cd level2-project

# 2. Initialize module
go mod init myproject.com/app

# 3. Create helpers package directory
mkdir helpers

# 4. Create helpers/greetings.go
cat > helpers/greetings.go << 'EOF'
package helpers

import "fmt"

// SayHello prints a greeting
func SayHello(name string) {
    fmt.Printf("Hello, %s!\n", name)
}

// GetGreeting returns a greeting string
func GetGreeting(name string) string {
    return "Hello, " + name + "!"
}
EOF

# 5. Create main.go
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "myproject.com/app/helpers"
)

func main() {
    helpers.SayHello("Alice")
    
    greeting := helpers.GetGreeting("Bob")
    fmt.Println(greeting)
}
EOF

# 6. Run the project
go run .
```

**Expected Output:**
```
Hello, Alice!
Hello, Bob!
```

### Verify
- ✅ helpers package created
- ✅ Functions are exported (capitalized)
- ✅ main.go imports and uses helpers
- ✅ Program runs without errors

---

## Exercise 2: Multiple Files in One Package

### Objective
Create a package with multiple Go files.

### Steps

```bash
# Create math package directory
mkdir math

# Create math/add.go
cat > math/add.go << 'EOF'
package math

func Add(a, b int) int {
    return a + b
}
EOF

# Create math/subtract.go
cat > math/subtract.go << 'EOF'
package math

func Subtract(a, b int) int {
    return a - b
}
EOF

# Create math/multiply.go
cat > math/multiply.go << 'EOF'
package math

func Multiply(a, b int) int {
    return a * b
}
EOF

# Update main.go
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "myproject.com/app/math"
)

func main() {
    fmt.Println("5 + 3 =", math.Add(5, 3))
    fmt.Println("5 - 3 =", math.Subtract(5, 3))
    fmt.Println("5 * 3 =", math.Multiply(5, 3))
}
EOF

# Run
go run .
```

**Expected Output:**
```
5 + 3 = 8
5 - 3 = 2
5 * 3 = 15
```

### File Structure
```
myproject/
├── go.mod
├── main.go
├── helpers/
│   └── greetings.go
└── math/
    ├── add.go
    ├── subtract.go
    └── multiply.go
```

---

## Exercise 3: Exported vs Unexported

### Objective
Learn the difference between exported and unexported identifiers.

### Steps

```bash
# Create utils package
mkdir utils

# Create utils/strings.go
cat > utils/strings.go << 'EOF'
package utils

import "strings"

// ReverseString is exported - can be used by other packages
func ReverseString(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

// CountVowels is exported
func CountVowels(s string) int {
    count := 0
    vowels := "aeiouAEIOU"
    for _, char := range s {
        if strings.Contains(vowels, string(char)) {
            count++
        }
    }
    return count
}

// capitalize is unexported - only used within this package
func capitalize(word string) string {
    return strings.ToUpper(word)
}

// isVowel is unexported - helper function
func isVowel(ch rune) bool {
    return strings.ContainsRune("aeiouAEIOU", ch)
}
EOF

# Update main.go
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "myproject.com/app/utils"
)

func main() {
    text := "hello golang"
    
    // ✅ These work - functions are exported
    reversed := utils.ReverseString(text)
    fmt.Println("Reversed:", reversed)
    
    vowelCount := utils.CountVowels(text)
    fmt.Println("Vowels:", vowelCount)
    
    // ❌ These would NOT work - functions are unexported
    // utils.capitalize("test")    // Compiler error!
    // utils.isVowel('a')          // Compiler error!
}
EOF

go run .
```

**Expected Output:**
```
Reversed: gnalop olleh
Vowels: 3
```

### Key Learning
- ✅ Exported functions: `ReverseString`, `CountVowels`
- ❌ Unexported functions: `capitalize`, `isVowel`
- Can't access unexported identifiers from other packages

---

## Exercise 4: Package-Level Variables

### Objective
Work with variables at package level.

### Steps

```bash
# Create database package
mkdir database

# Create database/db.go
cat > database/db.go << 'EOF'
package database

var connectionCount = 0
var isConnected = false

// GetConnectionCount returns how many connections exist
func GetConnectionCount() int {
    return connectionCount
}

// Connect simulates database connection
func Connect() bool {
    connectionCount++
    isConnected = true
    return isConnected
}

// IsConnected returns connection status
func IsConnected() bool {
    return isConnected
}
EOF

# Update main.go
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "myproject.com/app/database"
)

func main() {
    fmt.Println("Connected:", database.IsConnected())
    fmt.Println("Connections:", database.GetConnectionCount())
    
    database.Connect()
    fmt.Println("Connected:", database.IsConnected())
    fmt.Println("Connections:", database.GetConnectionCount())
}
EOF

go run .
```

**Expected Output:**
```
Connected: false
Connections: 0
Connected: true
Connections: 1
```

### Key Learning
- Package variables accessible only within their package
- Provide functions to access/modify them
- Use getters (GetXxx) for exported access

---

## Exercise 5: Using init() Functions

### Objective
Understand and use init() functions.

### Steps

```bash
# Create config package
mkdir config

# Create config/init.go
cat > config/init.go << 'EOF'
package config

import (
    "fmt"
    "log"
)

var AppName string
var Version string

func init() {
    fmt.Println("[init] Setting up configuration...")
    AppName = "MyApplication"
    Version = "1.0.0"
    fmt.Println("[init] Configuration ready!")
}

func GetAppInfo() string {
    return AppName + " v" + Version
}
EOF

# Create config/logger.go
cat > config/logger.go << 'EOF'
package config

import (
    "fmt"
)

var logCount = 0

func init() {
    fmt.Println("[init] Initializing logger...")
    logCount = 0
}

func Log(message string) {
    logCount++
    fmt.Printf("[Log %d] %s\n", logCount, message)
}
EOF

# Update main.go
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "myproject.com/app/config"
)

func init() {
    fmt.Println("[main init] Running main init")
}

func main() {
    fmt.Println("[main] Starting application...")
    fmt.Println(config.GetAppInfo())
    config.Log("Application started")
}
EOF

go run .
```

**Expected Output:**
```
[init] Setting up configuration...
[init] Configuration ready!
[init] Initializing logger...
[main init] Running main init
[main] Starting application...
MyApplication v1.0.0
[Log 1] Application started
```

### Key Learning
- ✅ init() runs automatically
- ✅ Multiple init() functions supported
- ✅ Runs before main()
- ✅ Good for setup and initialization

---

## Exercise 6: Import Aliases

### Objective
Use package aliases for conflict resolution and clarity.

### Steps

```bash
# Create statistics package
mkdir statistics

# Create statistics/calc.go
cat > statistics/calc.go << 'EOF'
package statistics

func Average(numbers []int) float64 {
    if len(numbers) == 0 {
        return 0
    }
    sum := 0
    for _, n := range numbers {
        sum += n
    }
    return float64(sum) / float64(len(numbers))
}
EOF

# Update main.go
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "math"                         // Standard library
    stats "myproject.com/app/statistics"  // Alias
)

func main() {
    // Using standard library math
    result1 := math.Sqrt(16)
    fmt.Println("Standard sqrt:", result1)
    
    // Using our statistics package
    numbers := []int{1, 2, 3, 4, 5}
    avg := stats.Average(numbers)
    fmt.Println("Average:", avg)
}
EOF

go run .
```

**Expected Output:**
```
Standard sqrt: 4
Average: 3
```

### Key Learning
- Use aliases for clarity
- Avoid naming conflicts
- Syntax: `alias "import/path"`

---

## Exercise 7: Organize a Real Project

### Objective
Create a well-organized project with multiple packages.

### Steps

```bash
# Create project structure
mkdir -p models
mkdir -p services
mkdir -p handlers

# Create models/user.go
cat > models/user.go << 'EOF'
package models

type User struct {
    ID   int
    Name string
    Email string
}

func NewUser(id int, name, email string) *User {
    return &User{
        ID:    id,
        Name:  name,
        Email: email,
    }
}
EOF

# Create services/user_service.go
cat > services/user_service.go << 'EOF'
package services

import (
    "myproject.com/app/models"
)

var users []*models.User

func CreateUser(id int, name, email string) *models.User {
    user := models.NewUser(id, name, email)
    users = append(users, user)
    return user
}

func GetUser(id int) *models.User {
    for _, user := range users {
        if user.ID == id {
            return user
        }
    }
    return nil
}

func ListUsers() []*models.User {
    return users
}
EOF

# Create handlers/user_handler.go
cat > handlers/user_handler.go << 'EOF'
package handlers

import (
    "fmt"
    "myproject.com/app/models"
    "myproject.com/app/services"
)

func PrintUser(user *models.User) {
    if user == nil {
        fmt.Println("User not found")
        return
    }
    fmt.Printf("User: %s (Email: %s)\n", user.Name, user.Email)
}

func PrintAllUsers() {
    users := services.ListUsers()
    fmt.Printf("Total users: %d\n", len(users))
    for _, user := range users {
        PrintUser(user)
    }
}
EOF

# Update main.go
cat > main.go << 'EOF'
package main

import (
    "myproject.com/app/handlers"
    "myproject.com/app/services"
)

func main() {
    // Create users
    services.CreateUser(1, "Alice", "alice@example.com")
    services.CreateUser(2, "Bob", "bob@example.com")
    services.CreateUser(3, "Charlie", "charlie@example.com")
    
    // Print all users
    handlers.PrintAllUsers()
    
    // Get specific user
    user := services.GetUser(2)
    handlers.PrintUser(user)
}
EOF

go run .
```

**Expected Output:**
```
Total users: 3
User: Alice (Email: alice@example.com)
User: Bob (Email: bob@example.com)
User: Charlie (Email: charlie@example.com)
User: Bob (Email: bob@example.com)
```

### Project Structure
```
myproject/
├── go.mod
├── main.go
├── models/
│   └── user.go
├── services/
│   └── user_service.go
└── handlers/
    └── user_handler.go
```

---

## Exercise 8: Avoiding Circular Imports

### Objective
Learn how to avoid circular import errors.

### Steps

```bash
# Create interfaces package (breaks circular dependency)
mkdir interfaces

# Create interfaces/repository.go
cat > interfaces/repository.go << 'EOF'
package interfaces

type Repository interface {
    Save(id int, name string) error
    Get(id int) (string, error)
}
EOF

# Create models/user.go (updated)
cat > models/user.go << 'EOF'
package models

type User struct {
    ID   int
    Name string
}

func NewUser(id int, name string) *User {
    return &User{ID: id, Name: name}
}
EOF

# Create storage/user_storage.go
cat > storage/user_storage.go << 'EOF'
package storage

import (
    "myproject.com/app/models"
)

var storage = make(map[int]*models.User)

func SaveUser(id int, name string) error {
    storage[id] = models.NewUser(id, name)
    return nil
}

func GetUser(id int) *models.User {
    return storage[id]
}
EOF

# Update main.go
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "myproject.com/app/storage"
)

func main() {
    storage.SaveUser(1, "Alice")
    user := storage.GetUser(1)
    fmt.Printf("User: %s\n", user.Name)
}
EOF

go run .
```

**Expected Output:**
```
User: Alice
```

### Key Learning
- ✅ Create interfaces package to break cycles
- ✅ Use dependency injection
- ✅ Extract common types

---

## Exercise 9: Build and Run

### Objective
Learn to build and run multi-package projects.

### Steps

```bash
# Format all files
go fmt ./...

# Build the project
go build

# Run the executable
./myproject

# Build with custom name
go build -o myapp

# Run custom executable
./myapp

# Clean up
go clean
```

### Verify
- ✅ `go fmt` formats all files
- ✅ `go build` creates executable
- ✅ Executable runs successfully
- ✅ `go clean` removes artifacts

---

## Exercise 10: Package Documentation

### Objective
Document your packages properly.

### Steps

```bash
# Create documented packages

# Update database/db.go with comments
cat > database/db.go << 'EOF'
// Package database provides database connection management.
package database

import "fmt"

var connectionCount = 0

// Connect establishes a database connection.
// It increments the connection counter and returns true on success.
func Connect() bool {
    connectionCount++
    fmt.Println("Connected to database")
    return true
}

// GetConnectionCount returns the current number of connections.
func GetConnectionCount() int {
    return connectionCount
}
EOF

# View documentation
go doc database
go doc database.Connect

# Serve documentation
godoc -http=:6060
# Visit: http://localhost:6060/pkg/myproject.com/app/
```

### Key Learning
- ✅ Package comments: `// Package name ...`
- ✅ Function comments: `// FunctionName ...`
- ✅ Use `go doc` to view
- ✅ Use `godoc` to serve HTML docs

---

## ✅ Exercise Completion Checklist

- [ ] Exercise 1: Create first package
- [ ] Exercise 2: Multiple files in package
- [ ] Exercise 3: Exported vs unexported
- [ ] Exercise 4: Package-level variables
- [ ] Exercise 5: Using init() functions
- [ ] Exercise 6: Import aliases
- [ ] Exercise 7: Organize real project
- [ ] Exercise 8: Avoid circular imports
- [ ] Exercise 9: Build and run
- [ ] Exercise 10: Package documentation

---

## 🎯 Bonus Challenges

### Challenge 1: Book Management System
Create packages for:
- `models` - Book struct
- `library` - Add/search books
- `handlers` - Display functions

### Challenge 2: Configuration System
Create a config package with:
- Multiple init() functions
- Package-level variables
- Getter functions

### Challenge 3: Logger Package
Create a logging package:
- Multiple log levels (Info, Warning, Error)
- Package-level configuration
- Exported functions

### Challenge 4: Multiple Packages
Reorganize your complete project into:
- `models` - Data types
- `services` - Business logic
- `handlers` - Output formatting
- `utils` - Helpers

---

## 🚀 Ready for Level 3?

Once you've completed these exercises:
- ✅ Can create and use packages
- ✅ Understand exported/unexported
- ✅ Can organize projects properly
- ✅ Know init() functions
- ✅ Can avoid circular imports

**Next:** Variables & Data Types

---

## 📞 Common Issues

**Issue:** Package not found
```
Solution: Check go.mod module name and import path
```

**Issue:** Exported function but still can't access
```
Solution: Ensure it starts with CAPITAL letter
```

**Issue:** Circular import error
```
Solution: Restructure packages or use interfaces
```

**Issue:** init() not running
```
Solution: Make sure package is imported
```

Great job! Keep building! 🎉
