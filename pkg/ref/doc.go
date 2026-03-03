// Package ref provides generic pointer and nil utilities that reduce
// boilerplate when working with optional values and pointer-based APIs.
//
// # Creating Pointers
//
// Ptr takes any value and returns a pointer to it. This eliminates the need
// for a temporary variable when you need to take the address of a literal or
// a function return value:
//
//	// Without ref.Ptr
//	name := "Alice"
//	req.Name = &name
//
//	// With ref.Ptr
//	req.Name = ref.Ptr("Alice")
//
// # Dereferencing Safely
//
// Deref dereferences a pointer and returns its value. When the pointer is
// nil, Deref returns the zero value of the pointed-to type instead of
// panicking:
//
//	var p *int
//	fmt.Println(ref.Deref(p)) // 0
//
// # Nil Checking
//
// IsNil reports whether a value is nil. Unlike a direct nil comparison,
// IsNil handles interface values that wrap a typed nil pointer, which would
// otherwise compare as non-nil:
//
//	var err error = (*os.PathError)(nil)
//	ref.IsNil(&err) // true
//
// The assert package uses ref.IsNil internally to implement AssertNil and
// AssertNotNil in a way that correctly handles interface boxing.
//
// All functions in this package are pure and safe for concurrent use.
package ref
