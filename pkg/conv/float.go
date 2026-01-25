package conv

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/sivaosorg/replify/pkg/encoding"
	"github.com/sivaosorg/replify/pkg/strutil"
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
// Section: Float conversion interface
// ///////////////////////////

// floatConverter is an interface for types that can convert themselves to float64.
type floatConverter interface {
	Float64() (float64, error)
}

// ///////////////////////////
// Section: Float64 conversion
// ///////////////////////////

// Float64 attempts to convert the given value to float64, returns the zero value
// and an error on failure.
//
// Parameters:
//   - from: The value to convert.
//
// Returns:
//   - The converted float64 value.
//   - An error if conversion fails.
func (c *Converter) Float64(from any) (float64, error) {
	if from == nil {
		if c.nilAsZero {
			return 0, nil
		}
		return 0, newConvError(from, "float64")
	}

	// Fast path for common types
	switch v := from.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case string:
		return c.stringToFloat64(v)
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case *float64:
		if v == nil {
			return 0, nil
		}
		return *v, nil
	case *float32:
		if v == nil {
			return 0, nil
		}
		return float64(*v), nil
	case *string:
		if v == nil {
			return 0, nil
		}
		return c.stringToFloat64(*v)
	}

	// Check for custom converter interface
	if conv, ok := from.(floatConverter); ok {
		return conv.Float64()
	}

	// Use reflection for other types
	return c.float64FromReflect(from)
}

// stringToFloat64 attempts to convert a string to float64.
//
// Parameters:
//   - v: The string to convert.
//
// Returns:
//   - The converted float64 value.
//   - An error if conversion fails.
func (c *Converter) stringToFloat64(v string) (float64, error) {
	if strutil.IsEmpty(v) {
		if c.emptyAsZero {
			return 0, nil
		}
		return 0, newConvErrorMsg("cannot convert empty string to float64")
	}

	if c.trimStrings {
		v = strings.TrimSpace(v)
	}

	// Try parsing as float
	if parsed, err := strconv.ParseFloat(v, 64); err == nil {
		return parsed, nil
	}

	// Try parsing as bool
	if parsed, err := c.stringToBool(v); err == nil {
		if parsed {
			return 1, nil
		}
		return 0, nil
	}

	return 0, newConvErrorf("cannot convert %q to float64", v)
}

// float64FromReflect converts a reflect.Value to float64.
//
// Parameters:
//   - from: The input value to convert.
//
// Returns:
//   - The converted float64 value.
//   - An error if the conversion fails.
func (c *Converter) float64FromReflect(from any) (float64, error) {
	value := indirectValue(reflect.ValueOf(from))
	if !value.IsValid() {
		if c.nilAsZero {
			return 0, nil
		}
		return 0, newConvError(from, "float64")
	}

	kind := value.Kind()
	switch {
	case kind == reflect.String:
		return c.stringToFloat64(value.String())
	case isKindFloat(kind):
		return value.Float(), nil
	case isKindInt(kind):
		return float64(value.Int()), nil
	case isKindUint(kind):
		return float64(value.Uint()), nil
	case isKindComplex(kind):
		return real(value.Complex()), nil
	case kind == reflect.Bool:
		if value.Bool() {
			return 1, nil
		}
		return 0, nil
	case isKindLength(kind):
		return float64(value.Len()), nil
	}

	return 0, newConvError(from, "float64")
}
