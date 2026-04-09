// Package hashy provides deterministic, structural hashing of arbitrary Go
// values, including structs, slices, maps, and primitive types.
//
// The package is built around a configurable hasher that traverses a value
// using reflection, feeds each field and element into a 64-bit FNV-1a hash
// function, and returns a reproducible uint64 digest. The same logical value
// always produces the same hash within a single binary; the hash is not
// stable across different Go versions or architectures.
//
// # Basic Usage
//
//	h, err := hashy.Hash(myStruct)
//	fmt.Printf("%016x\n", h)
//
//	// Multiple values are hashed as a tuple:
//	h, err = hashy.Hash(userID, role, timestamp)
//
// # Output Formats
//
// Hash returns a raw uint64. Convenience wrappers encode the result in
// common formats:
//
//	Hash256(v)    → SHA-256 of the uint64, as a hex string
//	Hash16Padded(v)    → zero-padded 16-character hex string
//	Hash16(v) → hexadecimal string
//	Hash10(v) → decimal string
//	Hash32(v) → base-32 string
//	Hash64(v)→ base-64 string
//
// # Configuration
//
// Hash behaviour can be tuned by passing an *Options value (built via
// NewOptions().WithTagName(...).WithZeroNil(true).Build()) as the final
// variadic argument:
//
//	opts := hashy.NewOptions().WithSlicesAsSets(true).Build()
//	h, err := hashy.Hash(mySlice, opts)
//
// Notable options include ZeroNil (treat nil pointers as zero values),
// IgnoreZeroValue (omit zero-value fields from the hash), SlicesAsSets
// (order-independent slice hashing), and UseStringer (use fmt.Stringer when
// available). The TagName option controls which struct tag is inspected for
// per-field directives such as "ignore" or "set".
//
// Structs may implement the Hashable, FieldSelector, or MapSelector
// interfaces to customise how they are hashed.
//
// hashy is safe for concurrent use when options are nil or were built with
// WithHasherFunc. Sharing options built with WithHasher across goroutines
// causes a data race because a single hash.Hash64 instance is not goroutine-safe.
package hashy
