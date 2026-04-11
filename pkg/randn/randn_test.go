package randn

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"testing"
)

func TestNewXID(t *testing.T) {
	id := NewXID()
	if id.IsZero() {
		t.Errorf("NewXID returned a nil ID")
	}

	str := id.String()
	if len(str) != 20 {
		t.Errorf("XID length must be 20, got %d", len(str))
	}

	id2 := NewXID()
	if id.Compare(id2) >= 0 {
		t.Errorf("XID should be sortable, %s >= %s", id.String(), id2.String())
	}
}

func TestUnmarshalXID(t *testing.T) {
	id := NewXID()
	str := id.String()

	var id2 XID
	err := id2.Unmarshal([]byte(str))
	if err != nil {
		t.Fatalf("UnmarshalText failed: %v", err)
	}

	if id.Compare(id2) != 0 {
		t.Errorf("Unmarshaled ID does not match original: %s != %s", id2.String(), id.String())
	}
}

func ExampleNewXID() {
	id := NewXID()
	fmt.Printf("XID: %s\n", id.String())
	fmt.Printf("Time: %s\n", id.Time().UTC())
}

// --- UUID tests ---

var uuidWithDashRe = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func TestUUID(t *testing.T) {
	uuid, err := UUID()
	if err != nil {
		t.Fatalf("UUID() returned error: %v", err)
	}
	if !uuidWithDashRe.MatchString(uuid) {
		t.Errorf("UUID() %q does not match RFC 4122 v4 pattern", uuid)
	}
}

func TestUUID_Uniqueness(t *testing.T) {
	seen := make(map[string]struct{}, 100)
	for i := range 100 {
		u, err := UUID()
		if err != nil {
			t.Fatalf("UUID() call %d returned error: %v", i, err)
		}
		if _, dup := seen[u]; dup {
			t.Fatalf("UUID() produced a duplicate: %q", u)
		}
		seen[u] = struct{}{}
	}
}

func TestUUIDSep(t *testing.T) {
	tests := []struct {
		delimiter string
		wantParts int
		wantLen   int // expected total hex chars (32) plus delimiter lengths
	}{
		{"-", 5, 32 + 4}, // standard UUID: 8+4+4+4+12 = 32 hex + 4 dashes
		{"", 1, 32},
		{"/", 5, 32 + 4},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("delim=%q", tc.delimiter), func(t *testing.T) {
			uuid, err := UUIDSep(tc.delimiter)
			if err != nil {
				t.Fatalf("UUIDSep(%q) returned error: %v", tc.delimiter, err)
			}
			if tc.delimiter != "" {
				parts := strings.Split(uuid, tc.delimiter)
				if len(parts) != tc.wantParts {
					t.Errorf("expected %d parts, got %d: %q", tc.wantParts, len(parts), uuid)
				}
			}
			if len(uuid) != tc.wantLen {
				t.Errorf("expected length %d, got %d: %q", tc.wantLen, len(uuid), uuid)
			}
		})
	}
}

// TestUUIDSep_RFC4122Bits verifies that version and variant bits are set correctly.
func TestUUIDSep_RFC4122Bits(t *testing.T) {
	for range 20 {
		uuid, err := UUID()
		if err != nil {
			t.Fatalf("UUID() returned error: %v", err)
		}
		// UUID format: 8-4-4-4-12; version nibble is at position 14 (0-indexed, after two dashes)
		// "xxxxxxxx-xxxx-Vxxx-..." → index 14 = V must be '4'
		if uuid[14] != '4' {
			t.Errorf("UUID version nibble should be '4', got %q in %q", uuid[14], uuid)
		}
		// Variant nibble is the first character of the 4th group (index 19)
		variantChar := uuid[19]
		if variantChar != '8' && variantChar != '9' && variantChar != 'a' && variantChar != 'b' {
			t.Errorf("UUID variant nibble should be 8/9/a/b, got %q in %q", variantChar, uuid)
		}
	}
}

// --- ID generation tests ---

func TestRandID(t *testing.T) {
	tests := []struct {
		length int
	}{
		{0},
		{1},
		{8},
		{16},
		{64},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("len=%d", tc.length), func(t *testing.T) {
			id := RandID(tc.length)
			if len(id) != tc.length {
				t.Errorf("RandID(%d) returned string of length %d", tc.length, len(id))
			}
			const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
			for _, c := range id {
				if !strings.ContainsRune(charset, c) {
					t.Errorf("RandID(%d) contains character %q not in charset", tc.length, c)
				}
			}
		})
	}
}

func TestCryptoID(t *testing.T) {
	id := CryptoID()
	if len(id) != 32 {
		t.Errorf("CryptoID() should return 32-char hex string, got length %d: %q", len(id), id)
	}
	for _, c := range id {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Errorf("CryptoID() contains non-hex character %q in %q", c, id)
		}
	}
}

func TestCryptoID_Uniqueness(t *testing.T) {
	seen := make(map[string]struct{}, 50)
	for i := range 50 {
		id := CryptoID()
		if _, dup := seen[id]; dup {
			t.Fatalf("CryptoID() produced a duplicate at iteration %d: %q", i, id)
		}
		seen[id] = struct{}{}
	}
}

func TestTimeID(t *testing.T) {
	id1 := TimeID()
	id2 := TimeID()
	if id1 == "" || id2 == "" {
		t.Error("TimeID() returned an empty string")
	}
	// They should almost certainly differ; this is a best-effort check.
	if id1 == id2 {
		t.Logf("TimeID() produced identical values on two calls (unlikely but possible): %q", id1)
	}
}

func TestRandUUID(t *testing.T) {
	uuid := RandUUID()
	if uuid == "" {
		t.Fatal("RandUUID() returned empty string unexpectedly")
	}
	if !uuidWithDashRe.MatchString(uuid) {
		t.Errorf("RandUUID() %q does not match RFC 4122 v4 pattern", uuid)
	}
}

func TestRandIDHex(t *testing.T) {
	tests := []struct{ byteLen int }{{4}, {8}, {16}, {32}}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("byteLen=%d", tc.byteLen), func(t *testing.T) {
			id := RandIDHex(tc.byteLen)
			if len(id) != tc.byteLen*2 {
				t.Errorf("RandIDHex(%d) expected length %d, got %d: %q", tc.byteLen, tc.byteLen*2, len(id), id)
			}
		})
	}
}

// --- Numeric random tests ---

func TestRandInt(t *testing.T) {
	for range 20 {
		if v := RandInt(); v < 0 {
			t.Errorf("RandInt() returned negative value: %d", v)
		}
	}
}

func TestRandIntr(t *testing.T) {
	tests := []struct {
		min, max int
	}{
		{0, 0}, // degenerate: min == max → return min
		{5, 5}, // degenerate
		{1, 10},
		{-10, 10},
		{-5, -1},
		{0, 1},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("[%d,%d]", tc.min, tc.max), func(t *testing.T) {
			for range 100 {
				v := RandIntr(tc.min, tc.max)
				if v < tc.min || v > tc.max {
					t.Errorf("RandIntr(%d, %d) = %d: out of range", tc.min, tc.max, v)
				}
			}
		})
	}
}

func TestRandIntr_MinGreaterThanMax(t *testing.T) {
	v := RandIntr(10, 5)
	if v != 10 {
		t.Errorf("RandIntr(10, 5) should return min=10, got %d", v)
	}
}

func TestRandUint32(t *testing.T) {
	seen := make(map[uint32]struct{}, 20)
	for range 20 {
		v := RandUint32()
		// Values are 24-bit; all should be within [0, 0xFFFFFF]
		if v > 0xFFFFFF {
			t.Errorf("RandUint32() = %d exceeds 24-bit range", v)
		}
		seen[v] = struct{}{}
	}
	if len(seen) < 15 {
		t.Errorf("RandUint32() produced too few unique values in 20 calls: %d", len(seen))
	}
}

func TestRandFt64(t *testing.T) {
	for range 20 {
		v := RandFt64()
		if v < 0.0 || v >= 1.0 {
			t.Errorf("RandFt64() = %v: out of [0.0, 1.0)", v)
		}
	}
}

func TestRandFt64r(t *testing.T) {
	tests := []struct{ lo, hi float64 }{{0, 1}, {-1, 1}, {5.5, 10.5}}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("[%.1f,%.1f)", tc.lo, tc.hi), func(t *testing.T) {
			for range 50 {
				v := RandFt64r(tc.lo, tc.hi)
				if v < tc.lo || v >= tc.hi {
					t.Errorf("RandFt64r(%.1f, %.1f) = %v: out of range", tc.lo, tc.hi, v)
				}
			}
		})
	}
}

func TestRandFt32(t *testing.T) {
	for range 20 {
		v := RandFt32()
		if v < 0.0 || v >= 1.0 {
			t.Errorf("RandFt32() = %v: out of [0.0, 1.0)", v)
		}
	}
}

func TestRandFt32r(t *testing.T) {
	tests := []struct{ lo, hi float32 }{{0, 1}, {-1, 1}, {5.5, 10.5}}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("[%.1f,%.1f)", tc.lo, tc.hi), func(t *testing.T) {
			for range 50 {
				v := RandFt32r(tc.lo, tc.hi)
				if v < tc.lo || v >= tc.hi {
					t.Errorf("RandFt32r(%.1f, %.1f) = %v: out of range", tc.lo, tc.hi, v)
				}
			}
		})
	}
}

func TestRandByte(t *testing.T) {
	tests := []struct{ count int }{{0}, {1}, {16}, {256}}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("count=%d", tc.count), func(t *testing.T) {
			b := RandByte(tc.count)
			if len(b) != tc.count {
				t.Errorf("RandByte(%d) returned slice of length %d", tc.count, len(b))
			}
		})
	}
}

// --- Concurrency safety tests ---

// TestRandIntr_Concurrent verifies RandIntr is free of data races when called
// from many goroutines simultaneously. Run with -race to detect races.
func TestRandIntr_Concurrent(t *testing.T) {
	const goroutines = 50
	const callsEach = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for range goroutines {
		go func() {
			defer wg.Done()
			for range callsEach {
				v := RandIntr(1, 1000)
				if v < 1 || v > 1000 {
					t.Errorf("RandIntr concurrent: value %d out of [1,1000]", v)
				}
			}
		}()
	}
	wg.Wait()
}

// TestRandID_Concurrent verifies RandID is free of data races.
func TestRandID_Concurrent(t *testing.T) {
	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for range goroutines {
		go func() {
			defer wg.Done()
			id := RandID(16)
			if len(id) != 16 {
				t.Errorf("RandID concurrent: expected length 16, got %d", len(id))
			}
		}()
	}
	wg.Wait()
}

// TestUUID_Concurrent verifies UUID generation is race-free.
func TestUUID_Concurrent(t *testing.T) {
	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for range goroutines {
		go func() {
			defer wg.Done()
			u, err := UUID()
			if err != nil {
				t.Errorf("UUID concurrent: unexpected error: %v", err)
				return
			}
			if !uuidWithDashRe.MatchString(u) {
				t.Errorf("UUID concurrent: %q does not match RFC 4122 v4 pattern", u)
			}
		}()
	}
	wg.Wait()
}
