# Level 7: Arrays - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

---

## Exercise 1: Declaring Arrays

**Objective:** Practice every way to declare an array

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level7-exercise1
cd ~/projects/level7-exercise1
go mod init level7.example/exercise1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== var Declaration (Zero-Valued Array) ===")
    var nums [5]int
    fmt.Println(nums)

    fmt.Println("\n=== Array Literal With Explicit Size ===")
    colors := [3]string{"Red", "Green", "Blue"}
    fmt.Println(colors)

    fmt.Println("\n=== Array Literal With Size Inference [...] ===")
    primes := [...]int{2, 3, 5, 7, 11}
    fmt.Println(primes)
    fmt.Printf("Type-inferred length: %d\n", len(primes))

    fmt.Println("\n=== Sparse Index Literal ===")
    sparse := [5]int{0: 1, 4: 5}
    fmt.Println(sparse)

    fmt.Println("\n=== Sparse Index Literal With [...] ===")
    sparseInferred := [...]int{2: 100, 5: 200}
    fmt.Println(sparseInferred)
    fmt.Printf("Inferred length: %d\n", len(sparseInferred))

    fmt.Println("\n=== var With Type Then Assign Later ===")
    var grades [4]float64
    grades[0] = 91.5
    grades[3] = 88.0
    fmt.Println(grades)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== var Declaration (Zero-Valued Array) ===
[0 0 0 0 0]

=== Array Literal With Explicit Size ===
[Red Green Blue]

=== Array Literal With Size Inference [...] ===
[2 3 5 7 11]
Type-inferred length: 5

=== Sparse Index Literal ===
[1 0 0 0 5]

=== Sparse Index Literal With [...] ===
[0 0 100 0 0 200]
Inferred length: 6

=== var With Type Then Assign Later ===
[91.5 0 0 88]
```

**Learning Objectives:**
- ✅ Declare arrays with `var`, literals, `[...]`, and sparse index literals
- ✅ Understand that `[...]` infers length from the highest index or element count
- ✅ See that unset elements default to the zero value

---

## Exercise 2: Indexing, len(), and Zero Values

**Objective:** Practice indexing and understand array bounds checking

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level7-exercise2
cd ~/projects/level7-exercise2
go mod init level7.example/exercise2
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== Indexing ===")
    fruits := [3]string{"apple", "banana", "cherry"}
    fmt.Println("fruits[0]:", fruits[0])
    fmt.Println("fruits[1]:", fruits[1])
    fmt.Println("fruits[2]:", fruits[2])

    fmt.Println("\n=== len() ===")
    fmt.Println("len(fruits):", len(fruits))

    fmt.Println("\n=== Modifying an Element ===")
    fruits[1] = "blueberry"
    fmt.Println(fruits)

    fmt.Println("\n=== Zero Values By Type ===")
    var ints [3]int
    var floats [3]float64
    var strs [3]string
    var bools [3]bool
    fmt.Println("ints:  ", ints)
    fmt.Println("floats:", floats)
    fmt.Println("strs:  ", strs)
    fmt.Println("bools: ", bools)

    fmt.Println("\n=== Out-of-Range Access Is Caught (commented out) ===")
    // fmt.Println(fruits[3])
    // A CONSTANT out-of-range index like the line above is a COMPILE error:
    //   invalid argument: index 3 out of bounds [0:3]
    badIndex := 3 // a variable index defeats the compile-time check
    fmt.Println("Accessing fruits[badIndex] would panic at RUNTIME instead:")
    fmt.Println("  panic: runtime error: index out of range [3] with length 3")
    _ = badIndex
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Indexing ===
fruits[0]: apple
fruits[1]: banana
fruits[2]: cherry

=== len() ===
len(fruits): 3

=== Modifying an Element ===
[apple blueberry cherry]

=== Zero Values By Type ===
ints:   [0 0 0]
floats: [0 0 0]
strs:   [  ]
bools:  [false false false]

=== Out-of-Range Access Is Caught (commented out) ===
Accessing fruits[badIndex] would panic at RUNTIME instead:
  panic: runtime error: index out of range [3] with length 3
```

**Learning Objectives:**
- ✅ Index into an array and modify elements
- ✅ Use len() to get an array's fixed size
- ✅ See the zero value for int, float64, string, and bool arrays
- ✅ Understand that a constant out-of-range index is a compile error, while a variable one panics at runtime

---

## Exercise 3: Arrays Are Value Types — Assignment Copies

**Objective:** Prove that assigning an array copies it, unlike a slice

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level7-exercise3
cd ~/projects/level7-exercise3
go mod init level7.example/exercise3
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== Array Assignment Copies the Whole Array ===")
    a := [3]int{1, 2, 3}
    b := a // COPY, not a reference
    b[0] = 99

    fmt.Println("a:", a)
    fmt.Println("b:", b)

    fmt.Println("\n=== Contrast: Slice Assignment Shares the Underlying Array ===")
    sa := []int{1, 2, 3}
    sb := sa // sb points at the SAME underlying array
    sb[0] = 99

    fmt.Println("sa:", sa)
    fmt.Println("sb:", sb)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Array Assignment Copies the Whole Array ===
a: [1 2 3]
b: [99 2 3]

=== Contrast: Slice Assignment Shares the Underlying Array ===
sa: [99 2 3]
sb: [99 2 3]
```

**Learning Objectives:**
- ✅ Prove with real output that `b := a` on an array creates an independent copy
- ✅ See the exact opposite behavior with the same syntax on a slice
- ✅ Understand this is THE defining difference between arrays and slices

---

## Exercise 4: Arrays Are Value Types — Function Parameters

**Objective:** Prove that passing an array to a function passes a copy

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level7-exercise4
cd ~/projects/level7-exercise4
go mod init level7.example/exercise4
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func zeroOutArray(arr [3]int) {
    for i := range arr {
        arr[i] = 0
    }
    fmt.Println("  inside zeroOutArray (param):", arr)
}

func zeroOutSlice(s []int) {
    for i := range s {
        s[i] = 0
    }
    fmt.Println("  inside zeroOutSlice (param):", s)
}

func main() {
    fmt.Println("=== Passing an Array to a Function (COPY) ===")
    arr := [3]int{1, 2, 3}
    fmt.Println("before call:", arr)
    zeroOutArray(arr)
    fmt.Println("after call: ", arr, "(unchanged! the function only zeroed its own copy)")

    fmt.Println("\n=== Contrast: Passing a Slice to a Function (SHARED) ===")
    s := []int{1, 2, 3}
    fmt.Println("before call:", s)
    zeroOutSlice(s)
    fmt.Println("after call: ", s, "(changed! slices share the underlying array)")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Passing an Array to a Function (COPY) ===
before call: [1 2 3]
  inside zeroOutArray (param): [0 0 0]
after call:  [1 2 3] (unchanged! the function only zeroed its own copy)

=== Contrast: Passing a Slice to a Function (SHARED) ===
before call: [1 2 3]
  inside zeroOutSlice (param): [0 0 0]
after call:  [0 0 0] (changed! slices share the underlying array)
```

**Learning Objectives:**
- ✅ Prove that a function receiving an array parameter cannot mutate the caller's array
- ✅ Prove that a function receiving a slice parameter CAN mutate the caller's data
- ✅ Understand why slices are the idiomatic choice for function parameters

---

## Exercise 5: Comparing Arrays with ==

**Objective:** Practice array comparison and understand type-size rules

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level7-exercise5
cd ~/projects/level7-exercise5
go mod init level7.example/exercise5
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== Comparing Equal Arrays ===")
    a := [3]int{1, 2, 3}
    b := [3]int{1, 2, 3}
    fmt.Println("a == b:", a == b)

    fmt.Println("\n=== Comparing Different Arrays (Same Type) ===")
    c := [3]int{1, 2, 99}
    fmt.Println("a == c:", a == c)

    fmt.Println("\n=== Comparing Arrays of Strings ===")
    names1 := [2]string{"Alice", "Bob"}
    names2 := [2]string{"Alice", "Bob"}
    names3 := [2]string{"Alice", "Charlie"}
    fmt.Println("names1 == names2:", names1 == names2)
    fmt.Println("names1 == names3:", names1 == names3)

    fmt.Println("\n=== Arrays as Map Keys (only possible because they're comparable) ===")
    seen := map[[2]int]string{}
    seen[[2]int{0, 0}] = "origin"
    seen[[2]int{1, 1}] = "diagonal"
    fmt.Println(seen[[2]int{0, 0}])
    fmt.Println(seen[[2]int{1, 1}])

    // NOTE: [3]int and [4]int are DIFFERENT types.
    // The following would NOT compile if uncommented:
    // d := [4]int{1, 2, 3, 4}
    // fmt.Println(a == d) // compile error: mismatched types [3]int and [4]int
    fmt.Println("\n[3]int and [4]int are different types - comparing them is a compile error")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Comparing Equal Arrays ===
a == b: true

=== Comparing Different Arrays (Same Type) ===
a == c: false

=== Comparing Arrays of Strings ===
names1 == names2: true
names1 == names3: false

=== Arrays as Map Keys (only possible because they're comparable) ===
origin
diagonal

[3]int and [4]int are different types - comparing them is a compile error
```

**Learning Objectives:**
- ✅ Use == and != to compare arrays element by element
- ✅ Use arrays as map keys, something slices cannot do
- ✅ Understand that array size is part of the type, so different sizes can't be compared

---

## Exercise 6: Multi-Dimensional Arrays

**Objective:** Declare, fill, and index a 2D array

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level7-exercise6
cd ~/projects/level7-exercise6
go mod init level7.example/exercise6
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== Declaring a 2D Array ===")
    var grid [3][3]int
    fmt.Println(grid)

    fmt.Println("\n=== Filling a 2D Array by Index ===")
    for row := 0; row < 3; row++ {
        for col := 0; col < 3; col++ {
            grid[row][col] = row*3 + col
        }
    }
    fmt.Println(grid)

    fmt.Println("\n=== Printing Row by Row ===")
    for _, row := range grid {
        fmt.Println(row)
    }

    fmt.Println("\n=== 2D Array Literal ===")
    board := [2][3]string{
        {"X", "O", "X"},
        {"O", "X", "O"},
    }
    fmt.Println(board)
    fmt.Println("board[0][1]:", board[0][1])
    fmt.Println("board[1][2]:", board[1][2])

    fmt.Println("\n=== len() on 2D Arrays ===")
    fmt.Println("len(board):", len(board), "(number of rows)")
    fmt.Println("len(board[0]):", len(board[0]), "(number of columns)")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Declaring a 2D Array ===
[[0 0 0] [0 0 0] [0 0 0]]

=== Filling a 2D Array by Index ===
[[0 1 2] [3 4 5] [6 7 8]]

=== Printing Row by Row ===
[0 1 2]
[3 4 5]
[6 7 8]

=== 2D Array Literal ===
[[X O X] [O X O]]
board[0][1]: O
board[1][2]: O

=== len() on 2D Arrays ===
len(board): 2 (number of rows)
len(board[0]): 3 (number of columns)
```

**Learning Objectives:**
- ✅ Declare a 2D array with `[N][M]T`
- ✅ Fill and index a 2D array with `grid[row][col]`
- ✅ Write a 2D array literal with nested `{}`
- ✅ Understand that `len()` on a 2D array gives the outer dimension only

---

## Exercise 7: Iterating Arrays with range

**Objective:** Practice range over 1D and 2D arrays, and confirm range copies values

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level7-exercise7
cd ~/projects/level7-exercise7
go mod init level7.example/exercise7
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== range Over a 1D Array ===")
    scores := [4]int{85, 92, 78, 95}
    for i, score := range scores {
        fmt.Printf("index=%d score=%d\n", i, score)
    }

    fmt.Println("\n=== range Ignoring the Index ===")
    sum := 0
    for _, score := range scores {
        sum += score
    }
    fmt.Println("sum:", sum)

    fmt.Println("\n=== range Over a 2D Array ===")
    grid := [2][2]int{
        {1, 2},
        {3, 4},
    }
    for rowIndex, row := range grid {
        for colIndex, value := range row {
            fmt.Printf("grid[%d][%d] = %d\n", rowIndex, colIndex, value)
        }
    }

    fmt.Println("\n=== range Copies Each Element (value semantics) ===")
    nums := [3]int{1, 2, 3}
    for _, v := range nums {
        v *= 100 // modifies the LOCAL copy only
        _ = v
    }
    fmt.Println("nums after range loop (unchanged):", nums)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== range Over a 1D Array ===
index=0 score=85
index=1 score=92
index=2 score=78
index=3 score=95

=== range Ignoring the Index ===
sum: 350

=== range Over a 2D Array ===
grid[0][0] = 1
grid[0][1] = 2
grid[1][0] = 3
grid[1][1] = 4

=== range Copies Each Element (value semantics) ===
nums after range loop (unchanged): [1 2 3]
```

**Learning Objectives:**
- ✅ Use range with index and value over a 1D array
- ✅ Nest range loops over a 2D array
- ✅ Confirm that mutating a range variable never affects the source array

---

## Exercise 8: Array vs Slice — Converting and Sharing Memory

**Objective:** See how slicing an array shares its memory, unlike copying it

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level7-exercise8
cd ~/projects/level7-exercise8
go mod init level7.example/exercise8
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== Converting an Array to a Slice ===")
    arr := [5]int{10, 20, 30, 40, 50}
    sliceOfArr := arr[:] // slices the WHOLE array, shares memory
    fmt.Println("arr:        ", arr)
    fmt.Println("sliceOfArr: ", sliceOfArr)

    fmt.Println("\n=== A Slice Derived From an Array Shares Its Memory ===")
    sliceOfArr[0] = 999
    fmt.Println("arr after modifying sliceOfArr[0]:", arr)
    fmt.Println("sliceOfArr:                        ", sliceOfArr)

    fmt.Println("\n=== A Partial Slice of an Array ===")
    partial := arr[1:3]
    fmt.Println("partial (arr[1:3]):", partial)

    fmt.Println("\n=== Arrays Passed by Value vs Slices Passed by Reference-Like Header ===")
    fmt.Println("len(arr):", len(arr), "cap does not apply to arrays directly")
    fmt.Println("len(sliceOfArr):", len(sliceOfArr), "cap(sliceOfArr):", cap(sliceOfArr))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Converting an Array to a Slice ===
arr:         [10 20 30 40 50]
sliceOfArr:  [10 20 30 40 50]

=== A Slice Derived From an Array Shares Its Memory ===
arr after modifying sliceOfArr[0]: [999 20 30 40 50]
sliceOfArr:                         [999 20 30 40 50]

=== A Partial Slice of an Array ===
partial (arr[1:3]): [20 30]

=== Arrays Passed by Value vs Slices Passed by Reference-Like Header ===
len(arr): 5 cap does not apply to arrays directly
len(sliceOfArr): 5 cap(sliceOfArr): 5
```

**Learning Objectives:**
- ✅ Convert an array to a slice with `arr[:]`
- ✅ See that a slice derived from an array shares its underlying memory
- ✅ Take a partial slice of an array with `arr[low:high]`
- ✅ Understand `len`/`cap` apply differently to arrays vs slices

---

## Exercise 9: Tic-Tac-Toe Board (2D Array Project)

**Objective:** Build a small tic-tac-toe board using a `[3][3]string` array

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level7-exercise9
cd ~/projects/level7-exercise9
go mod init level7.example/exercise9
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func printBoard(board [3][3]string) {
    for _, row := range board {
        for _, cell := range row {
            if cell == "" {
                fmt.Print(". ")
            } else {
                fmt.Print(cell + " ")
            }
        }
        fmt.Println()
    }
}

func checkWinner(board [3][3]string) string {
    lines := [][3][2]int{
        {{0, 0}, {0, 1}, {0, 2}},
        {{1, 0}, {1, 1}, {1, 2}},
        {{2, 0}, {2, 1}, {2, 2}},
        {{0, 0}, {1, 0}, {2, 0}},
        {{0, 1}, {1, 1}, {2, 1}},
        {{0, 2}, {1, 2}, {2, 2}},
        {{0, 0}, {1, 1}, {2, 2}},
        {{0, 2}, {1, 1}, {2, 0}},
    }

    for _, line := range lines {
        a := board[line[0][0]][line[0][1]]
        b := board[line[1][0]][line[1][1]]
        c := board[line[2][0]][line[2][1]]
        if a != "" && a == b && b == c {
            return a
        }
    }
    return ""
}

func main() {
    var board [3][3]string
    board[0][0] = "X"
    board[1][1] = "X"
    board[2][2] = "X"
    board[0][1] = "O"
    board[0][2] = "O"

    fmt.Println("=== Tic-Tac-Toe Board ===")
    printBoard(board)

    fmt.Println("\n=== Checking For a Winner ===")
    winner := checkWinner(board)
    if winner != "" {
        fmt.Printf("Winner: %s\n", winner)
    } else {
        fmt.Println("No winner yet")
    }

    fmt.Println("\n=== Board Is a Value Type: Copies Don't Affect the Original ===")
    boardCopy := board
    boardCopy[0][0] = "O"
    fmt.Println("original board[0][0]:", board[0][0])
    fmt.Println("boardCopy[0][0]:     ", boardCopy[0][0])
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Tic-Tac-Toe Board ===
X O O 
. X . 
. . X 

=== Checking For a Winner ===
Winner: X

=== Board Is a Value Type: Copies Don't Affect the Original ===
original board[0][0]: X
boardCopy[0][0]:      O
```

**Learning Objectives:**
- ✅ Model a fixed-size grid with a `[3][3]string` array
- ✅ Pass a 2D array to functions (as a copy) for reading
- ✅ Check diagonal/row/column win conditions with array indexing
- ✅ Reconfirm value semantics: copying a board never affects the original

---

## Exercise 10: Comprehensive Practice — Sudoku Row Validator & Scoreboard

**Objective:** Combine everything: fixed-size arrays, indexing, and range in realistic mini-programs

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level7-exercise10
cd ~/projects/level7-exercise10
go mod init level7.example/exercise10
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

// A Sudoku row is represented as a fixed-size [9]int array.
// Valid values are 1-9, each appearing at most once. 0 means empty.
func isValidRow(row [9]int) bool {
    var seen [10]bool // index 0 unused, 1-9 track seen digits
    for _, val := range row {
        if val == 0 {
            continue // empty cells are always allowed
        }
        if val < 1 || val > 9 {
            return false
        }
        if seen[val] {
            return false // duplicate digit
        }
        seen[val] = true
    }
    return true
}

func describeRow(label string, row [9]int) {
    valid := isValidRow(row)
    fmt.Printf("%s: %v -> valid=%v\n", label, row, valid)
}

func main() {
    fmt.Println("=== Sudoku Row Validator ([9]int arrays) ===")

    complete := [9]int{5, 3, 4, 6, 7, 8, 9, 1, 2}
    describeRow("complete, no duplicates", complete)

    partial := [9]int{5, 3, 0, 0, 7, 0, 0, 0, 0}
    describeRow("partial, with zeros    ", partial)

    duplicate := [9]int{5, 3, 4, 6, 7, 8, 9, 1, 5}
    describeRow("has duplicate 5        ", duplicate)

    var empty [9]int
    describeRow("all empty (zero value) ", empty)

    fmt.Println("\n=== Fixed-Size Scoreboard ===")
    names := [4]string{"Alice", "Bob", "Charlie", "Dana"}
    var scoreboard [4]int
    scoreboard[0] = 95
    scoreboard[1] = 82
    scoreboard[2] = 77
    scoreboard[3] = 88

    total := 0
    highScore := scoreboard[0]
    highScorer := names[0]
    for i, score := range scoreboard {
        total += score
        if score > highScore {
            highScore = score
            highScorer = names[i]
        }
    }

    for i, name := range names {
        fmt.Printf("%-10s %d\n", name, scoreboard[i])
    }
    fmt.Printf("Total: %d\n", total)
    fmt.Printf("Average: %.2f\n", float64(total)/float64(len(scoreboard)))
    fmt.Printf("High score: %s with %d\n", highScorer, highScore)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Sudoku Row Validator ([9]int arrays) ===
complete, no duplicates: [5 3 4 6 7 8 9 1 2] -> valid=true
partial, with zeros    : [5 3 0 0 7 0 0 0 0] -> valid=true
has duplicate 5        : [5 3 4 6 7 8 9 1 5] -> valid=false
all empty (zero value) : [0 0 0 0 0 0 0 0 0] -> valid=true

=== Fixed-Size Scoreboard ===
Alice      95
Bob        82
Charlie    77
Dana       88
Total: 342
Average: 85.50
High score: Alice with 95
```

**Learning Objectives:**
- ✅ Model a fixed-size domain rule (9-digit Sudoku row) with a `[9]int` array and a `[10]bool` lookup array
- ✅ Combine range, indexing, and helper functions on fixed-size arrays
- ✅ Build a small scoreboard program using parallel fixed-size arrays

---

## Bonus Challenges

### Challenge 1: Matrix Transpose

Write a function that transposes a `2x3` matrix (a `[2][3]int`) into a `3x2` matrix (a `[3][2]int`).

```bash
mkdir -p ~/projects/level7-bonus1
cd ~/projects/level7-bonus1
go mod init level7.example/bonus1
```

**Hints:**
- The transpose of `m[i][j]` goes to `result[j][i]`
- The result type is `[3][2]int` - note rows and columns swap in the type itself
- Verify by printing both matrices and checking `result[j][i] == m[i][j]`

### Challenge 2: Fixed-Size Circular Buffer

Build a ring buffer of capacity 5, backed by a `[5]int` array, that overwrites the oldest value once full.

```bash
mkdir -p ~/projects/level7-bonus2
cd ~/projects/level7-bonus2
go mod init level7.example/bonus2
```

**Hints:**
- Track a `start` index and a `count` of how many slots are filled
- `Push` writes at `(start+count) % len(data)`; once full, advance `start` to drop the oldest value
- Use the modulo operator (`%`) to wrap indices around the fixed-size array

### Challenge 3: Array-Based Checksum

Write a `checksum` function that takes a fixed `[8]byte` array and folds it into a `[4]byte` accumulator by adding each byte into `acc[i%4]`.

```bash
mkdir -p ~/projects/level7-bonus3
cd ~/projects/level7-bonus3
go mod init level7.example/bonus3
```

**Hints:**
- Declare the accumulator as `var acc [4]byte`
- Use `range` to walk the input array while indexing into the accumulator with `i % 4`
- Confirm that the same input always produces the same checksum (arrays are comparable!) with `==`

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Declare arrays with every available syntax, including sparse literals
✅ Index, assign, and use len() on both 1D and multi-dimensional arrays
✅ Prove - not just recite - that arrays copy on assignment and on function calls
✅ Contrast array value semantics with slice reference-like semantics
✅ Compare arrays with == and use them as map keys
✅ Build and iterate multi-dimensional arrays for grid-shaped data
✅ Decide when a fixed-size array is the right tool instead of a slice

---

## Next Level

Level 8: Slices
- Slices in depth: the pointer/length/capacity header
- append, growth, and re-slicing
- Why slices share memory - and when that surprises you

Great work! You're mastering Go! 🚀
