# Day 02 — Goroutine Lifecycle, Stack Growth, and Leak Detection

## Goal
Understand how goroutine stacks grow and shrink, how goroutines die, and how to detect and prevent goroutine leaks in production code.

---

## Step 1 — Read: Stack Growth Mechanics (30 min)

### How goroutine stacks work

Every goroutine starts with a small stack (currently 2KB in Go 1.21+, was 8KB in older versions).

**Stack growth (stack copying):**
1. On every function call, the compiler inserts a stack overflow check
2. If the remaining stack space is insufficient, the runtime allocates a **new, larger stack** (typically 2x)
3. All stack data is **copied** to the new stack — pointers on the stack are updated
4. The old stack is freed
5. Execution resumes on the new stack

**Why copying, not segmented stacks?**
- Go used to use segmented (linked) stacks (before Go 1.3)
- Problem: "hot split" — a loop straddling a segment boundary caused constant alloc/free
- Contiguous copying stacks solve this at the cost of O(stack_size) copy on growth

**Stack shrinking:**
- Happens at GC time
- If a goroutine's stack is using less than 1/4 of its capacity, the runtime may shrink it back down
- Shrinking also copies the stack to a smaller allocation

**Practical implication:**
```go
// This forces many stack growths — each recursive call needs more stack
func deepRecurse(n int) int {
    if n == 0 {
        return 0
    }
    var buf [1024]byte // 1KB local array
    _ = buf
    return deepRecurse(n-1) + 1
}
```

Run with `GODEBUG=gccheckmark=1` or look at stack allocations in pprof to observe.

### Goroutine lifecycle summary
```
go func() { ... }()
    ↓
_Grunnable (in run queue)
    ↓
_Grunning (executing on M+P)
    ↓ (channel block / mutex / timer)
_Gwaiting
    ↓ (unblocked)
_Grunnable
    ↓ (function returns)
_Gdead (available for reuse)
```

Goroutines in `_Gdead` are kept in a pool by the runtime and reused for new `go` statements — this is why goroutine creation is cheap after warm-up.

---

## Step 2 — Read: Goroutine Leaks (20 min)

A **goroutine leak** occurs when a goroutine is stuck in `_Gwaiting` forever — typically blocked on:
- A channel that nobody will ever send to (or receive from)
- A mutex that nobody will ever unlock
- A `select` with no viable case and no `default`
- An HTTP call with no timeout
- `time.After` in a loop (creates a new timer goroutine each iteration, old ones leak until they fire)

### Common leak patterns:

**Pattern 1: Abandoned channel receiver**
```go
func leak() {
    ch := make(chan int)
    go func() {
        v := <-ch  // blocks forever — nobody sends
        fmt.Println(v)
    }()
    // function returns, ch goes out of scope, goroutine leaks
}
```

**Pattern 2: Forgotten goroutine on error path**
```go
func fetch(urls []string) []string {
    ch := make(chan string)
    for _, url := range urls {
        go func(u string) {
            resp, err := http.Get(u)
            if err != nil {
                return  // goroutine exits, but ch never receives
            }
            ch <- resp.Status
        }(url)
    }
    // if we return early on error, remaining goroutines block on ch <- forever
    results := make([]string, len(urls))
    for i := range results {
        results[i] = <-ch
    }
    return results
}
```

**Pattern 3: time.After in a loop**
```go
func badRateLimit() {
    for {
        select {
        case <-time.After(1 * time.Second):  // creates a new timer goroutine every iteration
            doWork()
        }
    }
    // Fix: use time.NewTicker(1 * time.Second) outside the loop
}
```

---

## Step 3 — Install and Use goleak (30 min)

```bash
go get go.uber.org/goleak
```

### Using goleak in tests:

```go
package mypackage_test

import (
    "testing"
    "go.uber.org/goleak"
)

func TestMain(m *testing.M) {
    goleak.VerifyTestMain(m)  // checks for leaks after ALL tests
}

func TestSpecificFunction(t *testing.T) {
    defer goleak.VerifyNone(t)  // checks for leaks after this test
    
    leak()  // call the function under test
}
```

### Manual leak checking:

```go
func checkLeaks() {
    before := runtime.NumGoroutine()
    
    // ... do work ...
    
    after := runtime.NumGoroutine()
    if after > before {
        // print goroutine dump
        buf := make([]byte, 1<<20)
        n := runtime.Stack(buf, true)
        fmt.Printf("goroutine leak detected (+%d)\n%s\n", after-before, buf[:n])
    }
}
```

---

## Step 4 — Experiment: Stack Growth Profiling (30 min)

```go
package main

import (
    "fmt"
    "runtime"
    "runtime/debug"
)

// Measure stack size at various call depths
func measureStack(depth int) {
    if depth == 0 {
        var ms runtime.MemStats
        runtime.ReadMemStats(&ms)
        fmt.Printf("Stack in use: %d bytes\n", ms.StackInuse)
        return
    }
    var arr [512]int64  // 4KB on stack per frame
    _ = arr
    measureStack(depth - 1)
}

func main() {
    // Disable GC to get clean measurements
    debug.SetGCPercent(-1)
    
    for _, depth := range []int{10, 100, 1000} {
        fmt.Printf("\nDepth %d:\n", depth)
        go func(d int) {
            measureStack(d)
        }(depth)
        runtime.Gosched()
    }
    
    // Let goroutines run
    var input string
    fmt.Scanln(&input)
}
```

Run and observe how `StackInuse` scales with depth.

---

## Step 5 — Daily Coding Task: Robust Worker Pool with Graceful Shutdown (60 min)

Extend the Day 1 worker pool to handle:
1. Context-based cancellation
2. No goroutine leaks (verified with goleak)
3. Graceful drain (finish in-progress work before shutting down)

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"

    "golang.org/x/sync/errgroup"
)

type Task struct {
    ID int
    Fn func(ctx context.Context) error
}

type WorkerPool struct {
    tasks   chan Task
    workers int
}

func NewWorkerPool(workers int) *WorkerPool {
    return &WorkerPool{
        tasks:   make(chan Task, workers*10),
        workers: workers,
    }
}

// Run starts workers and blocks until ctx is cancelled and all tasks are drained.
func (p *WorkerPool) Run(ctx context.Context) error {
    g, ctx := errgroup.WithContext(ctx)

    for i := 0; i < p.workers; i++ {
        workerID := i
        g.Go(func() error {
            for {
                select {
                case task, ok := <-p.tasks:
                    if !ok {
                        fmt.Printf("worker %d: channel closed, exiting\n", workerID)
                        return nil
                    }
                    if err := task.Fn(ctx); err != nil {
                        return fmt.Errorf("task %d: %w", task.ID, err)
                    }
                case <-ctx.Done():
                    // drain remaining tasks before exiting
                    for {
                        select {
                        case task, ok := <-p.tasks:
                            if !ok {
                                return nil
                            }
                            _ = task.Fn(context.Background()) // finish with background ctx
                        default:
                            return ctx.Err()
                        }
                    }
                }
            }
        })
    }

    return g.Wait()
}

func (p *WorkerPool) Submit(task Task) {
    p.tasks <- task
}

func (p *WorkerPool) Close() {
    close(p.tasks)
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    pool := NewWorkerPool(4)

    var submittedCount int
    var mu sync.Mutex

    // Submit tasks in background
    go func() {
        for i := 0; ; i++ {
            select {
            case <-ctx.Done():
                pool.Close()
                return
            default:
                pool.Submit(Task{
                    ID: i,
                    Fn: func(ctx context.Context) error {
                        time.Sleep(200 * time.Millisecond)
                        mu.Lock()
                        submittedCount++
                        mu.Unlock()
                        return nil
                    },
                })
            }
        }
    }()

    if err := pool.Run(ctx); err != nil {
        fmt.Printf("pool error: %v\n", err)
    }

    fmt.Printf("total tasks completed: %d\n", submittedCount)
}
```

**Add goleak test:**
```go
package main

import (
    "context"
    "testing"
    "time"
    "go.uber.org/goleak"
)

func TestWorkerPoolNoLeak(t *testing.T) {
    defer goleak.VerifyNone(t)

    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()

    pool := NewWorkerPool(2)
    go func() {
        for i := 0; i < 10; i++ {
            pool.Submit(Task{ID: i, Fn: func(ctx context.Context) error {
                return nil
            }})
        }
        pool.Close()
    }()

    pool.Run(ctx)
}
```

Run:
```bash
go test -race -v ./...
```

---

## Step 6 — Experiment: time.After Leak (15 min)

Prove the `time.After` leak, then fix it:

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func leakyLoop(n int) {
    for i := 0; i < n; i++ {
        select {
        case <-time.After(1 * time.Hour):  // creates goroutine that leaks until timer fires
        default:
        }
    }
}

func fixedLoop(n int) {
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()
    for i := 0; i < n; i++ {
        select {
        case <-ticker.C:
        default:
        }
    }
}

func main() {
    before := runtime.NumGoroutine()
    leakyLoop(1000)
    after := runtime.NumGoroutine()
    fmt.Printf("Leaky: goroutines went from %d to %d (+%d)\n", before, after, after-before)

    before = runtime.NumGoroutine()
    fixedLoop(1000)
    after = runtime.NumGoroutine()
    fmt.Printf("Fixed: goroutines went from %d to %d (+%d)\n", before, after, after-before)
}
```

---

## Day 02 Checklist

- [ ] Can explain stack copying vs segmented stacks and why Go switched
- [ ] Know the 3 most common goroutine leak patterns
- [ ] Installed goleak and used it in a test
- [ ] Observed goroutine count via `runtime.NumGoroutine()`
- [ ] Built a worker pool with clean shutdown verified by goleak
- [ ] Proved and fixed the `time.After` leak

## Key Takeaway

Goroutine stacks grow by copying — cheap creation, O(n) growth cost. Leaks are always a blocked channel or a forgotten goroutine. Use `goleak` in every concurrent test and always give goroutines a way to exit via `context.Context` cancellation.
