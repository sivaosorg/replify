package conv

import "time"

// ///////////////////////////
// Section: Default converter instance
// ///////////////////////////

// defaultConverter is the package-level converter instance used by
// all top-level conversion functions.  It is safe for concurrent use.
var defaultConverter = NewConverter()

// ///////////////////////////
// Section: Generic functions (Go 1.18+)
// ///////////////////////////

// To converts the given value to the specified type T.
// Returns the zero value of T and an error if conversion fails.
//
// Example:
//
//	val, err := conv.To[int]("42")
//	// val -> 42
//
// Parameters:
//   - from: The source value to convert.
//
// Returns:
//   - The converted value of type T.
//   - An error if conversion fails.
func To[T any](from any) (T, error) {
	var zero T
	var result any
	var err error

	switch any(zero).(type) {
	case bool:
		result, err = defaultConverter.Bool(from)
	case string:
		result, err = defaultConverter.String(from)
	case int:
		result, err = defaultConverter.Int(from)
	case int8:
		result, err = defaultConverter.Int8(from)
	case int16:
		result, err = defaultConverter.Int16(from)
	case int32:
		result, err = defaultConverter.Int32(from)
	case int64:
		result, err = defaultConverter.Int64(from)
	case uint:
		result, err = defaultConverter.Uint(from)
	case uint8:
		result, err = defaultConverter.Uint8(from)
	case uint16:
		result, err = defaultConverter.Uint16(from)
	case uint32:
		result, err = defaultConverter.Uint32(from)
	case uint64:
		result, err = defaultConverter.Uint64(from)
	case float32:
		result, err = defaultConverter.Float32(from)
	case float64:
		result, err = defaultConverter.Float64(from)
	case time.Time:
		result, err = defaultConverter.Time(from)
	case time.Duration:
		result, err = defaultConverter.Duration(from)
	default:
		return zero, newConvError(from, "unsupported type")
	}

	if err != nil {
		return zero, err
	}
	return result.(T), nil
}
