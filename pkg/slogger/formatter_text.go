package slogger

import (
	"strconv"
	"strings"

	"github.com/sivaosorg/replify/pkg/encoding"
	"github.com/sivaosorg/replify/pkg/fj"
	"github.com/sivaosorg/replify/pkg/strutil"
)

// WithTimeFormat sets the time layout string used when formatting timestamps.
//
// Parameters:
//   - `fmt`: the Go time layout string (e.g. time.RFC3339Nano)
//
// Returns:
//
// the receiver, for method chaining.
func (f *TextFormatter) WithTimeFormat(fmt string) *TextFormatter {
	f.timeFormat = fmt
	return f
}

// WithDisableColors disables ANSI colour codes in the output.
// Useful when writing to files, pipes, or CI environments that do not
// interpret escape sequences.
//
// Returns:
//
// the receiver, for method chaining.
func (f *TextFormatter) WithDisableColors() *TextFormatter {
	f.disableColors = true
	return f
}

// WithDisableTimestamp omits the timestamp from formatted output.
// Useful when the surrounding infrastructure (systemd, Docker) adds its own
// timestamps.
//
// Returns:
//
// the receiver, for method chaining.
func (f *TextFormatter) WithDisableTimestamp() *TextFormatter {
	f.disableTimestamp = true
	return f
}

// WithEnableCaller appends the source file and line number (caller=file:line)
// to formatted output, aiding in debugging.
//
// Returns:
//
// the receiver, for method chaining.
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

	useColor := !f.disableColors && istty(f.output)

	// Add caller information if enabled.
	if e.caller != nil {
		f.WithEnableCaller()
	}

	if !f.disableTimestamp {
		b.WriteString(e.Time().Format(f.timeFormat))
		b.WriteByte(' ')
	}

	levelStr := levelPad(e.GetLevel())
	if useColor {
		b.WriteString(levelColor(e.GetLevel()))
		b.WriteString(colorBold)
		b.WriteString(levelStr)
		b.WriteString(colorReset)
	} else {
		b.WriteString(levelStr)
	}
	b.WriteByte(' ')

	if l := e.Logger(); l != nil && strutil.IsNotEmpty(l.name) {
		b.WriteByte('[')
		b.WriteString(l.name)
		b.WriteString("] ")
	}

	b.WriteString(e.Message())

	for _, fld := range e.Fields() {
		b.WriteByte(' ')
		b.WriteString(fld.key)
		b.WriteByte('=')
		v := fld.Value()
		if encoding.IsValidJSON(v) && fj.IsValidJSON(v) {
			b.WriteString(v)
		} else {
			if shouldQuoting(v) {
				b.WriteString(strconv.Quote(v))
			} else {
				b.WriteString(v)
			}
		}
	}

	if f.enableCaller {
		if c := e.Caller(); c != nil {
			b.WriteString(" caller=")
			b.WriteString(c.File())
			b.WriteByte(':')
			b.WriteString(strconv.Itoa(c.Line()))
		}
	}

	b.WriteByte('\n')
	return []byte(b.String()), nil
}
