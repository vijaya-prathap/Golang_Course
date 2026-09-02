# Level 13: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: Structs
```
Day 1:  What a struct is, defining structs, zero values
Day 2:  Struct literals: keyed vs positional
Day 3:  Value semantics: copy vs pointer sharing
Day 4:  Comparing structs (== and its limits)
Day 5:  Nested structs
Day 6:  Embedded structs, field promotion, shadowing
Day 7:  Anonymous structs, struct tags, NewX constructors
```

### Week 2: Practice & Application
```
Day 1:  Exercises 1-3 (literals, zero values, copy vs pointer)
Day 2:  Exercises 4-6 (comparison, nesting, embedding)
Day 3:  Exercises 7-9 (anonymous structs, tags, constructors)
Day 4:  Exercise 10 (comprehensive library inventory)
Day 5:  Bonus challenges
Day 6-7: Practice & review
```

---

## 🎯 Struct Memory Layout

A struct is one contiguous block of memory - its fields laid out in order, one after another. Not pointers to fields elsewhere - the fields themselves, inline.

```
type Person struct {
    Name string
    Age  int
    City string
}

alice := Person{Name: "Alice", Age: 30, City: "Lisbon"}

┌─────────────────────────────────────────────┐
│                Person "alice"                │
├───────────────────┬────────┬─────────────────┤
│   Name             │  Age   │   City          │
│   "Alice"          │   30   │   "Lisbon"      │
└───────────────────┴────────┴─────────────────┘
        one struct value = all fields, together
```

This is exactly why assigning a struct copies everything at once - there's no separate "header" pointing elsewhere the way a slice has. Copy the box, you copy everything in it.

---

## 🗂️ Struct Literal Forms Compared

| Form | Example | Order matters? | Can skip fields? | Verdict |
|------|---------|-----------------|-------------------|---------|
| **Keyed** | `Person{Name: "A", Age: 3}` | No | Yes (skipped = zero) | ✅ preferred |
| **Positional** | `Person{"A", 3, "X"}` | Yes | No (must supply all) | ❌ avoid |
| **Mixed** | `Person{"A", Age: 3}` | - | - | ⛔ compile error |
| **Zero-value** | `var p Person` | - | all fields zero | ✅ great default |
| **Pointer literal** | `&Person{Name: "A"}` | - | Yes | ✅ for shared values |

```go
// Verified compile errors:

Person{"Alice", Age: 30, City: "Lisbon"}
// mixture of field:value and value elements in struct literal

Person{"Alice", 30}   // (Person has 3 fields)
// too few values in struct literal of type Person
```

---

## 🔁 Value Copy vs Pointer Sharing

```
COPY (assignment / value parameter):

    original: Point{X:1, Y:2}
         │
         │  copied := original
         ▼
    copied:   Point{X:1, Y:2}   ← independent box, own memory

    copied.X = 100
    ┌────────────┐        ┌────────────┐
    │ original   │        │ copied     │
    │ X:1  Y:2   │        │ X:100 Y:2  │   ← original untouched
    └────────────┘        └────────────┘


SHARE (pointer):

    original: Point{X:1, Y:2}
         ▲
         │  ptr := &original
    ptr ─┘   (ptr just holds original's address)

    ptr.X = 999   →  writes THROUGH to original
    ┌────────────┐
    │ original   │ ←──── ptr points here, no separate copy
    │ X:999 Y:2  │
    └────────────┘
```

Verified with `go run`:

```
original: {1 2}
copied:   {100 2}   (independent - copy)

original: {999 2}
via ptr:  {999 2}   (same memory - share)
```

**Rule:** `struct` assignment/value-params = copy (like `[N]T` arrays). `*struct` = share (like any pointer, Level 12).

---

## ⚖️ Comparable vs Uncomparable Field Types

A struct is comparable with `==` only if **every** field is comparable.

| Field type | Comparable? |
|-----------|-------------|
| `string`, `int`, `float64`, `bool`, other numerics | ✅ yes |
| `array` (`[N]T`, if T is comparable) | ✅ yes |
| pointer (`*T`) | ✅ yes (compares addresses) |
| another struct (if all ITS fields are comparable) | ✅ yes |
| `slice` (`[]T`) | ❌ no |
| `map[K]V` | ❌ no |
| `func` | ❌ no |

```go
type Team struct {
    Name    string
    Members []string // <- one slice field poisons the whole struct
}

t1 == t2
// invalid operation: t1 == t2 (struct containing []string cannot be compared)
```

If you need to compare such a struct, either compare fields manually (loop over the slice) or use `reflect.DeepEqual` (Level 26: Testing).

---

## 🧬 Embedding and Field Promotion

```go
type Animal struct {
    Name string
    Legs int
}

type Dog struct {
    Animal   // embedded - no field name, just the type
    Breed string
}
```

```
┌─────────────────────────────┐
│           Dog                │
│  ┌─────────────────────┐    │
│  │   Animal (embedded)  │    │
│  │   Name: "Rex"         │    │  ← promoted: d.Name reaches here
│  │   Legs: 4             │    │  ← promoted: d.Legs reaches here
│  └─────────────────────┘    │
│   Breed: "Border Collie"     │
└─────────────────────────────┘

d.Name         → Rex          (promoted shortcut)
d.Animal.Name  → Rex          (full explicit path - always works)
```

### Shadowing

```go
type Cat struct {
    Animal
    Name string // outer field with the SAME name as a promoted one
}

┌─────────────────────────────┐
│           Cat                 │
│  ┌─────────────────────┐    │
│  │  Animal (embedded)    │    │
│  │  Name: "some animal"  │    │  ← still here, reached via c.Animal.Name
│  │  Legs: 4              │    │
│  └─────────────────────┘    │
│   Name: "Whiskers"           │  ← THIS wins for c.Name (shadows)
└─────────────────────────────┘

c.Name          → Whiskers     (outer field wins)
c.Animal.Name   → some animal  (inner field still reachable)
```

**Note:** Field promotion works only in one direction - the outer type gains access to the inner type's fields, never the reverse. And promotion is for access, not literals: you must still write `Dog{Animal: Animal{...}}` when constructing.

---

## 📊 Nested vs Embedded — Side by Side

| | Nested (named field) | Embedded (anonymous field) |
|---|---|---|
| Declaration | `Address Address` | `Animal` (just the type) |
| Access | `p.Address.City` (always full path) | `d.Name` (promoted) or `d.Animal.Name` (full path) |
| Relationship | "has a" | "made of / composed of" |
| Method promotion (Level 14) | none | inner type's methods promoted too |
| When to use | data grouping, no promotion desired | reusing behavior/fields as if declared directly |

---

## 🚨 Common Mistakes

### Mistake 1: Assuming Assignment Shares

```go
// ❌ WRONG ASSUMPTION - b is independent
b := alice
b.Age = 99
fmt.Println(alice.Age) // still 30

// ✅ RIGHT - use a pointer for sharing
b := &alice
```

### Mistake 2: Mutating Through a Value Parameter

```go
// ❌ WRONG - modifies a copy
func birthday(p Person) { p.Age++ }

// ✅ RIGHT
func birthday(p *Person) { p.Age++ }
```

### Mistake 3: Positional Literal Field Mix-Ups

```go
// ❌ Ambiguous to a reader - which is which?
r := Rect{5, 10}

// ✅ Self-documenting
r := Rect{Width: 10, Height: 5}
```

### Mistake 4: == on a Struct With a Slice/Map Field

```go
// ❌ does not compile
type Team struct{ Members []string }
t1 == t2
```

### Mistake 5: Forgetting the Embedded Type Key in a Literal

```go
// ❌ unknown field Name in struct literal of type Dog
d := Dog{Name: "Rex"}

// ✅ RIGHT
d := Dog{Animal: Animal{Name: "Rex"}}
```

### Mistake 6: Missing a Shadowed Field

```go
// c.Name and c.Animal.Name are DIFFERENT fields once Cat declares its own Name
```

---

## 📈 Progression Summary

### Understanding Level 13

Level 13 teaches how to group related data into one named, composite type:

1. **Definition & zero values** — every field defaults to its own zero value
2. **Literals** — keyed (preferred) vs positional (fragile)
3. **Value semantics** — copy on assignment/pass; pointer for sharing
4. **Comparison** — `==` works only when every field is comparable
5. **Composition** — nested structs (has-a) and embedded structs (field promotion)
6. **Extras** — anonymous structs, struct tags, NewX constructors

### Prerequisites for Level 14

Before moving to Level 14 (Methods), you need:

- ✅ Comfortable defining structs and creating them with keyed literals
- ✅ Can explain and demonstrate copy-on-assignment vs pointer sharing
- ✅ Know when `==` compiles on a struct and when it doesn't (and why)
- ✅ Comfortable with nested structs and embedded structs with promotion
- ✅ Understand shadowing when an outer field name matches a promoted one
- ✅ Can write a NewX constructor function

### Ready for Level 14?

Level 14 attaches **behavior** to the struct types you just learned to define:
- Methods with value vs pointer receivers (this level's copy-vs-pointer lesson applies directly)
- Method promotion through embedding (this level's field promotion, extended to methods)
- Method sets - the concept that leads straight into interfaces (Level 15)

---

## ✅ Checklist Before Level 14

- [ ] Can define a struct type and create it with every literal form
- [ ] Can state every field's zero value in a `var`-declared struct
- [ ] Can prove struct assignment copies, and pointer assignment shares
- [ ] Can explain why a struct with a slice/map field loses `==`
- [ ] Can build a nested struct and access/mutate deep fields
- [ ] Can embed a struct, use a promoted field, and explain shadowing
- [ ] Can write and use a NewX constructor
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The Box
A struct is fields laid out together in one block of memory - copying the struct copies the whole block.

### The Choice
Value struct = independent copy. Pointer to struct = shared mutation. The same choice you learned for parameters in Level 12, now applied to your own types.

### The Composition
Embedding promotes fields (and, from Level 14 on, methods) - Go's alternative to inheritance. No "is-a," just "has-a, and borrows its stuff."

### The Convention
No constructors in the language - `NewX(...) *X` is the idiom, used only when the zero value isn't good enough on its own.

---

## 📚 Next Level

Level 14: Methods
- Attaching behavior to your struct types
- Value receivers vs pointer receivers
- Method promotion through embedding
- Method sets - the doorway to interfaces

You've got structs down! Keep going! 🚀
