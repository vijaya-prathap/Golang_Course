# Level 4: Operators & Expressions - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

---

## Exercise 1: Arithmetic Operators

**Objective:** Practice all basic arithmetic operators

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level4-exercise1
cd ~/projects/level4-exercise1
go mod init level4.example/exercise1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    a := 17
    b := 5

    fmt.Println("=== Arithmetic Operators ===")
    fmt.Printf("a + b = %d\n", a+b)
    fmt.Printf("a - b = %d\n", a-b)
    fmt.Printf("a * b = %d\n", a*b)
    fmt.Printf("a / b = %d\n", a/b)
    fmt.Printf("a %% b = %d\n", a%b)

    fmt.Println("\n=== Unary Operators ===")
    x := 8
    fmt.Printf("x = %d\n", x)
    fmt.Printf("-x = %d\n", -x)
    fmt.Printf("+x = %d\n", +x)

    fmt.Println("\n=== Float Arithmetic ===")
    f1 := 10.5
    f2 := 3.2
    fmt.Printf("f1 + f2 = %.2f\n", f1+f2)
    fmt.Printf("f1 / f2 = %.4f\n", f1/f2)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Arithmetic Operators ===
a + b = 22
a - b = 12
a * b = 85
a / b = 3
a % b = 2

=== Unary Operators ===
x = 8
-x = -8
+x = 8

=== Float Arithmetic ===
f1 + f2 = 13.70
f1 / f2 = 3.2812
```

**Learning Objectives:**
- ✅ Use all arithmetic operators (+, -, *, /, %)
- ✅ Understand unary +/-
- ✅ See the difference between int and float arithmetic

---

## Exercise 2: Integer Division & Modulo

**Objective:** Understand integer division truncation and modulo pitfalls

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level4-exercise2
cd ~/projects/level4-exercise2
go mod init level4.example/exercise2
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func isEven(n int) bool {
    return n%2 == 0
}

func main() {
    fmt.Println("=== Integer Division Truncates ===")
    fmt.Printf("10 / 3 = %d\n", 10/3)
    fmt.Printf("10 / 4 = %d\n", 10/4)
    fmt.Printf("7 / 2 = %d\n", 7/2)

    fmt.Println("\n=== Fixing It With float64 ===")
    fmt.Printf("float64(10) / 3 = %.4f\n", float64(10)/3)
    fmt.Printf("float64(7) / 2 = %.4f\n", float64(7)/2)

    fmt.Println("\n=== Modulo For Even/Odd ===")
    for _, n := range []int{4, 7, 10, 15} {
        if isEven(n) {
            fmt.Printf("%d is even\n", n)
        } else {
            fmt.Printf("%d is odd\n", n)
        }
    }

    fmt.Println("\n=== Modulo For Wrapping ===")
    for hour := 22; hour <= 26; hour++ {
        fmt.Printf("hour %d -> %d on a 24-hour clock\n", hour, hour%24)
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
=== Integer Division Truncates ===
10 / 3 = 3
10 / 4 = 2
7 / 2 = 3

=== Fixing It With float64 ===
float64(10) / 3 = 3.3333
float64(7) / 2 = 3.5000

=== Modulo For Even/Odd ===
4 is even
7 is odd
10 is even
15 is odd

=== Modulo For Wrapping ===
hour 22 -> 22 on a 24-hour clock
hour 23 -> 23 on a 24-hour clock
hour 24 -> 0 on a 24-hour clock
hour 25 -> 1 on a 24-hour clock
hour 26 -> 2 on a 24-hour clock
```

**Learning Objectives:**
- ✅ Recognize integer division truncation
- ✅ Convert to float64 to get precise results
- ✅ Use modulo for even/odd checks
- ✅ Use modulo for wrapping/cyclic values

---

## Exercise 3: Comparison Operators

**Objective:** Compare numbers and strings

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level4-exercise3
cd ~/projects/level4-exercise3
go mod init level4.example/exercise3
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    a := 10
    b := 3

    fmt.Println("=== Numeric Comparisons ===")
    fmt.Printf("a == b: %v\n", a == b)
    fmt.Printf("a != b: %v\n", a != b)
    fmt.Printf("a > b: %v\n", a > b)
    fmt.Printf("a >= b: %v\n", a >= b)
    fmt.Printf("a < b: %v\n", a < b)
    fmt.Printf("a <= b: %v\n", a <= b)

    fmt.Println("\n=== Float Comparisons ===")
    fmt.Printf("5.5 > 3.2: %v\n", 5.5 > 3.2)

    fmt.Println("\n=== String Comparisons (Lexicographic) ===")
    fmt.Printf("\"apple\" == \"apple\": %v\n", "apple" == "apple")
    fmt.Printf("\"abc\" < \"xyz\": %v\n", "abc" < "xyz")
    fmt.Printf("\"apple\" > \"apricot\": %v\n", "apple" > "apricot")
    fmt.Printf("\"Apple\" == \"apple\": %v\n", "Apple" == "apple")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Numeric Comparisons ===
a == b: false
a != b: true
a > b: true
a >= b: true
a < b: false
a <= b: false

=== Float Comparisons ===
5.5 > 3.2: true

=== String Comparisons (Lexicographic) ===
"apple" == "apple": true
"abc" < "xyz": true
"apple" > "apricot": false
"Apple" == "apple": false
```

**Learning Objectives:**
- ✅ Use all comparison operators
- ✅ Understand comparisons return bool
- ✅ Understand lexicographic string comparison
- ✅ Understand string comparison is case-sensitive

---

## Exercise 4: Logical Operators

**Objective:** Combine boolean values with &&, ||, and !

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level4-exercise4
cd ~/projects/level4-exercise4
go mod init level4.example/exercise4
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func canDrive(age int, hasLicense bool) bool {
    return age >= 18 && hasLicense
}

func main() {
    fmt.Println("=== Basic Logical Operators ===")
    a := true
    b := false
    fmt.Printf("a && b: %v\n", a && b)
    fmt.Printf("a || b: %v\n", a || b)
    fmt.Printf("!a: %v\n", !a)
    fmt.Printf("!!a: %v\n", !!a)

    fmt.Println("\n=== Validation Example ===")
    fmt.Printf("age=20, license=true  -> canDrive: %v\n", canDrive(20, true))
    fmt.Printf("age=20, license=false -> canDrive: %v\n", canDrive(20, false))
    fmt.Printf("age=15, license=true  -> canDrive: %v\n", canDrive(15, true))

    fmt.Println("\n=== Combining Conditions ===")
    x, y, z := 5, 10, -2
    fmt.Printf("(x > 0) && (y > 0): %v\n", (x > 0) && (y > 0))
    fmt.Printf("(x < 0) || (z < 0): %v\n", (x < 0) || (z < 0))
    fmt.Printf("!(x == y): %v\n", !(x == y))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Basic Logical Operators ===
a && b: false
a || b: true
!a: false
!!a: true

=== Validation Example ===
age=20, license=true  -> canDrive: true
age=20, license=false -> canDrive: false
age=15, license=true  -> canDrive: false

=== Combining Conditions ===
(x > 0) && (y > 0): true
(x < 0) || (z < 0): true
!(x == y): true
```

**Learning Objectives:**
- ✅ Use &&, ||, and ! correctly
- ✅ Combine multiple conditions
- ✅ Apply logical operators to real validation logic

---

## Exercise 5: Bitwise Operators

**Objective:** Manipulate individual bits and build a permission-flag system

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level4-exercise5
cd ~/projects/level4-exercise5
go mod init level4.example/exercise5
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

// Bit flags with iota (pattern from Level 3)
const (
    Read   = 1 << iota // 1 (001)
    Write              // 2 (010)
    Execute            // 4 (100)
)

func main() {
    fmt.Println("=== Basic Bitwise Operators ===")
    fmt.Printf("5 & 3 = %d\n", 5&3)
    fmt.Printf("5 | 3 = %d\n", 5|3)
    fmt.Printf("5 ^ 3 = %d\n", 5^3)
    fmt.Printf("^5 = %d\n", ^5)

    fmt.Println("\n=== Shifts ===")
    fmt.Printf("5 << 1 = %d\n", 5<<1)
    fmt.Printf("5 << 2 = %d\n", 5<<2)
    fmt.Printf("5 >> 1 = %d\n", 5>>1)
    fmt.Printf("5 >> 2 = %d\n", 5>>2)

    fmt.Println("\n=== Permission Flags ===")
    perms := Read | Write
    fmt.Printf("Read: %d, Write: %d, Execute: %d\n", Read, Write, Execute)
    fmt.Printf("perms (Read|Write) = %d\n", perms)

    if perms&Read == Read {
        fmt.Println("Has read permission")
    }
    if perms&Execute == Execute {
        fmt.Println("Has execute permission")
    } else {
        fmt.Println("No execute permission")
    }

    perms |= Execute
    fmt.Printf("\nAfter granting execute, perms = %d\n", perms)
    if perms&Execute == Execute {
        fmt.Println("Has execute permission")
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
=== Basic Bitwise Operators ===
5 & 3 = 1
5 | 3 = 7
5 ^ 3 = 6
^5 = -6

=== Shifts ===
5 << 1 = 10
5 << 2 = 20
5 >> 1 = 2
5 >> 2 = 1

=== Permission Flags ===
Read: 1, Write: 2, Execute: 4
perms (Read|Write) = 3
Has read permission
No execute permission

After granting execute, perms = 7
Has execute permission
```

**Learning Objectives:**
- ✅ Use &, |, ^, and ^ (unary NOT)
- ✅ Use << and >> for shifting
- ✅ Build and check a bit-flag permission system

---

## Exercise 6: Assignment Operators

**Objective:** Practice compound assignment and increment/decrement

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level4-exercise6
cd ~/projects/level4-exercise6
go mod init level4.example/exercise6
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    x := 5
    fmt.Println("=== Compound Arithmetic Assignment ===")
    fmt.Printf("start: x = %d\n", x)

    x += 3
    fmt.Printf("x += 3 -> %d\n", x)

    x -= 2
    fmt.Printf("x -= 2 -> %d\n", x)

    x *= 2
    fmt.Printf("x *= 2 -> %d\n", x)

    x /= 3
    fmt.Printf("x /= 3 -> %d\n", x)

    x %= 3
    fmt.Printf("x %%= 3 -> %d\n", x)

    fmt.Println("\n=== Compound Bitwise Assignment ===")
    flags := 5
    fmt.Printf("start: flags = %d\n", flags)

    flags &= 6
    fmt.Printf("flags &= 6 -> %d\n", flags)

    flags |= 2
    fmt.Printf("flags |= 2 -> %d\n", flags)

    flags ^= 1
    fmt.Printf("flags ^= 1 -> %d\n", flags)

    flags <<= 1
    fmt.Printf("flags <<= 1 -> %d\n", flags)

    flags >>= 2
    fmt.Printf("flags >>= 2 -> %d\n", flags)

    fmt.Println("\n=== Increment / Decrement ===")
    counter := 5
    counter++
    fmt.Printf("counter++ -> %d\n", counter)
    counter--
    fmt.Printf("counter-- -> %d\n", counter)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Compound Arithmetic Assignment ===
start: x = 5
x += 3 -> 8
x -= 2 -> 6
x *= 2 -> 12
x /= 3 -> 4
x %= 3 -> 1

=== Compound Bitwise Assignment ===
start: flags = 5
flags &= 6 -> 4
flags |= 2 -> 6
flags ^= 1 -> 7
flags <<= 1 -> 14
flags >>= 2 -> 3

=== Increment / Decrement ===
counter++ -> 6
counter-- -> 5
```

**Learning Objectives:**
- ✅ Use compound assignment operators for arithmetic
- ✅ Use compound assignment operators for bitwise ops
- ✅ Use ++ and -- as statements (not expressions)

---

## Exercise 7: Operator Precedence

**Objective:** Predict, then verify, how precedence changes results

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level4-exercise7
cd ~/projects/level4-exercise7
go mod init level4.example/exercise7
```

2. Before running the code, write down what you think each line prints.

3. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== Arithmetic Precedence ===")
    fmt.Printf("2 + 3 * 4 = %d\n", 2+3*4)
    fmt.Printf("(2 + 3) * 4 = %d\n", (2+3)*4)
    fmt.Printf("2 + 3 * 4 - 5 = %d\n", 2+3*4-5)
    fmt.Printf("10 - 2 * 3 + 1 = %d\n", 10-2*3+1)

    fmt.Println("\n=== Comparison + Logical Precedence ===")
    fmt.Printf("5 > 3 && 10 < 20 = %v\n", 5 > 3 && 10 < 20)
    fmt.Printf("5 > 3 || 10 > 20 = %v\n", 5 > 3 || 10 > 20)
    fmt.Printf("true || false && false = %v\n", true || false && false)
    fmt.Printf("(true || false) && false = %v\n", (true || false) && false)

    fmt.Println("\n=== Mixed Precedence ===")
    x, y, z := 6, 11, 2
    fmt.Printf("x > 5 && y > 10 || z < 3 = %v\n", x > 5 && y > 10 || z < 3)
    fmt.Printf("(x > 5 && y > 10) || (z < 3) = %v\n", (x > 5 && y > 10) || (z < 3))
}
EOF
```

4. Run the program and compare with your predictions:

```bash
go run main.go
```

**Expected Output:**

```
=== Arithmetic Precedence ===
2 + 3 * 4 = 14
(2 + 3) * 4 = 20
2 + 3 * 4 - 5 = 9
10 - 2 * 3 + 1 = 5

=== Comparison + Logical Precedence ===
5 > 3 && 10 < 20 = true
5 > 3 || 10 > 20 = true
true || false && false = true
(true || false) && false = false

=== Mixed Precedence ===
x > 5 && y > 10 || z < 3 = true
(x > 5 && y > 10) || (z < 3) = true
```

**Learning Objectives:**
- ✅ Predict evaluation order before running code
- ✅ See how parentheses override default precedence
- ✅ Understand && binds tighter than ||

---

## Exercise 8: Short-Circuit Evaluation

**Objective:** Use short-circuiting to write safe guard conditions

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level4-exercise8
cd ~/projects/level4-exercise8
go mod init level4.example/exercise8
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func safeDivide(a, b int) {
    // Safe: b != 0 is checked FIRST, so a/b never runs when b is 0
    if b != 0 && a/b > 1 {
        fmt.Printf("%d / %d = %d, and result > 1\n", a, b, a/b)
    } else {
        fmt.Printf("skipped division for a=%d, b=%d (b is zero or result <= 1)\n", a, b)
    }
}

func loggedTrue(label string, value bool) bool {
    fmt.Printf("  evaluated: %s\n", label)
    return value
}

func main() {
    fmt.Println("=== Safe Guard With && ===")
    safeDivide(10, 2)
    safeDivide(10, 0)
    safeDivide(1, 5)

    fmt.Println("\n=== Proving Short-Circuit With && ===")
    fmt.Println("false && loggedTrue(...):")
    result := false && loggedTrue("right side", true)
    fmt.Printf("result = %v (right side was NOT evaluated)\n", result)

    fmt.Println("\ntrue && loggedTrue(...):")
    result = true && loggedTrue("right side", true)
    fmt.Printf("result = %v (right side WAS evaluated)\n", result)

    fmt.Println("\n=== Proving Short-Circuit With || ===")
    fmt.Println("true || loggedTrue(...):")
    result = true || loggedTrue("right side", false)
    fmt.Printf("result = %v (right side was NOT evaluated)\n", result)

    fmt.Println("\nfalse || loggedTrue(...):")
    result = false || loggedTrue("right side", false)
    fmt.Printf("result = %v (right side WAS evaluated)\n", result)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Safe Guard With && ===
10 / 2 = 5, and result > 1
skipped division for a=10, b=0 (b is zero or result <= 1)
skipped division for a=1, b=5 (b is zero or result <= 1)

=== Proving Short-Circuit With && ===
false && loggedTrue(...):
result = false (right side was NOT evaluated)

true && loggedTrue(...):
  evaluated: right side
result = true (right side WAS evaluated)

=== Proving Short-Circuit With || ===
true || loggedTrue(...):
result = true (right side was NOT evaluated)

false || loggedTrue(...):
  evaluated: right side
result = false (right side WAS evaluated)
```

**Learning Objectives:**
- ✅ Write a safe guard condition using && ordering
- ✅ Prove short-circuit behavior for && and ||
- ✅ Understand why operand order matters

---

## Exercise 9: Building Expressions (Mini Calculator)

**Objective:** Combine arithmetic and comparison operators into a small calculator

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level4-exercise9
cd ~/projects/level4-exercise9
go mod init level4.example/exercise9
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func calculate(a, b float64, op string) (float64, bool) {
    switch op {
    case "+":
        return a + b, true
    case "-":
        return a - b, true
    case "*":
        return a * b, true
    case "/":
        if b == 0 {
            return 0, false
        }
        return a / b, true
    default:
        return 0, false
    }
}

func main() {
    fmt.Println("=== Mini Calculator ===")

    type calc struct {
        a, b float64
        op   string
    }

    calculations := []calc{
        {10, 4, "+"},
        {10, 4, "-"},
        {10, 4, "*"},
        {10, 4, "/"},
        {10, 0, "/"},
    }

    for _, c := range calculations {
        result, ok := calculate(c.a, c.b, c.op)
        if ok {
            fmt.Printf("%.1f %s %.1f = %.4f\n", c.a, c.op, c.b, result)
        } else {
            fmt.Printf("%.1f %s %.1f = error (invalid operation)\n", c.a, c.op, c.b)
        }
    }

    fmt.Println("\n=== Using Expressions For Decisions ===")
    result, _ := calculate(15, 4, "/")
    isWholeNumber := result == float64(int(result))
    fmt.Printf("15 / 4 = %.4f, isWholeNumber = %v\n", result, isWholeNumber)

    result, _ = calculate(20, 4, "/")
    isWholeNumber = result == float64(int(result))
    fmt.Printf("20 / 4 = %.4f, isWholeNumber = %v\n", result, isWholeNumber)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Mini Calculator ===
10.0 + 4.0 = 14.0000
10.0 - 4.0 = 6.0000
10.0 * 4.0 = 40.0000
10.0 / 4.0 = 2.5000
10.0 / 0.0 = error (invalid operation)

=== Using Expressions For Decisions ===
15 / 4 = 3.7500, isWholeNumber = false
20 / 4 = 5.0000, isWholeNumber = true
```

**Learning Objectives:**
- ✅ Combine arithmetic and comparison into an expression
- ✅ Use expression results to make decisions
- ✅ Guard against division by zero with a boolean condition

---

## Exercise 10: Comprehensive Practice — BMI Calculator

**Objective:** Apply arithmetic, comparison, logical, and assignment operators in one realistic program

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level4-exercise10
cd ~/projects/level4-exercise10
go mod init level4.example/exercise10
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type person struct {
    name      string
    weightKg  float64
    heightM   float64
}

func bmi(weightKg, heightM float64) float64 {
    return weightKg / (heightM * heightM)
}

func category(bmiValue float64) string {
    switch {
    case bmiValue < 18.5:
        return "Underweight"
    case bmiValue < 25:
        return "Normal"
    case bmiValue < 30:
        return "Overweight"
    default:
        return "Obese"
    }
}

func main() {
    people := []person{
        {"Alice", 60, 1.65},
        {"Bob", 90, 1.80},
        {"Charlie", 50, 1.70},
        {"Diana", 75, 1.60},
    }

    fmt.Println("=== BMI Calculator ===\n")

    var normalCount, outsideNormalCount int
    var totalBMI float64

    for _, p := range people {
        b := bmi(p.weightKg, p.heightM)
        cat := category(b)
        totalBMI += b

        isNormal := b >= 18.5 && b < 25
        if isNormal {
            normalCount++
        } else {
            outsideNormalCount++
        }

        fmt.Printf("%-10s weight=%5.1fkg height=%.2fm BMI=%5.2f -> %s\n",
            p.name, p.weightKg, p.heightM, b, cat)
    }

    avgBMI := totalBMI / float64(len(people))

    fmt.Println()
    fmt.Printf("Average BMI: %.2f\n", avgBMI)
    fmt.Printf("Normal range: %d, Outside normal range: %d\n", normalCount, outsideNormalCount)

    if normalCount >= outsideNormalCount && normalCount > 0 {
        fmt.Println("Group summary: Majority within normal range")
    } else {
        fmt.Println("Group summary: Majority outside normal range")
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
=== BMI Calculator ===

Alice      weight= 60.0kg height=1.65m BMI=22.04 -> Normal
Bob        weight= 90.0kg height=1.80m BMI=27.78 -> Overweight
Charlie    weight= 50.0kg height=1.70m BMI=17.30 -> Underweight
Diana      weight= 75.0kg height=1.60m BMI=29.30 -> Overweight

Average BMI: 24.10
Normal range: 1, Outside normal range: 3
Group summary: Majority outside normal range
```

**Learning Objectives:**
- ✅ Combine all Level 4 operator categories in one program
- ✅ Use arithmetic expressions inside function calls
- ✅ Use comparison and logical operators to classify data
- ✅ Use compound assignment (+=) to accumulate totals

---

## Bonus Challenges

### Challenge 1: Permission System

Build a more complete file-permission system using bitwise operators (Read, Write, Execute) that supports granting, revoking, and checking multiple flags at once.

```bash
mkdir -p ~/projects/level4-bonus1
cd ~/projects/level4-bonus1
go mod init level4.example/bonus1
```

**Hints:**
- Reuse the `Read`/`Write`/`Execute` iota pattern from Exercise 5
- Add a `Revoke(perms, flag int) int` helper using `perms &^ flag` or `perms & ^flag`
- Print permissions as a human-readable string like `"rwx"` or `"r--"`

### Challenge 2: Precedence Puzzle Quiz

Write a program that prints a series of un-parenthesized expressions as strings, asks the user (in comments) to predict the result, then evaluates and reveals the real answer.

```bash
mkdir -p ~/projects/level4-bonus2
cd ~/projects/level4-bonus2
go mod init level4.example/bonus2
```

**Hints:**
- Include at least one expression mixing arithmetic and comparison
- Include at least one expression mixing && and ||
- Show both the "naive" (wrong) expected order and the actual Go evaluation order

### Challenge 3: Running Statistics

Create a program that processes a slice of numbers and tracks running min, max, sum, and average using compound assignment operators inside a loop.

```bash
mkdir -p ~/projects/level4-bonus3
cd ~/projects/level4-bonus3
go mod init level4.example/bonus3
```

**Hints:**
- Initialize `min` and `max` to the first element, not 0
- Use `+=` to accumulate the sum
- Use comparison operators inside the loop to update `min`/`max`

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Use all arithmetic operators, including integer division and modulo
✅ Use all comparison operators on numbers and strings
✅ Combine conditions with logical operators
✅ Manipulate bits and build flag systems with bitwise operators
✅ Use compound assignment and increment/decrement operators
✅ Predict how operator precedence affects results
✅ Rely on (and prove) short-circuit evaluation
✅ Build and evaluate real expressions in realistic programs

---

## Next Level

Level 5: Conditions (if/else/switch)
- Control flow with if statements
- Switch statements
- Making decisions in code

Great work! You're mastering Go! 🚀
