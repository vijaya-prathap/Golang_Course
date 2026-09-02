# Level 14: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: Methods
```
Day 1:  What a method is; declaring methods on structs
Day 2:  Value receivers and copy semantics
Day 3:  Pointer receivers and mutation
Day 4:  Choosing receivers + the consistency rule
Day 5:  Auto-address sugar and addressability
Day 6:  Method sets; methods on named non-struct types
Day 7:  Method promotion, method values & expressions
```

### Week 2: Practice & Application
```
Day 1:  Exercises 1-3 (declaring, value vs pointer receivers)
Day 2:  Exercises 4-5 (sugar, non-addressable compile error)
Day 3:  Exercises 6-7 (consistency refactor, named types)
Day 4:  Exercises 8-10 (promotion, method values, ShoppingCart)
Day 5:  Bonus challenges
Day 6-7: Practice & review
```

---

## 🎯 Method Anatomy

```
func (r Rectangle) Area() float64 {
│    │  │          │      │
│    │  │          │      └─ return type
│    │  │          └─ method name
│    │  └─ receiver TYPE (the type this method belongs to)
│    └─ receiver NAME (short: first letter of type; never this/self)
└─ still just a func!

Call syntax:
    rect.Area()          method: receiver first, reads naturally
    area(rect)           plain function: same power, different shape
```

---

## ⚖️ Value vs Pointer Receivers

| | Value receiver `(p Person)` | Pointer receiver `(p *Person)` |
|---|---|---|
| Receives | a **COPY** of the caller's value | a **pointer to the original** |
| Mutations | lost when the method returns | stick - caller sees them |
| Copy cost per call | whole struct | one pointer (8 bytes) |
| Good for | small, immutable, value-like types | mutation, large structs |
| In method set of `T` | ✅ yes | ❌ no |
| In method set of `*T` | ✅ yes | ✅ yes |

### Verified Behavior

```go
func (c Counter) Increment()  { c.Count++ }  // value receiver
func (p *Player) AddPoints(n int) { p.Score += n }  // pointer receiver

c := Counter{}
c.Increment()
c.Increment()
// c.Count is STILL 0 - each call incremented a discarded copy

b := Player{Score: 0}
b.AddPoints(10)
b.AddPoints(10)
// b.Score is 20 - both calls mutated the original
```

### Decision Tree

```
Does ANY method of this type mutate the receiver?
├─ YES → ALL methods get pointer receivers (consistency rule)
└─ NO
   ├─ Is the struct large (many/heavy fields)?
   │  └─ YES → pointer receivers (avoid copying)
   └─ Small, value-like type (Celsius, Point, RGB)?
      └─ YES → value receivers (copy semantics are a feature)
```

---

## 📦 Method Sets (the Bridge to Level 15)

```
type Counter struct{ n int }

func (c Counter)  Peek() int   { return c.n }   // value receiver
func (c *Counter) Increment()  { c.n++ }        // pointer receiver
```

| Method | In method set of `Counter`? | In method set of `*Counter`? |
|--------|------------------------------|-------------------------------|
| `Peek` (value recv) | ✅ | ✅ |
| `Increment` (pointer recv) | ❌ | ✅ |

```
Rule of thumb:  *T ⊇ T     (the pointer's method set is always ⩾ the value's)
```

**Why care?** Interfaces are satisfied by METHOD SETS, and there's no auto-address sugar there:

```go
var c Counter
var i Incrementer = c   // ❌ compile error:
// Counter does not implement Incrementer (method Increment has pointer receiver)

var i Incrementer = &c  // ✅ *Counter's method set includes Increment
```

---

## 🧲 Addressability: When the Sugar Works

```
v.PointerMethod()  is sugar for  (&v).PointerMethod()   — needs &v to exist!
p.ValueMethod()    is sugar for  (*p).ValueMethod()     — always fine

ADDRESSABLE (sugar works ✅)          NOT ADDRESSABLE (compile error ❌)
├─ local variables:   acct           ├─ map elements:      players["alice"]
├─ struct fields:     order.Buyer    ├─ function returns:  makeCounter()
├─ slice elements:    team[0]        └─ most temporaries/literals
└─ dereferenced ptrs: *ptr
```

The real error, verbatim:

```
./main.go:19:22: cannot call pointer method AddPoints on Player
```

Fixes: copy out → modify → store back, or store pointers (`map[string]*Player`).

---

## 🧬 Method Promotion (Embedding)

```
type Animal struct{ Name string }        type Dog struct {
func (a Animal) Describe() string            Animal   ← embedded (no field name)
func (a Animal) Speak() string               Breed string
                                         }
                                         func (d Dog) Speak() string  ← override

            d := Dog{...}
            │
            ├─ d.Describe()      → PROMOTED from Animal      "Rex is an animal"
            ├─ d.Speak()         → Dog's own version WINS    "Rex says Woof!"
            └─ d.Animal.Speak()  → shadowed version, reached "Rex makes a sound"
                                   explicitly through the field
```

- Promotion = composition, not inheritance: `d.Describe()` is short for `d.Animal.Describe()`
- A same-name outer method **shadows** the promoted one (no virtual dispatch)
- Promoted methods join the outer type's method set → embedding can satisfy interfaces (Level 15)

---

## 🔗 Method Values vs Method Expressions

| Form | Example | Type | Receiver |
|------|---------|------|----------|
| Method **value** | `f := p.Greet` | `func() string` | bound NOW (value recv: a snapshot copy) |
| Method **expression** | `g := Person.Greet` | `func(Person) string` | passed as 1st argument |
| Pointer method expression | `h := (*Person).Rename` | `func(*Person, string)` | passed as 1st argument |

```go
greet := p.Greet     // snapshots p (value receiver)
p.Name = "Alicia"
greet()              // still greets the OLD name!
```

---

## 🚨 Common Mistakes

### Mistake 1: Mutating Through a Value Receiver

```go
// ❌ WRONG - compiles, silently does nothing
func (c Counter) Increment() { c.Count++ }

// ✅ RIGHT
func (c *Counter) Increment() { c.Count++ }
```

### Mistake 2: Pointer Method on a Map Element

```go
// ❌ WRONG - map values are not addressable
players["alice"].AddPoints(5)
// cannot call pointer method AddPoints on Player

// ✅ RIGHT
p := players["alice"]
p.AddPoints(5)
players["alice"] = p
```

### Mistake 3: this / self Receiver Names

```go
// ❌ Not Go style
func (this *Account) Deposit(n float64) { }

// ✅ Go style
func (a *Account) Deposit(n float64) { }
```

### Mistake 4: Mixed Receivers on One Type

```go
// ❌ WRONG - Withdraw mutates a copy; deposit works, withdraw doesn't
func (b *BankAccount) Deposit(n float64) { b.Balance += n }
func (b BankAccount) Withdraw(n float64) { b.Balance -= n }

// ✅ RIGHT - all pointer receivers
```

### Mistake 5: Forgetting the Store-Back

```go
// ❌ WRONG - the map still holds the old value
p := players["alice"]
p.AddPoints(5)

// ✅ RIGHT
players["alice"] = p
```

---

## 📈 Progression Summary

### Understanding Level 14

Level 14 attaches the behavior (Level 11 functions) to the data (Level 13 structs):

1. **Receivers** — a method is a function bound to a type
2. **Value vs pointer** — copy semantics vs mutation, plus the consistency rule
3. **Addressability** — where the auto-& sugar works and where it can't
4. **Method sets** — T gets value-receiver methods; *T gets both
5. **Promotion** — embedded types donate their methods to the outer type

### Prerequisites for Level 15

Before moving to Level 15 (Interfaces), you need:

- ✅ Comfortable declaring value- and pointer-receiver methods
- ✅ Can predict whether a mutation sticks, before running the code
- ✅ Know the consistency rule and why it exists
- ✅ Can state the method-set rule for T vs *T from memory
- ✅ Understand method promotion through embedding

### Ready for Level 15?

Level 15 is where methods pay off. An interface is just a named list of method signatures - and any type whose METHOD SET covers the list satisfies it automatically:
- Implicit satisfaction (no "implements" keyword)
- Your ShoppingCart's String() → fmt.Stringer
- Why `var i Incrementer = c` failed but `= &c` worked

---

## ✅ Checklist Before Level 15

- [ ] Can declare a method with a value receiver and a pointer receiver
- [ ] Can prove (with a program) that value receivers copy
- [ ] Know when to choose pointer receivers: mutation, size, consistency
- [ ] Can explain why `players["alice"].AddPoints(5)` fails to compile
- [ ] Can fill in the method-set table for T and *T
- [ ] Can override a promoted method and still reach the embedded version
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The Receiver
A method is a function with a receiver — the receiver is a real parameter, named short (`r`), never `this`/`self`.

### The Copy Rule
Value receivers copy, pointer receivers point — a mutation through a value receiver is a silent no-op.

### The Consistency Rule
If one method needs a pointer receiver, give ALL of the type's methods pointer receivers.

### The Method Sets
`T` owns the value-receiver methods; `*T` owns them all — this single rule decides interface satisfaction in Level 15.

---

## 📚 Next Level

Level 15: Interfaces
- Defining interfaces and implicit satisfaction
- Method sets in action: which types satisfy which interfaces
- The empty interface and type assertions
- fmt.Stringer and the standard library's small-interface style

Your types have behavior — next they get contracts! 🚀
