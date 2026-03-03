# crontask

`crontask` is a production-grade cron and task scheduling sub-package for the
[replify](https://github.com/sivaosorg/replify) ecosystem. It combines the
battle-tested scheduler-loop design of
[robfig/cron](https://github.com/robfig/cron) with the expressive expression
syntax ideas of [gronx](https://github.com/adhocore/gronx), packaged in the
idiomatic style used throughout replify.

---

## Overview & Motivation

Most Go projects eventually need periodic background work — cache warming,
report generation, health checks, data synchronisation. The standard options
are either `time.Ticker` (too low-level) or a third-party cron library (often
too opinionated or too minimal). `crontask` occupies the middle ground:

- **Correct** — schedules fire at the right time across DST transitions, leap
  years, and arbitrary timezones.
- **Expressive** — standard five/six-field cron, `@alias` shortcuts, `@every`
  intervals, step expressions, named month/weekday tokens, and per-job jitter.
- **Observable** — every job exposes live metadata (last run, next run, run
  count, last error) and accepts hook interfaces for metrics and alerting.
- **Extensible** — option-function constructors make it easy to configure
  individual schedulers and jobs without a proliferating API surface.

---

## Installation

```bash
go get github.com/sivaosorg/replify/pkg/crontask
```

Requires Go 1.21 or later.

---

## Feature Comparison

| Feature                         | robfig/cron | gronx | crontask |
|---------------------------------|:-----------:|:-----:|:--------:|
| 5-field standard cron           | ✔           | ✔     | ✔        |
| 6-field with seconds            | ✔           | ✔     | ✔        |
| `@alias` shortcuts              | ✔           | ✔     | ✔        |
| `@every <duration>`             | ✔           | ✗     | ✔        |
| Named month/DOW tokens          | ✔           | ✔     | ✔        |
| Timezone per scheduler          | ✔           | ✗     | ✔        |
| Timezone per expression (TZ=)   | ✗           | ✗     | ✔        |
| Per-job retry + backoff         | ✗           | ✗     | ✔        |
| Per-job execution timeout       | ✗           | ✗     | ✔        |
| Per-job jitter                  | ✗           | ✗     | ✔        |
| Pre/post execution hooks        | ✗           | ✗     | ✔        |
| Success/failure callbacks       | ✗           | ✗     | ✔        |
| Job introspection (next/last)   | ✗           | ✗     | ✔        |
| Context propagation             | ✔           | ✗     | ✔        |
| Graceful shutdown               | ✔           | ✗     | ✔        |
| Expression-only validation      | ✗           | ✔     | ✔        |
| Zero external dependencies      | ✗           | ✔     | ✔*       |

\* `crontask` depends only on the `replify/pkg/randn` package for UUID
generation, which is part of the same module.

---

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/sivaosorg/replify/pkg/crontask"
)

func main() {
    s, err := crontask.New()
    if err != nil {
        log.Fatal(err)
    }

    _, err = s.Register("0 * * * *", func(ctx context.Context) error {
        fmt.Println("every hour:", time.Now())
        return nil
    })
    if err != nil {
        log.Fatal(err)
    }

    s.Start()

    // Shut down cleanly after 30 seconds.
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    if err := s.Shutdown(ctx); err != nil {
        log.Printf("shutdown timed out: %v", err)
    }
}
```

---

## Advanced Usage

### Basic Scheduling

```go
s, _ := crontask.New()

// Every minute.
s.Register("* * * * *", handler)

// At 09:00 on weekdays.
s.Register("0 9 * * 1-5", handler)

// First day of every month at midnight.
s.Register("0 0 1 * *", handler)

// Every 30 seconds (six-field format).
s.Register("*/30 * * * * *", handler)

s.Start()
```

### Alias Expressions

```go
s.Register(crontask.AliasHourly,   handler)  // "@hourly"
s.Register(crontask.AliasDaily,    handler)  // "@daily"
s.Register(crontask.AliasWeekly,   handler)  // "@weekly"
s.Register(crontask.AliasMonthly,  handler)  // "@monthly"
s.Register(crontask.AliasYearly,   handler)  // "@yearly"
s.Register(crontask.AliasWeekdays, handler)  // Monday–Friday at midnight
s.Register(crontask.AliasWeekends, handler)  // Saturday–Sunday at midnight
```

### Interval Expressions

```go
// Fire every 5 minutes, aligned to the Unix epoch.
s.Register("@every 5m", handler)

// Fire every 30 seconds.
s.Register("@every 30s", handler)
```

### Per-Expression Timezone

```go
// Fire at 09:00 New York time on weekdays.
s.Register("TZ=America/New_York 0 9 * * 1-5", handler)
```

### Jitter

Spread load across a fleet by adding a random delay before each execution:

```go
s.Register("@hourly", handler,
    crontask.WithJitter(5*time.Minute), // random delay in [0, 5m)
)
```

### Job Cancellation

```go
id, _ := s.Register("* * * * *", handler, crontask.WithJobName("my-job"))

// Later:
if err := s.Remove(id); err != nil {
    log.Printf("remove: %v", err)
}
```

### Retry and Backoff

```go
s.Register("@hourly", riskyHandler,
    crontask.WithMaxRetries(3),
    crontask.WithBackoff(crontask.ExponentialBackoff(time.Second)),
)
```

### Execution Timeout

```go
s.Register("@minutely", handler,
    crontask.WithTimeout(10*time.Second),
)
```

### Hooks

```go
type MyHooks struct {
    crontask.NoopHooks
}

func (h *MyHooks) OnFailure(_ context.Context, id string, _ time.Duration, err error) {
    log.Printf("ALERT: job %s failed: %v", id, err)
}

s.Register("@daily", handler, crontask.WithHooks(&MyHooks{}))
```

### Context-Based Shutdown

```go
// Register a signal handler.
sigCh := make(chan os.Signal, 1)
signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

s.Start()
<-sigCh

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
if err := s.Shutdown(ctx); err != nil {
    log.Printf("forced shutdown: %v", err)
}
```

### Job Introspection

```go
// List all registered jobs.
for _, job := range s.Jobs() {
    fmt.Printf("%-20s next=%-25v runs=%d err=%v\n",
        job.Name, job.NextRun, job.RunCount, job.LastErr)
}

// Compute the next 5 activation times for a specific job.
runs, err := s.NextRuns(id, time.Now(), 5)
```

### Scheduler-Level Error Handler

```go
s, _ := crontask.New(
    crontask.WithErrorHandler(func(id string, err error) {
        metrics.IncrCounter("cron.failures", map[string]string{"job": id})
    }),
)
```

---

## Expression Syntax Reference

### Standard Fields

```
┌───────────── minute        (0–59)
│ ┌───────────── hour         (0–23)
│ │ ┌───────────── day-of-month (1–31)
│ │ │ ┌───────────── month       (1–12 or jan–dec)
│ │ │ │ ┌───────────── day-of-week (0–7, 0 and 7 = Sunday, or sun–sat)
│ │ │ │ │
* * * * *
```

With optional leading seconds field (requires `crontask.WithSeconds()` or
simply provide six fields):

```
┌─────────────── second       (0–59)
│ ┌─────────────── minute      (0–59)
│ │ ...
* * * * * *
```

### Operators

| Operator | Description                              | Example       |
|----------|------------------------------------------|---------------|
| `*`      | Every value in the range                 | `*`           |
| `n`      | Exact value                              | `5`           |
| `n-m`    | Inclusive range                          | `1-5`         |
| `n-m/s`  | Range with step                          | `0-30/5`      |
| `*/s`    | Every s-th value across the full range   | `*/15`        |
| `a,b,c`  | Comma-separated list                     | `1,15,30`     |

### Aliases

| Alias         | Equivalent            |
|---------------|-----------------------|
| `@yearly`     | `0 0 1 1 *`           |
| `@annually`   | `0 0 1 1 *`           |
| `@monthly`    | `0 0 1 * *`           |
| `@weekly`     | `0 0 * * 0`           |
| `@daily`      | `0 0 * * *`           |
| `@midnight`   | `0 0 * * *`           |
| `@hourly`     | `0 * * * *`           |
| `@minutely`   | `* * * * *`           |
| `@weekdays`   | `0 0 * * 1-5`         |
| `@weekends`   | `0 0 * * 0,6`         |
| `@every <d>`  | Interval (e.g. `5m`)  |

---

## Migration Guide

### From robfig/cron

```go
// robfig/cron
c := cron.New()
c.AddFunc("0 * * * *", func() { /* ... */ })
c.Start()
defer c.Stop()

// crontask
s, _ := crontask.New()
s.Register("0 * * * *", func(ctx context.Context) error {
    /* ... */
    return nil
})
s.Start()
defer s.Stop()
```

Key differences:

- Job functions receive a `context.Context` and return an `error`.
- Use `Shutdown(ctx)` instead of `Stop()` when you need to wait for the loop
  to exit.
- `Register` returns a job ID you can later pass to `Remove`.

### From gronx

gronx is primarily an expression evaluator/validator. If you use it only to
check whether an expression is valid or to compute the next run time, swap it
out directly:

```go
// gronx
g := gronx.New()
if !g.IsValid(expr) { /* ... */ }

// crontask
if err := crontask.Validate(expr); err != nil { /* ... */ }
```

```go
// gronx — next tick
tasker := gronx.NewTasker()
nextTime, _ := tasker.NextTick(expr, false)

// crontask — next tick
sched, _ := crontask.Parse(expr)
nextTime := sched.Next(time.Now())
```

---

## Limitations & Design Trade-offs

- **No persistent storage** — Job state (run count, last error) is kept in
  memory only. If the process restarts, the history is lost. A persistence
  hook interface is the intended extension point.
- **No distributed coordination** — Two instances of the same service will
  both execute scheduled jobs. Use a distributed lock (e.g. Redis `SET NX`)
  inside your job function if exactly-once semantics are required.
- **Goroutine-per-invocation** — Each job dispatch spawns a goroutine. For
  very high-frequency intervals (`@every 1ms`) with slow jobs, goroutine
  accumulation is possible. Use `WithTimeout` to bound execution time.
- **No missed-job detection** — If a job fires while the scheduler is stopped,
  the missed activation is silently skipped. This is consistent with the
  behaviour of both robfig/cron and most Unix cron implementations.
- **DST transitions** — The scheduler uses `time.After` which is monotonic and
  unaffected by wall-clock adjustments. During a DST "spring forward" the
  skipped hour is not replayed; during a "fall back" the extra hour may cause
  one additional activation.

---

## License

This package is part of the replify module and is distributed under the same
license. See the repository root for details.
