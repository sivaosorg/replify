package slogger

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// Logger is the core logging type.
// All exported methods are safe for concurrent use.
type Logger struct {
	mu         sync.RWMutex
	level      atomic.Int32
	formatter  Formatter
	output     io.Writer
	hooks      *Hooks
	fields     []Field
	name       string
	caller     bool
	callerSkip int
	sampling   *sampler
}

// New creates a Logger using the provided functional options.
// Unset options fall back to production-safe defaults (InfoLevel, TextFormatter,
// os.Stderr output).
//
// Parameters:
//   - `opts`: zero or more functions that mutate an *Options before the logger
//     is built
//
// Returns:
//
// a ready-to-use *Logger.
func New(opts ...func(*Options)) *Logger {
	o := defaultOptions()
	for _, fn := range opts {
		fn(o)
	}
	l := &Logger{
		formatter:  o.Formatter,
		output:     o.Output,
		hooks:      NewHooks(),
		fields:     append([]Field(nil), o.Fields...),
		name:       o.Name,
		caller:     o.CallerReporter,
		callerSkip: o.CallerSkip,
	}
	l.level.Store(int32(o.Level))
	if o.SamplingOpts != nil {
		l.sampling = newSampler(*o.SamplingOpts)
	}
	return l
}

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
		formatter:  l.formatter,
		output:     l.output,
		hooks:      l.hooks,
		name:       l.name,
		caller:     l.caller,
		callerSkip: l.callerSkip,
		sampling:   l.sampling,
	}
	l.mu.RUnlock()
	child.level.Store(l.level.Load())
	merged := make([]Field, 0, len(l.fields)+len(fields))
	merged = append(merged, l.fields...)
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
// a *Entry with Logger and Context set.
func (l *Logger) WithContext(ctx context.Context) *Entry {
	return &Entry{Logger: l, Context: ctx}
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
	if l.name != "" {
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
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields
func (l *Logger) Trace(msg string, fields ...Field) {
	l.log(TraceLevel, msg, fields...)
}

// Debug logs a DEBUG-level message.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields
func (l *Logger) Debug(msg string, fields ...Field) {
	l.log(DebugLevel, msg, fields...)
}

// Info logs an INFO-level message.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields
func (l *Logger) Info(msg string, fields ...Field) {
	l.log(InfoLevel, msg, fields...)
}

// Warn logs a WARN-level message.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields
func (l *Logger) Warn(msg string, fields ...Field) {
	l.log(WarnLevel, msg, fields...)
}

// Error logs an ERROR-level message.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields
func (l *Logger) Error(msg string, fields ...Field) {
	l.log(ErrorLevel, msg, fields...)
}

// Fatal logs a FATAL-level message and calls os.Exit(1).
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields
func (l *Logger) Fatal(msg string, fields ...Field) {
	l.log(FatalLevel, msg, fields...)
	os.Exit(1)
}

// Panic logs a PANIC-level message and panics with msg.
//
// Parameters:
//   - `msg`: the log message
//   - `fields`: optional structured fields
func (l *Logger) Panic(msg string, fields ...Field) {
	l.log(PanicLevel, msg, fields...)
	panic(msg)
}

// Tracef logs a TRACE-level formatted message.
//
// Parameters:
//   - `format`: a fmt.Sprintf format string
//   - `args`: arguments for the format string
func (l *Logger) Tracef(format string, args ...interface{}) {
	l.log(TraceLevel, fmt.Sprintf(format, args...))
}

// Debugf logs a DEBUG-level formatted message.
//
// Parameters:
//   - `format`: a fmt.Sprintf format string
//   - `args`: arguments for the format string
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log(DebugLevel, fmt.Sprintf(format, args...))
}

// Infof logs an INFO-level formatted message.
//
// Parameters:
//   - `format`: a fmt.Sprintf format string
//   - `args`: arguments for the format string
func (l *Logger) Infof(format string, args ...interface{}) {
	l.log(InfoLevel, fmt.Sprintf(format, args...))
}

// Warnf logs a WARN-level formatted message.
//
// Parameters:
//   - `format`: a fmt.Sprintf format string
//   - `args`: arguments for the format string
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log(WarnLevel, fmt.Sprintf(format, args...))
}

// Errorf logs an ERROR-level formatted message.
//
// Parameters:
//   - `format`: a fmt.Sprintf format string
//   - `args`: arguments for the format string
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log(ErrorLevel, fmt.Sprintf(format, args...))
}

// Fatalf logs a FATAL-level formatted message and calls os.Exit(1).
//
// Parameters:
//   - `format`: a fmt.Sprintf format string
//   - `args`: arguments for the format string
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.log(FatalLevel, fmt.Sprintf(format, args...))
	os.Exit(1)
}

// Panicf logs a PANIC-level formatted message and panics with the formatted string.
//
// Parameters:
//   - `format`: a fmt.Sprintf format string
//   - `args`: arguments for the format string
func (l *Logger) Panicf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.log(PanicLevel, msg)
	panic(msg)
}

// log fires a log entry without a context.
func (l *Logger) log(level Level, msg string, fields ...Field) {
	l.logCtx(nil, level, msg, fields...)
}

// logCtx is the internal dispatch point for all log events.
func (l *Logger) logCtx(ctx context.Context, level Level, msg string, fields ...Field) {
	if !level.IsEnabled(Level(l.level.Load())) {
		return
	}
	if l.sampling != nil && !l.sampling.allow(msg) {
		return
	}

	e := acquireEntry(l)
	e.Time = time.Now()
	e.Level = level
	e.Message = msg
	e.Context = ctx

	// Merge fields: logger-bound first, then context fields, then call-site fields.
	total := len(l.fields) + len(fields)
	if ctx != nil {
		total += len(FieldsFromContext(ctx))
	}
	if cap(e.Fields) < total {
		e.Fields = make([]Field, 0, total)
	}
	e.Fields = append(e.Fields, l.fields...)
	if ctx != nil {
		e.Fields = append(e.Fields, FieldsFromContext(ctx)...)
	}
	e.Fields = append(e.Fields, fields...)

	if l.caller {
		e.Caller = l.getCaller(l.callerSkip)
	}

	l.mu.RLock()
	formatter := l.formatter
	output := l.output
	l.mu.RUnlock()

	data, err := formatter.Format(e)
	if err == nil {
		l.mu.Lock()
		_, _ = output.Write(data)
		l.mu.Unlock()
	}

	_ = l.hooks.Fire(level, e)

	releaseEntry(e)
}

// getCaller captures the source location of the log call site.
//
// Parameters:
//   - `skip`: extra frames to skip on top of the standard 4
//
// Returns:
//
// a *CallerInfo or nil when the location cannot be determined.
func (l *Logger) getCaller(skip int) *CallerInfo {
	pcs := make([]uintptr, 1)
	// skip: runtime.Callers, getCaller, logCtx, log/logCtx caller (Info/Debug/…), + extra
	n := runtime.Callers(4+skip, pcs)
	if n == 0 {
		return nil
	}
	frame, _ := runtime.CallersFrames(pcs).Next()
	return &CallerInfo{
		File:     trimFile(frame.File),
		Line:     frame.Line,
		Function: frame.Function,
	}
}

// trimFile shortens an absolute file path to just the last two path segments.
func trimFile(path string) string {
	const sep = '/'
	slash := -1
	count := 0
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == sep {
			count++
			if count == 2 {
				slash = i
				break
			}
		}
	}
	if slash >= 0 {
		return path[slash+1:]
	}
	return path
}
