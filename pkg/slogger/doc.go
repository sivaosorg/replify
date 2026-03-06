// Package slogger provides a lightweight, production-grade structured logging
// library for Go applications. It is built entirely on the Go standard library
// and imposes no external dependencies.
//
// All Logger methods are safe for concurrent use. Entry objects are pooled
// internally to reduce allocation pressure on hot paths.
//
// # Package Architecture
//
// The package is split across focused files:
//
//   - level.go         — Level type, ParseLevel, IsEnabled
//   - field.go         — Field type and typed constructor functions
//   - entry.go         — Entry, CallerInfo, and entry-level logging methods
//   - internal_pool.go — sync.Pool management for Entry objects
//   - formatter.go     — Formatter interface
//   - formatter_text.go — TextFormatter: human-readable key=value output
//   - formatter_json.go — JSONFormatter: single-line JSON output
//   - hook.go          — Hook interface and Hooks registry
//   - writer.go        — MultiWriter and Stdout/Stderr helpers
//   - color.go         — ANSI colour helpers and IsTTY
//   - options.go       — Options struct and defaultOptions
//   - context.go       — WithContextFields / FieldsFromContext
//   - sampling.go      — SamplingOptions and per-message rate limiting
//   - logger.go        — Logger: core type, New, With, Named, log dispatch
//   - global.go        — Package-level functions delegating to a global Logger
//   - rotation.go      — LevelFileWriter, LevelWriterHook, RotationOptions: per-level log file rotation with ZIP archiving
//   - type.go          — All struct/var definitions: Logger, Entry, CallerInfo, Hooks, sampler, TextFormatter, JSONFormatter, MultiWriter, Options, SamplingOptions, entryPool, global
//
// # Log Levels
//
// slogger.TraceLevel  // most verbose
// slogger.DebugLevel
// slogger.InfoLevel   // default minimum
// slogger.WarnLevel
// slogger.ErrorLevel
// slogger.FatalLevel  // logs then calls os.Exit(1)
// slogger.PanicLevel  // logs then panics
//
// # Creating a Logger
//
// Use New with optional functional options:
//
//	log := slogger.New(func(o *slogger.Options) {
//	    o.Level     = slogger.DebugLevel
//	    o.Formatter = slogger.NewJSONFormatter()
//	    o.Output    = os.Stdout
//	})
//
// # Structured Fields
//
// slogger.String("key", "value")
// slogger.Int("count", 42)
// slogger.Int64("id", 123456789)
// slogger.Float64("ratio", 3.14)
// slogger.Bool("ok", true)
// slogger.Err(err)
// slogger.Time("at", time.Now())
// slogger.Duration("elapsed", 500*time.Millisecond)
// slogger.Any("meta", someStruct)
//
// # Logging
//
// log.Info("server started", slogger.String("addr", ":8080"))
// log.Warn("slow query", slogger.Duration("took", d))
// log.Error("request failed", slogger.Err(err))
//
// Formatted variants:
//
// log.Infof("listening on :%d", port)
//
// # Child Loggers
//
// Attach persistent fields with With:
//
//	req := log.With(slogger.String("request_id", rid))
//	req.Info("handler called")
//
// Scope loggers with Named (dot-separated):
//
//	db := log.Named("db")       // name = "db"
//	rw := db.Named("reader")   // name = "db.reader"
//
// # Context-Aware Logging
//
// Store fields in a context and retrieve them at log time:
//
//	ctx := slogger.WithContextFields(ctx, slogger.String("trace_id", tid))
//	log.WithContext(ctx).Info("processing request")
//
// # Formatters
//
// TextFormatter (default) — human-readable:
//
//	f := slogger.NewTextFormatter(os.Stderr).
//	    WithTimeFormat(time.RFC3339).
//	    WithEnableCaller()
//
// JSONFormatter — machine-parseable:
//
//	f := slogger.NewJSONFormatter().
//	    WithTimeKey("timestamp").
//	    WithEnableCaller()
//
// # Hooks
//
// Hooks fire on matching levels; useful for alerting or metrics:
//
//	type alertHook struct{}
//	func (h *alertHook) Levels() []slogger.Level { return []slogger.Level{slogger.ErrorLevel} }
//	func (h *alertHook) Fire(e *slogger.Entry) error { /* send alert */ return nil }
//
//	log.AddHook(&alertHook{})
//
// # Sampling
//
// Prevent log storms by rate-limiting identical messages:
//
//	log := slogger.New(func(o *slogger.Options) {
//	    o.SamplingOpts = &slogger.SamplingOptions{
//	        First:      10,
//	        Period:     time.Second,
//	        Thereafter: 100,
//	    }
//	})
//
// # Global Logger
//
// A package-level logger is provided for convenience:
//
//	slogger.SetGlobalLogger(log)
//	slogger.Info("application ready")
//	slogger.Errorf("unexpected status: %d", code)
//
// # MultiWriter
//
// Fan output to multiple destinations:
//
//	w := slogger.NewMultiWriter(os.Stdout, logFile)
//
// All methods on Logger are safe for concurrent use.
package slogger
