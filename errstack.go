package replify

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strconv"
	"strings"
)

// stack.StackTrace converts the stack of program counters to a StackTrace.
//
// Usage:
// This method creates a more developer-friendly StackTrace type from the raw program counters.
//
// Example:
//
//	st := callers()
//	trace := st.StackTrace()
//	fmt.Printf("%+v", trace)
func (s *stack) StackTrace() StackTrace {
	f := make([]Frame, len(*s))
	for i := 0; i < len(f); i++ {
		f[i] = Frame((*s)[i])
	}
	return f
}

// Frame.MarshalText formats a Frame as a text string.
// The output is the same as fmt.Sprintf("%+v", f), but without newlines or tabs.
//
// Usage:
// Converts the Frame to a compact text representation, suitable for logging or serialization.
//
// Example:
//
//	frame := Frame(somePC)
//	text, err := frame.MarshalText()
//	if err == nil {
//	    fmt.Println(string(text))
//	}
func (f Frame) MarshalText() ([]byte, error) {
	name := f.name()
	if name == ErrUnknown {
		return []byte(name), nil
	}
	return []byte(fmt.Sprintf("%s %s:%d", name, f.file(), f.line())), nil
}

// Frame represents a program counter inside a stack frame.
// For historical reasons, if Frame is interpreted as a uintptr,
// its value represents the program counter + 1.

// StackTrace represents a stack of Frames from innermost (newest) to outermost (oldest).

// stack is a slice of uintptrs representing program counters.

// fundamental is an error type that holds a message and a stack trace,
// but does not include a specific caller.

// Frame.Format formats the frame according to the fmt.Formatter interface.
//
// Usage:
// The `verb` parameter controls the formatting output:
//   - %s: Source file name.
//   - %d: Source line number.
//   - %n: Function name.
//   - %v: Equivalent to %s:%d.
//
// Flags:
//   - %+s: Includes function name and path of the source file relative to the compile-time GOPATH.
//   - %+v: Combines %+s and %d (function name, source path, and line number).
//
// Example:
//
//	frame := Frame(somePC)
//	fmt.Printf("%+v", frame)
func (f Frame) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			io.WriteString(s, f.name())
			io.WriteString(s, "\n\t")
			io.WriteString(s, f.file())
		default:
			io.WriteString(s, path.Base(f.file()))
		}
	case 'd':
		io.WriteString(s, strconv.Itoa(f.line()))
	case 'n':
		io.WriteString(s, f.shortName(f.name()))
	case 'v':
		f.Format(s, 's')
		io.WriteString(s, ":")
		f.Format(s, 'd')
	}
}

// StackTrace.Format formats the stack trace according to the fmt.Formatter interface.
//
// Usage:
// The `verb` parameter controls the formatting output:
//   - %s: Lists source files for each Frame in the stack.
//   - %v: Lists source file and line number for each Frame in the stack.
//
// Flags:
//   - %+v: Prints filename, function name, and line number for each Frame.
//
// Example:
//
//	trace := StackTrace{frame1, frame2}
//	fmt.Printf("%+v", trace)
func (st StackTrace) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case s.Flag('+'):
			for _, f := range st {
				io.WriteString(s, "\n")
				f.Format(s, verb)
			}
		case s.Flag('#'):
			fmt.Fprintf(s, "%#v", []Frame(st))
		default:
			st.fmtSlice(s, verb)
		}
	case 's':
		st.fmtSlice(s, verb)
	}
}

// stack.Format formats the stack of program counters.
//
// Usage:
// Similar to StackTrace.Format, but operates on the raw program counters (`stack` type).
// Use %+v to print each Frame in the stack trace with details.
//
// Example:
//
//	st := callers()
//	fmt.Printf("%+v", st)
func (s *stack) Format(st fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case st.Flag('+'):
			for _, pc := range *s {
				f := Frame(pc)
				fmt.Fprintf(st, "\n%+v", f)
			}
		}
	}
}

// Callers captures the current call stack as a stack of program counters.
//
// Usage:
// Use this function to capture the stack trace of the current execution context.
//
// Example:
//
//	st := Callers()
//	trace := st.StackTrace()
//	fmt.Printf("%+v", trace)
func Callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stack = pcs[0:n]
	return &st
}

// Frame.pc calculates the actual program counter for the Frame.
// This compensates for the `+1` adjustment used historically.
//
// Usage:
// Used internally for resolving function and file details.
func (f Frame) pc() uintptr { return uintptr(f) - 1 }

// Frame.file returns the full path of the source file for the Frame.
//
// Usage:
// Retrieves the path to the source file containing the function corresponding to the Frame's program counter.
//
// Example:
//
//	filePath := frame.file()
func (f Frame) file() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return ErrUnknown
	}
	file, _ := fn.FileLine(f.pc())
	return file
}

// Frame.line returns the line number of the source file for the Frame.
//
// Usage:
// Retrieves the line number in the source file where the function resides.
//
// Example:
//
//	line := frame.line()
func (f Frame) line() int {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(f.pc())
	return line
}

// Frame.name returns the function name for the Frame.
//
// Usage:
// Resolves the name of the function corresponding to the Frame's program counter.
//
// Example:
//
//	funcName := frame.name()
func (f Frame) name() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	return fn.Name()
}

// shortName extracts the function name without its path prefix.
//
// Usage:
// Useful for displaying function names in a compact format.
//
// Example:
//
//	shortName := shortName("path/to/package.FunctionName")
func (f Frame) shortName(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}

// StackTrace.fmtSlice formats the StackTrace as a slice of Frames for `%s` or `%v`.
//
// Usage:
// Used internally by the Format method of StackTrace for compact output.
//
// Example:
// (Indirect usage through fmt.Printf or similar functions)
func (st StackTrace) fmtSlice(s fmt.State, verb rune) {
	io.WriteString(s, "[")
	for i, f := range st {
		if i > 0 {
			io.WriteString(s, " ")
		}
		f.Format(s, verb)
	}
	io.WriteString(s, "]")
}
