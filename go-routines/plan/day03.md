# Day 03 — Channel Internals and Select Mechanics

## Goal
Understand how channels are implemented, what makes them block, and how `select` chooses between cases — including subtle behaviors that cause bugs.

---

## Step 1 — Read: hchan Internals (30 min)

Channels in Go are implemented as `hchan` structs in `runtime/chan.go`. Key fields:

```go
type hchan struct {
    qcount   uint           // number of elements in the circular buffer
    dataqsiz uint           // capacity (size of circular buffer)
    buf      unsafe.Pointer // pointer to circular buffer (for buffered channels)
    elemsize uint16
    closed   uint32
    sendx    uint           // send index in buf
    recvx    uint           // recv index in buf
    recvq    waitq          // goroutines waiting to receive
    sendq    waitq          // goroutines waiting to send
    lock     mutex          // protects all fields
}
```

### What happens on a send `ch <- v`:

**Case 1: receiver is waiting in recvq**
- Take the receiver goroutine directly out of recvq
- Copy `v` directly into the receiver's stack (zero-copy!)
- Make the receiver runnable

**Case 2: buffer has space (buffered channel)**
- Copy `v` into `buf[sendx]`
- Increment `sendx`, `qcount`
- Return immediately

**Case 3: no receiver, no buffer space**
- Block the current goroutine
- Add it to `sendq` with a copy of `v`
- Schedule another goroutine

### What happens on a receive `v := <-ch`:

Mirror of send. If a sender is in `sendq`, copy from sender directly. If buffer has data, copy from buffer. Otherwise block.

### Why this design is fast:
- Direct stack-to-stack copy avoids extra allocation
- The lock is held only briefly (not during the actual goroutine execution)
- No heap allocation for the value being sent (value is copied, not boxed)

### Closing a channel:
- Sets `closed = 1`
- Wakes up all waiting receivers (they get the zero value + `ok=false`)
- Wakes up all waiting senders (they panic — send on closed channel)
- Rule: **only the sender should close**; closing from receiver side is almost always a bug

---

## Step 2 — Read: Select Mechanics (20 min)

`select` is not syntactic sugar — it's a runtime operation (`runtime.selectgo`).

### How select works:

1. **Lock all channels** in the select statement (in a fixed order to avoid deadlock)
2. **Scan** all cases for any that can proceed without blocking
3. If multiple cases are ready: **pick one uniformly at random**
4. If no case is ready and there's a `default`: execute `default`
5. If no case is ready and no `default`: block the goroutine on all channels simultaneously
6. When any channel becomes ready: wake up, **unlock all channels**, execute the ready case

### Key behaviors:

**Fairness — random selection:**
```go
ch1 := make(chan int, 100)
ch2 := make(chan int, 100)
// fill both channels
for i := 0; i < 100; i++ { ch1 <- i; ch2 <- i }

hits1, hits2 := 0, 0
for i := 0; i < 100; i++ {
    select {
    case <-ch1: hits1++
    case <-ch2: hits2++
    }
}
// hits1 ≈ hits2 ≈ 50 — not deterministic
```

**Non-blocking send/receive with default:**
```go
select {
case ch <- value:
    // sent successfully
default:
    // channel full or no receiver — don't block
}
```

**Nil channel blocks forever:**
```go
var ch chan int  // nil channel
select {
case <-ch:      // this case is NEVER selected
    fmt.Println("never")
case <-time.After(1 * time.Second):
    fmt.Println("timeout")
}
```
Use case: dynamically disable a select case by setting a channel to nil.

---

## Step 3 — Experiment: Channel Behavior Verification (45 min)

### Experiment A — Direct goroutine handoff

```go
package main

import (
    "fmt"
    "runtime"
)

func main() {
    ch := make(chan int)  // unbuffered

    go func() {
        // receiver goroutine — blocks until sender arrives
        v := <-ch
        fmt.Printf("received %d on goroutine\n", v)
    }()

    runtime.Gosched()  // let receiver goroutine block on ch

    // When we send, the value goes directly to the blocked goroutine's stack
    ch <- 42
    fmt.Println("sent")
}
```

### Experiment B — Prove select randomness

```go
package main

import "fmt"

func main() {
    ch1 := make(chan string, 1)
    ch2 := make(chan string, 1)

    counts := map[string]int{}

    for i := 0; i < 10000; i++ {
        ch1 <- "ch1"
        ch2 <- "ch2"
        select {
        case v := <-ch1:
            counts[v]++
            <-ch2 // drain ch2
        case v := <-ch2:
            counts[v]++
            <-ch1 // drain ch1
        }
    }

    fmt.Printf("ch1 selected: %d times (%.1f%%)\n", counts["ch1"], float64(counts["ch1"])/100)
    fmt.Printf("ch2 selected: %d times (%.1f%%)\n", counts["ch2"], float64(counts["ch2"])/100)
}
```

### Experiment C — Nil channel pattern (dynamic case disable)

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch1 := make(chan int, 5)
    ch2 := make(chan int, 5)

    for i := 0; i < 5; i++ {
        ch1 <- i
        ch2 <- i * 10
    }

    // Drain ch1 first, then fall back to ch2 only
    for {
        select {
        case v, ok := <-ch1:
            if !ok {
                ch1 = nil  // disable this case — nil channel never selects
                fmt.Println("ch1 exhausted, disabling case")
            } else {
                fmt.Printf("ch1: %d\n", v)
            }
        case v, ok := <-ch2:
            if !ok {
                fmt.Println("both channels exhausted")
                return
            }
            fmt.Printf("ch2: %d\n", v)
        case <-time.After(100 * time.Millisecond):
            fmt.Println("timeout — closing ch2")
            close(ch2)
        }
    }
}
```

---

## Step 4 — Daily Coding Task: Fan-out / Fan-in Pipeline (60 min)

Build a 3-stage pipeline:
1. **Generator**: emits work items
2. **Fan-out**: distributes to N parallel workers  
3. **Fan-in**: merges results into a single output channel

```go
package main

import (
    "context"
    "fmt"
    "math/rand"
    "sync"
    "time"
)

// Stage 1: generator
func generate(ctx context.Context, nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            select {
            case out <- n:
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}

// Stage 2: single worker that squares numbers
func square(ctx context.Context, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            // Simulate work
            time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
            select {
            case out <- n * n:
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}

// Fan-out: distribute input channel to N workers
func fanOut(ctx context.Context, in <-chan int, n int) []<-chan int {
    outputs := make([]<-chan int, n)
    for i := 0; i < n; i++ {
        outputs[i] = square(ctx, in)
    }
    return outputs
}

// Fan-in: merge multiple input channels into one output channel
func merge(ctx context.Context, inputs ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup

    output := func(c <-chan int) {
        defer wg.Done()
        for n := range c {
            select {
            case out <- n:
            case <-ctx.Done():
                return
            }
        }
    }

    wg.Add(len(inputs))
    for _, c := range inputs {
        go output(c)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    // Generate numbers 1..20
    nums := make([]int, 20)
    for i := range nums { nums[i] = i + 1 }

    // Pipeline: generate → fan-out to 4 workers → fan-in → print
    in := generate(ctx, nums...)
    workers := fanOut(ctx, in, 4)
    results := merge(ctx, workers...)

    count := 0
    for result := range results {
        fmt.Printf("result: %d\n", result)
        count++
    }
    fmt.Printf("total results: %d\n", count)
}
```

**Note the problem with this fan-out:** all workers read from the same `in` channel, so work is distributed naturally. But the `square` goroutines all compete for the same input — this is the correct fan-out pattern.

**Extend it:**
1. Add a rate limiter stage using `time.Ticker`
2. Add cancellation: cancel after receiving 10 results
3. Run with `-race` and fix any races

---

## Step 5 — Daily Coding Task 2: Token Bucket Rate Limiter (30 min)

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"
)

type RateLimiter struct {
    tokens   chan struct{}
    refillMu sync.Mutex
}

// NewRateLimiter creates a limiter allowing `rate` operations per second
// with a burst capacity of `burst`.
func NewRateLimiter(rate int, burst int) *RateLimiter {
    rl := &RateLimiter{
        tokens: make(chan struct{}, burst),
    }
    // Fill initial burst
    for i := 0; i < burst; i++ {
        rl.tokens <- struct{}{}
    }
    // Refill at `rate` per second
    go func() {
        ticker := time.NewTicker(time.Second / time.Duration(rate))
        defer ticker.Stop()
        for range ticker.C {
            select {
            case rl.tokens <- struct{}{}:
            default:
                // bucket full, discard token
            }
        }
    }()
    return rl
}

// Wait blocks until a token is available or ctx is done.
func (rl *RateLimiter) Wait(ctx context.Context) error {
    select {
    case <-rl.tokens:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

func main() {
    rl := NewRateLimiter(5, 3) // 5 req/sec, burst of 3

    ctx := context.Background()
    start := time.Now()

    for i := 0; i < 15; i++ {
        if err := rl.Wait(ctx); err != nil {
            fmt.Printf("cancelled: %v\n", err)
            break
        }
        fmt.Printf("request %d at %v\n", i, time.Since(start).Round(time.Millisecond))
    }
}
```

---

## Day 03 Checklist

- [ ] Can describe the `hchan` struct fields and what each does
- [ ] Know the 3 paths for a channel send (direct handoff / buffer / block)
- [ ] Know that `select` picks randomly among ready cases
- [ ] Know the nil channel trick for disabling select cases
- [ ] Implemented and tested the fan-out/fan-in pipeline with `-race`
- [ ] Implemented the token bucket rate limiter using channels

## Key Takeaway

Channels are a mutex-protected circular buffer with goroutine wait queues. The most important property: when a receiver is waiting, a send does a direct stack-to-stack copy — no buffer, no extra allocation. `select` fairness is random, not FIFO. Nil channels are a deliberate tool for disabling cases dynamically.
