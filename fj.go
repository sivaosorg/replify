package replify

import (
	"github.com/sivaosorg/replify/pkg/fj"
)

// BodyCtx parses the body of the wrapper as JSON and returns a fj.Context for the
// entire document. This is the entry point for all fj-based operations on the wrapper.
//
// If the body is nil or cannot be serialized, a zero-value fj.Context is returned.
// Callers can check presence with ctx.Exists().
//
// Example:
//
//	ctx := w.BodyCtx()
//	fmt.Println(ctx.Get("user.name").String())
func (w *wrapper) BodyCtx() fj.Context {
	return fj.Parse(jsonpass(w.data))
}

// QueryBody retrieves the value at the given fj dot-notation path from the wrapper's
// body. The body is serialized to JSON on each call; for repeated queries on the same
// body, use BodyCtx() once and chain calls on the returned Context.
//
// Parameters:
//   - path: A fj dot-notation path (e.g. "user.name", "items.#.id", "roles.0").
//
// Returns:
//   - A fj.Context for the matched value. Call .Exists() to check presence.
//
// Example:
//
//	name := w.QueryBody("user.name").String()
func (w *wrapper) QueryBody(path string) fj.Context {
	return fj.Get(jsonpass(w.data), path)
}

// QueryBodyMul evaluates multiple fj paths against the body in a single pass and
// returns one fj.Context per path in the same order.
//
// Parameters:
//   - paths: One or more fj dot-notation paths.
//
// Returns:
//   - A slice of fj.Context values, one per path.
//
// Example:
//
//	results := w.QueryBodyMul("user.id", "user.email", "roles.#")
func (w *wrapper) QueryBodyMul(paths ...string) []fj.Context {
	return fj.GetMul(jsonpass(w.data), paths...)
}

// ValidBody reports whether the body of the wrapper is valid JSON.
//
// Returns:
//   - true if the body serializes to well-formed JSON; false otherwise.
//
// Example:
//
//	if !w.ValidBody() {
//	    log.Println("body is not valid JSON")
//	}
func (w *wrapper) ValidBody() bool {
	return fj.IsValidJSON(jsonpass(w.data))
}

// SearchBody performs a full-tree scan of the body JSON and returns all scalar
// leaf values whose string representation contains the given keyword (case-sensitive
// substring match).
//
// Parameters:
//   - keyword: The substring to search for. An empty keyword matches every leaf.
//
// Returns:
//   - A slice of fj.Context values whose string representation contains keyword.
//
// Example:
//
//	hits := w.SearchBody("admin")
//	for _, h := range hits {
//	    fmt.Println(h.String())
//	}
func (w *wrapper) SearchBody(keyword string) []fj.Context {
	return fj.Search(jsonpass(w.data), keyword)
}

// SearchBodyMatch performs a full-tree wildcard scan of the body JSON and returns
// all scalar leaf values whose string representation matches the given pattern.
//
// The pattern supports '*' (any sequence) and '?' (single character) wildcards.
//
// Parameters:
//   - pattern: A wildcard pattern applied to leaf string values.
//
// Example:
//
//	hits := w.SearchBodyMatch("admin*")
func (w *wrapper) SearchBodyMatch(pattern string) []fj.Context {
	return fj.SearchMatch(jsonpass(w.data), pattern)
}

// SearchBodyByKey performs a full-tree scan of the body JSON and returns all values
// stored under any of the given key names, regardless of nesting depth.
//
// Parameters:
//   - keys: One or more exact object key names to look up.
//
// Example:
//
//	emails := w.SearchBodyByKey("email")
func (w *wrapper) SearchBodyByKey(keys ...string) []fj.Context {
	return fj.SearchByKey(jsonpass(w.data), keys...)
}

// SearchBodyByKeyPattern performs a full-tree wildcard scan of the body JSON and
// returns all values stored under object keys that match the given pattern.
//
// Parameters:
//   - keyPattern: A wildcard pattern applied to object key names.
//
// Example:
//
//	hits := w.SearchBodyByKeyPattern("user*")
func (w *wrapper) SearchBodyByKeyPattern(keyPattern string) []fj.Context {
	return fj.SearchByKeyPattern(jsonpass(w.data), keyPattern)
}

// BodyContains reports whether the value at the given path inside the body contains
// the target substring (case-sensitive).
//
// Returns false when the path does not exist.
//
// Example:
//
//	w.BodyContains("user.role", "admin")
func (w *wrapper) BodyContains(path, target string) bool {
	return fj.Contains(jsonpass(w.data), path, target)
}

// BodyContainsMatch reports whether the value at the given path inside the body
// matches the given wildcard pattern.
//
// Returns false when the path does not exist.
//
// Example:
//
//	w.BodyContainsMatch("user.email", "*@example.com")
func (w *wrapper) BodyContainsMatch(path, pattern string) bool {
	return fj.ContainsMatch(jsonpass(w.data), path, pattern)
}

// FindBodyPath returns the first dot-notation path in the body at which a scalar
// value equals the given string (exact, case-sensitive match).
//
// Returns "" when no leaf matches.
//
// Example:
//
//	path := w.FindBodyPath("alice@example.com")
func (w *wrapper) FindBodyPath(value string) string {
	return fj.FindPath(jsonpass(w.data), value)
}

// FindBodyPaths returns all dot-notation paths in the body at which a scalar value
// equals the given string.
//
// Example:
//
//	paths := w.FindBodyPaths("active")
func (w *wrapper) FindBodyPaths(value string) []string {
	return fj.FindPaths(jsonpass(w.data), value)
}

// FindBodyPathMatch returns the first dot-notation path in the body at which a scalar
// value matches the given wildcard pattern.
//
// Example:
//
//	path := w.FindBodyPathMatch("alice*")
func (w *wrapper) FindBodyPathMatch(pattern string) string {
	return fj.FindPathMatch(jsonpass(w.data), pattern)
}

// FindBodyPathsMatch returns all dot-notation paths in the body at which a scalar
// value matches the given wildcard pattern.
//
// Example:
//
//	paths := w.FindBodyPathsMatch("err*")
func (w *wrapper) FindBodyPathsMatch(pattern string) []string {
	return fj.FindPathsMatch(jsonpass(w.data), pattern)
}

// CountBody returns the number of elements at the given path in the body.
// For an array result it returns the array length; for a scalar it returns 1;
// for a missing path it returns 0.
//
// Example:
//
//	n := w.CountBody("items")
func (w *wrapper) CountBody(path string) int {
	return fj.Count(jsonpass(w.data), path)
}

// SumBody returns the sum of all numeric values at the given path in the body.
// Non-numeric elements are ignored. Returns 0 when no numbers are found.
//
// Example:
//
//	total := w.SumBody("items.#.price")
func (w *wrapper) SumBody(path string) float64 {
	return fj.Sum(jsonpass(w.data), path)
}

// MinBody returns the minimum numeric value at the given path in the body.
// Returns (0, false) when no numeric values are found.
//
// Example:
//
//	v, ok := w.MinBody("scores")
func (w *wrapper) MinBody(path string) (float64, bool) {
	return fj.Min(jsonpass(w.data), path)
}

// MaxBody returns the maximum numeric value at the given path in the body.
// Returns (0, false) when no numeric values are found.
//
// Example:
//
//	v, ok := w.MaxBody("scores")
func (w *wrapper) MaxBody(path string) (float64, bool) {
	return fj.Max(jsonpass(w.data), path)
}

// AvgBody returns the arithmetic mean of all numeric values at the given path in the
// body. Returns (0, false) when no numeric values are found.
//
// Example:
//
//	avg, ok := w.AvgBody("ratings")
func (w *wrapper) AvgBody(path string) (float64, bool) {
	return fj.Avg(jsonpass(w.data), path)
}

// CollectBodyFloat64 collects every value at the given path in the body that can be
// coerced to float64 (including string-encoded numbers). Non-numeric values are
// skipped.
//
// Example:
//
//	prices := w.CollectBodyFloat64("items.#.price")
func (w *wrapper) CollectBodyFloat64(path string) []float64 {
	return fj.CollectFloat64(jsonpass(w.data), path)
}

// FilterBody evaluates the given path in the body, treats the result as an array,
// and returns only those elements for which fn returns true.
//
// Example:
//
//	active := w.FilterBody("users", func(ctx fj.Context) bool {
//	    return ctx.Get("active").Bool()
//	})
func (w *wrapper) FilterBody(path string, fn func(fj.Context) bool) []fj.Context {
	return fj.Filter(jsonpass(w.data), path, fn)
}

// FirstBody evaluates the given path in the body and returns the first element for
// which fn returns true. Returns a zero-value fj.Context when not found.
//
// Example:
//
//	admin := w.FirstBody("users", func(ctx fj.Context) bool {
//	    return ctx.Get("role").String() == "admin"
//	})
func (w *wrapper) FirstBody(path string, fn func(fj.Context) bool) fj.Context {
	return fj.First(jsonpass(w.data), path, fn)
}

// DistinctBody evaluates the given path in the body and returns a deduplicated slice
// of values using each element's string representation as the equality key.
// First-occurrence order is preserved.
//
// Example:
//
//	tags := w.DistinctBody("tags")
func (w *wrapper) DistinctBody(path string) []fj.Context {
	return fj.Distinct(jsonpass(w.data), path)
}

// PluckBody evaluates the given path in the body (expected: array of objects) and
// returns a new object for each element containing only the specified fields.
//
// Example:
//
//	rows := w.PluckBody("users", "id", "email")
func (w *wrapper) PluckBody(path string, fields ...string) []fj.Context {
	return fj.Pluck(jsonpass(w.data), path, fields...)
}

// GroupByBody groups the elements at the given path in the body by the string value
// of keyField, using conv.String for key normalization.
//
// Example:
//
//	byRole := w.GroupByBody("users", "role")
func (w *wrapper) GroupByBody(path, keyField string) map[string][]fj.Context {
	return fj.GroupBy(jsonpass(w.data), path, keyField)
}

// SortBodyBy sorts the elements at the given path in the body by the value of
// keyField. Numeric fields are compared as float64; all others fall back to string
// comparison.
//
// Parameters:
//   - path:      A fj path resolving to an array.
//   - keyField:  The field to sort by. Pass "" to sort scalar arrays.
//   - ascending: Sort direction.
//
// Example:
//
//	sorted := w.SortBodyBy("products", "price", true)
func (w *wrapper) SortBodyBy(path, keyField string, ascending bool) []fj.Context {
	return fj.SortBy(jsonpass(w.data), path, keyField, ascending)
}
