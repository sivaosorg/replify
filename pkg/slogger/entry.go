package slogger

import (
	"context"
	"os"
	"time"
)

// CallerInfo holds the runtime source-location data captured when caller
// reporting is enabled on a Logger.
type CallerInfo struct {
	File     string
	Line     int
	Function string
}

// Entry represents a single in-flight log event.
// Entries are pooled internally; callers must not retain an Entry after the
// log call that produced it returns.
type Entry struct {
	Logger  *Logger
	Time    time.Time
	Level   Level
	Message string
	Fields  []Field
	Caller  *CallerInfo
	Context context.Context
}

// reset zeros all fields so the Entry can be reused from the pool.
func (e *Entry) reset() {
	e.Logger = nil
	e.Time = time.Time{}
	e.Level = TraceLevel
	e.Message = ""
	e.Fields = e.Fields[:0]
	e.Caller = nil
	e.Context = nil
}

// WithContext returns a shallow copy of the entry bound to ctx.
//
// Parameters:
//   - `ctx`: the context to attach to the new entry
//
// Returns:
//
// a new *Entry with the same fields as the receiver but with Context set to ctx.
func (e *Entry) WithContext(ctx context.Context) *Entry {
	cp := *e
	cp.Context = ctx
	return &cp
}

// Trace logs a TRACE-level message using the entry's logger and context.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields
func (e *Entry) Trace(msg string, fields ...Field) {
	e.Logger.logCtx(e.Context, TraceLevel, msg, fields...)
}

// Debug logs a DEBUG-level message using the entry's logger and context.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields
func (e *Entry) Debug(msg string, fields ...Field) {
	e.Logger.logCtx(e.Context, DebugLevel, msg, fields...)
}

// Info logs an INFO-level message using the entry's logger and context.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields
func (e *Entry) Info(msg string, fields ...Field) {
	e.Logger.logCtx(e.Context, InfoLevel, msg, fields...)
}

// Warn logs a WARN-level message using the entry's logger and context.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields
func (e *Entry) Warn(msg string, fields ...Field) {
	e.Logger.logCtx(e.Context, WarnLevel, msg, fields...)
}

// Error logs an ERROR-level message using the entry's logger and context.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields
func (e *Entry) Error(msg string, fields ...Field) {
	e.Logger.logCtx(e.Context, ErrorLevel, msg, fields...)
}

// Fatal logs a FATAL-level message using the entry's logger and context,
// then calls os.Exit(1).
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields
func (e *Entry) Fatal(msg string, fields ...Field) {
	e.Logger.logCtx(e.Context, FatalLevel, msg, fields...)
	os.Exit(1)
}

// Panic logs a PANIC-level message using the entry's logger and context,
// then panics with msg.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields
func (e *Entry) Panic(msg string, fields ...Field) {
	e.Logger.logCtx(e.Context, PanicLevel, msg, fields...)
	panic(msg)
}
