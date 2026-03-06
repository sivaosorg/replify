package slogger

// Formatter serialises a log Entry to a byte slice.
// Implementations must be safe for concurrent use from multiple goroutines.
type Formatter interface {
	Format(*Entry) ([]byte, error)
}
