# slogger

A lightweight, production-grade structured logging library for Go — zero external dependencies, pure standard library.

---

## Table of Contents

1. [Overview](#overview)
2. [Architecture Overview](#architecture-overview)
3. [Installation](#installation)
4. [Quick Start](#quick-start)
5. [API Usage Guidelines](#api-usage-guidelines)
   - [Levels](#levels)
   - [Logger Construction — `New`](#logger-construction--new)
   - [Child Loggers — `With` and `Named`](#child-loggers--with-and-named)
   - [Structured Fields](#structured-fields)
   - [Context-Aware Logging](#context-aware-logging)
   - [Formatters](#formatters)
   - [Hooks](#hooks)
   - [Sampling](#sampling)
   - [MultiWriter](#multiwriter)
   - [Log File Rotation](#log-file-rotation)
   - [Global Logger](#global-logger)
   - [Entry Accessors](#entry-accessors)
6. [Options Reference](#options-reference)
7. [Best Practices](#best-practices)
   - [Production Services](#production-services)
   - [Microservices and Distributed Systems](#microservices-and-distributed-systems)
   - [CLI Tools](#cli-tools)
   - [Background Workers](#background-workers)
   - [Observability Integration](#observability-integration)
8. [Real-World Examples](#real-world-examples)
   - [HTTP Middleware Logging](#http-middleware-logging)
   - [Database Layer Logging](#database-layer-logging)
   - [Rotation with Compression](#rotation-with-compression)
   - [Multi-destination Output](#multi-destination-output)
   - [Testing with Hooks](#testing-with-hooks)
9. [Concurrency Model](#concurrency-model)
10. [Performance Notes](#performance-notes)

---

## Overview

`slogger` provides structured, levelled logging for Go applications. It is designed for minimal allocation, predictable behaviour, and zero external dependencies. Every `Logger` method is safe for concurrent use. Log entries are pooled via `sync.Pool` to reduce GC pressure on hot paths.

**Key design goals:**

- Structured fields as first-class citizens — not string interpolation
- Predictable allocation profile via entry pooling
- Composable child loggers for service-oriented architectures
- Pluggable formatters and hook system for observability pipelines
- Per-level file rotation with ZIP archiving for long-running processes

---

## Architecture Overview

```
pkg/slogger/
│
├── type.go             — All struct definitions and package-level vars
│                         (Logger, Entry, CallerInfo, Hooks, Options, …)
│
├── level.go            — Level type, constants, ParseLevel, IsEnabled
├── field.go            — Field type and typed constructors (String, Int, Err, …)
│
├── entry.go            — Entry/CallerInfo methods and accessor functions
│                         reset(), WithContext(), Logger(), Time(), GetLevel(),
│                         Message(), Fields(), Caller(), Context()
│                         File(), Line(), Function()
│
├── internal_pool.go    — sync.Pool management: acquireEntry / releaseEntry
│
├── formatter.go        — Formatter interface
├── formatter_text.go   — TextFormatter: human-readable key=value output
├── formatter_json.go   — JSONFormatter: single-line JSON output
│
├── hook.go             — Hook interface and Hooks registry
├── writer.go           — MultiWriter, Stdout(), Stderr()
├── color.go            — ANSI colour helpers, IsTTY
│
├── options.go          — defaultOptions(), WithRotation()
├── context.go          — WithContextFields(), FieldsFromContext()
├── sampling.go         — per-message rate-limiting (sampler, allow)
│
├── rotation.go         — LevelFileWriter, rotatingFile, LevelWriterHook
│                         RotationOptions, NewLevelFileWriter, Rotate, Close
│
├── logger.go           — Logger: New, With, Named, log dispatch, logCtx
└── global.go           — Package-level functions delegating to global Logger
```

**Data-flow diagram:**

```
Caller
  │
  ▼
Logger.Info(msg, fields...)
  │
  ├─ level check (atomic.Int32)
  ├─ sampler.allow(msg)          ← optional
  │
  ├─ acquireEntry(logger)        ← pool
  │     set time, level, message, ctx
  │     merge logger.fields + ctx fields + call-site fields
  │     getCaller() if enabled
  │
  ├─ Formatter.Format(entry)     ← TextFormatter or JSONFormatter
  ├─ output.Write(data)
  │
  ├─ Hooks.Fire(level, entry)    ← LevelWriterHook, user hooks, …
  │
  └─ releaseEntry(entry)         ← back to pool
```

---

## Installation

```go
import "github.com/sivaosorg/replify/pkg/slogger"
```

No external dependencies — uses only the Go standard library.

---

## Quick Start

```go
package main

import (
    "os"
    "github.com/sivaosorg/replify/pkg/slogger"
)

func main() {
    log := slogger.New(func(o *slogger.Options) {
        o.Level     = slogger.DebugLevel
        o.Output    = os.Stdout
        o.Formatter = slogger.NewTextFormatter(os.Stdout)
        o.Name      = "myapp"
    })

    log.Info("server started", slogger.String("addr", ":8080"))
    log.Debug("config loaded", slogger.Int("workers", 4))
    log.Warn("slow query", slogger.Duration("took", 450*time.Millisecond))
    log.Error("request failed", slogger.Err(err))
}
```

---

## API Usage Guidelines

### Levels

```go
slogger.TraceLevel  // extremely verbose; for low-level tracing
slogger.DebugLevel  // development diagnostics
slogger.InfoLevel   // general operational events (default minimum)
slogger.WarnLevel   // recoverable anomalies
slogger.ErrorLevel  // errors that need attention
slogger.FatalLevel  // logs then calls os.Exit(1)
slogger.PanicLevel  // logs then panics
```

**Parsing levels from config:**

```go
lvl, err := slogger.ParseLevel(os.Getenv("LOG_LEVEL"))
if err != nil {
    lvl = slogger.InfoLevel
}
```

`ParseLevel` is case-insensitive and trims whitespace. It accepts `WARN` and `WARNING` interchangeably.

**Checking if a level is enabled:**

```go
if log.IsLevelEnabled(slogger.DebugLevel) {
    log.Debug("expensive computation result", slogger.Any("val", compute()))
}
```

---

### Logger Construction — `New`

```go
// Purpose: create a ready-to-use Logger with functional options
// When: at application startup; once per component/service
// Pros: zero-value-safe, safe to call concurrently

log := slogger.New(
    func(o *slogger.Options) {
        o.Level          = slogger.InfoLevel
        o.Output         = os.Stdout
        o.Formatter      = slogger.NewJSONFormatter()
        o.Name           = "api"
        o.CallerReporter = true
        o.CallerSkip     = 0
        o.Fields         = []slogger.Field{
            slogger.String("version", buildVersion),
            slogger.String("env", os.Getenv("APP_ENV")),
        }
    },
)
```

Unset fields fall back to safe defaults: `InfoLevel`, `TextFormatter`, `os.Stderr`.

**Dynamically updating level:**

```go
log.SetLevel(slogger.WarnLevel) // atomic; safe to call from any goroutine
```

---

### Child Loggers — `With` and `Named`

Child loggers inherit all parent settings and prepend fields to every entry.

```go
// With — bind persistent fields
reqLog := log.With(
    slogger.String("request_id", rid),
    slogger.String("user_id", uid),
)
reqLog.Info("handler entered")
reqLog.Error("handler failed", slogger.Err(err))
```

```go
// Named — scope logger with a dot-separated name
db  := log.Named("db")           // name = "db"
rdr := db.Named("reader")        // name = "db.reader"
rdr.Info("connection established")
// output: … [db.reader] connection established
```

- `With` and `Named` share the parent's `*Hooks` registry by reference.
- The child's level can be changed independently via `child.SetLevel(...)`.
- Both are safe to call concurrently and the returned logger is safe for concurrent use.

---

### Structured Fields

All field constructors return `slogger.Field` — a typed value stored inline with no heap allocation for most types.

| Constructor | Type | Example |
|---|---|---|
| `String(key, val)` | string | `slogger.String("service", "api")` |
| `Int(key, val)` | int→int64 | `slogger.Int("count", 42)` |
| `Int64(key, val)` | int64 | `slogger.Int64("id", 9999999)` |
| `Float64(key, val)` | float64 | `slogger.Float64("ratio", 0.95)` |
| `Bool(key, val)` | bool | `slogger.Bool("ok", true)` |
| `Err(err)` | error | `slogger.Err(err)` |
| `Time(key, val)` | time.Time | `slogger.Time("at", time.Now())` |
| `Duration(key, val)` | time.Duration | `slogger.Duration("took", d)` |
| `Any(key, val)` | interface{} | `slogger.Any("meta", struct{}{})` |

**`Err`** always uses the key `"error"`. Nil errors produce the value `"<nil>"`.

**`Any`** falls back to `fmt.Sprintf("%v", val)` in text mode, and `json.Marshal` in JSON mode.

---

### Context-Aware Logging

Fields can be embedded into a `context.Context` and extracted automatically at log time.

```go
// Embed fields
ctx := slogger.WithContextFields(ctx,
    slogger.String("trace_id", traceID),
    slogger.String("span_id", spanID),
)

// Extract and log — fields are merged automatically
log.WithContext(ctx).Info("processing request")
```

`WithContext` returns a `*Entry` that delegates to the logger with the given context. It is NOT pooled — use it as a one-shot entry or call the returned entry's `Info/Debug/…` methods directly.

```go
entry := log.WithContext(ctx)
entry.Info("begin")
entry.Warn("slow path", slogger.Duration("took", d))
```

Fields are merged in order: **logger-bound fields → context fields → call-site fields**. This means call-site fields can override context fields if they share a key (formatters do not deduplicate keys — the last value wins in JSON parsers).

---

### Formatters

#### TextFormatter

Human-readable output, ideal for local development and CLI tools.

```
2024-03-01T12:00:00Z INFO  [api] server started addr=:8080
```

```go
f := slogger.NewTextFormatter(os.Stdout).
    WithTimeFormat(time.RFC3339Nano). // default: time.RFC3339
    WithDisableColors().               // disable ANSI codes
    WithDisableTimestamp().            // omit timestamp
    WithEnableCaller()                 // append caller=file:line
```

Colours are automatically disabled when the output is not a TTY.

#### JSONFormatter

Machine-parseable, NDJSON output ideal for log aggregators (Loki, Splunk, Datadog).

```json
{"ts":"2024-03-01T12:00:00Z","level":"INFO","name":"api","msg":"server started","addr":":8080"}
```

```go
f := slogger.NewJSONFormatter().
    WithTimeKey("timestamp").
    WithLevelKey("severity").
    WithMessageKey("message").
    WithNameKey("logger").
    WithCallerKey("source").
    WithEnableCaller()
```

All keys default to sensible values (`ts`, `level`, `msg`, `caller`, `name`) but are fully overridable for compatibility with existing log schemas.

**Custom formatter:** implement the `Formatter` interface:

```go
type Formatter interface {
    Format(*Entry) ([]byte, error)
}
```

Access entry data via accessors: `e.GetLevel()`, `e.Message()`, `e.Fields()`, `e.Time()`, `e.Caller()`.

---

### Hooks

Hooks fire side-effects for matching log levels — useful for error alerting, metrics emission, or secondary outputs.

```go
type Hook interface {
    Levels() []Level
    Fire(*Entry) error
}
```

**Implementation rules:**
- `Levels()` must return a stable slice; it is called only during `AddHook`.
- `Fire(*Entry)` must not retain the `*Entry` after returning — entries are pooled.
- `Fire` must be safe for concurrent use.
- If `Fire` returns an error, the error is collected but does not prevent other hooks from running.

**Example — error alerting hook:**

```go
type PagerHook struct {
    client *pager.Client
}

func (h *PagerHook) Levels() []slogger.Level {
    return []slogger.Level{slogger.ErrorLevel, slogger.FatalLevel}
}

func (h *PagerHook) Fire(e *slogger.Entry) error {
    // Copy what you need before returning
    msg := e.Message()
    return h.client.Alert(msg)
}

log.AddHook(&PagerHook{client: pagerClient})
```

**Hooks are shared** across `With`/`Named` children. Adding a hook to a parent logger affects all derived children.

---

### Sampling

Prevent log storms from high-frequency identical messages.

```go
log := slogger.New(func(o *slogger.Options) {
    o.SamplingOpts = &slogger.SamplingOptions{
        First:      10,           // always log first 10 per period
        Period:     time.Second,  // reset window
        Thereafter: 100,          // then log every 100th; 0 = drop all
    }
})
```

Sampling is keyed by the exact message string. Different messages have independent counters. Counters reset automatically after `Period`. Setting `Thereafter = 0` drops all messages beyond `First`.

---

### MultiWriter

Fan output to multiple destinations simultaneously.

```go
logFile, _ := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

log := slogger.New(func(o *slogger.Options) {
    o.Output = slogger.NewMultiWriter(os.Stdout, logFile)
})
```

`MultiWriter.Write` calls every writer in order, returning the result of the first. All writers are attempted even if an earlier one fails.

---

### Log File Rotation

`LevelFileWriter` writes log entries to separate per-level files with automatic rotation and optional ZIP compression.

```go
// Create a rotating writer
lfw, err := slogger.NewLevelFileWriter(slogger.RotationOptions{
    Dir:      "/var/log/myapp",
    MaxBytes: 50 * 1024 * 1024,   // rotate at 50 MB
    MaxAge:   24 * time.Hour,      // or after 24 h, whichever comes first
    Compress: true,                // zip the rotated file (default: false)
})
if err != nil {
    log.Fatal("cannot create log writer", slogger.Err(err))
}
defer lfw.Close()

// Attach to logger via hook
log.AddHook(slogger.NewLevelWriterHook(lfw, slogger.NewJSONFormatter()))
```

Or use the `WithRotation` functional option:

```go
log := slogger.New(
    func(o *slogger.Options) {
        o.Level = slogger.InfoLevel
    },
    slogger.WithRotation(slogger.RotationOptions{
        Dir:      "logs",
        MaxBytes: 10 * 1024 * 1024,
        Compress: true,
    }),
)
```

**File layout:**

```
logs/
├── debug.log
├── info.log
├── warn.log
├── error.log
└── archived/
    └── 2024-03-01/
        ├── 20240301120000_info.zip
        └── 20240301130000_error.zip
```

**Level routing:**

| Entry Level | File |
|---|---|
| TraceLevel | debug.log |
| DebugLevel | debug.log |
| InfoLevel | info.log |
| WarnLevel | warn.log |
| ErrorLevel | error.log |
| FatalLevel | error.log |
| PanicLevel | error.log |

**Rotation triggers** (checked before each write):
1. `size + incoming > MaxBytes`
2. `time.Since(openedAt) > MaxAge` (when MaxAge > 0)

**Forced rotation:**

```go
if err := lfw.Rotate(); err != nil {
    // handle
}
```

---

### Global Logger

A package-level logger is provided for convenience in small programs or scripts.

```go
// Replace the global logger
slogger.SetGlobalLogger(log)

// Package-level functions delegate to the global logger
slogger.Info("application ready")
slogger.Errorf("unexpected status: %d", code)
slogger.Debug("debug info", slogger.String("key", "val"))
```

**When to use the global logger:**
- Small CLIs or one-file programs
- `init()` functions before a structured logger is available

**When NOT to use the global logger:**
- Libraries (use dependency-injected `*Logger` parameters instead)
- Services with multiple components (use `Named` / `With` loggers per component)

---

### Entry Accessors

`*Entry` is passed to `Hook.Fire` and is safe to read via its accessor methods. **Do not retain the entry after `Fire` returns** — entries are pooled and their fields will be overwritten.

```go
func (h *myHook) Fire(e *slogger.Entry) error {
    // Read what you need immediately
    level  := e.GetLevel()    // Level
    msg    := e.Message()     // string
    fields := e.Fields()      // []Field (slice header copy; data shared)
    ts     := e.Time()        // time.Time
    caller := e.Caller()      // *CallerInfo or nil
    ctx    := e.Context()     // context.Context or nil
    logger := e.Logger()      // *Logger

    if caller != nil {
        fmt.Printf("%s:%d in %s\n", caller.File(), caller.Line(), caller.Function())
    }
    return nil
}
```

---

## Options Reference

| Field | Type | Default | Description |
|---|---|---|---|
| `Level` | `Level` | `InfoLevel` | Minimum level that produces output |
| `Formatter` | `Formatter` | `TextFormatter(stderr)` | Entry serialiser |
| `Output` | `io.Writer` | `os.Stderr` | Primary output destination |
| `CallerReporter` | `bool` | `false` | Capture file/line/func for each entry |
| `CallerSkip` | `int` | `0` | Extra frames to skip (for wrapper libraries) |
| `Fields` | `[]Field` | `nil` | Fields attached to every entry |
| `Name` | `string` | `""` | Logger identifier shown as `[name]` |
| `SamplingOpts` | `*SamplingOptions` | `nil` | Per-message rate limiting config |
| `RotationOpts` | `*RotationOptions` | `nil` | Per-level file rotation config |

---

## Best Practices

### Production Services

```go
// ✓ Use JSON for machine-parseable output
// ✓ Read level from environment for runtime control
// ✓ Attach service-wide fields at construction
// ✓ Use Named loggers per subsystem

lvl, _ := slogger.ParseLevel(os.Getenv("LOG_LEVEL"))

log := slogger.New(func(o *slogger.Options) {
    o.Level     = lvl
    o.Output    = os.Stdout
    o.Formatter = slogger.NewJSONFormatter()
    o.Name      = "myservice"
    o.Fields    = []slogger.Field{
        slogger.String("version", version),
        slogger.String("env", os.Getenv("APP_ENV")),
    }
})

dbLog  := log.Named("db")
apiLog := log.Named("api")
```

### Microservices and Distributed Systems

```go
// ✓ Embed trace/span IDs in context at the service boundary
// ✓ Use WithContext to propagate them automatically

func handleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := slogger.WithContextFields(r.Context(),
        slogger.String("trace_id", r.Header.Get("X-Trace-Id")),
        slogger.String("request_id", newRequestID()),
    )
    processOrder(ctx, log)
}

func processOrder(ctx context.Context, log *slogger.Logger) {
    log.WithContext(ctx).Info("processing order")
    // trace_id and request_id appear automatically
}
```

### CLI Tools

```go
// ✓ Use TextFormatter with colours for interactive output
// ✓ Disable colours when not a TTY (handled automatically)
// ✓ Set level from --verbose flag

log := slogger.New(func(o *slogger.Options) {
    o.Output    = os.Stderr // keep stdout clean for machine output
    o.Formatter = slogger.NewTextFormatter(os.Stderr)
    if verbose {
        o.Level = slogger.DebugLevel
    }
})
```

### Background Workers

```go
// ✓ Use With to bind worker-specific fields
// ✓ Use sampling to handle high-throughput repetitive messages
// ✓ Log durations for every job

workerLog := log.With(slogger.Int("worker_id", id))

for job := range queue {
    start := time.Now()
    err := process(job)
    dur := time.Since(start)
    if err != nil {
        workerLog.Error("job failed",
            slogger.String("job_id", job.ID),
            slogger.Err(err),
            slogger.Duration("took", dur),
        )
    } else {
        workerLog.Info("job complete",
            slogger.String("job_id", job.ID),
            slogger.Duration("took", dur),
        )
    }
}
```

### Observability Integration

```go
// Hook-based Prometheus counter
type metricsHook struct {
    errorTotal prometheus.Counter
}

func (h *metricsHook) Levels() []slogger.Level {
    return []slogger.Level{slogger.ErrorLevel, slogger.FatalLevel}
}

func (h *metricsHook) Fire(e *slogger.Entry) error {
    h.errorTotal.Inc()
    return nil
}

log.AddHook(&metricsHook{errorTotal: errorCounter})
```

---

## Real-World Examples

### HTTP Middleware Logging

```go
func LoggingMiddleware(log *slogger.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            rid   := r.Header.Get("X-Request-Id")
            if rid == "" {
                rid = generateID()
            }

            ctx := slogger.WithContextFields(r.Context(),
                slogger.String("request_id", rid),
                slogger.String("method", r.Method),
                slogger.String("path", r.URL.Path),
            )

            lrw := &loggingResponseWriter{ResponseWriter: w, status: 200}
            next.ServeHTTP(lrw, r.WithContext(ctx))

            log.WithContext(ctx).Info("request complete",
                slogger.Int("status", lrw.status),
                slogger.Duration("took", time.Since(start)),
            )
        })
    }
}
```

### Database Layer Logging

```go
type DB struct {
    pool *sql.DB
    log  *slogger.Logger
}

func NewDB(dsn string, log *slogger.Logger) *DB {
    return &DB{
        pool: mustOpen(dsn),
        log:  log.Named("db"),
    }
}

func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
    start := time.Now()
    rows, err := db.pool.QueryContext(ctx, query, args...)
    dur := time.Since(start)

    if err != nil {
        db.log.WithContext(ctx).Error("query failed",
            slogger.Err(err),
            slogger.Duration("took", dur),
        )
        return nil, err
    }

    if dur > 100*time.Millisecond {
        db.log.WithContext(ctx).Warn("slow query",
            slogger.Duration("took", dur),
        )
    }
    return rows, nil
}
```

### Rotation with Compression

```go
func newProductionLogger(dir string) (*slogger.Logger, *slogger.LevelFileWriter, error) {
    lfw, err := slogger.NewLevelFileWriter(slogger.RotationOptions{
        Dir:      dir,
        MaxBytes: 100 * 1024 * 1024, // 100 MB per level file
        MaxAge:   12 * time.Hour,
        Compress: true,
    })
    if err != nil {
        return nil, nil, fmt.Errorf("log writer: %w", err)
    }

    hook := slogger.NewLevelWriterHook(
        lfw,
        slogger.NewJSONFormatter(),
        slogger.InfoLevel, slogger.WarnLevel, slogger.ErrorLevel,
    )

    log := slogger.New(func(o *slogger.Options) {
        o.Level     = slogger.InfoLevel
        o.Output    = os.Stdout
        o.Formatter = slogger.NewTextFormatter(os.Stdout)
    })
    log.AddHook(hook)

    return log, lfw, nil
}

// In main:
log, lfw, err := newProductionLogger("/var/log/myapp")
if err != nil {
    panic(err)
}
defer lfw.Close()
```

### Multi-destination Output

```go
// Write to stdout and a rolling file simultaneously
logFile, _ := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

log := slogger.New(func(o *slogger.Options) {
    o.Output    = slogger.NewMultiWriter(os.Stdout, logFile)
    o.Formatter = slogger.NewJSONFormatter()
})
```

### Testing with Hooks

```go
type captureHook struct {
    mu      sync.Mutex
    entries []struct {
        level slogger.Level
        msg   string
    }
}

func (h *captureHook) Levels() []slogger.Level {
    return []slogger.Level{
        slogger.TraceLevel, slogger.DebugLevel, slogger.InfoLevel,
        slogger.WarnLevel, slogger.ErrorLevel,
    }
}

func (h *captureHook) Fire(e *slogger.Entry) error {
    // Copy fields BEFORE returning - entry is pooled
    h.mu.Lock()
    h.entries = append(h.entries, struct {
        level slogger.Level
        msg   string
    }{e.GetLevel(), e.Message()})
    h.mu.Unlock()
    return nil
}

func TestMyService(t *testing.T) {
    var buf bytes.Buffer
    hook := &captureHook{}
    log := slogger.New(func(o *slogger.Options) {
        o.Level     = slogger.TraceLevel
        o.Output    = &buf
        o.Formatter = slogger.NewTextFormatter(&buf).WithDisableColors()
    })
    log.AddHook(hook)

    svc := NewService(log)
    svc.DoWork()

    hook.mu.Lock()
    defer hook.mu.Unlock()
    if len(hook.entries) == 0 {
        t.Fatal("expected at least one log entry")
    }
}
```

---

## Concurrency Model

- **`Logger`** fields are protected by a `sync.RWMutex` for `output` and `formatter`. The log level is an `atomic.Int32` for lock-free reads.
- **`Hooks`** uses a `sync.RWMutex`: reads (firing) are shared-lock, writes (adding hooks) are exclusive.
- **Entry pool** uses `sync.Pool` which is goroutine-safe by design.
- **`LevelFileWriter`** has a top-level `sync.Mutex` for writer lookup, and each `rotatingFile` has its own `sync.Mutex` for file I/O, so concurrent writes to different levels proceed in parallel.
- **`sampler`** uses a `sync.Map` for per-message buckets and a `sync.Mutex` per bucket.
- All exported `Logger` methods (including `SetLevel`, `SetOutput`, `SetFormatter`, `AddHook`) are safe to call from multiple goroutines.

---

## Performance Notes

- **Entry pooling**: `sync.Pool` recycles `Entry` values. Entries are initialised with a `fields` slice capacity of 8 to avoid the first few appends.
- **Field constructors**: scalar types (`String`, `Int`, `Bool`, `Float64`) store values inline in the `Field` struct with no heap allocation.
- **Level filtering**: the level check is a single atomic load — entries that don't pass never reach allocation.
- **Sampling**: `sync.Map` ensures lock-free reads for well-established message keys.
- **Text/JSON formatting**: both formatters use `strings.Builder` with no intermediate allocations beyond the final `[]byte`.
- **Rotation**: write-path locking is per-file (by level), so concurrent log calls at different levels do not contend.

---

*All Logger methods are safe for concurrent use from multiple goroutines.*
