package encoding

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

// Toggle to choose how to handle NaN/±Inf floats in *safe* variants.
// When true: produce "null" (JSON-safe). When false: treat as error.
const floatsUseNullForNonFinite = true

// Error messages for JSON operations.
//
//   - ErrNilInterface is returned when a nil interface is passed to a JSON function.
//   - ErrInvalidRawMessage is returned when an invalid json.RawMessage is passed to a JSON function.
//   - ErrNonFiniteFloat is returned when a non-finite float (NaN/Inf) is passed to a JSON function.
//   - ErrUnsupportedValue is returned when an unsupported value (e.g., non-nil func, chan, etc.) is passed to a JSON function.
//   - ErrMarshalPanicRecovered is returned when a panic occurs during JSON marshalling.
var (
	ErrNilInterface          = errors.New("nil interface input")
	ErrInvalidRawMessage     = errors.New("invalid json.RawMessage")
	ErrNonFiniteFloat        = errors.New("non-finite float (NaN/Inf)")
	ErrUnsupportedValue      = errors.New("unsupported value (e.g., non-nil func, chan, etc.)")
	ErrMarshalPanicRecovered = errors.New("json marshal panic recovered")
)

// Marshal converts a Go value into its JSON byte representation.
//
// This function marshals the input value `v` using the standard json library.
// The resulting JSON data is returned as a byte slice. If there is an error
// during marshalling, it returns the error.
//
// Parameters:
//   - `v`: The Go value to be marshalled into JSON.
//
// Returns:
//   - A byte slice containing the JSON representation of the input value.
//   - An error if the marshalling fails.
//
// Example:
//
//	jsonData, err := Marshal(myStruct)
func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// MarshalIndent converts a Go value to its JSON string representation with indentation.
//
// This function marshals the input value `v` into a formatted JSON string,
// allowing for easy readability by including a specified prefix and indentation.
// It returns the resulting JSON byte slice or an error if marshalling fails.
//
// Parameters:
//   - `v`: The Go value to be marshalled into JSON.
//   - `prefix`: A string that will be prefixed to each line of the output JSON.
//   - `indent`: A string used for indentation (typically a series of spaces or a tab).
//
// Returns:
//   - A byte slice containing the formatted JSON representation of the input value.
//   - An error if the marshalling fails.
//
// Example:
//
//	jsonIndented, err := MarshalIndent(myStruct, "", "    ")
func MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

// MarshalToString converts a Go value to its JSON string representation.
//
// This function utilizes the standard json library to marshal the input value `v`
// into a JSON string. If the marshalling is successful, it returns the resulting
// JSON string. If an error occurs during the process, it returns an error.
//
// Parameters:
//   - `v`: The Go value to be marshalled into JSON.
//
// Returns:
//   - A string containing the JSON representation of the input value.
//   - An error if the marshalling fails.
//
// Example:
//
//	jsonString, err := MarshalToString(myStruct)
func MarshalToString(v any) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Unmarshal parses JSON-encoded data and stores the result in the value pointed to by `v`.
//
// This function uses the standard json library to unmarshal JSON data
// (given as a byte slice) into the specified Go value `v`. If the unmarshalling
// is successful, it populates the value `v`. If an error occurs, it returns the error.
//
// Parameters:
//   - `data`: A byte slice containing JSON data to be unmarshalled.
//   - `v`: A pointer to the Go value where the unmarshalled data will be stored.
//
// Returns:
//   - An error if the unmarshalling fails.
//
// Example:
//
//	err := Unmarshal(jsonData, &myStruct)
func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// UnmarshalFromString parses JSON-encoded string and stores the result in the value pointed to by `v`.
//
// This function utilizes the standard json library to unmarshal JSON data
// from a string into the specified Go value `v`. If the unmarshalling is
// successful, it populates the value `v`. If an error occurs, it returns the error.
//
// Parameters:
//   - `str`: A string containing JSON data to be unmarshalled.
//   - `v`: A pointer to the Go value where the unmarshalled data will be stored.
//
// Returns:
//   - An error if the unmarshalling fails.
//
// Example:
//
//	err := UnmarshalFromString(jsonString, &myStruct)
func UnmarshalFromString(str string, v any) error {
	return json.Unmarshal([]byte(str), v)
}

// IsValidJSON checks if a given string is a valid JSON format.
//
// This function uses the json.Valid method from the standard json library
// to determine if the input string `s` is a valid JSON representation.
//
// Parameters:
//   - `s`: The string to be validated as JSON.
//
// Returns:
//   - A boolean indicating whether the input string is valid JSON.
func IsValidJSON(s string) bool {
	return json.Valid([]byte(s))
}

// IsValidJSONBytes checks if a given byte slice is a valid JSON format.
//
// This function uses the json.Valid method from the standard json library
// to determine if the input byte slice `data` is a valid JSON representation.
//
// Parameters:
//   - `data`: The byte slice to be validated as JSON.
//
// Returns:
//   - A boolean indicating whether the input byte slice is valid JSON.
func IsValidJSONBytes(data []byte) bool {
	return json.Valid(data)
}

// JSONSafe converts a Go value to its JSON string representation or returns an error if the marshalling fails.
// It uses a deferred function to recover from any panics that may occur during marshalling.
//
// Parameters:
//   - `data`: The Go value to be converted to JSON.
//
// Returns:
//   - A string containing the JSON representation of the input value.
//   - An error if the marshalling fails.
//
// Example:
//
//	 var myStruct = struct {
//	 	Name string
//	 	Age  int
//	 }{
//	 	Name: "John",
//	 	Age:  30,
//	 }
//		jsonString, err := JSONSafe(myStruct)
func JSONSafe(data any) string {
	return jsonSafe(data, false)
}

// JSONSafeToken converts a Go value to its JSON string representation or returns an error if the marshalling fails.
// It uses a deferred function to recover from any panics that may occur during marshalling.
//
// Parameters:
//   - `data`: The Go value to be converted to JSON.
//
// Returns:
//   - A string containing the JSON representation of the input value.
//   - An error if the marshalling fails.
//
// Example:
//
//	 var myStruct = struct {
//	 	Name string
//	 	Age  int
//	 }{
//	 	Name: "John",
//	 	Age:  30,
//	 }
//		jsonString, err := JSONSafeToken(myStruct)
func JSONSafeToken(data any) (string, error) {
	return jsonSafeToken(data, false)
}

// JSONSafePretty converts a Go value to its pretty-printed JSON string representation or returns an error if the marshalling fails.
// It uses a deferred function to recover from any panics that may occur during marshalling.
//
// Parameters:
//   - `data`: The Go value to be converted to JSON.
//
// Returns:
//   - A string containing the JSON representation of the input value.
//   - An error if the marshalling fails.
//
// Example:
//
//	 var myStruct = struct {
//	 	Name string
//	 	Age  int
//	 }{
//	 	Name: "John",
//	 	Age:  30,
//	 }
//		jsonString, err := JSONSafePretty(myStruct)
func JSONSafePretty(data any) string {
	return jsonSafe(data, true)
}

// JSONSafePrettyToken converts a Go value to its pretty-printed JSON string representation or returns an error if the marshalling fails.
// It uses a deferred function to recover from any panics that may occur during marshalling.
//
// Parameters:
//   - `data`: The Go value to be converted to JSON.
//
// Returns:
//   - A string containing the JSON representation of the input value.
//   - An error if the marshalling fails.
//
// Example:
//
//	 var myStruct = struct {
//	 	Name string
//	 	Age  int
//	 }{
//	 	Name: "John",
//	 	Age:  30,
//	 }
//		jsonString, err := JSONSafePrettyToken(myStruct)
func JSONSafePrettyToken(data any) (string, error) {
	return jsonSafeToken(data, true)
}

// marshalToStrRecover marshals a Go value to its JSON string representation or returns an error if the marshalling fails.
// It uses a deferred function to recover from any panics that may occur during marshalling.
//
// Parameters:
//   - `v`: The Go value to be marshalled into JSON.
//   - `pretty`: A boolean indicating whether the JSON should be pretty-printed.
//
// Returns:
//   - A string containing the JSON representation of the input value.
//   - An error if the marshalling fails.
//
// Example:
//
//	jsonString, err := marshalToStrRecover(myStruct, false)
func marshalToStrRecover(v any, pretty bool) (out string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%w: %v", ErrMarshalPanicRecovered, r)
			out = ""
		}
	}()

	var b []byte
	if pretty {
		b, err = MarshalIndent(v, "", "    ")
	} else {
		b, err = Marshal(v)
	}
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// realFrom extracts the real part of a complex number from a reflect.Value.
//
// Parameters:
//   - `v`: The reflect.Value to extract the real part from.
//
// Returns:
//   - The real part of the complex number as a float64.
//
// Example:
//
//	realPart := realFrom(reflect.ValueOf(complex(1.23, 4.56)))
func realFrom(v reflect.Value) float64 {
	switch v.Kind() {
	case reflect.Complex64:
		return float64(real(v.Complex()))
	case reflect.Complex128:
		return real(v.Complex())
	default:
		return 0
	}
}

// imagFrom extracts the imaginary part of a complex number from a reflect.Value.
//
// Parameters:
//   - `v`: The reflect.Value to extract the imaginary part from.
//
// Returns:
//   - The imaginary part of the complex number as a float64.
//
// Example:
//
//	imagPart := imagFrom(reflect.ValueOf(complex(1.23, 4.56)))
func imagFrom(v reflect.Value) float64 {
	switch v.Kind() {
	case reflect.Complex64:
		return float64(imag(v.Complex()))
	case reflect.Complex128:
		return imag(v.Complex())
	default:
		return 0
	}
}

// formatFloatJSON formats a float64 to its JSON string representation.
// It uses 'g' formatting like encoding/json and converts non-finite numbers to null.
//
// Parameters:
//   - `f`: The float64 value to format.
//   - `is32`: A boolean indicating whether the float is a float32 (true) or float64 (false).
//
// Returns:
//   - A string containing the JSON representation of the float.
//
// Example:
//
//	jsonFloat := formatFloatJSON(1.2345, false)
func formatFloatJSON(f float64, is32 bool) string {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		if floatsUseNullForNonFinite {
			return "null"
		}
		return ""
	}
	bitSize := 64
	if is32 {
		bitSize = 32
	}
	// 'g' format is used for general-purpose formatting. It uses the shortest
	// representation of the float, either in decimal or scientific notation.
	// The -1 precision means that the smallest number of digits necessary to
	// represent the float will be used.
	return strconv.FormatFloat(f, 'g', -1, bitSize)
}

// encodeComplexJSON encodes a complex number to its JSON string representation.
// It uses 'g' formatting like encoding/json and converts non-finite numbers to null.
//
// Parameters:
//   - `realPart`: The real part of the complex number.
//   - `imagPart`: The imaginary part of the complex number.
//   - `is32`: A boolean indicating whether the complex number is a complex64 (true) or complex128 (false).
//
// Returns:
//   - A string containing the JSON representation of the complex number.
//
// Example:
//
//	r := encodeComplexJSON(1.2345, 6.789, false)
func encodeComplexJSON(realPart, imagPart float64, is32 bool) string {
	// Use 'g' formatting like encoding/json; convert non-finite to null.
	r := formatFloatJSON(realPart, is32)
	i := formatFloatJSON(imagPart, is32)
	if r == "" || i == "" {
		return ""
	}
	// Format complex number as JSON object with "real" and "imag" fields.
	// This is done to avoid the default JSON representation of complex numbers,
	// which is not valid JSON.
	return `{"real":` + r + `,"imag":` + i + `}`
}

// encodeComplexJSONToken encodes a complex number to its JSON string representation.
// It uses 'g' formatting like encoding/json and converts non-finite numbers to null.
//
// Parameters:
//   - `realPart`: The real part of the complex number.
//   - `imagPart`: The imaginary part of the complex number.
//   - `is32`: A boolean indicating whether the complex number is a complex64 (true) or complex128 (false).
//
// Returns:
//   - A string containing the JSON representation of the complex number.
//   - An error if the marshalling fails.
//
// Example:
//
//	r := encodeComplexJSONToken(1.2345, 6.789, false)
func encodeComplexJSONToken(realPart, imagPart float64, is32 bool) (string, error) {
	r := formatFloatJSON(realPart, is32)
	i := formatFloatJSON(imagPart, is32)
	// If non-finite handling is set to error, formatFloatJSON returns "" and we should error.
	if r == "" || i == "" {
		// Consistent with float policy: if not using "null", report ErrNonFiniteFloat.
		if !floatsUseNullForNonFinite {
			return "", ErrNonFiniteFloat
		}
		// If using "null", r or i would be "null", so we can still construct the object.
		if r == "" {
			r = "null"
		}
		if i == "" {
			i = "null"
		}
	}
	return `{"real":` + r + `,"imag":` + i + `}`, nil
}

// jsonSafe converts a Go value to its JSON string representation or returns an error if the marshalling fails.
// It uses a deferred function to recover from any panics that may occur during marshalling.
//
// Parameters:
//   - `data`: The Go value to be converted to JSON.
//   - `pretty`: A boolean indicating whether the JSON should be pretty-printed.
//
// Returns:
//   - A string containing the JSON representation of the input value.
//   - An error if the marshalling fails.
//
// Example:
//
//	jsonString, err := jsonSafe(myStruct, false)
func jsonSafe(data any, pretty bool) string {
	if data == nil {
		return ""
	}

	// 1) Pass-through raw JSON if explicitly provided.
	if rm, ok := data.(json.RawMessage); ok {
		if rm == nil {
			return "null"
		}
		if !json.Valid(rm) {
			// Invalid raw JSON => treat as error sentinel.
			return ""
		}
		if !pretty {
			return string(rm)
		}
		var buf bytes.Buffer
		// json.Indent is a no-op for primitives; safe to call for any valid JSON value.
		if err := json.Indent(&buf, rm, "", "    "); err != nil {
			return ""
		}
		return buf.String()
	}

	v := reflect.ValueOf(data)

	// 2) Unwrap *all* interface layers; nil interface => "".
	for v.Kind() == reflect.Interface {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}
	data = v.Interface()

	// 3) Nil-able kinds (Ptr, Map, Slice, Func, Chan) when nil => "null".
	v = reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan:
		if v.IsNil() {
			return "null"
		}
	}

	// 4) Scalar fast-path as valid JSON tokens.
	switch v.Kind() {
	case reflect.String:
		// MUST quote/escape to be a valid JSON string token.
		b, err := json.Marshal(v.String())
		if err != nil {
			return ""
		}
		return string(b)

	case reflect.Bool:
		return strconv.FormatBool(v.Bool())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)

	case reflect.Uintptr:
		// Avoid precision loss & misleading semantics; emit as quoted hex address.
		return fmt.Sprintf("%q", fmt.Sprintf("0x%x", v.Uint()))

	case reflect.Float32, reflect.Float64:
		f := v.Float()
		if math.IsNaN(f) || math.IsInf(f, 0) {
			if floatsUseNullForNonFinite {
				return "null"
			}
			// Optional legacy behavior: treat as marshal error sentinel.
			return ""
		}
		bitSize := 64
		if v.Kind() == reflect.Float32 {
			bitSize = 32
		}
		// Match encoding/json's number formatting preference.
		return strconv.FormatFloat(f, 'g', -1, bitSize)

	case reflect.Complex64, reflect.Complex128:
		// JSON has no complex type; emit as an object {"real":..., "imag":...}.
		r, i := realFrom(v), imagFrom(v)
		return encodeComplexJSON(r, i, v.Kind() == reflect.Complex64)
	}

	// 5) Everything else: marshal with panic protection and optional pretty indent.
	s, err := marshalToStrRecover(data, pretty)
	if err != nil {
		return ""
	}
	return s
}

// jsonSafeToken converts a Go value to its JSON string representation or returns an error if the marshalling fails.
// It uses a deferred function to recover from any panics that may occur during marshalling.
//
// Parameters:
//   - `data`: The Go value to be converted to JSON.
//   - `pretty`: A boolean indicating whether the JSON should be pretty-printed.
//
// Returns:
//   - A string containing the JSON representation of the input value.
//   - An error if the marshalling fails.
//
// Example:
//
//	jsonString, err := jsonSafeToken(myStruct, false)
func jsonSafeToken(data any, pretty bool) (string, error) {
	if data == nil {
		return "", ErrNilInterface
	}

	// 1) Pass-through raw JSON if explicitly provided.
	if rm, ok := data.(json.RawMessage); ok {
		if rm == nil {
			return "null", nil
		}
		if !json.Valid(rm) {
			return "", ErrInvalidRawMessage
		}
		if !pretty {
			return string(rm), nil
		}
		var buf bytes.Buffer
		if err := json.Indent(&buf, rm, "", "    "); err != nil {
			return "", err
		}
		return buf.String(), nil
	}

	v := reflect.ValueOf(data)

	// 2) Unwrap *all* interface layers; nil interface => error.
	for v.Kind() == reflect.Interface {
		if v.IsNil() {
			return "", ErrNilInterface
		}
		v = v.Elem()
	}
	data = v.Interface()

	// 3) Nil-able kinds (Ptr, Map, Slice, Func, Chan) when nil => "null".
	v = reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan:
		if v.IsNil() {
			return "null", nil
		}
	}

	// 4) Scalar fast-path as valid JSON tokens.
	switch v.Kind() {
	case reflect.String:
		// MUST quote/escape to be a valid JSON string token.
		b, err := json.Marshal(v.String())
		if err != nil {
			return "", err
		}
		return string(b), nil

	case reflect.Bool:
		return strconv.FormatBool(v.Bool()), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10), nil

	case reflect.Uintptr:
		// Avoid precision loss & misleading semantics; emit as quoted hex address.
		return fmt.Sprintf("%q", fmt.Sprintf("0x%x", v.Uint())), nil

	case reflect.Float32, reflect.Float64:
		f := v.Float()
		if math.IsNaN(f) || math.IsInf(f, 0) {
			if floatsUseNullForNonFinite {
				return "null", nil
			}
			return "", ErrNonFiniteFloat
		}
		bitSize := 64
		if v.Kind() == reflect.Float32 {
			bitSize = 32
		}
		return strconv.FormatFloat(f, 'g', -1, bitSize), nil

	case reflect.Complex64, reflect.Complex128:
		// Encode as {"real":...,"imag":...}
		r, i := realFrom(v), imagFrom(v)
		s, err := encodeComplexJSONToken(r, i, v.Kind() == reflect.Complex64)
		return s, err
	}

	// 5) Everything else: marshal with panic protection and optional pretty indent.
	return marshalToStrRecover(data, pretty)
}
