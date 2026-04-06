# Day 10 — End-to-End Profiling and Tuning

## Goal
Bring everything together: write a non-trivial concurrent program, benchmark it, profile it across all dimensions (CPU, memory, goroutines, mutex contention, GC), tune it, and interpret the execution tracer.

---

## Step 1 — The Capstone Program (60 min)

Build a concurrent image resizing pipeline. This will exercise:
- Goroutine scheduling (fan-out work)
- GC pressure (large byte buffers)
- Mutex contention (result aggregation)
- Context cancellation (timeout handling)

```go
package main

import (
    "bytes"
    "context"
    "fmt"
    "image"
    "image/color"
    "image/jpeg"
    "math/rand"
    "sync"
    "sync/atomic"
    "time"

    "golang.org/x/sync/errgroup"
)

// generateImage creates a synthetic RGBA image
func generateImage(width, height int) image.Image {
    img := image.NewRGBA(image.Rect(0, 0, width, height))
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            img.Set(x, y, color.RGBA{
                R: uint8(rand.Intn(256)),
                G: uint8(rand.Intn(256)),
                B: uint8(rand.Intn(256)),
                A: 255,
            })
        }
    }
    return img
}

// resizeImage does a simple nearest-neighbor resize
func resizeImage(src image.Image, newWidth, newHeight int) image.Image {
    srcBounds := src.Bounds()
    srcW := srcBounds.Max.X - srcBounds.Min.X
    srcH := srcBounds.Max.Y - srcBounds.Min.Y
    
    dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
    for y := 0; y < newHeight; y++ {
        for x := 0; x < newWidth; x++ {
            srcX := x * srcW / newWidth
            srcY := y * srcH / newHeight
            dst.Set(x, y, src.At(srcBounds.Min.X+srcX, srcBounds.Min.Y+srcY))
        }
    }
    return dst
}

// encodeJPEG encodes image to JPEG bytes
func encodeJPEG(img image.Image, quality int) ([]byte, error) {
    var buf bytes.Buffer
    err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
    return buf.Bytes(), err
}

type ProcessResult struct {
    ID    int
    Bytes int
    Error error
}

type Pipeline struct {
    workers    int
    results    []ProcessResult
    mu         sync.Mutex
    processed  atomic.Int64
    totalBytes atomic.Int64
}

func NewPipeline(workers int) *Pipeline {
    return &Pipeline{workers: workers}
}

func (p *Pipeline) Run(ctx context.Context, imageCount int) error {
    type job struct {
        id    int
        image image.Image
    }

    jobs := make(chan job, p.workers*2)

    // Producer: generate synthetic images
    go func() {
        defer close(jobs)
        for i := 0; i < imageCount; i++ {
            img := generateImage(800, 600)
            select {
            case jobs <- job{id: i, image: img}:
            case <-ctx.Done():
                return
            }
        }
    }()

    // Workers: resize + encode
    g, ctx := errgroup.WithContext(ctx)
    g.SetLimit(p.workers)

    for job := range jobs {
        job := job
        g.Go(func() error {
            resized := resizeImage(job.image, 200, 150)  // thumbnail
            encoded, err := encodeJPEG(resized, 85)
            if err != nil {
                return fmt.Errorf("job %d: %w", job.id, err)
            }

            p.processed.Add(1)
            p.totalBytes.Add(int64(len(encoded)))
            return nil
        })
    }

    return g.Wait()
}

func (p *Pipeline) Stats() string {
    return fmt.Sprintf("processed=%d totalBytes=%dMB",
        p.processed.Load(),
        p.totalBytes.Load()>>20,
    )
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    start := time.Now()
    pipeline := NewPipeline(8)
    
    if err := pipeline.Run(ctx, 50); err != nil {
        fmt.Printf("error: %v\n", err)
        return
    }

    elapsed := time.Since(start)
    fmt.Printf("Done in %v\n", elapsed)
    fmt.Printf("Stats: %s\n", pipeline.Stats())
    fmt.Printf("Throughput: %.1f images/sec\n", float64(50)/elapsed.Seconds())
}
```

---

## Step 2 — CPU Profiling (30 min)

### Add pprof to the program:

```go
import (
    "net/http"
    _ "net/http/pprof"
    "runtime"
    "os"
    "runtime/pprof"
)

func main() {
    // CPU profile to file
    f, _ := os.Create("cpu.prof")
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
    
    // ... your program ...
}
```

Or use benchmark:
```bash
go test -bench=BenchmarkPipeline -cpuprofile=cpu.out -benchtime=10s .
go tool pprof cpu.out
```

### In pprof interactive mode:
```
(pprof) top10           # top 10 functions by CPU time
(pprof) top10 -cum      # top 10 by cumulative (includes callees)
(pprof) list resizeImage # show annotated source with CPU time per line
(pprof) web             # open browser flamegraph (requires graphviz)
(pprof) peek resizeImage # show callers and callees
```

### What to look for:
- `runtime.mallocgc` high? → excessive allocation
- `runtime.gcBgMarkWorker` high? → GC pressure
- `sync.(*Mutex).Lock` high? → mutex contention
- Your own function high? → algorithmic optimization opportunity

---

## Step 3 — Memory Profiling (20 min)

```go
f, _ := os.Create("mem.prof")
defer func() {
    runtime.GC()  // force GC for accurate inuse counts
    pprof.WriteHeapProfile(f)
    f.Close()
}()
```

```bash
go tool pprof -alloc_space mem.prof   # where are allocations happening?
go tool pprof -inuse_space mem.prof   # what's live right now?
```

### Questions to answer:
1. Which function allocates the most? (`-alloc_space`)
2. What's the allocation rate (total alloc / time)?
3. Can any allocation be moved to stack or pooled?

### Reducing allocations in the pipeline:

```go
// Before: new buffer every encode
var buf bytes.Buffer
jpeg.Encode(&buf, img, opts)
return buf.Bytes(), nil

// After: pool the buffer
var bufPool = sync.Pool{New: func() any { return new(bytes.Buffer) }}

buf := bufPool.Get().(*bytes.Buffer)
buf.Reset()
defer bufPool.Put(buf)
jpeg.Encode(buf, img, opts)
result := make([]byte, buf.Len())
copy(result, buf.Bytes())  // must copy before returning buf to pool
return result, nil
```

---

## Step 4 — Goroutine and Mutex Contention Profiling (20 min)

### Goroutine profile:
```bash
# While program is running
curl http://localhost:6060/debug/pprof/goroutine?debug=2 > goroutine.txt
# or
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

Look for goroutines stuck in `chan receive`, `sync.Mutex.Lock`, or `sync.WaitGroup.Wait` longer than expected.

### Mutex contention profile:
```go
// Enable mutex profiling (disabled by default — costs ~5% overhead)
runtime.SetMutexProfileFraction(1)  // profile 1 in N mutex events (1 = all)
```

```bash
go test -bench=. -mutexprofile=mutex.out .
go tool pprof mutex.out
(pprof) top   # which mutexes have highest contention?
```

### Block profile (goroutines blocked on channels/mutexes):
```go
runtime.SetBlockProfileRate(1)  // profile all blocking events
```

```bash
go test -bench=. -blockprofile=block.out .
go tool pprof block.out
```

---

## Step 5 — Execution Tracer (45 min)

The execution tracer is the most powerful tool for understanding concurrent behavior. It records every goroutine event.

```go
import "runtime/trace"

f, _ := os.Create("trace.out")
trace.Start(f)
defer trace.Stop()
```

```bash
go test -bench=BenchmarkPipeline -trace=trace.out .
go tool trace trace.out
```

### What the tracer shows:
- **Goroutine analysis**: how many goroutines are running/runnable/waiting at each moment
- **GC events**: when STW happens, how long, which goroutines are paused
- **Network/syscall events**: when goroutines are blocked on I/O
- **Scheduler latency**: time between a goroutine becoming runnable and actually running
- **Per-P view**: what each processor is doing over time

### Using the tracer UI:

1. Open `go tool trace trace.out` in browser
2. Click "Goroutine analysis" → see goroutine count over time
3. Click a specific goroutine → see its full timeline
4. Click "GC" events → see pause duration and frequency
5. Look for:
   - Goroutines spending too long in "runnable" (scheduling latency)
   - GC pauses making goroutines "stop"
   - Uneven distribution of work across Ps (one P idle while others busy)

### Adding trace annotations:

```go
import "runtime/trace"

func processImage(ctx context.Context, img image.Image) ([]byte, error) {
    ctx, task := trace.NewTask(ctx, "processImage")
    defer task.End()

    trace.WithRegion(ctx, "resize", func() {
        // ... resize work ...
    })

    trace.WithRegion(ctx, "encode", func() {
        // ... encode work ...
    })

    return result, nil
}
```

These annotations show up as labeled regions in the trace UI.

---

## Step 6 — GC Tuning Exercise (30 min)

Profile the pipeline's GC behavior and tune it:

```bash
# Baseline: what's the GC doing?
GODEBUG=gctrace=1 go run main.go 2>&1 | grep "^gc"
```

Measure:
1. GC cycle count
2. Average heap size ratio (live/peak)
3. GC CPU%

**Tuning experiments:**

```bash
# Experiment 1: less frequent GC (more memory, less CPU)
GOGC=400 go run main.go

# Experiment 2: memory limit instead of GOGC
GOMEMLIMIT=200MiB GOGC=off go run main.go

# Experiment 3: GOMEMLIMIT + high GOGC (recommended for services)
GOMEMLIMIT=200MiB GOGC=400 go run main.go
```

For each: record throughput (images/sec), peak memory, GC CPU%.

**Expected result:**
- `GOGC=off + GOMEMLIMIT` gives highest throughput for the burst, then GCs hard when limit is hit
- `GOGC=400` gives smoother behavior with more memory headroom
- `GOGC=100` (default) is usually a good baseline; tune based on your memory budget

---

## Step 7 — Benchmark Table and Optimization Decisions (20 min)

Write a comprehensive benchmark with sub-benchmarks:

```go
func BenchmarkPipeline(b *testing.B) {
    configs := []struct {
        name    string
        workers int
        images  int
    }{
        {"workers=1", 1, 20},
        {"workers=4", 4, 20},
        {"workers=8", 8, 20},
        {"workers=16", 16, 20},
    }

    for _, c := range configs {
        c := c
        b.Run(c.name, func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                ctx := context.Background()
                p := NewPipeline(c.workers)
                p.Run(ctx, c.images)
            }
        })
    }
}
```

```bash
go test -bench=BenchmarkPipeline -benchmem -count=3 .
```

**Interpreting results:**
- Does throughput scale linearly with workers? If not, where's the bottleneck?
- Does `allocs/op` decrease with pooling? By how much?
- At what worker count does adding more workers stop helping?

---

## Step 8 — Final Optimization Checklist

Apply what you've learned across all 10 days:

```
[ ] Run with -race: no data races
[ ] Run with goleak: no goroutine leaks
[ ] Check escape analysis: no surprising escapes in hot paths
[ ] Pool large buffers: bytes.Buffer, []byte slices
[ ] Pre-allocate slices: make([]T, 0, cap)
[ ] Context everywhere: all goroutines respect cancellation
[ ] Consistent lock ordering: no potential for deadlock
[ ] GOMAXPROCS tuned: matches physical CPU count (default is fine)
[ ] GOGC/GOMEMLIMIT tuned: based on profiling data
[ ] Benchmarks: before and after for every optimization
[ ] CPU profile: no obvious bottlenecks in hot path
[ ] Memory profile: allocs/op target achieved
[ ] Trace: goroutine scheduling looks even, GC pauses minimal
```

---

## Day 10 Checklist

- [ ] Built the image pipeline capstone program
- [ ] Generated and analyzed CPU profile — identified top CPU consumers
- [ ] Generated and analyzed memory profile — reduced allocs/op with sync.Pool
- [ ] Generated goroutine and mutex contention profiles
- [ ] Used the execution tracer and interpreted goroutine/GC timeline
- [ ] Ran GOGC/GOMEMLIMIT experiments and recorded results
- [ ] Benchmarked with multiple worker counts and identified optimal concurrency

## Key Takeaway

Profiling is the scientific method for performance engineering. Never optimize without measuring first. The tools form a hierarchy:
- `gctrace=1` → is GC a problem?
- `pprof CPU` → what is the CPU spending time on?
- `pprof heap` → where are allocations happening?
- `pprof mutex/block` → where is contention?
- `trace` → how is the scheduler behaving?

Start at the top. Only descend when the higher level confirms a problem. Most performance issues are found at the first or second level.

---

## What's Next After Day 10

You've covered the hard parts. The remaining frontier:

1. **cgo and CGO_ENABLED=0** — calling C from Go and its scheduler implications
2. **unsafe.Pointer** — when and how to use it safely
3. **The runtime source** — `src/runtime/proc.go`, `src/runtime/malloc.go`
4. **netpoller** — how Go handles non-blocking I/O using epoll/kqueue
5. **Delve debugger** — goroutine-aware debugging in production
6. **Continuous profiling** — pprof in production with sampling
