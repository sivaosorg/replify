package randn

import (
	"math/rand"
	"time"

	cr "crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

var r *rand.Rand // Package-level random generator

func init() {
	// Initialize the package-level random generator with a seed
	src := rand.NewSource(time.Now().UTC().UnixNano())
	r = rand.New(src)
}

// UUID generates a new universally unique identifier (UUID) using random data from /dev/urandom (Unix-based systems).
//
// This function opens the special file /dev/urandom to read 16 random bytes, which are then used to construct a UUID
// in the standard format (8-4-4-4-12 hex characters). It ensures that the file is properly closed after reading, even
// if an error occurs. If there's an error opening or reading from /dev/urandom, the function returns an appropriate error.
//
// UUID Format: The generated UUID is formatted as a string in the following structure:
// XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX, where X is a hexadecimal digit.
//
// Returns:
//   - A string representing the newly generated UUID.
//   - An error if there is an issue opening or reading from /dev/urandom.
//
// Example:
//
//	 uuid, err := UUID()
//		if err != nil {
//		    log.Fatalf("Failed to generate UUID: %v", err)
//		}
//	 fmt.Println("Generated UUID:", uuid)
//
// Notes:
//   - This function is designed for Unix-based systems. On non-Unix systems, this may not work because /dev/urandom
//     may not be available.
func UUID() (string, error) {
	dash := "-"
	return UUIDJoin(dash)
}

// UUIDJoin generates a new universally unique identifier (UUID) using random data from /dev/urandom
// (Unix-based systems) with a customizable delimiter.
//
// This function is similar to GenerateUUID but allows the user to specify a custom delimiter to separate
// different sections of the UUID. It opens the special file /dev/urandom to read 16 random bytes,
// which are then used to construct a UUID. The UUID is returned as a string in the format:
// XXXXXXXX<delimiter>XXXX<delimiter>XXXX<delimiter>XXXX<delimiter>XXXXXXXXXXXX, where X is a hexadecimal digit.
//
// Parameters:
//   - delimiter: A string used to separate sections of the UUID. Common choices are "-" or "" (no delimiter).
//
// Returns:
//   - A string representing the newly generated UUID with the specified delimiter.
//   - An error if there is an issue opening or reading from /dev/urandom.
//
// Example:
//
//	uuid, err := UUIDJoin("-")
//	if err != nil {
//	    log.Fatalf("Failed to generate UUID: %v", err)
//	}
//	fmt.Println("Generated UUID:", uuid)
//
// Notes:
//   - This function is designed for Unix-based systems. On non-Unix systems, it may not work because /dev/urandom
//     may not be available.
func UUIDJoin(delimiter string) (string, error) {
	file, err := os.Open("/dev/urandom")
	if err != nil {
		return "", fmt.Errorf("failed to open /dev/urandom: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Error closing file: %s\n", err)
		}
	}()
	b := make([]byte, 16)
	_, err = file.Read(b)
	if err != nil {
		return "", err
	}
	// Format the bytes as a UUID string with the specified delimiter.
	// The UUID is structured as XXXXXXXX<delimiter>XXXX<delimiter>XXXX<delimiter>XXXX<delimiter>XXXXXXXXXXXX.
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
// The function uses a custom random source seeded with the current Unix timestamp
// in nanoseconds to ensure that each call produces a unique sequence.
// This function is intended to generate random strings quickly and is not
// cryptographically secure.
//
// Example:
//
//	id := RandID(16)
//	fmt.Println("Generated Random ID:", id)
//
// Notes:
//   - This function is suitable for use cases where simple random IDs are needed.
//     However, for cryptographic purposes, consider using more secure random generation.
func RandID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano())) // Create a seeded random generator for unique results each call
	// Allocate a byte slice for the generated ID and populate it with random characters
	id := make([]byte, length)
	for i := range id {
		id[i] = charset[seededRand.Intn(len(charset))]
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
	bytes := make([]byte, 16)
	// Use crypto/rand.Read for cryptographically secure random byte generation.
	if _, err := cr.Read(bytes); err != nil {
		log.Fatalf("Failed to generate secure random bytes: %v", err)
		return ""
	}
	return hex.EncodeToString(bytes)
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

// RandUUID generates and returns a new UUID using the GenerateUUID function.
//
// If an error occurs during UUID generation (for example, if there is an issue reading from /dev/urandom),
// the function returns an empty string.
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

// RandIntr generates a random integer within the specified range, inclusive,
// and reseeds the random number generator with a new value each time it is called.
//
// The function first creates a new seed by combining the current UTC time in nanoseconds
// with a random integer from the custom random generator `r`. This helps to ensure that
// each call to `RandIntr` produces a different sequence of random numbers.
//
// If the provided `min` is greater than or equal to `max`, the function returns `min`
// as a default value. This behavior should be considered when using this function,
// as it may not be intuitive to return `min` in cases where the range is invalid.
//
// The function uses the reseeded random generator to generate a random integer in the
// inclusive range from `min` to `max`. The calculation ensures that both bounds are included
// in the result by using the formula: r.Intn(max-min+1) + min.
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
//	fmt.Println("Random number between 1 and 10 after reseeding:", randomNum)
func RandIntr(min, max int) int {
	// Reseed the custom random generator with a new seed
	x := time.Now().UTC().UnixNano() + int64(r.Int())
	r.Seed(x)
	// Ensure max is included in the range
	if min >= max {
		return min
	}
	return r.Intn(max-min+1) + min
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
