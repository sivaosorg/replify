package slogger

import (
	"context"
	"sync/atomic"
	"unsafe"
)

// global holds the package-level logger pointer.
// Access is guarded by atomic operations so callers need no external lock.
var global unsafe.Pointer

func init() {
	l := New()
	atomic.StorePointer(&global, unsafe.Pointer(l))
}

// SetGlobalLogger replaces the package-level logger.
//
// Parameters:
//   - `l`: the logger to use for all package-level calls
func SetGlobalLogger(l *Logger) {
	if l == nil {
		return
	}
	atomic.StorePointer(&global, unsafe.Pointer(l))
}

// GetGlobalLogger returns the current package-level logger.
//
// Returns:
//
// the active *Logger used by all package-level functions.
func GetGlobalLogger() *Logger {
	return (*Logger)(atomic.LoadPointer(&global))
}

// Trace logs a TRACE-level message via the global logger.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields
func Trace(msg string, fields ...Field) {
	GetGlobalLogger().log(TraceLevel, msg, fields...)
}

// Debug logs a DEBUG-level message via the global logger.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields
func Debug(msg string, fields ...Field) {
	GetGlobalLogger().log(DebugLevel, msg, fields...)
}

// Info logs an INFO-level message via the global logger.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields
func Info(msg string, fields ...Field) {
	GetGlobalLogger().log(InfoLevel, msg, fields...)
}

// Warn logs a WARN-level message via the global logger.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields
func Warn(msg string, fields ...Field) {
	GetGlobalLogger().log(WarnLevel, msg, fields...)
}

// Error logs an ERROR-level message via the global logger.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields
func Error(msg string, fields ...Field) {
	GetGlobalLogger().log(ErrorLevel, msg, fields...)
}

// Tracef logs a TRACE-level formatted message via the global logger.
//
// Parameters:
//   - `format`: a fmt.Sprintf format string
//   - `args`: arguments for the format string
func Tracef(format string, args ...interface{}) {
	GetGlobalLogger().Tracef(format, args...)
}

// Debugf logs a DEBUG-level formatted message via the global logger.
//
// Parameters:
//   - `format`: a fmt.Sprintf format string
//   - `args`: arguments for the format string
func Debugf(format string, args ...interface{}) {
	GetGlobalLogger().Debugf(format, args...)
}

// Infof logs an INFO-level formatted message via the global logger.
//
// Parameters:
//   - `format`: a fmt.Sprintf format string
//   - `args`: arguments for the format string
func Infof(format string, args ...interface{}) {
	GetGlobalLogger().Infof(format, args...)
}

// Warnf logs a WARN-level formatted message via the global logger.
//
// Parameters:
//   - `format`: a fmt.Sprintf format string
//   - `args`: arguments for the format string
func Warnf(format string, args ...interface{}) {
	GetGlobalLogger().Warnf(format, args...)
}

// Errorf logs an ERROR-level formatted message via the global logger.
//
// Parameters:
//   - `format`: a fmt.Sprintf format string
//   - `args`: arguments for the format string
func Errorf(format string, args ...interface{}) {
	GetGlobalLogger().Errorf(format, args...)
}

// WithContextFields returns a new context that carries the provided fields for
// use with the global logger.
//
// Parameters:
//   - `ctx`: the parent context
//   - `fields`: fields to embed in the returned context
//
// Returns:
//
// a derived context containing the fields.
func GlobalWithContextFields(ctx context.Context, fields ...Field) context.Context {
	return WithContextFields(ctx, fields...)
}
