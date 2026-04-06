# Day 07 — Escape Analysis and Reducing Allocations

## Goal
Understand what "escapes to heap" means, how the compiler decides, and how to reduce allocations in hot paths using escape analysis, sync.Pool, and pre-allocation.

---

## Step 1 — Read: Stack vs. Heap Allocation (30 min)

### The allocation decision

The Go compiler performs **escape analysis** to decide whether a variable lives on the stack or the heap.

**Stack allocation** (preferred):
- Fast: just decrement the stack pointer
- Automatically freed when the function returns — no GC involvement
- Cache-friendly: stack is in CPU cache

**Heap allocation** (escape to heap):
- Slower: requires `mallocgc` + size class lookup
- Increases GC pressure
- Necessary when the variable outlives the function

### What causes escape to heap:

**1. Pointer to local variable returned:**
```go
func newFoo() *Foo {
    f := Foo{x: 1}  // escapes: address taken, returned, outlives function
    return &f
}

func noEscape() Foo {
    f := Foo{x: 1}  // stays on stack: value returned, not pointer
    return f
}
```

**2. Value stored in interface:**
```go
var i interface{} = Foo{x: 1}  // escapes: interface{} requires heap allocation to store concrete type
```

**3. Slice/map/channel backing arrays:**
```go
s := make([]int, 100)  // backing array MAY escape if it outlives the function
```

**4. Variable captured by closure:**
```go
func makeAdder(x int) func(int) int {
    // x escapes: the closure outlives makeAdder
    return func(y int) int { return x + y }
}
```

**5. Large objects:**
```go
var bigArr [65536]byte  // large objects go to heap regardless
```

**6. Dynamic size:**
```go
n := getUserInput()
s := make([]int, n)  // escapes: size unknown at compile time (in many cases)
```

### How to inspect escape analysis:

```bash
go build -gcflags="-m" .          # basic
go build -gcflags="-m -m" .       # verbose with reasons
go build -gcflags="-m=2" .        # level 2 detail
```

Output format:
```
./main.go:12:13: &f escapes to heap
./main.go:8:6: f does not escape
./main.go:20:14: make([]int, n) escapes to heap
```

---

## Step 2 — Experiment: Escape Analysis Deep Dive (45 min)

Create `escape_test.go`:

```go
package main

import "testing"

// Does NOT escape — value type returned
func stackAlloc() int {
    x := 42
    return x
}

// Escapes — pointer returned
func heapAlloc() *int {
    x := 42
    return &x
}

// Does NOT escape (in many cases) — small buffer on stack
func smallSlice() []byte {
    s := make([]byte, 64)
    s[0] = 1
    return s  // wait — does this escape? Try it!
}

// Escapes — returned to caller who might store it
func largeSlice() []byte {
    return make([]byte, 65536)
}

// Interface boxing causes escape
func boxed(v int) interface{} {
    return v  // v escapes to heap due to interface boxing
}

// Closure capture causes escape
func makeClosure() func() int {
    x := 42
    return func() int { return x }  // x escapes
}

func BenchmarkStackAlloc(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = stackAlloc()
    }
}

func BenchmarkHeapAlloc(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = heapAlloc()
    }
}

func BenchmarkSmallSlice(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = smallSlice()
    }
}

func BenchmarkBoxed(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = boxed(42)
    }
}
```

```bash
# See what escapes
go build -gcflags="-m" ./...

# Benchmark to measure impact
go test -bench=. -benchmem .
```

`-benchmem` shows `allocs/op` — that's the key metric for heap pressure.

---

## Step 3 — Read: sync.Pool (20 min)

`sync.Pool` lets you reuse objects to reduce GC pressure. The pool holds objects that would otherwise be garbage collected.

```go
var bufPool = sync.Pool{
    New: func() any {
        return make([]byte, 0, 4096)  // create new 4KB buffer if pool is empty
    },
}

func processRequest(data []byte) {
    buf := bufPool.Get().([]byte)
    defer func() {
        buf = buf[:0]         // reset length but keep capacity
        bufPool.Put(buf)      // return to pool
    }()

    buf = append(buf, data...)
    // ... process buf ...
}
```

### sync.Pool pitfalls:

**1. Pool is cleared on GC** — don't use pool for objects with important cleanup (connections, file handles). They'll be GC'd unexpectedly.

**2. Must reset before Put** — callers share pooled objects. Always reset to zero/empty before returning.

**3. Type assertion required** — `Get()` returns `any`; you must type-assert, or use Go 1.18+ generics.

**4. Don't pool large objects** — the pool's memory benefit is reduced if pooled objects are huge.

**5. Pool doesn't help if concurrent Get() > pool size** — under high load the pool may not have enough objects and `New` will be called frequently.

---

## Step 4 — Daily Coding Task: JSON Handler Optimization (60 min)

Start with a naive HTTP handler that allocates heavily, then optimize it.

**Naive version:**
```go
package main

import (
    "encoding/json"
    "net/http"
    "testing"
)

type Request struct {
    Name  string `json:"name"`
    Value int    `json:"value"`
}

type Response struct {
    Result string `json:"result"`
    Code   int    `json:"code"`
}

func naiveHandler(w http.ResponseWriter, r *http.Request) {
    // Allocates: request body buffer, json decoder, decoded struct, response struct, encoder, response buffer
    var req Request
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), 400)
        return
    }
    resp := Response{
        Result: "processed: " + req.Name,
        Code:   200,
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func BenchmarkNaiveHandler(b *testing.B) {
    // ... benchmark setup
}
```

**Optimized version:**
```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "sync"
)

// Pool reusable decoders and buffers
var (
    decoderPool = sync.Pool{
        New: func() any { return json.NewDecoder(nil) },
    }
    bufPool = sync.Pool{
        New: func() any { return new(bytes.Buffer) },
    }
)

func optimizedHandler(w http.ResponseWriter, r *http.Request) {
    // Reuse decoder
    dec := decoderPool.Get().(*json.Decoder)
    dec.Reset(r.Body)
    defer decoderPool.Put(dec)

    var req Request
    if err := dec.Decode(&req); err != nil {
        http.Error(w, err.Error(), 400)
        return
    }

    // Reuse response buffer
    buf := bufPool.Get().(*bytes.Buffer)
    buf.Reset()
    defer bufPool.Put(buf)

    resp := Response{
        Result: "processed: " + req.Name,
        Code:   200,
    }

    enc := json.NewEncoder(buf)
    if err := enc.Encode(resp); err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))
    buf.WriteTo(w)
}
```

Benchmark both and compare `allocs/op`:
```bash
go test -bench=. -benchmem -count=3 .
```

Target: reduce `allocs/op` by at least 50%.

---

## Step 5 — Pre-allocation Patterns (30 min)

### Slice pre-allocation

```go
// BAD: repeated reallocation
func buildSlice(n int) []int {
    var s []int
    for i := 0; i < n; i++ {
        s = append(s, i)  // may reallocate O(log n) times
    }
    return s
}

// GOOD: pre-allocate known capacity
func buildSliceFast(n int) []int {
    s := make([]int, 0, n)  // length 0, capacity n
    for i := 0; i < n; i++ {
        s = append(s, i)  // never reallocates
    }
    return s
}
```

### String builder

```go
// BAD: O(n²) allocations
func join(words []string) string {
    result := ""
    for _, w := range words {
        result += w + " "  // new string allocation each iteration
    }
    return result
}

// GOOD: strings.Builder reuses a buffer
func joinFast(words []string) string {
    var sb strings.Builder
    sb.Grow(len(words) * 10)  // hint to pre-allocate
    for _, w := range words {
        sb.WriteString(w)
        sb.WriteByte(' ')
    }
    return sb.String()
}
```

### Map pre-sizing

```go
// BAD: map grows incrementally (rehashing allocates new backing array)
m := make(map[string]int)
for _, k := range keys {
    m[k] = 1
}

// GOOD: hint at expected size
m := make(map[string]int, len(keys))
for _, k := range keys {
    m[k] = 1
}
```

### Benchmark pre-allocation:

```go
func BenchmarkAppendNoPrealloc(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = buildSlice(10000)
    }
}

func BenchmarkAppendWithPrealloc(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = buildSliceFast(10000)
    }
}
```

---

## Step 6 — Using pprof to Find Allocation Hotspots (20 min)

```bash
# Run benchmark with memory profile
go test -bench=BenchmarkNaiveHandler -memprofile=mem.out -benchmem .

# Analyze allocations by site
go tool pprof -alloc_space mem.out
(pprof) top10
(pprof) list naiveHandler

# Focus on inuse vs. alloc
go tool pprof -inuse_space mem.out   # current live allocations
go tool pprof -alloc_space mem.out   # total allocations over time
```

**What to look for:**
- Functions with high `alloc_space` but low `inuse_space` = short-lived allocations = GC pressure
- Functions with high `inuse_space` = long-lived data = real memory consumption

---

## Day 07 Checklist

- [ ] Can predict what escapes to heap vs. stays on stack for common patterns
- [ ] Can use `go build -gcflags="-m"` and read escape analysis output
- [ ] Understand `sync.Pool` semantics and all 5 pitfalls
- [ ] Optimized the JSON handler and reduced `allocs/op` by ≥50%
- [ ] Know the pre-allocation patterns for slices, strings, and maps
- [ ] Can use pprof to find allocation hotspots in a benchmark

## Key Takeaway

Every heap allocation costs: malloc time + GC scan time + GC pressure on all goroutines. The compiler's escape analysis is your friend — if something doesn't need to live longer than its function, keep it on the stack. For objects that must be heap-allocated but are frequently created/discarded, `sync.Pool` is the right tool. Always measure with `-benchmem` before and after optimization.
