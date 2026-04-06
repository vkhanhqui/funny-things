# Day 04 — Go Memory Model and Happens-Before

## Goal
Understand the formal Go Memory Model, what "happens-before" means, which operations establish synchronization, and how to reason about what's safe to read without a race.

---

## Step 1 — Read: The Go Memory Model (60 min)

**Required reading:** https://go.dev/ref/mem (updated 2022 — read the full thing, ~30 min)

### Core concept: Happens-Before

A program has a **data race** if two goroutines access the same memory location concurrently, at least one is a write, and neither "happens-before" the other.

Happens-before (HB) is a partial order. If A HB B, then B is guaranteed to observe the effect of A. If neither A HB B nor B HB A, the reads and writes can be reordered by the compiler or CPU.

### What establishes happens-before:

**1. Within a goroutine:** all statements are sequentially ordered.
```go
x = 1
y = x + 1  // guaranteed to see x=1
```

**2. Goroutine creation:**
```go
x = 5
go func() {
    fmt.Println(x)  // guaranteed to see x=5 (go statement HB start of goroutine)
}()
```

**3. Channel send HB channel receive (for unbuffered):**
```go
var x int
ch := make(chan int)
go func() {
    x = 42
    ch <- 0  // send HB receive
}()
<-ch
fmt.Println(x)  // guaranteed to see x=42
```

**4. Channel close HB receive of zero value:**
```go
var x int
ch := make(chan int)
go func() {
    x = 42
    close(ch)  // close HB receive returning zero
}()
<-ch
fmt.Println(x)  // safe
```

**5. For buffered channel of capacity C: the Kth receive HB the (K+C)th send:**
```go
// This allows semaphore patterns
limit := make(chan struct{}, 3)  // max 3 concurrent
go func() {
    limit <- struct{}{}  // blocks if 3 already running
    defer func() { <-limit }()
    doWork()
}()
```

**6. sync.Mutex: Unlock HB subsequent Lock:**
```go
var mu sync.Mutex
var x int
mu.Lock()
x = 1
mu.Unlock()  // HB
// another goroutine:
mu.Lock()
fmt.Println(x)  // guaranteed to see x=1
mu.Unlock()
```

**7. sync.Once: the function passed to Do HB return from any Do:**
```go
var once sync.Once
var x int
once.Do(func() {
    x = 42
})
// All goroutines that call once.Do are guaranteed to see x=42 after it returns
```

**8. Goroutine exit — NOT ordered:**
```go
var x int
go func() {
    x = 42
}()
// main goroutine: x is NOT guaranteed to be 42 here
// goroutine exit does NOT happen-before anything
// you must use a channel or WaitGroup to synchronize
```

---

## Step 2 — Common Misconceptions (20 min)

### Misconception 1: Goroutines run sequentially unless proven otherwise

```go
var x, y int

go func() { x = 1 }()
go func() { y = 1 }()

// You CANNOT assume x==1 or y==1 here, even with time.Sleep
// Sleep is not a synchronization primitive
```

### Misconception 2: Reading a variable written in another goroutine is safe if you "check" first

```go
var done bool
var x int

go func() {
    x = 42
    done = true  // DATA RACE: no happens-before between goroutines
}()

for !done {}     // DATA RACE: spinning without synchronization
fmt.Println(x)   // x might be 0 even though done is true!
```

Fix with a channel or `atomic`:
```go
var x int
ch := make(chan struct{})
go func() {
    x = 42
    close(ch)  // close HB receive
}()
<-ch
fmt.Println(x)  // safe: 42
```

### Misconception 3: Package-level var initialization is synchronized

```go
// These run sequentially within a package (init order is defined),
// but across packages, init order follows import graph.
var a = computeA()
var b = computeB()  // guaranteed: computeA() finished before computeB() starts

// But: goroutines started in init() are NOT synchronized with main()
func init() {
    go func() { /* this might not have run when main() starts */ }()
}
```

---

## Step 3 — Experiment: Observe Data Races (45 min)

Write 5 programs that contain races. Verify each with `-race`. Then fix each.

**Race 1: counter without synchronization**
```go
var counter int
var wg sync.WaitGroup
for i := 0; i < 1000; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        counter++  // RACE
    }()
}
wg.Wait()
fmt.Println(counter)
```

**Race 2: map concurrent write**
```go
m := make(map[string]int)
var wg sync.WaitGroup
for i := 0; i < 100; i++ {
    i := i
    wg.Add(1)
    go func() {
        defer wg.Done()
        m[fmt.Sprintf("key%d", i)] = i  // RACE: maps are not safe for concurrent writes
    }()
}
wg.Wait()
```

**Race 3: slice append**
```go
var s []int
var wg sync.WaitGroup
for i := 0; i < 100; i++ {
    i := i
    wg.Add(1)
    go func() {
        defer wg.Done()
        s = append(s, i)  // RACE: slice header is not atomically updated
    }()
}
wg.Wait()
```

**Race 4: closure captures loop variable (Go <1.22)**
```go
var wg sync.WaitGroup
for i := 0; i < 5; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        fmt.Println(i)  // RACE in Go <1.22: i might have changed
    }()
}
wg.Wait()
// Fix: go func(i int) { ... }(i)
// In Go 1.22+: loop variable is per-iteration, no race
```

**Race 5: double-checked locking without atomic**
```go
type Singleton struct{ value int }
var instance *Singleton
var mu sync.Mutex

func getInstance() *Singleton {
    if instance == nil {  // RACE: read without lock
        mu.Lock()
        defer mu.Unlock()
        if instance == nil {
            instance = &Singleton{value: 42}
        }
    }
    return instance
}
// Fix: use sync.Once
```

For each race:
```bash
go run -race main.go
```

---

## Step 4 — Fix All Races (30 min)

Fix each of the 5 races using the appropriate synchronization:
- Race 1: `sync/atomic.AddInt64` or `sync.Mutex`
- Race 2: `sync.Mutex` around map access, or `sync.Map`
- Race 3: mutex-protected append
- Race 4: loop variable capture fix (or upgrade to Go 1.22)
- Race 5: `sync.Once`

---

## Step 5 — WaitGroup and Channel Synchronization Patterns (20 min)

```go
// Pattern 1: WaitGroup — wait for N goroutines
var wg sync.WaitGroup
results := make([]int, 10)
for i := 0; i < 10; i++ {
    i := i
    wg.Add(1)
    go func() {
        defer wg.Done()
        results[i] = expensive(i)  // safe: each i is a unique index
    }()
}
wg.Wait()
// All writes to results[] happened-before here

// Pattern 2: Channel as done signal
done := make(chan struct{})
var result int
go func() {
    result = 42
    close(done)  // HB
}()
<-done
fmt.Println(result)  // safe

// Pattern 3: errgroup for concurrent tasks with first-error
g, ctx := errgroup.WithContext(context.Background())
var mu sync.Mutex
var total int
for i := 0; i < 10; i++ {
    i := i
    g.Go(func() error {
        v, err := fetchValue(ctx, i)
        if err != nil {
            return err
        }
        mu.Lock()
        total += v
        mu.Unlock()
        return nil
    })
}
if err := g.Wait(); err != nil {
    log.Fatal(err)
}
fmt.Println(total)
```

---

## Day 04 Checklist

- [ ] Read the full Go Memory Model spec
- [ ] Can list all 8 happens-before rules from memory
- [ ] Know why `time.Sleep` is NOT a synchronization primitive
- [ ] Wrote and detected 5 data races with `-race`
- [ ] Fixed all 5 races with appropriate primitives
- [ ] Understand why goroutine exit is NOT ordered

## Key Takeaway

The memory model defines what you're allowed to assume. If you don't have an explicit happens-before relationship, the compiler and CPU are free to reorder your reads and writes — and they will. The race detector (`-race`) catches violations dynamically; it's not exhaustive but catches most real races. Run it always in tests.
