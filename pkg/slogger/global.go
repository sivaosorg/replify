package slogger

import (
	"context"
	"sync/atomic"
	"unsafe"
)

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
func Trace(msg string, fields ...Field) {
	GetGlobalLogger().dispatch(TraceLevel, msg, fields...)
}

// Debug logs a DEBUG-level message via the global logger.
func Debug(msg string, fields ...Field) {
	GetGlobalLogger().dispatch(DebugLevel, msg, fields...)
}

// Info logs an INFO-level message via the global logger.
func Info(msg string, fields ...Field) {
	GetGlobalLogger().dispatch(InfoLevel, msg, fields...)
}

// Warn logs a WARN-level message via the global logger.
func Warn(msg string, fields ...Field) {
	GetGlobalLogger().dispatch(WarnLevel, msg, fields...)
}

// Error logs an ERROR-level message via the global logger.
func Error(msg string, fields ...Field) {
	GetGlobalLogger().dispatch(ErrorLevel, msg, fields...)
}

// Tracef logs a TRACE-level formatted message via the global logger.
func Tracef(format string, args ...any) {
	GetGlobalLogger().Tracef(format, args...)
}

// Debugf logs a DEBUG-level formatted message via the global logger.
func Debugf(format string, args ...any) {
	GetGlobalLogger().Debugf(format, args...)
}

// Infof logs an INFO-level formatted message via the global logger.
func Infof(format string, args ...any) {
	GetGlobalLogger().Infof(format, args...)
}

// Warnf logs a WARN-level formatted message via the global logger.
func Warnf(format string, args ...any) {
	GetGlobalLogger().Warnf(format, args...)
}

// Errorf logs an ERROR-level formatted message via the global logger.
func Errorf(format string, args ...any) {
	GetGlobalLogger().Errorf(format, args...)
}

// GlobalWithContextFields returns a new context that carries the provided fields for
// use with the global logger.
func GlobalWithContextFields(ctx context.Context, fields ...Field) context.Context {
	return WithContextFields(ctx, fields...)
}
