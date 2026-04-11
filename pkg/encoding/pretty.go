package encoding

import (
	"bytes"
	"strconv"
)

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
