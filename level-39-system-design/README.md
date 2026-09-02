# Level 39: System Design - Complete Guide

## Introduction

Welcome to Level 39! Level 38 (Production Debugging) taught you how to figure out why a *running* system is misbehaving - reading a profile, chasing a goroutine leak, making sense of a panic in production. This level asks a different question, earlier in the timeline: **how do you design a system so it holds up once it's running at scale in the first place?**

Every level before this one has mostly asked "does this function work?" - does the loop terminate correctly, does the handler return the right JSON, does the mutex prevent the race. Those are still the right questions for a single component. System design is the mindset shift that happens when you zoom out from "does this function work" to "will this *system* hold up when 10,000 people hit it at once, and does it fail gracefully when a dependency goes down at 3 AM?" A single Go program can be flawlessly correct and still fall over under load, or turn one slow database into a full outage, simply because nobody thought about the shape of the system around that program.

**Be upfront about what this level is:** most of the topics here are architectural reasoning, not syntax. There's no compiler error for "your caching strategy has no invalidation plan" or "you sharded your database by the wrong key." Where a concept genuinely reduces to runnable Go code - an LRU-ish cache, a token-bucket rate limiter, a circuit breaker, a round-robin selector - this guide writes that code for real, runs it, and shows you the actual captured output. Where a concept is a judgment call about tradeoffs - how to shard a database, which load balancing algorithm to pick, whether to reach for Kafka - this guide says so plainly and walks through the reasoning with diagrams instead of pretending there's a program to run.

---

## Table of Contents

1. [What System Design Is (and the Mindset Shift)](#what-system-design-is-and-the-mindset-shift)
2. [Scalability Fundamentals](#scalability-fundamentals)
3. [Caching](#caching)
4. [Rate Limiting](#rate-limiting)
5. [Load Balancing](#load-balancing)
6. [Database Scaling](#database-scaling)
7. [Message Queues and Asynchronous Processing](#message-queues-and-asynchronous-processing)
8. [Circuit Breakers](#circuit-breakers)
9. [Monolith vs Microservices](#monolith-vs-microservices)
10. [Case Study: Designing a URL Shortener](#case-study-designing-a-url-shortener)
11. [Best Practices](#best-practices)
12. [Common Mistakes](#common-mistakes)

---

## What System Design Is (and the Mindset Shift)

Every level up through Level 38 lived inside the boundary of "one program, one process." Even Level 27-30's HTTP APIs and Level 35-36's containers were about getting *a* service built and packaged correctly. System design lives one level up: it's about how multiple instances of that service, a database, a cache, and whatever else sits between your code and your users **behave together** under real conditions - real traffic, real failures, real growth.

Three questions define the mindset shift:

1. **"Does it work?" becomes "does it work under load?"** A handler that correctly returns a user's profile in 5ms with one request per second might time out or fall over at 5,000 requests per second - not because the code is wrong, but because a database connection pool, a downstream API, or a single CPU core becomes the bottleneck.

2. **"Does it work?" becomes "what happens when a dependency fails?"** Your Go program might be perfect, but if it calls a payment API that's down, does your *entire* service hang waiting for a timeout, cascading the outage to users who weren't even touching payments? System design is largely about designing the failure modes on purpose, rather than discovering them during an incident.

3. **"Does it work?" becomes "does it work *tomorrow*, at 10x the traffic, without a rewrite?"** This doesn't mean over-engineering for scale you don't have (Section 11 pushes back on that directly) - it means understanding which parts of a design bend easily under growth and which parts will need to be replaced, so growth doesn't arrive as a surprise.

None of this replaces what you already know. A well-designed system is still built out of correctly-written functions, safely-shared state (Level 23), and well-structured layers (Level 30). System design is the layer of reasoning *around* those components - and it's the last major concept this course covers before Level 40 (Real-World Projects) asks you to actually build something end-to-end using everything from Levels 0-39.

---

## Scalability Fundamentals

**Scaling** means handling more load - more users, more requests, more data - without the system falling over. There are two directions to do it in.

### Vertical Scaling ("Scale Up")

Give the existing machine more resources: more CPU cores, more RAM, a faster disk.

```
Before:                          After:
┌─────────────────┐              ┌─────────────────────┐
│   1 server       │             │   1 (bigger) server  │
│   4 CPU, 8GB RAM │   ───►      │   32 CPU, 128GB RAM   │
└─────────────────┘              └─────────────────────┘
```

- **Pros:** simple - no code changes, no coordination between instances, no distributed-systems complexity.
- **Cons:** has a hard ceiling (there's a biggest machine you can buy), it's a single point of failure (that one machine going down takes everything with it), and bigger machines get disproportionately expensive.

### Horizontal Scaling ("Scale Out")

Add *more* machines running the same service, and spread load across them (Section 5 covers how).

```
Before:                          After:
┌─────────────┐                  ┌─────────┐ ┌─────────┐ ┌─────────┐
│  1 server    │    ───►         │server 1 │ │server 2 │ │server 3 │
└─────────────┘                  └─────────┘ └─────────┘ └─────────┘
                                        ▲          ▲          ▲
                                        └──────────┴──────────┘
                                          load balancer in front
```

- **Pros:** no theoretical ceiling (add another box), and no single machine failing takes down the whole service.
- **Cons:** real added complexity - now you need a load balancer, and every server needs to agree on shared state (or not need to at all - see below).

### Stateless vs Stateful Services

This is the detail that decides whether horizontal scaling is easy or painful, and it's the same fork in the road Level 31 (Dependency Injection) and Level 32 (Authentication & JWT) already put in front of you when comparing session-based and token-based auth.

- A **stateful** service keeps something in server memory between requests from the same client - a session, an in-progress upload, a WebSocket's internal buffer. If server 2 doesn't have the session that server 1 created, the request breaks. That forces either **sticky sessions** (always route the same client to the same server - which undermines load balancing) or a **shared session store** (Redis, a database) that every server reads from.
- A **stateless** service keeps nothing about a specific client in server memory. Every request carries everything needed to handle it - most commonly a self-contained JWT (Level 32) instead of a server-side session lookup. *Any* server can handle *any* request, which is exactly what horizontal scaling wants: spin up server 4, point the load balancer at it, done - no session data to replicate, no sticky routing required.

This is precisely why Level 32 flagged JWTs as the more horizontally-scalable choice over server-side sessions: statelessness isn't a JWT implementation detail, it's a system design property that JWTs happen to provide.

```
STATEFUL (sessions in server memory)       STATELESS (JWT / self-contained)

  client ──► server 1 (has session)          client ──► any server (verifies
                                                          JWT signature, no lookup)
  client ──► server 2 (NO session!) ✗
  → breaks, or needs sticky routing,          Any server can handle any request.
    or a shared session store                 Horizontal scaling "just works."
```

Not everything can be made stateless (a database has to hold state *somewhere*), but pushing state out of your application servers and into a dedicated store (a database, a cache, an object store) is what lets the application layer itself scale horizontally without drama.

---

## Caching

A cache stores a copy of data that's expensive to (re)compute or (re)fetch, so future requests for the same data can be served fast. Caching is one of the highest-leverage tools in system design - and one of the easiest to get wrong (Section 12 covers the classic failure mode).

### Cache-Aside vs Write-Through (Conceptual)

These are the two most common patterns for keeping a cache and its source of truth (usually a database) in sync. Both are pure design reasoning here - there's no single "output" a diagram like this produces, just a description of the request flow.

**Cache-aside (a.k.a. lazy loading):** the application checks the cache first; on a miss, it reads from the database and populates the cache itself.

```
READ:
  app -> cache: got it?
       cache MISS
  app -> database: fetch it
  app -> cache: store it (so next read is a HIT)
  app -> caller: return it

WRITE:
  app -> database: write it
  app -> cache: invalidate (delete) the stale entry
```

- Cache only holds what's actually been requested (no wasted space on unread data).
- A cache outage just means slower reads (everything falls through to the database) - not broken writes.
- Every miss pays a full round trip before the cache can help.

**Write-through:** every write goes through the cache, which writes to the database on the caller's behalf, so the cache is always up to date.

```
WRITE:
  app -> cache: write it
  cache -> database: write it (synchronously)
  cache -> app: acknowledge

READ:
  app -> cache: got it? (almost always a HIT, since writes kept it current)
```

- Reads are consistently fast, because the cache is never stale after a write.
- Every write now costs two operations (cache + database) instead of one, and a cache outage can block writes entirely depending on how it's wired.

Most real systems reach for cache-aside by default (Redis-in-front-of-Postgres is the classic shape) and use write-through only when read-after-write consistency really matters and the extra write latency is acceptable.

### A Real, Runnable LRU-ish Cache

This part **is** real, runnable Go code - written and executed in a scratch directory for this level, with genuine captured output below. It bridges Level 10's maps (for O(1) lookup) and Level 23's mutexes (for safe concurrent access), and implements a fixed-capacity cache with **least-recently-used (LRU) eviction**: when the cache is full, the item that hasn't been touched in the longest time is thrown out to make room.

```go
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
```

**Actual captured output (`go run main.go`, in `/private/tmp/level39-verify/cache/`):**

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

Walking through the eviction: capacity is 3, and `a`, `b`, `c` were `Put` in that order, then all three were `Get` (so `a`, `b`, `c` are touched most-to-least-recently in that read order - `c` most recent, `a` least). Touching `a` again with `Get("a")` moves it back to most-recently-used, which makes `b` the new least-recently-used entry. `Put("d", 4)` pushes the cache to 4 entries, over its capacity of 3, so the least-recently-used entry - `b` - gets evicted. That's exactly what the captured output shows: `b` comes back `<evicted>`, everything else survives.

---

## Rate Limiting

Rate limiting caps how many requests a client (or the system as a whole) can make in a given time window, protecting a service from being overwhelmed - whether by a genuine traffic spike, a buggy retry loop, or a deliberate abuse attempt.

### The Token-Bucket Algorithm

Picture a bucket that holds up to `capacity` tokens. Tokens refill continuously at a fixed rate. Every request must take one token to proceed; if the bucket is empty, the request is denied (or queued, depending on the design).

```
   refill rate: R tokens/sec
        │
        ▼
   ┌─────────┐
   │ o o o   │  <- tokens (capacity = 3, currently holds 2)
   │ o       │
   └─────────┘
        │
        ▼
   each request takes 1 token to proceed
   empty bucket -> request denied
```

This naturally allows **bursts** up to the bucket's capacity (if it's been idle and full, several requests can go through back-to-back), while still capping the *sustained* rate to the refill rate over time - a deliberately different behavior from the sliding-window approach in this level's bonus challenges, which allows no burst beyond the raw count limit.

### A Real, Runnable Token-Bucket Rate Limiter

Bridging Level 22's timers (`time.After`, timeouts) with Level 23's mutexes, here's a genuine token bucket, run for real with actual elapsed wall-clock time (not hand-computed):

```go
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
```

**Actual captured output (`go run main.go`, run twice for stability, both runs identical):**

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

The first 3 requests drain the full bucket (capacity 3) and succeed; requests 4 and 5 immediately after find it empty and are denied. After waiting 400ms at a 5-tokens/sec refill rate (~2 tokens' worth), 2 more requests succeed before the third is denied again. The "steady trickle" section shows the other side of the algorithm: with capacity 1, waiting the exact refill interval (100ms at 10 tokens/sec) between each request keeps the bucket topped up just in time, every single request.

---

## Load Balancing

A load balancer sits in front of multiple backend instances and decides which one handles each incoming request. This section is mostly a **conceptual comparison** - the algorithms matter far more for their tradeoffs than for any single implementation - with one small piece of real, runnable code for the simplest algorithm.

### Round-Robin

Hands out requests to backends in a fixed rotating order: 1, 2, 3, 1, 2, 3, ...

```
requests: 1  2  3  4  5  6
             │  │  │  │  │  │
             ▼  ▼  ▼  ▼  ▼  ▼
          [ A  B  C  A  B  C ]   <- backend picked for each request
```

- Dead simple, no state needed beyond "who's next."
- Blind to load: if backend A is handling a slow request, round-robin still sends it the next one anyway.

### Least-Connections

Sends each request to whichever backend currently has the fewest active connections.

```
Backend A: 5 active connections
Backend B: 1 active connection   <- next request goes here
Backend C: 3 active connections
```

- Adapts to backends that are slower or already busy, unlike round-robin.
- Needs the load balancer to actually track connection counts per backend - more state, more complexity than round-robin.

### Consistent Hashing

Maps both backends and request keys (e.g., a user ID or cache key) onto a conceptual ring of hash values; a request is routed to the first backend found going clockwise from its own hash position.

```
                    0
                    │
        backend C ──┼── backend A
           ╲         │         ╱
            ╲        │        ╱
             ╲       │       ╱
              request "user-42"
              hashes to here, so it
              routes to backend A
              (next one clockwise)
```

- The key property: when a backend is added or removed, only a small fraction of keys need to move to a new backend - not all of them, unlike a naive `hash(key) % num_backends` scheme where changing the backend count reshuffles almost everything.
- Essential for **stateful sharding** (routing the same cache key or user ID to the same backend consistently, so a cache or in-memory store stays useful) at the cost of real implementation complexity over round-robin. Bonus Challenge 2 sketches the ring conceptually.

### Comparison

| Algorithm | State Needed | Adapts to Load? | Best For |
|---|---|---|---|
| Round-robin | Almost none | No | Uniform, stateless backends |
| Least-connections | Per-backend connection counts | Yes | Backends with variable request cost |
| Consistent hashing | The hash ring | No (routes by key, not load) | Sharded caches/stores where the same key must land on the same backend |

### A Real, Runnable Round-Robin Selector

This is the one load-balancing piece with genuine, runnable code - round-robin's whole logic is "remember an index, advance it" - simple enough to implement and verify for real:

```go
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
```

**Actual captured output (`go run main.go`):**

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

Eight requests over three backends wrap around exactly as expected: 1-2-3, 1-2-3, 1-2 - the `next` index advances with `(r.next + 1) % len(r.backends)`, the same modulo pattern from Level 4's operators.

---

## Database Scaling

This section is **conceptual/design reasoning** - there's no small runnable program that demonstrates "a database read replica," because it requires an actual second running database server to be meaningful. It bridges Level 29's `*sql.DB` connection pool settings into a wider architectural picture.

### Read Replicas

A **primary** database handles all writes. One or more **read replicas** continuously receive a copy of the primary's data (via replication) and serve read queries, spreading read load across multiple machines.

```
                  writes
                    │
                    ▼
              ┌───────────┐
              │  Primary  │
              │ (database)│
              └───────────┘
                    │
        replication │ (data flows one way)
        ┌───────────┼───────────┐
        ▼           ▼           ▼
   ┌─────────┐ ┌─────────┐ ┌─────────┐
   │Replica 1│ │Replica 2│ │Replica 3│
   └─────────┘ └─────────┘ └─────────┘
        ▲           ▲           ▲
        └───────────┴───────────┘
              reads (from the app)
```

- Reads scale horizontally (add more replicas); writes still funnel through one primary.
- Replication has lag - a replica might be milliseconds (or more, under load) behind the primary, so a read immediately after a write can miss it (**read-your-own-writes** is a real problem this design introduces - often solved by routing a user's own post-write reads to the primary for a short window).
- Most real-world read-heavy services (far more reads than writes) get a lot of headroom from this pattern alone, well before sharding is needed.

### Sharding

Splits data **horizontally** across multiple independent database servers, each holding a subset of the rows (e.g., by hashing a user ID) - because unlike read replicas, a single sharded database has no full copy of the data anywhere; each shard is authoritative for its own slice.

```
              shard_key = hash(user_id) % 4

   user 101 ──┐                       ┌── Shard 0: users 0, 4, 8, ...
   user 102 ──┤   routing logic       ├── Shard 1: users 1, 5, 9, ...
   user 103 ──┤   (in the app or a   ├── Shard 2: users 2, 6, 10, ...
   user 104 ──┘    proxy layer)       └── Shard 3: users 3, 7, 11, ...
```

- Solves what read replicas can't: write throughput and total data volume beyond what one primary can hold or handle.
- Real cost: queries that need to join or aggregate across shards become much harder (or impossible without a fan-out to every shard), shard rebalancing when adding a shard is a genuinely hard operational problem (this is exactly where consistent hashing from Section 5 earns its keep), and it's real added operational complexity that most systems should defer until read replicas and vertical scaling are provably not enough.

### Connection Pooling

Regardless of read replicas or sharding, every application instance talking to a database needs a **connection pool** - Level 29 already introduced the exact mechanism (`db.SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime`). The system design angle is this: those per-instance settings interact with horizontal scaling directly. If each of 20 application instances opens up to 25 connections, that's up to 500 connections the database server must be able to accept - a limit that's easy to blow past by scaling the application layer without ever revisiting the database's own `max_connections` setting (Section 12 lists this as a common mistake for exactly this reason).

```
20 app instances × SetMaxOpenConns(25) = up to 500 connections
                                            │
                                            ▼
                              does the DB server's max_connections
                              actually allow 500? if not: connection
                              errors under load, even though no
                              single instance is misconfigured
```

---

## Message Queues and Asynchronous Processing

Sometimes a request shouldn't wait for all the work it triggers to finish before responding. **Decoupling** producers (things that create work) from consumers (things that do the work) via a queue lets the producer move on immediately, and lets consumers process at their own pace - including retrying failed work without the original caller ever knowing.

### Why Decouple?

```
SYNCHRONOUS (no queue):
  client -> app: place order
  app -> email service: send confirmation  (client waits for this too!)
  app -> client: 200 OK
  (total latency = order processing + email sending, and if email is
   down, the ENTIRE request fails even though the order was fine)

ASYNCHRONOUS (with a queue):
  client -> app: place order
  app -> queue: enqueue "send confirmation email" job
  app -> client: 200 OK   (returns immediately - doesn't wait on email)
  ... separately, whenever it gets to it ...
  email worker -> queue: dequeue job
  email worker -> email service: send confirmation
```

Decoupling this way means a slow or temporarily-down email service no longer blocks (or breaks) order placement, and email-sending capacity can scale independently of order-placement capacity.

### In-Process Channels as an Analogy

Level 21's channels are exactly this pattern, just within a single process: a goroutine sends work on a channel (the producer), and one or more other goroutines receive and process it (the consumers) - decoupled in time (the producer doesn't block waiting for the work to finish) even though they're all still in the same running program.

```go
jobs := make(chan Order, 100) // buffered channel = in-process queue

// producer
go func() {
    jobs <- order
}()

// consumer(s)
go func() {
    for order := range jobs {
        sendConfirmationEmail(order)
    }
}()
```

### When In-Process Channels Are Enough vs When You Need a Real Message Queue

This is a genuine design decision, not a syntax question - and it's worth being honest that a channel and a message queue solve the same *shape* of problem at very different scales of durability guarantee:

- **In-process channels suffice when:** the work can be safely lost if the process crashes (it's not catastrophic to drop a few in-flight jobs), everything happens within one process or one deployable unit, and there's no need for multiple independent services to consume the same stream of work.
- **A real message queue (Kafka, RabbitMQ, Amazon SQS - named here, not implemented, since a meaningful example needs an actual running broker) is warranted when:** work must survive a process crash or restart (the queue is a separate durable system, not in-memory), multiple *different* services need to consume the same events, you need at-least-once delivery guarantees with retry/dead-letter handling built in, or producers and consumers are deployed and scaled completely independently of each other.

```
IN-PROCESS CHANNEL                    REAL MESSAGE QUEUE (Kafka/RabbitMQ/SQS)
────────────────────                  ────────────────────────────────────
Lives inside one process              A separate, durable service
Lost if the process crashes           Survives crashes/restarts
One producer, one process's           Many producers/consumers, possibly
  consumers                             different services entirely
Zero extra infrastructure             Infrastructure to run and operate
Great for: fan-out inside a           Great for: order events feeding
  single service                        billing, shipping, and analytics
                                         as three independent consumers
```

The rule of thumb: reach for a channel first because it's already there and it's simple; reach for a real queue only when the problem actually needs durability or cross-service fan-out that a channel structurally cannot provide.

---

## Circuit Breakers

A circuit breaker protects a system from repeatedly calling a dependency that's already failing - stopping the pile-up of slow, doomed requests (and the cascading failures they cause) by "tripping" and failing fast instead.

### The Pattern: Closed / Open / Half-Open

```
        failures reach threshold
   ┌───────────────────────────────┐
   │                                 ▼
┌────────┐                     ┌────────┐
│ CLOSED │                     │  OPEN  │
│(normal)│                     │(fail   │
│        │                     │ fast)  │
└────────┘                     └────────┘
   ▲                                │
   │                     reset timeout elapses
   │                                ▼
   │                          ┌───────────┐
   │      trial succeeds      │ HALF-OPEN │
   └───────────────────────── │(one trial │
        trial fails, back     │   call)   │
        to OPEN (not shown:   └───────────┘
        arrow to OPEN above)        │
                                     └──► (failure) back to OPEN
```

- **Closed:** normal operation. Every call goes through to the real dependency. Failures are counted.
- **Open:** the failure threshold has been hit. Calls are rejected immediately (without even trying the dependency) until a reset timeout elapses - this is the "fail fast" behavior that stops pile-up.
- **Half-open:** after the timeout, exactly one trial call is allowed through. Success closes the circuit again (back to normal); failure reopens it (and restarts the timeout).

### A Real, Runnable Circuit Breaker

This trips after N consecutive failures and recovers via a real elapsed timeout - genuinely run, not hand-computed:

```go
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
```

**Actual captured output (`go run main.go`, run twice for stability, both runs identical):**

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

Notice calls 4 and 5 fail **instantly** with `circuit breaker is open` - `fn` is never even invoked while `Open`, which is the entire point: fail fast instead of piling up slow, doomed calls against a dependency that's already down.

---

## Monolith vs Microservices

This is a pure design-tradeoffs discussion - there's no code to run, because the "right answer" genuinely depends on team size, domain complexity, and operational maturity, not on a technical fact a program could demonstrate.

### What Each Means

- **Monolith:** one deployable unit. All the business logic - orders, users, payments, notifications - lives in one codebase, built and deployed together, typically talking to one database (or one database per module, still deployed as one unit).
- **Microservices:** the same business logic split into multiple independently deployable services, each usually with its own datastore, communicating over the network (HTTP, gRPC, or a message queue from Section 7).

```
MONOLITH                              MICROSERVICES

┌───────────────────────┐             ┌────────┐  ┌────────┐  ┌──────────┐
│   orders │ users │ ... │             │ orders │  │ users  │  │ payments │
│      (one process,     │             │service │  │service │  │ service  │
│       one deploy)       │            └────────┘  └────────┘  └──────────┘
└───────────────────────┘                  │            │            │
           │                                └── network calls between them ──
           ▼                                        (HTTP/gRPC/queue)
     one database
```

### The Honest Tradeoffs

| | Monolith | Microservices |
|---|---|---|
| **Operational complexity** | Low - one thing to deploy, one thing to monitor, one log stream | High - many services to deploy, monitor, version, and debug across network boundaries |
| **Independent scaling** | No - scaling means running more copies of the *whole* thing, even the parts under no load | Yes - scale only the service that's actually under load |
| **Independent deployment** | No - any change requires redeploying the whole application | Yes - teams can ship their own service on their own schedule |
| **Debugging a request** | Straightforward - one process, one stack trace | Harder - a request may cross several services; needs distributed tracing to follow |
| **Team autonomy** | Lower - shared codebase, shared release train | Higher - each team can own a service end-to-end |
| **Consistency/transactions** | Easy - one database, real transactions | Hard - data lives in separate stores; cross-service consistency needs careful design (often eventual consistency, sagas) |
| **Where failure spreads** | A crash can take down the whole process | A failure in one service can be contained (especially with circuit breakers, Section 8) - or cascade badly if not |

### The Common Mistake: "Microservices by Default"

It's worth pushing back on directly: **splitting into microservices before you have the problems microservices solve is a common, expensive mistake.** Microservices trade operational simplicity for independent scaling and independent deployment - that trade is only worth making once a team is actually large enough that different groups stepping on each other's deploys is a real, recurring problem, or a specific part of the system genuinely needs to scale (or fail) independently of the rest.

A small team building a new product with microservices from day one typically pays the full operational tax - service discovery, distributed tracing, network failure handling, multiple CI/CD pipelines, versioned APIs between services - for a scaling problem they don't have yet, while still building the product's actual features slower because every change now crosses a network boundary. The far more common - and far more defensible - path is: **start as a well-structured monolith** (the Repository-Service-Handler layering from Level 30 keeps the internal boundaries clean even inside one deployable unit), and **extract a service only when a specific, measured need for independent scaling or deployment shows up.** That extraction is far easier from a cleanly-layered monolith than a tangled one - which is exactly why Level 30's pattern matters here even for a project that never splits into services at all.

---

## Case Study: Designing a URL Shortener

This section ties Sections 2-9 together into one worked example, in the shape of a real system design walkthrough. It's design reasoning throughout - a URL shortener's *value* is in the design decisions, not in a program whose output can be captured (Exercise 10 asks you to write this up formally as a design document).

### Requirements

**Functional:**
- `POST /shorten` with a long URL returns a short code (e.g., `abc123`).
- `GET /{code}` redirects to the original long URL.

**Non-functional (the numbers that actually drive the design):**
- Reads (redirects) vastly outnumber writes (shortens) - a realistic ratio is 100:1 or higher, since a link gets clicked far more often than it gets created.
- Redirects must be fast (low tens of milliseconds) - nobody tolerates a slow redirect.
- Short codes must not collide, and once created, a mapping is essentially permanent.

### Working Through the Concepts

**API shape:** two endpoints, matching the functional requirements above - `POST /shorten {url: "..."}` returning `{code: "abc123", short_url: "https://sho.rt/abc123"}`, and `GET /{code}` issuing an HTTP redirect (Level 27's `http.Redirect`) to the stored long URL.

**Storage choice:** a simple key-value shape (`code -> long_url`) is all this needs - no relational joins in the hot path. A relational database (Level 29's patterns) with a unique index on `code` works fine at moderate scale; at very large scale, a dedicated key-value store is a natural fit precisely because the access pattern never needs anything more than "look up this one key." Given the 100:1+ read:write ratio, **read replicas (Section 6)** are the very first scaling lever to pull once a single primary can't keep up with redirect traffic - sharding is not warranted until data volume or write throughput, not read throughput, becomes the bottleneck.

**Short code generation:** two realistic options - encode an auto-incrementing ID (from the datastore) into base62, or generate a random string and check for collisions before committing. The auto-increment approach guarantees no collisions and is simple; the random approach avoids a single monotonic counter becoming a bottleneck at very high write volume, at the cost of needing a collision check.

**Caching:** this is the highest-leverage lever in the whole design. Given a 100:1+ read:write ratio, an in-memory cache (this level's LRU cache, or Redis in a real deployment) sitting in front of the datastore, using the **cache-aside pattern (Section 3)**, turns the overwhelming majority of redirects into a cache hit that never touches the datastore at all. Popular links (the ones driving most of the traffic, by nature of being popular) stay hot in cache; the long tail of rarely-clicked links falls through to the datastore, which is fine because they're rare.

**Rate limiting:** `POST /shorten` needs a **token bucket (Section 4)** per API key or IP address, both to stop abuse (someone scripting mass link creation) and to protect the write path, which is far more expensive per-request than a cached redirect. `GET /{code}` redirects can be rate-limited too, but far more generously, since legitimate traffic is exactly what that endpoint exists to serve.

**Resilience:** if the datastore or cache becomes slow or unavailable, a **circuit breaker (Section 8)** around those calls stops a struggling datastore from piling up slow requests and taking the whole redirect path down with it - failing a redirect fast (or serving a stale cached mapping, if the design allows it) beats hanging every request until a timeout.

**Scaling plan:** the redirect service itself should be **stateless (Section 2)** - no server-side session, no in-memory data that only one instance has - so it scales horizontally behind a **load balancer (Section 5)**, most naturally round-robin or least-connections, since there's no need to route a given user's requests to a specific instance. As traffic grows, the order of scaling levers is: (1) cache in front of the datastore, (2) horizontally scale the stateless application layer, (3) add read replicas once the primary datastore itself is the bottleneck, (4) shard only if data volume or write throughput genuinely exceeds what replicas can handle - deliberately in that order, matching Section 11's "start simple, scale what actually needs scaling."

### What This Sets Up for Level 40

Level 40 (Real-World Projects) is where this case study - and the rest of this course - stops being a design document and becomes running code. This URL-shortener design (or the rate-limited public API variant from Exercise 10's alternative) is exactly the shape of project Level 40 has you actually build end-to-end: a real HTTP API (Level 27-28), backed by a real database (Level 29) behind the Repository-Service-Handler pattern (Level 30), with the caching and rate-limiting pieces from this level wired in for real, not just designed on paper.

---

## Best Practices

### 1. Design for Failure, Not Just the Happy Path

Assume every dependency - a database, a downstream API, a cache - will be slow or unavailable at some point, and decide *on purpose* what happens then: a circuit breaker (Section 8) that fails fast, a cached stale value served instead of a fresh one, a queued retry instead of a lost request. A system that only has a plan for success has, implicitly, a plan for failure too - it's just an undesigned one, discovered during an incident.

### 2. Start Simple and Scale What Actually Needs Scaling

Don't add a cache, a queue, a second database replica, or a microservice split before there's evidence the simple version can't keep up. Every one of these tools (Sections 3, 6, 7, 9) trades simplicity for a specific capability - added complexity that isn't paying for anything yet is pure cost. Section 12 ("designing for scale you don't have yet") is this same principle stated as a mistake to avoid.

### 3. Measure Before Optimizing

Level 38 (Production Debugging) covers the tools - profiling, tracing, metrics - for finding out *where* a system is actually slow or failing, rather than guessing. Apply that here: before reaching for a cache, a read replica, or a service split, get a number that shows the current design's actual bottleneck. Optimizing a part of the system that isn't the bottleneck doesn't help, and it adds complexity (see Best Practice 2) for no measured benefit.

### 4. Make Statelessness the Default for Application Servers

Per Section 2, a stateless application layer scales horizontally with almost no coordination cost. Push state that must persist into a dedicated store (a database, a cache, an object store) rather than application server memory, and horizontal scaling becomes a matter of adding instances behind a load balancer - not a redesign.

### 5. Put a Rate Limit on Every Public Endpoint

Any endpoint reachable without your own infrastructure in between deserves a rate limit (Section 4) before it ships, not after the first abuse incident. It's cheap to add up front and expensive to retrofit under fire.

---

## Common Mistakes

### Mistake 1: Adopting Microservices Before There's a Problem They Solve

```
❌ WRONG: A 3-person team splits a new product into 8 microservices on
   day one, "because that's how scalable systems are built."
   Result: 8x the deploy pipelines, network calls where function calls
   used to be, and a small team spending more time on service
   plumbing than on the product.

✅ RIGHT: Start as a well-layered monolith (Level 30's
   Repository-Service-Handler pattern). Extract a service only when a
   specific, measured need for independent scaling or deployment
   actually shows up.
```

### Mistake 2: Caching Without an Invalidation Strategy

```
❌ WRONG: Add a cache in front of the database to make reads fast, but
   never invalidate an entry when the underlying row changes.
   Result: users see stale data indefinitely, and nobody can predict
   for how long - "there are only two hard things in computer
   science: cache invalidation and naming things."

✅ RIGHT: Pick a pattern up front (cache-aside with explicit
   invalidation on write, Section 3) and apply it consistently, or set
   a deliberate TTL when perfect freshness doesn't matter.
```

### Mistake 3: No Rate Limiting on Public APIs

```
❌ WRONG: Ship a public endpoint with no rate limit "because real users
   won't hit it that hard."
   Result: one buggy client retry loop, or one deliberate abuser,
   takes the endpoint - and often the whole service behind it - down.

✅ RIGHT: Rate-limit every public endpoint (Section 4) before it ships,
   sized to real expected traffic with headroom, not added reactively
   after an incident.
```

### Mistake 4: Ignoring Database Connection Limits Under Load

```
❌ WRONG: Scale an application from 4 instances to 40, each with
   db.SetMaxOpenConns(25) (Level 29), without ever checking the
   database server's own max_connections.
   Result: up to 1,000 possible connections hitting a database
   configured for far fewer - connection errors under load that look
   like a database problem but are actually a scaling-coordination
   problem (Section 6).

✅ RIGHT: Treat the database's connection ceiling as a shared budget
   across every application instance, and size each instance's pool
   (and instance count) against that shared ceiling, not in isolation.
```

### Mistake 5: Designing for Scale You Don't Have Yet

```
❌ WRONG: A service handling 50 requests/minute is built with
   sharded databases, a message queue between every internal step,
   and a microservice per feature, "in case it needs to handle
   millions of users someday."
   Result: enormous complexity paid for today, for a scale that may
   never arrive - and if it does, the actual bottleneck likely won't
   be the one this design guessed at months or years earlier.

✅ RIGHT: Build the simplest design that meets today's real
   requirements, informed by system design thinking (so the
   *easy* scaling levers - statelessness, a cache, a load balancer -
   are there when needed) without pre-building the expensive ones
   (sharding, a service mesh) before there's evidence they're needed.
```

---

## Summary

**The Mindset Shift:**
- From "does this function work?" to "will this system hold up at scale, and does it fail gracefully?"

**Scalability:**
- Vertical scaling (bigger machine) is simple but has a ceiling; horizontal scaling (more machines) has no ceiling but needs coordination
- Stateless services scale horizontally with far less friction than stateful ones - this is why Level 32 favored JWTs over server-side sessions

**The Building Blocks (real, runnable code in this level):**
- **Caching** - cache-aside vs write-through (design), plus a real LRU-ish cache with verified eviction
- **Rate limiting** - the token-bucket algorithm, plus a real implementation verified to allow bursts and deny overflow
- **Load balancing** - round-robin, least-connections, consistent hashing compared; round-robin implemented and verified
- **Circuit breakers** - closed/open/half-open, plus a real implementation verified to trip and recover

**The Building Blocks (design reasoning, no runnable output):**
- **Database scaling** - read replicas for read throughput, sharding for write throughput/data volume, connection pooling as a shared budget across instances
- **Message queues** - in-process channels for simple in-process decoupling; a real broker (Kafka/RabbitMQ/SQS) when durability or cross-service fan-out is genuinely needed
- **Monolith vs microservices** - an honest tradeoff, not a default; start as a monolith, extract services when there's a measured need

**The Case Study:**
- A URL shortener ties every concept together: stateless API servers behind a load balancer, a cache-aside cache absorbing the overwhelming majority of reads, a token bucket protecting the write path, a circuit breaker around the datastore, and read replicas as the first real scaling lever - exactly the shape of project Level 40 builds for real

---

## Next Steps

You now understand:
- ✅ The mindset shift from correctness to systems that hold up at scale and fail gracefully
- ✅ Vertical vs horizontal scaling, and why statelessness is what makes horizontal scaling easy
- ✅ Caching patterns (cache-aside, write-through) and a real, working LRU-ish cache
- ✅ The token-bucket algorithm and a real, working rate limiter
- ✅ Load balancing algorithms (round-robin, least-connections, consistent hashing) and a real round-robin selector
- ✅ Database scaling levers - read replicas, sharding, connection pooling as a shared budget
- ✅ When in-process channels suffice vs when a real message queue is warranted
- ✅ The circuit breaker pattern and a real, working implementation that trips and recovers
- ✅ Honest monolith-vs-microservices tradeoffs, and why "microservices by default" is a mistake
- ✅ How to walk through a system design case study end-to-end, from requirements to a scaling plan

**Next level:** Level 40 - Real-World Projects
- Take everything from Levels 0-39 and build complete, real projects end-to-end
- Apply the Repository-Service-Handler pattern, real databases, and this level's system design thinking together
- Move from "I understand the pieces" to "I built the whole thing"

You've made the shift from writing correct code to designing systems that survive contact with real traffic and real failures. That's the last conceptual leap before Level 40 puts it all to work! 🚀
