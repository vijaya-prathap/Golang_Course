# Level 25: Generics - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

---

## Exercise 1: Before/After - Duplication vs. any vs. Generics

**Objective:** See the exact problem generics solve, side by side

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level25-exercise1
cd ~/projects/level25-exercise1
go mod init level25.example/exercise1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

// ===== BEFORE: one function per type =====

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

// ===== BEFORE: one function using any (loses type safety) =====

func SumAny(nums []any) any {
    var total int
    for _, n := range nums {
        // Must type-assert before doing arithmetic - runtime risk, not compile-time safety
        total += n.(int)
    }
    return total
}

// ===== AFTER: one generic function works for both, safely =====

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
    ints := []int{1, 2, 3, 4, 5}
    floats := []float64{1.5, 2.5, 3.0}

    fmt.Println("=== Before: Duplicated Type-Specific Functions ===")
    fmt.Println("SumInts:", SumInts(ints))
    fmt.Println("SumFloats:", SumFloats(floats))

    fmt.Println("\n=== Before: any-Based Function (loses type safety) ===")
    anyInts := []any{1, 2, 3, 4, 5}
    fmt.Println("SumAny:", SumAny(anyInts))
    fmt.Println("Note: SumAny would panic at runtime if the slice held a float64 or string")

    fmt.Println("\n=== After: One Generic Function, Fully Type-Safe ===")
    fmt.Println("Sum(ints):", Sum(ints))
    fmt.Println("Sum(floats):", Sum(floats))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Before: Duplicated Type-Specific Functions ===
SumInts: 15
SumFloats: 7

=== Before: any-Based Function (loses type safety) ===
SumAny: 15
Note: SumAny would panic at runtime if the slice held a float64 or string

=== After: One Generic Function, Fully Type-Safe ===
Sum(ints): 15
Sum(floats): 7
```

**Learning Objectives:**
- ✅ See type-specific duplication side by side with a single generic function
- ✅ Understand why `any` compiles but loses the compiler's type checking
- ✅ Confirm one generic `Sum[T Number]` replaces both `SumInts` and `SumFloats`

---

## Exercise 2: Type Parameter Syntax and Inference

**Objective:** Declare a generic function and call it with and without explicit type arguments

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level25-exercise2
cd ~/projects/level25-exercise2
go mod init level25.example/exercise2
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

// Map applies f to every element of s, producing a new slice of (possibly) a different type.
func Map[T, U any](s []T, f func(T) U) []U {
    result := make([]U, len(s))
    for i, v := range s {
        result[i] = f(v)
    }
    return result
}

func main() {
    nums := []int{1, 2, 3, 4}

    fmt.Println("=== Calling With Type Inference (no explicit type args) ===")
    doubled := Map(nums, func(n int) int {
        return n * 2
    })
    fmt.Println("doubled:", doubled)

    fmt.Println("\n=== Calling With Explicit Type Arguments ===")
    labels := Map[int, string](nums, func(n int) string {
        return fmt.Sprintf("#%d", n)
    })
    fmt.Println("labels:", labels)

    fmt.Println("\n=== Mapping to a Different Type (int -> bool) ===")
    isEven := Map(nums, func(n int) bool {
        return n%2 == 0
    })
    fmt.Println("isEven:", isEven)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Calling With Type Inference (no explicit type args) ===
doubled: [2 4 6 8]

=== Calling With Explicit Type Arguments ===
labels: [#1 #2 #3 #4]

=== Mapping to a Different Type (int -> bool) ===
isEven: [false true false true]
```

**Learning Objectives:**
- ✅ Declare a two-type-parameter generic function
- ✅ Call it with type inference and with explicit type arguments
- ✅ See that the output type (`U`) can differ from the input type (`T`)

---

## Exercise 3: The any Constraint

**Objective:** See what any does and doesn't let you do with a value

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level25-exercise3
cd ~/projects/level25-exercise3
go mod init level25.example/exercise3
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

// Describe accepts absolutely anything - the loosest possible constraint.
func Describe(v any) string {
    switch x := v.(type) {
    case int:
        return fmt.Sprintf("int: %d (doubled: %d)", x, x*2)
    case string:
        return fmt.Sprintf("string: %q (length: %d)", x, len(x))
    case bool:
        return fmt.Sprintf("bool: %t", x)
    default:
        return fmt.Sprintf("unknown type: %v", x)
    }
}

// PrintAll takes a slice of any - it can hold mixed types, but you lose
// compile-time knowledge of what's inside.
func PrintAll(items []any) {
    for _, item := range items {
        fmt.Println(Describe(item))
    }
}

func main() {
    fmt.Println("=== any Accepts Anything ===")
    items := []any{42, "hello", true, 3.14}
    PrintAll(items)

    fmt.Println("\n=== Without a Type Assertion, You Can't Do Much ===")
    var v any = 10
    // fmt.Println(v + 1) // would NOT compile: invalid operation: v + 1 (mismatched types any and untyped int)
    n, ok := v.(int)
    if ok {
        fmt.Println("after asserting to int, v + 1 =", n+1)
    }
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== any Accepts Anything ===
int: 42 (doubled: 84)
string: "hello" (length: 5)
bool: true
unknown type: 3.14

=== Without a Type Assertion, You Can't Do Much ===
after asserting to int, v + 1 = 11
```

**Learning Objectives:**
- ✅ Use `any` as the loosest possible constraint/parameter type
- ✅ See that a type switch is required to do anything type-specific with an `any` value
- ✅ Confirm arithmetic directly on an `any` value does not compile (see the commented line)

---

## Exercise 4: The comparable Constraint

**Objective:** Write a generic Contains using comparable, then trigger the real compile error from skipping it

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level25-exercise4
cd ~/projects/level25-exercise4
go mod init level25.example/exercise4
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

// Contains requires comparable because it uses == to compare elements.
func Contains[T comparable](s []T, v T) bool {
    for _, item := range s {
        if item == v {
            return true
        }
    }
    return false
}

// IndexOf reuses the comparable constraint too.
func IndexOf[T comparable](s []T, v T) int {
    for i, item := range s {
        if item == v {
            return i
        }
    }
    return -1
}

func main() {
    ints := []int{10, 20, 30, 40}
    words := []string{"go", "rust", "python"}

    fmt.Println("=== Contains With ints ===")
    fmt.Println(Contains(ints, 30))
    fmt.Println(Contains(ints, 99))

    fmt.Println("\n=== Contains With strings ===")
    fmt.Println(Contains(words, "rust"))
    fmt.Println(Contains(words, "java"))

    fmt.Println("\n=== IndexOf ===")
    fmt.Println(IndexOf(words, "python"))
    fmt.Println(IndexOf(words, "java"))

    fmt.Println("\n=== comparable Also Works as a Map Key Type Parameter ===")
    countOccurrences := func(s []string) map[string]int {
        counts := make(map[string]int)
        for _, v := range s {
            counts[v]++
        }
        return counts
    }
    fmt.Println(countOccurrences([]string{"a", "b", "a"}))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Contains With ints ===
true
false

=== Contains With strings ===
true
false

=== IndexOf ===
2
-1

=== comparable Also Works as a Map Key Type Parameter ===
map[a:2 b:1]
```

4. Now trigger the real compile error. In a **separate** scratch directory, create a broken version that constrains with `any` instead of `comparable` but still uses `==`:

```bash
mkdir -p ~/projects/level25-exercise4-broken
cd ~/projects/level25-exercise4-broken
go mod init level25.example/exercise4broken

cat > main.go << 'EOF'
package main

import "fmt"

// BROKEN: uses == inside a generic function but constrains T with `any`,
// which does not guarantee the values support ==.
func Contains[T any](s []T, v T) bool {
    for _, item := range s {
        if item == v {
            return true
        }
    }
    return false
}

func main() {
    fmt.Println(Contains([]int{1, 2, 3}, 2))
}
EOF

go build .
```

**Expected Output (real compiler error):**

```
# level25.example/exercise4broken
./main.go:9:6: invalid operation: item == v (incomparable types in type set)
```

**Learning Objectives:**
- ✅ Write a generic `Contains`/`IndexOf` constrained with `comparable`
- ✅ Use a type parameter as a `map` key type
- ✅ Trigger and read the real compile error from forgetting `comparable`

---

## Exercise 5: A Custom Number Constraint

**Objective:** Define your own constraint locally and use it for Sum and Average

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level25-exercise5
cd ~/projects/level25-exercise5
go mod init level25.example/exercise5
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

// Number is a custom constraint defined locally - no external packages needed.
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

func main() {
    ints := []int{1, 2, 3, 4, 5}
    int64s := []int64{100, 200, 300}
    floats := []float64{1.5, 2.5, 3.0}

    fmt.Println("=== Sum Across Different Numeric Types ===")
    fmt.Println("Sum(ints):", Sum(ints))
    fmt.Println("Sum(int64s):", Sum(int64s))
    fmt.Println("Sum(floats):", Sum(floats))

    fmt.Println("\n=== Average Across Different Numeric Types ===")
    fmt.Printf("Average(ints): %.2f\n", Average(ints))
    fmt.Printf("Average(int64s): %.2f\n", Average(int64s))
    fmt.Printf("Average(floats): %.2f\n", Average(floats))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Sum Across Different Numeric Types ===
Sum(ints): 15
Sum(int64s): 600
Sum(floats): 7

=== Average Across Different Numeric Types ===
Average(ints): 3.00
Average(int64s): 200.00
Average(floats): 2.33
```

4. Now trigger the real compile error from violating the constraint. In a **separate** scratch directory:

```bash
mkdir -p ~/projects/level25-exercise5-broken
cd ~/projects/level25-exercise5-broken
go mod init level25.example/exercise5broken

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
    // BROKEN: string does not satisfy the Number constraint's type set.
    words := []string{"a", "b", "c"}
    fmt.Println(Sum(words))
}
EOF

go build .
```

**Expected Output (real compiler error):**

```
# level25.example/exercise5broken
./main.go:20:17: string does not satisfy Number (string missing in int | int64 | float64)
```

**Learning Objectives:**
- ✅ Define a custom constraint interface with a type set (`int | int64 | float64`)
- ✅ Use it to constrain both `Sum` and `Average`
- ✅ Trigger and read the real compile error from passing a type outside the type set

---

## Exercise 6: Map, Filter, Reduce - Revisited From Level 11

**Objective:** Rewrite Level 11's higher-order helpers as truly generic functions

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level25-exercise6
cd ~/projects/level25-exercise6
go mod init level25.example/exercise6
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

// Map, Filter, and Reduce - generic rewrites of the Level 11 higher-order
// helpers, which only worked on []int before generics existed.

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

func main() {
    nums := []int{1, 2, 3, 4, 5, 6}

    fmt.Println("=== Map: int -> int ===")
    doubled := Map(nums, func(n int) int { return n * 2 })
    fmt.Println(doubled)

    fmt.Println("\n=== Map: int -> string (different output type) ===")
    labels := Map(nums, func(n int) string { return fmt.Sprintf("n%d", n) })
    fmt.Println(labels)

    fmt.Println("\n=== Filter: keep evens ===")
    evens := Filter(nums, func(n int) bool { return n%2 == 0 })
    fmt.Println(evens)

    fmt.Println("\n=== Filter: works on strings too ===")
    words := []string{"go", "is", "fun", "and", "fast"}
    longWords := Filter(words, func(w string) bool { return len(w) > 2 })
    fmt.Println(longWords)

    fmt.Println("\n=== Reduce: sum ===")
    sum := Reduce(nums, 0, func(acc, n int) int { return acc + n })
    fmt.Println(sum)

    fmt.Println("\n=== Reduce: build a string (different accumulator type) ===")
    joined := Reduce(words, "", func(acc string, w string) string {
        if acc == "" {
            return w
        }
        return acc + "-" + w
    })
    fmt.Println(joined)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Map: int -> int ===
[2 4 6 8 10 12]

=== Map: int -> string (different output type) ===
[n1 n2 n3 n4 n5 n6]

=== Filter: keep evens ===
[2 4 6]

=== Filter: works on strings too ===
[fun and fast]

=== Reduce: sum ===
21

=== Reduce: build a string (different accumulator type) ===
go-is-fun-and-fast
```

**Learning Objectives:**
- ✅ Rewrite `Map`/`Filter`/`Reduce` as truly generic functions, not type-specific ones
- ✅ Use two type parameters (`T`, `U`) where input and output/accumulator types differ
- ✅ Reuse the same three functions across `int` and `string` element types

---

## Exercise 7: A Generic Stack[T]

**Objective:** Build a generic data structure with methods

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level25-exercise7
cd ~/projects/level25-exercise7
go mod init level25.example/exercise7
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

// Stack is a generic LIFO data structure that works for any element type.
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

func (s *Stack[T]) Len() int {
    return len(s.items)
}

func main() {
    fmt.Println("=== Stack[int] ===")
    var ints Stack[int]
    ints.Push(1)
    ints.Push(2)
    ints.Push(3)
    fmt.Println("Len:", ints.Len())

    top, _ := ints.Peek()
    fmt.Println("Peek:", top)

    for ints.Len() > 0 {
        v, _ := ints.Pop()
        fmt.Println("Pop:", v)
    }

    _, ok := ints.Pop()
    fmt.Println("Pop on empty, ok =", ok)

    fmt.Println("\n=== Stack[string] ===")
    var words Stack[string]
    words.Push("go")
    words.Push("generics")
    words.Push("rock")
    for words.Len() > 0 {
        v, _ := words.Pop()
        fmt.Println("Pop:", v)
    }
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Stack[int] ===
Len: 3
Peek: 3
Pop: 3
Pop: 2
Pop: 1
Pop on empty, ok = false

=== Stack[string] ===
Pop: rock
Pop: generics
Pop: go
```

**Learning Objectives:**
- ✅ Declare a generic struct type with `[T any]`
- ✅ Repeat the type parameter on every method receiver (`func (s *Stack[T]) Pop()`)
- ✅ Instantiate the same type as `Stack[int]` and `Stack[string]` independently

---

## Exercise 8: A Generic Pair[K, V]

**Objective:** Build a generic struct holding two independently-typed values

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level25-exercise8
cd ~/projects/level25-exercise8
go mod init level25.example/exercise8
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

// Pair holds two related values of independent types.
type Pair[K, V any] struct {
    Key   K
    Value V
}

func NewPair[K, V any](key K, value V) Pair[K, V] {
    return Pair[K, V]{Key: key, Value: value}
}

func (p Pair[K, V]) String() string {
    return fmt.Sprintf("(%v, %v)", p.Key, p.Value)
}

// Swap returns a new Pair with Key and Value swapped - only valid when
// K and V are the same type, so it needs its own type parameter T.
func Swap[T any](p Pair[T, T]) Pair[T, T] {
    return Pair[T, T]{Key: p.Value, Value: p.Key}
}

func main() {
    fmt.Println("=== Pair[string, int] ===")
    age := NewPair("Alice", 30)
    fmt.Println(age)
    fmt.Println("Key:", age.Key, "Value:", age.Value)

    fmt.Println("\n=== Pair[int, bool] ===")
    flag := NewPair(1, true)
    fmt.Println(flag)

    fmt.Println("\n=== Slice of Pairs ===")
    pairs := []Pair[string, int]{
        NewPair("a", 1),
        NewPair("b", 2),
        NewPair("c", 3),
    }
    for _, p := range pairs {
        fmt.Println(p)
    }

    fmt.Println("\n=== Swap (same-type Pair) ===")
    same := NewPair("first", "second")
    fmt.Println("before:", same)
    fmt.Println("after: ", Swap(same))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Pair[string, int] ===
(Alice, 30)
Key: Alice Value: 30

=== Pair[int, bool] ===
(1, true)

=== Slice of Pairs ===
(a, 1)
(b, 2)
(c, 3)

=== Swap (same-type Pair) ===
before: (first, second)
after:  (second, first)
```

**Learning Objectives:**
- ✅ Declare a generic struct with two independent type parameters (`K`, `V`)
- ✅ Implement `fmt.Stringer` on a generic type
- ✅ Recognize when a helper needs `Pair[T, T]` (same type twice) instead of `Pair[K, V]`

---

## Exercise 9: Type Inference - When It Works and When It Fails

**Objective:** Explore inference success and the one case where it genuinely fails

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level25-exercise9
cd ~/projects/level25-exercise9
go mod init level25.example/exercise9
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func Map[T, U any](s []T, f func(T) U) []U {
    result := make([]U, len(s))
    for i, v := range s {
        result[i] = f(v)
    }
    return result
}

// Zero's type parameter appears ONLY in the return type - there is no
// argument for the compiler to infer T from.
func Zero[T any]() T {
    var zero T
    return zero
}

// MakeSlice has the same problem: T only appears in the return type.
func MakeSlice[T any](n int) []T {
    return make([]T, n)
}

func main() {
    fmt.Println("=== Inference Succeeds: T Appears in a Parameter ===")
    nums := []int{1, 2, 3}
    // T and U are both inferred from the arguments - no need to write Map[int, string]
    labels := Map(nums, func(n int) string { return fmt.Sprintf("#%d", n) })
    fmt.Println(labels)

    fmt.Println("\n=== Same Call, With Explicit Type Arguments (also valid) ===")
    labels2 := Map[int, string](nums, func(n int) string { return fmt.Sprintf("#%d", n) })
    fmt.Println(labels2)

    fmt.Println("\n=== Inference FAILS: T Appears Only in the Return Type ===")
    // zero := Zero() // would NOT compile: cannot infer T (declared at ./main.go:5:11)
    zeroInt := Zero[int]()
    zeroString := Zero[string]()
    fmt.Println("Zero[int]():", zeroInt)
    fmt.Printf("Zero[string](): %q\n", zeroString)

    fmt.Println("\n=== Same Problem With MakeSlice ===")
    s := MakeSlice[float64](3)
    fmt.Println("MakeSlice[float64](3):", s)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Inference Succeeds: T Appears in a Parameter ===
[#1 #2 #3]

=== Same Call, With Explicit Type Arguments (also valid) ===
[#1 #2 #3]

=== Inference FAILS: T Appears Only in the Return Type ===
Zero[int](): 0
Zero[string](): ""

=== Same Problem With MakeSlice ===
MakeSlice[float64](3): [0 0 0]
```

4. Now trigger the real compile error by actually calling `Zero()` with no type argument. In a **separate** scratch directory:

```bash
mkdir -p ~/projects/level25-exercise9-broken
cd ~/projects/level25-exercise9-broken
go mod init level25.example/exercise9broken

cat > main.go << 'EOF'
package main

import "fmt"

func Zero[T any]() T {
    var zero T
    return zero
}

func main() {
    zero := Zero()
    fmt.Println(zero)
}
EOF

go build .
```

**Expected Output (real compiler error):**

```
# level25.example/exercise9broken
./main.go:11:14: in call to Zero, cannot infer T (declared at ./main.go:5:11)
```

**Learning Objectives:**
- ✅ Confirm inference works when a type parameter appears in a regular argument
- ✅ Confirm explicit type arguments always remain valid, even when inference would succeed
- ✅ Trigger and read the real "cannot infer T" compile error when a type parameter appears only in the return type

---

## Exercise 10: Comprehensive Practice - A Generic Cache[K, V]

**Objective:** Combine constraints, generic types, and methods into one realistic data structure

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level25-exercise10
cd ~/projects/level25-exercise10
go mod init level25.example/exercise10
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

// Cache is a generic in-memory key-value store. K must be comparable
// because it's used as a map key; V can be anything.
type Cache[K comparable, V any] struct {
    data map[K]V
}

func NewCache[K comparable, V any]() *Cache[K, V] {
    return &Cache[K, V]{data: make(map[K]V)}
}

func (c *Cache[K, V]) Set(key K, value V) {
    c.data[key] = value
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
    v, ok := c.data[key]
    return v, ok
}

func (c *Cache[K, V]) Delete(key K) {
    delete(c.data, key)
}

func (c *Cache[K, V]) Len() int {
    return len(c.data)
}

type User struct {
    Name string
    Age  int
}

func main() {
    fmt.Println("=== Cache[string, int] ===")
    scores := NewCache[string, int]()
    scores.Set("Alice", 95)
    scores.Set("Bob", 87)
    fmt.Println("Len:", scores.Len())

    if v, ok := scores.Get("Alice"); ok {
        fmt.Println("Alice:", v)
    }

    scores.Delete("Alice")
    if _, ok := scores.Get("Alice"); !ok {
        fmt.Println("Alice deleted, ok =", ok)
    }

    fmt.Println("\n=== Cache[int, User] (struct values) ===")
    users := NewCache[int, User]()
    users.Set(1, User{Name: "Charlie", Age: 28})
    users.Set(2, User{Name: "Dana", Age: 34})

    if u, ok := users.Get(1); ok {
        fmt.Printf("User 1: %+v\n", u)
    }

    if _, ok := users.Get(99); !ok {
        fmt.Println("User 99 not found, ok =", ok)
    }

    fmt.Println("Final cache size:", users.Len())
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Cache[string, int] ===
Len: 2
Alice: 95
Alice deleted, ok = false

=== Cache[int, User] (struct values) ===
User 1: {Name:Charlie Age:28}
User 99 not found, ok = false
Final cache size: 2
```

**Learning Objectives:**
- ✅ Combine `comparable` (for the key) and `any` (for the value) in one generic type
- ✅ Build a constructor generic function (`NewCache[K, V]`) alongside a generic type
- ✅ Confirm the same `Cache` type works with primitive and struct value types

---

## Bonus Challenges

### Challenge 1: Generic Binary Search

Write `BinarySearch[T cmp.Ordered](s []T, target T) (int, bool)` using the standard library's `cmp.Ordered` constraint (any type supporting `<`, `<=`, `>`, `>=`) against a sorted slice.

```bash
mkdir -p ~/projects/level25-bonus1
cd ~/projects/level25-bonus1
go mod init level25.example/bonus1
```

**Hints:**
- Import `"cmp"` - `cmp.Ordered` is a real, offline, standard-library constraint (verified against this course's Go toolchain)
- Track `lo, hi := 0, len(s)-1` and compare `s[mid]` against `target` each iteration
- Return `(-1, false)` when `lo > hi` without finding a match

### Challenge 2: Generic Linked List

Build a generic singly linked list: `type Node[T any] struct { Value T; Next *Node[T] }` and `type List[T any] struct { Head *Node[T] }` with a `PushFront` method and a `ToSlice` method for inspection.

```bash
mkdir -p ~/projects/level25-bonus2
cd ~/projects/level25-bonus2
go mod init level25.example/bonus2
```

**Hints:**
- `PushFront` just needs to set `l.Head = &Node[T]{Value: v, Next: l.Head}`
- `ToSlice` walks `for n := l.Head; n != nil; n = n.Next`

### Challenge 3: Generic Zip

Write `Zip[T, U any](a []T, b []U) []Pair[T, U]` that combines two slices into a slice of pairs, stopping at the shorter slice's length.

```bash
mkdir -p ~/projects/level25-bonus3
cd ~/projects/level25-bonus3
go mod init level25.example/bonus3
```

**Hints:**
- Reuse the `Pair[K, V]` type from Exercise 8
- Compute `n := min(len(a), len(b))` (by hand or with the `min` builtin) before looping

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Explain the exact duplication-vs-type-safety trade-off generics resolve
✅ Declare and call generic functions with and without explicit type arguments
✅ Use the built-in `any` and `comparable` constraints correctly
✅ Define your own constraints locally with type sets, fully offline
✅ Rewrite Level 11's `Map`/`Filter`/`Reduce` as truly generic functions
✅ Build generic types (`Stack[T]`, `Pair[K, V]`, `Cache[K, V]`) with methods
✅ Predict when type inference works and when you must specify type arguments
✅ Read and understand real generics-related compiler errors

---

## Next Level

Level 26: Testing
- Writing table-driven tests with the `testing` package
- `go test`, subtests, and test coverage
- Testing generic functions across multiple type instantiations
- Mocking and test doubles using the small-interface patterns from Level 15

Great work! You're mastering Go! 🚀
