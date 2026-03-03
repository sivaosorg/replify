// Package conv provides flexible, panic-free conversion between Go's core
// scalar types and a handful of stdlib time types.
//
// The package is built around the Converter type, which encapsulates
// conversion logic and is safe for concurrent use. A shared package-level
// Converter backs a set of top-level convenience functions so that callers
// need not manage a Converter instance directly:
//
//	n, err := conv.Int("42")          // string  → int
//	f, err := conv.Float64("3.14")    // string  → float64
//	b, err := conv.Bool("yes")        // string  → bool
//	t, err := conv.Time("2024-01-02") // string  → time.Time
//
// # Generic Conversion
//
// To[T] is a generic convenience wrapper introduced in Go 1.18 that infers
// the conversion target from the declared type variable:
//
//	val, err := conv.To[int64]("9000000000")
//
// Infer accepts a pointer to the destination variable and assigns the
// converted value directly, making it ideal for configuration unmarshalling:
//
//	var timeout time.Duration
//	err := conv.Infer(&timeout, "30s")
//
// # Must Variants
//
// MustBool, MustInt, MustString, and other Must-prefixed functions mirror
// their error-returning counterparts but panic on failure. They are intended
// for use in initialisation code where a conversion failure is a programming
// error rather than a runtime condition.
//
// # Supported Types
//
// bool, string, int / int8 / int16 / int32 / int64,
// uint / uint8 / uint16 / uint32 / uint64,
// float32 / float64, time.Time, and time.Duration.
//
// Conversion between any pair of supported types is handled automatically,
// including numeric widening and narrowing, string-to-number parsing, and
// several flexible time-layout heuristics. When a conversion cannot be
// performed, the error value carries a descriptive message that identifies
// both the source value and the target type.
package conv
