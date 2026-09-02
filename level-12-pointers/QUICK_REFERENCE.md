# Level 12: Pointers - Quick Reference Card

## 🚀 Core Syntax

    x := 42
    p := &x        // p has type *int; it stores x's address
    value := *p    // read the value at p: 42
    *p = 100       // write through p; x is now 100

Read the operators this way:

| Syntax | Meaning |
|---|---|
| &x | address of x |
| *p in an expression | value p points to |
| *int in a type | pointer to int |
| nil | a pointer with no target |

## 📌 Pointer Types and nil

    var ip *int
    var sp *string

    fmt.Println(ip == nil) // true

Always check a possibly-nil pointer before dereferencing it:

    if ip != nil {
        fmt.Println(*ip)
    }

Dereferencing nil panics:

    var p *int
    fmt.Println(*p) // panic

## 🔁 Value Parameter vs Pointer Parameter

    func doubleValue(n int) {
        n *= 2 // modifies a copy
    }

    func doublePointer(n *int) {
        *n *= 2 // modifies caller's value
    }

    x := 10
    doubleValue(x)    // x remains 10
    doublePointer(&x) // x becomes 20

Go still passes a copied pointer value. The copied address points to the same
underlying value.

## 🔄 Swap

    func swap(a, b *int) {
        *a, *b = *b, *a
    }

    x, y := 1, 2
    swap(&x, &y)
    // x == 2, y == 1

## 🧱 Struct Pointers

    type Player struct {
        name  string
        score int
    }

    p := &Player{name: "Alice", score: 100}
    p.score += 25  // automatic dereference

    // (*p).score += 25 also works, but p.score is idiomatic.

    func addPoints(p *Player, points int) {
        p.score += points
    }

## 🆕 new(T) vs &T{}

| Form | Result | Best use |
|---|---|---|
| new(int) | *int pointing to 0 | occasional pointer to a basic zero value |
| new(Point) | *Point pointing to zero value | valid, but unusual for structs |
| &Point{} | *Point pointing to zero Point | idiomatic struct allocation |
| &Point{x: 3, y: 4} | *Point with initialized fields | preferred constructor-style form |

    type Point struct{ x, y int }

    origin := &Point{}
    corner := &Point{x: 3, y: 4}

## 🏗️ Safe Pointer Constructors

    func newCounter() *Counter {
        c := Counter{}
        return &c // safe; Go handles lifetime/escape analysis
    }

Never worry that this returns an invalid pointer. Go keeps the value alive as
long as it remains reachable.

## 🧺 Slices and Maps

Element and entry changes are visible without pointers:

    func setFirst(nums []int) {
        nums[0] = 999
    }

    func addScore(scores map[string]int) {
        scores["Bob"] = 87
    }

Appending must return the new slice:

    func addItem(items []int, value int) []int {
        return append(items, value)
    }

    items = addItem(items, 42)

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---|---|---|
| Dereference nil | *p when p may be nil | check p != nil first |
| Expect value mutation | func f(x int) { x++ } | func f(x *int) { *x++ } |
| Noisy struct field syntax | (*p).score | p.score |
| Pointer to slice for append | func f(s *[]int) | return append(s, value) |
| Address of range copy | for _, v := range xs { &v } | for i := range xs { &xs[i] } |
| Pointer arithmetic | p++ | index a slice instead |
| new for initialized struct | p := new(Player); p.name = "A" | p := &Player{name: "A"} |

## ✅ Decision Guide

Use *T when:

- The function must mutate the caller's basic value or struct.
- Copying a large struct would be unnecessarily expensive.
- nil usefully means "absent" or "not configured."

Avoid *T when:

- The value is small and read-only.
- You only change elements of a slice or entries of a map.
- You append to a slice: return its new header instead.

## 🎓 Before Level 13

- [ ] I can explain &x and *p without looking.
- [ ] I can explain why pointers are still passed by value.
- [ ] I nil-check a pointer before dereferencing when nil is possible.
- [ ] I can write a function that mutates a caller's struct.
- [ ] I use p.field for fields of a struct pointer.
- [ ] I know why returning &localVar is safe in Go.
- [ ] I know element mutations of slices/maps need no pointer.
- [ ] I return a slice after append inside a helper function.
