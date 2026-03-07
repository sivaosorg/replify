package slogger

import (
	"context"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// ///////////////////////////////////////////////////////////////////////////
// Level types
// ///////////////////////////////////////////////////////////////////////////

// Level represents the severity of a log message.
type Level int32

const (
	// TraceLevel is the most verbose level; for fine-grained diagnostic output.
	TraceLevel Level = iota
	// DebugLevel is for debugging information during development.
	DebugLevel
	// InfoLevel is for general operational messages.
	InfoLevel
	// WarnLevel is for potentially harmful situations.
	WarnLevel
	// ErrorLevel is for error events that might still allow the application to continue.
	ErrorLevel
	// FatalLevel logs the message and then calls os.Exit(1).
	FatalLevel
	// PanicLevel logs the message and then panics.
	PanicLevel
)

// ///////////////////////////////////////////////////////////////////////////
// Field types
// ///////////////////////////////////////////////////////////////////////////

// FieldType enumerates the concrete value types that a Field may carry.
type FieldType int

const (
	// StringType holds a string value.
	StringType FieldType = iota
	// Int64Type holds an int64 value.
	Int64Type
	// Float64Type holds a float64 value.
	Float64Type
	// BoolType holds a bool value.
	BoolType
	// ErrorType holds an error value.
	ErrorType
	// TimeType holds a time.Time value.
	TimeType
	// DurationType holds a time.Duration value.
	DurationType
	// AnyType holds an arbitrary interface{} value.
	AnyType
)

// Field is a typed key-value pair attached to a log entry.
// All value variants are stored inline to minimise heap allocations.
type Field struct {
	Key      string
	Type     FieldType
	strVal   string
	intVal   int64
	floatVal float64
	boolVal  bool
	errVal   error
	timeVal  time.Time
	durVal   time.Duration
	anyVal   interface{}
}

// ///////////////////////////////////////////////////////////////////////////
// Formatter and Hook interfaces
// ///////////////////////////////////////////////////////////////////////////

// Formatter serialises a log Entry to a byte slice.
// Implementations must be safe for concurrent use from multiple goroutines.
type Formatter interface {
	Format(*Entry) ([]byte, error)
}

// Hook fires side-effects (e.g. alerting, metrics) for matching log levels.
// Implementations must be safe for concurrent use.
type Hook interface {
	// Levels returns the set of log levels this hook handles.
	Levels() []Level
	// Fire is called for each matching entry; it must not retain the entry.
	Fire(*Entry) error
}

// ///////////////////////////////////////////////////////////////////////////
// Rotation types
// ///////////////////////////////////////////////////////////////////////////

// RotationOptions configures the rotating file writer.
type RotationOptions struct {
	// Dir is the base log directory. Defaults to "logs".
	Dir string
	// MaxBytes is the maximum file size before rotation. Defaults to 10 MB.
	MaxBytes int64
	// MaxAge is the maximum age of a log file before rotation. Zero means no age-based rotation.
	MaxAge time.Duration
	// Compress controls whether rotated files are zipped. Defaults to true.
	Compress bool
}

// LevelFileWriter writes log entries to separate files per log level.
// It supports automatic rotation and ZIP compression of archived logs.
type LevelFileWriter struct {
	mu      sync.Mutex
	opts    RotationOptions
	writers map[Level]*rotatingFile
}

// rotatingFile manages a single level-specific log file with size/age rotation.
type rotatingFile struct {
	mu       sync.Mutex
	path     string
	file     *os.File
	size     int64
	maxBytes int64
	openedAt time.Time
	maxAge   time.Duration
	compress bool
	dir      string
	level    Level
}

// LevelWriterHook is a Hook that writes log entries to level-specific files.
// Use it with AddHook to enable automatic per-level file logging.
type LevelWriterHook struct {
	writer    *LevelFileWriter
	formatter Formatter
	levels    []Level
}

// ///////////////////////////////////////////////////////////////////////////
// Logger and Entry types
// ///////////////////////////////////////////////////////////////////////////

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
