# slogger

A **zero-dependency, production-grade structured logging library** for Go — built entirely on the Go standard library.

---

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Installation](#installation)
4. [Quick Start](#quick-start)
5. [Logging Modes](#logging-modes)
6. [Configuration](#configuration)
   - [Log Levels](#log-levels)
   - [Logger Construction](#logger-construction)
   - [Formatters](#formatters)
   - [Output Writers](#output-writers)
   - [Hooks](#hooks)
   - [Log Rotation](#log-rotation)
   - [Color Support](#color-support)
7. [Structured Fields](#structured-fields)
8. [Child Loggers](#child-loggers)
9. [Context-Aware Logging](#context-aware-logging)
10. [Sampling](#sampling)
11. [Global Logger](#global-logger)
12. [Best Practices](#best-practices)
13. [Performance Considerations](#performance-considerations)
14. [Real-World Examples](#real-world-examples)
15. [Options Reference](#options-reference)

---

## Overview

`slogger` provides structured, levelled logging for Go services and tools. It
was designed from first principles to satisfy three non-negotiable requirements
common in production Go deployments:

1. **Zero allocations on the hot path** — log entries are recycled through a
   `sync.Pool`, and field values are stored inline in the `Field` struct to
   avoid escaping to the heap.
2. **Composable, context-aware loggers** — child loggers created with `With`
   and `Named` inherit parent configuration while adding their own fields,
   enabling per-service and per-request scoping without global state.
3. **No external dependencies** — every capability (structured fields, JSON
   output, ANSI colours, file rotation, rate limiting, hooks) is implemented
   against the Go standard library alone.

### Design philosophy

`slogger` was written for the **replify ecosystem** but is fully usable as a
standalone package. Its API surface is deliberately narrow: one `Logger` type,
one `Entry` type, and a small set of typed field constructors. This keeps the
learning curve minimal and makes tooling (linters, static analysis, IDE
completion) maximally effective.

### Why slogger?

| Goal | slogger approach |
|---|---|
| Structured fields | Typed inline Field values, no interface{} boxing on common types |
| JSON output | Single-pass builder, no reflection |
| Coloured terminals | IsTTY detection, ANSI escape codes |
| File rotation | Size- and age-based, ZIP compression, date-bucketed archives |
| Rate limiting | Per-message sliding-window sampler |
| Side-effect hooks | Level-indexed registry called after every write |
| Context propagation | Fields stored in context.Context, merged at log time |

---

## Architecture

```
pkg/slogger/
│
├── constant.go        — All package constants (levels, field types, ANSI codes,
│                        default dirs, file names, timestamp formats, JSON keys)
│
├── type.go            — All type definitions, ordered highest→lowest abstraction
│                        Logger, Entry, CallerInfo, Field, Formatter (interface),
│                        Hook (interface), Hooks, Options, SamplingOptions,
│                        TextFormatter, JSONFormatter, MultiWriter,
│                        RotationOptions, LevelFileWriter, LevelWriterHook,
│                        sampler, rotatingFile, contextKey, entryPool, global
│
├── utilities.go       — All unexported helper functions:
│                        entry pool, trimFile, defaultOptions, newSampler,
│                        rotation helpers, color helpers, JSON/text helpers
│
├── level.go           — Level methods: String, ParseLevel, IsEnabled
├── field.go           — Field constructors: String, Int, Int64, Float64, …
│
├── entry.go           — Entry/CallerInfo methods and accessors
│                        reset, WithContext, Logger, Time, GetLevel, Message,
│                        Fields, Caller, Context, File, Line, Function
│                        Entry-level log methods: Trace … Panic
│
├── formatter.go       — (placeholder; Formatter interface in type.go)
├── formatter_text.go  — TextFormatter: human-readable key=value output
├── formatter_json.go  — JSONFormatter: single-line JSON output
│
├── hook.go            — Hooks registry: NewHooks, Add, Fire, Len
├── writer.go          — MultiWriter, Stdout, Stderr
├── color.go           — IsTTY terminal detection
├── options.go         — WithRotation functional option
├── context.go         — WithContextFields, FieldsFromContext
├── sampling.go        — sampler.allow (rate-limiting logic)
│
├── logger.go          — Logger: New, With, Named, SetLevel, GetLevel,
│                        SetOutput, SetFormatter, AddHook, IsLevelEnabled,
│                        Trace … Panic, Tracef … Panicf, log, logCtx, getCaller
│
├── global.go          — Package-level delegates to a global Logger
│                        SetGlobalLogger, GetGlobalLogger, Trace … Errorf,
│                        GlobalWithContextFields
│
└── rotation.go        — LevelFileWriter, LevelWriterHook:
                         NewLevelFileWriter, WriteLevel, Write, Close, Rotate,
                         NewLevelWriterHook, Levels, Fire
```

### Component deep-dive

#### Logger

`Logger` is the root type. It holds:

- **`level`** (`atomic.Int32`) — the minimum severity threshold; reads/writes
  are lock-free via `sync/atomic`.
- **`formatter`** (`Formatter`) — the serialiser that converts an `Entry` to
  `[]byte`. Protected by `mu` for safe runtime replacement.
- **`output`** (`io.Writer`) — the primary write destination. Protected by `mu`.
- **`hooks`** (`*Hooks`) — the level-indexed registry of side-effect handlers.
- **`fields`** (`[]Field`) — logger-bound fields prepended to every entry.
- **`name`** (`string`) — dot-separated logger name shown in output.
- **`caller`** / **`callerSkip`** — source-location capture settings.
- **`sampling`** (`*sampler`) — optional per-message rate limiter.

All exported `Logger` methods are goroutine-safe. The formatter and output are
double-protected: reading under a read-lock, writing under an exclusive lock.
The level itself uses atomic operations so that `SetLevel` does not block
concurrent log calls.

#### Entry

`Entry` is an in-flight log event. It is obtained from `sync.Pool` via
`acquireEntry` and returned via `releaseEntry` after the entry has been
formatted and written. The pool pre-allocates a `fields` slice of capacity
`defaultEntryFieldCap` to avoid reallocations for the common case of a small
number of fields per entry.

**Important:** callers must not retain an `*Entry` after the logging call that
produced it returns; the entry will be recycled and its fields overwritten.

#### Field

`Field` is a discriminated union: a `FieldType` tag plus a set of value slots
(string, int64, float64, bool, error, time.Time, time.Duration, interface{}).
Primitive types are stored inline, avoiding heap allocation. The `Value()`
method renders any variant as a `string` for formatters that need text output.

#### Formatter

`Formatter` is a single-method interface:

```go
type Formatter interface {
    Format(*Entry) ([]byte, error)
}
```

Two built-in implementations are provided:

| Formatter | Output style | Best for |
|---|---|---|
| `TextFormatter` | `timestamp LEVEL [name] message key=value` | Development, CLI tools |
| `JSONFormatter` | `{"ts":"…","level":"…","msg":"…","key":value}` | Production, log aggregators |

Both formatters are stateless once constructed and are safe for concurrent use.

#### Hook

A `Hook` fires side-effects on matching log levels after the entry has been
written to the primary output. Typical uses include:

- Sending high-severity alerts (Slack, PagerDuty)
- Emitting metrics (Prometheus counters)
- Shipping entries to a remote log aggregator

Hooks are registered per-logger with `Logger.AddHook` and stored in a
level-indexed registry (`Hooks`) protected by a read-write mutex.

#### Rotation

`LevelFileWriter` maintains four open file handles (debug, info, warn, error).
On each write it checks whether the current file exceeds `MaxBytes` or `MaxAge`
and triggers rotation if needed. Rotation:

1. Closes the active file.
2. Creates a date-based archive directory (`logs/archived/2006-01-02/`).
3. Moves or compresses (ZIP) the old file with a timestamp prefix.
4. Opens a fresh file at the original path.

The `LevelWriterHook` bridges the `Hook` interface and `LevelFileWriter`, so
rotation integrates seamlessly with the hook system.

---

## Installation

```bash
go get github.com/sivaosorg/replify
```

Import:

```go
import "github.com/sivaosorg/replify/pkg/slogger"
```

No additional dependencies are required.

---

## Quick Start

```go
package main

import (
    "os"

    "github.com/sivaosorg/replify/pkg/slogger"
)

func main() {
    // Create a logger with JSON output to stdout at DEBUG level.
    log := slogger.New(func(o *slogger.Options) {
        o.Level     = slogger.DebugLevel
        o.Formatter = slogger.NewJSONFormatter()
        o.Output    = os.Stdout
    })

    log.Info("server started",
        slogger.String("addr", ":8080"),
        slogger.String("env", "production"),
    )

    // Structured error logging.
    if err := doSomething(); err != nil {
        log.Error("operation failed",
            slogger.Err(err),
            slogger.String("op", "doSomething"),
        )
    }
}
```

---

## Logging Modes

### Structured logging

Structured logging attaches discrete, typed key-value pairs to every log
entry rather than embedding values inside a free-form string.

```go
log.Info("user login",
    slogger.String("user_id", uid),
    slogger.String("ip",      ip),
    slogger.Duration("latency", latency),
)
```

**JSON output:**
```json
{"ts":"2026-01-15T10:00:00Z","level":"INFO","msg":"user login","user_id":"u42","ip":"10.0.0.1","latency":"3.2ms"}
```

**Advantages:**
- Fields are indexable and queryable in log aggregators (Elasticsearch, Loki, Datadog).
- No string parsing needed downstream; every field has a predictable name and type.
- New fields can be added without changing parsers or dashboards.

**Disadvantages:**
- Slightly more verbose at the call site compared to `fmt.Sprintf`.
- Requires a log aggregator that understands structured logs to realise the full benefit.

### Unstructured (formatted) logging

`slogger` also provides `fmt.Sprintf`-style methods for situations where
structured fields would be cumbersome:

```go
log.Infof("user %q logged in from %s after %v", uid, ip, latency)
```

**Advantages:**
- Familiar pattern for Go developers.
- Concise for one-off messages.

**Disadvantages:**
- Values are embedded in the message string; downstream tooling cannot parse
  individual fields reliably.
- Harder to filter, alert on, or aggregate by specific values.

**Recommendation:** Use structured fields for any value you might want to query
or alert on. Reserve formatted messages for truly free-form diagnostic output
(e.g., startup banners, developer debugging).

---

## Configuration

### Log levels

Seven severity levels are defined, in increasing order:

| Constant | Numeric | Meaning |
|---|---|---|
| `TraceLevel` | 0 | Fine-grained diagnostic output; very verbose |
| `DebugLevel` | 1 | Developer debugging; enabled in test/staging |
| `InfoLevel`  | 2 | General operational messages; default minimum |
| `WarnLevel`  | 3 | Potentially harmful situations |
| `ErrorLevel` | 4 | Errors that do not stop the application |
| `FatalLevel` | 5 | Logs the message then calls `os.Exit(1)` |
| `PanicLevel` | 6 | Logs the message then panics |

Change the minimum level at any time:

```go
// Atomically — safe to call from any goroutine while the logger is in use.
log.SetLevel(slogger.WarnLevel)

// Check at runtime:
if log.IsLevelEnabled(slogger.DebugLevel) {
    log.Debug("expensive debug info", buildExpensiveField())
}
```

Parse a level from a string (e.g., an environment variable):

```go
lvl, err := slogger.ParseLevel(os.Getenv("LOG_LEVEL"))
if err != nil {
    lvl = slogger.InfoLevel
}
log.SetLevel(lvl)
```

### Logger construction

`New` accepts zero or more functional options:

```go
log := slogger.New(
    func(o *slogger.Options) {
        o.Level          = slogger.DebugLevel
        o.Formatter      = slogger.NewJSONFormatter()
        o.Output         = os.Stdout
        o.CallerReporter = true
        o.CallerSkip     = 0
        o.Name           = "my-service"
        o.Fields         = []slogger.Field{
            slogger.String("version", "1.2.3"),
        }
    },
)
```

If `New` is called with no options, the logger uses:
- `InfoLevel`
- `TextFormatter` writing to `os.Stderr`
- No caller capture
- No sampling, no rotation

### Formatters

#### TextFormatter

Human-readable output, ideal for development and CLI tools:

```
2026-01-15T10:00:00Z INFO  [api] server started addr=:8080 env=production
```

```go
f := slogger.NewTextFormatter(os.Stderr).
    WithTimeFormat(time.RFC3339Nano).  // custom timestamp layout
    WithEnableCaller().                // append caller=pkg/foo/bar.go:42
    WithDisableColor().               // disable ANSI codes (e.g. for files)
    WithDisableTimestamp()             // omit timestamp (when infra adds its own)
```

`TextFormatter` automatically detects whether its output writer is a terminal
using `IsTTY`. Color codes are only emitted when connected to a TTY.

#### JSONFormatter

Machine-parseable single-line JSON, ideal for production and log aggregators:

```json
{"ts":"2026-01-15T10:00:00Z","level":"INFO","name":"api","msg":"server started","addr":":8080"}
```

```go
f := slogger.NewJSONFormatter().
    WithTimeFormat(time.RFC3339Nano).
    WithTimeKey("timestamp").    // override default "ts"
    WithLevelKey("severity").    // override default "level"
    WithMessageKey("message").   // override default "msg"
    WithEnableCaller()           // add "caller":"pkg/foo/bar.go:42"
```

### Output writers

Any `io.Writer` can be used as the log destination:

```go
// Write to a file.
f, _ := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
log.SetOutput(f)

// Write to multiple destinations simultaneously.
w := slogger.NewMultiWriter(os.Stdout, f)
log.SetOutput(w)

// Replace output at runtime (safe; protected by an internal mutex).
log.SetOutput(os.Stdout)
```

Convenience helpers:

```go
slogger.Stdout() // returns os.Stdout as io.Writer
slogger.Stderr() // returns os.Stderr as io.Writer
```

### Hooks

Hooks are called after every log write for their registered levels.

```go
// Define a hook.
type errorAlertHook struct {
    // ... alerting client, metrics client, etc.
}

func (h *errorAlertHook) Levels() []slogger.Level {
    return []slogger.Level{slogger.ErrorLevel, slogger.FatalLevel}
}

func (h *errorAlertHook) Fire(e *slogger.Entry) error {
    // Send to alerting system, emit a Prometheus counter, etc.
    // Do NOT retain e after this function returns.
    return nil
}

// Register the hook.
log.AddHook(&errorAlertHook{})
```

**Performance note:** hooks are called synchronously on the logging goroutine.
If a hook performs I/O (HTTP requests, database writes), consider buffering or
dispatching to a background goroutine inside the hook's `Fire` implementation.

### Log rotation

Enable automatic per-level log file rotation:

```go
log := slogger.New(
    slogger.WithRotation(slogger.RotationOptions{
        Dir:      "logs",           // base directory; created if absent
        MaxBytes: 50 * 1024 * 1024, // rotate at 50 MiB
        MaxAge:   24 * time.Hour,   // or after 24 hours, whichever comes first
        Compress: true,             // zip rotated files
    }),
)
```

This creates four files: `logs/debug.log`, `logs/info.log`, `logs/warn.log`,
`logs/error.log`. When a file exceeds the size or age threshold, it is moved to
`logs/archived/2006-01-02/20060102150405_<level>.zip` (or `.log` without
compression).

Trigger rotation manually (e.g. on `SIGHUP`):

```go
// Obtain the writer from the hook registry if needed, or keep a reference:
lfw, _ := slogger.NewLevelFileWriter(opts)
hook := slogger.NewLevelWriterHook(lfw, slogger.NewJSONFormatter())
log.AddHook(hook)

// On SIGHUP:
_ = lfw.Rotate()
```

### Color support

Color is applied automatically when the output writer is a TTY. Colors are
disabled automatically for file writers, pipes, and CI environments.

Force-disable colors:

```go
f := slogger.NewTextFormatter(os.Stderr).WithDisableColor()
```

Color mapping:

| Level | Color |
|---|---|
| TRACE | Cyan |
| DEBUG | Blue |
| INFO  | Green |
| WARN  | Yellow |
| ERROR / FATAL / PANIC | Red |

---

## Structured Fields

All field constructors return a zero-allocation `Field` value:

```go
slogger.String("key", "value")
slogger.Int("count", 42)
slogger.Int64("id", int64(9007199254740993))
slogger.Float64("ratio", 3.14159)
slogger.Bool("ok", true)
slogger.Err(err)                          // key = "error"
slogger.Time("at", time.Now())
slogger.Duration("elapsed", 500*time.Millisecond)
slogger.Any("meta", anyValue)             // uses fmt.Sprintf("%v", val)
```

Fields are merged in this order when an entry is dispatched:
1. Logger-bound fields (from `With` or `Options.Fields`)
2. Context fields (from `WithContextFields`)
3. Call-site fields (passed directly to `Info`, `Error`, etc.)

---

## Child Loggers

### With — add persistent fields

```go
// All entries from reqLog include request_id and user_id.
reqLog := log.With(
    slogger.String("request_id", rid),
    slogger.String("user_id",    uid),
)
reqLog.Info("handler called")
reqLog.Error("validation failed", slogger.Err(err))
```

`With` is safe to call concurrently and returns a new `*Logger` that shares the
parent's hooks, formatter, and output writer but has its own independent field
list.

### Named — hierarchical logger names

```go
dbLog  := log.Named("db")             // name = "db"
rdLog  := dbLog.Named("reader")       // name = "db.reader"
wrLog  := dbLog.Named("writer")       // name = "db.writer"

rdLog.Info("query executed")
// → 2026-01-15T10:00:00Z INFO  [db.reader] query executed
```

Names are rendered in square brackets in `TextFormatter` output and as the
`name` field in `JSONFormatter` output.

---

## Context-Aware Logging

Store fields in a `context.Context` and retrieve them automatically at log time:

```go
// At request entry point (e.g. HTTP middleware):
ctx = slogger.WithContextFields(ctx,
    slogger.String("trace_id",   tid),
    slogger.String("request_id", rid),
    slogger.String("user_id",    uid),
)

// Deep in a service function — no need to thread a logger through call stacks:
log.WithContext(ctx).Info("database query",
    slogger.Duration("latency", elapsed),
)
// → fields include trace_id, request_id, user_id, and latency
```

Fields are appended, not replaced — successive calls to `WithContextFields`
accumulate fields in the context.

Retrieve fields from a context directly:

```go
fields := slogger.FieldsFromContext(ctx)
```

---

## Sampling

Rate-limit identical log messages to prevent log storms during error spikes:

```go
log := slogger.New(func(o *slogger.Options) {
    o.SamplingOpts = &slogger.SamplingOptions{
        First:      10,              // log the first 10 identical messages per second
        Period:     time.Second,     // sliding window duration
        Thereafter: 100,             // then log every 100th message
    }
})
```

Sampling is keyed on the exact message string. Each unique message maintains an
independent bucket so that one chatty message does not suppress others.

Setting `Thereafter` to `0` drops all messages after the first `First` within
the window — useful for suppressing completely repetitive events.

---

## Global Logger

A package-level logger is initialised at program start with production-safe
defaults (`InfoLevel`, `TextFormatter`, `os.Stderr`). All package-level
functions delegate to it:

```go
// Replace the global logger (e.g. during application bootstrap).
slogger.SetGlobalLogger(log)

// Use the global logger.
slogger.Info("application ready")
slogger.Warn("deprecation notice", slogger.String("api", "/v1/users"))
slogger.Errorf("unexpected status: %d", code)

// Store fields in context for the global logger.
ctx = slogger.GlobalWithContextFields(ctx, slogger.String("trace_id", tid))

// Retrieve the global logger when a reference is needed.
gl := slogger.GlobalLogger()
gl.AddHook(myHook)
```

**Recommendation:** replace the global logger once during bootstrap, then use
named or child loggers for all component-level logging.

---

## Best Practices

### Microservices

Always include a `service` or `app` field so entries can be distinguished in
a shared aggregator:

```go
log := slogger.New(func(o *slogger.Options) {
    o.Level     = slogger.InfoLevel
    o.Formatter = slogger.NewJSONFormatter()
    o.Output    = os.Stdout
    o.Fields    = []slogger.Field{
        slogger.String("service", "order-service"),
        slogger.String("version", version),
        slogger.String("env",     os.Getenv("ENV")),
    }
})
```

### Observability and distributed tracing

Use context-based field propagation to carry trace and span IDs without
threading the logger through every function signature:

```go
func handleRequest(ctx context.Context, r *http.Request) {
    ctx = slogger.WithContextFields(ctx,
        slogger.String("trace_id", extractTraceID(r)),
        slogger.String("span_id",  generateSpanID()),
    )

    // Pass ctx, not log, through the call stack.
    processOrder(ctx, order)
}

func processOrder(ctx context.Context, o Order) {
    log.WithContext(ctx).Info("processing order",
        slogger.String("order_id", o.ID),
    )
}
```

### CLI tools

Use `TextFormatter` with colors for interactive output. Disable timestamps when
the terminal already provides context:

```go
log := slogger.New(func(o *slogger.Options) {
    o.Level     = slogger.DebugLevel
    o.Formatter = slogger.NewTextFormatter(os.Stderr).WithDisableTimestamp()
})
```

### Background workers

Use `Named` to scope each worker's logs:

```go
for i := 0; i < numWorkers; i++ {
    workerLog := log.Named(fmt.Sprintf("worker-%d", i))
    go runWorker(ctx, workerLog)
}
```

### High-volume logging

For services that log tens of thousands of entries per second:

1. Use `JSONFormatter` — single-pass builder with no reflection.
2. Enable sampling for chatty message categories.
3. Route to a `MultiWriter` that includes a buffered writer (`bufio.NewWriter`)
   to reduce syscall frequency.
4. Use `With` to bind high-cardinality fields at the service level so they are
   not reallocated per entry.

```go
buf := bufio.NewWriterSize(logFile, 64*1024) // 64 KiB buffer
log.SetOutput(buf)
// Flush periodically or on shutdown:
defer buf.Flush()
```

---

## Performance Considerations

### Memory allocations

- **Entry pooling** — `Entry` objects are recycled through `sync.Pool`.
  At `InfoLevel`, a typical structured log call with three `String` fields
  performs **zero heap allocations**.
- **Inline field storage** — `Field` is a value type; primitive variants
  (string, int64, float64, bool, Duration) are stored directly in the struct
  without boxing.
- **Builder-based formatting** — both `TextFormatter` and `JSONFormatter`
  use `strings.Builder` with a single `[]byte` conversion at the end, avoiding
  intermediate allocations.

### Concurrency model

- `Logger.level` is an `atomic.Int32` — level reads/writes never block.
- `Logger.formatter` and `Logger.output` are protected by a `sync.RWMutex`:
  concurrent reads (format + write) hold a read lock; `SetFormatter`/`SetOutput`
  acquire an exclusive lock.
- `Hooks` uses a separate `sync.RWMutex` — hook dispatch does not block
  concurrent log calls on the same logger.
- `LevelFileWriter` uses a per-instance `sync.Mutex` — file writes are
  serialised per writer, but multiple loggers writing to different files are
  fully concurrent.

### Logger reuse strategies

- **Prefer child loggers over repeated field passing.** Instead of passing
  `slogger.String("request_id", rid)` to every log call, use
  `log.With(slogger.String("request_id", rid))` once at request entry.
- **Avoid `Any` for hot-path fields.** `AnyType` fields use `json.Marshal` or
  `fmt.Sprintf("%v")` which may allocate. Use typed constructors where possible.
- **Reuse `With` loggers across request lifetime.** Assign a request-scoped
  logger to the context or pass it as a function argument rather than
  constructing it repeatedly.

---

## Real-World Examples

### HTTP middleware logging

```go
func LoggingMiddleware(log *slogger.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            rid   := r.Header.Get("X-Request-ID")
            if rid == "" {
                rid = generateID()
            }

            // Bind request fields to a child logger and inject into context.
            reqLog := log.With(
                slogger.String("request_id", rid),
                slogger.String("method",     r.Method),
                slogger.String("path",       r.URL.Path),
                slogger.String("remote_ip",  r.RemoteAddr),
            )
            ctx := slogger.WithContextFields(r.Context(),
                slogger.String("request_id", rid),
            )

            reqLog.Info("request received")
            rw := newResponseWriter(w)
            next.ServeHTTP(rw, r.WithContext(ctx))

            reqLog.Info("request completed",
                slogger.Int("status",   rw.status),
                slogger.Duration("latency", time.Since(start)),
            )
        })
    }
}
```

### Database query logging

```go
func (r *UserRepository) FindByID(ctx context.Context, id string) (*User, error) {
    start := time.Now()
    user, err := r.db.QueryRowContext(ctx, queryFindByID, id).Scan(...)
    elapsed := time.Since(start)

    if elapsed > 100*time.Millisecond {
        log.WithContext(ctx).Warn("slow query",
            slogger.String("query",   "FindByID"),
            slogger.String("user_id", id),
            slogger.Duration("latency", elapsed),
        )
    }

    if err != nil {
        log.WithContext(ctx).Error("query failed",
            slogger.String("query", "FindByID"),
            slogger.Err(err),
        )
        return nil, err
    }

    log.WithContext(ctx).Debug("query ok",
        slogger.String("query", "FindByID"),
        slogger.Duration("latency", elapsed),
    )
    return user, nil
}
```

### Rotation with compression

```go
log := slogger.New(
    slogger.WithRotation(slogger.RotationOptions{
        Dir:      "/var/log/myapp",
        MaxBytes: 100 * 1024 * 1024, // 100 MiB per level file
        MaxAge:   6 * time.Hour,
        Compress: true,
    }),
    func(o *slogger.Options) {
        o.Level     = slogger.InfoLevel
        o.Formatter = slogger.NewJSONFormatter()
        o.Output    = os.Stdout // also write to stdout
    },
)
```

### Multi-destination output

```go
logFile, _ := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

log := slogger.New(func(o *slogger.Options) {
    o.Formatter = slogger.NewJSONFormatter()
    o.Output    = slogger.NewMultiWriter(
        os.Stdout,  // JSON to stdout for container log driver
        logFile,    // JSON to file for local inspection
    )
})
```

### Testing with hooks

Capture log entries in tests without mocking the entire logger:

```go
type captureHook struct {
    mu      sync.Mutex
    entries []*slogger.Entry
}

func (h *captureHook) Levels() []slogger.Level {
    return []slogger.Level{
        slogger.TraceLevel, slogger.DebugLevel, slogger.InfoLevel,
        slogger.WarnLevel, slogger.ErrorLevel,
    }
}

func (h *captureHook) Fire(e *slogger.Entry) error {
    // Copy the fields to avoid use-after-return of the pooled entry.
    h.mu.Lock()
    h.entries = append(h.entries, &slogger.Entry{})
    // Copy relevant data from e into the stored entry.
    _ = e
    h.mu.Unlock()
    return nil
}
```

---

## Options Reference

| Field | Type | Default | Description |
|---|---|---|---|
| `Level` | `Level` | `InfoLevel` | Minimum log level |
| `Formatter` | `Formatter` | `TextFormatter(stderr)` | Entry serialiser |
| `Output` | `io.Writer` | `os.Stderr` | Primary write destination |
| `CallerReporter` | `bool` | `false` | Capture source location |
| `CallerSkip` | `int` | `0` | Additional stack frames to skip |
| `Fields` | `[]Field` | nil | Logger-bound structured fields |
| `Name` | `string` | `""` | Logger name shown in output |
| `SamplingOpts` | `*SamplingOptions` | nil | Rate-limiting configuration |
| `RotationOpts` | `*RotationOptions` | nil | File rotation configuration |

### SamplingOptions

| Field | Type | Description |
|---|---|---|
| `First` | `int` | Number of identical messages always logged per Period |
| `Period` | `time.Duration` | Sliding window for the counter |
| `Thereafter` | `int` | Log every Nth message after First are exhausted; 0 = drop |

### RotationOptions

| Field | Type | Default | Description |
|---|---|---|---|
| `Dir` | `string` | `"logs"` | Base log directory |
| `MaxBytes` | `int64` | `10 MiB` | File size threshold for rotation |
| `MaxAge` | `time.Duration` | `0` (disabled) | Age threshold for rotation |
| `Compress` | `bool` | `false` | ZIP-compress rotated files |
