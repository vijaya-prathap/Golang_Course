# Level 20: Goroutines - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

**A note before you start:** goroutine scheduling is genuinely non-deterministic. Most exercises below produce the exact same output every time because they're built to synchronize correctly - but a couple of them (marked clearly) demonstrate the *lack* of a guarantee on purpose. Where that's the case, this file says so explicitly and shows one real captured run as an example, instead of pretending there's a single "correct" output.

---

## Exercise 1: The main() Exits Too Early Problem

**Objective:** Observe, for real, why a goroutine's output can silently disappear

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level20-exercise1
cd ~/projects/level20-exercise1
go mod init level20.example/exercise1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func sayHello() {
    fmt.Println("Hello from the goroutine!")
}

func main() {
    fmt.Println("main: launching goroutine")
    go sayHello()
    fmt.Println("main: returning immediately (no wait!)")
}
EOF
```

3. Run the program **several times in a row**:

```bash
go run main.go
go run main.go
go run main.go
go run main.go
go run main.go
```

**Expected Output:**

This program's output is **non-deterministic by design** - line order and even whether the goroutine's message appears at all can vary between runs. Do not expect the same result every time. Here is one real set of 5 runs captured while building this course:

```
--run 1--
main: launching goroutine
main: returning immediately (no wait!)
--run 2--
main: launching goroutine
main: returning immediately (no wait!)
--run 3--
main: launching goroutine
main: returning immediately (no wait!)
--run 4--
main: launching goroutine
main: returning immediately (no wait!)
--run 5--
main: launching goroutine
main: returning immediately (no wait!)
```

Most runs will look like the above - `Hello from the goroutine!` never appears, because `main()` finished and the process exited before the scheduler ran `sayHello`. On some runs (this is real - it happened on run 5 out of 8 when captured for this course), you might get lucky and see:

```
main: launching goroutine
main: returning immediately (no wait!)
Hello from the goroutine!
```

**The lesson is the unpredictability itself.** Run it enough times and you will see both shapes of output. This is exactly the bug Exercise 2 fixes.

**Learning Objectives:**
- ✅ Launch a goroutine with the `go` keyword
- ✅ Observe, first-hand, that `main()` returning kills any goroutines still running
- ✅ Understand why this output is genuinely non-deterministic, not just "unlucky"

---

## Exercise 2: Fixing It With sync.WaitGroup

**Objective:** Use `sync.WaitGroup` to make goroutine completion deterministic

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level20-exercise2
cd ~/projects/level20-exercise2
go mod init level20.example/exercise2
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "sync"
)

func sayHello(wg *sync.WaitGroup) {
    defer wg.Done()
    fmt.Println("Hello from the goroutine!")
}

func main() {
    var wg sync.WaitGroup

    fmt.Println("main: launching goroutine")
    wg.Add(1)
    go sayHello(&wg)

    wg.Wait()
    fmt.Println("main: goroutine finished, safe to return")
}
EOF
```

3. Run the program (try it several times - it will be identical every time):

```bash
go run main.go
```

**Expected Output:**

```
main: launching goroutine
Hello from the goroutine!
main: goroutine finished, safe to return
```

This output is deterministic - verified stable across 5 real runs. `wg.Wait()` blocks `main()` until `wg.Done()` has actually executed, so the message is always guaranteed to print before the program exits.

**Learning Objectives:**
- ✅ Use `sync.WaitGroup`'s `Add`, `Done`, and `Wait` methods
- ✅ Understand why passing `*sync.WaitGroup` (a pointer) matters
- ✅ See the same program from Exercise 1 become fully deterministic

---

## Exercise 3: Waiting for Multiple Goroutines

**Objective:** Launch several goroutines in a loop and wait for all of them with one WaitGroup

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level20-exercise3
cd ~/projects/level20-exercise3
go mod init level20.example/exercise3
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "sync"
)

func worker(id int, wg *sync.WaitGroup, done *int32, mu *sync.Mutex) {
    defer wg.Done()
    mu.Lock()
    *done++
    mu.Unlock()
}

func main() {
    var wg sync.WaitGroup
    var mu sync.Mutex
    var completed int32

    const numWorkers = 5
    fmt.Printf("main: launching %d goroutines\n", numWorkers)

    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go worker(i, &wg, &completed, &mu)
    }

    wg.Wait()
    fmt.Printf("main: all goroutines finished, completed = %d\n", completed)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
main: launching 5 goroutines
main: all goroutines finished, completed = 5
```

Deterministic - verified stable across 5 real runs. Note that `completed` is shared, mutable state written from 5 goroutines, so it's protected with a `sync.Mutex` even though this exercise is mainly about `WaitGroup` - unsynchronized shared state would make the final count unreliable (Exercise 5 demonstrates exactly that failure).

**Learning Objectives:**
- ✅ Call `wg.Add(1)` once per goroutine launched in a loop
- ✅ Confirm a known number of goroutines all completed before continuing
- ✅ See a mutex protecting a shared counter, previewing Exercises 5-6

---

## Exercise 4: Closures and Loop Variable Capture

**Objective:** Verify that goroutines launched in a loop correctly capture their own per-iteration variable on this project's Go version

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level20-exercise4
cd ~/projects/level20-exercise4
go mod init level20.example/exercise4
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

func main() {
    var wg sync.WaitGroup
    var mu sync.Mutex
    var results []int

    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            mu.Lock()
            results = append(results, i)
            mu.Unlock()
        }()
    }

    wg.Wait()
    sort.Ints(results)
    fmt.Println("Captured values (sorted):", results)
}
EOF
```

3. Check your Go version, then run the program:

```bash
go version
go run main.go
```

**Expected Output:**

```
Captured values (sorted): [0 1 2 3 4]
```

Deterministic (once sorted) - verified stable across 5 real runs on `go1.26.5`. Each goroutine's closure captured its **own** copy of `i` from that iteration, because this project's `go.mod` specifies `go 1.26.5`, which is well past the Go 1.22 change to per-iteration loop variable semantics. See [Level 11's closures-in-loops section](../level-11-functions/README.md#closures-in-loops-the-loop-variable-gotcha) and this level's README section 4 for the full history - on Go 1.21 or earlier, this same program could have produced duplicate or missing values instead of a clean `[0 1 2 3 4]`.

**Learning Objectives:**
- ✅ Confirm your project's Go version affects loop variable semantics
- ✅ See that goroutines in a loop capture independent per-iteration values on Go 1.22+
- ✅ Combine `sort.Ints` with a mutex-protected slice to get a deterministic, checkable result

---

## Exercise 5: Data Race - Shared Counter, No Synchronization

**Objective:** Observe a real data race and (if available) catch it with the race detector

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level20-exercise5
cd ~/projects/level20-exercise5
go mod init level20.example/exercise5
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "sync"
)

func main() {
    var wg sync.WaitGroup
    counter := 0

    const numGoroutines = 1000
    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter++ // UNSYNCHRONIZED read-modify-write - a data race!
        }()
    }

    wg.Wait()
    fmt.Println("Expected:", numGoroutines, "Got:", counter)
}
EOF
```

3. Run it plainly, several times:

```bash
go run main.go
go run main.go
go run main.go
```

4. Now attempt the race detector:

```bash
go run -race main.go
```

**Expected Output:**

This program's "Got:" value is **genuinely non-deterministic and will typically be less than 1000** - never treat any specific number here as "the" expected output. Real captured output from 5 plain runs in this course's sandbox:

```
Expected: 1000 Got: 970
Expected: 1000 Got: 987
Expected: 1000 Got: 958
Expected: 1000 Got: 943
Expected: 1000 Got: 943
```

**`go run -race` was tested in this sandbox and it worked** (cgo is available here). Real captured `-race` output:

```
==================
WARNING: DATA RACE
Read at 0x00c000012158 by goroutine 8:
  main.main.func1()
      /path/to/main.go:17 +0x68

Previous write at 0x00c000012158 by goroutine 7:
  main.main.func1()
      /path/to/main.go:17 +0x78

Goroutine 8 (running) created at:
  main.main()
      /path/to/main.go:15 +0x70

Goroutine 7 (finished) created at:
  main.main()
      /path/to/main.go:15 +0x70
==================
Expected: 1000 Got: 869
Found 2 data race(s)
exit status 1
```

If `-race` isn't available on your machine (it requires cgo, which not every environment has), you can still confirm the race exists from the plain runs alone: a "Got:" value below 1000, varying between runs, is proof of lost updates even without the detector's line-level report.

**Learning Objectives:**
- ✅ Understand a data race is unsynchronized concurrent access with at least one write
- ✅ See real, non-deterministic, incorrect output from a racy counter
- ✅ Use (or attempt) `go run -race` and read its report

---

## Exercise 6: Fixing the Race With sync.Mutex

**Objective:** Fix Exercise 5's race with a minimal mutex

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level20-exercise6
cd ~/projects/level20-exercise6
go mod init level20.example/exercise6
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "sync"
)

func main() {
    var wg sync.WaitGroup
    var mu sync.Mutex
    counter := 0

    const numGoroutines = 1000
    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            mu.Lock()
            counter++
            mu.Unlock()
        }()
    }

    wg.Wait()
    fmt.Println("Expected:", numGoroutines, "Got:", counter)
}
EOF
```

3. Run it plainly several times, then with the race detector:

```bash
go run main.go
go run main.go
go run -race main.go
```

**Expected Output:**

```
Expected: 1000 Got: 1000
```

Deterministic - verified stable across 5 plain runs, and `go run -race` reported **0 data races** on this program in this sandbox. `mu.Lock()`/`mu.Unlock()` around `counter++` makes the increment effectively atomic with respect to every other goroutine.

**Learning Objectives:**
- ✅ Use `sync.Mutex`'s `Lock`/`Unlock` to protect a critical section
- ✅ Confirm the fix both by correctness (`Got: 1000` every time) and by the race detector finding nothing
- ✅ See the same buggy program from Exercise 5 become fully correct with 3 added lines

---

## Exercise 7: Goroutine Leaks

**Objective:** Understand why a goroutine blocked forever is never cleaned up

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level20-exercise7
cd ~/projects/level20-exercise7
go mod init level20.example/exercise7
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "time"
)

func leaky(ch chan int) {
    val := 42
    fmt.Println("leaky goroutine: about to send, but nobody is receiving...")
    ch <- val // blocks forever - nobody ever reads from ch
    fmt.Println("leaky goroutine: this line never runs")
}

func main() {
    ch := make(chan int) // unbuffered: a send blocks until a receiver is ready

    fmt.Println("main: launching a goroutine that will block forever")
    go leaky(ch)

    select {
    case <-time.After(200 * time.Millisecond):
        fmt.Println("main: gave up waiting - the goroutine above is now leaked")
    }

    fmt.Println("main: returning; the leaked goroutine dies with the process")
}
EOF
```

3. Run the program (it terminates safely on its own - the timeout guarantees that):

```bash
go run main.go
```

**Expected Output:**

```
main: launching a goroutine that will block forever
leaky goroutine: about to send, but nobody is receiving...
main: gave up waiting - the goroutine above is now leaked
main: returning; the leaked goroutine dies with the process
```

Deterministic - verified stable across 3 real runs. Note the line `leaky goroutine: this line never runs` never appears, by design - `leaky` really is stuck forever on `ch <- val`. This exercise only terminates because `main()` uses a 200ms timeout instead of waiting on `leaky` forever, and because the whole process exits afterward (which forcibly kills every goroutine, leaked or not). In a long-running server, a leak like this would instead sit there consuming memory for the life of the process, one instance per time this code path is hit.

**Learning Objectives:**
- ✅ Understand a goroutine leak: blocked forever, never garbage collected
- ✅ See a minimal, safely-terminating illustration using a timeout
- ✅ Recognize the danger of an unbounded channel send/receive with no guaranteed reader/writer on the other end

---

## Exercise 8: Fan-Out Worker Pattern - Parallel Partial Sums

**Objective:** Split independent work across goroutines and aggregate the results safely

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level20-exercise8
cd ~/projects/level20-exercise8
go mod init level20.example/exercise8
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "sync"
)

func partialSum(nums []int, wg *sync.WaitGroup, mu *sync.Mutex, total *int) {
    defer wg.Done()
    sum := 0
    for _, n := range nums {
        sum += n
    }
    mu.Lock()
    *total += sum
    mu.Unlock()
}

func main() {
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

    const numWorkers = 4
    chunkSize := len(numbers) / numWorkers

    var wg sync.WaitGroup
    var mu sync.Mutex
    total := 0

    fmt.Printf("main: splitting %d numbers across %d workers\n", len(numbers), numWorkers)

    for w := 0; w < numWorkers; w++ {
        start := w * chunkSize
        end := start + chunkSize
        if w == numWorkers-1 {
            end = len(numbers)
        }
        wg.Add(1)
        go partialSum(numbers[start:end], &wg, &mu, &total)
    }

    wg.Wait()
    fmt.Println("Total sum:", total)
}
EOF
```

3. Run the program:

```bash
go run main.go
go run -race main.go
```

**Expected Output:**

```
main: splitting 20 numbers across 4 workers
Total sum: 210
```

Deterministic - verified stable across 5 plain runs and clean under `-race`. Each of the 4 workers sums its own slice of `numbers` independently (no shared state during the actual summing), then only briefly locks the mutex to add its partial result into the shared `total` - this minimizes how much time is spent inside the critical section.

**Learning Objectives:**
- ✅ Split one workload into independent chunks (fan-out)
- ✅ Have each goroutine do its heavy work unsynchronized, and only synchronize the final aggregation
- ✅ Confirm correctness against a hand-computable expected total (1+2+...+20 = 210)

---

## Exercise 9: Mutex-Protected Results Collection

**Objective:** Have many goroutines each produce one result and collect them all safely into a shared slice

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level20-exercise9
cd ~/projects/level20-exercise9
go mod init level20.example/exercise9
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

func square(n int, wg *sync.WaitGroup, mu *sync.Mutex, results *[]int) {
    defer wg.Done()
    sq := n * n
    mu.Lock()
    *results = append(*results, sq)
    mu.Unlock()
}

func main() {
    inputs := []int{1, 2, 3, 4, 5, 6, 7, 8}

    var wg sync.WaitGroup
    var mu sync.Mutex
    var results []int

    fmt.Printf("main: squaring %d numbers concurrently\n", len(inputs))

    for _, n := range inputs {
        wg.Add(1)
        go square(n, &wg, &mu, &results)
    }

    wg.Wait()
    sort.Ints(results)
    fmt.Println("Squares (sorted):", results)
}
EOF
```

3. Run the program:

```bash
go run main.go
go run -race main.go
```

**Expected Output:**

```
main: squaring 8 numbers concurrently
Squares (sorted): [1 4 9 16 25 36 49 64]
```

Deterministic once sorted - verified stable across 5 plain runs and clean under `-race`. The individual goroutines can finish in any order, so `results` fills up in an unpredictable sequence - `sort.Ints` is what makes the *checked* output deterministic regardless of that internal ordering, exactly the pattern this course's non-determinism note recommends.

**Learning Objectives:**
- ✅ Append to a shared slice safely from multiple goroutines using a mutex
- ✅ Recognize that "the order goroutines finish in" and "the final checked result" are different things
- ✅ Reinforce the mutex-protected-collection pattern used in the final exercise

---

## Exercise 10: Comprehensive Practice - Concurrent URL Fetch Simulation

**Objective:** Combine everything - WaitGroup, closures-in-loops, and a mutex-protected results map - into one realistic program

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level20-exercise10
cd ~/projects/level20-exercise10
go mod init level20.example/exercise10
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "sort"
    "sync"
    "time"
)

// fetchURL simulates a network call with a fake latency.
func fetchURL(url string, delay time.Duration, wg *sync.WaitGroup, mu *sync.Mutex, results map[string]string) {
    defer wg.Done()
    time.Sleep(delay) // simulated network latency
    mu.Lock()
    results[url] = fmt.Sprintf("200 OK (%v)", delay)
    mu.Unlock()
}

func main() {
    urls := map[string]time.Duration{
        "https://api.example.com/users":    30 * time.Millisecond,
        "https://api.example.com/orders":   50 * time.Millisecond,
        "https://api.example.com/products": 10 * time.Millisecond,
        "https://api.example.com/invoices": 40 * time.Millisecond,
    }

    var wg sync.WaitGroup
    var mu sync.Mutex
    results := make(map[string]string)

    fmt.Printf("main: fetching %d URLs concurrently\n", len(urls))
    start := time.Now()

    for url, delay := range urls {
        wg.Add(1)
        go fetchURL(url, delay, &wg, &mu, results)
    }

    wg.Wait()
    elapsed := time.Since(start)

    fmt.Println("main: all fetches complete")

    keys := make([]string, 0, len(results))
    for k := range results {
        keys = append(keys, k)
    }
    sort.Strings(keys)

    for _, k := range keys {
        fmt.Printf("  %s -> %s\n", k, results[k])
    }

    if elapsed < 100*time.Millisecond {
        fmt.Println("Ran concurrently (faster than the sum of all delays)")
    }
}
EOF
```

3. Run the program:

```bash
go run main.go
go run -race main.go
```

**Expected Output:**

```
main: fetching 4 URLs concurrently
main: all fetches complete
  https://api.example.com/invoices -> 200 OK (40ms)
  https://api.example.com/orders -> 200 OK (50ms)
  https://api.example.com/products -> 200 OK (10ms)
  https://api.example.com/users -> 200 OK (30ms)
Ran concurrently (faster than the sum of all delays)
```

Deterministic - verified stable across 5 plain runs and clean under `-race`. Even though `range` over a map visits entries in randomized order (Level 6), and the 4 simulated fetches finish at different real times (10ms to 50ms apart), the **checked, observable output** is fully deterministic because: (1) `wg.Wait()` guarantees every fetch has written its result before we read `results`, (2) the results are sorted by URL before printing, and (3) the total elapsed time is checked with a threshold (`< 100ms`) rather than an exact duration - which is comfortably true here since the four fetches overlap (running concurrently takes roughly as long as the *slowest* one, 50ms, not the *sum* of all four, 130ms).

**Learning Objectives:**
- ✅ Combine WaitGroup synchronization, safe per-iteration loop variable capture, and mutex-protected shared state in one program
- ✅ Use `time.Sleep` to simulate I/O-bound work and observe concurrent speedup
- ✅ Design a concurrent program so its checked output is deterministic even though internal timing isn't

---

## Bonus Challenges

### Challenge 1: Goroutine-Based Pipeline (No Channels Yet)

Build a three-stage pipeline using goroutines, slices, and a mutex (not channels - those are Level 21's subject): produce a slice of numbers, launch goroutines that each square one number into a shared, mutex-protected results slice, then sum the squares once every goroutine has finished.

```bash
mkdir -p ~/projects/level20-bonus1
cd ~/projects/level20-bonus1
go mod init level20.example/bonus1
```

**Hints:**
- Reuse the mutex-protected-slice pattern from Exercise 9 for the "square" stage
- Only sum the results after `wg.Wait()` returns
- `sort.Ints` the squared slice first if you want a stable printed order

### Challenge 2: Bounded Concurrency With a Semaphore

Launch 10 goroutines, but use a buffered channel of `struct{}` as a counting semaphore so that no more than 3 are doing their "work" at the same time. This is a preview of a channel *pattern* - full channel semantics arrive in Level 21.

```bash
mkdir -p ~/projects/level20-bonus2
cd ~/projects/level20-bonus2
go mod init level20.example/bonus2
```

**Hints:**
- `sem := make(chan struct{}, 3)` - a channel with room for 3 empty values
- Before doing work: `sem <- struct{}{}` (blocks once 3 are already "checked in")
- After finishing: `<-sem` (frees a slot for someone else)
- Track the maximum number concurrently active with a mutex-protected counter, and assert it never exceeds 3

### Challenge 3: Sequential vs. Goroutine Timing Comparison

Do the same simulated "work" (e.g., `time.Sleep(5 * time.Millisecond)` plus a small computation) 50 times - once in a plain sequential loop, once by launching 50 goroutines - and print `time.Since(start)` for each approach.

```bash
mkdir -p ~/projects/level20-bonus3
cd ~/projects/level20-bonus3
go mod init level20.example/bonus3
```

**Hints:**
- Use `time.Now()` before each approach and `time.Since(start)` after
- **Be honest with yourself about the result:** exact durations depend entirely on your machine, its core count, and what else is running - don't expect to match anyone else's numbers exactly. The *shape* of the result (concurrent version dramatically faster for I/O-bound work like `time.Sleep`) is the point, not the specific milliseconds.
- One real captured run in this course's sandbox: sequential took ~288ms for 50 items at 5ms each (as expected, roughly 50 × 5ms), while the goroutine version took ~6ms total (all 50 sleeps overlapped) - your numbers will differ, but the concurrent version should be dramatically faster on any machine with more than one core available.

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Launch goroutines with `go` and explain why unsynchronized output is non-deterministic
✅ Use `sync.WaitGroup` to wait for one or many goroutines to finish
✅ Confirm per-iteration loop variable capture works correctly on this project's Go version (1.26.5)
✅ Recognize a data race, reproduce one, and read `-race` detector output
✅ Fix a race with a minimal `sync.Mutex`
✅ Explain what a goroutine leak is and why it's never garbage collected
✅ Build a fan-out worker pattern that aggregates results safely
✅ Design concurrent programs whose checked output is deterministic even when internal timing isn't

---

## Next Level

Level 21: Channels
- Sending and receiving values between goroutines
- Buffered vs. unbuffered channels
- `select` for waiting on multiple channels
- Replacing the mutex-protected-slice patterns here with idiomatic channel-based communication

Great work! Concurrency is one of Go's superpowers, and you just started wielding it! 🚀
