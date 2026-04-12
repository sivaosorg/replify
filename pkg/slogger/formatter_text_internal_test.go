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
	assertEqual(t, ColorNever, f.ColorMode()) // Should also set ColorNever
}

func TestTextFormatter_WithEnableColor(t *testing.T) {
	t.Parallel()

	f := NewTextFormatter(nil).WithDisableColor().WithEnableColor()
	assertFalse(t, f.IsDisableColors())
	assertEqual(t, ColorAuto, f.ColorMode()) // Should reset to ColorAuto
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
// ColorMode Tests
// =============================================================================

func TestTextFormatter_WithColorMode(t *testing.T) {
	t.Parallel()

	t.Run("ColorNever disables colors", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		f := NewTextFormatter(&buf).WithColorMode(ColorNever)
		assertEqual(t, ColorNever, f.ColorMode())
		assertTrue(t, f.IsDisableColors())
	})

	t.Run("ColorAlways enables colors", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		f := NewTextFormatter(&buf).WithColorMode(ColorAlways)
		assertEqual(t, ColorAlways, f.ColorMode())
		assertFalse(t, f.IsDisableColors())
	})

	t.Run("ColorAuto is default", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		f := NewTextFormatter(&buf)
		assertEqual(t, ColorAuto, f.ColorMode())
	})
}

func TestTextFormatter_ColorMode_FormatOutput(t *testing.T) {
	t.Parallel()

	t.Run("ColorNever produces no ANSI codes", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		f := NewTextFormatter(&buf).WithColorMode(ColorNever)
		log := New(WithOutput(&buf), WithFormatter(f))

		entry := &Entry{
			logger:  log,
			time:    time.Now(),
			level:   InfoLevel,
			message: "test message",
		}

		data, err := f.Format(entry)
		assertNoError(t, err)
		out := string(data)

		// Verify no ANSI escape sequences
		assertNotContains(t, out, "\033[")
		assertNotContains(t, out, colorGreen)
		assertNotContains(t, out, colorReset)
		assertContains(t, out, "INFO")
	})

	t.Run("ColorAlways produces ANSI codes", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		f := NewTextFormatter(&buf).WithColorMode(ColorAlways)
		log := New(WithOutput(&buf), WithFormatter(f))

		entry := &Entry{
			logger:  log,
			time:    time.Now(),
			level:   InfoLevel,
			message: "test message",
		}

		data, err := f.Format(entry)
		assertNoError(t, err)
		out := string(data)

		// Verify ANSI escape sequences are present
		assertContains(t, out, "\033[")
		assertContains(t, out, colorReset)
	})

	t.Run("ColorAuto with non-TTY produces no ANSI codes", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer // bytes.Buffer is not a TTY
		f := NewTextFormatter(&buf).WithColorMode(ColorAuto)
		log := New(WithOutput(&buf), WithFormatter(f))

		entry := &Entry{
			logger:  log,
			time:    time.Now(),
			level:   WarnLevel,
			message: "warning message",
		}

		data, err := f.Format(entry)
		assertNoError(t, err)
		out := string(data)

		// bytes.Buffer is not a TTY, so no colours
		assertNotContains(t, out, "\033[")
		assertContains(t, out, "WARN")
	})
}

func TestTextFormatter_ColorMode_NilReceiver(t *testing.T) {
	t.Parallel()

	var f *TextFormatter
	assertEqual(t, ColorAuto, f.ColorMode())
}

// =============================================================================
// StripANSI Tests
// =============================================================================

func TestStripANSI(t *testing.T) {
	t.Parallel()

	t.Run("empty string", func(t *testing.T) {
		t.Parallel()
		assertEqual(t, "", stripANSI(""))
	})

	t.Run("no escape sequences", func(t *testing.T) {
		t.Parallel()
		input := "plain text message"
		assertEqual(t, input, stripANSI(input))
	})

	t.Run("strips color codes", func(t *testing.T) {
		t.Parallel()
		input := colorGreen + "INFO" + colorReset + " message"
		expected := "INFO message"
		assertEqual(t, expected, stripANSI(input))
	})

	t.Run("strips bold", func(t *testing.T) {
		t.Parallel()
		input := colorBold + "BOLD" + colorReset
		expected := "BOLD"
		assertEqual(t, expected, stripANSI(input))
	})

	t.Run("strips multiple sequences", func(t *testing.T) {
		t.Parallel()
		input := colorRed + colorBold + "ERROR" + colorReset + " " + colorYellow + "warning" + colorReset
		expected := "ERROR warning"
		assertEqual(t, expected, stripANSI(input))
	})

	t.Run("preserves unicode", func(t *testing.T) {
		t.Parallel()
		input := colorGreen + "日本語" + colorReset + " 🎉"
		expected := "日本語 🎉"
		assertEqual(t, expected, stripANSI(input))
	})

	t.Run("handles real log line", func(t *testing.T) {
		t.Parallel()
		// Simulate actual colored log output
		input := "2026-04-12 10:30:00 " + colorGreen + colorBold + "INFO " + colorReset + " application started"
		expected := "2026-04-12 10:30:00 INFO  application started"
		assertEqual(t, expected, stripANSI(input))
	})

	// =========================================================================
	// Advanced edge cases
	// =========================================================================

	t.Run("standalone ESC at end of string", func(t *testing.T) {
		t.Parallel()
		// ESC alone without CSI sequence should be preserved (not a valid ANSI sequence)
		input := "text\033"
		expected := "text\033"
		assertEqual(t, expected, stripANSI(input))
	})

	t.Run("ESC followed by non-bracket character", func(t *testing.T) {
		t.Parallel()
		// Non-CSI escape sequence (e.g., \033) followed by something other than [
		input := "text\033)text2"
		expected := "text\033)text2"
		assertEqual(t, expected, stripANSI(input))
	})

	t.Run("incomplete CSI sequence at end", func(t *testing.T) {
		t.Parallel()
		// CSI sequence that starts but has no terminator
		input := "text\033["
		expected := "text"
		assertEqual(t, expected, stripANSI(input))
	})

	t.Run("incomplete CSI with parameters at end", func(t *testing.T) {
		t.Parallel()
		// CSI sequence with parameters but no terminator
		input := "text\033[38;5;196"
		expected := "text"
		assertEqual(t, expected, stripANSI(input))
	})

	t.Run("256-color SGR sequence", func(t *testing.T) {
		t.Parallel()
		// 256-color foreground: ESC[38;5;<n>m
		input := "\033[38;5;196mRED\033[0m"
		expected := "RED"
		assertEqual(t, expected, stripANSI(input))
	})

	t.Run("24-bit true color SGR sequence", func(t *testing.T) {
		t.Parallel()
		// 24-bit RGB foreground: ESC[38;2;<r>;<g>;<b>m
		input := "\033[38;2;255;0;0mTRUECOLOR\033[0m"
		expected := "TRUECOLOR"
		assertEqual(t, expected, stripANSI(input))
	})

	t.Run("cursor movement sequences", func(t *testing.T) {
		t.Parallel()
		// Cursor up (A), down (B), forward (C), back (D)
		input := "\033[2Aup\033[3Bdown\033[4Cforward\033[5Dback"
		expected := "updownforwardback"
		assertEqual(t, expected, stripANSI(input))
	})

	t.Run("cursor position sequence", func(t *testing.T) {
		t.Parallel()
		// Cursor position: ESC[<row>;<col>H
		input := "start\033[10;20Hmiddle\033[1;1Hend"
		expected := "startmiddleend"
		assertEqual(t, expected, stripANSI(input))
	})

	t.Run("erase sequences", func(t *testing.T) {
		t.Parallel()
		// Erase in Display (J) and Erase in Line (K)
		input := "text\033[2Jcleared\033[Kline"
		expected := "textclearedline"
		assertEqual(t, expected, stripANSI(input))
	})

	t.Run("scroll sequences", func(t *testing.T) {
		t.Parallel()
		// Scroll up (S) and scroll down (T)
		input := "text\033[2Sscrolled\033[3T"
		expected := "textscrolled"
		assertEqual(t, expected, stripANSI(input))
	})

	t.Run("multiple ESC in a row", func(t *testing.T) {
		t.Parallel()
		// Multiple ESC characters where only the second forms a valid CSI
		input := "\033\033[31mred\033[0m"
		expected := "\033red"
		assertEqual(t, expected, stripANSI(input))
	})

	t.Run("string with only ANSI codes", func(t *testing.T) {
		t.Parallel()
		// String containing only escape sequences, no visible content
		input := "\033[31m\033[1m\033[0m"
		expected := ""
		assertEqual(t, expected, stripANSI(input))
	})

	t.Run("private mode sequences", func(t *testing.T) {
		t.Parallel()
		// Private mode sequences use ? after CSI (e.g., show/hide cursor)
		input := "\033[?25lhidden\033[?25h"
		expected := "hidden"
		assertEqual(t, expected, stripANSI(input))
	})

	t.Run("SGR with semicolons and multiple attributes", func(t *testing.T) {
		t.Parallel()
		// Multiple SGR attributes: bold;red;underline
		input := "\033[1;31;4mstyledtext\033[0m"
		expected := "styledtext"
		assertEqual(t, expected, stripANSI(input))
	})

	t.Run("CSI with @ terminator", func(t *testing.T) {
		t.Parallel()
		// Insert character: ESC[@
		input := "ab\033[2@cd"
		expected := "abcd"
		assertEqual(t, expected, stripANSI(input))
	})

	t.Run("CSI with tilde terminator", func(t *testing.T) {
		t.Parallel()
		// Function key sequences often end with ~
		input := "text\033[15~more"
		expected := "textmore"
		assertEqual(t, expected, stripANSI(input))
	})

	t.Run("mixed content with newlines and ANSI", func(t *testing.T) {
		t.Parallel()
		input := "\033[32mline1\033[0m\nline2\n\033[31mline3\033[0m"
		expected := "line1\nline2\nline3"
		assertEqual(t, expected, stripANSI(input))
	})

	t.Run("tab characters preserved with ANSI", func(t *testing.T) {
		t.Parallel()
		input := "\033[32mcol1\033[0m\tcol2"
		expected := "col1\tcol2"
		assertEqual(t, expected, stripANSI(input))
	})

	t.Run("very long ANSI sequence", func(t *testing.T) {
		t.Parallel()
		// Long parameter list
		input := "\033[0;1;2;3;4;5;6;7;8;9;10;11;12;13;14;15mtext\033[0m"
		expected := "text"
		assertEqual(t, expected, stripANSI(input))
	})

	t.Run("ESC at various positions", func(t *testing.T) {
		t.Parallel()
		// ESC at start, middle, end positions without forming valid CSI
		input := "\033 text \033"
		expected := "\033 text \033"
		assertEqual(t, expected, stripANSI(input))
	})

	t.Run("binary-like content with ESC", func(t *testing.T) {
		t.Parallel()
		// Ensure non-printable characters are preserved
		input := "text\x00\x01\033[31mred\033[0m\x02\x03"
		expected := "text\x00\x01red\x02\x03"
		assertEqual(t, expected, stripANSI(input))
	})
}

// =============================================================================
// cloneFormatterForFile Tests
// =============================================================================

func TestCloneFormatterForFile(t *testing.T) {
	t.Parallel()

	t.Run("TextFormatter gets ColorNever", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		original := NewTextFormatter(&buf).
			WithTimeFormat("2006-01-02").
			WithEnableCaller().
			WithColorMode(ColorAlways) // Even if ColorAlways, file should be ColorNever

		cloned := cloneFormatterForFile(original)
		tf, ok := cloned.(*TextFormatter)
		assertTrue(t, ok)
		assertEqual(t, ColorNever, tf.ColorMode())
		assertTrue(t, tf.IsDisableColors())
		assertEqual(t, "2006-01-02", tf.TimeFormat())
		assertTrue(t, tf.IsEnableCaller())
	})

	t.Run("JSONFormatter returned unchanged", func(t *testing.T) {
		t.Parallel()
		original := NewJSONFormatter().WithTimeFormat("2006-01-02")
		cloned := cloneFormatterForFile(original)

		// Should be the same instance for JSONFormatter
		assertEqual(t, original, cloned)
	})
}

// =============================================================================
// Integration Tests
// =============================================================================

func TestTextFormatter_FileOutput_NoColors(t *testing.T) {
	t.Parallel()

	// Simulate the scenario where TextFormatter is used for both stdout and file
	var stdoutBuf bytes.Buffer

	// Formatter for stdout with auto-detect (non-TTY in test)
	stdoutFormatter := NewTextFormatter(&stdoutBuf).WithColorMode(ColorAlways)

	// Formatter for file output should have ColorNever
	fileFormatter := cloneFormatterForFile(stdoutFormatter)

	log := New(WithOutput(&stdoutBuf), WithFormatter(stdoutFormatter))

	entry := &Entry{
		logger:  log,
		time:    time.Now(),
		level:   ErrorLevel,
		message: "error occurred",
	}

	// Format for stdout (colored)
	stdoutData, err := stdoutFormatter.Format(entry)
	assertNoError(t, err)

	// Format for file (no colors)
	fileData, err := fileFormatter.Format(entry)
	assertNoError(t, err)

	// Verify stdout has colors
	assertContains(t, string(stdoutData), "\033[")

	// Verify file has no colors
	assertNotContains(t, string(fileData), "\033[")

	// Both should contain the message
	assertContains(t, string(stdoutData), "error occurred")
	assertContains(t, string(fileData), "error occurred")
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

func BenchmarkStripANSI_NoEscape(b *testing.B) {
	// Common case: no escape sequences
	input := "2026-04-12 10:30:00 INFO  application started with key=value count=42"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = stripANSI(input)
	}
}

func BenchmarkStripANSI_WithColors(b *testing.B) {
	// Typical colored log line
	input := "2026-04-12 10:30:00 \033[32m\033[1mINFO \033[0m application started key=value"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = stripANSI(input)
	}
}

func BenchmarkStripANSI_Heavy(b *testing.B) {
	// Many ANSI sequences
	input := "\033[38;2;255;0;0m\033[1mERROR\033[0m \033[38;5;196mfailed\033[0m \033[33mwarning\033[0m \033[34mdebug\033[0m"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = stripANSI(input)
	}
}
