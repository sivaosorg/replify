package strchain

// Weaver defines the common interface for fluent string builder operations.
// Both StringWeaver (non-thread-safe, maximum performance) and SafeStringWeaver
// (thread-safe with mutex synchronization) implement this interface.
//
// Use Weaver when you need polymorphism — for example, a function that accepts
// either a StringWeaver or a SafeStringWeaver as a parameter.
//
// Methods that accept callback functions with concrete receiver types
// (When, Unless, Each) are not included in this interface because their
// function parameter types differ between StringWeaver and SafeStringWeaver.
// These methods are available only on the concrete types.
//
// Example:
//
//	var w strchain.Weaver = strchain.New()
//	result := w.Append("Hello").Space().Append("World").Build()
//
//	var sw strchain.Weaver = strchain.NewSafe()
//	result = sw.Append("Hello").Space().Append("World").Build()
type Weaver interface {

	// Append adds a string and returns the builder for chaining.
	Append(s string) Weaver

	// AppendF adds a formatted string (printf-style) and returns the builder.
	AppendF(format string, args ...any) Weaver

	// AppendByte adds a single byte and returns the builder.
	AppendByte(b byte) Weaver

	// AppendRune adds a single rune and returns the builder.
	AppendRune(r rune) Weaver

	// AppendBytes adds a byte slice and returns the builder.
	AppendBytes(b []byte) Weaver

	// AppendInt adds an integer and returns the builder.
	AppendInt(i int) Weaver

	// AppendInt8 adds an int8 and returns the builder.
	AppendInt8(i int8) Weaver

	// AppendInt16 adds an int16 and returns the builder.
	AppendInt16(i int16) Weaver

	// AppendInt32 adds an int32 and returns the builder.
	AppendInt32(i int32) Weaver

	// AppendInt64 adds an int64 and returns the builder.
	AppendInt64(i int64) Weaver

	// AppendUint adds a uint and returns the builder.
	AppendUint(i uint) Weaver

	// AppendUint8 adds a uint8 and returns the builder.
	AppendUint8(i uint8) Weaver

	// AppendUint16 adds a uint16 and returns the builder.
	AppendUint16(i uint16) Weaver

	// AppendUint32 adds a uint32 and returns the builder.
	AppendUint32(i uint32) Weaver

	// AppendUint64 adds a uint64 and returns the builder.
	AppendUint64(i uint64) Weaver

	// AppendUintptr adds a uintptr and returns the builder.
	AppendUintptr(i uintptr) Weaver

	// AppendFloat32 adds a float32 and returns the builder.
	AppendFloat32(f float32) Weaver

	// AppendFloat64 adds a float64 and returns the builder.
	AppendFloat64(f float64) Weaver

	// AppendBool adds a boolean value and returns the builder.
	AppendBool(b bool) Weaver

	// Space adds a single space character and returns the builder.
	Space() Weaver

	// Spaces adds n space characters and returns the builder.
	Spaces(n int) Weaver

	// Tab adds a single tab character and returns the builder.
	Tab() Weaver

	// Tabs adds n tab characters and returns the builder.
	Tabs(n int) Weaver

	// NewLine adds a single newline character and returns the builder.
	NewLine() Weaver

	// NewLines adds n newline characters and returns the builder.
	NewLines(n int) Weaver

	// Line adds a string followed by a newline and returns the builder.
	Line(s string) Weaver

	// LineF adds a formatted string followed by a newline and returns the builder.
	LineF(format string, args ...any) Weaver

	// Repeat adds a string n times and returns the builder.
	Repeat(s string, n int) Weaver

	// Join adds strings with a separator and returns the builder.
	Join(sep string, elements ...string) Weaver

	// AppendIf conditionally adds a string if the condition is true.
	AppendIf(condition bool, s string) Weaver

	// AppendIfF conditionally adds a formatted string if the condition is true.
	AppendIfF(condition bool, format string, args ...any) Weaver

	// Indent adds indentation (2 spaces per level) before appending text.
	Indent(level int, s string) Weaver

	// IndentLine adds indentation (2 spaces per level) before text and ends with a newline.
	IndentLine(level int, s string) Weaver

	// Wrap wraps text with a prefix and suffix.
	Wrap(prefix, text, suffix string) Weaver

	// Quote wraps text in double quotes.
	Quote(s string) Weaver

	// SingleQuote wraps text in single quotes.
	SingleQuote(s string) Weaver

	// Parenthesize wraps text in parentheses.
	Parenthesize(s string) Weaver

	// Bracket wraps text in square brackets.
	Bracket(s string) Weaver

	// Brace wraps text in curly braces.
	Brace(s string) Weaver

	// Comma adds a comma to the builder.
	Comma() Weaver

	// Dot adds a period to the builder.
	Dot() Weaver

	// Colon adds a colon to the builder.
	Colon() Weaver

	// Semicolon adds a semicolon to the builder.
	Semicolon() Weaver

	// Equals adds an equals sign to the builder.
	Equals() Weaver

	// Arrow adds an arrow (->) to the builder.
	Arrow() Weaver

	// FatArrow adds a fat arrow (=>) to the builder.
	FatArrow() Weaver

	// Grow grows the builder's capacity by n bytes.
	Grow(n int) Weaver

	// Reset clears the builder and returns it for reuse.
	Reset() Weaver

	// Len returns the current length of the built string.
	Len() int

	// Cap returns the current capacity of the builder.
	Cap() int

	// String returns the final built string.
	String() string

	// Build is an alias for String() for fluent API consistency.
	Build() string

	// Clone creates an independent copy of the current builder state.
	Clone() Weaver

	// Inspect executes a function with access to the current state without modification.
	Inspect(fn func(current string)) Weaver
}
