# Level 13: Structs - Complete Guide

## Introduction

Welcome to Level 13! You've mastered pointers (Level 12) - `&`, `*`, and the difference between copying a value and sharing it. Now it's time to learn **structs**: Go's way of grouping related pieces of data into one named type.

If you've heard that "Go doesn't have classes" - this level is where you learn what Go has instead. A struct is the **data half** of Go's object story. Level 14 (Methods) adds the **behavior half**. Together, structs + methods + interfaces (Level 15) are Go's complete answer to object-oriented programming - built on **composition, not inheritance**.

---

## Table of Contents

1. [What Is a Struct?](#what-is-a-struct)
2. [Defining Structs and Zero Values](#defining-structs-and-zero-values)
3. [Struct Literals: Keyed vs Positional](#struct-literals-keyed-vs-positional)
4. [Structs Are Value Types: Copying vs Sharing](#structs-are-value-types-copying-vs-sharing)
5. [Comparing Structs](#comparing-structs)
6. [Nested Structs](#nested-structs)
7. [Embedded Structs and Field Promotion](#embedded-structs-and-field-promotion)
8. [Anonymous Structs](#anonymous-structs)
9. [Struct Tags: A First Look](#struct-tags-a-first-look)
10. [Constructor Functions: The NewX Convention](#constructor-functions-the-newx-convention)
11. [Best Practices](#best-practices)
12. [Common Mistakes](#common-mistakes)

---

## What Is a Struct?

So far, related data has lived in separate variables:

```go
name := "Alice"
age := 30
city := "Lisbon"
```

That works for one person. But what about a slice of people? Three parallel slices you have to keep in sync? That's fragile. A **struct** groups the fields into a single value with one name:

```go
type Person struct {
    Name string
    Age  int
    City string
}
```

- `type Person struct` declares a new **named type** (just like `type Celsius float64` would, but with structure inside)
- Each **field** has a name and a type
- Fields with the same type can share a line: `X, Y int`
- A struct is a **composite type**: one value made of other values, accessed by *name* (unlike arrays/slices, which are accessed by *index*)

Think of a struct as a record, a row in a table, or a form with labeled boxes. It's Go's building block for modeling real things: users, servers, orders, books, coordinates.

---

## Defining Structs and Zero Values

Declare a struct variable with `var` and you get the **zero-value struct**: every field is set to its own zero value. No constructors, no `null` surprises - it just works.

```go
type Person struct {
    Name string
    Age  int
}

var p Person
fmt.Printf("%+v\n", p)                // {Name: Age:0}
fmt.Println(p.Name == "", p.Age == 0) // true true
```

Verified with `go run` - the zero values are exactly what Level 3 taught for each type:

| Field type | Zero value |
|-----------|------------|
| `string` | `""` |
| `int`, `float64` | `0` |
| `bool` | `false` |
| slices, maps, pointers | `nil` |
| nested struct | that struct's zero value (recursively!) |

A zero-value struct is immediately usable - assign to its fields and go. Many great Go types (like `bytes.Buffer` and `sync.Mutex` in later levels) are designed so their zero value is ready to use. Aim for that with your own types.

### Printing Structs

Three format verbs you'll use constantly (verified output):

```go
alice := Person{Name: "Alice", Age: 30, City: "Lisbon"}
fmt.Printf("%v\n", alice)  // {Alice 30 Lisbon}
fmt.Printf("%+v\n", alice) // {Name:Alice Age:30 City:Lisbon}
fmt.Printf("%#v\n", alice) // main.Person{Name:"Alice", Age:30, City:"Lisbon"}
```

`%+v` (values *with* field names) is the debugging workhorse.

---

## Struct Literals: Keyed vs Positional

There are two literal forms - and one of them is strongly preferred.

### Keyed Literals (✅ preferred)

```go
alice := Person{Name: "Alice", Age: 30, City: "Lisbon"}
```

- Order doesn't matter
- You can skip fields - **skipped fields get their zero value**
- Adding a new field to the struct later doesn't break this code

```go
carol := Person{Name: "Carol"}
fmt.Println(carol) // {Carol 0 }   <- Age and City are zero-valued
```

### Positional Literals (❌ avoid)

```go
bob := Person{"Bob", 25, "Porto"}
```

- Values must be in **exact field order**
- You must provide **every** field - or it won't compile:

```go
p := Person{"Alice", 30} // missing City - compile error!
```

```
too few values in struct literal of type Person
```

- If someone reorders or inserts a field in the struct definition, every positional literal silently means something different (or breaks)

You also can't mix the two forms in one literal (verified compile error):

```
mixture of field:value and value elements in struct literal
```

**Rule of thumb:** keyed literals everywhere, positional only for tiny, obvious structs that will never change (like `image.Point{2, 3}`).

### Pointer to a Fresh Struct

From Level 12, `&` in front of a literal gives you a pointer to a brand-new struct - the idiomatic way to create shared struct values:

```go
p := &Person{Name: "Alice", Age: 30} // p is a *Person
```

---

## Structs Are Value Types: Copying vs Sharing

This is the most important section in this level. **Structs are value types, like arrays (Level 7) - not like slices and maps.** Assigning a struct copies *all* of its fields.

Verified with `go run` on this Go toolchain:

```go
type Point struct {
    X int
    Y int
}

original := Point{X: 1, Y: 2}
copied := original // full copy of every field
copied.X = 100

fmt.Println(original) // {1 2}   <- unchanged!
fmt.Println(copied)   // {100 2}
```

```
original: {1 2}
copied:   {100 2}
```

The same is true when passing a struct to a function - the function gets a copy (exactly like value parameters in Level 12):

```go
func tryToMutate(p Point) {
    p.X = -1 // modifies the COPY only
}

tryToMutate(original)
fmt.Println(original) // {1 2} - still unchanged
```

### Pointer to Struct: Shared Mutation

When you *want* changes to be visible to the caller, pass a pointer - the bridge from Level 12:

```go
func mutate(p *Point) {
    p.Y = -1 // modifies the caller's struct
}

ptr := &original
ptr.X = 999 // shorthand for (*ptr).X = 999
mutate(&original)
fmt.Println(original) // {999 -1}
```

Note the huge convenience: **Go auto-dereferences struct pointers for field access**. You write `ptr.X`, never `(*ptr).X` - both work (verified), but nobody writes the second form.

### Which Should You Use?

| Situation | Use |
|-----------|-----|
| Small struct, no mutation needed | value (`Point`) |
| Function must modify the caller's struct | pointer (`*Point`) |
| Large struct passed around a lot | pointer (avoids big copies) |
| Struct stored in shared state | pointer |

Level 14 revisits this exact question for **method receivers** - the intuition you build here carries straight over.

---

## Comparing Structs

Structs support `==` and `!=` **when every field is comparable**. Two structs are equal when all corresponding fields are equal.

```go
a := Point{X: 1, Y: 2}
b := Point{X: 1, Y: 2}
c := Point{X: 5, Y: 2}

fmt.Println(a == b) // true
fmt.Println(a == c) // false
```

A handy idiom - checking for "still the zero value":

```go
var zero Point
fmt.Println(zero == Point{}) // true
```

And because comparable structs are, well, comparable, they work as **map keys** (verified):

```go
m := map[Point]string{
    {0, 0}: "origin",
    {1, 2}: "a",
}
fmt.Println(m[Point{1, 2}]) // a
```

### When Comparison Is a Compile Error

Slices, maps, and functions are **not** comparable (remember from Level 8: slices only compare to `nil`). A struct containing any of them loses `==` entirely:

```go
type Team struct {
    Name    string
    Members []string // a slice field - Team is no longer comparable
}

t1 := Team{Name: "Go", Members: []string{"Alice", "Bob"}}
t2 := Team{Name: "Go", Members: []string{"Alice", "Bob"}}
fmt.Println(t1 == t2) // does not compile!
```

The real compiler error, captured verbatim:

```
invalid operation: t1 == t2 (struct containing []string cannot be compared)
```

To compare such structs you compare field by field yourself (looping over the slices with Level 6 patterns), or use `reflect.DeepEqual` (you'll meet it in Level 26: Testing).

---

## Nested Structs

A struct field can be another struct - modeling has-a relationships (a Person *has an* Address):

```go
type Address struct {
    Street string
    City   string
    Zip    string
}

type Person struct {
    Name    string
    Age     int
    Address Address // nested struct field
}
```

### Initialization

Nested literals just nest (keyed, of course):

```go
alice := Person{
    Name: "Alice",
    Age:  30,
    Address: Address{
        Street: "12 Rose Lane",
        City:   "Lisbon",
        Zip:    "1100-048",
    },
}
```

### Access and Mutation - Chain the Dots

```go
fmt.Println(alice.Address.City) // Lisbon
alice.Address.City = "Porto"    // mutate deep fields directly
```

### Nested Structs Copy With Their Parent

Value semantics go all the way down (verified):

```go
bob := alice          // copies Person AND the Address inside it
bob.Address.City = "Faro"

fmt.Println(alice.Address.City) // Porto - alice is untouched
fmt.Println(bob.Address.City)   // Faro
```

---

## Embedded Structs and Field Promotion

Nesting with a **field name** (`Address Address`) is "has-a". Embedding drops the field name - just the type - and Go does something special:

```go
type Animal struct {
    Name string
    Legs int
}

type Dog struct {
    Animal // embedded (anonymous field) - no field name!
    Breed  string
}
```

The embedded struct's fields are **promoted**: you can access them directly on the outer struct, as if they were declared there.

```go
d := Dog{
    Animal: Animal{Name: "Rex", Legs: 4},
    Breed:  "Border Collie",
}

fmt.Println(d.Name)        // Rex  - promoted! (short form)
fmt.Println(d.Animal.Name) // Rex  - full path still works
d.Legs = 3                 // writes through the promotion too
```

Verified output:

```
d.Name (promoted):         Rex
d.Legs (promoted):         4
d.Animal.Name (full path): Rex
```

Two things to notice:

1. In the **literal**, the embedded field is keyed by its *type name* (`Animal: Animal{...}`) - promotion is for *access*, not initialization.
2. This is **composition, not inheritance**. A `Dog` is not an `Animal` subtype - it *contains* an `Animal` and borrows convenient access to its fields. There is no "is-a" in Go's type system.

### Shadowing: An Explicit Field Wins

If the outer struct declares a field with the same name, the outer field **shadows** the promoted one. The inner field is still there - reachable through the full path (verified):

```go
type Cat struct {
    Animal
    Name string // shadows the promoted Animal.Name
}

c := Cat{
    Animal: Animal{Name: "some animal", Legs: 4},
    Name:   "Whiskers",
}

fmt.Println(c.Name)        // Whiskers    - the outer field wins
fmt.Println(c.Animal.Name) // some animal - inner field still reachable
```

### Why Embedding Matters

Right now, embedding promotes *fields*. In **Level 14 you'll see it also promotes methods** - embed an `Animal` with a `Speak()` method, and `Dog` gets `Speak()` for free. That's Go's composition-over-inheritance story, and it's why embedding is everywhere in real Go code.

---

## Anonymous Structs

Sometimes a struct type is needed exactly once - too small and local to deserve a `type` declaration. Define it inline:

```go
server := struct {
    Host string
    Port int
}{
    Host: "localhost",
    Port: 8080,
}

fmt.Printf("%+v\n", server) // {Host:localhost Port:8080}
```

The type and the value appear together: `struct { ...fields... }{ ...values... }`.

The most common real-world use is a **slice of anonymous structs** - grouping rows of related data, the exact pattern used by table-driven tests (Level 26):

```go
tests := []struct {
    Input    int
    Expected int
}{
    {Input: 2, Expected: 4},
    {Input: 3, Expected: 9},
    {Input: 5, Expected: 25},
}

for _, tc := range tests {
    fmt.Println(tc.Input, tc.Expected)
}
```

**Rule of thumb:** if the struct is used in one place, anonymous is fine. The moment it appears twice, give it a name with `type`.

---

## Struct Tags: A First Look

You'll often see backtick strings after struct fields:

```go
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age,omitempty"`
}
```

Those are **struct tags** - metadata attached to fields:

- Syntax: a raw string in backticks, conventionally `key:"value"` pairs separated by spaces
- Tags do **nothing** by themselves - your program runs identically with or without them
- Libraries read them at runtime via **reflection** (`reflect.StructField.Tag`) and change *their* behavior: `encoding/json` uses `json:"name"` to decide what a field is called in JSON, database libraries use `db:"..."`, validation libraries use `validate:"..."`

That's all you need for now. **Level 19 (JSON) uses struct tags heavily** - this level just teaches you to recognize the syntax so it never looks mysterious.

---

## Constructor Functions: The NewX Convention

Go has no constructors. The convention instead: an ordinary function named `New<Type>` that validates/defaults the fields and returns a **pointer** to a ready-to-use struct:

```go
type BankAccount struct {
    Owner   string
    Balance float64
    Active  bool
}

func NewBankAccount(owner string, openingBalance float64) *BankAccount {
    if openingBalance < 0 {
        openingBalance = 0 // enforce a sensible default
    }
    return &BankAccount{
        Owner:   owner,
        Balance: openingBalance,
        Active:  true, // every account starts active
    }
}

acct := NewBankAccount("Alice", 250.75)
fmt.Printf("%+v\n", *acct) // {Owner:Alice Balance:250.75 Active:true}
```

Why this pattern?

- **One place** to put validation and defaults instead of scattering them at every creation site
- Returning `*T` (safe thanks to Go's escape analysis - Level 12) means every caller shares the same struct
- If the package's struct is the main thing it exports, the function is often just `New()` (e.g., `errors.New`, and `bufio.NewReader`)
- Constructors that can fail return `(*T, error)` - the pattern you'll master in Level 16

**When do you need one?** Only when the zero value isn't good enough. If `var b Buffer` works, don't write `NewBuffer`.

---

## Best Practices

### 1. Always Use Keyed Struct Literals

```go
// ✅ Good - survives field reordering and additions
p := Person{Name: "Alice", Age: 30}

// ❌ Fragile - breaks or silently changes meaning if the struct evolves
p := Person{"Alice", 30, "Lisbon"}
```

### 2. Design for a Useful Zero Value

```go
// ✅ Good - var c Counter is immediately usable
type Counter struct {
    Count int // zero value 0 is a valid starting count
}
```

When the zero value can't be valid (required fields, invariants), provide a `NewX` constructor - and make it the obvious way to build the type.

### 3. Prefer %+v for Debug Printing

```go
fmt.Printf("%+v\n", user) // field names included - readable weeks later
```

### 4. Use Pointers Deliberately, Not Reflexively

Small struct, read-only use? Pass the value. Mutation or big struct? Pass the pointer. Say what you mean - a `*Person` parameter tells the reader "I may modify this."

### 5. Reach for Embedding Only for Real "Made Of" Relationships

Embed when the outer type genuinely *is composed of* the inner one and you want its fields (and later, methods) at the top level. If you just need the data, a plain named field (`Address Address`) is clearer.

### 6. Keep Structs Focused

A struct with 25 fields is usually several structs in a trench coat. Group related fields into nested structs (`Address`, `Metadata`) - it documents structure for free.

---

## Common Mistakes

### Mistake 1: Expecting Assignment to Share (It Copies!)

```go
// ❌ WRONG ASSUMPTION - b is an independent copy
b := alice
b.Age = 99
fmt.Println(alice.Age) // still 30, NOT 99

// ✅ RIGHT - use a pointer when you want sharing
b := &alice
b.Age = 99
fmt.Println(alice.Age) // 99
```

### Mistake 2: Mutating a Struct Through a Value Parameter

```go
// ❌ WRONG - modifies a copy, caller sees nothing
func birthday(p Person) {
    p.Age++
}

// ✅ RIGHT - pointer parameter mutates the caller's struct
func birthday(p *Person) {
    p.Age++
}
```

### Mistake 3: Positional Literal Field Mix-Ups

```go
type Rect struct {
    Width, Height int
}

// ❌ COMPILES but is silently wrong if you meant Width=10
r := Rect{5, 10} // which is which? future readers must check the type

// ✅ RIGHT
r := Rect{Width: 10, Height: 5}
```

### Mistake 4: Using == on a Struct With Slice/Map Fields

```go
// ❌ WRONG - does not compile
type Team struct {
    Name    string
    Members []string
}
t1 == t2
// invalid operation: t1 == t2 (struct containing []string cannot be compared)

// ✅ RIGHT - compare field by field (loop over the slice), or
// use reflect.DeepEqual (Level 26)
```

### Mistake 5: Forgetting the Type Key When Initializing an Embedded Struct

```go
// ❌ WRONG - there is no field called "Name" directly in the literal
d := Dog{Name: "Rex"} // unknown field Name in struct literal of type Dog

// ✅ RIGHT - promotion works for ACCESS, not for literals
d := Dog{Animal: Animal{Name: "Rex", Legs: 4}, Breed: "Collie"}
fmt.Println(d.Name) // access via promotion is fine
```

### Mistake 6: Not Noticing a Shadowed Field

```go
// ❌ SURPRISE - c.Name is the outer field, the embedded one is hidden
type Cat struct {
    Animal
    Name string
}

// ✅ Be explicit when both exist - use the full path for the inner one
c.Name        // outer field
c.Animal.Name // embedded field
```

### Mistake 7: Confusing a Struct Type With a Struct Value

```go
// ❌ WRONG - Person is a type, not a value
fmt.Println(Person) // Person (type) is not an expression

// ✅ RIGHT - create a value of the type first
fmt.Println(Person{Name: "Alice"})
```

---

## Summary

**Defining and Creating:**
- `type Name struct { Field Type ... }` declares a named composite type
- `var p Person` gives a zero-value struct - every field zero, immediately usable
- Keyed literals `Person{Name: "A"}` are the norm; skipped fields are zero
- `&Person{...}` creates a struct and hands you a pointer to it

**Value Semantics:**
- Assignment and function calls **copy** structs (like arrays; unlike slices/maps)
- Use `*Person` when you want shared mutation - `ptr.Field` auto-dereferences
- `==` works when every field is comparable; slice/map/function fields remove it

**Composition:**
- Nested structs (named fields) model has-a; access with chained dots
- Embedded structs (anonymous fields) **promote** their fields to the outer struct
- An explicit outer field shadows a promoted one; the full path still reaches it
- Method promotion (Level 14) is what makes embedding really shine

**Conveniences:**
- Anonymous structs for one-off groupings and table-style data
- Struct tags = backtick metadata read by libraries via reflection (Level 19)
- `NewX(...) *X` constructor functions centralize validation and defaults

---

## Next Steps

You now understand:
- ✅ Defining structs and creating them with every literal form
- ✅ Zero-value structs and why they're a design goal
- ✅ Copy-on-assignment vs pointer sharing
- ✅ Struct comparison rules (and the compile error when fields aren't comparable)
- ✅ Nested structs, embedding, field promotion, and shadowing
- ✅ Anonymous structs, struct tags, and the NewX convention

**Next level:** Level 14 - Methods
- Attaching behavior to your struct types with methods
- Value receivers vs pointer receivers (your Level 12+13 knowledge pays off here)
- Method promotion through embedding
- Method sets and how they prepare you for interfaces (Level 15)

Your data now has structure. Next, it learns to do things! 🚀
