package mapsort

import (
	"sort"
	"time"
)

// Ordered is a constraint that permits any ordered type: any type that supports
// the operators < <= >= >. This includes all numeric types, strings, and types
// that are defined with one of these underlying types.
//
// This is defined locally to avoid external dependencies on golang.org/x/exp/constraints.
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

// item represents a key-value pair from a map.
// K is the key type (must be comparable) and V is the value type.
type item[K comparable, V any] struct {
	Key   K
	Value V
}

// LessFunc is a comparison function that returns true if x should be ordered before y.
type LessFunc[K comparable, V any] func(x, y item[K, V]) bool

// items is a slice of map key-value pairs that provides utility methods.
type items[K comparable, V any] []item[K, V]

// timeItem represents a key-value pair where the value is a time.Time.
// This is used for sorting maps with time.Time values.
type timeItem[K comparable] struct {
	Key   K
	Value time.Time
}

// timeItems is a slice of TimeItem for utility methods.
type timeItems[K comparable] []timeItem[K]

// Top returns a slice containing up to n items from the beginning.
// If n exceeds the slice length, all items are returned.
// This method is useful for getting top-N results after sorting.
func (items items[K, V]) Top(n int) items[K, V] {
	if n > len(items) {
		n = len(items)
	}
	return items[:n]
}

// ToMap converts the sorted items back to a map.
// Note: Maps in Go are unordered, so the sorting information is lost.
// This is useful when you need map operations after sorting and filtering.
//
// Example:
//
//	sorted := mapsort.ByValue(m).Top(10)
//	topMap := sorted.ToMap() // Get top 10 as a map
func (items items[K, V]) ToMap() map[K]V {
	result := make(map[K]V, len(items))
	for _, item := range items {
		result[item.Key] = item.Value
	}
	return result
}

// Keys returns a slice of all keys in order.
// Useful when you only need the sorted keys.
//
// Example:
//
//	sortedKeys := mapsort.ByValue(m).Keys()
func (items items[K, V]) Keys() []K {
	keys := make([]K, len(items))
	for i, item := range items {
		keys[i] = item.Key
	}
	return keys
}

// Values returns a slice of all values in order.
// Useful when you only need the sorted values.
//
// Example:
//
//	sortedValues := mapsort.ByKey(m).Values()
func (items items[K, V]) Values() []V {
	values := make([]V, len(items))
	for i, item := range items {
		values[i] = item.Value
	}
	return values
}

// SortFunc sorts a map using a custom comparison function.
// This provides maximum flexibility for custom sorting logic.
// The sorting is performed in-place on a pre-allocated slice for optimal performance.
//
// Example:
//
//	m := map[string]int{"a": 3, "b": 1, "c": 2}
//	sorted := mapsort.SortFunc(m, func(x, y mapsort.Item[string, int]) bool {
//	    return x.Value < y.Value
//	})
func SortFunc[K comparable, V any](m map[K]V, less LessFunc[K, V]) items[K, V] {
	items := make(items[K, V], 0, len(m))
	for k, v := range m {
		items = append(items, item[K, V]{Key: k, Value: v})
	}

	sort.Slice(items, func(i, j int) bool {
		return less(items[i], items[j])
	})

	return items
}

// SortKey sorts a map by its keys in ascending order.
// Uses compile-time type checking to ensure keys are comparable.
// Performance: O(n log n) where n is the number of map entries.
//
// Example:
//
//	m := map[string]int{"charlie": 3, "alice": 1, "bob": 2}
//	sorted := mapsort.SortKey(m) // Returns items ordered: alice, bob, charlie
func SortKey[K Ordered, V any](m map[K]V) items[K, V] {
	items := make(items[K, V], 0, len(m))
	for k, v := range m {
		items = append(items, item[K, V]{Key: k, Value: v})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Key < items[j].Key
	})

	return items
}

// SortKeyDesc sorts a map by its keys in descending order.
// Uses compile-time type checking to ensure keys are comparable.
// Performance: O(n log n) where n is the number of map entries.
//
// Example:
//
//	m := map[string]int{"charlie": 3, "alice": 1, "bob": 2}
//	sorted := mapsort.SortKeyDesc(m) // Returns items ordered: charlie, bob, alice
func SortKeyDesc[K Ordered, V any](m map[K]V) items[K, V] {
	items := make(items[K, V], 0, len(m))
	for k, v := range m {
		items = append(items, item[K, V]{Key: k, Value: v})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Key > items[j].Key
	})

	return items
}

// SortValue sorts a map by its values in ascending order.
// Uses compile-time type checking to ensure values are comparable.
// Performance: O(n log n) where n is the number of map entries.
//
// Example:
//
//	m := map[string]int{"a": 3, "b": 1, "c": 2}
//	sorted := mapsort.SortValue(m) // Returns items ordered by values: 1, 2, 3
func SortValue[K comparable, V Ordered](m map[K]V) items[K, V] {
	items := make(items[K, V], 0, len(m))
	for k, v := range m {
		items = append(items, item[K, V]{Key: k, Value: v})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Value < items[j].Value
	})

	return items
}

// SortValueDesc sorts a map by its values in descending order.
// Uses compile-time type checking to ensure values are comparable.
// Performance: O(n log n) where n is the number of map entries.
//
// Example:
//
//	m := map[string]int{"a": 3, "b": 1, "c": 2}
//	sorted := mapsort.SortValueDesc(m) // Returns items ordered by values: 3, 2, 1
func SortValueDesc[K comparable, V Ordered](m map[K]V) items[K, V] {
	items := make(items[K, V], 0, len(m))
	for k, v := range m {
		items = append(items, item[K, V]{Key: k, Value: v})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Value > items[j].Value
	})

	return items
}

// Top returns up to n items from the beginning of the slice.
// If n exceeds the slice length, all items are returned.
//
// This method is useful for getting top-N results after sorting.
//
// Example:
//
//	sorted := mapsort.ByValue(m).Top(5) // Get top 5 items
func (items timeItems[K]) Top(n int) timeItems[K] {
	if n > len(items) {
		n = len(items)
	}
	return items[:n]
}

// ToMap converts the sorted items back to a map.
// Note: Maps in Go are unordered, so the sorting information is lost.
//
// This is useful when you need map operations after sorting and filtering.
//
// Example:
//
//	sorted := mapsort.ByValue(m).Top(10)
//	topMap := sorted.ToMap() // Get top 10 as a map
func (items timeItems[K]) ToMap() map[K]time.Time {
	result := make(map[K]time.Time, len(items))
	for _, item := range items {
		result[item.Key] = item.Value
	}
	return result
}

// Keys returns a slice of all keys in order.
// Useful when you only need the sorted keys.
//
// Example:
//
//	sortedKeys := mapsort.ByValue(m).Keys()
func (items timeItems[K]) Keys() []K {
	keys := make([]K, len(items))
	for i, item := range items {
		keys[i] = item.Key
	}
	return keys
}

// Values returns a slice of all time.Time values in order.
// Useful when you only need the sorted values.
//
// Example:
//
//	sortedValues := mapsort.ByKey(m).Values()
func (items timeItems[K]) Values() []time.Time {
	values := make([]time.Time, len(items))
	for i, item := range items {
		values[i] = item.Value
	}
	return values
}

// SortTimeValue sorts a map with time.Time values in chronological order (earliest first).
// This specialized function handles time.Time which doesn't implement constraints.Ordered.
// Performance: O(n log n) where n is the number of map entries.
//
// Example:
//
//	m := map[string]time.Time{
//	    "event1": time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
//	    "event2": time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
//	}
//	sorted := mapsort.SortTimeValue(m) // event2 comes before event1
func SortTimeValue[K comparable](m map[K]time.Time) timeItems[K] {
	items := make(timeItems[K], 0, len(m))
	for k, v := range m {
		items = append(items, timeItem[K]{Key: k, Value: v})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Value.Before(items[j].Value)
	})

	return items
}

// SortTimeValueDesc sorts a map with time.Time values in reverse chronological order (latest first).
// This specialized function handles time.Time which doesn't implement constraints.Ordered.
// Performance: O(n log n) where n is the number of map entries.
//
// Example:
//
//	m := map[string]time.Time{
//	    "event1": time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
//	    "event2": time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
//	}
//	sorted := mapsort.SortTimeValueDesc(m) // event1 comes before event2
func SortTimeValueDesc[K comparable](m map[K]time.Time) timeItems[K] {
	items := make(timeItems[K], 0, len(m))
	for k, v := range m {
		items = append(items, timeItem[K]{Key: k, Value: v})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Value.After(items[j].Value)
	})

	return items
}

// Stable sorting variants for when order of equal elements matters
//
// SortKeyStable sorts a map by keys in ascending order, preserving order of equal elements.
// Note: Since map keys are unique, this behaves the same as ByKey but uses stable sort.
func SortKeyStable[K Ordered, V any](m map[K]V) items[K, V] {
	items := make(items[K, V], 0, len(m))
	for k, v := range m {
		items = append(items, item[K, V]{Key: k, Value: v})
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Key < items[j].Key
	})

	return items
}

// SortValueStable sorts a map by values in ascending order, preserving original order for equal values.
// Useful when multiple keys may have the same value and you want deterministic ordering.
//
// Example:
//
//	m := map[string]int{"a": 1, "b": 2, "c": 1, "d": 2}
//	sorted := mapsort.SortValueStable(m) // Equal values maintain iteration order
func SortValueStable[K comparable, V Ordered](m map[K]V) items[K, V] {
	items := make(items[K, V], 0, len(m))
	for k, v := range m {
		items = append(items, item[K, V]{Key: k, Value: v})
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Value < items[j].Value
	})

	return items
}
