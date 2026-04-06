# Day 05 — Atomics, Memory Ordering, and the Race Detector

## Goal
Understand when to use `sync/atomic` vs. mutex, what memory ordering Go's atomics guarantee, how the race detector works, and how to build lock-free data structures safely.

---

## Step 1 — Read: sync/atomic (30 min)

### What atomics guarantee

Go's atomic operations (`sync/atomic`) are:
- **Atomic**: the operation completes in one step, no partial reads/writes
- **Sequentially consistent**: all atomic operations appear in a single total order observed by all goroutines (since Go 1.19 formalized this)

Available operations:
```go
// Integer ops (int32, int64, uint32, uint64, uintptr)
atomic.AddInt64(&n, delta)
atomic.LoadInt64(&n)
atomic.StoreInt64(&n, val)
atomic.SwapInt64(&n, new)             // returns old
atomic.CompareAndSwapInt64(&n, old, new) // returns bool

// Pointer ops (Go 1.19+)
var p atomic.Pointer[MyStruct]
p.Store(&myStruct)
v := p.Load()     // returns *MyStruct
p.CompareAndSwap(old, new)

// Generic Value (any type, but type must be consistent)
var v atomic.Value
v.Store(myValue)
v.Load()
```

### When to use atomic vs. mutex

| Scenario | Use |
|----------|-----|
| Simple counter, flag, single int | `atomic` |
| Multiple fields updated together | `mutex` (must protect the invariant) |
| Read-heavy config pointer | `atomic.Pointer` |
| Complex data structure | `mutex` |
| Performance critical hot path | benchmark both |

### The classic mistake: atomics don't compose

```go
var a, b int64

// WRONG: even though each Load is atomic, the combination is not
if atomic.LoadInt64(&a) == atomic.LoadInt64(&b) {
    // another goroutine could change a or b between the two Loads
}

// RIGHT: use mutex to protect the invariant
mu.Lock()
eq := a == b
mu.Unlock()
```

---

## Step 2 — Read: Memory Ordering Deep Dive (20 min)

Modern CPUs reorder memory operations for performance. There are 4 memory ordering levels:
1. **Relaxed**: no ordering guarantees (just atomicity)
2. **Acquire**: a load that prevents reordering of subsequent reads/writes before it
3. **Release**: a store that prevents reordering of preceding reads/writes after it
4. **Sequential Consistency (SeqCst)**: total global order for all atomic ops

Go's `sync/atomic` operations are all **sequentially consistent** (as of the 2022 memory model spec). This is the strongest guarantee. It means:
- All goroutines agree on the order of ALL atomic operations
- No need to think about acquire/release in Go (unlike C++ `std::atomic` with explicit memory orders)

**Consequence:** Go atomics are safe but potentially slower than acquire/release in other languages. In practice, on modern x86 (which is TSO — total store order), the cost is negligible. On ARM, there's a real cost.

---

## Step 3 — Read: sync.Once Internals (15 min)

`sync.Once` is implemented with an atomic flag + mutex:
```go
type Once struct {
    done atomic.Uint32  // 0 = not done, 1 = done
    m    Mutex
}

func (o *Once) Do(f func()) {
    if o.done.Load() == 0 {  // fast path: atomic load
        o.doSlow(f)
    }
}

func (o *Once) doSlow(f func()) {
    o.m.Lock()
    defer o.m.Unlock()
    if o.done.Load() == 0 {  // double-check under lock
        defer o.done.Store(1)
        f()
    }
}
```

Key properties:
- After `Once.Do` returns, the function has completed — guaranteed by the mutex unlock HB subsequent lock
- The atomic `done` flag makes the fast path (already initialized) cheap — no mutex acquisition
- If `f` panics, `done` is still set to 1 — subsequent calls won't retry

---

## Step 4 — Experiment: Lock-Free Counter (30 min)

Compare mutex counter vs. atomic counter under contention:

```go
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
    "testing"
)

// Mutex counter
type MutexCounter struct {
    mu    sync.Mutex
    count int64
}

func (c *MutexCounter) Inc() { c.mu.Lock(); c.count++; c.mu.Unlock() }
func (c *MutexCounter) Get() int64 { c.mu.Lock(); defer c.mu.Unlock(); return c.count }

// Atomic counter
type AtomicCounter struct {
    count atomic.Int64
}

func (c *AtomicCounter) Inc() { c.count.Add(1) }
func (c *AtomicCounter) Get() int64 { return c.count.Load() }

func BenchmarkMutexCounter(b *testing.B) {
    c := &MutexCounter{}
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            c.Inc()
        }
    })
}

func BenchmarkAtomicCounter(b *testing.B) {
    c := &AtomicCounter{}
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            c.Inc()
        }
    })
}

func main() {
    results := testing.Benchmark(BenchmarkMutexCounter)
    fmt.Printf("Mutex:  %s\n", results)
    results = testing.Benchmark(BenchmarkAtomicCounter)
    fmt.Printf("Atomic: %s\n", results)
}
```

Run:
```bash
go test -bench=. -benchmem -cpu=1,2,4,8 .
```

Observe: atomics win at low goroutine counts; at high counts and high contention, the difference narrows because of cache-line bouncing (the atomic variable's cache line ping-pongs between CPU cores).

---

## Step 5 — Build: Lock-Free Stack with atomic.Pointer (45 min)

Treiber stack — a classic lock-free data structure:

```go
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

type node[T any] struct {
    value T
    next  *node[T]
}

// LockFreeStack is a lock-free LIFO stack using compare-and-swap.
type LockFreeStack[T any] struct {
    head atomic.Pointer[node[T]]
}

func (s *LockFreeStack[T]) Push(v T) {
    newNode := &node[T]{value: v}
    for {
        oldHead := s.head.Load()
        newNode.next = oldHead
        // If head is still oldHead, atomically replace it with newNode.
        // If another goroutine changed head, retry.
        if s.head.CompareAndSwap(oldHead, newNode) {
            return
        }
        // CAS failed: another goroutine pushed first, retry
    }
}

func (s *LockFreeStack[T]) Pop() (T, bool) {
    for {
        oldHead := s.head.Load()
        if oldHead == nil {
            var zero T
            return zero, false
        }
        // Try to move head to oldHead.next
        if s.head.CompareAndSwap(oldHead, oldHead.next) {
            return oldHead.value, true
        }
        // CAS failed: another goroutine popped, retry
    }
}

func main() {
    s := &LockFreeStack[int]{}

    var wg sync.WaitGroup
    // Push 1000 values concurrently
    for i := 0; i < 1000; i++ {
        i := i
        wg.Add(1)
        go func() {
            defer wg.Done()
            s.Push(i)
        }()
    }
    wg.Wait()

    // Pop all values
    count := 0
    for {
        _, ok := s.Pop()
        if !ok {
            break
        }
        count++
    }
    fmt.Printf("pushed 1000, popped %d\n", count)
}
```

Test with race detector:
```bash
go run -race main.go
```

**Understand the ABA problem:**
The Treiber stack above has an ABA problem in theory: if a node is popped and re-pushed before the CAS, the CAS succeeds incorrectly. In Go, this is usually safe because the GC ensures a popped node isn't immediately reused. In languages without GC, you need hazard pointers or tagged pointers. Understand this limitation.

---

## Step 6 — Read: How the Race Detector Works (20 min)

The race detector is based on **ThreadSanitizer (TSan)**. It:
1. Instruments every memory access and synchronization operation at compile time
2. At runtime, maintains a "shadow memory" for each memory location, tracking which goroutine last read/wrote it and what "vector clock" was in effect
3. If two goroutines access the same location without an ordering relationship, it reports a race

**Cost:**
- ~5–10x slowdown in CPU
- ~5–10x increase in memory usage
- Some false negatives (races it misses) — it only detects races that actually occur during the run

**Enable:**
```bash
go test -race ./...
go run -race main.go
go build -race .   # produces a race-detecting binary
```

**Always use in:**
- Unit and integration tests
- Local development
- CI pipelines

**Don't use in:**
- Production (too slow, too much memory)
- Exception: you can run a race-enabled binary in staging for a limited time

**Reading race output:**
```
==================
WARNING: DATA RACE
Write at 0x00c0000b4010 by goroutine 7:
  main.main.func1()
      /tmp/race.go:10 +0x40

Previous read at 0x00c0000b4010 by goroutine 6:
  main.main.func2()
      /tmp/race.go:15 +0x38

Goroutine 7 (running) created at:
  main.main()
      /tmp/race.go:8 +0x88
==================
```

The detector tells you: which variable, which goroutines, and where they were created. This is almost always enough to diagnose the issue.

---

## Step 7 — sync.Map vs. mutex-protected map (20 min)

```go
package main

import (
    "fmt"
    "sync"
    "testing"
)

// Mutex map
type SafeMap struct {
    mu sync.RWMutex
    m  map[string]int
}

func (m *SafeMap) Store(k string, v int) {
    m.mu.Lock()
    m.m[k] = v
    m.mu.Unlock()
}

func (m *SafeMap) Load(k string) (int, bool) {
    m.mu.RLock()
    v, ok := m.m[k]
    m.mu.RUnlock()
    return v, ok
}

// sync.Map
type SyncMap struct {
    m sync.Map
}

func (m *SyncMap) Store(k string, v int) { m.m.Store(k, v) }
func (m *SyncMap) Load(k string) (int, bool) {
    v, ok := m.m.Load(k)
    if !ok {
        return 0, false
    }
    return v.(int), true
}

func BenchmarkSafeMap(b *testing.B) {
    m := &SafeMap{m: make(map[string]int)}
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            m.Store("key", 1)
            m.Load("key")
        }
    })
}

func BenchmarkSyncMap(b *testing.B) {
    m := &SyncMap{}
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            m.Store("key", 1)
            m.Load("key")
        }
    })
}

func main() {
    fmt.Println("Run: go test -bench=. -cpu=1,4,8")
}
```

**When to use sync.Map:**
- Read-mostly with rare writes
- Many different goroutines each read/write different keys (reduces contention)

**When to use mutex + map:**
- Write-heavy workloads
- You need operations that touch multiple keys atomically
- You need to range over the map

---

## Day 05 Checklist

- [ ] Know the available atomic operations and their signatures
- [ ] Understand sequential consistency and why it matters in Go
- [ ] Know when atomics compose and when they don't
- [ ] Implemented the lock-free Treiber stack with CAS
- [ ] Benchmarked atomic vs. mutex counter at different goroutine counts
- [ ] Can read and interpret race detector output
- [ ] Know when to use `sync.Map` vs. mutex-protected map

## Key Takeaway

Go's atomics are sequentially consistent — the strongest guarantee. This makes them easier to reason about than C++ atomics with explicit memory orders, but you must still remember that atomics don't compose. For multi-field invariants, always use a mutex. The race detector is your best tool for finding concurrency bugs — run it in every test suite.
