package conv

import (
	"reflect"
)

// ///////////////////////////
// Section:  Infer conversion
// ///////////////////////////

// Infer will perform conversion by inferring the conversion operation from
// the base type of a pointer to a supported type.  The value is assigned
// directly so only an error is returned.
//
// Example:
//
//	var into int64
//	err := conv.Infer(&into, "42")
//	// into -> 42
//
//	var t time.Time
//	err := conv.Infer(&t, "2024-01-15")
//	// t -> 2024-01-15 00:00:00 +0000 UTC
//
// Parameters:
//   - into: A pointer to the target variable.
//   - from: The source value to convert.
//
// Returns:
//   - An error if the conversion fails or if `into` is not a settable pointer.
func (c *Converter) Infer(into, from any) error {
	// Get reflect.Value of target
	var value reflect.Value
	switch v := into.(type) {
	case reflect.Value:
		value = v
	default:
		value = reflect.ValueOf(into)
	}

	// Validate the value
	if !value.IsValid() {
		return newConvErrorf("%T is not a valid value", into)
	}

	// Must be a pointer
	if value.Kind() != reflect.Ptr {
		return newConvErrorf("cannot infer conversion for non-pointer %T", into)
	}

	// Dereference the pointer
	value = value.Elem()

	// Must be settable
	if !value.CanSet() {
		return newConvErrorf("cannot infer conversion for unchangeable %v (type %[1]T)", into)
	}

	// Perform the conversion
	result, err := c.inferValue(value, from)
	if err != nil {
		return err
	}

	// Set the value
	value.Set(reflect.ValueOf(result))
	return nil
}

// ///////////////////////////
// Section: Package-level Infer helpers
// ///////////////////////////

// TryInfer attempts to infer conversion, returning success status.
//
// Example:
//
//	val, ok := conv.TryInfer[int]("42")
//	// // val -> 42, ok -> true
//	val2, ok2 := conv.TryInfer[int]("not an int")
//	// // val2 -> 0, ok2 -> false
//
// Parameters:
//   - from: The source value to convert.
//
// Returns:
//   - The converted value.
func TryInfer[T any](from any) (T, bool) {
	var into T
	if err := defaultConverter.Infer(&into, from); err != nil {
		return into, false
	}
	return into, true
}

// inferValue infers the type from reflect.Value and performs conversion.
//
// Parameters:
//   - val: The reflect.Value representing the target type.
//   - from: The value to be converted.
//
// Returns:
//   - The converted value.
//   - An error if the conversion fails.
func (c *Converter) inferValue(val reflect.Value, from any) (any, error) {
	kind := val.Kind()
	valType := val.Type()

	// Check for special types first
	if valType == typeOfTime {
		return c.Time(from)
	}
	if valType == typeOfDuration {
		return c.Duration(from)
	}

	// Handle based on kind
	switch kind {
	case reflect.Bool:
		return c.Bool(from)
	case reflect.String:
		return c.String(from)
	case reflect.Int:
		return c.Int(from)
	case reflect.Int8:
		return c.Int8(from)
	case reflect.Int16:
		return c.Int16(from)
	case reflect.Int32:
		return c.Int32(from)
	case reflect.Int64:
		return c.Int64(from)
	case reflect.Uint:
		return c.Uint(from)
	case reflect.Uint8:
		return c.Uint8(from)
	case reflect.Uint16:
		return c.Uint16(from)
	case reflect.Uint32:
		return c.Uint32(from)
	case reflect.Uint64:
		return c.Uint64(from)
	case reflect.Float32:
		return c.Float32(from)
	case reflect.Float64:
		return c.Float64(from)
	case reflect.Slice:
		return c.inferSlice(val, from)
	case reflect.Map:
		return c.inferMap(val, from)
	case reflect.Ptr:
		return c.inferPointer(val, from)
	default:
		return nil, newConvErrorf("cannot infer conversion for %v (type %v)", val, valType)
	}
}

// inferSlice infers and converts a slice type.
//
// Parameters:
//   - val: The reflect.Value representing the target slice type.
//   - from: The value to be converted.
//
// Returns:
//   - The converted slice.
//   - An error if the conversion fails.
func (c *Converter) inferSlice(val reflect.Value, from any) (any, error) {
	elemType := val.Type().Elem()

	// Handle []byte specially
	if elemType.Kind() == reflect.Uint8 {
		s, err := c.String(from)
		if err != nil {
			return nil, err
		}
		return []byte(s), nil
	}

	// Handle []string
	if elemType.Kind() == reflect.String {
		return StringSlice(from)
	}

	// For other slice types, try to convert
	fromVal := reflect.ValueOf(from)
	if fromVal.Kind() != reflect.Slice && fromVal.Kind() != reflect.Array {
		return nil, newConvErrorf("cannot convert %T to slice", from)
	}

	result := reflect.MakeSlice(val.Type(), fromVal.Len(), fromVal.Len())
	for i := 0; i < fromVal.Len(); i++ {
		elemVal := reflect.New(elemType).Elem()
		converted, err := c.inferValue(elemVal, fromVal.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		result.Index(i).Set(reflect.ValueOf(converted))
	}

	return result.Interface(), nil
}

// inferMap infers and converts a map type.
//
// Parameters:
//   - val: The reflect.Value representing the target map type.
//   - from: The value to be converted.
//
// Returns:
//   - The converted map.
//   - An error if the conversion fails.
func (c *Converter) inferMap(val reflect.Value, from any) (any, error) {
	fromVal := reflect.ValueOf(from)
	if fromVal.Kind() != reflect.Map {
		return nil, newConvErrorf("cannot convert %T to map", from)
	}

	mapType := val.Type()
	keyType := mapType.Key()
	elemType := mapType.Elem()

	result := reflect.MakeMap(mapType)
	iter := fromVal.MapRange()
	for iter.Next() {
		// Convert key
		keyVal := reflect.New(keyType).Elem()
		convertedKey, err := c.inferValue(keyVal, iter.Key().Interface())
		if err != nil {
			return nil, err
		}

		// Convert value
		elemVal := reflect.New(elemType).Elem()
		convertedValue, err := c.inferValue(elemVal, iter.Value().Interface())
		if err != nil {
			return nil, err
		}

		result.SetMapIndex(reflect.ValueOf(convertedKey), reflect.ValueOf(convertedValue))
	}

	return result.Interface(), nil
}

// inferPointer infers and converts a pointer type.
//
// Parameters:
//   - val: The reflect.Value representing the target pointer type.
//   - from: The value to be converted.
//
// Returns:
//   - The converted pointer.
//   - An error if the conversion fails.
func (c *Converter) inferPointer(val reflect.Value, from any) (any, error) {
	elemType := val.Type().Elem()
	elemVal := reflect.New(elemType).Elem()

	converted, err := c.inferValue(elemVal, from)
	if err != nil {
		return nil, err
	}

	ptr := reflect.New(elemType)
	ptr.Elem().Set(reflect.ValueOf(converted))
	return ptr.Interface(), nil
}
