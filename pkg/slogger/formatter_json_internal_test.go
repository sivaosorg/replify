package slogger

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// =============================================================================
// JSONFormatter Constructor Tests
// =============================================================================

func TestNewJSONFormatter(t *testing.T) {
	t.Parallel()

	f := NewJSONFormatter()
	assertNotNil(t, f)
	assertEqual(t, defaultTimeFormat, f.TimeFormat())
	assertEqual(t, defaultJSONTimeKey, f.TimeKey())
	assertEqual(t, defaultJSONLevelKey, f.LevelKey())
	assertEqual(t, defaultJSONMessageKey, f.MessageKey())
	assertEqual(t, defaultJSONCallerKey, f.CallerKey())
	assertEqual(t, defaultJSONNameKey, f.NameKey())
	assertFalse(t, f.IsEnableColor())
}

// =============================================================================
// JSONFormatter Fluent API Tests
// =============================================================================

func TestJSONFormatter_WithTimeFormat(t *testing.T) {
	t.Parallel()

	f := NewJSONFormatter().WithTimeFormat("2006-01-02")
	assertEqual(t, "2006-01-02", f.TimeFormat())
}

func TestJSONFormatter_WithEnableCaller(t *testing.T) {
	t.Parallel()

	f := NewJSONFormatter().WithEnableCaller()
	assertTrue(t, f.IsEnableCaller())
}

func TestJSONFormatter_WithColor(t *testing.T) {
	t.Parallel()

	t.Run("enabled", func(t *testing.T) {
		t.Parallel()
		f := NewJSONFormatter().WithColor(true)
		assertTrue(t, f.IsEnableColor())
	})

	t.Run("disabled", func(t *testing.T) {
		t.Parallel()
		f := NewJSONFormatter().WithColor(true).WithColor(false)
		assertFalse(t, f.IsEnableColor())
	})
}

func TestJSONFormatter_WithEnableColor(t *testing.T) {
	t.Parallel()

	f := NewJSONFormatter().WithEnableColor()
	assertTrue(t, f.IsEnableColor())
}

func TestJSONFormatter_WithTimeKey(t *testing.T) {
	t.Parallel()

	f := NewJSONFormatter().WithTimeKey("timestamp")
	assertEqual(t, "timestamp", f.TimeKey())
}

func TestJSONFormatter_WithLevelKey(t *testing.T) {
	t.Parallel()

	f := NewJSONFormatter().WithLevelKey("severity")
	assertEqual(t, "severity", f.LevelKey())
}

func TestJSONFormatter_WithMessageKey(t *testing.T) {
	t.Parallel()

	f := NewJSONFormatter().WithMessageKey("message")
	assertEqual(t, "message", f.MessageKey())
}

func TestJSONFormatter_WithCallerKey(t *testing.T) {
	t.Parallel()

	f := NewJSONFormatter().WithCallerKey("source")
	assertEqual(t, "source", f.CallerKey())
}

func TestJSONFormatter_WithNameKey(t *testing.T) {
	t.Parallel()

	f := NewJSONFormatter().WithNameKey("logger")
	assertEqual(t, "logger", f.NameKey())
}

// =============================================================================
// JSONFormatter Accessor Tests
// =============================================================================

func TestJSONFormatter_TimeFormat(t *testing.T) {
	t.Parallel()

	t.Run("returns format", func(t *testing.T) {
		t.Parallel()
		f := &JSONFormatter{timeFormat: "2006-01-02"}
		assertEqual(t, "2006-01-02", f.TimeFormat())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var f *JSONFormatter
		assertEqual(t, "", f.TimeFormat())
	})
}

func TestJSONFormatter_IsEnableCaller(t *testing.T) {
	t.Parallel()

	t.Run("enabled", func(t *testing.T) {
		t.Parallel()
		f := &JSONFormatter{enableCaller: true}
		assertTrue(t, f.IsEnableCaller())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var f *JSONFormatter
		assertFalse(t, f.IsEnableCaller())
	})
}

func TestJSONFormatter_IsEnableColor(t *testing.T) {
	t.Parallel()

	t.Run("enabled", func(t *testing.T) {
		t.Parallel()
		f := &JSONFormatter{enableColor: true}
		assertTrue(t, f.IsEnableColor())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var f *JSONFormatter
		assertFalse(t, f.IsEnableColor())
	})
}

func TestJSONFormatter_TimeKey(t *testing.T) {
	t.Parallel()

	t.Run("returns key", func(t *testing.T) {
		t.Parallel()
		f := &JSONFormatter{timeKey: "timestamp"}
		assertEqual(t, "timestamp", f.TimeKey())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var f *JSONFormatter
		assertEqual(t, "", f.TimeKey())
	})
}

func TestJSONFormatter_LevelKey(t *testing.T) {
	t.Parallel()

	t.Run("returns key", func(t *testing.T) {
		t.Parallel()
		f := &JSONFormatter{levelKey: "severity"}
		assertEqual(t, "severity", f.LevelKey())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var f *JSONFormatter
		assertEqual(t, "", f.LevelKey())
	})
}

func TestJSONFormatter_MessageKey(t *testing.T) {
	t.Parallel()

	t.Run("returns key", func(t *testing.T) {
		t.Parallel()
		f := &JSONFormatter{messageKey: "message"}
		assertEqual(t, "message", f.MessageKey())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var f *JSONFormatter
		assertEqual(t, "", f.MessageKey())
	})
}

func TestJSONFormatter_CallerKey(t *testing.T) {
	t.Parallel()

	t.Run("returns key", func(t *testing.T) {
		t.Parallel()
		f := &JSONFormatter{callerKey: "source"}
		assertEqual(t, "source", f.CallerKey())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var f *JSONFormatter
		assertEqual(t, "", f.CallerKey())
	})
}

func TestJSONFormatter_NameKey(t *testing.T) {
	t.Parallel()

	t.Run("returns key", func(t *testing.T) {
		t.Parallel()
		f := &JSONFormatter{nameKey: "logger"}
		assertEqual(t, "logger", f.NameKey())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var f *JSONFormatter
		assertEqual(t, "", f.NameKey())
	})
}

func TestJSONFormatter_Color(t *testing.T) {
	t.Parallel()

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var f *JSONFormatter
		assertNil(t, f.Color())
	})
}

// =============================================================================
// JSONFormatter Format Tests
// =============================================================================

func TestJSONFormatter_Format(t *testing.T) {
	t.Parallel()

	t.Run("basic format", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		f := NewJSONFormatter()

		log := New(WithOutput(&buf), WithFormatter(f), WithLevel(InfoLevel))
		entry := &Entry{
			logger:  log,
			time:    time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC),
			level:   InfoLevel,
			message: "test message",
		}

		data, err := f.Format(entry)
		assertNoError(t, err)

		var m map[string]interface{}
		err = json.Unmarshal(data, &m)
		assertNoError(t, err)

		assertEqual(t, "INFO", m["level"])
		assertEqual(t, "test message", m["msg"])
	})

	t.Run("with fields", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		f := NewJSONFormatter()

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

		var m map[string]interface{}
		err = json.Unmarshal(data, &m)
		assertNoError(t, err)

		assertEqual(t, "value", m["key"])
		assertEqual(t, float64(42), m["count"])
	})

	t.Run("with caller", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		f := NewJSONFormatter().WithEnableCaller()

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

		var m map[string]interface{}
		err = json.Unmarshal(data, &m)
		assertNoError(t, err)

		assertEqual(t, "test.go:42", m["caller"])
	})

	t.Run("with logger name", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		f := NewJSONFormatter()

		log := New(WithOutput(&buf), WithFormatter(f), WithName("mylogger"))
		entry := &Entry{
			logger:  log,
			time:    time.Now(),
			level:   InfoLevel,
			message: "test",
		}

		data, err := f.Format(entry)
		assertNoError(t, err)

		var m map[string]interface{}
		err = json.Unmarshal(data, &m)
		assertNoError(t, err)

		assertEqual(t, "mylogger", m["name"])
	})

	t.Run("custom keys", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		f := NewJSONFormatter().
			WithTimeKey("timestamp").
			WithLevelKey("severity").
			WithMessageKey("message")

		log := New(WithOutput(&buf), WithFormatter(f))
		entry := &Entry{
			logger:  log,
			time:    time.Now(),
			level:   InfoLevel,
			message: "test",
		}

		data, err := f.Format(entry)
		assertNoError(t, err)

		var m map[string]interface{}
		err = json.Unmarshal(data, &m)
		assertNoError(t, err)

		_, ok := m["timestamp"]
		assertTrue(t, ok)
		assertEqual(t, "INFO", m["severity"])
		assertEqual(t, "test", m["message"])
	})

	t.Run("ends with newline", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		f := NewJSONFormatter()

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
			f := NewJSONFormatter()
			log := New(WithOutput(&buf), WithFormatter(f))

			entry := &Entry{
				logger:  log,
				time:    time.Now(),
				level:   lvl,
				message: "test",
			}

			data, err := f.Format(entry)
			assertNoError(t, err)

			var m map[string]interface{}
			err = json.Unmarshal(data, &m)
			assertNoError(t, err)
			assertEqual(t, expected[i], m["level"])
		}
	})
}

// =============================================================================
// JSONFormatter Field Type Tests
// =============================================================================

func TestJSONFormatter_FormatFieldTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		field    Field
		checkVal func(t *testing.T, m map[string]interface{})
	}{
		{
			name:  "string field",
			field: String("key", "value"),
			checkVal: func(t *testing.T, m map[string]interface{}) {
				assertEqual(t, "value", m["key"])
			},
		},
		{
			name:  "int64 field",
			field: Int64("key", 12345),
			checkVal: func(t *testing.T, m map[string]interface{}) {
				assertEqual(t, float64(12345), m["key"])
			},
		},
		{
			name:  "float64 field",
			field: Float64("key", 3.14),
			checkVal: func(t *testing.T, m map[string]interface{}) {
				assertEqual(t, 3.14, m["key"])
			},
		},
		{
			name:  "bool true field",
			field: Bool("key", true),
			checkVal: func(t *testing.T, m map[string]interface{}) {
				assertEqual(t, true, m["key"])
			},
		},
		{
			name:  "bool false field",
			field: Bool("key", false),
			checkVal: func(t *testing.T, m map[string]interface{}) {
				assertEqual(t, false, m["key"])
			},
		},
		{
			name:  "nil error field",
			field: Err(nil),
			checkVal: func(t *testing.T, m map[string]interface{}) {
				assertNil(t, m["error"])
			},
		},
		{
			name:  "duration field",
			field: Duration("key", 5*time.Second),
			checkVal: func(t *testing.T, m map[string]interface{}) {
				assertEqual(t, "5s", m["key"])
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			f := NewJSONFormatter()
			log := New(WithOutput(&buf), WithFormatter(f))

			entry := &Entry{
				logger:  log,
				time:    time.Now(),
				level:   InfoLevel,
				message: "test",
				fields:  []Field{tt.field},
			}

			data, err := f.Format(entry)
			assertNoError(t, err)

			var m map[string]interface{}
			err = json.Unmarshal(data, &m)
			assertNoError(t, err)

			tt.checkVal(t, m)
		})
	}
}

// =============================================================================
// JSONFormatter Implements Formatter Tests
// =============================================================================

func TestJSONFormatter_ImplementsFormatter(t *testing.T) {
	t.Parallel()

	var _ Formatter = (*JSONFormatter)(nil)
}

// =============================================================================
// JSONFormatter Edge Cases
// =============================================================================

func TestJSONFormatter_EmptyMessage(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	f := NewJSONFormatter()
	log := New(WithOutput(&buf), WithFormatter(f))

	entry := &Entry{
		logger:  log,
		time:    time.Now(),
		level:   InfoLevel,
		message: "",
	}

	data, err := f.Format(entry)
	assertNoError(t, err)

	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	assertNoError(t, err)
	assertEqual(t, "", m["msg"])
}

func TestJSONFormatter_UnicodeMessage(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	f := NewJSONFormatter()
	log := New(WithOutput(&buf), WithFormatter(f))

	entry := &Entry{
		logger:  log,
		time:    time.Now(),
		level:   InfoLevel,
		message: "日本語メッセージ🎉",
	}

	data, err := f.Format(entry)
	assertNoError(t, err)

	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	assertNoError(t, err)
	assertEqual(t, "日本語メッセージ🎉", m["msg"])
}

func TestJSONFormatter_SpecialCharactersInField(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	f := NewJSONFormatter()
	log := New(WithOutput(&buf), WithFormatter(f))

	entry := &Entry{
		logger:  log,
		time:    time.Now(),
		level:   InfoLevel,
		message: "test",
		fields:  []Field{String("data", "line1\nline2\ttab")},
	}

	data, err := f.Format(entry)
	assertNoError(t, err)

	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	assertNoError(t, err)
	assertEqual(t, "line1\nline2\ttab", m["data"])
}

// =============================================================================
// Benchmarks
// =============================================================================

func BenchmarkJSONFormatter_Format(b *testing.B) {
	var buf bytes.Buffer
	f := NewJSONFormatter()
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

func BenchmarkJSONFormatter_FormatWithCaller(b *testing.B) {
	var buf bytes.Buffer
	f := NewJSONFormatter().WithEnableCaller()
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
