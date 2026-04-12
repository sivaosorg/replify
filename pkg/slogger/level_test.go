package slogger

import (
	"strings"
	"testing"
)

// =============================================================================
// ParseLevel Tests
// =============================================================================

func TestParseLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    Level
		wantErr bool
	}{
		// Trace level
		{name: "trace lowercase", input: "trace", want: TraceLevel, wantErr: false},
		{name: "trace uppercase", input: "TRACE", want: TraceLevel, wantErr: false},
		{name: "trace mixed", input: "Trace", want: TraceLevel, wantErr: false},
		{name: "trace with spaces", input: "  trace  ", want: TraceLevel, wantErr: false},

		// Debug level
		{name: "debug lowercase", input: "debug", want: DebugLevel, wantErr: false},
		{name: "debug uppercase", input: "DEBUG", want: DebugLevel, wantErr: false},
		{name: "debug mixed", input: "DeBuG", want: DebugLevel, wantErr: false},

		// Info level
		{name: "info lowercase", input: "info", want: InfoLevel, wantErr: false},
		{name: "info uppercase", input: "INFO", want: InfoLevel, wantErr: false},

		// Warn level
		{name: "warn lowercase", input: "warn", want: WarnLevel, wantErr: false},
		{name: "warn uppercase", input: "WARN", want: WarnLevel, wantErr: false},
		{name: "warning lowercase", input: "warning", want: WarnLevel, wantErr: false},
		{name: "warning uppercase", input: "WARNING", want: WarnLevel, wantErr: false},

		// Error level
		{name: "error lowercase", input: "error", want: ErrorLevel, wantErr: false},
		{name: "error uppercase", input: "ERROR", want: ErrorLevel, wantErr: false},

		// Fatal level
		{name: "fatal lowercase", input: "fatal", want: FatalLevel, wantErr: false},
		{name: "fatal uppercase", input: "FATAL", want: FatalLevel, wantErr: false},

		// Panic level
		{name: "panic lowercase", input: "panic", want: PanicLevel, wantErr: false},
		{name: "panic uppercase", input: "PANIC", want: PanicLevel, wantErr: false},

		// Error cases
		{name: "empty string", input: "", want: TraceLevel, wantErr: true},
		{name: "whitespace only", input: "   ", want: TraceLevel, wantErr: true},
		{name: "unknown level", input: "unknown", want: TraceLevel, wantErr: true},
		{name: "invalid level", input: "verbose", want: TraceLevel, wantErr: true},
		{name: "numeric string", input: "123", want: TraceLevel, wantErr: true},
		{name: "special characters", input: "@#$%", want: TraceLevel, wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseLevel(tt.input)
			if tt.wantErr {
				assertError(t, err)
				assertContains(t, err.Error(), "unknown log level")
				return
			}
			assertNoError(t, err)
			assertEqual(t, tt.want, got)
		})
	}
}

// =============================================================================
// Level.IsEnabled Tests
// =============================================================================

func TestLevel_IsEnabled(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		level  Level
		min    Level
		result bool
	}{
		// Same level - should be enabled
		{name: "trace >= trace", level: TraceLevel, min: TraceLevel, result: true},
		{name: "debug >= debug", level: DebugLevel, min: DebugLevel, result: true},
		{name: "info >= info", level: InfoLevel, min: InfoLevel, result: true},
		{name: "warn >= warn", level: WarnLevel, min: WarnLevel, result: true},
		{name: "error >= error", level: ErrorLevel, min: ErrorLevel, result: true},
		{name: "fatal >= fatal", level: FatalLevel, min: FatalLevel, result: true},
		{name: "panic >= panic", level: PanicLevel, min: PanicLevel, result: true},

		// Higher level - should be enabled
		{name: "info >= trace", level: InfoLevel, min: TraceLevel, result: true},
		{name: "error >= info", level: ErrorLevel, min: InfoLevel, result: true},
		{name: "panic >= debug", level: PanicLevel, min: DebugLevel, result: true},
		{name: "fatal >= warn", level: FatalLevel, min: WarnLevel, result: true},

		// Lower level - should be disabled
		{name: "trace < info", level: TraceLevel, min: InfoLevel, result: false},
		{name: "debug < warn", level: DebugLevel, min: WarnLevel, result: false},
		{name: "info < error", level: InfoLevel, min: ErrorLevel, result: false},
		{name: "warn < fatal", level: WarnLevel, min: FatalLevel, result: false},
		{name: "error < panic", level: ErrorLevel, min: PanicLevel, result: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.level.IsEnabled(tt.min)
			assertEqual(t, tt.result, got)
		})
	}
}

// =============================================================================
// Level.String Tests
// =============================================================================

func TestLevel_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		level Level
		want  string
	}{
		{TraceLevel, "TRACE"},
		{DebugLevel, "DEBUG"},
		{InfoLevel, "INFO"},
		{WarnLevel, "WARN"},
		{ErrorLevel, "ERROR"},
		{FatalLevel, "FATAL"},
		{PanicLevel, "PANIC"},
		{Level(100), "UNKNOWN"},  // Unknown level
		{Level(-1), "UNKNOWN"},   // Negative level
		{Level(999), "UNKNOWN"},  // Very high level
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()
			got := tt.level.String()
			assertEqual(t, tt.want, got)
		})
	}
}

// =============================================================================
// Level Constants Tests
// =============================================================================

func TestLevelConstants(t *testing.T) {
	t.Parallel()

	// Verify level ordering
	t.Run("level ordering", func(t *testing.T) {
		t.Parallel()
		assertTrue(t, TraceLevel < DebugLevel)
		assertTrue(t, DebugLevel < InfoLevel)
		assertTrue(t, InfoLevel < WarnLevel)
		assertTrue(t, WarnLevel < ErrorLevel)
		assertTrue(t, ErrorLevel < FatalLevel)
		assertTrue(t, FatalLevel < PanicLevel)
	})

	// Verify TraceLevel is zero value
	t.Run("trace is zero", func(t *testing.T) {
		t.Parallel()
		assertEqual(t, Level(0), TraceLevel)
	})
}

// =============================================================================
// Level Edge Cases
// =============================================================================

func TestLevel_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("zero value level", func(t *testing.T) {
		t.Parallel()
		var level Level
		assertEqual(t, "TRACE", level.String())
		assertTrue(t, level.IsEnabled(TraceLevel))
	})

	t.Run("negative level", func(t *testing.T) {
		t.Parallel()
		level := Level(-1)
		assertEqual(t, "UNKNOWN", level.String())
	})

	t.Run("very high level", func(t *testing.T) {
		t.Parallel()
		level := Level(1000)
		assertEqual(t, "UNKNOWN", level.String())
		// Very high levels should always be enabled
		assertTrue(t, level.IsEnabled(PanicLevel))
	})
}

// =============================================================================
// ParseLevel Fuzz Test
// =============================================================================

func FuzzParseLevel(f *testing.F) {
	// Add seed corpus
	f.Add("trace")
	f.Add("DEBUG")
	f.Add("info")
	f.Add("WARN")
	f.Add("warning")
	f.Add("error")
	f.Add("fatal")
	f.Add("panic")
	f.Add("")
	f.Add("unknown")
	f.Add("   info   ")

	f.Fuzz(func(t *testing.T, input string) {
		level, err := ParseLevel(input)

		// Valid inputs
		normalized := strings.TrimSpace(strings.ToUpper(input))
		validLevels := map[string]Level{
			"TRACE":   TraceLevel,
			"DEBUG":   DebugLevel,
			"INFO":    InfoLevel,
			"WARN":    WarnLevel,
			"WARNING": WarnLevel,
			"ERROR":   ErrorLevel,
			"FATAL":   FatalLevel,
			"PANIC":   PanicLevel,
		}

		if expected, ok := validLevels[normalized]; ok {
			assertNoError(t, err)
			assertEqual(t, expected, level)
		} else {
			assertError(t, err)
		}
	})
}

// =============================================================================
// Level Benchmarks
// =============================================================================

func BenchmarkParseLevel(b *testing.B) {
	inputs := []string{"trace", "DEBUG", "info", "WARN", "error", "fatal", "panic"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseLevel(inputs[i%len(inputs)])
	}
}

func BenchmarkLevel_String(b *testing.B) {
	levels := []Level{TraceLevel, DebugLevel, InfoLevel, WarnLevel, ErrorLevel, FatalLevel, PanicLevel}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = levels[i%len(levels)].String()
	}
}

func BenchmarkLevel_IsEnabled(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = InfoLevel.IsEnabled(DebugLevel)
	}
}
