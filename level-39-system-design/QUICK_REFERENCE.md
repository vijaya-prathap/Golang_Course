# Level 39: Quick Reference Card

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir myapp && cd myapp
go mod init myapp

# Create main.go - a minimal token-bucket rate limiter
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "sync"
    "time"
)

type TokenBucket struct {
    mu         sync.Mutex
    capacity   float64
    tokens     float64
    refillRate float64
    lastRefill time.Time
}

func NewTokenBucket(capacity, refillRate float64) *TokenBucket {
    return &TokenBucket{capacity: capacity, tokens: capacity, refillRate: refillRate, lastRefill: time.Now()}
}

func (b *TokenBucket) Allow() bool {
    b.mu.Lock()
    defer b.mu.Unlock()
    now := time.Now()
    b.tokens += now.Sub(b.lastRefill).Seconds() * b.refillRate
    if b.tokens > b.capacity {
        b.tokens = b.capacity
    }
    b.lastRefill = now
    if b.tokens >= 1 {
        b.tokens--
        return true
    }
    return false
}

func main() {
    bucket := NewTokenBucket(3, 5)
    for i := 1; i <= 5; i++ {
        fmt.Printf("Request %d: allowed=%v\n", i, bucket.Allow())
    }
}
EOF

# Run
go run main.go
```

---

## 📋 The Mindset Shift

```
"Does this function work?"  ──►  "Will this system hold up at scale,
                                   and does it fail gracefully?"
```

---

## 📈 Scaling

```go
// Vertical:   bigger machine        - simple, has a ceiling
// Horizontal: more machines         - no ceiling, needs coordination

// Stateless services scale horizontally easily - no server-side
// session to replicate, any instance can handle any request.
// Stateful services need sticky sessions or a shared store (Redis/DB).
```

---

## 🗄️ Caching Patterns

```
Cache-aside:   app checks cache -> miss -> fetch DB -> populate cache
               write -> DB, then invalidate cache entry

Write-through: write -> cache -> cache writes DB synchronously
               read -> almost always a cache hit
```

```go
// LRU cache core shape (container/list + map + sync.Mutex)
type LRUCache struct {
    mu       sync.Mutex
    capacity int
    items    map[string]*list.Element
    order    *list.List // front = most recent, back = least recent
}
// Get: MoveToFront on hit
// Put: PushFront, evict order.Back() if over capacity
```

---

## 🪣 Token-Bucket Rate Limiting

```go
// capacity tokens, refills at refillRate tokens/sec
// Allow(): refill based on elapsed time, then consume 1 token if available
func (b *TokenBucket) Allow() bool {
    b.refill()
    if b.tokens >= 1 {
        b.tokens--
        return true
    }
    return false
}
```

Allows a **burst** up to capacity, then smooths to the refill rate. Contrast with a **sliding window** (counts actual requests in the last N ms, no pre-filled burst).

---

## ⚖️ Load Balancing

```go
// Round-robin: simplest, blind to load
func (r *RoundRobin) Next() string {
    backend := r.backends[r.next]
    r.next = (r.next + 1) % len(r.backends)
    return backend
}
```

| Algorithm | Adapts to Load? | Best For |
|---|---|---|
| Round-robin | No | Uniform request cost |
| Least-connections | Yes | Variable request cost |
| Consistent hashing | No | Key-to-backend stickiness (sharded caches) |

---

## 🗃️ Database Scaling

```
Read replicas: scale READS, writes still go through one primary
Sharding:      scale WRITES/data volume, splits rows across servers
               (no full copy anywhere - harder cross-shard queries)

Connection budget = instances × SetMaxOpenConns(n)
                     MUST stay under the DB server's max_connections
```

---

## 📬 Channel vs Message Queue

```
In-process channel:  same process, lost on crash, one consumer type
                      -> fine for in-process fan-out (e.g. thumbnail resize)

Real queue (Kafka/RabbitMQ/SQS): durable, survives crashes, many
                      independent consumer services
                      -> needed for cross-service fan-out or durability
```

---

## 🔌 Circuit Breaker

```go
// States: Closed -> (N failures) -> Open -> (timeout) -> HalfOpen
//         HalfOpen success -> Closed
//         HalfOpen failure -> Open (timeout restarts)

func (cb *CircuitBreaker) Call(fn func() error) error {
    if cb.state == Open && !timeoutElapsed {
        return ErrCircuitOpen // fail fast, fn is never called
    }
    // ... call fn, update state based on result
}
```

---

## 🏗️ Monolith vs Microservices

```
Start:   well-layered monolith (Repository-Service-Handler, Level 30)
Extract: only when there's a MEASURED need for independent
         scaling or independent deployment - not by default
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Premature microservices | Split into services on day one | Start monolith, extract on measured need |
| No cache invalidation | "It'll refresh eventually" | Explicit invalidation on write (+ TTL backstop) |
| No rate limiting | Public endpoint, no limit | Token bucket on every public endpoint |
| Ignoring connection limits | `instances × maxConns` unchecked | Treat DB's `max_connections` as a shared budget |
| Designing for scale you don't have | Sharding + queues for 50 req/min | Build for today's real requirements |

---

## 🎓 Before Next Level

Can you:
- [ ] Explain the mindset shift from "does it work" to "does it hold up at scale"?
- [ ] Explain why stateless services scale horizontally more easily?
- [ ] Build a working LRU cache, token-bucket limiter, round-robin selector, and circuit breaker?
- [ ] Compare load-balancing algorithms for a given workload?
- [ ] Explain read replicas vs sharding, and connection limits as a shared budget?
- [ ] Decide between an in-process channel and a real message queue?
- [ ] Make an honest monolith-vs-microservices call for a given scenario?

If YES → You're ready for Level 40!

---

## 📚 Next Level

Level 40: Real-World Projects
- Build complete, real projects end-to-end using everything from Levels 0-39
- Apply the Repository-Service-Handler pattern, real databases, and system design thinking together

You've got system design down! 💪
