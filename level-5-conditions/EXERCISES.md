# Level 5: Conditions - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

---

## Exercise 1: Basic if Statement

**Objective:** Write your first conditional logic

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level5-exercise1
cd ~/projects/level5-exercise1
go mod init level5.example/exercise1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    ages := []int{5, 17, 18, 25, 65}

    for _, age := range ages {
        fmt.Printf("Age %d: ", age)
        if age >= 18 {
            fmt.Println("adult")
        }
        if age < 18 {
            fmt.Println("minor")
        }
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
Age 5: minor
Age 17: minor
Age 18: adult
Age 25: adult
Age 65: adult
```

**Learning Objectives:**
- ✅ Write a basic `if` statement
- ✅ Understand Go requires the opening brace on the same line

---

## Exercise 2: if/else

**Objective:** Use else to handle the opposite case

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level5-exercise2
cd ~/projects/level5-exercise2
go mod init level5.example/exercise2
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    numbers := []int{4, 7, 10, 15, 22}

    for _, n := range numbers {
        if n%2 == 0 {
            fmt.Printf("%d is even\n", n)
        } else {
            fmt.Printf("%d is odd\n", n)
        }
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
4 is even
7 is odd
10 is even
15 is odd
22 is even
```

**Learning Objectives:**
- ✅ Use if/else for binary decisions
- ✅ Combine Level 4's modulo operator with a condition

---

## Exercise 3: if/else if/else Chains

**Objective:** Build a grade calculator with multiple thresholds

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level5-exercise3
cd ~/projects/level5-exercise3
go mod init level5.example/exercise3
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func grade(score int) string {
    if score >= 90 {
        return "A"
    } else if score >= 80 {
        return "B"
    } else if score >= 70 {
        return "C"
    } else if score >= 60 {
        return "D"
    } else {
        return "F"
    }
}

func main() {
    scores := []int{95, 82, 71, 60, 45}

    for _, s := range scores {
        fmt.Printf("Score %d -> Grade %s\n", s, grade(s))
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
Score 95 -> Grade A
Score 82 -> Grade B
Score 71 -> Grade C
Score 60 -> Grade D
Score 45 -> Grade F
```

**Learning Objectives:**
- ✅ Chain multiple conditions with else if
- ✅ Understand that the first true condition wins
- ✅ Order thresholds from highest to lowest

---

## Exercise 4: if With Short Init

**Objective:** Scope a helper variable to the if/else blocks only

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level5-exercise4
cd ~/projects/level5-exercise4
go mod init level5.example/exercise4
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
    inputs := []string{"42", "abc", "17", "3.14"}

    fmt.Println("=== Parsing With Short Init ===")
    for _, s := range inputs {
        if value, err := strconv.Atoi(s); err == nil {
            fmt.Printf("%q -> parsed as %d\n", s, value)
        } else {
            fmt.Printf("%q -> failed to parse: %v\n", s, err)
        }
    }

    fmt.Println("\n=== Map Lookup With Short Init ===")
    scores := map[string]int{"Alice": 95, "Bob": 87}
    names := []string{"Alice", "Charlie"}

    for _, name := range names {
        if score, ok := scores[name]; ok {
            fmt.Printf("%s's score: %d\n", name, score)
        } else {
            fmt.Printf("%s not found\n", name)
        }
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
=== Parsing With Short Init ===
"42" -> parsed as 42
"abc" -> failed to parse: strconv.Atoi: parsing "abc": invalid syntax
"17" -> parsed as 17
"3.14" -> failed to parse: strconv.Atoi: parsing "3.14": invalid syntax

=== Map Lookup With Short Init ===
Alice's score: 95
Charlie not found
```

**Learning Objectives:**
- ✅ Use the if-with-short-init pattern
- ✅ Understand the short-init variable is scoped to if/else only
- ✅ Apply the pattern to error checking and map lookups (comma-ok idiom)

---

## Exercise 5: Basic switch

**Objective:** Replace an if/else if chain with a cleaner switch

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level5-exercise5
cd ~/projects/level5-exercise5
go mod init level5.example/exercise5
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func describeDay(day string) string {
    switch day {
    case "Saturday", "Sunday":
        return "Weekend"
    case "Monday":
        return "Start of the work week"
    case "Friday":
        return "Almost the weekend"
    default:
        return "Midweek"
    }
}

func main() {
    days := []string{"Monday", "Wednesday", "Friday", "Saturday", "Sunday"}

    for _, d := range days {
        fmt.Printf("%s: %s\n", d, describeDay(d))
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
Monday: Start of the work week
Wednesday: Midweek
Friday: Almost the weekend
Saturday: Weekend
Sunday: Weekend
```

**Learning Objectives:**
- ✅ Write a basic switch statement
- ✅ Use multiple comma-separated values in one case
- ✅ Use default for the fallback case

---

## Exercise 6: switch With Short Init

**Objective:** Combine switch with a short init statement

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level5-exercise6
cd ~/projects/level5-exercise6
go mod init level5.example/exercise6
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func greeting(hour int) string {
    switch h := hour; {
    case h < 12:
        return "Good morning"
    case h < 18:
        return "Good afternoon"
    default:
        return "Good evening"
    }
}

func main() {
    hours := []int{6, 9, 13, 19, 23}

    for _, h := range hours {
        fmt.Printf("Hour %2d: %s\n", h, greeting(h))
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
Hour  6: Good morning
Hour  9: Good morning
Hour 13: Good afternoon
Hour 19: Good evening
Hour 23: Good evening
```

**Learning Objectives:**
- ✅ Use a switch with a short init statement
- ✅ Understand switch with no value after the semicolon acts like if/else if

---

## Exercise 7: Conditionless switch

**Objective:** Rewrite a grade calculator using a conditionless switch

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level5-exercise7
cd ~/projects/level5-exercise7
go mod init level5.example/exercise7
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func grade(score int) string {
    switch {
    case score >= 90:
        return "A"
    case score >= 80:
        return "B"
    case score >= 70:
        return "C"
    case score >= 60:
        return "D"
    default:
        return "F"
    }
}

func main() {
    scores := []int{95, 82, 71, 60, 45}

    for _, s := range scores {
        fmt.Printf("Score %d -> Grade %s\n", s, grade(s))
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
Score 95 -> Grade A
Score 82 -> Grade B
Score 71 -> Grade C
Score 60 -> Grade D
Score 45 -> Grade F
```

**Learning Objectives:**
- ✅ Rewrite an if/else if chain as a conditionless switch
- ✅ Compare readability between the two forms

---

## Exercise 8: fallthrough

**Objective:** Understand explicit fallthrough behavior

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level5-exercise8
cd ~/projects/level5-exercise8
go mod init level5.example/exercise8
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func describeTier(tier int) {
    fmt.Printf("Tier %d benefits: ", tier)
    switch tier {
    case 1:
        fmt.Print("basic support")
        fallthrough
    case 2:
        fmt.Print(", priority queue")
        fallthrough
    case 3:
        fmt.Print(", dedicated agent")
    default:
        // no extra benefits below tier 1 or above tier 3
    }
    fmt.Println()
}

func main() {
    fmt.Println("=== Without Fallthrough (Normal switch) ===")
    switch n := 2; n {
    case 1:
        fmt.Println("one")
    case 2:
        fmt.Println("two")
    case 3:
        fmt.Println("three")
    }

    fmt.Println("\n=== With Fallthrough ===")
    for tier := 1; tier <= 3; tier++ {
        describeTier(tier)
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
=== Without Fallthrough (Normal switch) ===
two

=== With Fallthrough ===
Tier 1 benefits: basic support, priority queue, dedicated agent
Tier 2 benefits: , priority queue, dedicated agent
Tier 3 benefits: , dedicated agent
```

**Learning Objectives:**
- ✅ Confirm Go's switch does NOT fall through by default
- ✅ Use fallthrough to explicitly chain case bodies
- ✅ See how fallthrough runs the next case unconditionally

---

## Exercise 9: Avoiding Deep Nesting With Early Returns

**Objective:** Refactor nested conditions into flat early returns

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level5-exercise9
cd ~/projects/level5-exercise9
go mod init level5.example/exercise9
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type user struct {
    active        bool
    verified      bool
    hasPermission bool
}

// Deeply nested version
func processNested(u user) string {
    if u.active {
        if u.verified {
            if u.hasPermission {
                return "OK"
            } else {
                return "No permission"
            }
        } else {
            return "Not verified"
        }
    } else {
        return "Inactive"
    }
}

// Flattened version using early returns
func processFlat(u user) string {
    if !u.active {
        return "Inactive"
    }
    if !u.verified {
        return "Not verified"
    }
    if !u.hasPermission {
        return "No permission"
    }
    return "OK"
}

func main() {
    users := []user{
        {active: false, verified: false, hasPermission: false},
        {active: true, verified: false, hasPermission: false},
        {active: true, verified: true, hasPermission: false},
        {active: true, verified: true, hasPermission: true},
    }

    fmt.Println("=== Comparing Nested vs Flat ===")
    for _, u := range users {
        nested := processNested(u)
        flat := processFlat(u)
        match := nested == flat
        fmt.Printf("%+v -> nested=%q flat=%q match=%v\n", u, nested, flat, match)
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
=== Comparing Nested vs Flat ===
{active:false verified:false hasPermission:false} -> nested="Inactive" flat="Inactive" match=true
{active:true verified:false hasPermission:false} -> nested="Not verified" flat="Not verified" match=true
{active:true verified:true hasPermission:false} -> nested="No permission" flat="No permission" match=true
{active:true verified:true hasPermission:true} -> nested="OK" flat="OK" match=true
```

**Learning Objectives:**
- ✅ Recognize deeply nested conditions
- ✅ Refactor them into flat early returns
- ✅ Confirm both versions behave identically

---

## Exercise 10: Comprehensive Practice — Ticket Pricing System

**Objective:** Combine if/else, short-init, and switch in one realistic program

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level5-exercise10
cd ~/projects/level5-exercise10
go mod init level5.example/exercise10
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type visitor struct {
    name       string
    age        int
    dayOfWeek  string
    isMember   bool
}

func basePrice(age int) float64 {
    switch {
    case age < 5:
        return 0
    case age < 13:
        return 8
    case age < 65:
        return 15
    default:
        return 10
    }
}

func dayMultiplier(day string) float64 {
    switch day {
    case "Saturday", "Sunday":
        return 1.5
    default:
        return 1.0
    }
}

func ticketPrice(v visitor) (float64, string) {
    price := basePrice(v.age) * dayMultiplier(v.dayOfWeek)

    if v.isMember {
        price *= 0.8 // 20% member discount
    }

    var note string
    if price == 0 {
        note = "free admission"
    } else if v.isMember {
        note = "member discount applied"
    } else {
        note = "standard pricing"
    }

    return price, note
}

func main() {
    visitors := []visitor{
        {"Baby Sam", 3, "Monday", false},
        {"Teen Lee", 10, "Saturday", false},
        {"Adult Kim", 30, "Wednesday", true},
        {"Senior Cho", 70, "Sunday", false},
    }

    fmt.Println("=== Ticket Pricing System ===\n")

    var total float64
    for _, v := range visitors {
        price, note := ticketPrice(v)
        total += price
        fmt.Printf("%-12s age=%-3d %-9s member=%-5v -> $%.2f (%s)\n",
            v.name, v.age, v.dayOfWeek, v.isMember, price, note)
    }

    fmt.Printf("\nTotal for group: $%.2f\n", total)

    if total > 50 {
        fmt.Println("Group rate: consider a group discount next time!")
    } else {
        fmt.Println("Group rate: not yet eligible for a group discount")
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
=== Ticket Pricing System ===

Baby Sam     age=3   Monday    member=false -> $0.00 (free admission)
Teen Lee     age=10  Saturday  member=false -> $12.00 (standard pricing)
Adult Kim    age=30  Wednesday member=true  -> $12.00 (member discount applied)
Senior Cho   age=70  Sunday    member=false -> $15.00 (standard pricing)

Total for group: $39.00
Group rate: not yet eligible for a group discount
```

**Learning Objectives:**
- ✅ Combine if/else, switch, and conditionless switch in one program
- ✅ Use functions that return decisions based on conditions
- ✅ Apply everything from Level 4 (arithmetic, comparison) and Level 5 (conditions) together

---

## Bonus Challenges

### Challenge 1: Rock-Paper-Scissors Judge

Write a program that takes two hardcoded moves ("rock", "paper", "scissors") and decides the winner using a switch statement.

```bash
mkdir -p ~/projects/level5-bonus1
cd ~/projects/level5-bonus1
go mod init level5.example/bonus1
```

**Hints:**
- Use a `switch player1 { case ...: switch player2 { ... } }` nested structure, or a lookup table
- Handle the tie case explicitly
- Validate that both inputs are one of the three valid moves

### Challenge 2: FizzBuzz With switch

Implement classic FizzBuzz (1 to 30) using a conditionless switch instead of if/else if.

```bash
mkdir -p ~/projects/level5-bonus2
cd ~/projects/level5-bonus2
go mod init level5.example/bonus2
```

**Hints:**
- Check divisibility by 15 first (most specific), then 3, then 5
- Use `n%15 == 0`, `n%3 == 0`, `n%5 == 0` as your case conditions

### Challenge 3: Simple State Machine

Model a traffic light (Red -> Green -> Yellow -> Red) using a switch that determines the next state from the current one.

```bash
mkdir -p ~/projects/level5-bonus3
cd ~/projects/level5-bonus3
go mod init level5.example/bonus3
```

**Hints:**
- Use a named type `type light string` with constants for each state
- Write a `next(l light) light` function with a switch
- Loop and print several transitions in a row

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Write if, if/else, and if/else if chains
✅ Use the short-init form to scope helper variables
✅ Write switch statements, including multi-value cases
✅ Use conditionless switch as a clean if/else if replacement
✅ Use fallthrough deliberately
✅ Refactor deeply nested conditions into flat early returns
✅ Combine conditions with Level 4's operators in realistic programs

---

## Next Level

Level 6: Loops
- for loops (Go's only loop keyword)
- while-style and infinite loops
- break and continue
- Looping over collections

Great work! You're mastering Go! 🚀
