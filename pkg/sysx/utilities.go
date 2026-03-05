package sysx

import (
	"fmt"
	"strings"
)

// ///////////////////////////
// Section: Internal string helpers
// ///////////////////////////

// isZero reports whether s is empty or consists entirely of whitespace.
func isZero(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// trimSpace returns s with all leading and trailing whitespace removed.
func trimSpace(s string) string {
	return strings.TrimSpace(s)
}

// parseBoolString parses a lowercase, trimmed string as a boolean.
//
// Recognised true  values: "1", "true", "yes", "on"
// Recognised false values: "0", "false", "no", "off"
//
// Returns (value, true) when the string is recognised, or (false, false) when
// it is not.
func parseBoolString(s string) (bool, bool) {
	switch s {
	case "1", "true", "yes", "on":
		return true, true
	case "0", "false", "no", "off":
		return false, true
	}
	return false, false
}

// splitLines splits s into individual lines by "\n", stripping any trailing
// "\r" (to handle "\r\n" line endings). A single trailing newline is consumed
// so that "a\nb\n" returns ["a","b"] rather than ["a","b",""]. An empty input
// returns nil.
func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	s = strings.TrimRight(s, "\n")
	parts := strings.Split(s, "\n")
	result := make([]string, len(parts))
	for i, p := range parts {
		result[i] = strings.TrimRight(p, "\r")
	}
	return result
}

// ///////////////////////////
// Section: Internal I/O helpers
// ///////////////////////////

// commandBuffer is a minimal zero-allocation-friendly byte accumulator used
// internally by Execute to capture stdout and stderr. It satisfies io.Writer.
type commandBuffer struct {
	buf strings.Builder
}

// Write appends p to the buffer.
func (b *commandBuffer) Write(p []byte) (int, error) {
	return b.buf.Write(p)
}

// String returns the accumulated content as a string.
func (b *commandBuffer) String() string {
	return b.buf.String()
}

// ///////////////////////////
// Section: Exported composite helpers
// ///////////////////////////

// UserInfo returns a formatted string containing the numeric user and group
// identifiers of the current process.
//
// The string has the form "uid=X gid=Y", where X and Y are the values
// returned by UID() and GID() respectively.
//
// Returns:
//
//	A non-empty string of the form "uid=<uid> gid=<gid>".
//
// Example:
//
//	fmt.Println(sysx.UserInfo()) // "uid=1000 gid=1000"
func UserInfo() string {
	return fmt.Sprintf("uid=%d gid=%d", UID(), GID())
}
