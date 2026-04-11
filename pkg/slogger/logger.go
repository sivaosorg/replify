package slogger

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/sivaosorg/replify/pkg/strutil"
)

// With returns a child Logger that inherits all settings and prepends the
// provided fields to every subsequent log entry.
//
// Parameters:
//   - `fields`: structured fields to bind to the child logger
//
// Returns:
//
// a new *Logger sharing configuration with the receiver.
func (l *Logger) With(fields ...Field) *Logger {
	l.mu.RLock()
	child := &Logger{
		wmu:        l.wmu, // share write mutex for thread-safe output
		formatter:  l.formatter,
		output:     l.output,
		hooks:      l.hooks,
		name:       l.name,
		caller:     l.caller,
		callerSkip: l.callerSkip,
		sampling:   l.sampling,
	}
	// Deep copy the parent's fields slice while holding the lock to prevent
	// race conditions when the parent's fields are modified concurrently.
	// Field is a value type, so copying the slice copies all field values.
	parentFields := make([]Field, len(l.fields))
	copy(parentFields, l.fields)
	l.mu.RUnlock()
	child.level.Store(l.level.Load())
	merged := make([]Field, 0, len(parentFields)+len(fields))
	merged = append(merged, parentFields...)
	merged = append(merged, fields...)
	child.fields = merged
	return child
}

// WithContext returns an Entry bound to ctx for context-aware logging.
//
// Parameters:
//   - `ctx`: the context to associate with the returned entry
//
// Returns:
//
// a *Entry with logger and ctx set.
func (l *Logger) WithContext(ctx context.Context) *Entry {
	return &Entry{logger: l, ctx: ctx}
}

// Named returns a child Logger whose name extends the parent name with a dot
// separator when the parent already has a name.
//
// Parameters:
//   - `name`: the child logger segment
//
// Returns:
//
// a new *Logger with the composed name.
func (l *Logger) Named(name string) *Logger {
	l.mu.RLock()
	child := &Logger{
		wmu:        l.wmu, // share write mutex for thread-safe output
		formatter:  l.formatter,
		output:     l.output,
		hooks:      l.hooks,
		fields:     append([]Field(nil), l.fields...),
		caller:     l.caller,
		callerSkip: l.callerSkip,
		sampling:   l.sampling,
	}
	l.mu.RUnlock()
	child.level.Store(l.level.Load())
	if strutil.IsNotEmpty(l.name) {
		child.name = l.name + "." + name
	} else {
		child.name = name
	}
	return child
}

// SetLevel atomically updates the minimum log level.
//
// Parameters:
//   - `level`: the new minimum level
func (l *Logger) SetLevel(level Level) {
	l.level.Store(int32(level))
}

// GetLevel returns the current minimum log level.
//
// Returns:
//
// the active Level.
func (l *Logger) GetLevel() Level {
	return Level(l.level.Load())
}

// SetOutput replaces the output writer.
//
// Parameters:
//   - `w`: the new destination writer
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	l.output = w
	l.mu.Unlock()
}

// SetFormatter replaces the entry formatter.
//
// Parameters:
//   - `f`: the new Formatter
func (l *Logger) SetFormatter(f Formatter) {
	l.mu.Lock()
	l.formatter = f
	l.mu.Unlock()
}

// AddHook registers a Hook on this logger.
//
// Parameters:
//   - `hook`: the Hook to add
func (l *Logger) AddHook(hook Hook) {
	l.hooks.Add(hook)
}

// IsLevelEnabled reports whether level is at or above the current minimum.
//
// Parameters:
//   - `level`: the level to test
//
// Returns:
//
// true when the level would produce output.
func (l *Logger) IsLevelEnabled(level Level) bool {
	return level.IsEnabled(Level(l.level.Load()))
}

// Trace logs a TRACE-level message.
func (l *Logger) Trace(msg string, fields ...Field) { l.dispatch(TraceLevel, msg, fields...) }

// Debug logs a DEBUG-level message.
func (l *Logger) Debug(msg string, fields ...Field) { l.dispatch(DebugLevel, msg, fields...) }

// Info logs an INFO-level message.
func (l *Logger) Info(msg string, fields ...Field) { l.dispatch(InfoLevel, msg, fields...) }

// Warn logs a WARN-level message.
func (l *Logger) Warn(msg string, fields ...Field) { l.dispatch(WarnLevel, msg, fields...) }

// Error logs an ERROR-level message.
func (l *Logger) Error(msg string, fields ...Field) { l.dispatch(ErrorLevel, msg, fields...) }

// Fatal logs a FATAL-level message and calls os.Exit(1).
func (l *Logger) Fatal(msg string, fields ...Field) {
	l.dispatch(FatalLevel, msg, fields...)
	os.Exit(1)
}

// Panic logs a PANIC-level message and panics with msg.
func (l *Logger) Panic(msg string, fields ...Field) {
	l.dispatch(PanicLevel, msg, fields...)
	panic(msg)
}

// Tracef logs a TRACE-level formatted message.
func (l *Logger) Tracef(format string, args ...any) {
	l.dispatch(TraceLevel, fmt.Sprintf(format, args...))
}

// Debugf logs a DEBUG-level formatted message.
func (l *Logger) Debugf(format string, args ...any) {
	l.dispatch(DebugLevel, fmt.Sprintf(format, args...))
}

// Infof logs an INFO-level formatted message.
func (l *Logger) Infof(format string, args ...any) {
	l.dispatch(InfoLevel, fmt.Sprintf(format, args...))
}

// Warnf logs a WARN-level formatted message.
func (l *Logger) Warnf(format string, args ...any) {
	l.dispatch(WarnLevel, fmt.Sprintf(format, args...))
}

// Errorf logs an ERROR-level formatted message.
func (l *Logger) Errorf(format string, args ...any) {
	l.dispatch(ErrorLevel, fmt.Sprintf(format, args...))
}

// Fatalf logs a FATAL-level formatted message and calls os.Exit(1).
func (l *Logger) Fatalf(format string, args ...any) {
	l.dispatch(FatalLevel, fmt.Sprintf(format, args...))
	os.Exit(1)
}

// Panicf logs a PANIC-level formatted message and panics.
func (l *Logger) Panicf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	l.dispatch(PanicLevel, msg)
	panic(msg)
}

// WithLevel sets the minimum log level and returns the logger for chaining.
// The update is atomic and safe to call concurrently.
// This is the fluent equivalent of SetLevel for use in builder-style
// logger configuration.
//
// Parameters:
//   - `level`: the new minimum severity level
//
// Returns:
//
// the receiver *Logger, enabling method chaining.
func (l *Logger) WithLevel(level Level) *Logger {
	l.level.Store(int32(level))
	return l
}

// WithFormatter replaces the entry formatter and returns the logger for chaining.
// The update is lock-protected and safe to call concurrently.
//
// Parameters:
//   - `f`: the new Formatter to use for serialising log entries
//
// Returns:
//
// the receiver *Logger, enabling method chaining.
func (l *Logger) WithFormatter(f Formatter) *Logger {
	l.mu.Lock()
	l.formatter = f
	l.mu.Unlock()
	return l
}

// WithOutput replaces the primary output writer and returns the logger for
// chaining. The update is lock-protected and safe to call concurrently.
//
// Parameters:
//   - `w`: the new destination writer
//
// Returns:
//
// the receiver *Logger, enabling method chaining.
func (l *Logger) WithOutput(w io.Writer) *Logger {
	l.mu.Lock()
	l.output = w
	l.mu.Unlock()
	return l
}

// WithCaller enables or disables automatic source-location capture and returns
// the logger for chaining.
//
// Parameters:
//   - `enabled`: true to capture file/line/function on every log call
//
// Returns:
//
// the receiver *Logger, enabling method chaining.
func (l *Logger) WithCaller(enabled bool) *Logger {
	l.mu.Lock()
	l.caller = enabled
	l.mu.Unlock()
	return l
}

// WithCallerSkip sets the number of additional stack frames to skip when
// capturing caller information and returns the logger for chaining.
// Useful when wrapping slogger inside a library adapter.
//
// Parameters:
//   - `skip`: the number of additional frames to skip beyond the default
//
// Returns:
//
// the receiver *Logger, enabling method chaining.
func (l *Logger) WithCallerSkip(skip int) *Logger {
	l.mu.Lock()
	l.callerSkip = skip
	l.mu.Unlock()
	return l
}

// WithSampling enables per-message rate limiting on the logger and returns it
// for chaining. Sampling prevents log storms by allowing only the first opts.First
// identical messages per opts.Period, then every opts.Thereafter-th message.
//
// Parameters:
//   - `opts`: the SamplingOptions that define the rate-limiting behaviour
//
// Returns:
//
// the receiver *Logger, enabling method chaining.
func (l *Logger) WithSampling(opts SamplingOptions) *Logger {
	l.mu.Lock()
	l.sampling = newSampler(opts)
	l.mu.Unlock()
	return l
}

// WithRotation enables per-level file rotation on the logger and returns it for
// chaining. It creates a LevelFileWriter from opts and registers a
// LevelWriterHook on the logger. If initialisation fails, a diagnostic message
// is written to stderr and the logger continues without rotation.
//
// Parameters:
//   - `opts`: the RotationOptions that configure directory, size limits, and compression
//
// Returns:
//
// the receiver *Logger, enabling method chaining.
func (l *Logger) WithRotation(opts RotationOptions) *Logger {
	lfw, err := newLevelFileWriter(opts)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "slogger: rotation setup failed: %v\n", err)
		return l
	}
	l.mu.RLock()
	formatter := l.formatter
	l.mu.RUnlock()
	l.hooks.Add(NewLevelWriterHook(lfw, formatter))
	return l
}

// getCaller captures the source location of the log call site.
//
// Parameters:
//   - `skip`: additional frames to skip beyond the internal logging frames
//
// Returns:
//
// a *CallerInfo describing the file, line, and function, or nil if the
// call site could not be determined.
func (l *Logger) getCaller(skip int) *CallerInfo {
	pcs := make([]uintptr, 1)
	n := runtime.Callers(4+skip, pcs)
	if n == 0 {
		return nil
	}
	frame, _ := runtime.CallersFrames(pcs).Next()
	return &CallerInfo{
		file:     trimFilepath(frame.File),
		line:     frame.Line,
		function: frame.Function,
	}
}

// dispatch fires a dispatch entry without a context.
//
// Parameters:
//   - `level`: the dispatch level
//   - `msg`: the message to dispatch
//   - `fields`: optional structured fields
func (l *Logger) dispatch(level Level, msg string, fields ...Field) {
	//nolint:gochecknoglobals
	//nolint:SA1012
	l.dispatchContext(nil, level, msg, fields...)
}

// dispatchContext is the internal dispatch point for all log events.
//
// Parameters:
//   - `ctx`: the context to associate with the entry
//   - `level`: the log level
//   - `msg`: the message to log
//   - `fields`: optional structured fields
func (l *Logger) dispatchContext(ctx context.Context, level Level, msg string, fields ...Field) {
	if !level.IsEnabled(Level(l.level.Load())) {
		return
	}

	// Check sampling.
	if l.sampling != nil && !l.sampling.allow(msg) {
		return
	}

	// Acquire entry.
	e := acquireEntry(l)
	e.time = time.Now()
	e.level = level
	e.message = msg
	e.ctx = ctx

	// Calculate total number of fields.
	total := len(l.fields) + len(fields)
	if ctx != nil {
		total += len(FieldsFromContext(ctx))
	}

	// Ensure enough capacity for all fields.
	if cap(e.fields) < total {
		e.fields = make([]Field, 0, total)
	}

	// Append fields in order: logger-bound, context, call-site.
	e.fields = append(e.fields, l.fields...)
	if ctx != nil {
		e.fields = append(e.fields, FieldsFromContext(ctx)...)
	}
	e.fields = append(e.fields, fields...)

	// Add caller information if enabled.
	if l.caller {
		e.caller = l.getCaller(l.callerSkip)
	}

	l.mu.RLock()
	formatter := l.formatter
	output := l.output
	l.mu.RUnlock()

	// Serialize entry.
	data, err := formatter.Format(e)
	if err == nil {
		// Use shared writeMu to synchronize writes across child loggers
		l.wmu.Lock()
		_, _ = output.Write(data)
		l.wmu.Unlock()
	}

	// Fire hooks.
	_ = l.hooks.Fire(level, e)

	// Release entry.
	releaseEntry(e)
}
