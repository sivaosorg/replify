// Package randn provides functions for generating random values, unique
// identifiers, and universally unique identifiers (UUIDs).
//
// The package is initialised with a package-level *rand.Rand seeded from the
// current wall-clock time, giving good statistical randomness for
// non-security-critical workloads. For cryptographically strong identifiers,
// the CryptoID function reads from crypto/rand.
//
// # Identifiers
//
//	randn.UUID()         // standard UUID, e.g. "550e8400-e29b-41d4-a716-446655440000"
//	randn.UUIDJoin("/")  // UUID with a custom delimiter
//	randn.RandUUID()     // UUID, empty string on error
//	randn.RandID(16)     // 16-character alphanumeric string (math/rand)
//	randn.CryptoID()     // 32-character hex string (crypto/rand)
//	randn.TimeID()       // nanosecond timestamp + random int, as a string
//	randn.NewXID()       // unique 20-character identifier (rs/xid port)
//
// # Numeric Values
//
//	randn.RandInt()            // random int
//	randn.RandIntr(min, max)   // random int in [min, max] inclusive; reseeds on each call
//	randn.RandFt64()           // random float64 in [0.0, 1.0)
//	randn.RandFt64r(lo, hi)    // random float64 in [lo, hi)
//	randn.RandFt32()           // random float32 in [0.0, 1.0)
//	randn.RandFt32r(lo, hi)    // random float32 in [lo, hi)
//	randn.RandByte(n)          // slice of n random bytes
//
// UUID and UUIDJoin read from /dev/urandom and are therefore available only
// on Unix-like systems. RandID and the numeric functions rely on math/rand
// and are not cryptographically secure; use CryptoID when security properties
// matter.
//
// All functions in this package are safe for concurrent use.
package randn
