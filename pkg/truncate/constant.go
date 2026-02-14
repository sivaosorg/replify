// Package truncate provides a set of OOP-style strategies to truncate strings
// with support for customisable omission markers and positional control.
package truncate

// Position represents the location at which the omission marker is placed
// when truncating a string. It determines which part of the original string
// is preserved and which part is cut off.
type Position int

const (
	// PositionEnd truncates from the end of the string, preserving the leading
	// characters and appending the omission marker at the tail.
	//
	// Example:
	//
	//	"Hello, World!" → "Hello…"
	PositionEnd Position = iota

	// PositionStart truncates from the start of the string, preserving the
	// trailing characters and prepending the omission marker at the front.
	//
	// Example:
	//
	//	"Hello, World!" → "…orld!"
	PositionStart

	// PositionMiddle truncates from the middle of the string, preserving
	// both the leading and trailing characters and inserting the omission
	// marker in between.
	//
	// Example:
	//
	//	"Hello, World!" → "Hel…d!"
	PositionMiddle
)

const (
	// DefaultOmission is the default omission marker used when truncating strings.
	// It is a single Unicode ellipsis character (U+2026).
	DefaultOmission = "…"
)
