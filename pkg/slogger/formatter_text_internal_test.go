package slogger

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

// =============================================================================
// TextFormatter Constructor Tests
// =============================================================================

func TestNewTextFormatter(t *testing.T) {
	t.Parallel()

	t.Run("creates formatter with output", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		f := NewTextFormatter(&buf)
		assertNotNil(t, f)
		assertEqual(t, &buf, f.Output())
	})

	t.Run("nil output", func(t *testing.T) {
		t.Parallel()
		f := NewTextFormatter(nil)
		assertNotNil(t, f)
		assertNil(t, f.Output())
	})
}

// =============================================================================
// TextFormatter Fluent API Tests
// =============================================================================

func TestTextFormatter_WithTimeFormat(t *testing.T) {
	t.Parallel()

	f := NewTextFormatter(nil).WithTimeFormat("2006-01-02")
	assertEqual(t, "2006-01-02", f.TimeFormat())
}

func TestTextFormatter_WithDisableColor(t *testing.T) {
	t.Parallel()

	f := NewTextFormatter(nil).WithDisableColor()
	assertTrue(t, f.IsDisableColors())
}

func TestTextFormatter_WithEnableColor(t *testing.T) {
	t.Parallel()

	f := NewTextFormatter(nil).WithDisableColor().WithEnableColor()
	assertFalse(t, f.IsDisableColors())
}

func TestTextFormatter_WithDisableTimestamp(t *testing.T) {
	t.Parallel()

	f := NewTextFormatter(nil).WithDisableTimestamp()
	assertTrue(t, f.IsDisableTimestamp())
}

func TestTextFormatter_WithEnableCaller(t *testing.T) {
	t.Parallel()

	f := NewTextFormatter(nil).WithEnableCaller()
	assertTrue(t, f.IsEnableCaller())
}

// =============================================================================
// TextFormatter Accessor Tests
// =============================================================================

func TestTextFormatter_TimeFormat(t *testing.T) {
	t.Parallel()

	t.Run("returns format", func(t *testing.T) {
		t.Parallel()
		f := &TextFormatter{timeFormat: "2006-01-02"}
		assertEqual(t, "2006-01-02", f.TimeFormat())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var f *TextFormatter
		assertEqual(t, "", f.TimeFormat())
	})
}

func TestTextFormatter_IsDisableColors(t *testing.T) {
	t.Parallel()

	t.Run("enabled", func(t *testing.T) {
		t.Parallel()
		f := &TextFormatter{disableColors: true}
		assertTrue(t, f.IsDisableColors())
	})

	t.Run("disabled", func(t *testing.T) {
		t.Parallel()
		f := &TextFormatter{disableColors: false}
		assertFalse(t, f.IsDisableColors())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var f *TextFormatter
		assertFalse(t, f.IsDisableColors())
	})
}

func TestTextFormatter_IsDisableTimestamp(t *testing.T) {
	t.Parallel()

	t.Run("enabled", func(t *testing.T) {
		t.Parallel()
		f := &TextFormatter{disableTimestamp: true}
		assertTrue(t, f.IsDisableTimestamp())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var f *TextFormatter
		assertFalse(t, f.IsDisableTimestamp())
	})
}

func TestTextFormatter_IsEnableCaller(t *testing.T) {
	t.Parallel()

	t.Run("enabled", func(t *testing.T) {
		t.Parallel()
		f := &TextFormatter{enableCaller: true}
		assertTrue(t, f.IsEnableCaller())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var f *TextFormatter
		assertFalse(t, f.IsEnableCaller())
	})
}

func TestTextFormatter_Output(t *testing.T) {
	t.Parallel()

	t.Run("returns output", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		f := &TextFormatter{output: &buf}
		assertEqual(t, &buf, f.Output())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var f *TextFormatter
		assertNil(t, f.Output())
	})
}

// =============================================================================
// TextFormatter Format Tests
// =============================================================================

func TestTextFormatter_Format(t *testing.T) {
	t.Parallel()

	t.Run("basic format", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		f := NewTextFormatter(&buf).WithDisableColor()

		log := New(WithOutput(&buf), WithFormatter(f), WithLevel(InfoLevel))
		entry := &Entry{
			logger:  log,
			time:    time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC),
			level:   InfoLevel,
			message: "test message",
		}

		data, err := f.Format(entry)
		assertNoError(t, err)
		out := string(data)

		assertContains(t, out, "INFO")
		assertContains(t, out, "test message")
		assertContains(t, out, "2024-06-15")
	})

	t.Run("with fields", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		f := NewTextFormatter(&buf).WithDisableColor()

		log := New(WithOutput(&buf), WithFormatter(f))
		entry := &Entry{
			logger:  log,
			time:    time.Now(),
			level:   InfoLevel,
			message: "test",
			fields:  []Field{String("key", "value"), Int("count", 42)},
		}

		data, err := f.Format(entry)
		assertNoError(t, err)
		out := string(data)

		assertContains(t, out, "key=value")
		assertContains(t, out, "count=42")
	})

	t.Run("without timestamp", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		f := NewTextFormatter(&buf).WithDisableColor().WithDisableTimestamp()

		log := New(WithOutput(&buf), WithFormatter(f))
		entry := &Entry{
			logger:  log,
			time:    time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC),
			level:   InfoLevel,
			message: "test",
		}

		data, err := f.Format(entry)
		assertNoError(t, err)
		out := string(data)

		assertNotContains(t, out, "2024-06-15")
	})

	t.Run("with caller", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		f := NewTextFormatter(&buf).WithDisableColor().WithEnableCaller()

		log := New(WithOutput(&buf), WithFormatter(f))
		entry := &Entry{
			logger:  log,
			time:    time.Now(),
			level:   InfoLevel,
			message: "test",
			caller:  &CallerInfo{file: "test.go", line: 42},
		}

		data, err := f.Format(entry)
		assertNoError(t, err)
		out := string(data)

		assertContains(t, out, "caller=test.go:42")
	})

	t.Run("with logger name", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		f := NewTextFormatter(&buf).WithDisableColor()

		log := New(WithOutput(&buf), WithFormatter(f), WithName("mylogger"))
		entry := &Entry{
			logger:  log,
			time:    time.Now(),
			level:   InfoLevel,
			message: "test",
		}

		data, err := f.Format(entry)
		assertNoError(t, err)
		out := string(data)

		assertContains(t, out, "[mylogger]")
	})

	t.Run("ends with newline", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		f := NewTextFormatter(&buf).WithDisableColor()

		log := New(WithOutput(&buf), WithFormatter(f))
		entry := &Entry{
			logger:  log,
			time:    time.Now(),
			level:   InfoLevel,
			message: "test",
		}

		data, err := f.Format(entry)
		assertNoError(t, err)
		assertTrue(t, strings.HasSuffix(string(data), "\n"))
	})

	t.Run("all levels", func(t *testing.T) {
		t.Parallel()
		levels := []Level{TraceLevel, DebugLevel, InfoLevel, WarnLevel, ErrorLevel, FatalLevel, PanicLevel}
		expected := []string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL", "PANIC"}

		for i, lvl := range levels {
			var buf bytes.Buffer
			f := NewTextFormatter(&buf).WithDisableColor()
			log := New(WithOutput(&buf), WithFormatter(f))

			entry := &Entry{
				logger:  log,
				time:    time.Now(),
				level:   lvl,
				message: "test",
			}

			data, err := f.Format(entry)
			assertNoError(t, err)
			assertContains(t, string(data), expected[i])
		}
	})

	t.Run("quoted values", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		f := NewTextFormatter(&buf).WithDisableColor()

		log := New(WithOutput(&buf), WithFormatter(f))
		entry := &Entry{
			logger:  log,
			time:    time.Now(),
			level:   InfoLevel,
			message: "test",
			fields:  []Field{String("msg", "hello world")}, // Contains space
		}

		data, err := f.Format(entry)
		assertNoError(t, err)
		out := string(data)

		// Value with space should be quoted
		assertContains(t, out, `"hello world"`)
	})
}

// =============================================================================
// TextFormatter Implements Formatter Tests
// =============================================================================

func TestTextFormatter_ImplementsFormatter(t *testing.T) {
	t.Parallel()

	var _ Formatter = (*TextFormatter)(nil)
}

// =============================================================================
// TextFormatter Edge Cases
// =============================================================================

func TestTextFormatter_EmptyMessage(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	f := NewTextFormatter(&buf).WithDisableColor()
	log := New(WithOutput(&buf), WithFormatter(f))

	entry := &Entry{
		logger:  log,
		time:    time.Now(),
		level:   InfoLevel,
		message: "",
	}

	data, err := f.Format(entry)
	assertNoError(t, err)
	assertContains(t, string(data), "INFO")
}

func TestTextFormatter_UnicodeMessage(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	f := NewTextFormatter(&buf).WithDisableColor()
	log := New(WithOutput(&buf), WithFormatter(f))

	entry := &Entry{
		logger:  log,
		time:    time.Now(),
		level:   InfoLevel,
		message: "日本語メッセージ🎉",
	}

	data, err := f.Format(entry)
	assertNoError(t, err)
	assertContains(t, string(data), "日本語メッセージ🎉")
}

func TestTextFormatter_SpecialCharactersInField(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	f := NewTextFormatter(&buf).WithDisableColor()
	log := New(WithOutput(&buf), WithFormatter(f))

	entry := &Entry{
		logger:  log,
		time:    time.Now(),
		level:   InfoLevel,
		message: "test",
		fields:  []Field{String("data", "line1\nline2")},
	}

	data, err := f.Format(entry)
	assertNoError(t, err)
	// Should be quoted or escaped
	assertNotEmpty(t, string(data))
}

// =============================================================================
// Benchmarks
// =============================================================================

func BenchmarkTextFormatter_Format(b *testing.B) {
	var buf bytes.Buffer
	f := NewTextFormatter(&buf).WithDisableColor()
	log := New(WithOutput(&buf), WithFormatter(f))

	entry := &Entry{
		logger:  log,
		time:    time.Now(),
		level:   InfoLevel,
		message: "benchmark message",
		fields:  []Field{String("key", "value"), Int("count", 42)},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = f.Format(entry)
	}
}

func BenchmarkTextFormatter_FormatWithCaller(b *testing.B) {
	var buf bytes.Buffer
	f := NewTextFormatter(&buf).WithDisableColor().WithEnableCaller()
	log := New(WithOutput(&buf), WithFormatter(f))

	entry := &Entry{
		logger:  log,
		time:    time.Now(),
		level:   InfoLevel,
		message: "benchmark message",
		caller:  &CallerInfo{file: "test.go", line: 42, function: "TestFunc"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = f.Format(entry)
	}
}
