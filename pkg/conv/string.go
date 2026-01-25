package conv

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

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
// Section: String conversion interface
// ///////////////////////////

// stringConverter is an interface for types that can convert themselves to string.
type stringConverter interface {
	String() string
}

// stringErrorConverter is an interface for types that can convert themselves to string with error.
type stringErrorConverter interface {
	String() (string, error)
}

// ///////////////////////////
// Section:  String conversion
// ///////////////////////////

// String returns the string representation from the given interface{} value.
// This function cannot fail for most types - it will use fmt.Sprintf as fallback.
//
// Parameters:
//   - from:  The value to convert.
//
// Returns:
//   - The converted string value.
//   - An error (currently always nil for compatibility, but check it for future-proofing).
func (c *Converter) String(from any) (string, error) {
	if from == nil {
		if c.nilAsZero {
			return "", nil
		}
		return "", newConvError(from, "string")
	}

	// Fast path for common types
	switch v := from.(type) {
	case string:
		return v, nil
	case *string:
		if v == nil {
			return "", nil
		}
		return *v, nil
	case []byte:
		return string(v), nil
	case *[]byte:
		if v == nil {
			return "", nil
		}
		return string(*v), nil
	case []rune:
		return string(v), nil
	case bool:
		return strconv.FormatBool(v), nil
	case *bool:
		if v == nil {
			return "", nil
		}
		return strconv.FormatBool(*v), nil
	case int:
		return strconv.Itoa(v), nil
	case *int:
		if v == nil {
			return "", nil
		}
		return strconv.Itoa(*v), nil
	case int8:
		return strconv.FormatInt(int64(v), 10), nil
	case int16:
		return strconv.FormatInt(int64(v), 10), nil
	case int32:
		return strconv.FormatInt(int64(v), 10), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case *int64:
		if v == nil {
			return "", nil
		}
		return strconv.FormatInt(*v, 10), nil
	case uint:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint64:
		return strconv.FormatUint(v, 10), nil
	case *uint64:
		if v == nil {
			return "", nil
		}
		return strconv.FormatUint(*v, 10), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case *float64:
		if v == nil {
			return "", nil
		}
		return strconv.FormatFloat(*v, 'f', -1, 64), nil
	case complex64:
		return fmt.Sprintf("%v", v), nil
	case complex128:
		return fmt.Sprintf("%v", v), nil
	case time.Time:
		return v.Format(time.RFC3339), nil
	case *time.Time:
		if v == nil {
			return "", nil
		}
		return v.Format(time.RFC3339), nil
	case time.Duration:
		return v.String(), nil
	case *time.Duration:
		if v == nil {
			return "", nil
		}
		return v.String(), nil
	case error:
		return v.Error(), nil
	case fmt.Stringer:
		return v.String(), nil
	}

	// Check for custom converter interfaces
	if conv, ok := from.(stringConverter); ok {
		return conv.String(), nil
	}
	if conv, ok := from.(stringErrorConverter); ok {
		return conv.String()
	}

	// Use reflection for other types
	return c.stringFromReflect(from)
}

// ///////////////////////////
// Section: String formatting
// ///////////////////////////

// StringSlice converts a slice of any type to a slice of strings.
//
// Parameters:
//   - from: The input value to be converted to a slice of strings.
//
// Returns:
//   - A slice of strings representing the converted values.
//   - An error if any conversion fails.
func StringSlice(from any) ([]string, error) {
	if from == nil {
		return nil, nil
	}

	// Fast path for string slice
	if v, ok := from.([]string); ok {
		return v, nil
	}

	value := reflect.ValueOf(from)
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		// Single value - wrap in slice
		s, err := defaultConverter.String(from)
		if err != nil {
			return nil, err
		}
		return []string{s}, nil
	}

	result := make([]string, value.Len())
	for i := 0; i < value.Len(); i++ {
		s, err := defaultConverter.String(value.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		result[i] = s
	}

	return result, nil
}

// ///////////////////////////
// Section: String utility functions
// ///////////////////////////

// StringOrEmpty returns the converted string or empty string if conversion fails.
//
// Parameters:
//   - from:  The value to convert.
//
// Returns:
//   - The converted string value, or empty string if conversion fails.
func StringOrEmpty(from any) string {
	v, _ := defaultConverter.String(from)
	return v
}

// Quote returns a double-quoted string safely escaped with Go syntax.
//
// Parameters:
//   - from:  The value to convert.
//
// Returns:
//   - The quoted string value.
func Quote(from any) string {
	s, _ := defaultConverter.String(from)
	return strconv.Quote(s)
}

// TrimSpace returns the string with all leading and trailing white space removed.
//
// Parameters:
//   - from:  The value to convert.
//
// Returns:
//   - The trimmed string value.
func TrimSpace(from any) string {
	s, _ := defaultConverter.String(from)
	return strings.TrimSpace(s)
}

// Join converts a slice of any type to a single string joined by the specified separator.
//
// Parameters:
//   - from: The input value to be converted to a joined string.
//   - sep:  The separator string to use between elements.
//
// Returns:
//   - A single string with all elements joined by the separator.
//   - An error if any conversion fails.
func Join(from any, sep string) (string, error) {
	slice, err := StringSlice(from)
	if err != nil {
		return "", err
	}
	return strings.Join(slice, sep), nil
}

// stringFromReflect converts a reflect.Value to a string based on its kind.
//
// Parameters:
//   - `from`: The input value to be converted to string.
//
// Returns:
//   - A string representation of the input value.
//   - An error if the conversion fails.
func (c *Converter) stringFromReflect(from any) (string, error) {
	value := indirectValue(reflect.ValueOf(from))
	if !value.IsValid() {
		if c.nilAsZero {
			return "", nil
		}
		return "", newConvError(from, "string")
	}

	kind := value.Kind()
	switch kind {
	case reflect.String:
		return value.String(), nil
	case reflect.Bool:
		return strconv.FormatBool(value.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10), nil
	case reflect.Float32:
		return strconv.FormatFloat(value.Float(), 'f', -1, 32), nil
	case reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'f', -1, 64), nil
	case reflect.Slice:
		if value.Type().Elem().Kind() == reflect.Uint8 {
			return string(value.Bytes()), nil
		}
	}

	// Fallback to fmt.Sprintf
	return fmt.Sprintf("%v", from), nil
}
