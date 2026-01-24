package conv

import (
	"reflect"
)

// ///////////////////////////
// Section: Kind checking functions
// ///////////////////////////

// isKindComplex returns true if the given Kind is a complex value.
//
// Parameters:
//   - k: The reflect.Kind to be checked.
//
// Returns:
//   - A boolean indicating whether the Kind is complex.
func isKindComplex(k reflect.Kind) bool {
	return k == reflect.Complex64 || k == reflect.Complex128
}

// isKindFloat returns true if the given Kind is a float value.
//
// Parameters:
//   - k: The reflect.Kind to be checked.
//
// Returns:
//   - A boolean indicating whether the Kind is a float.
func isKindFloat(k reflect.Kind) bool {
	return k == reflect.Float32 || k == reflect.Float64
}

// isKindInt returns true if the given Kind is a signed int value.
//
// Parameters:
//   - k: The reflect.Kind to be checked.
//
// Returns:
//   - A boolean indicating whether the Kind is a signed int.
func isKindInt(k reflect.Kind) bool {
	return k >= reflect.Int && k <= reflect.Int64
}

// isKindUint returns true if the given Kind is an unsigned int value.
//
// Parameters:
//   - k: The reflect.Kind to be checked.
//
// Returns:
//   - A boolean indicating whether the Kind is an unsigned int.
func isKindUint(k reflect.Kind) bool {
	return k >= reflect.Uint && k <= reflect.Uintptr
}

// isKindNumeric returns true if the given Kind is any numeric value.
//
// Parameters:
//   - k: The reflect.Kind to be checked.
//
// Returns:
//   - A boolean indicating whether the Kind is numeric.
func isKindNumeric(k reflect.Kind) bool {
	return (k >= reflect.Int && k <= reflect.Uint64) ||
		(k >= reflect.Float32 && k <= reflect.Complex128)
}

// isKindNillable returns true if the Kind can be nil.
//
// Parameters:
//   - k: The reflect.Kind to be checked.
//
// Returns:
//   - A boolean indicating whether the Kind is nillable.
func isKindNillable(k reflect.Kind) bool {
	return (k >= reflect.Chan && k <= reflect.Slice) || k == reflect.UnsafePointer
}

// isKindLength returns true if the Kind has a length.
//
// Parameters:
//   - k: The reflect.Kind to be checked.
//
// Returns:
//   - A boolean indicating whether the Kind has a length.
func isKindLength(k reflect.Kind) bool {
	return k == reflect.Array || k == reflect.Chan || k == reflect.Map ||
		k == reflect.Slice || k == reflect.String
}
