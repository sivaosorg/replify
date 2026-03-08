package slogger

import (
	"context"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/sivaosorg/replify/pkg/encoding"
)

// ///////////////////////////////////////////////////////////////////////////
// Level and FieldType — typed enumerations
// ///////////////////////////////////////////////////////////////////////////

// Level represents the severity of a log message.
// Severity increases with numeric value; TraceLevel is the least severe.
type Level int32

// FieldType enumerates the concrete value types that a Field may carry.
// Callers should use the typed constructor functions (String, Int64, Bool, …)
// rather than setting FieldType directly.
type FieldType int

// ///////////////////////////////////////////////////////////////////////////
// Logger — core type
// ///////////////////////////////////////////////////////////////////////////

// Logger is the core logging type.
// All exported methods are safe for concurrent use.
//
// A Logger is created with New and customised through functional options.
// Child loggers that inherit configuration can be derived with With and Named.
// The zero value is not usable; always create a Logger through New.
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

// ///////////////////////////////////////////////////////////////////////////
// Entry — in-flight log event
// ///////////////////////////////////////////////////////////////////////////

// Entry represents a single in-flight log event.
// Entries are pooled internally; callers must not retain an Entry after the
// log call that produced it returns.
//
// To obtain an Entry that carries a context, use Logger.WithContext.
// All log-level methods on Entry are safe for concurrent use.
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
//
// All fields are read-only after construction. Access them through the
// exported accessor methods File, Line, and Function.
type CallerInfo struct {
	file     string
	line     int
	function string
}

// ///////////////////////////////////////////////////////////////////////////
// Field — structured key-value data
// ///////////////////////////////////////////////////////////////////////////

// Field is a typed key-value pair attached to a log entry.
// All value variants are stored inline to minimise heap allocations.
//
// Use the typed constructors — String, Int, Int64, Float64, Bool, Err,
// Time, Duration, Any, Int8, Int16, Int32, Uint, Uint64, Float32, Timef, JSON —
// to build Field values; do not construct them directly.
// Access the key and type via the Key() and FieldType() accessor methods.
type Field struct {
	key       string
	typ       FieldType
	strVal    string
	intVal    int64
	uint64Val uint64
	floatVal  float64
	boolVal   bool
	errVal    error
	timeVal   time.Time
	durVal    time.Duration
	anyVal    any
}

// ///////////////////////////////////////////////////////////////////////////
// Formatter — serialisation interface
// ///////////////////////////////////////////////////////////////////////////

// Formatter serialises a log Entry to a byte slice.
// Implementations must be safe for concurrent use from multiple goroutines.
//
// Two built-in implementations are provided:
//   - TextFormatter: human-readable key=value lines
//   - JSONFormatter: single-line JSON objects
type Formatter interface {
	Format(*Entry) ([]byte, error)
}

// ///////////////////////////////////////////////////////////////////////////
// Hook — side-effect interface
// ///////////////////////////////////////////////////////////////////////////

// Hook fires side-effects (e.g. alerting, metrics, remote shipping) for
// matching log levels. Implementations must be safe for concurrent use.
//
// Register hooks on a Logger with Logger.AddHook.
// Hooks are called after the entry has been written to the primary output.
type Hook interface {
	// Levels returns the set of log levels this hook handles.
	Levels() []Level

	// Fire is called for each matching entry; it must not retain the entry.
	Fire(*Entry) error
}

// ///////////////////////////////////////////////////////////////////////////
// Hooks — level-indexed hook registry
// ///////////////////////////////////////////////////////////////////////////

// Hooks is a level-indexed registry of Hook instances.
// All methods are safe for concurrent use.
type Hooks struct {
	mu    sync.RWMutex
	hooks map[Level][]Hook
}

// ///////////////////////////////////////////////////////////////////////////
// Options — logger construction configuration
// ///////////////////////////////////////////////////////////////////////////

// Options configures a Logger at construction time.
// Pass functional options to New to override the defaults.
//
// All fields are optional; unset fields fall back to production-safe defaults
// (InfoLevel, TextFormatter, os.Stderr).
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
	First int `json:"first"`

	// Period is the window after which the counter resets.
	Period time.Duration `json:"period"`

	// Thereafter logs every Nth message after the First are exhausted.
	// Zero means drop all subsequent messages.
	Thereafter int `json:"thereafter"`
}

// ///////////////////////////////////////////////////////////////////////////
// Formatter implementations
// ///////////////////////////////////////////////////////////////////////////

// TextFormatter formats log entries as human-readable key=value lines.
//
// Output format:
//
// 2006-01-02 15:04:05.999999 INFO  message  key=value key2=value2\n
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
// {"ts":"2006-01-02 15:04:05.999999","level":"INFO","msg":"message","key":"value"}\n
type JSONFormatter struct {
	timeFormat   string
	enableCaller bool
	enableColor  bool
	timeKey      string
	levelKey     string
	messageKey   string
	callerKey    string
	nameKey      string
	color        *encoding.Style // Color style for log levels
}

// ///////////////////////////////////////////////////////////////////////////
// MultiWriter — fan-out writer
// ///////////////////////////////////////////////////////////////////////////

// MultiWriter fans log output to multiple io.Writer targets simultaneously.
// All Write calls are delivered to every registered writer in registration order.
type MultiWriter struct {
	writers []io.Writer
}

// ///////////////////////////////////////////////////////////////////////////
// Rotation types
// ///////////////////////////////////////////////////////////////////////////

// RotationOptions configures the rotating file writer.
// Pass RotationOptions to WithRotation to enable per-level log file rotation.
type RotationOptions struct {
	// Dir is the base log directory. Defaults to defaultLogDir ("logs").
	Dir string `json:"dir"`

	// MaxBytes is the maximum file size before rotation. Defaults to 10 MiB.
	MaxBytes int64 `json:"max_bytes"`

	// MaxAge is the maximum age of a log file before rotation. Zero means no age-based rotation.
	MaxAge time.Duration `json:"max_age"`

	// Compress controls whether rotated files are zipped. Defaults to false.
	Compress bool `json:"compress"`
}

// LevelFileWriter writes log entries to separate files per log level.
// It supports automatic rotation and ZIP compression of archived logs.
//
// Create instances with NewLevelFileWriter and register them on a logger using
// NewLevelWriterHook together with Logger.AddHook.
type LevelFileWriter struct {
	mu      sync.Mutex
	opts    RotationOptions
	writers map[Level]*rotatingFile
}

// LevelWriterHook is a Hook that writes log entries to level-specific files.
// Use it with AddHook to enable automatic per-level file logging.
type LevelWriterHook struct {
	writer    *LevelFileWriter
	formatter Formatter
	levels    []Level
}

// ///////////////////////////////////////////////////////////////////////////
// Internal — sampler and bucket
// ///////////////////////////////////////////////////////////////////////////

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

// ///////////////////////////////////////////////////////////////////////////
// Internal — rotating file handle
// ///////////////////////////////////////////////////////////////////////////

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

// ///////////////////////////////////////////////////////////////////////////
// Internal — context key
// ///////////////////////////////////////////////////////////////////////////

// contextKey is the unexported key type used to store log fields in a context.
// Using a named type prevents collisions with keys from other packages.
type contextKey struct{}

// ///////////////////////////////////////////////////////////////////////////
// Package-level variables
// ///////////////////////////////////////////////////////////////////////////

// entryPool is the package-level sync.Pool for recycling Entry objects.
// Capacity is pre-allocated to defaultEntryFieldCap to avoid early reallocations.
var entryPool = sync.Pool{New: func() any { return &Entry{fields: make([]Field, 0, defaultEntryFieldCap)} }}

// global holds the package-level logger pointer.
// Access is guarded by atomic operations so callers need no external lock.
var global unsafe.Pointer
