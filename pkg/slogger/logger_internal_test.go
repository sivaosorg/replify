package slogger

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"
	"time"
)

// =============================================================================
// Logger Construction Tests
// =============================================================================

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("default options", func(t *testing.T) {
		t.Parallel()
		log := New()
		assertNotNil(t, log)
		assertEqual(t, InfoLevel, log.Level())
		assertNotNil(t, log.Formatter())
		assertNotNil(t, log.Output())
		assertNotNil(t, log.Hooks())
	})

	t.Run("with custom level", func(t *testing.T) {
		t.Parallel()
		log := New(WithLevel(DebugLevel))
		assertEqual(t, DebugLevel, log.Level())
	})

	t.Run("with custom formatter", func(t *testing.T) {
		t.Parallel()
		formatter := NewJSONFormatter()
		log := New(WithFormatter(formatter))
		assertEqual(t, formatter, log.Formatter())
	})

	t.Run("with custom output", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		log := New(WithOutput(&buf))
		assertEqual(t, &buf, log.Output())
	})

	t.Run("with multiple options", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		formatter := NewJSONFormatter()
		log := New(
			WithLevel(TraceLevel),
			WithOutput(&buf),
			WithFormatter(formatter),
			WithCaller(true),
			WithName("test"),
		)
		assertEqual(t, TraceLevel, log.Level())
		assertEqual(t, &buf, log.Output())
		assertEqual(t, formatter, log.Formatter())
		assertTrue(t, log.IsCaller())
		assertEqual(t, "test", log.Name())
	})

	t.Run("with fields", func(t *testing.T) {
		t.Parallel()
		log := New(WithFields(String("app", "test"), Int("version", 1)))
		fields := log.Fields()
		assertLen(t, fields, 2)
	})

	t.Run("with sampling", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		log := New(
			WithOutput(&buf),
			WithLevel(TraceLevel),
			WithSamplingOpts(NewSamplingOptions().WithFirst(1).WithThereafter(0)),
		)
		assertNotNil(t, log)
	})
}

func TestNewLogger(t *testing.T) {
	t.Parallel()

	log := NewLogger()
	assertNotNil(t, log)
	assertEqual(t, InfoLevel, log.Level())
}

// =============================================================================
// Logger With Tests
// =============================================================================

func TestLogger_With(t *testing.T) {
	t.Parallel()

	t.Run("creates child with fields", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		parent := New(
			WithOutput(&buf),
			WithLevel(TraceLevel),
			WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
		)
		child := parent.With(String("component", "auth"))
		child.Info("test message")

		assertContains(t, buf.String(), "component=auth")
	})

	t.Run("inherits parent fields", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		parent := New(
			WithOutput(&buf),
			WithLevel(TraceLevel),
			WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
			WithFields(String("app", "test")),
		)
		child := parent.With(String("component", "auth"))
		child.Info("test message")

		assertContains(t, buf.String(), "app=test")
		assertContains(t, buf.String(), "component=auth")
	})

	t.Run("shares output", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		parent := New(
			WithOutput(&buf),
			WithLevel(TraceLevel),
			WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
		)
		child := parent.With(String("child", "true"))

		parent.Info("parent message")
		child.Info("child message")

		out := buf.String()
		assertContains(t, out, "parent message")
		assertContains(t, out, "child message")
		assertContains(t, out, "child=true")
	})

	t.Run("child level independent", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		parent := New(
			WithOutput(&buf),
			WithLevel(InfoLevel),
			WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
		)
		child := parent.With(String("child", "true"))
		child.SetLevel(DebugLevel)

		// Parent level should be unchanged
		assertEqual(t, InfoLevel, parent.Level())
		assertEqual(t, DebugLevel, child.Level())
	})

	t.Run("multiple fields", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		log := New(
			WithOutput(&buf),
			WithLevel(TraceLevel),
			WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
		)
		child := log.With(String("a", "1"), Int("b", 2), Bool("c", true))
		child.Info("test")

		out := buf.String()
		assertContains(t, out, "a=1")
		assertContains(t, out, "b=2")
		assertContains(t, out, "c=true")
	})
}

// =============================================================================
// Logger Named Tests
// =============================================================================

func TestLogger_Named(t *testing.T) {
	t.Parallel()

	t.Run("creates named child", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		log := New(
			WithOutput(&buf),
			WithLevel(InfoLevel),
			WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
		)
		named := log.Named("component")
		named.Info("test")

		assertContains(t, buf.String(), "[component]")
	})

	t.Run("extends parent name", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		log := New(
			WithOutput(&buf),
			WithLevel(InfoLevel),
			WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
			WithName("app"),
		)
		child := log.Named("db")
		grandchild := child.Named("reader")
		grandchild.Info("test")

		assertContains(t, buf.String(), "[app.db.reader]")
	})

	t.Run("empty parent name", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		log := New(
			WithOutput(&buf),
			WithLevel(InfoLevel),
			WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
		)
		child := log.Named("standalone")
		child.Info("test")

		assertContains(t, buf.String(), "[standalone]")
	})
}

// =============================================================================
// Logger Level Tests
// =============================================================================

func TestLogger_SetLevel(t *testing.T) {
	t.Parallel()

	t.Run("sets level", func(t *testing.T) {
		t.Parallel()
		log := New()
		log.SetLevel(DebugLevel)
		assertEqual(t, DebugLevel, log.Level())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var log *Logger
		assertNotPanics(t, func() {
			log.SetLevel(DebugLevel)
		})
	})
}

func TestLogger_Level(t *testing.T) {
	t.Parallel()

	t.Run("returns level", func(t *testing.T) {
		t.Parallel()
		log := New(WithLevel(WarnLevel))
		assertEqual(t, WarnLevel, log.Level())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var log *Logger
		assertEqual(t, InfoLevel, log.Level())
	})
}

func TestLogger_GetLevel(t *testing.T) {
	t.Parallel()

	log := New(WithLevel(ErrorLevel))
	assertEqual(t, ErrorLevel, log.GetLevel())
}

// =============================================================================
// Logger Output Tests
// =============================================================================

func TestLogger_SetOutput(t *testing.T) {
	t.Parallel()

	t.Run("changes output", func(t *testing.T) {
		t.Parallel()
		var buf1, buf2 bytes.Buffer
		log := New(
			WithOutput(&buf1),
			WithLevel(InfoLevel),
			WithFormatter(NewTextFormatter(&buf1).WithDisableColor()),
		)
		log.Info("first")
		log.SetOutput(&buf2)
		log.Info("second")

		assertContains(t, buf1.String(), "first")
		assertContains(t, buf2.String(), "second")
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var log *Logger
		assertNotPanics(t, func() {
			log.SetOutput(io.Discard)
		})
	})
}

func TestLogger_Output(t *testing.T) {
	t.Parallel()

	t.Run("returns output", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		log := New(WithOutput(&buf))
		assertEqual(t, &buf, log.Output())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var log *Logger
		assertNil(t, log.Output())
	})
}

// =============================================================================
// Logger Formatter Tests
// =============================================================================

func TestLogger_SetFormatter(t *testing.T) {
	t.Parallel()

	t.Run("changes formatter", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		log := New(
			WithOutput(&buf),
			WithLevel(InfoLevel),
			WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
		)
		log.SetFormatter(NewJSONFormatter())
		log.Info("test")

		// Should now be JSON
		assertContains(t, buf.String(), `"level":"INFO"`)
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var log *Logger
		assertNotPanics(t, func() {
			log.SetFormatter(NewJSONFormatter())
		})
	})
}

func TestLogger_Formatter(t *testing.T) {
	t.Parallel()

	t.Run("returns formatter", func(t *testing.T) {
		t.Parallel()
		formatter := NewJSONFormatter()
		log := New(WithFormatter(formatter))
		assertEqual(t, formatter, log.Formatter())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var log *Logger
		assertNil(t, log.Formatter())
	})
}

// =============================================================================
// Logger Name Tests
// =============================================================================

func TestLogger_Name(t *testing.T) {
	t.Parallel()

	t.Run("returns name", func(t *testing.T) {
		t.Parallel()
		log := New(WithName("mylogger"))
		assertEqual(t, "mylogger", log.Name())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var log *Logger
		assertEqual(t, "", log.Name())
	})
}

// =============================================================================
// Logger Caller Tests
// =============================================================================

func TestLogger_IsCaller(t *testing.T) {
	t.Parallel()

	t.Run("enabled", func(t *testing.T) {
		t.Parallel()
		log := New(WithCaller(true))
		assertTrue(t, log.IsCaller())
	})

	t.Run("disabled", func(t *testing.T) {
		t.Parallel()
		log := New(WithCaller(false))
		assertFalse(t, log.IsCaller())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var log *Logger
		assertFalse(t, log.IsCaller())
	})
}

func TestLogger_CallerSkip(t *testing.T) {
	t.Parallel()

	t.Run("returns skip", func(t *testing.T) {
		t.Parallel()
		log := New(WithCallerSkip(3))
		assertEqual(t, 3, log.CallerSkip())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var log *Logger
		assertEqual(t, 0, log.CallerSkip())
	})
}

// =============================================================================
// Logger Fields Tests
// =============================================================================

func TestLogger_Fields(t *testing.T) {
	t.Parallel()

	t.Run("returns fields", func(t *testing.T) {
		t.Parallel()
		log := New(WithFields(String("a", "1"), Int("b", 2)))
		fields := log.Fields()
		assertLen(t, fields, 2)
	})

	t.Run("returns copy", func(t *testing.T) {
		t.Parallel()
		log := New(WithFields(String("a", "1")))
		fields := log.Fields()
		fields[0] = String("modified", "value")
		// Original should be unchanged
		assertEqual(t, "a", log.Fields()[0].Key())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var log *Logger
		assertNil(t, log.Fields())
	})
}

// =============================================================================
// Logger Hooks Tests
// =============================================================================

func TestLogger_Hooks(t *testing.T) {
	t.Parallel()

	t.Run("returns hooks", func(t *testing.T) {
		t.Parallel()
		log := New()
		assertNotNil(t, log.Hooks())
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Parallel()
		var log *Logger
		assertNil(t, log.Hooks())
	})
}

func TestLogger_AddHook(t *testing.T) {
	t.Parallel()

	hook := &testHook{levels: []Level{InfoLevel}}
	log := New()
	log.AddHook(hook)
	assertEqual(t, 1, log.Hooks().Len(InfoLevel))
}

// =============================================================================
// Logger IsLevelEnabled Tests
// =============================================================================

func TestLogger_IsLevelEnabled(t *testing.T) {
	t.Parallel()

	log := New(WithLevel(WarnLevel))

	t.Run("below minimum", func(t *testing.T) {
		t.Parallel()
		assertFalse(t, log.IsLevelEnabled(DebugLevel))
		assertFalse(t, log.IsLevelEnabled(InfoLevel))
	})

	t.Run("at minimum", func(t *testing.T) {
		t.Parallel()
		assertTrue(t, log.IsLevelEnabled(WarnLevel))
	})

	t.Run("above minimum", func(t *testing.T) {
		t.Parallel()
		assertTrue(t, log.IsLevelEnabled(ErrorLevel))
		assertTrue(t, log.IsLevelEnabled(FatalLevel))
	})
}

// =============================================================================
// Logger Logging Methods Tests
// =============================================================================

func TestLogger_Trace(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(TraceLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
	)
	log.Trace("trace message", String("key", "value"))

	out := buf.String()
	assertContains(t, out, "TRACE")
	assertContains(t, out, "trace message")
	assertContains(t, out, "key=value")
}

func TestLogger_Debug(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(DebugLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
	)
	log.Debug("debug message")

	assertContains(t, buf.String(), "DEBUG")
	assertContains(t, buf.String(), "debug message")
}

func TestLogger_Info(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(InfoLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
	)
	log.Info("info message")

	assertContains(t, buf.String(), "INFO")
	assertContains(t, buf.String(), "info message")
}

func TestLogger_Warn(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(WarnLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
	)
	log.Warn("warn message")

	assertContains(t, buf.String(), "WARN")
	assertContains(t, buf.String(), "warn message")
}

func TestLogger_Error(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(ErrorLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
	)
	log.Error("error message")

	assertContains(t, buf.String(), "ERROR")
	assertContains(t, buf.String(), "error message")
}

func TestLogger_Panic(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(PanicLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
	)

	defer func() {
		r := recover()
		assertNotNil(t, r)
		assertEqual(t, "panic message", r)
		assertContains(t, buf.String(), "PANIC")
	}()

	log.Panic("panic message")
}

// =============================================================================
// Logger Formatted Logging Methods Tests
// =============================================================================

func TestLogger_Tracef(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(TraceLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
	)
	log.Tracef("trace %s %d", "test", 42)

	assertContains(t, buf.String(), "trace test 42")
}

func TestLogger_Debugf(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(DebugLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
	)
	log.Debugf("debug %s", "test")

	assertContains(t, buf.String(), "debug test")
}

func TestLogger_Infof(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(InfoLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
	)
	log.Infof("info %v", map[string]int{"a": 1})

	assertContains(t, buf.String(), "info map[a:1]")
}

func TestLogger_Warnf(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(WarnLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
	)
	log.Warnf("warn %d%%", 50)

	assertContains(t, buf.String(), "warn 50%")
}

func TestLogger_Errorf(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(ErrorLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
	)
	log.Errorf("error: %v", fmt.Errorf("test error"))

	assertContains(t, buf.String(), "error: test error")
}

func TestLogger_Panicf(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(PanicLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
	)

	defer func() {
		r := recover()
		assertNotNil(t, r)
		assertEqual(t, "panic 42", r)
	}()

	log.Panicf("panic %d", 42)
}

// =============================================================================
// Logger Fluent API Tests
// =============================================================================

func TestLogger_WithLevel(t *testing.T) {
	t.Parallel()

	log := NewLogger().WithLevel(DebugLevel)
	assertEqual(t, DebugLevel, log.Level())
}

func TestLogger_WithFormatter(t *testing.T) {
	t.Parallel()

	formatter := NewJSONFormatter()
	log := NewLogger().WithFormatter(formatter)
	assertEqual(t, formatter, log.Formatter())
}

func TestLogger_WithOutput(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := NewLogger().WithOutput(&buf)
	assertEqual(t, &buf, log.Output())
}

func TestLogger_WithCaller(t *testing.T) {
	t.Parallel()

	t.Run("enabled", func(t *testing.T) {
		t.Parallel()
		log := NewLogger().WithCaller(true)
		assertTrue(t, log.IsCaller())
	})

	t.Run("disabled", func(t *testing.T) {
		t.Parallel()
		log := NewLogger().WithCaller(false)
		assertFalse(t, log.IsCaller())
	})
}

func TestLogger_WithCallerSkip(t *testing.T) {
	t.Parallel()

	log := NewLogger().WithCallerSkip(5)
	assertEqual(t, 5, log.CallerSkip())
}

func TestLogger_WithSampling(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := NewLogger().
		WithOutput(&buf).
		WithLevel(TraceLevel).
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()).
		WithSampling(SamplingOptions{first: 1, thereafter: 0, period: time.Hour})

	// First message should be logged
	log.Info("test sampling")
	assertContains(t, buf.String(), "test sampling")

	// Second identical message should be dropped
	buf.Reset()
	log.Info("test sampling")
	assertEqual(t, "", buf.String())

	// Different message should still be logged
	log.Info("different message")
	assertContains(t, buf.String(), "different message")
}

func TestLogger_FluentChaining(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := NewLogger().
		WithLevel(DebugLevel).
		WithOutput(&buf).
		WithFormatter(NewJSONFormatter()).
		WithCaller(true).
		WithCallerSkip(1)

	assertEqual(t, DebugLevel, log.Level())
	assertEqual(t, &buf, log.Output())
	assertTrue(t, log.IsCaller())
	assertEqual(t, 1, log.CallerSkip())
}

// =============================================================================
// Logger WithContext Tests
// =============================================================================

func TestLogger_WithContext(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(InfoLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
	)

	ctx := context.Background()
	ctx = WithContextFields(ctx, String("request_id", "abc123"))

	entry := log.WithContext(ctx)
	entry.Info("test message")

	assertContains(t, buf.String(), "request_id=abc123")
}

// =============================================================================
// Logger Level Filtering Tests
// =============================================================================

func TestLogger_LevelFiltering(t *testing.T) {
	t.Parallel()

	t.Run("below minimum not logged", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		log := New(
			WithOutput(&buf),
			WithLevel(WarnLevel),
			WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
		)

		log.Debug("should not appear")
		log.Info("should not appear")

		assertNotContains(t, buf.String(), "should not appear")
	})

	t.Run("at and above minimum logged", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		log := New(
			WithOutput(&buf),
			WithLevel(WarnLevel),
			WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
		)

		log.Warn("warning message")
		log.Error("error message")

		assertContains(t, buf.String(), "warning message")
		assertContains(t, buf.String(), "error message")
	})
}

// =============================================================================
// Logger Concurrent Access Tests
// =============================================================================

func TestLogger_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(TraceLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
	)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func(n int) {
			defer wg.Done()
			log.Info(fmt.Sprintf("message-%d", n))
		}(i)
		go func() {
			defer wg.Done()
			_ = log.Level()
		}()
	}
	wg.Wait()

	// Count messages
	lines := strings.Count(buf.String(), "\n")
	assertGreaterOrEqual(t, lines, 50) // At least some messages should be logged
}

func TestLogger_ConcurrentWith(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(TraceLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
	)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			child := log.With(Int("n", n))
			child.Info("child message")
		}(i)
	}
	wg.Wait()

	// Should have 50 messages
	lines := strings.Count(buf.String(), "\n")
	assertEqual(t, 50, lines)
}

// =============================================================================
// Logger Edge Cases Tests
// =============================================================================

func TestLogger_EmptyMessage(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(InfoLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
	)
	log.Info("")

	assertContains(t, buf.String(), "INFO")
}

func TestLogger_VeryLongMessage(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(InfoLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
	)
	longMsg := strings.Repeat("a", 100000)
	log.Info(longMsg)

	assertContains(t, buf.String(), longMsg)
}

func TestLogger_UnicodeMessage(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(InfoLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
	)
	log.Info("日本語メッセージ🎉")

	assertContains(t, buf.String(), "日本語メッセージ🎉")
}

func TestLogger_SpecialCharacters(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(InfoLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
	)
	log.Info("line1\nline2\ttab")

	// Should be quoted or escaped
	assertContains(t, buf.String(), "line1")
}

// =============================================================================
// Logger Caller Info Tests
// =============================================================================

func TestLogger_CallerInfo(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(InfoLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor().WithEnableCaller()),
		WithCaller(true),
	)
	log.Info("test message")

	assertContains(t, buf.String(), "caller=")
	assertContains(t, buf.String(), ".go:")
}

// =============================================================================
// testHook helper
// =============================================================================

type testHook struct {
	mu     sync.Mutex
	levels []Level
	fired  []*Entry
}

func (h *testHook) Levels() []Level {
	return h.levels
}

func (h *testHook) Fire(e *Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.fired = append(h.fired, e)
	return nil
}

// =============================================================================
// Benchmarks
// =============================================================================

func BenchmarkLogger_Info(b *testing.B) {
	log := New(WithOutput(io.Discard), WithLevel(InfoLevel))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info("benchmark message")
	}
}

func BenchmarkLogger_InfoWithFields(b *testing.B) {
	log := New(WithOutput(io.Discard), WithLevel(InfoLevel))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info("benchmark message", String("key", "value"), Int("count", 42))
	}
}

func BenchmarkLogger_With(b *testing.B) {
	log := New(WithOutput(io.Discard), WithLevel(InfoLevel))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = log.With(String("key", "value"))
	}
}

func BenchmarkLogger_DisabledLevel(b *testing.B) {
	log := New(WithOutput(io.Discard), WithLevel(ErrorLevel))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Debug("this should not be logged")
	}
}

func BenchmarkLogger_Parallel(b *testing.B) {
	log := New(WithOutput(io.Discard), WithLevel(InfoLevel))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Info("parallel benchmark")
		}
	})
}
