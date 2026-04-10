package common

import "reflect"

// sliceTypeOf returns the reflect.Type of a slice whose element type matches the
// element type of v.  When v is already a slice its type is returned as-is;
// when v is an array, reflect.SliceOf(elem) is returned so that reflect.MakeSlice
// does not panic (reflect.MakeSlice only accepts slice types, never array types).
func sliceTypeOf(v reflect.Value) reflect.Type {
	if v.Kind() == reflect.Array {
		return reflect.SliceOf(v.Type().Elem())
	}
	return v.Type()
}
