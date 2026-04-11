package encoding

import "errors"

// Error messages for JSON operations.
//
//   - ErrNilInterface is returned when a nil interface is passed to a JSON function.
//   - ErrInvalidRawMessage is returned when an invalid json.RawMessage is passed to a JSON function.
//   - ErrNonFiniteFloat is returned when a non-finite float (NaN/Inf) is passed to a JSON function.
//   - ErrUnsupportedValue is returned when an unsupported value (e.g., non-nil func, chan, etc.) is passed to a JSON function.
//   - ErrMarshalPanicRecovered is returned when a panic occurs during JSON marshalling.
//   - ErrEmptyInput is returned when an empty or whitespace-only string is passed to a JSON function.
//   - ErrInvalidJSON is returned when a byte slice or string that does not constitute valid JSON is passed to a JSON function.
var (
	ErrNilInterface          = errors.New("nil interface input")
	ErrInvalidRawMessage     = errors.New("invalid json.RawMessage")
	ErrNonFiniteFloat        = errors.New("non-finite float (NaN/Inf)")
	ErrUnsupportedValue      = errors.New("unsupported value (e.g., non-nil func, chan, etc.)")
	ErrMarshalPanicRecovered = errors.New("json marshal panic recovered")
	ErrEmptyInput            = errors.New("empty input")
	ErrInvalidJSON           = errors.New("invalid JSON")
)
