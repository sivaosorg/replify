package randn

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"sync/atomic"
	"time"
)

func init() {
	for i := range len(decodeBytes) {
		decodeBytes[i] = 0xFF
	}
	for i := range len(encoding) {
		decodeBytes[encoding[i]] = byte(i)
	}

	// XOR the PID with a platform-specific offset to help ensure uniqueness in
	// containerized environments (e.g., Linux cgroups). On non-Linux platforms
	// pidContainerOffset returns 0 and this is a no-op.
	pid ^= pidContainerOffset()
}

// NewXID generates a new unique XID using the current time.
func NewXID() XID {
	return NewXIDWithTime(time.Now())
}

// NewXIDWithTime generates a new unique XID using the provided time.
func NewXIDWithTime(t time.Time) XID {
	var id XID
	binary.BigEndian.PutUint32(id[:], uint32(t.Unix()))
	id[4] = machineID[0]
	id[5] = machineID[1]
	id[6] = machineID[2]
	id[7] = byte(pid >> 8)
	id[8] = byte(pid)
	i := atomic.AddUint32(&idCounter, 1)
	id[9] = byte(i >> 16)
	id[10] = byte(i >> 8)
	id[11] = byte(i)
	return id
}

// IsZero returns true if the XID is the zero value.
func (id XID) IsZero() bool {
	return id == zeroXID
}

// String returns the base32 hex lowercased string representation of the ID.
// The output is 20 characters long.
func (id XID) String() string {
	text := make([]byte, encodedLen)
	encode(text, id[:])
	return string(text)
}

// Time extracts the time component from the ID.
func (id XID) Time() time.Time {
	secs := int64(binary.BigEndian.Uint32(id[0:4]))
	return time.Unix(secs, 0)
}

// Value implements the database/sql/driver.Valuer interface.
func (id XID) Value() (driver.Value, error) {
	if id.IsZero() {
		return nil, nil
	}
	return id.String(), nil
}

// Compare compares two XIDs. It returns -1 if id < other, 1 if id > other, and 0 if they are equal.
func (id XID) Compare(other XID) int {
	return bytes.Compare(id[:], other[:])
}

// Scan implements the database/sql.Scanner interface.
func (id *XID) Scan(value any) error {
	switch val := value.(type) {
	case string:
		return id.Unmarshal([]byte(val))
	case []byte:
		return id.Unmarshal(val)
	case nil:
		*id = zeroXID
		return nil
	default:
		return fmt.Errorf("randn: scanning unsupported type: %T", value)
	}
}

// Unmarshal decodes a 20-byte base32 hex lowercased representation into the XID.
func (id *XID) Unmarshal(text []byte) error {
	if len(text) != encodedLen {
		return ErrInvalidID
	}
	for _, c := range text {
		if decodeBytes[c] == 0xFF {
			return ErrInvalidID
		}
	}
	if !decode(id, text) {
		*id = zeroXID
		return ErrInvalidID
	}
	return nil
}
