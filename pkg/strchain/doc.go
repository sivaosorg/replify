// Package strchain provides a fluent, chainable string builder API built on top of
// Go's strings.Builder. It offers two implementations with identical APIs:
//
//   - [StringWeaver]: A high-performance, non-thread-safe builder for single-goroutine usage.
//     It wraps strings.Builder directly with zero overhead beyond the method call indirection.
//
//   - [SafeStringWeaver]: A thread-safe builder protected by sync.Mutex, suitable for
//     concurrent access from multiple goroutines. Each operation acquires a lock, adding
//     approximately 10–20ns of overhead per call.
//
// Both types implement the [Weaver] interface, which defines the full set of chainable
// methods. This allows writing functions that accept either implementation polymorphically.
//
// # Choosing Between StringWeaver and SafeStringWeaver
//
// Use [StringWeaver] (via [New], [NewWithCapacity], or [From]) when:
//   - The builder is confined to a single goroutine.
//   - Maximum throughput is required (zero synchronization overhead).
//
// Use [SafeStringWeaver] (via [NewSafe], [NewSafeWithCapacity], or [SafeFrom]) when:
//   - Multiple goroutines write to the same builder concurrently.
//   - Thread safety is required without external locking.
//
// # Fluent API
//
// All chainable methods return the [Weaver] interface, enabling expressive method chains:
//
//	result := strchain.New().
//	    Append("SELECT ").
//	    Join(", ", "id", "name", "email").
//	    Append(" FROM users").
//	    Build()
//
// # Callback Methods
//
// The [StringWeaver.When], [StringWeaver.Unless], and [StringWeaver.Each] methods accept
// callbacks with concrete receiver types. These methods are not part of the [Weaver] interface
// because their function signatures differ between implementations. They are available
// directly on each concrete type.
//
// # Polymorphism
//
// Functions that build strings can accept [Weaver] to work with either implementation:
//
//	func buildGreeting(w strchain.Weaver, name string) string {
//	    return w.Append("Hello, ").Append(name).Append("!").Build()
//	}
//
//	// Single-threaded
//	buildGreeting(strchain.New(), "Alice")
//
//	// Thread-safe
//	buildGreeting(strchain.NewSafe(), "Bob")
package strchain
