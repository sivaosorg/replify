package codegen

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

// generate is the internal function that performs random code generation.
// It uses crypto/rand.Int to provide cryptographically secure randomness.
// No mutex is required because crypto/rand is already thread-safe.
func generate(opts Options) (string, error) {
	chars := []rune(string(opts.Charset))
	charsetSize := big.NewInt(int64(len(chars)))

	var sb strings.Builder
	sb.Grow(len(opts.Prefix) + opts.Length + len(opts.Suffix))
	sb.WriteString(opts.Prefix)

	for i := 0; i < opts.Length; i++ {
		idx, err := rand.Int(rand.Reader, charsetSize)
		if err != nil {
			return "", fmt.Errorf("codegen: crypto/rand failed: %w", err)
		}
		sb.WriteRune(chars[idx.Int64()])
	}

	sb.WriteString(opts.Suffix)
	return sb.String(), nil
}

// validateOptions verifies that the provided Options are valid.
// It returns the first validation error encountered.
func validateOptions(o Options) error {
	if o.Length < 1 {
		return ErrInvalidLength
	}
	if len([]rune(string(o.Charset))) == 0 {
		return ErrEmptyCharset
	}
	return nil
}

// deduplicateRunes removes duplicate characters from a string while
// preserving the order of their first occurrence. Supports Unicode.
func deduplicateRunes(s string) string {
	seen := make(map[rune]struct{})
	runes := make([]rune, 0, len([]rune(s)))
	for _, r := range s {
		if _, ok := seen[r]; !ok {
			seen[r] = struct{}{}
			runes = append(runes, r)
		}
	}
	return string(runes)
}
