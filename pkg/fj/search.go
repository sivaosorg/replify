package fj

import (
	"math"
	"strings"
)

// Search performs a full-tree scan of the JSON document and returns all leaf values
// whose string representation contains the given keyword (case-sensitive substring
// match). Both string values and non-string scalar values (numbers, booleans, null)
// are compared against the keyword.
//
// The traversal is depth-first; array elements and object values are each visited
// recursively. Compound values (objects and arrays) are never returned as matches —
// only scalar leaves are considered.
//
// Parameters:
//   - json:    A well-formed JSON string to scan.
//   - keyword: The substring to search for. An empty keyword matches every leaf.
//
// Returns:
//   - A slice of Context values whose string representation contains keyword.
//     Returns an empty (non-nil) slice when there are no matches.
//
// Example:
//
//	json := `{"users":[{"name":"Alice"},{"name":"Bob"},{"name":"Charlie"}]}`
//	results := fj.Search(json, "Ali")
//	// results[0].String() == "Alice"
func Search(json, keyword string) []Context {
	return searchLeaves(nil, Parse(json), keyword)
}

// searchLeaves is the internal recursive worker for Search.
// It appends to `all` every scalar leaf whose String() contains keyword.
func searchLeaves(all []Context, node Context, keyword string) []Context {
	if node.IsArray() || node.IsObject() {
		node.Foreach(func(_, child Context) bool {
			all = searchLeaves(all, child, keyword)
			return true
		})
		return all
	}
	if !node.Exists() {
		return all
	}
	if isEmpty(keyword) || strings.Contains(node.String(), keyword) {
		all = append(all, node)
	}
	return all
}

// SearchByKey performs a full-tree scan of the JSON document and returns all Context
// values that are stored under any of the given key names. The search is recursive —
// it descends into nested objects and arrays at every depth level.
//
// Key matching is exact and case-sensitive. Multiple key names may be supplied; any
// value whose immediate parent key matches at least one of them is included.
//
// Parameters:
//   - json: A well-formed JSON string to scan.
//   - keys: One or more object key names to look up. If no keys are provided the
//     function returns an empty slice.
//
// Returns:
//   - A slice of Context values stored under the given keys, in depth-first order.
//
// Example:
//
//	json := `{"a":{"title":"Go"},"b":{"title":"Rust"},"c":{"other":"x"}}`
//	results := fj.SearchByKey(json, "title")
//	// len(results) == 2, results[0].String() == "Go"
func SearchByKey(json string, keys ...string) []Context {
	if len(keys) == 0 {
		return []Context{}
	}
	keySet := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		keySet[k] = struct{}{}
	}
	return searchByKeyRecursive(nil, Parse(json), keySet)
}

// searchByKeyRecursive is the internal recursive worker for SearchByKey.
func searchByKeyRecursive(all []Context, node Context, keySet map[string]struct{}) []Context {
	if node.IsObject() {
		node.Foreach(func(key, val Context) bool {
			if _, ok := keySet[key.String()]; ok {
				all = append(all, val)
			}
			// Recurse into value regardless of whether the key matched.
			if val.IsObject() || val.IsArray() {
				all = searchByKeyRecursive(all, val, keySet)
			}
			return true
		})
		return all
	}
	if node.IsArray() {
		node.Foreach(func(_, child Context) bool {
			all = searchByKeyRecursive(all, child, keySet)
			return true
		})
	}
	return all
}

// Contains reports whether the result of querying json at path contains the given
// target substring (case-sensitive). If the path does not exist, Contains returns
// false.
//
// Parameters:
//   - json:   A well-formed JSON string.
//   - path:   A fj dot-notation path.
//   - target: The substring to look for within the string representation of the
//     value found at path.
//
// Returns:
//   - true if the value exists and its string representation contains target.
//
// Example:
//
//	json := `{"msg":"hello world"}`
//	fj.Contains(json, "msg", "world") // true
//	fj.Contains(json, "msg", "xyz")   // false
func Contains(json, path, target string) bool {
	ctx := Get(json, path)
	if !ctx.Exists() {
		return false
	}
	return strings.Contains(ctx.String(), target)
}

// FindPath returns the first dot-notation path in the JSON document at which the
// given scalar value can be found (case-sensitive, exact string match against
// Context.String()). Object keys and array indices are joined by ".".
//
// FindPath only searches leaf (scalar) values. If no leaf matches, an empty string
// is returned.
//
// Parameters:
//   - json:  A well-formed JSON string.
//   - value: The scalar string value to locate.
//
// Returns:
//   - The dot-notation path of the first matching leaf, or "" when not found.
//
// Example:
//
//	json := `{"user":{"name":"Alice","age":30}}`
//	fj.FindPath(json, "Alice") // "user.name"
func FindPath(json, value string) string {
	path, _ := findPathRecursive(Parse(json), value, "")
	return path
}

// FindPaths returns the dot-notation paths for every leaf in the JSON document
// whose string representation exactly equals value (case-sensitive).
//
// Parameters:
//   - json:  A well-formed JSON string.
//   - value: The scalar string value to locate.
//
// Returns:
//   - All matching paths in depth-first order. Returns an empty slice when there
//     are no matches.
//
// Example:
//
//	json := `{"a":"x","b":{"c":"x","d":"y"}}`
//	fj.FindPaths(json, "x") // ["a", "b.c"]
func FindPaths(json, value string) []string {
	return findAllPathsRecursive(nil, Parse(json), value, "")
}

// findPathRecursive is the depth-first worker for FindPath.
// Returns the first matching path and a bool indicating whether it was found.
func findPathRecursive(node Context, value, prefix string) (string, bool) {
	if node.IsObject() {
		var found string
		var ok bool
		node.Foreach(func(key, child Context) bool {
			p := joinPath(prefix, key.String())
			if child.IsObject() || child.IsArray() {
				found, ok = findPathRecursive(child, value, p)
			} else if child.Exists() && child.String() == value {
				found, ok = p, true
			}
			return !ok
		})
		return found, ok
	}
	if node.IsArray() {
		var found string
		var ok bool
		idx := 0
		node.Foreach(func(_, child Context) bool {
			p := joinPath(prefix, itoa(idx))
			if child.IsObject() || child.IsArray() {
				found, ok = findPathRecursive(child, value, p)
			} else if child.Exists() && child.String() == value {
				found, ok = p, true
			}
			idx++
			return !ok
		})
		return found, ok
	}
	return "", false
}

// findAllPathsRecursive is the depth-first worker for FindPaths.
func findAllPathsRecursive(all []string, node Context, value, prefix string) []string {
	if node.IsObject() {
		node.Foreach(func(key, child Context) bool {
			p := joinPath(prefix, key.String())
			if child.IsObject() || child.IsArray() {
				all = findAllPathsRecursive(all, child, value, p)
			} else if child.Exists() && child.String() == value {
				all = append(all, p)
			}
			return true
		})
		return all
	}
	if node.IsArray() {
		idx := 0
		node.Foreach(func(_, child Context) bool {
			p := joinPath(prefix, itoa(idx))
			if child.IsObject() || child.IsArray() {
				all = findAllPathsRecursive(all, child, value, p)
			} else if child.Exists() && child.String() == value {
				all = append(all, p)
			}
			idx++
			return true
		})
	}
	return all
}

// joinPath concatenates a dot-notation prefix with a segment, inserting "." only
// when the prefix is non-empty.
func joinPath(prefix, segment string) string {
	if prefix == "" {
		return segment
	}
	return prefix + "." + segment
}

// itoa converts a non-negative integer to its decimal string representation without
// importing the strconv package.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}

// Count returns the number of elements returned by evaluating path against json.
// For a path that produces a JSON array the count equals the array length. For a
// path that produces a single scalar value, Count returns 1. For a missing or null
// result, Count returns 0.
//
// Parameters:
//   - json: A well-formed JSON string.
//   - path: A fj dot-notation path.
//
// Returns:
//   - The count of matching elements (≥ 0).
//
// Example:
//
//	json := `{"tags":["go","json","fast"]}`
//	fj.Count(json, "tags")   // 3
//	fj.Count(json, "tags.0") // 1
//	fj.Count(json, "missing")// 0
func Count(json, path string) int {
	ctx := Get(json, path)
	if !ctx.Exists() {
		return 0
	}
	if ctx.IsArray() {
		return len(ctx.Array())
	}
	return 1
}

// Sum returns the sum of all numeric values produced by evaluating path against
// json. Non-numeric results are silently ignored. Returns 0 when no numeric values
// are found.
//
// Parameters:
//   - json: A well-formed JSON string.
//   - path: A fj dot-notation path. The path may resolve to a JSON array of numbers
//     (e.g. "scores") or a single number (e.g. "scores.0").
//
// Returns:
//   - The sum as float64.
//
// Example:
//
//	json := `{"scores":[10,20,30]}`
//	fj.Sum(json, "scores") // 60.0
func Sum(json, path string) float64 {
	var total float64
	collectNumbers(json, path, func(n float64) { total += n })
	return total
}

// Min returns the minimum numeric value among all results produced by evaluating
// path against json. Non-numeric results are silently ignored.
//
// Parameters:
//   - json: A well-formed JSON string.
//   - path: A fj dot-notation path.
//
// Returns:
//   - The minimum value and true when at least one number is found.
//   - 0 and false when no numeric values are found.
//
// Example:
//
//	json := `{"scores":[10,20,5,30]}`
//	v, ok := fj.Min(json, "scores") // 5.0, true
func Min(json, path string) (float64, bool) {
	min := math.MaxFloat64
	found := false
	collectNumbers(json, path, func(n float64) {
		if n < min {
			min = n
		}
		found = true
	})
	if !found {
		return 0, false
	}
	return min, true
}

// Max returns the maximum numeric value among all results produced by evaluating
// path against json. Non-numeric results are silently ignored.
//
// Parameters:
//   - json: A well-formed JSON string.
//   - path: A fj dot-notation path.
//
// Returns:
//   - The maximum value and true when at least one number is found.
//   - 0 and false when no numeric values are found.
//
// Example:
//
//	json := `{"scores":[10,20,5,30]}`
//	v, ok := fj.Max(json, "scores") // 30.0, true
func Max(json, path string) (float64, bool) {
	max := -math.MaxFloat64
	found := false
	collectNumbers(json, path, func(n float64) {
		if n > max {
			max = n
		}
		found = true
	})
	if !found {
		return 0, false
	}
	return max, true
}

// Avg returns the arithmetic mean of all numeric values produced by evaluating path
// against json. Non-numeric results are silently ignored.
//
// Parameters:
//   - json: A well-formed JSON string.
//   - path: A fj dot-notation path.
//
// Returns:
//   - The average value and true when at least one number is found.
//   - 0 and false when no numeric values are found.
//
// Example:
//
//	json := `{"scores":[10,20,30]}`
//	v, ok := fj.Avg(json, "scores") // 20.0, true
func Avg(json, path string) (float64, bool) {
	var total float64
	var n int
	collectNumbers(json, path, func(v float64) {
		total += v
		n++
	})
	if n == 0 {
		return 0, false
	}
	return total / float64(n), true
}

// collectNumbers is an internal helper shared by Sum, Min, Max, and Avg.
// It visits every Context returned by path (treating a JSON array result as a
// sequence of individual values) and calls fn for each numeric one.
func collectNumbers(json, path string, fn func(float64)) {
	ctx := Get(json, path)
	if !ctx.Exists() {
		return
	}
	if ctx.IsArray() {
		ctx.Foreach(func(_, item Context) bool {
			if item.kind == Number {
				fn(item.Float64())
			}
			return true
		})
		return
	}
	if ctx.kind == Number {
		fn(ctx.Float64())
	}
}

// Filter evaluates path against json, treats the result as an array, and returns
// only those elements for which fn returns true.
//
// If path resolves to a non-array value (including a single scalar), that single
// value is treated as a one-element collection. If the path does not exist, an
// empty slice is returned.
//
// Parameters:
//   - json: A well-formed JSON string.
//   - path: A fj dot-notation path.
//   - fn:   A predicate applied to each element. Returning true keeps the element.
//
// Returns:
//   - A slice of matching Context values.
//
// Example:
//
//	json := `{"items":[1,2,3,4,5]}`
//	results := fj.Filter(json, "items", func(ctx fj.Context) bool {
//	    return ctx.Float64() > 2
//	})
//	// results holds Context values for 3, 4, 5
func Filter(json, path string, fn func(Context) bool) []Context {
	ctx := Get(json, path)
	if !ctx.Exists() {
		return []Context{}
	}
	items := ctx.Array()
	out := make([]Context, 0, len(items))
	for _, item := range items {
		if fn(item) {
			out = append(out, item)
		}
	}
	return out
}

// First evaluates path against json and returns the first element for which fn
// returns true. If no element matches, a zero-value Context is returned (use
// Context.Exists() to check).
//
// Parameters:
//   - json: A well-formed JSON string.
//   - path: A fj dot-notation path.
//   - fn:   A predicate applied to each element.
//
// Returns:
//   - The first matching Context, or a zero-value Context when not found.
//
// Example:
//
//	json := `{"items":[1,2,3,4,5]}`
//	ctx := fj.First(json, "items", func(c fj.Context) bool {
//	    return c.Float64() > 3
//	})
//	ctx.Float64() // 4
func First(json, path string, fn func(Context) bool) Context {
	ctx := Get(json, path)
	if !ctx.Exists() {
		return Context{}
	}
	for _, item := range ctx.Array() {
		if fn(item) {
			return item
		}
	}
	return Context{}
}

// Distinct evaluates path against json and returns a deduplicated slice of values
// using the string representation of each element as the equality key. Order is
// preserved (first occurrence wins).
//
// Parameters:
//   - json: A well-formed JSON string.
//   - path: A fj dot-notation path.
//
// Returns:
//   - A slice of unique Context values. Returns an empty slice when the path does
//     not exist.
//
// Example:
//
//	json := `{"tags":["go","json","go","fast","json"]}`
//	results := fj.Distinct(json, "tags")
//	// len(results) == 3: "go", "json", "fast"
func Distinct(json, path string) []Context {
	ctx := Get(json, path)
	if !ctx.Exists() {
		return []Context{}
	}
	seen := make(map[string]struct{})
	var out []Context
	for _, item := range ctx.Array() {
		key := item.String()
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			out = append(out, item)
		}
	}
	if out == nil {
		return []Context{}
	}
	return out
}

// Pluck evaluates path against json (expecting a JSON array of objects), and for
// each element builds a new JSON object containing only the specified fields. Fields
// absent from an element are omitted from its output object.
//
// Parameters:
//   - json:   A well-formed JSON string.
//   - path:   A fj dot-notation path resolving to an array of objects.
//   - fields: The field names to extract from each object.
//
// Returns:
//   - A slice of Context values, one per element in the source array, each
//     wrapping the projected JSON object. Returns an empty slice when path does
//     not exist or no fields are provided.
//
// Example:
//
//	json := `{"users":[
//	    {"id":1,"name":"Alice","email":"a@example.com"},
//	    {"id":2,"name":"Bob","email":"b@example.com"}
//	]}`
//	results := fj.Pluck(json, "users", "id", "name")
//	// results[0].String() == `{"id":1,"name":"Alice"}`
//	// results[1].String() == `{"id":2,"name":"Bob"}`
func Pluck(json, path string, fields ...string) []Context {
	if len(fields) == 0 {
		return []Context{}
	}
	ctx := Get(json, path)
	if !ctx.Exists() {
		return []Context{}
	}
	items := ctx.Array()
	out := make([]Context, 0, len(items))
	for _, item := range items {
		if !item.IsObject() {
			continue
		}
		var b strings.Builder
		b.WriteByte('{')
		wrote := 0
		for _, f := range fields {
			v := item.Get(f)
			if !v.Exists() {
				continue
			}
			if wrote > 0 {
				b.WriteByte(',')
			}
			b.WriteString(`"` + f + `":`)
			raw := v.unprocessed
			if raw == "" {
				raw = v.String()
			}
			b.WriteString(raw)
			wrote++
		}
		b.WriteByte('}')
		projected := Parse(b.String())
		out = append(out, projected)
	}
	return out
}
