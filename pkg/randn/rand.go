package randn

import (
	"math/rand"
	"time"

	cr "crypto/rand"
	"encoding/hex"
	"fmt"
)

// UUID generates a new universally unique identifier (UUID) (RFC 4122 version 4)
// using crypto/rand for cryptographically secure randomness.
//
// UUID Format: XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX (version 4, variant 1).
//
// Returns:
//   - A string representing the newly generated UUID.
//   - An error if the random source is unavailable.
//
// Example:
//
//	uuid, err := UUID()
//	if err != nil {
//	    log.Fatalf("Failed to generate UUID: %v", err)
//	}
//	fmt.Println("Generated UUID:", uuid)
func UUID() (string, error) {
	return UUIDSep("-")
}

// UUIDSep generates a new universally unique identifier (UUID) (RFC 4122 version 4)
// using crypto/rand, with a customizable delimiter between UUID sections.
//
// This function is cross-platform: it does not rely on /dev/urandom or any
// OS-specific file, and works correctly on Linux, macOS, Windows, and other
// platforms supported by the Go crypto/rand package.
//
// Parameters:
//   - delimiter: A string used to separate sections of the UUID. Common choices are "-" or "".
//
// Returns:
//   - A string representing the newly generated UUID with the specified delimiter.
//   - An error if crypto/rand fails to produce random bytes.
//
// Example:
//
//	uuid, err := UUIDSep("-")
//	if err != nil {
//	    log.Fatalf("Failed to generate UUID: %v", err)
//	}
//	fmt.Println("Generated UUID:", uuid)
func UUIDSep(delimiter string) (string, error) {
	b := make([]byte, 16)
	if _, err := cr.Read(b); err != nil {
		return "", fmt.Errorf("randn: failed to generate UUID random bytes: %w", err)
	}
	// Set RFC 4122 version 4 and variant 1 bits.
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant 1
	uuid := fmt.Sprintf("%x%s%x%s%x%s%x%s%x", b[0:4], delimiter, b[4:6], delimiter, b[6:8], delimiter, b[8:10], delimiter, b[10:])
	return uuid, nil
}

// RandID generates a random alphanumeric string of the specified length.
// This string includes uppercase letters, lowercase letters, and numbers, making it
// suitable for use as unique IDs or tokens.
//
// Parameters:
//   - length: The length of the random ID to generate. Must be a positive integer.
//
// Returns:
//   - A string of random alphanumeric characters with the specified length.
//
// The function uses the goroutine-safe global math/rand source (automatically seeded
// in Go 1.20+). This function is intended to generate random strings quickly and is
// not cryptographically secure.
//
// Example:
//
//	id := RandID(16)
//	fmt.Println("Generated Random ID:", id)
//
// Notes:
//   - This function is suitable for use cases where simple random IDs are needed.
//     However, for cryptographic purposes, consider using CryptoID instead.
func RandID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	id := make([]byte, length)
	for i := range id {
		id[i] = charset[rand.Intn(len(charset))]
	}
	return string(id)
}

// CryptoID generates a cryptographically secure random ID as a hexadecimal string.
// It uses 16 random bytes, which are then encoded to a hexadecimal string for easy representation.
//
// Returns:
//   - A string representing a secure random hexadecimal ID of 32 characters (since 16 bytes are used, and each byte
//     is represented by two hexadecimal characters).
//
// The function uses crypto/rand.Read to ensure cryptographic security in the generated ID, making it suitable for
// sensitive use cases such as API keys, session tokens, or any security-critical identifiers.
//
// Example:
//
//	id := CryptoID()
//	fmt.Println("Generated Crypto ID:", id)
//
// Notes:
//   - This function is suitable for use cases where high security is required in the generated ID.
//   - It is not recommended for use cases where deterministic or non-cryptographic IDs are preferred.
func CryptoID() string {
	b := make([]byte, 16)
	// Use crypto/rand.Read for cryptographically secure random byte generation.
	if _, err := cr.Read(b); err != nil {
		panic(fmt.Sprintf("randn: failed to generate secure random bytes: %v", err))
	}
	return hex.EncodeToString(b)
}

// TimeID generates a unique identifier based on the current Unix timestamp in nanoseconds,
// with an additional random integer to enhance uniqueness.
//
// This function captures the current time in nanoseconds since the Unix epoch and appends a random integer
// to ensure additional randomness and uniqueness, even if called in rapid succession. The result is returned
// as a string. This type of ID is well-suited for time-based ordering and can be useful for generating
// unique identifiers for logs, events, or non-cryptographic applications.
//
// Returns:
//   - A string representing the current Unix timestamp in nanoseconds, concatenated with a random integer.
//
// Example:
//
//	id := TimeID()
//	fmt.Println("Generated Timestamp ID:", id)
//
// Notes:
//   - This function provides a unique, time-ordered identifier, but it is not suitable for cryptographic use.
//   - The combination of the current time and a random integer is best suited for applications requiring
//     uniqueness and ordering, rather than secure identifiers.
func TimeID() string {
	return fmt.Sprintf("%d%d", time.Now().UnixNano(), rand.Int())
}

// RandUUID generates and returns a new UUID (RFC 4122 version 4).
//
// If an error occurs during UUID generation, the function returns an empty string.
//
// This function is useful when you want a simple UUID generation without handling errors directly.
// It abstracts away the error handling by returning an empty string in case of failure.
//
// Returns:
//   - A string representing the newly generated UUID.
//   - An empty string if an error occurs during UUID generation.
//
// Example:
//
//	uuid := RandUUID()
//
//	if uuid == "" {
//	    fmt.Println("Failed to generate UUID")
//	} else {
//
//	    fmt.Println("Generated UUID:", uuid)
//	}
func RandUUID() string {
	v, err := UUID()
	if err != nil {
		return ""
	}
	return v
}

// RandInt returns the next random int value.
//
// This function uses the rand package to generate a random int value.
// Returns:
//   - A random int value.
func RandInt() int {
	return rand.Int()
}

// RandIntr generates a random integer within the specified range [min, max], inclusive.
//
// If the provided min is greater than or equal to max, the function returns min
// as a default value.
//
// The function uses the goroutine-safe global math/rand source (automatically seeded
// in Go 1.20+) and is safe for concurrent use without external synchronization.
//
// Parameters:
//   - `min`: The lower bound of the random number range (inclusive).
//   - `max`: The upper bound of the random number range (inclusive).
//
// Returns:
//   - A random integer between `min` and `max`, including both bounds.
//
// Example:
//
//	randomNum := RandIntr(1, 10)
//	fmt.Println("Random number between 1 and 10:", randomNum)
func RandIntr(min, max int) int {
	if min >= max {
		return min
	}
	return rand.Intn(max-min+1) + min
}

// RandUint32 returns the next random uint32 value.
//
// This function uses the crypto/rand package to generate a random uint32 value.
//
// Returns:
//   - A random uint32 value.
//
// Example:
//
//	randomNum := RandUint32()
//	fmt.Println("Random number between 1 and 10 after reseeding:", randomNum)
func RandUint32() uint32 {
	b := make([]byte, 3)
	if _, err := cr.Reader.Read(b); err != nil {
		panic(fmt.Errorf("randn: cannot generate random number: %v;", err))
	}
	return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
}

// RandFt64 returns the next random float64 value in the range [0.0, 1.0).
//
// This function uses the rand package to generate a random float64 value.
// The generated value is uniformly distributed over the interval [0.0, 1.0).
//
// Returns:
//   - A random float64 value between 0.0 and 1.0.
func RandFt64() float64 {
	return rand.Float64()
}

// RandFt64r returns the next random float64 value bounded by the specified range.
//
// Parameters:
//   - `start`: The lower bound of the random float64 value (inclusive).
//   - `end`: The upper bound of the random float64 value (exclusive).
//
// Returns:
//   - A random float64 value uniformly distributed between `start` and `end`.
func RandFt64r(start float64, end float64) float64 {
	return rand.Float64()*(end-start) + start
}

// RandFt32 returns the next random float32 value in the range [0.0, 1.0).
//
// This function uses the rand package to generate a random float32 value.
// The generated value is uniformly distributed over the interval [0.0, 1.0).
//
// Returns:
//   - A random float32 value between 0.0 and 1.0.
func RandFt32() float32 {
	return rand.Float32()
}

// RandFt32r returns the next random float32 value bounded by the specified range.
//
// Parameters:
//   - `start`: The lower bound of the random float32 value (inclusive).
//   - `end`: The upper bound of the random float32 value (exclusive).
//
// Returns:
//   - A random float32 value uniformly distributed between `start` and `end`.
func RandFt32r(start float32, end float32) float32 {
	return rand.Float32()*(end-start) + start
}

// RandByte creates an array of random bytes with the specified length.
//
// Parameters:
//   - `count`: The number of random bytes to generate.
//
// Returns:
//   - A slice of random bytes with the specified length.
func RandByte(count int) []byte {
	a := make([]byte, count)
	for i := range a {
		a[i] = (byte)(rand.Int())
	}
	return a
}

// RandIDHex produces a cryptographically random identifier of the specified
// byte length, returned as a hex-encoded string. The resulting string length
// is 2x the requested byte count.
//
// This replaces the use of math/rand, providing collision-resistant
// identifiers suitable for concurrent operations.
//
// Example:
//
//	id := RandIDHex(8)
//	fmt.Println(len(id)) // 16
func RandIDHex(byteLen int) string {
	b := make([]byte, byteLen)
	if _, err := cr.Read(b); err != nil {
		// Fallback: this should never fail on modern systems, but if it does
		// we still need a unique-ish string to avoid panics.
		return fmt.Sprintf("fallback-%p", &b)
	}
	return hex.EncodeToString(b)
}
