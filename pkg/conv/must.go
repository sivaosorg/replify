package conv

import "time"

// MustTo converts the given value to type T, panicking if conversion fails.
//
// Example:
//
//	val := conv.MustTo[int]("42")
//	// val -> 42
//
// Parameters:
//   - from: The source value to convert.
//
// Returns:
//   - The converted value of type T.
func MustTo[T any](from any) T {
	v, err := To[T](from)
	if err != nil {
		panic(err)
	}
	return v
}

// MustBool returns the converted bool value or panics if conversion fails.
//
// Example:
//
//	val := conv.MustBool("true")
//	// val -> true
//
// Parameters:
//   - from: The source value to convert.
//
// Returns:
//   - The converted bool value.
func MustBool(from any) bool {
	v, err := defaultConverter.Bool(from)
	if err != nil {
		panic(err)
	}
	return v
}

// MustDuration returns the converted duration or panics if conversion fails.
//
// Example:
//
//	val := conv.MustDuration("1h15m")
//	// val -> 1h15m0s
//
// Parameters:
//   - from: The source value to convert.
//
// Returns:
//   - The converted time.Duration value.
func MustDuration(from any) time.Duration {
	v, err := defaultConverter.Duration(from)
	if err != nil {
		panic(err)
	}
	return v
}

// MustString returns the converted string or panics if conversion fails.
//
// Example:
//
//	val := conv.MustString(1001)
//	// val -> "1001"
//
// Parameters:
//   - from: The source value to convert.
//
// Returns:
//   - The converted string value.
func MustString(from any) string {
	v, err := defaultConverter.String(from)
	if err != nil {
		panic(err)
	}
	return v
}

// MustTime returns the converted time or panics if conversion fails.
//
// Example:
//
//	val := conv.MustTime("2024-12-25T10:00:00Z")
//	// val -> 2024-12-25 10:00:00 +0000 UTC
//
// Parameters:
//   - from: The source value to convert.
//
// Returns:
//   - The converted time.Time value.
func MustTime(from any) time.Time {
	v, err := defaultConverter.Time(from)
	if err != nil {
		panic(err)
	}
	return v
}

// MustInt returns the converted int or panics if conversion fails.
//
// Example:
//
//	val := conv.MustInt("456")
//	// val -> 456
//
// Parameters:
//   - from: The source value to convert.
//
// Returns:
//   - The converted int value.
func MustInt(from any) int {
	v, err := defaultConverter.Int(from)
	if err != nil {
		panic(err)
	}
	return v
}
