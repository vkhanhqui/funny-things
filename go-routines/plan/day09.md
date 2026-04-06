# Day 09 — Advanced Concurrency Patterns

## Goal
Master production-grade concurrency patterns: context propagation, errgroup, semaphores, singleflight, and circuit breakers.

---

## Step 1 — Read: context.Context Semantics (25 min)

### Rules for context:
1. **Always pass context as the first argument** — never store it in a struct field (exception: HTTP Request)
2. **Never pass nil context** — use `context.Background()` or `context.TODO()` instead
3. **Context cancellation is advisory** — the goroutine receiving it must check and honor it
4. **Don't store mutable values in context** — context values are for request-scoped immutable data (request ID, auth token), NOT for optional function parameters

### The three context constructors:

```go
// Root contexts — use at program entry points
ctx := context.Background()  // non-nil, empty — use in main, tests, top-level servers
ctx := context.TODO()        // non-nil, empty — placeholder when unsure what to use

// Derived contexts
ctx, cancel := context.WithCancel(parent)      // manual cancellation
ctx, cancel := context.WithTimeout(parent, 5*time.Second)  // auto-cancel after timeout
ctx, cancel := context.WithDeadline(parent, time.Now().Add(5*time.Second))  // same as timeout
ctx         := context.WithValue(parent, key, val)  // attach value (no cancel)
```

### Cancel is always a function — always call it:
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()  // ALWAYS defer cancel, even if ctx times out naturally
                // This releases the timer goroutine and context resources
```

### Propagating cancellation correctly:
```go
func fetchAll(ctx context.Context, urls []string) ([]string, error) {
    g, ctx := errgroup.WithContext(ctx)  // derived ctx cancels when any goroutine returns error
    
    results := make([]string, len(urls))
    for i, url := range urls {
        i, url := i, url
        g.Go(func() error {
            req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
            if err != nil {
                return err
            }
            resp, err := http.DefaultClient.Do(req)
            if err != nil {
                return err  // causes all other goroutines' ctx to be cancelled
            }
            defer resp.Body.Close()
            body, _ := io.ReadAll(resp.Body)
            results[i] = string(body)
            return nil
        })
    }
    
    if err := g.Wait(); err != nil {
        return nil, err
    }
    return results, nil
}
```

### context.WithValue — do it right:
```go
// BAD: string key causes collisions across packages
ctx = context.WithValue(ctx, "requestID", id)

// GOOD: unexported type as key prevents collisions
type contextKey struct{}
type requestIDKey struct{}

ctx = context.WithValue(ctx, requestIDKey{}, id)
id := ctx.Value(requestIDKey{}).(string)
```

---

## Step 2 — Read: errgroup (15 min)

`golang.org/x/sync/errgroup` runs a group of goroutines and returns the first non-nil error:

```go
import "golang.org/x/sync/errgroup"

g, ctx := errgroup.WithContext(context.Background())

g.Go(func() error {
    return doThing1(ctx)
})
g.Go(func() error {
    return doThing2(ctx)
})
g.Go(func() error {
    return doThing3(ctx)
})

// Blocks until all goroutines finish; returns first error (if any)
if err := g.Wait(); err != nil {
    log.Fatal(err)
}
```

**Limiting concurrency with errgroup:**
```go
g := new(errgroup.Group)
g.SetLimit(10)  // at most 10 goroutines running concurrently

for _, item := range items {
    item := item
    g.Go(func() error {
        return process(item)
    })
}
g.Wait()
```

---

## Step 3 — Read: Semaphore Pattern (15 min)

A semaphore limits the number of concurrent operations. In Go, implement with a buffered channel or `golang.org/x/sync/semaphore`:

### Buffered channel semaphore:
```go
const maxConcurrent = 10
sem := make(chan struct{}, maxConcurrent)

for _, task := range tasks {
    task := task
    sem <- struct{}{}    // acquire — blocks if 10 goroutines already running
    go func() {
        defer func() { <-sem }()  // release
        process(task)
    }()
}

// Wait for all to finish
for i := 0; i < cap(sem); i++ {
    sem <- struct{}{}
}
```

### weighted semaphore (x/sync):
```go
import "golang.org/x/sync/semaphore"

sem := semaphore.NewWeighted(10)  // max weight 10

for _, task := range tasks {
    task := task
    // Acquire weight 1; blocks until available
    if err := sem.Acquire(ctx, 1); err != nil {
        return err
    }
    go func() {
        defer sem.Release(1)
        process(task)
    }()
}

// Acquire all slots to wait for completion
sem.Acquire(ctx, 10)
```

---

## Step 4 — Read: singleflight (20 min)

`singleflight` collapses duplicate in-flight requests into a single call:

```go
import "golang.org/x/sync/singleflight"

var group singleflight.Group

func fetchUser(ctx context.Context, id string) (*User, error) {
    v, err, shared := group.Do(id, func() (any, error) {
        // Only ONE goroutine executes this for a given id at a time
        // All other goroutines with the same id block and receive this result
        return db.GetUser(ctx, id)
    })
    if err != nil {
        return nil, err
    }
    _ = shared  // true if result was shared with other callers
    return v.(*User), nil
}
```

**When to use singleflight:**
- Cache miss stampede: N goroutines all miss the cache for key X at the same time
- Expensive read-only computation triggered by multiple concurrent requests
- External API calls where duplicate requests are wasteful or rate-limited

**When NOT to use:**
- Writes — singleflight is for reads only (multiple writers to the same key is almost always wrong)
- Operations with different arguments that happen to use the same key

### DoChan — non-blocking variant:
```go
ch := group.DoChan(id, func() (any, error) {
    return db.GetUser(ctx, id)
})

select {
case res := <-ch:
    return res.Val.(*User), res.Err
case <-ctx.Done():
    return nil, ctx.Err()
}
```

---

## Step 5 — Daily Coding Task: URL Fetcher with singleflight + Cache (60 min)

Build a concurrent URL fetcher that:
1. Caches results for 30 seconds
2. Uses `singleflight` to deduplicate in-flight requests for the same URL
3. Limits concurrency to 5 parallel fetches
4. Respects `context.Context` cancellation

```go
package main

import (
    "context"
    "fmt"
    "io"
    "net/http"
    "sync"
    "sync/atomic"
    "time"

    "golang.org/x/sync/errgroup"
    "golang.org/x/sync/singleflight"
)

type CachedFetcher struct {
    mu      sync.RWMutex
    cache   map[string]cacheEntry
    ttl     time.Duration
    group   singleflight.Group
    sem     chan struct{}

    // metrics
    hits   atomic.Int64
    misses atomic.Int64
    deduped atomic.Int64
}

type cacheEntry struct {
    body      string
    fetchedAt time.Time
}

func NewCachedFetcher(ttl time.Duration, maxConcurrent int) *CachedFetcher {
    return &CachedFetcher{
        cache: make(map[string]cacheEntry),
        ttl:   ttl,
        sem:   make(chan struct{}, maxConcurrent),
    }
}

func (f *CachedFetcher) Fetch(ctx context.Context, url string) (string, error) {
    // 1. Check cache
    f.mu.RLock()
    if e, ok := f.cache[url]; ok && time.Since(e.fetchedAt) < f.ttl {
        f.mu.RUnlock()
        f.hits.Add(1)
        return e.body, nil
    }
    f.mu.RUnlock()
    f.misses.Add(1)

    // 2. Deduplicate concurrent fetches for the same URL
    v, err, shared := f.group.Do(url, func() (any, error) {
        // 3. Rate limit concurrent HTTP fetches
        select {
        case f.sem <- struct{}{}:
            defer func() { <-f.sem }()
        case <-ctx.Done():
            return nil, ctx.Err()
        }

        req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
        if err != nil {
            return nil, err
        }
        resp, err := http.DefaultClient.Do(req)
        if err != nil {
            return nil, err
        }
        defer resp.Body.Close()
        body, err := io.ReadAll(resp.Body)
        if err != nil {
            return nil, err
        }

        result := string(body[:min(len(body), 100)])

        // 4. Store in cache
        f.mu.Lock()
        f.cache[url] = cacheEntry{body: result, fetchedAt: time.Now()}
        f.mu.Unlock()

        return result, nil
    })
    
    if shared {
        f.deduped.Add(1)
    }
    if err != nil {
        return "", err
    }
    return v.(string), nil
}

func (f *CachedFetcher) Stats() string {
    return fmt.Sprintf("hits=%d misses=%d deduped=%d",
        f.hits.Load(), f.misses.Load(), f.deduped.Load())
}

func min(a, b int) int {
    if a < b { return a }
    return b
}

func main() {
    fetcher := NewCachedFetcher(30*time.Second, 5)
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    urls := []string{
        "https://httpbin.org/get",
        "https://httpbin.org/get",  // duplicate — should be deduped or cached
        "https://httpbin.org/delay/1",
        "https://httpbin.org/delay/1",
        "https://httpbin.org/status/200",
    }

    g, ctx := errgroup.WithContext(ctx)
    for _, url := range urls {
        url := url
        g.Go(func() error {
            body, err := fetcher.Fetch(ctx, url)
            if err != nil {
                return fmt.Errorf("%s: %w", url, err)
            }
            fmt.Printf("%.50s... (%d chars)\n", url, len(body))
            return nil
        })
    }

    if err := g.Wait(); err != nil {
        fmt.Printf("error: %v\n", err)
    }

    fmt.Println("\nStats:", fetcher.Stats())
    
    // Second round — should all be cache hits
    fmt.Println("\nSecond round (all cached):")
    for _, url := range urls[:3] {
        body, _ := fetcher.Fetch(ctx, url)
        _ = body
    }
    fmt.Println("Stats:", fetcher.Stats())
}
```

---

## Step 6 — Circuit Breaker Pattern (30 min)

A circuit breaker prevents cascading failures by stopping calls to a failing dependency.

**States:**
- **Closed**: requests flow through normally; failures are counted
- **Open**: requests are rejected immediately (fail fast) for a cooldown period
- **Half-Open**: one request is allowed through to test if the dependency recovered

```go
package main

import (
    "errors"
    "sync"
    "time"
)

type State int

const (
    StateClosed   State = iota
    StateOpen
    StateHalfOpen
)

var ErrCircuitOpen = errors.New("circuit breaker is open")

type CircuitBreaker struct {
    mu           sync.Mutex
    state        State
    failures     int
    maxFailures  int
    openedAt     time.Time
    cooldown     time.Duration
}

func NewCircuitBreaker(maxFailures int, cooldown time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        maxFailures: maxFailures,
        cooldown:    cooldown,
    }
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    cb.mu.Lock()
    switch cb.state {
    case StateOpen:
        if time.Since(cb.openedAt) < cb.cooldown {
            cb.mu.Unlock()
            return ErrCircuitOpen
        }
        // Cooldown expired — try half-open
        cb.state = StateHalfOpen
    }
    cb.mu.Unlock()

    err := fn()

    cb.mu.Lock()
    defer cb.mu.Unlock()

    if err != nil {
        cb.failures++
        if cb.state == StateHalfOpen || cb.failures >= cb.maxFailures {
            cb.state = StateOpen
            cb.openedAt = time.Now()
        }
        return err
    }

    // Success — reset
    cb.failures = 0
    cb.state = StateClosed
    return nil
}

func (cb *CircuitBreaker) State() State {
    cb.mu.Lock()
    defer cb.mu.Unlock()
    return cb.state
}
```

---

## Day 09 Checklist

- [ ] Know all 3 context constructors and when to use each
- [ ] Always defer `cancel()` — know why
- [ ] Know the correct unexported key type pattern for context values
- [ ] Used `errgroup.SetLimit` for bounded parallelism
- [ ] Implemented the buffered-channel semaphore
- [ ] Built the cached fetcher with singleflight deduplication
- [ ] Understand circuit breaker state transitions

## Key Takeaway

Production concurrency in Go is built from a few composable primitives: `context` for cancellation propagation, `errgroup` for parallel work with error collection, `singleflight` for deduplication, and a semaphore (buffered channel) for rate limiting. These patterns are not clever tricks — they're the standard way to write robust, leak-free concurrent services. Compose them cleanly, always with context cancellation threading through every layer.
