// Package coll provides generic collection types and functional utilities for
// working with slices and maps in Go.
//
// The package is organised into three layers:
//
//  1. Concrete collection types – HashMap, HashSet, and Stack – each backed
//     by the corresponding standard Go data structure and exposed through a
//     typed, method-based API.
//
//  2. Map utilities – standalone functions for transforming, merging,
//     filtering, flattening, and inverting maps with full type-parameter
//     support.
//
//  3. Slice utilities – higher-order functions such as Map, Filter, Reduce,
//     Chunk, and ZipWith that operate on slices of any element type.
//
// # Generic Collection Types
//
// HashMap[K, V comparable] is a typed wrapper around Go's built-in map that
// provides Put, Get, Remove, ContainsKey, KeySet, and Clear operations:
//
//	m := coll.NewHashMap[string, int]()
//	m.Put("hits", 42)
//	fmt.Println(m.Get("hits")) // 42
//
// HashSet[T comparable] stores unique elements and supports set-algebra
// operations including Intersection, Union, and Difference:
//
//	s := coll.NewHashSet("a", "b", "c")
//	s.Add("d")
//	fmt.Println(s.Contains("b")) // true
//
// Stack[T comparable] implements a last-in, first-out (LIFO) structure with
// Push, Pop, Peek, and IsEmpty:
//
//	stack := coll.NewStack[int]()
//	stack.Push(1)
//	stack.Push(2)
//	fmt.Println(stack.Pop()) // 2
//
// # Map Utilities
//
// MergeComp, DeepMerge, FlattenMap, UnflattenMap, PickComp, OmitComp, and
// InvertComp cover the most common map transformation patterns. Functions
// whose names end in Comp accept keys constrained to comparable, while their
// non-Comp counterparts accept map[any]V for looser typing.
//
// # Slice Utilities
//
// Map and ToSlice transform each element, ToMap indexes a slice by a key
// function, and GetOrDefault provides a safe fallback for map lookups.
//
// All functions are safe for concurrent use if the underlying collections are
// not mutated during the call. The collection types themselves are not
// goroutine-safe; external synchronisation is required when sharing instances
// across goroutines.
package coll
