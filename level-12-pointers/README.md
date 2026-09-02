# Level 12: Pointers - Complete Guide

## Introduction

Welcome to Level 12! You've completed the FUNDAMENTALS tier (Levels 0-11) - congratulations! This level begins the **OOP CONCEPTS** tier (Levels 12-16: Pointers, Structs, Methods, Interfaces, Error Handling), and pointers are its foundation.

A **pointer** is a variable that holds the *memory address* of another value. That one idea unlocks the rest of the tier: structs get mutated through pointers, methods choose between value and pointer receivers, and interfaces care about which one you used.

If you've heard scary stories about pointers in C - relax. Go keeps the useful part (efficient sharing and mutation) and removes the dangerous part: there is **no pointer arithmetic** in Go. You can't wander off the end of an array through a pointer, and the garbage collector keeps every pointed-at value alive for exactly as long as you need it.

---

## Table of Contents

1. [What Is a Pointer?](#what-is-a-pointer)
2. [The Address-Of and Dereference Operators](#the-address-of-and-dereference-operators)
3. [The Zero Value of a Pointer: nil](#the-zero-value-of-a-pointer-nil)
4. [Passing by Value vs Passing a Pointer](#passing-by-value-vs-passing-a-pointer)
5. [Pointers to Structs](#pointers-to-structs)
6. [new(T) vs &T{}](#newt-vs-t)
7. [Returning a Pointer to a Local Variable Is Safe](#returning-a-pointer-to-a-local-variable-is-safe)
8. [When You Do NOT Need Pointers](#when-you-do-not-need-pointers)
9. [nil-Safety Patterns](#nil-safety-patterns)
10. [Best Practices](#best-practices)
11. [Common Mistakes](#common-mistakes)

---

## What Is a Pointer?

Every variable lives somewhere in memory, and that somewhere has an **address**. A pointer is simply a variable whose value *is* an address.

```
Regular variable:          Pointer variable:

x ┌──────────┐             p ┌──────────────────┐
  │    42    │               │  0x735196f9a020  │ ──points at──▶ x
  └──────────┘               └──────────────────┘
  at address                 "the address where
  0x735196f9a020              x lives"
```

### Why Does Go Have Pointers?

1. **Mutation across function boundaries.** Go passes everything by value (Level 11) - functions get *copies*. A pointer lets a function reach back and modify the caller's actual variable.
2. **Efficient sharing.** Copying a large struct on every call is wasteful. Copying a pointer is one machine word, no matter how big the pointed-at value is.
3. **Expressing "optional" or "not set".** A pointer can be `nil`, which is different from any real value - useful for "this field was never provided".

### Why Go Pointers Are NOT C Pointers

| C pointers | Go pointers |
|------------|-------------|
| Pointer arithmetic (`p++`, `p + 5`) | **No pointer arithmetic** - it won't compile |
| Can point into freed memory | Garbage collector keeps pointed-at values alive |
| Casts between unrelated pointer types | Strict typing - `*int` and `*string` are incompatible |
| Manual `malloc`/`free` | Allocation is automatic; no `free` exists |

```go
p := &x
p++ // ❌ compile error: invalid operation: p++ (non-numeric type *int)
```

In Go, a pointer is a *reference to a value*, never a cursor you slide around memory.

---

## The Address-Of and Dereference Operators

Two operators do all the work:

- `&x` - **address-of**: gives you a pointer to `x`
- `*p` - **dereference**: follows the pointer `p` to the value it points at

```go
x := 42

p := &x         // p is a *int holding the address of x
fmt.Println(*p) // 42 - follow the pointer, read the value

*p = 100        // follow the pointer, WRITE the value
fmt.Println(x)  // 100 - x itself changed!
```

### Pointer Types

A pointer's type is `*T` - "pointer to T". The `*` is part of the type name.

```go
var ip *int    // pointer to an int
var sp *string // pointer to a string

fmt.Printf("%T %T\n", ip, sp)
// *int *string
```

`*int` and `*string` are completely different, incompatible types. You can't assign one to the other, and you can't point a `*int` at a float.

### Printing Pointers With %p

`%p` formats an address in hexadecimal. Plain `%v` (or `fmt.Println`) also prints the address for a non-nil pointer.

```go
x := 42
fmt.Printf("%p\n", &x)
// 0x735196f9a020   (the exact address differs on every run - that's normal!)
```

Addresses vary between runs and machines. You'll almost never care what the number *is* - only whether two pointers hold the *same* address.

### Reading the Symbols Aloud

```
p := &x     "p gets the address of x"
v := *p     "v gets the value at p"
*p = 9      "the value at p becomes 9"
var q *int  "q is a pointer to int"
```

Note the double duty: `*` in a **type** (`*int`) means "pointer to int"; `*` in an **expression** (`*p`) means "the value p points at".

---

## The Zero Value of a Pointer: nil

Like every Go type, pointers have a zero value: `nil`, meaning "points at nothing".

```go
var p *int

fmt.Println(p)        // <nil>
fmt.Println(p == nil) // true
```

### Dereferencing nil Panics

Following a pointer that points at nothing is a **runtime panic** - the program crashes:

```go
var p *int
fmt.Println(*p) // 💥 PANIC!
```

Here is the real crash, verified with `go run` on this Go toolchain (the hex numbers and the goroutine trace vary per run):

```
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x2 addr=0x0 pc=0x102b1b600]

goroutine 1 [running]:
main.main()
	/tmp/demo/main.go:7 +0x20
exit status 2
```

This is one of the most common runtime errors in real Go programs. Memorize its first line - when you see `invalid memory address or nil pointer dereference`, you know exactly what happened: something followed a `nil` pointer. Section 9 shows the defensive patterns that prevent it.

---

## Passing by Value vs Passing a Pointer

Level 11 taught that Go passes arguments **by value** - the function gets a copy. That means a function taking `x int` can never change the caller's variable:

```go
func doubleValue(n int) {
    n = n * 2 // modifies the COPY only
}

func doublePointer(n *int) {
    *n = *n * 2 // follows the pointer to the caller's variable
}

func main() {
    x := 10
    doubleValue(x)
    fmt.Println(x) // 10 - unchanged!

    y := 10
    doublePointer(&y)
    fmt.Println(y) // 20 - changed!
}
```

Verified output:

```
10
20
```

### What Actually Happens

```
doubleValue(x)                     doublePointer(&y)

caller:  x = 10                    caller:  y = 10  (at 0xc0000140a0)
            │ copy the VALUE                  │ copy the ADDRESS
            ▼                                 ▼
callee:  n = 10  (separate box)    callee:  n = 0xc0000140a0 ──▶ y
         n = 20  (copy changes,             *n = 20 (writes into the
         caller never sees it)               caller's actual box)
```

Even with a pointer, Go is *still* passing by value - it copies the pointer. But a copy of an address still points at the same place, so writes through it reach the original.

### The Classic Example: swap

```go
func swap(a, b *int) {
    *a, *b = *b, *a
}

x, y := 1, 2
swap(&x, &y)
fmt.Println(x, y) // 2 1
```

Without pointers, `swap` would shuffle its own private copies and accomplish nothing.

---

## Pointers to Structs

You'll see a full structs level next (Level 13), but pointers and structs are such a natural pair that they meet here first.

```go
type Player struct {
    name  string
    score int
}

player := Player{name: "Alice", score: 100}
p := &player // p is a *Player
```

### The Automatic-Dereference Sugar

Strictly speaking, to reach a field through a pointer you'd write `(*p).score`. Go lets you skip the ceremony - `p.score` automatically dereferences:

```go
fmt.Println((*p).score) // 100 - explicit form: works, but noisy
fmt.Println(p.score)    // 100 - automatic dereference: idiomatic

p.score = 150               // writes through the pointer
fmt.Println(player.score)   // 150 - the original struct changed
```

Both forms are identical in meaning. **Always write `p.field`** - you'll essentially never see `(*p).field` in real Go code.

### Mutating Structs Through Functions

This is the big payoff, and the pattern you'll use constantly from here on:

```go
func addPoints(p *Player, points int) {
    p.score += points
}

addPoints(&player, 25)
fmt.Println(player.score) // 175
```

A function taking `Player` (by value) could read the struct but its changes would vanish. Taking `*Player` makes the mutation real. In Level 14 this exact idea becomes *pointer receivers* on methods.

---

## new(T) vs &T{}

Go gives you two ways to allocate a value and get a pointer to it in one step.

### new(T)

The built-in `new(T)` allocates a zero value of type `T` and returns a `*T`:

```go
pi := new(int)   // *int pointing at 0
pp := new(Point) // *Point pointing at {x:0 y:0}

fmt.Println(*pi) // 0
```

`new` can allocate *any* type, but it can only ever give you the **zero value** - there's no way to pass initial data.

### &T{} (the composite-literal form)

For structs (and slices/maps), taking the address of a composite literal does the same allocation but lets you set fields:

```go
origin := &Point{}          // same result as new(Point)
corner := &Point{x: 3, y: 4} // allocated AND initialized

fmt.Printf("%+v\n", *corner) // {x:3 y:4}
```

### Comparison

| | `new(T)` | `&T{}` |
|---|----------|--------|
| Returns | `*T` | `*T` |
| Initial value | always the zero value | any fields you choose |
| Works for | any type (`new(int)` is fine) | composite types (structs, slices, maps, arrays) |
| Idiomatic for structs? | rarely used | ✅ **preferred** |

**The idiomatic choice:** use `&T{}` for structs - it does everything `new(T)` does, can also initialize, and looks like the struct literals you'll write everywhere else. `new` mostly appears when you need a pointer to a zeroed *basic* type (`new(int)`), which is uncommon.

---

## Returning a Pointer to a Local Variable Is Safe

If you've written C, this section will feel illegal. It isn't - in Go, it's a core idiom.

```go
type Counter struct {
    count int
}

func newCounter() *Counter {
    c := Counter{count: 0} // a LOCAL variable...
    return &c              // ...and we return its address. Totally safe!
}
```

In C, `c` would die when the function returns, leaving a dangling pointer. In Go, the compiler performs **escape analysis**: it notices that `&c` outlives the function, so it allocates `c` on the *heap* instead of the stack. The garbage collector then keeps it alive for as long as any pointer refers to it.

You can watch the compiler make this decision. Building the program above with `go build -gcflags=-m` prints (verified on this toolchain):

```
./main.go:10:2: moved to heap: c
```

That line is escape analysis saying: "this local escapes - I'll heap-allocate it."

### The Constructor Pattern

This is why Go code is full of `NewThing()` functions returning pointers:

```go
c1 := newCounter()
c2 := newCounter()

increment(c1)
increment(c1)
increment(c2)

fmt.Println(c1.count) // 2 - each call returned an INDEPENDENT counter
fmt.Println(c2.count) // 1
```

You never think about stack vs heap when writing Go - the compiler decides, and it always decides *correctly*. There is no dangling-pointer bug to write.

---

## When You Do NOT Need Pointers

This is where many learners over-apply their new tool. **Slices and maps are already reference-like** - they contain an internal pointer to their data (Level 8 called it the slice header). Passing them by value still shares the underlying data.

### Mutating Elements: Visible Without Pointers

```go
func setFirst(nums []int) {
    nums[0] = 999
}

func addScore(scores map[string]int, name string, score int) {
    scores[name] = score
}

nums := []int{1, 2, 3}
setFirst(nums)
fmt.Println(nums) // [999 2 3] - caller sees the change!

scores := map[string]int{"Alice": 95}
addScore(scores, "Bob", 87)
fmt.Println(scores) // map[Alice:95 Bob:87] - caller sees the change!
```

No `*[]int`, no `*map[string]int`. The copy of the slice header / map reference still points at the same underlying data, so element writes are shared. Passing `*[]int` or `*map[K]V` around is almost always a design smell.

### 🚨 The append Gotcha: Length Changes Are NOT Shared

Here's the classic confusion. Mutating *elements* is visible - but **appending inside a function is not**:

```go
func addItem(items []int) {
    items = append(items, 42) // reassigns the LOCAL copy of the header
    fmt.Println("inside addItem :", items, "len =", len(items))
}

nums := []int{1, 2, 3}
addItem(nums)
fmt.Println("caller sees    :", nums, "len =", len(nums))
```

Verified output:

```
inside addItem : [1 2 3 42] len = 4
caller sees    : [1 2 3] len = 3
```

Why? A slice is a small header: `{pointer, len, cap}`. The function receives a *copy* of that header. `append` returns a **new header** (bigger `len`, possibly a whole new backing array), and assigning it to `items` only updates the function's local copy. The caller's header still says `len = 3`.

It gets subtler: if the slice happens to have spare capacity, `append` writes the new element into the *shared* backing array - but the caller *still* can't see it, because the caller's `len` never changed. Verified:

```go
spare := make([]int, 3, 4) // len=3 cap=4
spare[0], spare[1], spare[2] = 1, 2, 3
addItem(spare)
fmt.Println("caller sees    :", spare, "len =", len(spare))
fmt.Println("but spare[:4] reveals:", spare[:4])
```

```
inside addItem : [1 2 3 42] len = 4
caller sees    : [1 2 3] len = 3
but spare[:4] reveals: [1 2 3 42]
```

The 42 was silently written into the shared array, invisible behind the caller's length. Depending on capacity, the append is either lost entirely or lurking invisibly - either way, **the caller cannot rely on it**.

### The Fix: Return the New Slice

```go
func addItemFixed(items []int) []int {
    return append(items, 42)
}

nums = addItemFixed(nums) // ✅ caller adopts the new header
```

This is exactly why the standard library's `append` itself is used as `s = append(s, v)` - the reassignment is the contract. (A `*[]int` parameter would also work, but returning the slice is the idiomatic fix.)

### Quick Rules

| You're passing... | Need a pointer to mutate? |
|-------------------|---------------------------|
| `int`, `string`, `bool`, `float64` | ✅ Yes - `*int` etc. |
| a struct | ✅ Yes - `*MyStruct` |
| slice **elements** (`s[i] = v`) | ❌ No - already shared |
| slice **length** (`append`) | ❌ Return the new slice instead |
| map entries (add/update/delete) | ❌ No - already shared |

---

## nil-Safety Patterns

Since dereferencing `nil` panics, defensive Go code checks before following a pointer:

### Pattern 1: Guard Before Dereferencing

```go
if p != nil {
    fmt.Println(*p)
} else {
    fmt.Println("no value")
}
```

### Pattern 2: nil-Safe Accessor With a Default

A helper that *accepts* nil and turns it into a sensible default:

```go
type Config struct {
    timeout int
}

func timeoutOrDefault(cfg *Config) int {
    if cfg == nil {
        return 30 // sensible default
    }
    return cfg.timeout
}

custom := &Config{timeout: 90}
var missing *Config // nil

fmt.Println(timeoutOrDefault(custom))  // 90
fmt.Println(timeoutOrDefault(missing)) // 30 - no panic!
```

Callers never need to nil-check first - the helper absorbs the danger once, in one place.

### Pattern 3: Early Return

```go
func process(p *Player) {
    if p == nil {
        return // or return an error - Level 16 covers that properly
    }
    p.score++
}
```

Panics from `nil` *can* be caught with `recover()` (you'll do it in Exercise 3), but recovery is a last-resort safety net - the real fix is always a check like the ones above.

---

## Best Practices

### 1. Prefer &T{} Over new(T) for Structs

```go
// ✅ Good - idiomatic, and can initialize fields
p := &Point{x: 3, y: 4}

// ❌ Works but unidiomatic, and can't initialize
p := new(Point)
p.x = 3
p.y = 4
```

### 2. Use p.field, Never (*p).field

```go
// ✅ Good - automatic dereference
p.score = 150

// ❌ Legal but needlessly noisy
(*p).score = 150
```

### 3. Don't Pass Pointers to Slices or Maps

```go
// ✅ Good - slices/maps already share their data
func addScore(scores map[string]int, name string, s int)
func setFirst(nums []int)

// ❌ Almost always unnecessary noise
func addScore(scores *map[string]int, name string, s int)
```

(The rare exception: a function that must *replace* the caller's whole slice - but prefer returning the new slice.)

### 4. Return the Slice From Functions That append

```go
// ✅ Good - the append contract
func addItem(items []int, v int) []int {
    return append(items, v)
}
items = addItem(items, 42)
```

### 5. Nil-Check at the Boundary, Not Everywhere

Check pointers where they *enter* your code (function tops, constructors), so the rest of the function can dereference freely. Sprinkling `if p != nil` around every single line is a sign the check belongs earlier.

### 6. Use Pointers for "Must Mutate" or "Too Big to Copy"

If a function only *reads* a small value, take it by value - it's simpler and can't surprise anyone. Reach for `*T` when the function must mutate, or the struct is large enough that copying matters.

---

## Common Mistakes

### Mistake 1: Dereferencing a Possibly-nil Pointer

```go
// ❌ WRONG - panics if p is nil
var p *int
fmt.Println(*p) // panic: runtime error: invalid memory address or nil pointer dereference

// ✅ RIGHT - guard first
if p != nil {
    fmt.Println(*p)
}
```

### Mistake 2: Expecting a Value Parameter to Mutate the Caller

```go
// ❌ WRONG - modifies a copy, caller sees nothing
func levelUp(score int) {
    score += 100
}

// ✅ RIGHT - take a pointer when you mean to mutate
func levelUp(score *int) {
    *score += 100
}
```

### Mistake 3: Appending Inside a Function and Expecting the Caller to See It

```go
// ❌ WRONG - the local header copy grows, the caller's doesn't
func addItem(items []int) {
    items = append(items, 42)
}

// ✅ RIGHT - return the new slice
func addItem(items []int) []int {
    return append(items, 42)
}
```

### Mistake 4: Confusing the Two Meanings of *

```go
var p *int  // * in a TYPE: "p is a pointer to int"
v := *p     // * in an EXPRESSION: "the value p points at"

// ❌ WRONG mental model: "var p *int declares a dereference"
// ✅ RIGHT: type position = pointer type; expression position = dereference
```

### Mistake 5: Taking a Pointer to a Loop Variable's Copy Unintentionally

```go
// ❌ RISKY - v is the loop's own variable, not the slice element
for _, v := range players {
    modify(&v) // mutates the copy v, NOT players[i]!
}

// ✅ RIGHT - point at the element itself
for i := range players {
    modify(&players[i])
}
```

### Mistake 6: Trying Pointer Arithmetic

```go
// ❌ WRONG - this is Go, not C; it won't compile
p := &x
p++          // invalid operation: p++ (non-numeric type *int)
q := p + 1   // invalid operation

// ✅ RIGHT - to walk a collection, index a slice; pointers only reference
```

### Mistake 7: Using new(T) When You Wanted Initial Values

```go
// ❌ WRONG-ish - new can only zero; needs follow-up assignments
p := new(Player)
p.name = "Alice"

// ✅ RIGHT - one clear expression
p := &Player{name: "Alice"}
```

---

## Summary

**The Core Idea:**
- A pointer holds the memory *address* of another value
- `&x` takes the address; `*p` follows it (to read OR write)
- Pointer types are `*T`; the zero value is `nil`; dereferencing `nil` panics
- Go has NO pointer arithmetic - pointers reference, they don't roam

**Mutation Across Functions:**
- Go passes copies; `func f(x int)` cannot change the caller's variable
- `func f(x *int)` can - a copied address still points at the original
- Structs pair naturally with pointers: `p.field` auto-dereferences

**Allocation:**
- `&T{...}` - idiomatic pointer-to-struct, with initial values
- `new(T)` - pointer to a zero value; rarely needed for structs
- Returning `&localVar` is SAFE - escape analysis heap-allocates it

**What Doesn't Need Pointers:**
- Slice element writes and map entries are shared without pointers
- But `append` inside a function is invisible to the caller (header copy) - return the new slice
- Guard with `if p != nil` before dereferencing anything that might be nil

---

## Next Steps

You now understand:
- ✅ What pointers are and why Go has them (without pointer arithmetic)
- ✅ `&`, `*`, pointer types, and `%p`
- ✅ nil pointers and the panic they cause when dereferenced
- ✅ Value vs pointer parameters, and mutation across function boundaries
- ✅ Struct pointers, auto-dereference, `new(T)` vs `&T{}`, and safe constructors
- ✅ When slices and maps make pointers unnecessary - and the append gotcha

**Next level:** Level 13 - Structs
- Defining your own types with `struct`
- Struct literals, field access, and nested structs
- Structs combined with pointers (you just learned these!)
- Comparing, copying, and embedding structs

Pointers are the hardest single idea in this tier - and you've got them. Everything from here builds on today. Keep going! 🚀
