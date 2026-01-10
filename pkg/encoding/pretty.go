package encoding

import (
	"bytes"
	"encoding/json"
	"sort"
	"strconv"
)

// OptionsConfig defines the configuration options for pretty-printing JSON data.
// It allows customization of width, prefix, indentation, and sorting of keys.
// These options control how the JSON output will be formatted.
//
// Fields:
//   - Width: The maximum column width for single-line arrays. This prevents arrays from becoming too wide.
//     Default is 80 characters.
//   - Prefix: A string that will be prepended to each line of the output. Useful for adding custom prefixes
//     or structuring the output with additional information. Default is an empty string.
//   - Indent: The string used for indentation in nested JSON structures. Default is two spaces ("  ").
//   - SortKeys: A flag indicating whether the keys in JSON objects should be sorted alphabetically. Default is false.
type OptionsConfig struct {
	// Width is an max column width for single line arrays
	// Default is 80
	Width int `json:"width"`
	// Prefix is a prefix for all lines
	// Default is an empty string
	Prefix string `json:"prefix"`
	// Indent is the nested indentation
	// Default is two spaces
	Indent string `json:"indent"`
	// SortKeys will sort the keys alphabetically
	// Default is false
	SortKeys bool `json:"sort_keys"`
}

// Style is the color style
type Style struct {
	Key, String, Number [2]string
	True, False, Null   [2]string
	Escape              [2]string
	Brackets            [2]string
	Append              func(dst []byte, c byte) []byte
}

// TerminalStyle is for terminals
var TerminalStyle *Style

// VSCodeDarkStyle is for VS Code dark theme
var VSCodeDarkStyle *Style

// DraculaStyle is for Dracula theme
var DraculaStyle *Style

// MonokaiStyle is for Monokai theme
var MonokaiStyle *Style

// SolarizedDarkStyle is for Solarized Dark theme
var SolarizedDarkStyle *Style

// MinimalGrayStyle is a minimal gray style
var MinimalGrayStyle *Style

// DefaultOptionsConfig is a pre-configured default set of options for pretty-printing JSON.
// This configuration uses a width of 80, an empty prefix, two-space indentation, and does not sort keys.
// It is used when no custom options are provided in the PrettyOptions function.
var DefaultOptionsConfig = &OptionsConfig{Width: 80, Prefix: "", Indent: "  ", SortKeys: false}

// jsonType represents the different types of JSON values.
//
// This enumeration defines constants representing various JSON data types, including `null`, `boolean`, `number`,
// `string`, and `JSON object or array`. These constants are used by the `getJsonType` function to identify the type
// of a given JSON value based on its first character.
type jsonType int

// sortCriteria represents the criteria for sorting JSON key-value pairs.
// It is used in the byKeyVal struct's isLess method to determine whether to compare by key or by value.
type sortCriteria int

// Constants for sorting by key or by value
// Used in the byKeyVal struct's isLess method.
const (
	sKey sortCriteria = 0
	sVal sortCriteria = 1
)

// Constants representing different JSON types.
// These constants are used to identify the type of a JSON value
// based on its first character.
const (
	jsonNull   jsonType = iota // Represents a JSON null value
	jsonFalse                  // Represents a JSON false boolean
	jNumber                    // Represents a JSON number
	jsonString                 // Represents a JSON string
	jsonTrue                   // Represents a JSON true boolean
	jsonJson                   // Represents a JSON object or array
)

// kvPairs represents the positions of a key-value pair in a JSON object.
// It contains the start and end indices for both the key and the value within the JSON byte slice.
type kvPairs struct {
	keyStart, keyEnd     int // Indices for the key, inclusive start and exclusive end
	valueStart, valueEnd int // Indices for the value, inclusive start and exclusive end
}

// kvSorter is a struct that provides a way to sort JSON key-value pairs.
// It contains the JSON data, a buffer to hold trimmed values, and a list of pairs to be sorted.
type kvSorter struct {
	sorted bool      // indicates whether the pairs are sorted
	json   []byte    // original JSON data
	buf    []byte    // buffer used for processing values
	pairs  []kvPairs // list of key-value pairs to sort
}

func init() {
	TerminalStyle = &Style{
		Key:      [2]string{"\x1B[1m\x1B[94m", "\x1B[0m"},
		String:   [2]string{"\x1B[32m", "\x1B[0m"},
		Number:   [2]string{"\x1B[33m", "\x1B[0m"},
		True:     [2]string{"\x1B[36m", "\x1B[0m"},
		False:    [2]string{"\x1B[36m", "\x1B[0m"},
		Null:     [2]string{"\x1B[2m", "\x1B[0m"},
		Escape:   [2]string{"\x1B[35m", "\x1B[0m"},
		Brackets: [2]string{"\x1B[1m", "\x1B[0m"},
		Append:   defaultStyleAppend,
	}

	VSCodeDarkStyle = &Style{
		Key:      [2]string{"\x1B[38;5;81m", "\x1B[0m"},  // blue
		String:   [2]string{"\x1B[38;5;114m", "\x1B[0m"}, // green
		Number:   [2]string{"\x1B[38;5;209m", "\x1B[0m"}, // orange
		True:     [2]string{"\x1B[38;5;80m", "\x1B[0m"},  // cyan
		False:    [2]string{"\x1B[38;5;80m", "\x1B[0m"},
		Null:     [2]string{"\x1B[38;5;244m", "\x1B[0m"}, // gray
		Escape:   [2]string{"\x1B[38;5;176m", "\x1B[0m"}, // purple
		Brackets: [2]string{"\x1B[38;5;250m", "\x1B[0m"},
		Append:   defaultStyleAppend,
	}

	DraculaStyle = &Style{
		Key:      [2]string{"\x1B[38;5;81m", "\x1B[0m"},  // cyan
		String:   [2]string{"\x1B[38;5;114m", "\x1B[0m"}, // green
		Number:   [2]string{"\x1B[38;5;212m", "\x1B[0m"}, // pink
		True:     [2]string{"\x1B[38;5;203m", "\x1B[0m"}, // red
		False:    [2]string{"\x1B[38;5;203m", "\x1B[0m"},
		Null:     [2]string{"\x1B[38;5;244m", "\x1B[0m"},
		Escape:   [2]string{"\x1B[38;5;176m", "\x1B[0m"}, // purple
		Brackets: [2]string{"\x1B[38;5;231m", "\x1B[0m"}, // white
		Append:   defaultStyleAppend,
	}

	MonokaiStyle = &Style{
		Key:      [2]string{"\x1B[38;5;81m", "\x1B[0m"},  // blue
		String:   [2]string{"\x1B[38;5;186m", "\x1B[0m"}, // yellow
		Number:   [2]string{"\x1B[38;5;208m", "\x1B[0m"}, // orange
		True:     [2]string{"\x1B[38;5;197m", "\x1B[0m"}, // pink
		False:    [2]string{"\x1B[38;5;197m", "\x1B[0m"},
		Null:     [2]string{"\x1B[38;5;244m", "\x1B[0m"},
		Escape:   [2]string{"\x1B[38;5;141m", "\x1B[0m"}, // purple
		Brackets: [2]string{"\x1B[38;5;252m", "\x1B[0m"},
		Append:   defaultStyleAppend,
	}

	SolarizedDarkStyle = &Style{
		Key:      [2]string{"\x1B[38;5;33m", "\x1B[0m"},  // blue
		String:   [2]string{"\x1B[38;5;64m", "\x1B[0m"},  // green
		Number:   [2]string{"\x1B[38;5;166m", "\x1B[0m"}, // orange
		True:     [2]string{"\x1B[38;5;37m", "\x1B[0m"},  // cyan
		False:    [2]string{"\x1B[38;5;37m", "\x1B[0m"},
		Null:     [2]string{"\x1B[38;5;244m", "\x1B[0m"},
		Escape:   [2]string{"\x1B[38;5;125m", "\x1B[0m"}, // violet
		Brackets: [2]string{"\x1B[38;5;250m", "\x1B[0m"},
		Append:   defaultStyleAppend,
	}

	MinimalGrayStyle = &Style{
		Key:      [2]string{"\x1B[37m", "\x1B[0m"},
		String:   [2]string{"\x1B[90m", "\x1B[0m"},
		Number:   [2]string{"\x1B[37m", "\x1B[0m"},
		True:     [2]string{"\x1B[37m", "\x1B[0m"},
		False:    [2]string{"\x1B[37m", "\x1B[0m"},
		Null:     [2]string{"\x1B[2m", "\x1B[0m"},
		Escape:   [2]string{"\x1B[90m", "\x1B[0m"},
		Brackets: [2]string{"\x1B[37m", "\x1B[0m"},
		Append:   defaultStyleAppend,
	}
}

// Pretty takes a JSON byte slice and returns a pretty-printed version of the JSON.
// It uses the default configuration options specified in DefaultOptionsConfig.
//
// Parameters:
//   - json: The JSON data to be pretty-printed.
//
// Returns:
//   - A byte slice containing the pretty-printed JSON data.
func Pretty(json []byte) []byte { return PrettyOptions(json, nil) }

// PrettyOptions takes a JSON byte slice and returns a pretty-printed version of the JSON,
// with customizable options specified by the `option` parameter.
//
// Parameters:
//   - json: The JSON data to be pretty-printed.
//   - option: A pointer to an OptionsConfig struct containing custom options for pretty-printing.
//     If nil, the default options (DefaultOptionsConfig) will be used.
//
// Returns:
//   - A byte slice containing the pretty-printed JSON data based on the specified options.
//
// Notes:
//   - If `option` is nil, it falls back to the default configuration (DefaultOptionsConfig).
//   - The `appendPrettyAny` function is called to format the JSON with the provided options.
//
// PrettyOptions is like Pretty but with customized options.
func PrettyOptions(json []byte, option *OptionsConfig) []byte {
	if option == nil {
		option = DefaultOptionsConfig
	}
	buf := make([]byte, 0, len(json))
	if len(option.Prefix) != 0 {
		buf = append(buf, option.Prefix...)
	}
	buf, _, _, _ = appendPrettyAny(buf, json, 0, true,
		option.Width, option.Prefix, option.Indent, option.SortKeys,
		0, 0, -1)
	if len(buf) > 0 {
		buf = append(buf, '\n')
	}
	return buf
}

// Ugly removes unwanted characters from a JSON byte slice, returning a cleaned-up copy.
//
// This function creates a new byte buffer to hold a "cleaned" version of the input JSON, then calls the
// `ugly` function to process the input `json` byte slice. The `ugly` function filters out non-printable
// characters (characters with ASCII values less than or equal to ' '), preserving quoted substrings within
// the JSON. Unlike `UglyInPlace`, this function does not modify the original input and instead returns a
// new byte slice.
//
// Parameters:
//   - `json`: A byte slice containing JSON data that may include unwanted characters or quoted substrings.
//
// Returns:
//   - A new byte slice with unwanted characters removed. The cleaned-up version of the input `json`.
//
// Example:
//
//	json := []byte(`hello "world" 1234`)
//	cleanedJson := Ugly(json)
//	// cleanedJson will be []byte{'h', 'e', 'l', 'l', 'o', ' ', '"', 'w', 'o', 'r', 'l', 'd', '"', ' ', '1', '2', '3', '4'}
//	// as the function preserves printable characters and properly handles quoted substrings.
//
// Notes:
//   - This function is useful when you need a cleaned copy of the original JSON data without modifying the original byte slice.
//   - The buffer created (`buf`) is pre-allocated with a capacity equal to the length of the input, optimizing memory allocation.
func Ugly(json []byte) []byte {
	buf := make([]byte, 0, len(json))
	return ugly(buf, json)
}

// UglyInPlace removes unwanted characters from a JSON byte slice in-place and returns the modified byte slice.
//
// This function is a wrapper around the `ugly` function, which processes a byte slice and removes non-printable characters
// (i.e., characters with ASCII values less than or equal to ' '), preserving quoted substrings. `UglyInPlace` calls `ugly`
// with the same slice as both the source (`src`) and the destination (`dst`), effectively performing the cleaning operation
// in place.
//
// Parameters:
//   - `json`: A byte slice containing the JSON data to clean up. It may include unwanted characters or quoted substrings.
//
// Returns:
//   - The modified `json` byte slice with unwanted characters removed.
//
// Example:
//
//	json := []byte(`hello "world" 1234`)
//	cleanedJson := UglyInPlace(json)
//	// cleanedJson will be []byte{'h', 'e', 'l', 'l', 'o', ' ', '"', 'w', 'o', 'r', 'l', 'd', '"', ' ', '1', '2', '3', '4'}
//	// as the function preserves printable characters and quoted substrings.
//
// Notes:
//   - This function is intended for cases where in-place modification of the input is acceptable.
//   - The underlying `ugly` function processes each character, handling escaped double quotes to avoid breaking quoted substrings.
func UglyInPlace(json []byte) []byte { return ugly(json, json) }

// Spec strips out comments and trailing commas and converts the input to a valid JSON format
// according to the official JSON specification (RFC 8259).
//
// This function calls the `spec` helper function to process the input `source` byte slice. It
// removes all single-line (`//`) and multi-line (`/* */`) comments, as well as any trailing commas
// that might be present in the input JSON. The result is a valid, parsable JSON byte slice that
// conforms to the RFC 8259 standard.
//
// Key characteristics of the function:
//   - The output will have the same length as the input source byte slice.
//   - All original line breaks (newlines) will be preserved at the same positions in the output.
//   - The function ensures that the cleaned JSON remains structurally valid and compliant with the
//     official specification, making it ready for parsing by external JSON parsers.
//
// This function is useful for scenarios where you need to preprocessed JSON-like data, removing
// comments and trailing commas while maintaining the correct formatting and offsets for later
// parsing and error reporting.
//
// Parameters:
//   - `source`: The input byte slice containing the raw JSON-like data, which may include
//     comments and trailing commas.
//
// Returns:
//   - A new byte slice containing the cleaned, valid JSON data with comments and trailing
//     commas removed, while preserving original formatting and line breaks.
//
// Example usage:
//
//	rawJSON := []byte(`{ // comment\n "key": "value", }`)
//	validJSON := Spec(rawJSON)
//	// validJSON will be cleaned and ready for parsing.
func Spec(source []byte) []byte {
	return spec(source, nil)
}

// SpecInPlace strips out comments and trailing commas from the input JSON-like data
// and modifies the input slice in-place, converting it to valid JSON format according to
// the official JSON specification (RFC 8259).
//
// This function behaves similarly to the `Spec` function, but instead of returning a new
// byte slice with the cleaned JSON, it modifies the original `source` byte slice directly.
//
// It removes all single-line (`//`) and multi-line (`/* */`) comments, as well as any trailing commas
// that might be present in the input. The result is a valid, parsable JSON byte slice that
// adheres to the RFC 8259 standard.
//
// Key characteristics of the function:
//   - The output is stored directly in the `source` byte slice, modifying it in place.
//   - The function ensures that the cleaned JSON remains structurally valid and compliant with the
//     official specification, making it ready for parsing by external JSON parsers.
//
// This function is useful when you want to modify the input data directly, avoiding the need
// for creating a new byte slice. It ensures that the original slice is cleaned while maintaining
// the correct formatting and line breaks.
//
// Parameters:
//   - `source`: The input byte slice containing the raw JSON-like data, which may include
//     comments and trailing commas. This slice will be modified in place.
//
// Returns:
//   - The same `source` byte slice, now containing the cleaned, valid JSON data with comments
//     and trailing commas removed, while preserving original formatting and line breaks.
//
// Example usage:
//
//	rawJSON := []byte(`{ // comment\n "key": "value", }`)
//	SpecInPlace(rawJSON)
//	// rawJSON will be cleaned in-place and ready for parsing.
func SpecInPlace(source []byte) []byte {
	return spec(source, source)
}

// Color takes a JSON source in the form of a byte slice and applies syntax highlighting based on the provided style.
// The function returns a new byte slice with the JSON source formatted according to the specified styles.
//
// Parameters:
//   - `source`: A byte slice containing the JSON content to be styled.
//   - `style`: A pointer to a `Style` struct that defines the styling for various components of the JSON content.
//     If `nil`, the function uses the default `TerminalStyle`.
//
// Returns:
//   - A byte slice with the styled JSON content. Each component (such as keys, values, numbers, etc.) is colored based on the provided style.
//
// The function processes the JSON source character by character and applies the corresponding styles as follows:
//   - String values are enclosed in double quotes and are styled using the `Key` and `String` style fields.
//   - Numbers, booleans (`true`, `false`), and `null` values are styled using the `Number`, `True`, `False`, and `Null` fields.
//   - Brackets (`{`, `}`, `[`, `]`) are styled using the `Brackets` field.
//   - Escape sequences within strings are styled using the `Escape` field.
//   - The function handles nested objects and arrays using a stack to track the current level of nesting.
//
// Example:
//
//	source := []byte(`{"name": "John", "age": 30, "active": true}`)
//	style := &Style{
//	  Key:   [2]string{"\033[1;34m", "\033[0m"},
//	  String: [2]string{"\033[1;32m", "\033[0m"},
//	  Number: [2]string{"\033[1;33m", "\033[0m"},
//	  True: [2]string{"\033[1;35m", "\033[0m"},
//	  False: [2]string{"\033[1;35m", "\033[0m"},
//	  Null: [2]string{"\033[1;35m", "\033[0m"},
//	  Escape: [2]string{"\033[1;31m", "\033[0m"},
//	  Brackets: [2]string{"\033[1;37m", "\033[0m"},
//	  Append: func(dst []byte, c byte) []byte { return append(dst, c) },
//	}
//	result := Color(source, style)
//	fmt.Println(string(result)) // Prints the styled JSON
//
// Notes:
//   - The function handles escape sequences (e.g., `\n`, `\"`) and ensures they are properly colored as part of strings.
//   - The `Append` function in the `Style` struct allows customization of how each character is appended, enabling more flexible formatting if needed.
func Color(source []byte, style *Style) []byte {
	if style == nil {
		style = TerminalStyle
	}
	appendStyle := style.Append
	if appendStyle == nil {
		appendStyle = func(dst []byte, c byte) []byte {
			return append(dst, c)
		}
	}
	type innerStack struct {
		kind byte
		key  bool
	}
	var destinationByte []byte
	var stack []innerStack
	for i := 0; i < len(source); i++ {
		if source[i] == '"' {
			key := len(stack) > 0 && stack[len(stack)-1].key
			if key {
				destinationByte = append(destinationByte, style.Key[0]...)
			} else {
				destinationByte = append(destinationByte, style.String[0]...)
			}
			destinationByte = appendStyle(destinationByte, '"')
			esc := false
			useEsc := 0
			for i = i + 1; i < len(source); i++ {
				if source[i] == '\\' {
					if key {
						destinationByte = append(destinationByte, style.Key[1]...)
					} else {
						destinationByte = append(destinationByte, style.String[1]...)
					}
					destinationByte = append(destinationByte, style.Escape[0]...)
					destinationByte = appendStyle(destinationByte, source[i])
					esc = true
					if i+1 < len(source) && source[i+1] == 'u' {
						useEsc = 5
					} else {
						useEsc = 1
					}
				} else if esc {
					destinationByte = appendStyle(destinationByte, source[i])
					if useEsc == 1 {
						esc = false
						destinationByte = append(destinationByte, style.Escape[1]...)
						if key {
							destinationByte = append(destinationByte, style.Key[0]...)
						} else {
							destinationByte = append(destinationByte, style.String[0]...)
						}
					} else {
						useEsc--
					}
				} else {
					destinationByte = appendStyle(destinationByte, source[i])
				}
				if source[i] == '"' {
					j := i - 1
					for ; ; j-- {
						if source[j] != '\\' {
							break
						}
					}
					if (j-i)%2 != 0 {
						break
					}
				}
			}
			if esc {
				destinationByte = append(destinationByte, style.Escape[1]...)
			} else if key {
				destinationByte = append(destinationByte, style.Key[1]...)
			} else {
				destinationByte = append(destinationByte, style.String[1]...)
			}
		} else if source[i] == '{' || source[i] == '[' {
			stack = append(stack, innerStack{source[i], source[i] == '{'})
			destinationByte = append(destinationByte, style.Brackets[0]...)
			destinationByte = appendStyle(destinationByte, source[i])
			destinationByte = append(destinationByte, style.Brackets[1]...)
		} else if (source[i] == '}' || source[i] == ']') && len(stack) > 0 {
			stack = stack[:len(stack)-1]
			destinationByte = append(destinationByte, style.Brackets[0]...)
			destinationByte = appendStyle(destinationByte, source[i])
			destinationByte = append(destinationByte, style.Brackets[1]...)
		} else if (source[i] == ':' || source[i] == ',') && len(stack) > 0 && stack[len(stack)-1].kind == '{' {
			stack[len(stack)-1].key = !stack[len(stack)-1].key
			destinationByte = append(destinationByte, style.Brackets[0]...)
			destinationByte = appendStyle(destinationByte, source[i])
			destinationByte = append(destinationByte, style.Brackets[1]...)
		} else {
			var kind byte
			if (source[i] >= '0' && source[i] <= '9') || source[i] == '-' || isNaNOrInf(source[i:]) {
				kind = '0'
				destinationByte = append(destinationByte, style.Number[0]...)
			} else if source[i] == 't' {
				kind = 't'
				destinationByte = append(destinationByte, style.True[0]...)
			} else if source[i] == 'f' {
				kind = 'f'
				destinationByte = append(destinationByte, style.False[0]...)
			} else if source[i] == 'n' {
				kind = 'n'
				destinationByte = append(destinationByte, style.Null[0]...)
			} else {
				destinationByte = appendStyle(destinationByte, source[i])
			}
			if kind != 0 {
				for ; i < len(source); i++ {
					if source[i] <= ' ' || source[i] == ',' || source[i] == ':' || source[i] == ']' || source[i] == '}' {
						i--
						break
					}
					destinationByte = appendStyle(destinationByte, source[i])
				}
				if kind == '0' {
					destinationByte = append(destinationByte, style.Number[1]...)
				} else if kind == 't' {
					destinationByte = append(destinationByte, style.True[1]...)
				} else if kind == 'f' {
					destinationByte = append(destinationByte, style.False[1]...)
				} else if kind == 'n' {
					destinationByte = append(destinationByte, style.Null[1]...)
				}
			}
		}
	}
	return destinationByte
}

// Len returns the number of key-value pairs in the byKeyVal struct.
// It is part of the sort.Interface implementation, allowing pairs to be sorted.
func (b *kvSorter) Len() int {
	return len(b.pairs)
}

// Less compares two pairs at indices i and j to determine if the pair at i should come before the pair at j.
// It first compares by key, and if the keys are equal, it compares by value.
func (b *kvSorter) Less(i, j int) bool {
	if b.isLess(i, j, sKey) {
		return true
	}
	if b.isLess(j, i, sKey) {
		return false
	}
	return b.isLess(i, j, sVal)
}

// Swap exchanges the positions of the pairs at indices i and j in the pairs slice.
// It also sets the sorted flag to true.
func (b *kvSorter) Swap(i, j int) {
	b.pairs[i], b.pairs[j] = b.pairs[j], b.pairs[i]
	b.sorted = true
}

// isLess compares two pairs by a specified criterion (key or value) and determines if one is less than the other.
// It trims whitespace from values and further processes them if they are strings or numbers.
func (a *kvSorter) isLess(i, j int, kind sortCriteria) bool {
	k1 := a.json[a.pairs[i].keyStart:a.pairs[i].keyEnd]
	k2 := a.json[a.pairs[j].keyStart:a.pairs[j].keyEnd]
	var v1, v2 []byte
	if kind == sKey {
		v1 = k1
		v2 = k2
	} else {
		v1 = bytes.TrimSpace(a.buf[a.pairs[i].valueStart:a.pairs[i].valueEnd])
		v2 = bytes.TrimSpace(a.buf[a.pairs[j].valueStart:a.pairs[j].valueEnd])
		if len(v1) >= len(k1)+1 {
			v1 = bytes.TrimSpace(v1[len(k1)+1:])
		}
		if len(v2) >= len(k2)+1 {
			v2 = bytes.TrimSpace(v2[len(k2)+1:])
		}
	}
	t1 := getJsonType(v1)
	t2 := getJsonType(v2)
	if t1 < t2 {
		return true
	}
	if t1 > t2 {
		return false
	}
	if t1 == jsonString {
		s1 := unescapeJSONString(v1)
		s2 := unescapeJSONString(v2)
		return string(s1) < string(s2)
	}
	if t1 == jNumber {
		n1, _ := strconv.ParseFloat(string(v1), 64)
		n2, _ := strconv.ParseFloat(string(v2), 64)
		return n1 < n2
	}
	return string(v1) < string(v2)
}

// ugly processes a source byte slice and removes unwanted characters, returning a new cleaned-up byte slice.
//
// This function processes the input `src` byte slice and appends characters to the `dst` byte slice based on certain criteria.
// It specifically filters out characters that are not printable (i.e., characters with ASCII values greater than `' '`).
// Additionally, it handles quoted substrings, ensuring that characters inside properly escaped double quotes are preserved.
// If an unescaped double quote is encountered, it will stop processing further characters.
//
// Parameters:
//   - `dst`: The destination slice of bytes where the cleaned characters will be appended.
//   - `src`: The source slice of bytes to process, which may contain unwanted characters and quoted substrings.
//
// Returns:
//   - A new byte slice (`dst`) with unwanted characters removed. The cleaned-up version of `src`.
//
// Example:
//
//	src := []byte(`hello "world" 1234`)
//	dst := ugly([]byte{}, src)
//	// dst will be []byte{'h', 'e', 'l', 'l', 'o', ' ', '"', 'w', 'o', 'r', 'l', 'd', '"', ' ', '1', '2', '3', '4'},
//	// as the function preserves only printable characters and properly handles quoted substrings.
//
// Notes:
//   - This function skips characters that are not printable (less than or equal to ASCII ' ').
//   - When encountering a double quote (`"`), the function ensures that it correctly handles escaped quotes, skipping characters
//     until a valid closing quote is found. If an odd number of backslashes precede the closing quote, it breaks the loop to avoid
//     incorrect parsing of the quotes.
func ugly(dst, src []byte) []byte {
	dst = dst[:0] // Reset destination slice to an empty state
	for i := 0; i < len(src); i++ {
		if src[i] > ' ' { // Only include characters that are printable
			dst = append(dst, src[i])
			if src[i] == '"' { // Handle quoted substring (double quotes)
				for i = i + 1; i < len(src); i++ {
					dst = append(dst, src[i])
					if src[i] == '"' {
						// Search backwards for the last non-escaped backslash
						j := i - 1
						for ; ; j-- {
							if src[j] != '\\' {
								break
							}
						}
						// If the number of consecutive backslashes is odd, break the loop
						if (j-i)%2 != 0 {
							break
						}
					}
				}
			}
		}
	}
	return dst
}

// isNaNOrInf checks if a byte slice represents a special numeric value: NaN (Not-a-Number) or Infinity.
//
// This function inspects the first character of the input byte slice to determine if it represents
// either NaN or Infinity, including variations such as `Inf`, `+Inf`, `inf`, `NaN`, and `nan`.
// The function returns `true` if the input matches any of these special values.
//
// Parameters:
//   - `src`: A byte slice to inspect, typically a numeric string.
//
// Returns:
//   - `true` if the byte slice represents NaN or Infinity, `false` otherwise.
//
// Example:
//
//	src1 := []byte("Inf")
//	src2 := []byte("NaN")
//	src3 := []byte("+Inf")
//	src4 := []byte("infinity")
//	result1 := isNaNOrInf(src1) // result1 will be true
//	result2 := isNaNOrInf(src2) // result2 will be true
//	result3 := isNaNOrInf(src3) // result3 will be true
//	result4 := isNaNOrInf(src4) // result4 will be false
//
// Notes:
//   - The function only inspects the first character (or first two characters for lowercase `nan`) to make a quick determination.
//   - It supports the variations `Inf`, `+Inf`, `inf`, `NaN`, and `nan` as valid representations.
func isNaNOrInf(src []byte) bool {
	return src[0] == 'i' || //Inf
		src[0] == 'I' || // inf
		src[0] == '+' || // +Inf
		src[0] == 'N' || // Nan
		(src[0] == 'n' && len(src) > 1 && src[1] != 'u') // nan
}

// getJsonType identifies the JSON type of a given byte slice based on its first character.
//
// This function analyzes the first character in the byte slice `v` to determine which JSON data type it represents.
// Based on the initial character, it categorizes the input as one of the following types:
// `jsonNull`, `jsonFalse`, `jsonTrue`, `jsonString`, `jNumber`, or `jsonJson` (indicating either a JSON object or array).
//
// Parameters:
//   - `v`: A byte slice representing a JSON value.
//
// Returns:
//   - A `jsonType` value that represents the JSON type of `v`, based on its first character.
//
// Example:
//
//	value1 := []byte(`"hello"`)
//	value2 := []byte("false")
//	value3 := []byte("123")
//	value4 := []byte("null")
//	value5 := []byte("[1, 2, 3]")
//	result1 := getJsonType(value1) // result1 will be jsonString
//	result2 := getJsonType(value2) // result2 will be jsonFalse
//	result3 := getJsonType(value3) // result3 will be jNumber
//	result4 := getJsonType(value4) // result4 will be jsonNull
//	result5 := getJsonType(value5) // result5 will be jsonJson
//
// Notes:
//   - If the byte slice is empty, the function returns `jsonNull`.
//   - The function uses the initial character of `v` to distinguish types, assuming `true`, `false`, and `null` are valid JSON values.
func getJsonType(v []byte) jsonType {
	if len(v) == 0 {
		return jsonNull
	}
	switch v[0] {
	case '"':
		return jsonString
	case 'f':
		return jsonFalse
	case 't':
		return jsonTrue
	case 'n':
		return jsonNull
	case '[', '{':
		return jsonJson
	default:
		return jNumber
	}
}

// unescapeJSONString extracts a JSON string from a byte slice, handling escaped characters when present.
//
// This function takes a JSON byte slice representing a string value and extracts the unescaped content.
// It iterates through the byte slice, checking for either escaped characters or the closing double quote (`"`).
// If an escape character (`\`) is detected, `unescapeJSONString` uses JSON unmarshalling to handle any escape sequences.
// If no escape character is encountered, it returns the substring between the opening and closing quotes.
//
// Parameters:
//   - `s`: A byte slice containing a JSON string, including the enclosing double quotes and potentially escaped characters.
//
// Returns:
//   - A byte slice with the unescaped content of the JSON string if valid, or `nil` if an error occurs.
//
// Example:
//
//	s := []byte(`"Hello, world!"`)
//	result := unescapeJSONString(s)
//	// result will be []byte{'H', 'e', 'l', 'l', 'o', ',', ' ', 'w', 'o', 'r', 'l', 'd', '!'}
//	s := []byte(`"Line1\nLine2"`)
//	result := unescapeJSONString(s)
//	// result will be []byte{'L', 'i', 'n', 'e', '1', '\n', 'L', 'i', 'n', 'e', '2'}
//
// Notes:
//   - If an escape sequence is encountered (`\`), JSON unmarshalling is used to correctly interpret it.
//   - If the byte slice does not contain a properly closed string, the function returns `nil`.
func unescapeJSONString(s []byte) []byte {
	for i := 1; i < len(s); i++ {
		if s[i] == '\\' {
			var str string
			json.Unmarshal(s, &str)
			return []byte(str)
		}
		if s[i] == '"' {
			return s[1:i]
		}
	}
	return nil
}

// sortPairs sorts JSON key-value pairs in a stable order, preserving original formatting,
// and returns the updated buffer containing the sorted pairs.
//
// This function takes the JSON data, a buffer to store formatted values, and a list of key-value pairs (`pairs`).
// If there are no pairs to sort, it directly returns the buffer. Otherwise, it initializes a `byKeyVal` struct
// with the JSON data, buffer, and pairs, and sorts them by key (and by value if keys are identical) using
// `sort.Stable`. After sorting, it constructs a new byte slice with the sorted pairs in order, each followed
// by a comma and newline, and replaces the original content in `buf`.
//
// Parameters:
//   - `json`: The original JSON data as a byte slice.
//   - `buf`: A byte slice that holds formatted key-value pairs and can be modified in-place.
//   - `pairs`: A slice of `pair` structs representing the key-value pairs in the JSON data.
//
// Returns:
//   - A byte slice containing the buffer with sorted key-value pairs, in stable order.
//
// Example:
//
//	json := []byte(`{"b":2, "a":1}`)
//	pairs := []pair{ /* initialized with positions of keys and values */ }
//	buf := make([]byte, len(json))
//	result := sortPairs(json, buf, pairs)
//	// result will contain sorted pairs by key.
//
// Notes:
//   - If `pairs` is empty, `buf` is returned unchanged.
//   - `sort.Stable` is used to ensure that pairs with identical keys maintain their original relative order.
//   - If `byKeyVal` marks pairs as unsorted, it skips replacing the original buffer.
func sortPairs(json, buf []byte, pairs []kvPairs) []byte {
	if len(pairs) == 0 {
		return buf
	}
	_valStart := pairs[0].valueStart
	_valEnd := pairs[len(pairs)-1].valueEnd
	_keyVal := kvSorter{false, json, buf, pairs}
	sort.Stable(&_keyVal)
	if !_keyVal.sorted {
		return buf
	}
	n := make([]byte, 0, _valEnd-_valStart)
	for i, p := range pairs {
		n = append(n, buf[p.valueStart:p.valueEnd]...)
		if i < len(pairs)-1 {
			n = append(n, ',')
			n = append(n, '\n')
		}
	}
	return append(buf[:_valStart], n...)
}

// appendPrettyString appends a JSON string value from the input JSON byte slice (`json`) to a buffer (`buf`),
// handling any escaped characters within the string, and returns the updated buffer and indices.
//
// This function begins at a given index `i` within a JSON byte slice `json`, assuming the current character is
// the start of a JSON string (`"`). It appends the entire string (from opening to closing quote) to the buffer `buf`,
// handling escaped quotes within the string. If an escape sequence (`\`) precedes a closing quote, it continues searching
// until it finds an unescaped closing quote, marking the end of the string.
//
// Parameters:
//   - `buf`: The destination byte slice to which the JSON string value is appended.
//   - `json`: The source JSON byte slice containing the entire JSON structure.
//   - `i`: The starting index in `json`, pointing to the beginning of the JSON string (initial quote).
//   - `nl`: The current newline position (used for pretty-printing in larger context).
//
// Returns:
//   - `buf`: The updated buffer, containing the appended JSON string.
//   - `i`: The updated index in `json` after the end of the string (right after the closing quote).
//   - `nl`: The unchanged newline position (for tracking in pretty-printing).
//   - `true`: A boolean flag indicating the function processed a string (for handling in other contexts).
//
// Example:
//
//	json := []byte(`"example \"string\" value"`)
//	buf := []byte{}
//	buf, i, nl, processed := appendPrettyString(buf, json, 0, 0)
//	// buf will contain `example \"string\" value`, i will point to the next index after the closing quote,
//	// nl remains the same, and processed is true.
//
// Notes:
//   - The function counts consecutive backslashes before each closing quote to determine if it is escaped.
//   - It appends the entire string (including quotes) to `buf` for easy integration in pretty-printing or formatting routines.
func appendPrettyString(buf, json []byte, i, nl int) ([]byte, int, int, bool) {
	s := i
	i++
	for ; i < len(json); i++ {
		if json[i] == '"' {
			var sc int
			for j := i - 1; j > s; j-- {
				if json[j] == '\\' {
					sc++
				} else {
					break
				}
			}
			if sc%2 == 1 {
				continue
			}
			i++
			break
		}
	}
	return append(buf, json[s:i]...), i, nl, true
}

// appendPrettyNumber appends a JSON number value from the input JSON byte slice (`json`) to a buffer (`buf`),
// and returns the updated buffer and indices.
//
// This function starts at a given index `i` within a JSON byte slice `json` (assuming the current character is the start
// of a JSON number). It appends the entire number to the buffer `buf` and handles all characters up to the next
// non-number character, such as spaces, commas, colons, or closing brackets/braces.
//
// Parameters:
//   - `buf`: The destination byte slice to which the JSON number value will be appended.
//   - `json`: The source JSON byte slice containing the entire JSON structure.
//   - `i`: The starting index in `json`, pointing to the first character of the JSON number.
//   - `nl`: The current newline position (used for pretty-printing in larger context).
//
// Returns:
//   - `buf`: The updated buffer, containing the appended JSON number value.
//   - `i`: The updated index in `json` after the number (right after the last character of the number).
//   - `nl`: The unchanged newline position (for tracking in pretty-printing).
//   - `true`: A boolean flag indicating that a number was processed successfully.
//
// Example:
//
//	json := []byte(`12345`)
//	buf := []byte{}
//	buf, i, nl, processed := appendPrettyNumber(buf, json, 0, 0)
//	// buf will contain `12345`, i will point to the next index after the number,
//	// nl remains unchanged, and processed will be true.
//
// Notes:
//   - The function scans for all characters that are part of the number (digits, decimal point, etc.) until it
//     encounters a character that is not part of a valid number, such as a space, comma, colon, or closing bracket/braces.
//   - It assumes that the number is well-formed and does not handle error cases like invalid numbers.
func appendPrettyNumber(buf, json []byte, i, nl int) ([]byte, int, int, bool) {
	s := i // Record the start index of the number
	i++    // Move past the initial digit (or minus sign if present)
	for ; i < len(json); i++ {
		// Break the loop if a non-number character is encountered (e.g., space, comma, colon, bracket, or brace)
		if json[i] <= ' ' || json[i] == ',' || json[i] == ':' || json[i] == ']' || json[i] == '}' {
			break
		}
	}
	// Append the number from the start index `s` to the updated index `i` (excluding non-number characters)
	return append(buf, json[s:i]...), i, nl, true
}

// appendPrettyAny processes the next JSON value in the input JSON byte slice (`json`) and appends it to the buffer (`buf`),
// while handling different types of JSON values (strings, numbers, objects, arrays, and literals).
// It returns the updated buffer and indices, as well as a boolean flag indicating whether a value was processed.
//
// This function is responsible for recognizing the type of the next JSON value and delegating the task of appending that value
// to the appropriate helper function (such as `appendPrettyString` for strings, `appendPrettyNumber` for numbers, and others
// for objects, arrays, and literals). It processes the JSON byte slice one element at a time and handles all value types
// correctly, ensuring that each value is pretty-printed if required.
//
// Parameters:
//   - `buf`: The destination byte slice to which the processed JSON value will be appended.
//   - `json`: The source JSON byte slice containing the entire JSON structure.
//   - `i`: The starting index in `json` from where the next value should be processed.
//   - `pretty`: A boolean flag indicating whether pretty-printing should be applied (i.e., adding newlines and indentation).
//   - `width`: The width used for pretty-printing (not used in this function, but passed for consistency in pretty-printing logic).
//   - `prefix`: A string prefix (used in pretty-printing to add leading indentation, not used here).
//   - `indent`: The string used for indentation (not used here but part of the pretty-printing configuration).
//   - `sortKeys`: A boolean flag indicating whether the keys in objects should be sorted (not used here but passed for consistency).
//   - `tabs`: The number of tabs for indentation (not used here but passed for consistency in pretty-printing logic).
//   - `nl`: The current newline position (used for pretty-printing, ensuring line breaks are maintained correctly).
//   - `max`: The maximum number of characters to pretty-print before breaking into a new line (not used in this function).
//
// Returns:
//   - `buf`: The updated buffer, containing the appended JSON value (pretty-printed if the `pretty` flag is true).
//   - `i`: The updated index in `json`, pointing to the position after the processed value.
//   - `nl`: The unchanged newline position (used for pretty-printing in the larger context).
//   - `true`: A boolean flag indicating that a JSON value was successfully processed.
//
// Example usage:
//
//	json := []byte(`{ "key1": 123, "key2": "value", "key3": [1, 2, 3] }`)
//	buf := []byte{}
//	buf, i, nl, processed := appendPrettyAny(buf, json, 0, true, 0, "", "  ", false, 0, 0, 0)
//	// buf will contain the pretty-printed JSON value, i will point to the next index,
//	// nl remains unchanged, and processed will be true.
//
// Notes:
//   - This function processes and appends various JSON data types, including strings, numbers, objects, arrays,
//     and literals (`true`, `false`, `null`).
//   - The function assumes the JSON is valid and well-formed; it does not handle parsing errors for invalid JSON.
func appendPrettyAny(buf, json []byte, i int, pretty bool, width int, prefix, indent string, sortKeys bool, tabs, nl, max int) ([]byte, int, int, bool) {
	for ; i < len(json); i++ {
		if json[i] <= ' ' {
			continue
		}
		if json[i] == '"' {
			return appendPrettyString(buf, json, i, nl)
		}
		if (json[i] >= '0' && json[i] <= '9') || json[i] == '-' || isNaNOrInf(json[i:]) {
			return appendPrettyNumber(buf, json, i, nl)
		}
		if json[i] == '{' {
			return appendPrettyObject(buf, json, i, '{', '}', pretty, width, prefix, indent, sortKeys, tabs, nl, max)
		}
		if json[i] == '[' {
			return appendPrettyObject(buf, json, i, '[', ']', pretty, width, prefix, indent, sortKeys, tabs, nl, max)
		}
		switch json[i] {
		case 't':
			return append(buf, 't', 'r', 'u', 'e'), i + 4, nl, true
		case 'f':
			return append(buf, 'f', 'a', 'l', 's', 'e'), i + 5, nl, true
		case 'n':
			return append(buf, 'n', 'u', 'l', 'l'), i + 4, nl, true
		}
	}
	return buf, i, nl, true
}

// appendPrettyObject processes the next JSON object or array in the input JSON byte slice (`json`)
// and appends it to the buffer (`buf`), while handling pretty-printing, sorting of object keys, and enforcing width constraints.
//
// This function handles the parsing and formatting of JSON objects (`{}`) and arrays (`[]`), adding appropriate indentation,
// newlines, and sorting of keys (if specified). It ensures that the resulting object or array is correctly formatted and
// inserted into the buffer, maintaining the structure and respecting pretty-printing preferences.
//
// It also handles the optional width limit (to control the length of single-line arrays) and can process objects with sorted keys,
// ensuring proper formatting of both simple and complex JSON objects.
//
// Parameters:
//   - `buf`: The destination byte slice to which the processed JSON object or array will be appended.
//   - `json`: The source JSON byte slice containing the entire JSON structure.
//   - `i`: The starting index in `json` from where the next object or array should be processed.
//   - `open`: The opening byte (either '{' for an object or '[' for an array).
//   - `close`: The closing byte (either '}' for an object or ']' for an array).
//   - `pretty`: A boolean flag indicating whether pretty-printing should be applied (i.e., adding newlines and indentation).
//   - `width`: The width used for pretty-printing (influences line breaks for large arrays).
//   - `prefix`: A string prefix used for leading indentation (used in pretty-printing).
//   - `indent`: The string used for indentation in pretty-printing.
//   - `sortKeys`: A boolean flag indicating whether the keys in objects should be sorted.
//   - `tabs`: The number of tabs for indentation (used in pretty-printing).
//   - `nl`: The current newline position, used for managing where to insert newlines during pretty-printing.
//   - `max`: The maximum number of characters to pretty-print before breaking into a new line (relevant for width-based formatting).
//
// Returns:
//   - `buf`: The updated buffer containing the appended JSON object or array (pretty-printed if the `pretty` flag is true).
//   - `i`: The updated index in `json`, pointing to the position after the processed object or array.
//   - `nl`: The updated newline position, adjusted for pretty-printing.
//   - `true` or `false`: A boolean flag indicating whether the processing was successful. The function returns `false` if
//     there is an issue (for example, exceeding the width limit or malformed data).
//
// Example usage:
//
//	json := []byte(`{ "key1": 123, "key2": "value", "key3": [1, 2, 3] }`)
//	buf := []byte{}
//	buf, i, nl, processed := appendPrettyObject(buf, json, 0, '{', '}', true, 80, "", "  ", false, 0, 0, -1)
//	// buf will contain the pretty-printed JSON object, i will point to the next index,
//	// nl will be adjusted for newlines, and processed will be true.
//
// Notes:
//   - This function handles both JSON objects and arrays, depending on the value of `open` and `close` (either '{', '}' for objects
//     or '[' and ']' for arrays).
//   - Pretty-printing is applied based on the `pretty` flag, including indentation and line breaks.
//   - If `sortKeys` is set to true, the keys in the JSON object will be sorted lexicographically before being appended to the buffer.
//   - The `max` value helps control the number of characters in a single line for arrays, ensuring that arrays are properly wrapped into multiple lines if necessary.
//   - The function can also handle arrays and objects with nested structures, ensuring the formatting remains correct throughout.
func appendPrettyObject(buf, json []byte, i int, open, close byte, pretty bool, width int, prefix, indent string, sortKeys bool, tabs, nl, max int) ([]byte, int, int, bool) {
	var ok bool
	if width > 0 {
		if pretty && open == '[' && max == -1 {
			// here we try to create a single line array
			max := width - (len(buf) - nl)
			if max > 3 {
				s1, s2 := len(buf), i
				buf, i, _, ok = appendPrettyObject(buf, json, i, '[', ']', false, width, prefix, "", sortKeys, 0, 0, max)
				if ok && len(buf)-s1 <= max {
					return buf, i, nl, true
				}
				buf = buf[:s1]
				i = s2
			}
		} else if max != -1 && open == '{' {
			return buf, i, nl, false
		}
	}
	buf = append(buf, open)
	i++
	var pairs []kvPairs
	if open == '{' && sortKeys {
		pairs = make([]kvPairs, 0, 8)
	}
	var n int
	for ; i < len(json); i++ {
		if json[i] <= ' ' {
			continue
		}
		if json[i] == close {
			if pretty {
				if open == '{' && sortKeys {
					buf = sortPairs(json, buf, pairs)
				}
				if n > 0 {
					nl = len(buf)
					if buf[nl-1] == ' ' {
						buf[nl-1] = '\n'
					} else {
						buf = append(buf, '\n')
					}
				}
				if buf[len(buf)-1] != open {
					buf = appendTabs(buf, prefix, indent, tabs)
				}
			}
			buf = append(buf, close)
			return buf, i + 1, nl, open != '{'
		}
		if open == '[' || json[i] == '"' {
			if n > 0 {
				buf = append(buf, ',')
				if width != -1 && open == '[' {
					buf = append(buf, ' ')
				}
			}
			var p kvPairs
			if pretty {
				nl = len(buf)
				if buf[nl-1] == ' ' {
					buf[nl-1] = '\n'
				} else {
					buf = append(buf, '\n')
				}
				if open == '{' && sortKeys {
					p.keyStart = i
					p.valueStart = len(buf)
				}
				buf = appendTabs(buf, prefix, indent, tabs+1)
			}
			if open == '{' {
				buf, i, nl, _ = appendPrettyString(buf, json, i, nl)
				if sortKeys {
					p.keyEnd = i
				}
				buf = append(buf, ':')
				if pretty {
					buf = append(buf, ' ')
				}
			}
			buf, i, nl, ok = appendPrettyAny(buf, json, i, pretty, width, prefix, indent, sortKeys, tabs+1, nl, max)
			if max != -1 && !ok {
				return buf, i, nl, false
			}
			if pretty && open == '{' && sortKeys {
				p.valueEnd = len(buf)
				if p.keyStart > p.keyEnd || p.valueStart > p.valueEnd {
					// bad data. disable sorting
					sortKeys = false
				} else {
					pairs = append(pairs, p)
				}
			}
			i--
			n++
		}
	}
	return buf, i, nl, open != '{'
}

// appendTabs appends indentation to the provided buffer (`buf`) based on the specified `prefix`, `indent`,
// and the number of `tabs` to insert.
//
// This function adds a specific number of tab or space-based indents to the `buf`, depending on the `indent` value.
// If the `indent` string contains exactly two spaces, it will append spaces for each tab; otherwise, it uses
// the provided `indent` string (which can represent any indenting character, such as tabs or custom strings).
// Additionally, if a `prefix` is provided, it will be prepended to the buffer before any indentation is added.
//
// Parameters:
//   - `buf`: The byte slice to which the indentations are appended.
//   - `prefix`: A string (byte slice) that will be prepended to `buf` before any indentation (optional).
//   - `indent`: The string to use for indentation, typically consisting of spaces or tabs (e.g., `"\t"` or `"  "`).
//   - `tabs`: The number of times the `indent` should be repeated to represent the desired level of indentation.
//
// Returns:
//   - The updated `buf` with the appropriate amount of indentation based on `tabs` and `indent`.
//
// Example:
//
//	buf := []byte{}
//	prefix := "  "
//	indent := "\t"
//	tabs := 3
//	buf = appendTabs(buf, prefix, indent, tabs)
//	// buf will be `{"  "\t\t\t` (prefix followed by 3 tab characters).
//
// Notes:
//   - If the `indent` string is exactly two spaces (`"  "`), the function will append two spaces for each `tab`.
//   - If the `indent` string is anything else, it will be appended `tabs` times.
func appendTabs(buf []byte, prefix, indent string, tabs int) []byte {
	if len(prefix) != 0 { // Append prefix if it's not an empty string
		buf = append(buf, prefix...)
	}
	// Check if the indent is exactly two spaces and append spaces for each tab
	if len(indent) == 2 && indent[0] == ' ' && indent[1] == ' ' {
		for i := 0; i < tabs; i++ {
			buf = append(buf, ' ', ' ')
		}
	} else {
		// Otherwise, append the custom indent string for each tab
		for i := 0; i < tabs; i++ {
			buf = append(buf, indent...)
		}
	}
	return buf
}

// hexDigit converts a numeric value to its corresponding hexadecimal character.
//
// This function takes a single byte `p`, which represents a numeric value, and converts it to its
// hexadecimal equivalent as a byte. The function assumes that the input `p` is in the range of 0 to 15
// (i.e., it represents a single hexadecimal digit).
//
// Parameters:
//   - `p`: A byte representing a numeric value between 0 and 15 (inclusive).
//
// Returns:
//   - A byte representing the corresponding hexadecimal character.
//   - If `p` is less than 10, it returns the ASCII character for the corresponding digit (0-9).
//   - If `p` is 10 or greater, it returns the lowercase ASCII character for the corresponding letter (a-f).
//
// Example:
//
//	hexDigit(0)  // returns '0'
//	hexDigit(9)  // returns '9'
//	hexDigit(10) // returns 'a'
//	hexDigit(15) // returns 'f'
//
// Notes:
//   - The function works only for values between 0 and 15 (inclusive).
//   - For input values greater than 15, the behavior is not defined and may lead to unexpected results.
func hexDigit(p byte) byte {
	// If p is less than 10, return the corresponding digit character ('0' to '9')
	switch {
	case p < 10:
		return p + '0' // Add ASCII value of '0' to convert to character
	default:
		// If p is 10 or greater, return the corresponding letter character ('a' to 'f')
		// Add ASCII value of 'a' to get 'a' to 'f'
		return (p - 10) + 'a'
	}
}

// spec processes a source byte slice (source) and removes or replaces comment sections,
// formatting them into a cleaned-up destination byte slice (destination).
// It handles both single-line (`//`) and multi-line (`/* */`) comments and strips them out,
// replacing them with spaces or newlines as appropriate.
//
// Parameters:
//   - `source`: The input byte slice containing the source code to process.
//   - `destination`: The output byte slice that will contain the cleaned code (without comments).
//     It is assumed to be initialized as an empty slice.
//
// Returns:
//   - A byte slice representing the cleaned source code with comments removed or replaced.
//   - Single-line comments (`//`) are replaced with spaces.
//   - Multi-line comments (`/* */`) are replaced with spaces and preserved newlines for line breaks.
//
// Example:
//
//	source := []byte("int x = 10; // initialize x\n/* multi-line\n comment */")
//	destination := spec(source, []byte{})
//	// destination will contain: "int x = 10;    \n   "
//
// Notes:
//   - The function handles the following types of comments:
//   - Single-line comments starting with `//` and ending with the newline (`\n`).
//   - Multi-line comments enclosed in `/* */`, even if they span multiple lines.
//   - Strings inside double quotes (`"`) and special characters like `}` or `]` are preserved as is,
//     with careful handling of quotes and escape sequences within the string.
func spec(source, destination []byte) []byte {
	destination = destination[:0]
	for i := 0; i < len(source); i++ {
		if source[i] == '/' {
			if i < len(source)-1 {
				if source[i+1] == '/' {
					destination = append(destination, ' ', ' ')
					i += 2
					for ; i < len(source); i++ {
						if source[i] == '\n' {
							destination = append(destination, '\n')
							break
						} else if source[i] == '\t' || source[i] == '\r' {
							destination = append(destination, source[i])
						} else {
							destination = append(destination, ' ')
						}
					}
					continue
				}
				if source[i+1] == '*' {
					destination = append(destination, ' ', ' ')
					i += 2
					for ; i < len(source)-1; i++ {
						if source[i] == '*' && source[i+1] == '/' {
							destination = append(destination, ' ', ' ')
							i++
							break
						} else if source[i] == '\n' || source[i] == '\t' ||
							source[i] == '\r' {
							destination = append(destination, source[i])
						} else {
							destination = append(destination, ' ')
						}
					}
					continue
				}
			}
		}
		destination = append(destination, source[i])
		if source[i] == '"' {
			for i = i + 1; i < len(source); i++ {
				destination = append(destination, source[i])
				if source[i] == '"' {
					j := i - 1
					for ; ; j-- {
						if source[j] != '\\' {
							break
						}
					}
					if (j-i)%2 != 0 {
						break
					}
				}
			}
		} else if source[i] == '}' || source[i] == ']' {
			for j := len(destination) - 2; j >= 0; j-- {
				if destination[j] <= ' ' {
					continue
				}
				if destination[j] == ',' {
					destination[j] = ' '
				}
				break
			}
		}
	}
	return destination
}

// defaultStyleAppend appends a byte `c` to the destination byte slice `dst`,
// handling special control characters by escaping them in Unicode format.
// If `c` is a control character (ASCII value less than 32) and not one of the
// common whitespace characters (`\r`, `\n`, `\t`, `\v`), it appends
// the Unicode escape sequence `\u00XX` to `dst`, where `XX` is the hexadecimal
// representation of `c`. Otherwise, it appends `c` directly to `dst`.
func defaultStyleAppend(dst []byte, c byte) []byte {
	if c < ' ' && (c != '\r' && c != '\n' && c != '\t' && c != '\v') {
		dst = append(dst, "\\u00"...)
		dst = append(dst, hexDigit((c>>4)&0xF))
		return append(dst, hexDigit(c&0xF))
	}
	return append(dst, c)
}
