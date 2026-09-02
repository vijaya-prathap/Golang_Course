# Level 23: Mutex, WaitGroup & Atomic - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

---

## Exercise 1: Fixing the Racy Counter With sync.Mutex

**Objective:** Reproduce Level 20's racy counter, then fix it with a mutex and confirm it's `-race`-clean

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level23-exercise1
cd ~/projects/level23-exercise1
go mod init level23.example/exercise1
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
    fmt.Println("=== Racy Version (No Synchronization) ===")
    racyCounter := 0
    var wg1 sync.WaitGroup
    const n = 1000
    for i := 0; i < n; i++ {
        wg1.Add(1)
        go func() {
            defer wg1.Done()
            racyCounter++ // data race: read-modify-write with no protection
        }()
    }
    wg1.Wait()
    fmt.Println("Expected:", n, "Got:", racyCounter)

    fmt.Println("\n=== Fixed Version (sync.Mutex) ===")
    var mu sync.Mutex
    fixedCounter := 0
    var wg2 sync.WaitGroup
    for i := 0; i < n; i++ {
        wg2.Add(1)
        go func() {
            defer wg2.Done()
            mu.Lock()
            fixedCounter++
            mu.Unlock()
        }()
    }
    wg2.Wait()
    fmt.Println("Expected:", n, "Got:", fixedCounter)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

4. Now run it with the race detector against just the mutex-protected half by running the whole program under `-race`:

```bash
go run -race main.go
```

**Expected Output (plain run - the racy count varies EVERY run, the fixed count never does):**

```
=== Racy Version (No Synchronization) ===
Expected: 1000 Got: 976

=== Fixed Version (sync.Mutex) ===
Expected: 1000 Got: 1000
```

**Expected Output (`-race` run):** the racy half reports `WARNING: DATA RACE` blocks and `Found 2 data race(s)` with a non-zero exit status; the fixed half still prints `Expected: 1000 Got: 1000` with no race warnings anywhere near it. `-race` was confirmed working in this project's sandbox (it needs cgo, which is available here) - if `go run -race` fails with a cgo-related error in your environment, the plain-run evidence above (racy counts that are always wrong, fixed counts that are always exactly 1000) already proves the fix works.

**Learning Objectives:**
- ✅ Reproduce a genuine data race and see it lose increments on every run
- ✅ Fix it with `sync.Mutex` and get an exact, deterministic result
- ✅ Confirm the fix with `go run -race`

---

## Exercise 2: The "Always Unlock" Pattern With defer

**Objective:** Protect a multi-field struct's critical section using `mu.Lock(); defer mu.Unlock()`

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level23-exercise2
cd ~/projects/level23-exercise2
go mod init level23.example/exercise2
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "sync"
)

type BankAccount struct {
    mu      sync.Mutex
    balance int
    history []string
}

func (a *BankAccount) Deposit(amount int) {
    a.mu.Lock()
    defer a.mu.Unlock() // always runs, even on a future early return/panic
    a.balance += amount
    a.history = append(a.history, fmt.Sprintf("deposit %d", amount))
}

func (a *BankAccount) Withdraw(amount int) bool {
    a.mu.Lock()
    defer a.mu.Unlock()
    if amount > a.balance {
        a.history = append(a.history, fmt.Sprintf("rejected withdrawal %d (insufficient funds)", amount))
        return false
    }
    a.balance -= amount
    a.history = append(a.history, fmt.Sprintf("withdraw %d", amount))
    return true
}

func (a *BankAccount) Balance() int {
    a.mu.Lock()
    defer a.mu.Unlock()
    return a.balance
}

func main() {
    account := &BankAccount{balance: 1000}
    var wg sync.WaitGroup

    for i := 0; i < 50; i++ {
        wg.Add(2)
        go func() {
            defer wg.Done()
            account.Deposit(10)
        }()
        go func() {
            defer wg.Done()
            account.Withdraw(5)
        }()
    }
    wg.Wait()

    fmt.Println("Final balance:", account.Balance())
    fmt.Println("Expected balance:", 1000+50*10-50*5)
    fmt.Println("History entries recorded:", len(account.history))
    fmt.Println("Expected history entries:", 100)
}
EOF
```

3. Run the program:

```bash
go run main.go
go run -race main.go
```

**Expected Output (identical for both, `-race` reports 0 races):**

```
Final balance: 1250
Expected balance: 1250
History entries recorded: 100
Expected history entries: 100
```

**Learning Objectives:**
- ✅ Protect two related fields (`balance` and `history`) with one mutex
- ✅ Use `defer mu.Unlock()` so an early `return` (the rejected-withdrawal path) still unlocks
- ✅ Confirm 100 concurrent operations against a shared struct produce an exact, `-race`-clean result

---

## Exercise 3: sync.RWMutex With Concurrent Readers and Writers

**Objective:** Allow many concurrent readers while writers still get exclusive access

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level23-exercise3
cd ~/projects/level23-exercise3
go mod init level23.example/exercise3
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "sort"
    "sync"
    "sync/atomic"
    "time"
)

type Cache struct {
    mu   sync.RWMutex
    data map[string]int
}

func (c *Cache) Get(key string) (int, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    time.Sleep(2 * time.Millisecond) // simulate read work
    v, ok := c.data[key]
    return v, ok
}

func (c *Cache) Set(key string, value int) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = value
}

func main() {
    cache := &Cache{data: map[string]int{"x": 1}}
    var wg sync.WaitGroup
    var readsCompleted atomic.Int64

    for i := 0; i < 20; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            if _, ok := cache.Get("x"); ok {
                readsCompleted.Add(1)
            }
        }()
    }

    var writeLog []string
    var logMu sync.Mutex
    for i := 0; i < 3; i++ {
        wg.Add(1)
        val := i + 2
        go func() {
            defer wg.Done()
            cache.Set("x", val)
            logMu.Lock()
            writeLog = append(writeLog, fmt.Sprintf("wrote x=%d", val))
            logMu.Unlock()
        }()
    }

    wg.Wait()
    sort.Strings(writeLog)
    fmt.Println("Reads completed:", readsCompleted.Load(), "(expected 20)")
    fmt.Println("Writes performed:", len(writeLog), "(expected 3)")
    finalVal, _ := cache.Get("x")
    fmt.Println("Final value of x is one of {2,3,4}:", finalVal == 2 || finalVal == 3 || finalVal == 4)
}
EOF
```

3. Run the program:

```bash
go run main.go
go run -race main.go
```

**Expected Output (stable across runs, `-race` reports 0 races):**

```
Reads completed: 20 (expected 20)
Writes performed: 3 (expected 3)
Final value of x is one of {2,3,4}: true
```

**Learning Objectives:**
- ✅ Use `RLock()`/`RUnlock()` to let 20 readers run concurrently
- ✅ Use `Lock()`/`Unlock()` for writers, which still get exclusive access
- ✅ See that the final value is nondeterministic (whichever writer ran last) but always one of the valid values

---

## Exercise 4: WaitGroup - Add() Timing Matters

**Objective:** Compare calling `Add()` before launching a goroutine vs. from inside it

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level23-exercise4
cd ~/projects/level23-exercise4
go mod init level23.example/exercise4
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "sync"
    "time"
)

func main() {
    fmt.Println("=== Correct: Add() called BEFORE launching each goroutine ===")
    var wg sync.WaitGroup
    var mu sync.Mutex
    correctCount := 0
    for i := 0; i < 5; i++ {
        wg.Add(1) // called from the loop, before `go`
        go func(id int) {
            defer wg.Done()
            mu.Lock()
            correctCount++
            mu.Unlock()
        }(i)
    }
    wg.Wait()
    fmt.Println("completed:", correctCount, "(always 5 - deterministic)")

    fmt.Println("\n=== Buggy: Add() called INSIDE the goroutine (races with Wait) ===")
    var wg2 sync.WaitGroup
    var mu2 sync.Mutex
    buggyCount := 0
    for i := 0; i < 5; i++ {
        go func(id int) {
            wg2.Add(1) // WRONG: main's Wait() below might already be checking the counter
            defer wg2.Done()
            time.Sleep(20 * time.Millisecond)
            mu2.Lock()
            buggyCount++
            mu2.Unlock()
        }(i)
    }
    wg2.Wait()
    fmt.Println("completed when Wait() returned:", buggyCount, "(NOT reliably 5 - this is the bug)")
}
EOF
```

3. Run the program 5 times in a row and compare:

```bash
for i in 1 2 3 4 5; do go run main.go; done
```

**Expected Output (the "correct" half is always 5; the "buggy" half is unreliable - verified across 13 total runs in this project's sandbox, EVERY run printed 0 for the buggy half):**

```
=== Correct: Add() called BEFORE launching each goroutine ===
completed: 5 (always 5 - deterministic)

=== Buggy: Add() called INSIDE the goroutine (races with Wait) ===
completed when Wait() returned: 0 (NOT reliably 5 - this is the bug)
```

The exact number in the buggy half is scheduler-dependent - it is not guaranteed to always be 0 on your machine (a separate test of the same pattern occasionally got 5 when a goroutine's `Add()` won the race), but it is never *reliably* 5 the way the correct version is.

**Learning Objectives:**
- ✅ See that `Add()` before `go` produces a deterministic, always-correct wait
- ✅ See that `Add()` inside the goroutine races with `Wait()` and can let `Wait()` return too early
- ✅ Understand that the danger is unreliability itself, not one specific wrong number

---

## Exercise 5: The WaitGroup-Passed-By-Value Bug

**Objective:** Trigger and understand a real deadlock caused by copying a WaitGroup

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level23-exercise5
cd ~/projects/level23-exercise5
go mod init level23.example/exercise5
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "sync"
)

// BUG: wg is a sync.WaitGroup VALUE parameter, so every call gets its own COPY.
// Done() decrements the copy inside the function, never the caller's original -
// so the caller's wg.Wait() below can never see the counter reach zero.
func processItem(wg sync.WaitGroup, item string) {
    defer wg.Done()
    fmt.Println("processing:", item)
}

func main() {
    var wg sync.WaitGroup
    wg.Add(1)
    go processItem(wg, "invoice-42") // wg is COPIED into processItem here
    wg.Wait()                        // blocks forever - the ORIGINAL counter is still 1
    fmt.Println("all done")
}
EOF
```

3. First, run `go vet` - it catches this mistake WITHOUT even running the program:

```bash
go vet .
```

4. Now run the program anyway (it deliberately hangs, then Go's own deadlock detector kills it - this is expected and safe):

```bash
go run main.go
```

**Expected Output (`go vet`):**

```
main.go:11:21: processItem passes lock by value: sync.WaitGroup contains sync.noCopy
main.go:19:17: call of processItem copies lock value: sync.WaitGroup contains sync.noCopy
```

**Expected Output (`go run main.go` - captured verbatim from 3 real runs, identical every time except memory addresses):**

```
processing: invoice-42
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [sync.WaitGroup.Wait]:
sync.runtime_SemacquireWaitGroup(0x1049cbad0?, 0x18?)
	/Users/macbookpro/.local/go/src/runtime/sema.go:114 +0x38
sync.(*WaitGroup).Wait(0x6a2e65b88020)
	/Users/macbookpro/.local/go/src/sync/waitgroup.go:206 +0xa8
main.main()
	/private/tmp/level23-verify/ex5/main.go:20 +0x74
exit status 2
```

`processItem` finishes and prints normally - it just decrements a copy nobody else can see. `main`'s own `wg` never reaches zero, so `wg.Wait()` blocks forever; once `processItem`'s goroutine exits, `main` is the only goroutine left and it's permanently stuck, so the Go runtime detects that **all** goroutines are asleep and crashes the process itself rather than hanging silently.

**Learning Objectives:**
- ✅ See `go vet` catch a by-value WaitGroup mistake statically
- ✅ Understand exactly why a copied WaitGroup's `Done()` doesn't affect the original
- ✅ Recognize Go's real `fatal error: all goroutines are asleep - deadlock!` output and what it means
- ✅ Confirm the fix is always to pass `*sync.WaitGroup` (a pointer)

---

## Exercise 6: A Typed Atomic Counter

**Objective:** Use `atomic.Int64` for an exact, lock-free counter under heavy concurrency

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level23-exercise6
cd ~/projects/level23-exercise6
go mod init level23.example/exercise6
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

func main() {
    var wg sync.WaitGroup
    var hits atomic.Int64

    const numGoroutines = 5000
    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            hits.Add(1)
        }()
    }
    wg.Wait()

    got := hits.Load()
    fmt.Println("Expected:", numGoroutines, "Got:", got)
    if got == numGoroutines {
        fmt.Println("MATCH: atomic counter is exact under concurrent increments")
    } else {
        fmt.Println("MISMATCH (should never happen with atomic.Int64.Add)")
    }
}
EOF
```

3. Run the program:

```bash
go run main.go
go run -race main.go
```

**Expected Output (stable across 5 plain runs AND a `-race` run):**

```
Expected: 5000 Got: 5000
MATCH: atomic counter is exact under concurrent increments
```

**Learning Objectives:**
- ✅ Use `atomic.Int64`'s `Add`/`Load` methods for a lock-free counter
- ✅ Confirm the count is exact across 5,000 concurrent increments, with no mutex anywhere
- ✅ Confirm it's `-race`-clean

---

## Exercise 7: Atomics vs Mutex - Choosing the Right Tool

**Objective:** Build one type that correctly uses atomics and one that correctly needs a mutex

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level23-exercise7
cd ~/projects/level23-exercise7
go mod init level23.example/exercise7
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

// Good fit for atomic: one independent counter, nothing else to keep in sync.
type Stats struct {
    requests atomic.Int64
    errors   atomic.Int64
}

// Good fit for mutex: two fields (balance and a running total of fees) that
// must always be updated TOGETHER, as one consistent unit.
type Ledger struct {
    mu            sync.Mutex
    balance       int
    totalFeesPaid int
}

func (l *Ledger) ChargeFee(fee int) {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.balance -= fee
    l.totalFeesPaid += fee // both fields move together - can't be two separate atomics safely
}

func main() {
    var stats Stats
    var wg sync.WaitGroup

    for i := 0; i < 200; i++ {
        wg.Add(1)
        isError := i%10 == 0
        go func() {
            defer wg.Done()
            stats.requests.Add(1)
            if isError {
                stats.errors.Add(1)
            }
        }()
    }
    wg.Wait()
    fmt.Println("=== atomic: independent counters ===")
    fmt.Println("requests:", stats.requests.Load(), "errors:", stats.errors.Load())

    ledger := &Ledger{balance: 10000}
    var wg2 sync.WaitGroup
    for i := 0; i < 200; i++ {
        wg2.Add(1)
        go func() {
            defer wg2.Done()
            ledger.ChargeFee(1)
        }()
    }
    wg2.Wait()
    fmt.Println("\n=== mutex: fields that must move together ===")
    fmt.Println("balance:", ledger.balance, "(expected 9800)")
    fmt.Println("totalFeesPaid:", ledger.totalFeesPaid, "(expected 200)")
    fmt.Println("balance + totalFeesPaid + net change matches original:",
        ledger.balance+ledger.totalFeesPaid == 10000)
}
EOF
```

3. Run the program:

```bash
go run main.go
go run -race main.go
```

**Expected Output (stable across 3 plain runs and a `-race` run):**

```
=== atomic: independent counters ===
requests: 200 errors: 20

=== mutex: fields that must move together ===
balance: 9800 (expected 9800)
totalFeesPaid: 200 (expected 200)
balance + totalFeesPaid + net change matches original: true
```

**Learning Objectives:**
- ✅ Use independent atomics for two counters (`requests`, `errors`) that never need to be consistent with each other
- ✅ Use one mutex to keep two related fields (`balance`, `totalFeesPaid`) consistent as a unit
- ✅ Explain why two separate atomics could not safely replace the mutex here

---

## Exercise 8: sync.Once Under Concurrent Calls

**Objective:** Verify a function runs exactly once even when 100 goroutines call it at once

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level23-exercise8
cd ~/projects/level23-exercise8
go mod init level23.example/exercise8
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

type Config struct {
    value string
}

var (
    once      sync.Once
    cfg       *Config
    initCalls atomic.Int32
)

func loadConfig() *Config {
    once.Do(func() {
        initCalls.Add(1)
        fmt.Println("loading configuration from disk (expensive, should happen once)...")
        cfg = &Config{value: "production"}
    })
    return cfg
}

func main() {
    var wg sync.WaitGroup
    results := make([]*Config, 100)

    for i := 0; i < 100; i++ {
        wg.Add(1)
        idx := i
        go func() {
            defer wg.Done()
            results[idx] = loadConfig()
        }()
    }
    wg.Wait()

    fmt.Println("initCalls (actual runs of the init function):", initCalls.Load(), "(expected 1)")

    allSamePointer := true
    for _, r := range results {
        if r != results[0] {
            allSamePointer = false
        }
    }
    fmt.Println("all 100 goroutines got the SAME *Config instance:", allSamePointer)
}
EOF
```

3. Run the program:

```bash
go run main.go
go run -race main.go
```

**Expected Output (stable across 5 plain runs and a `-race` run):**

```
loading configuration from disk (expensive, should happen once)...
initCalls (actual runs of the init function): 1 (expected 1)
all 100 goroutines got the SAME *Config instance: true
```

**Learning Objectives:**
- ✅ Use `sync.Once.Do` to guarantee one-time initialization
- ✅ Confirm the init function runs exactly once out of 100 concurrent callers
- ✅ Confirm every caller receives the identical initialized instance

---

## Exercise 9: A Real Lock-Ordering Deadlock

**Objective:** Trigger and read a genuine deadlock caused by acquiring two mutexes in opposite order

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level23-exercise9
cd ~/projects/level23-exercise9
go mod init level23.example/exercise9
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "sync"
    "time"
)

var mu1, mu2 sync.Mutex

// transferAtoB locks mu1 then mu2.
func transferAtoB(wg *sync.WaitGroup) {
    defer wg.Done()
    mu1.Lock()
    fmt.Println("transferAtoB: locked mu1")
    time.Sleep(100 * time.Millisecond)
    fmt.Println("transferAtoB: waiting for mu2")
    mu2.Lock()
    fmt.Println("transferAtoB: locked mu2")
    mu2.Unlock()
    mu1.Unlock()
}

// transferBtoA locks mu2 then mu1 - OPPOSITE ORDER. This is the bug.
func transferBtoA(wg *sync.WaitGroup) {
    defer wg.Done()
    mu2.Lock()
    fmt.Println("transferBtoA: locked mu2")
    time.Sleep(100 * time.Millisecond)
    fmt.Println("transferBtoA: waiting for mu1")
    mu1.Lock()
    fmt.Println("transferBtoA: locked mu1")
    mu1.Unlock()
    mu2.Unlock()
}

func main() {
    var wg sync.WaitGroup
    wg.Add(2)
    go transferAtoB(&wg)
    go transferBtoA(&wg)
    wg.Wait()
    fmt.Println("both transfers completed")
}
EOF
```

3. Run the program (it deliberately deadlocks - Go's runtime detector will catch it and crash the process automatically within about 100ms, this is expected and safe):

```bash
go run main.go
```

**Expected Output (captured verbatim from a real run in this project's sandbox - confirmed across 3 runs, always deadlocking):**

```
transferBtoA: locked mu2
transferAtoB: locked mu1
transferAtoB: waiting for mu2
transferBtoA: waiting for mu1
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [sync.WaitGroup.Wait]:
sync.runtime_SemacquireWaitGroup(0x103033ad0?, 0x80?)
	/Users/macbookpro/.local/go/src/runtime/sema.go:114 +0x38
sync.(*WaitGroup).Wait(0x3ef5046f20f0)
	/Users/macbookpro/.local/go/src/sync/waitgroup.go:206 +0xa8
main.main()
	/private/tmp/level23-verify/ex9/main.go:42 +0xb8

goroutine 7 [sync.Mutex.Lock]:
internal/sync.runtime_SemacquireMutex(0x3ef5046e2040?, 0x0?, 0x1e?)
	/Users/macbookpro/.local/go/src/runtime/sema.go:95 +0x28
internal/sync.(*Mutex).lockSlow(0x103194f58)
	/Users/macbookpro/.local/go/src/internal/sync/mutex.go:149 +0x170
internal/sync.(*Mutex).Lock(...)
	/Users/macbookpro/.local/go/src/internal/sync/mutex.go:70
sync.(*Mutex).Lock(...)
	/Users/macbookpro/.local/go/src/sync/mutex.go:46
main.transferAtoB(0x0?)
	/private/tmp/level23-verify/ex9/main.go:18 +0x184
created by main.main in goroutine 1
	/private/tmp/level23-verify/ex9/main.go:40 +0x70

goroutine 8 [sync.Mutex.Lock]:
internal/sync.runtime_SemacquireMutex(0x3ef5046e2040?, 0xa0?, 0x1e?)
	/Users/macbookpro/.local/go/src/runtime/sema.go:95 +0x28
internal/sync.(*Mutex).lockSlow(0x103194f50)
	/Users/macbookpro/.local/go/src/internal/sync/mutex.go:149 +0x170
internal/sync.(*Mutex).Lock(...)
	/Users/macbookpro/.local/go/src/internal/sync/mutex.go:70
sync.(*Mutex).Lock(...)
	/Users/macbookpro/.local/go/src/sync/mutex.go:46
main.transferBtoA(0x0?)
	/private/tmp/level23-verify/ex9/main.go:31 +0x184
created by main.main in goroutine 1
	/private/tmp/level23-verify/ex9/main.go:41 +0xb0
exit status 2
```

The exact interleaving of the first four "locked"/"waiting" print lines can vary between runs (which goroutine grabs its first lock first is a race), but the outcome is always the same: both workers end up waiting on a lock the other holds, and the runtime detects every goroutine (including `main`) is permanently blocked.

**Learning Objectives:**
- ✅ Construct a classic two-mutex, opposite-order deadlock
- ✅ Read a real goroutine dump and identify which goroutine is stuck where
- ✅ Understand why Go's runtime can prove this is a true deadlock (not just a slow program) and crash on it
- ✅ Know the fix: always acquire shared locks in the same global order

---

## Exercise 10: Comprehensive Practice - Thread-Safe Cache/Counter Service

**Objective:** Combine a mutex-protected map with atomic hit/miss counters in one service

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level23-exercise10
cd ~/projects/level23-exercise10
go mod init level23.example/exercise10
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

// CacheService is a thread-safe in-memory cache: a mutex protects the map
// itself (a single compound invariant: presence + value), while atomic
// counters track hits/misses lock-free since those are simple, independent tallies.
type CacheService struct {
    mu   sync.Mutex
    data map[string]int

    hits   atomic.Int64
    misses atomic.Int64
}

func NewCacheService() *CacheService {
    return &CacheService{data: make(map[string]int)}
}

func (c *CacheService) Set(key string, value int) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = value
}

func (c *CacheService) Get(key string) (int, bool) {
    c.mu.Lock()
    v, ok := c.data[key]
    c.mu.Unlock()

    if ok {
        c.hits.Add(1)
    } else {
        c.misses.Add(1)
    }
    return v, ok
}

func (c *CacheService) Stats() (hits, misses int64) {
    return c.hits.Load(), c.misses.Load()
}

func main() {
    cache := NewCacheService()
    cache.Set("alpha", 1)
    cache.Set("beta", 2)

    var wg sync.WaitGroup
    keys := []string{"alpha", "beta", "gamma", "delta"} // gamma/delta don't exist -> misses

    const workers = 200
    for i := 0; i < workers; i++ {
        wg.Add(1)
        key := keys[i%len(keys)]
        go func() {
            defer wg.Done()
            cache.Get(key)
        }()
    }
    wg.Wait()

    hits, misses := cache.Stats()
    fmt.Println("hits:", hits, "(expected 100)")
    fmt.Println("misses:", misses, "(expected 100)")
    fmt.Println("total lookups:", hits+misses, "(expected", workers, ")")
}
EOF
```

3. Run the program:

```bash
go run main.go
go run -race main.go
```

**Expected Output (stable across 3 plain runs and a `-race` run):**

```
hits: 100 (expected 100)
misses: 100 (expected 100)
total lookups: 200 (expected 200 )
```

**Learning Objectives:**
- ✅ Combine a mutex-protected map (the compound "does this key exist, and what's its value" invariant) with independent atomic counters (simple tallies) in one type
- ✅ Recognize why the counters don't need the same mutex as the map - they're independent of it
- ✅ Confirm exact, `-race`-clean results under 200 concurrent lookups

---

## Bonus Challenges

### Challenge 1: RWMutex vs Mutex - An Informal Benchmark

Build two versions of the same tiny read-heavy cache - one backed by `sync.RWMutex`, one by plain `sync.Mutex` - and time 50 concurrent goroutines each doing 2,000 reads against each version.

```bash
mkdir -p ~/projects/level23-bonus1
cd ~/projects/level23-bonus1
go mod init level23.example/bonus1
```

**Hints:**
- Use `time.Now()` / `time.Since()` around a `sync.WaitGroup`-bounded batch of reader goroutines
- **Note: exact timings vary by machine and load.** When this was actually run in this project's sandbox, plain `Mutex` (≈13.6ms) was measured *faster* than `RWMutex` (≈20ms) for 100,000 trivial map-lookup reads - the lesson is that `RWMutex`'s bookkeeping overhead isn't automatically worth it for cheap critical sections, not that one primitive is universally faster
- Print both durations and compare them on your own machine - don't expect to reproduce this project's exact numbers

### Challenge 2: A sync.Once-Based Singleton

Build a `GetDatabase()` function that lazily creates one `*Database` instance no matter how many goroutines call it concurrently, and confirm every caller gets the identical pointer.

```bash
mkdir -p ~/projects/level23-bonus2
cd ~/projects/level23-bonus2
go mod init level23.example/bonus2
```

**Hints:**
- A package-level `var once sync.Once` and `var instance *Database` is the classic shape
- Launch 20 goroutines all calling `GetDatabase()` and compare every result against the first one for pointer equality
- The "creating the one and only Database instance..." message should print exactly once

### Challenge 3: An Atomic-Based Spinlock (Learning Exercise Only)

Build a `SpinLock` type using `atomic.Bool.CompareAndSwap` in a loop instead of `sync.Mutex`, and confirm it correctly protects a counter under concurrent increments.

```bash
mkdir -p ~/projects/level23-bonus3
cd ~/projects/level23-bonus3
go mod init level23.example/bonus3
```

**Hints:**
- `Lock()` should loop calling `CompareAndSwap(false, true)` until it succeeds; call `runtime.Gosched()` in the loop so it yields instead of hammering the CPU
- `Unlock()` should just `Store(false)`
- **Important:** this is a learning exercise to understand what `CompareAndSwap` can build, not a recommendation - real code should almost always use `sync.Mutex`, which lets the OS/runtime park a blocked goroutine instead of burning CPU cycles spinning

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Fix a Level-20-style racy counter with `sync.Mutex` and confirm it's `-race`-clean
✅ Use `mu.Lock(); defer mu.Unlock()` to protect a multi-field critical section on every code path
✅ Use `sync.RWMutex` to let readers run concurrently while writers stay exclusive
✅ Recognize the WaitGroup `Add()`-timing bug and the WaitGroup-by-value bug - and reproduce a real captured deadlock from the latter
✅ Use typed `sync/atomic` types for exact, lock-free counters
✅ Choose between an atomic and a mutex based on how many related pieces of state must move together
✅ Guarantee one-time initialization under concurrency with `sync.Once`
✅ Reproduce and read a genuine lock-ordering deadlock from Go's own runtime detector
✅ Combine a mutex-protected map with atomic counters in one thread-safe service

---

## Next Level

Level 24: Context
- `context.Context` for cancellation, deadlines, and timeouts
- Passing request-scoped values safely
- Building on this level's synchronization primitives to control goroutine lifetimes explicitly

Great work! You've mastered Go's classic concurrency toolkit! 🚀
