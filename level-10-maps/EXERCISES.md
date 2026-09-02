# Level 10: Maps - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

---

## Exercise 1: Creating Maps

**Objective:** Practice all three ways to create a map, and see the nil-vs-empty distinction

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level10-exercise1
cd ~/projects/level10-exercise1
go mod init level10.example/exercise1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== Map Literal ===")
    ages := map[string]int{
        "Alice": 30,
        "Bob":   25,
    }
    fmt.Println(ages)

    fmt.Println("\n=== make(map[K]V) ===")
    scores := make(map[string]int)
    scores["Charlie"] = 88
    scores["Diana"] = 91
    fmt.Println(scores)

    fmt.Println("\n=== Empty Map Literal vs nil Map ===")
    empty := map[string]int{}
    var nilMap map[string]int
    fmt.Println("empty:", empty, "len:", len(empty), "is nil:", empty == nil)
    fmt.Println("nilMap:", nilMap, "len:", len(nilMap), "is nil:", nilMap == nil)

    fmt.Println("\n=== Reading From a nil Map Is Safe ===")
    value := nilMap["anything"]
    fmt.Println("reading a missing key from a nil map returns the zero value:", value)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Map Literal ===
map[Alice:30 Bob:25]

=== make(map[K]V) ===
map[Charlie:88 Diana:91]

=== Empty Map Literal vs nil Map ===
empty: map[] len: 0 is nil: false
nilMap: map[] len: 0 is nil: true

=== Reading From a nil Map Is Safe ===
reading a missing key from a nil map returns the zero value: 0
```

**Learning Objectives:**
- ✅ Create maps with a literal, with `make`, and via the `var` zero value
- ✅ See that a nil map and an empty map both print `map[]` but are NOT the same thing
- ✅ Confirm that reading from a nil map is safe and returns the zero value

---

## Exercise 2: The nil Map Write Panic

**Objective:** Trigger the real nil-map-write panic and recover from it safely

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level10-exercise2
cd ~/projects/level10-exercise2
go mod init level10.example/exercise2
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func writeToNilMap() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from panic:", r)
        }
    }()

    var m map[string]int // nil map - declared but never made
    fmt.Println("About to write to a nil map...")
    m["key"] = 1 // this line panics
    fmt.Println("This line never runs")
}

func main() {
    fmt.Println("=== Reading From a nil Map (safe) ===")
    var readMap map[string]int
    fmt.Println("value:", readMap["missing"])

    fmt.Println("\n=== Writing To a nil Map (panics, recovered) ===")
    writeToNilMap()

    fmt.Println("\n=== Program Continues After Recovery ===")
    fmt.Println("Fix: initialize with make() or a map literal before writing")
    fixed := make(map[string]int)
    fixed["key"] = 1
    fmt.Println("fixed map after write:", fixed)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Reading From a nil Map (safe) ===
value: 0

=== Writing To a nil Map (panics, recovered) ===
About to write to a nil map...
Recovered from panic: assignment to entry in nil map

=== Program Continues After Recovery ===
Fix: initialize with make() or a map literal before writing
fixed map after write: map[key:1]
```

4. **Bonus verification - see the real, uncaught panic.** Comment out the `defer`/`recover` block (or create a second, throwaway `main.go` with just `var m map[string]int; m["key"] = 1`) and run it with `go run main.go`. You should see something very close to:

```
panic: assignment to entry in nil map

goroutine 1 [running]:
main.main()
        /path/to/main.go:5 +0x34
exit status 2
```

The file path and hex offset will differ on your machine, but the panic message text - `assignment to entry in nil map` - is exact, fixed wording from the Go runtime.

**Learning Objectives:**
- ✅ Confirm reading from a nil map never panics
- ✅ Trigger and observe the exact nil-map-write panic message
- ✅ Use `recover()` to catch a panic and let the program continue
- ✅ Understand the fix: always initialize a map before writing to it

---

## Exercise 3: CRUD Operations

**Objective:** Practice Create, Read, Update, and Delete on a map

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level10-exercise3
cd ~/projects/level10-exercise3
go mod init level10.example/exercise3
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    inventory := make(map[string]int)

    fmt.Println("=== Insert ===")
    inventory["apples"] = 50
    inventory["bananas"] = 30
    inventory["cherries"] = 100
    fmt.Println(inventory)

    fmt.Println("\n=== Read ===")
    fmt.Println("apples:", inventory["apples"])
    fmt.Println("missing item (grapes):", inventory["grapes"]) // zero value, no panic

    fmt.Println("\n=== Update ===")
    inventory["apples"] = inventory["apples"] - 10 // sold 10 apples
    fmt.Println("apples after sale:", inventory["apples"])

    fmt.Println("\n=== Delete ===")
    delete(inventory, "bananas")
    fmt.Println(inventory)

    fmt.Println("\n=== Delete a Key That Doesn't Exist (no-op, no panic) ===")
    delete(inventory, "grapes")
    fmt.Println(inventory)

    fmt.Println("\n=== len() Counts Entries ===")
    fmt.Println("total distinct items:", len(inventory))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Insert ===
map[apples:50 bananas:30 cherries:100]

=== Read ===
apples: 50
missing item (grapes): 0

=== Update ===
apples after sale: 40

=== Delete ===
map[apples:40 cherries:100]

=== Delete a Key That Doesn't Exist (no-op, no panic) ===
map[apples:40 cherries:100]

=== len() Counts Entries ===
total distinct items: 2
```

**Learning Objectives:**
- ✅ Insert and update entries with `m[k] = v`
- ✅ Read entries, including missing ones (zero value, no panic)
- ✅ Delete entries safely, even keys that don't exist
- ✅ Use `len()` to count entries in a map

---

## Exercise 4: The comma-ok Idiom Deep Dive

**Objective:** Distinguish "key absent" from "key present with a zero value"

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level10-exercise4
cd ~/projects/level10-exercise4
go mod init level10.example/exercise4
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    stock := map[string]int{
        "apples":  50,
        "bananas": 0, // present, but quantity is genuinely zero
    }

    fmt.Println("=== Without comma-ok: Can't Tell Absent From Zero ===")
    fmt.Println("apples:", stock["apples"])
    fmt.Println("bananas:", stock["bananas"])
    fmt.Println("grapes:", stock["grapes"]) // absent - also prints 0!

    fmt.Println("\n=== With comma-ok: The Truth Comes Out ===")
    for _, item := range []string{"apples", "bananas", "grapes"} {
        qty, ok := stock[item]
        if ok {
            fmt.Printf("%-8s present, quantity = %d\n", item, qty)
        } else {
            fmt.Printf("%-8s NOT in stock map at all\n", item)
        }
    }

    fmt.Println("\n=== comma-ok Directly in an if Statement ===")
    if qty, ok := stock["bananas"]; ok {
        fmt.Printf("bananas found, quantity = %d (zero, but present!)\n", qty)
    }
    if _, ok := stock["grapes"]; !ok {
        fmt.Println("grapes confirmed absent from the map")
    }

    fmt.Println("\n=== comma-ok Alone (checking existence only) ===")
    _, hasApples := stock["apples"]
    fmt.Println("has apples key:", hasApples)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Without comma-ok: Can't Tell Absent From Zero ===
apples: 50
bananas: 0
grapes: 0

=== With comma-ok: The Truth Comes Out ===
apples   present, quantity = 50
bananas  present, quantity = 0
grapes   NOT in stock map at all

=== comma-ok Directly in an if Statement ===
bananas found, quantity = 0 (zero, but present!)
grapes confirmed absent from the map

=== comma-ok Alone (checking existence only) ===
has apples key: true
```

**Learning Objectives:**
- ✅ See firsthand that a bare index can't distinguish absent from zero
- ✅ Use `v, ok := m[k]` to get the truth about presence
- ✅ Use comma-ok directly inside an `if` statement
- ✅ Use `_, ok := m[k]` to check existence only

---

## Exercise 5: Map Key Requirements

**Objective:** Use valid comparable key types, then trigger the real compile error for an invalid one

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level10-exercise5
cd ~/projects/level10-exercise5
go mod init level10.example/exercise5
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type Point struct {
    X, Y int
}

func main() {
    fmt.Println("=== Comparable Keys Work Fine ===")

    // Basic types as keys
    byInt := map[int]string{1: "one", 2: "two"}
    fmt.Println("int keys:", byInt)

    // Array as a key (arrays ARE comparable when their element type is)
    byCoord := map[[2]int]string{
        {0, 0}: "origin",
        {1, 1}: "diagonal",
    }
    fmt.Println("array keys:", byCoord)

    // Struct as a key (comparable because all its fields are comparable)
    byPoint := map[Point]string{
        {X: 0, Y: 0}: "origin",
        {X: 3, Y: 4}: "some point",
    }
    fmt.Println("struct keys:", byPoint[Point{X: 3, Y: 4}])
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Comparable Keys Work Fine ===
int keys: map[1:one 2:two]
array keys: map[[0 0]:origin [1 1]:diagonal]
struct keys: some point
```

4. **Now trigger the invalid-key compile error for real.** Create a second file to see what happens when a key type isn't comparable:

```bash
mkdir -p ~/projects/level10-exercise5-badkey
cd ~/projects/level10-exercise5-badkey
go mod init level10.example/exercise5badkey

cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    bad := map[[]int]string{} // slices are not comparable
    fmt.Println(bad)
}
EOF

go build ./...
```

**Expected Output (compile error, captured verbatim from `go build`):**

```
# level10.example/exercise5badkey
./main.go:6:13: invalid map key type []int
```

The build fails with exit status 1 - this is caught entirely at compile time, before the program ever runs.

**Learning Objectives:**
- ✅ Use basic types, arrays, and structs as map keys
- ✅ Understand that comparability (not "hashability") is the requirement
- ✅ See the exact compiler error produced by an invalid (slice) key type

---

## Exercise 6: Maps as Sets

**Objective:** Build a set two ways and see why `map[T]struct{}` is the idiomatic choice

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level10-exercise6
cd ~/projects/level10-exercise6
go mod init level10.example/exercise6
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "unsafe"
)

func main() {
    fmt.Println("=== Set Using map[T]bool ===")
    visitedBool := map[string]bool{
        "home":    true,
        "about":   true,
        "pricing": true,
    }
    visitedBool["contact"] = true
    fmt.Println("visited 'home'?", visitedBool["home"])
    fmt.Println("visited 'blog'?", visitedBool["blog"]) // false: absent key, zero value

    fmt.Println("\n=== Set Using map[T]struct{} (the zero-memory idiom) ===")
    visitedSet := map[string]struct{}{
        "home":    {},
        "about":   {},
        "pricing": {},
    }
    visitedSet["contact"] = struct{}{}

    _, seenHome := visitedSet["home"]
    _, seenBlog := visitedSet["blog"]
    fmt.Println("visited 'home'?", seenHome)
    fmt.Println("visited 'blog'?", seenBlog)

    fmt.Println("\n=== Why struct{} Is Zero-Memory ===")
    var s struct{}
    fmt.Println("unsafe.Sizeof(bool(true)) =", unsafe.Sizeof(true), "bytes")
    fmt.Println("unsafe.Sizeof(struct{}{}) =", unsafe.Sizeof(s), "bytes")

    fmt.Println("\n=== Removing From a Set ===")
    delete(visitedSet, "about")
    _, stillThere := visitedSet["about"]
    fmt.Println("'about' still in set after delete?", stillThere)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Set Using map[T]bool ===
visited 'home'? true
visited 'blog'? false

=== Set Using map[T]struct{} (the zero-memory idiom) ===
visited 'home'? true
visited 'blog'? false

=== Why struct{} Is Zero-Memory ===
unsafe.Sizeof(bool(true)) = 1 bytes
unsafe.Sizeof(struct{}{}) = 0 bytes

=== Removing From a Set ===
'about' still in set after delete? false
```

**Learning Objectives:**
- ✅ Implement a set with `map[T]bool`
- ✅ Implement a set with `map[T]struct{}`
- ✅ Prove with `unsafe.Sizeof` that `struct{}` occupies zero bytes
- ✅ Use `delete()` to remove a member from a set

---

## Exercise 7: Maps of Structs vs Pointers to Structs

**Objective:** Trigger the real "not addressable" compile error, then apply both fixes

**Instructions:**

1. First, trigger the compile error for real:

```bash
mkdir -p ~/projects/level10-exercise7-bad
cd ~/projects/level10-exercise7-bad
go mod init level10.example/exercise7bad

cat > main.go << 'EOF'
package main

import "fmt"

type Player struct {
    Name  string
    Score int
}

func main() {
    players := map[string]Player{
        "alice": {Name: "Alice", Score: 0},
    }

    players["alice"].Score = 10 // does not compile
    fmt.Println(players)
}
EOF

go build ./...
```

**Expected Output (compile error, captured verbatim from `go build`):**

```
# level10.example/exercise7bad
./main.go:15:2: cannot assign to struct field players["alice"].Score in map
```

2. Now create the working exercise directory with both fixes:

```bash
mkdir -p ~/projects/level10-exercise7
cd ~/projects/level10-exercise7
go mod init level10.example/exercise7
```

3. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type Player struct {
    Name  string
    Score int
}

func main() {
    fmt.Println("=== The Gotcha (does not compile - see comment) ===")
    fmt.Println("players[\"alice\"].Score = 10")
    fmt.Println("  -> cannot assign to struct field players[\"alice\"].Score in map")
    fmt.Println("  because map values are not addressable")

    fmt.Println("\n=== Fix 1: Read, Modify, Reassign the Whole Struct ===")
    players := map[string]Player{
        "alice": {Name: "Alice", Score: 0},
    }
    p := players["alice"] // copy out
    p.Score = 10          // modify the copy
    players["alice"] = p  // write the whole struct back
    fmt.Println(players["alice"])

    fmt.Println("\n=== Fix 2: Use a Map of Pointers Instead ===")
    playerPtrs := map[string]*Player{
        "bob": {Name: "Bob", Score: 0},
    }
    playerPtrs["bob"].Score = 20 // compiles fine - dereferencing a pointer IS addressable
    fmt.Println(*playerPtrs["bob"])
}
EOF
```

4. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== The Gotcha (does not compile - see comment) ===
players["alice"].Score = 10
  -> cannot assign to struct field players["alice"].Score in map
  because map values are not addressable

=== Fix 1: Read, Modify, Reassign the Whole Struct ===
{Alice 10}

=== Fix 2: Use a Map of Pointers Instead ===
{Bob 20}
```

**Learning Objectives:**
- ✅ Trigger and read the exact "cannot assign to struct field ... in map" compile error
- ✅ Understand why map values are not addressable
- ✅ Fix it by reading, modifying, and reassigning the whole struct
- ✅ Fix it by storing pointers to structs instead of struct values

---

## Exercise 8: Nested Maps

**Objective:** Safely read and write a map of maps

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level10-exercise8
cd ~/projects/level10-exercise8
go mod init level10.example/exercise8
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== Nested Map Literal ===")
    population := map[string]map[string]int{
        "USA": {
            "New York": 8000000,
            "Chicago":  2700000,
        },
        "Japan": {
            "Tokyo": 14000000,
        },
    }
    fmt.Println(population["USA"]["New York"])
    fmt.Println(population["Japan"]["Tokyo"])

    fmt.Println("\n=== Reading a Missing Outer Key Is Safe (returns nil inner map) ===")
    fmt.Println(population["France"] == nil)
    fmt.Println(population["France"]["Paris"]) // reading from nil inner map: 0, no panic

    fmt.Println("\n=== Writing to a Missing Outer Key Panics (same nil-map rule, one level deeper) ===")
    func() {
        defer func() {
            if r := recover(); r != nil {
                fmt.Println("Recovered:", r)
            }
        }()
        population["France"]["Paris"] = 2100000 // panics: inner map is nil
    }()

    fmt.Println("\n=== The Safe Pattern: Initialize the Inner Map First ===")
    if _, exists := population["France"]; !exists {
        population["France"] = make(map[string]int)
    }
    population["France"]["Paris"] = 2100000
    fmt.Println(population["France"])

    fmt.Println("\n=== Building a Nested Map From Scratch ===")
    inventory := make(map[string]map[string]int)
    addStock := func(category, item string, qty int) {
        if inventory[category] == nil {
            inventory[category] = make(map[string]int)
        }
        inventory[category][item] += qty
    }
    addStock("produce", "apples", 50)
    addStock("produce", "apples", 25) // adds to existing
    addStock("produce", "bananas", 30)
    addStock("dairy", "milk", 40)

    fmt.Println("produce apples:", inventory["produce"]["apples"])
    fmt.Println("produce bananas:", inventory["produce"]["bananas"])
    fmt.Println("dairy milk:", inventory["dairy"]["milk"])
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Nested Map Literal ===
8000000
14000000

=== Reading a Missing Outer Key Is Safe (returns nil inner map) ===
true
0

=== Writing to a Missing Outer Key Panics (same nil-map rule, one level deeper) ===
Recovered: assignment to entry in nil map

=== The Safe Pattern: Initialize the Inner Map First ===
map[Paris:2100000]

=== Building a Nested Map From Scratch ===
produce apples: 75
produce bananas: 30
dairy milk: 40
```

**Learning Objectives:**
- ✅ Read through nested maps safely, even past a missing outer key
- ✅ See that writing through a missing outer key panics, same as any nil-map write
- ✅ Guard writes by initializing the inner map first
- ✅ Build a nested map incrementally with a small helper function

---

## Exercise 9: Word Frequency and an Inverted Index

**Objective:** Combine maps with strings (Level 9) to build a small search index

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level10-exercise9
cd ~/projects/level10-exercise9
go mod init level10.example/exercise9
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "sort"
    "strings"
)

func main() {
    documents := map[string]string{
        "doc1": "go is fun and go is fast",
        "doc2": "go maps are fast and easy",
        "doc3": "fun with maps in go",
    }

    fmt.Println("=== Word Frequency Per Document ===")
    docIDs := make([]string, 0, len(documents))
    for id := range documents {
        docIDs = append(docIDs, id)
    }
    sort.Strings(docIDs)

    wordCounts := make(map[string]map[string]int) // docID -> word -> count
    for _, id := range docIDs {
        words := strings.Fields(documents[id])
        counts := make(map[string]int)
        for _, w := range words {
            counts[w]++
        }
        wordCounts[id] = counts
        fmt.Printf("%s: %v word(s), \"go\" appears %d time(s)\n", id, len(words), counts["go"])
    }

    fmt.Println("\n=== Inverted Index: word -> documents containing it ===")
    invertedIndex := make(map[string]map[string]struct{}) // word -> set of docIDs
    for _, id := range docIDs {
        for word := range wordCounts[id] {
            if invertedIndex[word] == nil {
                invertedIndex[word] = make(map[string]struct{})
            }
            invertedIndex[word][id] = struct{}{}
        }
    }

    words := make([]string, 0, len(invertedIndex))
    for w := range invertedIndex {
        words = append(words, w)
    }
    sort.Strings(words)

    for _, w := range words {
        docsForWord := make([]string, 0, len(invertedIndex[w]))
        for id := range invertedIndex[w] {
            docsForWord = append(docsForWord, id)
        }
        sort.Strings(docsForWord)
        fmt.Printf("%-6s -> %v\n", w, docsForWord)
    }

    fmt.Println("\n=== Query: Which Documents Contain \"maps\"? ===")
    if docs, ok := invertedIndex["maps"]; ok {
        result := make([]string, 0, len(docs))
        for id := range docs {
            result = append(result, id)
        }
        sort.Strings(result)
        fmt.Println(result)
    }
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Word Frequency Per Document ===
doc1: 7 word(s), "go" appears 2 time(s)
doc2: 6 word(s), "go" appears 1 time(s)
doc3: 5 word(s), "go" appears 1 time(s)

=== Inverted Index: word -> documents containing it ===
and    -> [doc1 doc2]
are    -> [doc2]
easy   -> [doc2]
fast   -> [doc1 doc2]
fun    -> [doc1 doc3]
go     -> [doc1 doc2 doc3]
in     -> [doc3]
is     -> [doc1]
maps   -> [doc2 doc3]
with   -> [doc3]

=== Query: Which Documents Contain "maps"? ===
[doc2 doc3]
```

**Learning Objectives:**
- ✅ Build a nested map (`map[string]map[string]int`) from tokenized strings
- ✅ Build a second nested structure (an inverted index) from the first
- ✅ Use `map[string]struct{}` as a set of document IDs
- ✅ Sort keys for deterministic, repeatable output

---

## Exercise 10: Comprehensive Practice — Inventory Management System

**Objective:** Combine everything from this level into a small, realistic system

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level10-exercise10
cd ~/projects/level10-exercise10
go mod init level10.example/exercise10
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "sort"
)

// Inventory: category -> item -> quantity
type Inventory map[string]map[string]int

func (inv Inventory) AddStock(category, item string, qty int) {
    if inv[category] == nil {
        inv[category] = make(map[string]int)
    }
    inv[category][item] += qty
}

// RemoveStock returns false if the item doesn't exist or there isn't enough stock.
func (inv Inventory) RemoveStock(category, item string, qty int) bool {
    items, ok := inv[category]
    if !ok {
        return false
    }
    current, ok := items[item]
    if !ok || current < qty {
        return false
    }
    remaining := current - qty
    if remaining == 0 {
        delete(items, item) // clean up zero-quantity entries
    } else {
        items[item] = remaining
    }
    return true
}

func (inv Inventory) CategoryTotal(category string) int {
    total := 0
    for _, qty := range inv[category] {
        total += qty
    }
    return total
}

func (inv Inventory) GrandTotal() int {
    total := 0
    for category := range inv {
        total += inv.CategoryTotal(category)
    }
    return total
}

func (inv Inventory) Report() {
    categories := make([]string, 0, len(inv))
    for c := range inv {
        categories = append(categories, c)
    }
    sort.Strings(categories)

    for _, category := range categories {
        fmt.Printf("%s (total: %d)\n", category, inv.CategoryTotal(category))
        items := make([]string, 0, len(inv[category]))
        for item := range inv[category] {
            items = append(items, item)
        }
        sort.Strings(items)
        for _, item := range items {
            fmt.Printf("  %-10s %d\n", item, inv[category][item])
        }
    }
}

func main() {
    inv := make(Inventory)

    fmt.Println("=== Stocking the Warehouse ===")
    inv.AddStock("produce", "apples", 50)
    inv.AddStock("produce", "bananas", 30)
    inv.AddStock("produce", "apples", 20) // restock, adds to existing
    inv.AddStock("dairy", "milk", 40)
    inv.AddStock("dairy", "cheese", 15)
    inv.AddStock("bakery", "bread", 25)
    inv.Report()

    fmt.Println("\n=== Selling Items ===")
    ok := inv.RemoveStock("produce", "apples", 30)
    fmt.Println("sold 30 apples, success:", ok)

    ok = inv.RemoveStock("dairy", "cheese", 15) // sells out entire stock
    fmt.Println("sold all 15 cheese, success:", ok)

    ok = inv.RemoveStock("bakery", "bread", 100) // not enough stock
    fmt.Println("tried to sell 100 bread (only 25 in stock), success:", ok)

    ok = inv.RemoveStock("produce", "grapes", 1) // item never existed
    fmt.Println("tried to sell grapes (never stocked), success:", ok)

    fmt.Println("\n=== Inventory After Sales ===")
    inv.Report()

    fmt.Println("\n=== Checking a Specific Item With comma-ok ===")
    if qty, ok := inv["dairy"]["cheese"]; ok {
        fmt.Println("cheese still in stock:", qty)
    } else {
        fmt.Println("cheese is sold out and no longer tracked")
    }

    fmt.Println("\n=== Grand Total Across All Categories ===")
    fmt.Println("total units in warehouse:", inv.GrandTotal())
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Stocking the Warehouse ===
bakery (total: 25)
  bread      25
dairy (total: 55)
  cheese     15
  milk       40
produce (total: 100)
  apples     70
  bananas    30

=== Selling Items ===
sold 30 apples, success: true
sold all 15 cheese, success: true
tried to sell 100 bread (only 25 in stock), success: false
tried to sell grapes (never stocked), success: false

=== Inventory After Sales ===
bakery (total: 25)
  bread      25
dairy (total: 40)
  milk       40
produce (total: 70)
  apples     40
  bananas    30

=== Checking a Specific Item With comma-ok ===
cheese is sold out and no longer tracked

=== Grand Total Across All Categories ===
total units in warehouse: 135
```

**Learning Objectives:**
- ✅ Combine nested maps, comma-ok, delete, and CRUD into one realistic program
- ✅ Model a two-level lookup (category -> item -> quantity) with methods on a named map type
- ✅ Use comma-ok to reject invalid operations (insufficient stock, unknown items)
- ✅ Clean up zero-quantity entries with `delete` instead of leaving stale zero entries behind

---

## Bonus Challenges

### Challenge 1: Group a Slice of Structs by Field

Given a slice of `Employee{Name, Dept string}` structs, build a `map[string][]Employee` grouping employees by department.

```bash
mkdir -p ~/projects/level10-bonus1
cd ~/projects/level10-bonus1
go mod init level10.example/bonus1
```

**Hints:**
- `groups := make(map[string][]Employee)`
- For each employee: `groups[e.Dept] = append(groups[e.Dept], e)`
- This works even before any entry exists for a department - `append` to a nil slice value (the zero value of an absent map key) just creates a new one

### Challenge 2: A Simple LRU-ish Cache

Implement a small cache with a fixed capacity that evicts the **l**east **r**ecently **u**sed entry when it's full, using a `map[string]int` for storage plus a `[]string` to track access order.

```bash
mkdir -p ~/projects/level10-bonus2
cd ~/projects/level10-bonus2
go mod init level10.example/bonus2
```

**Hints:**
- Keep `data map[string]int` for values and `order []string` with the most-recently-used key at the end
- On `Get` or `Put` for an existing key, remove the key from `order` and re-append it to the end
- On `Put` for a new key when `len(data) >= capacity`, evict `order[0]` (the least recently used) before inserting

### Challenge 3: Set Operations - Union, Intersection, Difference

Given two `map[string]struct{}` sets, implement `union`, `intersection`, and `difference` functions.

```bash
mkdir -p ~/projects/level10-bonus3
cd ~/projects/level10-bonus3
go mod init level10.example/bonus3
```

**Hints:**
- `union(a, b)`: copy every key from `a` into a result, then every key from `b`
- `intersection(a, b)`: keep only keys from `a` that are also present in `b` (comma-ok against `b`)
- `difference(a, b)`: keep only keys from `a` that are NOT present in `b`
- Sort the resulting keys before printing - set iteration order is just as randomized as any other map

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Create maps with a literal, `make`, and understand the nil zero value
✅ Explain (and demonstrate) why reading a nil map is safe but writing panics
✅ Perform full CRUD operations, including the always-safe `delete`
✅ Use the comma-ok idiom to distinguish absent keys from zero-valued ones
✅ Identify which types are valid map keys, and read the compiler's rejection of invalid ones
✅ Build memory-efficient sets with `map[T]struct{}`
✅ Recognize and fix the "map values aren't addressable" struct-field gotcha
✅ Safely read and write nested maps
✅ Combine maps with strings to build a small inverted index
✅ Design a small, realistic system (an inventory tracker) around a nested map

---

## Next Level

Level 11: Functions
- Declaring and calling functions
- Multiple return values
- Variadic parameters
- Closures and first-class functions

Great work! You're mastering Go! 🚀
