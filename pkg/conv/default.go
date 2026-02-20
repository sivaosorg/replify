package conv

import "time"

// ToOrDefault converts the given value to type T, returning defaultValue if conversion fails.
//
// Example:
//
//	val := conv.ToOrDefault[int]("invalid", 100)
//	// val -> 100
//
// Parameters:
//   - from: The source value to convert.
//   - defaultValue: The value to return if conversion fails.
//
// Returns:
//   - The converted value of type T, or defaultValue if conversion fails.
func ToOrDefault[T any](from any, defaultValue T) T {
	if v, err := To[T](from); err == nil {
		return v
	}
	return defaultValue
}

// InferOrDefault attempts to infer conversion, using default value on failure.
//
// Example:
//
//	val := conv.InferOrDefault[int]("not an int", 100)
//	// val -> 100
//	val2 := conv.InferOrDefault[int]("42", 100)
//	// val2 -> 42
//
// Parameters:
//   - from: The source value to convert.
//   - defaultValue: The default value to return if conversion fails.
//
// Returns:
//   - The converted value or the default value if conversion fails.
func InferOrDefault[T any](from any, defaultValue T) T {
	var into T
	if err := defaultConverter.Infer(&into, from); err != nil {
		return defaultValue
	}
	return into
}

// BoolOrDefault returns the converted bool value or the provided default if conversion fails.
//
// Example:
//
//	val := conv.BoolOrDefault("invalid", true)
//	// val -> true
//
// Parameters:
//   - from: The source value to convert.
//   - defaultValue: The value to return if conversion fails.
//
// Returns:
//   - The converted bool value, or defaultValue if conversion fails.
func BoolOrDefault(from any, defaultValue bool) bool {
	if v, err := defaultConverter.Bool(from); err == nil {
		return v
	}
	return defaultValue
}

// DurationOrDefault returns the converted duration or the provided default if conversion fails.
//
// Example:
//
//	val := conv.DurationOrDefault("invalid", 30*time.Minute)
//	// val -> 30m0s
//
// Parameters:
//   - from: The source value to convert.
//   - defaultValue: The value to return if conversion fails.
//
// Returns:
//   - The converted time.Duration value, or defaultValue if conversion fails.
func DurationOrDefault(from any, defaultValue time.Duration) time.Duration {
	if v, err := defaultConverter.Duration(from); err == nil {
		return v
	}
	return defaultValue
}

// StringOrDefault returns the converted string or the provided default if conversion fails.
//
// Example:
//
//	val := conv.StringOrDefault(3.14159, "default")
//	// val -> "3.14159"
//
// Parameters:
//   - from: The source value to convert.
//   - defaultValue: The value to return if conversion fails.
//
// Returns:
//   - The converted string value, or defaultValue if conversion fails.
func StringOrDefault(from any, defaultValue string) string {
	if v, err := defaultConverter.String(from); err == nil {
		return v
	}
	return defaultValue
}

// TimeOrDefault returns the converted time or the provided default if conversion fails.
//
// Example:
//
//	val := conv.TimeOrDefault("invalid", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
//	// val -> 2024-01-01 00:00:00 +0000 UTC
//
// Parameters:
//   - from: The source value to convert.
//   - defaultValue: The value to return if conversion fails.
//
// Returns:
//   - The converted time.Time value, or defaultValue if conversion fails.
func TimeOrDefault(from any, defaultValue time.Time) time.Time {
	if v, err := defaultConverter.Time(from); err == nil {
		return v
	}
	return defaultValue
}

// IntOrDefault returns the converted int or the provided default if conversion fails.
//
// Example:
//
//	val := conv.IntOrDefault("invalid", 42)
//	// val -> 42
//
// Parameters:
//   - from: The source value to convert.
//   - defaultValue: The value to return if conversion fails.
//
// Returns:
//   - The converted int value, or defaultValue if conversion fails.
func IntOrDefault(from any, defaultValue int) int {
	if v, err := defaultConverter.Int(from); err == nil {
		return v
	}
	return defaultValue
}

// Int8OrDefault returns the converted int8 or the provided default if conversion fails.
//
// Example:
//
//	val := conv.Int8OrDefault("invalid", 42)
//	// val -> 42
//
// Parameters:
//   - from: The source value to convert.
//   - defaultValue: The value to return if conversion fails.
//
// Returns:
//   - The converted int8 value, or defaultValue if conversion fails.
func Int8OrDefault(from any, defaultValue int8) int8 {
	if v, err := defaultConverter.Int8(from); err == nil {
		return v
	}
	return defaultValue
}

// Int16OrDefault returns the converted int16 or the provided default if conversion fails.
//
// Example:
//
//	val := conv.Int16OrDefault("invalid", 42)
//	// val -> 42
//
// Parameters:
//   - from: The source value to convert.
//   - defaultValue: The value to return if conversion fails.
//
// Returns:
//   - The converted int16 value, or defaultValue if conversion fails.
func Int16OrDefault(from any, defaultValue int16) int16 {
	if v, err := defaultConverter.Int16(from); err == nil {
		return v
	}
	return defaultValue
}

// Int32OrDefault returns the converted int32 or the provided default if conversion fails.
//
// Example:
//
//	val := conv.Int32OrDefault("invalid", 42)
//	// val -> 42
//
// Parameters:
//   - from: The source value to convert.
//   - defaultValue: The value to return if conversion fails.
//
// Returns:
//   - The converted int32 value, or defaultValue if conversion fails.
func Int32OrDefault(from any, defaultValue int32) int32 {
	if v, err := defaultConverter.Int32(from); err == nil {
		return v
	}
	return defaultValue
}

// Int64OrDefault returns the converted int64 or the provided default if conversion fails.
//
// Example:
//
//	val := conv.Int64OrDefault("invalid", 1234567890)
//	// val -> 1234567890
//
// Parameters:
//   - from: The source value to convert.
//   - defaultValue: The value to return if conversion fails.
//
// Returns:
//   - The converted int64 value, or defaultValue if conversion fails.
func Int64OrDefault(from any, defaultValue int64) int64 {
	if v, err := defaultConverter.Int64(from); err == nil {
		return v
	}
	return defaultValue
}

// UintOrDefault returns the converted uint or the provided default if conversion fails.
//
// Parameters:
//   - from: The value to convert.
//   - defaultValue: The default value to return if conversion fails.
//
// Returns:
//   - The converted uint value or the default value.
func UintOrDefault(from any, defaultValue uint) uint {
	if v, err := defaultConverter.Uint(from); err == nil {
		return v
	}
	return defaultValue
}

// Uint8OrDefault returns the converted uint8 or the provided default if conversion fails.
//
// Parameters:
//   - from: The value to convert.
//   - defaultValue: The default value to return if conversion fails.
//
// Returns:
//   - The converted uint8 value or the default value.
func Uint8OrDefault(from any, defaultValue uint8) uint8 {
	if v, err := defaultConverter.Uint8(from); err == nil {
		return v
	}
	return defaultValue
}

// Uint16OrDefault returns the converted uint16 or the provided default if conversion fails.
//
// Parameters:
//   - from: The value to convert.
//   - defaultValue: The default value to return if conversion fails.
//
// Returns:
//   - The converted uint16 value or the default value.
func Uint16OrDefault(from any, defaultValue uint16) uint16 {
	if v, err := defaultConverter.Uint16(from); err == nil {
		return v
	}
	return defaultValue
}

// Uint32OrDefault returns the converted uint32 or the provided default if conversion fails.
//
// Parameters:
//   - from: The value to convert.
//   - defaultValue: The default value to return if conversion fails.
//
// Returns:
//   - The converted uint32 value or the default value.
func Uint32OrDefault(from any, defaultValue uint32) uint32 {
	if v, err := defaultConverter.Uint32(from); err == nil {
		return v
	}
	return defaultValue
}

// Uint64OrDefault returns the converted uint64 or the provided default if conversion fails.
//
// Example:
//
//	val := conv.Uint64OrDefault("invalid", 9876543210)
//	// val -> 9876543210
//
// Parameters:
//   - from: The source value to convert.
//   - defaultValue: The value to return if conversion fails.
//
// Returns:
//   - The converted uint64 value, or defaultValue if conversion fails.
func Uint64OrDefault(from any, defaultValue uint64) uint64 {
	if v, err := defaultConverter.Uint64(from); err == nil {
		return v
	}
	return defaultValue
}

// Float32OrDefault returns the converted float32 or the provided default if conversion fails.
//
// Parameters:
//   - from: The value to convert.
//   - defaultValue: The default value to return if conversion fails.
//
// Returns:
//   - The converted float32 value or the default value.
func Float32OrDefault(from any, defaultValue float32) float32 {
	if v, err := defaultConverter.Float32(from); err == nil {
		return v
	}
	return defaultValue
}

// Float64OrDefault returns the converted float64 or the provided default if conversion fails.
//
// Example:
//
//	val := conv.Float64OrDefault("invalid", 1.61803)
//	// val -> 1.61803
//
// Parameters:
//   - from: The source value to convert.
//   - defaultValue: The value to return if conversion fails.
//
// Returns:
//   - The converted float64 value, or defaultValue if conversion fails.
func Float64OrDefault(from any, defaultValue float64) float64 {
	if v, err := defaultConverter.Float64(from); err == nil {
		return v
	}
	return defaultValue
}
