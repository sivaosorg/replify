package slogger

import (
"context"
"os"
"time"
)

// reset zeros all fields so the Entry can be reused from the pool.
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
func (e *Entry) Logger() *Logger { return e.logger }

// Time returns the entry's timestamp.
func (e *Entry) Time() time.Time { return e.time }

// GetLevel returns the entry's log level.
func (e *Entry) GetLevel() Level { return e.level }

// Message returns the entry's log message.
func (e *Entry) Message() string { return e.message }

// Fields returns the entry's structured fields.
func (e *Entry) Fields() []Field { return e.fields }

// Caller returns the entry's caller information (may be nil).
func (e *Entry) Caller() *CallerInfo { return e.caller }

// Context returns the entry's associated context (may be nil).
func (e *Entry) Context() context.Context { return e.ctx }

// File returns the source file path.
func (c *CallerInfo) File() string { return c.file }

// Line returns the source line number.
func (c *CallerInfo) Line() int { return c.line }

// Function returns the fully-qualified function name.
func (c *CallerInfo) Function() string { return c.function }

// Trace logs a TRACE-level message using the entry's logger and context.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields
func (e *Entry) Trace(msg string, fields ...Field) {
e.logger.logCtx(e.ctx, TraceLevel, msg, fields...)
}

// Debug logs a DEBUG-level message using the entry's logger and context.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields
func (e *Entry) Debug(msg string, fields ...Field) {
e.logger.logCtx(e.ctx, DebugLevel, msg, fields...)
}

// Info logs an INFO-level message using the entry's logger and context.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields
func (e *Entry) Info(msg string, fields ...Field) {
e.logger.logCtx(e.ctx, InfoLevel, msg, fields...)
}

// Warn logs a WARN-level message using the entry's logger and context.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields
func (e *Entry) Warn(msg string, fields ...Field) {
e.logger.logCtx(e.ctx, WarnLevel, msg, fields...)
}

// Error logs an ERROR-level message using the entry's logger and context.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields
func (e *Entry) Error(msg string, fields ...Field) {
e.logger.logCtx(e.ctx, ErrorLevel, msg, fields...)
}

// Fatal logs a FATAL-level message using the entry's logger and context,
// then calls os.Exit(1).
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields
func (e *Entry) Fatal(msg string, fields ...Field) {
e.logger.logCtx(e.ctx, FatalLevel, msg, fields...)
os.Exit(1)
}

// Panic logs a PANIC-level message using the entry's logger and context,
// then panics with msg.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields
func (e *Entry) Panic(msg string, fields ...Field) {
e.logger.logCtx(e.ctx, PanicLevel, msg, fields...)
panic(msg)
}
