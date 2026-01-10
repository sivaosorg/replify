package coll

import (
	"fmt"

	"github.com/sivaosorg/replify/pkg/strutil"
)

// itemExists is an empty struct used as a placeholder for values in the `HashSet` map.
// Using an empty struct saves memory, as it does not allocate any space in Go.
var itemExists struct{}

// HashSet is a generic set data structure that stores unique elements of type `T`.
// It is implemented using a map where the keys are the elements and the values are
// empty structs (`struct{}`), which is an efficient way to represent a set in Go.
// The type `T` must be comparable, meaning it can be compared using == and !=.
type HashSet[T comparable] struct {
	items map[T]struct{}
}

// NewHashSet is a constructor function that creates a new `HashSet` and optionally adds
// any provided elements to the set.
// Parameters:
//   - `elements`: A variadic list of elements of type `T` to initialize the set with (optional).
//
// Returns:
//   - A pointer to the newly created `HashSet`.
//
// Example usage:
//
//	hashSet := NewHashSet(1, 2, 3) // Creates a set with initial elements 1, 2, and 3.
func NewHashSet[T comparable](elements ...T) *HashSet[T] {
	hash := &HashSet[T]{
		items: make(map[T]struct{}),
	}
	hash.AddAll(elements...)
	return hash
}

// Add inserts a new element into the HashSet. If the element already exists, no action is taken.
// Parameters:
//   - `element`: The element to be added to the set.
//
// Example usage:
//
//	hashSet.Add(4) // Adds the element 4 to the set.
func (hash *HashSet[T]) Add(element T) {
	hash.items[element] = itemExists
}

// AddAll inserts multiple elements into the HashSet. If any element already exists, it is ignored.
// Parameters:
//   - `elements`: A variadic list of elements of type `T` to add to the set.
//
// Example usage:
//
//	hashSet.AddAll(4, 5, 6) // Adds elements 4, 5, and 6 to the set.
func (hash *HashSet[T]) AddAll(elements ...T) {
	for _, e := range elements {
		hash.Add(e)
	}
}

// Remove deletes the specified element from the HashSet.
// If the element does not exist, no action is taken.
// Parameters:
//   - `element`: The element to be removed from the set.
//
// Example usage:
//
//	hashSet.Remove(5) // Removes element 5 from the set.
func (hash *HashSet[T]) Remove(element T) {
	delete(hash.items, element)
}

// RemoveAll deletes multiple elements from the HashSet.
// Parameters:
//   - `elements`: A variadic list of elements to be removed from the set.
//
// Example usage:
//
//	hashSet.RemoveAll(2, 4, 6) // Removes elements 2, 4, and 6 from the set.
func (hash *HashSet[T]) RemoveAll(elements ...T) {
	for _, e := range elements {
		hash.Remove(e)
	}
}

// Clear removes all elements from the HashSet, resetting it to an empty set.
//
// Example usage:
//
//	hashSet.Clear() // Clears all elements from the set.
func (hash *HashSet[T]) Clear() {
	hash.items = make(map[T]struct{})
}

// Size returns the number of elements currently stored in the HashSet.
// Returns:
//   - The number of elements in the set.
//
// Example usage:
//
//	size := hashSet.Size() // Returns the size of the set.
func (hash *HashSet[T]) Size() int {
	return len(hash.items)
}

// IsEmpty checks if the HashSet contains any elements.
// Returns:
//   - `true` if the set is empty, `false` otherwise.
//
// Example usage:
//
//	isEmpty := hashSet.IsEmpty() // Returns true if the set is empty.
func (hash *HashSet[T]) IsEmpty() bool {
	return len(hash.items) == 0
}

// Contains checks whether the specified element exists in the HashSet.
// Parameters:
//   - `element`: The element to check for existence in the set.
//
// Returns:
//   - `true` if the element is present in the set, `false` otherwise.
//
// Example usage:
//
//	exists := hashSet.Contains(3) // Checks if element 3 is in the set.
func (hash *HashSet[T]) Contains(element T) bool {
	_, ok := hash.items[element]
	return ok
}

// Intersection returns a new HashSet containing only the elements present in both the current set and another set.
// Parameters:
//   - `another`: Another `HashSet` to find the intersection with.
//
// Returns:
//   - A new `HashSet` containing the intersection of the two sets.
//
// Example usage:
//
//	resultSet := hashSet.Intersection(anotherSet) // Returns the intersection of two sets.
func (hash *HashSet[T]) Intersection(another *HashSet[T]) *HashSet[T] {
	result := NewHashSet[T]()
	if hash.Size() <= another.Size() {
		for item := range hash.items {
			if _, contains := another.items[item]; contains {
				result.Add(item)
			}
		}
		return result
	}
	for item := range another.items {
		if _, contains := hash.items[item]; contains {
			result.Add(item)
		}
	}
	return result
}

// Union returns a new HashSet that contains all elements from both the current set and another set.
// Parameters:
//   - `another`: Another `HashSet` to combine with the current set.
//
// Returns:
//   - A new `HashSet` containing the union of the two sets.
//
// Example usage:
//
//	resultSet := hashSet.Union(anotherSet) // Returns the union of two sets.
func (hash *HashSet[T]) Union(another *HashSet[T]) *HashSet[T] {
	result := NewHashSet[T]()
	for item := range hash.items {
		result.Add(item)
	}
	for item := range another.items {
		result.Add(item)
	}
	return result
}

// Difference returns a new HashSet containing the elements that are in the current set but not in another set.
// Parameters:
//   - `another`: Another `HashSet` to compare against.
//
// Returns:
//   - A new `HashSet` containing the difference between the current set and the other set.
//
// Example usage:
//
//	resultSet := hashSet.Difference(anotherSet) // Returns the difference between two sets.
func (hash *HashSet[T]) Difference(another *HashSet[T]) *HashSet[T] {
	result := NewHashSet[T]()
	for item := range hash.items {
		if _, exist := another.items[item]; !exist {
			result.Add(item)
		}
	}
	return result
}

// Slice converts the HashSet into a slice of elements.
// Returns:
//   - A slice containing all elements in the set.
//
// Example usage:
//
//	slice := hashSet.Slice() // Converts the set to a slice.
func (hash *HashSet[T]) Slice() []T {
	slices := make([]T, hash.Size())
	i := 0
	for val := range hash.items {
		slices[i] = val
		i++
	}
	return slices
}

// String returns a string representation of the HashSet, with elements separated by commas.
// Returns:
//   - A string that represents the elements in the set.
//
// Example usage:
//
//	str := hashSet.String() // Returns a string representation of the set.
func (hash *HashSet[T]) String() string {
	s := ""
	for val := range hash.items {
		if strutil.IsEmpty(s) {
			s = fmt.Sprintf("%v", val)
		} else {
			s = fmt.Sprintf("%s,%v", s, val)
		}
	}
	return s
}
