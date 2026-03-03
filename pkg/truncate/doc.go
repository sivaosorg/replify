// Package truncate provides Unicode-aware string truncation with
// configurable omission markers and positioning.
//
// All operations work on rune counts rather than byte lengths, ensuring that
// multi-byte UTF-8 characters are never split and that length limits are
// applied consistently regardless of the script being used.
//
// # Strategies
//
// The package exposes four ready-made Strategy implementations:
//
//	CutStrategy{}               → plain cut, no omission marker
//	CutEllipsisStrategy{}       → trailing ellipsis ("Hello, …")
//	CutEllipsisLeadingStrategy{} → leading ellipsis ("…World!")
//	EllipsisMiddleStrategy{}    → middle ellipsis ("Hel…ld!")
//
// Strategies are applied through the Apply convenience function:
//
//	result := truncate.Apply("Hello, World!", 8, truncate.NewCutEllipsisStrategy())
//	// result: "Hello, …"
//
// # Builder API
//
// For more control, construct a Truncator with the fluent builder:
//
//	t := truncate.NewTruncator().
//	    WithOmission("...").
//	    WithPosition(truncate.PositionMiddle).
//	    WithMaxLength(15).
//	    Build()
//
//	result := t.Truncate("A very long string")
//	// result: "A very...string"
//
// TruncateWithLength overrides the configured maximum length on a per-call
// basis, making it easy to reuse a single Truncator for different output
// widths.
//
// # Omission Positions
//
// Three position constants control where the omission marker is inserted:
//
//	PositionEnd    (default) – marker appended after the preserved head
//	PositionStart            – marker prepended before the preserved tail
//	PositionMiddle           – marker inserted between preserved head and tail
//
// truncate is used internally by strutil to implement its truncation helpers
// and is available as a public API for callers that need precise control over
// truncation behaviour.
package truncate
