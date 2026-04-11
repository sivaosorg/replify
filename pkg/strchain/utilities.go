package strchain

import (
	"encoding/json"
	"strconv"
)

// jsonEncodeString returns the RFC 8259-compliant JSON encoding of s,
// including the surrounding double-quote characters.
// In practice json.Marshal never errors for a plain string value in any known
// Go version; the strconv.Quote fallback is an unreachable defensive measure.
// Note: the fallback produces Go-style escaping (e.g., \x00 for null) which is
// NOT valid JSON — if reached it should be treated as a bug.
func jsonEncodeString(s string) string {
	if b, err := json.Marshal(s); err == nil {
		return string(b)
	}
	return strconv.Quote(s)
}
