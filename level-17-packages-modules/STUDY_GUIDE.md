# Level 17: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: Modules In Depth
```
Day 1:  Recap packages/modules; go.mod anatomy (module, go, require, replace, exclude)
Day 2:  Semantic versioning as Go modules interpret it (v0/v1/v2+)
Day 3:  Practical commands: go mod tidy, go list -m all, go mod why, go mod graph
Day 4:  Building a local multi-module project with replace
Day 5:  internal/ packages and the visibility they enforce
Day 6:  Package-level API design; avoiding circular imports
Day 7:  Multi-module workspaces with go.work
```

### Week 2: Practice & Application
```
Day 1:  Exercises 1-3 (go.mod anatomy, two-module project, tidy/why/graph)
Day 2:  Exercises 4-6 (internal/ allowed + error, circular imports)
Day 3:  Exercises 7-8 (API refactor, go.work workspace)
Day 4:  Exercises 9-10 (replace vs go.work, comprehensive practice)
Day 5:  Bonus challenges
Day 6-7: Practice & review
```

---

## 🗂️ go.mod Anatomy Diagram

```
module example.com/widgets          ← identity + import-path prefix
                                       for every package inside this module
go 1.26                             ← minimum language version required

require (                           ← direct + indirect dependencies,
    example.com/greetings v1.2.0        each pinned to an exact version
    example.com/legacy v0.9.0           (computed by Minimum Version Selection)
)

replace example.com/greetings       ← "use THIS instead" — main module
    => ../greetings                    only; a local path never touches
                                        the network

exclude example.com/legacy v0.8.0   ← "never select this exact version" —
                                        main module only
```

```
Minimum Version Selection (MVS), in one line:

  For every module in the graph, pick the HIGHEST version anything
  requires — never a version nobody asked for, never silently "latest".
```

---

## 🔢 Semantic Versioning Table

| Version range | Meaning | Compatibility promise | Import path |
|---|---|---|---|
| `v0.x.y` | Pre-release / unstable | None — anything can change | `example.com/widgets` |
| `v1.x.y` | Stable, first major release | Backward compatible within v1 | `example.com/widgets` |
| `v2.x.y` | Stable, breaking release | Backward compatible within v2 only | `example.com/widgets/v2` |
| `v3.x.y`+ | Same pattern, repeats | Backward compatible within vN only | `example.com/widgets/v3` |

```
Why the /vN suffix exists:

  v1 and v2+ are DECLARED incompatible, so Go treats them as
  DIFFERENT modules that can be imported side by side:

      example.com/widgets      (v1, still used transitively somewhere)
      example.com/widgets/v2   (v2, used directly)

  Both can appear in the same build at once — no conflict.
```

```
Pseudo-version, decoded:

  v0.0.0-00010101000000-000000000000
   │      │              └─ commit hash (abbreviated)
   │      └─ commit timestamp (UTC, YYYYMMDDHHMMSS)
   └─ base version (no matching tag exists yet)

  Go generates this automatically for a require line pointing at
  an untagged commit — you'll see it every time you `go mod edit
  -require` a local module that was never tagged.
```

---

## ⚙️ Practical Module Commands Cheat Sheet

| Command | Answers | Verified real output shape |
|---|---|---|
| `go mod init <path>` | "Create a new module here" | `go: creating new go.mod: module <path>` |
| `go mod tidy` | "Sync go.mod/go.sum with actual imports" | Silent when nothing needs to change |
| `go list -m all` | "What's the full resolved build list?" | One line per module, `=>` shows an active replace |
| `go mod why <pkg>` | "What import chain pulled this in?" | `# <pkg>` then the chain, top to bottom |
| `go mod graph` | "What's the raw dependency graph?" | `<requirer> <required>@<version>`, one edge per line |

---

## 🔗 replace vs. go.work Comparison Diagram

```
                    replace (in go.mod)              go.work (separate file)
                    ─────────────────────             ───────────────────────
Lives in            The consumer's go.mod             Its own go.work file
                    (normally committed)               (normally NOT committed)

Scope               Only the main module's             Every module listed with
                    own build                          `use`, for every command run
                                                        from the workspace

Needs a require     Yes — still needs SOME              No — workspace resolution
line too?           version placeholder                replaces the need entirely

Typical use         "Point this ONE dependency          "I'm actively developing 2+
                    at a fork or local path,            modules together and want
                    temporarily"                        them to see each other live"

Set up with         go mod edit -require=...            go work init
                    go mod edit -replace=...            go work use ./a ./b
```

```
The one thing BOTH respect equally:

  internal/ visibility. Neither replace nor go.work makes an
  internal/ package importable from outside its module tree —
  they only make packages FINDABLE across modules, never MORE
  VISIBLE. (Verified in Exercise 5 and again in Exercise 10.)
```

---

## 🔒 internal/ Visibility Diagram

```
example.com/owner/                    Allowed to import
├── go.mod                            example.com/owner/internal/secret:
├── cmd/
│   └── tool/
│       └── main.go   ─────────✅────→   anything rooted at example.com/owner
├── internal/                              (the parent of "internal")
│   └── secret/
│       └── secret.go
└── ...

example.com/outsider/                 NOT allowed — different module,
├── go.mod                            outside the owner/ tree:
└── main.go   ──────────────────✗────→   compile error, even with a
                                            replace pointing right at it:

  main.go:6:5: use of internal package
  example.com/owner/internal/secret not allowed
```

```
The rule, generalized:

  A package under .../internal/... is importable ONLY by code
  rooted at the directory that is internal's PARENT.

  This nests at any depth:
    myapp/internal/...            → visible to all of myapp
    myapp/cmd/tool/internal/...   → visible only to myapp/cmd/tool
```

---

## 🧩 Package Design Principles

```go
// ✅ Cohesive — one job, small surface
package textutil

func Shout(s string) string   { ... }   // exported: the API
func Reverse(s string) string { ... }   // exported: the API

// ❌ Low cohesion — grab-bag of unrelated helpers
package toolbox

func Shout(s string) string    { ... }  // string helper
func Max(a, b int) int         { ... }  // ...next to a math helper
func ParseDate(s string) error { ... }  // ...next to a date helper
```

```
Minimize the exported surface:

  Every exported name is a PROMISE. Before capitalizing a name,
  ask: does a caller outside this package actually need it, or
  is it an implementation detail?

  Small exported surface  → easy to document, easy to keep stable
  Large exported surface  → every refactor risks a breaking change
```

```
Avoiding circular imports (the DAG rule):

  Go's import graph — across packages AND across modules — must
  have NO CYCLES. If A needs something from B and B needs
  something from A, fix it with one of:

    1. Pass plain VALUES instead of importing the other's type
       (order.Cost float64, not the whole Order struct)
    2. Extract a shared package with NO dependencies of its own
    3. Define an interface in the "downstream" package that the
       "upstream" package satisfies
```

---

## 🚨 Common Mistakes

### Mistake 1: Committing a Local-Only replace Path

```
// ❌ WRONG - only works on your machine
replace example.com/greetings => /Users/yourname/dev/greetings
```

### Mistake 2: Forgetting the /vN Suffix on a Breaking Release

```
// ❌ WRONG - a breaking v2 release still using the v1 module path
module example.com/widgets

// ✅ RIGHT
module example.com/widgets/v2
```

### Mistake 3: Expecting go.work to Travel With the Repo

```
// ❌ WRONG ASSUMPTION - go.work paths are usually specific to YOUR checkout
```

### Mistake 4: Trying to Route Around internal/

```
// ❌ WRONG - no replace or go.work makes this importable from outside the tree
import "example.com/owner/internal/secret"
```

### Mistake 5: One Catch-All utils Package Forever

```
// ❌ WRONG - absorbs everything, becomes a dependency of nearly everything
package utils

// ✅ RIGHT - split by responsibility
package dateutil
package mathutil
package validate
```

### Mistake 6: Assuming go mod tidy Needs the Network for a replace'd Module

```
// ❌ WRONG ASSUMPTION
// Once replace/go.work resolves a path to a local directory, tidy/build/run
// read it straight off disk — no network call happens, ever, for that module.
```

---

## 📈 Progression Summary

### Understanding Level 17

Level 17 goes from "one module, several packages" (Level 2) to "several modules, each with a deliberate API," entirely offline:

1. **go.mod anatomy** — every directive and what it does
2. **Semantic versioning** — how Go encodes compatibility promises, including the v2+ import-path rule
3. **Module commands** — tidy, list -m all, why, graph, all run for real
4. **Local multi-module projects** — replace directives standing in for a real registry
5. **internal/** — compiler-enforced, module-tree-scoped visibility
6. **API design** — cohesion, minimal exported surface, no cycles
7. **go.work** — the cleaner tool for actively developing several local modules together

### Prerequisites for Level 18

Before moving to Level 18 (File Handling), you need:

- ✅ Comfortable reading and writing every go.mod directive
- ✅ Can build a two-module local project with a `replace` directive, offline
- ✅ Can read `go list -m all`, `go mod why`, and `go mod graph` output
- ✅ Understand why `internal/` imports fail outside their module tree
- ✅ Can split a low-cohesion package into focused ones
- ✅ Can set up and use a `go.work` workspace across local modules

### Ready for Level 18?

Level 18 shifts from code organization to real I/O — reading and writing files, working with paths and directories, and handling the file-specific errors that Level 16's error handling prepared you for.

---

## ✅ Checklist Before Level 18

- [ ] Can explain every line of a go.mod from memory
- [ ] Can explain v0 vs v1 vs v2+ and why v2+ needs a path suffix
- [ ] Have run `go mod tidy`, `go list -m all`, `go mod why`, and `go mod graph` for real
- [ ] Built a working two-module project wired together with `replace`, fully offline
- [ ] Reproduced the real `internal/` compile error from outside a module tree
- [ ] Caused and fixed a real circular import
- [ ] Refactored a monolithic package into two cohesive ones
- [ ] Set up a `go.work` workspace and built across it with no `replace`
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The Module
A module is a unit of versioning; a package is a unit of organization inside it. `go.mod` is the module's complete contract: identity, minimum Go version, dependencies, and any local overrides.

### The Version
Go modules use semantic versioning as a *promise*, not just a label — v0 promises nothing, v1 promises compatibility within v1, and v2+ needs its own import path because it's treated as a wholly separate module.

### The Boundary
`internal/` is the one visibility rule the compiler enforces *between* packages, not just within one. It nests at any depth and cannot be bypassed with `replace` or `go.work`.

### The Design
A good package does one job, exports the smallest surface that does it, and never creates an import cycle - across packages or across modules.

### The Workflow
`replace` is the quick, committable override for one dependency. `go.work` is the ongoing, local, uncommitted setup for developing several modules together. Both work entirely offline once they point at a directory on disk.

---

## 📚 Next Level

Level 18: File Handling
- Reading and writing files
- Working with paths and directories
- Buffered I/O
- Handling file-related errors

You've got real project structure down! Keep going! 🚀
