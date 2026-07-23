package codegen

import "sync"

// Charset represents the set of characters used for random code generation.
// You can use one of the predefined constants or define your own character
// set with WithCustomCharset.
type Charset string

// Options contains the complete configuration for a Generator.
// All fields have valid default values provided by defaultOptions.
type Options struct {
	// Length is the number of random characters in each generated code,
	// excluding the Prefix and Suffix.
	Length int

	// Charset is the character set used for random code generation.
	Charset Charset

	// Prefix is a static string prepended to every generated code.
	Prefix string

	// Suffix is a static string appended to every generated code.
	Suffix string
}

// Generator is the primary type for generating random codes.
// It uses crypto/rand to provide cryptographically secure randomness,
// making it suitable for order codes in fulfillment systems.
//
// Generator is safe for concurrent use by multiple goroutines.
// A single instance can be shared across the entire application.
type Generator struct {
	mu   sync.Mutex
	opts Options
}
