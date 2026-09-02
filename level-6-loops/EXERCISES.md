# Level 6: Loops - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

---

## Exercise 1: The Three-Component for Loop

**Objective:** Write a classic counting loop

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level6-exercise1
cd ~/projects/level6-exercise1
go mod init level6.example/exercise1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== Counting Up ===")
    for i := 1; i <= 5; i++ {
        fmt.Println(i)
    }

    fmt.Println("\n=== Counting Down ===")
    for i := 5; i >= 1; i-- {
        fmt.Println(i)
    }

    fmt.Println("\n=== Stepping By 2 ===")
    for i := 0; i <= 10; i += 2 {
        fmt.Println(i)
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
=== Counting Up ===
1
2
3
4
5

=== Counting Down ===
5
4
3
2
1

=== Stepping By 2 ===
0
2
4
6
8
10
```

**Learning Objectives:**
- ✅ Write a three-component for loop
- ✅ Count up, count down, and step by custom amounts

---

## Exercise 2: while-Style for Loop

**Objective:** Use for as a while loop

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level6-exercise2
cd ~/projects/level6-exercise2
go mod init level6.example/exercise2
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== while-Style Countdown ===")
    count := 5
    for count > 0 {
        fmt.Println(count)
        count--
    }
    fmt.Println("Liftoff!")

    fmt.Println("\n=== Doubling Until a Limit ===")
    value := 1
    for value < 100 {
        fmt.Println(value)
        value *= 2
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
=== while-Style Countdown ===
5
4
3
2
1
Liftoff!

=== Doubling Until a Limit ===
1
2
4
8
16
32
64
```

**Learning Objectives:**
- ✅ Use for with only a condition (Go's while)
- ✅ Understand there's no separate while keyword needed

---

## Exercise 3: Infinite Loop With break

**Objective:** Use an infinite loop with a manual exit condition

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level6-exercise3
cd ~/projects/level6-exercise3
go mod init level6.example/exercise3
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== Infinite Loop With break ===")
    count := 0
    for {
        if count >= 5 {
            break
        }
        fmt.Println(count)
        count++
    }
    fmt.Println("Done")

    fmt.Println("\n=== Simulated Retry Loop ===")
    attempt := 0
    maxAttempts := 3
    for {
        attempt++
        fmt.Printf("Attempt %d\n", attempt)
        if attempt >= maxAttempts {
            fmt.Println("Giving up after max attempts")
            break
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
=== Infinite Loop With break ===
0
1
2
3
4
Done

=== Simulated Retry Loop ===
Attempt 1
Attempt 2
Attempt 3
Giving up after max attempts
```

**Learning Objectives:**
- ✅ Write an infinite for loop
- ✅ Use break to exit when a condition is met inside the body
- ✅ Recognize the retry-loop pattern

---

## Exercise 4: continue

**Objective:** Skip iterations without exiting the loop

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level6-exercise4
cd ~/projects/level6-exercise4
go mod init level6.example/exercise4
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== Only Even Numbers ===")
    for i := 1; i <= 10; i++ {
        if i%2 != 0 {
            continue
        }
        fmt.Println(i)
    }

    fmt.Println("\n=== Skip Multiples of 3 ===")
    for i := 1; i <= 15; i++ {
        if i%3 == 0 {
            continue
        }
        fmt.Println(i)
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
=== Only Even Numbers ===
2
4
6
8
10

=== Skip Multiples of 3 ===
1
2
4
5
7
8
10
11
13
14
```

**Learning Objectives:**
- ✅ Use continue to skip an iteration
- ✅ Combine continue with modulo checks from Level 4

---

## Exercise 5: range Over Slices

**Objective:** Iterate over a slice with index and value

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level6-exercise5
cd ~/projects/level6-exercise5
go mod init level6.example/exercise5
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fruits := []string{"apple", "banana", "cherry", "date"}

    fmt.Println("=== Index and Value ===")
    for i, fruit := range fruits {
        fmt.Printf("%d: %s\n", i, fruit)
    }

    fmt.Println("\n=== Value Only ===")
    for _, fruit := range fruits {
        fmt.Println(fruit)
    }

    fmt.Println("\n=== Index Only ===")
    for i := range fruits {
        fmt.Println(i)
    }

    fmt.Println("\n=== Sum a Slice of Numbers ===")
    numbers := []int{10, 20, 30, 40}
    sum := 0
    for _, n := range numbers {
        sum += n
    }
    fmt.Println("Sum:", sum)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Index and Value ===
0: apple
1: banana
2: cherry
3: date

=== Value Only ===
apple
banana
cherry
date

=== Index Only ===
0
1
2
3

=== Sum a Slice of Numbers ===
Sum: 100
```

**Learning Objectives:**
- ✅ Use range to get both index and value
- ✅ Use _ to ignore the index or value
- ✅ Accumulate a sum with range

---

## Exercise 6: range Over Strings (Runes)

**Objective:** Understand that range gives runes, not bytes

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level6-exercise6
cd ~/projects/level6-exercise6
go mod init level6.example/exercise6
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== range Over an ASCII String ===")
    for i, r := range "Go!" {
        fmt.Printf("index=%d rune=%c (%d)\n", i, r, r)
    }

    fmt.Println("\n=== range Over a String With Multi-Byte Characters ===")
    for i, r := range "héllo" {
        fmt.Printf("index=%d rune=%c\n", i, r)
    }

    fmt.Println("\n=== Byte Length vs Rune Count ===")
    s := "héllo"
    fmt.Printf("len(s) [bytes] = %d\n", len(s))

    runeCount := 0
    for range s {
        runeCount++
    }
    fmt.Printf("rune count = %d\n", runeCount)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== range Over an ASCII String ===
index=0 rune=G (71)
index=1 rune=o (111)
index=2 rune=! (33)

=== range Over a String With Multi-Byte Characters ===
index=0 rune=h
index=1 rune=é
index=3 rune=l
index=4 rune=l
index=5 rune=o

=== Byte Length vs Rune Count ===
len(s) [bytes] = 6
rune count = 5
```

**Learning Objectives:**
- ✅ Use range to iterate a string's runes
- ✅ See that multi-byte runes cause the index to skip
- ✅ Understand len() counts bytes, not characters

---

## Exercise 7: range Over Maps

**Objective:** Iterate a map and produce a stable, sorted order

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level6-exercise7
cd ~/projects/level6-exercise7
go mod init level6.example/exercise7
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "sort"
)

func main() {
    scores := map[string]int{
        "Charlie": 92,
        "Alice":   95,
        "Bob":     87,
    }

    fmt.Println("=== Unsorted range (order not guaranteed) ===")
    count := 0
    for range scores {
        count++
    }
    fmt.Println("Total entries:", count)

    fmt.Println("\n=== Sorted By Name ===")
    names := make([]string, 0, len(scores))
    for name := range scores {
        names = append(names, name)
    }
    sort.Strings(names)

    for _, name := range names {
        fmt.Printf("%s: %d\n", name, scores[name])
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
=== Unsorted range (order not guaranteed) ===
Total entries: 3

=== Sorted By Name ===
Alice: 95
Bob: 87
Charlie: 92
```

**Learning Objectives:**
- ✅ Range over a map's keys and values
- ✅ Understand map iteration order is randomized
- ✅ Collect and sort keys for a predictable order

---

## Exercise 8: Labeled break and continue

**Objective:** Control outer loops from inside nested loops

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level6-exercise8
cd ~/projects/level6-exercise8
go mod init level6.example/exercise8
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("=== Unlabeled break (only exits inner loop) ===")
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if j == 1 {
                break
            }
            fmt.Printf("i=%d j=%d\n", i, j)
        }
    }

    fmt.Println("\n=== Labeled break (exits both loops) ===")
search:
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if i == 1 && j == 1 {
                fmt.Println("found target, stopping")
                break search
            }
            fmt.Printf("i=%d j=%d\n", i, j)
        }
    }

    fmt.Println("\n=== Labeled continue (skips to next outer iteration) ===")
outer:
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if j == 1 {
                continue outer
            }
            fmt.Printf("i=%d j=%d\n", i, j)
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
=== Unlabeled break (only exits inner loop) ===
i=0 j=0
i=1 j=0
i=2 j=0

=== Labeled break (exits both loops) ===
i=0 j=0
i=0 j=1
i=0 j=2
i=1 j=0
found target, stopping

=== Labeled continue (skips to next outer iteration) ===
i=0 j=0
i=1 j=0
i=2 j=0
```

**Learning Objectives:**
- ✅ See that unlabeled break only exits the innermost loop
- ✅ Use a labeled break to exit multiple nested loops
- ✅ Use a labeled continue to skip to the next outer iteration

---

## Exercise 9: Common Loop Patterns

**Objective:** Practice accumulate, find, filter, and count patterns together

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level6-exercise9
cd ~/projects/level6-exercise9
go mod init level6.example/exercise9
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    numbers := []int{4, 8, 15, 16, 23, 42}
    words := []string{"go", "is", "fun", "go", "is", "fast"}

    fmt.Println("=== Accumulate: Sum and Max ===")
    sum := 0
    max := numbers[0]
    for _, n := range numbers {
        sum += n
        if n > max {
            max = n
        }
    }
    fmt.Printf("Sum: %d, Max: %d\n", sum, max)

    fmt.Println("\n=== Find: First Number Over 20 ===")
    found := -1
    for _, n := range numbers {
        if n > 20 {
            found = n
            break
        }
    }
    fmt.Println("First over 20:", found)

    fmt.Println("\n=== Filter: Even Numbers Only ===")
    var evens []int
    for _, n := range numbers {
        if n%2 == 0 {
            evens = append(evens, n)
        }
    }
    fmt.Println("Evens:", evens)

    fmt.Println("\n=== Count: Word Frequency ===")
    counts := make(map[string]int)
    for _, w := range words {
        counts[w]++
    }
    fmt.Println("go:", counts["go"])
    fmt.Println("is:", counts["is"])
    fmt.Println("fast:", counts["fast"])
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Accumulate: Sum and Max ===
Sum: 108, Max: 42

=== Find: First Number Over 20 ===
First over 20: 23

=== Filter: Even Numbers Only ===
Evens: [4 8 16 42]

=== Count: Word Frequency ===
go: 2
is: 2
fast: 1
```

**Learning Objectives:**
- ✅ Combine accumulate, find, filter, and count patterns
- ✅ Build a new slice inside a loop with append
- ✅ Use a map as a frequency counter

---

## Exercise 10: Comprehensive Practice — Prime Number Finder

**Objective:** Combine nested loops, break, continue, and accumulation

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level6-exercise10
cd ~/projects/level6-exercise10
go mod init level6.example/exercise10
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func isPrime(n int) bool {
    if n < 2 {
        return false
    }
    for i := 2; i*i <= n; i++ {
        if n%i == 0 {
            return false
        }
    }
    return true
}

func main() {
    fmt.Println("=== Primes Under 50 ===")
    var primes []int
    for n := 2; n < 50; n++ {
        if !isPrime(n) {
            continue
        }
        primes = append(primes, n)
    }
    fmt.Println(primes)

    fmt.Println("\n=== Multiplication Table (1-5) ===")
    for row := 1; row <= 5; row++ {
        for col := 1; col <= 5; col++ {
            fmt.Printf("%3d", row*col)
        }
        fmt.Println()
    }

    fmt.Println("\n=== First Prime Pair Summing to a Target ===")
    target := 24
found:
    for _, a := range primes {
        for _, b := range primes {
            if a != b && a+b == target {
                fmt.Printf("%d + %d = %d\n", a, b, target)
                break found
            }
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
=== Primes Under 50 ===
[2 3 5 7 11 13 17 19 23 29 31 37 41 43 47]

=== Multiplication Table (1-5) ===
  1  2  3  4  5
  2  4  6  8 10
  3  6  9 12 15
  4  8 12 16 20
  5 10 15 20 25

=== First Prime Pair Summing to a Target ===
5 + 19 = 24
```

**Learning Objectives:**
- ✅ Combine a helper function, nested loops, and a labeled break
- ✅ Build a slice incrementally with continue guarding invalid entries
- ✅ Format a grid of numbers with a nested loop

---

## Bonus Challenges

### Challenge 1: FizzBuzz

Print numbers 1 to 30, but print "Fizz" for multiples of 3, "Buzz" for multiples of 5, and "FizzBuzz" for multiples of both.

```bash
mkdir -p ~/projects/level6-bonus1
cd ~/projects/level6-bonus1
go mod init level6.example/bonus1
```

**Hints:**
- Check divisibility by 15 first (both 3 and 5), then 3, then 5
- Use a single `for` loop from 1 to 30

### Challenge 2: Reverse a Slice In Place

Write a function that reverses a slice of ints without allocating a new slice, using a loop with two indices moving toward each other.

```bash
mkdir -p ~/projects/level6-bonus2
cd ~/projects/level6-bonus2
go mod init level6.example/bonus2
```

**Hints:**
- Use `for left, right := 0, len(s)-1; left < right; left, right = left+1, right-1`
- Swap `s[left]` and `s[right]` each iteration

### Challenge 3: Grid Search

Given a 2D grid (slice of slices) of ints, find the coordinates of a target value using nested loops and a labeled break.

```bash
mkdir -p ~/projects/level6-bonus3
cd ~/projects/level6-bonus3
go mod init level6.example/bonus3
```

**Hints:**
- Use `[][]int` for the grid
- Use a labeled break to stop both loops once found
- Print "not found" if the target never appears

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Write all three for loop forms (three-component, while-style, infinite)
✅ Use break and continue to control loop flow
✅ Range over slices, strings, and maps correctly
✅ Understand rune iteration and map iteration order
✅ Use labeled break/continue to control nested loops
✅ Apply accumulate, find, filter, and count patterns
✅ Combine loops with functions and helper checks in real programs

---

## Next Level

Level 7: Arrays
- Fixed-size arrays in depth
- Array vs slice semantics
- Multi-dimensional arrays
- When to use arrays vs slices

Great work! You're mastering Go! 🚀
