package slogger

import (
	"strings"
	"time"
)

// NewJSONFormatter returns a JSONFormatter with sensible production defaults.
//
// Returns:
//
// a *JSONFormatter with keys ts, level, msg, caller, and name.
func NewJSONFormatter() *JSONFormatter {
return &JSONFormatter{
timeFormat: time.RFC3339,
timeKey:    "ts",
levelKey:   "level",
messageKey: "msg",
callerKey:  "caller",
nameKey:    "name",
}
}

// WithTimeFormat overrides the timestamp layout.
func (f *JSONFormatter) WithTimeFormat(fmt string) *JSONFormatter {
f.timeFormat = fmt
return f
}

// WithEnableCaller includes caller information in the JSON output.
func (f *JSONFormatter) WithEnableCaller() *JSONFormatter {
f.enableCaller = true
return f
}

// WithTimeKey sets the JSON key for the timestamp field.
func (f *JSONFormatter) WithTimeKey(key string) *JSONFormatter {
f.timeKey = key
return f
}

// WithLevelKey sets the JSON key for the level field.
func (f *JSONFormatter) WithLevelKey(key string) *JSONFormatter {
f.levelKey = key
return f
}

// WithMessageKey sets the JSON key for the message field.
func (f *JSONFormatter) WithMessageKey(key string) *JSONFormatter {
f.messageKey = key
return f
}

// WithCallerKey sets the JSON key for the caller field.
func (f *JSONFormatter) WithCallerKey(key string) *JSONFormatter {
f.callerKey = key
return f
}

// WithNameKey sets the JSON key for the logger name field.
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

if l := e.Logger(); l != nil && l.name != "" {
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
writeJSONKey(&b, fld.Key)
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
