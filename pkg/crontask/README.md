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
  intervals, step expressions, named month/weekday tokens, business-oriented
  built-in aliases, runtime-registerable custom aliases, and per-job jitter.
- **Observable** — every job exposes live metadata (last run, next run, run
  count, last error) and accepts hook interfaces for metrics and alerting.
- **Extensible** — option-function constructors make it easy to configure
  individual schedulers and jobs without a proliferating API surface.
- **Utility-first** — standalone helpers (`IsValidCronExpr`, `IsDue`,
  `NextRun`, `NextRuns`, `Explain`, `MustParse`) work without a running
  scheduler, making the package useful in CLI tools and test helpers too.

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
| Business-oriented built-in aliases | ✗        | ✗     | ✔        |
| Runtime custom alias registration  | ✗        | ✗     | ✔        |
| Timezone per scheduler          | ✔           | ✗     | ✔        |
| Timezone per expression (TZ=)   | ✗           | ✗     | ✔        |
| Human-readable expression explanation | ✗     | ✗     | ✔        |
| Standalone utility functions    | ✗           | partial | ✔      |
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
s.Register(crontask.AliasHourly,         handler)  // "@hourly"
s.Register(crontask.AliasDaily,          handler)  // "@daily"
s.Register(crontask.AliasWeekly,         handler)  // "@weekly"
s.Register(crontask.AliasMonthly,        handler)  // "@monthly"
s.Register(crontask.AliasYearly,         handler)  // "@yearly"
s.Register(crontask.AliasWeekdays,       handler)  // Monday–Friday at midnight
s.Register(crontask.AliasWeekends,       handler)  // Saturday–Sunday at midnight
// Business-oriented aliases:
s.Register(crontask.AliasBusinessDaily,  handler)  // 09:00 weekdays
s.Register(crontask.AliasBusinessHourly, handler)  // top of each hour 09–17 weekdays
s.Register(crontask.AliasQuarterly,      handler)  // midnight, 1st of each quarter
s.Register(crontask.AliasSemiMonthly,    handler)  // midnight, 1st and 15th
s.Register(crontask.AliasWorkHours,      handler)  // every minute 09–17 weekdays
s.Register(crontask.AliasMarketOpen,     handler)  // 09:30 weekdays
s.Register(crontask.AliasMarketClose,    handler)  // 16:00 weekdays
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

## Standalone Utility Functions

These helpers work without a running scheduler and are safe for concurrent use.

### Validation

```go
// Boolean validity check.
ok := crontask.IsValidCronExpr("0 9 * * 1-5") // true

// Structured error for invalid input.
if err := crontask.ValidateCronExpr("bad expr"); err != nil {
    log.Println(err) // crontask: invalid expression "bad expr": expected 5 or 6 fields, got 2
}
```

### IsDue

```go
// Check whether a schedule is due right now (second-level granularity).
if crontask.IsDue("0 9 * * 1-5", time.Now()) {
    sendDailyReport()
}
```

### NextRun / NextRuns

```go
// First activation after now.
next, err := crontask.NextRun("0 9 * * 1-5", time.Now())

// Next 5 activations.
runs, err := crontask.NextRuns("0 9 * * 1-5", time.Now(), 5)
for _, r := range runs {
    fmt.Println(r.Format(time.RFC3339))
}
```

### MustParse

```go
// Panics for invalid expressions — safe for package-level variables.
var settlement = crontask.MustParse(crontask.AliasMarketClose)

fmt.Println(settlement.Raw())                  // "@marketClose"
fmt.Println(settlement.Next(time.Now()))       // next 16:00 weekday
fmt.Println(settlement.IsDue(time.Now()))      // true at 16:00 on a weekday
fmt.Println(settlement.NextN(time.Now(), 3))   // next 3 closings
```

### Custom Alias Registration

```go
// Register once at startup (e.g. in main or init).
if err := crontask.RegisterAlias("@nightly", "0 2 * * *"); err != nil {
    log.Fatal(err)
}

// Use everywhere the new alias is valid.
s.Register("@nightly", backupHandler)
next, _ := crontask.NextRun("@nightly", time.Now())
```

### Explain

Convert any expression to a natural English description:

```go
desc, _ := crontask.Explain("@every 5m")
// "Every 5 minutes"

desc, _ = crontask.Explain("*/30 * * * * *")
// "Every 30 seconds"

desc, _ = crontask.Explain("0 0 * * 1-5")
// "At 00:00, Monday through Friday"

desc, _ = crontask.Explain("0 9 * * 1-5")
// "At 09:00, Monday through Friday"

desc, _ = crontask.Explain("@marketClose")
// "At 16:00, Monday through Friday"

desc, _ = crontask.Explain("0 0 1 1,4,7,10 *")
// "At 00:00, on the 1st of each month, in January, April, July, and October"
```

---

## Real-World Examples

### Email Report Every Weekday at 9 AM

```go
s.Register("0 9 * * 1-5",
    func(ctx context.Context) error {
        return emailService.SendDailyDigest(ctx)
    },
    crontask.WithJobName("daily-email-report"),
    crontask.WithTimeout(30*time.Second),
    crontask.WithMaxRetries(2),
    crontask.WithBackoff(crontask.ConstantBackoff(5*time.Second)),
)
```

### Cleanup Job Every 6 Hours

```go
s.Register("0 */6 * * *",
    func(ctx context.Context) error {
        return store.PurgeExpiredSessions(ctx)
    },
    crontask.WithJobName("session-cleanup"),
    crontask.WithTimeout(2*time.Minute),
)
```

### Financial Settlement at Market Close

```go
s.Register(crontask.AliasMarketClose,
    func(ctx context.Context) error {
        return settlement.RunDailySettlement(ctx)
    },
    crontask.WithJobName("daily-settlement"),
    crontask.WithTimeout(10*time.Minute),
    crontask.WithMaxRetries(1),
    crontask.WithBackoff(crontask.ConstantBackoff(30*time.Second)),
)
```

### Retry Webhook Every 2 Minutes with Backoff

```go
s.Register("@every 2m",
    func(ctx context.Context) error {
        return webhookClient.DeliverPending(ctx)
    },
    crontask.WithJobName("webhook-retry"),
    crontask.WithMaxRetries(5),
    crontask.WithBackoff(crontask.ExponentialBackoff(time.Second)),
    crontask.WithTimeout(30*time.Second),
)
```

### Health Check Every 30 Seconds

```go
s.Register("@every 30s",
    func(ctx context.Context) error {
        return healthChecker.Ping(ctx)
    },
    crontask.WithJobName("health-check"),
    crontask.WithTimeout(5*time.Second),
)
```

### Business-Only Scheduling (Weekdays Only)

```go
// Use the built-in alias.
s.Register(crontask.AliasBusinessHourly,
    func(ctx context.Context) error {
        return metrics.CollectBusinessHourStats(ctx)
    },
    crontask.WithJobName("business-hour-metrics"),
)

// Or with an explicit TZ to align with a specific office.
s.Register("TZ=Europe/London 0 9-17 * * 1-5",
    func(ctx context.Context) error {
        return alerts.CheckSLABreaches(ctx)
    },
    crontask.WithJobName("sla-check-london"),
)
```

### Using Explain to Log the Schedule

```go
exprs := []string{"@businessDaily", "@every 30s", "0 0 1 1,4,7,10 *"}
for _, expr := range exprs {
    desc, err := crontask.Explain(expr)
    if err != nil {
        log.Printf("invalid expression %q: %v", expr, err)
        continue
    }
    log.Printf("registered job: %s (%s)", expr, desc)
}
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

With optional leading seconds field:

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

### Built-in Aliases

| Alias              | Equivalent                | Description                            |
|--------------------|---------------------------|----------------------------------------|
| `@yearly`          | `0 0 1 1 *`               | Once a year, midnight January 1st      |
| `@annually`        | `0 0 1 1 *`               | Synonym for @yearly                    |
| `@monthly`         | `0 0 1 * *`               | Midnight on the 1st each month         |
| `@weekly`          | `0 0 * * 0`               | Midnight on Sunday                     |
| `@daily`           | `0 0 * * *`               | Midnight every day                     |
| `@midnight`        | `0 0 * * *`               | Synonym for @daily                     |
| `@hourly`          | `0 * * * *`               | Top of each hour                       |
| `@minutely`        | `* * * * *`               | Every minute                           |
| `@weekdays`        | `0 0 * * 1-5`             | Midnight, Monday–Friday                |
| `@weekends`        | `0 0 * * 0,6`             | Midnight, Saturday and Sunday          |
| `@businessDaily`   | `0 9 * * 1-5`             | 09:00 weekdays                         |
| `@businessHourly`  | `0 9-17 * * 1-5`          | Top of each hour, 09–17, weekdays      |
| `@quarterly`       | `0 0 1 1,4,7,10 *`        | Midnight, first day of each quarter    |
| `@semiMonthly`     | `0 0 1,15 * *`            | Midnight on the 1st and 15th           |
| `@workhours`       | `* 9-17 * * 1-5`          | Every minute 09:00–17:59, weekdays     |
| `@marketOpen`      | `30 9 * * 1-5`            | 09:30 weekdays (US market open)        |
| `@marketClose`     | `0 16 * * 1-5`            | 16:00 weekdays (US market close)       |
| `@every <d>`       | Interval (e.g. `5m`)      | Fires every duration aligned to epoch  |

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
if !crontask.IsValidCronExpr(expr) { /* ... */ }
// or
if err := crontask.ValidateCronExpr(expr); err != nil { /* ... */ }
```

```go
// gronx — next tick
tasker := gronx.NewTasker()
nextTime, _ := tasker.NextTick(expr, false)

// crontask — next tick
nextTime, err := crontask.NextRun(expr, time.Now())
```

---

## Limitations & Design Trade-offs

### No Persistent Storage

**Current behaviour:** Job state (run count, last error) is kept in memory
only. If the process restarts, the history is lost.

**Mitigation:** Define a persistence adapter around the existing `Hooks`
interface:

```go
// PersistenceAdapter is a conceptual sketch — implement using your preferred
// storage backend (Redis, Postgres, BoltDB, etc.).
type PersistenceAdapter struct {
    crontask.NoopHooks
    db *sql.DB
}

func (p *PersistenceAdapter) OnSuccess(_ context.Context, id string, d time.Duration) {
    _, _ = p.db.Exec(
        "UPDATE cron_jobs SET last_run=NOW(), last_err=NULL, run_count=run_count+1 WHERE id=$1", id)
}

func (p *PersistenceAdapter) OnFailure(_ context.Context, id string, _ time.Duration, err error) {
    _, _ = p.db.Exec(
        "UPDATE cron_jobs SET last_run=NOW(), last_err=$2, run_count=run_count+1 WHERE id=$1", id, err.Error())
}

// Attach at registration time.
s.Register("@daily", handler, crontask.WithHooks(&PersistenceAdapter{db: db}))
```

At process start, read the persisted metadata and use it to decide whether
to skip the first activation (catch-up avoidance) or run immediately.

---

### No Distributed Coordination

**Current behaviour:** Two instances of the same service will both execute
scheduled jobs. This can cause duplicate work or data corruption.

**Mitigation 1 — distributed lock inside the job:**

```go
s.Register("@hourly", func(ctx context.Context) error {
    // Acquire a 60-second lease via Redis SET NX EX.
    acquired, err := redisClient.SetNX(ctx, "lock:hourly-report", "1", time.Minute).Result()
    if err != nil || !acquired {
        return nil // another instance has the lock; skip quietly
    }
    defer redisClient.Del(ctx, "lock:hourly-report")

    return generateReport(ctx)
})
```

**Mitigation 2 — lock middleware via Hooks:**

```go
type DistributedLockHook struct {
    crontask.NoopHooks
    redis  *redis.Client
    locked map[string]bool
    mu     sync.Mutex
}

func (h *DistributedLockHook) OnStart(ctx context.Context, id string) {
    acquired, _ := h.redis.SetNX(ctx, "lock:"+id, "1", 2*time.Minute).Result()
    h.mu.Lock()
    h.locked[id] = acquired
    h.mu.Unlock()
    // If not acquired, the job function should check ctx.Done() or the job
    // should return early by convention.
}

func (h *DistributedLockHook) OnComplete(_ context.Context, id string, _ time.Duration) {
    h.mu.Lock()
    if h.locked[id] {
        h.redis.Del(context.Background(), "lock:"+id)
    }
    delete(h.locked, id)
    h.mu.Unlock()
}
```

---

### Goroutine-per-Invocation

**Current behaviour:** Each job dispatch spawns a goroutine. For very
high-frequency intervals with slow jobs, goroutine accumulation is possible.

**Mitigation 1 — always set a timeout:**

```go
s.Register("@every 1s", handler,
    crontask.WithTimeout(500*time.Millisecond), // bound execution
)
```

**Mitigation 2 — semaphore inside the job function:**

```go
sem := make(chan struct{}, 4) // allow at most 4 concurrent executions

s.Register("@every 1s", func(ctx context.Context) error {
    select {
    case sem <- struct{}{}:
        defer func() { <-sem }()
    default:
        return nil // drop this tick; already at capacity
    }
    return doWork(ctx)
})
```

**Best practice:** use `@every` intervals no smaller than the expected
job duration, and always set `WithTimeout` to prevent goroutine leaks from
stuck jobs.

---

### No Missed-Job Detection

**Current behaviour:** If a job fires while the scheduler is stopped (e.g.
during a deployment), the missed activation is silently skipped. This matches
the behaviour of Unix cron and robfig/cron.

**Mitigation — external catch-up pattern:**

```go
// At startup, check whether the last recorded run is older than one period.
// If so, run the job immediately before starting the scheduler.

lastRun := db.GetLastRun("daily-report")
if time.Since(lastRun) > 24*time.Hour {
    if err := dailyReportJob(ctx); err != nil {
        log.Printf("catch-up run failed: %v", err)
    }
}

s.Register("@daily", dailyReportJob)
s.Start()
```

---

### DST Transitions

**Current behaviour:** The scheduler uses `time.After` which is based on the
monotonic clock and is therefore unaffected by wall-clock adjustments:

- **Spring forward (clocks skip an hour):** any activation whose scheduled
  time falls in the skipped hour is silently missed for that day.
- **Fall back (clocks repeat an hour):** the repeated hour may cause one
  additional activation compared to a typical day.

**Example — illustrating the spring-forward case:**

```go
// A job scheduled for 02:30 in a timezone that observes DST will not fire
// on the night clocks move from 02:00 → 03:00 (the time never exists).
tz, _ := time.LoadLocation("America/New_York")
s, _ := crontask.New(crontask.WithLocation(tz))
s.Register("30 2 * * *", handler) // silently skipped on DST spring-forward night
```

**Recommendation:** Avoid scheduling jobs at wall-clock times that fall
inside DST transition windows (typically 00:00–03:00 in the affected timezone).
Use UTC for high-reliability schedules.

---

## License

This package is part of the replify module and is distributed under the same
license. See the repository root for details.

