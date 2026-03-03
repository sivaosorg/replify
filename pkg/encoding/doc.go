// Package encoding provides utilities for marshalling, unmarshalling,
// validating, normalising, and pretty-printing JSON data.
//
// The package wraps the standard [encoding/json] library with additional
// safety, convenience, and formatting capabilities used throughout replify.
//
// # Marshalling and Unmarshalling
//
// MarshalJSONb and MarshalJSONs marshal any Go value to []byte or string
// respectively. MarshalJSONIndent produces human-readable output with a
// configurable prefix and indent string. UnmarshalBytes and UnmarshalJSON
// are thin wrappers around json.Unmarshal; their Safe variants additionally
// reject empty input and syntactically invalid JSON before attempting
// deserialisation.
//
// # Validation and Normalisation
//
// IsValidJSON and IsValidJSONBytes report whether a string or byte slice
// constitutes valid JSON. NormalizeJSON attempts to repair common
// malformations — leading UTF-8 BOMs, embedded null bytes, escaped
// structural quotes, and trailing commas — returning valid JSON or an error:
//
//	normalized, err := encoding.NormalizeJSON(`{\"key\": "value",}`)
//	// normalized: {"key": "value"}
//
// # Safe JSON Conversion
//
// JSON and JSONToken convert any Go value to a JSON string, handling nil
// interfaces, raw [json.RawMessage] pass-through, scalar fast paths, and
// non-finite floating-point values (NaN / ±Inf are mapped to "null").
// JSONPretty and JSONPrettyToken produce indented output. The Token variants
// return an explicit error instead of an empty string sentinel.
//
// # Pretty Printing
//
// Pretty and Color format an existing JSON byte slice with configurable
// indentation, key/value alignment, and optional ANSI terminal colouring.
// Several predefined Style values are provided: TerminalStyle,
// VSCodeDarkStyle, DraculaStyle, MonokaiStyle, SolarizedDarkStyle, and
// MinimalGrayStyle.
package encoding
