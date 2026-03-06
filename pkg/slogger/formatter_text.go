package slogger

import (
"io"
"strconv"
"strings"
"time"
)

// NewTextFormatter returns a TextFormatter that writes to output.
//
// Parameters:
//   - `output`: the destination writer (used only for TTY detection)
//
// Returns:
//
// a *TextFormatter with RFC3339 timestamps and colours enabled when output is
// a terminal.
func NewTextFormatter(output io.Writer) *TextFormatter {
return &TextFormatter{
timeFormat: time.RFC3339,
output:     output,
}
}

// WithTimeFormat sets the time layout string used when formatting timestamps.
func (f *TextFormatter) WithTimeFormat(fmt string) *TextFormatter {
f.timeFormat = fmt
return f
}

// WithDisableColors disables ANSI colour codes in the output.
func (f *TextFormatter) WithDisableColors() *TextFormatter {
f.disableColors = true
return f
}

// WithDisableTimestamp omits the timestamp from formatted output.
func (f *TextFormatter) WithDisableTimestamp() *TextFormatter {
f.disableTimestamp = true
return f
}

// WithEnableCaller appends caller information to formatted output.
func (f *TextFormatter) WithEnableCaller() *TextFormatter {
f.enableCaller = true
return f
}

// Format serialises e to a human-readable key=value byte slice.
//
// Parameters:
//   - `e`: the log entry to format
//
// Returns:
//
// the formatted bytes and any encoding error.
func (f *TextFormatter) Format(e *Entry) ([]byte, error) {
var b strings.Builder

useColor := !f.disableColors && IsTTY(f.output)

if !f.disableTimestamp {
b.WriteString(e.time.Format(f.timeFormat))
b.WriteByte(' ')
}

levelStr := levelPad(e.level)
if useColor {
b.WriteString(levelColor(e.level))
b.WriteString(colorBold)
b.WriteString(levelStr)
b.WriteString(colorReset)
} else {
b.WriteString(levelStr)
}
b.WriteByte(' ')

if e.logger != nil && e.logger.name != "" {
b.WriteByte('[')
b.WriteString(e.logger.name)
b.WriteString("] ")
}

b.WriteString(e.message)

for i := range e.fields {
b.WriteByte(' ')
b.WriteString(e.fields[i].Key)
b.WriteByte('=')
v := e.fields[i].Value()
if needsQuoting(v) {
b.WriteString(strconv.Quote(v))
} else {
b.WriteString(v)
}
}

if f.enableCaller && e.caller != nil {
b.WriteString(" caller=")
b.WriteString(e.caller.file)
b.WriteByte(':')
b.WriteString(strconv.Itoa(e.caller.line))
}

b.WriteByte('\n')
return []byte(b.String()), nil
}

// levelPad returns the level string padded to 5 characters.
func levelPad(l Level) string {
s := l.String()
for len(s) < 5 {
s += " "
}
return s
}

// needsQuoting reports whether a value string should be quoted in text output.
func needsQuoting(s string) bool {
for _, c := range s {
if c == ' ' || c == '=' || c == '"' || c == '\\' {
return true
}
}
return false
}
