package slogger

import (
"encoding/json"
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
writeJSONString(&b, e.time.Format(f.timeFormat))

b.WriteByte(',')
writeJSONKey(&b, f.levelKey)
b.WriteByte(':')
writeJSONString(&b, e.level.String())

if e.logger != nil && e.logger.name != "" {
b.WriteByte(',')
writeJSONKey(&b, f.nameKey)
b.WriteByte(':')
writeJSONString(&b, e.logger.name)
}

b.WriteByte(',')
writeJSONKey(&b, f.messageKey)
b.WriteByte(':')
writeJSONString(&b, e.message)

for i := range e.fields {
b.WriteByte(',')
writeJSONKey(&b, e.fields[i].Key)
b.WriteByte(':')
writeJSONValue(&b, &e.fields[i])
}

if f.enableCaller && e.caller != nil {
b.WriteByte(',')
writeJSONKey(&b, f.callerKey)
b.WriteByte(':')
writeJSONString(&b, e.caller.file+":"+itoa(e.caller.line))
}

b.WriteByte('}')
b.WriteByte('\n')
return []byte(b.String()), nil
}

// writeJSONKey writes a JSON-encoded object key to b.
func writeJSONKey(b *strings.Builder, key string) {
writeJSONString(b, key)
}

// writeJSONString writes a JSON-encoded string value to b using encoding/json.
func writeJSONString(b *strings.Builder, s string) {
enc, _ := json.Marshal(s)
b.Write(enc)
}

// writeJSONValue writes the JSON encoding of a field's value to b.
func writeJSONValue(b *strings.Builder, f *Field) {
switch f.Type {
case StringType:
writeJSONString(b, f.strVal)
case Int64Type:
b.WriteString(itoa64(f.intVal))
case Float64Type:
enc, _ := json.Marshal(f.floatVal)
b.Write(enc)
case BoolType:
if f.boolVal {
b.WriteString("true")
} else {
b.WriteString("false")
}
case ErrorType:
if f.errVal == nil {
b.WriteString("null")
} else {
writeJSONString(b, f.errVal.Error())
}
case TimeType:
writeJSONString(b, f.timeVal.Format(time.RFC3339))
case DurationType:
writeJSONString(b, f.durVal.String())
case AnyType:
enc, err := json.Marshal(f.anyVal)
if err != nil {
writeJSONString(b, f.Value())
} else {
b.Write(enc)
}
default:
writeJSONString(b, f.Value())
}
}

// itoa converts an int to its decimal string representation without importing strconv.
func itoa(n int) string {
return itoa64(int64(n))
}

// itoa64 converts an int64 to its decimal string representation.
func itoa64(n int64) string {
if n == 0 {
return "0"
}
neg := false
if n < 0 {
neg = true
n = -n
}
var buf [20]byte
pos := len(buf)
for n > 0 {
pos--
buf[pos] = byte('0' + n%10)
n /= 10
}
if neg {
pos--
buf[pos] = '-'
}
return string(buf[pos:])
}
