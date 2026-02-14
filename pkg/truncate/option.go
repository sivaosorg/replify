package truncate

// NewTruncator returns a new TruncatorBuilder initialised with sensible defaults:
//   - Omission: DefaultOmission ("…")
//   - Position: PositionEnd (truncate from the end)
//   - MaxLength: 0 (must be set explicitly before use)
//
// Returns:
//   - A pointer to a TruncatorBuilder ready for configuration.
//
// Example:
//
//	t := truncate.NewTruncator().
//	    WithMaxLength(20).
//	    Build()
func NewTruncator() *TruncatorBuilder {
	return &TruncatorBuilder{
		omission:  DefaultOmission,
		position:  PositionEnd,
		maxLength: 0,
	}
}

// WithOmission sets the omission marker that will be inserted where characters
// are removed during truncation. Common choices are "…" (single character) and
// "..." (three characters).
//
// Parameters:
//   - omission: The omission string to use.
//
// Returns:
//   - A pointer to the TruncatorBuilder for method chaining.
//
// Example:
//
//	builder := truncate.NewTruncator().WithOmission("...")
func (b *TruncatorBuilder) WithOmission(omission string) *TruncatorBuilder {
	b.omission = omission
	return b
}

// WithPosition sets the position at which the omission marker will be placed.
// Use one of PositionEnd, PositionStart, or PositionMiddle.
//
// Parameters:
//   - position: The truncation position.
//
// Returns:
//   - A pointer to the TruncatorBuilder for method chaining.
//
// Example:
//
//	builder := truncate.NewTruncator().WithPosition(truncate.PositionStart)
func (b *TruncatorBuilder) WithPosition(position Position) *TruncatorBuilder {
	b.position = position
	return b
}

// WithMaxLength sets the maximum number of runes the truncated string may contain.
// This includes the runes consumed by the omission marker itself.
//
// Parameters:
//   - maxLength: The maximum allowed rune count.
//
// Returns:
//   - A pointer to the TruncatorBuilder for method chaining.
//
// Example:
//
//	builder := truncate.NewTruncator().WithMaxLength(30)
func (b *TruncatorBuilder) WithMaxLength(maxLength int) *TruncatorBuilder {
	b.maxLength = maxLength
	return b
}

// Build creates and returns an immutable Truncator configured with the values
// accumulated by the builder. The returned Truncator is safe for concurrent use
// because all of its fields are set at construction time.
//
// Returns:
//   - A pointer to the newly created Truncator.
//
// Example:
//
//	t := truncate.NewTruncator().
//	    WithOmission("...").
//	    WithPosition(truncate.PositionMiddle).
//	    WithMaxLength(20).
//	    Build()
func (b *TruncatorBuilder) Build() *Truncator {
	return &Truncator{
		omission:  b.omission,
		position:  b.position,
		maxLength: b.maxLength,
	}
}
