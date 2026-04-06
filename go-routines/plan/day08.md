# Day 08 — Mutex, RWMutex, Deadlocks, and sync.Cond

## Goal
Master the full sync package: when to use each primitive, how to prevent deadlocks, and how to reason about lock ordering and contention.

---

## Step 1 — Read: sync.Mutex Internals (25 min)

`sync.Mutex` in Go is a **hybrid mutex** — it starts as a spinlock and falls back to OS parking when contention is high.

### States:
```
Locked bit (bit 0): is the mutex currently locked?
Woken bit (bit 1): is a waiting goroutine being woken?
Starving bit (bit 2): is the mutex in starvation mode?
Waiter count (bits 3+): how many goroutines are waiting?
```

### Normal mode vs. starvation mode:

**Normal mode** (default):
- Waiters are queued in FIFO
- When unlocked, the newly-unblocked goroutine competes with newly-arriving goroutines for the lock
- Newly-arriving goroutines often win (they're already running on CPU) — low latency
- Downside: a goroutine in the queue can starve if new goroutines keep arriving

**Starvation mode** (triggered if a goroutine waits >1ms):
- Ownership is transferred directly to the first waiter in queue
- New arrivals don't try to acquire the mutex — they go to the back of the queue
- Higher latency, but guaranteed progress for waiters
- Returns to normal mode when queue is drained

**Implication:** `sync.Mutex` is fair over time (no permanent starvation) but not FIFO in the short term. Don't assume lock-acquire order matches goroutine-creation order.

### sync.RWMutex

```go
var mu sync.RWMutex

// Multiple readers can hold simultaneously
mu.RLock()
defer mu.RUnlock()
v := sharedValue  // read

// Only one writer at a time; blocks all readers
mu.Lock()
defer mu.Unlock()
sharedValue = newValue  // write
```

**When RWMutex beats Mutex:**
- Read : write ratio >= 10:1
- Read-heavy config, cached data, read-mostly indices

**RWMutex gotcha — writer starvation prevention:**
When a writer is waiting, no new readers are admitted. This prevents writer starvation but means a burst of readers can delay a writer. Once the writer calls `Lock()`, all subsequent `RLock()` calls block until the writer finishes.

```go
// This can BLOCK if a writer is waiting:
mu.RLock()
// critical section
mu.RUnlock()
```

---

## Step 2 — Deadlock Analysis (30 min)

### Four conditions for deadlock (Coffman):
1. **Mutual exclusion**: only one goroutine can hold a resource
2. **Hold and wait**: a goroutine holds a lock and waits for another
3. **No preemption**: locks can't be forcibly taken
4. **Circular wait**: goroutine A holds what B needs, B holds what A needs

**Prevention: always acquire locks in the same global order.**

### Example deadlock:
```go
var mu1, mu2 sync.Mutex

// Goroutine 1
go func() {
    mu1.Lock()
    defer mu1.Unlock()
    mu2.Lock()       // waits for mu2
    defer mu2.Unlock()
}()

// Goroutine 2
go func() {
    mu2.Lock()
    defer mu2.Unlock()
    mu1.Lock()       // waits for mu1 — DEADLOCK
    defer mu1.Unlock()
}()
```

**Fix: always lock mu1 before mu2 in ALL goroutines.**

### Detecting deadlocks:
```bash
# Go runtime detects all-goroutines-blocked and panics:
fatal error: all goroutines are asleep - deadlock!

# For partial deadlocks (some goroutines still running), use pprof goroutine profile:
go tool pprof http://localhost:6060/debug/pprof/goroutine
(pprof) traces   # shows full stack traces for all goroutines
```

### Lock inversion (a subtler deadlock):
```go
type Cache struct {
    mu sync.Mutex
    db *Database
}

func (c *Cache) Refresh() {
    c.mu.Lock()
    defer c.mu.Unlock()
    // DB callback acquires db.mu internally:
    rows, _ := c.db.Query("SELECT ...")  // if db.Query also calls c.mu.Lock → deadlock
}
```
Always be aware of what lock the code you call under a lock might acquire.

---

## Step 3 — Experiment: Build and Debug a Deadlock (30 min)

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type Account struct {
    mu      sync.Mutex
    balance float64
    id      int
}

// BUGGY: can deadlock if two goroutines call Transfer(a, b) and Transfer(b, a) concurrently
func TransferBuggy(from, to *Account, amount float64) {
    from.mu.Lock()
    defer from.mu.Unlock()
    
    time.Sleep(1 * time.Millisecond)  // simulate work, increases deadlock probability
    
    to.mu.Lock()
    defer to.mu.Unlock()
    
    if from.balance < amount {
        return
    }
    from.balance -= amount
    to.balance += amount
}

// FIXED: always lock lower-ID account first
func TransferFixed(from, to *Account, amount float64) {
    // Establish a consistent lock order
    first, second := from, to
    if from.id > to.id {
        first, second = to, from
    }
    
    first.mu.Lock()
    defer first.mu.Unlock()
    second.mu.Lock()
    defer second.mu.Unlock()
    
    if from.balance < amount {
        return
    }
    from.balance -= amount
    to.balance += amount
}

func main() {
    alice := &Account{id: 1, balance: 1000}
    bob := &Account{id: 2, balance: 1000}

    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(2)
        go func() {
            defer wg.Done()
            TransferFixed(alice, bob, 10)  // swap to TransferBuggy to see deadlock
        }()
        go func() {
            defer wg.Done()
            TransferFixed(bob, alice, 10)
        }()
    }
    wg.Wait()
    
    fmt.Printf("Alice: %.0f, Bob: %.0f\n", alice.balance, bob.balance)
    fmt.Printf("Total: %.0f (should be 2000)\n", alice.balance+bob.balance)
}
```

Run both buggy and fixed versions. The buggy version will likely deadlock quickly.

---

## Step 4 — Read: sync.Cond (20 min)

`sync.Cond` provides a way to **broadcast a condition change** to multiple waiters efficiently — something channels can't easily do.

```go
type sync.Cond struct {
    L Locker  // usually a *sync.Mutex or *sync.RWMutex
}

cond.Wait()      // atomically releases L and suspends goroutine; re-acquires L on wake
cond.Signal()    // wakes one waiting goroutine
cond.Broadcast() // wakes ALL waiting goroutines
```

### Classic use case: condition variable pattern

```go
var mu sync.Mutex
var ready bool
cond := sync.NewCond(&mu)

// Waiter
go func() {
    mu.Lock()
    for !ready {           // ALWAYS check condition in a loop (spurious wakeups)
        cond.Wait()        // releases mu, sleeps, re-acquires mu on wake
    }
    // condition is now true
    mu.Unlock()
}()

// Notifier
mu.Lock()
ready = true
cond.Broadcast()  // wake all waiters
mu.Unlock()
```

### When to use sync.Cond vs. channels:

| | sync.Cond | Channel |
|--|-----------|---------|
| N waiters waiting for same condition | Broadcast() — perfect | Would need N channels or repeated sends |
| One waiter waiting for one event | Overkill | Simple `<-ch` |
| Pass data along with signal | Awkward | Natural |
| Cancel/timeout | Awkward | `select` with ctx |

### sync.Cond practical example:

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Pool that blocks when empty and broadcasts when items are added
type BlockingPool struct {
    mu    sync.Mutex
    cond  *sync.Cond
    items []string
}

func NewBlockingPool() *BlockingPool {
    p := &BlockingPool{}
    p.cond = sync.NewCond(&p.mu)
    return p
}

func (p *BlockingPool) Put(item string) {
    p.mu.Lock()
    p.items = append(p.items, item)
    p.cond.Signal()  // wake one waiter
    p.mu.Unlock()
}

func (p *BlockingPool) Get() string {
    p.mu.Lock()
    defer p.mu.Unlock()
    for len(p.items) == 0 {
        p.cond.Wait()  // blocks until Put() calls Signal()
    }
    item := p.items[0]
    p.items = p.items[1:]
    return item
}

func main() {
    pool := NewBlockingPool()

    // Consumer
    go func() {
        for i := 0; i < 5; i++ {
            item := pool.Get()
            fmt.Printf("got: %s\n", item)
        }
    }()

    // Producer (delayed)
    for i := 0; i < 5; i++ {
        time.Sleep(300 * time.Millisecond)
        pool.Put(fmt.Sprintf("item-%d", i))
    }
    
    time.Sleep(100 * time.Millisecond)
}
```

---

## Step 5 — Daily Coding Task: Concurrent Cache with RWMutex (60 min)

Build a TTL-aware concurrent cache:

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"
)

type entry struct {
    value     any
    expiresAt time.Time
}

type Cache struct {
    mu      sync.RWMutex
    items   map[string]entry
    ttl     time.Duration
}

func NewCache(ttl time.Duration) *Cache {
    c := &Cache{
        items: make(map[string]string),
        ttl:   ttl,
    }
    return c
}

func (c *Cache) Set(key string, value any) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items[key] = entry{
        value:     value,
        expiresAt: time.Now().Add(c.ttl),
    }
}

func (c *Cache) Get(key string) (any, bool) {
    c.mu.RLock()
    e, ok := c.items[key]
    c.mu.RUnlock()
    
    if !ok || time.Now().After(e.expiresAt) {
        return nil, false
    }
    return e.value, true
}

func (c *Cache) Delete(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    delete(c.items, key)
}

// StartCleanup runs a background goroutine to purge expired entries.
// It stops when ctx is cancelled.
func (c *Cache) StartCleanup(ctx context.Context, interval time.Duration) {
    go func() {
        ticker := time.NewTicker(interval)
        defer ticker.Stop()
        for {
            select {
            case <-ticker.C:
                c.evictExpired()
            case <-ctx.Done():
                return
            }
        }
    }()
}

func (c *Cache) evictExpired() {
    now := time.Now()
    c.mu.Lock()
    defer c.mu.Unlock()
    for k, e := range c.items {
        if now.After(e.expiresAt) {
            delete(c.items, k)
        }
    }
}

func (c *Cache) Len() int {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return len(c.items)
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    cache := NewCache(1 * time.Second)
    cache.StartCleanup(ctx, 500*time.Millisecond)
    
    // Concurrent reads and writes
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        i := i
        wg.Add(1)
        go func() {
            defer wg.Done()
            cache.Set(fmt.Sprintf("key%d", i), i*100)
        }()
    }
    wg.Wait()
    
    fmt.Printf("cache size: %d\n", cache.Len())
    
    // Wait for TTL expiry
    time.Sleep(1500 * time.Millisecond)
    fmt.Printf("after TTL: %d\n", cache.Len())
}
```

**Extend it:**
1. Add `GetOrSet(key, computeFn)` that avoids thundering herd using `singleflight`
2. Add cache metrics (hit rate, eviction count) with atomic counters
3. Run with `-race` to verify no data races

---

## Day 08 Checklist

- [ ] Know the difference between Mutex normal mode and starvation mode
- [ ] Know when RWMutex is faster than Mutex (read/write ratio)
- [ ] Can identify the 4 Coffman conditions for deadlock
- [ ] Can apply consistent lock ordering to prevent deadlocks
- [ ] Built and debugged the TransferBuggy deadlock
- [ ] Know when to use sync.Cond vs. channels
- [ ] Implemented the TTL cache with RWMutex and cleanup goroutine

## Key Takeaway

`sync.Mutex` is not simple — it has normal and starvation modes to balance latency vs. fairness. Deadlocks always involve circular waits; the fix is always consistent lock ordering. `sync.RWMutex` is only worth it when reads greatly outnumber writes — it has overhead too. `sync.Cond` is the right tool when N goroutines need to react to the same state change.
