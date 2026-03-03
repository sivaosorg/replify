// Package msort provides generic, order-aware iteration over Go maps.
//
// Go maps have no defined iteration order. msort converts a map into a
// sorted slice of key-value pairs so that callers can iterate entries in a
// predictable sequence, retrieve the top-N entries by key or value, or
// convert the sorted slice back to a map for subsequent operations.
//
// # Sorting Functions
//
//	msort.SortKey(m)           // ascending by key
//	msort.SortKeyDesc(m)       // descending by key
//	msort.SortValue(m)         // ascending by value
//	msort.SortValueDesc(m)     // descending by value
//	msort.SortTimeValue(m)     // ascending by time.Time value
//	msort.SortTimeValueDesc(m) // descending by time.Time value
//	msort.SortFunc(m, less)    // custom comparison function
//
// Stable variants (SortKeyStable, SortValueStable) preserve the original
// iteration order of entries with equal sort keys, which is useful for
// deterministic output when values may collide.
//
// # Result Operations
//
// Every sort function returns a typed slice that supports fluent chaining:
//
//	top5 := msort.SortValue(scores).Top(5)
//	keys := top5.Keys()      // []K in sorted order
//	vals := top5.Values()    // []V in sorted order
//	m    := top5.ToMap()     // back to map[K]V (order lost)
//
// # Type Constraints
//
// Key and value types must satisfy the Ordered constraint — any type whose
// underlying type is a numeric or string kind — except for SortTimeValue,
// which handles time.Time separately because it does not implement Ordered.
// SortFunc accepts any comparable key type with an arbitrary comparison
// function, providing maximum flexibility.
//
// All functions are pure (they do not modify the input map) and safe for
// concurrent use provided the map is not mutated during the call.
package msort
