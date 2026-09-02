# Level 13: Quick Reference Card

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir myapp && cd myapp
go mod init myapp

# Create main.go
cat > main.go << 'EOF'
package main

import "fmt"

type Person struct {
    Name string
    Age  int
}

func main() {
    alice := Person{Name: "Alice", Age: 30}
    fmt.Println(alice)

    bob := alice // copy, not shared
    bob.Age = 99
    fmt.Println(alice.Age, bob.Age) // 30 99
}
EOF

# Run
go run main.go
```

---

## 📋 Defining and Creating Structs

```go
type Person struct {
    Name string
    Age  int
    City string
}

var p Person                                  // zero value: {"" 0 ""}
alice := Person{Name: "Alice", Age: 30}       // keyed - preferred, can skip fields
bob := Person{"Bob", 25, "Porto"}             // positional - avoid, needs ALL fields
ptr := &Person{Name: "Carol"}                 // pointer to a fresh struct
```

---

## 🎯 Printing Structs

```go
fmt.Printf("%v\n", p)   // {Alice 30 Lisbon}
fmt.Printf("%+v\n", p)  // {Name:Alice Age:30 City:Lisbon}  <- use this for debugging
fmt.Printf("%#v\n", p)  // main.Person{Name:"Alice", Age:30, City:"Lisbon"}
```

---

## 🔁 Value Semantics: Copy vs Pointer

```go
b := a          // COPIES every field - b and a are independent
b.Field = x     // does NOT affect a

ptr := &a       // ptr points AT a - no copy
ptr.Field = x   // DOES affect a (auto-dereferenced: same as (*ptr).Field = x)

func f(p Person)  { }  // p is a copy - mutations don't escape
func g(p *Person) { }  // p shares the caller's struct - mutations DO escape
```

---

## ⚖️ Comparing Structs

```go
a == b   // true if every field is equal
a != b

var zero Point
a == zero          // check "is this the zero value?"
zero == Point{}    // also valid

m := map[Point]string{ {0, 0}: "origin" } // comparable structs work as map keys
```

```go
// ❌ Compile error - struct has a slice field
// invalid operation: t1 == t2 (struct containing []string cannot be compared)
```

---

## 🗂️ Nested vs Embedded Structs

```go
// Nested (named field) - "has a"
type Person struct {
    Name    string
    Address Address
}
p.Address.City            // always full path

// Embedded (anonymous field) - fields PROMOTED
type Dog struct {
    Animal
    Breed string
}
d.Name                    // promoted shortcut
d.Animal.Name             // full path, always works too

// Shadowing - outer field wins when names collide
type Cat struct {
    Animal
    Name string
}
c.Name          // outer field
c.Animal.Name   // inner (embedded) field
```

---

## 🧩 Anonymous Structs

```go
server := struct {
    Host string
    Port int
}{Host: "localhost", Port: 8080}

tests := []struct {
    Input, Expected int
}{
    {Input: 2, Expected: 4},
    {Input: 3, Expected: 9},
}
```

---

## 🏷️ Struct Tags

```go
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age,omitempty"`
}
// Syntax only here - read by reflection (reflect.StructField.Tag)
// Used heavily by encoding/json - see Level 19
```

---

## 🏗️ Constructor Convention

```go
func NewBankAccount(owner string, balance float64) *BankAccount {
    if balance < 0 {
        balance = 0
    }
    return &BankAccount{Owner: owner, Balance: balance, Active: true}
}

acct := NewBankAccount("Alice", 250.75)
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Assuming assignment shares | `b := a; b.X = 1` affects `a` | use `ptr := &a` to share |
| Mutating via value param | `func f(p Person) { p.X = 1 }` | `func f(p *Person) { p.X = 1 }` |
| Positional literal drift | `Rect{5, 10}` - which is which? | `Rect{Width: 10, Height: 5}` |
| == on slice/map field | `t1 == t2` won't compile | compare fields manually or `reflect.DeepEqual` |
| Embedded literal by field name | `Dog{Name: "Rex"}` | `Dog{Animal: Animal{Name: "Rex"}}` |
| Missing a shadowed field | assuming `c.Name` reaches the embedded one | use `c.Animal.Name` explicitly |

---

## 🎓 Before Next Level

Can you:
- [ ] Define a struct and build it with keyed, positional, and zero-value forms?
- [ ] Prove struct assignment copies while a pointer shares?
- [ ] Explain exactly when `==` on a struct fails to compile?
- [ ] Build a nested struct and an embedded struct, and explain promotion?
- [ ] Explain what shadowing does when field names collide?
- [ ] Write a NewX constructor?

If YES → You're ready for Level 14!

---

## 📚 Next Level

Level 14: Methods
- Value receivers vs pointer receivers
- Method promotion through embedding
- Method sets - the doorway to interfaces

You've got structs down! 💪
