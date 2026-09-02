# Level 15: Interfaces - Complete Guide

## Introduction

Welcome to Level 15! You've mastered methods (Level 14) - attaching behavior to your own types with value and pointer receivers. Now it's time for **interfaces** - the feature that ties methods, types, and polymorphism together, and arguably the single most important concept for writing idiomatic Go.

An interface describes **what a type can do**, not what it is. Any type that has the right methods satisfies an interface *automatically* - no `implements` keyword, no inheritance, no ceremony. This one design decision shapes how the entire Go ecosystem is built: small contracts, loosely coupled packages, and functions that work with any type that behaves correctly.

This is the fourth level of the **OOP CONCEPTS** tier (Levels 12-16). Everything you learned about method sets in Level 14 pays off here - and everything here sets up Level 16, because Go's error handling is built on one tiny interface called `error`.

---

## Table of Contents

1. [What Is an Interface?](#what-is-an-interface)
2. [Implicit Satisfaction](#implicit-satisfaction)
3. [Polymorphism: One Interface, Many Types](#polymorphism-one-interface-many-types)
4. [Interface Values Under the Hood](#interface-values-under-the-hood)
5. [fmt.Stringer](#fmtstringer)
6. [The Empty Interface: any](#the-empty-interface-any)
7. [Type Assertions](#type-assertions)
8. [Type Switches](#type-switches)
9. [Method Sets and Interface Satisfaction](#method-sets-and-interface-satisfaction)
10. [The nil Interface vs Typed-nil Trap](#the-nil-interface-vs-typed-nil-trap)
11. [Interface Composition](#interface-composition)
12. [Interface Design: Small Is Beautiful](#interface-design-small-is-beautiful)
13. [Best Practices](#best-practices)
14. [Common Mistakes](#common-mistakes)

---

## What Is an Interface?

An interface is a **named set of method signatures** - a contract that describes *behavior*, not data.

```go
type Speaker interface {
    Speak() string
}
```

This says: "a Speaker is anything with a `Speak()` method that returns a `string`." That's the whole contract.

Compare the two kinds of type definitions you now know:

```go
// A struct describes what a thing IS (data):
type Dog struct {
    Name string
    Age  int
}

// An interface describes what a thing can DO (behavior):
type Speaker interface {
    Speak() string
}
```

Three things to notice about interface declarations:

- **No fields.** Interfaces carry zero data - only method signatures.
- **No method bodies.** The interface says *what* must exist, never *how* it works.
- **No list of implementing types.** The interface doesn't know (or care) who satisfies it.

💡 A helpful mental model: a struct is a *noun*, an interface is a *job description*. Any noun that can do the job qualifies.

---

## Implicit Satisfaction

Here is the part that surprises people coming from Java or C#: **there is no `implements` keyword in Go.** A type satisfies an interface simply by having all of its methods.

```go
type Speaker interface {
    Speak() string
}

type Dog struct{ Name string }

func (d Dog) Speak() string { return d.Name + " says Woof!" }

type Robot struct{ ID int }

func (r Robot) Speak() string { return fmt.Sprintf("Robot-%d beeps", r.ID) }
```

`Dog` and `Robot` never mention `Speaker` anywhere. But because each has a `Speak() string` method, **both satisfy `Speaker` automatically**:

```go
speakers := []Speaker{Dog{Name: "Rex"}, Robot{ID: 7}}
for _, s := range speakers {
    fmt.Println(s.Speak())
}
// Rex says Woof!
// Robot-7 beeps
```

### Why Implicit Satisfaction Matters

This is not just saved typing - it **decouples packages**:

- A type can satisfy an interface defined in a package it has *never imported* - even one written years later.
- The consumer of behavior can define exactly the interface it needs, without asking type authors to opt in.
- Your `Dog` type satisfies `fmt.Stringer`, `Speaker`, and any future one-method interface with the right signature - all at once, with zero declarations.

```
Java/C#:  type  --"implements X"-->  interface   (type must know the interface)
Go:       type  --has the methods--> interface   (neither side knows the other)
```

This is why Go programmers say interfaces are *discovered*, not *designed up front*: you notice several types share behavior, and you name that behavior.

---

## Polymorphism: One Interface, Many Types

Here is the core demonstration of this level - one function that works with every shape ever written, past or future:

```go
package main

import (
    "fmt"
    "math"
)

type Shape interface {
    Area() float64
    Perimeter() float64
}

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}

type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// describe accepts ANY Shape - Circle, Rectangle, or types not yet written
func describe(s Shape) {
    fmt.Printf("%T: area=%.2f, perimeter=%.2f\n", s, s.Area(), s.Perimeter())
}

func main() {
    describe(Circle{Radius: 2})
    describe(Rectangle{Width: 3, Height: 4})
}
```

Output:

```
main.Circle: area=12.57, perimeter=12.57
main.Rectangle: area=12.00, perimeter=14.00
```

`describe` was written once, against the `Shape` contract. It has no idea circles or rectangles exist - it only knows "whatever I'm given can compute its area and perimeter." Add a `Triangle` next month and `describe` handles it without changing a single line. **That is polymorphism.**

Interfaces also work as element types, so you can mix concrete types in one slice:

```go
shapes := []Shape{
    Circle{Radius: 1},
    Rectangle{Width: 2, Height: 5},
    Circle{Radius: 3},
}
totalArea := 0.0
for _, s := range shapes {
    totalArea += s.Area() // calls the right method for each dynamic type
}
fmt.Printf("total area of %d shapes: %.2f\n", len(shapes), totalArea)
// total area of 3 shapes: 41.42
```

---

## Interface Values Under the Hood

This model explains *everything else in this level*, so slow down here.

An interface value is a **pair**: `(dynamic type, dynamic value)`.

```
var s Shape                 s = Circle{Radius: 2}

┌──────────────┐            ┌──────────────────┐
│ type:  nil   │            │ type:  Circle    │
│ value: nil   │            │ value: {2}       │
└──────────────┘            └──────────────────┘
  a nil interface             holds a Circle
```

```go
var s Shape // nil interface: dynamic type=nil, dynamic value=nil
fmt.Printf("type=%T value=%v s==nil? %v\n", s, s, s == nil)

s = Circle{Radius: 2} // now: dynamic type=main.Circle, dynamic value={2}
fmt.Printf("type=%T value=%v s==nil? %v\n", s, s, s == nil)
```

Output:

```
type=<nil> value=<nil> s==nil? true
type=main.Circle value={2} s==nil? false
```

Key consequences of the pair model:

- **Method calls dispatch on the dynamic type.** `s.Area()` looks at the pair, sees `Circle`, and runs `Circle`'s `Area` method. That's how one call site runs different code for different types.
- **`%T` prints the dynamic type** - which is why `describe` printed `main.Circle` and `main.Rectangle` above.
- **An interface is `== nil` only when BOTH halves are nil.** This single fact produces the typed-nil trap in section 10.
- **Calling a method on a nil interface panics** - there's no dynamic type to dispatch to:

```go
var s Shape
s.Area()
// panic: runtime error: invalid memory address or nil pointer dereference
```

---

## fmt.Stringer

You've been benefiting from interfaces since Level 1 without knowing it. The `fmt` package defines this tiny interface:

```go
type Stringer interface {
    String() string
}
```

Whenever `fmt.Println` (or `%v`, or `fmt.Sprintf`, ...) formats a value, it first checks: *does this value satisfy `Stringer`?* If yes, it calls `String()` and uses the result. Level 14 previewed the `String()` method - now you know the machinery behind it.

```go
type Book struct {
    Title string
    Pages int
}

func (b Book) String() string {
    return fmt.Sprintf("%q (%d pages)", b.Title, b.Pages)
}

func main() {
    b := Book{"The Go Programming Language", 380}
    fmt.Println(b)
    fmt.Printf("reading %v\n", b)
}
```

Output:

```
"The Go Programming Language" (380 pages)
reading "The Go Programming Language" (380 pages)
```

Without the `String()` method, the same `Println` would print `{The Go Programming Language 380}` - the default struct formatting.

📖 Notice what just happened: **your** type satisfied an interface defined in the **standard library**, and `fmt` called **your** code - even though `fmt` was written long before your `Book` existed. That is implicit satisfaction doing real work.

---

## The Empty Interface: any

An interface with zero methods is satisfied by *every* type - there are no requirements to fail.

```go
any             // the empty interface (alias for interface{}, since Go 1.18)
```

`any` can hold a value of any type:

```go
var x any
x = 42
fmt.Printf("value=%v type=%T\n", x, x) // value=42 type=int
x = "hello"
fmt.Printf("value=%v type=%T\n", x, x) // value=hello type=string
x = []int{1, 2, 3}
fmt.Printf("value=%v type=%T\n", x, x) // value=[1 2 3] type=[]int
```

This is how `fmt.Println` accepts anything - its real signature is:

```go
func Println(a ...any) (n int, err error)
```

### The Cost: You Lose Type Safety

Once a value goes into an `any`, the compiler only knows it's "something":

```go
var n any = 10
total := n + 5 // ❌ compile error - the static type of n is any, not int
```

To do anything type-specific, you must get the concrete type back out with a type assertion or type switch (next two sections). Every assertion is a runtime check that can fail - the compiler can no longer help you.

🚨 **`any` is appropriate only rarely**: truly generic containers of mixed values, decoding unknown JSON, print-style functions. If you catch yourself reaching for `any` for ordinary code, you usually want a *specific* small interface instead - and since Go 1.18, generics (Level 25) cover many former uses of `any` with full type safety.

---

## Type Assertions

A **type assertion** extracts the dynamic value from an interface, asserting it has a specific type.

### The comma-ok Form (safe - use this)

```go
var i any = "hello"

s, ok := i.(string)
fmt.Printf("s=%q ok=%v\n", s, ok)
// s="hello" ok=true

n, ok := i.(int)
fmt.Printf("n=%d ok=%v\n", n, ok)
// n=0 ok=false   (failed assertion gives the zero value, no panic)
```

If the assertion fails, `ok` is `false` and the value is the zero value of the asserted type. This never panics, so it's perfect for guarding behavior:

```go
if s, ok := v.(string); ok {
    fmt.Printf("%q is a string of length %d\n", s, len(s))
}
```

### The Single-Value Form (panics on failure)

```go
var i any = "hello"
s := i.(string) // ✅ works - i really does hold a string
n := i.(int)    // ❌ PANICS - i holds a string, not an int
```

The failed assertion crashes the program with this exact message:

```
panic: interface conversion: interface {} is string, not int
```

Use the single-value form only when a failed assertion is a programmer error that *should* crash. In normal control flow, always use comma-ok.

---

## Type Switches

When you need to handle *several* possible dynamic types, chaining assertions gets clumsy. A **type switch** branches on the dynamic type directly:

```go
func describe(i any) string {
    switch v := i.(type) {
    case nil:
        return "nil value"
    case int:
        return fmt.Sprintf("int %d (doubled: %d)", v, v*2)
    case string:
        return fmt.Sprintf("string %q (%d bytes)", v, len(v))
    case bool:
        return fmt.Sprintf("bool %v", v)
    case []int:
        sum := 0
        for _, n := range v {
            sum += n
        }
        return fmt.Sprintf("[]int with %d elements (sum: %d)", len(v), sum)
    default:
        return fmt.Sprintf("unhandled type %T", v)
    }
}
```

```go
values := []any{42, "gopher", true, []int{1, 2, 3}, 3.14, nil}
for _, v := range values {
    fmt.Println(describe(v))
}
```

Output:

```
int 42 (doubled: 84)
string "gopher" (6 bytes)
bool true
[]int with 3 elements (sum: 6)
unhandled type float64
nil value
```

Notes:

- `switch v := i.(type)` is special syntax - `.(type)` is only legal inside a `switch`.
- **Inside each `case`, `v` already has that case's concrete type** - no extra assertion needed (`v*2` works in the `int` case, `len(v)` in the `string` case).
- `case nil` matches a nil interface; `default` catches everything unlisted.

---

## Method Sets and Interface Satisfaction

Level 14's method-set rules become critically important with interfaces. Recall:

| Type | Method set contains |
|------|--------------------|
| `T` (value) | methods with **value receivers** only |
| `*T` (pointer) | methods with value receivers **AND** pointer receivers |

**A type satisfies an interface only if its method set contains ALL of the interface's methods.** So:

```go
type Incrementer interface {
    Increment()
    Value() int
}

type Counter struct {
    count int
}

func (c *Counter) Increment() { // pointer receiver
    c.count++
}

func (c Counter) Value() int { // value receiver
    return c.count
}
```

A `*Counter` satisfies `Incrementer` (its method set has both methods):

```go
c := Counter{}
var inc Incrementer = &c // ✅ works
inc.Increment()
```

But a plain `Counter` value does NOT - `Increment` is missing from its method set:

```go
c := Counter{}
var inc Incrementer = c // ❌ compile error!
```

The compiler rejects it with this exact message:

```
cannot use c (variable of struct type Counter) as Incrementer value in variable declaration: Counter does not implement Incrementer (method Increment has pointer receiver)
```

That last parenthetical - `(method Increment has pointer receiver)` - is the compiler telling you exactly what's wrong and hinting at the fix: **assign a pointer** (`&c`) instead.

The other direction never bites: if all methods use value receivers, both `T` and `*T` satisfy the interface:

```go
type Greeter interface{ Greet() string }

type Person struct{ Name string }

func (p Person) Greet() string { return "hi, " + p.Name } // value receiver

var g1 Greeter = Person{Name: "Ann"}  // ✅ value works
var g2 Greeter = &Person{Name: "Bob"} // ✅ pointer also works
```

💡 Rule of thumb: if *any* method of a type has a pointer receiver, store pointers (`&T{...}`) in interface values.

---

## The nil Interface vs Typed-nil Trap

This is the most famous interface gotcha in Go, and the pair model from section 4 explains it completely.

**An interface value is `== nil` only when BOTH its dynamic type AND dynamic value are nil.** An interface holding a nil *pointer* still has a dynamic type - so it is **not** nil:

```go
type MyErr struct{}

func (e *MyErr) Error() string { return "boom" }

var p *MyErr            // p is a nil *MyErr
var err error = p       // interface now holds (type=*MyErr, value=nil)
fmt.Println(p == nil)   // true
fmt.Println(err == nil) // false!  ← the trap
fmt.Printf("%T\n", err) // *main.MyErr
```

```
       nil interface                typed nil (NOT a nil interface!)
     ┌──────────────┐               ┌────────────────────┐
     │ type:  nil   │               │ type:  *MyErr      │
     │ value: nil   │               │ value: nil         │
     └──────────────┘               └────────────────────┘
       err == nil → true              err == nil → false
```

### The Classic Error-Return Trap

This bites hardest when a function declares a *concrete pointer* return type and its caller stores the result in an `error` interface:

```go
type ValidationError struct {
    Field string
}

func (e *ValidationError) Error() string {
    return "invalid field: " + e.Field
}

// ❌ BUGGY: declares the concrete pointer type as its return type
func validateBad(input string) *ValidationError {
    if input == "" {
        return &ValidationError{Field: "input"}
    }
    return nil // a nil *ValidationError - NOT a nil interface!
}

func main() {
    var err error = validateBad("valid input") // validation passed...
    fmt.Println("err == nil?", err == nil)
    if err != nil {
        fmt.Println("BUG: this branch runs even though validation passed!")
    }
}
```

Output (verified):

```
err == nil? false
BUG: this branch runs even though validation passed!
```

The `nil` returned by `validateBad` is a typed `(*ValidationError)(nil)`. Assigning it to `err` wraps it in an interface with a non-nil dynamic type - so `err != nil` is true and the error path runs on *success*.

### The Fix

Declare the **interface type** (`error`) as the return type, and return literal `nil`:

```go
// ✅ FIXED: declares the error interface as its return type
func validateGood(input string) error {
    if input == "" {
        return &ValidationError{Field: "input"}
    }
    return nil // a true nil interface: (type=nil, value=nil)
}
```

Now `validateGood("valid input") == nil` is `true`, and `err != nil` means what everyone expects. This is exactly why idiomatic Go functions return `error`, never `*SomeError` - a rule you'll lean on constantly in Level 16.

---

## Interface Composition

Interfaces can **embed other interfaces**, combining small contracts into bigger ones:

```go
type Reader interface {
    Read() string
}

type Writer interface {
    Write(data string)
}

// Composition: ReadWriter is Reader AND Writer combined
type ReadWriter interface {
    Reader
    Writer
}
```

A type satisfies `ReadWriter` by having *all* the methods of both embedded interfaces:

```go
type MemoryBuffer struct {
    data string
}

func (m *MemoryBuffer) Read() string       { return m.data }
func (m *MemoryBuffer) Write(data string)  { m.data += data }

buf := &MemoryBuffer{}
var w Writer = buf      // ✅
var r Reader = buf      // ✅
var rw ReadWriter = buf // ✅ satisfies the composed interface too
```

And because a `ReadWriter` *is* a `Reader` and *is* a `Writer`, it can be passed wherever either is needed.

### The Real-World Model: package io

This pattern is lifted straight from the standard library, where it's used everywhere:

```go
// from the io package:
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type ReadWriter interface {
    Reader
    Writer
}
```

Files, network connections, buffers, gzip streams, HTTP bodies - all of them satisfy `io.Reader` and/or `io.Writer`. That's why one `io.Copy(dst, src)` function can copy file→network, network→buffer, or anything→anything. Small interfaces, composed as needed, are the backbone of Go's standard library.

---

## Interface Design: Small Is Beautiful

Go's most quoted design proverb (Rob Pike):

> **"The bigger the interface, the weaker the abstraction."**

A one-method interface (`io.Reader`, `fmt.Stringer`, `error`) is satisfied by countless types and composes freely. A ten-method interface is satisfied by almost nothing and locks implementations in.

### Rule 1: Keep Interfaces Small

Most good Go interfaces have **one or two methods**. Notice the standard library's `-er` naming convention: `Reader`, `Writer`, `Stringer`, `Closer` - a small interface names one *capability*.

### Rule 2: Accept Interfaces, Return Concrete Types

```go
// ✅ Parameter: interface - callers can pass anything that fits
func report(w io.Writer, data []Record) error

// ✅ Return: concrete type - callers get full access to everything it offers
func NewBuffer() *bytes.Buffer
```

Accepting an interface makes your function flexible and testable (pass a real file in production, an in-memory buffer in tests). Returning a concrete type keeps you honest and lets callers decide which interface *they* care about.

### Rule 3: Define Interfaces Where They're USED, Not Where Types Are Defined

This one is the least obvious and the most Go-specific. Because satisfaction is implicit, the **consumer** can declare exactly the contract it needs:

```go
// ❌ Provider-side (Java style): storage package exports an interface
//    and hopes consumers want exactly these 8 methods

// ✅ Consumer-side (Go style): the report code declares what IT needs
type UserSource interface {
    FetchUsers() []string
}

func printReport(src UserSource) { ... }
```

Any type with a `FetchUsers() []string` method - a database wrapper, a CSV file reader, a test fake - works immediately, without importing the report package or declaring anything. The provider ships concrete types; each consumer defines its own minimal view of them.

### Rule 4: Don't Create Interfaces Preemptively

Start with concrete types. Introduce an interface only when you actually have (or genuinely foresee) a second implementation or a consumer that needs decoupling. An interface with exactly one implementation and one consumer is usually just indirection.

---

## Best Practices

### 1. Keep Interfaces Small (1-2 Methods)

```go
// ✅ Good - one capability, satisfied by many types
type Notifier interface {
    Notify(message string)
}

// ❌ Weak - a grab-bag almost nothing will satisfy
type MegaService interface {
    Notify(message string)
    Connect() error
    Close() error
    Stats() map[string]int
    Reset()
}
```

### 2. Accept Interfaces, Return Concrete Types

```go
// ✅ Good
func Fprint(w io.Writer, v any) error   // accepts an interface
func NewBuffer() *bytes.Buffer          // returns a concrete type
```

### 3. Define Interfaces at the Consumer

Declare the interface in the code that *calls* the methods, listing only the methods it actually calls. Providers export concrete types.

### 4. Return error, Never a Concrete Error Pointer

```go
// ✅ Good - a nil return is a true nil interface
func validate(input string) error

// ❌ Trap - callers comparing the result to nil will be wrong
func validate(input string) *ValidationError
```

### 5. Prefer comma-ok Assertions and Type Switches

```go
// ✅ Good - handles failure gracefully
if s, ok := v.(string); ok { ... }

// ❌ Risky - panics unless you KNOW the dynamic type
s := v.(string)
```

### 6. Store Pointers in Interfaces When Any Method Has a Pointer Receiver

```go
c := Counter{}
var inc Incrementer = &c // ✅ works regardless of receiver kinds
```

### 7. Follow the -er Naming Convention

Single-method interfaces are named after the method: `Read` → `Reader`, `Notify` → `Notifier`, `String` → `Stringer`.

---

## Common Mistakes

### Mistake 1: Assigning a Value When a Pointer Receiver Is Needed

```go
type Incrementer interface{ Increment() }
type Counter struct{ count int }
func (c *Counter) Increment() { c.count++ }

// ❌ WRONG - Counter's method set doesn't include Increment
var inc Incrementer = Counter{}
// compile error: Counter does not implement Incrementer
//                (method Increment has pointer receiver)

// ✅ RIGHT
var inc Incrementer = &Counter{}
```

### Mistake 2: Returning a Concrete Error Pointer (Typed nil)

```go
// ❌ WRONG - "return nil" here produces a NON-nil error interface in callers
func validate(s string) *ValidationError {
    return nil
}

// ✅ RIGHT - return the interface type; nil means nil
func validate(s string) error {
    return nil
}
```

### Mistake 3: Single-Value Assertion Without Certainty

```go
var i any = "hello"

// ❌ WRONG - panics: interface conversion: interface {} is string, not int
n := i.(int)

// ✅ RIGHT
n, ok := i.(int)
if !ok {
    // handle the mismatch
}
```

### Mistake 4: Reaching for any Instead of a Specific Interface

```go
// ❌ WRONG - caller can pass anything, function must type-switch on everything
func process(item any)

// ✅ RIGHT - state the behavior you actually need
func process(item Notifier)
```

### Mistake 5: Giant Provider-Side Interfaces

```go
// ❌ WRONG - provider exports a 10-method interface "just in case"
// ✅ RIGHT - provider exports concrete types;
//            each consumer defines the 1-2 method interface it needs
```

### Mistake 6: Calling a Method on a nil Interface

```go
var s Shape // nil interface - no dynamic type to dispatch to
s.Area()
// ❌ panic: runtime error: invalid memory address or nil pointer dereference

// ✅ RIGHT - assign a concrete value first, or guard with s != nil
```

### Mistake 7: Expecting an Interface Holding a nil Pointer to Equal nil

```go
var p *MyErr
var err error = p
// ❌ WRONG ASSUMPTION
fmt.Println(err == nil) // false! (dynamic type is *MyErr, so not nil)
```

---

## Summary

**The Contract:**
- An interface is a named set of method signatures - behavior, not data
- Satisfaction is implicit: any type with the right methods satisfies it, no `implements` keyword
- Implicit satisfaction decouples packages - consumers define contracts, providers ship types

**The Machinery:**
- An interface value is a `(dynamic type, dynamic value)` pair
- Method calls dispatch on the dynamic type - that's polymorphism
- An interface is `== nil` only when BOTH halves are nil - an interface holding a typed nil pointer is NOT nil
- Method sets decide satisfaction: `T` offers value-receiver methods; `*T` offers both

**The Tools:**
- `v, ok := i.(T)` - safe assertion; single-value form panics on mismatch
- `switch v := i.(type)` - branch on the dynamic type, `v` is concrete in each case
- `any` (= `interface{}`) holds anything but forfeits type safety - use sparingly
- Embed interfaces in interfaces to compose contracts (`Reader` + `Writer` = `ReadWriter`)

**The Design Rules:**
- The bigger the interface, the weaker the abstraction - keep them small
- Accept interfaces, return concrete types
- Define interfaces where they're used, not where types live

---

## Next Steps

You now understand:
- ✅ Interfaces as behavior contracts and implicit satisfaction
- ✅ Polymorphism through interface values and dynamic dispatch
- ✅ The (dynamic type, dynamic value) pair model
- ✅ fmt.Stringer, `any`, type assertions, and type switches
- ✅ Method sets, the pointer-receiver compile error, and the typed-nil trap
- ✅ Interface composition and small-interface design

**Next level:** Level 16 - Error Handling
- The `error` interface - one method, `Error() string` (you can already read it!)
- Creating, wrapping, and unwrapping errors
- `errors.Is`, `errors.As` (which uses type assertions under the hood!)
- Sentinel errors and custom error types - where the typed-nil rule pays off

Interfaces are the heart of idiomatic Go - and Level 16's entire error system is just one small interface used well. You're ready for it! 🚀
