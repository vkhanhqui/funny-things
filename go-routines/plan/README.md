# Go Mastery in 10 Days

Target: deep expertise in goroutines, Go's memory model, GC tuning, and concurrency safety.
Prereq: ~2 years Go experience. Basics are skipped.

## Structure

| Phase | Days | Focus |
|-------|------|-------|
| 1 | 1–3 | Goroutines & the GMP scheduler |
| 2 | 4–5 | Go Memory Model & the race detector |
| 3 | 6–7 | GC internals & tuning |
| 4 | 8–9 | Concurrency safety patterns |
| 5 | 10  | End-to-end profiling & integration |

## Daily Files

- [Day 01](day01.md) — GMP scheduler internals
- [Day 02](day02.md) — Goroutine lifecycle, stack growth, leak detection
- [Day 03](day03.md) — Channel internals & select mechanics
- [Day 04](day04.md) — Go Memory Model & happens-before
- [Day 05](day05.md) — Atomics, memory ordering, race detector
- [Day 06](day06.md) — GC internals & gctrace
- [Day 07](day07.md) — Escape analysis & reducing allocations
- [Day 08](day08.md) — Mutex, RWMutex, deadlocks, sync.Cond
- [Day 09](day09.md) — Advanced concurrency patterns
- [Day 10](day10.md) — End-to-end profiling & tuning

## Daily Toolkit Cheatsheet

```bash
# Race detector
go test -race ./...
go run -race main.go

# Escape analysis
go build -gcflags="-m" .
go build -gcflags="-m -m" .   # verbose, shows reason

# GC tracing
GODEBUG=gctrace=1 go run main.go

# Scheduler tracing
GODEBUG=schedtrace=1000 GOMAXPROCS=2 go run main.go

# CPU profile
go test -cpuprofile=cpu.out -bench=.
go tool pprof cpu.out

# Heap profile
go test -memprofile=mem.out -bench=.
go tool pprof -alloc_objects mem.out

# Execution tracer
go test -trace=trace.out -bench=.
go tool trace trace.out

# Goroutine leak check
go get go.uber.org/goleak
```

## Key References

- Go Memory Model: https://go.dev/ref/mem
- GC Guide: https://go.dev/doc/gc-guide
- pprof blog: https://go.dev/blog/pprof
- Runtime hacking notes: https://go.dev/src/runtime/HACKING.md
- golang.org/x/sync package: errgroup, singleflight, semaphore
