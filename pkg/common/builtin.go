package common

import (
	"bytes"
	"encoding/json"
	"reflect"
)

// DeepEqualComp compares two values of any comparable type to determine if they are deeply equal.
//
// This function uses the `reflect.DeepEqualComp` function from the `reflect` package to compare
// two values `a` and `b`. It checks for deep equality, meaning it considers nested structures,
// such as slices, maps, or structs, and compares them element-by-element or field-by-field.
// If the values are deeply equal, the function returns `true`; otherwise, it returns `false`.
//
// The function is generic, allowing it to work with any type `T` that is comparable, including
// basic types (e.g., integers, strings) as well as complex types with nested structures.
//
// Parameters:
//   - `a`: The first value to compare. It can be of any comparable type `T`.
//   - `b`: The second value to compare. It must be of the same type `T` as `a`.
//
// Returns:
//   - `true` if `a` and `b` are deeply equal; `false` otherwise.
//
// Example:
//
//	// Comparing two integer values
//	isEqual := DeepEqualComp(5, 5)
//	// isEqual will be true as both integers are equal
//
//	// Comparing two slices with the same elements
//	sliceA := []int{1, 2, 3}
//	sliceB := []int{1, 2, 3}
//	isEqual = DeepEqualComp(sliceA, sliceB)
//	// isEqual will be true as both slices have identical elements in the same order
//
//	// Comparing two different maps
//	mapA := map[string]int{"a": 1, "b": 2}
//	mapB := map[string]int{"a": 1, "b": 3}
//	isEqual = DeepEqualComp(mapA, mapB)
//	// isEqual will be false as the values for key "b" differ between the maps
func DeepEqualComp[T comparable](a, b T) bool {
	return reflect.DeepEqual(a, b)
}

// DeepEqual compares two values for equality via JSON serialization.
//
// This function serializes both input values `a` and `b` to JSON format using the `json.Marshal` function.
// It then compares the resulting JSON byte slices using `bytes.Equal`. If the serialized JSON representations
// are identical, the function returns `true`, indicating that the two values are considered equal in terms of their JSON representation.
//
// Parameters:
//   - `a`: The first value to compare. It can be of any type.
//   - `b`: The second value to compare. It can also be of any type.
//
// Returns:
//   - `true` if the JSON representations of `a` and `b` are equal; `false` otherwise.
//
// Example:
//
//	// Comparing two structs with the same data
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//	personA := Person{Name: "Alice", Age: 30}
//	personB := Person{Name: "Alice", Age: 30}
//	isEqual := DeepEqual(personA, personB)
//	// isEqual will be true as both structs serialize to the same JSON
//	// Comparing two different maps
//	mapA := map[string]int{"a": 1, "b": 2}
//	mapB := map[string]int{"a": 1, "b": 3}
//	isEqual = DeepEqual(mapA, mapB)
//	// isEqual will be false as the JSON representations differ
func DeepEqual(a, b any) bool {
	aJSON, err := json.Marshal(a)
	if err != nil {
		return false
	}

	bJSON, err := json.Marshal(b)
	if err != nil {
		return false
	}

	return bytes.Equal(aJSON, bJSON)
}

// IsScalarType checks whether the given value is a primitive type in Go.
//
// Primitive types include:
//   - Signed integers: int, int8, int16, int32, int64
//   - Unsigned integers: uint, uint8, uint16, uint32, uint64, uintptr
//   - Floating-point numbers: float32, float64
//   - Complex numbers: complex64, complex128
//   - Boolean: bool
//   - String: string
//
// The function first checks if the input value is `nil`, returning `false` if so. Otherwise, it uses reflection to determine
// the type of the value and compares it against known primitive types.
//
// Parameters:
//   - `value`: An interface{} that can hold any Go value. The function checks the type of this value.
//
// Returns:
//   - `true` if the value is of a primitive type.
//   - `false` if the value is `nil` or not a primitive type.
//
// Example:
//
//	primitive := 42
//	if IsScalarType(primitive) {
//	    fmt.Println("The value is a primitive type.")
//	} else {
//	    fmt.Println("The value is not a primitive type.")
//	}
func IsScalarType(value any) bool {
	if value == nil {
		return false
	}
	switch reflect.TypeOf(value).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.Bool, reflect.String:
		return true
	default:
		return false
	}
}

// IsEmptyValue checks whether the given reflect.Value is considered empty.
//
// An empty value is defined as:
//   - Zero length for arrays, maps, slices, and strings.
//   - False for booleans.
//   - Zero for numeric types (int, uint, float).
//   - Nil for interfaces and pointers.
//   - Zero value for structs.
//
// Parameters:
//   - v: The reflect.Value to check.
//
// Returns:
//   - true if the value is empty, false otherwise.
//
// Example:
//
//	val := reflect.ValueOf("")
//	if IsEmptyValue(val) {
//	    fmt.Println("The value is empty.")
//	} else {
//	    fmt.Println("The value is not empty.")
//	}
func IsEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Struct:
		return v.IsZero()
	}
	return false
}
