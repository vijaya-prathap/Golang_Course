# Level 10: Quick Reference Card

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir myapp && cd myapp
go mod init myapp

# Create main.go
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    scores := map[string]int{"Alice": 95, "Bob": 87}
    scores["Charlie"] = 92

    if v, ok := scores["Alice"]; ok {
        fmt.Println("Alice:", v)
    }

    delete(scores, "Bob")
    fmt.Println(scores)
}
EOF

# Run
go run main.go
```

---

## 📋 Creating Maps

```go
// Literal - ready immediately
m := map[string]int{"a": 1, "b": 2}

// make - empty, ready immediately
m := make(map[string]int)

// var - the DANGEROUS one: nil, readable but panics on write
var m map[string]int
```

---

## ⚡ The nil Map Rule

```go
var m map[string]int

_ = m["x"]           // ✅ safe, returns 0
v, ok := m["x"]       // ✅ safe, ok == false
len(m)                // ✅ safe, returns 0
for range m { }       // ✅ safe, 0 iterations
delete(m, "x")        // ✅ safe, no-op

m["x"] = 1            // ❌ PANIC: assignment to entry in nil map
```

**Fix:** always `m := make(map[K]V)` or use a literal before writing.

---

## 🔧 CRUD Syntax

```go
m[k] = v          // Create or Update (same syntax for both)
v := m[k]         // Read (zero value if missing)
v, ok := m[k]     // Read with presence check
delete(m, k)      // Delete (safe even if k is missing)
len(m)            // Count of entries
```

---

## ✅❌ comma-ok Idiom

```go
v, ok := m[k]
// ok == true  -> key IS in the map (v is the real stored value, even if it's zero)
// ok == false -> key is NOT in the map (v is just the zero value, meaningless)

if v, ok := m[k]; ok {
    fmt.Println("found:", v)
}

_, ok := m[k] // existence check only, discard the value
```

---

## 🔑 Map Key Types

```go
map[int]string{}        // ✅ basic types
map[[2]int]string{}     // ✅ arrays (if element type is comparable)
map[Point]string{}      // ✅ structs (if all fields are comparable)
map[*User]string{}      // ✅ pointers

map[[]int]string{}      // ❌ invalid map key type []int
map[map[string]int]int{} // ❌ maps can't be keys
map[func()]int{}        // ❌ funcs can't be keys
```

---

## 🧱 Maps as Sets

```go
// map[T]bool - simple, 1 byte per entry
set := map[string]bool{"a": true}
if set["a"] { }

// map[T]struct{} - idiomatic, ZERO bytes per entry
set := map[string]struct{}{"a": {}}
if _, ok := set["a"]; ok { }

delete(set, "a") // remove from either kind
```

---

## 🧩 Struct Values in Maps: The Gotcha

```go
type Player struct{ Score int }
players := map[string]Player{"a": {}}

players["a"].Score = 10
// ❌ cannot assign to struct field players["a"].Score in map

// ✅ Fix 1: reassign the whole struct
p := players["a"]
p.Score = 10
players["a"] = p

// ✅ Fix 2: use a map of pointers
ptrs := map[string]*Player{"a": {}}
ptrs["a"].Score = 10 // fine
```

---

## 🪆 Nested Maps

```go
m := make(map[string]map[string]int)

m["a"]["b"] = 1
// ❌ PANIC: assignment to entry in nil map (inner map doesn't exist yet)

// ✅ Guard the inner map first
if m["a"] == nil {
    m["a"] = make(map[string]int)
}
m["a"]["b"] = 1 // fine

m["missing"]["x"] // ✅ reading is always safe: 0, no panic
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Writing to a nil map | `var m map[string]int; m[k]=v` | `m := make(map[string]int)` |
| Bare index for presence | `if m[k] != 0 { }` | `if _, ok := m[k]; ok { }` |
| Slice/map/func as key | `map[[]int]string{}` | `map[[3]int]string{}` |
| Struct field via map | `m[k].Field = x` | `p := m[k]; p.Field = x; m[k] = p` |
| Nested write, no guard | `m[a][b] = 1` | init `m[a]` with `make()` first |
| Assuming map order | relying on iteration order | sort keys first |

---

## 🎓 Before Next Level

Can you:
- [ ] Create a map three ways and explain the nil zero value?
- [ ] Reproduce the exact nil-map-write panic message?
- [ ] Perform full CRUD on a map, including safe delete?
- [ ] Explain why comma-ok matters using a concrete example?
- [ ] Name three valid and three invalid map key types?
- [ ] Build a set with `map[T]struct{}` and explain why it beats `map[T]bool`?
- [ ] Explain and fix the map-of-struct-value addressability gotcha?
- [ ] Safely read and write a nested map?

If YES → You're ready for Level 11!

---

## 📚 Next Level

Level 11: Functions
- Declaring and calling functions
- Multiple return values
- Variadic parameters
- Closures and first-class functions

You've got maps down! 💪
