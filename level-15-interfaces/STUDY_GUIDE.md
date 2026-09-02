# Level 15: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: Interfaces
```
Day 1:  What an interface is; implicit satisfaction
Day 2:  Polymorphism with Shape; interface values as (type, value) pairs
Day 3:  fmt.Stringer; the empty interface (any)
Day 4:  Type assertions (comma-ok vs panic) and type switches
Day 5:  Method sets and interface satisfaction
Day 6:  The typed-nil trap; interface composition
Day 7:  Interface design: small interfaces, consumer-side contracts
```

### Week 2: Practice & Application
```
Day 1:  Exercises 1-3 (Shape, Stringer, any)
Day 2:  Exercises 4-5 (assertions, type switches)
Day 3:  Exercises 6-7 (method sets, typed nil)
Day 4:  Exercises 8-9 (composition, consumer-side design)
Day 5:  Exercise 10 (notification system)
Day 6:  Bonus challenges
Day 7:  Practice & review
```

---

## 🎯 What an Interface Is

```
STRUCT = data (what a thing IS)        INTERFACE = behavior (what a thing can DO)

type Dog struct {                      type Speaker interface {
    Name string                            Speak() string
    Age  int                           }
}
                                       - no fields
                                       - no method bodies
                                       - no list of implementers
```

---

## 🔌 Implicit Satisfaction

```
        NO "implements" KEYWORD - the methods ARE the declaration

   type Speaker interface {              type Dog struct{ Name string }
       Speak() string          ◄─────    func (d Dog) Speak() string { ... }
   }                        satisfies
                           automatically  type Robot struct{ ID int }
                               ◄─────    func (r Robot) Speak() string { ... }

   Dog and Robot NEVER mention Speaker.
   Having the right methods IS the implementation.

   Why it matters:
   ├─ types satisfy interfaces from packages they never imported
   ├─ consumers define exactly the contract they need
   └─ one type can satisfy many interfaces at once, forever
```

---

## 📦 Interface Values: the (type, value) Pair

Every interface value is a two-slot box. This model explains everything.

```
var s Shape                s = Circle{2}              var p *MyErr; var e error = p

┌───────────────┐          ┌───────────────┐          ┌────────────────────┐
│ type:  nil    │          │ type:  Circle │          │ type:  *MyErr      │
│ value: nil    │          │ value: {2}    │          │ value: nil         │
└───────────────┘          └───────────────┘          └────────────────────┘
  NIL INTERFACE              NORMAL VALUE                TYPED NIL (!)
  s == nil → true            s == nil → false            e == nil → FALSE
  method call → PANIC        method call → Circle's      the classic trap
```

**The rule:** an interface is `== nil` only when BOTH slots are nil.

```
method call s.Area()
    └─ looks at the type slot → finds Circle → runs Circle's Area()
       (this dispatch-on-dynamic-type IS polymorphism)

%T prints the type slot        %v prints the value slot
```

---

## 📊 Method Sets & Interface Satisfaction Table

Which method sets contain which methods (from Level 14):

| Receiver on method | In method set of `T`? | In method set of `*T`? |
|--------------------|----------------------|------------------------|
| `func (t T)` value receiver | ✅ yes | ✅ yes |
| `func (t *T)` pointer receiver | ❌ NO | ✅ yes |

Consequence for interfaces (satisfaction requires ALL interface methods):

| Interface needs | Assign `T` value | Assign `*T` pointer |
|-----------------|------------------|---------------------|
| only value-receiver methods | ✅ works | ✅ works |
| any pointer-receiver method | ❌ **compile error** | ✅ works |

The exact compile error to recognize:

```
Counter does not implement Incrementer (method Increment has pointer receiver)
```

```
Rule of thumb:
    ANY method with a pointer receiver?
        └─ store POINTERS (&T{...}) in interface values
```

---

## 🔍 Type Assertion / Type Switch Decision Tree

```
I have an interface value and need the concrete type back.
│
├─ Expecting exactly ONE possible type?
│  │
│  ├─ Mismatch is normal / must be handled gracefully?
│  │  └─ comma-ok assertion:      v, ok := i.(T)
│  │        ok=false → zero value, NO panic
│  │
│  └─ Mismatch means a bug that SHOULD crash?
│     └─ single-value assertion:  v := i.(T)
│           wrong type → panic: interface conversion: interface {} is string, not int
│
└─ Expecting SEVERAL possible types?
   └─ type switch:
          switch v := i.(type) {
          case int:      // v is an int here
          case string:   // v is a string here
          case nil:      // nil interface
          default:       // anything else (inspect with %T)
          }
```

---

## 🧩 Interface Composition

```
type Reader interface {          type Writer interface {
    Read() string                    Write(data string)
}                                }

            └───────────┬───────────┘
                        ▼
            type ReadWriter interface {
                Reader              ← embedded
                Writer              ← embedded
            }

A type with Read() AND Write() satisfies all three interfaces.
The io package works exactly this way:
    io.Reader + io.Writer = io.ReadWriter
    → one io.Copy(dst, src) copies file→network, buffer→file, anything→anything
```

---

## 🧭 Interface Design Principles

> **"The bigger the interface, the weaker the abstraction."** - Rob Pike

```
1. KEEP INTERFACES SMALL
   1-2 methods is normal.  error, io.Reader, fmt.Stringer: one method each.
   Name single-method interfaces with -er: Reader, Writer, Notifier, Stringer.

2. ACCEPT INTERFACES, RETURN CONCRETE TYPES
   func report(w io.Writer) error     ← parameter: interface (flexible callers)
   func NewBuffer() *bytes.Buffer     ← return: concrete (full access)

3. DEFINE INTERFACES WHERE THEY'RE USED
   ❌ provider exports a big interface and hopes it fits everyone
   ✅ each consumer declares the 1-2 methods IT actually calls
      (works because satisfaction is implicit)

4. DON'T CREATE INTERFACES PREEMPTIVELY
   Start concrete. Add the interface when a second implementation
   or a decoupling need actually appears.
```

---

## 🚨 Common Mistakes

### Mistake 1: Value Assigned Where a Pointer Receiver Is Needed

```go
func (c *Counter) Increment() { c.count++ }

// ❌ WRONG - compile error: method Increment has pointer receiver
var inc Incrementer = Counter{}

// ✅ RIGHT
var inc Incrementer = &Counter{}
```

### Mistake 2: Returning a Concrete Error Pointer

```go
// ❌ WRONG - callers get a typed nil; err != nil is TRUE on success
func validate(s string) *ValidationError { return nil }

// ✅ RIGHT - return the interface type
func validate(s string) error { return nil }
```

### Mistake 3: Unguarded Single-Value Assertion

```go
// ❌ WRONG - panics if i doesn't hold an int
n := i.(int)

// ✅ RIGHT
n, ok := i.(int)
```

### Mistake 4: Using any When a Small Interface Fits

```go
// ❌ WRONG - process must type-switch on the world
func process(item any)

// ✅ RIGHT - name the behavior you need
func process(item Notifier)
```

### Mistake 5: Calling a Method on a nil Interface

```go
var s Shape
// ❌ WRONG - panic: invalid memory address or nil pointer dereference
s.Area()
```

### Mistake 6: Comparing a Typed nil to nil

```go
var p *MyErr
var err error = p
// ❌ WRONG ASSUMPTION - prints false! (type slot holds *MyErr)
fmt.Println(err == nil)
```

---

## 📈 Progression Summary

### Understanding Level 15

Level 15 turns the methods of Level 14 into contracts:

1. **The contract** — an interface is a named set of method signatures
2. **The magic** — satisfaction is implicit; the methods are the declaration
3. **The machinery** — every interface value is a (dynamic type, dynamic value) pair
4. **The tools** — assertions, type switches, `any`, composition
5. **The judgment** — small interfaces, defined at the consumer, concrete returns

### Prerequisites for Level 16

Before moving to Level 16 (Error Handling), you need:

- ✅ Comfortable defining and satisfying interfaces implicitly
- ✅ Can explain the (type, value) pair and both nil cases
- ✅ Can use comma-ok assertions and type switches fluently
- ✅ Can predict satisfaction from method sets (T vs *T)
- ✅ Understand why functions return `error`, never `*MyError`

### Ready for Level 16?

Level 16 is built entirely on one interface you can already read:

```go
type error interface {
    Error() string
}
```

- Creating and wrapping errors (`errors.New`, `fmt.Errorf` with `%w`)
- `errors.Is` and `errors.As` - the latter is type assertion machinery at work
- Custom error types - where this level's typed-nil rule protects you

---

## ✅ Checklist Before Level 16

- [ ] Can define an interface and satisfy it without any `implements` keyword
- [ ] Can write a function that accepts an interface and works with multiple types
- [ ] Can implement fmt.Stringer and know who calls String()
- [ ] Can explain the (dynamic type, dynamic value) pair
- [ ] Can use `v, ok := i.(T)` and `switch v := i.(type)` from memory
- [ ] Can predict the pointer-receiver compile error before the compiler shows it
- [ ] Can explain why an interface holding (*T)(nil) is not == nil
- [ ] Know the three design rules: small, consumer-side, accept-interfaces/return-concrete
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The Contract
An interface describes behavior - a set of method signatures. Structs are nouns; interfaces are job descriptions.

### The Magic
Satisfaction is implicit. Having the methods IS implementing the interface - which decouples packages and lets consumers own their contracts.

### The Machinery
An interface value is a (dynamic type, dynamic value) pair. Dispatch, `%T`, nil checks, and the typed-nil trap all fall out of this one model.

### The Judgment
The bigger the interface, the weaker the abstraction. Keep them small, define them where they're used, accept interfaces, return concrete types.

---

## 📚 Next Level

Level 16: Error Handling
- The error interface - one method, `Error() string`
- Creating, wrapping, and inspecting errors
- errors.Is / errors.As and error chains
- Custom error types done safely (typed-nil aware)

You've mastered Go's most important abstraction! Keep going! 🚀
