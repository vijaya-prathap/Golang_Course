# Level 25: Generics - Complete Guide

## Introduction

Welcome to Level 25! You've mastered Context (Level 24) - cancellation signals, deadlines, and request-scoped values flowing cleanly through a call chain. Now it's time to learn **generics**: writing a single function or type that works correctly across many types, without giving up the compiler's type checking.

Before Go 1.18, you had exactly two choices when you wanted one piece of logic to work on more than one type:

1. **Write it once per type** (`SumInts`, `SumFloats`, `SumInt64s`, ...) - correct and fast, but you're copy-pasting the same loop over and over.
2. **Write it once using `any`/`interface{}`** - one function, but you lose the compiler's help. Every use requires a runtime type assertion, and a mismatched type is a `panic`, not a compile error.

Generics give you a third option: **one function, one implementation, full compile-time type safety** - the compiler checks every call site and generates the right code for each type you use it with. This level is marked **Advanced** because generics syntax is dense and the constraint system takes some getting used to - but the payoff is real: less duplicated code, and none of the runtime type-assertion risk that `any` carries.

---

## Table of Contents

1. [Why Generics: The Problem They Solve](#why-generics-the-problem-they-solve)
2. [Type Parameter Syntax](#type-parameter-syntax)
3. [The any Constraint: Loosest Possible, Least Useful](#the-any-constraint-loosest-possible-least-useful)
4. [The comparable Constraint](#the-comparable-constraint)
5. [Custom Constraints With Type Sets](#custom-constraints-with-type-sets)
6. [Generic Functions Revisited: Map, Filter, Reduce](#generic-functions-revisited-map-filter-reduce)
7. [Generic Types: Stack[T] and Pair[K, V]](#generic-types-stackt-and-pairk-v)
8. [Type Inference Rules](#type-inference-rules)
9. [When NOT to Use Generics](#when-not-to-use-generics)
10. [Best Practices](#best-practices)
11. [Common Mistakes](#common-mistakes)

---

## Why Generics: The Problem They Solve

### Before Generics: Duplication

```go
func SumInts(nums []int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

func SumFloats(nums []float64) float64 {
    total := 0.0
    for _, n := range nums {
        total += n
    }
    return total
}
```

Same logic, copy-pasted for every numeric type you need. Add `int64` support later and you copy-paste again.

### Before Generics: any (Loses Type Safety)

```go
func SumAny(nums []any) any {
    var total int
    for _, n := range nums {
        total += n.(int) // type assertion - panics if the slice holds anything else
    }
    return total
}
```

One function - but the compiler no longer checks that `nums` actually holds `int`s. Pass a `[]any{1, 2, "oops"}` and `SumAny` compiles fine and panics at runtime.

### After: One Generic Function

```go
type Number interface {
    int | int64 | float64
}

func Sum[T Number](nums []T) T {
    var total T
    for _, n := range nums {
        total += n
    }
    return total
}
```

One implementation. `Sum([]int{1, 2, 3})` and `Sum([]float64{1.5, 2.5})` both compile, both run, and passing a `[]string` is a **compile error**, not a runtime panic. This is the whole pitch of generics: write it once, keep the type safety.

---

## Type Parameter Syntax

A **type parameter** is declared in square brackets between the function name and its regular parameter list.

### Anatomy

```go
func Map[T, U any](s []T, f func(T) U) []U {
    result := make([]U, len(s))
    for i, v := range s {
        result[i] = f(v)
    }
    return result
}

//     ┌─ type parameter list
//     │
func Map[T, U any](s []T, f func(T) U) []U
//         │            │              │
//         │            │              └─ return type, built from U
//         │            └─ regular parameters, built from T and U
//         └─ constraint - both T and U may be any type
```

`T` and `U` are placeholders for real types, filled in at each call site. `any` here is the **constraint** - it says "T and U can be literally any type."

### Calling: Inference vs. Explicit Type Arguments

```go
nums := []int{1, 2, 3, 4}

// Type inference - the compiler figures out T=int and U=string from the arguments
labels := Map(nums, func(n int) string {
    return fmt.Sprintf("#%d", n)
})

// Same call, with explicit type arguments - never wrong, sometimes required
labels2 := Map[int, string](nums, func(n int) string {
    return fmt.Sprintf("#%d", n)
})
```

Both calls above produce identical results. Go infers type arguments whenever it can see them in the function's regular arguments, so most calls skip the `[T, U]` part entirely.

**Verified output** (see Exercise 2):

```
=== Calling With Type Inference (no explicit type args) ===
doubled: [2 4 6 8]

=== Calling With Explicit Type Arguments ===
labels: [#1 #2 #3 #4]
```

---

## The any Constraint: Loosest Possible, Least Useful

`any` is a built-in alias for `interface{}` - the empty interface, satisfied by every type. As a **constraint**, `any` means "accept literally anything," which is exactly as permissive (and exactly as limiting) as it sounds.

```go
func Describe(v any) string {
    switch x := v.(type) {
    case int:
        return fmt.Sprintf("int: %d", x)
    case string:
        return fmt.Sprintf("string: %q", x)
    default:
        return fmt.Sprintf("unknown type: %v", x)
    }
}
```

The catch: once a value is typed `any`, you can't do anything type-specific with it - no arithmetic, no field access, no method calls - until you narrow it back down with a type assertion or type switch.

```go
var v any = 10
fmt.Println(v + 1) // does NOT compile
```

**Verified compile error** (see `/private/tmp/level25-verify/readme-any` during authoring):

```
./main.go:7:14: invalid operation: v + 1 (mismatched types any and untyped int)
```

Contrast this with a real type parameter: `func Sum[T Number](nums []T) T` lets you write `total += n` directly, because the constraint `Number` promises the compiler that every possible `T` supports `+`. `any` promises nothing beyond "it's a value," so the compiler won't let you do anything beyond storing and retrieving it.

---

## The comparable Constraint

`comparable` is a second built-in constraint. It's satisfied by any type that supports `==` and `!=` - which is also exactly what you need to use a type as a **map key**.

```go
func Contains[T comparable](s []T, v T) bool {
    for _, item := range s {
        if item == v {
            return true
        }
    }
    return false
}
```

```go
Contains([]int{10, 20, 30}, 20)        // true
Contains([]string{"go", "rust"}, "go") // true
```

`comparable` is required the moment your generic function uses `==`/`!=` on a value of the type parameter, or uses the type parameter as a `map[K]V` key type. Leaving the constraint too loose (`any` instead of `comparable`) and then using `==` anyway is a compile error - see [Common Mistakes](#common-mistakes) for the exact wording.

---

## Custom Constraints With Type Sets

`any` and `comparable` are built in, but most real constraints are ones **you** define - a plain interface listing the types the type parameter is allowed to be. This is called a **type set**, and it needs no external package:

```go
// Defined locally - works fully offline, no imports beyond fmt for printing.
type Number interface {
    int | int64 | float64
}

func Sum[T Number](nums []T) T {
    var total T
    for _, n := range nums {
        total += n
    }
    return total
}

func Average[T Number](nums []T) float64 {
    if len(nums) == 0 {
        return 0
    }
    return float64(Sum(nums)) / float64(len(nums))
}
```

The `|` inside the interface lists the **type set** - the types allowed to satisfy `Number`. Any type not in that list fails at compile time:

**Verified compile error** (see Exercise 5):

```
./main.go:20:17: string does not satisfy Number (string missing in int | int64 | float64)
```

A type set can also use `~int` ("any type whose underlying type is `int`") to include user-defined types like `type Celsius int`, but a plain `int | int64 | float64` is enough for most everyday constraints.

---

## Generic Functions Revisited: Map, Filter, Reduce

Level 11 introduced `apply` and `filter` as higher-order functions, but they only worked on `[]int` (or whatever single type you hard-coded). With generics, they become truly reusable across any element type - and any output type:

```go
func Map[T, U any](s []T, f func(T) U) []U {
    result := make([]U, len(s))
    for i, v := range s {
        result[i] = f(v)
    }
    return result
}

func Filter[T any](s []T, keep func(T) bool) []T {
    var result []T
    for _, v := range s {
        if keep(v) {
            result = append(result, v)
        }
    }
    return result
}

func Reduce[T, U any](s []T, initial U, f func(U, T) U) U {
    acc := initial
    for _, v := range s {
        acc = f(acc, v)
    }
    return acc
}
```

Notice `Map` and `Reduce` both need **two** type parameters - the input element type `T` and the output/accumulator type `U`, which can legitimately differ (mapping `[]int` to `[]string`, or reducing `[]string` down to a single `int` count). `Filter` only needs one, because it returns the same type it was given.

**Verified output** (see Exercise 6):

```
=== Map: int -> string (different output type) ===
[n1 n2 n3 n4 n5 n6]

=== Reduce: build a string (different accumulator type) ===
go-is-fun-and-fast
```

---

## Generic Types: Stack[T] and Pair[K, V]

Type parameters aren't just for functions - `struct` types can carry them too, letting you build a data structure once and reuse it for any element type.

### A Generic Stack

```go
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(v T) {
    s.items = append(s.items, v)
}

func (s *Stack[T]) Pop() (T, bool) {
    var zero T
    if len(s.items) == 0 {
        return zero, false
    }
    last := len(s.items) - 1
    v := s.items[last]
    s.items = s.items[:last]
    return v, true
}

func (s *Stack[T]) Peek() (T, bool) {
    var zero T
    if len(s.items) == 0 {
        return zero, false
    }
    return s.items[len(s.items)-1], true
}
```

Every method on a generic type repeats the type parameter in the receiver: `func (s *Stack[T]) Pop()`, not `func (s *Stack) Pop()`. `var ints Stack[int]` and `var words Stack[string]` are two independent stacks - the compiler treats `Stack[int]` and `Stack[string]` as distinct types.

### A Generic Pair

```go
type Pair[K, V any] struct {
    Key   K
    Value V
}

func NewPair[K, V any](key K, value V) Pair[K, V] {
    return Pair[K, V]{Key: key, Value: value}
}
```

`Pair[string, int]{Key: "Alice", Value: 30}` bundles two independently-typed values together - useful anywhere you'd otherwise reach for a two-element `any` slice or a throwaway single-use struct.

**Verified output** (see Exercises 7 and 8):

```
=== Stack[int] ===
Pop: 3
Pop: 2
Pop: 1
Pop on empty, ok = false

=== Pair[string, int] ===
(Alice, 30)
```

---

## Type Inference Rules

Go infers type arguments **from the types of the regular arguments you pass** - it does not (and cannot) look at the return type to guess what you want.

```go
func Map[T, U any](s []T, f func(T) U) []U { ... }

// T=int, U=string - both inferred from `nums` and the function literal's signature
Map(nums, func(n int) string { return fmt.Sprintf("#%d", n) })
```

That works because `T` and `U` both show up somewhere in the parameter list. Now consider a function where the type parameter appears **only** in the return type:

```go
func Zero[T any]() T {
    var zero T
    return zero
}
```

`Zero()` takes no arguments at all, so there's nothing for the compiler to infer `T` from - even though the call itself would otherwise be unambiguous to a human reader. You must supply the type argument explicitly:

```go
zero := Zero() // does NOT compile
```

**Verified compile error** (see Exercise 9):

```
./main.go:11:14: in call to Zero, cannot infer T (declared at ./main.go:5:11)
```

```go
zeroInt := Zero[int]()       // works: 0
zeroString := Zero[string]() // works: ""
```

**Rule of thumb:** if every type parameter appears in at least one regular argument, you can usually drop the explicit type arguments. If a type parameter appears *only* in the return type (a common shape for constructors and "zero value" helpers), you must specify it explicitly at the call site.

---

## When NOT to Use Generics

Generics solve a specific problem: **the same algorithm or data structure, applied to arbitrary element types**. They are not a general-purpose tool for "flexibility," and reaching for them when a plain interface already does the job adds a type parameter with no real payoff.

Recall Level 15's design guidance: *"the bigger the interface, the weaker the abstraction,"* and prefer small, consumer-defined interfaces over premature abstraction. That guidance still applies here - and it tells you when to reach for an interface instead of a type parameter.

```go
// A plain interface already expresses "things that have an area" cleanly.
type Shape interface {
    Area() float64
}

type Circle struct{ Radius float64 }
type Rectangle struct{ Width, Height float64 }

func (c Circle) Area() float64    { return 3.14159 * c.Radius * c.Radius }
func (r Rectangle) Area() float64 { return r.Width * r.Height }

func TotalArea(shapes []Shape) float64 {
    total := 0.0
    for _, s := range shapes {
        total += s.Area()
    }
    return total
}
```

`TotalArea` needs **behavior** (something that can compute its own area), not an arbitrary element type - a `[]Shape` of mixed `Circle`s and `Rectangle`s already works today, in one function, with zero type parameters. Rewriting it as `TotalArea[T Shape](shapes []T) float64` would only work with a slice of a *single* concrete shape type - strictly less useful, for more syntax.

**The distinction:**

| Use a generic type parameter when... | Use an interface when... |
|---|---|
| You need the same algorithm/structure over arbitrary element types (`Stack[T]`, `Map`, `Sum`) | You need arbitrary **behavior** behind a common contract (something that can `Area()`, `Read()`, `String()`) |
| The element type doesn't need to do anything - it's just stored, compared, or transformed | Different concrete types genuinely behave differently, and callers need to mix them in one slice |
| You'd otherwise duplicate the same function/struct per type | A single method (or two) already captures what callers need |

Generics shine for algorithms and data structures over arbitrary element types. For polymorphism based on *what a value can do*, Level 15's small interfaces are usually the simpler, more idiomatic answer - don't generic-ize "just in case."

---

## Best Practices

### 1. Keep Constraints as Narrow as Truly Needed

```go
// ✅ Good - only requires what Sum actually uses (+)
type Number interface {
    int | int64 | float64
}
func Sum[T Number](nums []T) T { ... }

// ❌ Overly loose - any doesn't guarantee + is even defined
func Sum[T any](nums []T) T { ... } // won't compile once you try total += n
```

### 2. Don't Generic-ize Prematurely

Start with a concrete function (`func SumInts(nums []int) int`). Reach for a type parameter only once you actually need a second type, or clearly will soon. A generic function with exactly one call site and one concrete type is just indirection.

### 3. Prefer Named Constraint Types Over Inline Unions

```go
// ✅ Good - reusable, self-documenting, one place to extend the type set
type Number interface {
    int | int64 | float64
}
func Sum[T Number](nums []T) T { ... }
func Average[T Number](nums []T) float64 { ... }

// ❌ Works, but the union is invisible from the function signature and
// can't be reused by the next function that needs the same constraint
func Sum[T int | int64 | float64](nums []T) T { ... }
```

### 4. Repeat the Type Parameter on Every Method Receiver

```go
// ✅ Good
func (s *Stack[T]) Push(v T) { ... }

// ❌ Won't compile - a generic type's methods must restate its type parameters
func (s *Stack) Push(v T) { ... }
```

### 5. Use comparable Only When You Actually Compare or Key a Map

Don't default to `comparable` out of caution - use `any` unless the function body genuinely needs `==`/`!=` or a map key.

---

## Common Mistakes

### Mistake 1: Forgetting comparable and Using == Anyway

```go
// ❌ WRONG - any does not guarantee == is defined for every possible T
func Contains[T any](s []T, v T) bool {
    for _, item := range s {
        if item == v {
            return true
        }
    }
    return false
}
```

**Real compile error:**

```
./main.go:9:6: invalid operation: item == v (incomparable types in type set)
```

```go
// ✅ RIGHT - comparable guarantees == and != are defined
func Contains[T comparable](s []T, v T) bool {
    for _, item := range s {
        if item == v {
            return true
        }
    }
    return false
}
```

### Mistake 2: Overusing Generics Where a Plain Interface Would Be Simpler

```go
// ❌ Overkill - Shape is about behavior (Area()), not an arbitrary stored type
func TotalArea[T Shape](shapes []T) float64 { ... } // only works with ONE concrete shape type per call

// ✅ RIGHT - a plain interface slice mixes concrete types freely
func TotalArea(shapes []Shape) float64 { ... }
```

### Mistake 3: Passing a Type That Doesn't Satisfy a Custom Constraint

```go
type Number interface {
    int | int64 | float64
}

func Sum[T Number](nums []T) T { ... }

words := []string{"a", "b", "c"}
fmt.Println(Sum(words)) // string is not in Number's type set
```

**Real compile error:**

```
./main.go:20:17: string does not satisfy Number (string missing in int | int64 | float64)
```

### Mistake 4: Expecting Inference to Work When T Only Appears in the Return Type

```go
func Zero[T any]() T {
    var zero T
    return zero
}

zero := Zero() // cannot infer T - there's no argument to infer it from
```

**Real compile error:**

```
./main.go:11:14: in call to Zero, cannot infer T (declared at ./main.go:5:11)
```

```go
// ✅ RIGHT - supply the type argument explicitly
zero := Zero[int]()
```

---

## Summary

**Why Generics:**
- One function/type, many types, full compile-time safety - no per-type duplication, no `any` + runtime type assertions

**Syntax:**
- `func Name[T Constraint](params) ReturnType` - type parameters live in `[...]` right after the name
- Generic types repeat their type parameters on every method receiver: `func (s *Stack[T]) Push(v T)`

**Constraints:**
- `any` - accepts anything, but you can't do much without a type assertion
- `comparable` - required for `==`/`!=` or using the type parameter as a map key
- Custom (`type Number interface { int | int64 | float64 }`) - your own type set, defined locally, no external packages needed

**Inference:**
- Go infers type arguments from the regular arguments you pass
- If a type parameter only appears in the return type, you must specify it explicitly

**When to Reach for What:**
- Same algorithm/structure over arbitrary element types → generics (`Stack[T]`, `Sum[T Number]`)
- Arbitrary behavior behind a common contract → a small interface (Level 15), not a type parameter

---

## Next Steps

You now understand:
- ✅ Why generics exist and the exact problem (duplication vs. lost type safety) they solve
- ✅ Type parameter syntax, and how to call generic functions with and without explicit type arguments
- ✅ `any` and `comparable`, the two built-in constraints
- ✅ Defining your own constraints with type sets, entirely offline
- ✅ Generic functions (`Map`/`Filter`/`Reduce`) and generic types (`Stack[T]`, `Pair[K, V]`)
- ✅ When type inference works, when it fails, and what the real compiler errors look like
- ✅ When a plain interface is the better tool than a generic type parameter

**Next level:** Level 26: Testing
- Writing table-driven tests with the `testing` package
- `go test`, subtests, and test coverage
- Testing generic functions across multiple type instantiations
- Mocking and test doubles using the small-interface patterns from Level 15

Generics round out Go's type system toolbox - between concrete types, interfaces, and now type parameters, you can express almost any shape of reusable code. Next you'll learn to prove that code actually works! 🚀
