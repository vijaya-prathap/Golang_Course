# Level 10: Maps - Complete Guide

## Introduction

Welcome to Level 10! You've mastered arrays (Level 7), slices (Level 8), and strings & runes (Level 9). Now it's time to go deep on **maps** - Go's built-in hash table type for storing key-value pairs.

You've actually already brushed up against maps a few times:
- Level 3 introduced map literals as part of Go's composite types
- Level 5 used the comma-ok idiom (`v, ok := m[k]`) inside an `if` statement
- Level 6 ranged over maps, saw randomized iteration order, sorted keys for stable output, and used a map as a frequency counter

This level builds on all of that. Instead of just *using* maps, you'll learn exactly how they behave under the hood: what makes a value nil-safe to read but panic-prone to write, which types are even allowed as keys, why map values aren't addressable, and how to safely nest maps inside maps.

---

## Table of Contents

1. [What a Map Is](#what-a-map-is)
2. [Creating Maps](#creating-maps)
3. [The nil Map Gotcha](#the-nil-map-gotcha)
4. [CRUD Operations](#crud-operations)
5. [Checking Existence vs Checking for a Zero Value](#checking-existence-vs-checking-for-a-zero-value)
6. [Map Key Requirements](#map-key-requirements)
7. [Maps as Sets](#maps-as-sets)
8. [Maps of Structs vs Maps of Pointers to Structs](#maps-of-structs-vs-maps-of-pointers-to-structs)
9. [Nested Maps](#nested-maps)
10. [Best Practices](#best-practices)
11. [Common Mistakes](#common-mistakes)

---

## What a Map Is

A **map** is an unordered collection of key-value pairs. Internally, Go implements it as a **hash table**: each key is hashed to find a bucket, so lookup, insert, and delete are all fast - close to O(1) on average - regardless of how many entries the map holds.

```go
map[KeyType]ValueType
```

- `KeyType` must be a **comparable** type (more on this in [Map Key Requirements](#map-key-requirements))
- `ValueType` can be absolutely anything - a basic type, a struct, a slice, even another map
- The pairs have **no defined order** - Level 6 already showed you that ranging over a map gives you keys in randomized order

```go
ages := map[string]int{
    "Alice": 30,
    "Bob":   25,
}
fmt.Println(ages) // map[Alice:30 Bob:25]
```

Think of a map as a lookup table: given a key, it hands you the associated value almost instantly, no matter where that key lives internally.

---

## Creating Maps

There are three ways to bring a map into existence, and they behave very differently.

### 1. Map Literal

```go
ages := map[string]int{
    "Alice": 30,
    "Bob":   25,
}
```

Ready to use immediately - both read and write work.

### 2. make(map[K]V)

```go
scores := make(map[string]int)
scores["Charlie"] = 88
```

`make` allocates the underlying hash table and returns a usable, empty map. This is the standard way to create a map you'll populate in code (as opposed to one you can write out as a literal).

### 3. The Zero Value: nil

```go
var scores map[string]int
fmt.Println(scores)          // map[]
fmt.Println(scores == nil)   // true
fmt.Println(len(scores))     // 0
```

Just declaring a map variable with `var` (and no literal) gives you its **zero value**, which is `nil`. A nil map is NOT the same thing as an empty map created by `map[string]int{}` or `make(map[string]int)` - although they print and `len()` the same way, only the nil one reports `true` for `== nil`.

Verified:

```
=== Empty Map Literal vs nil Map ===
empty: map[] len: 0 is nil: false
nilMap: map[] len: 0 is nil: true
```

This distinction matters a lot, as the next section shows.

---

## The nil Map Gotcha

Here's the single most important thing to internalize about maps in Go:

> **Reading** from a nil map is completely safe and returns the zero value.
> **Writing** to a nil map panics.

### Reading Is Safe

```go
var m map[string]int // nil map
value := m["anything"]
fmt.Println(value) // 0 - no panic!
```

A nil map behaves like an empty map for every read-only operation: indexing, `len()`, `range` (zero iterations), and the comma-ok idiom all just work and report "not found."

### Writing Panics

```go
var m map[string]int // nil map
m["key"] = 1         // panic!
```

Running this for real produces exactly this output (captured verbatim from `go run`):

```
panic: assignment to entry in nil map

goroutine 1 [running]:
main.main()
        /path/to/main.go:5 +0x34
exit status 2
```

The exact file path and offset will differ on your machine, but the panic message itself - `assignment to entry in nil map` - is fixed wording from the Go runtime.

### Catching It (for demonstration - not a substitute for initializing your map!)

```go
func writeToNilMap() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from panic:", r)
        }
    }()
    var m map[string]int
    m["key"] = 1
}
```

Verified output:

```
Recovered from panic: assignment to entry in nil map
```

### The Fix

Always initialize a map with `make()` or a literal before writing to it:

```go
m := make(map[string]int) // or m := map[string]int{}
m["key"] = 1              // works fine
```

### Why the Asymmetry?

A nil map has no underlying hash table allocated at all. There's nothing to search, so a read can safely report "not found" (the zero value). But there's also nowhere to *put* a new entry, and Go doesn't silently allocate one for you on write - it panics instead, forcing you to be explicit about initialization. This is a deliberate design choice, not an oversight.

---

## CRUD Operations

Maps support the four operations you'd expect: **C**reate, **R**ead, **U**pdate, **D**elete.

```go
inventory := make(map[string]int)

// Create / Insert
inventory["apples"] = 50

// Read
fmt.Println(inventory["apples"])  // 50
fmt.Println(inventory["grapes"])  // 0 - missing key returns zero value, no panic

// Update (same syntax as insert - a map has no separate "update" operator)
inventory["apples"] = inventory["apples"] - 10

// Delete
delete(inventory, "apples")
```

A few details worth calling out:

- **Insert and update use the exact same syntax** (`m[k] = v`). Go doesn't distinguish "this key is new" from "this key already exists" at the syntax level - `m[k] = v` does whichever is needed.
- **`delete(m, k)` is always safe**, even if `k` isn't in the map. It's a no-op in that case - no panic, no error.
- **`len(m)` returns the number of key-value pairs**, not the number of "slots" or any notion of capacity - maps don't have a `cap()`.

Verified:

```
=== Delete a Key That Doesn't Exist (no-op, no panic) ===
map[apples:40 cherries:100]
```

---

## Checking Existence vs Checking for a Zero Value

Level 5 already showed you the syntax for the comma-ok idiom. Here's *why* it exists and exactly what problem it solves.

### The Problem

```go
stock := map[string]int{
    "apples":  50,
    "bananas": 0, // present, but genuinely out of stock
}

fmt.Println(stock["bananas"]) // 0
fmt.Println(stock["grapes"])  // 0 - also 0! But grapes was never added at all.
```

Indexing a map with a single return value can NEVER tell you whether the zero value you got back means "the key is here and its value happens to be zero" or "the key isn't here at all."

### The Solution: comma-ok

```go
qty, ok := stock["bananas"]
// qty = 0, ok = true  -> present, value is genuinely zero

qty, ok = stock["grapes"]
// qty = 0, ok = false -> absent, this 0 is just the zero value
```

`ok` is a `bool` that's `true` if the key was found, `false` otherwise. This works for maps of any value type - `ok` tells you about *presence*, completely independent of what the value is.

### Comma-ok Inside an if Statement

The idiomatic pattern scopes the temporary variables to the `if` block:

```go
if qty, ok := stock["bananas"]; ok {
    fmt.Println("bananas found, quantity =", qty)
}

if _, ok := stock["grapes"]; !ok {
    fmt.Println("grapes confirmed absent")
}
```

### Checking Existence Only

If you don't care about the value at all, discard it with `_`:

```go
_, hasApples := stock["apples"]
```

Verified output from all of the above:

```
=== With comma-ok: The Truth Comes Out ===
apples   present, quantity = 50
bananas  present, quantity = 0
grapes   NOT in stock map at all
```

**Rule of thumb:** any time zero is a value your program can legitimately store, and you need to distinguish "stored zero" from "never stored," reach for comma-ok.

---

## Map Key Requirements

Not every type can be a map key. Go requires key types to be **comparable** - meaning values of that type can be compared with `==` and `!=`.

### Valid Key Types

| Key type | Comparable? | Example |
|---|---|---|
| Basic types (`int`, `string`, `float64`, `bool`, ...) | ✅ Yes | `map[string]int` |
| Arrays (fixed-size, comparable element type) | ✅ Yes | `map[[2]int]string` |
| Structs (all fields comparable) | ✅ Yes | `map[Point]string` |
| Pointers | ✅ Yes (compares addresses) | `map[*User]bool` |
| Interfaces (if the dynamic value is comparable) | ✅ Yes | `map[error]int` |

### Invalid Key Types

| Key type | Comparable? | Why not |
|---|---|---|
| Slices (`[]T`) | ❌ No | Slices can't be compared with `==` |
| Maps (`map[K]V`) | ❌ No | Same reason - maps aren't comparable |
| Functions (`func(...)`) | ❌ No | Functions aren't comparable |

### Valid Keys in Action

```go
type Point struct{ X, Y int }

byInt := map[int]string{1: "one", 2: "two"}

byCoord := map[[2]int]string{
    {0, 0}: "origin",
    {1, 1}: "diagonal",
}

byPoint := map[Point]string{
    {X: 0, Y: 0}: "origin",
    {X: 3, Y: 4}: "some point",
}
```

Verified output:

```
int keys: map[1:one 2:two]
array keys: map[[0 0]:origin [1 1]:diagonal]
struct keys: some point
```

Arrays and structs are comparable (and therefore valid keys) *as long as every element/field they contain is itself comparable*. A struct containing a slice field, for instance, would NOT be comparable and couldn't be used as a key.

### The Compile Error for an Invalid Key

Trying to use a slice as a key doesn't even compile:

```go
bad := map[[]int]string{} // compile error
```

Running `go build` on this produces exactly this error (captured verbatim):

```
./main.go:6:13: invalid map key type []int
```

This is caught entirely at compile time - there's no way to accidentally ship code with a slice-keyed map.

---

## Maps as Sets

Go has no built-in `Set` type, but a map with a throwaway value makes an excellent one. Two common idioms:

### map[T]bool

```go
visited := map[string]bool{
    "home":  true,
    "about": true,
}
visited["contact"] = true

fmt.Println(visited["home"]) // true
fmt.Println(visited["blog"]) // false - absent key, zero value
```

Simple and readable, but every entry costs at least one byte for the `bool`, and it's easy to accidentally write `visited["x"] = false` instead of deleting the key - which leaves `x` looking "present but not visited," an ambiguous state a true set shouldn't have.

### map[T]struct{} (the idiomatic zero-memory set)

```go
visited := map[string]struct{}{
    "home":  {},
    "about": {},
}
visited["contact"] = struct{}{}

_, seenHome := visited["home"]
_, seenBlog := visited["blog"]
```

`struct{}` is the **empty struct** - a struct type with no fields. It carries no data, and critically, it occupies **zero bytes** of memory:

```go
fmt.Println(unsafe.Sizeof(true))       // 1
fmt.Println(unsafe.Sizeof(struct{}{})) // 0
```

Verified:

```
=== Why struct{} Is Zero-Memory ===
unsafe.Sizeof(bool(true)) = 1 bytes
unsafe.Sizeof(struct{}{}) = 0 bytes
```

With `map[T]struct{}`, presence in the map IS membership - there's no ambiguous "false but present" state, and it's the most memory-efficient set representation in Go. This is why `map[T]struct{}` is considered the idiomatic set type, even though the syntax (`struct{}{}`) looks a little unusual at first.

Both membership test and removal work the same way regardless of which idiom you pick:

```go
_, isMember := mySet[value]
delete(mySet, value)
```

---

## Maps of Structs vs Maps of Pointers to Structs

This is one of the sharpest edges in Go's map design, and it trips up almost everyone the first time.

### The Gotcha

```go
type Player struct {
    Name  string
    Score int
}

players := map[string]Player{
    "alice": {Name: "Alice", Score: 0},
}

players["alice"].Score = 10 // does not compile!
```

Running `go build` on this produces exactly this error (captured verbatim):

```
./main.go:15:2: cannot assign to struct field players["alice"].Score in map
```

### Why This Happens

Map values are **not addressable**. `players["alice"]` isn't a real memory location you can take the address of - it's more like the result of a function call that returns a *copy* of the stored struct. You can read the whole struct, but Go can't let you reach into that copy and mutate a single field, because the mutation would apply to a temporary copy and silently vanish - so instead of allowing confusing, silently-lost writes, Go refuses to compile the expression at all.

Contrast this with a slice: `mySlice[0].Score = 10` compiles fine, because a slice element genuinely IS addressable - it lives at a real, stable location inside the backing array.

### Fix 1: Read, Modify, Reassign the Whole Struct

```go
p := players["alice"] // copy out
p.Score = 10           // modify the copy
players["alice"] = p   // write the whole struct back
```

Verified output: `{Alice 10}`

### Fix 2: Use a Map of Pointers Instead

```go
playerPtrs := map[string]*Player{
    "bob": {Name: "Bob", Score: 0},
}
playerPtrs["bob"].Score = 20 // compiles fine!
```

This works because `playerPtrs["bob"]` yields a **pointer value** (a copy of the pointer, but pointers all point at the same underlying struct). Dereferencing that pointer to reach `.Score` is addressable - you're mutating the one struct the pointer refers to, not a throwaway copy.

Verified output: `{Bob 20}`

**Rule of thumb:** if you need to mutate individual fields of struct values stored in a map, store `*T` (pointers) instead of `T` (values).

---

## Nested Maps

A map's value type can be another map, letting you model two-level (or deeper) lookups: `map[string]map[string]int`.

### Reading Is Safe, Even Through Missing Outer Keys

```go
population := map[string]map[string]int{
    "USA":   {"New York": 8000000, "Chicago": 2700000},
    "Japan": {"Tokyo": 14000000},
}

fmt.Println(population["France"] == nil)   // true  - missing outer key -> nil inner map
fmt.Println(population["France"]["Paris"]) // 0     - reading a nil map is always safe
```

### Writing Through a Missing Outer Key Panics

This is the same nil-map-write rule from earlier, just one level deeper:

```go
population["France"]["Paris"] = 2100000 // panic: assignment to entry in nil map
```

Verified (caught with recover for demonstration):

```
Recovered: assignment to entry in nil map
```

### The Safe Pattern: Initialize the Inner Map First

```go
if _, exists := population["France"]; !exists {
    population["France"] = make(map[string]int)
}
population["France"]["Paris"] = 2100000
```

Or, written as a small reusable helper when building a nested map incrementally:

```go
inventory := make(map[string]map[string]int)

addStock := func(category, item string, qty int) {
    if inventory[category] == nil {
        inventory[category] = make(map[string]int)
    }
    inventory[category][item] += qty
}

addStock("produce", "apples", 50)
addStock("produce", "apples", 25) // adds to the existing 50
```

Verified: `inventory["produce"]["apples"]` is `75` after both calls.

**Rule of thumb:** with nested maps, always check (or unconditionally create) the inner map before writing through it. Reading never needs this guard - only writing does.

---

## Best Practices

### 1. Always Initialize Before Writing

```go
// ✅ Good
scores := make(map[string]int)
scores["Alice"] = 95

// ❌ Panics on the write
var scores map[string]int
scores["Alice"] = 95
```

### 2. Use comma-ok When Zero Is a Legitimate Value

```go
// ✅ Good - distinguishes "absent" from "present but zero"
if qty, ok := stock[item]; ok {
    process(qty)
}

// ❌ Ambiguous - can't tell 0 apart from "not found"
qty := stock[item]
```

### 3. Prefer map[T]struct{} for Sets

```go
// ✅ Good - zero memory per entry, no ambiguous false state
seen := make(map[string]struct{})
seen[id] = struct{}{}

// ⚠️ Works, but wastes memory and allows a confusing false-but-present state
seen := make(map[string]bool)
seen[id] = true
```

### 4. Use Pointers When You Need to Mutate Struct Fields In Place

```go
// ✅ Good
players := map[string]*Player{}
players["alice"].Score++

// ❌ Doesn't compile - map values aren't addressable
players := map[string]Player{}
players["alice"].Score++
```

### 5. Guard Inner Maps Before Writing to Nested Maps

```go
// ✅ Good
if grid[row] == nil {
    grid[row] = make(map[int]string)
}
grid[row][col] = "X"

// ❌ Panics if grid[row] hasn't been created yet
grid[row][col] = "X"
```

### 6. Sort Keys When You Need Stable, Repeatable Output

This carries forward directly from Level 6 - map iteration order is randomized on purpose, so anywhere output order matters (tests, reports, logs), collect and sort the keys first.

```go
keys := make([]string, 0, len(m))
for k := range m {
    keys = append(keys, k)
}
sort.Strings(keys)
```

---

## Common Mistakes

### Mistake 1: Writing to a nil Map

```go
// ❌ WRONG - panics: assignment to entry in nil map
var m map[string]int
m["key"] = 1

// ✅ RIGHT
m := make(map[string]int)
m["key"] = 1
```

### Mistake 2: Using a Bare Index to Check Existence

```go
// ❌ WRONG - can't tell "absent" from "present with zero value"
if m["key"] != 0 {
    // ...
}

// ✅ RIGHT
if _, ok := m["key"]; ok {
    // ...
}
```

### Mistake 3: Trying to Use a Slice (or Map, or Func) as a Key

```go
// ❌ WRONG - compile error: invalid map key type []int
bad := map[[]int]string{}

// ✅ RIGHT - use an array instead, if the size is fixed
ok := map[[3]int]string{}
```

### Mistake 4: Assigning Directly to a Struct Field Inside a Map

```go
type Player struct{ Score int }
players := map[string]Player{"alice": {}}

// ❌ WRONG - compile error: cannot assign to struct field ... in map
players["alice"].Score = 10

// ✅ RIGHT - reassign the whole struct, or use a map of pointers
p := players["alice"]
p.Score = 10
players["alice"] = p
```

### Mistake 5: Writing Through an Uninitialized Inner Map

```go
// ❌ WRONG - panics if grid["a"] hasn't been created yet
grid := make(map[string]map[string]int)
grid["a"]["b"] = 1

// ✅ RIGHT
if grid["a"] == nil {
    grid["a"] = make(map[string]int)
}
grid["a"]["b"] = 1
```

### Mistake 6: Assuming Map Iteration Order Is Stable

Carried forward from Level 6, but worth repeating - it bites people at every level:

```go
// ❌ WRONG ASSUMPTION - order varies between runs
for k := range myMap {
    fmt.Println(k)
}

// ✅ RIGHT - sort keys first if order matters
```

### Mistake 7: Assuming len() Tells You Something About Capacity

```go
// ❌ WRONG ASSUMPTION - maps have no cap(), unlike slices
m := make(map[string]int, 100) // 100 is a size *hint*, not a hard capacity
// m can grow past 100 entries just fine; there's no cap(m)
```

---

## Summary

**Creating Maps:**
- Literal: `map[string]int{"a": 1}`
- `make(map[string]int)` for an empty, ready-to-write map
- `var m map[string]int` gives you `nil` - readable, but panics on write

**The Central Asymmetry:**
- Reading a nil map: always safe, returns the zero value
- Writing to a nil map: always panics with `assignment to entry in nil map`

**CRUD:**
- Create/Update: `m[k] = v` (same syntax for both)
- Read: `v := m[k]` or `v, ok := m[k]`
- Delete: `delete(m, k)` - safe even if `k` isn't present

**comma-ok:**
- `v, ok := m[k]` - `ok` tells you presence, independent of the value
- Essential whenever zero is a value your program can legitimately store

**Map Keys:**
- Must be comparable: basic types, arrays, structs (of comparable fields), pointers
- NOT allowed: slices, maps, functions - these fail to compile

**Sets:**
- `map[T]struct{}` is the idiomatic, zero-memory set - membership IS presence

**Struct Values in Maps:**
- Map values aren't addressable - `m[k].Field = x` doesn't compile for `map[K]StructType`
- Fix: reassign the whole struct, or store `*StructType` instead

**Nested Maps:**
- Reading through a missing key is safe (`nil` inner map)
- Writing requires initializing the inner map first

---

## Next Steps

You now understand:
- ✅ What a map is and how it behaves as a hash table
- ✅ The three ways to create a map, and the special danger of the nil zero value
- ✅ Full CRUD operations, including the safe, always-works `delete`
- ✅ Why and when to reach for the comma-ok idiom
- ✅ Which types are valid map keys, and why slices/maps/funcs aren't
- ✅ How to build sets with `map[T]struct{}`
- ✅ Why map values aren't addressable, and how to work around it
- ✅ How to safely read and write nested maps

**Next level:** Level 11 - Functions
- Declaring and calling functions
- Multiple return values
- Variadic parameters
- Closures and first-class functions

You now have a complete toolkit of Go's core data structures - arrays, slices, strings, and maps. Keep going! 🚀
