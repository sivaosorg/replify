package strchain

import (
	"fmt"
	"strings"
	"sync"
)

// Compile-time interface satisfaction check.
var _ Weaver = (*SafeStringWeaver)(nil)

// SafeStringWeaver wraps strings.Builder with a fluent API for chainable operations.
// This implementation is THREAD-SAFE and can be used concurrently by multiple goroutines.
// Each method acquires a mutex lock to ensure safe concurrent access.
//
// For single-threaded usage where maximum performance is desired, use StringWeaver instead.
//
// Example:
//
//	sw := strchain.NewSafe().Append("Hello").Space().Append("World")
//	fmt.Println(sw.Build()) // Output: Hello World
type SafeStringWeaver struct {
	builder strings.Builder
	mu      sync.Mutex
}

// SafeFrom creates a new SafeStringWeaver initialized with the given string.
//
// Example:
//
//	sw := strchain.SafeFrom("Initial")
func SafeFrom(s string) *SafeStringWeaver {
	sw := NewSafe()
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.builder.WriteString(s)
	return sw
}

// NewSafe creates a new SafeStringWeaver instance.
//
// Example:
//
//	sw := strchain.NewSafe()
func NewSafe() *SafeStringWeaver {
	return &SafeStringWeaver{}
}

// NewSafeWithCapacity creates a new SafeStringWeaver with an initial capacity hint.
//
// Example:
//
//	sw := strchain.NewSafeWithCapacity(1024)
func NewSafeWithCapacity(capacity int) *SafeStringWeaver {
	sw := &SafeStringWeaver{}
	sw.builder.Grow(capacity)
	return sw
}

// Append adds a string and returns the builder for chaining.
//
// Example:
//
//	sw.Append("hello").Append(" world")
func (sw *SafeStringWeaver) Append(s string) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.builder.WriteString(s)
	return sw
}

// AppendF adds a formatted string (printf-style) and returns the builder.
//
// Example:
//
//	sw.AppendF("value: %d", 42)
func (sw *SafeStringWeaver) AppendF(format string, args ...any) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	fmt.Fprintf(&sw.builder, format, args...)
	return sw
}

// AppendByte adds a single byte and returns the builder.
//
// Example:
//
//	sw.AppendByte('a')
func (sw *SafeStringWeaver) AppendByte(b byte) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.builder.WriteByte(b)
	return sw
}

// AppendRune adds a single rune and returns the builder.
//
// Example:
//
//	sw.AppendRune('⌘')
func (sw *SafeStringWeaver) AppendRune(r rune) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.builder.WriteRune(r)
	return sw
}

// AppendBytes adds a byte slice and returns the builder.
//
// Example:
//
//	sw.AppendBytes([]byte("data"))
func (sw *SafeStringWeaver) AppendBytes(b []byte) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.builder.Write(b)
	return sw
}

// AppendInt adds an integer and returns the builder.
//
// Example:
//
//	sw.AppendInt(123)
func (sw *SafeStringWeaver) AppendInt(i int) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	fmt.Fprintf(&sw.builder, "%d", i)
	return sw
}

// AppendInt8 adds an int8 and returns the builder.
//
// Example:
//
//	sw.AppendInt8(8)
func (sw *SafeStringWeaver) AppendInt8(i int8) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	fmt.Fprintf(&sw.builder, "%d", i)
	return sw
}

// AppendInt16 adds an int16 and returns the builder.
//
// Example:
//
//	sw.AppendInt16(16)
func (sw *SafeStringWeaver) AppendInt16(i int16) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	fmt.Fprintf(&sw.builder, "%d", i)
	return sw
}

// AppendInt32 adds an int32 and returns the builder.
//
// Example:
//
//	sw.AppendInt32(32)
func (sw *SafeStringWeaver) AppendInt32(i int32) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	fmt.Fprintf(&sw.builder, "%d", i)
	return sw
}

// AppendInt64 adds an int64 and returns the builder.
//
// Example:
//
//	sw.AppendInt64(64)
func (sw *SafeStringWeaver) AppendInt64(i int64) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	fmt.Fprintf(&sw.builder, "%d", i)
	return sw
}

// AppendUint adds a uint and returns the builder.
//
// Example:
//
//	sw.AppendUint(123)
func (sw *SafeStringWeaver) AppendUint(i uint) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	fmt.Fprintf(&sw.builder, "%d", i)
	return sw
}

// AppendUint8 adds a uint8 and returns the builder.
//
// Example:
//
//	sw.AppendUint8(8)
func (sw *SafeStringWeaver) AppendUint8(i uint8) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	fmt.Fprintf(&sw.builder, "%d", i)
	return sw
}

// AppendUint16 adds a uint16 and returns the builder.
//
// Example:
//
//	sw.AppendUint16(16)
func (sw *SafeStringWeaver) AppendUint16(i uint16) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	fmt.Fprintf(&sw.builder, "%d", i)
	return sw
}

// AppendUint32 adds a uint32 and returns the builder.
//
// Example:
//
//	sw.AppendUint32(32)
func (sw *SafeStringWeaver) AppendUint32(i uint32) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	fmt.Fprintf(&sw.builder, "%d", i)
	return sw
}

// AppendUint64 adds a uint64 and returns the builder.
//
// Example:
//
//	sw.AppendUint64(64)
func (sw *SafeStringWeaver) AppendUint64(i uint64) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	fmt.Fprintf(&sw.builder, "%d", i)
	return sw
}

// AppendUintptr adds a uintptr and returns the builder.
//
// Example:
//
//	sw.AppendUintptr(0xdeadbeef)
func (sw *SafeStringWeaver) AppendUintptr(i uintptr) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	fmt.Fprintf(&sw.builder, "%d", i)
	return sw
}

// AppendFloat32 adds a float32 and returns the builder.
//
// Example:
//
//	sw.AppendFloat32(3.14)
func (sw *SafeStringWeaver) AppendFloat32(f float32) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	fmt.Fprintf(&sw.builder, "%f", f)
	return sw
}

// AppendFloat64 adds a float64 and returns the builder.
//
// Example:
//
//	sw.AppendFloat64(3.14159)
func (sw *SafeStringWeaver) AppendFloat64(f float64) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	fmt.Fprintf(&sw.builder, "%f", f)
	return sw
}

// AppendBool adds a boolean value and returns the builder.
//
// Example:
//
//	sw.AppendBool(true)
func (sw *SafeStringWeaver) AppendBool(b bool) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	fmt.Fprintf(&sw.builder, "%t", b)
	return sw
}

// Space adds a single space character and returns the builder.
//
// Example:
//
//	sw.Append("Hello").Space().Append("World")
func (sw *SafeStringWeaver) Space() Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.builder.WriteByte(' ')
	return sw
}

// Spaces adds n space characters and returns the builder.
//
// Example:
//
//	sw.Append("Key:").Spaces(4).Append("Value")
func (sw *SafeStringWeaver) Spaces(n int) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	for i := 0; i < n; i++ {
		sw.builder.WriteByte(' ')
	}
	return sw
}

// Tab adds a single tab character and returns the builder.
//
// Example:
//
//	sw.Append("Column1").Tab().Append("Column2")
func (sw *SafeStringWeaver) Tab() Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.builder.WriteByte('\t')
	return sw
}

// Tabs adds n tab characters and returns the builder.
//
// Example:
//
//	sw.Tabs(2).Append("Indented text")
func (sw *SafeStringWeaver) Tabs(n int) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	for i := 0; i < n; i++ {
		sw.builder.WriteByte('\t')
	}
	return sw
}

// NewLine adds a single newline character and returns the builder.
//
// Example:
//
//	sw.Append("Line 1").NewLine().Append("Line 2")
func (sw *SafeStringWeaver) NewLine() Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.builder.WriteByte('\n')
	return sw
}

// NewLines adds n newline characters and returns the builder.
//
// Example:
//
//	sw.Append("Paragraph 1").NewLines(2).Append("Paragraph 2")
func (sw *SafeStringWeaver) NewLines(n int) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	for i := 0; i < n; i++ {
		sw.builder.WriteByte('\n')
	}
	return sw
}

// Line adds a string followed by a newline and returns the builder.
//
// Example:
//
//	sw.Line("First Line").Line("Second Line")
func (sw *SafeStringWeaver) Line(s string) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.builder.WriteString(s)
	sw.builder.WriteByte('\n')
	return sw
}

// LineF adds a formatted string followed by a newline and returns the builder.
//
// Example:
//
//	sw.LineF("User ID: %d", 123)
func (sw *SafeStringWeaver) LineF(format string, args ...any) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	fmt.Fprintf(&sw.builder, format, args...)
	sw.builder.WriteByte('\n')
	return sw
}

// Repeat adds a string n times and returns the builder.
//
// Example:
//
//	sw.Repeat("-", 10) // adds "----------"
func (sw *SafeStringWeaver) Repeat(s string, n int) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	for i := 0; i < n; i++ {
		sw.builder.WriteString(s)
	}
	return sw
}

// Join adds strings with a separator and returns the builder.
//
// Example:
//
//	sw.Join(", ", "Apple", "Banana", "Cherry")
func (sw *SafeStringWeaver) Join(sep string, elements ...string) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	for i, elem := range elements {
		if i > 0 {
			sw.builder.WriteString(sep)
		}
		sw.builder.WriteString(elem)
	}
	return sw
}

// AppendIf conditionally adds a string if the condition is true.
//
// Example:
//
//	sw.AppendIf(isValid, " Validated")
func (sw *SafeStringWeaver) AppendIf(condition bool, s string) Weaver {
	if condition {
		sw.mu.Lock()
		defer sw.mu.Unlock()
		sw.builder.WriteString(s)
	}
	return sw
}

// AppendIfF conditionally adds a formatted string if the condition is true.
//
// Example:
//
//	sw.AppendIfF(debugMode, "Debug: %s", msg)
func (sw *SafeStringWeaver) AppendIfF(condition bool, format string, args ...any) Weaver {
	if condition {
		sw.mu.Lock()
		defer sw.mu.Unlock()
		fmt.Fprintf(&sw.builder, format, args...)
	}
	return sw
}

// When executes a function on the builder if the condition is true.
//
// Example:
//
//	sw.When(isHeader, func(s *SafeStringWeaver) {
//	    s.Append("# ").Line(title)
//	})
func (sw *SafeStringWeaver) When(condition bool, fn func(*SafeStringWeaver)) *SafeStringWeaver {
	if condition {
		fn(sw)
	}
	return sw
}

// Unless executes a function on the builder if the condition is false.
//
// Example:
//
//	sw.Unless(isMinified, func(s *SafeStringWeaver) {
//	    s.NewLine().Indent(1, "")
//	})
func (sw *SafeStringWeaver) Unless(condition bool, fn func(*SafeStringWeaver)) *SafeStringWeaver {
	if !condition {
		fn(sw)
	}
	return sw
}

// Each iterates over a slice and applies a function for each element.
//
// Example:
//
//	sw.Each([]string{"a", "b", "c"}, func(s *SafeStringWeaver, item string) {
//	    s.Append(item).Comma()
//	})
func (sw *SafeStringWeaver) Each(items []string, fn func(*SafeStringWeaver, string)) *SafeStringWeaver {
	for _, item := range items {
		fn(sw, item)
	}
	return sw
}

// Indent adds indentation (2 spaces per level) before appending text.
//
// Example:
//
//	sw.Indent(2, "Sub-item")
func (sw *SafeStringWeaver) Indent(level int, s string) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	for i := 0; i < level*2; i++ {
		sw.builder.WriteByte(' ')
	}
	sw.builder.WriteString(s)
	return sw
}

// IndentLine adds indentation (2 spaces per level) before text and ends with a newline.
//
// Example:
//
//	sw.IndentLine(1, "Point 1")
func (sw *SafeStringWeaver) IndentLine(level int, s string) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	for i := 0; i < level*2; i++ {
		sw.builder.WriteByte(' ')
	}
	sw.builder.WriteString(s)
	sw.builder.WriteByte('\n')
	return sw
}

// Wrap wraps text with a prefix and suffix.
//
// Example:
//
//	sw.Wrap("**", "Bold text", "**")
func (sw *SafeStringWeaver) Wrap(prefix, text, suffix string) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.builder.WriteString(prefix)
	sw.builder.WriteString(text)
	sw.builder.WriteString(suffix)
	return sw
}

// Quote wraps text in double quotes.
//
// Example:
//
//	sw.Quote("quoted string") // adds "quoted string"
func (sw *SafeStringWeaver) Quote(s string) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.builder.WriteByte('"')
	sw.builder.WriteString(s)
	sw.builder.WriteByte('"')
	return sw
}

// SingleQuote wraps text in single quotes.
//
// Example:
//
//	sw.SingleQuote("char") // adds 'char'
func (sw *SafeStringWeaver) SingleQuote(s string) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.builder.WriteByte('\'')
	sw.builder.WriteString(s)
	sw.builder.WriteByte('\'')
	return sw
}

// Parenthesize wraps text in parentheses.
//
// Example:
//
//	sw.Parenthesize("expression") // adds (expression)
func (sw *SafeStringWeaver) Parenthesize(s string) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.builder.WriteByte('(')
	sw.builder.WriteString(s)
	sw.builder.WriteByte(')')
	return sw
}

// Bracket wraps text in square brackets.
//
// Example:
//
//	sw.Bracket("index") // adds [index]
func (sw *SafeStringWeaver) Bracket(s string) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.builder.WriteByte('[')
	sw.builder.WriteString(s)
	sw.builder.WriteByte(']')
	return sw
}

// Brace wraps text in curly braces.
//
// Example:
//
//	sw.Brace("struct definition") // adds {struct definition}
func (sw *SafeStringWeaver) Brace(s string) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.builder.WriteByte('{')
	sw.builder.WriteString(s)
	sw.builder.WriteByte('}')
	return sw
}

// Comma adds a comma to the builder.
//
// Example:
//
//	sw.Append("item1").Comma().Append("item2")
func (sw *SafeStringWeaver) Comma() Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.builder.WriteByte(',')
	return sw
}

// Dot adds a period to the builder.
//
// Example:
//
//	sw.Append("End of sentence").Dot()
func (sw *SafeStringWeaver) Dot() Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.builder.WriteByte('.')
	return sw
}

// Colon adds a colon to the builder.
//
// Example:
//
//	sw.Append("Label").Colon().Space().Append("Value")
func (sw *SafeStringWeaver) Colon() Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.builder.WriteByte(':')
	return sw
}

// Semicolon adds a semicolon to the builder.
//
// Example:
//
//	sw.Append("x = 1").Semicolon()
func (sw *SafeStringWeaver) Semicolon() Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.builder.WriteByte(';')
	return sw
}

// Equals adds an equals sign to the builder.
//
// Example:
//
//	sw.Append("key").Equals().Append("value")
func (sw *SafeStringWeaver) Equals() Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.builder.WriteByte('=')
	return sw
}

// Arrow adds an arrow (->) to the builder.
//
// Example:
//
//	sw.Append("source").Arrow().Append("target")
func (sw *SafeStringWeaver) Arrow() Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.builder.WriteString("->")
	return sw
}

// FatArrow adds a fat arrow (=>) to the builder.
//
// Example:
//
//	sw.Append("map").FatArrow().Append("result")
func (sw *SafeStringWeaver) FatArrow() Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.builder.WriteString("=>")
	return sw
}

// Grow grows the builder's capacity by n bytes.
//
// Example:
//
//	sw.Grow(512)
func (sw *SafeStringWeaver) Grow(n int) Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.builder.Grow(n)
	return sw
}

// Reset clears the builder and returns it for reuse.
//
// Example:
//
//	sw.Reset().Append("Fresh start")
func (sw *SafeStringWeaver) Reset() Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.builder.Reset()
	return sw
}

// Len returns the current length of the built string.
//
// Example:
//
//	length := sw.Len()
func (sw *SafeStringWeaver) Len() int {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.builder.Len()
}

// Cap returns the current capacity of the builder.
//
// Example:
//
//	capacity := sw.Cap()
func (sw *SafeStringWeaver) Cap() int {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.builder.Cap()
}

// String returns the final built string.
//
// Example:
//
//	result := sw.String()
func (sw *SafeStringWeaver) String() string {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.builder.String()
}

// Build is an alias for String() for fluent API consistency.
//
// Example:
//
//	result := sw.Append("data").Build()
func (sw *SafeStringWeaver) Build() string {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.builder.String()
}

// Clone creates a thread-safe copy of the current builder state.
//
// Example:
//
//	newSw := sw.Clone()
func (sw *SafeStringWeaver) Clone() Weaver {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	clone := NewSafe()
	clone.builder.WriteString(sw.builder.String())
	return clone
}

// Inspect executes a function with access to the current state without modification.
//
// Example:
//
//	sw.Inspect(func(current string) {
//	    log.Printf("Current state: %s", current)
//	})
func (sw *SafeStringWeaver) Inspect(fn func(current string)) Weaver {
	sw.mu.Lock()
	current := sw.builder.String()
	sw.mu.Unlock()
	fn(current)
	return sw
}

// Builder returns the underlying strings.Builder for advanced operations.
//
// Example:
//
//	builder := sw.Builder()
func (sw *SafeStringWeaver) Builder() *strings.Builder {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return &sw.builder
}
