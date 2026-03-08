package slogger

import (
	"archive/zip"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sivaosorg/replify/pkg/encoding"
	"github.com/sivaosorg/replify/pkg/fj"
	"github.com/sivaosorg/replify/pkg/strutil"
	"github.com/sivaosorg/replify/pkg/sysx"
)

// ///////////////////////////////////////////////////////////////////////////
// Entry pool helpers
// ///////////////////////////////////////////////////////////////////////////

// acquireEntry retrieves an Entry from the pool and sets its logger field.
//
// Parameters:
//   - `l`: the logger that owns this entry
//
// Returns:
//
// a ready-to-use *Entry with logger set.
func acquireEntry(l *Logger) *Entry {
	e := entryPool.Get().(*Entry)
	e.logger = l
	return e
}

// releaseEntry resets e and returns it to the pool for reuse.
//
// Parameters:
//   - `e`: the entry to release
func releaseEntry(e *Entry) {
	e.reset()
	entryPool.Put(e)
}

// ///////////////////////////////////////////////////////////////////////////
// Logger internals
// ///////////////////////////////////////////////////////////////////////////

// trimFilepath shortens an absolute file path to just the last two path segments.
// For example "/home/user/project/pkg/foo/bar.go" becomes "pkg/foo/bar.go".
func trimFilepath(path string) string {
	const sep = '/'
	slash := -1
	count := 0
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == sep {
			count++
			if count == 2 {
				slash = i
				break
			}
		}
	}
	if slash >= 0 {
		return path[slash+1:]
	}
	return path
}

// levelFileName returns the log filename for the given level using the
// package-level file name constants.
//
// Parameters:
//   - `level`: the log level
//
// Returns:
//
// the log filename string for the given level.
func levelFileName(level Level) string {
	switch level {
	case DebugLevel:
		return logFileDebug
	case InfoLevel:
		return logFileInfo
	case WarnLevel:
		return logFileWarn
	default:
		return logFileError
	}
}

// compressToZip writes the contents of srcPath into a new ZIP archive at zipPath.
//
// Parameters:
//   - `srcPath`: the path to the file to compress
//   - `zipPath`: the path where the ZIP archive should be created
//
// Returns:
//
// an error if compression fails, nil otherwise.
//
// Example:
//
//	compressToZip("debug.log", "debug.log.zip")
func compressToZip(srcPath, zipPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	zf, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zf.Close()

	zw := zip.NewWriter(zf)
	defer zw.Close()

	w, err := zw.Create(filepath.Base(srcPath))
	if err != nil {
		return err
	}

	_, err = io.Copy(w, src)
	return err
}

// ///////////////////////////////////////////////////////////////////////////
// Color helpers
// ///////////////////////////////////////////////////////////////////////////

// levelColor returns the ANSI escape sequence associated with level l.
//
// Parameters:
//   - `l`: the log level
//
// Returns:
//
// an ANSI colour code string; Trace=cyan, Debug=blue, Info=green,
// Warn=yellow, Error/Fatal/Panic=red+bold.
//
// Example:
//
//	color := levelColor(InfoLevel) // "\x1b[32m"
func levelColor(l Level) string {
	switch l {
	case TraceLevel:
		return colorCyan
	case DebugLevel:
		return colorBlue
	case InfoLevel:
		return colorGreen
	case WarnLevel:
		return colorYellow
	case ErrorLevel, FatalLevel, PanicLevel:
		return colorRed
	default:
		return colorReset
	}
}

// ///////////////////////////////////////////////////////////////////////////
// Text formatter helpers
// ///////////////////////////////////////////////////////////////////////////

// levelPad returns the level string right-padded to levelPadWidth characters
// so that all levels align in column-based text output.
//
// Parameters:
//   - `l`: the log level
//
// Returns:
//
// the level string right-padded to levelPadWidth characters
//
// Example:
//
//	levelPad(InfoLevel) // "INFO    "
func levelPad(l Level) string {
	s := l.String()
	for len(s) < levelPadWidth {
		s += " "
	}
	return s
}

// shouldQuoting reports whether a value string should be quoted in text output.
// Strings that contain spaces, equals signs, double-quotes, or backslashes
// are wrapped in strconv.Quote to preserve round-trip fidelity.
//
// Parameters:
//   - `s`: the string to check
//
// Returns:
//
// true if the string should be quoted, false otherwise
//
// Example:
//
//	shouldQuoting("hello world") // true
//	shouldQuoting("hello") // false
func shouldQuoting(s string) bool {
	for _, c := range s {
		if c == ' ' || c == '=' || c == '"' || c == '\\' {
			return true
		}
	}
	return false
}

// ///////////////////////////////////////////////////////////////////////////
// JSON formatter helpers
// ///////////////////////////////////////////////////////////////////////////

// writeJSONKey writes a JSON-encoded object key to b.
//
// Parameters:
//   - `b`: the strings.Builder to write to
//   - `key`: the key to write
//
// Example:
//
//	var b strings.Builder
//	writeJSONKey(&b, "hello") // writes "hello"
func writeJSONKey(b *strings.Builder, key string) {
	writeJSONString(b, key)
}

// writeJSONString writes a JSON-encoded string value to b using encoding/json.
//
// Parameters:
//   - `b`: the strings.Builder to write to
//   - `s`: the string to write
//
// Example:
//
//	var b strings.Builder
//	writeJSONString(&b, "hello") // writes "\"hello\""
func writeJSONString(b *strings.Builder, s string) {
	if encoding.IsValidJSON(s) && fj.IsValidJSON(s) {
		b.Write([]byte(s))
		return
	}
	enc, _ := json.Marshal(s)
	b.Write(enc)
}

// writeJSONValue writes the JSON encoding of a field's value to b.
//
// Parameters:
//   - `b`: the strings.Builder to write to
//   - `f`: the field to write
//
// Example:
//
//	var b strings.Builder
//	writeJSONValue(&b, &Field{typ: StringType, strVal: "hello"})
func writeJSONValue(b *strings.Builder, f *Field) {
	switch f.typ {
	case StringType, JSONType:
		writeJSONString(b, f.strVal)
	case Int64Type, Int8Type, Int16Type, Int32Type:
		b.WriteString(itoa64(f.intVal))
	case UintType, Uint8Type, Uint16Type, Uint32Type, Uint64Type:
		b.WriteString(utoa64(f.uint64Val))
	case Float64Type:
		enc, _ := json.Marshal(f.floatVal)
		b.Write(enc)
	case Float32Type:
		enc, _ := json.Marshal(float32(f.floatVal))
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
	case TimefType:
		layout := f.strVal
		if strutil.IsEmpty(layout) {
			layout = time.RFC3339
		}
		writeJSONString(b, f.timeVal.Format(layout))
	case DurationType:
		writeJSONString(b, f.durVal.String())
	case AnyType:
		if s, ok := f.anyVal.(string); ok {
			writeJSONString(b, s)
		} else {
			enc, err := json.Marshal(f.anyVal)
			if err != nil {
				writeJSONString(b, f.Value())
			} else {
				b.Write(enc)
			}
		}
	default:
		writeJSONString(b, f.Value())
	}
}

// itoa converts an int to its decimal string representation.
//
// Parameters:
//   - `n`: the int to convert
//
// Returns:
//
// the decimal string representation of n
//
// Example:
//
//	itoa(123) // "123"
func itoa(n int) string {
	return itoa64(int64(n))
}

// utoa64 converts a uint64 to its decimal string representation.
//
// Parameters:
//   - `n`: the uint64 to convert
//
// Returns:
//
// the decimal string representation of n
//
// Example:
//
//	utoa64(123) // "123"
func utoa64(n uint64) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}

// itoa64 converts an int64 to its decimal string representation.
//
// Parameters:
//   - `n`: the int64 to convert
//
// Returns:
//
// the decimal string representation of n
//
// Example:
//
//	itoa64(123) // "123"
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

// istty reports whether w is connected to a terminal (character device).
//
// Parameters:
//   - `w`: the writer to test
//
// Returns:
//
// true when w is an *os.File whose device mode includes os.ModeCharDevice.
func istty(w io.Writer) bool {
	return sysx.IsTTY(w)
}
