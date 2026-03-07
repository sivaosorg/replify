package slogger

import (
"archive/zip"
"encoding/json"
"io"
"os"
"path/filepath"
"strings"
"time"
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

// trimFile shortens an absolute file path to just the last two path segments.
// For example "/home/user/project/pkg/foo/bar.go" becomes "pkg/foo/bar.go".
func trimFile(path string) string {
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

// defaultOptions returns production-ready defaults:
// InfoLevel, TextFormatter writing to os.Stderr, no caller capture.
//
// Returns:
//
// an *Options pre-populated with sensible production settings.
func defaultOptions() *Options {
out := os.Stderr
return &Options{
Level:     InfoLevel,
Output:    out,
Formatter: NewTextFormatter(out),
}
}

// ///////////////////////////////////////////////////////////////////////////
// Sampling helpers
// ///////////////////////////////////////////////////////////////////////////

// newSampler creates a sampler with the given options.
//
// Parameters:
//   - `opts`: the sampling configuration
//
// Returns:
//
// a ready-to-use *sampler.
func newSampler(opts SamplingOptions) *sampler {
return &sampler{opts: opts}
}

// ///////////////////////////////////////////////////////////////////////////
// Rotation helpers
// ///////////////////////////////////////////////////////////////////////////

// newRotatingFile creates and opens a rotating log file for the given level.
func newRotatingFile(opts RotationOptions, level Level) (*rotatingFile, error) {
rf := &rotatingFile{
path:     filepath.Join(opts.Dir, levelFileName(level)),
maxBytes: opts.MaxBytes,
maxAge:   opts.MaxAge,
compress: opts.Compress,
dir:      opts.Dir,
level:    level,
}
if err := rf.open(); err != nil {
return nil, err
}
return rf, nil
}

// levelFileName returns the log filename for the given level using the
// package-level file name constants.
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
func levelPad(l Level) string {
s := l.String()
for len(s) < levelPadWidth {
s += " "
}
return s
}

// needsQuoting reports whether a value string should be quoted in text output.
// Strings that contain spaces, equals signs, double-quotes, or backslashes
// are wrapped in strconv.Quote to preserve round-trip fidelity.
func needsQuoting(s string) bool {
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

// itoa converts an int to its decimal string representation.
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
