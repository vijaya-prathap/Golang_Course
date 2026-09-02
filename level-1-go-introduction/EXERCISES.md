# Level 1: Hands-On Exercises & Setup Guide

## 🎯 Practical Setup Exercises

Follow these exercises in order to ensure your Go environment is properly set up.

---

## Exercise 1: Verify Go Installation

### Objective
Confirm that Go is properly installed on your system.

### Steps

```bash
# 1. Check Go version
go version

# Expected output: go version go1.21.0 (or newer)

# 2. Check Go environment
go env

# 3. Check GOPATH
go env GOPATH

# 4. Check GOROOT
go env GOROOT
```

### What Should You See?

```bash
$ go version
go version go1.21.0 linux/amd64

$ go env GOPATH
/home/user/go

$ go env GOROOT
/usr/local/go
```

### Troubleshooting

**Problem:** "command not found: go"
**Solution:**
```bash
# Check if in PATH
export PATH=$PATH:/usr/local/go/bin

# On macOS with Homebrew
brew install go

# Restart terminal and try again
```

**Problem:** Version too old (< 1.16)
**Solution:**
```bash
# Update Go
brew upgrade go          # macOS
apt-get install golang   # Linux
# Or download from golang.org
```

---

## Exercise 2: Create Your First Project Directory

### Objective
Set up a proper Go project workspace.

### Steps

```bash
# 1. Create project directory
mkdir -p ~/workspace/golang-course
cd ~/workspace/golang-course

# 2. Verify you're in correct directory
pwd
# Output: /Users/username/workspace/golang-course
# (path may differ on your system)

# 3. List directory (should be empty)
ls -la
# Output: total 0
```

### Directory Structure
After this exercise:
```
~/workspace/golang-course/
└── (empty, ready for project)
```

---

## Exercise 3: Initialize a Go Module

### Objective
Create a Go module for dependency management.

### Steps

```bash
# 1. Navigate to project directory
cd ~/workspace/golang-course

# 2. Initialize module
go mod init example.com/hello

# 3. Verify go.mod was created
cat go.mod
```

### Expected Output

```go
module example.com/hello

go 1.21
```

### What Was Created?

```
~/workspace/golang-course/
└── go.mod          (Module definition file)
```

### Understanding Module Name

The module name `example.com/hello` consists of:
- `example.com` - Your organization/domain (can be anything for learning)
- `/hello` - Project name

**For real projects:**
```
github.com/username/project-name
```

### Alternative Module Names
```bash
# For learning
go mod init learn.golang/hello

# For your laptop (no domain needed)
go mod init hello-world

# GitHub style
go mod init github.com/yourusername/my-app
```

---

## Exercise 4: Create Your First Program

### Objective
Write and run the classic "Hello, World!" program.

### Steps

#### Step 1: Create main.go file

```bash
# Option A: Using cat (simple)
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
EOF

# Option B: Using your editor
code main.go  # VS Code
vim main.go   # Vim
nano main.go  # Nano
```

#### Step 2: View the file
```bash
cat main.go
```

Expected output:
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

### Directory Structure Now
```
~/workspace/golang-course/
├── go.mod
└── main.go
```

---

## Exercise 5: Run Your Program

### Objective
Execute your first Go program.

### Steps

```bash
# 1. Make sure you're in project directory
cd ~/workspace/golang-course

# 2. Run the program using go run
go run main.go

# 3. You should see:
# Hello, World!
```

### What Happened?
```
go run main.go
   ↓
Compile to temporary binary
   ↓
Execute the binary
   ↓
Print "Hello, World!"
   ↓
Clean up temporary binary
```

### Verify It Worked
- Did you see "Hello, World!" in the terminal?
- If YES → Great! Move to Exercise 6
- If NO → Check [Troubleshooting](#troubleshooting-exercise-5)

### Troubleshooting Exercise 5

**Problem:** "cannot find package fmt"
```
solution: package fmt is built-in, verify your go.mod exists
check: ls -la go.mod
```

**Problem:** "syntax error"
```
solution: Check your code matches exactly
check: Spelling, capitalization, curly braces
```

**Problem:** "no Go files in /path"
```
solution: Make sure main.go exists in current directory
check: ls -la main.go
```

---

## Exercise 6: Build an Executable

### Objective
Create a standalone executable binary.

### Steps

```bash
# 1. Build the program
go build

# 2. List directory contents
ls -la

# Expected output:
# -rwxr-xr-x  main
# -rw-r--r--  go.mod
# -rw-r--r--  main.go
# -rw-r--r--  go.sum (optional)
```

### Run the Executable

```bash
# macOS/Linux
./main

# Windows
main.exe

# Expected output:
# Hello, World!
```

### What Happened?
```
go build
   ↓
Compile to binary named 'main' (default from package main)
   ↓
Place in current directory
   ↓
Binary is ready to run
```

### Directory Structure Now
```
~/workspace/golang-course/
├── go.mod
├── go.sum        (created by build)
├── main          (compiled executable) ← NEW!
└── main.go
```

---

## Exercise 7: Customize Build Output

### Objective
Build with a custom output name.

### Steps

```bash
# 1. Build with custom name
go build -o hello

# 2. List directory
ls -la

# Expected output shows 'hello' instead of 'main'

# 3. Run the executable
./hello

# Output: Hello, World!
```

### Why Customize?

| Default | Customized | Reason |
|---------|-----------|--------|
| `go build` → `main` | `go build -o hello` | Better names for final executables |
| Confusing | Clear | User knows what it does |

### Useful Flags

```bash
# Custom name
go build -o myapp

# Windows executable
GOOS=windows GOARCH=amd64 go build -o myapp.exe

# Optimized (no debug info)
go build -ldflags="-s -w"

# Version embedded
go build -ldflags="-X main.Version=1.0.0"
```

---

## Exercise 8: Format Your Code

### Objective
Learn Go's code formatting standard.

### Steps

```bash
# 1. View current code
cat main.go

# 2. Run formatter
go fmt

# 3. Verify formatting
cat main.go
```

### Why gofmt?

Go enforces a single formatting style:
- Consistent across all Go code
- No debates about spaces vs tabs
- Readable and professional

### Before and After

**Before (improper spacing):**
```go
package   main
import   "fmt"
func main(  ){
fmt.Println(   "Hello, World!" )
}
```

**After (go fmt):**
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

### Format Multiple Files

```bash
# Format all Go files in directory
go fmt ./...

# See what would change (dry-run)
gofmt -d main.go

# Format and overwrite
gofmt -w main.go
```

---

## Exercise 9: Modify Your Program

### Objective
Make changes to your Go program to verify it works.

### Steps

```bash
# 1. Open main.go in your editor
code main.go

# 2. Modify the output
# Change: fmt.Println("Hello, World!")
# To:     fmt.Println("Hello, Golang!")
```

**Updated main.go:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Golang!")
}
```

```bash
# 3. Run to see changes
go run main.go

# Output: Hello, Golang!

# 4. Build new executable
go build -o hello

# 5. Run executable
./hello

# Output: Hello, Golang!
```

### Directory Structure
```
~/workspace/golang-course/
├── go.mod
├── go.sum
├── hello           (updated executable)
├── main            (old executable)
└── main.go         (modified)
```

---

## Exercise 10: Print Multiple Lines

### Objective
Understand Go's printing functions.

### Steps

```bash
# 1. Update main.go
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("Welcome to Go!")
    fmt.Println("This is line 2")
    fmt.Println("This is line 3")
    fmt.Println("")
    fmt.Println("Each println adds a newline")
}
EOF

# 2. Run the program
go run main.go
```

### Expected Output
```
Welcome to Go!
This is line 2
This is line 3

Each println adds a newline
```

### Different Printing Functions

```go
fmt.Print("No newline")          // No newline
fmt.Println("With newline")      // Adds newline
fmt.Printf("Format: %d\n", 42)   // Formatted output
```

---

## Exercise 11: Understand Packages

### Objective
Learn how packages work in Go.

### Steps

```bash
# 1. View current code (it uses the fmt package)
cat main.go

# 2. Check available fmt functions
go doc fmt

# 3. View specific function
go doc fmt.Println

# 4. Online documentation
# Visit: https://pkg.go.dev/fmt
```

### Important Packages

| Package | Purpose | Examples |
|---------|---------|----------|
| `fmt` | Input/output | Print, formatting |
| `math` | Math functions | Sin, Cos, Sqrt |
| `strings` | String manipulation | Split, Replace, Trim |
| `time` | Time and duration | Clock, Sleep, Timer |
| `os` | Operating system | File, Env, Args |

---

## Exercise 12: Environment Verification

### Objective
Verify your complete Go environment.

### Checklist Script

Run this script to verify everything:

```bash
#!/bin/bash

echo "=== Go Environment Check ==="
echo ""

echo "1. Go Version:"
go version
echo ""

echo "2. GOROOT (Go installation):"
go env GOROOT
echo ""

echo "3. GOPATH (Workspace):"
go env GOPATH
echo ""

echo "4. GOOS (Operating System):"
go env GOOS
echo ""

echo "5. GOARCH (Architecture):"
go env GOARCH
echo ""

echo "6. Project directory:"
pwd
echo ""

echo "7. Project files:"
ls -la
echo ""

echo "8. Module info:"
cat go.mod
echo ""

echo "=== All Checks Complete ==="
```

Save as `verify.sh`:
```bash
chmod +x verify.sh
./verify.sh
```

---

## Exercise 13: Clean Up Build Artifacts

### Objective
Learn to clean unnecessary files.

### Steps

```bash
# 1. See current files
ls -la

# 2. Clean build artifacts
go clean

# 3. Check what remains
ls -la
```

### What go clean Removes

```
Before:
├── go.mod
├── go.sum
├── hello       ← Removed
├── main        ← Removed
└── main.go

After:
├── go.mod
├── go.sum
└── main.go
```

### Useful Clean Commands

```bash
# Remove binaries only
go clean -i

# Remove build cache
go clean -cache

# Remove everything
go clean -r
```

---

## Exercise 14: Set Up Your Editor

### Objective
Configure your editor for Go development.

### For VS Code

```bash
# 1. Install VS Code (if needed)
brew install --cask visual-studio-code

# 2. Launch VS Code
code .

# 3. Install Go extension:
#    - Press Cmd+Shift+X (Extensions)
#    - Search "Go"
#    - Install "Go" by golang

# 4. Accept tool installation when prompted
#    - Click "Install All" for Go tools

# 5. Restart VS Code
```

### VS Code Features

After setup, you'll have:
- ✅ Syntax highlighting
- ✅ Code completion
- ✅ Go to definition
- ✅ Error squiggles
- ✅ Format on save
- ✅ Debug support

### For Other Editors

**Vim:**
```bash
git clone https://github.com/vim-jp/vim-go.git
```

**Sublime Text:**
```bash
# Install Package Control, then GoSublime
```

---

## Exercise 15: First Project Summary

### Objective
Complete your first Go project and document it.

### Create README.md

```bash
cat > README.md << 'EOF'
# My First Go Program

## Description
This is my first Go program, created while learning Go basics.

## Files
- `main.go` - Main program
- `go.mod` - Module definition
- `go.sum` - Dependency checksums
- `README.md` - This file

## How to Run

### Using go run (Development)
```bash
go run main.go
```

### Build and Run (Production)
```bash
go build
./main
```

## Output
```
Hello, Golang!
This is line 2
This is line 3

Each println adds a newline
```

## Learning Outcomes
- ✅ Installed Go
- ✅ Created Go project
- ✅ Wrote first program
- ✅ Ran and built program
- ✅ Formatted code
- ✅ Used fmt package

## Next Steps
- Learn Go syntax
- Create multiple files
- Use packages
- Handle errors
EOF

cat README.md
```

### Final Project Structure

```
my-go-project/
├── README.md       ← Documentation
├── go.mod          ← Module definition
├── go.sum          ← Dependency checksums
└── main.go         ← Program code
```

---

## ✅ Exercise Completion Checklist

Mark off each exercise as you complete it:

- [ ] Exercise 1: Verify Go Installation
- [ ] Exercise 2: Create Project Directory
- [ ] Exercise 3: Initialize Go Module
- [ ] Exercise 4: Create First Program
- [ ] Exercise 5: Run Your Program
- [ ] Exercise 6: Build an Executable
- [ ] Exercise 7: Customize Build Output
- [ ] Exercise 8: Format Your Code
- [ ] Exercise 9: Modify Your Program
- [ ] Exercise 10: Print Multiple Lines
- [ ] Exercise 11: Understand Packages
- [ ] Exercise 12: Environment Verification
- [ ] Exercise 13: Clean Up Build Artifacts
- [ ] Exercise 14: Set Up Your Editor
- [ ] Exercise 15: First Project Summary

---

## 🎯 Bonus Challenges

### Challenge 1: Print a Box
Create a program that prints a box with your name:

```
Expected output:
┌─────────────────┐
│  Your Name      │
└─────────────────┘
```

### Challenge 2: Multiple Programs
Create two separate Go files and run them both with `go run ./...`

### Challenge 3: Cross-Compile
Build an executable for a different operating system:
```bash
GOOS=windows GOARCH=amd64 go build -o hello.exe
```

### Challenge 4: Documentation
Visit https://pkg.go.dev/ and find 3 packages you want to use

### Challenge 5: Make an Alias
Add a bash alias for faster development:
```bash
alias gorun='go run main.go'
alias gobuild='go build && ./main'
```

---

## 🚀 Ready for Level 2?

Once you've completed all 15 exercises:
- ✅ Go is properly installed
- ✅ You can create projects
- ✅ You can run and build programs
- ✅ Your editor is configured
- ✅ You understand the Go workspace

**Next Level:** Level 2 - Go Program Structure
- Multiple packages
- Import statements
- Project organization
- Public and private identifiers

---

## 📞 Need Help?

Stuck on an exercise?

1. **Re-read the exercise steps**
2. **Check error messages** - they usually help!
3. **Verify file names** - typos are common
4. **Use `go doc`** - get help on packages
5. **Check your code formatting** - run `go fmt`

Remember: Errors are learning opportunities!

Happy coding! 🎉
