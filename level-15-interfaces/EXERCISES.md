# Level 15: Interfaces - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

---

## Exercise 1: Shape Polymorphism

**Objective:** Define an interface, satisfy it with two types, and write one function that works with both

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level15-exercise1
cd ~/projects/level15-exercise1
go mod init level15.example/exercise1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "math"
)

type Shape interface {
    Area() float64
    Perimeter() float64
}

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}

type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

func describe(s Shape) {
    fmt.Printf("%T: area=%.2f, perimeter=%.2f\n", s, s.Area(), s.Perimeter())
}

func main() {
    fmt.Println("=== One Function, Many Shapes ===")
    c := Circle{Radius: 2}
    r := Rectangle{Width: 3, Height: 4}
    describe(c)
    describe(r)

    fmt.Println("\n=== A Slice of Shapes ===")
    shapes := []Shape{
        Circle{Radius: 1},
        Rectangle{Width: 2, Height: 5},
        Circle{Radius: 3},
    }
    totalArea := 0.0
    for _, s := range shapes {
        totalArea += s.Area()
    }
    fmt.Printf("total area of %d shapes: %.2f\n", len(shapes), totalArea)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== One Function, Many Shapes ===
main.Circle: area=12.57, perimeter=12.57
main.Rectangle: area=12.00, perimeter=14.00

=== A Slice of Shapes ===
total area of 3 shapes: 41.42
```

Notice that `Circle` and `Rectangle` never mention `Shape` anywhere - they satisfy it purely by having the right methods. That's implicit satisfaction.

**Learning Objectives:**
- ✅ Define an interface as a set of method signatures
- ✅ Satisfy it implicitly - no implements keyword
- ✅ Write one polymorphic function that handles every Shape
- ✅ Store mixed concrete types in a []Shape and dispatch dynamically

---

## Exercise 2: fmt.Stringer

**Objective:** Implement `String() string` and watch fmt use it automatically

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level15-exercise2
cd ~/projects/level15-exercise2
go mod init level15.example/exercise2
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type Book struct {
    Title  string
    Author string
    Pages  int
}

func (b Book) String() string {
    return fmt.Sprintf("%q by %s (%d pages)", b.Title, b.Author, b.Pages)
}

type Temperature float64

func (t Temperature) String() string {
    return fmt.Sprintf("%.1f°C", float64(t))
}

func main() {
    fmt.Println("=== Without String(), fmt Uses Default Struct Formatting ===")
    type plainBook struct {
        Title  string
        Author string
        Pages  int
    }
    fmt.Println(plainBook{"The Go Programming Language", "Donovan & Kernighan", 380})

    fmt.Println("\n=== With String(), fmt.Println Calls It Automatically ===")
    b := Book{"The Go Programming Language", "Donovan & Kernighan", 380}
    fmt.Println(b)

    fmt.Println("\n=== It Works Anywhere fmt Formats a Value ===")
    fmt.Printf("currently reading: %v\n", b)
    temps := []Temperature{21.5, 23.0, 18}
    fmt.Println(temps)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Without String(), fmt Uses Default Struct Formatting ===
{The Go Programming Language Donovan & Kernighan 380}

=== With String(), fmt.Println Calls It Automatically ===
"The Go Programming Language" by Donovan & Kernighan (380 pages)

=== It Works Anywhere fmt Formats a Value ===
currently reading: "The Go Programming Language" by Donovan & Kernighan (380 pages)
[21.5°C 23.0°C 18.0°C]
```

Your types satisfied `fmt.Stringer` - an interface defined in the standard library, long before your code existed. `fmt` checked for it and called YOUR method. Even the slice printed with your formatting, element by element.

**Learning Objectives:**
- ✅ Implement the fmt.Stringer interface with a String() method
- ✅ See fmt.Println and %v call it automatically
- ✅ Attach a Stringer to a non-struct named type (Temperature)
- ✅ Understand this as implicit satisfaction across package boundaries

---

## Exercise 3: The Empty Interface (any)

**Objective:** Use `any` to hold values of different types and see what it costs

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level15-exercise3
cd ~/projects/level15-exercise3
go mod init level15.example/exercise3
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== any Holds a Value of Any Type ===")
    var x any
    x = 42
    fmt.Printf("value=%v type=%T\n", x, x)
    x = "hello"
    fmt.Printf("value=%v type=%T\n", x, x)
    x = 3.14
    fmt.Printf("value=%v type=%T\n", x, x)
    x = []int{1, 2, 3}
    fmt.Printf("value=%v type=%T\n", x, x)

    fmt.Println("\n=== A Heterogeneous Slice ===")
    mixed := []any{42, "gopher", true, 3.14}
    for i, v := range mixed {
        fmt.Printf("mixed[%d] = %v (%T)\n", i, v, v)
    }

    fmt.Println("\n=== The Cost: No Type Safety ===")
    var n any = 10
    // total := n + 5 // ❌ would NOT compile: any is not an int
    fmt.Println("n holds", n, "- but n + 5 is a compile error,")
    fmt.Println("because the static type of n is any, not int")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== any Holds a Value of Any Type ===
value=42 type=int
value=hello type=string
value=3.14 type=float64
value=[1 2 3] type=[]int

=== A Heterogeneous Slice ===
mixed[0] = 42 (int)
mixed[1] = gopher (string)
mixed[2] = true (bool)
mixed[3] = 3.14 (float64)

=== The Cost: No Type Safety ===
n holds 10 - but n + 5 is a compile error,
because the static type of n is any, not int
```

4. **See the cost for yourself.** Uncomment the `total := n + 5` line and run again - the compiler rejects it, because once a value is inside `any`, its concrete type is invisible to the compiler.

**Learning Objectives:**
- ✅ Understand that any (= interface{}) is satisfied by every type
- ✅ Build a heterogeneous slice with []any
- ✅ Use %T to inspect the dynamic type inside an interface
- ✅ Understand the trade-off: total flexibility, zero compile-time safety

---

## Exercise 4: Type Assertions - comma-ok vs Panic

**Objective:** Extract concrete values from interfaces safely, then trigger the real assertion panic

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level15-exercise4
cd ~/projects/level15-exercise4
go mod init level15.example/exercise4
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    var i any = "hello"

    fmt.Println("=== The comma-ok Form Never Panics ===")
    s, ok := i.(string)
    fmt.Printf("s=%q ok=%v\n", s, ok)

    n, ok := i.(int)
    fmt.Printf("n=%d ok=%v (failed assertion gives the zero value)\n", n, ok)

    fmt.Println("\n=== Guarding Behavior With comma-ok ===")
    values := []any{"go", 42, true, "gopher"}
    for _, v := range values {
        if s, ok := v.(string); ok {
            fmt.Printf("%q is a string of length %d\n", s, len(s))
        } else {
            fmt.Printf("%v (%T) is not a string, skipping\n", v, v)
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
=== The comma-ok Form Never Panics ===
s="hello" ok=true
n=0 ok=false (failed assertion gives the zero value)

=== Guarding Behavior With comma-ok ===
"go" is a string of length 2
42 (int) is not a string, skipping
true (bool) is not a string, skipping
"gopher" is a string of length 6
```

4. **Now trigger the real panic.** Create a second program that uses the single-value form on a wrong type:

```bash
mkdir -p ~/projects/level15-exercise4-panic
cd ~/projects/level15-exercise4-panic
go mod init level15.example/exercise4panic

cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    var i any = "hello"
    n := i.(int) // i holds a string, not an int - this panics
    fmt.Println(n)
}
EOF

go run main.go
```

**Expected Output (a real panic):**

```
panic: interface conversion: interface {} is string, not int

goroutine 1 [running]:
main.main()
        /path/to/main.go:7 +0x34
exit status 2
```

The file path and hex offset will differ on your machine, but the panic message - `interface conversion: interface {} is string, not int` - is the exact, fixed wording from the Go runtime. Note that it names both the actual dynamic type (`string`) and the type you wrongly asserted (`int`).

**Learning Objectives:**
- ✅ Use the comma-ok form to assert safely
- ✅ See that a failed comma-ok assertion yields the zero value and false
- ✅ Trigger and read the exact interface-conversion panic
- ✅ Know the rule: comma-ok in normal flow, single-value only when failure should crash

---

## Exercise 5: Type Switches

**Objective:** Branch on an interface's dynamic type with `switch v := i.(type)`

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level15-exercise5
cd ~/projects/level15-exercise5
go mod init level15.example/exercise5
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func describe(i any) string {
    switch v := i.(type) {
    case nil:
        return "nil value"
    case int:
        return fmt.Sprintf("int %d (doubled: %d)", v, v*2)
    case string:
        return fmt.Sprintf("string %q (%d bytes)", v, len(v))
    case bool:
        return fmt.Sprintf("bool %v", v)
    case []int:
        sum := 0
        for _, n := range v {
            sum += n
        }
        return fmt.Sprintf("[]int with %d elements (sum: %d)", len(v), sum)
    default:
        return fmt.Sprintf("unhandled type %T", v)
    }
}

func main() {
    fmt.Println("=== Type Switch in Action ===")
    values := []any{42, "gopher", true, []int{1, 2, 3}, 3.14, nil}
    for _, v := range values {
        fmt.Println(describe(v))
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
=== Type Switch in Action ===
int 42 (doubled: 84)
string "gopher" (6 bytes)
bool true
[]int with 3 elements (sum: 6)
unhandled type float64
nil value
```

Inside each case, `v` already has that case's concrete type: the `int` case computes `v*2`, the `string` case calls `len(v)`, and the `[]int` case ranges over it - no separate assertions needed.

**Learning Objectives:**
- ✅ Use switch v := i.(type) to branch on the dynamic type
- ✅ Use v with its concrete type inside each case
- ✅ Handle nil and default cases explicitly
- ✅ Prefer a type switch over chained comma-ok assertions

---

## Exercise 6: Method Sets - The Pointer-Receiver Compile Error

**Objective:** See exactly when a value satisfies an interface and capture the classic compile error when it doesn't

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level15-exercise6
cd ~/projects/level15-exercise6
go mod init level15.example/exercise6
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type Incrementer interface {
    Increment()
    Value() int
}

type Counter struct {
    count int
}

func (c *Counter) Increment() {
    c.count++
}

func (c Counter) Value() int {
    return c.count
}

func main() {
    fmt.Println("=== A *Counter Satisfies Incrementer ===")
    c := Counter{}
    var inc Incrementer = &c // pointer: method set has BOTH receiver kinds
    inc.Increment()
    inc.Increment()
    inc.Increment()
    fmt.Println("value after 3 increments:", inc.Value())
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== A *Counter Satisfies Incrementer ===
value after 3 increments: 3
```

4. **Now trigger the compile error for real.** Create a second program that assigns a `Counter` *value* instead of a pointer:

```bash
mkdir -p ~/projects/level15-exercise6bad
cd ~/projects/level15-exercise6bad
go mod init level15.example/exercise6bad

cat > main.go << 'EOF'
package main

import "fmt"

type Incrementer interface {
    Increment()
    Value() int
}

type Counter struct {
    count int
}

func (c *Counter) Increment() {
    c.count++
}

func (c Counter) Value() int {
    return c.count
}

func main() {
    c := Counter{}
    var inc Incrementer = c // ❌ value: Increment is NOT in Counter's method set
    inc.Increment()
    fmt.Println(inc.Value())
}
EOF

go build ./...
```

**Expected Output (compile error, captured verbatim from `go build`):**

```
# level15.example/exercise6bad
./main.go:24:27: cannot use c (variable of struct type Counter) as Incrementer value in variable declaration: Counter does not implement Incrementer (method Increment has pointer receiver)
```

The build fails with exit status 1. Read the parenthetical at the end - `(method Increment has pointer receiver)` - the compiler names the exact method and receiver kind that broke satisfaction. The fix is one character: assign `&c` instead of `c`.

**Learning Objectives:**
- ✅ Confirm *T satisfies interfaces via both value- and pointer-receiver methods
- ✅ Trigger the real "does not implement" compile error with a T value
- ✅ Read the compiler's hint naming the pointer-receiver method
- ✅ Apply the rule: any pointer receiver anywhere → store pointers in interfaces

---

## Exercise 7: The Typed-nil Interface Trap

**Objective:** Demonstrate that an interface holding a nil pointer is NOT nil, and fix the classic error-return bug

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level15-exercise7
cd ~/projects/level15-exercise7
go mod init level15.example/exercise7
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type ValidationError struct {
    Field string
}

func (e *ValidationError) Error() string {
    return "invalid field: " + e.Field
}

// ❌ BUGGY: declares the concrete pointer type as its return type
func validateBad(input string) *ValidationError {
    if input == "" {
        return &ValidationError{Field: "input"}
    }
    return nil // a nil *ValidationError - NOT a nil interface!
}

// ✅ FIXED: declares the error interface as its return type
func validateGood(input string) error {
    if input == "" {
        return &ValidationError{Field: "input"}
    }
    return nil // a true nil interface
}

func main() {
    fmt.Println("=== The Trap: a Typed nil Inside an Interface ===")
    var err error = validateBad("valid input") // validation passed...
    fmt.Println("err == nil?", err == nil)
    fmt.Printf("dynamic type: %T, dynamic value: %v\n", err, err)
    if err != nil {
        fmt.Println("BUG: this branch runs even though validation passed!")
    }

    fmt.Println("\n=== The Fix: Return the error Interface Type ===")
    err = validateGood("valid input")
    fmt.Println("err == nil?", err == nil)
    if err == nil {
        fmt.Println("correct: no error reported")
    }

    fmt.Println("\n=== Real Failures Still Work Normally ===")
    err = validateGood("")
    fmt.Println("err == nil?", err == nil)
    fmt.Println("message:", err)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== The Trap: a Typed nil Inside an Interface ===
err == nil? false
dynamic type: *main.ValidationError, dynamic value: <nil>
BUG: this branch runs even though validation passed!

=== The Fix: Return the error Interface Type ===
err == nil? true
correct: no error reported

=== Real Failures Still Work Normally ===
err == nil? false
message: invalid field: input
```

Study the second line: the dynamic *value* is `<nil>`, but the dynamic *type* is `*main.ValidationError` - and an interface is `== nil` only when BOTH halves are nil. That is the entire trap in one line of output.

**Learning Objectives:**
- ✅ Produce a typed nil: interface (type=*T, value=nil)
- ✅ Verify with real output that err != nil is true for a typed nil
- ✅ Explain the behavior with the (dynamic type, dynamic value) pair model
- ✅ Apply the fix: return the error interface type, and return literal nil

---

## Exercise 8: Interface Composition - Building a ReadWriter

**Objective:** Embed interfaces in interfaces, modeled on the io package

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level15-exercise8
cd ~/projects/level15-exercise8
go mod init level15.example/exercise8
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type Reader interface {
    Read() string
}

type Writer interface {
    Write(data string)
}

// Composition: ReadWriter is Reader AND Writer combined
type ReadWriter interface {
    Reader
    Writer
}

type MemoryBuffer struct {
    data string
}

func (m *MemoryBuffer) Read() string {
    return m.data
}

func (m *MemoryBuffer) Write(data string) {
    m.data += data
}

func copyData(dst Writer, src Reader) {
    dst.Write(src.Read())
}

func main() {
    fmt.Println("=== One Type, Three Interfaces ===")
    buf := &MemoryBuffer{}

    var w Writer = buf
    var r Reader = buf
    var rw ReadWriter = buf

    w.Write("hello, ")
    rw.Write("interfaces!")
    fmt.Println("read back:", r.Read())

    fmt.Println("\n=== A ReadWriter Works Wherever a Reader or Writer Is Needed ===")
    src := &MemoryBuffer{data: "copy me"}
    dst := &MemoryBuffer{}
    copyData(dst, src)
    fmt.Println("dst now holds:", dst.Read())
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== One Type, Three Interfaces ===
read back: hello, interfaces!

=== A ReadWriter Works Wherever a Reader or Writer Is Needed ===
dst now holds: copy me
```

This mirrors the real `io` package exactly: `io.ReadWriter` is nothing more than `io.Reader` and `io.Writer` embedded together, and `io.Copy` is a grown-up `copyData`.

**Learning Objectives:**
- ✅ Compose interfaces by embedding (Reader + Writer = ReadWriter)
- ✅ Satisfy the composed interface by satisfying all embedded methods
- ✅ Pass one concrete type as three different interface views
- ✅ Recognize the io.Reader/io.Writer/io.ReadWriter pattern

---

## Exercise 9: Define the Interface at the Consumer

**Objective:** Refactor toward Go style - the consumer declares the contract it needs; providers ship concrete types

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level15-exercise9
cd ~/projects/level15-exercise9
go mod init level15.example/exercise9
```

2. Create `main.go`. Imagine the "provider side" and "consumer side" living in separate packages - the comments mark the boundary:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

// --- provider side: concrete types only, NO interfaces defined here ---

type Database struct{}

func (Database) FetchUsers() []string {
    return []string{"alice", "bob", "carol"}
}

type CSVFile struct{}

func (CSVFile) FetchUsers() []string {
    return []string{"dave", "erin"}
}

// --- consumer side: the CONSUMER defines the interface it needs ---

type UserSource interface {
    FetchUsers() []string
}

func printReport(src UserSource) {
    users := src.FetchUsers()
    fmt.Printf("report: %d user(s)\n", len(users))
    for _, u := range users {
        fmt.Println(" -", u)
    }
}

func main() {
    fmt.Println("=== The Report Works With Any UserSource ===")
    printReport(Database{})
    fmt.Println()
    printReport(CSVFile{})
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== The Report Works With Any UserSource ===
report: 3 user(s)
 - alice
 - bob
 - carol

report: 2 user(s)
 - dave
 - erin
```

4. **Think through the decoupling.** `Database` and `CSVFile` never mention `UserSource` - in a real project they could live in a package that has never imported the report code. The consumer declared exactly the one method it calls, so ANY provider (including a test fake you write in 3 lines) plugs in. This only works because satisfaction is implicit.

**Learning Objectives:**
- ✅ Place the interface with the consumer, not the provider
- ✅ Keep the interface minimal - only the methods the consumer calls
- ✅ Swap providers without either provider knowing the interface exists
- ✅ Recognize why this makes testing trivial (fakes satisfy automatically)

---

## Exercise 10: Comprehensive Practice - Notification System

**Objective:** Combine polymorphism, fan-out over a []Notifier, open extension, and a type switch in one program

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level15-exercise10
cd ~/projects/level15-exercise10
go mod init level15.example/exercise10
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "strings"
)

type Notifier interface {
    Notify(message string)
}

type EmailNotifier struct {
    Address string
}

func (e EmailNotifier) Notify(message string) {
    fmt.Printf("[email -> %s] %s\n", e.Address, message)
}

type SMSNotifier struct {
    Number string
}

func (s SMSNotifier) Notify(message string) {
    fmt.Printf("[sms   -> %s] %s\n", s.Number, message)
}

type SlackNotifier struct {
    Channel string
}

func (s SlackNotifier) Notify(message string) {
    fmt.Printf("[slack -> %s] %s\n", s.Channel, strings.ToUpper(message))
}

type PagerNotifier struct {
    Team string
}

func (p PagerNotifier) Notify(message string) {
    fmt.Printf("[pager -> %s] URGENT: %s\n", p.Team, message)
}

func broadcast(message string, targets []Notifier) {
    fmt.Printf("broadcasting to %d notifier(s):\n", len(targets))
    for _, t := range targets {
        t.Notify(message)
    }
}

func main() {
    fmt.Println("=== Building the Notifier List ===")
    targets := []Notifier{
        EmailNotifier{Address: "team@example.com"},
        SMSNotifier{Number: "+1-555-0100"},
        SlackNotifier{Channel: "#deploys"},
    }
    for _, t := range targets {
        fmt.Printf("registered %T\n", t)
    }

    fmt.Println("\n=== Fan-Out: One Call, Many Deliveries ===")
    broadcast("deploy finished", targets)

    fmt.Println("\n=== Adding a Notifier Changes Nothing in broadcast ===")
    targets = append(targets, PagerNotifier{Team: "oncall"})
    broadcast("disk almost full", targets)

    fmt.Println("\n=== Inspecting Dynamic Types With a Type Switch ===")
    for _, t := range targets {
        switch n := t.(type) {
        case EmailNotifier:
            fmt.Println("email goes to", n.Address)
        case SMSNotifier:
            fmt.Println("sms goes to", n.Number)
        default:
            fmt.Printf("other notifier: %T\n", n)
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
=== Building the Notifier List ===
registered main.EmailNotifier
registered main.SMSNotifier
registered main.SlackNotifier

=== Fan-Out: One Call, Many Deliveries ===
broadcasting to 3 notifier(s):
[email -> team@example.com] deploy finished
[sms   -> +1-555-0100] deploy finished
[slack -> #deploys] DEPLOY FINISHED

=== Adding a Notifier Changes Nothing in broadcast ===
broadcasting to 4 notifier(s):
[email -> team@example.com] disk almost full
[sms   -> +1-555-0100] disk almost full
[slack -> #deploys] DISK ALMOST FULL
[pager -> oncall] URGENT: disk almost full

=== Inspecting Dynamic Types With a Type Switch ===
email goes to team@example.com
sms goes to +1-555-0100
other notifier: main.SlackNotifier
other notifier: main.PagerNotifier
```

`broadcast` is written once against the one-method `Notifier` contract. Each concrete type delivers its own way (Slack even shouts), and adding `PagerNotifier` required zero changes to `broadcast` - polymorphism keeping code open for extension.

**Learning Objectives:**
- ✅ Design a small, single-method interface for a real scenario
- ✅ Fan a message out over a []Notifier polymorphically
- ✅ Extend the system with a new type without touching existing functions
- ✅ Combine interfaces, dynamic dispatch, and a type switch in one program

---

## Bonus Challenges

### Challenge 1: sort.Interface

The standard library's `sort.Sort` works on ANY type that satisfies `sort.Interface` - three methods: `Len() int`, `Less(i, j int) bool`, and `Swap(i, j int)`. Define a `ByLength` type (based on `[]string`) that sorts strings by length, and sort a word list with it.

```bash
mkdir -p ~/projects/level15-bonus1
cd ~/projects/level15-bonus1
go mod init level15.example/bonus1
```

**Hints:**
- Declare `type ByLength []string`, then implement the three methods on it
- `Less(i, j int) bool` should return `len(b[i]) < len(b[j])`
- Call `sort.Sort(ByLength(words))` and print the slice before and after
- Bonus insight: this is consumer-side interface design - `sort` defined the contract it needs, and your type satisfied it implicitly

### Challenge 2: Plugin Registry

Build a tiny command dispatcher: a `Handler` interface with a `Handle(args string) string` method, and a registry `map[string]Handler`. Register a few handlers (e.g. `"greet"`, `"upper"`, `"len"`), then dispatch commands by name.

```bash
mkdir -p ~/projects/level15-bonus2
cd ~/projects/level15-bonus2
go mod init level15.example/bonus2
```

**Hints:**
- Each handler is its own struct type with a `Handle` method - state lives in the struct if needed
- Look up with comma-ok: `h, ok := registry[name]` - print a helpful message for unknown commands
- Loop over a slice of test commands like `{"greet Ann", "upper go", "unknown x"}` and dispatch each
- This is how real plugin systems, routers, and command-line tools are structured in Go

### Challenge 3: One Stack Interface, Two Implementations

Define a `Stack` interface with `Push(v int)`, `Pop() (int, bool)`, and `Len() int`. Implement it twice: once backed by a slice, once backed by a linked list (nodes with `value` and `next *node`). Write ONE test function `exercise(s Stack)` that pushes, pops, and prints - and run it against both implementations.

```bash
mkdir -p ~/projects/level15-bonus3
cd ~/projects/level15-bonus3
go mod init level15.example/bonus3
```

**Hints:**
- Both implementations need pointer receivers (they mutate) - so pass `&SliceStack{}` and `&LinkedStack{}` to `exercise`
- Pop's comma-ok style return handles the empty-stack case without panicking
- Slice version: append to push, reslice (`s.items[:len(s.items)-1]`) to pop
- The payoff: `exercise` cannot tell the implementations apart - that's the abstraction working

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Define interfaces and satisfy them implicitly - no implements keyword
✅ Write polymorphic functions that accept any satisfying type
✅ Implement fmt.Stringer and let fmt call your formatting code
✅ Use any for heterogeneous data - and articulate what it costs
✅ Assert types safely with comma-ok, and recognize the interface-conversion panic
✅ Branch on dynamic types with type switches
✅ Predict interface satisfaction from method sets, and fix the pointer-receiver compile error
✅ Explain and avoid the typed-nil interface trap
✅ Compose interfaces by embedding, io-style
✅ Design small, consumer-side interfaces like an experienced Go programmer

---

## Next Level

Level 16: Error Handling
- The error interface - one method, built on everything from this level
- Creating errors with errors.New and fmt.Errorf
- Wrapping, unwrapping, errors.Is and errors.As
- Custom error types - where typed-nil awareness pays off

Great work! You're mastering Go! 🚀
