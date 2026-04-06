# Day 06 — GC Internals and gctrace

## Goal
Understand how Go's garbage collector works, what triggers a GC cycle, how to observe it in production, and what causes GC pauses.

---

## Step 1 — Read: Go GC Algorithm (45 min)

Go uses a **concurrent, tricolor mark-and-sweep** garbage collector with a write barrier.

### The tricolor invariant

Every object is colored:
- **White**: not yet visited — candidates for collection
- **Gray**: discovered but not yet fully scanned
- **Black**: fully scanned — guaranteed reachable

**Invariant**: a black object must never point directly to a white object.

The GC maintains this by using a **write barrier**: whenever the running program writes a pointer, the write barrier colors the target object gray (ensures it won't be missed by the GC).

### GC phases:

```
1. Mark Setup (STW ~0.1ms)
   - Stop the world briefly
   - Enable write barriers on all goroutines
   - Scan stack roots, globals

2. Concurrent Mark (~80% of GC time)
   - World is RUNNING (your code runs concurrently with GC)
   - GC workers scan the gray set, turn objects black
   - Write barrier keeps marking correct as new pointers are written
   - GC assists: when an allocating goroutine is allocating faster than GC can mark,
     it's forced to do GC work proportional to how much it allocates (pacing)

3. Mark Termination (STW ~0.5ms)
   - Stop the world again
   - Final scan of stacks (goroutines may have written pointers during concurrent mark)
   - Disable write barriers

4. Concurrent Sweep
   - World is RUNNING
   - Sweep (free) white objects — overlaps with next allocation cycle
   - Go uses a span-based allocator (mspan): each span holds objects of one size class
```

### Key property: GC pauses are tiny
Modern Go GC has sub-millisecond STW pauses (typically <0.5ms). The "STW" concern from older versions is largely obsolete for most applications.

### GOGC — the primary knob

`GOGC` controls when the next GC cycle triggers:

```
next_heap_size = live_heap_after_last_GC * (1 + GOGC/100)
```

- Default: `GOGC=100` — GC triggers when heap doubles (100% growth)
- `GOGC=50` — more frequent GC, lower peak memory, higher CPU usage
- `GOGC=200` — less frequent GC, higher peak memory, lower CPU usage
- `GOGC=off` — disable GC (only for debugging, never production)

**Set programmatically:**
```go
import "runtime/debug"
debug.SetGCPercent(50)   // equivalent to GOGC=50
```

### GOMEMLIMIT (Go 1.19+)

`GOMEMLIMIT` is a soft limit on total memory the Go runtime will use. When memory usage approaches the limit, GC becomes more aggressive regardless of `GOGC`.

```bash
GOMEMLIMIT=256MiB go run main.go
```

```go
import "runtime/debug"
debug.SetMemoryLimit(256 * 1024 * 1024)
```

**Recommended pattern:** Set `GOGC=off` (or high) + `GOMEMLIMIT` to your memory budget. The runtime will GC just enough to stay under the limit. This avoids the problem of `GOGC=100` causing GC when there's plenty of free memory.

---

## Step 2 — Experiment: Observe GC with gctrace (45 min)

### Reading gctrace output

```bash
GODEBUG=gctrace=1 go run main.go
```

Output format:
```
gc 1 @0.008s 0%: 0.010+1.2+0.020 ms clock, 0.040+0.80/1.1/0+0.080 ms cpu, 4->4->2 MB, 5 MB goal, 0 MB stacks, 0 MB globals, 4 P
```

Breakdown:
- `gc 1` — GC cycle number
- `@0.008s` — time since program start
- `0%` — fraction of CPU time used by GC so far
- `0.010+1.2+0.020 ms clock` — wall time for (STW mark setup) + (concurrent mark) + (STW mark termination)
- `4->4->2 MB` — heap size at start of GC → end of marking → live after GC
- `5 MB goal` — target heap size for next GC trigger
- `4 P` — GOMAXPROCS

**What to look for:**
- STW times (first and last clock times) should be <1ms
- Heap ratio: if `start size` is much larger than `live after`, you're allocating a lot of short-lived objects
- GC CPU%: if consistently >10%, consider tuning GOGC or reducing allocations

### Experiment A — Allocation pressure

```go
package main

import (
    "fmt"
    "time"
)

func allocate() {
    // Allocate ~1MB of short-lived objects per call
    for i := 0; i < 1000; i++ {
        s := make([]byte, 1024)
        _ = s
    }
}

func main() {
    start := time.Now()
    for time.Since(start) < 5*time.Second {
        allocate()
    }
    fmt.Println("done")
}
```

Run:
```bash
GODEBUG=gctrace=1 go run main.go 2>&1 | head -30
```

Observe the GC frequency and heap sizes.

### Experiment B — GOGC impact

```bash
# More GC, lower memory
GOGC=25 GODEBUG=gctrace=1 go run main.go 2>&1 | grep "^gc " | wc -l

# Less GC, higher memory
GOGC=400 GODEBUG=gctrace=1 go run main.go 2>&1 | grep "^gc " | wc -l
```

Compare GC cycle counts and peak memory.

### Experiment C — GOMEMLIMIT vs GOGC

```go
package main

import (
    "fmt"
    "runtime"
    "runtime/debug"
    "time"
)

func printMemStats(label string) {
    var ms runtime.MemStats
    runtime.ReadMemStats(&ms)
    fmt.Printf("%s: HeapAlloc=%dMB, HeapSys=%dMB, NumGC=%d\n",
        label,
        ms.HeapAlloc/1024/1024,
        ms.HeapSys/1024/1024,
        ms.NumGC,
    )
}

func main() {
    // Keep 50MB of long-lived data
    longLived := make([][]byte, 50)
    for i := range longLived {
        longLived[i] = make([]byte, 1024*1024) // 1MB each
    }

    // Set a 100MB memory limit
    debug.SetMemoryLimit(100 * 1024 * 1024)
    debug.SetGCPercent(-1) // disable GOGC trigger, rely only on GOMEMLIMIT

    printMemStats("before allocation burst")
    
    for i := 0; i < 100; i++ {
        // Allocate 1MB short-lived, then discard
        buf := make([]byte, 1024*1024)
        _ = buf
    }
    
    printMemStats("after allocation burst")
    _ = longLived
    time.Sleep(100 * time.Millisecond)
    printMemStats("after sleep")
}
```

---

## Step 3 — Read: GC Assist and Pacing (20 min)

**GC assist** prevents runaway allocation from outpacing the GC:

If a goroutine allocates memory faster than the concurrent GC can mark, the allocating goroutine is paused and forced to do GC marking work proportional to the bytes it's allocating. This is invisible in normal code but shows up in profiling as time spent in `runtime.mallocgc`.

**When you see GC assist in pprof:**
- Your goroutine is allocating faster than GC can keep up
- Solutions: reduce allocation rate, increase `GOGC`, use `sync.Pool`, pre-allocate

**GC pacing algorithm:**
The runtime continuously adjusts the amount of work GC goroutines do to finish marking before the heap doubles. If background GC workers can't keep up, it starts forcing allocating goroutines to help (assist).

---

## Step 4 — Experiment: ReadMemStats and Monitoring (30 min)

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

type GCStats struct {
    NumGC        uint32
    PauseTotal   time.Duration
    LastPause    time.Duration
    HeapAlloc    uint64
    HeapInuse    uint64
    HeapReleased uint64
    StackInuse   uint64
    NumGoroutine int
}

func collectStats() GCStats {
    var ms runtime.MemStats
    runtime.ReadMemStats(&ms)  // NOTE: this stops the world briefly to get consistent snapshot
    return GCStats{
        NumGC:        ms.NumGC,
        PauseTotal:   time.Duration(ms.PauseTotalNs),
        LastPause:    time.Duration(ms.PauseNs[(ms.NumGC+255)%256]),
        HeapAlloc:    ms.HeapAlloc,
        HeapInuse:    ms.HeapInuse,
        HeapReleased: ms.HeapReleased,
        StackInuse:   ms.StackInuse,
        NumGoroutine: runtime.NumGoroutine(),
    }
}

func printStats(s GCStats) {
    fmt.Printf("GC cycles: %d | Last pause: %v | Total pause: %v\n",
        s.NumGC, s.LastPause, s.PauseTotal)
    fmt.Printf("Heap alloc: %dMB | In use: %dMB | Released: %dMB\n",
        s.HeapAlloc>>20, s.HeapInuse>>20, s.HeapReleased>>20)
    fmt.Printf("Stack: %dKB | Goroutines: %d\n",
        s.StackInuse>>10, s.NumGoroutine)
    fmt.Println("---")
}

func main() {
    // Background stats printer
    go func() {
        ticker := time.NewTicker(500 * time.Millisecond)
        defer ticker.Stop()
        for range ticker.C {
            printStats(collectStats())
        }
    }()

    // Workload: mix of short-lived and long-lived allocations
    var longLived [][]byte
    for i := 0; i < 30; i++ {
        // Short-lived
        for j := 0; j < 100; j++ {
            b := make([]byte, 10*1024) // 10KB
            _ = b
        }
        // Long-lived: accumulate 1MB every iteration
        longLived = append(longLived, make([]byte, 1024*1024))
        time.Sleep(100 * time.Millisecond)
    }

    _ = longLived
}
```

---

## Step 5 — GC Profiling with pprof (30 min)

```go
package main

import (
    "net/http"
    _ "net/http/pprof"  // registers /debug/pprof/ handlers
    "time"
)

func allocLoop() {
    for {
        s := make([]string, 1000)
        for i := range s {
            s[i] = fmt.Sprintf("item-%d", i)
        }
        time.Sleep(10 * time.Millisecond)
    }
}

func main() {
    go allocLoop()
    http.ListenAndServe(":6060", nil)
}
```

```bash
# Run the program
go run main.go &

# Capture heap profile
go tool pprof http://localhost:6060/debug/pprof/heap

# In pprof:
(pprof) top10
(pprof) list allocLoop
(pprof) web    # opens browser graph (requires graphviz)

# Useful pprof flags:
# -alloc_objects: total allocations (not just current)
# -alloc_space: total bytes allocated
# -inuse_objects: currently live objects
# -inuse_space: currently live bytes
go tool pprof -alloc_space http://localhost:6060/debug/pprof/heap
```

---

## Day 06 Checklist

- [ ] Can describe the 4 GC phases and which are STW vs. concurrent
- [ ] Know the tricolor invariant and why write barriers are needed
- [ ] Can read and interpret `gctrace=1` output
- [ ] Understand `GOGC` formula and its effect on frequency vs. memory
- [ ] Know `GOMEMLIMIT` and when to prefer it over `GOGC`
- [ ] Ran pprof heap profile and identified top allocators
- [ ] Understand GC assist and when it triggers

## Key Takeaway

Go's GC is designed for low latency, not throughput. STW pauses are typically <0.5ms. The main tuning knob is `GOGC` (frequency vs. memory tradeoff). For most services, the best modern approach is `GOMEMLIMIT` + high `GOGC` to avoid over-GCing when memory is available. If your GC CPU% is high, the root cause is almost always excessive allocation — fix that first.
