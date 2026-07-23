package codegen

import (
	"fmt"
)

// New creates and returns a new Generator with the provided options.
// If no options are specified, the following defaults are applied:
//   - Length:  8
//   - Charset: CharsetAlphanumeric
//   - Prefix:  "" (empty)
//   - Suffix:  "" (empty)
//
// Returns an error if any option is invalid:
//   - ErrInvalidLength: if Length < 1
//   - ErrEmptyCharset:  if Charset is empty
//
// Example:
//
//	g, err := codegen.New(
//	    codegen.WithLength(12),
//	    codegen.WithCharset(codegen.CharsetAlphanumericUpper),
//	    codegen.WithPrefix("ORD-"),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
func New(opts ...Option) (*Generator, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(&o)
	}
	if err := validateOptions(o); err != nil {
		return nil, err
	}
	return &Generator{opts: o}, nil
}

// MustNew behaves like New but panics if an error occurs.
// It is intended for initialization during application startup,
// where invalid configuration should be detected immediately (fail fast).
//
// Example:
//
//	var orderGen = codegen.MustNew(
//	    codegen.WithLength(12),
//	    codegen.WithCharset(codegen.CharsetAlphanumericUpper),
//	    codegen.WithPrefix("ORD-"),
//	)
func MustNew(opts ...Option) *Generator {
	g, err := New(opts...)
	if err != nil {
		panic(fmt.Sprintf("codegen: MustNew failed: %v", err))
	}
	return g
}

// Generate creates and returns a single random code using the current configuration.
// The total length of the returned string is:
// len(Prefix) + Length + len(Suffix).
//
// Uses crypto/rand to ensure unpredictability.
// Safe for concurrent use by multiple goroutines.
//
// Example:
//
//	code, err := g.Generate()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(code) // "ORD-A3BF9KP2XQ17"
func (g *Generator) Generate() (string, error) {
	// Copy the options under lock to avoid holding the lock
	// during the entire generation process.
	g.mu.Lock()
	opts := g.opts
	g.mu.Unlock()

	return generate(opts)
}

// GenerateN creates and returns a slice containing n independently
// generated random codes. Since each code is generated independently,
// duplicates are theoretically possible, although the probability is
// extremely low when using a sufficiently large charset and length.
//
// Returns ErrInvalidCount if n < 1.
// Safe for concurrent use by multiple goroutines.
//
// Example:
//
//	codes, err := g.GenerateN(100)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// codes contains a slice of 100 order codes
func (g *Generator) GenerateN(n int) ([]string, error) {
	if n < 1 {
		return nil, ErrInvalidCount
	}

	g.mu.Lock()
	opts := g.opts
	g.mu.Unlock()

	codes := make([]string, n)
	for i := range codes {
		code, err := generate(opts)
		if err != nil {
			return nil, err
		}
		codes[i] = code
	}
	return codes, nil
}

// SetOptions atomically updates the Generator configuration.
// If any provided option is invalid, the existing configuration
// remains unchanged and an error is returned (no partial update).
//
// Safe for concurrent use by multiple goroutines.
//
// Example:
//
//	err := g.SetOptions(
//	    codegen.WithLength(16),
//	    codegen.WithPrefix("INV-"),
//	)
func (g *Generator) SetOptions(opts ...Option) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Create a copy first for validation to avoid partial updates
	// if an error occurs.
	candidate := g.opts
	for _, opt := range opts {
		opt(&candidate)
	}
	if err := validateOptions(candidate); err != nil {
		return err
	}
	g.opts = candidate
	return nil
}

// Options returns a snapshot copy of the Generator's current configuration.
// Modifying the returned value does not affect the Generator.
//
// Safe for concurrent use by multiple goroutines.
func (g *Generator) Options() Options {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.opts
}

// Generate is a package-level convenience function that creates a temporary
// Generator with the provided options and generates a single code.
//
// It is suitable for one-off code generation.
// If you need to generate multiple codes, create a Generator with New
// and reuse it instead.
//
// Example:
//
//	code, err := codegen.Generate(
//	    codegen.WithLength(10),
//	    codegen.WithCharset(codegen.CharsetNumeric),
//	)
func Generate(opts ...Option) (string, error) {
	g, err := New(opts...)
	if err != nil {
		return "", err
	}
	return g.Generate()
}
