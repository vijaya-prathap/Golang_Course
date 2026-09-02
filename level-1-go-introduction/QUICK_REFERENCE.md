# Level 1: Quick Reference Card

## 🚀 Quick Start (5 Minutes)

```bash
# 1. Verify Go
go version

# 2. Create project
mkdir my-app && cd my-app

# 3. Initialize module
go mod init example.com/app

# 4. Write program
cat > main.go << 'EOF'
package main
import "fmt"
func main() {
    fmt.Println("Hello, Go!")
}
EOF

# 5. Run
go run main.go

# 6. Build
go build

# 7. Run executable
./main          # macOS/Linux
main.exe        # Windows
```

---

## 📋 Essential Commands

```bash
# Run & Test
go run main.go              # Compile & run
go run ./...                # Run all Go files
go test ./...               # Run tests

# Build
go build                    # Create executable
go build -o myapp           # Custom name
go clean                    # Remove artifacts

# Format & Quality
go fmt ./...                # Format all files
go vet ./...                # Find bugs
go mod tidy                 # Clean dependencies

# Modules
go mod init example.com/app # Create module
go get github.com/user/pkg  # Add dependency
go mod graph                # Show dependencies

# Documentation
go doc fmt                  # Package docs
go doc fmt.Println          # Function docs
go help                     # List commands
```

---

## 📁 Project Structure

```
project-name/
├── go.mod              # Module definition
├── main.go             # Program entry point
└── README.md           # Documentation
```

**go.mod contents:**
```go
module example.com/app

go 1.21
```

---

## 💻 Program Template

```go
package main           // Executable package

import "fmt"           // Imports

func main() {          // Entry point
    fmt.Println("Hello!")
}
```

---

## 🔧 Installation Verification

| Check | Command | Expected |
|-------|---------|----------|
| Installed? | `go version` | `go version go1.21...` |
| GOROOT | `go env GOROOT` | `/usr/local/go` |
| Can run? | `go help` | List of commands |

---

## 🔄 Development Loop

```
1. Write   → Edit main.go
2. Format  → go fmt ./...
3. Run     → go run main.go
4. Build   → go build
5. Test    → ./main (or main.exe)
```

---

## 📦 Built-in Packages

| Package | Use | Example |
|---------|-----|---------|
| `fmt` | Print/format | `fmt.Println()` |
| `os` | OS operations | `os.Exit()` |
| `time` | Time/date | `time.Now()` |
| `math` | Math functions | `math.Sqrt()` |
| `strings` | Text operations | `strings.Split()` |

---

## 🎯 Go Philosophy (Remember These!)

1. **Simplicity** - Simple syntax, 25 keywords
2. **Readability** - Code read > code written
3. **Explicitness** - No magic, clear intent
4. **Concurrency** - Goroutines built-in
5. **Productivity** - Fast compile, deploy

---

## ⚠️ Common Errors

| Error | Cause | Fix |
|-------|-------|-----|
| `cannot find package fmt` | Missing import | Add `import "fmt"` |
| `undefined: main` | Wrong package | Use `package main` |
| `syntax error` | Code error | Check brackets/commas |
| `file not found` | Wrong directory | Run from correct folder |
| `command not found: go` | Not installed | Install Go |

---

## 🛠️ Troubleshooting Quick Fixes

**Go not recognized:**
```bash
export PATH=$PATH:/usr/local/go/bin
```

**Can't run file:**
```bash
# Make executable
chmod +x main

# Run with ./
./main
```

**Syntax error:**
```bash
# Check formatting
go fmt main.go

# Check compilation
go build
```

**Module issues:**
```bash
# Recreate module
rm go.mod go.sum
go mod init example.com/app
```

---

## 📚 Files to Know

| File | Purpose |
|------|---------|
| `main.go` | Program code |
| `go.mod` | Module definition |
| `go.sum` | Dependency lock file |
| `README.md` | Documentation |

---

## 🔗 Essential Links

- **Install:** https://golang.org/dl
- **Docs:** https://pkg.go.dev
- **Playground:** https://play.golang.org
- **Effective Go:** https://golang.org/doc/effective_go

---

## 🎓 Checklist: Ready for Level 2?

- [ ] `go version` shows version
- [ ] Can run `go run main.go`
- [ ] Can build with `go build`
- [ ] Created at least 3 programs
- [ ] Understand `package main`
- [ ] Know `go fmt`, `go build`, `go run`
- [ ] Can modify programs
- [ ] Editor is set up (optional)

---

## 💡 Remember

- **Files must start with:** `package main` or `package name`
- **Executables need:** `func main() { }`
- **Imports go:** Inside the file, after package declaration
- **Curly braces:** Go style brace placement (same line)
- **Formatting:** Use `go fmt` - it's the standard

---

## Next: Level 2 - Go Program Structure

You'll learn:
- Multiple files and packages
- Proper project organization
- Import statements
- Public/private identifiers

🚀 **Keep coding!**
