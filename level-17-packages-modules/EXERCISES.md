# Level 17: Packages & Modules - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

**Everything in this file works fully offline.** No exercise calls `go get` against the network or needs `proxy.golang.org`. Multi-module exercises use a `replace` directive or a `go.work` file to point one local module at another by relative path on disk - that's what makes it possible to practice real multi-module mechanics without an internet connection.

---

## Exercise 1: go.mod Anatomy Walkthrough

**Objective:** See every directive a go.mod can hold, and inspect it as structured data

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level17-exercise1
cd ~/projects/level17-exercise1
go mod init level17.example/exercise1
```

2. Look at the go.mod that was just generated:

```bash
cat go.mod
```

3. Add a `require`, an `exclude`, and a `replace` directive by hand with `go mod edit` (these name a module that doesn't exist anywhere - that's fine, we're only inspecting syntax here, never building this one):

```bash
go mod edit -require=level17.example/sample@v1.2.0
go mod edit -exclude=level17.example/sample@v1.0.0
go mod edit -replace=level17.example/sample=../sample-not-fetched
cat go.mod
```

4. View the same file as structured JSON:

```bash
go mod edit -json
```

5. Clean the illustrative directives back out (this module won't be used again):

```bash
go mod edit -dropreplace=level17.example/sample
go mod edit -dropexclude=level17.example/sample@v1.0.0
go mod edit -droprequire=level17.example/sample
cat go.mod
```

**Expected Output:**

Step 2:
```
module level17.example/exercise1

go 1.26.5
```

Step 3:
```
module level17.example/exercise1

go 1.26.5

require level17.example/sample v1.2.0

exclude level17.example/sample v1.0.0

replace level17.example/sample => ../sample-not-fetched
```

Step 4:
```
{
	"Module": {
		"Path": "level17.example/exercise1"
	},
	"Go": "1.26.5",
	"Require": [
		{
			"Path": "level17.example/sample",
			"Version": "v1.2.0"
		}
	],
	"Exclude": [
		{
			"Path": "level17.example/sample",
			"Version": "v1.0.0"
		}
	],
	"Replace": [
		{
			"Old": {
				"Path": "level17.example/sample"
			},
			"New": {
				"Path": "../sample-not-fetched"
			}
		}
	],
	"Retract": null,
	"Tool": null,
	"Ignore": null
}
```

Step 5:
```
module level17.example/exercise1

go 1.26.5
```

(Your Go version may differ from `1.26.5` - `go mod init` always stamps in whatever version is installed on your machine.)

**Learning Objectives:**
- ✅ Identify every directive a go.mod can contain: `module`, `go`, `require`, `replace`, `exclude`
- ✅ Use `go mod edit` to script changes to go.mod without an editor
- ✅ Read a go.mod as structured JSON with `go mod edit -json`

---

## Exercise 2: A Two-Module Local Project With a replace Directive

**Objective:** Build a library module and an app module that consumes it, wired together entirely offline

**Instructions:**

1. Create the library module:

```bash
mkdir -p ~/projects/level17-exercise2/greetings
cd ~/projects/level17-exercise2/greetings
go mod init level17.example/exercise2/greetings
```

2. Create `greetings.go`:

```bash
cat > greetings.go << 'EOF'
// Package greetings provides simple greeting utilities.
package greetings

import "fmt"

// Hello returns a greeting for the given name.
func Hello(name string) string {
    if name == "" {
        name = "World"
    }
    return fmt.Sprintf("Hello, %s!", name)
}

// Goodbye returns a farewell for the given name.
func Goodbye(name string) string {
    if name == "" {
        name = "World"
    }
    return fmt.Sprintf("Goodbye, %s!", name)
}
EOF
```

3. Create the app module as a sibling directory:

```bash
mkdir -p ~/projects/level17-exercise2/app
cd ~/projects/level17-exercise2/app
go mod init level17.example/exercise2/app
```

4. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"

    "level17.example/exercise2/greetings"
)

func main() {
    fmt.Println(greetings.Hello("Gopher"))
    fmt.Println(greetings.Goodbye("Gopher"))
}
EOF
```

5. Wire `app` to the local `greetings` module - no network needed, `greetings` has never been published anywhere:

```bash
go mod edit -require=level17.example/exercise2/greetings@v0.0.0-00010101000000-000000000000
go mod edit -replace=level17.example/exercise2/greetings=../greetings
cat go.mod
```

6. Run it:

```bash
go run .
```

**Expected Output:**

Step 5 (`cat go.mod`):
```
module level17.example/exercise2/app

go 1.26.5

require level17.example/exercise2/greetings v0.0.0-00010101000000-000000000000

replace level17.example/exercise2/greetings => ../greetings
```

Step 6:
```
Hello, Gopher!
Goodbye, Gopher!
```

**Learning Objectives:**
- ✅ Create two separate modules on disk that depend on each other
- ✅ Use `go mod edit -require` and `-replace` to wire an unpublished local module into another module's build
- ✅ Confirm the whole thing builds and runs with zero network access

---

## Exercise 3: Inspecting the Build — tidy, why, and graph

**Objective:** Run the real module-inspection commands against the Exercise 2 setup and read their output

**Instructions:**

1. From the `app` module you built in Exercise 2:

```bash
cd ~/projects/level17-exercise2/app
```

2. Run `go mod tidy` verbosely - since go.mod already matches the imports exactly, expect no output:

```bash
go mod tidy -v
```

3. List the full build list:

```bash
go list -m all
```

4. Ask why `greetings` is in the build:

```bash
go mod why level17.example/exercise2/greetings
```

5. Print the raw dependency graph:

```bash
go mod graph
```

**Expected Output:**

Step 2: *(no output - nothing needed to change)*

Step 3:
```
level17.example/exercise2/app
level17.example/exercise2/greetings v0.0.0-00010101000000-000000000000 => ../greetings
```

Step 4:
```
# level17.example/exercise2/greetings
level17.example/exercise2/app
level17.example/exercise2/greetings
```

Step 5:
```
level17.example/exercise2/app go@1.26.5
level17.example/exercise2/app level17.example/exercise2/greetings@v0.0.0-00010101000000-000000000000
go@1.26.5 toolchain@go1.26.5
level17.example/exercise2/greetings@v0.0.0-00010101000000-000000000000 go@1.26.5
```

**Learning Objectives:**
- ✅ Run `go mod tidy` and recognize that no output means go.mod is already correct
- ✅ Read `go list -m all` to see the resolved build list, including an active `replace`
- ✅ Read `go mod why` as an import chain, top to bottom
- ✅ Read `go mod graph`'s `<requirer> <required>@<version>` edge format

---

## Exercise 4: The internal/ Directory — Allowed Import

**Objective:** See that code inside a module's own tree can freely import its `internal/` packages

**Instructions:**

1. Create the module and an `internal` package:

```bash
mkdir -p ~/projects/level17-exercise4/owner/internal/secret
mkdir -p ~/projects/level17-exercise4/owner/cmd/tool
cd ~/projects/level17-exercise4/owner
go mod init level17.example/exercise4/owner
```

2. Create `internal/secret/secret.go`:

```bash
cat > internal/secret/secret.go << 'EOF'
// Package secret holds logic private to the level17.example/exercise4/owner module tree.
package secret

// Value is only meant to be used inside level17.example/exercise4/owner.
func Value() string {
    return "top secret"
}
EOF
```

3. Create `cmd/tool/main.go`, which lives inside the same module tree:

```bash
cat > cmd/tool/main.go << 'EOF'
package main

import (
    "fmt"

    "level17.example/exercise4/owner/internal/secret"
)

func main() {
    fmt.Println(secret.Value())
}
EOF
```

4. Run it:

```bash
go run ./cmd/tool
```

**Expected Output:**

```
top secret
```

**Learning Objectives:**
- ✅ Create an `internal/` package
- ✅ Confirm code rooted at the parent of `internal` (here, the whole `owner` module) can import it freely

---

## Exercise 5: The internal/ Directory — Compile Error From Outside

**Objective:** Capture, verbatim, the real compiler error when a different module tries to import another module's `internal/` package

**Instructions:**

1. Create a second, separate module next to Exercise 4's `owner` module:

```bash
mkdir -p ~/projects/level17-exercise5/outsider
cd ~/projects/level17-exercise5/outsider
go mod init level17.example/exercise5/outsider
```

2. Create `main.go`, importing the *other* module's internal package:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"

    "level17.example/exercise4/owner/internal/secret"
)

func main() {
    fmt.Println(secret.Value())
}
EOF
```

3. Wire the two modules together on disk (so the import can at least be *found* - this does not change the visibility rule):

```bash
go mod edit -require=level17.example/exercise4/owner@v0.0.0-00010101000000-000000000000
go mod edit -replace=level17.example/exercise4/owner=../../level17-exercise4/owner
```

4. Try to build it:

```bash
go build ./...
```

**Expected Output:**

Step 4 (captured verbatim - this is a compile error, the command exits non-zero):
```
package level17.example/exercise5/outsider
	main.go:6:5: use of internal package level17.example/exercise4/owner/internal/secret not allowed
```

**Learning Objectives:**
- ✅ Confirm `internal/` visibility is enforced by the compiler, not by convention
- ✅ See that even a `replace` directive that makes the package *findable* on disk does not make it *importable*
- ✅ Recognize this exact error message if you ever see it in real projects

---

## Exercise 6: Avoiding Circular Imports

**Objective:** Cause a real import cycle, capture the compiler's error, then fix it

**Instructions:**

1. Create a module with two packages that will import each other:

```bash
mkdir -p ~/projects/level17-exercise6/order
mkdir -p ~/projects/level17-exercise6/shipping
cd ~/projects/level17-exercise6
go mod init level17.example/exercise6
```

2. Create `order/order.go`, which imports `shipping`:

```bash
cat > order/order.go << 'EOF'
package order

import "level17.example/exercise6/shipping"

type Order struct {
    ID   int
    Cost float64
}

func (o Order) EstimateDelivery() string {
    return shipping.Estimate(o)
}
EOF
```

3. Create `shipping/shipping.go`, which imports `order` right back - this is the cycle:

```bash
cat > shipping/shipping.go << 'EOF'
package shipping

import "level17.example/exercise6/order"

func Estimate(o order.Order) string {
    if o.Cost > 100 {
        return "free shipping"
    }
    return "standard shipping"
}
EOF
```

4. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"

    "level17.example/exercise6/order"
)

func main() {
    o := order.Order{ID: 1, Cost: 50}
    fmt.Println(o.EstimateDelivery())
}
EOF
```

5. Try to build it (expect a real compile error):

```bash
go build ./...
```

6. Fix it: `shipping` doesn't actually need the whole `Order` - it only needs the cost. Change `shipping.go` to take a plain `float64` instead of importing `order` at all:

```bash
cat > shipping/shipping.go << 'EOF'
package shipping

func Estimate(cost float64) string {
    if cost > 100 {
        return "free shipping"
    }
    return "standard shipping"
}
EOF
```

7. Update `order.go` to pass just the cost:

```bash
cat > order/order.go << 'EOF'
package order

import "level17.example/exercise6/shipping"

type Order struct {
    ID   int
    Cost float64
}

func (o Order) EstimateDelivery() string {
    return shipping.Estimate(o.Cost)
}
EOF
```

8. Update `main.go` to try both a cheap and an expensive order:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"

    "level17.example/exercise6/order"
)

func main() {
    o1 := order.Order{ID: 1, Cost: 50}
    fmt.Println(o1.EstimateDelivery())

    o2 := order.Order{ID: 2, Cost: 150}
    fmt.Println(o2.EstimateDelivery())
}
EOF
```

9. Run it:

```bash
go run .
```

**Expected Output:**

Step 5 (captured verbatim):
```
package level17.example/exercise6
	imports level17.example/exercise6/order from main.go
	imports level17.example/exercise6/shipping from order.go
	imports level17.example/exercise6/order from shipping.go: import cycle not allowed
```

Step 9:
```
standard shipping
free shipping
```

**Learning Objectives:**
- ✅ Produce and read a real `import cycle not allowed` error, including its import-chain trace
- ✅ Fix a circular import by having the "downstream" package depend on a plain value instead of the "upstream" package's type
- ✅ Confirm Go's import graph must always be a DAG - across packages, and across modules

---

## Exercise 7: Package API Design Refactor

**Objective:** Split a monolithic, low-cohesion package into two small, focused ones - with identical behavior before and after

**Instructions:**

1. Create the "before" version - a single `toolbox` package mixing unrelated string and math helpers:

```bash
mkdir -p ~/projects/level17-exercise7/toolbox
cd ~/projects/level17-exercise7
go mod init level17.example/exercise7
```

```bash
cat > toolbox/toolbox.go << 'EOF'
// Package toolbox mixes unrelated string and math helpers together.
package toolbox

import "strings"

func Shout(s string) string {
    return strings.ToUpper(s) + "!"
}

func Reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

func Max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func Min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
EOF
```

2. Create `main.go` using it, and run it:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"

    "level17.example/exercise7/toolbox"
)

func main() {
    fmt.Println(toolbox.Shout("go"))
    fmt.Println(toolbox.Reverse("go"))
    fmt.Println(toolbox.Max(3, 7))
    fmt.Println(toolbox.Min(3, 7))
}
EOF
go run .
```

3. Now refactor: split `toolbox` into two cohesive packages, `textutil` (string helpers) and `mathutil` (numeric helpers), and delete `toolbox` entirely:

```bash
mkdir -p textutil mathutil

cat > textutil/textutil.go << 'EOF'
// Package textutil provides string manipulation helpers.
package textutil

import "strings"

// Shout uppercases s and appends an exclamation mark.
func Shout(s string) string {
    return strings.ToUpper(s) + "!"
}

// Reverse returns s with its characters in reverse order.
func Reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}
EOF

cat > mathutil/mathutil.go << 'EOF'
// Package mathutil provides small numeric helpers.
package mathutil

// Max returns the larger of a and b.
func Max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

// Min returns the smaller of a and b.
func Min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
EOF

rm -rf toolbox
```

4. Update `main.go` to use the two new packages, and run it again:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"

    "level17.example/exercise7/mathutil"
    "level17.example/exercise7/textutil"
)

func main() {
    fmt.Println(textutil.Shout("go"))
    fmt.Println(textutil.Reverse("go"))
    fmt.Println(mathutil.Max(3, 7))
    fmt.Println(mathutil.Min(3, 7))
}
EOF
go run .
```

**Expected Output:**

Both step 2 and step 4 print the exact same thing:
```
GO!
og
7
3
```

**Learning Objectives:**
- ✅ Recognize a low-cohesion package (unrelated responsibilities bundled together)
- ✅ Split it into small, focused packages named for what they do
- ✅ Confirm a package-design refactor can be behavior-preserving - same output, better structure

---

## Exercise 8: A Multi-Module Workspace With go.work

**Objective:** Develop two local modules together using `go.work`, with no `require` or `replace` between them

**Instructions:**

1. Create the library module:

```bash
mkdir -p ~/projects/level17-exercise8/greetings
cd ~/projects/level17-exercise8/greetings
go mod init level17.example/exercise8/greetings
```

```bash
cat > greetings.go << 'EOF'
// Package greetings provides simple greeting utilities.
package greetings

import "fmt"

// Hello returns a greeting for the given name.
func Hello(name string) string {
    if name == "" {
        name = "World"
    }
    return fmt.Sprintf("Hello, %s!", name)
}
EOF
```

2. Create the app module - note there is **no** `require` line added for `greetings` at all:

```bash
mkdir -p ~/projects/level17-exercise8/app
cd ~/projects/level17-exercise8/app
go mod init level17.example/exercise8/app
```

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"

    "level17.example/exercise8/greetings"
)

func main() {
    fmt.Println(greetings.Hello("Workspace"))
}
EOF
```

3. Try to run it as-is - expect it to fail, since nothing has told this module where to find `greetings`:

```bash
go run .
```

4. Go up one level and create the workspace:

```bash
cd ~/projects/level17-exercise8
go work init
go work use ./greetings ./app
cat go.work
```

5. Now run the app from the workspace root - no `require`/`replace` was ever added:

```bash
go run ./app
```

6. Confirm the build list resolves cleanly:

```bash
go list -m all
```

7. Try `go build ./...` from the workspace root itself, then from inside the `app` directory:

```bash
go build ./...
cd app && go build ./... && ls
```

**Expected Output:**

Step 3 (captured verbatim):
```
main.go:6:5: no required module provides package level17.example/exercise8/greetings; to add it:
	go get level17.example/exercise8/greetings
```

Step 4 (`cat go.work`):
```
go 1.26.5

use (
	./app
	./greetings
)
```

Step 5:
```
Hello, Workspace!
```

Step 6:
```
level17.example/exercise8/app
level17.example/exercise8/greetings
```

Step 7, from the workspace root (captured verbatim - this is expected, the root is not itself a module):
```
pattern ./...: directory prefix . does not contain modules listed in go.work or their selected dependencies
```

Step 7, from inside `app` (succeeds, produces an `app` binary alongside `go.mod` and `main.go`):
```
go.mod  main.go  app
```

**Learning Objectives:**
- ✅ Use `go work init` and `go work use` to create a workspace over two local modules
- ✅ Confirm cross-module imports resolve with no `require`, `replace`, or `go.sum` entry needed between workspace members
- ✅ Recognize the workspace-root-is-not-a-module gotcha, and know to run per-member commands from inside each module

---

## Exercise 9: replace vs. go.work, Side by Side

**Objective:** Convert Exercise 2's `replace`-based project into a `go.work` workspace and confirm it behaves identically

**Instructions:**

1. Go back to the `app` module from Exercise 2:

```bash
cd ~/projects/level17-exercise2/app
cat go.mod
```

2. Remove the `require` and `replace` lines - `go.work` is about to take over that job:

```bash
go mod edit -dropreplace=level17.example/exercise2/greetings
go mod edit -droprequire=level17.example/exercise2/greetings
cat go.mod
```

3. Go up to the shared parent directory and create a workspace over both modules:

```bash
cd ~/projects/level17-exercise2
go work init
go work use ./greetings ./app
cat go.work
```

4. Run the app exactly as before:

```bash
go run ./app
```

**Expected Output:**

Step 1 (before conversion):
```
module level17.example/exercise2/app

go 1.26.5

require level17.example/exercise2/greetings v0.0.0-00010101000000-000000000000

replace level17.example/exercise2/greetings => ../greetings
```

Step 2 (after dropping both lines):
```
module level17.example/exercise2/app

go 1.26.5
```

Step 3:
```
go 1.26.5

use (
	./app
	./greetings
)
```

Step 4 (identical to Exercise 2's output, with no `require`/`replace` in `app/go.mod` at all):
```
Hello, Gopher!
Goodbye, Gopher!
```

**Learning Objectives:**
- ✅ Convert a `replace`-based two-module setup into a `go.work` workspace
- ✅ Confirm both approaches produce identical runtime behavior
- ✅ Understand `go.work` moves the "how do these local modules find each other" concern out of go.mod entirely

---

## Exercise 10: Comprehensive Practice — A Multi-Module Notes System

**Objective:** Combine a library module, an internal package, and a go.work workspace into one working system

**Instructions:**

1. Create the `notes` library module, with an `internal/store` package it uses to format output:

```bash
mkdir -p ~/projects/level17-exercise10/notes/internal/store
mkdir -p ~/projects/level17-exercise10/notes/cmd/dump
cd ~/projects/level17-exercise10/notes
go mod init level17.example/exercise10/notes
```

```bash
cat > internal/store/store.go << 'EOF'
// Package store formats note data. It is private to the
// level17.example/exercise10/notes module tree.
package store

import "fmt"

// Format renders a title and body as a single display string.
func Format(title, body string) string {
    return fmt.Sprintf("[%s] %s", title, body)
}
EOF
```

2. Create the public `notes` package, which uses `internal/store` internally:

```bash
cat > notes.go << 'EOF'
// Package notes provides a tiny in-memory note type.
package notes

import "level17.example/exercise10/notes/internal/store"

// Note is a single titled note.
type Note struct {
    Title string
    Body  string
}

// New creates a Note.
func New(title, body string) Note {
    return Note{Title: title, Body: body}
}

// String renders the note for display, using the internal
// store package to do the formatting.
func (n Note) String() string {
    return store.Format(n.Title, n.Body)
}
EOF
```

3. Create a small command *inside the same module tree* that reaches `internal/store` directly (allowed, since it's inside the tree):

```bash
cat > cmd/dump/main.go << 'EOF'
package main

import (
    "fmt"

    "level17.example/exercise10/notes/internal/store"
)

func main() {
    fmt.Println(store.Format("debug", "internal package reachable from within the tree"))
}
EOF
go run ./cmd/dump
```

4. Create the `app` module - it will only ever use `notes`' exported API:

```bash
mkdir -p ~/projects/level17-exercise10/app
cd ~/projects/level17-exercise10/app
go mod init level17.example/exercise10/app
```

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"

    "level17.example/exercise10/notes"
)

func main() {
    n := notes.New("Shopping", "milk, eggs, bread")
    fmt.Println(n.String())
}
EOF
```

5. Tie both modules together with a workspace (no `require`/`replace` needed):

```bash
cd ~/projects/level17-exercise10
go work init
go work use ./notes ./app
go run ./app
```

6. Confirm `internal/` visibility is still enforced *even inside the workspace* - try importing `notes/internal/store` from `app` and watch it fail:

```bash
cd app
cat > main.go << 'EOF'
package main

import (
    "fmt"

    "level17.example/exercise10/notes"
    "level17.example/exercise10/notes/internal/store"
)

func main() {
    n := notes.New("Shopping", "milk, eggs, bread")
    fmt.Println(n.String())
    fmt.Println(store.Format("x", "y"))
}
EOF
go build ./...
```

7. Revert `app/main.go` to the version from step 4 so the project is left in a working state.

**Expected Output:**

Step 3:
```
[debug] internal package reachable from within the tree
```

Step 5:
```
[Shopping] milk, eggs, bread
```

Step 6 (captured verbatim - `internal/` is enforced regardless of `go.work`):
```
package level17.example/exercise10/app
	main.go:7:5: use of internal package level17.example/exercise10/notes/internal/store not allowed
```

**Learning Objectives:**
- ✅ Combine a library module's public API with an `internal/` package used only by that module's own code
- ✅ Tie two local modules together with `go.work` end to end
- ✅ Confirm `internal/` visibility rules are enforced by the compiler regardless of `replace` or `go.work` - a workspace makes packages *findable* across modules, never *more visible*

---

## Bonus Challenges

### Challenge 1: A "stringutil" Library Module

Design a small, standalone `stringutil` module (its own `go.mod`, no dependencies) with a deliberately minimal exported API - two or three well-named functions (for example `Slugify`, `Truncate`, `WordCount`) and nothing else exported. Wire a throwaway `app` module to it with a `replace` directive and print a result from each function.

```bash
mkdir -p ~/projects/level17-bonus1/stringutil
cd ~/projects/level17-bonus1/stringutil
go mod init level17.example/bonus1/stringutil
```

**Hints:**
- Keep helper logic (like a rune-classification check) unexported inside the package
- Write a one-line package comment (`// Package stringutil ...`) above the `package` declaration
- Wire it to an `app` module the same way Exercise 2 did

### Challenge 2: Convert a replace Setup Into a go.work Workspace, From Scratch

Without re-reading Exercise 9, build a fresh two-module `replace`-based project (any pair of packages you like), confirm it runs, then convert it to a `go.work` workspace on your own and confirm the output is unchanged.

```bash
mkdir -p ~/projects/level17-bonus2/lib
cd ~/projects/level17-bonus2/lib
go mod init level17.example/bonus2/lib
```

**Hints:**
- `go mod edit -droprequire=<module>` and `go mod edit -dropreplace=<module>` remove exactly what `go.work` is about to replace
- `go work init` then `go work use <dir1> <dir2>` from the shared parent directory
- Compare `go list -m all` before and after - a workspace-resolved module shows no version, just its path

### Challenge 3: A Script That Runs go mod why Against Your Own Setup

Write a small shell script, `why.sh`, that takes a module path as its one argument and runs `go mod why -m <arg>` inside one of your multi-module exercises from this level, printing a clear header first.

```bash
mkdir -p ~/projects/level17-bonus3
cd ~/projects/level17-bonus3
```

**Hints:**
- `#!/bin/sh` at the top, `chmod +x why.sh` to make it runnable
- Use `"$1"` for the argument and check `[ -z "$1" ]` to print a usage message if it's missing
- Run it from inside `~/projects/level17-exercise2/app` against `level17.example/exercise2/greetings`

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Read and write every directive a go.mod can contain
✅ Build a real two-module project wired together with a `replace` directive, entirely offline
✅ Run and interpret `go mod tidy`, `go list -m all`, `go mod why`, and `go mod graph`
✅ Reproduce and read the exact compiler error for an unauthorized `internal/` import
✅ Cause, diagnose, and fix a real circular import
✅ Refactor a low-cohesion package into small, focused ones without changing behavior
✅ Stand up a `go.work` workspace across multiple local modules, with no `require`/`replace` needed between them
✅ Explain the tradeoffs between `replace` and `go.work` for local multi-module development

---

## Next Level

Level 18: File Handling
- Reading and writing files
- Working with paths and directories
- Buffered I/O
- Handling file-related errors

Great work! You're building software the way real, multi-module Go projects are built! 🚀
