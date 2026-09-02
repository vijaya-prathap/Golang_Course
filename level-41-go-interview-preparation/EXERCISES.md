# Level 41: Go Interview Preparation - Exercises

Complete all exercises in order. Exercises 1-4 verify classic Go gotchas for yourself; Exercises 5-10 are the kind of live-coding problems you'll actually be asked to solve in an interview.

---

## Exercise 1: Prove the nil-Interface Trap For Yourself

**Objective:** Reproduce, in your own program, the single most commonly asked Go gotcha - a nil pointer wrapped inside a non-nil interface

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level41-exercise1
cd ~/projects/level41-exercise1
go mod init level41.example/exercise1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type MyError struct{ msg string }

func (e *MyError) Error() string {
    if e == nil {
        return "<nil *MyError>"
    }
    return e.msg
}

func doWork(fail bool) *MyError {
    if fail {
        return &MyError{msg: "it broke"}
    }
    return nil // a nil *MyError
}

// The classic trap: returning the concrete nil pointer as the `error` interface
func doWorkWrapped(fail bool) error {
    var err *MyError // nil *MyError
    if fail {
        err = &MyError{msg: "it broke"}
    }
    return err // BUG: always returns a non-nil `error`, even when err is nil!
}

func main() {
    fmt.Println("=== Case 1: A True nil error ===")
    var err error
    fmt.Println("err == nil:", err == nil)

    fmt.Println("\n=== Case 2: nil *MyError Assigned Directly (still nil) ===")
    var p *MyError
    fmt.Println("p == nil:", p == nil)

    fmt.Println("\n=== Case 3: THE TRAP - nil *MyError Returned as error Interface ===")
    result := doWorkWrapped(false)
    fmt.Println("result == nil:", result == nil, "<- FALSE! This is the trap.")
    fmt.Printf("result's dynamic type: %T, dynamic value: %v\n", result, result)

    fmt.Println("\n=== Case 4: The Correct Pattern - Return nil Directly ===")
    fixedResult := doWork(false)
    var asError error
    if fixedResult != nil {
        asError = fixedResult
    }
    fmt.Println("asError == nil:", asError == nil, "<- TRUE, correctly nil")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Case 1: A True nil error ===
err == nil: true

=== Case 2: nil *MyError Assigned Directly (still nil) ===
p == nil: true

=== Case 3: THE TRAP - nil *MyError Returned as error Interface ===
result == nil: false <- FALSE! This is the trap.
result's dynamic type: *main.MyError, dynamic value: <nil *MyError>

=== Case 4: The Correct Pattern - Return nil Directly ===
asError == nil: true <- TRUE, correctly nil
```

**Learning Objectives:**
- ✅ Reproduce the nil-interface vs typed-nil trap from Level 15 with your own code
- ✅ Explain why `result == nil` is `false` in terms of an interface's type+value pair
- ✅ Apply the fix: never assign a concrete nil pointer directly to an interface-typed return value

---

## Exercise 2: Slice Aliasing and append() Surprises

**Objective:** Prove that slicing shares an underlying array, and that `append()` can silently overwrite a sibling slice

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level41-exercise2
cd ~/projects/level41-exercise2
go mod init level41.example/exercise2
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== Slicing Shares the Underlying Array ===")
    original := []int{1, 2, 3, 4, 5}
    view := original[1:3]
    fmt.Println("original:", original)
    fmt.Println("view:", view)

    view[0] = 999
    fmt.Println("after view[0] = 999:")
    fmt.Println("original:", original, "(mutated through the shared array!)")
    fmt.Println("view:", view)

    fmt.Println("\n=== append() Aliasing Surprise ===")
    base := make([]int, 3, 5) // len=3, cap=5 - room to grow without reallocating
    base[0], base[1], base[2] = 1, 2, 3
    fmt.Println("base:", base, "len:", len(base), "cap:", cap(base))

    branchA := append(base, 100) // fits in spare capacity, no reallocation
    fmt.Println("branchA:", branchA)

    branchB := append(base, 200) // ALSO fits in the same spare capacity slot!
    fmt.Println("branchB:", branchB)

    fmt.Println("branchA AFTER branchB append:", branchA, "<- silently changed to 200!")

    fmt.Println("\n=== The Safe Fix: copy() for an Independent Slice ===")
    safeBase := make([]int, 3, 5)
    safeBase[0], safeBase[1], safeBase[2] = 1, 2, 3
    independent := make([]int, len(safeBase))
    copy(independent, safeBase)
    independent = append(independent, 999)
    fmt.Println("safeBase:", safeBase)
    fmt.Println("independent:", independent, "(safeBase untouched)")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Slicing Shares the Underlying Array ===
original: [1 2 3 4 5]
view: [2 3]
after view[0] = 999:
original: [1 999 3 4 5] (mutated through the shared array!)
view: [999 3]

=== append() Aliasing Surprise ===
base: [1 2 3] len: 3 cap: 5
branchA: [1 2 3 100]
branchB: [1 2 3 200]
branchA AFTER branchB append: [1 2 3 200] <- silently changed to 200!

=== The Safe Fix: copy() for an Independent Slice ===
safeBase: [1 2 3]
independent: [1 2 3 999] (safeBase untouched)
```

**Learning Objectives:**
- ✅ Reproduce Level 8's slice-aliasing behavior for yourself
- ✅ Explain why two `append()` calls on the same base slice can silently clobber each other
- ✅ Apply `copy()` as the fix for needing a genuinely independent slice

---

## Exercise 3: Verify This Project's Loop Variable Semantics (Go 1.22+)

**Objective:** Confirm which loop-variable-capture behavior this project's actual Go version exhibits, rather than assuming

**Instructions:**

1. Check the course's Go version first:

```bash
cat /Users/macbookpro/go/Golang/go.mod
```

You should see `go 1.26.5` (or later) - well past the Go 1.22 semantics change.

2. Create working directory:

```bash
mkdir -p ~/projects/level41-exercise3
cd ~/projects/level41-exercise3
go mod init level41.example/exercise3
```

3. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "sort"
    "sync"
)

func main() {
    fmt.Println("=== Loop Variable Capture: This Project's Go Version ===")

    var wg sync.WaitGroup
    results := make([]int, 0, 3)
    var mu sync.Mutex

    for _, n := range []int{1, 2, 3} {
        wg.Add(1)
        go func() {
            defer wg.Done()
            mu.Lock()
            results = append(results, n)
            mu.Unlock()
        }()
    }
    wg.Wait()
    sort.Ints(results)
    fmt.Println("captured values (sorted):", results)
    fmt.Println("Go 1.22+ semantics: each iteration gets its OWN copy of n,")
    fmt.Println("so every goroutine captures a different value: [1 2 3], never [3 3 3]")

    fmt.Println("\n=== Classic for i := 0; i < n; i++ Also Gets Per-Iteration Scope ===")
    funcs := make([]func() int, 0, 3)
    for i := 0; i < 3; i++ {
        funcs = append(funcs, func() int { return i })
    }
    for _, f := range funcs {
        fmt.Print(f(), " ")
    }
    fmt.Println("<- prints 0 1 2 on Go 1.22+, would print 3 3 3 pre-1.22")
}
EOF
```

4. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Loop Variable Capture: This Project's Go Version ===
captured values (sorted): [1 2 3]
Go 1.22+ semantics: each iteration gets its OWN copy of n,
so every goroutine captures a different value: [1 2 3], never [3 3 3]

=== Classic for i := 0; i < n; i++ Also Gets Per-Iteration Scope ===
0 1 2 <- prints 0 1 2 on Go 1.22+, would print 3 3 3 pre-1.22
```

**Learning Objectives:**
- ✅ Check a project's `go.mod` language version before assuming which loop-capture semantics apply
- ✅ Confirm this course's toolchain (Go 1.26.5) exhibits Go 1.22+ per-iteration loop variables
- ✅ Explain what the pre-1.22 output would have been and why (one shared variable, read after the loop finished)

---

## Exercise 4: Map Iteration Order and Uncomparable Structs

**Objective:** Verify map iteration order randomization across multiple runs, and reproduce the compile error from comparing structs with slice fields

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level41-exercise4
cd ~/projects/level41-exercise4
go mod init level41.example/exercise4
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "reflect"
    "sort"
)

type Config struct {
    Name string
    Tags []string // a slice field makes Config uncomparable with ==
}

func main() {
    fmt.Println("=== Part 1: Map Iteration Order Is Randomized ===")
    m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5}
    for k := range m {
        fmt.Print(k, " ")
    }
    fmt.Println()

    fmt.Println("Sorted for a stable, repeatable order:")
    keys := make([]string, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    for _, k := range keys {
        fmt.Print(k, " ")
    }
    fmt.Println()

    fmt.Println("\n=== Part 2: Comparing Structs With Uncomparable Fields ===")
    a := Config{Name: "x", Tags: []string{"go"}}
    b := Config{Name: "x", Tags: []string{"go"}}
    fmt.Println("a == b would fail to compile: Config has a []string field")
    fmt.Println("reflect.DeepEqual(a, b):", reflect.DeepEqual(a, b))

    arr1 := [3]int{1, 2, 3}
    arr2 := [3]int{1, 2, 3}
    fmt.Println("arr1 == arr2 (arrays of comparable elements ARE comparable):", arr1 == arr2)
}
EOF
```

3. Run the program at least twice to see the map order change:

```bash
go run main.go
go run main.go
```

4. Now prove the compile error is real by adding a file that actually tries `a == b`:

```bash
cat > broken.go << 'EOF'
package main

func triggerCompileError() bool {
    a := Config{Name: "x", Tags: []string{"go"}}
    b := Config{Name: "x", Tags: []string{"go"}}
    return a == b
}
EOF
go build ./...
```

5. Remove the broken file so the module builds cleanly again:

```bash
rm broken.go
```

**Expected Output (step 3, one representative run - your exact letter order will vary):**

```
=== Part 1: Map Iteration Order Is Randomized ===
e a b c d 
Sorted for a stable, repeatable order:
a b c d e 

=== Part 2: Comparing Structs With Uncomparable Fields ===
a == b would fail to compile: Config has a []string field
reflect.DeepEqual(a, b): true
arr1 == arr2 (arrays of comparable elements ARE comparable): true
```

**Expected Output (step 4, the compile error):**

```
# level41.example/exercise4
./broken.go:6:9: invalid operation: a == b (struct containing []string cannot be compared)
```

**Learning Objectives:**
- ✅ Observe map iteration order changing across separate runs of identical code
- ✅ Reproduce the exact compile-time error Go gives for comparing uncomparable structs
- ✅ Use `reflect.DeepEqual` as the correct fix, and recognize that arrays of comparable elements remain `==`-comparable

---

## Exercise 5: Reverse a String, Correctly Handling Unicode

**Objective:** Implement string reversal that works on any Unicode input, not just ASCII, and prove it with tests

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level41-exercise5
cd ~/projects/level41-exercise5
go mod init level41.example/exercise5
```

2. Create `reverse.go`:

```bash
cat > reverse.go << 'EOF'
package main

import "fmt"

func reverseString(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

func main() {
    fmt.Println(reverseString("hello"))
    fmt.Println(reverseString("héllo"))
    fmt.Println(reverseString("日本語"))
    fmt.Println(reverseString("Go is fun! 🚀"))
}
EOF
```

3. Create `reverse_test.go`:

```bash
cat > reverse_test.go << 'EOF'
package main

import "testing"

func TestReverseString(t *testing.T) {
    cases := []struct {
        name, input, want string
    }{
        {"ascii", "hello", "olleh"},
        {"accented", "héllo", "olléh"},
        {"cjk", "日本語", "語本日"},
        {"emoji", "Go is fun! 🚀", "🚀 !nuf si oG"},
        {"empty", "", ""},
        {"single rune", "x", "x"},
    }
    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            got := reverseString(c.input)
            if got != c.want {
                t.Errorf("reverseString(%q) = %q, want %q", c.input, got, c.want)
            }
        })
    }
}
EOF
```

4. Run the program, then run the tests:

```bash
go run reverse.go
go test -v ./...
```

**Expected Output (`go run reverse.go`):**

```
olleh
olléh
語本日
🚀 !nuf si oG
```

**Expected Output (`go test -v ./...`):**

```
=== RUN   TestReverseString
=== RUN   TestReverseString/ascii
=== RUN   TestReverseString/accented
=== RUN   TestReverseString/cjk
=== RUN   TestReverseString/emoji
=== RUN   TestReverseString/empty
=== RUN   TestReverseString/single_rune
--- PASS: TestReverseString (0.00s)
    --- PASS: TestReverseString/ascii (0.00s)
    --- PASS: TestReverseString/accented (0.00s)
    --- PASS: TestReverseString/cjk (0.00s)
    --- PASS: TestReverseString/emoji (0.00s)
    --- PASS: TestReverseString/empty (0.00s)
    --- PASS: TestReverseString/single_rune (0.00s)
PASS
ok  	level41.example/exercise5	0.648s
```

**Learning Objectives:**
- ✅ Reverse a string by runes, not bytes, bridging Level 9's UTF-8 material
- ✅ Write table-driven tests with `t.Run` subtests, per Level 26's pattern
- ✅ Cover edge cases (empty string, single character) without being asked

---

## Exercise 6: Detect a Cycle in a Linked List

**Objective:** Implement Floyd's cycle-detection algorithm and prove it with tests covering both cyclic and acyclic lists

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level41-exercise6
cd ~/projects/level41-exercise6
go mod init level41.example/exercise6
```

2. Create `cycle.go`:

```bash
cat > cycle.go << 'EOF'
package main

import "fmt"

type Node struct {
    Val  int
    Next *Node
}

func hasCycle(head *Node) bool {
    slow, fast := head, head
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
        if slow == fast {
            return true
        }
    }
    return false
}

func buildList(vals []int) *Node {
    dummy := &Node{}
    cur := dummy
    for _, v := range vals {
        cur.Next = &Node{Val: v}
        cur = cur.Next
    }
    return dummy.Next
}

func main() {
    acyclic := buildList([]int{1, 2, 3})
    fmt.Println("acyclic:", hasCycle(acyclic))
}
EOF
```

3. Create `cycle_test.go`:

```bash
cat > cycle_test.go << 'EOF'
package main

import "testing"

func TestHasCycle(t *testing.T) {
    t.Run("no cycle", func(t *testing.T) {
        list := buildList([]int{1, 2, 3, 4, 5})
        if hasCycle(list) {
            t.Error("expected no cycle")
        }
    })

    t.Run("cycle back to middle", func(t *testing.T) {
        list := buildList([]int{1, 2, 3, 4, 5})
        cur, target := list, list
        for i := 0; i < 2; i++ {
            target = target.Next
        }
        for cur.Next != nil {
            cur = cur.Next
        }
        cur.Next = target
        if !hasCycle(list) {
            t.Error("expected a cycle")
        }
    })

    t.Run("empty list", func(t *testing.T) {
        if hasCycle(nil) {
            t.Error("expected no cycle on nil head")
        }
    })

    t.Run("single self-referencing node", func(t *testing.T) {
        n := &Node{Val: 1}
        n.Next = n
        if !hasCycle(n) {
            t.Error("expected a cycle")
        }
    })
}
EOF
```

4. Run the tests:

```bash
go test -v ./...
```

**Expected Output:**

```
=== RUN   TestHasCycle
=== RUN   TestHasCycle/no_cycle
=== RUN   TestHasCycle/cycle_back_to_middle
=== RUN   TestHasCycle/empty_list
=== RUN   TestHasCycle/single_self-referencing_node
--- PASS: TestHasCycle (0.00s)
    --- PASS: TestHasCycle/no_cycle (0.00s)
    --- PASS: TestHasCycle/cycle_back_to_middle (0.00s)
    --- PASS: TestHasCycle/empty_list (0.00s)
    --- PASS: TestHasCycle/single_self-referencing_node (0.00s)
PASS
ok  	level41.example/exercise6	0.534s
```

**Learning Objectives:**
- ✅ Implement Floyd's tortoise-and-hare algorithm for O(1)-space cycle detection
- ✅ Handle the `nil` head edge case explicitly, without being prompted
- ✅ Build a self-referencing list in a test to exercise the smallest possible cycle

---

## Exercise 7: First Non-Repeating Character

**Objective:** Find the first character in a string that appears exactly once, using a two-pass frequency-map approach

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level41-exercise7
cd ~/projects/level41-exercise7
go mod init level41.example/exercise7
```

2. Create `firstunique.go`:

```bash
cat > firstunique.go << 'EOF'
package main

import "fmt"

func firstNonRepeating(s string) int {
    counts := make(map[rune]int)
    runes := []rune(s)
    for _, r := range runes {
        counts[r]++
    }
    for i, r := range runes {
        if counts[r] == 1 {
            return i
        }
    }
    return -1
}

func main() {
    fmt.Println(firstNonRepeating("leetcode"))
}
EOF
```

3. Create `firstunique_test.go`:

```bash
cat > firstunique_test.go << 'EOF'
package main

import "testing"

func TestFirstNonRepeating(t *testing.T) {
    cases := []struct {
        name  string
        input string
        want  int
    }{
        {"first char unique", "leetcode", 0},
        {"all repeated", "aabb", -1},
        {"middle char unique", "swiss", 1},
        {"empty string", "", -1},
        {"unique at end", "aabbccd", 6},
    }
    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            got := firstNonRepeating(c.input)
            if got != c.want {
                t.Errorf("firstNonRepeating(%q) = %d, want %d", c.input, got, c.want)
            }
        })
    }
}
EOF
```

4. Run the tests:

```bash
go test -v ./...
```

**Expected Output:**

```
=== RUN   TestFirstNonRepeating
=== RUN   TestFirstNonRepeating/first_char_unique
=== RUN   TestFirstNonRepeating/all_repeated
=== RUN   TestFirstNonRepeating/middle_char_unique
=== RUN   TestFirstNonRepeating/empty_string
=== RUN   TestFirstNonRepeating/unique_at_end
--- PASS: TestFirstNonRepeating (0.00s)
    --- PASS: TestFirstNonRepeating/first_char_unique (0.00s)
    --- PASS: TestFirstNonRepeating/all_repeated (0.00s)
    --- PASS: TestFirstNonRepeating/middle_char_unique (0.00s)
    --- PASS: TestFirstNonRepeating/empty_string (0.00s)
    --- PASS: TestFirstNonRepeating/unique_at_end (0.00s)
PASS
ok  	level41.example/exercise7	0.536s
```

**Learning Objectives:**
- ✅ Recognize why a single pass with a map alone can't preserve "first in order"
- ✅ Combine a frequency map (Level 10) with a second ordered pass to solve it in O(n)
- ✅ Test the "no answer exists" case (`-1`) as carefully as the "answer exists" case

---

## Exercise 8: Implement an LRU Cache From Scratch

**Objective:** Build a fixed-capacity least-recently-used cache with O(1) `Get` and `Put`, using a map plus `container/list`

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level41-exercise8
cd ~/projects/level41-exercise8
go mod init level41.example/exercise8
```

2. Create `lru.go`:

```bash
cat > lru.go << 'EOF'
package main

import (
    "container/list"
    "fmt"
)

type LRUCache struct {
    capacity int
    items    map[int]*list.Element
    order    *list.List
}

type entry struct{ key, value int }

func NewLRUCache(capacity int) *LRUCache {
    return &LRUCache{capacity: capacity, items: make(map[int]*list.Element), order: list.New()}
}

func (c *LRUCache) Get(key int) (int, bool) {
    el, ok := c.items[key]
    if !ok {
        return 0, false
    }
    c.order.MoveToFront(el)
    return el.Value.(*entry).value, true
}

func (c *LRUCache) Put(key, value int) {
    if el, ok := c.items[key]; ok {
        el.Value.(*entry).value = value
        c.order.MoveToFront(el)
        return
    }
    if c.order.Len() >= c.capacity {
        back := c.order.Back()
        if back != nil {
            c.order.Remove(back)
            delete(c.items, back.Value.(*entry).key)
        }
    }
    el := c.order.PushFront(&entry{key: key, value: value})
    c.items[key] = el
}

func main() {
    c := NewLRUCache(2)
    c.Put(1, 100)
    fmt.Println(c.Get(1))
}
EOF
```

3. Create `lru_test.go`:

```bash
cat > lru_test.go << 'EOF'
package main

import "testing"

func TestLRUCache(t *testing.T) {
    c := NewLRUCache(2)
    c.Put(1, 100)
    c.Put(2, 200)

    if v, ok := c.Get(1); !ok || v != 100 {
        t.Fatalf("Get(1) = %d, %v; want 100, true", v, ok)
    }

    c.Put(3, 300) // key 2 is now least recently used -> evicted
    if _, ok := c.Get(2); ok {
        t.Error("expected key 2 to be evicted")
    }
    if v, ok := c.Get(3); !ok || v != 300 {
        t.Errorf("Get(3) = %d, %v; want 300, true", v, ok)
    }

    c.Put(4, 400) // key 1 is now least recently used -> evicted
    if _, ok := c.Get(1); ok {
        t.Error("expected key 1 to be evicted")
    }
    if v, ok := c.Get(4); !ok || v != 400 {
        t.Errorf("Get(4) = %d, %v; want 400, true", v, ok)
    }
}

func TestLRUCacheUpdateExisting(t *testing.T) {
    c := NewLRUCache(2)
    c.Put(1, 100)
    c.Put(1, 999) // update, not insert
    if v, ok := c.Get(1); !ok || v != 999 {
        t.Errorf("Get(1) = %d, %v; want 999, true", v, ok)
    }
}
EOF
```

4. Run the tests:

```bash
go test -v ./...
```

**Expected Output:**

```
=== RUN   TestLRUCache
--- PASS: TestLRUCache (0.00s)
=== RUN   TestLRUCacheUpdateExisting
--- PASS: TestLRUCacheUpdateExisting (0.00s)
PASS
ok  	level41.example/exercise8	0.544s
```

**Learning Objectives:**
- ✅ Compose a map (Level 10) with a doubly linked list (`container/list`) to get O(1) time for both operations
- ✅ Implement move-to-front-on-access and evict-from-the-back semantics correctly
- ✅ Test both eviction order and the "update an existing key" path, which is easy to get wrong

---

## Exercise 9: Merge Overlapping Intervals

**Objective:** Sort and merge a list of intervals so no two remaining intervals overlap

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level41-exercise9
cd ~/projects/level41-exercise9
go mod init level41.example/exercise9
```

2. Create `merge.go`:

```bash
cat > merge.go << 'EOF'
package main

import (
    "fmt"
    "sort"
)

func mergeIntervals(intervals [][2]int) [][2]int {
    if len(intervals) == 0 {
        return intervals
    }
    sorted := make([][2]int, len(intervals))
    copy(sorted, intervals)
    sort.Slice(sorted, func(i, j int) bool { return sorted[i][0] < sorted[j][0] })

    result := [][2]int{sorted[0]}
    for _, cur := range sorted[1:] {
        last := &result[len(result)-1]
        if cur[0] <= last[1] {
            if cur[1] > last[1] {
                last[1] = cur[1]
            }
        } else {
            result = append(result, cur)
        }
    }
    return result
}

func main() {
    fmt.Println(mergeIntervals([][2]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}))
}
EOF
```

3. Create `merge_test.go`:

```bash
cat > merge_test.go << 'EOF'
package main

import (
    "reflect"
    "testing"
)

func TestMergeIntervals(t *testing.T) {
    cases := []struct {
        name  string
        input [][2]int
        want  [][2]int
    }{
        {"classic overlapping", [][2]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}, [][2]int{{1, 6}, {8, 10}, {15, 18}}},
        {"touching intervals merge", [][2]int{{1, 4}, {4, 5}}, [][2]int{{1, 5}}},
        {"fully contained interval", [][2]int{{1, 4}, {2, 3}}, [][2]int{{1, 4}}},
        {"empty input", [][2]int{}, [][2]int{}},
        {"unsorted input", [][2]int{{5, 6}, {1, 2}}, [][2]int{{1, 2}, {5, 6}}},
    }
    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            got := mergeIntervals(c.input)
            if !reflect.DeepEqual(got, c.want) {
                t.Errorf("mergeIntervals(%v) = %v, want %v", c.input, got, c.want)
            }
        })
    }
}
EOF
```

4. Run the tests:

```bash
go test -v ./...
```

**Expected Output:**

```
=== RUN   TestMergeIntervals
=== RUN   TestMergeIntervals/classic_overlapping
=== RUN   TestMergeIntervals/touching_intervals_merge
=== RUN   TestMergeIntervals/fully_contained_interval
=== RUN   TestMergeIntervals/empty_input
=== RUN   TestMergeIntervals/unsorted_input
--- PASS: TestMergeIntervals (0.00s)
    --- PASS: TestMergeIntervals/classic_overlapping (0.00s)
    --- PASS: TestMergeIntervals/touching_intervals_merge (0.00s)
    --- PASS: TestMergeIntervals/fully_contained_interval (0.00s)
    --- PASS: TestMergeIntervals/empty_input (0.00s)
    --- PASS: TestMergeIntervals/unsorted_input (0.00s)
PASS
ok  	level41.example/exercise9	0.527s
```

**Learning Objectives:**
- ✅ Sort by interval start, then merge in a single linear pass
- ✅ Handle the "touching but not overlapping" edge case (`[1,4]` and `[4,5]` merge into `[1,5]`)
- ✅ Test with unsorted input to confirm the sort step is actually doing its job

---

## Exercise 10: Implement a Worker Pool From Scratch

**Objective:** Build a bounded-concurrency worker pool that processes jobs across a fixed number of goroutines

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level41-exercise10
cd ~/projects/level41-exercise10
go mod init level41.example/exercise10
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "sort"
    "sync"
)

type Job struct {
    ID    int
    Input int
}

type Result struct {
    JobID  int
    Output int
}

// worker pulls jobs off jobs until it's closed, sending each result to results.
func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
    defer wg.Done()
    for j := range jobs {
        results <- Result{JobID: j.ID, Output: j.Input * j.Input}
    }
}

func main() {
    const numWorkers = 4
    const numJobs = 10

    jobs := make(chan Job, numJobs)
    results := make(chan Result, numJobs)
    var wg sync.WaitGroup

    for w := 1; w <= numWorkers; w++ {
        wg.Add(1)
        go worker(w, jobs, results, &wg)
    }

    for i := 1; i <= numJobs; i++ {
        jobs <- Job{ID: i, Input: i}
    }
    close(jobs) // no more jobs - workers exit their range loop once drained

    go func() {
        wg.Wait()
        close(results) // safe to close only after every worker is done sending
    }()

    collected := make([]Result, 0, numJobs)
    for r := range results {
        collected = append(collected, r)
    }

    sort.Slice(collected, func(i, j int) bool { return collected[i].JobID < collected[j].JobID })
    for _, r := range collected {
        fmt.Printf("job %d -> %d\n", r.JobID, r.Output)
    }
    fmt.Println("total results:", len(collected))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
job 1 -> 1
job 2 -> 4
job 3 -> 9
job 4 -> 16
job 5 -> 25
job 6 -> 36
job 7 -> 49
job 8 -> 64
job 9 -> 81
job 10 -> 100
total results: 10
```

**Learning Objectives:**
- ✅ Build a worker pool with a jobs channel, a results channel, and a fixed number of workers, bridging Level 20/21
- ✅ Explain why `jobs` must be closed for workers' `range jobs` loops to exit
- ✅ Explain why `results` can only be closed after `wg.Wait()` confirms every worker is done sending

---

## Bonus Challenges

### Challenge 1: Rate Limiter From Scratch, Under Time Pressure

Implement a simple token-bucket rate limiter: a struct with an `Allow() bool` method that returns `true` up to N times per time window, then `false` until the window resets. Give yourself 20 minutes, as if this were a live interview.

```bash
mkdir -p ~/projects/level41-bonus1
cd ~/projects/level41-bonus1
go mod init level41.example/bonus1
```

**Hints:**
- Track a token count and a `time.Time` for the last refill
- On each `Allow()` call, check if enough time has passed to refill tokens (`time.Since(lastRefill)`), then decrement a token if one is available
- Protect the struct's state with a `sync.Mutex` if it might be called from multiple goroutines
- Think about what happens at exactly the boundary of a time window - is that an edge case you'd mention out loud to an interviewer?

### Challenge 2: Explain and Fix a Buggy Concurrent Snippet

Below is a broken counter. Explain out loud (or in comments) exactly what's wrong before you fix it - that's the actual interview skill being tested, not just producing a working version.

```go
func main() {
    counter := 0
    var wg sync.WaitGroup
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter++ // <- what's wrong here?
        }()
    }
    wg.Wait()
    fmt.Println(counter) // rarely prints 1000
}
```

```bash
mkdir -p ~/projects/level41-bonus2
cd ~/projects/level41-bonus2
go mod init level41.example/bonus2
```

**Hints:**
- `counter++` is not atomic - it's a read, an increment, and a write, and 1000 goroutines are doing all three concurrently
- Run it with `go run -race main.go` and read what the race detector reports
- Fix it two different ways: a `sync.Mutex` around the increment, and `sync/atomic`'s `atomic.Int64` (Level 23) - be ready to explain when you'd pick one over the other

### Challenge 3: Explain This Code (Using the Course's Own Level 15 Snippet)

Without running it first, write down what you think this prints and why - then run it and see if you were right.

```go
type Shape interface {
    Area() float64
}

type Circle struct {
    Radius float64
}

func (c *Circle) Area() float64 {
    return 3.14159 * c.Radius * c.Radius
}

func describe(s Shape) {
    if s == nil {
        fmt.Println("shape is nil")
        return
    }
    fmt.Println("shape is not nil, area:", s.Area())
}

func main() {
    var c *Circle
    describe(c)
}
```

```bash
mkdir -p ~/projects/level41-bonus3
cd ~/projects/level41-bonus3
go mod init level41.example/bonus3
```

**Hints:**
- This is Exercise 1's trap wearing different clothes - `c` is a nil `*Circle`, but it's passed into `describe` as a `Shape` interface
- Ask: does `s == nil` inside `describe` evaluate to `true` or `false`? Why?
- If `s == nil` is `false`, what happens when `s.Area()` is called on a nil `*Circle` receiver - does it panic? (Look back at how `MyError.Error()` in Exercise 1 handled a nil receiver, and compare)

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Reproduce and explain seven classic Go gotchas with your own verified code, not just recite them
✅ Confirm which loop-variable semantics a specific Go version exhibits instead of assuming
✅ Implement and test five classic live-coding challenges: Unicode string reversal, cycle detection, first non-repeating character, an LRU cache, and merging intervals
✅ Build a worker pool from scratch and explain its channel-closing protocol precisely
✅ Approach an unfamiliar concurrent bug or unfamiliar code snippet the way an interviewer expects: explain first, then fix

---

## You've Completed the Course

There is no Level 42 - this was the last set of exercises in the entire 42-level course. Every gotcha, pattern, and challenge above pulled from levels you've already completed; the only thing new here was compressing all of it into interview shape.

Go build something real with it. 🚀
