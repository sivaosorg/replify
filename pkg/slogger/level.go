package slogger

import (
	"errors"
	"strings"
)

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

// String returns the uppercase name of the level.
//
// Returns:
//
// one of "TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL", "PANIC",
// or "UNKNOWN" for unrecognised values.
func (l Level) String() string {
	switch l {
	case TraceLevel:
		return "TRACE"
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	case PanicLevel:
		return "PANIC"
	default:
		return "UNKNOWN"
	}
}

// ParseLevel converts a string to a Level in a case-insensitive manner.
//
// Parameters:
//   - `s`: the level name, e.g. "info", "WARN", "debug"
//
// Returns:
//
// the matching Level and a nil error, or zero and an error when the name is
// not recognised.
func ParseLevel(s string) (Level, error) {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "TRACE":
		return TraceLevel, nil
	case "DEBUG":
		return DebugLevel, nil
	case "INFO":
		return InfoLevel, nil
	case "WARN", "WARNING":
		return WarnLevel, nil
	case "ERROR":
		return ErrorLevel, nil
	case "FATAL":
		return FatalLevel, nil
	case "PANIC":
		return PanicLevel, nil
	default:
		return TraceLevel, errors.New("slogger: unknown log level: " + s)
	}
}

// IsEnabled reports whether the level should be logged given a minimum level.
//
// Parameters:
//   - `min`: the minimum level that must be reached for logging to occur
//
// Returns:
//
// true when l >= min.
func (l Level) IsEnabled(min Level) bool {
	return l >= min
}
