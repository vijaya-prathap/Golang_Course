# Level 25: Quick Reference Card

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir myapp && cd myapp
go mod init myapp

# Create main.go
cat > main.go << 'EOF'
package main

import "fmt"

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

func main() {
    fmt.Println(Sum([]int{1, 2, 3}))
    fmt.Println(Sum([]float64{1.5, 2.5}))
}
EOF

# Run
go run main.go
```

---

## 📋 Type Parameter Syntax

```go
// Function with type parameters
func Map[T, U any](s []T, f func(T) U) []U { ... }

// Calling - inferred (most common)
Map(nums, toString)

// Calling - explicit type arguments
Map[int, string](nums, toString)

// Generic struct type
type Stack[T any] struct {
    items []T
}

// Methods MUST repeat the type parameter
func (s *Stack[T]) Push(v T) { ... }
```

---

## 🔒 Constraints

```go
any            // accepts anything - needs a type assertion before use
comparable     // supports == / != - required for map keys too

// Custom constraint - defined locally, no imports needed
type Number interface {
    int | int64 | float64
}
```

```go
func Contains[T comparable](s []T, v T) bool {
    for _, item := range s {
        if item == v { return true }
    }
    return false
}
```

---

## 🧩 Generic Functions (Level 11, Revisited)

```go
func Map[T, U any](s []T, f func(T) U) []U {
    result := make([]U, len(s))
    for i, v := range s { result[i] = f(v) }
    return result
}

func Filter[T any](s []T, keep func(T) bool) []T {
    var result []T
    for _, v := range s {
        if keep(v) { result = append(result, v) }
    }
    return result
}

func Reduce[T, U any](s []T, initial U, f func(U, T) U) U {
    acc := initial
    for _, v := range s { acc = f(acc, v) }
    return acc
}
```

---

## 🗼 Generic Types

```go
// Stack[T]
type Stack[T any] struct{ items []T }
func (s *Stack[T]) Push(v T)        { s.items = append(s.items, v) }
func (s *Stack[T]) Pop() (T, bool)  { ... }

// Pair[K, V]
type Pair[K, V any] struct {
    Key   K
    Value V
}

// Cache[K, V]
type Cache[K comparable, V any] struct {
    data map[K]V
}
```

---

## 🎯 Type Inference

```go
// Works: T, U both appear in regular arguments
Map(nums, toString)

// Fails: T appears ONLY in the return type
func Zero[T any]() T { var z T; return z }
zero := Zero()          // ❌ cannot infer T
zero := Zero[int]()     // ✅ must specify explicitly
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| `==` without `comparable` | `func F[T any](a, b T) bool { return a == b }` | `func F[T comparable](a, b T) bool { return a == b }` |
| Generics for pure behavior | `func TotalArea[T Shape](s []T) float64` | `func TotalArea(s []Shape) float64` |
| Wrong constraint type | `Sum([]string{...})` with `Number` constraint | Pass only `int`/`int64`/`float64` slices |
| Return-only type param | `zero := Zero()` | `zero := Zero[int]()` |
| Method missing type param | `func (s *Stack) Push(v T)` | `func (s *Stack[T]) Push(v T)` |

---

## 🚨 Real Compile Errors (Verified)

```
invalid operation: item == v (incomparable types in type set)
string does not satisfy Number (string missing in int | int64 | float64)
in call to Zero, cannot infer T (declared at ./main.go:5:11)
invalid operation: v + 1 (mismatched types any and untyped int)
```

---

## 🎓 Before Next Level

Can you:
- [ ] Write `func Name[T Constraint](params) ReturnType` from memory?
- [ ] Explain when `comparable` is required vs. `any`?
- [ ] Define a custom constraint with a type set?
- [ ] Build a generic struct type with methods?
- [ ] Spot a call where type inference will fail?
- [ ] Decide between a generic type parameter and a small interface?

If YES → You're ready for Level 26!

---

## 📚 Next Level

Level 26: Testing
- Table-driven tests with the `testing` package
- `go test`, subtests, and coverage
- Testing generic functions across multiple type instantiations

You've got generics down! 💪
