package codegen

import "errors"

// Sentinel errors returned by the library.
// They can be compared directly using errors.Is or the == operator.
var (
	// ErrInvalidLength is returned when Length is less than 1.
	ErrInvalidLength = errors.New("codegen: length must be greater than 0")

	// ErrEmptyCharset is returned when Charset is an empty string.
	ErrEmptyCharset = errors.New("codegen: charset must not be empty")

	// ErrInvalidCount is returned when the n argument passed to GenerateN is less than 1.
	ErrInvalidCount = errors.New("codegen: count must be greater than 0")
)
