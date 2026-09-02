# Level 3: Variables & Data Types - Complete Guide

## Introduction

Welcome to Level 3! You've mastered Go's program structure (packages and organization). Now it's time to learn about **variables and data types** - the building blocks of all Go programs.

In Go, every variable has:
- A **name** (identifier)
- A **type** (what kind of data it holds)
- A **value** (the actual data)

This level teaches you how to declare variables, understand Go's basic types, convert between types, and use constants effectively.

---

## Table of Contents

1. [Variable Declaration](#variable-declaration)
2. [Basic Data Types](#basic-data-types)
3. [Zero Values](#zero-values)
4. [Type Conversion](#type-conversion)
5. [Constants](#constants)
6. [Named Types](#named-types)
7. [Composite Types Introduction](#composite-types-introduction)
8. [Best Practices](#best-practices)
9. [Common Mistakes](#common-mistakes)
10. [Troubleshooting](#troubleshooting)

---

## Variable Declaration

### Method 1: var with Type

```go
var name string
var age int
var price float64

// Initialize
var name string = "Alice"
var age int = 25
var price float64 = 19.99
```

**Use when:**
- Declaring package-level variables
- Declaring multiple variables
- Want explicit type

### Method 2: Short Declaration (:=)

```go
name := "Alice"
age := 25
price := 19.99

// Multiple variables
x, y := 10, 20
```

**Use when:**
- Inside functions
- Want type inference
- Declaring and initializing together

**Important:** `:=` only works inside functions!

```go
// ✅ Inside function
func main() {
    name := "Alice"  // OK
}

// ❌ Package level - WRONG!
name := "Alice"  // Compiler error!

// ✅ Package level - right way
var name = "Alice"  // or: var name string = "Alice"
```

### Method 3: var with Type Inference

```go
var name = "Alice"      // Type inferred from value
var age = 25            // int
var price = 19.99       // float64
```

**Useful when:**
- Want to be explicit about declaration
- Type is obvious from value
- At package level

### Method 4: Multiple Variables

```go
// Method 1: Multiple declarations
var a, b, c int = 1, 2, 3

// Method 2: Type once
var (
    name string
    age int
    score float64
)

// Method 3: Short declaration (function only)
x, y, z := 10, 20, 30
```

### Method 5: Blank Identifier

```go
// When you don't need a value
_, err := someFunction()  // Ignore error

// In loops
for _, value := range mySlice {
    fmt.Println(value)
}
```

### Reassignment vs Declaration

```go
// Declaration (first time)
name := "Alice"

// Reassignment (subsequent)
name = "Bob"

// ❌ Can't re-declare with :=
name := "Charlie"  // Error! Already declared
```

---

## Basic Data Types

### Integer Types

Go has several integer types:

```go
// Signed integers
var a int8 = -128      // Range: -128 to 127
var b int16 = -32768   // Range: -32768 to 32767
var c int32 = -2147483648
var d int64 = -9223372036854775808
var e int = 42         // Size depends on platform (32 or 64 bit)

// Unsigned integers
var f uint8 = 255      // Range: 0 to 255
var g uint16 = 65535
var h uint32 = 4294967295
var i uint64 = 18446744073709551615
var j uint = 42        // Unsigned, platform dependent

// Common usage
var count int = 0      // Most common for counting
var id uint = 1        // For positive-only values
var tiny int8 = 1      // When memory matters
```

**Which to use?**
- `int` for most cases (platform optimal)
- `uint` when you know value is positive
- `int32`, `int64` when you need specific size
- `int8`, `uint8` when memory is critical

### Floating-Point Types

```go
var a float32 = 3.14          // 32-bit float
var b float64 = 3.14159265    // 64-bit float (more precise)
var c = 3.14                  // Defaults to float64

// Infinity and NaN
var inf float64 = math.Inf(1)   // Positive infinity
var nan float64 = math.NaN()    // Not a number
```

**When to use:**
- `float64` for most calculations (default)
- `float32` when memory is critical
- Both have precision limitations

### String Type

```go
// String literals
var s1 string = "Hello"
var s2 string = "World"
s3 := "Go Programming"

// Empty string (zero value)
var empty string  // ""

// Multi-line strings (raw strings)
raw := `This is a
multi-line
string`

// String with escape sequences
escaped := "Line 1\nLine 2\tTabbed"

// String concatenation
greeting := "Hello" + " " + "World"

// String length
len(greeting)  // Returns 11

// Access character (returns byte/rune)
greeting[0]    // Returns 72 ('H')
```

**String operations:**
```go
// Import "strings" package for these

import "strings"

s := "Hello, World!"

strings.Contains(s, "World")      // true
strings.Count(s, "l")             // 3
strings.HasPrefix(s, "Hello")     // true
strings.ToUpper(s)                // "HELLO, WORLD!"
strings.ToLower(s)                // "hello, world!"
strings.Replace(s, "World", "Go", 1)  // "Hello, Go!"
strings.Split(s, ",")             // ["Hello" " World!"]
strings.Join([]string{"a", "b"}, "-")  // "a-b"
```

### Boolean Type

```go
var flag bool = true
var isValid bool = false

// Zero value is false
var empty bool  // false

// Boolean operations
a := true
b := false

a && b  // AND - false
a || b  // OR - true
!a      // NOT - false

// Common use
if isValid {
    fmt.Println("Valid!")
}
```

### Byte and Rune Types

```go
// Byte (alias for uint8)
var b byte = 65  // ASCII code for 'A'

// Rune (alias for int32)
var r rune = 65  // ASCII/Unicode code point

// Character literals
b := 'A'  // Rune type (not string!)
r := '€'  // Rune (supports Unicode)

// Convert character to string
s := string(r)  // "€"
```

---

## Zero Values

When you declare a variable without initializing it, Go assigns a **zero value**:

```go
var i int       // 0
var f float64   // 0.0
var s string    // "" (empty string)
var b bool      // false
var p *int      // nil (pointer)
```

**Why this matters:**
- Variables are always initialized
- No "undefined" state
- Predictable behavior

```go
func main() {
    var count int
    count += 1  // OK! count was 0
    fmt.Println(count)  // 1
}
```

---

## Type Conversion

Go does **not** allow automatic type conversion. You must convert explicitly:

```go
// ❌ WRONG - automatic conversion not allowed
var i int = 42
var f float64 = i  // Compiler error!

// ✅ RIGHT - explicit conversion
var i int = 42
var f float64 = float64(i)  // OK
```

### Converting Between Numeric Types

```go
var i int = 42
var f float64 = float64(i)      // int to float64
var u uint = uint(i)            // int to uint
var i8 int8 = int8(i)           // int to int8

// Back to int
var back int = int(f)           // float64 to int (loses decimal part!)
```

### Converting Between Strings and Numbers

```go
import "strconv"

// String to int
s := "42"
i, err := strconv.Atoi(s)  // Returns int and error
// i = 42

// String to int64
i64, err := strconv.ParseInt("42", 10, 64)  // Base 10, 64-bit
// i64 = 42

// String to float
f, err := strconv.ParseFloat("3.14", 64)
// f = 3.14

// Int to string
s := strconv.Itoa(42)  // "42"

// Any value to string
s := fmt.Sprintf("%v", 42)  // "42"
```

### Converting Strings to Bytes and Vice Versa

```go
// String to byte slice
s := "Hello"
b := []byte(s)  // [72 101 108 108 111]

// Byte slice to string
b := []byte{72, 101, 108, 108, 111}
s := string(b)  // "Hello"

// String to rune slice
s := "Hello"
r := []rune(s)  // [72 101 108 108 111]

// Rune slice to string
r := []rune{72, 101, 108, 108, 111}
s := string(r)  // "Hello"
```

---

## Constants

Constants are like variables but **cannot be changed** after declaration:

```go
// Declare constant
const pi = 3.14159
const maxConnections int = 100

// Multiple constants
const (
    Monday = 1
    Tuesday = 2
    Wednesday = 3
)

// Typed constants
const greeting string = "Hello"

// Untyped constants
const answer = 42  // Type determined by context
```

### Using Constants

```go
// No zero value needed - must initialize
const x = 10
x = 20  // ❌ Error! Constants are immutable

// Can use in calculations
const a = 5
const b = a * 2  // OK - constants calculated at compile time
```

### The iota Keyword

Useful for creating sequences:

```go
const (
    Sunday = iota    // 0
    Monday           // 1
    Tuesday          // 2
    Wednesday        // 3
    Thursday         // 4
    Friday           // 5
    Saturday         // 6
)

// With bit shifting (for flags)
const (
    Read   = 1 << iota  // 1 (2^0)
    Write              // 2 (2^1)
    Execute            // 4 (2^2)
)
```

### Constant vs Variable

```go
// Variable - can change
var count int = 0
count = 5  // OK

// Constant - cannot change
const maxCount int = 100
maxCount = 50  // ❌ Error!
```

**Use constants when:**
- Value should never change
- Preventing accidental modification
- Making intent clear
- Performance matters (compile-time evaluation)

---

## Named Types

You can create custom names for existing types:

```go
// Define a named type
type Age int
type Temperature float64
type UserId string

// Declare variables with named type
var myAge Age = 25
var temp Temperature = 98.6
var id UserId = "user123"
```

### Why Named Types?

1. **Clarity** - Code is more readable
2. **Prevent mistakes** - Different types can't be mixed
3. **Add methods** - Can add methods to named types

### Type vs Type Alias

```go
// Named type (Go 1.8 and earlier style)
type UserId int
u := UserId(42)
i := int(u)  // Must convert

// Type alias (Go 1.9+)
type UserId = int
u := UserId(42)
i := u  // No conversion needed - same type!
```

**Named types in practice:**

```go
type Age int
type Score float64

// Add methods to named types
func (a Age) IsAdult() bool {
    return a >= 18
}

func (s Score) Grade() string {
    if s >= 90 {
        return "A"
    } else if s >= 80 {
        return "B"
    }
    return "C"
}

// Use them
myAge := Age(25)
if myAge.IsAdult() {
    fmt.Println("Adult")
}

score := Score(95)
fmt.Println("Grade:", score.Grade())  // "Grade: A"
```

---

## Composite Types Introduction

Go has composite types for grouping data:

### Arrays

```go
// Array of 5 integers
var arr [5]int

// Initialize with values
arr := [5]int{1, 2, 3, 4, 5}

// Or let Go count
arr := [...]int{1, 2, 3, 4, 5}  // Length is 5

// Access elements
fmt.Println(arr[0])  // 1
arr[0] = 10          // Change element
```

### Slices

Slices are like arrays but with dynamic length:

```go
// Create slice
slice := []int{1, 2, 3, 4, 5}

// Empty slice
var empty []int

// Slice from array
arr := [5]int{1, 2, 3, 4, 5}
slice := arr[1:4]  // Elements at index 1, 2, 3

// Slice operations
len(slice)           // Length
cap(slice)           // Capacity
slice = append(slice, 6)  // Add element
```

### Maps

Maps are key-value pairs:

```go
// Create map
scores := map[string]int{
    "Alice": 90,
    "Bob": 85,
}

// Access value
alice := scores["Alice"]  // 90

// Add entry
scores["Charlie"] = 95

// Delete entry
delete(scores, "Bob")

// Check if key exists
if value, ok := scores["Alice"]; ok {
    fmt.Println(value)  // 90
}
```

### Structs

Structs group related data:

```go
type Person struct {
    Name string
    Age  int
    City string
}

// Create struct
p := Person{
    Name: "Alice",
    Age:  25,
    City: "New York",
}

// Access fields
fmt.Println(p.Name)  // Alice
p.Age = 26           // Change field
```

---

## Best Practices

### 1. Use Short Declarations Inside Functions

```go
// ✅ Good
func main() {
    name := "Alice"
    age := 25
}

// ❌ Verbose
func main() {
    var name string = "Alice"
    var age int = 25
}
```

### 2. Use var at Package Level

```go
// ✅ Good
var (
    defaultTimeout = 30 * time.Second
    maxRetries = 3
)

// ❌ Can't use := at package level
```

### 3. Be Explicit About Types When Needed

```go
// ✅ Clear intent
var count int = 0
var enabled bool = true

// Can also be:
count := 0
enabled := true
```

### 4. Use Constants for Values That Don't Change

```go
// ✅ Good
const MaxConnections = 100
const DefaultPort = 8080

// ❌ Mutable variable
var maxConnections int = 100  // Can be accidentally changed
```

### 5. Use Meaningful Variable Names

```go
// ✅ Clear
userName := "Alice"
isActive := true
maxRetries := 3

// ❌ Unclear
u := "Alice"
ia := true
mr := 3
```

---

## Common Mistakes

### Mistake 1: Automatic Type Conversion

```go
// ❌ WRONG
var i int = 42
var f float64 = i  // Compiler error!

// ✅ RIGHT
var f float64 = float64(i)
```

### Mistake 2: Using := at Package Level

```go
// ❌ WRONG
package main
count := 42  // Syntax error!

// ✅ RIGHT
package main
var count = 42
```

### Mistake 3: Re-declaring with :=

```go
// ❌ WRONG
name := "Alice"
name := "Bob"  // Error! Already declared

// ✅ RIGHT
name := "Alice"
name = "Bob"   // Reassignment
```

### Mistake 4: Forgetting Error Handling

```go
// ❌ WRONG
i, _ := strconv.Atoi("not a number")  // Panic!

// ✅ RIGHT
i, err := strconv.Atoi("42")
if err != nil {
    log.Fatal(err)
}
```

### Mistake 5: Confusing float32 and float64

```go
// ❌ WRONG
var f float32 = 3.14
var d float64 = f  // Compiler error!

// ✅ RIGHT
var d float64 = float64(f)
```

### Mistake 6: String Indexing Returns Byte/Rune

```go
// ❌ WRONG (assumes string index returns string)
s := "Hello"
c := s[0]  // Returns byte (72), not 'H'
fmt.Println(c)  // Prints: 72

// ✅ RIGHT
s := "Hello"
c := string(s[0])  // Convert to string: "H"
fmt.Println(c)     // Prints: H
```

---

## Troubleshooting

### "Cannot use X (type int) as type int64 in assignment"

**Problem:** Tried to assign int to int64 without conversion

**Solution:** Convert explicitly
```go
var i int = 42
var i64 int64 = int64(i)  // Correct
```

### "Undefined: strconv.Atoi"

**Problem:** Forgot to import strconv package

**Solution:** Add import
```go
import "strconv"
```

### "Short variable declarations only allowed inside function"

**Problem:** Used `:=` at package level

**Solution:** Use `var` instead
```go
// Package level
var count = 42

// Inside function
count := 42
```

### "String indexing produces a value of type byte, not rune"

**Problem:** Forgot that string indexing returns a byte value

**Solution:** Convert appropriately
```go
s := "Hello"
b := s[0]           // Returns byte (72)
c := string(s[0])   // Convert to string ("H")
r := []rune(s)[0]   // Get as rune ('H')
```

### "Cannot convert X to Y"

**Problem:** Tried to convert between incompatible types (e.g., string to int)

**Solution:** Use appropriate conversion functions
```go
s := "42"
i := strconv.Atoi(s)  // Use Atoi for string to int
```

---

## Summary

**Key Takeaways:**

1. **Three ways to declare:** `var`, `:=` (inside functions), `const`
2. **Basic types:** int, float64, string, bool, byte, rune
3. **Zero values:** Variables always initialized to zero value of their type
4. **Type conversion:** Must be explicit (no automatic conversion)
5. **Constants:** Use for unchangeable values
6. **Named types:** Create custom type names for clarity
7. **Composite types:** Arrays, slices, maps, structs (more detail in later levels)

---

## Next Steps

You now understand:
- ✅ Variable declaration
- ✅ Go's type system
- ✅ Type conversion
- ✅ Constants
- ✅ Named types

**Next level:** Level 4 - Operators & Expressions
- Arithmetic, comparison, logical operators
- Operator precedence
- Expression evaluation

You're building solid Go fundamentals! Keep going! 🚀
