package conv

import (
	"reflect"

	"github.com/sivaosorg/replify/pkg/encoding"
)

// jsonpass converts a Go value to its JSON string representation or returns the value directly if it is already a string.
//
// This function checks if the input data is a string; if so, it returns it directly.
// Otherwise, it marshals the input value `data` into a JSON string using the
// MarshalToString function. If an error occurs during marshalling, it returns an empty string.
//
// Parameters:
//   - `data`: The Go value to be converted to JSON, or a string to be returned directly.
//
// Returns:
//   - A string containing the JSON representation of the input value, or an empty string if an error occurs.
//
// Example:
//
//	jsonStr := jsonpass(myStruct)
func Jsonpass(data any) string {
	return encoding.JsonSafe(data)
}

// ///////////////////////////
// Section:  Slice conversion
// ///////////////////////////

// ToSlice converts the given value to a slice of the specified element type.
// If the input is not a slice/array, it wraps the value in a single-element slice.
//
// Parameters:
//   - from: The value to convert.
//
// Returns:
//   - A slice of interface{} containing the converted elements.
//   - An error if conversion fails.
func (c *Converter) ToSlice(from any) ([]any, error) {
	if from == nil {
		if c.nilAsZero {
			return nil, nil
		}
		return nil, newConvError(from, "[]any")
	}

	// Fast path for []any
	if v, ok := from.([]any); ok {
		return v, nil
	}

	value := reflect.ValueOf(from)
	kind := value.Kind()

	// Handle pointer
	if kind == reflect.Ptr {
		if value.IsNil() {
			return nil, nil
		}
		value = value.Elem()
		kind = value.Kind()
	}

	// If not a slice or array, wrap in single-element slice
	if kind != reflect.Slice && kind != reflect.Array {
		return []any{from}, nil
	}

	// Convert slice elements
	result := make([]any, value.Len())
	for i := 0; i < value.Len(); i++ {
		result[i] = value.Index(i).Interface()
	}

	return result, nil
}

// ToIntSlice converts the given value to a slice of int.
//
// Parameters:
//   - from:  The value to convert (slice, array, or single value).
//
// Returns:
//   - A slice of int.
//   - An error if conversion fails.
func (c *Converter) ToIntSlice(from any) ([]int, error) {
	if from == nil {
		if c.nilAsZero {
			return nil, nil
		}
		return nil, newConvError(from, "[]int")
	}

	// Fast path for []int
	if v, ok := from.([]int); ok {
		return v, nil
	}

	slice, err := c.ToSlice(from)
	if err != nil {
		return nil, err
	}

	result := make([]int, len(slice))
	for i, v := range slice {
		converted, err := c.Int(v)
		if err != nil {
			return nil, newConvErrorf("cannot convert element %d to int:  %v", i, err)
		}
		result[i] = converted
	}

	return result, nil
}

// ToInt64Slice converts the given value to a slice of int64.
//
// Parameters:
//   - from: The value to convert (slice, array, or single value).
//
// Returns:
//   - A slice of int64.
//   - An error if conversion fails.
func (c *Converter) ToInt64Slice(from any) ([]int64, error) {
	if from == nil {
		if c.nilAsZero {
			return nil, nil
		}
		return nil, newConvError(from, "[]int64")
	}

	// Fast path for []int64
	if v, ok := from.([]int64); ok {
		return v, nil
	}

	slice, err := c.ToSlice(from)
	if err != nil {
		return nil, err
	}

	result := make([]int64, len(slice))
	for i, v := range slice {
		converted, err := c.Int64(v)
		if err != nil {
			return nil, newConvErrorf("cannot convert element %d to int64: %v", i, err)
		}
		result[i] = converted
	}

	return result, nil
}

// ToFloat64Slice converts the given value to a slice of float64.
//
// Parameters:
//   - from: The value to convert (slice, array, or single value).
//
// Returns:
//   - A slice of float64.
//   - An error if conversion fails.
func (c *Converter) ToFloat64Slice(from any) ([]float64, error) {
	if from == nil {
		if c.nilAsZero {
			return nil, nil
		}
		return nil, newConvError(from, "[]float64")
	}

	// Fast path for []float64
	if v, ok := from.([]float64); ok {
		return v, nil
	}

	slice, err := c.ToSlice(from)
	if err != nil {
		return nil, err
	}

	result := make([]float64, len(slice))
	for i, v := range slice {
		converted, err := c.Float64(v)
		if err != nil {
			return nil, newConvErrorf("cannot convert element %d to float64: %v", i, err)
		}
		result[i] = converted
	}

	return result, nil
}

// ToStringSlice converts the given value to a slice of string.
//
// Parameters:
//   - from:  The value to convert (slice, array, or single value).
//
// Returns:
//   - A slice of string.
//   - An error if conversion fails.
func (c *Converter) ToStringSlice(from any) ([]string, error) {
	if from == nil {
		if c.nilAsZero {
			return nil, nil
		}
		return nil, newConvError(from, "[]string")
	}

	// Fast path for []string
	if v, ok := from.([]string); ok {
		return v, nil
	}

	slice, err := c.ToSlice(from)
	if err != nil {
		return nil, err
	}

	result := make([]string, len(slice))
	for i, v := range slice {
		converted, err := c.String(v)
		if err != nil {
			return nil, newConvErrorf("cannot convert element %d to string: %v", i, err)
		}
		result[i] = converted
	}

	return result, nil
}

// ToBoolSlice converts the given value to a slice of bool.
//
// Parameters:
//   - from:  The value to convert (slice, array, or single value).
//
// Returns:
//   - A slice of bool.
//   - An error if conversion fails.
func (c *Converter) ToBoolSlice(from any) ([]bool, error) {
	if from == nil {
		if c.nilAsZero {
			return nil, nil
		}
		return nil, newConvError(from, "[]bool")
	}

	// Fast path for []bool
	if v, ok := from.([]bool); ok {
		return v, nil
	}

	slice, err := c.ToSlice(from)
	if err != nil {
		return nil, err
	}

	result := make([]bool, len(slice))
	for i, v := range slice {
		converted, err := c.Bool(v)
		if err != nil {
			return nil, newConvErrorf("cannot convert element %d to bool: %v", i, err)
		}
		result[i] = converted
	}

	return result, nil
}

// ///////////////////////////
// Section: Package-level Slice functions
// ///////////////////////////

// IntSlice converts the given value to a slice of int.
//
// Parameters:
//   - from:  The value to convert (slice, array, or single value).
//
// Returns:
//   - A slice of int.
//   - An error if conversion fails.
func IntSlice(from any) ([]int, error) {
	return defaultConverter.ToIntSlice(from)
}

// Int64Slice converts the given value to a slice of int64.
//
// Parameters:
//   - from: The value to convert (slice, array, or single value).
//
// Returns:
//   - A slice of int64.
//   - An error if conversion fails.
func Int64Slice(from any) ([]int64, error) {
	return defaultConverter.ToInt64Slice(from)
}

// Float64Slice converts the given value to a slice of float64.
//
// Parameters:
//   - from: The value to convert (slice, array, or single value).
//
// Returns:
//   - A slice of float64.
//   - An error if conversion fails.
func Float64Slice(from any) ([]float64, error) {
	return defaultConverter.ToFloat64Slice(from)
}

// BoolSlice converts the given value to a slice of bool.
//
// Parameters:
//   - from:  The value to convert (slice, array, or single value).
//
// Returns:
//   - A slice of bool.
//   - An error if conversion fails.
func BoolSlice(from any) ([]bool, error) {
	return defaultConverter.ToBoolSlice(from)
}

// ///////////////////////////
// Section: Generic Slice functions
// ///////////////////////////

// Slice converts the given value to a slice of type T.
//
// Parameters:
//   - from: The value to convert (slice, array, or single value).
//
// Returns:
//   - A slice of type T.
//   - An error if conversion fails.
//
// Example:
//
//	val, err := conv.Slice[int]([]any{"1", 2, 3.0})
//	//	val -> []int{1, 2, 3}
func Slice[T any](from any) ([]T, error) {
	if from == nil {
		return nil, nil
	}

	// Fast path for []T
	if v, ok := from.([]T); ok {
		return v, nil
	}

	slice, err := defaultConverter.ToSlice(from)
	if err != nil {
		return nil, err
	}

	result := make([]T, len(slice))
	for i, v := range slice {
		converted, err := To[T](v)
		if err != nil {
			return nil, newConvErrorf("cannot convert element %d:  %v", i, err)
		}
		result[i] = converted
	}

	return result, nil
}

// SliceOrDefault converts to slice of T, returning default on failure.
//
// Parameters:
//   - from: The value to convert (slice, array, or single value).
//   - defaultValue: The default slice to return if conversion fails.
//
// Returns:
//   - A slice of type T or the default slice if conversion fails.
//
// Example:
//
//	val := conv.SliceOrDefault[int]("not a slice", []int{10, 20, 30})
//	//	val -> []int{10, 20, 30}
func SliceOrDefault[T any](from any, defaultValue []T) []T {
	if v, err := Slice[T](from); err == nil {
		return v
	}
	return defaultValue
}
