# Level 39: System Design - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

**A note on this level's exercises:** system design mixes two very different kinds of exercise, and this file is honest about which is which. Exercises marked **(Real, Runnable)** give you exact bash commands, a `main.go`, and a genuine **Expected Output** captured by actually running `go run main.go` in a scratch directory - verify your own run matches. Exercises marked **(Design Exercise)** have no program to run and no single correct output; instead they give you a scenario, then a **Discussion Points** section that walks through what a strong answer looks like, so you can compare your own reasoning against it.

---

## Exercise 1: LRU-ish Cache (Real, Runnable)

**Objective:** Implement a fixed-capacity, thread-safe cache with least-recently-used eviction

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level39-exercise1
cd ~/projects/level39-exercise1
go mod init level39.example/exercise1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "container/list"
    "fmt"
    "sync"
)

// entry is what each list.Element.Value holds.
type entry struct {
    key   string
    value int
}

// LRUCache is a fixed-capacity, thread-safe, least-recently-used cache.
// It combines a map (Level 10) for O(1) lookup with a doubly linked list
// (container/list) to track recency, and a sync.Mutex (Level 23) so it's
// safe under concurrent access.
type LRUCache struct {
    mu       sync.Mutex
    capacity int
    items    map[string]*list.Element
    order    *list.List // front = most recently used, back = least recently used
}

func NewLRUCache(capacity int) *LRUCache {
    return &LRUCache{
        capacity: capacity,
        items:    make(map[string]*list.Element),
        order:    list.New(),
    }
}

// Get returns the value for key and marks it as most recently used.
func (c *LRUCache) Get(key string) (int, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()

    el, ok := c.items[key]
    if !ok {
        return 0, false
    }
    c.order.MoveToFront(el)
    return el.Value.(*entry).value, true
}

// Put inserts or updates key, marking it most recently used. If the cache
// is over capacity afterward, the least-recently-used entry is evicted.
func (c *LRUCache) Put(key string, value int) {
    c.mu.Lock()
    defer c.mu.Unlock()

    if el, ok := c.items[key]; ok {
        el.Value.(*entry).value = value
        c.order.MoveToFront(el)
        return
    }

    el := c.order.PushFront(&entry{key: key, value: value})
    c.items[key] = el

    if c.order.Len() > c.capacity {
        oldest := c.order.Back()
        if oldest != nil {
            c.order.Remove(oldest)
            delete(c.items, oldest.Value.(*entry).key)
        }
    }
}

// Len reports how many entries are currently cached.
func (c *LRUCache) Len() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.order.Len()
}

func main() {
    fmt.Println("=== LRU Cache: Basic Put/Get ===")
    cache := NewLRUCache(3)
    cache.Put("a", 1)
    cache.Put("b", 2)
    cache.Put("c", 3)

    for _, k := range []string{"a", "b", "c"} {
        v, ok := cache.Get(k)
        fmt.Printf("Get(%q) = %d, found=%v\n", k, v, ok)
    }

    fmt.Println("\n=== LRU Cache: Eviction When Over Capacity ===")
    fmt.Println("Capacity is 3, cache holds a, b, c (all three were just Get, so c is most recent, a is least recent).")
    fmt.Println("Touching a again, then adding d:")
    cache.Get("a") // "a" moves to most-recently-used; "b" is now the least-recently-used
    cache.Put("d", 4)

    for _, k := range []string{"a", "b", "c", "d"} {
        v, ok := cache.Get(k)
        if ok {
            fmt.Printf("Get(%q) = %d, found=%v\n", k, v, ok)
        } else {
            fmt.Printf("Get(%q) = <evicted>, found=%v\n", k, ok)
        }
    }
    fmt.Println("Cache length:", cache.Len())

    fmt.Println("\n=== LRU Cache: Updating an Existing Key Refreshes Recency ===")
    cache2 := NewLRUCache(2)
    cache2.Put("x", 10)
    cache2.Put("y", 20)
    cache2.Put("x", 99) // update x -> most recently used again
    cache2.Put("z", 30) // over capacity -> evicts least recently used, which is now y

    _, xFound := cache2.Get("x")
    _, yFound := cache2.Get("y")
    _, zFound := cache2.Get("z")
    fmt.Printf("x found=%v, y found=%v, z found=%v\n", xFound, yFound, zFound)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== LRU Cache: Basic Put/Get ===
Get("a") = 1, found=true
Get("b") = 2, found=true
Get("c") = 3, found=true

=== LRU Cache: Eviction When Over Capacity ===
Capacity is 3, cache holds a, b, c (all three were just Get, so c is most recent, a is least recent).
Touching a again, then adding d:
Get("a") = 1, found=true
Get("b") = <evicted>, found=false
Get("c") = 3, found=true
Get("d") = 4, found=true
Cache length: 3

=== LRU Cache: Updating an Existing Key Refreshes Recency ===
x found=true, y found=false, z found=true
```

**Learning Objectives:**
- ✅ Implement a fixed-capacity cache combining a map and a doubly linked list
- ✅ Verify LRU eviction actually removes the least-recently-used entry, not just the oldest inserted
- ✅ See that touching a key with Get refreshes its recency, protecting it from eviction

---

## Exercise 2: Cache Strategy Design (Design Exercise)

**Objective:** Choose and justify a caching pattern for a given scenario

**Scenario:** You're adding caching to a product-catalog service. Product data (name, price, description) is read on every page view but only updated a few times a day by an internal admin tool. Page views outnumber updates by roughly 10,000 to 1. Product pages must never show a price that's more than a few seconds stale after an admin updates it.

**Design Considerations:**

- **Cache-aside vs write-through:** With a 10,000:1 read:write ratio, cache-aside (Section 3 of the README) is the natural fit - the cache only holds what's actually been requested, and the rare write path stays simple (write to the database, then invalidate the cache entry). Write-through would pay its extra write-path cost on every one of those rare admin updates for no real benefit, since reads are already going to be cache hits under cache-aside almost all the time.
- **Invalidation:** On every admin update, explicitly delete (or overwrite) the cached entry for that product ID, rather than relying purely on a TTL - this directly satisfies the "never stale for more than a few seconds" requirement, since the next read after an update is guaranteed to miss and refetch the fresh value. A TTL alone (e.g., 5 minutes) would occasionally leave a stale price visible for up to that TTL after an update - not acceptable here.
- **What if invalidation itself fails or is delayed?** A short TTL as a *backstop* behind explicit invalidation (e.g., invalidate on write, but also expire after 30 seconds regardless) bounds the worst case if an invalidation is ever missed - belt and suspenders, not a replacement for explicit invalidation.
- **Common Mistake 2 from the README applies directly here:** shipping the cache without deciding on an invalidation strategy at all - "it'll refresh eventually" - is exactly the failure mode this scenario's price-staleness requirement rules out.

**Learning Objectives:**
- ✅ Match a caching pattern (cache-aside vs write-through) to a read:write ratio and a freshness requirement
- ✅ Design an explicit invalidation strategy rather than relying on TTL alone
- ✅ Recognize "caching without an invalidation strategy" as a common, costly mistake

---

## Exercise 3: Token-Bucket Rate Limiter (Real, Runnable)

**Objective:** Implement and verify the token-bucket rate-limiting algorithm

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level39-exercise3
cd ~/projects/level39-exercise3
go mod init level39.example/exercise3
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

// TokenBucket implements the token-bucket rate-limiting algorithm.
// A bucket holds up to `capacity` tokens. Every request consumes one
// token. Tokens refill continuously at `refillRate` tokens per second.
// If the bucket is empty, the request is denied.
type TokenBucket struct {
    mu         sync.Mutex
    capacity   float64
    tokens     float64
    refillRate float64 // tokens added per second
    lastRefill time.Time
}

func NewTokenBucket(capacity float64, refillRate float64) *TokenBucket {
    return &TokenBucket{
        capacity:   capacity,
        tokens:     capacity, // start full
        refillRate: refillRate,
        lastRefill: time.Now(),
    }
}

// refill adds tokens based on how much time has passed since the last
// refill, capped at capacity. Must be called with mu held.
func (b *TokenBucket) refill() {
    now := time.Now()
    elapsed := now.Sub(b.lastRefill).Seconds()
    b.tokens += elapsed * b.refillRate
    if b.tokens > b.capacity {
        b.tokens = b.capacity
    }
    b.lastRefill = now
}

// Allow reports whether a request may proceed right now. If so, it
// consumes one token.
func (b *TokenBucket) Allow() bool {
    b.mu.Lock()
    defer b.mu.Unlock()

    b.refill()
    if b.tokens >= 1 {
        b.tokens--
        return true
    }
    return false
}

func main() {
    fmt.Println("=== Token Bucket: Burst Then Deny ===")
    fmt.Println("Capacity 3 tokens, refill rate 5 tokens/sec (bucket starts full)")
    bucket := NewTokenBucket(3, 5)

    for i := 1; i <= 5; i++ {
        allowed := bucket.Allow()
        fmt.Printf("Request %d: allowed=%v\n", i, allowed)
    }

    fmt.Println("\n=== Token Bucket: Waiting Lets Tokens Refill ===")
    fmt.Println("Waiting 400ms (at 5 tokens/sec, that refills ~2 tokens)...")
    time.Sleep(400 * time.Millisecond)

    for i := 1; i <= 3; i++ {
        allowed := bucket.Allow()
        fmt.Printf("Request %d after wait: allowed=%v\n", i, allowed)
    }

    fmt.Println("\n=== Token Bucket: Steady Trickle at the Refill Rate ===")
    fmt.Println("Capacity 1 token, refill rate 10 tokens/sec (1 token every 100ms)")
    steady := NewTokenBucket(1, 10)
    steady.Allow() // drain the initial full token

    for i := 1; i <= 4; i++ {
        time.Sleep(100 * time.Millisecond)
        allowed := steady.Allow()
        fmt.Printf("Request %d (after 100ms wait): allowed=%v\n", i, allowed)
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
=== Token Bucket: Burst Then Deny ===
Capacity 3 tokens, refill rate 5 tokens/sec (bucket starts full)
Request 1: allowed=true
Request 2: allowed=true
Request 3: allowed=true
Request 4: allowed=false
Request 5: allowed=false

=== Token Bucket: Waiting Lets Tokens Refill ===
Waiting 400ms (at 5 tokens/sec, that refills ~2 tokens)...
Request 1 after wait: allowed=true
Request 2 after wait: allowed=true
Request 3 after wait: allowed=false

=== Token Bucket: Steady Trickle at the Refill Rate ===
Capacity 1 token, refill rate 10 tokens/sec (1 token every 100ms)
Request 1 (after 100ms wait): allowed=true
Request 2 (after 100ms wait): allowed=true
Request 3 (after 100ms wait): allowed=true
Request 4 (after 100ms wait): allowed=true
```

**Learning Objectives:**
- ✅ Implement token-bucket refill based on real elapsed wall-clock time
- ✅ Verify a burst up to capacity is allowed, and overflow beyond it is denied
- ✅ Verify waiting lets tokens refill, and a steady request rate at the refill rate is sustained indefinitely

---

## Exercise 4: Load-Balancing Algorithm Comparison (Design Exercise + Real Round-Robin)

**Objective:** Compare load-balancing algorithms for a given traffic pattern, and implement the simplest one for real

**Scenario:** You operate 4 backend instances behind a load balancer. Two candidate workloads:

- **Workload A:** a stateless JSON API where every request costs roughly the same amount of backend work.
- **Workload B:** a report-generation endpoint where some requests take 50ms and others take 8 seconds, mixed unpredictably.

**Design Considerations:**

- **Workload A -> round-robin.** Since every request costs about the same, blindly rotating through backends (Section 5 of the README) distributes load evenly with the least mechanism possible - there's nothing round-robin's blindness to load actually costs you here.
- **Workload B -> least-connections.** Round-robin would happily send the next request to a backend that's already 7 seconds into an 8-second report, while a backend that just finished a 50ms request sits idle waiting for its next turn. Least-connections actively routes around backends that are currently tied up, which matters a lot when request cost varies this widely.
- **Where would consistent hashing fit instead?** Neither workload above needs it - it's the right tool specifically when the *same key* (a user ID, a cache key) must keep landing on the *same backend*, e.g., to keep a local, per-backend cache useful, or to shard stateful data across backends. Picking it for workload A or B would add real complexity (Bonus Challenge 2 sketches why) for no benefit, since neither workload needs key-to-backend stickiness.

3. Implement round-robin for real:

```bash
mkdir -p ~/projects/level39-exercise4
cd ~/projects/level39-exercise4
go mod init level39.example/exercise4
```

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "sync"
)

// RoundRobin cycles through a fixed list of backend addresses, handing
// out the next one on every call. It's the simplest load-balancing
// algorithm: no health awareness, no weighting, just "next in line."
type RoundRobin struct {
    mu       sync.Mutex
    backends []string
    next     int
}

func NewRoundRobin(backends []string) *RoundRobin {
    return &RoundRobin{backends: backends}
}

// Next returns the next backend address in rotation.
func (r *RoundRobin) Next() string {
    r.mu.Lock()
    defer r.mu.Unlock()

    backend := r.backends[r.next]
    r.next = (r.next + 1) % len(r.backends)
    return backend
}

func main() {
    fmt.Println("=== Round-Robin Load Balancer: 8 Requests Over 3 Backends ===")
    lb := NewRoundRobin([]string{"10.0.0.1:8080", "10.0.0.2:8080", "10.0.0.3:8080"})

    for i := 1; i <= 8; i++ {
        fmt.Printf("Request %d -> %s\n", i, lb.Next())
    }
}
EOF
```

```bash
go run main.go
```

**Expected Output:**

```
=== Round-Robin Load Balancer: 8 Requests Over 3 Backends ===
Request 1 -> 10.0.0.1:8080
Request 2 -> 10.0.0.2:8080
Request 3 -> 10.0.0.3:8080
Request 4 -> 10.0.0.1:8080
Request 5 -> 10.0.0.2:8080
Request 6 -> 10.0.0.3:8080
Request 7 -> 10.0.0.1:8080
Request 8 -> 10.0.0.2:8080
```

**Learning Objectives:**
- ✅ Match a load-balancing algorithm to whether request cost is uniform or highly variable
- ✅ Recognize when consistent hashing's key-to-backend stickiness is (and isn't) needed
- ✅ Implement and verify a real round-robin selector

---

## Exercise 5: Database Scaling Plan (Design Exercise)

**Objective:** Design a database scaling plan given growth projections and a connection-limit constraint

**Scenario:** A service currently runs 3 application instances, each with `db.SetMaxOpenConns(20)`, against a single Postgres primary configured for `max_connections = 100`. Read traffic is growing fast (10x expected over the next year) while writes stay roughly flat. Product wants to scale the application layer to 10 instances to keep up with read traffic.

**Design Considerations:**

- **The connection-limit math comes first.** 3 instances × 20 connections = 60, comfortably under 100. But scaling to 10 instances × 20 = 200 - double what the primary allows. This is Common Mistake 4 from the README happening in slow motion: scaling the application layer without revisiting the database's shared connection budget. Any scaling plan here has to address this *before* the instance count grows, not after connection errors start in production.
- **Read replicas are the correctly-targeted lever.** Since writes are flat and reads are the growth driver, adding read replicas (Section 6) - and routing read-only queries to them - scales exactly the traffic that's actually growing, without touching the primary's write path or its connection budget for writes.
- **Revised connection plan:** each of the 10 instances splits its budget - say, a smaller pool against the primary for writes (10 instances × 5 = 50 write-path connections) and a separate pool against a rotation of read replicas for reads. The primary's `max_connections` no longer has to absorb the full growth, because most of the new instances' connections land on replicas instead.
- **Why not shard instead?** Sharding (Section 6) solves write throughput and total data volume - neither of which is this scenario's actual bottleneck (writes are flat). Reaching for sharding here would add its real operational cost (cross-shard queries, rebalancing) to solve a problem this scenario doesn't have; it's a textbook case of Common Mistake 5, designing for scale that isn't the one actually arriving.

**Learning Objectives:**
- ✅ Treat a database's connection limit as a shared budget across every application instance, not a per-instance setting in isolation
- ✅ Match read replicas to a read-heavy growth pattern, rather than reaching for sharding by default
- ✅ Recognize when sharding's real costs would be paid for a problem the scenario doesn't actually have

---

## Exercise 6: Message Queue vs In-Process Channels (Design Exercise)

**Objective:** Decide whether an in-process channel suffices or a real message queue is warranted

**Scenario:** Two features to design:

- **Feature A:** when a user uploads a profile picture, resize it into three thumbnail sizes before serving it. This happens inside the same API service that received the upload.
- **Feature B:** when an order is placed, notify three separate downstream teams' systems (billing, shipping, and a third-party analytics platform) - each owned and deployed independently of the order service and of each other.

**Design Considerations:**

- **Feature A -> an in-process channel (Level 21) is enough.** The work (resizing) happens within the same process that received the upload, doesn't need to survive a crash of that process (the user can simply re-upload if the service crashes mid-resize, which is rare and low-cost), and has exactly one type of consumer. Reaching for Kafka or RabbitMQ here would mean standing up and operating a separate durable system to solve a problem a buffered channel already solves in a few lines.
- **Feature B -> a real message queue is warranted.** Three genuinely separate, independently-deployed services need to consume the same "order placed" event, each on its own schedule and without the order service knowing or caring about their individual availability. If the order service tried to call all three synchronously, a single slow or down consumer (e.g., the third-party analytics platform having an outage) would block or fail order placement entirely - exactly the coupling Section 7 argues against. A queue (named by example in the README: Kafka, RabbitMQ, SQS) durably holds the event so each consumer can process it independently, retry on failure, and even be added or removed without changing the order service at all.
- **The distinguishing question, stated generally:** does the work need to survive a process crash, and do multiple *independently deployed* consumers need the same stream of events? If either answer is yes, an in-process channel structurally cannot provide it (it dies with the process, and it only has one process's goroutines as consumers) - that's when the added operational cost of a real broker starts paying for something a channel cannot.

**Learning Objectives:**
- ✅ Distinguish "decoupling within one process" from "decoupling across independently deployed services"
- ✅ Identify durability-across-crashes and cross-service fan-out as the two signals that justify a real message queue
- ✅ Avoid reaching for a message queue's operational cost when an in-process channel already solves the actual problem

---

## Exercise 7: Circuit Breaker (Real, Runnable)

**Objective:** Implement and verify a circuit breaker that trips after repeated failures and recovers after a timeout

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level39-exercise7
cd ~/projects/level39-exercise7
go mod init level39.example/exercise7
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "errors"
    "fmt"
    "sync"
    "time"
)

type State int

const (
    Closed State = iota
    Open
    HalfOpen
)

func (s State) String() string {
    switch s {
    case Closed:
        return "CLOSED"
    case Open:
        return "OPEN"
    case HalfOpen:
        return "HALF-OPEN"
    default:
        return "UNKNOWN"
    }
}

var ErrCircuitOpen = errors.New("circuit breaker is open")

// CircuitBreaker trips to Open after `failureThreshold` consecutive
// failures, refuses all calls while Open, and after `resetTimeout` moves
// to HalfOpen to let exactly one trial call through. A successful trial
// call closes the circuit again; a failed one reopens it.
type CircuitBreaker struct {
    mu               sync.Mutex
    state            State
    failureThreshold int
    consecutiveFails int
    resetTimeout     time.Duration
    openedAt         time.Time
}

func NewCircuitBreaker(failureThreshold int, resetTimeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        state:            Closed,
        failureThreshold: failureThreshold,
        resetTimeout:     resetTimeout,
    }
}

// Call runs fn through the circuit breaker's protection.
func (cb *CircuitBreaker) Call(fn func() error) error {
    cb.mu.Lock()
    if cb.state == Open {
        if time.Since(cb.openedAt) >= cb.resetTimeout {
            cb.state = HalfOpen
            fmt.Println("  [breaker] reset timeout elapsed -> HALF-OPEN, allowing one trial call")
        } else {
            cb.mu.Unlock()
            return ErrCircuitOpen
        }
    }
    cb.mu.Unlock()

    err := fn()

    cb.mu.Lock()
    defer cb.mu.Unlock()

    if err != nil {
        cb.consecutiveFails++
        if cb.state == HalfOpen {
            // Trial call failed - back to Open, restart the timeout.
            cb.state = Open
            cb.openedAt = time.Now()
            fmt.Println("  [breaker] trial call failed -> back to OPEN")
        } else if cb.consecutiveFails >= cb.failureThreshold {
            cb.state = Open
            cb.openedAt = time.Now()
            fmt.Printf("  [breaker] %d consecutive failures reached -> OPEN\n", cb.consecutiveFails)
        }
        return err
    }

    // Success.
    if cb.state == HalfOpen {
        fmt.Println("  [breaker] trial call succeeded -> CLOSED")
    }
    cb.state = Closed
    cb.consecutiveFails = 0
    return nil
}

func (cb *CircuitBreaker) State() State {
    cb.mu.Lock()
    defer cb.mu.Unlock()
    return cb.state
}

func main() {
    fmt.Println("=== Circuit Breaker: Trips After 3 Consecutive Failures ===")
    failing := errors.New("downstream service unavailable")
    cb := NewCircuitBreaker(3, 300*time.Millisecond)

    unreliableCall := func(shouldFail bool) func() error {
        return func() error {
            if shouldFail {
                return failing
            }
            return nil
        }
    }

    // Three failing calls trip the breaker.
    for i := 1; i <= 3; i++ {
        err := cb.Call(unreliableCall(true))
        fmt.Printf("Call %d: err=%v, state=%s\n", i, err, cb.State())
    }

    fmt.Println("\n=== Circuit Breaker: Rejects Calls Immediately While Open ===")
    for i := 4; i <= 5; i++ {
        err := cb.Call(unreliableCall(false)) // wouldn't matter if it succeeds - breaker is open
        fmt.Printf("Call %d: err=%v, state=%s\n", i, err, cb.State())
    }

    fmt.Println("\n=== Circuit Breaker: Half-Open Trial After Reset Timeout ===")
    fmt.Println("Waiting 350ms for the reset timeout to elapse...")
    time.Sleep(350 * time.Millisecond)

    err := cb.Call(unreliableCall(false)) // this time the downstream call succeeds
    fmt.Printf("Call 6 (trial): err=%v, state=%s\n", err, cb.State())

    fmt.Println("\n=== Circuit Breaker: Back to Normal Operation ===")
    for i := 7; i <= 8; i++ {
        err := cb.Call(unreliableCall(false))
        fmt.Printf("Call %d: err=%v, state=%s\n", i, err, cb.State())
    }

    fmt.Println("\n=== Circuit Breaker: A Failed Trial Reopens It ===")
    cb2 := NewCircuitBreaker(2, 200*time.Millisecond)
    cb2.Call(unreliableCall(true))
    cb2.Call(unreliableCall(true))
    fmt.Println("After 2 failures, state:", cb2.State())
    time.Sleep(220 * time.Millisecond)
    err = cb2.Call(unreliableCall(true)) // trial call also fails
    fmt.Printf("Trial call also fails: err=%v, state=%s\n", err, cb2.State())
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Circuit Breaker: Trips After 3 Consecutive Failures ===
Call 1: err=downstream service unavailable, state=CLOSED
Call 2: err=downstream service unavailable, state=CLOSED
  [breaker] 3 consecutive failures reached -> OPEN
Call 3: err=downstream service unavailable, state=OPEN

=== Circuit Breaker: Rejects Calls Immediately While Open ===
Call 4: err=circuit breaker is open, state=OPEN
Call 5: err=circuit breaker is open, state=OPEN

=== Circuit Breaker: Half-Open Trial After Reset Timeout ===
Waiting 350ms for the reset timeout to elapse...
  [breaker] reset timeout elapsed -> HALF-OPEN, allowing one trial call
  [breaker] trial call succeeded -> CLOSED
Call 6 (trial): err=<nil>, state=CLOSED

=== Circuit Breaker: Back to Normal Operation ===
Call 7: err=<nil>, state=CLOSED
Call 8: err=<nil>, state=CLOSED

=== Circuit Breaker: A Failed Trial Reopens It ===
  [breaker] 2 consecutive failures reached -> OPEN
After 2 failures, state: OPEN
  [breaker] reset timeout elapsed -> HALF-OPEN, allowing one trial call
  [breaker] trial call failed -> back to OPEN
Trial call also fails: err=downstream service unavailable, state=OPEN
```

**Learning Objectives:**
- ✅ Implement the closed/open/half-open state machine
- ✅ Verify the breaker fails fast (rejects immediately, without calling the dependency) once open
- ✅ Verify a real elapsed timeout moves it to half-open, and that a trial call's outcome decides whether it closes or reopens

---

## Exercise 8: Monolith vs Microservices Tradeoff (Design Exercise)

**Objective:** Recommend an architecture for a given team and product stage, with honest tradeoffs

**Scenario:** A 4-person engineering team is building a new B2B SaaS product's first version. They've read about how companies like Netflix and Uber run hundreds of microservices, and are debating whether to start the same way.

**Design Considerations:**

- **Start as a well-layered monolith.** At 4 engineers building a *first* version, there's no independent-scaling need yet (there's no traffic yet at all) and no team-autonomy need yet (4 people can coordinate a single release train without difficulty). Splitting into microservices now means paying the full operational tax up front - multiple deploy pipelines, network calls where function calls would do, distributed tracing to debug a single request - while the team is also trying to find product-market fit, which is the worst possible time to slow feature velocity down.
- **Keep the internal boundaries clean anyway.** Using the Repository-Service-Handler pattern (Level 30) inside the monolith means the *seams* a future microservices split would need already exist as clean interfaces, even though everything currently deploys as one unit. This is the concrete version of "extraction is far easier from a cleanly-layered monolith than a tangled one" from the README's monolith-vs-microservices section.
- **What Netflix/Uber's situation is actually different about:** those are organizations with hundreds of engineers across many teams, each needing to ship independently without blocking each other, and traffic patterns where different parts of the system genuinely need very different scaling profiles. Neither condition holds for a 4-person team shipping a first version - copying the architecture without the organizational and traffic conditions that motivated it is Common Mistake 1 from the README, stated as a live scenario.
- **When would this team's advice change?** If the team grows to the point where multiple groups are regularly blocked on each other's release cadence, or a specific feature (e.g., a heavy report-generation job) develops a scaling or resource profile wildly different from the rest of the product, that's the signal - not a company size milestone or an industry trend - to extract that one piece into its own service.

**Learning Objectives:**
- ✅ Weigh operational complexity against independent scaling/deployment for a concrete team size and product stage
- ✅ Recognize when a well-known company's architecture doesn't transfer to a different team size or traffic profile
- ✅ Identify the actual signal (not a size milestone) that should trigger extracting a service from a monolith

---

## Exercise 9: Stateless Service Design (Design Exercise)

**Objective:** Redesign a stateful service to be stateless so it can scale horizontally

**Scenario:** A shopping-cart service keeps each user's in-progress cart in an in-memory map on whichever server first handled that user's session. It currently runs on a single server. Product wants to run 5 instances behind a load balancer to handle Black Friday traffic, but testing shows carts randomly "disappear" when a user's requests land on a different instance than the one holding their cart.

**Design Considerations:**

- **Diagnosis:** this is Section 2's stateful-service problem exactly - the cart lives in one server's memory, so any request that lands on a *different* server (which round-robin or least-connections load balancing will do, by design) finds no cart there. This isn't a bug in the cart logic; it's a mismatch between a stateful design and a horizontally-scaled deployment.
- **Option 1 - sticky sessions:** always route a given user to the same server. This "fixes" the symptom without removing the underlying coupling, and it undermines the load balancer's ability to spread load evenly (Section 5) - if that one user's server gets overloaded, their traffic can't be redirected elsewhere without losing their cart anyway.
- **Option 2 - a shared store (the better fix):** move cart state out of application-server memory and into a shared store (Redis, or a database table) that every instance reads from and writes to. Now any of the 5 instances can serve any request for any user's cart - the application layer becomes stateless (Section 2), which is exactly what makes horizontal scaling behind a load balancer work cleanly, matching how Level 32 favored JWTs over server-side sessions for the same underlying reason.
- **Tradeoff to acknowledge:** the shared store adds a network hop to every cart read/write that used to be an in-memory map lookup - a real latency cost, though a small one against something like Redis, and one that's a much better trade than sticky sessions' load-imbalance risk during the traffic spike this redesign exists to handle.

**Learning Objectives:**
- ✅ Diagnose a horizontal-scaling failure as a stateful-service problem
- ✅ Compare sticky sessions against a shared store as two different fixes, and explain why one is structurally better
- ✅ Connect this scenario back to the session-vs-JWT stateless tradeoff from Level 32

---

## Exercise 10: Comprehensive Design Exercise — URL Shortener Design Document (Design Exercise)

**Objective:** Produce a complete system design document for a URL shortener, tying together every concept from this level

**Instructions:** Write out a design document (in a text file, or just as your own notes) covering each section below. This mirrors the README's Section 10 case study - use it as your reference for what a strong answer looks like, but write your own version in your own words before comparing.

Cover, in order:

1. **Requirements** - functional (`POST /shorten`, `GET /{code}`) and non-functional (expected read:write ratio, latency expectations, no-collision guarantee).
2. **API shape** - the exact request/response shape for both endpoints.
3. **Storage choice** - what kind of datastore fits a `code -> long_url` lookup, and why a relational table with a unique index is enough at moderate scale.
4. **Short code generation** - compare auto-incrementing-ID-to-base62 against random-with-collision-check.
5. **Caching plan** - where a cache-aside cache (Section 3 / Exercise 1's LRU cache) fits, and why the read-heavy ratio makes this the highest-leverage lever available.
6. **Rate limiting plan** - a token bucket (Section 4 / Exercise 3) on `POST /shorten`, sized differently than on `GET /{code}`.
7. **Resilience plan** - where a circuit breaker (Section 8 / Exercise 7) belongs around the datastore.
8. **Scaling plan** - the order in which you'd pull scaling levers as traffic grows: stateless horizontal scaling behind a load balancer, then caching, then read replicas, then (only if truly necessary) sharding.

**Alternative prompt:** if you'd rather design a different system, do the same 8 sections for a **rate-limited public weather API** instead (`GET /weather?city=...`, called by many third-party developers, each issued an API key) - the same concepts apply: caching (weather data is the same for everyone querying the same city within a short window - a great cache-aside candidate), per-API-key rate limiting, and a scaling plan for a read-heavy, publicly abusable endpoint.

**Learning Objectives:**
- ✅ Apply caching, rate limiting, load balancing, database scaling, and circuit breakers together in one coherent design
- ✅ Justify each design decision against the system's actual requirements, not by default
- ✅ Produce a design document in the same shape a real system-design interview or a real architecture proposal would expect

---

## Bonus Challenges

### Challenge 1: Sliding-Window Rate Limiter (Real, Runnable)

Implement a sliding-window rate limiter as an alternative to Exercise 3's token bucket, and see how its behavior differs.

```bash
mkdir -p ~/projects/level39-bonus1
cd ~/projects/level39-bonus1
go mod init level39.example/bonus1
```

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "sync"
    "time"
)

// SlidingWindowLimiter allows at most `limit` requests within any
// rolling `window` duration. Unlike a token bucket (which allows a burst
// up to the bucket's capacity and then smooths out), a sliding window
// looks at exactly how many requests happened in the last `window` of
// real time - no bucket to pre-fill, no smoothing.
type SlidingWindowLimiter struct {
    mu         sync.Mutex
    limit      int
    window     time.Duration
    timestamps []time.Time
}

func NewSlidingWindowLimiter(limit int, window time.Duration) *SlidingWindowLimiter {
    return &SlidingWindowLimiter{
        limit:  limit,
        window: window,
    }
}

// Allow reports whether a request may proceed right now, and records it
// if so.
func (l *SlidingWindowLimiter) Allow() bool {
    l.mu.Lock()
    defer l.mu.Unlock()

    now := time.Now()
    cutoff := now.Add(-l.window)

    // Drop timestamps that have fallen outside the window.
    kept := l.timestamps[:0]
    for _, t := range l.timestamps {
        if t.After(cutoff) {
            kept = append(kept, t)
        }
    }
    l.timestamps = kept

    if len(l.timestamps) >= l.limit {
        return false
    }
    l.timestamps = append(l.timestamps, now)
    return true
}

func main() {
    fmt.Println("=== Sliding Window Limiter: Max 3 Requests Per 300ms ===")
    limiter := NewSlidingWindowLimiter(3, 300*time.Millisecond)

    for i := 1; i <= 4; i++ {
        allowed := limiter.Allow()
        fmt.Printf("Request %d: allowed=%v\n", i, allowed)
    }

    fmt.Println("\nWaiting 350ms so the window fully rolls past the first 3 requests...")
    time.Sleep(350 * time.Millisecond)

    allowed := limiter.Allow()
    fmt.Printf("Request 5 (after wait): allowed=%v\n", allowed)

    fmt.Println("\n=== Sliding Window vs Token Bucket: Partial Expiry ===")
    limiter2 := NewSlidingWindowLimiter(2, 200*time.Millisecond)
    fmt.Printf("Request A: allowed=%v\n", limiter2.Allow())
    time.Sleep(120 * time.Millisecond)
    fmt.Printf("Request B (120ms later): allowed=%v\n", limiter2.Allow())
    time.Sleep(100 * time.Millisecond)
    // At this point Request A (220ms ago) has fallen out of the 200ms window,
    // but Request B (100ms ago) is still inside it. Only A's slot freed up.
    fmt.Printf("Request C (220ms after A, 100ms after B): allowed=%v\n", limiter2.Allow())
    fmt.Printf("Request D (immediately after C): allowed=%v\n", limiter2.Allow())
}
EOF
```

```bash
go run main.go
```

**Expected Output:**

```
=== Sliding Window Limiter: Max 3 Requests Per 300ms ===
Request 1: allowed=true
Request 2: allowed=true
Request 3: allowed=true
Request 4: allowed=false

Waiting 350ms so the window fully rolls past the first 3 requests...
Request 5 (after wait): allowed=true

=== Sliding Window vs Token Bucket: Partial Expiry ===
Request A: allowed=true
Request B (120ms later): allowed=true
Request C (220ms after A, 100ms after B): allowed=true
Request D (immediately after C): allowed=false
```

**Hints:**
- Unlike the token bucket, there's no pre-filled "burst allowance" sitting ready - only the actual requests made in the last `window` count against the limit.
- Request C is allowed because only Request A (220ms old) has aged out of the 200ms window by the time C arrives - Request B (100ms old) is still counted, leaving exactly one free slot.

---

### Challenge 2: Sketch a Consistent-Hashing Ring (Design Exercise)

Sketch (on paper, or as an ASCII diagram like the README's) a consistent-hashing ring with 3 backends (A, B, C) and 5 keys hashed onto arbitrary positions around it. Then simulate removing backend B and show which keys need to move.

**Hints:**
- Each backend occupies one or more positions on the ring (real implementations use multiple virtual positions per backend to spread load evenly - look up "virtual nodes" if you want to go further).
- A key is served by the first backend found going clockwise from the key's own hash position.
- When B is removed, only the keys that were routed to B need to move (to whichever backend is now next clockwise) - every other key's backend is unaffected. Contrast this with `hash(key) % 3` becoming `hash(key) % 2`, which would reshuffle nearly every key.

---

### Challenge 3: Back-of-Envelope Capacity Estimation (Design Exercise)

Estimate the storage needs for a URL shortener holding 100 million short URLs. Show your work.

**Hints:**
- Estimate an average long URL length (a reasonable assumption: ~100 bytes) and a fixed short code length (e.g., 7 base62 characters).
- Don't forget metadata you'd realistically store per row: a creation timestamp, maybe a click counter, maybe an owner/API-key reference - estimate a reasonable total row size, not just the two URL fields.
- Multiply your total-bytes-per-row estimate by 100,000,000 and convert to GB, then sanity-check: does that number sound like "fits on a laptop," "needs a dedicated database server," or "needs sharding across multiple machines"? There's no single correct answer here - the point is showing realistic, order-of-magnitude reasoning, not hitting an exact number.

---

## What You've Learned

After completing these 10 exercises and the bonus challenges, you can:

✅ Implement and verify a real, working LRU-ish cache with eviction
✅ Implement and verify a real, working token-bucket rate limiter
✅ Implement and verify a real, working round-robin load-balancing selector
✅ Implement and verify a real, working circuit breaker that trips and recovers
✅ Choose a caching pattern (cache-aside vs write-through) and design an invalidation strategy for a given scenario
✅ Match a load-balancing algorithm to a workload's traffic pattern
✅ Build a database scaling plan that treats connection limits as a shared, system-wide budget
✅ Decide between an in-process channel and a real message queue based on durability and fan-out needs
✅ Make an honest, scenario-grounded monolith-vs-microservices recommendation
✅ Redesign a stateful service to be stateless for horizontal scaling
✅ Write a complete system design document that ties caching, rate limiting, load balancing, database scaling, and resilience together

---

## Next Level

Level 40: Real-World Projects
- Take everything from Levels 0-39 and build complete, real projects end-to-end
- Apply the Repository-Service-Handler pattern, real databases, and this level's system design thinking together
- Move from "I understand the pieces" to "I built the whole thing"

Great work! You've made the leap from writing correct code to designing systems that hold up under real conditions! 🚀
