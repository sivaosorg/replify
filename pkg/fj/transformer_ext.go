package fj

import (
	"fmt"
	"math"
	"strings"

	"github.com/sivaosorg/replify/pkg/strutil"
)

// applyProject returns a JSON object containing only the fields named in the `pick`
// list, optionally renaming them using the `rename` map.
//
// The `arg` must be a JSON object that may contain:
//   - `"pick"`: a JSON array of field names to include (all fields kept when absent).
//   - `"rename"`: a JSON object whose keys are original field names and whose values
//     are the desired output names.
//
// If the input is not a JSON object the input is returned unchanged.
//
// Example:
//
//	// projection only
//	fj.Get(`{"name":"Alice","age":30,"city":"NY"}`, `@project:{"pick":["name","age"]}`)
//	// → {"name":"Alice","age":30}
//
//	// rename only
//	fj.Get(`{"name":"Alice","age":30}`, `@project:{"rename":{"name":"fullName"}}`)
//	// → {"fullName":"Alice","age":30}
//
//	// projection + rename combined
//	fj.Get(`{"name":"Alice","age":30,"city":"NY"}`, `@project:{"pick":["name","age"],"rename":{"name":"fullName","age":"years"}}`)
//	// → {"fullName":"Alice","years":30}
//
// Performance: O(fields) with a single Foreach pass; builds output into a pre-grown
// byte slice, so only one allocation per call.
func applyProject(json, arg string) string {
	ctx := Parse(json)
	if !ctx.IsObject() {
		return json
	}
	var pick []string
	rename := make(map[string]string)
	if strutil.IsNotEmpty(arg) {
		Parse(arg).Foreach(func(key, value Context) bool {
			switch key.String() {
			case "pick":
				value.Foreach(func(_, v Context) bool {
					pick = append(pick, v.String())
					return true
				})
			case "rename":
				value.Foreach(func(k, v Context) bool {
					rename[k.String()] = v.String()
					return true
				})
			}
			return true
		})
	}
	pickSet := make(map[string]bool, len(pick))
	for _, f := range pick {
		pickSet[f] = true
	}
	out := make([]byte, 0, len(json))
	out = append(out, '{')
	var i int
	ctx.Foreach(func(key, value Context) bool {
		fieldName := key.String()
		if len(pickSet) > 0 && !pickSet[fieldName] {
			return true
		}
		outputName := fieldName
		if newName, ok := rename[fieldName]; ok {
			outputName = newName
		}
		if i > 0 {
			out = append(out, ',')
		}
		out = appendJSONString(out, outputName)
		out = append(out, ':')
		out = append(out, value.raw...)
		i++
		return true
	})
	out = append(out, '}')
	return strutil.SafeStr(out)
}

// applyFilter removes elements from a JSON array that do not satisfy a condition.
//
// The `arg` must be a JSON object with:
//   - `"key"` (required): the field name to test on each array element.
//   - `"value"` (required): the value to compare against.
//   - `"op"` (optional): comparison operator, one of:
//     `"eq"` (default), `"ne"`, `"gt"`, `"gte"`, `"lt"`, `"lte"`, `"contains"`.
//
// Elements that are not JSON objects are kept unchanged when the op cannot be
// evaluated (fail-open).  If the input is not a JSON array, it is returned
// unchanged.
//
// Example:
//
//	fj.Get(`[{"name":"Alice","age":30},{"name":"Bob","age":25}]`,
//	        `@filter:{"key":"age","op":"gt","value":28}`)
//	// → [{"name":"Alice","age":30}]
//
//	fj.Get(`[{"status":"active"},{"status":"inactive"}]`,
//	        `@filter:{"key":"status","value":"active"}`)
//	// → [{"status":"active"}]
//
// Performance: single Foreach pass; uses a pre-grown byte slice.
func applyFilter(json, arg string) string {
	ctx := Parse(json)
	if !ctx.IsArray() {
		return json
	}
	var key, op, rawVal string
	op = "eq"
	Parse(arg).Foreach(func(k, v Context) bool {
		switch k.String() {
		case "key":
			key = v.String()
		case "op":
			op = v.String()
		case "value":
			rawVal = v.raw
		}
		return true
	})
	if strutil.IsEmpty(key) {
		return json
	}
	cmpVal := Parse(rawVal)
	out := make([]byte, 0, len(json))
	out = append(out, '[')
	var i int
	ctx.Foreach(func(_, elem Context) bool {
		fieldVal := elem.Get(key)
		if !fieldVal.Exists() {
			return true
		}
		if matchesCondition(fieldVal, cmpVal, op) {
			if i > 0 {
				out = append(out, ',')
			}
			out = append(out, elem.raw...)
			i++
		}
		return true
	})
	out = append(out, ']')
	return strutil.SafeStr(out)
}

// matchesCondition tests whether `actual` satisfies `op` relative to `expected`.
// Numeric comparisons use float64; all others fall back to string comparison.
func matchesCondition(actual, expected Context, op string) bool {
	switch op {
	case "contains":
		return strings.Contains(actual.String(), expected.String())
	case "ne":
		return actual.raw != expected.raw
	case "gt", "gte", "lt", "lte":
		a, e := actual.Float64(), expected.Float64()
		switch op {
		case "gt":
			return a > e
		case "gte":
			return a >= e
		case "lt":
			return a < e
		case "lte":
			return a <= e
		}
	}
	// default: eq – raw JSON equality (works for strings, numbers, bools, null)
	return actual.raw == expected.raw
}

// applyDefault injects fallback values for fields that are absent or explicitly null
// in a JSON object.
//
// The `arg` must be a JSON object mapping field names to their default values.
// Fields that already exist with a non-null value are left untouched.
// Fields listed in `arg` that are missing from the input object are appended.
// If the input is not a JSON object, it is returned unchanged.
//
// Example:
//
//	fj.Get(`{"name":"Alice","role":null}`,
//	        `@default:{"role":"user","active":true}`)
//	// → {"name":"Alice","role":"user","active":true}
//
// Performance: two Foreach passes (one to collect existing, one to append defaults);
// single allocation for the output byte slice.
func applyDefault(json, arg string) string {
	ctx := Parse(json)
	if !ctx.IsObject() {
		return json
	}
	if strutil.IsEmpty(arg) {
		return json
	}
	defaults := Parse(arg)
	if !defaults.IsObject() {
		return json
	}
	// Collect the set of all fields that appear in the object (including null ones).
	// The second pass only appends keys that are entirely absent from the input.
	present := make(map[string]bool)
	ctx.Foreach(func(key, _ Context) bool {
		present[key.String()] = true
		return true
	})
	out := make([]byte, 0, len(json)+len(arg))
	out = append(out, '{')
	var i int
	// Emit original fields, substituting null with the default when available.
	ctx.Foreach(func(key, value Context) bool {
		if i > 0 {
			out = append(out, ',')
		}
		out = append(out, key.raw...)
		out = append(out, ':')
		if value.kind == Null {
			if dv := defaults.Get(key.String()); dv.Exists() {
				out = append(out, dv.raw...)
			} else {
				out = append(out, value.raw...)
			}
		} else {
			out = append(out, value.raw...)
		}
		i++
		return true
	})
	// Append default fields that were missing entirely.
	defaults.Foreach(func(key, value Context) bool {
		if !present[key.String()] {
			if i > 0 {
				out = append(out, ',')
			}
			out = append(out, key.raw...)
			out = append(out, ':')
			out = append(out, value.raw...)
			i++
		}
		return true
	})
	out = append(out, '}')
	return strutil.SafeStr(out)
}

// applyCoerce converts a scalar JSON value to the type specified by the `arg`.
//
// Supported target types (case-insensitive):
//   - `"string"`: converts the value to a JSON string.
//   - `"number"`: parses the value as a float64 and re-emits it; returns `null`
//     when conversion is not possible.
//   - `"bool"` / `"boolean"`: interprets truthy values (non-zero numbers, "true",
//     "1", "yes") as `true`, everything else as `false`.
//
// Objects and arrays are returned unchanged for any target type.
//
// Example:
//
//	fj.Get(`42`,    `@coerce:{"to":"string"}`)  // → "42"
//	fj.Get(`"99"`,  `@coerce:{"to":"number"}`)  // → 99
//	fj.Get(`1`,     `@coerce:{"to":"bool"}`)    // → true
//
// Performance: no heap allocations for number→string and bool conversions.
func applyCoerce(json, arg string) string {
	ctx := Parse(json)
	// Objects and arrays pass through unchanged.
	if ctx.IsObject() || ctx.IsArray() {
		return json
	}
	var to string
	Parse(arg).Foreach(func(key, value Context) bool {
		if key.String() == "to" {
			to = strings.ToLower(value.String())
		}
		return true
	})
	switch to {
	case "string":
		s := ctx.String()
		return string(appendJSONString(nil, s))
	case "number":
		f := ctx.Float64()
		if math.IsNaN(f) || math.IsInf(f, 0) {
			return "null"
		}
		return fmt.Sprintf("%g", f)
	case "bool", "boolean":
		if ctx.Bool() {
			return "true"
		}
		return "false"
	}
	return json
}

// applyCount returns the number of elements in a JSON array, or the number of
// key-value pairs in a JSON object, as a plain JSON integer.
//
// For scalar values (strings, numbers, booleans, null) the result is always 0.
//
// Example:
//
//	fj.Get(`[1,2,3]`,          `@count`) // → 3
//	fj.Get(`{"a":1,"b":2}`,    `@count`) // → 2
//	fj.Get(`"hello"`,          `@count`) // → 0
//
// Performance: single Foreach pass with no allocations.
func applyCount(json, arg string) string {
	ctx := Parse(json)
	if !ctx.IsArray() && !ctx.IsObject() {
		return "0"
	}
	var n int
	ctx.Foreach(func(_, _ Context) bool {
		n++
		return true
	})
	return fmt.Sprintf("%d", n)
}

// applyFirst returns the first element of a JSON array as a raw JSON value.
// Returns `null` if the array is empty or the input is not an array.
//
// Example:
//
//	fj.Get(`[10,20,30]`, `@first`) // → 10
//	fj.Get(`[]`,         `@first`) // → null
//
// Performance: early-exit Foreach; at most one element is examined.
func applyFirst(json, arg string) string {
	ctx := Parse(json)
	if !ctx.IsArray() {
		return "null"
	}
	var first string
	ctx.Foreach(func(_, value Context) bool {
		first = value.raw
		return false // stop after the first element
	})
	if first == "" {
		return "null"
	}
	return first
}

// applyLast returns the last element of a JSON array as a raw JSON value.
// Returns `null` if the array is empty or the input is not an array.
//
// Example:
//
//	fj.Get(`[10,20,30]`, `@last`) // → 30
//	fj.Get(`[]`,         `@last`) // → null
//
// Performance: full Foreach pass; the last raw value overwrites on each iteration.
func applyLast(json, arg string) string {
	ctx := Parse(json)
	if !ctx.IsArray() {
		return "null"
	}
	var last string
	ctx.Foreach(func(_, value Context) bool {
		last = value.raw
		return true
	})
	if last == "" {
		return "null"
	}
	return last
}

// applySum returns the arithmetic sum of all numeric values in a JSON array.
// Non-numeric elements (strings, objects, arrays, booleans, null) are skipped.
// Returns `0` when the input is not an array or the array contains no numbers.
//
// Example:
//
//	fj.Get(`[1,2,3,4]`,          `@sum`) // → 10
//	fj.Get(`[1.5,2.5,"x",null]`, `@sum`) // → 4
//
// Performance: single Foreach pass; no allocations beyond the return string.
func applySum(json, arg string) string {
	ctx := Parse(json)
	if !ctx.IsArray() {
		return "0"
	}
	var sum float64
	ctx.Foreach(func(_, value Context) bool {
		if value.kind == Number {
			sum += value.Float64()
		}
		return true
	})
	return formatNumber(sum)
}

// applyMin returns the minimum numeric value found in a JSON array.
// Non-numeric elements are skipped.  Returns `null` if the array is empty or
// contains no numbers.
//
// Example:
//
//	fj.Get(`[3,1,4,1,5]`, `@min`) // → 1
//
// Performance: single Foreach pass; no heap allocations.
func applyMin(json, arg string) string {
	ctx := Parse(json)
	if !ctx.IsArray() {
		return "null"
	}
	min := math.MaxFloat64
	found := false
	ctx.Foreach(func(_, value Context) bool {
		if value.kind == Number {
			f := value.Float64()
			if !found || f < min {
				min = f
				found = true
			}
		}
		return true
	})
	if !found {
		return "null"
	}
	return formatNumber(min)
}

// applyMax returns the maximum numeric value found in a JSON array.
// Non-numeric elements are skipped.  Returns `null` if the array is empty or
// contains no numbers.
//
// Example:
//
//	fj.Get(`[3,1,4,1,5]`, `@max`) // → 5
//
// Performance: single Foreach pass; no heap allocations.
func applyMax(json, arg string) string {
	ctx := Parse(json)
	if !ctx.IsArray() {
		return "null"
	}
	max := -math.MaxFloat64
	found := false
	ctx.Foreach(func(_, value Context) bool {
		if value.kind == Number {
			f := value.Float64()
			if !found || f > max {
				max = f
				found = true
			}
		}
		return true
	})
	if !found {
		return "null"
	}
	return formatNumber(max)
}

// applyPluck extracts a named field from every element of a JSON array,
// returning a new JSON array of the extracted values.  Elements that do not
// contain the named field are omitted from the result.
//
// The `arg` is the plain field name (not a JSON string).  Nested path
// expressions are supported (e.g. `"address.city"`).
//
// Example:
//
//	fj.Get(`[{"name":"Alice","age":30},{"name":"Bob","age":25}]`, `@pluck:name`)
//	// → ["Alice","Bob"]
//
//	fj.Get(`[{"addr":{"city":"NY"}},{"addr":{"city":"LA"}}]`, `@pluck:addr.city`)
//	// → ["NY","LA"]
//
// Performance: single Foreach pass; pre-grown output slice.
func applyPluck(json, arg string) string {
	ctx := Parse(json)
	if !ctx.IsArray() {
		return "[]"
	}
	fieldPath := strings.TrimSpace(arg)
	if strutil.IsEmpty(fieldPath) {
		return "[]"
	}
	out := make([]byte, 0, len(json))
	out = append(out, '[')
	var i int
	ctx.Foreach(func(_, elem Context) bool {
		val := elem.Get(fieldPath)
		if val.Exists() {
			if i > 0 {
				out = append(out, ',')
			}
			raw := val.raw
			if raw == "" {
				raw = "null"
			}
			out = append(out, raw...)
			i++
		}
		return true
	})
	out = append(out, ']')
	return strutil.SafeStr(out)
}

// formatNumber renders a float64 as a compact JSON number string.
// Integers are emitted without a decimal point.
func formatNumber(f float64) string {
	if f == math.Trunc(f) && !math.IsInf(f, 0) {
		return fmt.Sprintf("%d", int64(f))
	}
	return fmt.Sprintf("%g", f)
}

// init registers the extended transformers in the global registry.
func init() {
	// Structural / object transformers
	globalRegistry.Register("project", TransformerFunc(applyProject))

	// Array filtering and aggregation
	globalRegistry.Register("filter", TransformerFunc(applyFilter))
	globalRegistry.Register("count", TransformerFunc(applyCount))
	globalRegistry.Register("first", TransformerFunc(applyFirst))
	globalRegistry.Register("last", TransformerFunc(applyLast))
	globalRegistry.Register("sum", TransformerFunc(applySum))
	globalRegistry.Register("min", TransformerFunc(applyMin))
	globalRegistry.Register("max", TransformerFunc(applyMax))
	globalRegistry.Register("pluck", TransformerFunc(applyPluck))

	// Value normalization
	globalRegistry.Register("default", TransformerFunc(applyDefault))
	globalRegistry.Register("coerce", TransformerFunc(applyCoerce))
}
