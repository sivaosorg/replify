package codegen

// Option is a functional option used to configure a Generator.
// Use the WithXxx helper functions to create Options and pass them
// to New or SetOptions.
type Option func(*Options)

// defaultOptions returns the default configuration:
//   - Length: 8
//   - Charset: CharsetAlphanumeric
//   - Prefix: "" (empty)
//   - Suffix: "" (empty)
func defaultOptions() Options {
	return Options{
		Length:  8,
		Charset: CharsetAlphanumeric,
	}
}

// WithLength sets the number of random characters for each generated code.
// The value must be greater than 0; otherwise, New and SetOptions will
// return ErrInvalidLength.
//
// Example:
//
//	g, _ := codegen.New(codegen.WithLength(12))
func WithLength(length int) Option {
	return func(o *Options) {
		o.Length = length
	}
}

// WithCharset sets the character set used for random code generation.
// It is recommended to use one of the predefined Charset constants
// provided by this package.
//
// Example:
//
//	g, _ := codegen.New(codegen.WithCharset(codegen.CharsetNumeric))
func WithCharset(charset Charset) Option {
	return func(o *Options) {
		o.Charset = charset
	}
}

// WithCustomCharset sets a custom character set from an arbitrary string.
// Duplicate characters are removed to ensure a uniform distribution.
//
// Example:
//
//	// Use only characters that are easy to distinguish visually.
//	g, _ := codegen.New(codegen.WithCustomCharset("23456789ABCDEFGHJKLMNPQRSTUVWXYZ"))
func WithCustomCharset(chars string) Option {
	return func(o *Options) {
		o.Charset = Charset(deduplicateRunes(chars))
	}
}

// WithPrefix sets the static string prepended to every generated code.
// The Prefix is not included in Length.
//
// Example:
//
//	g, _ := codegen.New(codegen.WithPrefix("ORD-"))
//	// Generates: "ORD-A3BF9KP2"
func WithPrefix(prefix string) Option {
	return func(o *Options) {
		o.Prefix = prefix
	}
}

// WithSuffix sets the static string appended to every generated code.
// The Suffix is not included in Length.
//
// Example:
//
//	g, _ := codegen.New(codegen.WithSuffix("-VN"))
//	// Generates: "A3BF9KP2-VN"
func WithSuffix(suffix string) Option {
	return func(o *Options) {
		o.Suffix = suffix
	}
}
