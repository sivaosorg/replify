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

// MustInfer performs type inference and conversion, panicking on failure.
//
// Example:
//
//	val := conv.MustInfer[int]("42")
//	// val -> 42
//
// Parameters:
//   - from: The source value to convert.
//
// Returns:
//   - The converted value.
func MustInfer[T any](from any) T {
	var into T
	if err := defaultConverter.Infer(&into, from); err != nil {
		panic(err)
	}
	return into
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

// MustUint returns the converted uint or panics if conversion fails.
//
// Parameters:
//   - from: The value to convert.
//
// Returns:
//   - The converted uint value.
func MustUint(from any) uint {
	v, err := defaultConverter.Uint(from)
	if err != nil {
		panic(err)
	}
	return v
}

// MustFloat32 returns the converted float32 or panics if conversion fails.
//
// Parameters:
//   - from: The value to convert.
//
// Returns:
//   - The converted float32 value.
func MustFloat32(from any) float32 {
	v, err := defaultConverter.Float32(from)
	if err != nil {
		panic(err)
	}
	return v
}

// MustFloat64 returns the converted float64 or panics if conversion fails.
//
// Parameters:
//   - from: The value to convert.
//
// Returns:
//   - The converted float64 value.
func MustFloat64(from any) float64 {
	v, err := defaultConverter.Float64(from)
	if err != nil {
		panic(err)
	}
	return v
}

// MustParseJSON is a package-level helper that parses a JSON string into a variable of type T.
//
// Parameters:
//   - `jsonStr`: The JSON string to be parsed.
//
// Returns:
//   - A variable of type T populated with the parsed data.
//   - Panics if the parsing fails.
func MustParseJSON[T any](jsonStr string) T {
	result, err := ParseJSON[T](jsonStr)
	if err != nil {
		panic(err)
	}
	return result
}

// MustParseJSONBytes is a package-level helper that parses a JSON byte slice into a variable of type T.
//
// Parameters:
//   - `data`: The JSON byte slice to be parsed.
//
// Returns:
//   - A variable of type T populated with the parsed data.
//   - Panics if the parsing fails.
func MustParseJSONBytes[T any](data []byte) T {
	result, err := ParseJSONBytes[T](data)
	if err != nil {
		panic(err)
	}
	return result
}

// MustClone creates a deep copy of the given value using JSON serialization and deserialization,
// panicking on failure.
//
// Parameters:
//   - `v`: The value to be cloned.
//
// Returns:
//   - A deep copy of the input value.
//
// Example:
//
//	original := MyStruct{Field: "value"}
//	clone := conv.MustClone(original)
//	// use clone
func MustClone[T any](v T) T {
	result, err := Clone(v)
	if err != nil {
		panic(err)
	}
	return result
}
