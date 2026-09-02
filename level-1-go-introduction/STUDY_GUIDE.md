# Level 1: Study Guide & Quick Reference

## 📚 Learning Path

### Week 1: Understanding Go
```
Day 1-2:  What is Go & Why Go?
Day 2-3:  Philosophy & Design Principles
Day 3-4:  Installation & Environment Setup
Day 4-5:  Workspace & Module Structure
```

### Week 2: Hands-On Setup
```
Day 1:    Exercises 1-3 (Installation & Module)
Day 2:    Exercises 4-7 (Create & Build)
Day 3:    Exercises 8-11 (Format & Packages)
Day 4:    Exercises 12-15 (Verification & Summary)
```

---

## 🎯 Go Architecture Overview

```
┌─────────────────────────────────────────────────┐
│            YOUR GO PROGRAM                      │
│  ┌─────────────────────────────────────────┐   │
│  │  main.go (package main)                 │   │
│  │  - Entry point                          │   │
│  │  - func main() { ... }                  │   │
│  └─────────────────────────────────────────┘   │
└─────────────────────────────────────────────────┘
         ↓
┌─────────────────────────────────────────────────┐
│        STANDARD LIBRARY                         │
│  fmt, os, time, strings, math, ...              │
└─────────────────────────────────────────────────┘
         ↓
┌─────────────────────────────────────────────────┐
│      GO COMPILER & RUNTIME                      │
│  Compiles → Optimizes → Links → Executable     │
└─────────────────────────────────────────────────┘
         ↓
┌─────────────────────────────────────────────────┐
│      OPERATING SYSTEM                           │
│  macOS / Linux / Windows / ...                  │
└─────────────────────────────────────────────────┘
```

---

## 🔧 Installation Checklist

### Step-by-Step Verification

```
✅ Go is Installed
   └─ Command: go version

✅ GOROOT is Set
   └─ Command: go env GOROOT
   └─ Should show: /usr/local/go or similar

✅ PATH is Updated
   └─ Command: echo $PATH | grep go
   └─ Should include: /usr/local/go/bin

✅ Can Run go help
   └─ Command: go help
   └─ Should show list of commands

✅ Can Create Modules
   └─ Command: go mod init example.com/test
   └─ Should create: go.mod file

✅ Can Build Programs
   └─ Command: go build
   └─ Should create: executable binary

✅ Can Run Programs
   └─ Command: ./program-name
   └─ Should run without errors
```

---

## 📊 Go Commands Reference

### Most Used Commands

```bash
┌─────────────────────────────────────────┐
│  COMMAND        │  PURPOSE              │
├─────────────────────────────────────────┤
│  go run         │  Compile & run        │
│  go build       │  Compile to binary    │
│  go test        │  Run tests            │
│  go fmt         │  Format code          │
│  go mod         │  Manage modules       │
│  go get         │  Install packages     │
│  go doc         │  View documentation  │
│  go clean       │  Remove artifacts     │
└─────────────────────────────────────────┘
```

### Command Syntax

```go
// Compile and run directly
go run main.go

// Compile to executable in current dir
go build

// Compile to specific output
go build -o myapp

// Compile for different platform
GOOS=windows GOARCH=amd64 go build

// Format all files
go fmt ./...

// Initialize module
go mod init example.com/app

// Tidy dependencies
go mod tidy

// View package documentation
go doc fmt
```

---

## 📁 Project Structure

### Single File Project
```
my-app/
├── go.mod
├── main.go
└── README.md
```

### Multi-Package Project
```
my-app/
├── go.mod
├── go.sum
├── main.go              (package main, entry point)
├── config/
│   └── config.go        (package config)
├── utils/
│   └── helpers.go       (package utils)
└── README.md
```

### Key Files

#### go.mod
```go
module example.com/app

go 1.21

require (
    github.com/some/package v1.2.3
)
```

**Purpose:** Declares module name and dependencies

#### go.sum
```
github.com/some/package v1.2.3 h1:abc...
github.com/some/package v1.2.3/go.mod h1:def...
```

**Purpose:** Locks exact versions for reproducibility

#### main.go
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello")
}
```

**Purpose:** Entry point, must have `func main()`

---

## 🔍 Understanding Packages

### Package Declaration

```go
package main    // Executable program
// OR
package utils   // Library/reusable code
```

### Package Names
- `package main` - Only for programs with `func main()`
- `package other` - For libraries and modules
- **Must match directory name** (usually)

### Importing Packages

```go
// Single import
import "fmt"

// Multiple imports
import (
    "fmt"
    "os"
    "time"
)

// Alias import
import f "fmt"     // Use as f.Println()

// Blank import (for side effects)
import _ "database/sql/driver"
```

---

## 💾 go.mod Management

### Creating Module
```bash
go mod init example.com/myapp
```

Creates:
```go
module example.com/myapp

go 1.21
```

### Adding Dependency
```bash
go get github.com/username/package
```

Automatically updates go.mod:
```go
require (
    github.com/username/package v1.2.3
)
```

### Viewing Dependencies
```bash
go mod graph
```

### Cleaning Up Dependencies
```bash
go mod tidy
```

Removes unused dependencies, adds missing ones

### Vendor Dependencies (Optional)
```bash
go mod vendor
```

Copies all dependencies to local `vendor/` directory

---

## 🛠️ Development Workflow

### Day-to-Day Workflow

```
1. WRITE
   └─ Edit main.go in your editor

2. FORMAT
   └─ Run: go fmt ./...

3. TEST (mentally at this stage)
   └─ Trace through your code

4. RUN
   └─ Run: go run main.go

5. BUILD (when ready)
   └─ Run: go build

6. TEST (run executable)
   └─ Run: ./program-name
```

### Alias for Speed

Add to `~/.bashrc` or `~/.zshrc`:

```bash
# Quick commands
alias gr='go run main.go'
alias gb='go build && ./main'
alias gf='go fmt ./...'
alias gt='go test ./...'
alias gd='go doc'
```

Then use:
```bash
gr              # Instead of: go run main.go
gb              # Instead of: go build && ./main
gf              # Instead of: go fmt ./...
```

---

## 📚 Go vs Level 0

### Mapping Concepts

#### Level 0 → Go

| Level 0 | Go | Example |
|---------|-----|---------|
| Variables | `:=` assignment | `name := "Go"` |
| Data Types | Typed variables | `var age int = 25` |
| Sequence | Code order | Statements top-to-bottom |
| If-Else | `if/else` | `if x > 0 { ... }` |
| Loops | `for` loop | `for i := 0; i < 10; i++` |
| Functions | `func` | `func add(a, b int) int` |
| Print | `fmt.Print` | `fmt.Println("hi")` |

### Differences

| Aspect | Level 0 | Go |
|--------|---------|-----|
| Syntax | Pseudocode | Real language |
| Compilation | N/A | Required |
| Types | Basic | Strict (enforced) |
| Execution | Conceptual | Binary runs |

---

## ⚠️ Common Mistakes

### Mistake 1: Wrong Package Name
```go
// ❌ WRONG: Directory name doesn't match package
// File: helpers/helpers.go
package utils    // Mismatch!

// ✅ RIGHT: Directory and package match
// File: utils/utils.go
package utils
```

### Mistake 2: Missing Main Package
```go
// ❌ WRONG: Not an executable
package lib
func main() { }   // Invalid in package lib

// ✅ RIGHT: Only in package main
package main
func main() { }
```

### Mistake 3: Forgetting Imports
```go
// ❌ WRONG: Using fmt without importing
func main() {
    fmt.Println("Hi")  // Error!
}

// ✅ RIGHT: Import first
import "fmt"
func main() {
    fmt.Println("Hi")
}
```

### Mistake 4: Wrong go.mod Name
```bash
# ❌ WRONG: Complex name
go mod init my@project.app

# ✅ RIGHT: Simple, clear
go mod init github.com/username/project
```

### Mistake 5: Not Running from Right Directory
```bash
# ❌ WRONG: main.go is elsewhere
go run main.go    # File not found

# ✅ RIGHT: In same directory
cd ~/my-app
go run main.go
```

---

## 🚀 Quick Start Template

Copy-paste this to start a new project:

```bash
#!/bin/bash
# Setup new Go project

PROJECT_NAME=$1

if [ -z "$PROJECT_NAME" ]; then
    echo "Usage: ./setup-go-project.sh project-name"
    exit 1
fi

# Create directory
mkdir $PROJECT_NAME
cd $PROJECT_NAME

# Initialize module
go mod init github.com/yourname/$PROJECT_NAME

# Create main.go
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("Welcome to Go!")
}
EOF

# Create README
cat > README.md << 'EOF'
# $PROJECT_NAME

## Getting Started

### Run
```bash
go run main.go
```

### Build
```bash
go build
```

EOF

echo "✅ Project '$PROJECT_NAME' created!"
echo "cd $PROJECT_NAME && go run main.go"
```

Save as `setup-go-project.sh` and use:
```bash
chmod +x setup-go-project.sh
./setup-go-project.sh my-new-app
```

---

## 🎯 Checklist: Before Moving to Level 2

- [ ] Go is installed (verified with `go version`)
- [ ] Can create new projects with `go mod init`
- [ ] Can write and run basic Go programs
- [ ] Can build executables with `go build`
- [ ] Understand package structure
- [ ] Can format code with `go fmt`
- [ ] Understand difference between `go run` and `go build`
- [ ] Know key Go commands
- [ ] Editor is configured (optional but helpful)
- [ ] Can read compiler error messages
- [ ] Completed all 15 exercises

---

## 💡 Pro Tips

### Tip 1: Use Playground for Quick Tests
Visit https://play.golang.org to run Go code online without installation.

### Tip 2: Read Error Messages Carefully
Go error messages are usually very helpful:
```
./main.go:5:5: undefined: fmt    // fmt not imported
./main.go:7:19: undefined: name   // variable not declared
```

### Tip 3: Use go doc
```bash
go doc fmt              # Package docs
go doc fmt.Println      # Specific function
go doc -all fmt         # Everything
```

### Tip 4: Leverage Code Completion
With VS Code + Go extension, get:
- Function suggestions
- Parameter hints
- Documentation tooltips
- Error squiggles

### Tip 5: Keep Programs Simple
Start small, add features gradually:
```
✓ HelloWorld
✓ Print multiple lines
✓ Use variables
✓ Read user input
✓ Do calculations
```

---

## 📊 Environment Variables Cheat Sheet

```bash
# View all Go settings
go env

# Specific settings
go env GOROOT       # Where Go is installed
go env GOPATH       # Workspace directory
go env GOOS         # Operating system
go env GOARCH       # CPU architecture
go env GOBIN        # Binary installation directory

# Cross-compile
GOOS=windows GOARCH=amd64 go build    # For Windows
GOOS=darwin GOARCH=amd64 go build     # For Intel Mac
GOOS=darwin GOARCH=arm64 go build     # For Apple Silicon
GOOS=linux GOARCH=amd64 go build      # For Linux
```

---

## 🔗 Quick Links

- **Official Go Site:** https://golang.org
- **Go Packages:** https://pkg.go.dev
- **Go Playground:** https://play.golang.org
- **Effective Go:** https://golang.org/doc/effective_go
- **Go Blog:** https://go.dev/blog

---

## Next Steps

1. Complete all exercises
2. Write your own small program
3. Build and run it
4. Share with someone
5. **Move to Level 2: Go Program Structure**

---

## 🎓 Level 1 Summary

**What You Learned:**
- ✅ What Go is and why it matters
- ✅ How to install and verify Go
- ✅ Go's philosophy and design
- ✅ How to set up projects with modules
- ✅ How to write and run programs
- ✅ Basic Go commands

**What You Can Do:**
- ✅ Create new Go projects
- ✅ Write simple programs
- ✅ Compile and run executables
- ✅ Format code properly
- ✅ Use basic Go commands

**Ready For:**
- ✅ Level 2: Program Structure
  - Multiple files and packages
  - Proper project organization
  - Working with modules

🚀 **You've got this!**
