# Level 8: Slices - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

---

## Exercise 1: Basic Slice Creation and append()

**Objective:** Create slices four different ways and grow one with append()

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level8-exercise1
cd ~/projects/level8-exercise1
go mod init level8.example/exercise1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== Slice Literal ===")
    nums := []int{10, 20, 30}
    fmt.Println(nums, "len:", len(nums), "cap:", cap(nums))

    fmt.Println("\n=== make() With Length Only ===")
    zeros := make([]int, 4)
    fmt.Println(zeros, "len:", len(zeros), "cap:", cap(zeros))

    fmt.Println("\n=== make() With Length and Capacity ===")
    buf := make([]int, 2, 6)
    fmt.Println(buf, "len:", len(buf), "cap:", cap(buf))

    fmt.Println("\n=== append() Growing a Slice ===")
    nums = append(nums, 40)
    fmt.Println(nums, "len:", len(nums), "cap:", cap(nums))

    fmt.Println("\n=== append() With Multiple Values ===")
    nums = append(nums, 50, 60, 70)
    fmt.Println(nums, "len:", len(nums), "cap:", cap(nums))

    fmt.Println("\n=== append() One Slice Onto Another (... spread) ===")
    more := []int{80, 90}
    nums = append(nums, more...)
    fmt.Println(nums, "len:", len(nums), "cap:", cap(nums))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Slice Literal ===
[10 20 30] len: 3 cap: 3

=== make() With Length Only ===
[0 0 0 0] len: 4 cap: 4

=== make() With Length and Capacity ===
[0 0] len: 2 cap: 6

=== append() Growing a Slice ===
[10 20 30 40] len: 4 cap: 6

=== append() With Multiple Values ===
[10 20 30 40 50 60 70] len: 7 cap: 12

=== append() One Slice Onto Another (... spread) ===
[10 20 30 40 50 60 70 80 90] len: 9 cap: 12
```

**Learning Objectives:**
- ✅ Create slices with a literal, and with make() (with and without a separate capacity)
- ✅ Grow a slice with append(), including multiple values and spreading another slice
- ✅ Observe that cap() can grow ahead of len() to leave room for future appends

---

## Exercise 2: len() vs cap() Exploration

**Objective:** Watch cap() actually grow as you append, and see how a sub-slice's cap differs from its len

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level8-exercise2
cd ~/projects/level8-exercise2
go mod init level8.example/exercise2
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== Watching cap() Grow As We Append ===")
    var s []int
    for i := 0; i < 20; i++ {
        s = append(s, i)
        fmt.Printf("after append %2d: len=%2d cap=%2d\n", i, len(s), cap(s))
    }

    fmt.Println("\n=== len() vs cap() On a Sliced Slice ===")
    backing := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
    view := backing[2:5]
    fmt.Println("backing:", backing, "len:", len(backing), "cap:", cap(backing))
    fmt.Println("view := backing[2:5]:", view, "len:", len(view), "cap:", cap(view))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output** (this is the ACTUAL growth pattern captured from `go run` on this Go toolchain - your own run should match exactly, since slice growth is deterministic for a given Go version/platform):

```
=== Watching cap() Grow As We Append ===
after append  0: len= 1 cap= 4
after append  1: len= 2 cap= 4
after append  2: len= 3 cap= 4
after append  3: len= 4 cap= 4
after append  4: len= 5 cap= 8
after append  5: len= 6 cap= 8
after append  6: len= 7 cap= 8
after append  7: len= 8 cap= 8
after append  8: len= 9 cap=16
after append  9: len=10 cap=16
after append 10: len=11 cap=16
after append 11: len=12 cap=16
after append 12: len=13 cap=16
after append 13: len=14 cap=16
after append 14: len=15 cap=16
after append 15: len=16 cap=16
after append 16: len=17 cap=32
after append 17: len=18 cap=32
after append 18: len=19 cap=32
after append 19: len=20 cap=32

=== len() vs cap() On a Sliced Slice ===
backing: [0 1 2 3 4 5 6 7 8 9] len: 10 cap: 10
view := backing[2:5]: [2 3 4] len: 3 cap: 8
```

**Notice:** the very first append to a nil `[]int` jumps straight to `cap=4`, not `cap=1` - Go rounds allocations to internal size classes, it does not simply double from 1. Don't memorize "always doubles from 1" - the real rule is "roughly doubles for small slices, until it doesn't" (see the STUDY_GUIDE for the larger-scale pattern). Also notice `view`'s capacity (8) has nothing to do with its own length (3) - it reflects how much room is left in `backing` starting from index 2.

**Learning Objectives:**
- ✅ Observe real, actual cap() growth values instead of assuming a formula
- ✅ Understand that a sub-slice's capacity is measured from where IT starts, not from 0
- ✅ Recognize that len() and cap() answer different questions

---

## Exercise 3: Slicing Arrays and Slices

**Objective:** Practice slice expressions and see how the start position affects capacity

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level8-exercise3
cd ~/projects/level8-exercise3
go mod init level8.example/exercise3
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== Slicing an Array ===")
    arr := [6]int{10, 20, 30, 40, 50, 60}
    s := arr[1:4]
    fmt.Println("arr:", arr)
    fmt.Println("arr[1:4]:", s, "len:", len(s), "cap:", cap(s))

    fmt.Println("\n=== Slicing a Slice ===")
    nums := []int{1, 2, 3, 4, 5, 6, 7, 8}
    fmt.Println("nums[:3]:", nums[:3])
    fmt.Println("nums[3:]:", nums[3:])
    fmt.Println("nums[2:5]:", nums[2:5])
    fmt.Println("nums[:]:", nums[:])

    fmt.Println("\n=== cap() Depends on Where the Slice Starts ===")
    fmt.Println("nums:", nums, "len:", len(nums), "cap:", cap(nums))
    mid := nums[3:5]
    fmt.Println("mid := nums[3:5]:", mid, "len:", len(mid), "cap:", cap(mid))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Slicing an Array ===
arr: [10 20 30 40 50 60]
arr[1:4]: [20 30 40] len: 3 cap: 5

=== Slicing a Slice ===
nums[:3]: [1 2 3]
nums[3:]: [4 5 6 7 8]
nums[2:5]: [3 4 5]
nums[:]: [1 2 3 4 5 6 7 8]

=== cap() Depends on Where the Slice Starts ===
nums: [1 2 3 4 5 6 7 8] len: 8 cap: 8
mid := nums[3:5]: [4 5] len: 2 cap: 5
```

**Learning Objectives:**
- ✅ Use `[low:high]` slice expressions on both arrays and slices
- ✅ Understand that omitting `low` or `high` defaults to the start/end
- ✅ See that capacity always extends to the end of the underlying array, from the slice's own starting point

---

## Exercise 4: Aliasing - Sub-Slices Share the Underlying Array

**Objective:** Prove that mutating a sub-slice element affects the original slice

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level8-exercise4
cd ~/projects/level8-exercise4
go mod init level8.example/exercise4
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== Sub-Slices Share the Underlying Array ===")
    original := []int{1, 2, 3, 4, 5}
    sub := original[1:4]
    fmt.Println("original:", original)
    fmt.Println("sub := original[1:4]:", sub)

    fmt.Println("\n--- Mutating sub[0] ---")
    sub[0] = 999
    fmt.Println("sub:", sub)
    fmt.Println("original:", original, "(index 1 changed too!)")

    fmt.Println("\n=== Two Slices, Same Backing Array ===")
    base := []int{100, 200, 300, 400, 500}
    left := base[0:3]
    right := base[2:5]
    fmt.Println("left:", left)
    fmt.Println("right:", right)

    fmt.Println("\n--- left and right overlap at base[2] ---")
    left[2] = 9999
    fmt.Println("after left[2] = 9999:")
    fmt.Println("left:", left)
    fmt.Println("right:", right, "(right[0] changed too - same memory as left[2]/base[2])")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Sub-Slices Share the Underlying Array ===
original: [1 2 3 4 5]
sub := original[1:4]: [2 3 4]

--- Mutating sub[0] ---
sub: [999 3 4]
original: [1 999 3 4 5] (index 1 changed too!)

=== Two Slices, Same Backing Array ===
left: [100 200 300]
right: [300 400 500]

--- left and right overlap at base[2] ---
after left[2] = 9999:
left: [100 200 9999]
right: [9999 400 500] (right[0] changed too - same memory as left[2]/base[2])
```

**Learning Objectives:**
- ✅ Prove with real output that slicing does NOT copy data
- ✅ See that two slices can overlap partially and share only some elements
- ✅ Build the instinct to ask "does anything else alias this array?" before mutating through a slice

---

## Exercise 5: append() Reallocation Behavior

**Objective:** Prove exactly when append() mutates the shared array vs. when it allocates a new one, using pointer addresses as evidence

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level8-exercise5
cd ~/projects/level8-exercise5
go mod init level8.example/exercise5
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== Case 1: append() Within Capacity (mutates shared array) ===")
    base := make([]int, 3, 5)
    base[0], base[1], base[2] = 1, 2, 3
    view := base[:2]
    fmt.Printf("base:  %v len=%d cap=%d ptr=%p\n", base, len(base), cap(base), base)
    fmt.Printf("view:  %v len=%d cap=%d ptr=%p\n", view, len(view), cap(view), view)

    view = append(view, 999)
    fmt.Println("\n--- after view = append(view, 999) (still within cap) ---")
    fmt.Printf("base:  %v len=%d cap=%d ptr=%p\n", base, len(base), cap(base), base)
    fmt.Printf("view:  %v len=%d cap=%d ptr=%p\n", view, len(view), cap(view), view)
    fmt.Println("Same address -> base[2] was overwritten by view's append!")

    fmt.Println("\n=== Case 2: append() Beyond Capacity (reallocates) ===")
    small := make([]int, 2, 2)
    small[0], small[1] = 10, 20
    fmt.Printf("small before: %v len=%d cap=%d ptr=%p\n", small, len(small), cap(small), small)

    grown := append(small, 30)
    fmt.Printf("grown after:  %v len=%d cap=%d ptr=%p\n", grown, len(grown), cap(grown), grown)

    grown[0] = 111
    fmt.Println("\n--- after grown[0] = 111 ---")
    fmt.Println("small:", small, "(unchanged - different backing array now)")
    fmt.Println("grown:", grown)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output** (the exact hex addresses will differ every time you run it - what matters is that `base` and `view` show the SAME address in Case 1, and `small`/`grown` show DIFFERENT addresses in Case 2):

```
=== Case 1: append() Within Capacity (mutates shared array) ===
base:  [1 2 3] len=3 cap=5 ptr=0x671f46834030
view:  [1 2] len=2 cap=5 ptr=0x671f46834030

--- after view = append(view, 999) (still within cap) ---
base:  [1 2 999] len=3 cap=5 ptr=0x671f46834030
view:  [1 2 999] len=3 cap=5 ptr=0x671f46834030
Same address -> base[2] was overwritten by view's append!

=== Case 2: append() Beyond Capacity (reallocates) ===
small before: [10 20] len=2 cap=2 ptr=0x671f4680e080
grown after:  [10 20 30] len=3 cap=4 ptr=0x671f46844000

--- after grown[0] = 111 ---
small: [10 20] (unchanged - different backing array now)
grown: [111 20 30]
```

**Learning Objectives:**
- ✅ Use `%p` to print a slice's backing-array pointer as direct evidence
- ✅ Prove append() mutates the shared array when `len < cap`
- ✅ Prove append() reallocates (and disconnects from the original) when `len == cap`

---

## Exercise 6: copy() for Independent Copies

**Objective:** Use copy() to make a slice that is truly independent of its source

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level8-exercise6
cd ~/projects/level8-exercise6
go mod init level8.example/exercise6
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== copy() Makes an Independent Copy ===")
    src := []int{1, 2, 3, 4, 5}
    dst := make([]int, len(src))
    n := copy(dst, src)
    fmt.Println("src:", src)
    fmt.Println("dst:", dst, "(copied", n, "elements)")

    dst[0] = 999
    fmt.Println("\n--- after dst[0] = 999 ---")
    fmt.Println("src:", src, "(unaffected - dst has its own array)")
    fmt.Println("dst:", dst)

    fmt.Println("\n=== copy() Only Copies min(len(dst), len(src)) Elements ===")
    small := make([]int, 3)
    copied := copy(small, src)
    fmt.Println("small:", small, "(copied", copied, "elements from src)")

    fmt.Println("\n=== Idiomatic Full Copy Pattern ===")
    original := []string{"a", "b", "c"}
    clone := make([]string, len(original))
    copy(clone, original)
    clone[0] = "Z"
    fmt.Println("original:", original)
    fmt.Println("clone:", clone)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== copy() Makes an Independent Copy ===
src: [1 2 3 4 5]
dst: [1 2 3 4 5] (copied 5 elements)

--- after dst[0] = 999 ---
src: [1 2 3 4 5] (unaffected - dst has its own array)
dst: [999 2 3 4 5]

=== copy() Only Copies min(len(dst), len(src)) Elements ===
small: [1 2 3] (copied 3 elements from src)

=== Idiomatic Full Copy Pattern ===
original: [a b c]
clone: [Z b c]
```

**Learning Objectives:**
- ✅ Use copy() to get a slice with no aliasing back to the source
- ✅ Understand copy() is capped by min(len(dst), len(src))
- ✅ Apply the make() + copy() idiom for cloning a slice

---

## Exercise 7: The Three-Index Slice Expression

**Objective:** Use `s[low:high:max]` to prevent an append() from leaking into shared memory

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level8-exercise7
cd ~/projects/level8-exercise7
go mod init level8.example/exercise7
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== Regular Slice Expression Leaks Capacity ===")
    source := []int{1, 2, 3, 4, 5, 6, 7, 8}
    leaky := source[1:3]
    fmt.Println("source:", source)
    fmt.Println("leaky := source[1:3]:", leaky, "len:", len(leaky), "cap:", cap(leaky))

    leaky = append(leaky, 999)
    fmt.Println("\n--- after leaky = append(leaky, 999) ---")
    fmt.Println("leaky:", leaky)
    fmt.Println("source:", source, "(source[3] got clobbered!)")

    fmt.Println("\n=== Three-Index Slice Expression Caps Capacity ===")
    source2 := []int{1, 2, 3, 4, 5, 6, 7, 8}
    safe := source2[1:3:3] // low:high:max -> cap = max-low
    fmt.Println("source2:", source2)
    fmt.Println("safe := source2[1:3:3]:", safe, "len:", len(safe), "cap:", cap(safe))

    safe = append(safe, 999)
    fmt.Println("\n--- after safe = append(safe, 999) ---")
    fmt.Println("safe:", safe)
    fmt.Println("source2:", source2, "(unchanged - append forced a reallocation)")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Regular Slice Expression Leaks Capacity ===
source: [1 2 3 4 5 6 7 8]
leaky := source[1:3]: [2 3] len: 2 cap: 7

--- after leaky = append(leaky, 999) ---
leaky: [2 3 999]
source: [1 2 3 999 5 6 7 8] (source[3] got clobbered!)

=== Three-Index Slice Expression Caps Capacity ===
source2: [1 2 3 4 5 6 7 8]
safe := source2[1:3:3]: [2 3] len: 2 cap: 2

--- after safe = append(safe, 999) ---
safe: [2 3 999]
source2: [1 2 3 4 5 6 7 8] (unchanged - append forced a reallocation)
```

**Learning Objectives:**
- ✅ See a two-index slice's capacity "leak" and let append() corrupt neighboring data
- ✅ Use the three-index form to cap capacity exactly at the slice's own length
- ✅ Prove the capped slice forces a reallocation instead of aliasing

---

## Exercise 8: nil Slice vs Empty Slice

**Objective:** Understand the real difference between a nil slice and an empty slice

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level8-exercise8
cd ~/projects/level8-exercise8
go mod init level8.example/exercise8
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    var nilSlice []int
    emptySlice := []int{}

    fmt.Println("=== nil Slice vs Empty Slice ===")
    fmt.Println("nilSlice:", nilSlice, "len:", len(nilSlice), "cap:", cap(nilSlice))
    fmt.Println("emptySlice:", emptySlice, "len:", len(emptySlice), "cap:", cap(emptySlice))

    fmt.Println("\n=== Comparing to nil ===")
    fmt.Println("nilSlice == nil:", nilSlice == nil)
    fmt.Println("emptySlice == nil:", emptySlice == nil)

    fmt.Println("\n=== %v and %#v Formatting ===")
    fmt.Printf("nilSlice:   %v   %#v\n", nilSlice, nilSlice)
    fmt.Printf("emptySlice: %v   %#v\n", emptySlice, emptySlice)

    fmt.Println("\n=== Both Behave the Same With append() ===")
    nilSlice = append(nilSlice, 1)
    emptySlice = append(emptySlice, 1)
    fmt.Println("nilSlice after append:", nilSlice)
    fmt.Println("emptySlice after append:", emptySlice)

    fmt.Println("\n=== Both Behave the Same With range (zero iterations) ===")
    var anotherNil []int
    count := 0
    for range anotherNil {
        count++
    }
    fmt.Println("iterations over nil slice:", count)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== nil Slice vs Empty Slice ===
nilSlice: [] len: 0 cap: 0
emptySlice: [] len: 0 cap: 0

=== Comparing to nil ===
nilSlice == nil: true
emptySlice == nil: false

=== %v and %#v Formatting ===
nilSlice:   []   []int(nil)
emptySlice: []   []int{}

=== Both Behave the Same With append() ===
nilSlice after append: [1]
emptySlice after append: [1]

=== Both Behave the Same With range (zero iterations) ===
iterations over nil slice: 0
```

**Learning Objectives:**
- ✅ See that len/cap are identical (0) for nil and empty slices
- ✅ Prove `nilSlice == nil` is true while `emptySlice == nil` is false
- ✅ Confirm append(), range, and len() treat both identically

---

## Exercise 9: 2D Slices - Building a Grid Correctly

**Objective:** Build a 2D slice grid, first showing the shared-row bug, then fixing it

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level8-exercise9
cd ~/projects/level8-exercise9
go mod init level8.example/exercise9
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func buildGridWrong(rows, cols int) [][]int {
    row := make([]int, cols)
    grid := make([][]int, rows)
    for i := range grid {
        grid[i] = row // BUG: every row is the SAME underlying slice
    }
    return grid
}

func buildGridRight(rows, cols int) [][]int {
    grid := make([][]int, rows)
    for i := range grid {
        grid[i] = make([]int, cols) // each row gets its own backing array
    }
    return grid
}

func printGrid(grid [][]int) {
    for _, row := range grid {
        fmt.Println(row)
    }
}

func main() {
    fmt.Println("=== The Shared-Row Bug ===")
    wrong := buildGridWrong(3, 3)
    wrong[0][0] = 1
    fmt.Println("After setting wrong[0][0] = 1:")
    printGrid(wrong)
    fmt.Println("(every row shows the change - they all share ONE underlying array!)")

    fmt.Println("\n=== The Correct Way ===")
    right := buildGridRight(3, 3)
    right[0][0] = 1
    right[1][1] = 2
    right[2][2] = 3
    fmt.Println("After setting the diagonal:")
    printGrid(right)
    fmt.Println("(each row is independent)")

    fmt.Println("\n=== 2D Slice Literal ===")
    grid := [][]int{
        {1, 2, 3},
        {4, 5, 6},
        {7, 8, 9},
    }
    printGrid(grid)
    sum := 0
    for _, row := range grid {
        for _, v := range row {
            sum += v
        }
    }
    fmt.Println("sum of all elements:", sum)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== The Shared-Row Bug ===
After setting wrong[0][0] = 1:
[1 0 0]
[1 0 0]
[1 0 0]
(every row shows the change - they all share ONE underlying array!)

=== The Correct Way ===
After setting the diagonal:
[1 0 0]
[0 2 0]
[0 0 3]
(each row is independent)

=== 2D Slice Literal ===
[1 2 3]
[4 5 6]
[7 8 9]
sum of all elements: 45
```

**Learning Objectives:**
- ✅ Reproduce the classic shared-row 2D slice bug and see its actual (wrong) output
- ✅ Fix it by allocating a fresh backing array per row
- ✅ Use nested range to process a 2D slice

---

## Exercise 10: Comprehensive Practice - In-Memory Task List

**Objective:** Combine everything from this level in a small, realistic slice-backed data structure

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level8-exercise10
cd ~/projects/level8-exercise10
go mod init level8.example/exercise10
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

type Task struct {
    ID   int
    Name string
    Done bool
}

func addTask(tasks []Task, id int, name string) []Task {
    return append(tasks, Task{ID: id, Name: name, Done: false})
}

func completeTask(tasks []Task, id int) {
    for i := range tasks {
        if tasks[i].ID == id {
            tasks[i].Done = true
            return
        }
    }
}

func removeTask(tasks []Task, id int) []Task {
    result := tasks[:0]
    for _, t := range tasks {
        if t.ID != id {
            result = append(result, t)
        }
    }
    return result
}

func filterPending(tasks []Task) []Task {
    var pending []Task
    for _, t := range tasks {
        if !t.Done {
            pending = append(pending, t)
        }
    }
    return pending
}

func printTasks(label string, tasks []Task) {
    fmt.Println(label)
    if len(tasks) == 0 {
        fmt.Println("  (no tasks)")
        return
    }
    for _, t := range tasks {
        status := " "
        if t.Done {
            status = "x"
        }
        fmt.Printf("  [%s] #%d %s\n", status, t.ID, t.Name)
    }
}

func main() {
    var tasks []Task
    tasks = addTask(tasks, 1, "Write README")
    tasks = addTask(tasks, 2, "Verify exercises")
    tasks = addTask(tasks, 3, "Build study guide")
    tasks = addTask(tasks, 4, "Publish level")

    printTasks("=== All Tasks ===", tasks)

    completeTask(tasks, 2)
    completeTask(tasks, 4)
    printTasks("\n=== After Completing #2 and #4 ===", tasks)

    pending := filterPending(tasks)
    printTasks("\n=== Pending Tasks ===", pending)

    tasks = removeTask(tasks, 3)
    printTasks("\n=== After Removing #3 ===", tasks)

    fmt.Printf("\nTotal tasks: %d, Pending: %d\n", len(tasks), len(filterPending(tasks)))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== All Tasks ===
  [ ] #1 Write README
  [ ] #2 Verify exercises
  [ ] #3 Build study guide
  [ ] #4 Publish level

=== After Completing #2 and #4 ===
  [ ] #1 Write README
  [x] #2 Verify exercises
  [ ] #3 Build study guide
  [x] #4 Publish level

=== Pending Tasks ===
  [ ] #1 Write README
  [ ] #3 Build study guide

=== After Removing #3 ===
  [ ] #1 Write README
  [x] #2 Verify exercises
  [x] #4 Publish level

Total tasks: 3, Pending: 1
```

**Note:** `removeTask` uses the classic `tasks[:0]` in-place-filter pattern - `result` starts as a zero-length view over the SAME backing array as `tasks`, and each kept task is appended back into that same array. This works safely here because the write position in `result` never runs ahead of the read position in the range loop.

**Learning Objectives:**
- ✅ Combine structs, slices, and functions into a small realistic program
- ✅ Use append() to add, a filtered rebuild to remove, and a fresh slice to filter
- ✅ Practice passing slices to functions and returning the updated slice (the append() pattern)

---

## Bonus Challenges

### Challenge 1: Int-Slice Stack (push/pop)

Implement a stack of ints using only `append` and slicing - no separate package needed.

```bash
mkdir -p ~/projects/level8-bonus1
cd ~/projects/level8-bonus1
go mod init level8.example/bonus1
```

**Verified reference solution:**

```go
package main

import "fmt"

type IntStack struct {
    items []int
}

func (s *IntStack) Push(v int) {
    s.items = append(s.items, v)
}

func (s *IntStack) Pop() (int, bool) {
    if len(s.items) == 0 {
        return 0, false
    }
    last := len(s.items) - 1
    v := s.items[last]
    s.items = s.items[:last]
    return v, true
}

func (s *IntStack) Peek() (int, bool) {
    if len(s.items) == 0 {
        return 0, false
    }
    return s.items[len(s.items)-1], true
}

func (s *IntStack) IsEmpty() bool {
    return len(s.items) == 0
}

func main() {
    var s IntStack
    s.Push(1)
    s.Push(2)
    s.Push(3)
    fmt.Println("stack after pushing 1, 2, 3:", s.items)

    top, _ := s.Peek()
    fmt.Println("peek:", top)

    for !s.IsEmpty() {
        v, _ := s.Pop()
        fmt.Println("popped:", v, "remaining:", s.items)
    }

    _, ok := s.Pop()
    fmt.Println("pop on empty stack ok:", ok)
}
```

**Actual verified output:**

```
stack after pushing 1, 2, 3: [1 2 3]
peek: 3
popped: 3 remaining: [1 2]
popped: 2 remaining: [1]
popped: 1 remaining: []
pop on empty stack ok: false
```

**Hints:**
- `Push` is just `append`
- `Pop` reads the last element, then re-slices to `s.items[:last]` to shrink by one
- Guard against popping an empty stack

### Challenge 2: Efficiently Remove an Element From a Slice

Write two versions: one that preserves order (O(n)), and one that doesn't but runs in O(1).

```bash
mkdir -p ~/projects/level8-bonus2
cd ~/projects/level8-bonus2
go mod init level8.example/bonus2
```

**Verified reference solution:**

```go
package main

import "fmt"

// removeAt removes the element at index i, preserving order.
// It reuses the underlying array - no new allocation.
func removeAt(s []int, i int) []int {
    return append(s[:i], s[i+1:]...)
}

// removeAtFast removes the element at index i WITHOUT preserving order.
// O(1) instead of O(n) - swaps the last element into the gap.
func removeAtFast(s []int, i int) []int {
    s[i] = s[len(s)-1]
    return s[:len(s)-1]
}

func main() {
    fmt.Println("=== removeAt (order-preserving) ===")
    nums := []int{10, 20, 30, 40, 50}
    fmt.Println("before:", nums)
    nums = removeAt(nums, 2)
    fmt.Println("after removeAt(nums, 2):", nums)

    fmt.Println("\n=== removeAtFast (order NOT preserved, O(1)) ===")
    fast := []int{10, 20, 30, 40, 50}
    fmt.Println("before:", fast)
    fast = removeAtFast(fast, 1)
    fmt.Println("after removeAtFast(fast, 1):", fast)
}
```

**Actual verified output:**

```
=== removeAt (order-preserving) ===
before: [10 20 30 40 50]
after removeAt(nums, 2): [10 20 40 50]

=== removeAtFast (order NOT preserved, O(1)) ===
before: [10 20 30 40 50]
after removeAtFast(fast, 1): [10 50 30 40]
```

**Hints:**
- `append(s[:i], s[i+1:]...)` shifts everything after index `i` back by one
- The fast version overwrites index `i` with the last element, then shrinks by one - no shifting at all

### Challenge 3: In-Place Slice Reversal

Reverse a slice of ints without allocating a new slice.

```bash
mkdir -p ~/projects/level8-bonus3
cd ~/projects/level8-bonus3
go mod init level8.example/bonus3
```

**Verified reference solution:**

```go
package main

import "fmt"

func reverse(s []int) {
    for left, right := 0, len(s)-1; left < right; left, right = left+1, right-1 {
        s[left], s[right] = s[right], s[left]
    }
}

func main() {
    nums := []int{1, 2, 3, 4, 5}
    fmt.Println("before:", nums)
    reverse(nums)
    fmt.Println("after: ", nums)

    evenLen := []int{1, 2, 3, 4, 5, 6}
    fmt.Println("\nbefore:", evenLen)
    reverse(evenLen)
    fmt.Println("after: ", evenLen)

    single := []int{42}
    reverse(single)
    fmt.Println("\nsingle-element slice reversed:", single)

    var empty []int
    reverse(empty)
    fmt.Println("empty slice reversed:", empty)
}
```

**Actual verified output:**

```
before: [1 2 3 4 5]
after:  [5 4 3 2 1]

before: [1 2 3 4 5 6]
after:  [6 5 4 3 2 1]

single-element slice reversed: [42]
empty slice reversed: []
```

**Hints:**
- Use two indices moving toward each other: `left, right := 0, len(s)-1`
- Swap `s[left]` and `s[right]` each iteration with `s[left], s[right] = s[right], s[left]`
- Confirm it also handles single-element and empty slices without panicking

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Create slices with literals, make(), and slicing expressions
✅ Explain the real difference between len() and cap()
✅ Predict (and prove) exactly when append() mutates in place vs. reallocates
✅ Demonstrate that sub-slices alias their underlying array
✅ Use the three-index slice expression to prevent aliasing surprises
✅ Use copy() to create a truly independent slice
✅ Explain the nil-vs-empty-slice distinction and why it rarely matters in practice
✅ Build a correct 2D slice grid, avoiding the shared-row bug
✅ Combine slices, structs, and functions into a realistic small program

---

## Next Level

Level 9: Strings & Runes
- Strings as immutable byte sequences
- Runes and Unicode code points
- Converting between strings, []byte, and []rune
- Common string-processing patterns

Great work! You're mastering Go! 🚀
