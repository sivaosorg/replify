package truncate

import (
	"math"
	"unicode/utf8"
)

// Apply is a package-level convenience function that truncates a string using
// the provided Strategy. It delegates directly to the strategy's Truncate method.
//
// Parameters:
//   - str: The input string to truncate.
//   - length: The maximum number of runes allowed in the result.
//   - strategy: The truncation strategy to apply.
//
// Returns:
//   - A truncated string produced by the strategy.
//
// Example:
//
//	result := truncate.Apply("Hello, World!", 8, truncate.NewCutEllipsisStrategy())
//	// result: "Hello, …"
func Apply(str string, length int, strategy Strategy) string {
	return strategy.Truncate(str, length)
}

// Truncate applies the truncator's configured omission, position, and maxLength
// to produce a truncated version of the input string. If the string already fits
// within maxLength, it is returned unchanged.
//
// Parameters:
//   - str: The input string to truncate.
//
// Returns:
//   - A truncated string whose rune count does not exceed maxLength.
//
// Example:
//
//	t := NewTruncator().WithMaxLength(10).Build()
//	result := t.Truncate("Hello, World!") // "Hello, Wo…"
func (t *Truncator) Truncate(str string) string {
	return truncateCore(str, t.maxLength, t.omission, t.position)
}

// TruncateWithLength applies the truncator's configured omission and position
// to produce a truncated version of the input string using the given length
// instead of the pre-configured maxLength. This is useful when you want to
// reuse a Truncator with different length constraints.
//
// Parameters:
//   - str: The input string to truncate.
//   - length: The maximum number of runes allowed in the result.
//
// Returns:
//   - A truncated string whose rune count does not exceed the specified length.
//
// Example:
//
//	t := NewTruncator().WithOmission("...").Build()
//	result := t.TruncateWithLength("Hello, World!", 8) // "Hello..."
func (t *Truncator) TruncateWithLength(str string, length int) string {
	return truncateCore(str, length, t.omission, t.position)
}

// Truncate on CutStrategy truncates the string to the desired length without
// any omission marker, cutting from the end.
//
// Parameters:
//   - str: The input string to truncate.
//   - length: The maximum number of runes allowed in the result.
//
// Returns:
//   - A truncated string with no omission marker.
//
// Example:
//
//	s := CutStrategy{}
//	result := s.Truncate("Hello, World!", 5) // "Hello"
func (CutStrategy) Truncate(str string, length int) string {
	return truncateCore(str, length, "", PositionEnd)
}

// Truncate on CutEllipsisStrategy truncates the string from the end and
// appends the default ellipsis omission marker.
//
// Parameters:
//   - str: The input string to truncate.
//   - length: The maximum number of runes allowed in the result.
//
// Returns:
//   - A truncated string with the default ellipsis appended at the end.
//
// Example:
//
//	s := CutEllipsisStrategy{}
//	result := s.Truncate("Hello, World!", 8) // "Hello, …"
func (CutEllipsisStrategy) Truncate(str string, length int) string {
	return truncateCore(str, length, DefaultOmission, PositionEnd)
}

// Truncate on CutEllipsisLeadingStrategy truncates the string from the start
// and prepends the default ellipsis omission marker.
//
// Parameters:
//   - str: The input string to truncate.
//   - length: The maximum number of runes allowed in the result.
//
// Returns:
//   - A truncated string with the default ellipsis prepended at the start.
//
// Example:
//
//	s := CutEllipsisLeadingStrategy{}
//	result := s.Truncate("Hello, World!", 8) // "…World!"
func (CutEllipsisLeadingStrategy) Truncate(str string, length int) string {
	return truncateCore(str, length, DefaultOmission, PositionStart)
}

// Truncate on EllipsisMiddleStrategy truncates the string from the middle
// and inserts the default ellipsis omission marker between the preserved
// head and tail.
//
// Parameters:
//   - str: The input string to truncate.
//   - length: The maximum number of runes allowed in the result.
//
// Returns:
//   - A truncated string with the default ellipsis in the middle.
//
// Example:
//
//	s := EllipsisMiddleStrategy{}
//	result := s.Truncate("Hello, World!", 8) // "Hel…ld!"
func (EllipsisMiddleStrategy) Truncate(str string, length int) string {
	return truncateCore(str, length, DefaultOmission, PositionMiddle)
}

// truncateCore is the internal engine that implements all truncation logic.
// It is shared by every strategy and the Truncator struct to avoid code
// duplication. The function operates on runes to correctly handle multi-byte
// Unicode characters.
//
// Parameters:
//   - str: The input string to truncate.
//   - length: The maximum number of runes allowed in the result.
//   - omission: The omission marker string (may be empty for a plain cut).
//   - pos: The position at which the omission marker is placed.
//
// Returns:
//   - A truncated string that respects the given constraints.
func truncateCore(str string, length int, omission string, pos Position) string {
	if length < 1 {
		return ""
	}
	r := []rune(str)
	sLen := len(r)
	oLen := utf8.RuneCountInString(omission)
	// No truncation needed — string already fits.
	if length >= sLen {
		return str
	}
	// When the requested length is shorter than or equal to the omission
	// marker itself, fall back to a plain end-cut without the marker.
	if length <= oLen {
		return truncateEnd(r, length, "", 0)
	}
	switch pos {
	case PositionStart:
		return truncateStart(r, length, omission, oLen)
	case PositionMiddle:
		return truncateMiddle(r, length, omission, oLen)
	default:
		return truncateEnd(r, length, omission, oLen)
	}
}

// truncateStart preserves the trailing runes of the slice and prepends the
// omission marker, producing a result of exactly `length` runes.
//
// Parameters:
//   - r: The input string as a rune slice.
//   - length: The maximum number of runes allowed in the result.
//   - omission: The omission marker string.
//   - oLen: The rune count of the omission marker.
//
// Returns:
//   - A string with the omission marker at the start followed by the tail.
func truncateStart(r []rune, length int, omission string, oLen int) string {
	return omission + string(r[len(r)-length+oLen:])
}

// truncateEnd preserves the leading runes of the slice and appends the
// omission marker, producing a result of exactly `length` runes.
//
// Parameters:
//   - r: The input string as a rune slice.
//   - length: The maximum number of runes allowed in the result.
//   - omission: The omission marker string.
//   - oLen: The rune count of the omission marker.
//
// Returns:
//   - A string with the head followed by the omission marker.
func truncateEnd(r []rune, length int, omission string, oLen int) string {
	return string(r[:length-oLen]) + omission
}

// truncateMiddle preserves both the leading and trailing runes of the slice,
// inserting the omission marker in between. The split point is balanced so
// that the head and tail are as close in length as possible.
//
// If the requested length is too short to accommodate at least one character
// on each side of the omission marker, it falls back to a plain end-cut.
//
// Parameters:
//   - r: The input string as a rune slice.
//   - length: The maximum number of runes allowed in the result.
//   - omission: The omission marker string.
//   - oLen: The rune count of the omission marker.
//
// Returns:
//   - A string with the head, omission marker, and tail concatenated.
func truncateMiddle(r []rune, length int, omission string, oLen int) string {
	sLen := len(r)
	// Ensure at least one character on each side of the omission marker.
	if length < oLen+2 {
		return truncateEnd(r, length, "", oLen)
	}
	// Compute the number of characters to keep before the omission marker.
	// The parity of the original string length determines the rounding
	// direction so that the split stays visually balanced.
	var delta int
	if sLen%2 == 0 {
		delta = int(math.Ceil(float64(length-oLen) / 2))
	} else {
		delta = int(math.Floor(float64(length-oLen) / 2))
	}
	result := make([]rune, length)
	copy(result, r[0:delta])
	copy(result[delta:], []rune(omission))
	copy(result[delta+oLen:], r[sLen-length+oLen+delta:])
	return string(result)
}
