package strchain

import (
	"fmt"
	"strings"
	"sync"
	"testing"
)

// ---------------------------------------------------------------------------
// StringWeaver Tests (non-thread-safe)
// ---------------------------------------------------------------------------

func TestNew(t *testing.T) {
	sw := New()
	if sw == nil {
		t.Fatal("New() returned nil")
	}
	if sw.Len() != 0 {
		t.Errorf("New() should have length 0, got %d", sw.Len())
	}
}

func TestNewWithCapacity(t *testing.T) {
	sw := NewWithCapacity(1024)
	if sw == nil {
		t.Fatal("NewWithCapacity() returned nil")
	}
	if sw.Cap() < 1024 {
		t.Errorf("NewWithCapacity(1024) should have capacity >= 1024, got %d", sw.Cap())
	}
}

func TestFrom(t *testing.T) {
	sw := From("hello")
	if sw.String() != "hello" {
		t.Errorf("From(\"hello\") should produce \"hello\", got %q", sw.String())
	}
}

func TestStringWeaver_Append(t *testing.T) {
	result := New().Append("hello").Append(" ").Append("world").Build()
	if result != "hello world" {
		t.Errorf("expected \"hello world\", got %q", result)
	}
}

func TestStringWeaver_AppendF(t *testing.T) {
	result := New().AppendF("value: %d", 42).Build()
	if result != "value: 42" {
		t.Errorf("expected \"value: 42\", got %q", result)
	}
}

func TestStringWeaver_AppendByte(t *testing.T) {
	result := New().AppendByte('A').AppendByte('B').Build()
	if result != "AB" {
		t.Errorf("expected \"AB\", got %q", result)
	}
}

func TestStringWeaver_AppendRune(t *testing.T) {
	result := New().AppendRune('⌘').Build()
	if result != "⌘" {
		t.Errorf("expected \"⌘\", got %q", result)
	}
}

func TestStringWeaver_AppendBytes(t *testing.T) {
	result := New().AppendBytes([]byte("data")).Build()
	if result != "data" {
		t.Errorf("expected \"data\", got %q", result)
	}
}

func TestStringWeaver_AppendInt(t *testing.T) {
	result := New().AppendInt(123).Build()
	if result != "123" {
		t.Errorf("expected \"123\", got %q", result)
	}
}

func TestStringWeaver_AppendInt8(t *testing.T) {
	result := New().AppendInt8(8).Build()
	if result != "8" {
		t.Errorf("expected \"8\", got %q", result)
	}
}

func TestStringWeaver_AppendInt16(t *testing.T) {
	result := New().AppendInt16(16).Build()
	if result != "16" {
		t.Errorf("expected \"16\", got %q", result)
	}
}

func TestStringWeaver_AppendInt32(t *testing.T) {
	result := New().AppendInt32(32).Build()
	if result != "32" {
		t.Errorf("expected \"32\", got %q", result)
	}
}

func TestStringWeaver_AppendInt64(t *testing.T) {
	result := New().AppendInt64(64).Build()
	if result != "64" {
		t.Errorf("expected \"64\", got %q", result)
	}
}

func TestStringWeaver_AppendUint(t *testing.T) {
	result := New().AppendUint(123).Build()
	if result != "123" {
		t.Errorf("expected \"123\", got %q", result)
	}
}

func TestStringWeaver_AppendUint8(t *testing.T) {
	result := New().AppendUint8(8).Build()
	if result != "8" {
		t.Errorf("expected \"8\", got %q", result)
	}
}

func TestStringWeaver_AppendUint16(t *testing.T) {
	result := New().AppendUint16(16).Build()
	if result != "16" {
		t.Errorf("expected \"16\", got %q", result)
	}
}

func TestStringWeaver_AppendUint32(t *testing.T) {
	result := New().AppendUint32(32).Build()
	if result != "32" {
		t.Errorf("expected \"32\", got %q", result)
	}
}

func TestStringWeaver_AppendUint64(t *testing.T) {
	result := New().AppendUint64(64).Build()
	if result != "64" {
		t.Errorf("expected \"64\", got %q", result)
	}
}

func TestStringWeaver_AppendUintptr(t *testing.T) {
	result := New().AppendUintptr(100).Build()
	if result != "100" {
		t.Errorf("expected \"100\", got %q", result)
	}
}

func TestStringWeaver_AppendFloat32(t *testing.T) {
	result := New().AppendFloat32(3.14).Build()
	if !strings.HasPrefix(result, "3.14") {
		t.Errorf("expected result starting with \"3.14\", got %q", result)
	}
}

func TestStringWeaver_AppendFloat64(t *testing.T) {
	result := New().AppendFloat64(3.14159).Build()
	if !strings.HasPrefix(result, "3.14") {
		t.Errorf("expected result starting with \"3.14\", got %q", result)
	}
}

func TestStringWeaver_AppendBool(t *testing.T) {
	resultTrue := New().AppendBool(true).Build()
	resultFalse := New().AppendBool(false).Build()
	if resultTrue != "true" {
		t.Errorf("expected \"true\", got %q", resultTrue)
	}
	if resultFalse != "false" {
		t.Errorf("expected \"false\", got %q", resultFalse)
	}
}

func TestStringWeaver_Space(t *testing.T) {
	result := New().Append("a").Space().Append("b").Build()
	if result != "a b" {
		t.Errorf("expected \"a b\", got %q", result)
	}
}

func TestStringWeaver_Spaces(t *testing.T) {
	result := New().Append("a").Spaces(3).Append("b").Build()
	if result != "a   b" {
		t.Errorf("expected \"a   b\", got %q", result)
	}
}

func TestStringWeaver_Tab(t *testing.T) {
	result := New().Append("a").Tab().Append("b").Build()
	if result != "a\tb" {
		t.Errorf("expected \"a\\tb\", got %q", result)
	}
}

func TestStringWeaver_Tabs(t *testing.T) {
	result := New().Append("a").Tabs(2).Append("b").Build()
	if result != "a\t\tb" {
		t.Errorf("expected \"a\\t\\tb\", got %q", result)
	}
}

func TestStringWeaver_NewLine(t *testing.T) {
	result := New().Append("a").NewLine().Append("b").Build()
	if result != "a\nb" {
		t.Errorf("expected \"a\\nb\", got %q", result)
	}
}

func TestStringWeaver_NewLines(t *testing.T) {
	result := New().Append("a").NewLines(2).Append("b").Build()
	if result != "a\n\nb" {
		t.Errorf("expected \"a\\n\\nb\", got %q", result)
	}
}

func TestStringWeaver_Line(t *testing.T) {
	result := New().Line("hello").Line("world").Build()
	if result != "hello\nworld\n" {
		t.Errorf("expected \"hello\\nworld\\n\", got %q", result)
	}
}

func TestStringWeaver_LineF(t *testing.T) {
	result := New().LineF("id: %d", 42).Build()
	if result != "id: 42\n" {
		t.Errorf("expected \"id: 42\\n\", got %q", result)
	}
}

func TestStringWeaver_Repeat(t *testing.T) {
	result := New().Repeat("-", 5).Build()
	if result != "-----" {
		t.Errorf("expected \"-----\", got %q", result)
	}
}

func TestStringWeaver_Join(t *testing.T) {
	result := New().Join(", ", "a", "b", "c").Build()
	if result != "a, b, c" {
		t.Errorf("expected \"a, b, c\", got %q", result)
	}
}

func TestStringWeaver_AppendIf(t *testing.T) {
	resultTrue := New().Append("a").AppendIf(true, "b").Build()
	resultFalse := New().Append("a").AppendIf(false, "b").Build()
	if resultTrue != "ab" {
		t.Errorf("expected \"ab\", got %q", resultTrue)
	}
	if resultFalse != "a" {
		t.Errorf("expected \"a\", got %q", resultFalse)
	}
}

func TestStringWeaver_AppendIfF(t *testing.T) {
	result := New().AppendIfF(true, "n=%d", 42).Build()
	if result != "n=42" {
		t.Errorf("expected \"n=42\", got %q", result)
	}
	result = New().AppendIfF(false, "n=%d", 42).Build()
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestStringWeaver_When(t *testing.T) {
	sw := New()
	sw.Append("a")
	result := sw.When(true, func(sw *StringWeaver) {
		sw.Append("b")
	}).Build()
	if result != "ab" {
		t.Errorf("expected \"ab\", got %q", result)
	}
	sw2 := New()
	sw2.Append("a")
	result = sw2.When(false, func(sw *StringWeaver) {
		sw.Append("b")
	}).Build()
	if result != "a" {
		t.Errorf("expected \"a\", got %q", result)
	}
}

func TestStringWeaver_Unless(t *testing.T) {
	sw := New()
	sw.Append("a")
	result := sw.Unless(false, func(sw *StringWeaver) {
		sw.Append("b")
	}).Build()
	if result != "ab" {
		t.Errorf("expected \"ab\", got %q", result)
	}
	sw2 := New()
	sw2.Append("a")
	result = sw2.Unless(true, func(sw *StringWeaver) {
		sw.Append("b")
	}).Build()
	if result != "a" {
		t.Errorf("expected \"a\", got %q", result)
	}
}

func TestStringWeaver_Each(t *testing.T) {
	result := New().Each([]string{"a", "b", "c"}, func(sw *StringWeaver, item string) {
		sw.Append(item)
	}).Build()
	if result != "abc" {
		t.Errorf("expected \"abc\", got %q", result)
	}
}

func TestStringWeaver_Indent(t *testing.T) {
	result := New().Indent(2, "text").Build()
	if result != "    text" {
		t.Errorf("expected \"    text\", got %q", result)
	}
}

func TestStringWeaver_IndentLine(t *testing.T) {
	result := New().IndentLine(1, "item").Build()
	if result != "  item\n" {
		t.Errorf("expected \"  item\\n\", got %q", result)
	}
}

func TestStringWeaver_Wrap(t *testing.T) {
	result := New().Wrap("<", "tag", ">").Build()
	if result != "<tag>" {
		t.Errorf("expected \"<tag>\", got %q", result)
	}
}

func TestStringWeaver_Quote(t *testing.T) {
	result := New().Quote("text").Build()
	if result != `"text"` {
		t.Errorf("expected %q, got %q", `"text"`, result)
	}
}

func TestStringWeaver_SingleQuote(t *testing.T) {
	result := New().SingleQuote("text").Build()
	if result != "'text'" {
		t.Errorf("expected \"'text'\", got %q", result)
	}
}

func TestStringWeaver_Parenthesize(t *testing.T) {
	result := New().Parenthesize("expr").Build()
	if result != "(expr)" {
		t.Errorf("expected \"(expr)\", got %q", result)
	}
}

func TestStringWeaver_Bracket(t *testing.T) {
	result := New().Bracket("0").Build()
	if result != "[0]" {
		t.Errorf("expected \"[0]\", got %q", result)
	}
}

func TestStringWeaver_Brace(t *testing.T) {
	result := New().Brace("body").Build()
	if result != "{body}" {
		t.Errorf("expected \"{body}\", got %q", result)
	}
}

func TestStringWeaver_Comma(t *testing.T) {
	result := New().Append("a").Comma().Append("b").Build()
	if result != "a,b" {
		t.Errorf("expected \"a,b\", got %q", result)
	}
}

func TestStringWeaver_Dot(t *testing.T) {
	result := New().Append("end").Dot().Build()
	if result != "end." {
		t.Errorf("expected \"end.\", got %q", result)
	}
}

func TestStringWeaver_Colon(t *testing.T) {
	result := New().Append("key").Colon().Append("val").Build()
	if result != "key:val" {
		t.Errorf("expected \"key:val\", got %q", result)
	}
}

func TestStringWeaver_Semicolon(t *testing.T) {
	result := New().Append("x=1").Semicolon().Build()
	if result != "x=1;" {
		t.Errorf("expected \"x=1;\", got %q", result)
	}
}

func TestStringWeaver_Equals(t *testing.T) {
	result := New().Append("key").Equals().Append("val").Build()
	if result != "key=val" {
		t.Errorf("expected \"key=val\", got %q", result)
	}
}

func TestStringWeaver_Arrow(t *testing.T) {
	result := New().Append("a").Arrow().Append("b").Build()
	if result != "a->b" {
		t.Errorf("expected \"a->b\", got %q", result)
	}
}

func TestStringWeaver_FatArrow(t *testing.T) {
	result := New().Append("a").FatArrow().Append("b").Build()
	if result != "a=>b" {
		t.Errorf("expected \"a=>b\", got %q", result)
	}
}

func TestStringWeaver_Grow(t *testing.T) {
	sw := New()
	sw.Grow(512)
	if sw.Cap() < 512 {
		t.Errorf("after Grow(512), capacity should be >= 512, got %d", sw.Cap())
	}
}

func TestStringWeaver_Reset(t *testing.T) {
	sw := New()
	sw.Append("hello")
	sw.Reset()
	if sw.Len() != 0 {
		t.Errorf("after Reset(), length should be 0, got %d", sw.Len())
	}
	if sw.String() != "" {
		t.Errorf("after Reset(), string should be empty, got %q", sw.String())
	}
}

func TestStringWeaver_Len(t *testing.T) {
	sw := New()
	sw.Append("hello")
	if sw.Len() != 5 {
		t.Errorf("expected length 5, got %d", sw.Len())
	}
}

func TestStringWeaver_Cap(t *testing.T) {
	sw := NewWithCapacity(256)
	if sw.Cap() < 256 {
		t.Errorf("expected capacity >= 256, got %d", sw.Cap())
	}
}

func TestStringWeaver_String(t *testing.T) {
	result := New().Append("test").String()
	if result != "test" {
		t.Errorf("expected \"test\", got %q", result)
	}
}

func TestStringWeaver_Build(t *testing.T) {
	result := New().Append("test").Build()
	if result != "test" {
		t.Errorf("expected \"test\", got %q", result)
	}
}

func TestStringWeaver_Clone(t *testing.T) {
	original := New()
	original.Append("original")
	clone := original.Clone()

	// Clone should have the same content.
	if clone.Build() != "original" {
		t.Errorf("clone should contain \"original\", got %q", clone.Build())
	}

	// Mutating original should not affect clone.
	original.Append(" modified")
	if clone.Build() != "original" {
		t.Errorf("clone should still be \"original\" after mutating original, got %q", clone.Build())
	}
}

func TestStringWeaver_Inspect(t *testing.T) {
	var inspected string
	New().Append("hello").Inspect(func(current string) {
		inspected = current
	}).Append(" world")

	if inspected != "hello" {
		t.Errorf("Inspect should have captured \"hello\", got %q", inspected)
	}
}

// ---------------------------------------------------------------------------
// Weaver Interface Compliance
// ---------------------------------------------------------------------------

func TestStringWeaver_ImplementsWeaver(t *testing.T) {
	var _ Weaver = New()
	var _ Weaver = From("test")
	var _ Weaver = NewWithCapacity(64)
}

func TestSafeStringWeaver_ImplementsWeaver(t *testing.T) {
	var _ Weaver = NewSafe()
	var _ Weaver = SafeFrom("test")
	var _ Weaver = NewSafeWithCapacity(64)
}

// ---------------------------------------------------------------------------
// SafeStringWeaver Tests (thread-safe)
// ---------------------------------------------------------------------------

func TestNewSafe(t *testing.T) {
	sw := NewSafe()
	if sw == nil {
		t.Fatal("NewSafe() returned nil")
	}
	if sw.Len() != 0 {
		t.Errorf("NewSafe() should have length 0, got %d", sw.Len())
	}
}

func TestNewSafeWithCapacity(t *testing.T) {
	sw := NewSafeWithCapacity(1024)
	if sw == nil {
		t.Fatal("NewSafeWithCapacity() returned nil")
	}
	if sw.Cap() < 1024 {
		t.Errorf("NewSafeWithCapacity(1024) should have capacity >= 1024, got %d", sw.Cap())
	}
}

func TestSafeFrom(t *testing.T) {
	sw := SafeFrom("hello")
	if sw.String() != "hello" {
		t.Errorf("SafeFrom(\"hello\") should produce \"hello\", got %q", sw.String())
	}
}

func TestSafeStringWeaver_Append(t *testing.T) {
	result := NewSafe().Append("hello").Append(" ").Append("world").Build()
	if result != "hello world" {
		t.Errorf("expected \"hello world\", got %q", result)
	}
}

func TestSafeStringWeaver_AppendF(t *testing.T) {
	result := NewSafe().AppendF("value: %d", 42).Build()
	if result != "value: 42" {
		t.Errorf("expected \"value: 42\", got %q", result)
	}
}

func TestSafeStringWeaver_NumericAppends(t *testing.T) {
	tests := []struct {
		name     string
		build    func() string
		expected string
	}{
		{"AppendInt", func() string { return NewSafe().AppendInt(123).Build() }, "123"},
		{"AppendInt8", func() string { return NewSafe().AppendInt8(8).Build() }, "8"},
		{"AppendInt16", func() string { return NewSafe().AppendInt16(16).Build() }, "16"},
		{"AppendInt32", func() string { return NewSafe().AppendInt32(32).Build() }, "32"},
		{"AppendInt64", func() string { return NewSafe().AppendInt64(64).Build() }, "64"},
		{"AppendUint", func() string { return NewSafe().AppendUint(123).Build() }, "123"},
		{"AppendUint8", func() string { return NewSafe().AppendUint8(8).Build() }, "8"},
		{"AppendUint16", func() string { return NewSafe().AppendUint16(16).Build() }, "16"},
		{"AppendUint32", func() string { return NewSafe().AppendUint32(32).Build() }, "32"},
		{"AppendUint64", func() string { return NewSafe().AppendUint64(64).Build() }, "64"},
		{"AppendUintptr", func() string { return NewSafe().AppendUintptr(100).Build() }, "100"},
		{"AppendBool/true", func() string { return NewSafe().AppendBool(true).Build() }, "true"},
		{"AppendBool/false", func() string { return NewSafe().AppendBool(false).Build() }, "false"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.build()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSafeStringWeaver_Whitespace(t *testing.T) {
	tests := []struct {
		name     string
		build    func() string
		expected string
	}{
		{"Space", func() string { return NewSafe().Append("a").Space().Append("b").Build() }, "a b"},
		{"Spaces", func() string { return NewSafe().Append("a").Spaces(3).Append("b").Build() }, "a   b"},
		{"Tab", func() string { return NewSafe().Append("a").Tab().Append("b").Build() }, "a\tb"},
		{"Tabs", func() string { return NewSafe().Append("a").Tabs(2).Append("b").Build() }, "a\t\tb"},
		{"NewLine", func() string { return NewSafe().Append("a").NewLine().Append("b").Build() }, "a\nb"},
		{"NewLines", func() string { return NewSafe().Append("a").NewLines(2).Append("b").Build() }, "a\n\nb"},
		{"Line", func() string { return NewSafe().Line("hello").Build() }, "hello\n"},
		{"LineF", func() string { return NewSafe().LineF("id: %d", 42).Build() }, "id: 42\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.build()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSafeStringWeaver_Formatting(t *testing.T) {
	tests := []struct {
		name     string
		build    func() string
		expected string
	}{
		{"Quote", func() string { return NewSafe().Quote("text").Build() }, `"text"`},
		{"SingleQuote", func() string { return NewSafe().SingleQuote("text").Build() }, "'text'"},
		{"Parenthesize", func() string { return NewSafe().Parenthesize("expr").Build() }, "(expr)"},
		{"Bracket", func() string { return NewSafe().Bracket("0").Build() }, "[0]"},
		{"Brace", func() string { return NewSafe().Brace("body").Build() }, "{body}"},
		{"Wrap", func() string { return NewSafe().Wrap("<", "tag", ">").Build() }, "<tag>"},
		{"Indent", func() string { return NewSafe().Indent(2, "text").Build() }, "    text"},
		{"IndentLine", func() string { return NewSafe().IndentLine(1, "item").Build() }, "  item\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.build()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSafeStringWeaver_Punctuation(t *testing.T) {
	tests := []struct {
		name     string
		build    func() string
		expected string
	}{
		{"Comma", func() string { return NewSafe().Append("a").Comma().Append("b").Build() }, "a,b"},
		{"Dot", func() string { return NewSafe().Append("end").Dot().Build() }, "end."},
		{"Colon", func() string { return NewSafe().Append("k").Colon().Append("v").Build() }, "k:v"},
		{"Semicolon", func() string { return NewSafe().Append("x").Semicolon().Build() }, "x;"},
		{"Equals", func() string { return NewSafe().Append("k").Equals().Append("v").Build() }, "k=v"},
		{"Arrow", func() string { return NewSafe().Append("a").Arrow().Append("b").Build() }, "a->b"},
		{"FatArrow", func() string { return NewSafe().Append("a").FatArrow().Append("b").Build() }, "a=>b"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.build()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSafeStringWeaver_Repeat(t *testing.T) {
	result := NewSafe().Repeat("*", 5).Build()
	if result != "*****" {
		t.Errorf("expected \"*****\", got %q", result)
	}
}

func TestSafeStringWeaver_Join(t *testing.T) {
	result := NewSafe().Join(", ", "a", "b", "c").Build()
	if result != "a, b, c" {
		t.Errorf("expected \"a, b, c\", got %q", result)
	}
}

func TestSafeStringWeaver_AppendIf(t *testing.T) {
	resultTrue := NewSafe().Append("a").AppendIf(true, "b").Build()
	resultFalse := NewSafe().Append("a").AppendIf(false, "b").Build()
	if resultTrue != "ab" {
		t.Errorf("expected \"ab\", got %q", resultTrue)
	}
	if resultFalse != "a" {
		t.Errorf("expected \"a\", got %q", resultFalse)
	}
}

func TestSafeStringWeaver_AppendIfF(t *testing.T) {
	result := NewSafe().AppendIfF(true, "n=%d", 42).Build()
	if result != "n=42" {
		t.Errorf("expected \"n=42\", got %q", result)
	}
}

func TestSafeStringWeaver_When(t *testing.T) {
	sw := NewSafe()
	sw.Append("a")
	result := sw.When(true, func(sw *SafeStringWeaver) {
		sw.Append("b")
	}).Build()
	if result != "ab" {
		t.Errorf("expected \"ab\", got %q", result)
	}
}

func TestSafeStringWeaver_Unless(t *testing.T) {
	sw := NewSafe()
	sw.Append("a")
	result := sw.Unless(false, func(sw *SafeStringWeaver) {
		sw.Append("b")
	}).Build()
	if result != "ab" {
		t.Errorf("expected \"ab\", got %q", result)
	}
}

func TestSafeStringWeaver_Each(t *testing.T) {
	result := NewSafe().Each([]string{"a", "b", "c"}, func(sw *SafeStringWeaver, item string) {
		sw.Append(item)
	}).Build()
	if result != "abc" {
		t.Errorf("expected \"abc\", got %q", result)
	}
}

func TestSafeStringWeaver_Reset(t *testing.T) {
	sw := NewSafe()
	sw.Append("hello")
	sw.Reset()
	if sw.Len() != 0 {
		t.Errorf("after Reset(), length should be 0, got %d", sw.Len())
	}
}

func TestSafeStringWeaver_Clone(t *testing.T) {
	original := NewSafe()
	original.Append("original")
	clone := original.Clone()

	if clone.Build() != "original" {
		t.Errorf("clone should contain \"original\", got %q", clone.Build())
	}

	original.Append(" modified")
	if clone.Build() != "original" {
		t.Errorf("clone should still be \"original\" after mutating original, got %q", clone.Build())
	}
}

func TestSafeStringWeaver_Inspect(t *testing.T) {
	var inspected string
	NewSafe().Append("hello").Inspect(func(current string) {
		inspected = current
	}).Append(" world")

	if inspected != "hello" {
		t.Errorf("Inspect should have captured \"hello\", got %q", inspected)
	}
}

// ---------------------------------------------------------------------------
// Concurrency Tests
// ---------------------------------------------------------------------------

func TestSafeStringWeaver_ConcurrentAppend(t *testing.T) {
	sw := NewSafe()
	var wg sync.WaitGroup
	n := 100

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sw.Append("x")
		}()
	}

	wg.Wait()
	if sw.Len() != n {
		t.Errorf("expected length %d after %d concurrent appends, got %d", n, n, sw.Len())
	}
}

func TestSafeStringWeaver_ConcurrentLine(t *testing.T) {
	sw := NewSafe()
	var wg sync.WaitGroup
	n := 50

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sw.LineF("[%d] log entry", id)
		}(i)
	}

	wg.Wait()
	lines := strings.Split(strings.TrimRight(sw.Build(), "\n"), "\n")
	if len(lines) != n {
		t.Errorf("expected %d lines, got %d", n, len(lines))
	}
}

func TestSafeStringWeaver_ConcurrentClone(t *testing.T) {
	template := SafeFrom("base:")
	var wg sync.WaitGroup
	n := 50
	results := make([]string, n)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			clone := template.Clone()
			clone.AppendF("%d", id)
			results[id] = clone.Build()
		}(i)
	}

	wg.Wait()
	for i := 0; i < n; i++ {
		expected := fmt.Sprintf("base:%d", i)
		if results[i] != expected {
			t.Errorf("result[%d] = %q, expected %q", i, results[i], expected)
		}
	}
}

func TestSafeStringWeaver_ConcurrentReadWrite(t *testing.T) {
	sw := NewSafe()
	var wg sync.WaitGroup

	// Writers.
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sw.Append("x")
		}()
	}

	// Readers.
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = sw.Len()
			_ = sw.String()
		}()
	}

	wg.Wait()
	// No panic or race condition is the success criterion.
}

// ---------------------------------------------------------------------------
// Weaver Polymorphism Tests
// ---------------------------------------------------------------------------

func TestWeaver_Polymorphism(t *testing.T) {
	buildGreeting := func(w Weaver, name string) string {
		return w.Append("Hello, ").Append(name).Append("!").Build()
	}

	resultUnsafe := buildGreeting(New(), "Alice")
	resultSafe := buildGreeting(NewSafe(), "Bob")

	if resultUnsafe != "Hello, Alice!" {
		t.Errorf("StringWeaver via Weaver: expected \"Hello, Alice!\", got %q", resultUnsafe)
	}
	if resultSafe != "Hello, Bob!" {
		t.Errorf("SafeStringWeaver via Weaver: expected \"Hello, Bob!\", got %q", resultSafe)
	}
}

// ---------------------------------------------------------------------------
// Fluent Chaining Integration Tests
// ---------------------------------------------------------------------------

func TestStringWeaver_FluentChain_SQL(t *testing.T) {
	result := New().
		Append("SELECT ").
		Join(", ", "id", "name", "email").
		Append(" FROM users").
		Append(" WHERE active = true").
		Build()

	expected := "SELECT id, name, email FROM users WHERE active = true"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestStringWeaver_FluentChain_IndentedBlock(t *testing.T) {
	result := New().
		Line("{").
		IndentLine(1, `"name": "John",`).
		IndentLine(1, `"age": 30`).
		Append("}").
		Build()

	expected := "{\n  \"name\": \"John\",\n  \"age\": 30\n}"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestStringWeaver_FluentChain_Conditional(t *testing.T) {
	verbose := true
	result := New().
		Append("result").
		AppendIf(verbose, " (verbose mode)").
		Build()

	if result != "result (verbose mode)" {
		t.Errorf("expected \"result (verbose mode)\", got %q", result)
	}
}

// ---------------------------------------------------------------------------
// JSON string encoding compliance (Fix: json.Marshal instead of strconv.Quote)
// ---------------------------------------------------------------------------

func TestStringWeaver_JSONString_RFC8259(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple string", "hello", `"hello"`},
		{"embedded double quote", `he"llo`, `"he\"llo"`},
		{"backslash", `back\slash`, `"back\\slash"`},
		// Control char \x00: json.Marshal produces \u0000 (valid JSON);
		// strconv.Quote would produce \x00 which is NOT valid JSON.
		{"null byte", "a\x00b", `"a\u0000b"`},
		{"newline", "line\nbreak", `"line\nbreak"`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := New().JSONString(tc.input).Build()
			if got != tc.want {
				t.Errorf("JSONString(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestSafeStringWeaver_JSONString_RFC8259(t *testing.T) {
	// null byte should be encoded as \u0000 (RFC 8259), not \x00 (Go-style)
	got := NewSafe().JSONString("a\x00b").Build()
	want := `"a\u0000b"`
	if got != want {
		t.Errorf("SafeStringWeaver.JSONString(\"a\\x00b\") = %q, want %q", got, want)
	}
}

// ---------------------------------------------------------------------------
// FromPtr / SafeFromPtr tests (Fix: avoid strings.Builder copy-after-use)
// ---------------------------------------------------------------------------

func TestFromPtr_NilReturnsEmpty(t *testing.T) {
	sw := FromPtr(nil)
	if sw == nil {
		t.Fatal("FromPtr(nil) returned nil")
	}
	if sw.Len() != 0 {
		t.Errorf("FromPtr(nil) should produce empty builder, got len=%d", sw.Len())
	}
}

func TestFromPtr_CopiesContentFromUsedBuilder(t *testing.T) {
	var b strings.Builder
	b.WriteString("hello")

	sw := FromPtr(&b)
	if got := sw.String(); got != "hello" {
		t.Errorf("FromPtr(&usedBuilder) = %q, want %q", got, "hello")
	}
	// Mutating the original builder must not affect the copy.
	b.WriteString(" world")
	if got := sw.String(); got != "hello" {
		t.Errorf("original mutation leaked into FromPtr copy: got %q", got)
	}
}

func TestSafeFromPtr_NilReturnsEmpty(t *testing.T) {
	sw := SafeFromPtr(nil)
	if sw == nil {
		t.Fatal("SafeFromPtr(nil) returned nil")
	}
	if sw.Len() != 0 {
		t.Errorf("SafeFromPtr(nil) should produce empty builder, got len=%d", sw.Len())
	}
}

func TestSafeFromPtr_CopiesContentFromUsedBuilder(t *testing.T) {
	var b strings.Builder
	b.WriteString("safe")

	sw := SafeFromPtr(&b)
	if got := sw.String(); got != "safe" {
		t.Errorf("SafeFromPtr(&usedBuilder) = %q, want %q", got, "safe")
	}
	// Mutating the original builder must not affect the copy.
	b.WriteString("-modified")
	if got := sw.String(); got != "safe" {
		t.Errorf("original mutation leaked into SafeFromPtr copy: got %q", got)
	}
}

// ---------------------------------------------------------------------------
// SafeStringWeaver.AppendIf / AppendIfF consistent locking (Fix)
// ---------------------------------------------------------------------------

func TestSafeStringWeaver_AppendIf_Concurrent(t *testing.T) {
	sw := NewSafe()
	var wg sync.WaitGroup
	n := 200
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			sw.AppendIf(idx%2 == 0, "x")
		}(i)
	}
	wg.Wait()
	// Half of n goroutines (those with even idx) append "x".
	if sw.Len() != n/2 {
		t.Errorf("expected %d chars after concurrent AppendIf, got %d", n/2, sw.Len())
	}
}

func TestSafeStringWeaver_AppendIfF_Concurrent(t *testing.T) {
	sw := NewSafe()
	var wg sync.WaitGroup
	n := 100
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			sw.AppendIfF(idx%2 == 0, "y")
		}(i)
	}
	wg.Wait()
	if sw.Len() != n/2 {
		t.Errorf("expected %d chars after concurrent AppendIfF, got %d", n/2, sw.Len())
	}
}

// ---------------------------------------------------------------------------
// strconv-based append correctness (Fix: allocation-free numeric conversion)
// ---------------------------------------------------------------------------

func TestStringWeaver_AppendInt_Negative(t *testing.T) {
	if got := New().AppendInt(-42).Build(); got != "-42" {
		t.Errorf("AppendInt(-42) = %q, want \"-42\"", got)
	}
}

func TestStringWeaver_AppendFloat32_Compact(t *testing.T) {
	// strconv.FormatFloat with prec=-1 gives the shortest representation.
	got := New().AppendFloat32(3.14).Build()
	if got != "3.14" {
		t.Errorf("AppendFloat32(3.14) = %q, want \"3.14\"", got)
	}
}

func TestStringWeaver_AppendFloat64_Compact(t *testing.T) {
	got := New().AppendFloat64(2.718281828).Build()
	if got != "2.718281828" {
		t.Errorf("AppendFloat64(2.718281828) = %q, want \"2.718281828\"", got)
	}
}

// ---------------------------------------------------------------------------
// strings.Repeat-based methods correctness (Fix)
// ---------------------------------------------------------------------------

func TestStringWeaver_Repeat_ZeroNoop(t *testing.T) {
	if got := New().Repeat("-", 0).Build(); got != "" {
		t.Errorf("Repeat(\"-\", 0) should produce empty string, got %q", got)
	}
}

func TestStringWeaver_Spaces_ZeroNoop(t *testing.T) {
	if got := New().Spaces(0).Build(); got != "" {
		t.Errorf("Spaces(0) should produce empty string, got %q", got)
	}
}

func TestStringWeaver_Indent_ZeroLevelNoop(t *testing.T) {
	if got := New().Indent(0, "text").Build(); got != "text" {
		t.Errorf("Indent(0, \"text\") should produce \"text\", got %q", got)
	}
}

func TestStringWeaver_Indent_NegativeLevelNoop(t *testing.T) {
	if got := New().Indent(-1, "text").Build(); got != "text" {
		t.Errorf("Indent(-1, \"text\") should produce \"text\", got %q", got)
	}
}

// ---------------------------------------------------------------------------
// Benchmarks
// ---------------------------------------------------------------------------

func BenchmarkStringWeaver_Append(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sw := New()
		sw.Append("hello").Space().Append("world")
		_ = sw.Build()
	}
}

func BenchmarkSafeStringWeaver_Append(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sw := NewSafe()
		sw.Append("hello").Space().Append("world")
		_ = sw.Build()
	}
}

func BenchmarkNativeStringsBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var sb strings.Builder
		sb.WriteString("hello")
		sb.WriteByte(' ')
		sb.WriteString("world")
		_ = sb.String()
	}
}

func BenchmarkStringWeaver_LargeChain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sw := NewWithCapacity(256)
		sw.Append("SELECT ").
			Join(", ", "id", "name", "email", "created_at").
			Append(" FROM users").
			Append(" WHERE active = true").
			Append(" ORDER BY created_at DESC").
			Append(" LIMIT 100")
		_ = sw.Build()
	}
}

func BenchmarkSafeStringWeaver_LargeChain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sw := NewSafeWithCapacity(256)
		sw.Append("SELECT ").
			Join(", ", "id", "name", "email", "created_at").
			Append(" FROM users").
			Append(" WHERE active = true").
			Append(" ORDER BY created_at DESC").
			Append(" LIMIT 100")
		_ = sw.Build()
	}
}

func BenchmarkSafeStringWeaver_ConcurrentAppend(b *testing.B) {
	sw := NewSafeWithCapacity(b.N * 5)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sw.Append("hello")
		}
	})
}
