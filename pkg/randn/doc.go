// Package randn provides functions for generating random values, unique
// identifiers, and universally unique identifiers (UUIDs).
//
// UUID and UUIDSep use crypto/rand to produce RFC 4122 version-4 UUIDs and
// are safe on all platforms (Linux, macOS, Windows, and others).
// RandID and the numeric functions rely on the goroutine-safe global math/rand
// source (automatically seeded in Go 1.20+) and are not cryptographically
// secure; use CryptoID when security properties matter.
//
// # Identifiers
//
//	randn.UUID()         // standard UUID, e.g. "550e8400-e29b-41d4-a716-446655440000"
//	randn.UUIDSep("/")  // UUID with a custom delimiter
//	randn.RandUUID()     // UUID, empty string on error
//	randn.RandID(16)     // 16-character alphanumeric string (math/rand)
//	randn.CryptoID()     // 32-character hex string (crypto/rand)
//	randn.TimeID()       // nanosecond timestamp + random int, as a string
//	randn.NewXID()       // unique 20-character identifier (rs/xid port)
//
// # Numeric Values
//
//	randn.RandInt()            // random int
//	randn.RandIntr(min, max)   // random int in [min, max] inclusive
//	randn.RandFt64()           // random float64 in [0.0, 1.0)
//	randn.RandFt64r(lo, hi)    // random float64 in [lo, hi)
//	randn.RandFt32()           // random float32 in [0.0, 1.0)
//	randn.RandFt32r(lo, hi)    // random float32 in [lo, hi)
//	randn.RandByte(n)          // slice of n random bytes
//
// All functions in this package are safe for concurrent use.
package randn
