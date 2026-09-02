# Level 12: Pointers - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

---

## Exercise 1: Your First Pointer - & and *

**Objective:** Take an address, store it in a pointer, and read/write through it

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level12-exercise1
cd ~/projects/level12-exercise1
go mod init level12.example/exercise1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    x := 42

    fmt.Println("=== A Variable and Its Address ===")
    fmt.Println("Value of x:", x)
    fmt.Printf("Address of x: %p\n", &x)

    fmt.Println("\n=== A Pointer Stores That Address ===")
    p := &x
    fmt.Printf("p (the pointer) : %p\n", p)
    fmt.Println("*p (dereference):", *p)

    fmt.Println("\n=== Modifying x Through the Pointer ===")
    *p = 100
    fmt.Println("*p =", *p)
    fmt.Println("x  =", x)

    fmt.Println("\n=== Reading x Through the Pointer After a Direct Change ===")
    x = 7
    fmt.Println("x  =", x)
    fmt.Println("*p =", *p)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== A Variable and Its Address ===
Value of x: 42
Address of x: 0x735196f9a020

=== A Pointer Stores That Address ===
p (the pointer) : 0x735196f9a020
*p (dereference): 42

=== Modifying x Through the Pointer ===
*p = 100
x  = 100

=== Reading x Through the Pointer After a Direct Change ===
x  = 7
*p = 7
```

*Note: the `0x...` addresses WILL be different on your machine and on every run - that's expected. Every other line should match exactly.*

**Learning Objectives:**
- ✅ Take a variable's address with &
- ✅ Dereference a pointer with * to read and write
- ✅ See that x and *p are two names for the same memory

---

## Exercise 2: Pointer Types and the nil Zero Value

**Objective:** Explore *int and *string types, and what an unassigned pointer holds

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level12-exercise2
cd ~/projects/level12-exercise2
go mod init level12.example/exercise2
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== Pointer Types ===")
    age := 30
    name := "Alice"

    var agePtr *int = &age
    var namePtr *string = &name

    fmt.Printf("agePtr  has type %T and points at the value %d\n", agePtr, *agePtr)
    fmt.Printf("namePtr has type %T and points at the value %s\n", namePtr, *namePtr)

    fmt.Println("\n=== The Zero Value of a Pointer Is nil ===")
    var unset *int
    fmt.Println("unset:", unset)
    fmt.Println("unset == nil:", unset == nil)

    fmt.Println("\n=== Assigning Makes It Non-nil ===")
    unset = &age
    fmt.Println("unset == nil:", unset == nil)
    fmt.Println("*unset:", *unset)

    fmt.Println("\n=== Checking Before Dereferencing ===")
    var maybe *string
    if maybe != nil {
        fmt.Println("value:", *maybe)
    } else {
        fmt.Println("maybe is nil - not safe to dereference!")
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
=== Pointer Types ===
agePtr  has type *int and points at the value 30
namePtr has type *string and points at the value Alice

=== The Zero Value of a Pointer Is nil ===
unset: <nil>
unset == nil: true

=== Assigning Makes It Non-nil ===
unset == nil: false
*unset: 30

=== Checking Before Dereferencing ===
maybe is nil - not safe to dereference!
```

**Learning Objectives:**
- ✅ Declare typed pointers (*int, *string) and print their types with %T
- ✅ See that an unassigned pointer is nil
- ✅ Guard a dereference with an if p != nil check

---

## Exercise 3: The nil Pointer Panic - and Recovering From It

**Objective:** Trigger the famous nil-dereference panic safely, catch it with recover()

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level12-exercise3
cd ~/projects/level12-exercise3
go mod init level12.example/exercise3
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func safeDereference(p *int) (value int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("recovered from: %v", r)
        }
    }()
    return *p, nil
}

func main() {
    fmt.Println("=== Dereferencing a Valid Pointer ===")
    x := 42
    if v, err := safeDereference(&x); err == nil {
        fmt.Println("Got value:", v)
    }

    fmt.Println("\n=== Dereferencing a nil Pointer (recovered) ===")
    var p *int
    if _, err := safeDereference(p); err != nil {
        fmt.Println("Error:", err)
    }

    fmt.Println("\nThe program is still running - recover() caught the panic!")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Dereferencing a Valid Pointer ===
Got value: 42

=== Dereferencing a nil Pointer (recovered) ===
Error: recovered from: runtime error: invalid memory address or nil pointer dereference

The program is still running - recover() caught the panic!
```

**What if we DIDN'T recover?** Try this two-line variant in a scratch directory:

```go
var p *int
fmt.Println(*p)
```

The program crashes with this real output (captured from an actual run - the hex numbers and goroutine trace vary per run):

```
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x2 addr=0x0 pc=0x102b1b600]

goroutine 1 [running]:
main.main()
	/tmp/demo/main.go:7 +0x20
exit status 2
```

Memorize that first line - you WILL meet it again in real code.

**Learning Objectives:**
- ✅ Reproduce the exact nil-dereference runtime error
- ✅ Use defer + recover() (from Level 11) to survive a panic
- ✅ Understand that recovery is a safety net, not a substitute for nil checks

---

## Exercise 4: Passing by Value vs Passing a Pointer

**Objective:** Prove that `func(n int)` cannot mutate the caller's variable, but `func(n *int)` can

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level12-exercise4
cd ~/projects/level12-exercise4
go mod init level12.example/exercise4
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func doubleValue(n int) {
    n = n * 2
    fmt.Println("  inside doubleValue : n =", n)
}

func doublePointer(n *int) {
    *n = *n * 2
    fmt.Println("  inside doublePointer: *n =", *n)
}

func main() {
    fmt.Println("=== Passing by Value (a copy) ===")
    x := 10
    fmt.Println("before doubleValue : x =", x)
    doubleValue(x)
    fmt.Println("after doubleValue  : x =", x)

    fmt.Println("\n=== Passing a Pointer ===")
    y := 10
    fmt.Println("before doublePointer: y =", y)
    doublePointer(&y)
    fmt.Println("after doublePointer : y =", y)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Passing by Value (a copy) ===
before doubleValue : x = 10
  inside doubleValue : n = 20
after doubleValue  : x = 10

=== Passing a Pointer ===
before doublePointer: y = 10
  inside doublePointer: *n = 20
after doublePointer : y = 20
```

**Learning Objectives:**
- ✅ See that a value parameter doubles its copy while the caller's x stays 10
- ✅ See that a pointer parameter's write reaches the caller's y
- ✅ Understand that even the pointer is passed by value - but a copied address still points home

---

## Exercise 5: Swap Two Variables With Pointers

**Objective:** Write the classic swap function - impossible without pointers

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level12-exercise5
cd ~/projects/level12-exercise5
go mod init level12.example/exercise5
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func swapBroken(a, b int) {
    a, b = b, a
}

func swap(a, b *int) {
    *a, *b = *b, *a
}

func swapStrings(a, b *string) {
    *a, *b = *b, *a
}

func main() {
    fmt.Println("=== swapBroken (copies - does nothing) ===")
    x, y := 1, 2
    swapBroken(x, y)
    fmt.Println("x =", x, "y =", y)

    fmt.Println("\n=== swap (pointers - actually swaps) ===")
    swap(&x, &y)
    fmt.Println("x =", x, "y =", y)

    fmt.Println("\n=== Works for Any Type ===")
    first, second := "hello", "world"
    swapStrings(&first, &second)
    fmt.Println("first =", first, "second =", second)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== swapBroken (copies - does nothing) ===
x = 1 y = 2

=== swap (pointers - actually swaps) ===
x = 2 y = 1

=== Works for Any Type ===
first = world second = hello
```

**Learning Objectives:**
- ✅ Prove a value-based swap silently does nothing
- ✅ Swap through pointers using tuple assignment (*a, *b = *b, *a)
- ✅ Apply the same pattern to a second type (*string)

---

## Exercise 6: Pointers to Structs and Auto-Dereference

**Objective:** Mutate struct fields through a pointer, using both (*p).field and the p.field sugar

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level12-exercise6
cd ~/projects/level12-exercise6
go mod init level12.example/exercise6
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type Player struct {
    name  string
    score int
}

func addPoints(p *Player, points int) {
    p.score += points // automatic dereference - same as (*p).score += points
}

func main() {
    player := Player{name: "Alice", score: 100}
    p := &player

    fmt.Println("=== Explicit Dereference: (*p).field ===")
    fmt.Println("(*p).name  =", (*p).name)
    fmt.Println("(*p).score =", (*p).score)

    fmt.Println("\n=== Automatic Dereference: p.field (idiomatic) ===")
    fmt.Println("p.name  =", p.name)
    fmt.Println("p.score =", p.score)

    fmt.Println("\n=== Mutating a Field Through the Pointer ===")
    p.score = 150
    fmt.Println("player.score =", player.score)

    fmt.Println("\n=== Mutating Through a Function Taking *Player ===")
    addPoints(&player, 25)
    fmt.Println("player.score =", player.score)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Explicit Dereference: (*p).field ===
(*p).name  = Alice
(*p).score = 100

=== Automatic Dereference: p.field (idiomatic) ===
p.name  = Alice
p.score = 100

=== Mutating a Field Through the Pointer ===
player.score = 150

=== Mutating Through a Function Taking *Player ===
player.score = 175
```

**Learning Objectives:**
- ✅ Access struct fields through a pointer both ways - and prefer p.field
- ✅ Mutate a field through the pointer and see the original struct change
- ✅ Preview the *Player function pattern that becomes pointer receivers in Level 14

---

## Exercise 7: new(T) vs &T{}

**Objective:** Compare Go's two one-step "allocate and get a pointer" forms

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level12-exercise7
cd ~/projects/level12-exercise7
go mod init level12.example/exercise7
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type Point struct {
    x, y int
}

func main() {
    fmt.Println("=== new(T): a Pointer to a Zero Value ===")
    pi := new(int)
    fmt.Printf("pi has type %T, *pi = %d\n", pi, *pi)

    ps := new(string)
    fmt.Printf("ps has type %T, *ps = %q\n", ps, *ps)

    pp := new(Point)
    fmt.Printf("pp has type %T, *pp = %+v\n", pp, *pp)

    fmt.Println("\n=== &T{}: a Pointer With Chosen Initial Values ===")
    origin := &Point{}
    corner := &Point{x: 3, y: 4}
    fmt.Printf("origin: %+v\n", *origin)
    fmt.Printf("corner: %+v\n", *corner)

    fmt.Println("\n=== For Zero Values, the Two Are Equivalent ===")
    a := new(Point)
    b := &Point{}
    fmt.Println("*a == *b:", *a == *b)

    fmt.Println("\n=== Both Give You a Usable Pointer ===")
    corner.x = 30
    pp.y = 99
    fmt.Printf("corner: %+v\n", *corner)
    fmt.Printf("*pp   : %+v\n", *pp)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== new(T): a Pointer to a Zero Value ===
pi has type *int, *pi = 0
ps has type *string, *ps = ""
pp has type *main.Point, *pp = {x:0 y:0}

=== &T{}: a Pointer With Chosen Initial Values ===
origin: {x:0 y:0}
corner: {x:3 y:4}

=== For Zero Values, the Two Are Equivalent ===
*a == *b: true

=== Both Give You a Usable Pointer ===
corner: {x:30 y:4}
*pp   : {x:0 y:99}
```

**Learning Objectives:**
- ✅ Use new(T) and see it always yields a pointer to a zero value
- ✅ Use &T{} to allocate AND initialize in one expression
- ✅ Know why &T{} is the idiomatic choice for structs

---

## Exercise 8: A Constructor That Returns a Pointer

**Objective:** Return &localVariable from a function - safe in Go thanks to escape analysis

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level12-exercise8
cd ~/projects/level12-exercise8
go mod init level12.example/exercise8
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type Counter struct {
    count int
}

func newCounter() *Counter {
    c := Counter{count: 0} // a LOCAL variable...
    return &c              // ...safely returned - Go moves it to the heap
}

func increment(c *Counter) {
    c.count++
}

func main() {
    fmt.Println("=== Each Call Creates an Independent Counter ===")
    c1 := newCounter()
    c2 := newCounter()

    increment(c1)
    increment(c1)
    increment(c1)
    increment(c2)

    fmt.Println("c1.count =", c1.count)
    fmt.Println("c2.count =", c2.count)

    fmt.Println("\n=== The Pointer Stays Valid After newCounter Returns ===")
    c3 := newCounter()
    for i := 0; i < 10; i++ {
        increment(c3)
    }
    fmt.Println("c3.count =", c3.count)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Each Call Creates an Independent Counter ===
c1.count = 3
c2.count = 1

=== The Pointer Stays Valid After newCounter Returns ===
c3.count = 10
```

4. (Optional) Watch escape analysis in action:

```bash
go build -gcflags=-m 2>&1 | grep "moved to heap"
```

You'll see the compiler's actual decision (captured from a real run):

```
./main.go:10:2: moved to heap: c
```

Line 10 is `c := Counter{count: 0}` - the compiler noticed `&c` escapes the function and heap-allocated it. In C this pattern is a dangling-pointer bug; in Go it's THE standard constructor idiom.

**Learning Objectives:**
- ✅ Return a pointer to a local variable safely
- ✅ Verify that each constructor call yields an independent value
- ✅ See escape analysis ("moved to heap") with -gcflags=-m

---

## Exercise 9: Slices and Maps Don't Need Pointers

**Objective:** Prove element/entry mutations are shared without pointers - and expose the append gotcha

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level12-exercise9
cd ~/projects/level12-exercise9
go mod init level12.example/exercise9
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func setFirst(nums []int) {
    nums[0] = 999 // element mutation IS visible to the caller
}

func addScore(scores map[string]int, name string, score int) {
    scores[name] = score // map insertion IS visible to the caller
}

func appendBroken(nums []int) {
    nums = append(nums, 42) // only reassigns the LOCAL copy of the header
    fmt.Println("  inside appendBroken:", nums)
}

func appendFixed(nums []int) []int {
    return append(nums, 42)
}

func main() {
    fmt.Println("=== Mutating a Slice Element (visible - no pointer needed) ===")
    nums := []int{1, 2, 3}
    setFirst(nums)
    fmt.Println("after setFirst:", nums)

    fmt.Println("\n=== Adding a Map Entry (visible - no pointer needed) ===")
    scores := map[string]int{"Alice": 95}
    addScore(scores, "Bob", 87)
    fmt.Println("after addScore:", scores)

    fmt.Println("\n=== The append Gotcha (NOT visible) ===")
    fresh := []int{1, 2, 3}
    appendBroken(fresh)
    fmt.Println("after appendBroken:", fresh)

    fmt.Println("\n=== The Fix: Return the New Slice ===")
    fresh = appendFixed(fresh)
    fmt.Println("after appendFixed:", fresh)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Mutating a Slice Element (visible - no pointer needed) ===
after setFirst: [999 2 3]

=== Adding a Map Entry (visible - no pointer needed) ===
after addScore: map[Alice:95 Bob:87]

=== The append Gotcha (NOT visible) ===
  inside appendBroken: [1 2 3 42]
after appendBroken: [1 2 3]

=== The Fix: Return the New Slice ===
after appendFixed: [1 2 3 42]
```

Look closely at the gotcha: *inside* the function the slice really is `[1 2 3 42]` - but the caller still sees `[1 2 3]`. The function only grew its local copy of the slice header (Level 8!). The README's Section 8 shows the even sneakier spare-capacity variant.

**Learning Objectives:**
- ✅ Mutate slice elements and map entries through plain (non-pointer) parameters
- ✅ Reproduce the append-inside-a-function gotcha with real output
- ✅ Fix it idiomatically by returning the new slice

---

## Exercise 10: Comprehensive Practice — Bank Account

**Objective:** Combine constructors, *T parameters, struct mutation, and nil-free design in one program

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level12-exercise10
cd ~/projects/level12-exercise10
go mod init level12.example/exercise10
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type Account struct {
    owner   string
    balance float64
}

func newAccount(owner string, initial float64) *Account {
    return &Account{owner: owner, balance: initial}
}

func deposit(a *Account, amount float64) {
    a.balance += amount
    fmt.Printf("  %s deposited $%.2f -> balance $%.2f\n", a.owner, amount, a.balance)
}

func withdraw(a *Account, amount float64) bool {
    if amount > a.balance {
        fmt.Printf("  %s cannot withdraw $%.2f (balance only $%.2f)\n", a.owner, amount, a.balance)
        return false
    }
    a.balance -= amount
    fmt.Printf("  %s withdrew $%.2f -> balance $%.2f\n", a.owner, amount, a.balance)
    return true
}

func transfer(from, to *Account, amount float64) bool {
    if !withdraw(from, amount) {
        fmt.Println("  transfer cancelled")
        return false
    }
    deposit(to, amount)
    return true
}

func main() {
    alice := newAccount("Alice", 100)
    bob := newAccount("Bob", 50)

    fmt.Println("=== Deposits ===")
    deposit(alice, 25)

    fmt.Println("\n=== Withdrawals ===")
    withdraw(bob, 20)
    withdraw(bob, 500)

    fmt.Println("\n=== Transfer Alice -> Bob ===")
    transfer(alice, bob, 75)

    fmt.Println("\n=== Failed Transfer (insufficient funds) ===")
    transfer(bob, alice, 9999)

    fmt.Println("\n=== Final Balances ===")
    fmt.Printf("%s: $%.2f\n", alice.owner, alice.balance)
    fmt.Printf("%s: $%.2f\n", bob.owner, bob.balance)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Deposits ===
  Alice deposited $25.00 -> balance $125.00

=== Withdrawals ===
  Bob withdrew $20.00 -> balance $30.00
  Bob cannot withdraw $500.00 (balance only $30.00)

=== Transfer Alice -> Bob ===
  Alice withdrew $75.00 -> balance $50.00
  Bob deposited $75.00 -> balance $105.00

=== Failed Transfer (insufficient funds) ===
  Bob cannot withdraw $9999.00 (balance only $105.00)
  transfer cancelled

=== Final Balances ===
Alice: $50.00
Bob: $105.00
```

**Learning Objectives:**
- ✅ Build a constructor returning *Account with &T{} initialization
- ✅ Mutate shared state through *Account parameters (deposits/withdrawals persist!)
- ✅ Compose pointer-taking functions (transfer = withdraw + deposit on the same accounts)
- ✅ Preview the shape that becomes methods with pointer receivers in Level 14

---

## Bonus Challenges

### Challenge 1: Singly-Linked List Insertion

Build a tiny linked list using a `Node` struct (`value int` and `next *Node` fields). Write `insertAfter(n *Node, value int)` that splices a new node in after `n`, then walk the list from the head printing every value.

```bash
mkdir -p ~/projects/level12-bonus1
cd ~/projects/level12-bonus1
go mod init level12.example/bonus1
```

**Hints:**
- The `next` field of the LAST node is nil - that's how you know to stop walking
- Splice order matters: set `newNode.next = n.next` BEFORE `n.next = newNode`
- Walk with `for cur := head; cur != nil; cur = cur.next { ... }`
- A struct can refer to its own type only through a pointer (`next *Node`, not `next Node`)

### Challenge 2: In-Place Matrix Scaling

Given a `[][]float64` matrix, write `scaleCell(p *float64, factor float64)` that multiplies one cell through a pointer, then use nested loops with `&matrix[i][j]` to scale the whole matrix in place. Print the matrix before and after.

```bash
mkdir -p ~/projects/level12-bonus2
cd ~/projects/level12-bonus2
go mod init level12.example/bonus2
```

**Hints:**
- `&matrix[i][j]` is a pointer to the actual element - writes through it stick
- Remember the loop-variable trap: `for _, row := range matrix { for _, v := range row { &v } }` points at COPIES - index with `matrix[i][j]` instead
- Bonus-bonus: also write `scaleMatrix(m [][]float64, factor float64)` and explain why it needs no pointer at all

### Challenge 3: A nil-Safe Accessor Helper

Model an optional configuration: `type Config struct { timeout int }`. Write `timeoutOrDefault(cfg *Config) int` that returns 30 when `cfg` is nil and `cfg.timeout` otherwise. Call it with a real config, and with a nil pointer - no panics allowed.

```bash
mkdir -p ~/projects/level12-bonus3
cd ~/projects/level12-bonus3
go mod init level12.example/bonus3
```

**Hints:**
- The nil check must come FIRST - `cfg.timeout` before `if cfg == nil` defeats the purpose
- Test with `var missing *Config` (nil) and `custom := &Config{timeout: 90}`
- This "absorb nil once, in one place" helper is a very common real-world pattern

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Take addresses with & and dereference with * to read and write
✅ Declare pointer types (*int, *string, *Struct) and recognize nil as their zero value
✅ Reproduce AND recover from the nil-dereference panic
✅ Choose between value parameters (read-only copies) and pointer parameters (mutation)
✅ Mutate struct fields through pointers with the p.field auto-dereference sugar
✅ Allocate with new(T) and &T{}, and prefer &T{} for structs
✅ Write constructors that safely return pointers to locals (escape analysis)
✅ Explain why slices and maps rarely need pointers - and dodge the append gotcha

---

## Next Level

Level 13: Structs
- Defining your own types with `struct`
- Struct literals, field access, and nested structs
- Structs combined with pointers (you just learned these!)
- Comparing, copying, and embedding structs

Great work! You've mastered the hardest idea in the OOP tier! 🚀
