# Level 3: Quick Reference Card

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir myapp && cd myapp
go mod init myapp

# Create main.go
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    // Declaration methods
    var name string = "Alice"  // Method 1: explicit
    var age = 25               // Method 2: inferred
    city := "NYC"              // Method 3: short (functions only)
    
    // Constants
    const pi = 3.14159
    
    // Display
    fmt.Printf("Name: %s, Age: %d, City: %s, Pi: %.2f\n", name, age, city, pi)
}
EOF

# Run
go run main.go
```

---

## 📋 Essential Syntax

### Variable Declaration
```go
var x int = 5              // Explicit type
var x = 5                  // Type inference
x := 5                     // Short (functions only!)

var (
    a int
    b string
)                          // Multiple

x, y := 10, 20            // Multiple short
```

### Constants
```go
const pi = 3.14159
const maxUsers int = 100

const (
    Sun = iota     // 0
    Mon            // 1
    Tue            // 2
)
```

### Named Types
```go
type Age int
type Temperature float64

age := Age(25)
temp := Temperature(98.6)
```

---

## 🔢 Basic Types Reference

### Integers
```go
int, int8, int16, int32, int64    // Signed
uint, uint8, uint16, uint32, uint64  // Unsigned

// Default: int (platform dependent)
var count int = 42
```

### Floats
```go
float32  // 32-bit
float64  // 64-bit (default for decimal literals)

var pi float64 = 3.14159
```

### String, Bool, Byte, Rune
```go
string  // Text (immutable)
bool    // true or false
byte    // uint8 (ASCII)
rune    // int32 (Unicode)

var greeting string = "Hello"
var isActive bool = true
var char byte = 65  // 'A'
```

---

## 🔄 Type Conversion Cheat Sheet

### Numeric Conversions
```go
float64(42)           // int to float
int(3.14)             // float to int (loses decimal)
uint(42)              // int to uint
int64(42)             // int to int64
```

### String Conversions
```go
import "strconv"

// String to int
n, err := strconv.Atoi("42")      // Returns int

// Int to string
s := strconv.Itoa(42)             // Returns string

// String to float
f, err := strconv.ParseFloat("3.14", 64)

// Float to string
s := strconv.FormatFloat(3.14, 'f', 2, 64)

// String ↔ Bytes
b := []byte("hello")              // string to []byte
s := string(b)                    // []byte to string
```

---

## 💾 Zero Values

| Type | Zero Value |
|------|------------|
| int, int64, etc. | 0 |
| float32, float64 | 0.0 |
| string | "" |
| bool | false |
| pointer | nil |

```go
var i int           // 0
var f float64       // 0.0
var s string        // ""
var b bool          // false
```

---

## 🎯 Declaration Quick Guide

```
Inside function?
├─ YES → Use x := 5
└─ NO → Use var x = 5 or var x int = 5

Type obvious?
├─ YES → Use type inference
└─ NO → Be explicit

Multiple variables?
├─ YES → Use var() block or x, y := 1, 2
└─ NO → Single declaration
```

---

## 🔧 Common Operations

### String Operations
```go
import "strings"

s := "Hello, World!"

len(s)                           // 13
s[0]                            // 72 (byte)
s[0:5]                          // "Hello"

strings.ToUpper(s)              // "HELLO, WORLD!"
strings.ToLower(s)              // "hello, world!"
strings.Contains(s, "World")    // true
strings.Replace(s, "World", "Go", 1)  // "Hello, Go!"
```

### Integer Operations
```go
a, b := 10, 3

a + b       // 13
a - b       // 7
a * b       // 30
a / b       // 3 (integer division)
a % b       // 1 (modulo)
a == b      // false
a < b       // false
```

### Boolean Operations
```go
true && false   // false (AND)
true || false   // true (OR)
!true           // false (NOT)
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Mixing types | `f := i` | `f := float64(i)` |
| := at package | `x := 5` | `var x = 5` |
| Re-declaring | `x := 1; x := 2` | `x := 1; x = 2` |
| String index | `s[0] == 'H'` | `string(s[0]) == "H"` |
| Forgetting conversion | `int(str)` | `strconv.Atoi(str)` |

---

## 📚 Named Types Pattern

```go
// Define
type UserId int
type Score float64

// Use
u := UserId(42)
s := Score(95)

// Add methods
func (s Score) Grade() string {
    if s >= 90 { return "A" }
    if s >= 80 { return "B" }
    return "C"
}

// Call
fmt.Println(s.Grade())  // "A"
```

---

## 🌳 Type Decision Tree

```
Need a number?
├─ For counting (default) → int
├─ For decimals → float64
├─ For exact range → int32, int64, float32
└─ For positive only → uint

Need text?
├─ Whole text → string
├─ Individual char → rune (or byte)
└─ Binary data → []byte

Need grouping?
├─ Fixed size → [size]Type (array)
├─ Dynamic size → []Type (slice)
├─ Key-value → map[K]V
└─ Multiple fields → struct
```

---

## 📋 Composite Types Preview

### Arrays (Fixed)
```go
colors := [3]string{"Red", "Green", "Blue"}
colors[0] = "Yellow"
len(colors)  // 3
```

### Slices (Dynamic)
```go
nums := []int{1, 2, 3}
nums = append(nums, 4)
nums[0]
len(nums)
cap(nums)
```

### Maps (Key-Value)
```go
scores := map[string]int{"Alice": 95, "Bob": 87}
scores["Charlie"] = 92
delete(scores, "Bob")
val, ok := scores["Alice"]  // Check exists
```

### Structs
```go
type Person struct {
    Name string
    Age  int
}

p := Person{Name: "Alice", Age: 25}
p.Age = 26
```

---

## 🎓 Before Next Level

Can you:
- [ ] Declare variables three ways?
- [ ] List Go's basic types?
- [ ] Convert between types?
- [ ] Explain zero values?
- [ ] Create named types?
- [ ] Use constants?

If YES → You're ready for Level 4!

---

## 📚 Next Level

Level 4: Operators & Expressions
- Arithmetic: +, -, *, /, %
- Comparison: ==, !=, <, >, <=, >=
- Logical: &&, ||, !
- Bitwise: &, |, ^, <<, >>

You've got the type system! 💪
