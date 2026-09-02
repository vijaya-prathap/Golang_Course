# Level 13: Structs - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

---

## Exercise 1: Defining and Initializing Structs

**Objective:** Define a struct and create values with every literal form

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level13-exercise1
cd ~/projects/level13-exercise1
go mod init level13.example/exercise1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type Person struct {
    Name string
    Age  int
    City string
}

func main() {
    fmt.Println("=== Keyed Struct Literal (preferred) ===")
    alice := Person{Name: "Alice", Age: 30, City: "Lisbon"}
    fmt.Println(alice)

    fmt.Println("\n=== Positional Struct Literal (fragile - avoid) ===")
    bob := Person{"Bob", 25, "Porto"}
    fmt.Println(bob)

    fmt.Println("\n=== Partial Keyed Literal (missing fields are zero) ===")
    carol := Person{Name: "Carol"}
    fmt.Println(carol)

    fmt.Println("\n=== Zero-Value Struct via var ===")
    var nobody Person
    fmt.Println(nobody)

    fmt.Println("\n=== %v vs %+v vs %#v ===")
    fmt.Printf("%v\n", alice)
    fmt.Printf("%+v\n", alice)
    fmt.Printf("%#v\n", alice)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Keyed Struct Literal (preferred) ===
{Alice 30 Lisbon}

=== Positional Struct Literal (fragile - avoid) ===
{Bob 25 Porto}

=== Partial Keyed Literal (missing fields are zero) ===
{Carol 0 }

=== Zero-Value Struct via var ===
{ 0 }

=== %v vs %+v vs %#v ===
{Alice 30 Lisbon}
{Name:Alice Age:30 City:Lisbon}
main.Person{Name:"Alice", Age:30, City:"Lisbon"}
```

(Note the space inside `{Carol 0 }` - that's the empty `City` string printing as nothing.)

**Learning Objectives:**
- ✅ Define a struct type with named fields
- ✅ Create structs with keyed, positional, and partial keyed literals
- ✅ Print structs with %v, %+v, and %#v

---

## Exercise 2: Exploring the Zero-Value Struct

**Objective:** Verify that every field of a zero-value struct gets its own zero value

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level13-exercise2
cd ~/projects/level13-exercise2
go mod init level13.example/exercise2
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type Profile struct {
    Username string
    Age      int
    Score    float64
    Active   bool
    Tags     []string
    Settings map[string]string
    Referrer *Profile
}

func main() {
    var p Profile

    fmt.Println("=== The Whole Zero-Value Struct ===")
    fmt.Printf("%+v\n", p)

    fmt.Println("\n=== Each Field's Zero Value ===")
    fmt.Printf("Username: %q\n", p.Username)
    fmt.Printf("Age:      %d\n", p.Age)
    fmt.Printf("Score:    %g\n", p.Score)
    fmt.Printf("Active:   %t\n", p.Active)
    fmt.Printf("Tags:     %v (nil? %t)\n", p.Tags, p.Tags == nil)
    fmt.Printf("Settings: %v (nil? %t)\n", p.Settings, p.Settings == nil)
    fmt.Printf("Referrer: %v (nil? %t)\n", p.Referrer, p.Referrer == nil)

    fmt.Println("\n=== A Zero-Value Struct Is Immediately Usable ===")
    p.Username = "gopher"
    p.Age = 3
    p.Tags = append(p.Tags, "newbie")
    fmt.Printf("%+v\n", p)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== The Whole Zero-Value Struct ===
{Username: Age:0 Score:0 Active:false Tags:[] Settings:map[] Referrer:<nil>}

=== Each Field's Zero Value ===
Username: ""
Age:      0
Score:    0
Active:   false
Tags:     [] (nil? true)
Settings: map[] (nil? true)
Referrer: <nil> (nil? true)

=== A Zero-Value Struct Is Immediately Usable ===
{Username:gopher Age:3 Score:0 Active:false Tags:[newbie] Settings:map[] Referrer:<nil>}
```

**Learning Objectives:**
- ✅ See that every field of a zero-value struct gets its type's zero value
- ✅ Recognize that slice/map/pointer fields start as nil
- ✅ Use a zero-value struct immediately (including append on its nil slice)

---

## Exercise 3: Structs Are Value Types (Copy vs Pointer)

**Objective:** Prove that assignment copies a struct, and that pointers share it

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level13-exercise3
cd ~/projects/level13-exercise3
go mod init level13.example/exercise3
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type Point struct {
    X int
    Y int
}

func tryToMutate(p Point) {
    p.X = -1 // modifies the COPY only
}

func mutate(p *Point) {
    p.Y = -1 // modifies the caller's struct
}

func main() {
    fmt.Println("=== Assignment Copies the Whole Struct ===")
    original := Point{X: 1, Y: 2}
    copied := original
    copied.X = 100

    fmt.Println("original:", original)
    fmt.Println("copied:  ", copied)

    fmt.Println("\n=== Passing to a Function Copies Too ===")
    tryToMutate(original)
    fmt.Println("after tryToMutate:", original)

    fmt.Println("\n=== Pointer to Struct: Shared Mutation (Level 12!) ===")
    ptr := &original
    ptr.X = 999 // shorthand for (*ptr).X = 999
    fmt.Println("original:", original)
    fmt.Println("via ptr: ", *ptr)

    fmt.Println("\n=== Pointer Parameter Mutates the Caller's Struct ===")
    mutate(&original)
    fmt.Println("after mutate:", original)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Assignment Copies the Whole Struct ===
original: {1 2}
copied:   {100 2}

=== Passing to a Function Copies Too ===
after tryToMutate: {1 2}

=== Pointer to Struct: Shared Mutation (Level 12!) ===
original: {999 2}
via ptr:  {999 2}

=== Pointer Parameter Mutates the Caller's Struct ===
after mutate: {999 -1}
```

**Learning Objectives:**
- ✅ Prove that struct assignment copies every field
- ✅ Prove that value parameters receive an independent copy
- ✅ Mutate a struct through a pointer, with auto-dereferenced field access

---

## Exercise 4: Comparing Structs (and When You Can't)

**Objective:** Compare structs with ==, then capture the compile error for an uncomparable struct

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level13-exercise4
cd ~/projects/level13-exercise4
go mod init level13.example/exercise4
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type Point struct {
    X int
    Y int
}

func main() {
    a := Point{X: 1, Y: 2}
    b := Point{X: 1, Y: 2}
    c := Point{X: 5, Y: 2}

    fmt.Println("=== Comparing Structs With == ===")
    fmt.Println("a == b:", a == b)
    fmt.Println("a == c:", a == c)
    fmt.Println("a != c:", a != c)

    fmt.Println("\n=== Zero-Value Check With == ===")
    var zero Point
    fmt.Println("a == zero:      ", a == zero)
    fmt.Println("zero == Point{}:", zero == Point{})

    fmt.Println("\n=== Comparable Structs Work as Map Keys ===")
    visits := make(map[Point]int)
    visits[Point{1, 2}]++
    visits[Point{1, 2}]++
    visits[Point{3, 4}]++
    fmt.Println("visits to {1 2}:", visits[Point{1, 2}])
    fmt.Println("visits to {3 4}:", visits[Point{3, 4}])
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Comparing Structs With == ===
a == b: true
a == c: false
a != c: true

=== Zero-Value Check With == ===
a == zero:       false
zero == Point{}: true

=== Comparable Structs Work as Map Keys ===
visits to {1 2}: 2
visits to {3 4}: 1
```

4. Now **replace** `main.go` with a struct containing a slice field - this one will NOT compile:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type Team struct {
    Name    string
    Members []string
}

func main() {
    t1 := Team{Name: "Go", Members: []string{"Alice", "Bob"}}
    t2 := Team{Name: "Go", Members: []string{"Alice", "Bob"}}
    fmt.Println(t1 == t2)
}
EOF
```

5. Try to run it and read the error carefully:

```bash
go run main.go
```

**Expected Output (a compile error - this is the point!):**

```
# command-line-arguments
./main.go:13:17: invalid operation: t1 == t2 (struct containing []string cannot be compared)
```

**Learning Objectives:**
- ✅ Compare structs field-by-field with == and !=
- ✅ Use == against the zero value and use comparable structs as map keys
- ✅ See the real compile error when a struct contains an uncomparable field (slice)

---

## Exercise 5: Nested Structs

**Objective:** Build a struct containing another struct and work with deep fields

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level13-exercise5
cd ~/projects/level13-exercise5
go mod init level13.example/exercise5
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type Address struct {
    Street string
    City   string
    Zip    string
}

type Person struct {
    Name    string
    Age     int
    Address Address
}

func main() {
    fmt.Println("=== Initializing a Nested Struct ===")
    alice := Person{
        Name: "Alice",
        Age:  30,
        Address: Address{
            Street: "12 Rose Lane",
            City:   "Lisbon",
            Zip:    "1100-048",
        },
    }
    fmt.Printf("%+v\n", alice)

    fmt.Println("\n=== Accessing Nested Fields With Chained Dots ===")
    fmt.Println("City:", alice.Address.City)
    fmt.Println("Zip: ", alice.Address.Zip)

    fmt.Println("\n=== Mutating a Nested Field ===")
    alice.Address.City = "Porto"
    fmt.Println("Moved to:", alice.Address.City)

    fmt.Println("\n=== The Nested Struct Copies With Its Parent ===")
    bob := alice
    bob.Name = "Bob"
    bob.Address.City = "Faro"
    fmt.Println("alice.Address.City:", alice.Address.City)
    fmt.Println("bob.Address.City:  ", bob.Address.City)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Initializing a Nested Struct ===
{Name:Alice Age:30 Address:{Street:12 Rose Lane City:Lisbon Zip:1100-048}}

=== Accessing Nested Fields With Chained Dots ===
City: Lisbon
Zip:  1100-048

=== Mutating a Nested Field ===
Moved to: Porto

=== The Nested Struct Copies With Its Parent ===
alice.Address.City: Porto
bob.Address.City:   Faro
```

**Learning Objectives:**
- ✅ Initialize a nested struct with nested keyed literals
- ✅ Read and mutate deep fields with chained dot access
- ✅ Confirm value semantics apply to the nested struct too

---

## Exercise 6: Embedded Structs and Field Promotion

**Objective:** Embed a struct, use promoted fields, and observe shadowing

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level13-exercise6
cd ~/projects/level13-exercise6
go mod init level13.example/exercise6
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type Animal struct {
    Name string
    Legs int
}

type Dog struct {
    Animal // embedded (anonymous field) - Name and Legs are PROMOTED
    Breed  string
}

type Cat struct {
    Animal
    Name string // explicit field SHADOWS the promoted Animal.Name
}

func main() {
    fmt.Println("=== Embedding and Field Promotion ===")
    d := Dog{
        Animal: Animal{Name: "Rex", Legs: 4},
        Breed:  "Border Collie",
    }
    fmt.Println("d.Name (promoted):        ", d.Name)
    fmt.Println("d.Legs (promoted):        ", d.Legs)
    fmt.Println("d.Animal.Name (full path):", d.Animal.Name)
    fmt.Println("d.Breed:                  ", d.Breed)

    fmt.Println("\n=== Promoted Fields Are Writable Too ===")
    d.Legs = 3 // exactly the same as d.Animal.Legs = 3
    fmt.Println("d.Animal.Legs:", d.Animal.Legs)

    fmt.Println("\n=== Shadowing: the Outer Field Wins ===")
    c := Cat{
        Animal: Animal{Name: "some animal", Legs: 4},
        Name:   "Whiskers",
    }
    fmt.Println("c.Name (outer field):           ", c.Name)
    fmt.Println("c.Animal.Name (still reachable):", c.Animal.Name)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Embedding and Field Promotion ===
d.Name (promoted):         Rex
d.Legs (promoted):         4
d.Animal.Name (full path): Rex
d.Breed:                   Border Collie

=== Promoted Fields Are Writable Too ===
d.Animal.Legs: 3

=== Shadowing: the Outer Field Wins ===
c.Name (outer field):            Whiskers
c.Animal.Name (still reachable): some animal
```

**Learning Objectives:**
- ✅ Embed a struct as an anonymous field
- ✅ Read AND write promoted fields directly on the outer struct
- ✅ See an explicit field shadow a promoted one, with the full path still reachable

---

## Exercise 7: Anonymous Structs

**Objective:** Use inline one-off struct types, including a table of rows

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level13-exercise7
cd ~/projects/level13-exercise7
go mod init level13.example/exercise7
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== A One-Off Anonymous Struct ===")
    server := struct {
        Host string
        Port int
    }{
        Host: "localhost",
        Port: 8080,
    }
    fmt.Printf("%+v\n", server)
    fmt.Printf("address: %s:%d\n", server.Host, server.Port)

    fmt.Println("\n=== A Slice of Anonymous Structs (table style) ===")
    tests := []struct {
        Input    int
        Expected int
    }{
        {Input: 2, Expected: 4},
        {Input: 3, Expected: 9},
        {Input: 5, Expected: 26},
    }

    for _, tc := range tests {
        got := tc.Input * tc.Input
        status := "PASS"
        if got != tc.Expected {
            status = "FAIL"
        }
        fmt.Printf("square(%d) = %d (want %d) %s\n", tc.Input, got, tc.Expected, status)
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
=== A One-Off Anonymous Struct ===
{Host:localhost Port:8080}
address: localhost:8080

=== A Slice of Anonymous Structs (table style) ===
square(2) = 4 (want 4) PASS
square(3) = 9 (want 9) PASS
square(5) = 25 (want 26) FAIL
```

(The last row FAILs on purpose - the table says 26, but 5*5 is 25. This is exactly how table-driven tests catch bugs in Level 26!)

**Learning Objectives:**
- ✅ Define and initialize an anonymous struct in one expression
- ✅ Build a slice of anonymous structs as a data table
- ✅ Process table rows with the Level 6 range pattern

---

## Exercise 8: Struct Tags

**Objective:** Attach tags to fields and read them back via reflection

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level13-exercise8
cd ~/projects/level13-exercise8
go mod init level13.example/exercise8
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "reflect"
)

type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age,omitempty"`
}

func main() {
    fmt.Println("=== Tags Don't Change Normal Behavior ===")
    u := User{Name: "Alice", Email: "alice@example.com", Age: 30}
    fmt.Printf("%+v\n", u)

    fmt.Println("\n=== Reading Tags via Reflection ===")
    t := reflect.TypeOf(u)
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        fmt.Printf("%-5s json tag: %q\n", field.Name, field.Tag.Get("json"))
    }

    fmt.Println("\n=== The Raw Tag String ===")
    ageField, _ := t.FieldByName("Age")
    fmt.Println(ageField.Tag)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Tags Don't Change Normal Behavior ===
{Name:Alice Email:alice@example.com Age:30}

=== Reading Tags via Reflection ===
Name  json tag: "name"
Email json tag: "email"
Age   json tag: "age,omitempty"

=== The Raw Tag String ===
json:"age,omitempty"
```

**Learning Objectives:**
- ✅ Write struct tags with the backtick `key:"value"` syntax
- ✅ Confirm tags are inert metadata - normal code ignores them
- ✅ See how libraries (like encoding/json in Level 19) read tags via reflection

---

## Exercise 9: The Constructor Convention (NewX)

**Objective:** Write a NewX constructor that validates input and returns a pointer

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level13-exercise9
cd ~/projects/level13-exercise9
go mod init level13.example/exercise9
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type BankAccount struct {
    Owner   string
    Balance float64
    Active  bool
}

// NewBankAccount is Go's constructor convention: a New<Type> function
// that returns a pointer to a fully initialized struct.
func NewBankAccount(owner string, openingBalance float64) *BankAccount {
    if openingBalance < 0 {
        openingBalance = 0
    }
    return &BankAccount{
        Owner:   owner,
        Balance: openingBalance,
        Active:  true,
    }
}

func deposit(a *BankAccount, amount float64) {
    a.Balance += amount
}

func main() {
    fmt.Println("=== Using the Constructor ===")
    acct := NewBankAccount("Alice", 250.75)
    fmt.Printf("%+v\n", *acct)

    fmt.Println("\n=== The Constructor Enforces Sensible Defaults ===")
    weird := NewBankAccount("Bob", -100)
    fmt.Printf("%+v\n", *weird)

    fmt.Println("\n=== A Pointer Return Means Shared Mutation ===")
    deposit(acct, 49.25)
    fmt.Printf("after deposit: %+v\n", *acct)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Using the Constructor ===
{Owner:Alice Balance:250.75 Active:true}

=== The Constructor Enforces Sensible Defaults ===
{Owner:Bob Balance:0 Active:true}

=== A Pointer Return Means Shared Mutation ===
after deposit: {Owner:Alice Balance:300 Active:true}
```

**Learning Objectives:**
- ✅ Write a New<Type> constructor returning *Type
- ✅ Centralize validation and defaults in one place
- ✅ Combine the constructor with pointer-based mutation from Level 12

---

## Exercise 10: Comprehensive Practice — Library Inventory

**Objective:** Combine embedding, a constructor, and Level 6 loop patterns in one model

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level13-exercise10
cd ~/projects/level13-exercise10
go mod init level13.example/exercise10
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type Metadata struct {
    Publisher string
    Year      int
}

type Book struct {
    Metadata // embedded: Publisher and Year are promoted
    Title    string
    Author   string
    Copies   int
}

func NewBook(title, author, publisher string, year, copies int) *Book {
    if copies < 0 {
        copies = 0
    }
    return &Book{
        Metadata: Metadata{Publisher: publisher, Year: year},
        Title:    title,
        Author:   author,
        Copies:   copies,
    }
}

func main() {
    library := []Book{
        *NewBook("The Go Programming Language", "Donovan & Kernighan", "Addison-Wesley", 2015, 3),
        *NewBook("Learning Go", "Jon Bodner", "O'Reilly", 2021, 2),
        *NewBook("Go in Action", "Kennedy & Ketelsen", "Manning", 2015, 0),
        *NewBook("100 Go Mistakes", "Teiva Harsanyi", "Manning", 2022, 5),
    }

    fmt.Println("=== Full Catalog ===")
    for i, b := range library {
        fmt.Printf("%d. %q by %s (%s, %d) - %d copies\n",
            i+1, b.Title, b.Author, b.Publisher, b.Year, b.Copies)
    }

    fmt.Println("\n=== Filter: In-Stock Books ===")
    var inStock []string
    for _, b := range library {
        if b.Copies > 0 {
            inStock = append(inStock, b.Title)
        }
    }
    fmt.Println(inStock)

    fmt.Println("\n=== Accumulate: Total Copies ===")
    total := 0
    for _, b := range library {
        total += b.Copies
    }
    fmt.Println("total copies:", total)

    fmt.Println("\n=== Find: First Book From Manning ===")
    for _, b := range library {
        if b.Publisher == "Manning" { // promoted field access!
            fmt.Println("found:", b.Title)
            break
        }
    }

    fmt.Println("\n=== Count: Books Per Publisher ===")
    counts := make(map[string]int)
    for _, b := range library {
        counts[b.Publisher]++
    }
    fmt.Println("Addison-Wesley:", counts["Addison-Wesley"])
    fmt.Println("Manning:       ", counts["Manning"])
    fmt.Println("O'Reilly:      ", counts["O'Reilly"])
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Full Catalog ===
1. "The Go Programming Language" by Donovan & Kernighan (Addison-Wesley, 2015) - 3 copies
2. "Learning Go" by Jon Bodner (O'Reilly, 2021) - 2 copies
3. "Go in Action" by Kennedy & Ketelsen (Manning, 2015) - 0 copies
4. "100 Go Mistakes" by Teiva Harsanyi (Manning, 2022) - 5 copies

=== Filter: In-Stock Books ===
[The Go Programming Language Learning Go 100 Go Mistakes]

=== Accumulate: Total Copies ===
total copies: 10

=== Find: First Book From Manning ===
found: Go in Action

=== Count: Books Per Publisher ===
Addison-Wesley: 1
Manning:        2
O'Reilly:       1
```

**Learning Objectives:**
- ✅ Model real data with an embedded Metadata struct and promoted fields
- ✅ Build values through a NewX constructor
- ✅ Process a slice of structs with filter/accumulate/find/count loop patterns

---

## Bonus Challenges

### Challenge 1: 2D Vector Type

Build a `Vector` struct with `X, Y float64` and helper functions: `Add(a, b Vector) Vector`, `Scale(v Vector, factor float64) Vector`, and `Distance(a, b Vector) float64`. Print the results of a few operations.

```bash
mkdir -p ~/projects/level13-bonus1
cd ~/projects/level13-bonus1
go mod init level13.example/bonus1
```

**Hints:**
- `Add` and `Scale` return NEW vectors - value semantics make this natural
- For `Distance`, use `math.Hypot(b.X-a.X, b.Y-a.Y)` from the `math` package
- Because all fields are comparable, you can check `a == b` in your experiments

### Challenge 2: Config Struct With Defaults

Model a `ServerConfig` struct (`Host string`, `Port int`, `TimeoutSeconds int`, `Debug bool`) and write `NewServerConfig(host string) *ServerConfig` that fills in sensible defaults (port 8080, timeout 30, debug off) whenever the zero value wouldn't be useful.

```bash
mkdir -p ~/projects/level13-bonus2
cd ~/projects/level13-bonus2
go mod init level13.example/bonus2
```

**Hints:**
- Default `Host` to `"localhost"` when the caller passes an empty string
- Print the config with `%+v` before and after manual field overrides
- Ask yourself: which fields have a *usable* zero value (Debug!) and which don't (Port)?

### Challenge 3: Company Org Model

Model a company: a `Person` struct (`Name`, `Age`), an `Employee` struct that **embeds** `Person` and adds `Title string` and `Salary float64`, and a `Company` struct with a `Name` and a slice of employees. Print a roster, compute total payroll, and count employees per title.

```bash
mkdir -p ~/projects/level13-bonus3
cd ~/projects/level13-bonus3
go mod init level13.example/bonus3
```

**Hints:**
- Thanks to promotion you can write `e.Name` even though `Name` lives in `Person`
- Use `map[string]int` for the per-title count (Level 10 + Level 6 patterns)
- Bonus points: write `NewEmployee` and reject negative salaries

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Define struct types and create them with keyed, positional, and partial literals
✅ Predict every field of a zero-value struct
✅ Prove copy-on-assignment and choose pointers when sharing is needed
✅ Compare structs with == and explain exactly when comparison won't compile
✅ Build and traverse nested structs with chained dot access
✅ Use embedding, field promotion, and understand shadowing
✅ Reach for anonymous structs for one-off tables of data
✅ Read struct tag syntax and know who consumes tags (and when - Level 19)
✅ Write NewX constructors that validate and default fields

---

## Next Level

Level 14: Methods
- Attaching behavior to your struct types
- Value receivers vs pointer receivers
- Method promotion through embedding
- Method sets - the doorway to interfaces

Great work! Your data has structure now - next it gets behavior! 🚀
