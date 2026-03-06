package slogger

import (
	"io"
	"strconv"
	"strings"
	"time"
)

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
//
// Parameters:
//   - `fmt`: a Go time layout string
//
// Returns:
//
// the receiver for chaining.
func (f *TextFormatter) WithTimeFormat(fmt string) *TextFormatter {
	f.timeFormat = fmt
	return f
}

// WithDisableColors disables ANSI colour codes in the output.
//
// Returns:
//
// the receiver for chaining.
func (f *TextFormatter) WithDisableColors() *TextFormatter {
	f.disableColors = true
	return f
}

// WithDisableTimestamp omits the timestamp from formatted output.
//
// Returns:
//
// the receiver for chaining.
func (f *TextFormatter) WithDisableTimestamp() *TextFormatter {
	f.disableTimestamp = true
	return f
}

// WithEnableCaller appends caller information to formatted output.
//
// Returns:
//
// the receiver for chaining.
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
		b.WriteString(e.Time.Format(f.timeFormat))
		b.WriteByte(' ')
	}

	levelStr := levelPad(e.Level)
	if useColor {
		b.WriteString(levelColor(e.Level))
		b.WriteString(colorBold)
		b.WriteString(levelStr)
		b.WriteString(colorReset)
	} else {
		b.WriteString(levelStr)
	}
	b.WriteByte(' ')

	if e.Logger != nil && e.Logger.name != "" {
		b.WriteByte('[')
		b.WriteString(e.Logger.name)
		b.WriteString("] ")
	}

	b.WriteString(e.Message)

	for i := range e.Fields {
		b.WriteByte(' ')
		b.WriteString(e.Fields[i].Key)
		b.WriteByte('=')
		v := e.Fields[i].Value()
		if needsQuoting(v) {
			b.WriteString(strconv.Quote(v))
		} else {
			b.WriteString(v)
		}
	}

	if f.enableCaller && e.Caller != nil {
		b.WriteString(" caller=")
		b.WriteString(e.Caller.File)
		b.WriteByte(':')
		b.WriteString(strconv.Itoa(e.Caller.Line))
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
