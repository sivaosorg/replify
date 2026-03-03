// Package common provides shared runtime utilities used across the replify
// packages. It covers three broad areas: deep equality comparisons, I/O
// helper functions, and generic collection transformations.
//
// # Deep Equality
//
// DeepEqualComp uses reflect.DeepEqual to compare two values of the same
// comparable type. DeepEqual serialises both values to JSON and compares the
// resulting byte slices, making it a convenient option for loosely typed
// comparisons where field ordering in maps is not a concern.
//
// IsScalarType reports whether a value is one of Go's built-in primitive
// types (bool, numeric, or string). IsEmptyValue extends this to any
// reflect.Value, returning true when the value is a zero length collection,
// false boolean, zero numeric, nil pointer or interface, or zero struct.
//
// # I/O Helpers
//
// ReadAll and SlurpAll both drain an [io.Reader] into a string; ReadAll
// delegates to [io.Copy] while SlurpAll reads in 1 KiB chunks. SlurpLines
// returns each line as a separate element of a string slice, and SlurpLine
// concatenates all lines into a single string.
//
// TeeCopy reads from an [io.Reader] and writes to an [io.Writer]
// simultaneously, returning the data as a string. TeeTap is similar but
// discards the writer side, effectively logging a stream without
// forwarding it.
//
// # Generic Transformations
//
// Transform applies a predicate to every element of a slice, array, or map
// using reflection, returning a new slice of the mapped values. A rich set
// of higher-order slice functions — Filter, Reduce, ForEach, GroupBy,
// Partition, Chunk, Unique, Flatten, ZipWith, Any, All, and None — round out
// the functional programming toolkit available to replify consumers.
//
// All functions in this package are stateless and safe for concurrent use.
package common
