# slogger

A lightweight, production-grade structured logging library for Go — zero external dependencies, pure standard library.

## Features

- **Levelled logging**: TRACE, DEBUG, INFO, WARN, ERROR, FATAL, PANIC
- **Structured fields**: strongly-typed key-value pairs with minimal heap allocation
- **Two formatters out of the box**: human-readable `TextFormatter` and machine-parseable `JSONFormatter`
- **Context-aware logging**: embed fields in `context.Context` and extract them at log time
- **Hook system**: fire side-effects (metrics, alerts) on matching log levels
- **Per-message sampling**: prevent log storms with configurable rate limiting
- **Child loggers**: `With(fields...)` and `Named(name)` for scoped loggers
- **MultiWriter**: fan output to multiple `io.Writer` destinations
- **Entry pool**: `sync.Pool`-backed entry recycling for hot paths
- **Concurrent-safe**: all exported methods safe for use from multiple goroutines
- **ANSI colour support**: automatic TTY detection with per-level colours

## Installation

```go
import "github.com/sivaosorg/replify/pkg/slogger"
```

No external dependencies — Go standard library only.

## Quick start

```go
log := slogger.New(func(o *slogger.Options) {
    o.Level     = slogger.DebugLevel
    o.Formatter = slogger.NewJSONFormatter()
    o.Output    = os.Stdout
})

log.Info("server started", slogger.String("addr", ":8080"))
log.Warn("slow query",     slogger.Duration("took", 230*time.Millisecond))
log.Error("request failed", slogger.Err(err))
```

## Log Levels

| Constant       | Value | Behaviour            |
|----------------|-------|----------------------|
| `TraceLevel`   | 0     | verbose diagnostics  |
| `DebugLevel`   | 1     | development details  |
| `InfoLevel`    | 2     | operational messages |
| `WarnLevel`    | 3     | potential issues     |
| `ErrorLevel`   | 4     | recoverable errors   |
| `FatalLevel`   | 5     | logs then `os.Exit(1)` |
| `PanicLevel`   | 6     | logs then `panic(msg)` |

Parse a level from a string (case-insensitive):

```go
lvl, err := slogger.ParseLevel("warn")
```

## Structured Fields

```go
slogger.String("key", "value")
slogger.Int("count", 42)
slogger.Int64("id", 123456789)
slogger.Float64("ratio", 3.14)
slogger.Bool("ok", true)
slogger.Err(err)                              // key = "error"
slogger.Time("at", time.Now())
slogger.Duration("elapsed", 500*time.Millisecond)
slogger.Any("meta", anyValue)
```

## Child loggers

```go
// Attach persistent fields
reqLog := log.With(slogger.String("request_id", rid))
reqLog.Info("handler invoked")

// Scope by name (dot-separated hierarchy)
db     := log.Named("db")       // name = "db"
reader := db.Named("reader")    // name = "db.reader"
```

## Context-aware logging

```go
ctx := slogger.WithContextFields(ctx,
    slogger.String("trace_id", traceID),
    slogger.String("span_id",  spanID),
)

// Fields are automatically merged when logging via an entry
log.WithContext(ctx).Info("processing")
```

## Formatters

### TextFormatter (default)

```
2024-01-15T10:30:00Z INFO  server started addr=:8080
```

```go
f := slogger.NewTextFormatter(os.Stderr).
    WithTimeFormat(time.RFC3339).
    WithDisableColors().
    WithEnableCaller()
```

### JSONFormatter

```json
{"ts":"2024-01-15T10:30:00Z","level":"INFO","msg":"server started","addr":":8080"}
```

```go
f := slogger.NewJSONFormatter().
    WithTimeKey("timestamp").
    WithLevelKey("severity").
    WithEnableCaller()
```

## Hooks

```go
type alertHook struct{ fired bool }

func (h *alertHook) Levels() []slogger.Level {
    return []slogger.Level{slogger.ErrorLevel, slogger.FatalLevel}
}

func (h *alertHook) Fire(e *slogger.Entry) error {
    h.fired = true
    // send alert, update metrics, etc.
    return nil
}

log.AddHook(&alertHook{})
```

## Sampling

Prevent log storms by capping identical messages:

```go
log := slogger.New(func(o *slogger.Options) {
    o.SamplingOpts = &slogger.SamplingOptions{
        First:      10,           // always log the first 10 per period
        Period:     time.Second,  // reset window
        Thereafter: 100,          // then log every 100th
    }
})
```

## MultiWriter

```go
w := slogger.NewMultiWriter(os.Stdout, logFile)
log := slogger.New(func(o *slogger.Options) { o.Output = w })
```

## Global logger

```go
slogger.SetGlobalLogger(log)

slogger.Info("application ready")
slogger.Errorf("unexpected status: %d", code)
```

## Dynamic level control

```go
log.SetLevel(slogger.DebugLevel)
log.GetLevel()
log.IsLevelEnabled(slogger.TraceLevel)
```

## Caller reporting

```go
log := slogger.New(func(o *slogger.Options) {
    o.CallerReporter = true
    o.CallerSkip     = 0
})
// Output includes caller=pkg/server/handler.go:42
```

## Thread safety

All `Logger` methods are safe for concurrent use. The entry pool, atomic level
loads, and output mutex ensure correct behaviour under heavy parallelism.
