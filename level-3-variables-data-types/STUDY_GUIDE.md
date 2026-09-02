# Level 3: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: Variable Declaration & Basic Types
```
Day 1:  Variable declaration styles
Day 2:  Integer and float types
Day 3:  Strings and booleans
Day 4:  Zero values
Day 5:  Type conversion
Day 6:  Constants and iota
Day 7:  Named types
```

### Week 2: Practice & Application
```
Day 1:  Exercises 1-3 (Declarations)
Day 2:  Exercises 4-6 (Conversion, constants, named types)
Day 3:  Exercises 7-8 (Multiple variables, type safety)
Day 4:  Exercises 9-10 (Composite types)
Day 5:  Bonus challenges
Day 6-7: Practice & review
```

---

## 🎯 Variable Declaration Comparison

### Three Methods

| Method | Syntax | Where | Use Case |
|--------|--------|-------|----------|
| **var with type** | `var x int = 5` | Anywhere | Package level, explicit type |
| **var with inference** | `var x = 5` | Anywhere | Type obvious from value |
| **Short (:=)** | `x := 5` | Function only | Inside functions, quick |

### Decision Tree

```
Want to declare a variable?
├─ Inside a function?
│  └─ Use := (x := 5)
├─ At package level?
│  ├─ Type obvious?
│  │  └─ Use var with inference (var x = 5)
│  ├─ Type NOT obvious?
│  │  └─ Use var with type (var x int = 5)
└─ With explicit type always?
   └─ Use var with type
```

---

## 🔢 Go Type System

### Integer Types

```
Signed Integers:
int8   → -128 to 127
int16  → -32,768 to 32,767
int32  → -2,147,483,648 to 2,147,483,647
int64  → -9,223,372,036,854,775,808 to 9,223,372,036,854,775,807
int    → Platform dependent (32 or 64 bit)

Unsigned Integers:
uint8  → 0 to 255
uint16 → 0 to 65,535
uint32 → 0 to 4,294,967,295
uint64 → 0 to 18,446,744,073,709,551,615
uint   → Platform dependent (32 or 64 bit)

Most Common:
int    → For general counting/indexing
uint   → For positive values only
```

### Floating-Point Types

```
float32 → 32-bit (≈7 digits precision)
float64 → 64-bit (≈15 digits precision) ← Default

Use float64 unless memory is critical
```

### String, Byte, Rune

```
string → Text (immutable)
         len() returns byte count
         s[i] returns byte (not character)

byte   → uint8 (0-255)
         ASCII values

rune   → int32 (Unicode)
         Supports all Unicode characters

Important: s[0] returns byte, not string!
To get character: string(s[0]) or []rune(s)[0]
```

---

## 🔄 Type Conversion Flow

```
int ──→ float64: float64(x)
float64 → int: int(x)  [loses decimal]

string ──→ int: strconv.Atoi(s)
int ────→ string: strconv.Itoa(x)

string ──→ float: strconv.ParseFloat(s, 64)
float ──→ string: strconv.FormatFloat(x, 'f', 2, 64)

string ──→ []byte: []byte(s)
[]byte ──→ string: string(b)

string ──→ []rune: []rune(s)
[]rune ──→ string: string(r)
```

### Remember:
- ❌ NO automatic conversion
- ✅ MUST use explicit conversion
- ❌ Different types can't be mixed

```go
// Error!
var i int = 42
var f float64 = i  // Compiler error

// Correct!
var f float64 = float64(i)
```

---

## 📋 Zero Values Table

| Type | Zero Value |
|------|------------|
| int, int8, int64, etc. | 0 |
| uint, uint8, etc. | 0 |
| float32, float64 | 0.0 |
| string | "" (empty) |
| bool | false |
| pointer | nil |
| slice | nil |
| map | nil |
| array | All elements zero |

---

## ⭐ Constants vs Variables

```
Variable:
var x = 10      Can change
x = 20          ✅ OK

Constant:
const y = 10    Cannot change
y = 20          ❌ Error!

const z = 5
const result = z * 2  ✅ OK (compile-time)
```

### When to Use Constants

```
const when:
- Value should NEVER change (MaxSize)
- Want compiler to check (catch errors early)
- Performance matters (compile-time evaluation)
- Making intent clear (this is fixed)

Example good constants:
const MaxConnections = 100
const DefaultTimeout = 30
const APIVersion = "v2"
```

---

## 🔀 iota for Sequences

```go
const (
    Sun    = iota  // 0
    Mon           // 1
    Tue           // 2
    Wed           // 3
)

// With operations:
const (
    Read   = 1 << iota  // 1   (2^0)
    Write              // 2   (2^1)
    Exec               // 4   (2^2)
)

// Skip values:
const (
    A = iota       // 0
    _              // (skipped)
    C              // 2
)
```

---

## 🏷️ Named Types

### Basic Idea

```go
// Before: just int
var age = 25
var salary = 50000

// After: clear intent
type Age int
type Salary int

var age Age = 25
var salary Salary = 50000
```

### Benefits

```
1. Type Safety
   age and salary are different types
   Can't accidentally mix them

2. Methods
   func (a Age) IsAdult() bool { return a >= 18 }
   age.IsAdult()  ✅ Works

3. Documentation
   Type name explains what the value represents
```

### Pattern: Named Type + Methods

```go
type Temperature float64

func (t Temperature) Celsius() float64 {
    return float64(t)
}

func (t Temperature) Fahrenheit() float64 {
    return float64(t)*9/5 + 32
}

func (t Temperature) Kelvin() float64 {
    return float64(t) + 273.15
}

// Use:
temp := Temperature(25)
fmt.Println(temp.Fahrenheit())  // 77
```

---

## 🧮 Composite Types Preview

### Arrays (Fixed Size)

```go
var colors [3]string
colors[0] = "Red"

scores := [5]int{90, 85, 78, 92, 88}

len(scores)  // 5 (always fixed)
```

### Slices (Dynamic Size)

```go
numbers := []int{1, 2, 3, 4, 5}

len(numbers)      // 5 (current length)
cap(numbers)      // 5 (capacity)

numbers = append(numbers, 6)  // Add element

first3 := numbers[:3]  // Slice [1, 2, 3]
```

### Maps (Key-Value)

```go
person := map[string]string{
    "name": "Alice",
    "city": "NYC",
}

person["age"] = "25"  // Add entry

name := person["name"]  // "Alice"

delete(person, "age")   // Remove entry
```

### Structs (Grouped Data)

```go
type Person struct {
    Name string
    Age  int
    City string
}

p := Person{
    Name: "Alice",
    Age: 25,
    City: "NYC",
}

p.Age = 26  // Access field
```

---

## 📊 Type Conversion Decision Tree

```
Want to convert?
├─ Number to different number type?
│  └─ Use T(value): int(f), float64(i), etc.
├─ String to number?
│  ├─ String to int?
│  │  └─ Use strconv.Atoi() or ParseInt()
│  ├─ String to float?
│  │  └─ Use strconv.ParseFloat()
├─ Number to string?
│  ├─ Int to string?
│  │  └─ Use strconv.Itoa()
│  ├─ Float to string?
│  │  └─ Use strconv.FormatFloat() or fmt
├─ String to/from bytes?
│  ├─ String to []byte?
│  │  └─ []byte(s)
│  └─ []byte to string?
│     └─ string(b)
└─ String to/from runes?
   ├─ String to []rune?
   │  └─ []rune(s)
   └─ []rune to string?
      └─ string(r)
```

---

## 🚨 Common Mistakes

### Mistake 1: Forgetting Types Must Match

```go
// ❌ WRONG
var i int = 42
var f float64 = i

// ✅ RIGHT
var f float64 = float64(i)
```

### Mistake 2: Using := at Package Level

```go
// ❌ WRONG
package main
count := 42

// ✅ RIGHT
package main
var count = 42
```

### Mistake 3: Expecting Automatic int → float

```go
// ❌ WRONG
avg := (90 + 80 + 70) / 3  // Result is int!

// ✅ RIGHT
avg := float64((90 + 80 + 70)) / 3.0
```

### Mistake 4: String Indexing Returns Byte

```go
// ❌ WRONG (gets numeric value)
s := "Hello"
c := s[0]  // 72, not 'H'

// ✅ RIGHT (convert to string)
c := string(s[0])  // "H"
```

### Mistake 5: Re-declaring with :=

```go
// ❌ WRONG
x := 5
x := 10  // Error! Already declared

// ✅ RIGHT
x := 5
x = 10   // Reassignment
```

---

## 📈 Progression Summary

### Understanding Level 3

Level 3 teaches the foundation of Go's type system:

1. **Declaration** - How to create variables
2. **Types** - What kinds of data Go supports
3. **Conversion** - How to change between types
4. **Constants** - Fixed values
5. **Named Types** - Custom type names
6. **Composite Types** - Grouping data (preview)

### Prerequisites for Level 4

Before moving to Level 4 (Operators), you need:

- ✅ Understand var, :=, and const
- ✅ Know basic types (int, float64, string, bool)
- ✅ Can do type conversion
- ✅ Understand zero values
- ✅ Can create named types

### Ready for Level 4?

Level 4 teaches operators that work WITH these types:
- Arithmetic: +, -, *, /, %
- Comparison: ==, !=, <, >, <=, >=
- Logical: &&, ||, !
- Bitwise: &, |, ^, <<, >>

---

## ✅ Checklist Before Level 4

- [ ] Can declare variables three ways
- [ ] Understand all basic types
- [ ] Can convert between types
- [ ] Know difference between var and const
- [ ] Can create named types
- [ ] Understand zero values
- [ ] Completed 8+ exercises
- [ ] Can explain types to someone else

---

## 💡 Key Takeaways

### The Three Ways
1. `var x int = 5` - Package level, explicit
2. `var x = 5` - Anywhere, inferred
3. `x := 5` - Function only, inferred

### The Rule
No automatic type conversion. ALWAYS explicit.

### The Types
int, float64, string, bool + byte, rune

### The Methods
- Numeric: T(value)
- String: strconv.Atoi/ParseFloat
- Bytes: []byte(s) / string(b)

---

## 📚 Next Level

Level 4: Operators & Expressions
- Arithmetic operations
- Comparison operators
- Logical operators
- Bitwise operations
- Operator precedence

You've mastered Go's types! Keep going! 🚀
