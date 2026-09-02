# Level 14: Methods - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

---

## Exercise 1: Your First Methods

**Objective:** Declare methods on a struct and compare them with plain functions

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level14-exercise1
cd ~/projects/level14-exercise1
go mod init level14.example/exercise1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type Rectangle struct {
    Width  float64
    Height float64
}

// Area is a METHOD: a function with a receiver (r Rectangle)
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// area is a plain FUNCTION doing the same job
func area(r Rectangle) float64 {
    return r.Width * r.Height
}

func main() {
    rect := Rectangle{Width: 4, Height: 3}

    fmt.Println("=== Calling Methods ===")
    fmt.Println("Area:", rect.Area())
    fmt.Println("Perimeter:", rect.Perimeter())

    fmt.Println("\n=== Same Logic as a Plain Function ===")
    fmt.Println("area(rect):", area(rect))

    fmt.Println("\n=== Each Value Uses Its Own Fields ===")
    small := Rectangle{Width: 2, Height: 1}
    big := Rectangle{Width: 10, Height: 8}
    fmt.Println("small.Area():", small.Area())
    fmt.Println("big.Area():", big.Area())
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Calling Methods ===
Area: 12
Perimeter: 14

=== Same Logic as a Plain Function ===
area(rect): 12

=== Each Value Uses Its Own Fields ===
small.Area(): 2
big.Area(): 80
```

**Learning Objectives:**
- ✅ Declare a method with a receiver
- ✅ See that methods and functions have the same power, different call syntax
- ✅ Use a short, consistent receiver name (r, not this/self)

---

## Exercise 2: Value Receivers - Proof of the Copy

**Objective:** Prove that a value receiver operates on a copy

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level14-exercise2
cd ~/projects/level14-exercise2
go mod init level14.example/exercise2
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type Counter struct {
    Count int
}

// Value receiver: c is a COPY of the caller's Counter
func (c Counter) Increment() {
    c.Count++
    fmt.Println("  inside Increment: c.Count =", c.Count)
}

func main() {
    c := Counter{Count: 0}

    fmt.Println("=== One Call ===")
    fmt.Println("before:", c.Count)
    c.Increment()
    fmt.Println("after: ", c.Count)

    fmt.Println("\n=== Three More Calls ===")
    c.Increment()
    c.Increment()
    c.Increment()
    fmt.Println("final: ", c.Count)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== One Call ===
before: 0
  inside Increment: c.Count = 1
after:  0

=== Three More Calls ===
  inside Increment: c.Count = 1
  inside Increment: c.Count = 1
  inside Increment: c.Count = 1
final:  0
```

**Learning Objectives:**
- ✅ Prove the copy IS incremented inside the method (prints 1)
- ✅ Prove the caller's value is untouched (still 0 after every call)
- ✅ Recognize the silent no-op mutation bug

---

## Exercise 3: Pointer Receivers - Mutating the Original

**Objective:** Contrast value and pointer receivers side by side

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level14-exercise3
cd ~/projects/level14-exercise3
go mod init level14.example/exercise3
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type Player struct {
    Name  string
    Score int
}

// Value receiver: mutates a COPY (the change is lost)
func (p Player) AddPointsValue(n int) {
    p.Score += n
}

// Pointer receiver: mutates the ORIGINAL
func (p *Player) AddPointsPointer(n int) {
    p.Score += n
}

func main() {
    fmt.Println("=== Value Receiver (copy) ===")
    p1 := Player{Name: "Alice", Score: 100}
    p1.AddPointsValue(50)
    fmt.Println(p1.Name, "score:", p1.Score)

    fmt.Println("\n=== Pointer Receiver (original) ===")
    p2 := Player{Name: "Bob", Score: 100}
    p2.AddPointsPointer(50)
    fmt.Println(p2.Name, "score:", p2.Score)

    fmt.Println("\n=== Side by Side, Five Calls Each ===")
    a := Player{Name: "ValueGal", Score: 0}
    b := Player{Name: "PointerPal", Score: 0}
    for i := 0; i < 5; i++ {
        a.AddPointsValue(10)
        b.AddPointsPointer(10)
    }
    fmt.Println(a.Name, "final score:", a.Score)
    fmt.Println(b.Name, "final score:", b.Score)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Value Receiver (copy) ===
Alice score: 100

=== Pointer Receiver (original) ===
Bob score: 150

=== Side by Side, Five Calls Each ===
ValueGal final score: 0
PointerPal final score: 50
```

**Learning Objectives:**
- ✅ Declare a pointer-receiver method with (p *Player)
- ✅ See mutation stick through a pointer receiver and vanish through a value receiver
- ✅ Understand why mutating methods need pointer receivers

---

## Exercise 4: Auto-Address and Auto-Dereference

**Objective:** See Go's receiver sugar in both directions

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level14-exercise4
cd ~/projects/level14-exercise4
go mod init level14.example/exercise4
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type Account struct {
    Owner   string
    Balance float64
}

func (a *Account) Deposit(amount float64) {
    a.Balance += amount
}

func (a Account) Report() string {
    return fmt.Sprintf("%s: $%.2f", a.Owner, a.Balance)
}

func main() {
    fmt.Println("=== Value Variable, Pointer Method ===")
    acct := Account{Owner: "Alice", Balance: 100}
    acct.Deposit(50)    // sugar: Go rewrites this as (&acct).Deposit(50)
    (&acct).Deposit(25) // the explicit long form (never written in practice)
    fmt.Println(acct.Report())

    fmt.Println("\n=== Pointer Variable, Value Method ===")
    ptr := &Account{Owner: "Bob", Balance: 10}
    fmt.Println(ptr.Report())    // sugar: Go rewrites this as (*ptr).Report()
    fmt.Println((*ptr).Report()) // the explicit long form

    fmt.Println("\n=== Pointer Variable, Pointer Method ===")
    ptr.Deposit(90)
    fmt.Println(ptr.Report())
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Value Variable, Pointer Method ===
Alice: $175.00

=== Pointer Variable, Value Method ===
Bob: $10.00
Bob: $10.00

=== Pointer Variable, Pointer Method ===
Bob: $100.00
```

**Learning Objectives:**
- ✅ Call a pointer-receiver method on an addressable value (auto-address)
- ✅ Call a value-receiver method through a pointer (auto-dereference)
- ✅ Recognize the explicit forms (&v). and (*p). that the sugar replaces

---

## Exercise 5: The Non-Addressable Compile Error

**Objective:** Hit the real compile error where the auto-address sugar cannot help, then fix it both ways

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level14-exercise5
cd ~/projects/level14-exercise5
go mod init level14.example/exercise5
```

2. Create `main.go` (this version intentionally does NOT compile):

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type Player struct {
    Name  string
    Score int
}

func (p *Player) AddPoints(n int) {
    p.Score += n
}

func main() {
    players := map[string]Player{
        "alice": {Name: "Alice", Score: 10},
    }

    players["alice"].AddPoints(5) // map values are NOT addressable!

    fmt.Println(players)
}
EOF
```

3. Try to run it and read the error carefully:

```bash
go run main.go
```

**Expected Output (a compile error - this is the point!):**

```
# command-line-arguments
./main.go:19:22: cannot call pointer method AddPoints on Player
```

Go cannot take the address of a struct stored inside a map, so it has no `*Player` to pass to `AddPoints`.

4. Now replace `main.go` with the fixed version:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type Player struct {
    Name  string
    Score int
}

func (p *Player) AddPoints(n int) {
    p.Score += n
}

func main() {
    fmt.Println("=== Fix 1: Copy Out, Modify, Store Back ===")
    players := map[string]Player{
        "alice": {Name: "Alice", Score: 10},
    }
    p := players["alice"] // copy out (p is an addressable variable)
    p.AddPoints(5)        // works: Go can take &p
    players["alice"] = p  // store the modified copy back
    fmt.Println(players["alice"])

    fmt.Println("\n=== Fix 2: Store Pointers in the Map ===")
    roster := map[string]*Player{
        "bob": {Name: "Bob", Score: 20},
    }
    roster["bob"].AddPoints(5) // works: the map value IS a pointer
    fmt.Println(*roster["bob"])
}
EOF
```

5. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Fix 1: Copy Out, Modify, Store Back ===
{Alice 15}

=== Fix 2: Store Pointers in the Map ===
{Bob 25}
```

**Learning Objectives:**
- ✅ Recognize the "cannot call pointer method ... on ..." compile error
- ✅ Understand WHY map elements are not addressable
- ✅ Apply both fixes: copy-out/store-back and a map of pointers

---

## Exercise 6: The Consistency Rule Refactor

**Objective:** Debug a type with mixed receivers, then refactor to all-pointer receivers

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level14-exercise6
cd ~/projects/level14-exercise6
go mod init level14.example/exercise6
```

2. Create `main.go` (this version compiles but has a BUG - can you spot it?):

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type BankAccount struct {
    Owner   string
    Balance float64
}

// Pointer receiver: works correctly
func (b *BankAccount) Deposit(amount float64) {
    b.Balance += amount
}

// Value receiver: BUG - subtracts from a copy!
func (b BankAccount) Withdraw(amount float64) {
    b.Balance -= amount
}

func main() {
    acct := BankAccount{Owner: "Alice", Balance: 100}

    acct.Deposit(50)
    fmt.Println("after Deposit(50): ", acct.Balance)

    acct.Withdraw(30)
    fmt.Println("after Withdraw(30):", acct.Balance)
}
EOF
```

3. Run it and study the wrong result:

```bash
go run main.go
```

**Expected Output (note the bug - the withdrawal did nothing!):**

```
after Deposit(50):  150
after Withdraw(30): 150
```

4. Refactor: give ALL methods pointer receivers (the consistency rule). Replace `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type BankAccount struct {
    Owner   string
    Balance float64
}

// ALL methods now use pointer receivers - consistent!
func (b *BankAccount) Deposit(amount float64) {
    b.Balance += amount
}

func (b *BankAccount) Withdraw(amount float64) {
    b.Balance -= amount
}

func (b *BankAccount) Report() string {
    return fmt.Sprintf("%s has $%.2f", b.Owner, b.Balance)
}

func main() {
    acct := BankAccount{Owner: "Alice", Balance: 100}

    acct.Deposit(50)
    fmt.Println("after Deposit(50): ", acct.Balance)

    acct.Withdraw(30)
    fmt.Println("after Withdraw(30):", acct.Balance)

    fmt.Println(acct.Report())
}
EOF
```

5. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
after Deposit(50):  150
after Withdraw(30): 120
Alice has $120.00
```

**Learning Objectives:**
- ✅ Spot the silent mixed-receiver bug (no compile error, wrong behavior)
- ✅ Apply the consistency rule: one pointer receiver → all pointer receivers
- ✅ See that even read-only methods (Report) follow the rule for uniformity

---

## Exercise 7: Methods on Named Non-Struct Types

**Objective:** Attach methods to named float64 and string types

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level14-exercise7
cd ~/projects/level14-exercise7
go mod init level14.example/exercise7
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "strings"
)

type Celsius float64
type Fahrenheit float64

func (c Celsius) ToFahrenheit() Fahrenheit {
    return Fahrenheit(c*9/5 + 32)
}

func (c Celsius) IsFreezing() bool {
    return c <= 0
}

func (f Fahrenheit) ToCelsius() Celsius {
    return Celsius((f - 32) * 5 / 9)
}

type Path string

func (p Path) Base() string {
    s := string(p)
    if i := strings.LastIndex(s, "/"); i >= 0 {
        return s[i+1:]
    }
    return s
}

func (p Path) Ext() string {
    base := p.Base()
    if i := strings.LastIndex(base, "."); i >= 0 {
        return base[i:]
    }
    return ""
}

func main() {
    fmt.Println("=== Temperature Methods ===")
    boiling := Celsius(100)
    cold := Celsius(-5)
    fmt.Printf("%.0f°C = %.0f°F\n", float64(boiling), float64(boiling.ToFahrenheit()))
    fmt.Printf("%.0f°C freezing? %v\n", float64(cold), cold.IsFreezing())
    fmt.Printf("98.6°F = %.1f°C\n", float64(Fahrenheit(98.6).ToCelsius()))

    fmt.Println("\n=== Path Methods ===")
    p := Path("/home/user/report.pdf")
    fmt.Println("Base:", p.Base())
    fmt.Println("Ext: ", p.Ext())
    fmt.Println("Bare:", Path("notes").Ext() == "")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Temperature Methods ===
100°C = 212°F
-5°C freezing? true
98.6°F = 37.0°C

=== Path Methods ===
Base: report.pdf
Ext:  .pdf
Bare: true
```

**Learning Objectives:**
- ✅ Define methods on named float64 and string types
- ✅ Deepen Level 3's type Age int pattern into full mini-APIs
- ✅ Convert between named types explicitly (Celsius ↔ Fahrenheit)

---

## Exercise 8: Method Promotion and Overriding

**Objective:** Promote an embedded type's methods and override one

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level14-exercise8
cd ~/projects/level14-exercise8
go mod init level14.example/exercise8
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type Animal struct {
    Name string
}

func (a Animal) Describe() string {
    return a.Name + " is an animal"
}

func (a Animal) Speak() string {
    return a.Name + " makes a sound"
}

type Dog struct {
    Animal // embedded: Animal's methods are PROMOTED to Dog
    Breed  string
}

// A same-name method on the outer type OVERRIDES the promoted one
func (d Dog) Speak() string {
    return d.Name + " says Woof!"
}

func main() {
    d := Dog{Animal: Animal{Name: "Rex"}, Breed: "Beagle"}

    fmt.Println("=== Promoted Method ===")
    fmt.Println(d.Describe()) // defined on Animal, called on Dog

    fmt.Println("\n=== Overridden Method ===")
    fmt.Println(d.Speak()) // Dog's version wins

    fmt.Println("\n=== Reaching the Embedded Version Explicitly ===")
    fmt.Println(d.Animal.Speak()) // the shadowed Animal version

    fmt.Println("\n=== Promotion Works Through Pointers Too ===")
    pd := &Dog{Animal: Animal{Name: "Buddy"}, Breed: "Corgi"}
    fmt.Println(pd.Describe())
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Promoted Method ===
Rex is an animal

=== Overridden Method ===
Rex says Woof!

=== Reaching the Embedded Version Explicitly ===
Rex makes a sound

=== Promotion Works Through Pointers Too ===
Buddy is an animal
```

**Learning Objectives:**
- ✅ Call a promoted method directly on the outer type
- ✅ Override a promoted method with a same-name outer method
- ✅ Reach the shadowed embedded version via d.Animal.Speak()

---

## Exercise 9: Method Values and Method Expressions

**Objective:** Treat methods as values, with and without a bound receiver

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level14-exercise9
cd ~/projects/level14-exercise9
go mod init level14.example/exercise9
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type Person struct {
    Name string
}

func (p Person) Greet() string {
    return "Hello, I'm " + p.Name
}

func (p *Person) Rename(newName string) {
    p.Name = newName
}

func main() {
    p := Person{Name: "Alice"}

    fmt.Println("=== Method Value (receiver bound now) ===")
    greet := p.Greet // binds a COPY of p as the receiver
    fmt.Println(greet())

    fmt.Println("\n=== The Bound Receiver Is a Snapshot ===")
    p.Name = "Alicia"
    fmt.Println(greet())   // still the old copy!
    fmt.Println(p.Greet()) // a fresh call sees the new name

    fmt.Println("\n=== Method Expression (receiver as 1st argument) ===")
    greetFn := Person.Greet // func(Person) string
    fmt.Println(greetFn(Person{Name: "Bob"}))
    fmt.Println(greetFn(p))

    fmt.Println("\n=== Pointer Method Expression ===")
    renameFn := (*Person).Rename // func(*Person, string)
    renameFn(&p, "Allie")
    fmt.Println(p.Greet())
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Method Value (receiver bound now) ===
Hello, I'm Alice

=== The Bound Receiver Is a Snapshot ===
Hello, I'm Alice
Hello, I'm Alicia

=== Method Expression (receiver as 1st argument) ===
Hello, I'm Bob
Hello, I'm Alicia

=== Pointer Method Expression ===
Hello, I'm Allie
```

**Learning Objectives:**
- ✅ Bind a method value with p.Greet and call it later
- ✅ Prove a value-receiver method value snapshots its receiver
- ✅ Use method expressions that take the receiver as the first argument

---

## Exercise 10: Comprehensive Practice — ShoppingCart

**Objective:** Build a complete type with pointer-receiver methods, a constructor, and a String() method

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level14-exercise10
cd ~/projects/level14-exercise10
go mod init level14.example/exercise10
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type Item struct {
    Name     string
    Price    float64
    Quantity int
}

type ShoppingCart struct {
    Items []Item
}

func NewShoppingCart() *ShoppingCart {
    return &ShoppingCart{}
}

func (c *ShoppingCart) AddItem(name string, price float64, quantity int) {
    for i := range c.Items {
        if c.Items[i].Name == name {
            c.Items[i].Quantity += quantity
            return
        }
    }
    c.Items = append(c.Items, Item{Name: name, Price: price, Quantity: quantity})
}

func (c *ShoppingCart) RemoveItem(name string) bool {
    for i := range c.Items {
        if c.Items[i].Name == name {
            c.Items = append(c.Items[:i], c.Items[i+1:]...)
            return true
        }
    }
    return false
}

func (c *ShoppingCart) Total() float64 {
    total := 0.0
    for _, item := range c.Items {
        total += item.Price * float64(item.Quantity)
    }
    return total
}

// String() will make ShoppingCart satisfy fmt.Stringer in Level 15!
func (c *ShoppingCart) String() string {
    if len(c.Items) == 0 {
        return "Cart: (empty)"
    }
    s := "Cart:\n"
    for _, item := range c.Items {
        s += fmt.Sprintf("  %-6s x%d @ $%.2f = $%.2f\n",
            item.Name, item.Quantity, item.Price, item.Price*float64(item.Quantity))
    }
    s += fmt.Sprintf("  TOTAL: $%.2f", c.Total())
    return s
}

func main() {
    cart := NewShoppingCart()
    fmt.Println(cart)

    fmt.Println("\n=== Adding Items ===")
    cart.AddItem("apple", 0.50, 4)
    cart.AddItem("bread", 2.25, 1)
    cart.AddItem("milk", 3.10, 2)
    cart.AddItem("apple", 0.50, 2) // merges with the existing apples
    fmt.Println(cart)

    fmt.Println("\n=== Removing an Item ===")
    fmt.Println("removed bread?", cart.RemoveItem("bread"))
    fmt.Println("removed candy?", cart.RemoveItem("candy"))
    fmt.Println(cart)

    fmt.Println("\n=== Total ===")
    fmt.Printf("Total: $%.2f\n", cart.Total())
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
Cart: (empty)

=== Adding Items ===
Cart:
  apple  x6 @ $0.50 = $3.00
  bread  x1 @ $2.25 = $2.25
  milk   x2 @ $3.10 = $6.20
  TOTAL: $11.45

=== Removing an Item ===
removed bread? true
removed candy? false
Cart:
  apple  x6 @ $0.50 = $3.00
  milk   x2 @ $3.10 = $6.20
  TOTAL: $9.20

=== Total ===
Total: $9.20
```

**Learning Objectives:**
- ✅ Build a full type API: constructor + consistent pointer-receiver methods
- ✅ Mutate a slice field (append/remove) through pointer receivers
- ✅ Give a type a String() method that fmt.Println picks up automatically
- ✅ Preview fmt.Stringer - the first interface you'll meet in Level 15

---

## Bonus Challenges

### Challenge 1: Matrix Arithmetic

Build a `Matrix` type (`type Matrix [][]float64`) with methods `Add(other Matrix) Matrix`, `Scale(factor float64) Matrix`, and `Transpose() Matrix`. Each method should return a NEW matrix (value semantics - like the temperature types).

```bash
mkdir -p ~/projects/level14-bonus1
cd ~/projects/level14-bonus1
go mod init level14.example/bonus1
```

**Hints:**
- Write a helper `NewMatrix(rows, cols int) Matrix` that allocates the nested slices
- `Add` should check that dimensions match and return the sum element by element
- In `Transpose`, result[c][r] = m[r][c]
- These methods only READ the receiver, so value receivers are the right choice

### Challenge 2: StringSet

Build a `StringSet` type on `map[string]struct{}` with methods `Add(s string)`, `Has(s string) bool`, `Remove(s string)`, and `Len() int`.

```bash
mkdir -p ~/projects/level14-bonus2
cd ~/projects/level14-bonus2
go mod init level14.example/bonus2
```

**Hints:**
- `type StringSet map[string]struct{}` - the empty struct{} value takes zero bytes
- A map is reference-like (Level 10), so even a VALUE receiver can add/delete entries - but you still need a constructor (`func NewStringSet() StringSet { return StringSet{} }`) so the map is initialized
- `Has` uses the comma-ok form: `_, ok := s[key]; return ok`
- Bonus points: add a `Values() []string` method that returns sorted keys

### Challenge 3: Chainable Builder

Build a `PizzaBuilder` whose methods return `*PizzaBuilder` so calls chain fluently: `NewPizza().Size("large").Topping("olives").Topping("basil").Build()`.

```bash
mkdir -p ~/projects/level14-bonus3
cd ~/projects/level14-bonus3
go mod init level14.example/bonus3
```

**Hints:**
- Each setter mutates the receiver, then ends with `return b`
- Signature pattern: `func (b *PizzaBuilder) Size(s string) *PizzaBuilder`
- Because every method already returns `*PizzaBuilder`, the whole chain works on one underlying builder - no copies
- `Build()` returns the finished `Pizza` struct (and could validate required fields)

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Declare methods with value and pointer receivers
✅ Prove value receivers copy and pointer receivers mutate the original
✅ Use (and see through) the auto-address/auto-dereference sugar
✅ Recognize and fix the non-addressable "cannot call pointer method" error
✅ Apply the consistency rule: one pointer receiver → all pointer receivers
✅ Attach methods to named non-struct types like Celsius and Path
✅ Use method promotion through embedding, and override promoted methods
✅ Bind method values and use method expressions
✅ Design a full type API with a constructor and String() method

---

## Next Level

Level 15: Interfaces
- Defining interfaces and implicit satisfaction
- Method sets in action: which types satisfy which interfaces
- The empty interface and type assertions
- fmt.Stringer and the standard library's small-interface style

Great work! Your types have behavior now! 🚀
