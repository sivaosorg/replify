package conv

import (
	"reflect"
)

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
