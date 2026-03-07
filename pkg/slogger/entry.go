package slogger

import (
	"context"
	"os"
	"time"
)

// reset zeros all fields so the Entry can be reused from the pool.
//
// Parameters:
//   - `e`: the entry to reset
func (e *Entry) reset() {
	e.logger = nil
	e.time = time.Time{}
	e.level = TraceLevel
	e.message = ""
	e.fields = e.fields[:0]
	e.caller = nil
	e.ctx = nil
}

// WithContext returns a shallow copy of the entry bound to ctx.
// The original entry is not modified.
//
// Parameters:
//   - `ctx`: the context to attach to the new entry
//
// Returns:
//
// a new *Entry with the same fields as the receiver but with ctx set.
func (e *Entry) WithContext(ctx context.Context) *Entry {
	cp := *e
	cp.ctx = ctx
	return &cp
}

// Logger returns the entry's associated logger.
//
// Returns:
//
// the *Logger that created this entry, or nil for a detached entry.
func (e *Entry) Logger() *Logger { return e.logger }

// Time returns the entry's timestamp, set to the moment the log call was made.
//
// Returns:
//
// the time.Time at which this log entry was created.
func (e *Entry) Time() time.Time { return e.time }

// GetLevel returns the severity level of this log entry.
//
// Returns:
//
// the Level at which this entry was logged.
func (e *Entry) GetLevel() Level { return e.level }

// Message returns the primary log message string.
//
// Returns:
//
// the message string passed to the logging method.
func (e *Entry) Message() string { return e.message }

// Fields returns the structured key-value fields attached to this entry.
// The slice includes logger-bound fields, context fields, and call-site fields,
// in that order.
//
// Returns:
//
// a slice of Field values; may be empty but never nil.
func (e *Entry) Fields() []Field { return e.fields }

// Caller returns the source-location information captured when caller reporting
// is enabled on the parent Logger. Returns nil when caller reporting is disabled.
//
// Returns:
//
// the *CallerInfo describing the log call site, or nil.
func (e *Entry) Caller() *CallerInfo { return e.caller }

// Context returns the context associated with this entry, if any.
// The context is set via Logger.WithContext or Entry.WithContext.
//
// Returns:
//
// the context.Context attached to this entry, or nil.
func (e *Entry) Context() context.Context { return e.ctx }

// File returns the source file path of the log call site.
// The path is trimmed to the last two path segments for brevity.
//
// Returns:
//
// a short file path string such as "pkg/foo/bar.go".
func (c *CallerInfo) File() string { return c.file }

// Line returns the source line number of the log call site.
//
// Returns:
//
// the 1-based line number within the source file.
func (c *CallerInfo) Line() int { return c.line }

// Function returns the fully-qualified function name of the log call site.
//
// Returns:
//
// a string in the form "package.Function" or "package.Type.Method".
func (c *CallerInfo) Function() string { return c.function }

// Trace logs a TRACE-level message using the entry's logger and context.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields appended after the entry's own fields
func (e *Entry) Trace(msg string, fields ...Field) {
	e.logger.logCtx(e.ctx, TraceLevel, msg, fields...)
}

// Debug logs a DEBUG-level message using the entry's logger and context.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields appended after the entry's own fields
func (e *Entry) Debug(msg string, fields ...Field) {
	e.logger.logCtx(e.ctx, DebugLevel, msg, fields...)
}

// Info logs an INFO-level message using the entry's logger and context.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields appended after the entry's own fields
func (e *Entry) Info(msg string, fields ...Field) {
	e.logger.logCtx(e.ctx, InfoLevel, msg, fields...)
}

// Warn logs a WARN-level message using the entry's logger and context.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields appended after the entry's own fields
func (e *Entry) Warn(msg string, fields ...Field) {
	e.logger.logCtx(e.ctx, WarnLevel, msg, fields...)
}

// Error logs an ERROR-level message using the entry's logger and context.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields appended after the entry's own fields
func (e *Entry) Error(msg string, fields ...Field) {
	e.logger.logCtx(e.ctx, ErrorLevel, msg, fields...)
}

// Fatal logs a FATAL-level message using the entry's logger and context,
// then calls os.Exit(1). The exit is unconditional regardless of the
// configured minimum level.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields appended after the entry's own fields
func (e *Entry) Fatal(msg string, fields ...Field) {
	e.logger.logCtx(e.ctx, FatalLevel, msg, fields...)
	os.Exit(1)
}

// Panic logs a PANIC-level message using the entry's logger and context,
// then panics with msg as the panic value.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields appended after the entry's own fields
func (e *Entry) Panic(msg string, fields ...Field) {
	e.logger.logCtx(e.ctx, PanicLevel, msg, fields...)
	panic(msg)
}
