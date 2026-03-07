package slogger

// ///////////////////////////////////////////////////////////////////////////
// Log Level constants
// ///////////////////////////////////////////////////////////////////////////

// Log level severity constants, ordered from lowest to highest severity.
// The zero value (TraceLevel) is the most verbose level.
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
// Field type constants
// ///////////////////////////////////////////////////////////////////////////

// Field value-type constants enumerate the concrete types that a Field may carry.
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

	// Int8Type holds an int8 value stored as int64.
	Int8Type

	// Int16Type holds an int16 value stored as int64.
	Int16Type

	// Int32Type holds an int32 value stored as int64.
	Int32Type

	// UintType holds a uint value stored as uint64.
	UintType

	// Uint8Type holds a uint8 value stored as uint64.
	Uint8Type

	// Uint16Type holds a uint16 value stored as uint64.
	Uint16Type

	// Uint32Type holds a uint32 value stored as uint64.
	Uint32Type

	// Uint64Type holds a uint64 value.
	Uint64Type

	// Float32Type holds a float32 value stored as float64.
	Float32Type

	// TimefType holds a time.Time value formatted with a custom layout string.
	TimefType

	// JSONType holds a JSON-encoded string representation of an arbitrary value.
	JSONType
)

// ///////////////////////////////////////////////////////////////////////////
// ANSI color codes
// ///////////////////////////////////////////////////////////////////////////

// Terminal ANSI escape sequences used by TextFormatter when color output is
// enabled. All codes are standard 8-color VT100 sequences.
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorGreen  = "\033[32m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

// ///////////////////////////////////////////////////////////////////////////
// Default directories and file names
// ///////////////////////////////////////////////////////////////////////////

// Directory and file-name defaults used by the rotation subsystem.
const (
	// defaultLogDir is the base log directory created when RotationOptions.Dir
	// is empty.
	defaultLogDir = "logs"

	// defaultArchivedDir is the sub-directory inside defaultLogDir where
	// rotated log files are stored, organized further by date.
	defaultArchivedDir = "archived"

	// logFileDebug is the filename for DEBUG-level log output.
	logFileDebug = "debug.log"

	// logFileInfo is the filename for INFO-level log output.
	logFileInfo = "info.log"

	// logFileWarn is the filename for WARN-level log output.
	logFileWarn = "warn.log"

	// logFileError is the filename for ERROR and above log output.
	logFileError = "error.log"
)

// ///////////////////////////////////////////////////////////////////////////
// Default formatting and timestamp patterns
// ///////////////////////////////////////////////////////////////////////////

// Timestamp layout constants used by formatters and the rotation archiver.
const (
	// defaultTimeFormat is the default timestamp layout for both
	// TextFormatter and JSONFormatter. It is equivalent to time.RFC3339.
	defaultTimeFormat = "2006-01-02 15:04:05.999999" // "2006-01-02T15:04:05Z07:00"

	// archiveDateFormat is the date-based sub-directory layout used when
	// creating the archived log folder (e.g. "2006-01-02").
	archiveDateFormat = "2006-01-02"

	// archiveStampFormat is the timestamp prefix embedded in the rotated
	// archive filename (e.g. "20060102150405").
	archiveStampFormat = "20060102150405"

	// levelPadWidth is the fixed column width used when padding level names
	// in text-formatted output.
	levelPadWidth = 5
)

// ///////////////////////////////////////////////////////////////////////////
// Default JSON field keys
// ///////////////////////////////////////////////////////////////////////////

// Default JSON output key names used by JSONFormatter.
const (
	// defaultJSONTimeKey is the default key for the timestamp field.
	defaultJSONTimeKey = "ts"

	// defaultJSONLevelKey is the default key for the log level field.
	defaultJSONLevelKey = "level"

	// defaultJSONMessageKey is the default key for the log message field.
	defaultJSONMessageKey = "msg"

	// defaultJSONCallerKey is the default key for the caller location field.
	defaultJSONCallerKey = "caller"

	// defaultJSONNameKey is the default key for the logger name field.
	defaultJSONNameKey = "name"
)

// ///////////////////////////////////////////////////////////////////////////
// Default rotation sizes
// ///////////////////////////////////////////////////////////////////////////

// Size and capacity defaults for the rotation subsystem.
const (
	// defaultMaxBytes is the maximum log-file size before automatic rotation
	// is triggered. Default is 10 MiB.
	defaultMaxBytes int64 = 10 * 1024 * 1024

	// defaultEntryFieldCap is the initial field-slice capacity allocated for
	// each recycled Entry to avoid reallocations on the common case.
	defaultEntryFieldCap = 8
)
