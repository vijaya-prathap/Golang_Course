# Level 14: Quick Reference Card

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir myapp && cd myapp
go mod init myapp

# Create main.go
cat > main.go << 'EOF'
package main

import "fmt"

type Counter struct {
    Count int
}

func (c *Counter) Increment() {
    c.Count++
}

func (c Counter) Peek() int {
    return c.Count
}

func main() {
    c := Counter{}
    c.Increment()
    c.Increment()
    fmt.Println(c.Peek())
}
EOF

# Run
go run main.go
# prints: 2
```

---

## 📋 Method Declaration

```go
//    ┌ receiver: binds the func to a type
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

rect.Area()      // method call
// receiver name: short (r), consistent across methods, never this/self
// methods only on types declared in the SAME package
```

---

## ⚖️ Value vs Pointer Receivers

```go
func (p Person) Greet() string { }   // value:   p is a COPY
func (p *Person) Birthday()    { }   // pointer: p is the ORIGINAL

// value receiver mutation = silent no-op:
func (c Counter) Increment() { c.Count++ }  // ❌ caller sees nothing

// pointer receiver mutation sticks:
func (c *Counter) Increment() { c.Count++ } // ✅ caller sees it
```

**Choose pointer when:** the method mutates, the struct is large, or ANY other method already uses a pointer receiver (consistency rule: one pointer → all pointer).

**Choose value when:** small value-like type (Celsius, Point) with read-only methods.

---

## 🎁 Auto-Address / Auto-Dereference Sugar

```go
acct.Deposit(50)   // value calling pointer method  → (&acct).Deposit(50)
ptr.Report()       // pointer calling value method  → (*ptr).Report()

// Sugar needs ADDRESSABILITY - these FAIL to compile:
players["alice"].AddPoints(5)   // map element ❌
makeCounter().Increment()       // function return ❌
// error: cannot call pointer method AddPoints on Player

// Fixes:
p := players["alice"]; p.AddPoints(5); players["alice"] = p
roster := map[string]*Player{...}; roster["bob"].AddPoints(5)
```

---

## 📦 Method Sets (matters for Level 15!)

| Method receiver | In `T`'s set | In `*T`'s set |
|-----------------|--------------|----------------|
| `func (t T)` value | ✅ | ✅ |
| `func (t *T)` pointer | ❌ | ✅ |

```go
var i Incrementer = c    // ❌ Counter lacks pointer-recv Increment
var i Incrementer = &c   // ✅ *Counter has it
```

---

## 🧬 Embedding & Promotion

```go
type Dog struct {
    Animal          // embedded: Animal's methods PROMOTED to Dog
    Breed string
}

d.Describe()        // promoted from Animal
func (d Dog) Speak() string { }  // same name = OVERRIDES promoted Speak
d.Animal.Speak()    // reach the shadowed embedded version explicitly
```

---

## 🔗 Method Values & Expressions

```go
greet := p.Greet          // method VALUE: receiver bound now (snapshot for value recv)
greet()

greetFn := Person.Greet   // method EXPRESSION: func(Person) string
greetFn(p)

renameFn := (*Person).Rename  // func(*Person, string)
renameFn(&p, "Allie")
```

---

## 🏷️ Named Non-Struct Types

```go
type Celsius float64
func (c Celsius) ToFahrenheit() Fahrenheit { return Fahrenheit(c*9/5 + 32) }

type Path string
func (p Path) Base() string { /* ... */ }

Celsius(100).ToFahrenheit()   // 212
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Mutating via value receiver | `func (c Counter) Inc() { c.Count++ }` | `func (c *Counter) Inc()` |
| Pointer method on map element | `m["k"].AddPoints(5)` | copy out, store back (or `map[K]*V`) |
| this/self receiver names | `func (this *T) ...` | `func (t *T) ...` |
| Mixed receivers on one type | some `T`, some `*T` | all `*T` if any method needs one |
| Forgetting store-back | `p := m["k"]; p.Fix()` | add `m["k"] = p` |

---

## 🎓 Before Next Level

Can you:
- [ ] Declare value- and pointer-receiver methods from memory?
- [ ] Predict which mutations stick without running the code?
- [ ] State the consistency rule and the reason behind it?
- [ ] Explain why `players["alice"].AddPoints(5)` won't compile?
- [ ] Fill in the T vs *T method-set table?
- [ ] Override a promoted method and reach the embedded version?

If YES → You're ready for Level 15!

---

## 📚 Next Level

Level 15: Interfaces
- Defining interfaces and implicit satisfaction
- Method sets in action
- fmt.Stringer and the standard library's small-interface style

Your types have behavior now! 💪
