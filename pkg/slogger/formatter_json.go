package slogger

import (
	"strings"

	"github.com/sivaosorg/replify/pkg/strutil"
)

// WithTimeFormat overrides the timestamp layout used in JSON output.
//
// Parameters:
//   - `fmt`: the Go time layout string (e.g. time.RFC3339Nano)
//
// Returns:
//
// the receiver, for method chaining.
func (f *JSONFormatter) WithTimeFormat(fmt string) *JSONFormatter {
	f.timeFormat = fmt
	return f
}

// WithEnableCaller includes the caller source location (file:line) in the
// JSON output under the callerKey field.
//
// Returns:
//
// the receiver, for method chaining.
func (f *JSONFormatter) WithEnableCaller() *JSONFormatter {
	f.enableCaller = true
	return f
}

// WithTimeKey sets the JSON key for the timestamp field.
//
// Parameters:
//   - `key`: the desired JSON key name
//
// Returns:
//
// the receiver, for method chaining.
func (f *JSONFormatter) WithTimeKey(key string) *JSONFormatter {
	f.timeKey = key
	return f
}

// WithLevelKey sets the JSON key for the level field.
//
// Parameters:
//   - `key`: the desired JSON key name
//
// Returns:
//
// the receiver, for method chaining.
func (f *JSONFormatter) WithLevelKey(key string) *JSONFormatter {
	f.levelKey = key
	return f
}

// WithMessageKey sets the JSON key for the message field.
//
// Parameters:
//   - `key`: the desired JSON key name
//
// Returns:
//
// the receiver, for method chaining.
func (f *JSONFormatter) WithMessageKey(key string) *JSONFormatter {
	f.messageKey = key
	return f
}

// WithCallerKey sets the JSON key for the caller field.
//
// Parameters:
//   - `key`: the desired JSON key name
//
// Returns:
//
// the receiver, for method chaining.
func (f *JSONFormatter) WithCallerKey(key string) *JSONFormatter {
	f.callerKey = key
	return f
}

// WithNameKey sets the JSON key for the logger name field.
//
// Parameters:
//   - `key`: the desired JSON key name
//
// Returns:
//
// the receiver, for method chaining.
func (f *JSONFormatter) WithNameKey(key string) *JSONFormatter {
	f.nameKey = key
	return f
}

// Format serialises e to a single-line JSON byte slice.
//
// Parameters:
//   - `e`: the log entry to format
//
// Returns:
//
// the formatted bytes and any encoding error.
func (f *JSONFormatter) Format(e *Entry) ([]byte, error) {
	var b strings.Builder
	b.WriteByte('{')

	writeJSONKey(&b, f.timeKey)
	b.WriteByte(':')
	writeJSONString(&b, e.Time().Format(f.timeFormat))

	b.WriteByte(',')
	writeJSONKey(&b, f.levelKey)
	b.WriteByte(':')
	writeJSONString(&b, e.GetLevel().String())

	if l := e.Logger(); l != nil && strutil.IsNotEmpty(l.name) {
		b.WriteByte(',')
		writeJSONKey(&b, f.nameKey)
		b.WriteByte(':')
		writeJSONString(&b, l.name)
	}

	b.WriteByte(',')
	writeJSONKey(&b, f.messageKey)
	b.WriteByte(':')
	writeJSONString(&b, e.Message())

	for _, fld := range e.Fields() {
		b.WriteByte(',')
		writeJSONKey(&b, fld.key)
		b.WriteByte(':')
		writeJSONValue(&b, &fld)
	}

	if f.enableCaller {
		if c := e.Caller(); c != nil {
			b.WriteByte(',')
			writeJSONKey(&b, f.callerKey)
			b.WriteByte(':')
			writeJSONString(&b, c.File()+":"+itoa(c.Line()))
		}
	}

	b.WriteByte('}')
	b.WriteByte('\n')

	return []byte(b.String()), nil
}
