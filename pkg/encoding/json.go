package encoding

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sivaosorg/replify/pkg/strutil"
)

// MarshalJSONb converts a Go value into its JSON byte representation.
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
//	jsonData, err := MarshalJSONb(myStruct)
func MarshalJSONb(v any) ([]byte, error) {
	return json.Marshal(v)
}

// MarshalJSONs converts a Go value to its JSON string representation.
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
//	jsonString, err := MarshalJSONs(myStruct)
func MarshalJSONs(v any) (string, error) {
	data, err := MarshalJSONb(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// MarshalJSONIndent converts a Go value to its JSON string representation with indentation.
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
//	jsonIndented, err := MarshalJSONIndent(myStruct, "", "    ")
func MarshalJSONIndent(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

// UnmarshalBytes parses JSON-encoded data and stores the result in the value pointed to by `v`.
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
//	err := UnmarshalBytes(jsonData, &myStruct)
func UnmarshalBytes(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// SafeUnmarshalBytes parses JSON-encoded data and stores the result in the value pointed to by `v`.
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
//	err := SafeUnmarshalBytes(jsonData, &myStruct)
func SafeUnmarshalBytes(data []byte, v any) error {
	if len(data) == 0 {
		return fmt.Errorf("%w: JSON data must not be empty", ErrEmptyInput)
	}

	if !IsValidJSONBytes(data) {
		return fmt.Errorf("%w: input is not valid JSON", ErrInvalidJSON)
	}

	return UnmarshalBytes(data, v)
}

// UnmarshalJSON parses JSON-encoded string and stores the result in the value pointed to by `v`.
//
// This function utilizes the standard json library to unmarshal JSON data
// from a string into the specified Go value `v`. If the unmarshalling is
// successful, it populates the value `v`. If an error occurs, it returns the error.
//
// Parameters:
//   - `jsonStr`: A string containing JSON data to be unmarshalled.
//   - `v`: A pointer to the Go value where the unmarshalled data will be stored.
//
// Returns:
//   - An error if the unmarshalling fails.
//
// Example:
//
//	err := UnmarshalJSON(jsonString, &myStruct)
func UnmarshalJSON(jsonStr string, v any) error {
	return json.Unmarshal([]byte(jsonStr), v)
}

// SafeUnmarshalJSON parses JSON-encoded string and stores the result in the value pointed to by `v`.
//
// This function uses the standard json library to unmarshal JSON data
// (given as a string) into the specified Go value `v`. If the unmarshalling
// is successful, it populates the value `v`. If an error occurs, it returns the error.
//
// Parameters:
//   - `jsonStr`: A string containing JSON data to be unmarshalled.
//   - `v`: A pointer to the Go value where the unmarshalled data will be stored.
//
// Returns:
//   - An error if the unmarshalling fails.
//
// Example:
//
//	err := SafeUnmarshalJSON(jsonString, &myStruct)
func SafeUnmarshalJSON(jsonStr string, v any) error {
	if strutil.IsEmpty(jsonStr) {
		return fmt.Errorf("%w: JSON string must not be empty", ErrEmptyInput)
	}

	if !IsValidJSON(jsonStr) {
		return fmt.Errorf("%w: input is not valid JSON", ErrInvalidJSON)
	}

	return UnmarshalJSON(jsonStr, v)
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

// NormalizeJSON attempts to normalize a malformed JSON-like string into valid JSON.
//
// Normalization strategy — passes are applied in sequence; validity is checked
// after each pass that modifies the candidate. The function returns as soon as
// a pass produces a valid JSON string, so no unnecessary work is performed.
//
//  1. Empty / whitespace-only input → return error.
//  2. Already valid JSON → return unchanged (fast path, no allocation).
//  3. Pass 1 – strip a leading UTF-8 BOM (U+FEFF / 0xEF 0xBB 0xBF).
//  4. Pass 2 – remove embedded null bytes (0x00) which are invalid inside JSON text.
//  5. Pass 3 – unescape literal `\"` sequences to `"`.  This is the most common
//     artifact produced when JSON is stored in Go raw string literals or travels
//     through systems that double-escape structural quote characters.
//  6. Pass 4 – remove trailing commas before `}` or `]`.  These are produced by
//     some serializers and are not permitted by the JSON grammar.
//
// Passes are cumulative: each pass operates on the output of the previous one.
// The function does NOT silently corrupt already-valid JSON (step 2 guarantees
// this). Only inputs that fail the initial validation ever enter the pass chain.
//
// Parameters:
//   - s: The input string to normalize.
//
// Returns:
//   - A valid JSON string on success.
//   - An error if the input is empty/whitespace or cannot be normalized to valid JSON.
//
// Example:
//
//	normalized, err := NormalizeJSON(`{\"key\": "value"}`)
func NormalizeJSON(s string) (string, error) {
	if strutil.IsEmpty(s) {
		return "", fmt.Errorf("%w: JSON string must not be empty", ErrEmptyInput)
	}

	// Fast path: already valid JSON — return as-is with no allocation.
	if IsValidJSON(s) {
		return s, nil
	}

	candidate := s

	// Pass 1: Strip leading UTF-8 BOM (0xEF 0xBB 0xBF).
	if strings.HasPrefix(candidate, "\xEF\xBB\xBF") {
		candidate = candidate[3:]
		if IsValidJSON(candidate) {
			return candidate, nil
		}
	}

	// Pass 2: Remove embedded null bytes.
	if strings.Contains(candidate, "\x00") {
		candidate = strings.ReplaceAll(candidate, "\x00", "")
		if IsValidJSON(candidate) {
			return candidate, nil
		}
	}

	// Pass 3: Unescape literal \" → " (structural quote escape artifacts).
	if strings.Contains(candidate, `\"`) {
		candidate = strings.ReplaceAll(candidate, `\"`, `"`)
		if IsValidJSON(candidate) {
			return candidate, nil
		}
	}

	// Pass 4: Remove trailing commas before } or ] (invalid in JSON grammar).
	if noTrailing := normalizeTrailingCommaRe.ReplaceAllString(candidate, "$1"); noTrailing != candidate {
		candidate = noTrailing
		if IsValidJSON(candidate) {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("%w: input could not be repaired to valid JSON", ErrInvalidJSON)
}

// JSON converts a Go value to its JSON string representation or returns an error if the marshalling fails.
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
//		jsonString, err := JSON(myStruct)
func JSON(data any) string {
	return jsonSafe(data, false)
}

// JSONToken converts a Go value to its JSON string representation or returns an error if the marshalling fails.
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
//		jsonString, err := JSONToken(myStruct)
func JSONToken(data any) (string, error) {
	return jsonSafeToken(data, false)
}

// JSONPretty converts a Go value to its pretty-printed JSON string representation or returns an error if the marshalling fails.
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
//		jsonString, err := JSONPretty(myStruct)
func JSONPretty(data any) string {
	return jsonSafe(data, true)
}

// JSONPrettyToken converts a Go value to its pretty-printed JSON string representation or returns an error if the marshalling fails.
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
//		jsonString, err := JSONPrettyToken(myStruct)
func JSONPrettyToken(data any) (string, error) {
	return jsonSafeToken(data, true)
}
