# Level 17: Quick Reference Card

## 🚀 Quick Start: Two Local Modules, Wired Offline (5 Minutes)

```bash
# Library module
mkdir -p greetings && cd greetings
go mod init example.com/greetings
cat > greetings.go << 'EOF'
package greetings

import "fmt"

func Hello(name string) string {
    return fmt.Sprintf("Hello, %s!", name)
}
EOF
cd ..

# App module
mkdir -p app && cd app
go mod init example.com/app
cat > main.go << 'EOF'
package main

import (
    "fmt"

    "example.com/greetings"
)

func main() {
    fmt.Println(greetings.Hello("Gopher"))
}
EOF

# Wire it up - no network needed
go mod edit -require=example.com/greetings@v0.0.0-00010101000000-000000000000
go mod edit -replace=example.com/greetings=../greetings
go run .
```

---

## 📋 go.mod Directives

```
module example.com/widgets     // this module's identity + import prefix
go 1.26                        // minimum language version
require example.com/x v1.2.0   // a dependency and its resolved version
replace example.com/x => ../x  // redirect x to a local path (main module only)
exclude example.com/x v1.0.0   // blacklist one version (main module only)
```

---

## 🔢 Semantic Versioning

```
v0.x.y   → no compatibility promise
v1.x.y   → stable, compatible within v1
v2+.x.y  → stable, compatible within vN only — NEEDS a /vN suffix:

    module example.com/widgets/v2
    import "example.com/widgets/v2/render"
```

---

## ⚙️ Module Commands

```bash
go mod init <path>        # create a new go.mod
go mod tidy                # sync go.mod/go.sum with actual imports
go mod tidy -v             # same, but verbose (silent = nothing changed)
go list -m all             # print the full resolved build list
go mod why <pkg>           # why is this package/module in the build?
go mod why -m <module>     # same question, at the module level
go mod graph               # raw dependency graph: requirer required@version
go mod edit -require=...   # script a go.mod edit
go mod edit -replace=...   # script a go.mod edit
go mod edit -json          # print go.mod as structured JSON
```

---

## 🔒 internal/ Visibility

```
example.com/owner/
├── cmd/tool/        ← CAN import example.com/owner/internal/x
└── internal/x/

example.com/outsider/
└── main.go          ← CANNOT import example.com/owner/internal/x
                        (real compile error, even with replace/go.work):

    use of internal package example.com/owner/internal/x not allowed
```

Rule: visible only to code rooted at `internal`'s **parent directory**. Nests at any depth.

---

## 🧩 Package API Design

```go
// ✅ Cohesive, minimal surface
package cache
func New(capacity int) *Cache { ... }        // exported entry point
func evict(c *Cache) { ... }                 // unexported detail

// ❌ Grab-bag, everything exported
package toolbox
func Shout(s string) string { ... }
func Max(a, b int) int { ... }
```

```
Avoiding circular imports:
  1. Pass a plain value instead of the other package's type
  2. Extract a shared, dependency-free package
  3. Define an interface the "upstream" package satisfies
```

---

## 🗂️ go.work Workspace

```bash
mkdir workspace && cd workspace
# (greetings/ and app/ already exist as sibling modules)
go work init
go work use ./greetings ./app
go run ./app            # works - no require/replace needed between members
```

```
go.work:

    go 1.26

    use (
        ./app
        ./greetings
    )
```

⚠️ `go build ./...` from the **workspace root itself** fails - it isn't a module. Run it from inside a member directory, or target one: `go build ./app/...`.

---

## 🔀 replace vs. go.work

| | `replace` | `go.work` |
|---|---|---|
| Where | In go.mod (usually committed) | Separate file (usually not committed) |
| Scope | Main module only | Every `use`d module |
| Needs a `require` too | Yes | No |
| Best for | One-off local override | Ongoing multi-module development |

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Local-only replace path | `replace x => /Users/me/x` (committed) | Keep it local, or use `go.work` instead |
| Missing v2+ suffix | `module x` at v2.0.0 | `module x/v2` |
| Committing go.work | Sharing machine-specific paths | Treat it as local, uncommitted |
| Routing around internal/ | Trying replace/go.work to reach it | Not possible - ask for it to be exported |
| One giant utils package | Everything in `package utils` | Split into `dateutil`, `mathutil`, etc. |

---

## 🎓 Before Next Level

Can you:
- [ ] Explain every line of a go.mod from memory?
- [ ] Explain why v2+ modules need a `/vN` import-path suffix?
- [ ] Build two local modules wired together with `replace`, fully offline?
- [ ] Read `go list -m all`, `go mod why`, and `go mod graph` output?
- [ ] Reproduce the real `internal/` compile error?
- [ ] Split a monolithic package into cohesive ones?
- [ ] Set up a `go.work` workspace across two local modules?

If YES → You're ready for Level 18!

---

## 📚 Next Level

Level 18: File Handling
- Reading and writing files
- Working with paths and directories
- Buffered I/O

You've got real project structure down! 💪
