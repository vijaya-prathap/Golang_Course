# Level 17: Packages & Modules - Complete Guide

## Introduction

Welcome to Level 17! You've mastered error handling (Level 16) - you can now write Go code that fails gracefully. This level starts a new phase of the course: **applying Go to real, larger programs**.

Level 2 taught you the basics of packages, imports, and `go.mod` - one directory is one package, capitalization controls visibility, and `go mod init` creates a module. That was enough to organize a single project. This level goes deeper: how modules *really* work, how to build software out of **multiple, independently-versioned modules**, how to control what a module exposes to the outside world with `internal/`, how to design clean package APIs, and how to develop several local modules together with a **workspace** (`go.work`).

**A note on staying offline:** Real Go projects usually depend on modules published to a registry (`proxy.golang.org`) and fetched with `go get`. This course assumes your practice environment may have **no internet access**, so every exercise here is built entirely from **local, sibling modules on disk** - no `go get` against the network, ever. You'll use a `replace` directive (and later `go.work`) to point one local module at another by relative path. This teaches the exact same mechanics real multi-module projects use - it's how you'd develop two modules you own together before either is published - just without needing a network connection.

---

## Table of Contents

1. [Recap: One Package, One Module](#recap-one-package-one-module)
2. [go.mod Anatomy Deep Dive](#gomod-anatomy-deep-dive)
3. [Semantic Versioning As Go Modules Interpret It](#semantic-versioning-as-go-modules-interpret-it)
4. [Practical Module Commands](#practical-module-commands)
5. [Building a Local Multi-Module Setup](#building-a-local-multi-module-setup)
6. [Internal Packages](#internal-packages)
7. [Package-Level API Design](#package-level-api-design)
8. [Multi-Module Workspaces With go.work](#multi-module-workspaces-with-gowork)
9. [Best Practices](#best-practices)
10. [Common Mistakes](#common-mistakes)

---

## Recap: One Package, One Module

Quick bridge from Level 2, so the vocabulary is nailed down before going deeper:

- A **package** is a directory of `.go` files that share a `package` declaration. One directory = one package.
- A **module** is one or more packages that are versioned and released together, described by a single `go.mod` file at the root of the module.
- `go mod init <module-path>` creates that `go.mod` and turns the current directory into the root of a module.

```
mymodule/                     <- module root (has go.mod)
├── go.mod                    <- module example.com/mymodule
├── main.go                   <- package main
├── greetings/                <- package greetings (a package inside the module)
│   └── greetings.go
└── validate/                 <- package validate (another package inside the module)
    └── validate.go
```

**The key relationship:** a module is the *unit of versioning and distribution*; a package is the *unit of code organization inside it*. One module very often contains many packages (like the tree above), but Level 2 only ever showed you a single module. This level is about what happens once you have **more than one module** - and about being deliberate with what each package exposes.

---

## go.mod Anatomy Deep Dive

Every module has exactly one `go.mod` at its root. Here is a `go.mod` with every common directive:

```
module example.com/widgets

go 1.26

require (
    example.com/greetings v1.2.0
    example.com/legacy v0.9.0
)

replace example.com/greetings => ../greetings

exclude example.com/legacy v0.8.0
```

### Line by line

**`module example.com/widgets`**
Declares the module's path - both its unique identity and the import-path *prefix* every package inside it is imported under. If this module has a package in `internal/store`, other modules import it as `example.com/widgets/internal/store` (visibility rules aside - more on that in Section 6).

**`go 1.26`**
The minimum Go language version the module's code requires. This is not just a comment - the compiler enforces it: code can use language features introduced up to this version, and the `go` toolchain uses it to decide default behavior (like automatically enabling newer module-graph pruning rules). Running `go mod init` stamps in whatever Go version is installed on your machine.

**`require`**
Lists the *direct and indirect* dependencies this module needs, each pinned to an exact version. A single dependency can appear as one line (`require example.com/greetings v1.2.0`) or as a `require (...)` block for several. Go computes the final "build list" (which version of every dependency actually gets used) using **Minimum Version Selection (MVS)**: for each module, it picks the *highest* version requested by anything in the dependency graph - never a version nobody asked for, and never silently "latest".

**`replace`**
Says "wherever the build list says to use module X at version V, use *this* instead." The replacement can be:
- a different version of the same module (`replace example.com/greetings v1.2.0 => example.com/greetings v1.2.1`)
- a fork under a different path (`replace example.com/greetings => example.com/myfork v1.2.0`)
- a **local filesystem path** (`replace example.com/greetings => ../greetings`) - this is the one this entire level leans on, because a local path never touches the network.

`replace` only affects the **main module** - if your module is imported by someone else, your `replace` directives are ignored in their build. It's a tool for local development and testing, not for permanently redirecting a dependency in a published module.

**`exclude`**
Says "never select this exact version of this module, even if MVS would otherwise pick it" - typically used to blacklist one version known to be broken, forcing MVS to fall back to a different version some other requirement provides. Like `replace`, `exclude` only applies in the main module.

### What `go mod init` writes

Running `go mod init example.com/widgets` in an empty directory writes just the first two lines:

```
module example.com/widgets

go 1.26
```

`require`, `replace`, and `exclude` are added later - `require` usually by `go mod tidy` (Section 4) after you write imports, `replace`/`exclude` by hand or with `go mod edit`.

---

## Semantic Versioning As Go Modules Interpret It

Go modules are versioned with **semantic versioning**: `vMAJOR.MINOR.PATCH`, e.g. `v1.4.2`.

| Version | Meaning | Compatibility promise |
|---------|---------|------------------------|
| `v0.x.y` | Pre-release / unstable | **None.** Anything can change between minor or even patch releases. |
| `v1.x.y` | Stable, first major release | Backward compatible within `v1` - `v1.4.0` must not break code written against `v1.0.0`. |
| `vN.x.y` (N ≥ 2) | Stable, breaking release | Backward compatible within `vN`, but assumed **incompatible** with `vN-1`. |

### The rule that trips everyone up: v2+ needs a suffix in the import path

For `v0` and `v1`, the import path is just the module path: `example.com/widgets`. The moment a module reaches `v2.0.0` or higher, Go requires the **major version to appear in both the module path and every import path**:

```
module example.com/widgets/v2       // go.mod for the v2 release

go 1.26
```

```go
import "example.com/widgets/v2/render"   // callers must use the /v2 suffix
```

Why? Because `v1` and `v2+` are declared incompatible, Go treats `example.com/widgets` and `example.com/widgets/v2` as **entirely different modules** that can be imported *side by side* in the same program. That is what lets a large program depend on both an old `v1` dependency (transitively, through some other library) and a `v2` version it uses directly, without a version conflict - something most other languages' package managers cannot do cleanly.

We can't fetch a real published `v2` module without network access, so there's no runnable exercise for this - treat it as something to recognize when you see it in the wild: an import path ending in `/v2`, `/v3`, etc. is not a typo, it's Go's major-version encoding at work. Modules below `v2` never carry a version suffix.

### Pseudo-versions

You'll see versions that don't look hand-written, like `v0.0.0-00010101000000-000000000000`. This is a **pseudo-version** - Go's encoding for "a commit that has no tagged release", built from a timestamp and a commit hash. You'll see Go generate exactly this form later in this level when you add a `require` line for a local module that has no tags.

---

## Practical Module Commands

These commands were run for real, offline, against the two-module `greetings`/`app` setup you'll build in the next section. `GOPROXY=off` was set so nothing could reach the network even by accident - every command below succeeded without it.

### `go mod init`

Creates a new `go.mod`.

```
$ go mod init example.com/greetings
go: creating new go.mod: module example.com/greetings
```

### `go mod tidy`

Reconciles `go.mod`/`go.sum` with what your code actually imports: adds missing `require` entries, removes ones nothing imports anymore. Run after adding or removing an import.

```
$ go mod tidy -v
```

(No output means nothing needed to change - `go.mod` already matched the imports exactly.)

### `go list -m all`

Prints the full build list: every module used, and its resolved version. With a `replace` directive active, it shows the replacement too:

```
$ go list -m all
example.com/app
example.com/greetings v0.0.0-00010101000000-000000000000 => ../greetings
```

Read this as: "the build list resolved `example.com/greetings` to that pseudo-version, but the `replace` directive redirects it to `../greetings` on disk."

### `go mod why`

Explains *why* a module is in the build list - what import chain pulled it in.

```
$ go mod why example.com/greetings
# example.com/greetings
example.com/app
example.com/greetings
```

This reads top-to-bottom as a path: `example.com/app` imports `example.com/greetings`. For a deeper dependency chain, every hop would be listed. `go mod why -m <module>` asks the same question at the module level rather than the package level.

### `go mod graph`

Prints the raw dependency graph as edges: `<requirer> <required>@<version>`.

```
$ go mod graph
example.com/app example.com/greetings@v0.0.0-00010101000000-000000000000
example.com/app go@1.26.5
example.com/greetings@v0.0.0-00010101000000-000000000000 go@1.26.5
go@1.26.5 toolchain@go1.26.5
```

Every module's requirements appear as one line per edge - this is the raw material MVS uses to compute the final build list you see in `go list -m all`.

---

## Building a Local Multi-Module Setup

Here is the complete, verified recipe for two modules - a `greetings` library and an `app` that consumes it - wired together with `replace`, entirely offline. (Exercise 2 walks through this step by step with copy-paste commands; this section explains what's happening.)

### Step 1: The library module

```
$ mkdir greetings && cd greetings
$ go mod init example.com/greetings
```

```go
// greetings.go
package greetings

import "fmt"

func Hello(name string) string {
    if name == "" {
        name = "World"
    }
    return fmt.Sprintf("Hello, %s!", name)
}
```

### Step 2: The app module

```
$ mkdir app && cd app
$ go mod init example.com/app
```

```go
// main.go
package main

import (
    "fmt"

    "example.com/greetings"
)

func main() {
    fmt.Println(greetings.Hello("Gopher"))
}
```

### Step 3: Point app at greetings without the network

`app` imports a module that has never been published anywhere. Two edits to `app/go.mod` make that work:

```
$ go mod edit -require=example.com/greetings@v0.0.0-00010101000000-000000000000
$ go mod edit -replace=example.com/greetings=../greetings
```

`go mod edit` makes scripted, precise changes to `go.mod` without opening an editor - useful in scripts and exactly what's used here. The version string is a placeholder pseudo-version; since `replace` immediately redirects to the local path, its exact value barely matters (Go still requires *some* syntactically valid version for the `require` line to exist).

Resulting `app/go.mod`:

```
module example.com/app

go 1.26.5

require example.com/greetings v0.0.0-00010101000000-000000000000

replace example.com/greetings => ../greetings
```

### Step 4: Run it

```
$ go run .
Hello, Gopher!
```

No network call was made anywhere in this sequence - `go build`/`go run` resolved `example.com/greetings` straight to the `../greetings` directory on disk because of the `replace` line.

This is the general pattern for developing two modules you own together, offline, before either is tagged or published: give the consumer a `require` line (any valid version string) and a `replace` line pointing at the producer's path.

---

## Internal Packages

Go has one visibility mechanism inside a single package (capitalization - Level 2), and one visibility mechanism *between* packages in a module tree: the special directory name **`internal`**.

**The rule:** a package under a path containing an `internal/` segment can only be imported by code rooted at the directory that is the *parent of* `internal`. Anything outside that subtree - even another module entirely - gets a compile error if it tries.

```
example.com/owner/
├── go.mod
├── cmd/
│   └── tool/
│       └── main.go          <- CAN import example.com/owner/internal/secret
└── internal/
    └── secret/
        └── secret.go        <- package secret
```

`cmd/tool` is inside the `example.com/owner` tree (the parent of `internal`), so it's allowed. Verified for real:

```
$ go run ./cmd/tool
top secret
```

Now put a *second, separate module* next to it and have it try the same import:

```go
// example.com/outsider/main.go
package main

import (
    "fmt"

    "example.com/owner/internal/secret"
)

func main() {
    fmt.Println(secret.Value())
}
```

Even after adding a `require`/`replace` pair so the outsider module can physically find the `owner` module on disk, the build fails - captured verbatim:

```
$ go build ./...
package example.com/outsider
        main.go:6:5: use of internal package example.com/owner/internal/secret not allowed
```

This is enforced by the compiler itself, not a linter or convention - there is no way around it short of moving the package out of `internal/`. That's the whole point: `internal/` is how a module's author says "this is an implementation detail of *this* module tree; nothing outside it may depend on it," which means it can be changed or removed in any future release without it counting as a breaking change for anyone else.

**Where to put it:** a top-level `internal/` holds things private to the whole module; `cmd/tool/internal/` would hold things private to just that one command. The rule is always "visible only to the subtree rooted at `internal`'s parent" - it nests correctly at any depth.

---

## Package-Level API Design

Level 2 introduced exported vs. unexported names. Now that packages are the unit you ship and depend on across module boundaries, the *design* of a package's exported surface matters as much as any individual function's correctness.

### Cohesion: one package, one job

A package should have a single, nameable reason to exist. `textutil` (string helpers) and `mathutil` (numeric helpers) are cohesive; a `toolbox` package that mixes both is not - nothing unifies its contents except "stuff I needed once." Exercise 7 walks through refactoring exactly this: a monolithic `toolbox` package split into `textutil` and `mathutil`, both verified to produce identical output before and after.

### Minimize the exported surface

Every exported name is a promise: "other code may now depend on this, forever (or until the next major version)." Before exporting something, ask whether callers actually need it, or whether it's an implementation detail that happens to be capitalized out of habit.

```go
// ✅ Small, deliberate surface
package cache

func New(capacity int) *Cache { ... }   // exported: the public entry point
func (c *Cache) Get(key string) (string, bool) { ... }
func (c *Cache) Set(key, value string) { ... }

func evict(c *Cache) { ... }            // unexported: an implementation detail
type bucket struct{ ... }               // unexported: internal storage shape
```

```go
// ❌ Everything exported "just in case"
package cache

func New(capacity int) *Cache { ... }
func Evict(c *Cache) { ... }            // now every future refactor of eviction is a breaking change
type Bucket struct{ ... }               // internal storage shape is now part of the public API
```

A smaller exported surface is easier to document, easier to keep backward compatible, and easier for callers to learn - they only ever see the handful of names that matter.

### Avoiding circular imports (recap + a new example)

Level 2 showed the fix (extract a common package, use interfaces, or invert the dependency). The underlying rule stays the same at the module level too: Go's import graph - across packages *and* across modules - must be a **DAG** (no cycles), full stop. It's worth seeing once more with fresh names, because the fix generalizes:

```go
// order/order.go
package order

import "example.com/shop/shipping"   // order depends on shipping

type Order struct {
    ID   int
    Cost float64
}

func (o Order) EstimateDelivery() string {
    return shipping.Estimate(o.Cost)
}
```

```go
// shipping/shipping.go — DO NOT import "example.com/shop/order" here!
package shipping

// ❌ import "example.com/shop/order" would create order -> shipping -> order

func Estimate(cost float64) string {
    if cost > 100 {
        return "free shipping"
    }
    return "standard shipping"
}
```

If `shipping` ever needed order data, the fix is the same three options from Level 2: pass in only the plain values it needs (as above - `Estimate` takes a `float64`, not an `Order`), extract a shared `shoptypes` package with no dependencies of its own, or define an interface in `shipping` that `order` satisfies. Reach for the smallest fix that keeps the dependency pointing one direction.

### Small, focused packages vs. one giant package

A single sprawling `utils` or `common` package tends to grow forever, accumulate unrelated code, and become a package everything else depends on (so nothing can be changed safely). Prefer several small packages named for what they *do* (`textutil`, `mathutil`, `validate`, `cache`) over one named for what it *is not* (`utils`, `misc`, `common`). If you can't describe a package's purpose in one short sentence, it's a candidate for splitting.

---

## Multi-Module Workspaces With go.work

`replace` directives work, but they're a bit of a hack: they live in `go.mod`, which normally describes *dependencies*, not *"also, use my local copy while I'm developing this."* If you're actively developing two or more modules together, a **workspace** is the cleaner tool.

A workspace is a `go.work` file that lists several module directories and tells the `go` command "resolve imports between these by reading them straight off disk - skip the module cache and the network entirely for anything they need from each other."

### Setting one up

Given sibling directories `greetings/` and `app/` (each already a module, `app` importing `example.com/greetings` with **no** `require` or `replace` needed for it - the workspace supplies that):

```
$ go work init
$ go work use ./greetings ./app
```

The resulting `go.work`, captured for real:

```
go 1.26.5

use (
        ./app
        ./greetings
)
```

### What it buys you

Before creating `go.work`, `app`'s go.mod has no `require` line for `greetings` at all, and running it fails - real captured error:

```
$ go run .
main.go:6:5: missing go.sum entry for module providing package example.com/greetings (imported by example.com/app); to add:
        go get example.com/app
```

After `go work init` + `go work use`, from the workspace root:

```
$ go run ./app
Hello, Workspace!
```

It just works - **no `require`, no `replace`, no `go.sum` entry needed** for `example.com/greetings` in `app/go.mod` at all. `go list -m all` (run from the workspace root) confirms both modules resolve cleanly with no versions to fetch:

```
$ go list -m all
example.com/app
example.com/greetings
```

### A gotcha worth knowing

`go build ./...` and `go vet ./...` run from the **workspace root itself** (not inside a member module) fail, because the root isn't a module:

```
$ go build ./...
pattern ./...: directory prefix . does not contain modules listed in go.work or their selected dependencies
```

Run these from inside a member module directory (`cd app && go build ./...`), or target a member explicitly from the root (`go run ./app`, `go build ./app/...`) - both were verified to work above.

### replace vs. go.work

| | `replace` in go.mod | `go.work` |
|---|---|---|
| Lives in | The consuming module's `go.mod` (committed to version control) | A separate `go.work` file (usually **not** committed - it's a local dev setting) |
| Scope | Only affects the main module's own build | Affects every module listed with `use`, for every command run from the workspace |
| Typical use | "Temporarily point this one dependency at a fork or local path" | "I'm actively developing 2+ modules together and want them to see each other's changes instantly" |
| Needs a `require` line? | Yes, still needs a version placeholder | No - workspace resolution replaces the need for one between workspace members |

Both are legitimate, offline-friendly tools - `replace` is the right call for a one-off local override you might commit; `go.work` is the right call for ongoing multi-module development, and it's designed to be added to your local `.gitignore` rather than shared with collaborators (each developer can point their own workspace at whatever local checkouts they have).

---

## Best Practices

### 1. Keep `replace` and `go.work` for development, not release

A published module should never rely on `replace` to work correctly for its users - `replace` in your `go.mod` is ignored by anyone importing your module, so if your code doesn't work without it, it doesn't really work. Use `replace`/`go.work` while developing, remove them (or keep `go.work` untracked) before tagging a release.

### 2. Run `go mod tidy` before every commit

It keeps `go.mod`/`go.sum` in sync with what your code actually imports - no stray requirements, nothing missing. Treat an unclean `go mod tidy` diff as a signal something changed in your imports that you should look at.

### 3. Default to the smallest possible exported API

Start with everything unexported, and export only what callers actually need, when they need it. It's easy to export something later; it's a breaking change to unexport it.

### 4. Name packages for what they do, not what they contain

`validate`, `render`, `cache` describe a job. `utils`, `common`, `helpers` describe nothing - they're where code goes when nobody thought about where it belongs, and they tend to grow without bound.

### 5. Put a module's version-sensitive internals behind `internal/`

If a piece of code should never be depended on outside your module tree, put it under `internal/` immediately - don't wait until someone else has already imported it and you have to make it a breaking change to move it there.

### 6. Prefer `go.work` over hand-edited `replace` for multi-module development

If you're actively iterating on two or more local modules together, `go.work` avoids having to add and remove `replace` lines from `go.mod` (and forgetting to remove one before committing).

---

## Common Mistakes

### Mistake 1: Committing a `replace` That Points at a Local-Only Path

```
// ❌ WRONG - breaks the build for anyone who isn't you
replace example.com/greetings => /Users/yourname/dev/greetings
```

Anyone else building this module has no `/Users/yourname/dev/greetings` on their machine. Local-path replaces are for your own development loop - remove them (or use a relative path *within a shared repo layout*, or better, `go.work`) before sharing the module.

### Mistake 2: Forgetting the `/v2`+ Suffix When Bumping a Major Version

```
// ❌ WRONG - go.mod still says v1 path, but this is meant to be a breaking v2 release
module example.com/widgets

// ✅ RIGHT
module example.com/widgets/v2
```

Without the suffix, Go has no way to let old `v1` importers and new `v2` importers coexist in someone else's dependency graph - the version suffix *is* the mechanism that makes major versions distinct modules.

### Mistake 3: Expecting `go.work` to Be Shared Like `go.mod`

```
// ❌ WRONG ASSUMPTION - committing go.work and expecting it to travel with the repo
```

`go.work` almost always points at paths specific to *your* machine's checkout layout. Treat it like a personal, local dev setting (commonly `.gitignore`d), not a piece of the module's public definition the way `go.mod` is.

### Mistake 4: Trying to Import an `internal/` Package From Outside Its Tree

```
// ❌ WRONG - compile error, not a lint warning
import "example.com/owner/internal/secret"   // from a different module
```

There is no workaround short of the package's owner moving it out of `internal/`. If you find yourself needing something under someone else's `internal/`, that's a signal to ask for it to be exported, not to route around the restriction.

### Mistake 5: One Catch-All Package Instead of Several Focused Ones

```
// ❌ WRONG - "utils" absorbs everything, forever
package utils

func FormatDate(...) string { ... }
func Max(a, b int) int { ... }
func ValidateEmail(s string) bool { ... }
```

```
// ✅ RIGHT - split by responsibility
package dateutil   // FormatDate
package mathutil   // Max
package validate   // ValidateEmail
```

A catch-all package becomes a dependency of nearly everything else, which makes it risky to change and hard to reason about. Small, purpose-named packages keep dependencies legible.

### Mistake 6: Assuming `go mod tidy` Uses the Network Even With a Working `replace`/`go.work`

```
// ❌ WRONG ASSUMPTION - "go mod tidy will need the internet since it's a real module path"
```

Once a `replace` directive or `go.work` entry resolves a module path to a local directory, `go mod tidy`/`go build`/`go run` read it straight from disk - no network call happens for that module, regardless of whether its path looks like a normal registry path (`example.com/...`). This is exactly what makes the offline workflow in this level possible.

---

## Summary

**go.mod anatomy:**
- `module` - this module's identity and import-path prefix
- `go` - minimum language version
- `require` - dependencies and their versions (computed by MVS)
- `replace` - redirect a dependency, for local development (main module only)
- `exclude` - blacklist a specific version (main module only)

**Semantic versioning:**
- `v0.x.y` - no compatibility promise
- `v1.x.y` - stable, backward compatible within v1
- `v2+` - stable, but needs a `/vN` suffix in the module path and every import path

**Practical commands:** `go mod init`, `go mod tidy`, `go list -m all`, `go mod why`, `go mod graph` - all runnable, all offline, once local modules are wired together with `replace` or `go.work`.

**Multi-module development, offline:** a `require` line with any valid version, plus a `replace` line pointing at a relative path, is enough to develop two local modules together with zero network access.

**`internal/`:** a directory segment that makes everything beneath it importable only from the subtree rooted at its parent - enforced by the compiler, not a convention.

**Package API design:** cohesive, small, deliberately-exported packages named for what they do; avoid circular imports and catch-all packages.

**`go.work`:** the cleaner alternative to `replace` for actively developing multiple local modules together - no `require`/`replace`/`go.sum` entries needed between workspace members, at the cost of `go.work` being a local, not-typically-committed file.

---

## Next Steps

You now understand:
- ✅ Every line of a `go.mod` file, including `replace` and `exclude`
- ✅ How Go modules interpret semantic versioning, including the v2+ import-path rule
- ✅ How to run and read `go mod tidy`, `go list -m all`, `go mod why`, and `go mod graph`
- ✅ How to build a real multi-module project entirely offline with `replace`
- ✅ How `internal/` enforces module-tree-only visibility at compile time
- ✅ How to design cohesive, minimal, focused package APIs
- ✅ How to develop multiple local modules together with `go.work`

**Next level:** Level 18 - File Handling
- Reading and writing files
- Working with paths and directories
- Buffered I/O
- Handling file-related errors (building on Level 16)

You're now equipped to structure Go code the way real, multi-module projects do. Keep going! 🚀
