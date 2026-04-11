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
func (e *Entry) Logger() *Logger {
	if e == nil {
		return nil
	}
	return e.logger
}

// Time returns the entry's timestamp, set to the moment the log call was made.
//
// Returns:
//
// the time.Time at which this log entry was created.
func (e *Entry) Time() time.Time {
	if e == nil {
		return time.Time{}
	}
	return e.time
}

// Level returns the severity level of this log entry.
//
// Returns:
//
// the Level at which this entry was logged.
func (e *Entry) Level() Level {
	if e == nil {
		return TraceLevel
	}
	return e.level
}

// GetLevel returns the severity level of this log entry.
//
// Deprecated: Use Level() instead. This method will be removed in a future version.
//
// Returns:
//
// the Level at which this entry was logged.
func (e *Entry) GetLevel() Level { return e.Level() }

// Message returns the primary log message string.
//
// Returns:
//
// the message string passed to the logging method.
func (e *Entry) Message() string {
	if e == nil {
		return ""
	}
	return e.message
}

// Fields returns the structured key-value fields attached to this entry.
// The slice includes logger-bound fields, context fields, and call-site fields,
// in that order.
//
// Returns:
//
// a copy of the Field slice; may be empty or nil when the receiver or fields are nil.
func (e *Entry) Fields() []Field {
	if e == nil || e.fields == nil {
		return nil
	}
	result := make([]Field, len(e.fields))
	copy(result, e.fields)
	return result
}

// Caller returns the source-location information captured when caller reporting
// is enabled on the parent Logger. Returns nil when caller reporting is disabled.
//
// Returns:
//
// the *CallerInfo describing the log call site, or nil.
func (e *Entry) Caller() *CallerInfo {
	if e == nil {
		return nil
	}
	return e.caller
}

// Context returns the context associated with this entry, if any.
// The context is set via Logger.WithContext or Entry.WithContext.
//
// Returns:
//
// the context.Context attached to this entry, or nil.
func (e *Entry) Context() context.Context {
	if e == nil {
		return nil
	}
	return e.ctx
}

// File returns the source file path of the log call site.
// The path is trimmed to the last two path segments for brevity.
//
// Returns:
//
// a short file path string such as "pkg/foo/bar.go", or empty string if receiver is nil.
func (c *CallerInfo) File() string {
	if c == nil {
		return ""
	}
	return c.file
}

// Line returns the source line number of the log call site.
//
// Returns:
//
// the 1-based line number within the source file, or 0 if receiver is nil.
func (c *CallerInfo) Line() int {
	if c == nil {
		return 0
	}
	return c.line
}

// Function returns the fully-qualified function name of the log call site.
//
// Returns:
//
// a string in the form "package.Function" or "package.Type.Method", or empty string if receiver is nil.
func (c *CallerInfo) Function() string {
	if c == nil {
		return ""
	}
	return c.function
}

// Trace logs a TRACE-level message using the entry's logger and context.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields appended after the entry's own fields
func (e *Entry) Trace(msg string, fields ...Field) {
	if e.logger == nil {
		return
	}
	e.logger.dispatchContext(e.ctx, TraceLevel, msg, fields...)
}

// Debug logs a DEBUG-level message using the entry's logger and context.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields appended after the entry's own fields
func (e *Entry) Debug(msg string, fields ...Field) {
	if e.logger == nil {
		return
	}
	e.logger.dispatchContext(e.ctx, DebugLevel, msg, fields...)
}

// Info logs an INFO-level message using the entry's logger and context.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields appended after the entry's own fields
func (e *Entry) Info(msg string, fields ...Field) {
	if e.logger == nil {
		return
	}
	e.logger.dispatchContext(e.ctx, InfoLevel, msg, fields...)
}

// Warn logs a WARN-level message using the entry's logger and context.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields appended after the entry's own fields
func (e *Entry) Warn(msg string, fields ...Field) {
	if e.logger == nil {
		return
	}
	e.logger.dispatchContext(e.ctx, WarnLevel, msg, fields...)
}

// Error logs an ERROR-level message using the entry's logger and context.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields appended after the entry's own fields
func (e *Entry) Error(msg string, fields ...Field) {
	if e.logger == nil {
		return
	}
	e.logger.dispatchContext(e.ctx, ErrorLevel, msg, fields...)
}

// Fatal logs a FATAL-level message using the entry's logger and context,
// then calls os.Exit(1). The exit is unconditional regardless of the
// configured minimum level.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields appended after the entry's own fields
func (e *Entry) Fatal(msg string, fields ...Field) {
	if e.logger != nil {
		e.logger.dispatchContext(e.ctx, FatalLevel, msg, fields...)
	}
	os.Exit(1)
}

// Panic logs a PANIC-level message using the entry's logger and context,
// then panics with msg as the panic value.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields appended after the entry's own fields
func (e *Entry) Panic(msg string, fields ...Field) {
	if e.logger != nil {
		e.logger.dispatchContext(e.ctx, PanicLevel, msg, fields...)
	}
	panic(msg)
}
