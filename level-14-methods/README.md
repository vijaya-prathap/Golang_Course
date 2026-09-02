# Level 14: Methods - Complete Guide

## Introduction

Welcome to Level 14! You've mastered structs (Level 13). Now it's time to learn **methods** - how to attach behavior directly to your types.

You actually met methods briefly back in Level 3 (`type Age int` with an `IsAdult()` method). This level formalizes that idea: what a receiver really is, when it copies and when it doesn't, and the rules that will power **interfaces** in Level 15.

---

## Table of Contents

1. [What Is a Method?](#what-is-a-method)
2. [Declaring Methods on Structs](#declaring-methods-on-structs)
3. [Value Receivers: The Method Gets a Copy](#value-receivers-the-method-gets-a-copy)
4. [Pointer Receivers: The Method Gets the Original](#pointer-receivers-the-method-gets-the-original)
5. [Choosing Between Value and Pointer Receivers](#choosing-between-value-and-pointer-receivers)
6. [Auto-Address and Auto-Dereference](#auto-address-and-auto-dereference)
7. [Method Sets (Preview for Interfaces)](#method-sets-preview-for-interfaces)
8. [Methods on Named Non-Struct Types](#methods-on-named-non-struct-types)
9. [Methods and Embedding: Method Promotion](#methods-and-embedding-method-promotion)
10. [Method Values and Method Expressions](#method-values-and-method-expressions)
11. [Best Practices](#best-practices)
12. [Common Mistakes](#common-mistakes)

---

## What Is a Method?

A **method** is a function with a **receiver** - an extra parameter written before the function name that binds the function to a type.

```go
type Rectangle struct {
    Width  float64
    Height float64
}

// A METHOD: a function with a receiver (r Rectangle)
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// A plain FUNCTION doing the same job
func area(r Rectangle) float64 {
    return r.Width * r.Height
}
```

Both compute the same thing - the difference is call syntax and organization:

```go
rect := Rectangle{Width: 4, Height: 3}

fmt.Println(rect.Area()) // method call:   12
fmt.Println(area(rect))  // function call: 12
```

**Methods and functions have the same power.** Methods win on organization:

- Behavior lives next to the data it belongs to (`rect.Area()` reads naturally)
- Every type gets its own namespace (`Rectangle` and `Circle` can BOTH have `Area()`)
- Methods make a type satisfy **interfaces** (Level 15 - this is the big payoff)

---

## Declaring Methods on Structs

### Anatomy of a Method

```go
func (r Rectangle) Area() float64 {
     └─────┬─────┘ └─┬─┘  └──┬──┘
        receiver    name   return type

receiver: binds this function to the type Rectangle
          "r" is a real parameter - a Rectangle you can use inside the body
```

### Receiver Naming Conventions

The receiver name should be **short** (one or two letters, usually the first letter of the type) and **consistent** across all of the type's methods.

```go
// ✅ Good - short and consistent
func (r Rectangle) Area() float64      { return r.Width * r.Height }
func (r Rectangle) Perimeter() float64 { return 2 * (r.Width + r.Height) }

// ❌ Not idiomatic - Go does NOT use this/self
func (this Rectangle) Area() float64 { return this.Width * this.Height }
func (self Rectangle) Area() float64 { return self.Width * self.Height }
```

### A Complete Example

```go
type Person struct {
    Name string
}

func (p Person) Greet() string {
    return "Hello, I'm " + p.Name
}

p := Person{Name: "Alice"}
fmt.Println(p.Greet()) // Hello, I'm Alice
```

### One Rule: Same Package Only

You can only define methods on types declared in the **same package**. You can't attach methods to built-in types directly:

```go
// ❌ Compile error
func (s string) Shout() string { return s + "!" }
// cannot define new methods on non-local type string

// ✅ Define a named type first (see section 8)
type Loud string
func (s Loud) Shout() string { return string(s) + "!" }
```

---

## Value Receivers: The Method Gets a Copy

With a **value receiver** like `(p Person)`, the method operates on a **COPY** of the caller's value. Mutations inside the method do NOT affect the original.

```go
type Person struct {
    Name string
    Age  int
}

// Value receiver: p is a copy
func (p Person) BirthdayValue() {
    p.Age++ // increments the COPY only!
}

alice := Person{Name: "Alice", Age: 30}
alice.BirthdayValue()
fmt.Println(alice.Age) // 30 - unchanged!
```

### Proving It

The copy really does change inside the method - the change is just thrown away when the method returns:

```go
type Counter struct {
    Count int
}

func (c Counter) Increment() {
    c.Count++
    fmt.Println("  inside Increment: c.Count =", c.Count)
}

c := Counter{Count: 0}
fmt.Println("before:", c.Count)
c.Increment()
fmt.Println("after: ", c.Count)
```

Verified output:

```
before: 0
  inside Increment: c.Count = 1
after:  0
```

Call it a hundred times - the caller's `Count` stays 0. Each call copies a fresh `Counter`, increments the copy, and discards it.

---

## Pointer Receivers: The Method Gets the Original

With a **pointer receiver** like `(p *Person)`, the method receives a pointer to the caller's value - so mutations affect the **ORIGINAL**.

```go
// Pointer receiver: p points at the original
func (p *Person) Birthday() {
    p.Age++ // increments the ORIGINAL
}

bob := Person{Name: "Bob", Age: 30}
bob.Birthday()
fmt.Println(bob.Age) // 31 - changed!
```

### Side by Side

```go
type Player struct {
    Name  string
    Score int
}

func (p Player) AddPointsValue(n int)    { p.Score += n } // copy
func (p *Player) AddPointsPointer(n int) { p.Score += n } // original

a := Player{Name: "ValueGal", Score: 0}
b := Player{Name: "PointerPal", Score: 0}
for i := 0; i < 5; i++ {
    a.AddPointsValue(10)
    b.AddPointsPointer(10)
}
fmt.Println(a.Name, "final score:", a.Score)
fmt.Println(b.Name, "final score:", b.Score)
```

Verified output:

```
ValueGal final score: 0
PointerPal final score: 50
```

Five calls through the value receiver accomplished nothing. Five calls through the pointer receiver accumulated all 50 points.

---

## Choosing Between Value and Pointer Receivers

### Use a POINTER receiver when:

1. **The method needs to mutate the receiver** (the most common reason)
2. **The struct is large** - a value receiver copies the whole struct on every call
3. **Any other method of the type already uses a pointer receiver** (the consistency rule below)

### Use a VALUE receiver when:

1. The type is **small and value-like** (a point, a color, a temperature, `time.Time`)
2. The method **only reads** and the type has no pointer-receiver methods
3. You WANT copy semantics (immutability by construction)

### The Consistency Rule 🚨

> **If ANY method of a type needs a pointer receiver, give ALL of that type's methods pointer receivers.**

Mixing receivers on one type causes two problems:

- Readers can no longer tell at a glance whether a call can mutate
- The type's value and pointer **method sets** differ in confusing ways (section 7), which bites hard when interfaces arrive in Level 15

```go
// ❌ Mixed - Deposit mutates, Withdraw silently doesn't (see Common Mistakes)
func (b *BankAccount) Deposit(amount float64)  { b.Balance += amount }
func (b BankAccount) Withdraw(amount float64)  { b.Balance -= amount }

// ✅ Consistent - all pointer receivers
func (b *BankAccount) Deposit(amount float64)  { b.Balance += amount }
func (b *BankAccount) Withdraw(amount float64) { b.Balance -= amount }
func (b *BankAccount) Report() string          { return fmt.Sprintf("%s has $%.2f", b.Owner, b.Balance) }
```

---

## Auto-Address and Auto-Dereference

You may have noticed something odd above: `bob.Birthday()` worked even though `bob` is a **value** and `Birthday` wants a **pointer**. That's Go being helpful.

### Calling a Pointer Method on an Addressable Value ✅

If a value is **addressable** (Go can take its address), `v.PointerMethod()` is sugar for `(&v).PointerMethod()`:

```go
type Account struct {
    Owner   string
    Balance float64
}

func (a *Account) Deposit(amount float64) { a.Balance += amount }
func (a Account) Report() string          { return fmt.Sprintf("%s: $%.2f", a.Owner, a.Balance) }

acct := Account{Owner: "Alice", Balance: 100}
acct.Deposit(50)    // sugar: Go rewrites this as (&acct).Deposit(50)
(&acct).Deposit(25) // the explicit long form (never written in practice)
fmt.Println(acct.Report()) // Alice: $175.00
```

### Calling a Value Method Through a Pointer ✅

The reverse works too - `ptr.ValueMethod()` is sugar for `(*ptr).ValueMethod()`:

```go
ptr := &Account{Owner: "Bob", Balance: 10}
fmt.Println(ptr.Report())    // sugar: Go rewrites this as (*ptr).Report()
fmt.Println((*ptr).Report()) // the explicit long form
// both print: Bob: $10.00
```

### The Case That Does NOT Work ❌

The auto-address trick needs an **addressable** value. Two everyday values are NOT addressable:

**1. A struct value stored in a map:**

```go
func (p *Player) AddPoints(n int) { p.Score += n }

players := map[string]Player{
    "alice": {Name: "Alice", Score: 10},
}

players["alice"].AddPoints(5) // ❌ COMPILE ERROR
```

The real compiler output:

```
./main.go:19:22: cannot call pointer method AddPoints on Player
```

**2. A function's return value used directly:**

```go
func makeCounter() Counter { return Counter{} }

makeCounter().Increment() // ❌ COMPILE ERROR
```

```
./main.go:16:19: cannot call pointer method Increment on Counter
```

Go can't take the address of a map element (map internals may move entries around) or of a temporary return value - so there's no pointer to hand the method.

### The Fixes

```go
// Fix 1: copy out, modify, store back
p := players["alice"] // p is an addressable variable
p.AddPoints(5)
players["alice"] = p

// Fix 2: store pointers in the map
roster := map[string]*Player{"bob": {Name: "Bob", Score: 20}}
roster["bob"].AddPoints(5) // works: the map value IS a pointer
```

---

## Method Sets (Preview for Interfaces)

Every type has a **method set** - the set of methods you can rely on for interface satisfaction. The rule:

| Type | Method set includes |
|------|---------------------|
| `T` (value) | methods with **value** receivers only |
| `*T` (pointer) | methods with **value** receivers AND **pointer** receivers |

In other words: `*T`'s method set is always the bigger one.

### Why This Will Matter in Level 15

Interfaces are satisfied by method sets, and there's NO auto-address sugar there - an interface value stores a copy, so Go has no original to point at:

```go
type Incrementer interface {
    Increment()
}

type Counter struct {
    n int
}

func (c *Counter) Increment() { c.n++ } // pointer receiver

var c Counter
var i Incrementer = c // ❌ COMPILE ERROR
```

The real compiler output:

```
./main.go:19:25: cannot use c (variable of struct type Counter) as Incrementer value in variable declaration: Counter does not implement Incrementer (method Increment has pointer receiver)
```

The fix is to store the pointer - `*Counter`'s method set includes `Increment`:

```go
var i Incrementer = &c // ✅ works
i.Increment()
i.Increment()
fmt.Println(c.n) // 2
```

Don't worry about mastering interfaces yet - just remember: **value-receiver methods belong to both `T` and `*T`; pointer-receiver methods belong to `*T` only.** This is the single biggest reason the consistency rule exists.

---

## Methods on Named Non-Struct Types

Receivers aren't limited to structs - ANY named type declared in your package can have methods. Level 3 previewed this with `type Age int`; here's the pattern in full:

### Temperature Conversions

```go
type Celsius float64
type Fahrenheit float64

func (c Celsius) ToFahrenheit() Fahrenheit {
    return Fahrenheit(c*9/5 + 32)
}

func (c Celsius) IsFreezing() bool {
    return c <= 0
}

func (f Fahrenheit) ToCelsius() Celsius {
    return Celsius((f - 32) * 5 / 9)
}

boiling := Celsius(100)
fmt.Printf("%.0f°C = %.0f°F\n", float64(boiling), float64(boiling.ToFahrenheit()))
// 100°C = 212°F
fmt.Println(Celsius(-5).IsFreezing()) // true
```

### String Helpers

```go
type Path string

func (p Path) Base() string {
    s := string(p)
    if i := strings.LastIndex(s, "/"); i >= 0 {
        return s[i+1:]
    }
    return s
}

p := Path("/home/user/report.pdf")
fmt.Println(p.Base()) // report.pdf
```

Named types + methods turn a bare `float64` or `string` into a small, self-documenting API. The receiver-choice rules are identical: these small value-like types are the classic case for **value receivers**.

---

## Methods and Embedding: Method Promotion

Level 13 showed that embedding a struct **promotes its fields**. The same happens to its **methods** - this is Go's composition in action.

```go
type Animal struct {
    Name string
}

func (a Animal) Describe() string { return a.Name + " is an animal" }
func (a Animal) Speak() string    { return a.Name + " makes a sound" }

type Dog struct {
    Animal // embedded: Animal's methods are PROMOTED to Dog
    Breed  string
}

d := Dog{Animal: Animal{Name: "Rex"}, Breed: "Beagle"}
fmt.Println(d.Describe()) // Rex is an animal - defined on Animal, called on Dog!
```

### Overriding a Promoted Method

Define a same-name method on the outer type and it **shadows** the promoted one:

```go
func (d Dog) Speak() string { return d.Name + " says Woof!" }

fmt.Println(d.Speak())        // Rex says Woof!      - Dog's version wins
fmt.Println(d.Animal.Speak()) // Rex makes a sound   - the embedded version, reached explicitly
```

Note the difference from classic OOP inheritance: `d.Animal.Speak()` runs Animal's method on the **embedded Animal value** - there's no virtual dispatch, just a shorter name for a real field access.

### Promotion and Method Sets

Promoted methods join the outer type's method set (pointer-receiver promoted methods join `*Outer`'s). This means embedding is a legitimate way to satisfy interfaces in Level 15.

---

## Method Values and Method Expressions

Two brief but useful forms - methods are values too (Level 11 energy!).

### Method Value: `p.Greet` (receiver bound NOW)

```go
p := Person{Name: "Alice"}
greet := p.Greet // binds a COPY of p as the receiver
fmt.Println(greet()) // Hello, I'm Alice

p.Name = "Alicia"
fmt.Println(greet())   // Hello, I'm Alice  - still the old snapshot!
fmt.Println(p.Greet()) // Hello, I'm Alicia - a fresh call sees the new name
```

For a value receiver, the receiver is **copied at bind time**. (Binding a pointer-receiver method captures the pointer instead, so later changes ARE visible.)

### Method Expression: `Person.Greet` (receiver becomes the 1st argument)

```go
greetFn := Person.Greet // type: func(Person) string
fmt.Println(greetFn(Person{Name: "Bob"})) // Hello, I'm Bob

renameFn := (*Person).Rename // type: func(*Person, string)
renameFn(&p, "Allie")
```

Method values shine as callbacks (`button.OnClick(cart.Checkout)`); method expressions are rarer but great for passing "the method itself" to generic helpers.

---

## Best Practices

### 1. Keep Receiver Names Short and Consistent

```go
// ✅ Good - one letter, same across all methods
func (c *ShoppingCart) AddItem(...)    { }
func (c *ShoppingCart) RemoveItem(...) { }

// ❌ Avoid this/self, and don't vary the name per method
func (this *ShoppingCart) AddItem(...) { }
func (cart *ShoppingCart) RemoveItem(...) { }
```

### 2. Follow the Consistency Rule

If one method needs a pointer receiver, make them ALL pointer receivers. Uniform receivers keep behavior predictable and keep the type's method sets simple for Level 15.

### 3. Default to Pointer Receivers for Structs With Mutating Methods

Mutation is the clearest signal. Big structs (avoid copying) and the presence of other pointer methods (consistency) are the other two.

### 4. Use Value Receivers for Small, Value-Like Types

`Celsius`, `Point`, `RGB` - types you think of as immutable values read best with value receivers, matching how the standard library treats `time.Time`.

### 5. Pair Pointer-Receiver Types With Constructors Returning *T

```go
func NewShoppingCart() *ShoppingCart { return &ShoppingCart{} }
```

Callers get a pointer from the start, so every method call mutates the one true cart - no accidental copies.

### 6. Give Your Types a String() Method

```go
func (c *ShoppingCart) String() string { /* format the cart */ }
fmt.Println(cart) // fmt uses String() automatically
```

In Level 15 you'll learn this makes your type satisfy the `fmt.Stringer` interface - your first real interface win.

---

## Common Mistakes

### Mistake 1: Mutating Through a Value Receiver (Silent No-Op)

```go
// ❌ WRONG - compiles fine, does nothing
func (c Counter) Increment() { c.Count++ }

c := Counter{}
c.Increment()
fmt.Println(c.Count) // still 0!

// ✅ RIGHT
func (c *Counter) Increment() { c.Count++ }
```

This is the #1 methods bug - no error, no warning, just changes that vanish.

### Mistake 2: Calling a Pointer Method on a Non-Addressable Value

```go
// ❌ WRONG - map elements are not addressable
players["alice"].AddPoints(5)
// ./main.go:19:22: cannot call pointer method AddPoints on Player

// ✅ RIGHT - copy out and store back, or use map[string]*Player
p := players["alice"]
p.AddPoints(5)
players["alice"] = p
```

### Mistake 3: Naming the Receiver this or self

```go
// ❌ Not Go style
func (this *Account) Deposit(n float64) { this.Balance += n }

// ✅ Go style - short, first letter of the type
func (a *Account) Deposit(n float64) { a.Balance += n }
```

### Mistake 4: Mixing Value and Pointer Receivers on One Type

```go
// ❌ WRONG - Deposit works, Withdraw silently mutates a copy
func (b *BankAccount) Deposit(n float64)  { b.Balance += n }
func (b BankAccount) Withdraw(n float64)  { b.Balance -= n }

acct := BankAccount{Balance: 100}
acct.Deposit(50)   // Balance: 150
acct.Withdraw(30)  // Balance: STILL 150!
```

### Mistake 5: Forgetting the Store-Back in the Copy-Out Fix

```go
// ❌ WRONG - modified the copy, never updated the map
p := players["alice"]
p.AddPoints(5) // players["alice"] unchanged

// ✅ RIGHT
p := players["alice"]
p.AddPoints(5)
players["alice"] = p
```

### Mistake 6: Expecting a Method Value to Track Later Changes

```go
greet := p.Greet // value receiver: snapshots p NOW
p.Name = "Alicia"
greet() // still uses the OLD p
```

If you need live behavior, call the method fresh each time, or bind a pointer-receiver method.

---

## Summary

**Methods:**
- A method is a function with a receiver: `func (r Rectangle) Area() float64`
- Same power as functions - better organization, and the key to interfaces
- Receiver names: short, consistent, never `this`/`self`

**Receivers:**
- Value receiver `(p Person)` - method gets a COPY; mutations are lost
- Pointer receiver `(p *Person)` - method gets the ORIGINAL; mutations stick
- Consistency rule: one pointer receiver → all pointer receivers

**Addressability:**
- `v.PointerMethod()` is sugar for `(&v).PointerMethod()` - works on addressable values
- Map elements and function return values are NOT addressable → compile error

**Method Sets:**
- `T`: value-receiver methods only
- `*T`: value-receiver AND pointer-receiver methods
- Interfaces (Level 15) are satisfied by method sets - no auto-address sugar there

**Composition:**
- Embedding promotes methods; same-name outer methods override
- Method values bind the receiver now; method expressions take it as the first argument

---

## Next Steps

You now understand:
- ✅ Methods as functions with receivers
- ✅ Value vs pointer receivers and copy semantics
- ✅ The auto-address sugar and its addressability limits
- ✅ Method sets - the foundation for interfaces
- ✅ Method promotion through embedding

**Next level:** Level 15 - Interfaces
- Defining interfaces and implicit satisfaction
- Method sets in action: which types satisfy which interfaces
- The empty interface and type assertions
- fmt.Stringer and the standard library's small-interface style

Your types now have behavior - next they'll get contracts! 🚀
