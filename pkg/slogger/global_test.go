package slogger

import (
	"bytes"
	"io"
	"testing"
)

// =============================================================================
// SetGlobalLogger Tests
// =============================================================================

func TestSetGlobalLogger(t *testing.T) {
	// Not parallel because it modifies global state
	defer func() {
		// Restore default global logger
		SetGlobalLogger(New())
	}()

	t.Run("sets global logger", func(t *testing.T) {
		var buf bytes.Buffer
		log := New(
			WithOutput(&buf),
			WithLevel(InfoLevel),
			WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
		)
		SetGlobalLogger(log)

		Info("test message")
		assertContains(t, buf.String(), "test message")
	})

	t.Run("nil logger ignored", func(t *testing.T) {
		originalLogger := GlobalLogger()
		SetGlobalLogger(nil)
		assertEqual(t, originalLogger, GlobalLogger())
	})
}

// =============================================================================
// GlobalLogger Tests
// =============================================================================

func TestGlobalLogger(t *testing.T) {
	t.Run("returns logger", func(t *testing.T) {
		log := GlobalLogger()
		assertNotNil(t, log)
	})
}

// =============================================================================
// Global Logging Functions Tests
// =============================================================================

func TestGlobalTrace(t *testing.T) {
	// Not parallel because it modifies global state
	defer func() {
		SetGlobalLogger(New())
	}()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(TraceLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
	)
	SetGlobalLogger(log)

	Trace("trace message", String("key", "value"))

	out := buf.String()
	assertContains(t, out, "TRACE")
	assertContains(t, out, "trace message")
	assertContains(t, out, "key=value")
}

func TestGlobalDebug(t *testing.T) {
	defer func() {
		SetGlobalLogger(New())
	}()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(DebugLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
	)
	SetGlobalLogger(log)

	Debug("debug message")

	assertContains(t, buf.String(), "DEBUG")
	assertContains(t, buf.String(), "debug message")
}

func TestGlobalInfo(t *testing.T) {
	defer func() {
		SetGlobalLogger(New())
	}()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(InfoLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
	)
	SetGlobalLogger(log)

	Info("info message")

	assertContains(t, buf.String(), "INFO")
	assertContains(t, buf.String(), "info message")
}

func TestGlobalWarn(t *testing.T) {
	defer func() {
		SetGlobalLogger(New())
	}()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(WarnLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
	)
	SetGlobalLogger(log)

	Warn("warn message")

	assertContains(t, buf.String(), "WARN")
	assertContains(t, buf.String(), "warn message")
}

func TestGlobalError(t *testing.T) {
	defer func() {
		SetGlobalLogger(New())
	}()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(ErrorLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
	)
	SetGlobalLogger(log)

	Error("error message")

	assertContains(t, buf.String(), "ERROR")
	assertContains(t, buf.String(), "error message")
}

// =============================================================================
// Global Formatted Logging Functions Tests
// =============================================================================

func TestGlobalTracef(t *testing.T) {
	defer func() {
		SetGlobalLogger(New())
	}()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(TraceLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
	)
	SetGlobalLogger(log)

	Tracef("trace %s %d", "test", 42)

	assertContains(t, buf.String(), "trace test 42")
}

func TestGlobalDebugf(t *testing.T) {
	defer func() {
		SetGlobalLogger(New())
	}()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(DebugLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
	)
	SetGlobalLogger(log)

	Debugf("debug %s", "formatted")

	assertContains(t, buf.String(), "debug formatted")
}

func TestGlobalInfof(t *testing.T) {
	defer func() {
		SetGlobalLogger(New())
	}()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(InfoLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
	)
	SetGlobalLogger(log)

	Infof("info %d", 123)

	assertContains(t, buf.String(), "info 123")
}

func TestGlobalWarnf(t *testing.T) {
	defer func() {
		SetGlobalLogger(New())
	}()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(WarnLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
	)
	SetGlobalLogger(log)

	Warnf("warn %v", map[string]int{"a": 1})

	assertContains(t, buf.String(), "warn map[a:1]")
}

func TestGlobalErrorf(t *testing.T) {
	defer func() {
		SetGlobalLogger(New())
	}()

	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(ErrorLevel),
		WithFormatter(NewTextFormatter(&buf).WithDisableColor()),
	)
	SetGlobalLogger(log)

	Errorf("error: %s", "test error")

	assertContains(t, buf.String(), "error: test error")
}

// =============================================================================
// ApplyGlobalConfig Tests
// =============================================================================

func TestApplyGlobalConfig(t *testing.T) {
	defer func() {
		SetGlobalLogger(New())
	}()

	t.Run("applies config with default level", func(t *testing.T) {
		cfg := SloggerConfig{
			Level:     "invalid", // Should fall back to InfoLevel
			Formatter: "text",
		}
		err := ApplyGlobalConfig(cfg)
		assertNoError(t, err)
		assertEqual(t, InfoLevel, GlobalLogger().Level())
	})

	t.Run("applies debug level", func(t *testing.T) {
		cfg := SloggerConfig{
			Level:     "debug",
			Formatter: "text",
		}
		err := ApplyGlobalConfig(cfg)
		assertNoError(t, err)
		assertEqual(t, DebugLevel, GlobalLogger().Level())
	})

	t.Run("applies json formatter", func(t *testing.T) {
		cfg := SloggerConfig{
			Level:     "info",
			Formatter: "json",
		}
		err := ApplyGlobalConfig(cfg)
		assertNoError(t, err)
		_, isJSON := GlobalLogger().Formatter().(*JSONFormatter)
		assertTrue(t, isJSON)
	})

	t.Run("applies text formatter", func(t *testing.T) {
		cfg := SloggerConfig{
			Level:     "info",
			Formatter: "text",
		}
		err := ApplyGlobalConfig(cfg)
		assertNoError(t, err)
		_, isText := GlobalLogger().Formatter().(*TextFormatter)
		assertTrue(t, isText)
	})

	t.Run("applies caller config", func(t *testing.T) {
		cfg := SloggerConfig{
			Level:     "info",
			Formatter: "text",
			Caller:    CallerConfig{IsEnabled: true},
		}
		err := ApplyGlobalConfig(cfg)
		assertNoError(t, err)
		assertTrue(t, GlobalLogger().IsCaller())
	})

	t.Run("applies color config", func(t *testing.T) {
		cfg := SloggerConfig{
			Level:     "info",
			Formatter: "json",
			Color:     ColorConfig{IsEnabled: true},
		}
		err := ApplyGlobalConfig(cfg)
		assertNoError(t, err)
		jsonFmt, _ := GlobalLogger().Formatter().(*JSONFormatter)
		assertTrue(t, jsonFmt.IsEnableColor())
	})
}

// =============================================================================
// Benchmarks
// =============================================================================

func BenchmarkGlobalLogger(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GlobalLogger()
	}
}

func BenchmarkGlobalInfo(b *testing.B) {
	SetGlobalLogger(New(WithOutput(io.Discard)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info("benchmark message")
	}
}

func BenchmarkGlobalInfoWithFields(b *testing.B) {
	SetGlobalLogger(New(WithOutput(io.Discard)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info("benchmark message", String("key", "value"), Int("count", 42))
	}
}
