package common

import (
	"reflect"
	"sort"
)

// Transform applies a transformation function to each element of a collection (slice, array, or map) and returns
// a new collection with the transformed values.
//
// This function takes a collection of any type (using an empty `interface{}`), which can be a slice, array, or map,
// and a mapping function (`mapper`). The function applies the `mapper` function to each value in the collection,
// transforming it. For slices and arrays, the callback is applied to each element. For maps, the callback is applied
// to each key and value. The transformed elements (or key-value pairs) are collected into a new result, which is returned.
//
// The function uses reflection to handle different types of collections and constructs a new collection with the transformed values.
//
// Parameters:
//   - `collection`: The collection (slice, array, or map) to iterate over and transform. It can be of any type.
//   - `mapper`: A function that takes a value from the collection and transforms it. The function is applied to each
//     element of the collection (or key-value pair for maps).
//
// Returns:
//   - A new collection of the same type as the original collection, where each element has been transformed using the
//     provided `mapper` function.
//
// Example:
//
//	// Mapping a slice of integers to their squares
//	numbers := []int{1, 2, 3, 4}
//	squared := Transform(numbers, func(value interface{}) interface{} {
//		return value.(int) * value.(int)
//	})
//	// squared will be []int{1, 4, 9, 16}
//
//	// Mapping a map of strings to their lengths
//	words := map[string]string{"apple": "fruit", "carrot": "vegetable"}
//	lengths := Transform(words, func(value interface{}) interface{} {
//		return len(value.(string))
//	})
//	// lengths will be a new collection containing key-value pairs, where values represent string lengths.
//
// Notes:
//   - For slices and arrays, the `mapper` function is applied to each element.
//   - For maps, the `mapper` function is applied to both the key and the value. The resulting collection will include
//     both transformed keys and values.
//
// Limitations:
//   - The function creates a new collection based on the results of the `mapper` function, so it does not modify the
//     original collection.
func Transform(collection any, predicate func(value any) any) any {
	v := reflect.ValueOf(collection)
	result := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(predicate(v.Index(0).Interface()))), 0, 0)

	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			mappedValue := predicate(v.Index(i).Interface())
			result = reflect.Append(result, reflect.ValueOf(mappedValue))
		}
	} else if v.Kind() == reflect.Map {
		keys := v.MapKeys()
		for _, key := range keys {
			mappedKey := predicate(key.Interface())
			mappedValue := predicate(v.MapIndex(key).Interface())
			result = reflect.Append(result, reflect.ValueOf(mappedKey))
			result = reflect.Append(result, reflect.ValueOf(mappedValue))
		}
	}
	return result.Interface()
}

// Filter filters a collection (slice or array) based on a predicate function and returns a new collection
// containing only the elements that satisfy the condition specified by the predicate.
//
// This function takes a collection of any type (using an empty `interface{}`), which can be a slice or array,
// and a filtering predicate function (`predicate`). The function applies the `predicate` function to each element in
// the collection, and if the predicate returns true, the element is included in the new collection. The function only
// supports slices and arrays as input coll.
//
// The function uses reflection to handle slices and arrays and constructs a new collection containing only the
// elements that pass the filter condition.
//
// Parameters:
//   - `collection`: The collection (slice or array) to filter. It can be of any type, but only slices and arrays
//     are supported.
//   - `predicate`: A function that takes a value from the collection and returns a boolean indicating whether the
//     element should be included in the resulting collection. If it returns true, the element is included.
//
// Returns:
//   - A new collection of the same type as the original collection, containing only the elements that satisfy the
//     condition defined by the `predicate` function.
//
// Example:
//
//	// Filtering a slice of integers to get only even numbers
//	numbers := []int{1, 2, 3, 4, 5, 6}
//	evens := Filter(numbers, func(value interface{}) bool {
//		return value.(int)%2 == 0
//	})
//	// evens will be []int{2, 4, 6}
//
// Notes:
//   - This function only works with slices or arrays, and will return an empty collection if the input is of another type.
//
// Limitations:
//   - The function creates a new collection based on the results of the `predicate` function, so it does not modify
//     the original collection.
func Filter(collection any, predicate func(value any) bool) any {
	v := reflect.ValueOf(collection)
	result := reflect.MakeSlice(v.Type(), 0, 0)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i).Interface()
			if predicate(item) {
				result = reflect.Append(result, reflect.ValueOf(item))
			}
		}
	}
	return result.Interface()
}

// Iterate iterates over a collection (slice, array, or map) and applies a callback function on each element.
//
// This function takes a collection of any type (using an empty `interface{}`), which can be a slice, array, or map,
// and a callback function. The callback function is executed for each element in the collection. For slices and arrays,
// the callback receives the index and the corresponding value. For maps, the callback receives each key and value, with
// the key being passed first followed by the value. For slices and arrays, the index is passed, while for maps, it is
// passed as -1 (since maps are unordered).
//
// The function uses reflection to handle different types of collections and ensures that the correct value is passed
// to the callback.
//
// Parameters:
//   - `collection`: The collection (slice, array, or map) to iterate over. It can be of any type.
//   - `callback`: A function that takes two arguments: the index (for slices and arrays, -1 for maps) and the value
//     from the collection. The callback is executed for each element in the collection.
//
// Example:
//
//	// Iterating over a slice
//	numbers := []int{1, 2, 3, 4}
//	Iterate(numbers, func(index int, value interface{}) {
//		fmt.Printf("Index: %d, Value: %v\n", index, value)
//	})
//	// Output:
//	// Index: 0, Value: 1
//	// Index: 1, Value: 2
//	// Index: 2, Value: 3
//	// Index: 3, Value: 4
//
//	// Iterating over a map
//	colors := map[string]string{"red": "FF0000", "green": "00FF00", "blue": "0000FF"}
//	Iterate(colors, func(index int, value interface{}) {
//		fmt.Printf("Value: %v\n", value)
//	})
//	// Output:
//	// Value: red
//	// Value: FF0000
//	// Value: green
//	// Value: 00FF00
//	// Value: blue
//	// Value: 0000FF
//
// Notes:
//   - For slices and arrays, the callback will receive the index and the value from the collection.
//   - For maps, the callback will be executed twice per key-value pair: once with the key and once with the value,
//     since maps are unordered and the order of key-value pairs cannot be guaranteed.
func Iterate(collection any, callback func(index int, value any)) {
	v := reflect.ValueOf(collection)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			callback(i, v.Index(i).Interface())
		}
	} else if v.Kind() == reflect.Map {
		keys := v.MapKeys()
		for _, key := range keys {
			callback(-1, key.Interface())
			callback(-1, v.MapIndex(key).Interface())
		}
	}
}

// Reduce reduces a collection (slice or array) to a single value by applying a reducer function
// to each element, combining them into a single result.
//
// This function takes a collection of any type (using an empty `interface{}`), which can be a slice or array,
// and a reducer function (`reducer`). The reducer function is applied to each element of the collection, along with an
// accumulator value, and the result of the reducer is passed as the accumulator to the next iteration. This process
// continues until all elements are processed, resulting in a single accumulated value.
//
// The function uses reflection to handle slices and arrays and iterates over the collection to apply the reducer function.
//
// Parameters:
//   - `collection`: The collection (slice or array) to reduce. It can be of any type, but only slices and arrays
//     are supported.
//   - `reducer`: A function that takes two arguments: an accumulator and an element from the collection. It returns
//     a new accumulator value after combining the accumulator and the element.
//   - `initialValue`: The initial value of the accumulator, used as the starting point for the reduction process.
//
// Returns:
//   - A single value, which is the result of reducing the entire collection using the `reducer` function. The return
//     type is the same as the type of the `initialValue`.
//
// Example:
//
//	// Reducing a slice of integers by summing them
//	numbers := []int{1, 2, 3, 4, 5}
//	sum := Reduce(numbers, func(acc interface{}, value interface{}) interface{} {
//		return acc.(int) + value.(int)
//	}, 0)
//	// sum will be 15
//
// Notes:
//   - This function only works with slices or arrays, and will return the initial value if the input collection is empty.
//
// Limitations:
//   - The function creates a single accumulated result by repeatedly applying the `reducer` function to each element,
//     so it does not modify the original collection.
func Reduce(collection any, reducer func(acc any, value any) any, initialValue any) any {
	v := reflect.ValueOf(collection)
	accumulator := initialValue
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			accumulator = reducer(accumulator, v.Index(i).Interface())
		}
	}
	return accumulator
}

// ReduceRight performs a right-to-left reduction on a collection (slice or array) using a reducer function.
//
// This function takes a collection (slice or array), a reducer function, and an initial accumulator value.
// It iterates through the collection from right to left, applying the reducer function to each element and
// accumulating the result. The reducer function is called with two arguments: the current accumulator value
// and the current element in the collection. The final result is returned after processing all elements.
//
// Parameters:
//   - `collection`: A slice or array of any type to reduce.
//   - `reducer`: A function that takes the current accumulator value and the current element as input, and returns
//     the updated accumulator value. This function is applied from right to left on the collection.
//   - `initialValue`: The initial value for the accumulator, which will be passed as the first argument to the
//     reducer function during the first iteration.
//
// Returns:
//   - The final accumulated value after reducing the collection from right to left.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4}
//	result := ReduceRight(numbers, func(acc, value interface{}) interface{} {
//		return acc.(int) + value.(int) // Sum of elements from right to left
//	}, 0)
//	// result will be 10, as the reduction is (0 + 4) + (4 + 3) + (7 + 2) + (9 + 1) = 10
//
// Notes:
//   - The reduction starts from the rightmost element of the collection and proceeds towards the left.
//   - The function uses reflection to support collections of any type.
func ReduceRight(collection any, reducer func(acc, value any) any, initialValue any) any {
	v := reflect.ValueOf(collection)
	accumulator := initialValue
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := v.Len() - 1; i >= 0; i-- {
			accumulator = reducer(accumulator, v.Index(i).Interface())
		}
	}
	return accumulator
}

// Find searches for the first element in a collection (slice or array) that satisfies a given predicate
// function and returns it.
//
// This function takes a collection of any type (using an empty `interface{}`), which can be a slice or array,
// and a predicate function (`predicate`). The predicate function is applied to each element of the collection,
// and if it returns true for an element, that element is returned as the result. The function will return the first
// element that satisfies the condition, and if no elements satisfy the condition, it returns `nil`.
//
// The function uses reflection to handle slices and arrays, iterating over the collection to check each element
// with the provided predicate.
//
// Parameters:
//   - `collection`: The collection (slice or array) to search through. It can be of any type, but only slices and arrays
//     are supported.
//   - `predicate`: A function that takes a value from the collection and returns a boolean indicating whether the element
//     satisfies the condition. If it returns true, the element is returned.
//
// Returns:
//   - The first element from the collection that satisfies the `predicate` function. If no elements satisfy the condition,
//     it returns `nil`.
//
// Example:
//
//	// Finding the first even number in a slice of integers
//	numbers := []int{1, 2, 3, 4, 5}
//	even := Find(numbers, func(value interface{}) bool {
//		return value.(int)%2 == 0
//	})
//	// even will be 2 (the first even number)
//
// Notes:
//   - This function only works with slices or arrays. If the collection is of another type, it will return `nil`.
//   - The function returns the first element that matches the predicate and stops searching after finding it.
//
// Limitations:
//   - The function works only with slices and arrays. If no elements satisfy the predicate, the function will return `nil`,
//     even if the collection is non-empty.
func Find(collection any, predicate func(value any) bool) any {
	v := reflect.ValueOf(collection)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i).Interface()
			if predicate(item) {
				return item
			}
		}
	}
	return nil
}

// All checks whether all elements in a collection (slice or array) satisfy a given condition.
//
// This function takes a collection of any type (using an empty `interface{}`), which can be a slice or array,
// and a condition function (`condition`). The condition function is applied to each element of the collection, and
// the function returns `true` if all elements satisfy the condition. If any element does not satisfy the condition,
// the function returns `false`. If the collection is empty, the function returns `true` (since the condition is trivially satisfied).
//
// The function uses reflection to handle slices and arrays, iterating over the collection to check each element
// with the provided condition.
//
// Parameters:
//   - `collection`: The collection (slice or array) to check. It can be of any type, but only slices and arrays
//     are supported.
//   - `condition`: A function that takes a value from the collection and returns a boolean indicating whether the
//     element satisfies the condition. If it returns `false` for any element, the function immediately returns `false`.
//
// Returns:
//   - `true` if all elements in the collection satisfy the condition, otherwise `false`.
//   - If the collection is empty, the function returns `true`.
//
// Example:
//
//	// Checking if all elements in a slice of integers are positive
//	numbers := []int{1, 2, 3, 4, 5}
//	allPositive := All(numbers, func(value interface{}) bool {
//		return value.(int) > 0
//	})
//	// allPositive will be true
//
// Notes:
//   - This function only works with slices or arrays. If the collection is of another type, it will return `false`.
//   - The function returns `false` as soon as it finds an element that does not satisfy the condition, making it more
//     efficient for early termination.
//
// Limitations:
//   - The function works only with slices and arrays. If no elements satisfy the condition, the function will return `false`,
//     but if all elements are valid, it will return `true`. An empty collection is considered to trivially satisfy the condition.
func All(collection any, condition func(value any) bool) bool {
	v := reflect.ValueOf(collection)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			if !condition(v.Index(i).Interface()) {
				return false
			}
		}
		return true
	}
	return false
}

// Any checks whether any element in a collection (slice or array) satisfies a given condition.
//
// This function accepts a collection (slice or array) and a condition function. It returns `true` if at least one element in the collection satisfies the condition,
// and `false` if none do. The function stops iterating as soon as a matching element is found, making it more efficient for early termination.
//
// Parameters:
//   - `collection`: A slice or array of any type to check.
//   - `condition`: A function that takes a value from the collection and returns a boolean indicating whether the element satisfies the condition.
//
// Returns:
//   - `true` if any element satisfies the condition, `false` otherwise.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	anyNegative := Any(numbers, func(value interface{}) bool {
//	  return value.(int) < 0
//	})
//	// anyNegative will be false because no element is negative.
func Any(collection any, condition func(value any) bool) bool {
	v := reflect.ValueOf(collection)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			if condition(v.Index(i).Interface()) {
				return true
			}
		}
		return false
	}
	return false
}

// Count returns the number of elements in a collection (slice or array) that satisfy a given condition.
//
// This function takes a collection (slice or array) and a condition function. It iterates through the collection, applying the condition to each element.
// It returns the total count of elements that satisfy the condition.
//
// Parameters:
//   - `collection`: A slice or array of any type to check.
//   - `condition`: A function that checks if an element satisfies a condition. The function returns `true` for elements that match the condition, and `false` otherwise.
//
// Returns:
//   - The count of elements in the collection that satisfy the condition.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	countNegative := Count(numbers, func(value interface{}) bool {
//	  return value.(int) < 0
//	})
//	// countNegative will be 0, since no element is negative.
func Count(collection any, condition func(value any) bool) int {
	v := reflect.ValueOf(collection)
	count := 0
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			if condition(v.Index(i).Interface()) {
				count++
			}
		}
	}
	return count
}

// Remove returns a new collection (slice or array) where all elements that satisfy a given condition are removed.
//
// This function takes a collection (slice or array) and a condition function. It iterates through the collection, and for each element that does not satisfy the condition,
// it is added to a new result collection. The function returns a new collection that contains only the elements that do not match the condition.
//
// Parameters:
//   - `collection`: A slice or array of any type to process.
//   - `condition`: A function that checks if an element satisfies a condition. If the element matches the condition, it is removed from the result collection.
//
// Returns:
//   - A new collection (slice or array) with elements that do not satisfy the condition.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	result := Remove(numbers, func(value interface{}) bool {
//	  return value.(int) % 2 == 0 // Removes even numbers
//	})
//	// result will be []int{1, 3, 5}
func Remove(collection any, condition func(value any) bool) any {
	v := reflect.ValueOf(collection)
	result := reflect.MakeSlice(v.Type(), 0, 0)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i).Interface()
			if !condition(item) {
				result = reflect.Append(result, reflect.ValueOf(item))
			}
		}
	}
	return result.Interface()
}

// Sort sorts a collection (slice or array) in-place according to a custom comparison function.
//
// This function takes a collection (slice or array) and a `less` function. The `less` function should return `true` if the element at index `i` should come before the element at index `j`.
// The collection is sorted in-place, meaning the original collection is modified.
//
// Parameters:
//   - `collection`: A slice or array of any type to sort.
//   - `less`: A comparison function that takes two indices `i` and `j` and returns `true` if the element at index `i` should come before the element at index `j`.
//
// Returns:
//   - The collection is sorted in place. No value is returned.
//
// Example:
//
//	numbers := []int{5, 3, 1, 4, 2}
//	Sort(numbers, func(i, j int) bool {
//	  return numbers[i] < numbers[j] // Sort in ascending order
//	})
//	// numbers will be sorted to []int{1, 2, 3, 4, 5}
func Sort(collection any, less func(i, j int) bool) {
	v := reflect.ValueOf(collection)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		sort.SliceStable(collection, func(i, j int) bool {
			return less(i, j)
		})
	}
}

// Reverse reverses the order of elements in a collection (slice or array) in-place.
//
// This function takes a collection (slice or array) and reverses the elements by swapping elements at corresponding positions (i and length-i-1) until the middle of the collection is reached.
// The collection is modified in place, meaning no new collection is created.
//
// Parameters:
//   - `collection`: A slice or array of any type to reverse.
//
// Returns:
//   - The collection is reversed in place. No value is returned.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	Reverse(numbers)
//	// numbers will be reversed to []int{5, 4, 3, 2, 1}
func Reverse(collection any) {
	v := reflect.ValueOf(collection)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		length := v.Len()
		for i := 0; i < length/2; i++ {
			j := length - i - 1
			vi := v.Index(i).Interface()
			vj := v.Index(j).Interface()
			v.Index(i).Set(reflect.ValueOf(vj))
			v.Index(j).Set(reflect.ValueOf(vi))
		}
	}
}

// Unique returns a new collection (slice or array) containing only unique elements from the original collection.
//
// This function takes a collection (slice or array) and removes duplicate elements, returning a new collection with only the unique elements.
// The function uses a map to track elements that have already been encountered, ensuring that only the first occurrence of each element is included in the result.
//
// Parameters:
//   - `collection`: A slice or array of any type from which duplicates should be removed.
//
// Returns:
//   - A new collection (slice or array) containing only the unique elements from the original collection.
//
// Example:
//
//	numbers := []int{1, 2, 2, 3, 4, 4, 5}
//	result := Unique(numbers)
//	// result will be []int{1, 2, 3, 4, 5}
func Unique(collection any) any {
	v := reflect.ValueOf(collection)
	uniqueMap := make(map[any]struct{})
	result := reflect.MakeSlice(v.Type(), 0, 0)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i).Interface()
			if _, found := uniqueMap[item]; !found {
				uniqueMap[item] = struct{}{}
				result = reflect.Append(result, reflect.ValueOf(item))
			}
		}
	}
	return result.Interface()
}

// Contains checks if a given element exists within a collection (slice or array).
//
// This function takes a collection (slice or array) and an element, then iterates through the collection to see if any element matches the provided element.
// It uses `reflect.DeepEqual` to compare elements, which ensures that even complex types (like structs or slices) are compared correctly.
//
// Parameters:
//   - `collection`: A slice or array of any type to search within.
//   - `element`: The element to search for within the collection.
//
// Returns:
//   - `true` if the element is found within the collection, `false` otherwise.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	containsThree := Contains(numbers, 3)
//	// containsThree will be true because 3 is in the slice.
func Contains(collection any, element any) bool {
	v := reflect.ValueOf(collection)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			if reflect.DeepEqual(v.Index(i).Interface(), element) {
				return true
			}
		}
	}
	return false
}

// Difference returns a new collection (slice or array) containing elements from the first collection
// that are not present in the second collection.
//
// This function takes two collections (slices or arrays) and compares the elements of the first collection
// against the elements of the second collection. It returns a new collection with elements that appear in the
// first collection but are absent in the second collection. The function uses `Contains_N` to check for membership
// of each element of the first collection in the second collection.
//
// Parameters:
//   - `collection1`: The first slice or array of any type to compare.
//   - `collection2`: The second slice or array of any type to compare against.
//
// Returns:
//   - A new collection (slice or array) containing the elements from `collection1` that are not in `collection2`.
//
// Example:
//
//	numbers1 := []int{1, 2, 3, 4, 5}
//	numbers2 := []int{3, 4, 6}
//	result := Difference(numbers1, numbers2)
//	// result will be []int{1, 2, 5}, as these are the elements in numbers1 that are not in numbers2.
func Difference(collection1 any, collection2 any) any {
	v1 := reflect.ValueOf(collection1)
	result := reflect.MakeSlice(v1.Type(), 0, 0)
	if v1.Kind() == reflect.Slice || v1.Kind() == reflect.Array {
		for i := 0; i < v1.Len(); i++ {
			item := v1.Index(i).Interface()
			if !Contains(collection2, item) {
				result = reflect.Append(result, v1.Index(i))
			}
		}
	}
	return result.Interface()
}

// Intersection returns a new collection (slice or array) containing elements that are present in both coll.
//
// This function takes two collections (slices or arrays) and compares the elements of both coll.
// It returns a new collection with elements that appear in both the first and the second collection.
// The function uses `Contains_N` to check if an element from the first collection is also present in the second collection.
//
// Parameters:
//   - `collection1`: The first slice or array of any type to compare.
//   - `collection2`: The second slice or array of any type to compare against.
//
// Returns:
//   - A new collection (slice or array) containing the elements that are found in both `collection1` and `collection2`.
//
// Example:
//
//	numbers1 := []int{1, 2, 3, 4, 5}
//	numbers2 := []int{3, 4, 6}
//	result := Intersection(numbers1, numbers2)
//	// result will be []int{3, 4}, as these are the elements common to both numbers1 and numbers2.
func Intersection(collection1 any, collection2 any) any {
	v1 := reflect.ValueOf(collection1)
	result := reflect.MakeSlice(v1.Type(), 0, 0)
	if v1.Kind() == reflect.Slice || v1.Kind() == reflect.Array {
		for i := 0; i < v1.Len(); i++ {
			item := v1.Index(i).Interface()
			if Contains(collection2, item) {
				result = reflect.Append(result, v1.Index(i))
			}
		}
	}
	return result.Interface()
}

// Slice returns a new collection (slice or array) that is a sub-range of the input collection,
// starting from the specified `start` index and ending at the `end` index (exclusive).
//
// This function creates a sub-slice or sub-array from the given collection by copying elements from the
// original collection starting at the `start` index and ending just before the `end` index. If `start` or `end`
// is out of bounds, the function adjusts them to ensure they stay within the valid range for the collection.
// The returned collection will contain the elements between the adjusted `start` and `end` indices.
//
// Parameters:
//   - `collection`: A slice or array of any type to extract a subrange from.
//   - `start`: The starting index (inclusive) for the sub-slice or sub-array.
//   - `end`: The ending index (exclusive) for the sub-slice or sub-array.
//
// Returns:
//   - A new collection (slice or array) containing the elements from `start` to `end` (exclusive).
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5, 6, 7}
//	result := Slice(numbers, 2, 5)
//	// result will be []int{3, 4, 5}, as it extracts elements from index 2 to 4 (inclusive of 2, exclusive of 5).
//
// Notes:
//   - If `start` is less than 0, it is adjusted to 0 (beginning of the collection).
//   - If `end` is greater than the length of the collection, it is adjusted to the collection's length.
func Slice(collection any, start, end int) any {
	v := reflect.ValueOf(collection)
	result := reflect.MakeSlice(v.Type(), 0, 0)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		if start < 0 {
			start = 0
		}
		if end > v.Len() {
			end = v.Len()
		}
		for i := start; i < end; i++ {
			result = reflect.Append(result, v.Index(i))
		}
	}
	return result.Interface()
}

// SliceWithIndices returns a new collection (slice or array) containing elements from the input collection
// specified by the provided indices.
//
// This function creates a new collection by selecting elements from the original collection based on a list
// of indices provided in the `indices` slice. It ensures that only valid indices (within the bounds of the
// collection) are used to build the result collection.
//
// Parameters:
//   - `collection`: A slice or array of any type to extract elements from.
//   - `indices`: A slice of integers representing the indices of elements to include in the result collection.
//     Only valid indices (within the bounds of the collection) are considered.
//
// Returns:
//   - A new collection (slice or array) containing the elements from `collection` at the specified `indices`.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5, 6, 7}
//	indices := []int{1, 3, 5}
//	result := SliceWithIndices(numbers, indices)
//	// result will be []int{2, 4, 6}, as these are the elements at indices 1, 3, and 5 of the original slice.
//
// Notes:
//   - If an index in `indices` is out of bounds (less than 0 or greater than or equal to the collection's length),
//     it is ignored.
func SliceWithIndices(collection any, indices []int) any {
	v := reflect.ValueOf(collection)
	result := reflect.MakeSlice(v.Type(), 0, 0)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for _, index := range indices {
			if index >= 0 && index < v.Len() {
				result = reflect.Append(result, v.Index(index))
			}
		}
	}
	return result.Interface()
}

// Partition splits a collection (slice or array) into two parts based on a condition function.
//
// This function iterates through the elements of the input collection, applying the provided `condition` function to each element.
// It creates two separate collections: one containing elements for which the condition returns `true`, and the other containing elements
// for which the condition returns `false`. The function returns both coll.
//
// Parameters:
//   - `collection`: A slice or array of any type to partition.
//   - `condition`: A function that checks whether an element should go into the "true" partition or the "false" partition.
//     It returns `true` if the element should go into the "true" partition, and `false` otherwise.
//
// Returns:
//   - A tuple containing two new collections:
//   - The first collection contains elements for which the condition is `true`.
//   - The second collection contains elements for which the condition is `false`.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5, 6}
//	truePartition, falsePartition := Partition(numbers, func(value interface{}) bool {
//		return value.(int) % 2 == 0 // Partition into even and odd numbers
//	})
//	// truePartition will be []int{2, 4, 6} (even numbers)
//	// falsePartition will be []int{1, 3, 5} (odd numbers)
//
// Notes:
//   - The condition function is applied to each element in the collection.
//   - The returned collections are of the same type as the input collection (slice or array).
func Partition(collection any, condition func(value any) bool) (any, any) {
	v := reflect.ValueOf(collection)
	truePartition := reflect.MakeSlice(v.Type(), 0, 0)
	falsePartition := reflect.MakeSlice(v.Type(), 0, 0)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i).Interface()
			if condition(item) {
				truePartition = reflect.Append(truePartition, reflect.ValueOf(item))
			} else {
				falsePartition = reflect.Append(falsePartition, reflect.ValueOf(item))
			}
		}
	}
	return truePartition.Interface(), falsePartition.Interface()
}

// Zip combines multiple collections (slices or arrays) element-wise into a new collection of tuples.
//
// This function takes multiple collections (slices or arrays) as arguments and combines their elements
// into tuples. Each tuple consists of elements from the same index of each collection. The function returns
// a new collection (a slice of tuples), where each tuple contains the elements from the input collections
// at the same position. The length of the resulting collection is determined by the shortest input collection.
//
// Parameters:
//   - `collections`: A variadic list of slices or arrays of any type to combine. All collections must be
//     of the same length or shorter, and they will be combined element-wise until the shortest collection's length is reached.
//
// Returns:
//   - A slice of tuples (slices), where each tuple contains the elements from the input collections at the same index.
//     The length of the returned slice is the same as the shortest input collection.
//
// Example:
//
//	numbers := []int{1, 2, 3}
//	strings := []string{"a", "b", "c"}
//	result := Zip(numbers, strings)
//	// result will be [][]interface{}{{1, "a"}, {2, "b"}, {3, "c"}}
//
// Notes:
//   - If any of the input collections is not a slice or array, the function returns `nil`.
//   - If the collections have different lengths, the function will combine elements up to the length of the shortest collection.
func Zip(collections ...any) []any {
	minLength := -1
	for _, collection := range collections {
		v := reflect.ValueOf(collection)
		if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
			return nil
		}
		if minLength == -1 || v.Len() < minLength {
			minLength = v.Len()
		}
	}
	result := make([]any, minLength)
	for i := 0; i < minLength; i++ {
		tuple := make([]any, len(collections))
		for j, collection := range collections {
			v := reflect.ValueOf(collection)
			tuple[j] = v.Index(i).Interface()
		}
		result[i] = tuple
	}
	return result
}

// RotateLeft rotates the elements of a collection (slice or array) to the left by a specified number of positions.
//
// This function takes a collection (slice or array) and rotates its elements to the left by the given number
// of positions. Elements at the beginning of the collection are moved to the end. The rotation is done in-place
// (modifying the collection), and if the number of positions is negative, the rotation will be adjusted accordingly.
//
// Parameters:
//   - `collection`: A slice or array of any type to rotate.
//   - `positions`: The number of positions to rotate the collection to the left. A positive number rotates
//     elements left, while a negative number rotates right. The positions are normalized to be within the valid range
//     of the collection's length.
//
// Returns:
//   - A new collection (slice or array) where the elements have been rotated left by the specified number of positions.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	result := RotateLeft(numbers, 2)
//	// result will be []int{3, 4, 5, 1, 2}, as the collection is rotated left by 2 positions.
//
// Notes:
//   - If the number of positions is larger than the length of the collection, it is normalized using modulo
//     to ensure it rotates only the necessary number of positions.
//   - If the collection is not a slice or array, the original collection is returned unchanged.
func RotateLeft(collection any, positions int) any {
	v := reflect.ValueOf(collection)
	length := v.Len()
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		if positions < 0 {
			positions = (positions%length + length) % length
		} else {
			positions = positions % length
		}
		result := reflect.MakeSlice(v.Type(), length, length)
		for i := 0; i < length; i++ {
			result.Index((i - positions + length) % length).Set(v.Index(i))
		}
		return result.Interface()
	}
	return collection
}

// RotateRight rotates the elements of a collection (slice or array) to the right by a specified number of positions.
//
// This function takes a collection (slice or array) and rotates its elements to the right by the given number
// of positions. Elements at the end of the collection are moved to the beginning. The rotation is done in-place
// (modifying the collection), and if the number of positions is negative, the rotation will be adjusted accordingly.
//
// Parameters:
//   - `collection`: A slice or array of any type to rotate.
//   - `positions`: The number of positions to rotate the collection to the right. A positive number rotates
//     elements right, while a negative number rotates left. The positions are normalized to be within the valid range
//     of the collection's length.
//
// Returns:
//   - A new collection (slice or array) where the elements have been rotated right by the specified number of positions.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	result := RotateRight(numbers, 2)
//	// result will be []int{4, 5, 1, 2, 3}, as the collection is rotated right by 2 positions.
//
// Notes:
//   - If the number of positions is larger than the length of the collection, it is normalized using modulo
//     to ensure it rotates only the necessary number of positions.
//   - If the collection is not a slice or array, the original collection is returned unchanged.
func RotateRight(collection any, positions int) any {
	v := reflect.ValueOf(collection)
	length := v.Len()
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		if positions < 0 {
			positions = (-positions%length + length) % length
		} else {
			positions = positions % length
		}
		result := reflect.MakeSlice(v.Type(), length, length)
		for i := 0; i < length; i++ {
			result.Index((i + positions) % length).Set(v.Index(i))
		}
		return result.Interface()
	}
	return collection
}
