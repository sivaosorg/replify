package encoding

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/sivaosorg/replify/pkg/common"
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

// JsonSafe converts a Go value to its JSON string representation or returns the value directly if it is already a string.
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
//	jsonStr := JsonSafe(myStruct)
func JsonSafe(data any) string {
	if data == nil {
		return ""
	}
	value := reflect.ValueOf(data)
	if common.IsEmptyValue(value) {
		return ""
	}
	// Check for scalar types and convert directly
	// This avoids unnecessary JSON marshalling for simple types
	ok := common.IsScalarType(data)
	if ok {
		// return fmt.Sprintf("%v", data)
		return value.Interface().(string)
	}

	result, err := MarshalToString(data)
	if err != nil {
		return ""
	}
	return result
}

// JsonSafePretty converts a Go value to its pretty-printed JSON string representation or returns the value directly if it is already a string.
//
// This function checks if the input data is a string; if so, it returns it directly.
// Otherwise, it marshals the input value `data` into a formatted JSON string using
// the MarshalIndent function. If an error occurs during marshalling, it returns an empty string.
//
// Parameters:
//   - `data`: The Go value to be converted to pretty-printed JSON, or a string to be returned directly.
//
// Returns:
//   - A string containing the pretty-printed JSON representation of the input value, or an empty string if an error occurs.
//
// Example:
//
//	jsonPrettyStr := JsonSafePretty(myStruct)
func JsonSafePretty(data any) string {
	if data == nil {
		return ""
	}
	// Check for scalar types and convert directly
	// This avoids unnecessary JSON marshalling for simple types
	ok := common.IsScalarType(data)
	if ok {
		return fmt.Sprintf("%v", data)
	}

	result, err := MarshalIndent(data, "", "    ")
	if err != nil {
		return ""
	}
	return string(result)
}
