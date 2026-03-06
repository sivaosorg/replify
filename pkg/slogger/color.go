package slogger

import (
	"io"
	"os"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorGreen  = "\033[32m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

// IsTTY reports whether w is connected to a terminal (character device).
//
// Parameters:
//   - `w`: the writer to test
//
// Returns:
//
// true when w is an *os.File whose device mode includes os.ModeCharDevice.
func IsTTY(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}

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
