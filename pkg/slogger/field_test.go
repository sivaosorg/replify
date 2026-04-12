package slogger

import (
	"errors"
	"math"
	"strings"
	"testing"
	"time"
)

// =============================================================================
// Field Accessor Tests
// =============================================================================

func TestField_Key(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		key  string
	}{
		{name: "simple key", key: "user"},
		{name: "empty key", key: ""},
		{name: "key with spaces", key: "user name"},
		{name: "unicode key", key: "日本語"},
		{name: "special chars", key: "key!@#$%"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := String(tt.key, "value")
			assertEqual(t, tt.key, f.Key())
		})
	}
}

func TestField_Type(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		field Field
		want  FieldType
	}{
		{name: "string type", field: String("k", "v"), want: StringType},
		{name: "int64 type", field: Int64("k", 1), want: Int64Type},
		{name: "int type", field: Int("k", 1), want: Int64Type},
		{name: "int8 type", field: Int8("k", 1), want: Int8Type},
		{name: "int16 type", field: Int16("k", 1), want: Int16Type},
		{name: "int32 type", field: Int32("k", 1), want: Int32Type},
		{name: "uint type", field: Uint("k", 1), want: UintType},
		{name: "uint8 type", field: Uint8("k", 1), want: Uint8Type},
		{name: "uint16 type", field: Uint16("k", 1), want: Uint16Type},
		{name: "uint32 type", field: Uint32("k", 1), want: Uint32Type},
		{name: "uint64 type", field: Uint64("k", 1), want: Uint64Type},
		{name: "float32 type", field: Float32("k", 1.0), want: Float32Type},
		{name: "float64 type", field: Float64("k", 1.0), want: Float64Type},
		{name: "bool type", field: Bool("k", true), want: BoolType},
		{name: "error type", field: Err(errors.New("err")), want: ErrorType},
		{name: "time type", field: Time("k", time.Now()), want: TimeType},
		{name: "timef type", field: Timef("k", time.Now(), "2006-01-02"), want: TimefType},
		{name: "duration type", field: Duration("k", time.Second), want: DurationType},
		{name: "any type", field: Any("k", struct{}{}), want: AnyType},
		{name: "json type", field: JSON("k", map[string]int{"a": 1}), want: JSONType},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assertEqual(t, tt.want, tt.field.Type())
			// Also test deprecated FieldType() method
			assertEqual(t, tt.want, tt.field.FieldType())
		})
	}
}

func TestField_StringValue(t *testing.T) {
	t.Parallel()
	f := String("key", "hello world")
	assertEqual(t, "hello world", f.StringValue())
}

func TestField_IntValue(t *testing.T) {
	t.Parallel()
	f := Int64("key", 12345)
	assertEqual(t, int64(12345), f.IntValue())
}

func TestField_Uint64Value(t *testing.T) {
	t.Parallel()
	f := Uint64("key", 12345)
	assertEqual(t, uint64(12345), f.Uint64Value())
}

func TestField_FloatValue(t *testing.T) {
	t.Parallel()
	f := Float64("key", 3.14159)
	assertEqual(t, 3.14159, f.FloatValue())
}

func TestField_BoolValue(t *testing.T) {
	t.Parallel()

	t.Run("true", func(t *testing.T) {
		t.Parallel()
		f := Bool("key", true)
		assertTrue(t, f.BoolValue())
	})

	t.Run("false", func(t *testing.T) {
		t.Parallel()
		f := Bool("key", false)
		assertFalse(t, f.BoolValue())
	})
}

func TestField_ErrValue(t *testing.T) {
	t.Parallel()

	t.Run("non-nil error", func(t *testing.T) {
		t.Parallel()
		err := errors.New("test error")
		f := Err(err)
		assertEqual(t, err, f.ErrValue())
	})

	t.Run("nil error", func(t *testing.T) {
		t.Parallel()
		f := Err(nil)
		assertNil(t, f.ErrValue())
	})
}

func TestField_TimeValue(t *testing.T) {
	t.Parallel()
	now := time.Now()
	f := Time("key", now)
	assertEqual(t, now, f.TimeValue())
}

func TestField_DurationValue(t *testing.T) {
	t.Parallel()
	dur := 5 * time.Second
	f := Duration("key", dur)
	assertEqual(t, dur, f.DurationValue())
}

func TestField_AnyValue(t *testing.T) {
	t.Parallel()
	val := map[string]int{"a": 1, "b": 2}
	f := Any("key", val)
	assertEqual(t, val, f.AnyValue())
}

// =============================================================================
// Field Constructor Tests - Primitive Types
// =============================================================================

func TestString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		key      string
		val      string
		wantVal  string
	}{
		{name: "normal string", key: "msg", val: "hello", wantVal: "hello"},
		{name: "empty string", key: "empty", val: "", wantVal: ""},
		{name: "unicode string", key: "unicode", val: "日本語🎉", wantVal: "日本語🎉"},
		{name: "special chars", key: "special", val: "a\nb\tc", wantVal: "a\nb\tc"},
		{name: "very long string", key: "long", val: strings.Repeat("a", 10000), wantVal: strings.Repeat("a", 10000)},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := String(tt.key, tt.val)
			assertEqual(t, tt.key, f.Key())
			assertEqual(t, StringType, f.Type())
			assertEqual(t, tt.wantVal, f.Value())
		})
	}
}

func TestBool(t *testing.T) {
	t.Parallel()

	t.Run("true value", func(t *testing.T) {
		t.Parallel()
		f := Bool("enabled", true)
		assertEqual(t, "enabled", f.Key())
		assertEqual(t, BoolType, f.Type())
		assertEqual(t, "true", f.Value())
	})

	t.Run("false value", func(t *testing.T) {
		t.Parallel()
		f := Bool("enabled", false)
		assertEqual(t, "false", f.Value())
	})
}

// =============================================================================
// Field Constructor Tests - Signed Integer Types
// =============================================================================

func TestInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		val     int
		wantVal string
	}{
		{name: "positive", val: 42, wantVal: "42"},
		{name: "zero", val: 0, wantVal: "0"},
		{name: "negative", val: -100, wantVal: "-100"},
		{name: "max int", val: math.MaxInt, wantVal: itoa64(int64(math.MaxInt))},
		{name: "min int", val: math.MinInt, wantVal: itoa64(int64(math.MinInt))},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := Int("n", tt.val)
			assertEqual(t, Int64Type, f.Type())
			assertEqual(t, tt.wantVal, f.Value())
		})
	}
}

func TestInt8(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		val     int8
		wantVal string
	}{
		{name: "positive", val: 42, wantVal: "42"},
		{name: "zero", val: 0, wantVal: "0"},
		{name: "negative", val: -100, wantVal: "-100"},
		{name: "max", val: math.MaxInt8, wantVal: "127"},
		{name: "min", val: math.MinInt8, wantVal: "-128"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := Int8("n", tt.val)
			assertEqual(t, Int8Type, f.Type())
			assertEqual(t, tt.wantVal, f.Value())
		})
	}
}

func TestInt16(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		val     int16
		wantVal string
	}{
		{name: "positive", val: 1000, wantVal: "1000"},
		{name: "max", val: math.MaxInt16, wantVal: "32767"},
		{name: "min", val: math.MinInt16, wantVal: "-32768"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := Int16("n", tt.val)
			assertEqual(t, Int16Type, f.Type())
			assertEqual(t, tt.wantVal, f.Value())
		})
	}
}

func TestInt32(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		val     int32
		wantVal string
	}{
		{name: "positive", val: 100000, wantVal: "100000"},
		{name: "max", val: math.MaxInt32, wantVal: "2147483647"},
		{name: "min", val: math.MinInt32, wantVal: "-2147483648"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := Int32("n", tt.val)
			assertEqual(t, Int32Type, f.Type())
			assertEqual(t, tt.wantVal, f.Value())
		})
	}
}

func TestInt64(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		val     int64
		wantVal string
	}{
		{name: "positive", val: 9876543210, wantVal: "9876543210"},
		{name: "zero", val: 0, wantVal: "0"},
		{name: "negative", val: -9876543210, wantVal: "-9876543210"},
		{name: "max", val: math.MaxInt64, wantVal: "9223372036854775807"},
		{name: "min", val: math.MinInt64, wantVal: "-9223372036854775808"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := Int64("n", tt.val)
			assertEqual(t, Int64Type, f.Type())
			assertEqual(t, tt.wantVal, f.Value())
		})
	}
}

// =============================================================================
// Field Constructor Tests - Unsigned Integer Types
// =============================================================================

func TestUint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		val     uint
		wantVal string
	}{
		{name: "positive", val: 42, wantVal: "42"},
		{name: "zero", val: 0, wantVal: "0"},
		{name: "max", val: math.MaxUint, wantVal: utoa64(uint64(math.MaxUint))},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := Uint("n", tt.val)
			assertEqual(t, UintType, f.Type())
			assertEqual(t, tt.wantVal, f.Value())
		})
	}
}

func TestUint8(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		val     uint8
		wantVal string
	}{
		{name: "positive", val: 42, wantVal: "42"},
		{name: "zero", val: 0, wantVal: "0"},
		{name: "max", val: math.MaxUint8, wantVal: "255"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := Uint8("n", tt.val)
			assertEqual(t, Uint8Type, f.Type())
			assertEqual(t, tt.wantVal, f.Value())
		})
	}
}

func TestUint16(t *testing.T) {
	t.Parallel()

	f := Uint16("n", math.MaxUint16)
	assertEqual(t, Uint16Type, f.Type())
	assertEqual(t, "65535", f.Value())
}

func TestUint32(t *testing.T) {
	t.Parallel()

	f := Uint32("n", math.MaxUint32)
	assertEqual(t, Uint32Type, f.Type())
	assertEqual(t, "4294967295", f.Value())
}

func TestUint64(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		val     uint64
		wantVal string
	}{
		{name: "positive", val: 9876543210, wantVal: "9876543210"},
		{name: "zero", val: 0, wantVal: "0"},
		{name: "max", val: math.MaxUint64, wantVal: "18446744073709551615"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := Uint64("n", tt.val)
			assertEqual(t, Uint64Type, f.Type())
			assertEqual(t, tt.wantVal, f.Value())
		})
	}
}

// =============================================================================
// Field Constructor Tests - Floating-Point Types
// =============================================================================

func TestFloat32(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		val  float32
	}{
		{name: "positive", val: 3.14},
		{name: "zero", val: 0},
		{name: "negative", val: -3.14},
		{name: "very small", val: 0.000001},
		{name: "very large", val: 1e10},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := Float32("n", tt.val)
			assertEqual(t, Float32Type, f.Type())
			assertNotEmpty(t, f.Value())
		})
	}
}

func TestFloat64(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		val     float64
		wantVal string
	}{
		{name: "pi", val: 3.14, wantVal: "3.14"},
		{name: "zero", val: 0, wantVal: "0"},
		{name: "negative", val: -1.5, wantVal: "-1.5"},
		{name: "integer value", val: 42, wantVal: "42"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := Float64("n", tt.val)
			assertEqual(t, Float64Type, f.Type())
			assertEqual(t, tt.wantVal, f.Value())
		})
	}
}

// =============================================================================
// Field Constructor Tests - Time and Duration
// =============================================================================

func TestTime(t *testing.T) {
	t.Parallel()

	t.Run("normal time", func(t *testing.T) {
		t.Parallel()
		ts := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
		f := Time("at", ts)
		assertEqual(t, TimeType, f.Type())
		assertContains(t, f.Value(), "2024-01-15")
	})

	t.Run("zero time", func(t *testing.T) {
		t.Parallel()
		var ts time.Time
		f := Time("at", ts)
		assertEqual(t, "", f.Value())
	})
}

func TestTimef(t *testing.T) {
	t.Parallel()

	t.Run("custom format", func(t *testing.T) {
		t.Parallel()
		ts := time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC)
		f := Timef("at", ts, "2006/01/02")
		assertEqual(t, TimefType, f.Type())
		assertEqual(t, "2024/06/15", f.Value())
	})

	t.Run("empty format", func(t *testing.T) {
		t.Parallel()
		ts := time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC)
		f := Timef("at", ts, "")
		assertEqual(t, TimefType, f.Type())
		// Falls back to default format
		assertContains(t, f.Value(), "2024-06-15")
	})

	t.Run("zero time", func(t *testing.T) {
		t.Parallel()
		var ts time.Time
		f := Timef("at", ts, "2006-01-02")
		assertEqual(t, "", f.Value())
	})
}

func TestDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		val     time.Duration
		wantVal string
	}{
		{name: "milliseconds", val: 500 * time.Millisecond, wantVal: "500ms"},
		{name: "seconds", val: 5 * time.Second, wantVal: "5s"},
		{name: "minutes", val: 2 * time.Minute, wantVal: "2m0s"},
		{name: "hours", val: 3 * time.Hour, wantVal: "3h0m0s"},
		{name: "zero", val: 0, wantVal: "0s"},
		{name: "negative", val: -5 * time.Second, wantVal: "-5s"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := Duration("took", tt.val)
			assertEqual(t, DurationType, f.Type())
			assertEqual(t, tt.wantVal, f.Value())
		})
	}
}

// =============================================================================
// Field Constructor Tests - Error and Generic Types
// =============================================================================

func TestErr(t *testing.T) {
	t.Parallel()

	t.Run("non-nil error", func(t *testing.T) {
		t.Parallel()
		err := errors.New("test error")
		f := Err(err)
		assertEqual(t, "error", f.Key())
		assertEqual(t, ErrorType, f.Type())
		assertEqual(t, "test error", f.Value())
	})

	t.Run("nil error", func(t *testing.T) {
		t.Parallel()
		f := Err(nil)
		assertEqual(t, "error", f.Key())
		assertEqual(t, ErrorType, f.Type())
		assertEqual(t, "<nil>", f.Value())
	})
}

func TestAny(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		val     interface{}
		wantVal string
	}{
		{name: "slice", val: []int{1, 2, 3}, wantVal: "[1 2 3]"},
		{name: "map", val: map[string]int{"a": 1}, wantVal: "map[a:1]"},
		{name: "struct", val: struct{ Name string }{Name: "test"}, wantVal: "{test}"},
		{name: "nil", val: nil, wantVal: "<nil>"},
		{name: "string", val: "hello", wantVal: "hello"},
		{name: "int", val: 42, wantVal: "42"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := Any("data", tt.val)
			assertEqual(t, AnyType, f.Type())
			assertEqual(t, tt.wantVal, f.Value())
		})
	}
}

func TestJSON(t *testing.T) {
	t.Parallel()

	t.Run("valid json object", func(t *testing.T) {
		t.Parallel()
		data := map[string]int{"a": 1, "b": 2}
		f := JSON("data", data)
		assertEqual(t, JSONType, f.Type())
		// JSON encoding order may vary
		assertContains(t, f.Value(), `"a":1`)
		assertContains(t, f.Value(), `"b":2`)
	})

	t.Run("simple value", func(t *testing.T) {
		t.Parallel()
		f := JSON("num", 42)
		assertEqual(t, JSONType, f.Type())
		assertEqual(t, "42", f.Value())
	})

	t.Run("string value", func(t *testing.T) {
		t.Parallel()
		f := JSON("str", "hello")
		assertEqual(t, JSONType, f.Type())
	})
}

// =============================================================================
// Field.Value Tests (comprehensive)
// =============================================================================

func TestField_Value_AllTypes(t *testing.T) {
	t.Parallel()

	// Test default case (unknown type)
	t.Run("unknown type", func(t *testing.T) {
		t.Parallel()
		f := Field{key: "test", typ: FieldType(999)}
		assertEqual(t, "", f.Value())
	})
}

// =============================================================================
// Field Edge Cases
// =============================================================================

func TestField_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("empty key", func(t *testing.T) {
		t.Parallel()
		f := String("", "value")
		assertEqual(t, "", f.Key())
		assertEqual(t, "value", f.Value())
	})

	t.Run("very long value", func(t *testing.T) {
		t.Parallel()
		longStr := strings.Repeat("x", 100000)
		f := String("key", longStr)
		assertEqual(t, longStr, f.Value())
	})

	t.Run("unicode in value", func(t *testing.T) {
		t.Parallel()
		f := String("key", "こんにちは世界🌍")
		assertEqual(t, "こんにちは世界🌍", f.Value())
	})

	t.Run("special characters in value", func(t *testing.T) {
		t.Parallel()
		f := String("key", "line1\nline2\ttab")
		assertEqual(t, "line1\nline2\ttab", f.Value())
	})

	t.Run("null byte in value", func(t *testing.T) {
		t.Parallel()
		f := String("key", "before\x00after")
		assertEqual(t, "before\x00after", f.Value())
	})
}

// =============================================================================
// Field Benchmarks
// =============================================================================

func BenchmarkString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = String("key", "value")
	}
}

func BenchmarkInt64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Int64("key", 12345678)
	}
}

func BenchmarkFloat64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Float64("key", 3.14159)
	}
}

func BenchmarkField_Value(b *testing.B) {
	f := String("key", "test value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = f.Value()
	}
}

func BenchmarkJSON(b *testing.B) {
	data := map[string]interface{}{"key": "value", "num": 123}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = JSON("data", data)
	}
}
