package slogger

import (
	"context"
	"io"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// contextKey is the unexported key type used to store log fields in a context.
type contextKey struct{}

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

// Entry represents a single in-flight log event.
// Entries are pooled internally; callers must not retain an Entry after the
// log call that produced it returns.
type Entry struct {
	logger  *Logger
	time    time.Time
	level   Level
	message string
	fields  []Field
	caller  *CallerInfo
	ctx     context.Context
}

// CallerInfo holds the runtime source-location data captured when caller
// reporting is enabled on a Logger.
type CallerInfo struct {
	file     string
	line     int
	function string
}

// Hooks is a level-indexed registry of Hook instances.
type Hooks struct {
	mu    sync.RWMutex
	hooks map[Level][]Hook
}

// samplerBucket tracks message counts for a single message key.
type samplerBucket struct {
	mu      sync.Mutex
	count   uint64
	resetAt time.Time
}

// sampler applies per-message rate limiting.
type sampler struct {
	opts    SamplingOptions
	buckets sync.Map // string (message) -> *samplerBucket
}

// MultiWriter fans log output to multiple io.Writer targets simultaneously.
type MultiWriter struct {
	writers []io.Writer
}

// Options configures a Logger at construction time.
// Pass functional options to New to override the defaults.
type Options struct {
	// Level is the minimum level that will be logged.
	Level Level
	// Formatter serialises entries; defaults to TextFormatter on os.Stderr.
	Formatter Formatter
	// Output is the destination writer; defaults to os.Stderr.
	Output io.Writer
	// CallerReporter enables automatic source-location capture.
	CallerReporter bool
	// CallerSkip adds extra skip frames for library wrappers.
	CallerSkip int
	// Fields are attached to every entry emitted by the logger.
	Fields []Field
	// Name is a dot-separated logger identifier prepended to output.
	Name string
	// SamplingOpts, when non-nil, enables per-message rate limiting.
	SamplingOpts *SamplingOptions
	// RotationOpts, when non-nil, enables per-level file rotation.
	RotationOpts *RotationOptions
}

// SamplingOptions configures per-message rate limiting for a Logger.
// The first First messages per Period are always logged.
// After that, every Thereafter-th message is logged (0 means drop all remaining).
type SamplingOptions struct {
	// First is the number of identical messages always logged within Period.
	First int
	// Period is the window after which the counter resets.
	Period time.Duration
	// Thereafter logs every Nth message after the First are exhausted.
	// Zero means drop all subsequent messages.
	Thereafter int
}

// TextFormatter formats log entries as human-readable key=value lines.
//
// Output format:
//
//	2006-01-02T15:04:05Z07:00 INFO  message  key=value key2=value2\n
type TextFormatter struct {
	timeFormat       string
	disableColors    bool
	disableTimestamp bool
	enableCaller     bool
	output           io.Writer
}

// JSONFormatter formats log entries as single-line JSON objects.
//
// Output example:
//
//	{"ts":"2006-01-02T15:04:05Z07:00","level":"INFO","msg":"message","key":"value"}\n
type JSONFormatter struct {
	timeFormat   string
	enableCaller bool
	timeKey      string
	levelKey     string
	messageKey   string
	callerKey    string
	nameKey      string
}

// entryPool is the package-level sync.Pool for recycling Entry objects.
var entryPool = sync.Pool{New: func() any { return &Entry{fields: make([]Field, 0, 8)} }}

// global holds the package-level logger pointer.
// Access is guarded by atomic operations so callers need no external lock.
var global unsafe.Pointer
