package conv

import (
	"fmt"

	"github.com/sivaosorg/replify/pkg/strutil"
)

// ///////////////////////////
// Section:  Conversion errors
// ///////////////////////////

// ConvError represents a type-conversion failure and its diagnostic context.
// It includes the source value, the target type name, and an optional custom message.
type ConvError struct {
	From    any    // The source value
	To      string // The target type name
	Message string // Additional error message
}

// Error implements the error interface for ConvError.
//
// It returns the custom message if provided; otherwise, it constructs a default error message
// indicating the source value and target type.
//
// Returns:
//   - A string representing the error message.
func (e *ConvError) Error() string {
	if strutil.IsNotEmpty(e.Message) {
		return e.Message
	}
	return fmt.Sprintf("cannot convert %#v (type %[1]T) to %v", e.From, e.To)
}

// ///////////////////////////
// Section: Error checking
// ///////////////////////////

// IsConvError checks if the given error is of type ConvError.
//
// Parameters:
//   - `err`: The error to be checked.
//
// Returns:
//   - A boolean indicating whether the error is a ConvError.
func IsConvError(err error) bool {
	_, ok := err.(*ConvError)
	return ok
}

// AsConvError attempts to cast the given error to a ConvError.
//
// Parameters:
//   - `err`: The error to be cast.
//
// Returns:
//   - A pointer to the ConvError if the cast is successful.
func AsConvError(err error) (*ConvError, bool) {
	if e, ok := err.(*ConvError); ok {
		return e, true
	}
	return nil, false
}

// ///////////////////////////
// Section: Error constructors
// ///////////////////////////

// newConvError creates a new ConvError instance with the given source value and target type name.
//
// Parameters:
//   - `from`: The source value that failed to convert.
//   - `to`: The target type name to which the conversion was attempted.
//
// Returns:
//   - A pointer to a ConvError instance representing the conversion error.
func newConvError(from any, to string) error {
	return &ConvError{From: from, To: to}
}

// newConvErrorf creates a new ConvError instance with a formatted error message.
//
// Parameters:
//   - `format`: The format string for the error message.
//   - `args`: The arguments to be formatted into the message.
//
// Returns:
//   - A pointer to a ConvError instance representing the conversion error with the formatted message.
func newConvErrorf(format string, args ...any) error {
	return &ConvError{Message: fmt.Sprintf(format, args...)}
}

// newConvErrorMsg creates a new ConvError instance with a custom error message.
//
// Parameters:
//   - `msg`: The custom error message.
//
// Returns:
//   - A pointer to a ConvError instance representing the conversion error with the custom message.
func newConvErrorMsg(msg string) error {
	return &ConvError{Message: msg}
}
