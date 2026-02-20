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

// Infer will perform conversion by inferring the conversion operation from
// the base type of a pointer to a supported type. The value is assigned
// directly so only an error is returned.
//
// Example:
//
//	var into int64
//	err := conv.Infer(&into, "42")
//	// into -> 42
//
// Parameters:
//   - into: A pointer to the variable where the converted value will be stored.
//   - from: The source value to convert.
//
// Returns:
//   - An error if the conversion fails.
func Infer(into, from any) error {
	return defaultConverter.Infer(into, from)
}

// Bool will convert the given value to a bool, returns the default value of
// false if a conversion cannot be made.
//
// Supported string values for true:  "1", "t", "T", "true", "True", "TRUE", "y", "Y", "yes", "Yes", "YES"
//
// Supported string values for false: "0", "f", "F", "false", "False", "FALSE", "n", "N", "no", "No", "NO"
//
// Example:
//
//	val, err := conv.Bool("true")
//	// val -> true
//	val2, err := conv.Bool(0)
//	// val2 -> false
//
// Parameters:
//   - from: The source value to convert.
//
// Returns:
//   - The converted bool value.
func Bool(from any) (bool, error) {
	return defaultConverter.Bool(from)
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

// Duration will convert the given value to a time.Duration, returns the default
// value of 0 if a conversion cannot be made.
//
// Example:
//
//	val, err := conv.Duration("2h45m")
//	// val -> 2h45m0s
//
// Parameters:
//   - from: The source value to convert.
//
// Returns:
//   - The converted time.Duration value.
func Duration(from any) (time.Duration, error) {
	return defaultConverter.Duration(from)
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

// String will convert the given value to a string, returns the default value
// of "" if a conversion cannot be made.
//
// Example:
//
//	val, err := conv.String(12345)
//	// val -> "12345"
//
// Parameters:
//   - from: The source value to convert.
//
// Returns:
//   - The converted string value.
//   - An error if conversion fails.
func String(from any) (string, error) {
	return defaultConverter.String(from)
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

// Time will convert the given value to a time.Time, returns the empty struct
// time.Time{} if a conversion cannot be made.
//
// Example:
//
//	val, err := conv.Time("2024-01-02T15:04:05Z")
//	// val -> 2024-01-02 15:04:05 +0000 UTC
//
// Parameters:
//   - from: The source value to convert.
//
// Returns:
//   - The converted time.Time value.
//   - An error if conversion fails.
func Time(from any) (time.Time, error) {
	return defaultConverter.Time(from)
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

// Float32 will convert the given value to a float32, returns the default value
// of 0.0 if a conversion cannot be made.
//
// Example:
//
//	val, err := conv.Float32("3.14")
//	// val -> 3.14
//
// Parameters:
//   - from: The source value to convert.
//
// Returns:
//   - The converted float32 value.
//   - An error if conversion fails.
func Float32(from any) (float32, error) {
	return defaultConverter.Float32(from)
}

// Float64 will convert the given value to a float64, returns the default value
// of 0.0 if a conversion cannot be made.
//
// Example:
//
//	val, err := conv.Float64("2.71828")
//	// val -> 2.71828
//
// Parameters:
//   - from: The source value to convert.
//
// Returns:
//   - The converted float64 value.
//   - An error if conversion fails.
func Float64(from any) (float64, error) {
	return defaultConverter.Float64(from)
}

// Int will convert the given value to an int, returns the default value of 0 if
// a conversion cannot be made.
//
// Example:
//
//	val, err := conv.Int("123")
//	// val -> 123
//
// Parameters:
//   - from: The source value to convert.
//
// Returns:
//   - The converted int value.
//   - An error if conversion fails.
func Int(from any) (int, error) {
	return defaultConverter.Int(from)
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

// Int8 will convert the given value to an int8, returns the default value of 0
// if a conversion cannot be made.
//
// Example:
//
//	val, err := conv.Int8("127")
//	// val -> 127
//
// Parameters:
//   - from: The source value to convert.
//
// Returns:
//   - The converted int8 value.
//   - An error if conversion fails.
func Int8(from any) (int8, error) {
	return defaultConverter.Int8(from)
}

// Int16 will convert the given value to an int16, returns the default value of 0
// if a conversion cannot be made.
//
// Example:
//
//	val, err := conv.Int16("32000")
//	// val -> 32000
//
// Parameters:
//   - from: The source value to convert.
//
// Returns:
//   - The converted int16 value.
//   - An error if conversion fails.
func Int16(from any) (int16, error) {
	return defaultConverter.Int16(from)
}

// Int32 will convert the given value to an int32, returns the default value of 0
// if a conversion cannot be made.
//
// Example:
//
//	val, err := conv.Int32("2000000000")
//	// val -> 2000000000
//
// Parameters:
//   - from: The source value to convert.
//
// Returns:
//   - The converted int32 value.
//   - An error if conversion fails.
func Int32(from any) (int32, error) {
	return defaultConverter.Int32(from)
}

// Int64 will convert the given value to an int64, returns the default value of 0
// if a conversion cannot be made.
//
// Example:
//
//	val, err := conv.Int64("9000000000")
//	// val -> 9000000000
//
// Parameters:
//   - from: The source value to convert.
//
// Returns:
//   - The converted int64 value.
//   - An error if conversion fails.
func Int64(from any) (int64, error) {
	return defaultConverter.Int64(from)
}

// Uint will convert the given value to a uint, returns the default value of 0
// if a conversion cannot be made.
//
// Example:
//
//	val, err := conv.Uint("12345")
//	// val -> 12345
//
// Parameters:
//   - from: The source value to convert.
//
// Returns:
//   - The converted uint value.
//   - An error if conversion fails.
func Uint(from any) (uint, error) {
	return defaultConverter.Uint(from)
}

// Uint8 will convert the given value to a uint8, returns the default value of 0
// if a conversion cannot be made.
//
// Example:
//
//	val, err := conv.Uint8("255")
//	// val -> 255
//
// Parameters:
//   - from: The source value to convert.
//
// Returns:
//   - The converted uint8 value.
//   - An error if conversion fails.
func Uint8(from any) (uint8, error) {
	return defaultConverter.Uint8(from)
}

// Uint16 will convert the given value to a uint16, returns the default value of 0
// if a conversion cannot be made.
//
// Example:
//
//	val, err := conv.Uint16("65000")
//	// val -> 65000
//
// Parameters:
//   - from: The source value to convert.
//
// Returns:
//   - The converted uint16 value.
//   - An error if conversion fails.
func Uint16(from any) (uint16, error) {
	return defaultConverter.Uint16(from)
}

// Uint32 will convert the given value to a uint32, returns the default value of 0
// if a conversion cannot be made.
//
// Example:
//
//	val, err := conv.Uint32("4000000000")
//	// val -> 4000000000
//
// Parameters:
//   - from: The source value to convert.
//
// Returns:
//   - The converted uint32 value.
//   - An error if conversion fails.
func Uint32(from any) (uint32, error) {
	return defaultConverter.Uint32(from)
}

// Uint64 will convert the given value to a uint64, returns the default value of 0
// if a conversion cannot be made.
//
// Example:
//
//	val, err := conv.Uint64("9000000000")
//	// val -> 9000000000
//
// Parameters:
//   - from: The source value to convert.
//
// Returns:
//   - The converted uint64 value.
//   - An error if conversion fails.
func Uint64(from any) (uint64, error) {
	return defaultConverter.Uint64(from)
}
