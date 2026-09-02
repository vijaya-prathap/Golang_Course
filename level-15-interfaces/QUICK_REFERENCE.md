# Level 15: Quick Reference Card

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir myapp && cd myapp
go mod init myapp

# Create main.go
cat > main.go << 'EOF'
package main

import "fmt"

type Shape interface {
    Area() float64
}

type Circle struct{ Radius float64 }

func (c Circle) Area() float64 { return 3.14159 * c.Radius * c.Radius }

func describe(s Shape) {
    fmt.Printf("%T area=%.2f\n", s, s.Area())
}

func main() {
    describe(Circle{Radius: 2})
}
EOF

# Run
go run main.go
```

---

## 📋 Defining and Satisfying an Interface

```go
// Define: a set of method signatures - behavior, not data
type Shape interface {
    Area() float64
    Perimeter() float64
}

// Satisfy: IMPLICITLY - no "implements" keyword anywhere
type Circle struct{ Radius float64 }

func (c Circle) Area() float64      { return 3.14159 * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * 3.14159 * c.Radius }
// Circle now satisfies Shape automatically
```

---

## 📦 The (type, value) Pair

```go
var s Shape              // type=nil  value=nil   s == nil → true
s = Circle{Radius: 2}    // type=Circle value={2} s == nil → false

fmt.Printf("%T", s)      // prints the DYNAMIC TYPE
```

---

## 🗣️ fmt.Stringer

```go
type Stringer interface {
    String() string
}

func (b Book) String() string {
    return fmt.Sprintf("%q", b.Title)
}

fmt.Println(b) // fmt sees Stringer, calls String() automatically
```

---

## 📥 The Empty Interface: any

```go
var x any = 42        // holds ANY type
x = "hello"            // now holds a string

mixed := []any{42, "gopher", true} // heterogeneous slice

// Cost: no compile-time type safety once inside any
```

---

## 🔎 Type Assertions

```go
v, ok := i.(T)   // ✅ safe - comma-ok - ok=false + zero value on mismatch
v := i.(T)       // ⚠️  panics on mismatch - only when failure = bug

// Real panic text (verified):
// panic: interface conversion: interface {} is string, not int
```

---

## 🔀 Type Switches

```go
switch v := i.(type) {
case int:
    fmt.Println("int:", v)
case string:
    fmt.Println("string:", v)
case nil:
    fmt.Println("nil value")
default:
    fmt.Printf("other: %T\n", v)
}
```

---

## 📊 Method Sets and Satisfaction

| Method receiver | In method set of `T` | In method set of `*T` |
|------------------|----------------------|-------------------------|
| value `(t T)` | ✅ | ✅ |
| pointer `(t *T)` | ❌ | ✅ |

```go
func (c *Counter) Increment() { c.count++ } // pointer receiver

var inc Incrementer = &Counter{} // ✅ works
var inc Incrementer = Counter{}  // ❌ compile error:
// Counter does not implement Incrementer (method Increment has pointer receiver)
```

---

## ⚠️ The Typed-nil Trap

```go
func validate(s string) *ValidationError { // ❌ concrete pointer return type
    if s == "" { return &ValidationError{} }
    return nil // typed nil! wrapped in error, err != nil is TRUE
}

func validate(s string) error { // ✅ interface return type
    if s == "" { return &ValidationError{} }
    return nil // true nil interface
}
```

---

## 🧩 Interface Composition

```go
type Reader interface{ Read() string }
type Writer interface{ Write(data string) }

type ReadWriter interface {
    Reader
    Writer
}
// satisfied by any type with BOTH Read() and Write()
```

---

## 🧭 Design Rules

```go
// 1. Keep interfaces small (1-2 methods)
// 2. Accept interfaces, return concrete types
func report(w io.Writer) error
func NewBuffer() *bytes.Buffer

// 3. Define interfaces where they're USED (consumer-side)
type UserSource interface { FetchUsers() []string }
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Value where pointer receiver needed | `var i I = T{}` | `var i I = &T{}` |
| Concrete error return type | `func f() *MyErr` | `func f() error` |
| Unguarded assertion | `v := i.(T)` | `v, ok := i.(T)` |
| Reaching for any | `func f(x any)` | `func f(x Notifier)` |
| Nil interface method call | `var s Shape; s.Area()` | assign a value first |
| Typed nil == nil | assuming `err == nil` | check dynamic type with `%T` |

---

## 🎓 Before Next Level

Can you:
- [ ] Define an interface and satisfy it with zero "implements" syntax?
- [ ] Explain an interface value as a (dynamic type, dynamic value) pair?
- [ ] Implement fmt.Stringer and predict when fmt calls it?
- [ ] Use comma-ok assertions and type switches confidently?
- [ ] Predict the pointer-receiver compile error before seeing it?
- [ ] Explain why a typed nil inside an interface is not == nil?

If YES → You're ready for Level 16!

---

## 📚 Next Level

Level 16: Error Handling
- The error interface (one method: Error() string)
- errors.New, fmt.Errorf, %w wrapping
- errors.Is and errors.As

You've mastered Go's most important abstraction! 💪
