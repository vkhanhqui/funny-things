# Day 01 — GMP Scheduler Internals

## Goal
Understand how Go schedules goroutines onto OS threads, why it's fast, and how to observe and reason about scheduling behavior.

---

## Step 1 — Read: The GMP Model (45 min)

### Concepts to internalize:

**G — Goroutine**
- Lightweight user-space thread, starts with a 2KB–8KB stack
- Contains: stack pointer, program counter, goroutine state (running/runnable/waiting/dead)
- Created cheaply — benchmark `go func(){}()` vs. `pthread_create`

**M — Machine (OS thread)**
- Actual OS thread managed by the runtime
- An M must hold a P to run Go code
- An M without a P can still run: syscalls, cgo, GC helpers

**P — Processor (logical CPU)**
- Count set by `GOMAXPROCS` (default = number of CPU cores)
- Each P has a local run queue (LRQ) of up to 256 goroutines
- Global run queue (GRQ) handles overflow and newly created goroutines

### Scheduling loop (simplified):
```
for {
    g := findRunnable()   // LRQ → GRQ → work-steal from other P
    execute(g)            // run until block, syscall, or preemption
}
```

### Work stealing:
- When a P's LRQ is empty, it steals half the goroutines from a random other P's LRQ
- This keeps all P's busy without a central lock on every schedule

### Preemption (Go 1.14+ async preemption):
- Before 1.14: preemption only at function call sites (a tight loop could block a thread)
- Since 1.14: the runtime sends SIGURG to threads running long goroutines, forcing a preemption at any safe point
- Consequence: goroutines are preemptible even without function calls

**Read these:**
- https://go.dev/src/runtime/HACKING.md (skim — 15 min)
- https://go.dev/blog/sched (Dmitry Vyukov's original design notes)

---

## Step 2 — Experiment: Observe the Scheduler (45 min)

### Experiment A — schedtrace baseline

Create `main.go`:
```go
package main

import (
    "fmt"
    "runtime"
    "sync"
)

func main() {
    runtime.GOMAXPROCS(2)
    var wg sync.WaitGroup

    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            // CPU-bound work
            sum := 0
            for j := 0; j < 1_000_000; j++ {
                sum += j
            }
            fmt.Printf("goroutine %d done, sum=%d\n", id, sum)
        }(i)
    }
    wg.Wait()
}
```

Run with scheduler trace:
```bash
GODEBUG=schedtrace=200 GOMAXPROCS=2 go run main.go 2>&1 | grep "^SCHED"
```

**Understand each field in the output:**
```
SCHED 204ms: gomaxprocs=2 idleprocs=0 threads=4 spinningthreads=0 needspinning=0 idlethreads=1 runqueue=3 [2 1]
```
- `gomaxprocs=2` — 2 P's active
- `idleprocs=0` — no idle P's (all busy)
- `threads=4` — 4 OS threads (2 running + extras for syscalls/GC)
- `runqueue=3` — 3 goroutines in the global run queue
- `[2 1]` — LRQ sizes for P0 and P1

### Experiment B — GOMAXPROCS impact

Run the same program with different GOMAXPROCS and time it:
```bash
for n in 1 2 4 8; do
    echo "GOMAXPROCS=$n:"
    GOMAXPROCS=$n go run main.go
done
```
Observe: beyond the number of physical cores, more GOMAXPROCS doesn't help CPU-bound work.

### Experiment C — Syscall behavior

Create a program that blocks on I/O:
```go
package main

import (
    "net/http"
    "sync"
)

func main() {
    var wg sync.WaitGroup
    for i := 0; i < 20; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            http.Get("http://httpbin.org/delay/1")
        }()
    }
    wg.Wait()
}
```

Run with schedtrace and observe: `threads` count rises above `gomaxprocs` because the runtime creates new Ms to keep Ps busy while others are in syscalls.

---

## Step 3 — Experiment: GOMAXPROCS and Goroutine Count (30 min)

```go
package main

import (
    "fmt"
    "runtime"
    "sync/atomic"
    "time"
)

func main() {
    var count int64
    start := time.Now()

    for i := 0; i < 1_000_000; i++ {
        go func() {
            atomic.AddInt64(&count, 1)
        }()
    }

    // Wait for all goroutines
    for atomic.LoadInt64(&count) < 1_000_000 {
        runtime.Gosched()
    }

    fmt.Printf("1M goroutines completed in %v\n", time.Since(start))
    fmt.Printf("NumGoroutine at end: %d\n", runtime.NumGoroutine())
}
```

**Questions to answer:**
1. How long does it take to spawn and run 1M goroutines?
2. What's the peak memory usage? (`/usr/bin/time -l go run main.go` on mac)
3. How does this compare to spawning 1M OS threads? (you can't — default ulimit is ~10K)

---

## Step 4 — Read: Goroutine States (20 min)

A goroutine transitions through these states:
```
_Gidle       → freshly allocated
_Grunnable   → in a run queue, ready to run
_Grunning    → currently executing on an M
_Gsyscall    → in a syscall (M is blocked, P is released)
_Gwaiting    → blocked (channel, mutex, timer, GC, etc.)
_Gdead       → finished or being recycled
```

Key transitions to understand:
- `_Grunning → _Gsyscall`: P is immediately handed off to another M (or a new M is created)
- `_Grunning → _Gwaiting`: goroutine blocked on channel, P runs another goroutine
- `_Gwaiting → _Grunnable`: channel unblocked, goroutine put back in a run queue

---

## Step 5 — Daily Coding Task (60 min)

Build a CPU-core-saturating worker pool:

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

// WorkerPool runs tasks on exactly GOMAXPROCS goroutines
type WorkerPool struct {
    tasks chan func()
    wg    sync.WaitGroup
}

func NewWorkerPool(size int) *WorkerPool {
    p := &WorkerPool{
        tasks: make(chan func(), size*10),
    }
    for i := 0; i < size; i++ {
        p.wg.Add(1)
        go p.worker()
    }
    return p
}

func (p *WorkerPool) worker() {
    defer p.wg.Done()
    for task := range p.tasks {
        task()
    }
}

func (p *WorkerPool) Submit(task func()) {
    p.tasks <- task
}

func (p *WorkerPool) Shutdown() {
    close(p.tasks)
    p.wg.Wait()
}

func main() {
    numCPU := runtime.NumCPU()
    pool := NewWorkerPool(numCPU)

    start := time.Now()
    var mu sync.Mutex
    results := make([]int, 100)

    for i := 0; i < 100; i++ {
        i := i
        pool.Submit(func() {
            sum := 0
            for j := 0; j < 10_000_000; j++ {
                sum += j
            }
            mu.Lock()
            results[i] = sum
            mu.Unlock()
        })
    }

    pool.Shutdown()
    fmt.Printf("100 tasks on %d workers: %v\n", numCPU, time.Since(start))
    fmt.Printf("First result: %d\n", results[0])
}
```

**Extend it:**
1. Add metrics: tasks completed per second
2. Make pool size configurable at runtime via `GOMAXPROCS`
3. Use `GODEBUG=schedtrace=500` to observe the pool in action

---

## Day 01 Checklist

- [ ] Can explain GMP model from memory
- [ ] Know what `GOMAXPROCS` controls and its default
- [ ] Understand work-stealing: when it happens and why
- [ ] Know difference between pre-1.14 and 1.14+ preemption
- [ ] Can read and interpret `schedtrace` output
- [ ] Implemented and ran the worker pool
- [ ] Observed goroutine state transitions via `schedtrace`

## Key Takeaway

The scheduler is cooperative-ish (goroutines yield at channel ops, syscalls, function calls) but preemptive since 1.14. The GMP model avoids per-goroutine OS threads while keeping CPUs busy via work-stealing. `GOMAXPROCS` controls parallelism, not concurrency.
