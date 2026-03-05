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
