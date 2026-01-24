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
// Section: Pointer indirection
// ///////////////////////////

// indirectValue dereferences pointers until a non-pointer value is found.
// It handles recursive pointer types safely to prevent infinite loops.
//
// Parameters:
//   - val: The reflect.Value to indirect.
//
// Returns:
//   - The dereferenced reflect.Value.
func indirectValue(val reflect.Value) reflect.Value {
	var last uintptr
	for {
		if val.Kind() != reflect.Ptr {
			return val
		}

		// Check for nil pointer
		if val.IsNil() {
			return reflect.Value{}
		}

		// Check for circular reference
		ptr := val.Pointer()
		if ptr == last {
			return val
		}
		last, val = ptr, val.Elem()
	}
}

// indirect dereferences pointers until a non-pointer value is found.
//
// Returns the original value if it cannot be dereferenced.
//
// Parameters:
//   - value: The interface value to indirect.
//
// Returns:
//   - The dereferenced interface value.
func indirect(value any) any {
	for {
		val := reflect.ValueOf(value)
		if val.Kind() != reflect.Ptr {
			return value
		}

		res := reflect.Indirect(val)
		if !res.IsValid() || !res.CanInterface() {
			return value
		}

		// Check for circular reference
		if res.Kind() == reflect.Ptr && val.Pointer() == res.Pointer() {
			return value
		}

		value = res.Interface()
	}
}
