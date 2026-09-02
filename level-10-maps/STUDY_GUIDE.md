# Level 10: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: Maps In Depth
```
Day 1:  What a map is, creating maps (literal, make, nil)
Day 2:  The nil map read/write asymmetry
Day 3:  CRUD operations
Day 4:  comma-ok idiom deep dive
Day 5:  Map key requirements
Day 6:  Maps as sets
Day 7:  Maps of structs vs pointers, nested maps
```

### Week 2: Practice & Application
```
Day 1:  Exercises 1-3 (creation, nil panic, CRUD)
Day 2:  Exercises 4-6 (comma-ok, keys, sets)
Day 3:  Exercises 7-8 (struct gotcha, nested maps)
Day 4:  Exercises 9-10 (inverted index, inventory system)
Day 5:  Bonus challenges
Day 6-7: Practice & review
```

---

## 🎯 Ways to Create a Map

| Method | Syntax | Result |
|---|---|---|
| **Literal** | `map[string]int{"a": 1}` | Ready to read and write |
| **make** | `make(map[string]int)` | Empty, ready to read and write |
| **var (zero value)** | `var m map[string]int` | `nil` - readable, panics on write |

```
map[string]int{"a": 1}   ready immediately
make(map[string]int)     ready immediately, starts empty
var m map[string]int     m == nil -> read OK, write PANICS
```

---

## 🗺️ Map Operations Reference

| Operation | Syntax | Behavior on Missing Key | Behavior on nil Map |
|---|---|---|---|
| Insert / Update | `m[k] = v` | Creates the entry | **Panics** |
| Read (bare) | `v := m[k]` | Returns zero value | Returns zero value (safe) |
| Read (comma-ok) | `v, ok := m[k]` | `ok` is `false` | `ok` is `false` (safe) |
| Delete | `delete(m, k)` | No-op, no panic | No-op, no panic (safe) |
| Length | `len(m)` | n/a | Returns `0` (safe) |
| Iterate | `for k, v := range m` | n/a | Zero iterations (safe) |

**The one operation that panics on a nil map is a write (`m[k] = v`). Every other operation is completely safe.**

---

## ⚡ The nil Map Asymmetry (Diagram)

```
                    nil map (var m map[K]V, never made)
                              |
              +---------------+---------------+
              |                               |
           READ                            WRITE
              |                               |
     m[k]  -> zero value                m[k] = v
     v,ok:=m[k] -> ok=false                   |
     len(m) -> 0                              v
     range m -> 0 iterations          panic: assignment to
     delete(m,k) -> no-op              entry in nil map
              |
              v
        ALL SAFE - no
        allocation needed
        to answer "not found"
```

A nil map has no hash table allocated. Answering "is this key here?" needs no storage - the answer is always "no." But *storing* a new entry needs somewhere to put it, and Go refuses to allocate that for you silently - you must call `make()` (or use a literal) first.

---

## ✅❌ Comma-ok Truth Table

```go
v, ok := m[k]
```

| Key state | `ok` | `v` | Meaning |
|---|---|---|---|
| Key present, value is non-zero | `true` | the real value | Found, as expected |
| Key present, value IS the zero value | `true` | zero value | Found - the zero value was stored on purpose |
| Key absent | `false` | zero value | Not found - `v` is meaningless, always zero-valued |
| Key absent, map is `nil` | `false` | zero value | Not found - same as any absent key |

**The critical row is the second one.** Without `ok`, you cannot tell "present, value happens to be zero" apart from "absent." Both print the exact same zero value from a bare index.

```go
stock := map[string]int{"bananas": 0} // present, value 0

stock["bananas"]        // 0 - present
stock["grapes"]          // 0 - absent
// Identical output! Only comma-ok can tell them apart:

_, ok1 := stock["bananas"] // ok1 = true
_, ok2 := stock["grapes"]  // ok2 = false
```

---

## 🔑 Valid vs Invalid Map Key Types

| Key type | Valid? | Reason |
|---|---|---|
| `int`, `string`, `float64`, `bool`, ... | ✅ | Comparable with `==` |
| `[N]T` array (comparable `T`) | ✅ | Arrays are comparable if their elements are |
| `struct{...}` (all fields comparable) | ✅ | Structs are comparable if every field is |
| `*T` pointer | ✅ | Pointers compare by address |
| `error` / other comparable interfaces | ✅ | Compares dynamic type + value |
| `[]T` slice | ❌ | Slices are never comparable |
| `map[K]V` | ❌ | Maps are never comparable |
| `func(...)` | ❌ | Functions are never comparable |

```
Real compiler error for an invalid key (captured verbatim):

    ./main.go:6:13: invalid map key type []int

Real compiler error for the map-of-struct-value field gotcha (captured verbatim):

    ./main.go:15:2: cannot assign to struct field players["alice"].Score in map
```

### Decision Tree

```
Choosing a map key type?
├─ Is it a basic type (int, string, bool, float, ...)?
│  └─ ✅ Always fine
├─ Is it an array or struct?
│  └─ ✅ Fine, AS LONG AS every element/field is itself comparable
├─ Is it a slice, map, or function?
│  └─ ❌ Compile error - use an array instead of a slice,
│         or a key derived from the data (e.g., a string ID)
```

---

## 🧱 Maps as Sets

| Idiom | Memory per entry | Ambiguous "false but present" state? |
|---|---|---|
| `map[T]bool` | 1 byte (the bool) | Yes - `m[x] = false` looks different from absence, but shouldn't |
| `map[T]struct{}` | **0 bytes** | No - presence in the map IS membership |

```go
// map[T]bool
set["x"] = true       // add
set["x"] = false      // "removed"? Or just marked false? Ambiguous!
delete(set, "x")      // the ONLY unambiguous removal

// map[T]struct{}
set["x"] = struct{}{} // add
delete(set, "x")      // remove - the only state that exists is "present"
_, in := set["x"]     // membership test
```

```
unsafe.Sizeof(bool(true))    = 1 byte
unsafe.Sizeof(struct{}{})    = 0 bytes   <- verified with unsafe.Sizeof
```

---

## 🧩 Maps of Structs vs Pointers to Structs

```
map[string]Player          map[string]*Player
       |                            |
players["a"].Score = 10     players["a"].Score = 10
       |                            |
  COMPILE ERROR              COMPILES FINE
  (map value not                    |
   addressable)              (pointer dereference
                               IS addressable)
```

| You have | `m[k].Field = x` compiles? | Why |
|---|---|---|
| `map[K]StructType` | ❌ No | `m[k]` is a copy, not a real address |
| `map[K]*StructType` | ✅ Yes | `m[k]` is a pointer; dereferencing it reaches the real struct |
| `[]StructType` (slice) | ✅ Yes | Slice elements ARE addressable |

### The Two Fixes for `map[K]StructType`

```go
// Fix 1: read-modify-reassign
p := m[k]
p.Field = x
m[k] = p

// Fix 2: switch to a map of pointers
m := map[string]*Player{"a": {Score: 0}}
m["a"].Score = 10 // fine
```

---

## 🪆 Nested Maps

```
map[string]map[string]int
        |            |
     outer key    inner map (map[string]int)

READ through a missing outer key  -> safe, inner map is nil, gives zero value
WRITE through a missing outer key -> PANIC (nil-map-write, same rule, one level deeper)
```

```go
population["France"]["Paris"]       // 0, safe (France not in map)
population["France"]["Paris"] = 1   // PANIC - population["France"] is nil

// Safe write pattern:
if population["France"] == nil {
    population["France"] = make(map[string]int)
}
population["France"]["Paris"] = 1   // fine now
```

---

## 🚨 Common Mistakes

### Mistake 1: Writing to a nil Map

```go
// ❌ WRONG - panic: assignment to entry in nil map
var m map[string]int
m["key"] = 1

// ✅ RIGHT
m := make(map[string]int)
```

### Mistake 2: Bare Index Instead of comma-ok

```go
// ❌ WRONG - 0 could mean "absent" or "stored zero"
if m["key"] != 0 { }

// ✅ RIGHT
if _, ok := m["key"]; ok { }
```

### Mistake 3: Slice (or Map, or Func) as a Key

```go
// ❌ WRONG - invalid map key type []int
bad := map[[]int]string{}

// ✅ RIGHT - use a fixed-size array
ok := map[[3]int]string{}
```

### Mistake 4: Assigning to a Struct Field Through a Map

```go
// ❌ WRONG - cannot assign to struct field ... in map
players := map[string]Player{"a": {}}
players["a"].Score = 10

// ✅ RIGHT
p := players["a"]
p.Score = 10
players["a"] = p
```

### Mistake 5: Writing Through an Uninitialized Inner Map

```go
// ❌ WRONG - panics if grid["a"] was never created
grid := make(map[string]map[string]int)
grid["a"]["b"] = 1

// ✅ RIGHT
if grid["a"] == nil {
    grid["a"] = make(map[string]int)
}
grid["a"]["b"] = 1
```

### Mistake 6: Assuming Map Iteration Order Is Stable

```go
// ❌ WRONG ASSUMPTION - carried forward from Level 6
for k := range m {
    fmt.Println(k) // order varies between runs
}
```

---

## 📈 Progression Summary

### Understanding Level 10

Level 10 goes far beyond the map basics from Levels 3, 5, and 6:

1. **Creation** - literal, `make`, and the dangerous nil zero value
2. **The nil asymmetry** - safe reads, panicking writes
3. **CRUD** - insert, read, update, delete, all with exact semantics
4. **comma-ok** - the only reliable way to distinguish absent from zero
5. **Key requirements** - comparability, and what breaks the compiler
6. **Sets** - `map[T]struct{}` as the zero-memory idiom
7. **Addressability** - why `m[k].Field = x` fails for struct values
8. **Nested maps** - safe reads, guarded writes

### Prerequisites for Level 11

Before moving to Level 11 (Functions), you need:

- ✅ Comfortable creating maps all three ways and explaining the nil zero value
- ✅ Can explain (and reproduce) the nil-map-write panic
- ✅ Fluent with comma-ok for presence checks
- ✅ Know which types are valid map keys and why
- ✅ Can build a set with `map[T]struct{}`
- ✅ Understand why map values aren't addressable and how to work around it
- ✅ Can safely read and write nested maps

### Ready for Level 11?

Level 11 teaches how to package logic into reusable functions:
- Declaring and calling functions
- Multiple return values (you've already used this pattern via comma-ok!)
- Variadic parameters
- Closures and first-class functions

---

## ✅ Checklist Before Level 11

- [ ] Can create a map with a literal, with `make`, and explain what `var m map[K]V` gives you
- [ ] Can explain why reading a nil map is safe but writing panics, with the exact panic message
- [ ] Can perform insert, read, update, and delete on a map
- [ ] Can use comma-ok to distinguish an absent key from a zero-valued one
- [ ] Can name which types are valid and invalid map keys, and why
- [ ] Can build a set using `map[T]struct{}` and explain why it beats `map[T]bool`
- [ ] Can explain (and fix) the map-of-struct-value addressability gotcha
- [ ] Can safely initialize and write to a nested map
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The Central Rule
Reading a nil map is always safe. Writing to one always panics with `assignment to entry in nil map`.

### The Idiom
`v, ok := m[k]` is the only reliable way to know whether a key is really there.

### The Keys
Comparable types only - basic types, arrays, and structs work; slices, maps, and funcs don't compile.

### The Sets
`map[T]struct{}` costs zero bytes per entry and has no ambiguous "false but present" state.

### The Gotcha
Map values aren't addressable - mutate struct fields by reassigning the whole struct, or store pointers instead.

---

## 📚 Next Level

Level 11: Functions
- Declaring and calling functions
- Multiple return values
- Variadic parameters
- Closures and first-class functions

You've mastered maps! Keep going! 🚀
