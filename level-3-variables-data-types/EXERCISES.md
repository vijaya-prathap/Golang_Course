# Level 3: Variables & Data Types - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

---

## Exercise 1: Variable Declaration Styles

**Objective:** Learn the three ways to declare variables in Go

**Instructions:**

1. Create a file `main.go`:

```bash
mkdir -p ~/projects/level3-exercise1
cd ~/projects/level3-exercise1
go mod init level3.example/exercise1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

// Package-level variables (using var)
var globalName string = "Alice"
var globalAge int = 25

func main() {
    // Method 1: var with type
    var firstName string = "Bob"
    var age int = 30
    var height float64 = 5.9
    
    // Method 2: var with type inference
    var city = "New York"
    var temperature = 72.5
    
    // Method 3: short declaration (only inside functions)
    lastName := "Smith"
    salary := 50000
    
    // Display values
    fmt.Println("Method 1 (var with type):")
    fmt.Println("First Name:", firstName, "Age:", age, "Height:", height)
    
    fmt.Println("\nMethod 2 (var with inference):")
    fmt.Println("City:", city, "Temperature:", temperature)
    
    fmt.Println("\nMethod 3 (short declaration):")
    fmt.Println("Last Name:", lastName, "Salary:", salary)
    
    fmt.Println("\nPackage-level variables:")
    fmt.Println("Global Name:", globalName, "Global Age:", globalAge)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
Method 1 (var with type):
First Name: Bob Age: 30 Height: 5.9

Method 2 (var with inference):
City: New York Temperature: 72.5

Method 3 (short declaration):
Last Name: Smith Salary: 50000

Package-level variables:
Global Name: Alice Global Age: 25
```

**Learning Objectives:**
- ✅ Understand three variable declaration methods
- ✅ Know when to use each method
- ✅ Understand package-level vs function-level variables

---

## Exercise 2: Basic Data Types

**Objective:** Practice using Go's basic data types

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level3-exercise2
cd ~/projects/level3-exercise2
go mod init level3.example/exercise2
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    // Integers
    var count int = 42
    var positive uint = 100
    var small int8 = -50
    
    // Floating-point
    var pi float64 = 3.14159
    var ratio float32 = 2.5
    
    // String
    var greeting string = "Hello, Go!"
    
    // Boolean
    var isActive bool = true
    var isValid bool = false
    
    // Display all types
    fmt.Println("=== Integer Types ===")
    fmt.Println("int:", count)
    fmt.Println("uint:", positive)
    fmt.Println("int8:", small)
    
    fmt.Println("\n=== Floating-Point Types ===")
    fmt.Println("float64:", pi)
    fmt.Println("float32:", ratio)
    
    fmt.Println("\n=== String Type ===")
    fmt.Println("string:", greeting)
    fmt.Println("length:", len(greeting))
    fmt.Println("first character (byte):", greeting[0])
    
    fmt.Println("\n=== Boolean Type ===")
    fmt.Println("bool true:", isActive)
    fmt.Println("bool false:", isValid)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Integer Types ===
int: 42
uint: 100
int8: -50

=== Floating-Point Types ===
float64: 3.14159
float32: 2.5

=== String Type ===
string: Hello, Go!
length: 11
first character (byte): 72

=== Boolean Type ===
bool true: true
bool false: false
```

**Learning Objectives:**
- ✅ Understand integer types (int, uint, int8, etc.)
- ✅ Use floating-point types
- ✅ Work with strings
- ✅ Use booleans
- ✅ Understand type defaults

---

## Exercise 3: Zero Values

**Objective:** Understand Go's zero value concept

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level3-exercise3
cd ~/projects/level3-exercise3
go mod init level3.example/exercise3
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    // Declare without initialization
    var i int
    var f float64
    var s string
    var b bool
    
    fmt.Println("=== Zero Values ===")
    fmt.Println("int zero value:", i)
    fmt.Println("float64 zero value:", f)
    fmt.Println("string zero value: [" + s + "]")
    fmt.Println("bool zero value:", b)
    
    // Demonstrate using zero values
    fmt.Println("\n=== Using Zero Values ===")
    var counter int  // 0
    counter += 5
    fmt.Println("counter after +=5:", counter)
    
    var text string  // ""
    text = text + "Hello"
    fmt.Println("text after concatenation:", text)
    
    var flag bool  // false
    flag = !flag
    fmt.Println("flag after negation:", flag)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Zero Values ===
int zero value: 0
float64 zero value: 0
string zero value: []
bool zero value: false

=== Using Zero Values ===
counter after +=5: 5
text after concatenation: Hello
flag after negation: true
```

**Learning Objectives:**
- ✅ Understand zero values for all types
- ✅ Know that variables are always initialized
- ✅ See how zero values work in operations

---

## Exercise 4: Type Conversion

**Objective:** Practice explicit type conversion

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level3-exercise4
cd ~/projects/level3-exercise4
go mod init level3.example/exercise4
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "strconv"
)

func main() {
    fmt.Println("=== Numeric Type Conversion ===")
    
    // int to float64
    var i int = 42
    var f float64 = float64(i)
    fmt.Printf("int %d to float64 %.1f\n", i, f)
    
    // float64 to int (loses decimal)
    var pi float64 = 3.14159
    var rounded int = int(pi)
    fmt.Printf("float64 %.5f to int %d\n", pi, rounded)
    
    // int to uint
    var num int = 100
    var unsigned uint = uint(num)
    fmt.Printf("int %d to uint %d\n", num, unsigned)
    
    fmt.Println("\n=== String Conversions ===")
    
    // String to int
    str := "123"
    num, err := strconv.Atoi(str)
    if err == nil {
        fmt.Printf("string \"%s\" to int %d\n", str, num)
    }
    
    // Int to string
    intVal := 456
    strVal := strconv.Itoa(intVal)
    fmt.Printf("int %d to string \"%s\"\n", intVal, strVal)
    
    // String to float
    str2 := "3.14"
    floatVal, err := strconv.ParseFloat(str2, 64)
    if err == nil {
        fmt.Printf("string \"%s\" to float64 %.2f\n", str2, floatVal)
    }
    
    fmt.Println("\n=== Byte/String Conversions ===")
    
    // String to bytes
    text := "Hello"
    bytes := []byte(text)
    fmt.Printf("string \"%s\" to bytes %v\n", text, bytes)
    
    // Bytes to string
    bytesArray := []byte{72, 101, 108, 108, 111}
    backToString := string(bytesArray)
    fmt.Printf("bytes %v to string \"%s\"\n", bytesArray, backToString)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Numeric Type Conversion ===
int 42 to float64 42.0
float64 3.14159 to int 3
int 100 to uint 100

=== String Conversions ===
string "123" to int 123
int 456 to string "456"
string "3.14" to float64 3.14

=== Byte/String Conversions ===
string "Hello" to bytes [72 101 108 108 111]
bytes [72 101 108 108 111] to string "Hello"
```

**Learning Objectives:**
- ✅ Convert between numeric types
- ✅ Convert strings to numbers
- ✅ Convert numbers to strings
- ✅ Handle conversion errors
- ✅ Convert between strings and bytes

---

## Exercise 5: Constants

**Objective:** Learn to declare and use constants

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level3-exercise5
cd ~/projects/level3-exercise5
go mod init level3.example/exercise5
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

// Package-level constants
const AppName = "MyApp"
const Version = "1.0.0"
const MaxUsers int = 1000
const Pi float64 = 3.14159

// Using iota for sequences
const (
    Sunday = iota
    Monday
    Tuesday
    Wednesday
    Thursday
    Friday
    Saturday
)

// Bit flags with iota
const (
    Read   = 1 << iota  // 1 (2^0)
    Write              // 2 (2^1)
    Execute            // 4 (2^2)
)

func main() {
    fmt.Println("=== Constants ===")
    fmt.Println("App:", AppName)
    fmt.Println("Version:", Version)
    fmt.Println("Max Users:", MaxUsers)
    fmt.Println("Pi:", Pi)
    
    fmt.Println("\n=== Days of Week (iota) ===")
    fmt.Println("Sunday:", Sunday)
    fmt.Println("Monday:", Monday)
    fmt.Println("Wednesday:", Wednesday)
    fmt.Println("Saturday:", Saturday)
    
    fmt.Println("\n=== Permission Flags ===")
    fmt.Println("Read:", Read)
    fmt.Println("Write:", Write)
    fmt.Println("Execute:", Execute)
    
    // Combine flags
    permissions := Read | Write
    fmt.Println("Read + Write:", permissions)
    
    // Check if flag is set
    if (permissions & Read) == Read {
        fmt.Println("Has read permission")
    }
    if (permissions & Execute) == Execute {
        fmt.Println("Has execute permission")
    } else {
        fmt.Println("No execute permission")
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
=== Constants ===
App: MyApp
Version: 1.0.0
Max Users: 1000
Pi: 3.14159

=== Days of Week (iota) ===
Sunday: 0
Monday: 1
Wednesday: 3
Saturday: 6

=== Permission Flags ===
Read: 1
Write: 2
Execute: 4
Read + Write: 3
Has read permission
No execute permission
```

**Learning Objectives:**
- ✅ Declare constants with const
- ✅ Understand constant immutability
- ✅ Use iota for sequences
- ✅ Create bit flags
- ✅ Combine flags with bitwise OR

---

## Exercise 6: Named Types

**Objective:** Create and use custom type names

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level3-exercise6
cd ~/projects/level3-exercise6
go mod init level3.example/exercise6
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

// Define named types
type Age int
type Temperature float64
type UserId string
type Score float64

// Add methods to named types
func (a Age) IsAdult() bool {
    return a >= 18
}

func (t Temperature) Fahrenheit() string {
    return fmt.Sprintf("%.1f°F", t*9/5+32)
}

func (s Score) Grade() string {
    if s >= 90 {
        return "A"
    } else if s >= 80 {
        return "B"
    } else if s >= 70 {
        return "C"
    }
    return "F"
}

func main() {
    fmt.Println("=== Named Types ===")
    
    // Age
    age1 := Age(15)
    age2 := Age(25)
    fmt.Printf("Age %d is adult: %v\n", age1, age1.IsAdult())
    fmt.Printf("Age %d is adult: %v\n", age2, age2.IsAdult())
    
    // Temperature
    temp := Temperature(25)
    fmt.Printf("Temperature: %.1f°C = %s\n", temp, temp.Fahrenheit())
    
    // UserId
    user1 := UserId("alice123")
    user2 := UserId("bob456")
    fmt.Printf("Users: %s, %s\n", user1, user2)
    
    // Score with grade
    scores := []Score{95, 85, 72, 60}
    fmt.Println("\nScores and Grades:")
    for _, score := range scores {
        fmt.Printf("Score: %.0f → Grade: %s\n", score, score.Grade())
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
=== Named Types ===
Age 15 is adult: false
Age 25 is adult: true
Temperature: 25.0°C = 77.0°F

Users: alice123, bob456

Scores and Grades:
Score: 95 → Grade: A
Score: 85 → Grade: B
Score: 72 → Grade: C
Score: 60 → Grade: F
```

**Learning Objectives:**
- ✅ Define named types
- ✅ Add methods to named types
- ✅ Use named types for type safety
- ✅ Improve code clarity with named types

---

## Exercise 7: Working with Multiple Variables

**Objective:** Declare and manage multiple variables efficiently

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level3-exercise7
cd ~/projects/level3-exercise7
go mod init level3.example/exercise7
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

// Multiple package-level variables
var (
    AppVersion string = "2.0.1"
    MaxRetries int = 5
    Timeout float64 = 30.0
    Debug bool = false
)

func main() {
    fmt.Println("=== Multiple Variable Declaration ===")
    
    // Method 1: Multiple declarations
    var x, y, z int = 1, 2, 3
    fmt.Printf("x=%d, y=%d, z=%d\n", x, y, z)
    
    // Method 2: Multiple declarations same type
    var name, city, country string = "Alice", "NYC", "USA"
    fmt.Printf("Name: %s, City: %s, Country: %s\n", name, city, country)
    
    // Method 3: Short declaration multiple
    a, b := 10, 20
    fmt.Printf("a=%d, b=%d\n", a, b)
    
    // Method 4: Swap values
    fmt.Println("\nSwap values:")
    a, b = b, a
    fmt.Printf("After swap: a=%d, b=%d\n", a, b)
    
    // Method 5: Multiple return values
    fmt.Println("\nMultiple assignment from function:")
    result, remainder := divmod(17, 5)
    fmt.Printf("17 ÷ 5 = %d remainder %d\n", result, remainder)
    
    fmt.Println("\nPackage-level variables:")
    fmt.Printf("App v%s, MaxRetries=%d, Timeout=%.1fs, Debug=%v\n",
        AppVersion, MaxRetries, Timeout, Debug)
}

func divmod(a, b int) (int, int) {
    return a / b, a % b
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Multiple Variable Declaration ===
x=1, y=2, z=3
Name: Alice, City: NYC, Country: USA
a=10, b=20

Swap values:
After swap: a=20, b=10

Multiple assignment from function:
17 ÷ 5 = 3 remainder 2

Package-level variables:
App v2.0.1, MaxRetries=5, Timeout=30.0s, Debug=false
```

**Learning Objectives:**
- ✅ Declare multiple variables
- ✅ Swap variable values
- ✅ Handle multiple return values
- ✅ Use package-level variable blocks

---

## Exercise 8: Type Safety with Named Types

**Objective:** Understand how named types prevent mixing different types

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level3-exercise8
cd ~/projects/level3-exercise8
go mod init level3.example/exercise8
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

// Different named types based on int
type UserId int
type ProductId int
type Quantity int

func main() {
    fmt.Println("=== Type Safety Example ===")
    
    user := UserId(42)
    product := ProductId(100)
    qty := Quantity(5)
    
    fmt.Printf("User ID: %v\n", user)
    fmt.Printf("Product ID: %v\n", product)
    fmt.Printf("Quantity: %v\n", qty)
    
    // This works - same type
    var a UserId = 1
    var b UserId = 2
    sum := int(a) + int(b)
    fmt.Printf("\nUserId(1) + UserId(2) = %d\n", sum)
    
    // This requires explicit conversion - catches mistakes!
    fmt.Printf("\nProduct %d, Quantity %d\n", product, qty)
    
    // Mixing types requires explicit conversion (safer!)
    // productWithQty := product + qty  // ERROR! Different types
    productWithQty := ProductId(int(product) + int(qty))  // OK with conversion
    fmt.Printf("Product ID %d + Quantity %d = %d\n", product, qty, productWithQty)
    
    fmt.Println("\nNamed types prevent accidental mixing!")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Type Safety Example ===
User ID: 42
Product ID: 100
Quantity: 5

UserId(1) + UserId(2) = 3

Product 100, Quantity 5
Product ID 100 + Quantity 5 = 105

Named types prevent accidental mixing!
```

**Learning Objectives:**
- ✅ Understand type safety
- ✅ Create distinct types that can't be accidentally mixed
- ✅ Use named types to prevent bugs

---

## Exercise 9: Intro to Composite Types

**Objective:** Get first experience with arrays, slices, and maps

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level3-exercise9
cd ~/projects/level3-exercise9
go mod init level3.example/exercise9
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== Arrays ===")
    var colors [3]string
    colors[0] = "Red"
    colors[1] = "Green"
    colors[2] = "Blue"
    
    for i, color := range colors {
        fmt.Printf("%d: %s\n", i, color)
    }
    
    fmt.Println("\n=== Slices ===")
    numbers := []int{1, 2, 3, 4, 5}
    fmt.Printf("Slice: %v\n", numbers)
    fmt.Printf("Length: %d, Capacity: %d\n", len(numbers), cap(numbers))
    
    // Slice operations
    fmt.Printf("First 3: %v\n", numbers[:3])
    fmt.Printf("From index 2: %v\n", numbers[2:])
    
    // Append
    numbers = append(numbers, 6, 7)
    fmt.Printf("After append: %v\n", numbers)
    
    fmt.Println("\n=== Maps ===")
    scores := map[string]int{
        "Alice": 95,
        "Bob": 87,
        "Charlie": 92,
    }
    
    fmt.Println("Scores:")
    for name, score := range scores {
        fmt.Printf("%s: %d\n", name, score)
    }
    
    // Add
    scores["Diana"] = 89
    
    // Delete
    delete(scores, "Bob")
    
    fmt.Println("\nAfter adding Diana and removing Bob:")
    for name, score := range scores {
        fmt.Printf("%s: %d\n", name, score)
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
=== Arrays ===
0: Red
1: Green
2: Blue

=== Slices ===
Slice: [1 2 3 4 5]
Length: 5, Capacity: 5
First 3: [1 2 3]
From index 2: [3 4 5]
After append: [1 2 3 4 5 6 7]

=== Maps ===
Scores:
Alice: 95
Bob: 87
Charlie: 92

After adding Diana and removing Bob:
Alice: 95
Charlie: 92
Diana: 89
```

**Learning Objectives:**
- ✅ Use arrays
- ✅ Use slices
- ✅ Use maps
- ✅ Understand basic operations on composite types

---

## Exercise 10: Comprehensive Type Practice

**Objective:** Apply all Level 3 concepts in a realistic program

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level3-exercise10
cd ~/projects/level3-exercise10
go mod init level3.example/exercise10
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "strconv"
)

// Named types
type StudentName string
type Grade float64

// Constants
const (
    PassingGrade = 60.0
    PerfectScore = 100.0
)

func main() {
    fmt.Println("=== Student Grade Analyzer ===\n")
    
    // Student data
    students := map[StudentName]Grade{
        "Alice": 95,
        "Bob": 87,
        "Charlie": 92,
        "Diana": 58,
        "Eve": 78,
    }
    
    var passCount, failCount int = 0, 0
    var totalGrade float64 = 0
    
    fmt.Println("Student Grades:")
    fmt.Println(string(make([]byte, 40)))
    
    for name, grade := range students {
        status := "PASS"
        if grade < PassingGrade {
            status = "FAIL"
            failCount++
        } else {
            passCount++
        }
        
        totalGrade += float64(grade)
        fmt.Printf("%-15s %6.1f  %s\n", name, grade, status)
    }
    
    // Calculate average
    avgGrade := totalGrade / float64(len(students))
    avgStr := strconv.FormatFloat(avgGrade, 'f', 1, 64)
    
    // Display statistics
    fmt.Println(string(make([]byte, 40)))
    fmt.Printf("Total Students: %d\n", len(students))
    fmt.Printf("Passed: %d\n", passCount)
    fmt.Printf("Failed: %d\n", failCount)
    fmt.Printf("Average: %s\n", avgStr)
    
    // Performance indicator
    if avgGrade >= 90 {
        fmt.Println("Class Performance: Excellent! 🌟")
    } else if avgGrade >= 80 {
        fmt.Println("Class Performance: Good! ✓")
    } else {
        fmt.Println("Class Performance: Needs Improvement")
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
=== Student Grade Analyzer ===

Student Grades:
                                        
Alice               95.0  PASS
Bob                 87.0  PASS
Charlie             92.0  PASS
Diana               58.0  FAIL
Eve                 78.0  PASS
                                        
Total Students: 5
Passed: 4
Failed: 1
Average: 82.0
Class Performance: Good! ✓
```

**Learning Objectives:**
- ✅ Combine all Level 3 concepts
- ✅ Use named types in realistic program
- ✅ Work with constants
- ✅ Use type conversion
- ✅ Manage multiple variables and types

---

## Bonus Challenges

### Challenge 1: Temperature Converter

Create a program that converts between Celsius, Fahrenheit, and Kelvin.

```bash
mkdir -p ~/projects/level3-bonus1
cd ~/projects/level3-bonus1
go mod init level3.example/bonus1
```

**Hints:**
- Create named types: `type Celsius float64`, `type Fahrenheit float64`
- Add methods for conversion
- Handle edge cases (absolute zero)

### Challenge 2: Unit Converter

Create a program that converts between different units (meters to feet, kilograms to pounds, etc.).

```bash
mkdir -p ~/projects/level3-bonus2
cd ~/projects/level3-bonus2
go mod init level3.example/bonus2
```

**Hints:**
- Use constants for conversion factors
- Create named types for different units
- Add methods that return converted values

### Challenge 3: Grade Distribution

Create a program that analyzes grade distribution and creates a histogram.

```bash
mkdir -p ~/projects/level3-bonus3
cd ~/projects/level3-bonus3
go mod init level3.example/bonus3
```

**Hints:**
- Use a map to store grade ranges
- Use slices to store grades
- Use loops to build a text histogram

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Declare variables in multiple ways
✅ Understand Go's basic types
✅ Convert between types
✅ Work with constants
✅ Create named types
✅ Handle multiple variables
✅ Use composite types (arrays, slices, maps)
✅ Apply all concepts in realistic programs

---

## Next Level

Level 4: Operators & Expressions
- Arithmetic operators
- Comparison operators
- Logical operators
- Bitwise operators
- Operator precedence

Great work! You're mastering Go! 🚀
