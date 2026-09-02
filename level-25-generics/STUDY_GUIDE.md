# Level 25: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: Generics Fundamentals
```
Day 1:  Why generics - duplication vs. any vs. type parameters
Day 2:  Type parameter syntax, calling with/without explicit type args
Day 3:  any and comparable, the two built-in constraints
Day 4:  Custom constraints with type sets (Number)
Day 5:  Generic functions - Map, Filter, Reduce revisited
Day 6:  Generic types - Stack[T], Pair[K, V]
Day 7:  Type inference rules and when interfaces beat generics
```

### Week 2: Practice & Application
```
Day 1:  Exercises 1-3 (before/after, syntax, any)
Day 2:  Exercises 4-5 (comparable, custom constraints + real compile errors)
Day 3:  Exercises 6-7 (Map/Filter/Reduce, Stack[T])
Day 4:  Exercises 8-10 (Pair, inference, Cache[K, V])
Day 5:  Bonus challenges
Day 6-7: Practice & review
```

---

## 🎯 Type Parameter Syntax Anatomy

```
     ┌─ function name
     │      ┌─ type parameter list (square brackets)
     │      │
func Map [ T, U   any ] ( s []T, f func(T) U )  []U
              │     │        │         │           │
              │     │        │         │           └─ return type, built from U
              │     │        │         └─ regular parameter, built from T and U
              │     │        └─ regular parameter, built from T
              │     └─ constraint - applies to every name before it in the list
              └─ type parameter names (placeholders, filled in per call)
```

```
Generic type (struct), same idea, repeated on every method:

type Stack[T any] struct {        func (s *Stack[T]) Push(v T)
    items []T                     func (s *Stack[T]) Pop() (T, bool)
}                                  └─ T must be restated on the receiver
```

| Calling Style | Example | When |
|---|---|---|
| Inferred | `Map(nums, toString)` | Type parameters appear in the regular arguments |
| Explicit | `Map[int, string](nums, toString)` | Always valid; required when inference can't see enough |

---

## 🔒 Constraint Hierarchy

```
any               →  accepts literally anything
                     └─ can't do anything type-specific without a type assertion

comparable        →  any type supporting == / !=
                     └─ required for == inside a generic function, or as a map key type

custom (type set) →  YOUR OWN interface listing exact allowed types
                     └─ type Number interface { int | int64 | float64 }
                     └─ narrowest, most expressive, fully local (no imports needed)
```

| Constraint | Satisfied By | Lets You... | Example |
|---|---|---|---|
| `any` | every type | store, pass around, type-assert | `func Describe(v any) string` |
| `comparable` | types supporting `==`/`!=` | compare with `==`, use as map key | `func Contains[T comparable](s []T, v T) bool` |
| `type Number interface { int \| int64 \| float64 }` | exactly the listed types | arithmetic (`+`, `-`, etc.) | `func Sum[T Number](nums []T) T` |

**Rule:** pick the *narrowest* constraint that lets the function body compile - `any` when you truly don't touch the value's type, `comparable` the moment you need `==`, a custom type-set constraint the moment you need arithmetic or another type-specific operator.

---

## 🌳 Generics vs. Interfaces: Decision Tree

Ties back to Level 15's "small interfaces" guidance - the question is whether you need an arbitrary **element type** or arbitrary **behavior**.

```
Do different concrete types need to behave DIFFERENTLY behind one contract?
(e.g., Circle.Area() vs Rectangle.Area() - same method, different math)
│
├─ YES → use an INTERFACE (Level 15)
│         - one or two methods, defined at the consumer
│         - a single []Shape slice can mix Circle and Rectangle values
│
└─ NO → is the SAME algorithm/structure needed for arbitrary element types?
         (e.g., a stack that holds ints today, strings tomorrow)
         │
         ├─ YES → use a TYPE PARAMETER (this level)
         │         - Stack[T], Map[T, U], Sum[T Number]
         │         - the element type is just stored/compared/transformed,
         │           never asked to "do" anything itself
         │
         └─ NO (only one type, one call site) → use a CONCRETE type/function
                   - don't generic-ize "just in case"
```

| Question | Interface | Generic Type Parameter |
|---|---|---|
| What varies between uses? | Behavior (what a value can *do*) | The element type (what's *stored*) |
| Can one slice/collection mix types? | Yes (`[]Shape`) | No (`Stack[int]` ≠ `Stack[string]`) |
| Typical method count | 1-2 | N/A (constraint lists types, not methods, unless it embeds one) |
| Go proverb | "The bigger the interface, the weaker the abstraction" | "Don't generic-ize just in case" |

---

## 🗼 Generic Stack[T]: Memory & Operation Diagram

```
type Stack[T any] struct {
    items []T
}

Stack[int]{items: []int{10, 20, 30}}

    index:   0    1    2
            ┌────┬────┬────┐
   items:   │ 10 │ 20 │ 30 │   ← top of stack is the LAST element
            └────┴────┴────┘
                        ▲
                       top (index len-1)

Push(40):
            ┌────┬────┬────┬────┐
   items:   │ 10 │ 20 │ 30 │ 40 │   ← append grows the backing slice
            └────┴────┴────┴────┘

Pop() -> (40, true):
            ┌────┬────┬────┐
   items:   │ 10 │ 20 │ 30 │   ← slice re-sliced to items[:len-1]
            └────┴────┴────┘

Pop() on empty -> (zero value, false):
   items:   []                 ← returns T's zero value (0 for int, "" for
                                  string, nil for a pointer/slice/map), plus
                                  false so the caller can tell "empty" apart
                                  from "the zero value was actually stored"
```

```
Stack[int] and Stack[string] are DISTINCT types at compile time:

    var ints  Stack[int]      var words Stack[string]
    ints.Push(1)               words.Push("go")
    ints.Push(2)                words.Push("generics")

    ints.Push("oops")   ← does NOT compile: "oops" (string) is not int
```

---

## 🚨 Common Mistakes

### Mistake 1: Forgetting comparable, Using == Anyway

```go
// ❌ WRONG
func Contains[T any](s []T, v T) bool {
    for _, item := range s {
        if item == v { return true }   // any doesn't guarantee ==
    }
    return false
}
```
Real error: `invalid operation: item == v (incomparable types in type set)`

### Mistake 2: Overusing Generics for Behavior-Based Polymorphism

```go
// ❌ Overkill - only works with ONE concrete shape type per call
func TotalArea[T Shape](shapes []T) float64 { ... }

// ✅ RIGHT - a plain interface slice mixes concrete types freely
func TotalArea(shapes []Shape) float64 { ... }
```

### Mistake 3: Constraint Mismatch

```go
type Number interface { int | int64 | float64 }
func Sum[T Number](nums []T) T { ... }

Sum([]string{"a", "b"})   // string isn't in Number's type set
```
Real error: `string does not satisfy Number (string missing in int | int64 | float64)`

### Mistake 4: Expecting Inference When T Is Only in the Return Type

```go
func Zero[T any]() T { var z T; return z }
zero := Zero()   // nothing to infer T from
```
Real error: `in call to Zero, cannot infer T (declared at ./main.go:5:11)`

### Mistake 5: Forgetting to Restate the Type Parameter on a Method Receiver

```go
// ❌ WRONG - won't compile
func (s *Stack) Push(v T) { ... }

// ✅ RIGHT
func (s *Stack[T]) Push(v T) { ... }
```

---

## 📈 Progression Summary

### Understanding Level 25

Level 25 teaches how to write one function or type that safely works across many types:

1. **The problem** - duplication (one function per type) vs. `any` (loses type safety)
2. **Syntax** - `[T Constraint]`, inference vs. explicit type arguments
3. **Built-in constraints** - `any` (anything) and `comparable` (`==`/`!=`, map keys)
4. **Custom constraints** - your own type-set interfaces, fully offline
5. **Generic functions & types** - `Map`/`Filter`/`Reduce`, `Stack[T]`, `Pair[K, V]`, `Cache[K, V]`
6. **Judgment** - knowing when a small interface (Level 15) is the better tool

### Prerequisites for Level 26

Before moving to Level 26 (Testing), you need:

- ✅ Comfortable declaring and calling generic functions
- ✅ Know the difference between `any` and `comparable`, and when each is required
- ✅ Can write a custom constraint interface with a type set
- ✅ Can build a generic type (struct) with methods
- ✅ Can predict when type inference works and when you must specify type arguments
- ✅ Can explain when an interface is the better choice over a type parameter

### Ready for Level 26?

Level 26 teaches how to verify all of this code actually works, automatically:
- Table-driven tests with the `testing` package
- Running the same test table against multiple type instantiations of a generic function
- `go test`, subtests, and coverage
- Building test doubles with the small-interface patterns from Level 15

---

## ✅ Checklist Before Level 26

- [ ] Can write `func Name[T Constraint](params) ReturnType` from memory
- [ ] Can explain why `any` requires a type assertion before doing anything type-specific
- [ ] Can explain why `comparable` is required for `==` or a map key type parameter
- [ ] Can define a custom constraint interface with a type set (`int | int64 | float64`)
- [ ] Can build a generic struct type and repeat its type parameter on every method
- [ ] Can identify a call where type inference fails and knows how to fix it
- [ ] Can explain when a small interface beats a generic type parameter
- [ ] Completed 8+ exercises, including both real compile-error exercises

---

## 💡 Key Takeaways

### The Trade-Off
Generics give you the reusability of `any` without giving up the compiler's type checking - one implementation, many types, zero runtime type-assertion risk.

### The Syntax
Type parameters live in `[...]` right after the name, on both functions and types - and generic types must repeat their type parameters on every method receiver.

### The Constraints
`any` (anything) → `comparable` (`==`/`!=`, map keys) → your own type-set interface (arithmetic, custom operators) - pick the narrowest one that compiles.

### The Judgment Call
Generics are for arbitrary **element types** doing the same algorithm. Interfaces (Level 15) are for arbitrary **behavior** behind a shared contract. Don't reach for a type parameter just because it's new and interesting.

---

## 📚 Next Level

Level 26: Testing
- Table-driven tests with the `testing` package
- `go test`, subtests, and coverage
- Testing generic functions across multiple type instantiations
- Test doubles built from Level 15's small-interface patterns

You've got generics down! Keep going! 🚀
