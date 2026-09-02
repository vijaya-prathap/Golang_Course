# Level 12: Pointers - Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: Pointer Foundations

    Day 1: What pointers store; address-of & and dereference *
    Day 2: Pointer types and nil safety
    Day 3: Value parameters vs pointer parameters
    Day 4: Swapping and mutating values through pointers
    Day 5: Pointers to structs and automatic dereference
    Day 6: new(T), &T{}, constructors, and escape analysis
    Day 7: Slices/maps, append, and best practices

### Week 2: Practice & Application

    Day 1: Exercises 1-3: basic pointers, nil, and panic recovery
    Day 2: Exercises 4-6: mutation, swap, and struct pointers
    Day 3: Exercises 7-8: allocation and safe constructors
    Day 4: Exercises 9-10: slice/map rules and bank account
    Day 5: Bonus linked list, matrix, and nil-safe helper
    Day 6-7: Explain each diagram and review the checklist

---

## 🎯 What a Pointer Is

Every variable has a location in memory. A pointer stores that location:

    x
    ┌──────────┐
    │    42    │
    └──────────┘
    address: 0xA0

    p
    ┌──────────┐
    │  0xA0    │ ───────────────────────────▶ x's box
    └──────────┘

    p := &x       p receives the address of x
    *p            follows that address to the value 42
    *p = 100      writes through the address; x becomes 100

The hexadecimal address is an implementation detail. Its exact value changes
between runs. The important fact is whether two pointers point at the same box.

---

## 🧩 The Two Meanings of *

    var p *int
          │
          └─ type position: "pointer to int"

    value := *p
              │
              └─ expression position: "value at the address in p"

    *p = 9
     │
     └─ write to the pointed-at value

This is one symbol with two context-dependent meanings. In a declaration it
describes a type; in an expression it dereferences a pointer.

---

## 🚫 nil: A Pointer With No Target

    var p *int

    p
    ┌──────────┐
    │   nil    │ ──▶ nothing
    └──────────┘

    p == nil    // true

Trying to follow a nil pointer crashes:

    *p
     │
     └─ no target exists → runtime panic

Safe boundary check:

    func printValue(p *int) {
        if p == nil {
            fmt.Println("no value")
            return
        }
        fmt.Println(*p)
    }

The check comes before any p.field or *p operation. Once a function has
returned early on nil, the remaining body can use the pointer freely.

---

## 🔁 Why a Pointer Parameter Can Mutate

Go passes every argument by value. Compare the actual copies:

    doubleValue(x)

    caller                         function
    x = 10       ──copy value──▶   n = 10
                                      │
                                      └─ n = 20 changes only this box

    caller x remains 10

    doublePointer(&y)

    caller                         function
    y = 10 at address A
        ▲                             n = address A
        └────────write via *n─────────┘

    *n = 20 writes into y's original box

The pointer itself is copied, but an address copy still identifies the same
value. This is pass-by-value with a shared target, not pass-by-reference.

---

## 🔄 Swap With Pointers

    x = 1                           y = 2
    address A                       address B

    swap(&x, &y)
           │    │
           │    └─ b receives B
           └────── a receives A

    *a, *b = *b, *a

    x = 2                           y = 1

With value parameters, a and b would only be separate copies. A function must
receive addresses to exchange caller-owned values.

---

## 🧱 Pointers to Structs

    type Player struct {
        name  string
        score int
    }

    player := Player{name: "Alice", score: 100}
    p := &player

    p ──▶ ┌─────────────────────┐
          │ name:  "Alice"      │
          │ score: 100          │
          └─────────────────────┘

    p.score = 150

    player.score is now 150

Go automatically dereferences p when selecting a field:

    p.score       preferred spelling
    (*p).score    same effect, but visually noisy

This lets pointer-taking functions read and mutate fields cleanly:

    func addPoints(p *Player, points int) {
        p.score += points
    }

---

## 🆕 Allocation Forms

    new(Point)

    ┌───────────────┐
    │ Point{x:0,y:0}│
    └───────────────┘
             ▲
             │ pointer result

    &Point{x:3, y:4}

    ┌───────────────┐
    │ Point{x:3,y:4}│
    └───────────────┘
             ▲
             │ pointer result

Both forms produce a pointer. The difference is initialization:

| Form | Value initially stored | Typical choice |
|---|---|---|
| new(T) | Always T's zero value | occasional basic type use |
| &T{} | Zero struct value | idiomatic zero struct pointer |
| &T{fields} | chosen fields | idiomatic initialized struct pointer |

For structs, use &T{} because it remains clear as the value gains fields.

---

## 🏗️ Returning a Local Pointer Is Safe

    func newCounter() *Counter {
        c := Counter{}
        return &c
    }

    call newCounter
          │
          ▼
    c is created
          │
          ▼
    &c is returned
          │
          ▼
    caller retains a reachable pointer

The compiler sees that c must outlive the function. It arranges storage safely
(often described by the diagnostic "moved to heap"), and the garbage collector
keeps c alive while the returned pointer can be reached.

    c1 := newCounter()   independent Counter #1
    c2 := newCounter()   independent Counter #2

Each constructor call creates a distinct value.

---

## 🧺 Slices: Shared Elements, Copied Headers

A slice is a small descriptor:

    slice header
    ┌────────────┬─────┬─────┐
    │ data ptr   │ len │ cap │
    └─────┬──────┴─────┴─────┘
          │
          ▼
    backing array: [1 2 3 ...]

When a slice is passed to a function, the header is copied:

    caller header ─┐
                   ├─ both headers point to the same backing array
    function header┘

Therefore an element change is visible:

    func setFirst(nums []int) { nums[0] = 999 }

Both headers still point at the modified array.

---

## ⚠️ The append Gotcha

    caller:   header {data A, len 3, cap 4}
              backing array [1 2 3 _]

    function receives a COPY of the header

    items = append(items, 42)

    function: header {data A, len 4, cap 4}
              backing array [1 2 3 42]

    caller:   header {data A, len 3, cap 4}
              sees [1 2 3]

Even when append reused the backing array, the caller's length is still 3.
If append allocated another backing array, caller and function additionally
point at different arrays. In either case, the caller cannot rely on the
append unless it receives the returned header.

    func addItem(items []int, v int) []int {
        return append(items, v)
    }

    items = addItem(items, 42)

---

## 🗺️ Maps Also Need No Pointer for Entries

    scores := map[string]int{"Alice": 95}

    addScore(scores, "Bob", 87)

    scores now contains Alice and Bob

Maps are reference-like values. A function can add, update, and delete entries
using a normal map parameter. A pointer to a map is almost never needed.

---

## 🧭 Pointer Decision Tree

    Does this function need to update the caller's basic value or struct?
    ├─ YES → accept *T and document nil expectations.
    └─ NO
       │
       Does it only mutate elements of a slice or map entries?
       ├─ YES → use []T or map[K]V directly; no pointer.
       └─ NO
          │
          Does it append to a slice?
          ├─ YES → return the new slice.
          └─ NO → pass small read-only values by value.

---

## 🚨 Range Variable Address Trap

    for _, v := range players {
        modify(&v) // v is the iteration variable, not a slice element
    }

The address above is not an address within players. To mutate actual elements,
use indexes:

    for i := range players {
        modify(&players[i])
    }

This rule is about whether the pointer names the desired value, not merely
about modern loop-variable capture semantics.

---

## 📈 Progression Summary

1. Address and dereference: &x and *p connect a pointer with its target.
2. nil: a zero pointer has no target and must not be dereferenced.
3. Pointer parameters: copied addresses can mutate a caller's original value.
4. Struct pointers: p.field is convenient automatic dereference.
5. Allocation: &T{} is the normal struct-pointer form; new(T) creates zeroes.
6. Lifetime: returned pointers to locals are safe through escape analysis.
7. Slices/maps: element and entry writes are already shared; append returns a
   new slice header.

---

## ✅ Checklist Before Level 13

- [ ] I can distinguish &x, *p, and *int.
- [ ] I know why dereferencing nil panics and how to prevent it.
- [ ] I can explain pass-by-value even for pointer parameters.
- [ ] I can write a pointer-based swap.
- [ ] I write p.field for a struct pointer.
- [ ] I prefer &T{...} for a struct pointer.
- [ ] I know why &local is safe to return in Go.
- [ ] I can explain why append needs a returned slice.
- [ ] I avoid pointers to ordinary slices and maps.
