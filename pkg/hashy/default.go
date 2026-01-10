package hashy

import "hash/fnv"

// DefaultOptions returns default options for hashing.
//
// Returns:
//   - A pointer to a newly created `Options` instance with the default values.
func DefaultOptions() *hashOptions {
	return &hashOptions{
		Hasher:  fnv.New64a(),
		TagName: "hash",
	}
}

// NewOptions creates a new options builder with defaults.
//
// Returns:
//   - A pointer to a newly created `OptionsBuilder` instance with the default values.
func NewOptions() *OptionsBuilder {
	return &OptionsBuilder{
		opts: hashOptions{
			Hasher:  fnv.New64a(),
			TagName: "hash",
		},
	}
}
