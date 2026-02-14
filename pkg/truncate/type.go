package truncate

// Strategy is an interface for truncation strategies.
// Any type that implements this interface can be used with the Apply
// convenience function to truncate strings in a uniform manner.
//
// Example:
//
//	var strategy truncate.Strategy = truncate.NewCutEllipsisStrategy()
//	result := strategy.Truncate("Hello, World!", 8) // "Hello, …"
type Strategy interface {
	// Truncate applies the truncation strategy to the given string,
	// returning a new string whose visual length does not exceed the
	// specified number of runes.
	//
	// Parameters:
	//   - str: The input string to truncate.
	//   - length: The maximum number of runes allowed in the result.
	//
	// Returns:
	//   - A truncated string that fits within the specified length.
	Truncate(str string, length int) string
}

// Truncator is the core engine that performs string truncation according to
// its configured omission marker, position, and maximum length.
//
// A Truncator is created through the fluent builder API:
//
//	t := truncate.NewTruncator().
//	    WithOmission("...").
//	    WithPosition(truncate.PositionMiddle).
//	    WithMaxLength(20).
//	    Build()
//	result := t.Truncate("A very long string that needs truncation")
//
// Fields:
//   - omission: The string inserted where characters are removed (e.g. "…").
//   - position: Where the omission marker is placed (start, middle, or end).
//   - maxLength: The maximum number of runes the result may contain.
type Truncator struct {
	omission  string
	position  Position
	maxLength int
}

// TruncatorBuilder provides a fluent interface for constructing a Truncator.
// It allows chaining configuration methods before calling Build to produce
// an immutable Truncator instance.
//
// Example:
//
//	builder := truncate.NewTruncator().
//	    WithOmission("...").
//	    WithPosition(truncate.PositionEnd).
//	    WithMaxLength(15)
//	t := builder.Build()
type TruncatorBuilder struct {
	omission  string
	position  Position
	maxLength int
}

// CutStrategy simply truncates the string to the desired length without
// adding any omission marker. Characters beyond the maximum length are
// silently dropped from the end.
//
// Example:
//
//	s := truncate.CutStrategy{}
//	result := s.Truncate("Hello, World!", 5) // "Hello"
type CutStrategy struct{}

// CutEllipsisStrategy truncates the string from the end and appends the
// default ellipsis omission marker (DefaultOmission).
//
// Example:
//
//	s := truncate.CutEllipsisStrategy{}
//	result := s.Truncate("Hello, World!", 8) // "Hello, …"
type CutEllipsisStrategy struct{}

// CutEllipsisLeadingStrategy truncates the string from the start and
// prepends the default ellipsis omission marker (DefaultOmission).
//
// Example:
//
//	s := truncate.CutEllipsisLeadingStrategy{}
//	result := s.Truncate("Hello, World!", 8) // "…World!"
type CutEllipsisLeadingStrategy struct{}

// EllipsisMiddleStrategy truncates the string from the middle and inserts
// the default ellipsis omission marker (DefaultOmission) between the
// preserved head and tail.
//
// Example:
//
//	s := truncate.EllipsisMiddleStrategy{}
//	result := s.Truncate("Hello, World!", 8) // "Hel…ld!"
type EllipsisMiddleStrategy struct{}
