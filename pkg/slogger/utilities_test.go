package slogger

import (
	"math"
	"strings"
	"testing"
)

// =============================================================================
// trimFilepath Tests
// =============================================================================

func TestTrimFilepath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "long unix path",
			input: "/home/user/project/pkg/slogger/logger.go",
			want:  "slogger/logger.go",
		},
		{
			name:  "windows path",
			input: "C:\\Users\\user\\project\\pkg\\slogger\\logger.go",
			want:  "slogger\\logger.go",
		},
		{
			name:  "short path single segment",
			input: "logger.go",
			want:  "logger.go",
		},
		{
			name:  "two segments",
			input: "slogger/logger.go",
			want:  "slogger/logger.go",
		},
		{
			name:  "empty path",
			input: "",
			want:  "",
		},
		{
			name:  "root only",
			input: "/",
			want:  "/",
		},
		{
			name:  "trailing slash",
			input: "/home/user/project/",
			want:  "project/",
		},
		{
			name:  "one slash",
			input: "pkg/file.go",
			want:  "pkg/file.go",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := trimFilepath(tt.input)
			assertEqual(t, tt.want, got)
		})
	}
}

// =============================================================================
// levelColor Tests
// =============================================================================

func TestLevelColor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		level Level
		want  string
	}{
		{TraceLevel, colorCyan},
		{DebugLevel, colorBlue},
		{InfoLevel, colorGreen},
		{WarnLevel, colorYellow},
		{ErrorLevel, colorRed},
		{FatalLevel, colorRed},
		{PanicLevel, colorRed},
		{Level(100), colorReset}, // Unknown level
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.level.String(), func(t *testing.T) {
			t.Parallel()
			got := levelColor(tt.level)
			assertEqual(t, tt.want, got)
		})
	}
}

// =============================================================================
// levelPad Tests
// =============================================================================

func TestLevelPad(t *testing.T) {
	t.Parallel()

	levels := []Level{TraceLevel, DebugLevel, InfoLevel, WarnLevel, ErrorLevel, FatalLevel, PanicLevel}

	for _, level := range levels {
		level := level
		t.Run(level.String(), func(t *testing.T) {
			t.Parallel()
			got := levelPad(level)
			// All padded strings should have the same length
			assertGreaterOrEqual(t, len(got), levelPadWidth)
		})
	}
}

// =============================================================================
// shouldQuoting Tests
// =============================================================================

func TestShouldQuoting(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "simple string", input: "hello", want: false},
		{name: "contains space", input: "hello world", want: true},
		{name: "contains equals", input: "key=value", want: true},
		{name: "contains quote", input: `say "hi"`, want: true},
		{name: "contains backslash", input: `path\to`, want: true},
		{name: "empty string", input: "", want: false},
		{name: "unicode no special", input: "日本語", want: false},
		{name: "unicode with space", input: "日本語 テスト", want: true},
		{name: "numbers only", input: "12345", want: false},
		{name: "special chars", input: "@#$%^&*()", want: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := shouldQuoting(tt.input)
			assertEqual(t, tt.want, got)
		})
	}
}

// =============================================================================
// itoa Tests
// =============================================================================

func TestItoa(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input int
		want  string
	}{
		{name: "zero", input: 0, want: "0"},
		{name: "positive", input: 42, want: "42"},
		{name: "negative", input: -42, want: "-42"},
		{name: "large positive", input: 123456789, want: "123456789"},
		{name: "large negative", input: -123456789, want: "-123456789"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := itoa(tt.input)
			assertEqual(t, tt.want, got)
		})
	}
}

// =============================================================================
// itoa64 Tests
// =============================================================================

func TestItoa64(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input int64
		want  string
	}{
		{name: "zero", input: 0, want: "0"},
		{name: "positive", input: 42, want: "42"},
		{name: "negative", input: -42, want: "-42"},
		{name: "max int64", input: math.MaxInt64, want: "9223372036854775807"},
		{name: "min int64", input: math.MinInt64, want: "-9223372036854775808"},
		{name: "large positive", input: 9876543210, want: "9876543210"},
		{name: "large negative", input: -9876543210, want: "-9876543210"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := itoa64(tt.input)
			assertEqual(t, tt.want, got)
		})
	}
}

// =============================================================================
// utoa64 Tests
// =============================================================================

func TestUtoa64(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input uint64
		want  string
	}{
		{name: "zero", input: 0, want: "0"},
		{name: "positive", input: 42, want: "42"},
		{name: "large", input: 9876543210, want: "9876543210"},
		{name: "max uint64", input: math.MaxUint64, want: "18446744073709551615"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := utoa64(tt.input)
			assertEqual(t, tt.want, got)
		})
	}
}

// =============================================================================
// writeJSONKey Tests
// =============================================================================

func TestWriteJSONKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		key  string
		want string
	}{
		{name: "simple key", key: "message", want: `"message"`},
		{name: "empty key", key: "", want: `""`},
		{name: "key with special chars", key: "msg\nkey", want: `"msg\nkey"`},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var b strings.Builder
			writeJSONKey(&b, tt.key)
			assertEqual(t, tt.want, b.String())
		})
	}
}

// =============================================================================
// writeJSONString Tests
// =============================================================================

func TestWriteJSONString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "simple string", input: "hello", want: `"hello"`},
		{name: "empty string", input: "", want: `""`},
		{name: "string with newline", input: "line1\nline2", want: `"line1\nline2"`},
		{name: "string with tab", input: "col1\tcol2", want: `"col1\tcol2"`},
		{name: "string with quote", input: `say "hi"`, want: `"say \"hi\""`},
		{name: "unicode", input: "日本語🎉", want: `"日本語🎉"`},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var b strings.Builder
			writeJSONString(&b, tt.input)
			assertEqual(t, tt.want, b.String())
		})
	}
}

// =============================================================================
// writeJSONValue Tests
// =============================================================================

func TestWriteJSONValue(t *testing.T) {
	t.Parallel()

	t.Run("string type", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		f := String("key", "value")
		writeJSONValue(&b, &f)
		assertEqual(t, `"value"`, b.String())
	})

	t.Run("int64 type", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		f := Int64("key", 12345)
		writeJSONValue(&b, &f)
		assertEqual(t, "12345", b.String())
	})

	t.Run("uint64 type", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		f := Uint64("key", 12345)
		writeJSONValue(&b, &f)
		assertEqual(t, "12345", b.String())
	})

	t.Run("float64 type", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		f := Float64("key", 3.14)
		writeJSONValue(&b, &f)
		assertEqual(t, "3.14", b.String())
	})

	t.Run("bool true", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		f := Bool("key", true)
		writeJSONValue(&b, &f)
		assertEqual(t, "true", b.String())
	})

	t.Run("bool false", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		f := Bool("key", false)
		writeJSONValue(&b, &f)
		assertEqual(t, "false", b.String())
	})

	t.Run("error nil", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		f := Err(nil)
		writeJSONValue(&b, &f)
		assertEqual(t, "null", b.String())
	})
}

// =============================================================================
// Fuzz Tests
// =============================================================================

func FuzzItoa64(f *testing.F) {
	f.Add(int64(0))
	f.Add(int64(1))
	f.Add(int64(-1))
	f.Add(int64(math.MaxInt64))
	f.Add(int64(math.MinInt64))
	f.Add(int64(123456789))

	f.Fuzz(func(t *testing.T, n int64) {
		result := itoa64(n)

		// Verify result is non-empty
		if len(result) == 0 {
			t.Error("itoa64 returned empty string")
		}

		// Verify negative numbers start with -
		if n < 0 && result[0] != '-' {
			t.Errorf("negative number %d doesn't start with -: %s", n, result)
		}

		// Verify non-negative numbers don't start with -
		if n >= 0 && result[0] == '-' {
			t.Errorf("non-negative number %d starts with -: %s", n, result)
		}
	})
}

func FuzzUtoa64(f *testing.F) {
	f.Add(uint64(0))
	f.Add(uint64(1))
	f.Add(uint64(math.MaxUint64))
	f.Add(uint64(123456789))

	f.Fuzz(func(t *testing.T, n uint64) {
		result := utoa64(n)

		// Verify result is non-empty
		if len(result) == 0 {
			t.Error("utoa64 returned empty string")
		}

		// Verify result contains only digits
		for _, c := range result {
			if c < '0' || c > '9' {
				t.Errorf("utoa64(%d) contains non-digit: %s", n, result)
			}
		}
	})
}

// =============================================================================
// Benchmarks
// =============================================================================

func BenchmarkTrimFilepath(b *testing.B) {
	path := "/home/user/project/pkg/slogger/logger.go"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = trimFilepath(path)
	}
}

func BenchmarkLevelColor(b *testing.B) {
	levels := []Level{TraceLevel, DebugLevel, InfoLevel, WarnLevel, ErrorLevel, FatalLevel, PanicLevel}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = levelColor(levels[i%len(levels)])
	}
}

func BenchmarkLevelPad(b *testing.B) {
	levels := []Level{TraceLevel, DebugLevel, InfoLevel, WarnLevel, ErrorLevel, FatalLevel, PanicLevel}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = levelPad(levels[i%len(levels)])
	}
}

func BenchmarkShouldQuoting(b *testing.B) {
	strings := []string{"hello", "hello world", "key=value", `say "hi"`}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = shouldQuoting(strings[i%len(strings)])
	}
}

func BenchmarkItoa64(b *testing.B) {
	nums := []int64{0, 42, -42, math.MaxInt64, math.MinInt64}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = itoa64(nums[i%len(nums)])
	}
}

func BenchmarkUtoa64(b *testing.B) {
	nums := []uint64{0, 42, math.MaxUint64}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = utoa64(nums[i%len(nums)])
	}
}

func BenchmarkWriteJSONString(b *testing.B) {
	var builder strings.Builder
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder.Reset()
		writeJSONString(&builder, "benchmark test string")
	}
}
