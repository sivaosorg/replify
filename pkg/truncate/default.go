package truncate

// NewCutStrategy creates and returns a CutStrategy, which truncates a string
// to the desired length without adding any omission marker.
//
// Returns:
//   - A Strategy that performs a plain cut.
//
// Example:
//
//	s := truncate.NewCutStrategy()
//	result := s.Truncate("Hello, World!", 5) // "Hello"
func NewCutStrategy() Strategy {
	return CutStrategy{}
}

// NewCutEllipsisStrategy creates and returns a CutEllipsisStrategy, which
// truncates a string from the end and appends the default ellipsis marker.
//
// Returns:
//   - A Strategy that cuts from the end with an ellipsis.
//
// Example:
//
//	s := truncate.NewCutEllipsisStrategy()
//	result := s.Truncate("Hello, World!", 8) // "Hello, …"
func NewCutEllipsisStrategy() Strategy {
	return CutEllipsisStrategy{}
}

// NewCutEllipsisLeadingStrategy creates and returns a CutEllipsisLeadingStrategy,
// which truncates a string from the start and prepends the default ellipsis marker.
//
// Returns:
//   - A Strategy that cuts from the start with an ellipsis.
//
// Example:
//
//	s := truncate.NewCutEllipsisLeadingStrategy()
//	result := s.Truncate("Hello, World!", 8) // "…World!"
func NewCutEllipsisLeadingStrategy() Strategy {
	return CutEllipsisLeadingStrategy{}
}

// NewEllipsisMiddleStrategy creates and returns an EllipsisMiddleStrategy,
// which truncates the string from the middle and inserts the default ellipsis
// marker between the preserved head and tail.
//
// Returns:
//   - A Strategy that places an ellipsis in the middle.
//
// Example:
//
//	s := truncate.NewEllipsisMiddleStrategy()
//	result := s.Truncate("Hello, World!", 8) // "Hel…ld!"
func NewEllipsisMiddleStrategy() Strategy {
	return EllipsisMiddleStrategy{}
}
