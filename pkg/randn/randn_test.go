package randn

import (
	"fmt"
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
