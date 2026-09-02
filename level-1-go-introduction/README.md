# Level 1: Go Introduction & Environment

## Overview

Welcome to Go! In this level, you'll set up your development environment and understand why Go is special. We'll explore Go's philosophy, install the necessary tools, and prepare for writing your first Go program.

---

## 📚 Table of Contents

1. [What is Go?](#what-is-go)
2. [Why Choose Go?](#why-choose-go)
3. [Go Philosophy & Design](#go-philosophy--design)
4. [Setting Up Your Environment](#setting-up-your-environment)
5. [Understanding the Go Workspace](#understanding-the-go-workspace)
6. [Your First Program](#your-first-program)
7. [Running Go Programs](#running-go-programs)
8. [Go Tools & Commands](#go-tools--commands)
9. [Troubleshooting](#troubleshooting)
10. [Key Concepts from Level 0 → Go](#key-concepts-from-level-0--go)

---

## What is Go?

### Go Definition
**Go** (also called **Golang**) is a compiled, statically-typed programming language created by **Google** in 2007 by Rob Pike, Ken Thompson, and Robert Griesemer.

### Quick Facts
- **Created:** 2007 at Google
- **Open Source:** Yes (since 2009)
- **Type System:** Statically typed (but feels dynamic)
- **Compilation:** Compiled to machine code
- **Performance:** Fast, similar to C/C++
- **Syntax:** Simple and clean
- **Current Version:** 1.21+ (as of 2024)

### Go's Logo: The Gopher
Go uses a cute gopher mascot as its logo. The gopher has become iconic in the Go community!

---

## Why Choose Go?

### Go Strengths

#### 1. **Simplicity**
```
✓ Clean syntax
✓ Minimal keywords (25 keywords vs 50+ in other languages)
✓ Fast learning curve
✓ Easy to read and understand
```

#### 2. **Performance**
```
✓ Compiled to native machine code
✓ Fast execution (comparable to C/C++)
✓ Low memory footprint
✓ Efficient concurrency
```

#### 3. **Concurrency**
```
✓ Goroutines (lightweight threads)
✓ Channels for communication
✓ Built-in concurrency primitives
✓ Handle millions of goroutines
```

#### 4. **Production-Ready**
```
✓ Used by Google, Docker, Kubernetes, etc.
✓ Excellent standard library
✓ Cross-platform compilation
✓ Single binary deployment
```

#### 5. **Developer Experience**
```
✓ Fast compilation
✓ Built-in formatting (gofmt)
✓ Excellent error messages
✓ Active community
```

### Use Cases for Go

| Use Case | Why Go | Examples |
|----------|--------|----------|
| **Backend APIs** | Fast, concurrent, simple | REST APIs, microservices |
| **DevOps Tools** | Static binary, fast | Docker, Kubernetes, Terraform |
| **Cloud Services** | Efficient, scalable | Cloud platforms, services |
| **CLI Tools** | Single binary, fast | Command-line utilities |
| **Systems Software** | Performance, control | System tools, networking |
| **Databases** | Concurrent, efficient | InfluxDB, CockroachDB |

---

## Go Philosophy & Design

### Go's Core Principles

#### 1. **Simplicity Over Features**
Go intentionally omits features found in other languages:
- ❌ No inheritance (use composition instead)
- ❌ No generics (until Go 1.18)
- ❌ No exceptions (use error handling)
- ❌ No complex syntax

This makes Go easier to learn and read.

#### 2. **Explicit Over Implicit**
```go
// Go prefers explicit code
if err != nil {
    return err  // Handle errors explicitly
}

// Over implicit features like exceptions:
try {
    // Python would do this
} catch {
    
}
```

#### 3. **Readability is Important**
Go's motto: "Code is read much more often than it is written"

- Single way to do things
- Enforced code formatting (gofmt)
- Clear error handling
- Minimal syntax surprises

#### 4. **Concurrency as a First-Class Citizen**
Goroutines and channels make concurrent programming easy.

#### 5. **Composition Over Inheritance**
```go
// Go uses composition
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Instead of inheritance hierarchies
```

---

## Setting Up Your Environment

### System Requirements

#### Minimum Requirements
- **OS:** macOS, Linux, Windows, FreeBSD
- **Disk Space:** ~500 MB for Go installation
- **RAM:** 1 GB minimum (2+ GB recommended)
- **Processor:** Any modern processor

#### Recommended Setup
- macOS 10.12+ or Linux (Ubuntu 18.04+) or Windows 10+
- 2+ GB RAM
- Text editor or IDE
- Terminal/Command line access

### Installation Steps

#### On macOS

**Using Homebrew (Recommended):**
```bash
# Install Go using Homebrew
brew install go

# Verify installation
go version
```

**Using Installer:**
1. Download from https://golang.org/dl/
2. Download macOS installer (pkg)
3. Run the installer
4. Follow the prompts

#### On Linux (Ubuntu/Debian)

**Using Package Manager:**
```bash
sudo apt-get update
sudo apt-get install golang-go

# Or for latest version:
wget https://go.dev/dl/go1.21.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.linux-amd64.tar.gz
```

#### On Windows

1. Download installer from https://golang.org/dl/
2. Download Windows installer (msi)
3. Run installer
4. Follow the installation wizard
5. Restart your computer

### Verifying Installation

After installation, verify Go is working:

```bash
# Check Go version
go version
# Output: go version go1.21.0 (or newer)

# Check Go environment
go env
# Shows GOROOT, GOPATH, and other settings

# Quick check
go help
# Shows available Go commands
```

### Setting Up Your First Workspace

Go uses a workspace structure. Create a directory for your Go projects:

```bash
# Create workspace directory
mkdir -p ~/go-workspace/golang-course

# Navigate to it
cd ~/go-workspace/golang-course

# Create subdirectories (optional, but good practice)
mkdir src
mkdir bin
mkdir pkg
```

---

## Understanding the Go Workspace

### Go Project Structure

Modern Go uses **modules** (recommended, default in Go 1.11+).

### Workspace Layout (Modern - With Modules)

```
golang-course/
├── go.mod              # Module definition
├── go.sum              # Dependency checksums
├── main.go             # Your code
├── hello/
│   └── hello.go        # Package files
└── README.md           # Documentation
```

### GOPATH vs Go Modules

#### Traditional (GOPATH)
```
$GOPATH/
├── src/
│   ├── github.com/user/project1/
│   └── github.com/user/project2/
├── bin/                # Compiled binaries
└── pkg/                # Compiled packages
```

**Note:** GOPATH is deprecated for new projects.

#### Modern (Go Modules) - RECOMMENDED
```
project-folder/
├── go.mod              # Module declaration
├── go.sum              # Locked dependencies
├── main.go
└── package/
    └── file.go
```

### Key Differences

| Aspect | GOPATH | Go Modules |
|--------|--------|-----------|
| Dependency Management | Manual | Automatic |
| Version Control | Limited | Full support |
| Multiple Versions | Hard | Easy |
| Recommended | ❌ No | ✅ Yes |

---

## Your First Program

### The Classic "Hello, World!"

Let's create your first Go program.

#### Step 1: Create Project Directory

```bash
mkdir ~/golang-course
cd ~/golang-course
```

#### Step 2: Initialize Go Module

```bash
go mod init example.com/hello
```

This creates `go.mod`:
```go
module example.com/hello

go 1.21
```

#### Step 3: Create main.go

Create a file named `main.go`:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

### Understanding the Code

```go
package main          // Package declaration (required for executable)

import "fmt"          // Import package (fmt = formatting)

func main() {         // Entry point (required for executable)
    fmt.Println("Hello, World!")  // Print with newline
}
```

**Line by line:**
1. `package main` - Declares this is the main package (executable program)
2. `import "fmt"` - Imports the fmt package for printing
3. `func main()` - The main function (program starts here)
4. `fmt.Println()` - Prints text followed by a newline

### Key Concepts

#### Packages
- Every Go file must start with a package declaration
- `package main` means this is an executable program
- Other packages are imported libraries

#### Imports
- Import statement includes external packages
- `fmt` package provides formatting functions
- Can import multiple packages

#### Functions
- `func keyword` declares a function
- `main()` is special - it's where the program starts
- No function parameters needed here

#### The Go Way
- No semicolons (added automatically)
- Lowercase for unexported (private)
- Uppercase for exported (public)
- Curly braces on same line (enforced by gofmt)

---

## Running Go Programs

### Method 1: Run Directly (Development)

```bash
go run main.go
```

**Output:**
```
Hello, World!
```

This compiles and runs in one command. Useful for development and testing.

### Method 2: Build Then Run (Production)

```bash
# Step 1: Build (creates executable)
go build

# Step 2: Run the executable
./main          # On macOS/Linux
main.exe        # On Windows
```

**Why two steps?**
- Build creates a reusable binary
- Binary runs faster (pre-compiled)
- Can distribute the binary
- Useful for production deployment

### Method 3: Install (System-wide)

```bash
go install ./...
```

This compiles and places the executable in `$GOPATH/bin` or `$GOBIN`.

### Comparison

| Method | Command | Speed | Use Case |
|--------|---------|-------|----------|
| **Run** | `go run` | Slower | Development |
| **Build** | `go build` | Faster | Production |
| **Install** | `go install` | Faster | System tools |

---

## Go Tools & Commands

### Essential Go Commands

#### 1. **go run** - Run without building
```bash
go run main.go
go run ./...        # Run all Go files in directory
```

#### 2. **go build** - Compile to executable
```bash
go build            # Creates binary in current directory
go build -o myapp   # Specify output name
```

#### 3. **go mod** - Manage dependencies
```bash
go mod init <module>           # Initialize module
go mod add <package>           # Add dependency
go mod tidy                    # Clean dependencies
go mod download                # Download dependencies
```

#### 4. **go get** - Download packages
```bash
go get github.com/user/package
go get -u ./...                # Update all dependencies
```

#### 5. **go fmt** - Format code
```bash
go fmt ./...       # Format all files in directory
gofmt -w main.go   # Format and overwrite
```

#### 6. **go test** - Run tests
```bash
go test ./...      # Run all tests
go test -v ./...   # Verbose output
```

#### 7. **go clean** - Remove build artifacts
```bash
go clean           # Remove binaries and test files
go clean -cache    # Clear build cache
```

#### 8. **go doc** - View documentation
```bash
go doc fmt.Println         # See function documentation
go doc fmt                 # See package documentation
godoc -http=:6060         # Start documentation server
```

### Understanding go.mod and go.sum

#### go.mod
Declares module name and required Go version.

Example:
```go
module example.com/hello

go 1.21

require (
    github.com/google/uuid v1.3.0
)
```

**Contents:**
- Module name (import path)
- Go version minimum
- Direct dependencies

#### go.sum
Records exact versions for reproducible builds.

Example:
```
github.com/google/uuid v1.3.0 h1:fnBlHpVsqrTYL61...
github.com/google/uuid v1.3.0/go.mod h1:6W4AAJF5ea...
```

**Purpose:**
- Security (verify package integrity)
- Reproducibility (same versions each build)
- Version locking

### Useful Command Combinations

```bash
# Initialize and run
go mod init example.com/app && go run main.go

# Format and test
go fmt ./... && go test ./...

# Build and display info
go build && ls -lh

# Clean and rebuild
go clean && go build
```

---

## Troubleshooting

### Common Issues

#### Issue 1: "command not found: go"

**Causes:**
- Go not installed
- Go not in PATH
- Terminal needs restart after installation

**Solutions:**
```bash
# Check if installed
which go           # macOS/Linux
where go           # Windows

# Verify version
go version

# If not in PATH on macOS/Linux, add to ~/.bashrc or ~/.zshrc:
export PATH=$PATH:/usr/local/go/bin
```

#### Issue 2: "cannot find package"

**Example error:**
```
cannot find module providing package fmt
```

**Causes:**
- Module not initialized
- Typo in import
- Dependency not installed

**Solutions:**
```bash
# Initialize module first
go mod init example.com/hello

# Or check imports are correct
# fmt is built-in, no installation needed
```

#### Issue 3: "syntax error near" or "unexpected"

**Causes:**
- Typo in code
- Missing curly brace
- Wrong capitalization

**Solutions:**
```bash
# Check syntax
go fmt main.go    # Format (catches many errors)
go build           # Compile (catches syntax errors)

# Read error message carefully - it shows line number
```

#### Issue 4: Executable won't run

**On macOS/Linux:**
```bash
# Make executable
chmod +x ./main

# Run with ./
./main

# Not just: main
```

**On Windows:**
```bash
# Run with extension
main.exe

# Or full path
C:\path\to\main.exe
```

#### Issue 5: "go: command not found" on Windows

**Causes:**
- Installation incomplete
- Restart needed after installation

**Solutions:**
```bash
# Restart PowerShell or Command Prompt
# Then verify
go version

# Or check installation
C:\Program Files\Go\bin\go.exe version
```

### Debugging Tips

```bash
# See what Go is doing
go build -x                # Show all commands executed

# Verbose mode (tests)
go test -v ./...

# Print environment
go env GOROOT GOPATH GOOS GOARCH

# Clean slate
go clean -cache -testcache
rm -rf go.mod go.sum      # If needed
```

---

## Key Concepts from Level 0 → Go

### Mapping Level 0 to Go

#### Sequence → Code Flow
```go
// Execute in order
fmt.Println("Step 1")
fmt.Println("Step 2")
fmt.Println("Step 3")
```

#### Selection → If Statements
```go
if age >= 18 {
    fmt.Println("You are an adult")
} else {
    fmt.Println("You are a minor")
}
```

#### Repetition → Loops
```go
// For loop
for i := 1; i <= 10; i++ {
    fmt.Println(i)
}
```

#### Functions → Functions
```go
func add(a, b int) int {
    return a + b
}
```

#### Variables → Variables
```go
name := "John"      // Type inferred
age := 25
price := 99.99
```

#### Data Types → Go Types
```go
var i int = 42           // Integer
var f float64 = 3.14     // Float
var s string = "hello"   // String
var b bool = true        // Boolean
```

### Differences from Level 0 Pseudocode

| Concept | Level 0 | Go | Difference |
|---------|---------|-----|-----------|
| **Syntax** | Pseudocode | Real language | Go has specific syntax |
| **Types** | Mentioned | Enforced | Go is strictly typed |
| **Compilation** | N/A | Required | Must compile to run |
| **Libraries** | Not discussed | Packages | Import to use features |
| **Error Handling** | Not detailed | Explicit | Go requires error checks |

---

## Environment Variables

### Key Go Environment Variables

#### GOROOT
Location of Go installation.

```bash
go env GOROOT
# Output: /usr/local/go
```

#### GOPATH
Workspace directory (for Go modules, less important).

```bash
go env GOPATH
# Output: $HOME/go (default)
```

#### GOBIN
Where `go install` places executables.

```bash
go env GOBIN
# Usually: $GOPATH/bin
```

#### GOOS and GOARCH
Operating system and architecture.

```bash
go env GOOS GOARCH
# Output: linux amd64 (on Linux 64-bit)
# Output: darwin arm64 (on Apple Silicon Mac)
```

### Cross-Compilation

Go can compile for different platforms:

```bash
# Compile for Windows (from macOS/Linux)
GOOS=windows GOARCH=amd64 go build

# Compile for macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build

# Compile for Linux
GOOS=linux GOARCH=amd64 go build
```

---

## IDE Setup (Optional but Recommended)

### Recommended Editors

#### 1. **VS Code** (Recommended)
```bash
# Install VS Code
# Install Go extension: golang.go

# Inside VS Code:
# Go will suggest installing tools
# Accept the installation
```

#### 2. **GoLand** (JetBrains)
- Full-featured IDE for Go
- Free trial available
- Excellent code completion

#### 3. **Vim/Neovim**
- With vim-go plugin
- For terminal lovers

#### 4. **Sublime Text**
- Lightweight
- With GoSublime plugin

### VS Code Setup Steps

1. **Install VS Code**
   ```bash
   brew install --cask visual-studio-code
   ```

2. **Install Go Extension**
   - Open VS Code
   - Go to Extensions (Ctrl+Shift+X / Cmd+Shift+X)
   - Search "Go"
   - Install "Go" by golang

3. **Configure Go Tools**
   - VS Code will prompt to install tools
   - Click "Install All"
   - Wait for installation

4. **Verify**
   - Create a `.go` file
   - VS Code should show syntax highlighting
   - Code completion should work

---

## ✅ Level 1 Checklist

By the end of this level, you should:

- [ ] Understand what Go is and why it's useful
- [ ] Know Go's core philosophy
- [ ] Have Go installed and verified
- [ ] Understand Go workspace and project structure
- [ ] Created a `go.mod` file
- [ ] Written and run your first Go program
- [ ] Understand the difference between `go run` and `go build`
- [ ] Know key Go commands (run, build, fmt, test, mod)
- [ ] Set up your preferred editor
- [ ] Know how to troubleshoot common issues
- [ ] Understand how Level 0 concepts map to Go

---

## 🎯 Quick Reference: Go Setup Checklist

```bash
# Step 1: Verify Go installation
go version

# Step 2: Create project directory
mkdir ~/my-go-project
cd ~/my-go-project

# Step 3: Initialize module
go mod init example.com/hello

# Step 4: Create main.go
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
EOF

# Step 5: Run program
go run main.go

# Step 6: Build executable
go build

# Step 7: Run executable
./main          # macOS/Linux
main.exe        # Windows
```

---

## 🚀 Ready for Level 2?

Once you've completed Level 1:
- ✅ Your environment is set up
- ✅ You can run Go programs
- ✅ You understand the Go workspace
- ✅ You're ready to learn Go syntax

**Next Level:** Level 2 - Go Program Structure
- Package organization
- Import statements
- Main function deep dive
- Multiple files in a project

---

## 📚 Additional Resources

- **Official Go Site:** https://golang.org
- **Go Documentation:** https://pkg.go.dev/
- **Go Playground:** https://play.golang.org (write code online)
- **Effective Go:** https://golang.org/doc/effective_go (best practices)
- **Go Blog:** https://go.dev/blog (announcements and articles)

---

## 💡 Key Takeaways

1. **Go is Simple:** Easy syntax, fast to learn
2. **Go is Fast:** Compiled language, high performance
3. **Go is Practical:** Built for real-world problems
4. **Go is Concurrent:** Goroutines for parallel processing
5. **Go is Productive:** Fast compilation and deployment

---

## Next Steps

1. **Install Go** (if not already done)
2. **Create a project directory**
3. **Initialize with `go mod init`**
4. **Write "Hello, World!" program**
5. **Practice `go run` and `go build`**
6. **Set up your editor**
7. **Move to Level 2**

---

## Questions? 

If you encounter issues:
1. Check the [Troubleshooting](#troubleshooting) section
2. Review the setup steps
3. Check the [Quick Reference](#-quick-reference-go-setup-checklist)
4. Visit [golang.org](https://golang.org) for official docs

**You're now ready to write Go code!** 🚀
